package for_game

import (
	"encoding/json"
	"game_server/easygo"

	"github.com/astaxie/beego/logs"
	v20190321 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cms/v20190321"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	v20190711 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
)

//文本内容检测 content:base64后的字符串
//talker：说话人
//teamId：所在群id，私聊则为0
type DetailResultObj struct {
	EvilLabel  string   //0正常，1可疑
	EvilType   int32    //恶意类型
	Keywords   []string //命中的关键词
	Score      int32    //命中的模型分值
	Suggestion string   //建议值:Block：打击,Review：待复审,Normal：正常
}
type TextModerationRsp struct {
	EvilFlag     int32
	EvilType     int32
	DetailResult []DetailResultObj
	Keywords     []string
}

//敏感词屏蔽:100正常 ，content为base64编码字符串
//20001：政治 20002：色情 20006：涉毒违法 20007：谩骂 20103 性感 20105：广告引流 24001：暴恐
func TextModeration(content string) *TextModerationRsp {
	c, err := v20190321.NewClient(common.NewCredential("AKIDOYR4Xst9ZicIwJrJ7ex2rPlqgY9VbIj1", "DDCVaPxb2evwpgJfleZEm4RXPAe7KOCk"), "ap-guangzhou", profile.NewClientProfile())
	easygo.PanicError(err)
	req := v20190321.NewTextModerationRequest()
	req.Content = easygo.NewString(content)
	resp, err := c.TextModeration(req)
	if err != nil {
		return nil
	}
	data := resp.Response.Data
	js, err := json.Marshal(data)
	easygo.PanicError(err)
	result := &TextModerationRsp{}
	err = json.Unmarshal(js, &result)
	easygo.PanicError(err)
	return result
}

type ImageCommon struct {
	EvilType int32    //类型
	HitFlag  int32    //判定：0正常，1可疑
	Keywords []string //关键词明细
	Labels   []string //标签
	Score    int32    //得分
}

//涉黄
type ImagePornDetect struct {
	ImageCommon
}

//性感
type ImageHotDetect struct {
	ImageCommon
}

//违法
type ImageIllegalDetect struct {
	ImageCommon
}
type ImageRrectF struct {
	Cx     float32 //logo横坐标
	Cy     float32 //logo纵坐标
	Height float32 //logo图片高读
	Rotate float32 //logo图标中信旋转
	Width  float32 //logo图标宽度
}
type ImageLogo struct {
	RrectF     ImageRrectF
	Confidence float32
	Name       string
}
type ImageLogoDetail struct {
	AppLogoDetail []ImageLogo
}

//涉政
type ImagePolityDetect struct {
	ImageCommon
	PolityLogoDetail []ImageLogo
	FaceNames        []string
	PolityItems      []string
}
type ImageCodePosition struct {
	FloatX float32 //二维码边界点X轴坐标
	FloatY float32 //二维码边界点Y轴坐标
}
type ImageCodeDetail struct {
	CodePosition []ImageCodePosition //二维码在图片中的位置，由边界点的坐标表示
	CodeCharset  string              //二维码文本的编码格式
	CodeText     string              //二维码的文本内容
	CodeType     int32               //二维码的类型：1:ONED_BARCODE，2:QRCOD，3:WXCODE，4:PDF417，5:DATAMATRIX
}

//图片二维码详情
type ImageCodeDetect struct {
	ModerationDetail []ImageCodeDetail //从图片中检测到的二维码，可能为多个
	ModerationCode   int32             //检测是否成功，0：成功，-1：出错
}

//手机模型识别
type ImagePhoneDetect struct {
	ImageCommon
}

//暴恐
type ImageTerrorDetect struct {
	ImageCommon
}
type ImageModerationRsp struct {
	EvilFlag      int32              //0正常，1可疑
	EvilType      int32              //类型:100
	CodeDetect    ImageCodeDetect    //二维码
	HotDetect     ImageHotDetect     //性感图
	IllegalDetect ImageIllegalDetect //违法图
	LogoDetect    ImageLogoDetail    //logo图
	//OCRDetect     string          //OCR图
	PhoneDetect  string            //手机检测
	PolityDetect ImagePolityDetect //涉证
	PornDetect   ImagePornDetect   //涉黄
	//Similar       string          //相似度
	TerrorDetect ImageTerrorDetect //暴力
}

//敏感图屏蔽:100正常
////20001：政治 20002：色情 20006：涉毒违法 20007：谩骂 20103 性感 20105：广告引流 24001：暴恐
func ImageModeration(url string) *ImageModerationRsp {
	c, err := v20190321.NewClient(common.NewCredential("AKIDOYR4Xst9ZicIwJrJ7ex2rPlqgY9VbIj1", "DDCVaPxb2evwpgJfleZEm4RXPAe7KOCk"), "ap-guangzhou", profile.NewClientProfile())
	easygo.PanicError(err)
	req := v20190321.NewImageModerationRequest()
	req.FileUrl = easygo.NewString(url)
	resp, err := c.ImageModeration(req)
	if err != nil {
		return nil
	}
	data := resp.Response.Data
	js1, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	result := &ImageModerationRsp{}
	err = json.Unmarshal(js1, &result)
	return result

}

