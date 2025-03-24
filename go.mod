module github.com/vincentfree/opentelemetry

go 1.23

require (
	github.com/rs/zerolog v1.34.0
	github.com/sirupsen/logrus v1.9.3
	github.com/vincentfree/opentelemetry/otellogrus v0.0.2
	github.com/vincentfree/opentelemetry/otelslog v0.0.0-00010101000000-000000000000
	github.com/vincentfree/opentelemetry/otelzerolog v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.34.0
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/bridges/otelslog v0.9.0 // indirect
	go.opentelemetry.io/contrib/bridges/otelzerolog v0.0.0-20240809024635-0c3fcdf3c470 // indirect
	go.opentelemetry.io/otel/log v0.10.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)

replace github.com/vincentfree/opentelemetry/otelzerolog => ./otelzerolog

replace github.com/vincentfree/opentelemetry/otelslog => ./otelslog
