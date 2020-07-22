package sdl2

import (
	"image/color"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
)

type surfaceState struct {
	x          int
	y          int
	color      color.Color
	brightness float64
	effect     d2enum.DrawEffect
}
