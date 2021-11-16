package hall

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
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
	"strings"
	"time"

	"software.sslmate.com/src/go-pkcs12"

	"github.com/astaxie/beego/logs"
	"github.com/axgle/mahonia"
)

/*
	鹏聚支付服务入口
*/

type WebPengJuPay struct {
	MerchantNo string   //商户号
	Url        string   //接口地址
	HeadFormat string   //头格式
	BankList   []string //支持银行列表
}

//************************************************请求代付结构******************************
type INFO_DF struct {
	TRX_CODE   string //交易代码:100005
	VERSION    string //版本:04
	DATA_TYPE  string //数据格式 2：xml 格式
	LEVEL      string //处理级别:0 实时处理
	USER_NAME  string //用户名:同 MERCHANT_ID
	USER_PASS  string //用户密码:无需提供
	REQ_SN     string //交易流水号
	SIGNED_MSG string //签名信息
}
type TRANS_SUMObj struct {
	BUSINESS_CODE string //业务代码
	MERCHANT_ID   string //商户代码
	SUBMIT_TIME   string //提交时间
	TOTAL_ITEM    string //总记录数
	TOTAL_SUM     string //总金额
}
type TRANS_DETAILObj struct {
	SN                 string //记录序号
	E_USER_CODE        string //用户编号
	BANK_CODE          string //银行代码
	ACCOUNT_TYPE       string //账号类型
	ACCOUNT_NO         string //账号
	ACCOUNT_NAME       string //账号名
	PROVINCE           string //开户行所在省
	CITY               string //开户行所在市
	BANK_NAME          string //开户行名称
	BANK_TYPE          string //联行号
	ACCOUNT_PROP       string //账号属性
	AMOUNT             string //金额
	CURRENCY           string //货币类型
	PROTOCOL           string //协议号
	PROTOCOL_USERID    string //协议用户编号
	ID_TYPE            string //开户证件类型
	ID                 string //证件号
	TEL                string //手机号/小灵通
	RECKON_ACCOUNT     string //清算账号
	RECKON_CURRENCY    string //清算货币类型
	TRADE_NOTIFY_URL   string //交易异步通知地址
	REFUND_NOTIFY_URLL string //退回异步通知地址
	TERMINAL_NO        string //终端号 由我方出号时分配
	CUST_USERID        string //自定义用户号
	REMARK             string //备注
	RESERVE1           string //保留域 1
	RESERVE2           string //保留域 2
}
type QUERY_TRANSObj struct {
	QUERY_SN     string //要查询的交易流水
	QUERY_REMARK string //查询备注
}
type TRANS_DETAILSObj struct {
	TRANS_DETAIL TRANS_DETAILObj
}
type BODY_DF struct {
	TRANS_SUM     TRANS_SUMObj
	TRANS_DETAILS TRANS_DETAILSObj
}

//代付请求结构
type DF_GHT_REQ struct {
	INFO INFO_DF
	BODY BODY_DF
}

//请求返回信息
type INFO_RESP struct {
	TRX_CODE   string //交易代码:100005
	VERSION    string //版本:04
	DATA_TYPE  string //数据格式 2：xml 格式
	REQ_SN     string //交易流水号
	RET_CODE   string //返回代码
	ERR_MSG    string //错误信息
	SIGNED_MSG string //签名信息
}
type RET_DETAILObj struct {
	SN           string //记录序号
	ACCOUNT_NO   string //账号
	ACCOUNT_NAME string //账号名
	AMOUNT       string //金额
	CUST_USERID  string //自定义用户号
	REMARK       string //备注
	RET_CODE     string //返回码
	ERR_MSG      string //错误文本
	RESERVE1     string //保留域 1
	RESERVE2     string //保留域 2
}
type RET_DETAILS struct {
	RET_DETAIL RET_DETAILObj
}
type BODY_RESP struct {
	RET_DETAILS RET_DETAILS
	QUERY_TRANS QUERY_TRANSObj
}

//代付响应结构
type GHT_RESP struct {
	INFO INFO_RESP
	BODY BODY_RESP
}

////************************************************请求查询结构*****************************
type BODY_CX struct {
	QUERY_TRANS QUERY_TRANSObj
}

//
//查询请求结构
type CX_GHT_REQ struct {
	INFO INFO_DF
	BODY BODY_CX
}

func NewWebPengJuPay() *WebPengJuPay {
	p := &WebPengJuPay{}
	p.Init()
	return p
}

//秘钥 090b8a929133cef439abf2e27fa2215c  机构号86038810商户号15989289979
func (self *WebPengJuPay) Init() {
	self.MerchantNo = "000000000103285" //商户号
	//self.MerchantNo = "000000000100641"                  //测试商户号
	//self.Url = "https://120.31.132.118:8181/e/merchant/" //utf-8
	self.Url = "https://dsf.sicpay.com/e/merchant/" //正式
	self.HeadFormat = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	for code, _ := range for_game.BankPayNo {
		self.BankList = append(self.BankList, code)
	}
}

