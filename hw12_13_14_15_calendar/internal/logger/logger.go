package logger

import (
	"github.com/sirupsen/logrus"
)

type logger struct {
	logger *logrus.Logger
}

func (l logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}
