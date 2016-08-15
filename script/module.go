package script

import (
	"image"

	"github.com/anthonynsimon/bild"
	"github.com/yuin/gopher-lua"
)

type scriptCtx struct {
	exports map[string]lua.LGFunction
	rgba    *image.RGBA
}

func newScriptCtx(rgba *image.RGBA) *scriptCtx {
	c := scriptCtx{
		exports: make(map[string]lua.LGFunction),
		rgba:    rgba,
	}

	c.exports["grayscale"] = c.luaGrayscale
	c.exports["boxBlur"] = c.luaBoxBlur
	c.exports["gaussianBlue"] = c.luaGaussianBlur
	c.exports["edgeDetection"] = c.luaEdgeDetection
	c.exports["emboss"] = c.luaEmboss
	c.exports["invert"] = c.luaInvert
	c.exports["flipH"] = c.luaFlipH
	c.exports["flipV"] = c.luaFlipV

	c.exports["getInputImage"] = c.luaGetImage
	c.exports["setOutputImage"] = c.luaSetImage

	return &c
}

func (sc *scriptCtx) load(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), sc.exports)
	L.Push(mod)
	return 1
}

func (sc *scriptCtx) luaGetImage(L *lua.LState) int {
	ud := L.NewUserData()
	ud.Value = sc.rgba
	L.Push(ud)
	return 1
}

func (sc *scriptCtx) luaSetImage(L *lua.LState) int {
	if L.GetTop() == 1 {
		data := L.Get(1).(*lua.LUserData)
		sc.rgba = data.Value.(*image.RGBA)
	}
	return 1
}

func (sc *scriptCtx) luaGrayscale(L *lua.LState) int {
	if L.GetTop() == 1 {
		data := L.Get(1).(*lua.LUserData)
		data.Value = bild.Grayscale(data.Value.(image.Image))
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaGaussianBlur(L *lua.LState) int {
	if L.GetTop() == 2 {
		data := L.Get(1).(*lua.LUserData)
		radius := float64(L.Get(2).(lua.LNumber))
		data.Value = bild.GaussianBlur(data.Value.(image.Image), radius)
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaBoxBlur(L *lua.LState) int {
	if L.GetTop() == 2 {
		data := L.Get(1).(*lua.LUserData)
		radius := float64(L.Get(2).(lua.LNumber))
		data.Value = bild.BoxBlur(data.Value.(image.Image), radius)
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaEdgeDetection(L *lua.LState) int {
	if L.GetTop() == 2 {
		data := L.Get(1).(*lua.LUserData)
		radius := float64(L.Get(2).(lua.LNumber))
		data.Value = bild.EdgeDetection(data.Value.(image.Image), radius)
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaEmboss(L *lua.LState) int {
	if L.GetTop() == 1 {
		data := L.Get(1).(*lua.LUserData)
		data.Value = bild.Emboss(data.Value.(image.Image))
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaInvert(L *lua.LState) int {
	if L.GetTop() == 1 {
		data := L.Get(1).(*lua.LUserData)
		data.Value = bild.Invert(data.Value.(image.Image))
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaFlipH(L *lua.LState) int {
	if L.GetTop() == 1 {
		data := L.Get(1).(*lua.LUserData)
		data.Value = bild.FlipH(data.Value.(image.Image))
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaFlipV(L *lua.LState) int {
	if L.GetTop() == 1 {
		data := L.Get(1).(*lua.LUserData)
		data.Value = bild.FlipV(data.Value.(image.Image))
		L.Push(data)
	}
	return 1
}
