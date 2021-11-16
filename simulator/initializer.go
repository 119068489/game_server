package main

import (
	"game_server/for_game"
)

type Initializer struct {
	for_game.Initializer
}

func NewInitializer() *Initializer {
	p := &Initializer{}
	p.Init(p)
	return p
}
