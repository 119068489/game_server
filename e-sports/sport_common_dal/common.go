package sport_common_dal

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/protobuf/proto"
	"github.com/astaxie/beego/logs"
)

func NextId(tbn string) int64 {
	id := for_game.NextId(tbn)
	return id
}
func GetC(table string) (*mgo.Collection, func()) {
	return easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, table)
}

//发送给指定大厅指定玩家发送消息
func SendMsgToHallClientNewEx(pServerInfoMgr *for_game.ServerInfoManager, playerId int64, sid int32, methodName string, msg easygo.IMessage) bool {
	var srv *share_message.ServerInfo
	if pServerInfoMgr == nil {
		logs.Info("服务器列表不能为空 ")
		return false
	}
	if sid == 0 {
		srv = pServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_HALL)
	} else {
		srv = pServerInfoMgr.GetServerInfo(sid)
	}
	if srv == nil {
		logs.Info(" 找不到大厅服务器 SID：%d ", sid)
		return false
	}
	var msgByte []byte
	if msg != nil {
		b, err := msg.Marshal()
		easygo.PanicError(err)
		msgByte = b
	} else {
		msgByte = []byte{}
	}

	req := &share_message.MsgToClient{
		PlayerIds: []int64{playerId},
		RpcName:   easygo.NewString(methodName),
		MsgName:   easygo.NewString(proto.MessageName(msg)),
		Msg:       msgByte,
	}
	_, err := for_game.SendToServerEx(srv, "RpcMsgToHallClient", req)

	if err != nil {
		return false
	}

	return true
}

//发送订单消息给指定用户
func sendGameOrderMsgToPlayer(pServerInfoMgr *for_game.ServerInfoManager, playerId int64, data *share_message.TableESPortsGameOrderSysMsg) bool {

	isSend := false
	pbase := for_game.GetRedisPlayerBase(playerId)
	var unPayed []*share_message.TableESPortsGameOrderSysMsg
	var payed []*share_message.TableESPortsGameOrderSysMsg
	if pbase != nil && pbase.GetIsOnLine() && pbase.GetSid() > 0 {

		//在线的处理
		if data.GetBetResult() == for_game.GAME_GUESS_BET_RESULT_1 { //已结算的 就需要删除掉以前的数据
			unPayed = []*share_message.TableESPortsGameOrderSysMsg{data}
		} else {
			payed = []*share_message.TableESPortsGameOrderSysMsg{data}
		}
		rd := &client_hall.ESPortsSysMsgList{
			PlayerId: easygo.NewInt64(playerId),
			UnPayed:  unPayed,
			Payed:    payed,
		}
		logs.Info("发送推送系统竞猜 RpcESportNewSysMessage 2", rd, "playerId", playerId, "sid:", pbase.GetSid())
		isSend = SendMsgToHallClientNewEx(pServerInfoMgr, playerId, pbase.GetSid(), "RpcESportNewSysMessage", rd)
		/*if isSend {
			if data.GetBetResult() != for_game.GAME_GUESS_BET_RESULT_1 { //已结算的 就需要删除掉以前的数据
				DeleteTableESPortsGameOrderEndSysMsgByOrderId(playerId, data.GetOrderId())
			}
		}*/
	}
	return isSend
}

//发送多個訂單订单消息给指定用户
func sendGameMultipleOrderMsgToPlayer(pServerInfoMgr *for_game.ServerInfoManager, playerId int64, data []*share_message.TableESPortsGameOrderSysMsg) bool {

	isSend := false
	pbase := for_game.GetRedisPlayerBase(playerId)
	unPayed := []*share_message.TableESPortsGameOrderSysMsg{}
	payed := []*share_message.TableESPortsGameOrderSysMsg{}
	if pbase != nil && pbase.GetIsOnLine() && pbase.GetSid() > 0 {

		for _, v := range data {
			if v.GetBetResult() == for_game.GAME_GUESS_BET_RESULT_1 { //已结算的 就需要删除掉以前的数据
				unPayed = append(unPayed, v)
			} else {
				payed = append(payed, v)
			}
		}

		rd := &client_hall.ESPortsSysMsgList{
			PlayerId: easygo.NewInt64(playerId),
			UnPayed:  unPayed,
			Payed:    payed,
		}
		logs.Info("发送推送系统多個竞猜訂單 RpcESportNewSysMessage 2", rd, "playerId", playerId, "sid:", pbase.GetSid())
		isSend = SendMsgToHallClientNewEx(pServerInfoMgr, playerId, pbase.GetSid(), "RpcESportNewSysMessage", rd)
		/*if isSend {
			if data.GetBetResult() != for_game.GAME_GUESS_BET_RESULT_1 { //已结算的 就需要删除掉以前的数据
				DeleteTableESPortsGameOrderEndSysMsgByOrderId(playerId, data.GetOrderId())
			}
		}*/
	}
	return isSend
}

