module github.com/vincentfree/opentelemetry/otelzerolog

go 1.23

toolchain go1.23.5

require (
	github.com/rs/zerolog v1.33.0
	go.opentelemetry.io/contrib/bridges/otelzerolog v0.0.0-20240809024635-0c3fcdf3c470
	go.opentelemetry.io/otel v1.34.0
	go.opentelemetry.io/otel/log v0.10.0
	go.opentelemetry.io/otel/trace v1.34.0
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)

retract v0.1.0
