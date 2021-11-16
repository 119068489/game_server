package brower_backstage

import (
	"game_server/easygo"
	"game_server/easygo/base"
)

type _ = base.NoReturn

type IBrower2Wish interface {
	RpcBrowerTest(reqMsg *base.Empty) *base.Empty
	RpcBrowerTest_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryWishBoxList(reqMsg *WishBoxListRequest) *WishBoxList
	RpcQueryWishBoxList_(reqMsg *WishBoxListRequest) (*WishBoxList, easygo.IRpcInterrupt)
	RpcUpdateWishBox(reqMsg *WishBox) *base.Empty
	RpcUpdateWishBox_(reqMsg *WishBox) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryWishBoxGoodsItemList(reqMsg *ListRequest) *WishBoxGoodsItemList
	RpcQueryWishBoxGoodsItemList_(reqMsg *ListRequest) (*WishBoxGoodsItemList, easygo.IRpcInterrupt)
	RpcQueryWishBoxWinCfgList(reqMsg *ListRequest) *WishBoxWinCfgList
	RpcQueryWishBoxWinCfgList_(reqMsg *ListRequest) (*WishBoxWinCfgList, easygo.IRpcInterrupt)
	RpcUpdateWishBoxWinCfgList(reqMsg *WishBoxWinCfgList) *base.Empty
	RpcUpdateWishBoxWinCfgList_(reqMsg *WishBoxWinCfgList) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryWishGoodsList(reqMsg *WishBoxGoodsListRequest) *WishBoxGoods
	RpcQueryWishGoodsList_(reqMsg *WishBoxGoodsListRequest) (*WishBoxGoods, easygo.IRpcInterrupt)
	RpcUpdateWishGoods(reqMsg *WishBoxGoods) *base.Empty
	RpcUpdateWishGoods_(reqMsg *WishBoxGoods) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryWishGoodsBrandList(reqMsg *ListRequest) *WishGoodsBrandList
	RpcQueryWishGoodsBrandList_(reqMsg *ListRequest) (*WishGoodsBrandList, easygo.IRpcInterrupt)
	RpcUpdateWishGoodsBrand(reqMsg *WishGoodsBrand) *base.Empty
	RpcUpdateWishGoodsBrand_(reqMsg *WishGoodsBrand) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryWishGoodsTypeList(reqMsg *ListRequest) *WishGoodsTypeList
	RpcQueryWishGoodsTypeList_(reqMsg *ListRequest) (*WishGoodsTypeList, easygo.IRpcInterrupt)
	RpcUpdateWishGoodsType(reqMsg *WishGoodsType) *base.Empty
	RpcUpdateWishGoodsType_(reqMsg *WishGoodsType) (*base.Empty, easygo.IRpcInterrupt)
	RpcQueryWishDeliveryOrderList(reqMsg *ListRequest) *WishDeliveryOrderList
	RpcQueryWishDeliveryOrderList_(reqMsg *ListRequest) (*WishDeliveryOrderList, easygo.IRpcInterrupt)
	RpcQueryWishRecycleOrderList(reqMsg *ListRequest) *WishRecycleOrderList
	RpcQueryWishRecycleOrderList_(reqMsg *ListRequest) (*WishRecycleOrderList, easygo.IRpcInterrupt)
}

type Brower2Wish struct {
	Sender easygo.IMessageSender
}

func (self *Brower2Wish) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Brower2Wish) RpcBrowerTest(reqMsg *base.Empty) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcBrowerTest", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Wish) RpcBrowerTest_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBrowerTest", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Wish) RpcQueryWishBoxList(reqMsg *WishBoxListRequest) *WishBoxList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBoxList)
}

func (self *Brower2Wish) RpcQueryWishBoxList_(reqMsg *WishBoxListRequest) (*WishBoxList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBoxList), e
}
func (self *Brower2Wish) RpcUpdateWishBox(reqMsg *WishBox) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishBox", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Wish) RpcUpdateWishBox_(reqMsg *WishBox) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishBox", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Wish) RpcQueryWishBoxGoodsItemList(reqMsg *ListRequest) *WishBoxGoodsItemList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxGoodsItemList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBoxGoodsItemList)
}