//获取支持的银行卡列表
func (self *WebPengJuPay) GetSupportBankList() []*client_hall.BankData {
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

//加上<GHT></GHT>
func (self *WebPengJuPay) XMLFormatGHT(xml, name string) string {
	return strings.Replace(xml, name, "GHT", -1)
}

//发起代付
func (self *WebPengJuPay) RechargeOrder(order *share_message.Order) *base.Fail {
	//充值的钱减去手续费
	amount := easygo.AnytoA(order.GetAmount() + order.GetTax() + order.GetRealTax())
	rq := DF_GHT_REQ{
		INFO: INFO_DF{
			TRX_CODE:  "100005",
			VERSION:   "04",
			USER_NAME: self.MerchantNo,
			REQ_SN:    order.GetOrderId(),
		},
		BODY: BODY_DF{
			TRANS_SUM: TRANS_SUMObj{
				BUSINESS_CODE: "05100",
				MERCHANT_ID:   self.MerchantNo,
				TOTAL_ITEM:    "1",
				TOTAL_SUM:     amount,
			},
			TRANS_DETAILS: TRANS_DETAILSObj{
				TRANS_DETAIL: TRANS_DETAILObj{
					SN:           "1001",
					BANK_CODE:    order.GetBankCode(),
					ACCOUNT_TYPE: order.GetAccountType(),
					ACCOUNT_NO:   order.GetAccountNo(),
					ACCOUNT_PROP: order.GetAccountProp(),
					ACCOUNT_NAME: order.GetAccountName(),
					AMOUNT:       amount,
				},
			},
		},
	}
	logs.Info(rq)
	ght := self.HttpsReq("DF_GHT_REQ", rq)
	if ght != nil {
		logs.Info("结果:", ght.BODY.RET_DETAILS.RET_DETAIL.ERR_MSG)
		//启动定时查询订单:5秒后查询
		fun := func() {
			self.CheckOrder(ght.INFO.REQ_SN, time.Second*5)
		}
		easygo.AfterFunc(time.Second*5, fun)
		if ght.BODY.RET_DETAILS.RET_DETAIL.RET_CODE == "0000" {
			//交易成功处理逻辑
			return easygo.NewFailMsg(ght.BODY.RET_DETAILS.RET_DETAIL.ERR_MSG, for_game.FAIL_MSG_CODE_SUCCESS)
		} else {
			return easygo.NewFailMsg(ght.BODY.RET_DETAILS.RET_DETAIL.ERR_MSG, for_game.FAIL_MSG_CODE_1007)
		}
	}
	return easygo.NewFailMsg("RedisCache", for_game.FAIL_MSG_CODE_1007)
}

//生成内部订单
func (self *WebPengJuPay) CreateOrder(payData *client_hall.WithdrawInfo, player *Player, setting *share_message.PaymentSetting) string {
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
		PayChannel:  easygo.NewInt32(for_game.PAY_CHANNEL_PENGJU),
		BankCode:    easygo.NewString(payData.GetBankCode()),
		AccountType: easygo.NewString(payData.GetAccountType()),
		AccountNo:   easygo.NewString(payData.GetAccountNo()),
		BankInfo:    easygo.NewString(payData.GetAccountNo()),
		AccountName: easygo.NewString(payData.GetAccountName()),
		AccountProp: easygo.NewString("0"),
		PayType:     easygo.NewInt32(for_game.PAY_TYPE_BANKCARD),
		OrderDate:   easygo.NewString(for_game.GetCurTimeString2()),
	}
	obj := for_game.CreateRedisOrder(order)
	return obj.GetOrderId()
}

//查询提现订单
func (self *WebPengJuPay) CheckOrder(orderId string, t time.Duration) {
	t += 10 * time.Second
	if t > 86400*time.Second {
		logs.Info("订单超过一天无结果，转人工:", orderId)
		return
	}
	rq := CX_GHT_REQ{
		INFO: INFO_DF{
			TRX_CODE:  "200001",
			VERSION:   "04",
			DATA_TYPE: "2",
			USER_NAME: self.MerchantNo,
			REQ_SN:    orderId,
		},
		BODY: BODY_CX{
			QUERY_TRANS: QUERY_TRANSObj{
				QUERY_SN:     orderId,
				QUERY_REMARK: "查询交易",
			},
		},
	}
	logs.Info(rq)
	ght := self.HttpsReq("CX_GHT_REQ", rq)
	if ght != nil {
		ret := ght.BODY.RET_DETAILS.RET_DETAIL.RET_CODE
		retNum := easygo.AtoInt32(ret)
		order := for_game.GetRedisOrderObj(ght.INFO.REQ_SN)
		if order == nil {
			panic("无效的交易订单:" + ght.INFO.REQ_SN)
		}
		if ret == for_game.FAIL_MSG_CODE_SUCCESS {
			//交易处理完毕
			logs.Info("代付完成")
			order.SetPayStatus(for_game.PAY_ST_FINISH)
			order.SetStatus(for_game.ORDER_ST_FINISH)
			order.SetExtendValue("交易完成")
			HandleAfterRecharge(ght.INFO.REQ_SN)
		} else if ret == "0001" || ret == "0002" {
			logs.Info("代付取消关闭")
			//0001：交易数据处理失败，0002：交易撤销
			order.SetPayStatus(for_game.PAY_ST_WAITTING)
			//order.PayStatus = easygo.NewInt32(for_game.PAY_ST_CANCEL)
			//order.Status = easygo.NewInt32(for_game.ORDER_ST_CANCEL)
			order.SetChanneltype(for_game.CHANNEL_MAN_MAKE)
			order.SetExtendValue(ght.BODY.RET_DETAILS.RET_DETAIL.ERR_MSG)
		} else if retNum >= 2000 && retNum < 3000 {
			//10秒后重新查询
			logs.Info("代付处理中，重新查询")
			fun := func() {
				self.CheckOrder(orderId, t)
			}
			easygo.AfterFunc(t, fun)
		}
	}
	return
}

