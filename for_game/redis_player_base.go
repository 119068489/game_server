package for_game

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/client_hall"
	"game_server/pb/client_server"
	"game_server/pb/share_message"
	"math/rand"
	"runtime/debug"
	"sort"
	"time"

	"github.com/astaxie/beego/logs"

	//"reflect"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/garyburd/redigo/redis"
)

const (
	//列表key
	PLAYER_PHOTO_LIST     = "Photo"       //玩家相片
	PLAYER_TEAM_IDS       = "TeamIds"     //玩家群id
	PLAYER_FRIEND_IDS     = "FriendList"  //玩家好友id
	PLAYER_BANK_INFO      = "BankInfo"    //银行卡列表
	PLYAER_BLACK_LIST     = "BlackList"   //黑名单列表
	PLYAER_COLLECT_INFO   = "CollectInfo" //收藏的信息
	PLYAER_LABEL_LIST     = "Label"       //标签列表
	PLYAER_CUSTOM_TAG     = "Customtag"   //自定义标签列表
	PLAYER_ATTENTION_LIST = "Attention"   //社交广场关注列表
	PLAYER_DYNAMIC_LIST   = "Dynamic"     //自己发布的动态id列表
	PLAYER_EMOTICON       = "Emoticon"    //表情
	PLAYER_FANS           = "Fans"        //粉丝
	//PLAYER_POINTS         = "Points"      // 玩家的数组坐标.
	PLAYER_POINTS = "Points1" // 玩家的数组坐标.

	//结构key
	PLAYER_SETTING  = "PlayerSetting" //个人设置信息
	PLAYER_CALLINFO = "CallInfo"      //正在拨打电话的信息

	PLAYER_EXIST_LIST = "player_base_exist_list" //redis key值存在列表

	PLAYER_SESSIONS = "Sessions" //会话列表

	PLAYER_PERSONALITY_TAGS = "PersonalityTags" //个性化标签

	PLAYER_EXIST_TIME = 600 * 1000 //毫秒，key值存在时间

)
const (
	DefaultPageSize = 50
)

const (
	PLAYER_SEX_ALL  = 0 // 全部
	PLAYER_SEX_BOY  = 1 // 男
	PLAYER_SEX_GIRL = 2 // 女
)

const (
	NEAR_SORT_DISTANCE = 1 // 距离优先
	NEAR_SORT_ONLINE   = 2 // 在线优先

)

const (
	NEAR_ALL_NUM       = 200    // 附近的人总数 200
	NEAR_OPERATION_NUM = 30     // 运营号30人
	NEAR_DISTANCE      = 100000 // 100公里
	NEAR_NEAR_DISTANCE = 50000  // 附近的人推荐好友50公里以外
	//NEAR_DISTANCE = 4000 // 100公里
	NEAR_OPERATIONAL_DISTANCE = 50 // 50米内
	//NEAR_OPERATIONAL_DISTANCE = 5000 // 10公里
	NEAR_BEGIN_ROBOT_DISTANCE = 10000 // 10公里以外
	//NEAR_BEGIN_ROBOT_DISTANCE = 2000 // 10公里

	NEAR_EXPIRE      = 600 // 10分钟过期
	NEAR_DYNAMIC_NUM = 3   // 3条动态

	// 引导区插入的区间
	NEAR_INSERT_LEAD_BEGIN = 3 // 3开始
	NEAR_INSERT_LEAD_END   = 8 // 8结束

)

const (
	REDIS_NEAR_PLAYER = "redis_near_player"
)

const (
	ONLINE_STATUS_NEW        = 1 // 在线
	ONLINE_STATUS_ONLINE_NEW = 2 // 在线新人
	ONLINE_STATUS_LOG_OUT    = 3 // 刚刚下线
	ONLINE_STATUS_OFFLINE    = 4 // 离线
)

const (
	NEAR_INFO_DATE_TYPE_NOMAL = 1 //   普通信息
	NEAR_INFO_DATE_TYPE_LEAD  = 2 // 引导区域信息
)

const (
	NEAR_DEFAULT_LNG = 113.922867 // x 深圳南山
	NEAR_DEFAULT_LAT = 22.53661   // y
)

// 附近的人内容类型
const (
	NEAR_CONTENT_TYPE_COMMENT = 1 // 1- 普通内容; 2-图片或表情
	NEAR_CONTENT_TYPE_EMOTI   = 2 // 2-图片或表情
)

type IRedisPlayerBase interface {
	GetClientEndpoint() IClientEndpoint
	GetMyInfo() *client_server.PlayerMsg
}
type RedisPlayerBaseObj struct {
	Id PLAYER_ID //玩家id
	//非数据库存储字段
	LastOSMTime int64 //上次获取短信验证码时间
	RedisBase
}

//玩家基本信息
type PlayerBaseEx struct {
	PlayerId              int64   `json:"_id"` //玩家id
	Password              string  //密码
	NickName              string  //昵称
	HeadIcon              string  //头像
	Sex                   int32   //性别 1男 2女
	Gold                  int64   //玩家携带的金币
	IsRobot               bool    //是否是机器人
	LastOnLineTime        int64   //最后上线时间
	Email                 string  //邮箱
	PeopleId              string  //身份证id
	Account               string  //闲聊号
	Phone                 string  //手机号
	CreateTime            int64   //注册时间
	IsOnline              bool    //是否在线
	RealName              string  //真实姓名
	Signature             string  //个性签名
	Provice               string  //省
	City                  string  //市
	Area                  string  //区
	IsRecommend           bool    //是否推荐
	LastLogOutTime        int64   //最后下线时间
	LoginTimes            int32   //登录次数
	OnlineTime            int64   //在线时长
	X                     float64 //X坐标
	Y                     float64 //Y坐标
	DeviceType            int32   //设备类型 1 IOS，2 Android，3 PC
	IsNearBy              bool    //是否有附近的人的打招呼消息
	Channel               string  //渠道 注册来源
	Types                 int32   //用户类型 1普通用户，2客服用户
	LastAssistantTime     int64   //最后一次取畅聊助手的时间
	WXOpenId              string
	WXSessionKey          string
	WXUnionid             string
	ClearLocalLogTime     int64  //清除本地聊天记录时间
	ComplaintTime         int64  //上次投诉意见的时间
	IsVisitor             bool   //是否是游客
	Sid                   int32  //所在大厅sid
	TodayOnlineTime       int64  //当日在线时长
	CreateIP              string //注册IP
	LastLoginIP           string //最后登陆IP
	ApiUrl                string //用户后台api地址
	SecretKey             string //用户后台apiKey
	FreeTimes             int32  //提现免手续费次数:默认初始3次
	AutoLoginTime         int64  //自动登录设置时间
	AutoLoginToken        string //自动登录的token
	Status                int32  //用户状态 0 正常 1冻结
	DeviceCode            string //最后登录设备码
	IsRecommendOver       bool   //是否选择了推荐好友和群信息
	Note                  string //备注
	GrabTag               int32  //抓取标签
	Token                 string //登陆token
	PayPassword           string //支付密码
	Version               string //当前客户端版本
	Brand                 string //登陆设备
	Zan                   int64  //社交广场被点赞数
	IsReloadSquare        bool   //是否加载过自己的社交广场数据
	VerifiedTime          int64  //实名认证时间
	Coin                  int64  //硬币
	BCoin                 int64  //绑定硬币:活动或者赠送
	AreaCode              string //电话号码区号
	BackgroundImageURL    string // 我的主页的背景图
	CheckNum              int32  //过审广场动态数量 大于50 自动审核 违规重置为0
	IsBrowse2Square       bool   //是否连续浏览2层屏内容,如果是false,后端不处理,直接查询,如果是true,后端埋点
	FirstAddSquareDynamic bool   //该用户是否已第一次发布动态
	IsCheckChatLog        bool   //聊天记录查询白名单
	YoungPassWord         string //青少年模式
	IsLoadedAllSessions   bool   //是否加载过全部会话
	IsCanRoam             bool   //是否漫游

	VCZanNum          int64  //名片被点赞数
	BgImageUrl        string //背景图
	MixId             int64  //录音ID
	Constellation     int32  //星座ID
	BsVCZanNum        int64  //后台名片点赞数
	ReadLoveMeLogTime int64  //读取我喜欢列表的时间
	ConstellationTime int64  //修改星座的时间
	ESportCoin        int64  //电竞币
	Operator          string //修改操作人
	BanOverTime       int64  //封禁结束时间
}

//对外方法=============================================
func NewRedisPlayerBase(playerId int64, data ...*share_message.PlayerBase) *RedisPlayerBaseObj {
	p := &RedisPlayerBaseObj{
		Id: playerId,
	}
	obj := append(data, nil)[0]

	return p.Init(obj)
}

func (self *RedisPlayerBaseObj) Init(obj *share_message.PlayerBase) *RedisPlayerBaseObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	self.Sid = PlayerBaseMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		PlayerBaseMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = self.QueryPlayerBase(self.Id)
			if obj == nil {
				return nil
			}
		}
		self.SetRedisPlayerBase(obj)
	}
	//	logs.Info("初始化新的PlayerBase管理器:", self.Id)
	return self
}
func (self *RedisPlayerBaseObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisPlayerBaseObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_PLAYER_BASE, self.Id)
}

//定时更新数据
func (self *RedisPlayerBaseObj) UpdateData() { //override
	if !self.IsExistKey() {
		PlayerBaseMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > PLAYER_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		PlayerBaseMgr.Delete(self.Id) // 释放对象
		return
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}

func (self *RedisPlayerBaseObj) InitRedis() { //override
	obj := self.QueryPlayerBase(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisPlayerBase(obj)
}
func (self *RedisPlayerBaseObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisPlayerBase()
	return data
}
func (self *RedisPlayerBaseObj) SaveOtherData() { //override
	//保存表情数据
	self.SaveEmoticons()
}

//通过playerId从mongo中读取登录玩家数据
func (self *RedisPlayerBaseObj) QueryPlayerBase(id PLAYER_ID) *share_message.PlayerBase {
	data := self.QueryMongoData(id)
	if data != nil {
		var player share_message.PlayerBase
		StructToOtherStruct(data, &player)
		return &player
	}
	//redis中获取数据初始化
	if self.IsExistKey() {
		return self.GetRedisPlayerBase()
	}
	return nil
}
func (self *RedisPlayerBaseObj) SetRedisPlayerBase(obj *share_message.PlayerBase) {
	PlayerBaseMgr.Store(obj.GetPlayerId(), self)
	self.AddToExistList(obj.GetPlayerId())
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	//玩家基础信息
	player := &PlayerBaseEx{}
	StructToOtherStruct(obj, player)
	err := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), player)
	easygo.PanicError(err)
	//其他玩家数据设置到redis
	self.SetPlayerOtherData(obj, false)
	//好友数据
	self.SetPlayerFriendsData(nil, false)
	//表情数据
	self.SetPlayerEmoticonData(nil, false)

}

//设置玩家好友数据
func (self *RedisPlayerBaseObj) SetPlayerFriendsData(ids []int64, save ...bool) {
	if ids == nil {
		friend := GetFriendBase(self.Id)
		if friend != nil {
			ids = friend.GetFriendIds()
		}
	}
	self.SetStringValueToRedis(PLAYER_FRIEND_IDS, ids, save...)
}

//设置玩家表情数据
func (self *RedisPlayerBaseObj) SetPlayerEmoticonData(emots []*share_message.PlayerEmoticon, save ...bool) {
	if emots == nil {
		emots = GetEmoticonFromMongo(self.Id)
	}
	self.SetStringValueToRedis(PLAYER_EMOTICON, emots, save...)
}

//玩家其他数据写入redis
func (self *RedisPlayerBaseObj) SetPlayerOtherData(obj *share_message.PlayerBase, save ...bool) {
	//	logs.Info("设置玩家其他数据")
	self.SetStringValueToRedis(PLAYER_SETTING, obj.GetPlayerSetting(), save...)
	self.SetStringValueToRedis(PLAYER_PHOTO_LIST, obj.GetPhoto(), save...)
	self.SetStringValueToRedis(PLAYER_TEAM_IDS, obj.GetTeamIds(), save...)
	self.SetStringValueToRedis(PLAYER_ATTENTION_LIST, obj.GetAttenList(), save...)
	self.SetStringValueToRedis(PLAYER_FANS, obj.GetFansList(), save...)
	self.SetStringValueToRedis(PLAYER_DYNAMIC_LIST, obj.GetDynamicList(), save...)
	self.SetStringValueToRedis(PLAYER_BANK_INFO, obj.GetBankInfo(), save...)
	self.SetStringValueToRedis(PLYAER_BLACK_LIST, obj.GetBlackList(), save...)
	self.SetStringValueToRedis(PLYAER_COLLECT_INFO, obj.GetCollectInfo(), save...)
	self.SetStringValueToRedis(PLAYER_CALLINFO, obj.GetCallInfo(), save...)
	self.SetStringValueToRedis(PLYAER_LABEL_LIST, obj.GetLabel(), save...)
	self.SetStringValueToRedis(PLYAER_CUSTOM_TAG, obj.GetCustomTag(), save...)
	if obj.Points1 == nil {
		points1 := &share_message.GeoJson{
			Type:        easygo.NewString("Point"),
			Coordinates: []float64{obj.GetX(), obj.GetY()},
		}
		obj.Points1 = points1
	}
	self.SetStringValueToRedis(PLAYER_POINTS, obj.GetPoints1(), save...)
	self.SetStringValueToRedis(PLAYER_PERSONALITY_TAGS, obj.GetPersonalityTags(), save...)

}

