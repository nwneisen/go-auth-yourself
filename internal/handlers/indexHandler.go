package handlers

import (
	"bytes"
	"fmt"
	"html/template"

	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/responses"
	"nwneisen/go-proxy-yourself/pkg/server/handlers"
)

// IndexHandler
type IndexHandler struct {
	*handlers.BaseHandler
}

// NewIndexHandler creates a new handler
func NewIndexHandler() handlers.Handler {
	return IndexHandler{
		BaseHandler: handlers.NewBaseHandler(),
	}
}

// Get returns the index.html page
func (i IndexHandler) Get() *responses.Response {
	logger.Info("Index %s handler called", i.Request().Method)
	indexHTML := "web/index.html"

	response := CheckForSAML(i.Request())
	if response.GetCode() == 200 {
		return response
	}

	// page, err := ioutil.ReadFile("web/index.html")
	// if err != nil {
	// 	msg := fmt.Sprintf("could not read index html file: %v", err.Error())
	// 	logger.Debug(msg)
	// 	return responses.InternalServerError(msg)
	// }

	tmpl, err := template.ParseFiles(indexHTML)
	if err != nil {
		msg := fmt.Sprintf("could not read index html file: %v", err)
		return responses.InternalServerError(msg)
	}

	routes, err := config.Routes()
	if err != nil {
		responses.NotFound(fmt.Sprintf("could not get routes: %v", err))
	}

	var doc bytes.Buffer
	tmpl.Execute(&doc, routes)
	// t, err := template.New("index").ParseFiles("web/index.html")
	// if err != nil {
	// 	msg := fmt.Sprintf("could not read index html file: %v", err.Error())
	// 	logger.Debug(msg)
	// 	return responses.InternalServerError(msg)
	// }
	// err = t.Execute(&doc, i.Config().Routes)
	// if err != nil {
	// 	msg := fmt.Sprintf("could not execute index template: %v", err.Error())
	// 	logger.Debug(msg)
	// 	return responses.InternalServerError(msg)
	// }
	page := doc.String()
	return responses.OK(string(page))
}

// // Index old handler that is no longer really used
// func (i *Index) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	io.WriteString(w, "<a href=/oauth>oauth</a><br>")
// 	io.WriteString(w, "<a href=/saml>saml</a><br>")

// 	// host := req.Host

// 	// if _, ok := h.config.Routes[host]; ok {
// 	// 	// message := fmt.Sprintf("Routing from %s to %s", host, route.EgressHostname)
// 	// 	// logger.Info(w, message)

// 	// 	// h.idpAuthFlow(w, req, route)
// 	// 	logger.Info("Main handler called")

// 	// 	// if req.Referer() != "https://test.nneisen.local/" {
// 	// 	// 	h.googleOAuthFlow(w, req, route)
// 	// 	// }
// 	// }

// 	// h.dumpReq(w, req)
// }

// Post is the default POST method for all handlers
func (h IndexHandler) Post() *responses.Response {

	req := h.Request()
	referers, ok := req.Header["Referer"]
	if !ok {
		return responses.BadRequest("Referer header not found")
	}

	return responses.OK(fmt.Sprintf("referer found: %s", referers[0]))
}
