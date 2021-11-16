package for_game

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"time"
)

//在某个某个比赛房间的用户
type RedisESportBpsClickLogObj struct {
	PlayerId int64
	RedisBase
	IsDeleteRedis bool //是否删除redis
}
type RedisESportBpsClick struct {
	//埋點處理
	MenuId       int32
	PageType     int32
	DataId       int64
	ExId         int64
	TimeKey      int64
	BeginTime    int64
	ExTabId      int64
	LabelId      int64
	ActCount     int32
	NavigationId int32
	DataType     int32
	ButtonId     int32
	//埋點處理
}

func (this *RedisESportBpsClick) GetKey(plyid int64) string {
	return GetBpsClickIdkey(plyid, this.TimeKey, this.PageType, this.MenuId, this.LabelId, this.ExTabId, this.DataId, this.ExId, this.NavigationId, this.DataType, this.ButtonId)
}

func GetBpsClickIdkey(plyid, TimeKey int64, PageType, MenuId int32, LabelId, ExTabId, DataId, ExId int64, navigationId int32, dataType int32, buttonId int32) string {
	ksy := fmt.Sprintf("%d_%d_%d_%d_%d_%d_%d_%d_%d_%d_%d",
		plyid, TimeKey, PageType, MenuId, LabelId, ExTabId, DataId, ExId, navigationId, dataType, buttonId)
	return ksy
}
func GetBpsClickIdkeyEx(plyid int64, PageType, MenuId int32, LabelId, ExTabId, DataId, ExId int64, navigation int32, dataType int32, buttonId int32) string {
	ksy := fmt.Sprintf("%d_%d_%d_%d_%d_%d_%d_%d_%d_%d",
		plyid, PageType, MenuId, LabelId, ExTabId, DataId, ExId, navigation, dataType, buttonId)
	return ksy
}

const REDIS_Click_LOG_EXIST_TIME = 1000 * 600                   //毫秒，key值存在时间
var REDIS_ESPORT_BPS_Click_LOG_KEY = ESportExN("bps_click_log") //redis内存中存在的key

func (self *RedisESportBpsClickLogObj) Init(playerId int64) *RedisESportBpsClickLogObj {
	self.PlayerId = playerId
	self.RedisBase.Init(self, self.PlayerId, easygo.MongoMgr, MONGODB_NINGMENG_LOG, REDIS_ESPORT_BPS_Click_LOG_KEY)
	self.Sid = ESportBpsClickMgr.GetSid()
	self.AddToExistList(self.PlayerId)
	ESportBpsClickMgr.Store(self.PlayerId, self)
	self.IsDeleteRedis = false
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)

	return self
}
func (self *RedisESportBpsClickLogObj) GetId() interface{} { //override
	return self.PlayerId
}
func (self *RedisESportBpsClickLogObj) GetKeyId() string { //override
	return MakeRedisKey(REDIS_ESPORT_BPS_Click_LOG_KEY, self.PlayerId)
}
func (self *RedisESportBpsClickLogObj) UpdateData() { //override

	self.SaveToMongo()

	if self.IsDeleteRedis { //用户退出的时候控制

		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.PlayerId)
			self.DelRedisKey() //redis删除
		}
		ESportBpsClickMgr.Delete(self.PlayerId) // 释放对象
		return

	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}

//重写保存方法
func (self *RedisESportBpsClickLogObj) SaveToMongo() {

	logsData := self.GetAllItemBpsClickList()
	//logs.Info("data:", logsData)
	var saveData []interface{}
	nowt := easygo.GetToday0ClockTimestamp()
	for _, log := range logsData {
		if nowt > log.GetCreateTime() { //刪除昨天的
			_, err := easygo.RedisMgr.GetC().Hdel(self.GetKeyId(), log.GetId())
			if err != nil {
				logs.Error(err)
			}
		} else {
			saveData = append(saveData, bson.M{"_id": log.GetId()}, log)
		}

	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoLogMgr, MONGODB_NINGMENG_LOG, TABLE_ESPORTS_BPS_CLICK_LOG, saveData)
	}

	self.SetSaveStatus(false)
}

//当检测到redis key不存在时
func (self *RedisESportBpsClickLogObj) InitRedis() { //override

}
func (self *RedisESportBpsClickLogObj) GetRedisSaveData() interface{} { //override
	return nil
}
func (self *RedisESportBpsClickLogObj) SaveOtherData() { //override

}

//进入页面
func (self *RedisESportBpsClickLogObj) SaveKeyData(key, data string) {
	self.SetOneValue(key, data)
}

func (self *RedisESportBpsClickLogObj) GetKeyData(key string) string {
	var val string
	self.GetOneValue(key, &val)
	return val
}

func (self *RedisESportBpsClickLogObj) GetAllItemBpsClickList() map[string]*share_message.TableESPortsBpsClick {

	mps := self.GetBpsClickStr()
	savelist := make(map[string]*share_message.TableESPortsBpsClick)
	for k, v := range mps {
		it := &RedisESportBpsClick{}
		_ = json.Unmarshal([]byte(v), &it)
		vit := self.BpsClickItemToTableStruct(k, it)
		savelist[k] = vit
	}
	return savelist
}

func (self *RedisESportBpsClickLogObj) GetBpsClickStr() map[string]string {

	values, err := StrkeyStringMap(easygo.RedisMgr.GetC().HGetAll(self.GetKeyId()))
	easygo.PanicError(err)
	return values
}

