package client_server

import (
	"game_server/easygo"
)

var UpRpc = map[string]easygo.Pair{
	"RpcHeartbeat":      easygo.Pair{"client_server.NTP", "client_server.NTP"},
	"RpcTFToServer":     easygo.Pair{"client_server.ClientInfo", "base.NoReturn"},
	"RpcBtnClick":       easygo.Pair{"client_server.BtnClickInfo", "base.Empty"},
	"RpcPageRegLogLoad": easygo.Pair{"client_server.PageRegLogLoad", "base.Empty"},
}

var DownRpc = map[string]easygo.Pair{
	"RpcToast":               easygo.Pair{"client_server.ToastMsg", "base.NoReturn"},
	"RpcBroadCastMsg":        easygo.Pair{"share_message.BroadCastMsg", "base.NoReturn"},
	"RpcPlayerAttrChange":    easygo.Pair{"client_server.PlayerMsg", "base.NoReturn"},
	"RpcStopBroad":           easygo.Pair{"client_server.BroadIdReq", "base.NoReturn"},
	"RpcPlayerTimeoutBeKick": easygo.Pair{"client_server.PlayerTimeoutBeKick", "base.NoReturn"},
}
