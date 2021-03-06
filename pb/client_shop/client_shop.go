package client_shop

import (
	"game_server/easygo"
)

var UpRpc = map[string]easygo.Pair{
	"RpcLogin":                       easygo.Pair{"client_shop.LoginMsg", "base.Empty"},
	"RpcLogOut":                      easygo.Pair{"base.Empty", "base.Empty"},
	"RpcShopItemUpload":              easygo.Pair{"share_message.ShopItemUploadInfo", "share_message.ShopItemUploadResult"},
	"RpcShopItemEdit":                easygo.Pair{"share_message.ShopItemEditInfo", "share_message.ShopItemUploadResult"},
	"RpcShopItemDelete":              easygo.Pair{"share_message.ShopItemID", "share_message.ShopItemDeleteResult"},
	"RpcShopItemSoldOut":             easygo.Pair{"share_message.ShopItemID", "share_message.ShopItemSoldOutResult"},
	"RpcShopItemList":                easygo.Pair{"share_message.ShopInfo", "share_message.ItemList"},
	"RpcShopItemInfo":                easygo.Pair{"share_message.ShopItemInfo", "share_message.ShopItemShowDetail"},
	"RpcReceiveAddressList":          easygo.Pair{"base.Empty", "share_message.ReceiveAddressList"},
	"RpcReceiveAddressEdit":          easygo.Pair{"share_message.ReceiveAddressInfo", "share_message.ReceiveAddressEditResult"},
	"RpcReceiveAddressAdd":           easygo.Pair{"share_message.ReceiveAddress", "share_message.ReceiveAddressAddResult"},
	"RpcReceiveAddressDelete":        easygo.Pair{"share_message.ReceiveAddressID", "share_message.ReceiveAddressRemoveResult"},
	"RpcDeliverAddressList":          easygo.Pair{"base.Empty", "share_message.DeliverAddressList"},
	"RpcDeliverAddressEdit":          easygo.Pair{"share_message.DeliverAddressInfo", "share_message.DeliverAddressEditResult"},
	"RpcDeliverAddressAdd":           easygo.Pair{"share_message.DeliverAddress", "share_message.DeliverAddressAddResult"},
	"RpcDeliverAddressDelete":        easygo.Pair{"share_message.DeliverAddressID", "share_message.DeliverAddressRemoveResult"},
	"RpcShopItemCommentUpload":       easygo.Pair{"share_message.UploadComment", "share_message.UploadCommentResult"},
	"RpcShopItemCommentList":         easygo.Pair{"share_message.ShopCommentList", "share_message.ShopCommentListResult"},
	"RpcLikeComment":                 easygo.Pair{"share_message.LikeCommentInfo", "share_message.LikeCommentResult"},
	"RpcCartInfo":                    easygo.Pair{"base.Empty", "share_message.CartItemInfoList"},
	"RpcAddItemToCart":               easygo.Pair{"share_message.ShopItemID", "share_message.AddCartResult"},
	"RpcSubItemToCart":               easygo.Pair{"share_message.ShopItemID", "share_message.SubCartResult"},
	"RpcRemoveItemFromCart":          easygo.Pair{"share_message.ItemIdList", "share_message.RemoveCartResult"},
	"RpcStoreInfo":                   easygo.Pair{"base.Empty", "share_message.StoreItemList"},
	"RpcAddItemToStore":              easygo.Pair{"share_message.ShopItemID", "share_message.AddStoreResult"},
	"RpcRemoveItemFromStore":         easygo.Pair{"share_message.ShopItemID", "share_message.RemoveStoreResult"},
	"RpcBatchAddItemToStore":         easygo.Pair{"share_message.ItemIdList", "share_message.BatchAddStoreResult"},
	"RpcCreateOrder":                 easygo.Pair{"share_message.BuyItemInfo", "share_message.BuyItemResult"},
	"RpcOrderList":                   easygo.Pair{"share_message.OrderInfo", "share_message.OrderItemList"},
	"RpcOrderInfo":                   easygo.Pair{"share_message.OrderDetailInfoPara", "share_message.OrderDetailInfoShow"},
	"RpcCancelOrder":                 easygo.Pair{"share_message.OrderID", "share_message.CancelOrderResult"},
	"RpcDeleteOrder":                 easygo.Pair{"share_message.OrderID", "share_message.DeleteOrderResult"},
	"RpcSettlementBtn":               easygo.Pair{"share_message.SettlementInfo", "share_message.SettlementResult"},
	"RpcItemListForMyReleaseOnline":  easygo.Pair{"share_message.MyReleaseInfo", "share_message.ItemListForMyRelease"},
	"RpcItemListForMyReleaseOffline": easygo.Pair{"share_message.MyReleaseInfo", "share_message.ItemListForMyRelease"},
	"RpcItemSearch":                  easygo.Pair{"share_message.SearchInfo", "share_message.SearchResult"},
	"RpcDelayReceiveItem":            easygo.Pair{"share_message.OrderID", "share_message.DelayReceiveResult"},
	"RpcEditOrderAddress":            easygo.Pair{"share_message.EditOrderAddress", "share_message.EditOrderAddressResult"},
	"RpcEditDeliverAddress":          easygo.Pair{"share_message.EditDeliverAddress", "share_message.EditDeliverAddressResult"},
	"RpcConfirmReceive":              easygo.Pair{"share_message.OrderID", "share_message.ConfirmReceiveResult"},
	"RpcShopItemEvaluteUpload":       easygo.Pair{"share_message.UploadEvalute", "share_message.UploadCommentResult"},
	"RpcExpressCodeUpload":           easygo.Pair{"share_message.ExpressInfo", "share_message.ExpressCodeResult"},
	"RpcExpressComInfos":             easygo.Pair{"base.Empty", "share_message.ExpressComInfosResult"},
	"RpcExpressInfos":                easygo.Pair{"share_message.QueryExpressInfo", "share_message.QueryExpressInfosResult"},
	"RpcNotifySendItem":              easygo.Pair{"share_message.OrderID", "share_message.NotifySendItemResult"},
	"RpcShopItemMessage":             easygo.Pair{"share_message.ShopItemMessageListInfo", "share_message.ShopItemMessageList"},
	"RpcShopMessageFlgUpd":           easygo.Pair{"share_message.MessageIdList", "base.NoReturn"},
	"RpcGetShopOrderNotifyInfos":     easygo.Pair{"share_message.PlayerID", "share_message.ShopOrderIdList"},
	"RpcShopOrderNotifyFlgUpd":       easygo.Pair{"share_message.ShopOrderNotifyFlgUpdInfo", "base.NoReturn"},
	"RpcShopUploadAuth":              easygo.Pair{"share_message.PlayerID", "share_message.UploadAuthResult"},
	"RpcShopUploadAuthConfirm":       easygo.Pair{"share_message.PlayerID", "base.Empty"},
}

var DownRpc = map[string]easygo.Pair{
	"RpcPlayerLoginResponse": easygo.Pair{"client_server.AllPlayerMsg", "base.NoReturn"},
}
