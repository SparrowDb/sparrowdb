package script

import (
	"image"
	"image/color"
	"image/draw"

	lua "github.com/yuin/gopher-lua"
)

const (
	luaImageRGBATypeName = "sparrowRGBA"
)

// SparrowRGBA editable instance of an image
type SparrowRGBA struct {
	RGBA image.RGBA
}

func newRGBA(L *lua.LState) int {
	if L.GetTop() == 1 {
		udata := L.Get(1).(*lua.LUserData)
		data := udata.Value.(*SparrowImage)
		bounds := data.Img.Bounds()

		img := image.NewRGBA(bounds)
		draw.Draw(img, bounds, data.Img, bounds.Min, draw.Src)
		rgba := &SparrowRGBA{*img}

		ud := L.NewUserData()
		ud.Value = rgba
		L.SetMetatable(ud, L.GetTypeMetatable(luaImageRGBATypeName))
		L.Push(ud)
		return 1
	}
	return 0
}

func checkRGBA(L *lua.LState) *SparrowRGBA {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*SparrowRGBA); ok {
		return v
	}
	L.ArgError(1, "SparrowRGBA expected")
	return nil
}

func imgGetPixel(L *lua.LState) int {
	rgba := checkRGBA(L)
	if L.GetTop() == 3 {
		x := int(L.Get(2).(lua.LNumber))
		y := int(L.Get(3).(lua.LNumber))
		r, g, b, a := rgba.RGBA.At(x, y).RGBA()
		tbl := L.NewTable()
		tbl.RawSetH(lua.LString("red"), lua.LNumber(r))
		tbl.RawSetH(lua.LString("green"), lua.LNumber(g))
		tbl.RawSetH(lua.LString("blue"), lua.LNumber(b))
		tbl.RawSetH(lua.LString("alpha"), lua.LNumber(a))
		L.Push(tbl)
		return 1
	}
	L.Push(lua.LBool(false))
	return 0
}

func imgSetPixel(L *lua.LState) int {
	rgba := checkRGBA(L)
	if L.GetTop() == 7 {
		x := int(L.Get(2).(lua.LNumber))
		y := int(L.Get(3).(lua.LNumber))
		r := uint8(L.Get(4).(lua.LNumber))
		g := uint8(L.Get(5).(lua.LNumber))
		b := uint8(L.Get(6).(lua.LNumber))
		a := uint8(L.Get(7).(lua.LNumber))
		rgba.RGBA.SetRGBA(x, y, color.RGBA{r, g, b, a})
		L.Push(lua.LBool(true))
		return 1
	}
	L.Push(lua.LBool(false))
	return 0
}

func imgBounds(L *lua.LState) int {
	rgba := checkRGBA(L)
	b := rgba.RGBA.Bounds().Size()
	tbl := L.NewTable()
	tbl.RawSetH(lua.LString("width"), lua.LNumber(b.X))
	tbl.RawSetH(lua.LString("height"), lua.LNumber(b.Y))
	L.Push(tbl)
	return 1
}

func registerRGBAType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaImageRGBATypeName)
	L.SetGlobal(luaImageRGBATypeName, mt)
	L.SetField(mt, "new", L.NewFunction(newRGBA))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"getPixel": imgGetPixel,
		"setPixel": imgSetPixel,
		"bounds":   imgBounds,
	}))
}