type TenCentSMS struct {
}

func (*TenCentSMS) SendMessageCode(phone, code string, areaCode string) bool {
	panic("实现方法")
}

//发送短信,t:短信请求类型：1短信验证码，2短信内容(暂时删除t参数)
// phone 手机号,前面需要加国家吗,如 +8613800000000
//templateParam:模板参数
// isSuccess 是否注销成功的模板
//isInternational是否国际
func (*TenCentSMS) SendMessageCodeEx(phone, templateParam string, isSuccess, isInternational bool) bool {
	logs.Info("发送模板短信,参数:phone: %s,templateParam: %s,成功或失败的模板: %v", phone, templateParam, isSuccess)
	c, err := v20190711.NewClient(common.NewCredential(easygo.YamlCfg.GetValueAsString("TC_SMS_SECRETID"), easygo.YamlCfg.GetValueAsString("TC_SMS_SECRETKEY")), "ap-guangzhou", profile.NewClientProfile())
	easygo.PanicError(err)
	req := v20190711.NewSendSmsRequest()
	req.PhoneNumberSet = append(req.PhoneNumberSet, easygo.NewString(phone))
	var templateId string
	if isSuccess {
		templateId = easygo.YamlCfg.GetValueAsString("TC_SMS_LOGOUT_SUCCESS_TEMPLATE_ID")
	} else {
		templateId = easygo.YamlCfg.GetValueAsString("TC_SMS_LOGOUT_FAILED_TEMPLATE_ID")
	}
	smsSdkAppid := easygo.YamlCfg.GetValueAsString("TC_SMS_SDK_APPID")
	sign := easygo.YamlCfg.GetValueAsString("TC_SMS_SIGN")
	req.TemplateID = easygo.NewString(templateId)
	req.SmsSdkAppid = easygo.NewString(smsSdkAppid)
	req.Sign = easygo.NewString(sign)
	if templateParam != "" {
		req.TemplateParamSet = append(req.TemplateParamSet, easygo.NewString(templateParam))
	}
	resp, err := c.SendSms(req)
	if err != nil {
		logs.Error("phone:%s code:%s reason:%s", phone, templateParam, err.Error())
		return false
	}
	if len(resp.Response.SendStatusSet) > 0 {
		for _, status := range resp.Response.SendStatusSet {
			if *status.Code != "Ok" {
				logs.Error("发送注销短信失败,phone: %s,Code: %s,Message: %s", *status.PhoneNumber, *status.Code, *status.Message)
				continue
			}
			logs.Info("发送注销短信成功,phone: %s", *status.PhoneNumber)
		}
	}
	return true
}

// 发送提示短信 templateParam 时间
func (*TenCentSMS) SendWarningSMS(phone, templateParam []string) {
	credential := common.NewCredential(
		easygo.YamlCfg.GetValueAsString("TC_SMS_SECRETID"),
		easygo.YamlCfg.GetValueAsString("TC_SMS_SECRETKEY"),
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	/* SDK 默认用 TC3-HMAC-SHA256 进行签名，非必要请不要修改该字段 */
	cpf.SignMethod = "HmacSHA1"
	client, _ := v20190711.NewClient(credential, "ap-guangzhou", cpf)
	request := v20190711.NewSendSmsRequest()
	request.SmsSdkAppid = common.StringPtr(easygo.YamlCfg.GetValueAsString("TC_SMS_SDK_APPID"))
	/* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名，可登录 [短信控制台] 查看签名信息 */
	request.Sign = common.StringPtr(easygo.YamlCfg.GetValueAsString("TC_SMS_SIGN"))
	request.ExtendCode = common.StringPtr("0")
	/* 模板参数: 若无模板参数，则设置为空*/
	//request.TemplateParamSet = common.StringPtrs([]string{easygo.Stamp2Str(time.Now().Unix())})
	request.TemplateParamSet = common.StringPtrs(templateParam)
	request.TemplateID = common.StringPtr(easygo.YamlCfg.GetValueAsString("TC_SMS_WARNING_TEMPLATE_ID"))
	/* 下发手机号码，采用 e.164 标准，+[国家或地区码][手机号]
	 * 例如+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	request.PhoneNumberSet = common.StringPtrs(phone)
	// 通过 client 对象调用想要访问的接口，需要传入请求对象
	response, err := client.SendSms(request)
	if err != nil {
		logs.Error("SendWarningSMS phone:%v templateParam:%v reason:%s", phone, templateParam, err.Error())
		return
	}
	if len(response.Response.SendStatusSet) > 0 {
		for _, status := range response.Response.SendStatusSet {
			if *status.Code != "Ok" {
				logs.Error("SendWarningSMS 发送预警功能短信失败,phone: %s,Code: %s,Message: %s", *status.PhoneNumber, *status.Code, *status.Message)
				continue
			}
			logs.Info("SendWarningSMS 发送预警功能短信成功,phone: %s", *status.PhoneNumber)
		}
	}
}
