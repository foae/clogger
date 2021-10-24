package clogger

import (
	"context"
	"runtime"
	"sync"
	"time"
)

type Eventful interface {
	Set(key string, value interface{}) Eventful
	SetOnErr(key string, value interface{}) Eventful
	SetLabel(key string, value interface{}) Eventful
}

// Event describes a single action that happens at a given time.
// Typically, it has a short life span where it gathers field, labels and log entries.
// It must be accompanied by its End method to mark the once of its lifecycle, usually from a defer function call.
type Event struct {
	message   string
	timestamp time.Time

	// Severity takes a default Info level and gets
	// (1) raised if any child log is >
	// (2) skipped if the configured log level is <
	severity Severity

	// An event might gather log entries throughout its lifecycle
	logs []*LogEntry

	// An event might carry information on its own that we don't necessarily
	// want it to be passed down to its child logs (as with Labels).
	fields *fieldCollection

	// Holds fields that will be outputted only in case an error is encountered.
	// It will be skipped from output otherwise.
	errFields *fieldCollection

	// An event might be marked with important data by using Labels.
	// A label gets written/passed down to each and every log entry that it holds.
	labels *fieldCollection

	once sync.Once
}

func NewEvent(ctx context.Context, message string) (context.Context, *Event) {
	ev := &Event{
		message:   message,
		timestamp: time.Now(),
		severity:  SeverityInfo,
		logs:      make([]*LogEntry, 0),
		labels:    newFieldCollection(),
		fields:    newFieldCollection(),
		errFields: newFieldCollection(),
		once:      sync.Once{},
	}

	// Apply all decorators (modifiers) registered for this event
	for _, opt := range logger().EventOptions() {
		opt(ctx, ev)
	}

	return context.WithValue(ctx, eventKey, ev), ev
}

// Set registers a key/value pair for the current event.
// These pairs are not being passed down to child log entries.
func (ev *Event) Set(key string, value interface{}) Eventful {
	ev.fields.add(key, value)
	return ev
}

// SetOnErr registers a key/value pair for the current event that will be logged (outputted)
// only in case the severity is raised over the Info threshold e.g. by registering an error log entry.
// fields Set with this method get logged only if an error occurs during the event's lifecycle.
// Otherwise, these are discarded.
func (ev *Event) SetOnErr(key string, value interface{}) Eventful {
	ev.errFields.add(key, value)
	return ev
}

// SetLabel registers a key/value pair as a label. This pair will then be passed down
// and stamped/written onto all child log entries that this event might collect
// throughout its lifecycle.
func (ev *Event) SetLabel(key string, value interface{}) Eventful {
	ev.labels.add(key, value)
	return ev
}

// End signals the once of the lifecycle to the event.
// It will apply all gathered Labels onto all child log entries,
// and it will finally output using the configured logger instance.
// Safe to be called multiple times.
func (ev *Event) End() {
	ev.once.Do(func() {
		// Process all child log entries
		for _, entry := range ev.logs {
			if ev.labels.len() > 0 {
				// Transfer the event's Labels, if any are defined, onto each child log entry.
				// It will overwrite existing key/value pairs from the log entry's fields.
				entry.fields.merge(ev.labels)
			}

			// Raise severity level, if any of the child logs contains
			// a severity level greater than the event's.
			if entry.severity > ev.severity {
				ev.severity = entry.severity
			}

			// Output all child logs
			logger().StreamLogEntry(entry)
		}

		// The logger instance takes care of encoding
		// and streaming the contents of an event.
		logger().StreamEvent(ev)

	})
}

func eventFromCtx(ctx context.Context) *Event {
	event, ok := ctx.Value(eventKey).(*Event)
	if ok {
		return event
	}

	// Passed ctx does not contain an event.
	// Create one and warn.
	cctx, ev := NewEvent(ctx, "Unnamed event")
	ctx = cctx
	ev.severity = SeverityWarn

	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	With("caller_name", frame.Function).
		Errorf(context.Background(), "Called (%s) but there was no event found in context", frame.Function)

	return ev
}
