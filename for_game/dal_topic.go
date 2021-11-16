/**
话题相关的db操作
*/
package for_game

import (
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"time"

	"github.com/akqp2019/mgo/bson"
)

const (
	TOPIC_CLASS_BS   = 1 // 官方添加的话题的类别
	TOPIC_CLASS_USER = 2 // 用户自定义的话题类别.
)

const (
	TOPIC_STATUS_OPEN  = 1 // 开启
	TOPIC_STATUS_CLOSE = 2 // 关闭
)
const (
	DYNAMIC_TOPIC_NUM = 5 // 动态添加话题时,最多5个话题.
)
const (
	TOPIC_MAIN_PAGE_SIZE = 6 // 话题主页为6条动态.
)

// 话题头部话题列表
const (
	TOPIC_HEAD_TOPIC_RECOMD_NUM = 1 // 推荐话题1条
	TOPIC_HEAD_TOPIC_HOT_NUM    = 1 // 热门的话题1条
	TOPIC_HEAD_TOPIC_ALL_NUM    = 5 // 总共5条话题,减去上面2条剩下的就是普通,普通的话题 3条
)

const (
	DynamicInc = 5
	CommentInc = 1
	LikeInc    = 1
	DynamicMax = 20
	CommentMax = 15
	LikeMax    = 15
	DevoteMin  = 10
)

/**
操作话题的操作数
num:参与数,正数表示添加,负数表示减少.
*/
func OperateTopicParticipationNum(content string, num int64) {
	if content == "" {
		return
	}
	topics := CheckTopic(content)
	//取前5个话题
	if len(topics) > DYNAMIC_TOPIC_NUM {
		topics = topics[:DYNAMIC_TOPIC_NUM]
	}
	for _, topicName := range topics {
		topic := GetTopicByNameFromDB(topicName)
		if topic == nil {
			continue
		}
		// 给该话题加参与数
		IncTopicParticipationNumToDB(topic.GetId(), num)
	}
}

/**
操作话题的浏览数
num:参与数,正数表示添加,负数表示减少.
*/
func OperateTopicViewNum(content string, num int64) {
	if content == "" {
		return
	}
	topics := CheckTopic(content)
	//取前5个话题
	if len(topics) > DYNAMIC_TOPIC_NUM {
		topics = topics[:DYNAMIC_TOPIC_NUM]
	}
	for _, topicName := range topics {
		topic := GetTopicByNameFromDB(topicName)
		if topic == nil {
			continue
		}
		// 给该话题加参与数
		IncTopicViewingNumToDB(topic.GetId(), num)
	}
}

/**
解析动态内容,抽取话题
content: 动态内容.
num:参与数,正数表示添加,负数表示减少.
*/
func ParseContentTopic(content string, num int64) []int64 {
	// 1.获取话题 2.遍历查找话题是否存在.3,插入话题.添加参与分数.
	if content == "" {
		return nil
	}
	topics := CheckTopic(content)
	topicIdList := make([]int64, 0)
	//取前5个话题
	if len(topics) > DYNAMIC_TOPIC_NUM {
		topics = topics[:DYNAMIC_TOPIC_NUM]
	}
	for _, topicName := range topics {
		if topic := GetTopicByNameNoStatusFromDB(topicName); topic != nil {
			topicIdList = append(topicIdList, topic.GetId())
			// 给该话题加参与数
			IncTopicParticipationNumToDB(topic.GetId(), num)
		} else { // 新增话题
			// 新增话题类别,自定义,..,id为0
			topicType := &share_message.TopicType{
				Id:         easygo.NewInt64(0),
				Name:       easygo.NewString("自定义"),
				CreateTime: easygo.NewInt64(GetMillSecond()),
				TopicClass: easygo.NewInt32(TOPIC_CLASS_USER), // 用户自定义
				Status:     easygo.NewInt32(TOPIC_STATUS_OPEN),
				UpdateTime: nil,
			}
			UpsetTopicTypeToDB(topicType)
			tpcId := NextId(TABLE_TOPIC)
			tpc := &share_message.Topic{
				Id:                  easygo.NewInt64(tpcId),
				TopicTypeId:         easygo.NewInt64(topicType.GetId()),
				Name:                easygo.NewString(topicName),
				TopicClass:          easygo.NewInt32(TOPIC_CLASS_USER),
				CreateTime:          easygo.NewInt64(GetMillSecond()),
				Status:              easygo.NewInt32(TOPIC_STATUS_OPEN),
				ParticipationNum:    easygo.NewInt64(1),
				FansNum:             easygo.NewInt64(0),
				ViewingNum:          easygo.NewInt64(0),
				AddViewingNum:       easygo.NewInt64(0),
				AddParticipationNum: easygo.NewInt64(0),
				AddFansNum:          easygo.NewInt64(0),
				IsOpen:              easygo.NewBool(true),
			}
			UpsetTopicToDB(tpc)
			topicIdList = append(topicIdList, tpcId)
		}
	}
	return topicIdList
}

