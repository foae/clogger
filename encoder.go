package clogger

type Encoder interface {
	EncodeLogEntry(entry *LogEntry) ([]byte, error)
	EncodeEvent(event *Event) ([]byte, error)
}
