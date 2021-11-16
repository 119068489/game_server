package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"
)

// 代表连接到浏览器的一个连接
type IBrowerEndpoint interface {
	easygo.IEndpointWithWebSocket
	brower_backstage.IBackstage2Brower

	SetUser(user *share_message.Manager)
	GetUser() *share_message.Manager
}

type BrowerEndpoint struct {
	easygo.EndpointWithWebSocket
	brower_backstage.Backstage2Brower

	User *share_message.Manager // 关联的用户对象
}

func NewBrowerEndpoint(endpointId ENDPOINT_ID) *BrowerEndpoint {
	p := &BrowerEndpoint{}
	p.Init(p, endpointId)
	return p
}

func (be *BrowerEndpoint) Init(me IBrowerEndpoint, endpointId ENDPOINT_ID) {
	services := map[string]interface{}{
		// 处理 rpc 请求的 service 类添写到这里
		"ServiceForBrower": &ServiceForBrower{},
	}

	be.EndpointWithWebSocket.Init(me, services, brower_backstage.UpRpc, brower_backstage.DownRpc, endpointId) // 写死
	be.Backstage2Brower.Init(me)
}

func (be *BrowerEndpoint) CreateFailMsgWhenPanic(o interface{}, methodName string, requestId uint64) *base.Fail { // override
	s := fmt.Sprintf("操作失败:%v", o)
	return &base.Fail{Reason: &s}
}

func (be *BrowerEndpoint) GetContexForDealRequest(methodName string, requestId uint64) (interface{}, interface{}) { // override
	if methodName == "RpcLogin" || methodName == "RpcGetCode" || methodName == "RpcSignin" {
		return be, be
	}
	if be.User == nil {
		return be, fmt.Sprintf("收到 %s,非法操作，请先登录", methodName)
	}
	return be, be.User
}

func (be *BrowerEndpoint) SetUser(user *share_message.Manager) {
	be.User = user
}

func (be *BrowerEndpoint) GetUser() *share_message.Manager {
	return be.User
}
