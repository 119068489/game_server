package for_game

import (
	"fmt"
	"game_server/easygo"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

//商城购买的时候检测库存和去库存
func ShopCheckAndSubStock(orderId int64) string {

	//取得订单的订单id和订单对应的商品id
	var bill *share_message.TableBill = &share_message.TableBill{}

	colBill, closeFunBill := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_BILLS)
	defer closeFunBill()

	colOrder, closeFunOrder := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_ORDERS)
	defer closeFunOrder()

	colItem, closeFunItem := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_ITEMS)
	defer closeFunItem()

	//点卡的类型,支付完成后处理比较安全,支付开始只处理上面中商品的库存
	//colCard, closeFunCard := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_POINT_CARD)
	//defer closeFunCard()

	eBill := colBill.Find(bson.M{"_id": orderId}).One(bill)
	if eBill != nil && eBill != mgo.ErrNotFound {
		logs.Error(eBill)
		return "操作失败"
	}

	//说明不是从购物车结算过来的支付操作
	if eBill == mgo.ErrNotFound {

		var shopOrder share_message.TableShopOrder = share_message.TableShopOrder{}
		//需要从订单表中取得数据

		eOrder := colOrder.Find(bson.M{"_id": orderId}).One(&shopOrder)
		if eOrder != nil && eOrder != mgo.ErrNotFound {
			logs.Error(eOrder)
			return "操作失败"
		}

		if eOrder == mgo.ErrNotFound {
			logs.Error(eOrder)
			return "订单不存在"
		}

		if shopOrder.GetItems() != nil {

			if shopOrder.GetItems().GetItemType() == SHOP_POINT_CARD_CATEGORY {
				if shopOrder.GetState() == SHOP_ORDER_EVALUTE {
					return "重复支付,刷新重试！"
				}
			} else {
				if shopOrder.GetState() == SHOP_ORDER_WAIT_SEND {
					return "重复支付,刷新重试！"
				}
			}
		} else {
			return "订单无商品信息！"
		}

		if shopOrder.GetState() != SHOP_ORDER_WAIT_PAY {
			return "操作失败,刷新重试！"
		}

		//判断库存和减少库存开始
		//1判断商品是否下架
		var item share_message.TableShopItem = share_message.TableShopItem{}

		eItem := colItem.Find(bson.M{"_id": shopOrder.GetItems().GetItemId(), "state": SHOP_ITEM_SALE}).One(&item)

		if eItem != nil && eItem != mgo.ErrNotFound {
			logs.Error(eItem)
			return "操作失败"
		}

		if eItem == mgo.ErrNotFound {
			s := fmt.Sprintf("订单%v对应的商品%v已经下架", shopOrder.GetOrderId(), shopOrder.GetItems().GetItemId())
			logs.Error(s)
			return "商品已经下架"
		}

		//2库存判断
		if item.GetStockCount() < shopOrder.GetItems().GetCount() {
			s := fmt.Sprintf("该订单%v对应的商品%v库存不足", shopOrder.GetOrderId(), item.GetItemId())
			logs.Error(s)
			return "商品库存不足"
		}

		//如果是点卡 再判断一次实际导入库中的库存是否足够
		if item.GetType() != nil && item.GetType().GetType() == SHOP_POINT_CARD_CATEGORY {

			//通过物品的个数取得卡密个数
			pointCardList := GetPointCardByBuyInfos(item.GetPlayerAccount(), item.GetPointCardName(), shopOrder.GetItems().GetCount())
			if nil != pointCardList && (len(pointCardList) == 0 || int32(len(pointCardList)) < shopOrder.GetItems().GetCount()) {
				s := fmt.Sprintf("该订单%v对应的商品%v库存不足", shopOrder.GetOrderId(), item.GetItemId())
				logs.Error(s)
				return "商品库存不足"
			}
		}

		//3减少库存,只减少商品表中的库存,导入库中的库存在支付成功后才会做更新和绑定到订单
		s := fmt.Sprintf("订单%v 当前物品%v原库存%v, 该扣除数量%v",
			shopOrder.GetOrderId(),
			item.GetItemId(),
			item.GetStockCount(),
			shopOrder.GetItems().GetCount())
		WriteFile("shop_order.log", s)
		eItem2 := colItem.Update(bson.M{"_id": item.GetItemId()}, bson.M{"$inc": bson.M{"stock_count": -shopOrder.GetItems().GetCount()}})

		if eItem2 != nil && eItem2 != mgo.ErrNotFound {
			logs.Error(eItem2)
			return "操作失败"
		}

		if eItem2 == mgo.ErrNotFound {
			s := fmt.Sprintf("订单%v对应的商品%v不存在", shopOrder.GetOrderId(), item.GetItemId())
			logs.Error(s)
			return "商品不存在"
		}

		//点卡的类型,支付完成后处理比较安全,支付开始只处理上面中商品的库存
		////如果是点卡,减少导入库中的库存,更新和绑定到订单
		//if item.GetType() != nil && item.GetType().GetType() == SHOP_POINT_CARD_CATEGORY {
		//
		//	s = fmt.Sprintf("发起支付的时候 点卡订单%v 绑定点卡信息到订单开始", shopOrder.GetOrderId())
		//
		//	order := share_message.TableShopOrder{}
		//	var itemPointCards []*share_message.ShopPointCardInfo = make([]*share_message.ShopPointCardInfo, 0)
		//	var tempCardIds []int64 = []int64{}
		//
		//	eOrderQuery := colOrder.Find(bson.M{"_id": shopOrder.GetOrderId()}).One(&order)
		//
		//	if eOrderQuery != nil && eOrderQuery != mgo.ErrNotFound {
		//		logs.Error(eOrderQuery)
		//		return "操作失败"
		//	}
		//	if eOrderQuery == mgo.ErrNotFound {
		//		logs.Error(eOrderQuery)
		//		s := fmt.Sprintf("订单%v不存在", shopOrder.GetOrderId())
		//		logs.Error(s)
		//		return "订单不存在"
		//	} else {
		//		if order.GetItems().GetPointCardInfos() == nil || len(order.GetItems().GetPointCardInfos()) <= 0 {
		//			//通过物品的个数取得卡密个数
		//			pointCardList := GetPointCardByBuyInfos(order.GetSponsorAccount(),
		//				order.GetItems().GetPointCardName(),
		//				order.GetItems().GetCount())
		//			for nil != pointCardList && len(pointCardList) > 0 {
		//				for _, valuePoint := range pointCardList {
		//					itemPointCard := share_message.ShopPointCardInfo{
		//						CardId:       easygo.NewInt64(valuePoint.GetCardId()),
		//						CardNo:       easygo.NewString(valuePoint.GetCardNo()),
		//						CardPassword: easygo.NewString(valuePoint.GetCardPassword()),
		//						Key:          easygo.NewString(valuePoint.GetKey()),
		//					}
		//
		//					itemPointCards = append(itemPointCards, &itemPointCard)
		//					tempCardIds = append(tempCardIds, valuePoint.GetCardId())
		//				}
		//			}
		//		} else {
		//			for _, valuePoint := range order.GetItems().GetPointCardInfos() {
		//				itemPointCard := share_message.ShopPointCardInfo{
		//					CardId:       easygo.NewInt64(valuePoint.GetCardId()),
		//					CardNo:       easygo.NewString(valuePoint.GetCardNo()),
		//					CardPassword: easygo.NewString(valuePoint.GetCardPassword()),
		//					Key:          easygo.NewString(valuePoint.GetKey()),
		//				}
		//
		//				itemPointCards = append(itemPointCards, &itemPointCard)
		//				tempCardIds = append(tempCardIds, valuePoint.GetCardId())
		//			}
		//		}
		//
		//		//更新导入库中点卡状态和订单id
		//		_, errCard := colCard.UpdateAll(bson.M{"_id": bson.M{"$in": tempCardIds}}, bson.M{"$set": bson.M{"card_status": SHOP_POINT_CARD_SELLOUT, "order_no": shopOrder.GetOrderId()}})
		//
		//		//绑定点卡到订单
		//		eOrder := colOrder.Update(
		//			bson.M{"_id": shopOrder.GetOrderId()},
		//			bson.M{"$set": bson.M{"items.pointCardInfos": itemPointCards}})
		//
		//		if eOrder != nil || errCard != nil {
		//			logs.Error(eOrder)
		//			logs.Error(errCard)
		//			return "操作失败"
		//		}
		//	}
		//
		//	s = fmt.Sprintf("发起支付的时候 点卡订单%v 绑定点卡信息到订单结束", shopOrder.GetOrderId())
		//
		//	WriteFile("shop_order.log", s)
		//}

		//下架判断操作
		soldOutErr := ItemSoldOut(&shopOrder)
		if soldOutErr != "" {
			return soldOutErr
		}

		//说明是从购物车结算过来的支付操作
	} else {

		if bill.GetState() == SHOP_ORDER_WAIT_SEND {
			return "重复支付,请刷新"
		}

		if bill.GetState() != SHOP_ORDER_WAIT_PAY {
			return "操作失败,刷新重试！"
		}

		//取得各个子订单对应的商品的id
		var shopOrderList []*share_message.TableShopOrder
		//需要从订单表中取得数据
		eOrder := colOrder.Find(bson.M{"_id": bson.M{"$in": bill.GetOrderList()}}).All(&shopOrderList)
		if eOrder != nil && eOrder != mgo.ErrNotFound {
			logs.Error(eOrder)
			return "操作失败！"
		}

		if eOrder == mgo.ErrNotFound {
			logs.Error(eOrder)
			return "订单不存在"
		}
		//对于每一个商品都判断
		for _, value := range shopOrderList {

			if value.GetItems() != nil {

				if value.GetItems().GetItemType() == SHOP_POINT_CARD_CATEGORY {
					if value.GetState() == SHOP_ORDER_EVALUTE {
						return "存在订单重复支付,刷新重试！"
					}
				} else {
					if value.GetState() == SHOP_ORDER_WAIT_SEND {
						return "存在订单重复支付,刷新重试！"
					}
				}
			} else {
				return "存在无商品信息的订单！"
			}

			if value.GetState() != SHOP_ORDER_WAIT_PAY {

				return "操作失败,刷新重试！"
			}
			//判断库存和减少库存开始
			//1判断商品是否下架
			var item share_message.TableShopItem = share_message.TableShopItem{}

			eItem := colItem.Find(bson.M{"_id": value.GetItems().GetItemId(), "state": SHOP_ITEM_SALE}).One(&item)

			if eItem != nil && eItem != mgo.ErrNotFound {
				logs.Error(eItem)
				return "操作失败"
			}

			if eItem == mgo.ErrNotFound {
				s := fmt.Sprintf("订单%v对应的商品%v已经下架", value.GetOrderId(), value.GetItems().GetItemId())
				logs.Error(s)
				rts := fmt.Sprintf("商品%v已经下架", value.GetItems().GetName())
				return rts
			}

			//2库存判断
			if item.GetStockCount() < value.GetItems().GetCount() {
				s := fmt.Sprintf("该订单%v对应的商品%v库存不足", value.GetOrderId(), item.GetItemId())
				logs.Error(s)
				srt := fmt.Sprintf("商品%v库存不足", item.GetName())
				return srt
			}

			//如果是点卡 再判断一次实际导入库中的库存是否足够
			if item.GetType() != nil && item.GetType().GetType() == SHOP_POINT_CARD_CATEGORY {

				//通过物品的个数取得卡密个数
				pointCardList := GetPointCardByBuyInfos(item.GetPlayerAccount(), item.GetPointCardName(), value.GetItems().GetCount())
				if nil != pointCardList && (len(pointCardList) == 0 || int32(len(pointCardList)) < value.GetItems().GetCount()) {
					s := fmt.Sprintf("该订单%v对应的商品%v库存不足", value.GetOrderId(), item.GetItemId())
					logs.Error(s)
					return "商品库存不足"
				}
			}

			//3减少库存,减少商品表中的库存
			s := fmt.Sprintf("订单%v 当前物品%v原库存%v, 该扣除数量%v",
				value.GetOrderId(),
				item.GetItemId(),
				item.GetStockCount(),
				value.GetItems().GetCount())
			WriteFile("shop_order.log", s)
			eItem2 := colItem.Update(bson.M{"_id": item.GetItemId()}, bson.M{"$inc": bson.M{"stock_count": -value.GetItems().GetCount()}})

			if eItem2 != nil && eItem2 != mgo.ErrNotFound {
				logs.Error(eItem2)
				return "操作失败"
			}

			if eItem2 == mgo.ErrNotFound {
				s := fmt.Sprintf("订单%v对应的商品%v不存在", value.GetOrderId(), item.GetItemId())
				logs.Error(s)
				srt := fmt.Sprintf("商品%v不存在", item.GetName())
				return srt
			}

			//点卡的类型,支付完成后处理比较安全,支付开始只处理上面中商品的库存
			////如果是点卡,减少导入库中的库存,更新和绑定到订单
			//if item.GetType() != nil && item.GetType().GetType() == SHOP_POINT_CARD_CATEGORY {
			//
			//	s = fmt.Sprintf("发起支付的时候 点卡订单%v 绑定点卡信息到订单开始", value.GetOrderId())
			//
			//	order := share_message.TableShopOrder{}
			//	var itemPointCards []*share_message.ShopPointCardInfo = make([]*share_message.ShopPointCardInfo, 0)
			//	var tempCardIds []int64 = []int64{}
			//
			//	eOrderQuery := colOrder.Find(bson.M{"_id": value.GetOrderId()}).One(&order)
			//
			//	if eOrderQuery != nil && eOrderQuery != mgo.ErrNotFound {
			//		logs.Error(eOrderQuery)
			//		return "操作失败"
			//	}
			//	if eOrderQuery == mgo.ErrNotFound {
			//		logs.Error(eOrderQuery)
			//		s := fmt.Sprintf("订单%v不存在", value.GetOrderId())
			//		logs.Error(s)
			//		return "订单不存在"
			//	} else {
			//		if order.GetItems().GetPointCardInfos() == nil || len(order.GetItems().GetPointCardInfos()) <= 0 {
			//			//通过物品的个数取得卡密个数
			//			pointCardList := GetPointCardByBuyInfos(order.GetSponsorAccount(),
			//				order.GetItems().GetPointCardName(),
			//				order.GetItems().GetCount())
			//			for nil != pointCardList && len(pointCardList) > 0 {
			//				for _, valuePoint := range pointCardList {
			//					itemPointCard := share_message.ShopPointCardInfo{
			//						CardId:       easygo.NewInt64(valuePoint.GetCardId()),
			//						CardNo:       easygo.NewString(valuePoint.GetCardNo()),
			//						CardPassword: easygo.NewString(valuePoint.GetCardPassword()),
			//						Key:          easygo.NewString(valuePoint.GetKey()),
			//					}
			//
			//					itemPointCards = append(itemPointCards, &itemPointCard)
			//					tempCardIds = append(tempCardIds, valuePoint.GetCardId())
			//				}
			//			}
			//		} else {
			//			for _, valuePoint := range order.GetItems().GetPointCardInfos() {
			//				itemPointCard := share_message.ShopPointCardInfo{
			//					CardId:       easygo.NewInt64(valuePoint.GetCardId()),
			//					CardNo:       easygo.NewString(valuePoint.GetCardNo()),
			//					CardPassword: easygo.NewString(valuePoint.GetCardPassword()),
			//					Key:          easygo.NewString(valuePoint.GetKey()),
			//				}
			//
			//				itemPointCards = append(itemPointCards, &itemPointCard)
			//				tempCardIds = append(tempCardIds, valuePoint.GetCardId())
			//			}
			//		}
			//
			//		//更新导入库中点卡状态和订单id
			//		_, errCard := colCard.UpdateAll(bson.M{"_id": bson.M{"$in": tempCardIds}}, bson.M{"$set": bson.M{"card_status": SHOP_POINT_CARD_SELLOUT, "order_no": value.GetOrderId()}})
			//
			//		//绑定点卡到订单
			//		eOrder := colOrder.Update(
			//			bson.M{"_id": value.GetOrderId()},
			//			bson.M{"$set": bson.M{"items.pointCardInfos": itemPointCards}})
			//
			//		if eOrder != nil || errCard != nil {
			//			logs.Error(eOrder)
			//			logs.Error(errCard)
			//			return "操作失败"
			//		}
			//	}
			//
			//	s = fmt.Sprintf("发起支付的时候 点卡订单%v 绑定点卡信息到订单结束", value.GetOrderId())
			//
			//	WriteFile("shop_order.log", s)
			//}

			//下架判断操作
			soldOutErr := ItemSoldOut(value)
			if soldOutErr != "" {
				return soldOutErr
			}
		}

	}
	return ""
}

