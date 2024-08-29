package observability

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Configure(sentryDsn *string, json bool, logLevel zapcore.Level) (*zap.Logger, error) {
	if sentryDsn != nil {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn: *sentryDsn,
		}); err != nil {
			return nil, err
		}
	}

	if json {
		loggerConfig := zap.NewProductionConfig()
		loggerConfig.Level.SetLevel(logLevel)

		return loggerConfig.Build(
			zap.AddCaller(),
			zap.AddStacktrace(zap.ErrorLevel),
			zap.WrapCore(ZapSentryAdapter(EnvironmentProduction)),
		)
	} else {
		loggerConfig := zap.NewDevelopmentConfig()
		loggerConfig.Level.SetLevel(logLevel)
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		return loggerConfig.Build(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	}
}
