package index

import (
	"bytes"
	"fmt"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/server/handlers"
	"nwneisen/go-proxy-yourself/pkg/server/responses"
	"text/template"
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

// Get returns the index.html page
func (i Index) Get() *responses.Response {
	i.Log().Info("Index %s handler called", i.Request().Method)

	// page, err := ioutil.ReadFile("web/index.html")
	// if err != nil {
	// 	msg := fmt.Sprintf("could not read index html file: %v", err.Error())
	// 	i.Log().Error(msg)
	// 	return responses.InternalServerError(msg)
	// }

	tmpl, err := template.ParseFiles("web/index.html")
	if err != nil {
		msg := fmt.Sprintf("could not read index html file: %v", err.Error())
		i.Log().Error(msg)
		return responses.InternalServerError(msg)
	}

	var doc bytes.Buffer
	// tmpl.Execute(&doc, i.Config().Routes)
	t, err := template.New("index").ParseFiles("web/index.html")
	err = t.Execute(&doc, i.Config().Routes)
	page := doc.String()

	routes := fmt.Sprintf("%+v", page)

	return responses.OK(string(routes))
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
