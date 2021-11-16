package main

import (
	. "bufio"
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/hall"
	"game_server/pb/share_message"
	"os"
	"time"

	"github.com/akqp2019/mgo"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

func GetLimitPlayerBase() {
	for {
		base := for_game.GetLimitPlayerBase()
		var saveData []interface{}
		fmt.Println("长度_----------->", len(base))
		for _, v := range base {
			po := &share_message.GeoJson{
				Type:        easygo.NewString("Point"),
				Coordinates: []float64{v.GetX(), v.GetY()},
			}
			v.Points1 = po
			saveData = append(saveData, bson.M{"_id": v.GetPlayerId()}, v)
			obj := for_game.GetRedisPlayerBase(v.GetPlayerId())
			obj.SetPoints(po)

		}
		for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE, saveData)
		if len(base) < 5000 {
			break
		}
	}
}

//初始化服务器配置
func main() {
	initializer := hall.NewInitializer()
	defer func() { // 若是异常了,确保异步日志有成功写盘
		logger := initializer.GetBeeLogger()
		if logger != nil {
			logger.Flush()
		}
	}()
	dict := easygo.KWAT{
		"logName":  "hall",
		"yamlPath": "config_hall.yaml",
	}
	initializer.Execute(dict)
	hall.Initialize(false)
	InitPlayerConstellation()
	//testOrder()
	//hall.SaveRedisData()
	/*	GenPriSession()
		GenTeamSession()
		GenPlayerBase()
		//测试服用正式服数据需要清空下
		GenCleanPlayerToken()*/

	//initMongoIndex() // 添加索引

	//ChangeRandHeadIcon()
	//GetLimitPlayerBase()
	// SetTeamReadId()
	//for_game.ReSetRedPacketStatistics()
	//UpdateTeamChatLogTalker()
	//停服前先保存数据
	//hall.SaveRedisData()
	//time.Sleep(time.Second * 10)
	//changePersonnalChat()
	//logs.Info("迁移个人日志完成")
	//changeTeamChat()
	//logs.Info("群聊迁移完成")
	//for_game.StatisticsRedPacket()
	//logs.Info("红包统计完成")
	//changeTeamMembers()
	//logs.Info("群成员修改完成")
	//changeNextId()
	//logs.Info("log_login 自增迁移完成")
	//////数据表增加索引
	//initMongoIndex()
	//logs.Info("增加查询mongo索引完成")
	//changeWechatAccount()
	//ChangeAccount()
}
func testOrder() {
	for {
		input := NewScanner(os.Stdin)
		input.Scan()
		s := input.Text()
		if s != "" {
			easygo.Spawn(doCmd, s)
		}
	}
}
func doCmd(orderId string) {
	hall.PWebHuiChaoPay.ReqCheckPayOrder(orderId, 1*time.Second)
}

//初始化玩家星座
func InitPlayerConstellation() {
	ids := []int64{1887443575, 1887437297, 1887436238, 1887436239, 1887436261}
	n := 1
	for _, id := range ids {
		tags := make([]int32, 0)
		for i := 0; i < 6; i++ {
			tags = append(tags, int32(for_game.RandInt(1, 30)))
		}
		player := for_game.GetRedisPlayerBase(id)
		player.SetPersonalityTags(tags)
		player.SetBgImageUrl("https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/bon8.jpeg")

		player.SetMixId(int64(n))
		player.SaveToMongo()
		n += 1
	}
	logs.Info("初始化完毕")
}
func SetTeamReadId() {
	logs.Info("开始设置群已读id")
	var teamData []*share_message.TeamData
	col, closeFun := hall.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAM_DATA)
	defer closeFun()
	err := col.Find(bson.M{}).All(&teamData)
	easygo.PanicError(err)
	var teamMembers []*share_message.PersonalTeamData
	for _, team := range teamData {
		teamObj := for_game.GetRedisTeamObj(team.GetId())
		memberObj := for_game.GetRedisTeamPersonalObj(team.GetId())
		teamMembers = memberObj.GetRedisTeamPersonal()
		var saveMembers []interface{}
		for _, m := range teamMembers {
			if m.GetReadId() > teamObj.GetSessionLogMaxId() {
				logId := teamObj.GetSessionLogMaxId()
				if logId > 100 {
					logId = logId - 100
				} else if logId > 50 {
					logId = logId - 50
				}
				memberObj.ReadTeamChatLog(m.GetPlayerId(), logId)
				m.ReadId = easygo.NewInt64(logId)
				saveMembers = append(saveMembers, bson.M{"_id": m.GetId()}, m)
			}
		}
		for_game.UpsertAll(hall.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_TEAMMEMBER, saveMembers)
		if len(teamMembers) > 500 {
			time.Sleep(time.Second)
		}
	}
	logs.Info("开始设置群已读id  完成--------》》》")
}

