package client_hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/share_message"
)

type _ = base.NoReturn

type IShopClient2Hall interface {
	RpcShopItemUpload(reqMsg *share_message.ShopItemUploadInfo) *share_message.ShopItemUploadResult
	RpcShopItemUpload_(reqMsg *share_message.ShopItemUploadInfo) (*share_message.ShopItemUploadResult, easygo.IRpcInterrupt)
	RpcShopItemEdit(reqMsg *share_message.ShopItemEditInfo) *share_message.ShopItemUploadResult
	RpcShopItemEdit_(reqMsg *share_message.ShopItemEditInfo) (*share_message.ShopItemUploadResult, easygo.IRpcInterrupt)
	RpcShopItemDelete(reqMsg *share_message.ShopItemID) *share_message.ShopItemDeleteResult
	RpcShopItemDelete_(reqMsg *share_message.ShopItemID) (*share_message.ShopItemDeleteResult, easygo.IRpcInterrupt)
	RpcShopItemSoldOut(reqMsg *share_message.ShopItemID) *share_message.ShopItemSoldOutResult
	RpcShopItemSoldOut_(reqMsg *share_message.ShopItemID) (*share_message.ShopItemSoldOutResult, easygo.IRpcInterrupt)
	RpcShopItemList(reqMsg *share_message.ShopInfo) *share_message.ItemList
	RpcShopItemList_(reqMsg *share_message.ShopInfo) (*share_message.ItemList, easygo.IRpcInterrupt)
	RpcShopItemInfo(reqMsg *share_message.ShopItemInfo) *share_message.ShopItemShowDetail
	RpcShopItemInfo_(reqMsg *share_message.ShopItemInfo) (*share_message.ShopItemShowDetail, easygo.IRpcInterrupt)
	RpcReceiveAddressList(reqMsg *base.Empty) *share_message.ReceiveAddressList
	RpcReceiveAddressList_(reqMsg *base.Empty) (*share_message.ReceiveAddressList, easygo.IRpcInterrupt)
	RpcReceiveAddressEdit(reqMsg *share_message.ReceiveAddressInfo) *share_message.ReceiveAddressEditResult
	RpcReceiveAddressEdit_(reqMsg *share_message.ReceiveAddressInfo) (*share_message.ReceiveAddressEditResult, easygo.IRpcInterrupt)
	RpcReceiveAddressAdd(reqMsg *share_message.ReceiveAddress) *share_message.ReceiveAddressAddResult
	RpcReceiveAddressAdd_(reqMsg *share_message.ReceiveAddress) (*share_message.ReceiveAddressAddResult, easygo.IRpcInterrupt)
	RpcReceiveAddressDelete(reqMsg *share_message.ReceiveAddressID) *share_message.ReceiveAddressRemoveResult
	RpcReceiveAddressDelete_(reqMsg *share_message.ReceiveAddressID) (*share_message.ReceiveAddressRemoveResult, easygo.IRpcInterrupt)
	RpcDeliverAddressList(reqMsg *base.Empty) *share_message.DeliverAddressList
	RpcDeliverAddressList_(reqMsg *base.Empty) (*share_message.DeliverAddressList, easygo.IRpcInterrupt)
	RpcDeliverAddressEdit(reqMsg *share_message.DeliverAddressInfo) *share_message.DeliverAddressEditResult
	RpcDeliverAddressEdit_(reqMsg *share_message.DeliverAddressInfo) (*share_message.DeliverAddressEditResult, easygo.IRpcInterrupt)
	RpcDeliverAddressAdd(reqMsg *share_message.DeliverAddress) *share_message.DeliverAddressAddResult
	RpcDeliverAddressAdd_(reqMsg *share_message.DeliverAddress) (*share_message.DeliverAddressAddResult, easygo.IRpcInterrupt)
	RpcDeliverAddressDelete(reqMsg *share_message.DeliverAddressID) *share_message.DeliverAddressRemoveResult
	RpcDeliverAddressDelete_(reqMsg *share_message.DeliverAddressID) (*share_message.DeliverAddressRemoveResult, easygo.IRpcInterrupt)
	RpcShopItemCommentUpload(reqMsg *share_message.UploadComment) *share_message.UploadCommentResult
	RpcShopItemCommentUpload_(reqMsg *share_message.UploadComment) (*share_message.UploadCommentResult, easygo.IRpcInterrupt)
	RpcShopItemCommentList(reqMsg *share_message.ShopCommentList) *share_message.ShopCommentListResult
	RpcShopItemCommentList_(reqMsg *share_message.ShopCommentList) (*share_message.ShopCommentListResult, easygo.IRpcInterrupt)
	RpcLikeComment(reqMsg *share_message.LikeCommentInfo) *share_message.LikeCommentResult
	RpcLikeComment_(reqMsg *share_message.LikeCommentInfo) (*share_message.LikeCommentResult, easygo.IRpcInterrupt)
	RpcCartInfo(reqMsg *base.Empty) *share_message.CartItemInfoList
	RpcCartInfo_(reqMsg *base.Empty) (*share_message.CartItemInfoList, easygo.IRpcInterrupt)
	RpcAddItemToCart(reqMsg *share_message.ShopItemID) *share_message.AddCartResult
	RpcAddItemToCart_(reqMsg *share_message.ShopItemID) (*share_message.AddCartResult, easygo.IRpcInterrupt)
	RpcSubItemToCart(reqMsg *share_message.ShopItemID) *share_message.SubCartResult
	RpcSubItemToCart_(reqMsg *share_message.ShopItemID) (*share_message.SubCartResult, easygo.IRpcInterrupt)
	RpcRemoveItemFromCart(reqMsg *share_message.ItemIdList) *share_message.RemoveCartResult
	RpcRemoveItemFromCart_(reqMsg *share_message.ItemIdList) (*share_message.RemoveCartResult, easygo.IRpcInterrupt)
	RpcStoreInfo(reqMsg *base.Empty) *share_message.StoreItemList
	RpcStoreInfo_(reqMsg *base.Empty) (*share_message.StoreItemList, easygo.IRpcInterrupt)
	RpcAddItemToStore(reqMsg *share_message.ShopItemID) *share_message.AddStoreResult
	RpcAddItemToStore_(reqMsg *share_message.ShopItemID) (*share_message.AddStoreResult, easygo.IRpcInterrupt)
	RpcRemoveItemFromStore(reqMsg *share_message.ShopItemID) *share_message.RemoveStoreResult
	RpcRemoveItemFromStore_(reqMsg *share_message.ShopItemID) (*share_message.RemoveStoreResult, easygo.IRpcInterrupt)
	RpcBatchAddItemToStore(reqMsg *share_message.ItemIdList) *share_message.BatchAddStoreResult
	RpcBatchAddItemToStore_(reqMsg *share_message.ItemIdList) (*share_message.BatchAddStoreResult, easygo.IRpcInterrupt)
	RpcCreateOrder(reqMsg *share_message.BuyItemInfo) *share_message.BuyItemResult
	RpcCreateOrder_(reqMsg *share_message.BuyItemInfo) (*share_message.BuyItemResult, easygo.IRpcInterrupt)
	RpcOrderList(reqMsg *share_message.OrderInfo) *share_message.OrderItemList
	RpcOrderList_(reqMsg *share_message.OrderInfo) (*share_message.OrderItemList, easygo.IRpcInterrupt)
	RpcOrderInfo(reqMsg *share_message.OrderDetailInfoPara) *share_message.OrderDetailInfoShow
	RpcOrderInfo_(reqMsg *share_message.OrderDetailInfoPara) (*share_message.OrderDetailInfoShow, easygo.IRpcInterrupt)
	RpcCancelOrder(reqMsg *share_message.OrderID) *share_message.CancelOrderResult
	RpcCancelOrder_(reqMsg *share_message.OrderID) (*share_message.CancelOrderResult, easygo.IRpcInterrupt)
	RpcDeleteOrder(reqMsg *share_message.OrderID) *share_message.DeleteOrderResult
	RpcDeleteOrder_(reqMsg *share_message.OrderID) (*share_message.DeleteOrderResult, easygo.IRpcInterrupt)
	RpcSettlementBtn(reqMsg *share_message.SettlementInfo) *share_message.SettlementResult
	RpcSettlementBtn_(reqMsg *share_message.SettlementInfo) (*share_message.SettlementResult, easygo.IRpcInterrupt)
	RpcItemListForMyReleaseOnline(reqMsg *share_message.MyReleaseInfo) *share_message.ItemListForMyRelease
	RpcItemListForMyReleaseOnline_(reqMsg *share_message.MyReleaseInfo) (*share_message.ItemListForMyRelease, easygo.IRpcInterrupt)
	RpcItemListForMyReleaseOffline(reqMsg *share_message.MyReleaseInfo) *share_message.ItemListForMyRelease
	RpcItemListForMyReleaseOffline_(reqMsg *share_message.MyReleaseInfo) (*share_message.ItemListForMyRelease, easygo.IRpcInterrupt)
	RpcItemSearch(reqMsg *share_message.SearchInfo) *share_message.SearchResult
	RpcItemSearch_(reqMsg *share_message.SearchInfo) (*share_message.SearchResult, easygo.IRpcInterrupt)
	RpcDelayReceiveItem(reqMsg *share_message.OrderID) *share_message.DelayReceiveResult
	RpcDelayReceiveItem_(reqMsg *share_message.OrderID) (*share_message.DelayReceiveResult, easygo.IRpcInterrupt)
	RpcEditOrderAddress(reqMsg *share_message.EditOrderAddress) *share_message.EditOrderAddressResult
	RpcEditOrderAddress_(reqMsg *share_message.EditOrderAddress) (*share_message.EditOrderAddressResult, easygo.IRpcInterrupt)
	RpcEditDeliverAddress(reqMsg *share_message.EditDeliverAddress) *share_message.EditDeliverAddressResult
	RpcEditDeliverAddress_(reqMsg *share_message.EditDeliverAddress) (*share_message.EditDeliverAddressResult, easygo.IRpcInterrupt)
	RpcConfirmReceive(reqMsg *share_message.OrderID) *share_message.ConfirmReceiveResult
	RpcConfirmReceive_(reqMsg *share_message.OrderID) (*share_message.ConfirmReceiveResult, easygo.IRpcInterrupt)
	RpcShopItemEvaluteUpload(reqMsg *share_message.UploadEvalute) *share_message.UploadCommentResult
	RpcShopItemEvaluteUpload_(reqMsg *share_message.UploadEvalute) (*share_message.UploadCommentResult, easygo.IRpcInterrupt)
	RpcExpressCodeUpload(reqMsg *share_message.ExpressInfo) *share_message.ExpressCodeResult
	RpcExpressCodeUpload_(reqMsg *share_message.ExpressInfo) (*share_message.ExpressCodeResult, easygo.IRpcInterrupt)
	RpcExpressComInfos(reqMsg *base.Empty) *share_message.ExpressComInfosResult
	RpcExpressComInfos_(reqMsg *base.Empty) (*share_message.ExpressComInfosResult, easygo.IRpcInterrupt)
	RpcExpressInfos(reqMsg *share_message.QueryExpressInfo) *share_message.QueryExpressInfosResult
	RpcExpressInfos_(reqMsg *share_message.QueryExpressInfo) (*share_message.QueryExpressInfosResult, easygo.IRpcInterrupt)
	RpcNotifySendItem(reqMsg *share_message.OrderID) *share_message.NotifySendItemResult
	RpcNotifySendItem_(reqMsg *share_message.OrderID) (*share_message.NotifySendItemResult, easygo.IRpcInterrupt)
	RpcShopItemMessage(reqMsg *share_message.ShopItemMessageListInfo) *share_message.ShopItemMessageList
	RpcShopItemMessage_(reqMsg *share_message.ShopItemMessageListInfo) (*share_message.ShopItemMessageList, easygo.IRpcInterrupt)
	RpcShopMessageFlgUpd(reqMsg *share_message.MessageIdList)
	RpcGetShopOrderNotifyInfos(reqMsg *share_message.PlayerID) *share_message.ShopOrderIdList
	RpcGetShopOrderNotifyInfos_(reqMsg *share_message.PlayerID) (*share_message.ShopOrderIdList, easygo.IRpcInterrupt)
	RpcShopOrderNotifyFlgUpd(reqMsg *share_message.ShopOrderNotifyFlgUpdInfo)
	RpcShopUploadAuth(reqMsg *share_message.PlayerID) *share_message.UploadAuthResult
	RpcShopUploadAuth_(reqMsg *share_message.PlayerID) (*share_message.UploadAuthResult, easygo.IRpcInterrupt)
	RpcShopUploadAuthConfirm(reqMsg *share_message.PlayerID) *base.Empty
	RpcShopUploadAuthConfirm_(reqMsg *share_message.PlayerID) (*base.Empty, easygo.IRpcInterrupt)
}

