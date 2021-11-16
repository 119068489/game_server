package main

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/hall"
	"game_server/login"
	"game_server/pb/client_hall"
	"game_server/pb/h5_wish"
	"game_server/pb/share_message"
	"game_server/wish"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"

	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/alipay"
)

func init() {
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
	hall.Initialize()
	wish.Initialize()
	////启动etcd
	//hall.PClient3KVMgr.StartClintTV3()
	//defer hall.PClient3KVMgr.Close() //关闭etcd
	////把已启动的服务增加到内存管理
	//for_game.InitExistServer(hall.PClient3KVMgr, hall.PServerInfoMgr, hall.PServerInfo)
	//hall.PWebApiForServer = hall.NewWebHttpForServer(hall.PServerInfo.GetServerApiPort())

}

func TesTestUpsetTopicToDB() {
	t := &share_message.Topic{
		Id:          easygo.NewInt64(5),
		Name:        easygo.NewString("#盐焗鸡#"),
		TopicTypeId: easygo.NewInt64(2),
		TopicClass:  easygo.NewInt32(1),
	}
	for_game.UpsetTopicToDB(t)
}
func TestUpsetTopicTypeToDB() {
	t := &share_message.TopicType{
		Id:         easygo.NewInt64(1),
		TopicClass: easygo.NewInt32(1),
		Name:       easygo.NewString("时尚美妆"),
	}
	for_game.UpsetTopicTypeToDB(t)
}
func TestGetTopicByNameFromDB() {
	tp := for_game.GetTopicByNameFromDB("#我是话题1#")
	if tp == nil {
		return
	}
	fmt.Println(tp.GetId(), tp.GetTopicTypeId(), tp.GetName())
}
func TestGetBSTopicTypeListFormDB() {
	tp := for_game.GetBSTopicTypeListByClassFormDB(for_game.TOPIC_CLASS_BS)

	fmt.Println(tp)
}

func TestInsertPlayerAttentionToDB() {
	pat := &share_message.PlayerAttentionTopic{
		Id:         easygo.NewInt64(1),
		TopicId:    easygo.NewInt64(1),
		PlayerId:   easygo.NewInt64(666),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
	}
	for_game.InsertPlayerAttentionToDB(pat)
}

func TestDelPlayerAttentionFromDB() {
	for_game.DelPlayerAttentionFromDB(666, 1)
}
func TestUpdateTopicFansToDB() {
	for_game.UpdateTopicFansToDB(1, -1)
}

func TestGetTopicPlayerAttentionListByPid() {
	topics, count := for_game.GetTopicPlayerAttentionListByPidFromDB(1887439538, 2, 2)
	logs.Info("count--------------", count)
	bb, _ := json.Marshal(topics)
	fmt.Println("----->", string(bb))

}
func TestGetRangeTopicByTypeIdFromDB() {
	db := for_game.GetRangeTopicByTypeIdFromDB(1, []int64{1}, 5)
	bb, _ := json.Marshal(db)
	fmt.Println("----->", string(bb))
}

func TestGetIsHotDynamicByTopicIdFromDB() {
	data, count := for_game.GetIsHotDynamicByTopicIdFromDB(10, 1, 1, 1)
	fmt.Println("条数--------->", count)
	bytes, _ := json.Marshal(data)
	fmt.Println("------->", string(bytes))
}
func TestGetSortTimeDynamicByTopicIdFromDB() {
	data, count := for_game.GetSortTimeDynamicByTopicIdFromDB(2, 1, 1)
	fmt.Println("条数--------->", count)
	bytes, _ := json.Marshal(data)
	fmt.Println("------->", string(bytes))
}

func TestCheckPlayerIsAttentionTopic() {
	topic := &share_message.Topic{
		Id: easygo.NewInt64(1),
	}
	for_game.CheckPlayerIsAttentionTopic(1887439539, topic)
	fmt.Printf("------>%+v", topic)
}

func TestLikeSearchTopicByNameFromDB() {
	list := for_game.LikeSearchTopicByNameFromDB("唧唧")
	bytes, _ := json.Marshal(list)
	fmt.Println("---------->", string(bytes))
}
func TestIncTopicHotScoreNumToDB() {
	for_game.IncTopicHotScoreNumToDB(1, -4)
}

func TestGetHighHotScoreTopicFromDB() {
	list := for_game.GetHighHotScoreTopicFromDB(3)
	bytes, _ := json.Marshal(list)
	fmt.Println("----------->", string(bytes))
}

func TestGetRangeIsRecommendTopicFromDB() {
	list := for_game.GetRangeIsRecommendTopicFromDB(2, []int64{5, 6})
	//list := for_game.GetRangeIsRecommendTopicFromDB(2, nil)
	bytes, _ := json.Marshal(list)
	fmt.Println("----------->", string(bytes))
}

func TestGetBSTopicTypeListByTypeIdPageFormDB() {
	topics, i := for_game.GetBSTopicTypeListByTypeIdPageFormDB(2, 2, 1)
	bytes, _ := json.Marshal(topics)
	fmt.Println("count----->", i)
	fmt.Println("---->", string(bytes))

}

func TestGetHotTopicFromDB() {
	list := for_game.GetHotTopicFromDB(2)
	bytes, _ := json.Marshal(list)
	fmt.Println("--->", string(bytes))
}

func TestGetRecommendTopicByPageFromDB() {
	topics, count := for_game.GetRecommendTopicByPageFromDB(1, 1)
	fmt.Println("count----->", count)
	bytes, _ := json.Marshal(topics)
	fmt.Println("list--------->", string(bytes))
}

func TestFlushTopic() {
	lists, count := for_game.FlushTopic(1, 10)
	fmt.Println("---->", count)
	bytes, _ := json.Marshal(lists)
	fmt.Println("list----->", string(bytes))
}

func TestGetPlayerAttentionTopicsByPidFromDB() {
	tps := for_game.GetPlayerAttentionTopicsByPidFromDB(1887439538)
	bytes, _ := json.Marshal(tps)
	fmt.Println("----->", string(bytes))
}

func TestGetHotTopicByPageFromDB() {
	topics, i := for_game.GetHotTopicByPageFromDB(2, 1)
	fmt.Println("count----------->", i)
	bytes, _ := json.Marshal(topics)

	fmt.Println("topics----------->", string(bytes))
}

func TestGetRandPlayerByDynamicCountFromDB() {
	list := for_game.GetRandPlayerByDynamicCountFromDB(3, 5)
	bytes, _ := json.Marshal(list)
	fmt.Println("------->", string(bytes))
}

func TestGetRandPlayer() {
	player := for_game.GetRandPlayer(3)
	bytes, _ := json.Marshal(player)
	fmt.Println("------>", string(bytes))
}

func TestGetAttentionRecommendPlayer() {
	player := for_game.GetAttentionRecommendPlayer()
	bytes, _ := json.Marshal(player)
	fmt.Println("bytes------->", string(bytes))
}

func TestGetSquareAttentionData() {
	req := &client_hall.SquareAttentionReq{
		Page:               easygo.NewInt64(1),
		PageSize:           nil,
		HasAttentionTopic:  nil,
		HasAttentionPlayer: nil,
	}
	// 110114   1887439541
	resp := for_game.GetSquareAttentionData(1887439541, 100, req)
	fmt.Println("动态条数---->", resp.GetCount())
	bytes, _ := json.Marshal(resp.GetTopicList())
	fmt.Println("关注的话题------->", string(bytes))
	i3, _ := json.Marshal(resp.GetDynamicList())
	fmt.Println("关注人的动态---->", string(i3))
	marshal, _ := json.Marshal(resp.GetTopicTypeList())
	fmt.Println("onetypes------>", string(marshal))
	i2, _ := json.Marshal(resp.GetPlayerList())
	fmt.Println("players------>", string(i2))
}

func TestGetSortTimeDynamicByTopicIdListFromDB() {
	data, count := for_game.GetSortTimeDynamicByTopicIdListFromDB(2, 3, []int64{3, 7, 8, 5})
	fmt.Println("count---------->", count)
	bytes, _ := json.Marshal(data)
	fmt.Println("----------->", string(bytes))
}

func TestUpdateDynamicSenderTypeToDB() {
	for_game.UpdateDynamicSenderTypeToDB(259, 1)
}

func TestGetDeviceHotDynamicByTopicIdFromDB() {
	data, count := for_game.GetDeviceHotDynamicByTopicIdFromDB(1, 2, 1)
	fmt.Println("count------>", count)
	bytes, _ := json.Marshal(data)
	fmt.Println("------------->", string(bytes))
}

func TestGetDevicePlayerHotDynamicByTopicIdFromDB() {
	data, count := for_game.GetDevicePlayerHotDynamicByTopicIdFromDB(2, 6, 1, 3)
	fmt.Println("count------>", count)
	bytes, _ := json.Marshal(data)
	fmt.Println("------------->", string(bytes))
}

func TestGetHighHotScoreTopicByTypeIdFromDB() {
	db, count := for_game.GetHighHotScoreTopicByTypeIdFromDB(1, 3, 2)
	bytes, _ := json.Marshal(db)
	fmt.Println("count---->", count)
	fmt.Println("---->", string(bytes))

}

func TestGetDynamicCommentInfoByPage() {
	sys := &share_message.SysParameter{
		CommentHotScore: easygo.NewInt32(1),
		CommentHotCount: easygo.NewInt32(3),
	}
	page := for_game.GetDynamicCommentInfoByPage(1213, 1, 10, nil, sys)
	bytes, _ := json.Marshal(page)
	fmt.Println("==========>", string(bytes))
}

func TestGetSecondaryCommentNum() {
	num := for_game.GetSecondaryCommentNum(1059)
	fmt.Println("count------->", num)
}

func TestWishAddress() {
	/*address := &h5_wish.WishAddress{
		Name:      easygo.NewString("ddd"),
		Detail:    easygo.NewString("cccc"),
		IfDefault: easygo.NewBool(false),
		Phone:     easygo.NewString("1321371234567"),
		AddressId: easygo.NewInt64(123456),
		Region:    easygo.NewString("广东省-广州市-天河区"),
	}
	err := wish.AddAddress(1887436008, address)
	if err != nil {
		fmt.Println("err", err)
	} else {
		fmt.Println("test ")
	}*/
	address := &h5_wish.WishAddress{
		Name:      easygo.NewString("ffff"),
		Detail:    easygo.NewString("xxxx"),
		IfDefault: easygo.NewBool(true),
		Phone:     easygo.NewString("1321327"),
		AddressId: easygo.NewInt64(1614222634),
		//Region:    easygo.NewString("广东省-广州市-天河区"),
	}
	err := wish.EditAddress(1887436008, address)
	if err != nil {
		fmt.Println("err", err)
	} else {
		fmt.Println("test ")
	}
	lst := wish.GetAddressListByUid(1887436008)
	for _, v := range lst {
		fmt.Println(v.GetName())
	}

}

/*func TestCollection() {

	lst, count,_ := wish.GetWishDataList(888888, 0*1, 1)

	for _, v := range lst {
		fmt.Println(v)
	}
	fmt.Println("count: ", count)

}*/

func TestGetRangeHotTopicFromDB() {
	topic := for_game.GetRangeHotTopicFromDB(1, nil)
	bytes, _ := json.Marshal(topic)
	fmt.Println("----->", string(bytes))
}

func TestLock() {
	itemLockKeys := []string{"redis_lock_666666"}
	errLock := easygo.RedisMgr.GetC().DoBatchRedisLockWithRetry(itemLockKeys, 10)
	errLock1 := easygo.RedisMgr.GetC().DoBatchRedisLockWithRetry(itemLockKeys, 10)
	if errLock != nil {
		fmt.Println("errLock------->", errLock.Error())
	}
	if errLock1 != nil {
		fmt.Println("errLock1------->", errLock1.Error())
	}
}

// 计算橘距离
func TestGetDistance() {
	dist := for_game.GetDistance(23.13953, 23.140620, 113.33969, 113.336020)
	fmt.Println("dist ----->", dist)
}

func TestGetNearInfoFromDB() {
	req := &client_hall.LocationInfoNewReq{
		X:   easygo.NewFloat64(113.336020),
		Y:   easygo.NewFloat64(23.140620),
		Sex: easygo.NewInt32(0),
	}
	pb := for_game.GetNearInfoFromDB([]int64{666}, for_game.ACCOUNT_TYPES_PT, req, 30)
	bytes, _ := json.Marshal(pb)
	fmt.Println("---------->", string(bytes))

}

