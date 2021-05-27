package sentry

import (
	"github.com/getsentry/raven-go"
	wrapper "github.com/go-errors/errors"
	"github.com/rxdn/gdl/rest/request"
	"os"
	"time"
)

var project string

func constructErrorPacket(e error) *raven.Packet {
	return constructPacket(e, raven.ERROR)
}

func constructPacket(e error, level raven.Severity) *raven.Packet {
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

	return &raven.Packet{
		Message:    e.Error(),
		Extra:      extra,
		Project:    project,
		Timestamp:  raven.Timestamp(time.Now()),
		Level:      level,
		ServerName: hostname,
	}
}

func constructLogPacket(msg string, extra map[string]interface{}) *raven.Packet {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "null"
	}

	return &raven.Packet{
		Message:    msg,
		Extra:      extra,
		Project:    project,
		Timestamp:  raven.Timestamp(time.Now()),
		Level:      raven.INFO,
		ServerName: hostname,
	}
}
