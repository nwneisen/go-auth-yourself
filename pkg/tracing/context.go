package tracing

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/propagation"
)

type contextKey string

const traceIDKey contextKey = "trace_id"

func TraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}
	return ""
}

func AttachTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func ExtractTraceID(r *http.Request) string {
	traceID := r.Header.Get("X-Custom-Trace-Id")
	if traceID == "" {
		traceID = r.Header.Get("uber-trace-id")
	}
	return traceID
}

func InjectTraceID(w http.ResponseWriter, traceID string) {
	w.Header().Set("X-Custom-Trace-Id", traceID)
	w.Header().Set("uber-trace-id", traceID)
}

func GetPropagator() propagation.TextMapPropagator {
	return propagation.TraceContext{}
}