type ShopClient2Hall struct {
	Sender easygo.IMessageSender
}

func (self *ShopClient2Hall) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *ShopClient2Hall) RpcShopItemUpload(reqMsg *share_message.ShopItemUploadInfo) *share_message.ShopItemUploadResult {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemUpload", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ShopItemUploadResult)
}

func (self *ShopClient2Hall) RpcShopItemUpload_(reqMsg *share_message.ShopItemUploadInfo) (*share_message.ShopItemUploadResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemUpload", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ShopItemUploadResult), e
}
func (self *ShopClient2Hall) RpcShopItemEdit(reqMsg *share_message.ShopItemEditInfo) *share_message.ShopItemUploadResult {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemEdit", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ShopItemUploadResult)
}

func (self *ShopClient2Hall) RpcShopItemEdit_(reqMsg *share_message.ShopItemEditInfo) (*share_message.ShopItemUploadResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemEdit", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ShopItemUploadResult), e
}
func (self *ShopClient2Hall) RpcShopItemDelete(reqMsg *share_message.ShopItemID) *share_message.ShopItemDeleteResult {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemDelete", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ShopItemDeleteResult)
}

func (self *ShopClient2Hall) RpcShopItemDelete_(reqMsg *share_message.ShopItemID) (*share_message.ShopItemDeleteResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemDelete", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ShopItemDeleteResult), e
}
func (self *ShopClient2Hall) RpcShopItemSoldOut(reqMsg *share_message.ShopItemID) *share_message.ShopItemSoldOutResult {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemSoldOut", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ShopItemSoldOutResult)
}

