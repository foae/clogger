package clogger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/compute/metadata"
)

var (
	projectID = ""
	once      sync.Once
)

func init() {
	once.Do(func() {
		pID, err := metadata.ProjectID()
		if err != nil {
			Errorf(context.Background(), "Unable to extract the GCP project ID from the env – your logs might not be displayed correctly in Stackdriver")
		}

		projectID = pID
	})
}

/*
	https://cloud.google.com/logging/docs/structured-logging
	{
	  "severity":"ERROR",
	  "message":"There was an error in the application.",
	  "httpRequest":{
		"requestMethod":"GET"
	  },
	  "time":"2020-10-12T07:20:50.52Z",
	  "logging.googleapis.com/labels":{
		"user_label_1":"value_1",
		"user_label_2":"value_2"
	  },
	  "logging.googleapis.com/sourceLocation":{
		"file":"get_data.py",
		"line":"142",
		"function":"getData"
	  },
	  "logging.googleapis.com/spanId":"000000000000004a",
	  "logging.googleapis.com/trace":"projects/my-projectid/traces/06796866738c859f2f19b7cfb3214824",
	  "logging.googleapis.com/trace_sampled":false
	}

*/
type StackdriverEncoder struct{}

func (j *StackdriverEncoder) EncodeLogEntry(entry *LogEntry) ([]byte, error) {
	fields := make(map[string]interface{})
	for _, field := range entry.fields.fields() {
		switch field.key {
		case opencensusSpanID, otelSpanID:
			fields["logging.googleapis.com/spanId"] = field.value
		case opencensusTraceID, otelTraceID:
			// "projects/my-projectid/traces/06796866738c859f2f19b7cfb3214824"
			// "projects/_my-project-id_/traces/_trace-id_"
			fields["logging.googleapis.com/trace"] = fmt.Sprintf("projects/%s/traces/%s", projectID, field.value)
		case opencensusSampled, otelSampled:
			fields["logging.googleapis.com/trace_sampled"] = field.value
		default:
			fields[field.key] = field.value
		}
	}

	// TODO: sourceLocation

	fields["message"] = entry.message
	fields["timestamp"] = entry.timestamp.Format(time.RFC3339Nano)
	fields["severity"] = entry.severity.String()

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(fields)

	return buf.Bytes(), err
}

func (j *StackdriverEncoder) EncodeEvent(event *Event) ([]byte, error) {
	fields := make(map[string]interface{})

	// Fields
	for _, field := range event.fields.fields() {
		switch field.key {
		case opencensusSpanID, otelSpanID:
			fields["logging.googleapis.com/spanId"] = field.value
			v, _ := field.value.(string)
			fields["logging.googleapis.com/trace_sampled"] = v != ""
		case opencensusTraceID, otelTraceID:
			fields["logging.googleapis.com/trace"] = field.value
			v, _ := field.value.(string)
			fields["logging.googleapis.com/trace_sampled"] = v != ""
		default:
			fields["fields"] = map[string]interface{}{
				field.key: field.value,
			}
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
		fields["logging.googleapis.com/labels"] = map[string]interface{}{
			field.key: field.value,
		}
	}

	// TODO: sourceLocation

	fields["message"] = event.message
	fields["timestamp"] = event.timestamp.Format(time.RFC3339Nano)
	fields["latencySeconds"] = time.Since(event.timestamp).String()
	fields["severity"] = event.severity.String()

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(fields)

	return buf.Bytes(), err
}
