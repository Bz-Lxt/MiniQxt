package logger

import (
	"log/slog"
	"os"
	"strings"
)

var L *slog.Logger

func Init(env string) {
	level := slog.LevelInfo
	if strings.EqualFold(env, "development") || strings.EqualFold(env, "dev") {
		level = slog.LevelDebug
	}
	L = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(L)
}

func Debug(msg string, args ...any) {
	if L == nil {
		return
	}
	L.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	if L == nil {
		return
	}
	L.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	if L == nil {
		return
	}
	L.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	if L == nil {
		return
	}
	L.Error(msg, args...)
}
