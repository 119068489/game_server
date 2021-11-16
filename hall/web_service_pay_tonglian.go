package hall

import (
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
	"time"

	"github.com/astaxie/beego/logs"
)

/*
通联支付模块逻辑
*/
type WebTongLianPay struct {
	AppId         string //小程序appId
	OpenId        string //小程序openId
	MerchantNo    string //商户号
	MerchantAppId string //商户appid
	NotifyUrl     string //回调地址
	Key           string //密钥
	Url           string //接口地址
}

//订单查询
type CheckOrder struct {
	Cusid     string `json:"cusid"`     //商户号
	AppId     string `json:"appid"`     //应用ID
	Version   string `json:"version"`   //版本号 ，默认填11
	Reqsn     string `json:"reqsn"`     //商户订单号
	Trxid     string `json:"trxid"`     //平台交易流水
	Randomstr string `json:"randomstr"` //随机字符串
	Sign      string `json:"sign"`      //签名
}
type CheckOrderData struct {
	Retcode   string `json:"retcode"`   //返回码
	Retmsg    string `json:"retmsg"`    //返回说明码
	Cusid     string `json:"cusid"`     //商户号
	AppId     string `json:"appid"`     //应用ID
	Trxid     string `json:"trxid"`     //平台交易流水
	Chnltrxid string `json:"chnltrxid"` //支付渠道单号
	Reqsn     string `json:"reqsn"`     //商户订单号
	Trxcode   string `json:"trxcode"`   //交易类型
	Trxamt    string `json:"trxamt"`    //交易金额
	Trxstatus string `json:"trxstatus"` //交易状态
	Acct      string `json:"acct"`      //支付平台用户标识
	Fintime   string `json:"fintime"`   //交易完成时间
	Randomstr string `json:"randomstr"` //随机字符串
	Errmsg    string `json:"errmsg"`    //错误原因
	Cmid      string `json:"cmid"`      //渠道子商户号
	Chnlid    string `json:"chnlid"`    //渠道号
	Initamt   string `json:"initamt"`   //原交易金额
	Fee       string `json:"fee"`       //手续费
	Sign      string `json:"sign"`      //签名
}

//返回前端支付信息，payinfo前端拉起
type PayInfo struct {
	PayType        int32  `json:"payType"`
	PreparePayInfo string `json:"preparePayInfo"` //支付payinfo
	OrderId        string `json:"orderId"`        //订单id
}

type ParseMDData struct {
	Data    PayInfo `json:"data"`
	Result  bool    `json:"result"`
	Code    int32   `json:"code"`
	Message string  `json:"message"`
}

//撤销订单结构
type CloseOrderObj struct {
	Cusid     string `json:"cusid"`     //商户号
	AppId     string `json:"appid"`     //应用ID
	Version   string `json:"version"`   //版本号 ，默认填11
	Oldreqsn  string `json:"oldreqsn"`  //原交易单号
	Oldtrxid  string `json:"oldtrxid"`  //原交易流水
	Randomstr string `json:"randomstr"` //随机字符串
	SignType  string `json:"signtype"`  //签名方式
	Sign      string `json:"sign"`      //签名
}

//撤销订单返回数据
type CloseOrderResult struct {
	Retcode   string `json:"retcode"`   //返回码
	Retmsg    string `json:"retmsg"`    //返回说明码
	Cusid     string `json:"cusid"`     //商户号
	AppId     string `json:"appid"`     //应用ID
	Trxstatus string `json:"trxstatus"` //交易状态
	Errmsg    string `json:"errmsg"`    //错误原因
	Randomstr string `json:"randomstr"` //随机字符串
	Sign      string `json:"sign"`      //签名
}

func NewWebTongLianPay() *WebTongLianPay {
	p := &WebTongLianPay{}
	p.Init()
	return p
}

//秘钥 090b8a929133cef439abf2e27fa2215c  机构号86038810商户号15989289979
func (self *WebTongLianPay) Init() {
	self.AppId = "wx414f06b4024739c8"
	self.Key = "64336113WEIwei"                             //md5key：正式:64336113WEIwei  allinpay888
	self.MerchantNo = "56058104816WR4N"                     //正式:"56058104816WR4N"
	self.MerchantAppId = "00184603"                         //正式:"00184603"
	self.Url = "https://vsp.allinpay.com/apiweb/unitorder/" //正式
	self.NotifyUrl = "http://103.254.208.205:2601/tonglian"
}

