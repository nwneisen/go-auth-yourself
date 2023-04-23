package handlers

import (
	"io"
	"net/http"

	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/responses"
)

// HandlerWrapper wraps a Handler to create high level calling methods
type HandlerWrapper struct {
	Handler
}

// NewHandlerWrapper creates a new HandlerWrapper
func NewHandlerWrapper(handle Handler) *HandlerWrapper {
	return &HandlerWrapper{
		Handler: handle,
	}
}

// ServeHTTP is the main entry point for the handlers
func (h *HandlerWrapper) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.SetRequest(req)

	// Figure out what type of request it is
	var response *responses.Response
	if h.Request().Method == "GET" {
		response = h.Get()
	} else if h.Request().Method == "POST" {
		response = h.Post()
	} else if h.Request().Method == "PUT" {
		response = h.Put()
	} else if h.Request().Method == "DELETE" {
		response = h.Delete()
	} else {
		msg := "unknown method"
		logger.Info(msg)
		response = responses.BadRequest(msg)
	}

	// Write the response
	w.WriteHeader(response.GetCode())
	io.WriteString(w, response.GetBody())
}

// Handler is the interface for all handlers
type Handler interface {
	Get() *responses.Response
	Post() *responses.Response
	Put() *responses.Response
	Delete() *responses.Response
	SetRequest(request *http.Request)
	Request() *http.Request
}

// BaseHandler is the base struct for all handlers to share methods
type BaseHandler struct {
	request *http.Request
}

// NewBaseHandler creates a new BaseHandler
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

// Get is the default GET method for all handlers
func (h BaseHandler) Get() *responses.Response {
	msg := "GET method not implemented"
	logger.Info(msg)
	return responses.BadRequest(msg)
}

// Post is the default POST method for all handlers
func (h BaseHandler) Post() *responses.Response {
	msg := "POST method not implemented"
	logger.Info(msg)
	return responses.BadRequest(msg)
}

// Put is the default PUT method for all handlers
func (h BaseHandler) Put() *responses.Response {
	msg := "PUT method not implemented"
	logger.Info(msg)
	return responses.BadRequest(msg)
}

// Delete is the default DELETE method for all handlers
func (h BaseHandler) Delete() *responses.Response {
	msg := "DELETE method not implemented"
	logger.Info(msg)
	return responses.BadRequest(msg)
}

// Request returns the request for the handlers
func (h BaseHandler) Request() *http.Request {
	return h.request
}

// SetRequest sets the request for the handlers
func (h *BaseHandler) SetRequest(req *http.Request) {
	h.request = req
}

// // dumpReq is for debugging and sends all of the request data to the browser
// func (h *BaseHandler) dumpReq(w http.ResponseWriter, req *http.Request) {
// 	// values := req.URL.Query()
// 	// if authCode, ok := values["code"]; ok {
// 	// 	h.googleAuthToken(w, authCode[0])
// 	// }

// 	for key, value := range req.Header {
// 		header := fmt.Sprintf("\n%q:%q", key, value[0])
// 		io.WriteString(w, header)
// 	}
// 	io.WriteString(w, "\n\n")

// 	for key, value := range req.URL.Query() {
// 		query := fmt.Sprintf("\n%q:%q", key, value[0])
// 		io.WriteString(w, query)
// 	}
// }
