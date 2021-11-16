package for_game

import (
	"encoding/base64"
	"encoding/json"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

/*
群成员内存数据管理
*/

const (
	NORMAL_TEAMMANAGE = 20   //普通群管理人数
	VIP_TEAMMANAGE    = 50   //vip群管理人数
	NORMAL_TEAMMERBER = 500  //普通群总人数
	VIP_TEAMMEMBER    = 2000 //vip群总人数
	TEAM_OWNER        = 1    //群主
	TEAM_MANAGER      = 2    //管理员
	TEAM_MASSES       = 3    //群众
	TEAM_UNUSE        = 4    //不在群中
)

const REDIS_TEAMPERSONAL = "team_teampersonal"
const REDIS_TEAMPERSONAL_EXIST_LIST = "team_teampersonal_exist_list"
const REDIS_TEAMPERSONAL_EXIST_TIME = 1000 * 600 //毫秒，key值存在时间

type RedisTeamPersonalObj struct {
	Id int64 //群id
	RedisBase
}

func NewRedisTeamPersonalObj(id int64, data ...[]*share_message.PersonalTeamData) *RedisTeamPersonalObj {
	//logs.Info("初始化群成员:", id)
	p := &RedisTeamPersonalObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	p.Init(obj)
	return p
}
func (self *RedisTeamPersonalObj) Init(obj []*share_message.PersonalTeamData) *RedisTeamPersonalObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_TEAMMEMBER)
	self.Sid = TeamPersonalMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		TeamPersonalMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = self.QueryTeamPersonal()
			//if obj == nil {
			//	return nil
			//}
		}
		self.SetRedisTeamPersonal(obj)
	}
	//logs.Info("初始化新的TeamPersonal管理器:", self.Id)
	return self
}
func (self *RedisTeamPersonalObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisTeamPersonalObj) GetKeyId() string { //override
	return MakeRedisKey(REDIS_TEAMPERSONAL, self.Id)
}

//重写保存方法
func (self *RedisTeamPersonalObj) SaveToMongo() {
	members := self.GetRedisTeamPersonal()
	if len(members) > 0 {
		saveList := make([]*share_message.PersonalTeamData, 0)
		for _, m := range members {
			if m.GetIsSave() {
				m.IsSave = easygo.NewBool(false)
				saveList = append(saveList, m)
				self.UpdateTeamMemberInfo(m)
			}
		}
		if len(saveList) > 0 {
			var data []interface{}
			for _, v := range saveList {
				b1 := bson.M{"_id": v.GetId()}
				b2 := v
				data = append(data, b1, b2)
			}
			UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_TEAMMEMBER, data)
		}
	}
	self.SetSaveStatus(false)
}

//定时更新数据
func (self *RedisTeamPersonalObj) UpdateData() { //override
	if !self.IsExistKey() {
		TeamPersonalMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > REDIS_TEAMPERSONAL_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		TeamPersonalMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}

func (self *RedisTeamPersonalObj) InitRedis() { //override
	obj := self.QueryTeamPersonal()
	if obj == nil {
		return
	}
	self.SetRedisTeamPersonal(obj)
}
func (self *RedisTeamPersonalObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisTeamPersonal()
	return data
}
func (self *RedisTeamPersonalObj) SaveOtherData() { //override
	//保存表情数据
}

//通过playerId从mongo中读取登录玩家数据
func (self *RedisTeamPersonalObj) QueryTeamPersonal() []*share_message.PersonalTeamData {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAMMEMBER)
	defer closeFun()
	var members []*share_message.PersonalTeamData
	queryBson := bson.M{"TeamId": self.Id}
	err := col.Find(queryBson).All(&members)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return members
}

//初始化用
func (self *RedisTeamPersonalObj) SetRedisTeamPersonal(members []*share_message.PersonalTeamData) {
	TeamPersonalMgr.Store(self.Id, self)
	self.AddToExistList(self.Id)
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	teamInfo := make(map[int64]string)
	for _, info := range members {
		s, _ := json.Marshal(info)
		teamInfo[info.GetPlayerId()] = string(s)
	}
	if len(teamInfo) == 0 {
		return
	}
	err2 := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), teamInfo)
	easygo.PanicError(err2)

}

