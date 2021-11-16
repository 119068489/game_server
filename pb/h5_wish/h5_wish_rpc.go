package h5_wish

import (
	"game_server/easygo"
	"game_server/easygo/base"
)

type _ = base.NoReturn

type IWishClient2Hall interface {
	RpcLogin(reqMsg *LoginReq) *LoginResp
	RpcLogin_(reqMsg *LoginReq) (*LoginResp, easygo.IRpcInterrupt)
	RpcQueryBox(reqMsg *QueryBoxReq) *QueryBoxResp
	RpcQueryBox_(reqMsg *QueryBoxReq) (*QueryBoxResp, easygo.IRpcInterrupt)
	RpcQueryBoxProductName(reqMsg *base.Empty) *BoxProductNameResp
	RpcQueryBoxProductName_(reqMsg *base.Empty) (*BoxProductNameResp, easygo.IRpcInterrupt)
	RpcSearchFound(reqMsg *base.Empty) *BoxProductNameResp
	RpcSearchFound_(reqMsg *base.Empty) (*BoxProductNameResp, easygo.IRpcInterrupt)
	RpcProductShow(reqMsg *base.Empty) *ProductShowResp
	RpcProductShow_(reqMsg *base.Empty) (*ProductShowResp, easygo.IRpcInterrupt)
	RpcGetCoin(reqMsg *base.Empty) *GetCoinResp
	RpcGetCoin_(reqMsg *base.Empty) (*GetCoinResp, easygo.IRpcInterrupt)
	RpcGetUseInfo(reqMsg *UserInfoReq) *UserInfoResp
	RpcGetUseInfo_(reqMsg *UserInfoReq) (*UserInfoResp, easygo.IRpcInterrupt)
	RpcAddCoin(reqMsg *AddCoinReq) *AddCoinResp
	RpcAddCoin_(reqMsg *AddCoinReq) (*AddCoinResp, easygo.IRpcInterrupt)
	RpcAddGold(reqMsg *AddGoldReq) *AddGoldResp
	RpcAddGold_(reqMsg *AddGoldReq) (*AddGoldResp, easygo.IRpcInterrupt)
	RpcHomeMessage(reqMsg *base.Empty) *HomeMessageResp
	RpcHomeMessage_(reqMsg *base.Empty) (*HomeMessageResp, easygo.IRpcInterrupt)
	RpcProtector(reqMsg *base.Empty) *ProtectorDataResp
	RpcProtector_(reqMsg *base.Empty) (*ProtectorDataResp, easygo.IRpcInterrupt)
	RpcMenu(reqMsg *base.Empty) *MenuResp
	RpcMenu_(reqMsg *base.Empty) (*MenuResp, easygo.IRpcInterrupt)
	RpcProductBrand(reqMsg *base.Empty) *ProductBrandListResp
	RpcProductBrand_(reqMsg *base.Empty) (*ProductBrandListResp, easygo.IRpcInterrupt)
	RpcSearchBox(reqMsg *SearchBoxReq) *SearchBoxResp
	RpcSearchBox_(reqMsg *SearchBoxReq) (*SearchBoxResp, easygo.IRpcInterrupt)
	RpcBrandList(reqMsg *base.Empty) *BrandListResp
	RpcBrandList_(reqMsg *base.Empty) (*BrandListResp, easygo.IRpcInterrupt)
	RpcProductTypeList(reqMsg *base.Empty) *TypeListResp
	RpcProductTypeList_(reqMsg *base.Empty) (*TypeListResp, easygo.IRpcInterrupt)
	RpcGetRandProduct(reqMsg *DareReq) *RandProductResp
	RpcGetRandProduct_(reqMsg *DareReq) (*RandProductResp, easygo.IRpcInterrupt)
	RpcGetDareMessage(reqMsg *base.Empty) *DareMessageResp
	RpcGetDareMessage_(reqMsg *base.Empty) (*DareMessageResp, easygo.IRpcInterrupt)
	RpcDefenderCarousel(reqMsg *base.Empty) *DefenderMsgResp
	RpcDefenderCarousel_(reqMsg *base.Empty) (*DefenderMsgResp, easygo.IRpcInterrupt)
	RpcGotWishCarousel(reqMsg *base.Empty) *GotWishPlayerResp
	RpcGotWishCarousel_(reqMsg *base.Empty) (*GotWishPlayerResp, easygo.IRpcInterrupt)
	RpcDareRecommend(reqMsg *base.Empty) *DareRecommendResp
	RpcDareRecommend_(reqMsg *base.Empty) (*DareRecommendResp, easygo.IRpcInterrupt)
	RpcRankings(reqMsg *base.Empty) *RankingResp
	RpcRankings_(reqMsg *base.Empty) (*RankingResp, easygo.IRpcInterrupt)
	RpcMyRecord(reqMsg *base.Empty) *MyRecordResp
	RpcMyRecord_(reqMsg *base.Empty) (*MyRecordResp, easygo.IRpcInterrupt)
	RpcMyDare(reqMsg *MyDareReq) *MyDareResp
	RpcMyDare_(reqMsg *MyDareReq) (*MyDareResp, easygo.IRpcInterrupt)
	RpcBoxInfo(reqMsg *BoxReq) *BoxResp
	RpcBoxInfo_(reqMsg *BoxReq) (*BoxResp, easygo.IRpcInterrupt)
	RpcDareList(reqMsg *DareReq) *DareResp
	RpcDareList_(reqMsg *DareReq) (*DareResp, easygo.IRpcInterrupt)
	RpcProductDetail(reqMsg *ProductDetailReq) *ProductDetail
	RpcProductDetail_(reqMsg *ProductDetailReq) (*ProductDetail, easygo.IRpcInterrupt)
	RpcDareRecord(reqMsg *DareRecordReq) *DareRecordResp
	RpcDareRecord_(reqMsg *DareRecordReq) (*DareRecordResp, easygo.IRpcInterrupt)
	RpcWish(reqMsg *WishReq) *WishResp
	RpcWish_(reqMsg *WishReq) (*WishResp, easygo.IRpcInterrupt)
	RpcDoDare(reqMsg *DoDareReq) *DoDareResp
	RpcDoDare_(reqMsg *DoDareReq) (*DoDareResp, easygo.IRpcInterrupt)
	RpcBoxList(reqMsg *base.Empty) *BoxListResp
	RpcBoxList_(reqMsg *base.Empty) (*BoxListResp, easygo.IRpcInterrupt)
	RpcBoxProduct(reqMsg *DareReq) *BoxProductResp
	RpcBoxProduct_(reqMsg *DareReq) (*BoxProductResp, easygo.IRpcInterrupt)
	RpcTryOnce(reqMsg *base.Empty) *base.Empty
	RpcTryOnce_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt)
	RpcMyWish(reqMsg *MyWishReq) *ProductResp
	RpcMyWish_(reqMsg *MyWishReq) (*ProductResp, easygo.IRpcInterrupt)
	RpcMyAllWish(reqMsg *MyWishReq) *ProductResp
	RpcMyAllWish_(reqMsg *MyWishReq) (*ProductResp, easygo.IRpcInterrupt)
	RpcGetCollectionBox(reqMsg *DataPageReq) *MyCollectedBoxResp
	RpcGetCollectionBox_(reqMsg *DataPageReq) (*MyCollectedBoxResp, easygo.IRpcInterrupt)
	RpcGetAllCollectionBox(reqMsg *base.Empty) *MyCollectedBoxResp
	RpcGetAllCollectionBox_(reqMsg *base.Empty) (*MyCollectedBoxResp, easygo.IRpcInterrupt)
	RpcCollectionBox(reqMsg *CollectionBoxReq) *DefaultResp
	RpcCollectionBox_(reqMsg *CollectionBoxReq) (*DefaultResp, easygo.IRpcInterrupt)
	RpcGetWishBoxList(reqMsg *DataPageReq) *MyCollectedBoxResp
	RpcGetWishBoxList_(reqMsg *DataPageReq) (*MyCollectedBoxResp, easygo.IRpcInterrupt)
	RpcGetAllWishBoxList(reqMsg *base.Empty) *MyCollectedBoxResp
	RpcGetAllWishBoxList_(reqMsg *base.Empty) (*MyCollectedBoxResp, easygo.IRpcInterrupt)
	RpcDelWishBox(reqMsg *WishBoxReq) *DefaultResp
	RpcDelWishBox_(reqMsg *WishBoxReq) (*DefaultResp, easygo.IRpcInterrupt)
	RpcExchangeBox(reqMsg *WishBoxReq) *DefaultResp
	RpcExchangeBox_(reqMsg *WishBoxReq) (*DefaultResp, easygo.IRpcInterrupt)
	RpcRecycleGoods(reqMsg *WishBoxReq) *DefaultResp
	RpcRecycleGoods_(reqMsg *WishBoxReq) (*DefaultResp, easygo.IRpcInterrupt)
	RpcGetAddressList(reqMsg *DataPageReq) *AddressListResp
	RpcGetAddressList_(reqMsg *DataPageReq) (*AddressListResp, easygo.IRpcInterrupt)
	RpcAddAddress(reqMsg *WishAddress) *DefaultResp
	RpcAddAddress_(reqMsg *WishAddress) (*DefaultResp, easygo.IRpcInterrupt)
	RpcEditAddress(reqMsg *WishAddress) *DefaultResp
	RpcEditAddress_(reqMsg *WishAddress) (*DefaultResp, easygo.IRpcInterrupt)
	RpcRemoveAddress(reqMsg *RemoveAddressReq) *DefaultResp
	RpcRemoveAddress_(reqMsg *RemoveAddressReq) (*DefaultResp, easygo.IRpcInterrupt)
	RpcGetUnReadWishNum(reqMsg *base.Empty) *JustNumberResp
	RpcGetUnReadWishNum_(reqMsg *base.Empty) (*JustNumberResp, easygo.IRpcInterrupt)
	RpcToExchangeWishNum(reqMsg *base.Empty) *JustNumberResp
	RpcToExchangeWishNum_(reqMsg *base.Empty) (*JustNumberResp, easygo.IRpcInterrupt)
	RpcAreaPostage(reqMsg *base.Empty) *PostageResp
	RpcAreaPostage_(reqMsg *base.Empty) (*PostageResp, easygo.IRpcInterrupt)
	RpcRecycleRatio(reqMsg *base.Empty) *JustNumberResp
	RpcRecycleRatio_(reqMsg *base.Empty) (*JustNumberResp, easygo.IRpcInterrupt)
	RpcGetConfig(reqMsg *base.Empty) *ConfigResp
	RpcGetConfig_(reqMsg *base.Empty) (*ConfigResp, easygo.IRpcInterrupt)
	RpcGetUserIdBankCards(reqMsg *base.Empty) *BankCardResp
	RpcGetUserIdBankCards_(reqMsg *base.Empty) (*BankCardResp, easygo.IRpcInterrupt)
	RpcRecycleDesc(reqMsg *base.Empty) *DefaultResp
	RpcRecycleDesc_(reqMsg *base.Empty) (*DefaultResp, easygo.IRpcInterrupt)
	RpcSetBoxPreSale(reqMsg *PresaleReq) *base.Empty
	RpcSetBoxPreSale_(reqMsg *PresaleReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcSoldOutBox(reqMsg *DealBoxReq) *base.Empty
	RpcSoldOutBox_(reqMsg *DealBoxReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcReplenishBox(reqMsg *DealBoxReq) *base.Empty
	RpcReplenishBox_(reqMsg *DealBoxReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcSetBoxExpress(reqMsg *SetExpressInfoReq) *base.Empty
	RpcSetBoxExpress_(reqMsg *SetExpressInfoReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetBoxExpress(reqMsg *GetExpressInfoReq) *GetExpressInfoResp
	RpcGetBoxExpress_(reqMsg *GetExpressInfoReq) (*GetExpressInfoResp, easygo.IRpcInterrupt)
	RpcCoinToDiamond(reqMsg *CoinToDiamondReq) *CoinToDiamondResq
	RpcCoinToDiamond_(reqMsg *CoinToDiamondReq) (*CoinToDiamondResq, easygo.IRpcInterrupt)
	RpcDiamondRechargeList(reqMsg *base.Empty) *DiamondRechargeResp
	RpcDiamondRechargeList_(reqMsg *base.Empty) (*DiamondRechargeResp, easygo.IRpcInterrupt)
	RpcDiamondChangeLogList(reqMsg *DiamondChangeLogReq) *DiamondChangeLogResp
	RpcDiamondChangeLogList_(reqMsg *DiamondChangeLogReq) (*DiamondChangeLogResp, easygo.IRpcInterrupt)
	RpcGetPriceSection(reqMsg *base.Empty) *PriceSectionResp
	RpcGetPriceSection_(reqMsg *base.Empty) (*PriceSectionResp, easygo.IRpcInterrupt)
	RpcBatchDare(reqMsg *BatchDareReq) *BatchDareResp
	RpcBatchDare_(reqMsg *BatchDareReq) (*BatchDareResp, easygo.IRpcInterrupt)
	RpcPlayCfg(reqMsg *base.Empty) *PlayCfgResp
	RpcPlayCfg_(reqMsg *base.Empty) (*PlayCfgResp, easygo.IRpcInterrupt)
	RpcSumNum(reqMsg *SumReq) *SumNumResp
	RpcSumNum_(reqMsg *SumReq) (*SumNumResp, easygo.IRpcInterrupt)
	RpcSumMoney(reqMsg *SumMoneyReq) *SumMoneyResp
	RpcSumMoney_(reqMsg *SumMoneyReq) (*SumMoneyResp, easygo.IRpcInterrupt)
	RpcGive(reqMsg *GiveReq) *GiveResp
	RpcGive_(reqMsg *GiveReq) (*GiveResp, easygo.IRpcInterrupt)
	RpcActPoolList(reqMsg *base.Empty) *ActPoolResp
	RpcActPoolList_(reqMsg *base.Empty) (*ActPoolResp, easygo.IRpcInterrupt)
	RpcActPoolRule(reqMsg *ActPoolRuleReq) *ActPoolRuleResp
	RpcActPoolRule_(reqMsg *ActPoolRuleReq) (*ActPoolRuleResp, easygo.IRpcInterrupt)
	RpcActName(reqMsg *ActNameReq) *ActNameResp
	RpcActName_(reqMsg *ActNameReq) (*ActNameResp, easygo.IRpcInterrupt)
	RpcActOpenStatus(reqMsg *base.Empty) *ActOpenStatusResp
	RpcActOpenStatus_(reqMsg *base.Empty) (*ActOpenStatusResp, easygo.IRpcInterrupt)
	RpcRechargeActStatus(reqMsg *base.Empty) *ActOpenStatusResp
	RpcRechargeActStatus_(reqMsg *base.Empty) (*ActOpenStatusResp, easygo.IRpcInterrupt)
	RpcReportWishLog(reqMsg *TypeReq) *base.Empty
	RpcReportWishLog_(reqMsg *TypeReq) (*base.Empty, easygo.IRpcInterrupt)
	RpcGetRechargeAct(reqMsg *base.Empty) *RechargeActResp
	RpcGetRechargeAct_(reqMsg *base.Empty) (*RechargeActResp, easygo.IRpcInterrupt)
	RpcGetRechargeLogs(reqMsg *DataPageReq) *RechargeLogsResp
	RpcGetRechargeLogs_(reqMsg *DataPageReq) (*RechargeLogsResp, easygo.IRpcInterrupt)
}

type WishClient2Hall struct {
	Sender easygo.IMessageSender
}

func (self *WishClient2Hall) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------

func (self *WishClient2Hall) RpcLogin(reqMsg *LoginReq) *LoginResp {
	msg, e := self.Sender.CallRpcMethod("RpcLogin", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*LoginResp)
}

func (self *WishClient2Hall) RpcLogin_(reqMsg *LoginReq) (*LoginResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcLogin", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*LoginResp), e
}
func (self *WishClient2Hall) RpcQueryBox(reqMsg *QueryBoxReq) *QueryBoxResp {
	msg, e := self.Sender.CallRpcMethod("RpcQueryBox", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*QueryBoxResp)
}

func (self *WishClient2Hall) RpcQueryBox_(reqMsg *QueryBoxReq) (*QueryBoxResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryBox", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*QueryBoxResp), e
}
func (self *WishClient2Hall) RpcQueryBoxProductName(reqMsg *base.Empty) *BoxProductNameResp {
	msg, e := self.Sender.CallRpcMethod("RpcQueryBoxProductName", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BoxProductNameResp)
}

func (self *WishClient2Hall) RpcQueryBoxProductName_(reqMsg *base.Empty) (*BoxProductNameResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcQueryBoxProductName", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BoxProductNameResp), e
}
func (self *WishClient2Hall) RpcSearchFound(reqMsg *base.Empty) *BoxProductNameResp {
	msg, e := self.Sender.CallRpcMethod("RpcSearchFound", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BoxProductNameResp)
}

func (self *WishClient2Hall) RpcSearchFound_(reqMsg *base.Empty) (*BoxProductNameResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSearchFound", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BoxProductNameResp), e
}
func (self *WishClient2Hall) RpcProductShow(reqMsg *base.Empty) *ProductShowResp {
	msg, e := self.Sender.CallRpcMethod("RpcProductShow", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ProductShowResp)
}

func (self *WishClient2Hall) RpcProductShow_(reqMsg *base.Empty) (*ProductShowResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcProductShow", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ProductShowResp), e
}
func (self *WishClient2Hall) RpcGetCoin(reqMsg *base.Empty) *GetCoinResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetCoin", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GetCoinResp)
}

func (self *WishClient2Hall) RpcGetCoin_(reqMsg *base.Empty) (*GetCoinResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetCoin", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GetCoinResp), e
}
func (self *WishClient2Hall) RpcGetUseInfo(reqMsg *UserInfoReq) *UserInfoResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetUseInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*UserInfoResp)
}

