package sdl2

import (
	"image"
	"image/color"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/veandco/go-sdl2/sdl"
)

type mouseEventInfo struct {
	state bool
	time  uint32
}

type Renderer struct {
	window       *sdl.Window
	renderer     *sdl.Renderer
	fullscreen   bool
	cursorPosX   int
	cursorPosY   int
	stateStack   []surfaceState
	stateCurrent surfaceState
	mouseState   []mouseEventInfo
}

func CreateRenderer() (*Renderer, error) {
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS); err != nil {
		return nil, err
	}

	window, err := sdl.CreateWindow("", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 640, 480, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)

	if err != nil {
		return nil, err
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	//renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)

	if err != nil {
		return nil, err
	}

	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.SetIntegerScale(false)
	renderer.SetLogicalSize(800, 600)
	window.SetMinimumSize(800, 600)

	result := &Renderer{
		window:       window,
		renderer:     renderer,
		stateCurrent: surfaceState{effect: d2enum.DrawEffectNone},
		mouseState:   make([]mouseEventInfo, 10),
	}

	sdl.ShowCursor(0)

	return result, nil
}

func (r *Renderer) GetRendererName() string {
	return "SDL2"
}

func (r *Renderer) SetWindowIcon(fileName string) {

}

func (r *Renderer) Run(f func(d2interface.Surface) error, width, height int, title string) error {
	var running = true

	var event sdl.Event

	r.window.SetTitle(title)
	r.window.SetSize(int32(width), int32(height))

	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseMotionEvent:
				r.cursorPosX = int(t.X)
				r.cursorPosY = int(t.Y)
			case *sdl.MouseButtonEvent:
				r.mouseState[t.Button].time = sdl.GetTicks()
				r.mouseState[t.Button].state = t.Type == sdl.MOUSEBUTTONDOWN
			}
		}

		sdlMutex.Lock()
		if err := r.renderer.SetDrawColor(0, 0, 0, 0); err != nil {
			return err
		}

		if err := r.renderer.Clear(); err != nil {
			return err
		}
		sdlMutex.Unlock()

		if err := f(r); err != nil {
			return err
		}

		sdlMutex.Lock()
		r.renderer.Present()
		sdlMutex.Unlock()
	}

	return nil
}

func (r *Renderer) IsDrawingSkipped() bool {
	return false
}

func (r *Renderer) CreateSurface(surface d2interface.Surface) (d2interface.Surface, error) {
	panic("implement me")
}

func (r *Renderer) NewSurface(width, height int, filter d2enum.Filter) (d2interface.Surface, error) {
	return createSurface(r, width, height, surfaceState{effect: d2enum.DrawEffectNone})
}

func (r *Renderer) IsFullScreen() bool {
	return r.fullscreen
}

func (r *Renderer) SetFullScreen(fullScreen bool) {
	if fullScreen == r.fullscreen {
		return
	}

	if fullScreen {
		r.window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
	} else {
		r.window.SetFullscreen(0)
	}
}

func (r *Renderer) SetVSyncEnabled(vsync bool) {

}

func (r *Renderer) GetVSyncEnabled() bool {
	return false
}

func (r *Renderer) GetCursorPos() (int, int) {
	return r.cursorPosX, r.cursorPosY
}

func (r *Renderer) CurrentFPS() float64 {
	return 60.0
}

func (r *Renderer) Clear(color color.Color) error {
	rr, g, b, a := color.RGBA()

	r.renderer.SetDrawColor(uint8(rr), uint8(g), uint8(b), uint8(a))
	r.renderer.FillRect(nil)

	return nil
}

func (r *Renderer) DrawRect(width, height int, color color.Color) {
	rr, g, b, a := color.RGBA()

	r.renderer.SetDrawColor(uint8(rr), uint8(g), uint8(b), uint8(a))

	rect := sdl.Rect{
		X: int32(r.stateCurrent.x),
		Y: int32(r.stateCurrent.y),
		W: int32(width),
		H: int32(height),
	}

	r.renderer.DrawRect(&rect)
}

func (r *Renderer) DrawLine(x, y int, color color.Color) {
	rr, g, b, a := color.RGBA()

	r.renderer.SetDrawColor(uint8(rr), uint8(g), uint8(b), uint8(a))
	r.renderer.DrawLine(int32(r.stateCurrent.x), int32(r.stateCurrent.y), int32(r.stateCurrent.x+x), int32(r.stateCurrent.y+y))
}

func (r *Renderer) DrawTextf(format string, params ...interface{}) {

}

