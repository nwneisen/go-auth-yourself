package saml

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/server/handlers"
	"nwneisen/go-proxy-yourself/pkg/server/responses"
	"time"
)

type SAML struct {
	*handlers.BaseHandler
}

func NewSaml(config *config.Config, logger *logger.Logger) handlers.Handler {
	return SAML{
		BaseHandler: handlers.NewBaseHandler(config, logger),
	}
}

// Index
func (s SAML) Get() *responses.Response {
	host := s.Request().Host

	if _, ok := s.Config().Routes[host]; ok {
		// message := fmt.Sprintf("Routing from %s to %s", host, route.EgressHostname)
		// log.Println(w, message)

		// h.idpAuthFlow(w, req, route)
		s.Log().Info("Saml route called")

		// if req.Referer() != "https://test.nneisen.local/" {
		// 	h.googleOAuthFlow(w, req, route)
		// }
	}

	// h.dumpReq(w, req)
	return responses.OK("Saml route called")
}

func (s SAML) idpAuthFlow(w http.ResponseWriter, req *http.Request, route config.Route) {
	if value, ok := req.Header["Referer"]; ok {
		// Request came from IDP
		referer := value[0]

		msg := fmt.Sprintf("Referred from %s\n", referer)
		io.WriteString(w, msg)

		// h.addCookie(&w)

		s.Log().Info("Adding cookie")
		cookie := &http.Cookie{
			Name:   "token",
			Value:  "some_token",
			MaxAge: 300,
		}
		http.SetCookie(w, cookie)
		s.finalRedirect(w, req, &route)
	} else {
		// Request did not come from IDP
		s.idpAuth(w)
	}
}

// idpAuth performs a browser redirect to the identity provider
func (s SAML) idpAuth(w http.ResponseWriter) {
	s.Log().Info("Checking auth with IDP")
	page, err := ioutil.ReadFile("web/okta-redirect.html")
	if err != nil {
		s.Log().Error(err.Error())
	}

	io.WriteString(w, string(page))
}

// finalRedirect sends the user to the service provider
func (s SAML) finalRedirect(w http.ResponseWriter, req *http.Request, route *config.Route) {
	s.Log().Info("Sending final redirect")

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
