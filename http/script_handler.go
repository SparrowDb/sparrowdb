package http

import (
	"net/http"
	"path"
	"strings"

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
	}
}
