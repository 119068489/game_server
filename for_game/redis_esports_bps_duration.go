package for_game

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"sync"
)

//在某个某个比赛房间的用户
type RedisESportBpsDurationLogObj struct {
	PlayerId int64
	RedisBase
	BpsDuration   sync.Map
	IsDeleteRedis bool //是否删除redis
}
type RedisESportBpsDuration struct {
	//埋點處理
	Id           string
	MenuId       int32
	PageType     int32
	DataId       int64
	ExId         int64
	TimeKey      int64
	BeginTime    int64
	EndTime      int64
	ExTabId      int64
	LabelId      int64
	Duration     int64
	NavigationId int32
	//埋點處理
}

func (this *RedisESportBpsDuration) BpsKey(playerId int64) string {
	key := GetBpsDurationIdkey(playerId, this.PageType, this.MenuId, this.LabelId, this.ExTabId, this.DataId, this.ExId, this.NavigationId)
	return key
}
func GetBpsDurationIdkey(plyid int64, PageType, MenuId int32, LabelId, ExTabId, DataId, ExId int64, navigationId int32) string {
	ksy := fmt.Sprintf("%d_%d_%d_%d_%d_%d_%d_%d",
		plyid, PageType, MenuId, LabelId, ExTabId, DataId, ExId, navigationId)
	return ksy
}

const REDIS_DURATION_MIN_TIME int64 = 2                               //停留时长最小保留时间，小于就直接删除，单位秒
const REDIS_DURATION_LOG_EXIST_TIME = 1000 * 600                      //毫秒，key值存在时间
var REDIS_ESPORT_BPS_DURATION_LOG_KEY = ESportExN("bps_duration_log") //redis内存中存在的key

func (self *RedisESportBpsDurationLogObj) Init(playerId int64) *RedisESportBpsDurationLogObj {
	self.PlayerId = playerId
	self.RedisBase.Init(self, self.PlayerId, easygo.MongoMgr, MONGODB_NINGMENG_LOG, REDIS_ESPORT_BPS_DURATION_LOG_KEY)
	self.Sid = ESportBpsDurationMgr.GetSid()
	self.AddToExistList(self.PlayerId)
	ESportBpsDurationMgr.Store(self.PlayerId, self)
	self.IsDeleteRedis = false
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)

	return self
}
func (self *RedisESportBpsDurationLogObj) GetId() interface{} { //override
	return self.PlayerId
}
func (self *RedisESportBpsDurationLogObj) GetKeyId() string { //override
	return MakeRedisKey(REDIS_ESPORT_BPS_DURATION_LOG_KEY, self.PlayerId)
}
func (self *RedisESportBpsDurationLogObj) UpdateData() { //override

	self.SaveToMongo()
	if self.IsDeleteRedis { //用户退出的时候控制
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.PlayerId)
			self.DelRedisKey() //redis删除
		}
		ESportBpsDurationMgr.Delete(self.PlayerId) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisESportBpsDurationLogObj) GetSaveList() []*share_message.TableESPortsBpsDuration {

	savelist := []*share_message.TableESPortsBpsDuration{}
	self.BpsDuration.Range(func(k, v interface{}) bool {
		if v != nil {
			logs.Info(v)
			it := v.(*share_message.TableESPortsBpsDuration)
			savelist = append(savelist, it)
		}
		return true
	})
	return savelist
}

//重写保存方法
func (self *RedisESportBpsDurationLogObj) SaveToMongo() {

	logsData := self.GetSaveList()
	//logs.Info("data:", logsData)
	var saveData []interface{}
	for _, log := range logsData {
		saveData = append(saveData, bson.M{"_id": log.GetId()}, log)
		self.BpsDuration.Delete(log.GetId())
	}
	if len(saveData) > 0 {
		logs.Info("埋点停留保存")
		logs.Info(saveData)
		UpsertAll(easygo.MongoLogMgr, MONGODB_NINGMENG_LOG, TABLE_ESPORTS_BPS_DURATION_LOG, saveData)
	}

	self.SetSaveStatus(false)
}

//当检测到redis key不存在时
func (self *RedisESportBpsDurationLogObj) InitRedis() { //override

}
func (self *RedisESportBpsDurationLogObj) GetRedisSaveData() interface{} { //override
	return nil
}
func (self *RedisESportBpsDurationLogObj) SaveOtherData() { //override

}
func (self *RedisESportBpsDurationLogObj) MakePlyerKey(pageType int32) string {
	return fmt.Sprintf("%d", pageType)
}

//进入页面
func (self *RedisESportBpsDurationLogObj) EnterPage(pageType int32, data string) {
	self.SetOneValue(self.MakePlyerKey(pageType), data)
}

//删除页面
func (self *RedisESportBpsDurationLogObj) DeletePage(pageType int32) {
	_, err := easygo.RedisMgr.GetC().Hdel(self.GetKeyId(), self.MakePlyerKey(pageType))
	easygo.PanicError(err)
}

func (self *RedisESportBpsDurationLogObj) GetItemBpsDuration(pageType int32) map[int32]*RedisESportBpsDuration {

	mps := self.GetBpsDurationStr()
	savelist := make(map[int32]*RedisESportBpsDuration)
	for k, v := range mps {
		it := &RedisESportBpsDuration{}
		_ = json.Unmarshal([]byte(v), &it)
		if it.PageType >= pageType {
			savelist[int32(k)] = it //savelist = append(savelist, it)
		}
	}
	return savelist
}

func (self *RedisESportBpsDurationLogObj) GetBpsDurationStr() map[int64]string {

	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(self.GetKeyId()))
	easygo.PanicError(err)
	return values
}

