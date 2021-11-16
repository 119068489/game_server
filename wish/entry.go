package wish

import (
	"flag"
	"game_server/easygo"
	"game_server/for_game"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/astaxie/beego/logs"
)

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
		"yamlPath": "config_wish.yaml",
	}
	initializer.Execute(dict)

	serverState := "当前服务器状态:"
	if for_game.IS_FORMAL_SERVER {
		serverState += "正式服"
	} else {
		serverState += "测试服"
	}
	log.Println(serverState)

	Initialize()
	InitializeDependDB()
	//启动etcd
	PClient3KVMgr.StartClintTV3()
	defer PClient3KVMgr.Close() //关闭etcd

	for_game.TimerMgr = for_game.NewSquareTimerMgr()

	PPayChannelMgr = for_game.NewPayChannelManager()
	PSysParameterMgr = for_game.NewSysParameterManager()

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
			//TODO:服务器关闭逻辑处理
			logs.Info("许愿池服务器关闭", PServerInfo.GetSid())
			logs.Info("停服更新许愿池数据开始")
			SavePoolDataToMongo() // 停服更新许愿池数据
			logs.Info("停服更新许愿池数据结束")
			PServerInfoMgr.DelServerInfo(PServerInfo.GetSid())
			time.Sleep(time.Second * 1)
			os.Exit(1)
		default:
			logs.Debug("error sign", s)
		}
	}
}
