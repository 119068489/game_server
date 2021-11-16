package easygo

import (
	"time"

	"github.com/robfig/cron"
)

type CronJob struct {
	Job                *cron.Cron
	TenMinEvent        *Event
	HalfHourEvent      *Event
	HourEvent          *Event
	DayEvent           *Event
	WeekEvent          *Event // 每周一五点执行
	MonthEvent         *Event // 每月1号五点执行
	DayEightClockEvent *Event // 每天八点执行
	DayFiveClockEvent  *Event // 每天五点执行
}

func NewCronJob() *CronJob {
	p := &CronJob{}
	p.Init()
	return p
}

func GetTimeData() (year int, month int, day int, hour int, Minute int, week int) {
	t := time.Now()
	return t.Year(), (int)(t.Month()), t.Day(), t.Hour(), t.Minute(), (int)(t.Weekday())
}
func newWithSecond() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}
func (self *CronJob) Init() {
	self.Job = newWithSecond()
	self.TenMinEvent = NewEvent()
	self.HalfHourEvent = NewEvent()
	self.HourEvent = NewEvent()
	self.DayEvent = NewEvent()
	self.WeekEvent = NewEvent()
	self.MonthEvent = NewEvent()
	self.DayEightClockEvent = NewEvent()
	self.DayFiveClockEvent = NewEvent()
	//秒 分 时 日 月 周

	f1 := func() {
		self.HalfHourEvent.Trigger(GetTimeData())
	}
	self.Job.AddFunc("0 */30 * * * *", f1) //整30分钟触发

	f2 := func() { self.HourEvent.Trigger(GetTimeData()) }
	self.Job.AddFunc("0 0 * * * *", f2) //整点触发

	f3 := func() { self.DayEvent.Trigger(GetTimeData()) }
	self.Job.AddFunc("0 0 0 * * *", f3) //整天0点触发

	f4 := func() { self.WeekEvent.Trigger(GetTimeData()) }
	self.Job.AddFunc("0 0 5 * * 1", f4)

	f5 := func() { self.MonthEvent.Trigger(GetTimeData()) }
	self.Job.AddFunc("0 0 5 1 * *", f5)

	f6 := func() {
		self.TenMinEvent.Trigger(GetTimeData())
	}
	self.Job.AddFunc("0 */10 * * * *", f6) //整10分钟触发

	f7 := func() {
		self.DayEightClockEvent.Trigger(GetTimeData())
	}
	self.Job.AddFunc("0 0 8 * * *", f7) // 每天八点执行
	f8 := func() {
		self.DayFiveClockEvent.Trigger(GetTimeData())
	}
	self.Job.AddFunc("0 0 5 * * *", f8) // 每天五点执行
}

func (self *CronJob) Serve() {
	self.InitJob()           //起服挂载上所有方法
	Spawn(self.StartCronJob) //携程起定时器
}
func (self *CronJob) StartCronJob() {
	self.Job.Start()
	defer self.Job.Stop() //关闭计划任务, 但是不能关闭已经在执行中的任务.
	select {}
}

func (self *CronJob) InitJob() {

}

var Cronjob *CronJob //定时器任务
func init() {
	Cronjob = NewCronJob() //定时器任务
	Cronjob.Serve()
}

func TenMinEventAddHandler(function interface{}) {
	Cronjob.TenMinEvent.AddHandler(function)
}

func HalfHourEventAddHandler(function interface{}) {
	Cronjob.HalfHourEvent.AddHandler(function)
}

func HourEventAddHandler(function interface{}) {
	Cronjob.HourEvent.AddHandler(function)
}

func DayEventAddHandler(function interface{}) {
	Cronjob.DayEvent.AddHandler(function)
}
func WeekEventAddHandler(function interface{}) {
	Cronjob.WeekEvent.AddHandler(function)
}

func MonthEventAddHandler(function interface{}) {
	Cronjob.MonthEvent.AddHandler(function)
}

func DayEightClockEventAddHandler(function interface{}) {
	Cronjob.DayEightClockEvent.AddHandler(function)
}

func DayFiveClockEventAddHandler(function interface{}) {
	Cronjob.DayEightClockEvent.AddHandler(function)
}
