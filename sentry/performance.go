package sentry

import (
	"context"
	"github.com/getsentry/sentry-go"
)

func StartSpan(context context.Context, operationName string) *sentry.Span {
	return sentry.StartSpan(context, operationName)
}
