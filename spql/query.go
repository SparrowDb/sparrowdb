package spql

import "encoding/json"

// Query holds query parsed content
type Query struct {
	Action  string
	Method  string
	Params  interface{}
	Filters map[string]string
}

// QueryRequest holds query parsed content
type QueryRequest struct {
	Action  string            `json:"type"`
	Params  *json.RawMessage  `json:"params"`
	Filters map[string]string `json:"filters"`
}

// ParseQuery parses a query string and returns QueryObject
func (qr *QueryRequest) ParseQuery() (Query, error) {
	q := Query{}

	q.Action = qr.Action

	if perr := parseStmt(&q, qr.Params); perr != nil {
		return q, perr
	}

	return q, nil
}
