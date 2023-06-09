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
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/config"
	checkoutservice "gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/api/checkout"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/client/grpc/kube"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/client/grpc/loms"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/client/grpc/productservice"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository/transactor"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/db"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/grpc/client/wrapper"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/grpc/server/metrics"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/grpc/server/tracing"
	logger "gitlab.ozon.dev/rragusskiy/homework-1/lib/logger/grpc"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/logger/zerolog"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	jaegerConfig, err := jaegercfg.FromEnv()
	if err != nil {
		return
	}
	jaegerConfig.Sampler = &jaegercfg.SamplerConfig{
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

	lomsConn, err := wrapper.NewClient(cfg.LOMS.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err, "error while dialing loms grpc server")
	}
	lomsClient := loms.NewClient(lomsConn)

	productConn, err := wrapper.NewClient(cfg.ProductService.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err, "error while dialing product service grpc server")
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	var productClient *productservice.Client
	// Execute if app is run inside kubernetes cluster.
	if cfg.Service.Environment == "kubernetes" {
		// Create a new Kubernetes client.
		kubeClient, err := kube.NewInClusterClient(cfg.Kubernetes, zerolog.New(
			os.Stdout,
			cfg.Log.Level,
			fmt.Sprintf("%v-kubeClient", cfg.Service.Name),
		))
		if err != nil {
			log.Fatal(err, "error while creating kubernetes client")
		}

		// Create a product service with dynamic rate limit.
		// Rate limit is changed based on the number of active replicas.
		productClient = productservice.NewClient(productConn, cfg.ProductService.Token,
			ratelimit.New(ctx, cfg.RateLimit.Rate,
				ratelimit.WithBurst(cfg.RateLimit.Burst),
				ratelimit.WithReplicas(kubeClient, cfg.Kubernetes.UpdateInterval),
			),
		)
	} else {
		// Create a product service without dynamic rate limit
		// as we don't expect horizontal scaling.
		productClient = productservice.NewClient(productConn, cfg.ProductService.Token,
			ratelimit.New(ctx, cfg.RateLimit.Rate,
				ratelimit.WithBurst(cfg.RateLimit.Burst),
			),
		)
	}

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

	model := domain.New(cfg.Service, repo, tx, lomsClient, productClient, lomsClient)

	grpchealth.RegisterHealthServer(s, health.NewServer())
	reflection.Register(s)
	checkout.RegisterCheckoutServer(s, checkoutservice.New(model))
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

	<-ctx.Done()
	cancel()

	log.Info("shutting down: checkout service")
	s.GracefulStop()
}
