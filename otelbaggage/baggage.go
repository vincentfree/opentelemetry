package otelbaggage

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var (
	members    []string
	propagator propagation.TextMapPropagator
)

type Option func(*config)

type config struct {
	members     []string
	propagators []propagation.TextMapPropagator
}

func DefaultOptions(options ...Option) {
	cfg := &config{}

	for _, o := range options {
		o(cfg)
	}

	if cfg.members != nil {
		members = cfg.members
	}

	if cfg.propagators != nil {
		// use composite text map propagator to handle one or multiple propagators as one
		p := propagation.NewCompositeTextMapPropagator(cfg.propagators...)
		propagator = p
	} else {
		propagator = propagation.Baggage{}
	}
	p := otel.GetTextMapPropagator()
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(p, propagator))
}

// WithMembers adds member keys to a list, this list is used to extract values from the baggage.
// When an application receives baggage values,
// then all baggage values with keys that correspond with this list can be injected into an active trace.Span.
func WithMembers(memberKeys ...string) Option {
	return func(c *config) {
		c.members = memberKeys
	}
}

func WithPropagators(propagators ...propagation.TextMapPropagator) Option {
	return func(c *config) {
		c.propagators = propagators
	}
}

// InjectIntoSpan extracts baggage.Baggage metadata from the context and carrier,
// formats it as key-value pair attributes, and infuses these attributes into the trace span.
// Returns the updated span with the newly incorporated attributes based on the members set in this package.
// when
func InjectIntoSpan(ctx context.Context, span trace.Span, carrier propagation.TextMapCarrier) trace.Span {
	b := extractBaggage(ctx, carrier)
	ms := extractMembers(b)

	attrs := make([]attribute.KeyValue, len(ms))

	for _, m := range ms {
		attrs = append(attrs, toKeyValue(m))
	}
	span.SetAttributes(attrs...)

	return span
}

func InjectIntoBaggage(ctx context.Context, members []baggage.Member) (baggage.Baggage, error) {
	var errs error
	b := baggage.FromContext(ctx)
	for _, m := range members {
		nb, err := b.SetMember(m)
		if err != nil {
			_ = errors.Join(errs, err)
		}
		b = nb
	}

	return b, errs
}

// todo make function to add the key values for baggage initially

// extractMembers extracts and returns a slice of baggage.Member objects from the input
// baggage.Baggage object which contain non-empty values and match the members.
func extractMembers(b baggage.Baggage) []baggage.Member {
	ms := make([]baggage.Member, len(members))
	for _, m := range members {
		member := b.Member(m)
		if member.String() == "" {
			continue
		}
		ms = append(ms, member)
	}
	return ms
}

func extractBaggage(ctx context.Context, carrier propagation.TextMapCarrier) baggage.Baggage {
	nCtx := propagator.Extract(ctx, carrier)
	return baggage.FromContext(nCtx)
}

func toKeyValue(member baggage.Member) attribute.KeyValue {
	return attribute.String(member.Key(), member.Value())
}
