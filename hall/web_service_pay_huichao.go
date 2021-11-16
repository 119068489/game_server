package hall

import (
	"crypto"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

/*
	汇潮支付服务入口
*/

type WebHuiChaoPay struct {
	MerchantNo    string   // 商户号
	Url           string   // 请求支付地址  http post请求
	CallBackUrl   string   // 支付回调地址
	ReturnURL     string   // 前端跳回页面
	CallBackUrlDF string   // 代付回调地址
	HeadFormat    string   // xml头部
	WXCompanyNo   string   // 微信支付商户标签（入驻获取）
	ZFBCompanyNo  string   // 支付宝商户标签（入驻获取）
	WXAppId       string   // 微信小程序appid
	BankList      []string //支持的银行卡列表
}
type IMerchantInfoObj interface {
}
type ZFBMerchantInfoObj struct {
	IMerchantInfoObj
	MerName      string // 商户名称
	ShortName    string // 商户简称，会显示在订单上
	ContactName  string // 联系人，支付宝必填
	ServicePhone string // 客服电话
	Business     string // 行业类别：支付宝参考11.1，微信参考：11.2
	Mcc          string // mcc码 参考11.1 支付宝必须
	ContactTag   string // 商户联系人业务 支付宝必须
	ContactType  string // 联系人类型 支付宝必须
	City         string // 城市 支付宝必须
	District     string // 区县 支付宝必须
	Address      string // 地址 支付宝必须
	Province     string // 省份 支付宝必须
	CardNo       string // 对公结算银行卡号 支付宝必须
	CardName     string // 对公结算银行卡持卡人姓名  支付宝必须
	//ContactPhone  string // 联系人电话 花呗必传
	//ContactMobile string // 联系人手机号 花呗必传
	//BusinessLicense     string // 商户证件编号 花呗必传
	//BusinessLicenseType string // 商户证件类型 花呗必传
	//IdCardNo            string // 联系人身份证 花呗必传
}
type WXMerchantInfoObj struct {
	IMerchantInfoObj
	MerName      string // 商户名称
	ShortName    string // 商户简称，会显示在订单上
	ContactName  string // 联系人，支付宝必填
	ServicePhone string // 客服电话
	Business     string // 行业类别：支付宝参考11.1，微信参考：11.2
	//=====================微信参数=======================
	SubAppID string // 微信商户appid 小程序必传
	PayPath  string // 微信授权目录路径 微信小程序支付方式均必传(路径末尾一定要加符号/)
}

//入驻请求
type ScanMerchantInRequest struct {
	MerNo        string // 商户号
	Version      string // 版本 1.0
	PayType      string // 签名
	ChannelNo    string // 渠道号:微信渠道号/支付宝 PID
	RandomStr    string // 随机字符串
	SignInfo     string // 支付类型:微信线下支付 WXZF,WXZF_ONLINE,支付宝线下支付:ZFBZF
	NotifyUrl    string // 异步通知地址商户编号
	MerchantInfo IMerchantInfoObj
	CompanyNo    string //
}

func NewWebHuiChaoPay() *WebHuiChaoPay {
	p := &WebHuiChaoPay{}
	p.Init()
	return p
}
func (self *WebHuiChaoPay) Init() {
	self.Url = easygo.YamlCfg.GetValueAsString("HUICHAO_URL")
	self.MerchantNo = easygo.YamlCfg.GetValueAsString("HUICHAO_MERCHANTNO")
	self.CallBackUrl = easygo.YamlCfg.GetValueAsString("HUICHAO_PAY_CALLBACK_URL") // 测试www.lemonchat.cn/testhuichaopay
	self.ReturnURL = ""
	self.CallBackUrlDF = easygo.YamlCfg.GetValueAsString("HUICHAO_PAY_DF_CALLBACK_URL")
	self.HeadFormat = "<?xml version=\"1.0\" encoding=\"utf-8\"?>"
	self.WXCompanyNo = easygo.YamlCfg.GetValueAsString("HUICHAO_WXCOMPANYNO")
	self.ZFBCompanyNo = easygo.YamlCfg.GetValueAsString("HUICHAO_ZFBCOMPANYNO")
	//self.WXAppId = "wx414f06b4024739c8" // 旧主体的 appid
	self.WXAppId = easygo.YamlCfg.GetValueAsString("HUICHAO_WXAPPID") // 新主体的 appid
	//汇潮支持的银行
	self.BankList = []string{"ICBC", "ABC", "CCB", "BOC", "BOCOM", "BOS", "BCCB", "CEB", "CIB", "CMB", "CMBC", "CNCB", "GDB", "HXB", "PAB", "PSBC", "SPDB"}
}

//获取支持的银行卡列表
func (self *WebHuiChaoPay) GetSupportBankList() []*client_hall.BankData {
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

//签名消息体，并返回base64数据
func (self *WebHuiChaoPay) SignMsg(strData string) string {
	block, _ := pem.Decode([]byte(for_game.HCHAOPrivateKey))
	if block == nil {
		logs.Info("block=nil")
		return ""
	}
	private, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	easygo.PanicError(err)

	h := crypto.Hash.New(crypto.SHA1) // 进行SHA1的散列
	h.Write([]byte(strData))
	hashed := h.Sum(nil)
	// 进行rsa加密签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, private.(*rsa.PrivateKey), crypto.SHA1, hashed)

	sign := base64.StdEncoding.EncodeToString(signature)
	logs.Info("sign:", sign)
	return sign
}

// 代付签名消息体，并返回base64数据
func (self *WebHuiChaoPay) DFSignMsg(strData string) string {
	block, _ := pem.Decode([]byte(for_game.HCHAODFPrivateKey))
	if block == nil {
		logs.Info("block=nil")
		return ""
	}
	private, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	easygo.PanicError(err)

	h := crypto.Hash.New(crypto.SHA1) // 进行SHA1的散列
	h.Write([]byte(strData))
	hashed := h.Sum(nil)
	// 进行rsa加密签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, private.(*rsa.PrivateKey), crypto.SHA1, hashed)

	sign := base64.StdEncoding.EncodeToString(signature)
	logs.Info("sign:", sign)
	return sign
}

type ScanMerchantInQueryResponse struct {
	RespCode         string
	RespMsg          string
	MerNo            string
	Version          string
	PayType          string
	CompanyNo        string
	RandomStr        string
	ResultCode       string
	ChannelCompanyNo string
}

// http post 请求
func (self *WebHuiChaoPay) ReqHttps(uurl string, data interface{}) []byte {
	xmlData := self.MakeXMLParams(data)
	logs.Info("发送请求:url=", uurl, xmlData)
	params := make(url.Values)
	params.Add("requestDomain", xmlData)
	logs.Info("requestDomain------->", params)
	resp, err := http.PostForm(uurl, params)

	easygo.PanicError(err)
	defer resp.Body.Close()
	result, err1 := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err1)
	str := for_game.Base64DecodeStr(string(result))
	logs.Info("返回数据:", str)
	return result
}

// 测试下单成功后 查询订单,后续删除
func (self *WebHuiChaoPay) ReqHttpsTestQuery(uurl string, data interface{}) []byte {
	xmlData := self.MakeXMLParams(data)
	logs.Info("发送请求:url=", uurl, xmlData)
	params := make(url.Values)
	params.Add("requestDomain", xmlData)
	logs.Info("requestDomain------->", params)
	resp, err := http.PostForm(uurl, params)

	easygo.PanicError(err)
	defer resp.Body.Close()
	result, err1 := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err1)
	str := for_game.Base64DecodeStr(string(result))
	logs.Info("下单后,查询订单返回的数据:--------->", str)
	return result
}

