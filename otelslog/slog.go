// Copyright 2023 Vincent Free
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

package otelslog

import (
	"context"
	otelslogger "go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

// LogOption takes a logConfig struct and applies changes.
// It can be passed to the SetLogOptions function to configure a logConfig struct.
type LogOption func(*logConfig)

type logConfig struct {
	attributes       []attribute.KeyValue
	serviceName      string
	traceId          string
	spanId           string
	attributePrefix  string
	handler          *otelslogger.Handler
	overwriteHandler slog.Handler
	bridgeDisabled   bool
	handlerOptions   *slog.HandlerOptions
}

type Logger struct {
	*slog.Logger
	// defaultTraceId has a default trace ID key in the logs
	defaultTraceId string
	// defaultSpanId has a default span ID key in the logs
	defaultSpanId string
	// defaultServiceName is empty by default, when no value is set the service name won't be used with a default in the logs
	defaultServiceName string
	// defaultAttributes contains a global set of attribute.KeyValue's that will be added to very structured log. when the slice is empty they won't be added
	defaultAttributes []attribute.KeyValue
	// defaultAttrPrefix
	defaultAttrPrefix string
}

func defaultLogger() *Logger {
	return &Logger{
		Logger:             slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		defaultTraceId:     "traceID",
		defaultSpanId:      "spanID",
		defaultServiceName: "",
		defaultAttributes:  []attribute.KeyValue(nil),
		defaultAttrPrefix:  "trace.attribute",
	}
}

var (
	_logger = defaultLogger()
)

//	// _traceId has a default trace ID key in the logs
//	_traceId = "traceID"
//	// _spanId has a default span ID key in the logs
//	_spanId = "spanID"
//	// _serviceName is empty by default, when no value is set the service name won't be used with a default in the logs
//	_serviceName string
//	// _attributes contains a global set of attribute.KeyValue's that will be added to very structured log. when the slice is empty they won't be added
//	_attributes = []attribute.KeyValue(nil)
//
//	// _attrPrefix x
//	_attrPrefix = "trace.attribute"
//)

// SetLogOptions takes LogOption's and overwrites library defaults
func SetLogOptions(options ...LogOption) {
	_logger = initLogger(options)
	slog.SetDefault(_logger.Logger)
}

func initLogger(options []LogOption) *Logger {
	// initialize an empty logConfig
	config := &logConfig{}
	logger := defaultLogger()
	for _, option := range options {
		option(config)
	}

	if config.traceId != "" {
		logger.defaultTraceId = config.traceId
	}

	if config.spanId != "" {
		logger.defaultSpanId = config.spanId
	}

	if config.serviceName != "" {
		logger.defaultServiceName = config.serviceName
	}

	if config.attributePrefix != "" {
		ap, _ := strings.CutSuffix(config.attributePrefix, ".")
		logger.defaultAttrPrefix = ap
	}

	if len(config.attributes) > 0 {
		logger.defaultAttributes = append(logger.defaultAttributes, config.attributes...)
	}

	if config.handler == nil && !config.bridgeDisabled {
		var name string
		if config.serviceName != "" {
			name = config.serviceName
		} else {
			name = "myApp"
		}

		config.handler = otelslogger.NewHandler(name, otelslogger.WithLoggerProvider(noop.NewLoggerProvider()))
	}

	if !config.bridgeDisabled {
		if config.overwriteHandler == nil {
			logger.Logger = slog.New(config.handler)
		} else {
			logger.Logger = slog.New(config.overwriteHandler)
		}

	}

	return logger
}

// WithTraceID overwrites the default 'traceID' field in the structured logs with your own key
func WithTraceID(traceID string) LogOption {
	return func(c *logConfig) {
		c.traceId = traceID
	}
}

// WithSpanID overwrites the default 'spanID' field in the structured logs with your own key
func WithSpanID(spanID string) LogOption {
	return func(c *logConfig) {
		c.spanId = spanID
	}
}

// WithServiceName adds a service name to the field 'service.name' in your structured logs
func WithServiceName(serviceName string) LogOption {
	return func(c *logConfig) {
		c.serviceName = serviceName
	}
}

// WithAttributePrefix updates the default 'trace.attribute' attribute prefix
func WithAttributePrefix(prefix string) LogOption {
	return func(c *logConfig) {
		c.attributePrefix = prefix
	}
}

// WithAttributes adds global attributes that will be added to all structured logs. attributes have a prefix followed by the key of the attribute.
//
// Example: if the attribute is of type string and the key is: 'http.method' then in the log it uses the default(but over-writable) 'trace.attribute' followed by 'http.method' so the end result is: 'trace.attribute.http.method'
func WithAttributes(attributes ...attribute.KeyValue) LogOption {
	return func(c *logConfig) {
		c.attributes = append(c.attributes, attributes...)
	}
}

func WithHandlerOptions(options *slog.HandlerOptions) LogOption {
	return func(c *logConfig) {
		c.handlerOptions = options
	}
}

// AddTracingContext lets you add the trace context to a structured log
func AddTracingContext(span trace.Span, err ...error) []slog.Attr {
	return _logger.addTracingContext(span, err...)
}

// AddTracingContextWithAttributes lets you add the trace context to a structured log, including attribute.KeyValue's to extend the log
func AddTracingContextWithAttributes(span trace.Span, attributes []attribute.KeyValue, err ...error) []slog.Attr {
	return _logger.addTraceContextWithAttributes(span, attributes, err...)

}

func (l *Logger) addTracingContext(span trace.Span, err ...error) []slog.Attr {
	a := []attribute.KeyValue(nil)
	return l.addTraceContextWithAttributes(span, a, err...)
}

func (l *Logger) addTraceContextWithAttributes(span trace.Span, attributes []attribute.KeyValue, err ...error) []slog.Attr {
	var result []slog.Attr
	result = handleError(span, err, result)
	result = l.addTraceContextToLog(span, result)
	result = l.addServiceName(result)

	// _attributes are a global set of attributes added at initialization
	attributes = append(attributes, l.defaultAttributes...)

	// add attributes when global or passed in attributes are > 0
	result = append(result, l.ConvertToSlogFormat(attributes)...)

	return result
}

// WithTracingContext is a method for the Logger struct which takes a context.Context
// and log parameters, including a span from distributed tracing and (optional) error information.
// When logging without an error, pass a nil
func (l *Logger) WithTracingContext(ctx context.Context, level slog.Level, msg string, span trace.Span, err error, attrs ...slog.Attr) {
	attrs = append(attrs, AddTracingContext(span, err)...)
	l.LogAttrs(ctx, level, msg, attrs...)
}

// WithTracingContextAndAttributes is a method for the Logger struct which takes a context.Context
// and log parameters, including a span from distributed tracing, open telemetry attributes in the attribute.KeyValue format and (optional) error information.
// When logging without an error, pass a nil
func (l *Logger) WithTracingContextAndAttributes(ctx context.Context, level slog.Level, msg string, span trace.Span, err error, attributes []attribute.KeyValue, attrs ...slog.Attr) {
	attrs = append(attrs, l.addTraceContextWithAttributes(span, attributes, err)...)
	l.LogAttrs(ctx, level, msg, attrs...)
}

// ConvertToSlogFormat converts a list of attribute.KeyValue into the slog.Attr format
// and appends them to "result". The different types of attribute.KeyValue's
// are converted accordingly.
func (l *Logger) ConvertToSlogFormat(attributes []attribute.KeyValue) []slog.Attr {
	result := []slog.Attr(nil)
	if len(attributes) > 0 {
		for _, attr := range attributes {
			switch attr.Value.Type() {
			case attribute.STRING:
				result = append(result, slog.String(l.defaultAttrPrefix+"."+string(attr.Key), attr.Value.AsString()))
			case attribute.FLOAT64:
				result = append(result, slog.Float64(l.defaultAttrPrefix+"."+string(attr.Key), attr.Value.AsFloat64()))
			case attribute.BOOL:
				result = append(result, slog.Bool(l.defaultAttrPrefix+"."+string(attr.Key), attr.Value.AsBool()))
			case attribute.INT64:
				result = append(result, slog.Int64(l.defaultAttrPrefix+"."+string(attr.Key), attr.Value.AsInt64()))
			case attribute.BOOLSLICE:
				s := attr.Value.AsBoolSlice()
				vals := []slog.Attr(nil)
				for i, b := range s {
					key := l.defaultAttrPrefix + "." + string(attr.Key) + "." + strconv.Itoa(i)
					vals = append(vals, slog.Bool(key, b))
				}
				result = append(result, vals...)

			case attribute.INT64SLICE:
				s := attr.Value.AsInt64Slice()
				vals := []slog.Attr(nil)
				for i, b := range s {
					key := l.defaultAttrPrefix + "." + string(attr.Key) + "." + strconv.Itoa(i)
					vals = append(vals, slog.Int64(key, b))
				}
				result = append(result, vals...)
			case attribute.FLOAT64SLICE:
				s := attr.Value.AsFloat64Slice()
				vals := []slog.Attr(nil)
				for i, b := range s {
					key := l.defaultAttrPrefix + "." + string(attr.Key) + "." + strconv.Itoa(i)
					vals = append(vals, slog.Float64(key, b))
				}
				result = append(result, vals...)
			case attribute.STRINGSLICE:
				s := attr.Value.AsStringSlice()
				vals := []slog.Attr(nil)
				for i, b := range s {
					key := l.defaultAttrPrefix + "." + string(attr.Key) + "." + strconv.Itoa(i)
					vals = append(vals, slog.String(key, b))
				}
				result = append(result, vals...)
			default:
				if attr.Value.Type() != attribute.INVALID {
					//l.Debug().Caller().Str("attribute_key", string(attr.Key)).Interface("attribute_value", attr.Value).Msg("invalid attribute value type")
					l.LogAttrs(nil, slog.LevelDebug, "invalid attribute value type",
						slog.String("attribute_key", string(attr.Key)),
						slog.Any("attribute_value", attr.Value),
					)
				}
				continue
			}
		}
	}
	return result
}

func (l *Logger) addServiceName(result []slog.Attr) []slog.Attr {
	// set service.name if the value isn't empty
	if l.defaultServiceName != "" {
		result = append(result, slog.String("service.name", l.defaultServiceName))
	}
	return result
}

func (l *Logger) addTraceContextToLog(span trace.Span, result []slog.Attr) []slog.Attr {
	result = append(result,
		slog.String(l.defaultTraceId, span.SpanContext().TraceID().String()),
		slog.String(l.defaultSpanId, span.SpanContext().SpanID().String()),
	)
	return result
}

func handleError(span trace.Span, err []error, result []slog.Attr) []slog.Attr {
	if len(err) > 0 && err[0] != nil {
		span.RecordError(err[0])
		span.SetStatus(codes.Error, err[0].Error())
		result = append(result, slog.String("error", err[0].Error()))
	}
	return result
}

// New initializes a new Logger instance with the provided LogOptions.
// The Logger struct contains a wrapped slog.Logger, along with default configuration for trace and span IDs,
// service name, and default attributes.
// The options parameter takes any number of LogOptions, which are functions that modify the default settings for the logger.
func New(options ...LogOption) *Logger {
	return initLogger(options)
}

// NewWithHandler initializes a new Logger instance with a specified slog.Handler.
// If no handler is provided (i.e., handler is nil), it calls the New() function
// to create a Logger with the default setup.
//
// Deprecated: replaced by option WithOtelBridge, this creates a handler that sends logs to the otel endpoint in the otlp logRecord format.
//
//	Instead of NewWithHandler, try using the New function to get a logger with added functionality for appending trace info
func NewWithHandler(handler slog.Handler) *Logger {
	// If handler is nil, create a default Logger with json logging to standard output
	if handler == nil {
		return New(nil)
	}

	return &Logger{Logger: slog.New(handler)}
}
