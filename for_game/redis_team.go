package for_game

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/pb/client_hall"
	"game_server/pb/client_server"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"

	"github.com/garyburd/redigo/redis"
)

/*
群reids内存数据管理
*/

const (
	INVITE_UNTREATED = 1 //未处理
	INVITE_ACCEPT    = 2 //同意
	INVITE_TIMEOUT   = 3 //超时
	INVITE_REFUSE    = 4 //拒绝
)

const (
	REDIS_TEAM_MEMBER       = "MemberList"     //群成员列表
	REDIS_TEAM_MANAGE       = "ManageList"     //群管理员列表
	REDIS_TEAM_OUTTEAM      = "OutTeamInfo"    //退群玩家列表
	REDIS_TEAM_INVITE       = "InviteInfo"     //群邀请
	REDIS_TEAM_OPERATORINFO = "OperatorInfo"   //群操作
	REDIS_TEAM_SETTING      = "MessageSetting" //群设置
	REDIS_TEAMCHAT_MAXLOGID = "MaxLogId"       //最大logId
	REDIS_TEAM_EXIST_LIST   = "team_exist_list"
	REDIS_TEAM_EXIST_TIME   = 1000 * 600 //毫秒，key值存在时间
)

type RedisTeamObj struct {
	Id int64 //群id
	RedisBase
}
type RedisTeam struct {
	TeamId       int64  `json:"_id"` //群id
	Name         string //群名称
	HeadUrl      string //群头像
	GongGao      string //群公告
	Owner        int64  //群主ID
	QRCode       string //收款二维码链接
	CreateTime   int64  //创建时间
	LastTalkTime int64  //最后聊天时间
	MaxMember    int32  //最大人数
	TeamChat     string //群聊号
	//MessageSetting string //群管理设置
	Status          int32  //群状态 0正常，1解散
	IsRecommend     bool   //是否推荐
	DissolveTime    int64  //解散时间
	RefreshTime     int64  //刷新群收款二维码时间
	AdminId         int64  //后台建群管理员ID
	CreateName      string //创建群的人的名字
	Level           int32  //群等级，越高级的群权限越多
	LogMaxId        int64  //群聊天id
	WelcomeWord     string //欢迎语
	Topic           string //话题
	TopicDesc       string //话题简介
	DynamicId       int64  //上次发布一条动态的id
	LastDynamicTime int64  //上一次发布动态时间
}

func NewRedisTeamObj(Id int64, data ...*share_message.TeamData) *RedisTeamObj {
	p := &RedisTeamObj{
		Id: Id,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}
func (self *RedisTeamObj) Init(obj *share_message.TeamData) *RedisTeamObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_TEAM_DATA)
	self.Sid = TeamMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		TeamMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		self.SetSaveStatus(true)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = self.QueryTeam()
			if obj == nil {
				return nil
			}
		}
		self.SetRedisTeam(obj)
	}

	//	logs.Info("初始化新的Team管理器:", self.Id)
	return self
}
func (self *RedisTeamObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisTeamObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_TEAM_DATA, self.Id)
}

//定时更新数据
func (self *RedisTeamObj) UpdateData() { //override
	if !self.IsExistKey() {
		TeamMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > REDIS_TEAM_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		TeamMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}

func (self *RedisTeamObj) InitRedis() { //override
	obj := self.QueryTeam()
	if obj == nil {
		return
	}
	self.SetRedisTeam(obj)
}
func (self *RedisTeamObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisTeam()
	return data
}
func (self *RedisTeamObj) SaveOtherData() { //override
	//保存表情数据
}

//查询当前最大的logid
func (self *RedisTeamObj) GetMongoMaxLogId() int64 {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_TEAM_CHAT_LOG)
	defer closeFun()
	var log *share_message.TeamChatLog
	err := col.Find(bson.M{"TeamId": self.Id}).Sort("-TeamLogId").One(&log)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return 0
	}
	return log.GetTeamLogId()
}

//通过playerId从mongo中读取登录玩家数据
func (self *RedisTeamObj) QueryTeam() *share_message.TeamData {
	data := self.QueryMongoData(self.Id)
	if data != nil {
		var team share_message.TeamData
		StructToOtherStruct(data, &team)
		if team.GetLogMaxId() == 0 {
			team.LogMaxId = easygo.NewInt64(self.GetMongoMaxLogId())
		}
		return &team
	}
	return nil
}

func (self *RedisTeamObj) SetRedisTeam(obj *share_message.TeamData) {
	//logs.Info("初始化群id:", self.Id)
	TeamMgr.Store(obj.GetId(), self)
	self.AddToExistList(obj.GetId())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		//logs.Info("群数据已经存在")
		return
	}
	//玩家基础信息
	team := &RedisTeam{}
	StructToOtherStruct(obj, team)
	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), team)
	easygo.PanicError(err)
	//其他群数据设置到redis
	self.SetTeamOtherData(obj, false)
	//初始化聊天obj
	_ = TeamChatLogMgr.GetRedisTeamChatLogObj(self.Id)
	//logs.Info("chatObj:", chatObj)
	//if chatObj == nil {
	//	self.UpdateTeamChatLogMaxId(0) //新群从1开始
	//}
	//logs.Info("群初始化完毕:", self.Id)
}

