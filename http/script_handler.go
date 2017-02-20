package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	_error "errors"

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

	content, err := readScriptFile(pname)
	if err != nil {
		resp.AddError(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.AddContent("script", content)
	c.JSON(http.StatusOK, resp)
}
