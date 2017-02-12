package script

import (
	"image"
	"image/draw"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/transform"
	lua "github.com/yuin/gopher-lua"
)

const (
	luaImageTypeName         = "sparrowImage"
	luaImageTypeInstanceName = "imageCtx"
)

// SparrowImage main image effects
type SparrowImage struct {
	Name string
	Type string
	Img  image.Image
}

func (s *SparrowImage) registerType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaImageTypeName)
	L.SetGlobal(luaImageTypeName, mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"name":          s.getName,
		"type":          s.getType,
		"grayscale":     s.imgGrayscale,
		"gaussianBlur":  s.imgGaussianBlur,
		"boxBlur":       s.imgBoxBlur,
		"edgeDetection": s.imgEdgeDetection,
		"emboss":        s.imgEmboss,
		"invert":        s.imgInvert,
		"flipH":         s.imgFlipH,
		"flipV":         s.imgFlipV,
		"rotate":        s.imgRotate,
		"resize":        s.imgResize,
		"crop":          s.imgCrop,
		"bounds":        s.imgBounds,
		"setOutput":     s.imgSetOutput,
	}))

	// put an instance of SparrowImage with filled attrs in context
	ud := L.NewUserData()
	ud.Value = s
	L.SetMetatable(ud, L.GetTypeMetatable(luaImageTypeName))
	L.Push(ud)
	L.SetGlobal(luaImageTypeInstanceName, ud)
}

func (s *SparrowImage) getName(L *lua.LState) int {
	L.Push(lua.LString(s.Name))
	return 1
}

func (s *SparrowImage) getType(L *lua.LState) int {
	L.Push(lua.LString(s.Type))
	return 1
}

func (s *SparrowImage) imgGrayscale(L *lua.LState) int {
	src := effect.Grayscale(s.Img)
	bounds := src.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, src, bounds.Min, draw.Src)
	s.Img = img
	L.Push(lua.LBool(true))
	return 1
}

func (s *SparrowImage) imgGaussianBlur(L *lua.LState) int {
	if L.GetTop() == 2 {
		radius := float64(L.Get(2).(lua.LNumber))
		s.Img = blur.Gaussian(s.Img, radius)
		L.Push(lua.LBool(true))
		return 1
	}
	L.Push(lua.LBool(false))
	return 0
}

func (s *SparrowImage) imgBoxBlur(L *lua.LState) int {
	if L.GetTop() == 2 {
		radius := float64(L.Get(2).(lua.LNumber))
		s.Img = blur.Box(s.Img, radius)
		L.Push(lua.LBool(true))
		return 1
	}
	L.Push(lua.LBool(false))
	return 0
}

func (s *SparrowImage) imgEdgeDetection(L *lua.LState) int {
	if L.GetTop() == 2 {
		radius := float64(L.Get(2).(lua.LNumber))
		s.Img = effect.EdgeDetection(s.Img, radius)
		L.Push(lua.LBool(true))
		return 1
	}
	L.Push(lua.LBool(false))
	return 0
}

func (s *SparrowImage) imgEmboss(L *lua.LState) int {
	if L.GetTop() == 1 {
		s.Img = effect.Emboss(s.Img)
		L.Push(lua.LBool(true))
		return 1
	}
	L.Push(lua.LBool(false))
	return 0
}

func (s *SparrowImage) imgInvert(L *lua.LState) int {
	if L.GetTop() == 1 {
		s.Img = effect.Invert(s.Img)
		L.Push(lua.LBool(true))
		return 1
	}
	L.Push(lua.LBool(false))
	return 0
}

func (s *SparrowImage) imgFlipH(L *lua.LState) int {
	if L.GetTop() == 1 {
		s.Img = transform.FlipH(s.Img)
		L.Push(lua.LBool(true))
		return 1
	}
	L.Push(lua.LBool(false))
	return 0
}

func (s *SparrowImage) imgFlipV(L *lua.LState) int {
	if L.GetTop() == 1 {
		s.Img = transform.FlipV(s.Img)
		L.Push(lua.LBool(true))
		return 1
	}
	L.Push(lua.LBool(false))
	return 0
}

func (s *SparrowImage) imgRotate(L *lua.LState) int {
	if L.GetTop() == 2 {
		val := float64(L.Get(2).(lua.LNumber))
		s.Img = transform.Rotate(s.Img, val, nil)
		L.Push(lua.LBool(true))
		return 1
	}
	L.Push(lua.LBool(false))
	return 0
}

func (s *SparrowImage) imgResize(L *lua.LState) int {
	if L.GetTop() == 3 {
		w := int(L.Get(2).(lua.LNumber))
		h := int(L.Get(3).(lua.LNumber))
		s.Img = transform.Resize(s.Img, w, h, transform.Linear)
		L.Push(lua.LBool(true))
		return 1
	}
	L.Push(lua.LBool(false))
	return 0
}

func (s *SparrowImage) imgCrop(L *lua.LState) int {
	if L.GetTop() == 5 {
		x1 := int(L.Get(2).(lua.LNumber))
		y1 := int(L.Get(3).(lua.LNumber))
		x2 := int(L.Get(4).(lua.LNumber))
		y2 := int(L.Get(5).(lua.LNumber))
		s.Img = transform.Crop(s.Img, image.Rect(x1, y1, x2, y2))
		L.Push(lua.LBool(true))
		return 1
	}
	L.Push(lua.LBool(false))
	return 0
}

func (s *SparrowImage) imgBounds(L *lua.LState) int {
	b := s.Img.Bounds().Size()
	tbl := L.NewTable()
	tbl.RawSetH(lua.LString("width"), lua.LNumber(b.X))
	tbl.RawSetH(lua.LString("height"), lua.LNumber(b.Y))
	L.Push(tbl)
	L.Push(lua.LBool(false))
	return 1
}

func (s *SparrowImage) imgSetOutput(L *lua.LState) int {
	if L.GetTop() == 2 {
		data := L.Get(2).(*lua.LUserData)
		rgba := data.Value.(*SparrowRGBA)
		s.Img = rgba.RGBA.SubImage(rgba.RGBA.Rect)
		L.Push(lua.LBool(true))
		return 1
	}
	L.Push(lua.LBool(false))
	return 0
}