//发起充值：
//1数据库先生成订单，然后再调起请求
//微信小程序充值
func (self *WebTongLianPay) RechargeOrder(web *WebHttpServer, payData *share_message.PayOrderInfo, openId string) ([]byte, *base.Fail) {
	//参数设置
	amount := easygo.AtoFloat64(payData.GetAmount())
	money := int64(amount * 100)
	self.MerchantNo = easygo.YamlCfg.GetValueAsString("TL_PAY_MERCHANTNO")
	self.MerchantAppId = easygo.YamlCfg.GetValueAsString("TL_PAY_MERCHANTAPPID")
	self.NotifyUrl = easygo.YamlCfg.GetValueAsString("TL_PAY_NOTIFY_URL")
	self.AppId = easygo.YamlCfg.GetValueAsString("TL_PAY_APPID")
	self.Key = easygo.YamlCfg.GetValueAsString("TL_PAY_KEY")
	param := &share_message.RechargeTLOrder{
		Cusid:     easygo.NewString(self.MerchantNo),    //商户号
		Appid:     easygo.NewString(self.MerchantAppId), //商户appid
		NotifyUrl: easygo.NewString(self.NotifyUrl),
		Body:      easygo.NewString(payData.GetProduceName()),
		Randomstr: easygo.NewString(for_game.RandString(8)),
		Paytype:   easygo.NewString("W06"),
		Trxamt:    easygo.NewString(easygo.AnytoA(money)),
		Validtime: easygo.NewString("5"), //订单存在时间
		Acct:      easygo.NewString(openId),
		SubAppid:  easygo.NewString(self.AppId),
		Signtype:  easygo.NewString("MD5"),
	}
	//订单号生成
	orderNo := self.CreateOrder(payData, openId)
	logs.Info("微信支付,内部订单号为: ------->", orderNo)
	param.Reqsn = easygo.NewString(orderNo)

	urlVals := self.MakePostParams(param)
	resData := self.HttpsReq("pay", urlVals)
	web.LogPay(string(resData))
	//通知前端结果前端
	msg := &share_message.RechargeTLOrderResult{}
	err := json.Unmarshal(resData, msg)
	if err != nil {
		return nil, easygo.NewFailMsg("第三方支付失败:", err.Error())
	}
	if msg.GetRetcode() != "SUCCESS" {
		return nil, easygo.NewFailMsg("第三方支付失败:" + msg.GetRetcode())
	}
	order := for_game.GetRedisOrderObj(msg.GetReqsn())
	if order == nil {
		return nil, easygo.NewFailMsg("无效的支付订单:" + msg.GetReqsn())
	}
	order.SetExternalNo(msg.GetTrxid())
	pData := &ParseMDData{
		Data: PayInfo{
			PayType:        payData.GetPayType(),
			PreparePayInfo: msg.GetPayinfo(),
			OrderId:        msg.GetReqsn(),
		},
		Result:  true,
		Code:    200,
		Message: msg.GetRetmsg(),
	}
	d, err := json.Marshal(pData)
	if err != nil {
		return nil, easygo.NewFailMsg("json.Marshal(pData):" + err.Error())
	}
	return d, nil
}

//发起查询支付订单结果
func (self *WebTongLianPay) CheckOrderResult(web *WebHttpServer, orderId string) string {
	backStr := ""
	order := for_game.GetRedisOrderObj(orderId)
	if order == nil {
		backStr = "无效的订单id:" + orderId
		return backStr
	}
	checkOrder := &CheckOrder{
		Cusid:     self.MerchantNo,
		AppId:     self.MerchantAppId,
		Version:   "11",
		Reqsn:     orderId,
		Trxid:     order.GetExternalNo(),
		Randomstr: for_game.RandString(8),
	}
	urlVals := self.MakePostParams(checkOrder)
	resData := self.HttpsReq("query", urlVals)
	web.LogPay("查询交易结果:" + string(resData))
	backData := &CheckOrderData{}
	err := json.Unmarshal(resData, backData)
	logs.Info("backData:", backData)
	easygo.PanicError(err)
	if backData.Retcode == "SUCCESS" {
		if backData.Trxstatus == "2008" || backData.Trxstatus == "2000" {
			return self.CloseOrder(web, order)
		}
	} else {
		logs.Info("查询订单失败:", backData.Retmsg)
		backStr = backData.Retmsg
	}
	return backStr
}

