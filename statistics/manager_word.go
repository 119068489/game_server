package statistics

import (
	b64 "encoding/base64"
	"encoding/json"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"strings"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

type WordManager struct {
	Words    []*share_message.CrawlWords
	GrabTags []*share_message.GrabTag
}

const WORK_TIME = time.Second * 300

func NewWordManager() *WordManager {
	p := &WordManager{}
	p.Init()
	return p
}
func (self *WordManager) Init() {
	//启动定时任务:10秒后开始工作，每10分钟处理一次
	self.Words = self.GetCrawlWords()
	self.GrabTags = self.GetGrabTags()
	easygo.AfterFunc(time.Second*2, self.Update)
}

//获取聊天最大记录id
func (self *WordManager) GetMaxChatLog(key string) int64 {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ID_GENERATOR)
	defer closeFun()
	var identity for_game.Identity
	err := col.Find(bson.M{"_id": key}).One(&identity)
	if err != nil {
		return 0
	}
	return int64(*identity.Value)
}

//获取系统日志
func (self *WordManager) GetSystemLog() *share_message.SystemLog {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SYSTEM_LOG)
	defer closeFun()
	var sysInfo *share_message.SystemLog
	err := col.Find(bson.M{"_id": 1}).One(&sysInfo)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		//数据库还没存在
		return &share_message.SystemLog{
			Id:             easygo.NewInt64(1),
			WordPerSonalId: easygo.NewInt64(0),
			WordTeamId:     easygo.NewInt64(0),
		}
	}
	return sysInfo
}

//获取抓取词库
func (self *WordManager) GetCrawlWords() []*share_message.CrawlWords {
	var words []*share_message.CrawlWords
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CRAWL_WORDS)
	defer closeFun()
	err := col.Find(bson.M{}).All(&words)
	easygo.PanicError(err)
	return words
}

//获取标签库
func (self *WordManager) GetGrabTags() []*share_message.GrabTag {
	var tags []*share_message.GrabTag
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_GRABTAG)
	defer closeFun()
	err := col.Find(bson.M{}).All(&tags)
	easygo.PanicError(err)
	return tags
}

//定时查询
func (self *WordManager) Update() {
	self.DealChatLogs()
	easygo.AfterFunc(WORK_TIME, self.Update)
}

//处理聊天关键词统计
func (self *WordManager) DealChatLogs() {
	logs.Info("聊天内容关键词抓取处理开始------>>>")
	self.Words = self.GetCrawlWords()
	self.GrabTags = self.GetGrabTags()
	sysLog := self.GetSystemLog()
	personalMaxLog := sysLog.GetWordPerSonalId()
	teamMaxLog := sysLog.GetWordTeamId()
	logs.Info("personalLog:", sysLog.GetWordPerSonalId())
	logs.Info("teamLog:", sysLog.GetWordTeamId())
	personLogs := self.GetChatLogs(sysLog.GetWordPerSonalId(), for_game.TABLE_PERSONAL_CHAT_LOG)
	logs.Info("personlogs:", len(personLogs))
	teamLogs := self.GetChatLogs(sysLog.GetWordTeamId(), for_game.TABLE_TEAM_CHAT_LOG)
	logs.Info("teamLogs:", len(teamLogs))
	lgs := make(map[PLAYER_ID][]interface{})
	//个人聊天信息整合
	for _, v := range personLogs {
		lg, ok := v.(bson.M)
		if !ok {
			continue
		}
		js, _ := json.Marshal(lg)
		var log *share_message.PersonalChatLog
		_ = json.Unmarshal(js, &log)
		lgs[log.GetTalker()] = append(lgs[log.GetTalker()], log)
		if log.GetLogId() > personalMaxLog {
			personalMaxLog = log.GetLogId()
		}
	}
	//群聊天信息整合
	for _, v := range teamLogs {
		lg, ok := v.(bson.M)
		if !ok {
			continue
		}
		js, _ := json.Marshal(lg)
		var log *share_message.TeamChatLog
		_ = json.Unmarshal(js, &log)
		lgs[log.GetTalker()] = append(lgs[log.GetTalker()], log)
		if log.GetLogId() > teamMaxLog {
			teamMaxLog = log.GetLogId()
		}
	}
	if len(lgs) == 0 {
		logs.Info("没有可以处理记录")
		return
	}
	self.CheckChatLogs(lgs)
	//记录当前处理了的日志
	sysLog.WordPerSonalId = easygo.NewInt64(personalMaxLog)
	sysLog.WordTeamId = easygo.NewInt64(teamMaxLog)
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SYSTEM_LOG)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": sysLog.GetId()}, bson.M{"$set": sysLog})
	easygo.PanicError(err)
	logs.Info("关键词抓取处理完成------>>>")

}

