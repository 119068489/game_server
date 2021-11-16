package hall

import (
	"game_server/easygo"
	"game_server/for_game"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/astaxie/beego/logs"
)

/*
	汇聚支付回调服务入口
*/

//订单支付完成，发货处理
//4.7接口
func (self *WebHttpServer) HuiChaoEntry(w http.ResponseWriter, r *http.Request) {
	logs.Info("------------->", "订单支付完成,回调开始")
	self.NextJobId()
	self.LogPay("HuiChaoEntry 支付回调")
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body) //获取post的数据
	self.LogPay(string(data))
	logs.Info("data:", string(data))
	params, err := url.ParseQuery(string(data))
	easygo.PanicError(err)
	logs.Info("p:", params)
	errMsg := "fail"
	defer func() {
		r := recover()
		if r != nil {
			easygo.LogPanicAndStack(r, 2)
		}
		w.Write([]byte(errMsg)) //处理成功返回
		self.LogPay("HuiChaoEntry 成功结果:" + errMsg)
	}()
	if !self.CheckAddress(r.RemoteAddr) {
		errMsg = "访问过于频繁"
		return
	}
	src := "MerNo=" + PWebHuiChaoPay.MerchantNo + "&BillNo=" + params.Get("BillNo") +
		"&OrderNo=" + params.Get("OrderNo") + "&Amount=" + params.Get("Amount") + "&Succeed=" + params.Get("Succeed")
	logs.Info("签名源文:", src)
	if !PWebHuiChaoPay.VerifySign(src, params.Get("SignInfo"), for_game.HCHAOPublic) {
		logs.Info("验签失败")
		return
	}

	// todo 判断订单状态,然后执行发货.
	//开启定时器查单
	orderId := params.Get("BillNo")
	order := for_game.GetRedisOrderObj(orderId)
	if order != nil && order.GetPayStatus() == for_game.PAY_ST_FINISH {
		if order.GetStatus() == for_game.PAY_ST_FINISH {
			errMsg = "ok"
			logs.Info("订单已经处理过了:", orderId)
			return
		}
		//如果订单不是已完成状态，把订单设置成已完成
		order.SetPayStatus(for_game.PAY_ST_FINISH)
	}
	fun := func() {
		PWebHuiChaoPay.ReqCheckPayOrder(orderId, time.Second)
	}
	easygo.AfterFunc(0, fun)
	errMsg = "ok"
	return
	//============================================================
	/*
		//处理支付完成逻辑
		//res := false
		orderId := params.Get("BillNo")
		order := for_game.GetRedisOrder(orderId)
		if order == nil {
			errMsg = PAY_FAIL + ":找不到订单"
			return
		}
		st := order.GetStatus()
		if st != for_game.ORDER_ST_WAITTING {
			errMsg = "ok"
			return
		}

		//支付成功
		if params.Get("Succeed") == "88" {
			order.PayStatus = easygo.NewInt32(for_game.PAY_ST_FINISH)
		} else {
			order.PayStatus = easygo.NewInt32(for_game.PAY_ST_CANCEL)
			order.Status = easygo.NewInt32(for_game.ORDER_ST_CANCEL)
			res = false
			errMsg = "ok"
		}
		//修改订单
		for_game.RedisCreateOrder(order, false)
		//支付完成的订单发货
		if order.GetPayStatus() == for_game.PAY_ST_FINISH {
			RechargeGoldToPlayer(order.GetPlayerId(), order.GetOrderId())
			errMsg = "ok"
			res = true

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
						serverId := PlayerOnlineMgr.GetPlayerServerId(order.GetPlayerId())
						BroadCastMsgToOtherHallClient(serverId, []int64{order.GetPlayerId()}, "RpcRechargeMoneyFinish", msg)
					}
				}
			}
		}*/
}

// 代付完成,回调
// HuiChaoDFEntry 5.6接口 对接口直接给第三方返回成功.真正的订单状态逻辑处理在代付逻辑里通过查单的形式实现了
func (self *WebHttpServer) HuiChaoDFEntry(w http.ResponseWriter, r *http.Request) {
	self.NextJobId()
	self.LogPay("HuiChaoEntry 支付回调")
	merBillNo := r.PostFormValue("MerBillNo")
	cardNo := r.PostFormValue("CardNo")
	amount := r.PostFormValue("Amount")
	succeed := r.PostFormValue("Succeed")
	billNo := r.PostFormValue("BillNo")
	signInfo := r.PostFormValue("SignInfo")
	logs.Info("汇潮代付回调参数 merBillNo: %s,cardNo: %s,amount: %s,succeed: %s,billNo: %s,signInfo: %s", merBillNo, cardNo, amount, succeed, billNo, signInfo)
	errMsg := "fail"
	defer func() {
		r := recover()
		if r != nil {
			easygo.LogPanicAndStack(r, 2)
		}
		w.Write([]byte(errMsg)) //处理成功返回
		self.LogPay("HuiChaoDFEntry 成功结果:" + errMsg)
	}()
	// 通过订单号查询订单信息
	src := "MerNo=" + PWebHuiChaoPay.MerchantNo + "&MerBillNo=" + merBillNo +
		"&CardNo=" + cardNo + "&Amount=" + amount +
		"&Succeed=" + succeed + "&BillNo=" + billNo
	logs.Info("签名源文:", src)
	if !PWebHuiChaoPay.VerifySign(src, signInfo, for_game.HCHAOPublic) {
		logs.Info("验签失败")
		return
	}

	// 直接返回OK
	errMsg = "ok"
	return
}
