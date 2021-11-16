// 大厅给游戏客户端连接的 server

package login

import (
	"game_server/for_game"
	"game_server/pb/client_login"
	"github.com/astaxie/beego/logs"

	"game_server/easygo"
	"game_server/pb/client_server"
)

type TcpServerForClient struct {
	easygo.TcpSocketServer
}

func NewTcpServerForClient(address string) *TcpServerForClient {
	p := &TcpServerForClient{}
	p.Init(p, address, "Login服务器 for TCP")
	return p
}

func (self *TcpServerForClient) CreateEndpoint(endpointId ENDPOINT_ID) easygo.IEndpointWithSocket { // override
	endpointId = GenClientEndpointId() // 另外生成一个全局唯一的 id ,因为允许 H5 客户端和 U3D 客户端同时进游戏

	upRpc := easygo.CombineRpcMap(client_login.UpRpc, client_server.UpRpc)
	downRpc := easygo.CombineRpcMap(client_login.DownRpc, client_server.DownRpc)

	ep := NewGameClientTcpEndpoint(upRpc, downRpc, endpointId)
	event := ep.GetDisconnectedEvent()
	event.AddHandler(self.OnDisConnect)
	ClientEpMgr.Store(endpointId, ep)
	ep.SetConType(for_game.CONN_TYPE_TCP)
	return ep
}

func (self *TcpServerForClient) OnDisConnect(ep IGameClientEndpoint) { // override
	epId := ep.GetEndpointId()
	ClientEpMgr.Delete(epId)
	player := ep.GetPlayer()
	if player != nil {
		playerId := player.GetPlayerId()
		ClientEpMp.Delete(playerId)
		HandleAfterLogout(ep, player)
	}
	logs.Info("游戏客户端到大厅的连接断开了。epId=%d,addr=%v", epId, ep.GetAddr())

}

func HandleAfterLogout(ep IGameClientEndpoint, player *Player) {
}

//-------------------------------------------------

func init() {
	PackageInitialized.AddHandler(afterHallInitialized1)
}

func afterHallInitialized1() {
	//address := for_game.MakeAddress(PServerInfo.GetIp(), PServerInfo.GetTcpClientPort())
	//默认监听本机
	address := for_game.MakeAddress("0.0.0.0", PServerInfo.GetClientTCPPort())
	server := NewTcpServerForClient(address)
	ServeFunctions = append(ServeFunctions, server.Serve)
}
