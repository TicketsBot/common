package sentry

import (
	"context"
	"github.com/getsentry/sentry-go"
)

func StartSpan(operationName string) *sentry.Span {
	return sentry.StartSpan(context.Background(), operationName)
}