//群组其他数据写入redis
func (self *RedisTeamObj) SetTeamOtherData(obj *share_message.TeamData, save ...bool) {
	self.SetStringValueToRedis(REDIS_TEAM_MEMBER, obj.GetMemberList(), save...)
	self.SetStringValueToRedis(REDIS_TEAM_MANAGE, obj.GetManagerList(), save...)
	self.SetStringValueToRedis(REDIS_TEAM_OUTTEAM, obj.GetOutTeam(), save...)
	self.SetStringValueToRedis(REDIS_TEAM_INVITE, obj.GetInvite(), save...)
	self.SetStringValueToRedis(REDIS_TEAM_OPERATORINFO, obj.GetOperatorInfo(), save...)
	self.SetStringValueToRedis(REDIS_TEAM_SETTING, obj.GetMessageSetting(), save...)
}

//获取群组其他数据
func (self *RedisTeamObj) GetTeamOtherData(obj *share_message.TeamData) {
	self.GetStringValueToRedis(REDIS_TEAM_MEMBER, &obj.MemberList)
	self.GetStringValueToRedis(REDIS_TEAM_MANAGE, &obj.ManagerList)
	self.GetStringValueToRedis(REDIS_TEAM_OUTTEAM, &obj.OutTeam)
	self.GetStringValueToRedis(REDIS_TEAM_INVITE, &obj.Invite)
	self.GetStringValueToRedis(REDIS_TEAM_OPERATORINFO, &obj.OperatorInfo)
	self.GetStringValueToRedis(REDIS_TEAM_SETTING, &obj.MessageSetting)
}

//获取群信息
func (self *RedisTeamObj) GetRedisTeam() *share_message.TeamData {
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	easygo.PanicError(err)
	var team RedisTeam
	err = redis.ScanStruct(value, &team)
	easygo.PanicError(err)
	newTeam := &share_message.TeamData{}
	StructToOtherStruct(team, newTeam)
	//其他数据
	self.GetTeamOtherData(newTeam)
	if newTeam.GetHeadUrl() == "" {
		owner := GetRedisPlayerBase(self.GetTeamOwner())
		if owner != nil {
			newTeam.HeadUrl = easygo.NewString(owner.GetHeadIcon())
			self.SetHeadUrl(owner.GetHeadIcon())
		}
	}
	return newTeam
}

//对外方法
func GetRedisTeamObj(id int64, team ...*share_message.TeamData) *RedisTeamObj {
	return TeamMgr.GetRedisTeamObj(id, team...)
}

//群管理员处理
func (self *RedisTeamObj) GetTeamManageList() []int64 {
	var manageList []int64
	self.GetStringValueToRedis(REDIS_TEAM_MANAGE, &manageList)
	return manageList
}

//增加管理员
func (self *RedisTeamObj) AddTeamManage(playerIds []int64) {
	list := self.GetTeamManageList()
	list = append(list, playerIds...)
	self.SetStringValueToRedis(REDIS_TEAM_MANAGE, list)
}

func (self *RedisTeamObj) DelTeamManageList(plst []int64) {
	list := self.GetTeamManageList()
	for _, p := range plst {
		if easygo.Contain(list, p) {
			list = easygo.Del(list, p).([]int64)
		}
	}
	self.SetStringValueToRedis(REDIS_TEAM_MANAGE, list)
}

//群成员列表
func (self *RedisTeamObj) GetTeamMemberList() []int64 {
	var memberList []int64
	self.GetStringValueToRedis(REDIS_TEAM_MEMBER, &memberList)
	return memberList
}

//增加群成员
func (self *RedisTeamObj) AddTeamMember(playerIds []int64) {
	list := self.GetTeamMemberList()
	list = append(list, playerIds...)
	self.SetStringValueToRedis(REDIS_TEAM_MEMBER, list)
	session := GetRedisChatSessionObj(easygo.AnytoA(self.Id))
	session.AddRedisChatSessionPlayers(playerIds)
}

func (self *RedisTeamObj) DelTeamMemberList(plst []int64) {
	list := self.GetTeamMemberList()
	for _, p := range plst {
		if easygo.Contain(list, p) {
			list = easygo.Del(list, p).([]int64)
		}
	}
	self.SetStringValueToRedis(REDIS_TEAM_MEMBER, list)
	session := GetRedisChatSessionObj(easygo.AnytoA(self.Id))
	session.DelRedisChatSessionPlayers(plst)
}

//获取群设置信息
func (self *RedisTeamObj) GetTeamMessageSetting() *share_message.MessageSetting {
	var setting *share_message.MessageSetting
	self.GetStringValueToRedis(REDIS_TEAM_SETTING, &setting)
	return setting
}

//获取群人数上限
func (self *RedisTeamObj) GetTeamMaxMember() int32 {
	var val int32
	self.GetOneValue("MaxMember", &val)
	return val
}

