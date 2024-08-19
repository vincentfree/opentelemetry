package providerconfig

import (
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type SignalProcessor interface {
	AsyncTraceProcessor(...trace.BatchSpanProcessorOption) trace.SpanProcessor
	SyncTraceProcessor() trace.SpanProcessor
	AsyncLogProcessor(...log.BatchProcessorOption) log.Processor
	SyncLogProcessor(...log.SimpleProcessorOption) log.Processor
	MetricProcessor(...metric.PeriodicReaderOption) metric.Reader
}
