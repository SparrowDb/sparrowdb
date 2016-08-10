package spql

import "encoding/json"

// QueryResult holds the result of statement
type QueryResult struct {
	Database string        `json:"database"`
	Values   []interface{} `json:"values"`
	Error    []string      `json:"error"`
}

// Value returns query result as json
func (qr *QueryResult) Value() []byte {
	b, _ := json.Marshal(qr)
	return b
}

// AddErrorStr adds error message as string
func (qr *QueryResult) AddErrorStr(text string) {
	qr.Error = append(qr.Error, text)
}

// AddValue adds object to query return values
func (qr *QueryResult) AddValue(i interface{}) {
	qr.Values = append(qr.Values, i)
}

// NewQueryResult returns new QueryResult
func NewQueryResult() *QueryResult {
	return &QueryResult{}
}