// 新增或修改话题
func UpsetTopicToDB(topic *share_message.Topic) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": topic.GetId()}, bson.M{"$set": topic})
	easygo.PanicError(err)
}

// 新增或修改话题类型
func UpsetTopicTypeToDB(t *share_message.TopicType) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_TYPE)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": t.GetId()}, bson.M{"$set": t})
	easygo.PanicError(err)
}

// 根据话题名字获取话题id 和类别id
func GetTopicByNameFromDB(name string) *share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	var tp *share_message.Topic
	err := col.Find(bson.M{"Name": name, "Status": TOPIC_STATUS_OPEN}).Select(bson.M{"_id": 1, "TopicTypeId": 1, "Name": 1, "TopicMaster": 1}).One(&tp)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return tp
}

// 根据话题名字获取话题id 和类别id
func GetTopicByNameNoStatusFromDB(name string) *share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	var tp *share_message.Topic
	err := col.Find(bson.M{"Name": name}).Select(bson.M{"_id": 1, "TopicTypeId": 1, "Name": 1, "TopicMaster": 1}).One(&tp)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return tp
}

// 根据话题类型查找是官方还是个人自定义的话题类别
func GetBSTopicTypeListByClassFormDB(topicClass int) []*share_message.TopicType {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_TYPE)
	defer closeFun()
	topicType := make([]*share_message.TopicType, 0)
	err := col.Find(bson.M{"TopicClass": topicClass, "Status": TOPIC_STATUS_OPEN}).Sort("-Sort").All(&topicType)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return topicType
}

// 根据话题类别查找话题列表
func GetBSTopicTypeListByTypeIdFormDB(topicTypeId int64) []*share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	topics := make([]*share_message.Topic, 0)
	err := col.Find(bson.M{"TopicTypeId": topicTypeId, "Status": TOPIC_STATUS_OPEN}).All(&topics)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	return topics
}

// 根据类别查询话题,分页查询.
func GetBSTopicTypeListByTypeIdPageFormDB(topicTypeId int64, page, pageSize int) ([]*share_message.Topic, int) {
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	list := make([]*share_message.Topic, 0)

	queryBson := bson.M{}
	queryBson["TopicTypeId"] = topicTypeId
	queryBson["Status"] = TOPIC_STATUS_OPEN
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	err1 := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(err1)

	return list, count
}

// 根据话题id获取话题.
func GetTopicByIdFromDB(id int64) *share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	var topic *share_message.Topic
	err := col.Find(bson.M{"_id": id, "Status": TOPIC_STATUS_OPEN}).One(&topic)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	return topic
}

// 根据话题id获取话题,不使用状态
func GetTopicByIdNoStatusFromDB(id int64) *share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	var topic *share_message.Topic
	err := col.Find(bson.M{"_id": id}).One(&topic)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	return topic
}

func InsertPlayerAttentionToDB(pat *share_message.PlayerAttentionTopic) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_PLAYER_ATTENTION)
	defer closeFun()
	err := col.Insert(pat)
	easygo.PanicError(err)
}

func DelPlayerAttentionFromDB(pid, tid int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_PLAYER_ATTENTION)
	defer closeFun()
	err := col.Remove(bson.M{"PlayerId": pid, "TopicId": tid})
	easygo.PanicError(err)
}

