package pathfinder

import (
	"container/heap"
	"math"

	"gocnc/interfaces"
	"gocnc/layers"
)

type Pathfinder struct {
	Terrain *layers.Hitmap
	Ground  *layers.Hitmap
}

func New() *Pathfinder {
	return &Pathfinder{}
}

func (pf *Pathfinder) FindPath(fromX, fromY, toX, toY float64, speed float64) []interfaces.TilePoint {
	start := interfaces.TilePoint{X: int(math.Round(fromX)), Y: int(math.Round(fromY))}
	goal := interfaces.TilePoint{X: int(math.Round(toX)), Y: int(math.Round(toY))}

	if start == goal {
		return nil
	}

	if pf.Terrain == nil || pf.Ground == nil {
		return nil
	}

	if !pf.Terrain.TerrainPassable(goal.X, goal.Y) || !pf.Ground.GroundPassable(goal.X, goal.Y) {
		goal = pf.nearestPassable(start, goal)
	}

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &node{tile: start, g: 0, f: pf.heuristic(start, goal)})

	cameFrom := map[interfaces.TilePoint]interfaces.TilePoint{}
	gScore := map[interfaces.TilePoint]float64{start: 0}
	closed := map[interfaces.TilePoint]bool{}

	limit := pf.Terrain.Width * pf.Terrain.Height * 2
	iter := 0

	for pq.Len() > 0 && iter < limit {
		iter++
		current := heap.Pop(pq).(*node)

		if current.tile == goal {
			return pf.reconstructPath(cameFrom, current.tile)
		}

		closed[current.tile] = true

		for _, n := range pf.neighbors(current.tile) {
			if closed[n] {
				continue
			}
			if !pf.isPassable(n) {
				continue
			}

			tentativeG := gScore[current.tile] + pf.moveCost(current.tile, n, speed)
			if existingG, ok := gScore[n]; !ok || tentativeG < existingG {
				gScore[n] = tentativeG
				cameFrom[n] = current.tile
				heap.Push(pq, &node{tile: n, g: tentativeG, f: tentativeG + pf.heuristic(n, goal)})
			}
		}
	}

	return nil
}

func (pf *Pathfinder) neighbors(p interfaces.TilePoint) []interfaces.TilePoint {
	dirs := [][2]int{
		{0, -1}, {0, 1}, {-1, 0}, {1, 0},
		{-1, -1}, {1, -1}, {-1, 1}, {1, 1},
	}
	result := make([]interfaces.TilePoint, 0, 8)
	for _, d := range dirs {
		nx, ny := p.X+d[0], p.Y+d[1]
		if nx >= 0 && nx < pf.Terrain.Width && ny >= 0 && ny < pf.Terrain.Height {
			result = append(result, interfaces.TilePoint{X: nx, Y: ny})
		}
	}
	return result
}

func (pf *Pathfinder) isPassable(p interfaces.TilePoint) bool {
	return pf.Terrain.TerrainPassable(p.X, p.Y) && pf.Ground.GroundPassable(p.X, p.Y)
}

func (pf *Pathfinder) moveCost(a, b interfaces.TilePoint, speed float64) float64 {
	if speed <= 0 {
		return math.Inf(1)
	}
	dist := 1.0
	if a.X != b.X && a.Y != b.Y {
		dist = 1.414
	}
	return dist / speed
}

func (pf *Pathfinder) heuristic(a, b interfaces.TilePoint) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return math.Sqrt(dx*dx+dy*dy) / 10.0
}

func (pf *Pathfinder) nearestPassable(from, to interfaces.TilePoint) interfaces.TilePoint {
	xStep := 1
	if from.X > to.X {
		xStep = -1
	}
	yStep := 1
	if from.Y > to.Y {
		yStep = -1
	}

	for x := to.X; x >= 0 && x < pf.Terrain.Width; x += xStep {
		for y := to.Y; y >= 0 && y < pf.Terrain.Height; y += yStep {
			if pf.isPassable(interfaces.TilePoint{X: x, Y: y}) {
				return interfaces.TilePoint{X: x, Y: y}
			}
		}
		if x == from.X {
			break
		}
	}
	return from
}

func (pf *Pathfinder) reconstructPath(cameFrom map[interfaces.TilePoint]interfaces.TilePoint, current interfaces.TilePoint) []interfaces.TilePoint {
	var path []interfaces.TilePoint
	path = append(path, current)
	for {
		prev, ok := cameFrom[current]
		if !ok {
			break
		}
		path = append([]interfaces.TilePoint{prev}, path...)
		current = prev
	}
	if len(path) > 1 {
		path = path[1:]
	}
	return path
}

type node struct {
	tile  interfaces.TilePoint
	g, f  float64
	index int
}

type priorityQueue []*node

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].f < pq[j].f }
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*node)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[:n-1]
	return item
}
