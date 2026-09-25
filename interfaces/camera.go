package interfaces

import "github.com/hajimehoshi/ebiten/v2"

type Camera interface {
	ScreenToWorld(sx, sy int) (float64, float64)
	WorldToScreen(wx, wy float64) (float64, float64)
	ApplyTransform(opts *ebiten.DrawImageOptions)
	Zoom() float64
}
