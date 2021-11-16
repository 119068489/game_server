package hall

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego/logs"
)

/*
	汇聚支付回调服务入口
*/

//订单支付完成，发货处理
//4.7接口
func (self *WebHttpServer) HuiJuEntry(w http.ResponseWriter, r *http.Request) {
	self.NextJobId()
	self.LogPay("HuiJuEntry 支付回调")
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body) //获取post的数据
	self.LogPay(string(data))
	var errMsg string
	defer func() {
		r := recover()
		if r != nil {
			easygo.LogPanicAndStack(r, 2)
		}
		w.Write([]byte(errMsg)) //处理成功返回
		self.LogPay("HuiJuEntry 成功结果:" + errMsg)
	}()
	params, err := PWebHuiJuPay.CheckData(data)
	if err != nil {
		errMsg = PAY_FAIL + ":" + err.GetReason()
		return
	}
	//处理支付完成逻辑
	res := false
	orderId := params.GetString("mch_order_no")
	order := for_game.GetRedisOrderObj(orderId)
	if order == nil {
		errMsg = PAY_FAIL + ":找不到订单"
		return
	}
	order.SetExternalNo(params.GetString("jp_order_no"))
	st := params.GetString("order_status")
	if st == HJ_ORDER_STATUS_SUCCESS {
		order.SetPayStatus(for_game.PAY_ST_FINISH)
	} else if st == HJ_ORDER_STATUS_IN_PROCESSING { //交易处理中
		return
	} else {
		order.SetPayStatus(for_game.PAY_ST_CANCEL)
		order.SetStatus(for_game.ORDER_ST_CANCEL)
		//存储外部订单号
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
			}
		}
	}
}