// 代付请求
func (self *WebHuiChaoPay) ReqHttpsDF(uurl string, data interface{}) []byte {
	xmlData := self.MakeXMLParams(data)
	logs.Info("发送请求:url=", uurl, xmlData)
	params := make(url.Values)
	params.Add("transData", xmlData)
	resp, err := http.PostForm(uurl, params)
	easygo.PanicError(err)
	defer resp.Body.Close()
	result, err1 := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err1)
	b, err := base64.StdEncoding.DecodeString(string(result))
	logs.Info("返回数据:", string(b))
	return b
}

// ReqHttpsDFQuery 代付订单查询请求
func (self *WebHuiChaoPay) ReqHttpsDFQuery(uurl string, data interface{}) []byte {
	xmlData := self.MakeXMLParams(data)
	logs.Info("发送请求:url=", uurl, xmlData)
	params := make(url.Values)
	params.Add("requestDomain", xmlData)
	logs.Info("requestDomain------->", params)
	resp, err := http.PostForm(uurl, params)
	easygo.PanicError(err)
	defer resp.Body.Close()
	result, err1 := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err1)
	b, err := base64.StdEncoding.DecodeString(string(result))
	logs.Info("err------------------.", err)
	logs.Info("返回数据:", string(b))
	return b
}
func (self *WebHuiChaoPay) ReqHttpsDFCheckBankQuery(uurl string, data interface{}) []byte {
	xmlData := self.MakeXMLParams(data)
	logs.Info("发送请求:url=", uurl, xmlData)
	params := make(url.Values)
	params.Add("requestDomain", xmlData)
	logs.Info("requestDomain------->", params)
	resp, err := http.PostForm(uurl, params)
	easygo.PanicError(err)
	defer resp.Body.Close()
	result, err1 := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err1)
	logs.Info("返回数据:", string(result))
	return result
}

// 平台公钥对AESKey进行加密
func (self *WebHuiChaoPay) PubEncode(pubKey, data string) string {
	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		logs.Info("block=nil")
		return ""
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	easygo.PanicError(err)
	// 进行rsa加密
	signature, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), []byte(data))
	return base64.StdEncoding.EncodeToString(signature)
}

// 私钥进行签名
func (self *WebHuiChaoPay) RsaSign(priKey, data string) string {
	block, _ := pem.Decode([]byte(priKey))
	if block == nil {
		logs.Info("block=nil")
		return ""
	}
	private, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	easygo.PanicError(err)
	h := crypto.Hash.New(crypto.SHA1) //进行SHA1的散列
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	// 进行rsa加密签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, private.(*rsa.PrivateKey), crypto.SHA1, hashed)
	sign := base64.StdEncoding.EncodeToString(signature)
	return sign
}

// 一麻袋验签 代付的公钥是独立的.
func (self *WebHuiChaoPay) VerifySign(src, sign, public string) bool {
	//block, _ := pem.Decode([]byte(for_game.HCHAOPublic))
	block, _ := pem.Decode([]byte(public))
	if block == nil {
		logs.Info("block=nil")
		return false
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	easygo.PanicError(err)
	// 进行rsa加密
	h := crypto.Hash.New(crypto.SHA1) //进行SHA1的散列
	h.Write([]byte(src))
	hashed := h.Sum(nil)
	relSign, err := base64.StdEncoding.DecodeString(sign)
	easygo.PanicError(err)
	err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hashed, relSign)
	if err != nil {
		return false
	}
	return true
}
func (self *WebHuiChaoPay) ReqHttpsYL(uurl string, key, data string) []byte {
	timeStr := for_game.GetCurTimeString()
	params := make(url.Values)
	params.Add("requestTime", timeStr)
	params.Add("version", "1.0")
	params.Add("merchantNo", self.MerchantNo)
	params.Add("requestData", data)
	params.Add("encryptKey", self.PubEncode(for_game.HCHAOPublic, key))
	s := data + timeStr + self.MerchantNo
	strSign := base64.StdEncoding.EncodeToString([]byte(s))
	sign := self.RsaSign(for_game.HCHAOPrivateKey, strSign)
	params.Add("sign", sign)
	logs.Info("请求参数:", params)
	resp, err := http.PostForm(uurl, params)
	easygo.PanicError(err)
	defer resp.Body.Close()
	result, err1 := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err1)
	//b, err := base64.StdEncoding.DecodeString(string(result))
	//logs.Info("返回数据:", string(b))
	return result
}

// 请求参数组装
func (self *WebHuiChaoPay) MakeXMLParams(data interface{}) string {
	xmlData := self.HeadFormat
	switch value := data.(type) {
	case map[string]string:
		xmlData += " <root tx=\"1001\">"
		for k, v := range value {
			xmlData += "<" + k + ">" + v + "</" + k + ">"
		}
		xmlData += "</root>"
		logs.Info("xmlData:", xmlData)

	default:
		xl, err := xml.MarshalIndent(data, "", "")
		easygo.PanicError(err)
		xmlData += string(xl)
		logs.Info("xml:", xmlData)
	}
	b64Xml := base64.StdEncoding.EncodeToString([]byte(xmlData))
	return b64Xml
	//if t == reflect.Int32 {
	//
	//} else {
	//
	//}

}

//===============================聚合扫码支付接口==============================
// 10.1商户入驻，支付宝入驻成功不需要重复调用
func (self *WebHuiChaoPay) ZFBMerchantIn() {
	url := self.Url + "scanpay/merchantIn"
	logs.Info("url:", url)
	//支付宝参数
	zfbInfo := ZFBMerchantInfoObj{
		MerName:      "广州诺金信息科技有限公司",     // 商户名称
		ShortName:    "柠檬畅聊",             // 商户简称，会显示在订单上
		ContactName:  "梁耿",               // 联系人，支付宝必填
		ServicePhone: "4009669990",       // 客服电话
		Business:     "2015050700000000", // 行业类别：支付宝参考11.1，微信参考：11.2
		Mcc:          "5812",             // mcc码 参考11.1 支付宝必须
		ContactTag:   "06",               // 商户联系人业务 支付宝必须
		ContactType:  "AGENT",            // 联系人类型 支付宝必须
		City:         "广州市",              // 城市 支付宝必须
		District:     "白云区",              // 区县 支付宝必须
		Address:      "西槎路465号A栋3A层之G19", // 地址 支付宝必须
		Province:     "广东省",              // 省份 支付宝必须
		CardNo:       "120917141510601",  // 对公结算银行卡号 支付宝必须
		CardName:     "王晨辉",              // 对公结算银行卡持卡人姓名  支付宝必须
		//	ContactPhone:        "", // 联系人电话 花呗必传
		//	ContactMobile:       "", // 联系人手机号 花呗必传
		//	BusinessLicense:     "", // 商户证件编号 花呗必传
		//	BusinessLicenseType: "", // 商户证件类型 花呗必传
		//	IdCardNo:            "", // 联系人身份证 花呗必传
	}
	randStr := for_game.RandString(16)
	req := &ScanMerchantInRequest{
		MerNo:        self.MerchantNo,
		Version:      "1.0",
		PayType:      "ZFBZF",
		ChannelNo:    "2088831825958854", // 微信渠道号或支付宝pid:2088831825958854
		RandomStr:    randStr,
		NotifyUrl:    self.CallBackUrl,
		MerchantInfo: zfbInfo,
	}
	signStr := "MerNo=" + self.MerchantNo + "&PayType=" + req.PayType + "&RandomStr=" + req.RandomStr
	signInfo := self.SignMsg(signStr)
	logs.Info("签名信息:", signInfo)
	req.SignInfo = signInfo
	self.ReqHttps(url, req)
}

