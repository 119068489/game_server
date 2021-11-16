package hall

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/alipay"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

/*
	支付宝api接口
*/

type AliApiMgr struct {
	LoginUrl string
	Appid    string
}

func NewAliApiMgr() *AliApiMgr {
	p := &AliApiMgr{}
	p.Init()
	return p
}

//https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code
func (self *AliApiMgr) Init() {
	self.LoginUrl = "https://openapi.alipay.com/gateway.do"
	self.Appid = "2021001171677322"
	//self.LoginUrl = "https://openapi.alipaydev.com/gateway.do" //沙箱
	//self.Appid = "2016102700770128"                            //沙箱
}
func (self *AliApiMgr) GetSourceStr(values url.Values) string {
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
	return buf.String()
}

//签名
func (self *AliApiMgr) MakeSign(str string) string {
	h := crypto.Hash.New(crypto.SHA256) //进行SHA1的散列
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	block, _ := pem.Decode(for_game.AliGamePrivateKey)
	if block == nil {
		logs.Info("block=nil")
		return ""
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	easygo.PanicError(err)
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv.(*rsa.PrivateKey), crypto.SHA256, hashed)
	logs.Info("SignMsg:", base64.StdEncoding.EncodeToString(signature))
	return base64.StdEncoding.EncodeToString(signature)
}

/**
AliAPILogin 登录验证
支付宝请求失败返回: {"error_response":{"code":"40002","msg":"Invalid Arguments","sub_code":"isv.code-invalid","sub_msg":"授权码code无效"},"sign":"Fg9nEKYQGAouawPq0kz88IeGEWtm+Xcdienw8V5rMlPTnNEQ/TCT4u6w6cjxt4MrTv+22qJefuJlQtUSkN6qzNe/6fojrnxn17ySz9wP7zUgrXkGQ2/sTZp8TP9AIZ2sQzpU+/3CUr+dXyfqbAahyaTxnFvzDFS+tyi/uSs5LRD8757+dPOmnRTmna5DNZREGQy5XxLcI5yYath3bB1AsWfyZLidCLkCXFb+lEPVEMBzeP9AyUz8j4Otf0nyUKdzhpwZuuSTH69XURSAcay3izwr7YKj6ZVAnXXgJjFrqjKF6hnsYtw1sLx0mTltklcZp4jERde++Fj2WP75tYIR7g=="}
支付宝请求成功返回: {"alipay_system_oauth_token_response":{"access_token":"authusrB94739c226e9c400196af5bc4740a9X01","alipay_user_id":"20880063008624437006283650118501","expires_in":1296000,"re_expires_in":2592000,"refresh_token":"authusrBff1f0be8cae1436dace07e6e7cd7cX01","user_id":"2088902487861012"},"sign":"W/ZQgCAX2ExObxWYWire/+vfGfRqiA4IIBZC9oUzNLdHF+rCC2aZN1ERqymUomqBVjaFsFSB0S0XzqbTFJNf+JkQpwvYwJr0mCudA2mjolpnD/jIZJy+to/w20Kw/wp6vvdIcZ/pXGULi2s98Ndu1mnmC3phZNfNGgAUHnzG1QVeo2ap/yT9iul04oMsiRHddcAWHKY44SE8Bpj4k6MdqBufuBeC4fbfbPDH/4UfXjN3Dj6FJNAjy5n9+gofXnzECOP0Jvj63I12/ag273d1FJYLVppCuyw+TUGbnEOiGh7ZtNsDWnBmZRPv2uG7DvsiogEEuZYM0Fm3X1fUMLoDPw=="}
*/
func (self *AliApiMgr) GetAliAPIError(code int32, err string) []byte {
	backMsg := &share_message.RechargeOrderResult{
		Result:  easygo.NewBool(false),
		Code:    easygo.NewInt32(code),
		Message: easygo.NewString(err),
	}
	d, _ := json.Marshal(backMsg)
	return d
}
func (self *AliApiMgr) AliAPILogin(web *WebHttpServer, code string, payData *share_message.PayOrderInfo) ([]byte, *base.Fail) {
	//web.LogPay(payData)
	logs.Info("paydata:", payData)
	//channel := payData.GetPayId() //默认是秒到
	if payData.GetPaySence() == 3 {
		logs.Info("扫码支付，不需要登录")
		return PWebTongTongPay.RechargeOrder(web, payData, "")
	}
	vals := url.Values{}

	vals.Set("app_id", self.Appid)
	vals.Set("method", "alipay.system.oauth.token")
	vals.Set("format", "JSON")
	vals.Set("charset", "UTF-8")
	vals.Set("sign_type", "RSA2")
	vals.Set("timestamp", for_game.GetCurTimeString2())
	vals.Set("version", "1.0")
	vals.Add("grant_type", "authorization_code")
	vals.Set("code", code)
	bizContent := "{\"grant_type\":\"authorization_code\",\"code\":" + code + "}"
	vals.Set("biz_content", bizContent)
	sign := self.MakeSign(self.GetSourceStr(vals))
	vals.Set("sign", sign)

	logs.Info("请求参数:", vals)
	resData := self.HttpsReq(vals)
	logs.Info("支付宝支付返回:", string(resData))

	var r *share_message.AliLoginResult
	err := json.Unmarshal(resData, &r)
	if err != nil {
		logs.Error(err)
		return nil, easygo.NewFailMsg("第三方支付返回数据异常:", err.Error())
	}

	if r.ErrorResponse.GetCode() != "" { // 调用支付宝登录不成功
		return nil, easygo.NewFailMsg(r.ErrorResponse.GetMsg())
	}

	channel := payData.GetPayId()
	switch channel { // 支付宝目前只有汇潮
	case for_game.PAY_CHANNEL_HUICHAO_ZFB: // 支付宝目前只有汇潮
		return PWebHuiChaoPay.RechargeOrder(web, payData, r.AlipaySystemOauthTokenResponse.GetUserId())
	case for_game.PAY_CHANNEL_TONGTONG_ZFB:
		return PWebTongTongPay.RechargeOrder(web, payData, r.AlipaySystemOauthTokenResponse.GetUserId())
	default:
		return nil, easygo.NewFailMsg("找不到支付渠道:" + easygo.AnytoA(channel))
	}

}

