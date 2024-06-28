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
	"fmt"
	"github.com/vincentfree/opentelemetry/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log/slog"
	"math"
	"os"
)

func ExampleAddTracingContext() {
	ctx := context.TODO()
	tracer := otel.Tracer("otelslog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	// pass span to AddTracingContext
	slog.LogAttrs(ctx, slog.LevelInfo, "in case of a success", otelslog.AddTracingContext(span)...)

	// or in the case of an error
	err := errors.New("example error")
	slog.LogAttrs(ctx, slog.LevelError, "in case of a failure", otelslog.AddTracingContext(span, err)...)
}

func ExampleAddTracingContextWithAttributes() {
	ctx := context.TODO()
	tracer := otel.Tracer("otelslog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	attributes := []attribute.KeyValue{
		attribute.String("exampleKey", "exampleValue"),
		attribute.Bool("isValid", true),
	}

	// pass span to AddTracingContext
	slog.LogAttrs(ctx, slog.LevelInfo, "in case of a success", otelslog.AddTracingContextWithAttributes(span, attributes)...)

	// or in the case of an error
	err := errors.New("example error")
	slog.LogAttrs(ctx, slog.LevelError, "in case of a failure", otelslog.AddTracingContextWithAttributes(span, attributes, err)...)
}

func ExampleWithAttributes() {
	option := otelslog.WithAttributes(attribute.String("test", "value"), attribute.Bool("isValid", true))
	otelslog.SetLogOptions(option)

	tracer := otel.Tracer("otelslog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	// pass span to AddTracingContext
	slog.LogAttrs(nil, slog.LevelInfo, "in case of a success", otelslog.AddTracingContext(span)...)

	// or in the case of an error
	err := errors.New("example error")
	slog.LogAttrs(nil, slog.LevelError, "in case of a failure", otelslog.AddTracingContext(span, err)...)
}

func ExampleWithAttributePrefix() {
	otelslog.SetLogOptions(otelslog.WithAttributePrefix("prefix"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleWithServiceName() {
	otelslog.SetLogOptions(otelslog.WithServiceName("example-service"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleWithSpanID() {
	otelslog.SetLogOptions(otelslog.WithSpanID("span-id"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleWithTraceID() {
	otelslog.SetLogOptions(otelslog.WithTraceID("trace-id"))
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

	logger := otelslog.New()
	// pass span to AddTracingContext
	logger.WithTracingContext(nil, slog.LevelInfo, "in case of a success", span, nil)
	slog.LogAttrs(nil, slog.LevelInfo, "in case of a success", otelslog.AddTracingContext(span)...)

	// or in the case of an error
	err := errors.New("example error")
	logger.WithTracingContext(nil, slog.LevelError, "in case of a failure", span, err)
}

func ExampleLogger_WithTracingContextAndAttributes() {
	attributes := []attribute.KeyValue{
		attribute.String("exampleKey", "exampleValue"),
		attribute.Bool("isValid", true),
	}

	tracer := otel.Tracer("otelslog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	logger := otelslog.New()
	// pass span to AddTracingContext
	logger.WithTracingContextAndAttributes(nil, slog.LevelInfo, "in case of a success", span, nil, attributes)

	// or in the case of an error
	err := errors.New("example error")
	logger.WithTracingContextAndAttributes(nil, slog.LevelError, "in case of a failure", span, err, attributes)
}

func ExampleNew() {
	tracer := otel.Tracer("otelslog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	logger := otelslog.New()
	// pass span to AddTracingContext
	logger.WithTracingContext(nil, slog.LevelInfo, "in case of a success", span, nil)
}

func ExampleNewWithHandler() {
	tracer := otel.Tracer("otelslog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	logger := otelslog.NewWithHandler(slog.NewTextHandler(os.Stdout, nil))
	// pass span to AddTracingContext
	logger.WithTracingContext(nil, slog.LevelInfo, "in case of a success", span, nil)
}

func ExampleConvertToSlogFormat() {
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
	attrs = otelslog.ConvertToSlogFormat(attributes, attrs)
	fmt.Println(attrs)
}
