package logger

import (
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
}

func New() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.Printf("[INFO] "+msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.Printf("[ERROR] "+msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Printf("[DEBUG] "+msg, args...)
}
