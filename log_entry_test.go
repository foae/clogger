package clogger

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_exampleLog(t *testing.T) {
	l := NewDefaultLogger()
	SetGlobal(l)

	ctx := context.Background()
	assert.NotPanics(t, func() {
		Debug(ctx, "A debug statement happened.")
		Info(nil, "An informational message.")
	})

}

func Test_exampleLogWithFields(t *testing.T) {
	l := NewDefaultLogger()
	SetGlobal(l)

	ctx := context.Background()
	assert.NotPanics(t, func() {
		With("key1", "value1").
			With("key2", "value2").
			Debug(ctx, "A debug statement happened.")

		With("key3", "value3").
			Info(nil, "An informational message.")
	})
}

func TestLogEntryOptions(t *testing.T) {
	l := NewDefaultLogger()
	l.SetLogEntryOptions(WithExampleEntryOption(), WithOpenTelemetrySpan(), WithOpenTelemetrySpan())
	SetGlobal(l)

	assert.Equal(t, 3, len(l.LogEntryOptions()))
	assert.Equal(t, 3, len(l.LogEntryOptions()))

	for _, opt := range l.LogEntryOptions() {
		typ := reflect.TypeOf(opt)
		assert.Equal(t, "func", typ.Kind().String(), "Expecting a func (callable)")
		assert.Equal(t, "LogEntryOption", typ.Name())
	}

}
