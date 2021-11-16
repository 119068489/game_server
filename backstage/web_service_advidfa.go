// 广告设备监听处理

package backstage

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"net/http"
	"net/url"
	"strconv"
)

func (self *WebHttpServer) AdvIdfaEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	params := r.Form

	t := params.Get("t") //请求数据类型
	switch t {
	case "ks": //快手
		self.ApiKsIdfa(w, params)
	case "qtt": //趣头条
		self.ApiQttIdfa(w, params)
	default:
		OutputJson(w, 0, "fail", nil)
	}
}

//快手idfa上报API
//安卓：https://testapi.lemonchat.cn/advidfa?t=ks&code=__ANDROIDID2__&callback=__CALLBACK__&advid=__DID__&os=__OS__&ip=__IP__&scenesid=__CSITE__&ts=__TS__
//IOS：https://testapi.lemonchat.cn/advidfa?t=ks&code=__IDFA2__&callback=__CALLBACK__&advid=__DID__&os=__OS__&ip=__IP__&scenesid=__CSITE__&ts=__TS__
func (self *WebHttpServer) ApiKsIdfa(w http.ResponseWriter, params url.Values) {
	code := params.Get("code")         //md5后的设备码
	advid := params.Get("advid")       //广告计划id
	os := params.Get("os")             //系统类型
	ip := params.Get("ip")             //ip
	scenesid := params.Get("scenesid") //广告场景
	callback := params.Get("callback") //回调链接
	ts := params.Get("ts")             //时间毫秒

	advId, _ := strconv.ParseInt(advid, 10, 64)
	oS, _ := strconv.ParseInt(os, 10, 64)
	scenesId, _ := strconv.ParseInt(scenesid, 10, 64)
	createTime, _ := strconv.ParseInt(ts, 10, 64)

	msg := &share_message.KsPosAdvIdfa{
		CodeMd5:    easygo.NewString(code),
		CreateTime: easygo.NewInt64(createTime),
		AdvId:      easygo.NewInt64(advId),
		OsType:     easygo.NewInt32(oS),
		Ip:         easygo.NewString(ip),
		ScenesId:   easygo.NewInt32(scenesId),
		Callback:   easygo.NewString(callback),
		IsActive:   easygo.NewBool(false),
		IsRegister: easygo.NewBool(false),
		Platform:   easygo.NewString("ks"),
	}

	for_game.SavePosDeviceAdvIdfa(msg, false)

	OutputJson(w, 1, "success", nil)
}

//趣头条idfa上报API
// https://testapi.lemonchat.cn/advidfa?t=qtt&imeimd5=__IMEIMD5__&andid=__ANDROIDIDMD5__&idfa=__IDFA__&callback=__CALLBACK_URL__&advid=__CID__&os=__OS__&ip=__IP__&ts=__TSMS__&oaid=__OAID__
func (self *WebHttpServer) ApiQttIdfa(w http.ResponseWriter, params url.Values) {
	andid := params.Get("andid")       //安卓md5后的设备码
	idfa := params.Get("idfa")         //苹果idfa原值
	advid := params.Get("advid")       //广告计划id
	os := params.Get("os")             //系统类型
	ip := params.Get("ip")             //ip
	callback := params.Get("callback") //回调链接
	ts := params.Get("ts")             //时间毫秒

	advId, _ := strconv.ParseInt(advid, 10, 64)
	oS, _ := strconv.ParseInt(os, 10, 64)
	createTime, _ := strconv.ParseInt(ts, 10, 64)

	var code string
	switch os {
	case "0": //安卓
		code = andid
	case "1": //苹果
		code = for_game.Md5(idfa)
	}

	msg := &share_message.KsPosAdvIdfa{
		CodeMd5:    easygo.NewString(code),
		CreateTime: easygo.NewInt64(createTime),
		AdvId:      easygo.NewInt64(advId),
		OsType:     easygo.NewInt32(oS),
		Ip:         easygo.NewString(ip),
		Callback:   easygo.NewString(callback),
		IsActive:   easygo.NewBool(false),
		IsRegister: easygo.NewBool(false),
		Platform:   easygo.NewString("qtt"),
	}

	for_game.SavePosDeviceAdvIdfa(msg, false)

	OutputJson(w, 1, "success", nil)
}
