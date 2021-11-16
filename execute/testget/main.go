package main

import (
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/hall"
	"game_server/login"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func main() {
	//testYunTongShangPay()
	//testListInsert()
	//testGet()
	//testRegual()
	/*params := make(url.Values)
	order := "1101160470113254579521"
	sign := for_game.Md5(order + "yyhu2020@999.")
	params.Add("orderId", order)
	params.Add("sign", sign)
	resp, err := http.PostForm("http://127.0.0.1:2601/checkorder", params)
	easygo.PanicError(err)
	defer resp.Body.Close()
	result, err1 := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err1)
	logs.Info("返回结果:", string(result))*/
}

//云通商支付测试
func testYunTongShangPay() {
	initializer := hall.NewInitializer()
	defer func() { // 若是异常了,确保异步日志有成功写盘
		logger := initializer.GetBeeLogger()
		if logger != nil {
			logger.Flush()
		}
	}()
	dict := easygo.KWAT{
		"logName":  "login",
		"yamlPath": "config_login.yaml",
	}
	initializer.Execute(dict)
	hall.Initialize()
	payData := &share_message.PayOrderInfo{
		ProduceName: easygo.NewString("充值"),
		Amount:      easygo.NewString("1"),
		PayWay:      easygo.NewInt32(1),
		PlayerId:    easygo.NewInt64(1887436042),
		PayType:     easygo.NewInt32(1),
		PayId:       easygo.NewInt32(12),
	}
	resp, err := hall.PWebYunTongShangPay.RechargeOrder(nil, payData)
	if err != nil {
		logs.Error("err:", err.GetReason())
	}
	logs.Info("resp:", resp)
	time.Sleep(time.Second * 100)
}

//插入数据
func testListInsert() {
	l := []int32{1, 2, 3, 4, 5, 6, 7}
	a := int32(99)
	l2 := easygo.Insert(l, 3, a).([]int32)
	logs.Info("l:", l)
	logs.Info("l2:", l2)
}
func testYouMiCallBack() {
	testUrl := "http%3A%2F%2Fcb-api.ymapp.com%2Ftask%2Fcallback%3Fs%3DdUAaeq2p11ywtzc1VAf6xM8Vkjl6H60OI4z7QEa1I4zjqfYp90SEyuZWgId7mw0cF6i7H6yL8FeDl7z3AOtQdbSPzI0wBlTzDUNyebFJFmvwVudQ7z0byFJdFmfKRFue2EtdFSfp7alM2OXsvNgLZvmgqlaL98qAmD%26d%3D3CE2FEE7-D424-40C3-893E-581AC9CF0B1B%26i%3D400.27"
	a, _ := url.PathUnescape(testUrl)
	logs.Info("urlL:", a)
	Url, err := url.Parse(a)
	if err != nil {
		logs.Error("GET请求解析url错误:\r\n%v", err)
	}
	resp, err := http.Get(Url.String())
	data, err := ioutil.ReadAll(resp.Body)
	logs.Info("==============" + string(data))
	if err != nil {
		logs.Error("GET请求失败,错误信息:\r\n%v", err)

	} else {
		_, err := util.JsonDecode(([]byte)(data))

		if nil == err {

		} else {
			logs.Error("GET请求失败,错误信息:\r\n%v", err)

		}

	}
}
func testGet() {
	//testURL := "https://newtest.lemonchat.cn/click"
	testURL := "http://192.168.150.233:4501/youmi"
	param := url.Values{}
	param.Set("source", "youmi")
	param.Set("appid", "com.silbermond.tktalks")
	param.Set("idfa", "8306904B-CD22-4190-962F-124E5B033F0F")
	param.Set("imei", "")
	param.Set("callback_url", "http://cb-api.ymapp.com/task/callback?s=dUAaeq2p11ywtzc1VAf6xM8Vkjl6H60OI4z7QEa1I4zjqfYp90SEyuZWgId7mw0cF6i7H6yL8FeDl7z3AOtQdbSPzI0wBlTzDUNyebFJFmvwVudQ7z0byFJdFmfKRFue2EtdFSfp7alM2OXsvNgLZvmgqlaL98qAmD&d=3CE2FEE7-D424-40C3-893E-581AC9CF0B1B&i=400.27")
	Url, err := url.Parse(testURL)
	if err != nil {
		logs.Error("GET请求解析url错误:\r\n%v", err)
	}
	////如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = param.Encode()
	//发送请求
	resp, err := http.Get(Url.String())
	data, err := ioutil.ReadAll(resp.Body)
	logs.Info("==============" + string(data))
	if err != nil {
		logs.Error("GET请求失败,错误信息:\r\n%v", err)

	} else {

		_, err := util.JsonDecode(([]byte)(data))

		if nil == err {

		} else {
			logs.Error("GET请求失败,错误信息:\r\n%v", err)

		}
	}
}
func testRegual() {
	//逗号 顿号 空格 感叹号 省略号 问号 书名号 - 双引号 冒号 分号
	a := for_game.CheckTopic(`#！……？《我》-“1aA":;"#`)
	if a != nil {
		logs.Info("返回结果:", a)
	}
}
func testRedis() {
	initializer := login.NewInitializer()
	defer func() { // 若是异常了,确保异步日志有成功写盘
		logger := initializer.GetBeeLogger()
		if logger != nil {
			logger.Flush()
		}
	}()
	dict := easygo.KWAT{
		"logName":  "login",
		"yamlPath": "config_login.yaml",
	}
	initializer.Execute(dict)
	fun := func(n int) {
		order := for_game.GetRedisOrderObj("1100159831794551122787")
		order.SetStatus(int32(n))
	}
	login.Initialize()
	for i := 0; i < 10000; i++ {
		go fun(i)
		time.Sleep(3 * time.Second)
	}

	time.Sleep(99999999999999)
	//a := []int64{}
	//err := easygo.RedisMgr.GetC().LPush("testlist", a)
	//easygo.PanicError(err)
	//b, err := easygo.RedisMgr.GetC().LRange("testlist", 0, -1)
	//easygo.PanicError(err)
	//var data []int64
	//for_game.InterfersToInt64s(b, &data)
	//logs.Info(data)
	//obj := for_game.GetRedisPlayerBase(1887436027)
	//if obj == nil {
	//	logs.Info("无效的obj")
	//	return
	//}
	//player := obj.GetRedisPlayerBase()
	//logs.Info("player11:", player.GetPlayerSetting())
	//logs.Info("player11:", player.GetLabel())
	//logs.Info("player11:", player.GetCallInfo())
	//logs.Info("player11:", player.GetTeamIds())
	//logs.Info("obj:", obj)
	//obj.SetStatus(2)
	//time.Sleep(11 * time.Second)
	//logs.Info("status:", obj.GetPhoto())
	//obj1 := for_game.GetRedisPlayerBase(1887436027)
	//logs.Info("obj1:", obj1)
	//time.Sleep(11 * time.Second)
	//obj2 := for_game.GetRedisPlayerBase(1887436027)
	//logs.Info("obj2:", obj2)
	//time.Sleep(99999999999999)
	//orderObj1 := for_game.GetRedisOrderObj("1100159831794551122787")
	//if orderObj1 == nil {
	//	logs.Info("无效的空订单")
	//	return
	//}
	//logs.Info(orderObj1)
	//order := orderObj.GetRedisOrder()
	//order.OrderId = easygo.NewString("1100159831794551122789")
	//logs.Info("oreder:", order)
	//newObj := for_game.CreateRedisOrder(order)
	//logs.Info(newObj)
	//newObj.SetStatus(for_game.ORDER_ST_FINISH)
	//newObj.SaveToMongo()
	//orderObj.SaveOneRedisDataToMongo("Status", for_game.ORDER_ST_CANCEL)
}
func yeziApiTest() {

	////注意10.3的时候必须调用10.1权限
	//GetRollRole()

	////请求地址
	yeziURL := "http://47.91.198.97:8082/"
	//app_key := "pS8vq5uw"
	//timestamp := time.Now().Unix()
	//app_secret := "ejeGhQ8Qu28Maol2QoIHKUSAPFsfjLfE"
	//
	//sign := for_game.Md5(app_key := "pS8vq5uw")

	//初始化参数
	param := url.Values{}
	var timeStr string = strconv.FormatInt(time.Now().Unix(), 10)

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("app_key", "pS8vq5uw")
	param.Set("timestamp", timeStr)
	param.Set("event_id", "24")
	param.Set("app_secret", "ejeGhQ8Qu28Maol2QoIHKUSAPFsfjLfE")
	////项目列表开始//////////////////////////////
	//param.Set("resource", "event")
	//param.Set("func", "events")
	////项目列表结束============================

	//赛事列表开始//////////////////////////////
	////DOTA2
	//param.Set("resource", "event")
	//param.Set("func", "matches")
	//param.Set("event_id", "24")
	//CSGO
	//param.Set("resource", "event")
	//param.Set("func", "matches")
	//param.Set("event_id", "205")
	//赛事列表结束============================

	////特殊竞猜开始//////////////////////////////
	//param.Set("resource", "specialGuess")
	//param.Set("func", "getList")
	////特殊竞猜结束============================

	//返回未开始的比赛列表开始//////////////////////////////
	param.Set("resource", "game")
	param.Set("func", "lists")
	//返回未开始的比赛列表结束============================

	////6.1返回比赛详情(两组对决)开始//////////////////////////////
	//param.Set("resource", "game")
	//param.Set("func", "info")
	//param.Set("id", "51833")
	////6.1返回比赛详情(两组对决)结束============================

	////6.4返回比赛详情(两组对决)开始//////////////////////////////
	//param.Set("resource", "game")
	//param.Set("func", "info_extend")
	//param.Set("game_id", "51833")
	////6.4返回比赛详情(两组对决)结束============================

	////8.1两队胜败统计接口开始//////////////////////////////
	//param.Set("resource", "game")
	//param.Set("func", "continueWin")
	//param.Set("game_id", "51833")
	////8.1两队胜败统计接口结束============================

	////8.2 两队天敌克制统计接口开始//////////////////////////////
	//param.Set("resource", "game")
	//param.Set("func", "teamOpponent")
	//param.Set("game_id", "51833")
	////8.2 两队天敌克制统计接口结束============================

	////8.2 两队天敌克制统计接口开始//////////////////////////////
	//param.Set("resource", "game")
	//param.Set("func", "teamOpponent")
	//param.Set("game_id", "51833")
	////8.2 两队天敌克制统计接口结束============================

	////9.0 比赛动态信息接口（早盘）开始//////////////////////////////
	//param.Set("resource", "dynamic")
	//param.Set("func", "game")
	//param.Set("game_id", "51833")
	////9.0 比赛动态信息接口（早盘）结束============================

	////TODO
	////10.3 获取滚盘比赛动态信息开始//////////////////////////////
	//param.Set("resource", "roll")
	//param.Set("func", "get_info")
	//param.Set("game_id", "51852")
	////10.3 获取滚盘比赛动态信息结束============================

	Url, err := url.Parse(yeziURL)
	if err != nil {
		logs.Error("野子科技GET请求解析url错误:\r\n%v", err)
	}
	////如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = param.Encode()
	sign := for_game.Md5(Url.RawQuery)

	//重新设置url
	//初始化参数
	param1 := url.Values{}

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param1.Set("app_key", "pS8vq5uw")
	param1.Set("timestamp", timeStr)
	param1.Set("event_id", "24")

	////项目列表开始//////////////////////////////
	//param1.Set("resource", "event")
	//param1.Set("func", "events")
	////项目列表结束============================

	//赛事列表开始//////////////////////////////
	////DOTA2
	//param1.Set("resource", "event")
	//param1.Set("func", "matches")
	//param1.Set("event_id", "24")
	//CSGO
	//param1.Set("resource", "event")
	//param1.Set("func", "matches")
	//param1.Set("event_id", "205")
	//赛事列表结束============================

	//返回未开始的比赛列表开始//////////////////////////////
	param1.Set("resource", "game")
	param1.Set("func", "lists")
	//返回未开始的比赛列表结束============================

	////特殊竞猜开始//////////////////////////////
	//param1.Set("resource", "specialGuess")
	//param1.Set("func", "getList")
	////特殊竞猜结束============================

	////6.1返回比赛详情(两组对决)开始//////////////////////////////
	//param1.Set("resource", "game")
	//param1.Set("func", "info")
	//param1.Set("id", "51833")
	////6.1返回比赛详情(两组对决)结束============================

	////6.4返回比赛详情(两组对决)开始//////////////////////////////
	//param1.Set("resource", "game")
	//param1.Set("func", "info_extend")
	//param1.Set("game_id", "51833")
	////6.4返回比赛详情(两组对决)结束============================

	////8.1两队胜败统计接口开始//////////////////////////////
	//param1.Set("resource", "game")
	//param1.Set("func", "continueWin")
	//param1.Set("game_id", "51833")
	////8.1两队胜败统计接口结束============================

	////8.2 两队天敌克制统计接口开始//////////////////////////////
	//param1.Set("resource", "game")
	//param1.Set("func", "teamOpponent")
	//param1.Set("game_id", "51833")
	////8.2 两队天敌克制统计接口结束============================

	////9.0 比赛动态信息接口（早盘）开始//////////////////////////////
	//param1.Set("resource", "dynamic")
	//param1.Set("func", "game")
	//param1.Set("game_id", "51833")
	////9.0 比赛动态信息接口（早盘）结束============================

	////TODO
	////10.3 获取滚盘比赛动态信息开始//////////////////////////////
	//param1.Set("resource", "roll")
	//param1.Set("func", "get_info")
	//param1.Set("game_id", "51852")
	////10.3 获取滚盘比赛动态信息结束============================

	param1.Set("sign", sign)
	//发送请求
	data, err := Get(yeziURL, param1)
	logs.Info("==============" + string(data))
	if err != nil {
		logs.Error("野子科技GET请求失败,错误信息:\r\n%v", err)

	} else {

		_, err := util.JsonDecode(([]byte)(data))

		if nil == err {

		} else {
			logs.Error("野子科技GET请求失败,错误信息:\r\n%v", err)

		}

	}
}

func GetRollRole() {
	////请求地址
	yeziURL := "http://47.91.198.97:8082/"

	//初始化参数
	param := url.Values{}
	var timeStr string = strconv.FormatInt(time.Now().Unix(), 10)

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("app_key", "pS8vq5uw")
	param.Set("timestamp", timeStr)
	param.Set("event_id", "24")
	param.Set("app_secret", "ejeGhQ8Qu28Maol2QoIHKUSAPFsfjLfE")

	//10.1 使用滚盘开始//////////////////////////////
	param.Set("resource", "roll")
	param.Set("func", "use_roll")
	param.Set("game_id", "51852")
	//10.1 使用滚盘结束============================

	Url, err := url.Parse(yeziURL)
	if err != nil {
		logs.Error("野子科技GET请求解析url错误:\r\n%v", err)
	}
	////如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = param.Encode()
	sign := for_game.Md5(Url.RawQuery)

	//重新设置url
	//初始化参数
	param1 := url.Values{}

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param1.Set("app_key", "pS8vq5uw")
	param1.Set("timestamp", timeStr)
	param1.Set("event_id", "24")

	//10.1 使用滚盘开始//////////////////////////////
	param1.Set("resource", "roll")
	param1.Set("func", "use_roll")
	param1.Set("game_id", "51852")
	//10.1 使用滚盘结束============================

	param1.Set("sign", sign)
	//发送请求
	data, err := Get(yeziURL, param1)
	logs.Info("==============" + string(data))
	if err != nil {
		logs.Error("野子科技GET请求失败,错误信息:\r\n%v", err)

	} else {

		_, err := util.JsonDecode(([]byte)(data))

		if nil == err {

		} else {
			logs.Error("野子科技GET请求失败,错误信息:\r\n%v", err)

		}

	}
}

// get 网络请求
func Get(apiURL string, params url.Values) (rs []byte, err error) {
	var Url *url.URL
	Url, err = url.Parse(apiURL)
	if err != nil {
		logs.Error("野子科技GET请求解析url错误:\r\n%v", err)
		return nil, err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	resp, err := http.Get(Url.String())
	logs.Info(resp.Body)
	if err != nil {
		logs.Error("野子科技GET请求err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