func (self *WishClient2Hall) RpcGetUseInfo_(reqMsg *UserInfoReq) (*UserInfoResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetUseInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*UserInfoResp), e
}
func (self *WishClient2Hall) RpcAddCoin(reqMsg *AddCoinReq) *AddCoinResp {
	msg, e := self.Sender.CallRpcMethod("RpcAddCoin", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AddCoinResp)
}

func (self *WishClient2Hall) RpcAddCoin_(reqMsg *AddCoinReq) (*AddCoinResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddCoin", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AddCoinResp), e
}
func (self *WishClient2Hall) RpcAddGold(reqMsg *AddGoldReq) *AddGoldResp {
	msg, e := self.Sender.CallRpcMethod("RpcAddGold", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AddGoldResp)
}

func (self *WishClient2Hall) RpcAddGold_(reqMsg *AddGoldReq) (*AddGoldResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddGold", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AddGoldResp), e
}
func (self *WishClient2Hall) RpcHomeMessage(reqMsg *base.Empty) *HomeMessageResp {
	msg, e := self.Sender.CallRpcMethod("RpcHomeMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*HomeMessageResp)
}

func (self *WishClient2Hall) RpcHomeMessage_(reqMsg *base.Empty) (*HomeMessageResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcHomeMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*HomeMessageResp), e
}
func (self *WishClient2Hall) RpcProtector(reqMsg *base.Empty) *ProtectorDataResp {
	msg, e := self.Sender.CallRpcMethod("RpcProtector", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ProtectorDataResp)
}

