package for_game

import "game_server/pb/share_message"

//--------------------------- 表名统一放这里 ---------------------------
const (
	//tool============
	TABLE_TOOL_WISH_BOX           = "tool_wish_box"           //盲盒基础数据表
	TABLE_TOOL_WISH_BOX_ITEM      = "tool_wish_box_item"      //物品基础数据表
	TABLE_TOOL_WISH_POOL_PUMP_LOG = "tool_wish_pool_pump_log" //水池抽水日志
	TABLE_TOOL_PLAYER_WISH_ITEM   = "tool_player_wish_item"   //挑战产生物品
	TABLE_TOOL_WISH_POOL_LOG      = "tool_wish_pool_log"      // 水池流水表
	//===========许愿池=============
	TABLE_DIAMOND_CHANGELOG = "log_diamond_change" //硬币变化log
	TABLE_WISH_PLAYER       = "wish_player"        //许愿池玩家
	TABLE_WISH_MENU         = "wish_menu"          //许愿池主页菜单
	TABLE_WISH_BRAND        = "wish_brand"         //许愿池物品品牌
	TABLE_WISH_ITEM_TYPE    = "wish_item_type"     //许愿池物品类型
	TABLE_WISH_STYLE        = "wish_style"         //许愿池物品款式
	TABLE_WISH_ITEM         = "wish_item"          //许愿池物品
	TABLE_WISH_BOX_ITEM     = "wish_box_item"      //许愿池盲盒商品
	TABLE_WISH_BOX          = "wish_box"           //许愿池盲盒
	TABLE_WISH_LOG          = "wish_log"           //许愿池挑战记录
	TABLE_WISH_WHITE        = "wish_white"         //白名单
	TABLE_PLAYER_WISH_ITEM  = "player_wish_item"   // 玩家物品列表(挑战产生物品)
	TABLE_PLAYER_WISH_DATA  = "player_wish_data"   // 玩家愿望数据
	//TABLE_WISH_TOP_LOG              = "wish_top_log"              //  玩家排行榜列表.
	TABLE_PLAYER_WISH_COLLECTION    = "player_wish_collection"    //  玩家收藏记录
	TABLE_WISH_OCCUPIED             = "wish_occupied"             // 盲盒挑战占领时长表
	TABLE_WISH_GUARDIAN_DIAMOND_LOG = "wish_guardian_diamond_log" //守护者获得的钻石流水
	TABLE_WISH_GUARDIAN_TOP_LOG     = "wish_guardian_top_log"     // 玩家排行榜列表
	TABLE_WISH_SUM_OCCUPIED         = "wish_sum_occupied"         // 汇总占领时长排行榜总表

	TABLE_PLAYER_EXCHANGE_LOG = "player_exchange_log" //  玩家兑换记录
	TABLE_WISH_RECYCLE_ORDER  = "wish_recycle_order"  //  玩家回收记录

	TABLE_WISH_BOX_ITEM_WIN_CFG        = "wish_box_item_win_cfg"        //盲盒商品中奖配置
	TABLE_WISH_POOL                    = "wish_pool"                    // 盲盒水池
	TABLE_WISH_POOL_CFG                = "wish_pool_cfg"                // 盲盒水池
	TABLE_WISH_PRICE_SECTION           = "wish_price_section"           // 许愿池 搜索价格区间
	TABLE_WISH_MAIL_SECTION            = "wish_mail_section"            // 许愿池 邮寄参数
	TABLE_WISH_RECYCLE_SECTION         = "wish_recycle_section"         // 许愿池 物品回收参数
	TABLE_WISH_PAY_WARN_CFG            = "wish_pay_warn_cfg"            // 许愿池 物品回收参数
	TABLE_WISH_CURRENCY_CONVERSION_CFG = "wish_currency_conversion_cfg" // 许愿池 货币换算参数设置
	TABLE_WISH_RECYCLE_REASON          = "wish_recycle_reason"          // 许愿池 物品回收理由
	TABLE_WISH_GUARDIAN_CFG            = "wish_wish_guardian_Cfg"       // 许愿池 守护者收益设置
	TABLE_WISH_POOL_CFG_ONCE           = "wish_pool_cfg_once"           // 许愿池水池配置初始化一次
	TABLE_WISH_RECYCLE_NOTE_CFG        = "wish_recycle_note"            // 许愿池回收说明
	//===========许愿池活动============
	TABLE_WISH_ACT_POOL                  = "wish_act_pool"                  // 许愿池 活动奖池管理
	TABLE_WISH_ACT_POOL_RULE             = "wish_act_pool_rule"             // 许愿池 活动奖池累计规则
	TABLE_WISH_COIN_RECHARGE_ACT_CFG     = "wish_coin_recharge_act_cfg"     //充值活动配置表
	TABLE_WISH_PLAYER_ACTIVITY           = "wish_player_activity"           // 玩家活动数据
	TABLE_WISH_PLAYER_ACCESS_LOG         = "wish_player_access_log"         //许愿池埋点报表
	TABLE_WISH_PLAYER_RECHARGE_ACT_LOG   = "wish_player_recharge_act_log"   //充值活动用户记录
	TABLE_WISH_PLAYER_RECHARGE_ACT_FIRST = "wish_player_recharge_act_first" //充值活动用户首充记录
	TABLE_REPORT_WISH_LOG                = "report_wish_log"                //许愿池埋点报表
	TABLE_WISH_BURYING_POINT_LOG         = "wish_burying_point_log"         //许愿池埋点日志
	TABLE_WISH_ACTIVITY_PRIZE_LOG        = "wish_activity_prize_log"        //许愿池活动奖项日志记录

	TABLE_WISH_WEEK_TOP         = "wish_week_top"         // 周排名
	TABLE_WISH_MONTH_TOP        = "wish_month_top"        // 月排名
	TABLE_WISH_DAY_ACTIVITY_LOG = "wish_day_activity_log" // 玩家抽奖记录,去重天数的

	//===========许愿池后台============
	TABLE_WISH_OTHER_CFG               = "wish_other_cfg"               //  许愿池其他配置.
	TABLE_REPORT_WISH_POOL             = "report_wish_pool"             //  许愿池报表.
	TABLE_REPORT_WISH_POOL_WEEK        = "report_wish_pool_week"        //  许愿池报表(周表).
	TABLE_REPORT_WISH_POOL_MONTH       = "report_wish_pool_month"       //  许愿池报表(月表).
	TABLE_REPORT_WISH_BOX              = "report_wish_box"              //  许愿池盲盒报表.
	TABLE_REPORT_WISH_BOX_WEEK         = "report_wish_box_week"         //  许愿池盲盒报表.(周表)
	TABLE_REPORT_WISH_BOX_MONTH        = "report_wish_box_month"        //  许愿池盲盒报表.(月表)
	TABLE_REPORT_WISH_BOX_DETAIL       = "report_wish_box_detail"       //  许愿池盲盒详情报表.
	TABLE_REPORT_WISH_BOX_DETAIL_WEEK  = "report_wish_box_detail_week"  //  许愿池盲盒详情报表（周表）.
	TABLE_REPORT_WISH_BOX_DETAIL_MONTH = "report_wish_box_detail_month" //  许愿池盲盒详情报表（月表）.
	TABLE_REPORT_WISH_BOX_TEMP         = "report_wish_box_temp"         //  许愿池盲盒报表临时表.
	TABLE_REPORT_WISH_BOX_DETAIL_TEMP  = "report_wish_box_detail_temp"  //  许愿池盲盒详情报表临时表.
	TABLE_REPORT_WISH_ITEM             = "report_wish_item"             //  许愿池商品报表.
	TABLE_REPORT_WISH_ACTIVITY         = "report_wish_activity"         //  许愿池活动报表.
	TABLE_WISH_POOL_PUMP_LOG           = "wish_pool_pump_log"           //  水池抽水日志.
	TABLE_WISH_POOL_LOG                = "wish_pool_log"                //  水池流水日志.
	TABLE_WISH_COOL_DOWN_CONFIG        = "wish_cool_down_config"        //  水池冷却配置表
	TABLE_WISH_DIAMOND_RECHARGE        = "wish_diamond_recharge"        //钻石配置表
	TABLE_WISH_PAYPLAYER_LOCATION_LOG  = "wish_log_payplayer_location"  //付费用户地理位置分布日志
)

