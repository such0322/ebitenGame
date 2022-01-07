package keyboard

import (
	"bytes"
	kbres "ebitenGame/resources/images/keyboard"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	_ "image/png"
	"log"
	"strings"
)

const (
	ScreenWidth  = 320
	ScreenHeight = 240
)

const (
	offsetX = 24
	offsetY = 40
)

var keyboardImage *ebiten.Image

func init() {
	img, _, err := image.Decode(bytes.NewReader(kbres.Keyboard_png))
	if err != nil {
		log.Fatal(err)
	}
	keyboardImage = ebiten.NewImageFromImage(img)
}

type Game struct {
	keys []ebiten.Key
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var op ebiten.DrawImageOptions
	op.GeoM.Translate(offsetX, offsetY)
	op.ColorM.Scale(0.5, 0.5, 0.5, 1)
	screen.DrawImage(keyboardImage, &op)

	op = ebiten.DrawImageOptions{}
	for _, key := range g.keys {
		op.GeoM.Reset()
		r, ok := KeyRect(key)
		if !ok {
			continue
		}
		op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
		op.GeoM.Translate(offsetX, offsetY)
		screen.DrawImage(keyboardImage.SubImage(r).(*ebiten.Image), &op)
	}

	keyStrs := []string{}
	for _, p := range g.keys {
		keyStrs = append(keyStrs, p.String())
	}
	ebitenutil.DebugPrint(screen, strings.Join(keyStrs, ", "))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
