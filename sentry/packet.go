package sentry

import (
	"github.com/getsentry/sentry-go"
	wrapper "github.com/go-errors/errors"
	"github.com/rxdn/gdl/rest/request"
	"os"
	"time"
)

func constructErrorPacket(e error, tags map[string]string) *sentry.Event {
	return constructPacket(e, sentry.LevelError, tags)
}

func constructPacket(e error, level sentry.Level, tags map[string]string) *sentry.Event {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "null"
	}

	extra := map[string]interface{}{
		"stack": wrapper.New(e).ErrorStack(),
	}

	if restError, ok := e.(request.RestError); ok {
		extra["status_code"] = restError.StatusCode
		extra["message"] = restError.Error()
		extra["url"] = restError.Url
		extra["raw"] = string(restError.Raw)
	}

	return &sentry.Event{
		Message:    e.Error(),
		Extra:      extra,
		Timestamp:  time.Now(),
		Level:      level,
		ServerName: hostname,
		Tags:       tags,
	}
}

func constructLogPacket(msg string, extra map[string]interface{}, tags map[string]string) *sentry.Event {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "null"
	}

	return &sentry.Event{
		Message:    msg,
		Extra:      extra,
		Timestamp:  time.Now(),
		Level:      sentry.LevelInfo,
		ServerName: hostname,
		Tags:       tags,
	}
}