func (self *WishClient2Hall) RpcProtector_(reqMsg *base.Empty) (*ProtectorDataResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcProtector", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ProtectorDataResp), e
}
func (self *WishClient2Hall) RpcMenu(reqMsg *base.Empty) *MenuResp {
	msg, e := self.Sender.CallRpcMethod("RpcMenu", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MenuResp)
}

func (self *WishClient2Hall) RpcMenu_(reqMsg *base.Empty) (*MenuResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcMenu", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MenuResp), e
}
func (self *WishClient2Hall) RpcProductBrand(reqMsg *base.Empty) *ProductBrandListResp {
	msg, e := self.Sender.CallRpcMethod("RpcProductBrand", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ProductBrandListResp)
}

func (self *WishClient2Hall) RpcProductBrand_(reqMsg *base.Empty) (*ProductBrandListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcProductBrand", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ProductBrandListResp), e
}
func (self *WishClient2Hall) RpcSearchBox(reqMsg *SearchBoxReq) *SearchBoxResp {
	msg, e := self.Sender.CallRpcMethod("RpcSearchBox", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SearchBoxResp)
}

func (self *WishClient2Hall) RpcSearchBox_(reqMsg *SearchBoxReq) (*SearchBoxResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSearchBox", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SearchBoxResp), e
}
func (self *WishClient2Hall) RpcBrandList(reqMsg *base.Empty) *BrandListResp {
	msg, e := self.Sender.CallRpcMethod("RpcBrandList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BrandListResp)
}

func (self *WishClient2Hall) RpcBrandList_(reqMsg *base.Empty) (*BrandListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBrandList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BrandListResp), e
}
func (self *WishClient2Hall) RpcProductTypeList(reqMsg *base.Empty) *TypeListResp {
	msg, e := self.Sender.CallRpcMethod("RpcProductTypeList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*TypeListResp)
}

func (self *WishClient2Hall) RpcProductTypeList_(reqMsg *base.Empty) (*TypeListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcProductTypeList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*TypeListResp), e
}
func (self *WishClient2Hall) RpcGetRandProduct(reqMsg *DareReq) *RandProductResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetRandProduct", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*RandProductResp)
}

