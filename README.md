# OpenTelemetry http

[![Go](https://github.com/vincentfree/opentelemetry-http/actions/workflows/go.yml/badge.svg)](https://github.com/vincentfree/opentelemetry-http/actions/workflows/go.yml)
[![CodeQL](https://github.com/vincentfree/opentelemetry-http/actions/workflows/codeql.yml/badge.svg)](https://github.com/vincentfree/opentelemetry-http/actions/workflows/codeql.yml)
[![Dependency Review](https://github.com/vincentfree/opentelemetry-http/actions/workflows/dependency-review.yml/badge.svg)](https://github.com/vincentfree/opentelemetry-http/actions/workflows/dependency-review.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/vincentfree/opentelemetry-http/otelmiddleware.svg)](https://pkg.go.dev/github.com/vincentfree/opentelemetry-http/otelmiddleware)

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
