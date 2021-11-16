package main

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/client_server"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"log"

	"github.com/astaxie/beego/logs"
)

var _ = fmt.Sprintf
var _ = log.Println

//==========================================================================
//RpcLogin,msg=Account:"5555" Token:"7MsJoDLu98YgC3BE5555" RegistrationId:"" Channel:"" login_type:1 device_type:2 Type:2 PlayerId:1887436001
func (self *Command) RpcLogin(reqMsg *client_hall.LoginMsg) {
	logs.Info("登陆大厅------------")
	resp := hallConnect.FetchEndpoint().RpcLogin(reqMsg)
	logs.Info("resp------>", resp)
	//self.RpcAddSquareCommentZan("0", "0")
	//self.RpcLocationInfoNew()
}
func (self *Command) RpcBindBankCode() {
	msg := &client_hall.BankMessage{
		UserName:   nil,
		IdType:     nil,
		IdNo:       easygo.NewString("44178119941003022X"),
		BankCardNo: easygo.NewString("6214633131067889708"),
		MobileNo:   easygo.NewString("13168180383"),
		ExpireDate: nil,
		Cvv:        nil,
		BankCode:   easygo.NewString("GCB"),
		JPOrderNo:  nil,
		OrderNo:    nil,
		MsgCode:    nil,
		SignNo:     nil,
		Provice:    easygo.NewString("广东省"),
		City:       easygo.NewString("广州市"),
	}
	hallConnect.FetchEndpoint().RpcBindBankCode(msg)
}
func (self *Command) RpcSetPassword() {
	msg := &client_server.PasswordInfo{
		Password: easygo.NewString("111111"),
		Type:     easygo.NewInt32(1),
	}
	hallConnect.FetchEndpoint().RpcSetPassword(msg)
}
func (self *Command) RpcAddBank() {
	msg := &client_hall.BankMessage{
		UserName:   nil,
		IdType:     nil,
		IdNo:       easygo.NewString("44178119941003022X"),
		BankCardNo: easygo.NewString("6214633131067889708"),
		MobileNo:   easygo.NewString("13168180383"),
		ExpireDate: nil,
		Cvv:        nil,
		BankCode:   easygo.NewString("GCB"),
		JPOrderNo:  nil,
		OrderNo:    easygo.NewString("202007151515555"),
		MsgCode:    easygo.NewString("111111"),
		SignNo:     nil,
		Provice:    easygo.NewString("广东省"),
		City:       easygo.NewString("河源市"),
	}
	hallConnect.FetchEndpoint().RpcAddBank(msg)
}

func (self *Command) RpcWithdrawRequest() {
	msg := &client_hall.WithdrawInfo{
		BankCode:    easygo.NewString("GCB"),
		AccountType: easygo.NewString("00"),
		AccountNo:   easygo.NewString("6214633131067889708"),
		AccountName: easygo.NewString("黄家茵"),
		AccountProp: easygo.NewString("0"),
		Amount:      easygo.NewInt64(1.00),
		Result:      nil,
		OrderId:     easygo.NewString("2201159563214689749110"),
		//StartTime:   easygo.NewInt64(),
		//Tax:         nil,
	}

	hallConnect.FetchEndpoint().RpcWithdrawRequest(msg)
}

/*func (self *Command) RpcFlushSquareDynamic() {
	req := &client_hall.NewVersionFlushInfo{
		Type:     easygo.NewInt32(666),
		AdvId:    easygo.NewInt64(777),
		Page:     easygo.NewInt32(888),
		PageSize: easygo.NewInt32(999),
	}
	hallConnect.FetchEndpoint().RpcNewVersionFlushSquareDynamic(req)

}
*/

