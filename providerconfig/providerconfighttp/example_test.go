// Copyright 2024 Vincent Free
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package providerconfighttp_test

import (
	"github.com/vincentfree/opentelemetry/providerconfig"
	"github.com/vincentfree/opentelemetry/providerconfig/providerconfighttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/trace"
)

func ExampleNew() {
	signalProcessor := providerconfighttp.New(
		providerconfighttp.WithCollectorEndpoint("0.0.0.0", 9898),
	)
	provider := providerconfig.New(providerconfig.WithApplicationName("example-app"),
		providerconfig.WithApplicationVersion("0.1.0"),
		providerconfig.WithSignalProcessor(signalProcessor),
	)
	// set providers
	otel.SetTracerProvider(provider.TraceProvider())
	otel.SetMeterProvider(provider.MetricProvider())
	global.SetLoggerProvider(provider.LogProvider())

	// shutdown all providers
	provider.ShutdownAll()
}

func ExampleWithCollectorEndpoint() {
	providerconfighttp.New(
		providerconfighttp.WithCollectorEndpoint("0.0.0.0", 9898),
	)
}

func ExampleWithTraceOptions() {
	providerconfighttp.New(
		providerconfighttp.WithTraceOptions(otlptracehttp.WithCompression(otlptracehttp.GzipCompression)),
	)
}

func ExampleWithSpanProcessorOptions() {
	providerconfighttp.New(
		providerconfighttp.WithSpanProcessorOptions(trace.WithMaxExportBatchSize(512)),
	)
}

func ExampleWithLogOptions() {
	providerconfighttp.New(
		providerconfighttp.WithLogOptions(otlploghttp.WithCompression(otlploghttp.GzipCompression)),
	)
}

func ExampleWithBatchProcessorOptions() {
	providerconfighttp.New(
		providerconfighttp.WithBatchProcessorOptions(log.WithExportMaxBatchSize(512)),
	)
}
