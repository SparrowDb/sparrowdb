package spql

import (
	"encoding/json"
	"strings"

	"github.com/SparrowDb/sparrowdb/auth"
	"github.com/SparrowDb/sparrowdb/errors"
)

// UserStmt holds user parsed arguments from query
type UserStmt struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetTokenFromRequest extracts token from user query request
func GetTokenFromRequest(queryStr string) (string, bool) {
	p := Parser{query: queryStr}

	if _, err := parseSpql(&p); err != nil {
		return "", false
	}

	return p.rawQuery.Token, true
}

// ValidateUserToken validates user token from request
func ValidateUserToken(str string) bool {
	if rt, got := GetTokenFromRequest(str); got == true {
		return auth.IsLogged(rt)
	}
	return false
}

// ParseUserStmt parses user query and returns QueryObject
func ParseUserStmt(queryStr string) (*QueryResult, error) {
	p := Parser{query: queryStr}
	q, err := parseSpql(&p)
	if err != nil {
		return nil, err
	}

	// Parse query params
	qr, err := parseUser(q, p.rawQuery.Params)
	if err != nil {
		return nil, err
	}

	return qr, nil
}

func parseUser(q *Query, r *json.RawMessage) (*QueryResult, error) {
	var err error

	stmt := &UserStmt{}
	json.Unmarshal(*r, stmt)

	qr := QueryResult{Database: "sparrow_user"}

	switch strings.ToLower(q.Action) {
	case "login":
		{
			u := auth.User{Username: stmt.Username, Password: stmt.Password}
			if logged, token := auth.Authenticate(u); logged == true {
				qr.AddValue(token)
			} else {
				qr.AddErrorStr(errors.ErrLogin.Error())
			}
		}
	}

	return &qr, err
}
