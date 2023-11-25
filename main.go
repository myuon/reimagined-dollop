package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/myuon/reimagined-dollup/app"
)

type Player struct {
	X float64
	Y float64
}

var PlayerImage *ebiten.Image

type Game struct {
	Player       Player
	KeysPressing map[ebiten.Key]int
}

func init() {
	PlayerImage = ebiten.NewImage(16, 16)
	PlayerImage.Fill(color.RGBA{0xff, 0x0, 0xff, 0xff})
}

func (g *Game) Update() error {
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

	if g.KeysPressing[ebiten.KeyLeft] > 0 {
		g.Player.X -= 5 * app.EaseOutSine(float64(g.KeysPressing[ebiten.KeyLeft])/5)
	}
	if g.KeysPressing[ebiten.KeyRight] > 0 {
		g.Player.X += 5 * app.EaseOutSine(float64(g.KeysPressing[ebiten.KeyRight])/5)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")

	option := ebiten.DrawImageOptions{}
	option.GeoM.Translate(g.Player.X, g.Player.Y)

	screen.DrawImage(PlayerImage, &option)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{
		Player: Player{
			X: 160,
			Y: 120,
		},
		KeysPressing: make(map[ebiten.Key]int),
	}); err != nil {
		log.Fatal(err)
	}
}
