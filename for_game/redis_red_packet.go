package for_game

import (
	b64 "encoding/base64"
	"encoding/json"
	"game_server/easygo"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

/*
redis红包数据
*/
const (
	//列表key
	REDPACKET_PACKETS    = "redpacket_packets"     //小红包面额
	REDPACKET_PLAYERLIST = "PlayerList"            //指定可以抢红包的人
	REDPACKET_LOGS       = "Logs"                  //红包日志
	REDPACKET_EXIST_LIST = "redpacket_exist_lists" //修改的红包列表
	REDPACKET_EXIST_TIME = 1000 * 600              //redis的key删除时间:毫秒
)

type RedisRedPacketObj struct {
	Id REDPACKET_ID //红包Id
	RedisBase
}

type RedPacketEx struct {
	Id         int64  `json:"_id"`
	Type       int32  //1私聊发包，2，拼手气红包
	Sender     int64  //发红包的人
	SenderName string //发送人名字
	SenderHead string //发送人头像
	TargetId   int64  //接收人或所在群id
	TotalMoney int64  //红包总金额
	TotalCount int32  //红包总个数
	PerMoney   int64  //每个包的值，如果为0表示拼手气，否则是普通红包
	Content    string //发包留言
	CreateTime int64  //发包时间
	CurMoney   int64  //当前红包余额
	CurCount   int32  //当前剩余红包个数
	State      int32  //红包状态:1可领取，2已领完，3已过期(24小时未领取退回)
	PayWay     int32  //支付方式:99 零钱 1微信 2支付宝 3银行卡
	OrderId    string //交易id
	LuckId     int64  //手气最佳人物id
	OverTime   int64  //领取完时间
	Sex        int32  //发送者性别
}

//对外方法=============================================
func NewRedisRedPacket(id REDPACKET_ID, redPacket ...*share_message.RedPacket) *RedisRedPacketObj {
	p := &RedisRedPacketObj{
		Id: id,
	}
	obj := append(redPacket, nil)[0]
	return p.Init(obj)
}

func (self *RedisRedPacketObj) Init(obj *share_message.RedPacket) *RedisRedPacketObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_RED_PACKET)
	self.Sid = RedPacketMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		RedPacketMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = self.QueryRedPacket(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisRedPacket(obj)
	}
	logs.Info("初始化新的redPacket管理器:", self.Id)
	return self
}
func (self *RedisRedPacketObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisRedPacketObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_RED_PACKET, self.Id)
}

//定时更新数据
func (self *RedisRedPacketObj) UpdateData() { //override
	if !self.IsExistKey() {
		RedPacketMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > REDPACKET_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		RedPacketMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisRedPacketObj) InitRedis() { //override
	obj := self.QueryRedPacket(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisRedPacket(obj)
}
func (self *RedisRedPacketObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisRedPacket()
	return data
}
func (self *RedisRedPacketObj) SaveRedPacketLogs() {
	logList := self.GetRedPacketLogs()
	if logList == nil || len(logList) == 0 {
		return
	}
	var data []interface{}
	for _, v := range logList {
		if v.GetIsSave() {
			v.IsSave = easygo.NewBool(false)
			data = append(data, bson.M{"_id": v.GetId()}, v)
		}
	}
	if len(data) > 0 {
		UpsertAll(easygo.MongoLogMgr, MONGODB_NINGMENG_LOG, TABLE_RED_PACKET_LOG, data)
	}
}
func (self *RedisRedPacketObj) SaveOtherData() { //override
	//存储日志数据
	self.SaveRedPacketLogs()
}
func (self *RedisRedPacketObj) QueryRedPacket(id REDPACKET_ID) *share_message.RedPacket {
	data := self.QueryMongoData(id)
	if data != nil {
		var redPacket share_message.RedPacket
		StructToOtherStruct(data, &redPacket)
		return &redPacket
	}
	return nil
}
func (self *RedisRedPacketObj) SetRedisRedPacket(obj *share_message.RedPacket) {
	//增加到管理器
	RedPacketMgr.Store(obj.GetId(), self)
	self.AddToExistList(obj.GetId())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	redPacket := &RedPacketEx{}
	StructToOtherStruct(obj, redPacket)
	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), redPacket)
	easygo.PanicError(err)
	//设置红包其他数据
	self.SetPackets(obj.GetPackets())
	self.SetRedPacketOtherData(obj, false)
	//红包日志数据
	self.SetRedPacketLogs(nil, false)
}
func (self *RedisRedPacketObj) SetPackets(packets []int64) {
	if packets != nil && len(packets) > 0 {
		err := easygo.RedisMgr.GetC().LPush(MakeRedisKey(REDPACKET_PACKETS, self.Id), packets)
		easygo.PanicError(err)
	}
}
func (self *RedisRedPacketObj) GetPackets() []int64 {
	val, err := easygo.RedisMgr.GetC().LRange(MakeRedisKey(REDPACKET_PACKETS, self.Id), 0, -1)
	easygo.PanicError(err)
	var data []int64
	InterfersToInt64s(val, &data)
	return data
}
func (self *RedisRedPacketObj) PopPackets() int64 {
	val, err := easygo.RedisMgr.GetC().LPop(MakeRedisKey(REDPACKET_PACKETS, self.Id))
	if err != nil {
		logs.Error("PopPackets:", err.Error())
		return 0
	}
	return easygo.AtoInt64(val)
}
func (self *RedisRedPacketObj) SetRedPacketOtherData(obj *share_message.RedPacket, save ...bool) {
	self.SetStringValueToRedis(REDPACKET_PLAYERLIST, obj.GetPlayerList(), save...)
}
func (self *RedisRedPacketObj) GetRedPacketOtherData(obj *share_message.RedPacket) {
	self.GetStringValueToRedis(REDPACKET_PLAYERLIST, &obj.PlayerList)
}

