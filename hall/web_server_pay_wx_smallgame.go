package hall

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/share_message"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

/*
	微信小程序支付模块服务入口
*/

const PAY_SUCCESS = "success"
const PAY_FAIL = "fail"

//key="yyhu2020@999."
//sign=md5(id+code+num+key)
//http://game.lemonchat.cn/getcode?id=8&code=8989&num=100&sign=009
//id="{"PlayerId":18925001,"Num":"100","PayId":1,"PayWay":1,"PaySence":1}"
func (self *WebHttpServer) WXPayEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	self.NextJobId()
	self.LogPay("web请求:WXPayEntry")
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
			self.LogPay("web请求返回:" + PAY_SUCCESS)
		}
	}()
	if !self.CheckAddress(r.RemoteAddr) {
		errCode = 1000
		errMsg = "访问过于频繁"
		return
	}
	//if !for_game.IS_FORMAL_SERVER {
	//	errCode = 1001
	//	errMsg = "测试服不能充值，请联系客服"
	//	return
	//}

	sysConfig := PSysParameterMgr.GetSysParameter(for_game.LIMIT_PARAMETER)
	if !sysConfig.GetIsRecharge() {
		errCode = 1001
		errMsg = "充值入口已关闭，请联系客服"
		return
	}

	err := r.ParseForm() //解析参数，默认是不会解析的
	if err != nil {
		errCode = 1001
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
	id := params.Get("id")
	code := params.Get("code")
	num := params.Get("num")
	str := id + code + num + "yyhu2020@999."
	newSign := for_game.Md5(str)
	if newSign != sign {
		errCode = 1002
		errMsg = "签名错误"
		return
	}
	//检查id值是否有效
	pData := &share_message.PayOrderInfo{}
	payInfo := self.ParseDecode(id)
	logs.Info("wxpay111请求参数:", string(payInfo))
	err1 := json.Unmarshal(payInfo, pData)
	if err1 != nil {
		errCode = 1003
		errMsg = "数据异常"
		return
	}
	logs.Info("wxpay请求参数:", pData)

	if pData.GetPayWay() == for_game.PAY_TYPE_SHOP {

		//定义锁数组
		var itemLockKeys []string = []string{}
		var orderLockKeys []string = []string{}
		var orderLockPayIngKeys []string = []string{}

		//取得订单的订单id和订单对应的商品id
		var bill *share_message.TableBill = &share_message.TableBill{}

		colBill, closeFunBill := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_BILLS)
		defer closeFunBill()

		colOrder, closeFunOrder := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
		defer closeFunOrder()

		var pay_gold int64 = 0

		eBill := colBill.Find(bson.M{"_id": pData.GetPayTargetId()}).One(bill)
		if eBill != nil && eBill != mgo.ErrNotFound {
			logs.Error(eBill)
			errCode = 1010
			errMsg = "操作失败"
			return
		}

		if eBill == mgo.ErrNotFound {
			order := share_message.TableShopOrder{}
			eOrder := colOrder.Find(bson.M{"_id": pData.GetPayTargetId()}).One(&order)

			if eOrder != nil && eOrder != mgo.ErrNotFound {
				logs.Error(eOrder)
				errCode = 1010
				errMsg = "操作失败,刷新重试！"
				return
			}

			if eOrder == mgo.ErrNotFound {
				logs.Error(eOrder)
				errCode = 1012
				errMsg = "订单不存在"
				return
			}

			if order.GetItems() != nil {

				if order.GetItems().GetItemType() == for_game.SHOP_POINT_CARD_CATEGORY {
					if order.GetState() == for_game.SHOP_ORDER_EVALUTE {
						errCode = 1013
						errMsg = "重复支付,刷新重试！"
						return
					}
				} else {
					if order.GetState() == for_game.SHOP_ORDER_WAIT_SEND {
						errCode = 1013
						errMsg = "重复支付,刷新重试！"
						return
					}
				}
			} else {
				errCode = 1017
				errMsg = "订单无商品信息！"
				return
			}

			if order.GetState() != for_game.SHOP_ORDER_WAIT_PAY {
				errCode = 1010
				errMsg = "操作失败,刷新重试！"
				return
			}
			tempItemLockKey := for_game.MakeRedisKey(for_game.SHOP_ITEM_PAY_MUTEX, order.GetItems().GetItemId())
			itemLockKeys = []string{tempItemLockKey}

			tempOrderLockKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_WAIT_PAY_MUTEX, order.GetOrderId())
			orderLockKeys = []string{tempOrderLockKey}

			tempOrderLockPayIngKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_PAYING_MUTEX, order.GetOrderId())
			orderLockPayIngKeys = []string{tempOrderLockPayIngKey}

			pay_gold = int64(order.GetItems().GetPrice() * order.GetItems().GetCount())
		} else {

			if bill.GetState() == for_game.SHOP_ORDER_WAIT_SEND {
				errCode = 1013
				errMsg = "重复支付,刷新重试！"
				return
			}
			if bill.GetState() != for_game.SHOP_ORDER_WAIT_PAY {
				errCode = 1010
				errMsg = "操作失败,刷新重试！"
				return
			}

			//取得各个子订单对应的商品的id
			var shopOrderList []*share_message.TableShopOrder
			//需要从订单表中取得数据
			eOrder := colOrder.Find(bson.M{"_id": bson.M{"$in": bill.GetOrderList()}}).All(&shopOrderList)
			if eOrder != nil && eOrder != mgo.ErrNotFound {
				logs.Error(eOrder)
				errCode = 1010
				errMsg = "操作失败,刷新重试！"
				return
			}

			if eOrder == mgo.ErrNotFound {
				logs.Error(eOrder)
				errCode = 1012
				errMsg = "订单不存在"
				return
			}
			for _, value := range shopOrderList {

				if value.GetItems() != nil {

					if value.GetItems().GetItemType() == for_game.SHOP_POINT_CARD_CATEGORY {
						if value.GetState() == for_game.SHOP_ORDER_EVALUTE {
							errCode = 1013
							errMsg = "存在订单重复支付,刷新重试！"
							return
						}
					} else {
						if value.GetState() == for_game.SHOP_ORDER_WAIT_SEND {
							errCode = 1013
							errMsg = "存在订单重复支付,刷新重试！"
							return
						}
					}
				} else {
					errCode = 1017
					errMsg = "存在无商品信息的订单！"
					return
				}

				if value.GetState() != for_game.SHOP_ORDER_WAIT_PAY {
					errCode = 1010
					errMsg = "操作失败,刷新重试！"
					return
				}
				tempItemLockKey := for_game.MakeRedisKey(for_game.SHOP_ITEM_PAY_MUTEX, value.GetItems().GetItemId())
				itemLockKeys = append(itemLockKeys, tempItemLockKey)

				tempOrderLockKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_WAIT_PAY_MUTEX, value.GetOrderId())
				orderLockKeys = append(orderLockKeys, tempOrderLockKey)

				tempOrderLockPayIngKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_PAYING_MUTEX, value.GetOrderId())
				orderLockPayIngKeys = append(orderLockPayIngKeys, tempOrderLockPayIngKey)
			}

			pay_gold = int64(bill.GetPrice())
		}

		//响应第三方支付回调以及用户一些异常操作、进入即给与10s的分布式过期锁
		//这里不要主动清理、要让key自动失效
		errLockFirst := easygo.RedisMgr.GetC().DoBatchRedisLockNoRetry(orderLockPayIngKeys, 10)

		//未取得锁就直接不做了
		if errLockFirst != nil {
			s := fmt.Sprintf("WXPayEntry 取得订单多key分布式不重试锁失败,redis keys is %v", orderLockPayIngKeys)
			logs.Error(s)
			logs.Error(errLockFirst)
			errCode = 1013
			errMsg = "该订单支付未处理完成,10s后刷新重试"
			return
		}

		//取得分布式锁开始1、取得订单对应的商品的分布式锁(阻塞重试）2、取得订单对应的每个订单的分布式锁(不需要重试）
		//1、取得订单对应的商品的分布式锁，此锁需要重试，直到重试次数结束提示退出
		errLock := easygo.RedisMgr.GetC().DoBatchRedisLockWithRetry(itemLockKeys, 10)
		defer easygo.RedisMgr.GetC().DoBatchRedisUnlock(itemLockKeys)

		//如果重试后还未取得锁就直接不做了
		if errLock != nil {
			s := fmt.Sprintf("WXPayEntry 取得商品多key分布式重试锁失败,redis keys is %v", itemLockKeys)
			logs.Error(s)
			logs.Error(errLock)
			errCode = 1018
			errMsg = "下单失败,刷新重试"
			return
		}

		//因为订单失效、取消等操作恢复库存这里必须加订单
		//2、取得订单的分布式锁，此锁不重试,有一个取不到就返回
		errLock2 := easygo.RedisMgr.GetC().DoBatchRedisLockNoRetry(orderLockKeys, 10)
		defer easygo.RedisMgr.GetC().DoBatchRedisUnlock(orderLockKeys)

		//未取得锁就直接不做了
		if errLock2 != nil {
			s := fmt.Sprintf("WXPayEntry 取得订单多key分布式不重试锁失败,redis keys is %v", orderLockKeys)
			logs.Error(s)
			logs.Error(errLock2)
			errCode = 1016
			errMsg = "下单失败,刷新重试"
			return
		}

		if pay_gold <= 0 {
			errCode = 1011
			errMsg = "钱是0？"
			return
		}

		pData.Amount = easygo.NewString(easygo.AnytoA(float64(pay_gold) / 100))

		//判断库存和减少库存操作
		rst := for_game.ShopCheckAndSubStock(pData.GetPayTargetId())
		if rst != "" {
			logs.Error(rst)
			errCode = 1014
			errMsg = rst
			return
		}
	}

	player := for_game.GetRedisPlayerBase(pData.GetPlayerId())
	//青少年保护模式
	if player.GetYoungPassWord() != "" {
		errCode = 1004
		errMsg = "青少年模式下无法充值"
		return
	}
	// 充值预警通知.
	pAmount := int64(easygo.AtoFloat64(pData.GetAmount()) * 100)
	easygo.Spawn(CheckWarningSMS, for_game.GOLD_CHANGE_TYPE_IN, pAmount)
	//
	//微信登陆
	if pData.GetPayType() == for_game.PAY_TYPE_WEIXIN {
		logs.Info("微信支付")
		s := make([]byte, 0)
		var fail *base.Fail
		if pData.GetPayId() == for_game.PAY_CHANNEL_YTS_WX { //不需要走小程序登录
			s, fail = PWebYunTongShangPay.RechargeOrder(self, pData)
		} else {
			s, fail = PWXApiMgr.WXAPILogin(self, code, pData) //需要小程序openId的支付
		}
		if fail != nil { //失败,还没调起小程序
			logs.Info("微信支付失败", fail.GetCode(), fail.GetReason())
			errCode = 1015
			errMsg = "第三方支付请求失败"
			//如果是商城恢复库存
			if pData.GetPayWay() == for_game.PAY_TYPE_SHOP {
				rst := for_game.ShopRecoverStock(pData.GetPayTargetId())
				if rst != "" {
					logs.Error(rst)
					errCode = 1014
					errMsg = rst
					return
				}
			}
		} else {
			if string(s) != "" {
				errCode = 1003
				errMsg = string(s)
				logs.Info("微信支付成功:", string(s))
			}
		}
	} else if pData.GetPayType() == for_game.PAY_TYPE_ZHIFUBAO {
		logs.Info("支付宝支付")
		s, fail := PAliApiMgr.AliAPILogin(self, code, pData)
		if fail != nil { //失败,还没调起小程序
			logs.Info("支付宝支付失败:", fail.GetCode(), fail.GetReason())
			errCode = 1015
			errMsg = "第三方支付请求失败"
			//如果是商城恢复库存
			if pData.GetPayWay() == for_game.PAY_TYPE_SHOP {
				rst := for_game.ShopRecoverStock(pData.GetPayTargetId())
				if rst != "" {
					logs.Error(rst)
					errCode = 1014
					errMsg = rst
					return
				}
			}
		} else {
			if string(s) != "" {
				errCode = 1003
				errMsg = string(s)
				logs.Info("支付宝支付成功", string(s))
			}
		}
	}

}
func (self *WebHttpServer) ParseDecode(data string) []byte {
	enEscapeUrl, e := url.QueryUnescape(data)
	easygo.PanicError(e)
	if len(enEscapeUrl) == 0 {
		return []byte("")
	}
	enEscapeUrl = strings.Replace(enEscapeUrl, " ", "", -1)
	decodeBytes, err := base64.StdEncoding.DecodeString(enEscapeUrl)
	easygo.PanicError(err)
	return decodeBytes
}

