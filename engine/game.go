package engine

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"gocnc"
	"gocnc/entities"
	entities_buildings "gocnc/entities/buildings"
	entities_turrets "gocnc/entities/turrets"
	"gocnc/interfaces"
	"gocnc/levels"
	"gocnc/pathfinder"
	"gocnc/renderer"
	"gocnc/world"
)

type Game struct {
	Camera     *Camera
	World      *world.World
	Renderer   *renderer.Renderer
	Input      *InputManager
	Pathfinder *pathfinder.Pathfinder
	Audio      *AudioManager

	Sprites map[string]*ebiten.Image

	screenW int
	screenH int
}

func NewGame() *Game {
	tilesW, tilesH := 40, 40

	cam := NewCamera(entities.TileSize)
	cam.SetViewport(960, 540)
	cam.SetWorldSize(tilesW, tilesH)

	w := world.NewWorld(tilesW, tilesH)
	w.PlayerTeam = world.TeamGDI

	gdi := world.NewPlayer(1, world.TeamGDI)
	nod := world.NewPlayer(2, world.TeamNOD)
	w.AddPlayer(gdi)
	w.AddPlayer(nod)

	pf := pathfinder.New()
	pf.Terrain = w.TerrainHitmap
	pf.Ground = w.GroundHitmap

	input := NewInputManager(cam, int(w.PlayerTeam))
	input.World = w
	input.Pathfinder = pf

	rend := renderer.NewRenderer(cam)

	return &Game{
		Camera:     cam,
		World:      w,
		Renderer:   rend,
		Input:      input,
		Pathfinder: pf,
		Sprites:    make(map[string]*ebiten.Image),
	}
}

func (g *Game) LoadSprites() error {
	paths := map[string]string{
		"construction_yard": "entities/buildings/ConstructionYard.png",
		"light_tank":        "entities/vehicles/LightTank.png",
		"gun_turret":        "entities/turrets/GunTurret.png",
		"tree":              "entities/objects/Tree.png",
	}
	var err error
	g.Sprites, err = renderer.LoadSprites(gocnc.EmbeddedFS, paths)
	return err
}