func (self *Command) RpcNewVersionFlushSquareDynamic() {
	/*	req := &client_square.NewVersionFlushInfo{
			Type:     easygo.NewInt32(1),
			AdvId:    easygo.NewInt64(0),
			Page:     easygo.NewInt32(1),
			PageSize: easygo.NewInt32(5),
			//PlayerId: easygo.NewInt64(1887439538),
		}
		_, err := hallConnect.FetchEndpoint().CallRpcMethod("RpcNewVersionFlushSquareDynamic", req, easygo.SERVER_TYPE_SQUARE)
		easygo.PanicError(err)*/
	//hallConnect.FetchEndpoint().RpcNewVersionFlushSquareDynamic(req)
}

//func (self *Command) RpcAddSquareDynamic() {
//	req := &share_message.DynamicData{
//		LogId:            easygo.NewInt64(878),
//		PlayerId:         easygo.NewInt64(1887439538),
//		Content:          easygo.NewString("测试新版本的api接口"),
//		ClientUniqueCode: easygo.NewString("2"),
//	}
//	hallConnect.FetchEndpoint().RpcAddSquareDynamic(req)
//}

/*func (self *Command) RpcDelSquareDynamic() {
	req := &client_server.RequestInfo{
		Id: easygo.NewInt64(878),
	}
	hallConnect.FetchEndpoint().RpcDelSquareDynamicApi(req)
}
*/

/*func (self *Command) RpcDelNewFriendList() {
	req := &client_hall.DelNewFriendListReq{
		PlayerIds: []int64{1887436829},
	}
	hallConnect.FetchEndpoint().CallRpcMethod("RpcDelNewFriendList", req)
}*/

func (self *Command) RpcGetAllTopic() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcGetAllTopic", nil, common)

}
func (self *Command) RpcAttentionTopic() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}
	req := &client_hall.AttentionTopicReq{
		Id:      []int64{3, 5},
		Operate: easygo.NewInt32(2),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcAttentionTopic", req, common)

}
func (self *Command) RpcMyAttentionTopicList() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}
	req := &client_hall.MyAttentionTopicListReq{
		Page:     easygo.NewInt64(2),
		PageSize: easygo.NewInt64(2),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcMyAttentionTopicList", req, common)

}
func (self *Command) RpcRpcGetTopicDetailReq() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}
	req := &client_hall.TopicDetailReq{
		Id: easygo.NewInt64(1),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcGetTopicDetailReq", req, common)

}
func (self *Command) RpcGetTopicParticipateList() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}
	req := &client_hall.TopicParticipateListReq{
		Id:       easygo.NewInt64(1),
		Page:     easygo.NewInt64(1),
		PageSize: easygo.NewInt64(10),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcGetTopicParticipateList", req, common)

}
func (self *Command) RpcGetTopicMainPageList() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}
	req := &client_hall.TopicMainPageListReq{
		Id:       easygo.NewInt64(0),
		Page:     easygo.NewInt64(1),
		PageSize: easygo.NewInt64(10),
		ReqType:  easygo.NewInt32(1),
		Name:     easygo.NewString("#找cp助脱单#"),
	}
	backMsg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetTopicMainPageList", req, common)
	logs.Info("返回 RpcGetTopicMainPageList 数据:", backMsg)
}

//获取动态详情: 18618
func (self *Command) RpcGetDynamicInfoNew(id string) {

	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}
	req := &client_server.IdInfo{
		Id: easygo.NewInt64(id),
	}
	backMsg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetDynamicInfoNew", req, common)
	logs.Info("返回 RpcGetDynamicInfoNew 数据:", backMsg)
}

//获取动态主评论:GetDynamicMainCommentNew 18618 1 10
func (self *Command) RpcGetDynamicMainCommentNew(id, page, size string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}
	req := &client_server.IdInfo{
		Id:       easygo.NewInt64(id),
		Page:     easygo.NewInt32(page),
		PageSize: easygo.NewInt32(size),
		HotList:  []int64{24080, 24070},
	}
	backMsg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetDynamicMainCommentNew", req, common)
	logs.Info("返回 RpcGetDynamicMainCommentNew 数据:", backMsg)
}

