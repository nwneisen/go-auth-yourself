package callbacks

import (
	"net/http"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
)

// Handlers is a generic handler for none specific routes
type Callbacks struct {
	config *config.Config
	logger *logger.Logger
}

// NewCallbacks creates a new callback handler
func NewCallbacks(config *config.Config, logger *logger.Logger) *Callbacks {
	return &Callbacks{config, logger}
}

// ServeHTTP handles the request by passing it to the real handler
func (c *Callbacks) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.logger.Info("Callback handler called")
}