func (self *WishClient2Hall) RpcGetRandProduct_(reqMsg *DareReq) (*RandProductResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetRandProduct", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*RandProductResp), e
}
func (self *WishClient2Hall) RpcGetDareMessage(reqMsg *base.Empty) *DareMessageResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetDareMessage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DareMessageResp)
}

func (self *WishClient2Hall) RpcGetDareMessage_(reqMsg *base.Empty) (*DareMessageResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetDareMessage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DareMessageResp), e
}
func (self *WishClient2Hall) RpcDefenderCarousel(reqMsg *base.Empty) *DefenderMsgResp {
	msg, e := self.Sender.CallRpcMethod("RpcDefenderCarousel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DefenderMsgResp)
}

func (self *WishClient2Hall) RpcDefenderCarousel_(reqMsg *base.Empty) (*DefenderMsgResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDefenderCarousel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DefenderMsgResp), e
}
func (self *WishClient2Hall) RpcGotWishCarousel(reqMsg *base.Empty) *GotWishPlayerResp {
	msg, e := self.Sender.CallRpcMethod("RpcGotWishCarousel", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GotWishPlayerResp)
}

func (self *WishClient2Hall) RpcGotWishCarousel_(reqMsg *base.Empty) (*GotWishPlayerResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGotWishCarousel", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GotWishPlayerResp), e
}
func (self *WishClient2Hall) RpcDareRecommend(reqMsg *base.Empty) *DareRecommendResp {
	msg, e := self.Sender.CallRpcMethod("RpcDareRecommend", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DareRecommendResp)
}