//保存自上的所有停留買點 包括自己
func (self *RedisESportBpsDurationLogObj) SaveBpsDuration(pageType int32) {

	savelist := self.GetItemBpsDuration(pageType)
	for _, v := range savelist {
		logs.Info("保存", v)
		self.SaveBpsDurationItem(v)
		self.DeletePage(v.PageType)
	}

}

//保存停留时间
func (self *RedisESportBpsDurationLogObj) BeginCurrentBpsDuration(pageType int32, menuId int32, labelId, exTabId int64, dataId int64, exId int64, navigationId int32) {
	b := self.EndLastBpsDuration(pageType, menuId, labelId, exTabId, dataId, exId, navigationId) //保存结束上一次
	if b {
		//如果上一个路径和这次路径一样~直接返回
		return
	}
	data := RedisESportBpsDuration{
		MenuId:       menuId,
		PageType:     pageType,
		DataId:       dataId,
		ExId:         exId,
		TimeKey:      easygo.GetToday0ClockTimestamp(),
		BeginTime:    easygo.NowTimestamp(),
		EndTime:      0,
		ExTabId:      exTabId,
		LabelId:      labelId,
		NavigationId: navigationId,
	}
	s, _ := json.Marshal(data)
	self.EnterPage(pageType, string(s))
}
func (self *RedisESportBpsDurationLogObj) SavePage(pageType int32, data *RedisESportBpsDuration) {

	if data != nil {
		s, _ := json.Marshal(data)
		self.EnterPage(pageType, string(s))
	}
}

//保存停留时间
func (self *RedisESportBpsDurationLogObj) SaveBpsDurationItem(data *RedisESportBpsDuration) {

	now := easygo.NowTimestamp()
	durationSecond := now - data.BeginTime
	if durationSecond > 24*60*60 { //超过24小时的就属于异常数据
		self.DeletePage(data.PageType)
		return
	}
	if durationSecond >= REDIS_DURATION_MIN_TIME { //小于2秒的直接删除，不保存
		plyid := self.PlayerId
		//_ := GetRedisESportBpsDurationLogObj(self.PlayerId)
		tis := easygo.NowTimestamp()
		id := fmt.Sprintf("%d_%d_%d", plyid, data.PageType, tis)
		data.EndTime = tis
		data.Duration = durationSecond
		data.Id = id
		log := &share_message.TableESPortsBpsDuration{
			Id:           easygo.NewString(id),
			TimeKey:      easygo.NewInt64(data.TimeKey),
			PlayerId:     easygo.NewInt64(plyid),
			CreateTime:   easygo.NewInt64(data.BeginTime),
			EndTime:      easygo.NewInt64(now),
			Duration:     easygo.NewInt64(data.Duration),
			MenuId:       easygo.NewInt32(data.MenuId),
			LabelId:      easygo.NewInt64(data.LabelId),
			ExTabId:      easygo.NewInt64(data.ExTabId),
			DataId:       easygo.NewInt64(data.DataId),
			ExId:         easygo.NewInt64(data.ExId),
			PageType:     easygo.NewInt32(data.PageType),
			NavigationId: easygo.NewInt32(data.NavigationId),
		}
		logs.Info("停留了", durationSecond, "秒")
		self.AddESPortsBpsDuration(log) //加入保存队列
	}
	self.DeletePage(data.PageType)
}
func (self *RedisESportBpsDurationLogObj) AddESPortsBpsDuration(data *share_message.TableESPortsBpsDuration) {
	self.BpsDuration.Store(data.GetId(), data)
}

//返回会否同样的key
func (self *RedisESportBpsDurationLogObj) EndLastBpsDuration(pageType int32, menuId int32, labelId, exTabId int64, dataId int64, exId int64, navigationId int32) bool {
	list := self.GetItemBpsDuration(pageType)
	if len(list) < 1 {
		return false
	}
	it := list[pageType]
	if it != nil {
		curkey := GetBpsDurationIdkey(self.PlayerId, pageType, menuId, labelId, exTabId, dataId, exId, navigationId)
		if it.BpsKey(self.PlayerId) != curkey {
			self.SaveBpsDuration(pageType)
			return false
		} else {
			return true
		}
	}
	return false
}

//保存停留时间
func (self *RedisESportBpsDurationLogObj) EndCurrentBpsDuration(pageType int32) { //, menuId int32, labelId, exTabId int64, dataId int64, exId int64) {

	//list := self.GetItemBpsDuration(pageType)
	//if len(list) < 1 {
	//	return
	//	}
	//it := list[pageType]
	//if it != nil {
	//logs.Info("当前埋点路径", it)
	//logs.Info("新的埋点路径", "pageType", pageType, "menuId", menuId, "labelId", labelId, "exTabId", exTabId, "dataId", dataId, "exId", exId)
	//curkey := GetBpsClickIdkeyEx(self.PlayerId, pageType, menuId, labelId, exTabId, dataId, exId)
	//if it.BpsKey(self.PlayerId) == curkey {
	self.SaveBpsDuration(pageType)
	//}

	//}
}

func NewRedisESportBpsDurationLogObj(id int64) *RedisESportBpsDurationLogObj {
	obj := RedisESportBpsDurationLogObj{}
	return obj.Init(id)
}

//对外方法，获取对象，如果为nil表示redis内存不存在，数据库也不存在
func GetRedisESportBpsDurationLogObj(id int64) *RedisESportBpsDurationLogObj {
	obj, ok := ESportBpsDurationMgr.Load(id)
	if ok && obj != nil {
		return obj.(*RedisESportBpsDurationLogObj)
	} else {
		return NewRedisESportBpsDurationLogObj(id)
	}
}