func (self *RedisTeamPersonalObj) GetRedisTeamPersonal() []*share_message.PersonalTeamData {
	members := []*share_message.PersonalTeamData{}
	m, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(self.GetKeyId()))
	easygo.PanicError(err)
	for _, s := range m {
		var obj *share_message.PersonalTeamData
		json.Unmarshal([]byte(s), &obj)
		members = append(members, obj)
	}
	return members
}
func GetRedisTeamPersonalObj(id int64, data ...[]*share_message.PersonalTeamData) *RedisTeamPersonalObj {
	return TeamPersonalMgr.GetRedisTeamPersonalObj(id, data...)
}

func (self *RedisTeamPersonalObj) ExistTeamMember(playerId int64) bool {
	b, err := easygo.RedisMgr.GetC().HExists(self.GetKeyId(), easygo.AnytoA(playerId))
	easygo.PanicError(err)
	return b
}

//增加群成员
func (self *RedisTeamPersonalObj) AddRedisTeamPersonal(members []*share_message.PersonalTeamData) {
	teamInfo := make(map[int64]string)
	for _, info := range members {
		s, _ := json.Marshal(info)
		teamInfo[info.GetPlayerId()] = string(s)
	}
	if len(teamInfo) == 0 {
		return
	}
	err2 := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), teamInfo)
	easygo.PanicError(err2)
	self.SetSaveSid()
	self.SetSaveStatus(true)
}
func GetDefaultSetting() *share_message.PersonalTeamSetting {
	setting := &share_message.PersonalTeamSetting{
		IsTopChat:   easygo.NewBool(false),
		IsNoDisturb: easygo.NewBool(false),
		IsSaveAdd:   easygo.NewBool(false),
	}
	return setting
}

//增加群成员
func AddTeamPersonData(teamId, maxId int64, plst []int64, pos int32, reason string, msg *share_message.TeamChannel, id ...int64) bool {
	logs.Info("==========AddTeamPersonData=================", teamId, plst, pos, maxId)
	obj := TeamPersonalMgr.GetRedisTeamPersonalObj(teamId)
	owner := append(id, 0)[0]
	ti := GetMillSecond()
	st := base64.StdEncoding.EncodeToString([]byte(reason))
	var teamMembers []*share_message.PersonalTeamData
	//chatObj := GetRedisTeamChatLog(teamId)
	//maxId := chatObj.GetTeamMaxChatLogId()
	for _, pid := range plst {
		id := NextId(TABLE_TEAMMEMBER)
		player := GetRedisPlayerBase(pid)
		if player == nil {
			continue
		}
		player.AddPlayerSession(easygo.AnytoA(teamId))
		msg := &share_message.PersonalTeamData{
			Id:       easygo.NewInt64(id),
			TeamId:   easygo.NewInt64(teamId),
			PlayerId: easygo.NewInt64(pid),
			//NickName:    easygo.NewString(player.GetNickName()),
			NickName:    easygo.NewString(""), //群昵称
			Setting:     GetDefaultSetting(),
			ReadId:      easygo.NewInt64(maxId),
			Time:        easygo.NewInt64(ti),
			Channel:     easygo.NewString(st),
			TeamChannel: msg,
			Status:      easygo.NewInt32(1),
			IsSave:      easygo.NewBool(true),
			TeamName:    easygo.NewString(""),
		}
		if pid == owner {
			msg.Position = easygo.NewInt32(1) //群主
			msg.Channel = easygo.NewString("")
		} else {
			msg.Position = easygo.NewInt32(pos)
		}
		teamMembers = append(teamMembers, msg)
	}
	obj.AddRedisTeamPersonal(teamMembers)
	obj.SaveToMongo()
	return true
}

