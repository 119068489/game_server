package hall

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

/*
	云通商支付服务入口
*/

type WebYunTongShangPay struct {
	WXAppId     string   // wx appid
	ZFBAppId    string   // zfb appid
	Url         string   // 请求支付地址  http post请求
	CallBackUrl string   // 支付回调地址
	ReturnURL   string   // 前端跳回页面
	SecretKey   string   //md5密钥
	BankList    []string //支持的银行卡列表
}

func NewWebYunTongShangPay() *WebYunTongShangPay {
	p := &WebYunTongShangPay{}
	p.Init()
	return p
}
func (self *WebYunTongShangPay) Init() {
	self.WXAppId = easygo.YamlCfg.GetValueAsString("YTS_WXAPPID")
	self.Url = easygo.YamlCfg.GetValueAsString("YTS_URL")
	self.CallBackUrl = easygo.YamlCfg.GetValueAsString("YTS_PAY_CALLBACK_URL") // 测试www.lemonchat.cn/testhuichaopay
	self.SecretKey = easygo.YamlCfg.GetValueAsString("YTS_MD5_KEY")
	//汇潮支持的银行
	self.BankList = []string{"ICBC", "ABC", "CCB", "BOC", "BOCOM", "BOS", "BCCB", "CEB", "CIB", "CMB", "CMBC", "CNCB", "GDB", "HXB", "PAB", "PSBC", "SPDB"}
}

//获取支持的银行卡列表
func (self *WebYunTongShangPay) GetSupportBankList() []*client_hall.BankData {
	banks := make([]*client_hall.BankData, 0)
	for _, code := range self.BankList {
		if name, ok := for_game.BankName[code]; ok {
			bank := &client_hall.BankData{
				Code: easygo.NewString(code),
				Name: easygo.NewString(name),
			}
			banks = append(banks, bank)
		}
	}
	return banks
}

//签名消息体
func (self *WebYunTongShangPay) SignMsg(params url.Values) string {
	str := self.Encode(params) + "&secret=" + self.SecretKey
	logs.Info("SignMsg 加密源串:", str)
	sign := strings.ToUpper(for_game.Md5(str))
	return sign
}

//校验签名
func (self *WebYunTongShangPay) VerifySign(params YTSResp) bool {
	signParams := url.Values{}
	signParams.Add("amount", easygo.AnytoA(params.Amount))
	signParams.Add("sys_order_id", params.SysOrderID)
	signParams.Add("order_id", params.OrderID)
	signParams.Add("app_id", params.AppID)
	signParams.Add("type", params.Type)
	sign := params.Sign
	newSign := self.SignMsg(signParams)
	return sign == newSign
}
func (self *WebYunTongShangPay) Encode(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(v)
		}
	}
	return buf.String()
}

// http post 请求
func (self *WebYunTongShangPay) HttpsReq(uurl string, params url.Values) []byte {
	logs.Info("请求参数:", params.Encode())
	resp, err := http.PostForm(uurl, params)
	easygo.PanicError(err)
	defer resp.Body.Close()
	result, err1 := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err1)
	logs.Info("返回数据:", string(result))
	return result
}

func (self *WebYunTongShangPay) CreateOrder(payData *share_message.PayOrderInfo) string {
	amount := easygo.AtoFloat64(payData.GetAmount())
	tax := int64(0)
	changeGold := int64(amount * 100)
	player := GetPlayerObj(payData.GetPlayerId())
	channel := for_game.PAY_CHANNEL_YTS_WX
	note := for_game.GetRechargeNote(payData.GetPayWay())
	order := &share_message.Order{
		PlayerId:    easygo.NewInt64(player.GetPlayerId()),
		Account:     easygo.NewString(player.GetAccount()),
		NickName:    easygo.NewString(player.GetNickName()),
		RealName:    easygo.NewString(player.GetRealName()),
		SourceType:  easygo.NewInt32(for_game.GOLD_TYPE_CASH_IN),
		ChangeType:  easygo.NewInt32(for_game.GOLD_CHANGE_TYPE_IN),
		Channeltype: easygo.NewInt32(for_game.CHANNEL_OTHER),
		CurGold:     easygo.NewInt64(player.GetGold()),
		ChangeGold:  easygo.NewInt64(changeGold),
		Gold:        easygo.NewInt64(player.GetGold() + changeGold - tax),
		Amount:      easygo.NewInt64(changeGold),
		CreateTime:  easygo.NewInt64(for_game.GetMillSecond()),
		CreateIP:    easygo.NewString(player.GetLastLoginIP()),
		Status:      easygo.NewInt32(for_game.ORDER_ST_WAITTING),
		PayStatus:   easygo.NewInt32(for_game.PAY_ST_WAITTING),
		Note:        easygo.NewString(note),
		Tax:         easygo.NewInt64(-tax),
		Operator:    easygo.NewString("system"),
		PayChannel:  easygo.NewInt32(channel),
		PayType:     easygo.NewInt32(payData.GetPayType()),
		PayWay:      easygo.NewInt32(payData.GetPayWay()),
		PayTargetId: easygo.NewInt64(payData.GetPayTargetId()),
		//PayOpenId:   easygo.NewString(openId),
		TotalCount:  easygo.NewInt32(payData.GetTotalCount()),
		Content:     easygo.NewString(payData.GetContent()),
		ExtendValue: easygo.NewString(payData.GetExtendValue()),
		BankInfo:    easygo.NewString(payData.GetPayBankNo()),
	}
	obj := for_game.CreateRedisOrder(order)
	return obj.GetOrderId()
}

