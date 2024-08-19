package providerconfighttp_test

import (
	"github.com/vincentfree/opentelemetry/providerconfig"
	"github.com/vincentfree/opentelemetry/providerconfig/providerconfighttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/trace"
)

func ExampleNew() {
	signalProcessor := providerconfighttp.New(
		providerconfighttp.WithCollectorEndpoint("0.0.0.0", 9898),
	)
	otelConfig := providerconfig.New(providerconfig.WithApplicationName("example-app"),
		providerconfig.WithApplicationVersion("0.1.0"),
		providerconfig.WithSignalProcessor(signalProcessor),
	)
	// set providers
	otel.SetTracerProvider(otelConfig.TraceProvider())
	otel.SetMeterProvider(otelConfig.MetricProvider())
	global.SetLoggerProvider(otelConfig.LogProvider())

	// shutdown all providers
	otelConfig.ShutdownAll()
}

func ExampleWithCollectorEndpoint() {
	providerconfighttp.New(
		providerconfighttp.WithCollectorEndpoint("0.0.0.0", 9898),
	)
}

func ExampleWithTraceOptions() {
	providerconfighttp.New(
		providerconfighttp.WithTraceOptions(otlptracehttp.WithCompression(otlptracehttp.GzipCompression)),
	)
}

func ExampleWithSpanProcessorOptions() {
	providerconfighttp.New(
		providerconfighttp.WithSpanProcessorOptions(trace.WithMaxExportBatchSize(512)),
	)
}

func ExampleWithLogOptions() {
	providerconfighttp.New(
		providerconfighttp.WithLogOptions(otlploghttp.WithCompression(otlploghttp.GzipCompression)),
	)
}

func ExampleWithBatchProcessorOptions() {
	providerconfighttp.New(
		providerconfighttp.WithBatchProcessorOptions(log.WithExportMaxBatchSize(512)),
	)
}
