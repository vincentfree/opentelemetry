package otelmiddleware_test

import (
	"log"
	"net/http"

	"github.com/vincentfree/opentelemetry-http/otelmiddleware"
)

func Example() {
	// create a new ServeMux
	serve := http.NewServeMux()
	// add a new route to the ServeMux
	serve.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello World"))
		if err != nil {
			// handle error
		}
	}))
	// create the Trace middleware and decorate the ServeMux routes with this middleware.
	handler := otelmiddleware.TraceWithOptions(otelmiddleware.WithServiceName("ExampleService"))(serve)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
