package easygo

import (
	"time"
)

const A_DAY_SECOND = 24 * 3600
const TIME_ZONE_OFFSET = 8 * 3600 // 时间戳为 0 时的北京时间是 1970/1/1 8:0:0
//指定时间0点的时间戳
func Get0ClockTimestamp(timeNew int64) int64 {
	if timeNew > 1000000000000 {
		timeNew = timeNew / 1000
	}
	if timeNew < 0 {
		return 0
	}
	return timeNew - (timeNew+TIME_ZONE_OFFSET)%A_DAY_SECOND // 已经测试过了

	//timeStr := time.Unix(timeNew, 0).Format("2006-01-02")
	//loc, err := time.LoadLocation("Local")
	//PanicError(err)
	//t, err2 := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 00:00:00", loc)
	//PanicError(err2)
	//timeNumber := t.Unix()
	//
	//return timeNumber
}

//指定时间0点的时间戳
func Get24ClockTimestamp(timeNew int64) int64 {
	if timeNew > 1000000000000 {
		timeNew = timeNew / 1000
	}
	return timeNew - (timeNew+TIME_ZONE_OFFSET)%A_DAY_SECOND + A_DAY_SECOND
}

//指定时间0点的时间戳返回毫秒
func Get0ClockMillTimestamp(timeNew int64) int64 {
	return Get0ClockTimestamp(timeNew) * 1000
}

//指定时间0点的时间戳返回毫秒
func Get24ClockMillTimestamp(timeNew int64) int64 {
	return Get24ClockTimestamp(timeNew) * 1000
}

//今天0点时的时间戳
func GetToday0ClockTimestamp() int64 {
	return Get0ClockTimestamp(time.Now().Unix())
}

//今天24点时的时间戳
func GetToday24ClockTimestamp() int64 {
	return GetToday0ClockTimestamp() + (24 * 3600)
}

//昨天0点时的时间戳
func GetYesterday0ClockTimestamp() int64 {
	return GetYesterday24ClockTimestamp() - A_DAY_SECOND
}

//昨天24点时的时间戳
func GetYesterday24ClockTimestamp() int64 {
	return GetToday0ClockTimestamp()
}

//明天0点时的时间戳
func GetTomorrow0ClockTimestamp() int64 {
	return GetToday24ClockTimestamp()
}

//是否是今天的时间
func IsTodayTimestamp(timestamp int64) bool {
	if timestamp > 1000000000000 {
		timestamp = timestamp / 1000
	}
	return GetToday0ClockTimestamp() <= timestamp && timestamp < GetTomorrow0ClockTimestamp()
}

func NowTimestamp() int64 {
	return time.Now().Unix()
}

//获取本周开始时间戳
func GetWeek0ClockTimestamp() int64 {
	now := time.Now()

	offset := int64(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	first := time.Now().Unix() + (86400 * offset)
	first -= (first + TIME_ZONE_OFFSET) % A_DAY_SECOND
	return first
}

//获取本周12点的时间戳
func GetWeek12ClockTimestamp() int64 {
	now := time.Now()

	offset := int64(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	first := time.Now().Unix() + (86400 * offset)
	first -= (first + TIME_ZONE_OFFSET) % A_DAY_SECOND
	first = first + (A_DAY_SECOND / 2)
	return first
}

//指定时间获取当周开始时间戳
func GetWeek0ClockOfTimestamp(t int64) int64 {
	if t > 1000000000000 {
		t = t / 1000
	}

	dt := Stamp2Time(t)

	offset := int64(time.Monday - dt.Weekday())
	if offset > 0 {
		offset = -6
	}

	first := t + (86400 * offset)
	first -= (first + TIME_ZONE_OFFSET) % A_DAY_SECOND
	return first
}

//获取本月开始时间戳
func GetMonth0ClockTimestamp() int64 {
	now := time.Now()
	now = now.AddDate(0, 0, -now.Day()+1)
	nowTime := now.Unix()
	nowTime -= (nowTime + TIME_ZONE_OFFSET) % A_DAY_SECOND
	return nowTime
}

//获取本月12点时间戳
func GetMonth12ClockTimestamp() int64 {
	now := time.Now()
	now = now.AddDate(0, 0, -now.Day()+1)
	nowTime := now.Unix()
	nowTime -= (nowTime + TIME_ZONE_OFFSET) % A_DAY_SECOND
	nowTime = nowTime + (A_DAY_SECOND / 2)
	return nowTime
}

//指定时间获取当月开始时间戳
func GetMonth0ClockOfTimestamp(t int64) int64 {
	if t > 1000000000000 {
		t = t / 1000
	}

	dt := Stamp2Time(t)

	dt = dt.AddDate(0, 0, -dt.Day()+1)
	dtTime := dt.Unix()
	dtTime -= (dtTime + TIME_ZONE_OFFSET) % A_DAY_SECOND
	return dtTime
}

/**字符串->时间戳 formatTimeStr 时间字符串"2006-01-02 15:04:05", isMs 是否毫秒 */
func GetTimeStrToTimestamp(formatTimeStr string, isMs ...bool) int64 {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, formatTimeStr, loc) //使用模板在对应时区转化为time.time类型
	isMs = append(isMs, false)
	if isMs[0] {
		return theTime.Unix() * 1000
	}
	return theTime.Unix()
}

/**字符串->时间对象*/
func Str2Time(formatTimeStr string) time.Time {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, formatTimeStr, loc) //使用模板在对应时区转化为time.time类型

	return theTime

}

