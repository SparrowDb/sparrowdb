package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	govalidator "gopkg.in/asaskevich/govalidator.v4"

	_error "errors"

	"os"

	"github.com/SparrowDb/sparrowdb/auth"
	"github.com/SparrowDb/sparrowdb/errors"
	"github.com/SparrowDb/sparrowdb/script"
	"github.com/SparrowDb/sparrowdb/util"
	"github.com/gin-gonic/gin"
)

func scriptList() ([]string, error) {
	scriptPath, err := script.GetScriptPath()
	if err != nil {
		return nil, err
	}

	files, err := util.ListDirectory(scriptPath)
	if err != nil {
		return nil, err
	}

	scriptName := []string{}

	for _, scriptDir := range files {
		sname := path.Base(scriptDir)
		n := strings.LastIndexByte(sname, '.')
		if n > 0 {
			sname = sname[:n]
		}
		scriptName = append(scriptName, sname)
	}

	return scriptName, nil
}

func readScriptFile(name string) (string, error) {
	scriptPath, err := script.GetScriptPath()
	if err != nil {
		return "", err
	}

	fpath := filepath.Join(scriptPath, name+".lua")
	if exists, err := util.Exists(fpath); exists == false {
		if err != nil {
			return "", nil
		}
		return "", _error.New(fmt.Sprintf(errors.ErrScriptNotExists.Error(), name))
	}

	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func getScriptList(c *gin.Context) {
	resp := NewResponse()
	pname := c.Param("name")

	if pname == "_all" {
		scripts, err := scriptList()
		if err != nil {
			resp.AddError(err)
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		resp.AddContent("scripts", scripts)
		c.JSON(http.StatusOK, resp)
		return
	}

	if hasPermission(c, auth.RoleScriptManager) == false {
		resp.AddError(errors.ErrNoPrivilege)
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	content, err := readScriptFile(pname)
	if err != nil {
		resp.AddError(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.AddContent("script", content)
	c.JSON(http.StatusOK, resp)
}

func saveScript(c *gin.Context) {
	resp := NewResponse()

	if hasPermission(c, auth.RoleScriptManager) == false {
		resp.AddError(errors.ErrNoPrivilege)
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	scriptName := c.Param("name")

	if r := (govalidator.IsAlphanumeric(scriptName) && govalidator.IsByteLength(scriptName, 3, 50)); r == false {
		resp.AddError(errors.ErrScriptInvalidName)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var scriptInfo struct {
		Content string `json:"content"`
	}

	if err := c.BindJSON(&scriptInfo); err != nil {
		resp.AddError(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	scriptPath, err := script.GetScriptPath()
	if err != nil {
		resp.AddError(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	fpath := filepath.Join(scriptPath, scriptName+".lua")
	if err := ioutil.WriteFile(fpath, []byte(scriptInfo.Content), 0644); err != nil {
		resp.AddError(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.AddContent("script", scriptName)
	c.JSON(http.StatusOK, resp)
}

func deleteScript(c *gin.Context) {
	resp := NewResponse()
	scriptName := c.Param("name")

	if r := (govalidator.IsAlphanumeric(scriptName) && govalidator.IsByteLength(scriptName, 3, 50)); r == false {
		resp.AddError(errors.ErrScriptInvalidName)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	scriptPath, err := script.GetScriptPath()
	if err != nil {
		resp.AddError(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	fpath := filepath.Join(scriptPath, scriptName+".lua")
	if err := os.Remove(fpath); err != nil {
		resp.AddError(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.AddContent("script", scriptName)
	c.JSON(http.StatusOK, resp)
}