//推送多个比赛竞猜数据 ,  betResult = 订单状态
func PushGameMultipleOrderSysMsg(pServerInfoMgr *for_game.ServerInfoManager, playerId int64, list []*share_message.TableESPortsGameOrderSysMsg) *client_hall.ESportCommonResult {

	rd := &client_hall.ESportCommonResult{}

	if playerId < 1 {
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString("用户Id不能为0")
		return rd
	}

	if pServerInfoMgr == nil {
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString("服务器列表不能为空")
		return rd
	}

	isSend := sendGameMultipleOrderMsgToPlayer(pServerInfoMgr, playerId, list)
	if !isSend { //如果没有发送出去,就存库
		//if betResult == for_game.GAME_GUESS_BET_RESULT_1 {
		//	func(ls1 []*share_message.TableESPortsGameOrderSysMsg) {
		//		col, closeFun := GetC(for_game.TABLE_ESPORTS_GAME_ORDER_SYS_MSG)
		//		defer closeFun()
		//		for _, v := range ls1 {
		//			CreateTableESPortsGameOrderSysMsgEx(col, v)
		//		}
		//	}(list)
		//
		//} else {
		//	func(ls2 []*share_message.TableESPortsGameOrderSysMsg) {
		//		col, closeFun := GetC(for_game.TABLE_ESPORTS_GAME_ORDER_SYS_MSG_E)
		//		defer closeFun()
		//		for _, v := range ls2 {
		//			CreateTableESPortsGameEndOrderSysMsgEx(col, v)
		//		}
		//	}(list)
		//}

		unPayed := make([]*share_message.TableESPortsGameOrderSysMsg, 0)
		payed := make([]*share_message.TableESPortsGameOrderSysMsg, 0)

		for _, v := range list {
			v.CreateTime = easygo.NewInt64(easygo.NowTimestamp())
			v.UpdateTime = easygo.NewInt64(easygo.NowTimestamp())
			if v.GetBetResult() == for_game.GAME_GUESS_BET_RESULT_1 { //已结算的 就需要删除掉以前的数据
				unPayed = append(unPayed, v)
			} else {
				payed = append(payed, v)
			}
		}

		if unPayed != nil && len(unPayed) > 0 {
			CreateTableESPortsGameOrderSysMsgEx(unPayed)
		}

		if payed != nil && len(payed) > 0 {
			CreateTableESPortsGameEndOrderSysMsgEx(payed)
		}

	}
	rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
	rd.Msg = easygo.NewString("操作结束")
	return rd
}

//推送比赛竞猜数据
func PushGameOrderSysMsg(pServerInfoMgr *for_game.ServerInfoManager, data *share_message.TableESPortsGameOrderSysMsg) *client_hall.ESportCommonResult {

	rd := &client_hall.ESportCommonResult{}
	playerId := data.GetPlayerId()
	if playerId < 1 {
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString("用户Id不能为0")
		return rd
	}

	if pServerInfoMgr == nil {
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString("服务器列表不能为空")
		return rd
	}

	isSend := sendGameOrderMsgToPlayer(pServerInfoMgr, playerId, data)
	//if data.GetBetResult() == for_game.GAME_GUESS_BET_RESULT_1 { // 待开奖的话，
	//	CreateTableESPortsGameEndOrderSysMsg(data)
	//}
	if !isSend { //如果没有发送出去,就存库
		if data.GetBetResult() == for_game.GAME_GUESS_BET_RESULT_1 {
			CreateTableESPortsGameOrderSysMsg(data)
		} else {
			CreateTableESPortsGameEndOrderSysMsg(data)
			//	UpdateTableESPortsGameOrderEndSysMsg(data.GetOrderId(), data.GetBetResult(), data.GetResultAmount())
		}
	}
	rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
	rd.Msg = easygo.NewString("操作结束")
	return rd
}

//随机选择一台发送
func SendMsgToIdOtherServer(pServerInfoMgr *for_game.ServerInfoManager, t int32, methodName string, msg easygo.IMessage) (easygo.IMessage, *base.Fail) {
	if pServerInfoMgr == nil {
		return nil, easygo.NewFailMsg("服务器列表不能为空")
	}

	srv := pServerInfoMgr.GetIdelServer(t)
	if srv != nil {
		return SendMsgToServerNew(pServerInfoMgr, srv.GetSid(), methodName, msg)
	}
	s := fmt.Sprintf("找不到类型type为%v的服务器", t)
	return nil, easygo.NewFailMsg(s)
}

//指定某个用户在特定的大厅发送
func SendMsgToServerNewEx(pServerInfoMgr *for_game.ServerInfoManager, pid int64, methodName string, msg easygo.IMessage) (easygo.IMessage, *base.Fail) {

	player := for_game.GetRedisPlayerBase(pid)
	if player == nil {
		return nil, easygo.NewFailMsg("无效的玩家id")
	}
	return SendMsgToServerNew(pServerInfoMgr, player.GetSid(), methodName, msg)
}

//服务器间通讯通用
func SendMsgToServerNew(pServerInfoMgr *for_game.ServerInfoManager, sid int32, methodName string, msg easygo.IMessage, pid ...int64) (easygo.IMessage, *base.Fail) {
	var srv *share_message.ServerInfo
	if pServerInfoMgr == nil {
		return nil, easygo.NewFailMsg("服务器列表不能为空")
	}
	if sid == 0 {
		srv = pServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_HALL)
	} else {
		srv = pServerInfoMgr.GetServerInfo(sid)
	}

	if srv == nil {
		return nil, easygo.NewFailMsg("无效的服务器id =" + easygo.AnytoA(sid))
	}
	var msgByte []byte
	if msg != nil {
		b, err := msg.Marshal()
		easygo.PanicError(err)
		msgByte = b
	} else {
		msgByte = []byte{}
	}

	playerId := append(pid, 0)[0]
	req := &share_message.MsgToServer{
		PlayerId: easygo.NewInt64(playerId),
		RpcName:  easygo.NewString(methodName),
		MsgName:  easygo.NewString(proto.MessageName(msg)),
		Msg:      msgByte,
	}
	return for_game.SendToServerEx(srv, "RpcMsgToOtherServer", req)
}
