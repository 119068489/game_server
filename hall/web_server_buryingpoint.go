package hall

import (
	"encoding/base64"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

/*
 埋点服务入口
*/

//累加总数
var mIncrEvent = map[int32]string{
	for_game.BP_VC_MAIN_BG_ENTER:     "MainEnterNum",
	for_game.BP_VC_MAIN_BG_CHCEKCARD: "MainReadCardNum",
	for_game.BP_VC_MAIN_DJ_FRESH:     "MainFreshNum",
	for_game.BP_VC_LZMP_BG_ENTER:     "LZMPEnterNum",
	for_game.BP_VC_LZMP_DJ_LY:        "LZMPlyNum",
	for_game.BP_VC_LZMP_DJ_SANGCHUAN: "LZMPcg",
	for_game.BP_VC_LZMP_DJ_QUXIAO:    "LZMPqx",
	for_game.BP_VC_LZMP_DJ_SZJYMP:    "LZMPjymp",
	for_game.BP_VC_SXHW_WXH_DJ_BF:    "SXHWwxhBoFang",
	for_game.BP_VC_SXHW_WXH_DJ_LT:    "SXHWwxhChat",
	for_game.BP_VC_SXHW_WXH_DJ_TX:    "SXHWwxhHead",
}

//持续事件累加
var mIncrLastTimeEvent = map[int32]string{
	for_game.BP_VC_MAIN_BG_EXIT: "MainOnLineTime",
	for_game.BP_VC_LZMP_BG_EXIT: "LZMPOnlineTime",
	//for_game.BP_VC_LZMP_DJ_EXIT: "LZMPOnlineTime",
	for_game.BP_VC_SSGD_BG_EXIT: "SSGDOnlineTime",
}

//埋点服务入口
func (self *WebHttpServer) BuryingPointEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	self.NextJobId()
	var errStr string
	defer func() {
		if errStr != "" {
			logs.Error(errStr)
			w.Write([]byte(errStr))
		}
		r.Body.Close()
	}()
	com := r.Header.Get("Common")
	if com == "" {
		errStr = "Common参数为空"
		return
	}
	by, err := base64.StdEncoding.DecodeString(com)
	if err != nil {
		errStr = "Common参数非base64格式:" + err.Error()
		return
	}
	//	logs.Info("str:", string(by))
	common := &base.Common{}
	err = common.Unmarshal(by)
	if err != nil {
		errStr = "Common参数值不对:" + err.Error()
		return
	}
	logs.Info("common:", common)
	player := for_game.GetRedisPlayerBase(common.GetUserId())
	if player == nil {
		errStr = "无效的玩家id"
		return
	}
	if player.GetToken() != common.GetToken() {
		errStr = "token校验失败"
		return
	}
	body, _ := ioutil.ReadAll(r.Body) //获取post的数据
	if len(body) == 0 {
		errStr = "Data参数为空"
		return
	}
	params, err := url.ParseQuery(string(body))
	if err != nil {
		errStr = "解析数据异常"
		return
	}

	//logs.Info("data:", params.Get("Data"))

	by1, err := base64.StdEncoding.DecodeString(params.Get("Data"))
	if err != nil {
		errStr = "Data参数非base64格式:" + err.Error()
		return
	}
	re := &client_hall.BuryingPointList{}
	err = re.Unmarshal(by1)
	if err != nil {
		errStr = "Data参数值不对:" + err.Error()
		return
	}
	logs.Info("data:", re)
	if len(re.GetData()) > 0 {
		saveData := make([]interface{}, 0)
		mData := make(map[string]int64)
		for _, d := range re.GetData() {

			data := &share_message.BuryingPointLog{
				Id:        easygo.NewInt64(for_game.NextId(for_game.TABLE_BURYING_POINT_LOG)),
				PlayerId:  easygo.NewInt64(common.GetUserId()),
				EventType: easygo.NewInt32(d.GetEventType()),
				Time:      easygo.NewInt64(d.GetTime()),
				TargetId:  easygo.NewInt64(d.GetTargetId()),
				LastTime:  easygo.NewInt64(d.GetLastTime()),
			}
			saveData = append(saveData, bson.M{"_id": data.GetId()}, data)
			if file, ok := mIncrEvent[d.GetEventType()]; ok {
				mData[file] += 1
			}
			if file, ok := mIncrLastTimeEvent[d.GetEventType()]; ok {
				mData[file] += d.GetLastTime()
			}

		}
		for_game.UpsertAll(MongoLogMgr, for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_BURYING_POINT_LOG, saveData)

		//处理redis数据
		UpdateRedisReport(mData)
	}
	errStr = "succeed"
}

//处理reids报表更新
func UpdateRedisReport(m map[string]int64) {
	vcObj := for_game.GetRedisVCBuryingPointReportObj(time.Now().Unix())
	for k, v := range m {
		vcObj.IncrFileVal(k, v)
	}
}
