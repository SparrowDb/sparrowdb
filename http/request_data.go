package http

import "net/http"

// RequestData holds http.ResponseWriter and http.Request
type RequestData struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	params         []string
}
