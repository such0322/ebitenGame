package pacman

import "github.com/hajimehoshi/ebiten/v2/audio"

type Audio struct{
	ctx *audio.Context
	players *AudioPlayers
}

type AudioPlayers struct {
	Beginning *audio.Player
	Chomp     *audio.Player
	Death     *audio.Player
	EatFlask  *audio.Player
	EatGhost  *audio.Player
	ExtraPac  *audio.Player
}