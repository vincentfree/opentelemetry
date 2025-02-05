module github.com/vincentfree/opentelemetry/otellogrus

go 1.23

toolchain go1.23.5

require (
	github.com/sirupsen/logrus v1.9.3
	go.opentelemetry.io/otel v1.34.0
	go.opentelemetry.io/otel/trace v1.34.0
	golang.org/x/exp v0.0.0-20230728194245-b0cb94b80691
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)

retract v0.0.1
