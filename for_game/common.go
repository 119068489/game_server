package for_game

import (
	"bytes"
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	jpushclient "game_server/easygo/jpush"
	"game_server/easygo/util"
	"game_server/pb/brower_backstage"
	"game_server/pb/client_server"
	"game_server/pb/share_message"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"

	mahonia "github.com/axgle/mahonia"
	"github.com/gocarina/gocsv"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/garyburd/redigo/redis"

	"github.com/astaxie/beego/logs"

	crand "crypto/rand"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
)

type PRIMARY_KEY = interface{}
type PLAYER_ID = int64
type SITE = string
type ENDPOINT_ID = easygo.ENDPOINT_ID
type SERVER_ID = int32
type REDPACKET_ID = int64
type TRANSFERMONEY_ID = int64
type ORDER_ID = string

type GET_C = func() (c *mgo.Collection, fun func())
type IPrimaryKey interface {
	Keys() []interface{}
}

const SERVER_NAME = "nmcl" //服务器项目名
const ETCD_SERVER_PATH = "/" + SERVER_NAME + "/"
const MSG_FLAG = 10085 //服务器间通讯标识
const VOICE_CARD_FRESH_MAX_NUM = 100

const STRANGER_MAX_TALK_NUM = 3

var AES_KEY = []byte("bzkj2019qwertyui")
var MD5_KEY = string("bzkj@999.")

var AES_KEY_SHOP_CARD = []byte("nmcl2020shop!@#$card")

var NOT_NEED_LOGIN_PB = []string{"RpcTFToServer", "RpcGetAllGames", "RpcPhoneRegister", "RpcPhoneFindPwd", "RpcGetSiteList", "RpcAccountCancel"}
var IS_LABA_TEST_GAME bool

//后台操作类型 //新增类型以后要在InitManageLogTypes()中增加
const (
	LOGIN_BACKSTAGE   = "登录后台"
	SIGNOUT_BACKSTAGE = "退出后台"
	USER_MANAGE       = "管理员管理"
	BSPLAYER_MANAGE   = "用户管理"
	ROLE_MANAGE       = "权限管理"
	TEAM_MANAGE       = "群管理"
	SITE_MANAGE       = "站点管理"
	SYS_MAIL          = "系统邮件"
	SYS_MSG           = "系统消息"
	PAY_MANAGE        = "支付管理"
	FEATURES_MANAGE   = "功能管理"
	SHOP_MANAGE       = "商城管理"
	OPERATION_MANAGE  = "运营管理"
	WAITER_MANAGE     = "客服管理"
	SQUARE_MANAGE     = "社交广场管理"
	TOPIC_MANAGE      = "话题管理"
	ADV_MANAGE        = "广告管理"
	MALL_MANAGE       = "虚拟商城管理"
	COINS_MANAGE      = "硬币管理"
	PROPS_MANAGE      = "道具管理"

	WISH_MANAGE    = "许愿池管理"
	LOVE_MANAGE    = "恋爱交友匹配管理"
	ESPORTS_MANAGE = "电竞管理"
)

const (
	MessageCheckTime   = 300        //邮件检查有效期时间
	Min_Robot_PlayerId = 1770000000 //最小机器人id
	Max_Robot_PlayerId = 1870000000 //最大机器人id
)

//毫秒时间
const (
	ONE_HOUR_MILLSECOND    int64 = 3600000
	SIX_HOUR_MILLSECOND    int64 = 3600000 * 6
	TWELVE_HOUR_MILLSECOND int64 = 3600000 * 12
)

//文章评论状态
const (
	ARTICLE_COMMENT_SHOW int32 = 1 //可显示
	ARTICLE_COMMENT_HIDE int32 = 2 //隐藏 删除
)

//文章api请求数据类型
const (
	ARTICLE_API_READ        = "1" //读文章
	ARTICLE_API_ZAN         = "2" //点赞
	ARTICLE_API_COMMENT     = "3" //评论
	ARTICLE_API_GET_COMMENT = "4" //获取评论
)

//语音作品状态:0-未审核,1-已发布,3-已删除
const (
	VC_STATUS_UNCHECK int32 = 0 //0-未审核
	VC_STATUS_CHECKED int32 = 1 //已发布
	VC_STATUS_DELETE  int32 = 3 //已删除
)
const DAY_SECOND int64 = 86400 //毫秒
//const DAY_SECOND int64 = 120 //秒 600

//连接类型
const (
	CONN_TYPE_TCP       int = iota //TCP连接
	CONN_TYPE_WEBSOCKET            //websocket连接
)

//用户设备类型
const (
	TYPE_IOS     int32 = 1
	TYPE_ANDROID int32 = 2
	TYPE_PC      int32 = 3
)

//提示错误码:后续往后递增增加
const (
	FAIL_MSG_CODE_SUCCESS = "0000" //处理成功
	FAIL_MSG_CODE_1001    = "1001" //登录token错误
	FAIL_MSG_CODE_1002    = "1002" //黑名单消息拒收
	FAIL_MSG_CODE_1003    = "1003" //订单不存在
	FAIL_MSG_CODE_1004    = "1004" //金币改变失败
	FAIL_MSG_CODE_1005    = "1005" //不存在的玩家ID
	FAIL_MSG_CODE_1006    = "1006" //金币不足
	FAIL_MSG_CODE_1007    = "1007" //代付请求下单失败
	FAIL_MSG_CODE_1008    = "1008" //账号已被注册
	FAIL_MSG_CODE_1009    = "1009" //语音或者视频通话时 正好同时互相呼叫
	FAIL_MSG_CODE_1010    = "1010" //语音或者视频通话时 对方正在跟别人通话
	FAIL_MSG_CODE_1011    = "1011" //服务器停服中。。。
	FAIL_MSG_CODE_1012    = "1012" //银行卡信息不全
	FAIL_MSG_CODE_1013    = "1013" //账号注销中，继续登录将会取消注销
	FAIL_MSG_CODE_1014    = "1014" //登录过于频繁，请稍后继续尝试
	FAIL_MSG_CODE_1015    = "1015" //银联支付失败提示
	FAIL_MSG_CODE_1016    = "1016" // 转账/发送红包,账号已注销
	FAIL_MSG_CODE_1017    = "1017" // 语音匹配已达上线

)

//聊天文本内容类型
const (
	TALK_CONTENT_SYSTEM         int32 = iota //系统提示0
	TALK_CONTENT_WORD                        //文字1
	TALK_CONTENT_SOUND                       //语音2
	TALK_CONTENT_IMAGE                       //图片3
	TALK_CONTENT_REDPACKET                   //红包4
	TALK_CONTENT_TRANSFER_MONEY              //转账5
	TALK_CONTENT_TEAM_CARD                   //群名片6
	TALK_CONTENT_PERSONAL_CARD               //个人名片7
	TALK_CONTENT_REDPACKET_LOG               //领红包提示8
	TALK_CONTENT_WITHDRAW                    //撤回信息9
	TALK_CONTENT_AUDIO                       //语音通话10
	TALK_CONTENT_VIDEO                       //视屏通话11
	TALK_CONTENT_TRAMSFER_LOG                //领转账提示12
	TALK_CONTENT_CITE                        //引用13
	TALK_CONTENT_GROUPNOTICE                 //群公告14
	TALK_CONTENT_WEBARTICLE                  //web文章15
	TALK_CONTENT_VIDEO1                      //视频 16
	TALK_CONTENT_GROUPBAN                    //群禁言 17
	TALK_CONTENT_EMOTION                     //自定义表情包 18
	TALK_CONTENT_SAY_HI                      //打招呼格式 19
	TALK_CONTENT_SAY_HI_WORD                 //打招呼语 20
	TALK_CONTENT_ESP_VEDIO                   //电竞视频 21
	TALK_CONTENT_DYNAMIC                     //动态推送 22
	TALK_CONTENT_SHARE_TOPIC                 //分享话题23
)
const (
	TALK_STATUS_NORMAL int32 = iota //正常
	TALK_STATUS_DELETE              //删除
)

//默认违禁图
const (
	BAN_PICTURE_BIG   = "https://im-resource-1253887233.cos.accelerate.myqcloud.com/audit/weiguitu1.png"
	BAN_PICTURE_SMALL = "https://im-resource-1253887233.cos.accelerate.myqcloud.com/audit/weiguitu2.png"
)

//服务器类型
const (
	SERVER_TYPE_LOGIN              int32 = 1  //登录服
	SERVER_TYPE_HALL               int32 = 2  //大厅服
	SERVER_TYPE_BACKSTAGE          int32 = 3  //后台服
	SERVER_TYPE_SHOP               int32 = 4  //商场服
	SERVER_TYPE_STATISTICS         int32 = 5  //统计服
	SERVER_TYPE_SQUARE             int32 = 6  //社交广场服
	SERVER_TYPE_SPORT_APPLY        int32 = 8  //电竞接口服务器
	SERVER_TYPE_SPORT_CRAWL        int32 = 9  //电竞爬虫服务器
	SERVER_TYPE_SPORT_API          int32 = 10 //电竞第三方接口服务
	SERVER_TYPE_SPORT_LOTTERY_CSGO int32 = 11 //CSGO开奖
	SERVER_TYPE_SPORT_LOTTERY_WZRY int32 = 12 //王者荣耀开奖
	SERVER_TYPE_SPORT_LOTTERY_DOTA int32 = 13 //DOTA开奖
	SERVER_TYPE_SPORT_LOTTERY_LOL  int32 = 14 //LOL开奖
	SERVER_TYPE_WISH               int32 = 15 //许愿池
)

//红包类型
const (
	RED_PACKET_PERSONAL    int32 = 1 //私人红包
	RED_PACKET_TEAM_LUCKEY int32 = 2 //拼手气红包
	RED_PACKET_TEAM_NOMAL  int32 = 3 //普通红包
)

const (
	PAY_TYPE_GOLD          = 99 //零钱付款
	PAY_TYPE_WEIXIN        = 1  //微信支付
	PAY_TYPE_ZHIFUBAO      = 2  //支付宝支付
	PAY_TYPE_BANKCARD      = 3  //银行卡付款
	PAY_TYPE_BACKSTAGE_IN  = 4  //  后台操作入款
	PAY_TYPE_BACKSTAGE_OUT = 5  //  后台操作出款

)

//红包状态
const (
	PACKET_MONEY_OPEN    int32 = 1 //可领取
	PACKET_MONEY_FNISH   int32 = 2 //已领取
	PACKET_MONEY_TIMEOUT int32 = 3 //超时
)

//转账状态
const (
	TRANSFER_MONEY_OPEN  int32 = 1 //可领取
	TRANSFER_MONEY_FNISH int32 = 2 //已领取
	TRANSFER_MONEY_BACK  int32 = 3 //(主动/超时)退还
)

//注销账号状态
const (
	ACCOUNT_CANCEL_WAITING int32 = 0 //待处理
	ACCOUNT_CANCEL_FINISH  int32 = 1 //已完成
	ACCOUNT_CANCEL_REFULE  int32 = 2 //已拒绝
	ACCOUNT_CANCEL_CANCEL  int32 = 3 //已取消
)

//最小红包面值
const RED_PACKET_MIN_VALUE float64 = 1
const (
	GOLD_CHANGE_TYPE_IN  int32 = 1 //增加
	GOLD_CHANGE_TYPE_OUT int32 = 2 //减少
)

//充值方式
const (
	CHANNEL_SYSTEM   int32 = 0 //系统
	CHANNEL_MAN_MAKE int32 = 1 //人工
	CHANNEL_OTHER    int32 = 2 //第三方
)

//充值类型
const (
	PAY_TYPE_MONEY              int32 = 1 //零钱充值
	PAY_TYPE_REPACKET_PERSIONAL int32 = 2 //发个人红包充值
	PAY_TYPE_REPACKET_TEAM      int32 = 3 //发群红包充值
	PAY_TYPE_TRANSFER           int32 = 4 //转账充值
	PAY_TYPE_SHOP               int32 = 5 //商店购买充值
	PAY_TYPE_CODE               int32 = 6 //二维码付款充值
	PAY_TYPE_TEAMCODE           int32 = 7 //群二维码付款充值
	PAY_TYPE_COIN               int32 = 8 //充值兑换硬币
	PAY_TYPE_COIN_ITEM          int32 = 9 //充值购买硬币商品

)

func GetRechargeNote(t int32) string {
	switch t {
	case PAY_TYPE_MONEY:
		return "充值支付-零钱"
	case PAY_TYPE_REPACKET_PERSIONAL:
		return "充值支付-个人红包"
	case PAY_TYPE_REPACKET_TEAM:
		return "充值支付-群红包"
	case PAY_TYPE_TRANSFER:
		return "充值支付-转账"
	case PAY_TYPE_SHOP:
		return "充值支付-商场"
	case PAY_TYPE_CODE:
		return "充值支付-二维码"
	case PAY_TYPE_TEAMCODE:
		return "充值支付-群二维码"
	case PAY_TYPE_COIN:
		return "充值支付-硬币"
	case PAY_TYPE_COIN_ITEM:
		return "充值支付-硬币商品"
	default:
		return "充值支付"
	}

}

//支付渠道
const (
	PAY_CHANNEL_MIAODAO      int32 = 1  //秒到
	PAY_CHANNEL_TONGLIAN     int32 = 2  //通联	（微信）
	PAY_CHANNEL_PENGJU       int32 = 3  //鹏聚代付 （）
	PAY_CHANNEL_HUIJU        int32 = 4  //汇聚支付 （银联）
	PAY_CHANNEL_HUIJU_DF     int32 = 5  //汇聚代付付
	PAY_CHANNEL_HUICHAO_WX   int32 = 6  //汇潮支付(微信)
	PAY_CHANNEL_HUICHAO_ZFB  int32 = 7  //汇潮支付(支付宝)
	PAY_CHANNEL_HUICHAO_YL   int32 = 8  //汇潮支付(银联)
	PAY_CHANNEL_HUICHAO_DF   int32 = 9  //汇潮代付付
	PAY_CHANNEL_TONGTONG_WX  int32 = 10 //统统付微信支付
	PAY_CHANNEL_TONGTONG_ZFB int32 = 11 //统统付支付宝支付
	PAY_CHANNEL_YTS_WX       int32 = 12 //云通商微信支付
)

//Order充值订单状态
const (
	ORDER_ST_WAITTING int32 = 0 //等待处理
	ORDER_ST_FINISH   int32 = 1 //已完成
	ORDER_ST_AUDIT    int32 = 2 //已审核
	ORDER_ST_CANCEL   int32 = 3 //已取消
	ORDER_ST_REFUSE   int32 = 4 //拒绝
)

//青少年保护模式开关
const (
	YOUNG_STATUS_OPEN  = 1 //开
	YOUNG_STATUS_CLOSE = 2 //关
)

//PayStatus支付状态
const (
	PAY_ST_WAITTING int32 = 0 //待支付
	PAY_ST_FINISH   int32 = 1 //已支付
	PAY_ST_CANCEL   int32 = 2 //已取消
	PAY_ST_DOING    int32 = 3 //第三方处理中
	PAY_ST_REFUSE   int32 = 4 //已超时
	PAY_ST_FAIL     int32 = 5 //支付失败
)

// OrderType 订单类型
const (
	ORDER_ST_IM   int32 = 0 // IM
	ORDER_ST_WISH int32 = 1 // 许愿池
)