func GetTeamChannel(s string) (bool, string) {
	by, _ := base64.StdEncoding.DecodeString(s)
	text := string(by)
	var channel string
	var b bool
	if len(text) >= 6 && text[:6] == "客服" { //兼容代码
		text = "群主邀请入群"
		channel = base64.StdEncoding.EncodeToString([]byte(text))
		b = true
	} else {
		channel = s
	}
	return b, channel
}
func (self *RedisTeamPersonalObj) GetLoginTeamMember(ids []int64, playerId int64) ([]*share_message.PersonalTeamData, *share_message.PersonalTeamData) {
	if !self.IsExistKey() {
		self.InitRedis()
	}
	members := make([]*share_message.PersonalTeamData, 0)
	myData := &share_message.PersonalTeamData{}
	m := self.GetRedisTeamPersonal()
	for _, obj := range m {
		b, channel := GetTeamChannel(obj.GetChannel())
		if b {
			obj.Channel = easygo.NewString(channel)
		}
		if easygo.Contain(ids, obj.GetPlayerId()) {
			members = append(members, obj)
		}
		if playerId == obj.GetPlayerId() {
			myData = obj
		}
	}
	//处理刷新玩家信息到群成员
	mpIds := append(ids, playerId)
	pMap := GetAllPlayerBase(mpIds, false)
	for _, m := range members {
		m = GetOneTeamMember(m, pMap[m.GetPlayerId()])
	}
	myData = GetOneTeamMember(myData, pMap[playerId])
	return members, myData
}
func GetAllTeamMember(teamId int64, playerId ...PLAYER_ID) []*share_message.PersonalTeamData {
	pid := append(playerId, 0)[0]
	members := make([]*share_message.PersonalTeamData, 0)
	memberObj := GetRedisTeamPersonalObj(teamId)
	m := memberObj.GetRedisTeamPersonal()
	teamObj := GetRedisTeamObj(teamId)
	pList := teamObj.GetTeamMemberList()
	//处理刷新玩家信息到群成员
	pMap := GetAllPlayerBase(pList, false)
	for _, obj := range m {
		//var obj *share_message.PersonalTeamData
		//json.Unmarshal([]byte(s), &obj)
		b, channel := GetTeamChannel(obj.GetChannel())
		if b {
			obj.Channel = easygo.NewString(channel)
		}
		//logs.Info("data:", pMap[obj.GetPlayerId()])
		obj.OperatorInfoPer = memberObj.GetRedisOperatorInfoPer(obj.GetPlayerId(), pMap[obj.GetPlayerId()])
		if pid == obj.GetPlayerId() {
			//检查是否是自己被封禁了
			if infos := obj.GetOperatorInfoPer(); infos != nil {
				for _, info := range infos {
					/*	if info.GetCloseTime() <= util.GetMilliTime() && obj.GetStatus() != 1 { //已过禁言时间
						obj.Status = easygo.NewInt32(1)
						obj.IsSave = easygo.NewBool(true)
						memberObj.UpdateTeamMemberInfo(obj)
						memberObj.DelRedisOperatorInfoPer(info.GetFlag(), obj.GetPlayerId())
					}*/

					if info.GetCloseTime() > util.GetMilliTime() || obj.GetStatus() == 1 { //已过禁言时间
						continue
					}
					obj.Status = easygo.NewInt32(1)
					obj.IsSave = easygo.NewBool(true)
					memberObj.UpdateTeamMemberInfo(obj)
					memberObj.DelRedisOperatorInfoPer(info.GetFlag(), obj.GetPlayerId())

				}
			}
		}
		//obj = GetOneTeamMember(obj)
		members = append(members, obj)
	}
	for _, m := range members {
		m = GetOneTeamMember(m, pMap[m.GetPlayerId()])
	}
	return members
}

func (self *RedisTeamPersonalObj) GetTeamMember(playerId int64, player ...*share_message.PlayerBase) *share_message.PersonalTeamData {
	p := append(player, nil)[0]
	member := &share_message.PersonalTeamData{}
	b := self.ExistTeamMember(playerId)
	if !b {
		//玩家已经被踢出群
		member = GetPerTeamDataForMongoDB(self.Id, playerId)
		if member != nil {
			return member
		}
		logs.Error("玩家不在群里", playerId, self.Id)
		//清理成员残留数据
		base := GetRedisPlayerBase(playerId)
		if base != nil {
			base.DelTeamId(self.Id)
		}
		teamObj := GetRedisTeamObj(self.Id)
		if teamObj != nil {
			teamObj.DelTeamMemberList([]int64{playerId})
		}
		return nil
	}
	val, err := easygo.RedisMgr.GetC().HGet(self.GetKeyId(), easygo.AnytoA(playerId))
	if err != nil {
		logs.Error("玩家信息获取不到,去数据库查询")
		member = GetPerTeamDataForMongoDB(self.Id, playerId)
	} else {
		_ = json.Unmarshal(val, &member)
	}
	b, channel := GetTeamChannel(member.GetChannel())
	if b {
		member.Channel = easygo.NewString(channel)
	}
	if p == nil {
		return member
	}
	member = GetOneTeamMember(member, p)
	//member.OperatorInfoPer = self.GetRedisOperatorInfoPer(member.GetPlayerId())
	return member
}

