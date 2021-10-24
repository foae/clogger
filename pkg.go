package clogger

import (
	"context"
)

// SetGlobal ...
func SetGlobal(l Logger) {
	log = l
}

func logger() Logger {
	return log
}

/*
	Setters (Eventful interface)
*/
// Set decorates the event's fields with the given key/value pair. These fields are outputted only in the
// event entry itself, and it is not passed down to its child log entries (use SetLabel for that).
// Set Should only be used if you have started an event in your call chain and passed ctx.
// If no event was found, it will create (ad-hoc) an "Unnamed event" with its severity raised to Warning.
// All subsequent calls to the Eventful interface will be attached to this event,
// whether it existed or was created ad-hoc.
func Set(ctx context.Context, key string, value interface{}) Eventful {
	return eventFromCtx(ctx).Set(key, value)
}

// SetOnErr decorates the event's error fields with the given key/value pair. Error fields are skipped
// (not outputted) if "nothing happens", meaning if no log entry raises the severity level over the
// configured threshold. It is generally useful to log additional fields but only when an error occurs.
// SetOnErr Should only be used if you have started an event in your call chain and passed ctx.
// If no event was found, it will create (ad-hoc) an "Unnamed event" with its severity raised to Warning.
// All subsequent calls to the Eventful interface will be attached to this event,
// whether it existed or was created ad-hoc.
func SetOnErr(ctx context.Context, key string, value interface{}) Eventful {
	return eventFromCtx(ctx).SetOnErr(key, value)
}

// SetLabel decorates the event's labels with the given key/value pair. Labels pass down their
// fields (key/value pairs) to all the event's child log entries.
// SetLabel Should only be used if you have started an event in your call chain and passed ctx.
// If no event was found, it will create (ad-hoc) an "Unnamed event" with its severity raised to Warning.
// All subsequent calls to the Eventful interface will be attached to this event,
// whether it existed or was created ad-hoc.
func SetLabel(ctx context.Context, key string, value interface{}) Eventful {
	return eventFromCtx(ctx).SetLabel(key, value)
}

// With registers a set of fields
func With(key string, value interface{}) Loggable {
	entry := newLogEntry()
	entry.fields.add(key, value)
	return entry
}

// Debug creates a new log entry with the given severity.
func Debug(ctx context.Context, msg string) {
	newLogEntry().Debug(ctx, msg)
}

// Debugf creates a new log entry with the given severity.
func Debugf(ctx context.Context, msg string, args ...interface{}) {
	newLogEntry().Debugf(ctx, msg, args...)
}

func Info(ctx context.Context, msg string) {
	newLogEntry().Info(ctx, msg)
}

// Infof creates a new log entry with the given severity.
func Infof(ctx context.Context, msg string, args ...interface{}) {
	newLogEntry().Infof(ctx, msg, args...)
}

// Warn creates a new log entry with the given severity.
func Warn(ctx context.Context, msg string) {
	newLogEntry().Warn(ctx, msg)
}

// Warnf creates a new log entry with the given severity.
func Warnf(ctx context.Context, msg string, args ...interface{}) {
	newLogEntry().Warnf(ctx, msg, args...)
}

func Error(ctx context.Context, msg string) {
	newLogEntry().Error(ctx, msg)
}

// Errorf creates a new log entry with the given severity.
func Errorf(ctx context.Context, msg string, args ...interface{}) {
	newLogEntry().Errorf(ctx, msg, args...)
}

// Fatal creates a new log entry with the given severity.
func Fatal(ctx context.Context, msg string) {
	newLogEntry().Fatal(ctx, msg)
}

// Fatalf creates a new log entry with the given severity.
func Fatalf(ctx context.Context, msg string, args ...interface{}) {
	newLogEntry().Fatalf(ctx, msg, args...)
}
