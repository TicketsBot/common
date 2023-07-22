package observability

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
	"os"
)

type Environment string

func (e Environment) String() string {
	return string(e)
}

const (
	EnvironmentProduction  Environment = "production"
	EnvironmentStaging     Environment = "staging"
	EnvironmentDevelopment Environment = "development"
)

func ZapSentryAdapter(environment Environment) func(core zapcore.Core) zapcore.Core {
	return func(core zapcore.Core) zapcore.Core {
		return zapcore.RegisterHooks(core, func(entry zapcore.Entry) error {
			if entry.Level == zapcore.ErrorLevel {
				hostname, _ := os.Hostname()

				sentry.CaptureEvent(&sentry.Event{
					Environment: environment.String(),
					Extra: map[string]any{
						"caller": entry.Caller.String(),
						"stack":  entry.Stack,
					},
					Level:      sentry.LevelError,
					Message:    entry.Message,
					ServerName: hostname,
					Timestamp:  entry.Time,
					Logger:     entry.LoggerName,
				})
			}

			return nil
		})
	}
}
