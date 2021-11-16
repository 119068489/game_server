package for_game

import (
	"game_server/easygo"
	"io/ioutil"
	"net/http"
)

/*
	广告终端设备行为汇报
*/

//types 1激活,2注册
func AdvIdfaPosDevCallBack(platform string, types int, eventTime int64, callbake string) {
	switch platform {
	case "ks":
		callbake = callbake + "&event_type=" + easygo.AnytoA(types) + "&event_time=" + easygo.AnytoA(eventTime/1000)
	case "qtt":
		callbake = callbake + "&op2=" + easygo.AnytoA(types) + "&opt_active_time=" + easygo.AnytoA(eventTime/1000)
	}

	res, err := http.Get(callbake)
	easygo.PanicError(err)
	_, err1 := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	easygo.PanicError(err1)
}
