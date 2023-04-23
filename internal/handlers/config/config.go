package config

import (
	"net/http"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
)

// Handlers is a generic handler for none specific routes
type Config struct {
	config *config.Config
	logger *logger.Logger
}

// NewConfig creates a new Config handler
func NewConfig(config *config.Config, logger *logger.Logger) *Config {
	return &Config{config, logger}
}

// ServeHTTP handles the request by passing it to the real handler
func (c *Config) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.logger.Info("Callback handler called")
}
