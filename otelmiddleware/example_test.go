package otelmiddleware_test

import (
	"github.com/vincentfree/opentelemetry-http/otelmiddleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"net/http"
)

type exampleHandler func(http.ResponseWriter, *http.Request)

func (th exampleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	th(w, r)
}

var (
	eh = exampleHandler(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello World"))
	})
)

func ExampleWithTracer() {
	// returns a function that excepts a http.Handler.
	handler := otelmiddleware.TraceWithOptions(otelmiddleware.WithTracer(otel.Tracer("example-tracer")))
	// pass a http.Handler to extend it with Tracing functionality.
	http.Handle("/", handler(eh))
}

func ExampleWithServiceName() {
	// returns a function that excepts a http.Handler.
	handler := otelmiddleware.TraceWithOptions(otelmiddleware.WithServiceName("exampleService"))
	// pass a http.Handler to extend it with Tracing functionality.
	http.Handle("/", handler(eh))
}

func ExampleWithAttributes() {
	// returns a function that excepts a http.Handler.
	handler := otelmiddleware.TraceWithOptions(otelmiddleware.WithAttributes(attribute.String("example", "value")))
	// pass a http.Handler to extend it with Tracing functionality.
	http.Handle("/", handler(eh))
}

func ExampleWithPropagator() {
	// returns a function that excepts a http.Handler.
	handler := otelmiddleware.TraceWithOptions(otelmiddleware.WithPropagator(otel.GetTextMapPropagator()))
	// pass a http.Handler to extend it with Tracing functionality.
	http.Handle("/", handler(eh))
}

func ExampleTraceWithOptions() {
	// returns a function that excepts a http.Handler.
	handler := otelmiddleware.TraceWithOptions(otelmiddleware.WithServiceName("exampleService"))
	// pass a http.Handler to extend it with Tracing functionality.
	http.Handle("/", handler(eh))
}

func ExampleTrace() {
	http.Handle("/", otelmiddleware.Trace(eh))
}