const (
	LOGINREGISTER_PASSWDLOGIN    = 1 //密码登录
	LOGINREGISTER_MESSAGELOGIN   = 2 //验证码登录
	LOGINREGISTER_ONEKEYLOGIN    = 4 //一键登录
	LOGINREGISTER_WECHATLOGIN    = 5 //微信登录
	LOGINREGISTER_AUTOLOGIN      = 6 //自动登录
	LOGINREGISTER_PHONEREGISTER  = 7 //手机号码注册
	LOGINREGISTER_ONEKEYREGISTER = 8 //一键登录注册
	LOGINREGISTER_WECHATREGISTER = 9 //微信登录注册
)

//金币类型
const (
	//入款项
	GOLD_TYPE_CASH_AFIN           int32 = 100 //人工入款
	GOLD_TYPE_CASH_IN             int32 = 101 //充值
	GOLD_TYPE_GET_REDPACKET       int32 = 102 //收红包
	GOLD_TYPE_GET_TRANSFER_MONEY  int32 = 103 //转入
	GOLD_TYPE_GET_MONEY           int32 = 104 //二维码收款
	GOLD_TYPE_REDPACKET_OVERTIME  int32 = 111 //红包退款
	GOLD_TYPE_TRANSFER_MONEY_OVER int32 = 112 //转账退款
	GOLD_TYPE_BACK_MONEY          int32 = 113 //商家退款
	GOLD_TYPE_CASH_OUT_BACK       int32 = 114 //取消提现退款
	GOLD_TYPE_SHOP_ITEM_MONEY     int32 = 115 //商城卖家货款

	//出款项
	GOLD_TYPE_CASH_AFOUT          int32 = 200 //人工出款
	GOLD_TYPE_CASH_OUT            int32 = 201 //提现
	GOLD_TYPE_SEND_REDPACKET      int32 = 202 //发红包
	GOLD_TYPE_SEND_TRANSFER_MONEY int32 = 203 //转出
	GOLD_TYPE_PAY_MONEY           int32 = 204 //二维码付款
	GOLD_TYPE_FINE_MONEY          int32 = 215 //罚没
	GOLD_TYPE_EXTRA_MONEY         int32 = 216 //手续费
	GOLD_TYPE_SHOP_MONEY          int32 = 217 //商城消费
	GOLD_TYPE_EXCHANGE_COIN       int32 = 219 //兑换硬币
	GOLD_TYPE_XN_SHOP_MONEY       int32 = 220 //虚拟商城消费

	COIN_TYPE_SYSTEM_IN         int32 = 500 //系统赠送
	COIN_TYPE_EXCHANGE_IN       int32 = 501 //兑换
	COIN_TYPE_PLAYER_IN         int32 = 502 //被投币
	COIN_TYPE_ACT_PRIZE         int32 = 503 //活动奖励
	COIN_TYPE_WISH_ADD          int32 = 505 //许愿池回收
	COIN_TYPE_WISH_PRO_ADD      int32 = 506 //许愿池守护者收益
	COIN_TYPE_WISH_PLATFORM_ADD int32 = 507 // 许愿池平台回收.
	COIN_TYPE_WISH_DARE_BACK    int32 = 508 // 许愿池抽奖返利

	COIN_TYPE_SYSTEM_OUT          int32 = 600 //系统回收
	COIN_TYPE_SHOP_OUT            int32 = 601 //商场消费
	COIN_TYPE_PLAYER_OUT          int32 = 602 //投币
	COIN_TYPE_EXPIRED_OUT         int32 = 603 //过期回收
	COIN_TYPE_CONFISCATE_OUT      int32 = 604 //系统罚没
	COIN_TYPE_WISH_PAY            int32 = 605 // 许愿池兑换钻石
	COIN_TYPE_ESPORT_EXCHANGE_OUT int32 = 606 //硬币电竞兑换
	// 许愿池 700   800
	DIAMOND_TYPE_EXCHANGE_IN      int32 = 700 // 兑换
	DIAMOND_TYPE_WISH_DARE_BACK   int32 = 701 // 抽奖返利
	DIAMOND_TYPE_WISH_FAILD       int32 = 702 // 许愿池抽奖失败钻石返回
	DIAMOND_TYPE_WISH_GUARDIAN_IN int32 = 703 // 守护者收益
	DIAMOND_TYPE_WISH_BACK        int32 = 704 // 回收
	DIAMOND_TYPE_POSTAGE_FAILD    int32 = 705 // 邮费返回
	DIAMOND_TYPE_BACK_GIVE        int32 = 706 // 后台赠送
	DIAMOND_TYPE_WISH_ACT         int32 = 707 // 许愿池活动所得

	DIAMOND_TYPE_PLAYER_OUT      int32 = 800 // 抽奖
	DIAMOND_TYPE_POSTAGE_OUT     int32 = 801 // 运费
	DIAMOND_TYPE_WISH_BACK_FAILD int32 = 802 // 回收失败返回
	DIAMOND_TYPE_BACK_RECYCLE    int32 = 803 // 后台扣除

	ESPORTCOIN_TYPE_GUESS_BACK_IN    int32 = 900 //竞猜返还
	ESPORTCOIN_TYPE_EXCHANGE_IN      int32 = 901 //电竞币电竞兑换
	ESPORTCOIN_TYPE_EXCHANGE_GIVE_IN int32 = 902 //电竞币赠送

	ESPORTCOIN_TYPE_GUESS_BET_OUT int32 = 1000 //竞猜投注

)

//虚拟商场相关定义
const (
	//道具类型
	COIN_PROPS_TYPE_LB   int32 = 1 //礼包
	COIN_PROPS_TYPE_GJ   int32 = 2 //挂件
	COIN_PROPS_TYPE_QP   int32 = 3 //气泡
	COIN_PROPS_TYPE_MP   int32 = 4 //铭牌
	COIN_PROPS_TYPE_QTX  int32 = 5 //群特效
	COIN_PROPS_TYPE_MZBS int32 = 6 //名字变色

	//道具使用类型
	COIN_PROPS_USETYPE_COMSUME   int32 = 1 //消耗型
	COIN_PROPS_USETYPE_EQUIPMENT int32 = 2 //装备型

	//支付类型
	COIN_PROPS_BUYWAY_COIN  int32 = 1 //硬币购买
	COIN_PROPS_BUYWAY_MONEY int32 = 2 //零钱购买

	COIN_PROPS_FOREVER = -1 //永久道具

	//背包道具状态
	COIN_BAG_ITEM_UNUSE   = 1 //未使用
	COIN_BAG_ITEM_USED    = 2 //使用中
	COIN_BAG_ITEM_EXPIRED = 3 //过期
	//道具获得类型 1-购买 2-系统赠送(做任务获得);3-玩家赠送;4-活动获得
	COIN_ITEM_GETTYPE_BUY         = 1 //购买
	COIN_ITEM_GETTYPE_SEND        = 2 //系统赠送
	COIN_ITEM_GETTYPE_PLAYER_SEND = 3 //玩家赠送
	COIN_ITEM_GETTYPE_ACTIVITY    = 4 //活动获得
	COIN_ITEM_GETTYPE_BACK        = 5 //系统回收

	COIN_EXPIRATION_TIME = 90 * 86400 //绑定硬币90天过期
	// COIN_EXPIRATION_TIME = 86400 //绑定硬币90天过期

	BCOIN_STATUS_UNUSE      = 0 //未使用
	BCOIN_STATUS_USED       = 1 //已使用
	BCOIN_STATUS_EXPIRATION = 2 //已过期

	//虚拟商品状态
	COIN_PRODUCT_STATUS_UP   = 1 //上架
	COIN_PRODUCT_STATUS_DOWN = 2 //下架
	COIN_PRODUCT_STATUS_DEL  = 3 //删除
)

const (
	ADV_LOG_OP_TYPE_1 = 1 // 展示次数
	ADV_LOG_OP_TYPE_2 = 2 // 展示人数
	ADV_LOG_OP_TYPE_3 = 3 // 点击次数
	ADV_LOG_OP_TYPE_4 = 4 // 点击人数
)

const (
	APK_CODE_ANDROID = 100
	APK_CODE_IOS     = 101
)

const (
	MAX_SAY_MESSAGE_NUM = 100              // 最多100条
	TIME_BEFORE         = 3 * 86400 * 1000 // 3天前
)

var BankName = map[string]string{
	"SRCB":      "深圳农村商业银行",
	"BGB":       "广西北部湾银行",
	"SHRCB":     "上海农村商业银行",
	"BJBANK":    "北京银行",
	"WHCCB":     "威海市商业银行",
	"BOZK":      "周口银行",
	"KORLABANK": "库尔勒市商业银行",
	"SPABANK":   "平安银行",
	"SDEB":      "顺德农商银行",
	"HURCB":     "湖北省农村信用社",
	"WRCB":      "无锡农村商业银行",
	"BOCY":      "朝阳银行",
	"CZBANK":    "浙商银行",
	"HDBANK":    "邯郸银行",
	"BOC":       "中国银行",
	"BOD":       "东莞银行",
	"CCB":       "中国建设银行",
	"ZYCBANK":   "遵义市商业银行",
	"SXCB":      "绍兴银行",
	"GZRCU":     "贵州省农村信用社",
	"ZJKCCB":    "张家口市商业银行",
	"BOJZ":      "锦州银行",
	"BOP":       "平顶山银行",
	"HKB":       "汉口银行",
	"SPDB":      "上海浦东发展银行",
	"NXRCU":     "宁夏黄河农村商业银行",
	"NYNB":      "广东南粤银行",
	"GRCB":      "广州农商银行",
	"BOSZ":      "苏州银行",
	"HZCB":      "杭州银行",
	"HSBK":      "衡水银行",
	"HBC":       "湖北银行",
	"JXBANK":    "嘉兴银行",
	"HRXJB":     "华融湘江银行",
	"BODD":      "丹东银行",
	"AYCB":      "安阳银行",
	"EGBANK":    "恒丰银行",
	"CDB":       "国家开发银行",
	"TCRCB":     "江苏太仓农村商业银行",
	"NJCB":      "南京银行",
	"ZZBANK":    "郑州银行",
	"DYCB":      "德阳商业银行",
	"YBCCB":     "宜宾市商业银行",
	"SCRCU":     "四川省农村信用",
	"KLB":       "昆仑银行",
	"LSBANK":    "莱商银行",
	"YDRCB":     "尧都农商行",
	"CCQTGB":    "重庆三峡银行",
	"FDB":       "富滇银行",
	"JSRCU":     "江苏省农村信用联合社",
	"JNBANK":    "济宁银行",
	"CMB":       "招商银行",
	"JINCHB":    "晋城银行JCBANK",
	"FXCB":      "阜新银行",
	"WHRCB":     "武汉农村商业银行",
	"HBYCBANK":  "湖北银行宜昌分行",
	"TZCB":      "台州银行",
	"TACCB":     "泰安市商业银行",
	"XCYH":      "许昌银行",
	"CEB":       "中国光大银行",
	"NXBANK":    "宁夏银行",
	"HSBANK":    "徽商银行",
	"JJBANK":    "九江银行",
	"NHQS":      "农信银清算中心",
	"MTBANK":    "浙江民泰商业银行",
	"LANGFB":    "廊坊银行",
	"ASCB":      "鞍山银行",
	"KSRB":      "昆山农村商业银行",
	"YXCCB":     "玉溪市商业银行",
	"DLB":       "大连银行",
	"DRCBCL":    "东莞农村商业银行",
	"GCB":       "广州银行",
	"NBBANK":    "宁波银行",
	"BOYK":      "营口银行",
	"SXRCCU":    "陕西信合",
	"GLBANK":    "桂林银行",
	"BOQH":      "青海银行",
	"CDRCB":     "成都农商银行",
	"QDCCB":     "青岛银行",
	"HKBEA":     "东亚银行",
	"HBHSBANK":  "湖北银行黄石分行",
	"WZCB":      "温州银行",
	"TRCB":      "天津农商银行",
	"QLBANK":    "齐鲁银行",
	"GDRCC":     "广东省农村信用社联合社",
	"ZJTLCB":    "浙江泰隆商业银行",
	"GZB":       "赣州银行",
	"GYCB":      "贵阳市商业银行",
	"CQBANK":    "重庆银行",
	"DAQINGB":   "龙江银行",
	"CGNB":      "南充市商业银行",
	"SCCB":      "三门峡银行",
	"CSRCB":     "常熟农村商业银行",
	"SHBANK":    "上海银行",
	"JLBANK":    "吉林银行",
	"CZRCB":     "常州农村信用联社",
	"BANKWF":    "潍坊银行",
	"ZRCBANK":   "张家港农村商业银行",
	"FJHXBC":    "福建海峡银行",
	"ZJNX":      "浙江省农村信用社联合社",
	"LZYH":      "兰州银行",
	"JSB":       "晋商银行",
	"BOHAIB":    "渤海银行",
	"CZCB":      "浙江稠州商业银行",
	"YQCCB":     "阳泉银行",
	"SJBANK":    "盛京银行",
	"XABANK":    "西安银行",
	"BSB":       "包商银行",
	"JSBANK":    "江苏银行",
	"FSCB":      "抚顺银行",
	"HNRCU":     "河南省农村信用",
	"COMM":      "交通银行",
	"XTB":       "邢台银行",
	"CITIC":     "中信银行",
	"HXBANK":    "华夏银行",
	"HNRCC":     "湖南省农村信用社",
	"DYCCB":     "东营市商业银行",
	"ORBANK":    "鄂尔多斯银行",
	"BJRCB":     "北京农村商业银行",
	"XYBANK":    "信阳银行",
	"ZGCCB":     "自贡市商业银行",
	"CDCB":      "成都银行",
	"HANABANK":  "韩亚银行",
	"CMBC":      "中国民生银行",
	"LYBANK":    "洛阳银行",
	"GDB":       "广东发展银行",
	"ZBCB":      "齐商银行",
	"CBKF":      "开封市商业银行",
	"H3CB":      "内蒙古银行",
	"CIB":       "兴业银行",
	"CRCBANK":   "重庆农村商业银行",
	"SZSBK":     "石嘴山银行",
	"DZBANK":    "德州银行",
	"SRBANK":    "上饶银行",
	"LSCCB":     "乐山市商业银行",
	"JXRCU":     "江西省农村信用",
	"ICBC":      "中国工商银行",
	"JZBANK":    "晋中市商业银行",
	"HZCCB":     "湖州市商业银行",
	"NHB":       "南海农村信用联社",
	"XXBANK":    "新乡银行",
	"JRCB":      "江苏江阴农村商业银行",
	"YNRCC":     "云南省农村信用社",
	"ABC":       "中国农业银行",
	"GXRCU":     "广西省农村信用",
	"PSBC":      "中国邮政储蓄银行",
	"BZMD":      "驻马店银行",
	"ARCU":      "安徽省农村信用社",
	"GSRCU":     "甘肃省农村信用",
	"LYCB":      "辽阳市商业银行",
	"JLRCU":     "吉林农信",
	"URMQCCB":   "乌鲁木齐市商业银行",
	"XLBANK":    "中山小榄村镇银行",
	"CSCB":      "长沙银行",
	"JHBANK":    "金华银行",
	"BHB":       "河北银行",
	"NBYZ":      "鄞州银行",
	"LSBC":      "临商银行",
	"BOCD":      "承德银行",
	"SDRCU":     "山东农信",
	"NCB":       "南昌银行",
	"TCCB":      "天津银行",
	"WJRCB":     "吴江农商银行",
	"CBBQS":     "城市商业银行资金清算中心",
	"HBRCU":     "河北省农村信用社",
}

