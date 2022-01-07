package gameRes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
)

type Game struct {
	camera Camera
	world  World
	keys   []ebiten.Key
	runner *Runner
}

func NewGame() *Game {
	var g Game
	runner := NewRunner()
	g.runner = runner
	g.world = NewWorld()

	return &g
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	g.runner.Update(g.keys)
	g.camera.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.world.Draw(screen)
	g.camera.Render(g.world.Img, screen)
	g.runner.Draw(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
