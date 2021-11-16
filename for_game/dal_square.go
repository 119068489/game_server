package for_game

import (
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

// 新增或修改社交广场动态置顶定时任务管理器
func UpsetSquareTopTimerMgrToDB(mgr *share_message.BackstageNotifyTopReq) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_TOP_TIMER_MGR)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": mgr.GetLogId()}, bson.M{"$set": mgr})
	easygo.PanicError(err)
}

// 获取所有社交广场动态置顶定时任务管理器
func GetAllSquareTopTimeMgrFromDB() []*share_message.BackstageNotifyTopReq {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_TOP_TIMER_MGR)
	defer closeFun()

	var lst []*share_message.BackstageNotifyTopReq
	err := col.Find(bson.M{}).All(&lst)
	easygo.PanicError(err)
	return lst
}

func DelSquareTopTimeMgrByIdFromDB(id int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_TOP_TIMER_MGR)
	defer closeFun()
	err := col.Remove(bson.M{"_id": id})
	easygo.PanicError(err)
}

const (
	DataLimit = 100
)

//==============动态相关数据============================
/**
xiong
倒序获取动态列表,传入限制条数
*/
func GetDynamicListByLimitFromDB(limit int) []*share_message.DynamicData {
	var lst []*share_message.DynamicData
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFunc()
	err := col.Find(bson.M{"Statue": DYNAMIC_STATUE_COMMON}).Sort("-_id").Limit(limit).All(&lst) //起服加载最新100条
	easygo.PanicError(err)
	return lst
}

/**
xiong
根据玩家id和状态查找动态列表
*/
func GetDynamicListByPlayerIdAndStatusFromDB(playerId int64, status int) ([]*share_message.DynamicData, error) {
	var lst []*share_message.DynamicData
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFunc()
	err := col.Find(bson.M{"PlayerId": playerId, "Statue": status}).All(&lst)
	return lst, err
}

/**
xiong
从mongo中加载小于logid的100条动态
*/
func GetDynamicListByLtLogIdAndLimitFromDB(logId int64, status, limit int) ([]*share_message.DynamicData, error) {
	var lst []*share_message.DynamicData
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFunc()
	err := col.Find(bson.M{"_id": bson.M{"$lt": logId}, "Statue": status}).Sort("-_id").Limit(limit).All(&lst) //起服加载最新100条
	return lst, err
}

/**
xiong
根据动态id列表获取动态列表.
*/
func GetDynamicListByLogIdsAndStatusFromDB(logIds []int64, status int) ([]*share_message.DynamicData, error) {
	var lst []*share_message.DynamicData
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFunc()
	err := col.Find(bson.M{"_id": bson.M{"$in": logIds}, "Statue": status}).All(&lst)
	return lst, err
}

/**
xiong
获取单个动态信息
*/
func GetDynamicByLogIdFromDB(logId int64) (*share_message.DynamicData, error) {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFunc()
	var dynamicData *share_message.DynamicData
	err := col.Find(bson.M{"_id": logId}).One(&dynamicData)
	return dynamicData, err
}

//获取用户指定话题的动态数
func GetPlayerTopicDynamicCount(playerId, topicId int64) int {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFunc()
	topicIdArr := make([]int64, 0)
	topicIdArr = append(topicIdArr, topicId)
	whereBson := bson.M{"PlayerId": playerId, "TopicId": bson.M{"$elemMatch": bson.M{"$in": topicIdArr}}}
	count, _ := col.Find(whereBson).Count()
	return count
}

/**
xiong
插入单个动态insert
*/
func InsertDynamicToDB(msg *share_message.DynamicData) error {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFunc()
	_, err := col.Upsert(bson.M{"_id": msg.GetLogId()}, msg)
	return err
	// return col.Insert(msg)
}

// 更新senderType字段
func UpdateDynamicSenderTypeToDB(logId int64, senderType int32) {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFunc()
	err := col.Update(bson.M{"_id": logId}, bson.M{"$set": bson.M{"SenderType": senderType}})
	if err != nil || err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
}

/**
xiong
更新动态状态
*/
func UpSertDynamicStatusByLogIdFromDB(logId int64, status int32, note string) error {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFunc()
	_, err := col.Upsert(bson.M{"_id": logId}, bson.M{"$set": bson.M{"Statue": status, "Note": note}})
	return err
}