func (r *Renderer) GetSize() (width, height int) {
	w, h := r.renderer.GetLogicalSize()
	return int(w), int(h)
}

func (r *Renderer) GetDepth() int {
	return len(r.stateStack)
}

func (r *Renderer) Pop() {
	count := len(r.stateStack)
	if count == 0 {
		panic("empty stack")
	}

	r.stateCurrent = r.stateStack[count-1]
	r.stateStack = r.stateStack[:count-1]
}

func (r *Renderer) PopN(n int) {
	for i := 0; i < n; i++ {
		r.Pop()
	}
}

func (r *Renderer) PushColor(color color.Color) {
	r.stateStack = append(r.stateStack, r.stateCurrent)
	r.stateCurrent.color = color
}

func (r *Renderer) PushEffect(effect d2enum.DrawEffect) {
	r.stateStack = append(r.stateStack, r.stateCurrent)
	r.stateCurrent.effect = effect
}

func (r *Renderer) PushFilter(filter d2enum.Filter) {
	r.stateStack = append(r.stateStack, r.stateCurrent)
}

func (r *Renderer) PushTranslation(x, y int) {
	r.stateStack = append(r.stateStack, r.stateCurrent)
	r.stateCurrent.x += x
	r.stateCurrent.y += y
}

func (r *Renderer) PushBrightness(brightness float64) {
	r.stateStack = append(r.stateStack, r.stateCurrent)
	r.stateCurrent.brightness = brightness
}

func (r *Renderer) Render(sfc d2interface.Surface) error {
	sdlMutex.Lock()
	defer sdlMutex.Unlock()

	s := sfc.(*surface)

	rect := sdl.Rect{
		X: int32(r.stateCurrent.x),
		Y: int32(r.stateCurrent.y),
		W: int32(s.width),
		H: int32(s.height),
	}

	switch r.stateCurrent.effect {
	case d2enum.DrawEffectPctTransparency25:
		s.texture.SetAlphaMod(192)
	case d2enum.DrawEffectPctTransparency50:
		s.texture.SetAlphaMod(128)
	case d2enum.DrawEffectPctTransparency75:
		s.texture.SetAlphaMod(64)
	case d2enum.DrawEffectModulate:
		s.texture.SetAlphaMod(255)
		s.texture.SetBlendMode(sdl.BLENDMODE_ADD)
	case d2enum.DrawEffectBurn:
	case d2enum.DrawEffectNormal:
	case d2enum.DrawEffectMod2XTrans:
	case d2enum.DrawEffectMod2X:
	case d2enum.DrawEffectNone:
		s.texture.SetAlphaMod(255)
		s.texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	}

	err := r.renderer.Copy(s.texture, nil, &rect)

	if err != nil {
		return err
	}

	return nil
}

func (r *Renderer) RenderSection(sfc d2interface.Surface, bound image.Rectangle) error {
	return nil
	sdlMutex.Lock()
	defer sdlMutex.Unlock()

	s := sfc.(*surface)

	rect := sdl.Rect{
		X: int32(r.stateCurrent.x),
		Y: int32(r.stateCurrent.y),
		W: int32(bound.Dx()),
		H: int32(bound.Dy()),
	}

	srcRect := sdl.Rect{
		X: int32(bound.Min.X),
		Y: int32(bound.Min.Y),
		W: int32(bound.Dx()),
		H: int32(bound.Dy()),
	}

	switch r.stateCurrent.effect {
	case d2enum.DrawEffectPctTransparency25:
		s.texture.SetAlphaMod(192)
	case d2enum.DrawEffectPctTransparency50:
		s.texture.SetAlphaMod(128)
	case d2enum.DrawEffectPctTransparency75:
		s.texture.SetAlphaMod(64)
	case d2enum.DrawEffectModulate:
		s.texture.SetAlphaMod(255)
		s.texture.SetBlendMode(sdl.BLENDMODE_ADD)
	case d2enum.DrawEffectBurn:
	case d2enum.DrawEffectNormal:
	case d2enum.DrawEffectMod2XTrans:
	case d2enum.DrawEffectMod2X:
	case d2enum.DrawEffectNone:
		s.texture.SetAlphaMod(255)
		s.texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	}

	err := r.renderer.Copy(s.texture, &srcRect, &rect)

	if err != nil {
		return err
	}

	return nil
}

func (r *Renderer) ReplacePixels(pixels []byte) error {
	panic("implement me")
}

func (r *Renderer) Screenshot() *image.RGBA {
	panic("implement me")
}
