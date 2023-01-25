package main

import (
	"fmt"
	"image/color"
	"log"

	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	table      Table
	clock      float64
	lastupdate int64
	round      int
	ready      bool
}

var cellnum int
var cellsize int
var speed float64

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
			t.cells[i] = append(t.cells[i], false)
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

func (g *Game) Setup() {
	g.clock = 1000 / speed
	g.table.Fill()
	g.round = 0
	g.ready = true
}

func (g *Game) Update() error {
	if !g.ready {
		g.Setup()
		return nil
	}
	if g.round == 0 {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			x = x / cellsize
			y = y / cellsize
			g.table.cells[x][y] = true
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			x, y := ebiten.CursorPosition()
			x = x / cellsize
			y = y / cellsize
			g.table.cells[x][y] = false
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.round = 1
			g.lastupdate = time.Now().UnixMilli()
		}
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
				ebitenutil.DrawRect(screen, float64(i*cellsize)+1, float64(j*cellsize)+1, float64(cellsize-2), float64(cellsize-2), color.White)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return cellnum * cellsize, cellnum * cellsize
}

func main() {

	var input int
	fmt.Println("Number of cells in a row: ")
	fmt.Scanln(&input)
	cellnum = input
	fmt.Println("Pixel size of a cell: ")
	fmt.Scanln(&input)
	cellsize = input
	var s float64
	fmt.Println("Amount of game ticks per second: ")
	fmt.Scanln(&s)
	speed = s

	ebiten.SetWindowSize(cellnum*cellsize, cellnum*cellsize)
	ebiten.SetWindowTitle("Game of Life")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