// 10.1商户入驻，微信入驻成功不需要重复调用
func (self *WebHuiChaoPay) WXMerchantIn() {
	url := self.Url + "scanpay/merchantIn"
	logs.Info("url:", url)
	//微信参数
	wxInfo := WXMerchantInfoObj{
		MerName:      "广州诺金信息科技有限公司",       // 商户名称
		ShortName:    "柠檬畅聊",               // 商户简称，会显示在订单上
		ContactName:  "梁耿",                 // 联系人，支付宝必填
		ServicePhone: "4009669990",         // 客服电话
		Business:     "545",                // 行业类别：支付宝参考11.1，微信参考：11.2
		SubAppID:     "wx414f06b4024739c8", // 微信商户appid 小程序必传
		PayPath:      "api.chihuoqun.org/", // 微信授权目录路径 微信小程序支付方式均必传(路径末尾一定要加符号/)
	}
	randStr := for_game.RandString(16)
	req := &ScanMerchantInRequest{
		MerNo:        self.MerchantNo,
		Version:      "1.0",
		PayType:      "WXZF",
		ChannelNo:    "385996806", // 微信渠道号或支付宝pid:2088831825958854
		RandomStr:    randStr,
		NotifyUrl:    self.CallBackUrl,
		MerchantInfo: wxInfo,
	}
	signStr := "MerNo=" + self.MerchantNo + "&PayType=" + req.PayType + "&RandomStr=" + req.RandomStr
	signInfo := self.SignMsg(signStr)
	logs.Info("签名信息:", signInfo)
	req.SignInfo = signInfo
	self.ReqHttps(url, req)
}

//入驻查询
type ScanMerchantInQueryRequest struct {
	MerNo     string // 商户号
	Version   string // 版本 1.0
	CompanyNo string // 扫码商户标示
	PayType   string // 签名
	ChannelNo string // 渠道号:微信渠道号/支付宝 PID
	RandomStr string // 随机字符串
	SignInfo  string //
}

//10.2商户入驻查询
func (self *WebHuiChaoPay) QueryMerchantIn(channelNo, companyNo, payType string) {
	url := self.Url + "scanpay/merchantInQuery"
	logs.Info("url:", url)
	randStr := for_game.RandString(16)
	req := &ScanMerchantInQueryRequest{
		MerNo:     self.MerchantNo,
		Version:   "1.0",
		PayType:   payType,
		ChannelNo: channelNo, // 微信渠道号或支付宝pid:2088831825958854
		RandomStr: randStr,
		CompanyNo: companyNo,
	}
	signStr := "MerNo=" + self.MerchantNo + "&CompanyNo=" + companyNo + "&PayType=" + req.PayType + "&RandomStr=" + req.RandomStr
	signInfo := self.SignMsg(signStr)
	req.SignInfo = signInfo
	self.ReqHttps(url, req)
}

//10.3商户入驻修改
func (self *WebHuiChaoPay) UpdateMerchantIn(companyNo, channelNo, payType string) {
	url := self.Url + "scanpay/merchantUpdate"
	randStr := for_game.RandString(16)
	//zfbInfo := ZFBMerchantInfoObj{
	//	MerName:      "广州诺金信息科技有限公司",     // 商户名称
	//	ShortName:    "柠檬畅聊",             // 商户简称，会显示在订单上
	//	ContactName:  "梁耿",               // 联系人，支付宝必填
	//	ServicePhone: "4009669990",       // 客服电话
	//	Business:     "2015050700000000", // 行业类别：支付宝参考11.1，微信参考：11.2
	//	Mcc:          "5812",             // mcc码 参考11.1 支付宝必须
	//	ContactTag:   "06",               // 商户联系人业务 支付宝必须
	//	ContactType:  "AGENT",            // 联系人类型 支付宝必须
	//	City:         "广州市",              // 城市 支付宝必须
	//	District:     "白云区",              // 区县 支付宝必须
	//	Address:      "西槎路465号A栋3A层之G19", // 地址 支付宝必须
	//	Province:     "广东省",              // 省份 支付宝必须
	//	CardNo:       "120917141510601",  // 对公结算银行卡号 支付宝必须
	//	CardName:     "王晨辉",              // 对公结算银行卡持卡人姓名  支付宝必须
	//	//	ContactPhone:        "", // 联系人电话 花呗必传
	//	//	ContactMobile:       "", // 联系人手机号 花呗必传
	//	//	BusinessLicense:     "", // 商户证件编号 花呗必传
	//	//	BusinessLicenseType: "", // 商户证件类型 花呗必传
	//	//	IdCardNo:            "", // 联系人身份证 花呗必传
	//}
	wxInfo := WXMerchantInfoObj{
		MerName:      "广州诺金信息科技有限公司",       // 商户名称
		ShortName:    "柠檬畅聊",               // 商户简称，会显示在订单上
		ContactName:  "梁耿",                 // 联系人，支付宝必填
		ServicePhone: "4009669990",         // 客服电话
		Business:     "545",                // 行业类别：支付宝参考11.1，微信参考：11.2
		SubAppID:     "wx905904883fc074d3", // 微信商户appid 小程序必传
		//SubAppID: "wxca76ab480f42ecb8", // 微信商户appid 小程序必传
		PayPath: "api.chihuoqun.org/", // 微信授权目录路径 微信小程序支付方式均必传(路径末尾一定要加符号/)
	}
	req := &ScanMerchantInRequest{
		CompanyNo:    companyNo,
		MerNo:        self.MerchantNo,
		Version:      "1.0",
		PayType:      payType,
		ChannelNo:    channelNo, // 微信渠道号或支付宝pid:2088831825958854
		RandomStr:    randStr,
		MerchantInfo: wxInfo,
	}
	//CompanyNo=CompanyNo&MerNo=MerNo&PayType=PayType&RandomStr=RandomStr
	signStr := "CompanyNo=" + companyNo + "&MerNo=" + self.MerchantNo + "&PayType=" + payType + "&RandomStr=" + randStr
	signInfo := self.SignMsg(signStr)
	logs.Info("签名信息:", signInfo)
	req.SignInfo = signInfo
	self.ReqHttps(url, req)
}

