package hall

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

/*
	汇聚支付服务入口
*/

type WebHuiJuPay struct {
	MerchantNo string //商户号
	//PubKey        string //解密公钥
	//PriKey        string //加密私钥
	PassWord      string //加密密钥
	Url           string //请求支付地址  http post请求
	DaiFuUrl      string //请求代付提现地址 http post 请求
	CallBackUrl   string //支付回调地址
	CallBackUrlDF string //代付回调地址
	MD5Key        string //代付密钥
}

//源明文不能存在字样
var HUIJU_EXT_STR = []string{"0x0d", "0x0a", "0x4f", "=", "|", "%", "《", "<", "+", "￥", "$", "\r\n", "！", "@", "#", "￥", "*", "&", "NULL", "<html>", "if", "while", "or"}

const (
	HJ_ORDER_STATUS_SUCCESS       = "P1000" //交易成功
	HJ_ORDER_STATUS_FAIL          = "P2000" //交易失败
	HJ_ORDER_STATUS_IN_PROCESSING = "P3000" //交易处理中
	HJ_ORDER_STATUS_CANCEL        = "P4000" //订单已取消
	HJ_ORDER_STATUS_CONTROL       = "P5000" //风控阻断
	HJ_ORDER_STATUS_CLOSE         = "P6000" //订单已关闭

	HJ_OSM_DATE_OSM_TIME = 6 //汇聚每天绑卡短信请求次数
)

var SUPPORT_BANK_LIST = []string{"ICBC", "BOC", "CITIC", "SHBANK", "CEB", "CMBC", "SPABANK", "CCB", "SPDB", "PSBC", "GDB"}

//代付签名顺序
var DF_HMAC_LIST = []string{"userNo", "productCode", "requestTime", "merchantOrderNo", "receiverAccountNoEnc",
	"receiverNameEnc", "receiverAccountType", "receiverBankChannelNo", "paidAmount", "currency", "isChecked", "paidDesc",
	"paidUse", "callbackUrl", "firstProductCode"}

var DF_HMAC_RSP_LIST = []string{"statusCode", "message", "status", "errorCode", "errorDesc", "userNo", "merchantOrderNo", "platformSerialNo", "receiverAccountNoEnc", "receiverNameEnc", "paidAmount", "fee"}

func NewWebHuiJuPay() *WebHuiJuPay {
	p := &WebHuiJuPay{}
	p.Init()
	return p
}
func (self *WebHuiJuPay) Init() {
	self.Url = "https://api.joinpay.com/fastpay"
	self.DaiFuUrl = "https://www.joinpay.com/payment/pay/singlePay"
	self.MerchantNo = "888109000009160"
	self.PassWord = "1928374655647382" //长度16的密钥串
	self.CallBackUrl = "http://127.0.0.1:2601/huiju"
	self.CallBackUrlDF = "127.0.0.1:2601/huijuDF"
	self.MD5Key = "5cd119f4a41e4b179414f4ac0e704055"
}

func (self *WebHuiJuPay) SignMsg(data string) string {
	hashMd5 := md5.Sum([]byte(data))
	hashed := hashMd5[:]
	block, _ := pem.Decode(for_game.HJPrivateKey)
	if block == nil {
		logs.Info("block=nil")
		return ""
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	easygo.PanicError(err)
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv.(*rsa.PrivateKey), crypto.MD5, hashed)
	logs.Info("SignMsg:", base64.StdEncoding.EncodeToString(signature))
	return base64.StdEncoding.EncodeToString(signature)
}
func (self *WebHuiJuPay) VerifySign(data, sign string) bool {
	bSign, err := base64.StdEncoding.DecodeString(sign)
	block, _ := pem.Decode([]byte(for_game.HJPublicKey))
	if block == nil {
		return false
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	easygo.PanicError(err)
	hashMd5 := md5.Sum([]byte(data))
	hashed := hashMd5[:]
	res := rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.MD5, hashed, bSign)
	return res == nil
}

//平台公钥对AES密钥进行加密
func (self *WebHuiJuPay) PubEncode(pubKey, data string) string {
	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		logs.Info("block=nil")
		return ""
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	easygo.PanicError(err)
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), []byte(data))
	logs.Info("PubEncode:", base64.StdEncoding.EncodeToString(encryptedData))
	return base64.StdEncoding.EncodeToString(encryptedData)
}