//银行代号和支付数字编号
var BankPayNo = map[string]string{
	"SRCB":      "314",
	"SHRCB":     "314",
	"WHCCB":     "313",
	"KORLABANK": "313",
	"SDEB":      "314",
	"HURCB":     "402",
	"WRCB":      "314",
	"CZBANK":    "316",
	"BOC":       "104",
	"CCB":       "105",
	"ZYCBANK":   "313",
	"GZRCU":     "402",
	"ZJKCCB":    "313",
	"SPDB":      "310",
	"NXRCU":     "314",
	"GRCB":      "314",
	"EGBANK":    "315",
	"CDB":       "201",
	"TCRCB":     "314",
	"DYCB":      "313",
	"YBCCB":     "313",
	"SCRCU":     "402",
	"YDRCB":     "314",
	"JSRCU":     "402",
	"CMB":       "308",
	"WHRCB":     "314",
	"TACCB":     "313",
	"CEB":       "303",
	"NHQS":      "402",
	"MTBANK":    "313",
	"KSRB":      "314",
	"YXCCB":     "313",
	"DRCBCL":    "314",
	"CDRCB":     "314",
	"HKBEA":     "502",
	"TRCB":      "314",
	"GDRCC":     "402",
	"ZJTLCB":    "313",
	"GYCB":      "313",
	"CGNB":      "313",
	"CSRCB":     "314",
	"CZRCB":     "402",
	"ZRCBANK":   "314",
	"CZCB":      "313",
	"HNRCU":     "402",
	"COMM":      "301",
	"CITIC":     "302",
	"HXBANK":    "304",
	"HNRCC":     "402",
	"DYCCB":     "313",
	"BJRCB":     "314",
	"ZGCCB":     "313",
	"CMBC":      "305",
	"GDB":       "306",
	"CBKF":      "313",
	"CIB":       "309",
	"CRCBANK":   "314",
	"LSCCB":     "313",
	"JXRCU":     "402",
	"ICBC":      "102",
	"JZBANK":    "313",
	"HZCCB":     "313",
	"NHB":       "402",
	"JRCB":      "314",
	"YNRCC":     "402",
	"ABC":       "103",
	"GXRCU":     "402",
	"PSBC":      "403",
	"GSRCU":     "402",
	"LYCB":      "313",
	"JLRCU":     "402",
	"URMQCCB":   "313",
	"SDRCU":     "402",
	"WJRCB":     "314",
	"CBBQS":     "313",
	"HBRCU":     "402",
}

//入账集合
var InSouceType = []int32{GOLD_TYPE_CASH_IN, GOLD_TYPE_GET_REDPACKET, GOLD_TYPE_REDPACKET_OVERTIME, GOLD_TYPE_GET_TRANSFER_MONEY, GOLD_TYPE_TRANSFER_MONEY_OVER, GOLD_TYPE_GET_MONEY, GOLD_TYPE_BACK_MONEY, GOLD_TYPE_SHOP_ITEM_MONEY}

//出账集合
var OutSouceType = []int32{GOLD_TYPE_CASH_OUT, GOLD_TYPE_SEND_REDPACKET, GOLD_TYPE_SEND_TRANSFER_MONEY, GOLD_TYPE_PAY_MONEY, GOLD_TYPE_FINE_MONEY, GOLD_TYPE_EXTRA_MONEY}

const (
	_GOLD_ATTR int32 = iota
	_NICK_NAME_ATTR
	_HEAD_ICON_ATTR
	VIP_LEVEL_ATTR
	_SAFE_BOX_ATTR
	_ADDRESS_ATTR
	_IP_ATTR
)
const (
	SMS_TYPE_CODE    = 1 //验证码类型模板
	SMS_TYPE_MESSAGE = 2 //消息类型模板
)
const (
	CLIENT_CODE_LOGIN         = 1 //登录发送验证码
	CLIENT_CODE_REGISTER      = 2 //注册发送验证码
	CLIENT_CODE_PAYPASSWORD   = 3 //设置支付密码短信验证
	CLIENT_CODE_FORGETLOGIN   = 4 //忘记登录密码短信验证
	CLIENT_CODE_BINDBANK      = 5 //绑定银行卡短信验证
	CLIENT_CODE_PLAYERMESSAGE = 6 //第一次登陆完善信息
	CLIENT_CODE_BINDPHONE     = 7 //绑定手机验证码
	CLIENT_CODE_CANCELACCOUNT = 8 //注销账号短信验证码
)

//账号状态
const (
	ACCOUNT_NORMAL       = 0 //正常
	ACCOUNT_USER_FROZEN  = 1 //用户冻结
	ACCOUNT_ADMIN_FROZEN = 2 //后台冻结
	ACCOUNT_CANCELING    = 3 //注销中
	ACCOUNT_CANCELED     = 4 //已注销
)

//账号类型
const (
	ACCOUNT_TYPES_PT   = 1 //1普通用户
	ACCOUNT_TYPES_YXYY = 2 //营销运营
	ACCOUNT_TYPES_SC   = 3 //商城
	ACCOUNT_TYPES_GLYY = 4 //管理运营
	ACCOUNT_TYPES_GFYY = 5 //官方运营
	ACCOUNT_TYPES_CSYY = 6
)

//用户类型 1普通用户,2营销运营,3商城账号,4管理运营,5官方运营
const (
	PLAYER_NORMAL   = 1
	PLAYER_MARKET   = 2
	PLAYER_SHOP     = 3
	PLAYER_MANAGE   = 4
	PLAYER_OFFICIAL = 5
)

const (
	CHAT_TYPE_PRIVATE = 1 //私聊
	CHAT_TYPE_TEAM    = 2 //群聊
	CHAT_TYPE_GROUP   = 3 //讨论组
	CHAT_TYPE_NEARBY  = 4 //附近的人
)

//聊天失败原因
const (
	CHAT_REFUSE_TYPE_1  = 1  //陌生人发送失败
	CHAT_REFUSE_TYPE_2  = 2  //黑名单发送失败
	CHAT_REFUSE_TYPE_3  = 3  //群禁言
	CHAT_REFUSE_TYPE_4  = 4  //敏感词
	CHAT_REFUSE_TYPE_5  = 5  //账号已注销
	CHAT_REFUSE_TYPE_6  = 6  //陌生人说话超3次
	CHAT_REFUSE_TYPE_7  = 7  //不允许陌生人打招呼
	CHAT_REFUSE_TYPE_8  = 8  //不允许群聊打招呼
	CHAT_REFUSE_TYPE_9  = 9  //不允许二维码打招呼
	CHAT_REFUSE_TYPE_10 = 10 //不允许名片打招呼
	CHAT_REFUSE_TYPE_11 = 11 //你的好友达上限
	CHAT_REFUSE_TYPE_12 = 12 //对方好友达上限
	CHAT_REFUSE_TYPE_13 = 13 //陌生人说话次数大于1次，小于3次
)
const (
	TYPE_REPORT_NOTICE   = 0 //广告公告
	TYPE_ACTIVITY_NOTICE = 1 //活动公告
	TYPE_GAME_REPORT     = 2 //游戏公告
)

const (
	ORDER_NODEAL   = 0 //0未处理，1待审核(待放款)，2已完成，3已取消，4已拒绝
	ORDER_WAIT     = 1
	ORDER_COMPLETE = 2
	ORDER_CANCEL   = 3
	ORDER_REFUSE   = 4
)

//给大厅发送类型
const (
	SEND_ALL_HALL = 1 //发送所有大厅
	SEND_ONE_HALL = 2 //随机发送一台大厅
)

const (
	SHOP_ORDER_WAIT_PAY         = 0 //  商城订单状态  0待付款 1超时 2取消 3待发货 4待收货 5已完成 6评价 7后台取消
	SHOP_ORDER_EXPIRE           = 1
	SHOP_ORDER_CANCEL           = 2
	SHOP_ORDER_WAIT_SEND        = 3
	SHOP_ORDER_WAIT_RECEIVE     = 4
	SHOP_ORDER_FINISH           = 5
	SHOP_ORDER_EVALUTE          = 6
	SHOP_ORDER_BACKSTAGE_CANCLE = 7
)

const (
	SHOP_ITEM_SALE       = 0 //商城物品状态 0 上架 1下架 2删除 3审核中 4审核失败
	SHOP_ITEM_SOLD_OUT   = 1
	SHOP_ITEM_DELETE     = 2
	SHOP_ITEM_IN_AUDIT   = 3
	SHOP_ITEM_FAIL_AUDIT = 4
)

const (
	SHOP_COMMENT_NO_REPLY = 0 //商城留言状态 0 未恢复 1恢复 2删除
	SHOP_COMMENT_REPLY    = 1
	SHOP_COMMENT_DELETE   = 2
)

const (
	SHOP_COMMENT_LEVEL_COMMON = 0 //商城留言等级 0 普通留言 1差评 2中评 3好评
	SHOP_COMMENT_LEVEL_BAD    = 1
	SHOP_COMMENT_LEVEL_MID    = 2
	SHOP_COMMENT_LEVEL_GOOD   = 3
)

//商城api请求数据类型
const (
	SHOP_API_GOODS       = "1" //查询商品详情
	SHOP_API_OTHER_GOODS = "2" //查询同一个卖家的同类商品
	SHOP_API_PLACE_ORDER = "3" //下订单
	SHOP_API_PAY         = "4" //支付
)

//商城api请求数据类型
const (
	SHOP_API_ITEM_DETAIL  = "1" //查询商品详情数据包装(链接打开等动作)
	SHOP_API_NOW_BUY      = "2" //立即购买动作数据包装
	SHOP_API_FIRST_PAY    = "3" //抢先支付按钮数据包装(创建用户信息以及创建订单动作,这个两个动作放在一起)
	SHOP_API_SERACH       = "4" //搜索按钮
	SHOP_API_ORDER_DETAIL = "5" //订单详情
)

//商城点卡状态
const (
	SHOP_POINT_CARD_SALE    = 1 //待售
	SHOP_POINT_CARD_SELLOUT = 2 //售罄
)

//商城点卡小分类的值45
const (
	SHOP_POINT_CARD_CATEGORY = 45 //商城点卡小分类的值45
)

const (
	//商城订单id分布式生成用的redis key
	SHOP_CREATE_ORDER_ID = "redis_shop:create_order_id"

	//商城在是商城待付款的时候各个物品之间的锁(需要用重试锁),取得锁的时候以物品为单位,需要加上物品id
	//付款的时候
	SHOP_ITEM_PAY_MUTEX = "redis_shop:item_pay_mutex"

	//商城在是商城待付款的时候各个订单之间的锁(不重试锁),取得锁的时候以订单为单位,需要加上订单id
	//取消,超时取消,后台取消
	//(其实这些操作跟支付是竞争的,目前如果并发的时候优先支付完成,支付那里不锁这个锁)
	SHOP_ORDER_WAIT_PAY_MUTEX = "redis_shop:order_wait_pay_mutex"

	//支付中、以及重复支付用
	SHOP_ORDER_PAYING_MUTEX = "redis_shop:order_paying_mutex"

	//商城在是收货,延长收货和自动收货的时候各个订单之间的锁(不重试锁),取得锁的时候以订单为单位,需要加上订单id
	SHOP_ORDER_RECEIVE_MUTEX = "redis_shop:order_receive_mutex"

	//商城在是发货,后台发货的时候各个订单之间的锁(不重试锁),取得锁的时候以订单为单位,需要加上订单id
	SHOP_ORDER_SEND_MUTEX = "redis_shop:order_send_mutex"

	//商城自动收货 集群下服务器之间的不重试锁,保证同一时间在各个商城服务器之间只有一台在运行
	SHOP_AUTO_RECEIVE_MUTEX_SERVER = "redis_shop:auto_receive_mutex_server"
)

const ( //群消息通知类型
	TIME_CLEAR             = 1  //定时清理
	READ_CLEAR             = 2  //阅后
	SCREENSHOT_NOTICE      = 3  //截屏通知
	STOP_TALK              = 4  //全员禁言
	STOP_ADDFRIEND         = 5  //禁止群成员互加好友
	TEAM_INVITE            = 6  //群聊邀请确认
	TEAM_NAME              = 7  //群名称
	TEAM_GONGGAO           = 8  //群公告
	INVITE_PLAYER          = 9  //邀请入群
	DEL_PLAYER             = 10 //踢出群
	ADD_MANAGER            = 11 //增加管理員
	DEL_MANAGER            = 12 //刪除管理員
	STOP_REDPACKET         = 13 //禁止领取零钱红包
	CHANGE_OWNER           = 14 //转让群主
	EXIT_PLAYER            = 15 //退出群
	ACTIVE_ADDTEAM         = 16 //主动进群
	TEAM_RECOMMEND         = 17 //群推荐
	REQUEST_ADDTEAM        = 18 //申请进群
	CHANGE_TEAMMONEYCODE   = 19 //修改群收款二维码
	WITHDRAW_MESSAGE       = 20 //撤回信息
	STOP_ADDTEAM           = 21 //禁止主动进群
	BACKSTAGE_BAN_TEAM     = 22 //群封禁
	BACKSTAGE_BAN_TEAM_MEM = 23 //群成员封禁
	ADV_TEAM_MEM           = 24 // 点击广告连接进群
	WELCOME_WORD           = 25 //群欢迎语开关
	WELCOME_WORD_MANAGER   = 26 //群欢迎语管理员权限开关
	EDIT_WELCOME_WORD      = 27 //编辑群欢迎语
	TEAM_HEAD              = 28 // 群头像设置.
	TOPIC_TEAM_DESC        = 29 // 话题群简介
)

const (
	NORMAL   = 0 //正常
	DISSOLVE = 1 //int64 解散
	BANNED   = 2 //封禁
)

const (
	SYSTEM  = 1 //系统
	OWNER   = 2 //群主
	MANAGER = 3 //群管理员
)

// 发送短信验证码类型
const (
	SmsTypeLogin               = 1 // 登录发送验证码
	SmsTypeRegister            = 2 // 注册发送验证码
	SmsTypeSetPayPassword      = 3 // 设置支付密码短信验证
	SmsTypeForgotLoginPassword = 4 // 忘记登录密码短信验证
	SmsTypeBindBankCard        = 5 // 绑定银行卡短信验证
	SmsTypeFirstLogin          = 6 // 第一次登陆完善信息
)

// 系统参数表中id常量
const (
	AVATAR_PARAMETER    = "avatar_parameter"    // 头像参数
	INTEREST_PARAMETER  = "interest_parameter"  // 兴趣标签
	LIMIT_PARAMETER     = "limit_parameter"     // 系统功能转账参数
	OBJ_MODERATIONS     = "obj_moderations"     // 系统功能屏蔽控制参数
	SQUAREHOT_PARAMETER = "squarehot_parameter" // 动态热门参数
	WARNING_PARAMETER   = "warning_parameter"   // 支付预警参数
	TOPICHOT_PARAMETER  = "topichot_parameter"  // 话题热门参数
	PUSH_PARAMETER      = "push_parameter"      //极光推送管理
	ESPORT_PARAMETER    = "esport_parameter"    //电竞系统控制参数
	COMMON_PARAMETER    = "common_parameter"    //通用参数
)

