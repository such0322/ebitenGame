package pacman

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func spaceReleased() bool {
	return inpututil.IsKeyJustReleased(ebiten.KeySpace)
}

func upKeyPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyUp)
}

func downKeyPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyDown)
}

func leftKeyPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyLeft)
}

func rightKeyPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyRight)
}
