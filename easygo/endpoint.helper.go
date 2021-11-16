package easygo

import (
	"net"
)

type ResponseInfo struct {
	RpcOk bool
	Bytes []byte
}

// 把 net.Conn 和 websocket.Conn 的都有的函数都抽取在一块
type EasygoConn interface {
	RemoteAddr() net.Addr
	Close() error
	// 有需要再补齐剩下函数
}
