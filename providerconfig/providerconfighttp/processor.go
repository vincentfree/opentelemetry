package providerconfighttp

import (
	"context"
	"github.com/vincentfree/opentelemetry/providerconfig"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"log/slog"
	"os"
)

var (
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelError}))
)

func New(options ...Option) providerconfig.SignalProcessor {
	ctx := context.Background()
	cfg := &httpConfig{}

	for _, opt := range options {
		opt(cfg)
	}

	traceExporter, err := otlptracehttp.New(ctx, cfg.traceOptions...)
	handleErr(err, "failed to create trace exporter")
	metricExporter, err := otlpmetrichttp.New(ctx, cfg.metricOptions...)
	handleErr(err, "failed to create metric exporter")
	logExporter, err := otlploghttp.New(ctx, cfg.logOptions...)
	handleErr(err, "failed to create log exporter")

	return &httpProvider{
		traceExporter:          traceExporter,
		metricExporter:         metricExporter,
		logExporter:            logExporter,
		batchProcessorOptions:  cfg.batchProcessorOptions,
		simpleProcessorOptions: cfg.simpleProcessorOptions,
		periodicReaderOptions:  cfg.periodicReaderOptions,
		spanProcessorOptions:   cfg.spanProcessorOptions,
	}
}

type httpProvider struct {
	traceExporter          *otlptrace.Exporter
	metricExporter         *otlpmetrichttp.Exporter
	logExporter            *otlploghttp.Exporter
	batchProcessorOptions  []log.BatchProcessorOption
	simpleProcessorOptions []log.SimpleProcessorOption
	periodicReaderOptions  []metric.PeriodicReaderOption
	spanProcessorOptions   []trace.BatchSpanProcessorOption
}

func (g httpProvider) AsyncTraceProcessor(option ...trace.BatchSpanProcessorOption) trace.SpanProcessor {
	opts := append(g.spanProcessorOptions, option...)
	return trace.NewBatchSpanProcessor(g.traceExporter, opts...)
}

func (g httpProvider) SyncTraceProcessor() trace.SpanProcessor {
	return trace.NewSimpleSpanProcessor(g.traceExporter)
}

func (g httpProvider) AsyncLogProcessor(option ...log.BatchProcessorOption) log.Processor {
	opts := append(g.batchProcessorOptions, option...)
	return log.NewBatchProcessor(g.logExporter, opts...)
}

func (g httpProvider) SyncLogProcessor(option ...log.SimpleProcessorOption) log.Processor {
	opts := append(g.simpleProcessorOptions, option...)
	return log.NewSimpleProcessor(g.logExporter, opts...)
}

func (g httpProvider) MetricProcessor(option ...metric.PeriodicReaderOption) metric.Reader {
	opts := append(g.periodicReaderOptions, option...)
	return metric.NewPeriodicReader(g.metricExporter, opts...)
}

func handleErr(err error, message string) {
	if err != nil {
		logger.Error(message, slog.Any("error", err))
		os.Exit(1)
	}
}
