// 侦听电竞数据请求

package hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
	"reflect"
)

//===================================================================

type ServiceForESports struct {
	Service reflect.Value
}

//通知大厅硬币变化
func (self *ServiceForESports) RpcESportSendChangeCoins(common *base.Common, reqMsg *share_message.ESportCoinRecharge) easygo.IMessage {
	logs.Info("=========RpcESportSendChangeCoins============= ", reqMsg)

	if err := NotifyAddCoin(reqMsg.GetPlayerId(), reqMsg.GetRechargeCoin(), reqMsg.GetNote(), reqMsg.GetSourceType(), reqMsg.GetExtendLog()); err != "" {

		return &client_hall.ESportCommonResult{
			Code: easygo.NewInt32(for_game.C_DEDUCT_MONEY_FAIL),
			Msg:  easygo.NewString(err),
		}
	}

	return &client_hall.ESportCommonResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}
}

//通知大厅电竞币变化
func (self *ServiceForESports) RpcESportSendChangeESportCoins(common *base.Common, reqMsg *share_message.ESportCoinRecharge) easygo.IMessage {
	logs.Info("=========RpcESportSendChangeESportCoins============= ", reqMsg)

	if err := NotifyAddESportCoin(reqMsg.GetPlayerId(), reqMsg.GetRechargeCoin(), reqMsg.GetNote(), reqMsg.GetSourceType(), reqMsg.GetExtendLog()); err != "" {

		if reqMsg.GetSourceType() == for_game.ESPORTCOIN_TYPE_GUESS_BET_OUT {
			return &client_hall.ESportCommonResult{
				Code: easygo.NewInt32(for_game.C_DEDUCT_MONEY_FAIL),
				Msg:  easygo.NewString(err),
			}
		} else if reqMsg.GetSourceType() == for_game.ESPORTCOIN_TYPE_GUESS_BACK_IN {
			return &client_hall.ESportCommonResult{
				Code: easygo.NewInt32(for_game.C_SETTLEMENT_MONEY_FAIL),
				Msg:  easygo.NewString(err),
			}
		}
	}

	return &client_hall.ESportCommonResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}
}
