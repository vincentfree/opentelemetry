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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	attributes      []attribute.KeyValue
	serviceName     string
	traceId         string
	spanId          string
	attributePrefix string
}

type Logger struct {
	*slog.Logger
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
func AddTracingContext(span trace.Span, err ...error) []slog.Attr {
	a := []attribute.KeyValue(nil)
	return AddTracingContextWithAttributes(span, a, err...)
}

// AddTracingContextWithAttributes lets you add the trace context to a structured log, including attribute.KeyValue's to extend the log
func AddTracingContextWithAttributes(span trace.Span, attributes []attribute.KeyValue, err ...error) []slog.Attr {
	var result []slog.Attr
	if len(err) > 0 && err[0] != nil {
		span.RecordError(err[0])
		span.SetStatus(codes.Error, err[0].Error())
		result = append(result, slog.String("error", err[0].Error()))
	}

	result = append(result,
		slog.String(_traceId, span.SpanContext().TraceID().String()),
		slog.String(_spanId, span.SpanContext().SpanID().String()),
	)
	// set service.name if the value isn't empty
	if _serviceName != "" {
		result = append(result, slog.String("service.name", _serviceName))
	}

	attrs := attributes
	attrs = append(attrs, _attributes...)

	// add attributes when global or passed attributes are > 0
	if len(attrs) > 0 {
		for _, attr := range attrs {
			switch attr.Value.Type() {
			case attribute.STRING:
				result = append(result, slog.String(_attrPrefix+"."+string(attr.Key), attr.Value.AsString()))
			case attribute.FLOAT64:
				result = append(result, slog.Float64(_attrPrefix+"."+string(attr.Key), attr.Value.AsFloat64()))
			case attribute.BOOL:
				result = append(result, slog.Bool(_attrPrefix+"."+string(attr.Key), attr.Value.AsBool()))
			case attribute.INT64:
				result = append(result, slog.Int64(_attrPrefix+"."+string(attr.Key), attr.Value.AsInt64()))
			case attribute.BOOLSLICE:
				s := attr.Value.AsBoolSlice()
				vals := []slog.Attr(nil)
				for i, b := range s {
					key := _attrPrefix + "." + string(attr.Key) + "." + strconv.FormatInt(int64(i), 10)
					vals = append(vals, slog.Bool(key, b))
				}
				result = append(result, vals...)

			case attribute.INT64SLICE:
				s := attr.Value.AsInt64Slice()
				vals := []slog.Attr(nil)
				for i, b := range s {
					key := _attrPrefix + "." + string(attr.Key) + "." + strconv.FormatInt(int64(i), 10)
					vals = append(vals, slog.Int64(key, b))
				}
				result = append(result, vals...)
			case attribute.FLOAT64SLICE:
				s := attr.Value.AsFloat64Slice()
				vals := []slog.Attr(nil)
				for i, b := range s {
					key := _attrPrefix + "." + string(attr.Key) + "." + strconv.FormatInt(int64(i), 10)
					vals = append(vals, slog.Float64(key, b))
				}
				result = append(result, vals...)
			case attribute.STRINGSLICE:
				s := attr.Value.AsStringSlice()
				vals := []slog.Attr(nil)
				for i, b := range s {
					key := _attrPrefix + "." + string(attr.Key) + "." + strconv.FormatInt(int64(i), 10)
					vals = append(vals, slog.String(key, b))
				}
				result = append(result, vals...)
			}
		}
	}

	return result
}

func New() *Logger {
	return &Logger{slog.New(slog.NewJSONHandler(os.Stdout, nil))}
}

func NewWithHandler(handler slog.Handler) *Logger {
	if handler == nil {
		return New()
	}
	return &Logger{slog.New(handler)}
}

func (l Logger) WithTracingContext(ctx context.Context, level slog.Level, msg string, span trace.Span, err error, attrs ...slog.Attr) {
	attrs = append(attrs, AddTracingContext(span, err)...)
	l.LogAttrs(ctx, level, msg, attrs...)
}

func (l Logger) WithTracingContextAndAttributes(ctx context.Context, level slog.Level, msg string, span trace.Span, err error, attributes []attribute.KeyValue, attrs ...slog.Attr) {
	attrs = append(attrs, AddTracingContextWithAttributes(span, attributes, err)...)
	l.LogAttrs(ctx, level, msg, attrs...)
}