//商城购买的时候出错恢复库存
func ShopRecoverStock(orderId int64) string {

	//取得订单的订单id和订单对应的商品id
	var bill *share_message.TableBill = &share_message.TableBill{}

	colBill, closeFunBill := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_BILLS)
	defer closeFunBill()

	colOrder, closeFunOrder := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_ORDERS)
	defer closeFunOrder()

	colItem, closeFunItem := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_ITEMS)
	defer closeFunItem()

	eBill := colBill.Find(bson.M{"_id": orderId}).One(bill)
	if eBill != nil && eBill != mgo.ErrNotFound {
		logs.Error(eBill)
		return "操作失败"
	}

	//说明不是从购物车结算过来的支付操作
	if eBill == mgo.ErrNotFound {

		var shopOrder share_message.TableShopOrder = share_message.TableShopOrder{}
		//需要从订单表中取得数据

		eOrder := colOrder.Find(bson.M{"_id": orderId}).One(&shopOrder)
		if eOrder != nil && eOrder != mgo.ErrNotFound {
			logs.Error(eOrder)
			return "操作失败"
		}

		if eOrder == mgo.ErrNotFound {
			logs.Error(eOrder)
			return "订单不存在"
		}

		//恢复库存开始
		var item share_message.TableShopItem = share_message.TableShopItem{}

		eItem := colItem.Find(bson.M{"_id": shopOrder.GetItems().GetItemId()}).One(&item)

		if eItem != nil && eItem != mgo.ErrNotFound {
			logs.Error(eItem)
			return "操作失败"
		}

		if eItem == mgo.ErrNotFound {
			s := fmt.Sprintf("订单%v对应的商品%v不存在", shopOrder.GetOrderId(), shopOrder.GetItems().GetItemId())
			logs.Error(s)
			return "商品不存在"
		}
		//恢复库存
		eItem2 := colItem.Update(bson.M{"_id": item.GetItemId()}, bson.M{"$inc": bson.M{"stock_count": shopOrder.GetItems().GetCount()}})

		if eItem2 != nil && eItem2 != mgo.ErrNotFound {
			logs.Error(eItem2)
			return "操作失败"
		}

		if eItem2 == mgo.ErrNotFound {
			s := fmt.Sprintf("订单%v对应的商品%v不存在", shopOrder.GetOrderId(), shopOrder.GetItems().GetItemId())
			logs.Error(s)
			return "商品不存在"
		}

		//点卡的放到支付完成后做跟减少库存对应
		////如果是点卡,且已经将点卡信息绑定成功到订单的,恢复导入库中的库存
		//if shopOrder.GetItems() != nil &&
		//	shopOrder.GetItems().GetItemType() == SHOP_POINT_CARD_CATEGORY &&
		//	shopOrder.GetItems().GetPointCardInfos() != nil && len(shopOrder.GetItems().GetPointCardInfos()) > 0 {
		//	errPoint := PointCardRecoverStock(shopOrder)
		//	if errPoint != "" {
		//		logs.Error(errPoint)
		//		return "恢复点卡库存失败"
		//	}
		//}

		//上架判断操作
		onSaleErr := ItemForOnSale(&shopOrder)
		if onSaleErr != "" {
			return onSaleErr
		}
		//说明是从购物车结算过来的支付操作
	} else {

		//取得各个子订单对应的商品的id
		var shopOrderList []*share_message.TableShopOrder
		//需要从订单表中取得数据
		eOrder := colOrder.Find(bson.M{"_id": bson.M{"$in": bill.GetOrderList()}}).All(&shopOrderList)
		if eOrder != nil && eOrder != mgo.ErrNotFound {
			logs.Error(eOrder)
			return "操作失败！"
		}

		if eOrder == mgo.ErrNotFound {
			logs.Error(eOrder)
			return "订单不存在"
		}
		//对于每一个商品都判断
		for _, value := range shopOrderList {

			//恢复库存开始
			var item share_message.TableShopItem = share_message.TableShopItem{}

			eItem := colItem.Find(bson.M{"_id": value.GetItems().GetItemId()}).One(&item)

			if eItem != nil && eItem != mgo.ErrNotFound {
				logs.Error(eItem)
				return "操作失败"
			}

			if eItem == mgo.ErrNotFound {
				s := fmt.Sprintf("订单%v对应的商品%v不存在", value.GetOrderId(), value.GetItems().GetItemId())
				logs.Error(s)
				rts := fmt.Sprintf("商品%v不存在", value.GetItems().GetName())
				return rts
			}

			//恢复库存
			eItem2 := colItem.Update(bson.M{"_id": item.GetItemId()}, bson.M{"$inc": bson.M{"stock_count": value.GetItems().GetCount()}})

			if eItem2 != nil && eItem2 != mgo.ErrNotFound {
				logs.Error(eItem2)
				return "操作失败"
			}

			if eItem2 == mgo.ErrNotFound {
				s := fmt.Sprintf("订单%v对应的商品%v不存在", value.GetOrderId(), item.GetItemId())
				logs.Error(s)
				srt := fmt.Sprintf("商品%v不存在", item.GetName())
				return srt
			}

			//点卡的放到支付完成后做跟减少库存对应
			////如果是点卡,且已经将点卡信息绑定成功到订单的,恢复导入库中的库存
			//if value.GetItems() != nil &&
			//	value.GetItems().GetItemType() == SHOP_POINT_CARD_CATEGORY &&
			//	value.GetItems().GetPointCardInfos() != nil && len(value.GetItems().GetPointCardInfos()) > 0 {
			//
			//	errPoint := PointCardRecoverStock(*value)
			//	if errPoint != "" {
			//		logs.Error(errPoint)
			//		return "恢复点卡库存失败"
			//	}
			//}

			//上架判断操作
			onSaleErr := ItemForOnSale(value)
			if onSaleErr != "" {
				return onSaleErr
			}
		}

	}
	return ""
}

