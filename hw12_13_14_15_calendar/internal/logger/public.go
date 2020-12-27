package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

func New(logLevel string, output io.Writer, fileName string) (Logger, error) {
	log := logrus.New()

	result := logger{
		logger: log,
	}

	if logLevel != "" {
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			return result, fmt.Errorf("failed to parse log level: %w", err)
		}
		log.SetLevel(level)
	}

	if output != nil {
		log.SetOutput(output)
	} else if fileName != "" {
		fileName, err := filepath.Abs(fileName)
		if err != nil {
			return result, fmt.Errorf("failed to open log file: %w", err)
		}
		if err = os.MkdirAll(filepath.Dir(fileName), 0775); err != nil {
			return result, fmt.Errorf("failed to open log file: %w", err)
		}
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return result, fmt.Errorf("failed to open log file: %w", err)
		}
		log.SetOutput(file)
	}

	return result, nil
}
