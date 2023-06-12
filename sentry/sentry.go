package sentry

import (
	"github.com/TicketsBot/common/utils"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger = nil

type Options sentry.ClientOptions

func Initialise(options Options) (err error) {
	if err := sentry.Init(sentry.ClientOptions(options)); err != nil {
		return err
	}

	if options.Debug {
		logger = logrus.New()
	}

	return
}

// log raw error
func Error(err error) string {
	eventId := sentry.CaptureEvent(constructErrorPacket(err, nil))

	if logger != nil {
		logger.Error(err.Error())
	}

	return string(utils.ValueOrZero(eventId))
}

func ErrorWithContext(err error, ctx ErrorContext) string {
	eventId := sentry.CaptureEvent(constructErrorPacket(err, ctx.ToMap()))

	if logger != nil {
		fields := make(map[string]interface{})
		for k, v := range ctx.ToMap() {
			fields[k] = v
		}

		logger.WithFields(fields).Error(err.Error())
	}

	return string(utils.ValueOrZero(eventId))
}

func Log(msg string, extra map[string]interface{}) string {
	eventId := sentry.CaptureEvent(constructLogPacket(msg, extra, nil))

	if logger != nil {
		logger.WithFields(extra).Info(msg)
	}

	return string(utils.ValueOrZero(eventId))
}

func LogWithTags(msg string, extra map[string]interface{}, tags map[string]string) string {
	eventId := sentry.CaptureEvent(constructLogPacket(msg, extra, tags))

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

	return string(utils.ValueOrZero(eventId))
}

func LogWithContext(err error, ctx ErrorContext) string {
	eventId := sentry.CaptureEvent(constructPacket(err, sentry.LevelInfo, ctx.ToMap()))

	if logger != nil {
		fields := make(map[string]interface{})
		for k, v := range ctx.ToMap() {
			fields[k] = v
		}

		logger.WithFields(fields).Info(err.Error())
	}

	return string(utils.ValueOrZero(eventId))
}
