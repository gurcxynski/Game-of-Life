package main

import (
	"image/color"
	"log"

	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	table      Table
	clock      float64
	lastupdate int64
	round      int
}

type Table struct {
	cells [][]bool
}

type Pair struct{
	x int
	y int
	}

func (t *Table) Fill() {
	for i := 0; i < 64; i += 1 {
		t.cells = append(t.cells, []bool{})
		for j := 0; j < 64; j += 1 {
			t.cells[i] = append(t.cells[i], false)
		}
	}
}

func (table *Table) Neighbors(x, y int) int {
	sum := 0
	var neighbors []Pair {
		{x-1, y-1}, {x, y-1}, {x+1, y-1},
		{x-1, y}, 			, {x+1, y},
		{x-1, y+1}, {x, y+1}, {x+1, y+1}
	}
	var valid []bool
	for pair := range neighbors{
		if pair.x >= 0 && pair.x <= 63 && pair.y >= 0 && pair.y <= 63{
			valid = append(valid, table.cells[x][y])
		} 
	}
	for _, v := valid{
		if v{
			sum += 1
		}
	}
	return sum
}

func (g *Game) Start() {
	g.clock = 0.1
	g.table.Fill()
	g.lastupdate = time.Now().UnixMilli()
	g.table.cells[5][5] = true
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
	for i := 0; i < 64; i += 1 {
		for j := 0; j < 64; j += 1 {
			if g.table.cells[i][j] {
				ebitenutil.DrawRect(screen, float64(i*10)+1, float64(j*10)+1, 8, 8, color.White)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 640
}

func main() {
	ebiten.SetWindowSize(640, 640)
	ebiten.SetWindowTitle("Game of Life")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