//获取群主id
func (self *RedisTeamObj) GetTeamOwner() int64 {
	var val int64
	self.GetOneValue("Owner", &val)
	return val
}

//获取最近一次群发布的动态id
func (self *RedisTeamObj) GetDynamicId() int64 {
	var val int64
	self.GetOneValue("DynamicId", &val)
	return val
}

//设置最近群发布的动态id
func (self *RedisTeamObj) SetDynamicId(id int64) {
	self.SetOneValue("DynamicId", id)
	self.SaveOneRedisDataToMongo("DynamicId", id)
}

//获取最近一次群发布的动态id
func (self *RedisTeamObj) GetLastDynamicTime() int64 {
	var val int64
	self.GetOneValue("LastDynamicTime", &val)
	return val
}

//设置最近群发布的动态id
func (self *RedisTeamObj) SetLastDynamicTime(t int64) {
	self.SetOneValue("LastDynamicTime", t)
	self.SaveOneRedisDataToMongo("LastDynamicTime", t)
}

//获取群主id
func (self *RedisTeamObj) GetTeamName() string {
	var val string
	self.GetOneValue("Name", &val)
	return val
}

//获取群号
func (self *RedisTeamObj) GetTeamChat() string {
	var val string
	self.GetOneValue("TeamChat", &val)
	return val
}

// 获取群头像
func (self *RedisTeamObj) GetHeadUrl() string {
	var val string
	self.GetOneValue("HeadUrl", &val)
	return val
}

// 设置群头像
func (self *RedisTeamObj) SetHeadUrl(headUrl string) {
	self.SetOneValue("HeadUrl", headUrl)
	//会话群头像对应修改
	session := GetRedisChatSessionObj(easygo.AnytoA(self.Id))
	session.SetSessionHeadUrl(headUrl)
}

//获取话题
func (self *RedisTeamObj) GetTopic() string {
	var val string
	self.GetOneValue("Topic", &val)
	return val
}

//设置话题
func (self *RedisTeamObj) SetTopic(desc string) {
	self.SetOneValue("Topic", desc)
}

//获取话题简介
func (self *RedisTeamObj) GetTopicDesc() string {
	var val string
	self.GetOneValue("TopicDesc", &val)
	return val
}

//设置话题简介
func (self *RedisTeamObj) SetTopicDesc(desc string) {
	self.SetOneValue("TopicDesc", desc)
}
func (self *RedisTeamObj) GetSessionLogMaxId() int64 {
	session := GetRedisChatSessionObj(easygo.AnytoA(self.Id))
	return session.GetMaxLogId()
}

//获取群收款二维码
func (self *RedisTeamObj) GetTeamQRCode() string {
	var val string
	self.GetOneValue("QRCode", &val)
	return val
}

func (self *RedisTeamObj) GetTeamStatus() int32 {
	var val int32
	self.GetOneValue("Status", &val)
	return val
}

//获取群收款二维码刷新时间
func (self *RedisTeamObj) GetTeamRefreshTime() int64 {
	var val int64
	self.GetOneValue("RefreshTime", &val)
	return val
}

//检查群人数上限
func (self *RedisTeamObj) CheckTeamMaxMember(num int) bool {
	max := self.GetTeamMaxMember()
	memberList := self.GetTeamMemberList()
	if len(memberList)+num > int(max) {
		return false
	}
	return true
}

//检查群管理员人数上限
func (self *RedisTeamObj) CheckTeamMaxManager(num int) bool {
	manageList := self.GetTeamManageList()
	if len(manageList)+num > 10 {
		return false
	}
	return true
}

//设置欢迎语
func (self *RedisTeamObj) SetWelcomeWord(word string) {
	self.SetOneValue("WelcomeWord", word)
}

//获取欢迎语
func (self *RedisTeamObj) GetWelcomeWord() string {
	var val string
	self.GetOneValue("WelcomeWord", &val)
	return val
}

//=================outteam=======================

//获取所有退群信息
func (self *RedisTeamObj) GetTeamOutTeamInfo() []*share_message.OutTeamInfo {
	var val []*share_message.OutTeamInfo
	self.GetStringValueToRedis(REDIS_TEAM_OUTTEAM, &val)
	return val
}

//更新退群信息
func (self *RedisTeamObj) UpsertTeamOutTeamInfo(playerId int64) {
	outList := self.GetTeamOutTeamInfo()
	isExist := false
	for _, v := range outList {
		if v.GetPlayerId() == playerId {
			v.Time = easygo.NewInt64(GetMillSecond())
			isExist = true
			break
		}
	}
	if !isExist {
		msg := &share_message.OutTeamInfo{
			PlayerId: easygo.NewInt64(playerId),
			Time:     easygo.NewInt64(GetMillSecond()),
		}
		outList = append(outList, msg)
	}
	self.SetStringValueToRedis(REDIS_TEAM_OUTTEAM, outList)
}

//=================invite=======================

