package hall

import (
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/akqp2019/mgo/bson"

	"github.com/astaxie/beego/logs"
)

/*
	硬币H5服务入口
*/

//首页返回
type IndexData struct {
	Cions     *int64                      `json:"Cions"`     //用户硬币
	BuyCions  *int64                      `json:"BuyCions"`  //用户购买硬币
	GiftCions *int64                      `json:"GiftCions"` //用户获赠硬币
	AdvList   []*share_message.AdvSetting `json:"AdvList"`   //广告列表
}

//硬币记录返回
type LogsData struct {
	LogList []*LogData `json:"LogList"` //日志列表
	Total   int        `json:"Total"`   //日志数量
}

type LogData struct {
	Title      *string `json:"Title"`
	Coin       *int64  `json:"Coin"`
	CreateTime *int64  `json:"CreateTime"`
}

//硬币H5服务入口
func (self *WebHttpServer) CoinsEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	self.NextJobId()
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body) //获取post的数据
	if len(data) == 0 {
		self.OutputJson(w, 0, "参数异常", nil)
		return
	}
	params, err := url.ParseQuery(string(data))
	if err != nil {
		logs.Info("data:", string(data))
		self.OutputJson(w, 0, "解析数据异常", nil)
		return
	}
	// logs.Error("lastTime:", self.LastTime)
	if self.LastTime == easygo.NowTimestamp() {
		time.Sleep(time.Second * 1)
	}
	self.UpdateLastTime()

	switch params.Get("t") {
	case "index", "": //硬币首页
		self.CoinsIndex(w, params, r.RemoteAddr)
	case "logs": //硬币记录
		self.CoinsLog(w, params)
	default:
		self.OutputJson(w, 0, "请求类型错误", nil)
		return
	}
}

//硬币首页
func (self *WebHttpServer) CoinsIndex(w http.ResponseWriter, params url.Values, ip string) {
	logs.Debug("H5硬币页请求IP:" + ip)
	ok, err := self.CheckPlayer(params)
	if !ok {
		self.OutputJson(w, 0, err, nil)
		return
	}

	playerId := easygo.StringToInt64noErr(params.Get("id"))
	base := for_game.GetRedisPlayerBase(playerId)
	if base == nil {
		self.OutputJson(w, 0, "用户不存在", nil)
		return
	}

	advList := for_game.QueryAdvListToDB(for_game.ADV_LOCATION_BANNER_COIN, ip) //5硬币页广告
	data := &IndexData{
		Cions:     easygo.NewInt64(base.GetAllCoin()),
		BuyCions:  easygo.NewInt64(base.GetCoin()),
		GiftCions: easygo.NewInt64(base.GetBCoin()),
		AdvList:   advList,
	}
	self.OutputJson(w, 1, "SUCCESS", data)
}

