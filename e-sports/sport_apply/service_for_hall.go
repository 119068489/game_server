package sport_apply

import (
	dal "game_server/e-sports/sport_common_dal"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
	"reflect"
	"time"
)

//===================================================================

type ServiceForHall struct {
	Service reflect.Value
}

func GetOrCraeteESPortsPlayer(plyId int64) *for_game.RedisESportPlayerObj {
	pinfo := for_game.GetRedisESportPlayerObj(plyId)
	if pinfo == nil {
		pdata := &share_message.TableESPortsPlayer{
			Id:                easygo.NewInt64(plyId),
			Status:            easygo.NewInt32(for_game.ESPORT_PLAYER_STATUS_1),
			LastPullTime:      easygo.NewInt64(0),
			CurrentRoomLiveId: easygo.NewInt64(0),
			CreateTime:        easygo.NewInt64(easygo.NowTimestamp()),
			LastLoginTime:     easygo.NewInt64(easygo.NowTimestamp()),
		}
		code, _ := dal.CreateTableESPortsPlayer(pdata)

		if code == for_game.C_OPT_SUCCESS {
			pinfo = for_game.NewRedisESportPlayerObj(plyId, pdata)
			logs.Info("创建电竞用户: %d", plyId)
		}
	}
	return pinfo
}

// 玩家上线通知
func (self *ServiceForHall) RpcESPortsPlayerOnLine(common *base.Common, reqMsg *share_message.PlayerState) easygo.IMessage {
	logs.Info("RpcESPortsPlayerOnLine ", reqMsg)
	///PlayerOnlineMgr.PlayerOnline(reqMsg.GetPlayerId(), reqMsg.GetServerId())
	//plyId := reqMsg.GetPlayerId()

	/*if pinfo != nil {
		list, code := dal.GetTableESportAllSysList(pinfo.GetLastPullTime())
		logs.Info("GetTableESportAllSysList", list)
		if code == for_game.C_OPT_SUCCESS {

			if list != nil || len(list) > 0 {
				rd := &client_hall.ESPortsSysMsgList{
					PlayerId:   easygo.NewInt64(plyId),
					DataType:   easygo.NewInt32(1),
					SysMsgList: list,
				}

				//go func() {
				//time.Sleep(time.Second * 2)

				logs.Info("发送系统消息 RpcESportNewSysMessage", rd, "plyId", plyId, "sid:", reqMsg.GetServerId())
				SendMsgToHallClientNewEx(plyId, reqMsg.GetServerId(), "RpcESportNewSysMessage", rd)
				//测试期间每次都推
				pinfo.SetLastPullTime(easygo.NowTimestamp())
				//}()

			}

		} else {
			logs.Info("获取系统消息数据失败 用户ID：", plyId)
		}
	} else {
		logs.Info("NewRedisESportPlayerObj", "nil")
	}*/
	return nil
}

func handleOffLine(plyId int64) {
	logs.Info("5秒后处理用户：%d离线 ", plyId)
	time.Sleep(time.Second * 5) //5秒后处理离线 未免网络波动
	obj := for_game.GetRedisPlayerBase(plyId)
	b := false
	if obj != nil {
		b = obj.GetIsOnLine()
	}
	if !b {
		//如果已离线
		pinfo := for_game.GetRedisESportPlayerObj(plyId)
		if pinfo != nil {
			liveId := pinfo.GetCurrentRoomLiveId()
			if liveId > 0 {
				room := for_game.GetRedisLiveRoomPlayerObj(liveId)
				room.LeaveRoom(plyId)
				pinfo.SetCurrentRoomLiveId(0) //离线自动离开
				logs.Info("5秒后处理用户：%d离线 设置离开放映厅%d ", plyId, liveId)
			}
			pinfo.SaveToMongoEx()
			//for_game.SaveRedisESportPlayerToMongo()
		}

		dpsd := for_game.GetRedisESportBpsDurationLogObj(plyId)
		if dpsd != nil {
			dpsd.IsDeleteRedis = true
			dpsd.EndCurrentBpsDuration(for_game.ESPORT_BPS_PAGE_TYPE_1) //, 0, 0, 0, 0, 0) //結算電競所有停留時長
			dpsd.UpdateData()

		}
	}
}

// 玩家离线通知
func (self *ServiceForHall) RpcESPortsPlayerOffLine(common *base.Common, reqMsg *share_message.PlayerState) easygo.IMessage {
	logs.Info("RpcESPortsPlayerOffLine ", reqMsg)
	plyId := reqMsg.GetPlayerId()
	easygo.Spawn(handleOffLine, plyId) //离线处理
	//PlayerOnlineMgr.PlayerOffline(reqMsg.GetPlayerId())
	return nil
}

//推送比赛竞猜数据
func (self *ServiceForHall) RpcESPortsPushGameOrderSysMsg(common *base.Common, reqMsg *share_message.TableESPortsGameOrderSysMsg) easygo.IMessage {
	logs.Info("RpcESPortsPushGameOrderSysMsg reqMsg", reqMsg)
	rd := dal.PushGameOrderSysMsg(PServerInfoMgr, reqMsg)
	logs.Info("RpcESPortsPushGameOrderSysMsg rd", rd)
	return rd
}

func (self *ServiceForHall) RpcESportDataStatusInfo(common *base.Common, reqMsg *client_hall.ESportDataStatusInfo) *client_hall.ESportCommonResult {
	logs.Info("RpcESportDataStatusInfo reqMsg", reqMsg)
	rd := &client_hall.ESportCommonResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
	}
	if reqMsg.GetMenuId() == for_game.ESPORTMENU_LIVE {
		liveId := reqMsg.GetDataId()
		obj := for_game.GetRedisLiveRoomPlayerObj(liveId)
		if obj != nil {
			//if reqMsg.GetStatus() != for_game.ESPORTS_NEWS_STATUS_1 {
			plyIds := obj.GetPlayerIds()
			lenl := len(plyIds)
			if lenl > 0 {
				logs.Info("發了%d個", lenl)
				SendMsgToHallClientNew(plyIds, "RpcESportDataStatusInfo", reqMsg)
			} else {
				logs.Info("房間沒人", lenl)
			}
			//}
		}
	}
	return rd
}