//10.4查询渠道商户认证状态

func (self *WebHuiChaoPay) QueryWeiXinAuthorize() {
	url := self.Url + "scanpay/weiXinAuthorizeQuery"
	logs.Info("url:", url)
}

//10.5支付接口
//payType:WxJsapi_OffLine或者AliJsapiPay_OffLine
func (self *WebHuiChaoPay) ScanPay(payType, amount, companyNo string) {
	url := self.Url + "pay/scanpay"
	logs.Info("url:", url)
}
func (self *WebHuiChaoPay) CreateOrder(payData *share_message.PayOrderInfo, openId string) string {
	amount := easygo.AtoFloat64(payData.GetAmount())
	tax := int64(0)
	changeGold := int64(amount * 100)
	player := GetPlayerObj(payData.GetPlayerId())
	channel := for_game.PAY_CHANNEL_HUICHAO_WX
	if payData.GetPayType() == for_game.PAY_TYPE_ZHIFUBAO {
		channel = for_game.PAY_CHANNEL_HUICHAO_ZFB
	} else if payData.GetPayType() == for_game.PAY_TYPE_BANKCARD {
		channel = for_game.PAY_CHANNEL_HUICHAO_YL
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
		PayOpenId:   easygo.NewString(openId),
		TotalCount:  easygo.NewInt32(payData.GetTotalCount()),
		Content:     easygo.NewString(payData.GetContent()),
		ExtendValue: easygo.NewString(payData.GetExtendValue()),
		BankInfo:    easygo.NewString(payData.GetPayBankNo()),
	}
	obj := for_game.CreateRedisOrder(order)
	return obj.GetOrderId()
}

//支付请求参数
type AggregatePayRequest struct {
	MerchantNo      string //商户号
	MerchantOrderNo string //商户平台订单号
	PayType         string //WxJsapi_OffLine(微信线下公众号支付)、AliJsapiPay_OffLine（支付宝线下服务窗支付）,云闪付AppAggrePay
	Amount          string //支付金额：100.00
	Subject         string //支付界面主题 ：收银台， AppAggrePay可为空
	Desc            string //订单描述 互联网支付
	CompanyNo       string //companyNo
	RandomStr       string //随机串
	SignInfo        string //签名
	AdviceUrl       string //异步通知地址
	SubAppid        string //小程序appid
	UserId          string //微信、支付宝不可为空 微信openid/支付宝userid
}

func (self *WebHuiChaoPay) RechargeOrder(web *WebHttpServer, payData *share_message.PayOrderInfo, openId string) ([]byte, *base.Fail) {
	//区别支付渠道
	if payData.GetPayType() == for_game.PAY_TYPE_WEIXIN {
		return self.WXAggregatePay(web, payData, openId)
	} else if payData.GetPayType() == for_game.PAY_TYPE_ZHIFUBAO {
		return self.ZFBAggregatePay(web, payData, openId)
	} else {
		return self.YLAggregatePay(web, payData, openId)
	}
}

//支付返回
type AggregatePayResponse struct {
	RespCode        string
	RespMsg         string
	OrderNo         string
	MerchantOrderNo string
	PayType         string
	Amount          string
	PayStr          string
	MerchantNo      string
	Subject         string
	Desc            string
}

//10.6聚合支付接口
func (self *WebHuiChaoPay) WXAggregatePay(web *WebHttpServer, payData *share_message.PayOrderInfo, openId string) ([]byte, *base.Fail) {
	payType := "WxJsapi_OffLine"
	subAppid := self.WXAppId
	orderNo := self.CreateOrder(payData, openId)
	req := &AggregatePayRequest{
		MerchantNo:      self.MerchantNo,
		MerchantOrderNo: easygo.AnytoA(orderNo),
		PayType:         payType,
		Amount:          payData.GetAmount(), //单位元，2位小数
		Subject:         "柠檬畅聊支付",
		//Desc:            payData.GetProduceName(),
		Desc:      "用户支付", // 因他们不喜欢用户备注发给第三方,先暂时写死
		CompanyNo: self.WXCompanyNo,
		RandomStr: for_game.RandString(16),
		AdviceUrl: self.CallBackUrl,
		SubAppid:  subAppid,
		UserId:    openId,
	}
	bytes, _ := json.Marshal(req)
	web.LogPay("AggregatePayRequest--------->" + string(bytes))
	rsp := self.AggregatePay(req)
	data := AggregatePayResponse{}
	// 对返回的数据进行base64解码
	reply := for_game.Base64DecodeStr(string(rsp))
	logs.Info("汇潮微信支付返回的数据--->", reply)
	err := xml.Unmarshal([]byte(reply), &data)
	if err != nil {
		return nil, easygo.NewFailMsg("第三方支付放回数据解析异常")
	}
	logs.Info("data: %+v", data)
	if data.RespCode != "0000" {
		return nil, easygo.NewFailMsg("第三方支付失败:" + data.RespCode)
	}
	order := for_game.GetRedisOrderObj(data.MerchantOrderNo)
	if order == nil {
		return nil, easygo.NewFailMsg("无效的支付订单:" + data.MerchantOrderNo)
	}
	pData := &ParseMDData{
		Data: PayInfo{
			PayType:        0,
			PreparePayInfo: data.PayStr,
			OrderId:        data.MerchantOrderNo,
		},
		Result:  true,
		Code:    200,
		Message: data.RespMsg,
	}
	d, err := json.Marshal(pData)
	if err != nil {
		return nil, easygo.NewFailMsg("json.Marshal(pData):" + err.Error())
	}
	//开启定时器查单
	fun := func() {
		self.ReqCheckPayOrder(orderNo, time.Second)
	}
	easygo.AfterFunc(time.Second*5, fun)
	return d, nil

}

