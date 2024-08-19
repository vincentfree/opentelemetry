package providerconfignoop

import (
	"context"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/trace"
	"providerconfig"
)

type noopProcessor struct{}

func NewNoopProcessor() providerconfig.SignalProcessor {
	return &noopProcessor{}
}
func (n noopProcessor) AsyncTraceProcessor(option ...trace.BatchSpanProcessorOption) trace.SpanProcessor {
	return &noopSpanProcessor{}
}

func (n noopProcessor) SyncTraceProcessor() trace.SpanProcessor {
	return &noopSpanProcessor{}
}

func (n noopProcessor) AsyncLogProcessor(option ...log.BatchProcessorOption) log.Processor {
	return &noopLogProcessor{}
}

func (n noopProcessor) SyncLogProcessor(option ...log.SimpleProcessorOption) log.Processor {
	return &noopLogProcessor{}
}

// MetricProcessor can't implement the interface for metric.Reader so the metric.NewManualReader is returned
func (n noopProcessor) MetricProcessor(option ...metric.PeriodicReaderOption) metric.Reader {
	return metric.NewManualReader()
}

type noopSpanProcessor struct{}

func (n noopSpanProcessor) OnStart(_ context.Context, _ trace.ReadWriteSpan) {
	return
}

func (n noopSpanProcessor) OnEnd(_ trace.ReadOnlySpan) {
	return
}

func (n noopSpanProcessor) Shutdown(_ context.Context) error {
	return nil
}

func (n noopSpanProcessor) ForceFlush(_ context.Context) error {
	return nil
}

type noopLogProcessor struct{}

func (n noopLogProcessor) OnEmit(ctx context.Context, record log.Record) error {
	return nil
}

func (n noopLogProcessor) Enabled(ctx context.Context, record log.Record) bool {
	return true
}

func (n noopLogProcessor) Shutdown(ctx context.Context) error {
	return nil
}

func (n noopLogProcessor) ForceFlush(ctx context.Context) error {
	return nil
}

type noopMetricReader struct{}

func (n noopMetricReader) register(producer interface{}) {
}

func (n noopMetricReader) temporality(kind metric.InstrumentKind) metricdata.Temporality {
	return 0
}

func (n noopMetricReader) aggregation(kind metric.InstrumentKind) metric.Aggregation {
	return metric.AggregationDefault{}
}

func (n noopMetricReader) Collect(ctx context.Context, rm *metricdata.ResourceMetrics) error {
	return nil
}

func (n noopMetricReader) Shutdown(ctx context.Context) error {
	return nil
}