const (
	NEAR_MESSAGE_NORMAL = 1 // 正常的附近的人消息
	NEAR_MESSAGE_DELETE = 2 //删除的附近的人消息
)

const (
	//虚拟商城
	PUSH_ITEM_101 int32 = 101 //硬币过期提醒-您有平台赠送硬币明日0点即将过期
	PUSH_ITEM_102 int32 = 102 //道具过期提醒-您的%v即将过期,请及时查看
	//社交广场
	PUSH_ITEM_201 int32 = 201 //点赞我的动态-%v点赞了我的社交广场动态
	PUSH_ITEM_202 int32 = 202 //评论我的动态-%v评论了我的社交广场动态
	PUSH_ITEM_203 int32 = 203 //回复我的动态-%v回复了我在社交广场的评论
	PUSH_ITEM_204 int32 = 204 //我关注的人发布新动态-%v发布了新的社交广场动态
)

const (
	/*
		1-协议页面同意按钮点击
		2-协议页面不同意按钮点击
		3-手机登录注册按钮点击次数
		4-微信登录注册点击次数
		5-本机号码一键登录次数
		6-其他号码登录点击次数
		7-注册登录页2返回键次数
		8-获取验证码次数
		9-重新获取验证码按钮点击次数
		10-兴趣墙确定按钮点击次数 s大厅服务器
		11-兴趣墙返回键点击次数 s
		12-推荐页面跳过按钮点击次数 s
		13-推荐页面下一步按钮点击次数 s
		14-进入柠檬畅聊按钮 s
		15-输入手机号页返回次数
		16-输入验证码页返回次数
		17-信息页返回次数
	*/
	WelcomeAgree            int32 = 1
	WelcomeNoAgree          int32 = 2
	PhoneLoginPV            int32 = 3
	WeixinLoginPV           int32 = 4
	OneClickLoginPV         int32 = 5
	OtherNumberLoginPV      int32 = 6
	LoginPage2ReturnPV      int32 = 7
	VerificationCodePV      int32 = 8
	VerificationCodePvAgain int32 = 9
	InterestSurePV          int32 = 10
	InterestReturnPV        int32 = 11
	RecommendSkipPV         int32 = 12
	RecommendNextPV         int32 = 13
	InNmBtn                 int32 = 14
	InPhoneBack             int32 = 15
	InCodeBack              int32 = 16
	InfoBack                int32 = 17
)
const (
	VC_CARD_TYPE_NOMARL       = 0 //正常卡片
	VC_CARD_TYPE_VOICE        = 1 //需要补充录音
	VC_CARD_TYPE_PERSIONALTAG = 2 //需要补充个性标签
	VC_CARD_TYPE_ADV          = 3 //广告
)

//恋爱匹配交友关注方式
const (
	VC_ATTENTION_LIKE = 1 // 关注类型：点赞
	VC_ATTENTION_HI   = 2 //关注类型：SayHi
)
const VC_MAX_ZAN_NUM = 6 //最大点赞数

//恋爱匹配交友关注对象
const (
	VC_ATTENTION_TO_ME    = 1 // 关注我
	VC_ATTENTION_TO_OTHER = 2 //关注其他人
)

//音频短片类型
const (
	VC_VOIDE_ALL       = 0 // 全部
	VC_VOIDE_SOLILOQUY = 1 // 独白
	VC_VOIDE_DUB       = 2 //电影配音
	VC_VOIDE_SING      = 3 //唱一唱
)

//搜索音频短片类型
const (
	VC_VOIDE_NAME    = 1 //作品名
	VC_VOIDE_CONTENT = 2 //台本
	VC_VOIDE_AUTHOR  = 3 //作者
)

//搜索音频短片类型
const (
	VC_VOIDE_AUDIT      = 0 //未审核
	VC_VOIDE_PASS_AUDIT = 1 //通过
	VC_VOIDE_REFUSED    = 2 //拒绝
	VC_VOIDE_DEL        = 3 //删除

)

// 有些表是联合主键,需要取出来
func TransPrimaryKey(primaryKey PRIMARY_KEY) []interface{} {
	if o, ok := primaryKey.(IPrimaryKey); ok {
		return o.Keys()
	} else {
		return []interface{}{primaryKey}
	}
}

// 用得太多了，包装一下
func RpcToast(ep client_server.IServer2Client, text string) {
	// logs.Debug("RpcToastAndFail:%s", text)
	msg := &client_server.ToastMsg{Text: &text}
	if ep != nil {
		ep.RpcToast(msg)
	}
}

// 用得太多了，包装一下
func RpcToastAndFail(ep client_server.IServer2Client, text string) *base.Fail {
	// logs.Debug("RpcToastAndFail:%s", text)
	msg := &client_server.ToastMsg{Text: &text}
	if ep != nil {
		ep.RpcToast(msg)
	}
	return easygo.FailMsg
}

func RpcToastAndPanic(err error, ep client_server.IServer2Client, text string) {
	if err != nil {
		// logs.Debug("RpcToastAndFail:%s", text)
		msg := &client_server.ToastMsg{Text: &text}
		if ep != nil {
			ep.RpcToast(msg)
		}
		panic(err)
	}
}

func GetRandAccount(head string, id int64) string {
	center := strconv.FormatInt(id, 16) //10 yo 16
	return head + center
}

//md5字符串
func Md5(d string) string {
	tokenMd5 := md5.New()
	tokenMd5.Write([]byte(d))
	token := hex.EncodeToString(tokenMd5.Sum(nil)) //
	return token
}

//生成密码
func CreatePasswd(pass string, slat string) string {
	str := pass + slat
	h := md5.New()
	io.WriteString(h, slat)
	io.WriteString(h, str)
	passwd := fmt.Sprintf("%x", h.Sum(nil))
	return passwd
}

//生成随机字符串
func RandString(n int) string {
	t := time.Now()
	h := md5.New()
	io.WriteString(h, "JY.COM")
	io.WriteString(h, t.String())
	str := fmt.Sprintf("%x", h.Sum(nil))
	str = string([]rune(str)[:n])
	return str
}

// 生成随机数  [min, max)
func RandInt(min, max int) int {
	if min == max { //兼容相等
		return min
	}
	if min > max {
		panic("随机数生成参数错误")
	}
	//rand.Seed(time.Now().UnixNano())
	//返回的数值是 min~max-1 内的随机值
	//return rand.Intn(max-min) + min
	needNum := max - min
	result := rand.Intn(needNum)
	return result + min
}

// 判断数组中是否包含数字参数
func IsContains(nub int64, con []int64) (index int) {
	index = -1
	for i := 0; i < len(con); i++ {
		if con[i] == nub {
			index = i
			return
		}
	}
	return
}

// 判断数组中是否包含字符串参数
func IsContainsStr(str string, con []string) (index int) {
	index = -1
	for i := 0; i < len(con); i++ {
		if con[i] == str {
			index = i
			return
		}
	}
	return
}

//aes  加密
func AesEncrypt(ciphertext, key []byte) string {
	pkey := PaddingLeft(key, '0', aes.BlockSize)
	block, err := aes.NewCipher(pkey) //选择加密算法
	if err != nil {
		return ""
	}
	blockSize := block.BlockSize()

	ciphertext = PKCS7Padding(ciphertext, blockSize)
	blockModel := cipher.NewCBCEncrypter(block, pkey)
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, []byte(ciphertext))
	//plantText = PKCS7UnPadding(plantText)
	//return hex.EncodeToString(plantText)
	return base64.StdEncoding.EncodeToString(plantText)
}

//aes  解密
func AesDecrypt(inData string, key []byte) ([]byte, error) {
	ciphertext, _ := base64.StdEncoding.DecodeString(inData)
	pkey := PaddingLeft(key, '0', aes.BlockSize)
	block, err := aes.NewCipher(pkey) //选择加密算法
	if err != nil {
		return nil, err
	}
	blockModel := cipher.NewCBCDecrypter(block, pkey)
	plantText := make([]byte, len(ciphertext))
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("input not full blocks")
	}
	blockModel.CryptBlocks(plantText, []byte(ciphertext))
	plantText = PKCS7UnPadding(plantText)
	return plantText, nil
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	log.Printf("the len is %d and xlen is %d", length, unpadding)
	return plantText[:(length - unpadding)]
}

func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PaddingLeft(ori []byte, pad byte, length int) []byte {
	if len(ori) >= length {
		return ori[:length]
	}
	pads := bytes.Repeat([]byte{pad}, length-len(ori))
	return append(pads, ori...)
}

// 大厅与子游戏的游戏客户端endpoint 基接口
type IClientEndpoint interface {
	client_server.IServer2Client
	easygo.IEndpointWithSocket
}

//初始化管理员
func InitUsers(users []*share_message.Manager) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_MANAGER)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)

	if count == 0 {
		var il []interface{}
		for _, rr := range users {
			salt := RandString(8)
			rr.Id = easygo.NewInt64(int64(10000 + NextId(TABLE_MANAGER)))
			rr.Password = easygo.NewString(CreatePasswd("123456", salt))
			rr.Salt = easygo.NewString(salt)
			il = append(il, rr)
		}

		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)
	}
}

//截取字符串str的前面n个字符,并判断是有需要在截取的字符串后面添加字符串
func CutOutStr(num int, str string, tail string) string {

	var newStr string
	if num < len(str) {
		newStr = str[0 : num+1]
		if tail != "" {
			newStr = newStr + tail
		}
	} else {
		newStr = str
	}
	return newStr
}

// 再加一个目录
func WriteFile(file string, text string, v ...interface{}) {

	easygo.WriteFile("logs/"+file, text, v...)
}

func GetTimeForString(t int32) string { //时间戳转换为字符串日期
	stime := time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
	return stime
}

//md5加密
func EncryptWithMd5(origData string) string {
	h := md5.New()
	h.Write([]byte(origData))
	digest := h.Sum(nil)
	return hex.EncodeToString(digest)
}

//请求网络资源
func HttpRequest(method, reqUrl string, requestData map[string]string) (string, int) {
	data := url.Values{}
	for k, v := range requestData {
		data.Set(k, v)
	}
	client := &http.Client{}
	request, err := http.NewRequest(method, reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return "NewRequest error", 0
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request)
	if err != nil {
		return "第三方响应失败！", 0
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	// if err1 != nil {
	// }
	return string(body), response.StatusCode

}

//post方式请求网络资源
func HttpPostForm(reqUrl string, requestData map[string]interface{}) ([]byte, int) {
	data := url.Values{}
	for k, v := range requestData {
		switch w := v.(type) {
		case int:
			data.Set(k, strconv.Itoa(w))
		case int32:
			data.Set(k, strconv.Itoa(int(w)))
		case int64:
			data.Set(k, strconv.FormatInt(w, 10))

		case string:
			data.Set(k, w)
		default:
			panic("不支持的类型，请自行添加！")
		}

	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 15 * time.Second}
	response, err := client.PostForm(reqUrl, data)
	if err != nil {
		return []byte(fmt.Sprintf("NewRequest error:%v", err.Error())), 0
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	// if err1 != nil {
	// }
	return body, response.StatusCode

}

//请求网络资源
func HttpRequestWithJson(method, reqUrl string, postData string) (string, int) {
	client := &http.Client{}
	request, err := http.NewRequest(method, reqUrl, strings.NewReader(postData))
	if err != nil {
		return "NewRequest error", 0
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return "第三方响应失败！", 0
	}
	defer response.Body.Close()
	body, err1 := ioutil.ReadAll(response.Body)
	if err1 != nil {
	}
	return string(body), response.StatusCode

}

//传入map，返给类似这样的字符串attach=site_internal&opstate=0&ovalue=50.00
func GetSortStrByMap(keyValues map[string]interface{}, splitStr string, b bool) []byte {
	var keys, signData []string
	for key, value := range keyValues {
		if !b && value == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := keyValues[key]
		signData = append(signData, fmt.Sprintf(`%v=%v`, key, value))

	}
	signStr := strings.Join(signData, splitStr)
	return []byte(signStr)

}

//originalData带签名的内容，pubKey ca证书内容，signData已经签名的密串
func RsaVerySignWithCA(originalData, pubKey, signData []byte) error {
	block, _ := pem.Decode(pubKey)
	ca, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}
	pub := ca.PublicKey
	hash := sha1.New()
	hash.Write(originalData)
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), signData)
}

//通过玩家站点获取配置站点
func GetConfigSite(dbSite string) string {
	//找到站点并且不混服运营的
	if dbSite == MONGODB_NINGMENG || dbSite == "site_robot" || dbSite == "site_peiwan" {
		return MONGODB_NINGMENG
	}
	return dbSite
}

//func GetSortStrBySortMap(keyValues map[string]interface{}, splitStr string) []byte {
//	var keys, signData []string
//	for key := range keyValues {
//		keys = append(keys, key)
//	}
//	sort.Strings(keys)
//	for _, key := range keys {
//		value := keyValues[key]
//		signData = append(signData, fmt.Sprintf(`%v=%v`, key, value))
//
//	}
//	signStr := strings.Join(signData, splitStr)
//	return []byte(signStr)
//
//}
//获取所有随机名字
func GetRobotName(sex int) []string {
	var namelist []string
	if sex == 1 {
		namelist = ManRandName
	} else {
		namelist = WomanRandName
	}
	for i := len(namelist) - 1; i > 0; i-- { //随机排序
		num := rand.Intn(i + 1)
		namelist[i], namelist[num] = namelist[num], namelist[i]
	}
	return namelist
}

//随机一定数量的随机名字
func GetManyRobotName(sex, num int) []string {
	names := GetRobotName(sex)
	return names[:num]
}

//随机一定数量的随机头像
func GetManyRobotHeadIcon(num, sex int) []int {
	var max int
	if sex == 1 {
		max = int(QuerySysParameterById(AVATAR_PARAMETER).GetMavatarCount())
	} else {
		max = int(QuerySysParameterById(AVATAR_PARAMETER).GetWavatarCount())
	}

	headlist := rand.Perm(max)
	return headlist[:num]
}

//随机获取真实头像
func GetRandRealHeadIcon(sex int) string {
	var mark string
	var icon int = 1
	if sex == 1 {
		mark = "mavatar"
	} else {
		mark = "wavatar"
	}
	headList := GetManyRobotHeadIcon(1, sex)
	icon = headList[0]
	head := fmt.Sprintf("https://im-resource-1253887233.file.myqcloud.com/prod/%s/%d.png", mark, icon)
	return head
}

//随机获取群头像
func GetRandTeamHeadIcon(n ...int) string {
	n = append(n, 7)
	icon := util.RandIntn(n[0])
	head := fmt.Sprintf("https://im-resource-1253887233.file.myqcloud.com/defaulticon/%s%d.png", "group_", icon)
	return head
}

//随机获取匹配引导语
func GetRandMatchGuide(guide string) string {
	m := []bson.M{
		{"$match": bson.M{"_id": bson.M{"$ne": guide}}},
		{"$sample": bson.M{"size": 1}},
	}
	one := FindPipeOne(MONGODB_NINGMENG, TABLE_MATCH_GUIDE, m)
	if one == nil {
		return ""
	}
	return one.(bson.M)["_id"].(string)
}

//随机获取SayHi
func GetRandSayHi() string {
	m := []bson.M{
		{"$sample": bson.M{"size": 1}},
	}
	one := FindPipeOne(MONGODB_NINGMENG, TABLE_SAY_HI, m)
	if one == nil {
		return ""
	}
	return one.(bson.M)["_id"].(string)
}