//获取玩家其他redis数据
func (self *RedisPlayerBaseObj) GetPlayerOtherData(obj *share_message.PlayerBase) {
	self.GetStringValueToRedis(PLAYER_SETTING, &obj.PlayerSetting)
	self.GetStringValueToRedis(PLAYER_PHOTO_LIST, &obj.Photo)
	self.GetStringValueToRedis(PLAYER_TEAM_IDS, &obj.TeamIds)
	self.GetStringValueToRedis(PLAYER_ATTENTION_LIST, &obj.AttenList)
	self.GetStringValueToRedis(PLAYER_FANS, &obj.FansList)
	self.GetStringValueToRedis(PLAYER_DYNAMIC_LIST, &obj.DynamicList)
	self.GetStringValueToRedis(PLAYER_BANK_INFO, &obj.BankInfo)
	self.GetStringValueToRedis(PLYAER_BLACK_LIST, &obj.BlackList)
	self.GetStringValueToRedis(PLYAER_COLLECT_INFO, &obj.CollectInfo)
	self.GetStringValueToRedis(PLAYER_CALLINFO, &obj.CallInfo)
	self.GetStringValueToRedis(PLYAER_LABEL_LIST, &obj.Label)
	self.GetStringValueToRedis(PLYAER_CUSTOM_TAG, &obj.CustomTag)
	if obj.Points1 == nil {
		geoJson := &share_message.GeoJson{
			Type:        easygo.NewString("Point"),
			Coordinates: []float64{obj.GetX(), obj.GetY()},
		}
		obj.Points1 = geoJson
		//obj.Points = []float64{obj.GetX(), obj.GetY()}
	}
	//数组去重
	if obj.AttenList != nil {
		obj.AttenList = RemoveRepeatedElementInt64(obj.GetAttenList())
	}
	if obj.FansList != nil {
		obj.FansList = RemoveRepeatedElementInt64(obj.GetFansList())
	}
	if obj.DynamicList != nil {
		obj.DynamicList = RemoveRepeatedElementInt64(obj.GetDynamicList())
	}
	if obj.Photo != nil {
		obj.Photo = RemoveRepeatedElementString(obj.GetPhoto())
	}
	if obj.CustomTag != nil {
		obj.CustomTag = RemoveRepeatedElementInt32(obj.GetCustomTag())
	}
	self.GetStringValueToRedis(PLAYER_POINTS, &obj.Points1)
	self.GetStringValueToRedis(PLAYER_PERSONALITY_TAGS, &obj.PersonalityTags)
}

//获取玩家信息
func (self *RedisPlayerBaseObj) GetRedisPlayerBase() *share_message.PlayerBase {
	value, err := easygo.RedisMgr.GetC().HGetAll(self.GetKeyId())
	if len(value) == 0 {
		return nil
	}
	easygo.PanicError(err)
	var base PlayerBaseEx
	err = redis.ScanStruct(value, &base)
	easygo.PanicError(err)
	var newBase *share_message.PlayerBase
	StructToOtherStruct(base, &newBase)
	//其他数据
	self.GetPlayerOtherData(newBase)
	if newBase.GetPlayerId() == 0 {
		debug.PrintStack()
		return nil
	}
	return newBase
}

//修改表情包
func (self *RedisPlayerBaseObj) ModifyRedisEmoticon(emoticon *share_message.PlayerEmoticon) bool {
	emots := self.GetEmoticons()
	isNew := true
	for _, v := range emots {
		if v.GetTypeId() == emoticon.GetTypeId() {
			isNew = false
			v.Data = easygo.NewString(emoticon.GetData())
			v.IsSave = easygo.NewBool(true)
			break
		}
	}
	if isNew {
		//新增
		emoticon.PlayerId = easygo.NewInt64(self.Id)
		emoticon.Id = easygo.NewInt64(NextId(TABLE_PLAYER_EMOTICON))
		emoticon.IsSave = easygo.NewBool(true)
		emots = append(emots, emoticon)
	}
	self.SetPlayerEmoticonData(emots)
	return true
}

//删除表情包
func (self *RedisPlayerBaseObj) DelRedisEmoticon(id int32) bool {
	emots := self.GetEmoticons()
	for i, v := range emots {
		if v.GetTypeId() == id {
			emots = append(emots[:i], emots[i+1:]...)
			break
		}
	}
	self.SetPlayerEmoticonData(emots)
	//直接删除数据库
	DelEmoticonFromMongo(self.Id, id)
	return true
}

//增加好友
func (self *RedisPlayerBaseObj) AddRedisPlayerFriends(id int64) {
	ids := self.GetFriends()
	logs.Info("添加好友前:", self.Id, ids)
	ids = append(ids, id)
	self.SetPlayerFriendsData(ids)
	logs.Info("添加好友后:", self.Id, ids)
}

//删除好友
func (self *RedisPlayerBaseObj) DelRedisPlayerFriends(id int64) {
	ids := self.GetFriends()
	ids = easygo.Del(ids, id).([]int64)
	self.SetPlayerFriendsData(ids)
}

//删除关注id
func (self *RedisPlayerBaseObj) DelAttention(id int64) {
	ids := self.GetAttention()
	ids = easygo.Del(ids, id).([]int64)
	self.SetStringValueToRedis(PLAYER_ATTENTION_LIST, ids)
}

//增加关注列表
func (self *RedisPlayerBaseObj) AddAttention(id int64) {
	lst := self.GetAttention()
	if util.Int64InSlice(id, lst) {
		return
	}
	lst = append(lst, id)
	self.SetStringValueToRedis(PLAYER_ATTENTION_LIST, lst)
}

//删除粉丝id
func (self *RedisPlayerBaseObj) DelFans(id int64) {
	fans := self.GetFans()
	fans = easygo.Del(fans, id).([]int64)
	self.SetStringValueToRedis(PLAYER_FANS, fans)
}

//增加粉丝列表
func (self *RedisPlayerBaseObj) AddFans(id int64) {
	lst := self.GetFans()
	if util.Int64InSlice(id, lst) {
		return
	}
	lst = append(lst, id)
	self.SetStringValueToRedis(PLAYER_FANS, lst)
}

// 判断操作者不是自己,
// 判断操作者的关注类表是否有self.id,如果有,又是第一页,则置顶
// 如果操作者是自己,是第一页,置顶
// GetRedisPlayerDynamicListByPage 获取动态id列表
//func (self *RedisPlayerBaseObj) GetRedisPlayerDynamicListByPage(opId int64, page, pageSize int32, logId ...int64) map[string]interface{} {
//	dynamicIdList := self.GetDynamicList()      // 未排序
//	easygo.SortSliceInt64(dynamicIdList, false) // 降序排序
//	var inAttentions bool
//	if self.GetPlayerId() != opId {
//		// 获得当前操作者,获取当前操作者的关注列表,判断self是否在关注列表中
//		opPlayer := GetRedisPlayerBase(opId)
//		opAttentions := opPlayer.GetAttention()
//		inAttentions = util.Int64InSlice(self.GetPlayerId(), opAttentions)
//	} else {
//		inAttentions = true
//	}
//	// 排除所有置顶的id
//	newIds := make([]int64, 0) // 没有置顶的动态id列表,前面几条是置顶的条数
//	// 如果是第一页,封装置顶消息在前面几条
//	if page == 1 && inAttentions {
//
//		ofTopIds := self.GetTopByKeyNameFromRedis(REDIS_SQUARE_OF_TOP_DYNAMIC) // 官方置顶
//		bsTopIds := self.GetTopByKeyNameFromRedis(REDIS_SQUARE_BS_TOP_DYNAMIC) // 后台置顶
//		appTopIds := self.GetTopByKeyNameFromRedis(REDIS_SQUARE_TOP_DYNAMIC)   // app置顶
//		if len(ofTopIds) > 0 {
//			// 时间最新的最靠前
//			sortSlice := SortDynamicSliceByTime(ofTopIds)
//			newIds = append(newIds, sortSlice...)
//		}
//		if len(bsTopIds) > 0 {
//			// 时间最新的最靠前
//			sortSlice := SortDynamicSliceByTime(bsTopIds)
//			newIds = append(newIds, sortSlice...)
//		}
//		if len(appTopIds) > 0 {
//			// 时间最新的最靠前
//			sortSlice := SortDynamicSliceByTime(appTopIds)
//			newIds = append(newIds, sortSlice...)
//		}
//		// 取前3条置顶消息
//		if len(newIds) > PLAYER_TOP_NUM {
//			newIds = newIds[0:PLAYER_TOP_NUM]
//		}
//	}
//	noTopIds := make([]int64, 0) // 倒序排序使用
//	// 获取该用户的所有置顶动态id列表
//	allTopIds := self.GetPlayerAllTopIds()
//	noTopIds = util.Slice1DelSlice2(dynamicIdList, allTopIds)
//	if len(noTopIds) == 0 {
//		noTopIds = append(noTopIds, dynamicIdList...)
//	}
//	easygo.SortSliceInt64(noTopIds, false) //降序排序
//	newIds = append(newIds, noTopIds...)
//	// 分页
//	paginatorMap := Paginator(page, pageSize, newIds)
//	return paginatorMap
//}

// 判断操作者不是自己,
// 判断操作者的关注类表是否有self.id,如果有,又是第一页,则置顶
// 如果操作者是自己,是第一页,置顶
// GetRedisPlayerDynamicListByPage 获取动态id列表
func (self *RedisPlayerBaseObj) GetRedisPlayerDynamicListByPage(opId int64, page, pageSize int32) ([]*share_message.DynamicData, int) {
	//noTopDynamicIdList, count := GetNoTopDynamicByPIDsFromDB(int(page), int(pageSize), []int64{self.GetPlayerId()})
	maxLogIdKey := MakeNewString(opId, self.GetPlayerId())
	noTopDynamicIdList, count := GetNoTopDynamicByPIDs(opId, int(page), int(pageSize), []int64{self.GetPlayerId()}, maxLogIdKey)
	var inAttentions bool
	if self.GetPlayerId() != opId {
		// 获得当前操作者,获取当前操作者的关注列表,判断self是否在关注列表中
		opPlayer := GetRedisPlayerBase(opId)
		opAttentions := opPlayer.GetAttention()
		inAttentions = util.Int64InSlice(self.GetPlayerId(), opAttentions)
	} else {
		inAttentions = true
	}
	dsList := make([]*share_message.DynamicData, 0)
	// 如果是第一页,封装置顶消息在前面几条
	if page == 1 && inAttentions {
		//获取后台置顶的动态列表(in)
		bsTopDynamicList := GetBSTopDynamicListByIDsFromDB(opId, []int64{self.GetPlayerId()})
		if len(bsTopDynamicList) > 0 {
			slice := GetDynamicSliceByRandFromSlice(bsTopDynamicList, PLAYER_TOP_NUM)
			//// 时间最新的最靠前
			sortSlice := SortDynamicSliceByTime1(slice)
			dsList = append(dsList, sortSlice...)
		}

		//获取app置顶的动态列表(in)
		appTopDynamicList := GetAppTopDynamicListByIDsFromDB(opId, []int64{self.GetPlayerId()})
		if len(appTopDynamicList) > 0 {
			slice := GetDynamicSliceByRandFromSlice(appTopDynamicList, APP_TOP_NUM)
			//// 时间最新的最靠前
			sortSlice := SortDynamicSliceByTime1(slice)
			dsList = append(dsList, sortSlice...)
		}
		// 取前3条置顶消息
		if len(dsList) > PLAYER_TOP_NUM {
			dsList = dsList[0:PLAYER_TOP_NUM]
		}
	}
	dsList = append(dsList, noTopDynamicIdList...)
	return dsList, count
}

//增加动态id列表
func (self *RedisPlayerBaseObj) AddRedisPlayerDynamicList(id int64) {
	ids := self.GetDynamicList()
	if util.Int64InSlice(id, ids) {
		return
	}
	ids = append(ids, id)
	self.SetStringValueToRedis(PLAYER_DYNAMIC_LIST, ids)
}

//删除动态id
func (self *RedisPlayerBaseObj) DelRedisPlayerDynamicList(id int64) {
	ids := self.GetDynamicList()
	ids = easygo.Del(ids, id).([]int64)
	self.SetStringValueToRedis(PLAYER_DYNAMIC_LIST, ids)
}

//删除指定群
func (self *RedisPlayerBaseObj) DelTeamId(id int64) {
	ids := self.GetTeamIds()
	ids = easygo.Del(ids, id).([]int64)
	self.SetStringValueToRedis(PLAYER_TEAM_IDS, ids)
}

//增加群id
func (self *RedisPlayerBaseObj) AddTeamId(id int64) {
	teamIds := self.GetTeamIds()
	if util.Int64InSlice(id, teamIds) {
		return
	}
	teamIds = append(teamIds, id)
	self.SetStringValueToRedis(PLAYER_TEAM_IDS, teamIds)
}

//清空银行卡
func (self *RedisPlayerBaseObj) ClearBankId() {
	self.SetStringValueToRedis(PLAYER_BANK_INFO, []*share_message.BankInfo{})
}

//添加银行卡
func (self *RedisPlayerBaseObj) AddBankInfo(bank *share_message.BankInfo) {
	bankInfos := self.GetBankInfos()
	bankInfos = append(bankInfos, bank)
	self.SetStringValueToRedis(PLAYER_BANK_INFO, bankInfos)
}

//删除银行卡
func (self *RedisPlayerBaseObj) DelBankInfo(bankId string) {
	bankInfos := self.GetBankInfos()
	for i, b := range bankInfos {
		if b.GetBankId() == bankId {
			bankInfos = append(bankInfos[:i], bankInfos[i+1:]...) //easygo.Del(bankInfos, b).([]*share_message.BankInfo)
			break
		}
	}
	self.SetStringValueToRedis(PLAYER_BANK_INFO, bankInfos)
}

//获取某张银行卡信息
func (self *RedisPlayerBaseObj) GetBankInfo(bankId string) *share_message.BankInfo {
	bankInfos := self.GetBankInfos()
	for _, b := range bankInfos {
		if b.GetBankId() == bankId {
			return b
		}
	}
	return nil
}

//修改卡信息
func (self *RedisPlayerBaseObj) ModifyBankInfo(reqMsg *client_hall.BankMessage) {
	bankInfos := self.GetBankInfos()
	for _, b := range bankInfos {
		if b.GetBankId() == reqMsg.GetBankCardNo() {
			b.Provice = reqMsg.Provice
			b.City = reqMsg.City
			b.Area = reqMsg.Area
			break
		}
	}
	self.SetStringValueToRedis(PLAYER_BANK_INFO, bankInfos)
}

func (self *RedisPlayerBaseObj) SetRedisBlackList(blackList []int64) {
	self.SetStringValueToRedis(PLYAER_BLACK_LIST, blackList)
}

