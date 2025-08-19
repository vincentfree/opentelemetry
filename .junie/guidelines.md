# Development Guidelines for OpenTelemetry Extensions

This document provides guidelines and instructions for developing and testing the OpenTelemetry extensions project.

## Build/Configuration Instructions

### Project Structure

The project is organized as a multi-module Go repository with the following structure:

- `cmd/`: Main applications and performance tests
- `otel*`: Go libraries for different logging frameworks:
  - `otellogrus`: Integration with Logrus
  - `otelzerolog`: Integration with Zerolog
  - `otelslog`: Integration with Go's standard slog
  - `otelmiddleware`: HTTP middleware for tracing
- `providerconfig`: Configuration utilities for OpenTelemetry providers
  - `providerconfighttp`: HTTP-specific provider configuration
  - `providerconfiggrpc`: gRPC-specific provider configuration
  - `providerconfignoop`: No-op provider configuration

### Building the Project

Each module in the project needs to be built separately. Use the following command in each module directory:

```bash
go build -v ./...
```

For example:

```bash
cd otelmiddleware
go build -v ./...
cd ../otelzerolog
go build -v ./...
# and so on for each module
```

### Dependencies

The project uses Go modules for dependency management. Each module has its own `go.mod` file that specifies its dependencies. The root `go.mod` file includes dependencies for the main project.

To update dependencies for a specific module:

```bash
cd <module-directory>
go get -u ./...
go mod tidy
```

## Testing Information

### Running Tests

Tests can be run for each module separately using the standard Go test command:

```bash
cd <module-directory>
go test -v ./...
```

For example:

```bash
cd otellogrus
go test -v ./...
```

### Adding New Tests

When adding new tests, follow these guidelines:

1. **Test File Naming**: Test files should be named with a `_test.go` suffix.
2. **Test Function Naming**: Test functions should be named with a `Test` prefix followed by the name of the functionality being tested.
3. **Example Tests**: Example tests should be named with an `Example` prefix and placed in a file with an `example_test.go` suffix.

### Test Helpers

The project includes several test helpers for common operations:

- Capturing log output
- Converting logs to maps for assertion
- Checking attribute values
- Verifying trace and span IDs

Here's an example of how to use these helpers in a test:

```go
func TestLogWithTracing(t *testing.T) {
    // Capture log output
    output := captureLog(t, func(logger *logrus.Logger) {
        // Log with tracing context
        logger.WithFields(otellogrus.AddTracingContext(span)).Info("test message")
    })

    // Convert log output to map for assertions
    logMap := logToMap(t, output)

    // Verify trace information is present
    attributeKeyCheck(t, logMap, "trace_id")
    attributeKeyCheck(t, logMap, "span_id")
}
```

## Additional Development Information

### Code Style

The project follows standard Go formatting and conventions:

- Use `gofmt` or `go fmt` to format your code.
- Follow the conventions outlined in [Effective Go](https://go.dev/doc/effective_go).
- Use meaningful variable and function names.
- Add comments for exported functions and types.

### Performance Considerations

The project includes performance benchmarks in the `cmd/perf_test.go` file. When making changes that might affect performance, run these benchmarks to ensure that performance is not degraded:

```bash
cd cmd
go test -bench=.
```

Different logging libraries have different performance characteristics, as shown in the benchmark results in the README.md file. Consider these differences when choosing a logging library for your application.

### Release Process

To release a new version of a module:

1. Update the `VERSIONS.md` file with the new version number.
2. Create a Git tag for the release:
   ```bash
   git tag -a <module-name>/v<version> -sm 'message'
   ```
   For example:
   ```bash
   git tag -a otelslog/v0.0.1 -sm 'Initial release'
   ```
3. Push the tag to the remote repository:
   ```bash
   git push origin <module-name>/v<version>
   ```
4. The Go proxy will automatically pick up the new tag and make it available on `pkg.go.dev`.

You can trigger the Go proxy to update with the following command:

```bash
curl https://proxy.golang.org/github.com/vincentfree/opentelemetry/<module-name>/@v/v<version>.info
```