//计算base64图片的流大小
func CountImgBase64Size(imageByte []byte) int {
	imageBaseStr := string(imageByte)
	if strings.Index(imageBaseStr, ",") < 0 {
		return 0
	}
	base64ImgSplice := strings.Split(imageBaseStr, ",")
	imageQz := base64ImgSplice[0]
	if strings.Index(imageQz, "image") < 0 {
		return 0
	}
	base64Str := base64ImgSplice[1]
	eaq := strings.Index(base64Str, "=")
	if eaq < 0 {
		return 0
	}

	newBase64Str := base64Str[0:eaq]
	strLen := len(newBase64Str)
	return int(strLen - (strLen/8)*2)
}

//*****************************开元AES加密解密***********************BEGIN//
func KYAesDecrypt(crypted, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("err is:", err)
	}
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	fmt.Println("decrypt :", string(origData))
	return origData
}

func KYAesEncrypt(src, key string) []byte {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println("key error1", err)
	}
	if src == "" {
		fmt.Println("plain content empty")
	}
	ecb := NewECBEncrypter(block)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	// base64Encode
	//fmt.Println("encrypt:", base64.StdEncoding.EncodeToString(crypted))

	return crypted
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// remove last byte unpadding
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

//*****************************开元AES加密解密***********************END//
//t:请求短信类型
//phone：电话号码，是国际电话一定要带上区号,国内号码不需要
//isInternational 是否国际电话，true |false
func SendCodeToClientUser(t int32, phone string, areaCode string) string {
	var code string
	// 检测手机号是否是000开头的,如果是,生成固定的四位数短信验证码(应该是测试时使用.).
	if phone[0:3] == "000" {
		code = "0258"
	} else {
		if t == SmsTypeBindBankCard {
			code = strconv.Itoa(RandInt(100000, 1000000)) //生成6位随机验证码
		} else {
			code = strconv.Itoa(RandInt(1000, 10000)) //生成4位随机验证码
		}
		// 使用阿里的运营商
		if !NewSMSInst(SMS_BUSINESS_ALI).SendMessageCode(phone, code, areaCode) {
			return ""
		}
	}
	MessageMarkInfo.AddMessageMarkInfo(t, phone, code)
	return code
}

func GetMillSecond() int64 {
	return time.Now().UnixNano() / 1e6
}
func MakeAddress(ip string, port int32) string {
	return ip + ":" + easygo.AnytoA(port)
}

//读取服务器配置
func ReadServerInfoByYaml() *share_message.ServerInfo {
	srvId := easygo.YamlCfg.GetValueAsInt("SERVER_ID")
	srvName := easygo.YamlCfg.GetValueAsString("SERVER_NAME")
	srvType := easygo.YamlCfg.GetValueAsInt("SERVER_TYPE")
	srvExternalIp := easygo.YamlCfg.GetValueAsString("SERVER_ADDR")
	srvInternalIp := easygo.YamlCfg.GetValueAsString("SERVER_ADDR_INTERNAL")
	srvClientWSPort := easygo.YamlCfg.GetValueAsInt("LISTEN_PORT_FOR_CLIENT")
	srvClientTCPPort := easygo.YamlCfg.GetValueAsInt("LISTEN_PORT_FOR_CLIENT_TCP")
	srvClientApiPort := easygo.YamlCfg.GetValueAsInt("LISTEN_PORT_FOR_WEB_API_CLIENT")
	srvServerApiPort := easygo.YamlCfg.GetValueAsInt("LISTEN_PORT_FOR_WEB_API_SERVER")

	srvWebApiPort := easygo.YamlCfg.GetValueAsInt("LISTEN_ADDR_FOR_WEB_API")
	srvBackstagePort := easygo.YamlCfg.GetValueAsInt("LISTEN_PORT_FOR_BACKSTAGE_API")
	srvVersion := easygo.YamlCfg.GetValueAsString("VERSION_NUMBER")
	server := &share_message.ServerInfo{
		Sid:              easygo.NewInt32(srvId),
		Name:             easygo.NewString(srvName),
		Type:             easygo.NewInt32(srvType),
		ExternalIp:       easygo.NewString(srvExternalIp),
		InternalIP:       easygo.NewString(srvInternalIp),
		ClientWSPort:     easygo.NewInt32(srvClientWSPort),
		ClientTCPPort:    easygo.NewInt32(srvClientTCPPort),
		ClientApiPort:    easygo.NewInt32(srvClientApiPort),
		ServerApiPort:    easygo.NewInt32(srvServerApiPort),
		WebApiPort:       easygo.NewInt32(srvWebApiPort),
		BackStageApiPort: easygo.NewInt32(srvBackstagePort),
		Version:          easygo.NewString(srvVersion),
	}
	return server
}

type AuthPeopleIdResult struct {
	Realname string `json:"realname"`
	Idcard   string `json:"idcard"`
	Res      int    `json:"res"`
}

type AuthPeopleId struct {
	Error_code int                `json:"error_code"`
	Reason     string             `json:"reason"`
	Result     AuthPeopleIdResult `json:"result"`
}

//身份认证
func AuthPeopleIdName(id, name string) bool {
	k := easygo.YamlCfg.GetValueAsString("JH_AUTHPEOPLEIDKEY")
	params := url.Values{"idcard": {id}, "realname": {name}, "key": {k}}
	resp, err := http.PostForm("http://op.juhe.cn/idcard/query", params)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var obj AuthPeopleId
	rep, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		easygo.PanicError(err1)
	}
	json.Unmarshal([]byte(string(rep)), &obj)
	if obj.Error_code != 0 {
		return false
	}
	if obj.Result.Res != 1 {
		return false
	}

	return true
}

type AuthBankIdResult struct {
	Jobid    string `json:"jobid"`
	Bankcard string `json:"bankcard"`
	Realname string `json:"realname"`
	Idcard   string `json:"idcard"`
	Res      int    `json:"res"`
	Message  string `json:"message"`
	Mobile   string `json:"mobile"`
}

type AuthBankId struct {
	Error_code int              `json:"error_code"`
	Reason     string           `json:"reason"`
	Result     AuthBankIdResult `json:"result"`
}

// AuthBankIdName 验证绑卡人银行卡信息是否有误
func AuthBankIdName(bankId, name, id, phone string) (bool, string) {
	k := easygo.YamlCfg.GetValueAsString("JH_AUTHBANKIDKEY")
	params := url.Values{"idcard": {id}, "realname": {name}, "key": {k}, "bankcard": {bankId}, "mobile": {phone}}
	resp, err := http.PostForm("http://v.juhe.cn/verifybankcard4/query", params)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var obj AuthBankId
	rep, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		easygo.PanicError(err1)
	}
	json.Unmarshal([]byte(string(rep)), &obj)
	logs.Error("error:", obj)
	if obj.Error_code != 0 {
		return false, obj.Result.Message
	}
	if obj.Result.Message != "验证成功" {
		return false, obj.Result.Message
	}
	return true, obj.Result.Message
}

type BankInfo struct {
	Bank      string `json:"bank"`
	Vaildated bool   `json:"vaildated"`
	CardType  string `json:"card_type"`
	Key       string `json:"key"`
	Message   []int  `json:"message"`
	Stat      string `json:"stat"`
}

func GetBankCodeForBankId(id string) string {
	s := fmt.Sprintf("https://ccdcapi.alipay.com/validateAndCacheCardInfo.json?_input_charset=utf-8&cardNo=%s&cardBinCheck=true", id)
	resp, err := http.Get(s)
	if err != nil {
		panic(err)
	}
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		panic(err1)
	}
	var obj BankInfo
	json.Unmarshal([]byte(string(body)), &obj)
	if obj.Stat != "ok" || obj.Bank == "" {
		return ""
	}
	return obj.Bank
}

//返回单位为：千米
func GetDistance(lat1, lat2, lng1, lng2 float64) float64 {
	radius := 6371000.0 //6378137.0
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius
}

//字符串string转[]byte数组
func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//数组转字符串
func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func GetGoldChangeNote(t int32, pid PLAYER_ID, kwarg map[string]interface{}) string {
	var reason string
	switch t {
	case GOLD_TYPE_CASH_AFIN:
		tt := kwarg["Type"].(int32)
		if tt == PAY_TYPE_BACKSTAGE_IN {
			reason = "人工入款"
		} else if tt == PAY_TYPE_BACKSTAGE_OUT {
			reason = "人工出款"
		}
	case GOLD_TYPE_CASH_IN:
		t := kwarg["Type"].(int32)
		if t == 1 {
			reason = "零钱充值-来自微信支付"
		} else if t == 2 {
			reason = "零钱充值-来自支付宝支付"
		} else if t == 3 {
			code := kwarg["BankId"].(string)
			name := kwarg["BankName"].(string)
			reason = fmt.Sprintf("零钱充值-来自%s(%s)", name, code)
		}

	case GOLD_TYPE_CASH_OUT:
		code := kwarg["BankId"].(string)
		name := kwarg["BankName"].(string)
		reason = fmt.Sprintf("零钱提现-到%s(%s)", name, code)
	case GOLD_TYPE_GET_REDPACKET:
		base := GetRedisPlayerBase(pid)
		reason = fmt.Sprintf("柠檬红包-来自%s", base.GetNickName())
	case GOLD_TYPE_CASH_OUT_BACK:
		reason = "退还零钱-提现被取消"
	case GOLD_TYPE_SEND_REDPACKET:
		redtype := kwarg["RedType"].(int32)
		if redtype == 1 {
			base := GetRedisPlayerBase(pid)
			reason = fmt.Sprintf("柠檬红包-发给%s", base.GetNickName())
		} else {
			reason = "柠檬红包-发出群红包"
		}
	case GOLD_TYPE_SHOP_MONEY:
		reason = fmt.Sprintf("好物-付给%s", kwarg["SponsorName"])
	case GOLD_TYPE_BACK_MONEY:
		reason = fmt.Sprintf("商城订单%s退款", easygo.AnytoA(kwarg["ShopOrderID"]))
	case GOLD_TYPE_REDPACKET_OVERTIME:
		redtype := kwarg["RedType"].(int32)
		if redtype == 1 {
			base := GetRedisPlayerBase(pid)
			reason = fmt.Sprintf("柠檬红包-发给%s", base.GetNickName())
		} else {
			reason = "柠檬红包-发出群红包"
		}
	case GOLD_TYPE_TRANSFER_MONEY_OVER:
		base := GetRedisPlayerBase(pid)
		reason = fmt.Sprintf("转账退还-转给%s", base.GetNickName())
	case GOLD_TYPE_GET_MONEY:
		base := GetRedisPlayerBase(pid)
		reason = fmt.Sprintf("二维码收款-来自%s", base.GetNickName())
	case GOLD_TYPE_PAY_MONEY:
		base := GetRedisPlayerBase(pid)
		reason = fmt.Sprintf("扫二维码付款-给%s", base.GetNickName())
	case GOLD_TYPE_GET_TRANSFER_MONEY:
		base := GetRedisPlayerBase(pid)
		reason = fmt.Sprintf("转账-来自%s", base.GetNickName())
	case GOLD_TYPE_SEND_TRANSFER_MONEY:
		base := GetRedisPlayerBase(pid)
		reason = fmt.Sprintf("转账-转给%s", base.GetNickName())
	case GOLD_TYPE_FINE_MONEY:
	case GOLD_TYPE_EXTRA_MONEY:
	case GOLD_TYPE_SHOP_ITEM_MONEY:
		reason = fmt.Sprintf("好物-来自%s", kwarg["ReceiverName"])
	case GOLD_TYPE_EXCHANGE_COIN:
		reason = "兑换硬币"
	case COIN_TYPE_SHOP_OUT:
		reason = fmt.Sprintf("商场消费-%s*%s", kwarg["Name"], kwarg["Num"])
	}
	return reason
}

//非充值提现类下订单  返回订单号
//func PlaceOrder(playerId PLAYER_ID, changeGold int64, sourceType int32, Id ...string) (string, *base.Fail) {
//	bankId := append(Id, "")[0]
//	player := GetRedisPlayerBase(playerId)
//	if player == nil {
//		s := fmt.Sprintf("玩家Id: %d 不存在", playerId)
//		return "", easygo.NewFailMsg(s, FAIL_MSG_CODE_1005)
//	}
//
//	if player.GetGold()+changeGold < 0 {
//		s := fmt.Sprintf("金额不足")
//		return "", easygo.NewFailMsg(s, FAIL_MSG_CODE_1006)
//	}
//
//	var reason string
//	switch sourceType {
//	case GOLD_TYPE_CASH_IN:
//		reason = "充值成功"
//	case GOLD_TYPE_GET_REDPACKET:
//		reason = "收红包"
//	case GOLD_TYPE_GET_TRANSFER_MONEY:
//		reason = "转入成功"
//	case GOLD_TYPE_GET_MONEY:
//		reason = "收款成功"
//	case GOLD_TYPE_REDPACKET_OVERTIME:
//		reason = "红包退款"
//	case GOLD_TYPE_TRANSFER_MONEY_OVER:
//		reason = "转账退款"
//	case GOLD_TYPE_BACK_MONEY:
//		reason = "商家退款"
//	case GOLD_TYPE_SHOP_ITEM_MONEY:
//		reason = "商城卖家货款"
//	case GOLD_TYPE_CASH_OUT:
//		reason = "提现"
//	case GOLD_TYPE_SEND_REDPACKET:
//		reason = "发红包"
//	case GOLD_TYPE_SEND_TRANSFER_MONEY:
//		reason = "转出成功"
//	case GOLD_TYPE_PAY_MONEY:
//		reason = "付款成功"
//	case GOLD_TYPE_FINE_MONEY:
//		reason = "罚没"
//	case GOLD_TYPE_EXTRA_MONEY:
//		reason = "手续费"
//	case GOLD_TYPE_SHOP_MONEY:
//		reason = "商家消费"
//	}
//
//	changeType := 1
//	if changeGold < 0 {
//		changeType = 2
//	}
//
//	tax := int64(0)
//	if sourceType == GOLD_TYPE_CASH_OUT {
//		tax = int64(float64(changeGold) * 0.008)
//	}
//
//	order := &share_message.Order{
//		PlayerId:    easygo.NewInt64(player.GetPlayerId()),
//		Account:     easygo.NewString(player.GetAccount()),
//		NickName:    easygo.NewString(player.GetNickName()),
//		RealName:    easygo.NewString(player.GetRealName()),
//		SourceType:  easygo.NewInt32(sourceType),
//		ChangeType:  easygo.NewInt32(changeType),
//		Channeltype: easygo.NewInt32(0),
//		CurGold:     easygo.NewInt64(player.GetGold()),
//		ChangeGold:  easygo.NewInt64(changeGold),
//		Gold:        easygo.NewInt64(player.GetGold() + changeGold),
//		Amount:      easygo.NewInt64(changeGold),
//		CreateTime:  easygo.NewInt64(GetMillSecond()),
//		CreateIP:    easygo.NewString("127.0.0.1"),
//		Status:      easygo.NewInt32(0),
//		PayStatus:   easygo.NewInt32(1),
//		Note:        easygo.NewString(reason),
//		Tax:         easygo.NewInt64(-tax),
//		Operator:    easygo.NewString("system"),
//		BankInfo:    easygo.NewString(bankId),
//	}
//	orderId := RedisCreateOrder(order, true)
//	return orderId, easygo.NewFailMsg("", FAIL_MSG_CODE_SUCCESS)
//}

