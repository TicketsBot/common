package sentry

import (
	go_errors "errors"
	"github.com/getsentry/raven-go"
	"github.com/go-errors/errors"
	"github.com/rxdn/gdl/rest/request"
	"os"
	"time"
)

var project string

func constructErrorPacket(e *errors.Error) *raven.Packet {
	return constructPacket(e, raven.ERROR)
}

func constructPacket(e *errors.Error, level raven.Severity) *raven.Packet {
	hostname, err := os.Hostname(); if err != nil {
		hostname = "null"
		Error(err)
	}

	extra := map[string]interface{}{
		"stack": e.ErrorStack(),
	}

	var restError *request.RestError
	if go_errors.As(e, &restError) {
		extra["error_code"] = restError.ErrorCode
		extra["message"] = restError.Message
	}

	return &raven.Packet{
		Message: e.Error(),
		Extra: extra,
		Project: project,
		Timestamp: raven.Timestamp(time.Now()),
		Level: level,
		ServerName: hostname,
	}
}


