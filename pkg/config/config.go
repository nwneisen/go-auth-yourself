package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"nwneisen/go-proxy-yourself/internal/fields"
	"nwneisen/go-proxy-yourself/pkg/logger"
)

var globalConfig *fields.Root

// InitConfig initializes the configuration using Viper
func InitConfig(configLocation string) error {
	// If configLocation is empty, use XDG default path
	if configLocation == "" {
		// Try to get the config file path using XDG standard
		configDir := os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("could not get user home directory: %w", err)
			}
			configDir = filepath.Join(homeDir, ".config")
		}

		configLocation = filepath.Join(configDir, "go-proxy-yourself", "config.yaml")
	}

	// Set up Viper configuration
	viper.SetConfigFile(configLocation)

	// Set default values using viper
	viper.SetDefault("http_port", "80")
	viper.SetDefault("https_port", "443")
	viper.SetDefault("server_cert", "bin/server.cert")
	viper.SetDefault("server_key", "bin/server.key")

	// Set default route structure
	viper.SetDefault("routes.simple_app.egress_hostname", "localhost")
	viper.SetDefault("routes.simple_app.port", "8081")

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create empty config
			logger.Warn("config file does not exist, creating empty config\n")
			// Initialize with defaults
			globalConfig = fields.EmptyRoot()
			SaveConfig(configLocation)
			return nil
		} else {
			return fmt.Errorf("could not read config: %w", err)
		}
	}

	// Load config into our fields.Root structure
	globalConfig = fields.EmptyRoot()

	// Set the values from Viper
	globalConfig.HttpPort = viper.GetString("http_port")
	globalConfig.HttpsPort = viper.GetString("https_port")
	globalConfig.ServerCert = viper.GetString("server_cert")
	globalConfig.ServerKey = viper.GetString("server_key")

	return nil
}

// Return a map of all routes
func Routes() (map[string]*fields.Route, error) {
	return globalConfig.Routes, nil
}

// Return an individual route by hostname
func Route(hostname string) (*fields.Route, error) {
	logger.Info("%v", globalConfig.Routes)

	route, ok := globalConfig.Routes[hostname]
	if !ok {
		return nil, fmt.Errorf("route not found in config: %s", hostname)
	}
	return route, nil
}

func HttpPort() string {
	return viper.GetString("http_port")
}

func HttpsPort() string {
	return viper.GetString("https_port")
}

// SaveConfig saves the configuration to a file using Viper
func SaveConfig(filePath string) error {
	logger.Debug("saving the config to %s\n", filePath)

	// Set values in Viper before saving
	viper.Set("http_port", globalConfig.HttpPort)
	viper.Set("https_port", globalConfig.HttpsPort)
	viper.Set("server_cert", globalConfig.ServerCert)
	viper.Set("server_key", globalConfig.ServerKey)

	// Save the config file
	err := viper.WriteConfigAs(filePath)
	if err != nil {
		return fmt.Errorf("could not save the config to %s:%w", filePath, err)
	}

	return nil
}

// GetConfig returns the global configuration
func GetConfig() *fields.Root {
	return globalConfig
}

// EmptyConfig returns an empty configuration for initialization
func EmptyConfig() *fields.Root {
	return fields.EmptyRoot()
}
