package main

type BasePlayer struct {
	g  *Game
	in InputHandler
}

func NewBasePlayer(g *Game, in InputHandler) *BasePlayer {
	return &BasePlayer{
		g:  g,
		in: in,
	}
}

type HumanPlayer struct {
	BasePlayer
}
