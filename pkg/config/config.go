package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"nwneisen/go-proxy-yourself/internal/fields"
	"nwneisen/go-proxy-yourself/pkg/logger"
)

var globalConfig *fields.Root

const (
	DEFAULT_DEV_LOG = "configs/dev.yaml"
)

func InitConfig(configLocation string) error {
	globalConfig = &fields.Root{}
	LoadConfig(configLocation)
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
	route, ok := globalConfig.Routes[hostname]
	if !ok {
		err := fmt.Errorf("route not found in config: %s", hostname)
		logger.Error(err.Error())
		return nil, err
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
func SaveConfig(filePath string) {
	logger.Info("saving the config to %s\n", filePath)
	logger.Error("not implemented")
}

// LoadConfig loads the configuration from a file
func LoadConfig(filePath string) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Info("%v\n", err)
	}

	err = yaml.Unmarshal(yamlFile, globalConfig)
	if err != nil {
		logger.Fatal("%v\n", err)
	}
}
