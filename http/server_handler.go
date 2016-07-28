package http

import (
	"bytes"
	"io"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sparrowdb/db"
	"github.com/sparrowdb/model"
	"github.com/sparrowdb/monitor"
	"github.com/sparrowdb/spql"
)

// ServeHandler holds main http methods
type ServeHandler struct {
	dbManager     *db.DBManager
	queryExecutor *spql.QueryExecutor
}

func (sh *ServeHandler) writeResponse(request *RequestData, result *spql.QueryResult) {
	request.responseWriter.Write(result.Value())
}

func (sh *ServeHandler) writeError(request *RequestData, query string, errs ...error) {
	result := &spql.QueryResult{}
	for _, v := range errs {
		result.AddErrorStr(v.Error())
	}
	result.AddValue(strings.Replace(query, "\n", "", -1))
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
		log.Fatalf("ERROR: Nil query result")
		return
	}

	monitor.IncHTTPQueries()
	request.responseWriter.Header().Set("Content-Type", "application/json")
	sh.writeResponse(request, results)
}

func (sh *ServeHandler) get(request *RequestData) {
	if len(request.params) != 2 {
		request.responseWriter.WriteHeader(404)
		return
	}

	dbname := request.params[0]
	key := request.params[1]

	result := <-sh.dbManager.GetData(dbname, key)

	if result == nil {
		request.responseWriter.WriteHeader(404)
		request.responseWriter.Write([]byte("ERROOOOOOOOOOOOOOOOOOOOOOOOO"))
		return
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
		sto.InsertData(&model.DataDefinition{
			Key: request.request.FormValue("key"),

			// get file extension and remove dot before ext name
			Ext: filepath.Ext(fhandler.Filename)[1:],

			Size: uint32(len(buf.Bytes())),
			Buf:  buf.Bytes(),
		})

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