//保存社交动态
func SaveSquareDynamic(id int64) {
	b := GetRedisDynamic(id)
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFunc()

	_, err := col.Upsert(bson.M{"_id": id}, b)
	easygo.PanicError(err)
}

// 分页查询出不包含置顶的动态列表
func GetNoTopDynamicListByPageFromDB(opId int64, page, pageSize int, dynamicIds ...int64) ([]*share_message.DynamicData, int) {
	var list []*share_message.DynamicData
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()

	queryBson := bson.M{}
	if player := GetRedisPlayerBase(opId); player != nil {
		//青少年保护模式不显示视频
		if player.GetYoungPassWord() != "" {
			//queryBson["Video"] = nil
			queryBson = bson.M{"$or": []bson.M{{"Video": ""}, {"Video": nil}}}

		}
	}
	did := append(dynamicIds, 0)[0]

	//queryBson := bson.M{"IsTop": bson.M{"$ne": true}, "IsBsTop": bson.M{"$ne": true}, "Statue": 0}
	queryBson["IsTop"] = bson.M{"$ne": true}
	queryBson["IsBsTop"] = bson.M{"$ne": true}
	queryBson["Statue"] = 0
	if did > 0 {
		queryBson["_id"] = bson.M{"$lte": did}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	err1 := query.Sort("-SendTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(err1)
	return list, count
}

// 查询出一批人的动态列表
func GetNoTopDynamicByPIDsFromDB(opId int64, page, pageSize int, ids []int64, dynamicIds ...int64) ([]*share_message.DynamicData, int) {
	var list []*share_message.DynamicData
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()

	queryBson := bson.M{}
	if player := GetRedisPlayerBase(opId); player != nil {
		//青少年保护模式不显示视频
		if player.GetYoungPassWord() != "" {
			//queryBson["Video"] = nil
			queryBson = bson.M{"$or": []bson.M{{"Video": ""}, {"Video": nil}}}

		}
	}
	did := append(dynamicIds, 0)[0]
	//queryBson := bson.M{"PlayerId": bson.M{"$in": ids}, "IsTop": bson.M{"$ne": true}, "IsBsTop": bson.M{"$ne": true}, "Statue": 0}
	queryBson["PlayerId"] = bson.M{"$in": ids}
	queryBson["IsTop"] = bson.M{"$ne": true}
	queryBson["IsBsTop"] = bson.M{"$ne": true}
	queryBson["Statue"] = 0
	if did > 0 {
		queryBson["_id"] = bson.M{"$lte": did}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	err1 := query.Sort("-SendTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(err1)
	return list, count
}

// 分页去重 查询出一批人的动态列表
func GetNoTopDynamicByPIDs(opId int64, page, pageSize int, ids []int64, key string) ([]*share_message.DynamicData, int) {
	var dynamicId int64
	if page > 1 { //
		dynamicId = GetFirstPageMaxLogIdFromRedis(key)
	}
	list, count := GetNoTopDynamicByPIDsFromDB(opId, page, pageSize, ids, dynamicId)

	if page == 1 { // 找到最大的id
		var maxId int64
		for _, v := range list {
			if v.GetLogId() > maxId {
				maxId = v.GetLogId()
			}
		}
		// 设置id最大的值进redis
		SetFirstPageMaxLogIdToRedis(key, easygo.AnytoA(maxId))
	}
	return list, count
}

// 获取置顶人的后台置顶动态
func GetBSTopDynamicListByIDsFromDB(opId int64, ids []int64) []*share_message.DynamicData {
	list := make([]*share_message.DynamicData, 0)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	queryBson := bson.M{}
	if player := GetRedisPlayerBase(opId); player != nil {
		//青少年保护模式不显示视频
		if player.GetYoungPassWord() != "" {
			//queryBson["Video"] = nil
			queryBson = bson.M{"$or": []bson.M{{"Video": ""}, {"Video": nil}}}

		}
	}
	queryBson = bson.M{"IsBsTop": true, "Statue": 0}
	if len(ids) > 0 {
		queryBson = bson.M{"PlayerId": bson.M{"$in": ids}, "IsBsTop": true, "Statue": 0} // 包含自己关注的人的
	}

	if err := col.Find(queryBson).All(&list); err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	return list
}

// 获取置顶人的APP置顶动态
func GetAppTopDynamicListByIDsFromDB(opId int64, ids []int64) []*share_message.DynamicData {
	list := make([]*share_message.DynamicData, 0)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	queryBson := bson.M{}
	if player := GetRedisPlayerBase(opId); player != nil {
		//青少年保护模式不显示视频
		if player.GetYoungPassWord() != "" {
			//queryBson["Video"] = nil
			queryBson = bson.M{"$or": []bson.M{{"Video": ""}, {"Video": nil}}}

		}
	}
	queryBson = bson.M{"IsTop": true, "Statue": 0}
	if len(ids) > 0 {
		queryBson = bson.M{"PlayerId": bson.M{"$in": ids}, "IsTop": true, "Statue": 0}
	}

	if err := col.Find(queryBson).All(&list); err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	return list
}

// 添加动态的热门分数
func UpsetDynamicScore(logId int64, score int32) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": logId}, bson.M{"$inc": bson.M{"HostScore": score}})
	easygo.PanicError(err)
}

func GetDynamicByStatusSFromDB(id int64, status []int) *share_message.DynamicData {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	data := new(share_message.DynamicData)
	queryBson := bson.M{"_id": id, "Statue": bson.M{"$in": status}}
	if err := col.Find(queryBson).One(&data); err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return data
}

// 取动态,带有图片别切最热的3个动态
func GetHasPhotoHotDynamicToDB(pid int64, num int) []*share_message.DynamicData {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	ds := make([]*share_message.DynamicData, 0)
	queryBson := bson.M{"PlayerId": pid, "Statue": DYNAMIC_STATUE_COMMON, "Photo": bson.M{"$elemMatch": bson.M{"$ne": nil}}}
	query := col.Find(queryBson)

	err1 := query.Sort("-SendTime").Limit(num).All(&ds)
	easygo.PanicError(err1)
	return ds
}

// 举报次数加1
func UpdateComplaintNum(logId int64, num int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	err := col.Update(bson.M{"_id": logId}, bson.M{"$inc": bson.M{"ReportCount": num}})
	easygo.PanicError(err)
}

//==============动态相关数据============================

//==============评论相关数据============================
/**
xiong
根据动态id列表和状态获取评论列表
*/
func GetCommentListByLogIdsAndStatusFromDB(logIds []int64, status int) ([]*share_message.CommentData, error) {
	var commentInfo []*share_message.CommentData
	col1, closeFunc1 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_COMMENT)
	defer closeFunc1()
	err := col1.Find(bson.M{"LogId": bson.M{"$in": logIds}, "Statue": status}).All(&commentInfo)
	return commentInfo, err
}

func GetCommentListByPlayerIdFromDB(playerId int64) ([]*share_message.CommentData, error) {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_COMMENT)
	defer closeFunc()
	var lst []*share_message.CommentData
	err := col.Find(bson.M{"PlayerId": playerId}).All(&lst)
	return lst, err
}

func GetAttentionListByPlayerIdFromDB(playerId int64) ([]*share_message.AttentionData, error) {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_ATTENTION)
	defer closeFunc()
	var lst []*share_message.AttentionData
	err := col.Find(bson.M{"PlayerId": playerId}).All(&lst)
	return lst, err
}

