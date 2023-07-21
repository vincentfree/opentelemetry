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

package otelmiddleware_test

import (
	"github.com/vincentfree/opentelemetry/otelmiddleware"
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
