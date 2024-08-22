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

package otelslog_test

import (
	"context"
	"errors"
	"github.com/vincentfree/opentelemetry/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log/slog"
	"math"
	"os"
)

func ExampleAddTracingContext() {
	tracer := otel.Tracer("otelslog/example")
	_, span := tracer.Start(context.Background(), "example-span")
	logger := otelslog.New(otelslog.WithProvidedHandler(slog.NewJSONHandler(os.Stdout, timeRemoved)))
	// pass span to AddTracingContext
	logger.LogAttrs(nil, slog.LevelInfo, "in case of a success", otelslog.AddTracingContext(span)...)

	// or in the case of an error
	err := errors.New("example error")
	logger.LogAttrs(nil, slog.LevelError, "in case of a failure", otelslog.AddTracingContext(span, err)...)
	// Output: {"level":"INFO","msg":"in case of a success","traceID":"00000000000000000000000000000000","spanID":"0000000000000000"}
	// {"level":"ERROR","msg":"in case of a failure","error":"example error","traceID":"00000000000000000000000000000000","spanID":"0000000000000000"}
}

func ExampleAddTracingContextWithAttributes() {
	tracer := otel.Tracer("otelslog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	attributes := []attribute.KeyValue{
		attribute.String("exampleKey", "exampleValue"),
		attribute.Bool("isValid", true),
	}
	logger := otelslog.New(otelslog.WithProvidedHandler(slog.NewJSONHandler(os.Stdout, timeRemoved)))

	// pass span to AddTracingContext
	logger.LogAttrs(nil, slog.LevelInfo, "in case of a success", otelslog.AddTracingContextWithAttributes(span, attributes)...)

	// or in the case of an error
	err := errors.New("example error")
	logger.LogAttrs(nil, slog.LevelError, "in case of a failure", otelslog.AddTracingContextWithAttributes(span, attributes, err)...)
	// Output: {"level":"INFO","msg":"in case of a success","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.exampleKey":"exampleValue","trace.attribute.isValid":true}
	// {"level":"ERROR","msg":"in case of a failure","error":"example error","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.exampleKey":"exampleValue","trace.attribute.isValid":true}
}

