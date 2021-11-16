package for_game

// 短信运营商
const (
	SMS_BUSINESS_TC  = "tencent" // 腾讯
	SMS_BUSINESS_ALI = "ali"     // 阿里
)

/**
发送短信的接口
*/
type ISMS interface {
	/*
		phone: 手机号
		code: 发送的消息内容 如888888
		areaCode: 国家码.
	*/
	SendMessageCode(phone, code string, areaCode string) bool
	/**
	发送注销短信,t:短信请求类型：1短信验证码，2短信内容(暂时删除t参数)
	 phone 手机号,前面需要加国家吗,如 +8613800000000
	templateParam:模板参数
	 isSuccess 是否注销成功的模板
	isInternational是否国际
	*/
	SendMessageCodeEx(phone, templateParam string, isSuccess, isInternational bool) bool

	/**
	支付预警功能短信
	phone: 手机号数组,如腾讯云 []string{+8613711112222}， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号
	templateParam: 模板参数 腾讯云 显示时间.
	*/
	SendWarningSMS(phone, templateParam []string)
}

// 初始化短信运营商
func NewSMSInst(business string) ISMS {
	switch business {
	case SMS_BUSINESS_TC:
		return &TenCentSMS{}
	case SMS_BUSINESS_ALI:
		return &AliYunSMS{}
	}
	return nil
}
