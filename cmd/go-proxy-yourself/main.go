package main

import (
	"nwneisen/go-proxy-yourself/internal/handlers/index"
	"nwneisen/go-proxy-yourself/pkg/server"
)

func main() {
	// // Add http redirect handler
	// handlers := handlers.NewHandlers(config, logger)
	// oAuth := oauth.NewOAuth(config, logger)
	// saml := saml.NewSaml(config, logger)
	// callbacks := callbacks.NewCallbacks(config, logger)

	// // Add https handlers
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", handlers.Index)
	// mux.HandleFunc("/saml", saml.Index)
	// mux.Handle("/oauth", oAuth)
	// mux.Handle("/callback", callbacks)
	// // mux.Handle("/config", config)

	run()
}

func run() {
	server := server.NewServer()
	server.AddHandler("/", index.NewIndex)
	server.Start()
}