//设置红包日志数据
func (self *RedisRedPacketObj) SetRedPacketLogs(logs []*share_message.RedPacketLog, save ...bool) {
	if logs == nil {
		logs = QueryRedPacketLogs(self.Id)
	}
	self.SetStringValueToRedis(REDPACKET_LOGS, logs, save...)
}

//获取红包日志数据
func (self *RedisRedPacketObj) GetRedPacketLogs() []*share_message.RedPacketLog {
	logs := make([]*share_message.RedPacketLog, 0)
	self.GetStringValueToRedis(REDPACKET_LOGS, &logs)
	return logs
}
func (self *RedisRedPacketObj) AddRedPacketLog(log *share_message.RedPacketLog) {
	logs := self.GetRedPacketLogs()
	logs = append(logs, log)
	self.SetStringValueToRedis(REDPACKET_LOGS, logs)
}
func QueryRedPacketLogs(redPacketId REDPACKET_ID) []*share_message.RedPacketLog {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_RED_PACKET_LOG)
	defer closeFun()
	queryBson := bson.M{"RedPacketId": redPacketId}
	var lst []*share_message.RedPacketLog
	err := col.Find(queryBson).All(&lst)
	easygo.PanicError(err)
	return lst
}
func (self *RedisRedPacketObj) GetRedisRedPacket() *share_message.RedPacket {
	if !self.IsExistKey() {
		//如果key值不存在，先获取
		self.InitRedis()
	}
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	var obj RedPacketEx
	err = redis.ScanStruct(value, &obj)
	easygo.PanicError(err)
	newObj := &share_message.RedPacket{}
	StructToOtherStruct(obj, newObj)
	newObj.Packets = self.GetPackets()
	//其他数据获取
	self.GetRedPacketOtherData(newObj)
	newObj.Logs = self.GetRedPacketLogs()
	return newObj
}

//获取指定可以抢红包的人
func (self *RedisRedPacketObj) GetPlayerList() []int64 {
	playerList := make([]int64, 0)
	self.GetStringValueToRedis(REDPACKET_PLAYERLIST, &playerList)
	return playerList
}
func (self *RedisRedPacketObj) GetState() int32 {
	var val int32
	self.GetOneValue("State", &val)
	return val
}

func (self *RedisRedPacketObj) AddCurMoney(val int64) int64 {
	return self.IncrOneValue("CurMoney", val)
}
func (self *RedisRedPacketObj) AddCurCount(val int64) int64 {
	return self.IncrOneValue("CurCount", val)
}
func (self *RedisRedPacketObj) SetLuckId(val int64) {
	self.SetOneValue("LuckId", val)
}

func (self *RedisRedPacketObj) SetOverTime(val int64) {
	self.SetOneValue("OverTime", val)
}

