package handlers

import (
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/responses"
	"nwneisen/go-proxy-yourself/pkg/server/handlers"
)

// ConfigHandler for working with the config
type ConfigHandler struct {
	*handlers.BaseHandler
}

// NewConfigHandler creates a new Config handler
func NewConfigHandler() handlers.Handler {
	return ConfigHandler{
		BaseHandler: handlers.NewBaseHandler(),
	}
}

// Get returns the index.html page
func (h ConfigHandler) Get() *responses.Response {
	logger.Info("Index %s handler called", h.Request().Method)

	emptyConfig := config.EmptyConfig()
	return responses.JsonOK(emptyConfig.YAML())
}
