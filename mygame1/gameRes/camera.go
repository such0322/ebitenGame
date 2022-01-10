package gameRes

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
	"math"
)

type Camera struct {
	ViewPort   f64.Vec2
	Position   f64.Vec2
	ZoomFactor int
	Rotation   int
}

func NewCamera() *Camera {
	return &Camera{}
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %d, S: %d",
		c.Position, c.Rotation, c.ZoomFactor,
	)
}

func (c *Camera) viewportCenter() f64.Vec2 {
	return f64.Vec2{
		c.ViewPort[0] * 0.5 / 3,
		c.ViewPort[1] * 0.5 / 3,
	}
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position[0], -c.Position[1])
	m.Translate(-c.viewportCenter()[0], -c.viewportCenter()[1])
	m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
	return m
}

func (c *Camera) Render(world, screen *ebiten.Image) {
	screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: c.worldMatrix(),
	})
}

func (c *Camera) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		c.Position[0] -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		c.Position[0] += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		c.Position[1] -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		c.Position[1] += 1
	}
	//if ebiten.IsKeyPressed(ebiten.KeyQ) {
	//	if c.ZoomFactor > -2400 {
	//		c.ZoomFactor -= 1
	//	}
	//}
	//if ebiten.IsKeyPressed(ebiten.KeyE) {
	//	if c.ZoomFactor < 2400 {
	//		c.ZoomFactor += 1
	//	}
	//}
	//if ebiten.IsKeyPressed(ebiten.KeyR) {
	//	c.Rotation += 1
	//}
	//if ebiten.IsKeyPressed(ebiten.KeySpace) {
	//	c.Reset()
	//}
}
