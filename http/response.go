package http

import "encoding/json"

// Response HTTP response
type Response struct {
	Database string                 `json:"database"`
	Content  map[string]interface{} `json:"content"`
	Error    []string               `json:"error"`
}

// Value returns query result as json
func (r *Response) Value() []byte {
	b, _ := json.Marshal(r)
	return b
}

// AddErrorStr adds error message as string
func (r *Response) AddErrorStr(text string) {
	r.Error = append(r.Error, text)
}

// AddError adds error message
func (r *Response) AddError(err error) {
	r.Error = append(r.Error, err.Error())
}

// AddContent adds object to query return values
func (r *Response) AddContent(key string, val interface{}) {
	r.Content[key] = val
}

// NewResponse returns new Response
func NewResponse() *Response {
	return &Response{
		Content: make(map[string]interface{}, 0),
		Error:   make([]string, 0),
	}
}