//获取所有进群请求信息
func (self *RedisTeamObj) GetTeamInviteInfo() []*share_message.InviteInfo {
	var val []*share_message.InviteInfo
	self.GetStringValueToRedis(REDIS_TEAM_INVITE, &val)
	return val
}

//查询进群请求是否有效
func (self *RedisTeamObj) GetTeamVaildInviteState(playerId int64) bool {
	self.CheckTeamInviteTeamInfo()
	for _, info := range self.GetTeamInviteInfo() {
		if info.GetPlayerId() == playerId && info.GetState() == INVITE_UNTREATED {
			return true
		}
	}
	return false
}

func (self *RedisTeamObj) UpdateAllTeamInviteInfo(newList []*share_message.InviteInfo) {
	self.GetStringValueToRedis(REDIS_TEAM_INVITE, newList)
}

//检查进群请求信息
func (self *RedisTeamObj) CheckTeamInviteTeamInfo() {
	inviteInfo := self.GetTeamInviteInfo()
	t := time.Now().Unix()
	var b, b1 bool
	var newList []*share_message.InviteInfo
	for _, info := range inviteInfo { //如果申请信息超过7天就删除
		if t-info.GetTime() > 7*86400 {
			b = true
			continue
		}
		newList = append(newList, info)
	}

	for _, info := range newList { //如果申请信息超过3天就设置状态为过期
		if info.GetState() == INVITE_UNTREATED && t-info.GetTime() > 3*86400 {
			info.State = easygo.NewInt32(INVITE_TIMEOUT)
			b1 = true
		}
	}
	if b || b1 {
		self.SetStringValueToRedis(REDIS_TEAM_INVITE, newList)
	}
}

//增加加群请求
func (self *RedisTeamObj) AddTeamInvite(pid int64, inviteName string, plst []int64, reason string, msg *share_message.TeamChannel) {
	self.CheckTeamInviteTeamInfo()
	var newlst []*share_message.InviteInfo
	lst := self.GetTeamInviteInfo()
	for _, info := range lst { //获取所有未通过的请求
		if info.GetState() == INVITE_UNTREATED {
			newlst = append(newlst, info)
		}
	}

	var isSave bool
	for _, id := range plst {
		var b bool
		for _, info := range newlst {
			if info.GetPlayerId() == id {
				b = true
				break
			}
		}
		if b {
			continue
		}
		info := GetRedisPlayerBase(id)
		var logId int64
		if len(lst) == 0 {
			logId = 1
		} else {
			logId = lst[len(lst)-1].GetLogId() + 1
		}
		msg := &share_message.InviteInfo{
			PlayerId:    easygo.NewInt64(id),
			HeadIcon:    easygo.NewString(info.GetHeadIcon()),
			State:       easygo.NewInt32(INVITE_UNTREATED),
			Channel:     easygo.NewString(reason),
			InviteName:  easygo.NewString(inviteName),
			InviteId:    easygo.NewInt64(pid),
			TeamChannel: msg,
			Name:        easygo.NewString(info.GetNickName()),
			Time:        easygo.NewInt64(time.Now().Unix()),
			Account:     easygo.NewString(info.GetAccount()),
			LogId:       easygo.NewInt64(logId),
		}
		lst = append(lst, msg)
		isSave = true
	}
	if isSave {
		self.SetStringValueToRedis(REDIS_TEAM_INVITE, lst)
	}
}

func (self *RedisTeamObj) GetTeamInviteStateForLogId(logId int64) int32 {
	self.CheckTeamInviteTeamInfo()
	for _, info := range self.GetTeamInviteInfo() {
		if info.GetLogId() == logId {
			return info.GetState()
		}
	}
	return 0
}

func (self *RedisTeamObj) AcceptTeamRequest(logId int64) (string, int64, *share_message.TeamChannel) {
	var reason string
	var pid int64
	var channel *share_message.TeamChannel
	m := &share_message.InviteInfo{}
	inviteList := self.GetTeamInviteInfo()
	for _, msg := range inviteList {
		if msg.GetLogId() == logId {
			msg.State = easygo.NewInt32(INVITE_ACCEPT)
			reason = msg.GetChannel()
			channel = msg.GetTeamChannel()
			pid = msg.GetInviteId()
			m = msg
			break
		}
	}
	if m != nil {
		self.UpdateTeamInviteInfo(m)
	}
	return reason, pid, channel
}

//添加/修改一个邀请
func (self *RedisTeamObj) UpdateTeamInviteInfo(msg *share_message.InviteInfo) {
	inviteList := self.GetTeamInviteInfo()
	isFind := false
	for k, info := range inviteList {
		if info.GetLogId() == msg.GetLogId() {
			inviteList[k] = msg
			isFind = true
			break
		}
	}
	if !isFind {
		inviteList = append(inviteList, msg)
	}
	self.SetStringValueToRedis(REDIS_TEAM_INVITE, inviteList)
}

func (self *RedisTeamObj) UpdateTeamInviteState(logId int64, state int32) {
	inviteList := self.GetTeamInviteInfo()
	for _, info := range inviteList {
		if info.GetLogId() == logId {
			info.State = easygo.NewInt32(state)
			break
		}
	}
	self.SetStringValueToRedis(REDIS_TEAM_INVITE, inviteList)
}

