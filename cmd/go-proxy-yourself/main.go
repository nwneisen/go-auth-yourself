package main

import (
	"net/http"
	"nwneisen/go-proxy-yourself/internal/handlers"
	oauth "nwneisen/go-proxy-yourself/internal/handlers/oAuth"
	"nwneisen/go-proxy-yourself/internal/handlers/saml"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/tracer"
)

func main() {
	// Read the configs
	config := config.NewConfig()
	config.LoadConfig("configs/default.yaml")

	// Setup the logger
	logger := logger.NewLogger()
	logger.Info("Creating logger")

	// Add http redirect handler
	handlers := handlers.NewHandlers(config, logger)
	oAuth := oauth.NewOAuth(config, logger)
	saml := saml.NewSaml(config, logger)

	// Add https handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Index)
	mux.HandleFunc("/saml", saml.Index)
	mux.Handle("/oauth", oAuth)

	wrappedMux := addMiddleware(mux, logger)

	// Start listening for requests
	logger.Info("Listening for requests on port %s and %s", config.HttpPort, config.HttpsPort)
	go http.ListenAndServe(":"+config.HttpPort, http.HandlerFunc(handlers.RedirectToHTTPS()))
	http.ListenAndServeTLS(":"+config.HttpsPort, "server.cert", "server.key", wrappedMux)
}

// Add middlewares
func addMiddleware(h http.Handler, l *logger.Logger) http.Handler {
	tracerMux := tracer.NewTracer(h, l)

	return tracerMux
}
