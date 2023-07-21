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
Package otelzerolog provides a function to extend structured logs using zerolog with the Open Telemetry trace related context.

The github.com/rs/zerolog zerolog.Event is decorated with standard metadata extracted from the trace.SpanContext, a traceID, spanID and additional information is injected into a log.

The initialization uses file level configuration to set defaults for the function to use. SetLogOptions can overwrite the defaults.

When the configuration is done AddTracingContext and AddTracingContextWithAttributes decorate zerolog logs with data from the trace context.
A zeroLog.Event can be passed by using log.Info().Func(AddTracingContext(span)).Msg("") for example.

Functions

	func SetLogOptions(options ...LogOption)
	func WithTraceID(traceID string) LogOption
	func WithSpanID(spanID string) LogOption
	func WithServiceName(serviceName string) LogOption
	func WithAttributePrefix(prefix string) LogOption
	func WithAttributes(attributes ...attribute.KeyValue) LogOption
	func AddTracingContext(span trace.Span, err ...error) func(event *zerolog.Event)
	func AddTracingContextWithAttributes(span trace.Span, attributes []attribute.KeyValue, err ...error) func(event *zerolog.Event)

Types

	type LogOption func(*logConfig)

Structs

		type logConfig struct {
			attributes      []attribute.KeyValue
			serviceName     string
			traceId         string
			spanId          string
			attributePrefix string
	}

import "github.com/vincentfree/opentelemetry/otelzerolog"
*/
package otelzerolog
