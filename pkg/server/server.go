package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"nwneisen/go-proxy-yourself/internal/handlers/index"
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
	serve := &Server{
		config: config,
		logger: logger,
		mux:    mux,
	}

	root := index.NewIndex(serve.config, serve.logger)
	serve.AddHandler("/", &root)

	return serve
}

// Start listening for requests
func (s *Server) Start() {
	s.logger.Info("listening for requests on port %s and %s", s.config.HttpPort, s.config.HttpsPort)
	go http.ListenAndServe(":"+s.config.HttpPort, http.HandlerFunc(s.RedirectToHTTPS()))
	http.ListenAndServeTLS(":"+s.config.HttpsPort, "server.cert", "server.key", s.mux)
}

// Add a handler to the server
func (s *Server) AddHandler(path string, handler *handlers.Handler) {
	newHandler := handlers.NewHandlerWrapper(s.config, s.logger, handler)
	s.mux.Handle(path, newHandler)
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
