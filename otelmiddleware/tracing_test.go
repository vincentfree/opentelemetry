package otelmiddleware

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type testHandler func(http.ResponseWriter, *http.Request)

func (th testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	th(w, r)
}

func TestTraceFunctions(t *testing.T) {
	text := "Hi, test user"
	testCases := []struct {
		desc          string
		serverHandler testHandler
		fn            func(next http.Handler) http.Handler
	}{
		{
			desc: "Standard Trace handler",
			serverHandler: testHandler(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(text))
				if err != nil {
					panic(err)
				}
			}),
			fn: Trace,
		},
		{
			desc: "Trace handler With option: WithServiceName",
			serverHandler: testHandler(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(text))
				if err != nil {
					panic(err)
				}
			}),
			fn: TraceWithOptions(WithServiceName("otherName")),
		},
		{
			desc: "Trace handler With option: WithPropagator",
			serverHandler: testHandler(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(text))
				if err != nil {
					panic(err)
				}
			}),
			fn: TraceWithOptions(WithPropagator(otel.GetTextMapPropagator())),
		},
		{
			desc: "Trace handler With option: WithTracer",
			serverHandler: testHandler(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(text))
				if err != nil {
					panic(err)
				}
			}),
			fn: TraceWithOptions(WithTracer(otel.Tracer("test-tracer"))),
		},
		{
			desc: "Trace handler With option: WithAttributes",
			serverHandler: testHandler(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(text))
				if err != nil {
					panic(err)
				}
			}),
			fn: TraceWithOptions(WithAttributes(attribute.String("test", "test"))),
		},
		{
			desc: "Trace handler With all options",
			serverHandler: testHandler(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(text))
				if err != nil {
					panic(err)
				}
			}),
			fn: TraceWithOptions(
				WithTracer(otel.Tracer("test-tracer")),
				WithPropagator(otel.GetTextMapPropagator()),
				WithServiceName("testName"),
				WithAttributes(attribute.String("test", "test")),
			),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			traceHandler := tC.fn(tC.serverHandler)
			server := httptest.NewServer(traceHandler)
			server.Config.Handler = traceHandler
			defer server.Close()
			client := server.Client()
			response, err := client.Get(fmt.Sprintf("%s/", server.URL))
			if err != nil {
				t.Fatalf("calling the server failed due to: %v", err)
			}
			if response.StatusCode != 200 {
				t.Errorf("call dit not have the expected 200 statuscode but was: %d", response.StatusCode)
			}
			scanner := bufio.NewScanner(response.Body)
			scanner.Scan()
			result := scanner.Text()
			if result != text {
				t.Errorf("The response body should be: '%s', but was: '%s'", text, result)
			}
		})
	}

}
