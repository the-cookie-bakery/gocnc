package engine

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	X, Y     float64
	zoom     float64
	Width    int
	Height   int
	TileSize float64

	panSpeed    float64
	zoomSpeed   float64
	minZoom     float64
	maxZoom     float64
	edgePan     bool
	edgeMargin  int
	worldWidth  int
	worldHeight int

	midPrevX int
	midPrevY int
	midDrag  bool
}

func NewCamera(tileSize float64) *Camera {
	return &Camera{
		TileSize:   tileSize,
		zoom:       1.0,
		panSpeed:   300,
		zoomSpeed:  0.1,
		minZoom:    0.5,
		maxZoom:    3.0,
		edgePan:    true,
		edgeMargin: 30,
		midPrevX:   -1,
		midPrevY:   -1,
	}
}

func (c *Camera) SetViewport(w, h int) {
	c.Width = w
	c.Height = h
}

func (c *Camera) SetWorldSize(tilesW, tilesH int) {
	c.worldWidth = tilesW
	c.worldHeight = tilesH
}

func (c *Camera) Update(dt float64) {
	speed := c.panSpeed * dt / c.zoom

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		c.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		c.Y += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		c.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		c.X += speed
	}

	x, y := ebiten.Wheel()
	if x != 0 {
		c.zoom += x * c.zoomSpeed
	}
	if y != 0 {
		c.zoom += y * c.zoomSpeed
	}
	c.zoom = math.Max(c.minZoom, math.Min(c.maxZoom, c.zoom))

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		mx, my := ebiten.CursorPosition()
		if c.midPrevX >= 0 {
			dx := float64(c.midPrevX-mx) / c.zoom
			dy := float64(c.midPrevY-my) / c.zoom
			c.X += dx
			c.Y += dy
		}
		c.midPrevX, c.midPrevY = mx, my
		c.midDrag = true
	} else {
		c.midPrevX = -1
		c.midPrevY = -1
		c.midDrag = false
	}

	if c.edgePan {
		mx, my := ebiten.CursorPosition()
		if my >= 0 && my < c.edgeMargin {
			c.Y -= speed
		}
		if my >= c.Height-c.edgeMargin && my < c.Height {
			c.Y += speed
		}
		if mx >= 0 && mx < c.edgeMargin {
			c.X -= speed
		}
		if mx >= c.Width-c.edgeMargin && mx < c.Width {
			c.X += speed
		}
	}

	c.clampViewport()
}

func (c *Camera) clampViewport() {
	worldPxW := float64(c.worldWidth) * c.TileSize
	worldPxH := float64(c.worldHeight) * c.TileSize
	viewW := float64(c.Width) / c.zoom
	viewH := float64(c.Height) / c.zoom

	minX := math.Min(0, worldPxW-viewW)
	maxX := math.Max(0, worldPxW-viewW)
	minY := math.Min(0, worldPxH-viewH)
	maxY := math.Max(0, worldPxH-viewH)

	c.X = math.Max(minX, math.Min(c.X, maxX))
	c.Y = math.Max(minY, math.Min(c.Y, maxY))
}

func (c *Camera) ScreenToWorld(sx, sy int) (float64, float64) {
	wx := float64(sx)/c.zoom + c.X
	wy := float64(sy)/c.zoom + c.Y
	return wx / c.TileSize, wy / c.TileSize
}

func (c *Camera) WorldToScreen(wx, wy float64) (float64, float64) {
	sx := wx*c.TileSize - c.X
	sy := wy*c.TileSize - c.Y
	return sx, sy
}

func (c *Camera) ApplyTransform(opts *ebiten.DrawImageOptions) {
	opts.GeoM.Scale(c.zoom, c.zoom)
}

func (c *Camera) Zoom() float64 {
	return c.zoom
}