//获取聊天内容日志
func (self *WordManager) GetChatLogs(start int64, tab string) []interface{} {
	var logs []interface{}
	col, closeFun := easygo.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, tab)
	defer closeFun()
	err := col.Find(bson.M{"_id": bson.M{"$gt": start}, "Type": for_game.TALK_CONTENT_WORD}).All(&logs)
	easygo.PanicError(err)
	return logs
}

//聊天内容检测
func (self *WordManager) CheckChatLogs(chatlogs map[PLAYER_ID][]interface{}) {
	if len(chatlogs) == 0 {
		logs.Info("需要处理的log长度为0")
		return
	}
	var keys []int64
	for pid, _ := range chatlogs {
		keys = append(keys, pid)
	}
	//keys, _ := easygo.GetMapKeysValues(chatlogs)
	//logs.Info("说话玩家id:", keys)
	//查询玩家标签信息
	var playerWords []*share_message.PlayerCrawlWords
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_CRAWL_WORDS)
	defer closeFun()
	err := col.Find(bson.M{"_id": bson.M{"$in": keys}}).All(&playerWords)
	easygo.PanicError(err)
	pWords := make(map[int64]*share_message.PlayerCrawlWords)
	for _, v := range playerWords {
		pWords[v.GetId()] = v
	}
	logs.Info("pWords:", len(pWords))
	//查询玩家数据
	var players []*share_message.PlayerBase
	col1, closeFun1 := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun1()
	err = col1.Find(bson.M{"_id": bson.M{"$in": keys}}).All(&players)
	easygo.PanicError(err)
	pPlayers := make(map[int64]*share_message.PlayerBase)
	for _, v := range players {
		pPlayers[v.GetPlayerId()] = v
	}
	logs.Info("pPlayers:", len(pPlayers))
	//统计聊天日志
	for playerId, vals := range chatlogs {
		player := pPlayers[playerId]
		if player == nil {
			logs.Info("找不到玩家信息:", playerId)
			continue
		}
		for _, v := range vals {
			s := ""
			log, ok := v.(*share_message.PersonalChatLog) //强转
			if !ok {
				log1, ok1 := v.(*share_message.TeamChatLog) //强转
				if !ok1 {
					continue
				}
				s = log1.GetContent()
			} else {
				s = log.GetContent()
			}
			if s != "" {
				content, err := b64.StdEncoding.DecodeString(s)
				if err != nil {
					continue
				}
				myWord := pWords[playerId]
				if myWord == nil {
					words := self.CreateNewWords()
					tags := self.CreateNewTags()
					myWord = &share_message.PlayerCrawlWords{
						Id:    easygo.NewInt64(playerId),
						Words: words,
						Tags:  tags,
					}
					pWords[playerId] = myWord
				}
				self.CheckWord(myWord, string(content))
			}
		}

		////暂时一个个写入，后续改批量修改，新值写入mongo
		//_, _ = col.Upsert(bson.M{"_id": playerId}, bson.M{"$set": pWords[playerId]})
		////检测玩家抓取标签
		////logs.Info("统计完成:", pWords[playerId])
		tag := self.GetMaxCrabTag(pWords[playerId])
		//修改玩家标签
		player.GrabTag = easygo.NewInt32(tag)
		//if player.GetIsOnline() {
		//	base := for_game.GetRedisPlayerBase(playerId)
		//	if base != nil {
		//		base.SetGrabTag(player.GetGrabTag())
		//	}
		//} else {
		//	_, _ = col1.Upsert(bson.M{"_id": playerId}, bson.M{"$set": player})
		//}
	}
	//批量修改抓取日志信息

	var data []interface{}
	for pid, v := range pWords {
		b1 := bson.M{"_id": pid}
		b2 := bson.M{"$set": v}
		data = append(data, b1, b2)
	}
	for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_CRAWL_WORDS, data)

	///批量修改玩家的抓取标签
	var pOffLine []*share_message.PlayerBase
	for pid, p := range pPlayers {
		if p.GetIsOnline() { //在线玩家直接修改redis内存值
			base := for_game.GetRedisPlayerBase(pid)
			if base != nil {
				base.SetGrabTag(p.GetGrabTag())
			}
		} else {
			pOffLine = append(pOffLine, p)
		}
	}
	//数据库批量修改不在线的玩家
	var data1 []interface{}
	for _, p := range pOffLine {
		b1 := bson.M{"_id": p.GetPlayerId()}
		b2 := bson.M{"$set": p}
		data1 = append(data1, b1, b2)
	}
	for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE, data1)
}

