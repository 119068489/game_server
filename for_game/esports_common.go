package for_game

import (
	"encoding/json"
	"fmt"
	"game_server/easygo/util"
	"game_server/pb/share_message"
	"io/ioutil"
	"time"

	"github.com/astaxie/beego/logs"
)

//异常代号表
const (
	C_SYS_ERROR             int32 = 0  ////0系统异常
	C_OPT_SUCCESS           int32 = 1  ////1成功
	C_DISABLE               int32 = 3  ////3被禁用
	C_PASSWORD_ERROR        int32 = 2  //// 2密码错误，
	C_INFO_ERROR            int32 = 5  ////5输入信息有误
	C_SYSTEM_SLEEP          int32 = 4  //4 服务器维护
	C_INFO_EXISTS           int32 = 6  ////6信息已存在
	C_INFO_NOT_EXISTS       int32 = 7  ////7信息不存在
	C_ACCOUNT_LOGIN         int32 = 8  ////账号在其他登录
	C_NOT_LOGIN             int32 = 9  ////账号没登录登录
	C_NOT_POWER             int32 = 10 //没权限
	C_OVER_QUOTA_ONE        int32 = 11 //单次投注超额
	C_DEDUCT_MONEY_FAIL     int32 = 12 //扣硬币、电竞币余额不足
	C_SETTLEMENT_MONEY_FAIL int32 = 13 //兑换获得、结算返回电竞币失败
	C_OVER_QUOTA_ONE_DAY    int32 = 14 //单日投注超额
	C_VIOLATE_CONTENT       int32 = 15 //文字、图片违规
	C_ACTIVE_STOP           int32 = 16 //活动停止
)

//常用状态
const (
	//电竞新闻资讯视频状态 0未发布 1已发布 2已删除 3已禁用 4已过期
	ESPORTS_NEWS_STATUS_0 int32 = 0 //未发布
	ESPORTS_NEWS_STATUS_1 int32 = 1
	ESPORTS_NEWS_STATUS_2 int32 = 2
	ESPORTS_NEWS_STATUS_3 int32 = 3
	ESPORTS_NEWS_STATUS_4 int32 = 4

	//比赛发布状态
	GAME_RELEASE_FLAG_1 int32 = 1 //未发布
	GAME_RELEASE_FLAG_2 int32 = 2 //已发布

	//视频类型
	ESPORTS_VIDEO_TYPE_1 int32 = 1 //视频
	ESPORTS_VIDEO_TYPE_2 int32 = 2 //直播(放映厅)

	//系统消息 0-未推送,1-已推送，2-已过期
	ESPORTS_STATUS_0 int32 = 0
	ESPORTS_STATUS_1 int32 = 1
	ESPORTS_STATUS_2 int32 = 2

	//电竞标签类别
	// 1 行为标签  2 系统标签 3 游戏标签
	ESPORTS_LABEL_TYPE_1 int32 = 1
	ESPORTS_LABEL_TYPE_2 int32 = 2
	ESPORTS_LABEL_TYPE_3 int32 = 3

	//评论状态
	ESPORTS_COMM_STATUS_1 int32 = 1 //正常
	ESPORTS_COMM_STATUS_2 int32 = 2 //前端删除
	ESPORTS_COMM_STATUS_3 int32 = 3 //后台删除

	//轮播图启用状态 1正常 2禁止
	ESPORTS_BANNER_STATUS_1 int32 = 1
	ESPORTS_BANNER_STATUS_2 int32 = 2
)

//=======================api和app对应的配置项目=========开始======================================
//自定义电竞内嵌版项目LABEL值(值固定不变,可以减少新增)
const (
	ESPORTS_LABEL_WZRY  int32 = 10001 //王者荣耀
	ESPORTS_LABEL_DOTA2 int32 = 10002 //dota2
	ESPORTS_LABEL_LOL   int32 = 10003 //英雄联盟lol
	ESPORTS_LABEL_CSGO  int32 = 10004 //CSGO
	ESPORTS_LABEL_OTHER int32 = 20001 //其他
)

//自定义电竞内嵌版项目LABEL对应的名字(固定配置)
const (
	ESPORTS_LABEL_WZRY_NAME  string = "王者荣耀"  //王者荣耀
	ESPORTS_LABEL_DOTA2_NAME string = "DOTA2" //dota2
	ESPORTS_LABEL_LOL_NAME   string = "英雄联盟"  //英雄联盟lol
	ESPORTS_LABEL_CSGO_NAME  string = "CSGO"  //CSGO
	ESPORTS_LABEL_ORDER_NAME string = "其他"    //CSGO
)

