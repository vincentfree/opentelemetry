package main

import (
	"context"
	"github.com/rs/zerolog/log"

	// "github.com/sirupsen/logrus"
	// "github.com/vincentfree/opentelemetry/otellogrus"
	"github.com/vincentfree/opentelemetry/otelslog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
	"go.opentelemetry.io/otel"
	"log/slog"
)

func main() {
	ctx := context.Background()
	_, span := otel.Tracer("test").Start(context.Background(), "mainService")
	// zerolog implementation
	log.Info().Func(otelzerolog.AddTracingContext(span)).Msg("logging with zerolog")

	// slog implementation
	sLogger := otelslog.New()
	sLogger.WithTracingContext(ctx, slog.LevelInfo, "logging with slog", span, nil)

	// logrus implementation
	// lrLogger := otellogrus.New(otellogrus.WithLevel(logrus.InfoLevel), otellogrus.WithFormatter(&logrus.JSONFormatter{}))

}
