package clogger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	/*
		Set* without creating an event should not affect ctx,
		but should return a new event.
		Package calls
	*/
	t.Run("Set* has no effect on missing event from ctx", func(t *testing.T) {
		ctx := context.Background()
		Set(ctx, "set_key", "set_val").
			SetOnErr("set_err_key", "set_err_value").
			SetLabel("set_label_key", "set_label_value")

		assert.NotNil(t, ctx)

		ev := eventFromCtx(ctx)
		ev.End()
		assert.NotNil(t, ev)
		assert.NotEqual(t, "set_val", ev.fields.retrieve("set_key"))
		assert.NotEqual(t, "set_err_value", ev.fields.retrieve("set_err_key"))
		assert.NotEqual(t, "set_label_value", ev.fields.retrieve("set_label_key"))
	})

	/*
		Test with (a properly) created event
		Method calls
	*/
	t.Run("Method calls: event is persisted and populated", func(t *testing.T) {
		ctx := context.Background()
		cctx, ev := NewEvent(ctx, "A new event happened.")
		ev.Set("set_key", "set_value").
			SetOnErr("set_err_key", "set_err_value").
			SetLabel("set_label_key", "set_label_value")
		ev.End()

		assert.NotNil(t, cctx)
		assert.NotNil(t, ev)
		assert.Equal(t, "set_value", ev.fields.retrieve("set_key"))
		assert.Equal(t, "set_err_value", ev.errFields.retrieve("set_err_key"))
		assert.Equal(t, "set_label_value", ev.labels.retrieve("set_label_key"))
	})

	/*
		Test with (a properly) created event
		Package calls
	*/
	t.Run("Package calls: event is persisted and populated", func(t *testing.T) {
		ctx := context.Background()
		cctx, ev := NewEvent(ctx, "A new event happened.")
		Set(cctx, "set_key", "set_value").
			SetOnErr("set_err_key", "set_err_value").
			SetLabel("set_label_key", "set_label_value")
		ev.End()

		assert.NotNil(t, cctx)
		assert.NotNil(t, ev)
		assert.Equal(t, "set_value", ev.fields.retrieve("set_key"))
		assert.Equal(t, "set_err_value", ev.errFields.retrieve("set_err_key"))
		assert.Equal(t, "set_label_value", ev.labels.retrieve("set_label_key"))
	})
}
