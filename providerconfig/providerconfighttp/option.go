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

package providerconfighttp

import (
	"fmt"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type httpConfig struct {
	traceOptions           []otlptracehttp.Option
	metricOptions          []otlpmetrichttp.Option
	logOptions             []otlploghttp.Option
	batchProcessorOptions  []log.BatchProcessorOption
	simpleProcessorOptions []log.SimpleProcessorOption
	periodicReaderOptions  []metric.PeriodicReaderOption
	spanProcessorOptions   []trace.BatchSpanProcessorOption
}

type Option func(*httpConfig)

func WithTraceOptions(options ...otlptracehttp.Option) Option {
	return func(gc *httpConfig) {
		if gc.traceOptions == nil || len(gc.traceOptions) == 0 {
			gc.traceOptions = options
		} else {
			gc.traceOptions = append(gc.traceOptions, options...)
		}
	}
}

func WithMetricOptions(options ...otlpmetrichttp.Option) Option {
	return func(gc *httpConfig) {
		if gc.metricOptions == nil || len(gc.metricOptions) == 0 {
			gc.metricOptions = options
		} else {
			gc.metricOptions = append(gc.metricOptions, options...)
		}
	}
}

func WithLogOptions(options ...otlploghttp.Option) Option {
	return func(gc *httpConfig) {
		if gc.logOptions == nil || len(gc.logOptions) == 0 {
			gc.logOptions = options
		} else {
			gc.logOptions = append(gc.logOptions, options...)
		}
	}
}

func WithSpanProcessorOptions(options ...trace.BatchSpanProcessorOption) Option {
	return func(c *httpConfig) {
		c.spanProcessorOptions = options
	}
}

func WithPeriodicReaderOptions(options ...metric.PeriodicReaderOption) Option {
	return func(c *httpConfig) {
		c.periodicReaderOptions = options
	}
}

func WithSimpleProcessorOptions(options ...log.SimpleProcessorOption) Option {
	return func(c *httpConfig) {
		c.simpleProcessorOptions = options
	}
}

func WithBatchProcessorOptions(options ...log.BatchProcessorOption) Option {
	return func(c *httpConfig) {
		c.batchProcessorOptions = options
	}
}

// WithCollectorEndpoint handled by providerconfig.New
func WithCollectorEndpoint(url string, port uint16) Option {
	return func(gc *httpConfig) {
		if len(gc.traceOptions) == 0 {
			gc.traceOptions = []otlptracehttp.Option{otlptracehttp.WithEndpoint(fmt.Sprintf("%s:%d", url, port))}
		} else {
			gc.traceOptions = append(gc.traceOptions, otlptracehttp.WithEndpoint(fmt.Sprintf("%s:%d", url, port)))
		}

		if len(gc.metricOptions) == 0 {
			gc.metricOptions = []otlpmetrichttp.Option{otlpmetrichttp.WithEndpoint(fmt.Sprintf("%s:%d", url, port))}
		} else {
			gc.metricOptions = append(gc.metricOptions, otlpmetrichttp.WithEndpoint(fmt.Sprintf("%s:%d", url, port)))
		}

		if len(gc.logOptions) == 0 {
			gc.logOptions = []otlploghttp.Option{otlploghttp.WithEndpoint(fmt.Sprintf("%s:%d", url, port))}
		} else {
			gc.logOptions = append(gc.logOptions, otlploghttp.WithEndpoint(fmt.Sprintf("%s:%d", url, port)))
		}
	}
}