func GetTeamPlayerPos(teamId, playerId int64) int32 {
	self := GetRedisTeamPersonalObj(teamId)
	if !self.ExistTeamMember(playerId) {
		return TEAM_UNUSE
	}
	member := self.GetTeamMember(playerId)
	return member.GetPosition()
}

func SetTeamPlayerPos(teamId, playerId int64, pos int32) {
	self := GetRedisTeamPersonalObj(teamId)
	if !self.ExistTeamMember(playerId) {
		logs.Error("群成员不存在")
		return
	}
	member := self.GetTeamMember(playerId)
	member.Position = easygo.NewInt32(pos)
	member.IsSave = easygo.NewBool(true)
	self.UpdateTeamMemberInfo(member)
	//设置职位后，马上存储到mongo
	self.SaveToMongo()
}

func (self *RedisTeamPersonalObj) GetTeamMemberReName(playerId int64) string {
	if !self.ExistTeamMember(playerId) {
		return ""
	}
	member := self.GetTeamMember(playerId)
	return member.GetNickName()
}

func (self *RedisTeamPersonalObj) GetTeamMemberChannel(playerId int64) string {
	if !self.ExistTeamMember(playerId) {
		return ""
	}
	member := self.GetTeamMember(playerId)
	return member.GetChannel()
}
func (self *RedisTeamPersonalObj) GetTeamPosition(playerId int64) int32 {
	if !self.ExistTeamMember(playerId) {
		return TEAM_UNUSE
	}
	member := self.GetTeamMember(playerId)
	return member.GetPosition()
}

//获取群成员
func GetOneTeamMember(msg *share_message.PersonalTeamData, p ...*share_message.PlayerBase) *share_message.PersonalTeamData {
	base := append(p, nil)[0]
	if base == nil {
		return msg
	}

	photo := ""
	photo_list := base.GetPhoto()
	if len(photo_list) != 0 {
		photo = photo_list[0]
	}
	name := msg.GetNickName()
	//if name == "" {
	//	name = base.GetNickName()
	//}
	msg.NickName = easygo.NewString(name)
	msg.HeadIcon = easygo.NewString(base.GetHeadIcon())
	msg.Sex = easygo.NewInt32(base.GetSex())
	msg.Account = easygo.NewString(base.GetAccount())
	msg.Signture = easygo.NewString(base.GetSignature())
	msg.Photo = easygo.NewString(photo)
	msg.PerNickName = easygo.NewString(base.GetNickName())
	msg.Types = easygo.NewInt32(base.GetTypes())
	return msg
}

func (self *RedisTeamPersonalObj) GetTeamReadChatLogId(playerId int64) int64 {
	if !self.ExistTeamMember(playerId) {
		return 0
	}
	member := self.GetTeamMember(playerId)
	return member.GetReadId()
}

func (self *RedisTeamPersonalObj) ReadTeamChatLog(playerId, logId int64) {
	if !self.ExistTeamMember(playerId) {
		return
	}
	member := self.GetTeamMember(playerId)
	member.ReadId = easygo.NewInt64(logId)
	member.IsSave = easygo.NewBool(true)
	self.UpdateTeamMemberInfo(member)
	//logs.Info("member:", playerId, member.GetReadId())
}
func (self *RedisTeamPersonalObj) GetTeamQTX(playerId int64) bool {
	if !self.ExistTeamMember(playerId) {
		return false
	}
	member := self.GetTeamMember(playerId)
	return member.GetQTX()
}
func (self *RedisTeamPersonalObj) SetTeamQTX(playerId int64, b bool) {
	if !self.ExistTeamMember(playerId) {
		return
	}
	member := self.GetTeamMember(playerId)
	member.QTX = easygo.NewBool(b)
	self.UpdateTeamMemberInfo(member)
}
func (self *RedisTeamPersonalObj) DelTeamPersonData(playerIds []int64) {
	for _, pid := range playerIds {
		if !self.ExistTeamMember(pid) {
			continue
		}
		m := self.GetTeamMember(pid)
		col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAMMEMBER)
		defer closeFun()
		if m.GetId() > 0 {
			count, err := col.Find(bson.M{"_id": m.GetId()}).Count()
			easygo.PanicError(err)
			if count > 0 {
				err := col.RemoveId(m.GetId()) //数据库删除
				easygo.PanicError(err)
			}
		}
		b, err1 := easygo.RedisMgr.GetC().Hdel(self.GetKeyId(), easygo.AnytoA(pid))
		easygo.PanicError(err1)
		if b {
			teamObj := GetRedisTeamObj(self.Id)
			teamObj.UpsertTeamOutTeamInfo(pid)
		}
	}
}

