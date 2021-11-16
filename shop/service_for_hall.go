package shop

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_server"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"reflect"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

//===================================================================

type ServiceForHall struct {
	Service reflect.Value
}

func (self *ServiceForHall) RpcHall2ShopMsg(common *base.Common, reqMsg *share_message.MsgToServer) easygo.IMessage {
	logs.Debug("RpcHall2ShopMsg", reqMsg)

	return nil
}

func (self *ServiceForHall) RpcPayCallBack(common *base.Common, reqMsg *share_message.OrderID) easygo.IMessage {
	logs.Debug("RpcPayCallBack", reqMsg)
	easygo.PCall(ShopInstance.BuyCallBack, reqMsg.GetOrderId())
	return easygo.EmptyMsg
}

//func (self *ServiceForHall) RpcBlackListCallBack(ep IHallEndpoint, ctx interface{}, reqMsg *shop_hall.BlackList) easygo.IMessage {
//	logs.Debug("RpcBlackListCallBack")
//	player := PlayerMgr.LoadPlayer(reqMsg.GetPlayerId())
//	if player != nil {
//		player.BlackList = reqMsg.GetList()
//	}
//	return nil
//}

//func (self *ServiceForHall) RpcPeopleAuthCallBack(ep IHallEndpoint, ctx interface{}, reqMsg *shop_hall.AuthInfo) easygo.IMessage {
//	logs.Debug("RpcPeopleAuthCallBack")
//	player := PlayerMgr.LoadPlayer(reqMsg.GetPlayerId())
//	if player != nil {
//		*player.PeopleId = reqMsg.GetPeopleId()
//		*player.RealName = reqMsg.GetName()
//	}
//	return easygo.EmptyMsg
//}

