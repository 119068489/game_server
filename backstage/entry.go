// 大厅服务
package backstage

import (
	"flag"
	"game_server/easygo"
	"game_server/for_game"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/astaxie/beego/logs"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func Entry(flagSet *flag.FlagSet, args []string) {
	initializer := for_game.NewInitializer()
	defer func() { // 若是异常了,确保异步日志有成功写盘
		logger := initializer.GetBeeLogger()
		if logger != nil {
			logger.Flush()
		}
	}()

	dict := easygo.KWAT{
		"logName":  SERVER_NAME,
		"yamlPath": "config_backstage.yaml",
	}
	initializer.Execute(dict)

	Initialize()

	InitData()                                          //初始化预设数据到数据库
	TimedTask()                                         //定时任务
	for_game.DelRedisWaiterForSid(PServerInfo.GetSid()) //清理本服不在线客服
	for_game.DelRedisAdminForSid(PServerInfo.GetSid())  //清理本服不在线管理员

	//启动etcd
	PClient3KVMgr.StartClintTV3()
	defer PClient3KVMgr.Close() //关闭etcd
	var ServeFunctions = []func(){}
	ServeFunctions = append(ServeFunctions, SignHandle)
	ServeFunctions = append(ServeFunctions, WebApiServreMgr.Serve)
	//web api初始化
	PWebApiForClient = NewWebHttpForClient(PServerInfo.GetClientApiPort())
	PWebApiForServer = NewWebHttpForServer(PServerInfo.GetServerApiPort())
	ServeFunctions = append(ServeFunctions, PWebApiForClient.Serve)
	ServeFunctions = append(ServeFunctions, PWebApiForServer.Serve)
	ServeFunctions = append(ServeFunctions, Server4Brower.Serve)
	//etcd注册和发现
	for_game.InitExistServer(PClient3KVMgr, PServerInfoMgr, PServerInfo)
	jobs := []easygo.IGoroutine{}

	for _, f := range ServeFunctions {
		job := easygo.Spawn(f)
		jobs = append(jobs, job)
	}
	easygo.JoinAllJobs(jobs...)
}
func SignHandle() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM)
	for {
		s := <-c
		switch s {
		case syscall.SIGTERM:
			//TODO:服务器关闭逻辑处理
			logs.Info("后台服务器关闭", PServerInfo)
			PServerInfoMgr.DelServerInfo(PServerInfo.GetSid())
			time.Sleep(time.Second * 10)
			os.Exit(1)
		default:
			logs.Debug("error sign", s)
		}
	}
}
