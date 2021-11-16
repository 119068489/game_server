package for_game

import (
	"game_server/easygo"

	"github.com/akqp2019/mgo"
)

//=======================================

type PeriodDataBlock struct {
	MongoProduct `bson:"-"`
	Persistence  `bson:"-"`

	PlayerId int64 `bson:"-"` // 玩家 id

	DayPeriod      *DayPeriodData      `bson:"d,omitempty"`  // 天变量
	WeekPeriod     *WeekPeriodData     `bson:"w,omitempty"`  // 周变量
	MonthPeriod    *MonthPeriodData    `bson:"m,omitempty"`  // 月变量
	HaltYearPeriod *HalfYearPeriodData `bson:"hy,omitempty"` // 半年变量
	DirtyData      *PeriodDataBlock    `bson:"-"`
}

//抽象  子类实现
//func NewPeriodDataBlock(playerId PLAYER_ID, site string) *PeriodDataBlock {
//
//}
func NewPeriodDataBlock(playerId PLAYER_ID) *PeriodDataBlock {
	p := &PeriodDataBlock{}
	p.Init(playerId)
	return p
}

func (self *PeriodDataBlock) Init(playerId PLAYER_ID) {
	self.MongoProduct.Init(self, playerId, "周期数据")
	kwargs1 := easygo.KWAT{
		"DirtyEventHandler": self.DirtyEventHandler,
	}
	self.Persistence.Init(self, kwargs1)

	self.PlayerId = playerId
	keepPeriod := 1
	kwargs2 := easygo.KWAT{
		"DirtyEventHandler": self.DirtyEventHandler,
		"Locker":            self.Mutex,
	}
	self.DayPeriod = NewDayPeriodData(keepPeriod, kwargs2)
	//self.WeekPeriod = NewWeekPeriodData(keepPeriod, kwargs2)  //周变量周计算有问题 弃用
	self.MonthPeriod = NewMonthPeriodData(keepPeriod, kwargs2)
	self.HaltYearPeriod = NewHaltYearPeriodData(keepPeriod, kwargs2)
}

func (self *PeriodDataBlock) DirtyEventHandler(isAll ...bool) {
	self.SaveToDB(isAll...)

}
func (self *PeriodDataBlock) GetPersistenceObj() IPersistence { // override
	return self
}
func (self *PeriodDataBlock) GetDirtyData() interface{} { //override
	return self.DirtyData
}
func (self *PeriodDataBlock) CleanDirtyData() { //override
	self.DirtyData = &PeriodDataBlock{}
}
func (self *PeriodDataBlock) GetC() (c *mgo.Collection, fun func()) { // override
	return easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_PERIOD)
}
func (self *PeriodDataBlock) GetPlayerId() int64 {
	return self.PlayerId
}

// 对外接口
func GetPlayerPeriod(pid PLAYER_ID, kwargs ...easygo.KWAT) *PeriodDataBlock {
	obj := NewPeriodDataBlock(pid)
	if !obj.LoadFromDB(kwargs...) {
		obj.InsertToDB(kwargs...)
	}
	return obj
}