//获取动态子评论:24032
func (self *Command) RpcGetDynamicSecondaryCommentNew(id, mainId string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}
	req := &client_server.IdInfo{
		Id:     easygo.NewInt64(id),
		MainId: easygo.NewInt64(mainId),
	}
	backMsg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetDynamicSecondaryCommentNew", req, common)
	logs.Info("返回 RpcGetDynamicSecondaryCommentNew 数据:", backMsg)
}

//对主评论进行点赞18617,24032    dynamicId 1159  mainCommentId 1004
func (self *Command) RpcAddSquareCommentZan(id, commentId string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}
	req := &share_message.CommentDataZan{
		DynamicId: easygo.NewInt64(1159),
		CommentId: easygo.NewInt64(1004),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcAddSquareCommentZan", req, common)
}

func (self *Command) RpcGetTopicTypeList() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}

	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcGetTopicTypeList", nil, common)

}
func (self *Command) RpcGetTopicList() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}

	req := &client_hall.TopicListReq{
		TopicTypeId: easygo.NewInt64(1),
		Page:        easygo.NewInt64(1),
		PageSize:    easygo.NewInt64(10),
	}

	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcGetTopicList", req, common)

}
func (self *Command) RpcSearchHotTopic() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}

	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcSearchHotTopic", nil, common)

}
func (self *Command) RpcFlushTopic() {

	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}
	req := &client_hall.FlushTopicReq{
		Page:     easygo.NewInt64(2),
		PageSize: easygo.NewInt64(1),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcFlushTopic", req, common)

}
func (self *Command) RpcAttentionRecommendPlayer() {

	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}

	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcAttentionRecommendPlayer", nil, common)

}

func (self *Command) RpcBsOpShopOrder() {

	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}

	req := &server_server.ShopOrderRequest{
		OrderId: easygo.NewInt64(666888),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcBsOpShopOrder", req, common)

}
func (self *Command) RpcFlushSquareDynamicTopic() {

	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}

	req := &share_message.FlushSquareDynamicTopicReq{
		AdvId:    easygo.NewInt64(0),
		Page:     easygo.NewInt32(1),
		PageSize: easygo.NewInt32(10),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcFlushSquareDynamicTopic", req, common)

}
func (self *Command) RpcSquareAttention() {

	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}

	req := &client_hall.SquareAttentionReq{
		Page:               easygo.NewInt64(1),
		PageSize:           easygo.NewInt64(10),
		HasAttentionTopic:  easygo.NewBool(true),
		HasAttentionPlayer: easygo.NewBool(true),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcSquareAttention", req, common)

}
func (self *Command) RpcTopicHotDynamicParticipatePlayer() {

	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}

	req := &client_hall.TopicParticipateListReq{
		Id:       easygo.NewInt64(2),
		Page:     easygo.NewInt64(1),
		PageSize: easygo.NewInt64(10),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcTopicHotDynamicParticipatePlayer", req, common)

}
func (self *Command) RpcTopicHead() {

	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(6),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcTopicHead", nil, common)

}
func (self *Command) RpcNewTeamSetting() {

	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(1),
	}

	req := &client_hall.NewTeamSettingReq{
		TeamId: easygo.NewInt64(18801181),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcNewTeamSetting", req, common)

}

//获取虚拟道具配置
func (self *Command) RpcGetPropsItems() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.6"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}

	backMsg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetPropsItems", nil, common)
	logs.Info("返回 RpcGetPropsItems 数据:", backMsg)
}
func (self *Command) RpcGetCoinRechargeList(way string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.6"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.CoinRechargeList{
		Way: easygo.NewInt32(way),
	}
	backMsg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetCoinRechargeList", msg, common)
	logs.Info("返回 RpcGetCoinRechargeList 数据:", backMsg)
}
func (self *Command) RpcGetCoinShopList(t string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.6"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.CoinShopList{
		Type: easygo.NewInt32(t),
	}
	backMsg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetCoinShopList", msg, common)
	logs.Info("返回 RpcGetCoinShopList 数据:", backMsg)
}
func (self *Command) RpcDeleteNearByInfo() {

	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(1),
	}

	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcDeleteNearByInfo", nil, common)

}
func (self *Command) RpcLocationInfoNew() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(1),
	}
	reqMsg := &client_hall.LocationInfoNewReq{
		X:        easygo.NewFloat64(113.336020),
		Y:        easygo.NewFloat64(23.140620),
		Sex:      easygo.NewInt32(0),
		Page:     easygo.NewInt64(4),
		PageSize: easygo.NewInt64(10),
		Sort:     easygo.NewInt32(for_game.NEAR_SORT_DISTANCE),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcLocationInfoNew", reqMsg, common)

}