//对拼接的参数字符串再次在头部和尾部拼接分配给商户的 key
func (self *WebHuiJuPay) GetSourceStr(values easygo.KWAT, isUtf8 ...bool) string {
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
		if k == "sign" || k == "aesKey" || k == "aes_key" || k == "sec_key" { //sign不参与签名源串
			continue
		}
		v := values[k]
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(easygo.AnytoA(v))
	}
	return buf.String()
}

//代付提现拼接
func (self *WebHuiJuPay) GetDFHMac(values easygo.KWAT, keys []string) string {
	logs.Info("GetSourceStrDF:", values)
	if values == nil {
		return ""
	}
	//var buf strings.Builder
	var data map[string]interface{}
	if values["data"] != nil {
		data = values["data"].(map[string]interface{})
	}
	src := ""
	for _, k := range keys {
		v := values[k]
		if k == "hmac" { //sign不参与签名源串
			continue
		}
		if v == nil {
			if data == nil {
				continue
			} else {
				if data[k] != nil {
					if k == "fee" || k == "paidAmount" {
						v = fmt.Sprintf("%.2f", data[k])
						data[k] = v
						src += easygo.AnytoA(v)
					} else {
						src += easygo.AnytoA(data[k])
					}
				}
			}
		} else {
			src += easygo.AnytoA(v)
		}
	}
	logs.Info("签名源串:", src)
	hmac := for_game.Md5(src + self.MD5Key)
	logs.Info("签名:", hmac)
	return hmac
}

func (self *WebHuiJuPay) CheckData(data []byte) (easygo.KWAT, *base.Fail) {
	params := make(easygo.KWAT)
	err := json.Unmarshal([]byte(data), &params)
	easygo.PanicError(err)
	srcStr := self.GetSourceStr(params)
	if params.GetString("resp_code") != "SUCCESS" {
		return nil, easygo.NewFailMsg("响应结果:请求失败")
	}
	if self.VerifySign(srcStr, params["sign"].(string)) {
		logs.Info("响应成功验签")
		if params.GetString("biz_code") == "JS000000" {
			//业务响应码:成功
			jsData := params.GetString("data")
			logs.Info("jsData:", jsData)
			mpData := make(easygo.KWAT)
			err := json.Unmarshal([]byte(jsData), &mpData)
			easygo.PanicError(err)
			return mpData, nil
		} else {
			//业务响应失败
			logs.Info("error:", params.GetString("biz_msg"))
			if params.GetString("biz_code") == "JS100008" {
				//风控描述修改
				return nil, easygo.NewFailMsg("银行风控中,请稍后再试")
			}
			return nil, easygo.NewFailMsg(params.GetString("biz_msg"))
		}

	} else {
		logs.Info("响应数据验签失败")
		return nil, easygo.NewFailMsg("响应数据验签失败")
	}
}

func (self *WebHuiJuPay) CheckDataDF(data []byte) (easygo.KWAT, *base.Fail) {
	params := make(easygo.KWAT)
	err := json.Unmarshal(data, &params)
	hmac := self.GetDFHMac(params, DF_HMAC_RSP_LIST)
	jsData := easygo.AnytoA(params["data"])
	mpData := make(easygo.KWAT)
	err = json.Unmarshal([]byte(jsData), &mpData)
	easygo.PanicError(err)
	if hmac != mpData.GetString("hmac") {
		return mpData, easygo.NewFailMsg("响应数据验签失败", "2000")
	}
	return mpData, easygo.NewFailMsg(params.GetString("message"), params.GetString("statusCode"))
}

//生成支付订单
func (self *WebHuiJuPay) CreateOrder(payData *share_message.PayOrderInfo, player *Player) string {
	amount := easygo.AtoFloat64(payData.GetAmount())
	tax := int64(0)
	changeGold := int64(amount * 100)
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
		PayChannel:  easygo.NewInt32(for_game.PAY_CHANNEL_HUIJU),
		PayType:     easygo.NewInt32(payData.GetPayType()),
		PayWay:      easygo.NewInt32(payData.GetPayWay()),
		PayTargetId: easygo.NewInt64(payData.GetPayTargetId()),
		TotalCount:  easygo.NewInt32(payData.GetTotalCount()),
		Content:     easygo.NewString(payData.GetContent()),
		ExtendValue: easygo.NewString(payData.GetExtendValue()),
		BankInfo:    easygo.NewString(payData.GetPayBankNo()),
		OrderDate:   easygo.NewString(for_game.GetCurTimeString2()),
	}
	obj := for_game.CreateRedisOrder(order)
	return obj.GetOrderId()
}