func TestMakeRobot() {
	for_game.MakeRobot(20, 1)
}

type People struct {
	Name string
	Sex  int
}

func TestSetNearInfoToRedis() {
	p := make([]*People, 0)
	for i := 0; i < 3; i++ {
		p = append(p, &People{
			Name: "name" + easygo.AnytoA(i),
			Sex:  i,
		})
	}

	SetNearInfoToRedis1("aaabbbaaa", p, 30)
}

// expire 秒
func SetNearInfoToRedis1(key string, pl []*People, expire int64) {
	if isExist := for_game.ExistZAdd(key); isExist { // 存在的话先删除掉.
		for_game.DelZAdd(key)
	}
	m := make(map[int]string)
	for _, v := range pl {
		bytes, _ := json.Marshal(v)
		m[v.Sex] = string(bytes)
	}
	err1 := easygo.RedisMgr.GetC().ZAdd(key, m)
	easygo.PanicError(err1)
	err1 = easygo.RedisMgr.GetC().Expire(key, expire)
	easygo.PanicError(err1)
}

func TestExistZAdd() {
	fmt.Println("b-------->", for_game.ExistZAdd("aaabbbaaa"))
}
func TestDelZAdd() {
	for_game.DelZAdd("aaabbbaaa")
}

func TestGetNearByPageFromRedis() {
	page := 1
	pageSize := 2
	start, end := for_game.MakeRedisPage(page, pageSize, 3)
	fmt.Println(start, end)
	res := for_game.GetNearByPageFromRedis("aaabbbaaa", start, end)
	p := make([]*People, 0)
	marshal, _ := json.Marshal(res)
	json.Unmarshal(marshal, &p)
	for _, v := range res {
		fmt.Println("===========", string(v.([]byte)))
	}
	bytes, _ := json.Marshal(p)
	fmt.Println("=========>", string(bytes))
}

func GetGetHasPhotoHotDynamicToDB() {
	db := for_game.GetHasPhotoHotDynamicToDB(1887439539, 3)
	bytes, _ := json.Marshal(db)
	fmt.Println("---------->", string(bytes))
}

// 附近的人
func TestGetLocationInfoNew() {
	req := &client_hall.LocationInfoNewReq{
		X:        easygo.NewFloat64(121.497859),
		Y:        easygo.NewFloat64(31.247684),
		Sex:      easygo.NewInt32(0),
		Page:     easygo.NewInt64(1),
		PageSize: easygo.NewInt64(20),
		Sort:     easygo.NewInt32(for_game.NEAR_SORT_DISTANCE),
	}
	resp := for_game.GetLocationInfoNew(1887439538, true, req)
	bytes, _ := json.Marshal(resp)
	fmt.Println("1-------->", string(bytes))
	time.Sleep(500 * time.Millisecond)
	req = &client_hall.LocationInfoNewReq{
		X:        easygo.NewFloat64(121.497859),
		Y:        easygo.NewFloat64(31.247684),
		Sex:      easygo.NewInt32(0),
		Page:     easygo.NewInt64(2),
		PageSize: easygo.NewInt64(20),
		Sort:     easygo.NewInt32(for_game.NEAR_SORT_DISTANCE),
	}
	resp = for_game.GetLocationInfoNew(1887439538, false, req)
	bytes, _ = json.Marshal(resp)
	fmt.Println("2-------->", string(bytes))
	time.Sleep(500 * time.Millisecond)
	req = &client_hall.LocationInfoNewReq{
		X:        easygo.NewFloat64(121.497859),
		Y:        easygo.NewFloat64(31.247684),
		Sex:      easygo.NewInt32(0),
		Page:     easygo.NewInt64(3),
		PageSize: easygo.NewInt64(20),
		Sort:     easygo.NewInt32(for_game.NEAR_SORT_DISTANCE),
	}
	resp = for_game.GetLocationInfoNew(1887439538, false, req)
	bytes, _ = json.Marshal(resp)
	fmt.Println("3-------->", string(bytes))
	time.Sleep(500 * time.Millisecond)
	req = &client_hall.LocationInfoNewReq{
		X:        easygo.NewFloat64(121.497859),
		Y:        easygo.NewFloat64(31.247684),
		Sex:      easygo.NewInt32(0),
		Page:     easygo.NewInt64(4),
		PageSize: easygo.NewInt64(20),
		Sort:     easygo.NewInt32(for_game.NEAR_SORT_DISTANCE),
	}
	resp = for_game.GetLocationInfoNew(1887439538, false, req)
	bytes, _ = json.Marshal(resp)
	fmt.Println("4-------->", string(bytes))
	time.Sleep(500 * time.Millisecond)
	req = &client_hall.LocationInfoNewReq{
		X:        easygo.NewFloat64(121.497859),
		Y:        easygo.NewFloat64(31.247684),
		Sex:      easygo.NewInt32(0),
		Page:     easygo.NewInt64(5),
		PageSize: easygo.NewInt64(20),
		Sort:     easygo.NewInt32(for_game.NEAR_SORT_DISTANCE),
	}
	resp = for_game.GetLocationInfoNew(1887439538, false, req)
	bytes, _ = json.Marshal(resp)
	fmt.Println("5-------->", string(bytes))
}

func TestSaveNearLead() {
	for i := 1; i < 5; i++ {
		id := for_game.NextId(for_game.TABLE_NEAR_LEAD)
		reqMsg := &share_message.NearSet{
			Name: easygo.NewString(fmt.Sprintf("引导%d", i)),
			Icon: easygo.NewString(fmt.Sprintf("www.baidu%d.com", i)),
			//SkipUrl: easygo.NewString(fmt.Sprintf("www.baidu%d.com", i)),
			Weights: easygo.NewInt32(100 - i*10),
			Status:  easygo.NewInt32(1),
		}
		queryBson := bson.M{"_id": id}
		updateBson := bson.M{"$set": reqMsg}
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_NEAR_LEAD, queryBson, updateBson, true)
	}

}

func TestGetAllNearLeadFromDB() {
	db := for_game.GetAllNearLeadFromDB()
	bytes, _ := json.Marshal(db)
	fmt.Println("0----------->", string(bytes))
}
func TestInsertLeadData() {
	//for_game.InsertLeadData()
}
func TestGetNearRecommend() {
	bases, i := for_game.GetNearRecommend(1887439538, 1, 5, 113.336020, 23.140620)
	bytes, _ := json.Marshal(bases)
	fmt.Println("------->", i)
	fmt.Println("------->", string(bytes))
}

func TestAddSessionList() {
	for i := 100; i < 103; i++ {
		req := &share_message.NearSessionList{
			Content:         easygo.NewString(fmt.Sprintf("测试%d", i)),
			Status:          easygo.NewInt32(1),
			IsRead:          easygo.NewBool(false),
			CreateTime:      easygo.NewInt64(for_game.GetMillSecond()),
			ContentType:     easygo.NewInt32(1),
			ReceivePlayerId: easygo.NewInt64(1887439519),
		}
		for_game.AddSessionList(1887439539, 1887439519, req)
		time.Sleep(500 * time.Millisecond)
	}
}

func TestGetHaveUnRead() {
	//fmt.Println("count------->", for_game.GetHaveUnRead(1887439539))

}

func TestGetSessionFromDB() {
	db := for_game.GetSessionFromDB("1887439538_1887439539")
	bytes, _ := json.Marshal(db)
	fmt.Println("data-------->", string(bytes))

}

func TestUpdateSessionStatusByIds() {
	for_game.UpdateSessionStatusById("1887439538_1887439539")
}
func TestUpdateNearMessageStatusById() {
	for_game.UpdateNearMessageStatusById(1887439538, 1887439539)
}
func TestDelNearMessage() {
	for_game.DelNearMessage([]int64{1887439538, 1887439520}, 1887439539)
}

func TestGetMessageLogByIdFromDB() {
	db := for_game.GetMessageLogByIdFromDB(18)
	bytes, _ := json.Marshal(db)
	fmt.Println("-------->", string(bytes))
}

func TestUpdateIsRead() {
	for_game.UpdateIsRead(18)
}

func TestGetNearSessionListByPageFromDB() {
	lists, count := for_game.GetNearSessionListByPageFromDB(1887439539, 1, 10, 1607050674349)
	fmt.Println("count---------->", count)
	bytes, _ := json.Marshal(lists)
	fmt.Println("---------->", string(bytes))
}

func TestNearSessionList() {
	req := &client_hall.NearSessionListReq{
		Page:      easygo.NewInt64(1),
		PageSize:  easygo.NewInt64(10),
		QueryTime: nil,
	}
	sessionList := for_game.NearSessionList(1887439539, req)

	bytes, _ := json.Marshal(sessionList)
	fmt.Println("------->", string(bytes))
}

func TestPush() {
	ids := for_game.GetJGIds([]int64{1887439538})
	logs.Info("回收------>", ids)
	m := for_game.PushMessage{
		Title:       "过期提示",
		Content:     "您有平台赠送硬币明日0点即将过期",
		ContentType: for_game.JG_TYPE_BACKSTAGE_ASS,
		JumpObject:  3,
	}
	for_game.JGSendMessage(ids, m)
}

func TestUpdateComplaintNum() {
	for_game.UpdateComplaintNum(1215, 2)
}

func TestGetSignature() {
	fmt.Println(for_game.GetRandSignature())
}

func GetChatLog(logId int64) []*share_message.PersonalChatLog {
	col2, closeFun2 := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, "personal_chat_log")
	defer closeFun2()
	chatLogs := make([]*share_message.PersonalChatLog, 0)
	queryBson := bson.M{}
	if logId > 0 {
		queryBson["_id"] = bson.M{"$gt": logId}
	}
	err := col2.Find(queryBson).Sort("_id").Limit(5000).All(&chatLogs)
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

func GetPlayerName(id int64) string {
	col, closeFun := easygo.MongoMgr.GetDB(for_game.MONGODB_NINGMENG)
	defer closeFun()
	result := bson.M{}
	av := fmt.Sprintf("getPlayerName(%d)", id) //获取自增id的自定义函数
	logs.Info(av)
	err := col.Run(bson.M{"eval": av}, &result)
	easygo.PanicError(err)
	return easygo.AnytoA(result["retval"])
}

func TestFindWishBoxByNum() {
	boxes, e := for_game.FindWishBoxByNum(3)
	if e != nil {
		return
	}
	bytes, _ := json.Marshal(boxes)
	fmt.Println("---->", string(bytes))
}

