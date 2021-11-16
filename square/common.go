package square

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/share_message"
	"sync/atomic"

	"github.com/akqp2019/protobuf/proto"
)

type ITEM_ID = int32
type PLAYER_ID = int64
type ENDPOINT_ID = easygo.ENDPOINT_ID
type DB_NAME = string
type SERVER_ID = int32
type INSTANCE_ID = int32
type GAME_TYPE = int32
type SITE = string
type LEVEL = int32

var _ClientEndpointId ENDPOINT_ID = 0
var GetPlayerPeriod = for_game.GetPlayerPeriod
var PackageInitialized = easygo.NewEvent()

// 生成游戏客户端的 endpoint id (U3D 客户端 + H5 客户端)
func GenClientEndpointId() ENDPOINT_ID {
	v := atomic.AddInt32(&_ClientEndpointId, 1) // 溢出后自动回转
	return v
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
	return PWebApiForServer.SendToServer(srv, "RpcMsgToOtherServer", req)
}

//广播给指定类型服务器
func BroadCastMsgToServerNew(t int32, methodName string, msg easygo.IMessage, pid ...int64) {
	servers := PServerInfoMgr.GetAllServers(t)
	for _, srv := range servers {
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

//发送给指定大厅指定玩家发送消息
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
