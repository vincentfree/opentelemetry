# OpenTelemetry extensions


[![Go](https://github.com/vincentfree/opentelemetry/actions/workflows/go.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/go.yml)
[![CodeQL](https://github.com/vincentfree/opentelemetry/actions/workflows/codeql.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/codeql.yml)
[![Dependency Review](https://github.com/vincentfree/opentelemetry/actions/workflows/dependency-review.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/dependency-review.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/vincentfree/opentelemetry/otelmiddleware.svg)](https://pkg.go.dev/github.com/vincentfree/opentelemetry/otelmiddleware)

These libraries are ment as extensions for the Open Telemetry project. 
They provide functionality that makes working with traces, and logs easier.

Currently, there is support for:

* http severs through [otelmiddleware](otelmiddleware/README.md)
* logging with [zerolog](otelzerolog/README.md), [slog](otelslog/README.md), [logrus](otellogrus/README.md)

More extensions might follow for other logging libraries and more.
