package script

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/yuin/gopher-lua"
)

const (
	luaSparrowModuleName = "sparrowdb"
)

// Execute executes script that is in scripts folder
func Execute(script string, b []byte) ([]byte, error) {
	// check script file
	pwd, _ := os.Getwd()
	scriptpath := filepath.Join(pwd, "scripts", script+".lua")

	if _, err := os.Stat(scriptpath); err != nil {
		if os.IsNotExist(err) == false {
			return nil, err
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

	buf := image.NewRGBA(img.Bounds())
	draw.Draw(buf, buf.Bounds(), img, image.Point{0, 0}, draw.Src)

	// load lua modules
	sc := newScriptCtx(buf)
	L.PreloadModule(luaSparrowModuleName, sc.load)

	if err := L.DoFile(scriptpath); err != nil {
		return nil, err
	}

	nb := new(bytes.Buffer)
	png.Encode(nb, sc.rgba)
	return nb.Bytes(), nil
}