func (self *ShopClient2Hall) RpcShopItemSoldOut_(reqMsg *share_message.ShopItemID) (*share_message.ShopItemSoldOutResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemSoldOut", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ShopItemSoldOutResult), e
}
func (self *ShopClient2Hall) RpcShopItemList(reqMsg *share_message.ShopInfo) *share_message.ItemList {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ItemList)
}

func (self *ShopClient2Hall) RpcShopItemList_(reqMsg *share_message.ShopInfo) (*share_message.ItemList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ItemList), e
}
func (self *ShopClient2Hall) RpcShopItemInfo(reqMsg *share_message.ShopItemInfo) *share_message.ShopItemShowDetail {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ShopItemShowDetail)
}

func (self *ShopClient2Hall) RpcShopItemInfo_(reqMsg *share_message.ShopItemInfo) (*share_message.ShopItemShowDetail, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ShopItemShowDetail), e
}
func (self *ShopClient2Hall) RpcReceiveAddressList(reqMsg *base.Empty) *share_message.ReceiveAddressList {
	msg, e := self.Sender.CallRpcMethod("RpcReceiveAddressList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ReceiveAddressList)
}

func (self *ShopClient2Hall) RpcReceiveAddressList_(reqMsg *base.Empty) (*share_message.ReceiveAddressList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReceiveAddressList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ReceiveAddressList), e
}
func (self *ShopClient2Hall) RpcReceiveAddressEdit(reqMsg *share_message.ReceiveAddressInfo) *share_message.ReceiveAddressEditResult {
	msg, e := self.Sender.CallRpcMethod("RpcReceiveAddressEdit", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ReceiveAddressEditResult)
}