//硬币记录
func (self *WebHttpServer) CoinsLog(w http.ResponseWriter, params url.Values) {
	logs.Info("=========CoinsLog=============")
	ok, err := self.CheckPlayer(params)
	if !ok {
		self.OutputJson(w, 0, err, nil)
		return
	}

	playerId := easygo.StringToInt64noErr(params.Get("id"))
	pt := easygo.StringToIntnoErr(params.Get("pt"))
	page := easygo.StringToIntnoErr(params.Get("page"))
	pageSize := easygo.StringToIntnoErr(params.Get("pagesize"))

	list, Total := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_COINCHANGELOG, bson.M{"PlayerId": playerId, "PayType": pt}, pageSize, page, "-CreateTime")
	var logList []*LogData
	for _, li := range list {
		logData := &LogData{}
		switch int32(li.(bson.M)["SourceType"].(int)) {
		case for_game.COIN_TYPE_SYSTEM_IN:
			logData.Title = easygo.NewString("系统赠送")
			logData.Coin = easygo.NewInt64(li.(bson.M)["ChangeCoin"])
			logData.CreateTime = easygo.NewInt64(li.(bson.M)["CreateTime"])
		case for_game.COIN_TYPE_EXCHANGE_IN:
			logData.Title = easygo.NewString("兑换")
			logData.Coin = easygo.NewInt64(li.(bson.M)["ChangeCoin"])
			logData.CreateTime = easygo.NewInt64(li.(bson.M)["CreateTime"])
		case for_game.COIN_TYPE_PLAYER_IN:
			playerId := li.(bson.M)["Extend.PlayerId"]
			title := "获得动态投币"
			id, ok := playerId.(int64)
			if ok {
				base := for_game.GetRedisPlayerBase(id)
				if base != nil {
					title = fmt.Sprintf("%s投币了你的动态", base.GetNickName())
				}
			}
			logData.Title = easygo.NewString(title)
			logData.Coin = easygo.NewInt64(li.(bson.M)["ChangeCoin"])
			logData.CreateTime = easygo.NewInt64(li.(bson.M)["CreateTime"])
		case for_game.COIN_TYPE_ACT_PRIZE:
			logData.Title = easygo.NewString("活动奖励")
			logData.Coin = easygo.NewInt64(li.(bson.M)["ChangeCoin"])
			logData.CreateTime = easygo.NewInt64(li.(bson.M)["CreateTime"])
		case for_game.COIN_TYPE_SYSTEM_OUT:
			logData.Title = easygo.NewString("系统回收")
			logData.Coin = easygo.NewInt64(li.(bson.M)["ChangeCoin"])
			logData.CreateTime = easygo.NewInt64(li.(bson.M)["CreateTime"])
		case for_game.COIN_TYPE_SHOP_OUT:
			logsExtend := li.(bson.M)["Extend"]
			productId := logsExtend.(bson.M)["RedPacketId"]
			title := "兑换道具"
			id, ok := productId.(float64)
			if ok {
				product := for_game.GetCoinShopItem(int64(id))
				if product != nil {
					title = fmt.Sprintf("兑换【%s】", product.GetName())
				}
			}
			logData.Title = easygo.NewString(title)
			logData.Coin = easygo.NewInt64(li.(bson.M)["ChangeCoin"])
			logData.CreateTime = easygo.NewInt64(li.(bson.M)["CreateTime"])
		case for_game.COIN_TYPE_PLAYER_OUT:
			playerId := li.(bson.M)["Extend.PlayerId"]
			title := "投币动态"
			id, ok := playerId.(int64)
			if ok {
				base := for_game.GetRedisPlayerBase(id)
				if base != nil {
					title = fmt.Sprintf("你投币了%s的动态", base.GetNickName())
				}
			}
			logData.Title = easygo.NewString(title)
			logData.Coin = easygo.NewInt64(li.(bson.M)["ChangeCoin"])
			logData.CreateTime = easygo.NewInt64(li.(bson.M)["CreateTime"])
		case for_game.COIN_TYPE_EXPIRED_OUT:
			logData.Title = easygo.NewString("过期回收")
			logData.Coin = easygo.NewInt64(li.(bson.M)["ChangeCoin"])
			logData.CreateTime = easygo.NewInt64(li.(bson.M)["CreateTime"])
		case for_game.COIN_TYPE_ESPORT_EXCHANGE_OUT:
			logData.Title = easygo.NewString("电竞兑换")
			logData.Coin = easygo.NewInt64(li.(bson.M)["ChangeCoin"])
			logData.CreateTime = easygo.NewInt64(li.(bson.M)["CreateTime"])
		case for_game.COIN_TYPE_WISH_PAY:
			logData.Title = easygo.NewString("许愿池兑换钻石")
			logData.Coin = easygo.NewInt64(li.(bson.M)["ChangeCoin"])
			logData.CreateTime = easygo.NewInt64(li.(bson.M)["CreateTime"])
		default:
			logData.Title = easygo.NewString("未知")
			logData.Coin = easygo.NewInt64(li.(bson.M)["ChangeCoin"])
			logData.CreateTime = easygo.NewInt64(li.(bson.M)["CreateTime"])
		}
		logList = append(logList, logData)
	}

	data := &LogsData{
		LogList: logList,
		Total:   Total,
	}
	self.OutputJson(w, 1, "SUCCESS", data)
}

//验证用户身份
func (self *WebHttpServer) CheckPlayer(params url.Values) (bool, string) {
	playerId := easygo.StringToInt64noErr(params.Get("id"))
	base := for_game.GetRedisPlayerBase(playerId)
	if base == nil {
		return false, "用户不存在"
	}

	if !base.GetIsOnLine() {
		return false, "用户不存在"
	}

	token := params.Get("token")
	if token == "" || token != for_game.Md5(base.GetToken()) {
		return false, "token错误"
	}

	return true, ""
}
