package sdl2

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	mousePressThreshold = 25
)

func (r *Renderer) CursorPosition() (x int, y int) {
	return r.cursorPosX, r.cursorPosY
}

func (r *Renderer) InputChars() []rune {
	return []rune{}
}

func (r *Renderer) IsKeyPressed(key d2enum.Key) bool {
	return false
}

func (r *Renderer) IsKeyJustPressed(key d2enum.Key) bool {
	return false
}

func (r *Renderer) IsKeyJustReleased(key d2enum.Key) bool {
	return false
}

func (r *Renderer) IsMouseButtonPressed(button d2enum.MouseButton) bool {
	switch button {
	case d2enum.MouseButtonLeft:
		return r.mouseState[1].state
	case d2enum.MouseButtonRight:
		return r.mouseState[2].state
	case d2enum.MouseButtonMiddle:
		return r.mouseState[3].state
	default:
		return false
	}
}

func (r *Renderer) IsMouseButtonJustPressed(button d2enum.MouseButton) bool {
	switch button {
	case d2enum.MouseButtonLeft:
		return r.mouseState[1].state == true && (sdl.GetTicks()-r.mouseState[1].time) < mousePressThreshold
	case d2enum.MouseButtonRight:
		return r.mouseState[2].state == true && (sdl.GetTicks()-r.mouseState[2].time) < mousePressThreshold
	case d2enum.MouseButtonMiddle:
		return r.mouseState[3].state == true && (sdl.GetTicks()-r.mouseState[3].time) < mousePressThreshold
	default:
		return false
	}
}

func (r *Renderer) IsMouseButtonJustReleased(button d2enum.MouseButton) bool {
	switch button {
	case d2enum.MouseButtonLeft:
		return r.mouseState[1].state == false && (sdl.GetTicks()-r.mouseState[1].time) < mousePressThreshold
	case d2enum.MouseButtonRight:
		return r.mouseState[2].state == false && (sdl.GetTicks()-r.mouseState[2].time) < mousePressThreshold
	case d2enum.MouseButtonMiddle:
		return r.mouseState[3].state == false && (sdl.GetTicks()-r.mouseState[3].time) < mousePressThreshold
	default:
		return false
	}
}

func (r *Renderer) KeyPressDuration(key d2enum.Key) int {
	return 0
}
