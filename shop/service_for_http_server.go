package shop

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
	"reflect"

	"github.com/gin-gonic/gin"
)

type WebHttpForServer struct {
	for_game.WebHttpServer
	Service reflect.Value
}

func NewWebHttpForServer(port int32) *WebHttpForServer {
	p := &WebHttpForServer{}
	p.Init(port)
	return p
}

func (self *WebHttpForServer) Init(port int32) {
	services := map[string]interface{}{
		SERVER_NAME: self,
		"hall":      &ServiceForHall{},
	}
	//TODO:分发消息定义
	upRpc := easygo.CombineRpcMap(client_hall.UpRpc, server_server.UpRpc)
	self.WebHttpServer.Init(port, services, upRpc)
	self.InitRoute()
}

//初始化路由
func (self *WebHttpForServer) InitRoute() {
	self.R.POST("/api", self.ApiEntry)
}

//api入口，路由分发  RpcLogin bysf
func (self *WebHttpForServer) ApiEntry(c *gin.Context) {
	data, b := c.Get("Data")
	if !b {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(1, easygo.NewFailMsg("err ApiEntry 1")))
		return
	}
	request, ok := data.(*base.Request)
	if !ok {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(request.GetRequestId(), easygo.NewFailMsg("err ApiEntry 2")))
		return
	}
	com, b := c.Get("Common")
	if !b {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(request.GetRequestId(), easygo.NewFailMsg("err ApiEntry 3")))
		return
	}
	common, ok := com.(*base.Common)
	if !ok {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(request.GetRequestId(), easygo.NewFailMsg("err ApiEntry 4")))
		return
	}
	result := self.WebHttpServer.DealRequest(0, request, common)
	_, _ = c.Writer.Write(for_game.PacketProtoMsg(request.GetRequestId(), result))

}

//TODO 消息接收分发
func (self *WebHttpForServer) RpcMsgToOtherServer(common *base.Common, reqMsg *share_message.MsgToServer) easygo.IMessage {
	logs.Info("收到其他服务器的请求:", common, reqMsg)
	methodName := reqMsg.GetRpcName()
	var method reflect.Value
	for _, service := range self.Services {
		method = service.MethodByName(methodName)
		if method.IsValid() {
			break
		}
	}
	var msg easygo.IMessage
	if reqMsg.GetMsgName() != "" {
		msg = easygo.NewMessage(reqMsg.GetMsgName())
		err := msg.Unmarshal(reqMsg.GetMsg())
		easygo.PanicError(err)
	}
	if !method.IsValid() || method.Kind() != reflect.Func {
		logs.Info("无效的rpc请求，找不到methodName:", methodName)
		return nil
	}
	args := make([]reflect.Value, 0, 3)
	args = append(args, reflect.ValueOf(common))
	args = append(args, reflect.ValueOf(msg))
	backMsg := method.Call(args) // 分发到指定的rpc
	if backMsg != nil {
		bb, ok := backMsg[0].Interface().(easygo.IMessage)
		if ok {
			return bb
		}
	}
	return nil
}
