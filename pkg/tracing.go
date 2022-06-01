package otelmiddleware

import (
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

const version = "0.0.1"

type TraceOption func(*traceConfig)

type traceConfig struct {
	serviceName string
	tracer      trace.Tracer
	propagator  propagation.TextMapPropagator
}

func TraceWithOptions(opt ...TraceOption) func(next http.Handler) http.Handler {
	config := &traceConfig{}
	for _, o := range opt {
		o(config)
	}
	if config.tracer == nil {
		config.tracer = otel.Tracer("otel-tracer", trace.WithInstrumentationVersion(version))
	}
	if config.propagator == nil {
		config.propagator = otel.GetTextMapPropagator()
	}
	if config.serviceName == "" {
		config.serviceName = "TracedApplication"
	}
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestCtx := r.Context()
			ctx := config.propagator.Extract(requestCtx, propagation.HeaderCarrier(r.Header))
			opts := []trace.SpanStartOption{
				trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
				trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
				trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(r.Host, extractRoute(r.RequestURI), r)...),
				trace.WithSpanKind(trace.SpanKindServer),
			}
			spanName := extractRoute(r.RequestURI)
			if spanName == "" {
				spanName = fmt.Sprintf("HTTP %s route not found", r.Method)
			}
			spanName = fmt.Sprintf("%s /%s", r.Method, spanName)
			ctx, span := config.tracer.Start(ctx, spanName, opts...)
			defer span.End()

			// pass the span through the request context
			r = r.WithContext(ctx)

			// serve the request to the next middleware
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func Trace(next http.Handler) http.Handler {
	return TraceWithOptions(nil)(next)
}

func extractRoute(uri string) string {
	return uri[1:]
}

func WithTracer(tracer trace.Tracer) TraceOption {
	return func(c *traceConfig) {
		c.tracer = tracer
	}
}

func WithPropagator(p propagation.TextMapPropagator) TraceOption {
	return func(c *traceConfig) {
		c.propagator = p
	}
}

func WithServiceName(service string) TraceOption {
	return func(c *traceConfig) {
		c.serviceName = service
	}
}
