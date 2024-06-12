package sentry

import (
	"context"
	"github.com/getsentry/sentry-go"
)

type Span = sentry.Span

func StartTransaction(ctx context.Context, operation string, options ...sentry.SpanOption) *sentry.Span {
	return sentry.StartTransaction(ctx, operation, options...)
}

func StartSpan(ctx context.Context, operation string, options ...sentry.SpanOption) *sentry.Span {
	return sentry.StartSpan(ctx, operation, options...)
}

func WithSpan0(ctx context.Context, operation string, f func(span *sentry.Span)) {
	span := sentry.StartSpan(ctx, operation)
	f(span)
	span.Finish()
}

func WithSpan1[T any](ctx context.Context, operation string, f func(span *sentry.Span) T) T {
	span := sentry.StartSpan(ctx, operation)
	r1 := f(span)
	span.Finish()
	return r1
}

func WithSpan2[T any, U any](ctx context.Context, operation string, f func(span *sentry.Span) (T, U)) (T, U) {
	span := sentry.StartSpan(ctx, operation)
	r1, r2 := f(span)
	span.Finish()
	return r1, r2
}
