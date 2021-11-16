package hall

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/for_game"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"time"
)

/*
	云通商支付回调服务入口
*/

type YTSResp struct {
	Amount     float64 `json:"amount"`
	SysOrderID string  `json:"sys_order_id"`
	CreateTime string  `json:"create_time"`
	Sign       string  `json:"sign"`
	Type       string  `json:"type"`
	OrderID    string  `json:"order_id"`
	AppID      string  `json:"app_id"`
	PayTime    string  `json:"pay_time"`
}

//订单支付完成，回调发货处理:处理完成则返回success，失败则其他
func (self *WebHttpServer) YunTongShangEntry(w http.ResponseWriter, r *http.Request) {
	logs.Info("------------->", "YunTongShangEntry 订单支付完成,回调开始")
	self.NextJobId()
	self.LogPay("YunTongShangEntry 支付回调")
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body) //获取post的数据
	self.LogPay(string(data))
	logs.Info("data:", string(data))
	params := YTSResp{}
	errMsg := "fail"
	err := json.Unmarshal(data, &params)
	if err != nil {
		logs.Error("数据解析错误")
		errMsg = "数据解析错误"
		return
	}
	logs.Info("params:", params)

	defer func() {
		r := recover()
		if r != nil {
			easygo.LogPanicAndStack(r, 2)
		}
		w.Write([]byte(errMsg)) //处理成功返回
		self.LogPay("YunTongShangEntry 成功结果:" + errMsg)
	}()
	if !PWebYunTongShangPay.VerifySign(params) {
		errMsg = "fail:签名校验失败"
		easygo.NewString("签名校验失败")
		return
	}
	//处理支付完成逻辑
	method := params.Type
	orderId := params.OrderID
	order := for_game.GetRedisOrderObj(orderId)
	if order == nil {
		errMsg = PAY_FAIL + ":找不到订单"
		return
	}
	if method == "payment.success" {
		//支付成功回调
		if order.GetPayStatus() == for_game.PAY_ST_FINISH {
			if order.GetStatus() == for_game.ORDER_ST_FINISH {
				logs.Error("已经处理过的订单")
				errMsg = "success"
				return
			}
		}
	}
	//启动查询定时器处理订单
	fun := func() {
		PWebYunTongShangPay.ReqCheckPayOrder(orderId, time.Second)
	}
	easygo.AfterFunc(0, fun)
	errMsg = "success"
}
