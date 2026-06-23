package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func newTestLogger(buf *bytes.Buffer) *Logger {
	return &Logger{
		logger: log.New(buf, "", log.LstdFlags),
	}
}

func TestLogger_Info(t *testing.T) {
	buf := &bytes.Buffer{}
	l := newTestLogger(buf)

	l.Info("hello %s", "world")

	out := buf.String()

	if out == "" {
		t.Fatalf("expected output")
	}

	if !strings.Contains(out, "[INFO]") {
		t.Fatalf("expected INFO level")
	}

	if !strings.Contains(out, "hello world") {
		t.Fatalf("expected formatted message")
	}
}

func TestLogger_Error(t *testing.T) {
	buf := &bytes.Buffer{}
	l := newTestLogger(buf)

	l.Error("something failed")

	out := buf.String()

	if !strings.Contains(out, "[ERROR]") {
		t.Fatalf("expected ERROR level")
	}
}

func TestLogger_Debug(t *testing.T) {
	buf := &bytes.Buffer{}
	l := newTestLogger(buf)

	l.Debug("debug message")

	out := buf.String()

	if !strings.Contains(out, "[DEBUG]") {
		t.Fatalf("expected DEBUG level")
	}
}
