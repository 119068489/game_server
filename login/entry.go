// 大厅服务
package login

import (
	"flag"
	"game_server/easygo"
	"game_server/for_game"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/astaxie/beego/logs"
)

func Entry(flagSet *flag.FlagSet, args []string) {
	initializer := NewInitializer()
	defer func() { // 若是异常了,确保异步日志有成功写盘
		logger := initializer.GetBeeLogger()
		if logger != nil {
			logger.Flush()
		}
	}()
	dict := easygo.KWAT{
		"logName":  SERVER_NAME,
		"yamlPath": "config_login.yaml",
	}
	initializer.Execute(dict)
	serverState := "当前服务器状态:"
	if for_game.IS_FORMAL_SERVER {
		serverState += "正式服"
	} else {
		serverState += "测试服"
	}
	logs.Info(serverState)
	Initialize()
	InitializeDependDB()
	//PlayerIdGen.ReadEverySiteDB(MongoMgr.GetMongoSessions())
	PSysParameterMgr = for_game.NewSysParameterManager()
	//信号监控
	ServeFunctions = append(ServeFunctions, SignHandle)
	//web api初始化
	PWebApiForClient = NewWebHttpForClient(PServerInfo.GetClientApiPort())
	PWebApiForServer = NewWebHttpForServer(PServerInfo.GetServerApiPort())
	ServeFunctions = append(ServeFunctions, PWebApiForClient.Serve)
	ServeFunctions = append(ServeFunctions, PWebApiForServer.Serve)
	jobs := []easygo.IGoroutine{}
	for _, f := range ServeFunctions {
		job := easygo.Spawn(f)
		jobs = append(jobs, job)
	}
	logs.Info(PServerInfo.GetName() + "启动成功")
	//启动etcd
	PClient3KVMgr.StartClintTV3()
	defer PClient3KVMgr.Close() //关闭etcd
	//etcd注册和发现
	for_game.InitExistServer(PClient3KVMgr, PServerInfoMgr, PServerInfo)
	//redis初始化
	easygo.JoinAllJobs(jobs...)

}

var ServeFunctions = []func(){}

func SignHandle() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM)
	for {
		s := <-c
		switch s {
		case syscall.SIGTERM:
			//TODO:服务器关闭逻辑处理
			IsStopServer = true
			logs.Info("登录服务器关闭:", PServerInfo.GetSid())
			PServerInfoMgr.DelServerInfo(PServerInfo.GetSid())
			time.Sleep(time.Second * 1)
			os.Exit(1)
		default:
			logs.Debug("error sign", s)
		}
	}
}
