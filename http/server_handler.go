package http

import (
	"bytes"
	"errors"
	"io"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sparrowdb/db"
	"github.com/sparrowdb/model"
	"github.com/sparrowdb/monitor"
	"github.com/sparrowdb/spql"
	"github.com/sparrowdb/util/uuid"
)

// ServeHandler holds main http methods
type ServeHandler struct {
	dbManager     *db.DBManager
	queryExecutor *spql.QueryExecutor
}

var (
	errDatabaseNotFound = errors.New("Database not found")
	errWrongRequest     = errors.New("Wrong HTTP request")
	errEmptyQueryResult = errors.New("Empty query result")
	errWrongToken       = errors.New("Wrong token")
	errInsertImage      = errors.New("Could not insert images")
)

func (sh *ServeHandler) writeResponse(request *RequestData, result *spql.QueryResult) {
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
		sh.writeError(request, qStr, errEmptyQueryResult)
		return
	}

	monitor.IncHTTPQueries()
	request.responseWriter.Header().Set("Content-Type", "application/json")
	sh.writeResponse(request, results)
}

func (sh *ServeHandler) get(request *RequestData) {
	if len(request.params) < 2 {
		sh.writeError(request, "{}", errWrongRequest)
		return
	}

	dbname := request.params[0]
	key := request.params[1]

	// Check if database exists
	sto, ok := sh.dbManager.GetDatabase(dbname)
	if !ok {
		sh.writeError(request, "{}", errDatabaseNotFound)
		return
	}

	// Async get requested data
	result := <-sh.dbManager.GetData(dbname, key)

	// Check if found requested data
	if result == nil {
		sh.writeError(request, "{}", errEmptyQueryResult)
		return
	}

	// Token verification if enabled
	if sto.Descriptor.TokenActive {
		if len(request.params) != 3 {
			sh.writeError(request, "{}", errWrongRequest)
			return
		}
		token := request.params[2]

		if token != result.Token {
			sh.writeError(request, "{}", errWrongToken)
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
		log.Println(err)
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, file)

	dbname := request.request.FormValue("dbname")
	sto, ok := sh.dbManager.GetDatabase(dbname)

	if ok {
		var token string

		if sto.Descriptor.TokenActive {
			token = uuid.TimeUUID().String()
		}

		err := sto.InsertData(&model.DataDefinition{
			Key:   request.request.FormValue("key"),
			Token: token,

			// get file extension and remove dot before ext name
			Ext: filepath.Ext(fhandler.Filename)[1:],

			Size: uint32(len(buf.Bytes())),
			Buf:  buf.Bytes(),
		})

		if err != nil {
			sh.writeError(request, "{}", errInsertImage)
		}

		monitor.IncHTTPUploads()
	}
}

// NewServeHandler returns new ServeHandler
func NewServeHandler(dbm *db.DBManager, queryExecutor *spql.QueryExecutor) *ServeHandler {
	return &ServeHandler{
		dbManager:     dbm,
		queryExecutor: queryExecutor,
	}
}
