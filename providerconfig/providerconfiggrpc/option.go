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

package providerconfiggrpc

import (
	"fmt"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type grpcConfig struct {
	traceOptions           []otlptracegrpc.Option
	metricOptions          []otlpmetricgrpc.Option
	logOptions             []otlploggrpc.Option
	batchProcessorOptions  []log.BatchProcessorOption
	simpleProcessorOptions []log.SimpleProcessorOption
	periodicReaderOptions  []metric.PeriodicReaderOption
	spanProcessorOptions   []trace.BatchSpanProcessorOption
}

type Option func(*grpcConfig)

func WithTraceOptions(options ...otlptracegrpc.Option) Option {
	return func(gc *grpcConfig) {
		if gc.traceOptions == nil || len(gc.traceOptions) == 0 {
			gc.traceOptions = options
		} else {
			gc.traceOptions = append(gc.traceOptions, options...)
		}
	}
}

func WithMetricOptions(options ...otlpmetricgrpc.Option) Option {
	return func(gc *grpcConfig) {
		if gc.metricOptions == nil || len(gc.metricOptions) == 0 {
			gc.metricOptions = options
		} else {
			gc.metricOptions = append(gc.metricOptions, options...)
		}
	}
}

func WithLogOptions(options ...otlploggrpc.Option) Option {
	return func(gc *grpcConfig) {
		if gc.logOptions == nil || len(gc.logOptions) == 0 {
			gc.logOptions = options
		} else {
			gc.logOptions = append(gc.logOptions, options...)
		}
	}
}

func WithSpanProcessorOptions(options ...trace.BatchSpanProcessorOption) Option {
	return func(c *grpcConfig) {
		c.spanProcessorOptions = options
	}
}

func WithPeriodicReaderOptions(options ...metric.PeriodicReaderOption) Option {
	return func(c *grpcConfig) {
		c.periodicReaderOptions = options
	}
}

func WithSimpleProcessorOptions(options ...log.SimpleProcessorOption) Option {
	return func(c *grpcConfig) {
		c.simpleProcessorOptions = options
	}
}

func WithBatchProcessorOptions(options ...log.BatchProcessorOption) Option {
	return func(c *grpcConfig) {
		c.batchProcessorOptions = options
	}
}

// WithCollectorEndpoint handled by providerconfig.New
// Has no effect when WithGRPCConn is set for any of the Options
func WithCollectorEndpoint(url string, port uint16) Option {
	return func(gc *grpcConfig) {
		if gc.traceOptions == nil || len(gc.traceOptions) == 0 {
			gc.traceOptions = []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%d", url, port))}
		} else {
			gc.traceOptions = append(gc.traceOptions, otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%d", url, port)))
		}

		if gc.metricOptions == nil || len(gc.metricOptions) == 0 {
			gc.metricOptions = []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(fmt.Sprintf("%s:%d", url, port))}
		} else {
			gc.metricOptions = append(gc.metricOptions, otlpmetricgrpc.WithEndpoint(fmt.Sprintf("%s:%d", url, port)))
		}

		if gc.logOptions == nil || len(gc.logOptions) == 0 {
			gc.logOptions = []otlploggrpc.Option{otlploggrpc.WithEndpoint(fmt.Sprintf("%s:%d", url, port))}
		} else {
			gc.logOptions = append(gc.logOptions, otlploggrpc.WithEndpoint(fmt.Sprintf("%s:%d", url, port)))
		}
	}
}