//发起支付接口
func (self *WebYunTongShangPay) RechargeOrder(web *WebHttpServer, payData *share_message.PayOrderInfo) ([]byte, *base.Fail) {
	player := for_game.GetRedisPlayerBase(payData.GetPlayerId())
	if player == nil {
		return nil, easygo.NewFailMsg("无效的用户")
	}
	ips := strings.Split(player.GetLastLoginIP(), ":")
	logs.Error("ip:", ips)
	orderNo := self.CreateOrder(payData)
	params := url.Values{}
	params.Add("app_id", self.WXAppId)
	params.Add("goods_title", "柠檬支付")
	params.Add("goods_desc", payData.GetProduceName())
	params.Add("order_id", orderNo)
	params.Add("amount", payData.GetAmount())
	params.Add("type", "wx")
	params.Add("method", "wap")
	params.Add("notify_url", self.CallBackUrl)
	sign := self.SignMsg(params)
	if len(ips) > 0 {
		params.Add("ip", ips[0]) //必传
	} else {
		params.Add("ip", "0.0.0.0") //必传
	}
	params.Add("sign", sign)
	//订单号生成
	res := self.HttpsReq(self.Url+"pay", params)
	logs.Info("支付返回:", string(res))
	data := easygo.KWAT{}
	er := json.Unmarshal([]byte(res), &data)
	if er != nil {
		logs.Error("解析数据异常:err ", er.Error())
		return nil, easygo.NewFailMsg("json.Marshal(pData):" + er.Error())
	}
	if data.GetString("code") != "000001" {
		//下单成功
		return nil, easygo.NewFailMsg("第三方支付下单请求失败")
	}
	//payInfo := data.GetString("data")
	//pData := &ParseMDData{
	//	Data: PayInfo{
	//		PreparePayInfo: payInfo},
	//	Result:  true,
	//	Code:    200,
	//	Message: data.GetString("result"),
	//}
	//d, err := json.Marshal(pData)
	//if err != nil {
	//	return nil, easygo.NewFailMsg("json.Marshal(pData):" + err.Error())
	//}
	//启动线程查询订单
	fun := func() {
		self.ReqCheckPayOrder(orderNo, time.Second)
	}
	//10秒后开始查询结构
	easygo.AfterFunc(time.Second*10, fun)
	return res, nil
}
func (self *WebYunTongShangPay) GetPerTime(t time.Duration) time.Duration {
	if t >= 0 && t < 60*time.Second {
		return 10 * time.Second //大于0秒小于1分钟，每10秒执行一次
	} else if t >= 60*time.Second && t < 300*time.Second {
		return 30 * time.Second //大于1分钟小于5分钟，每30秒执行一次
	} else if t >= 300*time.Second && t < 1800*time.Second {
		return 60 * time.Second //大于5分钟小于30分钟，每60秒执行一次
	}
	return 600 * time.Second //5分钟执行一次
}

type YTSQueryResp struct {
	Code   string  `json:"code"`
	Method string  `json:"method"`
	Result string  `json:"result"`
	Data   YTSData `json:"data"`
}
type YTSData struct {
	OrderInfo string `json:"order_info"`
}

