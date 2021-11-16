package wish

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/share_message"
	"sync/atomic"

	"github.com/akqp2019/protobuf/proto"
	"github.com/astaxie/beego/logs"
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

const INIT_WISH_PLAYER_ID int64 = 18800000

//进入许愿池渠道
const (
	CHANNEL_NINGMANG = 1001 //柠檬畅聊im
	CHANNEL_VOICE    = 1002 //语音项目
	CHANNEL_WEB      = 1003 //网页
)

// 许愿的常量
const (
	WISH_OP_TYPE_ADD  = 1 //许愿
	WISH_OP_TYPE_EDIT = 2 //修改愿望
	WISH_OP_TYPE_DEL  = 3 //删除愿望
)

const (
	WISH_DARE_RECORED   = 1 //挑战记录
	WISH_DARE_HOLD_TIME = 2 // 占领时长
)

const (
	WISH_DARE    = 1 // 挑战赛
	WISH_NO_DARE = 2 // 非挑战赛
)

const (
	WISH_SALEOUT = 0 // 下架
	WISH_ONSALE  = 1 // 上架
	WISH_ALL     = 2 //全部
)

const (
	COIN = 1 //守护者每次被挑战增加的硬币数
)

const (
	RECYCLE_CHECKING = 0 //待审核
	RECYCLE_RECYCLED = 1 //已回收
)

const (
	EXCHANGE_ERR_STR = "兑换失败，请联系客服" //一般兑换失败返回信息
	RECYCLE_ERR_STR  = "回收失败，请联系客服" //一般回收失败返回信息
)

// 生成游戏客户端的 endpoint id (U3D 客户端 + H5 客户端)
func GenClientEndpointId() ENDPOINT_ID {
	v := atomic.AddInt32(&_ClientEndpointId, 1) // 溢出后自动回转
	return v
}

//随机发送指定服务器类型
//此方法只对从im进入调用，其他不走这里pid为第三方玩家id
func SendMsgToIdelServer(t int32, methodName string, msg easygo.IMessage, pid ...int64) (easygo.IMessage, *base.Fail) {
	playerId := append(pid, 0)[0]
	var srv *share_message.ServerInfo
	if playerId != 0 {
		player := for_game.GetWishPlayerInfo(playerId)
		if player != nil {
			srv = PServerInfoMgr.GetServerInfo(player.GetHallSid())
		}
	}
	if srv == nil {
		srv = PServerInfoMgr.GetIdelServer(t)
	}
	if srv == nil {
		return nil, easygo.NewFailMsg("无法找到指定类型服务器")
	}

	return SendMsgToServerNew(srv.GetSid(), methodName, msg, pid...)
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

// rpc 返回的公共方法
func RpcReturnCommon(rpcName string, resp easygo.IMessage, err error) easygo.IMessage {
	if err != nil {
		logs.Error("======接口名字为 %s 的 service 处理有误,err: %s", rpcName, err.Error())
		if err.Error() == "扣费失败" {
			return easygo.NewFailMsg("硬币不足")
		}
		//return easygo.NewFailMsg("参数有误")
		return easygo.NewFailMsg(err.Error())
	}
	return resp
}
