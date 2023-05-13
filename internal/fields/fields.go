package fields

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"nwneisen/go-proxy-yourself/pkg/logger"
)

// EmptyRoot returns a root with default values
func EmptyRoot() *Root {
	return &Root{
		HttpPort:  "80",
		HttpsPort: "443",
		Routes: map[string]*Route{
			"simple.app": {
				EgressHostname: "localhost",
				Port:           "8081",
				SAML: map[string]*SAMLProvider{
					"test": {
						URL: "https://localhost:8443/saml",
					},
				},
				OAuth: map[string]*OAuthProvider{
					"google": {},
				},
			},
		},
	}
}

// Root top level config layout
type Root struct {
	HttpPort   string `yaml:"httpPort" json:"httpPort" default:"80"`
	HttpsPort  string `yaml:"httpsPort" json:"httpsPort" default:"443"`
	ServerCert string `yaml:"serverCert" json:"serverCert" default:"bin/server.cert"`
	ServerKey  string `yaml:"serverKey" json:"serverKey" default:"bin/server.key"`

	Routes map[string]*Route `yaml:"routes" json:"routes"`
}

// String returns a string representation of the Root
func (r *Root) String() string {
	return fmt.Sprintf("Root{httpPort:%s, httpsPort:%s, serverCert:%s, serverKey:%s, routes:%v}",
		r.HttpPort,
		r.HttpsPort,
		r.ServerCert,
		r.ServerKey,
		r.Routes,
	)
}

// JSON returns a JSON representation of the Root
func (r *Root) JSON() string {
	b, err := json.Marshal(r)
	if err != nil {
		logger.Fatal("%v\n", err)
		return ""
	}
	return string(b)
}

// UnmarshalJSON unmarshals the root config values
func (r *Root) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, r)
	if err != nil {
		logger.Fatal("%v\n", err)
	}

	return nil
}

// YAML returns a YAML representation of the Root
func (r *Root) YAML() string {
	b, err := yaml.Marshal(r)
	if err != nil {
		logger.Fatal("%v\n", err)
		return ""
	}
	return string(b)
}

// UnmarshalYAML unmarshals the root config values
func (r *Root) UnmarshalYAML(bytes []byte) error {
	err := yaml.Unmarshal(bytes, r)
	if err != nil {
		logger.Fatal("%v\n", err)
	}

	return nil
}

// Validate the root config values
func (r *Root) Validate() error {
	// Check if the ServerCert file exists
	if _, err := os.Stat(r.ServerCert); err != nil {
		return fmt.Errorf("server cert does not exist at %s: %s", r.ServerCert, err)
	}

	// Check if the ServerKey file exists
	if _, err := os.Stat(r.ServerKey); err != nil {
		return fmt.Errorf("server key does not exist at %s: %s", r.ServerKey, err)
	}

	return nil
}

// Possible redirect routes and their information
type Route struct {
	EgressHostname string `yaml:"egressHostname" json:"egressHostname" default:"localhost"`
	Port           string `yaml:"port" json:"port" default:"8443"`

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
	Metadata string `yaml:"metadata" json:"metadata"`
}

// String returns a string representation of the SAMLProvider
func (sp *SAMLProvider) String() string {
	return fmt.Sprintf("SAML{url:%s, issuer:%s, certPath:%v, metadata:%v}",
		sp.URL,
		sp.Issuer,
		sp.CertPath,
		sp.Metadata,
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
