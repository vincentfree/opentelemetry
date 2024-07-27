module github.com/vincentfree/opentelemetry

go 1.22

require (
	github.com/rs/zerolog v1.33.0
	github.com/sirupsen/logrus v1.9.3
	github.com/vincentfree/opentelemetry/otellogrus v0.0.2
	github.com/vincentfree/opentelemetry/otelslog v0.0.3
	github.com/vincentfree/opentelemetry/otelzerolog v0.0.11
	go.opentelemetry.io/otel v1.28.0
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
)

replace (
	github.com/vincentfree/opentelemetry/otelzerolog => ./otelzerolog
)