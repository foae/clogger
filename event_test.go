package clogger

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventMethodsPkg(t *testing.T) {
	ctx, ev := NewEvent(context.Background(), "Something happened. Bam!")
	defer ev.End()

	require.NotNil(t, ev)
	require.NotNil(t, ev.fields)
	require.NotNil(t, ev.labels)
	require.NotNil(t, ev.errFields)
	assert.NotEmpty(t, ev.message)

	/*
		Test with package calls
	*/
	require.NotPanics(t, func() {
		Set(ctx, "set_key", "set_value")
		assert.Equal(t, ev.fields.retrieve("set_key").(string), "set_value")

		SetOnErr(ctx, "error_key", "error_value")
		assert.Equal(t, ev.errFields.retrieve("error_key").(string), "error_value")

		SetLabel(ctx, "label_key", "label_value")
		assert.Equal(t, ev.labels.retrieve("label_key").(string), "label_value")
	})
	ev.End()

	require.NotNil(t, ctx)
	assert.NoError(t, ctx.Err())
}

func TestEventMethods(t *testing.T) {
	ctx, ev := NewEvent(context.Background(), "Something happened. Bam!")
	defer ev.End()

	require.NotNil(t, ev)
	require.NotNil(t, ev.fields)
	require.NotNil(t, ev.labels)
	require.NotNil(t, ev.errFields)
	assert.NotEmpty(t, ev.message)

	/*
		Test with method calls
	*/
	require.NotPanics(t, func() {
		ev.Set("set_key", "set_value")
		assert.Equal(t, ev.fields.retrieve("set_key").(string), "set_value")

		ev.SetOnErr("error_key", "error_value")
		assert.Equal(t, ev.errFields.retrieve("error_key").(string), "error_value")

		ev.SetLabel("label_key", "label_value")
		assert.Equal(t, ev.labels.retrieve("label_key").(string), "label_value")
	})
	ev.End()

	require.NotNil(t, ctx)
	assert.NoError(t, ctx.Err())
}

func TestEventOptions(t *testing.T) {
	l := NewDefaultLogger()
	l.SetEventOptions(WithExampleEventOption(), WithOpenCensusTrace(), WithOpenTelemetryTrace())
	SetGlobal(l)

	assert.Equal(t, 3, len(l.EventOptions()))
	assert.Equal(t, 3, len(l.EventOptions()))

	for _, opt := range l.EventOptions() {
		typ := reflect.TypeOf(opt)
		assert.Equal(t, "func", typ.Kind().String(), "Expecting a func (callable)")
		assert.Equal(t, "EventOption", typ.Name())
	}

}

func TestNewEvent(t *testing.T) {
	ctx, ev := NewEvent(context.Background(), "Ooops!")
	require.NotNil(t, ctx)
	require.NotNil(t, ev)
	require.NotPanics(t, func() {
		assert.NotNil(t, ctx.Value(eventKey).(*Event))
		assert.Equal(t, ev, ctx.Value(eventKey).(*Event), "Ctx should contain the same event")
	})

	assert.Equal(t, SeverityInfo, ev.severity)
	assert.Equal(t, "Ooops!", ev.message)

	assert.NotNil(t, ev.fields)
	assert.NotNil(t, ev.errFields)
	assert.NotNil(t, ev.labels)

}
