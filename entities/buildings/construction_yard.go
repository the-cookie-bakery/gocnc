package buildings

import (
	"github.com/hajimehoshi/ebiten/v2"

	"gocnc/entities"
)

const (
	ConstYardHP       = 400
	ConstYardPowerOut = 30
)

type ConstructionYard struct {
	*entities.BuildingBase
}

func NewConstructionYard(team int, x, y float64, sprite *ebiten.Image) *ConstructionYard {
	be := entities.NewBaseEntity("construction_yard", "building", x, y, 1, 1, ConstYardHP, team, 1)
	be.Spr = sprite
	be.W = 3
	be.H = 2
	be.SetSpriteDimensions(81)
	return &ConstructionYard{
		BuildingBase: &entities.BuildingBase{
			BaseEntity: be,
			GridW:      3,
			GridH:      2,
			Built:      true,
			PowerOut:   ConstYardPowerOut,
		},
	}
}
