# OpenTelemetry extensions - otelzerolog

| Home                 | Related                                                                  |
|----------------------|--------------------------------------------------------------------------|
| [Home](../README.md) | [otelslog](../otelslog/README.md), [otellogrus](../otellogrus/README.md) |

----

[![Go](https://github.com/vincentfree/opentelemetry/actions/workflows/go.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/go.yml)
[![CodeQL](https://github.com/vincentfree/opentelemetry/actions/workflows/codeql.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/codeql.yml)
[![Dependency Review](https://github.com/vincentfree/opentelemetry/actions/workflows/dependency-review.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/dependency-review.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/vincentfree/opentelemetry/otelmiddleware.svg)](https://pkg.go.dev/github.com/vincentfree/opentelemetry/otelzerolog)

Package `otelzerolog` provides a function to extend structured logs using zerolog with the Open Telemetry trace related
context.

The `github.com/rs/zerolog` `zerolog.Event` is decorated with standard metadata extracted from the `trace.SpanContext`,
a traceID, spanID and additional information is injected into a log.

The initialization uses file level configuration to set defaults for the function to use. `SetLogOptions` can overwrite
the defaults.

When the configuration is done `AddTracingContext` and `AddTracingContextWithAttributes` decorate `zerolog` logs with
data from the trace context.
A `zeroLog.Event` can be passed by using `log.Info().Func(AddTracingContext(span)).Msg("")` for example.

### Functions

```go
func SetLogOptions(options ...LogOption)
func WithTraceID(traceID string) LogOption
func WithSpanID(spanID string) LogOption
func WithServiceName(serviceName string) LogOption
func WithAttributePrefix(prefix string) LogOption
func WithAttributes(attributes ...attribute.KeyValue) LogOption
func AddTracingContext(span trace.Span, err ...error) func (event *zerolog.Event)
func AddTracingContextWithAttributes(span trace.Span, attributes []attribute.KeyValue, err ...error) func (event *zerolog.Event)
```

### Types

```go
type LogOption func (*logConfig)
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