//兑换硬币
func (self *Command) RpcCoinRecharge(id string) {
	reqMsg := &client_hall.CoinRechargeResp{
		Id: easygo.NewInt64(id),
	}
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcCoinRecharge", reqMsg, common)
	logs.Info("兑换结果:", msg)
}

//购买虚拟道具
func (self *Command) RpcBuyCoinItem(id, num, way string) {
	reqMsg := &client_hall.BuyCoinItem{
		Id:  easygo.NewInt64(id),
		Num: easygo.NewInt32(num),
		Way: easygo.NewInt32(way),
	}
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcBuyCoinItem", reqMsg, common)
	logs.Info("购买结果:", msg)
}

//使用道具
func (self *Command) RpcUseCoinItem(id, way string) {
	reqMsg := &client_hall.UseCoinItem{
		Id:  easygo.NewInt64(id),
		Way: easygo.NewInt32(way),
	}
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcUseCoinItem", reqMsg, common)
	logs.Info("使用道具结果:", msg)
}

//获取装备情况
func (self *Command) RpcGetPlayerEquipment(id string) {
	reqMsg := &client_hall.EquipmentReq{
		Id: easygo.NewInt64(id),
	}
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetPlayerEquipment", reqMsg, common)
	logs.Info("装备情况:", msg)
}

//获取背包信息
func (self *Command) RpcGetPlayerBagItems() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetPlayerBagItems", nil, common)
	logs.Info("背包道具:", msg)
}

//广播群特效
func (self *Command) RpcBroadCastQTX() {
	send := &client_hall.BroadCastQTX{
		TeamId:   easygo.NewInt64(18800946),
		PlayerId: easygo.NewInt64(1887555881),
		PropsId:  easygo.NewInt64(50001),
	}
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	_, _ = hallConnect.FetchEndpoint().CallRpcMethod("RpcBroadCastQTX", send, common)
}

//请求玩家会话数据
func (self *Command) RpcGetSessionData() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetSessionData", nil, common)
	logs.Info("RpcGetSessionData:", msg)
}

//请求其他数据
func (self *Command) RpcGetPlayerOtherDataNew() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetPlayerOtherDataNew", nil, common)
	logs.Info("RpcGetPlayerOtherDataNew:", msg)
}
func (self *Command) RpcGetPlayerFriends() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetPlayerFriends", nil, common)
	logs.Info("RpcGetPlayerFriends:", msg)
}

