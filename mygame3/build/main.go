package main

import (
	"ebitenGame/mygame3/pacman"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	game, err := pacman.NewGame()
	if err != nil {
		log.Fatalln(err)
	}

	ebiten.SetWindowSize(pacman.ScreenWidth/2, pacman.ScreenHeight/2)
	ebiten.SetWindowTitle("moz PACMAN Demo")

	if err = ebiten.RunGame(game); err != nil {
		log.Fatalln(err)
	}
}