func TestFindMaxPriceBoxItemByBoxId() {
	item, e := for_game.FindMaxPriceBoxItemByBoxId(1, nil)
	if e != nil {
		fmt.Println("---.", e.Error())
		return
	}
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestQueryWishItemById() {
	item, e := for_game.QueryWishItemById(1)
	if e != nil {
		fmt.Println("---.", e.Error())
		return
	}
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestProductShowService() {
	item, e := wish.ProductShowService(18801001)
	if e != nil {
		fmt.Println("---.", e.Error())
		return
	}
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetCoinService() {
	item, e := wish.GetCoinService()
	if e != nil {
		fmt.Println("---.", e.Error())
		return
	}
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestInsertPlayerWishItemToDb() {
	item := &share_message.PlayerWishItem{
		Id:              easygo.NewInt64(2),
		PlayerId:        easygo.NewInt64(1887440612),
		ChallengeItemId: easygo.NewInt64(2),
		CreateTime:      easygo.NewInt64(time.Now().Unix()),
		Status:          easygo.NewInt32(0),
		WishBoxId:       easygo.NewInt64(1),
	}
	e := for_game.InsertPlayerWishItemToDb(item)
	if e != nil {
		fmt.Println("---.", e.Error())
		return
	}
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetRandProducts() {

	item := for_game.GetRandProducts(1)

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetWishBoxItemByIdsFromDB() {

	item, err := for_game.GetWishBoxItemByIdsFromDB([]int64{13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27})
	if err != nil {
		return
	}

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetWishItemByIdsFromDB() {

	item, err := for_game.GetWishItemByIdsFromDB([]int64{1, 2})
	if err != nil {
		return
	}

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}

func TestGetPlayWishItem() {

	playerWishItem, err := for_game.GetPlayWishItemByPlayerId(1887436008, 0, 0, 20, "-CreateTime")
	if err != nil {
		return
	}

	fmt.Println("--------->", playerWishItem)

}

func TestTryOne() {
	err := wish.TryOnceService(1)
	if err != nil {
		panic(err)
	}
}

func TestMyWishService() {
	var req = h5_wish.MyWishReq{
		Page:     easygo.NewInt32(1),
		PageSize: easygo.NewInt32(10),
		Type:     easygo.NewInt32(1),
	}
	playerWishItem, err := wish.MyWishService(1887436008, &req)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range playerWishItem {
		fmt.Println(v)
		fmt.Println(v.GetBoxName())
	}

}
func TestProductInfoFromPlayerWishItem() {
	m := make(map[int64]int64)
	m[1] = 1887440612
	m[2] = 1887440612
	item, procuct := wish.ProductInfoFromPlayerWishItem(m)
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))
	bytess, _ := json.Marshal(procuct)
	fmt.Println("--------->", string(bytess))

}

func TestHomeMessageService() {
	item, err := wish.HomeMessageService()
	if err != nil {
		return
	}
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))
}

func TestGetRandWishItem() {
	item := for_game.GetRandWishItem(2)

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetRangWishLogByGuardian() {
	item := for_game.GetRangWishLogByGuardian(1887444994)

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestProtectorService() {
	item := wish.ProtectorService()

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetMenuList() {
	item := for_game.GetMenuList()

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestMenuService() {
	item, _ := wish.MenuService()

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}

func TestGetWishTypeList() {
	item, _ := for_game.GetNotHotWishTypeList([]int64{}, 0)

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestBrandListService() {
	item, _ := wish.BrandListService()
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))
}
func TestProductTypeListService() {
	item, _ := wish.ProductTypeListService()

	fmt.Println("--------->", item.GetTypeList())

}
func TestGetRangHotWishBrandByNum() {
	item := for_game.GetRangeHotWishBrandByNum(3)

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetRangWishBrandByNum() {
	item := for_game.GetRangeWishBrandByNum(10)

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestProductBrandListService() {
	item, _ := wish.ProductBrandListService()

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestDareRecommendService() {
	item, _ := wish.DareRecommendService()

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}

func TestRankingsService() {
	item, _ := wish.RankingsService()

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetAllWishLogByPid() {
	item := for_game.GetAllWishLogByPid(1887444994)

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestMyRecordService() {
	item, _ := wish.MyRecordService(1887444994)

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetAllWishLogByPage() {
	item, count := for_game.GetAllWishLogByPage(1887444994, 2, 1)

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))
	fmt.Println("count--------->", count)

}
func TestGetWishBoxsByIdsFromDB() {
	item := for_game.GetWishBoxsByIdsFromDB([]int64{1, 2})
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetAllWishItem() {
	item := for_game.GetAllWishItem()
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetAllWishBox() {
	item := for_game.GetAllWishBox()
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestQueryBoxProductNameService() {
	item, _ := wish.QueryBoxProductNameService()
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetRangeRecommendWishItem() {
	item := for_game.GetRangeRecommendWishItem(2)
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestSearchFoundService() {
	item, _ := wish.SearchFoundService()
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetWishBoxItemByItemId() {
	item := for_game.GetWishBoxItemByItemId(2)
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}

func TestGetRangeWishOccupied() {
	item := for_game.GetRangeWishOccupied(10)
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetOneSuccessDarer() {
	item := for_game.GetOneSuccessDarer(10)
	fmt.Println(item == nil)
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetRangeWishLogByBoxId() {
	item := for_game.GetRangeWishLogByBoxId(1, 10)
	fmt.Println(item == nil)
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestGetAllWishOccupiedByPage() {
	item, count := for_game.GetAllWishOccupiedByPage(1887444994, 1, 10)

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))
	fmt.Println("count--------->", count)

}
func TestFindMaxPriceBoxItemByBoxIdAndLv() {
	item := for_game.FindMaxPriceBoxItemByBoxIdAndLv(2, 1)
	rand := for_game.RandInt(0, len(item))

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))
	bytes1, _ := json.Marshal(item[rand])
	fmt.Println("--------->", string(bytes1))

}

func TestGetBoxListByPage() {
	req := &h5_wish.SearchBoxReq{
		Complex: easygo.NewInt32(0),
		//ProductStatus:  easygo.NewInt32(0),
		//MinPrice: easygo.NewInt64(1000),
		//MaxPrice:       nil,
		//WishBrandId:    easygo.NewInt64(2),
		//WishItemTypeId: easygo.NewInt64(1),
	}
	item, count := for_game.GetBoxListByPage(req, 1, 10)
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))
	fmt.Println("count--------->", count)

}
func TestQueryBoxService() {
	req := &h5_wish.QueryBoxReq{
		Id:   easygo.NewInt64(2),
		Type: easygo.NewInt32(2),
	}
	item, _ := wish.QueryBoxService(req)
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestMyDareService() {
	req := &h5_wish.MyDareReq{
		Page:     easygo.NewInt32(2),
		PageSize: easygo.NewInt32(2),
	}
	item, _ := wish.MyDareService(1887444994, req)

	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))

}
func TestSendMsgToIdelServer() {
	reqMsg := &h5_wish.AddCoinReq{
		UserId: easygo.NewInt64(1887444994),
		Coin:   easygo.NewInt64(-1),
	}
	resp, err := wish.SendMsgToIdelServer(for_game.SERVER_TYPE_HALL, "RpcAddCoin", reqMsg)
	if err != nil {
		logs.Error("err:", err.GetReason())
		return
	}

	bytes, _ := json.Marshal(resp)
	fmt.Println("--------->", string(bytes))

}

func TestInsertWishLog() {
	item := &share_message.WishLog{
		Id:              easygo.NewInt64(2),
		WishBoxId:       easygo.NewInt64(1),
		DareId:          easygo.NewInt64(1887436931),
		DareName:        easygo.NewString("70"),
		BeDareId:        easygo.NewInt64(1887437109),
		BeDareName:      easygo.NewString("木槿暖夏"),
		CreateTime:      easygo.NewInt64(time.Now().Unix()),
		Result:          easygo.NewBool(false),
		ChallengeItemId: easygo.NewInt64(1),
	}
	e := for_game.InsertWishLogToDB(item)
	if e != nil {
		fmt.Println("---.", e.Error())
		return
	}
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))
}
func TestAddOccupied() {
	item := &share_message.WishOccupied{
		Id:           nil,
		WishBoxId:    nil,
		NickName:     nil,
		HeadUrl:      nil,
		PlayerId:     nil,
		CreateTime:   nil,
		EndTime:      nil,
		OccupiedTime: nil,
		Status:       nil,
		CoinNum:      easygo.NewInt32(31),
	}
	e := for_game.AddOccupied(item)
	if e != nil {
		fmt.Println("---.", e.Error())
		return
	}
	bytes, _ := json.Marshal(item)
	fmt.Println("--------->", string(bytes))
}
func TestAddWishPoolPumpLog() {
	item := &share_message.WishPoolPumpLog{
		BoxId:  easygo.NewInt64(1),
		PoolId: easygo.NewInt64(1),
		Price:  easygo.NewInt64(20),
	}
	for_game.AddWishPoolPumpLog(item)

}
func TestAddWishPoolLog() {
	item := &share_message.WishPoolLog{
		BoxId:       easygo.NewInt64(1),
		PoolId:      easygo.NewInt64(1),
		PlayerId:    easygo.NewInt64(111),
		BeforeValue: easygo.NewInt64(10),
		AfterValue:  easygo.NewInt64(30),
		Value:       easygo.NewInt64(20),
		Type:        easygo.NewInt64(1),
	}
	for_game.AddWishPoolLog(item)

}

//func TestUpsetWishTopToDB() {
//	data := &share_message.WishTopLog{
//		Id:       easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_TOP_LOG)),
//		PlayerId: easygo.NewInt64(1887444994),
//		HeadIcon: easygo.NewString("https://im-resource-1253887233.cos.accelerate.myqcloud.com/defaulticon/girl_4.png"),
//
//		ThatDayTime: easygo.NewInt64(easygo.GetToday0ClockTimestamp()),
//		CreateTime:  easygo.NewInt64(time.Now().Unix()),
//	}
//	for_game.UpsetWishTopToDB(data, 1, 1)
//}
func TestSyncMap() {
	var sMap sync.Map
	sMap.Store(12, "张三")
	value, ok := sMap.Load(12)
	logs.Info("==========>", &sMap, value, ok)

}

func TestAddWishPool() {
	bigWin := &share_message.WishPoolStatus{
		MaxValue: easygo.NewInt64(100),
		MinValue: easygo.NewInt64(80),
	}
	smallWin := &share_message.WishPoolStatus{
		MaxValue: easygo.NewInt64(80),
		MinValue: easygo.NewInt64(60),
	}
	common := &share_message.WishPoolStatus{
		MaxValue: easygo.NewInt64(60),
		MinValue: easygo.NewInt64(40),
	}
	bigLoss := &share_message.WishPoolStatus{
		MaxValue: easygo.NewInt64(20),
		MinValue: easygo.NewInt64(0),
	}
	smallLoss := &share_message.WishPoolStatus{
		MaxValue: easygo.NewInt64(40),
		MinValue: easygo.NewInt64(20),
	}
	req := &share_message.WishPool{
		PoolLimit:    easygo.NewInt64(10000),
		InitialValue: easygo.NewInt64(5000),
		IncomeValue:  easygo.NewInt64(5000),
		Recycle:      easygo.NewInt64(8000), // 回收阀值
		Commission:   easygo.NewInt64(1000),
		StartAward:   easygo.NewInt64(7000),
		CloseAward:   easygo.NewInt64(4000),
		Name:         easygo.NewString("水池1号"),
		IsOpenAward:  easygo.NewBool(false),
		IsDefault:    easygo.NewBool(false),
		BigLoss:      bigLoss,
		SmallLoss:    smallLoss,
		Common:       common,
		BigWin:       bigWin,
		SmallWin:     smallWin,
	}
	for_game.AddWishPool(req)
}

// 抽水逻辑测试
//func TestChou(pid int, box *share_message.WishBox) {
//
//	//st := for_game.GetMillSecond()
//	productId := wish.RandWishBoxItem(int64(pid), box)
//	//et := for_game.GetMillSecond()
//	//logs.Debug("抽奖耗时为:", et-st)
//	logs.Debug("productId=--------------->:", productId)
//}

func TestWishService() {
	req := &h5_wish.WishReq{
		BoxId:     easygo.NewInt64(1),
		ProductId: easygo.NewInt64(1),
		OpType:    easygo.NewInt32(1),
	}
	wish.WishService(req, 1887440774)
}
func TestUpWishDataById() {
	UpPlayerWishData := &share_message.PlayerWishData{
		Status:     easygo.NewInt32(1),
		FinishTime: easygo.NewInt64(time.Now().Unix()),
	}
	for_game.UpWishDataById(108, UpPlayerWishData) //修改许愿
}
func TestGetPoolStatus() {
	fmt.Println("水池的状态为--->", wish.GetPoolStatus(1))
}

func TestBoxInfoService() {
	resp, _ := wish.BoxInfoService(666, 4, 2)
	bytes, _ := json.Marshal(resp)
	fmt.Println("---------->", string(bytes))
}
func TestProductDetailService() {
	resp, _ := wish.ProductDetailService(2)
	bytes, _ := json.Marshal(resp)
	fmt.Println("---------->", string(bytes))
}
func TestGetWishPlayerByPid() {
	resp := for_game.GetWishPlayerByPid(1887440774)
	bytes, _ := json.Marshal(resp)
	fmt.Println("---------->", string(bytes))
}
func TestCollectionBoxService() {
	req := &h5_wish.CollectionBoxReq{
		IdList: []int64{int64(2)},
		OpType: easygo.NewInt32(1),
	}
	wish.CollectionBoxService(222, req)

}

func TestGetPollInfoFromRedis() {
	for_game.GetPollInfoFromRedis(1)
}

func TestAddWishCoolDownConfigFromDB() {
	data := &share_message.WishCoolDownConfig{
		IsOpen:          easygo.NewBool(true),
		ContinuousTime:  easygo.NewInt64(60),
		ContinuousTimes: easygo.NewInt64(10),
		CoolDownTime:    easygo.NewInt64(30),
		DayLimit:        easygo.NewInt64(40),
	}
	for_game.AddWishCoolDownConfigFromDB(data)
}

func TestGetWishCoolDownConfigFromDB() {
	data := for_game.GetWishCoolDownConfigFromDB()
	bytes, _ := json.Marshal(data)
	fmt.Println("------->", string(bytes))
}

func TestGetCurrencyCfg() {
	cfg := for_game.GetCurrencyCfg()
	bytes, _ := json.Marshal(cfg)
	fmt.Println("--------->", string(bytes))
}
func TestCoinToMoney() {
	cfg, _ := wish.CoinToMoney(5)
	bytes, _ := json.Marshal(cfg)
	fmt.Println("--------->", string(bytes))
}
func TestGetIsGuardianWishBoxList() {
	result := for_game.GetIsGuardianWishBoxList()
	bytes, _ := json.Marshal(result)
	fmt.Println("--------->", string(bytes))
}

func TestInsertDiamondRecharge() {
	data := &share_message.DiamondRecharge{
		Id:         easygo.NewInt64(for_game.NextId("TABLE_WISH_DIAMOND_RECHARGE")),
		Diamond:    easygo.NewInt64(50),
		CoinPrice:  easygo.NewInt64(1000),
		MonthFirst: easygo.NewInt64(0),
		Rebate:     easygo.NewInt32(70),
		StartTime:  easygo.NewInt64(time.Now().Unix()),
		EndTime:    easygo.NewInt64(time.Now().Unix() + 3600),
		Status:     easygo.NewInt32(1),
		Sort:       easygo.NewInt32(7),
		DisPrice:   easygo.NewInt64(700),
	}
	for_game.InsertDiamondRecharge(data)
}

func TestGetDiamondRechargeList() {
	list := for_game.GetDiamondRechargeList()
	bytes, _ := json.Marshal(list)
	fmt.Println("------>", string(bytes))
}
func BatchDareCH() {
	reqMsg := &h5_wish.DoDareReq{
		DareType:  easygo.NewInt32(2),
		WishBoxId: easygo.NewInt64(83),
	}
	// 如果是挑战赛,如果没有许愿,不可以发起挑战
	wishData, _ := for_game.GetWishDataByStatus(18801002, reqMsg.GetWishBoxId(), for_game.WISH_CHALLENGE_WAIT)

	//if wishData.GetId() == 0 {
	//	logs.Error("用户是挑战赛,没有许愿,用户id为: %d,盲盒id为: %d", 18801001, reqMsg.GetWishBoxId())
	//	return
	//}
	wish.BatchDareCH(reqMsg, 18801002, wishData)
}

func TestGetWishGuardianCfg() {
	cfg := for_game.GetWishGuardianCfg()
	bytes, _ := json.Marshal(cfg)
	fmt.Println("----------->", string(bytes))
}
func TestGetGuardianCoinNumList() {
	list := for_game.GetGuardianCoinNumList(1619506193774, 1619506242494)
	bytes, _ := json.Marshal(list)
	fmt.Println("----------->", string(bytes))
}
func TestGetWishSumOccupied() {
	list, c := for_game.GetWishSumOccupied(1, 2, 1)
	bytes, _ := json.Marshal(list)
	fmt.Println("c----------->", c)
	fmt.Println("----------->", string(bytes))
}

func TestUpdateWishSumOccupied() {
	data := &share_message.WishSumOccupied{
		WishBoxId: easygo.NewInt64(2),
		NickName:  easygo.NewString("李四"),
		HeadUrl:   easygo.NewString("www.google.com"),
		PlayerId:  easygo.NewInt64(777),
	}
	for_game.UpdateWishSumOccupied(data, 10)
}

func TestInsertWishActPool() {
	data := &share_message.WishActPool{
		Id:         easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_ACT_POOL)),
		Name:       easygo.NewString("测试奖池2"),
		BoxNum:     easygo.NewInt32(3),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
		BoxIds:     []int64{48, 55, 59},
	}
	for_game.InsertWishActPool(data)
}

func TestInsertWishActPoolRule() {
	data := &share_message.WishActPoolRule{
		Id: easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_ACT_POOL_RULE)),
		//WishActPoolId: easygo.NewInt64(1),
		Key:        easygo.NewInt32(4), // 5天,100个砖石
		Diamond:    easygo.NewInt64(200),
		WishItemId: easygo.NewInt64(5),
		AwardType:  easygo.NewInt32(1),
		Type:       easygo.NewInt32(3),
	}
	for_game.InsertWishActPoolRule(data)
}

func TestGetWishActPoolByBoxId() {
	cfg := for_game.GetWishActPoolByBoxId(18)
	bytes, _ := json.Marshal(cfg)
	fmt.Println("----------->", string(bytes))
}
func TestGeWishActPoolRuleByPId() {
	cfg := for_game.GeWishActPoolRuleByPId(1)
	bytes, _ := json.Marshal(cfg)
	fmt.Println("----------->", string(bytes))
}

func TestUpsertWishPlayerActivity() {
	data := &share_message.WishPlayerActivity{
		PlayerId: easygo.NewInt64(666),

		UpdateTime: easygo.NewInt64(for_game.GetMillSecond()),
	}
	for_game.UpsertWishPlayerActivity(data)
}

func TestWishActService() {
	wish.WishActService(18877004, 18, 2094, 666)
}

func TestSumWeekMonth() {
	wish.SumWeekMonth(2)
}
func TestGetWishActPoolList() {
	list := for_game.GetWishActPoolList()
	bytes, _ := json.Marshal(list)
	fmt.Println("------->", string(bytes))
}
func TestGetWishPlayerActivityByPage() {
	list, c := for_game.GetWishPlayerActivityByPage(2, 1, 1)
	bytes, _ := json.Marshal(list)
	fmt.Println("count------->", c)
	fmt.Println("------->", string(bytes))
}
func TestGetWishActPoolRuleListByTypeKey() {
	page := 2
	pageSize := 2
	var key int
	if page > 1 {
		key = (page - 1) * pageSize
	} else {
		key = page - 1
	}
	fmt.Println("key------->", key)
	list := for_game.GetWishActPoolRuleListByTypeKey(3, key, pageSize)
	bytes, _ := json.Marshal(list)
	fmt.Println("------->", string(bytes))
}
func TestSumMoneyService() {
	reqMsg := &h5_wish.SumMoneyReq{
		DataType: easygo.NewInt64(2),
		Page:     easygo.NewInt32(2),
		PageSize: easygo.NewInt32(1),
	}
	result := wish.SumMoneyService(666, reqMsg)
	bytes, _ := json.Marshal(result)
	fmt.Println("------->", string(bytes))
}
func TestGetWishActivityPrizeLogById() {
	prizeLog := for_game.GetWishActivityPrizeLogById(15)
	bytes, _ := json.Marshal(prizeLog)
	fmt.Println("------->", string(bytes))
}
func TestGetWishPlayerActivityTop() {
	prizeLog := for_game.GetWishPlayerActivityTop(3, 2)
	bytes, _ := json.Marshal(prizeLog)
	fmt.Println("------->", string(bytes))
}
func TestGeWishActivityPrizeLogByIds() {
	prizeLog := for_game.GeWishActivityPrizeLogByIds([]int64{31, 32})
	bytes, _ := json.Marshal(prizeLog)
	fmt.Println("------->", string(bytes))
}

func TestGeWishActivityPrizeLogByPId() {
	prizeLog := for_game.GeWishActivityPrizeLogByPId(18877004, 5)
	bytes, _ := json.Marshal(prizeLog)
	fmt.Println("------->", string(bytes))
}

func TestGetActivityByTypes() {
	prizeLog := for_game.GetActivityByTypes([]int32{3, 4, 5, 6})
	bytes, _ := json.Marshal(prizeLog)
	fmt.Println("------->", string(bytes))
}
func TestFlushLocalInfo() {
	req := &client_hall.LocationInfoNewReq{
		X:          nil,
		Y:          nil,
		Sex:        easygo.NewInt32(2),
		Sort:       nil,
		Page:       nil,
		PageSize:   nil,
		IsNewFlush: nil,
	}
	prizeLog := for_game.FlushLocalInfo(1887570293, []int64{}, req)
	bytes, _ := json.Marshal(prizeLog)
	fmt.Println("------->", string(bytes))
}

func TestAddReportWishLogService() {
	wish.AddReportWishLogService(int64(123456), for_game.WISH_REPORT_ACCESS_WISH)
}

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

func TestUpdateLastWeekMonth() {
	wish.UpdateLastWeekMonth(4)
}

func TestGetWishPlayerActivityByNum() {
	num := for_game.GetWishPlayerActivityByNum(10, 4)
	bytes, _ := json.Marshal(num)
	fmt.Println("-----------> ", string(bytes))
}
func TestGetWishWhiteList() {
	num := for_game.GetWishWhiteList()
	bytes, _ := json.Marshal(num)
	fmt.Println("-----------> ", string(bytes))
}
func TestCheckIsFreeze() {
	num := wish.CheckIsFreeze(18940006)
	fmt.Println("-----------> ", num)
}

func TestInsertWishWhite() {
	data := &share_message.WishWhite{
		Id:       easygo.NewInt64(18801003),
		NickName: easygo.NewString("1003"),
		Account:  easygo.NewString(""),
		PlayerId: easygo.NewInt64(1887571808),
		Note:     easygo.NewString("1003"),
	}
	for_game.InsertWishWhite(data)
}
func main() {
	//BatchDareCH()
	//TestGetWishGuardianCfg()
	//TestProtectorService()

	//TestGetWishSumOccupied()
	//fmt.Println(easygo.IsTodayTimestamp(for_game.GetMillSecond() - 24*3600000))
	//TestInsertWishActPoolRule()
	//TestAddReportWishLogService()
	//TestGetWishActivityPrizeLogById()
	//TestUpdatePlayerActZeroByField()
	//TestGetWishPlayerActivityTop()
	//TestWishActService()
	//TestGetWishPlayerActivityByNum()
	//TestGetWishWhiteList()
	//fmt.Println(easygo.GetToday0ClockTimestamp())
	//TestGeWishActivityPrizeLogByPId()
	//TestProtectorService()
	TestCheckIsFreeze()
	time.Sleep(10 * time.Hour)
	/*sortWeightList := []int32{1671, 1965, 3715, 4774, 5570, 8356, 11143, 27872, 66973, 66973, 66973, 66973, 66973, 83758, 90568, 90568, 90568, 111771, 134215, 134215}

	var count int32
	for _, v := range sortWeightList {
		count += v
	}
	rate := make([]float32, 0)
	for _, key := range sortWeightList {
		// 计算百分比，保留两位小数
		f := easygo.Decimal(easygo.AtoFloat64(easygo.AnytoA(key))/easygo.AtoFloat64(easygo.AnytoA(count)), 6)
		rate = append(rate, float32(f))
	}
	logs.Info("rate------------->", rate)
	var count0, count1 int
	for b := 0; b < 100000; b++ {

		index := for_game.WeightedRandomIndex(rate)
		if index == 0 {
			count0++
		}
		if index == 1 {
			count1++
		}
		logs.Info("index------------->", index)
	}
	logs.Info("count0------------->", count0)
	logs.Info("count1------------->", count1)*/
	//===========================================================
	//keys, err := easygo.RedisMgr.GetC().Scan("chat_session") //存储所有关注的信息
	//easygo.PanicError(err)
	//logs.Info("查出keys:", len(keys), keys)
	//data := make([]string, 0)
	//for_game.InterfersToStrings(val, data)
	//easygo.PanicError(err)
	//var log []int64
	//for_game.InterfersToInt64s(val, &log)
	//logs.Info("-------->>>", keys)
	//for_game.InitPropsToMap()
	//fmt.Println(1606904702000 - 3*86400*1000)
	//TestPush()
	//GenPriSession()
	//GenTeamSession()
	//
	//name := GetPlayerName(1887444994)
	//fmt.Println("ame--------->", name)
	//TestInsertPlayerWishItemToDb()
	//TestBrandListService()

	//err := easygo.RedisMgr.GetC().Set("aaaa", "1")
	//easygo.PanicError(err)
	//err = easygo.RedisMgr.GetC().Set("bbbb", "2")
	//err = easygo.RedisMgr.GetC().Set("cccc", "3")
	//err = easygo.RedisMgr.GetC().Set("dddd", "3")
	//keys := []interface{}{"aaaa", "bbbb", "cccc", "dddd", "ffff"}
	//
	//logs.Info("keys:", keys, keys[0])
	//b, err := easygo.RedisMgr.GetC().Delete(keys...)
	//easygo.PanicError(err)
	//if b {
	//	logs.Info("删除成功")
	//} else {
	//	logs.Info("删除失败")
	//}
	//time.Sleep(time.Second * 10000)
	//
	//hall.PWebHuiChaoPay.QueryMerchantIn("386852991", "sweep-b90d814d534a4b219ab1fe0983f248e9", "WXZF")
	//hall.PWebHuiChaoPay.UpdateMerchantIn("sweep-b90d814d534a4b219ab1fe0983f248e9", "386852991", "WXZF")
	time.Sleep(10000 * time.Second)
	//hall.PWebHuiChaoPay.QueryMerchantIn("386852991", "sweep-b90d814d534a4b219ab1fe0983f248e9", "WXZF")
	////hall.PWebHuiChaoPay.UpdateMerchantIn("sweep-b90d814d534a4b219ab1fe0983f248e9", "386852991", "WXZF")
	//for_game.GetTeamLogDelPlayers(18801082, 1615617943130)

	//for_game.WeightedRandomIndex() 权重
	//for i := 0; i < 10; i++ {
	//	go func() { TestSyncMap() }()
	//}
	//TestAddWishPoolLog()

	//水池的数据
	//TestAddWishPool()
	//==================抽奖==============

	//==================抽奖==============
	//for_game.UpdatePollPriceToRedis(1, 0, "IncomeValue")

	//js := `{"_id":1,"PoolLimit":10000,"InitialValue":5000,"IncomeValue":5000,"Recycle":8000,"Commission":1000,"StartAward":7000,"CloseAward":4000,"Name":"水池1号","CreateTime":1614686459,"IsOpenAward":false,"IsDefault":false,"BigLoss":{"MaxValue":20,"MinValue":0,"AddWeight":444},"SmallLoss":{"MaxValue":40,"MinValue":20,"AddWeight":555},"Common":{"MaxValue":60,"MinValue":40,"AddWeight":666},"BigWin":{"MaxValue":100,"MinValue":80,"AddWeight":888},"SmallWin":{"MaxValue":80,"MinValue":60,"AddWeight":777}}`
	//var jj *for_game.WishPoolEX
	//json.Unmarshal([]byte(js), &jj)
	//logs.Info("==--==%+v", jj)

	//==========================

	//=======================

	//for_game.SetPlayerCoolDownTimeFromRedis(1, 30)
	//fmt.Println(for_game.IsExistPlayerCoolDownTimeFromRedis(1))
	/*	req := &h5_wish.SearchBoxReq{
			Complex: easygo.NewInt32(0),
			//ProductStatus:  easygo.NewInt32(0),
			//MinPrice: easygo.NewInt64(1000),
			//MaxPrice:       nil,
			//WishBrandId:    easygo.NewInt64(2),
			//WishItemTypeId: easygo.NewInt64(1),
		}
		ids, _ := wish.SearchBoxService(req)
		bytes, _ := json.Marshal(ids)
		fmt.Println("------>", string(bytes))*/
	/*	box, _ := for_game.GetWishBox(21)
		status := wish.CheckBoxStatus(box)
		logs.Error("err", status)*/

	//easygo.Spawn(wish.TaskUpBox)

	//t := for_game.GetMillSecond()
	/*	m := int(time.Unix(1613725510000/1000, 0).Month())
		fmt.Println("--------->", m)

		fmt.Println("--------->", int(time.Unix(for_game.GetMillSecond()/1000, 0).Month()))*/

	//fmt.Println("水池的状态为------------->", wish.GetPoolStatus(2))

	//GetGetHasPhotoHotDynamicToDB()
	//TestGetLocationInfoNew()
	//fmt.Println(for_game.GetDistance(22.55329, 22.55329, 113.88308, 113.88308))
	//i, i2 := SliceByPage(2, 2, 3)
	//fmt.Println(i, i2)

	//	for_game.RedisLuckyPlayer.IncrLuckyCountToRedis(1887439519, -1)
	//fmt.Println(for_game.RedisLuckyPlayer.GetLuckyCountFromRedis(1887439519))
	//lp := for_game.GetLuckyPlayerFromDB(1887439519)
	//fmt.Println(lp.GetLuckyCount())
	//li := for_game.GetSysPropsRateFromDB()

	//for_game.GetLocationInfo1()
	//for_game.LuckyCardD()
	// 判断今天是否分享过人了
	//period := for_game.GetPlayerPeriod(20201001)
	//period.DayPeriod.AddInteger("和", 1)
	//fmt.Println(period.DayPeriod.FetchInt("和"))
	//log.Printf("%+v", for_game.GetSysPropsByIdFromDB(1))
	//fmt.Println(for_game.GetPlayerByPhone("18504210016"))
	//for_game.ReloadSysPropsToDayProps()
	//fmt.Println(for_game.GetDayPropsById(1))
	//for_game.InitPropsToMap()
	//testUpsetPlayerPropsToDB()
	//for_game.LuckyCard(1887439519, 1)
	//fmt.Println(easygo.GetToday0ClockTimestamp())
	//fmt.Println(for_game.GetPlayerPropsByPidAndProIdFromDB(1887439519, 4))

	//err := for_game.GiveCard(1887439519, 1887439520, 6)
	//if err != nil {
	//	fmt.Println(err.GetReason())
	//	return
	//}
	//time.Sleep(6 * time.Second)
	//m := make(map[int64]int64)
	//m[time.Now().Unix()] = 1887439519
	//for_game.RedisLuckyPlayer.SetFullLuckyPlayerInfoToRedis(m)
	//for_game.UpsetLuckyPlayerToDB(1887439519, 0, time.Now().Unix(), true)

	//s := for_game.GetLogListByPid(1887439519, 2)
	//bytes, _ := json.Marshal(s)
	//fmt.Println(string(bytes))

	//fmt.Println(easygo.GetToday0ClockTimestamp())
	//for_game.UpsetSysFullCount(1)
	//period := for_game.GetPlayerPeriod(20201001)
	//period.HaltYearPeriod.AddInteger("FullCount", 3)
	//fmt.Println(period.HaltYearPeriod.FetchInt("FullCount"))

	//for_game.RedisLuckyPlayer.SetFullCountToRedis(1)
	//fmt.Println(for_game.RedisLuckyPlayer.GetFullCountFromRedis())

	//adv := &share_message.AdvSetting{
	//	Id:         easygo.NewInt64(3),
	//	Title:      easygo.NewString("社交广场动态低权重"),
	//	Types:      easygo.NewInt32(1),
	//	Location:   easygo.NewInt32(1),
	//	Status:     easygo.NewInt32(1),
	//	CreateTime: easygo.NewInt64(time.Now().Unix()),
	//	StartTime:  easygo.NewInt64(for_game.GetMillSecond()),
	//	EndTime:    easygo.NewInt64(1602729911000),
	//	JumpUrl:    easygo.NewString("www.baidu1.com"),
	//	TxtSource:  easygo.NewString("社交1"),
	//	Weights:    easygo.NewInt32(45),
	//}
	//for_game.UpdateAdvListToDB(adv) // 设置广告数据
	//adv2 := &share_message.AdvSetting{
	//	Id:         easygo.NewInt64(4),
	//	Title:      easygo.NewString("社交广场动态低权重"),
	//	Types:      easygo.NewInt32(1),
	//	Location:   easygo.NewInt32(1),
	//	Status:     easygo.NewInt32(1),
	//	CreateTime: easygo.NewInt64(time.Now().Unix()),
	//	StartTime:  easygo.NewInt64(for_game.GetMillSecond()),
	//	EndTime:    easygo.NewInt64(1602729911000),
	//	JumpUrl:    easygo.NewString("www.google1.com"),
	//	TxtSource:  easygo.NewString("社交1"),
	//	Weights:    easygo.NewInt32(30),
	//}
	//for_game.UpdateAdvListToDB(adv2) // 设置广告数据

	//db := for_game.GetADVSetMap(1)
	//ss, ok := db[1]
	//bytes, _ := json.Marshal(ss)
	//fmt.Println("----->", len(db))
	//fmt.Println(for_game.GetAllTrueZan(1887439519))

	/*	// 绑定关联关系
		pr := &share_message.LuckyPlayerRelated{
			Id:          easygo.NewString(fmt.Sprintf("%s_%s", easygo.AnytoA(1887439538), "110114")),
			PlayerId:    easygo.NewInt64(1887439538),
			FriendPhone: easygo.NewString("110114"),
			RelatedTime: easygo.NewInt64(time.Now().Unix()),
		}
		if err := for_game.UpsetLuckyPlayerRelatedToDB(pr); err != nil {
			easygo.PanicError(err)
		}
	*/
	/*	data := &share_message.DynamicData{
			LogId:    easygo.NewInt64(714),
			PlayerId: easygo.NewInt64(1887439519),
		}

		data2 := &share_message.DynamicData{
			LogId:    easygo.NewInt64(711),
			PlayerId: easygo.NewInt64(1887439519),
		}
		data3 := &share_message.DynamicData{
			LogId:    easygo.NewInt64(710),
			PlayerId: easygo.NewInt64(1887439519),
		}
		bytes, _ := json.Marshal(data)
		bytes2, _ := json.Marshal(data2)
		bytes3, _ := json.Marshal(data3)
		m := make(map[int64]string)
		m[data.GetLogId()] = string(bytes)
		m[data2.GetLogId()] = string(bytes2)
		m[data3.GetLogId()] = string(bytes3)
		err := easygo.RedisMgr.GetC().HMSet("redis_square_bs_top_dynamic", m)
		easygo.PanicError(err)*/
	//for_game.UpsetLuckyPlayerToDB(101, 2, 0, true)
	//for_game.UpsetIsNewLuckyToDB(101, true)

	//data, i := for_game.GetNoTopDynamicListByPageFromDB(2, 30)
	//data, i := for_game.GetNoTopDynamicByPIDsFromDB(1, 30, []int64{1887439519})
	//data := for_game.GetBSTopDynamicListByIDsFromDB([]int64{1887439519, 1887439520})
	//data := for_game.GetAppTopDynamicListByIDsFromDB([]int64{1887439519, 1887439520})
	//bytes, _ := json.Marshal(data)
	//fmt.Println("------>", string(bytes))
	//fmt.Println("------>", i)

	//data1 := for_game.GetDynamicSliceByRandFromSlice(data, 2)
	//bytes1, _ := json.Marshal(data1)
	//fmt.Println("------>", string(bytes1))
	//
	//data2 := for_game.SortDynamicSliceByTime1(data1)
	//bytes2, _ := json.Marshal(data2)
	//fmt.Println("------>", string(bytes2))

	/*	dynamic1 := for_game.GetRedisNewDynamic1(2, 1887439519, 1, 10)
		bytes, _ := json.Marshal(dynamic1)
		fmt.Println("--->", string(bytes))*/

	//commentList, err := for_game.GetCommentListByPlayerIdFromDB(1887439550)
	//easygo.PanicError(err)
	//fmt.Println("------>", len(commentList))

	//dynamic1 := for_game.GetRedisNewDynamic1(1, 1, 20, 1887439519, 1)
	//bytes, _ := json.Marshal(dynamic1)
	//fmt.Println("---------->", string(bytes))

	//fmt.Println(for_game.GetAdvLogByPidAndOpFromDB(1887436049, 1) != nil)
	//keys, fromRedis := for_game.GetSomeZanInfoFromRedis()
	//fmt.Println(keys)
	//bytes, _ := json.Marshal(fromRedis)
	//fmt.Println(string(bytes))
	//for_game.UpsetAllZanDataToDB(fromRedis)
	//for_game.DelZanDatasFromRedis(keys)
	//fmt.Println(for_game.DelZanDataFromRedisAndDB(811, 35))

	//data := for_game.GetDynamicByStatusSFromDB(835, []int{for_game.DYNAMIC_STATUE_COMMON, for_game.DYNAMIC_STATUE_UNPUBLISHED})
	//fmt.Println(data.GetLogId())
	//fmt.Println(1601296140 - 1601297753)

	//for_game.SetFirstPageMaxLogIdToRedis("123123123", "456456")
	//fmt.Println(for_game.GetFirstPageMaxLogIdFromRedis("1887439539_1887439538"))
	//notifyTopReq := &share_message.BackstageNotifyTopReq{
	//	LogId:       easygo.NewInt64(666),
	//	TopOverTime: easygo.NewInt64(300),
	//}
	//// 通知广场服务器
	//hall.BroadCastMsgToOtherSquare("RpcHallNotifyTop", notifyTopReq)
	//fmt.Println("=============>", for_game.DelSquareDynamicById(999888))
	//fmt.Println(for_game.SetRechargeCountToRedis(666, 3, 20)) //3 8
	//time.Sleep(1 * time.Second)
	//fmt.Println(for_game.SetRechargeCountToRedis(777, 3, 20)) // 3 7
	//time.Sleep(1 * time.Second)
	//fmt.Println(for_game.SetRechargeCountToRedis(666, 2, 20)) //5 10
	//time.Sleep(1 * time.Second)
	//fmt.Println(for_game.SetRechargeCountToRedis(777, 1, 20)) // 4 8
	//time.Sleep(1 * time.Second)

	/*	hall.PSysParameterMgr = for_game.NewSysParameterManager()
		for i := 0; i < 1; i++ {
			time.Sleep(1000 * time.Millisecond)
			hall.CheckWarningSMS(2, 101)
		}
	*/
	//hall.CheckWarningSMS(666, 2, 10)
	//TestProductShowService()
	time.Sleep(10000000)
}

func TestGetAttention() {
	ids := for_game.GetAttentionPlayers(1, 1, 10, for_game.VC_ATTENTION_TO_OTHER)
	fmt.Println(ids)
}

func TestAddAttention() {
	for i := 9; i < 15; i++ {
		for_game.AddAttentionLog(int64(i), int64(i+1), for_game.VC_ATTENTION_HI)
		time.Sleep(1000)
	}
	//for_game.AddAttentionLog(3, 4, for_game.VC_ATTIENTION_TO_OTHER)
}

// 内部方法,抽卡,返回卡片id
func luckyCard() int64 {
	// 获取权重
	var v []*share_message.Props
	for_game.SysPropsData.RateMap.Range(func(key, value interface{}) bool {
		if easygo.GetToday0ClockTimestamp() == key.(int64) {
			v = value.([]*share_message.Props)
		}
		return true
	})
	if len(v) == 0 {
		logs.Error("luckyCard 获取权重出错,可能时间有问题,当天的整点时间为:", easygo.GetToday0ClockTimestamp())
		return -1
	}
	rate := make([]float32, 0)
	for _, value := range v {
		rate = append(rate, float32(value.GetRate())/100)
	}
	logs.Info("luckyCard 权重列表--->", rate)
	// 一万次
	var index int
	for i := 0; i < 10000; i++ {
		index = for_game.WeightedRandomIndex(rate)
	}

	id := for_game.SysPropsData.PropsSlice[index] // 抽到的卡的id
	id = 4
	if id == for_game.ID_QU { // 只有趣字才有控制
		if count := for_game.IncrDayProps(id, -1); count < 0 {
			//return luckyCard()
			return 2
		}
	}
	logs.Info("luckyCard 卡片索引为----->", index)
	logs.Info("luckyCard 卡片id为----->", id)
	return id
}

// 从数组中随机选取指定条数的数据,返回新的数组
func GetSliceByRandFromSlice(s []int64, count int) []int64 {
	result := make([]int64, 0)
	if len(s) == 0 {
		return result
	}
	if len(s) == count {
		return s
	}
	//rand.Seed(time.Now().Unix())
	rand.Seed(util.GetMilliTime())
	for i := 0; i < len(s); i++ {
		data := s[rand.Intn(len(s))]
		// 判断不在result中,就记录进去
		if !util.Int64InSlice(data, result) {
			result = append(result, data)
		}
		if len(result) == count { // 指定获取的长度
			break
		}
	}
	// 数据源不足指定条,结果不等于源数据
	if (len(s) < count) && (len(result) < len(s)) {
		return GetSliceByRandFromSlice(s, count)
	}

	if (len(s) > count) && (len(result) < count) {
		return GetSliceByRandFromSlice(s, count)
	}

	return result
}

func Insert() {
	i := A()
	var charuIndex int64
	for in := 0; in < len(i); in++ {
		rand := for_game.RangeRand(5, 10) // 得出广告的插入位置
		fmt.Println(rand)
		charuIndex += int64(in) + rand
		//fmt.Println("in--------->", in)
		if int(charuIndex) > len(i) {
			fmt.Println("============> 超过最大索引了")
			break
		}
		i = easygo.Insert(i, int(charuIndex), 10000+in).([]int)
		fmt.Println("charuIndex --------->", charuIndex)
	}
	fmt.Println("i --------->", i)
}

//pr := &share_message.LuckyPlayerRelated{
//	Id:          easygo.NewString(fmt.Sprintf("%s_%s", easygo.AnytoA(1887439519), easygo.AnytoA(530))),
//	PlayerId:    easygo.NewInt64(1887439519),
//	FriendPhone: easygo.NewString("530"),
//	RelatedTime: easygo.NewInt64(time.Now().Unix()),
//}
//if err := for_game.UpsetLuckyPlayerRelatedToDB(pr); err != nil {
//	easygo.PanicError(err)
//}

func A() []int {
	i := make([]int, 0)
	for in := 1; in <= 20; in++ {
		i = append(i, in)
	}
	fmt.Println("----------------", for_game.GetRelatedByPhoneFromDB("550"))
	return i
}

//testRedis()
//opt := redis.DialPassword("redis2020")
//conn, err := redis.Dial("tcp", "127.0.0.1:6379", opt)
//if err != nil {
//	fmt.Println("Connect to redis error", err)
//	return
//}
//defer conn.Close()
//ac := &for_game.RedisPlayerAccount{
//	PlayerId:    1887436004,
//	Account:     "601",
//	Email:       "",
//	Password:    "8d70d8ab2768f232ebe874175065ead3",
//	Token:       "O71VLlLcU3N5A7lS",
//	PayPassword: "",
//	CreateTime:  1585016972749,
//}
//res, err := conn.Do("hmset", redis.Args{"struct1"}.AddFlat(ac)...)
//easygo.PanicError(err)
//logs.Info("res:", res)
//value1, err := redis.Values(conn.Do("hgetall", "struct1"))
//logs.Info("value:", value1)
//object := &for_game.RedisPlayerAccount{}
//redis.ScanStruct(value1, object)
//logs.Info("最后：", object)
//testRedis()
//changePersonnalChat()
//changeTeamChat()
//ch1 := make(chan int,0)
//ch2 := make(chan int)
//go func(){
//	for i := 1;i < 11;i++{
//		ch1 <- i
//	}
//	fmt.Println("func1发送完毕")
//	close(ch1)
//}()
//
//go func(){
//	for{
//		v,ok := <- ch1
//		if !ok{
//			break
//		}
//		ch2 <- v * v
//	}
//	close(ch2)
//}()
//
//for v := range ch2{
//	fmt.Println(v)
//}
//ch1 <- 5
//fmt.Println("发送成功")
//testhuichao()
//TestTransferFixed()
//fmt.Println(for_game.AuthBankIdName("6214633131067889708", "黄家茵", "44178119941003022X", "13168180383"))
//logs.Error("RpcAddBank 增加银行卡,验证平台发送的短信验证码验证失败,手机号为: %s,验证码为: %s,短信类型为: %d",
//	//	"who.GetPhone()", "reqMsg.GetMsgCode()", for_game.CLIENT_CODE_BINDBANK)
//TestAliPayQuery()
//testhuichao()

//TestYemadaiPayQuery()
//testhuichao()
//TestM()
//Wechat()
//testLuckySignIn(1887439519)
//testGetPlayerPropsList(1887439519)
//testUpsetPlayerPropsToDB()
//testUpsetLuckyPlayerToDB()
//}

//func main() {
//	//logs.Info(easygo.Stamp2StrExt(1595592380943))
//
//	//testRedis()
//	//opt := redis.DialPassword("redis2020")
//	//conn, err := redis.Dial("tcp", "127.0.0.1:6379", opt)
//	//if err != nil {
//	//	fmt.Println("Connect to redis error", err)
//	//	return
//	//}
//	//defer conn.Close()
//	//ac := &for_game.RedisPlayerAccount{
//	//	PlayerId:    1887436004,
//	//	Account:     "601",
//	//	Email:       "",
//	//	Password:    "8d70d8ab2768f232ebe874175065ead3",
//	//	Token:       "O71VLlLcU3N5A7lS",
//	//	PayPassword: "",
//	//	CreateTime:  1585016972749,
//	//}
//	//res, err := conn.Do("hmset", redis.Args{"struct1"}.AddFlat(ac)...)
//	//easygo.PanicError(err)
//	//logs.Info("res:", res)
//	//value1, err := redis.Values(conn.Do("hgetall", "struct1"))
//	//logs.Info("value:", value1)
//	//object := &for_game.RedisPlayerAccount{}
//	//redis.ScanStruct(value1, object)
//	//logs.Info("最后：", object)
//	//testRedis()
//	//changePersonnalChat()
//	//changeTeamChat()
//	//ch1 := make(chan int,0)
//	//ch2 := make(chan int)
//	//go func(){
//	//	for i := 1;i < 11;i++{
//	//		ch1 <- i
//	//	}
//	//	fmt.Println("func1发送完毕")
//	//	close(ch1)
//	//}()
//	//
//	//go func(){
//	//	for{
//	//		v,ok := <- ch1
//	//		if !ok{
//	//			break
//	//		}
//	//		ch2 <- v * v
//	//	}
//	//	close(ch2)
//	//}()
//	//
//	//for v := range ch2{
//	//	fmt.Println(v)
//	//}
//	//ch1 <- 5
//	//fmt.Println("发送成功")
//	//testhuichao()
//	//TestTransferFixed()
//	//fmt.Println(for_game.AuthBankIdName("6214633131067889708", "黄家茵", "44178119941003022X", "13168180383"))
//	//logs.Error("RpcAddBank 增加银行卡,验证平台发送的短信验证码验证失败,手机号为: %s,验证码为: %s,短信类型为: %d",
//	//	//	"who.GetPhone()", "reqMsg.GetMsgCode()", for_game.CLIENT_CODE_BINDBANK)
//	//TestAliPayQuery()
//	//testhuichao()
//
//	//TestYemadaiPayQuery()
//	//testhuichao()
//	//TestM()
//	Wechat()
//}

func testhuichao() {
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
	hall.Initialize()

	//启动http web服务器
	//payData := &share_message.PayOrderInfo{
	//	PlayerId:    easygo.NewInt64(1887436001),
	//	ProduceName: easygo.NewString("测试信银联支付"),
	//	Amount:      easygo.NewString("0.01"),
	//	PayId:       easygo.NewInt32(for_game.PAY_CHANNEL_HUICHAO_YL),
	//	PayType:     easygo.NewInt32(3),
	//	PaySence:    easygo.NewInt32(5),
	//	PayWay:      easygo.NewInt32(1),
	//	PayBankNo:   easygo.NewString("6212262011022171826"),
	//}
	//hall.WebServreMgr = hall.NewWebHttpServer()
	//hall.WebServreMgr.HuiChaoEntry(nil, nil)
	//hall.PWXApiMgr.WXAPILogin(hall.WebServreMgr, "0432oOD91qG1aN1ld2F9152wD912oODv", payData)
	//aa := hall.NewWebHuiChaoPay()
	//data := "aa123"
	//s := aa.RsaEncode(for_game.TestPrivate, data)
	//logs.Info("明文:", data)
	//logs.Info("密文:", s)
	//
	//bs := base64.StdEncoding.EncodeToString([]byte(data))
	//s1 := aa.RsaEncode(for_game.TestPrivate, bs)
	//logs.Info("明文base64后加密后的密文:", s1)
	//player := hall.GetPlayerObj(1887436001)
	//if player != nil {
	//rsp, _ := hall.PWebHuiChaoPay.ReqPaySMSApi(payData, player)
	//if rsp != nil {
	//	logs.Info("orderNo:", rsp.GetString("mch_order_no"))
	//}
	//req := &client_hall.BankPaySMS{
	//	SMS:     easygo.NewString("807855"),
	//	OrderNo: easygo.NewString("1101159536251251463293"),
	//}
	//hall.PWebHuiChaoPay.ReqSMSPayApi(req)
	//hall.PWebHuiChaoPay.ReqCheckPayOrder("1101159536251251463293")
	//}

	//hall.PWebHuiChaoPay.QueryMerchantIn("386852991", "sweep-b90d814d534a4b219ab1fe0983f248e9", "WXZF")
	//hall.PWebHuiChaoPay.UpdateMerchantIn("sweep-b90d814d534a4b219ab1fe0983f248e9", "385996806", "WXZF")

	//fmt.Println(for_game.GetRedisSquareLogIds())
	//
	//logs.Info(easygo.RedisMgr.GetC().HGetAll(for_game.MakeNewString(for_game.REDIS_SQUARE_ATTENTION, 1887437594)))
	//for_game.SendMessageCodeEx("+8613570213647", "1234", true, false)
	//fmt.Println("------------", for_game.SendWarningSMS([]string{"+8613570213647"}, "1234"))
}

func send(ch chan int) {
	ch <- 10
	fmt.Println("发送成功")
}

//迁移个人聊天记录
//func changePersonnalChat() {
//	initializer := hall.NewInitializer()
//	defer func() { // 若是异常了,确保异步日志有成功写盘
//		logger := initializer.GetBeeLogger()
//		if logger != nil {
//			logger.Flush()
//		}
//	}()
//	dict := easygo.KWAT{
//		"logName":  "hall",
//		"yamlPath": "config_hall.yaml",
//	}
//	initializer.Execute(dict)
//	hall.Initialize()
//	//先查出id_generator
//	col, closeFun := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_PERSON_LOG, for_game.TABLE_ID_GENERATOR)
//	defer closeFun()
//	val := []easygo.KWAT{}
//	err := col.Find(bson.M{}).All(&val)
//	easygo.PanicError(err)
//	i := int64(1)
//	count := 0
//	col2, closeFun2 := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, "personal_chat_log")
//	defer closeFun2()
//	for _, v := range val {
//		k := v.GetString("_id")
//		va := v.GetString("Value")
//		logs.Info("k,v:", k, va)
//		data := make([]*share_message.PersonalChatLog, easygo.Atoi(va))
//		col1, closeFun1 := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_PERSON_LOG, k)
//		err = col1.Find(bson.M{}).All(&data)
//		closeFun1()
//		easygo.PanicError(err)
//		logs.Info("data:", data)
//		var saveData []interface{}
//		for j := 0; j < len(data); j++ {
//			data[j].LogId = easygo.NewInt64(i)
//			saveData = append(saveData, data[j])
//			i += 1
//
//		}
//		err = col2.Insert(saveData...)
//		easygo.PanicError(err)
//		count += easygo.Atoi(va)
//	}
//	logs.Info("总记录条数:", count)
//}

//迁移群聊天记录
//func changeTeamChat() {
//	initializer := hall.NewInitializer()
//	defer func() { // 若是异常了,确保异步日志有成功写盘
//		logger := initializer.GetBeeLogger()
//		if logger != nil {
//			logger.Flush()
//		}
//	}()
//	dict := easygo.KWAT{
//		"logName":  "hall",
//		"yamlPath": "config_hall.yaml",
//	}
//	initializer.Execute(dict)
//	hall.Initialize()
//	//先查出id_generator
//	col, closeFun := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_TEAM_LOG, for_game.TABLE_ID_GENERATOR)
//	defer closeFun()
//	val := []easygo.KWAT{}
//	err := col.Find(bson.M{}).All(&val)
//	easygo.PanicError(err)
//	i := int64(1)
//	count := 0
//	col2, closeFun2 := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, "team_chat_log")
//	defer closeFun2()
//	for _, v := range val {
//		k := v.GetString("_id")
//		va := v.GetString("Value")
//		logs.Info("k,v:", k, va)
//		data := make([]*share_message.TeamChatLog, easygo.Atoi(va))
//		col1, closeFun1 := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_TEAM_LOG, k)
//		err = col1.Find(bson.M{}).All(&data)
//		closeFun1()
//		easygo.PanicError(err)
//		logs.Info("data:", data)
//		var saveData []interface{}
//		for j := 0; j < len(data); j++ {
//			data[j].TeamLogId = easygo.NewInt64(data[j].LogId)
//			data[j].LogId = easygo.NewInt64(i)
//			saveData = append(saveData, data[j])
//			i += 1
//
//		}
//		err = col2.Insert(saveData...)
//		easygo.PanicError(err)
//		count += easygo.Atoi(va)
//	}
//	logs.Info("总记录条数:", count)
//	identity := &for_game.Identity{
//		Key:   easygo.NewString("team_chat_log"),
//		Value: easygo.NewUint64(count),
//	}
//	col3, closeFun3 := hall.MongoLogMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ID_GENERATOR)
//	defer closeFun3()
//	col3.Insert(&identity)
//}

func testRedis() {
	initializer := hall.NewInitializer()
	defer func() { // 若是异常了,确保异步日志有成功写盘
		logger := initializer.GetBeeLogger()
		if logger != nil {
			logger.Flush()
		}
	}()
	dict := easygo.KWAT{
		"logName":  "login",
		"yamlPath": "config_login.yaml",
	}
	initializer.Execute(dict)
	login.Initialize()
	player := for_game.GetPlayerById(1887436001)
	logs.Info(player)
	js, err := json.Marshal(player)
	easygo.PanicError(err)
	logs.Info("js:", string(js))
	psData := easygo.KWAT{}
	err = json.Unmarshal(js, &psData)
	//a := make(map[string]interface{})
	err = easygo.RedisMgr.GetC().HMSet("mynew_test", psData)
	easygo.PanicError(err)
	all, err := redis.Values(easygo.RedisMgr.GetC().HGetAll("mynew_test"))
	easygo.PanicError(err)
	logs.Info("all:", all)
	//newData := map[string]string{}
	newData, err := redis.StringMap(all, err)
	easygo.PanicError(err)
	logs.Info("newData:", newData["PlayerSetting"])
	//js1, err := json.Marshal(newData)

	//err = json.Unmarshal(js1, base)
	//logs.Info(base)
	//a := reflect.ValueOf(psData["PlayerSetting"])
	//logs.Info("a:", a)
	//t := reflect.TypeOf(a)
	//logs.Info("t:", t.Kind())
	//logs.Info("a:", a.Kind())
	//logs.Info("type:", reflect.TypeOf(a), a["IsNewMessage"])
	//b, ok := a.(*easygo.KWAT)
	//if !ok {
	//	logs.Info("err")
	//}
	//logs.Info("psData:", a["IsNewMessage"])
	//for_game.StatisticsRedPacket()

	//log := share_message.GoldChangeLog{
	//	LogId:      easygo.NewInt64(1),
	//	PlayerId:   easygo.NewInt64(111),
	//	ChangeGold: easygo.NewInt64(100),
	//	SourceType: easygo.NewInt32(1),
	//	PayType:    easygo.NewInt32(2),
	//	Note:       easygo.NewString(""),
	//
	//	CurGold:    easygo.NewInt64(10),
	//	Gold:       easygo.NewInt64(90),
	//	CreateTime: easygo.NewInt64(1111111111111111),
	//}
	//
	//extend := &share_message.RechargeExtend{
	//	Channeltype: easygo.NewInt32(1),
	//	PayChannel:  easygo.NewInt32(2),
	//}
	//
	//m := for_game.CommonGold{
	//	GoldChangeLog: log,
	//	Extend:        extend,
	//}
	//chatInfo := make(map[int64]string)
	//s, _ := json.Marshal(m)
	//chatInfo[log.GetLogId()] = string(s)
	//err1 := easygo.RedisMgr.GetC().HMSet("11111111111111", chatInfo)
	//easygo.PanicError(err1)
	//
	//var m1 for_game.CommonGold
	//b, err2 := easygo.RedisMgr.GetC().HGet("11111111111111", "1")
	//easygo.PanicError(err2)
	//err3 := json.Unmarshal(b, &m1)
	//easygo.PanicError(err3)
	//
	//m2 := &for_game.GoldChangeLog{}
	//for_game.StructToOtherStruct(m1, m2)
	//logs.Info(m2, m2.Extend)
	//for name, value := range m2 {
	//	v := easygo.NewInt32(value)
	//	logs.Info(*v)
	//	logs.Info(name, value, reflect.TypeOf(name), reflect.TypeOf(value))
	//}
	//err := easygo.RedisMgr.GetC().SAdd("111111", 1)
	//easygo.PanicError(err)
	//
	//err1 := easygo.RedisMgr.GetC().SAdd("111111", 1)
	//easygo.PanicError(err1)
	//
	//err2 := easygo.RedisMgr.GetC().SAdd("111111", 2)
	//easygo.PanicError(err2)
	//
	//err3 := easygo.RedisMgr.GetC().SAdd("111111", 3)
	//easygo.PanicError(err3)
	//
	//ids := []int64{}
	//value, err := easygo.RedisMgr.GetC().Smembers("111111")
	//easygo.PanicError(err)
	//for_game.InterfersToInt64s(value, &ids)
	//logs.Info(ids)

	//_, err4 := easygo.RedisMgr.GetC().Delete("111111")
	//easygo.PanicError(err4)
	//
	//b, err5 := easygo.RedisMgr.GetC().Exist("111111")
	//easygo.PanicError(err5)
	//logs.Info(b)
	//s := ""
	//b, _ := json.Marshal(s)
	//logs.Info(b)
	//var s1 string
	//err := json.Unmarshal(b, &s1)
	//logs.Info(err, s1)

	//err := easygo.RedisMgr.GetC().StringSet("1", "111")
	//easygo.PanicError(err)
	//
	//var s string
	//err1 := easygo.RedisMgr.GetC().StringGet("1", &s)
	//easygo.PanicError(err1)
	//logs.Info(s)
	//for_game.InitRedisIds(easygo.MongoMgr, for_game.MONGODB_NINGMENG)
	//logs.Info(for_game.GetRedisCurrentId(for_game.REDIS_TEAM_ID))
	//info := make(map[int64]string)
	//info[1] = "1"
	//info[2] = "2"
	//info[3] = "3"
	//info[4] = "4"
	//
	//err := easygo.RedisMgr.GetC().HMSet("2", info)
	//easygo.PanicError(err)

	//type people struct {
	//	Name  string
	//	Age   int32
	//	High  int32
	//	Weigh int32
	//}
	//
	//p := &people{Name: "fnnnnnn", Age: 18, High: 180, Weigh: 120}
	//s, _ := json.Marshal(p)
	//logs.Info(s, string(s))
	//err1 := easygo.RedisMgr.GetC().HMSet("4", p)
	//easygo.PanicError(err1)
	////
	//err2 := easygo.RedisMgr.GetC().HSet("4", "Age", 18)
	//easygo.PanicError(err2)
	//
	//v, err3 := easygo.RedisMgr.GetC().HGet("4", "Age")
	//easygo.PanicError(err3)
	//logs.Info(v, string(v))
	//for _, v1 := range v {
	//	logs.Info(string(v1.([]byte)))
	//}
	//logs.Info(p1)

	//v, err1 := easygo.RedisMgr.GetC().HMGet("2", "1")
	//easygo.PanicError(err1)
	//acct := 0
	//_, err := redis.Scan(v, &acct)
	//easygo.PanicError(err)
	//logs.Info(acct)

	//var lst []*share_message.TeamChatLog
	//logsId := []string{"1"}
	//values, err := easygo.RedisMgr.GetC().HMGet(for_game.MakeNewString(for_game.REDIS_TABLE_TEAMCHAT, "18800004"), logsId...)
	//easygo.PanicError(err)
	//for _, m := range values {
	//	log := &share_message.TeamChatLog{}
	//	json.Unmarshal(m.([]byte), &log)
	//	lst = append(lst, log)
	//}
	//logs.Info(lst)

	//team := &for_game.RedisTeam{
	//	TeamId: 1111111,
	//}
	//err := easygo.RedisMgr.GetC().HSet(for_game.REDIS_TABLE_TEAM, easygo.AnytoA(1), info)
	//easygo.PanicError(err)
	//
	//value, err1 := easygo.RedisMgr.GetC().HGet(for_game.REDIS_TABLE_TEAM, easygo.AnytoA(1))
	//easygo.PanicError(err1)
	//info1 := make(map[int64]string)
	//
	//logs.Info(value)
	//value, err := easygo.RedisMgr.GetC().HKeys(for_game.REDIS_TABLE_TEAM)
	//easygo.PanicError(err)
	//for _, v := range value {
	//	logs.Info(string(v.([]byte)))
	//}

	//logs.Info(v)
	//player := hall.GetPlayerObj(1887436745)
	//for_game.NewRedisPlayerBase(&player.PlayerBase.PlayerBase.PlayerBase)
	//base := for_game.GetRedisPlayerBase(1887436745)
	//logs.Info("playerbase:", base)
	//setting := for_game.GetRedisPlayerSetting(1887436745)
	//logs.Info("setting:", setting)
	//photo := for_game.GetredisPlayerPhoto(1887436745)
	//logs.Info("photo:", photo)
	//teamids := for_game.GetredisPlayerTeamIds(1887436745)
	//logs.Info("teamids:", teamids)
	//banksInfo := for_game.GetredisBankInfo(1887436745)
	//logs.Info("banksInfo:", banksInfo)
	//blackList := for_game.GetredisBlackList(1887436745)
	//logs.Info("blackList:", blackList)
	//collects := for_game.GetredisCollectInfo(1887436745)
	//logs.Info("collects:", collects)
}

func testRedisAccount() {
	//account := &share_message.PlayerAccount{
	//	PlayerId:    easygo.NewInt64(1887436001),
	//	Account:     easygo.NewString("601"),
	//	Email:       easygo.NewString(""),
	//	PayPassword: easygo.NewString("8d70d8ab2768f232ebe874175065ead3"),
	//	Password:    easygo.NewString("O71VLlLcU3N5A7lS"),
	//	Token:       easygo.NewString("O71VLlLcU3N5A7lS"),
	//	CreateTime:  easygo.NewInt64(1585016972749),
	//}
	//player := for_game.GetPlayerAccount(1887436001)
	////player := for_game.GetAccountByPhone("601")
	//logs.Info("player:", player)
	//for_game.AddPlayerAccount(&player.PlayerAccount)
	//res := for_game.GetAccountByPhone(player.GetAccount())
	//logs.Info("res:", res)
	//b := for_game.CheckPlayerByAccount(player.GetAccount())
	//logs.Info("b:", b, player.GetAccount())

}

// 测试代付
func TestTransferFixed() {
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
	hall.Initialize()

	hc := hall.NewWebHuiChaoPay()
	//hc.TransferFixed()
	hc.TransferQueryFixed("2201159563214689749109", 1*time.Second)
	//hc.CheckBalance()
}

func TestAliPayQuery() {
	appID := "2016102700770128"
	privateKey := `MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCN8VHpXP1ZSJSUZeQtZZNDQrHqr88nn7A+NmeKrsuu5PkQ2CvOc0deBLbL42Znz1PqaVK3ioHDOCTJHrX+xBsmDfhqE0s2QUoJDWrhqGw+HRQ8anlgUQ2C6UYKsII26jDYPFywtB86ZrykOZ8Cd+su0CmGjnE9htCOdjEql5jmmmTWdz8FZePLN7XNja8LbumVTG5TBzOornOm5Cz3IqrEqlJMMS+mO7yxpasb2JOXgrt+lyE3MRAv7wS9pCPgtwC0TveLN6FYe/lvb8kHnq9ROq0phBJlMdHpZ79CT41PpQoNTAxVLN5bFmbdO6syjgn0pHRLvLIBaiLprVdqRnSfAgMBAAECggEAeebyjhSKkI9A62HGYSaHHpC88+0hX8pJNmTK79PGoeGL9edxV9CxThGGW/xkCmuIih0CKRcO8nXZQdDaRH5vQnNlENSZF3Ni/ftD+6EFtSKMKobWzt1NWUy2FqAYdMkUQeE1SZyn5SQuhmvmH9yVYpLr1t+maUzK+E6RUx729bQCb2OMkfP8tZtfcMLWiFwk6U+p9zo7WlIzxJqpRyu2KodQtlo94Qsx0X0EIcHkR72MqdTOCdC7KcWd04FMSW9NPfLEejykbpKiF9T8HR1mK2dC0em5eSTLsYcknduWUgiBD3S2gM9qWFIDcWFChcV5yTs73iNULUHcRYJuMvWjYQKBgQDQ92MuQFz4CEFLretNSx2v8W22cwhIFBkGBYSLo1yAbV2UxOuSmi9JVkEP6My3iETzDFGh9tXAQxbfNqFEcwQ8t423maNvhqybU47nJn/F84uxNuFZOH9fK14RD9QSRpWtwnTTZoOGPlR7U7+qLqFLnEICLoUVR8UP8tp9avr3dQKBgQCt5AlCP79DApWs20Vqoh9HQqjnw+LRXHMpFTfnw5F9Qa3g7N5v94osKbZE7oCdqrIuzjboEfNnrKd/x6y6TSuYIquECYQwv+4ci/RLhOjTYnY/NCErq83qoGfuUxWGEkoX1KE7yB56OeuC6MdJ9Q+Ms4LJ6NLh/CaavqXu6dPNQwKBgET1TnKF5Ogo+Ts7Moo4PpzAJD9wKIx4rWVSTtIx36W18YrVjRO889vUrfXNEjmCq5Y1O38iUJl4ykRw57kJ550NyaOL/OYh4DYF1gOrrcCqRS/+91CVF1tVmV4yBf7d8ij8IcddbgvP59sm4PoNF0c3UoUbyukh3QMNVlLLCfS9AoGBAJQQ5EFg/n8UqFYzr3wI6BFJlYEjrvMOgZCt3JigUjYRwvkPOKimYyUPr4AqhaG7Q1XPibk578SLo2SOpWlNZJ16iAk6ATFxfFMaaL4VQhsccAuJW+VPuVrbkyO/40fyMtzv1QqOcEUrJHqns2oqHT91axx5/3cluclyJOC2gf75AoGAChJrZtH2lYsB9XFLHQnr1YYe/AFWmDNNWYIySvoNdduxbamdmFMtqXQDtu9G328EcZE74p6RdhXRlJe8hkwxPBDdR6VZfMvlEjpKtSnzQWwG7tYE5Zn5NP5LJ/YPPenwQMXM099kwFo9Z3YxMjM7o6D72wu27JGiw19BhkPoFuQ=`
	client := alipay.NewClient(appID, privateKey, false)

	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("trade_no", "2020071822001461011420004215")
	bm.Set("query_options", "TRADE_SETTLE_INFO")
	query, err := client.TradeQuery(bm)
	if err != nil {
		logs.Info("err------->", err.Error())

		return
	}
	logs.Info("------------>%+v", query)
}

func TestYemadaiPayQuery() {
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
	hall.Initialize()

	//hc := hall.NewWebHuiChaoPay()

	//hc.ReqCheckPayOrder("1101159711174990835432")
}

func TestM() {
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
	hall.Initialize()

	var arr []int64
	for i := int64(1); i <= int64(41); i++ {
		arr = append(arr, i)
	}
	logs.Info(for_game.Paginator(4, 10, arr))
}

func Wechat() {

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
	hall.Initialize()

	//for_game.GetWechatUserInfo("37_3HcdQ3JLVdTdrXAvtf2kBss405sAOI7ZMOeZaLQ-RxwfBXnUX71TvtOP8ydwXNxnS0rHx7PUo79le7ULnXQFIqgrKqEch2qNJ2igizn_cD8", "otIOqwXtEGzXjUiW1c99mqMN04r8")
}

//func testLuckySignIn(playerId int64) {
//	if err := for_game.LuckySignIn(playerId); err != nil {
//		logs.Info(err.String())
//	}
//}

func testGetPlayerPropsList(pid int64) {
	props, message := for_game.GetPlayerPropsList(pid)
	if message != nil {
		logs.Info("----------->", message.String())
		return
	}
	logs.Info("============>", props)
}

func testUpsetPlayerPropsToDB() {

	for_game.UpsetPlayerPropsToDB(1887439519, 1, 1)
}

func testUpsetLuckyPlayerToDB() {
	for_game.UpsetLuckyPlayerToDB(1887439519, 3, 0)
}

func testUpsetLuckyPlayerRelatedToDB() {
	pr := &share_message.LuckyPlayerRelated{
		PlayerId:    easygo.NewInt64(1887439519),
		FriendPhone: easygo.NewString(1887439520),
		RelatedTime: easygo.NewInt64(time.Now().Unix()),
	}
	for_game.UpsetLuckyPlayerRelatedToDB(pr)
}
func testIncBgVoiceTag() {
	for_game.IncBgVoiceTag([]int32{1, 2, 3})
}
func testGetVoiceTags() {
	top, later := for_game.GetVoiceTags(int32(0), 3, -1)
	for _, v := range top {
		logs.Info("id: %v, Name: %v, IsHot: %v", v.GetId(), v.GetName(), true)

	}
	for _, v := range later {
		logs.Info("id: %v, Name: %v, IsHot: %v", v.GetId(), v.GetName(), false)
	}
}

type TestLable struct {
	Id              int64
	Label           []int64
	PersonalityTags []int64
}

func TestInsertLable() {
	saveData := make([]interface{}, 0)
	for i := 0; i < 10000; i++ {
		lable := make([]int64, 0)
		start := for_game.RandInt(1, 6)
		time.Sleep(10 * time.Millisecond)
		end := for_game.RandInt(2, 6)
		lable = append(lable, int64(start))
		if start != end {
			lable = append(lable, int64(end))
		}
		personalityTags := make([]int64, 0)
		start1 := for_game.RandInt(1, 6)
		time.Sleep(10 * time.Millisecond)
		end2 := for_game.RandInt(2, 6)
		time.Sleep(10 * time.Millisecond)
		end3 := for_game.RandInt(4, 9)
		personalityTags = append(personalityTags, int64(start1))
		if start1 != end2 {
			personalityTags = append(personalityTags, int64(end2))
		}
		if end3 != start1 && end3 != end2 {
			personalityTags = append(personalityTags, int64(end3))
		}
		t := &TestLable{
			Id:              int64(i),
			Label:           lable,
			PersonalityTags: personalityTags,
		}
		saveData = append(saveData, bson.M{"_id": t.Id}, t)
	}

	if len(saveData) > 0 {
		for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, "test_lable", saveData)
	}

}

//// 内部方法,抽卡,返回卡片id
//func luckyCard() int64 {
//	// 获取权重
//	var v []*share_message.Props
//	for_game.SysPropsData.RateMap.Range(func(key, value interface{}) bool {
//		if easygo.GetToday0ClockTimestamp() == key.(int64) {
//			v = value.([]*share_message.Props)
//		}
//		return true
//	})
//	if len(v) == 0 {
//		logs.Error("luckyCard 获取权重出错,可能时间有问题,当天的整点时间为:", easygo.GetToday0ClockTimestamp())
//		return -1
//	}
//	rate := make([]float32, 0)
//	for _, value := range v {
//		rate = append(rate, float32(value.GetRate())/100)
//	}
//	logs.Info("luckyCard 权重列表--->", rate)
//	// 一万次
//	var index int
//	for i := 0; i < 10000; i++ {
//		index = for_game.WeightedRandomIndex(rate)
//	}
//
//	id := for_game.SysPropsData.PropsSlice[index] // 抽到的卡的id
//	id = 2
//	if id == for_game.ID_QU { // 只有趣字才有控制
//		if count := for_game.IncrDayProps(id, -1); count < 0 {
//			return luckyCard()
//		}
//	}
//	logs.Info("luckyCard 卡片索引为----->", index)
//	logs.Info("luckyCard 卡片id为----->", id)
//	return id
//}
