package index

import (
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/server/handlers"
	"nwneisen/go-proxy-yourself/pkg/server/responses"
)

// Index handler
type Index struct {
	*handlers.BaseHandler
}

// NewIndex creates a new handler
func NewIndex(config *config.Config, logger *logger.Logger) handlers.Handler {
	return Index{
		BaseHandler: handlers.NewBaseHandler(config, logger),
	}
}

func (i Index) Get() *responses.Response {
	i.Log().Info("Index %s handler called", i.Request().Method)
	return responses.OK("GET method of the index handler")
}

// // Index old handler that is no longer really used
// func (i *Index) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	io.WriteString(w, "<a href=/oauth>oauth</a><br>")
// 	io.WriteString(w, "<a href=/saml>saml</a><br>")

// 	// host := req.Host

// 	// if _, ok := h.config.Routes[host]; ok {
// 	// 	// message := fmt.Sprintf("Routing from %s to %s", host, route.EgressHostname)
// 	// 	// log.Println(w, message)

// 	// 	// h.idpAuthFlow(w, req, route)
// 	// 	h.logger.Info("Main handler called")

// 	// 	// if req.Referer() != "https://test.nneisen.local/" {
// 	// 	// 	h.googleOAuthFlow(w, req, route)
// 	// 	// }
// 	// }

// 	// h.dumpReq(w, req)
// }