// 根据话题id获取话题.
func GetPlayerAttentionFromDB(id, pid int64) *share_message.PlayerAttentionTopic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_PLAYER_ATTENTION)
	defer closeFun()
	var tp *share_message.PlayerAttentionTopic
	err := col.Find(bson.M{"PlayerId": pid, "TopicId": id}).One(&tp)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return tp
}

// 递增粉丝数
func UpdateTopicFansToDB(id int64, fans int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	err := col.Update(bson.M{"_id": id}, bson.M{"$inc": bson.M{"FansNum": fans}})
	easygo.PanicError(err)
}

// 递增参与数
func IncTopicParticipationNumToDB(id int64, num int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	err := col.Update(bson.M{"_id": id}, bson.M{"$inc": bson.M{"ParticipationNum": num}})
	easygo.PanicError(err)
}

// 递增浏览数
func IncTopicViewingNumToDB(id int64, num int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	err := col.Update(bson.M{"_id": id}, bson.M{"$inc": bson.M{"ViewingNum": num}})
	easygo.PanicError(err)
}

// 递增话题分数值
func IncTopicHotScoreNumToDB(id int64, num int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	err := col.Update(bson.M{"_id": id}, bson.M{"$inc": bson.M{"HotScore": num}})
	easygo.PanicError(err)
}

// 分页获取玩家关注的列表
func GetTopicPlayerAttentionListByPidFromDB(pid int64, page, pageSize int) ([]*share_message.PlayerAttentionTopic, int) {
	list := make([]*share_message.PlayerAttentionTopic, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_PLAYER_ATTENTION)
	defer closeFun()
	queryBson := bson.M{}
	queryBson["PlayerId"] = pid

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	err1 := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(err1)
	return list, count
}

// 根据类别随机选取n条话题
func GetRangeTopicByTypeIdFromDB(topicTypeId int64, topicIds []int64, num int) []*share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	list := make([]*share_message.Topic, 0)
	var queryBson bson.M
	queryBson = bson.M{"TopicTypeId": topicTypeId, "Status": TOPIC_STATUS_OPEN}
	if len(topicIds) > 0 {
		queryBson = bson.M{"TopicTypeId": topicTypeId, "Status": TOPIC_STATUS_OPEN, "_id": bson.M{"$nin": topicIds}}
	}

	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

// 根据类别随机选取n条话题
func GetRangeCommentTopicByTypeIdFromDB(topicTypeId int64, topicIds []int64, num int) []*share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	list := make([]*share_message.Topic, 0)
	var queryBson bson.M
	queryBson = bson.M{"TopicTypeId": topicTypeId, "Status": TOPIC_STATUS_OPEN, "IsRecommend": bson.M{"$ne": true}, "IsHot": bson.M{"$ne": true}}
	if len(topicIds) > 0 {
		queryBson = bson.M{"TopicTypeId": topicTypeId, "Status": TOPIC_STATUS_OPEN, "IsRecommend": bson.M{"$ne": true}, "IsHot": bson.M{"$ne": true}, "_id": bson.M{"$nin": topicIds}}
	}

	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

// 根据类别随机选取n条热门话题
func GetRangeHotTopicByTypeIdFromDB(topicTypeId int64, topicIds []int64, num int) []*share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	list := make([]*share_message.Topic, 0)
	var queryBson bson.M
	queryBson = bson.M{"TopicTypeId": topicTypeId, "Status": TOPIC_STATUS_OPEN, "IsHot": true}
	if len(topicIds) > 0 {
		queryBson = bson.M{"TopicTypeId": topicTypeId, "Status": TOPIC_STATUS_OPEN, "IsHot": true, "_id": bson.M{"$nin": topicIds}}
	}

	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

