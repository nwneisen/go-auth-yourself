package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"nwneisen/go-proxy-yourself/internal/fields"
	"nwneisen/go-proxy-yourself/pkg/logger"

	"gopkg.in/yaml.v2"
)

var globalConfig *fields.Root

const (
	DEFAULT_DEV_LOG = "configs/dev.yaml"
)

// InitConfig initializes the configuration
func InitConfig(configLocation string) error {
	globalConfig = fields.EmptyRoot()

	// Create an empty config if the file does not exist
	if _, err := os.Stat(configLocation); os.IsNotExist(err) {
		logger.Warn("config file does not exist, creating empty config\n")
		SaveConfig(configLocation)
		return nil
	}

	err := LoadConfig(configLocation)
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	return nil
}

// EmptyConfig returns a config with default values
func EmptyConfig() *fields.Root {
	return fields.EmptyRoot()
}

// Return a map of all routes
func Routes() (map[string]*fields.Route, error) {
	routes := globalConfig.Routes
	return routes, nil
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
	return globalConfig.HttpPort
}

func HttpsPort() string {
	return globalConfig.HttpsPort
}

// SaveConfig saves the configuration to a file
func SaveConfig(filePath string) error {
	logger.Debug("saving the config to %s\n", filePath)

	// Save the config file
	err := ioutil.WriteFile(filePath, []byte(globalConfig.YAML()), 0644)
	if err != nil {
		return fmt.Errorf("could not save the config to %s:%w", filePath, err)
	}

	return nil
}

// LoadConfig loads the configuration from a file
func LoadConfig(filePath string) error {
	logger.Debug("loading the config from %s\n", filePath)
	globalConfig := EmptyConfig()

	// Read the config file
	yamlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Unmarshal the config file
	// err = globalConfig.UnmarshalYAML(yamlBytes)
	// if err != nil {
	// 	return fmt.Errorf("%w", err)
	// }

	err = yaml.Unmarshal(yamlBytes, globalConfig)
	if err != nil {
		logger.Fatal("%w\n", err)
	}

	logger.Debug(globalConfig.String())

	return nil
}
