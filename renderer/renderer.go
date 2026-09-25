package renderer

import (
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"gocnc/entities"
	"gocnc/interfaces"
)

type Renderer struct {
	Camera interfaces.Camera
	Tick   float64
	screenW int
	screenH int
}

func NewRenderer(cam interfaces.Camera) *Renderer {
	return &Renderer{Camera: cam}
}

func (r *Renderer) SetScreenSize(w, h int) {
	r.screenW = w
	r.screenH = h
}

func (r *Renderer) Draw(screen *ebiten.Image, entities []interfaces.Entity, selected []interfaces.Entity,
	cursorX, cursorY int, hover interfaces.Entity, sidebar func(*ebiten.Image)) {

	r.Tick += 1.0 / 60.0
	screen.Fill(color.RGBA{34, 34, 34, 255})

	vx, vy, vw, vh := r.visibleRegion()
	r.drawGrid(screen, vx, vy, vw, vh)

	var vis []interfaces.Entity
	for _, e := range entities {
		if !e.IsAlive() {
			continue
		}
		ex, ey := e.Position()
		if r.isVisible(ex, ey, vx, vy, vw, vh) {
			vis = append(vis, e)
		}
	}

	sort.SliceStable(vis, func(i, j int) bool {
		a, b := vis[i], vis[j]
		if a.Layer() != b.Layer() {
			return a.Layer() < b.Layer()
		}
		_, ay := a.Position()
		_, by := b.Position()
		return ay < by
	})

	for _, e := range vis {
		r.drawEntity(screen, e)
	}

	for _, s := range selected {
		if s.IsAlive() {
			r.drawSelection(screen, s)
		}
	}

	if sidebar != nil {
		sidebar(screen)
	}

	r.drawCursor(screen, cursorX, cursorY, hover)
}

func (r *Renderer) visibleRegion() (float64, float64, float64, float64) {
	sw, sh := r.screenW, r.screenH
	if sw <= 0 {
		sw = 960
	}
	if sh <= 0 {
		sh = 540
	}
	tlX, tlY := r.Camera.ScreenToWorld(0, 0)
	brX, brY := r.Camera.ScreenToWorld(sw, sh)
	pad := 2.0
	return tlX - pad, tlY - pad, brX - tlX + pad*2, brY - tlY + pad*2
}

func (r *Renderer) isVisible(ex, ey, vx, vy, vw, vh float64) bool {
	return ex+2 >= vx && ex-2 <= vx+vw && ey+2 >= vy && ey-2 <= vy+vh
}

func (r *Renderer) drawGrid(screen *ebiten.Image, vx, vy, vw, vh float64) {
	tile := entities.TileSize
	zoom := r.Camera.Zoom()
	col := color.RGBA{50, 50, 50, 80}

	txStart := int(math.Floor(vx))
	txEnd := int(math.Ceil(vx + vw + 1))
	tyStart := int(math.Floor(vy))
	tyEnd := int(math.Ceil(vy + vh + 1))

	for tx := txStart; tx <= txEnd; tx++ {
		for ty := tyStart; ty <= tyEnd; ty++ {
			sx, sy := r.Camera.WorldToScreen(float64(tx), float64(ty))
			size := float32(tile * zoom)
			vector.StrokeRect(screen,
				float32(sx), float32(sy),
				size, size,
				0.5, col, false,
			)
		}
	}
}

func (r *Renderer) drawEntity(screen *ebiten.Image, e interfaces.Entity) {
	cx, cy := e.Center()
	sx, sy := r.Camera.WorldToScreen(cx, cy)
	opts := &ebiten.DrawImageOptions{}

	fw := e.SpriteFrameW()
	fh := e.SpriteFrameH()
	opts.GeoM.Translate(-fw/2, -fh/2)
	r.Camera.ApplyTransform(opts)
	opts.GeoM.Translate(sx, sy)

	e.Render(screen, opts)
}

func (r *Renderer) drawSelection(screen *ebiten.Image, e interfaces.Entity) {
	cx, cy := e.Center()
	sx, sy := r.Camera.WorldToScreen(cx, cy)
	zoom := r.Camera.Zoom()

	teamCol := color.RGBA{0, 180, 80, 220}
	if e.Team() == 2 {
		teamCol = color.RGBA{255, 50, 0, 220}
	}

	radius := float32(14 * zoom)
	vector.StrokeCircle(screen, float32(sx), float32(sy), radius+1, 1.5, color.RGBA{0, 0, 0, 180}, false)
	vector.StrokeCircle(screen, float32(sx), float32(sy), radius, 1.5, teamCol, false)

	hpPct := float64(e.Health()) / float64(e.MaxHealth())
	barW := float32(24 * zoom)
	barH := float32(3 * zoom)
	barY := float32(sy) + radius + 3

	vector.DrawFilledRect(screen, float32(sx)-barW/2, barY, barW, barH, color.RGBA{0, 0, 0, 180}, false)
	vector.DrawFilledRect(screen, float32(sx)-barW/2+1, barY+1, barW-2, barH-2, color.RGBA{80, 0, 0, 220}, false)
	if hpPct > 0 {
		vector.DrawFilledRect(screen, float32(sx)-barW/2+1, barY+1, (barW-2)*float32(hpPct), barH-2, color.RGBA{0, 200, 0, 220}, false)
	}
}

func (r *Renderer) drawCursor(screen *ebiten.Image, x, y int, hover interfaces.Entity) {
	anim := math.Mod(r.Tick*4, 2*math.Pi)
	alpha := uint8(128 + 127*math.Abs(math.Sin(anim)))

	col := color.RGBA{255, 255, 255, alpha}

	if hover != nil && hover.IsAlive() {
		if hover.Team() == 1 {
			col = color.RGBA{0, 220, 80, alpha}
		} else if hover.Team() == 2 {
			col = color.RGBA{255, 80, 80, alpha}
		}
	}

	size := float32(8.0)
	vector.StrokeLine(screen, float32(x)-size, float32(y), float32(x)+size, float32(y), 1.5, col, false)
	vector.StrokeLine(screen, float32(x), float32(y)-size, float32(x), float32(y)+size, 1.5, col, false)
	vector.StrokeCircle(screen, float32(x), float32(y), size*0.75, 1, col, false)
}
