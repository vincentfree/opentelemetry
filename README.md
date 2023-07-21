# OpenTelemetry Helper functions

[![Go](https://github.com/vincentfree/opentelemetry/actions/workflows/go.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/go.yml)
[![CodeQL](https://github.com/vincentfree/opentelemetry/actions/workflows/codeql.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/codeql.yml)
[![Dependency Review](https://github.com/vincentfree/opentelemetry/actions/workflows/dependency-review.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/dependency-review.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/vincentfree/opentelemetry/otelmiddleware.svg)](https://pkg.go.dev/github.com/vincentfree/opentelemetry/otelmiddleware)

Open Telemetry http middleware. This package provides an instrumentation for middleware that can be used to trace HTTP requests.

## otelmiddleware

Package otelmiddleware provides middleware for wrapping `http.Server` handlers with Open Telemetry tracing support.

The `trace.Span` is decorated with standard metadata extracted from the `http.Request` injected into the middleware. the basic information is extracted using the OpenTelemetry semconv package.

When a span gets initialized it uses the following slice of `trace.SpanStartOption`

```go
opts := []trace.SpanStartOption{
    trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
    trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
    trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(r.Host, extractRoute(r.RequestURI), r)...),
    trace.WithAttributes(semconv.HTTPClientAttributesFromHTTPRequest(r)...),
    trace.WithAttributes(semconv.TelemetrySDKLanguageGo),
    trace.WithSpanKind(trace.SpanKindClient),
}
```

The slice can be extended using the `WithAttributes` `TraceOption` function.

After these options are applied a new span is created and the middleware will pass the `http.ResponseWriter` and `http.Request` to the next `http.Handler`.

### Functions

```go
func TraceWithOptions(opt ...TraceOption) func(next http.Handler) http.Handler
func Trace(next http.Handler) http.Handler
func WithAttributes(attributes ...attribute.KeyValue) TraceOption
func WithPropagator(p propagation.TextMapPropagator) TraceOption
func WithServiceName(serviceName string) TraceOption
func WithTracer(tracer trace.Tracer) TraceOption
```

### Types

```go
type TraceOption func(*traceConfig)
```

### Structs

```go
type traceConfig struct {
    tracer trace.Tracer
    propagator propagation.TextMapPropagator
    attributes []attribute.KeyValue
    serviceName string
}
```

## otelzerolog

Package `otelzerolog` provides a function to extend structured logs using zerolog with the Open Telemetry trace related context.

The `github.com/rs/zerolog` `zerolog.Event` is decorated with standard metadata extracted from the `trace.SpanContext`, a traceID, spanID and additional information is injected into a log.

The initialization uses file level configuration to set defaults for the function to use. `SetLogOptions` can overwrite the defaults.

When the configuration is done `AddTracingContext` and `AddTracingContextWithAttributes` decorate `zerolog` logs with data from the trace context.
A `zeroLog.Event` can be passed by using `log.Info().Func(AddTracingContext(span)).Msg("")` for example.

### Functions

```go
func SetLogOptions(options ...LogOption)
func WithTraceID(traceID string) LogOption
func WithSpanID(spanID string) LogOption
func WithServiceName(serviceName string) LogOption
func WithAttributePrefix(prefix string) LogOption
func WithAttributes(attributes ...attribute.KeyValue) LogOption
func AddTracingContext(span trace.Span, err ...error) func(event *zerolog.Event)
func AddTracingContextWithAttributes(span trace.Span, attributes []attribute.KeyValue, err ...error) func(event *zerolog.Event)
```

### Types

```go
type LogOption func(*logConfig)
```

### Structs


```go
type logConfig struct {
	attributes      []attribute.KeyValue
	serviceName     string
	traceId         string
	spanId          string
	attributePrefix string
}
```
