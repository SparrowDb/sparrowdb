package spql

import (
	"encoding/json"
	"errors"
	"strings"
)

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

// parseStmt Parse query string to rawQuery
func parseStmt(p *Parser) (*Query, error) {
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
	q, err := parseStmt(p)
	if err != nil {
		return nil, err
	}

	// Parse query params
	perr := p.parse(q, p.rawQuery.Params)
	if perr != nil {
		return nil, perr
	}

	return q, nil
}

func (p *Parser) parse(q *Query, r *json.RawMessage) error {
	var err error

	switch strings.ToLower(q.Action) {
	case "create_database":
		q.Method = "CreateDatabase"
		err = ParseCreateDatabaseStmt(q, r)
	case "drop_database":
		q.Method = "DropDatabase"
		err = ParseDropDatabaseStmt(q, r)
	case "show_databases":
		q.Method = "ShowDatabases"
		err = nil
	case "delete":
		q.Method = "Delete"
		err = ParseDeleteStmt(q, r)
	case "select":
		q.Method = "Select"
		err = ParseSelectStmt(q, r)
	default:
		err = errors.New("Invalid query")
	}

	return err
}

// NewParser returns new Parser
func NewParser(query string) *Parser {
	return &Parser{
		query: query,
	}
}

/*
func ParseString(str string) map[string]interface{} {
	var msgMapTemplate interface{}
	json.Unmarshal([]byte(str), &msgMapTemplate)
	return msgMapTemplate.(map[string]interface{})
}*/