//---------------------------修改参数------------------------------
func (self *RedisTeamObj) SetTeamQRCode(s string) {
	self.SetOneValue("QRCode", s)
	self.SetOneValue("RefreshTime", time.Now().Unix())
}

func (self *RedisTeamObj) SetTeamLastTalkTime(t int64) {
	self.SetOneValue("LastTalkTime", t)
}

func (self *RedisTeamObj) SetTeamOwner(playerId int64) {
	self.SetOneValue("Owner", playerId)
}

func (self *RedisTeamObj) SetTeamIsRecommend(b bool) {
	self.SetOneValue("IsRecommend", b)
}

func (self *RedisTeamObj) SetTeamMaxMember(num int32) {
	self.SetOneValue("MaxMember", num)
}

func (self *RedisTeamObj) SetTeamNikeName(s string) {
	self.SetOneValue("Name", s)
	//会话群昵称对应修改
	session := GetRedisChatSessionObj(easygo.AnytoA(self.Id))
	session.SetTeamName(s)
}

func (self *RedisTeamObj) SetTeamGongGao(s string) {
	self.SetOneValue("GongGao", s)
}

func (self *RedisTeamObj) SetTeamStatus(status int32) {
	self.SetOneValue("Status", status)
	self.SaveOneRedisDataToMongo("Status", status)
}

func (self *RedisTeamObj) SetTeamDissolveTime(t int64) {
	self.SetOneValue("DissolveTime", t)
}

func (self *RedisTeamObj) SetTeamLevel(lv int32) {
	self.SetOneValue("Level", lv)
}

//获取自增后id
func (self *RedisTeamObj) GetNextLogMaxId() int64 {
	return self.IncrOneValue("LogMaxId", 1)
}

//获取当前id
func (self *RedisTeamObj) GetLogMaxId() int64 {
	var val int64
	self.GetOneValue("LogMaxId", &val)
	return val
}

func (self *RedisTeamObj) UpdateTeamMessageSetting(setting *share_message.MessageSetting) {
	self.SetStringValueToRedis(REDIS_TEAM_SETTING, setting)
}

func (self *RedisTeamObj) SetTeamIsOpenTeamMoneyCode(b bool) {
	setting := self.GetTeamMessageSetting()
	setting.IsOpenTeamMoneyCode = easygo.NewBool(b)
	self.UpdateTeamMessageSetting(setting)
}

func (self *RedisTeamObj) SetIsTimeClean(b bool) {
	setting := self.GetTeamMessageSetting()
	setting.IsTimeClean = easygo.NewBool(b)
	self.UpdateTeamMessageSetting(setting)
}

func (self *RedisTeamObj) SetTeamIsReadClean(b bool) {
	setting := self.GetTeamMessageSetting()
	setting.IsReadClean = easygo.NewBool(b)
	self.UpdateTeamMessageSetting(setting)
}

func (self *RedisTeamObj) SetTeamIsScreenShotNotify(b bool) {
	setting := self.GetTeamMessageSetting()
	setting.IsScreenShotNotify = easygo.NewBool(b)
	self.UpdateTeamMessageSetting(setting)
}

func (self *RedisTeamObj) SetTeamIsStopTalk(b bool) {
	setting := self.GetTeamMessageSetting()
	setting.IsStopTalk = easygo.NewBool(b)
	self.UpdateTeamMessageSetting(setting)
}

func (self *RedisTeamObj) SetTeamIsAddFriend(b bool) {
	setting := self.GetTeamMessageSetting()
	setting.IsAddFriend = easygo.NewBool(b)
	self.UpdateTeamMessageSetting(setting)
}

func (self *RedisTeamObj) SetTeamIsInvite(b bool) {
	setting := self.GetTeamMessageSetting()
	setting.IsInvite = easygo.NewBool(b)
	self.UpdateTeamMessageSetting(setting)
}

func (self *RedisTeamObj) SetTeamIsStopAddTeam(b bool) {
	setting := self.GetTeamMessageSetting()
	setting.IsStopAddTeam = easygo.NewBool(b)
	self.UpdateTeamMessageSetting(setting)
}

func (self *RedisTeamObj) SetTeamIsBan(teamId int64, b bool, t int64) {
	setting := self.GetTeamMessageSetting()
	setting.IsBan = easygo.NewBool(b)
	setting.UnBanTime = easygo.NewInt64(t)
	self.UpdateTeamMessageSetting(setting)
}
func (self *RedisTeamObj) SetTeamIsBan2(teamId int64, b bool) {
	setting := self.GetTeamMessageSetting()
	setting.IsBan = easygo.NewBool(b)
	self.UpdateTeamMessageSetting(setting)
}

//---------------------------------------------------------------
func GetTeamMsgForHall(teamId int64) *client_server.TeamMsg {
	team := GetRedisTeamObj(teamId)
	setInfo := &share_message.TeamSetting{
		PersonSetting:  GetDefaultSetting(),
		MessageSetting: team.GetTeamMessageSetting(),
	}
	msg := &client_server.TeamMsg{
		Team:    team.GetRedisTeam(),
		Members: GetAllTeamMember(teamId),
		Setting: setInfo,
	}
	return msg
}

