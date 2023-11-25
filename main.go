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

var PlayerWidth float64 = 16

type Player struct {
	X         float64
	Y         float64
	Vy        float64
	JumpCount int
}

var RectangleWidth float64 = 32

type Rectangle struct {
	X float64
	Y float64
}

func (r Rectangle) IsOutOfScreen() bool {
	return r.Y > float64(ScreenHeight)
}

func (r Rectangle) IsPlayerOnTop(p Player) bool {
	rl := p.X - PlayerWidth/2
	rr := p.X + PlayerWidth/2
	return (r.X-RectangleWidth/2 <= rl && rl <= r.X+RectangleWidth/2 || r.X-RectangleWidth/2 <= rr && rr <= r.X+RectangleWidth/2) && r.Y-RectangleWidth/2 <= p.Y+PlayerWidth/2 && p.Y+PlayerWidth/2 <= r.Y
}

func (r Rectangle) IsPlayerOnBottom(p Player) bool {
	rl := p.X - PlayerWidth/2
	rr := p.X + PlayerWidth/2
	return (r.X-RectangleWidth/2 <= rl && rl <= r.X+RectangleWidth/2 || r.X-RectangleWidth/2 <= rr && rr <= r.X+RectangleWidth/2) && p.X-PlayerWidth/2 <= r.X+RectangleWidth/2 && r.Y <= p.Y-PlayerWidth/2 && p.Y-PlayerWidth/2 <= r.Y+RectangleWidth/2
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
	PlayerImage = ebiten.NewImage(int(PlayerWidth), int(PlayerWidth))
	PlayerImage.Fill(color.RGBA{0xff, 0x0, 0xff, 0xff})

	RectangleImage = ebiten.NewImage(int(RectangleWidth), int(RectangleWidth))
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
		x := float64(rand.Intn(ScreenWidth))

		g.Rectangles = append(g.Rectangles, Rectangle{
			X: x,
			Y: 0,
		})
	}

	remove := []int{}
	for i, r := range g.Rectangles {
		rv := 2
		if g.Counter > 200 {
			rv = 0
		}

		g.Rectangles[i] = Rectangle{
			X: r.X,
			Y: r.Y + float64(rv),
		}

		if g.Rectangles[i].IsOutOfScreen() {
			remove = append(remove, i)
		}

		if g.Rectangles[i].IsPlayerOnTop(g.Player) {
			g.Player.Y = r.Y - RectangleWidth/2 - PlayerWidth/2
		} else if g.Rectangles[i].IsPlayerOnBottom(g.Player) {
			g.Player.Y = r.Y + RectangleWidth/2 + PlayerWidth/2
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
	option.GeoM.Translate(g.Player.X-PlayerWidth/2, g.Player.Y-PlayerWidth/2)

	screen.DrawImage(PlayerImage, &option)

	for _, r := range g.Rectangles {
		option.GeoM.Reset()
		option.GeoM.Translate(float64(r.X)-RectangleWidth/2, float64(r.Y)-RectangleWidth/2)

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