func (self *WebHuiChaoPay) ZFBAggregatePay(web *WebHttpServer, payData *share_message.PayOrderInfo, openId string) ([]byte, *base.Fail) {
	payType := "AliJsapiPay_OffLine"
	orderNo := self.CreateOrder(payData, openId)
	logs.Info("支付宝支付,内部订单号:-------------->", orderNo)
	req := &AggregatePayRequest{
		MerchantNo:      self.MerchantNo,
		MerchantOrderNo: easygo.AnytoA(orderNo),
		//MerchantOrderNo: easygo.AnytoA("1101160729084223486205"),
		PayType: payType,
		Amount:  payData.GetAmount(), //单位元，2位小数
		Subject: "柠檬畅聊支付",
		//Desc:            payData.GetProduceName(),
		Desc:      "用户支付", // 因他们不喜欢用户备注发给第三方,先暂时写死
		CompanyNo: self.ZFBCompanyNo,
		RandomStr: for_game.RandString(16),
		AdviceUrl: self.CallBackUrl,
		UserId:    openId,
	}
	rsp := self.AggregatePay(req)

	web.LogPay(string(rsp))
	rspBytes, err1 := base64.StdEncoding.DecodeString(string(rsp))
	easygo.PanicError(err1)
	logs.Info("ZFBAggregatePay 支付宝聚合支付返回结果: %s", string(rspBytes))
	data := AggregatePayResponse{}
	err := xml.Unmarshal(rspBytes, &data)
	if err != nil {
		return nil, easygo.NewFailMsg("第三方支付放回数据解析异常:" + err.Error())
	}
	logs.Info("data:", data)
	respDataByte, err1 := json.Marshal(data)
	if err1 != nil {
		return nil, easygo.NewFailMsg("第三方支付放回数据解析异常:" + err1.Error())
	}
	if data.RespCode != "0000" {
		return nil, easygo.NewFailMsg("第三方支付失败:" + data.RespCode)
	}
	order := for_game.GetRedisOrderObj(data.MerchantOrderNo)
	if order == nil {
		return nil, easygo.NewFailMsg("无效的支付订单:" + data.MerchantOrderNo)
	}
	order.SetExternalNo(data.OrderNo)
	pData := &ParseMDData{
		Data: PayInfo{
			PayType:        payData.GetPayType(),
			PreparePayInfo: string(respDataByte), // 支付串
			OrderId:        data.MerchantOrderNo, // 商户订单号
		},
		Result:  true,
		Code:    200,
		Message: data.RespMsg,
	}
	d, err := json.Marshal(pData)
	if err != nil {
		return nil, easygo.NewFailMsg("json.Marshal(pData):", err.Error())
	}
	//开启定时器查单
	fun := func() {
		self.ReqCheckPayOrder(orderNo, time.Second)
	}
	easygo.AfterFunc(time.Second*5, fun)
	return d, nil
}
func (self *WebHuiChaoPay) YLAggregatePay(web *WebHttpServer, payData *share_message.PayOrderInfo, openId string) ([]byte, *base.Fail) {
	payType := "AppAggrePay"
	orderNo := self.CreateOrder(payData, openId)
	req := &AggregatePayRequest{
		MerchantNo:      self.MerchantNo,
		MerchantOrderNo: easygo.AnytoA(orderNo),
		PayType:         payType,
		Amount:          payData.GetAmount(), //单位元，2位小数
		Subject:         "柠檬畅聊支付",
		Desc:            payData.GetProduceName(),
		CompanyNo:       "",
		RandomStr:       for_game.RandString(16),
		AdviceUrl:       self.CallBackUrl,
		UserId:          "",
	}
	rsp := self.AggregatePay(req)
	data := AggregatePayResponse{}
	err := xml.Unmarshal(rsp, &data)
	easygo.PanicError(err)
	logs.Info("data:", data)
	//开启定时器查单
	fun := func() {
		self.ReqCheckPayOrder(orderNo, time.Second)
	}
	easygo.AfterFunc(time.Second*5, fun)
	return rsp, nil
}

//10.6聚合支付接口
func (self *WebHuiChaoPay) AggregatePay(req *AggregatePayRequest) []byte {
	url := self.Url + "pay/aggregatePay"
	logs.Info("url:", url)
	signStr := "AdviceUrl=" + req.AdviceUrl + "&Amount=" + req.Amount + "&MerchantNo=" + req.MerchantNo + "&MerchantOrderNo=" + req.MerchantOrderNo + "&PayType=" + req.PayType + "&RandomStr=" + req.RandomStr
	signInfo := self.SignMsg(signStr)
	req.SignInfo = signInfo
	return self.ReqHttps(url, req)

}

//10.7微信配置appid或支付目录
func (self *WebHuiChaoPay) SubMerchantConf() {
	url := self.Url + "scanpay/subMerchantConf"
	logs.Info("url:", url)
}

//10.8订单查询接口
func (self *WebHuiChaoPay) MerchantBatchQueryAPI() {
	url := self.Url + "merchantBatchQueryAPI"
	logs.Info("url:", url)

}

//10.10（银联）行业码支付获取用户识别码(UserId)接口
//redirectURL=商户用于接收用户识别码的地址
func (self *WebHuiChaoPay) QueryUserCode() {
	url := self.Url + "queryUserCode?redirectURL="
	logs.Info("url:", url)

}

//10.11支付宝/微信获取用户userid
func (self *WebHuiChaoPay) ScanPayConf() {
	url := self.Url + "scanpay/scanpayConf"
	logs.Info("url:", url)
}

//===============================银联支付接口(认证支付)==============================
type AttestationPayReq struct {
	MerchantOrderNo string `json:"merchantOrderNo"` //商户订单号
	Amount          string `json:"amount"`          //金额
	Products        string `json:"products"`        //产品名称
	Remark          string `json:"remark"`          //备注
	CardNo          string `json:"cardNo"`          //卡号
	BankName        string `json:"bankName"`        //银行名称
	CardType        string `json:"cardType"`        //卡类型debit(借记卡)， credit(贷记卡)
	AccountName     string `json:"accountName"`     // 账户名
	//	cvn2            string //安全码,可空，卡类型贷记卡时必填
	//	validate        string //有效期，可空
	IdentifNo string `json:"identifNo"` //证件号
	IdentType string `json:"identType"` //证件类型
	Phone     string `json:"phone"`     //手机号码
	NotifyUrl string `json:"notifyUrl"` //回调通知地址
}

//ECB PKCS5Padding
func (self *WebHuiChaoPay) AesEncrypt(src, key string) string {
	newKey, err := base64.StdEncoding.DecodeString(key)
	if src == "" {
		return ""
	}

	block, err := aes.NewCipher(newKey)
	easygo.PanicError(err)
	ecb := for_game.NewECBEncrypter(block)
	content := []byte(src)
	content = for_game.PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return base64.StdEncoding.EncodeToString(crypted)
}

//7.1认证支付接口(银联支付) 银联短信验证码请求
func (self *WebHuiChaoPay) ReqPaySMSApi(payData *share_message.PayOrderInfo, player *Player) (easygo.KWAT, *base.Fail) {
	url1 := self.Url + "rp/attestationPay"
	logs.Info("url:", url1)
	orderNo := self.CreateOrder(payData, "")
	bankInfo := player.GetBankInfo(payData.GetPayBankNo())
	req := AttestationPayReq{
		MerchantOrderNo: orderNo,
		Amount:          payData.GetAmount(),
		Products:        payData.GetProduceName(),
		Remark:          "柠檬畅聊充值",
		CardNo:          payData.GetPayBankNo(),
		BankName:        bankInfo.GetBankName(),
		CardType:        "debit",
		AccountName:     player.GetRealName(),
		IdentifNo:       player.GetPeopleId(),
		IdentType:       "IDCard",
		Phone:           bankInfo.GetBankPhone(),
		NotifyUrl:       self.CallBackUrl,
	}
	logs.Info("req:", req)
	jsData, err := json.Marshal(req)
	easygo.PanicError(err)
	key := self.GenerateAesKey()
	requestData := self.AesEncrypt(base64.StdEncoding.EncodeToString(jsData), key)
	rsp := self.ReqHttpsYL(url1, key, requestData)
	logs.Info("银行卡支付返回:", string(rsp))
	mpData := make(easygo.KWAT)
	err = json.Unmarshal(rsp, &mpData)
	easygo.PanicError(err)
	logs.Info(mpData)
	if mpData.GetString("code") == "SUCCESS" {
		logs.Info("支付请求成功")
		mpData.Add("mch_order_no", orderNo)
		return mpData, nil
	}
	return nil, easygo.NewFailMsg(mpData.GetString("message"))
}
func (self *WebHuiChaoPay) GenerateAesKey() string {
	r := easygo.RandString(16)
	key := base64.StdEncoding.EncodeToString([]byte(r))
	return key
}

