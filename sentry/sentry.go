package sentry

import (
	"github.com/getsentry/raven-go"
	"github.com/go-errors/errors"
)

type Options struct {
	Dsn     string
	Project string
}

func Initialise(options Options) (err error) {
	if err = raven.SetDSN(options.Dsn); err != nil {
		return
	}

	project = options.Project

	return
}

// log raw error
func Error(e error) {
	wrapped := errors.New(e)
	raven.Capture(constructErrorPacket(wrapped), nil)
}

func LogWithContext(e error, ctx ErrorContext) {
	wrapped := errors.New(e)
	raven.Capture(constructPacket(wrapped, raven.INFO), ctx.ToMap())
}

func ErrorWithContext(e error, ctx ErrorContext) {
	wrapped := errors.New(e)
	raven.Capture(constructErrorPacket(wrapped), ctx.ToMap())
}
