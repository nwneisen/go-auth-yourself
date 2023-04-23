package handlers

import (
	"io"
	"io/ioutil"
	"net/http"

	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/responses"
	"nwneisen/go-proxy-yourself/pkg/server/handlers"
)

type SAMLHandler struct {
	*handlers.BaseHandler
}

func NewSamlHandler() handlers.Handler {
	return SAMLHandler{
		BaseHandler: handlers.NewBaseHandler(),
	}
}

// Index
func (s SAMLHandler) Get() *responses.Response {
	// host := s.Request().Host

	// if _, ok := config.Routes[host]; ok {
	// 	// message := fmt.Sprintf("Routing from %s to %s", host, route.EgressHostname)
	// 	// logger.Ifno(w, message)

	// 	// h.idpAuthFlow(w, req, route)
	// 	logger.Info("Saml route called")

	// 	// if req.Referer() != "https://test.nneisen.local/" {
	// 	// 	h.googleOAuthFlow(w, req, route)
	// 	// }
	// }

	// h.dumpReq(w, req)
	return responses.OK("Saml route called")
}

// func (s SAMLHandler) idpAuthFlow(w http.ResponseWriter, req *http.Request, route config.Route) {
// 	if value, ok := req.Header["Referer"]; ok {
// 		// Request came from IDP
// 		referer := value[0]

// 		msg := fmt.Sprintf("Referred from %s\n", referer)
// 		io.WriteString(w, msg)

// 		// h.addCookie(&w)

// 		logger.Info("Adding cookie")
// 		cookie := &http.Cookie{
// 			Name:   "token",
// 			Value:  "some_token",
// 			MaxAge: 300,
// 		}
// 		http.SetCookie(w, cookie)
// 		// s.finalRedirect(w, req, &route)
// 	} else {
// 		// Request did not come from IDP
// 		s.idpAuth(w)
// 	}
// }

// idpAuth performs a browser redirect to the identity provider
func (s SAMLHandler) idpAuth(w http.ResponseWriter) {
	logger.Info("Checking auth with IDP")
	page, err := ioutil.ReadFile("web/okta-redirect.html")
	if err != nil {
		logger.Error(err.Error())
	}

	io.WriteString(w, string(page))
}

// finalRedirect sends the user to the service provider
// func (s SAMLHandler) finalRedirect(w http.ResponseWriter, req *http.Request, route *config.Route) {
// 	logger.Info("Sending final redirect")

// 	// body, err := ioutil.ReadAll(req.Body)
// 	// if err != nil {
// 	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	// // create a new url from the raw RequestURI sent by the client
// 	url := fmt.Sprintf("%s:%s%s", route.EgressHostname, route.Port, req.RequestURI)

// 	// // proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))
// 	// proxyReq, err := http.NewRequest("Get", url, bytes.NewReader(body))

// 	// // We may want to filter some headers, otherwise we could just use a shallow copy
// 	// // proxyReq.Header = req.Header
// 	// proxyReq.Header = make(http.Header)
// 	// for h, val := range req.Header {
// 	// 	proxyReq.Header[h] = val
// 	// }

// 	client := http.Client{
// 		Timeout: time.Duration(1) * time.Second,
// 	}
// 	client.Get(url)

// 	// resp, err := client.Do(proxyReq)
// 	// if err != nil {
// 	// 	http.Error(w, err.Error(), http.StatusBadGateway)
// 	// 	return
// 	// }
// 	// defer resp.Body.Close()
// 	// io.Copy(w, resp.Body)
// }
