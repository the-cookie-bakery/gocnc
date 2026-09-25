package entities

import (
	"math"

	"gocnc/interfaces"
)

type UnitBase struct {
	*BaseEntity

	SpeedVal float64
	Rotation float64

	path    []interfaces.TilePoint
	pathIdx int
	moving  bool
	order   *interfaces.Order
	target  interfaces.Entity

	RangeVal float64
	DmgVal   int
	Ammo     int
	MaxAmmo  int
	FireRate float64
	fireTimer float64
	reloading bool
}

func (u *UnitBase) Speed() float64    { return u.SpeedVal }
func (u *UnitBase) Range() float64    { return u.RangeVal }
func (u *UnitBase) WeaponDamage() int { return u.DmgVal }

func (u *UnitBase) SetPath(path []interfaces.TilePoint) {
	u.path = path
	u.pathIdx = 0
	u.moving = len(path) > 0
}

func (u *UnitBase) Path() []interfaces.TilePoint { return u.path }

func (u *UnitBase) CurrentOrder() *interfaces.Order { return u.order }

func (u *UnitBase) SetOrder(o *interfaces.Order) {
	u.order = o
	u.moving = false
	u.path = nil
	u.pathIdx = 0
	u.target = nil
	if o != nil {
		u.target = o.Target
	}
}

func (u *UnitBase) IsMoving() bool { return u.moving || (len(u.path) > 0 && u.pathIdx < len(u.path)) }

func (u *UnitBase) Stop() {
	u.moving = false
	u.path = nil
	u.pathIdx = 0
}

func (u *UnitBase) Target() interfaces.Entity { return u.target }

func (u *UnitBase) SetTarget(target interfaces.Entity) {
	u.target = target
}

func (u *UnitBase) CanAttack() bool {
	return u.RangeVal > 0 && u.DmgVal > 0 && u.MaxAmmo > 0 && !u.reloading
}

func (u *UnitBase) Update(dt float64) {
	u.updateMovement(dt)
	u.updateCombat(dt)
}

func (u *UnitBase) updateMovement(dt float64) {
	if !u.moving || u.path == nil || u.pathIdx >= len(u.path) {
		if u.moving {
			u.moving = false
		}
		return
	}

	target := u.path[u.pathIdx]
	tx := float64(target.X) + 0.5
	ty := float64(target.Y) + 0.5
	dx := tx - u.X
	dy := ty - u.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist < 0.05 {
		u.X = tx
		u.Y = ty
		u.pathIdx++
		if u.pathIdx >= len(u.path) {
			u.moving = false
			u.path = nil
			u.pathIdx = 0
		}
		return
	}

	step := u.SpeedVal * dt
	if step > dist {
		step = dist
	}
	u.X += dx / dist * step
	u.Y += dy / dist * step

	u.Rotation = math.Atan2(dx, -dy) * 180.0 / math.Pi
	if u.Rotation < 0 {
		u.Rotation += 360
	}
}

func (u *UnitBase) updateCombat(dt float64) {
	if u.target == nil || !u.CanAttack() {
		return
	}

	if !u.target.IsAlive() {
		u.target = nil
		return
	}

	tx, ty := u.target.Position()
	dx := tx - u.X
	dy := ty - u.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	u.Rotation = math.Atan2(dx, -dy) * 180.0 / math.Pi
	if u.Rotation < 0 {
		u.Rotation += 360
	}

	if dist > u.RangeVal {
		return
	}

	if u.reloading {
		u.fireTimer += dt
		if u.fireTimer >= u.FireRate*2 {
			u.reloading = false
			u.Ammo = u.MaxAmmo
			u.fireTimer = 0
		}
		return
	}

	u.fireTimer += dt
	if u.fireTimer >= u.FireRate {
		u.fireTimer = 0
		u.Ammo--
		u.target.Damage(u.DmgVal)
		if u.Ammo <= 0 {
			u.reloading = true
			u.fireTimer = 0
		}
	}
}

func (u *UnitBase) MoveTo(x, y float64) {
	u.SetOrder(&interfaces.Order{Type: interfaces.OrderMove, X: x, Y: y})
}

func (u *UnitBase) AttackTarget(target interfaces.Entity) {
	u.SetOrder(&interfaces.Order{Type: interfaces.OrderAttack, Target: target})
	u.target = target
}
