package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"nwneisen/go-proxy-yourself/internal/fields"
	"nwneisen/go-proxy-yourself/pkg/config"
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

// SAMLHandler handles requests for SAML authentication
func (s SAMLHandler) Get() *responses.Response {
	host := s.Request().Host
	logger.Debug("starting SAML auth for %s", host)

	route, err := config.Route(host)
	if err != nil {
		message := fmt.Sprintf("route not found for %s: %v", host, err)
		logger.Error(message)
		return responses.NotFound(message)
	}

	fmt.Println(route)

	page, err := getRedirectHTML(route)
	if err != nil {
		message := fmt.Sprintf("could not get redirect page: %v", err)
		logger.Error(message)
		return responses.NotFound(message)
	}

	// 	// if req.Referer() != "https://test.nneisen.local/" {
	// 	// 	h.googleOAuthFlow(w, req, route)
	// 	// }
	// }

	// h.dumpReq(w, req)

	logger.Info(fmt.Sprintf("routing from %s to %s", host, route.EgressHostname))
	return responses.OK(page)
}

// CheckForSAML handles requests for SAML authentication
func CheckForSAML(request *http.Request) *responses.Response {
	host := request.Host
	logger.Debug("starting SAML auth for %s", host)

	route, err := config.Route(host)
	if err != nil {
		message := fmt.Sprintf("route not found for %s: %v", host, err)
		logger.Error(message)
		return responses.NotFound(message)
	}

	page, err := getRedirectHTML(route)
	if err != nil {
		message := fmt.Sprintf("could not get redirect page: %v", err)
		logger.Error(message)
		return responses.NotFound(message)
	}

	logger.Info(fmt.Sprintf("routing from %s to %s", host, route.EgressHostname))
	return responses.OK(page)
}

func FinalRedirect(request *http.Request) *responses.Response {
	host := request.Host
	logger.Debug("starting SAML auth for %s", host)

	route, err := config.Route(host)
	if err != nil {
		message := fmt.Sprintf("route not found for %s: %v", host, err)
		logger.Error(message)
		return responses.NotFound(message)
	}

	page, err := getRedirectHTML(route)
	if err != nil {
		message := fmt.Sprintf("could not get redirect page: %v", err)
		logger.Error(message)
		return responses.NotFound(message)
	}

	logger.Info(fmt.Sprintf("routing from %s to %s", host, route.EgressHostname))
	return responses.OK(page)
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

// Post is the default POST method for all handlers
func (s SAMLHandler) Post() *responses.Response {
	msg := "POST to SAML received"
	logger.Info(msg)
	return responses.BadRequest(msg)
}

// getRedirectHTML returns the HTML page that will redirect the user to the service provider
func getRedirectHTML(route *fields.Route) (string, error) {
	htmlPath := "web/client-redirects/okta.html"

	tmpl, err := template.ParseFiles(htmlPath)
	if err != nil {
		return "", fmt.Errorf("could not parse %s: %w", htmlPath, err)
	}

	var doc bytes.Buffer
	tmpl.Execute(&doc, route.SAML)

	return doc.String(), nil
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
