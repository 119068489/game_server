package client_esports

import (
	"game_server/easygo"
)

var UpRpc = map[string]easygo.Pair{
	"RpcLogin":  easygo.Pair{"client_esport.LoginMsg", "base.Empty"},
	"RpcLogOut": easygo.Pair{"base.Empty", "base.Empty"},
}

var DownRpc = map[string]easygo.Pair{
	"RpcPushBroadcast": easygo.Pair{"client_esport.BroadcastMsg", "base.NoReturn"},
}
