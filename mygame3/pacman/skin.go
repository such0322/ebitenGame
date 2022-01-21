package pacman

import (
	"ebitenGame/mygame3/assets"
	"ebitenGame/mygame3/spritetools"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image/color"
	"strconv"
)

const (
	MaxScoreView = 999999999
	MaxLifes     = 7
)

func SkinView(skin *ebiten.Image, powers *assets.Powers, arcadeFont *truetype.Font) (func(state gameState, data *Data) (*ebiten.Image, error), error) {
	fontface := truetype.NewFace(arcadeFont, &truetype.Options{
		Size:    28,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	width, height := skin.Size()
	view := ebiten.NewImage(width, height)
	life := spritetools.ScaleSprite(powers.Life, 0.5, 0.5)

	return func(state gameState, data *Data) (*ebiten.Image, error) {
		view.Clear()
		view.DrawImage(skin, &ebiten.DrawImageOptions{})

		switch state {
		case GameStart:
			fallthrough
		case GamePause:
			fallthrough
		case GameOver:
			if data != nil {
				score := data.score
				lifes := data.lifes
				if score > MaxScoreView {
					score = MaxScoreView
				}
				numstr := strconv.Itoa(score)
				text.Draw(view, numstr, fontface, 682-(len(numstr)*27), 64, color.White)
				if lifes > MaxLifes {
					lifes = MaxLifes
				}
				ops := &ebiten.DrawImageOptions{}
				width, _ := life.Size()
				for i := 0; i < lifes; i++ {
					ops.GeoM.Reset()
					ops.GeoM.Translate(float64(650-(width*i)), 80)
					view.DrawImage(life, ops)
				}
			}
		}
		return view, nil
	}, nil
}
