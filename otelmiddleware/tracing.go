package otelmiddleware

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

// version is used as the instrumentation version.
const version = "0.0.6"

// TraceOption takes a traceConfig struct and applies changes.
// It can be passed to the TraceWithOptions function to configure a traceConfig struct.
type TraceOption func(*traceConfig)

// traceConfig contains all the configuration for the library.
type traceConfig struct {
	serviceName string
	tracer      trace.Tracer
	propagator  propagation.TextMapPropagator
	attributes  []attribute.KeyValue
}

// TraceWithOptions takes TraceOption's and initializes a new trace.Span.
func TraceWithOptions(opt ...TraceOption) func(next http.Handler) http.Handler {
	// initialize an empty traceConfig.
	config := &traceConfig{}

	// apply the configuration passed to the function.
	for _, o := range opt {
		o(config)
	}
	// check for the traceConfig.tracer if absent use a default value.
	if config.tracer == nil {
		config.tracer = otel.Tracer("otel-tracer", trace.WithInstrumentationVersion(version))
	}
	// check for the traceConfig.propagator if absent use a default value.
	if config.propagator == nil {
		config.propagator = otel.GetTextMapPropagator()
	}
	// check for the traceConfig.serviceName if absent use a default value.
	if config.serviceName == "" {
		config.serviceName = "TracedApplication"
	}
	// the handler that initializes the trace.Span.
	return func(next http.Handler) http.Handler {
		// assign the handler which creates the OpenTelemetry trace.Span.
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestCtx := r.Context()
			// extract the OpenTelemetry span context from the context.Context object.
			ctx := config.propagator.Extract(requestCtx, propagation.HeaderCarrier(r.Header))
			ctxSpan := trace.SpanFromContext(ctx)
			var span trace.Span
			if ctxSpan.IsRecording() {
				span = ctxSpan
				attr := []attribute.KeyValue(nil)
				attr = append(attr, semconv.NetAttributesFromHTTPRequest("tcp", r)...)
				attr = append(attr, semconv.EndUserAttributesFromHTTPRequest(r)...)
				attr = append(attr, semconv.HTTPServerAttributesFromHTTPRequest(r.Host, extractRoute(r.RequestURI), r)...)

				if len(config.attributes) > 0 {
					attr = append(attr, config.attributes...)
				}

				span.SetAttributes(attr...)
			} else {
				// the standard trace.SpanStartOption options whom are applied to every server handler.
				opts := []trace.SpanStartOption{
					trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
					trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
					trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(r.Host, extractRoute(r.RequestURI), r)...),
					trace.WithSpanKind(trace.SpanKindServer),
				}

				// check for the traceConfig.attributes if present apply them to the trace.Span.
				if len(config.attributes) > 0 {
					opts = append(opts, trace.WithAttributes(config.attributes...))
				}
				// extract the route name which is used for setting a usable name of the span.
				spanName := extractRoute(r.RequestURI)
				if spanName == "" {
					spanName = fmt.Sprintf("HTTP %s route not found", r.Method)
				}

				// create a good name to recognize where the span originated.
				spanName = fmt.Sprintf("%s /%s", r.Method, spanName)

				// start the actual trace.Span.
				ctx, span = config.tracer.Start(ctx, spanName, opts...)
			}
			// end span after the function ends
			defer span.End()

			// pass the span through the request context.
			r = r.WithContext(ctx)

			// serve the request to the next middleware.
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

// Trace uses the TraceWithOptions without additional options, this is a shorthand for TraceWithOptions().
func Trace(next http.Handler) http.Handler {
	return TraceWithOptions()(next)
}

// extract the route name.
func extractRoute(uri string) string {
	return uri[1:]
}

// WithTracer is a TraceOption to inject your own trace.Tracer.
func WithTracer(tracer trace.Tracer) TraceOption {
	return func(c *traceConfig) {
		c.tracer = tracer
	}
}

// WithPropagator is a TraceOption to inject your own propagation.
func WithPropagator(p propagation.TextMapPropagator) TraceOption {
	return func(c *traceConfig) {
		c.propagator = p
	}
}

// WithServiceName is a TraceOption to inject your own serviceName.
func WithServiceName(serviceName string) TraceOption {
	return func(c *traceConfig) {
		c.serviceName = serviceName
	}
}

// WithAttributes is a TraceOption to inject your own attributes.
// Attributes are applied to the trace.Span.
func WithAttributes(attributes ...attribute.KeyValue) TraceOption {
	return func(c *traceConfig) {
		c.attributes = attributes
	}
}
