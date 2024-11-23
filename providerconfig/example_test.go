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

package providerconfig_test

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vincentfree/opentelemetry/providerconfig"
	"github.com/vincentfree/opentelemetry/providerconfig/providerconfignoop"
	prombridge "go.opentelemetry.io/contrib/bridges/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"time"
)

// initializing through main and manually setting the providers.
//
// providers can also be set directly through providerconfig.Option's like:
// providerconfig.WithInitTraces()
// providerconfig.WithInitMetrics()
// providerconfig.WithInitLogs()
func ExampleNew() {
	provider := providerconfig.New(
		providerconfig.WithApplicationName("example-app"),
		providerconfig.WithApplicationVersion("0.1.0"),
		providerconfig.WithExecutionType(providerconfig.Async),
		providerconfig.WithSignalProcessor(providerconfignoop.NewNoopProcessor()),
	)

	// traces
	otel.SetTracerProvider(provider.TraceProvider())

	// metrics
	otel.SetMeterProvider(provider.MetricProvider())

	// logs
	global.SetLoggerProvider(provider.LogProvider())
}

func ExampleWithApplicationName() {
	providerconfig.New(
		providerconfig.WithApplicationName("example-app"),
	)
}

func ExampleWithApplicationVersion() {
	providerconfig.New(
		providerconfig.WithApplicationVersion("1.0.0"),
	)
}

func ExampleWithExecutionType() {
	providerconfig.New(
		providerconfig.WithExecutionType(providerconfig.Async),
	)
}

func ExampleWithSignalProcessor() {
	providerconfig.New(
		// Example processor, in a real scenario, the http or grpc processors should be used.
		// Either is a separate import that needs to be added to the modules
		providerconfig.WithSignalProcessor(providerconfignoop.NewNoopProcessor()),
	)
}

func ExampleWithResourceOptions() {
	providerconfig.New(
		providerconfig.WithResourceOptions(resource.WithContainer(),
			resource.WithHost(),
		),
	)
}

func ExampleWithTraceProviderOptions() {
	providerconfig.New(
		providerconfig.WithTraceProviderOptions(sdktrace.WithSampler(sdktrace.AlwaysSample())),
	)
}

func ExampleWithLogProviderOptions() {
	providerconfig.New(
		providerconfig.WithLogProviderOptions(sdklog.WithAttributeCountLimit(15)),
	)
}

func ExampleWithMetricProviderOptions() {
	providerconfig.New(
		providerconfig.WithPeriodicReaderOptions(sdkmetric.WithTimeout(30 * time.Second)),
	)
}

func ExampleWithDisabledSignals() {
	providerconfig.New(
		providerconfig.WithDisabledSignals(
			false, // traces
			true,  // metrics
			true,  // logs
		),
	)
}

func ExampleWithTracePropagator() {
	providerconfig.New(
		providerconfig.WithTracePropagator(
			propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			),
		),
	)
}

func ExampleWithPrometheusBridge() {
	providerconfig.New(
		providerconfig.WithApplicationName("example-app"),
		providerconfig.WithApplicationVersion("0.1.0"),
		providerconfig.WithExecutionType(providerconfig.Async),
		providerconfig.WithSignalProcessor(providerconfignoop.NewNoopProcessor()),
		// options come from the opentelemetry bridge(renamed to prombridge) library
		// (naming conflicts with the prometheus library)
		providerconfig.WithPrometheusBridge(prombridge.WithGatherer(prometheus.NewRegistry())),
	)
}

func ExampleProvider_ShutdownAll() {
	provider := providerconfig.New(
		providerconfig.WithApplicationName("example-app"),
		providerconfig.WithApplicationVersion("0.0.1"),
	)

	// shutdown all the otel providers
	provider.ShutdownAll()
}

func ExampleProvider_ShutdownByType() {
	provider := providerconfig.New(
		providerconfig.WithApplicationName("example-app"),
		providerconfig.WithApplicationVersion("0.0.1"),
	)

	// shutdown just traces
	provider.ShutdownByType(providerconfig.TraceHook)

	// shutdown just metrics
	provider.ShutdownByType(providerconfig.MetricHook)

	// shutdown just logs
	provider.ShutdownByType(providerconfig.LogHook)
}
