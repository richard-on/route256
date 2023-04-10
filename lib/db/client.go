package db

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/db/metrics"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/tracing/utils"
)

type DB struct {
	pool   *pgxpool.Pool
	tracer opentracing.Tracer
}

// NewDBClient creates a new DB client with basic traces and metrics.
func NewDBClient(pool *pgxpool.Pool, tracer opentracing.Tracer) *DB {
	return &DB{
		pool:   pool,
		tracer: tracer,
	}
}

const (
	ComponentDB = "db"
	Type        = "sql"
	Instance    = "postgres"

	OperationQuery   = "query"
	OperationExec    = "exec"
	OperationBegin   = "begin"
	OperationBeginTx = "beginTx"
)

func (db *DB) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	metrics.RequestCounter.WithLabelValues(OperationQuery, query).Inc()

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}

	span := db.tracer.StartSpan(
		OperationQuery,
		opentracing.ChildOf(parentCtx),
		opentracing.Tag{Key: string(ext.Component), Value: ComponentDB},
		opentracing.Tag{Key: string(ext.DBType), Value: Type},
		opentracing.Tag{Key: string(ext.DBInstance), Value: Instance},
		opentracing.Tag{Key: string(ext.DBStatement), Value: query},
	)
	defer span.Finish()

	ctx = utils.InjectSpanContext(ctx, opentracing.GlobalTracer(), span)

	now := time.Now()
	rows, err := db.pool.Query(ctx, query, args...)
	elapsed := time.Since(now)

	if err != nil {
		otgrpc.SetSpanTags(span, err, true)
		span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
	}

	errStr := "nil"
	if err != nil {
		errStr = err.Error()
	}

	metrics.ResponseCounter.WithLabelValues(OperationQuery, query, parseTag(rows.CommandTag()), errStr).Inc()
	metrics.HistogramResponseTime.
		WithLabelValues(OperationQuery, query, parseTag(rows.CommandTag()), errStr).Observe(elapsed.Seconds())

	return rows, err
}

func (db *DB) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	metrics.RequestCounter.WithLabelValues(OperationExec, sql).Inc()

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}

	span := db.tracer.StartSpan(
		OperationExec,
		opentracing.ChildOf(parentCtx),
		opentracing.Tag{Key: string(ext.Component), Value: ComponentDB},
		opentracing.Tag{Key: string(ext.DBType), Value: Type},
		opentracing.Tag{Key: string(ext.DBInstance), Value: Instance},
		opentracing.Tag{Key: string(ext.DBStatement), Value: sql},
	)
	defer span.Finish()

	ctx = utils.InjectSpanContext(ctx, opentracing.GlobalTracer(), span)

	now := time.Now()
	tag, err := db.pool.Exec(ctx, sql, arguments...)
	elapsed := time.Since(now)

	if err != nil {
		otgrpc.SetSpanTags(span, err, true)
		span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
	}

	errStr := "nil"
	if err != nil {
		errStr = err.Error()
	}

	metrics.ResponseCounter.WithLabelValues(OperationExec, sql, parseTag(tag), errStr).Inc()
	metrics.HistogramResponseTime.
		WithLabelValues(OperationExec, sql, parseTag(tag), errStr).Observe(elapsed.Seconds())

	return tag, err
}

func (db *DB) Begin(ctx context.Context) (pgx.Tx, error) {
	metrics.TXCounter.WithLabelValues(OperationBegin).Inc()

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}

	span := db.tracer.StartSpan(
		OperationBegin,
		opentracing.ChildOf(parentCtx),
		opentracing.Tag{Key: string(ext.Component), Value: ComponentDB},
		opentracing.Tag{Key: string(ext.DBType), Value: Type},
		opentracing.Tag{Key: string(ext.DBInstance), Value: Instance},
	)
	defer span.Finish()

	tx, err := db.pool.Begin(ctx)

	if err != nil {
		otgrpc.SetSpanTags(span, err, true)
		span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
	}

	return tx, err
}

func (db *DB) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	metrics.TXCounter.WithLabelValues(OperationBeginTx).Inc()

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}

	span := db.tracer.StartSpan(
		OperationBeginTx,
		opentracing.ChildOf(parentCtx),
		opentracing.Tag{Key: string(ext.Component), Value: ComponentDB},
		opentracing.Tag{Key: string(ext.DBType), Value: Type},
		opentracing.Tag{Key: string(ext.DBInstance), Value: Instance},
	)
	defer span.Finish()

	tx, err := db.pool.BeginTx(ctx, txOptions)

	if err != nil {
		otgrpc.SetSpanTags(span, err, true)
		span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
	}

	return &Tx{tx, db.tracer}, err
}

func parseTag(tag pgconn.CommandTag) string {
	var op string
	switch {
	case tag.Update():
		op = "update"
	case tag.Delete():
		op = "delete"
	case tag.Select():
		op = "select"
	case tag.Insert():
		op = "insert"
	default:
		op = "unspecified"
	}

	return op
}
