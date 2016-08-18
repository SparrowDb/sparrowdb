package http

import (
	"bytes"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sparrowdb/db"
	"github.com/sparrowdb/errors"
	"github.com/sparrowdb/model"
	"github.com/sparrowdb/monitor"
	"github.com/sparrowdb/script"
	"github.com/sparrowdb/slog"
	"github.com/sparrowdb/spql"
	"github.com/sparrowdb/util/uuid"
)

// ServeHandler holds main http methods
type ServeHandler struct {
	dbManager     *db.DBManager
	queryExecutor *spql.QueryExecutor
}

func (sh *ServeHandler) writeResponse(request *RequestData, result *spql.QueryResult) {
	request.responseWriter.Header().Set("Content-Type", "application/json")
	request.responseWriter.Write(result.Value())
}

func (sh *ServeHandler) writeError(request *RequestData, query string, errs ...error) {
	result := &spql.QueryResult{}
	for _, v := range errs {
		result.AddErrorStr(v.Error())
	}

	result.AddValue(strings.Replace(query, "\n", "", -1))
	request.responseWriter.WriteHeader(404)
	request.responseWriter.Write(result.Value())
}

func (sh *ServeHandler) serveQuery(request *RequestData) {
	body := request.request.Body

	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	qStr := buf.String()

	p := spql.NewParser(qStr)
	q, err := p.ParseQuery()
	if err != nil {
		sh.writeError(request, qStr, err)
		return
	}

	results := <-sh.queryExecutor.ExecuteQuery(q)

	if results == nil {
		sh.writeError(request, qStr, errors.ErrEmptyQueryResult)
		return
	}

	monitor.IncHTTPQueries()
	sh.writeResponse(request, results)
}

func (sh *ServeHandler) get(request *RequestData) {
	if len(request.params) < 2 {
		sh.writeError(request, "{}", errors.ErrWrongRequest)
		return
	}

	dbname := request.params[0]
	key := request.params[1]

	// Check if database exists
	sto, ok := sh.dbManager.GetDatabase(dbname)
	if !ok {
		sh.writeError(request, "{}", errors.ErrDatabaseNotFound)
		return
	}

	// Async get requested data
	result := <-sh.dbManager.GetData(dbname, key)

	// Check if found requested data or DataDefinition is deleted
	if result == nil || result.Status == model.DataDefinitionRemoved {
		sh.writeError(request, "{}", errors.ErrEmptyQueryResult)
		return
	}

	// Token verification if enabled
	if sto.Descriptor.TokenActive {
		if len(request.params) != 3 {
			sh.writeError(request, "{}", errors.ErrWrongRequest)
			return
		}
		token := request.params[2]

		if token != result.Token {
			sh.writeError(request, "{}", errors.ErrWrongToken)
			return
		}
	}

	request.responseWriter.Header().Set("Content-Type", "image/"+result.Ext)
	request.responseWriter.Header().Set("Content-Length", strconv.Itoa(int(result.Size)))
	request.responseWriter.Write(result.Buf)
}

func (sh *ServeHandler) upload(request *RequestData) {
	file, fhandler, err := request.request.FormFile("uploadfile")
	if err != nil {
		slog.Fatalf(err.Error())
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, file)

	dbname := request.request.FormValue("dbname")
	sto, ok := sh.dbManager.GetDatabase(dbname)

	b := buf.Bytes()

	if ok {
		// checks if user request needs script execution
		if scriptName := request.request.FormValue("script"); len(strings.TrimSpace(scriptName)) > 0 {
			if b, err = script.Execute(scriptName, b); err != nil {
				sh.writeError(request, "{}", errors.ErrInsertImage)
				return
			}
		}

		dataKey := request.request.FormValue("key")
		dataToken := uuid.TimeUUID().String()

		// create new DataDefinition with requested values
		df := &model.DataDefinition{
			Key: dataKey,

			// default store UUID to keep information of insert time
			// and eliminates attacks aimed at guessing valid URLs for photos
			Token: dataToken,

			// get file extension and remove dot before ext name
			Ext: filepath.Ext(fhandler.Filename)[1:],

			Size: uint32(len(b)),

			// Default status 1 (Active)
			Status: model.DataDefinitionActive,

			Buf: b,
		}

		// try to insert image in database
		if err := sto.InsertData(df); err != nil {
			sh.writeError(request, "{}", errors.ErrInsertImage)
			return
		}

		// write ok response
		result := &spql.QueryResult{Database: dbname}
		result.AddValue(df.QueryResult())
		sh.writeResponse(request, result)

		// increment upload statistics
		monitor.IncHTTPUploads()
	} else {
		sh.writeError(request, "{}", errors.ErrDatabaseNotFound)
	}
}

// NewServeHandler returns new ServeHandler
func NewServeHandler(dbm *db.DBManager, queryExecutor *spql.QueryExecutor) *ServeHandler {
	return &ServeHandler{
		dbManager:     dbm,
		queryExecutor: queryExecutor,
	}
}
