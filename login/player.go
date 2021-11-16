package login

import (
	"game_server/easygo"
	"time"
)

//=======================================

type Player struct {
	//for_game.Product

	Mutex     easygo.RLock
	WebToken  string
	GameToken string

	Account  string
	PlayerId PLAYER_ID

	OnlineTime int64
	LoginIp    string //登陆游戏ip

	LastStampMap map[interface{}]int64 //最近一次请求时间[key]timeout
}

func NewPlayer(playerId PLAYER_ID, site string) *Player {
	p := &Player{}
	p.Init(playerId, site)
	return p
}

func (self *Player) Init(playerId PLAYER_ID, site string) {
	//self.Product.Init(self, playerId, "玩家外壳")
	self.PlayerId = playerId
	self.LastStampMap = make(map[interface{}]int64)
}

func (self *Player) GetPlayerId() PLAYER_ID {
	return self.PlayerId
}

func (self *Player) SetAccount(account string) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	self.Account = account
}
func (self *Player) GetAccount() string {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	return self.Account
}

func (self *Player) SetLastStamp(key interface{}, expiredSecond int64) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.LastStampMap[key] = time.Now().Unix() + expiredSecond
}

func (self *Player) GetLastStamp(key interface{}) int64 {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	t, ok := self.LastStampMap[key]
	if ok {
		return t
	}
	return 0
}

// rpc是否请求太快
func (self *Player) CheckSetTooFast(rpcName string) bool {
	now := time.Now().Unix()
	lastStamp := self.GetLastStamp(rpcName)
	if lastStamp > now {
		return true
	} else {
		self.SetLastStamp(rpcName, 3)
	}
	return false
}
