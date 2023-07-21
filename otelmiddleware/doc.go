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
Package otelmiddleware provides middleware for wrapping http.Server handlers with Open Telemetry tracing support.

The trace.Span is decorated with standard metadata extracted from the http.Request injected into the middleware.
the basic information is extracted using the OpenTelemetry semconv package.

When a span gets initialized it uses the following slice of trace.SpanStartOption

	opts := []trace.SpanStartOption{
		trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", request)...),
		trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(request)...),
		trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(request.Host, extractRoute(request.RequestURI), request)...),
		trace.WithSpanKind(trace.SpanKindServer),
	}

The slice can be extended using the WithAttributes TraceOption function.

After these options are applied a new span is created and the middleware will pass the http.ResponseWriter and http.Request to the next http.Handler.

Functions

	func TraceWithOptions(opt ...TraceOption) func(next http.Handler) http.Handler
	func Trace(next http.Handler) http.Handler
	func WithAttributes(attributes ...attribute.KeyValue) TraceOption
	func WithPropagator(p propagation.TextMapPropagator) TraceOption
	func WithServiceName(serviceName string) TraceOption
	func WithTracer(tracer trace.Tracer) TraceOption

Types

	type TraceOption func(*traceConfig)

Structs

	type traceConfig struct {
		tracer trace.Tracer
		propagator propagation.TextMapPropagator
		attributes []attribute.KeyValue
		serviceName string
	}
*/
package otelmiddleware // import "github.com/vincentfree/opentelemetry/otelmiddleware"
