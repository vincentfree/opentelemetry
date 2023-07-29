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

package otellogrus

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/exp/constraints"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestSetLogOptions(t *testing.T) {
	serviceName := "test-service"
	SetLogOptions(WithServiceName(serviceName))
	_, span := otel.Tracer("test").Start(context.Background(), serviceName)

	out := captureLog(t, func(logger *logrus.Logger) {
		logger.WithFields(AddTracingContext(span)).Info("test")
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
	out := captureLog(t, func(logger *logrus.Logger) {
		logger.WithFields(AddTracingContext(span)).Info("test")
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
	out := captureLog(t, func(logger *logrus.Logger) {
		logger.WithFields(AddTracingContext(span)).Info("test")
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
	out := captureLog(t, func(logger *logrus.Logger) {
		logger.WithFields(AddTracingContext(span)).Info("test")
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
	out := captureLog(t, func(logger *logrus.Logger) {
		logger.WithFields(AddTracingContext(span)).Info("test")
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
	out := captureLog(t, func(logger *logrus.Logger) {
		logger.WithFields(AddTracingContextWithAttributes(span, localAttributes)).Info("test")
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

func TestLogWithError(t *testing.T) {
	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")
	// when a log with AddTracingContext is preformed
	err := errors.New("error")
	out := captureLog(t, func(logger *logrus.Logger) {
		logger.WithFields(AddTracingContext(span, err)).Info("test")
	})

	data := logToMap(t, out)
	if _, ok := data["error"]; !ok {
		t.Errorf("the error was to injected int the log, msg: %s", data)
	}
}

func TestWithLevel(t *testing.T) {
	logger := New(WithLevel(logrus.WarnLevel))
	if logger.Level != logrus.WarnLevel {
		t.Errorf("logger level should be set to : %s but was %s", logrus.WarnLevel, logger.Level)
	}
}

func TestWithFormatter(t *testing.T) {
	logger := New(WithFormatter(&logrus.JSONFormatter{}))
	typeName := "*logrus.JSONFormatter"
	actualType := reflect.TypeOf(logger.Formatter).String()

	if actualType != typeName {
		t.Errorf("the formatting type should be: %s, but was: %s", typeName, actualType)
	}
}

func TestLogger_WithTracingContext(t *testing.T) {
	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")
	out := captureWithLogger(t, func(logger *Logger) {
		logger.WithTracingContext(span).Info("test")
	})
	data := logToMap(t, out)

	idCheck(t, "traceID", data["traceID"], 32)
	idCheck(t, "spanID", data["spanID"], 16)
}

func TestLogger_WithTracingContextAndAttributes(t *testing.T) {
	attrs := []attribute.KeyValue{attribute.String("test", "value")}
	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")
	out := captureWithLogger(t, func(logger *Logger) {
		logger.WithTracingContextAndAttributes(span, attrs).Info("test")
	})
	data := logToMap(t, out)

	idCheck(t, "traceID", data["traceID"], 32)
	idCheck(t, "spanID", data["spanID"], 16)

	attributeKeyCheck(t, data, "trace.attribute.test")
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

func captureLog(t *testing.T, fn func(logger *logrus.Logger)) []byte {
	r, w, _ := os.Pipe()
	logger := &logrus.Logger{
		Out:          w,
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.DebugLevel,
		Formatter:    &logrus.JSONFormatter{},
		ReportCaller: false,
		BufferPool:   nil,
	}
	fn(logger)
	err := w.Close()
	if err != nil {
		t.Errorf("should not fail while closing the Pipe file")
	}
	out, _ := io.ReadAll(r)
	return out
}

func captureWithLogger(t *testing.T, fn func(logger *Logger)) []byte {
	r, w, _ := os.Pipe()
	logger := New(WithFormatter(&logrus.JSONFormatter{}))
	logger.Out = w
	fn(logger)
	err := w.Close()
	if err != nil {
		t.Errorf("should not fail while closing the Pipe file")
	}
	out, _ := io.ReadAll(r)
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
