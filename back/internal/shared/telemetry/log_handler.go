package telemetry

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// LogHandler wraps a slog handler: a record logged within a traced request carries the
// trace and span identifiers, and is copied as an event onto the current span, where
// trace viewers show it. Error records also mark the span as failed.
type LogHandler struct {
	slog.Handler
}

// NewLogHandler returns a handler tying the records of next to the current span.
func NewLogHandler(next slog.Handler) *LogHandler {
	return &LogHandler{Handler: next}
}

func (h *LogHandler) Handle(ctx context.Context, record slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return h.Handler.Handle(ctx, record)
	}

	attributes := []attribute.KeyValue{attribute.String("level", record.Level.String())}
	record.Attrs(func(attr slog.Attr) bool {
		attributes = append(attributes, attribute.String(attr.Key, attr.Value.String()))
		return true
	})
	span.AddEvent(record.Message, trace.WithTimestamp(record.Time), trace.WithAttributes(attributes...))
	if record.Level >= slog.LevelError {
		span.SetStatus(codes.Error, record.Message)
	}

	record = record.Clone()
	record.AddAttrs(
		slog.String("trace_id", span.SpanContext().TraceID().String()),
		slog.String("span_id", span.SpanContext().SpanID().String()),
	)
	return h.Handler.Handle(ctx, record)
}

func (h *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewLogHandler(h.Handler.WithAttrs(attrs))
}

func (h *LogHandler) WithGroup(name string) slog.Handler {
	return NewLogHandler(h.Handler.WithGroup(name))
}
