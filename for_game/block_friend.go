package for_game

import (
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/astaxie/beego/logs"
)

const (
	ADDFRIEND_NODEAL  int32 = 1 //未处理
	ADDFRIEND_ACCEPT  int32 = 2 //接受
	ADDFRIEND_OUTTIME int32 = 3 //过期
	MAX_FRIEND              = 5000
)

type IFriend interface {
	IMongoProduct
	IEasyPersist
	DirtyEventHandler(isAll ...bool)
}

type FriendBase struct {
	MongoProduct                  `bson:"-"`
	EasyPersist                   `bson:",inline,omitempty"`
	share_message.FriendBase      `bson:",inline,omitempty"` //好友数据
	share_message.AllAddPlayerMsg `bson:",inline,omitempty"` //好友请求数据
	PlayerId                      PLAYER_ID                  `bson:"-"`
	Me                            IFriend                    `bson:"-"`
	DirtyData                     *FriendBase                `bson:"-"`
}

func NewFriendBase(playerId PLAYER_ID) *FriendBase {
	p := &FriendBase{}
	p.Init(p, playerId)
	return p
}

func (self *FriendBase) Init(me IFriend, playerId PLAYER_ID) {
	self.Me = me
	self.MongoProduct.Init(self, playerId, "好友表")
	kwargs1 := easygo.KWAT{
		"DirtyEventHandler": self.DirtyEventHandler,
	}
	self.EasyPersist.Init(self, kwargs1)
	self.PlayerId = playerId
}

func (self *FriendBase) DirtyEventHandler(isAll ...bool) {
	self.SaveToDB(isAll...)
}

func (self *FriendBase) GetPersistenceObj() IPersistence { // override
	return self
}

func (self *FriendBase) GetDirtyData() interface{} { //override
	return self.DirtyData
}
func (self *FriendBase) CleanDirtyData() { //override
	self.DirtyData = &FriendBase{}
}

func (self *FriendBase) GetC() (c *mgo.Collection, fun func()) { // override
	return easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_FRIEND)
}

func (self *FriendBase) OnBron(kwargs ...easygo.KWAT) {

}

func (self *FriendBase) GetFriend(pid PLAYER_ID) *share_message.FriendInfo {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for _, msg := range self.Friends {
		if msg.GetPlayerId() == pid {
			return msg
		}
	}
	return nil
}
func (self *FriendBase) GetFriendIds() []PLAYER_ID {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var plst []PLAYER_ID
	for _, m := range self.Friends {
		plst = append(plst, m.GetPlayerId())
	}
	return plst
}

func (self *FriendBase) GetFriendsReName(pid PLAYER_ID) string {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for _, info := range self.Friends {
		if info.GetPlayerId() == pid {
			return info.GetReName()
		}
	}
	return ""
}
func (self *FriendBase) SetFriendReName(pid PLAYER_ID, name string) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for _, info := range self.Friends {
		if info.GetPlayerId() == pid {
			info.ReName = easygo.NewString(name)
			break
		}
	}
	self.MarkDirty()
}
func (self *FriendBase) SetFriendSetting(pid PLAYER_ID, t int32, value bool) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for _, player := range self.Friends {
		if player.GetPlayerId() == pid {
			if t == 1 {
				player.Setting.IsTopChat = easygo.NewBool(value)
			} else if t == 2 {
				player.Setting.IsNoDisturb = easygo.NewBool(value)
			} else if t == 3 {
				player.Setting.IsAfterReadClear = easygo.NewBool(value)
			} else if t == 4 {
				player.Setting.IsScreenShotNotify = easygo.NewBool(value)
			} else {
				panic("不存在的类型")
			}
			break
		}
	}
	self.MarkDirty()
}

func (self *FriendBase) AddFriend(pid PLAYER_ID, t int32) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for _, info := range self.Friends {
		if info.GetPlayerId() == pid {
			return
		}
	}
	msg := &share_message.FriendInfo{
		PlayerId:   easygo.NewInt64(pid),
		ReName:     easygo.NewString(""),
		Setting:    self.GetDefalutFriendSetting(),
		Type:       easygo.NewInt32(t),
		CreateTime: easygo.NewInt64(GetMillSecond()),
	}
	self.Friends = append(self.Friends, msg)
	base := GetRedisPlayerBase(self.PlayerId)
	base.AddRedisPlayerFriends(pid)
	self.MarkDirty()
}

func (self *FriendBase) GetDefalutFriendSetting() *share_message.FriendSetting {
	msg := &share_message.FriendSetting{
		IsTopChat:          easygo.NewBool(false),
		IsNoDisturb:        easygo.NewBool(false),
		IsAfterReadClear:   easygo.NewBool(false),
		IsScreenShotNotify: easygo.NewBool(false),
	}
	return msg
}

