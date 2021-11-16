package shop

import (
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//----------------------------------
// 常用快递调用－ 聚合数据
// 在线接口文档：http://www.juhe.cn/docs/43
//----------------------------------
const APPKEY = "f3f5f417f33fee85e04dd140fd3a228c" //您申请的APPKEY

//1.常用快递查询API
func ExpressQuery(com string, no string, sPhoneLastFour string, rPhoneLastFour string) (map[string]interface{}, string) {
	//请求地址
	juheURL := "http://v.juhe.cn/exp/index"

	//初始化参数
	param := url.Values{}

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("com", com)    //需要查询的快递公司编号
	param.Set("no", no)      //需要查询的订单号
	param.Set("key", APPKEY) //应用APPKEY(应用详细页查询)
	//param.Set("dtype","") //返回数据的格式,xml或json，默认json

	//顺丰的时候设置发件人和收件人手机号
	if com == "sf" {
		if rPhoneLastFour != "" && len(rPhoneLastFour) == 4 {
			param.Set("receiverPhone", rPhoneLastFour) //收件人手机号后四位，顺丰快递需要提供senderPhone和receiverPhone其中一个
		} else if sPhoneLastFour != "" && len(sPhoneLastFour) == 4 {
			param.Set("senderPhone", sPhoneLastFour) //寄件人手机号后四位，顺丰快递需要提供senderPhone和receiverPhone其中一个
		}
	}

	//发送请求
	data, err := Post(juheURL, param)
	if err != nil {
		logs.Error("快递信息请求失败,错误信息:\r\n%v", err)
		return nil, EXPRESS_QUERY_ERROR_CODE_999
	} else {

		netReturn, err := util.JsonDecode(([]byte)(data))

		if nil == err {
			return netReturn, ""
		} else {
			logs.Error("快递信息请求失败,错误信息:\r\n%v", err)
			return nil, EXPRESS_QUERY_ERROR_CODE_999
		}

	}
	return nil, EXPRESS_QUERY_ERROR_CODE_999
}

// get 网络请求
func Get(apiURL string, params url.Values) (rs []byte, err error) {
	var Url *url.URL
	Url, err = url.Parse(apiURL)
	if err != nil {
		logs.Error("快递GET请求解析url错误:\r\n%v", err)
		return nil, err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	resp, err := http.Get(Url.String())
	if err != nil {
		logs.Error("快递GET请求err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// post 网络请求 ,params 是url.Values类型
func Post(apiURL string, params url.Values) (rs []byte, err error) {
	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		logs.Error("快递POST请求err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func GetExpressComInfos() *share_message.ExpressComInfosResult {

	var comInfos []*share_message.ExpressCom = []*share_message.ExpressCom{}
	var commonUseComInfos []*share_message.ExpressCom = []*share_message.ExpressCom{}

	var expressName string
	//B字母开头
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString(""),
		Name: easygo.NewString("B"),
	})

	//这个为了页面显示效果放这里是百世快递（接口那里是(汇通)百世快递）
	expressName, _ = GetExpressNamePhone("ht")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("ht"),
		Name: easygo.NewString(expressName),
	})
	commonUseComInfos = append(commonUseComInfos, &share_message.ExpressCom{
		Code: easygo.NewString("ht"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("bsky")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("bsky"),
		Name: easygo.NewString(expressName),
	})

	//D字母开头
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString(""),
		Name: easygo.NewString("D"),
	})

	expressName, _ = GetExpressNamePhone("db")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("db"),
		Name: easygo.NewString(expressName),
	})

	//E字母开头
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString(""),
		Name: easygo.NewString("E"),
	})

	expressName, _ = GetExpressNamePhone("ems")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("ems"),
		Name: easygo.NewString(expressName),
	})
	commonUseComInfos = append(commonUseComInfos, &share_message.ExpressCom{
		Code: easygo.NewString("ems"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("emsg")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("emsg"),
		Name: easygo.NewString(expressName),
	})

	////H字母开头
	//comInfos = append(comInfos, &share_message.ExpressCom{
	//	Code: easygo.NewString(""),
	//	Name: easygo.NewString("H"),
	//})
	//
	//expressName, _ = GetExpressNamePhone("ht")
	//comInfos = append(comInfos, &share_message.ExpressCom{
	//	Code: easygo.NewString("ht"),
	//	Name: easygo.NewString(expressName),
	//})
	//commonUseComInfos = append(commonUseComInfos, &share_message.ExpressCom{
	//	Code: easygo.NewString("ht"),
	//	Name: easygo.NewString(expressName),
	//})

	//J字母开头
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString(""),
		Name: easygo.NewString("J"),
	})

	expressName, _ = GetExpressNamePhone("jd")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("jd"),
		Name: easygo.NewString(expressName),
	})

	//S字母开头
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString(""),
		Name: easygo.NewString("S"),
	})

	expressName, _ = GetExpressNamePhone("sf")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("sf"),
		Name: easygo.NewString(expressName),
	})
	commonUseComInfos = append(commonUseComInfos, &share_message.ExpressCom{
		Code: easygo.NewString("sf"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("sto")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("sto"),
		Name: easygo.NewString(expressName),
	})
	commonUseComInfos = append(commonUseComInfos, &share_message.ExpressCom{
		Code: easygo.NewString("sto"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("suning")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("suning"),
		Name: easygo.NewString(expressName),
	})

	//T字母开头
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString(""),
		Name: easygo.NewString("T"),
	})

	expressName, _ = GetExpressNamePhone("tt")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("tt"),
		Name: easygo.NewString(expressName),
	})
	commonUseComInfos = append(commonUseComInfos, &share_message.ExpressCom{
		Code: easygo.NewString("tt"),
		Name: easygo.NewString(expressName),
	})

	//Y字母开头
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString(""),
		Name: easygo.NewString("Y"),
	})

	expressName, _ = GetExpressNamePhone("yt")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("yt"),
		Name: easygo.NewString(expressName),
	})
	commonUseComInfos = append(commonUseComInfos, &share_message.ExpressCom{
		Code: easygo.NewString("yt"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("yd")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("yd"),
		Name: easygo.NewString(expressName),
	})
	commonUseComInfos = append(commonUseComInfos, &share_message.ExpressCom{
		Code: easygo.NewString("yd"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("youzheng")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("youzheng"),
		Name: easygo.NewString(expressName),
	})
	commonUseComInfos = append(commonUseComInfos, &share_message.ExpressCom{
		Code: easygo.NewString("youzheng"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("yzgn")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("yzgn"),
		Name: easygo.NewString(expressName),
	})

	//Z字母开头
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString(""),
		Name: easygo.NewString("Z"),
	})

	expressName, _ = GetExpressNamePhone("zto")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("zto"),
		Name: easygo.NewString(expressName),
	})
	commonUseComInfos = append(commonUseComInfos, &share_message.ExpressCom{
		Code: easygo.NewString("zto"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("zhongyou")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("zhongyou"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("ztoky")
	comInfos = append(comInfos, &share_message.ExpressCom{
		Code: easygo.NewString("ztoky"),
		Name: easygo.NewString(expressName),
	})

	return &share_message.ExpressComInfosResult{
		Result:            easygo.NewInt32(0),
		Msg:               easygo.NewString(""),
		CommonUseComInfos: commonUseComInfos,
		ComInfos:          comInfos,
	}
}

func GetNewestExpressStatus(statusName string) *share_message.Express_Status {

	var newestStatus share_message.Express_Status
	if statusName == "" {
		return nil
	}
	switch statusName {

	case "PENDING": //待查询
		newestStatus = share_message.Express_Status_PENDING
	case "NO_RECORD": //无记录
		newestStatus = share_message.Express_Status_NO_RECORD
	case "ERROR": //查询异常
		newestStatus = share_message.Express_Status_ERROR
	case "IN_TRANSIT": //运输中
		newestStatus = share_message.Express_Status_IN_TRANSIT
	case "DELIVERING": //派送中
		newestStatus = share_message.Express_Status_DELIVERING
	case "SIGNED": //已签收
		newestStatus = share_message.Express_Status_SIGNED
	case "REJECTED": //拒签
		newestStatus = share_message.Express_Status_REJECTED
	case "PROBLEM": //疑难件
		newestStatus = share_message.Express_Status_PROBLEM
	case "INVALID": //无效件
		newestStatus = share_message.Express_Status_INVALID
	case "TIMEOUT": //超时件
		newestStatus = share_message.Express_Status_TIMEOUT
	case "FAILED": //派送失败
		newestStatus = share_message.Express_Status_FAILED
	case "SEND_BACK": //退回
		newestStatus = share_message.Express_Status_SEND_BACK
	case "TAKING": //揽件
		newestStatus = share_message.Express_Status_TAKING
	}
	return &newestStatus
}

func GetMiddleExpressStatus(remark string) *share_message.Express_Status {

	var middleStatus share_message.Express_Status
	if remark == "" {
		return nil
	} else {
		// NO_RECORD = 1; //无记录
		// ERROR = 2; //查询异常
		//以上两个状态不设置和 PENDING = 0; 待查询的值一样都是0
		//以下设置为了更好的保留物流的中间状态以及放到数据库中
		//========= 运输中,派送中,已签收,揽件,拒签页面会显示图标 其他状态默认图标=======
		if strings.Contains(remark, "疑难") ||
			strings.Contains(remark, "疑件") ||
			strings.Contains(remark, "难件") ||
			strings.Contains(remark, "无法派送") ||
			strings.Contains(remark, "无法按时派送") ||
			strings.Contains(remark, "无法送件") ||
			strings.Contains(remark, "无法按时送件") ||
			strings.Contains(remark, "无法派件") ||
			strings.Contains(remark, "无法按时派件") ||
			strings.Contains(remark, "没法派送") ||
			strings.Contains(remark, "没法按时派送") ||
			strings.Contains(remark, "没法送件") ||
			strings.Contains(remark, "没法按时送件") ||
			strings.Contains(remark, "没法派件") ||
			strings.Contains(remark, "没法按时派件") {

			middleStatus = share_message.Express_Status_PROBLEM //疑难件
			return &middleStatus

		} else if strings.Contains(remark, "派送失败") ||
			strings.Contains(remark, "送件失败") ||
			strings.Contains(remark, "派件失败") ||
			strings.Contains(remark, "失败") {

			middleStatus = share_message.Express_Status_FAILED //派送失败
			return &middleStatus

		} else if strings.Contains(remark, "无效件") ||
			strings.Contains(remark, "无效") {

			middleStatus = share_message.Express_Status_INVALID //无效件
			return &middleStatus

		} else if strings.Contains(remark, "超时件") ||
			strings.Contains(remark, "超时") ||
			strings.Contains(remark, "超过时间") ||
			strings.Contains(remark, "超过时限") {

			middleStatus = share_message.Express_Status_TIMEOUT //超时件
			return &middleStatus

		} else if strings.Contains(remark, "退回") ||
			strings.Contains(remark, "退件") ||
			strings.Contains(remark, "退到") {

			middleStatus = share_message.Express_Status_SEND_BACK //退回
			return &middleStatus

			//========= 运输中,派送中,已签收,揽件,拒签页面会显示图标 其他状态默认图标=======
		} else if strings.Contains(remark, "运输") ||
			strings.Contains(remark, "离开") ||
			strings.Contains(remark, "发往") ||
			strings.Contains(remark, "到达") ||
			strings.Contains(remark, "下一站") ||
			strings.Contains(remark, "经转") {

			middleStatus = share_message.Express_Status_IN_TRANSIT //运输中
			return &middleStatus

		} else if strings.Contains(remark, "派送") ||
			strings.Contains(remark, "派件") ||
			strings.Contains(remark, "送件") {

			middleStatus = share_message.Express_Status_DELIVERING //派送中
			return &middleStatus
		} else if strings.Contains(remark, "拒绝") ||
			strings.Contains(remark, "拒收") ||
			strings.Contains(remark, "拒签收") {

			middleStatus = share_message.Express_Status_REJECTED //拒签
			return &middleStatus

		} else if strings.Contains(remark, "已签收") ||
			strings.Contains(remark, "签收人") ||
			strings.Contains(remark, "代签") ||
			strings.Contains(remark, "签收") {

			middleStatus = share_message.Express_Status_SIGNED //已签收
			return &middleStatus
		} else if strings.Contains(remark, "揽件") ||
			strings.Contains(remark, "揽收") ||
			strings.Contains(remark, "收件") {

			middleStatus = share_message.Express_Status_TAKING //揽件
			return &middleStatus
		}
	}

	return &middleStatus
}

func GetExpressNamePhone(com string) (string, string) {

	switch com {
	case "sf":
		return "顺丰", "95338"
	case "sto":
		return "申通", "400-889-5543"
	case "yt":
		return "圆通", "95554"
	case "yd":
		return "韵达", "95546"
	case "tt":
		return "天天", "400-188-8888"
	case "ems":
		return "EMS", "11183"
	case "zto":
		return "中通", "95311"
	case "ht":
		return "百世快递（汇通）", "95320"
	case "db":
		return "德邦", "95353"
	case "jd":
		return "京东快递", "400-603-3600"
	case "zjs":
		return "宅急送", "4006-789-000"
	case "emsg":
		return "EMS国际", "11183"
	case "yzgn":
		return "邮政国内（挂号信）", "11183"
	case "ztky":
		return "中铁快运", "95572"
	case "zhongyou":
		return "中邮物流", "11183"
	case "ztoky":
		return "中通快运", "4000-270-270"
	case "youzheng":
		return "邮政快递", "11185"
	case "bsky":
		return "百世快运", "400-8856-561"
	case "suning":
		return "苏宁快递", "95315"
	default:
		return "", ""
	}
}