//生成代付订单
//生成内部订单
func (self *WebHuiJuPay) CreateDFOrder(payData *client_hall.WithdrawInfo, player *Player, setting *share_message.PaymentSetting) string {
	amount := payData.GetAmount()
	tax := payData.GetTax()
	changeGold := amount
	order := &share_message.Order{
		PlayerId:    easygo.NewInt64(player.GetPlayerId()),
		Account:     easygo.NewString(player.GetAccount()),
		NickName:    easygo.NewString(player.GetNickName()),
		RealName:    easygo.NewString(player.GetRealName()),
		SourceType:  easygo.NewInt32(for_game.GOLD_TYPE_CASH_OUT),
		ChangeType:  easygo.NewInt32(for_game.GOLD_CHANGE_TYPE_OUT),
		Channeltype: easygo.NewInt32(for_game.CHANNEL_OTHER),
		CurGold:     easygo.NewInt64(player.GetGold()),
		ChangeGold:  easygo.NewInt64(-changeGold),
		Gold:        easygo.NewInt64(player.GetGold() - changeGold - tax),
		Amount:      easygo.NewInt64(changeGold),
		CreateTime:  easygo.NewInt64(for_game.GetMillSecond()),
		CreateIP:    easygo.NewString(player.GetLastLoginIP()),
		Status:      easygo.NewInt32(for_game.ORDER_ST_WAITTING),
		PayStatus:   easygo.NewInt32(for_game.PAY_ST_WAITTING),
		Note:        easygo.NewString("提现"),
		Tax:         easygo.NewInt64(-tax),
		PlatformTax: easygo.NewInt64(-setting.GetPlatformTax()),
		RealTax:     easygo.NewInt64(-setting.GetRealTax()),
		Operator:    easygo.NewString("system"),
		PayChannel:  easygo.NewInt32(for_game.PAY_CHANNEL_HUIJU_DF),
		BankInfo:    easygo.NewString(payData.GetAccountNo()),
		BankCode:    easygo.NewString(payData.GetBankCode()),
		AccountType: easygo.NewString(payData.GetAccountType()),
		AccountNo:   easygo.NewString(payData.GetAccountNo()),
		AccountName: easygo.NewString(payData.GetAccountName()),
		AccountProp: easygo.NewString("201"),
		PayType:     easygo.NewInt32(for_game.PAY_TYPE_BANKCARD),
		OrderDate:   easygo.NewString(for_game.GetCurTimeString2()),
	}
	obj := for_game.CreateRedisOrder(order)
	return obj.GetOrderId()
}

//4.1请求签约短信码:fastPay.agreement.signSms
func (self *WebHuiJuPay) ReqSignSMSApi(reqMsg *client_hall.BankMessage) (easygo.KWAT, *base.Fail) {
	order := &share_message.Order{}
	id := for_game.RedisCreateOrderNo(order.GetChangeType(), order.GetSourceType())
	data := make(easygo.KWAT)
	data.Add("mch_order_no", id)
	data.Add("order_amount", "0.01")
	data.Add("mch_req_time", for_game.GetCurTimeString2())
	data.Add("payer_name", for_game.HJAesEncrypt(reqMsg.GetUserName(), self.PassWord))
	data.Add("id_type", reqMsg.GetIdType())
	data.Add("id_no", for_game.HJAesEncrypt(reqMsg.GetIdNo(), self.PassWord))
	data.Add("bank_card_no", for_game.HJAesEncrypt(reqMsg.GetBankCardNo(), self.PassWord))
	data.Add("mobile_no", for_game.HJAesEncrypt(reqMsg.GetMobileNo(), self.PassWord))
	if reqMsg.ExpireDate != nil {
		data.Add("expire_date", for_game.HJAesEncrypt(reqMsg.GetExpireDate(), self.PassWord))
	}
	if reqMsg.ExpireDate != nil {
		data.Add("expire_date", for_game.HJAesEncrypt(reqMsg.GetCvv(), self.PassWord))
	}
	js, err := json.Marshal(data)
	easygo.PanicError(err)
	result := self.HttpsReq(self.Url, "fastPay.agreement.signSms", string(js))
	logs.Info("ReqSignSMSApi 请求返回:", string(result))
	return self.CheckData(result)
}

