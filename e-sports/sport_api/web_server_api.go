// 为[浏览器]提供的API服务

package sport_api

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

type Result struct {
	Ret    int
	Reason string
	Data   interface{}
}

type WebHttpServer struct {
	Service reflect.Value
}

func NewWebHttpServer() *WebHttpServer {
	p := &WebHttpServer{}
	p.Init()
	return p
}
func (self *WebHttpServer) Init() {

}
func (self *WebHttpServer) Serve() {
	address := for_game.MakeAddress("0.0.0.0", PServerInfo.GetWebApiPort())

	logs.Info("(API 服务) 开始监听: %v", address)

	http.HandleFunc("/notice", self.YeZiApiEntry) //野子科技回调api
	err := http.ListenAndServe(address, nil)      // 第 2 个参数可以传 nil。传进去的 self 必须有实现一个 ServeHTTP 函数
	easygo.PanicError(err)
}

func (self *WebHttpServer) YeZiApiEntry(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	var reqStr string
	var sign string

	if nil != r.Header && len(r.Header) > 0 {
		reqStr = r.Header.Get("xxe-request")
		sign = r.Header.Get("xxe-sign")
	}
	defer r.Body.Close()
	b, errBody := ioutil.ReadAll(r.Body)
	if errBody != nil {
		logs.Error(errBody)
		s := fmt.Sprintf("=====野子科技回调通知解析Body失败,错误信息:%v======", errBody)
		for_game.WriteFile("ye_zi_api.log", s)
		easygo.PanicError(errBody)
	}

	//异步处理业务
	easygo.Spawn(self.YeZiApiDealSpawn, reqStr, sign, b)

	w.Write([]byte("success"))
}

