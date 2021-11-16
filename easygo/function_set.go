package easygo

import (
	"game_server/easygo/base"
	"time"
)

// function set
type IFunctionSet interface {
	GetTimestamp() int64
	MakeRequestPacket(methodName string, reqMsg IMessage, requestId uint64, timeout time.Duration, common ...*base.Common) []byte
}
type FunctionSet struct {
	Me              IFunctionSet
	TimestampOffset int64 // 为了方便测试，时间戳偏移量
}

func NewFunctionSet() *FunctionSet {
	p := &FunctionSet{}
	p.Init(p)
	return p
}

func (self *FunctionSet) Init(me IFunctionSet) {
	self.Me = me
	self.TimestampOffset = 0
}

func (self *FunctionSet) GetTimestamp() int64 {
	return time.Now().UnixNano()/1e6 + self.TimestampOffset // 转成毫秒
}

func (self *FunctionSet) SetTimestampOffset(offset int64) {
	self.TimestampOffset = offset
}

func (self *FunctionSet) MakeRequestPacket(methodName string, reqMsg IMessage, requestId uint64, timeout time.Duration, common ...*base.Common) []byte { //# requestId==AUTO_REQUEST_ID表示需生成id
	// requestId=0, timeout=0
	// 字段值若是为默认值,则不设置,节省流量
	if reqMsg == nil {
		reqMsg = EmptyMsg
	}
	com := append(common, nil)[0]

	request := &base.Request{}
	request.MethodName = &methodName
	request.Serialized = Marshal(reqMsg)
	request.Common = com
	if requestId != 0 {
		request.RequestId = &requestId
	}
	if timeout != 0 {
		t := uint32(timeout / time.Millisecond)
		request.Timeout = &t
		stamp := self.Me.GetTimestamp()
		request.Timestamp = &stamp
	}
	t := base.PacketType_TYPE_REQUEST
	packet := base.Packet{Type: &t, Serialized: Marshal(request)}
	return Marshal(&packet)
}