//红包统计重新生成

//玩家支付密码迁移到玩家身上
func ChangeAccount() {
	var accountList []*share_message.PlayerAccount
	col, closeFun := hall.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_ACCOUNT)
	defer closeFun()
	err := col.Find(nil).All(&accountList)
	easygo.PanicError(err)

	col1, closeFun1 := hall.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun1()
	for _, account := range accountList {
		err := col1.Update(bson.M{"_id": account.GetPlayerId()}, bson.M{"$set": bson.M{"PayPassword": account.GetPayPassword()}})
		if err != nil {
			logs.Info(err)
		}
	}
}

func initMongoIndex() {
	//for_game.EnsureIndexKey(easygo.MongoLogMgr, for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_TEAM_CHAT_LOG, "TeamId")
	//for_game.EnsureIndexKey(easygo.MongoLogMgr, for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PERSONAL_CHAT_LOG, "TargetId")
	//for_game.EnsureIndexKey(easygo.MongoLogMgr, for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_RED_PACKET_LOG, "RedPacketId", "PlayerId")
	//for_game.EnsureIndexKey(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_GOLDCHANGELOG, "PlayerId")
	//for_game.EnsureIndexKey(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_ORDER, "PlayerId")
	//for_game.EnsureIndexKey(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_ACCOUNT, "Account")
	//for_game.EnsureIndexKey(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_CART, "player_id")
	//for_game.EnsureIndexKey(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_RED_PACKET, "Sender")
	//for_game.EnsureIndexKey(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_TEAMMEMBER, "TeamId")

	//for_game.EnsureIndexKey(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_ATTENTION, "PlayerId")
	//for_game.EnsureIndexKey(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_COMMENT, "LogId")
	//for_game.EnsureIndexKey(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC, "PlayerId", "TopicId", "Statue", "HostScore")
	//for_game.EnsureIndexKey(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_COMMENT_ZAN, "DynamicId")
	//for_game.EnsureIndexKey(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_ZAN, "DynamicId", "PlayerId")
	//for_game.EnsureIndexKey(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC, "IsRecommend", "Status")
	for_game.EnsureIndexKey(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE, "Sex")
	logs.Info("添加索引完成")
}

//迁移个人聊天记录
func changePersonnalChat() {
	//先查出id_generator
	col, closeFun := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_PERSON_LOG, for_game.TABLE_ID_GENERATOR)
	defer closeFun()
	val := []easygo.KWAT{}
	err := col.Find(bson.M{}).All(&val)
	easygo.PanicError(err)
	i := int64(1)
	count := 0
	col2, closeFun2 := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, "personal_chat_log")
	defer closeFun2()
	for _, v := range val {
		k := v.GetString("_id")
		va := v.GetString("Value")
		//logs.Info("k,v:", k, va)
		data := make([]*share_message.PersonalChatLog, easygo.Atoi(va))
		col1, closeFun1 := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_PERSON_LOG, k)
		err = col1.Find(bson.M{}).All(&data)
		closeFun1()
		easygo.PanicError(err)
		//logs.Info("data:", data)
		var saveData []interface{}
		for j := 0; j < len(data); j++ {
			data[j].LogId = easygo.NewInt64(i)
			saveData = append(saveData, data[j])
			i += 1

		}
		err = col2.Insert(saveData...)
		easygo.PanicError(err)
		//err = col1.DropCollection()
		//easygo.PanicError(err)
		count += easygo.Atoi(va)
	}
	logs.Info("总记录条数:", count)
	identity := &for_game.Identity{
		Key:   easygo.NewString("personal_chat_log"),
		Value: easygo.NewUint64(count),
	}
	col3, closeFun3 := hall.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ID_GENERATOR)
	defer closeFun3()
	col3.Insert(&identity)
}