func (self *ServiceForHall) RpcModifyPlayerMsg(common *base.Common, reqMsg *client_server.ChangePlayerInfo) easygo.IMessage {
	//商城消息中的不维护
	t := reqMsg.GetType()
	switch t {
	case 1:
		player := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
		if player != nil {
			player.SetNickName(reqMsg.GetValue1())
		}
		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
		_, e := col.UpdateAll(bson.M{"player_id": reqMsg.GetPlayerId()}, bson.M{"$set": bson.M{"nickname": reqMsg.GetValue1()}})
		closeFun()
		if e != nil && e != mgo.ErrNotFound {
			logs.Error(e)
		}

		col_order, closeFun_order := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
		_, e_sponsor_order := col_order.UpdateAll(bson.M{"sponsor_id": reqMsg.GetPlayerId()}, bson.M{"$set": bson.M{"sponsor_nickname": reqMsg.GetValue1()}})
		_, e_receiver_order := col_order.UpdateAll(bson.M{"receiver_id": reqMsg.GetPlayerId()}, bson.M{"$set": bson.M{"receiver_nickname": reqMsg.GetValue1()}})
		closeFun_order()
		if e_sponsor_order != nil && e_sponsor_order != mgo.ErrNotFound {
			logs.Error(e_sponsor_order)
		}

		if e_receiver_order != nil && e_receiver_order != mgo.ErrNotFound {
			logs.Error(e_receiver_order)
		}

		col_com, closeFun_com := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ITEM_COMMENT)
		_, e_com := col_com.UpdateAll(bson.M{"player_id": reqMsg.GetPlayerId()}, bson.M{"$set": bson.M{"nickname": reqMsg.GetValue1()}})
		closeFun_com()
		if e_com != nil && e_com != mgo.ErrNotFound {
			logs.Error(e_com)
		}
	case 2:
		player := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
		if player != nil {
			player.SetHeadIcon(reqMsg.GetValue1())
		}
		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
		_, e := col.UpdateAll(bson.M{"player_id": reqMsg.GetPlayerId()}, bson.M{"$set": bson.M{"avatar": reqMsg.GetValue1()}})
		closeFun()
		if e != nil && e != mgo.ErrNotFound {
			logs.Error(e)
		}

		col_order, closeFun_order := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
		_, e_sponsor_order := col_order.UpdateAll(bson.M{"sponsor_id": reqMsg.GetPlayerId()}, bson.M{"$set": bson.M{"sponsor_avatar": reqMsg.GetValue1()}})
		_, e_receiver_order := col_order.UpdateAll(bson.M{"receiver_id": reqMsg.GetPlayerId()}, bson.M{"$set": bson.M{"receiver_avatar": reqMsg.GetValue1()}})
		closeFun_order()
		if e_sponsor_order != nil && e_sponsor_order != mgo.ErrNotFound {
			logs.Error(e_sponsor_order)
		}

		if e_receiver_order != nil && e_receiver_order != mgo.ErrNotFound {
			logs.Error(e_receiver_order)
		}

		col_com, closeFun_com := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ITEM_COMMENT)
		_, e_com := col_com.UpdateAll(bson.M{"player_id": reqMsg.GetPlayerId()}, bson.M{"$set": bson.M{"avatar": reqMsg.GetValue1()}})
		closeFun_com()
		if e_com != nil && e_com != mgo.ErrNotFound {
			logs.Error(e_com)
		}
	case 3:
		player := for_game.GetRedisPlayerBase(reqMsg.GetPlayerId())
		if player != nil {
			player.SetSex(reqMsg.GetValue())
		}
		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
		_, e := col.UpdateAll(bson.M{"player_id": reqMsg.GetPlayerId()}, bson.M{"$set": bson.M{"sex": reqMsg.GetValue()}})
		closeFun()
		if e != nil && e != mgo.ErrNotFound {
			logs.Error(e)
		}

		col_order, closeFun_order := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
		_, e_sponsor_order := col_order.UpdateAll(bson.M{"sponsor_id": reqMsg.GetPlayerId()}, bson.M{"$set": bson.M{"sponsor_sex": reqMsg.GetValue()}})
		_, e_receiver_order := col_order.UpdateAll(bson.M{"receiver_id": reqMsg.GetPlayerId()}, bson.M{"$set": bson.M{"receiver_sex": reqMsg.GetValue()}})
		closeFun_order()
		if e_sponsor_order != nil && e_sponsor_order != mgo.ErrNotFound {
			logs.Error(e_sponsor_order)
		}

		if e_receiver_order != nil && e_receiver_order != mgo.ErrNotFound {
			logs.Error(e_receiver_order)
		}

		col_com, closeFun_com := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ITEM_COMMENT)
		_, e_com := col_com.UpdateAll(bson.M{"player_id": reqMsg.GetPlayerId()}, bson.M{"$set": bson.M{"sex": reqMsg.GetValue()}})
		closeFun_com()
		if e_com != nil && e_com != mgo.ErrNotFound {
			logs.Error(e_com)
		}
	}
	return easygo.EmptyMsg
}

