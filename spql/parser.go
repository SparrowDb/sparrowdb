package spql

import "encoding/json"

// Parser holds query parser definitions
type Parser struct {
	query    string
	rawQuery struct {
		Token   string           `json:"token"`
		Action  string           `json:"type"`
		Params  *json.RawMessage `json:"params"`
		Filters *json.RawMessage `json:"filters"`
	}
}

// parseSpql Parse query string to rawQuery
func parseSpql(p *Parser) (*Query, error) {
	// Parse query string to rawQuery
	err := json.Unmarshal([]byte(p.query), &p.rawQuery)
	if err != nil {
		return nil, err
	}

	query := Query{
		Action: p.rawQuery.Action,
	}

	return &query, nil
}

// ParseQuery parses a query string and returns QueryObject
func (p *Parser) ParseQuery() (*Query, error) {
	q, err := parseSpql(p)
	if err != nil {
		return nil, err
	}

	// Parse query params
	perr := parseStmt(q, p.rawQuery.Params)
	if perr != nil {
		return nil, perr
	}

	return q, nil
}

// NewParser returns new Parser
func NewParser(query string) *Parser {
	return &Parser{
		query: query,
	}
}