/*时间戳->字符串*/
func Stamp2Str(stamp int64) string {
	if stamp > 1000000000000 {
		stamp = stamp / 1000
	}
	timeLayout := "2006-01-02 15:04:05"
	str := time.Unix(stamp, 0).Format(timeLayout)
	return str
}

/*时间戳->字符串:xxxx年xx月xx日*/
func Stamp2StrExt(stamp int64) string {
	if stamp > 1000000000000 {
		stamp = stamp / 1000
	}
	timeLayout := "2006年01月02日"
	str := time.Unix(stamp, 0).Format(timeLayout)
	return str
}

/*时间戳->时间对象*/
func Stamp2Time(stamp int64) time.Time {
	stampStr := Stamp2Str(stamp)
	timer := Str2Time(stampStr)
	return timer
}

// 时间定时器封装
type Timer struct {
	*time.Timer
	StartTimestamp int64         // 启动时间
	PauseTimestamp int64         // 暂停时间
	D              time.Duration // 延迟时间
}

// 暂停
func (self *Timer) Pause() {
	self.PauseTimestamp = time.Now().Unix()
	self.Stop()
}

// 恢复执行
func (self *Timer) Resume() {
	if self.PauseTimestamp > self.StartTimestamp {
		dis := self.D - time.Duration(self.PauseTimestamp-self.StartTimestamp)*time.Second
		self.Reset(dis)
		self.PauseTimestamp = 0
		self.StartTimestamp = time.Now().Unix()
		self.D = dis

	} else {
		panic("请先使用 Pause 暂停")
	}
}

// 定时。对外接口
func AfterFunc(d time.Duration, f func()) *Timer {
	newF := func() {
		defer RecoverAndLog()
		f()
	}
	timer := time.AfterFunc(d, newF)
	return &Timer{
		Timer:          timer,
		StartTimestamp: time.Now().Unix(),
		D:              d,
	}
}

//获取上月开始结束时间戳
func GetUpMouthStartEnd() (int64, int64) {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start := thisMonth.AddDate(0, -1, 0).Unix()
	end := thisMonth.AddDate(0, 0, 0).Unix()

	return start, end - 1
}

//获取本月开始结束时间戳
func GetMouthStartEnd(year int, month time.Month) (int64, int64) {
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start := thisMonth.AddDate(0, 0, 0).Unix()
	end := thisMonth.AddDate(0, 1, 0).Unix()

	return start, end - 1
}

//获取综合报表生成时间列表
func GetMakeReport(startTime int64) []int64 {
	sTime := time.Unix(startTime, 0)
	eTime := time.Now()
	d := sTime.Sub(eTime)
	f := d.Hours() / 24
	days := int(-f)

	var list []int64
	for i := 0; i <= days; i++ {
		if i > 0 {
			startTime = startTime + 86400
		}

		list = append(list, startTime)
	}

	return list
}

//获取2个时间戳之间的相差天数
func GetDifferenceDay(one, tow int64) int32 {
	if one > 1000000000000 {
		one = one / 1000
	}

	if tow > 1000000000000 {
		tow = tow / 1000
	}

	one0 := Get0ClockTimestamp(one)
	tow0 := Get0ClockTimestamp(tow)
	days := int64(tow0-one0) / A_DAY_SECOND
	return int32(days)
}

//检测是否是同一天
func CheckTheSameDay(one, tow int64) bool {
	if one > 1000000000000 {
		one = one / 1000
	}

	if tow > 1000000000000 {
		tow = tow / 1000
	}

	one0 := Get0ClockTimestamp(one)
	tow0 := Get0ClockTimestamp(tow)
	return tow0 == one0
}