//微信小程序客户端取消订单
func (self *WebHttpServer) WXPayResultEntry(w http.ResponseWriter, r *http.Request) {
	self.NextJobId()
	self.LogPay("web请求:WXPayResultEntry")
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
			self.LogPay("web请求返回:" + PAY_SUCCESS)
		}
	}()
	if !self.CheckAddress(r.RemoteAddr) {
		errCode = 1000
		errMsg = "访问过于频繁"
		return
	}
	err := r.ParseForm() //解析参数，默认是不会解析的
	if err != nil {
		errCode = 1001
		errMsg = err.Error()
		return
	}
	var params url.Values
	if r.Method == "GET" {
		params = r.Form
	} else if r.Method == "POST" {
		params = r.PostForm

	}
	var sign, orderId, retMsg string
	if r.Method == "GET" {
		params = r.Form
		self.LogPay(params)
		sign = params.Get("sign")
		orderId = params.Get("orderId")
		retMsg = params.Get("retMsg")
		logs.Info("orderId:", orderId)
		logs.Info("retMsg:", retMsg)
	} else if r.Method == "POST" {
		params = r.PostForm
		// 前端那边没发提供上面的提交方式
		self.LogPay(params)
		if len(params) == 0 {
			var m map[string]string
			defer r.Body.Close()
			all, err := ioutil.ReadAll(r.Body)
			easygo.PanicError(err)
			json.Unmarshal(all, &m)
			sign = m["sign"]
			orderId = m["orderId"]
			retMsg = m["retMsg"]
		} else {
			sign = params.Get("sign")
			orderId = params.Get("orderId")
			retMsg = params.Get("retMsg")
		}
	}

	//result := params.Get("result") // result 字段暂时没有用到,后续用到在说,result:{"result":"partner=\"\"&extern_token=\"RZ41KJ0afA3h2xvVpDL2EU5N4OU5Xn01mobilegwRZ41\"&biz_type=\"\"&biz_sub_type=\"\"&trade_no=\"2020072022001461011420303814\"&app_name=\"alipay\"&display_pay_result=\"true\"&appenv=\"appid=alipay^system=ios^version=10.1.98.6010\"&success=\"true\"","resultCode":"9000","callbackUrl":"","memo":"","extendInfo":{"isDisplayResult":true}}
	str := orderId + retMsg + "yyhu2020@999."
	logs.Info("str:", str)
	newSign := for_game.Md5(str)
	logs.Info("新签名:", newSign, sign)
	if newSign != sign {
		errCode = 1002
		errMsg = "签名错误"
		return
	}
	order := for_game.GetRedisOrderObj(orderId)
	if order == nil {
		errCode = 1005
		errMsg = "无效的订单id"
		return
	}
	//支付失败时，需要关闭订单
	if retMsg == "false" {
		// 从订单里面获取渠道,然后根据渠道去校验订单
		logs.Info("前端调用支付失败,redis中的订单数据为:-------> %+v", order)
		//通知前端交易取消
		RechargeMoneyFinish(order)
		var errStr string
		switch order.GetPayChannel() {
		case for_game.PAY_CHANNEL_HUICHAO_ZFB:
			logs.Info("支付结果前端调用--------------->", "汇潮支付宝")
			////通知前端交易取消
			//RechargeMoneyFinish(order)
			//修改订单
			order.SetPayStatus(for_game.PAY_ST_CANCEL)
			order.SetStatus(for_game.ORDER_ST_CANCEL)
			// 失败以后,也去查单,避免用户在支付宝待支付页面上点击支付
			// 异步查单
			//logs.Info("去第三方查询订单状态,orderId: %s", orderId)
			//fun := func() {
			//	PWebHuiChaoPay.ReqCheckPayOrder(orderId, 1*time.Second)
			//}
			//easygo.AfterFunc(time.Second*1, fun)

		case for_game.PAY_CHANNEL_TONGLIAN:
			logs.Info("支付结果前端调用--------------->", "通联微信")
			errStr = PWebTongLianPay.CheckOrderResult(self, orderId)
		case for_game.PAY_CHANNEL_HUICHAO_WX:
			logs.Info("支付结果前端调用--------------->", "汇潮微信")
			order.SetPayStatus(for_game.PAY_ST_CANCEL)
			order.SetStatus(for_game.ORDER_ST_CANCEL)
			// 失败以后,也去查单,避免用户在支付宝待支付页面上点击支付
			// 异步查单
			//logs.Info("去第三方查询订单状态,orderId: %s", orderId)
			//fun := func() {
			//	PWebHuiChaoPay.ReqCheckPayOrder(orderId, 1*time.Second)
			//}
			//easygo.AfterFunc(time.Second*1, fun)
		case for_game.PAY_CHANNEL_TONGTONG_WX, for_game.PAY_CHANNEL_TONGTONG_ZFB:
			logs.Info("支付结果前端调用--------------->", "统统付取消支付")
			order.SetPayStatus(for_game.PAY_ST_CANCEL)
			order.SetStatus(for_game.ORDER_ST_CANCEL)
		}
		if errStr != "" {
			errCode = 1003
			errMsg = errStr
		}

		//如果是商城恢复库存
		if order.GetPayWay() == for_game.PAY_TYPE_SHOP {
			rst := for_game.ShopRecoverStock(order.GetPayTargetId())
			if rst != "" {
				logs.Error(rst)
				errCode = 1014
				errMsg = rst
				return
			}
		}
	} else {
		if order.GetPayChannel() == for_game.PAY_CHANNEL_HUICHAO_WX || order.GetPayChannel() == for_game.PAY_CHANNEL_HUICHAO_ZFB {
			// 异步查单
			logs.Info("去第三方查询订单状态,orderId: %s", orderId)
			fun := func() {
				PWebHuiChaoPay.ReqCheckPayOrder(orderId, 1*time.Second)
			}
			easygo.AfterFunc(time.Second*1, fun)
		}
	}

}
