// Copyright 2024 Vincent Free
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package providerconfig

import (
	"context"
	"errors"
	"go.opentelemetry.io/contrib/bridges/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"log/slog"
	"os"
)

var (
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelError}))
)

type Option func(*config)
type Options []Option

type config struct {
	applicationName       string
	applicationVersion    string
	resourceOptions       []resource.Option
	prometheusBridge      bool
	metricInit            bool
	traceInit             bool
	logInit               bool
	disableTraces         bool
	disableMetrics        bool
	disableLogs           bool
	signalProcessor       SignalProcessor
	tracePropagator       propagation.TextMapPropagator
	executionType         Execution
	traceProviderOptions  []sdktrace.TracerProviderOption
	logProviderOptions    []sdklog.LoggerProviderOption
	periodicReaderOptions []sdkmetric.PeriodicReaderOption
	prometheusOptions     []prometheus.Option
}

func WithApplicationName(applicationName string) Option {
	return func(c *config) {
		c.applicationName = applicationName
	}
}

func WithApplicationVersion(applicationVersion string) Option {
	return func(c *config) {
		c.applicationVersion = applicationVersion
	}
}

func WithResourceOptions(resourceOptions ...resource.Option) Option {
	return func(c *config) {
		c.resourceOptions = append(c.resourceOptions, resourceOptions...)
	}
}

// WithTraceProviderOptions accepts TracerProviderOptions
func WithTraceProviderOptions(options ...sdktrace.TracerProviderOption) Option {
	return func(c *config) {
		c.traceProviderOptions = options
	}
}

// WithPeriodicReaderOptions accepts PeriodicReaderOptions
func WithPeriodicReaderOptions(options ...sdkmetric.PeriodicReaderOption) Option {
	return func(c *config) {
		c.periodicReaderOptions = options
	}
}

// WithLogProviderOptions accepts LoggerProviderOptions
func WithLogProviderOptions(options ...sdklog.LoggerProviderOption) Option {
	return func(c *config) {
		c.logProviderOptions = options
	}
}

// WithPrometheusBridge enables the Prometheus bridge in the configuration.
// The Prometheus bridge allows exporting metrics from the Prometheus instrumentation
// and forward them over OTLP to an endpoint.
// If the bridge is not enabled, prometheus metrics will not be exported over OTLP.
//
// # Accepts prometheus.Option's
//
// The Prometheus bridge is disabled by default.
func WithPrometheusBridge(options ...prometheus.Option) Option {
	return func(c *config) {
		c.prometheusBridge = true
		c.prometheusOptions = options
	}
}

// WithInitMetrics sets the global metric provider.
//
// If this function is not used, the user has to set the global provider or use it directly.
// Can be globally set using the following code:
//
//	otel.SetMeterProvider(otelConfig.MetricProvider())
func WithInitMetrics() Option {
	return func(c *config) {
		c.metricInit = true
	}
}

// WithInitTraces sets the global trace provider.
//
// If this function is not used, the user has to set the global provider or use it directly.
// Can be globally set using the following code:
//
//	otel.SetTracerProvider(otelConfig.TraceProvider())
func WithInitTraces() Option {
	return func(c *config) {
		c.traceInit = true
	}
}

// WithInitLogs sets the global log provider.
//
// If this function is not used, the user has to set the global provider or use it directly.
// Can be globally set using the following code:
//
//	global.SetLoggerProvider(otelConfig.LogProvider())
func WithInitLogs() Option {
	return func(c *config) {
		c.logInit = true
	}
}

// WithInitSignals sets all three observability signals by calling their setter functions.
//
// The setter functions would normally have to be set manually using the following lines of code:
//
//	 // traces
//		otel.SetTracerProvider(otelConfig.TraceProvider())
//
//		// metrics
//		otel.SetMeterProvider(otelConfig.MetricProvider())
//
//		// logs
//		global.SetLoggerProvider(otelConfig.LogProvider())
//
// With this Option these will be preformed for the user
func WithInitSignals() Option {
	return func(c *config) {
		c.traceInit = true
		c.metricInit = true
		c.logInit = true
	}
}

func WithExecutionType(executionType Execution) Option {
	return func(c *config) {
		c.executionType = executionType
	}
}

// WithTracePropagator overwrites the default trace propagators set by the library.
// The trace propagator(s) are used to propagate trace context across distributed systems.
// The trace context object and other metadata can be injected and extracted based on this configuration.
//
// If this function is not used, the defaults will be set.
// Default propagators are: propagation.TraceContext and propagation.Baggage
func WithTracePropagator(propagator propagation.TextMapPropagator) Option {
	return func(c *config) {
		c.tracePropagator = propagator
	}
}

func WithDisabledSignals(disableTraces, disableMetrics, disableLogs bool) Option {
	return func(c *config) {
		c.disableTraces = disableTraces
		c.disableMetrics = disableMetrics
		c.disableLogs = disableLogs
	}
}

// WithSignalProcessor expects an implementation of the SignalProcessor
// interface. There are two implementations provided by this library as separate
// modules to reduce the number of imported dependencies. This limits to what's
// actually used by the end user.
//
// The implementations can be found in the packages:
//   - github.com/vincentfree/opentelemetry/providerconfiggrpc
//   - github.com/vincentfree/opentelemetry/providerconfighttp
//
// Both packages contain a new function with their respective options.
func WithSignalProcessor(signalProcessor SignalProcessor) Option {
	return func(c *config) {
		c.signalProcessor = signalProcessor
	}
}

