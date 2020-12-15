package sentry

import (
	"github.com/getsentry/raven-go"
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
func Error(err error) {
	raven.Capture(constructErrorPacket(err), nil)
}

func ErrorWithContext(err error, ctx ErrorContext) {
	raven.Capture(constructErrorPacket(err), ctx.ToMap())
}

func Log(msg string, extra map[string]interface{}) {
	raven.Capture(constructLogPacket(msg, extra), nil)
}

func LogWithContext(err error, ctx ErrorContext) {
	raven.Capture(constructPacket(err, raven.INFO), ctx.ToMap())
}
