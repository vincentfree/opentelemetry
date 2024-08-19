package providerconfig_test

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"providerconfig"
	"providerconfig/providerconfignoop"
)

// initializing through main and manually setting the providers.
//
// providers can also be set directly through providerconfig.Option's like:
// providerconfig.WithInitTraces()
// providerconfig.WithInitMetrics()
// providerconfig.WithInitLogs()
func ExampleNew() {
	otelConfig := providerconfig.New(
		providerconfig.WithApplicationName("example-app"),
		providerconfig.WithApplicationVersion("0.1.0"),
		providerconfig.WithExecutionType(providerconfig.Async),
		providerconfig.WithSignalProcessor(providerconfignoop.NewNoopProcessor()),
	)

	// traces
	otel.SetTracerProvider(otelConfig.Providers.TraceProvider)

	// metrics
	otel.SetMeterProvider(otelConfig.Providers.MetricProvider)

	// logs
	global.SetLoggerProvider(otelConfig.Providers.LogProvider)
}

func ExampleWithApplicationName() {
	providerconfig.New(
		providerconfig.WithApplicationName("example-app"),
	)
}

func ExampleWithApplicationVersion() {
	providerconfig.New(
		providerconfig.WithApplicationVersion("1.0.0"),
	)
}

func ExampleWithExecutionType() {
	providerconfig.New(
		providerconfig.WithExecutionType(providerconfig.Async),
	)
}

func ExampleWithSignalProcessor() {
	providerconfig.New(
		// Example processor, in a real scenario, the http or grpc processors should be used.
		// Either is a separate import that needs to be added to the modules
		providerconfig.WithSignalProcessor(providerconfignoop.NewNoopProcessor()),
	)
}

func ExampleWithCollectorProtocol() {
	providerconfig.New(
		providerconfig.WithCollectorProtocol(providerconfig.Http),
	)
}

func ExampleWithCollectorUrl() {
	providerconfig.New(
		providerconfig.WithCollectorUrl("0.0.0.0"),
	)
}

func ExampleWithPort() {
	providerconfig.New(
		providerconfig.WithPort(8080),
	)
}

func ExampleWithProtocolAndPort() {
	providerconfig.New(
		providerconfig.WithProtocolAndPort(providerconfig.Grpc, 443),
	)
}

func ExampleWithResourceOptions() {
	providerconfig.New(
		providerconfig.WithResourceOptions(resource.WithContainer(),
			resource.WithHost(),
		),
	)
}

func ExampleWithDisabledSignals() {
	providerconfig.New(
		providerconfig.WithDisabledSignals(
			false, // traces
			true,  // metrics
			true,  // logs
		),
	)
}

func ExampleWithTracePropagator() {
	providerconfig.New(
		providerconfig.WithTracePropagator(
			propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			),
		),
	)
}
