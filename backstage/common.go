package backstage

import (
	"encoding/base64"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"

	"github.com/akqp2019/protobuf/proto"

	"github.com/astaxie/beego/logs"
)

type SITE = string
type USER_ID = int64
type USER_ACCOUNT = string
type PLAYER_ID = int64
type PLAYER_IDS = []int64
type TEAM_ID = int64
type ENDPOINT_ID = int32
type DB_NAME = string
type SERVER_ID = int32
type INSTANCE_ID = int32
type LEVEL = int32
type TYPEID = int32 //后台用户类型
type SHOP_ORDER_ID = int64
type SHOP_ITEM_ID = int64
type TIME_64 = int64 //int64时间类型
type MAP_STRING = map[string]string

//客服类型
const (
	WAITER_FAQ int32 = 3
	WAITER_ORD int32 = 4
)

//快手广告投放场景  1-优选广告，2-信息流广告（旧投放场景，含上下滑大屏广告），3-视频播放页广告，6-上下滑大屏广告，7-信息流广告（不含上下滑大屏广告）
const (
	KS_ADV_SCENES1 = 1
	KS_ADV_SCENES2 = 2
	KS_ADV_SCENES3 = 3
	KS_ADV_SCENES6 = 6
	KS_ADV_SCENES7 = 7
)

//用户投诉类型集合
var COMPLAINT_OTHER_TYPE = []int32{3, 6, 10, 12, 13, 19, 20, 22, 23, 24, 25}
var COMPLAINT_OTHER_REASON = []int32{4, 5}

//电竞资讯类型
const (
	//资讯类型
	ES_NEWS   = 1 //新闻
	ES_VIDEO  = 2 //视频
	ES_SYSMSG = 3 //系统消息

	//资讯发布类型
	ES_SEND_TYPE_NOW    = 1 //立即发送
	ES_SEND_TYPE_FUTURE = 2 //定时未来发送
)

/**
api 接口参数校验
m 参数map
method: 方法名
*/
func VerifyParams(m map[string]string, method string) string {
	if len(m) == 0 {
		logs.Error("api 接口名: %s 中,所有请求参数为空", method)
		return "参数为空"
	}
	for k, v := range m {
		if v == "" {
			logs.Error("api 接口名: %s 中,请求参数 %s 为空", method, k)
			return fmt.Sprintf("参数 %s 不能为空", k)
		}
	}
	return ""
}

//检查文字内容是否违规
func CheckOjbScore(content string, types int32) *brower_backstage.CheckScoreResponse {
	msg := &brower_backstage.CheckScoreResponse{}
	switch types {
	case 1:
		txt := base64.StdEncoding.EncodeToString([]byte(content))
		result := for_game.TextModeration(txt)
		if result == nil {
			return nil
		}
		msg.EvilFlag = easygo.NewInt32(result.EvilFlag)
		msg.EvilType = easygo.NewInt32(result.EvilType)
		if result.EvilType == 100 {
			msg.Score = easygo.NewInt32(0)
		} else {
			List := result.DetailResult
			for _, item := range List {
				if item.EvilType == result.EvilType {
					msg.Score = easygo.NewInt32(item.Score)
					break
				}
			}
		}
	case 2:
		result := for_game.ImageModeration(content)
		if result == nil {
			return nil
		}

		msg.EvilFlag = easygo.NewInt32(result.EvilFlag)
		msg.EvilType = easygo.NewInt32(result.EvilType)

		switch result.EvilType {
		case 20001:
			msg.Score = easygo.NewInt32(result.PolityDetect.Score)
		case 20002:
			msg.Score = easygo.NewInt32(result.PornDetect.Score)
		case 20006:
			msg.Score = easygo.NewInt32(result.IllegalDetect.Score)
		case 20007:
			msg.Score = easygo.NewInt32(0)
		case 20103:
			msg.Score = easygo.NewInt32(result.HotDetect.Score)
		case 24001:
			msg.Score = easygo.NewInt32(result.TerrorDetect.Score)
		default:
			msg.Score = easygo.NewInt32(0)
		}

	}

	return msg
}

//后台给大厅发送消息分三种：
//1、广播给所有大厅
//2、随机一台大厅发送，大厅收到特殊处理

//广播给所有大厅
//func BroadCastToAllHall(methodName string, msg easygo.IMessage) {
//	halls := PServerInfoMgr.GetAllServers(for_game.SERVER_TYPE_HALL)
//	for _, hall := range halls {
//		if hall != nil {
//			if msg == nil {
//				msg = &base.Empty{}
//			}
//			PWebApiForServer.SendToServer(hall, methodName, msg)
//		} else {
//			logs.Info("不存在的服务器连接:", hall)
//		}
//	}
//}

