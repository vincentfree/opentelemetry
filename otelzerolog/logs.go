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

package otelzerolog

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	otelzlog "go.opentelemetry.io/contrib/bridges/otelzerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/trace"
	"os"
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
	hook            *otelzlog.Hook
	zeroLogFeatures []func(zerolog.Context) zerolog.Context
	bridgeDisabled  bool
}

type Logger struct {
	zerolog.Logger
	// _traceId has a default trace ID key in the logs
	defaultTraceId string
	// _spanId has a default span ID key in the logs
	defaultSpanId string
	// _serviceName is empty by default, when no value is set the service name won't be used with a default in the logs
	defaultServiceName string
	// _attributes contains a global set of attribute.KeyValue's that will be added to very structured log. when the slice is empty they won't be added
	defaultAttributes []attribute.KeyValue
	// _attrPrefix x
	defaultAttrPrefix string
}

func defaultLogger() Logger {
	logger := Logger{}
	logger.Logger = log.Logger
	logger.defaultTraceId = "traceID"
	logger.defaultSpanId = "spanID"
	logger.defaultAttrPrefix = "trace.attribute"
	return logger
}

var (
	_logger     = defaultLogger()
	emptyLogger = Logger{}
)

type LogOptions []LogOption

func AsOtelLogger(logger zerolog.Logger) Logger {
	_logger.Logger = logger
	return _logger
}

// SetLogOptions takes LogOption's and overwrites library defaults
//
// Deprecated: due to naming, replaced by SetGlobalLogger and New.
// SetGlobalLogger is a direct replacement while New returns a zerolog.Logger with the same configuration.
func SetLogOptions(options ...LogOption) {
	SetGlobalLogger(options...)
}

// SetGlobalLogger sets the global logger for the library, using the provided LogOption entries.
//
// This function initializes a new logger using the LogOptions provided,
// and sets it as the global zerolog logger in the log package.
//
//	This functionality serves as a backwards compatible feature to how this library worked prior to version 0.1.0,
//	Please use the New function and use the logger provided to make use of all of its features.
//
// The LogOptions allow you to customize various aspects of the logger, such as the service name,
// trace and span ID names in the log entry, the attribute prefix, and zero log features like appending timestamps,
// setting the caller, etc.
func SetGlobalLogger(options ...LogOption) {
	_logger = initLogger(options)
	log.Logger = _logger.Logger
}

// initLogger initializes and configures a Logger based on the provided LogOptions.
// It applies the LogOptions to a logConfig struct, which holds various configuration parameters.
// The logger is finally returned and can be used for logging.
func initLogger(options LogOptions) Logger {
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

	if config.hook == nil && !config.bridgeDisabled {
		var name string
		if config.serviceName != "" {
			name = config.serviceName
		} else {
			name = "myApp"
		}
		config.hook = otelzlog.NewHook(name, otelzlog.WithLoggerProvider(noop.NewLoggerProvider()))
	}

	l := zerolog.New(os.Stdout).With().Logger()

	if !config.bridgeDisabled {
		l = l.Hook(config.hook)
	}

	if len(config.zeroLogFeatures) != 0 {
		for _, feature := range config.zeroLogFeatures {
			l = feature(l.With()).Logger()
		}
	}
	logger.Logger = l
	return logger
}

// New returns a zerolog.Logger configured with the provided LogOptions.
//
// The LogOptions are applied to a logConfig struct, which holds various configuration parameters.
// The logger is finally returned and can be used for logging.
func New(options ...LogOption) Logger {
	return initLogger(options)
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
func (l Logger) AddTracingContext(span trace.Span, err ...error) func(event *zerolog.Event) {
	a := []attribute.KeyValue(nil)
	return l.AddTracingContextWithAttributes(span, a, err...)
}

// AddTracingContextWithAttributes lets you add the trace context to a structured log, including attribute.KeyValue's to extend the log
func (l Logger) AddTracingContextWithAttributes(span trace.Span, attributes []attribute.KeyValue, err ...error) func(event *zerolog.Event) {
	return func(event *zerolog.Event) {

		if len(err) > 0 {
			span.RecordError(err[0])
			span.SetStatus(codes.Error, err[0].Error())
			event.Err(err[0])
		}

		e := event.Str(l.defaultTraceId, span.SpanContext().TraceID().String()).Str(l.defaultSpanId, span.SpanContext().SpanID().String())
		// set service.name if the value isn't empty
		if l.defaultServiceName != "" {
			e.Str("service.name", l.defaultServiceName)
		}

		attrs := attributes
		attrs = append(attrs, l.defaultAttributes...)

		// add attributes when global or passed attributes are > 0
		if len(attrs) > 0 {
			for _, attr := range attrs {
				switch attr.Value.Type() {
				case attribute.STRING:
					e.Str(l.defaultAttrPrefix+"."+string(attr.Key), attr.Value.AsString())
				case attribute.FLOAT64:
					e.Float64(l.defaultAttrPrefix+"."+string(attr.Key), attr.Value.AsFloat64())
				case attribute.BOOL:
					e.Bool(l.defaultAttrPrefix+"."+string(attr.Key), attr.Value.AsBool())
				case attribute.INT64:
					e.Int64(l.defaultAttrPrefix+"."+string(attr.Key), attr.Value.AsInt64())
				case attribute.BOOLSLICE:
					e.Bools(l.defaultAttrPrefix+"."+string(attr.Key), attr.Value.AsBoolSlice())
				case attribute.INT64SLICE:
					e.Ints64(l.defaultAttrPrefix+"."+string(attr.Key), attr.Value.AsInt64Slice())
				case attribute.FLOAT64SLICE:
					e.Floats64(l.defaultAttrPrefix+"."+string(attr.Key), attr.Value.AsFloat64Slice())
				case attribute.STRINGSLICE:
					e.Strs(l.defaultAttrPrefix+"."+string(attr.Key), attr.Value.AsStringSlice())
				default:
					if attr.Value.Type() != attribute.INVALID {
						l.Debug().Caller().Str("attribute_key", string(attr.Key)).Interface("attribute_value", attr.Value).Msg("invalid attribute value type")
					}
					continue
				}
			}
		}
	}
}