//保存停留时间
func (self *RedisESportBpsClickLogObj) BpsClick(rm *client_hall.ESPortsBpsClickRequest) { //) {

	pageType := rm.GetPageType()
	menuId := rm.GetMenuId()
	labelId := rm.GetLabelId()
	exTabId := rm.GetExTabId()
	dataId := rm.GetDataId()
	exId := rm.GetExId()
	navigationId := rm.GetNavigationId()
	dataType := rm.GetDataType()
	buttonId := rm.GetButtonId()
	timekey := easygo.GetToday0ClockTimestamp()
	key := GetBpsClickIdkey(self.PlayerId, timekey, pageType, menuId, labelId, exTabId, dataId, exId, navigationId, dataType, buttonId)
	strdata := self.GetKeyData(key)
	var data = &RedisESportBpsClick{}
	if strdata != "" && len(strdata) > 0 {
		_ = json.Unmarshal([]byte(strdata), &data)
	} else {
		data = &RedisESportBpsClick{
			MenuId:       menuId,
			PageType:     pageType,
			DataId:       dataId,
			ExId:         exId,
			TimeKey:      easygo.GetToday0ClockTimestamp(),
			BeginTime:    easygo.NowTimestamp(),
			ExTabId:      exTabId,
			LabelId:      labelId,
			ActCount:     0,
			NavigationId: navigationId,
			DataType:     dataType,
			ButtonId:     buttonId,
		}
	}
	data.ActCount++
	s, _ := json.Marshal(data)
	self.SaveKeyData(key, string(s))
}

//时间戳转时间
func unixToStr(timeUnix int64) string {
	layout := "2006-01-02 15:04:05"
	timeStr := time.Unix(timeUnix, 0).Format(layout)
	return timeStr
}

//保存停留时间
func (self *RedisESportBpsClickLogObj) BpsClickEx(reqMsg *client_hall.ESPortsBpsClickListRequest) { //) {
	rmList := reqMsg.BpsDataList
	for _, v := range rmList {

		beginTime := easygo.NowTimestamp() - (reqMsg.GetClientTime() - v.GetClientClickTime()) //公式： 服务器时间 - （客户端时间 - 客户端点击时间） = 点击的服务器时间
		timekey := easygo.Get0ClockTimestamp(beginTime)                                        //获取服务器时间0点
		logs.Info("timekey:", unixToStr(timekey), "beginTime:", unixToStr(beginTime),
			"reqMsg.GetClientTime():", unixToStr(reqMsg.GetClientTime()), "v.GetClientClickTime():", unixToStr(v.GetClientClickTime()), "相隔秒：", reqMsg.GetClientTime()-v.GetClientClickTime())
		rm := v.BpsData
		if rm == nil {
			continue
		}
		pageType := rm.GetPageType()
		menuId := rm.GetMenuId()
		labelId := rm.GetLabelId()
		exTabId := rm.GetExTabId()
		dataId := rm.GetDataId()
		exId := rm.GetExId()
		navigationId := rm.GetNavigationId()
		dataType := rm.GetDataType()
		buttonId := rm.GetButtonId()

		key := GetBpsClickIdkey(self.PlayerId, timekey, pageType, menuId, labelId, exTabId, dataId, exId, navigationId, dataType, buttonId)
		strdata := self.GetKeyData(key)
		var data = &RedisESportBpsClick{}
		if strdata != "" && len(strdata) > 0 {
			_ = json.Unmarshal([]byte(strdata), &data)
		} else {
			data = &RedisESportBpsClick{
				MenuId:       menuId,
				PageType:     pageType,
				DataId:       dataId,
				ExId:         exId,
				TimeKey:      timekey, //easygo.GetToday0ClockTimestamp(),
				BeginTime:    beginTime,
				ExTabId:      exTabId,
				LabelId:      labelId,
				ActCount:     0,
				NavigationId: navigationId,
				DataType:     dataType,
				ButtonId:     buttonId,
			}
		}
		data.ActCount++
		s, _ := json.Marshal(data)
		self.SaveKeyData(key, string(s))
	}
}

//保存停留时间
func (self *RedisESportBpsClickLogObj) BpsClickItemToTableStruct(key string, data *RedisESportBpsClick) *share_message.TableESPortsBpsClick {

	pretime := data.BeginTime
	curmId := data.MenuId
	curdId := data.DataId
	cureid := data.ExId
	timekey := data.TimeKey
	log := &share_message.TableESPortsBpsClick{
		Id:           easygo.NewString(key),
		TimeKey:      easygo.NewInt64(timekey),
		PlayerId:     easygo.NewInt64(self.PlayerId),
		CreateTime:   easygo.NewInt64(pretime),
		MenuId:       easygo.NewInt32(curmId),
		LabelId:      easygo.NewInt64(data.LabelId),
		ExTabId:      easygo.NewInt64(data.ExTabId),
		DataId:       easygo.NewInt64(curdId),
		ExId:         easygo.NewInt64(cureid),
		PageType:     easygo.NewInt32(data.PageType),
		ActCount:     easygo.NewInt32(data.ActCount),
		NavigationId: easygo.NewInt32(data.NavigationId),
		DataType:     easygo.NewInt32(data.DataType),
		ButtonId:     easygo.NewInt32(data.ButtonId),
	}
	return log

}

func NewRedisESportBpsClickLogObj(id int64) *RedisESportBpsClickLogObj {
	obj := RedisESportBpsClickLogObj{}
	return obj.Init(id)
}

//对外方法，获取对象，如果为nil表示redis内存不存在，数据库也不存在
func GetRedisESportBpsClickLogObj(id int64) *RedisESportBpsClickLogObj {
	obj, ok := ESportBpsClickMgr.Load(id)
	if ok && obj != nil {
		return obj.(*RedisESportBpsClickLogObj)
	} else {
		return NewRedisESportBpsClickLogObj(id)
	}
}
