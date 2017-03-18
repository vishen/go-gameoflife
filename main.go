package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Cell struct {
	// square Vertex Array Object
	drawable uint32

	alive     bool
	aliveNext bool

	X int
	Y int
}

func NewCell(x, y int) *Cell {
	// Copy the `square` data into a new `points` source
	points := make([]float32, len(square), len(square))
	copy(points, square)

	for i := 0; i < len(points); i++ {
		var position float32
		var size float32
		switch i % 3 {
		case 0:
			// Are we an 'X' coordinate
			size = 1.0 / float32(COLUMNS)
			position = float32(x) * size
		case 1:
			// Are we an 'Y' coordinate
			size = 1.0 / float32(ROWS)
			position = float32(y) * size
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = (position * 2) - 1
		} else {
			points[i] = ((position + size) * 2) - 1
		}
	}

	// Vertex Buffer Object
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// 4 * len(points) -> We are using slice float32 and a 32-bit float has 4 bytes
	number_bytes := 4 * len(points)
	gl.BufferData(gl.ARRAY_BUFFER, number_bytes, gl.Ptr(points), gl.STATIC_DRAW)

	// Vertex Array Object
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return &Cell{X: x, Y: y, drawable: vao}
}

func (c *Cell) checkState(cells [][]*Cell) {
	c.alive = c.aliveNext
	c.aliveNext = c.alive

	liveCount := c.liveNeighbors(cells)
	if c.alive {
		// 1. Any live cell with fewer than two live neighbours dies, as if caused by underpopulation.
		if liveCount < 2 {
			c.aliveNext = false
		}

		// 2. Any live cell with two or three live neighbours lives on to the next generation.
		if liveCount == 2 || liveCount == 3 {
			c.aliveNext = true
		}

		// 3. Any live cell with more than three live neighbours dies, as if by overpopulation.
		if liveCount > 3 {
			c.aliveNext = false
		}
	} else {
		// 4. Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
		if liveCount == 3 {
			c.aliveNext = true
		}
	}
}

func (c *Cell) liveNeighbors(cells [][]*Cell) int {
	var liveCount int
	add := func(x, y int) {
		// If we're at an edge, check the other side of the board.
		if x == len(cells) {
			x = 0
		} else if x == -1 {
			x = len(cells) - 1
		}
		if y == len(cells[x]) {
			y = 0
		} else if y == -1 {
			y = len(cells[x]) - 1
		}

		if cells[x][y].alive {
			liveCount++
		}
	}

	add(c.X-1, c.Y)   // To the left
	add(c.X+1, c.Y)   // To the right
	add(c.X, c.Y+1)   // up
	add(c.X, c.Y-1)   // down
	add(c.X-1, c.Y+1) // top-left
	add(c.X+1, c.Y+1) // top-right
	add(c.X-1, c.Y-1) // bottom-left
	add(c.X+1, c.Y-1) // bottom-right

	return liveCount
}

func makeCells() [][]*Cell {
	// TODO(): Move this into a `type` or something or change it so it is just a single array
	//  will make it much cleaner
	rand.Seed(time.Now().UnixNano())
	cells := make([][]*Cell, ROWS, ROWS)
	for x := 0; x < ROWS; x++ {
		for y := 0; y < COLUMNS; y++ {
			c := NewCell(x, y)

			c.alive = rand.Float64() < THRESHOLD
			c.aliveNext = c.alive

			cells[x] = append(cells[x], c)
		}
	}
	return cells
}

func draw(cells [][]*Cell) {
	// Remove anything from the window
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	var cell *Cell
	for x := range cells {
		for y := range cells[x] {
			cell = cells[x][y]
			if cell.alive {
				gl.BindVertexArray(cell.drawable)
				gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
			}
		}
	}

	glfw.PollEvents()
	window.SwapBuffers()

}

const (
	WIDTH  = 500
	HEIGHT = 500

	ROWS    = 10
	COLUMNS = 10

	THRESHOLD = 0.15
)

var (
	window  *glfw.Window
	program uint32
)

func main() {
	// Ensure that we will always execute in the same OS thread
	runtime.LockOSThread()

	initGlfw()
	defer glfw.Terminate()

	initOpenGL()

	cells := makeCells()

	ticker := time.NewTicker(time.Second / 5)

	for !window.ShouldClose() {

		for x := range cells {
			for y := range cells[x] {
				c := cells[x][y]
				c.checkState(cells)
			}
		}

		draw(cells)

		<-ticker.C

	}
}