//增加评论
func InsertCommentToDB(msg *share_message.CommentData) error {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_COMMENT)
	defer closeFunc()
	return col.Insert(msg)
}

// 查找单条评论
func GetCommentDataFromDB(cid int64) *share_message.CommentData {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_COMMENT)
	defer closeFunc()
	var c *share_message.CommentData
	err := col.Find(bson.M{"_id": cid, "Statue": DYNAMIC_COMMENT_STATUE_COMMON}).One(&c)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return c
}

//修改评论赞数,isComment:是否是评论回复加分
func IncrCommentZanNumToDB(commentId int64, val int64, isComment ...bool) *share_message.CommentData {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_COMMENT)
	defer closeFunc()
	isCom := append(isComment, false)[0]
	update := bson.M{"$inc": bson.M{"Score": val}}
	if !isCom {
		update = bson.M{"$inc": bson.M{"ZanNum": val, "Score": val}}
	}
	result := &share_message.CommentData{}
	_, err := col.Find(bson.M{"_id": commentId}).Apply(mgo.Change{
		Update:    update,
		Upsert:    true,
		ReturnNew: true,
	}, &result)
	easygo.PanicError(err)

	return result
}

//取消对评论的点赞，删除点赞记录
func DelDynamicCommentZan(Id int64) {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_COMMENT_ZAN)
	defer closeFunc()
	err := col.RemoveId(Id)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
}

