// 大厅给游戏客户端连接的 server

package hall

import (
	"game_server/for_game"

	"github.com/astaxie/beego/logs"

	"game_server/easygo"
	"game_server/pb/client_hall"
	"game_server/pb/client_server"
)

type TcpServerForClient struct {
	easygo.TcpSocketServer
}

func NewTcpServerForClient(address string) *TcpServerForClient {
	p := &TcpServerForClient{}
	p.Init(p, address, "大厅 for TCP")
	return p
}

func (self *TcpServerForClient) CreateEndpoint(endpointId ENDPOINT_ID) easygo.IEndpointWithSocket { // override
	endpointId = GenClientEndpointId() // 另外生成一个全局唯一的 id ,因为允许 H5 客户端和 U3D 客户端同时进游戏

	upRpc := easygo.CombineRpcMap(client_hall.UpRpc, client_server.UpRpc)
	downRpc := easygo.CombineRpcMap(client_hall.DownRpc, client_server.DownRpc)
	ep := NewGameClientTcpEndpoint(upRpc, downRpc, endpointId)
	event := ep.GetDisconnectedEvent()
	event.AddHandler(self.OnDisConnect)
	ClientEpMgr.Store(endpointId, ep)
	return ep
}

func (self *TcpServerForClient) OnDisConnect(ep IGameClientEndpoint) { // override
	logs.Info("==============TcpServerForClient.OnDisConnect===============")
	epId := ep.GetEndpointId()
	ClientEpMgr.Delete(epId)
	player := ep.GetPlayer()
	if player != nil {
		HandleAfterLogout(ep, player)
	}
	logs.Info("游戏客户端到大厅的连接断开了。epId=%d,addr=%v", epId, ep.GetAddr())

}

func HandleAfterLogout(ep IGameClientEndpoint, player *Player) {
	if ep.GetFlag() { //顶号连接
		logs.Error("玩家顶号连接，不释放玩家内存")
		return
	}
	playerId := player.GetPlayerId()
	ClientEpMp.Delete(playerId)
	PlayerOnlineMgr.PlayerOffline(playerId)

	player.SetIsOnLine(false) //设置玩家在线状态为下线
	player.SetSid(0)          //下线了设置sid
	player.UpdateLogoutTimestamp()
	player.UpdateOnlineTime()
	//player.SetShopServerId(0)
	player.UpdateLastLoginIP(ep.GetConnection().RemoteAddr().String())
	NotifyPlayerOffLine(playerId)
	PlayerMgr.Delete(player.GetPlayerId())
	fun := func() {
		for_game.AddOnlineTimeLog(playerId)                                //更新玩家在线时长日志
		for_game.MakePlayerBehaviorReport(2, playerId, nil, nil, nil, nil) //生成用户行为报表的1次会话用户数 已优化到Redis
	}
	easygo.Spawn(fun)
	player.SaveToMongo()
	//player.DelRedisKey()
	//PlayerMgr.Delete(playerId)
	logs.Info("玩家离线了:", playerId, ep)
}

//-------------------------------------------------

func init() {
	PackageInitialized.AddHandler(afterHallInitialized1)
}

func afterHallInitialized1() {
	//address := for_game.MakeAddress(PServerInfo.GetIp(), PServerInfo.GetTcpClientPort())
	address := for_game.MakeAddress("0.0.0.0", PServerInfo.GetClientTCPPort())
	server := NewTcpServerForClient(address)
	ServeFunctions = append(ServeFunctions, server.Serve)
}
