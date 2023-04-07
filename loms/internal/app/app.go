// Package app initializes server, establishes connection with database and other dependencies.
// It then runs the service, accepting incoming connections and listening for shutdown signals.
package app

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	jaegerconf "github.com/uber/jaeger-client-go/config"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/db"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/grpc/server/metrics"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/grpc/server/tracing"
	logger "gitlab.ozon.dev/rragusskiy/homework-1/lib/logger/grpc"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/logger/zerolog"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/config"
	lomsservice "gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/api/loms"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/message/broker/kafka"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/message/sender"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/transactor"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Run creates and runs the service using provided config.
func Run(cfg *config.Config) {
	log := zerolog.New(
		os.Stdout,
		cfg.Log.Level,
		cfg.Service.Name,
	)
	log.Info("config and logger init success")

	jaegerConfig, err := jaegerconf.FromEnv()
	if err != nil {
		return
	}
	jaegerConfig.Sampler = &jaegerconf.SamplerConfig{
		Type:  cfg.SamplerType,
		Param: cfg.SamplerParam,
	}

	closer, err := jaegerConfig.InitGlobalTracer(cfg.Service.Name)
	if err != nil {
		log.Fatal(err, "cannot init tracing")
	}
	defer func(closer io.Closer) {
		err = closer.Close()
		if err != nil {
			log.Fatal(err, "cannot close tracing")
		}
	}(closer)

	listener, err := net.Listen("tcp", net.JoinHostPort("", cfg.GRPC.Port))
	if err != nil {
		log.Fatal(err, "error while creating listener")
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			logger.UnaryServerInterceptor(
				zerolog.New(
					os.Stdout,
					cfg.Log.Level,
					fmt.Sprintf("%v-grpc", cfg.Service.Name),
				),
			),
			metrics.UnaryServerInterceptor(),
			tracing.UnaryServerInterceptor(opentracing.GlobalTracer()),
		)),
	)

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	pgConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DB))
	if err != nil {
		log.Fatal(err, "parsing database config")
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgConfig)
	if err != nil {
		log.Fatal(err, "connecting to database")
	}
	defer pool.Close()

	dbClient := db.NewDBClient(pool, opentracing.GlobalTracer())
	tx := transactor.New(dbClient)
	repo := repository.New(tx, tx, zerolog.New(
		os.Stdout,
		cfg.Log.Level,
		fmt.Sprintf("%v-postgres", cfg.Service.Name),
	))

	producer, err := kafka.NewAsyncProducer(cfg.Kafka)
	if err != nil {
		log.Fatal(err, "creating kafka producer")
	}

	// We'll need a separate context for sender.
	// This context is not cancelled on signal. It is cancelled when sender has been gracefully closed.
	senderCtx, cancelSender := context.WithCancel(context.Background())

	// Sender results will be passed into these channels.
	successChan := make(chan int64)
	failChan := make(chan int64)
	defer func() {
		close(successChan)
		close(failChan)
	}()
	statusSender := sender.NewStatusSender(
		producer,
		cfg.Kafka,
		sender.WithSuccessFunc(func(id int64) {
			successChan <- id
		}),
		sender.WithFailFunc(func(id int64) {
			failChan <- id
		}),
	)

	model := domain.New(cfg.Service, repo, tx, statusSender)

	grpchealth.RegisterHealthServer(s, health.NewServer())
	reflection.Register(s)
	loms.RegisterLOMSServer(s, lomsservice.New(model))
	go func() {
		err = s.Serve(listener)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatal(err, "error while running grpc server")
		}
	}()
	log.Infof("grpc server listening at %v", listener.Addr())

	http.Handle("/metrics", metrics.New())
	go func() {
		err = http.ListenAndServe(net.JoinHostPort("", cfg.Metrics.Port), nil)
		if err != nil {
			log.Fatal(err, "error while listening http")
		}
	}()

	// Start a separate goroutine to check and cancel unpaid orders.
	unpaidErrChan := make(chan error)
	go model.MonitorUnpaid(ctx, unpaidErrChan)
	go func() {
		for err = range unpaidErrChan {
			log.Error(err, "cancelling unpaid orders")
		}
	}()

	// Start a separate goroutine to send enqueued messages to a broker.
	unsentErrChan := make(chan error)
	go model.MonitorUnsent(ctx, unsentErrChan)
	go func() {
		for err = range unsentErrChan {
			log.Error(err, "sending message to a broker")
		}
	}()

	// Start a separate goroutine to monitor Success and Errors channels
	// and update outbox based on these results.
	monitorErrChan := make(chan error)
	go model.MonitorSenderResult(senderCtx, successChan, failChan, monitorErrChan)
	go func() {
		for err = range monitorErrChan {
			log.Error(err, "updating message status")
		}
	}()

	<-ctx.Done()
	// First, close the sender and drain its channels.
	statusSender.Close()
	// Then, cancel contexts.
	cancelSender()
	cancel()

	log.Info("shutting down: loms service")
	s.GracefulStop()
}
