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
	"log/slog"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	protocolReg = regexp.MustCompile("^https?://")
	endpointReg = regexp.MustCompile("^(https?://.+:\\d{2,5}|.+:\\d{2,5})$")
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

// WithInsecure appends the WithInsecure options for the trace, metric and log providers.
func WithInsecure() Option {
	return func(c *grpcConfig) {
		c.traceOptions = append(c.traceOptions, otlptracegrpc.WithInsecure())
		c.logOptions = append(c.logOptions, otlploggrpc.WithInsecure())
		c.metricOptions = append(c.metricOptions, otlpmetricgrpc.WithInsecure())
	}
}

// WithCollectorEndpoint handled by providerconfig.New
// Has no effect when WithGRPCConn is set for any of the Options
func WithCollectorEndpoint(endpoint string) Option {
	err := validateEndpoint(endpoint)
	if err != nil {
		panic(fmt.Errorf("endpoint did not pass validation. %w", err))
	}

	return func(gc *grpcConfig) {
		if len(gc.traceOptions) == 0 {
			gc.traceOptions = []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(endpoint)}
		} else {
			gc.traceOptions = append(gc.traceOptions, otlptracegrpc.WithEndpoint(endpoint))
		}

		if len(gc.metricOptions) == 0 {
			gc.metricOptions = []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(endpoint)}
		} else {
			gc.metricOptions = append(gc.metricOptions, otlpmetricgrpc.WithEndpoint(endpoint))
		}

		if len(gc.logOptions) == 0 {
			gc.logOptions = []otlploggrpc.Option{otlploggrpc.WithEndpoint(endpoint)}
		} else {
			gc.logOptions = append(gc.logOptions, otlploggrpc.WithEndpoint(endpoint))
		}
	}
}

func validateEndpoint(endpoint string) error {
	if !endpointReg.MatchString(endpoint) {
		return fmt.Errorf("invalid endpoint: %s", endpoint)
	}

	if strings.Contains(endpoint, "://") {
		x := strings.Split(endpoint, "://")
		switch strings.ToLower(x[0]) {
		case "https", "http":
			break
		default:
			return fmt.Errorf("invalid endpoint: %s, protocol was invalid: %s", endpoint, x[0])
		}
	}

	ne := protocolReg.ReplaceAllString(endpoint, "")
	i := strings.LastIndex(ne, ":")
	host := ne[:i]
	p := ne[i+1:]
	port, err := strconv.Atoi(p)
	if err != nil {
		return err
	}
	if port > math.MaxUint16 || port <= 0 {
		logger.Error("invalid port value", slog.Uint64("port_range_max", math.MaxUint16), slog.Uint64("port_range_min", 0), slog.Int("port_provided", port))
		return fmt.Errorf("invalid port value: %d", port)
	}

	logger.Debug("provided collector endpoint", slog.String("provided_host", host), slog.Int("provided_port", port))
	return nil
}
