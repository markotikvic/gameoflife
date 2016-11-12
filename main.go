package main

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"math/rand"
	"time"
)

var palette = []color.Color{color.White, color.Black}

const (
	imageSize = 500  // image canvas covers [-imageSize...+imageSize]
	gridSize  = 100
	cellSize  = 5
	delay     = 25   // delay between frames in 10ms units
	nCycles   = 1000 // number of gif frames
)

type Grid struct {
	image  *image.Paletted
	cells  [gridSize][gridSize]Cell
}

type Cell struct {
	x0, y0, x1, y1 int
	alive bool
}

func main() {
	imageName := "game_of_life.gif"
	if len(os.Args) > 1 {
		imageName = os.Args[1]
	}
	gameOfLife(imageName)
}

func blankGrid() *Grid {
	grid := new(Grid)
	rect := image.Rect(0, 0, imageSize, imageSize)
	grid.image = image.NewPaletted(rect, palette)
	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			grid.cells[x][y] = Cell{
				x*cellSize,
				y*cellSize,
				(x+1)*cellSize,
				(y+1)*cellSize,
				false,
			}
		}
	}
	return grid
}

func randomizeGrid(grid *Grid) {
	rand.Seed(time.Now().UnixNano())
	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			if rand.Float32() > 0.5 {
				resurectCell(grid, x, y)
			} else {
				killOffCell(grid, x, y)
			}
		}
	}
}

func applyPatern(grid *Grid, patern string) {
	switch patern {
	case "":
	}
}

func resurectCell(grid *Grid, x, y int) {
	x0 := grid.cells[x][y].x0
	x1 := grid.cells[x][y].x1
	y0 := grid.cells[x][y].y0
	y1 := grid.cells[x][y].y1
	for i := x0; i < x1; i++ {
		for j := y0; j < y1; j++ {
			grid.image.Set(i, j, color.Black)
		}
	}
	grid.cells[x][y].alive = true
}

func killOffCell(grid *Grid, x, y int) {
	x0 := grid.cells[x][y].x0
	x1 := grid.cells[x][y].x1
	y0 := grid.cells[x][y].y0
	y1 := grid.cells[x][y].y1
	for i := x0; i < x1; i++ {
		for j := y0; j < y1; j++ {
			grid.image.Set(i, j, color.White)
		}
	}
	grid.cells[x][y].alive = false
}

func isAlive(grid *Grid, x, y int) bool {
	return grid.cells[x][y].alive
}

func liveNeighboursCount(grid *Grid, x, y int) int {
	total := 0
	var neighbours = []image.Point{
		{-1, -1},
		{ 0, -1},
		{ 1, -1},

		{-1,  1},
		{ 0,  1},
		{ 1,  1},

		{-1,  0},
		{ 1,  0},
	}

	for _, n := range neighbours {
		if isAlive(grid, x + n.X, y + n.Y) {
			total++
		}
	}
	return total
}

func gameTick(old, new *Grid) {
	for x := 1; x < gridSize - 1; x++ {
		for y := 1; y < gridSize - 1; y++ {
			liveCellsCount := liveNeighboursCount(old, x, y)
			if isAlive(old, x, y) {
				// over population and underpopulation
				if liveCellsCount > 3 || liveCellsCount < 2 {
					killOffCell(new, x, y)
				} else {
					// leave cell alive
					resurectCell(new, x, y)
				}
			} else {
				// just right
				if liveCellsCount == 3 {
					resurectCell(new, x, y)
				}
			}
		}
	} 
}

func gameOfLife(imageName string) {
	oldGrid := blankGrid()
	randomizeGrid(oldGrid)
	anim := gif.GIF{LoopCount: nCycles}
	for i := 0; i < nCycles; i++ {
		newGrid := blankGrid()
		gameTick(oldGrid, newGrid)
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, newGrid.image)
		oldGrid = newGrid
	}
	f, _ := os.OpenFile(imageName, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, &anim)
}
