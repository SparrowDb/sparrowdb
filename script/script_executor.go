package script

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/yuin/gopher-lua"
)

const (
	luaSparrowModuleName = "sparrowdb"
)

// Execute executes script that is in scripts folder
func Execute(scriptpath string, b []byte) ([]byte, error) {
	// check script file
	if _, err := os.Stat(scriptpath); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}

	// lua interpreter
	L := lua.NewState()
	defer L.Close()

	// image bytes to RGBA
	img, _, _ := image.Decode(bytes.NewReader(b))
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
