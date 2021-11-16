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
)

/*
	微信api接口
*/
const WXLOGIN_BUSY = -1       //系统繁忙，此时请开发者稍候再试
const WXLOGIN_SUCCESS = 0     //请求成功
const WXLOGIN_INVALID = 40029 //code 无效
const WXLOGIN_BAN = 45011     //频率限制，每个用户每分钟100次

type WXApiMgr struct {
	LoginUrl string
	//Appid     string
	//Secret    string
	GrantType string
}

func NewWXApiMgr() *WXApiMgr {
	p := &WXApiMgr{}
	p.Init()
	return p
}

//https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code
func (self *WXApiMgr) Init() {
	self.LoginUrl = "https://api.weixin.qq.com/sns/jscode2session"
	//self.Appid = "wx414f06b4024739c8" // 旧主体的appid
	//self.Secret = "645661043e2b05624c3e1bb9bec71c3b" // 旧主体的 Secret
	//self.Appid = "wx905904883fc074d3"                // 新主体的 appid
	//self.Secret = "98844f416ebfdbd23a7238c592e08a30" // 新主体的  Secret
	self.GrantType = "authorization_code"
}

func (self *WXApiMgr) GetWXAPIError(code int32, err string) []byte {
	backMsg := &share_message.RechargeOrderResult{
		Result:  easygo.NewBool(false),
		Code:    easygo.NewInt32(code),
		Message: easygo.NewString(err),
	}
	d, _ := json.Marshal(backMsg)
	return d
}

//登录验证
func (self *WXApiMgr) WXAPILogin(web *WebHttpServer, code string, payData *share_message.PayOrderInfo) ([]byte, *base.Fail) {
	web.LogPay(payData)
	channel := payData.GetPayId() //默认是秒到
	vals := url.Values{}
	appKey := for_game.MakeNewString("WEIXIN_SGAME_APPID", payData.GetApkCode())
	secretKey := for_game.MakeNewString("WEIXIN_SGAME_SECRET", payData.GetApkCode())
	appId := easygo.YamlCfg.GetValueAsString(appKey)
	appSecret := easygo.YamlCfg.GetValueAsString(secretKey)
	vals.Add("appid", appId)
	vals.Add("secret", appSecret)
	vals.Add("js_code", code)
	vals.Add("grant_type", self.GrantType)
	resData, err := self.HttpsReq(vals)
	if err != nil {
		return nil, easygo.NewFailMsg("微信登录请求失败:" + err.Error())
	}
	web.LogPay("WX登录请求返回:" + string(resData))
	msg := &share_message.WXLoginResult{}
	err = json.Unmarshal(resData, msg)
	if err != nil {
		return nil, easygo.NewFailMsg("解析微信登录返回数据失败:" + err.Error())
	}
	if msg.GetErrcode() != WXLOGIN_SUCCESS { // 调用微信登录不成功
		return nil, easygo.NewFailMsg("微信授权登录不成功")
	}

	// 调用微信成功
	msg.Errcode = easygo.NewInt32(WXLOGIN_SUCCESS)
	if channel == for_game.PAY_CHANNEL_MIAODAO {
		return PWebMiaoDaoPay.RechargeOrder(web, payData, msg.GetOpenid())
	} else if channel == for_game.PAY_CHANNEL_TONGLIAN {
		return PWebTongLianPay.RechargeOrder(web, payData, msg.GetOpenid())
	} else if channel == for_game.PAY_CHANNEL_HUICHAO_WX {
		return PWebHuiChaoPay.RechargeOrder(web, payData, msg.GetOpenid())
	} else if channel == for_game.PAY_CHANNEL_TONGTONG_WX {
		return PWebTongTongPay.RechargeOrder(web, payData, msg.GetOpenid())
	} else {
		return nil, easygo.NewFailMsg("error 找不到支付渠道:" + easygo.AnytoA(channel))
	}
}

//发起Get请求
func (self *WXApiMgr) HttpsReq(q url.Values) ([]byte, error) {
	u, _ := url.Parse(self.LoginUrl)
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	result, err1 := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err1 != nil {
		return nil, err1
	}
	return result, nil
}

var PWXApiMgr = NewWXApiMgr()
