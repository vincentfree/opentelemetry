module github.com/vincentfree/opentelemetry

go 1.22

require (
	github.com/rs/zerolog v1.33.0
	github.com/sirupsen/logrus v1.9.3
	github.com/vincentfree/opentelemetry/otellogrus v0.0.2
	github.com/vincentfree/opentelemetry/otelslog v0.0.0-00010101000000-000000000000
	github.com/vincentfree/opentelemetry/otelzerolog v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.30.0
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	go.opentelemetry.io/contrib/bridges/otelslog v0.3.0 // indirect
	go.opentelemetry.io/contrib/bridges/otelzerolog v0.0.0-20240726205214-b0584291236a // indirect
	go.opentelemetry.io/otel/log v0.4.0 // indirect
	go.opentelemetry.io/otel/metric v1.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.30.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
)

replace github.com/vincentfree/opentelemetry/otelzerolog => ./otelzerolog

replace github.com/vincentfree/opentelemetry/otelslog => ./otelslog
