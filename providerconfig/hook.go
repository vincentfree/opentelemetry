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

package providerconfig

type SignalHookName string

const (
	TraceHook  SignalHookName = "trace"
	MetricHook SignalHookName = "metric"
	LogHook    SignalHookName = "log"
)

type ShutdownHook func()
type ShutdownHooks map[SignalHookName]ShutdownHook

func NewShutdownHooks(fns ...func() (SignalHookName, ShutdownHook)) ShutdownHooks {
	sdh := make(ShutdownHooks, len(fns))
	for _, fn := range fns {
		name, hook := fn()
		sdh[name] = hook
	}
	return sdh
}

func ShutDownPair(name SignalHookName, fn ShutdownHook) func() (SignalHookName, ShutdownHook) {
	return func() (SignalHookName, ShutdownHook) {
		return name, fn
	}
}

func (h ShutdownHooks) ShutdownAll() {
	for _, hook := range h {
		hook()
	}
}

func (h ShutdownHooks) ShutdownByType(hookType SignalHookName) bool {
	if hook, exists := h[hookType]; exists {
		hook()
		return true
	}
	return false
}