//迁移群聊天记录
func changeTeamChat() {

	//先查出id_generator
	col, closeFun := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_TEAM_LOG, for_game.TABLE_ID_GENERATOR)
	defer closeFun()
	val := []easygo.KWAT{}
	err := col.Find(bson.M{}).All(&val)
	easygo.PanicError(err)
	i := int64(1)
	count := 0
	col2, closeFun2 := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, "team_chat_log")
	defer closeFun2()
	for _, v := range val {
		k := v.GetString("_id")
		va := v.GetString("Value")
		//logs.Info("k,v:", k, va)
		data := make([]*share_message.TeamChatLog, easygo.Atoi(va))
		col1, closeFun1 := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_TEAM_LOG, k)
		err = col1.Find(bson.M{}).All(&data)
		closeFun1()
		easygo.PanicError(err)
		//logs.Info("data:", data)
		var saveData []interface{}
		for j := 0; j < len(data); j++ {
			data[j].TeamLogId = easygo.NewInt64(data[j].LogId)
			data[j].LogId = easygo.NewInt64(i)
			saveData = append(saveData, data[j])
			i += 1

		}
		err = col2.Insert(saveData...)
		easygo.PanicError(err)
		//err = col1.DropCollection()
		//easygo.PanicError(err)
		count += easygo.Atoi(va)
	}
	logs.Info("总记录条数:", count)
	identity := &for_game.Identity{
		Key:   easygo.NewString("team_chat_log"),
		Value: easygo.NewUint64(count),
	}
	col3, closeFun3 := hall.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ID_GENERATOR)
	defer closeFun3()
	col3.Insert(&identity)
}

//迁移个人群信息转换
func changeTeamMembers() {
	logs.Info("开始转换------")
	col, closeFun := hall.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAMMEMBER_DATA)
	defer closeFun()
	val := []*share_message.AllTeamData{}
	err := col.Find(bson.M{}).All(&val)
	easygo.PanicError(err)
	var newMembers []interface{}
	for _, team := range val {
		for _, member := range team.GetData() {
			member.TeamId = easygo.NewInt64(team.GetTeamId())
			member.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_TEAMMEMBER))
			newMembers = append(newMembers, member)
		}
	}
	col1, closeFun1 := hall.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAMMEMBER)
	defer closeFun1()
	err1 := col1.Insert(newMembers...)
	easygo.PanicError(err1)
	//插入成员id
	identity := &for_game.Identity{
		Key:   easygo.NewString("team_members"),
		Value: easygo.NewUint64(for_game.CurrentId(for_game.TABLE_TEAMMEMBER)),
	}
	col3, closeFun3 := hall.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ID_GENERATOR)
	defer closeFun3()
	col3.Insert(&identity)
	logs.Info("转换完成。。。。。。。")
}

//登录信息迁移
func changeNextId() {
	col, closeFun := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_ID_GENERATOR)
	defer closeFun()
	var val []interface{}
	err := col.Find(bson.M{}).All(&val)
	easygo.PanicError(err)
	col3, closeFun3 := hall.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ID_GENERATOR)
	defer closeFun3()
	err = col3.Insert(val...)
	easygo.PanicError(err)
}