func GetTeamSettingInfo(teamId, playerId int64) *client_server.TeamMsg {
	team := GetRedisTeamObj(teamId)
	memberObj := GetRedisTeamPersonalObj(teamId)
	personSetting := memberObj.GetTeamPersonSetting(playerId)
	setInfo := &share_message.TeamSetting{
		PersonSetting:  personSetting,
		MessageSetting: team.GetTeamMessageSetting(),
	}
	msg := &client_server.TeamMsg{
		Team:    team.GetRedisTeam(),
		Members: GetAllTeamMember(teamId),
		Setting: setInfo,
	}
	return msg
}
func GetTeamSettingInfoEx(teamId, playerId int64) *share_message.TeamSetting {
	team := GetRedisTeamObj(teamId)
	memberObj := GetRedisTeamPersonalObj(teamId)
	personSetting := memberObj.GetTeamPersonSetting(playerId)
	setInfo := &share_message.TeamSetting{
		PersonSetting:  personSetting,
		MessageSetting: team.GetTeamMessageSetting(),
	}
	return setInfo
}

func GetTeamActivityInfo(teamId int64) []*client_server.TeamActivity {
	var msg []*client_server.TeamActivity
	team := GetRedisTeamObj(teamId)
	nowTime := time.Now().Unix()
	info := make(map[int][]int64)
	for _, pid := range team.GetTeamMemberList() {
		base := GetRedisPlayerBase(pid)
		lastTime := base.GetLastOnLineTime()
		if nowTime-lastTime >= 30*24*3600 { //一个月不活跃
			lst := info[30]
			lst = append(lst, pid)
			info[30] = lst
		} else if nowTime-lastTime >= 7*24*3600 { //一周不活跃
			lst := info[7]
			lst = append(lst, pid)
			info[7] = lst
		} else if nowTime-lastTime >= 3*24*3600 { //三天不活跃
			lst := info[3]
			lst = append(lst, pid)
			info[3] = lst
		}
	}
	for day, lst := range info {
		m := &client_server.TeamActivity{
			Day:      easygo.NewInt32(day),
			PlayerId: lst,
		}
		msg = append(msg, m)
	}
	return msg
}

func GetTeamManageSetting(teamId int64) *client_server.TeamManagerSetting {
	//activityInfo := GetTeamActivityInfo(teamId)todo 活跃列表
	//OutPlayerInfo := GetTeamOutPlayerInfo(teamId) todo 退群列表
	team := GetRedisTeamObj(teamId)
	msg := &client_server.TeamManagerSetting{
		ManageList:     team.GetTeamManageList(),
		MessageSetting: team.GetTeamMessageSetting(),
		//ActivityInfo:   activityInfo,
		//OutPlayerInfo:  OutPlayerInfo,
		InviteInfo: team.GetTeamInviteInfo(),
		TeamId:     easygo.NewInt64(teamId),
		Members:    GetAllTeamMember(teamId),
	}
	return msg
}
func GetTeamManageSettingEx(teamId int64) *client_server.TeamManagerSetting {
	//activityInfo := GetTeamActivityInfo(teamId)todo 活跃列表
	//OutPlayerInfo := GetTeamOutPlayerInfo(teamId) todo 退群列表
	team := GetRedisTeamObj(teamId)
	msg := &client_server.TeamManagerSetting{
		ManageList:     team.GetTeamManageList(),
		MessageSetting: team.GetTeamMessageSetting(),
		//ActivityInfo:   activityInfo,
		//OutPlayerInfo:  OutPlayerInfo,
		InviteInfo: team.GetTeamInviteInfo(),
		//TeamId:     easygo.NewInt64(teamId),
		//Members:    GetAllTeamMember(teamId),
	}
	return msg
}

func GetTeamAllMsg(teamId int64, playerObj *RedisPlayerBaseObj) *client_server.TeamMsg { //登录获取群数据
	team := GetRedisTeamObj(teamId)
	teamData := team.GetRedisTeam()
	ids := []int64{}
	ids = append(ids, teamData.GetOwner())
	if teamData.GetManagerList() != nil {
		ids = append(ids, teamData.GetManagerList()...)
	}
	if len(ids) > 15 {
		ids = ids[:15]
	} else {
		members := teamData.GetMemberList()
		//剔除管理员和群主
		if members != nil {
			for _, id := range members {
				if len(ids) >= 15 {
					break
				}
				if easygo.Contain(ids, id) {
					continue
				}
				ids = append(ids, id)
			}
		}
	}
	memberObj := GetRedisTeamPersonalObj(teamId)
	//logs.Info("memberObj:", memberObj)
	members, myData := memberObj.GetLoginTeamMember(ids, playerObj.Id)
	msg := &client_server.TeamMsg{
		Team:    teamData,
		Members: members,
		MyData:  myData,
	}
	//logs.Info("teamData:", teamData)
	//logs.Info("members:", members)
	//logs.Info("myData:", myData)
	setting := memberObj.GetTeamPersonSetting(playerObj.Id)
	msg1 := &share_message.TeamSetting{
		PersonSetting:  setting,
		MessageSetting: team.GetTeamMessageSetting(),
	}
	//logs.Info("PersonSetting:", setting)
	//logs.Info("MessageSetting", team.GetTeamMessageSetting())
	msg.Setting = msg1
	chatObj := TeamChatLogMgr.GetRedisTeamChatLogObj(teamId)
	if chatObj != nil {
		msg.TeamChatInfo = chatObj.GetTeamUnReadChatLogs(team.GetSessionLogMaxId(), playerObj)
	}
	return msg
}

