# OpenTelemetry extensions


[![Go](https://github.com/vincentfree/opentelemetry/actions/workflows/go.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/go.yml)
[![CodeQL](https://github.com/vincentfree/opentelemetry/actions/workflows/codeql.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/codeql.yml)
[![Dependency Review](https://github.com/vincentfree/opentelemetry/actions/workflows/dependency-review.yml/badge.svg)](https://github.com/vincentfree/opentelemetry/actions/workflows/dependency-review.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/vincentfree/opentelemetry/otelmiddleware.svg)](https://pkg.go.dev/github.com/vincentfree/opentelemetry)

These libraries are ment as extensions for the Open Telemetry project. 
They provide functionality that makes working with traces, and logs easier.

Currently, there is support for:

* http severs through [otelmiddleware](otelmiddleware/README.md)
* logging with [zerolog](otelzerolog/README.md), [slog](otelslog/README.md), [logrus](otellogrus/README.md)

More extensions might follow for other logging libraries and more.

## Logging benchmark results

```shell
goos: darwin
goarch: arm64
pkg: github.com/vincentfree/opentelemetry/cmd
BenchmarkLogrus-10                    1000000        1054 ns/op     1064 B/op       22 allocs/op
BenchmarkLogrusTrace-10                507514        2255 ns/op     2456 B/op       38 allocs/op
BenchmarkLogrusTraceWithAttr-10        252018        4879 ns/op     3786 B/op       69 allocs/op
BenchmarkSlog-10                      2514039       477.0 ns/op        0 B/op        0 allocs/op
BenchmarkSlogTrace-10                 1460439       821.5 ns/op      176 B/op        5 allocs/op
BenchmarkSlogTraceWithAttr-10          463041        2538 ns/op     1592 B/op       32 allocs/op
BenchmarkZerolog-10                  22955486       51.40 ns/op        0 B/op        0 allocs/op
BenchmarkZerologTrace-10              6346458       195.8 ns/op      128 B/op        3 allocs/op
BenchmarkZerologTraceWithAttr-10      1254688       952.4 ns/op      312 B/op       12 allocs/op
```

> pre Go 1.21(slog was based on the golang.org/x/exp version not log/slog)
