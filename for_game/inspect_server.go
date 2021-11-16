package for_game

import (
	"fmt"
	"game_server/easygo"
	"net"
	"net/http"
	"net/http/pprof"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
)

type IInspectServer interface {
	http.Handler
	Serve()
	Home(w http.ResponseWriter, r *http.Request)
	GetTitle() string
	GetBody() string
}

type InspectServer struct {
	http.Server
	http.ServeMux
	Me      IInspectServer
	Address string
	Name    string
}

// 地址长成这样 localhost:8080
func NewInspectServer(address string, name string) *InspectServer {
	p := &InspectServer{}
	p.Init(p, address, name)
	return p
}

func (self *InspectServer) Init(me IInspectServer, address string, name string) {
	self.Me = me
	self.Address = address
	self.Name = name
	self.Server.Addr = address
	self.Server.Handler = me // me/self 从 ServeMux 继承了 ServeHTTP 函数

	self.HandleFunc("/debug/pprof/", pprof.Index)
	self.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	self.HandleFunc("/debug/pprof/profile", pprof.Profile)
	self.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	self.HandleFunc("/debug/pprof/trace", pprof.Trace)

	self.HandleFunc("/", self.Me.Home)
}

func (self *InspectServer) GetTitle() string {
	panic("请在子类 override 此方法")
}

func (self *InspectServer) GetBody() string {
	stats := easygo.MongoStats()
	return "MongoDB 驱动统计：<br/>" + strings.Join(stats, "<br/>")
}

func (self *InspectServer) Home(w http.ResponseWriter, r *http.Request) {
	title := self.Me.GetTitle()
	body := self.Me.GetBody()
	content := fmt.Sprintf(_HTML, title, body)
	w.Write([]byte(content))
}

func (self *InspectServer) Serve() {
	s := fmt.Sprintf("(%s) 开始监听: %v", self.Name, self.Address)
	logs.Info(s)
	self.ListenAndServe()
	// http.HandleFunc("/", self.Home)
	// err := http.ListenAndServe(self.Address, self) // 第 2 个参数可以传 nil。传进去的 self 必须有实现一个 ServeHTTP 函数
}

func (self *InspectServer) GetPort() int32 {
	_, sport, err := net.SplitHostPort(self.Address) //获取客户端IP
	easygo.PanicError(err)
	p, err := strconv.Atoi(sport)
	easygo.PanicError(err)
	return int32(p)
}

var _HTML = `
<html>
<body>
<h3>%s</h3>
%s
</body>
</html>
`