//id查询商品是否存在
func QueryShopItemById(Id int64) *share_message.TableShopItem {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_ITEMS)
	defer closeFun()
	rp := &share_message.TableShopItem{}

	errc := col.Find(bson.M{"_id": Id}).One(rp)
	if errc != nil && errc != mgo.ErrNotFound {
		panic(errc)
	}
	if errc == mgo.ErrNotFound {
		return nil
	}
	return rp
}

//通过物品的个数取得卡密个数(用于库存判断等)
func GetPointCardByBuyInfos(account, pointCardName string, count int32) []*share_message.TableShopPointCard {

	list := make([]*share_message.TableShopPointCard, 0)

	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_POINT_CARD)
	defer closeFun()

	e := col.Find(bson.M{"seller_account": account,
		"card_name":   pointCardName,
		"card_status": SHOP_POINT_CARD_SALE}).Sort("_id").Limit(int(count)).All(&list)

	if e != nil {
		logs.Error("通过物品的个数取得卡密个数(用于库存判断等) err:", e)
	}

	return list
}

//通过order中的信息恢复导入库中的库存等信息
//orderPara传入之前上层调用的地方需要判断各项值
func PointCardRecoverStock(orderPara share_message.TableShopOrder) string {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_POINT_CARD)
	defer closeFun()

	var tempCardIds []int64 = make([]int64, 0)

	for _, valuePoint := range orderPara.GetItems().GetPointCardInfos() {

		tempCardIds = append(tempCardIds, valuePoint.GetCardId())
	}

	if len(tempCardIds) > 0 {
		_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": tempCardIds}}, bson.M{"$set": bson.M{"card_status": SHOP_POINT_CARD_SALE, "order_no": 0}})

		if err != nil {
			logs.Error(err)
			return "操作失败"
		}
	}

	//将点卡绑定到订单的信息清空
	colOrder, closeFunOrder := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_ORDERS)
	defer closeFunOrder()
	var itemPointCards []*share_message.ShopPointCardInfo = make([]*share_message.ShopPointCardInfo, 0)
	//解绑点卡
	eOrder := colOrder.Update(
		bson.M{"_id": orderPara.GetOrderId()},
		bson.M{"$set": bson.M{"items.pointCardInfos": itemPointCards, "items.PointCardName": ""}})

	if eOrder != nil {
		logs.Error(eOrder)

		return "操作失败"
	}
	return ""
}

