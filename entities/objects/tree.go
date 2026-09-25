package objects

import (
	"github.com/hajimehoshi/ebiten/v2"

	"gocnc/entities"
)

type Tree struct {
	*entities.StaticEntity
}

func NewTree(x, y float64, sprite *ebiten.Image) *Tree {
	be := entities.NewBaseEntity("tree", "object", x, y, 1, 1, 1, 0, 0)
	be.Spr = sprite
	be.SetSpriteDimensions(0)
	return &Tree{
		StaticEntity: &entities.StaticEntity{
			BaseEntity:     be,
			BlocksMovement: true,
		},
	}
}
