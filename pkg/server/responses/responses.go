package responses

import (
	"fmt"
	"net/http"
)

type Response struct {
	body    string
	code    int
	headers map[string]string
}

func (r *Response) GetBody() string {
	return r.body
}

func (r *Response) GetCode() int {
	return r.code
}

func (r *Response) GetHeaders() map[string]string {
	return r.headers
}

// OK returns a 200 response
func OK(body string) *Response {
	return &Response{
		body:    body,
		code:    http.StatusOK,
		headers: map[string]string{},
	}
}

// TempRedirect returns a 307 response
func TempRedirect(msg string) *Response {
	return &Response{
		body:    fmt.Sprintf("redirecting: %v", msg),
		code:    http.StatusTemporaryRedirect,
		headers: map[string]string{},
	}
}

// BadRequest returns a 400 response
func BadRequest(msg string) *Response {
	return &Response{
		body:    fmt.Sprintf("bad Request: %v", msg),
		code:    http.StatusBadRequest,
		headers: map[string]string{},
	}
}

// NotFound returns a 404 response
func NotFound(msg string) *Response {
	return &Response{
		body:    fmt.Sprintf("not Found: %v", msg),
		code:    http.StatusNotFound,
		headers: map[string]string{},
	}
}

// InternalServerError returns a 500 response
func InternalServerError(msg string) *Response {
	return &Response{
		body:    fmt.Sprintf("internal error: %v", msg),
		code:    http.StatusInternalServerError,
		headers: map[string]string{},
	}
}
