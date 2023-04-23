package callbacks

import (
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/server/handlers"
	"nwneisen/go-proxy-yourself/pkg/server/responses"
)

// Handlers is a generic handler for none specific routes
type Callbacks struct {
	*handlers.BaseHandler
}

// NewCallbacks creates a new callback handler
func NewCallbacks(config *config.Config, logger *logger.Logger) handlers.Handler {
	return Callbacks{
		BaseHandler: handlers.NewBaseHandler(config, logger),
	}
}

// ServeHTTP handles the request by passing it to the real handler
func (c Callbacks) Get() *responses.Response {
	c.Log().Info("Callback handler called")
	return responses.OK("Callback handler called")
}
