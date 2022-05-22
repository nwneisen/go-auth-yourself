package saml

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"nwneisen/go-proxy-yourself/pkg/config"
	"time"
)

// OktaProvider for a okta specific provider
type OktaProvider struct{}

// NewOktaProvider constructs a new GoogleProvider structures
func NewOktaProvider() *OktaProvider {
	return &OktaProvider{}
}

// Begin starts the SAML authentication process
func (p *OktaProvider) Begin(w http.ResponseWriter) {
	// page, err := ioutil.ReadFile("web/google-redirect.html")
	// if err != nil {
	// 	return
	// }

	// io.WriteString(w, string(page))
}

func (h *OktaProvider) idpAuthFlow(w http.ResponseWriter, req *http.Request, route config.Route) {
	if value, ok := req.Header["Referer"]; ok {
		// Request came from IDP
		referer := value[0]

		msg := fmt.Sprintf("Referred from %s\n", referer)
		io.WriteString(w, msg)

		// h.addCookie(&w)

		cookie := &http.Cookie{
			Name:   "token",
			Value:  "some_token",
			MaxAge: 300,
		}
		http.SetCookie(w, cookie)
		h.finalRedirect(w, req, &route)
	} else {
		// Request did not come from IDP
		h.idpAuth(w)
	}
}

// idpAuth performs a browser redirect to the identity provider
func (h *OktaProvider) idpAuth(w http.ResponseWriter) {
	page, err := ioutil.ReadFile("web/okta-redirect.html")
	if err != nil {
		// h.logger.Error(err.Error())
		fmt.Println(err.Error())
	}

	io.WriteString(w, string(page))
}

// finalRedirect sends the user to the service provider
func (h *OktaProvider) finalRedirect(w http.ResponseWriter, req *http.Request, route *config.Route) {
	// body, err := ioutil.ReadAll(req.Body)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// // create a new url from the raw RequestURI sent by the client
	url := fmt.Sprintf("%s:%s%s", route.EgressHostname, route.Port, req.RequestURI)

	// // proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))
	// proxyReq, err := http.NewRequest("Get", url, bytes.NewReader(body))

	// // We may want to filter some headers, otherwise we could just use a shallow copy
	// // proxyReq.Header = req.Header
	// proxyReq.Header = make(http.Header)
	// for h, val := range req.Header {
	// 	proxyReq.Header[h] = val
	// }

	client := http.Client{
		Timeout: time.Duration(1) * time.Second,
	}
	client.Get(url)

	// resp, err := client.Do(proxyReq)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadGateway)
	// 	return
	// }
	// defer resp.Body.Close()
	// io.Copy(w, resp.Body)
}