//野子科技event_id(固定配置)
const (
	YEZI_ESPORTS_EVENT_WZRY string = "254" //王者荣耀
	//YEZI_ESPORTS_EVENT_DOTA2 string = "24"  //dota2
	YEZI_ESPORTS_EVENT_LOL string = "11" //英雄联盟lol
	//YEZI_ESPORTS_EVENT_CSGO  string = "205" //CSGO
)

//自定义Api接口来源
const (
	ESPORTS_API_ORIGIN_ID_YEZI int32 = 90001 //野子科技
)

//通过来源id取得来源名称
var ApiOriginIdToNameMap = map[int32]string{
	ESPORTS_API_ORIGIN_ID_YEZI: "野子科技", //野子科技
}

//app自己的标签转化为API过来的来源id的map
//key app项目游戏LABELID值
//value API过来的来源id(对应api的来源id)
//前端app传入label查询是哪个api来源id(可变配置)
// TODO
var LabelToESportApiOriginIdMap = map[int32]int32{
	ESPORTS_LABEL_WZRY: ESPORTS_API_ORIGIN_ID_YEZI, //王者荣耀(野子科技ORIGIN_ID)
	//ESPORTS_LABEL_DOTA2: ESPORTS_API_ORIGIN_ID_YEZI, //dota2(野子科技ORIGIN_ID)
	ESPORTS_LABEL_LOL: ESPORTS_API_ORIGIN_ID_YEZI, //英雄联盟lol(野子科技ORIGIN_ID)
	//ESPORTS_LABEL_CSGO:  ESPORTS_API_ORIGIN_ID_YEZI, //CSGO(野子科技ORIGIN_ID)
}

//API过来的EventId转化为app自己的标签的map
//key api项目游戏EventId值
//value app自己的标签(对应api的值,测试环境和生产环境发布前都要改)
//key可配置
//通过该map确定某一api是否在我们的配置中(可变配置)
// TODO
var EventIdToESportLabelMap = map[string]int32{
	YEZI_ESPORTS_EVENT_WZRY: ESPORTS_LABEL_WZRY, //王者荣耀(野子科技eventID)
	//YEZI_ESPORTS_EVENT_DOTA2: ESPORTS_LABEL_DOTA2, //dota2(野子科技eventID)
	YEZI_ESPORTS_EVENT_LOL: ESPORTS_LABEL_LOL, //英雄联盟lol(野子科技eventID)
	//YEZI_ESPORTS_EVENT_CSGO:  ESPORTS_LABEL_CSGO,  //CSGO(野子科技eventID)
}

//app自己的标签ID转化为游戏项目名称的map(固定配置)
//key app项目游戏LABELID值
//value 对应的游戏项目的名字
// TODO
var LabelToESportNameMap = map[int32]string{
	ESPORTS_LABEL_WZRY:  ESPORTS_LABEL_WZRY_NAME,  //王者荣耀
	ESPORTS_LABEL_DOTA2: ESPORTS_LABEL_DOTA2_NAME, //dota2
	ESPORTS_LABEL_LOL:   ESPORTS_LABEL_LOL_NAME,   //英雄联盟lol
	ESPORTS_LABEL_CSGO:  ESPORTS_LABEL_CSGO_NAME,  //CSGO
	ESPORTS_LABEL_OTHER: ESPORTS_LABEL_ORDER_NAME,
}

//=======================api和app对应的配置项目=========结束======================================