func (self *Command) RpcGetSessionChat(id, start, end string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.SessionChatData{
		SessionId: easygo.NewString(id),
		StartId:   easygo.NewInt64(start),
		EndId:     easygo.NewInt64(end),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetSessionChat", msg, common)
	logs.Info("RpcGetSessionChat:", resp)
}
func (self *Command) RpcChatNew(id, s string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &share_message.Chat{
		SourceId:    easygo.NewInt64(1887436624),
		SessionId:   easygo.NewString(id),
		TargetId:    easygo.NewInt64(id),
		Content:     easygo.NewString(s),
		ChatType:    easygo.NewInt32(2),
		ContentType: easygo.NewInt32(1),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcChatNew", msg, common)
	logs.Info("RpcChatNew:", resp)
}
func (self *Command) RpcCheckIsTeamMember() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.CheckTeamMember{
		PlayerId: easygo.NewInt64(1887438007),
		TeamId:   easygo.NewInt64(18800051),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcCheckIsTeamMember", msg, common)
	logs.Info("返回:", resp)
}

func (self *Command) RpcZanVoiceCard(playerId string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.PlayerInfoReq{
		PlayerId: easygo.NewInt64(playerId),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcZanVoiceCard", msg, common)
	logs.Info("返回:", resp)
}
func (self *Command) RpcSayHiToPlayer(playerId string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.PlayerInfoReq{
		PlayerId: easygo.NewInt64(playerId),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcSayHiToPlayer", msg, common)
	logs.Info("返回:", resp)
}
func (self *Command) RpcGetLoveMeNewNum() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}

	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetLoveMeNewNum", nil, common)
	logs.Info("返回:", resp)
}
func (self *Command) RpcReadLoveMeLog() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}

	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcReadLoveMeLog", nil, common)
	logs.Info("返回:", resp)
}
func (self *Command) RpcGetLoveMeList(page string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.LoveMeReq{
		Page: easygo.NewInt32(page),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetLoveMeList", msg, common)

	data := resp.(*client_hall.LoveMeResp)
	for _, v := range data.GetCards() {
		logs.Info("PlayerId:%v NickName:%v HeadUrl:%v Sex:%v Constellation:%v MatchingDegree:%v ZanNum:%v VoiceUrl:%v IsOnLine:%v BgUrl:%v Content:%v AttentionTime:%v AttentionType:%v ",
			v.GetPlayerId(), v.GetNickName(), v.GetHeadUrl(), v.GetSex(), v.GetConstellation(), v.GetMatchingDegree(), v.GetZanNum(), v.GetVoiceUrl(), v.GetIsOnLine(), v.GetBgUrl(), v.GetContent(), v.GetAttentionTime(), v.GetAttentionType())
	}
}
func (self *Command) RpcGetMyLoveList(page string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.LoveMeReq{
		Page: easygo.NewInt32(page),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetMyLoveList", msg, common)
	data := resp.(*client_hall.MyLoveResp)
	for _, v := range data.GetCards() {
		logs.Info("PlayerId:%v NickName:%v HeadUrl:%v Sex:%v Constellation:%v MatchingDegree:%v ZanNum:%v VoiceUrl:%v IsOnLine:%v BgUrl:%v Content:%v AttentionTime:%v AttentionType:%v ",
			v.GetPlayerId(), v.GetNickName(), v.GetHeadUrl(), v.GetSex(), v.GetConstellation(), v.GetMatchingDegree(), v.GetZanNum(), v.GetVoiceUrl(), v.GetIsOnLine(), v.GetBgUrl(), v.GetContent(), v.GetAttentionTime(), v.GetAttentionType())
	}
}
func (self *Command) RpcGetVoiceTags(t string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.SearchVoiceVideoReq{
		Type: easygo.NewInt32(t),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetVoiceTags", msg, common)
	logs.Info("RpcGetVoiceTags返回:", resp)

}

//制作素材
func (self *Command) RpcMakeVoiceVideo() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.VoiceVideo{
		Maker:    easygo.NewString("无名"),
		Name:     easygo.NewString("抱一抱"),
		TagIds:   []int32{11},
		Content:  easygo.NewString("有时候觉得如约而至是多美好的词，等得很苦，却从不辜负"),
		MusicUrl: easygo.NewString("https://cdn.happyyuyin.com/overdub/14009.mp3"),
		ImageUrl: easygo.NewString("https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/jkl17.jpeg"),
		Type:     easygo.NewInt32(3),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcMakeVoiceVideo", msg, common)
	logs.Info("RpcMakeVoiceVideo返回:", resp)
}

//制作卡片
func (self *Command) RpcMixVoiceVideo(id string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.MixVoiceVideo{
		BgId:       easygo.NewInt64(id),
		MyVoiceUrl: easygo.NewString("https://cdn.happyyuyin.com/overdub/14009.mp3"),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcMixVoiceVideo", msg, common)
	logs.Info("RpcMixVoiceVideo返回:", resp)
}

func (self *Command) RpcGetHotEpisode(reqType string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.SearchVoiceVideoReq{
		Type: easygo.NewInt32(reqType),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetHotEpisode", msg, common)
	logs.Info("RpcGetHotEpisode返回:", resp)
}

func (self *Command) RpcGetMayLikeEpisode(reqType string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.SearchVoiceVideoReq{
		Type: easygo.NewInt32(reqType),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetMayLikeEpisode", msg, common)
	logs.Info("RpcGetHotEpisode返回:", resp)
}

func (self *Command) RpcGetVoiceProduct(id, page string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.VoiceProduct{
		TabId: easygo.NewInt32(id),
		Page:  easygo.NewInt32(page),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetVoiceProduct", msg, common)
	logs.Info("RpcGetHotEpisode返回:", resp)
}

//获取名片列表
func (self *Command) RpcGetVoiceCards() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetVoiceCards", nil, common)
	logs.Info("RpcGetVoiceCards返回:", resp)
}

func (self *Command) RpcCoinRechargeAct(level, giveType, pwd, payWay string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(for_game.SERVER_TYPE_HALL),
	}
	msg := &client_hall.RechargeActReq{
		ActCfgId: easygo.NewInt64(level),
		GiveType: easygo.NewInt32(giveType),
		PassWord: easygo.NewString(pwd),
		PayWay:   easygo.NewInt32(payWay),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcCoinRechargeAct", msg, common)
	logs.Info("RpcCoinRechargeAct:", resp)
}

//充值接口调用
func (self *Command) RpcRechargeMoney() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &share_message.PayOrderInfo{
		ProduceName: easygo.NewString("充值"),
		Amount:      easygo.NewString("1"),
		PayWay:      easygo.NewInt32(1),
		PlayerId:    easygo.NewInt64(playerId),
		PayType:     easygo.NewInt32(1),
		PayId:       easygo.NewInt32(12),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcRechargeMoney", msg, common)
	logs.Info("RpcCoinRechargeAct:", resp)
}

//获取支持的支付渠道
func (self *Command) RpcGetSupportPayChannel(t string) {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(2),
	}
	msg := &client_hall.PayChannels{
		Type: easygo.NewInt32(t),
	}

	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetSupportPayChannel", msg, common)
	logs.Info("RpcGetSupportPayChannel:", resp)
}

//请求菜单id
func (self *Command) RpcGetAllMainMenu() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(for_game.SERVER_TYPE_HALL),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetAllMainMenu", nil, common)
	logs.Info("RpcGetAllMainMenu:", resp)
}

//请求悬浮框
func (self *Command) RpcGetAllTipAdvs() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(for_game.SERVER_TYPE_HALL),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcGetAllTipAdvs", nil, common)
	logs.Info("RpcGetAllTipAdvs:", resp)
}

// 话题贡献榜
func (self *Command) RpcTopicDevoteList() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(for_game.SERVER_TYPE_SQUARE),
	}
	msg := &client_hall.TopicDevoteListReq{
		DataType:  easygo.NewInt64(3),
		Page:      easygo.NewInt64(1),
		PageSize:  easygo.NewInt64(10),
		TopicName: easygo.NewString("#找cp助脱单#"),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcTopicDevoteList", msg, common)
	logs.Info("RpcTopicDevoteList:", resp)
}

// 获取申请话题主条件
func (self *Command) RpcTopicMasterCondition() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(for_game.SERVER_TYPE_SQUARE),
	}
	msg := &client_hall.TopicMasterConditionReq{
		TopicId: easygo.NewInt64(1),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcTopicMasterCondition", msg, common)
	logs.Info("RpcTopicMasterCondition:", resp)
}

// 申请话题主
func (self *Command) RpcApplyTopicMaster() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(for_game.SERVER_TYPE_SQUARE),
	}
	msg := &client_hall.ApplyTopicMasterReq{
		TopicId:        easygo.NewInt64(4),
		IsManageExp:    easygo.NewBool(true),
		Reason:         easygo.NewString("我就是晚安话题主"),
		ContactDetails: easygo.NewString("1380438943"),
		TopicName:      easygo.NewString("#晚安#"),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcApplyTopicMaster", msg, common)
	logs.Info("RpcApplyTopicMaster:", resp)
}

