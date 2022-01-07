package main

import (
	"ebitenGame/mygame1/gameRes"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	game := gameRes.NewGame()
	ebiten.SetWindowSize(gameRes.ScreenWidth, gameRes.ScreenHeight)
	ebiten.SetWindowTitle("moz game demo")

	err := ebiten.RunGame(game)
	if err != nil {
		log.Fatal()
	}
}
