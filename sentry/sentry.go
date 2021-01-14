package sentry

import (
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

type Options struct {
	Dsn     string
	Project string
	Debug   bool
}

var logger *logrus.Logger = nil

func Initialise(options Options) (err error) {
	if err = raven.SetDSN(options.Dsn); err != nil {
		return
	}

	project = options.Project

	if options.Debug {
		logger = logrus.New()
	}

	return
}

// log raw error
func Error(err error) {
	raven.Capture(constructErrorPacket(err), nil)

	if logger != nil {
		logger.Error(err.Error())
	}
}

func ErrorWithContext(err error, ctx ErrorContext) {
	raven.Capture(constructErrorPacket(err), ctx.ToMap())

	if logger != nil {
		fields := make(map[string]interface{})
		for k, v := range ctx.ToMap() {
			fields[k] = v
		}

		logger.WithFields(fields).Error(err.Error())
	}
}

func Log(msg string, extra map[string]interface{}) {
	raven.Capture(constructLogPacket(msg, extra), nil)

	if logger != nil {
		logger.WithFields(extra).Info(msg)
	}
}

func LogWithTags(msg string, extra map[string]interface{}, tags map[string]string) {
	raven.Capture(constructLogPacket(msg, extra), tags)

	if logger != nil {
		fields := make(map[string]interface{})
		for k, v := range extra {
			fields[k] = v
		}
		for k, v := range tags {
			fields[k] = v
		}

		logger.WithFields(fields).Info(msg)
	}
}

func LogWithContext(err error, ctx ErrorContext) {
	raven.Capture(constructPacket(err, raven.INFO), ctx.ToMap())

	if logger != nil {
		fields := make(map[string]interface{})
		for k, v := range ctx.ToMap() {
			fields[k] = v
		}

		logger.WithFields(fields).Info(err.Error())
	}
}
