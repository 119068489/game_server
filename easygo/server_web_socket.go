// websocket 服务器
package easygo

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strings"
)

type IWebSocketServer interface {
	ISocketServer
	GetCrtAndKey() (crtFile string, keyFile string)
}

type WebSocketServer struct {
	SocketServer
	Me   IWebSocketServer
	Name string
}

// 地址长成这样 localhost:8080
func NewWebSocketServer(address string, name string) *WebSocketServer {
	p := &WebSocketServer{}
	p.Init(p, address, name)
	return p
}

func (self *WebSocketServer) Init(me IWebSocketServer, address string, name string) {
	self.Me = me
	self.Name = name
	self.SocketServer.Init(me, address)
}

//uid:连接唯一标识
//sid:服务器类型: 1 登录服， 2 大厅服，3 后台服，4 商场服，5统计服， 6社交广场服
func (self *WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query()
	var uid, sid string
	if u != nil {
		uid = u.Get("uid")
		sid = u.Get("sid")
		logs.Info("连接传入:", uid, sid)
	}

	upgrader := websocket.Upgrader{} // use default options
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	PanicError(err)
	//暂时只对大厅校验
	if oldCon, ok := self.ConnMap.Load(uid); ok {
		//出现自己顶号现象，标志
		logs.Error("存在相同连接:", uid, conn.RemoteAddr())
		//关闭旧的连接
		if oldCon1, ok1 := oldCon.(*EndpointBase); ok1 {
			oldCon1.SetFlag(true)
			oldCon1.Shutdown()
			logs.Info("新连接覆盖旧的连接------", oldCon1)
		}
	}
	xForwarded := r.Header.Get("X-Forwarded-For")
	if xForwarded != "" {
		address := strings.Split(xForwarded, ",")
		realAddress := address[len(address)-1]
		logs.Debug(" X-Forwarded-For 走转发服务器时玩家的真实ip;", realAddress)
		if realAddress != "" {
			tcpAddr := conn.RemoteAddr().(*net.TCPAddr)
			tcpAddr.IP = net.ParseIP(realAddress)
		}
	}
	s := fmt.Sprintf("在 (%s) 的 %v 接收到连接: %s", self.Name, self.Address, conn.RemoteAddr())
	logs.Debug(s)
	self.Me.HandleConnection(conn, uid)
}

func (self *WebSocketServer) GetCrtAndKey() (crtFile string, keyFile string) {
	return "", ""
}

func (self *WebSocketServer) Serve() { // impletement
	s := fmt.Sprintf("(%s) 开始监听: %v", self.Name, self.Address)
	logs.Info(s)

	// http.HandleFunc("/", self.Home)        // 下面的 ListenAndServe 传了第 2 个参数，这里的就会不起作用

	var err error
	crtFile, keyFile := self.Me.GetCrtAndKey()
	if crtFile != "" && keyFile != "" {
		err = http.ListenAndServeTLS(self.Address, crtFile, keyFile, self)
	} else {
		err = http.ListenAndServe(self.Address, self) // 第 2 个参数可以传 nil。传进去的 self 必须有实现一个 ServeHTTP 函数
	}
	PanicError(err)
}
