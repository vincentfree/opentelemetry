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

package providerconfiggrpc_test

import (
	"github.com/vincentfree/opentelemetry/providerconfig"
	"github.com/vincentfree/opentelemetry/providerconfig/providerconfiggrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

func ExampleNew() {
	signalProcessor := providerconfiggrpc.New(
		providerconfiggrpc.WithCollectorEndpoint("0.0.0.0", 9898),
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
	providerconfiggrpc.New(
		providerconfiggrpc.WithCollectorEndpoint("0.0.0.0", 9898),
	)
}

func ExampleWithSpanProcessorOptions() {
	providerconfiggrpc.New(
		providerconfiggrpc.WithSpanProcessorOptions(trace.WithMaxExportBatchSize(512)),
	)
}

func ExampleWithBatchProcessorOptions() {
	providerconfiggrpc.New(
		providerconfiggrpc.WithBatchProcessorOptions(log.WithExportMaxBatchSize(512)),
	)
}

func ExampleWithSimpleProcessorOptions() {
	providerconfiggrpc.New(
		providerconfiggrpc.WithSimpleProcessorOptions( /*currently there are no options for the simpleProcessor*/ ),
	)
}

func ExampleWithPeriodicReaderOptions() {
	providerconfiggrpc.New(
		providerconfiggrpc.WithPeriodicReaderOptions(metric.WithInterval(5 * time.Second)),
	)
}

func ExampleWithTraceOptions() {
	providerconfiggrpc.New(
		providerconfiggrpc.WithTraceOptions(otlptracegrpc.WithCompressor("gzip")),
	)
}

func ExampleWithLogOptions() {
	providerconfiggrpc.New(
		providerconfiggrpc.WithLogOptions(
			otlploggrpc.WithCompressor("gzip"),
			otlploggrpc.WithTimeout(5*time.Second),
		),
	)
}

func ExampleWithMetricOptions() {
	providerconfiggrpc.New(
		providerconfiggrpc.WithMetricOptions(otlpmetricgrpc.WithCompressor("gzip")),
	)
}
