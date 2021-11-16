package for_game

import (
	"game_server/easygo"
	"log"
	"time"
)

var _ = log.Println

var _StandardTime time.Time

func init() {
	loc, err := time.LoadLocation("Local")
	easygo.PanicError(err)
	t, err2 := time.ParseInLocation("2006-01-02 15:04:05", "2020-04-01 00:00:00", loc)
	easygo.PanicError(err2)

	_StandardTime = t
}

type IPeriodSeqFunctionSet interface {
	GetMinuteNo() int
	GetHourNo() int
	GetDayNo() int
	GetWeekNo() int
	GetMonthNo() int
	GetHalfYearNo() int
}
type PeriodSeqFunctionSet struct {
	Me IPeriodSeqFunctionSet
}

func NewPeriodSeqFunctionSet() *PeriodSeqFunctionSet {
	p := &PeriodSeqFunctionSet{}
	p.Init(p)
	return p
}

func (self *PeriodSeqFunctionSet) Init(me IPeriodSeqFunctionSet) {
	self.Me = me
}

func (self *PeriodSeqFunctionSet) GetMinuteNo() int {
	t := time.Now()
	t1 := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
	t2 := time.Date(_StandardTime.Year(), _StandardTime.Month(), _StandardTime.Day(), _StandardTime.Hour(), _StandardTime.Minute(), 0, 0, time.Local)
	return int(t1.Sub(t2).Minutes())
}

func (self *PeriodSeqFunctionSet) GetHourNo() int {
	t := time.Now()
	t1 := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.Local)
	t2 := time.Date(_StandardTime.Year(), _StandardTime.Month(), _StandardTime.Day(), _StandardTime.Hour(), 0, 0, 0, time.Local)
	return int(t1.Sub(t2).Hours())
}

func (self *PeriodSeqFunctionSet) GetDayNo() int {
	t := time.Now()
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	t2 := time.Date(_StandardTime.Year(), _StandardTime.Month(), _StandardTime.Day(), 0, 0, 0, 0, time.Local)
	return int(t1.Sub(t2).Hours() / 24)
}

func (self *PeriodSeqFunctionSet) GetWeekNo() int {
	t := time.Now()
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	t2 := time.Date(_StandardTime.Year(), _StandardTime.Month(), _StandardTime.Day(), 0, 0, 0, 0, time.Local)
	return int(t1.Sub(t2).Hours() / 24 / 7)
}

func (self *PeriodSeqFunctionSet) GetMonthNo() int {
	t := time.Now()
	t1 := time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, time.Local)
	t2 := time.Date(_StandardTime.Year(), _StandardTime.Month(), 0, 0, 0, 0, 0, time.Local)
	return int(t1.Sub(t2).Hours() / 24 / 30)
}

func (self *PeriodSeqFunctionSet) GetHalfYearNo() int {
	t := time.Now()
	t1 := time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, time.Local)
	t2 := time.Date(_StandardTime.Year(), _StandardTime.Month(), 0, 0, 0, 0, 0, time.Local)
	return int(t1.Sub(t2).Hours() / 24 / 30 / 6)
}

var PeriodSeq IPeriodSeqFunctionSet
