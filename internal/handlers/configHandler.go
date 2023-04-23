package handlers

import (
	"net/http"

	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/server/handlers"
)

// Handlers is a generic handler for none specific routes
type ConfigHandler struct {
	*handlers.BaseHandler
}

// NewConfig creates a new Config handler
func NewConfig() *ConfigHandler {
	return &ConfigHandler{}
}

// ServeHTTP handles the request by passing it to the real handler
func (c *ConfigHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger.Info("Callback handler called")
}
