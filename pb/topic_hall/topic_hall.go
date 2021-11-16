package topic_hall

import (
	"game_server/easygo"
)

var UpRpc = map[string]easygo.Pair{
	"RpcTopic2HallClient": easygo.Pair{"share_message.MsgToClient", "base.NoReturn"},
}

var DownRpc = map[string]easygo.Pair{
	"RpcHall2Topic": easygo.Pair{"share_message.MsgToClient", "base.NoReturn"},
}
