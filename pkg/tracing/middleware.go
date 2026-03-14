package tracing

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Middleware func(http.Handler) http.Handler

func TracingMiddleware(serviceName string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			startTime := time.Now()

			ctx, span := otel.Tracer(serviceName).Start(ctx, r.URL.Path,
				trace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.url", r.URL.Path),
					attribute.String("http.user_agent", r.UserAgent()),
				),
			)
			defer span.End()

			propagator := otel.GetTextMapPropagator()
			propagator.Inject(ctx, propagation.HeaderCarrier(r.Header))

			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r.WithContext(ctx))

			span.SetAttributes(attribute.Int("http.status_code", wrapped.statusCode))
			span.SetAttributes(attribute.Float64("http.duration", time.Since(startTime).Seconds()))

			r.Header.Set("X-Custom-Trace-Id", span.SpanContext().TraceID().String())
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
