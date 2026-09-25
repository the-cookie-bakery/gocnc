package layers

const (
	TerrainClear   = 0
	TerrainBlocked = 1
	TerrainRoad    = 2
)

const (
	GroundEmpty    = 0
	GroundUnit     = 1
	GroundBuilding = 2
)

type Hitmap struct {
	Grid   [][]int
	Width  int
	Height int
}

func NewHitmap(w, h int) *Hitmap {
	grid := make([][]int, h)
	for y := range grid {
		grid[y] = make([]int, w)
	}
	return &Hitmap{Grid: grid, Width: w, Height: h}
}

func (hm *Hitmap) Get(x, y int) int {
	if x < 0 || x >= hm.Width || y < 0 || y >= hm.Height {
		return TerrainBlocked
	}
	return hm.Grid[y][x]
}

func (hm *Hitmap) Set(x, y int, val int) {
	if x >= 0 && x < hm.Width && y >= 0 && y < hm.Height {
		hm.Grid[y][x] = val
	}
}

func (hm *Hitmap) Fill(x, y, w, h, val int) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			hm.Set(x+dx, y+dy, val)
		}
	}
}

func (hm *Hitmap) Clear() {
	for y := range hm.Grid {
		for x := range hm.Grid[y] {
			hm.Grid[y][x] = 0
		}
	}
}

func (hm *Hitmap) IsPassable(x, y int) bool {
	return hm.Get(x, y) == 0
}

func (hm *Hitmap) TerrainPassable(x, y int) bool {
	v := hm.Get(x, y)
	return v == TerrainClear || v == TerrainRoad
}

func (hm *Hitmap) GroundPassable(x, y int) bool {
	return hm.Get(x, y) == GroundEmpty
}
