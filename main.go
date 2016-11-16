package main

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"math/rand"
	"time"
)

const (
	imageSize = 500
	gridSize  = 100
	cellSize  = 5
	delay     = 25
	nCycles   = 250
)

var palette = []color.Color{color.White, color.Black}

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
				cellLives(grid, x, y)
			} else {
				cellDies(grid, x, y)
			}
		}
	}
}

func cellLives(grid *Grid, x, y int) {
	cell := grid.cells[x][y]
	for i := cell.x0; i < cell.x1; i++ {
		for j := cell.y0; j < cell.y1; j++ {
			grid.image.Set(i, j, color.Black)
		}
	}
	grid.cells[x][y].alive = true
}

func cellDies(grid *Grid, x, y int) {
	cell := grid.cells[x][y]
	for i := cell.x0; i < cell.x1; i++ {
		for j := cell.y0; j < cell.y1; j++ {
			grid.image.Set(i, j, color.White)
		}
	}
	grid.cells[x][y].alive = false
}

func cellAlive(grid *Grid, x, y int) bool {
	return grid.cells[x][y].alive
}

func countNeighbours(grid *Grid, x, y int) int {
	total := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if !(i == 0 && j == 0) {
				if cellAlive(grid, x + i, y + j) {
					total++
				}
			}
		}
	}
	return total
}

func gameTick(old, new *Grid) {
	for x := 1; x < gridSize - 1; x++ {
		for y := 1; y < gridSize - 1; y++ {
			n := countNeighbours(old, x, y)
			if cellAlive(old, x, y) {
				if n > 3 || n < 2 {
					cellDies(new, x, y)
				} else {
					cellLives(new, x, y)
				}
			} else {
				if n == 3 {
					cellLives(new, x, y)
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
