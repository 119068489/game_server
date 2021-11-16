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
恋爱匹配埋点报表
*/
const ( //redis内存中存在的key
	BURYINT_POINT_EXIST_TIME = 1000 * 600 //redis的key删除时间:毫秒
)

//埋点事件类型定义
const (
	BP_VC_MAIN_BG_ENTER      int32 = 1001 //进入匹配主页
	BP_VC_MAIN_BG_EXIT       int32 = 1002 //退出匹配主页
	BP_VC_MAIN_BG_CHCEKCARD  int32 = 1003 //查看名片
	BP_VC_MAIN_DJ_SAYHI      int32 = 1004 //点击sayHi
	BP_VC_MAIN_DJ_XIHUAN     int32 = 1005 //点击喜欢
	BP_VC_MAIN_DJ_FRESH      int32 = 1006 //点击刷新
	BP_VC_MAIN_DJ_GRTX       int32 = 1007 //点击个人头像
	BP_VC_MAIN_DJ_LZMP       int32 = 1008 //点击录制名片
	BP_VC_MAIN_DJ_FINISH     int32 = 1009 //点击完成
	BP_VC_MAIN_DJ_MPTX       int32 = 1010 //点击名片头像
	BP_VC_MAIN_DJ_YYY        int32 = 1011 //点击摇一摇
	BP_VC_LZMP_BG_ENTER      int32 = 1012 //进入录制页
	BP_VC_LZMP_BG_EXIT       int32 = 1013 //退出录制页
	BP_VC_LZMP_DJ_EXIT       int32 = 1014 //点击返回
	BP_VC_LZMP_DJ_DUBAI      int32 = 1015 //点击独白
	BP_VC_LZMP_DJ_DYPY       int32 = 1016 //点击电影配音
	BP_VC_LZMP_DJ_CYC        int32 = 1017 //点击唱一唱
	BP_VC_LZMP_DJ_LY         int32 = 1018 //点击录音
	BP_VC_LZMP_DJ_SSGD       int32 = 1019 //点击搜索更多
	BP_VC_LZMP_DJ_SANGCHUAN  int32 = 1020 //点击上传
	BP_VC_LZMP_DJ_QUXIAO     int32 = 1021 //点击取消
	BP_VC_LZMP_DJ_SZJYMP     int32 = 1022 //点击设置交友名片
	BP_VC_SXHW_BG_ENTER      int32 = 1023 //进入谁喜欢我页
	BP_VC_SXHW_BG_EXIT       int32 = 1024 //退出谁喜欢我页
	BP_VC_SXHW_DJ_XHW        int32 = 1025 //点击喜欢我
	BP_VC_SXHW_DJ_WXH        int32 = 1026 //点击我喜欢
	BP_VC_SXHW_XHW_DJ_HF     int32 = 1027 //点击喜欢我页的回复
	BP_VC_SXHW_XHW_DJ_SAYHI  int32 = 1028 //点击喜欢我页的sayHi
	BP_VC_SXHW_XHW_DJ_XIHUAN int32 = 1029 //点击喜欢我页的喜欢
	BP_VC_SXHW_WXH_DJ_BF     int32 = 1030 //点击我喜欢页的播放
	BP_VC_SXHW_WXH_DJ_LT     int32 = 1031 //点击我喜欢页的聊天
	BP_VC_SXHW_WXH_DJ_TX     int32 = 1032 //点击我喜欢页的头像
	BP_VC_SSGD_BG_ENTER      int32 = 1033 //进入搜索更多页
	BP_VC_SSGD_BG_EXIT       int32 = 1034 //退出搜索更多页
	BP_VC_SSGD_DJ_SC         int32 = 1035 //点击搜索更多页的上传
	BP_VC_SSGD_SC_DJ_FH      int32 = 1036 //点击上传页的返回
	BP_VC_SSGD_SC_DJ_TJ      int32 = 1037 //点击上传页的提交
	BP_WISH_ACCESS           int32 = 3000 //访问许愿池
	BP_WISH_EXCHANGE         int32 = 3001 //访问兑换钻石页
	BP_WISH_VEXCHANGE        int32 = 3002 //成功兑换钻石
	BP_WISH_DARE             int32 = 3003 //点击挑战页
	BP_WISH_CHALLENGE        int32 = 3004 // 点击挑战盲盒
)

