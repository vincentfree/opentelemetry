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
	"github.com/rs/zerolog/log"
	"github.com/vincentfree/opentelemetry/otelzerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func ExampleAddTracingContext() {
	tracer := otel.Tracer("otelzerolog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	log.Info().Func(otelzerolog.AddTracingContext(span)).Msg("in case of a success")
	// or in the case of an error
	err := errors.New("example error")
	log.Error().Func(otelzerolog.AddTracingContext(span, err)).Msg("in case of a failure")
}

func ExampleAddTracingContextWithAttributes() {
	tracer := otel.Tracer("otelzerolog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	attributes := []attribute.KeyValue{
		attribute.String("exampleKey", "exampleValue"),
		attribute.Bool("isValid", true),
	}

	log.Info().Func(otelzerolog.AddTracingContextWithAttributes(span, attributes)).Msg("in case of a success")
	// or in the case of an error
	err := errors.New("example error")
	log.Error().Func(otelzerolog.AddTracingContextWithAttributes(span, attributes, err)).Msg("in case of a failure")
}

func ExampleWithAttributes() {
	option := otelzerolog.WithAttributes(attribute.String("test", "value"), attribute.Bool("isValid", true))
	otelzerolog.SetLogOptions(option)

	tracer := otel.Tracer("otelzerolog/example")
	_, span := tracer.Start(context.Background(), "example-span")

	log.Info().Func(otelzerolog.AddTracingContext(span)).Msg("in case of a success")

	// or in the case of an error
	err := errors.New("example error")
	log.Error().Func(otelzerolog.AddTracingContext(span, err)).Msg("in case of a failure")
}

func ExampleWithAttributePrefix() {
	otelzerolog.SetLogOptions(otelzerolog.WithAttributePrefix("prefix"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleWithServiceName() {
	otelzerolog.SetLogOptions(otelzerolog.WithServiceName("example-service"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleWithSpanID() {
	otelzerolog.SetLogOptions(otelzerolog.WithSpanID("span-id"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleWithTraceID() {
	otelzerolog.SetLogOptions(otelzerolog.WithTraceID("trace-id"))
	// use AddTracingContext or AddTracingContextWithAttributes
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
}