//指定在线玩家所在服通知
func SendToPlayer(playerId PLAYER_ID, methodName string, msg easygo.IMessage) easygo.IMessage {
	base := for_game.GetRedisPlayerBase(playerId)
	if base == nil {
		return easygo.NewFailMsg("用户不存在")
	}
	sid := base.GetSid()
	return ChooseOneHall(sid, methodName, msg)
}

//指定大厅通知，sid=0则随机获取一台
func ChooseOneHall(sid SERVER_ID, methodName string, msg easygo.IMessage) easygo.IMessage {
	var hall *share_message.ServerInfo
	if sid == 0 {
		hall = PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_HALL)
	} else {
		hall = PServerInfoMgr.GetServerInfo(sid)
	}
	if hall != nil {
		resp, err := SendMsgToServerNew(hall.GetSid(), methodName, msg)
		if err != nil {
			return err
		}
		return resp

	} else {
		logs.Info("不存在的服务器连接:", hall)
	}
	return nil
}

//通知其他后台服务器，用户上线
func NotifyPlayerOnLine(user *share_message.Manager) {
	msg := &share_message.PlayerState{
		PlayerId: easygo.NewInt64(user.GetId()),
		ServerId: easygo.NewInt32(PServerInfo.GetSid()),
	}
	logs.Info("通知其他后台服务器，用户上线,msg------->%+v", user.GetId())
	//通知其他后台服务器
	BroadCastMsgToServerNew(for_game.SERVER_TYPE_BACKSTAGE, "RpcUserOnLine", msg)
}

//服务器间通讯通用
func SendMsgToServerNew(sid SERVER_ID, methodName string, msg easygo.IMessage, pid ...int64) (easygo.IMessage, *base.Fail) {
	srv := PServerInfoMgr.GetServerInfo(sid)
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
	return PWebApiForServer.SendToServer(srv, "RpcMsgToOtherServer", req, playerId)
}

//广播给指定类型服务器
func BroadCastMsgToServerNew(t int32, methodName string, msg easygo.IMessage, pid ...int64) {
	servers := PServerInfoMgr.GetAllServers(t)
	for _, srv := range servers {
		if srv.GetSid() == PServerInfo.GetSid() {
			continue
		}
		if srv == nil {
			continue
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
		PWebApiForServer.SendToServer(srv, "RpcMsgToOtherServer", req)
	}
}

//通知给指定类型随机一台服务器
func SendMsgRandToServerNew(t int32, methodName string, msg easygo.IMessage, pid ...int64) (easygo.IMessage, *base.Fail) {
	srv := PServerInfoMgr.GetIdelServer(t)
	if srv == nil {
		return nil, easygo.NewFailMsg("服务器不存在")
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

	return PWebApiForServer.SendToServer(srv, "RpcMsgToOtherServer", req)
}

//直接发送给前端玩家的消息 methodName-ClientToHallName
func SendMsgToHallClientNew(playerIds []int64, methodName string, msg easygo.IMessage) {
	serversInfo := make(map[int32][]int64)
	for _, pid := range playerIds { //群发 每个人都发
		player := for_game.GetRedisPlayerBase(pid)
		if player == nil {
			continue
		}
		serversInfo[player.GetSid()] = append(serversInfo[player.GetSid()], pid)
	}
	for sid, pList := range serversInfo {
		srv := PServerInfoMgr.GetServerInfo(sid)
		if srv == nil {
			continue
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
			PlayerIds: pList,
			RpcName:   easygo.NewString(methodName),
			MsgName:   easygo.NewString(proto.MessageName(msg)),
			Msg:       msgByte,
		}
		PWebApiForServer.SendToServer(srv, "RpcMsgToHallClient", req)
	}
}

// 向许愿池发起接口调用
func ChooseOneWish(sid SERVER_ID, methodName string, msg easygo.IMessage) easygo.IMessage {
	var wish *share_message.ServerInfo
	if sid == 0 {
		wish = PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_WISH)
	} else {
		wish = PServerInfoMgr.GetServerInfo(sid)
	}
	if wish != nil {
		resp, err := SendMsgToServerNew(wish.GetSid(), methodName, msg)
		if err != nil {
			return err
		}
		return resp
	} else {
		logs.Info("不存在的服务器连接:", wish)
		return easygo.NewFailMsg("找不到许愿池服务器")
	}
}

// 向许愿池发起接口调用
func ChooseRpcWish(sid SERVER_ID, methodName string, msg easygo.IMessage) easygo.IMessage {
	var wish *share_message.ServerInfo
	if sid == 0 {
		wish = PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_WISH)
	} else {
		wish = PServerInfoMgr.GetServerInfo(sid)
	}
	if wish != nil {
		_, err := SendMsgToServerNew(wish.GetSid(), methodName, msg)
		if err != nil {
			return err
		}
		return nil

	}
	logs.Info("不存在的服务器连接:", wish)
	return &base.Fail{
		Reason: easygo.NewString("不存在的服务器连接"),
		Code:   easygo.NewString("不存在的服务器连接"),
	}
}
