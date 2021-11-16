package for_game

import (
	"encoding/csv"
	"game_server/easygo"
	"github.com/astaxie/beego/logs"
	"io"
	"log"
	"os"
)

/*
csv数据管理基类模板
*/

type CsvBaseObj struct {
	Data     map[string]*easygo.CsvData
	FilePath string
	Mutex    easygo.RLock
}

func NewCsvBase(file string) *CsvBaseObj {
	p := &CsvBaseObj{}
	p.Init(file)
	return p
}
func CreateCsvBase(file string) *CsvBaseObj {
	p := &CsvBaseObj{}
	p.Init(file, true)
	return p
}

func (self *CsvBaseObj) Init(file string, isCreate ...bool) {
	create := append(isCreate, false)[0]
	self.Data = make(map[string]*easygo.CsvData, 0)
	self.FilePath = file
	if !create {
		self.LoadFile()
	} else {
		self.CreateFile()
	}

}

//加载指定文件数据
func (self *CsvBaseObj) LoadFile() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	csvfile, err := os.Open(self.FilePath)
	if err != nil {
		logs.Error("Couldn't open the csv file", err)
	}
	defer csvfile.Close()
	// Parse the file
	r := csv.NewReader(csvfile)
	// Iterate through the records
	keys := make([]string, 0)
	cn := 0
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if cn == 0 {
			keys = record
		} else {
			data := easygo.NewCsvData()
			for i := 0; i < len(record); i++ {
				data.Add(keys[i], record[i])
			}
			self.Data[record[0]] = data
		}
		cn += 1
	}
	logs.Info("加载csv配置完毕:", self.FilePath)
}

//创建指定文件
func (self *CsvBaseObj) CreateFile() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	csvfile, err := os.OpenFile(self.FilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		logs.Error("Couldn't create the csv file", err)
	}
	defer csvfile.Close()
	logs.Info("创建模块")
}

//修改csv文件指定key的字段值:存在编码问题有待完善
//func (self *CsvBaseObj) ModifyFileData(key, name, val string) {
//	self.Mutex.Lock()
//	defer self.Mutex.Unlock()
//	csvfile, err := os.OpenFile(self.FilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
//	if err != nil {
//		logs.Error("Couldn't open the csv file", err)
//	}
//	defer csvfile.Close()
//	// Parse the file
//	r := csv.NewReader(csvfile)
//	// Iterate through the records
//	keys := make([]string, 0)
//	// Read all records from csv
//	records, err := r.ReadAll()
//	if err != nil {
//		logs.Error("配置文件读取出错:", self.FilePath)
//		return
//	}
//	if len(records) <= 0 {
//		return
//	}
//	keys = records[0]
//	for i := 1; i < len(records); i++ {
//		if records[i][0] == key {
//			for j := 0; j < len(keys); j++ {
//				if keys[j] == name {
//					records[i][j] = val
//					break
//				}
//			}
//			break
//		}
//	}
//	//重新写入csv
//	newFile, err := os.Create(self.FilePath)
//	defer newFile.Close()
//	w := csv.NewWriter(newFile)
//	w.WriteAll(records)
//	w.Flush()
//}

//增加csv文件指定key的字段值
//func (self *CsvBaseObj) AddFileData(data []string) {
//	self.Mutex.Lock()
//	defer self.Mutex.Unlock()
//	csvFile, err := os.OpenFile(self.FilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
//	if err != nil {
//		logs.Error("Couldn't open the csv file", err)
//	}
//	defer csvFile.Close()
//	// Parse the file
//	//重新写入csv
//	w := csv.NewWriter(csvFile)
//	logs.Info("file,data", data)
//	w.Write(data)
//	w.Flush()
//	err1 := w.Error()
//	if err1 != nil {
//		logs.Error("err:", err1.Error())
//	}
//}

//删除csv文件指定key的字段值
//func (self *CsvBaseObj) DeleteFileData(key string) {
//	self.Mutex.Lock()
//	defer self.Mutex.Unlock()
//	csvfile, err := os.OpenFile(self.FilePath, os.O_RDWR|os.O_SYNC, os.ModePerm)
//	if err != nil {
//		logs.Error("Couldn't open the csv file", err)
//	}
//	defer csvfile.Close()
//	// Parse the file
//	r := csv.NewReader(csvfile)
//	// Read all records from csv
//	records, err := r.ReadAll()
//	if err != nil {
//		logs.Error("配置文件读取出错:", self.FilePath)
//		return
//	}
//	if len(records) <= 0 {
//		return
//	}
//	newRecords := make([][]string, 0)
//	for i := 0; i < len(records); i++ {
//		logs.Info("%s-%s", records[i][0], key)
//		if records[i][0] == key {
//			continue
//		}
//		newRecords = append(newRecords, records[i])
//	}
//	logs.Info(" newRecords:", newRecords)
//	//重新写入csv
//	newFile, err := os.Create(self.FilePath)
//	defer newFile.Close()
//	w := csv.NewWriter(newFile)
//	w.WriteAll(newRecords)
//	w.Flush()
//}

//获取指定数据
func (self *CsvBaseObj) GetCsvData(key string) *easygo.CsvData {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return self.Data[key]
}
