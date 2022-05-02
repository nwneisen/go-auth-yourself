package tracer

import (
	"log"
	"net/http"
	"nwneisen/go-proxy-yourself/pkg/logger"
)

// TODO implement tracing

// Tracer is a middleware for adding the tracer spans
type Tracer struct {
	handler http.Handler
	logger  *logger.Logger
}

// ServeHTTP handles the request by passing it to the real handler
func (t *Tracer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start tracing span")
	// t.dump(r)
	t.handler.ServeHTTP(w, r)
	log.Printf("Finish tracing span")
}

// NewTracer constructs a new Tracer middleware handler
func NewTracer(handlerToWrap http.Handler, logger *logger.Logger) *Tracer {
	return &Tracer{handlerToWrap, logger}
}

// dump request information for debugging
func (t *Tracer) dump(req *http.Request) {
	t.logger.Info("Header values:")
	for key, value := range req.Header {
		t.logger.Info("%q:%q", key, value[0])
	}

	t.logger.Info("Query values:")
	for key, value := range req.URL.Query() {
		t.logger.Info("%q:%q", key, value[0])
	}
}