//
func (self *WebPengJuPay) SignMsg(strData string) string {
	//读取pfx证书
	strMsg := strings.Replace(strData, "<SIGNED_MSG></SIGNED_MSG>", "", -1)
	logs.Info("签名源串:", strMsg)
	pfx, err := ioutil.ReadFile("000000000103285.pfx")
	easygo.PanicError(err)
	private, _, _, err := pkcs12.DecodeChain(pfx, "123456")
	easygo.PanicError(err)
	h := crypto.Hash.New(crypto.SHA1) //进行SHA1的散列
	h.Write([]byte(strMsg))
	hashed := h.Sum(nil)
	// 进行rsa加密签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, private.(*rsa.PrivateKey), crypto.SHA1, hashed)

	sign := hex.EncodeToString(signature)
	logs.Info("sign:", sign)
	backData := strings.Replace(strData, "<SIGNED_MSG></SIGNED_MSG>", "<SIGNED_MSG>"+sign+"</SIGNED_MSG>", -1)
	return backData
}

func (self *WebPengJuPay) VerifySign(strXml string) bool {
	iStart := strings.Index(strXml, "<SIGNED_MSG>")
	if iStart != -1 {
		end := strings.Index(strXml, "</SIGNED_MSG>")
		signedMsg := strXml[iStart+12 : end]
		sign, err := hex.DecodeString(signedMsg)
		easygo.PanicError(err)
		strMsg := strXml[:iStart] + strXml[end+13:]
		//读取crt证书解析
		cer, err := ioutil.ReadFile("GHT_Root.crt")
		easygo.PanicError(err)
		_, restPEMBlock := pem.Decode(cer)
		x509Cert, err := x509.ParseCertificate(restPEMBlock)
		easygo.PanicError(err)
		h := crypto.Hash.New(crypto.SHA1) //进行SHA1的散列
		h.Write([]byte(strMsg))
		hashed := h.Sum(nil)
		easygo.PanicError(err)
		res := rsa.VerifyPKCS1v15(x509Cert.PublicKey.(*rsa.PublicKey), crypto.SHA1, hashed, sign)
		return res == nil
	}
	return false
}
func (self *WebPengJuPay) ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// 接口POST请求:XML格式请求
func (self *WebPengJuPay) HttpsReq(method string, ght interface{}) *GHT_RESP {
	WebServreMgr.NextJobId()
	WebServreMgr.LogPay("鹏聚提现请求:" + self.Url + ",method:" + method)
	xl, err := xml.MarshalIndent(ght, "", "")
	easygo.PanicError(err)
	sxl := strings.Replace(string(xl), method, "GHT", -1)
	msg := self.HeadFormat + sxl
	body := self.SignMsg(msg)
	logs.Info("发送请求:url=", self.Url, body)
	ssl := &tls.Config{
		//Certificates:       []tls.Certificate{nil},
		InsecureSkipVerify: true,
	}

	ssl.Rand = rand.Reader

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: ssl,
		},
	}
	WebServreMgr.LogPay("鹏聚提现body:" + body)
	// build a new request, but not doing the POST yet
	req, err := http.NewRequest("POST", self.Url, bytes.NewReader([]byte(body)))
	if err != nil {
		fmt.Println(err)
	}
	// set the Header here
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	// now POST it
	resp, err := client.Do(req)
	easygo.PanicError(err)
	if err != nil {
		fmt.Println(err)

	}
	result, err2 := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	easygo.PanicError(err2)
	res := string(result)
	WebServreMgr.LogPay("鹏聚提现返回:" + res)
	logs.Info("返回数据:", res)
	if self.VerifySign(res) {
		logs.Info("验签通过")
		ght := &GHT_RESP{}
		err := xml.Unmarshal(result, ght)
		easygo.PanicError(err)
		logs.Info("ght", ght)
		if ght.INFO.RET_CODE == "0000" {
			logs.Info(ght.INFO.ERR_MSG)
			return ght
		}
	} else {
		logs.Info("验签失败")
	}
	return nil
}

var PWebPengJuPay = NewWebPengJuPay()