type JGLoginInfo struct {
	Id      int64  `json:"id"`
	Exid    string `json:"exID"`
	Code    int32  `json:"code"`
	Content string `json:"content"`
	Phone   string `json:"phone"`
}

func GetJGOneKeyLoginPhone(token string, apkCode int32) string {
	var phone, key, value string
	appKey := MakeNewString("JG_APPKEY", 100)
	secretkey := MakeNewString("JG_SECRET", 100)
	if apkCode != 0 {
		appKey = MakeNewString("JG_APPKEY", apkCode)
		secretkey = MakeNewString("JG_SECRET", apkCode)
	}
	key = easygo.YamlCfg.GetValueAsString(appKey)
	value = easygo.YamlCfg.GetValueAsString(secretkey)
	data := []byte(fmt.Sprintf(`{"loginToken":"%s"}`, token))
	buffer := bytes.NewBuffer(data)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.verification.jpush.cn/v1/web/loginTokenVerify", buffer)
	easygo.PanicError(err)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(key, value)
	resp, err := client.Do(req)
	easygo.PanicError(err)
	var obj JGLoginInfo
	rep, err1 := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err1)
	json.Unmarshal([]byte(string(rep)), &obj)
	fmt.Println(obj.Code, obj.Content, obj.Id, obj.Exid)
	if obj.Code != 8000 {
		return ""
	}
	if obj.Phone == "" {
		return ""
	}
	secret, _ := base64.StdEncoding.DecodeString(obj.Phone)
	b, err2 := RsaDecrypt(secret, JGPrivateKey) //解密
	easygo.PanicError(err2)
	phone = string(b)
	logs.Info("============GetJGOneKeyLoginPhone", phone)
	return phone
}

// rsa加密
func RsaEncrypt(origData []byte, pubKey string) ([]byte, error) {
	//解密pem格式的公钥
	publicKey := []byte(pubKey)
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(crand.Reader, pub, origData)
}

// pkcs8解密
func RsaDecrypt(secret []byte, priKey string) ([]byte, error) {
	privateKey := []byte(priKey)
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(crand.Reader, priv.(*rsa.PrivateKey), secret)
}

func GetWeChatInfo(code string, apkCode int32) (string, string, string) {
	appKey := MakeNewString("WEIXIN_APPID", apkCode)
	secretKey := MakeNewString("WEIXIN_APPSECRET", apkCode)
	appId := easygo.YamlCfg.GetValueAsString(appKey)
	appSecret := easygo.YamlCfg.GetValueAsString(secretKey)
	wechatToken, openId, unionId := GetWechatAuth(appId, appSecret, code)
	return wechatToken, openId, unionId
}

func GetWechatAuth(appId, appSecret, token string) (string, string, string) {
	url := easygo.YamlCfg.GetValueAsString("WEIXIN_URL")
	resp, err := http.Get(fmt.Sprintf(url, appId, appSecret, token))
	easygo.PanicError(err)
	defer resp.Body.Close()
	data, err1 := ioutil.ReadAll(resp.Body) //把  body 内容读入字符串 s
	easygo.PanicError(err1)
	type WeiXin struct {
		OpenId       string `json:"openid"`
		ErrCode      int32  `json:"errcode"` // 错误码
		Access_token string `json:"access_token"`
		Errmsg       string `json:"errmsg"`
		UnionId      string `json:"unionid"` // 微信唯一id
	}
	var weixininfo WeiXin
	err2 := json.Unmarshal(data, &weixininfo)
	easygo.PanicError(err2)
	errCode := weixininfo.ErrCode
	if errCode != 0 {
		logs.Info(errCode, weixininfo.Errmsg)
		return "", "", ""
	}
	return weixininfo.Access_token, weixininfo.OpenId, weixininfo.UnionId
}

func GetWeChatUserInfo(token, openId string) (string, int32, string) {
	url := easygo.YamlCfg.GetValueAsString("WEIXIN_USERURL")
	resp, err := http.Get(fmt.Sprintf(url, token, openId))
	easygo.PanicError(err)
	defer resp.Body.Close()
	data, err1 := ioutil.ReadAll(resp.Body) //把  body 内容读入字符串 s
	easygo.PanicError(err1)
	type UserInfo struct {
		Nickname   string `json:"nickname"`
		Sex        int32  `json:"sex"`
		Headimgurl string `json:"headimgurl"`
		Errcode    int32  `json:"errcode"`
		Errmsg     string `json:"errmsg"`
		UnionId    string `json:"unionid"`
	}
	var userInfo UserInfo
	err2 := json.Unmarshal(data, &userInfo)
	easygo.PanicError(err2)
	if userInfo.Errcode != 0 {
		logs.Info(userInfo.Errcode, userInfo.Errmsg)
		return "", 0, ""
	}
	return userInfo.UnionId, userInfo.Sex, userInfo.Headimgurl
}

//获取mongo个人表名
func GetMongoTableName(key interface{}, tbName string) string {
	if key == nil {
		panic("如果没有key请不要调用 GetMongoTableName 方法")
	}
	return tbName + "_" + easygo.AnytoA(key)
}

type PushMessage struct {
	Title       string //标题
	Content     string //内容
	ContentType string //私聊还是群聊
	TargetId    string //目标id
	ChatType    string //聊天类型
	Msg         string //结构体json字符串
	JumpObject  int32  //跳转对象 1 主界面，2 柠檬团队，3 柠檬助手,10 -硬币页面,11-我的物品页
	JumpUrl     string //跳转URL
	JumpType    int32  //跳转类型 1外部跳转，2内部跳转
	Icon        string //文章封面
	ArticleType int32  //跳转类型
	Location    int32  //跳转位置：1 主界面，2 柠檬助手，3 柠檬团队，0附近的人，5社交广场-主界面，6社交广场-新增关注，7社交广场-指定动态：通过填写动态ID指定，8好物-主界面，9好物-指定商品：通过填写商品ID指定,,10群-指定群id,11社交广场发布页,12零钱,13话题-指定话题,14-指定的动态评论,15-话题主界面
	OrderId     string //订单id
	OperaType   int32  //0买家  2卖家
	ObjectId    int64  //对象
	PlayerId    int64  //玩家id
	ArticleId   int64  //文章id
	CommentId   int64  // 评论的id
	ItemId      int32  //推送项目id: PUSH_ITEM_101,PUSH_ITEM_102
}

const (
	JG_TYPE_PERSONALCHAT     = "personal_chat"       //私聊
	JG_TYPE_TEAMCHAT         = "team_chat"           //群聊
	JG_TYPE_SHOP             = "shop_notice"         //商城
	JG_TYPE_SHOP_MESSAGE     = "shop_message"        //商城留言
	JG_TYPE_BACKSTAGE        = "backstage_notice"    //后台
	JG_TYPE_BACKSTAGE_ASS    = "backstage_assistant" //后台小助手
	JG_TYPE_HALL             = "hall_notice"         // 大厅
	JG_TYPE_SQUARE           = "square_notice"       // 广场
	JG_TYPE_BACKSTAGE_ESPORT = "backstage_esport"    //后台电竞推送
)

//func JGSendMessage(ids []string, info PushMessage) {
//	//Platform
//	var pf jpushclient.Platform
//	pf.Add(jpushclient.ANDROID)
//	pf.Add(jpushclient.IOS)
//	pf.Add(jpushclient.WINPHONE)
//	//pf.All()
//
//	//Audience
//	var ad jpushclient.Audience
//	//var msg jpushclient.Message
//	//ad.SetTag(s)
//	//ad.SetID(s)
//	//ad.All()
//	t := info.ContentType
//	m := make(map[string]interface{})
//	if t == JG_TYPE_PERSONALCHAT || t == JG_TYPE_TEAMCHAT {
//		m["ContentType"] = t
//		m["TargetId"] = info.TargetId
//		m["ChatType"] = info.ChatType
//		if info.Msg != "" {
//			m["Msg"] = info.Msg
//		}
//		ad.SetAlias(ids)
//	} else if t == JG_TYPE_BACKSTAGE {
//		ad.SetAlias(ids)
//		m["JumpObject"] = info.JumpObject
//		m["Title"] = info.Title     //标题
//		m["Content"] = info.Content //文章概要
//		m["ContentType"] = t
//	} else if t == JG_TYPE_SHOP {
//		ad.SetAlias(ids)
//		m["Title"] = info.Title     //标题
//		m["Content"] = info.Content //文章概要
//		m["ContentType"] = t
//	} else if t == JG_TYPE_BACKSTAGE_ASS {
//		ad.SetAlias(ids)
//		m["ContentType"] = t
//		m["TargetId"] = -4                  //目标ID
//		m["JumpUrl"] = info.JumpUrl         //文章地址
//		m["Icon"] = info.Icon               //文章图标
//		m["ArticleType"] = info.ArticleType //跳转类型
//		m["Location"] = info.Location       //跳转位置
//		m["Title"] = info.Title             //标题
//		m["Content"] = info.Content         //文章概要
//		logs.Info("进入")
//	}
//	//Notice
//	var notice jpushclient.Notice
//	notice.SetAlert("柠檬畅聊")
//	notice.SetAndroidNotice(&jpushclient.AndroidNotice{Alert: info.Content, Title: info.Title, Extras: m, Uri_activity: "com.silbermond.tktalk.NotificationActivity", Uri_action: "com.silbermond.tktalk.NotificationActivity.oppo"})
//	notice.SetIOSNotice(&jpushclient.IOSNotice{Alert: info.Content, Sound: "sound.caf", Extras: m, ContentAvailable: true, Badge: 1})
//	notice.SetWinPhoneNotice(&jpushclient.WinPhoneNotice{Alert: "WinPhoneNotice"})
//
//	var option jpushclient.Option
//	option.SetApns(IS_FORMAL_SERVER)
//
//	payload := jpushclient.NewPushPayLoad()
//	payload.SetPlatform(&pf)
//	payload.SetAudience(&ad)
//	//payload.SetMessage(&msg)
//	payload.SetNotice(&notice)
//	payload.SetOptions(&option)
//
//	bytes, _ := payload.ToBytes()
//
//	appKey := easygo.YamlCfg.GetValueAsString("JG_APPKEY")
//	secret := easygo.YamlCfg.GetValueAsString("JG_SECRET")
//	//push
//	c := jpushclient.NewPushClient(secret, appKey)
//	c.Send(bytes)
//}

//	@description   mob推送(v3)
func JGSendMessage(ids []string, info PushMessage, sysp ...*share_message.SysParameter) {
	//logs.Info("推送入口------------->")
	//	logs.Info("ids:", ids)
	if len(sysp) > 0 {
		pushSet := sysp[0].PushSet
		for _, ps := range pushSet {
			if ps.GetObjId() == info.ItemId && !ps.GetIsPush() {
				return
			}
		}
	}

	if len(ids) == 0 {
		logs.Info("发送别名为空:", ids, info)
		return
	}
	key := easygo.YamlCfg.GetValueAsString("MOB_APPKEY")
	secret := easygo.YamlCfg.GetValueAsString("MOB_SECRET")
	url := easygo.YamlCfg.GetValueAsString("MOB_URL")

	var push jpushclient.Push
	push.SetSource("webapi")
	push.SetAppkey(key)

	var pushNotify jpushclient.PushNotify
	pushNotify.SetPlats([]int32{1, 2})
	pushNotify.SetContent(info.Content)
	pushNotify.SetTitle(info.Title)
	pushNotify.SetType(1)
	//正式服
	if IS_FORMAL_SERVER {
		pushNotify.SetIosProduction(1)
	} else {
		pushNotify.SetIosProduction(0)
	}

	var iosNotify jpushclient.IosNotify
	iosNotify.SetBadge(1)
	iosNotify.SetBadgeType(2)

	pushNotify.SetIosNotify(&iosNotify)
	extras := make([]map[string]interface{}, 0)
	t := info.ContentType
	if t == JG_TYPE_PERSONALCHAT || t == JG_TYPE_TEAMCHAT {
		extras = append(extras, map[string]interface{}{
			"key":   "ContentType",
			"value": info.ContentType,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "TargetId",
			"value": info.TargetId,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "ChatType",
			"value": info.ChatType,
		})
		if info.Msg != "" {
			extras = append(extras, map[string]interface{}{
				"key":   "Msg",
				"value": info.Msg,
			})
		}
	} else if t == JG_TYPE_BACKSTAGE {
		extras = append(extras, map[string]interface{}{
			"key":   "JumpObject",
			"value": info.JumpObject,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "Title",
			"value": info.Title,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "Content",
			"value": info.Content,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "ContentType",
			"value": info.ContentType,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "ObjectId",
			"value": info.ObjectId,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "TargetId",
			"value": info.TargetId,
		})
	} else if t == JG_TYPE_SHOP {
		extras = append(extras, map[string]interface{}{
			"key":   "Title",
			"value": info.Title,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "Content",
			"value": info.Content,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "ContentType",
			"value": t,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "OrderId",
			"value": info.OrderId,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "OperaType",
			"value": info.OperaType,
		})
	} else if t == JG_TYPE_BACKSTAGE_ASS {
		extras = append(extras, map[string]interface{}{
			"key":   "ContentType",
			"value": t,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "TargetId",
			"value": -4,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "JumpUrl",
			"value": info.JumpUrl,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "Icon",
			"value": info.Icon,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "ArticleType",
			"value": info.ArticleType,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "Location",
			"value": info.Location,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "Title",
			"value": info.Title,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "Content",
			"value": info.Content,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "ObjectId",
			"value": info.ObjectId,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "ArticleId",
			"value": info.ArticleId,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "PlayerId",
			"value": info.PlayerId,
		})
	} else if t == JG_TYPE_HALL {
		extras = append(extras, map[string]interface{}{
			"key":   "ContentType",
			"value": t,
		})

		extras = append(extras, map[string]interface{}{
			"key":   "JumpObject",
			"value": info.JumpObject,
		})
	} else if t == JG_TYPE_SQUARE {
		extras = append(extras, map[string]interface{}{
			"key":   "ContentType",
			"value": t,
		})

		extras = append(extras, map[string]interface{}{
			"key":   "JumpObject",
			"value": info.JumpObject,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "Title",
			"value": info.Title,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "Content",
			"value": info.Content,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "ObjectId",
			"value": info.ObjectId,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "Location",
			"value": info.Location,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "CommentId",
			"value": info.CommentId,
		})
	} else if t == JG_TYPE_BACKSTAGE_ESPORT {
		extras = append(extras, map[string]interface{}{
			"key":   "ContentType",
			"value": t,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "JumpUrl",
			"value": info.JumpUrl,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "Icon",
			"value": info.Icon,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "ArticleType",
			"value": info.ArticleType,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "Location",
			"value": info.Location,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "Title",
			"value": info.Title,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "Content",
			"value": info.Content,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "ObjectId",
			"value": info.ObjectId,
		})
		extras = append(extras, map[string]interface{}{
			"key":   "JumpType",
			"value": info.JumpType,
		})
	}

	//extras = append(extras, m)
	pushNotify.SetExtrasMapList(extras)
	push.SetPushNotify(&pushNotify)
	//设置推送用户
	var target jpushclient.PushTarget
	target.SetTarget(2)

	lenData := len(ids)
	n := 1
	if lenData > 1000 {
		n = lenData / 1000
		if lenData%1000 != 0 {
			n = n + 1
		}
	}
	//数据量过大，分批处理,每次处理1000条
	for i := 0; i < n; i++ {
		start := i * 1000
		end := (1 + i) * 1000
		if i == n-1 {
			end = lenData
		}
		cls := ids[start:end]
		target.SetAlias(cls)
		//logs.Info("发送别名:", target)
		//logs.Info("发送别名:", target)
		push.SetPushTarget(&target)
		sign := GenerateSign(&push, secret)
		ret := new(jpushclient.Login)
		//logs.Info("准备发送推送 HttpPostBody,----------->url: %s,sign: %s,key: %s", url, sign, key)
		postBody, err := HttpPostBody(url, &push, sign, key)
		if err != nil {
			logs.Error("推送 HttpPostBody 失败,err: ----> %s", err.Error())
		}
		json.Unmarshal(postBody, &ret)
		value, _ := json.Marshal(push)
		if ret.Status == 200 {
			//logs.Info("mob推送成功：", ret)
			WriteFile("mob.log", "mob推送成功:", string(value), target)
		} else {
			logs.Info("推送失败：", ret, string(value))
			WriteFile("mob.log", "mob推送失败:", ret.Error)
		}
	}

}

