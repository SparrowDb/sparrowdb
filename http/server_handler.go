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
	"github.com/SparrowDb/sparrowdb/cluster"
	"github.com/SparrowDb/sparrowdb/db"
	"github.com/SparrowDb/sparrowdb/errors"
	"github.com/SparrowDb/sparrowdb/model"
	"github.com/SparrowDb/sparrowdb/monitor"
	"github.com/SparrowDb/sparrowdb/script"
	"github.com/SparrowDb/sparrowdb/slog"
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

	tk, _ := auth.Authenticate(user, sh.dbManager.Config.UserExpire)
	c.JSON(200, gin.H{
		"token": tk,
	})
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

func (sh *ServeHandler) get(c *gin.Context) {
	resp := NewResponse()
	resp.Database = c.Param("dbname")
	key := c.Param("key")

	// Check if database exists
	sto, ok := sh.dbManager.GetDatabase(resp.Database)
	if !ok {
		resp.AddError(errors.ErrDatabaseNotFound)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Async get requested data
	result := <-sh.dbManager.GetData(resp.Database, key)

	// Check if found requested data or DataDefinition is deleted
	if result == nil || result.Status == model.DataDefinitionRemoved {
		resp.AddError(errors.ErrEmptyQueryResult)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Token verification if enabled
	if sto.Descriptor.TokenActive {
		token := c.DefaultQuery("token", "")

		if token == "" {
			resp.AddError(errors.ErrWrongRequest)
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if token != result.Token {
			resp.AddError(errors.ErrWrongToken)
			c.JSON(http.StatusBadRequest, resp)
			return
		}
	}

	c.Writer.Header().Add("Content-Type", "image/"+result.Ext)
	c.Writer.Header().Add("Content-Length", strconv.Itoa(int(result.Size)))
	c.Writer.Write(result.Buf)
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
		slog.Errorf(err.Error())
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

	if sh.dbManager.Config.EnableCluster {
		cluster.PublishData(*df, resp.Database)
	}

	// increment upload statistics
	monitor.IncHTTPUploads()
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

// NewServeHandler returns new ServeHandler
func NewServeHandler(dbm *db.DBManager) *ServeHandler {
	return &ServeHandler{
		dbManager: dbm,
	}
}
