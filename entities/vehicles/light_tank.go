package vehicles

import (
	"github.com/hajimehoshi/ebiten/v2"

	"gocnc/entities"
	"gocnc/interfaces"
)

const (
	LightTankHP     = 200
	LightTankSpeed  = 4.0
	LightTankRange  = 5.0
	LightTankDamage = 30
	LightTankAmmo   = 5
	LightTankRate   = 0.4
)

var lightTankSchema = struct {
	Idle [2]int
	Move [2]int
}{
	Idle: [2]int{0, 31},
	Move: [2]int{32, 63},
}

type LightTank struct {
	*entities.UnitBase
}

func NewLightTank(team int, x, y float64, sprite *ebiten.Image) *LightTank {
	be := entities.NewBaseEntity("light_tank", "vehicle", x, y, 1, 1, LightTankHP, team, 1)
	be.Spr = sprite
	be.SetSpriteDimensions(0)
	return &LightTank{
		UnitBase: &entities.UnitBase{
			BaseEntity: be,
			SpeedVal:   LightTankSpeed,
			RangeVal:   LightTankRange,
			DmgVal:     LightTankDamage,
			Ammo:       LightTankAmmo,
			MaxAmmo:    LightTankAmmo,
			FireRate:   LightTankRate,
		},
	}
}

var _ interfaces.Moveable = (*LightTank)(nil)
var _ interfaces.Attacker = (*LightTank)(nil)
