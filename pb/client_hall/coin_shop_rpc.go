package client_hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/share_message"
)

type _ = base.NoReturn

type ICoinShopClient2Hall interface {
	RpcGetPropsItems(reqMsg *base.Empty) *PropsItemList
	RpcGetPropsItems_(reqMsg *base.Empty) (*PropsItemList, easygo.IRpcInterrupt)
	RpcGetCoinRechargeList(reqMsg *CoinRechargeList) *CoinRechargeList
	RpcGetCoinRechargeList_(reqMsg *CoinRechargeList) (*CoinRechargeList, easygo.IRpcInterrupt)
	RpcGetCoinShopList(reqMsg *CoinShopList) *CoinShopList
	RpcGetCoinShopList_(reqMsg *CoinShopList) (*CoinShopList, easygo.IRpcInterrupt)
	RpcCoinRecharge(reqMsg *CoinRechargeReq) *CoinRechargeResp
	RpcCoinRecharge_(reqMsg *CoinRechargeReq) (*CoinRechargeResp, easygo.IRpcInterrupt)
	RpcBuyCoinItem(reqMsg *BuyCoinItem) *BuyCoinItem
	RpcBuyCoinItem_(reqMsg *BuyCoinItem) (*BuyCoinItem, easygo.IRpcInterrupt)
	RpcUseCoinItem(reqMsg *UseCoinItem) *UseCoinItem
	RpcUseCoinItem_(reqMsg *UseCoinItem) (*UseCoinItem, easygo.IRpcInterrupt)
	RpcGetPlayerEquipment(reqMsg *EquipmentReq) *EquipmentReq
	RpcGetPlayerEquipment_(reqMsg *EquipmentReq) (*EquipmentReq, easygo.IRpcInterrupt)
	RpcGetPlayerBagItems(reqMsg *BagItems) *BagItems
	RpcGetPlayerBagItems_(reqMsg *BagItems) (*BagItems, easygo.IRpcInterrupt)
	RpcCoinRechargeAct(reqMsg *RechargeActReq) *CoinRechargeResp
	RpcCoinRechargeAct_(reqMsg *RechargeActReq) (*CoinRechargeResp, easygo.IRpcInterrupt)
}

type CoinShopClient2Hall struct {
	Sender easygo.IMessageSender
}

func (self *CoinShopClient2Hall) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *CoinShopClient2Hall) RpcGetPropsItems(reqMsg *base.Empty) *PropsItemList {
	msg, e := self.Sender.CallRpcMethod("RpcGetPropsItems", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PropsItemList)
}

func (self *CoinShopClient2Hall) RpcGetPropsItems_(reqMsg *base.Empty) (*PropsItemList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPropsItems", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PropsItemList), e
}
func (self *CoinShopClient2Hall) RpcGetCoinRechargeList(reqMsg *CoinRechargeList) *CoinRechargeList {
	msg, e := self.Sender.CallRpcMethod("RpcGetCoinRechargeList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CoinRechargeList)
}

func (self *CoinShopClient2Hall) RpcGetCoinRechargeList_(reqMsg *CoinRechargeList) (*CoinRechargeList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetCoinRechargeList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CoinRechargeList), e
}
func (self *CoinShopClient2Hall) RpcGetCoinShopList(reqMsg *CoinShopList) *CoinShopList {
	msg, e := self.Sender.CallRpcMethod("RpcGetCoinShopList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CoinShopList)
}

func (self *CoinShopClient2Hall) RpcGetCoinShopList_(reqMsg *CoinShopList) (*CoinShopList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetCoinShopList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CoinShopList), e
}
func (self *CoinShopClient2Hall) RpcCoinRecharge(reqMsg *CoinRechargeReq) *CoinRechargeResp {
	msg, e := self.Sender.CallRpcMethod("RpcCoinRecharge", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CoinRechargeResp)
}

func (self *CoinShopClient2Hall) RpcCoinRecharge_(reqMsg *CoinRechargeReq) (*CoinRechargeResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCoinRecharge", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CoinRechargeResp), e
}
func (self *CoinShopClient2Hall) RpcBuyCoinItem(reqMsg *BuyCoinItem) *BuyCoinItem {
	msg, e := self.Sender.CallRpcMethod("RpcBuyCoinItem", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BuyCoinItem)
}

func (self *CoinShopClient2Hall) RpcBuyCoinItem_(reqMsg *BuyCoinItem) (*BuyCoinItem, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBuyCoinItem", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BuyCoinItem), e
}
func (self *CoinShopClient2Hall) RpcUseCoinItem(reqMsg *UseCoinItem) *UseCoinItem {
	msg, e := self.Sender.CallRpcMethod("RpcUseCoinItem", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*UseCoinItem)
}

func (self *CoinShopClient2Hall) RpcUseCoinItem_(reqMsg *UseCoinItem) (*UseCoinItem, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUseCoinItem", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*UseCoinItem), e
}
func (self *CoinShopClient2Hall) RpcGetPlayerEquipment(reqMsg *EquipmentReq) *EquipmentReq {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerEquipment", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*EquipmentReq)
}

func (self *CoinShopClient2Hall) RpcGetPlayerEquipment_(reqMsg *EquipmentReq) (*EquipmentReq, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerEquipment", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*EquipmentReq), e
}
func (self *CoinShopClient2Hall) RpcGetPlayerBagItems(reqMsg *BagItems) *BagItems {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerBagItems", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BagItems)
}

func (self *CoinShopClient2Hall) RpcGetPlayerBagItems_(reqMsg *BagItems) (*BagItems, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPlayerBagItems", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BagItems), e
}
func (self *CoinShopClient2Hall) RpcCoinRechargeAct(reqMsg *RechargeActReq) *CoinRechargeResp {
	msg, e := self.Sender.CallRpcMethod("RpcCoinRechargeAct", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CoinRechargeResp)
}

func (self *CoinShopClient2Hall) RpcCoinRechargeAct_(reqMsg *RechargeActReq) (*CoinRechargeResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCoinRechargeAct", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CoinRechargeResp), e
}

// ==========================================================
type IHall2CoinShopClient interface {
	RpcDelBagItem(reqMsg *share_message.PlayerBagItem)
	RpcModifyBagItem(reqMsg *BagItems)
	RpcModifyEquipment(reqMsg *EquipmentReq)
	RpcNewBagItemsTip(reqMsg *NewBagItemsTip)
}

type Hall2CoinShopClient struct {
	Sender easygo.IMessageSender
}

func (self *Hall2CoinShopClient) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Hall2CoinShopClient) RpcDelBagItem(reqMsg *share_message.PlayerBagItem) {
	self.Sender.CallRpcMethod("RpcDelBagItem", reqMsg)
}
func (self *Hall2CoinShopClient) RpcModifyBagItem(reqMsg *BagItems) {
	self.Sender.CallRpcMethod("RpcModifyBagItem", reqMsg)
}
func (self *Hall2CoinShopClient) RpcModifyEquipment(reqMsg *EquipmentReq) {
	self.Sender.CallRpcMethod("RpcModifyEquipment", reqMsg)
}
func (self *Hall2CoinShopClient) RpcNewBagItemsTip(reqMsg *NewBagItemsTip) {
	self.Sender.CallRpcMethod("RpcNewBagItemsTip", reqMsg)
}