func (self *ShopClient2Hall) RpcReceiveAddressEdit_(reqMsg *share_message.ReceiveAddressInfo) (*share_message.ReceiveAddressEditResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReceiveAddressEdit", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ReceiveAddressEditResult), e
}
func (self *ShopClient2Hall) RpcReceiveAddressAdd(reqMsg *share_message.ReceiveAddress) *share_message.ReceiveAddressAddResult {
	msg, e := self.Sender.CallRpcMethod("RpcReceiveAddressAdd", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ReceiveAddressAddResult)
}

func (self *ShopClient2Hall) RpcReceiveAddressAdd_(reqMsg *share_message.ReceiveAddress) (*share_message.ReceiveAddressAddResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReceiveAddressAdd", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ReceiveAddressAddResult), e
}
func (self *ShopClient2Hall) RpcReceiveAddressDelete(reqMsg *share_message.ReceiveAddressID) *share_message.ReceiveAddressRemoveResult {
	msg, e := self.Sender.CallRpcMethod("RpcReceiveAddressDelete", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ReceiveAddressRemoveResult)
}

func (self *ShopClient2Hall) RpcReceiveAddressDelete_(reqMsg *share_message.ReceiveAddressID) (*share_message.ReceiveAddressRemoveResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReceiveAddressDelete", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ReceiveAddressRemoveResult), e
}
func (self *ShopClient2Hall) RpcDeliverAddressList(reqMsg *base.Empty) *share_message.DeliverAddressList {
	msg, e := self.Sender.CallRpcMethod("RpcDeliverAddressList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.DeliverAddressList)
}

func (self *ShopClient2Hall) RpcDeliverAddressList_(reqMsg *base.Empty) (*share_message.DeliverAddressList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeliverAddressList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.DeliverAddressList), e
}
func (self *ShopClient2Hall) RpcDeliverAddressEdit(reqMsg *share_message.DeliverAddressInfo) *share_message.DeliverAddressEditResult {
	msg, e := self.Sender.CallRpcMethod("RpcDeliverAddressEdit", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.DeliverAddressEditResult)
}