//打开领取一次红包值：
func (self *RedisRedPacketObj) RedisOpen(playerId PLAYER_ID, ti int64) (int64, string) {
	if self.GetState() > PACKET_MONEY_OPEN {
		logs.Info("红包已经超时了")
		return 0, "红包已经超时了"
	}
	pList := self.GetPlayerList()
	if pList != nil {
		//属于指定领取玩家的红包
		if len(pList) > 0 && !easygo.Contain(pList, playerId) {
			return 0, "红包已经领完"
		}
	}
	logsList := self.GetRedPacketLogs()
	for _, p := range logsList {
		if p.GetPlayerId() == playerId {
			logs.Info("玩家已经领取过红包了")
			return 0, "玩家已经领取过红包了"
		}
	}
	//获取一个红包值
	val := self.PopPackets()
	if val == 0 {
		return 0, "手慢了，红包已被领完"
	}
	money := self.AddCurMoney(-val)
	cnt := self.AddCurCount(-1)
	if money < 0 || cnt < 0 {
		panic("怎么会出现money或者红包数为负数的呢")
	}
	//领取日志
	base := GetRedisPlayerBase(playerId)
	name := GetTeamReName(self.GetType(), self.GetTargetId(), playerId, base.GetNickName())
	logId := NextId(TABLE_RED_PACKET_LOG)
	log := &share_message.RedPacketLog{
		Id:          easygo.NewInt64(logId),
		RedPacketId: easygo.NewInt64(self.Id),
		PlayerId:    easygo.NewInt64(playerId),
		NickName:    easygo.NewString(name),
		HeadUrl:     easygo.NewString(base.GetHeadIcon()),
		Sex:         easygo.NewInt32(base.GetSex()),
		Money:       easygo.NewInt64(val),
		CreateTime:  easygo.NewInt64(ti),
		IsSave:      easygo.NewBool(true), //标识日志需要存储
	}
	self.AddRedPacketLog(log)
	easygo.Spawn(func() { MakePlayerBehaviorReport(5, 0, nil, nil, log, nil) }) //生成用户行为报表领红包字段 已优化到Redis

	if cnt == 0 && money == 0 { //红包被领取完成
		self.SetState(PACKET_MONEY_FNISH)
		if self.GetType() != 1 {
			var maxId, maxGold, RecTime int64
			for _, log := range self.GetRedPacketLogs() {
				gold := log.GetMoney()
				if maxId == 0 {
					maxId = log.GetPlayerId()
					maxGold = gold
					RecTime = log.GetCreateTime()
				} else {
					if gold > maxGold {
						maxId = log.GetPlayerId()
						maxGold = gold
						RecTime = log.GetCreateTime()
					}
				}
			}
			self.SetLuckId(maxId)
			easygo.Spawn(UpdateRedisRedPacketTotal, maxId, RecTime, val, REDPACKET_STATISTICS_LUCK)
		}
		self.SetOverTime(time.Now().Unix())
	}
	return val, ""
}
func CreateRedisRedPacket(msg *share_message.RedPacket) *RedisRedPacketObj {
	return RedPacketMgr.CreateRedisRedPacket(msg)
}

//获取传输给前端的字符串结构
func (self *RedisRedPacketObj) GetRedisSendRedPacketBase64Str() string {
	msg := self.GetRedisRedPacket()
	redPacket := &share_message.RedPacket{
		Id:         easygo.NewInt64(msg.GetId()),
		Type:       easygo.NewInt32(msg.GetType()),
		Sender:     easygo.NewInt64(msg.GetSender()),
		SenderName: easygo.NewString(msg.GetSenderName()),
		SenderHead: easygo.NewString(msg.GetSenderHead()),
		TargetId:   easygo.NewInt64(msg.GetTargetId()),
		TotalMoney: easygo.NewInt64(msg.GetTotalMoney()),
		TotalCount: easygo.NewInt32(msg.GetTotalCount()),
		PerMoney:   easygo.NewInt64(msg.GetPerMoney()),
		CreateTime: easygo.NewInt64(msg.GetCreateTime()),
		Content:    easygo.NewString(msg.GetContent()),
		CurMoney:   easygo.NewInt64(msg.GetCurMoney()),
		CurCount:   easygo.NewInt32(msg.GetCurCount()),
		State:      easygo.NewInt32(msg.GetState()),
		Logs:       self.GetRedPacketLogs(),
		OverTime:   easygo.NewInt64(msg.GetOverTime()),
	}
	content, err := json.Marshal(redPacket)
	easygo.PanicError(err)
	sEnc := b64.StdEncoding.EncodeToString(content)
	return sEnc
}

//获取传输给前端的字符串结构
func GetRedisOpenRedPacketBase64Str(msg *client_hall.OpenRedPacket) string {
	content, err := json.Marshal(msg)
	easygo.PanicError(err)
	sEnc := b64.StdEncoding.EncodeToString(content)
	return sEnc
}

//获取红包发送者
func (self *RedisRedPacketObj) GetSender() int64 {
	var val int64
	self.GetOneValue("Sender", &val)
	return val
}

//获取红包发送者
func (self *RedisRedPacketObj) GetSenderName() string {
	var val string
	self.GetOneValue("SenderName", &val)
	return val
}

//获取红包类型
func (self *RedisRedPacketObj) GetType() int32 {
	var val int32
	self.GetOneValue("Type", &val)
	return val
}