func (self *WishClient2Hall) RpcDareRecommend_(reqMsg *base.Empty) (*DareRecommendResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDareRecommend", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DareRecommendResp), e
}
func (self *WishClient2Hall) RpcRankings(reqMsg *base.Empty) *RankingResp {
	msg, e := self.Sender.CallRpcMethod("RpcRankings", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*RankingResp)
}

func (self *WishClient2Hall) RpcRankings_(reqMsg *base.Empty) (*RankingResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRankings", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*RankingResp), e
}
func (self *WishClient2Hall) RpcMyRecord(reqMsg *base.Empty) *MyRecordResp {
	msg, e := self.Sender.CallRpcMethod("RpcMyRecord", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MyRecordResp)
}

func (self *WishClient2Hall) RpcMyRecord_(reqMsg *base.Empty) (*MyRecordResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcMyRecord", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MyRecordResp), e
}
func (self *WishClient2Hall) RpcMyDare(reqMsg *MyDareReq) *MyDareResp {
	msg, e := self.Sender.CallRpcMethod("RpcMyDare", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MyDareResp)
}

func (self *WishClient2Hall) RpcMyDare_(reqMsg *MyDareReq) (*MyDareResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcMyDare", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MyDareResp), e
}
func (self *WishClient2Hall) RpcBoxInfo(reqMsg *BoxReq) *BoxResp {
	msg, e := self.Sender.CallRpcMethod("RpcBoxInfo", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BoxResp)
}

func (self *WishClient2Hall) RpcBoxInfo_(reqMsg *BoxReq) (*BoxResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBoxInfo", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BoxResp), e
}
func (self *WishClient2Hall) RpcDareList(reqMsg *DareReq) *DareResp {
	msg, e := self.Sender.CallRpcMethod("RpcDareList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DareResp)
}

func (self *WishClient2Hall) RpcDareList_(reqMsg *DareReq) (*DareResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDareList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DareResp), e
}
func (self *WishClient2Hall) RpcProductDetail(reqMsg *ProductDetailReq) *ProductDetail {
	msg, e := self.Sender.CallRpcMethod("RpcProductDetail", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ProductDetail)
}

func (self *WishClient2Hall) RpcProductDetail_(reqMsg *ProductDetailReq) (*ProductDetail, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcProductDetail", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ProductDetail), e
}
func (self *WishClient2Hall) RpcDareRecord(reqMsg *DareRecordReq) *DareRecordResp {
	msg, e := self.Sender.CallRpcMethod("RpcDareRecord", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DareRecordResp)
}

func (self *WishClient2Hall) RpcDareRecord_(reqMsg *DareRecordReq) (*DareRecordResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDareRecord", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DareRecordResp), e
}
func (self *WishClient2Hall) RpcWish(reqMsg *WishReq) *WishResp {
	msg, e := self.Sender.CallRpcMethod("RpcWish", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*WishResp)
}

func (self *WishClient2Hall) RpcWish_(reqMsg *WishReq) (*WishResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcWish", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*WishResp), e
}
func (self *WishClient2Hall) RpcDoDare(reqMsg *DoDareReq) *DoDareResp {
	msg, e := self.Sender.CallRpcMethod("RpcDoDare", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DoDareResp)
}

func (self *WishClient2Hall) RpcDoDare_(reqMsg *DoDareReq) (*DoDareResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDoDare", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DoDareResp), e
}
func (self *WishClient2Hall) RpcBoxList(reqMsg *base.Empty) *BoxListResp {
	msg, e := self.Sender.CallRpcMethod("RpcBoxList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BoxListResp)
}

func (self *WishClient2Hall) RpcBoxList_(reqMsg *base.Empty) (*BoxListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBoxList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BoxListResp), e
}
func (self *WishClient2Hall) RpcBoxProduct(reqMsg *DareReq) *BoxProductResp {
	msg, e := self.Sender.CallRpcMethod("RpcBoxProduct", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BoxProductResp)
}