type AttestationPaySmSReq struct {
	RequestNo string `json:"requestNo"` //请求流水
	MsgCode   string `json:"msgCode"`   //短信验证码
}

//7.2 银联短信支付
func (self *WebHuiChaoPay) ReqSMSPayApi(reqMsg *client_hall.BankPaySMS) (easygo.KWAT, *base.Fail) {
	url1 := self.Url + "rp/payByMsg"
	logs.Info("url:", url1)
	data := &AttestationPaySmSReq{
		RequestNo: reqMsg.GetOrderNo(),
		MsgCode:   reqMsg.GetSMS(),
	}

	jsData, err := json.Marshal(&data)
	easygo.PanicError(err)
	logs.Info("data:", string(jsData))
	key := self.GenerateAesKey()
	requestData := self.AesEncrypt(base64.StdEncoding.EncodeToString(jsData), key)
	rsp := self.ReqHttpsYL(url1, key, requestData)
	logs.Info("银行卡支付返回:", string(rsp))
	mpData := make(easygo.KWAT)
	err = json.Unmarshal(rsp, &mpData)
	easygo.PanicError(err)
	logs.Info(mpData)
	if mpData.GetString("code") == "SUCCESS" {
		logs.Info("支付请求成功")
		return mpData, nil
	}
	return nil, easygo.NewFailMsg(mpData.GetString("message"), for_game.FAIL_MSG_CODE_1015)
}

//7.3查询订单 返回成功示例: <?xml version="1.0" encoding="UTF-8" standalone="yes"?><root><merCode>50592</merCode><beginDate></beginDate><endDate></endDate><resultCount>1</resultCount><pageIndex></pageIndex><pageSize>100</pageSize><resultCode>00</resultCode><list><orderNumber>1101159711174990835432</orderNumber><orderDate>2020-07-30 19:43:01</orderDate><orderAmount>0.02</orderAmount><bankNo>1582096467</bankNo><channelNo>272020073022001443091428535100</channelNo><orderStatus>1</orderStatus><gouduiStatus>1</gouduiStatus><refundStatus>0</refundStatus></list></root>
// 失败: <?xml version="1.0" encoding="UTF-8" standalone="yes"?><root><merCode>50592</merCode><beginDate></beginDate><endDate></endDate><resultCount>1</resultCount><pageIndex></pageIndex><pageSize>100</pageSize><resultCode>00</resultCode><list><orderNumber>1101159711212973978134</orderNumber><orderDate>2020-07-30 19:49:21</orderDate><orderAmount>1.0</orderAmount><bankNo>null</bankNo><channelNo>null</channelNo><orderStatus>3</orderStatus><gouduiStatus>0</gouduiStatus><refundStatus>0</refundStatus></list></root>
func (self *WebHuiChaoPay) GetPerTime(t time.Duration) time.Duration {
	if t >= 0 && t < 300*time.Second {
		return 5 * time.Second //大于0秒小于5分钟，每10秒执行一次
	} else if t >= 300*time.Second && t < 900*time.Second {
		return 30 * time.Second //大于5分钟小于15分钟，每10秒执行一次
	} else if t >= 900*time.Second && t < 1800*time.Second {
		return 60 * time.Second //大于15分钟小于30分钟，每30秒执行一次
	}
	return 600 * time.Second //5分钟执行一次
}
func (self *WebHuiChaoPay) ReqCheckPayOrder(orderId string, t time.Duration) {
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
		logs.Info("汇潮支付手动查询订单超过半小时无结果，订单状态和支付状态转取消,orderId:--->", orderId)
		return
	}
	t += perTime

	url := self.Url + "merchantBatchQueryAPI"
	logs.Info("url:", url)

	sign := self.SignMsg(self.MerchantNo)
	beginTime := "" // 20200729114516
	endTime := ""   // 20200729122700
	m := make(map[string]string)
	m["merCode"] = self.MerchantNo
	m["orderNumber"] = orderId
	m["beginTime"] = beginTime
	m["endTime"] = endTime
	m["sign"] = sign
	replyByge := self.ReqHttpsTestQuery(url, m)
	var reply YemadaiPayQueryReply
	err := xml.Unmarshal(replyByge, &reply)
	if err != nil {
		logs.Error("汇潮支付手动查询订单反序列化失败,orderId: %s,err: %s", orderId, err.Error())
		return
	}
	logs.Info("汇潮支付手动查询订单返回结果------------>%+v", reply)
	if reply.List.OrderStatus == "1" && reply.List.GouduiStatus == "1" { // 如果成功, 执行发货功能
		order.SetPayStatus(for_game.PAY_ST_FINISH)
		//order.SetStatus(for_game.ORDER_ST_FINISH)
		RechargeGoldToPlayer(order.GetPlayerId(), order.GetOrderId())
		res := true

		//如果充值，通知前端充值结果
		if order.GetSourceType() == for_game.GOLD_TYPE_CASH_IN {
			msg := &share_message.RechargeFinish{
				Amount:        easygo.NewInt64(order.GetAmount()),
				TradeNo:       easygo.NewString(order.GetOrderId()),
				PayFinishTime: easygo.NewInt64(order.GetOverTime()),
				Result:        easygo.NewBool(res),
			}
			SendMsgToHallClientNew(order.GetPlayerId(), "RpcRechargeMoneyFinish", msg)
		}
		logs.Info("汇潮支付手动查询订单,并内部处理发货成功,orderId: %s", orderId)
	} else {
		// 2 的倍数 秒后重新查询
		logs.Info("汇潮支付手动查询订单处理中，重新查询,orderId: %s", orderId)
		fun := func() {
			self.ReqCheckPayOrder(orderId, t)
		}
		easygo.AfterFunc(perTime, fun)
	}
}

// YemadaiPayQueryReply 下单成功后查询订单返回结构体
type YemadaiPayQueryReply struct {
	MerCode     string                   `xml:"merCode"`
	BeginDate   string                   `xml:"beginDate"`
	EndDate     string                   `xml:"endDate"`
	ResultCount string                   `xml:"resultCount"`
	PageIndex   string                   `xml:"pageIndex"`
	PageSize    string                   `xml:"pageSize"`
	ResultCode  string                   `xml:"resultCode"`
	List        YemadaiPayQueryReplyList `xml:"list"`
}

type YemadaiPayQueryReplyList struct {
	OrderNumber  string `xml:"orderNumber"`
	OrderDate    string `xml:"orderDate"`
	OrderAmount  string `xml:"orderAmount"`
	BankNo       string `xml:"bankNo"`
	ChannelNo    string `xml:"channelNo"`
	OrderStatus  string `xml:"orderStatus"`
	GouduiStatus string `xml:"gouduiStatus"`
	RefundStatus string `xml:"refundStatus"`
}