func changeWechatAccount() {
	col, closeFun := hall.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_ACCOUNT)
	defer closeFun()
	var lst []*share_message.PlayerAccount
	err := col.Find(bson.M{"$where": "this.Account.length>11"}).All(&lst)
	easygo.PanicError(err)
	logs.Info(lst)
	ids := []int64{}
	info1 := make(map[int64]string)
	for _, m := range lst {
		ids = append(ids, m.GetPlayerId())
		openId := m.GetOpenId()
		info1[m.GetPlayerId()] = openId
		m.Account = easygo.NewString(openId)
		_, err1 := col.Upsert(bson.M{"_id": m.GetPlayerId()}, m)
		if err1 != nil {
			panic(err1)
		}
	}

	col1, closeFun1 := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun1()
	for pid, openid := range info1 {
		err := col1.Update(bson.M{"_id": pid}, bson.M{"$set": bson.M{"Phone": openid}})
		if err != nil {
			logs.Error(err)
			continue
		}
	}
}

//处理群聊天记录，把说话人的头像和昵称赋值上
func UpdateTeamChatLogTalker() {
	col, closeFun := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_TEAM_CHAT_LOG)
	defer closeFun()
	var lst []*share_message.TeamChatLog
	err := col.Find(bson.M{}).All(&lst)
	easygo.PanicError(err)
	ids := []int64{}
	for _, log := range lst {
		if easygo.Contain(ids, log.GetTalker()) {
			continue
		}
		ids = append(ids, log.GetTalker())
	}
	pMap := for_game.GetAllPlayerBase(ids, false)

	var saveData []interface{}
	for _, log := range lst {
		p, ok := pMap[log.GetTalker()]
		if !ok {
			continue
		}
		log.TalkerHeadUrl = easygo.NewString(p.GetHeadIcon())
		log.TalkerName = easygo.NewString(p.GetNickName())
		b1 := bson.M{"_id": log.GetLogId()}
		saveData = append(saveData, b1, log)
	}
	for_game.UpsertAll(easygo.MongoLogMgr, for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_TEAM_CHAT_LOG, saveData)
	logs.Info("聊天记录头像刷新完毕")
}

//修改玩家手机号
func UpdataPlayerPhone() {
	id := int64(1887530532)
	areaCode := "855"
	phone := "0962340200"
	// id := int64(1887550224)
	// areaCode := "86"
	// phone := "13888"
	base := for_game.GetRedisPlayerBase(id)
	account := for_game.GetRedisAccountObj(id)
	if base != nil {
		account.SetAccount(phone)
		base.SetAreaCode(areaCode)
		base.SetPhone(phone)
	}
	account.SaveToMongo()
	base.SaveToMongo()
	logs.Info("更新手机号码完毕")
}

//批量修改运营账号的头像
func ChangeRandHeadIcon() {
	or1 := bson.M{"HeadIcon": bson.M{"$regex": bson.RegEx{Pattern: "girl_", Options: "im"}}}
	or2 := bson.M{"HeadIcon": bson.M{"$regex": bson.RegEx{Pattern: "boy_", Options: "im"}}}
	queryBson := bson.M{"$or": []bson.M{or1, or2}, "Types": bson.M{"$gt": 1}}
	falseList, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE, queryBson, 0, 0)
	for _, fli := range falseList {
		head := for_game.GetRandRealHeadIcon(fli.(bson.M)["Sex"].(int))
		playerid := fli.(bson.M)["_id"].(int64)
		one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC, bson.M{"PlayerId": playerid, "Photo": bson.M{"$ne": nil}})
		pmg := for_game.GetRedisPlayerBase(playerid)
		if pmg != nil {
			pmg.SetHeadIcon(head)
			if one == nil {
				pmg.SetStatus(2)
			}
		}
	}
}

//===================================私聊增加会话id========================================

