package hall

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/astaxie/beego/logs"
)

/*
	统统付回调服务入口
*/

//订单支付完成，发货处理
func (self *WebHttpServer) TongTongPayEntry(w http.ResponseWriter, r *http.Request) {
	logs.Info("------------->", "订单支付完成,回调开始")
	self.NextJobId()
	self.LogPay("TongTongPayEntry 支付回调")
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body) //获取post的数据
	//data := for_game.GBKToUtf8(string(bData))
	self.LogPay(string(data))
	logs.Info("data:", string(data))

	var errCode int = 0
	var errMsg string
	defer func() {
		logs.Info("处理结果:", errCode, errMsg)
		r := recover()
		if r != nil {
			easygo.LogPanicAndStack(r, 2)
		}
		if errCode != 0 {
			self.LogPay("TongTongPayEntry 处理结果:" + errMsg)
			self.SendErrorCode(w, errCode, errMsg) //处理失败返回
		} else {
			self.LogPay("TongTongPayEntry 处理结果:" + errMsg)
			w.Write([]byte("SUCCESS")) //处理成功返回
		}
	}()
	signData, err := url.ParseQuery(string(data))
	if err != nil {
		errCode = 1001
		errMsg = "解析数据异常"
		return
	}
	if !self.CheckAddress(r.RemoteAddr) {
		errCode = 1002
		errMsg = "访问过于频繁"
		return
	}
	if !PWebTongTongPay.CheckSign(signData) {
		errCode = 1003
		errMsg = "验证签名错误"
		return
	}
	orderId := signData.Get("app_id")
	order := for_game.GetRedisOrderObj(orderId)
	if order == nil {
		errCode = 1004
		errMsg = "无效的订单id:" + orderId
		return
	}
	//已支付，或者已完成的订单不再处理
	if order.GetStatus() >= 1 || order.GetPayStatus() >= 1 {
		errCode = 0
		logs.Error("订单已处理:" + orderId)
		errMsg = "SUCCESS"
		return
	}
	errMsg = "SUCCESS"
	//处理支付完成逻辑
	res := false
	code := signData.Get("errorcode")
	if code == "0000" {
		order.SetPayStatus(for_game.PAY_ST_FINISH)
	} else {
		order.SetPayStatus(for_game.PAY_ST_CANCEL)
		order.SetStatus(for_game.ORDER_ST_CANCEL)
		res = false
	}
	//支付完成的订单发货
	if order.GetPayStatus() == for_game.PAY_ST_FINISH {
		RechargeGoldToPlayer(order.GetPlayerId(), order.GetOrderId())
		res = true
	}
	//如果充值，通知前端充值结果
	ep := ClientEpMgr.LoadEndpointByPid(order.GetPlayerId())
	if order.GetSourceType() == for_game.GOLD_TYPE_CASH_IN {
		msg := &share_message.RechargeFinish{
			Amount:        easygo.NewInt64(order.GetAmount()),
			TradeNo:       easygo.NewString(order.GetOrderId()),
			PayFinishTime: easygo.NewInt64(order.GetOverTime()),
			Result:        easygo.NewBool(res),
		}
		if ep != nil {
			ep.RpcRechargeMoneyFinish(msg)
			logs.Info("通知前端支付结果:", msg)
		} else {
			if PlayerOnlineMgr.CheckPlayerIsOnLine(order.GetPlayerId()) {
				SendMsgToHallClientNew(order.GetPlayerId(), "RpcRechargeMoneyFinish", msg)
				logs.Info("通知前端支付结果:", msg)
			}
		}
	}

}

//查询订单状态
func (self *WebHttpServer) CheckOrderEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	self.NextJobId()
	self.LogPay("web请求:CheckOrderEntry")
	var errCode int = 0
	var errMsg string
	defer func() {
		r := recover()
		if r != nil {
			easygo.LogPanicAndStack(r, 2)
		}
		if errCode != 0 {
			self.SendErrorCode(w, errCode, errMsg) //处理失败返回
		}
	}()
	if !self.CheckAddress(r.RemoteAddr) {
		errCode = 1001
		errMsg = "访问过于频繁"
		return
	}

	err := r.ParseForm() //解析参数，默认是不会解析的
	if err != nil {
		errCode = 1002
		errMsg = err.Error()
		return
	}
	var params url.Values
	if r.Method == "GET" {
		params = r.Form
	} else if r.Method == "POST" {
		params = r.PostForm

	}
	self.LogPay(params)
	sign := params.Get("sign")
	id := params.Get("orderId")
	str := id + "yyhu2020@999."
	newSign := for_game.Md5(str)
	if newSign != sign {
		errCode = 1003
		errMsg = "签名错误"
		return
	}
	//检查id值是否有效
	logs.Info("wxpay请求参数:", params)

	orderObj := for_game.GetRedisOrderObj(id)
	if orderObj == nil {
		errCode = 1004
		errMsg = "无效的订单id"
		return
	}
	//返回订单信息给前端
	order := orderObj.GetRedisOrder()
	if order.GetStatus() == for_game.ORDER_ST_FINISH {
		errCode = 2000
		errMsg = easygo.AnytoA(order.GetPayTargetId())
	} else if order.GetStatus() == for_game.ORDER_ST_WAITTING {
		errCode = 2001
		errMsg = "待支付"
	} else {
		errCode = 2002
		errMsg = "已取消"
	}
}
