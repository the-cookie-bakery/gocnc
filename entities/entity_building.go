package entities

import "gocnc/interfaces"

type BuildingBase struct {
	*BaseEntity

	GridW int
	GridH int

	BuildTime  float64
	buildTimer float64
	Built      bool

	PowerIn  int
	PowerOut int

	RangeVal float64
	DmgVal   int
	FireRate float64
	fireTimer float64
	Ammo     int
	MaxAmmo  int
	reloading bool
	Target   interfaces.Entity
}

func (b *BuildingBase) Center() (float64, float64) {
	return b.X + float64(b.GridW)/2, b.Y + float64(b.GridH)/2
}

func (b *BuildingBase) OccupiesTiles() []interfaces.TilePoint {
	var tiles []interfaces.TilePoint
	for ty := 0; ty < b.GridH; ty++ {
		for tx := 0; tx < b.GridW; tx++ {
			tiles = append(tiles, interfaces.TilePoint{X: int(b.X) + tx, Y: int(b.Y) + ty})
		}
	}
	return tiles
}

func (b *BuildingBase) Update(dt float64) {
	if !b.Built {
		b.buildTimer += dt
		if b.buildTimer >= b.BuildTime {
			b.Built = true
		}
		b.Frame = int(b.buildTimer / b.BuildTime * float64(b.Spr.Bounds().Dx()/int(TileSize)))
		return
	}
	b.updateCombat(dt)
}

func (b *BuildingBase) SetTarget(t interfaces.Entity) { b.Target = t }

func (b *BuildingBase) updateCombat(dt float64) {
	if b.Target == nil || b.RangeVal <= 0 || b.DmgVal <= 0 {
		return
	}
	if !b.Target.IsAlive() {
		b.Target = nil
		return
	}

	tx, ty := b.Target.Position()
	dx := tx - b.CenterX()
	dy := ty - b.CenterY()
	rangeSq := b.RangeVal * b.RangeVal

	if dx*dx+dy*dy > rangeSq {
		return
	}

	if b.reloading {
		b.fireTimer += dt
		if b.fireTimer >= b.FireRate*2 {
			b.reloading = false
			b.Ammo = b.MaxAmmo
			b.fireTimer = 0
		}
		return
	}

	b.fireTimer += dt
	if b.fireTimer >= b.FireRate {
		b.fireTimer = 0
		b.Ammo--
		b.Target.Damage(b.DmgVal)
		if b.Ammo <= 0 {
			b.reloading = true
			b.fireTimer = 0
		}
	}
}

func (b *BuildingBase) CenterX() float64 { return b.X + float64(b.GridW)/2 }
func (b *BuildingBase) CenterY() float64 { return b.Y + float64(b.GridH)/2 }