const (
	//电竞菜单定义
	ESPORTMENU_SYS_MSG    = 100 //系统消息 100= 系统消息
	ESPORTMENU_REALTIME   = 301 //资讯
	ESPORTMENU_RECREATION = 302 //娱乐
	ESPORTMENU_GAME       = 303 //比赛
	ESPORTMENU_LIVE       = 304 //放映厅
	ESPORTMENU_SHOP       = 305 //电竞币商城

	//im在电竞中的模块定义 1 消息、2 通讯录、3 广场、4电竞、 5我的
	ESPORT_MODLE_1 = 1 // 消息
	ESPORT_MODLE_2 = 2 // 通讯录
	ESPORT_MODLE_3 = 3 // 广场
	ESPORT_MODLE_4 = 4 // 电竞模块
	ESPORT_MODLE_5 = 5 // 我的
)
const (
	//电竞简要记录定义
	ESPORT_FLOW_VIDEO_HISTORY       = 102 //娱乐视频播放历史记录
	ESPORT_FLOW_LIVE_HISTORY        = 103 //放映厅直播放历史记录
	ESPORT_FLOW_LIVE_FOLLOW_HISTORY = 104 //放映厅直播关注
	/*
		//1000、人单次进入电竞时长
		//2000、菜单
		//3000、自定标签
		//4000、扩展tabId
		//5000、内容页
		//6000、内容页子页
	*/
	//电竞埋点内容类型定义
	ESPORT_BPS_PAGE_TYPE_1 = 1000 //im在电竞中的模块定义
	ESPORT_BPS_PAGE_TYPE_2 = 2000 //电竞顶上菜单
	ESPORT_BPS_PAGE_TYPE_3 = 3000 //电竞标签，包括自定义标签、默认标签、游戏标签等
	ESPORT_BPS_PAGE_TYPE_4 = 4000 //扩展tabId
	ESPORT_BPS_PAGE_TYPE_5 = 5000 //内容页
	ESPORT_BPS_PAGE_TYPE_6 = 6000 //内容页子页
)

