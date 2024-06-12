package sentry

import (
	"context"
	"github.com/getsentry/sentry-go"
)

func StartTransaction(ctx context.Context, operation string, options ...sentry.SpanOption) *sentry.Span {
	return sentry.StartTransaction(ctx, operation, options...)
}

func StartSpan(ctx context.Context, operation string, options ...sentry.SpanOption) *sentry.Span {
	return sentry.StartSpan(ctx, operation, options...)
}