func GetTeamReName(t int32, targetId, playerId int64, str string) string {
	var name string
	if t == 2 {
		memberObj := GetRedisTeamPersonalObj(targetId)
		name = memberObj.GetTeamMemberReName(playerId)
		if name == "" {
			name = str
		}
	} else {
		name = str
	}
	return name
}

//-----------------------------------------------------

//func UpdateTeamDataToMongoDB() {
//	teamIds := []int64{}
//	value, err1 := easygo.RedisMgr.GetC().Smembers(REDIS_TABLE_TEAM_ID)
//	easygo.PanicError(err1)
//	InterfersToInt64s(value, &teamIds)
//	for _, teamId := range teamIds {
//		SaveTeamDataToMongoDB(teamId)
//	}
//	_, err2 := easygo.RedisMgr.GetC().Delete(REDIS_TABLE_TEAM_ID)
//	easygo.PanicError(err2)
//}

//func SaveTeamDataToMongoDB(teamId int64) { //定时存储数据到mongodb
//	logs.Info("=========UpdateTeamDataToMongoDB=========", teamId)
//	//teamdata
//	teamData := GetTeamData(teamId)
//	SetTeamDataForMongoDB(teamData)
//
//	//teamchatlog
//	chatLogs := GetTeamSaveChatLogs(teamId)
//	if len(chatLogs) > 0 {
//		SetTeamChatLogForMongoDB(teamId, chatLogs)
//	}
//
//	//teammember
//	teamMember := GetAllTeamMember(teamId)
//	SetTeamMemberForMongoDB(teamMember)
//}

//func GetTeamDataForMongoDB(teamId int64) *share_message.TeamData {
//	var team *share_message.TeamData
//	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAM_DATA)
//	defer closeFun()
//	err := col.Find(bson.M{"_id": teamId}).One(&team)
//	if err != mgo.ErrNotFound && err != nil {
//		easygo.PanicError(err)
//	}
//	return team
//}

//func SetTeamDataForMongoDB(team *share_message.TeamData) {
//	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAM_DATA)
//	defer closeFun()
//	err := col.Update(bson.M{"_id": team.GetId()}, team)
//	easygo.PanicError(err)
//}

//func SetTeamMemberForMongoDB(teamMembers []*share_message.PersonalTeamData) {
//	col1, closeFun1 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAMMEMBER)
//	defer closeFun1()
//	for _, v := range teamMembers {
//		//err2 := col1.Update(bson.M{"TeamId": v.GetTeamId(), "PlayerId": v.GetPlayerId()}, bson.M{"$set": v})
//		err2 := col1.Update(bson.M{"_id": v.GetId()}, bson.M{"$set": v})
//		if err2 != nil && err2 != mgo.ErrNotFound {
//			easygo.PanicError(err2)
//		}
//	}
//}

//func SetTeamChatLogForMongoDB(teamId int64, chatLogs []interface{}) {
//	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_TEAM_CHAT_LOG)
//	defer closeFun()
//	err1 := col.Insert(chatLogs...)
//	easygo.PanicError(err1)
//	maxId := GetTeamMaxChatLogId(teamId)
//	UpdateTeamSaveLogId(teamId, maxId)
//}

//func UpdateAllForList(teamId int64) {
//	team := GetTeamDataForMongoDB(teamId)
//
//}

func (self *RedisTeamObj) SetRedisOperatorInfo(obj *share_message.OperatorInfo) {
	infos := self.GetRedisOperatorInfo()
	infos = append(infos, obj)
	self.SetStringValueToRedis(REDIS_TEAM_OPERATORINFO, infos)
	//js, err := json.Marshal(obj)
	//easygo.PanicError(err)
	//err1 := easygo.RedisMgr.GetC().HSet(self.GetKeyId(), easygo.AnytoA(obj.GetFlag()), string(js))
	//easygo.PanicError(err1)
}
func (self *RedisTeamObj) SetTeamUnBanTime(teamId int64, t int64) {
	setting := self.GetTeamMessageSetting()
	setting.UnBanTime = easygo.NewInt64(t)
	self.UpdateTeamMessageSetting(setting)
}

func (self *RedisTeamObj) GetRedisOperatorInfo() []*share_message.OperatorInfo {
	var inviteList []*share_message.OperatorInfo
	self.GetStringValueToRedis(REDIS_TEAM_OPERATORINFO, &inviteList)
	return inviteList
}

