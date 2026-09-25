package world

import (
	"math"
	"sort"

	"gocnc/interfaces"
	"gocnc/layers"
)

const (
	LayerTerrain = 0
	LayerGround  = 1
	LayerSky     = 2
	tileSize     = 24.0
)

type World struct {
	TilesW int
	TilesH int

	TerrainHitmap *layers.Hitmap
	GroundHitmap  *layers.Hitmap

	entities []interfaces.Entity
	players  []*Player

	PlayerTeam Team
}

func NewWorld(tilesW, tilesH int) *World {
	return &World{
		TilesW:        tilesW,
		TilesH:        tilesH,
		TerrainHitmap: layers.NewHitmap(tilesW, tilesH),
		GroundHitmap:  layers.NewHitmap(tilesW, tilesH),
	}
}

func (w *World) AddPlayer(p *Player) {
	w.players = append(w.players, p)
}

func (w *World) Player() *Player {
	for _, p := range w.players {
		if p.Team == w.PlayerTeam {
			return p
		}
	}
	return nil
}

func (w *World) Entities() []interfaces.Entity {
	return w.entities
}

func (w *World) AddEntity(e interfaces.Entity) {
	w.entities = append(w.entities, e)
}

func (w *World) EntityAt(tx, ty float64) interfaces.Entity {
	ix := int(math.Floor(tx))
	iy := int(math.Floor(ty))
	px := tx * tileSize
	py := ty * tileSize

	for i := len(w.entities) - 1; i >= 0; i-- {
		e := w.entities[i]
		if !e.IsAlive() {
			continue
		}
		for _, tp := range e.OccupiesTiles() {
			if tp.X == ix && tp.Y == iy {
				return e
			}
		}

		ex, ey := e.Center()
		fw := e.SpriteFrameW()
		fh := e.SpriteFrameH()
		ecx := ex * tileSize
		ecy := ey * tileSize
		if px >= ecx-fw/2 && px <= ecx+fw/2 && py >= ecy-fh/2 && py <= ecy+fh/2 {
			return e
		}
	}
	return nil
}

func (w *World) IsPassable(tx, ty float64) bool {
	ix := int(math.Floor(tx))
	iy := int(math.Floor(ty))
	return w.TerrainHitmap.TerrainPassable(ix, iy) && w.GroundHitmap.GroundPassable(ix, iy)
}

func (w *World) Update(dt float64) {
	for _, e := range w.entities {
		if e.IsAlive() {
			e.Update(dt)
		}
	}

	w.rebuildHitmaps()
	w.cleanupDead()
	w.recheckPaths()
}

func (w *World) rebuildHitmaps() {
	w.GroundHitmap.Clear()
	for _, e := range w.entities {
		if !e.IsAlive() {
			continue
		}
		val := layers.GroundUnit
		if e.Type() == "building" || e.Type() == "turret" {
			val = layers.GroundBuilding
		}
		for _, tp := range e.OccupiesTiles() {
			w.GroundHitmap.Set(tp.X, tp.Y, val)
		}
	}
}

func (w *World) recheckPaths() {
	for _, e := range w.entities {
		mover, ok := e.(interfaces.Moveable)
		if !ok || !mover.IsMoving() {
			continue
		}
		path := mover.Path()
		if path == nil || len(path) == 0 {
			continue
		}
		blocked := false
		for _, tp := range path {
			if !w.GroundHitmap.GroundPassable(tp.X, tp.Y) {
				blocked = true
				break
			}
		}
		if blocked {
			mover.Stop()
		}
	}
}

func (w *World) cleanupDead() {
	alive := w.entities[:0]
	for _, e := range w.entities {
		if e.IsAlive() {
			alive = append(alive, e)
		}
	}
	w.entities = alive
}

func (w *World) SortForRendering() {
	sort.SliceStable(w.entities, func(i, j int) bool {
		a, b := w.entities[i], w.entities[j]
		if a.Layer() != b.Layer() {
			return a.Layer() < b.Layer()
		}
		_, ay := a.Position()
		_, by := b.Position()
		return ay < by
	})
}

func (w *World) PlaceWall(x, y, bw, bh int) {
	w.TerrainHitmap.Fill(x, y, bw, bh, layers.TerrainBlocked)
}
