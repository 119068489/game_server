package for_game

import (
	"game_server/easygo"
	"sync"

	mgo "github.com/akqp2019/mgo"
)

type ProductInfo struct {
	CreateFun func(string) IProduct // 创建实例的方法
	Mgr       *sync.Map             // 存储实例的管理器
}

type IInitializer interface {
	easygo.IInitializer

	// 这些虚函数是为了给机会其他模块定义不同的实例
}

type Initializer struct {
	easygo.Initializer
	Me          IInitializer
	ProductInfo []*ProductInfo
}

func NewInitializer() *Initializer {
	p := &Initializer{}
	p.Init(p)
	return p
}

func (self *Initializer) Init(me IInitializer) {
	self.Me = me
	self.Initializer.Init(me)
}

func (self *Initializer) Execute(dict easygo.KWAT) {
	self.Initializer.Execute(dict)
	IS_FORMAL_SERVER = easygo.YamlCfg.GetValueAsBool("IS_FORMAL_SERVER")
	SERVER_CENTER_ADDR = easygo.YamlCfg.GetValueAsString("SERVER_CENTER_ADDR")

	IS_TFSERVER = easygo.YamlCfg.GetValueAsBool("IS_TFSERVER")
	MessageMarkInfo = NewMessageMarkInfoMgr()
	PeriodSeq = self.CreatePeriodSequence()
	//PDirtyWordsMgr = NewDirtyWordMgr()
}

// override
func (self *Initializer) MiscInit() {
	self.Initializer.MiscInit()

	mgo.SetStats(true) // 开启 mgo 各种对象计数。不然后面调用 GetStats() 会指针异常
}

func (self *Initializer) ReCreateConfig() {
}

//增加站点配置 {不要新增了，直接 ReCreate 来得快}
// func (self *Initializer) AddSiteMgrObject(siteName string) {
// 	for _, info := range self.ProductInfo {
// 		creator, mgr := info.CreateFun, info.Mgr
// 		config := creator(siteName)
// 		if !config.LoadFromDB() {
// 			config.InsertToDB()
// 		}
// 		mgr.Store(siteName, config)
// 	}
// }
//func (self *Initializer) CreateRedisPool(address string, db int) IRedisPool {
//	return NewRedisPool(address, db)
//}

// override
// func (self *Initializer) CreateYamlConfig(dict easygo.KWAT) easygo.IYamlConfig {

// 	path, ok := dict["yamlPath"]
// 	if !ok {
// 		panic("我需要一个 yamlPath 参数")
// 	}

// 	return NewYamlConfig(path.(string))
// }

func (self *Initializer) CreatePeriodSequence() IPeriodSeqFunctionSet {
	return NewPeriodSeqFunctionSet()
}

// func (self *Initializer) CreateSiteStyleList() *StyleList {
// 	return NewStyleList()
// }