func (self *RedisTeamPersonalObj) GetTeamPersonSetting(playerId int64) *share_message.PersonalTeamSetting {
	if !self.ExistTeamMember(playerId) {
		return nil
	}
	member := self.GetTeamMember(playerId)
	return member.GetSetting()
}

func (self *RedisTeamPersonalObj) UpdateTeamPersonalSetting(pid int64, t int32, value interface{}) {
	member := self.GetTeamMember(pid)
	switch t {
	case 1:
		member.Setting.IsTopChat = easygo.NewBool(value.(bool))
	case 2:
		member.Setting.IsNoDisturb = easygo.NewBool(value.(bool))
	case 3:
		member.Setting.IsSaveAdd = easygo.NewBool(value.(bool))
	case 4:
		member.NickName = easygo.NewString(value.(string))
	case 5:
		member.TeamName = easygo.NewString(value.(string))
	}
	member.IsSave = easygo.NewBool(true)
	self.UpdateTeamMemberInfo(member)
	if t == 3 {
		self.SaveToMongo()
	}
}

//更新信息进redis
func (self *RedisTeamPersonalObj) UpdateTeamMemberInfo(member *share_message.PersonalTeamData) {
	//member.IsSave = easygo.NewBool(true)
	s, _ := json.Marshal(member)
	err := easygo.RedisMgr.GetC().HSet(self.GetKeyId(), easygo.AnytoA(member.GetPlayerId()), string(s))
	easygo.PanicError(err)
	self.SetSaveSid()
	self.SetSaveStatus(true)
}

//更新封禁人信息
func (self *RedisTeamPersonalObj) SetRedisOperatorInfoPer(playerId int64, obj *share_message.OperatorInfoPer) {
	member := self.GetTeamMember(playerId)
	member.OperatorInfoPer = append(member.OperatorInfoPer, obj)
	js, err := json.Marshal(member)
	easygo.PanicError(err)
	err1 := easygo.RedisMgr.GetC().HSet(self.GetKeyId(), easygo.AnytoA(playerId), string(js))
	easygo.PanicError(err1)
}

//获取封禁人信息
func (self *RedisTeamPersonalObj) GetRedisOperatorInfoPer(playerId int64, player ...*share_message.PlayerBase) []*share_message.OperatorInfoPer {
	p := append(player, nil)[0]
	teamObj := GetRedisTeamObj(self.Id)
	m := teamObj.GetRedisOperatorInfo()
	member := self.GetTeamMember(playerId, p)
	inviteList := member.GetOperatorInfoPer()
	for _, s := range m {
		var obj *share_message.OperatorInfoPer
		StructToOtherStruct(s, &obj)
		inviteList = append(inviteList, obj)
	}
	return inviteList
}

//删除封禁人信息
func (self *RedisTeamPersonalObj) DelRedisOperatorInfoPer(flag int32, playerId int64) bool {

	member := self.GetTeamMember(playerId)
	member.OperatorInfoPer = nil
	js, err := json.Marshal(member)
	easygo.PanicError(err)
	err1 := easygo.RedisMgr.GetC().HSet(self.GetKeyId(), easygo.AnytoA(playerId), string(js))
	easygo.PanicError(err1)
	return true
}
func SaveRedisTeamPersonalToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_TEAMMEMBER, &ids)
	memberList := make([]*share_message.PersonalTeamData, 0)
	for _, id := range ids {
		obj := GetRedisTeamPersonalObj(id)
		if obj != nil {
			members := obj.GetRedisTeamPersonal()
			memberList = append(memberList, members...)
			obj.SetSaveStatus(false)
		}
	}
	if len(memberList) > 0 {
		saveData := make([]interface{}, 0)
		for _, m := range memberList {
			saveData = append(saveData, bson.M{"_id": m.GetId()}, m)
		}
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_TEAMMEMBER, saveData)
	}
}
