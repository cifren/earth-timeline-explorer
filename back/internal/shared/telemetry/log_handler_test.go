package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestLogHandler(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	tracer := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder)).Tracer("test")
	var output bytes.Buffer
	logger := slog.New(NewLogHandler(slog.NewJSONHandler(&output, nil)))

	ctx, span := tracer.Start(context.Background(), "request")
	logger.ErrorContext(ctx, "cache read failed", "key", "countries:1950s")
	span.End()

	var line map[string]any
	if err := json.Unmarshal(output.Bytes(), &line); err != nil {
		t.Fatal(err)
	}
	if line["trace_id"] != span.SpanContext().TraceID().String() || line["span_id"] != span.SpanContext().SpanID().String() {
		t.Fatalf("log line without the span identifiers: %v", line)
	}
	ended := recorder.Ended()[0]
	events := ended.Events()
	if len(events) != 1 || events[0].Name != "cache read failed" {
		t.Fatalf("span events = %v", events)
	}
	if ended.Status().Code != codes.Error {
		t.Fatalf("span status = %v, want error", ended.Status())
	}
}

func TestLogHandlerOutsideATrace(t *testing.T) {
	var output bytes.Buffer
	logger := slog.New(NewLogHandler(slog.NewJSONHandler(&output, nil)))

	logger.InfoContext(context.Background(), "server listening")

	if bytes.Contains(output.Bytes(), []byte("trace_id")) {
		t.Fatalf("untraced log carries a trace identifier: %s", output.Bytes())
	}
}
