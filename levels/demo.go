package levels

import (
	"github.com/hajimehoshi/ebiten/v2"

	entities_buildings "gocnc/entities/buildings"
	entities_objects "gocnc/entities/objects"
	entities_turrets "gocnc/entities/turrets"
	entities_vehicles "gocnc/entities/vehicles"
	"gocnc/interfaces"
	"gocnc/layers"
	"gocnc/world"
)

func SetupDemo(w *world.World, sprites map[string]*ebiten.Image) []interfaces.Entity {
	w.TerrainHitmap.Fill(0, 0, 1, w.TilesH, layers.TerrainBlocked)
	w.TerrainHitmap.Fill(w.TilesW-1, 0, 1, w.TilesH, layers.TerrainBlocked)
	w.TerrainHitmap.Fill(0, 0, w.TilesW, 1, layers.TerrainBlocked)
	w.TerrainHitmap.Fill(0, w.TilesH-1, w.TilesW, 1, layers.TerrainBlocked)

	for _, def := range []struct{ x, y int }{
		{3, 3}, {36, 3}, {3, 36}, {36, 36}, {20, 20},
		{15, 10}, {25, 10}, {15, 30}, {25, 30}, {10, 20},
	} {
		w.TerrainHitmap.Set(def.x, def.y, layers.TerrainBlocked)
	}

	var units []interfaces.Entity

	cy := entities_buildings.NewConstructionYard(int(world.TeamGDI), 5, 5, sprites["construction_yard"])
	w.AddEntity(cy)
	w.PlaceWall(int(cy.X), int(cy.Y), cy.GridW, cy.GridH)

	lt1 := entities_vehicles.NewLightTank(int(world.TeamGDI), 8, 5, sprites["light_tank"])
	w.AddEntity(lt1)

	lt2 := entities_vehicles.NewLightTank(int(world.TeamGDI), 8, 7, sprites["light_tank"])
	w.AddEntity(lt2)

	gt1 := entities_turrets.NewGunTurret(int(world.TeamNOD), 30, 10, sprites["gun_turret"])
	w.AddEntity(gt1)

	gt2 := entities_turrets.NewGunTurret(int(world.TeamNOD), 30, 15, sprites["gun_turret"])
	w.AddEntity(gt2)

	gt3 := entities_turrets.NewGunTurret(int(world.TeamNOD), 30, 20, sprites["gun_turret"])
	w.AddEntity(gt3)

	elt := entities_vehicles.NewLightTank(int(world.TeamNOD), 25, 12, sprites["light_tank"])
	w.AddEntity(elt)
	elt.SetOrder(interfaces.NewPatrolOrder(28, 8, 25, 12))

	elt2 := entities_vehicles.NewLightTank(int(world.TeamNOD), 25, 14, sprites["light_tank"])
	w.AddEntity(elt2)
	elt2.SetOrder(interfaces.NewPatrolOrder(28, 20, 25, 14))

	for _, pos := range []struct{ x, y int }{
		{11, 8}, {12, 8}, {13, 8}, {12, 9},
		{15, 18}, {16, 18}, {17, 18},
		{22, 22}, {23, 22}, {22, 23},
		{8, 25}, {9, 25}, {10, 25}, {9, 26},
		{28, 5}, {29, 5}, {28, 6},
	} {
		w.AddEntity(entities_objects.NewTiberium(float64(pos.x), float64(pos.y), 1000, sprites["tree"]))
	}

	for _, pos := range []struct{ x, y int }{
		{5, 12}, {6, 12}, {5, 13},
		{10, 15}, {10, 16}, {11, 15},
		{20, 8}, {21, 8}, {20, 9},
		{15, 28}, {16, 28},
		{32, 8}, {33, 8},
		{32, 25}, {33, 25}, {32, 26},
	} {
		tree := entities_objects.NewTree(float64(pos.x), float64(pos.y), sprites["tree"])
		w.AddEntity(tree)
		w.PlaceWall(pos.x, pos.y, 1, 1)
	}

	w.SortForRendering()
	return units
}
