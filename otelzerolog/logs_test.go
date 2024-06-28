// Copyright 2023 Vincent Free
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
	"context"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/exp/constraints"
	"io"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestSetLogOptions(t *testing.T) {
	serviceName := "test-service"
	SetLogOptions(WithServiceName(serviceName))
	_, span := otel.Tracer("test").Start(context.Background(), serviceName)

	out := captureLog(t, func(logger zerolog.Logger) {
		logger.Info().Func(AddTracingContext(span)).Msg("test")
	})

	data := logToMap(t, out)

	idCheck(t, "traceID", data["traceID"], 32)
	idCheck(t, "spanID", data["spanID"], 16)

	if v, ok := data["service.name"].(string); ok {
		if v != serviceName {
			t.Errorf("expected %s, but was %s", serviceName, v)
		}
	}
}

func TestWithSpanID(t *testing.T) {
	id := "testSpanID"
	// given a new span ID name
	SetLogOptions(WithSpanID(id))

	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")
	// when a log with AddTracingContext is preformed
	out := captureLog(t, func(logger zerolog.Logger) {
		logger.Info().Func(AddTracingContext(span)).Msg("test")
	})

	data := logToMap(t, out)
	if _, ok := data[id]; !ok {
		t.Errorf("the log should have had a overwritten spanID but the field %s was not found. data: %s", id, data)
	}
}

func TestWithTraceID(t *testing.T) {
	id := "testTraceID"
	// given a new span ID name
	SetLogOptions(WithTraceID(id))

	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")
	// when a log with AddTracingContext is preformed
	out := captureLog(t, func(logger zerolog.Logger) {
		logger.Info().Func(AddTracingContext(span)).Msg("test")
	})

	data := logToMap(t, out)
	if _, ok := data[id]; !ok {
		t.Errorf("the log should have had a overwritten traceID but the field %s was not found. data: %s", id, data)
	}
}

func TestWithAttributes(t *testing.T) {
	SetLogOptions(WithAttributes(attribute.String("test", "value"), attribute.Bool("isValid", true)))

	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")
	// when a log with AddTracingContext is preformed
	out := captureLog(t, func(logger zerolog.Logger) {
		logger.Info().Func(AddTracingContext(span)).Msg("test")
	})

	data := logToMap(t, out)

	if v, ok := data["trace.attribute.test"].(string); !ok {
		t.Errorf("attribute was not found")
	} else {
		if v != "value" {
			t.Error("the value of the attribute did not match")
		}
	}
	if v, ok := data["trace.attribute.isValid"].(bool); !ok {
		t.Errorf("attribute was not found")
	} else {
		if !v {
			t.Error("the value of the attribute did not match")
		}
	}
}

func TestWithAttributePrefix(t *testing.T) {
	SetLogOptions(WithAttributes(attribute.String("test", "value")), WithAttributePrefix("testing"))

	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")
	// when a log with AddTracingContext is preformed
	out := captureLog(t, func(logger zerolog.Logger) {
		logger.Info().Func(AddTracingContext(span)).Msg("test")
	})

	data := logToMap(t, out)

	if v, ok := data["testing.test"].(string); !ok {
		t.Errorf("attribute was not found")
	} else {
		if v != "value" {
			t.Error("the value of the attribute did not match")
		}
	}

	// reset prefix
	SetLogOptions(WithAttributePrefix("trace.attribute"))
}

func TestAddTracingContextWithAttributes(t *testing.T) {
	SetLogOptions(WithAttributes(attribute.String("test", "value")))
	localAttributes := []attribute.KeyValue{
		attribute.Float64("localFloat64", 42.0),
		attribute.Int64("localInt64", 42),
		attribute.BoolSlice("localBoolSlice", []bool{true}),
		attribute.Int64Slice("localInt64Slice", []int64{42}),
		attribute.Float64Slice("localFloat64Slice", []float64{42.0}),
		attribute.StringSlice("localStringSlice", []string{"test"}),
	}
	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")
	// when a log with AddTracingContext is preformed
	out := captureLog(t, func(logger zerolog.Logger) {
		logger.Info().Func(AddTracingContextWithAttributes(span, localAttributes)).Msg("test")
	})

	data := logToMap(t, out)

	if v, ok := data["trace.attribute.test"].(string); !ok {
		t.Errorf("attribute was not found")
	} else {
		if v != "value" {
			t.Error("the value of the attribute did not match")
		}
	}

	attributeCheck(t, data["trace.attribute.localFloat64"], 42.0)

	// although the function inject an int64 in the map it's seen as a float64
	attributeCheck(t, data["trace.attribute.localInt64"], 42.0)

	attributeKeyCheck(t, data, "trace.attribute.localBoolSlice")
	attributeKeyCheck(t, data, "trace.attribute.localInt64Slice")
	attributeKeyCheck(t, data, "trace.attribute.localFloat64Slice")
	attributeKeyCheck(t, data, "trace.attribute.localStringSlice")
}