type VCBuryingPointReportObj struct {
	Id int64
	RedisBase
}
type VCBuryingPointReportEx struct {
	Id                 int64 `json:"_id"`
	MainEnterPeopleNum int64
	MainEnterNum       int64
	MainReadCardNum    int64
	MainSayHiNum       int64
	MainZanNum         int64
	MainFreshNum       int64
	MainHeadNum        int64
	MainRecordNum      int64
	MainFinishNum      int64
	MainHeadCardNum    int64
	MainHandShake      int64
	MainOnLineTime     int64
	LZMPEnterPeopleNum int64
	LZMPEnterNum       int64
	LZMPBackNum        int64
	LZMPDuBai          int64
	LZMPdypy           int64
	LZMPcyc            int64
	LZMPly             int64
	LZMPlyNum          int64
	LZMPssgd           int64
	LZMPqx             int64
	LZMPcg             int64
	LZMPjymp           int64
	LZMPOnlineTime     int64
	SXHWEnterPeopleNum int64
	SXHWxhw            int64
	SXHWwxh            int64
	SXHWxhwHuiFu       int64
	SXHWxhwSayHi       int64
	SXHWxhwZan         int64
	SXHWwxhBoFang      int64
	SXHWwxhChat        int64
	SXHWwxhHead        int64
	SSGDEnterPeopleNum int64
	SSGDsc             int64
	SSGDscBackNum      int64
	SSGDscTJNum        int64
	SSGDOnlineTime     int64
}

func NewRedisVCBuryingPointReportObj(id int64, data ...*share_message.VCBuryingPointReport) *VCBuryingPointReportObj {
	p := &VCBuryingPointReportObj{
		Id: id,
	}
	obj := append(data, nil)[0]
	return p.Init(obj)
}

func (self *VCBuryingPointReportObj) Init(obj *share_message.VCBuryingPointReport) *VCBuryingPointReportObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_VC_BURYING_POINT_REPORT)
	self.Sid = VCBuryingPointMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		VCBuryingPointMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = QueryVCBuryingPointReport(self.Id, 0)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisVCBuryingPointReport(obj)
	}

	logs.Info("初始化新的BuryingPoint管理器:", self.Id)
	return self
}

func (self *VCBuryingPointReportObj) GetId() interface{} { //override
	return self.Id
}
func (self *VCBuryingPointReportObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_VC_BURYING_POINT_REPORT, self.Id)
}
func (self *VCBuryingPointReportObj) UpdateData() { //override
	if !self.IsExistKey() {
		VCBuryingPointMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > BURYINT_POINT_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		VCBuryingPointMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}

func (self *VCBuryingPointReportObj) InitRedis() { //override
	obj := QueryVCBuryingPointReport(self.Id, 0)
	if obj == nil {
		return
	}
	self.SetRedisVCBuryingPointReport(obj)
}

func (self *VCBuryingPointReportObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisVCBuryingPointReport()
	return data
}

func (self *VCBuryingPointReportObj) SaveOtherData() { //override
}

func (self *VCBuryingPointReportObj) SetRedisVCBuryingPointReport(obj *share_message.VCBuryingPointReport) {
	//增加到管理器
	VCBuryingPointMgr.Store(obj.GetId(), self)
	self.AddToExistList(obj.GetId())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	ex := &VCBuryingPointReportEx{}
	StructToOtherStruct(obj, ex)
	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), ex)
	easygo.PanicError(err)

}

func (self *VCBuryingPointReportObj) GetRedisVCBuryingPointReport() *share_message.VCBuryingPointReport {

	obj := &VCBuryingPointReportEx{}
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	easygo.PanicError(err)
	err = redis.ScanStruct(value, obj)
	easygo.PanicError(err)
	newObj := &share_message.VCBuryingPointReport{}
	StructToOtherStruct(obj, newObj)
	return newObj
}

//增加指定值
func (self *VCBuryingPointReportObj) IncrFileVal(file string, val int64) {
	self.IncrOneValue(file, val)
}

//管理器
func GetRedisVCBuryingPointReportObj(id int64) *VCBuryingPointReportObj {
	querytime := easygo.Get0ClockTimestamp(id)
	return VCBuryingPointMgr.GetRedisVCBuryingPointReportObj(querytime)
}

//更新报表
func (self *VCBuryingPointReportObj) UpdateRedisVCBuryingPointReport(obj *share_message.VCBuryingPointReport) {
	ex := &VCBuryingPointReportEx{}
	StructToOtherStruct(obj, ex)
	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), ex)
	easygo.PanicError(err)
}

//查询指定时间的埋点报表
func QueryVCBuryingPointReport(querytime int64, role ...int32) *share_message.VCBuryingPointReport {
	querytime = easygo.Get0ClockTimestamp(querytime)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_VC_BURYING_POINT_REPORT)
	defer closeFun()
	var obj *share_message.VCBuryingPointReport
	err := col.Find(bson.M{"_id": querytime}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//批量保存需要存储的数据
func SaveRedisVCBuryingPointReportToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_VC_BURYING_POINT_REPORT, &ids)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisVCBuryingPointReportObj(id)
		if obj != nil {
			data := obj.GetRedisVCBuryingPointReport()
			saveData = append(saveData, bson.M{"_id": data.GetId()}, data)
			obj.SetSaveStatus(false)
		}
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_VC_BURYING_POINT_REPORT, saveData)
	}
}
