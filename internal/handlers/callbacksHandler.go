package handlers

import (
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/responses"
	"nwneisen/go-proxy-yourself/pkg/server/handlers"
)

// Handlers is a generic handler for none specific routes
type CallbacksHandler struct {
	*handlers.BaseHandler
}

// NewCallbacks creates a new callback handler
func NewCallbacksHandler() handlers.Handler {
	return CallbacksHandler{
		BaseHandler: handlers.NewBaseHandler(),
	}
}

// ServeHTTP handles the request by passing it to the real handler
func (c CallbacksHandler) Get() *responses.Response {
	logger.Info("Callback handler called")
	return responses.OK("Callback handler called")
}
