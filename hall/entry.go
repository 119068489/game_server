// 大厅服务
package hall

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
		"yamlPath": "config_hall.yaml",
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
	InitCoinfigToRedis()
	//支付渠道管理初始化
	PPayChannelMgr = for_game.NewPayChannelManager()
	//初始化系统参数设置
	PSysParameterMgr = for_game.NewSysParameterManager()
	//启动cron
	StartCron()

	//启动http web服务器
	WebServreMgr = NewWebHttpServer()
	ServeFunctions = append(ServeFunctions, WebServreMgr.Serve)
	//web api初始化
	PWebApiForClient = NewWebHttpForClient(PServerInfo.GetClientApiPort())
	PWebApiForServer = NewWebHttpForServer(PServerInfo.GetServerApiPort())
	ServeFunctions = append(ServeFunctions, PWebApiForClient.Serve)
	ServeFunctions = append(ServeFunctions, PWebApiForServer.Serve)
	//信号监控
	ServeFunctions = append(ServeFunctions, SignHandle)

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
			IsStopServer = true
			logs.Info("大厅服务器关闭:", PServerInfo)
			SaveRedisData()
			time.Sleep(time.Second * 5)
			os.Exit(1)
		default:
			logs.Debug("error sign", s)
		}
	}
}