func GetChatLog(logId int64) []*share_message.PersonalChatLog {
	col2, closeFun2 := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, "personal_chat_log")
	defer closeFun2()
	chatLogs := make([]*share_message.PersonalChatLog, 0)
	queryBson := bson.M{}
	if logId > 0 {
		queryBson["_id"] = bson.M{"$gt": logId}
	}
	err := col2.Find(queryBson).Sort("_id").Limit(5000).Sort("_id").All(&chatLogs)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return chatLogs
}

// 生成key
func makeKey(talk, target int64) string {
	var key string
	if talk > target {
		key = for_game.MakeNewString(target, talk)
	} else {
		key = for_game.MakeNewString(talk, target)
	}
	return key
}

// 生成私聊聊天会话,添加sessionId进去表
func GenPriSession() {
	// 一千条一千条的查询
	var maxId int64
	idMap := make(map[string]int64)
	//playerSession := make(map[int64][]string)
	for {
		chatLog := GetChatLog(maxId)
		logs.Info(" 处私聊记录:", maxId, len(chatLog))
		if len(chatLog) == 0 {
			break
		}
		var saveData1 []interface{}
		var saveData22 []interface{}
		var saveData2 []*share_message.ChatSession

		for _, v := range chatLog {
			mKey := makeKey(v.GetTalker(), v.GetTargetId())
			// 自增的id
			maxLogId, ok1 := idMap[mKey]
			if !ok1 {
				maxLogId = 1
				idMap[mKey] = maxLogId
			} else {
				maxLogId = maxLogId + 1
				idMap[mKey] = maxLogId
			}

			//var readId1 int64
			v.SessionId = easygo.NewString(mKey)
			v.TalkLogId = easygo.NewInt64(maxLogId)
			saveData1 = append(saveData1, bson.M{"_id": v.GetLogId()}, v)

			data2 := &share_message.ChatSession{
				Id:        easygo.NewString(mKey),
				Type:      easygo.NewInt32(1),
				PlayerIds: []int64{v.GetTalker(), v.GetTargetId()},
				MaxLogId:  easygo.NewInt64(maxLogId),
				ReadInfo:  nil,
			}
			saveData2 = append(saveData2, data2)

		}
		for_game.UpsertAll(easygo.MongoLogMgr, for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PERSONAL_CHAT_LOG, saveData1)

		for _, d22 := range saveData2 {
			readInfoList := make([]*share_message.ReadLogInfo, 0)
			ids := d22.GetPlayerIds()
			for _, idv := range ids {
				//playerSession[idv] = append(playerSession[idv], d22.GetId())
				readId := GetCurReadId(d22.GetId(), idv, saveData1)
				readLogInfo := &share_message.ReadLogInfo{
					PlayerId: easygo.NewInt64(idv),
					ReadId:   easygo.NewInt64(readId),
				}
				readInfoList = append(readInfoList, readLogInfo)
			}
			d22.ReadInfo = readInfoList
			saveData22 = append(saveData22, bson.M{"_id": d22.GetId()}, d22)
		}
		if len(saveData22) > 0 {
			for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_MONGODB_CHAT_SESSION, saveData22)
		}
		maxId = chatLog[len(chatLog)-1].GetLogId()
		logs.Info(" 处私聊记录完成:", maxId, len(chatLog))
	}
	for_game.EnsureIndexKey(easygo.MongoLogMgr, for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PERSONAL_CHAT_LOG, "SessionId")
	////玩家会话列表
	//saveList11 := make([]interface{}, 0)
	//for pid, data := range playerSession {
	//	session := &share_message.PlayerChatSession{
	//		PlayerId:   easygo.NewInt64(pid),
	//		SessionIds: data,
	//	}
	//	saveList11 = append(saveList11, bson.M{"_id": pid}, session)
	//}
	//if len(saveList11) > 0 {
	//	for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_CHAT_SESSION, saveList11)
	//}

}
func GetCurReadId(sessionId string, playerId int64, logData []interface{}) int64 {
	var readId, maxId int64
	for _, lo := range logData {
		log, ok := lo.(*share_message.PersonalChatLog)
		if !ok {
			continue
		}
		if sessionId == log.GetSessionId() {
			if log.GetTargetId() == playerId && readId < log.GetTalkLogId() && !log.GetIsRead() {
				readId = log.GetTalkLogId() - 1
				break
			}
			if maxId < log.GetTalkLogId() {
				maxId = log.GetTalkLogId()
			}
		}
	}
	if readId == 0 {
		readId = maxId
	}
	return readId
}

