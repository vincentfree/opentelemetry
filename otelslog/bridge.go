// Copyright 2024 Vincent Free <vincentfree@outlook.com>
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

package otelslog

import (
	otelslogger "go.opentelemetry.io/contrib/bridges/otelslog"
	"log/slog"
)

func WithOtelBridge(name string, options ...otelslogger.Option) LogOption {
	return func(c *logConfig) {
		c.handler = otelslogger.NewHandler(name, options...)
	}
}

func WithOtelBridgeDisabled() LogOption {
	return func(c *logConfig) {
		c.bridgeDisabled = true
	}
}

// WithProvidedHandler is a LogOption that sets overwrites the otel based slog.Handler
// and uses the provided slog.Handler instead.
func WithProvidedHandler(handler slog.Handler) LogOption {
	return func(c *logConfig) {
		c.overwriteHandler = handler
	}
}
