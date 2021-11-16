package hall

import (
	"game_server/easygo"
	"game_server/for_game"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"strings"
)

/*
	鹏聚支付回调服务入口
*/

//订单支付完成，发货处理
func (self *WebHttpServer) PengJuEntry(w http.ResponseWriter, r *http.Request) {
	self.NextJobId()
	self.LogPay("web请求:TongLianEntry")
	defer r.Body.Close()
	var errCode int = 0
	var errMsg string

	defer func() {
		r := recover()
		if r != nil {
			easygo.LogPanicAndStack(r, 2)
		}
		if errCode != 0 {
			self.SendErrorCode(w, errCode, errMsg) //处理失败返回
		} else {
			w.Write([]byte(PAY_SUCCESS)) //处理成功返回
			self.LogPay("TongLianEntry 成功处理")
		}
	}()
	if !self.CheckAddress(r.RemoteAddr) {
		errCode = 1003
		errMsg = "访问过于频繁"
		return
	}
	data, _ := ioutil.ReadAll(r.Body) //获取post的数据
	self.LogPay(string(data))
	params := self.ParseTLData(data)
	logs.Info("params:", params)

	s := PWebTongLianPay.GetSourceStr(params)
	newSign := strings.ToUpper(for_game.Md5(s))
	if newSign != params.Get("sign") {
		errCode = 1003
		errMsg = "签名验证失败"
		return
	}
}
