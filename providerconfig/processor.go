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

import (
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type SignalProcessor interface {
	AsyncTraceProcessor(...trace.BatchSpanProcessorOption) trace.SpanProcessor
	SyncTraceProcessor() trace.SpanProcessor
	AsyncLogProcessor(...log.BatchProcessorOption) log.Processor
	SyncLogProcessor(...log.SimpleProcessorOption) log.Processor
	MetricProcessor(...metric.PeriodicReaderOption) metric.Reader
}
