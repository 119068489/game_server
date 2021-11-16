package for_game

import (
	"game_server/easygo"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

/*
文章阅读报表内存数据管理
*/
const (
	ARTICLEREPORT_EXIST_LIST = "articleReport_exist_list" //redis内存中存在的key
	ARTICLEREPORT_EXIST_TIME = 1000 * 600                 //redis的key删除时间:毫秒
)

type ArticleReportObj struct {
	Id int64
	RedisBase
}
type RedisArticleReport struct {
	Id         int64 `json:"_id"`
	CreateTime int64
	PushPlayer int64
	Clicks     int64
	Jumps      int64
}

func NewRedisArticleReport(id int64, data ...*share_message.ArticleReport) *ArticleReportObj {
	p := &ArticleReportObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}

func (self *ArticleReportObj) Init(obj *share_message.ArticleReport) *ArticleReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_ARTICLE_REPORT)
	self.Sid = ArticleReportMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		ArticleReportMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = QueryArticleReport(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisArticleReport(obj)
	}

	logs.Info("初始化新的ArticleReport管理器:", self.Id)
	return self
}
func (self *ArticleReportObj) GetId() interface{} { //override
	return self.Id
}
func (self *ArticleReportObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_ARTICLE_REPORT, self.Id)
}
func (self *ArticleReportObj) UpdateData() { //override
	if !self.IsExistKey() {
		ArticleReportMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > ARTICLEREPORT_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		ArticleReportMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *ArticleReportObj) InitRedis() { //override
	obj := QueryArticleReport(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisArticleReport(obj)
}
func (self *ArticleReportObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisArticleReport()
	return data
}
func (self *ArticleReportObj) SaveOtherData() { //override
}
func (self *ArticleReportObj) SetRedisArticleReport(obj *share_message.ArticleReport) {

	//增加到管理器
	ArticleReportMgr.Store(obj.GetId(), self)
	self.AddToExistList(obj.GetId())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}

	redisObj := &RedisArticleReport{}
	StructToOtherStruct(obj, redisObj)

	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), redisObj)
	easygo.PanicError(err)

}
func (self *ArticleReportObj) GetRedisArticleReport() *share_message.ArticleReport {

	obj := &RedisArticleReport{}
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.ArticleReport{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

//设置RedisArticleReport字段值
func SetRedisArticleReportFildVal(id, val int64, fild string) {
	obj := ArticleReportMgr.GetRedisArticleReport(id)
	obj.IncrOneValue(fild, val)
}

//获取数据
func GetRedisArticleReport(id int64) *ArticleReportObj {
	obj := ArticleReportMgr.GetRedisArticleReport(id)
	return obj
}

//id查询文章报表
func QueryArticleReport(id int64) *share_message.ArticleReport {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ARTICLE_REPORT)
	defer closeFun()
	var obj *share_message.ArticleReport
	err := col.Find(bson.M{"_id": id}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//批量保存需要存储的数据
func SaveRedisArticleReportToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_ARTICLE_REPORT, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisArticleReport(id)
		if obj != nil {
			data := obj.GetRedisArticleReport()
			saveData = append(saveData, bson.M{"_id": data.GetId()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_ARTICLE_REPORT, saveData)
	}
}
