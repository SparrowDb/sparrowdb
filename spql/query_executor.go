package spql

import (
	"fmt"
	"reflect"

	"github.com/sparrowdb/db"
	"github.com/sparrowdb/model"
)

// QueryExecutor holds query executor data
type QueryExecutor struct {
	dbManager *db.DBManager
}

// ExecuteQuery execute query in channel
func (qe *QueryExecutor) ExecuteQuery(query *Query) <-chan *QueryResult {
	results := make(chan *QueryResult)
	go qe.executeQuery(query, results)
	return results
}

func (qe *QueryExecutor) executeQuery(query *Query, results chan *QueryResult) {
	defer close(results)

	inputs := make([]reflect.Value, 2)
	inputs[0] = reflect.ValueOf(query)
	inputs[1] = reflect.ValueOf(results)

	reflect.ValueOf(qe).MethodByName(query.Method).Call(inputs)
}

// CreateDatabase process create database from query string
func (qe *QueryExecutor) CreateDatabase(query *Query, results chan *QueryResult) {
	qp := query.Params.(*CreateDatabaseStmt)

	databaseCfg := db.DatabaseDescriptor{
		Name:           qp.Name,
		MaxDataLogSize: qp.MaxDataLogSize,
		MaxCacheSize:   qp.MaxCacheSize,
		BloomFilterFp:  qp.BloomFilterFp,
		CronExp:        qp.CronExp,
		Path:           qp.Path,
	}

	err := qe.dbManager.CreateDatabase(databaseCfg)
	qr := QueryResult{Database: databaseCfg.Name}

	if err != nil {
		qr.AddErrorStr(err.Error())
	}

	results <- &qr
}

// DropDatabase process drop database from query string
func (qe *QueryExecutor) DropDatabase(query *Query, results chan *QueryResult) {
	qp := query.Params.(*DropDatabaseStmt)

	err := qe.dbManager.DropDatabase(qp.Name)
	qr := QueryResult{Database: qp.Name}

	if err != nil {
		qr.AddErrorStr(err.Error())
	}

	results <- &qr
}

// ShowDatabases process show databases from query string
func (qe *QueryExecutor) ShowDatabases(query *Query, results chan *QueryResult) {
	n := qe.dbManager.GetDatabasesNames()
	qr := QueryResult{}

	for _, v := range n {
		qr.AddValue(v)
	}

	results <- &qr
}

// Delete delets entry from database with tombstone
func (qe *QueryExecutor) Delete(query *Query, results chan *QueryResult) {
	qp := query.Params.(*DeleteStmt)
	qr := QueryResult{}

	if db, ok := qe.dbManager.GetDatabase(qp.Name); ok == true {
		result := <-qe.dbManager.GetData(qp.Name, qp.Key)

		// Check if found requested data or DataDefinition is deleted
		if result == nil || result.Status == model.DataDefinitionRemoved {
			qr.AddErrorStr(fmt.Sprintf("Image %s not found in %s", qp.Key, qp.Name))
		} else {
			tbs := model.NewTombstone(result)
			db.InsertData(tbs)
		}
	} else {
		qr.AddErrorStr(fmt.Sprintf("Database %s not found", qp.Name))
	}

	results <- &qr
}

// NewQueryExecutor returns new QueryExecutor
func NewQueryExecutor(dbm *db.DBManager) *QueryExecutor {
	return &QueryExecutor{
		dbManager: dbm,
	}
}
