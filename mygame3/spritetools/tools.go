package spritetools

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

func GetSprite(width, height int, xoffset, yoffset int, src *ebiten.Image) *ebiten.Image {
	sprite := ebiten.NewImage(width, height)
	rect := image.Rect(xoffset, yoffset, xoffset+width, yoffset+height)
	ops := &ebiten.DrawImageOptions{}
	sprite.DrawImage(src.SubImage(rect).(*ebiten.Image), ops)
	return sprite
}

func ScaleSprite(src *ebiten.Image, x, y float64) *ebiten.Image {
	spriteW, spriteH := src.Size()
	sSprite := ebiten.NewImage(int(float64(spriteW)*x), int(float64(spriteH)*y))
	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Scale(x, y)
	sSprite.DrawImage(src, ops)
	return sSprite
}
