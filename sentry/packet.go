package sentry

import (
	"github.com/getsentry/raven-go"
	"github.com/go-errors/errors"
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

	return &raven.Packet{
		Message: e.Error(),
		Extra: extra,
		Project: project,
		Timestamp: raven.Timestamp(time.Now()),
		Level: level,
		ServerName: hostname,
	}
}


