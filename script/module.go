package script

import (
	"image"
	"image/draw"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/transform"
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
	c.exports["rotate"] = c.luaRotate
	c.exports["resize"] = c.luaResize
	c.exports["crop"] = c.luaCrop

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
		src := effect.Grayscale(data.Value.(image.Image))

		bounds := src.Bounds()
		img := image.NewRGBA(bounds)
		draw.Draw(img, bounds, src, bounds.Min, draw.Src)

		data.Value = img
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaGaussianBlur(L *lua.LState) int {
	if L.GetTop() == 2 {
		data := L.Get(1).(*lua.LUserData)
		radius := float64(L.Get(2).(lua.LNumber))
		data.Value = blur.Gaussian(data.Value.(image.Image), radius)
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaBoxBlur(L *lua.LState) int {
	if L.GetTop() == 2 {
		data := L.Get(1).(*lua.LUserData)
		radius := float64(L.Get(2).(lua.LNumber))
		data.Value = blur.Box(data.Value.(image.Image), radius)
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaEdgeDetection(L *lua.LState) int {
	if L.GetTop() == 2 {
		data := L.Get(1).(*lua.LUserData)
		radius := float64(L.Get(2).(lua.LNumber))
		data.Value = effect.EdgeDetection(data.Value.(image.Image), radius)
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaEmboss(L *lua.LState) int {
	if L.GetTop() == 1 {
		data := L.Get(1).(*lua.LUserData)
		data.Value = effect.Emboss(data.Value.(image.Image))
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaInvert(L *lua.LState) int {
	if L.GetTop() == 1 {
		data := L.Get(1).(*lua.LUserData)
		data.Value = effect.Invert(data.Value.(image.Image))
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaFlipH(L *lua.LState) int {
	if L.GetTop() == 1 {
		data := L.Get(1).(*lua.LUserData)
		data.Value = transform.FlipH(data.Value.(image.Image))
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaFlipV(L *lua.LState) int {
	if L.GetTop() == 1 {
		data := L.Get(1).(*lua.LUserData)
		data.Value = transform.FlipV(data.Value.(image.Image))
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaRotate(L *lua.LState) int {
	if L.GetTop() == 2 {
		data := L.Get(1).(*lua.LUserData)
		val := float64(L.Get(2).(lua.LNumber))
		data.Value = transform.Rotate(data.Value.(image.Image), val, nil)
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaResize(L *lua.LState) int {
	if L.GetTop() == 3 {
		data := L.Get(1).(*lua.LUserData)
		w := int(L.Get(2).(lua.LNumber))
		h := int(L.Get(3).(lua.LNumber))
		data.Value = transform.Resize(data.Value.(image.Image), w, h, transform.Linear)
		L.Push(data)
	}
	return 1
}

func (sc *scriptCtx) luaCrop(L *lua.LState) int {
	if L.GetTop() == 5 {
		data := L.Get(1).(*lua.LUserData)
		x1 := int(L.Get(2).(lua.LNumber))
		y1 := int(L.Get(3).(lua.LNumber))
		x2 := int(L.Get(4).(lua.LNumber))
		y2 := int(L.Get(5).(lua.LNumber))
		data.Value = transform.Crop(data.Value.(image.Image), image.Rect(x1, y1, x2, y2))
		L.Push(data)
	}
	return 1
}
