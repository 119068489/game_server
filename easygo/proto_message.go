package easygo

import (
	"game_server/easygo/base"
	// "sync"
	// "time"
	"reflect"
	//"fmt"

	"github.com/akqp2019/protobuf/proto"
)

//从各个 *.proto 生成的 *.go 文件中找出的规律
type IMessage interface {
	Reset()
	String() string
	ProtoMessage()
	Descriptor() ([]byte, []int)

	Marshal() (dAtA []byte, err error)
	MarshalTo(dAtA []byte) (int, error)
	Unmarshal(dAtA []byte) error
}

type IMessageSender interface {
	CallRpcMethod(methodName string, msg IMessage, common ...*base.Common) (IMessage, IRpcInterrupt)
}

func NewMessage(name string) IMessage {
	t := proto.MessageType(name)
	msg := reflect.New(t.Elem()).Interface()
	return msg.(IMessage)
}
func Marshal(msg IMessage) []byte {
	bytes, err := msg.Marshal()
	if err != nil {
		panic(err)
	}
	return bytes
}
func Unmarshal(msg IMessage, data []byte) {
	err := msg.Unmarshal(data)
	if err != nil {
		panic(err)
	}
}
