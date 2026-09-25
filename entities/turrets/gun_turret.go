package turrets

import (
	"github.com/hajimehoshi/ebiten/v2"

	"gocnc/entities"
)

const (
	GunTurretHP     = 200
	GunTurretRange  = 6.0
	GunTurretDamage = 20
	GunTurretAmmo   = 10
	GunTurretRate   = 0.3
)

type GunTurret struct {
	*entities.BuildingBase
}

func NewGunTurret(team int, x, y float64, sprite *ebiten.Image) *GunTurret {
	be := entities.NewBaseEntity("gun_turret", "turret", x, y, 1, 1, GunTurretHP, team, 1)
	be.Spr = sprite
	be.SetSpriteDimensions(0)
	return &GunTurret{
		BuildingBase: &entities.BuildingBase{
			BaseEntity: be,
			GridW:      1,
			GridH:      1,
			Built:      true,
			RangeVal:   GunTurretRange,
			DmgVal:     GunTurretDamage,
			Ammo:       GunTurretAmmo,
			MaxAmmo:    GunTurretAmmo,
			FireRate:   GunTurretRate,
		},
	}
}
