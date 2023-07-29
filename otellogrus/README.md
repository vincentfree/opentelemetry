# OpenTelemetry extensions - otellogrus

| Home                 | Related                                                                    |
|----------------------|----------------------------------------------------------------------------|
| [Home](../README.md) | [otelslog](../otelslog/README.md), [otelzerolog](../otelzerolog/README.md) |

----

[![Go](https://github.com/vincentfree/opentelemetry/actions/workflows/go.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/go.yml)
[![CodeQL](https://github.com/vincentfree/opentelemetry/actions/workflows/codeql.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/codeql.yml)
[![Dependency Review](https://github.com/vincentfree/opentelemetry/actions/workflows/dependency-review.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/dependency-review.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/vincentfree/opentelemetry/otelmiddleware.svg)](https://pkg.go.dev/github.com/vincentfree/opentelemetry/otellogrus)

Package `otellogrus` provides a function to extend structured logs using logrus with the Open Telemetry trace related context.

The `github.com/sirupsen/logrus` `logrus` logs are decorated with standard metadata extracted from the `trace.SpanContext`, a traceID, spanID and additional information is injected into a log.

The initialization uses file level configuration to set defaults for the function to use. `SetLogOptions` can overwrite the defaults.

When the configuration is done `AddTracingContext` and `AddTracingContextWithAttributes` decorate `logrus` logs with data from the trace context.
Adding trace context ata to logs can be achieved by using `logrus.WithFields(AddTracingContext(span)).Info("test")` for example.

### Functions

```go
func SetLogOptions(options ...LogOption)
func WithTraceID(traceID string) LogOption
func WithSpanID(spanID string) LogOption
func WithServiceName(serviceName string) LogOption
func WithAttributePrefix(prefix string) LogOption
func WithAttributes(attributes ...attribute.KeyValue) LogOption
func AddTracingContext(span trace.Span, err ...error) logrus.Fields
func AddTracingContextWithAttributes(span trace.Span, attributes []attribute.KeyValue, err ...error) logrus.Fields
func WithLevel(level logrus.Level) LoggerOption
func WithFormatter(formatter logrus.Formatter) LoggerOption
func New(options ...LoggerOption) *Logger
func (l Logger) WithTracingContext(span trace.Span, err ...error) *logrus.Entry
func (l Logger) WithTracingContextAndAttributes(span trace.Span, attributes []attribute.KeyValue, err ...error) *logrus.Entry
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
type Logger struct {
    *logrus.Logger
}

type loggerConfig struct {
    formatter *logrus.Formatter
    level     logrus.Level
}
```
