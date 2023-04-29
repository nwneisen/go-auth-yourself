package fields

import (
	"encoding/json"
	"fmt"
	"nwneisen/go-proxy-yourself/pkg/logger"

	"gopkg.in/yaml.v2"
)

// EmptyRoot returns a root with default values
func EmptyRoot() *Root {
	return &Root{
		HttpPort:  "80",
		HttpsPort: "443",
		Routes: map[string]*Route{
			"example.com": &Route{
				EgressHostname: "example.com",
				Port:           "80",
				SAML: map[string]*SAMLProvider{
					"example.com": &SAMLProvider{
						URL:      "https://example.com/saml",
						Issuer:   "https://example.com",
						CertPath: "/path/to/cert",
					},
				},
				OAuth: map[string]*OAuthProvider{
					"example.com": &OAuthProvider{
						ClientId:     "example.com",
						ClientSecret: "example.com",
					},
				},
			},
		},
	}
}

// Root top level config layout
type Root struct {
	HttpPort  string `yaml:"httpPort" json:"httpPort" default:"80"`
	HttpsPort string `yaml:"httpsPort" json:"httpsPort" default:"443"`

	Routes map[string]*Route `yaml:"routes" json:"routes"`
}

func (r *Root) JSON() string {
	b, err := json.Marshal(r)
	if err != nil {
		logger.Fatal("%v\n", err)
		return ""
	}
	return string(b)
}

func (r *Root) YAML() string {
	b, err := yaml.Marshal(r)
	if err != nil {
		logger.Fatal("%v\n", err)
		return ""
	}
	return string(b)
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