// 根据类别随机选取n条推荐话题
func GetRangeRecommendTopicByTypeIdFromDB(topicTypeId int64, topicIds []int64, num int) []*share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	list := make([]*share_message.Topic, 0)
	var queryBson bson.M
	queryBson = bson.M{"IsRecommend": true, "TopicTypeId": topicTypeId, "Status": TOPIC_STATUS_OPEN}
	if len(topicIds) > 0 {
		queryBson = bson.M{"IsRecommend": true, "TopicTypeId": topicTypeId, "Status": TOPIC_STATUS_OPEN, "_id": bson.M{"$nin": topicIds}}
	}

	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

// 模糊搜索话题.
func LikeSearchTopicByNameFromDB(name string) []*share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	topicList := make([]*share_message.Topic, 0)
	queryBson := bson.M{
		"Name":   bson.M{"$regex": name},
		"Status": TOPIC_STATUS_OPEN,
	}
	if err := col.Find(queryBson).All(&topicList); err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	return topicList
}

// 查询热门分数最高的话题
func GetHighHotScoreTopicFromDB(hotScore int64) []*share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	list := make([]*share_message.Topic, 0)

	queryBson := bson.M{"Status": TOPIC_STATUS_OPEN}
	m := []bson.M{
		{"$match": queryBson},
		{"$addFields": bson.M{"Sort": bson.M{"$add": []string{"$ViewingNum", "$FansNum", "$ParticipationNum", "$AddViewingNum", "$AddParticipationNum", "$AddFansNum"}}}},
		{"$match": bson.M{"Sort": bson.M{"$gte": hotScore}}}, // 达标的热门分
		{"$sort": bson.M{"Sort": -1, "CreateTime": 1}},
		{"$project": bson.M{"Sort": 0}},
		{"$limit": 5},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

// 根据话题类别,话题的总分分页查询话题
func GetHighHotScoreTopicByTypeIdFromDB(typeId int64, page, pageSize int) ([]*share_message.Topic, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	list := make([]*share_message.Topic, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)
	queryBson := bson.M{"Status": TOPIC_STATUS_OPEN, "TopicTypeId": typeId}
	m := []bson.M{
		{"$match": queryBson},
		{"$addFields": bson.M{"Sort": bson.M{"$add": []string{"$ViewingNum", "$FansNum", "$ParticipationNum", "$AddViewingNum", "$AddParticipationNum", "$AddFansNum"}}}},
		{"$sort": bson.M{"Sort": -1, "CreateTime": 1}},
		{"$project": bson.M{"Sort": 0}},
		{"$skip": curPage * pageSize}, {"$limit": pageSize},
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	query1 := col.Pipe(m)
	err = query1.All(&list)
	easygo.PanicError(err)
	return list, count
}

// 随机抽取根据指定条数的热门数据
func GetHotTopicFromDB(num int) []*share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	list := make([]*share_message.Topic, 0)

	queryBson := bson.M{"Status": TOPIC_STATUS_OPEN, "IsHot": true}

	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

// 随机抽取根据指定条数的热门数据,去除某些
func GetRangeHotTopicFromDB(num int, ids []int64) []*share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	list := make([]*share_message.Topic, 0)

	queryBson := bson.M{"Status": TOPIC_STATUS_OPEN, "IsHot": true, "IsRecommend": bson.M{"$ne": true}}
	if len(ids) > 0 {
		queryBson = bson.M{"Status": TOPIC_STATUS_OPEN, "IsHot": true, "IsRecommend": bson.M{"$ne": true}, "_id": bson.M{"$nin": ids}}
	}
	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

//  随机选取n条普通话题,去掉某些
func GetRangeCommonTopicFromDB(num int, topicIds []int64) []*share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	list := make([]*share_message.Topic, 0)

	queryBson := bson.M{}
	queryBson = bson.M{"Status": TOPIC_STATUS_OPEN, "TopicClass": TOPIC_CLASS_BS, "IsRecommend": bson.M{"$ne": true}, "IsHot": bson.M{"$ne": true}}
	if len(topicIds) > 0 {
		queryBson = bson.M{"Status": TOPIC_STATUS_OPEN, "TopicClass": TOPIC_CLASS_BS, "IsRecommend": bson.M{"$ne": true}, "IsHot": bson.M{"$ne": true}, "_id": bson.M{"$nin": topicIds}}
	}

	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

//  随机选取n条推荐的话题
func GetRangeIsRecommendTopicFromDB(num int, topicIds []int64) []*share_message.Topic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	list := make([]*share_message.Topic, 0)

	queryBson := bson.M{}
	queryBson = bson.M{"IsRecommend": true, "Status": TOPIC_STATUS_OPEN}
	if len(topicIds) > 0 {
		queryBson = bson.M{"IsRecommend": true, "Status": TOPIC_STATUS_OPEN, "_id": bson.M{"$nin": topicIds}}
	}

	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

// 获取推荐的话题.
func GetRecommendTopicByPageFromDB(page, pageSize int) ([]*share_message.Topic, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	list := make([]*share_message.Topic, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)
	queryBson := bson.M{}
	queryBson["IsRecommend"] = true
	queryBson["Status"] = TOPIC_STATUS_OPEN
	m := []bson.M{
		{"$match": queryBson},
		{"$addFields": bson.M{"Sort": bson.M{"$add": []string{"$ViewingNum", "$FansNum", "$ParticipationNum", "$AddViewingNum", "$AddParticipationNum", "$AddFansNum"}}}},
		//{"$sort": bson.M{"Sort": -1, "CreateTime": 1}},
		{"$sort": bson.M{"Sort": -1}},
		{"$project": bson.M{"Sort": 0}},
		{"$skip": curPage * pageSize}, {"$limit": pageSize},
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	query1 := col.Pipe(m)
	err = query1.All(&list)
	easygo.PanicError(err)
	return list, count

}

/*// 获取推荐的话题.
func GetRecommendTopicByPageFromDB(page, pageSize int) ([]*share_message.Topic, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	curPage := easygo.If(page > 1, page-1, 0).(int)
	list := make([]*share_message.Topic, 0)

	queryBson := bson.M{}
	queryBson["IsRecommend"] = true
	queryBson["Status"] = TOPIC_STATUS_OPEN
	query := col.Find(queryBson)

	count, err := query.Count()
	easygo.PanicError(err)

	err1 := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(err1)

	return list, count
}
*/
// 根据用户id获取自己关注的话题
func GetPlayerAttentionTopicsByPidFromDB(pid int64) []*share_message.PlayerAttentionTopic {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_PLAYER_ATTENTION)
	defer closeFun()
	tps := make([]*share_message.PlayerAttentionTopic, 0)
	err := col.Find(bson.M{"PlayerId": pid}).All(&tps)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return tps
}

// 分页获取热门动态.
func GetHotTopicByPageFromDB(page, pageSize int) ([]*share_message.Topic, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	curPage := easygo.If(page > 1, page-1, 0).(int)
	list := make([]*share_message.Topic, 0)

	queryBson := bson.M{}
	queryBson["IsHot"] = true
	queryBson["Status"] = TOPIC_STATUS_OPEN
	queryBson["TopicClass"] = TOPIC_CLASS_BS

	m := []bson.M{
		{"$match": queryBson},
		{"$addFields": bson.M{"Sort": bson.M{"$add": []string{"$ViewingNum", "$FansNum", "$ParticipationNum", "$AddViewingNum", "$AddParticipationNum", "$AddFansNum"}}}},
		{"$sort": bson.M{"Sort": -1, "CreateTime": 1}},
		{"$project": bson.M{"Sort": 0}},
		{"$skip": curPage * pageSize}, {"$limit": pageSize},
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	query1 := col.Pipe(m)
	err = query1.All(&list)
	easygo.PanicError(err)
	return list, count

}

// 随机抽取置顶个数的运营号 todo 暂时先放在这里,和话题统一放
func GetRandPlayer(num int) []*share_message.PlayerBase {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	list := make([]*share_message.PlayerBase, 0)

	queryBson := bson.M{}
	queryBson = bson.M{"Status": 0, "Types": ACCOUNT_TYPES_YXYY, "HeadIcon": bson.M{"$ne": ""}} //有效的运营号

	m := []bson.M{
		{"$match": bson.M{"FansList.0": bson.M{"$exists": 1}}},
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)
	return list
}

// 设置指定日期的用户话题贡献度(日榜)
func SetSpecifyTopicPlayerDevoteDay(playerId, topicId, devote int64, topicName string, year, month, day int32) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_PLAYER_DEVOTE_DAY)
	defer closeFun()
	err := col.Update(
		bson.M{"TopicId": topicId, "PlayerId": playerId, "Year": year, "Month": month, "Day": day},
		bson.M{"$inc": bson.M{"Devote": devote}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	if err != nil && err.Error() == mgo.ErrNotFound.Error() && devote > 0 {
		tps := &share_message.TopicPlayerDevoteDay{
			Id:         easygo.NewInt64(NextId(TABLE_TOPIC_PLAYER_DEVOTE_DAY)),
			TopicId:    easygo.NewInt64(topicId),
			PlayerId:   easygo.NewInt64(playerId),
			Year:       easygo.NewInt32(year),
			Month:      easygo.NewInt32(month),
			Day:        easygo.NewInt32(day),
			Devote:     easygo.NewInt64(devote),
			TopicName:  easygo.NewString(topicName),
			CreateTime: easygo.NewInt64(time.Now().Unix()),
		}
		col.Insert(tps)
	}
	return err
}

// 设置指定日期的用户话题贡献度(月榜)
func SetSpecifyTopicPlayerDevoteMonth(playerId, topicId, devote int64, topicName string, year, month int32) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_PLAYER_DEVOTE_MONTH)
	defer closeFun()
	err := col.Update(
		bson.M{"TopicId": topicId, "PlayerId": playerId, "Year": year, "Month": month},
		bson.M{"$inc": bson.M{"Devote": devote}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	if err != nil && err.Error() == mgo.ErrNotFound.Error() && devote > 0 {
		tps := &share_message.TopicPlayerDevoteDay{
			Id:         easygo.NewInt64(NextId(TABLE_TOPIC_PLAYER_DEVOTE_MONTH)),
			TopicId:    easygo.NewInt64(topicId),
			PlayerId:   easygo.NewInt64(playerId),
			Year:       easygo.NewInt32(year),
			Month:      easygo.NewInt32(month),
			Devote:     easygo.NewInt64(devote),
			TopicName:  easygo.NewString(topicName),
			CreateTime: easygo.NewInt64(time.Now().Unix()),
		}
		col.Insert(tps)
	}
	return err
}

// 设置用户话题贡献度(日榜)
func SetTopicPlayerDevoteDay(playerId, topicId, devote int64, topicName string) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_PLAYER_DEVOTE_DAY)
	defer closeFun()
	year, month, day := GetDateYMD()
	err := col.Update(
		bson.M{"TopicId": topicId, "PlayerId": playerId, "Year": year, "Month": month, "Day": day},
		bson.M{"$inc": bson.M{"Devote": devote}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	if err != nil && err.Error() == mgo.ErrNotFound.Error() && devote > 0 {
		tps := &share_message.TopicPlayerDevoteDay{
			Id:         easygo.NewInt64(NextId(TABLE_TOPIC_PLAYER_DEVOTE_DAY)),
			TopicId:    easygo.NewInt64(topicId),
			PlayerId:   easygo.NewInt64(playerId),
			Year:       easygo.NewInt32(year),
			Month:      easygo.NewInt32(month),
			Day:        easygo.NewInt32(day),
			Devote:     easygo.NewInt64(devote),
			TopicName:  easygo.NewString(topicName),
			CreateTime: easygo.NewInt64(time.Now().Unix()),
		}
		col.Insert(tps)
	}
	return err
}

// 设置用户话题贡献度(月榜)
func SetTopicPlayerDevoteMonth(playerId, topicId, devote int64, topicName string) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_PLAYER_DEVOTE_MONTH)
	defer closeFun()
	year, month, _ := GetDateYMD()
	err := col.Update(
		bson.M{"TopicId": topicId, "PlayerId": playerId, "Year": year, "Month": month},
		bson.M{"$inc": bson.M{"Devote": devote}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	if err != nil && err.Error() == mgo.ErrNotFound.Error() && devote > 0 {
		tps := &share_message.TopicPlayerDevoteDay{
			Id:         easygo.NewInt64(NextId(TABLE_TOPIC_PLAYER_DEVOTE_MONTH)),
			TopicId:    easygo.NewInt64(topicId),
			PlayerId:   easygo.NewInt64(playerId),
			Year:       easygo.NewInt32(year),
			Month:      easygo.NewInt32(month),
			Devote:     easygo.NewInt64(devote),
			TopicName:  easygo.NewString(topicName),
			CreateTime: easygo.NewInt64(time.Now().Unix()),
		}
		col.Insert(tps)
	}
	return err
}

// 设置用户话题贡献度(总榜)
func SetTopicPlayerDevoteTotal(playerId, topicId, devote int64, topicName string) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_PLAYER_DEVOTE_TOTAL)
	defer closeFun()
	err := col.Update(
		bson.M{"TopicId": topicId, "PlayerId": playerId},
		bson.M{"$inc": bson.M{"Devote": devote}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	if err != nil && err.Error() == mgo.ErrNotFound.Error() && devote > 0 {
		tps := &share_message.TopicPlayerDevoteDay{
			Id:         easygo.NewInt64(NextId(TABLE_TOPIC_PLAYER_DEVOTE_TOTAL)),
			TopicId:    easygo.NewInt64(topicId),
			PlayerId:   easygo.NewInt64(playerId),
			Devote:     easygo.NewInt64(devote),
			TopicName:  easygo.NewString(topicName),
			CreateTime: easygo.NewInt64(time.Now().Unix()),
		}
		col.Insert(tps)
	}
	return err
}

//获取每日贡献度排行榜
func GetTopicPlayerDevoteDayList(topicName string, topicId int64, page, pageSize int) ([]*share_message.TopicDevote, int) {
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_PLAYER_DEVOTE_DAY)
	defer closeFun()
	year, month, day := GetDateYMD()
	list := make([]*share_message.TopicPlayerDevoteDay, 0)
	where := bson.M{}
	if topicName != "" {
		where = bson.M{"TopicName": topicName, "Year": year, "Month": month, "Day": day, "Devote": bson.M{"$gt": DevoteMin}}
	} else {
		where = bson.M{"TopicId": topicId, "Year": year, "Month": month, "Day": day, "Devote": bson.M{"$gt": DevoteMin}}
	}
	err := col.Find(where).Skip(curPage * pageSize).Limit(pageSize).Sort("-Devote").All(&list)
	count, _ := col.Find(where).Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	var ids []int64
	for _, v := range list {
		ids = append(ids, v.GetPlayerId())
	}
	topicDevoteList := make([]*share_message.TopicDevote, 0)
	playerInfoList := GetPlayerListByIds(ids)
	for _, v := range list {
		for _, v1 := range playerInfoList {
			if v.GetPlayerId() == v1.GetPlayerId() {
				topicDevoteData := &share_message.TopicDevote{}
				topicDevoteData.PlayerId = easygo.NewInt64(v1.GetPlayerId())
				topicDevoteData.NickName = easygo.NewString(v1.GetNickName())
				topicDevoteData.HeadIcon = easygo.NewString(v1.GetHeadIcon())
				topicDevoteData.Sex = easygo.NewInt64(v1.GetSex())
				topicDevoteData.Devote = easygo.NewInt64(v.GetDevote())
				topicDevoteList = append(topicDevoteList, topicDevoteData)
				break
			}
		}
	}
	return topicDevoteList, count
}

//获取每月贡献度排行榜
func GetTopicPlayerDevoteMonthList(topicName string, topicId int64, page, pageSize int) ([]*share_message.TopicDevote, int) {
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_PLAYER_DEVOTE_MONTH)
	defer closeFun()
	year, month, _ := GetDateYMD()
	list := make([]*share_message.TopicPlayerDevoteDay, 0)
	where := bson.M{}
	if topicName != "" {
		where = bson.M{"TopicName": topicName, "Year": year, "Month": month, "Devote": bson.M{"$gt": DevoteMin}}
	} else {
		where = bson.M{"TopicId": topicId, "Year": year, "Month": month, "Devote": bson.M{"$gt": DevoteMin}}
	}

	err := col.Find(where).Skip(curPage * pageSize).Limit(pageSize).Sort("-Devote").All(&list)
	count, _ := col.Find(where).Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	var ids []int64
	for _, v := range list {
		ids = append(ids, v.GetPlayerId())
	}
	topicDevoteList := make([]*share_message.TopicDevote, 0)
	playerInfoList := GetPlayerListByIds(ids)
	for _, v := range list {
		for _, v1 := range playerInfoList {
			if v.GetPlayerId() == v1.GetPlayerId() {
				topicDevoteData := &share_message.TopicDevote{}
				topicDevoteData.PlayerId = easygo.NewInt64(v1.GetPlayerId())
				topicDevoteData.NickName = easygo.NewString(v1.GetNickName())
				topicDevoteData.HeadIcon = easygo.NewString(v1.GetHeadIcon())
				topicDevoteData.Sex = easygo.NewInt64(v1.GetSex())
				topicDevoteData.Devote = easygo.NewInt64(v.GetDevote())
				topicDevoteList = append(topicDevoteList, topicDevoteData)
				break
			}
		}
	}
	return topicDevoteList, count
}

//获取总贡献度排行榜
func GetTopicPlayerDevoteTotalList(topicName string, topicId int64, page, pageSize int) ([]*share_message.TopicDevote, int) {
	if page == 0 {
		page = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_PLAYER_DEVOTE_TOTAL)
	defer closeFun()
	list := make([]*share_message.TopicPlayerDevoteDay, 0)
	where := bson.M{}
	if topicName != "" {
		where = bson.M{"TopicName": topicName, "Devote": bson.M{"$gt": DevoteMin}}
	} else {
		where = bson.M{"TopicId": topicId, "Devote": bson.M{"$gt": DevoteMin}}
	}

	err := col.Find(where).Skip(curPage * pageSize).Limit(pageSize).Sort("-Devote").All(&list)
	count, _ := col.Find(where).Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	var ids []int64
	for _, v := range list {
		ids = append(ids, v.GetPlayerId())
	}
	topicDevoteList := make([]*share_message.TopicDevote, 0)
	playerInfoList := GetPlayerListByIds(ids)
	for _, v := range list {
		for _, v1 := range playerInfoList {
			if v.GetPlayerId() == v1.GetPlayerId() {
				topicDevoteData := &share_message.TopicDevote{}
				topicDevoteData.PlayerId = easygo.NewInt64(v1.GetPlayerId())
				topicDevoteData.NickName = easygo.NewString(v1.GetNickName())
				topicDevoteData.HeadIcon = easygo.NewString(v1.GetHeadIcon())
				topicDevoteData.Sex = easygo.NewInt64(v1.GetSex())
				topicDevoteData.Devote = easygo.NewInt64(v.GetDevote())
				topicDevoteList = append(topicDevoteList, topicDevoteData)
				break
			}
		}
	}
	return topicDevoteList, count
}

//添加申请话题主记录
func AddApplyTopicMaster(data *share_message.ApplyTopicMaster) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_APPLY_TOPIC_MASTER)
	defer closeFun()
	data.Id = easygo.NewInt64(NextId(TABLE_TOPIC_APPLY_TOPIC_MASTER))
	data.CreateTime = easygo.NewInt64(util.GetTime())
	data.Status = easygo.NewInt32(0)
	err := col.Insert(data)
	return err
}

//申请修改话题信息
func EditTopicInfo(data *share_message.ApplyEditTopicInfo) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_APPLY_LOG)
	defer closeFun()
	data.CreateTime = easygo.NewInt64(util.GetTime())
	data.Status = easygo.NewInt32(0)
	err := col.Insert(data)
	return err
}

//退出话题主
func QuitTopicMaster(topicId, playerId int64) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC)
	defer closeFun()
	where := bson.M{"_id": topicId, "TopicMaster": playerId}
	err := col.Update(where, bson.M{"$set": bson.M{"IsOpen": true, "TopicMaster": int64(0), "Owner": ""}})
	return err
}
