package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// var globalInstance Config

// func InitGlobalInstance(configType string, configLocation string) error {
// 	var err error
// 	globalInstance, err = NewConfig(name, cfg)
// 	return err
// }

// Config top level config layout
type Config struct {
	HttpPort  string `yaml:"httpPort" json:"httpPort" default:"80"`
	HttpsPort string `yaml:"httpsPort" json:"httpsPort" default:"443"`

	Routes map[string]*Route `yaml:"routes" json:"routes"`
}

// Possible redirect routes and their information
type Route struct {
	EgressHostname string `yaml:"egressHostname" json:"egressHostname"`
	Port           string `yaml:"port" json:"port"`

	SAML  map[string]*SAMLProvider  `yaml:"saml,omitempty" json:"saml,omitempty"`
	OAuth map[string]*OAuthProvider `yaml:"oAuth,omitempty" json:"oAuth,omitempty"`
}

// String returns a string representation of the Route
func (r *Route) String() string {
	return fmt.Sprintf("Route{egressHostname:%s, Port:%s, SAML:%v, OAuth:%v}",
		r.EgressHostname,
		r.Port,
		r.SAML,
		r.OAuth,
	)
}

// SAMLProvider is the information needed to connect to a SAML IDP provider
type SAMLProvider struct {
	URL      string `yaml:"url" json:"url"`
	Issuer   string `yaml:"issuer" json:"issuer"`
	CertPath string `yaml:"certPath" json:"certPath"`
}

// String returns a string representation of the SAMLProvider
func (sp *SAMLProvider) String() string {
	return fmt.Sprintf("SAML{url:%s, issuer:%s, certPath:%v}",
		sp.URL,
		sp.Issuer,
		sp.CertPath,
	)
}

// OAuthProvider is the information needed to connect to an OAuth provider
type OAuthProvider struct {
	ClientId     string `yaml:"clientId" json:"clientId"`
	ClientSecret string `yaml:"clientSecret" json:"clientSecret"`
}

// String returns a string representation of the OAuthProvider
func (oap *OAuthProvider) String() string {
	return fmt.Sprintf("OAuth{clientId:%s, clientSecret:%s}",
		oap.ClientId,
		oap.ClientSecret,
	)
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

func (c *Config) JSON() string {
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("%v\n", err)
		return ""
	}
	return string(b)
}

func (c *Config) YAML() string {
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("%v\n", err)
		return ""
	}
	return string(b)
}