func (self *ShopClient2Hall) RpcDeliverAddressEdit_(reqMsg *share_message.DeliverAddressInfo) (*share_message.DeliverAddressEditResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeliverAddressEdit", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.DeliverAddressEditResult), e
}
func (self *ShopClient2Hall) RpcDeliverAddressAdd(reqMsg *share_message.DeliverAddress) *share_message.DeliverAddressAddResult {
	msg, e := self.Sender.CallRpcMethod("RpcDeliverAddressAdd", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.DeliverAddressAddResult)
}

func (self *ShopClient2Hall) RpcDeliverAddressAdd_(reqMsg *share_message.DeliverAddress) (*share_message.DeliverAddressAddResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeliverAddressAdd", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.DeliverAddressAddResult), e
}
func (self *ShopClient2Hall) RpcDeliverAddressDelete(reqMsg *share_message.DeliverAddressID) *share_message.DeliverAddressRemoveResult {
	msg, e := self.Sender.CallRpcMethod("RpcDeliverAddressDelete", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.DeliverAddressRemoveResult)
}

func (self *ShopClient2Hall) RpcDeliverAddressDelete_(reqMsg *share_message.DeliverAddressID) (*share_message.DeliverAddressRemoveResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeliverAddressDelete", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.DeliverAddressRemoveResult), e
}
func (self *ShopClient2Hall) RpcShopItemCommentUpload(reqMsg *share_message.UploadComment) *share_message.UploadCommentResult {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemCommentUpload", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.UploadCommentResult)
}

func (self *ShopClient2Hall) RpcShopItemCommentUpload_(reqMsg *share_message.UploadComment) (*share_message.UploadCommentResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemCommentUpload", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.UploadCommentResult), e
}
func (self *ShopClient2Hall) RpcShopItemCommentList(reqMsg *share_message.ShopCommentList) *share_message.ShopCommentListResult {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemCommentList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ShopCommentListResult)
}

func (self *ShopClient2Hall) RpcShopItemCommentList_(reqMsg *share_message.ShopCommentList) (*share_message.ShopCommentListResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemCommentList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ShopCommentListResult), e
}
func (self *ShopClient2Hall) RpcLikeComment(reqMsg *share_message.LikeCommentInfo) *share_message.LikeCommentResult {
	msg, e := self.Sender.CallRpcMethod("RpcLikeComment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.LikeCommentResult)
}

func (self *ShopClient2Hall) RpcLikeComment_(reqMsg *share_message.LikeCommentInfo) (*share_message.LikeCommentResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcLikeComment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.LikeCommentResult), e
}
func (self *ShopClient2Hall) RpcCartInfo(reqMsg *base.Empty) *share_message.CartItemInfoList {
	msg, e := self.Sender.CallRpcMethod("RpcCartInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.CartItemInfoList)
}

func (self *ShopClient2Hall) RpcCartInfo_(reqMsg *base.Empty) (*share_message.CartItemInfoList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCartInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.CartItemInfoList), e
}
func (self *ShopClient2Hall) RpcAddItemToCart(reqMsg *share_message.ShopItemID) *share_message.AddCartResult {
	msg, e := self.Sender.CallRpcMethod("RpcAddItemToCart", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.AddCartResult)
}

func (self *ShopClient2Hall) RpcAddItemToCart_(reqMsg *share_message.ShopItemID) (*share_message.AddCartResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddItemToCart", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.AddCartResult), e
}
func (self *ShopClient2Hall) RpcSubItemToCart(reqMsg *share_message.ShopItemID) *share_message.SubCartResult {
	msg, e := self.Sender.CallRpcMethod("RpcSubItemToCart", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.SubCartResult)
}

