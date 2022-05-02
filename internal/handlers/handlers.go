package handlers

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
)

// Handlers is a generic handler for none specific routes
type Handlers struct {
	config *config.Config
	logger *logger.Logger
}

// NewHandlers creates a new handler
func NewHandlers(config *config.Config, logger *logger.Logger) *Handlers {
	return &Handlers{config, logger}
}

// Index old handler that is no longer really used
func (h *Handlers) Index(w http.ResponseWriter, req *http.Request) {
	host := req.Host

	if _, ok := h.config.Routes[host]; ok {
		// message := fmt.Sprintf("Routing from %s to %s", host, route.EgressHostname)
		// log.Println(w, message)

		// h.idpAuthFlow(w, req, route)
		h.logger.Info("Main handler called")

		// if req.Referer() != "https://test.nneisen.local/" {
		// 	h.googleOAuthFlow(w, req, route)
		// }
	}

	h.dumpReq(w, req)
}

// dumpReq is for debugging and sends all of the request data to the browser
func (h *Handlers) dumpReq(w http.ResponseWriter, req *http.Request) {
	// values := req.URL.Query()
	// if authCode, ok := values["code"]; ok {
	// 	h.googleAuthToken(w, authCode[0])
	// }

	for key, value := range req.Header {
		header := fmt.Sprintf("\n%q:%q", key, value[0])
		io.WriteString(w, header)
	}
	io.WriteString(w, "\n\n")

	for key, value := range req.URL.Query() {
		query := fmt.Sprintf("\n%q:%q", key, value[0])
		io.WriteString(w, query)
	}
}

// RedirectToHTTPS sends all HTTP requests to HTTPS
func (h *Handlers) RedirectToHTTPS() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		host, _, err := net.SplitHostPort(req.Host)
		if err != nil {
			log.Println(err)
			host = req.Host
		}
		target := fmt.Sprintf("https://%s:%s%s", host, h.config.HttpsPort, req.URL.Path)

		if len(req.URL.RawQuery) > 0 {
			target += "?" + req.URL.RawQuery
		}
		h.logger.Info("redirect to: %s", target)
		http.Redirect(w, req, target,
			// see comments below and consider the codes 308, 302, or 301
			http.StatusTemporaryRedirect)
	}
}