func (self *WishClient2Hall) RpcBoxProduct_(reqMsg *DareReq) (*BoxProductResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBoxProduct", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BoxProductResp), e
}
func (self *WishClient2Hall) RpcTryOnce(reqMsg *base.Empty) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcTryOnce", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *WishClient2Hall) RpcTryOnce_(reqMsg *base.Empty) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcTryOnce", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *WishClient2Hall) RpcMyWish(reqMsg *MyWishReq) *ProductResp {
	msg, e := self.Sender.CallRpcMethod("RpcMyWish", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ProductResp)
}

func (self *WishClient2Hall) RpcMyWish_(reqMsg *MyWishReq) (*ProductResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcMyWish", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ProductResp), e
}
func (self *WishClient2Hall) RpcMyAllWish(reqMsg *MyWishReq) *ProductResp {
	msg, e := self.Sender.CallRpcMethod("RpcMyAllWish", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ProductResp)
}

func (self *WishClient2Hall) RpcMyAllWish_(reqMsg *MyWishReq) (*ProductResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcMyAllWish", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ProductResp), e
}
func (self *WishClient2Hall) RpcGetCollectionBox(reqMsg *DataPageReq) *MyCollectedBoxResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetCollectionBox", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MyCollectedBoxResp)
}

func (self *WishClient2Hall) RpcGetCollectionBox_(reqMsg *DataPageReq) (*MyCollectedBoxResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetCollectionBox", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MyCollectedBoxResp), e
}
func (self *WishClient2Hall) RpcGetAllCollectionBox(reqMsg *base.Empty) *MyCollectedBoxResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllCollectionBox", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MyCollectedBoxResp)
}

func (self *WishClient2Hall) RpcGetAllCollectionBox_(reqMsg *base.Empty) (*MyCollectedBoxResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllCollectionBox", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MyCollectedBoxResp), e
}
func (self *WishClient2Hall) RpcCollectionBox(reqMsg *CollectionBoxReq) *DefaultResp {
	msg, e := self.Sender.CallRpcMethod("RpcCollectionBox", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DefaultResp)
}

func (self *WishClient2Hall) RpcCollectionBox_(reqMsg *CollectionBoxReq) (*DefaultResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCollectionBox", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DefaultResp), e
}
func (self *WishClient2Hall) RpcGetWishBoxList(reqMsg *DataPageReq) *MyCollectedBoxResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishBoxList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MyCollectedBoxResp)
}

func (self *WishClient2Hall) RpcGetWishBoxList_(reqMsg *DataPageReq) (*MyCollectedBoxResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetWishBoxList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MyCollectedBoxResp), e
}
func (self *WishClient2Hall) RpcGetAllWishBoxList(reqMsg *base.Empty) *MyCollectedBoxResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllWishBoxList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*MyCollectedBoxResp)
}

func (self *WishClient2Hall) RpcGetAllWishBoxList_(reqMsg *base.Empty) (*MyCollectedBoxResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetAllWishBoxList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*MyCollectedBoxResp), e
}
func (self *WishClient2Hall) RpcDelWishBox(reqMsg *WishBoxReq) *DefaultResp {
	msg, e := self.Sender.CallRpcMethod("RpcDelWishBox", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DefaultResp)
}

func (self *WishClient2Hall) RpcDelWishBox_(reqMsg *WishBoxReq) (*DefaultResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDelWishBox", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DefaultResp), e
}
func (self *WishClient2Hall) RpcExchangeBox(reqMsg *WishBoxReq) *DefaultResp {
	msg, e := self.Sender.CallRpcMethod("RpcExchangeBox", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DefaultResp)
}

func (self *WishClient2Hall) RpcExchangeBox_(reqMsg *WishBoxReq) (*DefaultResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcExchangeBox", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DefaultResp), e
}
func (self *WishClient2Hall) RpcRecycleGoods(reqMsg *WishBoxReq) *DefaultResp {
	msg, e := self.Sender.CallRpcMethod("RpcRecycleGoods", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DefaultResp)
}

func (self *WishClient2Hall) RpcRecycleGoods_(reqMsg *WishBoxReq) (*DefaultResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRecycleGoods", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DefaultResp), e
}
func (self *WishClient2Hall) RpcGetAddressList(reqMsg *DataPageReq) *AddressListResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetAddressList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*AddressListResp)
}

func (self *WishClient2Hall) RpcGetAddressList_(reqMsg *DataPageReq) (*AddressListResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetAddressList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*AddressListResp), e
}
func (self *WishClient2Hall) RpcAddAddress(reqMsg *WishAddress) *DefaultResp {
	msg, e := self.Sender.CallRpcMethod("RpcAddAddress", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DefaultResp)
}

func (self *WishClient2Hall) RpcAddAddress_(reqMsg *WishAddress) (*DefaultResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAddAddress", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DefaultResp), e
}
func (self *WishClient2Hall) RpcEditAddress(reqMsg *WishAddress) *DefaultResp {
	msg, e := self.Sender.CallRpcMethod("RpcEditAddress", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DefaultResp)
}

func (self *WishClient2Hall) RpcEditAddress_(reqMsg *WishAddress) (*DefaultResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcEditAddress", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DefaultResp), e
}
func (self *WishClient2Hall) RpcRemoveAddress(reqMsg *RemoveAddressReq) *DefaultResp {
	msg, e := self.Sender.CallRpcMethod("RpcRemoveAddress", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DefaultResp)
}

func (self *WishClient2Hall) RpcRemoveAddress_(reqMsg *RemoveAddressReq) (*DefaultResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRemoveAddress", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DefaultResp), e
}
func (self *WishClient2Hall) RpcGetUnReadWishNum(reqMsg *base.Empty) *JustNumberResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetUnReadWishNum", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*JustNumberResp)
}

