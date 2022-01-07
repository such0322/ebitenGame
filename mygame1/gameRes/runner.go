package gameRes

import (
	"bytes"
	"ebitenGame/resources/images"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
	"log"
)

const (
	runnerOX     = 0
	runnerOY     = 32
	runnerWidth  = 32
	runnerHeight = 32

	runnerSpeedSkip = 10
	runnerSpeedRun  = 5
	runnerSpeedStop = 15

	runnerActSkip = 5
	runnerActRun  = 8
	runnerActStop = 4
)

const (
	runnerPngSkip = iota
	runnerPngRun
	runnerPngStop
)

type Runner struct {
	XState int //前后状态 -1,0,1
	YState int
	Img    *ebiten.Image

	count int
}

func NewRunner() *Runner {
	var r Runner
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	r.Img = ebiten.NewImageFromImage(img)
	return &r
}

func (r *Runner) Update(keys []ebiten.Key) error {
	r.count++
	r.XState = 0
	for _, key := range keys {
		if key == ebiten.KeyA {
			r.XState = -1
		} else if key == ebiten.KeyD {
			r.XState = 1
		}
	}

	return nil
}

func (r *Runner) Draw(screen *ebiten.Image) {
	if r.XState != 0 {
		r.Action(screen, runnerActRun, runnerSpeedRun, runnerPngRun)
		return
	}
	if r.YState != 0 {
		r.Action(screen, runnerActSkip, runnerSpeedSkip, runnerPngSkip)
		return
	}
	r.Action(screen, runnerActStop, runnerSpeedStop, runnerPngStop)
	return
}

func (r *Runner) Action(screen *ebiten.Image, act, speed, pngRow int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(runnerWidth)/2, -float64(runnerHeight)/2)
	op.GeoM.Translate(ScreenWidth/2, ScreenHeight/2)

	i := (r.count / speed) % act
	sx, sy := runnerOX+i*runnerWidth, runnerOY*pngRow
	img := r.Img.SubImage(image.Rect(sx, sy, sx+runnerWidth, sy+runnerHeight)).(*ebiten.Image)
	if r.XState == -1 {
	} else {
	}
	screen.DrawImage(img, op)
}
