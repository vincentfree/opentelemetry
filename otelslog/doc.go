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

/*
Package otelslog provides a function to extend structured logs using slog with the Open Telemetry trace related context.

Currently, slog is offered through golang.org/x/exp/slog slog.Logger is decorated with standard metadata extracted from the trace.SpanContext, a traceID, spanID and additional information is injected into a log.

The initialization uses file level variable configuration to set defaults for the functions to use. SetLogOptions can overwrite the defaults.

When the configuration is done AddTracingContext and AddTracingContextWithAttributes decorate slog logs with data from the trace context.

To add trace context data to logging, the context can be passed by using slog.LogAttrs(nil, slog.LevelInfo, "this is a log", otelslog.AddTracingContext(span)...) for example.
The use of slog.LogAttrs is advised due to AddTracingContext and AddTracingContextWithAttributes returning []slog.Attr which slog.LogAttrs accepts as a type.
Other functions accept ...any which in my tests resulted in !BADKEY entries.

Next to using native slog, this package also offers a Logger which extends the slog.Logger with its own functions to simplify working with slog.Logger's.

The Logger can be used as follows:

	logger := otelslog.New()
	// pass span to AddTracingContext
	logger.WithTracingContext(nil, slog.LevelInfo, "in case of a success", span, nil)
	err := errors.New("example error"))
	// error case with attributes
	logger.WithTracingContextAndAttributes(ctx, slog.LevelError, "in case of a failure", span, err, attributes)

Functions

	func SetLogOptions(options ...LogOption)
	func WithTraceID(traceID string) LogOption
	func WithSpanID(spanID string) LogOption
	func WithServiceName(serviceName string) LogOption
	func WithAttributePrefix(prefix string) LogOption
	func WithAttributes(attributes ...attribute.KeyValue) LogOption
	func AddTracingContext(span trace.Span, err ...error) func(event *zerolog.Event)
	func AddTracingContextWithAttributes(span trace.Span, attributes []attribute.KeyValue, err ...error) func(event *zerolog.Event)
	func New() *Logger
	func NewWithHandler(handler slog.Handler) *Logger
	func (l Logger) WithTracingContext(ctx context.Context, level slog.Level, msg string, span trace.Span, err error, attrs ...slog.Attr)
	func (l Logger) WithTracingContextAndAttributes(ctx context.Context, level slog.Level, msg string, span trace.Span, err error, attributes []attribute.KeyValue, attrs ...slog.Attr)

Types

	type LogOption func(*logConfig)

Structs

		type Logger struct {
			*slog.Logger
		}

		type logConfig struct {
			attributes      []attribute.KeyValue
			serviceName     string
			traceId         string
			spanId          string
			attributePrefix string
	}

import "github.com/vincentfree/opentelemetry/otelslog"
*/
package otelslog