//删除评论相关的赞
func DelDynamicCommentZanByCommentId(commentId int64) {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_COMMENT_ZAN)
	defer closeFunc()
	err := col.Remove(bson.M{"CommentId": commentId})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
}

//增加对评论的点赞
func AddDynamicCommentZan(zan *share_message.CommentDataZan) {
	id := NextId(TABLE_SQUARE_COMMENT_ZAN)
	zan.Id = easygo.NewInt64(id)
	zan.CreateTime = easygo.NewInt64(GetMillSecond())
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_COMMENT_ZAN)
	defer closeFunc()
	err := col.Insert(zan)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
}

//通过页码获取主评论
func GetDynamicCommentInfoByPage(id int64, page, pageSize int, hotList []int64, params *share_message.SysParameter) *share_message.CommentList {
	list := make([]*share_message.CommentData, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_COMMENT)
	defer closeFun()
	queryBson := bson.M{}
	queryBson["LogId"] = id
	queryBson["BelongId"] = 0
	queryBson["Statue"] = DYNAMIC_COMMENT_STATUE_COMMON
	var hotComment []*share_message.CommentData
	if page == 1 {
		//先查热门
		queryBson["Score"] = bson.M{"$gte": params.GetCommentHotScore()}
		query := col.Find(queryBson)
		err := query.Sort("-Score", "-CreateTime").Limit(int(params.GetCommentHotCount())).All(&hotComment)
		easygo.PanicError(err)
		hotList = make([]int64, 0)
		for _, com := range hotComment {
			hotList = append(hotList, com.GetId())
		}
		delete(queryBson, "Score")
		pageSize = pageSize - len(hotList)
	}
	//再查剩下的
	queryBson["_id"] = bson.M{"$nin": hotList}
	query := col.Find(queryBson)
	err := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	if err != nil && err.Error() == mgo.ErrNotFound.Error() {
		panic(err)
	}
	if page == 1 {
		list = append(hotComment, list...)
	}
	msg := &share_message.CommentList{
		CommentInfo: list,
		HotList:     hotList,
	}
	return msg
}

//获取评论获得赞数
func GetDynamicCommentZanNum(commentId int64) int64 {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_COMMENT_ZAN)
	defer closeFunc()
	n, err := col.Find(bson.M{"CommentId": commentId}).Count()
	easygo.PanicError(err)
	return int64(n)

}

//==============评论相关数据============================

//==============赞相关数据============================
/**
xiong
根据operateId查找赞列表数据
*/
func GetZanDataListByOperateIdFromDB(operateId int64) ([]*share_message.ZanData, error) {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_ZAN)
	defer closeFunc()
	var lst []*share_message.ZanData
	err := col.Find(bson.M{"OperateId": operateId}).All(&lst)
	return lst, err
}

/**
xiong
根据动态id列表查找赞列表信息
*/
func GetZanDataListByDynamicIdsFromDB(logIds []int64) ([]*share_message.ZanData, error) {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_ZAN)
	defer closeFunc()
	var lst []*share_message.ZanData
	err := col.Find(bson.M{"DynamicId": bson.M{"$in": logIds}}).All(&lst)
	return lst, err
}

// 获取所有的假赞数
func GetAllTrueZan(pid int64) int64 {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFunc()
	queryBson := bson.M{"PlayerId": pid}
	//聚合查询
	m := []bson.M{
		{"$match": queryBson},
		{"$group": bson.M{"_id": "$PlayerId", "TrueZan": bson.M{"$sum": "$TrueZan"}}},
	}
	query := col.Pipe(m)
	res := []bson.M{}
	err := query.All(&res)
	easygo.PanicError(err)
	if len(res) == 0 {
		return 0
	}
	return easygo.AtoInt64(easygo.AnytoA(res[0]["TrueZan"]))
}