func (self *ShopClient2Hall) RpcSubItemToCart_(reqMsg *share_message.ShopItemID) (*share_message.SubCartResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSubItemToCart", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.SubCartResult), e
}
func (self *ShopClient2Hall) RpcRemoveItemFromCart(reqMsg *share_message.ItemIdList) *share_message.RemoveCartResult {
	msg, e := self.Sender.CallRpcMethod("RpcRemoveItemFromCart", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.RemoveCartResult)
}

func (self *ShopClient2Hall) RpcRemoveItemFromCart_(reqMsg *share_message.ItemIdList) (*share_message.RemoveCartResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRemoveItemFromCart", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.RemoveCartResult), e
}
func (self *ShopClient2Hall) RpcStoreInfo(reqMsg *base.Empty) *share_message.StoreItemList {
	msg, e := self.Sender.CallRpcMethod("RpcStoreInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.StoreItemList)
}

func (self *ShopClient2Hall) RpcStoreInfo_(reqMsg *base.Empty) (*share_message.StoreItemList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcStoreInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.StoreItemList), e
}
func (self *ShopClient2Hall) RpcAddItemToStore(reqMsg *share_message.ShopItemID) *share_message.AddStoreResult {
	msg, e := self.Sender.CallRpcMethod("RpcAddItemToStore", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.AddStoreResult)
}

func (self *ShopClient2Hall) RpcAddItemToStore_(reqMsg *share_message.ShopItemID) (*share_message.AddStoreResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddItemToStore", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.AddStoreResult), e
}
func (self *ShopClient2Hall) RpcRemoveItemFromStore(reqMsg *share_message.ShopItemID) *share_message.RemoveStoreResult {
	msg, e := self.Sender.CallRpcMethod("RpcRemoveItemFromStore", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.RemoveStoreResult)
}

func (self *ShopClient2Hall) RpcRemoveItemFromStore_(reqMsg *share_message.ShopItemID) (*share_message.RemoveStoreResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRemoveItemFromStore", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.RemoveStoreResult), e
}
func (self *ShopClient2Hall) RpcBatchAddItemToStore(reqMsg *share_message.ItemIdList) *share_message.BatchAddStoreResult {
	msg, e := self.Sender.CallRpcMethod("RpcBatchAddItemToStore", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.BatchAddStoreResult)
}

func (self *ShopClient2Hall) RpcBatchAddItemToStore_(reqMsg *share_message.ItemIdList) (*share_message.BatchAddStoreResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBatchAddItemToStore", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.BatchAddStoreResult), e
}
func (self *ShopClient2Hall) RpcCreateOrder(reqMsg *share_message.BuyItemInfo) *share_message.BuyItemResult {
	msg, e := self.Sender.CallRpcMethod("RpcCreateOrder", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.BuyItemResult)
}

func (self *ShopClient2Hall) RpcCreateOrder_(reqMsg *share_message.BuyItemInfo) (*share_message.BuyItemResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCreateOrder", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.BuyItemResult), e
}
func (self *ShopClient2Hall) RpcOrderList(reqMsg *share_message.OrderInfo) *share_message.OrderItemList {
	msg, e := self.Sender.CallRpcMethod("RpcOrderList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.OrderItemList)
}

func (self *ShopClient2Hall) RpcOrderList_(reqMsg *share_message.OrderInfo) (*share_message.OrderItemList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOrderList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.OrderItemList), e
}
func (self *ShopClient2Hall) RpcOrderInfo(reqMsg *share_message.OrderDetailInfoPara) *share_message.OrderDetailInfoShow {
	msg, e := self.Sender.CallRpcMethod("RpcOrderInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.OrderDetailInfoShow)
}

func (self *ShopClient2Hall) RpcOrderInfo_(reqMsg *share_message.OrderDetailInfoPara) (*share_message.OrderDetailInfoShow, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcOrderInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.OrderDetailInfoShow), e
}
func (self *ShopClient2Hall) RpcCancelOrder(reqMsg *share_message.OrderID) *share_message.CancelOrderResult {
	msg, e := self.Sender.CallRpcMethod("RpcCancelOrder", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.CancelOrderResult)
}

func (self *ShopClient2Hall) RpcCancelOrder_(reqMsg *share_message.OrderID) (*share_message.CancelOrderResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCancelOrder", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.CancelOrderResult), e
}
func (self *ShopClient2Hall) RpcDeleteOrder(reqMsg *share_message.OrderID) *share_message.DeleteOrderResult {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteOrder", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.DeleteOrderResult)
}

