package logger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
	"github.com/wigfri/mustore/domain/services"
)

const (
	envLocal       = "local"
	envDevelopment = "dev"
	envProduction  = "prod"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(timeStr, level, msg, color.WhiteString(string(b)))

	return nil
}

func NewPrettyHandler(
	out io.Writer,
	opts PrettyHandlerOptions,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}

type Logger struct {
	log *slog.Logger
}

func Init(env string) services.Logger {
	var logger Logger

	switch env {
	case envProduction:
		logger.log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envDevelopment:
		opts := PrettyHandlerOptions{SlogOpts: slog.HandlerOptions{Level: slog.LevelDebug}}
		logger.log = slog.New(
			NewPrettyHandler(os.Stdout, opts))
	case envLocal:
		opts := PrettyHandlerOptions{SlogOpts: slog.HandlerOptions{Level: slog.LevelDebug}}
		logger.log = slog.New(
			NewPrettyHandler(os.Stdout, opts))
	}

	return &logger
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log.Info(msg, args...)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log.Debug(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log.Error(msg, args...)
}

func (l *Logger) Warning(msg string, args ...interface{}) {
	l.log.Warn(msg, args...)
}
func (l *Logger) With(args ...any) *slog.Logger {
	return l.log.With(args)
}
