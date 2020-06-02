package sentry

type ErrorContext interface {
	ToMap() map[string]string
}