func (self *ShopClient2Hall) RpcDeleteOrder_(reqMsg *share_message.OrderID) (*share_message.DeleteOrderResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDeleteOrder", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.DeleteOrderResult), e
}
func (self *ShopClient2Hall) RpcSettlementBtn(reqMsg *share_message.SettlementInfo) *share_message.SettlementResult {
	msg, e := self.Sender.CallRpcMethod("RpcSettlementBtn", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.SettlementResult)
}

func (self *ShopClient2Hall) RpcSettlementBtn_(reqMsg *share_message.SettlementInfo) (*share_message.SettlementResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSettlementBtn", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.SettlementResult), e
}
func (self *ShopClient2Hall) RpcItemListForMyReleaseOnline(reqMsg *share_message.MyReleaseInfo) *share_message.ItemListForMyRelease {
	msg, e := self.Sender.CallRpcMethod("RpcItemListForMyReleaseOnline", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ItemListForMyRelease)
}

func (self *ShopClient2Hall) RpcItemListForMyReleaseOnline_(reqMsg *share_message.MyReleaseInfo) (*share_message.ItemListForMyRelease, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcItemListForMyReleaseOnline", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ItemListForMyRelease), e
}
func (self *ShopClient2Hall) RpcItemListForMyReleaseOffline(reqMsg *share_message.MyReleaseInfo) *share_message.ItemListForMyRelease {
	msg, e := self.Sender.CallRpcMethod("RpcItemListForMyReleaseOffline", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ItemListForMyRelease)
}

func (self *ShopClient2Hall) RpcItemListForMyReleaseOffline_(reqMsg *share_message.MyReleaseInfo) (*share_message.ItemListForMyRelease, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcItemListForMyReleaseOffline", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ItemListForMyRelease), e
}
func (self *ShopClient2Hall) RpcItemSearch(reqMsg *share_message.SearchInfo) *share_message.SearchResult {
	msg, e := self.Sender.CallRpcMethod("RpcItemSearch", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.SearchResult)
}

func (self *ShopClient2Hall) RpcItemSearch_(reqMsg *share_message.SearchInfo) (*share_message.SearchResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcItemSearch", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.SearchResult), e
}
func (self *ShopClient2Hall) RpcDelayReceiveItem(reqMsg *share_message.OrderID) *share_message.DelayReceiveResult {
	msg, e := self.Sender.CallRpcMethod("RpcDelayReceiveItem", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.DelayReceiveResult)
}

func (self *ShopClient2Hall) RpcDelayReceiveItem_(reqMsg *share_message.OrderID) (*share_message.DelayReceiveResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelayReceiveItem", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.DelayReceiveResult), e
}
func (self *ShopClient2Hall) RpcEditOrderAddress(reqMsg *share_message.EditOrderAddress) *share_message.EditOrderAddressResult {
	msg, e := self.Sender.CallRpcMethod("RpcEditOrderAddress", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.EditOrderAddressResult)
}

func (self *ShopClient2Hall) RpcEditOrderAddress_(reqMsg *share_message.EditOrderAddress) (*share_message.EditOrderAddressResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditOrderAddress", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.EditOrderAddressResult), e
}
func (self *ShopClient2Hall) RpcEditDeliverAddress(reqMsg *share_message.EditDeliverAddress) *share_message.EditDeliverAddressResult {
	msg, e := self.Sender.CallRpcMethod("RpcEditDeliverAddress", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.EditDeliverAddressResult)
}

func (self *ShopClient2Hall) RpcEditDeliverAddress_(reqMsg *share_message.EditDeliverAddress) (*share_message.EditDeliverAddressResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditDeliverAddress", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.EditDeliverAddressResult), e
}
func (self *ShopClient2Hall) RpcConfirmReceive(reqMsg *share_message.OrderID) *share_message.ConfirmReceiveResult {
	msg, e := self.Sender.CallRpcMethod("RpcConfirmReceive", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ConfirmReceiveResult)
}

func (self *ShopClient2Hall) RpcConfirmReceive_(reqMsg *share_message.OrderID) (*share_message.ConfirmReceiveResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcConfirmReceive", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ConfirmReceiveResult), e
}
func (self *ShopClient2Hall) RpcShopItemEvaluteUpload(reqMsg *share_message.UploadEvalute) *share_message.UploadCommentResult {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemEvaluteUpload", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.UploadCommentResult)
}

