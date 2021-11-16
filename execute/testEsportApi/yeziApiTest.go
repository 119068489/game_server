package main

import (
	"game_server/easygo/util"
	"game_server/for_game"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func main() {

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
