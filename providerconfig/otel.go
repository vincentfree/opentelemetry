package internal

import (
	"context"
	"fmt"
	"github.com/vincentfree/opentelemetry/otelslog"
	otelslogger "go.opentelemetry.io/contrib/bridges/otelslog"
	prommetric "go.opentelemetry.io/contrib/bridges/prometheus"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdkLog "go.opentelemetry.io/otel/sdk/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"os"
	"time"
)

const (
	grpcPort = 4137
	httpPort = 4318
)

type OtelConfig struct {
	Tracer    trace.Tracer
	Logger    *otelslog.Logger
	Providers *Providers `yaml:"providers"`
}

type Providers struct {
	TraceProvider  sdkTrace.TracerProvider
	MetricProvider sdkMetric.MeterProvider
	LogProvider    sdkLog.LoggerProvider
	Hooks          ShutdownHooks
}

func NewResource(applicationName, applicationVersion string, resources ...resource.Option) *resource.Resource {
	resList := make([]resource.Option, len(resources))
	resList = append(resList, resources...)
	resList = append(resList, resource.WithAttributes(
		semconv.ServiceNameKey.String(applicationName),
		semconv.ServiceVersionKey.String(applicationVersion),
	))
	r, _ := resource.New(context.Background(), resList...)

	return r
}

// InitProvider initializes an OTLP exporter, and configures the corresponding trace, log and
// metric providers.
func New(applicationName, version, collectorUrl string) *OtelConfig {
	ctx := context.Background()

	res := NewResource(applicationName, version)

	dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	// TODO add timeout to connection

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", collectorUrl, grpcPort), dialOption)
	handleErr(err, "failed to create gRPC connection to collector")
	// Set up a trace exporter
	traceCtx, traceCancel := context.WithTimeout(ctx, time.Second*10)
	defer traceCancel()
	logCtx, logCancel := context.WithTimeout(ctx, time.Second*10)
	defer logCancel()
	metricCtx, metricCancel := context.WithTimeout(ctx, time.Second*10)
	defer metricCancel()

	traceExporter, err := otlptracegrpc.New(traceCtx, otlptracegrpc.WithGRPCConn(conn))
	handleErr(err, "failed to create trace exporter")

	logExporterHttp, logErr := otlploghttp.New(logCtx, otlploghttp.WithInsecure(), otlploghttp.WithEndpoint(fmt.Sprintf("%s:%d", collectorUrl, httpPort)))
	handleErr(logErr, "failed to create log exporter")

	metricExporter, metricErr := otlpmetrichttp.New(metricCtx, otlpmetrichttp.WithInsecure(), otlpmetrichttp.WithEndpoint(fmt.Sprintf("%s:%d", collectorUrl, httpPort)))
	handleErr(metricErr, "failed to create metric exporter")

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdkTrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdkTrace.NewTracerProvider(
		sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
		sdkTrace.WithResource(res),
		sdkTrace.WithSpanProcessor(bsp),
	)

	logProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		//sdklog.WithProcessor(sdklog.NewSimpleProcessor(logExporter)),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporterHttp)),
	)

	global.SetLoggerProvider(logProvider)
	otelslog.SetLogOptions(otelslog.WithTraceID("traceID"))
	logger := otelslog.NewWithHandler(otelslogger.NewHandler(applicationName, otelslogger.WithLoggerProvider(logProvider)))
	slog.SetDefault(logger.Logger)

	//Metrics
	bridge := prommetric.NewMetricProducer()
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			metric.WithProducer(bridge)),
		),
	)

	otel.SetMeterProvider(meterProvider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, jaeger.Jaeger{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	tracer := otel.Tracer(applicationName)
	hooks := NewShutdownHooks(
		ShutDownPair(TraceHook, func() {
			handleErr(tracerProvider.Shutdown(ctx), "failed to shutdown TracerProvider")
		}),
		ShutDownPair(MetricHook, func() {
			handleErr(metricExporter.Shutdown(ctx), "failed to shutdown MetricExporter")
		}),
		ShutDownPair(logHook, func() {
			handleErr(logProvider.Shutdown(ctx), "failed to shutdown LogProvider")
		}),
	)
	return &OtelConfig{
		Tracer: tracer,
		Logger: logger,
		Providers: &Providers{
			TraceProvider:  sdkTrace.TracerProvider{},
			MetricProvider: sdkMetric.MeterProvider{},
			LogProvider:    sdklog.LoggerProvider{},
			Hooks:          hooks,
		},
	}
}

func handleErr(err error, message string) {
	if err != nil {
		slog.Error(message, slog.Any("error", err))
		os.Exit(1)
	}
}
