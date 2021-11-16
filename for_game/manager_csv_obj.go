package for_game

import (
	"game_server/easygo"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"strings"
	"sync"
)

const SERVER_CONFIG_PATH = "./config/" //服务器配置文件路径
//管理Csv对象
type CsvObjManager struct {
	sync.Map
	Mutex easygo.RLock
}

func NewCsvObjManager() *CsvObjManager {
	p := &CsvObjManager{}
	p.InitCsvPath(SERVER_CONFIG_PATH)
	return p
}

//初始化/重载指定路径下csv数据到内存
func (self *CsvObjManager) InitCsvPath(path string) {
	dir_list, e := ioutil.ReadDir(path)
	if e != nil {
		logs.Error("加载文件目录出错")
		return
	}
	for _, v := range dir_list {
		if strings.Contains(v.Name(), ".csv") {
			//只处理csv文本
			fileName := strings.Split(v.Name(), ".csv")
			pathFile := path + v.Name()
			logs.Info("初始化csv数据:", pathFile)
			csvObj := NewCsvBase(pathFile)
			self.Store(fileName[0], csvObj)
		}
	}
}

//初始化/重载指定文件名csv数据到内存
func (self *CsvObjManager) InitCsvFile(fileName string) {
	pathFile := SERVER_CONFIG_PATH + fileName + ".csv"
	logs.Info("初始化csv数据:", pathFile)
	csvObj := NewCsvBase(pathFile)
	self.Store(fileName, csvObj)
}

func (self *CsvObjManager) LoadCsvBaseObj(fileName string) *CsvBaseObj {
	value, ok := self.Load(fileName)
	if ok {
		return value.(*CsvBaseObj)
	}
	return nil
}
func (self *CsvObjManager) CreateCsvBaseObj(fileName string) *CsvBaseObj {
	pathFile := SERVER_CONFIG_PATH + fileName + ".csv"
	csvObj := CreateCsvBase(pathFile)
	self.Store(fileName, csvObj)
	return csvObj
}

//对外方法，获取csvbase对象，
func (self *CsvObjManager) GetCsvBaseObj(fileName string) *CsvBaseObj {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadCsvBaseObj(fileName)
	return obj
}
func (self *CsvObjManager) GetCsvData(fileName, key string) *easygo.CsvData {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	obj := self.LoadCsvBaseObj(fileName)
	return obj.GetCsvData(key)
}
