package clogger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newLoggerWithDecorators(t *testing.T) {
	l := NewDefaultLogger()
	l.SetEventOptions(WithOpenTelemetryTrace(), WithOpenCensusTrace(), WithExampleEventOption())
	l.SetLogEntryOptions(WithOpenTelemetrySpan(), WithOpenCensusSpan(), WithExampleEntryOption())
	SetGlobal(l)
	defer SetGlobal(NewDefaultLogger())

	assert.Equal(t, 3, len(l.EventOptions()))
	assert.Equal(t, 3, len(l.LogEntryOptions()))

	ctx := context.Background()
	Info(ctx, "An Info log entry. Part 1.")
	Info(ctx, "An Info log entry. Part 2.")

}

func Test_Logger(t *testing.T) {
	SetGlobal(noopLogger)
	defer SetGlobal(NewDefaultLogger())

	ctx, ev := NewEvent(context.Background(), "A new test logger instance")
	Debug(ctx, "Debug called")
	Debugf(ctx, "Debugf with (%s) called", "argument")
	With("logger_key", "logger_value").Info(ctx, "Info With")
	With("logger_key_debug", "logger_value_debug").Debug(ctx, "Debug With")
	ev.End()

	assert.Equal(t, ev, ctx.Value(eventKey).(*Event))
	assert.NotPanics(t, func() {
		for _, l := range ev.logs {
			if l.fields.retrieve("logger_key") == "logger_value" {
				return
			}
		}

		panic("could not find required key/value pair")
	})

}

func Test_LoggerTerminalEnc(t *testing.T) {
	l := NewDefaultLogger()
	l.SetEncoder(&TerminalEncoder{})
	SetGlobal(l)
	defer SetGlobal(NewDefaultLogger())

	ctx, ev := NewEvent(context.Background(), "A new test logger instance with terminal encoding")
	Debug(ctx, "Debug called")
	Debugf(ctx, "Debugf with (%s) called", "argument")
	With("logger_key", "logger_value").Info(ctx, "Info With")
	With("logger_key_debug", "logger_value_debug").Debug(ctx, "Debug With")
	ev.Set("event_key", "event_value")
	ev.End()

	assert.Equal(t, ev, ctx.Value(eventKey).(*Event))
	assert.NotPanics(t, func() {
		for _, l := range ev.logs {
			if l.fields.retrieve("logger_key_debug") == "logger_value_debug" {
				return
			}
		}

		panic("could not find required key/value pair")
	})

}