func UpsetAllZanDataToDB(zs []*share_message.ZanData) {
	var data []interface{}
	for _, v := range zs {
		b1 := bson.M{"_id": v.GetLogId()}
		b2 := v
		data = append(data, b1, b2)
	}
	UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_SQUARE_ZAN, data)
}

func DelZanDataToDB(zanId int64) error {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_ZAN)
	defer closeFunc()
	queryBson := bson.M{"_id": zanId}
	err := col.Remove(queryBson)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
		return err
	}
	return nil
}

func GetZanDataByIdAndDynamicId(id, dynamicId int64) *share_message.ZanData {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_ZAN)
	defer closeFunc()
	queryBson := bson.M{"_id": id, "DynamicId": dynamicId}
	var zanData *share_message.ZanData

	err := col.Find(queryBson).One(&zanData)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return zanData

}
func GetZanDataListByDynamicId(dynamicId int64) []*share_message.ZanData {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_ZAN)
	defer closeFunc()
	queryBson := bson.M{"DynamicId": dynamicId}
	var zanDataList []*share_message.ZanData

	err := col.Find(queryBson).All(&zanDataList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return zanDataList

}
func GetZanDataListByPid(pid int64) []*share_message.ZanData {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_ZAN)
	defer closeFunc()
	queryBson := bson.M{"PlayerId": pid}
	var zanDataList []*share_message.ZanData

	err := col.Find(queryBson).All(&zanDataList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return zanDataList

}

//==============热门相关数据============================
//查询热门分大于指定分值得动态
func GetHotDynamic() []*share_message.DynamicData {
	list := make([]*share_message.DynamicData, 0)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()

	queryBson := bson.M{"HostScore": bson.M{"$gt": 0}}
	if err := col.Find(queryBson).All(&list); err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	return list
}

