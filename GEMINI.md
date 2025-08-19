# Gemini Code Assist Configuration

This file helps Gemini Code Assist understand your project's structure and conventions.

## Project Information

*   **Language:** Go
*   **Frameworks:** None
*   **Build System:** Go Modules
*   **Test Framework:** Go Test

## Build and Test Commands

The following commands are used to build and test the project:

*   **Build:**
    *   `go build -v ./...` (run in each of the following directories: `otelmiddleware`, `otelzerolog`, `otelslog`, `otellogrus`, `providerconfig`, `providerconfig/providerconfighttp`, `providerconfig/providerconfiggrpc`)
*   **Test:**
    *   `go test -v ./...` (run in each of the following directories: `otelmiddleware`, `otelzerolog`, `otelslog`, `otellogrus`, `providerconfig`, `providerconfig/providerconfighttp`, `providerconfig/providerconfiggrpc`)

## File and Directory Structure

*   `cmd/`: Main applications
*   `otel*`: Go libraries
*   `go.mod`, `go.sum`: Go module definitions

## Coding Style

*   Follows standard Go formatting (`gofmt`).
*   Uses the conventions outlined in [Effective Go](https://go.dev/doc/effective_go).

## Releasing to pkg.go.dev

To publish a new version of a library to `pkg.go.dev`, follow these steps:

1.  **Tag the commit:** Create a new Git tag for the release.
    ```bash
    git tag otelmiddleware/v1.2.3
    ```
2.  **Push the tag:** Push the new tag to the remote repository. The tag name must match the module name.
    ```bash
    git push origin otelmiddleware/v1.2.3
    ```
3.  **Update the proxy:** The Go proxy will automatically pick up the new tag and make it available on `pkg.go.dev`. This may take a few minutes. You can check the status on the [Go proxy status page](https://proxy.golang.org/).
