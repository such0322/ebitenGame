package main

import (
	"ebitenGame/examples/keyboard/keyboard"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(keyboard.ScreenWidth*2, keyboard.ScreenHeight*2)
	ebiten.SetWindowTitle("moz keyboard demo")

	err := ebiten.RunGame(&keyboard.Game{})
	if err != nil {
		fmt.Println(err)
	}
}
