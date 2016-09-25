package spql

import (
	"fmt"
	"reflect"

	"github.com/sparrowdb/backup"
	"github.com/sparrowdb/db"
	"github.com/sparrowdb/errors"
	"github.com/sparrowdb/model"
	"github.com/sparrowdb/slog"
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
	defer func() {
		if r := recover(); r != nil {
			slog.Errorf("%s", r)
		}
	}()

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
		SnapshotPath:   qp.SnapshotPath,
	}

	err := qe.dbManager.CreateDatabase(databaseCfg)
	qr := QueryResult{Database: databaseCfg.Name}

	if err != nil {
		qr.AddErrorStr(err.Error())
	}

	results <- &qr
}

// DropDatabase process drop database from results <- qrquery string
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
			qr.AddErrorStr(fmt.Sprintf(errDataNotFound, qp.Key, qp.Name))
		} else {
			tbs := model.NewTombstone(result)
			db.InsertData(tbs)
		}
	} else {
		qr.AddErrorStr(errors.ErrDatabaseNotFound.Error())
	}

	results <- &qr
}

// Select do query in database
func (qe *QueryExecutor) Select(query *Query, results chan *QueryResult) {
	qp := query.Params.(*SelectStmt)
	qr := QueryResult{Database: qp.Name}

	if db, ok := qe.dbManager.GetDatabase(qp.Name); ok {
		qe.doSelect(qp, &qr, db, results)
	} else {
		qr.AddErrorStr(errors.ErrDatabaseNotFound.Error())
		results <- &qr
	}
}

func (qe *QueryExecutor) doSelect(qp *SelectStmt, qr *QueryResult, db *db.Database, result chan *QueryResult) {
	// empty means query all
	if qp.Key == "" {

	} else {
		if d, ok := db.GetDataByKey(qp.Key); ok {
			qr.AddValue(d.QueryResult())
			result <- qr
		}
	}
}

// CreateSnapshot process to create snapshot of database
func (qe *QueryExecutor) CreateSnapshot(query *Query, results chan *QueryResult) {
	qp := query.Params.(*CreateSnapshotStmt)

	//err := qe.dbManager.DropDatabase(qp.Name)
	qr := QueryResult{Database: qp.Name}

	if db, ok := qe.dbManager.GetDatabase(qp.Name); ok == true {
		err := backup.CreateSnapshot(db.Descriptor.Path, db.Descriptor.SnapshotPath)
		if err != nil {
			qr.AddErrorStr(errors.ErrCreateDatabase.Error())
		}
	} else {
		qr.AddErrorStr(errors.ErrDatabaseNotFound.Error())
	}

	results <- &qr
}

// NewQueryExecutor returns new QueryExecutor
func NewQueryExecutor(dbm *db.DBManager) *QueryExecutor {
	return &QueryExecutor{
		dbManager: dbm,
	}
}