//增加收藏
func (self *RedisPlayerBaseObj) AddCollectInfo(collect *share_message.CollectInfo) *share_message.CollectInfo {
	var num int32
	collects := self.GetCollectInfo()
	if len(collects) == 0 {
		num = 1
	} else {
		num = collects[len(collects)-1].GetIndex() + 1
	}

	collect.Index = easygo.NewInt32(num)
	collect.Time = easygo.NewInt64(GetMillSecond())
	collects = append(collects, collect)
	self.SetStringValueToRedis(PLYAER_COLLECT_INFO, collects)
	return collect
}
func (self *RedisPlayerBaseObj) DelCollectInfo(index int32) {
	collects := self.GetCollectInfo()
	for i, msg := range collects {
		if index == msg.GetIndex() {
			collects = append(collects[:i], collects[i+1:]...)
			//collects = easygo.Del(collects, msg).([]*share_message.CollectInfo)
			self.SetStringValueToRedis(PLYAER_COLLECT_INFO, collects)
			break
		}
	}
}

//修改兴趣标签列表
func (self *RedisPlayerBaseObj) SetRedisLabelList(list []int32) {
	self.SetStringValueToRedis(PLYAER_LABEL_LIST, list)
}

//修改自定义标签列表
func (self *RedisPlayerBaseObj) SetRedisCustomTag(list []int32) {
	self.SetStringValueToRedis(PLYAER_CUSTOM_TAG, list)
	self.SaveOneRedisDataToMongo(PLYAER_CUSTOM_TAG, list)
}

