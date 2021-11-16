package hall

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/astaxie/beego/logs"
)

/*
	通联支付回调服务入口
*/

//订单支付完成，发货处理
func (self *WebHttpServer) TongLianEntry(w http.ResponseWriter, r *http.Request) {
	self.NextJobId()
	self.LogPay("web请求:TongLianEntry")
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body) //获取post的数据
	self.LogPay(string(data))
	params := self.ParseTLData(data)
	logs.Info("params:", params)
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

	s := PWebTongLianPay.GetSourceStr(params)
	newSign := strings.ToUpper(for_game.Md5(s))
	if newSign != params.Get("sign") {
		errCode = 1003
		errMsg = "签名验证失败"
		return
	}
	orderId := params.Get("cusorderid")
	order := for_game.GetRedisOrderObj(orderId)
	if order == nil {
		errCode = 1004
		errMsg = "无效的订单id:" + orderId
		return
	}
	//已支付，或者已完成的订单不再处理
	if order.GetStatus() >= 1 || order.GetPayStatus() >= 1 {
		errCode = 1005
		errMsg = "订单已处理:" + orderId
		return
	}
	//处理支付完成逻辑
	res := false
	st := params.Get("trxstatus")
	if st == "0000" {
		order.SetPayStatus(for_game.PAY_ST_FINISH)
	} else if st == "2008" || st == "2000" { //交易处理中
		//建议每隔一段时间(10秒)查询交易
		return
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
