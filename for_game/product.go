package for_game

import (
	"fmt"
	"game_server/easygo"
	"log"
)

var _ = fmt.Sprintf
var _ = log.Println

type IProduct interface {
	easygo.IFinalable

	GetPrimaryKey() PRIMARY_KEY
	InsertToDB(kwargs ...easygo.KWAT)
	SaveToDB(isSet ...bool) bool
	DeleteFromDB()
	LoadFromDB(kwargs ...easygo.KWAT) bool

	// SetFactory(factory IFactory)
	// GetFactory() IFactory

	// AddManager(manager IProductManager)
	// RemoveManager(manager IProductManager)
	// ManagerAmount() int
	GetPersistenceObj() IPersistence
	GetChineseName() string    // 中文名，用于调试
	GetDirtyData() interface{} //获取更新数据
	CleanDirtyData()           //重置更新数据
}

type Product struct {
	easygo.Finalable

	Me IProduct

	ProductMutex easygo.RLock

	PrimaryKey  PRIMARY_KEY
	ChineseName string
	// ProductManagers []IProductManager
	// Factory         IFactory
}

/* 抽象类，不能实例化
func NewProduct(primaryKey PRIMARY_KEY) *Product {
}*/

func (self *Product) Init(me IProduct, primaryKey PRIMARY_KEY, chineseName string) {
	self.Me = me
	self.ChineseName = chineseName
	self.Finalable.Init(me)
	self.PrimaryKey = primaryKey
}

func (self *Product) GetPrimaryKey() PRIMARY_KEY {
	return self.PrimaryKey
}

// func (self *Product) AddManager(manager IProductManager) {
// 	self.ProductMutex.Lock()
// 	defer self.ProductMutex.Unlock()

// 	self.ProductManagers = append(self.ProductManagers, manager)
// }

// func (self *Product) RemoveManager(manager IProductManager) {
// 	self.ProductMutex.Lock()
// 	defer self.ProductMutex.Unlock()

// 	i := easygo.Index(self.ProductManagers, manager)
// 	if i == -1 {
// 		panic("为什么是负一呢")
// 	}
// 	self.ProductManagers = append(self.ProductManagers[:i], self.ProductManagers[i+1:]...)
// }

// func (self *Product) ManagerAmount() int {
// 	self.ProductMutex.Lock()
// 	defer self.ProductMutex.Unlock()

// 	return len(self.ProductManagers)
// }

// func (self *Product) SetFactory(factory IFactory) {
// 	self.Factory = factory
// }
// func (self *Product) GetFactory() IFactory {
// 	return self.Factory
// }

func (self *Product) GetChineseName() string { // 中文名，用于调试
	return self.ChineseName
}

func (self *Product) InsertToDB(kwargs ...easygo.KWAT) { //, insertValues ...interface{}
	panic("抽象方法，请在子类实现")
}

func (self *Product) SaveToDB(isSet ...bool) bool {
	panic("抽象方法，请在子类实现")
}

func (self *Product) DeleteFromDB() {
	panic("抽象方法，请在子类实现")
}

func (self *Product) LoadFromDB(kwargs ...easygo.KWAT) bool {
	panic("抽象方法，请在子类实现")
}

func (self *Product) GetPersistenceObj() IPersistence {
	panic("抽象方法，请在子类实现")
}
func (self *Product) GetDirtyData() interface{} {
	panic("抽象方法，请在子类实现")
}
func (self *Product) CleanDirtyData() {
	panic("抽象方法，请在子类实现")
}
