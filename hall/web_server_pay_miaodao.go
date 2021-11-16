package hall

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

/*
	秒到支付模块回调服务入口
*/

//检测字段的完整性
func (self *WebHttpServer) CheckParams(values url.Values) string {
	if values.Get("outTradeNo") == "" {
		return "outTradeNo 字段值不能为空"
	}
	if values.Get("tradeNo") == "" {
		return "tradeNo 字段值不能为空"
	}
	if values.Get("payStatus") == "" {
		return "payStatus 字段值不能为空"
	}
	if values.Get("amount") == "" {
		return "amount 字段值不能为空"
	}
	if values.Get("payFinishTime") == "" {
		return "payFinishTime 字段值不能为空"
	}
	if values.Get("sign") == "" {
		return "sign 字段值不能为空"
	}
	return ""
}
func (self *WebHttpServer) ParseData(data []byte) url.Values {
	values := make(url.Values)
	strData := strings.Trim(string(data), "{}")
	strList := strings.Split(strData, ",")
	for _, v := range strList {
		newList := strings.SplitN(v, ":", 2)
		if len(newList) != 2 {
			continue
		}
		key := strings.Trim(newList[0], "\"")
		val := strings.Trim(newList[1], "\"")
		values.Add(key, val)
	}
	return values
}
func (self *WebHttpServer) ParseTLData(data []byte) url.Values {
	values := make(url.Values)
	strData := strings.Trim(string(data), "{}")
	strList := strings.Split(strData, "&")
	for _, v := range strList {
		newList := strings.SplitN(v, "=", 2)
		if len(newList) != 2 {
			continue
		}
		key := strings.Trim(newList[0], "\"")
		val := strings.Trim(newList[1], "\"")
		values.Add(key, val)
	}
	return values
}

//订单支付完成，发货处理
func (self *WebHttpServer) MiaoDaoEntry(w http.ResponseWriter, r *http.Request) {
	//创建工作id
	self.NextJobId()
	self.LogPay("web请求:MiaoDaoEntry")
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body) //获取post的数据
	self.LogPay(data)
	params := self.ParseData(data)
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
			self.LogPay("MiaoDaoEntry 成功处理")
		}
	}()
	//参与签名的字段
	signValues := make(url.Values)
	signValues.Add("outTradeNo", params.Get("outTradeNo"))
	signValues.Add("tradeNo", params.Get("tradeNo"))
	signValues.Add("payStatus", params.Get("payStatus"))
	signValues.Add("amount", params.Get("amount"))

	s := PWebMiaoDaoPay.GetSourceStr(signValues)
	newSign := strings.ToUpper(for_game.Md5(s))
	if newSign != params.Get("sign") {
		errCode = 1003
		errMsg = "签名验证失败"
		return
	}
	orderId := params.Get("outTradeNo")
	order := for_game.GetRedisOrderObj(orderId)
	if order == nil {
		errCode = 1004
		errMsg = "无效的订单id:" + orderId
		return
	}
	//已支付，或者已完成的订单不再处理
	if order.GetStatus() >= for_game.ORDER_ST_FINISH || order.GetPayStatus() >= for_game.PAY_ST_FINISH {
		errCode = 1005
		errMsg = "订单已处理:" + orderId
		return
	}
	//处理支付完成逻辑
	st := params.Get("payStatus")
	if st == "FINISH" {
		order.SetPayStatus(for_game.PAY_ST_FINISH)
	} else if st == "CLOSE" {
		order.SetPayStatus(for_game.PAY_ST_CANCEL)
	}
	order.SetOverTime(time.Now().Unix())
	res := false
	//支付完成的订单发货
	if order.GetPayStatus() == for_game.PAY_ST_FINISH {
		RechargeGoldToPlayer(order.GetPlayerId(), order.GetOrderId())
		order.SetStatus(for_game.ORDER_ST_FINISH)
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
			}
		}
	}
}
