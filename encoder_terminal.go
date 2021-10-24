package clogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

var (
	timeFormatTerm = "2006/01/02 - 15:04:05"
)

type TerminalEncoder struct{}

func (t *TerminalEncoder) EncodeLogEntry(entry *LogEntry) ([]byte, error) {
	fields := make(map[string]interface{})
	for _, field := range entry.fields.fields() {
		fields[field.key] = field.value
	}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(fields)

	out := fmt.Sprintf("%s\t %s\t %s | %s",
		time.Now().Format(timeFormatTerm),
		entry.severity,
		entry.message,
		buf.Bytes(),
	)

	return []byte(out), err
}

func (t *TerminalEncoder) EncodeEvent(event *Event) ([]byte, error) {
	fields := make(map[string]interface{})

	// Fields
	for _, field := range event.fields.fields() {
		fields["fields"] = map[string]interface{}{
			field.key: field.value,
		}
	}

	// Err fields
	for _, field := range event.errFields.fields() {
		fields["errors"] = map[string]interface{}{
			field.key: field.value,
		}
	}

	// Labels
	for _, field := range event.labels.fields() {
		fields["labels"] = map[string]interface{}{
			field.key: field.value,
		}
	}

	fields["elapsed"] = time.Since(event.timestamp).String()

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(fields)

	out := fmt.Sprintf("%s\t %s\t %s | %s",
		time.Now().Format(timeFormatTerm),
		event.severity,
		event.message,
		buf.Bytes(),
	)

	return []byte(out), err
}
