package sdl2

import (
	"image"
	"image/color"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
)

type surface struct {
	texture      *sdl.Texture
	stateStack   []surfaceState
	pixelData    []byte
	stateCurrent surfaceState
	width        int
	height       int
	renderer     *Renderer
}

func createSurface(r *Renderer, width, height int, currentState ...surfaceState) (*surface, error) {
	sdlMutex.Lock()
	defer sdlMutex.Unlock()

	texture, err := r.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_TARGET, int32(width), int32(height))

	if err != nil {
		return nil, err
	}

	state := surfaceState{effect: d2enum.DrawEffectNone}
	if len(currentState) > 0 {
		state = currentState[0]
	}

	result := &surface{
		width:        width,
		height:       height,
		texture:      texture,
		stateCurrent: state,
		renderer:     r,
	}

	return result, nil
}

func (s *surface) Clear(color color.Color) error {
	panic("implement me")
}

func (s *surface) DrawRect(width, height int, color color.Color) {
	panic("implement me")
}

func (s *surface) DrawLine(x, y int, color color.Color) {
	panic("implement me")
}

func (s *surface) DrawTextf(format string, params ...interface{}) {
	panic("implement me")
}

func (s *surface) GetSize() (width, height int) {
	return s.width, s.height
}

func (s *surface) GetDepth() int {
	return len(s.stateStack)
}

func (s *surface) Pop() {
	count := len(s.stateStack)
	if count == 0 {
		panic("empty stack")
	}

	s.stateCurrent = s.stateStack[count-1]
	s.stateStack = s.stateStack[:count-1]
}

func (s *surface) PopN(n int) {
	for i := 0; i < n; i++ {
		s.Pop()
	}
}

func (s *surface) PushColor(color color.Color) {
	s.stateStack = append(s.stateStack, s.stateCurrent)
	s.stateCurrent.color = color
}

func (s *surface) PushEffect(effect d2enum.DrawEffect) {
	s.stateStack = append(s.stateStack, s.stateCurrent)
	s.stateCurrent.effect = effect
}

func (s *surface) PushFilter(filter d2enum.Filter) {
	panic("implement me")
}

func (s *surface) PushTranslation(x, y int) {
	s.stateStack = append(s.stateStack, s.stateCurrent)
	s.stateCurrent.x += x
	s.stateCurrent.y += y
}

func (s *surface) PushBrightness(brightness float64) {
	s.stateStack = append(s.stateStack, s.stateCurrent)
	s.stateCurrent.brightness = brightness
}

func (s *surface) Render(sfc d2interface.Surface) error {
	sdlMutex.Lock()
	defer sdlMutex.Unlock()

	target := sfc.(*surface)

	rect := sdl.Rect{
		X: int32(s.stateCurrent.x),
		Y: int32(s.stateCurrent.y),
		W: int32(target.width),
		H: int32(target.height),
	}

	switch target.stateCurrent.effect {
	case d2enum.DrawEffectPctTransparency25:
		target.texture.SetAlphaMod(192)
	case d2enum.DrawEffectPctTransparency50:
		target.texture.SetAlphaMod(128)
	case d2enum.DrawEffectPctTransparency75:
		target.texture.SetAlphaMod(64)
	case d2enum.DrawEffectModulate:
		target.texture.SetAlphaMod(255)
		target.texture.SetBlendMode(sdl.BLENDMODE_MOD)
	case d2enum.DrawEffectBurn:
	case d2enum.DrawEffectNormal:
	case d2enum.DrawEffectMod2XTrans:
	case d2enum.DrawEffectMod2X:
	case d2enum.DrawEffectNone:
		target.texture.SetAlphaMod(255)
		target.texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	}

	s.renderer.renderer.SetRenderTarget(s.texture)
	err := s.renderer.renderer.Copy(target.texture, nil, &rect)
	s.renderer.renderer.SetRenderTarget(nil)

	if err != nil {
		return err
	}

	return nil
}

func (s *surface) RenderSection(surface d2interface.Surface, bound image.Rectangle) error {
	panic("implement me")
}

func (s *surface) ReplacePixels(pixels []byte) error {
	sdlMutex.Lock()
	defer sdlMutex.Unlock()

	s.pixelData = pixels
	return s.texture.Update(nil, pixels, s.width*4)
}

func (s *surface) Screenshot() *image.RGBA {
	panic("implement me")
}