func (self *AliApiMgr) AliGetUserInfo(web *WebHttpServer, code string, payData *share_message.PayOrderInfo) []byte {
	//web.LogPay(payData)
	logs.Info("paydata:", payData)
	//channel := payData.GetPayId() //默认是秒到
	vals := url.Values{}
	vals.Add("app_id", self.Appid)
	vals.Add("method", "alipay.user.info.share")
	vals.Add("format", "JSON")
	vals.Add("charset", "UTF-8")
	vals.Add("sign_type", "RSA2")
	vals.Add("timestamp", for_game.GetCurTimeString2())
	vals.Add("version", "1.0")
	//vals.Add("grant_type", "authorization_code")
	vals.Add("app_auth_token", code)
	bizContent := "{\"grant_type\":\"authorization_code\",\"code\":" + code + "}"
	vals.Add("biz_content", bizContent)
	sign := self.MakeSign(self.GetSourceStr(vals))
	vals.Add("sign", sign)
	logs.Info("请求参数:", vals)
	resData := self.HttpsReq(vals)
	logs.Info("支付宝支付返回:", string(resData))
	//web.LogPay("支付宝登录请求返回:" + string(resData))
	//msg := &share_message.WXLoginResult{}
	//err := json.Unmarshal(resData, msg)
	//easygo.PanicError(err)
	//if msg.GetErrcode() == WXLOGIN_SUCCESS {
	//	msg.Errcode = easygo.NewInt32(WXLOGIN_SUCCESS)
	//	if channel == for_game.PAY_CHANNEL_MIAODAO {
	//		return PWebMiaoDaoPay.RechargeOrder(web, payData, msg.GetOpenid())
	//	} else if channel == for_game.PAY_CHANNEL_TONGLIAN {
	//		return PWebTongLianPay.RechargeOrder(web, payData, msg.GetOpenid())
	//
	//	} else if channel == for_game.PAY_CHANNEL_HUICHAO_WX {
	//		return PWebHuiChaoPay.RechargeOrder(web, payData, msg.GetOpenid())
	//	} else {
	//		panic("error 找不到支付渠道:" + easygo.AnytoA(channel))
	//	}
	//} else {
	//	backMsg := &share_message.RechargeOrderResult{
	//		Result:  easygo.NewBool(false),
	//		Code:    easygo.NewInt32(msg.GetErrcode()),
	//		Message: easygo.NewString(msg.GetErrmsg()),
	//	}
	//	d, _ := json.Marshal(backMsg)
	//	return d
	//}
	return resData
}

//发起Get请求
func (self *AliApiMgr) HttpsReq(q url.Values) []byte {
	u, _ := url.Parse(self.LoginUrl)
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	easygo.PanicError(err)
	result, err1 := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	easygo.PanicError(err1)
	return result
}

var PAliApiMgr = NewAliApiMgr()

/**
TradeQuery 交易订单查询
tradeNo 前端取支付宝返回的字段 PayStr
isProd 正式环境true,沙箱 false
*/
func (self *AliApiMgr) TradeQuery(tradeNo string, t time.Duration) string {
	// 前期先简单用下面的初始化client,后期其他交易如果需要使用这个第三方库时统一初始化即可
	var isProd bool
	if for_game.IS_FORMAL_SERVER { // 正式服
		isProd = true
	}
	client := alipay.NewClient(self.Appid, for_game.AliPrivateKeyNoPreEnd, isProd)
	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("trade_no", tradeNo)
	client.PrivateKeyType = alipay.PKCS8
	query, err := client.TradeQuery(bm)
	if err != nil {
		// {"code":"20000","msg":"Service Currently Unavailable","sub_code":"aop.ACQ.SYSTEM_ERROR","sub_msg":"系统异常"}
		var m map[string]string
		err1 := json.Unmarshal([]byte(err.Error()), &m)
		easygo.PanicError(err1)
		return m["sub_msg"]
	}
	logs.Info("------------>%+v", query)
	// todo 待处理返回数据.

	return "" // todo
}
