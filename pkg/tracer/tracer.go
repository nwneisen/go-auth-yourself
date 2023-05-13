package tracer

import (
	"net/http"
	"nwneisen/go-proxy-yourself/pkg/logger"
)

// TODO implement tracing

// Tracer is a middleware for adding the tracer spans
type Tracer struct {
	handler http.Handler
}

// ServeHTTP handles the request by passing it to the real handler
func (t *Tracer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Info("Start tracing span")
	t.dump(r)
	t.handler.ServeHTTP(w, r)
	logger.Info("Finish tracing span")
}

// NewTracer constructs a new Tracer middleware handler
func NewTracer(handlerToWrap http.Handler) *Tracer {
	return &Tracer{handlerToWrap}
}

// dump request information for debugging
func (t *Tracer) dump(req *http.Request) {
	logger.Info("Header values:")
	for key, value := range req.Header {
		logger.Info("%q:%q", key, value[0])
	}

	logger.Info("Query values:")
	for key, value := range req.URL.Query() {
		logger.Info("%q:%q", key, value[0])
	}
}
