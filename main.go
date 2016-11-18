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
	cellSize = 5
	delay    = 25
)

var palette = []color.Color{color.White, color.Black}

var nCycles = flag.Int("n", 100, "number of life cycles")
var outName = flag.String("o", "life.gif", "output file name")
var gridX   = flag.Int("x", 100, "grid width")
var gridY   = flag.Int("y", 100, "grid height")

var imageX  = 500
var imageY  = 500

type Cell struct {
	x0, y0, x1, y1 int
	alive          bool
}

type Game struct {
	cells [][]Cell
	prev  [][]Cell
}

func main() {
	flag.Parse()
	setImageSize()
	gameOfLife()
}

func setImageSize() {
	imageX = (*gridX) * 5
	imageY = (*gridY) * 5
}

func newGame() *Game {
	game := &Game{}
	game.cells = make([][]Cell, *gridX)
	game.prev = make([][]Cell, *gridX)
	for x := 0; x < *gridX; x++ {
		game.cells[x] = make([]Cell, *gridY)
		game.prev[x] = make([]Cell, *gridY)
		for y := 0; y < *gridY; y++ {
			game.cells[x][y] = Cell{
				x * cellSize,
				y * cellSize,
				(x + 1) * cellSize,
				(y + 1) * cellSize,
				false,
			}
		}
	}
	return game
}

func randomGame() *Game {
	g := newGame()
	randomizeGame(g)
	return g
}

func randomizeGame(g *Game) {
	rand.Seed(time.Now().UnixNano())
	for x := 0; x < *gridX; x++ {
		for y := 0; y < *gridY; y++ {
			if rand.Float32() > 0.5 {
				g.cellLives(x, y)
			} else {
				g.cellDies(x, y)
			}
		}
	}
}

func (g *Game) cellLives(x, y int) {
	g.cells[x][y].alive = true
}

func (g *Game) cellDies(x, y int) {
	g.cells[x][y].alive = false
}

func (g *Game) cellAlive(x, y int) bool {
	return g.cells[x][y].alive
}

func (g *Game) neighbourAlive(x, y int) bool {
	return g.prev[x][y].alive
}

func outOfBounds(x, y int) bool {
	return (x < 0 || y < 0) || (x > *gridX - 1 || y > *gridY - 1)
}

func (g *Game) saveState() {
	for x := 0; x < *gridX; x++ {
		for y := 0; y < *gridY; y++ {
			g.prev[x][y] = g.cells[x][y]
		}
	}
}

func (g *Game) countNeighbours(x, y int) int {
	total := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if !(i == 0 && j == 0) && !outOfBounds(x + i, y + j) {
				if g.neighbourAlive(x + i, y + j) {
					total++
				}
			}
		}
	}
	return total
}

func fillCell(img *image.Paletted, c Cell) {
	for i := c.x0; i < c.x1; i++ {
		for j := c.y0; j < c.y1; j++ {
			img.Set(i, j, color.Black)
		}
	}
}

func clearCell(img *image.Paletted, c Cell) {
	for i := c.x0; i < c.x1; i++ {
		for j := c.y0; j < c.y1; j++ {
			img.Set(i, j, color.White)
		}
	}
}

func (g *Game) tick() *image.Paletted {
	rect := image.Rect(0, 0, imageX, imageY)
	image := image.NewPaletted(rect, palette)
	g.saveState()
	for x := 0; x < *gridX; x++ {
		for y := 0; y < *gridY; y++ {
			n := g.countNeighbours(x, y)
			if g.cellAlive(x, y) {
				if n > 3 || n < 2 {
					g.cellDies(x, y)
					clearCell(image, g.cells[x][y])
				} else {
					g.cellLives(x, y)
					fillCell(image, g.cells[x][y])
				}
			} else {
				if n == 3 {
					g.cellLives(x, y)
					fillCell(image, g.cells[x][y])
				}
			}
		}
	}
	return image
}

func gameOfLife() {
	game := randomGame()
	anim := gif.GIF{LoopCount: int(*nCycles)}
	for i := 0; i < *nCycles; i++ {
		image := game.tick()
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, image)
	}
	f, _ := os.OpenFile(*outName, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, &anim)
}
