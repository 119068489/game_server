// 大厅给游戏客户端连接的 server

package login

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/client_login"
	"game_server/pb/client_server"

	"github.com/astaxie/beego/logs"
)

type WebServerForClient struct {
	easygo.WebSocketServer
}

func NewWebServerForClient(address string) *WebServerForClient {
	p := &WebServerForClient{}
	p.Init(p, address, "Login服务器 for H5")
	return p
}

func (self *WebServerForClient) CreateEndpoint(endpointId ENDPOINT_ID) easygo.IEndpointWithSocket { // override
	endpointId = GenClientEndpointId() // 另外生成一个全局唯一的 id ,因为允许 H5 客户端和 U3D 客户端同时进游戏

	upRpc := easygo.CombineRpcMap(client_login.UpRpc, client_server.UpRpc)
	downRpc := easygo.CombineRpcMap(client_login.DownRpc, client_server.DownRpc)

	ep := NewGameClientWebEndpoint(upRpc, downRpc, endpointId)
	logs.Info("新连接:", ep.GetEndpointId())
	event := ep.GetDisconnectedEvent()
	event.AddHandler(self.OnDisConnect)
	ClientEpMgr.Store(endpointId, ep)
	ep.SetConType(for_game.CONN_TYPE_WEBSOCKET)
	return ep
}

func (self *WebServerForClient) OnDisConnect(ep IGameClientEndpoint) { // override
	logs.Info("释放uid:", ep.GetUid())
	self.ConnMap.Delete(ep.GetUid())
	epId := ep.GetEndpointId()
	ClientEpMgr.Delete(epId)
	player := ep.GetPlayer()
	if player != nil {
		playerId := player.GetPlayerId()
		ClientEpMp.Delete(playerId)
		HandleAfterLogout(ep, player)
	}
	logs.Info("客户端到大厅的连接断开了。epId=%d", epId)
}

func (self *WebServerForClient) GetCrtAndKey() (crtFile string, keyFile string) { // override
	crtPath := easygo.YamlCfg.GetValueAsString("TLS_CRT_FILE_PATH", "")
	keyPath := easygo.YamlCfg.GetValueAsString("TLS_KEY_FILE_PATH", "")
	return crtPath, keyPath
}

//-------------------------------------------------

func init() {
	PackageInitialized.AddHandler(afterHallInitialized2)
}

func afterHallInitialized2() {
	//address := for_game.MakeAddress(PServerInfo.GetIp(), PServerInfo.GetClientPort())
	//默认监听本机
	address := for_game.MakeAddress("0.0.0.0", PServerInfo.GetClientWSPort())
	server := NewWebServerForClient(address)
	ServeFunctions = append(ServeFunctions, server.Serve)
}
