package main

import (
	"image/color"
	"log"
	"math/rand"

	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const cellnum = 256
const cellsize = 3
const speed = 1000

type Game struct {
	table      Table
	clock      float64
	lastupdate int64
	round      int
}

type Table struct {
	cells [][]bool
}

type Pair struct {
	x int
	y int
}

func (t *Table) Fill() {
	for i := 0; i < cellnum; i += 1 {
		t.cells = append(t.cells, []bool{})
		for j := 0; j < cellnum; j += 1 {
			t.cells[i] = append(t.cells[i], rand.Int()%2 == 0)
		}
	}
}

func (table *Table) Neighbors(x, y int) int {
	sum := 0
	neighbors := []Pair{
		{x - 1, y - 1}, {x, y - 1}, {x + 1, y - 1},
		{x - 1, y}, {x + 1, y},
		{x - 1, y + 1}, {x, y + 1}, {x + 1, y + 1},
	}
	for _, pair := range neighbors {
		if pair.x >= 0 && pair.x < cellnum && pair.y >= 0 && pair.y < cellnum && table.cells[pair.x][pair.y] {
			sum += 1
		}
	}
	return sum
}

func (g *Game) Start() {
	g.clock = 1000 / speed
	g.table.Fill()
	g.lastupdate = time.Now().UnixMilli()
}

func (g *Game) Update() error {
	if g.round == 0 {
		g.Start()
		g.round += 1
		return nil
	}
	if elapsed := math.Abs(float64(time.Now().UnixMilli() - g.lastupdate)); elapsed < g.clock {
		return nil
	}
	g.lastupdate = time.Now().UnixMilli()
	newtable := Table{}
	newtable.Fill()
	for x, row := range g.table.cells {
		for y, cell := range row {
			sum := g.table.Neighbors(x, y)
			newtable.cells[x][y] = cell
			if sum < 2 {
				newtable.cells[x][y] = false
			}
			if sum > 3 {
				newtable.cells[x][y] = false
			}
			if sum == 3 && !cell {
				newtable.cells[x][y] = true
			}
		}
	}
	g.table = newtable
	g.round += 1
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < cellnum; i += 1 {
		for j := 0; j < cellnum; j += 1 {
			if g.table.cells[i][j] {
				ebitenutil.DrawRect(screen, float64(i*cellsize)+1, float64(j*cellsize)+1, cellsize-2, cellsize-2, color.White)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return cellnum * cellsize, cellnum * cellsize
}

func main() {
	ebiten.SetWindowSize(cellnum*cellsize, cellnum*cellsize)
	ebiten.SetWindowTitle("Game of Life")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
