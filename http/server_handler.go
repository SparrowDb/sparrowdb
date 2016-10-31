package http

import (
	"bytes"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/SparrowDb/sparrowdb/auth"
	"github.com/SparrowDb/sparrowdb/db"
	"github.com/SparrowDb/sparrowdb/errors"
	"github.com/SparrowDb/sparrowdb/model"
	"github.com/SparrowDb/sparrowdb/monitor"
	"github.com/SparrowDb/sparrowdb/script"
	"github.com/SparrowDb/sparrowdb/slog"
	"github.com/SparrowDb/sparrowdb/spql"
	"github.com/SparrowDb/sparrowdb/util/uuid"
)

// ServeHandler holds main http methods
type ServeHandler struct {
	dbManager     *db.DBManager
	queryExecutor *spql.QueryExecutor
}

const (
	statusOk         = 200
	statusBadRequest = 400
	statusNotFound   = 404
	statusConflict   = 409
)

func (sh *ServeHandler) writeResponse(status int, request *RequestData, result *spql.QueryResult) {
	request.responseWriter.Header().Set("Content-Type", "application/json")
	request.responseWriter.WriteHeader(status)
	request.responseWriter.Write(result.Value())
}

func (sh *ServeHandler) user(request *RequestData) {
	body := request.request.Body

	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	qStr := buf.String()

	qr, err := spql.ParseUserStmt(qStr)
	if err != nil {
		qr.AddErrorStr(err.Error())
		sh.writeResponse(statusBadRequest, request, qr)
		return
	}
	sh.writeResponse(statusOk, request, qr)
}

func (sh *ServeHandler) serveQuery(request *RequestData) {
	body := request.request.Body

	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	qStr := buf.String()

	results := &spql.QueryResult{}

	if sh.dbManager.Config.AuthenticationActive {
		userToken, _ := spql.GetTokenFromRequest(qStr)
		if !auth.IsLogged(userToken) {
			results.AddErrorStr(errors.ErrInvalidToken.Error())
			sh.writeResponse(statusBadRequest, request, results)
			return
		}
	}

	p := spql.NewParser(qStr)
	q, err := p.ParseQuery()
	if err != nil {
		results.AddErrorStr(err.Error())
		sh.writeResponse(statusBadRequest, request, results)
		return
	}

	results = <-sh.queryExecutor.ExecuteQuery(q)

	if results == nil {
		results.AddErrorStr(errors.ErrEmptyQueryResult.Error())
		sh.writeResponse(statusBadRequest, request, results)
		return
	}

	monitor.IncHTTPQueries()
	sh.writeResponse(statusOk, request, results)
}

func (sh *ServeHandler) get(request *RequestData) {
	results := &spql.QueryResult{}

	if len(request.params) < 2 {
		results.AddErrorStr(errors.ErrWrongRequest.Error())
		sh.writeResponse(statusBadRequest, request, results)
		return
	}

	dbname := request.params[0]
	key := request.params[1]

	// Check if database exists
	sto, ok := sh.dbManager.GetDatabase(dbname)
	if !ok {
		results.AddErrorStr(errors.ErrDatabaseNotFound.Error())
		sh.writeResponse(statusBadRequest, request, results)
		return
	}

	// Async get requested data
	result := <-sh.dbManager.GetData(dbname, key)

	// Check if found requested data or DataDefinition is deleted
	if result == nil || result.Status == model.DataDefinitionRemoved {
		results.AddErrorStr(errors.ErrEmptyQueryResult.Error())
		sh.writeResponse(statusBadRequest, request, results)
		return
	}

	// Token verification if enabled
	if sto.Descriptor.TokenActive {
		if len(request.params) != 3 {
			results.AddErrorStr(errors.ErrWrongRequest.Error())
			sh.writeResponse(statusBadRequest, request, results)
			return
		}
		token := request.params[2]

		if token != result.Token {
			results.AddErrorStr(errors.ErrWrongToken.Error())
			sh.writeResponse(statusBadRequest, request, results)
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
		slog.Errorf(err.Error())
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, file)

	dbname := request.request.FormValue("dbname")
	sto, ok := sh.dbManager.GetDatabase(dbname)

	b := buf.Bytes()
	results := &spql.QueryResult{}

	if !ok {
		results.AddErrorStr(errors.ErrDatabaseNotFound.Error())
		sh.writeResponse(statusBadRequest, request, results)
		return
	}
	results.Database = dbname
	dataKey := request.request.FormValue("key")

	var dataRev uint32

	if _rev := request.request.FormValue("rev"); len(strings.TrimSpace(_rev)) > 0 {
		_dataRev, err := strconv.Atoi(_rev)
		if err != nil {
			results.AddErrorStr(err.Error())
			sh.writeResponse(statusBadRequest, request, results)
			return
		}
		dataRev = uint32(_dataRev)
	}

	// checks if user request needs script execution
	if scriptName := request.request.FormValue("script"); len(strings.TrimSpace(scriptName)) > 0 {
		if b, err = script.Execute(scriptName, b); err != nil {
			results.AddErrorStr(errors.ErrInsertImage.Error())
			sh.writeResponse(statusBadRequest, request, results)
			return
		}
	}

	if isValidKey := spql.ValidateDatabaseName.MatchString(dataKey); !isValidKey {
		results.AddErrorStr(errors.ErrImageInvalidKey.Error())
		sh.writeResponse(statusBadRequest, request, results)
		return
	}

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

		Revision: 0,

		Version: make([]uint32, 0),

		Buf: b,
	}

	// try to insert image in database
	if _, err := sto.InsertCheckRevision(df, dataRev); err != nil {
		results.AddErrorStr(err.Error())
		sh.writeResponse(statusConflict, request, results)
		return
	}

	// write ok response
	results.AddValue(df.QueryResult())
	sh.writeResponse(statusOk, request, results)

	// increment upload statistics
	monitor.IncHTTPUploads()
}

// NewServeHandler returns new ServeHandler
func NewServeHandler(dbm *db.DBManager, queryExecutor *spql.QueryExecutor) *ServeHandler {
	return &ServeHandler{
		dbManager:     dbm,
		queryExecutor: queryExecutor,
	}
}
