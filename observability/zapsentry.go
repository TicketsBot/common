package observability

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
	"os"
)

func ZapSentryAdapter(environment string) func(core zapcore.Core) zapcore.Core {
	return func(core zapcore.Core) zapcore.Core {
		return zapcore.RegisterHooks(core, func(entry zapcore.Entry) error {
			if entry.Level == zapcore.ErrorLevel {
				hostname, _ := os.Hostname()

				sentry.CaptureEvent(&sentry.Event{
					Environment: environment,
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
