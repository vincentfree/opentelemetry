package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/vincentfree/opentelemetry/otellogrus"
	"log/slog"

	"github.com/vincentfree/opentelemetry/otelslog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx := context.Background()
	_, span := otel.Tracer("test").Start(context.Background(), "mainService")
	// zerolog implementation
	logger := otelzerolog.New(otelzerolog.WithOtelBridgeDisabled(), otelzerolog.WithServiceName("mainService"))
	logger.Info().Func(logger.AddTracingContext(span)).Msg("logging with zerolog")

	// slog implementation
	sLogger := otelslog.New(otelslog.WithOtelBridgeDisabled(), otelslog.WithServiceName("mainService"))
	sLogger.WithTracingContext(ctx, slog.LevelInfo, "logging with slog", span, nil)

	// logrus implementation
	lrLogger := otellogrus.New(otellogrus.WithLevel(logrus.InfoLevel), otellogrus.WithFormatter(&logrus.JSONFormatter{}))
	lrLogger.WithTracingContext(span).Info("logging with logrus")

}
