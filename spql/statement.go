package spql

import (
	"encoding/json"
	"strings"

	"github.com/SparrowDb/sparrowdb/errors"
)

var (
	commandTable map[string]func(*Query, *json.RawMessage) error
)

func init() {
	commandTable = make(map[string]func(*Query, *json.RawMessage) error)

	commandTable["create_database"] = ParseCreateDatabaseStmt
	commandTable["drop_database"] = ParseDropDatabaseStmt
	commandTable["show_databases"] = ParseShowDatabaseStmt
	commandTable["delete"] = ParseDeleteStmt
	commandTable["select"] = ParseSelectStmt
	commandTable["create_snapshot"] = ParseCreateSnapshotStmt
}

func parseStmt(q *Query, r *json.RawMessage) error {
	var err error
	cmd := strings.ToLower(q.Action)

	if f, ok := commandTable[cmd]; ok == true {
		err = f(q, r)
	} else {
		err = errors.ErrInvalidQueryAction
	}

	return err
}