func ExampleWithAttributes() {
	logger := otelslog.New(
		otelslog.WithAttributes(attribute.String("test", "value"), attribute.Bool("isValid", true)),
		otelslog.WithProvidedHandler(slog.NewTextHandler(os.Stdout, timeRemoved)),
	)

	tracer := otel.Tracer("otelslog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	// pass span to AddTracingContext
	logger.LogAttrs(nil, slog.LevelInfo, "in case of a success", otelslog.AddTracingContext(span)...)

	// or in the case of an error
	err := errors.New("example error")
	logger.LogAttrs(nil, slog.LevelError, "in case of a failure", otelslog.AddTracingContext(span, err)...)
	// Output: level=INFO msg="in case of a success" traceID=00000000000000000000000000000000 spanID=0000000000000000 trace.attribute.test=value trace.attribute.isValid=true
	// level=ERROR msg="in case of a failure" error="example error" traceID=00000000000000000000000000000000 spanID=0000000000000000 trace.attribute.test=value trace.attribute.isValid=true
}

func ExampleWithAttributePrefix() {
	otelslog.New(otelslog.WithAttributePrefix("prefix"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleWithServiceName() {
	otelslog.New(otelslog.WithServiceName("example-service"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleWithSpanID() {
	otelslog.New(otelslog.WithSpanID("span-id"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleWithTraceID() {
	otelslog.New(otelslog.WithTraceID("trace-id"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleSetLogOptions() {
	option := otelslog.WithAttributes(attribute.String("test", "value"), attribute.Bool("isValid", true))
	// use of SetLogOptions
	otelslog.SetLogOptions(option)

	// set up tracer
	tracer := otel.Tracer("otelslog/example")
	_, span := tracer.Start(context.Background(), "example-span")
	// pass span to AddTracingContext
	slog.LogAttrs(nil, slog.LevelInfo, "in case of a success", otelslog.AddTracingContext(span)...)

	// or in the case of an error
	err := errors.New("example error")
	slog.LogAttrs(nil, slog.LevelError, "in case of a failure", otelslog.AddTracingContext(span, err)...)
}

func ExampleLogger_WithTracingContext() {
	tracer := otel.Tracer("otelslog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	logger := otelslog.New(otelslog.WithProvidedHandler(slog.NewTextHandler(os.Stdout, timeRemoved)))
	// pass span to AddTracingContext
	logger.WithTracingContext(nil, slog.LevelInfo, "in case of a success", span, nil)
	slog.LogAttrs(nil, slog.LevelInfo, "in case of a success", otelslog.AddTracingContext(span)...)

	// or in the case of an error
	err := errors.New("example error")
	logger.WithTracingContext(nil, slog.LevelError, "in case of a failure", span, err)
	// Output: level=INFO msg="in case of a success" traceID=00000000000000000000000000000000 spanID=0000000000000000
	// level=ERROR msg="in case of a failure" error="example error" traceID=00000000000000000000000000000000 spanID=0000000000000000
}

func ExampleLogger_WithTracingContextAndAttributes() {
	attributes := []attribute.KeyValue{
		attribute.String("exampleKey", "exampleValue"),
		attribute.Bool("isValid", true),
	}

	tracer := otel.Tracer("otelslog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	logger := otelslog.New(otelslog.WithProvidedHandler(slog.NewJSONHandler(os.Stdout, timeRemoved)))
	// pass span to AddTracingContext
	logger.WithTracingContextAndAttributes(nil, slog.LevelInfo, "in case of a success", span, nil, attributes)

	// or in the case of an error
	err := errors.New("example error")
	logger.WithTracingContextAndAttributes(nil, slog.LevelError, "in case of a failure", span, err, attributes)
	// Output: {"level":"INFO","msg":"in case of a success","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.exampleKey":"exampleValue","trace.attribute.isValid":true}
	// {"level":"ERROR","msg":"in case of a failure","error":"example error","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.exampleKey":"exampleValue","trace.attribute.isValid":true}
}

func ExampleNew() {
	tracer := otel.Tracer("otelslog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	logger := otelslog.New(otelslog.WithProvidedHandler(slog.NewJSONHandler(os.Stdout, timeRemoved)))
	// pass span to AddTracingContext
	logger.WithTracingContext(nil, slog.LevelInfo, "in case of a success", span, nil)
	// Output: {"level":"INFO","msg":"in case of a success","traceID":"00000000000000000000000000000000","spanID":"0000000000000000"}
}

func ExampleNewWithHandler() {
	tracer := otel.Tracer("otelslog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	logger := otelslog.NewWithHandler(slog.NewTextHandler(os.Stdout, nil))
	// pass span to AddTracingContext
	logger.WithTracingContext(nil, slog.LevelInfo, "in case of a success", span, nil)
}

var timeRemoved = &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
	if a.Key == "time" {
		return slog.Attr{}
	}
	return a
}}

func ExampleConvertToSlogFormat() {
	logger := otelslog.New(otelslog.WithProvidedHandler(slog.NewJSONHandler(os.Stdout, timeRemoved)))
	attributes := []attribute.KeyValue{
		attribute.String("stringExample", "this is an example string"),
		attribute.Float64("float64Example", 42.0),
		attribute.Int64("int64Example", 42),
		attribute.Bool("boolExample", true),
		attribute.BoolSlice("boolSliceExample", []bool{true, false, true}),
		attribute.Int64Slice("int64SliceExample", []int64{42, math.MaxInt64}),
		attribute.Float64Slice("float64SliceExample", []float64{42.0, math.Pi}),
		attribute.StringSlice("stringSliceExample", []string{"test", "values"}),
	}
	attrs := []slog.Attr{slog.String("init", "attr")}
	attrs = append(attrs, logger.ConvertToSlogFormat(attributes)...)
	logger.LogAttrs(nil, slog.LevelInfo, "test", attrs...)
	// Output: {"level":"INFO","msg":"test","init":"attr","trace.attribute.stringExample":"this is an example string","trace.attribute.float64Example":42,"trace.attribute.int64Example":42,"trace.attribute.boolExample":true,"trace.attribute.boolSliceExample.0":true,"trace.attribute.boolSliceExample.1":false,"trace.attribute.boolSliceExample.2":true,"trace.attribute.int64SliceExample.0":42,"trace.attribute.int64SliceExample.1":9223372036854775807,"trace.attribute.float64SliceExample.0":42,"trace.attribute.float64SliceExample.1":3.141592653589793,"trace.attribute.stringSliceExample.0":"test","trace.attribute.stringSliceExample.1":"values"}
}