func (self *ServiceForHall) RpcBsOpShopOrder(common *base.Common, reqMsg *server_server.ShopOrderRequest) easygo.IMessage {
	//后台取消订单(目前后台管理只有待支付的才能取消,待发货和待收货的暂时代码先放着)
	if reqMsg.GetTypes() == 1 {

		var orderQuery share_message.TableShopOrder = share_message.TableShopOrder{}

		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
		defer closeFun()

		eQuery := col.Find(bson.M{"_id": reqMsg.GetOrderId()}).One(&orderQuery)

		if eQuery != nil {
			logs.Error(eQuery)
			return easygo.EmptyMsg
		}
		//订单状态  0待付款 1超时 2取消 3待发货 4待收货 5已完成 6评价 7后台取消
		//0待付款取消
		if orderQuery.GetState() == for_game.SHOP_ORDER_WAIT_PAY {

			e := col.Update(bson.M{"_id": reqMsg.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_PAY},
				bson.M{"$set": bson.M{"state": for_game.SHOP_ORDER_BACKSTAGE_CANCLE, "update_time": time.Now().Unix()}})
			if e != nil && e != mgo.ErrNotFound {
				logs.Error(e)
				return easygo.EmptyMsg
			}

			if e == mgo.ErrNotFound {
				s := fmt.Sprintf("%v待付款的订单才能取消", reqMsg.GetOrderId())
				logs.Error(e)
				logs.Error(s)
				return easygo.EmptyMsg
			}

			//判断是不是拉起了微信或支付宝支付后的订单但是没有在小程序点击取消的订单
			countCle, errCle := for_game.GetPayOrderListByShopOrderId(reqMsg.GetOrderId())
			if errCle == "" {
				//执行恢复库存,有几个支付订单就做几个
				for i := 0; i < countCle; i++ {
					//恢复库存并上架判断操作
					recoverStockErr := for_game.ShopRecoverStock(reqMsg.GetOrderId())
					if recoverStockErr != "" {
						logs.Error(recoverStockErr)
						return easygo.EmptyMsg
					}
				}
			} else {
				logs.Error(errCle)
				return easygo.EmptyMsg
			}
		}

		//订单状态  0待付款 1超时 2取消 3待发货 4待收货 5已完成 6评价 7后台取消
		//3待发货 4待收货 后台可以取消
		// 3待发货 4待收货的需要退款给买家和恢复库存但不自动上架
		if orderQuery.GetState() == for_game.SHOP_ORDER_WAIT_SEND ||
			orderQuery.GetState() == for_game.SHOP_ORDER_WAIT_RECEIVE {

			//更新需要加上状态,分开来判断
			if orderQuery.GetState() == for_game.SHOP_ORDER_WAIT_SEND {
				e := col.Update(bson.M{"_id": reqMsg.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_SEND},
					bson.M{"$set": bson.M{"state": for_game.SHOP_ORDER_BACKSTAGE_CANCLE, "update_time": time.Now().Unix()}})
				if e != nil && e != mgo.ErrNotFound {
					logs.Error(e)
					return easygo.EmptyMsg
				}

				if e == mgo.ErrNotFound {
					s := fmt.Sprintf("%v待发货的订单才能取消", reqMsg.GetOrderId())
					logs.Error(e)
					logs.Error(s)
					return easygo.EmptyMsg
				}
			}

			//更新需要加上状态,分开来判断
			if orderQuery.GetState() == for_game.SHOP_ORDER_WAIT_RECEIVE {

				e := col.Update(bson.M{"_id": reqMsg.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_RECEIVE},
					bson.M{"$set": bson.M{"state": for_game.SHOP_ORDER_BACKSTAGE_CANCLE, "update_time": time.Now().Unix()}})
				if e != nil && e != mgo.ErrNotFound {
					logs.Error(e)
					return easygo.EmptyMsg
				}

				if e == mgo.ErrNotFound {
					s := fmt.Sprintf("%v待收货的订单才能取消", reqMsg.GetOrderId())
					logs.Error(e)
					logs.Error(s)
					return easygo.EmptyMsg
				}
			}

			// 自动退款给买家
			SendMsgToServerNewEx(orderQuery.GetReceiverId(),
				"RpcShopPaySeller",
				&share_message.PaySellerInfo{
					OrderId:    orderQuery.OrderId,
					Money:      easygo.NewInt32(orderQuery.Items.GetPrice() * orderQuery.Items.GetCount()),
					Sponsor_Id: orderQuery.SponsorId,
					ReceiverId: orderQuery.ReceiverId,
					PayType:    easygo.NewInt32(1)})

			easygo.Spawn(func(itemIdPara int64) {
				//取消时候修改商品表中真实和虚假的付款数
				SubPayCnt(itemIdPara)
			}, orderQuery.GetItems().GetItemId())

			//如果是点卡恢复导入库中的点卡状态
			if orderQuery.GetItems() != nil &&
				orderQuery.GetItems().GetItemType() == for_game.SHOP_POINT_CARD_CATEGORY &&
				nil != orderQuery.GetItems().GetPointCardInfos() &&
				len(orderQuery.GetItems().GetPointCardInfos()) > 0 {
				//第一次立即执行
				RecoverDataToShopPointCard(orderQuery, 0)
			}

			//商品库存恢复以及上架
			recoverStockErr := for_game.ShopRecoverStock(reqMsg.GetOrderId())
			if recoverStockErr != "" {
				logs.Error(recoverStockErr)
				return easygo.EmptyMsg
			}
		}

		//后台提交物流信息
	} else if reqMsg.GetTypes() == 2 {

		var order *share_message.TableShopOrder = &share_message.TableShopOrder{}

		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
		defer closeFun()

		e := col.Find(bson.M{"_id": reqMsg.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_SEND}).One(order)

		if e == mgo.ErrNotFound {
			logs.Error(ORDER_EXPRESS_UPLOAD_REPEATED)
			return easygo.EmptyMsg
		}

		if e != nil {
			logs.Error(e, reqMsg.GetOrderId())
			return easygo.EmptyMsg
		}

		e = col.Update(
			bson.M{"_id": reqMsg.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_SEND},
			bson.M{"$set": bson.M{"state": for_game.SHOP_ORDER_WAIT_RECEIVE,
				"receive_time": time.Now().Unix() + 7*24*3600,
				"send_time":    time.Now().Unix(),
				"update_time":  time.Now().Unix()}})

		if e == mgo.ErrNotFound {
			s := fmt.Sprintf("%v待发货的订单才能发货", reqMsg.GetOrderId())
			logs.Error(e)
			logs.Error(s)
			return easygo.EmptyMsg
		}

		if e != nil {
			logs.Error(e)
			return easygo.EmptyMsg
		}

		easygo.Spawn(func(orderPara *share_message.TableShopOrder) {

			if nil != orderPara {

				var content string = MESSAGE_TO_BUYER_SEND
				typeValue := share_message.BuySell_Type_Buyer

				ShopInstance.InsMessageNotify(
					easygo.NewString(content),
					&typeValue,
					orderPara)

				var jgContent string = MESSAGE_TO_BUYER_SEND_PUSH
				ShopInstance.JGMessageNotify(jgContent, orderPara.GetReceiverId(), orderPara.GetOrderId(), typeValue)

				//后台发货 商城订单红点推送
				SendMsgToHallClientNew([]int64{orderPara.GetReceiverId()}, "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
					OrderId: easygo.NewInt64(orderPara.GetOrderId())})
				/*	SendToPlayer(orderPara.GetReceiverId(), "RpcShopOrderNotify",
					&share_message.ShopOrderNotifyInfoWithWho{
						PlayerId: easygo.NewInt64(orderPara.GetReceiverId()),
						OrderId:  easygo.NewInt64(orderPara.GetOrderId()),
					})*/

				//后台发货 商城订单红点推送
				SendMsgToHallClientNew([]int64{orderPara.GetSponsorId()}, "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
					OrderId: easygo.NewInt64(orderPara.GetOrderId())})
				/*		SendToPlayer(orderPara.GetSponsorId(), "RpcShopOrderNotify",
						&share_message.ShopOrderNotifyInfoWithWho{
							PlayerId: easygo.NewInt64(orderPara.GetSponsorId()),
							OrderId:  easygo.NewInt64(orderPara.GetOrderId()),
						})*/
			} else {
				logs.Debug("提交物流后发通知,缺少订单")
			}
		}, order)
		//后台取得物流信息
	} else if reqMsg.GetTypes() == 5 {

		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
		defer closeFun()
		order := share_message.TableShopOrder{}
		e := col.Find(bson.M{"_id": reqMsg.GetOrderId()}).One(&order)

		if e == mgo.ErrNotFound {
			return easygo.EmptyMsg
		}

		if e != nil {
			logs.Error(e)
			return easygo.EmptyMsg
		}
		var errCode string = ""
		var expressName string = ""
		var expressPhone string = ""
		var expressBodyList []*share_message.QueryExpressBody = []*share_message.QueryExpressBody{}

		//点卡的时候可能会出现
		if "" == order.GetExpressCom() || "" == order.GetExpressCode() {
			SendToIdelOtherServer(for_game.SERVER_TYPE_HALL, "RpcSendShopOrderExpress", &share_message.QueryExpressInfosResult{
				Result:       easygo.NewInt32(0),
				Msg:          easygo.NewString(""),
				ExpressInfos: expressBodyList,
				ExpressPhone: easygo.NewString(expressPhone),
				ExpressName:  easygo.NewString(expressName),
				UserId:       easygo.NewInt64(reqMsg.GetUserId()),
			})

		}

		expressName, expressPhone = GetExpressNamePhone(order.GetExpressCom())

		expressBodyList, errCode, _ = GetExpressInfos(
			order.GetOrderId(),
			order.GetExpressCom(),
			order.GetExpressCode(),
			order.GetDeliverAddress().GetPhone(),
			order.GetReceiveAddress().GetPhone(),
			0,
		)
		//这里出错不直接通知客户端返回一个空物流打印下err,包括数据库出错
		if errCode != "" {
			if errCode == EXPRESS_QUERY_ERROR_CODE_998 {

				s := fmt.Sprintf("查询缓存数据库出错：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.GetDeliverAddress().GetPhone(),
					order.GetReceiveAddress().GetPhone())
				logs.Error(s)

			} else if errCode == EXPRESS_QUERY_ERROR_CODE_999 {

				s := fmt.Sprintf("快递查询请求快递接口出错：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.DeliverAddress.GetPhone(),
					order.ReceiveAddress.GetPhone())

				logs.Error(s)

			} else if errCode == EXPRESS_QUERY_ERROR_CODE_1 {

				s := fmt.Sprintf("快递查询快递公司错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.DeliverAddress.GetPhone(),
					order.ReceiveAddress.GetPhone())

				logs.Error(s)

			} else if errCode == EXPRESS_QUERY_ERROR_CODE_2 {

				s := fmt.Sprintf("快递查询运单号错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.DeliverAddress.GetPhone(),
					order.ReceiveAddress.GetPhone())

				logs.Error(s)

			} else if errCode == EXPRESS_QUERY_ERROR_CODE_3 {

				s := fmt.Sprintf("快递查询失败：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.DeliverAddress.GetPhone(),
					order.ReceiveAddress.GetPhone())

				logs.Error(s)

			} else if errCode == EXPRESS_QUERY_ERROR_CODE_4 {

				s := fmt.Sprintf("快递查询查不到物流信息：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.DeliverAddress.GetPhone(),
					order.ReceiveAddress.GetPhone())

				logs.Error(s)

			} else if errCode == EXPRESS_QUERY_ERROR_CODE_5 {

				s := fmt.Sprintf("快递查询寄件人或收件人手机尾号错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.DeliverAddress.GetPhone(),
					order.ReceiveAddress.GetPhone())

				logs.Error(s)

			} else {
				s := fmt.Sprintf("快递查询其他错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.DeliverAddress.GetPhone(),
					order.ReceiveAddress.GetPhone())

				logs.Error(s)
			}
		}

		SendToIdelOtherServer(for_game.SERVER_TYPE_HALL, "RpcSendShopOrderExpress", &share_message.QueryExpressInfosResult{
			Result:       easygo.NewInt32(0),
			Msg:          easygo.NewString(""),
			ExpressInfos: expressBodyList,
			ExpressPhone: easygo.NewString(expressPhone),
			ExpressName:  easygo.NewString(expressName),
			UserId:       easygo.NewInt64(reqMsg.GetUserId()),
		})
	}

	return easygo.EmptyMsg
}

//如果是点卡的时候恢复导入库中点卡的状态
func RecoverDataToShopPointCard(orderPara share_message.TableShopOrder, t time.Duration) {
	t += 2 * time.Second    //现在是2秒间隔一次，间隔多久自定义
	if t > 10*time.Second { //现在是执行10秒跳出，多久跳出自定义
		return
	}

	err := for_game.PointCardRecoverStock(orderPara)

	if err != "" {
		logs.Error(err)
		fun := func() {
			RecoverDataToShopPointCard(orderPara, t)
		}
		easygo.AfterFunc(t, fun)
	}
}