//查询订单接口
func (self *WebYunTongShangPay) ReqCheckPayOrder(orderId string, t time.Duration) {
	perTime := self.GetPerTime(t)
	if for_game.CheckOrderDoing(orderId) {
		logs.Info("订单正在锁住处理，稍后再试:", orderId)
		return
	}
	for_game.LockOrderDoing(orderId)         //订单上锁
	defer for_game.UnLockOrderDoing(orderId) //处理完解锁
	order := for_game.GetRedisOrderObj(orderId)
	if order == nil {
		logs.Info("无效的订单号:", orderId)
		return
	}
	if order.GetStatus() == for_game.ORDER_ST_FINISH {
		logs.Info("已经处理过的订单:", orderId)
		return
	}
	if order.GetPayStatus() == for_game.PAY_ST_CANCEL {
		logs.Info("订单已取消:", orderId)
		return
	}

	if t > 3600*time.Second { // 超过半小时,说明此订单有问题,直接取消状态,后续人工查单处理.
		//   订单状态和支付状态改成取消状态
		order.SetPayStatus(for_game.PAY_ST_CANCEL)
		order.SetStatus(for_game.ORDER_ST_CANCEL)
		logs.Info("云通商付手动查询订单超过半小时无结果，订单状态和支付状态转取消,orderId:--->", orderId)
		return
	}
	t += perTime

	params := url.Values{}
	params.Add("app_id", self.WXAppId)
	params.Add("method", "query")
	params.Add("order_id", orderId)
	sign := self.SignMsg(params)
	params.Add("sign", sign)
	replyByge := self.HttpsReq(self.Url+"query", params)
	logs.Error("查询订单结果:", string(replyByge))
	resplay := YTSQueryResp{}
	err := json.Unmarshal(replyByge, &resplay)
	if err != nil {
		logs.Error("err")
	}
	logs.Info("云通商支付手动查询订单返回结果------------>%+v", resplay)
	if resplay.Result == "success" { // 如果成功, 执行发货功能
		if len(resplay.Data.OrderInfo) > 0 {
			orderInfo := easygo.KWAT{}
			_ = json.Unmarshal([]byte(resplay.Data.OrderInfo), &orderInfo)
			paySt := orderInfo.GetString("status")
			if paySt == "success" || paySt == "fail" {
				res := false
				if paySt == "success" { //支付成功
					order.SetPayStatus(for_game.PAY_ST_FINISH)
					//order.SetStatus(for_game.ORDER_ST_CANCEL)
					RechargeGoldToPlayer(order.GetPlayerId(), order.GetOrderId())
					res = true
					logs.Info("云通商发货成功:", orderId)
				} else if paySt == "fail" { //支付失败
					order.SetPayStatus(for_game.PAY_ST_FAIL)
					order.SetStatus(for_game.ORDER_ST_CANCEL)
					order.SaveToMongo()
				}
				//如果充值，通知前端充值结果
				if order.GetSourceType() == for_game.GOLD_TYPE_CASH_IN {
					msg := &share_message.RechargeFinish{
						Amount:        easygo.NewInt64(order.GetAmount()),
						TradeNo:       easygo.NewString(order.GetOrderId()),
						PayFinishTime: easygo.NewInt64(order.GetOverTime()),
						Result:        easygo.NewBool(res),
					}
					logs.Debug("ReqCheckPayOrder 通知前端支付结果:", msg)
					SendMsgToHallClientNew(order.GetPlayerId(), "RpcRechargeMoneyFinish", msg)
				}
				logs.Info("云通商支付手动查询订单结束,orderId:", orderId, paySt)
				return //跳出循环
			} else {
				logs.Info("云通商支付订单待支付或者支付中:", orderId)
			}
		}
	} else {
		logs.Info("云通商支付手动查询订单失败,orderId: %s", orderId)
	}
	// 2 的倍数 秒后重新查询
	fun := func() {
		self.ReqCheckPayOrder(orderId, t)
	}
	easygo.AfterFunc(perTime, fun)
}

//退款接口 TODO 暂不处理
func (self *WebYunTongShangPay) CancelPayOrder(orderId string) {
	params := url.Values{}
	params.Add("app_id", self.WXAppId)
	params.Add("method", "refund")
	params.Add("order_id", orderId)
	sign := self.SignMsg(params)
	params.Add("sign", sign)
	resData := self.HttpsReq(self.Url+"refund", params)
	logs.Info("退款响应:", resData)
}
