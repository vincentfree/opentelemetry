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

package otelzerolog_test

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vincentfree/opentelemetry/otelzerolog"
	otelzlog "go.opentelemetry.io/contrib/bridges/otelzerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log/noop"
)

func ExampleAddTracingContext() {
	tracer := otel.Tracer("otelzerolog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	logger := otelzerolog.New()

	logger.Info().Func(otelzerolog.AddTracingContext(span)).Msg("in case of a success")
	// or in the case of an error
	err := errors.New("example error")
	logger.Error().Func(otelzerolog.AddTracingContext(span, err)).Msg("in case of a failure")

	// Output: {"level":"info","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","message":"in case of a success"}
	// {"level":"error","error":"example error","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","message":"in case of a failure"}
}

func ExampleAddTracingContextWithAttributes() {
	tracer := otel.Tracer("otelzerolog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	attributes := []attribute.KeyValue{
		attribute.String("exampleKey", "exampleValue"),
		attribute.Bool("isValid", true),
	}

	logger := otelzerolog.New()

	logger.Info().Func(otelzerolog.AddTracingContextWithAttributes(span, attributes)).Msg("in case of a success")
	// or in the case of an error
	err := errors.New("example error")
	logger.Error().Func(otelzerolog.AddTracingContextWithAttributes(span, attributes, err)).Msg("in case of a failure")

	// Output: {"level":"info","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.exampleKey":"exampleValue","trace.attribute.isValid":true,"message":"in case of a success"}
	// {"level":"error","error":"example error","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.exampleKey":"exampleValue","trace.attribute.isValid":true,"message":"in case of a failure"}
}

func ExampleWithAttributes() {
	option := otelzerolog.WithAttributes(attribute.String("test", "value"), attribute.Bool("isValid", true))
	otelzerolog.SetGlobalLogger(option)

	tracer := otel.Tracer("otelzerolog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	log.Info().Func(otelzerolog.AddTracingContext(span)).Msg("in case of a success")
	// or in the case of an error
	err := errors.New("example error")
	log.Error().Func(otelzerolog.AddTracingContext(span, err)).Msg("in case of a failure")
	// Output: {"level":"info","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.test":"value","trace.attribute.isValid":true,"message":"in case of a success"}
	// {"level":"error","error":"example error","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.test":"value","trace.attribute.isValid":true,"message":"in case of a failure"}
}

func ExampleWithAttributePrefix() {
	otelzerolog.SetGlobalLogger(otelzerolog.WithAttributePrefix("prefix"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleWithServiceName() {
	otelzerolog.SetGlobalLogger(otelzerolog.WithServiceName("example-service"))
	// use AddTracingContext or AddTracingContextWithAttributes

	tracer := otel.Tracer("otelzerolog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	log.Info().Func(otelzerolog.AddTracingContext(span)).Msg("successful message")
	// Output: {"level":"info","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","service.name":"example-service","message":"successful message"}
}

func ExampleWithSpanID() {
	otelzerolog.SetGlobalLogger(otelzerolog.WithSpanID("span-id"))
	// use AddTracingContext or AddTracingContextWithAttributes

	tracer := otel.Tracer("otelzerolog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	log.Info().Func(otelzerolog.AddTracingContext(span)).Msg("successful message")
	// Output: {"level":"info","traceID":"00000000000000000000000000000000","span-id":"0000000000000000","message":"successful message"}
}

func ExampleWithTraceID() {
	otelzerolog.SetGlobalLogger(otelzerolog.WithTraceID("trace-id"))
	// use AddTracingContext or AddTracingContextWithAttributes

	tracer := otel.Tracer("otelzerolog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	log.Info().Func(otelzerolog.AddTracingContext(span)).Msg("successful message")
	// Output: {"level":"info","trace-id":"00000000000000000000000000000000","spanID":"0000000000000000","message":"successful message"}
}

func ExampleWithZeroLogFeatures() {
	otelzerolog.SetGlobalLogger(otelzerolog.WithZeroLogFeatures(zerolog.Context.Stack))
	log.Info().Msg("successful message")
	// Output: {"level":"info","message":"successful message"}
}

func ExampleWithOtelBridgeDisabled() {
	logger := otelzerolog.New(otelzerolog.WithOtelBridgeDisabled())
	logger.Info().Msg("successful message")
	// Output: {"level":"info","message":"successful message"}
}

func ExampleWithOtelBridge() {
	logger := otelzerolog.New(otelzerolog.WithOtelBridge("example", otelzlog.WithVersion("v0.0.1"), otelzlog.WithLoggerProvider(noop.NewLoggerProvider())))
	logger.Info().Msg("successful message")
	// Output: {"level":"info","message":"successful message"}
}

func ExampleNew() {
	logger := otelzerolog.New()
	// or
	logger = otelzerolog.New(otelzerolog.WithServiceName("example-service"))
	logger.Info().Msg("successful message")
	// Output: {"level":"info","message":"successful message"}
}

func ExampleSetGlobalLogger() {
	otelzerolog.SetGlobalLogger(otelzerolog.WithServiceName("example-service"))
	log.Info().Msg("successful message")
	// Output: {"level":"info","message":"successful message"}
}

func ExampleSetLogOptions() {
	option := otelzerolog.WithAttributes(attribute.String("test", "value"), attribute.Bool("isValid", true))
	// use of SetLogOptions
	otelzerolog.SetLogOptions(option)

	// set up tracer
	tracer := otel.Tracer("otelzerolog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	// pass span to AddTracingContext
	log.Info().Func(otelzerolog.AddTracingContext(span)).Msg("in case of a success")

	// or in the case of an error
	err := errors.New("example error")
	log.Error().Func(otelzerolog.AddTracingContext(span, err)).Msg("in case of a failure")
	// Output: {"level":"info","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.test":"value","trace.attribute.isValid":true,"message":"in case of a success"}
	// {"level":"error","error":"example error","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.test":"value","trace.attribute.isValid":true,"message":"in case of a failure"}

}
