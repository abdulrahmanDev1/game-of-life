package main

import "github.com/gdamore/tcell/v2"

type Grid struct {
	Cells [][]bool
}

type GridApp struct {
	screen     tcell.Screen
	grid       *Grid
	gridWidth  int
	gridHeight int
	paused     bool
	changed    bool
}

func NewGrid(width, height int) *Grid {
	cells := make([][]bool, height)
	for i := range cells {
		cells[i] = make([]bool, width)
	}
	return &Grid{Cells: cells}
}

func (g *Grid) Resize(width, height int) {
	newCells := make([][]bool, height)
	for i := range newCells {
		newCells[i] = make([]bool, width)
	}
	for y := 0; y < min(height, len(g.Cells)); y++ {
		for x := 0; x < min(width, len(g.Cells[y])); x++ {
			newCells[y][x] = g.Cells[y][x]
		}
	}
	g.Cells = newCells
}

func (g *Grid) ToggleCell(x, y int) {
	g.Cells[y][x] = !g.Cells[y][x]
}
