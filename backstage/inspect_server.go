//
package backstage

import (
	"fmt"
	"game_server/for_game"
	"strings"
	"sync/atomic"
)

type InspectServer struct {
	for_game.InspectServer
}

// 地址长成这样 localhost:8080
func NewInspectServer(address string) *InspectServer {
	p := &InspectServer{}
	p.Init(address)
	return p
}

func (is *InspectServer) Init(address string) {
	is.InspectServer.Init(is, address, "后台探查服务")
}

// override
func (is *InspectServer) GetTitle() string {
	return "后台服务器"
}

// override
func (is *InspectServer) GetBody() string {
	body1 := is.InspectServer.GetBody()
	//写入日志
	numEp := BrowerEpMp.Length()
	//numMgr := BrowerEpMgr.Length()

	numCount := atomic.LoadInt32(&LoginCount)

	list := []string{
		fmt.Sprintf("<p>玩家连接数,用玩家 id 关联的：%d,用 endpoint id 关联的: %d</p>", numEp, numEp),
		fmt.Sprintf("<p>正在处理中的登录数 %d (数量多说明很忙，特别是数据库很忙)<p>", numCount),
		"<p><a href='debug/pprof'>pprof 管理界面</a><p>",
	}
	body2 := strings.Join(list, "")
	return body1 + body2
}

var LoginCount int32 // 正在登录的数量