//获取玩家抓取标签最大的id
func (self *WordManager) GetMaxCrabTag(playerWords *share_message.PlayerCrawlWords) int32 {
	var tag int32
	count := int64(0)
	tags := make(map[int32]int64)
	//统计每个标签出现个数
	for _, v := range playerWords.Words {
		tags[v.GetGrabTag()] += v.GetCount()
	}
	//找出最多个数的标签
	for k, v := range tags {
		if v > count {
			count = v
			tag = k
		}
	}
	//重置标签数据存储
	playerWords.Tags = []*share_message.GrabTag{}
	for _, v := range self.GrabTags {
		if t, ok := tags[v.GetId()]; ok {
			newTag := &share_message.GrabTag{
				Id:    easygo.NewInt32(v.GetId()),
				Name:  easygo.NewString(v.GetName()),
				Count: easygo.NewInt64(t),
			}
			playerWords.Tags = append(playerWords.Tags, newTag)
		}
	}
	return tag
}

//创建新的数据库存储抓取词
func (self *WordManager) CreateNewWords() []*share_message.CrawlWords {
	words := []*share_message.CrawlWords{}
	//for _, w := range self.Words {
	//	craw := &share_message.CrawlWords{
	//		Id:      easygo.NewInt32(w.GetId()),
	//		Name:    easygo.NewString(w.GetName()),
	//		GrabTag: easygo.NewInt32(w.GetGrabTag()),
	//		Count:   easygo.NewInt64(0),
	//	}
	//	words = append(words, craw)
	//}
	return words
}

//创建新的数据库抓取标签
func (self *WordManager) CreateNewTags() []*share_message.GrabTag {
	tags := []*share_message.GrabTag{}
	//for _, v := range self.GrabTags {
	//	tag := &share_message.GrabTag{
	//		Id:    easygo.NewInt32(v.GetId()),
	//		Name:  easygo.NewString(v.GetName()),
	//		Count: easygo.NewInt64(0),
	//	}
	//	tags = append(tags, tag)
	//}
	return tags
}

//检测抓取词统计
func (self *WordManager) CheckWord(pWord *share_message.PlayerCrawlWords, content string) {
	wordIds := make(map[int32]*share_message.CrawlWords)
	for _, v := range pWord.GetWords() {
		wordIds[v.GetId()] = v
	}
	for _, w := range self.Words {
		//抓取词统计
		if v, ok := wordIds[w.GetId()]; ok {
			if strings.Contains(content, w.GetName()) {
				v.Count = easygo.NewInt64(v.GetCount() + 1)
			}
		} else {
			if strings.Contains(content, w.GetName()) {
				newWord := &share_message.CrawlWords{
					Id:      easygo.NewInt32(w.GetId()),
					Name:    easygo.NewString(w.GetName()),
					GrabTag: easygo.NewInt32(w.GetGrabTag()),
					Count:   easygo.NewInt64(1),
				}
				pWord.Words = append(pWord.GetWords(), newWord)
				wordIds[newWord.GetId()] = newWord
			}
		}
	}
}
