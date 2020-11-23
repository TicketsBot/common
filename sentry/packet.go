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
		extra["error_code"] = restError.ErrorCode
		extra["message"] = restError.Message
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
