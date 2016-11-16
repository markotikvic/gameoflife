package main

import (
	"flag"
	"image"
	"image/color"
	"image/gif"
	"math/rand"
	"os"
	"time"
)

const (
	imageSize = 500
	gridSize  = 100
	cellSize  = 5
	delay     = 25
)

var palette = []color.Color{color.White, color.Black}

var nCycles = flag.Uint64("n", 100, "number of life cycles")
var outName = flag.String("o", "life.gif", "output file name")

type Cell struct {
	x0, y0, x1, y1 int
	alive          bool
}

type Cells [gridSize][gridSize]Cell

func main() {
	flag.Parse()
	gameOfLife()
}

func newCells() *Cells {
	cells := new(Cells)
	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			cells[x][y] = Cell{
				x * cellSize,
				y * cellSize,
				(x + 1) * cellSize,
				(y + 1) * cellSize,
				false,
			}
		}
	}
	return cells
}

func randomCells() *Cells {
	g := newCells()
	randomizeCells(g)
	return g
}

func randomizeCells(cells *Cells) {
	rand.Seed(time.Now().UnixNano())
	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			if rand.Float32() > 0.5 {
				cellLives(&cells[x][y])
			} else {
				cellDies(&cells[x][y])
			}
		}
	}
}

func cellLives(c *Cell) {
	c.alive = true
}

func cellDies(c *Cell) {
	c.alive = false
}

func cellAlive(c *Cell) bool {
	return c.alive
}

func cellInGrid(x, y int) bool {
	return ((x > 0 && x < gridSize) && (y > 0 && y < gridSize))
}

func countNeighbours(cells *Cells, x, y int) int {
	total := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if !(i == 0 && j == 0) && cellInGrid(x+i, y+j) {
				if cellAlive(&cells[x+i][y+j]) {
					total++
				}
			}
		}
	}
	return total
}

func putBlackSpot(img *image.Paletted, c Cell) {
	for i := c.x0; i < c.x1; i++ {
		for j := c.y0; j < c.y1; j++ {
			img.Set(i, j, color.Black)
		}
	}
}

func putWhiteSpot(img *image.Paletted, c Cell) {
	for i := c.x0; i < c.x1; i++ {
		for j := c.y0; j < c.y1; j++ {
			img.Set(i, j, color.White)
		}
	}
}

func (cells *Cells) tick() *image.Paletted {
	rect := image.Rect(0, 0, imageSize, imageSize)
	image := image.NewPaletted(rect, palette)
	prev := *cells
	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			n := countNeighbours(&prev, x, y)
			if cellAlive(&prev[x][y]) {
				if n > 3 || n < 2 {
					cellDies(&cells[x][y])
					putWhiteSpot(image, cells[x][y])
				} else {
					cellLives(&cells[x][y])
					putBlackSpot(image, cells[x][y])
				}
			} else {
				if n == 3 {
					cellLives(&cells[x][y])
					putBlackSpot(image, cells[x][y])
				}
			}
		}
	}
	return image
}

func gameOfLife() {
	cells := randomCells()
	anim := gif.GIF{LoopCount: int(*nCycles)}
	for i := uint64(0); i < *nCycles; i++ {
		image := cells.tick()
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, image)
	}
	f, _ := os.OpenFile(*outName, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, &anim)
}
