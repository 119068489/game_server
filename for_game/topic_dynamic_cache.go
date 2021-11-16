package for_game

import (
	"game_server/easygo"
	"game_server/easygo/util"
	"github.com/astaxie/beego/logs"
	"sync"
)

//话题最新动态列表缓存结构
//topicDynamicList map[话题ID][]动态ID
//topicDynamicCount map[话题ID]话题动态数总数量
//TopicDynamicTopCount map[话题ID]话题置顶动态数
//TopicDynamicUnTopCount map[话题ID]话题非置顶动态数
//TopicDynamicCacheTime map[话题ID]缓存时间戳（秒）
//CheckStatus map[话题ID]是否查询中 true 是
//CheckTime map[话题ID]查询操作的开始时间
//DefCacheNumber 默认的缓存数据量
//CheckTimeOffset 查询操作的过期时间（秒）
type topicNewDynamicCacheStructure struct {
	TopicDynamicList       map[int64][]int64
	TopicDynamicCount      map[int64]int
	TopicDynamicTopCount   map[int64]int
	TopicDynamicUnTopCount map[int64]int
	TopicDynamicCacheTime  map[int64]int64
	CheckStatus            map[int64]bool
	CheckTime              map[int64]int64
	LockChan               map[int64]chan struct{}
	Mu                     sync.Mutex
	DefCacheNumber         int
	CheckTimeOffset        int64
}

var topicNewDynamicCacheGlobal = &topicNewDynamicCacheStructure{
	TopicDynamicList:       make(map[int64][]int64),
	TopicDynamicCount:      make(map[int64]int),
	TopicDynamicTopCount:   make(map[int64]int),
	TopicDynamicUnTopCount: make(map[int64]int),
	TopicDynamicCacheTime:  make(map[int64]int64),
	CheckStatus:            make(map[int64]bool),
	CheckTime:              make(map[int64]int64),
	LockChan:               make(map[int64]chan struct{}),
	DefCacheNumber:         1000,
	CheckTimeOffset:        60,
}

var topicHotDynamicCacheGlobal = &topicNewDynamicCacheStructure{
	TopicDynamicList:       make(map[int64][]int64),
	TopicDynamicCount:      make(map[int64]int),
	TopicDynamicTopCount:   make(map[int64]int),
	TopicDynamicUnTopCount: make(map[int64]int),
	TopicDynamicCacheTime:  make(map[int64]int64),
	CheckStatus:            make(map[int64]bool),
	CheckTime:              make(map[int64]int64),
	LockChan:               make(map[int64]chan struct{}),
	DefCacheNumber:         1000,
	CheckTimeOffset:        60,
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) GetDefCacheNumber() int {
	return topicNewDynamicCache.DefCacheNumber
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) GetCheckTimeOffset() int64 {
	return topicNewDynamicCache.CheckTimeOffset
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) CreateLockChan(topicID int64) {
	topicNewDynamicCache.Mu.Lock()
	if _, ok := topicNewDynamicCache.LockChan[topicID]; !ok {
		topicNewDynamicCache.LockChan[topicID] = make(chan struct{}, 1)
	}
	topicNewDynamicCache.Mu.Unlock()
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) SetTopicDynamicList(topicID int64, topicNewDynamicIds []int64) {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	topicNewDynamicCache.TopicDynamicList[topicID] = topicNewDynamicIds
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) GetTopicDynamicList(topicID int64) []int64 {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	return topicNewDynamicCache.TopicDynamicList[topicID]
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) SetTopicDynamicCount(topicID int64, count int) {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	topicNewDynamicCache.TopicDynamicCount[topicID] = count
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) GetTopicDynamicCount(topicID int64) int {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	return topicNewDynamicCache.TopicDynamicCount[topicID]
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) SetTopicDynamicTopCount(topicID int64, count int) {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	topicNewDynamicCache.TopicDynamicTopCount[topicID] = count
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) GetTopicDynamicTopCount(topicID int64) int {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	return topicNewDynamicCache.TopicDynamicTopCount[topicID]
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) SetTopicDynamicUnTopCount(topicID int64, count int) {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	topicNewDynamicCache.TopicDynamicUnTopCount[topicID] = count
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) GetTopicDynamicUnTopCount(topicID int64) int {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	return topicNewDynamicCache.TopicDynamicUnTopCount[topicID]
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) SetTopicDynamicCacheTime(topicID int64, time int64) {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	topicNewDynamicCache.TopicDynamicCacheTime[topicID] = time
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) GetTopicDynamicCacheTime(topicID int64) int64 {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	return topicNewDynamicCache.TopicDynamicCacheTime[topicID]
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) SetCheckStatus(topicID int64, status bool) {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	topicNewDynamicCache.CheckStatus[topicID] = status
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) GetCheckStatus(topicID int64) bool {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	return topicNewDynamicCache.CheckStatus[topicID]
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) SetCheckTime(topicID int64) {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	topicNewDynamicCache.CheckTime[topicID] = util.GetTime()
}

func (topicNewDynamicCache *topicNewDynamicCacheStructure) GetCheckTime(topicID int64) int64 {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	return topicNewDynamicCache.CheckTime[topicID]
}

//判断查询是否正在进行中
//是：true
func (topicNewDynamicCache *topicNewDynamicCacheStructure) IsTopicNewDynamicQuerying(topicID int64) bool {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	time := util.GetTime()
	if topicNewDynamicCache.CheckStatus[topicID] == true && topicNewDynamicCache.CheckTime[topicID]+topicNewDynamicCache.GetCheckTimeOffset() > time {
		return true
	}
	return false
}