func (self *ShopClient2Hall) RpcShopItemEvaluteUpload_(reqMsg *share_message.UploadEvalute) (*share_message.UploadCommentResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemEvaluteUpload", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.UploadCommentResult), e
}
func (self *ShopClient2Hall) RpcExpressCodeUpload(reqMsg *share_message.ExpressInfo) *share_message.ExpressCodeResult {
	msg, e := self.Sender.CallRpcMethod("RpcExpressCodeUpload", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ExpressCodeResult)
}

func (self *ShopClient2Hall) RpcExpressCodeUpload_(reqMsg *share_message.ExpressInfo) (*share_message.ExpressCodeResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcExpressCodeUpload", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ExpressCodeResult), e
}
func (self *ShopClient2Hall) RpcExpressComInfos(reqMsg *base.Empty) *share_message.ExpressComInfosResult {
	msg, e := self.Sender.CallRpcMethod("RpcExpressComInfos", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ExpressComInfosResult)
}

func (self *ShopClient2Hall) RpcExpressComInfos_(reqMsg *base.Empty) (*share_message.ExpressComInfosResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcExpressComInfos", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ExpressComInfosResult), e
}
func (self *ShopClient2Hall) RpcExpressInfos(reqMsg *share_message.QueryExpressInfo) *share_message.QueryExpressInfosResult {
	msg, e := self.Sender.CallRpcMethod("RpcExpressInfos", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.QueryExpressInfosResult)
}

func (self *ShopClient2Hall) RpcExpressInfos_(reqMsg *share_message.QueryExpressInfo) (*share_message.QueryExpressInfosResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcExpressInfos", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.QueryExpressInfosResult), e
}
func (self *ShopClient2Hall) RpcNotifySendItem(reqMsg *share_message.OrderID) *share_message.NotifySendItemResult {
	msg, e := self.Sender.CallRpcMethod("RpcNotifySendItem", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.NotifySendItemResult)
}

func (self *ShopClient2Hall) RpcNotifySendItem_(reqMsg *share_message.OrderID) (*share_message.NotifySendItemResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcNotifySendItem", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.NotifySendItemResult), e
}
func (self *ShopClient2Hall) RpcShopItemMessage(reqMsg *share_message.ShopItemMessageListInfo) *share_message.ShopItemMessageList {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ShopItemMessageList)
}

func (self *ShopClient2Hall) RpcShopItemMessage_(reqMsg *share_message.ShopItemMessageListInfo) (*share_message.ShopItemMessageList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShopItemMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ShopItemMessageList), e
}
func (self *ShopClient2Hall) RpcShopMessageFlgUpd(reqMsg *share_message.MessageIdList) {
	self.Sender.CallRpcMethod("RpcShopMessageFlgUpd", reqMsg)
}
func (self *ShopClient2Hall) RpcGetShopOrderNotifyInfos(reqMsg *share_message.PlayerID) *share_message.ShopOrderIdList {
	msg, e := self.Sender.CallRpcMethod("RpcGetShopOrderNotifyInfos", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.ShopOrderIdList)
}

func (self *ShopClient2Hall) RpcGetShopOrderNotifyInfos_(reqMsg *share_message.PlayerID) (*share_message.ShopOrderIdList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetShopOrderNotifyInfos", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.ShopOrderIdList), e
}
func (self *ShopClient2Hall) RpcShopOrderNotifyFlgUpd(reqMsg *share_message.ShopOrderNotifyFlgUpdInfo) {
	self.Sender.CallRpcMethod("RpcShopOrderNotifyFlgUpd", reqMsg)
}
func (self *ShopClient2Hall) RpcShopUploadAuth(reqMsg *share_message.PlayerID) *share_message.UploadAuthResult {
	msg, e := self.Sender.CallRpcMethod("RpcShopUploadAuth", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*share_message.UploadAuthResult)
}

func (self *ShopClient2Hall) RpcShopUploadAuth_(reqMsg *share_message.PlayerID) (*share_message.UploadAuthResult, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShopUploadAuth", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*share_message.UploadAuthResult), e
}
func (self *ShopClient2Hall) RpcShopUploadAuthConfirm(reqMsg *share_message.PlayerID) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcShopUploadAuthConfirm", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *ShopClient2Hall) RpcShopUploadAuthConfirm_(reqMsg *share_message.PlayerID) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcShopUploadAuthConfirm", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}

// ==========================================================
type IHall2ShopClient interface {
}

type Hall2ShopClient struct {
	Sender easygo.IMessageSender
}

func (self *Hall2ShopClient) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------
