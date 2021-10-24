package clogger

import (
	"context"

	oc "go.opencensus.io/trace"
	"go.opentelemetry.io/otel"
)

const (
	opencensusSpanID  = "oc_span_id"
	opencensusTraceID = "oc_trace_id"
	opencensusSampled = "oc_sampled"

	otelSpanID  = "otel_span_id"
	otelTraceID = "otel_trace_id"
	otelSampled = "otel_sampled"
)

type LogEntryOption func(ctx context.Context, entry *LogEntry)

type EventOption func(ctx context.Context, event *Event)

func WithOpenTelemetryTrace() EventOption {
	return func(ctx context.Context, event *Event) {
		_, span := otel.Tracer("").Start(ctx, "clogger.Event.WithOpenTelemetryTrace")
		span.AddEvent(event.message)
		defer span.End()

		// TODO: offer the possibility to pass a Tracer/TraceProvider

		event.fields.add(otelSpanID, span.SpanContext().SpanID().String())   // TODO: the logger should be aware how to encode this key
		event.fields.add(otelTraceID, span.SpanContext().TraceID().String()) // TODO: the logger should be aware how to encode this key
		event.fields.add(otelSampled, span.SpanContext().IsSampled())        // TODO: the logger should be aware how to encode this key
	}
}

func WithOpenTelemetrySpan() LogEntryOption {
	return func(ctx context.Context, entry *LogEntry) {
		_, span := otel.Tracer("").Start(ctx, "clogger.LogEntry.WithOpenTelemetryTrace")
		defer span.End()

		// TODO: offer the possibility to pass a Tracer/TraceProvider

		entry.fields.add(otelSpanID, span.SpanContext().SpanID().String())   // TODO: the logger should be aware how to encode this key
		entry.fields.add(otelTraceID, span.SpanContext().TraceID().String()) // TODO: the logger should be aware how to encode this key
	}
}

func WithOpenCensusTrace() EventOption {
	return func(ctx context.Context, event *Event) {
		_, span := oc.StartSpan(ctx, "clogger.Event.WithOpenCensusTrace")
		defer span.End()

		event.fields.add(opencensusSpanID, span.SpanContext().SpanID.String())   // TODO: the logger should be aware how to encode this key
		event.labels.add(opencensusTraceID, span.SpanContext().TraceID.String()) // TODO: the logger should be aware how to encode this key
		event.labels.add(opencensusSampled, span.SpanContext().IsSampled())      // TODO: the logger should be aware how to encode this key
	}
}

func WithOpenCensusSpan() LogEntryOption {
	return func(ctx context.Context, entry *LogEntry) {
		_, span := oc.StartSpan(ctx, "clogger.LogEntry.WithOpenCensusSpan")
		defer span.End()

		entry.fields.add(opencensusSpanID, span.SpanContext().SpanID.String())   // TODO: the logger should be aware how to encode this key
		entry.fields.add(opencensusTraceID, span.SpanContext().TraceID.String()) // TODO: the logger should be aware how to encode this key
	}
}

func WithExampleEventOption() EventOption {
	return func(ctx context.Context, event *Event) {
		event.fields.add("example_event_id", "example_event_value")
	}
}

func WithExampleEntryOption() LogEntryOption {
	return func(ctx context.Context, entry *LogEntry) {
		entry.fields.add("example_entry_id", "example_entry_value")
	}
}