const (

	//比赛状态 0 未开始，1 进行中，2 已结束
	//注意:数据库得到的2状态可以直接用, 0、1状态要结合开始时间通过方法getGameStatus确定
	GAME_STATUS_0 string = "0"
	GAME_STATUS_1 string = "1"
	GAME_STATUS_2 string = "2"

	//比赛状态描述 0 默认，1 正常结束，2 延时结束, 3 取消结束
	GAME_STATUS_TYPE_0 string = "0"
	GAME_STATUS_TYPE_1 string = "1"
	GAME_STATUS_TYPE_2 string = "2"
	GAME_STATUS_TYPE_3 string = "3"

	//比赛维度
	GAME_DIMENSION_1 string = "1" //两对战
	GAME_DIMENSION_2 string = "2" //多组混战
	GAME_DIMENSION_3 string = "3" //多个人混战

	//是否有推流直播 0 否，1 是
	GAME_HAVE_LIVE_0 int32 = 0
	GAME_HAVE_LIVE_1 int32 = 1

	//是否有赔率 0 否，1 是
	GAME_IS_BET_0 int32 = 0
	GAME_IS_BET_1 int32 = 1

	//是否有滚盘 0 否，1 是
	GAME_HAVE_ROLL_0 int32 = 0
	GAME_HAVE_ROLL_1 int32 = 1

	//1:早盘 2:滚盘
	GAME_IS_MORN_ROLL_1 int32 = 1
	GAME_IS_MORN_ROLL_2 int32 = 2

	//后台开启投注内容的盘口开关1:封盘 2：开盘
	//判断0是为了防止数据库第一次未入库的时候用、0也是封盘
	GAME_APP_GUESS_FLAG_0 int32 = 0
	GAME_APP_GUESS_FLAG_1 int32 = 1
	GAME_APP_GUESS_FLAG_2 int32 = 2

	//后台开启投注内容的显示开关1:不显示 2：显示
	//判断0是为了防止数据库第一次未入库的时候用、0也是封盘
	GAME_APP_GUESS_VIEW_FLAG_0 int32 = 0
	GAME_APP_GUESS_VIEW_FLAG_1 int32 = 1
	GAME_APP_GUESS_VIEW_FLAG_2 int32 = 2

	//投注状态 1进行中(待结算)、2完成(成功、失败)、3无效(返还)、4违规(扣除)
	GAME_GUESS_BET_STATUS_1 string = "1"
	GAME_GUESS_BET_STATUS_2 string = "2"
	GAME_GUESS_BET_STATUS_3 string = "3"
	GAME_GUESS_BET_STATUS_4 string = "4"

	//投注结果 1待结算、2成功、3失败、4返还、5扣除
	GAME_GUESS_BET_RESULT_1 string = "1"
	GAME_GUESS_BET_RESULT_2 string = "2"
	GAME_GUESS_BET_RESULT_3 string = "3"
	GAME_GUESS_BET_RESULT_4 string = "4"
	GAME_GUESS_BET_RESULT_5 string = "5"

	//注单无效理由
	GAME_GUESS_BET_DISABLE_REASON_1 string = "比赛异常" //比赛异常:比赛无结果、开赛时间提前、比赛人员变更后api通知取消结束等比赛异常情况
	GAME_GUESS_BET_DISABLE_REASON_2 string = "时段无效" //时段无效：封盘截止提前量时段中下注
	GAME_GUESS_BET_DISABLE_REASON_3 string = "人工操作" //人工操作：后台手动判定无效

	//内部服务用
	GAME_GUESS_BET_DISABLE_CODE_1 int32 = 1 //比赛异常无效
	GAME_GUESS_BET_DISABLE_CODE_2 int32 = 2 //时段无效

	//api初始化与回调标识
	INIT_CALLBACK_FLAG_1 int32 = 1 //初始化
	INIT_CALLBACK_FLAG_2 int32 = 2 //回调

	//api过来的建议投注状态(0 关闭，1 开放 ，3 暂停)
	GAME_GUESS_ODDS_STATUS_0 string = "0"
	GAME_GUESS_ODDS_STATUS_1 string = "1"
	GAME_GUESS_ODDS_STATUS_3 string = "3"

	//竞猜项是否有结果(0 否 1 是)
	GAME_GUESS_ITEM_STATUS_0 string = "0"
	GAME_GUESS_ITEM_STATUS_1 string = "1"

	//该竞猜项是否达成(0未达成，1达成 , -1 无结果)
	GAME_GUESS_ITEM_WIN_0     string = "0"
	GAME_GUESS_ITEM_WIN_1     string = "1"
	GAME_GUESS_ITEM_WIN_NORST string = "-1"

	//页面用:投注状态(1:可投注,2:封盘)(通过api过来的建议投注状态组合判断)
	GAME_GUESS_ITEM_ODDS_STATUS_1 string = "1"
	GAME_GUESS_ITEM_ODDS_STATUS_2 string = "2"

	//竞猜类型
	GAME_GUESS_ITEM_BET_TYPE_1 string = "1" //队伍
	GAME_GUESS_ITEM_BET_TYPE_2 string = "2" //选手
	GAME_GUESS_ITEM_BET_TYPE_3 string = "3" //自定义
	GAME_GUESS_ITEM_BET_TYPE_5 string = "5" //让分
	GAME_GUESS_ITEM_BET_TYPE_6 string = "6" //战队数组
	GAME_GUESS_ITEM_BET_TYPE_7 string = "7" //大小数组

	//group_flag
	GAME_GUESS_ITEM_GROUP_FLAG_1 string = "1" //1:+、>
	GAME_GUESS_ITEM_GROUP_FLAG_2 string = "2" //2:-、<

	LIVE_ROOM_OPT_1 int32 = 1 //房间说话
	LIVE_ROOM_OPT_2 int32 = 2 //进入房间

	ESPORT_PLAYER_STATUS_1 int32 = 1 //用户状态正常
	ESPORT_PLAYER_STATUS_2 int32 = 2 //用户状态禁用

	//是否有放映厅 0否,1是(详情也用于判断放映厅按钮)
	GAME_HAVE_VIDEO_HALL_STATUS_0 int32 = 0
	GAME_HAVE_VIDEO_HALL_STATUS_1 int32 = 1

	//是否开奖 0：未开奖 1:已开奖
	GAME_IS_LOTTERY_0 int32 = 0
	GAME_IS_LOTTERY_1 int32 = 1

	//电竞比赛、赔率等的Redis失效时间(3小时)
	ESPORT_GAME_REDIS_EXPIRE_TIME int64 = 10800 //过期时间

	//金额结算用
	LOTTERY_FLAG_1 int32 = 1 //开奖无效
	LOTTERY_FLAG_2 int32 = 2 //开奖成功
)

const (
	//兑换类型1-日常赠送,2-首充赠送,3-活动赠送
	ESPORT_EXCHANGE_TYPE_1 int32 = 1
	ESPORT_EXCHANGE_TYPE_2 int32 = 2
	ESPORT_EXCHANGE_TYPE_3 int32 = 3
)

//电竞币兑换中间用类型
const (
	ESPORTS_EXCHANGE_NORMAL int32 = 1 //正常额度兑换
	ESPORTS_EXCHANGE_WHITE  int32 = 2 //白名单兑换
	ESPORTS_EXCHANGE_FIRSRT int32 = 3 //首充
	ESPORTS_EXCHANGE_DAY    int32 = 4 //日常
	ESPORTS_EXCHANGE_ACTIVE int32 = 5 //活动
)

