package spql

import (
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	log.Println("TestParser--")

	queryStr := `{
            "type": "create_database",
            "params": {
                "name": "testedb",
                "max_cache_size": "214141"
            }
    }`

	query, err := NewParser(queryStr).ParseQuery()
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return
	}

	qp := query.Params.(*CreateDatabaseStmt)
	log.Printf("SUCCESS: %s", qp)
}
