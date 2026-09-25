package engine

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"gocnc/interfaces"
)

type InputManager struct {
	Camera     interfaces.Camera
	World      interface {
		EntityAt(tx, ty float64) interfaces.Entity
		IsPassable(tx, ty float64) bool
		Entities() []interfaces.Entity
	}
	Pathfinder interface {
		FindPath(fromX, fromY, toX, toY, speed float64) []interfaces.TilePoint
	}
	Selected   []interfaces.Entity
	PlayerTeam int

	DragStartX, DragStartY int
	dragging               bool
	isDragSelect           bool
	shiftSelect            bool

	CursorHover interfaces.Entity
	hoverX, hoverY float64
}

func NewInputManager(cam interfaces.Camera, playerTeam int) *InputManager {
	return &InputManager{
		Camera:     cam,
		PlayerTeam: playerTeam,
	}
}

func (im *InputManager) Update() {
	mx, my := ebiten.CursorPosition()
	wx, wy := im.Camera.ScreenToWorld(mx, my)
	im.hoverX, im.hoverY = wx, wy
	im.CursorHover = im.World.EntityAt(wx, wy)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		im.DragStartX, im.DragStartY = mx, my
		im.dragging = false
		im.isDragSelect = false
		im.shiftSelect = ebiten.IsKeyPressed(ebiten.KeyShift)
		im.handleLeftClick(mx, my)
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		dx := mx - im.DragStartX
		dy := my - im.DragStartY
		if dx*dx+dy*dy > 100 {
			if !im.dragging {
				im.dragging = true
			}
			im.isDragSelect = true
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if im.isDragSelect {
			im.finishDrag()
		}
		im.dragging = false
		im.isDragSelect = false
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		im.handleRightClick(wx, wy)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		im.Selected = nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		for _, s := range im.Selected {
			if m, ok := s.(interfaces.Moveable); ok {
				m.Stop()
			}
		}
	}
}

func (im *InputManager) handleLeftClick(mx, my int) {
	wx, wy := im.Camera.ScreenToWorld(mx, my)
	entity := im.World.EntityAt(wx, wy)
	if entity == nil || entity.Team() != im.PlayerTeam {
		return
	}
	if im.shiftSelect {
		already := false
		for _, s := range im.Selected {
			if s == entity {
				already = true
				break
			}
		}
		if !already {
			im.Selected = append(im.Selected, entity)
		}
	} else {
		im.Selected = []interfaces.Entity{entity}
	}
}

func (im *InputManager) handleRightClick(wx, wy float64) {
	if len(im.Selected) == 0 {
		return
	}
	hover := im.CursorHover

	if hover != nil && hover.Team() != im.PlayerTeam && hover.Type() != "object" && hover.IsAlive() {
		for _, s := range im.Selected {
			if atkr, ok := s.(interfaces.Attacker); ok {
				atkr.SetTarget(hover)
				atkr.SetOrder(interfaces.NewAttackOrder(hover))
			}
		}
		return
	}

	ix := int(math.Floor(wx))
	iy := int(math.Floor(wy))

	for _, s := range im.Selected {
		if mover, ok := s.(interfaces.Moveable); ok {
			fromX, fromY := mover.Position()
			path := im.Pathfinder.FindPath(fromX, fromY, wx, wy, mover.Speed())
			if path != nil {
				mover.SetOrder(interfaces.NewMoveOrder(wx, wy))
				mover.SetPath(path)
			} else {
				mover.SetOrder(interfaces.NewMoveOrder(wx, wy))
				mover.SetPath([]interfaces.TilePoint{{X: ix, Y: iy}})
			}
		}
	}
}

func (im *InputManager) finishDrag() {
	if !im.shiftSelect {
		im.Selected = nil
	}
	mx, my := ebiten.CursorPosition()
	x1, y1 := float64(im.DragStartX), float64(im.DragStartY)
	x2, y2 := float64(mx), float64(my)
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	for _, e := range im.World.Entities() {
		if !e.IsAlive() || e.Team() != im.PlayerTeam {
			continue
		}
		ex, ey := e.Center()
		sx, sy := im.Camera.WorldToScreen(ex, ey)
		if sx >= x1 && sx <= x2 && sy >= y1 && sy <= y2 {
			already := false
			for _, s := range im.Selected {
				if s == e {
					already = true
					break
				}
			}
			if !already {
				im.Selected = append(im.Selected, e)
			}
		}
	}
}

func (im *InputManager) IsDragging() bool { return im.isDragSelect }

func (im *InputManager) DragRect() (int, int, int, int) {
	mx, my := ebiten.CursorPosition()
	return im.DragStartX, im.DragStartY, mx, my
}

func (im *InputManager) HoverTile() (float64, float64) {
	return im.hoverX, im.hoverY
}
