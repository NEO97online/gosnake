package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	gridSize     = 20
	numColumns   = screenWidth / gridSize
	numRows      = screenHeight / gridSize
)

const (
	dirNone = iota
	dirLeft
	dirRight
	dirDown
	dirUp
)

type Position struct {
	X int
	Y int
}

type Game struct {
	body          []Position
	apple         Position
	timer         int
	speed         int
	score         int
	bestScore     int
	moveDirection int
	snakeColor    color.Color
	appleColor    color.Color
}

func (g *Game) resetBody() {
	startPos := Position{X: numColumns / 2, Y: numRows / 2}
	for i := 0; i < len(g.body); i++ {
		g.body[i].X = startPos.X
		g.body[i].Y = startPos.Y
	}
}

func (g *Game) reset() {
	g.apple.X = 3
	g.apple.Y = 3
	g.speed = 15
	g.score = 0
	g.body = g.body[:3]
	g.moveDirection = dirNone
	g.resetBody()
}

func (g *Game) Update(screen *ebiten.Image) error {
	// handle key inputs
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyS) && g.moveDirection != dirUp:
		g.moveDirection = dirDown
	case inpututil.IsKeyJustPressed(ebiten.KeyW) && g.moveDirection != dirDown:
		g.moveDirection = dirUp
	case inpututil.IsKeyJustPressed(ebiten.KeyD) && g.moveDirection != dirLeft:
		g.moveDirection = dirRight
	case inpututil.IsKeyJustPressed(ebiten.KeyA) && g.moveDirection != dirRight:
		g.moveDirection = dirLeft
	}

	// check if snake speed interval has passed, if so, move the snake
	if g.timer%(60/g.speed) == 0 && g.moveDirection != dirNone {
		// move snake body
		for i := len(g.body) - 1; i > 0; i-- {
			g.body[i].X = g.body[i-1].X
			g.body[i].Y = g.body[i-1].Y
		}

		// move snake head
		switch g.moveDirection {
		case dirLeft:
			g.body[0].X -= 1
		case dirRight:
			g.body[0].X += 1
		case dirUp:
			g.body[0].Y -= 1
		case dirDown:
			g.body[0].Y += 1
		}

		// nom yourself :(
		for i := len(g.body) - 1; i > 0; i-- {
			if g.body[i].X == g.body[0].X && g.body[i].Y == g.body[0].Y {
				g.reset()
				return nil
			}
		}

		// wrap across screen
		if g.body[0].X < 0 {
			g.body[0].X = numColumns - 1
		}
		if g.body[0].X > numColumns-1 {
			g.body[0].X = 0
		}
		if g.body[0].Y < 0 {
			g.body[0].Y = numRows - 1
		}
		if g.body[0].Y > numRows-1 {
			g.body[0].Y = 0
		}

		// nom the apple
		if g.body[0].X == g.apple.X && g.body[0].Y == g.apple.Y {
			g.apple.X = rand.Intn(numColumns - 1)
			g.apple.Y = rand.Intn(numRows - 1)
			g.body = append(g.body, Position{X: g.body[0].X, Y: g.body[0].Y})
			g.score++
			if g.score > g.bestScore {
				g.bestScore = g.score
			}
		}

	}

	g.timer++

	return nil
}

func drawTile(screen *ebiten.Image, x, y int, c color.Color) {
	ebitenutil.DrawRect(screen, float64(x*gridSize), float64(y*gridSize), gridSize, gridSize, c)
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw snake body
	for i := 0; i < len(g.body); i++ {
		pos := g.body[i]
		drawTile(screen, pos.X, pos.Y, g.snakeColor)
	}

	// Draw apple
	drawTile(screen, g.apple.X, g.apple.Y, g.appleColor)

	// Draw GUI
	if g.moveDirection == dirNone {
		ebitenutil.DebugPrint(screen, "Press WASD to start")
	} else {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f Score: %d Best Score: %d", ebiten.CurrentFPS(), g.score, g.bestScore))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func newGame() *Game {
	g := &Game{
		speed:         15,
		body:          make([]Position, 3),
		apple:         Position{X: 3, Y: 3},
		moveDirection: dirNone,
		score:         0,
		bestScore:     0,
		snakeColor:    color.RGBA{0, 0xff, 0, 0xff},
		appleColor:    color.RGBA{0xff, 0xff, 0xff, 0xff},
	}
	g.resetBody()
	return g
}

func main() {
	fmt.Println("Starting snake game... have fun :)")

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake")
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
