package spql

import "encoding/json"

// Query holds query parsed content
type Query struct {
	Action string
	Method string
	Params interface{}
}

// QueryRequest holds query parsed content
type QueryRequest struct {
	Action string           `json:"type"`
	Params *json.RawMessage `json:"params"`
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