//===================================私聊增加会话id========================================
//===================================群聊增加会话id========================================
//获取群信息
func GetTeamInfo(ids []int64) map[int64]*share_message.TeamData {
	col, closeFun := hall.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TEAM_DATA)
	defer closeFun()
	teams := make([]*share_message.TeamData, 0)
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&teams)
	easygo.PanicError(err)
	data := make(map[int64]*share_message.TeamData)
	for _, team := range teams {
		data[team.GetId()] = team
	}
	return data
}
func GetTeamIds(logs []*share_message.TeamChatLog) []int64 {
	ids := make([]int64, 0)
	for _, log := range logs {
		if easygo.Contain(ids, log.GetTeamId()) {
			continue
		}
		ids = append(ids, log.GetTeamId())
	}
	return ids
}

// 生成群聊聊天会话,添加sessionId进去表
func GenTeamSession() {
	// 一千条一千条的查询
	var logId int64
	m := make(map[string]*share_message.ChatSession)

	for {
		tLogs := GetTeamChatLog(logId)
		logs.Info(" 处理群记录:", logId, len(tLogs))
		if len(tLogs) == 0 {
			break
		}
		logId = tLogs[len(tLogs)-1].GetLogId()
		//var saveData2 []interface{}
		var updateData []interface{} // 修改状态
		//var teamData1 []interface{}  // 群数据
		//群信息
		teamIds := GetTeamIds(tLogs)
		teamMap := GetTeamInfo(teamIds)
		for _, v := range tLogs {
			teamKey := easygo.AnytoA(v.GetTeamId())
			session, ok := m[teamKey]
			maxLogId := v.GetTeamLogId()
			if !ok {
				teamData := teamMap[v.GetTeamId()]
				// 获得群成员
				//teamData := for_game.GetRedisTeamObj(v.GetTeamId())

				pids := make([]int64, 0)
				var headURL string
				var maxLogId int64
				if teamData != nil {
					pids = append(pids, teamData.GetMemberList()...)
					headURL = teamData.GetHeadUrl()
					if headURL == "" { // 群主的头像
						owner := for_game.GetRedisPlayerBase(teamData.GetOwner())
						if owner != nil {
							headURL = owner.GetHeadIcon()
						}
					}
					//tt:= teamData.GetRedisTeam()

					//teamData1 = append(teamData1, bson.M{"_id": teamData.GetId()}, bson.M{"$set": bson.M{"SessionId": teamKey}})

				}
				sessionName := v.GetTeamName()
				if sessionName == "" {
					// 取前三个群成员的昵称
					var count int
					for _, pid := range pids {
						if count > 2 {
							break
						}
						pb := for_game.GetRedisPlayerBase(pid)
						if pb != nil {
							sessionName += pb.GetNickName()
							if count != 2 {
								sessionName += "、"
							}
						}

						count++
					}

				}
				session = &share_message.ChatSession{
					Id:             easygo.NewString(teamKey),
					Type:           easygo.NewInt32(2),
					PlayerIds:      pids,
					SessionName:    easygo.NewString(sessionName),
					SessionHeadUrl: easygo.NewString(headURL),
					MaxLogId:       easygo.NewInt64(maxLogId),
					TeamName:       easygo.NewString(v.GetTeamName()),
				}
				m[teamKey] = session

				// team_session
				//saveData2 = append(saveData2, bson.M{"_id": teamKey}, share_message.ChatSession{
				//	Id:             easygo.NewString(teamKey),
				//	Type:           easygo.NewInt32(2),
				//	PlayerIds:      pids,
				//	SessionName:    easygo.NewString(sessionName),
				//	SessionHeadUrl: easygo.NewString(headURL),
				//	MaxLogId:       easygo.NewInt64(maxLogId),
				//	TeamName:       easygo.NewString(v.GetTeamName()),
				//})

			}
			if maxLogId > session.GetMaxLogId() {
				session.MaxLogId = easygo.NewInt64(maxLogId)
			}
			v.Status = easygo.NewInt32(0)
			v.SessionId = easygo.NewString(teamKey)
			updateData = append(updateData, bson.M{"_id": v.GetLogId()}, v)

		}

		if len(updateData) > 0 {
			for_game.UpsertAll(easygo.MongoLogMgr, for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_TEAM_CHAT_LOG, updateData)
		}
		logs.Info(" 处理群记录完成:", logId, len(tLogs))
	}
	if len(m) > 0 {
		saveData2 := make([]interface{}, 0)
		for k, v := range m {
			saveData2 = append(saveData2, bson.M{"_id": k}, v)
		}
		for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_MONGODB_CHAT_SESSION, saveData2)
	}
	for_game.EnsureIndexKey(easygo.MongoLogMgr, for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_TEAM_CHAT_LOG, "SessionId")
}

