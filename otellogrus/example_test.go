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

package otellogrus_test

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/vincentfree/opentelemetry/otellogrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func ExampleAddTracingContext() {
	tracer := otel.Tracer("otellogrus/example")
	_, span := tracer.Start(context.Background(), "example-span")

	logrus.WithFields(otellogrus.AddTracingContext(span)).Info("in case of a success")
	// or in the case of an error
	err := errors.New("example error")
	logrus.WithFields(otellogrus.AddTracingContext(span, err)).Info("in case of a failure")
}

func ExampleAddTracingContextWithAttributes() {
	tracer := otel.Tracer("otellogrus/example")
	_, span := tracer.Start(context.Background(), "example-span")

	attributes := []attribute.KeyValue{
		attribute.String("exampleKey", "exampleValue"),
		attribute.Bool("isValid", true),
	}

	logrus.WithFields(otellogrus.AddTracingContextWithAttributes(span, attributes)).Info("in case of a success")
	// or in the case of an error
	err := errors.New("example error")
	logrus.WithFields(otellogrus.AddTracingContextWithAttributes(span, attributes, err)).Info("in case of a failure")
}

func ExampleWithAttributes() {
	option := otellogrus.WithAttributes(attribute.String("test", "value"), attribute.Bool("isValid", true))
	otellogrus.SetLogOptions(option)

	tracer := otel.Tracer("otellogrus/example")
	_, span := tracer.Start(context.Background(), "example-span")

	logrus.WithFields(otellogrus.AddTracingContext(span)).Info("in case of a success")

	// or in the case of an error
	err := errors.New("example error")
	logrus.WithFields(otellogrus.AddTracingContext(span, err)).Info("in case of a failure")
}

func ExampleWithAttributePrefix() {
	otellogrus.SetLogOptions(otellogrus.WithAttributePrefix("prefix"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleWithServiceName() {
	otellogrus.SetLogOptions(otellogrus.WithServiceName("example-service"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleWithSpanID() {
	otellogrus.SetLogOptions(otellogrus.WithSpanID("span-id"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleWithTraceID() {
	otellogrus.SetLogOptions(otellogrus.WithTraceID("trace-id"))
	// use AddTracingContext or AddTracingContextWithAttributes
}

func ExampleSetLogOptions() {
	option := otellogrus.WithAttributes(attribute.String("test", "value"), attribute.Bool("isValid", true))
	// use of SetLogOptions
	otellogrus.SetLogOptions(option)

	// set up tracer
	tracer := otel.Tracer("otellogrus/example")
	_, span := tracer.Start(context.Background(), "example-span")

	logrus.WithFields(otellogrus.AddTracingContext(span)).Info("in case of a success")

	// or in the case of an error
	err := errors.New("example error")
	logrus.WithFields(otellogrus.AddTracingContext(span, err)).Info("in case of a failure")
}

func ExampleNew() {
	logger := otellogrus.New(otellogrus.WithLevel(logrus.ErrorLevel), otellogrus.WithFormatter(&logrus.JSONFormatter{}))
	logger.Info("message")
}

func ExampleLogger_WithTracingContext() {
	tracer := otel.Tracer("otellogrus/example")
	_, span := tracer.Start(context.Background(), "example-span")

	logger := otellogrus.New(otellogrus.WithLevel(logrus.ErrorLevel), otellogrus.WithFormatter(&logrus.JSONFormatter{}))
	// the logger returned by the otellogrus lib extends logrus but only does it if you use the logger.WithTracingContext(...) function as the first entry
	// the function returns a logrus.Entry which has not been extended
	logger.WithTracingContext(span).Info("example message")

	err := errors.New("example error")
	// with error is not used due to WithTracingContext doing the same thing internally
	logger.WithTracingContext(span, err).Error("example error message")
}

func ExampleLogger_WithTracingContextAndAttributes() {
	attributes := []attribute.KeyValue{
		attribute.String("exampleKey", "exampleValue"),
		attribute.Bool("isValid", true),
	}

	tracer := otel.Tracer("otellogrus/example")
	_, span := tracer.Start(context.Background(), "example-span")

	logger := otellogrus.New(otellogrus.WithLevel(logrus.ErrorLevel), otellogrus.WithFormatter(&logrus.JSONFormatter{}))
	// the logger returned by the otellogrus lib extends logrus but only does it if you use the logger.WithTracingContext(...) function as the first entry
	// the function returns a logrus.Entry which has not been extended
	logger.WithTracingContextAndAttributes(span, attributes).Info("example message")

	err := errors.New("example error")
	// with error is not used due to WithTracingContextAndAttributes doing the same thing internally
	logger.WithTracingContextAndAttributes(span, attributes, err).Error("example error message")

}
