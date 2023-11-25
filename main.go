package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/myuon/reimagined-dollup/app"
)

var ScreenWidth int = 320
var ScreenHeight int = 240

type Player struct {
	X         float64
	Y         float64
	Vy        float64
	JumpCount int
}

type Rectangle struct {
	X int
	Y int
}

func (r Rectangle) IsOutOfScreen() bool {
	return r.Y > ScreenHeight
}

var PlayerImage *ebiten.Image
var RectangleImage *ebiten.Image
var FloorY float64 = 220

type Game struct {
	Counter      int
	Player       Player
	KeysPressing map[ebiten.Key]int
	Rectangles   []Rectangle
}

func init() {
	PlayerImage = ebiten.NewImage(16, 16)
	PlayerImage.Fill(color.RGBA{0xff, 0x0, 0xff, 0xff})

	RectangleImage = ebiten.NewImage(32, 32)
	RectangleImage.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})
}

func (g *Game) Update() error {
	g.Counter++

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.KeysPressing[ebiten.KeyLeft]++
	} else {
		g.KeysPressing[ebiten.KeyLeft] = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.KeysPressing[ebiten.KeyRight]++
	} else {
		g.KeysPressing[ebiten.KeyRight] = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.KeysPressing[ebiten.KeyUp]++
	} else {
		g.KeysPressing[ebiten.KeyUp] = 0
	}

	if g.KeysPressing[ebiten.KeyLeft] > 0 {
		g.Player.X -= 4 * app.EaseOutSine(float64(g.KeysPressing[ebiten.KeyLeft])/4)
	}
	if g.KeysPressing[ebiten.KeyRight] > 0 {
		g.Player.X += 4 * app.EaseOutSine(float64(g.KeysPressing[ebiten.KeyRight])/4)
	}
	if g.KeysPressing[ebiten.KeyUp] == 1 && g.Player.JumpCount < 2 {
		g.Player.Vy = -12
		g.Player.JumpCount++
	}

	if g.Player.Y < FloorY {
		g.Player.Vy += 0.85
	}

	g.Player.Y += g.Player.Vy
	if g.Player.Y > FloorY {
		g.Player.Y = FloorY
		g.Player.Vy = 0
		g.Player.JumpCount = 0
	}

	if len(g.Rectangles) < 10 && g.Counter%30 == 0 {
		x := rand.Intn(ScreenWidth)

		g.Rectangles = append(g.Rectangles, Rectangle{
			X: x,
			Y: 0,
		})
	}

	remove := []int{}
	for i, r := range g.Rectangles {
		g.Rectangles[i] = Rectangle{
			X: r.X,
			Y: r.Y + 2,
		}

		if g.Rectangles[i].IsOutOfScreen() {
			remove = append(remove, i)
		}
	}

	if len(remove) > 0 {
		newRectangles := []Rectangle{}
		for i, r := range g.Rectangles {
			shouldRemove := false
			for _, j := range remove {
				if i == j {
					shouldRemove = true
					break
				}
			}

			if !shouldRemove {
				newRectangles = append(newRectangles, r)
			}
		}

		g.Rectangles = newRectangles
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Y=%v Vy=%v", g.Player.Y, g.Player.Vy))

	option := ebiten.DrawImageOptions{}
	option.GeoM.Translate(g.Player.X, g.Player.Y)

	screen.DrawImage(PlayerImage, &option)

	for _, r := range g.Rectangles {
		option.GeoM.Reset()
		option.GeoM.Translate(float64(r.X), float64(r.Y))

		screen.DrawImage(RectangleImage, &option)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{
		Counter: 0,
		Player: Player{
			X:         float64(ScreenWidth) / 2,
			Y:         FloorY,
			Vy:        0,
			JumpCount: 0,
		},
		KeysPressing: make(map[ebiten.Key]int),
		Rectangles:   []Rectangle{},
	}); err != nil {
		log.Fatal(err)
	}
}
