package main

import (
	"fmt"
	"game_server/easygo"
	"game_server/pb/client_hall"
	"game_server/pb/client_login"
	"log"
	"time"

	"github.com/astaxie/beego/logs"
)

var _ = fmt.Sprintf
var _ = log.Println

func (self *Command) RpcLoginHall(param ...string) {

	msg := client_login.LoginMsg{
		Account:  easygo.NewString(param[0]),
		Password: easygo.NewString(param[1]),
		Type:     easygo.NewInt32(2),
	}
	backMsg := loginConnect.FetchEndpoint().RpcLoginHall(&msg)
	logs.Info("登录请求返回：", backMsg)
	loginConnect.SetIsStop(true)
	loginConnect.FetchEndpoint().Shutdown()
	//addr := backMsg.GetAddress()
	addr := "127.0.0.1:2001"
	hallConnect = NewHallConnector(addr)
	hallConnect.ConnectOnce()
	time.Sleep(1 * time.Second)
	logs.Info(backMsg)
	self.RpcLogin(&client_hall.LoginMsg{
		Account:        backMsg.Account,
		Token:          backMsg.Token,
		RegistrationId: easygo.NewString(""),
		Channel:        easygo.NewString(""),
		LoginType:      easygo.NewInt32(1),
		DeviceType:     easygo.NewInt32(2),
		Type:           easygo.NewInt32(2),
		PlayerId:       backMsg.PlayerId,
		VersionNumber:  nil,
		Brand:          nil,
	})

	//self.RpcNewVersionFlushSquareDynamic() // 刷新动态
	//self.RpcAddSquareDynamic() // 添加动态
	//self.RpcDelSquareDynamic() // 添加动态
	//self.RpcDelNewFriendList() // 删除我的朋友列表中的某些人.
}
