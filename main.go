package main

import (
	"flag"
	"fmt"
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

type cell struct {
	x0, y0, x1, y1 int
	alive          bool
}

type GOL struct {
	gridw, gridh int
	ncycles      int
	cells        [][]cell
	prev         [][]cell
}

func main() {
	var (
		outname               string
		ncycles, gridw, gridh int
	)
	flag.StringVar(&outname, "o", "life.gif", "output file name")
	flag.IntVar(&ncycles, "n", 100, "number of life cycles")
	flag.IntVar(&gridw, "w", 100, "grid width")
	flag.IntVar(&gridh, "h", 100, "grid height")
	flag.Parse()

	g := newGOL(gridw, gridh, ncycles)
	g.randomize()
	animation := g.run()
	if err := saveAnimation(animation, outname); err != nil {
		fmt.Fprintf(os.Stderr, "error saving animation: %s\n", err.Error())
		return
	}
	fmt.Printf("game saved in %s\n", outname)
}

func newGOL(w, h, n int) *GOL {
	game := &GOL{gridw: w, gridh: h, ncycles: n}
	game.cells = make([][]cell, w)
	game.prev = make([][]cell, h)
	for x := 0; x < w; x++ {
		game.cells[x] = make([]cell, w)
		game.prev[x] = make([]cell, h)
		for y := 0; y < h; y++ {
			game.cells[x][y] = cell{
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

func (g *GOL) randomize() {
	rand.Seed(time.Now().UnixNano())
	for x := 0; x < g.gridw; x++ {
		for y := 0; y < g.gridh; y++ {
			if rand.Float32() > 0.5 {
				g.activateCell(x, y)
			} else {
				g.killCell(x, y)
			}
		}
	}
}

func (g *GOL) activateCell(x, y int) {
	g.cells[x][y].alive = true
}

func (g *GOL) killCell(x, y int) {
	g.cells[x][y].alive = false
}

func (g *GOL) isCellAlive(x, y int) bool {
	return g.cells[x][y].alive
}

func (g *GOL) isNeigbourAlive(x, y int) bool {
	return g.prev[x][y].alive
}

func (g *GOL) outOfBounds(x, y int) bool {
	return (x < 0 || y < 0) || (x > g.gridw-1 || y > g.gridh-1)
}

func (g *GOL) saveState() {
	for x := 0; x < g.gridw; x++ {
		for y := 0; y < g.gridh; y++ {
			g.prev[x][y] = g.cells[x][y]
		}
	}
}

func (g *GOL) countNeighbours(x, y int) int {
	total := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if !(i == 0 && j == 0) && !g.outOfBounds(x+i, y+j) {
				if g.isNeigbourAlive(x+i, y+j) {
					total++
				}
			}
		}
	}
	return total
}

func colorCell(img *image.Paletted, c cell, clr color.Gray16) {
	for i := c.x0; i < c.x1; i++ {
		for j := c.y0; j < c.y1; j++ {
			img.Set(i, j, clr)
		}
	}
}

func drawMesh(img *image.Paletted, w, h int) {
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			if i%cellSize == 0 || j%cellSize == 0 {
				img.Set(i, j, color.Black)
			}
		}
	}
}

func (g *GOL) tick() *image.Paletted {
	rect := image.Rect(0, 0, g.gridw*cellSize, g.gridh*cellSize)
	image := image.NewPaletted(rect, palette)
	g.saveState()
	for x := 0; x < g.gridw; x++ {
		for y := 0; y < g.gridh; y++ {
			n := g.countNeighbours(x, y)
			if g.isCellAlive(x, y) {
				if n > 3 || n < 2 {
					g.killCell(x, y)
					colorCell(image, g.cells[x][y], color.White)
				} else {
					g.activateCell(x, y)
					colorCell(image, g.cells[x][y], color.Black)
				}
				continue
			}

			if n == 3 {
				g.activateCell(x, y)
				colorCell(image, g.cells[x][y], color.Black)
			}
		}
	}
	return image
}

func (g *GOL) run() gif.GIF {
	anim := gif.GIF{LoopCount: g.ncycles}
	for i := 0; i < g.ncycles; i++ {
		image := g.tick()
		drawMesh(image, g.gridw*cellSize, g.gridh*cellSize)
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, image)
	}
	return anim
}

func saveAnimation(anim gif.GIF, filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	gif.EncodeAll(f, &anim)
	f.Close()
	return nil
}