//	@description   mob推送(v2)
func MobSendMessageV2(ids []string, info PushMessage) {

	key := easygo.YamlCfg.GetValueAsString("MOB_APPKEY")
	secret := easygo.YamlCfg.GetValueAsString("MOB_SECRET")
	url := easygo.YamlCfg.GetValueAsString("MOB_URL_V2")

	var push jpushclient.PushV2
	push.SetAppkey(key)
	push.SetPlats([]int32{1, 2})
	push.SetTarget(2)
	push.SetContent(info.Content)
	push.SetType(1)
	push.SetAlias(ids)
	push.SetAndroidTitle(info.Title)
	push.SetIosTitle(info.Title)
	push.SetIosBadge(1)

	t := info.ContentType
	m := make(map[string]interface{})
	if t == JG_TYPE_PERSONALCHAT || t == JG_TYPE_TEAMCHAT {
		m["ContentType"] = t
		m["TargetId"] = info.TargetId
		m["ChatType"] = info.ChatType
		if info.Msg != "" {
			m["Msg"] = info.Msg
		}
	} else if t == JG_TYPE_BACKSTAGE {
		m["JumpObject"] = info.JumpObject
		m["Title"] = info.Title     //标题
		m["Content"] = info.Content //文章概要
		m["ContentType"] = t
		m["ObjectId"] = info.ObjectId
	} else if t == JG_TYPE_SHOP {
		m["Title"] = info.Title     //标题
		m["Content"] = info.Content //文章概要
		m["ContentType"] = t
	} else if t == JG_TYPE_BACKSTAGE_ASS {
		m["ContentType"] = t
		m["TargetId"] = -4                  //目标ID
		m["JumpUrl"] = info.JumpUrl         //文章地址
		m["Icon"] = info.Icon               //文章图标
		m["ArticleType"] = info.ArticleType //跳转类型
		m["Location"] = info.Location       //跳转位置
		m["Title"] = info.Title             //标题
		m["Content"] = info.Content         //文章概要
		m["ObjectId"] = info.ObjectId
		logs.Info("进入")
	}

	push.SetExtras(m)
	sign := GenerateSignV2(&push, secret)

	ret := new(jpushclient.Login)
	postBody, _ := HttpPostBodyV2(url, &push, sign, key)
	json.Unmarshal(postBody, &ret)
	//if ret.Status == 200{
	//	logs.Info("mob推送成功")
	//}else{
	//	logs.Info(ret.Error)
	//}
}

//获取当前时间字符串:YYYYMMDDHHmmSS
func GetCurTimeString() string {
	t := time.Now()
	s := t.Format("2006-01-02 15:04:05")
	s = strings.Replace(s, "-", "", -1)
	s = strings.Replace(s, ":", "", -1)
	s = strings.Replace(s, " ", "", -1)
	return s
}

//组装字符串
func MakeNewString(keys ...interface{}) string {
	res := ""
	for k, v := range keys {
		res += easygo.AnytoA(v)
		if k < (len(keys) - 1) {
			res += "_"
		}
	}
	return res
}

//组装redis Key
func MakeRedisKey(keys ...interface{}) string {
	res := ""
	for k, v := range keys {
		res += easygo.AnytoA(v)
		if k < (len(keys) - 1) {
			res += ":"
		}
	}
	return res
}

//两个相似字段的结构体，相同字段值数据互转
func StructToOtherStruct(src interface{}, dest interface{}) {
	js, err := json.Marshal(src)
	easygo.PanicError(err)
	err = json.Unmarshal(js, dest)
	easygo.PanicError(err)
}
func StructToMap(src interface{}, dest interface{}) {
	js, err := json.Marshal(src)
	easygo.PanicError(err)
	err = json.Unmarshal(js, dest)
	easygo.PanicError(err)
}

//[]uint8数组转int64
func InterfersToInt64s(src []interface{}, dest *[]int64) {
	for i := range src {
		v := string(src[i].([]uint8))
		*dest = append(*dest, easygo.AtoInt64(v))
	}
	return
}

func InterfersToInt64(src []interface{}) []int64 {
	var dest []int64
	for _, val := range src {
		if v, ok := val.(int64); ok {
			dest = append(dest, v)
		}
	}
	return dest
}

//[]uint8数组转int32
func InterfersToInt32s(src []interface{}, dest *[]int32) {
	for i := range src {
		v := string(src[i].([]uint8))
		*dest = append(*dest, easygo.AtoInt32(v))
	}
	return
}

//[]uint8数组转[]string
func InterfersToStrings(src []interface{}, dest *[]string) {
	for i := range src {
		v := string(src[i].([]uint8))
		*dest = append(*dest, v)
	}
	return
}

func Int64StringMap(result interface{}, err error) (map[int64]string, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: StringMap expects even number of values result")
	}
	m := make(map[int64]string, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		value, okValue := values[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redigo: StringMap key not a bulk string value")
		}
		k := string(key)
		k1 := easygo.AtoInt64(k)
		m[k1] = string(value)
	}
	return m, nil
}
func ObjListExistStrKey(result interface{}, err error, ikey string) (bool, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return false, err
	}
	if len(values)%2 != 0 {
		return false, errors.New("redigo: StringMap expects even number of values result")
	}

	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		_, okValue := values[i+1].([]byte)
		if !okKey || !okValue {
			return false, errors.New("redigo: StringMap key not a bulk string value")
		}
		k := string(key)
		if k == ikey {
			return true, err
		}
	}
	return false, err
}
func ObjListToStrKeyList(result interface{}, err error) ([]string, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: StringMap expects even number of values result")
	}
	lst := []string{}
	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		_, okValue := values[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redigo: StringMap key not a bulk string value")
		}
		k := string(key)
		lst = append(lst, k)
	}
	return lst, err
}

func StrkeyStringMap(result interface{}, err error) (map[string]string, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: StringMap expects even number of values result")
	}
	m := make(map[string]string, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		value, okValue := values[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redigo: StringMap key not a bulk string value")
		}
		k := string(key)
		m[k] = string(value)
	}
	return m, nil
}

//redis obj 转int64 数组
func ObjInt64List(result interface{}, err error) ([]int64, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: StringMap expects even number of values result")
	}
	list := []int64{}
	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		_, okValue := values[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redigo: StringMap key not a bulk string value")
		}
		k := string(key)
		k1 := easygo.AtoInt64(k)
		//m[k1] = string(value)
		list = append(list, k1)
	}
	return list, nil
}

func GetCurTimeString2() string {
	t := time.Now()
	s := t.Format("2006-01-02 15:04:05")
	return s
}

// GetCurTimeString3 获取时间格式为 20060102150405 字符串
func GetCurTimeString3() string {
	t := time.Now()
	s := t.Format("20060102150405")
	return s
}

//批量存储数据:
func UpsertAll(mongo easygo.IMongoDBManager, db, tab string, data []interface{}) {
	//reids存储数据写入磁盘，恢复数据用
	WriteFile("redis_log_"+util.GetYMD(), tab+":", data)
	col, closeFun := mongo.GetC(db, tab)
	defer closeFun()
	lenData := len(data)
	n := lenData / 1000 //取整
	r := lenData % 1000 //求余
	if r == 0 {
		n = n - 1
	}
	//数据量过大，分批处理,每次处理1000条
	for i := 0; i <= n; i++ {
		start := i * 1000
		end := (1 + i) * 1000
		if i == n {
			end = len(data)
		}
		bulk := col.Bulk()
		saveData := data[start:end]
		lenSaveData := len(saveData)
		if lenSaveData%2 != 0 {
			logs.Info("dataIndex:", start, end)
			panic("mongo存储数据长度异常:" + easygo.AnytoA(lenSaveData))
		}
		bulk.Upsert(saveData...)
		_, err := bulk.Run()
		if err != nil {
			easygo.PanicError(err)
		}
	}
}

// @description	MObTech加密方法V3
// @param	request	map	"请求体"
// @param	secret	string	"密钥"
func GenerateSign(push *jpushclient.Push, secret string) string {
	ret := ""
	b, _ := json.Marshal(push)
	ret = string(b) + secret
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(ret))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// @description	MObTech加密方法V2
// @param	request	map	"请求体"
// @param	secret	string	"密钥"
func GenerateSignV2(push *jpushclient.PushV2, secret string) string {
	ret := ""
	b, _ := json.Marshal(push)
	ret = string(b) + secret
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(ret))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// @description	MobTech请求函数
// @param	msg	[]byte	"请求体"
// @param	sign	string	"签名"
// @param	key	string	"密钥"
func HttpPostBody(url string, push *jpushclient.Push, sign string, key string) ([]byte, error) {

	msg, _ := json.Marshal(push)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(msg))

	req.Header.Set("key", key)
	req.Header.Set("sign", sign)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return []byte(""), err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// @description	MobTech请求函数
// @param	msg	[]byte	"请求体"
// @param	sign	string	"签名"
// @param	key	string	"密钥"
func HttpPostBodyV2(url string, push *jpushclient.PushV2, sign string, key string) ([]byte, error) {

	msg, _ := json.Marshal(push)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(msg))

	req.Header.Set("key", key)
	req.Header.Set("sign", sign)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return []byte(""), err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

//通用的base64解码字符串
func Base64DecodeStr(str string) string {
	newStr := str
	b, err := base64.StdEncoding.DecodeString(str)
	if err == nil {
		newStr = string(b)
	}
	return newStr
}

//
////通用base64转码
//func Base64EncodeStr(str string) string {
//	_, err := base64.StdEncoding.DecodeString(str)
//	if err != nil {
//		newStr := base64.StdEncoding.EncodeToString([]byte(str))
//		return newStr
//	}
//	return str
//}
//随机获取默认头像 @sex 性别 1男，2女
func GetDefaultHeadicon(sex int) string {
	icon := rand.Intn(5) + 1
	headIconUrl := ""
	switch sex {
	case 2:
		headIconUrl = fmt.Sprintf("https://im-resource-1253887233.cos.accelerate.myqcloud.com/defaulticon/girl_%d.png", icon)
	case 1:
		headIconUrl = fmt.Sprintf("https://im-resource-1253887233.cos.accelerate.myqcloud.com/defaulticon/boy_%d.png", icon)
	default:
		headIconUrl = fmt.Sprintf("https://im-resource-1253887233.cos.accelerate.myqcloud.com/defaulticon/girl_%d.png", icon)
	}
	return headIconUrl
}

//修改转发服务器列表
func UpdateTfserver(data *TFserver, jsonFile string) {
	out, _ := json.MarshalIndent(data, "", "  ")
	err := ioutil.WriteFile(jsonFile, out, 0644)
	println(err)
}

//查询转发服务器管理文件
func FindTfserver(url string) *TFserver {
	ps := &TFserver{}

	data, err := ioutil.ReadFile(url)
	if err != nil {
		logs.Info(err)
	}

	// 这里参数要指定为变量的地址
	err = json.Unmarshal(data, &ps)
	if err != nil {
		logs.Info(err)
	}

	return ps
}

//查询客户端版本管理文件
func FindVersion(url string) *brower_backstage.VersionData {
	ps := &brower_backstage.VersionData{}

	data, err := ioutil.ReadFile(url)
	if err != nil {
		logs.Info(err)
	}

	// 这里参数要指定为变量的地址
	err = json.Unmarshal(data, &ps)
	if err != nil {
		logs.Info(err)
	}

	return ps
}

//修改客户端版本管理文件
func UpdateAll(data *brower_backstage.VersionData, jsonFile string) {
	out, _ := json.MarshalIndent(data, "", "  ")
	err := ioutil.WriteFile(jsonFile, out, 0644)
	println(err)
}

// Paginator 对社交广场动态的的id进行分页处理(因存进去的时候是sadd,后续又要对这个列表进行删除对应的元素,所以LPUSH没法实现)
func Paginator(page, pageSize int32, arr []int64) map[string]interface{} {
	m := make(map[string]interface{})
	arrLength := int32(len(arr))
	totalCount := arrLength            // 总条数
	pageCount := totalCount / pageSize // 页数
	// 如果有余数则添加一页
	if totalCount%pageSize != 0 {
		pageCount++
	}
	var newArr []int64
	m["totalCount"] = totalCount
	m["pageCount"] = pageCount
	m["arr"] = newArr
	startIndex := (page - 1) * pageSize // 偏移量
	if startIndex > int32(len(arr)) {
		return m
	}
	endIndex := startIndex + pageSize
	if arrLength < endIndex { // 如果末端索引大于数组长度
		endIndex = arrLength
	}
	newArr = append(newArr, arr[startIndex:endIndex]...)
	m["arr"] = newArr
	return m
}

// CheckMessageCode 校验平台内部发送的验证码
func CheckMessageCode(phone, code string, t int32) easygo.IMessage {
	if !IS_FORMAL_SERVER { //是测试服直接返回
		return nil
	}
	data := MessageMarkInfo.GetMessageMarkInfo(t, phone)
	if data == nil {
		return easygo.NewFailMsg("验证码不存在")
	}
	if data.Mark != code {
		return easygo.NewFailMsg("验证码不正确")
	}
	return nil
}

