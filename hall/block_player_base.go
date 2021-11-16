// 玩家基础数据块
package hall

//
//type PlayerBase struct {
//	//for_game.PlayerBase `bson:",inline,omitempty"`
//	CallId int64 `bson:"-"` //是否正在通话中
//	// 这里不允许加需要存档的数据，如果要存档的，要加到基类去
//}
//
//func NewPlayerBase(playerId PLAYER_ID) *PlayerBase {
//	p := &PlayerBase{}
//	p.Init(playerId)
//	return p
//}
//
//func (self *PlayerBase) Init(playerId PLAYER_ID) {
//	self.PlayerBase.Init(self, playerId)
//	if !self.LoadFromDB() {
//		panic(fmt.Sprintf("无效的玩家ID=", playerId))
//	}
//}
//
//func (self *PlayerBase) GetClientEndpoint() for_game.IClientEndpoint { // override
//	return ClientEpMp.LoadEndpoint(self.GetPlayerId())
//}
//func (self *PlayerBase) GetEndpoint() IGameClientEndpoint {
//	return ClientEpMp.LoadEndpoint(self.GetPlayerId())
//}
//func (self *PlayerBase) DirtyEventHandler(isAll ...bool) { // override
//	// easygo.Spawn(PlayerBaseFtr.SaveProductToDB, *self.PlayerId) // todo 即时存盘，以后再改为延迟存盘
//	self.SaveToDB(isAll...)
//}
