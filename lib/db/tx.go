package db

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/db/metrics"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/tracing/utils"
)

type Tx struct {
	tx     pgx.Tx
	tracer opentracing.Tracer
}

const (
	TagUnspecified = "unspecified"

	OperationCommit   = "commit"
	OperationRollback = "rollback"

	OperationQueryRow  = "query_row"
	OperationQueryFunc = "query_func"
)

func (tx *Tx) Commit(ctx context.Context) error {
	metrics.TXCounter.WithLabelValues(OperationCommit).Inc()

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}

	span := tx.tracer.StartSpan(
		OperationCommit,
		opentracing.ChildOf(parentCtx),
		opentracing.Tag{Key: string(ext.Component), Value: ComponentDB},
		opentracing.Tag{Key: string(ext.DBType), Value: Type},
		opentracing.Tag{Key: string(ext.DBInstance), Value: Instance},
	)
	defer span.Finish()

	ctx = utils.InjectSpanContext(ctx, opentracing.GlobalTracer(), span)

	err := tx.tx.Commit(ctx)

	if err != nil {
		otgrpc.SetSpanTags(span, err, true)
		span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
	}

	return err
}

func (tx *Tx) Rollback(ctx context.Context) error {
	metrics.TXCounter.WithLabelValues(OperationRollback).Inc()

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}

	span := tx.tracer.StartSpan(
		OperationRollback,
		opentracing.ChildOf(parentCtx),
		opentracing.Tag{Key: string(ext.Component), Value: ComponentDB},
		opentracing.Tag{Key: string(ext.DBType), Value: Type},
		opentracing.Tag{Key: string(ext.DBInstance), Value: Instance},
	)
	defer span.Finish()

	ctx = utils.InjectSpanContext(ctx, opentracing.GlobalTracer(), span)

	err := tx.tx.Rollback(ctx)

	if err != nil {
		otgrpc.SetSpanTags(span, err, true)
		span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
	}

	return err
}

func (tx *Tx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	metrics.RequestCounter.WithLabelValues(OperationQuery, sql).Inc()

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}

	span := tx.tracer.StartSpan(
		OperationQuery,
		opentracing.ChildOf(parentCtx),
		opentracing.Tag{Key: string(ext.Component), Value: ComponentDB},
		opentracing.Tag{Key: string(ext.DBType), Value: Type},
		opentracing.Tag{Key: string(ext.DBInstance), Value: Instance},
		opentracing.Tag{Key: string(ext.DBStatement), Value: sql},
	)
	defer span.Finish()

	ctx = utils.InjectSpanContext(ctx, opentracing.GlobalTracer(), span)

	now := time.Now()
	rows, err := tx.tx.Query(ctx, sql, args...)
	elapsed := time.Since(now)

	if err != nil {
		otgrpc.SetSpanTags(span, err, true)
		span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
	}

	errStr := "nil"
	if err != nil {
		errStr = err.Error()
	}

	metrics.ResponseCounter.WithLabelValues(OperationQuery, sql, parseTag(rows.CommandTag()), errStr).Inc()
	metrics.HistogramResponseTime.
		WithLabelValues(OperationQuery, sql, parseTag(rows.CommandTag()), errStr).Observe(elapsed.Seconds())

	return rows, err
}

func (tx *Tx) Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error) {
	metrics.RequestCounter.WithLabelValues(OperationExec, sql).Inc()

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}

	span := tx.tracer.StartSpan(
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
	tag, err := tx.tx.Exec(ctx, sql, arguments...)
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

func (tx *Tx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	metrics.RequestCounter.WithLabelValues(OperationQueryRow, sql).Inc()

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}

	span := tx.tracer.StartSpan(
		OperationQueryRow,
		opentracing.ChildOf(parentCtx),
		opentracing.Tag{Key: string(ext.Component), Value: ComponentDB},
		opentracing.Tag{Key: string(ext.DBType), Value: Type},
		opentracing.Tag{Key: string(ext.DBInstance), Value: Instance},
		opentracing.Tag{Key: string(ext.DBStatement), Value: sql},
	)
	defer span.Finish()

	ctx = utils.InjectSpanContext(ctx, opentracing.GlobalTracer(), span)

	now := time.Now()
	row := tx.tx.QueryRow(ctx, sql, args...)
	elapsed := time.Since(now)

	errStr := "nil"

	metrics.ResponseCounter.WithLabelValues(OperationQueryRow, sql, TagUnspecified, errStr).Inc()
	metrics.HistogramResponseTime.
		WithLabelValues(OperationQueryRow, sql, TagUnspecified, errStr).Observe(elapsed.Seconds())

	return row
}

func (tx *Tx) QueryFunc(ctx context.Context, sql string,
	args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	metrics.RequestCounter.WithLabelValues(OperationQueryFunc, sql).Inc()

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}

	span := tx.tracer.StartSpan(
		OperationQueryFunc,
		opentracing.ChildOf(parentCtx),
		opentracing.Tag{Key: string(ext.Component), Value: ComponentDB},
		opentracing.Tag{Key: string(ext.DBType), Value: Type},
		opentracing.Tag{Key: string(ext.DBInstance), Value: Instance},
		opentracing.Tag{Key: string(ext.DBStatement), Value: sql},
	)
	defer span.Finish()

	ctx = utils.InjectSpanContext(ctx, opentracing.GlobalTracer(), span)

	now := time.Now()
	tag, err := tx.tx.QueryFunc(ctx, sql, args, scans, f)
	elapsed := time.Since(now)

	if err != nil {
		otgrpc.SetSpanTags(span, err, true)
		span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
	}

	errStr := "nil"
	if err != nil {
		errStr = err.Error()
	}

	metrics.ResponseCounter.WithLabelValues(OperationQueryFunc, sql, parseTag(tag), errStr).Inc()
	metrics.HistogramResponseTime.
		WithLabelValues(OperationQueryFunc, sql, parseTag(tag), errStr).Observe(elapsed.Seconds())

	return tag, err
}

func (tx *Tx) CopyFrom(ctx context.Context, tableName pgx.Identifier,
	columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	// No metrics and traces.
	return tx.tx.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

func (tx *Tx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	// No metrics and traces.
	return tx.tx.SendBatch(ctx, b)
}

func (tx *Tx) LargeObjects() pgx.LargeObjects {
	// No metrics and traces.
	return tx.tx.LargeObjects()
}

func (tx *Tx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	// No metrics and traces.
	return tx.tx.Prepare(ctx, name, sql)
}

func (tx *Tx) Conn() *pgx.Conn {
	// No metrics and traces.
	return tx.tx.Conn()
}

func (tx *Tx) Begin(ctx context.Context) (pgx.Tx, error) {
	// No metrics and traces.
	return tx.tx.Begin(ctx)
}

func (tx *Tx) BeginFunc(ctx context.Context, f func(pgx.Tx) error) (err error) {
	// No metrics and traces.
	return tx.tx.BeginFunc(ctx, f)
}
