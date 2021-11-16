package for_game

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/share_message"
	"log"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
)

var _ = fmt.Sprintf
var _ = log.Println
var _ = easygo.Underline

const (
	WAITER_MESSAGE_NEW   = 0 //未接待
	WAITER_MESSAGE_ING   = 1 //正在接待
	WAITER_MESSAGE_END   = 2 //已结束
	WAITER_MESSAGE_GRADE = 3 //已评价
)

const (
	WAITER_MSG_TYPE_C = 1 //玩家发送的消息
	WAITER_MSG_TYPE_S = 2 //客服发送的消息
)

const (
	FAQ_TITLE_SEARCH = 1 //标题搜索
	FAQ_KEY_SEARCH   = 2 //关键词搜索
)

//玩家Id查询单条消息
func QueryIMmessageByPid(id int64, status ...int32) *share_message.IMmessage {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_WAITER_MESSAGE)
	defer closeFun()

	var msg *share_message.IMmessage
	queryBson := bson.M{"PlayerId": id}
	if len(status) > 0 {
		queryBson["Status"] = status[0]
	}
	query := col.Find(queryBson)
	err := query.One(&msg)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return msg
}

//消息Id查询单条消息
func QueryIMmessageByMid(id int64, status ...int32) *share_message.IMmessage {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_WAITER_MESSAGE)
	defer closeFun()

	var msg *share_message.IMmessage
	queryBson := bson.M{"_id": id}
	if len(status) > 0 {
		queryBson["Status"] = status[0]
	}
	query := col.Find(queryBson)
	err := query.One(&msg)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return msg
}

//保存消息到数据库
func SaveMessageToDB(reqMsg *share_message.IMmessage) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_WAITER_MESSAGE)
	defer closeFun()

	setBson := bson.M{"Snew": reqMsg.GetSnew(), "Cnew": reqMsg.GetCnew(), "UpdateTime": reqMsg.GetUpdateTime()}
	if reqMsg.WaiterId != nil && reqMsg.GetWaiterId() != 0 {
		setBson["WaiterId"] = reqMsg.GetWaiterId()
	}
	if reqMsg.WaiterName != nil && reqMsg.GetWaiterName() != "" {
		setBson["WaiterName"] = reqMsg.GetWaiterName()
	}

	content := *reqMsg.Content[0]
	pushBson := bson.M{"Content": content}

	im := &share_message.IMmessage{}
	errc := col.Find(bson.M{"_id": reqMsg.GetId()}).One(&im)
	if errc != nil && errc != mgo.ErrNotFound {
		panic(errc)
	}
	if errc == mgo.ErrNotFound {
		err := col.Insert(reqMsg)
		easygo.PanicError(err)
	} else {
		err := col.Update(bson.M{"_id": reqMsg.GetId()}, bson.M{"$push": pushBson, "$set": setBson}) //向数据文档中追加数据，并改变消息阅读状态为未读
		easygo.PanicError(err)
	}
}

//更新消息的已读状态 types 1 玩家已读， 2 客服已读
func UpdateMessageRead(mid int64, types int32) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_WAITER_MESSAGE)
	defer closeFun()

	setBson := bson.M{}
	switch types {
	case 1:
		setBson["Cnew"] = 0
	case 2:
		setBson["Snew"] = 0
	}
	err := col.Update(bson.M{"_id": mid}, bson.M{"$set": setBson})
	easygo.PanicError(err)
}

//给客服评价
func GradeToWaiter(mid int64, grade int32) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_WAITER_MESSAGE)
	defer closeFun()

	status := WAITER_MESSAGE_GRADE
	setBson := bson.M{"Status": status, "Grade": grade, "UpdateTime": util.GetMilliTime()}

	err := col.Update(bson.M{"_id": mid}, bson.M{"$set": setBson})
	easygo.PanicError(err)
}

//结束消息
func OverWaiterMessage(mid int64) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_WAITER_MESSAGE)
	defer closeFun()

	setBson := bson.M{"Status": WAITER_MESSAGE_END, "UpdateTime": util.GetMilliTime()}

	err := col.Update(bson.M{"_id": mid}, bson.M{"$set": setBson})
	easygo.PanicError(err)
}

//客服绩效查询
func QueryWaiterPerformance(waiterId int64) *share_message.WaiterPerformance {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WAITER_PERFORMANCE)
	defer closeFun()

	siteOne := &share_message.WaiterPerformance{}
	err := col.Find(bson.M{"WaiterId": waiterId, "CreateTime": easygo.GetMonth0ClockTimestamp()}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//修改客服绩效
func UpdateWaiterPerformance(waiterId int64, grade int32) {
	wf := QueryWaiterPerformance(waiterId)
	if wf == nil {
		wf = &share_message.WaiterPerformance{
			Id:         easygo.NewInt64(NextId(TABLE_WAITER_PERFORMANCE)),
			CreateTime: easygo.NewInt64(easygo.GetMonth0ClockTimestamp()),
			WaiterId:   easygo.NewInt64(waiterId),
			ConNum:     easygo.NewInt32(1),     //接待次数
			GradeNum:   easygo.NewInt32(1),     //评分次数
			SumGrade:   easygo.NewInt32(grade), //总分
		}
	} else {
		wf.ConNum = easygo.NewInt32(wf.GetConNum() + 1)
		wf.GradeNum = easygo.NewInt32(wf.GetGradeNum() + 1)
		wf.SumGrade = easygo.NewInt32(wf.GetSumGrade() + grade)
	}

	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WAITER_PERFORMANCE)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": wf.GetId()}, bson.M{"$set": wf})
	easygo.PanicError(err)
}

//常见问题搜索
func QueryWaiterFAQList(key string, types int) []*share_message.WaiterFAQ {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WAITER_FAQ)
	defer closeFun()

	// pageSize := int(reqMsg.GetPageSize())
	// curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{"Status": 1}
	if key != "" {
		switch types {
		case FAQ_TITLE_SEARCH:
			queryBson["Title"] = key
		case FAQ_KEY_SEARCH:
			var keys []string
			keys = append(keys, key)
			queryBson["KeyWord"] = bson.M{"$in": keys}
		}
	}

	var list []*share_message.WaiterFAQ
	query := col.Find(queryBson)
	errc := query.Sort("-Clicks").Skip(0).Limit(10).All(&list)
	easygo.PanicError(errc)

	return list
}

//查询常见问题
func OpenFaqById(id int32) *share_message.WaiterFAQ {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WAITER_FAQ)
	defer closeFun()

	one := &share_message.WaiterFAQ{}
	err := col.FindId(id).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	err = col.Update(bson.M{"_id": id}, bson.M{"$inc": bson.M{"Clicks": 1}})
	easygo.PanicError(err)

	return one

}
