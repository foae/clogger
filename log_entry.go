package clogger

import (
	"context"
	"fmt"
	"os"
	"time"
)

type Loggable interface {
	With(key string, value interface{}) Loggable

	Debug(ctx context.Context, message string)
	Debugf(ctx context.Context, message string, args ...interface{})

	Info(ctx context.Context, message string)
	Infof(ctx context.Context, message string, args ...interface{})

	Warn(ctx context.Context, message string)
	Warnf(ctx context.Context, message string, args ...interface{})

	Error(ctx context.Context, message string)
	Errorf(ctx context.Context, message string, args ...interface{})

	Fatal(ctx context.Context, message string)
	Fatalf(ctx context.Context, message string, args ...interface{})
}

// LogEntry defines the structure of a single, unitary log entry.
type LogEntry struct {
	message   string
	timestamp time.Time

	// Set by the caller indirectly through
	// one of the methods Debug/Info/Warn etc.
	severity Severity

	// Any information logged by the caller
	// will be stored in a field collection
	fields *fieldCollection

	// Keep track if this log entry
	// is part of an Event's lifecycle
	eventful bool
}

func newLogEntry() *LogEntry {
	return &LogEntry{
		message:   "",
		timestamp: time.Now(),
		severity:  SeverityDebug,
		fields:    newFieldCollection(),
		eventful:  false,
	}
}

func (e *LogEntry) log(ctx context.Context, sev Severity, msg string) *LogEntry {
	e.message = msg
	e.severity = sev

	// Guard against nil contexts
	if ctx == nil {
		ctx = context.TODO()
	}

	// Apply any custom decorators onto this LogEntry
	e.apply(ctx)

	// Check whether this is part of a greater event
	if event, ok := ctx.Value(eventKey).(*Event); ok {
		// Mark this log entry as such
		e.eventful = true

		// Add this log entry as a child log,
		// attached to the found event.
		event.logs = append(event.logs, e)

		// If the log's severity is greater than the
		if sev > event.severity {
			event.severity = sev
		}
	}

	return e
}

// dispatch will output the LogEntry or return early
// if it's part of a bigger event.
func (e *LogEntry) dispatch() {
	// If this is part of a bigger event, return early.
	if e.eventful {
		return
	}

	// Otherwise, it's an isolated LogEntry;
	// output directly
	logger().StreamLogEntry(e)
}

func (e *LogEntry) apply(ctx context.Context) {
	for _, opt := range logger().LogEntryOptions() {
		opt(ctx, e)
	}
}

// With ...
func (e *LogEntry) With(key string, value interface{}) Loggable {
	if e == nil {
		e = newLogEntry()
	}

	e.fields.add(key, value)
	return e
}

// Debug ...
func (e *LogEntry) Debug(ctx context.Context, message string) {
	e.log(ctx, SeverityDebug, message).dispatch()
}

// Debugf ...
func (e *LogEntry) Debugf(ctx context.Context, message string, args ...interface{}) {
	e.log(ctx, SeverityDebug, fmt.Sprintf(message, args...)).dispatch()
}

// Info ...
func (e *LogEntry) Info(ctx context.Context, message string) {
	e.log(ctx, SeverityInfo, message).dispatch()
}

// Infof ...
func (e *LogEntry) Infof(ctx context.Context, message string, args ...interface{}) {
	e.log(ctx, SeverityInfo, fmt.Sprintf(message, args...)).dispatch()
}

// Warn ...
func (e *LogEntry) Warn(ctx context.Context, message string) {
	e.log(ctx, SeverityWarn, message).dispatch()
}

// Warnf ...
func (e *LogEntry) Warnf(ctx context.Context, message string, args ...interface{}) {
	e.log(ctx, SeverityWarn, fmt.Sprintf(message, args...)).dispatch()
}

// Error ...
func (e *LogEntry) Error(ctx context.Context, message string) {
	e.log(ctx, SeverityError, message).dispatch()
}

// Errorf ...
func (e *LogEntry) Errorf(ctx context.Context, message string, args ...interface{}) {
	e.log(ctx, SeverityError, fmt.Sprintf(message, args...)).dispatch()
}

// Fatal ...
func (e *LogEntry) Fatal(ctx context.Context, message string) {
	e.log(ctx, SeverityCritical, message).dispatch()
	os.Exit(1)
}

// Fatalf ...
func (e *LogEntry) Fatalf(ctx context.Context, message string, args ...interface{}) {
	e.log(ctx, SeverityCritical, fmt.Sprintf(message, args...)).dispatch()
	os.Exit(1)
}