//商品下架
func ItemSoldOut(order *share_message.TableShopOrder) string {

	if order == nil || order.GetOrderId() == 0 || order.GetItems() == nil {
		return "订单不存在"
	}
	//下架
	var item share_message.TableShopItem = share_message.TableShopItem{}

	colItem, closeFunItem := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_ITEMS)
	defer closeFunItem()

	eItem := colItem.Find(bson.M{"_id": order.GetItems().GetItemId()}).One(&item)

	if eItem != nil && eItem != mgo.ErrNotFound {
		logs.Error(eItem)
		return "操作失败"
	}

	if eItem == mgo.ErrNotFound {
		s := fmt.Sprintf("订单%v对应的商品%v不存在", order.GetOrderId(), item.GetItemId())
		logs.Error(s)
		return s
	}

	//此时在发起支付的时候已经减少了库存这里只要判断数据库中的库存就可以下架
	if item.GetStockCount() <= 0 {

		//下架商品
		err := colItem.Update(bson.M{"_id": item.GetItemId()}, bson.M{"$set": bson.M{"state": SHOP_ITEM_SOLD_OUT, "sold_out_time": time.Now().Unix()}})

		if err != nil && err != mgo.ErrNotFound {
			logs.Error(err)
			return "操作失败"
		}

		if err == mgo.ErrNotFound {
			s := fmt.Sprintf("订单%v对应的商品%v不存在", order.GetOrderId(), item.GetItemId())
			logs.Error(s)
			return s
		}
	}

	return ""
}