func FuzzAddTracingContextWithAttributes(f *testing.F) {
	f.Add("value", 42.0, int64(42), true, int64(42), 42.0, "test", "value")
	f.Add("test_value", 2.2, int64(4242), false, int64(10), 2.1001, "15", "@$#!")

	f.Fuzz(func(t *testing.T, s string, lf float64, li int64, bs bool, is int64, fs float64, s1 string, s2 string) {
		SetLogOptions(WithAttributes(attribute.String("test", s)))
		localAttributes := []attribute.KeyValue{
			attribute.Float64("localFloat64", lf),
			attribute.Int64("localInt64", li),
			attribute.BoolSlice("localBoolSlice", []bool{bs}),
			attribute.Int64Slice("localInt64Slice", []int64{is, li}),
			attribute.Float64Slice("localFloat64Slice", []float64{fs}),
			attribute.StringSlice("localStringSlice", []string{s1, s2, s}),
		}

		_, span := otel.Tracer("test"+strconv.Itoa(int(rand.Int31n(10000000)))).Start(context.Background(), "serviceName")
		log.Info().Func(AddTracingContextWithAttributes(span, localAttributes)).Msg("test")
		// when a log with AddTracingContext is preformed
		//_ = captureLog(t, func(logger zerolog.Logger) {
		//})
		span.End()

		//data := logToMap(t, out)

		//if v, ok := data["trace.attribute.test"].(string); !ok {
		//	t.Errorf("attribute was not found")
		//} else {
		//	if v != s {
		//		t.Error("the value of the attribute did not match")
		//	}
		//}

		//attributeCheck(t, data["trace.attribute.localFloat64"], lf)

		// although the function inject an int64 in the map it's seen as a float64
		//attributeCheck(t, data["trace.attribute.localInt64"], li)

		//attributeKeyCheck(t, data, "trace.attribute.test")
		//attributeKeyCheck(t, data, "trace.attribute.localFloat64")
		//attributeKeyCheck(t, data, "trace.attribute.localInt64")
		//attributeKeyCheck(t, data, "trace.attribute.localBoolSlice")
		//attributeKeyCheck(t, data, "trace.attribute.localInt64Slice")
		//attributeKeyCheck(t, data, "trace.attribute.localFloat64Slice")
		//attributeKeyCheck(t, data, "trace.attribute.localStringSlice")

	})
}

func TestLogWithError(t *testing.T) {
	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")
	// when a log with AddTracingContext is preformed
	err := errors.New("error")
	out := captureLog(t, func(logger zerolog.Logger) {
		logger.Info().Func(AddTracingContext(span, err)).Msg("test")
	})

	data := logToMap(t, out)
	if _, ok := data["error"]; !ok {
		t.Errorf("the error was to injected int the log, msg: %s", data)
	}

}

func attributeCheck[T constraints.Ordered](t *testing.T, data any, checkValue T) {
	if v, ok := data.(T); !ok {
		t.Errorf("attribute was not found, incoming field %s", data)
	} else {
		if v != checkValue {
			t.Error("the value of the attribute did not match")
		}
	}
}

func attributeKeyCheck(t *testing.T, data map[string]any, field string) {
	if _, ok := data[field]; !ok {
		t.Errorf("attribute was not found")
	}
}

func captureLog(t *testing.T, fn func(logger zerolog.Logger)) []byte {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	logger := zerolog.New(w)
	fn(logger)
	err := w.Close()
	if err != nil {
		t.Errorf("should not fail while closing the Pipe file")
	}
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	return out
}

func logToMap(t *testing.T, out []byte) map[string]any {
	var data map[string]any
	err := json.Unmarshal(out, &data)
	if err != nil {
		t.Error("unable to unmarshal json log into map")
	}
	return data
}

func idCheck(t *testing.T, name string, value any, length int) {
	if v, ok := value.(string); ok {
		if len(v) != length {
			t.Errorf("%s should be %d log but was: %d", name, length, len(v))
		}

	} else {
		t.Errorf("%s should be in the log", name)
	}
}
