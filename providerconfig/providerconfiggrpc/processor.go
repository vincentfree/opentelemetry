package providerconfiggrpc

import (
	"context"
	"github.com/vincentfree/opentelemetry/providerconfig"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
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
	cfg := &grpcConfig{}

	for _, opt := range options {
		opt(cfg)
	}

	traceExporter, err := otlptracegrpc.New(ctx, cfg.traceOptions...)
	handleErr(err, "failed to create trace exporter")
	metricExporter, err := otlpmetricgrpc.New(ctx, cfg.metricOptions...)
	handleErr(err, "failed to create metric exporter")
	logExporter, err := otlploggrpc.New(ctx, cfg.logOptions...)
	handleErr(err, "failed to create log exporter")

	return &grpcProvider{
		traceExporter:          traceExporter,
		metricExporter:         metricExporter,
		logExporter:            logExporter,
		batchProcessorOptions:  cfg.batchProcessorOptions,
		simpleProcessorOptions: cfg.simpleProcessorOptions,
		periodicReaderOptions:  cfg.periodicReaderOptions,
		spanProcessorOptions:   cfg.spanProcessorOptions,
	}
}

type grpcProvider struct {
	traceExporter          *otlptrace.Exporter
	metricExporter         *otlpmetricgrpc.Exporter
	logExporter            *otlploggrpc.Exporter
	batchProcessorOptions  []log.BatchProcessorOption
	simpleProcessorOptions []log.SimpleProcessorOption
	periodicReaderOptions  []metric.PeriodicReaderOption
	spanProcessorOptions   []trace.BatchSpanProcessorOption
}

func (g grpcProvider) AsyncTraceProcessor(option ...trace.BatchSpanProcessorOption) trace.SpanProcessor {
	opts := append(g.spanProcessorOptions, option...)
	return trace.NewBatchSpanProcessor(g.traceExporter, opts...)
}

func (g grpcProvider) SyncTraceProcessor() trace.SpanProcessor {
	return trace.NewSimpleSpanProcessor(g.traceExporter)
}

func (g grpcProvider) AsyncLogProcessor(option ...log.BatchProcessorOption) log.Processor {
	opts := append(g.batchProcessorOptions, option...)
	return log.NewBatchProcessor(g.logExporter, opts...)
}

func (g grpcProvider) SyncLogProcessor(option ...log.SimpleProcessorOption) log.Processor {
	opts := append(g.simpleProcessorOptions, option...)
	return log.NewSimpleProcessor(g.logExporter, opts...)
}

func (g grpcProvider) MetricProcessor(option ...metric.PeriodicReaderOption) metric.Reader {
	opts := append(g.periodicReaderOptions, option...)
	return metric.NewPeriodicReader(g.metricExporter, opts...)
}

func handleErr(err error, message string) {
	if err != nil {
		logger.Error(message, slog.Any("error", err))
		os.Exit(1)
	}
}