// 话题主修改话题信息
func (self *Command) RpcTopicMasterEdit() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(for_game.SERVER_TYPE_SQUARE),
	}
	msg := &client_hall.TopicMasterEditReq{
		TopicId:     easygo.NewInt64(1),
		HeadURL:     easygo.NewString("话题头像url"),
		Description: easygo.NewString("这是话题介绍"),
		BgUrl:       easygo.NewString("话题背景url"),
		TopicRule:   easygo.NewString("话题规则"),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcTopicMasterEdit", msg, common)
	logs.Info("RpcTopicMasterEdit:", resp)
}

// 话题置顶
func (self *Command) RpcTopicTop() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(for_game.SERVER_TYPE_SQUARE),
	}
	msg := &client_hall.TopicTopReq{
		TopicId: easygo.NewInt64(2),
		LogId:   easygo.NewInt64(29671),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcTopicTop", msg, common)
	logs.Info("RpcTopicTop:", resp)
}

// 取消话题置顶
func (self *Command) RpcTopicTopCancel() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(for_game.SERVER_TYPE_SQUARE),
	}
	msg := &client_hall.TopicTopCancelReq{
		TopicId: easygo.NewInt64(2),
		LogId:   easygo.NewInt64(29671),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcTopicTopCancel", msg, common)
	logs.Info("RpcTopicTopCancel:", resp)
}

