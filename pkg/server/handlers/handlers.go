package handlers

import (
	"io"
	"net/http"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/server/responses"
)

type HandlerWrapper struct {
	config *config.Config
	logger *logger.Logger
	next   Handler

	request *http.Request
}

func NewHandlerWrapper(config *config.Config, logger *logger.Logger, handle *Handler) *HandlerWrapper {
	return &HandlerWrapper{
		config: config,
		logger: logger,
		next:   *handle,
	}
}

func (h *HandlerWrapper) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.request = req

	var response *responses.Response
	if h.request.Method == "GET" {
		response = h.next.Get()
	} else if h.request.Method == "POST" {
		response = h.next.Post()
	} else if h.request.Method == "PUT" {
		response = h.next.Put()
	} else if h.request.Method == "DELETE" {
		response = h.next.Delete()
	} else {
		msg := "unknown method"
		h.logger.Error(msg)
		response = responses.BadRequest(msg)
	}

	w.WriteHeader(response.GetCode())
	io.WriteString(w, response.GetBody())
}

type Handler interface {
	Get() *responses.Response
	Post() *responses.Response
	Put() *responses.Response
	Delete() *responses.Response
}

type BaseHandler struct {
	config *config.Config
	logger *logger.Logger
}

func NewBaseHandler(config *config.Config, logger *logger.Logger) *BaseHandler {
	return &BaseHandler{
		config: config,
		logger: logger,
	}
}

func (h *BaseHandler) Get() *responses.Response {
	msg := "GET method not implemented"
	h.logger.Info(msg)
	return responses.BadRequest(msg)
}

func (h *BaseHandler) Post() *responses.Response {
	msg := "POST method not implemented"
	h.logger.Error(msg)
	return responses.BadRequest(msg)
}

func (h *BaseHandler) Put() *responses.Response {
	msg := "PUT method not implemented"
	h.logger.Error(msg)
	return responses.BadRequest(msg)
}

func (h *BaseHandler) Delete() *responses.Response {
	msg := "DELETE method not implemented"
	h.logger.Error(msg)
	return responses.BadRequest(msg)
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