func (self *Brower2Wish) RpcQueryWishBoxGoodsItemList_(reqMsg *ListRequest) (*WishBoxGoodsItemList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxGoodsItemList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBoxGoodsItemList), e
}
func (self *Brower2Wish) RpcQueryWishBoxWinCfgList(reqMsg *ListRequest) *WishBoxWinCfgList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxWinCfgList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBoxWinCfgList)
}

func (self *Brower2Wish) RpcQueryWishBoxWinCfgList_(reqMsg *ListRequest) (*WishBoxWinCfgList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishBoxWinCfgList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBoxWinCfgList), e
}
func (self *Brower2Wish) RpcUpdateWishBoxWinCfgList(reqMsg *WishBoxWinCfgList) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishBoxWinCfgList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Wish) RpcUpdateWishBoxWinCfgList_(reqMsg *WishBoxWinCfgList) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishBoxWinCfgList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Wish) RpcQueryWishGoodsList(reqMsg *WishBoxGoodsListRequest) *WishBoxGoods {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishGoodsList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishBoxGoods)
}

func (self *Brower2Wish) RpcQueryWishGoodsList_(reqMsg *WishBoxGoodsListRequest) (*WishBoxGoods, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishGoodsList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishBoxGoods), e
}
func (self *Brower2Wish) RpcUpdateWishGoods(reqMsg *WishBoxGoods) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishGoods", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Wish) RpcUpdateWishGoods_(reqMsg *WishBoxGoods) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishGoods", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Wish) RpcQueryWishGoodsBrandList(reqMsg *ListRequest) *WishGoodsBrandList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishGoodsBrandList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishGoodsBrandList)
}

func (self *Brower2Wish) RpcQueryWishGoodsBrandList_(reqMsg *ListRequest) (*WishGoodsBrandList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishGoodsBrandList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishGoodsBrandList), e
}
func (self *Brower2Wish) RpcUpdateWishGoodsBrand(reqMsg *WishGoodsBrand) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishGoodsBrand", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Wish) RpcUpdateWishGoodsBrand_(reqMsg *WishGoodsBrand) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishGoodsBrand", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Wish) RpcQueryWishGoodsTypeList(reqMsg *ListRequest) *WishGoodsTypeList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishGoodsTypeList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishGoodsTypeList)
}

func (self *Brower2Wish) RpcQueryWishGoodsTypeList_(reqMsg *ListRequest) (*WishGoodsTypeList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishGoodsTypeList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishGoodsTypeList), e
}
func (self *Brower2Wish) RpcUpdateWishGoodsType(reqMsg *WishGoodsType) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishGoodsType", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Brower2Wish) RpcUpdateWishGoodsType_(reqMsg *WishGoodsType) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcUpdateWishGoodsType", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *Brower2Wish) RpcQueryWishDeliveryOrderList(reqMsg *ListRequest) *WishDeliveryOrderList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishDeliveryOrderList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishDeliveryOrderList)
}

func (self *Brower2Wish) RpcQueryWishDeliveryOrderList_(reqMsg *ListRequest) (*WishDeliveryOrderList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishDeliveryOrderList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishDeliveryOrderList), e
}
func (self *Brower2Wish) RpcQueryWishRecycleOrderList(reqMsg *ListRequest) *WishRecycleOrderList {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishRecycleOrderList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishRecycleOrderList)
}

func (self *Brower2Wish) RpcQueryWishRecycleOrderList_(reqMsg *ListRequest) (*WishRecycleOrderList, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryWishRecycleOrderList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishRecycleOrderList), e
}

// ==========================================================
type IBackstage2Wish interface {
	RpcBackstageTest(reqMsg *base.Empty) *base.Empty
	RpcBackstageTest_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt)
}

type Backstage2Wish struct {
	Sender easygo.IMessageSender
}

func (self *Backstage2Wish) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *Backstage2Wish) RpcBackstageTest(reqMsg *base.Empty) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcBackstageTest", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *Backstage2Wish) RpcBackstageTest_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBackstageTest", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
