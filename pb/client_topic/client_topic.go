package client_topic

import (
	"game_server/easygo"
)

var UpRpc = map[string]easygo.Pair{
	"RpcLogin": easygo.Pair{"client_topic.LoginMsg", "base.Empty"},
}

var DownRpc = map[string]easygo.Pair{
	"RpcLoginResp": easygo.Pair{"client_topic.LoginMsg", "base.Empty"},
}
