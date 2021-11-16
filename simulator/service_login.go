package main

import (
	"fmt"
	"game_server/easygo"
	"game_server/pb/client_server"
	"log"
)

var _ = fmt.Sprintf
var _ = log.Println

//==========================================================================

type ServiceForLogin struct {
}

func (self *ServiceForLogin) RpcToast(ep ILoginEndpoint, ctx easygo.IEndpointBase, reqMsg *client_server.ToastMsg) easygo.IMessage {
	log.Printf("RpcToast, msg=%v\n", reqMsg.GetText())
	return nil
}