func (self *WishClient2Hall) RpcGetUnReadWishNum_(reqMsg *base.Empty) (*JustNumberResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetUnReadWishNum", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*JustNumberResp), e
}
func (self *WishClient2Hall) RpcToExchangeWishNum(reqMsg *base.Empty) *JustNumberResp {
	msg, e := self.Sender.CallRpcMethod("RpcToExchangeWishNum", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*JustNumberResp)
}

func (self *WishClient2Hall) RpcToExchangeWishNum_(reqMsg *base.Empty) (*JustNumberResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcToExchangeWishNum", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*JustNumberResp), e
}
func (self *WishClient2Hall) RpcAreaPostage(reqMsg *base.Empty) *PostageResp {
	msg, e := self.Sender.CallRpcMethod("RpcAreaPostage", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PostageResp)
}

func (self *WishClient2Hall) RpcAreaPostage_(reqMsg *base.Empty) (*PostageResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcAreaPostage", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PostageResp), e
}
func (self *WishClient2Hall) RpcRecycleRatio(reqMsg *base.Empty) *JustNumberResp {
	msg, e := self.Sender.CallRpcMethod("RpcRecycleRatio", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*JustNumberResp)
}

func (self *WishClient2Hall) RpcRecycleRatio_(reqMsg *base.Empty) (*JustNumberResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRecycleRatio", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*JustNumberResp), e
}
func (self *WishClient2Hall) RpcGetConfig(reqMsg *base.Empty) *ConfigResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetConfig", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ConfigResp)
}

func (self *WishClient2Hall) RpcGetConfig_(reqMsg *base.Empty) (*ConfigResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetConfig", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ConfigResp), e
}
func (self *WishClient2Hall) RpcGetUserIdBankCards(reqMsg *base.Empty) *BankCardResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetUserIdBankCards", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BankCardResp)
}

func (self *WishClient2Hall) RpcGetUserIdBankCards_(reqMsg *base.Empty) (*BankCardResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetUserIdBankCards", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BankCardResp), e
}
func (self *WishClient2Hall) RpcRecycleDesc(reqMsg *base.Empty) *DefaultResp {
	msg, e := self.Sender.CallRpcMethod("RpcRecycleDesc", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DefaultResp)
}

func (self *WishClient2Hall) RpcRecycleDesc_(reqMsg *base.Empty) (*DefaultResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRecycleDesc", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DefaultResp), e
}
func (self *WishClient2Hall) RpcSetBoxPreSale(reqMsg *PresaleReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSetBoxPreSale", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *WishClient2Hall) RpcSetBoxPreSale_(reqMsg *PresaleReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSetBoxPreSale", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *WishClient2Hall) RpcSoldOutBox(reqMsg *DealBoxReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSoldOutBox", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *WishClient2Hall) RpcSoldOutBox_(reqMsg *DealBoxReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSoldOutBox", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *WishClient2Hall) RpcReplenishBox(reqMsg *DealBoxReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcReplenishBox", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *WishClient2Hall) RpcReplenishBox_(reqMsg *DealBoxReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReplenishBox", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *WishClient2Hall) RpcSetBoxExpress(reqMsg *SetExpressInfoReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcSetBoxExpress", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *WishClient2Hall) RpcSetBoxExpress_(reqMsg *SetExpressInfoReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSetBoxExpress", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *WishClient2Hall) RpcGetBoxExpress(reqMsg *GetExpressInfoReq) *GetExpressInfoResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetBoxExpress", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GetExpressInfoResp)
}

func (self *WishClient2Hall) RpcGetBoxExpress_(reqMsg *GetExpressInfoReq) (*GetExpressInfoResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetBoxExpress", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GetExpressInfoResp), e
}
func (self *WishClient2Hall) RpcCoinToDiamond(reqMsg *CoinToDiamondReq) *CoinToDiamondResq {
	msg, e := self.Sender.CallRpcMethod("RpcCoinToDiamond", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*CoinToDiamondResq)
}

func (self *WishClient2Hall) RpcCoinToDiamond_(reqMsg *CoinToDiamondReq) (*CoinToDiamondResq, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcCoinToDiamond", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*CoinToDiamondResq), e
}
func (self *WishClient2Hall) RpcDiamondRechargeList(reqMsg *base.Empty) *DiamondRechargeResp {
	msg, e := self.Sender.CallRpcMethod("RpcDiamondRechargeList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DiamondRechargeResp)
}

func (self *WishClient2Hall) RpcDiamondRechargeList_(reqMsg *base.Empty) (*DiamondRechargeResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDiamondRechargeList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DiamondRechargeResp), e
}
func (self *WishClient2Hall) RpcDiamondChangeLogList(reqMsg *DiamondChangeLogReq) *DiamondChangeLogResp {
	msg, e := self.Sender.CallRpcMethod("RpcDiamondChangeLogList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*DiamondChangeLogResp)
}

func (self *WishClient2Hall) RpcDiamondChangeLogList_(reqMsg *DiamondChangeLogReq) (*DiamondChangeLogResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcDiamondChangeLogList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*DiamondChangeLogResp), e
}
func (self *WishClient2Hall) RpcGetPriceSection(reqMsg *base.Empty) *PriceSectionResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetPriceSection", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PriceSectionResp)
}

func (self *WishClient2Hall) RpcGetPriceSection_(reqMsg *base.Empty) (*PriceSectionResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetPriceSection", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PriceSectionResp), e
}
func (self *WishClient2Hall) RpcBatchDare(reqMsg *BatchDareReq) *BatchDareResp {
	msg, e := self.Sender.CallRpcMethod("RpcBatchDare", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*BatchDareResp)
}

