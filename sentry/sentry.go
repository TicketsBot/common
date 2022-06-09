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
func Error(err error) string {
	eventId, _ := raven.Capture(constructErrorPacket(err), nil)

	if logger != nil {
		logger.Error(err.Error())
	}

	return eventId
}

func ErrorWithContext(err error, ctx ErrorContext) string {
	eventId, _ := raven.Capture(constructErrorPacket(err), ctx.ToMap())

	if logger != nil {
		fields := make(map[string]interface{})
		for k, v := range ctx.ToMap() {
			fields[k] = v
		}

		logger.WithFields(fields).Error(err.Error())
	}

	return eventId
}

func Log(msg string, extra map[string]interface{}) string {
	eventId, _ := raven.Capture(constructLogPacket(msg, extra), nil)

	if logger != nil {
		logger.WithFields(extra).Info(msg)
	}

	return eventId
}

func LogWithTags(msg string, extra map[string]interface{}, tags map[string]string) string {
	eventId, _ := raven.Capture(constructLogPacket(msg, extra), tags)

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

	return eventId
}

func LogWithContext(err error, ctx ErrorContext) string {
	eventId, _ := raven.Capture(constructPacket(err, raven.INFO), ctx.ToMap())

	if logger != nil {
		fields := make(map[string]interface{})
		for k, v := range ctx.ToMap() {
			fields[k] = v
		}

		logger.WithFields(fields).Info(err.Error())
	}

	return eventId
}
