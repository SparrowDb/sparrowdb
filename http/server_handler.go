package http

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"fmt"

	"github.com/SparrowDb/sparrowdb/auth"
	"github.com/SparrowDb/sparrowdb/db"
	"github.com/SparrowDb/sparrowdb/errors"
	"github.com/SparrowDb/sparrowdb/model"
	"github.com/SparrowDb/sparrowdb/monitor"
	"github.com/SparrowDb/sparrowdb/script"
	"github.com/SparrowDb/sparrowdb/slog"
	"github.com/SparrowDb/sparrowdb/spql"
	"github.com/SparrowDb/sparrowdb/util/uuid"
	"github.com/gin-gonic/gin"
)

// ServeHandler holds main http methods
type ServeHandler struct {
	dbManager     *db.DBManager
	queryExecutor *spql.QueryExecutor
}

func (sh *ServeHandler) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (sh *ServeHandler) userLogin(c *gin.Context) {
	var user auth.User
	c.BindJSON(&user)

	tk, _ := auth.Authenticate(user, sh.dbManager.Config.UserExpire)
	c.JSON(200, gin.H{
		"token": tk,
	})
}

func (sh *ServeHandler) serveQuery(c *gin.Context) {
	var qr spql.QueryRequest
	results := &spql.QueryResult{}

	if err := c.BindJSON(&qr); err != nil {
		results.AddErrorStr(err.Error())
		c.JSON(http.StatusBadRequest, results)
		return
	}

	q, err := qr.ParseQuery()
	if err != nil {
		results.AddErrorStr(err.Error())
		c.JSON(http.StatusBadRequest, results)
		return
	}

	results = <-sh.queryExecutor.ExecuteQuery(&q)
	if results == nil {
		results := &spql.QueryResult{}
		results.AddErrorStr(errors.ErrEmptyQueryResult.Error())
		c.JSON(http.StatusBadRequest, results)
		return
	}

	monitor.IncHTTPQueries()

	c.JSON(http.StatusOK, results)
}

func (sh *ServeHandler) get(c *gin.Context) {
	results := &spql.QueryResult{}
	dbname := c.Param("dbname")
	key := c.Param("key")

	// Check if database exists
	sto, ok := sh.dbManager.GetDatabase(dbname)
	if !ok {
		results.AddErrorStr(errors.ErrDatabaseNotFound.Error())
		c.JSON(http.StatusBadRequest, results)
		return
	}

	// Async get requested data
	result := <-sh.dbManager.GetData(dbname, key)

	// Check if found requested data or DataDefinition is deleted
	if result == nil || result.Status == model.DataDefinitionRemoved {
		results.AddErrorStr(errors.ErrEmptyQueryResult.Error())
		c.JSON(http.StatusBadRequest, results)
		return
	}

	// Token verification if enabled
	if sto.Descriptor.TokenActive {
		token := c.DefaultQuery("token", "")

		if token == "" {
			results.AddErrorStr(errors.ErrWrongRequest.Error())
			c.JSON(http.StatusBadRequest, results)
			return
		}

		if token != result.Token {
			results.AddErrorStr(errors.ErrWrongToken.Error())
			c.JSON(http.StatusBadRequest, results)
			return
		}
	}

	c.Writer.Header().Add("Content-Type", "image/"+result.Ext)
	c.Writer.Header().Add("Content-Length", strconv.Itoa(int(result.Size)))
	c.Writer.Write(result.Buf)
}

func (sh *ServeHandler) upload(c *gin.Context) {
	file, fhandler, err := c.Request.FormFile("uploadfile")
	if err != nil {
		slog.Errorf(err.Error())
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, file)

	dbname := c.PostForm("dbname")

	upsert := false
	if _upsert := c.DefaultPostForm("upsert", "false"); _upsert == "true" {
		upsert = true
	}

	sto, ok := sh.dbManager.GetDatabase(dbname)

	b := buf.Bytes()
	results := spql.QueryResult{}

	if !ok {
		results.AddErrorStr(errors.ErrDatabaseNotFound.Error())
		c.JSON(http.StatusBadRequest, results)
		return
	}
	results.Database = dbname
	dataKey := c.PostForm("key")

	// checks if user request needs script execution
	if scriptName := c.PostForm("script"); len(strings.TrimSpace(scriptName)) > 0 {
		if b, err = script.Execute(scriptName, b); err != nil {
			results.AddErrorStr(fmt.Sprintf(errors.ErrScriptNotExists.Error(), scriptName))
			c.JSON(http.StatusBadRequest, results)
			return
		}
	}

	if isValidKey := spql.ValidateDatabaseName.MatchString(dataKey); !isValidKey {
		results.AddErrorStr(errors.ErrImageInvalidKey.Error())
		c.JSON(http.StatusBadRequest, results)
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

		Buf: b,
	}

	// try to insert image in database
	if _, err := sto.InsertCheckRevision(df, upsert); err != nil {
		results.AddErrorStr(err.Error())
		c.JSON(http.StatusConflict, results)
		return
	}

	// write ok response
	results.AddValue(df.QueryResult())
	c.JSON(http.StatusOK, results)

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