func (self *WishClient2Hall) RpcBatchDare_(reqMsg *BatchDareReq) (*BatchDareResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcBatchDare", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*BatchDareResp), e
}
func (self *WishClient2Hall) RpcPlayCfg(reqMsg *base.Empty) *PlayCfgResp {
	msg, e := self.Sender.CallRpcMethod("RpcPlayCfg", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*PlayCfgResp)
}

func (self *WishClient2Hall) RpcPlayCfg_(reqMsg *base.Empty) (*PlayCfgResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcPlayCfg", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*PlayCfgResp), e
}
func (self *WishClient2Hall) RpcSumNum(reqMsg *SumReq) *SumNumResp {
	msg, e := self.Sender.CallRpcMethod("RpcSumNum", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SumNumResp)
}

func (self *WishClient2Hall) RpcSumNum_(reqMsg *SumReq) (*SumNumResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSumNum", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SumNumResp), e
}
func (self *WishClient2Hall) RpcSumMoney(reqMsg *SumMoneyReq) *SumMoneyResp {
	msg, e := self.Sender.CallRpcMethod("RpcSumMoney", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*SumMoneyResp)
}

func (self *WishClient2Hall) RpcSumMoney_(reqMsg *SumMoneyReq) (*SumMoneyResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcSumMoney", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*SumMoneyResp), e
}
func (self *WishClient2Hall) RpcGive(reqMsg *GiveReq) *GiveResp {
	msg, e := self.Sender.CallRpcMethod("RpcGive", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*GiveResp)
}

func (self *WishClient2Hall) RpcGive_(reqMsg *GiveReq) (*GiveResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGive", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*GiveResp), e
}
func (self *WishClient2Hall) RpcActPoolList(reqMsg *base.Empty) *ActPoolResp {
	msg, e := self.Sender.CallRpcMethod("RpcActPoolList", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ActPoolResp)
}

func (self *WishClient2Hall) RpcActPoolList_(reqMsg *base.Empty) (*ActPoolResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcActPoolList", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ActPoolResp), e
}
func (self *WishClient2Hall) RpcActPoolRule(reqMsg *ActPoolRuleReq) *ActPoolRuleResp {
	msg, e := self.Sender.CallRpcMethod("RpcActPoolRule", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ActPoolRuleResp)
}

func (self *WishClient2Hall) RpcActPoolRule_(reqMsg *ActPoolRuleReq) (*ActPoolRuleResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcActPoolRule", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ActPoolRuleResp), e
}
func (self *WishClient2Hall) RpcActName(reqMsg *ActNameReq) *ActNameResp {
	msg, e := self.Sender.CallRpcMethod("RpcActName", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ActNameResp)
}

func (self *WishClient2Hall) RpcActName_(reqMsg *ActNameReq) (*ActNameResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcActName", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ActNameResp), e
}
func (self *WishClient2Hall) RpcActOpenStatus(reqMsg *base.Empty) *ActOpenStatusResp {
	msg, e := self.Sender.CallRpcMethod("RpcActOpenStatus", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ActOpenStatusResp)
}

func (self *WishClient2Hall) RpcActOpenStatus_(reqMsg *base.Empty) (*ActOpenStatusResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcActOpenStatus", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ActOpenStatusResp), e
}
func (self *WishClient2Hall) RpcRechargeActStatus(reqMsg *base.Empty) *ActOpenStatusResp {
	msg, e := self.Sender.CallRpcMethod("RpcRechargeActStatus", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*ActOpenStatusResp)
}

func (self *WishClient2Hall) RpcRechargeActStatus_(reqMsg *base.Empty) (*ActOpenStatusResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcRechargeActStatus", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*ActOpenStatusResp), e
}
func (self *WishClient2Hall) RpcReportWishLog(reqMsg *TypeReq) *base.Empty {
	msg, e := self.Sender.CallRpcMethod("RpcReportWishLog", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*base.Empty)
}

func (self *WishClient2Hall) RpcReportWishLog_(reqMsg *TypeReq) (*base.Empty, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcReportWishLog", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*base.Empty), e
}
func (self *WishClient2Hall) RpcGetRechargeAct(reqMsg *base.Empty) *RechargeActResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetRechargeAct", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*RechargeActResp)
}

func (self *WishClient2Hall) RpcGetRechargeAct_(reqMsg *base.Empty) (*RechargeActResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetRechargeAct", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*RechargeActResp), e
}
func (self *WishClient2Hall) RpcGetRechargeLogs(reqMsg *DataPageReq) *RechargeLogsResp {
	msg, e := self.Sender.CallRpcMethod("RpcGetRechargeLogs", reqMsg)
	easygo.PanicError(e)
	if msg == nil {
		return nil
	}
	return msg.(*RechargeLogsResp)
}

func (self *WishClient2Hall) RpcGetRechargeLogs_(reqMsg *DataPageReq) (*RechargeLogsResp, easygo.IRpcInterrupt) {
	msg, e := self.Sender.CallRpcMethod("RpcGetRechargeLogs", reqMsg)
	if msg == nil {
		return nil, e
	}
	return msg.(*RechargeLogsResp), e
}

// ==========================================================
type IHall2WishClient interface {
}

type Hall2WishClient struct {
	Sender easygo.IMessageSender
}

func (self *Hall2WishClient) Init(sender easygo.IMessageSender) {
	self.Sender = sender
}

//-------------------------------
