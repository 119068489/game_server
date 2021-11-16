package for_game

import (
	"fmt"
	"game_server/easygo"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/astaxie/beego/logs"
)

type AliYunSMS struct{}

func (*AliYunSMS) SendMessageCode(phone, code string, areaCode string) bool {
	accesskeyId := easygo.YamlCfg.GetValueAsString("ALIYU_ACCESSKEYID")
	accesskeySecret := easygo.YamlCfg.GetValueAsString("ALIYU_ACCESSKEYSECRET")
	area := easygo.YamlCfg.GetValueAsString("ALIYU_AREA")
	signName := easygo.YamlCfg.GetValueAsString("ALIYU_SIGNNAME")
	templateCode := easygo.YamlCfg.GetValueAsString("ALIYU_TEMPLEATECODE")
	if areaCode != "" { //国际号码带上国际编号
		if areaCode == "855" { //柬埔寨号码特殊处理
			phone = areaCode + phone[1:] //去掉第一位
		} else {
			phone = areaCode + phone
		}
		signName = easygo.YamlCfg.GetValueAsString("ALIYU_SIGNNAME_INTERNATIONAL")
		templateCode = easygo.YamlCfg.GetValueAsString("ALIYU_TEMPLEATECODE_INTERNATIONAL")
	}
	logs.Info("发送号码:", phone, templateCode, signName)
	//accesskeyId := "LTAI4FdUu8fXogvv9penpDCK"
	//accesskeySecret := "mUe3dSYqxCAo3X5HvSdiInQt7690aG"
	//area := "cn-hangzhou"
	//signName := "柠檬畅聊"
	//templateCode := "SMS_180350717"

	client, err := sdk.NewClientWithAccessKey(area, accesskeyId, accesskeySecret)
	if err != nil {
		panic(err)
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = area
	request.QueryParams["PhoneNumbers"] = phone
	request.QueryParams["SignName"] = signName
	request.QueryParams["TemplateCode"] = templateCode
	request.QueryParams["TemplateParam"] = fmt.Sprintf("{\"code\":\"%s\"}", code)

	_, err1 := client.ProcessCommonRequest(request)
	if err1 != nil {
		logs.Error("phone:%s code:%s reason:%s", phone, code, err1.Error())
		return false
	}
	//logs.Error("phone:%s code:%s reason:%s", phone, code, response.GetHttpContentString())
	return true
}
func (*AliYunSMS) SendMessageCodeEx(phone, templateParam string, isSuccess, isInternational bool) bool {
	panic("实现方法")
}
func (*AliYunSMS) SendWarningSMS(phone, templateParam []string) {
	panic("实现方法")
}
