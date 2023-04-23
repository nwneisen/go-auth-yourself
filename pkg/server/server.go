package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/server/handlers"
)

// Server struct
type Server struct {
	config *config.Config
	logger *logger.Logger
	mux    *http.ServeMux
}

// Create a new server instance
func NewServer() *Server {
	// Read the configs
	config := config.NewConfig()
	config.LoadConfig("configs/dev.yaml")

	// Setup the logger
	logger := logger.NewLogger()
	logger.Info("logger created")

	// Setup the mux
	mux := http.NewServeMux()
	logger.Info("mux created")

	// Create the server
	return &Server{
		config: config,
		logger: logger,
		mux:    mux,
	}
}

// Start listening for requests
func (s *Server) Start() {
	s.logger.Info("listening for requests on port %s and %s", s.config.HttpPort, s.config.HttpsPort)
	go http.ListenAndServe(":"+s.config.HttpPort, http.HandlerFunc(s.RedirectToHTTPS()))
	http.ListenAndServeTLS(":"+s.config.HttpsPort, "server.cert", "server.key", s.mux)
}

// Add a handler to the server
func (s *Server) AddHandler(path string, newHandlerFunc func(config *config.Config, logger *logger.Logger) handlers.Handler) {
	handler := newHandlerFunc(s.config, s.logger)
	wrappedHandler := handlers.NewHandlerWrapper(s.config, s.logger, handler)
	s.mux.Handle(path, wrappedHandler)
}

// RedirectToHTTPS sends all HTTP requests to HTTPS
func (s *Server) RedirectToHTTPS() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		host, _, err := net.SplitHostPort(req.Host)
		if err != nil {
			log.Println(err)
			host = req.Host
		}
		target := fmt.Sprintf("https://%s:%s%s", host, s.config.HttpsPort, req.URL.Path)

		if len(req.URL.RawQuery) > 0 {
			target += "?" + req.URL.RawQuery
		}
		s.logger.Info("redirect to: %s", target)
		http.Redirect(w, req, target,
			// see comments below and consider the codes 308, 302, or 301
			http.StatusTemporaryRedirect)
	}
}

// Add middleware to the server
// Needs refactoring, untested
// func (s *Server) AddTracerMiddleware() http.Handler {
// 	s.mux = tracer.NewTracer(s.mux, s.logger)
// }
