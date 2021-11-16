package hall

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestWebHuiChaoPay_RsaSign(t *testing.T) {
	//pa := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><yemadai><code>0000</code><transfer><mertransferID>2201159563214689749109</mertransferID><amount>1.00</amount><state>00</state><date>2020-07-13 17:20:15</date><memo>成功[00000000000]</memo></transfer></yemadai>`
	pa := `<?xml version="1.0" encoding="utf-8"?>
<yemadai> 
  <errCode>0000</errCode>  
  <transferList> 
    <resCode>0000</resCode>  
    <transId>test10001</transId>  
    <accountName>王五</accountName>  
    <cardNo>6222021001067998889</cardNo>  
    <amount>10.00</amount>  
    <remark>测试转账</remark>  
    <secureCode>A13230D0CBFD964621B26984D513D13F</secureCode> 
  </transferList> 
</yemadai>
`
	m := new(huiChaoDfResp)
	err := xml.Unmarshal([]byte(pa), &m)
	if err != nil {
		logs.Info(err.Error())
		return
	}
	fmt.Printf(" %+v\n", m)
}

func TestCallBackSign(t *testing.T) {
	signParam := "MerNo=50592&MerBillNo=2201159563214689749110&CardNo=6214633131067889708&Amount=0.01&Succeed=00&BillNo=2201159563214689749109"
	logs.Info(RsaSign(for_game.HCHAODFPrivateKey, signParam))
}

// 私钥进行签名
func RsaSign(priKey, data string) string {
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

func TestQeury(t *testing.T) {

	signInfo := "oiFKtZkrjdCL94HxdWc7dO198trjT6+sCuNT9PbJ2a0BroODb5FSghJ9SoBUqkpDvit5tAYdRI07eFqgAzEorqxvR+NxV73GTY5yYoJegjAI8r8jUVG6vPSoI/rw3e0O9yoWHX/S2OJciNif7mYIVp9+XoYxSy2YU2E5Ej6QKJE="
	merBillNo := "2201159691942175151198"
	cardNo := "6214633131067889708"
	amount := "0.98"
	succeed := "00"
	billNo := "GATEWAY0033696350"
	// 通过订单号查询订单信息
	src := "MerNo=" + "50592" + "&MerBillNo=" + merBillNo +
		"&CardNo=" + cardNo + "&Amount=" + amount +
		"&Succeed=" + succeed + "&BillNo=" + billNo
	//b := PWebHuiChaoPay.VerifySign(src, signInfo, for_game.HCHAODFPublic)
	b := PWebHuiChaoPay.VerifySign(src, signInfo, for_game.HCHAOPublic)
	fmt.Println(b)

	//hc := NewWebHuiChaoPay()
	//result := hc.ReqCheckPayOrder("1101159623866410321061")
	//marshal, err := json.Marshal(result)
	//if err != nil {
	//	logs.Info(err.Error())
	//	return
	//}
	//logs.Info("==============>", string(marshal))

}
func TestA(t *testing.T) {
	s := `PD94bWwgdmVyc2lvbj0nMS4wJyBlbmNvZGluZz0nVVRGLTgnIHN0YW5kYWxvbmU9J3llcyc/PjxBZ2dyZWdhdGVQYXlSZXNwb25zZT48UmVzcENvZGU+MTAwOTwvUmVzcENvZGU+PFJlc3BNc2c+c3ViX21jaF9pZOS4jnN1Yl9hcHBpZOS4jeWMuemFjTwvUmVzcE1zZz48L0FnZ3JlZ2F0ZVBheVJlc3BvbnNlPg==`
	fmt.Println(for_game.Base64DecodeStr(s))
}

func TestQueryMerchantIn(t *testing.T) {
	initializer := NewInitializer()
	defer func() { // 若是异常了,确保异步日志有成功写盘
		logger := initializer.GetBeeLogger()
		if logger != nil {
			logger.Flush()
		}
	}()
	dict := easygo.KWAT{
		"logName":  "hall",
		"yamlPath": "../../../config_hall.yaml",
	}
	initializer.Execute(dict)
	Initialize()

	PWebHuiChaoPay.QueryMerchantIn("386852991", "sweep-b90d814d534a4b219ab1fe0983f248e9", "WXZF")
}

func TestAb(t *testing.T) {
	Id := for_game.NextId(for_game.TABLE_SQUARE_COMMENT)
	fmt.Println(Id)
}

func TestWebHuiChaoPayQueryReply(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><root><merCode>50592</merCode><beginDate></beginDate><endDate></endDate><resultCount>1</resultCount><pageIndex></pageIndex><pageSize>100</pageSize><resultCode>00</resultCode><list><orderNumber>1101159711174990835432</orderNumber><orderDate>2020-07-30 19:43:01</orderDate><orderAmount>0.02</orderAmount><bankNo>1582096467</bankNo><channelNo>272020073022001443091428535100</channelNo><orderStatus>1</orderStatus><gouduiStatus>1</gouduiStatus><refundStatus>0</refundStatus></list></root>`
	var reply YemadaiPayQueryReply

	err := xml.Unmarshal([]byte(x), &reply)
	if err != nil {
		fmt.Println("--------------->", err.Error())
		return
	}
	fmt.Printf("==============>%+v", reply.List.OrderNumber)
}