//获取话题动态缓存时间
func (topicNewDynamicCache *topicNewDynamicCacheStructure) GetTopicNewDynamicCacheTime(topicID int64) int64 {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	return topicNewDynamicCache.TopicDynamicCacheTime[topicID]
}

//获取话题缓存列表
func (topicNewDynamicCache *topicNewDynamicCacheStructure) GetTopicNewDynamicCacheIds(page, pageSize int, topicID int64) []int64 {
	topicNewDynamicCache.LockChan[topicID] <- struct{}{}
	defer func() {
		<-topicNewDynamicCache.LockChan[topicID]
	}()
	curPage := easygo.If(page > 1, page-1, 0).(int)
	startNum := curPage * pageSize
	endNum := startNum + pageSize
	if startNum > topicNewDynamicCache.TopicDynamicCount[topicID] {
		return []int64{}
	}
	if topicNewDynamicCache.TopicDynamicList == nil {
		return []int64{}
	}
	if _, ok := topicNewDynamicCache.TopicDynamicList[topicID]; !ok {
		return []int64{}
	}

	if endNum > topicNewDynamicCache.TopicDynamicCount[topicID] {
		endNum = topicNewDynamicCache.TopicDynamicCount[topicID]
	}

	return topicNewDynamicCache.TopicDynamicList[topicID][startNum:endNum]
}

//设置最新话题缓存
func (topicNewDynamicCache *topicNewDynamicCacheStructure) SetTopicNewDynamicCacheList(topicID int64, countTotal int) bool {
	isTopicNewDynamicQuerying := topicNewDynamicCache.IsTopicNewDynamicQuerying(topicID)
	if isTopicNewDynamicQuerying {
		return true
	}
	//标志为查询中
	topicNewDynamicCache.SetCheckStatus(topicID, true)
	topicNewDynamicCache.SetCheckTime(topicID)
	defer func() {
		topicNewDynamicCache.SetCheckStatus(topicID, false)
		//logs.Info("topicNewDynamicCache:",topicNewDynamicCache)
	}()
	topicNewDynamicTopList := GetTopicNewDynamicTopList(topicID)
	topicNewDynamicIds := make([]int64, 0, topicNewDynamicCache.GetDefCacheNumber()/10)
	for _, v := range topicNewDynamicTopList {
		topicNewDynamicIds = append(topicNewDynamicIds, v.GetLogId())
	}
	topNum := len(topicNewDynamicIds)
	topicNewDynamicCache.SetTopicDynamicTopCount(topicID, topNum)
	if topNum < topicNewDynamicCache.GetDefCacheNumber() {
		dynamicIdMapping := make(map[int64]struct{})
		for _, v := range topicNewDynamicIds {
			dynamicIdMapping[v] = struct{}{}
		}
		topicNewDynamicList := GetTopicNewDynamicList(0, topicNewDynamicCache.GetDefCacheNumber(), topicID)
		//合并top和非top的动态id
		for _, v := range topicNewDynamicList {
			if _, ok := dynamicIdMapping[v.GetLogId()]; !ok {
				topicNewDynamicIds = append(topicNewDynamicIds, v.GetLogId())
			}
		}
	}
	topicNewDynamicCache.SetTopicDynamicUnTopCount(topicID, len(topicNewDynamicIds)-topNum)
	topicNewDynamicCache.SetTopicDynamicCount(topicID, countTotal)
	topicNewDynamicCache.SetTopicDynamicList(topicID, topicNewDynamicIds)
	return true
}

//设置热门话题缓存
func (topicNewDynamicCache *topicNewDynamicCacheStructure) SetTopicHotDynamicCacheList(topicID int64, hotScore int32, countTotal int) bool {
	//如此缓存正在查询中，则直接返回成功
	isTopicNewDynamicQuerying := topicNewDynamicCache.IsTopicNewDynamicQuerying(topicID)
	if isTopicNewDynamicQuerying {
		return true
	}
	//标志为查询中
	topicNewDynamicCache.SetCheckStatus(topicID, true)
	topicNewDynamicCache.SetCheckTime(topicID)
	defer func() {
		topicNewDynamicCache.SetCheckStatus(topicID, false)
		logs.Info("topicHotDynamicCache:", topicNewDynamicCache)
	}()
	topicNewDynamicTopList := GetTopicNewDynamicTopList(topicID)
	topicNewDynamicIds := make([]int64, 0, topicNewDynamicCache.GetDefCacheNumber()/10)
	for _, v := range topicNewDynamicTopList {
		topicNewDynamicIds = append(topicNewDynamicIds, v.GetLogId())
	}
	topNum := len(topicNewDynamicIds)
	topicNewDynamicCache.SetTopicDynamicTopCount(topicID, topNum)
	if topNum < topicNewDynamicCache.GetDefCacheNumber() {
		dynamicIdMapping := make(map[int64]struct{})
		for _, v := range topicNewDynamicIds {
			dynamicIdMapping[v] = struct{}{}
		}
		topicNewDynamicList := GetTopicHotDynamicList(hotScore, 0, topicNewDynamicCache.GetDefCacheNumber(), topicID)
		//合并top和非top的动态id
		for _, v := range topicNewDynamicList {
			if _, ok := dynamicIdMapping[v.GetLogId()]; !ok {
				topicNewDynamicIds = append(topicNewDynamicIds, v.GetLogId())
			}
		}
	}
	topicNewDynamicCache.SetTopicDynamicUnTopCount(topicID, len(topicNewDynamicIds)-topNum)
	topicNewDynamicCache.SetTopicDynamicCount(topicID, countTotal)
	topicNewDynamicCache.SetTopicDynamicList(topicID, topicNewDynamicIds)
	return true
}
