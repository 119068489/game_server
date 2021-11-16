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

	"github.com/astaxie/beego/logs"
)

/*
秒到支付模块逻辑
*/
type WebMiaoDaoPay struct {
	AppId       string //小程序appId
	OpenId      string //小程序openId
	MerchantNo  string //商户号
	AgencyCode  string //机构号
	NotifyUrl   string //回调地址
	CallBackUrl string //页面回调地址(前端)
	Key         string //密钥
	Url         string //接口地址
	PayPath     string //下单接口支持：小程序、公众号、服务窗
}

func NewWebMiaoDaoPay() *WebMiaoDaoPay {
	p := &WebMiaoDaoPay{}
	p.Init()
	return p
}

//秘钥 090b8a929133cef439abf2e27fa2215c  机构号86038810商户号15989289979
func (self *WebMiaoDaoPay) Init() {

	self.Url = "http://open.miaodaochina.com"
	self.PayPath = "/unify/mini/pay"
	//self.NotifyUrl = "http://api.chihuoqun.org/miaodao"
	//self.NotifyUrl = "http://47.112.223.121:2601/miaodao"
}

type MDResult struct {
	code    int32
	data    string `json:"totalMoney"`
	message string
	result  bool
}

//发起充值：
//1数据库先生成订单，然后再调起请求
//微信小程序充值
func (self *WebMiaoDaoPay) RechargeOrder(web *WebHttpServer, payData *share_message.PayOrderInfo, openId string) ([]byte, *base.Fail) {
	//参数设置
	self.AppId = easygo.YamlCfg.GetValueAsString("MD_PAY_APPID")
	self.Key = easygo.YamlCfg.GetValueAsString("MD_PAY_KEY")
	self.MerchantNo = easygo.YamlCfg.GetValueAsString("MD_PAY_MERCHANTNO")
	self.AgencyCode = easygo.YamlCfg.GetValueAsString("MD_PAY_AGENCYCODE")
	self.NotifyUrl = easygo.YamlCfg.GetValueAsString("MD_PAY_NOTIFY_URL")
	param := &share_message.RechargeOrder{
		MerchantNo:       easygo.NewString(self.MerchantNo),
		AppId:            easygo.NewString(self.AppId),
		AgencyCode:       easygo.NewString(self.AgencyCode),
		NotifyUrl:        easygo.NewString(self.NotifyUrl),
		OpenId:           easygo.NewString(openId),
		GatewayPayMethod: easygo.NewString("MINIPROGRAM"),
		PiType:           easygo.NewString("WX"),
		TotalAmount:      easygo.NewString(payData.GetAmount()),
		ProductName:      easygo.NewString(payData.GetProduceName()),
	}
	//订单号生成

	orderNo := self.CreateOrder(payData, openId)
	param.OutTradeNo = easygo.NewString(orderNo)

	urlVals := self.MakePostParams(param)
	resData := self.HttpsReq(self.PayPath, urlVals)
	web.LogPay(string(resData))
	return resData, nil
}

//下订单
func (self *WebMiaoDaoPay) CreateOrder(payData *share_message.PayOrderInfo, openId string) string {
	tax := int64(0)
	amount := easygo.AtoFloat64(payData.GetAmount())
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
		CreateIP:    easygo.NewString(""),
		Status:      easygo.NewInt32(for_game.ORDER_ST_WAITTING),
		PayStatus:   easygo.NewInt32(for_game.PAY_ST_WAITTING),
		Note:        easygo.NewString(note),
		Tax:         easygo.NewInt64(-tax),
		Operator:    easygo.NewString("system"),
		PayChannel:  easygo.NewInt32(for_game.PAY_CHANNEL_MIAODAO),
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
func (self *WebMiaoDaoPay) GetSourceStr(values url.Values) string {
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
		if k == "sign" { //sign不参与签名源串
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
	return self.Key + buf.String() + self.Key
}

func (self *WebMiaoDaoPay) MakePostParams(order interface{}) url.Values {
	values := make(url.Values)
	data, err := json.Marshal(order)
	easygo.PanicError(err)
	strData := strings.Trim(string(data), "{}")
	strList := strings.Split(strData, ",")
	for _, v := range strList {
		newList := strings.SplitN(v, ":", 2)
		key := strings.Trim(newList[0], "\"")
		val := strings.Trim(newList[1], "\"")
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
func (self *WebMiaoDaoPay) HttpsReq(method string, params url.Values) []byte {
	logs.Info("发送请求:url=", self.Url+method)
	resp, err := http.PostForm(self.Url+method, params)
	easygo.PanicError(err)
	defer resp.Body.Close()
	result, err1 := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err1)
	return result
}

//var PWebMiaoDaoPay = NewWebMiaoDaoPay()
