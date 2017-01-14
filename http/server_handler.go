package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	govalidator "gopkg.in/asaskevich/govalidator.v4"

	"github.com/SparrowDb/sparrowdb/auth"
	"github.com/SparrowDb/sparrowdb/db"
	"github.com/SparrowDb/sparrowdb/errors"
	"github.com/SparrowDb/sparrowdb/model"
	"github.com/SparrowDb/sparrowdb/script"
	"github.com/SparrowDb/sparrowdb/util/uuid"
	"github.com/gin-gonic/gin"
)

// ServeHandler holds main http methods
type ServeHandler struct {
	dbManager *db.DBManager
}

func (sh *ServeHandler) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (sh *ServeHandler) userLogin(c *gin.Context) {
	var user auth.User
	c.BindJSON(&user)

	tk, ok := auth.Authenticate(user, sh.dbManager.Config.UserExpire)
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"token": tk,
		})
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func (sh *ServeHandler) createDatabase(c *gin.Context) {
	resp := NewResponse()
	resp.Database = c.Param("dbname")

	var req model.CreateDatabase
	c.BindJSON(&req)

	databaseCfg := db.DatabaseDescriptor{
		Name:           resp.Database,
		MaxDataLogSize: req.MaxDataLogSize,
		MaxCacheSize:   req.MaxCacheSize,
		BloomFilterFp:  req.BloomFilterFp,
		CronExp:        req.CronExp,
		Path:           req.Path,
		SnapshotPath:   req.SnapshotPath,
	}

	if _, err := govalidator.ValidateStruct(databaseCfg); err != nil {
		resp.AddError(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := sh.dbManager.CreateDatabase(databaseCfg); err != nil {
		resp.AddError(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.AddContent(databaseCfg.Name, req)
	c.JSON(http.StatusOK, resp)
}

func (sh *ServeHandler) dropDatabase(c *gin.Context) {
	resp := NewResponse()
	resp.Database = c.Param("dbname")

	if govalidator.IsAlphanumeric(resp.Database) && govalidator.IsByteLength(resp.Database, 3, 50) {
		if err := sh.dbManager.DropDatabase(resp.Database); err != nil {
			resp.AddError(err)
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		resp.AddContent(resp.Database, "ok")
		c.JSON(http.StatusOK, resp)
	} else {
		resp.AddError(errors.ErrInvalidName)
		c.JSON(http.StatusBadRequest, resp)
		return
	}
}

func (sh *ServeHandler) getDatabaseInfo(resp *Response) int {
	if r := (govalidator.IsAlphanumeric(resp.Database) && govalidator.IsByteLength(resp.Database, 3, 50)); r == false {
		resp.AddError(errors.ErrInvalidName)
		return http.StatusBadRequest
	}

	if db, ok := sh.dbManager.GetDatabase(resp.Database); ok == true {
		resp.AddContent("config", map[string]interface{}{
			"max_datalog_size":           db.Descriptor.MaxDataLogSize,
			"max_cache_size":             db.Descriptor.MaxCacheSize,
			"bloomfilter_fpp":            db.Descriptor.BloomFilterFp,
			"dataholder_cron_compaction": db.Descriptor.CronExp,
			"path":           db.Descriptor.Path,
			"snapshot_path":  db.Descriptor.SnapshotPath,
			"generate_token": db.Descriptor.TokenActive,
			"read_only":      db.Descriptor.ReadOnly,
		})
		resp.AddContent("statistics", db.Info())
		return http.StatusOK
	}

	resp.AddError(errors.ErrDatabaseNotFound)
	return http.StatusBadRequest
}

func (sh *ServeHandler) getDatabaseList(resp *Response) int {
	resp.AddContent("_all", sh.dbManager.GetDatabasesNames())
	return http.StatusBadRequest
}

func (sh *ServeHandler) infoDatabase(c *gin.Context) {
	resp := NewResponse()
	dbname := c.Param("dbname")
	resp.Database = dbname

	if dbname == "_all" {
		sh.getDatabaseList(resp)
	} else {
		sh.getDatabaseInfo(resp)
	}

	c.IndentedJSON(200, resp)
}

func (sh *ServeHandler) uploadData(c *gin.Context) {
	resp := NewResponse()
	resp.Database = c.Param("dbname")

	if r := (govalidator.IsAlphanumeric(resp.Database) && govalidator.IsByteLength(resp.Database, 3, 50)); r == false {
		resp.AddError(errors.ErrInvalidName)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	dataKey := c.Param("key")
	if r := (govalidator.IsAlphanumeric(dataKey) && govalidator.IsByteLength(dataKey, 1, 50)); r == false {
		resp.AddError(errors.ErrImageInvalidKey)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	file, fhandler, err := c.Request.FormFile("uploadfile")
	if err != nil {
		resp.AddError(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, file)

	upsert := false
	if _upsert := c.DefaultPostForm("upsert", "false"); _upsert == "true" {
		upsert = true
	}

	sto, ok := sh.dbManager.GetDatabase(resp.Database)

	b := buf.Bytes()

	if !ok {
		resp.AddError(errors.ErrDatabaseNotFound)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// checks if user request needs script execution
	if scriptName := c.PostForm("script"); len(strings.TrimSpace(scriptName)) > 0 {
		if b, err = script.Execute(scriptName, b); err != nil {
			resp.AddErrorStr(fmt.Sprintf(errors.ErrScriptNotExists.Error(), scriptName))
			c.JSON(http.StatusBadRequest, resp)
			return
		}
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
	if _, err := sto.InsertCheckUpsert(df, upsert); err != nil {
		resp.AddError(err)
		c.JSON(http.StatusConflict, resp)
		return
	}

	// write ok response
	resp.AddContent("data", df.QueryResult())
	c.JSON(http.StatusOK, resp)
}

func (sh *ServeHandler) deleteData(c *gin.Context) {
	resp := NewResponse()
	status := http.StatusBadRequest
	resp.Database = c.Param("dbname")

	if r := (govalidator.IsAlphanumeric(resp.Database) && govalidator.IsByteLength(resp.Database, 3, 50)); r == false {
		resp.AddError(errors.ErrInvalidName)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	dataKey := c.Param("key")
	if r := (govalidator.IsAlphanumeric(dataKey) && govalidator.IsByteLength(dataKey, 1, 50)); r == false {
		resp.AddError(errors.ErrImageInvalidKey)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if db, ok := sh.dbManager.GetDatabase(resp.Database); ok == true {
		storedDf, found := db.GetDataByKey(dataKey)

		// Check if data is in index
		if found {
			// check if data is already marked as tombstone
			if storedDf.Status == model.DataDefinitionRemoved {
				resp.AddErrorStr(fmt.Sprintf("Key %s not found in %s", dataKey, resp.Database))
			} else {
				tbs := model.NewTombstone(storedDf)
				db.InsertCheckUpsert(tbs, true)
				resp.AddContent(resp.Database, "ok")
				status = http.StatusOK
			}
		} else {
			resp.AddErrorStr(fmt.Sprintf("Key %s not found in %s", dataKey, resp.Database))
		}
	} else {
		resp.AddError(errors.ErrDatabaseNotFound)
	}

	// write ok response
	c.JSON(status, resp)
}

func (sh *ServeHandler) getData(dbname, key, token string) (*model.DataDefinition, error) {
	// Check if database exists
	sto, ok := sh.dbManager.GetDatabase(dbname)
	if !ok {
		return nil, errors.ErrDatabaseNotFound
	}

	// Async get requested data
	result := <-sh.dbManager.GetData(dbname, key)

	// Check if found requested data or DataDefinition is deleted
	if result == nil || result.Status == model.DataDefinitionRemoved {
		return nil, errors.ErrEmptyQueryResult
	}

	// Token verification if enabled
	if sto.Descriptor.TokenActive {
		if token == "" {
			return nil, errors.ErrWrongRequest
		}

		if token != result.Token {
			return nil, errors.ErrWrongToken
		}
	}

	return result, nil
}

func (sh *ServeHandler) get(c *gin.Context) {
	resp := NewResponse()
	resp.Database = c.Param("dbname")
	key := c.Param("key")
	token := c.Param("token")

	df, err := sh.getData(resp.Database, key, token)
	if err != nil {
		resp.AddError(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.Writer.Header().Add("Content-Type", "image/"+df.Ext)
	c.Writer.Header().Add("Content-Length", strconv.Itoa(int(df.Size)))
	c.Writer.Write(df.Buf)
}

func (sh *ServeHandler) getDataInfo(c *gin.Context) {
	resp := NewResponse()
	resp.Database = c.Param("dbname")
	key := c.Param("key")
	token := c.Param("token")

	if key == "_keys" {

		sto, ok := sh.dbManager.GetDatabase(resp.Database)

		if !ok {
			resp.AddError(errors.ErrDatabaseNotFound)
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		c.JSON(http.StatusOK, sto.Keys())

		return
	}

	df, err := sh.getData(resp.Database, key, token)
	if err != nil {
		resp.AddError(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.AddContent("data", df.QueryResult())
	c.IndentedJSON(http.StatusOK, resp)
}

// NewServeHandler returns new ServeHandler
func NewServeHandler(dbm *db.DBManager) *ServeHandler {
	return &ServeHandler{
		dbManager: dbm,
	}
}
