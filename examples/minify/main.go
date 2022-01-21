package main

import (
	"bytes"
	"ebitenGame/resources/images"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	_ "image/jpeg"
	"log"
	"math"
)

const (
	screenWidth  = 800
	screenHeight = 480
)

var gophersImage *ebiten.Image

type Game struct {
	rotate  bool
	clip    bool
	counter int
}

func (g *Game) Update() error {
	g.counter++
	if g.counter == 480 {
		g.counter = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.rotate = !g.rotate
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.clip = !g.clip
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	s := 1.5 / math.Pow(1.01, float64(g.counter))
	msg := fmt.Sprintf(`Minifying images (Nearest filter vs Linear filter):
Press R to rotate the images.
Press C to clip the images.
Scale: %0.2f`, s)
	ebitenutil.DebugPrint(screen, msg)

	clippedGophersImage := gophersImage.SubImage(image.Rect(100, 100, 200, 200)).(*ebiten.Image)
	for i, f := range []ebiten.Filter{ebiten.FilterNearest, ebiten.FilterLinear} {
		w, h := gophersImage.Size()

		op := &ebiten.DrawImageOptions{}
		if g.rotate {
			op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
			op.GeoM.Rotate(float64(g.counter) / 300 * 2 * math.Pi)
			op.GeoM.Translate(float64(w)/2, float64(h)/2)
		}
		op.GeoM.Scale(s, s)
		op.GeoM.Translate(32+float64(i*w)*s+float64(i*4), 64)
		op.Filter = f
		if g.clip {
			screen.DrawImage(clippedGophersImage, op)
		} else {
			screen.DrawImage(gophersImage, op)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	img, _, err := image.Decode(bytes.NewReader(images.Gophers_jpg))
	if err != nil {
		log.Fatal(err)
	}
	gophersImage = ebiten.NewImageFromImage(img)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("moz Minify Demo")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
