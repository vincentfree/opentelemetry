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
	"log/slog"
	"math"
	"regexp"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

var (
	protocolReg = regexp.MustCompile("^https?://")
	endpointReg = regexp.MustCompile("^(https?://\\w+:\\d{2,5}|\\w+:\\d{2,5})$")
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
func WithCollectorEndpoint(endpoint string) Option {
	err := validateEndpoint(endpoint)
	if err != nil {
		panic(fmt.Errorf("endpoint did not pass validation. %w", err))
	}

	return func(gc *httpConfig) {
		if len(gc.traceOptions) == 0 {
			gc.traceOptions = []otlptracehttp.Option{otlptracehttp.WithEndpoint(endpoint)}
		} else {
			gc.traceOptions = append(gc.traceOptions, otlptracehttp.WithEndpoint(endpoint))
		}

		if len(gc.metricOptions) == 0 {
			gc.metricOptions = []otlpmetrichttp.Option{otlpmetrichttp.WithEndpoint(endpoint)}
		} else {
			gc.metricOptions = append(gc.metricOptions, otlpmetrichttp.WithEndpoint(endpoint))
		}

		if len(gc.logOptions) == 0 {
			gc.logOptions = []otlploghttp.Option{otlploghttp.WithEndpoint(endpoint)}
		} else {
			gc.logOptions = append(gc.logOptions, otlploghttp.WithEndpoint(endpoint))
		}
	}
}

// WithInsecure appends the WithInsecure options for the trace, metric and log providers.
func WithInsecure() Option {
	return func(gc *httpConfig) {
		gc.traceOptions = append(gc.traceOptions, otlptracehttp.WithInsecure())
		gc.metricOptions = append(gc.metricOptions, otlpmetrichttp.WithInsecure())
		gc.logOptions = append(gc.logOptions, otlploghttp.WithInsecure())
	}
}

func validateEndpoint(endpoint string) error {
	if !endpointReg.MatchString(endpoint) {
		return fmt.Errorf("invalid endpoint: %s", endpoint)
	}
	ne := protocolReg.ReplaceAllString(endpoint, "")
	hostPort := strings.SplitN(ne, ":", 2)
	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		return err
	}
	if port > math.MaxUint16 || port <= 0 {
		logger.Error("invalid port value", slog.Uint64("port_range_max", math.MaxUint16), slog.Uint64("port_range_min", 0), slog.Int("port_provided", port))
		return fmt.Errorf("invalid port value: %d", port)
	}

	logger.Debug("provided collector endpoint", slog.String("provided_host", hostPort[0]), slog.Int("provided_port", port))
	return nil
}
