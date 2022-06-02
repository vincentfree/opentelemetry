/*
Package otelmiddleware provides middleware for wrapping http.Server handlers with Open Telemetry tracing support.

The trace.Span is decorated with standard meta data extracted from the http.Request injected into the middleware.
the basic information is extraced using the OpenTelemetry semconv package.

When a span gets initialized it uses the following slice of trace.SpanStartOption

	opts := []trace.SpanStartOption{
				trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", request)...),
				trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(request)...),
				trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(request.Host, extractRoute(request.RequestURI), request)...),
				trace.WithSpanKind(trace.SpanKindServer),
	}

After these options are applied a new span is created and the middleware will pass the http.ResponseWriter and http.Request to the next http.Handler.

*/
package otelmiddleware
