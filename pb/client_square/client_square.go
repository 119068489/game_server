package client_square

import (
	"game_server/easygo"
)

var UpRpc = map[string]easygo.Pair{
	"RpcLogin":                         easygo.Pair{"client_square.LoginMsg", "base.Empty"},
	"RpcFlushSquareDynamic":            easygo.Pair{"client_square.FlushInfo", "base.Empty"},
	"RpcNewVersionFlushSquareDynamic":  easygo.Pair{"client_square.NewVersionFlushInfo", "base.Empty"},
	"RpcAddSquareDynamic":              easygo.Pair{"share_message.DynamicData", "client_server.RequestInfo"},
	"RpcDelSquareDynamic":              easygo.Pair{"client_server.RequestInfo", "base.Empty"},
	"RpcZanOperateSquareDynamic":       easygo.Pair{"client_server.ZanInfo", "base.Empty"},
	"RpcAddCommentSquareDynamic":       easygo.Pair{"share_message.CommentData", "share_message.CommentData"},
	"RpcDelCommentSquareDynamic":       easygo.Pair{"client_server.IdInfo", "base.Empty"},
	"RpcAttentioPlayer":                easygo.Pair{"client_server.AttenInfo", "base.Empty"},
	"RpcGetDynamicInfo":                easygo.Pair{"client_server.IdInfo", "share_message.DynamicData"},
	"RpcGetDynamicMainComment":         easygo.Pair{"client_server.IdInfo", "share_message.CommentList"},
	"RpcGetDynamicSecondaryComment":    easygo.Pair{"client_server.IdInfo", "share_message.CommentList"},
	"RpcGetDynamicInfoNew":             easygo.Pair{"client_server.IdInfo", "share_message.DynamicData"},
	"RpcGetDynamicMainCommentNew":      easygo.Pair{"client_server.IdInfo", "share_message.CommentList"},
	"RpcGetDynamicSecondaryCommentNew": easygo.Pair{"client_server.IdInfo", "share_message.CommentList"},
	"RpcGetSquareMessage":              easygo.Pair{"client_server.IdInfo", "client_square.MessageMainInfo"},
	"RpcGetPlayerZanInfo":              easygo.Pair{"client_server.RequestInfo", "client_square.ZanList"},
	"RpcGetPlayerAttentionInfo":        easygo.Pair{"client_server.RequestInfo", "client_square.AttentionList"},
	"RpcDynamicTop":                    easygo.Pair{"client_square.DynamicTopReq", "base.Empty"},
	"RpcReadPlayerInfo":                easygo.Pair{"client_square.UnReadInfo", "base.Empty"},
	"RpcLogOut":                        easygo.Pair{"base.Empty", "base.Empty"},
	"RpcFirstLoginSquare":              easygo.Pair{"client_square.FirstLoginSquareReq", "client_square.FirstLoginSquareReply"},
	"RpcAdvDetail":                     easygo.Pair{"share_message.AdvSetting", "client_square.AdvDetailReply"},
	"RpcAddAdvLog":                     easygo.Pair{"share_message.AdvLogReq", "base.Empty"},
}

var DownRpc = map[string]easygo.Pair{
	"RpcSquareAllDynamic":           easygo.Pair{"client_square.AllInfo", "base.NoReturn"},
	"RpcNewMessage":                 easygo.Pair{"client_square.NewUnReadMessageResp", "base.NoReturn"},
	"RpcNoNewMessage":               easygo.Pair{"base.Empty", "base.NoReturn"},
	"RpcNewVersionSquareAllDynamic": easygo.Pair{"client_square.NewVersionAllInfo", "base.NoReturn"},
}
