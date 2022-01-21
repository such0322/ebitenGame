package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input int

const (
	KeyNone Input = iota
	KeyDefaultOrGraveyard
	KeyOpponentOrGraveyard
	KeyGraveyard
	KeyLeft
	KeyRight
	KeyUp
	KeyDown
	KeyQuit
)

type InputHandler interface {
	Read() Input
	Cancel()
}

type KeyboardInput struct {
	ch   chan Input
	done chan struct{}
}

func (k *KeyboardInput) Read() Input {
	var in Input
	select {
	case in = <-k.ch:
	default:
		in = KeyNone
	}
	return in
}

func (k *KeyboardInput) Cancel() {
	k.done <- struct{}{}
}

func NewKBInput() *KeyboardInput {
	var in KeyboardInput
	ch := make(chan Input)
	done := make(chan struct{})

	go func(i KeyboardInput) {
		for {
			select {
			case <-done:
				return
			default:
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				ch <- KeyDefaultOrGraveyard
			}
			if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
				ch <- KeyOpponentOrGraveyard
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyG) {
				ch <- KeyGraveyard
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
				ch <- KeyLeft
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
				ch <- KeyRight
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
				ch <- KeyQuit
			}

		}
	}(in)

	in.ch = ch
	in.done = done
	return &in

}
