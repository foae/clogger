package clogger

import (
	"os"
	"sync"
)

var (
	log Logger = NewDefaultLogger()
)

type Logger interface {
	StreamLogEntry(e *LogEntry)
	SetLogEntryOptions(opts ...LogEntryOption)
	LogEntryOptions() []LogEntryOption

	StreamEvent(ev *Event)
	SetEventOptions(opts ...EventOption)
	EventOptions() []EventOption

	SetEncoder(enc Encoder)
}

type DefaultLogger struct {
	eventOpts    []EventOption
	logEntryOpts []LogEntryOption
	mu           sync.Mutex
	enc          Encoder
}

func (l *DefaultLogger) SetEncoder(enc Encoder) {
	l.enc = enc
}

func NewDefaultLogger() Logger {
	return &DefaultLogger{
		eventOpts:    make([]EventOption, 0),
		logEntryOpts: make([]LogEntryOption, 0),
		mu:           sync.Mutex{},
		enc:          &JSONEncoder{},
	}
}

func (l *DefaultLogger) SetEventOptions(opts ...EventOption) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.eventOpts = opts
}

func (l *DefaultLogger) SetLogEntryOptions(opts ...LogEntryOption) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logEntryOpts = opts
}

func (l *DefaultLogger) EventOptions() []EventOption {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.eventOpts
}

func (l *DefaultLogger) LogEntryOptions() []LogEntryOption {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.logEntryOpts
}

func (l *DefaultLogger) StreamLogEntry(entry *LogEntry) {
	output, err := l.enc.EncodeLogEntry(entry)
	if err != nil {
		return
	}

	switch {
	case entry.severity > SeverityInfo:
		_, _ = os.Stderr.Write(output)
	default:
		_, _ = os.Stdout.Write(output)
	}
}

func (l *DefaultLogger) StreamEvent(event *Event) {
	output, err := l.enc.EncodeEvent(event)
	if err != nil {
		return
	}

	switch {
	case event.severity > SeverityInfo:
		_, _ = os.Stderr.Write(output)
	default:
		_, _ = os.Stdout.Write(output)
	}
}
