package spql

import "encoding/json"

// QueryResult holds the result of statement
type QueryResult struct {
	Database string        `json:"database"`
	Values   []interface{} `json:"values"`
	Error    []string      `json:"error"`
}

func (qr *QueryResult) Value() []byte {
	b, _ := json.Marshal(qr)
	return b
}

func (qr *QueryResult) AddErrorStr(text string) {
	qr.Error = append(qr.Error, text)
}

func (qr *QueryResult) AddValue(i interface{}) {
	qr.Values = append(qr.Values, i)
}

func NewQueryResult() *QueryResult {
	return &QueryResult{}
}