//4.2请求短信绑卡签约:
func (self *WebHuiJuPay) ReqSMSSignApi(reqMsg *client_hall.BankMessage) (easygo.KWAT, *base.Fail) {
	data := make(easygo.KWAT)
	data.Add("mch_order_no", reqMsg.GetOrderNo())
	data.Add("sms_code", reqMsg.GetMsgCode())
	js, err := json.Marshal(data)
	easygo.PanicError(err)
	result := self.HttpsReq(self.Url, "fastPay.agreement.smsSign", string(js))
	logs.Info("ReqSMSSignApi 请求返回:", string(result))
	return self.CheckData(result)
}

//4.3无短信支付
func (self *WebHuiJuPay) ReqFastPayApi(reqMsg *share_message.PayOrderInfo, player *Player) (easygo.KWAT, *base.Fail) {
	//生成订单
	orderId := self.CreateOrder(reqMsg, player)
	amount := fmt.Sprintf("%.2f", easygo.AtoFloat64(reqMsg.GetAmount()))
	data := make(easygo.KWAT)
	data.Add("mch_order_no", orderId)
	data.Add("order_amount", amount)
	data.Add("mch_req_time", for_game.GetCurTimeString2())
	data.Add("order_desc", string(reqMsg.GetProduceName()))
	data.Add("callback_url", self.CallBackUrl)
	data.Add("bank_card_no", for_game.HJAesEncrypt(reqMsg.GetPayBankNo(), self.PassWord))
	js, err := json.Marshal(data)
	logs.Info("js:", string(js))
	easygo.PanicError(err)
	result := self.HttpsReq(self.Url, "fastPay.agreement.pay", string(js))
	logs.Info("ReqFastPayApi 请求返回:", string(result))
	return self.CheckData(result)
}

//4.4请求支付短信码
func (self *WebHuiJuPay) ReqPaySMSApi(reqMsg *share_message.PayOrderInfo, player *Player) (easygo.KWAT, *base.Fail) {
	//生成订单
	orderId := self.CreateOrder(reqMsg, player)
	amount := fmt.Sprintf("%.2f", easygo.AtoFloat64(reqMsg.GetAmount()))
	data := make(easygo.KWAT)
	data.Add("mch_order_no", orderId)
	data.Add("order_amount", amount)
	data.Add("mch_req_time", for_game.GetCurTimeString2())
	data.Add("order_desc", string(reqMsg.GetProduceName()))
	data.Add("callback_url", self.CallBackUrl)
	data.Add("bank_card_no", for_game.HJAesEncrypt(reqMsg.GetPayBankNo(), self.PassWord))
	js, err := json.Marshal(data)
	logs.Info("js:", string(js))
	easygo.PanicError(err)
	result := self.HttpsReq(self.Url, "fastPay.agreement.paySms", string(js))
	logs.Info("ReqPaySMSApi 请求返回:", string(result))
	return self.CheckData(result)
}

//4.5请求短信码支付
func (self *WebHuiJuPay) ReqSMSPayApi(reqMsg *client_hall.BankPaySMS) (easygo.KWAT, *base.Fail) {
	data := make(easygo.KWAT)
	data.Add("mch_order_no", reqMsg.GetOrderNo())
	data.Add("mch_req_time", for_game.GetCurTimeString2())
	data.Add("sms_code", reqMsg.GetSMS())
	js, err := json.Marshal(data)
	easygo.PanicError(err)
	result := self.HttpsReq(self.Url, "fastPay.agreement.smsPay", string(js))
	logs.Info("ReqSMSPayApi 请求返回:", string(result))
	return self.CheckData(result)
}

//4.6银行卡解约
func (self *WebHuiJuPay) ReqUnSignBankApi(bankId string) (easygo.KWAT, *base.Fail) {
	data := make(easygo.KWAT)
	order := &share_message.Order{}
	id := for_game.RedisCreateOrderNo(order.GetChangeType(), order.GetSourceType())
	data.Add("mch_order_no", id)
	data.Add("mch_req_time", for_game.GetCurTimeString2())
	data.Add("bank_card_no", for_game.HJAesEncrypt(bankId, self.PassWord))
	js, err := json.Marshal(data)
	easygo.PanicError(err)
	result := self.HttpsReq(self.Url, "fastPay.agreement.unSign", string(js))
	logs.Info("ReqSMSPayApi 请求返回:", string(result))
	return self.CheckData(result)
}

