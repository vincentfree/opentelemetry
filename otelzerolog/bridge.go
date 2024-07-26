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

package otelzerolog

import (
	"github.com/rs/zerolog"
	otelzlog "go.opentelemetry.io/contrib/bridges/otelzerolog"
)

// WithOtelBridge is a LogOption that sets up an OpenTelemetry bridge for logging.
// It takes a name string and optional otelzerolog options. It creates a new hook
// that gets used to create a new logger or set the global logger with the otel bridge enabled.
//
// Examples contain the noop loggerProvider,
// when using the bridge for actual use cases please initialize a log provider as described here:
// https://opentelemetry.io/docs/languages/go/instrumentation/#logs-sdk
func WithOtelBridge(name string, options ...otelzlog.Option) LogOption {
	return func(c *logConfig) {
		c.hook = otelzlog.NewHook(name, options...)
	}
}

// WithZeroLogFeatures is a LogOption that adds zeroLog features to the log config.
// It takes a variadic parameter of functions that operate on a zerolog.Context and returns a zerolog.Context.
// These functions modify the context to add additional functionality to the zeroLog logger.
// The provided features are appended to the existing zeroLog features in the log config.
// This option allows for the customization of the zeroLog logger with additional features.
func WithZeroLogFeatures(features ...func(zerolog.Context) zerolog.Context) LogOption {
	return func(c *logConfig) {
		c.zeroLogFeatures = append(c.zeroLogFeatures, features...)
	}
}

// WithOtelBridgeDisabled is a LogOption that disables the OpenTelemetry bridge for logging.
// When this option is used, the logger created or modified by SetGlobalLogger will not have the OpenTelemetry bridge enabled.
func WithOtelBridgeDisabled() LogOption {
	return func(c *logConfig) {
		c.bridgeDisabled = true
	}
}
