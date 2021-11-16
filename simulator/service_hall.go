package main

import (
	"fmt"
	"game_server/easygo"
	"game_server/pb/client_hall"
	"game_server/pb/client_server"
	"log"

	"github.com/astaxie/beego/logs"
)

var _ = fmt.Sprintf
var _ = log.Println

//==========================================================================

type ServiceForHall struct {
}

func (self *ServiceForHall) RpcLogin(ep IHallEndpoint, ctx easygo.IEndpointBase, reqMsg *client_server.PlayerMsg) easygo.IMessage {
	logs.Info("测试登陆返回----------")
	//log.Printf("RpcPlayerAttrInit, msg=%v\n", reqMsg)

	return nil
}

func (self *ServiceForHall) RpcPlayerLoginResponse(ep IHallEndpoint, ctx easygo.IEndpointBase, reqMsg *client_server.AllPlayerMsg) easygo.IMessage {
	logs.Info("====RpcPlayerLoginResponse=====", reqMsg)
	playerId = reqMsg.GetMyself().GetPlayerId()
	//HeartBeat()
	return nil
}

func (self *ServiceForHall) RpcNewVersionSquareAllDynamic(ep IHallEndpoint, ctx easygo.IEndpointBase, reqMsg *client_hall.NewVersionAllInfo) easygo.IMessage {
	logs.Info("RpcNewVersionSquareAllDynamic, msg=%v\n", reqMsg)
	return nil
}
func (self *ServiceForHall) RpcToast(ep IHallEndpoint, ctx easygo.IEndpointBase, reqMsg *client_server.ToastMsg) easygo.IMessage {
	logs.Info("服务器报错了:", reqMsg)
	return nil
}
func (self *ServiceForHall) RpcPlayerAttrChange(ep IHallEndpoint, ctx easygo.IEndpointBase, reqMsg *client_server.PlayerMsg) easygo.IMessage {
	logs.Info("RpcPlayerAttrChange:", reqMsg)
	return nil
}
func (self *ServiceForHall) RpcModifyBagItem(ep IHallEndpoint, ctx easygo.IEndpointBase, reqMsg *client_hall.BagItems) easygo.IMessage {
	logs.Info("RpcModifyBagItem:", reqMsg)
	return nil
}
func (self *ServiceForHall) RpcModifyEquipment(ep IHallEndpoint, ctx easygo.IEndpointBase, reqMsg *client_hall.EquipmentReq) easygo.IMessage {
	logs.Info("RpcModifyEquipment:", reqMsg)
	return nil
}
func (self *ServiceForHall) RpcBroadCastQTXResp(ep IHallEndpoint, ctx easygo.IEndpointBase, reqMsg *client_hall.BroadCastQTX) easygo.IMessage {
	logs.Info("RpcBroadCastQTXResp:", reqMsg)
	return nil
}

/*func (self *ServiceForHall) RpcAddSquareDynamic(ep IHallEndpoint, ctx easygo.IEndpointBase, reqMsg *client_hall.NewVersionAllInfo) easygo.IMessage {
	logs.Info("RpcAddSquareDynamic, msg=%v\n", reqMsg)
	return nil
}
*/
