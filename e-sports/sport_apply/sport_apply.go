package sport_apply

import (
	"game_server/for_game"
	"game_server/pb/share_message"
)

type GuessOddsNumsSort []*share_message.GameGuessOddsNumObject

type SportApply struct {
}

func NewSportApply() *SportApply {
	p := &SportApply{}
	p.Init()
	return p
}
func (self *SportApply) Init() {
	for_game.InitRedisCreateBetOrderId()
}
