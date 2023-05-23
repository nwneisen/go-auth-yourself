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
func (h CallbacksHandler) Get() *responses.Response {
	logger.Info("Callback handler called")
	return responses.OK("Callback handler called")
}

// Post is the default POST method for all handlers
func (h CallbacksHandler) Post() *responses.Response {
	msg := "POST callback received"
	logger.Info(msg)
	return responses.BadRequest(msg)
}