func GetTeamChatLog(logId int64) []*share_message.TeamChatLog {
	col2, closeFun2 := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_TEAM_CHAT_LOG)
	defer closeFun2()
	chatLogs := make([]*share_message.TeamChatLog, 0)
	queryBson := bson.M{}
	if logId > 0 {
		queryBson["_id"] = bson.M{"$gt": logId}
	}
	err := col2.Find(queryBson).Sort("_id").Limit(5000).Sort("_id").All(&chatLogs)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return chatLogs
}

//===================================群聊增加会话id========================================

//===================================生成玩家会话列表======================================
func GetPlayerInfo(pid int64) []*share_message.PlayerBase {
	col, closeFun := hall.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()
	datas := make([]*share_message.PlayerBase, 0)
	queryBson := bson.M{}
	if pid > 0 {
		queryBson["_id"] = bson.M{"$gt": pid}
	}
	err := col.Find(queryBson).Sort("_id").Limit(5000).Sort("_id").All(&datas)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return datas
}

// 生成私聊聊天会话,添加sessionId进去表
func GenPlayerBase() {
	// 一千条一千条的查询
	var maxId int64 = 1887436000
	//增加索引
	for_game.EnsureIndexKey(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_MONGODB_CHAT_SESSION, "PlayerIds")
	for {
		paleyrList := GetPlayerInfo(maxId)
		logs.Info("处理玩家会话数据:", maxId, len(paleyrList))
		if len(paleyrList) == 0 {
			break
		}
		var saveData []interface{}
		for _, p := range paleyrList {
			sessions := for_game.GetSessionDataByPlayerId(p.GetPlayerId())
			ids := make([]string, 0)
			for _, s := range sessions {
				ids = append(ids, s.GetId())
			}
			saveData = append(saveData, bson.M{"_id": p.GetPlayerId()}, &share_message.PlayerChatSession{
				PlayerId:   easygo.NewInt64(p.GetPlayerId()),
				SessionIds: ids,
			})
		}
		if len(saveData) > 0 {
			for_game.UpsertAll(hall.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_CHAT_SESSION, saveData)
		}
		maxId = paleyrList[len(paleyrList)-1].GetPlayerId()
		logs.Info("处理玩家会话数据完成:", maxId, len(paleyrList))
	}

}

//===================================生成玩家会话列表======================================
func GenCleanPlayerToken() {
	col, closeFun := hall.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()
	_, err := col.UpdateAll(bson.M{"Token": bson.M{"$ne": ""}}, bson.M{"$set": bson.M{"Token": ""}})
	easygo.PanicError(err)
}
