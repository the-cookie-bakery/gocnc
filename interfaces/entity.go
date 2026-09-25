package interfaces

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type TilePoint struct {
	X, Y int
}

type Entity interface {
	Update(dt float64)
	Render(screen *ebiten.Image, opts *ebiten.DrawImageOptions)

	Position() (x, y float64)
	SetPosition(x, y float64)
	Center() (float64, float64)
	Bounds() image.Rectangle
	OccupiesTiles() []TilePoint

	Health() int
	MaxHealth() int
	Damage(amount int)
	Heal(amount int)
	IsAlive() bool

	Team() int
	Type() string
	Layer() int
	Name() string
	Sprite() *ebiten.Image
	SpriteFrame() int
	SpriteFrameW() float64
	SpriteFrameH() float64
}

type Moveable interface {
	Entity

	Speed() float64
	SetPath(path []TilePoint)
	Path() []TilePoint
	CurrentOrder() *Order
	SetOrder(o *Order)
	IsMoving() bool
	Stop()
}

type Attacker interface {
	Moveable

	Range() float64
	SetTarget(target Entity)
	Target() Entity
	WeaponDamage() int
	CanAttack() bool
}

type OrderType int

const (
	OrderNone OrderType = iota
	OrderMove
	OrderAttack
	OrderPatrol
	OrderGuard
	OrderHarvest
	OrderDeploy
)

type Order struct {
	Type   OrderType
	Target Entity
	X, Y   float64
}

func NewMoveOrder(x, y float64) *Order {
	return &Order{Type: OrderMove, X: x, Y: y}
}

func NewAttackOrder(target Entity) *Order {
	return &Order{Type: OrderAttack, Target: target}
}

func NewPatrolOrder(fromX, fromY, toX, toY float64) *Order {
	return &Order{Type: OrderPatrol, X: toX, Y: toY}
}

func NewGuardOrder(target Entity) *Order {
	return &Order{Type: OrderGuard, Target: target}
}
