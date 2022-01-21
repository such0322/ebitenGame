package pacman

import (
	"ebitenGame/mygame3/assets"
	"ebitenGame/mygame3/spritetools"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	MazeViewSize = 1536
	CellSize     = 64
)

func MazeView(walls *assets.Walls) (func(state gameState, data *Data) (*ebiten.Image, error), error) {
	icWallSide := spritetools.ScaleSprite(walls.InActiveSide, 1.0, 1.0)

	icWallCorner := spritetools.ScaleSprite(walls.InActiveCorner, 1.0, 1.0)

	mazeView := ebiten.NewImage(CellSize*Columns, MazeViewSize)

	var lastGrid [][Columns][4]rune

	return func(state gameState, data *Data) (*ebiten.Image, error) {
		if cp, eq := deepEqual(lastGrid, data.grid); eq {
			return mazeView, nil
		} else {
			lastGrid = cp
		}
		mazeView.Clear()
		ops := &ebiten.DrawImageOptions{}
		for i := 0; i < len(data.grid); i++ {
			for j := 0; j < len(data.grid[i]); j++ {
				side := icWallSide
				corner := icWallCorner
				cellWalls := data.grid[i][j]
				if cellWalls[0] == 'N' {
					ops.GeoM.Reset()
					ops.GeoM.Translate(float64(j*CellSize)+12, float64(MazeViewSize-((i*CellSize)+CellSize)))
					mazeView.DrawImage(side, ops)
				}
				if cellWalls[1] == 'E' {
					ops.GeoM.Reset()
					ops.GeoM.Rotate(1.5708)
					ops.GeoM.Translate(float64(j*CellSize)+CellSize, float64(MazeViewSize-((i+CellSize)+52)))
					mazeView.DrawImage(side, ops)
				}
				if cellWalls[2] == 'S' {
					ops.GeoM.Reset()
					ops.GeoM.Translate(float64(j*CellSize)+12, float64(MazeViewSize-((i*CellSize)+12)))
					mazeView.DrawImage(side, ops)
				}
				if cellWalls[3] == 'W' {
					ops.GeoM.Reset()
					ops.GeoM.Rotate(1.5708)
					ops.GeoM.Translate(float64(j*CellSize)+12, float64(MazeViewSize-((i*CellSize)+52)))
					mazeView.DrawImage(side, ops)
				}

				ops.GeoM.Reset()
				ops.GeoM.Translate(float64(j*CellSize)+52, float64(MazeViewSize-((i*CellSize)+CellSize)))
				mazeView.DrawImage(corner, ops)

				ops.GeoM.Reset()
				ops.GeoM.Translate(float64(j*CellSize), float64(MazeViewSize-((i*CellSize)+CellSize)))
				mazeView.DrawImage(corner, ops)

				ops.GeoM.Reset()
				ops.GeoM.Translate(float64(j*CellSize)+52, float64(MazeViewSize-((i*CellSize)+12)))
				mazeView.DrawImage(corner, ops)

				ops.GeoM.Reset()
				ops.GeoM.Translate(float64(j*CellSize), float64(MazeViewSize-((i*CellSize)+12)))
				mazeView.DrawImage(corner, ops)
			}
		}
		return mazeView, nil
	}, nil

}

func deepEqual(previous, next [][Columns][4]rune) ([][Columns][4]rune, bool) {
	deepCopy := func(src [][Columns][4]rune) [][Columns][4]rune {
		cp := make([][Columns][4]rune, 0)
		for i := 0; i < len(next); i++ {
			row := [Columns][4]rune{}
			for j := 0; j < Columns; j++ {
				row[j] = next[i][j]
			}
			cp = append(cp, row)
		}
		return cp
	}
	if len(previous) != len(next) {
		return deepCopy(next), false
	}
	for i := 0; i < len(previous); i++ {
		if previous[i] != next[i] {
			return deepCopy(next), false
		}
		for j := 0; j < Columns; j++ {
			if previous[i][j] != next[i][j] {
				return deepCopy(next), false
			}
		}
	}
	return next, true
}
