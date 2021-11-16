package hall

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"reflect"

	"github.com/astaxie/beego/logs"

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
		"hall":      &ServiceForHall{},
		"backstage": &ServiceForBackStage{},
		"shop":      &ServiceForShop{},
		"square":    &ServiceForSquare{},
		"esports":   &ServiceForESports{},
		"server":    self,
	}
	//TODO:分发消息定义
	upRpc := easygo.CombineRpcMap(client_hall.DownRpc, server_server.UpRpc)
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
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(1, easygo.NewFailMsg("err ApiEntry 2")))
		return
	}
	com, b := c.Get("Common")
	if !b {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(1, easygo.NewFailMsg("err ApiEntry 3")))
		return
	}
	common, ok := com.(*base.Common)
	if !ok {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(1, easygo.NewFailMsg("err ApiEntry 4")))
		return
	}
	result := self.WebHttpServer.DealRequest(0, request, common)
	if _, ok := result.(*base.NoReturn); ok {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(1, result))
		return
	}
	_, _ = c.Writer.Write(for_game.PacketProtoMsg(1, result))
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
	args := make([]reflect.Value, 0, 2)
	args = append(args, reflect.ValueOf(common))
	if msg == nil {
		msg = easygo.EmptyMsg
	}
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

//直接发送给指定客户端
func (self *WebHttpForServer) RpcMsgToHallClient(common *base.Common, reqMsg *share_message.MsgToClient) easygo.IMessage {
	for _, pid := range reqMsg.GetPlayerIds() {
		//直接推送给客户端
		ep := ClientEpMp.LoadEndpoint(pid)
		if ep == nil {
			continue
		}
		var msg easygo.IMessage
		if reqMsg.GetMsgName() != "" {
			msg = easygo.NewMessage(reqMsg.GetMsgName())
			err := msg.Unmarshal(reqMsg.GetMsg())
			easygo.PanicError(err)
		}
		_, err1 := ep.CallRpcMethod(reqMsg.GetRpcName(), msg)
		easygo.PanicError(err1)

	}
	return nil
}
