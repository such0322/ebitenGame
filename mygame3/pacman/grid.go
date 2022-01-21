package pacman

import (
	"image/color"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"ebitenGame/mygame3/assets"
	"ebitenGame/mygame3/spritetools"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

)

const GridViewSize = 1024

var GrayColor = color.RGBA{236, 240, 241, 255.0}

func GridView(
	characters *assets.Characters,
	powers *assets.Powers,
	arcadeFont *truetype.Font,
	mazeView func(state gameState, data *Data) (*ebiten.Image, error),
) (func(state gameState, data *Data) (*ebiten.Image, error), error) {
	fontface := truetype.NewFace(arcadeFont, &truetype.Options{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	limeAlpha := color.RGBA{250, 233, 8, 200}

	dot := ebiten.NewImage(8, 8)
	dot.Fill(limeAlpha)

	pacman := spritetools.ScaleSprite(characters.Pacman, 0.5, 0.5)
	ghost1 := spritetools.ScaleSprite(characters.Ghost1, 0.5, 0.5)
	ghost2 := spritetools.ScaleSprite(characters.Ghost2, 0.5, 0.5)
	ghost3 := spritetools.ScaleSprite(characters.Ghost3, 0.5, 0.5)
	ghost4 := spritetools.ScaleSprite(characters.Ghost4, 0.5, 0.5)

	left := spritetools.ScaleSprite(powers.Life, 0.5, 0.5)
	invinci := spritetools.ScaleSprite(powers.Invincibility, 0.5, 0.5)
	view := ebiten.NewImage(64*Columns, GridViewSize)

	return func(state gameState, data *Data) (*ebiten.Image, error) {
		view.Clear()
		view.Fill(color.Black)
		ops := &ebiten.DrawImageOptions{}
		switch state {
		case GameLoading:
			text.Draw(view, "PRESS SPACE", fontface, 320-176, 512-(10*32), color.White)
			text.Draw(view, "TO BEGIN", fontface, 320-128, 512+(10), color.White)
		case GameStart, GamePause, GameOver:
			mazeView, err := mazeView(state, data)
			if err != nil {
				return nil, err
			}
			ops.GeoM.Reset()
			ops.GeoM.Translate(0, -(float64(len(data.grid)*CellSize) - (GridViewSize + data.gridOffsetY)))
			view.DrawImage(mazeView, ops)

			for i := 0; i < len(data.active); i++ {
				for j := 0; j < Columns; j++ {
					if !data.active[i][j] {
						ops.GeoM.Reset()
						ops.GeoM.Translate(float64((j*CellSize)+30), -(float64(((i*CellSize)+(CellSize/2))+2) - (GridViewSize + data.gridOffsetY)))
						view.DrawImage(dot, ops)
					}
				}
			}

			for i := 0; i < len(data.powers); i++ {
				power := data.powers[i]
				powerImg := left
				if power.kind == Invincibility {
					powerImg = invinci
				}
				pwidth, pheight := powerImg.Size()
				ops.GeoM.Reset()
				ops.GeoM.Translate(float64((data.powers[i].cellX*CellSize)+pwidth/2), -(float64(((data.powers[i].cellY*CellSize)+(CellSize/2))+pheight/2) - (GridViewSize + data.gridOffsetY)))
				view.DrawImage(powerImg, ops)
			}
			ops.GeoM.Reset()
			pwidth, pheight := pacman.Size()
			switch data.pacman.direction {
			case North:
				ops.GeoM.Rotate(-1.5708)
				ops.GeoM.Translate(data.pacman.posX-float64(pwidth/2), GridViewSize-(data.pacman.posY-float64(pheight-(pheight/2))))
			case East:
				ops.GeoM.Translate(data.pacman.posX-float64(pwidth/2), GridViewSize-(data.pacman.posY+float64(pheight/2)))
			case South:
				ops.GeoM.Rotate(1.5708)
				ops.GeoM.Translate(data.pacman.posX+float64(pwidth/2), GridViewSize-(data.pacman.posY+float64(pheight/2)))
			case West:
				ops.GeoM.Rotate(3.14159)
				ops.GeoM.Translate(data.pacman.posX+float64(pwidth/2), GridViewSize-(data.pacman.posY-float64(pheight-(pheight/2))))
			}
			view.DrawImage(pacman, ops)
			for i := 0; i < len(data.ghosts); i++ {
				ghost := data.ghosts[i]
				ghostImg := ghost1
				switch ghost.kind {
				case Ghost2:
					ghostImg = ghost2
				case Ghost3:
					ghostImg = ghost3
				case Ghost4:
					ghostImg = ghost4
				}
				gwidth, gheight := ghostImg.Size()
				ops.GeoM.Reset()
				if data.invincible {
					ops.ColorM.ChangeHSV(0, 0, 1)
				}
				ops.GeoM.Translate(data.ghosts[i].posX-float64(gwidth/2), (GridViewSize+data.gridOffsetY)-(data.ghosts[i].posY+float64(gheight-(gheight/2))))
				view.DrawImage(ghostImg, ops)
			}
			if state == GamePause {
				back := ebiten.NewImage(389, 130)
				back.Fill(color.Black)
				text.Draw(back, "GAME PAUSED", fontface, 24, 65-(10), color.White)
				text.Draw(back, "PRESS SPACE", fontface, 24, 65+(10+31), color.White)
				ops.GeoM.Reset()
				ops.GeoM.Translate(320-(389/2), 512-(130/2))
				view.DrawImage(back, ops)
			} else if state == GameOver {
				back := ebiten.NewImage(389, 130)
				back.Fill(color.Black)
				text.Draw(back, "GAME OVER", fontface, 56, 65-(10), color.White)
				text.Draw(back, "PRESS SPACE", fontface, 24, 65+(10+31), color.White)
				ops.GeoM.Reset()
				ops.GeoM.Translate(320-(389/2), 512-(130/2))
				view.DrawImage(back, ops)
			}
		}
		return view, nil
	}, nil

}
