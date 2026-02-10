package server

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/server/handlers"
	"nwneisen/go-proxy-yourself/internal/handlers"
)

// Server struct
type Server struct {
	mux *http.ServeMux
}

// Create a new server instance
func NewServer() *Server {
	// Setup the logger
	logger.InitLogging()
	
	routes, err := config.Routes()
	if err != nil {
		logger.Warn("Failed to get routes: %v", err)
		routes = make(map[string]*fields.Route)
	}
	logger.Debug("%v", routes)

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
	server.AddHandler("/oauth", handlers.NewOAuthHandler)
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
	wrappedHandler := handlers.NewHandlerWrapper(handler)
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
		target := fmt.Sprintf("https://%s:%s%s", host, config.HttpsPort(), req.URL.Path)

		if len(req.URL.RawQuery) > 0 {
			target += "?" + req.URL.RawQuery
		}
		logger.Info("redirect to: %s", target)
		http.Redirect(w, req, target,
			// see comments below and consider the codes 308, 302, or 301
			http.StatusTemporaryRedirect)
	}
}

// Add middleware to the server
// Needs refactoring, untested
// func (s *Server) AddTracerMiddleware() http.Handler {
// 	s.mux = tracer.NewTracer(s.mux, logger)
// }