func (g *Game) LoadAudio() error {
	g.Audio = NewAudioManager(gocnc.EmbeddedFS, 44100)
	paths := map[string]string{
		"move":   "actions/UnitConfirmMove1.ogg",
		"attack": "actions/AttackWithTurret.ogg",
		"select": "actions/UnitReporting.ogg",
	}
	for name, path := range paths {
		if err := g.Audio.Load(name, path, gocnc.EmbeddedFS); err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) SetupDemoLevel() {
	levels.SetupDemo(g.World, g.Sprites)
}

func (g *Game) autoTargetTurrets() {
	for _, e := range g.World.Entities() {
		if e.Type() != "turret" || !e.IsAlive() {
			continue
		}
		cx, cy := e.Center()
		rangeVal := 6.0
		var closest interfaces.Entity
		closestDist := 999999.0

		for _, other := range g.World.Entities() {
			if !other.IsAlive() || other.Team() == e.Team() || other.Type() == "object" {
				continue
			}
			ox, oy := other.Position()
			dx := ox - cx
			dy := oy - cy
			dist := dx*dx + dy*dy
			if dist < rangeVal*rangeVal && dist < closestDist {
				closestDist = dist
				closest = other
			}
		}
		if closest != nil {
			if gt, ok := e.(*entities_turrets.GunTurret); ok {
				gt.SetTarget(closest)
			}
		}
	}
}

func (g *Game) processOrders() {
	for _, e := range g.World.Entities() {
		mover, ok := e.(interfaces.Moveable)
		if !ok || mover.IsMoving() {
			continue
		}
		order := mover.CurrentOrder()
		if order == nil {
			continue
		}
		px, py := mover.Position()
		switch order.Type {
		case interfaces.OrderPatrol:
			newOrder := interfaces.NewPatrolOrder(order.X, order.Y, px, py)
			path := g.Pathfinder.FindPath(px, py, order.X, order.Y, mover.Speed())
			mover.SetOrder(newOrder)
			if path != nil {
				mover.SetPath(path)
			}
		case interfaces.OrderAttack:
			if order.Target != nil && order.Target.IsAlive() {
				tx, ty := order.Target.Position()
				path := g.Pathfinder.FindPath(px, py, tx, ty, mover.Speed())
				if path != nil {
					mover.SetPath(path)
				}
			}
		case interfaces.OrderGuard:
			if order.Target != nil && order.Target.IsAlive() {
				tx, ty := order.Target.Position()
				dx := tx - px
				dy := ty - py
				if dx*dx+dy*dy > 9 {
					path := g.Pathfinder.FindPath(px, py, tx, ty, mover.Speed())
					if path != nil {
						mover.SetPath(path)
					}
				}
			}
		}
	}
}

func (g *Game) Update() error {
	g.Input.Update()

	dt := 1.0 / 60.0
	g.Camera.Update(dt)

	g.World.Update(dt)

	g.autoTargetTurrets()
	g.processOrders()

	player := g.World.Player()
	if player != nil {
		player.PowerIn = 0
		player.PowerOut = 0
		for _, e := range g.World.Entities() {
			if e.Type() == "building" && e.Team() == int(g.World.PlayerTeam) && e.IsAlive() {
				if b, ok := e.(*entities_buildings.ConstructionYard); ok && b.PowerOut > 0 {
					player.PowerIn += b.PowerOut
				}
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Camera.SetViewport(g.screenW, g.screenH)
	g.Renderer.SetScreenSize(g.screenW, g.screenH)

	cx, cy := ebiten.CursorPosition()
	g.Renderer.Draw(screen, g.World.Entities(), g.Input.Selected,
		cx, cy, g.Input.CursorHover,
		func(scr *ebiten.Image) {
			g.drawSidebar(scr)
		})

	if g.Input.IsDragging() {
		x1, y1, x2, y2 := g.Input.DragRect()
		fx1 := float32(min(x1, x2))
		fy1 := float32(min(y1, y2))
		fx2 := float32(max(x1, x2))
		fy2 := float32(max(y1, y2))
		vector.StrokeRect(screen, fx1, fy1, fx2-fx1, fy2-fy1, 1, color.RGBA{255, 255, 255, 160}, false)
	}
}

func (g *Game) drawSidebar(screen *ebiten.Image) {
	barX := float32(g.screenW - 160)
	barW := float32(160)
	barH := float32(g.screenH)

	vector.DrawFilledRect(screen, barX, 0, barW, barH, color.RGBA{20, 20, 20, 230}, false)
	vector.StrokeRect(screen, barX, 0, barW, barH, 1, color.RGBA{80, 80, 80, 200}, false)

	p := g.World.Player()

	labelc := color.RGBA{200, 200, 200, 255}
	green := color.RGBA{100, 255, 100, 255}
	red := color.RGBA{255, 100, 100, 255}
	gray := color.RGBA{150, 150, 150, 200}
	dim := color.RGBA{100, 100, 100, 150}

	g.drawText(screen, fmt.Sprintf("$%d", p.Credits), int(barX)+10, 25, green)

	g.drawText(screen, "MAP:", int(barX)+10, 55, labelc)
	g.drawText(screen, "  40x40 tiles", int(barX)+10, 70, gray)

	g.drawText(screen, "KEYS:", int(barX)+10, 100, labelc)
	g.drawText(screen, "  WASD/Arrows Pan", int(barX)+10, 115, gray)
	g.drawText(screen, "  Wheel Zoom", int(barX)+10, 128, gray)
	g.drawText(screen, "  LeftClick Select", int(barX)+10, 141, gray)
	g.drawText(screen, "  RightClick Move", int(barX)+10, 154, gray)
	g.drawText(screen, "  S Stop  Esc Clear", int(barX)+10, 167, gray)

	g.drawText(screen, "POWER:", int(barX)+10, 200, labelc)
	pwrCol := green
	if p.PowerOut > p.PowerIn {
		pwrCol = red
	}
	g.drawText(screen, fmt.Sprintf("  %d/%d", p.PowerOut, p.PowerIn), int(barX)+10, 215, pwrCol)

	g.drawText(screen, "SELECTED:", int(barX)+10, 250, labelc)
	selCount := len(g.Input.Selected)
	if selCount == 0 {
		g.drawText(screen, "  none", int(barX)+10, 265, gray)
	} else if selCount == 1 {
		s := g.Input.Selected[0]
		g.drawText(screen, "  "+s.Name(), int(barX)+10, 265, green)
		g.drawText(screen, fmt.Sprintf("  HP:%d/%d", s.Health(), s.MaxHealth()), int(barX)+10, 280, gray)
	} else {
		g.drawText(screen, fmt.Sprintf("  %d units", selCount), int(barX)+10, 265, green)
	}

	g.drawText(screen, "GoCNC", int(barX)+10, int(barH)-15, dim)
}

func (g *Game) drawText(screen *ebiten.Image, str string, x, y int, clr color.Color) {
	d := &font.Drawer{
		Dst:  screen,
		Src:  image.NewUniform(clr),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(str)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screenW = outsideWidth
	g.screenH = outsideHeight
	g.Camera.SetViewport(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}
