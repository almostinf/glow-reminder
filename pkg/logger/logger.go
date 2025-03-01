package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -package mocks -destination mocks/logger_mocks.go github.com/almostinf/glow-reminder/pkg/logger Logger

// Logger is an interface describes the available logger methods.
type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	logrus.SetOutput(os.Stdout)
}

type logrusLogger struct {
	logger *logrus.Logger
}

// NewLogrusLogger creates a new instance logrus logger that implements Logger interface
// with given config level (panic, fatal, error, warn, warning, info, debug, trace).
func NewLogrusLogger(cfg Config) (Logger, error) {
	logger := logrus.New()

	logrusLevel, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse logrus level: %w", err)
	}

	logger.SetLevel(logrusLevel)

	return &logrusLogger{
		logger: logger,
	}, nil
}

// Debug logs a message at level Debug with the given fields.
func (l *logrusLogger) Debug(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Debug(msg)
}

// Info logs a message at level Info with the given fields.
func (l *logrusLogger) Info(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Info(msg)
}

// Warn logs a message at level Warn with the given fields.
func (l *logrusLogger) Warn(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Warn(msg)
}

// Error logs a message at level Error with the given fields.
func (l *logrusLogger) Error(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Error(msg)
}

// Fatal logs a message at level Fatal with the given fields.
func (l *logrusLogger) Fatal(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Fatal(msg)
}