// 话题排行榜规则说明
func (self *Command) RpcTopicLeaderBoardDescription() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(for_game.SERVER_TYPE_SQUARE),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcTopicLeaderBoardDescription", nil, common)
	logs.Info("RpcTopicLeaderBoardDescription:", resp)
}

// 退出话题主
func (self *Command) RpcQuitTopicMaster() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(for_game.SERVER_TYPE_SQUARE),
	}
	msg := &client_hall.QuitTopicMasterReq{
		TopicId: easygo.NewInt64(1),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcQuitTopicMaster", msg, common)
	logs.Info("RpcQuitTopicMaster:", resp)
}

// 话题主删除话题中的动态
func (self *Command) RpcTopicMasterDelDynamic() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(for_game.SERVER_TYPE_SQUARE),
	}
	msg := &client_hall.TopicMasterDelDynamicReq{
		TopicId:      easygo.NewInt64(39),
		LogId:        easygo.NewInt64(29682),
		DelReasonId:  easygo.NewInt32(1),
		DelReasonMsg: easygo.NewString("删除一个动态"),
		TopicName:    easygo.NewString("#找cp助脱单#"),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcTopicMasterDelDynamic", msg, common)
	logs.Info("RpcTopicMasterDelDynamic:", resp)
}

func (self *Command) RpcTopicTeamDynamicList() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(for_game.SERVER_TYPE_HALL),
	}
	msg := &client_hall.TopicTeamDynamicReq{}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcTopicTeamDynamicList", msg, common)
	data := resp.(*client_hall.TopicTeamDynamicResp)
	logs.Info("TopicTeamDynamicResp==================>", data)
}

//创建话题群组
func (self *Command) RpcCreateTopicTeam() {
	common := &base.Common{
		Version:    easygo.NewString("2.7.5"),
		UserId:     nil,
		Token:      nil,
		Flag:       nil,
		ServerType: easygo.NewInt32(for_game.SERVER_TYPE_HALL),
	}
	msg := &client_hall.CreateTeam{
		Topic:     easygo.NewString("#晚安#"),
		TopicDesc: easygo.NewString("大家晚安!"),
		TeamName:  easygo.NewString("晚安"),
	}
	resp, _ := hallConnect.FetchEndpoint().CallRpcMethod("RpcCreateTopicTeam", msg, common)
	data := resp.(*client_hall.TopicTeamDynamicResp)
	logs.Info("RpcCreateTopicTeam==================>", data)
}