//处理支付结果
func (self *WebHuiChaoPay) DealBankPayResult(data easygo.KWAT, delay ...time.Duration) {
	code := data.GetString("code")
	orderId := data.GetString("requestNo")
	order := for_game.GetRedisOrderObj(orderId)

	if order == nil {
		logs.Info("无效的支付订单")
		return
	}
	if code == "SUCCESS" {
		logs.Info("支付请求完成，等待异步通知结果")
		return
	}
}

type yemadai struct {
	AccountNumber string       `xml:"accountNumber" json:"accountNumber"` // 一麻袋系统的数字帐户(商户号)
	NotifyURL     string       `xml:"notifyURL"     json:"notifyURL"`     // 退回通知，转账失败，退回通知地址
	Tt            string       `xml:"tt"            json:"tt"`            // 是否需要加急处理(加急手续费多4元)0(普通)
	SignType      string       `xml:"signType"      json:"signType"`      // 签名方式，暂只支持RSA
	TransferList  TransferList `xml:"transferList"  json:"transferList"`

	// --------------代付查询--------------
	MerchantNumber string `xml:"merchantNumber" json:"merchantNumber"`
	Sign           string `xml:"sign"           json:"sign"`
	RequestTime    string `xml:"requestTime"    json:"requestTime"`
	MertransferID  string `xml:"mertransferID" json:"mertransferID"` // 订单号

}
type TransferList struct {
	TransId     string `xml:"transId"     json:"transId"`     // 交易流水ID
	BankCode    string `xml:"bankCode"    json:"bankCode"`    // 银行名称(传中文简称)
	Provice     string `xml:"provice"     json:"provice"`     // 开户省份(固定的正确的省市信息即可)
	City        string `xml:"city"        json:"city"`        // 开户市
	BranchName  string `xml:"branchName"  json:"branchName"`  // 支行名称,仅填写支行信息,如:太平路支行。仅需填写支行名称不需要银行名称
	AccountName string `xml:"accountName" json:"accountName"` // 开户名称(银行卡的预留姓名)
	IdNo        string `xml:"idNo"        json:"idNo"`        // 身份证
	CardNo      string `xml:"cardNo"      json:"cardNo"`      // 卡号
	Phone       string `xml:"phone"       json:"phone"`       // 手机号
	Amount      string `xml:"amount"      json:"amount"`      // 金额,比如:100.08;90.00(元)
	Remark      string `xml:"remark"      json:"remark"`      // 备注
	SecureCode  string `xml:"secureCode"  json:"secureCode"`  // 签名信息
}

// ===============================代付API=============================
// 代付返回数据:  <?xml version="1.0" encoding="UTF-8" standalone="yes"?><yemadai><errCode>0000</errCode><transferList><resCode>0000</resCode><transId>2201159563214689749109</transId><accountName>黄家茵</accountName><cardNo>6214633131067889708</cardNo><amount>1.00</amount><remark>测试转账</remark><secureCode>EfQF3StAcf3cSdFb5TsXuZWhZyf1H1QmPCNXfY8JQOti/yBg/G1Bsqa69Ozyu5L0bmyPdbF6S6ttPCwOJPndb4G61LK4ZZ9Q34qLOay70bu5n0AKvzzh/yCgs9Mz7He7SnLN+WqT7hhxawf4Np91hCXZ4CSjwoHsAyKqzgbnN6w=</secureCode></transferList></yemadai>
// TransferFixed 5.1代付接口
func (self *WebHuiChaoPay) TransferFixed(order *share_message.Order) *base.Fail {
	money := float64(order.GetAmount()+order.GetTax()+order.GetRealTax()) / 100.0
	amount := fmt.Sprintf("%.2f", money)

	url := self.Url + "transfer/transferFixed"
	orderNo := order.GetOrderId()
	player := GetPlayerObj(order.GetPlayerId())         // 获取玩家
	bankInfo := player.GetBankInfo(order.GetBankInfo()) // 获取银行卡信息
	signParam := self.HuiChaoDFSignParam(orderNo, self.MerchantNo, bankInfo.GetBankId(), easygo.AnytoA(amount))
	secureCode := self.DFSignMsg(signParam)
	req := &yemadai{
		AccountNumber: self.MerchantNo,
		NotifyURL:     self.CallBackUrlDF,
		Tt:            "0",
		SignType:      "RSA",
		TransferList: TransferList{
			TransId:     orderNo,
			BankCode:    bankInfo.GetBankName(), // 银行名称
			Provice:     bankInfo.GetProvice(),
			City:        bankInfo.GetCity(),
			AccountName: player.GetRealName(),
			CardNo:      bankInfo.GetBankId(),
			Phone:       player.GetPhone(),
			Amount:      amount,
			Remark:      order.GetNote(),
			SecureCode:  secureCode,
		},
	}

	df := self.ReqHttpsDF(url, req)
	// 解析代付返回的数据
	dfResult, errMsg := self.ParsDFResp(df)
	logs.Info("代付返回的数据", dfResult)
	if errMsg.GetCode() != "0000" {
		//受理失败
		return easygo.NewFailMsg("受理失败")
	}
	// 启动定时查询订单:5秒后查询
	fun := func() {
		self.TransferQueryFixed(orderNo, time.Second*5)
	}
	easygo.AfterFunc(time.Second*5, fun)
	return easygo.NewFailMsg("提现下单成功", for_game.FAIL_MSG_CODE_SUCCESS)
}

// todo 下面方法暂时注释,用来测试对接代付接口回调使用.测试完成会删除.
//func (self *WebHuiChaoPay) TransferFixed() {
//	url := self.Url + "transfer/transferFixed"
//	orderNo := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_OUT, for_game.GOLD_TYPE_CASH_OUT)
//
//	//player := GetPlayerObj(order.GetPlayerId())         // 获取玩家
//	//bankInfo := player.GetBankInfo(order.GetBankInfo()) // 获取银行卡信息
//	signParam := self.HuiChaoDFSignParam(orderNo, self.MerchantNo, "6214633131067889708", "0.01")
//	secureCode := self.DFSignMsg(signParam)
//	req := &yemadai{
//		AccountNumber: self.MerchantNo,
//		NotifyURL:     self.CallBackUrl,
//		Tt:            "0",
//		SignType:      "RSA",
//		TransferList: TransferList{
//			TransId:     orderNo,
//			BankCode:    "广州银行", // 银行名称
//			Provice:     "广东省",
//			City:        "广州市",
//			AccountName: "黄家茵",
//			CardNo:      "6214633131067889708",
//			Phone:       "13168180383",
//			Amount:      "0.01",
//			Remark:      "测试0.01",
//			SecureCode:  secureCode,
//		},
//	}
//
//	df := self.ReqHttpsDF(url, req)
//	logs.Info("ReqHttpsDF----->", string(df))
//}

// HuiChaoDFSignParam 拼装汇潮代付签名参数
func (self *WebHuiChaoPay) HuiChaoDFSignParam(transId, accountNumber, cardNo, amount string) string {
	var sb strings.Builder
	sb.WriteString("transId=")
	sb.WriteString(transId)
	sb.WriteString("&")
	sb.WriteString("accountNumber=")
	sb.WriteString(accountNumber)
	sb.WriteString("&")
	sb.WriteString("cardNo=")
	sb.WriteString(cardNo)
	sb.WriteString("&")
	sb.WriteString("amount=")
	sb.WriteString(amount)
	return sb.String()
}