//4.8订单查询接口
func (self *WebHuiJuPay) ReqCheckPayOrder(orderId string, delay time.Duration) {
	delay += 5 * time.Second
	if delay > 86400*time.Second {
		logs.Info("订单超过一天无结果，转人工:", orderId)
		return
	}
	data := make(easygo.KWAT)
	data.Add("mch_order_no", orderId)
	data.Add("org_mch_req_time", for_game.GetCurTimeString2())
	js, err := json.Marshal(data)
	easygo.PanicError(err)
	result := self.HttpsReq("https://api.joinpay.com/query", "fastPay.query", string(js))
	logs.Info("ReqSMSPayApi 请求返回:", string(result))
	backData, _ := self.CheckData(result)
	self.DealBankPayResult(backData, delay)

}

//处理支付结果
func (self *WebHuiJuPay) DealBankPayResult(data easygo.KWAT, delay ...time.Duration) {
	t := append(delay, 0)[0]
	orderSt := data.GetString("order_status")
	orderId := data.GetString("mch_order_no")
	order := for_game.GetRedisOrderObj(orderId)
	res := false
	if order == nil {
		logs.Info("无效的支付订单")
		return
	}
	if order.GetStatus() == for_game.ORDER_ST_FINISH {
		logs.Info("订单处理已完成")
		return
	}
	if orderSt == "P1000" {
		//交易成功
		order.SetPayStatus(for_game.PAY_ST_FINISH)
	} else if orderSt == "P3000" {
		//处理中，间隔查询订单
		f := func() {
			PWebHuiJuPay.ReqCheckPayOrder(data.GetString("mch_order_no"), time.Second*5)
		}
		easygo.AfterFunc(t, f)
		return
	} else {
		//交易失败，订单取消，订单已关闭
		order.SetPayStatus(for_game.PAY_ST_CANCEL)
		order.SetStatus(for_game.ORDER_ST_CANCEL)
	}
	//支付完成的订单发货
	if order.GetPayStatus() == for_game.PAY_ST_FINISH {
		RechargeGoldToPlayer(order.GetPlayerId(), order.GetOrderId())
		res = true
	}
	//如果充值，通知前端充值结果
	ep := ClientEpMgr.LoadEndpointByPid(order.GetPlayerId())
	if order.GetSourceType() == for_game.GOLD_TYPE_CASH_IN {
		msg := &share_message.RechargeFinish{
			Amount:        easygo.NewInt64(order.GetAmount()),
			TradeNo:       easygo.NewString(order.GetOrderId()),
			PayFinishTime: easygo.NewInt64(order.GetOverTime()),
			Result:        easygo.NewBool(res),
		}
		if ep != nil {
			ep.RpcRechargeMoneyFinish(msg)
			logs.Info("通知前端支付结果:", msg)
		} else {
			if PlayerOnlineMgr.CheckPlayerIsOnLine(order.GetPlayerId()) {
				SendMsgToHallClientNew(order.GetPlayerId(), "RpcRechargeMoneyFinish", msg)
			}
		}
	}
}

//单笔代付接口
func (self *WebHuiJuPay) ReqSinglePay(order *share_message.Order) *base.Fail {
	money := float64(order.GetAmount()+order.GetTax()+order.GetRealTax()) / 100.0
	amount := fmt.Sprintf("%.2f", money)
	data := make(easygo.KWAT)
	data.Add("userNo", self.MerchantNo)             //商户号
	data.Add("productCode", "BANK_PAY_DAILY_ORDER") //产品类型:普通代付,朝夕付,任意付,组合付
	data.Add("requestTime", order.GetOrderDate())   //
	data.Add("merchantOrderNo", order.GetOrderId())
	data.Add("receiverAccountNoEnc", order.GetAccountNo())
	data.Add("receiverNameEnc", order.GetAccountName())
	data.Add("receiverAccountType", "201") //204对公
	data.Add("paidAmount", amount)
	data.Add("currency", "201")
	data.Add("isChecked", "202") //201审核，202不审核
	data.Add("paidDesc", "代付说明")
	data.Add("paidUse", "204") //用途:工资资金-201，活动经费-202，养老金-203，货款-204，劳务费-205，保险理财-206，其他-请联系客服
	data.Add("callbackUrl", "")
	//data.Add("hmac", "签名")
	result := self.HttpsReqDF(self.DaiFuUrl, "", data)
	logs.Info("ReqSinglePay 请求返回:", string(result))
	backData, errMsg := self.CheckDataDF(result)
	logs.Info("返回数:", backData)
	if errMsg.GetCode() != "2001" {
		//受理失败
		msg := backData.GetString("errorDesc")
		return easygo.NewFailMsg(msg)
	}
	//启动定时查询订单:5秒后查询
	fun := func() {
		self.CheckOrder(backData.GetString("merchantOrderNo"), time.Second*5)
	}
	easygo.AfterFunc(time.Second*5, fun)
	return easygo.NewFailMsg("提现下单成功", for_game.FAIL_MSG_CODE_SUCCESS)
}