//支付下订单
func (self *WebTongLianPay) CreateOrder(payData *share_message.PayOrderInfo, openId string) string {
	amount := easygo.AtoFloat64(payData.GetAmount())
	tax := int64(0)
	changeGold := int64(amount * 100)
	player := GetPlayerObj(payData.GetPlayerId())
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
		PayChannel:  easygo.NewInt32(for_game.PAY_CHANNEL_TONGLIAN),
		PayType:     easygo.NewInt32(payData.GetPayType()),
		PayWay:      easygo.NewInt32(payData.GetPayWay()),
		PayTargetId: easygo.NewInt64(payData.GetPayTargetId()),
		PayOpenId:   easygo.NewString(openId),
		TotalCount:  easygo.NewInt32(payData.GetTotalCount()),
		Content:     easygo.NewString(payData.GetContent()),
		ExtendValue: easygo.NewString(payData.GetExtendValue()),
	}
	obj := for_game.CreateRedisOrder(order)
	return obj.GetOrderId()
}

//获取加密源串
//签名值组装:参数按照参数名的 ASCII 码顺序从小到大排序,然后使用&字符拼接成字符串
//对拼接的参数字符串再次在头部和尾部拼接分配给商户的 key
func (self *WebTongLianPay) GetSourceStr(values url.Values) string {
	if values == nil {
		return ""
	}
	//加上MD5key
	values.Add("key", self.Key)
	var buf strings.Builder
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if k == "sign" { //sign不参与签名源串
			continue
		}
		vs := values[k]
		for _, v := range vs {
			if v == "" {
				continue
			}
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(v)
		}
	}
	//去掉MD5key
	values.Del("key")
	return buf.String()
}

func (self *WebTongLianPay) MakePostParams(order interface{}) url.Values {
	values := make(url.Values)
	data, err := json.Marshal(order)
	easygo.PanicError(err)
	strData := strings.Trim(string(data), "{}")
	strList := strings.Split(strData, ",")
	for _, v := range strList {
		newList := strings.SplitN(v, ":", 2)
		key := strings.Trim(newList[0], "\"")
		val := strings.Trim(newList[1], "\"")
		if key == "sign" { //sign不参与签名
			continue
		}
		values.Add(key, val)
	}
	s := self.GetSourceStr(values)
	logs.Info("签名源串:", s)
	sign := strings.ToUpper(for_game.Md5(s))
	logs.Info("签名:", sign)
	values.Add("sign", sign)
	return values
}

// 接口POST请求:
func (self *WebTongLianPay) HttpsReq(method string, params url.Values) []byte {
	logs.Info("发送请求:url=", self.Url+method, params)
	resp, err := http.PostForm(self.Url+method, params)
	easygo.PanicError(err)
	defer resp.Body.Close()
	result, err1 := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err1)
	return result
}

//关闭交易
func (self *WebTongLianPay) CloseOrder(web *WebHttpServer, order *for_game.RedisOrderObj) string {
	backStr := ""
	checkOrder := &CloseOrderObj{
		Cusid:     self.MerchantNo,
		AppId:     self.MerchantAppId,
		Version:   "11",
		Oldreqsn:  order.GetOrderId(),
		Oldtrxid:  order.GetExternalNo(),
		Randomstr: for_game.RandString(8),
		SignType:  "MD5",
	}
	urlVals := self.MakePostParams(checkOrder)
	resData := self.HttpsReq("close", urlVals)
	web.LogPay("关闭交易结果:" + string(resData))
	backData := &CloseOrderResult{}
	err := json.Unmarshal(resData, backData)
	easygo.PanicError(err)
	logs.Info("backData:", backData)
	if backData.Retcode == "SUCCESS" {
		if backData.Trxstatus == "0000" {
			//正在处理中尚未完成
			order.SetPayStatus(for_game.PAY_ST_CANCEL)
			order.SetStatus(for_game.ORDER_ST_CANCEL)
			order.SetOverTime(time.Now().Unix())
			//通知前端交易取消
			RechargeMoneyFinish(order)
			logs.Info("订单取消处理完成")
		}
	} else {
		backStr = backData.Retmsg
	}

	return backStr
}

var PWebTongLianPay = NewWebTongLianPay()

// 订单付款后,前端通知后端接口支付结果,后端通知前端订单结果
func RechargeMoneyFinish(order *for_game.RedisOrderObj) {
	ep := ClientEpMgr.LoadEndpointByPid(order.GetPlayerId())
	if order.GetSourceType() == for_game.GOLD_TYPE_CASH_IN {
		msg := &share_message.RechargeFinish{
			Amount:        easygo.NewInt64(order.GetAmount()),
			TradeNo:       easygo.NewString(order.GetOrderId()),
			PayFinishTime: easygo.NewInt64(order.GetOverTime()),
			Result:        easygo.NewBool(false),
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
