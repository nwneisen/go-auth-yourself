package main

import (
	"nwneisen/go-proxy-yourself/internal/handlers"
	"nwneisen/go-proxy-yourself/pkg/server"
)

// main point of entry
func main() {
	// TODO handle args
	run()
}

// run starts the server
func run() {
	server := server.NewServer()

	server.AddHandler("/", handlers.NewIndexHandler)
	server.AddHandler("/config", handlers.NewConfigHandler)
	server.AddHandler("/oauth", handlers.NewOAuthHandler)
	server.AddHandler("/saml", handlers.NewSamlHandler)
	server.AddHandler("/callback", handlers.NewCallbacksHandler)

	server.Start()
}
