package entities

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	"gocnc/interfaces"
)

const TileSize = 24.0

type BaseEntity struct {
	X, Y       float64
	W, H       float64
	HP         int
	MaxHP      int
	TeamID     int
	EType      string
	ename      string
	ELayer     int
	Alive      bool
	Spr        *ebiten.Image
	Frame      int
	FrameCount int
	FrameW     float64
	FrameH     float64
}

func NewBaseEntity(name, etype string, x, y, w, h float64, hp int, team int, layer int) *BaseEntity {
	e := &BaseEntity{
		ename:  name,
		EType:  etype,
		X:      x,
		Y:      y,
		W:      w,
		H:      h,
		HP:     hp,
		MaxHP:  hp,
		TeamID: team,
		ELayer: layer,
		Alive:  true,
	}
	return e
}

func (e *BaseEntity) Update(dt float64) {}

func (e *BaseEntity) Render(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if e.Spr == nil {
		return
	}
	fw := e.FrameW
	fh := e.FrameH
	if fw <= 0 || fh <= 0 {
		return
	}
	x := float64(e.Frame) * fw
	sub := e.Spr.SubImage(image.Rect(int(x), 0, int(x)+int(fw), int(fh))).(*ebiten.Image)
	screen.DrawImage(sub, opts)
}

func (e *BaseEntity) Position() (float64, float64) { return e.X, e.Y }
func (e *BaseEntity) SetPosition(x, y float64)     { e.X, e.Y = x, y }
func (e *BaseEntity) Center() (float64, float64)   { return e.X, e.Y }

func (e *BaseEntity) Bounds() image.Rectangle {
	return image.Rect(int(e.X), int(e.Y), int(e.X+e.W)+1, int(e.Y+e.H)+1)
}

func (e *BaseEntity) OccupiesTiles() []interfaces.TilePoint {
	var tiles []interfaces.TilePoint
	for ty := int(e.Y); ty < int(e.Y+e.H); ty++ {
		for tx := int(e.X); tx < int(e.X+e.W); tx++ {
			tiles = append(tiles, interfaces.TilePoint{X: tx, Y: ty})
		}
	}
	return tiles
}

func (e *BaseEntity) Health() int     { return e.HP }
func (e *BaseEntity) MaxHealth() int  { return e.MaxHP }
func (e *BaseEntity) IsAlive() bool   { return e.Alive }
func (e *BaseEntity) Team() int       { return e.TeamID }
func (e *BaseEntity) Type() string           { return e.EType }
func (e *BaseEntity) Name() string           { return e.ename }
func (e *BaseEntity) Layer() int             { return e.ELayer }
func (e *BaseEntity) Sprite() *ebiten.Image  { return e.Spr }
func (e *BaseEntity) SpriteFrame() int     { return e.Frame }
func (e *BaseEntity) SpriteFrameW() float64 { return e.FrameW }
func (e *BaseEntity) SpriteFrameH() float64 { return e.FrameH }

func (e *BaseEntity) Damage(amount int) {
	e.HP -= amount
	if e.HP <= 0 {
		e.HP = 0
		e.Alive = false
	}
}

func (e *BaseEntity) Heal(amount int) {
	e.HP += amount
	if e.HP > e.MaxHP {
		e.HP = e.MaxHP
	}
}

func (e *BaseEntity) SetSpriteDimensions(frameCount int) {
	if e.Spr == nil {
		return
	}
	iw := float64(e.Spr.Bounds().Dx())
	ih := float64(e.Spr.Bounds().Dy())
	e.FrameCount = frameCount
	if frameCount > 0 {
		e.FrameW = iw / float64(frameCount)
	} else {
		e.FrameW = ih
	}
	e.FrameH = ih
}