func initConfig(options Options) *config {
	cfg := &config{}
	for _, option := range options {
		option(cfg)
	}
	if cfg.applicationName == "" {
		panic("application name is required, use the 'providerconfig.WithApplicationName' option")
	}
	if cfg.applicationVersion == "" {
		panic("application version is required, use the 'providerconfig.WithApplicationVersion' option")
	}

	if cfg.tracePropagator == nil {
		cfg.tracePropagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	}

	if !cfg.executionType.IsValid() {
		logger.Info("using default async processors for execution of signals. Use 'providerconfig.WithExecutionType' option to override the execution type.")
		cfg.executionType = Async
	}

	return cfg
}

func newResource(applicationName, applicationVersion string, resources ...resource.Option) *resource.Resource {
	resList := make([]resource.Option, 0, len(resources))
	resList = append(resList, resources...)
	resList = append(resList, resource.WithAttributes(
		semconv.ServiceNameKey.String(applicationName),
		semconv.ServiceVersionKey.String(applicationVersion),
	))
	r, err := resource.New(context.Background(), resList...)
	if errors.Is(err, resource.ErrPartialResource) || errors.Is(err, resource.ErrSchemaURLConflict) {
		logger.Warn("partial error while building the resource used for otel providers", slog.Any("error", err))
	} else if err != nil {
		logger.Error("Unrecoverable error while building the resource used for otel providers", slog.Any("error", err))
		os.Exit(1)
	}

	return r
}

// New initializes an OTLP exporter, and configures the corresponding trace, log and
// metric providers.
//
// Although the function does not specify required Options,
// WithApplicationName and WithApplicationVersion are required.
//
// Next to Application information, WithSignalProcessor is also required.
// There are two external modules that provide implementations, one fork gRPC and one for HTTP.
//
// The implementations can be found in these packages:
//   - github.com/vincentfree/opentelemetry/providerconfiggrpc
//   - github.com/vincentfree/opentelemetry/providerconfighttp
func New(options ...Option) Provider {
	ctx := context.Background()
	cfg := initConfig(options)

	res := newResource(cfg.applicationName, cfg.applicationVersion, cfg.resourceOptions...)

	tracerProvider := setupTraceProvider(cfg, res)

	logProvider := setupLogProvider(cfg, res)

	meterProvider := setupMetricProvider(cfg, res)

	if cfg.metricInit && !cfg.disableMetrics {
		otel.SetMeterProvider(meterProvider)
	}

	if cfg.traceInit && !cfg.disableTraces {
		otel.SetTextMapPropagator(cfg.tracePropagator)
		otel.SetTracerProvider(tracerProvider)
	}

	if cfg.logInit && !cfg.disableLogs {
		global.SetLoggerProvider(logProvider)
	}

	hooks := NewShutdownHooks(
		ShutDownPair(TraceHook, traceProviderHook(ctx, tracerProvider)),
		ShutDownPair(MetricHook, metricProviderHook(ctx, meterProvider)),
		ShutDownPair(LogHook, logProviderHook(ctx, logProvider)),
	)

	return &providers{
		traceProvider:  tracerProvider,
		metricProvider: meterProvider,
		logProvider:    logProvider,
		hooks:          hooks,
	}
}

func setupMetricProvider(cfg *config, res *resource.Resource) *sdkmetric.MeterProvider {
	var metricOptions []sdkmetric.PeriodicReaderOption

	if cfg.periodicReaderOptions != nil {
		metricOptions = append(metricOptions, cfg.periodicReaderOptions...)
	}

	if cfg.prometheusBridge {
		bridge := prometheus.NewMetricProducer(cfg.prometheusOptions...)
		metricOptions = append(metricOptions, sdkmetric.WithProducer(bridge))
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(cfg.signalProcessor.MetricProcessor(metricOptions...)),
	)
	return meterProvider
}

func setupLogProvider(cfg *config, res *resource.Resource) *sdklog.LoggerProvider {
	var logProcessor sdklog.Processor
	switch cfg.executionType {
	case Async:
		logProcessor = cfg.signalProcessor.AsyncLogProcessor()
	case Sync:
		logProcessor = cfg.signalProcessor.SyncLogProcessor()
	}

	logProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(logProcessor),
	)
	return logProvider
}

func setupTraceProvider(cfg *config, res *resource.Resource) *sdktrace.TracerProvider {
	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	var bsp sdktrace.SpanProcessor
	switch cfg.executionType {
	case Sync:
		bsp = cfg.signalProcessor.SyncTraceProcessor()
	case Async:
		bsp = cfg.signalProcessor.AsyncTraceProcessor()
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	return tracerProvider
}

func handleErr(err error, message string) {
	if err != nil {
		logger.Error(message, slog.Any("error", err))
		os.Exit(1)
	}
}
func traceProviderHook(ctx context.Context, provider *sdktrace.TracerProvider) ShutdownHook {
	return func() {
		handleErr(provider.Shutdown(ctx), "failed to shutdown TracerProvider")
	}
}

func metricProviderHook(ctx context.Context, provider *sdkmetric.MeterProvider) ShutdownHook {
	return func() {
		handleErr(provider.Shutdown(ctx), "failed to shutdown MetricProvider")
	}
}

func logProviderHook(ctx context.Context, provider *sdklog.LoggerProvider) ShutdownHook {
	return func() {
		handleErr(provider.Shutdown(ctx), "failed to shutdown logProvider")
	}
}
