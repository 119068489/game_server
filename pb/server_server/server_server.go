package server_server

import (
	"game_server/easygo"
)

var UpRpc = map[string]easygo.Pair{
	"RpcMsgToHallClient":  easygo.Pair{"share_message.MsgToClient", "base.NoReturn"},
	"RpcMsgToOtherServer": easygo.Pair{"share_message.MsgToServer", "share_message.MsgToServer"},
}

var DownRpc = map[string]easygo.Pair{}