func (self *FriendBase) DelFriend(pid int64) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var lst []*share_message.FriendInfo
	for _, info := range self.Friends {
		if info.GetPlayerId() == pid {
			continue
		}
		lst = append(lst, info)
	}
	self.Friends = lst
	base := GetRedisPlayerBase(self.PlayerId)
	base.DelRedisPlayerFriends(pid)
	self.MarkDirty()
}

func (self *FriendBase) CheckMaxNum() bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if len(self.Friends) >= MAX_FRIEND {
		return false
	}
	return true
}

func (self *FriendBase) CheckAddFriendRequest() { //把超过7天的好友请求删除  客户端本地存库
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	t := time.Now().Unix()
	var lst []*share_message.AddPlayerRequest
	var b bool
	for _, info := range self.GetAddPlayerRequest() {
		if t-info.GetTime() > 7*86400 {
			b = true
			continue
		}
		lst = append(lst, info)
	}
	if b {
		self.AddPlayerRequest = lst
		self.MarkDirty()
	}
}

func (self *FriendBase) AddFriendRequest(msg *share_message.AddPlayerRequest) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.CheckAddFriendRequest()
	for _, info := range self.AddPlayerRequest {
		if info.GetPlayerId() == msg.GetPlayerId() {
			num := ADDFRIEND_NODEAL
			info.Result = &num
		}
	}
	self.AddPlayerRequest = append(self.AddPlayerRequest, msg)
	self.MarkDirty()
}

func (self *FriendBase) AgreeAddFriend(pid PLAYER_ID) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.CheckAddFriendRequest()
	for _, info := range self.AddPlayerRequest {
		if info.GetPlayerId() != pid {
			continue
		}
		num := ADDFRIEND_ACCEPT
		info.Result = &num
	}
	self.MarkDirty()
}

// 删除好友请求列表中的数据,彻底删除
func (self *FriendBase) DelAddFriend(pid []int64) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.CheckAddFriendRequest()
	if len(self.AddPlayerRequest) == 0 {
		return
	}
	newReqMap := make(map[int64]*share_message.AddPlayerRequest)
	for _, v := range self.AddPlayerRequest {
		newReqMap[v.GetPlayerId()] = v
	}
	// 相等的id
	ids := make([]int64, 0)
	for _, v := range pid {
		for _, info := range self.AddPlayerRequest {
			if info.GetPlayerId() != v {
				continue
			}
			ids = append(ids, v)
		}
	}
	// 从map中删除相等的id
	for _, id := range ids {
		delete(newReqMap, id)
	}
	// 从新封装进新的数组
	result := make([]*share_message.AddPlayerRequest, 0)
	for _, value := range newReqMap {
		result = append(result, value)
	}
	self.AddPlayerRequest = result
	self.MarkDirty()

}

func (self *FriendBase) ReadFriendRequest(ids []int64) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.CheckAddFriendRequest()
	for _, id := range ids {
		for _, msg := range self.AddPlayerRequest {
			if msg.GetPlayerId() == id {
				msg.IsRead = easygo.NewBool(true)
			}
		}
	}
	self.MarkDirty()
}

//
//func (self *FriendBase) GetAllFriendRequestForOne() *share_message.AllAddPlayerMsg {
//	self.Mutex.Lock()
//	defer self.Mutex.Unlock()
//	self.CheckAddFriendRequest()
//	var lst []*share_message.AddPlayerRequest
//	var pInfo = make(map[PLAYER_ID]int64)
//	t := time.Now().Unix()
//	for _, info := range self.AddPlayerRequest {
//		pid := info.GetPlayerId()
//		result := info.GetResult()
//		if result == ADDFRIEND_NODEAL {
//			if t-info.GetTime() > 86400 {
//				info.Result = easygo.NewInt32(ADDFRIEND_OUTTIME) // 未处理的超过一天标记已过期
//			} else if util.Int64InSlice(pid, self.GetFriendIds()) {
//				info.Result = easygo.NewInt32(ADDFRIEND_ACCEPT) // 已同意
//			}
//		}
//
//		pInfo[pid] = info.GetId()
//	}
//
//	for _, id := range pInfo {
//		for _, info := range self.AddPlayerRequest {
//			if id == info.GetId() {
//				lst = append(lst, info)
//				break
//			}
//		}
//	}
//
//	for _, msg := range lst {
//		pid := msg.GetPlayerId()
//		base := GetRedisPlayerBase(pid)
//		if base == nil {
//			continue
//		}
//		msg.NickName = easygo.NewString(base.GetNickName())
//		msg.HeadIcon = easygo.NewString(base.GetHeadIcon())
//		msg.Account = easygo.NewString(base.GetAccount())
//		msg.Phone = easygo.NewString(base.GetPhone())
//		msg.Sex = easygo.NewInt32(base.GetSex())
//	}
//
//	msg := &share_message.AllAddPlayerMsg{
//		AddPlayerRequest: lst,
//	}
//	return msg
//}

