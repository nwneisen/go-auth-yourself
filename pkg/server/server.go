package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"nwneisen/go-proxy-yourself/internal/handlers"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/metrics"
	serverHandlers "nwneisen/go-proxy-yourself/pkg/server/handlers"
	"nwneisen/go-proxy-yourself/pkg/tracing"
)

// Server struct
type Server struct {
	mux *http.ServeMux
}

// Create a new server instance
func NewServer() *Server {
	// Setup the logger
	logger.InitLogging()

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "go-auth-yourself"
	}

	jaegerURL := tracing.GetJaegerURL()
	if tracing.IsEnabled() {
		_, err := tracing.InitTracer(serviceName, jaegerURL)
		if err != nil {
			logger.Warn("Failed to initialize tracing: %v", err)
		}
		defer tracing.Shutdown()
	}

	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "9090"
	}
	go func() {
		if err := metrics.SetupMetricsServer(metricsPort); err != nil {
			logger.Warn("Failed to start metrics server: %v", err)
		}
	}()

	// Setup the mux
	mux := http.NewServeMux()
	logger.Info("mux created")

	// Create the server
	server := &Server{
		mux: mux,
	}

	// Add handlers
	server.AddHandler("/", handlers.NewIndexHandler)
	server.AddHandler("/config", handlers.NewConfigHandler)
	server.AddHandler("/oauth", handlers.NewOAuth)
	server.AddHandler("/saml", handlers.NewSamlHandler)
	server.AddHandler("/callback", handlers.NewCallbacksHandler)

	return server
}

// Start listening for requests
func (s *Server) Start() {

	logger.Info("listening for requests on port %s and %s", config.HttpPort(), config.HttpsPort())
	go http.ListenAndServe(":"+config.HttpPort(), http.HandlerFunc(s.RedirectToHTTPS()))
	http.ListenAndServeTLS(":"+config.HttpsPort(), "bin/server.cert", "bin/server.key", s.mux)
}

// Add a handler to the server
func (s *Server) AddHandler(path string, newHandlerFunc func() handlers.Handler) {
	handler := newHandlerFunc()
	wrappedHandler := serverHandlers.NewHandlerWrapper(handler)
	s.mux.Handle(path, wrappedHandler)
}

// RedirectToHTTPS sends all HTTP requests to HTTPS
func (s *Server) RedirectToHTTPS() func(w http.ResponseWriter, req *http.Request) {
	middleware := tracing.TracingMiddleware("go-auth-yourself")
	wrappedHandler := middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		host, _, err := net.SplitHostPort(req.Host)
		if err != nil {
			log.Println(err)
			host = req.Host
		}
		target := fmt.Sprintf("https://%s:%s%s", host, config.HttpsPort(), req.URL.Path)

		if len(req.URL.RawQuery) > 0 {
			target += "?" + req.URL.RawQuery
		}
		logger.Info("redirect to: %s", target)
		http.Redirect(w, req, target,
			http.StatusTemporaryRedirect)
	}))
	return func(w http.ResponseWriter, req *http.Request) {
		wrappedHandler.ServeHTTP(w, req)
	}
}

// Add middleware to the server
// Needs refactoring, untested
// func (s *Server) AddTracerMiddleware() http.Handler {
// 	s.mux = tracer.NewTracer(s.mux, logger)
// }
