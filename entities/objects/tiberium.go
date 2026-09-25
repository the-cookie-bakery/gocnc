package objects

import (
	"github.com/hajimehoshi/ebiten/v2"

	"gocnc/entities"
)

type Tiberium struct {
	*entities.StaticEntity
	Amount int
}

func NewTiberium(x, y float64, amount int, sprite *ebiten.Image) *Tiberium {
	be := entities.NewBaseEntity("tiberium", "object", x, y, 1, 1, 1, 0, 0)
	be.Spr = sprite
	be.SetSpriteDimensions(0)
	return &Tiberium{
		StaticEntity: &entities.StaticEntity{
			BaseEntity:     be,
			BlocksMovement: false,
		},
		Amount: amount,
	}
}
