package hall

import (
	"encoding/base64"
	"encoding/json"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/share_message"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/astaxie/beego/logs"
)

/*
统统付支付模块逻辑
*/

const TONG_PAY_WX = "H5_WXJSAPI"     //"H5_WXJSAPI"
const TONG_PAY_ZFB = "API_ZFBQRCODE" //API_ZFBQRCODE

type WebTongTongPay struct {
	Account     string //小程序appId
	SysCode     string //小程序openId
	MD5Key      string //密钥
	Url         string //接口地址
	ApiPay      string //支付api接口
	ApiCheck    string //查询api接口
	CallBackUrl string //回调地址
}

func NewWebTongTongPay() *WebTongTongPay {
	p := &WebTongTongPay{}
	p.Init()
	return p
}
func (self *WebTongTongPay) Init() {
	self.Account = easygo.YamlCfg.GetValueAsString("TTP_MERCHANTNO")
	self.SysCode = easygo.YamlCfg.GetValueAsString("TTP_SYSCODE")
	self.MD5Key = easygo.YamlCfg.GetValueAsString("TTP_MD5_KEY")
	self.Url = easygo.YamlCfg.GetValueAsString("TTP_URL")
	self.ApiPay = "web_pay.api"
	self.ApiCheck = "get_pay_result.api"
	self.CallBackUrl = easygo.YamlCfg.GetValueAsString("TTP_PAY_CALLBACK_URL")
}

//签名消息体，并返回base64数据

func (self *WebTongTongPay) SignMsg(values url.Values) string {
	if values == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(values))
	for k := range values {

		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if k == "signature" { //signature不参与签名源串
			continue
		}
		vs := values[k]
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(v)
		}
	}
	s := buf.String() + "&key=" + self.MD5Key
	sign := for_game.Md5(s)
	return strings.ToUpper(sign)
}
func (self *WebTongTongPay) CheckSign(values url.Values) bool {
	sign := self.SignMsg(values)
	if sign == values.Get("signature") {
		return true
	}
	return false
}

//发起充值：
//微信或者支付宝小程序支付
func (self *WebTongTongPay) RechargeOrder(web *WebHttpServer, payData *share_message.PayOrderInfo, OpenId string) ([]byte, *base.Fail) {
	orderId, ip := self.CreateOrder(payData, OpenId)
	payMode := TONG_PAY_WX
	if payData.GetPayType() == for_game.PAY_TYPE_ZHIFUBAO {
		payMode = TONG_PAY_ZFB
	}
	amount := payData.GetAmount()
	ama := easygo.AtoFloat64(amount) * 100
	amountStr := easygo.AnytoA(ama)

	logs.Info("tongtongpay 订单金额为----->%s,单位为分", amountStr)

	params := make(url.Values)
	params.Set("syscode", self.SysCode)                    //签名需要
	params.Set("trans_time", for_game.GetCurTimeString3()) //签名需要
	params.Set("account", self.Account)                    //签名需要
	//params.Set("amount", "1")                              //签名需要
	params.Set("amount", amountStr)            //签名需要
	params.Set("pay_mode", payMode)            //签名需要
	params.Set("app_id", orderId)              //签名需要
	params.Set("notify_url", self.CallBackUrl) //签名需要
	sign := self.SignMsg(params)
	params.Set("aging", "1")
	params.Set("terminal_ip", ip)
	params.Set("subject", for_game.Utf8ToGBK("柠檬充值"))
	params.Set("body", for_game.Utf8ToGBK(payData.GetProduceName()))
	params.Set("terminal_device", "4")
	params.Set("openid", OpenId) //微信支付宝用户标识
	params.Set("auth_app_id", easygo.YamlCfg.GetValueAsString("TTP_WXAPPID"))
	params.Set("signature", sign)
	res := self.ReqHttps(self.Url+self.ApiPay, params)
	logs.Info("支付返回:", string(res))
	signData := make(url.Values)
	respData := make(easygo.KWAT)
	err := json.Unmarshal(res, &respData)
	easygo.PanicError(err)
	for k, v := range respData {
		if k != "errorcode" && k != "errormessage" {
			signData.Add(k, v.(string))
		}
	}
	orderId = respData.GetString("app_id")
	eCode := respData.GetString("errorcode")
	if eCode != "0000" {
		return nil, easygo.NewFailMsg("第三方支付下单请求失败")
	}
	if !self.CheckSign(signData) {
		return nil, easygo.NewFailMsg("第三方支付验签失败")
	}
	order := for_game.GetRedisOrderObj(orderId)
	if order == nil {
		return nil, easygo.NewFailMsg("无效的订单")
	}
	//支付宝采用扫码支付，直接返回url地址
	payInfo := respData.GetString("pay_url")
	urlType := respData.GetString("url_type")
	if urlType == "8" {
		//微信支付小程序支付要解析payInfo
		payUrl, err := base64.StdEncoding.DecodeString(respData.GetString("pay_url"))
		if err != nil {
			return nil, easygo.NewFailMsg("base64 Decode err:" + err.Error())
		}
		payInfo1 := for_game.GBKToUtf8(string(payUrl))
		v1, err := url.ParseQuery(payInfo1)
		if err != nil {
			return nil, easygo.NewFailMsg("ParseQuery err:" + err.Error())
		}
		v2 := easygo.KWAT{}
		for k, v := range v1 {
			v2.Add(k, v[0])
		}
		payInfo2, err := json.Marshal(v2)
		if err != nil {
			return nil, easygo.NewFailMsg("json.Marshal(v2):" + err.Error())
		}
		payInfo = string(payInfo2)
	}
	logs.Info("payInfo:", payInfo)
	pData := &ParseMDData{
		Data: PayInfo{
			PayType:        0,
			PreparePayInfo: payInfo,
			OrderId:        orderId,
		},
		Result:  true,
		Code:    200,
		Message: respData.GetString("errormessage"),
	}
	d, err := json.Marshal(pData)
	if err != nil {
		return nil, easygo.NewFailMsg("json.Marshal(pData):" + err.Error())
	}
	logs.Info("第三方请求结果:", string(d))
	//启动定时器，检测订单是否完成
	//fun := func() {
	//	//5分钟支付时间，如果没完成支付，主动取消订单
	//	self.CheckOrder(orderId)
	//}
	//easygo.AfterFunc(5*60*time.Second, fun)
	return d, nil
}