// 新版本. 只需要未读和未处理的.
func (self *FriendBase) GetNewVersionAllFriendRequestForOne() *share_message.AllAddPlayerMsg {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.CheckAddFriendRequest()
	var lst []*share_message.AddPlayerRequest
	var pInfo = make(map[PLAYER_ID]int64)
	t := time.Now().Unix()
	for _, info := range self.AddPlayerRequest {
		pid := info.GetPlayerId()
		result := info.GetResult()
		if result == ADDFRIEND_NODEAL {
			if t-info.GetTime() > 3*86400 {
				info.Result = easygo.NewInt32(ADDFRIEND_OUTTIME) // 未处理的超过一天标记已过期
			} else if util.Int64InSlice(pid, self.GetFriendIds()) {
				info.Result = easygo.NewInt32(ADDFRIEND_ACCEPT) // 已同意
			}
		}

		pInfo[pid] = info.GetId()
	}

	for _, info := range self.AddPlayerRequest {
		if !info.GetIsRead() || info.GetResult() == ADDFRIEND_NODEAL {
			//只去未读和未处理的
			lst = append(lst, info)
		}
	}

	for _, msg := range lst {
		pid := msg.GetPlayerId()
		base := GetRedisPlayerBase(pid)
		if base == nil {
			continue
		}
		msg.NickName = easygo.NewString(base.GetNickName())
		msg.HeadIcon = easygo.NewString(base.GetHeadIcon())
		msg.Account = easygo.NewString(base.GetAccount())
		msg.Phone = easygo.NewString(base.GetPhone())
		msg.Sex = easygo.NewInt32(base.GetSex())
		msg.Types = easygo.NewInt32(base.GetTypes())
	}

	msg := &share_message.AllAddPlayerMsg{
		AddPlayerRequest: lst,
	}
	return msg
}

//================================================================================================================
// 对外接口，获取玩家好友数据
func GetFriendBase(pid PLAYER_ID, kwargs ...easygo.KWAT) *FriendBase {
	obj := NewFriendBase(pid)
	if !obj.LoadFromDB(kwargs...) {
		if len(kwargs) > 0 {
			obj.InsertToDB(kwargs...)
		} else {
			logs.Error("无效的玩家ID=" + easygo.AnytoA(pid))
			return nil
		}
	}
	return obj
}

func GetFriendsReName(pid, otherId PLAYER_ID) string {
	base := GetFriendBase(pid)
	if base == nil {
		return ""
	}
	friend := base.GetFriend(otherId)
	if friend == nil {
		return ""
	}
	return friend.GetReName()
}

//func TransFriendData() {
//	col, colfunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, "player_friend")
//	defer colfunc()
//
//	col1, colfunc1 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_FRIEND)
//	defer colfunc1()
//
//	col2, colfunc2 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_FRIEND_REQUEST)
//	defer colfunc2()
//
//	type FriendBase struct {
//		Friend   share_message.FriendBase      `bson:",inline,omitempty"` //好友数据
//		Request  share_message.AllAddPlayerMsg `bson:",inline,omitempty"` //好友请求数据
//		PlayerId int64                         `bson:"_id"`
//	}
//
//	var lst []*FriendBase
//
//	friendInfo := make(map[int64][]*share_message.FriendInfo)
//	reqInfo := make(map[int64][]*share_message.AddPlayerRequest)
//	var playerIds []int64
//	err := col.Find(nil).All(&lst)
//	easygo.PanicError(err)
//	for _, m := range lst {
//		pid := m.PlayerId
//		playerIds = append(playerIds, pid)
//		friends := m.Friend.Friends
//		if len(friends) > 0 {
//			friendInfo[pid] = friends
//		}
//		reqs := m.Request.AddPlayerRequest
//		if len(reqs) > 0 {
//			reqInfo[pid] = reqs
//		}
//	}
//
//	var friendList []interface{}
//	var reqList []interface{}
//	for _, pid := range playerIds {
//		lst := friendInfo[pid]
//		if len(lst) > 0 {
//			for _, m := range lst {
//				m.MyId = easygo.NewInt64(pid)
//				m.LogId = easygo.NewInt64(GetRedisNextId(REDIS_FRIENDS))
//				friendList = append(friendList, m)
//			}
//		}
//		lst1 := reqInfo[pid]
//		if len(lst1) > 0 {
//			for _, m := range lst1 {
//				m.MyId = easygo.NewInt64(pid)
//				m.LogId = easygo.NewInt64(GetRedisNextId(REDIS_FRIENDS_REQUEST))
//				reqList = append(reqList, m)
//			}
//		}
//	}
//	err1 := col1.Insert(friendList...)
//	easygo.PanicError(err1)
//	err2 := col2.Insert(reqList...)
//	easygo.PanicError(err2)
//}
