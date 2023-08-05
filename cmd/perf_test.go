package main

import (
	"bufio"
	"bytes"
	"context"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"github.com/vincentfree/opentelemetry/otellogrus"
	"github.com/vincentfree/opentelemetry/otelslog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/exp/slog"
	"math"
	"testing"
)

func BenchmarkLogrus(b *testing.B) {
	w := testWriter()

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = w
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("bench")
	}

	_ = w.Flush()
}

func BenchmarkLogrusTrace(b *testing.B) {
	w := testWriter()
	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = w
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.WithFields(otellogrus.AddTracingContext(span)).Info("bench")
	}

	_ = w.Flush()
}

func BenchmarkLogrusTraceWithAttr(b *testing.B) {
	w := testWriter()
	attrs := attributes()
	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = w
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.WithFields(otellogrus.AddTracingContextWithAttributes(span, attrs)).Info("bench")
	}

	_ = w.Flush()
}

func BenchmarkSlog(b *testing.B) {
	w := testWriter()
	logger := slog.New(slog.NewJSONHandler(w, nil))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("bench")
	}

	_ = w.Flush()
}
func BenchmarkSlogTrace(b *testing.B) {
	w := testWriter()
	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")
	logger := slog.New(slog.NewJSONHandler(w, nil))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.LogAttrs(nil, slog.LevelInfo, "bench", otelslog.AddTracingContext(span)...)
	}

	_ = w.Flush()
}
func BenchmarkSlogTraceWithAttr(b *testing.B) {
	w := testWriter()
	attrs := attributes()
	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")
	logger := slog.New(slog.NewJSONHandler(w, nil))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.LogAttrs(nil, slog.LevelInfo, "bench", otelslog.AddTracingContextWithAttributes(span, attrs)...)
	}

	_ = w.Flush()
}

func BenchmarkZerolog(b *testing.B) {
	w := testWriter()
	logger := zerolog.New(w)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info().Msg("bench")
	}

	_ = w.Flush()
}
func BenchmarkZerologTrace(b *testing.B) {
	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")

	w := testWriter()
	logger := zerolog.New(w)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Func(otelzerolog.AddTracingContext(span)).Msg("bench")
	}

	_ = w.Flush()
}
func BenchmarkZerologTraceWithAttr(b *testing.B) {
	attrs := attributes()
	_, span := otel.Tracer("test").Start(context.Background(), "serviceName")

	w := testWriter()
	logger := zerolog.New(w)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info().Func(otelzerolog.AddTracingContextWithAttributes(span, attrs)).Msg("bench")
	}

	_ = w.Flush()
}

func attributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.Float64("localFloat64", 42.0),
		attribute.Int64("localInt64", 42),
		attribute.BoolSlice("localBoolSlice", []bool{true}),
		attribute.Int64Slice("localInt64Slice", []int64{42}),
		attribute.Float64Slice("localFloat64Slice", []float64{42.0}),
		attribute.StringSlice("localStringSlice", []string{"test"}),
	}
}

func testWriter() *bufio.Writer {
	by := &bytes.Buffer{}
	by.Grow(math.MaxInt32)
	return bufio.NewWriter(by)
}