//商品上架
func ItemForOnSale(order *share_message.TableShopOrder) string {

	if order == nil || order.GetOrderId() == 0 || order.GetItems() == nil {
		return "订单不存在"
	}
	//上架
	var item share_message.TableShopItem = share_message.TableShopItem{}

	colItem, closeFunItem := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_ITEMS)
	defer closeFunItem()

	eItem := colItem.Find(bson.M{"_id": order.GetItems().GetItemId()}).One(&item)

	if eItem != nil && eItem != mgo.ErrNotFound {
		logs.Error(eItem)
		return "操作失败"
	}

	if eItem == mgo.ErrNotFound {
		s := fmt.Sprintf("订单%v对应的商品%v不存在", order.GetOrderId(), item.GetItemId())
		logs.Error(s)
		return s
	}

	//判断是否需要上架
	if item.GetStockCount() > 0 && item.GetState() == SHOP_ITEM_SOLD_OUT {

		//上架商品
		err := colItem.Update(bson.M{"_id": item.GetItemId()}, bson.M{"$set": bson.M{"state": SHOP_ITEM_SALE}})

		if err != nil && err != mgo.ErrNotFound {
			logs.Error(err)
			return "操作失败"
		}

		if err == mgo.ErrNotFound {
			s := fmt.Sprintf("订单%v对应的商品%v不存在", order.GetOrderId(), item.GetItemId())
			logs.Error(s)
			return s
		}
	}

	return ""
}

