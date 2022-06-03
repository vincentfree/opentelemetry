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

After these options are applied a new span is created and the middleware will pass the http.ResponseWriter and http.Request to the next http.Handler.

# Functions
	func TraceWithOptions(opt ...TraceOption) func(next http.Handler) http.Handler
	func Trace(next http.Handler) http.Handler
	func WithAttributes(attributes ...attribute.KeyValue) TraceOption
	func WithPropagator(p propagation.TextMapPropagator) TraceOption
	func WithServiceName(serviceName string) TraceOption
	func WithTracer(tracer trace.Tracer) TraceOption

# Types
	type TraceOption func(*traceConfig)

# Structs
	type traceConfig struct {
		tracer trace.Tracer
		propagator propagation.TextMapPropagator
		attributes []attribute.KeyValue
		serviceName string
	}


# Examples
	// create a new ServeMux
	serve := http.NewServeMux()
	// add a new route to the ServeMux
	serve.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	}))
	// create the Trace middleware and decorate the ServeMux routes with this middleware.
	handler := TraceWithOptions(WithServiceName("ExampleService"))(serve)
	http.ListenAndServe(":8080", handler)

*/
package otelmiddleware
