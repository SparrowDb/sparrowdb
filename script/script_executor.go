package script

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/SparrowDb/sparrowdb/errors"
	lua "github.com/yuin/gopher-lua"
)

const (
	luaSparrowModuleName = "sparrowdb"
)

// GetScriptPath returns script absolute path
func GetScriptPath() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", errors.ErrReadDir
	}
	return filepath.Join(pwd, "scripts"), nil
}

// Execute executes script that is in scripts folder
func Execute(script, key string, b []byte) ([]byte, error) {
	// check if image is supported
	if IsSupportedFileType(b) == false {
		return nil, errors.ErrNotSupportedFileType
	}

	// check script file
	sp, err := GetScriptPath()
	if err != nil {
		return nil, err
	}

	scriptpath := filepath.Join(sp, script+".lua")
	if _, err := os.Stat(scriptpath); err != nil {
		if os.IsNotExist(err) == false {
			return nil, errors.ErrScriptNotExists
		}
	}

	// lua interpreter
	L := lua.NewState()
	defer L.Close()

	// image bytes to RGBA
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	// register pixel editor
	registerRGBAType(L)

	// register sparrowdb image effect
	si := &SparrowImage{key, "png", img}
	si.registerType(L)

	if err := L.DoFile(scriptpath); err != nil {
		return nil, err
	}

	nb := new(bytes.Buffer)
	png.Encode(nb, si.Img)
	return nb.Bytes(), nil
}