const (
	WISH_DOWN_STATUS        = 0 // 盲盒下架
	WISH_PUT_ON_STATUS      = 1 // 盲盒上架
	WISH_ADD_PRODUCT_STATUS = 2 // 积极补货中
)
const (
	WISH_OCCUPIED_STATUS_UP   = 1 // 占领中
	WISH_OCCUPIED_STATUS_DOWN = 2 // 占领结束
)
const (
	WISH_PLAYER_TYPE_NORMAL = 0 // 普通用户
	WISH_PLAYER_TYPE_ROBOT  = 1 // 机器人
)

// 水池状态
const (
	POOL_STATUS_BIGLOSS   = 1 // 大亏
	POOL_STATUS_SMALLLOSS = 2 // 小亏
	POOL_STATUS_COMMON    = 3 // 普通
	POOL_STATUS_BIGWIN    = 4 // 大盈
	POOL_STATUS_SMALLWIN  = 5 // 小盈
)

const (
	ITEM_REWARDLV_1 = 1 // 小奖
	ITEM_REWARDLV_2 = 2 // 大奖
)

const (
	WISH_POOL_LOG_TYPE_1 = 1 // 1-抽水扣除
	WISH_POOL_LOG_TYPE_2 = 2 //  2-挑战增加
	WISH_POOL_LOG_TYPE_3 = 3 //  3-得到物品后扣除
)

//玩家许愿挑战状态
const (
	WISH_CHALLENGE_WAIT    = 0 //待挑战
	WISH_CHALLENGE_SUCCESS = 1 //挑战成功
	WISH_CHALLENGE_FAIL    = 2 //取消
)

const (
	WISH_DARE    = 1 // 挑战赛
	WISH_NO_DARE = 2 // 非挑战赛
)

//赠送币种
const (
	RECHARGE_GIVE_DIAMOND = 1 // 钻石
	RECHARGE_GIVE_ESCOIN  = 2 // 电竞币种
)

const (
	WISH_ERR_1 = "您已许愿过此愿望"
)

//钻石日志结构
type DiamondLog struct {
	share_message.DiamondChangeLog `bson:",inline,omitempty"`
	Extend                         *share_message.GoldExtendLog `bson:"Extend,omitempty"` //扩展数据
}
type CommonDiamond struct {
	share_message.DiamondChangeLog `bson:",inline,omitempty"` // inline 类型不能用指针
	Extend                         interface{}                `bson:"Extend,omitempty"`
}
