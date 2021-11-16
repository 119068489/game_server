package for_game

import (
	"game_server/easygo"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	REDPACKET_STATISTICS_RECV int32 = 1 //接收红包
	REDPACKET_STATISTICS_LUCK int32 = 2 //手气最佳
	REDPACKET_STATISTICS_SEND int32 = 3 //发送红包
)

const (
	REDPACKET_TOTAL_EXIT_LIST = "redpacket_total_exist_list"
	REDPACKET_TOTAL_EXIT_TIME = 1000 * 600
)

type RedisRedPacketTotalObj struct {
	Id string //玩家+每月开始时间组合:"1887437024_1585670400"
	RedisBase
}
type RedisRedPacketTotal struct {
	Id             string `json:"_id"` //
	PlayerId       int64  //玩家id
	RecTotalMoney  int64  //收到红包总钱数
	RecCount       int64  //收到的红包个数
	LuckCnt        int64  //最佳手气次数
	SendTotalMoney int64  //发送红包总钱数
	SendCount      int64  //发送的红包个数
	CreateTime     int64  //本月一号的时间戳
}

//写入redis
func NewRedPacketTotal(id string, totals ...*share_message.RedPacketTotalInfo) *RedisRedPacketTotalObj {

	p := &RedisRedPacketTotalObj{
		Id: id,
	}
	obj := append(totals, nil)[0]
	return p.Init(obj)
}
func (self *RedisRedPacketTotalObj) Init(obj *share_message.RedPacketTotalInfo) *RedisRedPacketTotalObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_REDPACKET_STATISTICS)
	self.Sid = RedPacketTotalMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		RedPacketTotalMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = self.QueryRedPacketTotal()
			if obj == nil {
				return nil
			}
		}
		self.SetRedisRedPacketTotal(obj)
	}

	logs.Info("初始化新的redpacketTotal管理器:", self.Id)
	return self
}

func (self *RedisRedPacketTotalObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisRedPacketTotalObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_REDPACKET_STATISTICS, self.Id)
}

