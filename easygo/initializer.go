package easygo

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/astaxie/beego/logs"
	"github.com/sasha-s/go-deadlock"
)

type IInitializer interface {
	Execute(dict KWAT)
	MiscInit()
	SetDeadLockOptions()
	SetBeegoLogs(dict KWAT)
	GetBeeLogger() *logs.BeeLogger
	CreateFunctionSet() IFunctionSet
	CreateYamlConfig(dict KWAT) IYamlConfig
	CreateMongoMgr() IMongoDBManager
	CreateRedisMgr(yaml IYamlConfig) IRedisManager
}

type Initializer struct {
	Me IInitializer
	// 此类不要带状态
}

func NewInitializer() *Initializer {
	p := &Initializer{}
	p.Init(p)
	return p
}

func (self *Initializer) Init(me IInitializer) {
	self.Me = me
}

func (self *Initializer) Execute(dict KWAT) {
	self.Me.MiscInit()
	self.Me.SetBeegoLogs(dict)
	self.Me.SetDeadLockOptions()
	FuncSet = self.Me.CreateFunctionSet()

	YamlCfg = self.Me.CreateYamlConfig(dict)

	//master数据库
	MongoMgr = self.Me.CreateMongoMgr()
	admin := "admin"                                                                  // MongoDB 自带的库
	MongoMgr.AddDbSession(admin, YamlCfg.GetSpecificInfoOrDefault(admin, "ningmeng")) // 这个 ningmeng 先写死吧，后面优化
	MongoMgr.InitSomeDB(YamlCfg.GetMongoDBDsnMaster())
	//slave数据库
	MongoLogMgr = self.Me.CreateMongoMgr()                                                        // MongoDB 自带的库
	MongoLogMgr.AddDbSession(admin, YamlCfg.GetSpecificInfoOrDefaultSlave(admin, "ningmeng_log")) // 这个 ningmeng_log 先写死吧，后面优化
	MongoLogMgr.InitSomeDB(YamlCfg.GetMongoDBDsnSlave())
	//redis数据库
	RedisMgr = self.Me.CreateRedisMgr(YamlCfg)

}

func (self *Initializer) CreateFunctionSet() IFunctionSet {
	return NewFunctionSet()
}

func (self *Initializer) CreateYamlConfig(dict KWAT) IYamlConfig {
	path, ok := dict["yamlPath"]
	if !ok {
		panic("我需要一个 yamlPath 参数")
	}
	return NewYamlConfig(path.(string))
}

func (self *Initializer) CreateMongoMgr() IMongoDBManager {
	return NewMongoDBManager()
}
func (self *Initializer) CreateRedisMgr(yaml IYamlConfig) IRedisManager {
	return NewRedisManager(yaml)
}

// 杂 7 杂 8 初始化
func (self *Initializer) MiscInit() {
	log.SetOutput(os.Stdout)           // log.PrintXXX 默认输出到 stderr ,这样不好，改到 stdout
	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
	runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪
}

func (self *Initializer) SetDeadLockOptions() {
	// deadlock.Opts.Disable = true // 死锁检查 + 超时检查
	// deadlock.Opts.DisableLockOrderDetection = true // 死锁检查
	deadlock.Opts.DeadlockTimeout = 0             // 默认 30 秒,0 表示不检查 // 4 * time.Second
	deadlock.Opts.OnPotentialDeadlock = func() {} // 默认 exit()
}

func (self *Initializer) SetBeegoLogs(dict KWAT) {
	logName, ok := dict["logName"]
	if !ok {
		panic("我需要一个 logName 参数")
	}

	var err error
	err = logs.SetLogger(logs.AdapterConsole, fmt.Sprintf(`{"level":%d}`, logs.LevelDebug))
	PanicError(err)

	config := `{"filename":"logs/%s.log", "separate":["error", "warning", "info", "debug"], "level":%d, "daily":true, "rotate":true, "maxdays":15,"perm":"777"}`
	config = fmt.Sprintf(config, logName, logs.LevelDebug)
	err = logs.SetLogger(logs.AdapterMultiFile, config)
	PanicError(err)

	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
	logs.Async()
}

func (self *Initializer) GetBeeLogger() *logs.BeeLogger {
	// TODO: 不要使用 logs 的包变量 beeLogger，要 new 一个新的 loger 出来
	return logs.GetBeeLogger()
}

var (
	FuncSet     IFunctionSet
	YamlCfg     IYamlConfig
	MongoMgr    IMongoDBManager //主游戏DB管理器
	MongoLogMgr IMongoDBManager //游戏日志管理器
	RedisMgr    IRedisManager   //redis存储
)