//各个商城订单待发货的时候取消动作过商城订单id查询待支付和已超时的支付订单,确定恢复库存用的订单
func GetPayOrderListByShopOrderId(orderId int64) (int, string) {

	//1商城子订单取消时要恢复的库存数
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ORDER)
	defer closeFun()

	colBill, closeFunBill := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SHOP_BILLS)
	defer closeFunBill()

	query := col.Find(bson.M{"PayTargetId": orderId, "PayWay": PAY_TYPE_SHOP, "$or": []bson.M{
		bson.M{"PayStatus": PAY_ST_WAITTING},
		bson.M{"PayStatus": PAY_ST_CANCEL},
		bson.M{"PayStatus": PAY_ST_REFUSE}}})
	count, err := query.Count()

	if err != nil && err != mgo.ErrNotFound {

		logs.Error(err)

		return 0, "操作失败"
	}

	var bill *share_message.TableBill
	//2判断传入的商城子订单是否存在bill中订单的支付订单也需要当作库存恢复
	errBill := colBill.Find(bson.M{"order_list": orderId, "state": SHOP_ORDER_WAIT_PAY}).One(bill)
	if errBill != nil && errBill != mgo.ErrNotFound {

		logs.Error(errBill)

		return 0, "操作失败"
	}

	var count2 int
	if errBill != mgo.ErrNotFound && bill != nil {

		queryByBill := col.Find(bson.M{"PayTargetId": bill.GetOrderId(), "PayWay": PAY_TYPE_SHOP, "$or": []bson.M{
			bson.M{"PayStatus": PAY_ST_WAITTING}, bson.M{"PayStatus": PAY_ST_CANCEL},
			bson.M{"PayStatus": PAY_ST_REFUSE}}})
		countByBill, errByBill := queryByBill.Count()
		count2 = countByBill
		if errByBill != nil && errByBill != mgo.ErrNotFound {

			logs.Error(errByBill)

			return 0, "操作失败"
		}

	}

	return count + count2, ""
}