//获取当前红包余额
func (self *RedisRedPacketObj) GetCurMoney() int64 {
	var val int64
	self.GetOneValue("CurMoney", &val)
	return val
}

//获取接收人或所在群id
func (self *RedisRedPacketObj) GetTargetId() int64 {
	var val int64
	self.GetOneValue("TargetId", &val)
	return val
}

//获取支付方式:99 零钱 1微信 2支付宝 3银行卡
func (self *RedisRedPacketObj) GetPayWay() int32 {
	var val int32
	self.GetOneValue("PayWay", &val)
	return val
}

func (self *RedisRedPacketObj) GetCreateTime() int64 {
	var val int64
	self.GetOneValue("CreateTime", &val)
	return val
}

func (self *RedisRedPacketObj) GetLuckId() int64 {
	var val int64
	self.GetOneValue("LuckId", &val)
	return val
}

//修改状态
func (self *RedisRedPacketObj) SetState(v int32) {
	self.SetOneValue("State", v)
}

func GetAllRedPackForRedis(playerId int64, t int32, year, month int) map[int64]*share_message.RedPacket {
	var year2, month2 int
	if month == 12 {
		year2 = year + 1
		month2 = 1
	} else {
		year2 = year
		month2 = month + 1
	}
	t1 := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).Unix()
	t2 := time.Date(year2, time.Month(month2), 1, 0, 0, 0, 0, time.Local).Unix()

	info := make(map[int64]*share_message.RedPacket)
	lst := make([]int64, 0)
	GetAllRedisExistList(TABLE_RED_PACKET, &lst)

	for _, id := range lst {
		redPacket := GetRedisRedPacket(id)
		if redPacket.GetCreateTime() < int64(t1) || redPacket.GetCreateTime() > int64(t2) {
			continue
		}
		if t == 1 { //如果是收红包
			if redPacket.GetType() == 1 { //如果是个人红包
				if redPacket.GetTargetId() == playerId { //并且收包人是自己
					info[redPacket.Id] = redPacket.GetRedisRedPacket()
				}
			} else { //如果是群红包
				openList := redPacket.GetRedPacketLogs()
				for _, log := range openList {
					if log.GetPlayerId() == playerId { //如果自己领取过红包
						info[redPacket.Id] = redPacket.GetRedisRedPacket()
						break
					}
				}
			}
		} else { //如果是发红包
			if redPacket.GetSender() == playerId {
				info[redPacket.Id] = redPacket.GetRedisRedPacket()
			}
		}
	}
	return info
}

// 封装外部方法，获取红包的信息
func GetRedisRedPacket(id REDPACKET_ID) *RedisRedPacketObj {
	return RedPacketMgr.GetRedisRedPacketObj(id)
}

//检测玩家是否有发送未领取完的红包
func CheckPlayerRedPacket(playerId PLAYER_ID) bool {
	//redis检测
	ids := make([]int64, 0)
	GetAllRedisExistList(TABLE_RED_PACKET, &ids)
	for _, id := range ids {
		base := GetRedisRedPacket(id)
		redPacket := base.GetRedisRedPacket()
		if redPacket.GetSender() == playerId && redPacket.GetState() == PACKET_MONEY_OPEN {
			return false
		}
	}
	//数据库检测
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_RED_PACKET)
	defer closeFun()
	queryBson := bson.M{"Sender": playerId, "State": PACKET_MONEY_OPEN}
	lst := []*share_message.RedPacket{}
	err := col.Find(queryBson).All(&lst)
	if err != nil && err != mgo.ErrNotFound {
		return false
	}
	if len(lst) > 0 {
		return false
	}
	return true
}

//停服保存处理，保存需要存储的数据
func SaveRedisRedPacketToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_RED_PACKET, &ids)
	saveData := make([]interface{}, 0)
	logList := make([]*share_message.RedPacketLog, 0)
	for _, id := range ids {
		obj := GetRedisRedPacket(id)
		if obj != nil {
			data := obj.GetRedisRedPacket()
			saveData = append(saveData, bson.M{"_id": data.GetId()}, data)
			obj.SetSaveStatus(false)
			logs := obj.GetRedPacketLogs()
			logList = append(logList, logs...)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_RED_PACKET, saveData)
	}
	if len(logList) > 0 {
		//领取日志保存
		saveLogs := make([]interface{}, 0)
		for _, it := range logList {
			saveLogs = append(saveLogs, bson.M{"_id": it.GetId()}, it)
		}
		UpsertAll(easygo.MongoLogMgr, MONGODB_NINGMENG_LOG, TABLE_RED_PACKET_LOG, saveLogs)
	}
}