// 根据话题id查找动态列表,根据热门分倒序
func GetHotDynamicByTopicIdFromDB(topicID int64, page, pageSize int) ([]*share_message.DynamicData, int) {
	list := make([]*share_message.DynamicData, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	queryBson := bson.M{}
	queryBson["TopicId"] = topicID
	queryBson["Statue"] = DYNAMIC_STATUE_COMMON

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	err1 := query.Sort("-HostScore").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(err1)
	return list, count
}

// 根据话题id查找动态列表,根据热门分倒序,根据用户去重
func GetDeviceHotDynamicByTopicIdFromDB(topicID int64, page, pageSize int) ([]*share_message.DynamicData, int) {
	list := make([]*share_message.DynamicData, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()

	m := []bson.M{
		{"$match": bson.M{"TopicId": topicID}},
		{"$group": bson.M{"_id": "$PlayerId"}},
		{"$project": bson.M{"PlayerId": 1}},
		{"$sort": bson.M{"HostScore": 1}},
		{"$skip": curPage * pageSize}, {"$limit": pageSize},
	}
	query := col.Pipe(m)
	err1 := query.All(&list)

	// 获取条数
	m1 := []bson.M{
		{"$match": bson.M{"TopicId": topicID}},
		{"$group": bson.M{"_id": "$PlayerId"}},
		{"$group": bson.M{"_id": "$PlayerId", "total": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"HostScore": -1}},
	}

	query1 := col.Pipe(m1)
	rst := make([]bson.M, 0)
	e := query1.All(&rst)
	var sum int = 0
	if nil == e {
		if rst != nil && len(rst) > 0 {
			sum = (rst[0]["total"]).(int)
		}
	} else {
		sum = 0
	}
	easygo.PanicError(err1)
	return list, sum
}

// 根据话题id查询已经是热门的动态,倒序排序.
func GetIsHotDynamicByTopicIdFromDB(hotScore int32, page, pageSize int, topicID int64) ([]*share_message.DynamicData, int) {
	list := make([]*share_message.DynamicData, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()

	queryBson := bson.M{"HostScore": bson.M{"$gte": hotScore}}
	queryBson["TopicId"] = topicID
	queryBson["Statue"] = 0
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	startNum := curPage * pageSize
	endNum := startNum + pageSize
	topicHotDynamicCache := topicHotDynamicCacheGlobal
	if endNum > topicHotDynamicCache.GetDefCacheNumber() {
		err1 := query.Sort("-HostScore").Skip(curPage * pageSize).Limit(pageSize).All(&list)
		easygo.PanicError(err1)
		logs.Info("hot-list:", list)
	} else {
		topicHotDynamicCache.CreateLockChan(topicID) //必须要先创建锁
		topicDynamicIds := make([]int64, 0)
		listCache := make([]*share_message.DynamicData, 0)
		if count > 0 && topicHotDynamicCache.GetTopicDynamicCount(topicID) != count {
			topicHotDynamicCache.SetTopicHotDynamicCacheList(topicID, hotScore, count)
			topicDynamicIds = topicHotDynamicCache.GetTopicNewDynamicCacheIds(page, pageSize, topicID)
		} else {
			//如果置顶数发生变化，更新缓存
			topicDynamicTopCount := GetTopicNewDynamicTopCount(topicID)
			topicDynamicTopCountCacheNum := topicHotDynamicCache.GetTopicDynamicTopCount(topicID)
			if topicDynamicTopCount > 0 && topicDynamicTopCount != topicDynamicTopCountCacheNum || (topicDynamicTopCount == 0 && topicDynamicTopCountCacheNum != 0) {
				topicHotDynamicCache.SetTopicHotDynamicCacheList(topicID, hotScore, count)
			}
			topicDynamicIds = topicHotDynamicCache.GetTopicNewDynamicCacheIds(page, pageSize, topicID)
		}
		err1 := col.Find(bson.M{"_id": bson.M{"$in": topicDynamicIds}}).All(&listCache)
		//进行排序
		for _, v := range topicDynamicIds {
			for _, v1 := range listCache {
				if v1.GetLogId() == v {
					list = append(list, v1)
					break
				}
			}
		}
		easygo.PanicError(err1)
		logs.Info("hot-list---Cache:", list)
	}
	return list, count
}

//获取热门的动态列表
func GetTopicHotDynamicList(hotScore int32, startSize, pageSize int, topicID int64) []*share_message.DynamicData {
	list := make([]*share_message.DynamicData, 0)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	queryBson := bson.M{"HostScore": bson.M{"$gte": hotScore}}
	queryBson["TopicId"] = topicID
	queryBson["Statue"] = DYNAMIC_STATUE_COMMON

	query := col.Find(queryBson)
	err1 := query.Sort("-HostScore").Select(bson.M{"_id": 1}).Skip(startSize).Limit(pageSize).All(&list)
	easygo.PanicError(err1)
	return list
}

//获取置顶的动态数量
func GetTopicNewDynamicTopCount(topicID int64) int {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	queryBson := bson.M{}
	queryBson["TopicId"] = topicID
	queryBson["Statue"] = DYNAMIC_STATUE_COMMON
	queryBson["TopicTopSet.IsTopicTop"] = true
	queryBson["TopicTopSet.TopicId"] = topicID
	m := []bson.M{
		{"$project": bson.M{"_id": 1, "TopicId": 1, "Statue": 1, "SendTime": 1, "TopicTopSet": 1}},
		{"$unwind": "$TopicTopSet"},
		{"$match": queryBson},
		{"$count": "total"},
	}
	query := col.Pipe(m)
	var listTop []interface{}
	err := query.All(&listTop)
	easygo.PanicError(err)
	if len(listTop) < 1 {
		return 0
	}
	return listTop[0].(bson.M)["total"].(int)
}

//获取置顶的最新动态
func GetTopicNewDynamicTopList(topicID int64) []*share_message.DynamicData {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	queryBson := bson.M{}
	queryBson["TopicId"] = topicID
	queryBson["Statue"] = DYNAMIC_STATUE_COMMON
	queryBson["TopicTopSet.IsTopicTop"] = true
	queryBson["TopicTopSet.TopicId"] = topicID
	m := []bson.M{
		{"$project": bson.M{"_id": 1, "TopicId": 1, "Statue": 1, "SendTime": 1, "TopicTopSet": 1}},
		{"$unwind": "$TopicTopSet"},
		{"$match": queryBson},
		{"$sort": bson.M{"TopicTopSet.TopicTopTime": -1}},
	}
	query := col.Pipe(m)
	listTop := make([]*share_message.DynamicData, 0)
	err := query.All(&listTop)
	easygo.PanicError(err)
	return listTop
}

//获取最新的动态列表
func GetTopicNewDynamicList(startSize, pageSize int, topicID int64) []*share_message.DynamicData {
	list := make([]*share_message.DynamicData, 0)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	queryBson := bson.M{}
	queryBson["TopicId"] = topicID
	queryBson["Statue"] = DYNAMIC_STATUE_COMMON

	query := col.Find(queryBson)
	err1 := query.Sort("-SendTime").Select(bson.M{"_id": 1}).Skip(startSize).Limit(pageSize).All(&list)
	easygo.PanicError(err1)
	return list
}

// 根据话题id查询最新的动态,倒序排序.
func GetSortTimeDynamicByTopicIdFromDB(page, pageSize int, topicID int64) ([]*share_message.DynamicData, int) {
	list := make([]*share_message.DynamicData, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()
	queryBson := bson.M{}
	queryBson["TopicId"] = topicID
	queryBson["Statue"] = DYNAMIC_STATUE_COMMON
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	startNum := curPage * pageSize
	endNum := startNum + pageSize
	topicNewDynamicCache := topicNewDynamicCacheGlobal
	if endNum > topicNewDynamicCache.GetDefCacheNumber() {
		err1 := query.Sort("-SendTime").Skip(startNum).Limit(pageSize).All(&list)
		easygo.PanicError(err1)
		logs.Info("new-list:", list)
	} else {
		topicNewDynamicCache.CreateLockChan(topicID) //必须要先创建锁
		topicDynamicIds := make([]int64, 0)
		listCache := make([]*share_message.DynamicData, 0)
		if count > 0 && topicNewDynamicCache.GetTopicDynamicCount(topicID) != count {
			topicNewDynamicCache.SetTopicNewDynamicCacheList(topicID, count)
			topicDynamicIds = topicNewDynamicCache.GetTopicNewDynamicCacheIds(page, pageSize, topicID)
		} else {
			//如果置顶数发生变化，更新缓存
			topicDynamicTopCount := GetTopicNewDynamicTopCount(topicID)
			TopicDynamicTopCountCacheNum := topicNewDynamicCache.GetTopicDynamicTopCount(topicID)
			if topicDynamicTopCount > 0 && topicDynamicTopCount != TopicDynamicTopCountCacheNum || (topicDynamicTopCount == 0 && TopicDynamicTopCountCacheNum != 0) {
				topicNewDynamicCache.SetTopicNewDynamicCacheList(topicID, count)
			}
			topicDynamicIds = topicNewDynamicCache.GetTopicNewDynamicCacheIds(page, pageSize, topicID)
		}
		err1 := col.Find(bson.M{"_id": bson.M{"$in": topicDynamicIds}}).All(&listCache)
		//进行排序
		for _, v := range topicDynamicIds {
			for _, v1 := range listCache {
				if v1.GetLogId() == v {
					list = append(list, v1)
					break
				}
			}
		}
		easygo.PanicError(err1)
		//logs.Info("new-list---Cache:",list)
	}

	return list, count
}

/**
随机选取动态大于3条的用户
gtNum:动态大于的条数
size:用户个数.
*/
func GetRandPlayerByDynamicCountFromDB(gtNum, size int) []*share_message.DynamicData {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()

	list := make([]*share_message.DynamicData, 0)

	m := []bson.M{
		{"$match": bson.M{"Statue": DYNAMIC_STATUE_COMMON}},
		{"$group": bson.M{"_id": "$PlayerId", "count": bson.M{"$sum": 1}}},
		{"$match": bson.M{"count": bson.M{"$gt": gtNum}}},
		{"$sample": bson.M{"size": size}}, //随机查询
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)
	return list
}

// 根据话题id列表查询已经是热门的动态,倒序排序.
func GetSortTimeDynamicByTopicIdListFromDB(page, pageSize int, topicIDs []int64) ([]*share_message.DynamicData, int) {
	list := make([]*share_message.DynamicData, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()

	queryBson := bson.M{}
	queryBson["TopicId"] = bson.M{"$in": topicIDs}
	queryBson["Statue"] = DYNAMIC_STATUE_COMMON

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	err1 := query.Sort("-SendTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(err1)
	return list, count
}

// 根据话题id查找超过热门分的动态列表,根据用户去重
func GetDevicePlayerHotDynamicByTopicIdFromDB(topicID int64, hotScore int32, page, pageSize int) ([]*share_message.DynamicData, int) {
	list := make([]*share_message.DynamicData, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()

	m := []bson.M{
		{"$match": bson.M{"HostScore": bson.M{"$gte": hotScore}, "Statue": 0, "TopicId": topicID}},
		{"$group": bson.M{"_id": "$PlayerId", "HostScore": bson.M{"$sum": "$HostScore"}}},
		{"$sort": bson.M{"HostScore": -1}},
		{"$skip": curPage * pageSize}, {"$limit": pageSize},
	}

	query := col.Pipe(m)
	err1 := query.All(&list)

	m1 := []bson.M{
		{"$match": bson.M{"HostScore": bson.M{"$gte": hotScore}, "Statue": 0, "TopicId": topicID}},
		{"$group": bson.M{"_id": "$PlayerId"}},
		{"$group": bson.M{"_id": "$PlayerId", "total": bson.M{"$sum": 1}}},
	}

	query1 := col.Pipe(m1)
	rst := make([]bson.M, 0)
	e := query1.All(&rst)
	var sum int = 0
	if nil == e {
		if rst != nil && len(rst) > 0 {
			sum = (rst[0]["total"]).(int)
		}
	} else {
		sum = 0
	}
	easygo.PanicError(err1)
	return list, sum
}

func GetSecondaryCommentNum(belongId int64) int {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_COMMENT)
	defer closeFunc()
	query := bson.M{"BelongId": belongId, "Statue": DYNAMIC_COMMENT_STATUE_COMMON}
	n, err := col.Find(query).Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return n
}

// 根据用户id列表找到动态列表,分页,按时间倒序   [1887530735 (64条),1887567020(4条)]
func GetDynamicListByPids(page, pageSize int, pids []int64) ([]*share_message.DynamicData, int) {
	list := make([]*share_message.DynamicData, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()

	queryBson := bson.M{}
	queryBson["Statue"] = DYNAMIC_STATUE_COMMON
	queryBson["PlayerId"] = bson.M{"$in": pids}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	err1 := query.Sort("-SendTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(err1)
	return list, count
}

//从一波玩家中找出最新的动态
func GetNewOneDynamicByPids(pids []int64) *share_message.DynamicData {
	var data *share_message.DynamicData
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFun()

	queryBson := bson.M{}
	queryBson["Statue"] = DYNAMIC_STATUE_COMMON
	queryBson["PlayerId"] = bson.M{"$in": pids}

	err := col.Find(queryBson).Sort("-SendTime").One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetNewOneDynamicByPids")
	}
	return data
}

//获取话题下的动态信息
func GetTopicDynamicByLogId(topicId, logId int64) (*share_message.DynamicData, error) {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFunc()
	var dynamicData *share_message.DynamicData
	err := col.Find(bson.M{"_id": logId, "TopicId": topicId}).One(&dynamicData)
	return dynamicData, err
}

//更新动态绑定的话题
func UpTopicDynamicByLogId(topicId, logId int64, topicArr []int64) error {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SQUARE_DYNAMIC)
	defer closeFunc()
	err := col.Update(bson.M{"_id": logId, "TopicId": topicId}, bson.M{"$set": bson.M{"TopicId": topicArr}})
	return err
}

//话题主删除动态操作日志表
func AddTopicMasterDelDynamicLog(data *share_message.TopicMasterDelDynamicLog) error {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOPIC_MASTER_DEL_DYNAMIC_LOG)
	defer closeFunc()
	data.CreateTime = easygo.NewInt64(util.GetTime())
	return col.Insert(data)
}
