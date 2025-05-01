package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 640
	screenHeight = 480
	boxSize      = 50
	boxSpeed     = 4
	maxMisses    = 5
)

type Game struct {
	boxX         float64
	boxY         float64
	dirX         float64
	dirY         float64
	score        int
	misses       int
	gameOver     bool
	mousePressed bool
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{240, 240, 240, 255})

	box := ebiten.NewImage(boxSize, boxSize)
	box.Fill(color.RGBA{255, 0, 0, 255})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.boxX, g.boxY)
	screen.DrawImage(box, op)

	msg := fmt.Sprint("Score: %d Misses: %d/%d", g.score, g.misses, maxMisses)
	if g.gameOver {
		msg = "Game Over! Click to exit."
	}
	text.Draw(screen, msg, basicfont.Face7x13, 10, 20, color.Black)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}

	g.boxX += g.dirX
	g.boxY += g.dirY

	if g.boxX < 0 || g.boxX+boxSize > screenWidth {
		g.dirX *= -1
	}
	if g.boxY < 0 || g.boxY+boxSize > screenHeight {
		g.dirY *= -1
	}

	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if pressed && !g.mousePressed {
		x, y := ebiten.CursorPosition()
		if float64(x) >= g.boxX && float64(x) <= g.boxX+boxSize &&
			float64(y) >= g.boxY && float64(y) <= g.boxY+boxSize {
			g.score++
			g.teleportBox()
		} else {
			g.misses++
			if g.misses >= maxMisses {
				g.gameOver = true
			}
		}
	}
	g.mousePressed = pressed

	return nil

}

func (g *Game) teleportBox() {
	g.boxX = float64(rand.Intn(screenWidth - boxSize))
	g.boxY = float64(rand.Intn(screenHeight - boxSize))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		boxX: 100,
		boxY: 100,
		dirX: boxSpeed,
		dirY: boxSpeed,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Click & Dodge")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
