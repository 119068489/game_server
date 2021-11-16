package for_game

import (
	"fmt"
	"game_server/easygo"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"log"
)

var _ = fmt.Sprintf
var _ = log.Println

type IMongoProduct interface {
	IProduct
	GetC() (c *mgo.Collection, fun func())
}

type MongoProduct struct {
	Product
	Me IMongoProduct

	GetColFunc GET_C
}

/* 抽象类，不能实例化
func NewMongoProduct(primaryKey PRIMARY_KEY) *MongoProduct {
}*/

func (self *MongoProduct) Init(me IMongoProduct, primaryKey PRIMARY_KEY, name string, getColFunc ...GET_C) {
	self.Me = me
	if len(getColFunc) > 0 {
		self.GetColFunc = getColFunc[0]
	}
	self.Product.Init(me, primaryKey, name)
	//重置初始化更新数据
	self.Me.CleanDirtyData()
}

func (self *MongoProduct) GetC() (c *mgo.Collection, fun func()) {
	if self.GetColFunc == nil {
		panic("请在子类实现 GetC 函数或者传参进来")
	}
	return self.GetColFunc()
}

func (self *MongoProduct) InsertToDB(kwargs ...easygo.KWAT) {
	c, f := self.Me.GetC()
	defer f()

	obj := self.Me.GetPersistenceObj()
	obj.OnBorn(kwargs...)
	// e := c.Insert(obj)
	info, e := c.Upsert(bson.M{"_id": self.Me.GetPrimaryKey()}, obj)
	easygo.PanicError(e)
	_ = info
}

func (self *MongoProduct) DeleteFromDB() { // implement
	c, f := self.Me.GetC()
	defer f()

	e := c.RemoveId(self.Me.GetPrimaryKey())
	easygo.PanicError(e)
}

func (self *MongoProduct) SaveToDB(isSet ...bool) bool { // implement
	isAll := append(isSet, false)[0]
	obj := self.Me.GetPersistenceObj()
	obj.ClearDirtyFlags()
	c, f := self.Me.GetC()
	defer f()
	mutex := obj.GetLocker()
	mutex.Lock() // 下面的 UpdateId 是有访问 obj 的全部成员变量的，所以要上锁(已经发生过在 UpdateId 导致程序退出)
	defer mutex.Unlock()
	if isAll {
		dirtyData := self.Me.GetDirtyData() //获取更新数据
		self.Me.CleanDirtyData()            // 重置更新数据
		data := bson.M{"$set": dirtyData}
		e := c.UpdateId(self.Me.GetPrimaryKey(), data) // 使用set，更新指定字段值
		easygo.PanicError(e)
	} else {
		e := c.UpdateId(self.Me.GetPrimaryKey(), obj) // 不用 $set, 全部覆盖
		easygo.PanicError(e)
	}
	return true
}
func (self *MongoProduct) LoadFromDB(kwargs ...easygo.KWAT) bool { // implement
	c, f := self.Me.GetC()
	defer f()
	obj := self.Me.GetPersistenceObj()
	mutex := obj.GetLocker()
	mutex.Lock() // 下面的 UpdateId 是有访问 obj 的全部成员变量的，所以要上锁
	defer mutex.Unlock()
	e := c.Find(bson.M{"_id": self.Me.GetPrimaryKey()}).One(obj)
	if e != nil && e != mgo.ErrNotFound {
		panic(e)
	}
	if e == nil {
		obj.OnLoad()
	}
	return e == nil
}