//定时更新数据
func (self *RedisRedPacketTotalObj) UpdateData() { //override
	if !self.IsExistKey() {
		RedPacketTotalMgr.Delete(self.Id) // 释放对象
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > REDPACKET_TOTAL_EXIT_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		RedPacketTotalMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisRedPacketTotalObj) InitRedis() { //override
	obj := self.QueryRedPacketTotal()
	if obj == nil {
		return
	}
	self.SetRedisRedPacketTotal(obj)
}
func (self *RedisRedPacketTotalObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisRedPacketTotal()
	return data
}
func (self *RedisRedPacketTotalObj) SaveOtherData() { //override
}

func (self *RedisRedPacketTotalObj) QueryRedPacketTotal() *share_message.RedPacketTotalInfo {
	data := self.QueryMongoData(self.Id)
	if data != nil {
		var info share_message.RedPacketTotalInfo
		StructToOtherStruct(data, &info)
		return &info
	}
	return nil
}
func (self *RedisRedPacketTotalObj) SetRedisRedPacketTotal(obj *share_message.RedPacketTotalInfo) {
	//增加到管理器
	RedPacketTotalMgr.Store(obj.GetId(), self)
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//重新激活定时器
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	total := &RedisRedPacketTotal{}
	StructToOtherStruct(obj, total)

	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), total)
	easygo.PanicError(err)
	self.AddToExistList(obj.GetId())
}
func (self *RedisRedPacketTotalObj) GetRedisRedPacketTotal() *share_message.RedPacketTotalInfo {
	obj := &RedisRedPacketTotal{}
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.RedPacketTotalInfo{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

//设置成员值
func (self *RedisRedPacketTotalObj) IncrSendTotalMoney(val int64) {
	self.IncrOneValue("SendTotalMoney", val)
}
func (self *RedisRedPacketTotalObj) GetSendTotalMoney() int64 {
	var val int64
	self.GetOneValue("SendTotalMoney", &val)
	return val
}
func (self *RedisRedPacketTotalObj) IncrSendCount(val int64) {
	self.IncrOneValue("SendCount", val)
}
func (self *RedisRedPacketTotalObj) GetSendCount() int64 {
	var val int64
	self.GetOneValue("SendCount", &val)
	return val
}
func (self *RedisRedPacketTotalObj) IncrRecTotalMoney(val int64) {
	self.IncrOneValue("RecTotalMoney", val)
}
func (self *RedisRedPacketTotalObj) GetRecTotalMoney() int64 {
	var val int64
	self.GetOneValue("RecTotalMoney", &val)
	return val
}
func (self *RedisRedPacketTotalObj) IncrRecCount(val int64) {
	self.IncrOneValue("RecCount", val)
}
func (self *RedisRedPacketTotalObj) GetRecCount() int64 {
	var val int64
	self.GetOneValue("RecCount", &val)
	return val
}
func (self *RedisRedPacketTotalObj) IncrLuckCnt(val int64) {
	self.IncrOneValue("LuckCnt", val)
}
func (self *RedisRedPacketTotalObj) GetLuckCnt() int64 {
	var val int64
	self.GetOneValue("LuckCnt", &val)
	return val
}

//获取当前id
func MakeRedisRedPacketTotalId(pid, ti int64) string {
	year, month, _ := time.Unix(ti/1e3, 0).Date()
	start, _ := easygo.GetMouthStartEnd(year, time.Month(month))
	return MakeNewString(pid, start)
}

//对外方法
func GetRedisRedPacketTotal(id string) *RedisRedPacketTotalObj {
	return RedPacketTotalMgr.GetRedisRedPackeTotalObj(id)
}
func GetRedisRedPacketTotalEx(pid int64, year, month int) *RedisRedPacketTotalObj {
	start, _ := easygo.GetMouthStartEnd(year, time.Month(month))
	id := MakeNewString(pid, start)
	obj := RedPacketTotalMgr.GetRedisRedPackeTotalObj(id)
	if obj == nil {
		data := &share_message.RedPacketTotalInfo{
			Id:             easygo.NewString(id),
			PlayerId:       easygo.NewInt64(pid),
			CreateTime:     easygo.NewInt64(start),
			RecTotalMoney:  easygo.NewInt64(0),
			RecCount:       easygo.NewInt64(0),
			LuckCnt:        easygo.NewInt64(0),
			SendTotalMoney: easygo.NewInt64(0),
			SendCount:      easygo.NewInt64(0),
		}
		obj = NewRedPacketTotal(id, data)
	}
	return obj
}

//新增数据
func CreateRedPacketTotalInfo(pid, ti int64) *RedisRedPacketTotalObj {
	year, month, _ := time.Unix(ti/1e3, 0).Date()
	start, _ := easygo.GetMouthStartEnd(year, time.Month(month))
	id := MakeNewString(pid, start)
	msg := &share_message.RedPacketTotalInfo{
		Id:             easygo.NewString(id),
		PlayerId:       easygo.NewInt64(pid),
		CreateTime:     easygo.NewInt64(start),
		SendTotalMoney: easygo.NewInt64(0),
		SendCount:      easygo.NewInt64(0),
		RecTotalMoney:  easygo.NewInt64(0),
		RecCount:       easygo.NewInt64(0),
		LuckCnt:        easygo.NewInt64(0),
	}
	obj := NewRedPacketTotal(id, msg)
	obj.SetSaveStatus(true) //标志数据要存储
	return obj
}

//更新统计值
//pid:玩家id，ti:更新时间，value:更新值，t:更新类型
func UpdateRedisRedPacketTotal(pid, ti, value int64, t int32) {
	RedPacketTotalMgr.UpdateRedisRedPacketTotal(pid, ti, value, t)
}

//删库重新生成
func StatisticsRedPacket() {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_REDPACKET_STATISTICS)
	defer closeFun()

	var revList, sendList []*GoldLog
	col1, closeFun1 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_GOLDCHANGELOG)
	defer closeFun1()
	err := col1.Find(bson.M{"SourceType": GOLD_TYPE_GET_REDPACKET}).All(&revList)
	easygo.PanicError(err)
	err1 := col1.Find(bson.M{"SourceType": GOLD_TYPE_SEND_REDPACKET}).All(&sendList)
	easygo.PanicError(err1)

	info := make(map[int64]map[int64]*share_message.RedPacketTotalInfo)

	var redpackets []*share_message.RedPacket
	packets := make(map[int64]*share_message.RedPacket)
	col2, closeFun2 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_RED_PACKET)
	defer closeFun2()
	err = col2.Find(bson.M{}).All(&redpackets)
	easygo.PanicError(err)
	for _, v := range redpackets {
		packets[v.GetId()] = v
	}
	for _, log := range revList {
		pid := log.GetPlayerId()
		if m, ok := info[pid]; ok { //有pid
			t := log.GetCreateTime()
			year, month, _ := time.Unix(t/1e3, 0).Date()
			start, _ := easygo.GetMouthStartEnd(year, month)
			if n, ok := m[start]; ok { //有这个月的记录
				n.RecCount = easygo.NewInt64(n.GetRecCount() + 1)
				n.RecTotalMoney = easygo.NewInt64(n.GetRecTotalMoney() + log.GetChangeGold())
				extend := log.Extend
				redpacket := packets[extend.GetRedPacketId()]
				luckId := redpacket.GetLuckId()
				if luckId == pid {
					n.LuckCnt = easygo.NewInt64(n.GetLuckCnt() + 1)
				}
				m[start] = n
				info[pid] = m
			} else { //没有这个月的记录
				id := MakeRedisRedPacketTotalId(pid, t)
				n := &share_message.RedPacketTotalInfo{
					Id:            easygo.NewString(id),
					PlayerId:      easygo.NewInt64(pid),
					RecTotalMoney: easygo.NewInt64(log.GetChangeGold()),
					RecCount:      easygo.NewInt64(1),
				}
				extend := log.Extend
				redpacket := packets[extend.GetRedPacketId()]
				luckId := redpacket.GetLuckId()
				if luckId == pid {
					n.LuckCnt = easygo.NewInt64(1)
				}
				t := log.GetCreateTime()
				year, month, _ := time.Unix(t/1e3, 0).Date()
				start, _ := easygo.GetMouthStartEnd(year, month)
				n.CreateTime = easygo.NewInt64(start)
				m[start] = n
				info[pid] = m
			}
		} else { //没有pid
			t := log.GetCreateTime()
			id := MakeRedisRedPacketTotalId(pid, t)
			n := &share_message.RedPacketTotalInfo{
				Id:            easygo.NewString(id),
				PlayerId:      easygo.NewInt64(pid),
				RecTotalMoney: easygo.NewInt64(log.GetChangeGold()),
				RecCount:      easygo.NewInt64(1),
			}
			extend := log.Extend
			redpacket := packets[extend.GetRedPacketId()]
			luckId := redpacket.GetLuckId()
			if luckId == pid {
				n.LuckCnt = easygo.NewInt64(1)
			}

			year, month, _ := time.Unix(t/1e3, 0).Date()
			start, _ := easygo.GetMouthStartEnd(year, month)
			n.CreateTime = easygo.NewInt64(start)
			in := map[int64]*share_message.RedPacketTotalInfo{
				start: n,
			}
			info[pid] = in
		}
	}

	for _, log := range sendList {
		pid := log.GetPlayerId()
		if m, ok := info[pid]; ok { //有pid
			t := log.GetCreateTime()
			year, month, _ := time.Unix(t/1e3, 0).Date()
			start, _ := easygo.GetMouthStartEnd(year, month)
			if n, ok := m[start]; ok { //有这个月的记录
				n.SendCount = easygo.NewInt64(n.GetSendCount() + 1)
				n.SendTotalMoney = easygo.NewInt64(n.GetSendTotalMoney() + (-log.GetChangeGold()))
				m[start] = n
				info[pid] = m
			} else { //没有这个月的记录
				id := MakeRedisRedPacketTotalId(pid, t)
				n := &share_message.RedPacketTotalInfo{
					Id:             easygo.NewString(id),
					PlayerId:       easygo.NewInt64(pid),
					SendTotalMoney: easygo.NewInt64(-log.GetChangeGold()),
					SendCount:      easygo.NewInt64(1),
				}
				year, month, _ := time.Unix(t/1e3, 0).Date()
				start, _ := easygo.GetMouthStartEnd(year, month)
				n.CreateTime = easygo.NewInt64(start)
				m[start] = n
				info[pid] = m
			}
		} else { //没有pid
			t := log.GetCreateTime()
			id := MakeRedisRedPacketTotalId(pid, t)
			n := &share_message.RedPacketTotalInfo{
				Id:             easygo.NewString(id),
				PlayerId:       easygo.NewInt64(pid),
				SendTotalMoney: easygo.NewInt64(-log.GetChangeGold()),
				SendCount:      easygo.NewInt64(1),
			}
			year, month, _ := time.Unix(t/1e3, 0).Date()
			start, _ := easygo.GetMouthStartEnd(year, month)
			n.CreateTime = easygo.NewInt64(start)
			in := map[int64]*share_message.RedPacketTotalInfo{
				start: n,
			}
			info[pid] = in
		}
	}
	//logs.Info("=================", info)
	var lst []interface{}
	for _, m := range info {
		for _, log := range m {
			lst = append(lst, log)
		}

	}
	err2 := col.Insert(lst...)
	easygo.PanicError(err2)
}

//红包统计表
func ReSetRedPacketStatistics() {
	var val []*share_message.RedPacketTotalInfo
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_REDPACKET_STATISTICS)
	defer closeFun()
	err := col.Find(bson.M{}).All(&val)
	easygo.PanicError(err)
	err = col.DropCollection()
	easygo.PanicError(err)
	var saveData []interface{}
	for _, v := range val {
		if v.GetPlayerId() == 0 {
			continue
		}
		v.Id = easygo.NewString(MakeNewString(v.GetPlayerId(), v.GetCreateTime()))
		saveData = append(saveData, bson.M{"_id": v.Id}, v)
	}
	UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_REDPACKET_STATISTICS, saveData)
}

//停服保存处理，保存需要存储的数据
func SaveRedPacketTotalToMongoDB() {
	ids := []string{}
	GetAllRedisSaveList(TABLE_REDPACKET_STATISTICS, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisRedPacketTotal(id)
		if obj != nil {
			data := obj.GetRedisRedPacketTotal()
			saveData = append(saveData, bson.M{"_id": data.GetId()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_REDPACKET_STATISTICS, saveData)
	}
}
