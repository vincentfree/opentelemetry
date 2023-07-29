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

package otellogrus

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"strings"
)

// LogOption takes a logConfig struct and applies changes.
// It can be passed to the SetLogOptions function to configure a logConfig struct.
type LogOption func(*logConfig)
type LoggerOption func(*loggerConfig)

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

var (
	// _traceId has a default trace ID key in the logs
	_traceId = "traceID"
	// _spanId has a default span ID key in the logs
	_spanId = "spanID"
	// _serviceName is empty by default, when no value is set the service name won't be used with a default in the logs
	_serviceName string
	// _attributes contains a global set of attribute.KeyValue's that will be added to very structured log. when the slice is empty they won't be added
	_attributes = []attribute.KeyValue(nil)

	// _attrPrefix x
	_attrPrefix = "trace.attribute"
)

// SetLogOptions takes LogOption's and overwrites library defaults
func SetLogOptions(options ...LogOption) {
	// initialize an empty logConfig
	config := &logConfig{}

	for _, option := range options {
		option(config)
	}

	if config.traceId != "" {
		_traceId = config.traceId
	}

	if config.spanId != "" {
		_traceId = config.spanId
	}

	if config.serviceName != "" {
		_serviceName = config.serviceName
	}

	if config.attributePrefix != "" {
		ap, _ := strings.CutSuffix(config.attributePrefix, ".")
		_attrPrefix = ap
	}

	if len(config.attributes) > 0 {
		_attributes = append(_attributes, config.attributes...)
	}
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

// AddTracingContext lets you add the trace context to a structured log
func AddTracingContext(span trace.Span, err ...error) logrus.Fields {
	a := []attribute.KeyValue(nil)
	return AddTracingContextWithAttributes(span, a, err...)
}

// AddTracingContextWithAttributes lets you add the trace context to a structured log, including attribute.KeyValue's to extend the log
func AddTracingContextWithAttributes(span trace.Span, attributes []attribute.KeyValue, err ...error) logrus.Fields {
	fields := logrus.Fields{}
	if len(err) > 0 {
		span.RecordError(err[0])
		span.SetStatus(codes.Error, err[0].Error())
		fields["error"] = err[0]
	}
	fields[_traceId] = span.SpanContext().TraceID().String()
	fields[_spanId] = span.SpanContext().SpanID().String()
	// set service.name if the value isn't empty
	if _serviceName != "" {
		fields["service.name"] = _serviceName
	}

	attrs := attributes
	attrs = append(attrs, _attributes...)

	// add attributes when global or passed attributes are > 0
	if len(attrs) > 0 {
		for _, attr := range attrs {
			switch attr.Value.Type() {
			case attribute.STRING:
				fields[_attrPrefix+"."+string(attr.Key)] = attr.Value.AsString()
			case attribute.FLOAT64:
				fields[_attrPrefix+"."+string(attr.Key)] = attr.Value.AsFloat64()
			case attribute.BOOL:
				fields[_attrPrefix+"."+string(attr.Key)] = attr.Value.AsBool()
			case attribute.INT64:
				fields[_attrPrefix+"."+string(attr.Key)] = attr.Value.AsInt64()
			case attribute.BOOLSLICE:
				fields[_attrPrefix+"."+string(attr.Key)] = attr.Value.AsBoolSlice()
			case attribute.INT64SLICE:
				fields[_attrPrefix+"."+string(attr.Key)] = attr.Value.AsInt64Slice()
			case attribute.FLOAT64SLICE:
				fields[_attrPrefix+"."+string(attr.Key)] = attr.Value.AsFloat64Slice()
			case attribute.STRINGSLICE:
				fields[_attrPrefix+"."+string(attr.Key)] = attr.Value.AsStringSlice()
			}
		}
	}
	return fields
}

// WithLevel sets the log level of the Logger instance.
//
// It receives a logrus.Level value as a parameter.
//
// Args:
//
//	level (logrus.Level): The level of logging required.
//
// Returns:
//
//	LoggerOption: Returns a function which modifies the 'level' attribute of a loggerConfig instance.
func WithLevel(level logrus.Level) LoggerOption {
	return func(c *loggerConfig) {
		c.level = level
	}
}

// WithFormatter sets a formatter.
// The formatter is used by the Logger.
//
// It receives a logrus.Formatter value as a parameter.
//
// Args:
//
//	formatter (logrus.Formatter): The formatter required for logging.
//
// Returns:
//
//	LoggerOption: Returns a function which modifies the 'formatter' attribute of a loggerConfig instance.
func WithFormatter(formatter logrus.Formatter) LoggerOption {
	return func(c *loggerConfig) {
		c.formatter = &formatter
	}
}

/*
New creates a new Logger with given options.

Parameters:
  - options (LoggerOption): Variadic configuration options for Logger.

Return:
  - A pointer to new Logger.
*/
func New(options ...LoggerOption) *Logger {
	l := &Logger{logrus.New()}
	c := &loggerConfig{level: logrus.InfoLevel}
	for _, opt := range options {
		opt(c)
	}

	if c.formatter != nil {
		l.Formatter = *c.formatter
	}
	// always set the level, the default is info
	l.Level = c.level

	return l
}

// WithTracingContext is a method on the Logger type. It uses the
// AddTracingContext helper function to gather tracing context from
// the given span and error, then creates a logrus.Entry with
// the context using the WithFields method of the Logger.
//
// The span is of the type trace.Span that provides the tracing information.
//
// A variadic parameter of error values is also provided which can be optionally used.
// only the first error is actually used, the variadic nature of this parameter is used to make it optional.
//
// It ultimately returns a pointer to a logrus.Entry populated with the tracing context.
func (l Logger) WithTracingContext(span trace.Span, err ...error) *logrus.Entry {
	return l.WithFields(AddTracingContext(span, err...))
}

// WithTracingContextAndAttributes is a method on the Logger type. Similar to WithTracingContext,
// this method uses a helper function (AddTracingContextWithAttributes in this case) to gather tracing
// context and attributes from the given span and error. It then creates a logrus.Entry with
// the context and attributes using the WithFields method of the Logger.
//
// The span is of the type trace.Span that provides the tracing information.
//
// Additionally, this method takes an array of attribute.KeyValue pairs (attributes). Each pair
// contains key-value attribute information that is added to the logging.
//
// Lastly, a variadic parameter of error values is provided which can be optionally used.
// Despite the parameter being variadic, only the first error is actually used; the variadic nature
// of this parameter is used to make it optional.
//
// It ultimately returns a pointer to a logrus.Entry populated with the tracing context and attributes.
func (l Logger) WithTracingContextAndAttributes(span trace.Span, attributes []attribute.KeyValue, err ...error) *logrus.Entry {
	return l.WithFields(AddTracingContextWithAttributes(span, attributes, err...))
}
