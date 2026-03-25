package services

import "log/slog"

type Logger interface {
	Info(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warning(msg string, args ...interface{})
	With(args ...any) *slog.Logger
}
