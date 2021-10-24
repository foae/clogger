package clogger

import (
	"bytes"
	"encoding/json"
	"time"
)

type JSONEncoder struct{}

func (j *JSONEncoder) EncodeLogEntry(entry *LogEntry) ([]byte, error) {
	fields := make(map[string]interface{})
	for _, field := range entry.fields.fields() {
		fields[field.key] = field.value
	}

	fields["message"] = entry.message
	fields["timestamp"] = entry.timestamp.Format(time.RFC3339Nano)
	fields["severity"] = entry.severity.String()

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(fields)

	return buf.Bytes(), err
}

func (j *JSONEncoder) EncodeEvent(event *Event) ([]byte, error) {
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

	fields["message"] = event.message
	fields["timestamp"] = event.timestamp.Format(time.RFC3339Nano)
	fields["elapsed"] = time.Since(event.timestamp).String()
	fields["severity"] = event.severity.String()

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(fields)

	return buf.Bytes(), err
}