//支付5分钟后查询订单
func (self *WebTongTongPay) CheckOrder(orderId string) {
	params := make(url.Values)
	params.Set("syscode", self.SysCode)                    //签名需要
	params.Set("trans_time", for_game.GetCurTimeString3()) //签名需要
	params.Set("account", self.Account)                    //签名需要
	params.Set("app_id", orderId)                          //签名需要
	sign := self.SignMsg(params)
	params.Set("signature", sign) //签名
	res := self.ReqHttps(self.Url+self.ApiCheck, params)
	respData := make(easygo.KWAT)
	err := json.Unmarshal(res, &respData)
	easygo.PanicError(err)
	logs.Info("查询结果:", orderId, respData)
	//if respData.GetString("result") == "3" {
	//	//取消待支付的订单
	//	self.CancelOrder(respData.GetString("app_id"), respData.GetString("trans_id"))
	//}
}

//取消订单接口:pay_cancel.api
func (self *WebTongTongPay) CancelOrder(orderId, transId string) {
	params := make(url.Values)
	params.Set("syscode", self.SysCode)                    //签名需要
	params.Set("trans_time", for_game.GetCurTimeString3()) //签名需要
	params.Set("account", self.Account)                    //签名需要
	params.Set("app_id", orderId)                          //签名需要
	params.Set("trans_id", transId)                        //签名需要
	sign := self.SignMsg(params)
	params.Set("signature", sign) //签名
	res := self.ReqHttps(self.Url+"pay_cancel.api", params)
	respData := make(easygo.KWAT)
	err := json.Unmarshal(res, &respData)
	easygo.PanicError(err)
	logs.Info("取消结果:", orderId, respData)
	if respData.GetString("errorcode") == "0000" {
		if respData.GetString("result") == "1" {
			order := for_game.GetRedisOrderObj(respData.GetString("app_id"))
			if order != nil {
				order.SetPayStatus(for_game.PAY_ST_CANCEL)
				order.SetStatus(for_game.ORDER_ST_CANCEL)
				order.SaveToMongo()
			}
		}
	}
}

// http post 请求
func (self *WebTongTongPay) ReqHttps(uurl string, params url.Values) []byte {
	//logs.Info("发送请求:url=", uurl, params)
	resp, err := http.PostForm(uurl, params)
	easygo.PanicError(err)
	defer resp.Body.Close()
	result, err1 := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err1)
	//logs.Info("返回数据:", for_game.GBKToUtf8(string(result)))
	return result
}
func (self *WebTongTongPay) CreateOrder(payData *share_message.PayOrderInfo, openId string) (string, string) {
	tax := int64(0)
	amount := easygo.AtoFloat64(payData.GetAmount())
	changeGold := int64(amount * 100)
	player := GetPlayerObj(payData.GetPlayerId())
	channel := for_game.PAY_CHANNEL_TONGTONG_WX
	if payData.GetPayType() == for_game.PAY_TYPE_ZHIFUBAO {
		channel = for_game.PAY_CHANNEL_TONGTONG_ZFB
	}
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
		CreateIP:    easygo.NewString(""),
		Status:      easygo.NewInt32(for_game.ORDER_ST_WAITTING),
		PayStatus:   easygo.NewInt32(for_game.PAY_ST_WAITTING),
		Note:        easygo.NewString(note),
		Tax:         easygo.NewInt64(-tax),
		Operator:    easygo.NewString("system"),
		PayChannel:  easygo.NewInt32(channel),
		PayType:     easygo.NewInt32(payData.GetPayType()),
		PayWay:      easygo.NewInt32(payData.GetPayWay()),
		PayTargetId: easygo.NewInt64(payData.GetPayTargetId()),
		PayOpenId:   easygo.NewString(openId),
		TotalCount:  easygo.NewInt32(payData.GetTotalCount()),
		Content:     easygo.NewString(payData.GetContent()),
		ExtendValue: easygo.NewString(payData.GetExtendValue()),
	}
	obj := for_game.CreateRedisOrder(order)
	return obj.GetOrderId(), player.GetLastLoginIP()
}
