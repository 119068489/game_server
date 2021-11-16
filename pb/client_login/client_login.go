package client_login

import (
	"game_server/easygo"
)

var UpRpc = map[string]easygo.Pair{
	"RpcLoginHall":           easygo.Pair{"client_login.LoginMsg", "client_login.LoginResult"},
	"RpcClientGetCode":       easygo.Pair{"client_server.GetCodeRequest", "base.Empty"},
	"RpcCheckMessageCode":    easygo.Pair{"client_server.CodeResponse", "base.Empty"},
	"RpcForgetLoginPassword": easygo.Pair{"client_login.LoginMsg", "base.Empty"},
	"RpcCheckAccountVaild":   easygo.Pair{"client_server.CheckInfo", "client_server.CheckInfo"},
	"RpcAccountCancel":       easygo.Pair{"client_login.AccountCancel", "client_login.AccountCancel"},
}

var DownRpc = map[string]easygo.Pair{}