// HuiChaoQueryDFSignParam 拼装汇潮代付查单签名参数
func (self *WebHuiChaoPay) HuiChaoQueryDFSignParam(merchantNumber, requestTime string) string {
	var sb strings.Builder
	sb.WriteString(merchantNumber)
	sb.WriteString("&")
	sb.WriteString(requestTime)
	return sb.String()
}

//生成代付订单
//生成内部订单
func (self *WebHuiChaoPay) CreateDFOrder(payData *client_hall.WithdrawInfo, player *Player, setting *share_message.PaymentSetting) string {
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
		Note:        easygo.NewString("汇潮代付"),
		Tax:         easygo.NewInt64(-tax),
		PlatformTax: easygo.NewInt64(-setting.GetPlatformTax()),
		RealTax:     easygo.NewInt64(-setting.GetRealTax()),
		Operator:    easygo.NewString("system"),
		PayChannel:  easygo.NewInt32(for_game.PAY_CHANNEL_HUICHAO_DF),
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

// TransferQueryFixed 6.1代付订单查询内部逻辑
func (self *WebHuiChaoPay) TransferQueryFixed(orderId string, t time.Duration) {
	t += 10 * time.Second
	if t > 86400*time.Second {
		logs.Info("订单超过一天无结果，转人工:", orderId)
		return
	}
	hcResp := self.transferQueryFixedToHuiChao(orderId)
	logs.Info("transferQueryFixedToHuiChao------>", string(hcResp))
	// 解析查单返回的参数
	result, errMsg := self.ParsDataDFQuery(hcResp)
	logs.Info("查单返回结果:", result, errMsg.GetCode(), errMsg.GetReason())
	order := for_game.GetRedisOrderObj(orderId)
	switch errMsg.GetCode() {
	case "11": // 00-成功;11-失败;22-处理中
		logs.Info("汇潮代付查单,订单状态的结果为失败状态")
		order.SetChanneltype(for_game.CHANNEL_MAN_MAKE)
		order.SetExtendValue(result.Transfer.Memo)
	case "00":
		logs.Info("代付完成")
		order.SetPayStatus(for_game.PAY_ST_FINISH)
		order.SetStatus(for_game.ORDER_ST_FINISH)
		order.SetExtendValue("交易完成")
		HandleAfterRecharge(orderId)
	default:
		//定时再次请求
		f := func() {
			self.TransferQueryFixed(orderId, t)
		}
		easygo.AfterFunc(t, f)
	}

}

// YemadaiQeury 查询订单结构体
type YemadaiQeury struct {
	MerCode     string `xml:"merCode"`
	OrderNumber string `xml:"orderNumber"`
	BeginTime   string `xml:"beginTime"`
	EndTime     string `xml:"endTime"`
	Sign        string `xml:"sign"`
	Tx          string `xml:"tx"`
}

//10.9订单退款接口
func (self *WebHuiChaoPay) MerchantRefundAPI(orderId string) {
	url := self.Url + "merchantRefundAPI"
	logs.Info("url:", url)

}

// transferQueryFixedToHuiChao 去汇潮查询订单
func (self *WebHuiChaoPay) transferQueryFixedToHuiChao(orderId string) []byte {
	p := "transfer/transferQueryFixed"
	url := fmt.Sprintf("%s%s", self.Url, p)
	logs.Info("url:", url)

	// 订单信息
	requestTime := for_game.GetCurTimeString3()
	signParam := self.HuiChaoQueryDFSignParam(self.MerchantNo, requestTime)
	sign := self.DFSignMsg(signParam)

	req := &yemadai{
		MerchantNumber: self.MerchantNo,
		SignType:       "RSA",
		Sign:           sign,
		RequestTime:    requestTime,
		MertransferID:  orderId,
	}
	return self.ReqHttpsDFQuery(url, req)
}

// 代付查询返回结构体
type huiChaoDFQueryResp struct {
	Code     string                   `xml:"code"`
	Transfer huiChaoDFQueryRespDetail `xml:"transfer"`
}
type huiChaoDFQueryRespDetail struct {
	MertransferID string `xml:"mertransferID"`
	Amount        string `xml:"amount"`
	State         string `xml:"state"`
	Date          string `xml:"date"`
	Memo          string `xml:"memo"`
}

// ParsDataDF 解析代付查询返回结果. 响应示例: <?xml version="1.0" encoding="UTF-8" standalone="yes"?><yemadai><code>0000</code><transfer><mertransferID>2201159563214689749109</mertransferID><amount>1.00</amount><state>00</state><date>2020-07-13 17:20:15</date><memo>成功[00000000000]</memo></transfer></yemadai>
func (self *WebHuiChaoPay) ParsDataDFQuery(data []byte) (*huiChaoDFQueryResp, *base.Fail) {
	resp := new(huiChaoDFQueryResp)
	err := xml.Unmarshal(data, &resp)
	easygo.PanicError(err)
	return resp, easygo.NewFailMsg(resp.Transfer.Memo, resp.Transfer.State)
}

// huiChaoDfResp 汇潮代付返回的结构体
type huiChaoDfResp struct {
	ErrCode      string              `xml:"errCode"`
	TransferList huiChaoDfRespDetail `xml:"transferList"`
}

// huiChaoDfRespDetail 汇潮代付返回结构体明细
type huiChaoDfRespDetail struct {
	ResCode     string `xml:"resCode"`
	TransId     string `xml:"transId"`
	AccountName string `xml:"accountName"`
	CardNo      string `xml:"cardNo"`
	Amount      string `xml:"amount"`
	Remark      string `xml:"remark"`
	SecureCode  string `xml:"secureCode"`
}

// ParsDFResp 解析代付返回的数据
func (self *WebHuiChaoPay) ParsDFResp(data []byte) (*huiChaoDfResp, *base.Fail) {
	resp := new(huiChaoDfResp)
	err := xml.Unmarshal(data, &resp)
	easygo.PanicError(err)
	return resp, easygo.NewFailMsg(resp.TransferList.Remark, resp.ErrCode)
}

// CheckBalance 7.1查询账户余额接口
func (self *WebHuiChaoPay) CheckBalance() {
	url := self.Url + "checkBalance"
	logs.Info("url:", url)
	requestTime := for_game.GetCurTimeString3()
	sign := self.DFSignMsg(fmt.Sprintf("%s%s", self.MerchantNo, requestTime))
	// 封装请求参数
	req := &CheckBalanceRequest{
		MerNo:       self.MerchantNo,
		RequestTime: requestTime,
		SignType:    "RSA",
		SignInfo:    sign,
	}
	// 发送请求
	checkBalanceResp := self.ReqHttpsDFCheckBankQuery(url, req)
	// 结果处理
	logs.Info("checkBalanceResp---->", string(checkBalanceResp))
}

type CheckBalanceRequest struct {
	MerNo       string `xml:"MerNo"`
	RequestTime string `xml:"RequestTime"`
	SignType    string `xml:"signType"`
	SignInfo    string `xml:"SignInfo"`
}
