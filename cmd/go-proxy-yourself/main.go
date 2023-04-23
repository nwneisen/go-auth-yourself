package main

import (
	"nwneisen/go-proxy-yourself/internal/handlers/callbacks"
	"nwneisen/go-proxy-yourself/internal/handlers/index"
	"nwneisen/go-proxy-yourself/internal/handlers/oauth"
	"nwneisen/go-proxy-yourself/internal/handlers/saml"
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
	server.AddHandler("/", index.NewIndex)
	server.AddHandler("/oauth", oauth.NewOAuth)
	server.AddHandler("/saml", saml.NewSaml)
	server.AddHandler("/callback", callbacks.NewCallbacks)
	server.Start()
}
