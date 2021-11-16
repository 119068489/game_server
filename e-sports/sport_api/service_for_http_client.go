package sport_api

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"github.com/gin-gonic/gin"
)

type WebHttpForClient struct {
	for_game.WebHttpServer
}

func NewWebHttpForClient(port int32) *WebHttpForClient {
	p := &WebHttpForClient{}
	p.Init(port)
	return p
}

func (self *WebHttpForClient) Init(port int32) {
	services := map[string]interface{}{
		SERVER_NAME: self,
	}
	self.WebHttpServer.Init(port, services, client_hall.UpRpc)
	self.InitRoute()
}

//初始化路由
func (self *WebHttpForClient) InitRoute() {
	self.R.POST("/api", self.ApiEntry)
}

//api入口，路由分发  RpcLogin bysf
func (self *WebHttpForClient) ApiEntry(c *gin.Context) {
	data, b := c.Get("Data")
	if !b {
		c.Writer.Write(for_game.PacketProtoMsg(1, easygo.NewFailMsg("err ApiEntry 1")))
		return
	}
	request, ok := data.(*base.Request)
	if !ok {
		c.Writer.Write(for_game.PacketProtoMsg(1, easygo.NewFailMsg("err ApiEntry 2")))
		return
	}
	com, b := c.Get("Common")
	if !b {
		c.Writer.Write(for_game.PacketProtoMsg(1, easygo.NewFailMsg("err ApiEntry 3")))
		return
	}
	common, ok := com.(*base.Common)
	if !ok {
		c.Writer.Write(for_game.PacketProtoMsg(1, easygo.NewFailMsg("err ApiEntry 4")))
		return
	}
	result := self.WebHttpServer.DealRequest(0, request, common)

	c.Writer.Write(for_game.PacketProtoMsg(1, result))

}