//TODO 滚盘、早盘、实时数据由于后台要修改可能要考虑用分布式锁
func (self *WebHttpServer) YeZiApiDealSpawn(reqStr string, sign string, b []byte) {

	//s := fmt.Sprintf("=======YeZiApiDealSpawn野子回调开始==========:参数reqStr=%v,sign=%v,callBackBody=%v",
	//	reqStr, sign, string(b))
	//
	//for_game.WriteFile("ye_zi_api.log", s)

	var funcParam string
	var gameIdParam string
	var typeParam string
	var timestampParam string
	var eventIdParam string

	u, err := url.Parse(reqStr)
	if err == nil {
		urlParam := u.Path
		m, err1 := url.ParseQuery(urlParam)
		if err1 == nil {
			//有可能无序
			funcParam = m.Get("func")
			gameIdParam = m.Get("game_id")
			timestampParam = m.Get("timestamp")
			eventIdParam = m.Get("event_id")
			if funcParam != "" && funcParam != "live_info" {
				typeParam = m.Get("type")
			}
		} else {
			logs.Error(err1)
			s := fmt.Sprintf("======野子科技回调通知接口解析Header参数失败,错误信息:%v=====", err1)
			for_game.WriteFile("ye_zi_api.log", s)
			easygo.PanicError(err1)
		}
	} else {
		logs.Error(err)
		s := fmt.Sprintf("======野子科技回调通知接口解析Header参数失败,错误信息:%v====", err)
		for_game.WriteFile("ye_zi_api.log", s)
		easygo.PanicError(err)
	}

	//重新排序
	param := url.Values{}
	param.Set("func", funcParam)
	param.Set("game_id", gameIdParam)
	param.Set("timestamp", timestampParam)
	param.Set("event_id", eventIdParam)
	param.Set("app_secret", yeZiAppSecret)
	if funcParam != "" && funcParam != "live_info" {
		param.Set("type", typeParam)
	}
	//解密
	signLoc := GetYeZiApiSign(param)
	//认证签名正确后处理业务
	if sign == signLoc {

		//配置的不是野子科技的项目参数传入的eventId
		if !ISAppLabel(eventIdParam) {
			s := fmt.Sprintf("================传入的eventId:%v不在游戏配置列表中、直接返回=================", eventIdParam)
			for_game.WriteFile("ye_zi_api.log", s)
			return
		}

		appLabelId := for_game.EventIdToESportLabelMap[eventIdParam]
		apiOrigin := for_game.ESPORTS_API_ORIGIN_ID_YEZI
		//判断标准
		switch funcParam {
		//两组对战比赛变更
		case "game_info":
			//转化为结构体
			callBackBody := CallBackBody{}
			errBody1 := json.Unmarshal(b, &callBackBody)

			if errBody1 != nil {
				logs.Error(errBody1)
				s := fmt.Sprintf("====野子科技回调比赛详情通知将Body转化为结构体失败,错误信息:%v", errBody1)
				for_game.WriteFile("ye_zi_api.log", s)
				easygo.PanicError(errBody1)
			}
			//处理比赛详情
			CallBackDealGameDetail(appLabelId, apiOrigin, gameIdParam, eventIdParam, callBackBody.UpdateTime)
		//早盘数据变更
		case "bet_info":
			//转化为结构体
			callBackBody := CallBackBody{}
			errBody1 := json.Unmarshal(b, &callBackBody)

			if errBody1 != nil {
				logs.Error(errBody1)
				s := fmt.Sprintf("====野子科技回调早盘通知将Body转化为结构体失败,错误信息:%v", errBody1)
				for_game.WriteFile("ye_zi_api.log", s)
				easygo.PanicError(errBody1)
			}
			CallBackDealGameGuessMorn(appLabelId, apiOrigin, gameIdParam, eventIdParam, callBackBody.UpdateTime)
		//滚盘 数据变更
		case "roll_bet_info":
			//转化为结构体
			callBackBody := CallBackBody{}
			errBody1 := json.Unmarshal(b, &callBackBody)

			if errBody1 != nil {
				logs.Error(errBody1)
				s := fmt.Sprintf("====野子科技回调滚盘通知将Body转化为结构体失败,错误信息:%v", errBody1)
				for_game.WriteFile("ye_zi_api.log", s)
				easygo.PanicError(errBody1)
			}
			CallBackUseGameGuessRoll(appLabelId, apiOrigin, gameIdParam, eventIdParam)
			CallBackDealGameGuessRoll(appLabelId, apiOrigin, gameIdParam, eventIdParam, callBackBody.UpdateTime)
		//冲正回调
		case "reversal":
			//转化为结构体
			callBackBody := CallBackBody{}
			errBody1 := json.Unmarshal(b, &callBackBody)

			if errBody1 != nil {
				logs.Error(errBody1)
				s := fmt.Sprintf("====野子科技冲正回调通知将Body转化为结构体失败,错误信息:%v", errBody1)
				for_game.WriteFile("ye_zi_api.log", s)
				easygo.PanicError(errBody1)
			}
			CallBackDealReversal(appLabelId, apiOrigin, gameIdParam, eventIdParam, callBackBody.UpdateTime, callBackBody.BetId)
		//游戏实时数据
		case "live_info":

			if appLabelId == for_game.ESPORTS_LABEL_LOL {
				//转化为结构体
				callBackBody := share_message.TableESPortsLOLRealTimeData{}
				errBody1 := json.Unmarshal(b, &callBackBody)

				if errBody1 != nil {
					logs.Error(errBody1)
					s := fmt.Sprintf("====野子科技回调LOL游戏实时数据将Body转化为结构体失败,错误信息:%v", errBody1)
					for_game.WriteFile("ye_zi_api.log", s)
					easygo.PanicError(errBody1)
				}
				CallBackLOLRealTimeData(appLabelId, apiOrigin, gameIdParam, eventIdParam, timestampParam, &callBackBody)
			} else if appLabelId == for_game.ESPORTS_LABEL_WZRY {
				//转化为结构体
				callBackBody := share_message.TableESPortsWZRYRealTimeData{}
				errBody1 := json.Unmarshal(b, &callBackBody)

				if errBody1 != nil {
					logs.Error(errBody1)
					s := fmt.Sprintf("====野子科技回调WZRY游戏实时数据将Body转化为结构体失败,错误信息:%v", errBody1)
					for_game.WriteFile("ye_zi_api.log", s)
					easygo.PanicError(errBody1)
				}
				CallBackWZRYRealTimeData(appLabelId, apiOrigin, gameIdParam, eventIdParam, timestampParam, &callBackBody)
			}

		//默认不处理任何业务
		default:
		}

	} else {
		logs.Error("=====野子科技签名验证失败===")
		s := fmt.Sprintf("=====野子科技签名验证失败====,参数:func=%v,===game_id=%v,===type=%v,===timestamp=%v,===event_id=%v,===sign=%v",
			funcParam, gameIdParam, typeParam, timestampParam, eventIdParam, sign)
		for_game.WriteFile("ye_zi_api.log", s)
	}

	s := fmt.Sprintf("=======YeZiApiDealSpawn野子回调  结束==========")

	for_game.WriteFile("ye_zi_api.log", s)
}