//增量标志位
const (
	ESPORTS_INCREMENT_FLAG_1 int32 = 1 //全量
	ESPORTS_INCREMENT_FLAG_2 int32 = 2 //处理当前时间前的数据
)

//竞猜项目中对应的num对应的局数(尽可能多定义)
//key 局数
//value 局数名称
var GuessNumMap = map[string]string{
	"0":  "全场",
	"1":  "第一局",
	"2":  "第二局",
	"3":  "第三局",
	"4":  "第四局",
	"5":  "第五局",
	"6":  "第六局",
	"7":  "第七局",
	"8":  "第八局",
	"9":  "第九局",
	"10": "第十局",
	"11": "第十一局",
	"12": "第十二局",
	"13": "第十三局",
	"14": "第十四局",
	"15": "第十五局",
}

//比赛进行中状态不推送，用户需要根据比赛时间自行处理
//如果数据库中不是2 已经结束的状态 传入比赛的开始时间取得比赛状态
//返回0 未开始   1 进行中"2021-01-31 17:59:00" 2已结束
func GetGameStatus(beginTime string, dbGameStatus string) string {
	if dbGameStatus == GAME_STATUS_2 {
		return dbGameStatus
	}

	if "" == beginTime {
		return GAME_STATUS_0
	}
	// string转化为时间，layout必须为 "2006-01-02 15:04:05"
	times, _ := time.ParseInLocation("2006-01-02 15:04:05", beginTime, time.Local)
	timeUnix := times.Unix()

	nowTime := time.Now().Unix()
	if timeUnix > nowTime {

		return GAME_STATUS_0
	} else {
		return GAME_STATUS_1
	}
}

//  获取当天开始时间时间戳 如:2021-02-01 00:00:00的时间戳
func GetGameTodayStartTime() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return t.Unix()
}

//获取当天开始时间字符串 如:2021-02-01 00:00:00
func GetGameTodayStartTimeStr() string {

	return util.FormatUnixTime(GetGameTodayStartTime())
}

//获取当天最后一秒的时间字符串 如:2021-02-01 23:59:59
func GetGameTodayEndTimeStr() string {

	return util.FormatUnixTime(GetGameTodayEndTime())
}

//  获取当天最后一秒的时间戳 如:2021-02-01 23:59:59时间戳
func GetGameTodayEndTime() int64 {
	timeEnd := GetGameTodayStartTime() + (24*60*60 - 1)
	return timeEnd
}

func GetGameTimeStrToInt64(timeStr string) int64 {
	if "" == timeStr {
		return 0
	}
	// string转化为时间，layout必须为 "2006-01-02 15:04:05"
	times, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	timeUnix := times.Unix()

	return timeUnix
}

//拼接比赛展示名称
func GetMatchVSName(matchName, matchStage, bo, teamAName, teamBName string) string {
	return fmt.Sprintf("%s%s -BO%s %s VS %s", matchName, matchStage, bo, teamAName, teamBName)
}

//电竞redis 名字键前缀
func ESportExN(table string) string {
	return fmt.Sprintf("redis_esport:%s", table)
}

//取得比赛头部数据redis组合key的内容
func GetGameHeadGroupKey(apiOrigin int32, appLabelId int32, gameId string) string {
	return fmt.Sprintf("%v_%v_%v", apiOrigin, appLabelId, gameId)
}

//查询英雄
func FindHeroInfos(url string) map[string]*share_message.RealTimeHeroObject {

	heroMap := make(map[string]*share_message.RealTimeHeroObject)

	data, err := ioutil.ReadFile(url)
	if err != nil {
		logs.Error(err)
	}

	err = json.Unmarshal(data, &heroMap)
	if err != nil {
		logs.Error(err)
	}

	return heroMap
}

//查询装备
func FindEquipInfos(url string) map[string]*share_message.RealTimeItemObject {
	equipMap := make(map[string]*share_message.RealTimeItemObject)

	data, err := ioutil.ReadFile(url)
	if err != nil {
		logs.Error(err)
	}

	err = json.Unmarshal(data, &equipMap)
	if err != nil {
		logs.Error(err)
	}

	return equipMap
}