//获取创建时间
func (self *RedisPlayerBaseObj) GetCreateTime() int64 {
	var val int64
	self.GetOneValue("CreateTime", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetCreateTime() {
	t := GetMillSecond()
	self.SetOneValue("CreateTime", t)
}

func (self *RedisPlayerBaseObj) GetIsLoadedAllSessions() bool {
	var val bool
	self.GetOneValue("IsLoadedAllSessions", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetIsLoadedAllSessions(b bool) {
	self.SetOneValue("IsLoadedAllSessions", b)
	self.SaveOneRedisDataToMongo("IsLoadedAllSessions", b)
}

func (self *RedisPlayerBaseObj) GetIsCanRoam() bool {
	var val bool
	self.GetOneValue("IsCanRoam", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetIsCanRoam(b bool) {
	self.SetOneValue("IsCanRoam", b)
	self.SaveOneRedisDataToMongo("IsCanRoam", b)
}

//更新支付密码
func (self *RedisPlayerBaseObj) SetPayPassword(pass string) {
	pwd := Md5(pass)
	self.SetOneValue("PayPassword", pwd)
}
func (self *RedisPlayerBaseObj) GetPayPassword() string {
	var val string
	self.GetOneValue("PayPassword", &val)
	return val
}

//获取最后登出时间
func (self *RedisPlayerBaseObj) GetLastLogOutTime() int64 {
	var val int64
	self.GetOneValue("LastLogOutTime", &val)
	return val
}

//设置最后登出时间
func (self *RedisPlayerBaseObj) SetLastLogOutTime(v int64) {
	self.SetOneValue("LastLogOutTime", v)
}

//获取当前客户端版本号
func (self *RedisPlayerBaseObj) GetVersion() string {
	var val string
	self.GetOneValue("Version", &val)
	return val
}

//设置当前客户端版本号
func (self *RedisPlayerBaseObj) SetVersion(v string) {
	self.SetOneValue("Version", v)
}

//获取当前客户端登陆型号
func (self *RedisPlayerBaseObj) GetBrand() string {
	var val string
	self.GetOneValue("Version", &val)
	return val
}

//设置当前客户端登陆型号
func (self *RedisPlayerBaseObj) SetBrand(v string) {
	self.SetOneValue("Brand", v)
}

//获取最后登录时间
func (self *RedisPlayerBaseObj) GetLastOnLineTime() int64 {
	var val int64
	self.GetOneValue("LastOnLineTime", &val)
	return val
}

//设置最后登录时间
func (self *RedisPlayerBaseObj) SetLastOnLineTime(v int64) {
	self.SetOneValue("LastOnLineTime", v)
}

//获取
func (self *RedisPlayerBaseObj) GetPhone() string {
	var val string
	self.GetOneValue("Phone", &val)
	return val
}

//设置手机号
func (self *RedisPlayerBaseObj) SetPhone(v string) {
	self.SetOneValue("Phone", v)
}

//获取昵称
func (self *RedisPlayerBaseObj) GetNickName() string {
	var val string
	self.GetOneValue("NickName", &val)
	return val
}

//设置昵称
func (self *RedisPlayerBaseObj) SetNickName(v string) {
	//存时进行base64
	self.SetOneValue("NickName", v)
	self.SaveOneRedisDataToMongo("NickName", v)
}

//获取账号
func (self *RedisPlayerBaseObj) GetAccount() string {
	var val string
	self.GetOneValue("Account", &val)
	return val
}

//设置账号
func (self *RedisPlayerBaseObj) SetAccount(v string) {
	self.SetOneValue("Account", v)
}

//是否设置了支付密码
func (self *RedisPlayerBaseObj) GetIsPayPassword() bool {
	if self.GetPayPassword() == "" {
		return false
	}
	return true
}

//是否设置登录密码
func (self *RedisPlayerBaseObj) GetIsLoginPassword() bool {
	account := GetRedisAccountObj(self.GetPlayerId())
	if account.GetPassword() == "" {
		return false
	}
	return true
}

//获取头像
func (self *RedisPlayerBaseObj) GetHeadIcon() string {
	var val string
	self.GetOneValue("HeadIcon", &val)
	return val
}

//设置头像
func (self *RedisPlayerBaseObj) SetHeadIcon(v string) {
	self.SetOneValue("HeadIcon", v)
	self.SaveOneRedisDataToMongo("HeadIcon", v)
}

//获取API地址
func (self *RedisPlayerBaseObj) GetApiUrl() string {
	var val string
	self.GetOneValue("ApiUrl", &val)
	return val
}

//设置API地址
func (self *RedisPlayerBaseObj) SetApiUrl(v string) {
	self.SetOneValue("ApiUrl", v)
}

//获取API KEY
func (self *RedisPlayerBaseObj) GetSecretKey() string {
	var val string
	self.GetOneValue("SecretKey", &val)
	return val
}

//设置API KEY
func (self *RedisPlayerBaseObj) SetSecretKey(v string) {
	self.SetOneValue("SecretKey", v)
}

// 设置背景图
func (self *RedisPlayerBaseObj) SetBackgroundImageURL(v string) {
	self.SetOneValue("BackgroundImageURL", v)
}

// 获取背景图
func (self *RedisPlayerBaseObj) GetBackgroundImageURL() string {
	var val string
	self.GetOneValue("BackgroundImageURL", &val)
	return val
}

//获取性别
func (self *RedisPlayerBaseObj) GetSex() int32 {
	var val int32
	self.GetOneValue("Sex", &val)
	return val
}

//设置性别
func (self *RedisPlayerBaseObj) SetSex(v int32) {
	self.SetOneValue("Sex", v)
}

//获取前端清除日记时间
func (self *RedisPlayerBaseObj) GetClearLocalLogTime() int64 {
	var val int64
	self.GetOneValue("ClearLocalLogTime", &val)
	return val
}

//设置前端清除日记时间
func (self *RedisPlayerBaseObj) SetClearLocalLogTime() {
	t := GetMillSecond()
	self.SetOneValue("ClearLocalLogTime", t)
}

//获取上次投诉时间
func (self *RedisPlayerBaseObj) GetComplaintTime() int64 {
	var val int64
	self.GetOneValue("ComplaintTime", &val)
	return val
}

//设置上次投诉时间
func (self *RedisPlayerBaseObj) SetComplaintTime() {
	t := util.GetMilliTime()
	self.SetOneValue("ComplaintTime", t)
}

//获取真实姓名
func (self *RedisPlayerBaseObj) GetRealName() string {
	var val string
	self.GetOneValue("RealName", &val)
	return val
}

//设置真实姓名
func (self *RedisPlayerBaseObj) SetRealName(v string) {
	self.SetOneValue("RealName", v)
}

//获取身份证
func (self *RedisPlayerBaseObj) GetPeopleId() string {
	var val string
	self.GetOneValue("PeopleId", &val)
	return val
}

//设置身份证
func (self *RedisPlayerBaseObj) SetPeopleId(v string) {
	self.SetOneValue("PeopleId", v)
	self.SaveOneRedisDataToMongo("PeopleId", v)
}

//获取实名认证时间
func (self *RedisPlayerBaseObj) GetAuthTime() int64 {
	var val int64
	self.GetOneValue("VerifiedTime", &val)
	return val
}

//设置实名认证时间
func (self *RedisPlayerBaseObj) SetAuthTime(v int64) {
	self.SetOneValue("VerifiedTime", v)
	self.SaveOneRedisDataToMongo("VerifiedTime", v)
}

//获取状态
func (self *RedisPlayerBaseObj) GetStatus() int32 {
	var val int32
	self.GetOneValue("Status", &val)
	return val
}

//修改状态
func (self *RedisPlayerBaseObj) SetStatus(v int32) {
	self.SetOneValue("Status", v)
}

//获取点赞总数
func (self *RedisPlayerBaseObj) GetZan() int64 {
	var val int64
	self.GetOneValue("Zan", &val)
	return val
}

func (self *RedisPlayerBaseObj) AddZan(num int64) int64 {
	val := self.GetZan()
	val += num
	self.SetOneValue("Zan", val)
	return val
}

////修改点赞总数
//func (self *RedisPlayerBaseObj) SetZan(v int64) {
//	self.SetOneValue("Zan", v)
//}

//获取备注
func (self *RedisPlayerBaseObj) GetNote() string {
	var val string
	self.GetOneValue("Note", &val)
	return val
}

//修改备注
func (self *RedisPlayerBaseObj) SetNote(v string) {
	self.SetOneValue("Note", v)
}

//获取邮箱
func (self *RedisPlayerBaseObj) GetEmail() string {
	var val string
	self.GetOneValue("Email", &val)
	return val
}

//修改邮箱
func (self *RedisPlayerBaseObj) SetEmail(v string) {
	self.SetOneValue("Email", v)
}

//获取个性签名
func (self *RedisPlayerBaseObj) GetSignature() string {
	var val string
	self.GetOneValue("Signature", &val)
	if val == "" {
		val = "这个人很懒，什么都没留下!"
		self.SetSignature(val)
	}
	return val
}

//修改个性签名
func (self *RedisPlayerBaseObj) SetSignature(v string) {
	self.SetOneValue("Signature", v)
}

//获取省份
func (self *RedisPlayerBaseObj) GetProvice() string {
	var val string
	self.GetOneValue("Provice", &val)
	return val
}

//修改省份
func (self *RedisPlayerBaseObj) SetProvice(v string) {
	self.SetOneValue("Provice", v)
}

//获取城市
func (self *RedisPlayerBaseObj) GetCity() string {
	var val string
	self.GetOneValue("City", &val)
	return val
}

//修改城市
func (self *RedisPlayerBaseObj) SetCity(v string) {
	self.SetOneValue("City", v)
}

//获取所在区
func (self *RedisPlayerBaseObj) GetArea() string {
	var val string
	self.GetOneValue("Area", &val)
	return val
}

//修改所在区
func (self *RedisPlayerBaseObj) SetArea(v string) {
	self.SetOneValue("Area", v)
}

//修改电话号码区号
func (self *RedisPlayerBaseObj) SetAreaCode(v string) {
	self.SetOneValue("AreaCode", v)
}

//获取电话号码区号
func (self *RedisPlayerBaseObj) GetAreaCode() string {
	var val string
	self.GetOneValue("AreaCode", &val)
	return val
}

//获取在线状态
func (self *RedisPlayerBaseObj) GetIsOnLine() bool {
	var val bool
	self.GetOneValue("IsOnline", &val)
	return val
}

//设置在线状态
func (self *RedisPlayerBaseObj) SetIsOnLine(v bool) {
	self.SetOneValue("IsOnline", v)
}

//获取极光id
func (self *RedisPlayerBaseObj) GetRegistrationId() string {
	var val string
	self.GetOneValue("RegistrationId", &val)
	return val
}

//修改极光id
func (self *RedisPlayerBaseObj) SetRegistrationId(v string) {
	self.SetOneValue("RegistrationId", v)
}

//获取渠道
func (self *RedisPlayerBaseObj) GetChannel() string {
	var val string
	self.GetOneValue("Channel", &val)
	return val
}

//设置渠道
func (self *RedisPlayerBaseObj) SetChannel(v string) {
	self.SetOneValue("Channel", v)
}

//获取金币
func (self *RedisPlayerBaseObj) GetGold() int64 {
	var val int64
	self.GetOneValue("Gold", &val)
	return val
}

//设置金币
func (self *RedisPlayerBaseObj) IncrGold(gold int64) int64 {
	return self.IncrOneValue("Gold", gold)
}

//获取过审动态数量方法
func (self *RedisPlayerBaseObj) GetCheckNum() int64 {
	var val int64
	self.GetOneValue("CheckNum", &val)
	return val
}

//设置过审动态数量方法
func (self *RedisPlayerBaseObj) IncrCheckNum(num int64) int64 {
	return self.IncrOneValue("CheckNum", num)
}

//需要通知玩家钱变化的话，请勿使用这个
func (self *RedisPlayerBaseObj) TryAddGold(value int64, reason string, sourceType int32, extendLog interface{}) string {
	return self.AddGold(value, reason, sourceType, extendLog, true)
}
func (self *RedisPlayerBaseObj) AddGold(value int64, reason string, sourceType int32, extendLog interface{}, try ...bool) string {
	newGold := self.IncrGold(value)
	if newGold < 0 {
		//减少钱加回来
		self.IncrGold(-value)
		b := append(try, false)[0]
		if b {
			return "零钱不足"
		}
	}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": self.GetPlayerId()}, bson.M{"$inc": bson.M{"Gold": value}})
	easygo.PanicError(err)
	//self.SetGold(newGold)
	self.AddGoldLog(value, reason, sourceType, newGold-value, newGold, extendLog)
	return ""
}
func (self *RedisPlayerBaseObj) AddGoldLog(value int64, reason string, sourceType int32, oldGold int64, newGold int64, extendLog interface{}) {
	if value == 0 {
		return
	}
	st := GettSourceTypeById(sourceType)
	logId := NextId(TABLE_GOLDCHANGELOG)
	log := share_message.GoldChangeLog{
		LogId:      &logId,
		PlayerId:   easygo.NewInt64(self.Id),
		ChangeGold: &value,
		SourceType: &sourceType,
		PayType:    st.Type,
		Note:       &reason,

		CurGold:    &oldGold,
		Gold:       &newGold,
		CreateTime: easygo.NewInt64(GetMillSecond()),
	}
	m := &CommonGold{
		GoldChangeLog: log,
		Extend:        extendLog,
	}
	AddGoldChangeLog(m)
	//col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_GOLDCHANGELOG)
	//defer closeFun()
	//err := col.Insert(m)
	//easygo.PanicError(err)

	fun := func() {
		MakeInOutCashSumReport(sourceType, value)
	}
	easygo.Spawn(fun) //写出入款汇总报表
}

//获取硬币
func (self *RedisPlayerBaseObj) GetAllCoin() int64 {
	var coin, bCoin int64
	self.GetOneValue("Coin", &coin)
	self.GetOneValue("BCoin", &bCoin)
	return coin + bCoin
}

func (self *RedisPlayerBaseObj) GetCoin() int64 {
	var coin int64
	self.GetOneValue("Coin", &coin)
	return coin
}
func (self *RedisPlayerBaseObj) GetBCoin() int64 {
	var coin int64
	self.GetOneValue("BCoin", &coin)
	return coin
}

//设置硬币
func (self *RedisPlayerBaseObj) IncrCoin(coin int64) int64 {
	return self.IncrOneValue("Coin", coin)
}
func (self *RedisPlayerBaseObj) IncrBCoin(coin int64) int64 {
	return self.IncrOneValue("BCoin", coin)
}

//绑定硬币获得记录
func (self *RedisPlayerBaseObj) AddBCoinLog(sourceType int32, coin int64) {
	logId := NextId(TABLE_PLAYER_BCOIN_LOG)
	t := time.Now().Unix()
	log := share_message.PlayerBCoinLog{
		Id:         easygo.NewInt64(logId),
		PlayerId:   easygo.NewInt64(self.Id),
		Way:        easygo.NewInt32(sourceType),
		BCoin:      easygo.NewInt64(coin),
		CurBCoin:   easygo.NewInt64(coin),
		CreateTime: easygo.NewInt64(t),
		OverTime:   easygo.NewInt64(t + COIN_EXPIRATION_TIME),
		Status:     easygo.NewInt32(BCOIN_STATUS_UNUSE),
	}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BCOIN_LOG)
	defer closeFun()
	err := col.Insert(log)
	easygo.PanicError(err)
}

//对绑定硬币状态修改
func (self *RedisPlayerBaseObj) UpdateBCoinLog(coin int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BCOIN_LOG)
	defer closeFun()
	var logs []*share_message.PlayerBCoinLog
	err := col.Find(bson.M{"PlayerId": self.Id, "Status": BCOIN_STATUS_UNUSE}).Sort("OverTime").All(&logs)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	var saveLogs []interface{}
	for _, log := range logs {
		if coin < 0 {
			if log.GetCurBCoin()+coin >= 0 {
				log.CurBCoin = easygo.NewInt64(log.GetCurBCoin() + coin)
				saveLogs = append(saveLogs, bson.M{"_id": log.GetId()}, log)
				break
			}
			coin += log.GetCurBCoin()
			log.Status = easygo.NewInt32(BCOIN_STATUS_USED)
			log.CurBCoin = easygo.NewInt64(0)
			saveLogs = append(saveLogs, bson.M{"_id": log.GetId()}, log)
		}
	}
	if len(saveLogs) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_BCOIN_LOG, saveLogs)
	}
}

//硬币分绑定和不绑定，非绑定的由兑换获得，其他均为绑定的，且有时效，一定时间后过去
//isUseCoin 是否使用非绑定硬币:true使用非绑定硬币，false不使用
func (self *RedisPlayerBaseObj) AddCoin(value int64, reason string, sourceType int32, extendLog interface{}, isUseCoin ...bool) string {
	var bCoin, coin, newBCoin, newCoin int64
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	if value > 0 {
		//获得硬币
		if sourceType == COIN_TYPE_EXCHANGE_IN {
			//充值获得的
			//redis增加
			newCoin = self.IncrCoin(value)
			coin = value
			newBCoin = self.GetBCoin()
		} else {
			//其他渠道获得硬币加到bcoin上
			newBCoin = self.IncrBCoin(value)
			bCoin = value
			newCoin = self.GetCoin()
			//TODO 写入新表记录过期时间
			self.AddBCoinLog(sourceType, value)
		}
	} else {
		isUse := append(isUseCoin, true)[0]
		switch sourceType {
		case COIN_TYPE_SYSTEM_OUT: //后台回收绑定硬币
			newBCoin = self.IncrBCoin(value)
			if newBCoin < 0 {
				//不使用非绑定硬币，绑定硬币直接扣完就行
				newBCoin = self.IncrBCoin(-newBCoin) //绑定硬币变成0
				bCoin = value - newBCoin
			} else {
				bCoin = value
			}
			newCoin = self.GetCoin()
		case COIN_TYPE_CONFISCATE_OUT: //后台回收充值硬币
			newCoin = self.IncrCoin(value)
			if newCoin < 0 {
				newCoin = self.IncrCoin(-newCoin)
				coin = value - newCoin
			} else {
				coin = value
			}
			newBCoin = self.GetBCoin()
		default:
			//消费硬币，优先优先消费绑定硬币，再消费非绑定的
			newBCoin = self.IncrBCoin(value)
			if newBCoin < 0 {
				//绑定硬币不够，扣非绑定
				if isUse {
					newCoin = self.IncrCoin(newBCoin)
					if newCoin < 0 {
						//非绑定也不够，则消费失败
						//减少的钱加回来
						self.IncrCoin(-newBCoin)
						self.IncrBCoin(-value)
						return "硬币不足"
					}
					bCoin = value - newBCoin
					coin = newBCoin
					newBCoin = self.IncrBCoin(-newBCoin) //绑定硬币变成0
				} else {
					//不使用非绑定硬币，绑定硬币直接扣完就行
					newBCoin = self.IncrBCoin(-newBCoin) //绑定硬币变成0
					bCoin = value - newBCoin
					newCoin = self.GetCoin()
				}
			} else {
				bCoin = value
				newCoin = self.GetCoin()
			}
		}
		if bCoin != 0 && isUse {
			self.UpdateBCoinLog(bCoin)
		}
	}
	//绑定硬币减少
	if bCoin != 0 {
		_, err := col.Upsert(bson.M{"_id": self.GetPlayerId()}, bson.M{"$inc": bson.M{"BCoin": bCoin}})
		easygo.PanicError(err)
	}
	//非绑定硬币减少
	if coin != 0 {
		_, err := col.Upsert(bson.M{"_id": self.GetPlayerId()}, bson.M{"$inc": bson.M{"Coin": coin}})
		easygo.PanicError(err)
	}
	//写流水日志
	self.AddCoinLog(value, reason, sourceType, newCoin-coin, newCoin, newBCoin-bCoin, newBCoin, extendLog)
	return ""
}
func (self *RedisPlayerBaseObj) AddCoinLog(value int64, reason string, sourceType int32, oldCoin, newCoin, oldBCoin, newBCoin int64, extendLog interface{}) {
	if value == 0 {
		return
	}
	st := GettSourceTypeById(sourceType)
	logId := NextId(TABLE_COINCHANGELOG)
	log := share_message.CoinChangeLog{
		LogId:      easygo.NewInt64(logId),
		PlayerId:   easygo.NewInt64(self.Id),
		ChangeCoin: easygo.NewInt64(value),
		SourceType: easygo.NewInt32(sourceType),
		PayType:    st.Type,
		Note:       easygo.NewString(reason),
		CurCoin:    easygo.NewInt64(oldCoin),
		Coin:       easygo.NewInt64(newCoin),
		CurBCoin:   easygo.NewInt64(oldBCoin),
		BCoin:      easygo.NewInt64(newBCoin),
		CreateTime: easygo.NewInt64(GetMillSecond()),
	}
	m := &CommonCoin{
		CoinChangeLog: log,
		Extend:        extendLog,
	}
	AddCoinChangeLog(m)

	//fun := func() {
	//	MakeInOutCashSumReport(sourceType, value)
	//}
	//easygo.Spawn(fun) //写出入款汇总报表
}

//获取电竞币
func (self *RedisPlayerBaseObj) GetESportCoin() int64 {
	var val int64
	self.GetOneValue("ESportCoin", &val)
	return val
}

//设置电竞币
func (self *RedisPlayerBaseObj) IncrESportCoin(eSportCoin int64) int64 {
	return self.IncrOneValue("ESportCoin", eSportCoin)
}

//需要通知玩家电竞币变化的话，请勿使用这个
func (self *RedisPlayerBaseObj) TryAddESportCoin(value int64, reason string, sourceType int32, extendLog interface{}) string {
	return self.AddESportCoin(value, reason, sourceType, extendLog, true)
}
func (self *RedisPlayerBaseObj) AddESportCoin(value int64, reason string, sourceType int32, extendLog interface{}, try ...bool) string {
	newESportCoin := self.IncrESportCoin(value)
	if newESportCoin < 0 {
		//减少电竞币加回来
		self.IncrESportCoin(-value)
		b := append(try, false)[0]
		if b {
			return "电竞币不足"
		}
	}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": self.GetPlayerId()}, bson.M{"$inc": bson.M{"ESportCoin": value}})
	easygo.PanicError(err)
	self.AddESportCoinLog(value, reason, sourceType, newESportCoin-value, newESportCoin, extendLog)
	return ""
}
func (self *RedisPlayerBaseObj) AddESportCoinLog(value int64, reason string, sourceType int32, oldESportCoin int64, newESportCoin int64, extendLog interface{}) {
	if value == 0 {
		return
	}
	st := GettSourceTypeById(sourceType)
	logId := NextId(TABLE_ESPORTCHANGELOG)
	log := share_message.ESportCoinChangeLog{
		LogId:            easygo.NewInt64(logId),
		PlayerId:         easygo.NewInt64(self.Id),
		ChangeESportCoin: easygo.NewInt64(value),
		SourceType:       easygo.NewInt32(sourceType),
		PayType:          st.Type,
		Note:             easygo.NewString(reason),

		CurESportCoin: easygo.NewInt64(oldESportCoin),
		ESportCoin:    easygo.NewInt64(newESportCoin),
		CreateTime:    easygo.NewInt64(GetMillSecond()),
	}
	m := &CommonESportCoin{
		ESportCoinChangeLog: log,
		Extend:              extendLog,
	}
	AddESportCoinChangeLog(m)
}

//是否游客
func (self *RedisPlayerBaseObj) GetIsVisitor() bool {
	var val bool
	self.GetOneValue("IsVisitor", &val)
	return val
}

//提现免密次数
func (self *RedisPlayerBaseObj) GetFreeTimes() int32 {
	var val int32
	self.GetOneValue("FreeTimes", &val)
	return val
}

func (self *RedisPlayerBaseObj) GetIsBindWechat() bool {
	account := GetRedisAccountObj(self.GetPlayerId())
	if account.GetOpenId() == "" {
		return false
	}
	return true
}

//是否是附近人
func (self *RedisPlayerBaseObj) GetIsNearBy() bool {
	var val bool
	self.GetOneValue("IsNearBy", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetIsNearBy(val bool) {
	self.SetOneValue("IsNearBy", val)
}

//是否是附近人
func (self *RedisPlayerBaseObj) GetIsRecommendOver() bool {
	var val bool
	self.GetOneValue("IsRecommendOver", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetIsRecommendOver(b bool) {
	self.SetOneValue("IsRecommendOver", b)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	base := &share_message.PlayerBase{IsRecommendOver: easygo.NewBool(b)}
	err := col.Update(bson.M{"_id": self.Id}, bson.M{"$set": base})
	if err != nil {
		logs.Error("err:", err, base)
	}
}

//提现免密次数
func (self *RedisPlayerBaseObj) GetPlayerId() int64 {
	return self.Id
}

//获取密码
func (self *RedisPlayerBaseObj) GetPassword() string {
	var val string
	self.GetOneValue("Password", &val)
	return val
}
func (self *RedisPlayerBaseObj) GetLastAssistantTime() int64 {
	var val int64
	self.GetOneValue("LastAssistantTime", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetLastAssistantTime(val int64) {
	self.SetOneValue("LastAssistantTime", val)
}
func (self *RedisPlayerBaseObj) GetDeviceType() int32 {
	var val int32
	self.GetOneValue("DeviceType", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetDeviceType(val int32) {
	self.SetOneValue("DeviceType", val)
	self.SaveOneRedisDataToMongo("DeviceType", val)
}
func (self *RedisPlayerBaseObj) SetApkCode(val int32) {
	self.SetOneValue("ApkCode", val)
}
func (self *RedisPlayerBaseObj) GetApkCode() int32 {
	var val int32
	self.GetOneValue("ApkCode", &val)
	return val
}
func (self *RedisPlayerBaseObj) GetIsRobot() bool {
	var val bool
	self.GetOneValue("IsRobot", &val)
	return val
}
func (self *RedisPlayerBaseObj) UpdateLogoutTimestamp() {
	time := GetMillSecond()
	self.SetLastLogOutTime(time)
}
func (self *RedisPlayerBaseObj) SetOnlineTime(val int64) {
	self.SetOneValue("OnlineTime", val)
}
func (self *RedisPlayerBaseObj) GetSid() int32 {
	var sid int32
	self.GetOneValue("Sid", &sid)
	return sid
}
func (self *RedisPlayerBaseObj) SetSid(val int32) {
	self.SetOneValue("Sid", val)
	self.SaveToMongo()
}
func (self *RedisPlayerBaseObj) GetOnlineTime() int64 {
	var val int64
	self.GetOneValue("OnlineTime", &val)
	return val
}
func (self *RedisPlayerBaseObj) GetTodayOnlineTime() int64 {
	var val int64
	self.GetOneValue("TodayOnlineTime", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetTodayOnlineTime(val int64) {
	self.SetOneValue("TodayOnlineTime", val)
}
func (self *RedisPlayerBaseObj) UpdateOnlineTime() {
	times := self.GetLastLogOutTime() - self.GetLastOnLineTime()
	self.SetOnlineTime(self.GetOnlineTime() + times)
	var day_start int64 = easygo.GetToday0ClockTimestamp() * 1000
	var todayOnline int64
	if self.GetLastOnLineTime() > day_start {
		todayOnline = self.GetTodayOnlineTime() + times
	} else {
		todayOnline = self.GetLastLogOutTime() - day_start
	}
	self.SetTodayOnlineTime(todayOnline)
}
func (self *RedisPlayerBaseObj) UpdateLastLoginIP(ip string) {
	self.SetOneValue("LastLoginIP", ip)
}
func (self *RedisPlayerBaseObj) GetLastLoginIP() string {
	var val string
	self.GetOneValue("LastLoginIP", &val)
	return val
}
func (self *RedisPlayerBaseObj) GetLoginTimes() int32 {
	var val int32
	self.GetOneValue("LoginTimes", &val)
	return val
}
func (self *RedisPlayerBaseObj) AddLoginTimes() {
	self.IncrOneValue("LoginTimes", 1)
}
func (self *RedisPlayerBaseObj) SetX(val float64) {
	self.SetOneValue("X", val)

}
func (self *RedisPlayerBaseObj) SetY(val float64) {
	self.SetOneValue("Y", val)
}
func (self *RedisPlayerBaseObj) SetStation(x, y float64) {
	self.SetX(x)
	self.SetY(y)
	//self.SetPoints([]float64{x, y})
	geoJson := &share_message.GeoJson{
		Type:        easygo.NewString("Point"),
		Coordinates: []float64{x, y},
	}
	self.SetPoints(geoJson)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()

	//base := &share_message.PlayerBase{X: easygo.NewFloat64(x), Y: easygo.NewFloat64(y), Points: []float64{x, y}}
	base := &share_message.PlayerBase{X: easygo.NewFloat64(x), Y: easygo.NewFloat64(y), Points1: geoJson}
	err := col.Update(bson.M{"_id": self.Id}, bson.M{"$set": base})
	easygo.PanicError(err)
}
func (self *RedisPlayerBaseObj) SetAutoLoginToken(val string) {
	self.SetOneValue("AutoLoginToken", val)
}

func (self *RedisPlayerBaseObj) GetAutoLoginToken() string {
	var val string
	self.GetOneValue("AutoLoginToken", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetAutoLoginTime(val int64) {
	self.SetOneValue("AutoLoginTime", val)
}
func (self *RedisPlayerBaseObj) GetAutoLoginTime() int64 {
	var val int64
	self.GetOneValue("AutoLoginTime", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetDeviceCode(val string) {
	self.SetOneValue("DeviceCode", val)
}

func (self *RedisPlayerBaseObj) GetDeviceCode() string {
	var val string
	self.GetOneValue("DeviceCode", &val)
	return val
}

func (self *RedisPlayerBaseObj) GetGrabTag() int32 {
	var val int32
	self.GetOneValue("GrabTag", &val)
	return val
}

//设置抓取词标签
func (self *RedisPlayerBaseObj) SetGrabTag(v int32) {
	self.SetOneValue("GrabTag", v)
	self.SaveOneRedisDataToMongo("GrabTag", v)
}

///========================playersetting=============================

//设置安全码
func (self *RedisPlayerBaseObj) SetSafePassword(v string) {
	setting := self.GetPlayerSetting()
	setting.SafePassword = easygo.NewString(Md5(v))
	if v == "" {
		setting.IsSafePassword = easygo.NewBool(false)
	} else {
		setting.IsSafePassword = easygo.NewBool(true)
	}
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

func (self *RedisPlayerBaseObj) GetSafePassword() string {
	setting := self.GetPlayerSetting()
	return setting.GetSafePassword()
}

//获取是否设置了安全码
//func (self *RedisPlayerBaseObj) SetIsSafePassword(v bool) {
//	self.SetSettingValue("IsSafePassword", v)
//}
//
//新消息通知开关
func (self *RedisPlayerBaseObj) SetIsNewMessage(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsNewMessage = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

//新消息通知开关
func (self *RedisPlayerBaseObj) GetIsNewMessage() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsNewMessage()
}

//声音开关
func (self *RedisPlayerBaseObj) SetIsMusic(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsMusic = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

//震动开关
func (self *RedisPlayerBaseObj) SetIsShake(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsShake = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

//好友认证
func (self *RedisPlayerBaseObj) SetIsAddFriend(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsAddFriend = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

func (self *RedisPlayerBaseObj) GetIsAddFriend() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsAddFriend()
}

//手机号添加
func (self *RedisPlayerBaseObj) SetIsPhone(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsPhone = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

//手机号
func (self *RedisPlayerBaseObj) GetIsPhone() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsPhone()
}

//闲聊号添加
func (self *RedisPlayerBaseObj) SetIsAccount(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsAccount = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

func (self *RedisPlayerBaseObj) GetIsAccount() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsAccount()
}

//群聊添加
func (self *RedisPlayerBaseObj) SetIsTeamChat(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsTeamChat = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}
func (self *RedisPlayerBaseObj) GetIsTeamChat() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsTeamChat()
}

//二维码添加
func (self *RedisPlayerBaseObj) SetIsCode(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsCode = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

func (self *RedisPlayerBaseObj) GetIsCode() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsCode()
}

//名片添加
func (self *RedisPlayerBaseObj) SetIsCard(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsCard = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

func (self *RedisPlayerBaseObj) GetIsCard() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsCard()
}

//安全防护
func (self *RedisPlayerBaseObj) SetIsSafeProtect(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsSafeProtect = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

//消息预览
func (self *RedisPlayerBaseObj) SetIsMessageShow(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsMessageShow = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

func (self *RedisPlayerBaseObj) GetIsMessageShow() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsMessageShow()
}

// 是否开启社交广场
func (self *RedisPlayerBaseObj) SetIsOpenSquare(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsOpenSquare = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

func (self *RedisPlayerBaseObj) GetIsOpenSquare() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsOpenSquare()
}

// 是否开启点赞或评论我的动态开关
func (self *RedisPlayerBaseObj) SetIsOpenZanOrComment(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsOpenZanOrComment = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

func (self *RedisPlayerBaseObj) GetIsOpenZanOrComment() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsOpenZanOrComment()
}

// 是否开启回复我的评论开关
func (self *RedisPlayerBaseObj) SetIsOpenRecoverComment(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsOpenRecoverComment = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

func (self *RedisPlayerBaseObj) GetIsOpenRecoverComment() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsOpenRecoverComment()
}

// 是否开启我关注的人发布新动态开关
func (self *RedisPlayerBaseObj) SetIsOpenMyAttention(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsOpenMyAttention = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

func (self *RedisPlayerBaseObj) GetIsOpenMyAttention() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsOpenMyAttention()
}

// 是否开启人气动态推荐
func (self *RedisPlayerBaseObj) SetIsOpenRecommend(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsOpenRecommend = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

func (self *RedisPlayerBaseObj) GetIsOpenRecommend() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsOpenRecommend()
}

//是否禁止陌生人打招呼
func (self *RedisPlayerBaseObj) GetIsBanSayHi() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsBanSayHi()
}
func (self *RedisPlayerBaseObj) SetIsBanSayHi(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsBanSayHi = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

// 是否开启硬币商城
func (self *RedisPlayerBaseObj) SetIsOpenCoinShop(v bool) {
	setting := self.GetPlayerSetting()
	setting.IsOpenCoinShop = easygo.NewBool(v)
	self.SetStringValueToRedis(PLAYER_SETTING, setting)
}

func (self *RedisPlayerBaseObj) GetIsOpenCoinShop() bool {
	setting := self.GetPlayerSetting()
	return setting.GetIsOpenCoinShop()
}

//更新token码
func (self *RedisPlayerBaseObj) SetToken(token string) {
	self.SetOneValue("Token", token)
	self.SaveOneRedisDataToMongo("Token", token)
}
func (self *RedisPlayerBaseObj) GetToken() string {
	var val string
	self.GetOneValue("Token", &val)
	return val
}

func (self *RedisPlayerBaseObj) SetIsRelodSquare(b bool) {
	self.SetOneValue("IsReloadSquare", b)
}

func (self *RedisPlayerBaseObj) GetIsRelodSquare() bool {
	var val bool
	self.GetOneValue("IsReloadSquare", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetIsBrowse2Square(b bool) {
	self.SetOneValue("IsBrowse2Square", b)
}

func (self *RedisPlayerBaseObj) GetIsBrowse2Square() bool {
	var val bool
	self.GetOneValue("IsBrowse2Square", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetFirstAddSquareDynamic(b bool) {
	self.SetOneValue("FirstAddSquareDynamic", b)
}

func (self *RedisPlayerBaseObj) GetFirstAddSquareDynamic() bool {
	var val bool
	self.GetOneValue("FirstAddSquareDynamic", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetIsCheckChatLog(b bool) {
	self.SetOneValue("IsCheckChatLog", b)
}

func (self *RedisPlayerBaseObj) GetIsCheckChatLog() bool {
	var val bool
	self.GetOneValue("IsCheckChatLog", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetYoungPassWord(val string) {
	self.SetOneValue("YoungPassWord", val)
}

func (self *RedisPlayerBaseObj) GetYoungPassWord() string {
	var val string
	self.GetOneValue("YoungPassWord", &val)
	return val
}

/*func (self *RedisPlayerBaseObj) SetPoints(points []float64) {
	self.SetStringValueToRedis(PLAYER_POINTS, points)

}*/
func (self *RedisPlayerBaseObj) SetPoints(geoJson *share_message.GeoJson) {
	self.SetStringValueToRedis(PLAYER_POINTS, geoJson)

}

//func (self *RedisPlayerBaseObj) GeRedisPoints() []float64 {
//	var val []float64
//	self.GetStringValueToRedis(PLYAER_CUSTOM_TAG, &val)
//	return val
//}
func (self *RedisPlayerBaseObj) GeRedisPoints() *share_message.GeoJson {
	logs.Info("报错的用户id为----->", self.GetPlayerId())
	var val *share_message.GeoJson
	self.GetStringValueToRedis(PLAYER_POINTS, &val)
	return val
}

func (self *RedisPlayerBaseObj) GetTypes() int32 {
	var val int32
	self.GetOneValue("Types", &val)
	return val
}
func (self *RedisPlayerBaseObj) SetTypes(t int32) {
	self.SetOneValue("Types", t)
}

//获取玩家个人设置
func (self *RedisPlayerBaseObj) GetPlayerSetting() *share_message.PlayerSetting {
	var val *share_message.PlayerSetting
	self.GetStringValueToRedis(PLAYER_SETTING, &val)
	return val
}

//获取玩家相册
func (self *RedisPlayerBaseObj) GetPhoto() []string {
	var val []string
	self.GetStringValueToRedis(PLAYER_PHOTO_LIST, &val)
	return val
}

//设置相册
func (self *RedisPlayerBaseObj) SetPhoto(val []string) {
	self.SetStringValueToRedis(PLAYER_PHOTO_LIST, val)
}

//获取玩家群列表
func (self *RedisPlayerBaseObj) GetTeamIds() []int64 {
	var val []int64
	self.GetStringValueToRedis(PLAYER_TEAM_IDS, &val)
	return val
}

//获取玩家关注列表
func (self *RedisPlayerBaseObj) GetAttention() []int64 {
	var val []int64
	self.GetStringValueToRedis(PLAYER_ATTENTION_LIST, &val)
	return val
}
func (self *RedisPlayerBaseObj) GetFans() []int64 {
	var val []int64
	self.GetStringValueToRedis(PLAYER_FANS, &val)
	return val
}
func (self *RedisPlayerBaseObj) GetDynamicList() []int64 {
	var val []int64
	self.GetStringValueToRedis(PLAYER_DYNAMIC_LIST, &val)
	return val
}
func (self *RedisPlayerBaseObj) GetBankInfos() []*share_message.BankInfo {
	var val []*share_message.BankInfo
	self.GetStringValueToRedis(PLAYER_BANK_INFO, &val)
	return val
}
func (self *RedisPlayerBaseObj) GetFriends() []int64 {
	var val []int64
	self.GetStringValueToRedis(PLAYER_FRIEND_IDS, &val)
	return val
}
func (self *RedisPlayerBaseObj) GetBlackList() []int64 {
	var val []int64
	self.GetStringValueToRedis(PLYAER_BLACK_LIST, &val)
	return val
}
func (self *RedisPlayerBaseObj) GetCollectInfo() []*share_message.CollectInfo {
	var val []*share_message.CollectInfo
	self.GetStringValueToRedis(PLYAER_COLLECT_INFO, &val)
	return val
}
func (self *RedisPlayerBaseObj) GetLabelList() []int32 {
	var val []int32
	self.GetStringValueToRedis(PLYAER_LABEL_LIST, &val)
	return val
}
func (self *RedisPlayerBaseObj) GetCustomTag() []int32 {
	var val []int32
	self.GetStringValueToRedis(PLYAER_CUSTOM_TAG, &val)
	return val
}
func (self *RedisPlayerBaseObj) GetEmoticons() []*share_message.PlayerEmoticon {
	var val []*share_message.PlayerEmoticon
	self.GetStringValueToRedis(PLAYER_EMOTICON, &val)
	return val
}

//声音名片点赞数
func (self *RedisPlayerBaseObj) GetVCZanNum() int64 {
	var val int64
	self.GetOneValue("VCZanNum", &val)
	return val
}

//增加声音名片点赞数
func (self *RedisPlayerBaseObj) AddVCZanNum(num int64) {
	self.IncrOneValue("VCZanNum", num)
	//self.SaveOneRedisDataToMongo("VCZanNum", new)
}

//设置声音名片点赞数
func (self *RedisPlayerBaseObj) SetVCZanNum(num int64) {
	self.SetOneValue("VCZanNum", num)
}

//声音名片背景图
func (self *RedisPlayerBaseObj) GetBgImageUrl() string {
	var val string
	self.GetOneValue("BgImageUrl", &val)
	return val
}

//设置声音名片背景图
func (self *RedisPlayerBaseObj) SetBgImageUrl(val string) {
	self.SetOneValue("BgImageUrl", val)
}

//声音名片录音ID
func (self *RedisPlayerBaseObj) GetMixId() int64 {
	var val int64
	self.GetOneValue("MixId", &val)
	return val
}

//设置声音名片录音ID
func (self *RedisPlayerBaseObj) SetMixId(val int64) {
	oldId := self.GetMixId()
	if oldId != 0 {
		ModifyMyMixVideoUseing(oldId, false)
	}
	self.SetOneValue("MixId", val)
	self.SaveOneRedisDataToMongo("MixId", val)
	if val != 0 {
		ModifyMyMixVideoUseing(val, true)
	}
}

//声音名片录音url
func (self *RedisPlayerBaseObj) GetMixVoiceUrl() string {
	id := self.GetMixId()
	data := GetVoiceCardInfo(id)
	return data.GetMixVoiceUrl()
}

//删除名片
func (self *RedisPlayerBaseObj) DelVoiceCard(mixId int64) int64 {
	ModifyMyMixVideoStatus(mixId, VC_STATUS_DELETE)
	if mixId == self.GetMixId() {
		newMix := GetNewMyMixVideo(self.Id)
		if newMix != nil {
			self.SetMixId(newMix.GetId())
		} else { //删除最后一张
			self.SetMixId(0)
		}
	}
	return self.GetMixId()
}

//个人星座
func (self *RedisPlayerBaseObj) GetConstellation() int32 {
	var val int32
	self.GetOneValue("Constellation", &val)
	return val
}
func (self *RedisPlayerBaseObj) GetConstellationStr() string {
	id := self.GetConstellation()
	return GetPlayerConstellationStr(id)

}

//设置个人星座
func (self *RedisPlayerBaseObj) SetConstellation(val int32) {
	self.SetOneValue("Constellation", val)
}

//后台卡片虚拟点赞
func (self *RedisPlayerBaseObj) GetBsVCZanNum() int64 {
	var val int64
	self.GetOneValue("BsVCZanNum", &val)
	return val
}

//后台卡片虚拟点赞
func (self *RedisPlayerBaseObj) SetBsVCZanNum(val int64) {
	self.SetOneValue("BsVCZanNum", val)
}

//读取喜欢我列表的时间
func (self *RedisPlayerBaseObj) GetReadLoveMeLogTime() int64 {
	var val int64
	self.GetOneValue("ReadLoveMeLogTime", &val)
	return val
}

//设置喜欢我列表的时间
func (self *RedisPlayerBaseObj) SetReadLoveMeLogTime(val int64) {
	self.SetOneValue("ReadLoveMeLogTime", val)
}

//获取修改星座时间
func (self *RedisPlayerBaseObj) GetConstellationTime() int64 {
	var val int64
	self.GetOneValue("ConstellationTime", &val)
	return val
}

//设置修改星座时间
func (self *RedisPlayerBaseObj) SetConstellationTime(val int64) {
	self.SetOneValue("ConstellationTime", val)
}

//设置修改操作人
func (self *RedisPlayerBaseObj) SetOperator(val string) {
	self.SetOneValue("Operator", val)
}

//获取封禁结束时间
func (self *RedisPlayerBaseObj) GetBanOverTime() int64 {
	var val int64
	self.GetOneValue("BanOverTime", &val)
	return val
}

//设置封禁结束时间
func (self *RedisPlayerBaseObj) SetBanOverTime(val int64) {
	self.SetOneValue("BanOverTime", val)
}

//获取个性化标签
func (self *RedisPlayerBaseObj) GetPersonalityTags() []int32 {
	var val []int32
	self.GetStringValueToRedis(PLAYER_PERSONALITY_TAGS, &val)
	return val
}
func (self *RedisPlayerBaseObj) GetPersonalityTagsStr() []*client_hall.PersonTag {
	tags := self.GetPersonalityTags()
	dbTags := GetPlayerPersonalityTags(tags)
	msg := make([]*client_hall.PersonTag, 0)
	for _, d := range dbTags {
		tag := &client_hall.PersonTag{
			Id:   easygo.NewInt32(d.GetId()),
			Name: easygo.NewString(d.GetName()),
		}
		msg = append(msg, tag)
	}
	return msg
}

//设置个性化标签
func (self *RedisPlayerBaseObj) SetPersonalityTags(tags []int32) {
	self.SetStringValueToRedis(PLAYER_PERSONALITY_TAGS, tags)
}

func (self *RedisPlayerBaseObj) AddPlayerSession(id string) {
	b, _ := easygo.RedisMgr.GetC().HExists(self.GetKeyId(), PLAYER_SESSIONS)
	sessions := make([]string, 0)
	if !b {
		playerSession := GetMySessions(self.Id)
		sessions = append(playerSession.GetSessionIds(), id)
	} else {
		self.GetStringValueToRedis(PLAYER_SESSIONS, &sessions)
		sessions = append(sessions, id)
	}
	self.SetStringValueToRedis(PLAYER_SESSIONS, sessions)
	SaveMySessions(self.Id, sessions)
}
func (self *RedisPlayerBaseObj) SavePlayerSessions() {
	b, _ := easygo.RedisMgr.GetC().HExists(self.GetKeyId(), PLAYER_SESSIONS)
	if !b {
		return
	}
	sessions := make([]string, 0)
	self.GetStringValueToRedis(PLAYER_SESSIONS, &sessions)
	SaveMySessions(self.Id, sessions)
}

//获取玩家会话列表
func (self *RedisPlayerBaseObj) GetPlayerSessions() []string {
	b, _ := easygo.RedisMgr.GetC().HExists(self.GetKeyId(), PLAYER_SESSIONS)
	if !b {
		playerSession := GetMySessions(self.Id)
		self.SetStringValueToRedis(PLAYER_SESSIONS, playerSession.GetSessionIds())
		return playerSession.GetSessionIds()
	}
	sessions := make([]string, 0)
	self.GetStringValueToRedis(PLAYER_SESSIONS, &sessions)
	return sessions
}

//删除玩家会话列表
func (self *RedisPlayerBaseObj) DeletePlayerSessions(ids []string) {
	b, _ := easygo.RedisMgr.GetC().HExists(self.GetKeyId(), PLAYER_SESSIONS)
	sessions := make([]string, 0)
	if b {
		self.GetStringValueToRedis(PLAYER_SESSIONS, &sessions)
		for _, id := range ids {
			sessions = easygo.Del(sessions, id).([]string)
		}
		self.SetStringValueToRedis(PLAYER_SESSIONS, sessions)
		SaveMySessions(self.Id, sessions)
		logs.Info("玩家退出群会话", self.Id)
	} else {
		//redis不存在，直接修改数据库
		DeleteMySessions(self.Id, ids)
		logs.Info("玩家退出群会话", self.Id)
	}

}
func (self *RedisPlayerBaseObj) GetMyInfo() *client_server.PlayerMsg {
	b := len(self.GetPersonalityTags()) > 0
	msg := &client_server.PlayerMsg{
		PlayerId:           easygo.NewInt64(self.Id),
		Gold:               easygo.NewInt64(self.GetGold()),
		HeadIcon:           easygo.NewString(self.GetHeadIcon()),
		NickName:           easygo.NewString(self.GetNickName()),
		Sex:                easygo.NewInt32(self.GetSex()),
		Account:            easygo.NewString(self.GetAccount()),
		Phone:              easygo.NewString(self.GetPhone()),
		Email:              easygo.NewString(self.GetEmail()),
		PeopleID:           easygo.NewString(self.GetPeopleId()),
		BankInfo:           self.GetBankInfos(),
		Signature:          easygo.NewString(self.GetSignature()),
		Provice:            easygo.NewString(self.GetProvice()),
		City:               easygo.NewString(self.GetCity()),
		Area:               easygo.NewString(self.GetArea()),
		IsPayPassword:      easygo.NewBool(self.GetIsPayPassword()),
		PlayerSetting:      self.GetPlayerSetting(),
		RealName:           easygo.NewString(self.GetRealName()),
		BlackList:          self.GetBlackList(),
		ClearLocalLogTime:  easygo.NewInt64(self.GetClearLocalLogTime()),
		IsVisitor:          easygo.NewBool(self.GetIsVisitor()),
		IsLoginPassword:    easygo.NewBool(self.GetIsLoginPassword()),
		FreeTimes:          easygo.NewInt32(self.GetFreeTimes()),
		IsBindWechat:       easygo.NewBool(self.GetIsBindWechat()),
		Emoticons:          self.GetEmoticons(),
		LabelInfo:          GetLabelInfo(self.GetLabelList()),
		AreaCode:           easygo.NewString(self.GetAreaCode()),
		BackgroundImageURL: easygo.NewString(self.GetBackgroundImageURL()),
		Coin:               easygo.NewInt64(self.GetCoin()),
		BCoin:              easygo.NewInt64(self.GetBCoin()),
		YoungPassWord:      easygo.NewString(self.GetYoungPassWord()),
		Types:              easygo.NewInt32(self.GetTypes()),
		IsCanRoam:          easygo.NewBool(self.GetIsCanRoam()),
		Constellation:      easygo.NewInt32(self.GetConstellation()),
		MixId:              easygo.NewInt64(self.GetMixId()),
		ESportCoin:         easygo.NewInt64(self.GetESportCoin()),
		IsSetPersonalTags:  easygo.NewBool(b),
		//Diamond:           easygo.NewInt64(diamond),
	}

	self.GetStringValueToRedis(PLAYER_SETTING, &msg.PlayerSetting)
	if self.GetPayPassword() != "" {
		msg.IsPayPassword = easygo.NewBool(true)
	}
	return msg
}

///========================playersetting=============================

///=======================callinfo===============================
//获取是否与人通话
func (self *RedisPlayerBaseObj) GetCallInfo() *share_message.CallInfo {
	var val *share_message.CallInfo
	self.GetStringValueToRedis(PLAYER_CALLINFO, &val)
	return val
}
func (self *RedisPlayerBaseObj) GetCallInfoPlayerId() int64 {
	callInfo := self.GetCallInfo()
	return callInfo.GetPlayerId()
}

func (self *RedisPlayerBaseObj) SetCallInfo(pid int64, msg *client_hall.SpecialChatInfo) {
	info := &share_message.CallInfo{}
	if msg != nil {
		s, _ := json.Marshal(msg)
		info = &share_message.CallInfo{
			PlayerId: easygo.NewInt64(pid),
			SMsg:     easygo.NewString(string(s)),
		}
	}
	self.SetStringValueToRedis(PLAYER_CALLINFO, info)
}

//检测语音是否是等待接受状态
func (self *RedisPlayerBaseObj) CheckCallInfo() bool {
	data := self.GetDetailCallInfo()
	if data == nil {
		return false
	}
	return data.GetOperate() != 0
}
func (self *RedisPlayerBaseObj) GetDetailCallInfo() *client_hall.SpecialChatInfo {
	callInfo := self.GetCallInfo()
	if callInfo == nil {
		return nil
	}
	var data *client_hall.SpecialChatInfo
	err := json.Unmarshal([]byte(callInfo.GetSMsg()), &data)
	if err != nil {
		return nil
	}
	return data
}

//======================callinfo========================

func (self *RedisPlayerBaseObj) UpdateLogInTimestamp() {
	//如果最后登录的时间不是今天重置今日在线时长为0
	if self.GetLastOnLineTime() < easygo.GetToday0ClockTimestamp()*1000 {
		self.SetTodayOnlineTime(0)
	}

	time := GetMillSecond()
	self.SetLastOnLineTime(time)
}
func (self *RedisPlayerBaseObj) UpdateLoginTimes() {
	self.AddLoginTimes()
}

//检测是否存在reids数据
func (self *RedisPlayerBaseObj) CheckIsExist(key string) bool {
	res, err := easygo.RedisMgr.GetC().Exist(key)
	easygo.PanicError(err)
	return res
}
func (self *RedisPlayerBaseObj) SetAutoLoginInfo(token string) {
	self.SetAutoLoginToken(token)
	self.SetAutoLoginTime(time.Now().Unix())
}

//保存表情到数据
func (self *RedisPlayerBaseObj) SaveEmoticons() {
	saveEmots := self.GetEmoticons()
	//保存表情数据
	if len(saveEmots) > 0 {
		var data []interface{}
		for _, v := range saveEmots {
			if v.GetIsSave() {
				b1 := bson.M{"_id": v.GetId()}
				v.IsSave = easygo.NewBool(false)
				data = append(data, b1, v)
			}

		}
		if len(data) > 0 {
			UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_EMOTICON, data)
		}
	}
}

// 封装外部方法，获取玩家的信息
func GetRedisPlayerBase(playerId PLAYER_ID, player ...*share_message.PlayerBase) *RedisPlayerBaseObj {
	return PlayerBaseMgr.GetRedisPlayerBaseObj(playerId, player...)
}

//完成注销账号
func (self *RedisPlayerBaseObj) CancelAccountFinish() {
	//增加到注销列表
	if self.GetPhone() == "" {
		logs.Info("注销的账号为电话为空:")
	}
	var phone, areaCode string
	phone = self.GetPhone()
	log := &share_message.CancelAccountList{
		Phone:      easygo.NewString(phone),
		FinishTime: easygo.NewInt64(GetMillSecond()),
	}
	SaveCancelPhone(log)
	//手机号设置为空，修改redis数据
	account := GetRedisAccountObj(self.GetPlayerId())
	if account != nil {
		account.SetAccount("")
		account.SetOpenId("")
		account.SetUnionId("")
		//删除账号playerid关联
		account.SaveToMongo()
		//account.DelRedisPlayerAccountPhone()

	}
	self.SetPhone("")
	self.SetStatus(ACCOUNT_CANCELED)
	//修改mongo数据
	self.SaveToMongo()
	self.SaveEmoticons()

	if self.GetAreaCode() == "" {
		areaCode = "+86"
		// 用腾讯运营商
		easygo.Spawn(NewSMSInst(SMS_BUSINESS_TC).SendMessageCodeEx, fmt.Sprintf("%s%s", areaCode, phone), phone, true, false)
	}
	logs.Info("保存注销玩家成功", self.Id)
}

// 获取指定人的置顶动态id列表
func (self *RedisPlayerBaseObj) GetPlayerAllTopIds() []int64 {
	ids := make([]int64, 0)
	bsTopIds := self.GetTopByKeyNameFromRedis(REDIS_SQUARE_BS_TOP_DYNAMIC) // 后台置顶
	appTopIds := self.GetTopByKeyNameFromRedis(REDIS_SQUARE_TOP_DYNAMIC)   // app置顶
	ofTopIds := self.GetTopByKeyNameFromRedis(REDIS_SQUARE_OF_TOP_DYNAMIC) // 官方置顶
	ids = append(ids, bsTopIds...)
	ids = append(ids, appTopIds...)
	ids = append(ids, ofTopIds...)
	return ids
}

// 获取个人某个置顶动态id列表
func (self *RedisPlayerBaseObj) GetTopByKeyNameFromRedis(keyName string) []int64 {
	values, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(keyName))
	easygo.PanicError(err)
	ids := make([]int64, 0)
	for _, v := range values {
		dynamic := &share_message.DynamicData{}
		_ = json.Unmarshal([]byte(v), &dynamic)
		if dynamic.GetPlayerId() == 0 {
			continue
		}
		if dynamic.GetPlayerId() == self.GetPlayerId() {
			ids = append(ids, dynamic.GetLogId())
		}
	}
	return ids
}

//获取玩家之间的共同标签
func (self *RedisPlayerBaseObj) GetCommonTags(pid int64) []string {
	tags := make([]string, 0)
	player := GetRedisPlayerBase(pid)
	if player == nil {
		return tags
	}
	ids := make([]int32, 0)
	for _, l1 := range self.GetPersonalityTags() {
		for _, l2 := range player.GetPersonalityTags() {
			if l1 == l2 {
				ids = append(ids, l1)
				break
			}
		}
		if len(ids) >= 3 { //只显示3条
			break
		}
	}
	if len(ids) > 0 {
		tagsData := GetPlayerPersonalityTags(ids)
		for _, d := range tagsData {
			s := fmt.Sprintf("你们都是%s", d.GetName())
			tags = append(tags, s)
		}
	}
	//已经够了，返回
	if len(tags) >= 3 {
		return tags
	}
	//星座相同
	if self.GetConstellation() != 0 && self.GetConstellation() == player.GetConstellation() {
		s := fmt.Sprintf("你们都是%s", GetConfigConstellationSortName(self.GetConstellation()))
		tags = append(tags, s)
	}
	return tags
}

//获取与指定玩家的匹配度
func (self *RedisPlayerBaseObj) GetMatchingDegree(pid int64) int32 {
	val := int32(0)
	pParam := &share_message.MatchingDegreeParam{
		OnLine:         easygo.NewInt32(46),
		OffLine:        easygo.NewInt32(40),
		SameSex:        easygo.NewInt32(20),
		UnSameSex:      easygo.NewInt32(35),
		LabelMax:       easygo.NewInt32(24),
		LabelPer:       easygo.NewInt32(8),
		PersonalTagMax: easygo.NewInt32(24),
		PersonalTagPer: easygo.NewInt32(4),
		Constellation:  easygo.NewInt32(5),
	}
	//1用户是否在线
	player := GetRedisPlayerBase(pid)
	if player == nil {
		return val
	}
	if player.GetIsOnLine() {
		val += pParam.GetOnLine()
	} else {
		val += pParam.GetOffLine()
	}
	//2 性别、年龄差距
	if self.GetSex() == player.GetSex() {
		val += pParam.GetSameSex()
	} else {
		val += pParam.GetUnSameSex()
	}
	//3 兴趣爱好,最多+15分
	addLabelVal := int32(0)
	for _, l1 := range self.GetLabelList() {
		if addLabelVal >= pParam.GetLabelMax() {
			break
		}
		for _, l2 := range player.GetLabelList() {
			if addLabelVal >= pParam.GetLabelMax() {
				break
			}
			if l1 == l2 {
				addLabelVal += pParam.GetLabelPer()
			}
		}
	}
	val += addLabelVal
	//4 个性化标签 最多26分
	addTagVal := int32(0)
	for _, l1 := range self.GetPersonalityTags() {
		if addTagVal >= pParam.GetPersonalTagMax() {
			break
		}
		for _, l2 := range player.GetPersonalityTags() {
			if addTagVal >= pParam.GetPersonalTagMax() {
				break
			}
			if l1 == l2 {
				addTagVal += pParam.GetPersonalTagPer()
				break
			}
		}
	}
	val += addTagVal
	//5 星座相同+5分
	if self.GetConstellation() == player.GetConstellation() {
		val += pParam.GetConstellation()
	}
	if val > 100 {
		return 100
	}
	return val
}

//获取sayHi数据
func (self *RedisPlayerBaseObj) GetSayHiData(pid int64) string {
	player := GetRedisPlayerBase(pid)
	if player == nil {
		return ""
	}
	data := make([]*client_hall.PlayerCommonData, 0)
	w1 := GetRandMatchGuide("")
	myData := &client_hall.PlayerCommonData{
		PlayerId:      easygo.NewInt64(self.Id),
		HeadUrl:       easygo.NewString(self.GetHeadIcon()),
		Tags:          self.GetPersonalityTagsStr(),
		Constellation: easygo.NewInt32(self.GetConstellation()),
		GuideWord:     easygo.NewString(w1),
		Labels:        self.GetPersonalityTags(),
		Sex:           easygo.NewInt32(self.GetSex()),
	}
	w2 := GetRandMatchGuide(w1)
	playerData := &client_hall.PlayerCommonData{
		PlayerId:      easygo.NewInt64(player.Id),
		HeadUrl:       easygo.NewString(player.GetHeadIcon()),
		Tags:          player.GetPersonalityTagsStr(),
		Constellation: easygo.NewInt32(player.GetConstellation()),
		GuideWord:     easygo.NewString(w2),
		Labels:        player.GetPersonalityTags(),
		Sex:           easygo.NewInt32(self.GetSex()),
	}
	data = append(data, myData, playerData)
	msg := &client_hall.SayHiData{
		PlayerData:     data,
		MatchingDegree: easygo.NewInt32(self.GetMatchingDegree(pid)),
	}
	b, err := json.Marshal(msg)
	easygo.PanicError(err)
	content := base64.StdEncoding.EncodeToString(b)
	return content
}

//获取玩家声音名片信息
func (self *RedisPlayerBaseObj) GetVoiceCardData() *client_hall.VoiceCard {
	card := &client_hall.VoiceCard{
		PlayerId:        easygo.NewInt64(self.Id),
		NickName:        easygo.NewString(self.GetNickName()),
		HeadUrl:         easygo.NewString(self.GetHeadIcon()),
		Sex:             easygo.NewInt32(self.GetSex()),
		ZanNum:          easygo.NewInt32(self.GetVCZanNum()),
		PersonalityTags: self.GetPersonalityTagsStr(),
		BgUrl:           easygo.NewString(self.GetBgImageUrl()),
		VoiceUrl:        easygo.NewString(self.GetMixVoiceUrl()),
	}
	return card
}

//从mongo中查询玩家数据
func GetPlayerById(id PLAYER_ID) *share_message.PlayerBase {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	var obj *share_message.PlayerBase
	err := col.Find(bson.M{"_id": id}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//通过id批量查询玩家
func GetAllPlayerBase(ids []int64, isRedis ...bool) map[int64]*share_message.PlayerBase {
	pNotInit := []int64{} //没有加载数据到内存的玩家列表
	pMap := make(map[int64]*share_message.PlayerBase)
	if len(ids) == 0 {
		return pMap
	}
	b := append(isRedis, true)[0]
	//如果在reids内存，则取内存值
	if b {
		for _, id := range ids {
			b, err := easygo.RedisMgr.GetC().Exist(MakeRedisKey(TABLE_PLAYER_BASE, easygo.AnytoA(id)))
			easygo.PanicError(err)
			if !b {
				pNotInit = append(pNotInit, id)
			} else {
				base := GetRedisPlayerBase(id)
				p := base.GetRedisPlayerBase()
				pMap[id] = p
			}
		}
	} else {
		pNotInit = ids
	}
	//否则从数据库中查取
	if len(pNotInit) > 0 {
		col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
		defer closeFun()
		players := make([]*share_message.PlayerBase, 0)
		err := col.Find(bson.M{"_id": bson.M{"$in": pNotInit}}).All(&players)
		easygo.PanicError(err)
		for _, player := range players {
			pMap[player.GetPlayerId()] = player
		}
	}
	return pMap
}

//从mongo中查询玩家数据
func GetPlayerByPhone(phone string) *share_message.PlayerBase {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	var obj *share_message.PlayerBase
	err := col.Find(bson.M{"Phone": phone}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

func GetPlayerByNickName(nickName string) *share_message.PlayerBase {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	var obj *share_message.PlayerBase
	err := col.Find(bson.M{"NickName": nickName}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

func GetPlayerByAccount(account string) *share_message.PlayerBase {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	var obj *share_message.PlayerBase
	err := col.Find(bson.M{"Account": account}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

// 取附近的运营号.
func GetNearOperation(x, y float64) []*share_message.PlayerBase {
	//col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	//defer closeFun()
	var operations []*share_message.PlayerBase
	//col.Find(bson.M{""})
	return operations
}

// 2.7.6 附近的人
func GetLocationInfoNew(pid int64, b bool, reqMsg *client_hall.LocationInfoNewReq) *client_hall.LocationInfoNewResp {
	//	begin := GetMillSecond()
	page := reqMsg.GetPage()
	//pageSize := reqMsg.GetPageSize()
	pageSize := 10
	if page == 0 {
		page = DEFAULT_PAGE
	}
	result := make([]*share_message.PlayerBase, 0)
	player := GetRedisPlayerBase(pid)
	friends := player.GetFriends()
	friends = append(friends, pid)
	var resultCount int
	// 先判断redis中是否存在
	redisKey := MakeRedisKey(REDIS_NEAR_PLAYER, pid)

	if reqMsg.GetIsNewFlush() || !ExistZAdd(redisKey) || b { // 新的刷新,重新查询数据.
		// 取200人,优先真实用户,运营号30人,剩下的是机器人.
		begin := GetMillSecond()
		pl := FlushLocalInfo(pid, friends, reqMsg)
		end := GetMillSecond()
		logs.Warn("查数据的时间--------->%d 毫秒", end-begin)
		m := make([]interface{}, 0)

		for _, pp := range pl {
			m = append(m, pp)
		}
		// 分页取数据
		resultInterface, count := SliceByPage(int(page), int(pageSize), m)
		resultCount = count

		for _, v := range resultInterface {
			p := v.(*share_message.PlayerBase)
			result = append(result, p)
		}

		// 异步存放redis
		fun := func() {
			SetNearInfoToRedis(redisKey, pl, NEAR_EXPIRE)
		}
		easygo.Spawn(fun)
	} else {
		// 分页取redis
		if start, end := MakeRedisPage(int(page), int(pageSize), NEAR_ALL_NUM); end != 0 { // end为0时 ,分页没数据
			fromRedis := GetNearByPageFromRedis(redisKey, start, end)
			for _, v := range fromRedis {
				var p *share_message.PlayerBase
				_ = json.Unmarshal(v.([]byte), &p)
				if p != nil {
					result = append(result, p)
				}
			}
		}
		resultCount = NEAR_ALL_NUM

	}

	localInfo := make([]*client_hall.LocationInfoNew, 0)
	// 封装数据6
	for _, v := range result {

		ds := make([]*share_message.DynamicData, 0)
		if !v.GetIsRobot() {
			// 取动态,带有图片别切最热的3个动态
			ds1 := GetHasPhotoHotDynamicToDB(v.GetPlayerId(), NEAR_DYNAMIC_NUM)
			logIdList := make([]int64, 0)
			for _, v1 := range ds1 {
				logIdList = append(logIdList, v1.GetLogId())

			}
			mm := GetRedisDynamicZanNumEx(logIdList)
			for _, v1 := range ds1 {
				v1.Zan = easygo.NewInt32(mm[v1.GetLogId()] + v1.GetTrueZan())
				ds = append(ds, v1)
			}
		}

		attentionList := player.GetAttention()
		GetRedisSomeDynamic1Ex(pid, v, ds, attentionList) // 遍历动态,判断是否已关注.
		localInfo = append(localInfo, &client_hall.LocationInfoNew{
			PlayerId:     easygo.NewInt64(v.GetPlayerId()),
			HeadIcon:     easygo.NewString(v.GetHeadIcon()),
			NickName:     easygo.NewString(v.GetNickName()),
			Sex:          easygo.NewInt32(v.GetSex()),
			OnlineStatus: easygo.NewInt32(v.GetOnlineStatus()),
			Signature:    easygo.NewString(v.GetSignature()),
			Distance:     easygo.NewFloat64(v.GetDistance()),
			IsFriend:     nil,
			DynamicList:  ds,
			DataType:     easygo.NewInt32(NEAR_INFO_DATE_TYPE_NOMAL),
			IsRobot:      easygo.NewBool(v.GetIsRobot()),
			Account:      easygo.NewString(v.GetAccount()),
			Types:        easygo.NewInt32(v.GetTypes()),
		})

	}
	//	end := GetMillSecond()
	// 封装引导数
	localInfo = insertLeadData(localInfo)
	return &client_hall.LocationInfoNewResp{
		Count:        easygo.NewInt32(resultCount),
		LocationInfo: localInfo,
	}
}

// 插入引导数据
func insertLeadData(infos []*client_hall.LocationInfoNew) []*client_hall.LocationInfoNew {
	// 查询引导数据
	leadList := GetAllNearLeadFromDB()
	// 没有权重,直接返回
	if len(leadList) == 0 {
		return infos
	}
	weightsMap := make(map[int32]*share_message.NearSet)
	rate := make([]float32, 0)
	for _, v := range leadList {
		weightsMap[v.GetWeights()] = v
		rate = append(rate, easygo.AtoFloat32(easygo.AnytoA(v.GetWeights()))/100)
	}
	// 计算权重
	var index int
	for i := 0; i < 10000; i++ {
		index = WeightedRandomIndex(rate)
	}
	weight := int32(rate[index] * 100)

	lead := weightsMap[weight]

	// 计算插入的位置
	insertIndex := RangeRand(NEAR_INSERT_LEAD_BEGIN, NEAR_INSERT_LEAD_END)
	info := &client_hall.LocationInfoNew{
		NearSet:  lead,
		DataType: easygo.NewInt32(NEAR_INFO_DATE_TYPE_LEAD),
	}
	// 插入引导数据
	if int(insertIndex) > len(infos) {
		insertIndex = int64(len(infos)) - 1
	}
	return easygo.Insert(infos, int(insertIndex), info).([]*client_hall.LocationInfoNew)
}
func checkOnLineStatus(p *share_message.PlayerBase) int32 {
	if p.GetIsOnline() { // 是否在线
		// 是否新人在线
		if GetMillSecond()-p.GetCreateTime() < 7*86400*1000 { // 注册不满7天
			return ONLINE_STATUS_ONLINE_NEW // 新人在线
		}
		return ONLINE_STATUS_NEW // 在线
	}
	if GetMillSecond()-p.GetLastLogOutTime() > 3600*100 {
		return ONLINE_STATUS_OFFLINE // 离线超过1小时
	}
	return ONLINE_STATUS_LOG_OUT // 刚刚
}

// 附近的人运营号抽取
var OperationPhones = []string{
	"10010001477",
	"10010001478",
	"10010001479",
	"10010001480",
	"10010001481",
	"10010001482",
	"10010001483",
	"10010001484",
	"10010001485",
	"10010001486",
	"10010001487",
	"10010001488",
	"10010001489",
	"10010001490",
	"10010001491",
	"10010001492",
	"10010001493",
	"10010001494",
	"10010001495",
	"10010001496",
	"10010001497",
	"10010001498",
	"10010001499",
	"10010001500",
	"10010001501",
	"10010001502",
	"10010001503",
	"10010001504",
	"10010001505",
	"10010001506",
	"10010001507",
	"10010001508",
	"10010001509",
	"10010001510",
	"10010001511",
	"10010001512",
	"10010001513",
	"10010001514",
	"10010001515",
	"10010001516",
	"10010001517",
	"10010001518",
	"10010001519",
	"10010001520",
	"10010001521",
	"10010001522",
	"10010001523",
	"10010001524",
	"10010001525",
	"10010001526",
	"10010001527",
	"10010001528",
	"10010001529",
	"10010001530",
	"10010001531",
	"10010001532",
	"10010001533",
	"10010001534",
	"10010001535",
	"10010001536",
	"10010001537",
	"10010001538",
	"10010001539",
	"10010001540",
	"10010001541",
	"10010001542",
	"10010001543",
	"10010001544",
	"10010001545",
	"10010001546",
	"10010001547",
	"10010001548",
	"10010001549",
	"10010001550",
	"10010001551",
	"10010001552",
	"10010001553",
	"10010001554",
	"10010001555",
	"10010001556",
	"10010001557",
	"10010001558",
	"10010001559",
	"10010001560",
	"10010001561",
	"10010001562",
	"10010001563",
	"10010001564",
	"10010001565",
	"10010001566",
	"10010001567",
	"10010001568",
	"10010001569",
	"10010001570",
	"10010001571",
	"10010001572",
	"10010001573",
	"10010001574",
	"10010001575",
	"10010001576",
}

// 根据经纬度查询200号人
func FlushLocalInfo(pid int64, friends []int64, reqMsg *client_hall.LocationInfoNewReq) []*share_message.PlayerBase {
	pl := make([]*share_message.PlayerBase, 0)
	// 先查运营号,50米内,30个
	phones := util.RandSliceFromSlice(OperationPhones, NEAR_OPERATION_NUM)
	//operationList := GetNearInfoFromDB(friends, ACCOUNT_TYPES_YXYY, reqMsg, NEAR_OPERATION_NUM)
	operationList := GetOperationByPhones(pid, phones, reqMsg)
	// 生成随机距离给客服
	for _, v := range operationList {
		distance := RangeRand(0, NEAR_OPERATIONAL_DISTANCE) // 50米内
		v.Distance = easygo.NewFloat64(distance)
	}
	pl = append(pl, operationList...)
	// 查询真实用户(不包含自己)
	realPlayerList := GetNearInfoFromDB(friends, ACCOUNT_TYPES_PT, reqMsg, NEAR_ALL_NUM-len(operationList))
	pl = append(pl, realPlayerList...)

	var robotNum int // 机器人数量
	if len(pl) < NEAR_ALL_NUM {
		robotNum = NEAR_ALL_NUM - len(pl)
	}

	// 封装机器人
	if robotNum > 0 {
		ps := MakeRobot(robotNum, reqMsg.GetSex())
		if ps != nil {
			pl = append(pl, ps...)
		}
	}
	if reqMsg.GetSort() == NEAR_SORT_DISTANCE {
		sort.Slice(pl, func(i, j int) bool {
			return pl[i].GetDistance() < pl[j].GetDistance() // 升序
		})
	}
	if reqMsg.GetSort() == NEAR_SORT_ONLINE {
		onLine := make([]*share_message.PlayerBase, 0)  // 在线的
		disLine := make([]*share_message.PlayerBase, 0) // 不是在线的
		for _, v := range pl {
			status := checkOnLineStatus(v)
			v.OnlineStatus = easygo.NewInt32(status)
			if status == 1 {
				onLine = append(onLine, v)
			} else {
				disLine = append(disLine, v)
			}
		}

		sort.Slice(onLine, func(i, j int) bool {
			return onLine[i].GetDistance() < onLine[j].GetDistance() // 升序
		})
		sort.Slice(disLine, func(i, j int) bool {
			return disLine[i].GetDistance() < disLine[j].GetDistance() // 升序
		})

		npl := make([]*share_message.PlayerBase, 0)
		npl = append(npl, onLine...)
		npl = append(npl, disLine...)
		// 排序
		for i := 0; i < len(npl); i++ {
			npl[i].NearSort = easygo.NewFloat64(i)
		}
		return npl
	}
	for i := 0; i < len(pl); i++ {
		pl[i].NearSort = easygo.NewFloat64(i)
		status := checkOnLineStatus(pl[i])
		pl[i].OnlineStatus = easygo.NewInt32(status)
	}
	return pl
}

/*func FlushLocalInfo(friends []int64, reqMsg *client_hall.LocationInfoNewReq) []*share_message.PlayerBase {
	pl := make([]*share_message.PlayerBase, 0)
	// 查询真实用户(不包含自己)
	realPlayerList := GetNearInfoFromDB(friends, ACCOUNT_TYPES_PT, reqMsg, NEAR_ALL_NUM)
	pl = append(pl, realPlayerList...)
	// 查询运营号
	var operationNum int                    // 运营号数量
	if len(realPlayerList) < NEAR_ALL_NUM { // 小于200人,30号运管号
		operationNum = NEAR_ALL_NUM - len(realPlayerList)
		if NEAR_ALL_NUM-len(realPlayerList) >= NEAR_OPERATION_NUM {
			operationNum = NEAR_OPERATION_NUM
		}

		operationList := GetNearInfoFromDB(friends, ACCOUNT_TYPES_YXYY, reqMsg, operationNum)
		// 生成随机距离给客服
		for _, v := range operationList {
			distance := RangeRand(100, NEAR_OPERATIONAL_DISTANCE) // 10公里内
			v.Distance = easygo.NewFloat64(distance)
		}
		pl = append(pl, operationList...)
	}

	var robotNum int // 机器人数量
	if len(pl) < NEAR_ALL_NUM {
		robotNum = NEAR_ALL_NUM - len(pl)
	}

	// 封装机器人
	if robotNum > 0 {
		ps := MakeRobot(robotNum, reqMsg.GetSex())
		if ps != nil {
			pl = append(pl, ps...)
		}
	}
	if reqMsg.GetSort() == NEAR_SORT_DISTANCE {
		sort.Slice(pl, func(i, j int) bool {
			return pl[i].GetDistance() < pl[j].GetDistance() // 升序
		})
	}
	if reqMsg.GetSort() == NEAR_SORT_ONLINE {
		onLine := make([]*share_message.PlayerBase, 0)  // 在线的
		disLine := make([]*share_message.PlayerBase, 0) // 不是在线的
		for _, v := range pl {
			status := checkOnLineStatus(v)
			v.OnlineStatus = easygo.NewInt32(status)
			if status == 1 {
				onLine = append(onLine, v)
			} else {
				disLine = append(disLine, v)
			}
		}

		sort.Slice(onLine, func(i, j int) bool {
			return onLine[i].GetDistance() < onLine[j].GetDistance() // 升序
		})
		sort.Slice(disLine, func(i, j int) bool {
			return disLine[i].GetDistance() < disLine[j].GetDistance() // 升序
		})

		npl := make([]*share_message.PlayerBase, 0)
		npl = append(npl, onLine...)
		npl = append(npl, disLine...)
		// 排序
		for i := 0; i < len(npl); i++ {
			npl[i].NearSort = easygo.NewFloat64(i)
		}
		return npl
	}
	for i := 0; i < len(pl); i++ {
		pl[i].NearSort = easygo.NewFloat64(i)
		status := checkOnLineStatus(pl[i])
		pl[i].OnlineStatus = easygo.NewInt32(status)
	}
	return pl
}
*/
// 生成机器人
func MakeRobot(robotNum int, reqSex int32) []*share_message.PlayerBase {
	if robotNum == 0 {
		return nil
	}
	pl := make([]*share_message.PlayerBase, 0)
	rand.Seed(time.Now().Unix())
	manHeadList := GetManyRobotHeadIcon(robotNum, 1)  //男随机头像列表
	girlHeadList := GetManyRobotHeadIcon(robotNum, 2) //女头像随机列表
	manNameList := GetManyRobotName(1, robotNum)
	girlNameList := GetManyRobotName(2, robotNum)
	sex := reqSex
	for i := 0; i < robotNum; i++ {
		if reqSex == PLAYER_SEX_ALL {
			sex = int32(RandInt(1, 3))
		}
		var mark, name string
		var icon int
		if sex == PLAYER_SEX_BOY {
			mark = "mavatar"
			icon = manHeadList[i]
			name = manNameList[i]
		} else {
			mark = "wavatar"
			icon = girlHeadList[i]
			name = girlNameList[i]
		}
		head := fmt.Sprintf("https://im-resource-1253887233.file.myqcloud.com/prod/%s/%d.png", mark, icon)

		dis := RangeRand(NEAR_BEGIN_ROBOT_DISTANCE, NEAR_DISTANCE)
		id := RandInt(Min_Robot_PlayerId, Max_Robot_PlayerId)
		//account := GetRandAccount("lm", int64(id))
		account := "lm***"
		// 获取个性签名
		signature := GetRandSignature()
		pl = append(pl, &share_message.PlayerBase{
			PlayerId:  easygo.NewInt64(id),
			Distance:  easygo.NewFloat64(dis),
			NickName:  easygo.NewString(name),
			Sex:       easygo.NewInt32(sex),
			HeadIcon:  easygo.NewString(head),
			Signature: easygo.NewString(signature),
			Account:   easygo.NewString(account),
			IsRobot:   easygo.NewBool(true),
		})
	}
	return pl
}

// expire 秒
func SetNearInfoToRedis(key string, pl []*share_message.PlayerBase, expire int64) {
	if isExist := ExistZAdd(key); isExist { // 存在的话先删除掉.
		DelZAdd(key)
	}
	m := make(map[float64]string)
	for _, v := range pl {
		bytes, _ := json.Marshal(v)
		m[v.GetNearSort()] = string(bytes)
	}
	err1 := easygo.RedisMgr.GetC().ZAdd(key, m)
	easygo.PanicError(err1)
	err1 = easygo.RedisMgr.GetC().Expire(key, expire)
	easygo.PanicError(err1)
}

func GetNearByPageFromRedis(key string, page, pageSize int) []interface{} {
	value, err1 := easygo.RedisMgr.GetC().ZRange(key, page, pageSize)
	if err1 != nil {
		logs.Error(err1)
	}
	return value
}

func ExistZAdd(key string) bool {
	b, err := easygo.RedisMgr.GetC().Exist(key)
	easygo.PanicError(err)
	return b
}

func DelZAdd(key string) {
	_, err := easygo.RedisMgr.GetC().Delete(key)
	easygo.PanicError(err)

}

func GetNearRecommend(pid, page, pageSize int64, x, y float64) ([]*client_hall.NearRecommend, int) {
	base := GetRedisPlayerBase(pid)
	area := base.GetArea()
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	nrs := make([]*client_hall.NearRecommend, 0)
	pl, count := GetNearRecommendPlayer(pid, int(page), int(pageSize), x, y, area)
	friends := base.GetFriends()
	for _, v := range pl {
		equipmentObj := GetRedisPlayerEquipmentObj(v.GetPlayerId())
		equipment := equipmentObj.GetEquipmentForClient()
		status := checkOnLineStatus(v)
		nrs = append(nrs, &client_hall.NearRecommend{
			PlayerId:     easygo.NewInt64(v.GetPlayerId()),
			NickName:     easygo.NewString(v.GetNickName()),
			HeadIcon:     easygo.NewString(v.GetHeadIcon()),
			Sex:          easygo.NewInt32(v.GetSex()),
			Signature:    easygo.NewString(v.GetSignature()),
			OnlineStatus: easygo.NewInt32(status),
			GJId:         easygo.NewInt64(equipment.GetGJ().GetPropsId()),
			QPId:         easygo.NewInt64(equipment.GetQP().GetPropsId()),
			MPId:         easygo.NewInt64(equipment.GetMP().GetPropsId()),
			QTXId:        easygo.NewInt64(equipment.GetQTX().GetPropsId()),
			MZBSId:       easygo.NewInt64(equipment.GetMZBS().GetPropsId()),
			IsFriend:     easygo.NewBool(util.Int64InSlice(v.GetPlayerId(), friends)),
			Distance:     easygo.NewFloat64(v.GetDistance()),
			Types:        easygo.NewInt32(v.GetTypes()),
		})
	}
	return nrs, count
}

//停服保存处理，保存需要存储的数据
func SaveRedisPlayerBaseToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_PLAYER_BASE, &ids)
	Emots := make([]*share_message.PlayerEmoticon, 0)
	saveData := make([]interface{}, 0)
	for _, id := range ids {
		obj := GetRedisPlayerBase(id)
		if obj != nil {
			data := obj.GetRedisPlayerBase()
			saveData = append(saveData, bson.M{"_id": data.GetPlayerId()}, data)
			emot := obj.GetEmoticons()
			Emots = append(Emots, emot...)
			obj.SetSaveStatus(false)
		}
	}
	if len(Emots) > 0 {
		saveEmots := make([]interface{}, 0)
		for _, it := range Emots {
			saveEmots = append(saveEmots, bson.M{"_id": it.GetId()}, it)
		}
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_EMOTICON, saveEmots)
	}
	if len(saveData) > 0 {
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_BASE, saveData)
	}
}