func (self *RedisTeamObj) DelRedisOperatorInfo(flag int32) bool {
	infos := self.GetRedisOperatorInfo()
	for i, info := range infos {
		if info.GetFlag() == flag {
			infos = append(infos[:i], infos[i+1:]...)
			//infos = easygo.Del(infos, info).([]*share_message.OperatorInfo)
			break
		}
	}
	self.SetStringValueToRedis(REDIS_TEAM_OPERATORINFO, infos)
	return true
}

//模拟群主发送欢迎语
func (self *RedisTeamObj) GetSendWelComeWord(playerIds []int64) *share_message.Chat {
	content := ""
	notice := &share_message.NoticeInfo{
		PlayerId: playerIds,
	}
	for i, pid := range playerIds {
		if i > 2 {
			ss := fmt.Sprintf(".....等%d人,", len(playerIds))
			content += ss
			break
		}
		player := GetRedisPlayerBase(pid)
		s := fmt.Sprintf("@%s ", player.GetNickName())
		content += s
	}
	content += self.GetWelcomeWord()
	b64Content := base64.StdEncoding.EncodeToString([]byte(content))
	chatLog := &share_message.Chat{
		SessionId:   easygo.NewString(self.Id),
		SourceId:    easygo.NewInt64(self.GetTeamOwner()),
		TargetId:    easygo.NewInt64(self.GetId()),
		Content:     easygo.NewString(b64Content),
		ChatType:    easygo.NewInt32(CHAT_TYPE_TEAM),
		ContentType: easygo.NewInt32(TALK_CONTENT_WORD),
		NoticeInfo:  notice,
	}
	equiptment := GetRedisPlayerEquipmentObj(self.GetTeamOwner())
	if equiptment != nil {
		eq := equiptment.GetEquipmentForClient()
		chatLog.QPId = easygo.NewInt64(eq.GetQP().GetPropsId())
	}
	return chatLog
}

//新的话题群组动态
func (self *RedisTeamObj) GetTopicTeamDynamic() *share_message.Chat {
	//检测最近群成员的动态
	dynamic := GetNewOneDynamicByPids(self.GetTeamMemberList())
	if dynamic != nil && dynamic.GetLogId() != self.GetDynamicId() {
		//组装发送群里的消息
		content, _ := base64.StdEncoding.DecodeString(dynamic.GetContent()) //特殊处理，前端需要对动态内容先反base64解密后再传输
		msg := client_hall.ChatDynamic{
			LogId:             easygo.NewInt64(dynamic.GetLogId()),              //动态id
			PlayerId:          easygo.NewInt64(dynamic.GetPlayerId()),           //玩家id
			HeadIcon:          easygo.NewString(dynamic.GetHeadIcon()),          //头像
			Sex:               easygo.NewInt32(dynamic.GetSex()),                //性别
			Content:           easygo.NewString(string(content)),                //文本内容
			Photo:             dynamic.GetPhoto(),                               //图片
			Voice:             easygo.NewString(dynamic.GetVoice()),             //语音
			Video:             easygo.NewString(dynamic.GetVideo()),             //视频
			VoiceTime:         easygo.NewInt64(dynamic.GetVoiceTime()),          //录音时长
			NickName:          easygo.NewString(dynamic.GetNickName()),          //昵称
			VideoThumbnailURL: easygo.NewString(dynamic.GetVideoThumbnailURL()), //视频缩略图
			SendTime:          easygo.NewInt64(dynamic.GetSendTime()),           //发布时间
			TopicId:           dynamic.GetTopicId(),                             //话题id
			TopicList:         dynamic.GetTopicList(),                           //话题列表
		}
		player := GetRedisPlayerBase(dynamic.GetPlayerId())
		if player != nil {
			msg.Sex = easygo.NewInt32(player.GetSex())
			msg.NickName = easygo.NewString(player.GetNickName())
			msg.HeadIcon = easygo.NewString(player.GetHeadIcon())
		}
		b, err := json.Marshal(msg)
		if err != nil {
			return nil
		}
		b64Content := base64.StdEncoding.EncodeToString(b)
		chatLog := &share_message.Chat{
			SessionId:   easygo.NewString(self.Id),
			SourceId:    easygo.NewInt64(self.GetTeamOwner()),
			TargetId:    easygo.NewInt64(self.GetId()),
			Content:     easygo.NewString(b64Content),
			ChatType:    easygo.NewInt32(CHAT_TYPE_TEAM),
			ContentType: easygo.NewInt32(TALK_CONTENT_DYNAMIC),
		}
		self.SetDynamicId(dynamic.GetLogId()) //更新id
		return chatLog
	}
	return nil
}

//批量保存需要存储的数据
func SaveRedisTeamToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_TEAM_DATA, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisTeamObj(id)
		if obj != nil {
			data := obj.GetRedisTeam()
			saveData = append(saveData, bson.M{"_id": data.GetId()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_TEAM_DATA, saveData)
	}
}
