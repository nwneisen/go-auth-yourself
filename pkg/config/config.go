package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config top level config layout
type Config struct {
	HttpPort  string `yaml:"httpPort"`
	HttpsPort string `yaml:"httpsPort"`

	Routes map[string]Route `yaml:"routes"`
}

// Possible redirect routes and their information
type Route struct {
	EgressHostname string `yaml:"egressHostname"`
	Port           string `yaml:"port"`

	IdpSsoUrl   string `yaml:"idpSsoUrl"`
	IdpIssuer   string `yaml:"idpIssuer"`
	IdpCertPath string `yaml:"idpCertPath"`

	GoogleClientId     string `yaml:"googleClientId"`
	GoogleClientSecret string `yaml:"googleClientSecret"`
}

// NewConfig creates a new config structure
func NewConfig() *Config {
	return &Config{}
}

// SaveConfig saves the configuration to a file
func (c *Config) SaveConfig(filePath string) {
	log.Printf("Saving the config to %s\n", filePath)
}

// LoadConfig loads the configuration from a file
func (c *Config) LoadConfig(filePath string) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("%v\n", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
