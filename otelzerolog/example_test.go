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

func ExampleLogger_AddTracingContext() {
	tracer := otel.Tracer("otelzerolog/ExampleAddTracingContext")
	_, span := tracer.Start(context.Background(), "example-span")
	defer span.End()
	logger := otelzerolog.New()

	logger.Info().Func(logger.AddTracingContext(span)).Msg("in case of a success")
	// or in the case of an error
	err := errors.New("example error")
	logger.Error().Func(logger.AddTracingContext(span, err)).Msg("in case of a failure")

	// Output: {"level":"info","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","message":"in case of a success"}
	// {"level":"error","error":"example error","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","message":"in case of a failure"}
}

func ExampleLogger_AddTracingContextWithAttributes() {
	tracer := otel.Tracer("otelzerolog/ExampleAddTracingContextWithAttributes")
	_, span := tracer.Start(context.Background(), "example-span")
	defer span.End()
	attributes := []attribute.KeyValue{
		attribute.String("exampleKey", "exampleValue"),
		attribute.Bool("isValid", true),
	}

	logger := otelzerolog.New()

	logger.Info().Func(logger.AddTracingContextWithAttributes(span, attributes)).Msg("in case of a success")
	// or in the case of an error
	err := errors.New("example error")
	logger.Error().Func(logger.AddTracingContextWithAttributes(span, attributes, err)).Msg("in case of a failure")

	// Output: {"level":"info","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.exampleKey":"exampleValue","trace.attribute.isValid":true,"message":"in case of a success"}
	// {"level":"error","error":"example error","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.exampleKey":"exampleValue","trace.attribute.isValid":true,"message":"in case of a failure"}
}

func ExampleWithAttributes() {
	option := otelzerolog.WithAttributes(attribute.String("test", "value"), attribute.Bool("isValid", true))
	logger := otelzerolog.New(option)

	tracer := otel.Tracer("otelzerolog/ExampleWithAttributes")
	_, span := tracer.Start(context.Background(), "example-span")
	defer span.End()
	logger.Info().Func(logger.AddTracingContext(span)).Msg("in case of a success")
	// or in the case of an error
	err := errors.New("example error")
	logger.Error().Func(logger.AddTracingContext(span, err)).Msg("in case of a failure")
	// Output: {"level":"info","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.test":"value","trace.attribute.isValid":true,"message":"in case of a success"}
	// {"level":"error","error":"example error","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.test":"value","trace.attribute.isValid":true,"message":"in case of a failure"}
}

func ExampleWithAttributePrefix() {
	tracer := otel.Tracer("otelzerolog/ExampleWithServiceName")
	_, span := tracer.Start(context.Background(), "example-span")
	defer span.End()
	logger := otelzerolog.New(otelzerolog.WithAttributePrefix("example-test"))
	// use AddTracingContext or AddTracingContextWithAttributes
	logger.Info().
		Func(logger.AddTracingContextWithAttributes(span, []attribute.KeyValue{
			attribute.String("example", "value"),
		})).
		Msg("success")
	// Output: {"level":"info","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","example-test.example":"value","message":"success"}
}

func ExampleWithServiceName() {
	logger := otelzerolog.New(otelzerolog.WithServiceName("example-service"))
	// use AddTracingContext or AddTracingContextWithAttributes

	tracer := otel.Tracer("otelzerolog/ExampleWithServiceName")
	_, span := tracer.Start(context.Background(), "ExampleWithServiceName")
	defer span.End()
	logger.Info().Func(logger.AddTracingContext(span)).Msg("successful message")
	// Output: {"level":"info","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","service.name":"example-service","message":"successful message"}
}

func ExampleWithSpanID() {
	logger := otelzerolog.New(otelzerolog.WithSpanID("span-id"))
	// use AddTracingContext or AddTracingContextWithAttributes

	tracer := otel.Tracer("otelzerolog/ExampleWithSpanID")
	_, span := tracer.Start(context.Background(), "example-span")
	defer span.End()

	logger.Info().Func(logger.AddTracingContext(span)).Msg("successful message")
	// Output: {"level":"info","traceID":"00000000000000000000000000000000","span-id":"0000000000000000","message":"successful message"}
}

func ExampleWithTraceID() {
	logger := otelzerolog.New(otelzerolog.WithTraceID("trace-id"))
	// use AddTracingContext or AddTracingContextWithAttributes

	tracer := otel.Tracer("otelzerolog/ExampleWithTraceID")
	_, span := tracer.Start(context.Background(), "example-span")
	defer span.End()

	logger.Info().Func(logger.AddTracingContext(span)).Msg("successful message")
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

func ExampleSetLogOptions() {
	otelzerolog.SetLogOptions(otelzerolog.WithServiceName("example-service"))
	log.Info().Msg("successful message")
	// Output: {"level":"info","message":"successful message"}
}

func ExampleSetGlobalLogger() {
	option := otelzerolog.WithAttributes(attribute.String("test", "value"), attribute.Bool("isValid", true))
	// use of SetLogOptions
	otelzerolog.SetGlobalLogger(option)

	// set up tracer
	tracer := otel.Tracer("otelzerolog/ExampleSetGlobalLogger")
	_, span := tracer.Start(context.Background(), "example-span")
	defer span.End()

	// pass span to AddTracingContext
	// expects that the otelzerolog.SetLogOptions or SetGlobalLogger has been used
	log.Info().Func(otelzerolog.AsOtelLogger(log.Logger).AddTracingContext(span)).Msg("in case of a success")

	// or in the case of an error
	err := errors.New("example error")
	log.Error().Func(otelzerolog.AsOtelLogger(log.Logger).AddTracingContext(span, err)).Msg("in case of a failure")
	// Output: {"level":"info","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.test":"value","trace.attribute.isValid":true,"message":"in case of a success"}
	// {"level":"error","error":"example error","traceID":"00000000000000000000000000000000","spanID":"0000000000000000","trace.attribute.test":"value","trace.attribute.isValid":true,"message":"in case of a failure"}
}