//查询订单:2201158845919906581141
func (self *WebHuiJuPay) CheckOrder(orderId string, t time.Duration) {
	t += 10 * time.Second
	if t > 86400*time.Second {
		logs.Info("订单超过一天无结果，转人工:", orderId)
		return
	}
	data := make(easygo.KWAT)
	data.Add("userNo", self.MerchantNo)  //商户号
	data.Add("merchantOrderNo", orderId) //订单id
	result := self.HttpsReqDF("https://www.joinpay.com/payment/pay/singlePayQuery", "", data)
	logs.Info("CheckOrder 请求返回:", string(result))
	backData, errMsg := self.CheckDataDF(result)
	logs.Info("返回数:", backData, errMsg.GetCode(), errMsg.GetReason())
	if errMsg.GetCode() != "2001" {
		//定时再次请求
		f := func() {
			self.CheckOrder(orderId, t)
		}
		easygo.AfterFunc(t, f)
	}
	//请求成功的数据
	status := backData.GetString("status")
	order := for_game.GetRedisOrderObj(orderId)
	if order == nil {
		return
	}
	switch status {
	case "205": //交易成功
		logs.Info("代付完成")
		order.SetPayStatus(for_game.PAY_ST_FINISH)
		order.SetStatus(for_game.ORDER_ST_FINISH)
		order.SetExternalNo(backData.GetString("platformSerialNo"))
		order.SetExtendValue("交易完成")
		HandleAfterRecharge(orderId)
	case "204", "208", "214": //交易失败、订单已取消、订单不存在
		order.SetChanneltype(for_game.CHANNEL_MAN_MAKE)
		order.SetExtendValue(backData.GetString("errorDesc"))
		order.SetExternalNo(backData.GetString("platformSerialNo"))
	default:
		//定时再次请求
		f := func() {
			self.CheckOrder(orderId, t)
		}
		easygo.AfterFunc(t, f)
	}
}

//发起支付请求
func (self *WebHuiJuPay) HttpsReq(url, method string, data string) []byte {
	//生成支付日志id
	WebServreMgr.NextJobId()
	WebServreMgr.LogPay("绑卡/充值请求:" + url + ",method:" + method)
	q := easygo.KWAT{}
	q.Add("method", method)
	q.Add("version", "1.0")
	q.Add("data", data)
	q.Add("rand_str", for_game.RandString(32))
	q.Add("sign_type", "2")
	q.Add("mch_no", self.MerchantNo)
	q.Add("sec_key", self.PubEncode(for_game.HJPublicKey, self.PassWord)) //不参与签名
	srcStr := self.GetSourceStr(q, true)
	logs.Info("加密源串:", srcStr)
	sign := self.SignMsg(srcStr)
	q.Add("sign", sign)
	body, err := json.Marshal(q)
	easygo.PanicError(err)
	logs.Info("发送body:", string(body))
	WebServreMgr.LogPay("绑卡/充值body:" + string(body))
	client := &http.Client{}
	// build a new request, but not doing the POST yet
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	easygo.PanicError(err)
	// set the Header here
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	// now POST it
	resp, err := client.Do(req)
	defer resp.Body.Close()
	easygo.PanicError(err)
	result, err := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err)
	WebServreMgr.LogPay("绑卡/充值返回:" + string(result))
	return result
}
func (self *WebHuiJuPay) HttpsReqDF(url, method string, data easygo.KWAT) []byte {
	WebServreMgr.NextJobId()
	WebServreMgr.LogPay("提现请求:" + url + method)
	hmac := self.GetDFHMac(data, DF_HMAC_LIST)
	data.Add("hmac", hmac)
	body, err := json.Marshal(data)
	logs.Info("发送body:", string(body))
	WebServreMgr.LogPay("提现body:" + string(body))
	client := &http.Client{}
	// build a new request, but not doing the POST yet
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	easygo.PanicError(err)
	// set the Header here
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	// now POST it
	resp, err := client.Do(req)
	defer resp.Body.Close()
	easygo.PanicError(err)
	result, err := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err)
	WebServreMgr.LogPay("提现返回:" + string(body))
	return result
}

//fa
var PWebHuiJuPay = NewWebHuiJuPay()