// 处理置顶定时任务逻辑公共方法
func ProcessTopTimer(reqMsg *share_message.BackstageNotifyTopReq) easygo.IMessage {
	logs.Info("===========处理定时任务公共方法=============,reqMsg: ", reqMsg)
	if reqMsg.GetLogId() == 0 || reqMsg.GetTopOverTime() < 0 {
		logs.Error("置顶消息参数有误,reqMsg--->", reqMsg)
		return easygo.NewFailMsg("参数有误")
	}
	//dynamic := GetRedisDynamic(reqMsg.GetLogId())
	// 从数据库读取
	dynamic := GetDynamicByStatusSFromDB(reqMsg.GetLogId(), []int{DYNAMIC_STATUE_COMMON, DYNAMIC_STATUE_UNPUBLISHED})
	if dynamic == nil {
		logs.Error("置顶消息,动态不存在,logId: ", reqMsg.GetLogId())
		return easygo.NewFailMsg("置顶消息,动态不存在")
	}
	topReq := &share_message.BackstageNotifyTopReq{
		LogId:       reqMsg.LogId,
		TopOverTime: easygo.NewInt64(dynamic.GetTopOverTime()),
	}
	// 判断是后台置顶还是app置顶
	switch true {
	case reqMsg.GetIsBsTop():
		topReq.IsBsTop = easygo.NewBool(true)
	case reqMsg.GetIsBsTop():
		topReq.IsTop = easygo.NewBool(true)
	default:
		return easygo.NewFailMsg("不是app置顶也不是后台置顶,不处置顶定时任务")
	}

	// 存放定时任务管理器的 map 是在同一个社交广场.
	// 计算定时任务时间,启动定时任务
	localTime := GetMillSecond() / 1000
	t := dynamic.GetTopOverTime() - localTime
	if t < 0 {
		logs.Error("定时任务时间计算结果为负数,t--->", t)
		return easygo.NewFailMsg("定时任务时间计算结果为负数")
	}
	// 封装定时任务的动态

	topTimer := AfterCancelTop(time.Duration(t)*time.Second, dynamic.GetLogId())
	// 存放定时任务进管理器
	TimerMgr.TimerMap.Store(reqMsg.GetLogId(), topTimer)
	// 把 logId 和到期时间存放进数据库,以宕机时加载使用,table: square_top_timer_mgr
	UpsetSquareTopTimerMgrToDB(topReq)
	return nil
}

/*
  方法：WeightedRandomIndex
  功能：按照指定的一组权重随机返回数组索引
  参数：weights []float32 权重切片
  返回：加权随机索引index，index是 0 ~ len(weights)-1 之间的一个整数
  示例如下：
  按权重[0.1, 0.2, 0.3, 0.4]随机调用1000次该方法，返回0,1,2,3的次数将接近于1:2:3:4
  	var weights = []float32{0.1, 0.2, 0.3, 0.4}
  	var result [4]int
  	rand.Seed(time.Now().Unix())
  	for i := 0; i < 1000; i++ {
  		result[WeightedRandomIndex(weights)]++
  	}
	fmt.Printf("%v\n", result)
  输出：
    [112 174 304 410]
*/
func WeightedRandomIndex(weights []float32) int {
	if len(weights) == 1 {
		return 0
	}
	var sum float32 = 0.0
	for _, w := range weights {
		sum += w
	}
	r := rand.Float32() * sum
	var t float32 = 0.0 // 0.6
	for i, w := range weights {
		t += w
		if t > r {
			return i
		}
	}
	return len(weights) - 1
}

// 生成区间[-m, n]的安全随机数
func RangeRand(min, max int64) int64 {
	if min > max {
		panic("the min is greater than max!")
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := crand.Int(crand.Reader, big.NewInt(max+1+i64Min))
		return result.Int64() - i64Min
	} else {
		result, _ := crand.Int(crand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}

func GenerateAesKey() string {
	r := easygo.RandString(16)
	key := base64.StdEncoding.EncodeToString([]byte(r))
	return key
}

// 从数组中随机选取指定条数的数据,返回新的数组
func GetSliceByRandFromSlice(s []int64, count int) []int64 {
	result := make([]int64, 0)
	if len(s) == 0 {
		return result
	}
	if len(s) <= count {
		return s
	}
	ids := make([]int, 0)
	tempIds := make([]int, 0)
	m := make(map[int64]int64)
	for i := 0; i < len(s); i++ {
		tempIds = append(tempIds, i)
		m[s[i]] = s[i]
	}
	for i := 0; i < count; i++ {
		id := RandInt(0, len(tempIds))
		tempIds = easygo.Del(tempIds, id).([]int)
		ids = append(ids, id)
	}
	for _, i := range ids {
		//result = append(result, s[i])
		result = append(result, m[int64(i)])
	}
	return result
}
func DeleteMulti(objects []string) {
	u, _ := url.Parse("https://im-resource-1253887233.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  "AKIDOYR4Xst9ZicIwJrJ7ex2rPlqgY9VbIj1",
			SecretKey: "DDCVaPxb2evwpgJfleZEm4RXPAe7KOCk",
		},
	})
	obs := []cos.Object{}
	for _, v := range objects {
		keys := strings.Split(v, ".myqcloud.com/")
		if len(keys) != 2 {
			continue
		}
		obs = append(obs, cos.Object{Key: keys[1]})
	}
	logs.Info("请求的图片：", obs)
	opt := &cos.ObjectDeleteMultiOptions{
		Objects: obs,
		// 布尔值，这个值决定了是否启动 Quiet 模式
		// 值为 true 启动 Quiet 模式，值为 false 则启动 Verbose 模式，默认值为 false
		//Quiet: true,
	}
	result, rsp, err := client.Object.DeleteMulti(context.Background(), opt)
	if err != nil {
		logs.Error(err)
	}
	logs.Info("DeleteMulti-->>:", result.Errors)
	logs.Info("DeleteMulti-->>:", rsp.Status, rsp.StatusCode)
}

//初始化etcd已存在的服务器
func InitExistServer(pClient3KVMgr *Client3KVManager, pServerInfoMgr *ServerInfoManager, server *share_message.ServerInfo) {
	client := pClient3KVMgr.GetClient()
	kv := pClient3KVMgr.GetClientKV()
	kvs, err := kv.Get(context.TODO(), ETCD_SERVER_PATH, clientv3.WithPrefix())
	easygo.PanicError(err)
	logs.Info("已经启动的服务器:", kvs.Kvs)
	for _, srv := range kvs.Kvs {
		s := &share_message.ServerInfo{}
		err1 := json.Unmarshal(srv.Value, s)
		easygo.PanicError(err1)
		if s.GetSid() == server.GetSid() {
			continue
		}
		logs.Info("服务器:", s.GetSid())
		pServerInfoMgr.AddServerInfo(s)
	}
	//监视login服务器变化
	WatchToServer(client, ETCD_SERVER_PATH, pServerInfoMgr)
}

//监听服务器的变化
func WatchToServer(clt *clientv3.Client, key string, pServerInfoMgr *ServerInfoManager) {
	easygo.Spawn(func() {
		wc := clt.Watch(context.TODO(), key, clientv3.WithPrefix())
		for v := range wc {
			for _, e := range v.Events {
				//logs.Info("type:%v kv:%v  prevKey:%v  value:%v\n ", e.Type, string(e.Kv.Key), e.PrevKv, e.Kv.Value)
				switch e.Type {
				case mvccpb.DELETE: //删除
					//关闭无效连接
					params := strings.Split(string(e.Kv.Key), "/")
					sid := easygo.AtoInt32(params[3])
					pServerInfoMgr.DelServerInfo(sid)
					logs.Info("remove ServerInfo:id=", sid)
				case mvccpb.PUT: //增加
					//如果已经连接
					ss := &share_message.ServerInfo{}
					if err := json.Unmarshal(e.Kv.Value, ss); err != nil {
						logs.Info("WatchToLogin err", err)
						continue
					}
					pServerInfoMgr.AddServerInfo(ss)
					logs.Info("add ServerInfo:", ss)
				}
			}
		}
	})
}

// 推送聊天内容,数字替换成表情.
func ReplaceEmotionStr(s string) string {
	emotion := []string{
		"😃", "😄", "😆", "🤣", "😂", "😋", "😍", "😘", "😚", "😜",
		"🤪", "🤗", "🤭", "🤫", "🤔", "🤐", "🤨", "😏", "😒", "🙄",
		"😬", "😌", "😪", "🤤", "😴", "🤧", "😷", "🤮", "🤒", "😵",
		"🤓", "🤠", "😯", "😳", "🥺", "😦", "😨", "😧", "😭", "😱",
		"😣", "😫", "😤", "😡", "👿", "🥶", "🤩", "👻", "🤬", "😐",
		"😑", "🤥", "🤑", "🤢", "💩", "💞", "🤯", "☠️", "😺", "🙈",
		"💌", "💘", "💯", "👋", "💣", "🤡", "💢", "👹", "💋", "🥳",
	}

	//s := "adoifuad[11]dadpaa[2][11]daidoiad[3][72]"
	r := `(\[#[0-9]+\])+`
	r2 := `([0-9])+`
	// 先匹配取到[[数字1][数字2]]
	compile, _ := regexp.Compile(r)
	allString := compile.FindAllString(s, -1)
	bytes, _ := json.Marshal(allString)
	s1 := string(bytes)
	// 再次匹配 [数字1 数字2]
	compile1, _ := regexp.Compile(r2)
	allString2 := compile1.FindAllString(s1, -1)
	if len(allString2) == 0 {
		return s
	}
	// 遍历第二次取到的内容
	// 把表情封装进map
	eMap := make(map[string]string)
	for _, v := range allString2 {
		// 判断是否大于70
		i, _ := strconv.Atoi(v)
		if i > 70 {
			continue
		}

		fmt.Println(i-1, emotion[i-1])
		eMap[fmt.Sprintf("%s%d%s", "[#", i, "]")] = emotion[i-1]

	}
	// 遍历map,把表情替换到对应的字符串里
	for k, v := range eMap {
		s = strings.ReplaceAll(s, k, v)
	}
	logs.Info("ssss===>>", s)
	return s
}

//utf8转gbk字符串
func Utf8ToGBK(s string) string {
	enc := mahonia.NewEncoder("gbk")
	output := enc.ConvertString(s)
	return output
}

//gbk转utf8
func GBKToUtf8(s string) string {
	dec := mahonia.NewDecoder("gbk")
	output := dec.ConvertString(s)
	return output
}

//获取客户端ip
func GetUserIp(r *http.Request) string {
	ip := ClientPublicIP(r)
	if ip == "" {
		ip = ClientIP(r)
	}
	return ip
}

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// ClientPublicIP 尽最大努力实现获取客户端公网 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" && !HasLocalIPddr(ip) {
			return ip
		}
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" && !HasLocalIPddr(ip) {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		if !HasLocalIPddr(ip) {
			return ip
		}
	}

	return ""
}

// RemoteIP 通过 RemoteAddr 获取 IP 地址， 只是一个快速解析方法。
func RemoteIP(r *http.Request) string {
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// HasLocalIPddr 检测 IP 地址字符串是否是内网地址
func HasLocalIPddr(ip string) bool {
	return HasLocalIP(net.ParseIP(ip))
}

// HasLocalIP 检测 IP 地址是否是内网地址
func HasLocalIP(ip net.IP) bool {
	// for _, network := range localNetworks {
	// 	if network.Contains(ip) {
	// 		return true
	// 	}
	// }

	return ip.IsLoopback()
}

//匹配话题，返回话题数组，nil为没有匹配到
func CheckTopic(s string) []string {
	reg := regexp.MustCompile(`#[\p{Han}0-9a-zA-Z ,，.。!！"“”:;<>《》?？；：、…-]{2,16}#`)
	if reg == nil {
		logs.Error("正则表达式错误")
		return nil
	}
	result := reg.FindAllString(s, -1)
	return result
}

// 切片分页
func SliceByPage(page, pageSize int, data []interface{}) ([]interface{}, int) {
	lenData := len(data)
	totalPage := lenData / pageSize //取整
	r := lenData % pageSize         //求余
	if r == 0 {
		totalPage = totalPage - 1
	}

	for i := 0; i <= totalPage; i++ {
		if i == page-1 {
			start := i * pageSize
			end := (1 + i) * pageSize
			if i == totalPage {
				end = len(data)
			}
			saveData := data[start:end]
			return saveData, lenData
		}

	}
	return nil, 0
}

// 切片分页 start end
func MakeRedisPage(page, pageSize, sliceLen int) (int, int) {
	lenData := sliceLen
	totalPage := lenData / pageSize //取整
	r := lenData % pageSize         //求余
	if r == 0 {
		totalPage = totalPage - 1
	}

	for i := 0; i <= totalPage; i++ {
		if i == page-1 {
			start := i * pageSize
			end := (1 + i) * pageSize
			if i == totalPage {
				end = sliceLen
			}

			return start, end - 1
		}

	}
	return 0, 0
}

//解析服务器间返回的报错参数
func ParseReturnDataErr(reqMsg easygo.IMessage) *base.Fail {
	if nil != reqMsg {
		fail, ok := reqMsg.(*base.Fail)
		if ok {
			return fail
		}
	}
	return nil
}

//列表去重
func RemoveRepeatedElementInt64(arr []int64) (newArr []int64) {
	newArr = make([]int64, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
func RemoveRepeatedElementInt32(arr []int32) (newArr []int32) {
	newArr = make([]int32, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
func RemoveRepeatedElementString(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

//组装私聊会话id
func MakeSessionKey(talk, target int64) string {
	var key string
	if talk > target {
		key = MakeNewString(target, talk)
	} else {
		key = MakeNewString(talk, target)
	}
	return key
}

/**
page:如果有,就返回当前的,如果没有就设置默认的
pageSize:如果有,就返回当前的,如果没有就设置默认的
*/
func MakePageAndPageSize(page, pageSize int32) (int, int) {
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	return int(page), int(pageSize)
}

//获取当天0点时间戳
func FirstSecondTime() int64 {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	time := startTime.Unix()
	return time
}

/**
csvName: csv的文件名
writeData 保存的数据,结构体数组
*/
func WriteCsv(csvName string, writeData interface{}) error {
	clientsFile, err := os.OpenFile(csvName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()
	if _, err := clientsFile.Seek(1, 2); err != nil { // Go to the start of the file
		return err
	}
	return gocsv.MarshalFile(writeData, clientsFile) // Use this to save the CSV back to the file
}

//向语音合成服务器发起合成请求合成请求:
// input_voice_url1=背景
// input_voice_url2=录制的
func MakeVoiceVideo(bgUrl, myUrl string, bgVolume, mixVolume int32) string {
	data := make(easygo.KWAT, 0)
	data.Add("input_voice_url1", UrlEncode(bgUrl))
	data.Add("input_voice_volume1", bgVolume)
	data.Add("input_voice_url2", UrlEncode(myUrl))
	data.Add("input_voice_volume2", mixVolume)
	body, err := json.Marshal(data)
	logs.Info("发送body:", string(body))
	client := &http.Client{}
	// build a new request, but not doing the POST yet
	sUrl := easygo.YamlCfg.GetValueAsString("VOICE_VEDIO_MIX_URL")
	req, err := http.NewRequest("POST", sUrl, bytes.NewReader(body))
	if err != nil {
		logs.Error(err)
		return ""
	}
	// set the Header here
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	// now POST it
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		logs.Error(err)
		return ""
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return ""
	}
	logs.Info("语音合成返回:" + string(result))
	res := make(easygo.KWAT, 0)
	err = json.Unmarshal(result, &res)
	if err != nil {
		logs.Error(err)
		return ""
	}
	return res.GetString("output_voice_url")
}

//对指定存储桶url进行urlencode处理
func UrlEncode(u string) string {
	format := "https://im-resource-1253887233.file.myqcloud.com/backstage/match/voice/"
	params := strings.Split(u, format)
	if len(params) != 2 {
		return u
	}
	return format + url.PathEscape(params[1])
}

//获取当天的剩余秒数
func SurplusTime() int64 {
	currentTime := time.Now()
	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())
	return endTime.Unix() - currentTime.Unix()
}

//获取当前年月日
func GetDateYMD() (int32, int32, int32) {
	time := time.Now()
	year := easygo.AtoInt32(time.Format("2006"))
	month := easygo.AtoInt32(time.Format("01"))
	day := easygo.AtoInt32(time.Format("02"))
	return year, month, day
}
