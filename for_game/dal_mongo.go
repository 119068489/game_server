package for_game

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/client_hall"
	"game_server/pb/client_login"
	"game_server/pb/client_server"
	"game_server/pb/share_message"
	"log"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

var _ = fmt.Sprintf
var _ = log.Println
var _ = easygo.Underline

var REDIS_LABEL_INFO = "redis_label_info"

//数据库管理工具方法=======================================================

//复杂创建索引
func EnsureIndex(dbName, tableName string, keys []string) {
	col, closeFun := easygo.MongoMgr.GetC(dbName, tableName)
	defer closeFun()

	index := mgo.Index{
		Key:        keys,
		Unique:     true, //设置成true，每个document只能有一个索引
		DropDups:   true, //设置成true，如果有文档的索引和之前的存在一样就会被删除，而不是报错
		Background: true, //设置为true，其他的连接可以访问集合，即使索引还没有建立，但是建立索引的连接是阻塞的，直到索引被建立完。
		Sparse:     true, //设置为true，则只有包含提供的关键字段的文档才会被包含在索引中。当使用稀疏索引进行排序时，只会返回已索引的文档。
	}
	err := col.EnsureIndex(index)
	easygo.PanicError(err)
}

//简单创建索引
func EnsureIndexKey(mongo easygo.IMongoDBManager, dbName, tableName string, key ...string) {
	col, closeFun := mongo.GetC(dbName, tableName)
	defer closeFun()
	err := col.EnsureIndexKey(key...)
	easygo.PanicError(err)
}

//查询当前表的现有索引
func Indexes(dbName, tableName string) []mgo.Index {
	col, closeFun := easygo.MongoMgr.GetC(dbName, tableName)
	defer closeFun()
	indexs, err := col.Indexes()
	easygo.PanicError(err)
	return indexs
}

//删除索引
func DropIndex(dbName, tableName string, key []string) {
	col, closeFun := easygo.MongoMgr.GetC(dbName, tableName)
	defer closeFun()
	err := col.DropIndex(key...)
	easygo.PanicError(err)
}

//查询某个表的某个字段集合 返回查询的全部非空字段值   for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE, "_id", bson.M{})
func Distinct(dbName, tableName, field string, find bson.M) []interface{} {
	var col *mgo.Collection
	var closeFun func()

	if dbName == MONGODB_NINGMENG_LOG {
		col, closeFun = easygo.MongoLogMgr.GetC(dbName, tableName)
	} else {
		col, closeFun = easygo.MongoMgr.GetC(dbName, tableName)
	}
	defer closeFun()
	var result []interface{}
	err := col.Find(find).Distinct(field, &result)
	easygo.PanicError(err)
	return result
}

//通用一步查询修改数据方法
func FindAndModify(dbName, tableName string, find, update interface{}, upsert bool) interface{} {
	var col *mgo.Collection
	var closeFun func()
	if dbName == MONGODB_NINGMENG_LOG {
		col, closeFun = easygo.MongoLogMgr.GetC(dbName, tableName)
	} else {
		col, closeFun = easygo.MongoMgr.GetC(dbName, tableName)
	}
	defer closeFun()
	var returnStruct interface{}
	_, err := col.Find(find).Apply(mgo.Change{
		Update:    update,
		Upsert:    upsert,
		ReturnNew: true,
	}, &returnStruct)

	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return returnStruct
}

//查询一个数据对象
func FindOne(dbName, tableName string, find bson.M) interface{} {
	var col *mgo.Collection
	var closeFun func()
	if dbName == MONGODB_NINGMENG_LOG {
		col, closeFun = easygo.MongoLogMgr.GetC(dbName, tableName)
	} else {
		col, closeFun = easygo.MongoMgr.GetC(dbName, tableName)
	}
	defer closeFun()
	var returnStruct interface{}
	err := col.Find(find).One(&returnStruct)

	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return returnStruct
}

func DelAllMgo(dbName, tableName string, delBson bson.M) {
	var col *mgo.Collection
	var closeFun func()
	if dbName == MONGODB_NINGMENG_LOG {
		col, closeFun = easygo.MongoLogMgr.GetC(dbName, tableName)
	} else {
		col, closeFun = easygo.MongoMgr.GetC(dbName, tableName)
	}
	defer closeFun()
	_, err := col.RemoveAll(delBson)
	easygo.PanicError(err)
}

func InsertMgo(dbName, tableName string, ls interface{}) error {
	var col *mgo.Collection
	var closeFun func()
	if dbName == MONGODB_NINGMENG_LOG {
		col, closeFun = easygo.MongoLogMgr.GetC(dbName, tableName)
	} else {
		col, closeFun = easygo.MongoMgr.GetC(dbName, tableName)
	}
	defer closeFun()
	// 批量插入数据库
	err := col.Insert(ls)
	return err
}

func InsertAllMgo(dbName, tableName string, ls ...interface{}) {
	var col *mgo.Collection
	var closeFun func()
	if dbName == MONGODB_NINGMENG_LOG {
		col, closeFun = easygo.MongoLogMgr.GetC(dbName, tableName)
	} else {
		col, closeFun = easygo.MongoMgr.GetC(dbName, tableName)
	}
	defer closeFun()
	// 批量插入数据库
	err := col.Insert(ls...)
	easygo.PanicError(err)
}

func UpdateAllMgo(dbName, tableName string, findBson, updateBson bson.M) (*mgo.ChangeInfo, error) {
	var col *mgo.Collection
	var closeFun func()
	if dbName == MONGODB_NINGMENG_LOG {
		col, closeFun = easygo.MongoLogMgr.GetC(dbName, tableName)
	} else {
		col, closeFun = easygo.MongoMgr.GetC(dbName, tableName)
	}
	defer closeFun()
	info, err := col.UpdateAll(findBson, updateBson)
	easygo.PanicError(err)
	return info, err
}

//查询所有数据对象
func FindAll(dbName, tableName string, find bson.M, pageSize, curPage int, sort ...string) ([]interface{}, int) {
	var col *mgo.Collection
	var closeFun func()
	if dbName == MONGODB_NINGMENG_LOG {
		col, closeFun = easygo.MongoLogMgr.GetC(dbName, tableName)
	} else {
		col, closeFun = easygo.MongoMgr.GetC(dbName, tableName)
	}
	defer closeFun()
	var returnStruct []interface{}
	query := col.Find(find)
	count, err := query.Count()
	easygo.PanicError(err)
	curPage = easygo.If(curPage > 0, curPage-1, 0).(int)
	query.Skip(curPage * pageSize).Limit(pageSize)
	if len(sort) > 0 && sort[0] != "" {
		query.Sort(sort...)
	}
	errc := query.All(&returnStruct)
	easygo.PanicError(errc)

	return returnStruct, count
}

//管道查询所有数据对象
func FindPipeAll(dbName, tableName string, m []bson.M, pageSize, curPage int) []interface{} {
	var col *mgo.Collection
	var closeFun func()
	if dbName == MONGODB_NINGMENG_LOG {
		col, closeFun = easygo.MongoLogMgr.GetC(dbName, tableName)
	} else {
		col, closeFun = easygo.MongoMgr.GetC(dbName, tableName)
	}
	defer closeFun()
	if curPage > 0 {
		curPage = easygo.If(curPage > 0, curPage-1, 0).(int)
		m = append(m, bson.M{"$skip": curPage * pageSize})
	}
	if pageSize > 0 {
		m = append(m, bson.M{"$limit": pageSize})
	}

	var returnStruct []interface{}
	query := col.Pipe(m)
	err := query.All(&returnStruct)
	easygo.PanicError(err)

	return returnStruct
}

//管道查询单个数据对象
func FindPipeOne(dbName, tableName string, m []bson.M) interface{} {
	var col *mgo.Collection
	var closeFun func()
	if dbName == MONGODB_NINGMENG_LOG {
		col, closeFun = easygo.MongoLogMgr.GetC(dbName, tableName)
	} else {
		col, closeFun = easygo.MongoMgr.GetC(dbName, tableName)
	}
	defer closeFun()

	var returnStruct interface{}
	query := col.Pipe(m)
	err := query.One(&returnStruct)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return returnStruct
}

//查询需要的聚合数量
func FindPipeAllCount(dbName, tableName string, m []bson.M) int64 {
	var col *mgo.Collection
	var closeFun func()
	if dbName == MONGODB_NINGMENG_LOG {
		col, closeFun = easygo.MongoLogMgr.GetC(dbName, tableName)
	} else {
		col, closeFun = easygo.MongoMgr.GetC(dbName, tableName)
	}
	defer closeFun()

	query := col.Pipe(m)
	var one *share_message.PipeIntCount
	err := query.One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return 0
	}

	return one.GetCount()
}

//查询对象数量
func FindAllCount(dbName, tableName string, find bson.M) int {
	var col *mgo.Collection
	var closeFun func()
	if dbName == MONGODB_NINGMENG_LOG {
		col, closeFun = easygo.MongoLogMgr.GetC(dbName, tableName)
	} else {
		col, closeFun = easygo.MongoMgr.GetC(dbName, tableName)
	}
	defer closeFun()

	query := col.Find(find)
	count, err := query.Count()
	easygo.PanicError(err)

	return count
}

//初始化通用方法 dbName(数据库名), tableName(表名),isIncrease(是否自增Id), reason(初始化说明)
func InitToMongo(dbName, tableName, reason string, isIncrease bool, queryBson bson.M, req []interface{}) {
	var col *mgo.Collection
	var closeFun func()
	if dbName == MONGODB_NINGMENG_LOG {
		col, closeFun = easygo.MongoLogMgr.GetC(dbName, tableName)
	} else {
		col, closeFun = easygo.MongoMgr.GetC(dbName, tableName)
	}
	defer closeFun()
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	//新增初始化数据
	if count == 0 {
		if len(req) > 0 {
			err := col.Insert(req...)
			easygo.PanicError(err)

			if isIncrease {
				DelAllMgo(MONGODB_NINGMENG, TABLE_ID_GENERATOR, bson.M{"_id": tableName})

				identity := &Identity{
					Key:   easygo.NewString(tableName),
					Value: easygo.NewUint64(len(req)),
				}
				InsertAllMgo(MONGODB_NINGMENG, TABLE_ID_GENERATOR, identity)
			}

			s := fmt.Sprintf(reason+"[%s]DB", dbName)
			logs.Info(s)
		} else {
			if tableName == TABLE_COIN_PRODUCT {
				req = InitShopProductItems()
				err := col.Insert(req...)
				easygo.PanicError(err)
			}
		}
		return
	}
	//更新初始化数据
	if count < len(req) {
		if tableName == TABLE_SOURCETYPE {
			_, err := col.RemoveAll(bson.M{})
			easygo.PanicError(err)

			err = col.Insert(req...)
			easygo.PanicError(err)

			if isIncrease {
				err := col.RemoveId(tableName)
				easygo.PanicError(err)

				identity := &Identity{
					Key:   easygo.NewString(tableName),
					Value: easygo.NewUint64(len(req)),
				}
				var identitys []interface{}
				identitys = append(identitys, identity)

				InsertAllMgo(MONGODB_NINGMENG, TABLE_ID_GENERATOR, identitys)
			}

			s := fmt.Sprintf("更新"+reason+"[%s]DB", dbName)
			logs.Info(s)
		}
	}
}

func InitShopProductItems() []interface{} {
	its := GetPropsItemsCfg()
	var items []interface{}
	for _, it := range its.GetItems() {
		id := NextId(TABLE_COIN_PRODUCT)
		product1 := &share_message.CoinProduct{
			Id:            easygo.NewInt64(id),
			Coin:          easygo.NewInt64(700),
			CreateTime:    easygo.NewInt64(GetMillSecond()),
			EffectiveTime: easygo.NewInt64(7),
			Name:          easygo.NewString(it.GetName()),
			Price:         easygo.NewInt64(0),
			ProductNum:    easygo.NewInt64(-1),
			PropsIcon:     easygo.NewString("https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1607407809000.png"),
			PropsId:       easygo.NewInt64(it.GetId()),
			PropsType:     easygo.NewInt32(it.GetPropsType()),
			Sort:          easygo.NewInt32(id),
			Status:        easygo.NewInt32(2), //1上架，2下架
		}
		id = NextId(TABLE_COIN_PRODUCT)
		product2 := &share_message.CoinProduct{
			Id:            easygo.NewInt64(NextId(TABLE_COIN_PRODUCT)),
			Coin:          easygo.NewInt64(2800),
			CreateTime:    easygo.NewInt64(GetMillSecond()),
			EffectiveTime: easygo.NewInt64(-1),
			Name:          easygo.NewString(it.GetName()),
			Price:         easygo.NewInt64(0),
			ProductNum:    easygo.NewInt64(-1),
			PropsIcon:     easygo.NewString("https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1607407809000.png"),
			PropsId:       easygo.NewInt64(it.GetId()),
			PropsType:     easygo.NewInt32(it.GetPropsType()),
			Sort:          easygo.NewInt32(id),
			Status:        easygo.NewInt32(2), //1上架，2下架
		}
		items = append(items, product1, product2)
	}
	return items
}

//生成金币变化类型
func InitSourceType(req []*share_message.SourceType) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_SOURCETYPE)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)

	var il []interface{}
	for _, rr := range req {
		il = append(il, rr)
	}

	if count == 0 {
		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)

		s := fmt.Sprintf("初始化金币变化源类型到[%s]DB", share)
		logs.Info(s)
	}
}

//生成系统参数
func InitSysParameter(req []*share_message.SysParameter) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_SYS_PARAMETER)
	defer closeFun()
	var params []*share_message.SysParameter
	err := col.Find(bson.M{}).All(&params)
	easygo.PanicError(err)
	count, err := col.Count()
	easygo.PanicError(err)

	var il []interface{}
	for _, rr := range req {
		il = append(il, rr)
	}
	if count == 0 {
		//DeleteSysParameter()
		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)

		s := fmt.Sprintf("初始化系统参数到[%s]DB", share)
		logs.Info(s)
	} else if len(req) > count {
		var temList []string
		for _, p := range params {
			temList = append(temList, p.GetId())
		}
		//新增的加入到数据库
		for _, p := range req {
			if !easygo.Contain(temList, p.GetId()) {
				_ = col.Insert(p) //批量插入到数据库
			}
		}
	}
}

//生成硬币生成道具
func InitCoinShopPropsItems(req []*share_message.PropsItem) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_PROPS_ITEM)
	defer closeFun()
	var params []*share_message.PropsItem
	err := col.Find(bson.M{}).All(&params)
	easygo.PanicError(err)
	count, err := col.Count()
	easygo.PanicError(err)

	var il []interface{}
	for _, rr := range req {
		il = append(il, rr)
	}
	if count == 0 {
		if len(il) == 0 {
			logs.Error("获取道具库初始化数据失败")
		} else {
			//DeleteSysParameter()
			err := col.Insert(il...) //批量插入到数据库
			easygo.PanicError(err)

			s := fmt.Sprintf("初始化硬币商场道具到[%s]DB", share)
			logs.Info(s)
		}
	}
}

//初始化标签
func InitInterestTag(req []*share_message.InterestTag) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_INTERESTTAG)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)

	if count == 0 {
		var il []interface{}
		for _, rr := range req {
			if rr.Id == nil && rr.GetId() == 0 {
				rr.Id = easygo.NewInt32(NextId(TABLE_INTERESTTAG))
			}
			il = append(il, rr)
		}

		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)

		s := fmt.Sprintf("初始化标签到[%s]DB", share)
		logs.Info(s)
	}
}

//初始化话题类别
func InitTopicType(req []*share_message.TopicType) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_TOPIC_TYPE)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)

	if count == 0 {
		var il []interface{}
		for _, rr := range req {
			if rr.Id == nil && rr.GetId() == 0 {
				rr.Id = easygo.NewInt64(NextId(TABLE_TOPIC_TYPE))
			}
			il = append(il, rr)
		}

		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)

		s := fmt.Sprintf("初始化话题类别[%s]DB", share)
		logs.Info(s)
	}
}

//初始化客服类型
func InitManageTypes(req []*share_message.ManagerTypes) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_MANAGER_TYPES)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)

	if count == 0 {
		var il []interface{}
		for _, rr := range req {
			if rr.Id == nil && rr.GetId() == 0 {
				rr.Id = easygo.NewInt32(NextId(TABLE_MANAGER_TYPES))
			}
			il = append(il, rr)
		}

		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)

		s := fmt.Sprintf("初始化客服类型到[%s]DB", share)
		logs.Info(s)
	}
}

//通用额度配置
func InitGeneralQuota(req *share_message.GeneralQuota) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_GENERAL_QUOTA)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)

	if count == 0 {
		var il []interface{}
		il = append(il, req)

		err := col.Insert(il...) //插入到数据库
		easygo.PanicError(err)

		s := fmt.Sprintf("初始化通用额度配置到[%s]DB", share)
		logs.Info(s)
	}
}

//初始化支付类型
func InitPayType(req []*share_message.PayType) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_PAYTYPE)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)

	if count == 0 {
		var il []interface{}
		for _, rr := range req {
			if rr.Id == nil && rr.GetId() == 0 {
				rr.Id = easygo.NewInt32(NextId(TABLE_PAYTYPE))
			}
			il = append(il, rr)
		}

		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)

		s := fmt.Sprintf("初始化支付类型到[%s]DB", share)
		logs.Info(s)
	}
}

//初始化支付场景
func InitPayScene(req []*share_message.PayScene) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_PAYSCENE)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)

	if count == 0 {
		var il []interface{}
		for _, rr := range req {
			if rr.Id == nil && rr.GetId() == 0 {
				rr.Id = easygo.NewInt32(NextId(TABLE_PAYSCENE))
			}
			il = append(il, rr)
		}

		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)

		s := fmt.Sprintf("初始化支付场景到[%s]DB", share)
		logs.Info(s)
	}
}

//初始化支付设定
func InitPaySetting(req []*share_message.PaymentSetting) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_PAYMENTSETTING)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)

	if count == 0 {
		var il []interface{}
		for _, rr := range req {
			if rr.Id == nil && rr.GetId() == 0 {
				rr.Id = easygo.NewInt32(NextId(TABLE_PAYMENTSETTING))
			}

			il = append(il, rr)
		}

		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)

		s := fmt.Sprintf("初始化支付管理设置到[%s]DB", share)
		logs.Info(s)
	}
}

//初始化支付平台
func InitPaymentPlatform(req []*share_message.PaymentPlatform) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_PAYMENTPLATFORM)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)

	if count == 0 {
		var il []interface{}
		for _, rr := range req {
			if rr.Id == nil && rr.GetId() == 0 {
				rr.Id = easygo.NewInt32(NextId(TABLE_PAYMENTPLATFORM))
			}
			il = append(il, rr)
		}

		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)

		s := fmt.Sprintf("初始化支付平台到[%s]DB", share)
		logs.Info(s)
	}
}

//初始化支付平台通道
func InitPlatformChannel(req []*share_message.PlatformChannel) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_PLATFORM_CHANNEL)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)

	if count == 0 {
		var il []interface{}
		for _, rr := range req {
			if rr.Id == nil && rr.GetId() == 0 {
				rr.Id = easygo.NewInt32(NextId(TABLE_PLATFORM_CHANNEL))
			}
			il = append(il, rr)
		}

		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)

		s := fmt.Sprintf("初始化支付平台通道到[%s]DB", share)
		logs.Info(s)
	}
}

//初始化活动
func InitActivity(req []*share_message.Activity) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_LUCKY_ACTIVITY)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)

	if count == 0 {
		var il []interface{}
		for _, rr := range req {
			if rr.Id == nil && rr.GetId() == 0 {
				rr.Id = easygo.NewInt64(NextId(TABLE_LUCKY_ACTIVITY))
			}

			il = append(il, rr)
		}

		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)

		s := fmt.Sprintf("初始化活动设置到[%s]DB", share)
		logs.Info(s)
	}
}

//初始化道具
func InitProps(req []*share_message.Props) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_LUCKY_PROPS)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)

	if count == 0 {
		var il []interface{}
		for _, rr := range req {
			if rr.Id == nil && rr.GetId() == 0 {
				rr.Id = easygo.NewInt64(NextId(TABLE_LUCKY_PROPS))
			}

			il = append(il, rr)
		}

		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)

		s := fmt.Sprintf("初始化道具设置到[%s]DB", share)
		logs.Info(s)
	}
}

//初始化抽卡概率
func InitPropsRate(req []*share_message.PropsRate) {
	share := MONGODB_NINGMENG
	col, closeFun := easygo.MongoMgr.GetC(share, TABLE_LUCKY_PROPS_RATE)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)

	if count == 0 {
		var il []interface{}
		for _, rr := range req {
			il = append(il, rr)
		}

		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)

		s := fmt.Sprintf("初始化抽卡概率设置到[%s]DB", share)
		logs.Info(s)
	}
}

//创建账号 types用户类型，1普通用户，2陪聊用户   isOnline 是否在线
func CreateAccount(data *share_message.CreateAccountData) (bool, PLAYER_ID) {
	var pwd string
	if data.GetPassWord() != "" { //如果密码是不空字符串代表着是后台创建用户
		pwd = Md5(data.GetPassWord())
	} else {
		pwd = Md5("0000")
	}
	playerId := NextId(TABLE_PLAYER_ACCOUNT)
	playerAccount := &share_message.PlayerAccount{
		PlayerId:    easygo.NewInt64(playerId),
		Account:     easygo.NewString(data.GetPhone()),
		Email:       easygo.NewString(""),
		Password:    easygo.NewString(pwd),
		PayPassword: easygo.NewString(""),
		CreateTime:  easygo.NewInt64(GetMillSecond()),
		AreaCode:    easygo.NewString(data.GetAreaCode()),
		IsBind:      easygo.NewBool(true), //新建的账号不给再绑
	}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_ACCOUNT)
	defer closeFun()
	//账号表注册
	_, err := col.Upsert(bson.M{"_id": playerId}, bson.M{"$set": playerAccount})
	easygo.PanicError(err)
	CreateRedisAccount(playerAccount)
	b := CreatePlayer(playerId, data)
	//等待创建完角色再返回
	if b {
		period := GetPlayerPeriod(playerId)
		period.HaltYearPeriod.Set(CHANGE_PHONE, true)
	}
	return b, playerId
}

//微信登录创建账号
func CreateAccountForWechat(account, name, headIcon, unionId string, sex int32, areaCode string) PLAYER_ID {
	playerId := NextId(TABLE_PLAYER_ACCOUNT)
	playerAccount := &share_message.PlayerAccount{
		PlayerId: easygo.NewInt64(playerId),
		Account:  easygo.NewString(account),
		Email:    easygo.NewString(""),
		Password: easygo.NewString(Md5("0000")),
		//Token:       easygo.NewString(""),
		PayPassword: easygo.NewString(""),
		OpenId:      easygo.NewString(""),
		UnionId:     easygo.NewString(unionId),
		CreateTime:  easygo.NewInt64(GetMillSecond()),
		AreaCode:    easygo.NewString(areaCode),
		IsBind:      easygo.NewBool(true), //新建的账号不给再绑
	}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_ACCOUNT)
	defer closeFun()
	//账号表注册
	_, err := col.Upsert(bson.M{"_id": playerId}, bson.M{"$set": playerAccount})
	easygo.PanicError(err)
	//redis存储账号信息
	CreateRedisAccount(playerAccount)
	lmAccount := GetRandAccount("lm", playerId)
	//数据库写入base64编码名字
	//	bs64Name := base64.StdEncoding.EncodeToString([]byte(name))
	player := &share_message.PlayerBase{
		PlayerId:        easygo.NewInt64(playerId),
		Account:         easygo.NewString(lmAccount),
		Email:           easygo.NewString(""),
		NickName:        easygo.NewString(name),
		HeadIcon:        easygo.NewString(headIcon),
		Sex:             easygo.NewInt32(sex),
		Gold:            easygo.NewInt64(0),
		IsRobot:         easygo.NewBool(false),
		PeopleId:        easygo.NewString(""),
		Phone:           easygo.NewString(account),
		RealName:        easygo.NewString(""),
		IsOnline:        easygo.NewBool(true),
		Signature:       easygo.NewString(""),
		Provice:         easygo.NewString(""),
		City:            easygo.NewString(""),
		PlayerSetting:   GetPlayerDefalutSetting(),
		IsNearBy:        easygo.NewBool(false),
		Types:           easygo.NewInt32(1),
		IsRecommend:     easygo.NewBool(false),
		Channel:         easygo.NewString(""),
		IsRecommendOver: easygo.NewBool(true),
		AreaCode:        easygo.NewString(areaCode),
	}

	col1, closeFun1 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun1()
	//账号表注册
	_, err3 := col1.Upsert(bson.M{"_id": playerId}, bson.M{"$set": player})
	easygo.PanicError(err3)
	data := easygo.KWAT{
		"PlayerId": playerId,
	}
	GetFriendBase(playerId, data) //创建新的好友对象
	return playerId
}

func GetPlayerDefalutSetting() *share_message.PlayerSetting {
	msg := &share_message.PlayerSetting{
		IsSafePassword: easygo.NewBool(false),
		IsNewMessage:   easygo.NewBool(true),
		IsMusic:        easygo.NewBool(true),
		IsShake:        easygo.NewBool(true),
		IsAddFriend:    easygo.NewBool(true),
		IsPhone:        easygo.NewBool(true),
		IsAccount:      easygo.NewBool(true),
		IsTeamChat:     easygo.NewBool(true),
		IsCode:         easygo.NewBool(true),
		IsCard:         easygo.NewBool(true),
		IsSafeProtect:  easygo.NewBool(false),
		SafePassword:   easygo.NewString(""),
		IsTouch:        easygo.NewBool(false),
		IsMessageShow:  easygo.NewBool(true),
	}
	return msg
}

//创建角色  types用户类型，1普通用户，2陪聊用户
func CreatePlayer(playerId int64, data *share_message.CreateAccountData) bool {
	lmAccount := GetRandAccount("lm", playerId)
	ctype := int32(1)
	if data.GetTypes() > 0 {
		ctype = data.GetTypes()
	}
	player := &share_message.PlayerBase{
		PlayerId:        easygo.NewInt64(playerId),
		Account:         easygo.NewString(lmAccount),
		Email:           easygo.NewString(""),
		NickName:        easygo.NewString(""),
		HeadIcon:        easygo.NewString(""),
		Sex:             easygo.NewInt32(0),
		Gold:            easygo.NewInt64(0),
		IsRobot:         easygo.NewBool(false),
		PeopleId:        easygo.NewString(""),
		Phone:           easygo.NewString(data.GetPhone()),
		RealName:        easygo.NewString(""),
		IsOnline:        easygo.NewBool(data.GetIsOnline()),
		Signature:       easygo.NewString(""),
		Provice:         easygo.NewString(""),
		City:            easygo.NewString(""),
		PlayerSetting:   GetPlayerDefalutSetting(),
		IsNearBy:        easygo.NewBool(false),
		Types:           easygo.NewInt32(ctype),
		IsRecommend:     easygo.NewBool(false),
		Channel:         easygo.NewString(""),
		CreateIP:        easygo.NewString(data.GetIp()),
		IsRecommendOver: easygo.NewBool(true),
		IsVisitor:       easygo.NewBool(data.GetIsVisitor()),
		AreaCode:        easygo.NewString(data.GetAreaCode()),
	}

	//TODO:h5 邮箱注册临时对应
	if strings.Contains(data.GetPhone(), "@") {
		player.Phone = easygo.NewString("")
		player.Email = easygo.NewString(data.GetPhone())
	}

	if ctype > 1 {
		player.Sex = easygo.NewInt32(2)
		player.IsRecommendOver = easygo.NewBool(false)
	}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	//账号表注册
	_, err := col.Upsert(bson.M{"_id": playerId}, bson.M{"$set": player})
	easygo.PanicError(err)

	fData := easygo.KWAT{
		"PlayerId": playerId,
	}
	GetFriendBase(playerId, fData) //创建新的好友对象
	return true
}

//根据IP查询冻结IP列表
func GetFreezeByIp(ip string) *share_message.FreezeIpList {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_FREEZEIP)
	defer closeFun()

	freeze := &share_message.FreezeIpList{}
	err := col.Find(bson.M{"LoginIP": ip}).One(&freeze)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return freeze
}

func LoginIpAuth(ip string) bool {
	freeInfo := GetFreezeByIp(ip)
	if freeInfo != nil {
		return freeInfo.GetLoginAuth()
	}
	return true
}

func CreateIpAuth(ip string) bool { //IP是否能创建角色
	freeInfo := GetFreezeByIp(ip)
	if freeInfo != nil {
		return freeInfo.GetRegisterAuth()
	}
	return true
}

//根据IP查询冻结IP列表
func GetFreezeByAccount(account string) *share_message.FreezeAccountList {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_FREEZEACCOUNT)
	defer closeFun()

	freeze := &share_message.FreezeAccountList{}
	err := col.Find(bson.M{"Account": account}).One(&freeze)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return freeze
}

func AddFreezeAccount(account string) {
	Id := NextId(TABLE_FREEZEACCOUNT)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_FREEZEACCOUNT)
	defer closeFun()

	freeze := &share_message.FreezeAccountList{
		Id:        easygo.NewInt64(Id),
		Account:   easygo.NewString(account),
		LoginAuth: easygo.NewBool(true),
	}
	err := col.Insert(freeze)
	easygo.PanicError(err)
}

func DelFreezeAccount(account string) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_FREEZEACCOUNT)
	defer closeFun()

	err := col.Remove(bson.M{"Account": account})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
}

func LoginAccountAuth(account string) bool {
	freeInfo := GetFreezeByAccount(account)
	if freeInfo != nil {
		return freeInfo.GetLoginAuth()
	}
	return true
}

//设备类型查询用户
func QueryPlayersByDeviceType(deviceType int32) []*share_message.PlayerBase {
	var lst []*share_message.PlayerBase
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	//只有正常的用户才推送
	queryBson := bson.M{"Status": ACCOUNT_NORMAL}
	if deviceType > 0 {
		queryBson["DeviceType"] = deviceType
	}
	err := col.Find(queryBson).All(&lst)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}

	//ids := []int64{}
	//for _, value := range lst {
	//	ids = append(ids, value.GetPlayerId())
	//}
	//
	//result := QueryPlayersByIds(ids)
	return lst
}

//设备号查询玩家
func QueryPlayersByDeviceCode(deviceCode string) []*share_message.PlayerBase {
	var lst []*share_message.PlayerBase
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	queryBson := bson.M{"DeviceCode": deviceCode}
	err := col.Find(queryBson).All(&lst)
	easygo.PanicError(err)

	return lst
}

//查询推送对象列表
func QueryPlayersByOfPush(deviceType int32, lable []int32, customTag []int32, grabTag []int32) []*share_message.PlayerBase {
	logs.Info("lable:", lable)
	logs.Info("customTag:", customTag)
	logs.Info("grabTag:", grabTag)
	var lst []*share_message.PlayerBase
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	q := make([]bson.M, 0)
	if len(lable) > 0 {
		q = append(q, bson.M{"Label": bson.M{"$elemMatch": bson.M{"$in": lable}}})
	}
	if len(customTag) > 0 {
		q = append(q, bson.M{"CustomTag": bson.M{"$elemMatch": bson.M{"$in": customTag}}})
	}
	if len(grabTag) > 0 {
		q = append(q, bson.M{"GrabTag": bson.M{"$in": grabTag}})
	}
	queryBson := bson.M{}
	if len(q) > 0 {
		queryBson = bson.M{"$or": q}
	}
	if deviceType > 0 {
		queryBson["DeviceType"] = deviceType
	}

	err := col.Find(queryBson).All(&lst)
	easygo.PanicError(err)

	return lst
}

//根据playerid查询玩家信息
func QueryPlayersByIds(playerIds []int64) []*share_message.PlayerAccount {
	var playerAccount []*share_message.PlayerAccount
	col1, closeFun1 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_ACCOUNT)
	defer closeFun1()
	err := col1.Find(bson.M{"_id": bson.M{"$in": playerIds}}).Select(bson.M{"Token": 1}).All(&playerAccount)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	return playerAccount
}

func QueryPlayersById(playerId int64) *share_message.PlayerBase {
	var playerBase *share_message.PlayerBase
	col1, closeFun1 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun1()
	err := col1.Find(bson.M{"_id": playerId}).Select(bson.M{"Token": 1}).One(&playerBase)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	return playerBase
}

//生成唯一订单ID
// func CreateOrderId(order *share_message.Order) string {
// 	r := int64(RandInt(10000, 99999))
// 	orderNo := easygo.AnytoA(order.GetChangeType()) + easygo.AnytoA(order.GetSourceType()) + easygo.AnytoA(GetMillSecond()) + easygo.AnytoA(NextId(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_ORDER)+1000000000) + easygo.AnytoA(r)
// 	return orderNo
// }

//更新订单(下订单):
// func SetOrder(order *share_message.Order) string {
// 	var orderNo string
// 	if order.OrderId == nil || order.GetOrderId() == "" {
// 		orderNo = CreateOrderId(order)
// 	} else {
// 		orderNo = order.GetOrderId()
// 	}

// 	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ORDER)
// 	defer closeFun()
// 	_, err := col.Upsert(bson.M{"_id": orderNo}, bson.M{"$set": order})

// 	easygo.PanicError(err)
// 	return orderNo
// }

//用户Id查询订单列表
func GetOrderListByPlayerId(playerId int64) ([]*share_message.Order, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ORDER)
	defer closeFun()

	var list []*share_message.Order
	query := col.Find(bson.M{"PlayerId": playerId})
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.All(&list)
	easygo.PanicError(errc)

	return list, count
}

//订单号查询订单
// func GetOrderById(id string) *share_message.Order {
// 	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ORDER)
// 	defer closeFun()

// 	order := &share_message.Order{}
// 	err := col.Find(bson.M{"_id": id}).One(order)
// 	if err != nil && err != mgo.ErrNotFound {
// 		panic(err)
// 	}
// 	if err == mgo.ErrNotFound {
// 		return nil
// 	}
// 	return order
// }

func GetPlayerIdForAccount(account string) PLAYER_ID { //通过柠檬号搜索玩家id
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	var obj *share_message.PlayerBase
	err := col.Find(bson.M{"Account": account}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return 0
	}
	return obj.GetPlayerId()
}

func GetPlayerInfoForUnionId(unionId string) *share_message.PlayerAccount { //通过微信openid搜索玩家id
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_ACCOUNT)
	defer closeFun()
	var obj *share_message.PlayerAccount
	err := col.Find(bson.M{"UnionId": unionId}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

func GetGoldLogInfo(pid PLAYER_ID, page, num int) []*GoldLog {
	var lst []*GoldLog
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_GOLDCHANGELOG)
	defer closeFun()
	err := col.Find(bson.M{"PlayerId": pid}).Sort("-CreateTime").Skip((page - 1) * num).Limit(num).All(&lst)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	return lst
}

func GetLocationInfo(pid int64, msg *client_hall.LocationInfo) *client_hall.AllLocationPlayerInfo {
	t := msg.GetType()
	province := msg.GetProvince()
	city := msg.GetCity()
	area := msg.GetArea()
	x := msg.GetX()
	y := msg.GetY()
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	var lst []*share_message.PlayerBase
	var query bson.M
	val := []interface{}{0.0, nil}
	if t == 0 {
		query = bson.M{"_id": bson.M{"$ne": pid}, "Types": 1, "X": bson.M{"$nin": val}, "Y": bson.M{"$nin": val}}
	} else if t == 1 {
		query = bson.M{"_id": bson.M{"$ne": pid}, "Types": 1, "Sex": 1, "X": bson.M{"$nin": val}, "Y": bson.M{"$nin": val}}
	} else if t == 2 {
		query = bson.M{"_id": bson.M{"$ne": pid}, "Types": 1, "Sex": 2, "X": bson.M{"$nin": val}, "Y": bson.M{"$nin": val}}
	} else {
		panic(fmt.Sprintf("不明类型t：%d", t))
	}

	err := col.Find(query).All(&lst)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	type People struct {
		PlayerId PLAYER_ID
		Distance float64
	}
	Info := []People{} //附近的人位置信息
	for _, m := range lst {
		// 判断是否注销的账号
		if m.GetStatus() == ACCOUNT_CANCELED {
			continue
		}
		x1 := m.GetX()
		y1 := m.GetY()
		dis1 := GetDistance(y, y1, x, x1)
		n := People{
			PlayerId: m.GetPlayerId(),
			Distance: dis1,
		}

		Info = append(Info, n)
	}

	sort.Slice(Info, func(i, j int) bool {
		return Info[i].Distance < Info[j].Distance // 升序
	})

	player := GetFriendBase(pid)       //获取玩家好友数据
	friendlst := player.GetFriendIds() //获取玩家好友列表数据

	var pInfo []*client_hall.LocationPlayerInfo
	num := len(Info) //附近的人长度
	if num > 200 {   //如果大于200人就取200人
		plst := Info[:200]
		for _, m := range plst {
			id := m.PlayerId
			base := GetRedisPlayerBase(id)
			var b bool
			if util.Int64InSlice(id, friendlst) { //如果好友出现在附近的人
				b = true
			}
			photolst := base.GetPhoto() //获取附近的人相册列表
			var photo string
			if len(photolst) != 0 {
				photo = photolst[0] //拿第一张照片
			}
			var dis float64
			if m.Distance < 100 { //如果距离小于100
				dis = 100
			} else {
				dis = m.Distance //否则采用真实距离
			}
			msg := &client_hall.LocationPlayerInfo{
				PlayerId:  easygo.NewInt64(id),
				Distance:  easygo.NewFloat64(dis),
				Name:      easygo.NewString(base.GetNickName()),
				Sex:       easygo.NewInt32(base.GetSex()),
				HeadIcon:  easygo.NewString(base.GetHeadIcon()),
				Signature: easygo.NewString(base.GetSignature()),
				IsFriend:  easygo.NewBool(b),
				Province:  easygo.NewString(base.GetProvice()),
				City:      easygo.NewString(base.GetCity()),
				Account:   easygo.NewString(base.GetAccount()),
				Photo:     easygo.NewString(photo),
				Area:      easygo.NewString(base.GetArea()),
				Types:     easygo.NewInt32(base.GetTypes()),
			}
			pInfo = append(pInfo, msg)
		}
	} else {
		rand.Seed(time.Now().Unix())
		manHeadList := GetManyRobotHeadIcon(num, 1)  //男随机头像列表
		girlHeadList := GetManyRobotHeadIcon(num, 2) //女头像随机列表
		manNameList := GetManyRobotName(1, num)
		girlNameList := GetManyRobotName(2, num)
		for i := 0; i < num; i++ { //增加机器人
			var sex int
			if t == 0 {
				sex = RandInt(1, 3)
			} else if t == 1 {
				sex = 1
			} else {
				sex = 2
			}
			var mark, name string
			var icon int
			if sex == 1 {
				mark = "mavatar"
				icon = manHeadList[i]
				name = manNameList[i]
			} else {
				mark = "wavatar"
				icon = girlHeadList[i]
				name = girlNameList[i]
			}
			head := fmt.Sprintf("https://im-resource-1253887233.file.myqcloud.com/prod/%s/%d.png", mark, icon)

			dis := RandInt(100, 10000)
			id := RandInt(Min_Robot_PlayerId, Max_Robot_PlayerId)
			account := GetRandAccount("lm", int64(id))
			labellst := GetRedisLabelInfo()
			labelIds := []int32{}
			if len(labellst) > 0 {
				if len(labellst) >= 2 {
					labelIds = append(labelIds, labellst[0].GetId(), labellst[1].GetId())
				} else {
					labelIds = append(labelIds, labellst[0].GetId())
				}
			}
			msg := &client_hall.LocationPlayerInfo{
				PlayerId:  easygo.NewInt64(id),
				Distance:  easygo.NewFloat64(dis),
				Name:      easygo.NewString(name),
				Sex:       easygo.NewInt32(sex),
				HeadIcon:  easygo.NewString(head),
				Signature: easygo.NewString(""),
				IsFriend:  easygo.NewBool(false),
				Province:  easygo.NewString(province),
				City:      easygo.NewString(city),
				Account:   easygo.NewString(account),
				Photo:     easygo.NewString(""),
				Area:      easygo.NewString(area),
				Types:     easygo.NewInt32(ACCOUNT_TYPES_PT), // 普通用户
			}
			if len(labelIds) > 0 {
				msg.LabelInfo = GetLabelInfo(labelIds)
			}
			pInfo = append(pInfo, msg)
		}
		for _, m := range Info { //获取数据库数据
			id := m.PlayerId
			base := GetRedisPlayerBase(id)
			var b bool
			if util.Int64InSlice(id, friendlst) {
				b = true
			}
			photolst := base.GetPhoto()
			var photo string
			if len(photolst) != 0 {
				photo = photolst[0]
			}

			msg := &client_hall.LocationPlayerInfo{
				PlayerId:  easygo.NewInt64(id),
				Distance:  easygo.NewFloat64(m.Distance),
				Name:      easygo.NewString(base.GetNickName()),
				Sex:       easygo.NewInt32(base.GetSex()),
				HeadIcon:  easygo.NewString(base.GetHeadIcon()),
				Signature: easygo.NewString(base.GetSignature()),
				IsFriend:  easygo.NewBool(b),
				Province:  easygo.NewString(base.GetProvice()),
				City:      easygo.NewString(base.GetCity()),
				Account:   easygo.NewString(base.GetAccount()),
				Photo:     easygo.NewString(photo),
				Area:      easygo.NewString(base.GetArea()),
				Types:     easygo.NewInt32(base.GetTypes()),
			}
			pInfo = append(pInfo, msg)
		}
	}

	allList := GetPlayerForChat(pid, int(t), 10) //获取所有客服号的人物id
	for _, pid := range allList {
		player := GetRedisPlayerBase(pid)
		dis := RandInt(100, 10000)
		var photo string
		photos := player.GetPhoto()
		if len(photos) != 0 {
			photo = photos[0]
		}
		msg := &client_hall.LocationPlayerInfo{
			PlayerId:  easygo.NewInt64(pid),
			Distance:  easygo.NewFloat64(dis),
			Name:      easygo.NewString(player.GetNickName()),
			Sex:       easygo.NewInt32(player.GetSex()),
			HeadIcon:  easygo.NewString(player.GetHeadIcon()),
			Signature: easygo.NewString(player.GetSignature()),
			IsFriend:  easygo.NewBool(false),
			Province:  easygo.NewString(province),
			Account:   easygo.NewString(player.GetAccount()),
			City:      easygo.NewString(city),
			Photo:     easygo.NewString(photo),
			Area:      easygo.NewString(area),
			Types:     easygo.NewInt32(player.GetTypes()),
		}
		pInfo = append(pInfo, msg)
	}

	sort.Slice(pInfo, func(i, j int) bool {
		return pInfo[i].GetDistance() < pInfo[j].GetDistance() // 升序
	})

	allInfo := &client_hall.AllLocationPlayerInfo{
		PlayerInfo: pInfo,
	}
	return allInfo
}

// 新版本获取附近的人
//func GetLocationInfo1(pid int64, msg *client_hall.LocationInfo) {
func GetLocationInfo1() {
	// 20公里 20000
	// 60人真实用户,20客服,坐标需要排序.先设置值.,我再弄机器人.最后才查出来.
	// 20个机器人.
	/**
	1 判断只看男生,还是只看女生,还是查看全部
	2.判断是否开启了附近的人,判断x,y,是不为0
	3.判断是否够100人
	4.不够100人,从客服中抽取20人,对这20人进行
	//=========================
	// 1判断够不够100人
	2.弄客服的坐标进数据库.
	3通过坐标排序查出来
	4,生成机器人的距离.
	5.排序
	6.结束返回.
	//===========够100人=======
	直接查询.
	查出来的列表进行排序.
	*/
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	m := []bson.M{
		{
			"$geoNear": bson.M{
				"near":               []float64{116.416637, 39.922705},
				"spherical":          true,
				"distanceMultiplier": 6378137,
				"maxDistance":        100.0 / 6378137.0, // 20公里
				"distanceField":      "Distance",
			},
		},
	}
	var players []*share_message.PlayerBase
	marshal, _ := json.Marshal(m)
	logs.Info("m==============>", string(marshal))
	query := col.Pipe(m)
	err := query.All(&players)
	easygo.PanicError(err)

	bytes, _ := json.Marshal(players)
	logs.Info("players---------->", len(players), string(bytes))
	dis := make([]float64, 0)
	pid := make([]int64, 0)
	for _, v := range players {
		dis1 := GetDistance(v.GetY(), 39.922705, v.GetX(), 116.416637)
		dis = append(dis, dis1)
		pid = append(pid, v.GetPlayerId())
	}
	logs.Info("dis----------->", dis)
	logs.Info("pid----------->", pid)
}

func GetAllNearByInfo(pid PLAYER_ID) []*client_hall.NearByMessage {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_LOG)
	defer closeFun()
	var obj *share_message.NearByGreetInfo
	var lst []*client_hall.NearByMessage
	err := col.Find(bson.M{"_id": pid}).One(&obj)
	if err != nil {
		if err.Error() != mgo.ErrNotFound.Error() {
			easygo.PanicError(err)
		}
		return lst
	}

	for _, m := range obj.GetNearByInfo() {
		id := m.GetPlayerId()
		base := GetRedisPlayerBase(id)
		var photo string
		photolst := base.GetPhoto()
		if len(photolst) != 0 {
			photo = photolst[0]
		}
		msg := &client_hall.NearByMessage{
			PlayerId:  easygo.NewInt64(id),
			Content:   easygo.NewString(m.GetContent()),
			NickName:  easygo.NewString(base.GetNickName()),
			HeadIcon:  easygo.NewString(base.GetHeadIcon()),
			Sex:       easygo.NewInt32(base.GetSex()),
			Provice:   easygo.NewString(base.GetProvice()),
			City:      easygo.NewString(base.GetCity()),
			Photo:     easygo.NewString(photo),
			Account:   easygo.NewString(base.GetAccount()),
			Time:      easygo.NewInt64(m.GetTime()),
			IsAdd:     easygo.NewBool(m.GetIsAdd()),
			Signature: easygo.NewString(base.GetSignature()),
			Types:     easygo.NewInt32(base.GetTypes()),
		}
		lst = append(lst, msg)
	}
	return lst
}

func GetAllNewNearByInfo(pid PLAYER_ID) []*client_hall.NearByMessage {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_LOG)
	defer closeFun()
	var obj *share_message.NearByGreetInfo
	err := col.Find(bson.M{"_id": pid}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	var lst []*client_hall.NearByMessage
	if err == mgo.ErrNotFound {
		obj = &share_message.NearByGreetInfo{
			PlayerId:   easygo.NewInt64(pid),
			NearByInfo: []*share_message.NearByInfo{},
		}
		col.Insert(&obj)
		return lst
	}

	alllst := obj.GetNearByInfo()
	t := time.Now().Unix()
	var newlst []*share_message.NearByInfo
	for _, m := range alllst {
		if t-m.GetTime() > 3*86400 {
			continue
		}
		newlst = append(newlst, m)
	}

	for _, m := range newlst {
		if m.GetIsRead() {
			continue
		}
		m.IsRead = easygo.NewBool(true)
		id := m.GetPlayerId()
		base := GetRedisPlayerBase(id)
		var photo string
		photolst := base.GetPhoto()
		if len(photolst) != 0 {
			photo = photolst[0]
		}
		msg := &client_hall.NearByMessage{
			PlayerId:  easygo.NewInt64(id),
			Content:   easygo.NewString(m.GetContent()),
			NickName:  easygo.NewString(base.GetNickName()),
			HeadIcon:  easygo.NewString(base.GetHeadIcon()),
			Sex:       easygo.NewInt32(base.GetSex()),
			Provice:   easygo.NewString(base.GetProvice()),
			City:      easygo.NewString(base.GetCity()),
			Photo:     easygo.NewString(photo),
			Account:   easygo.NewString(base.GetAccount()),
			Time:      easygo.NewInt64(m.GetTime()),
			IsAdd:     easygo.NewBool(m.GetIsAdd()),
			Signature: easygo.NewString(base.GetSignature()),
			Types:     easygo.NewInt32(base.GetTypes()),
		}
		lst = append(lst, msg)
	}

	_, err1 := col.Upsert(bson.M{"_id": pid}, bson.M{"$set": bson.M{"NearByInfo": newlst}})
	if err1 != nil {
		easygo.PanicError(err1)
	}

	return lst
}

func AddNearByInfo(pid, id int64, content string) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_LOG)
	defer closeFun()

	var obj *share_message.NearByGreetInfo
	err := col.Find(bson.M{"_id": pid}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}

	if err == mgo.ErrNotFound {
		obj = &share_message.NearByGreetInfo{
			PlayerId:   easygo.NewInt64(pid),
			NearByInfo: []*share_message.NearByInfo{},
		}
		col.Insert(&obj)
	}

	lst := obj.GetNearByInfo()
	msg := &share_message.NearByInfo{
		PlayerId: easygo.NewInt64(id),
		Content:  easygo.NewString(content),
		IsAdd:    easygo.NewBool(false),
		Time:     easygo.NewInt64(time.Now().Unix()),
	}

	if len(lst) >= 50 {
		lst = lst[1:]
		lst = append(lst, msg)
		_, err1 := col.Upsert(bson.M{"_id": pid}, bson.M{"$set": bson.M{"NearByInfo": lst}})
		if err1 != nil {
			easygo.PanicError(err1)
		}
	} else {
		err1 := col.Update(bson.M{"_id": pid}, bson.M{"$push": bson.M{"NearByInfo": msg}})
		if err1 != nil {
			easygo.PanicError(err1)
		}
	}

}

func AgreeNearByInfo(pid, id int64) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_LOG)
	defer closeFun()

	var obj *share_message.NearByGreetInfo
	err := col.Find(bson.M{"_id": pid}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	alllst := obj.GetNearByInfo()
	for _, m := range alllst {
		if m.GetPlayerId() == id {
			m.IsAdd = easygo.NewBool(true)
		}
	}
	_, err1 := col.Upsert(bson.M{"_id": pid}, bson.M{"$set": bson.M{"NearByInfo": alllst}})
	if err1 != nil {
		easygo.PanicError(err1)
	}
}

func DelNearByInfo(pid PLAYER_ID) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_LOG)
	defer closeFun()

	var obj *share_message.NearByGreetInfo
	err := col.Find(bson.M{"_id": pid}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return
	}
	_, err1 := col.Upsert(bson.M{"_id": pid}, bson.M{"$set": bson.M{"NearByInfo": []*share_message.NearByInfo{}}})
	if err1 != nil {
		easygo.PanicError(err1)
	}
}

//获取陪聊ID数组  sex性别，nmb取得条数 数据库随机查询
func GetPlayerForChat(pid int64, sex int, num int, labels ...[]int32) []int64 {
	lst := append(labels, []int32{})[0]
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()

	queryBson := bson.M{"Types": 2, "_id": bson.M{"$ne": pid}}
	if sex > 0 {
		queryBson["Sex"] = sex
	}
	if len(lst) > 0 {
		queryBson["Label"] = bson.M{"$elemMatch": bson.M{"$in": lst}}
	}
	t := GetMillSecond()
	queryBson["LastOnLineTime"] = bson.M{"$gte": t - 3*86400*1000}
	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
		{"$project": bson.M{"_id": 1}},
	}

	query := col.Pipe(m)
	var list []*share_message.PlayerBase
	err := query.All(&list)
	easygo.PanicError(err)
	ids := make([]int64, 0)
	ids2 := make([]int64, 0)
	for _, info := range list {
		ids = append(ids, info.GetPlayerId())
		ids2 = append(ids2, info.GetPlayerId())
	}

	if len(list) != num { //如果标签匹配的客服不够 就随机取客服
		need := num - len(list)
		ids2 = append(ids2, pid)
		delete(queryBson, "Label")
		queryBson["_id"] = bson.M{"$nin": ids2}
		m := []bson.M{
			{"$match": queryBson},
			{"$sample": bson.M{"size": need}},
			{"$project": bson.M{"_id": 1}},
		}

		query := col.Pipe(m)
		var list []*share_message.PlayerBase
		err := query.All(&list)
		easygo.PanicError(err)
		for _, info := range list {
			ids = append(ids, info.GetPlayerId())
		}
	}
	return ids
}

func GetMoneyOrderInfo(pid PLAYER_ID, t []int32, count int, year, month int) []*GoldLog {
	var isAll bool
	if len(t) == 0 {
		isAll = true
	}
	var year2, month2 int
	if month == 12 {
		year2 = year + 1
		month2 = 1
	} else {
		year2 = year
		month2 = month + 1
	}
	t1 := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).UnixNano() / 1e6
	t2 := time.Date(year2, time.Month(month2), 1, 0, 0, 0, 0, time.Local).UnixNano() / 1e6

	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_GOLDCHANGELOG)
	defer closeFun()
	var lst []*GoldLog
	if isAll { //是否是全部
		err := col.Find(bson.M{"PlayerId": pid, "CreateTime": bson.M{"$gte": t1, "$lte": t2}}).Limit(count).Sort("-CreateTime").All(&lst)
		if err != nil && err != mgo.ErrNotFound {
			easygo.PanicError(err)
		}
	} else {
		err := col.Find(bson.M{"PlayerId": pid, "CreateTime": bson.M{"$gte": t1, "$lte": t2}, "SourceType": bson.M{"$in": t}}).Limit(count).Sort("-CreateTime").All(&lst)
		if err != nil && err != mgo.ErrNotFound {
			easygo.PanicError(err)
		}
	}
	return lst
}

func GetCashOrderInfo(pid PLAYER_ID, page, num int) []*GoldLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_GOLDCHANGELOG)
	defer closeFun()
	var lst []*GoldLog
	err := col.Find(bson.M{"PlayerId": pid}).Skip((page - 1) * num).Limit(num).Sort("-CreateTime").All(&lst)
	if err == mgo.ErrNotFound {
		return lst
	}
	easygo.PanicError(err)
	return lst
}

func GetRedPacketOrderInfo(pid int64, gt int32, year, month int) []*GoldLog {
	var year2, month2 int
	if month == 12 {
		year2 = year + 1
		month2 = 1
	} else {
		year2 = year
		month2 = month + 1
	}
	t1 := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).UnixNano() / 1e6
	t2 := time.Date(year2, time.Month(month2), 1, 0, 0, 0, 0, time.Local).UnixNano() / 1e6

	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_GOLDCHANGELOG)
	defer closeFun()
	var allList []*GoldLog
	err1 := col.Find(bson.M{"PlayerId": pid, "CreateTime": bson.M{"$gte": t1, "$lte": t2}, "SourceType": gt}).Sort("-CreateTime").All(&allList)
	if err1 != mgo.ErrNotFound && err1 != nil {
		easygo.PanicError(err1)
	}
	return allList
}

func GetRedPacketForPageInfo(pid int64, gt int32, num, year, month int) []*GoldLog {
	var year2, month2 int
	if month == 12 {
		year2 = year + 1
		month2 = 1
	} else {
		year2 = year
		month2 = month + 1
	}
	t1 := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).UnixNano() / 1e6
	t2 := time.Date(year2, time.Month(month2), 1, 0, 0, 0, 0, time.Local).UnixNano() / 1e6
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_GOLDCHANGELOG)
	defer closeFun()
	var pageList []*GoldLog
	err := col.Find(bson.M{"PlayerId": pid, "CreateTime": bson.M{"$gte": t1, "$lte": t2}, "SourceType": gt}).Limit(num).Sort("-CreateTime").All(&pageList)
	if err != mgo.ErrNotFound && err != nil {
		easygo.PanicError(err)
	}
	return pageList
}

//增加投诉建议入库
func AddPlayerComplaint(msg *share_message.PlayerComplaint) error {
	logId := NextId(TABLE_PLAYER_COMPLAINT)
	msg.Id = easygo.NewInt64(logId)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_COMPLAINT)
	defer closeFun()
	err := col.Insert(msg)
	easygo.PanicError(err)
	return err
}

func CheckPeopleIdIsValid(id string) bool {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	var obj *share_message.PlayerBase
	err := col.Find(bson.M{"PeopleId": id}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return true
	}
	return false
}

func GetUnGetMoneyInfo(id int64, t int32) bool {
	var b bool
	var redList []*share_message.RedPacket
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_RED_PACKET)
	defer closeFun()
	err := col.Find(bson.M{"TargetId": id, "State": PACKET_MONEY_OPEN}).All(&redList)
	if err != mgo.ErrNotFound && err != nil {
		easygo.PanicError(err)
	}
	if len(redList) != 0 {
		b = true
	}
	if t == CHAT_TYPE_PRIVATE {
		var traList []*share_message.TransferMoney
		col1, closeFun1 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TRANSFER_MONEY)
		defer closeFun1()
		err := col1.Find(bson.M{"TargetId": id, "State": TRANSFER_MONEY_OPEN}).All(&traList)
		if err != mgo.ErrNotFound && err != nil {
			easygo.PanicError(err)
		}
		if len(traList) != 0 {
			b = true
		}
	}
	return b
}

//查询当前可用的系统参数设置
func QuerySysParameterList() []*share_message.SysParameter {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SYS_PARAMETER)
	defer closeFun()

	queryBson := bson.M{}
	var list []*share_message.SysParameter
	query := col.Find(queryBson)
	errc := query.Sort("-_id").All(&list)
	easygo.PanicError(errc)
	return list
}

//写埋点注册登录日志 loginType=0-重连,1-登录
func AddStatisticsInfo(t int32, pid int64, loginType int32, c ...int32) {
	var tableName string
	tableName = TABLE_LOGIN_REGISTER_LOG //埋点登录注册日志
	msg := &share_message.LoginRegisterInfo{
		Type:     easygo.NewInt32(t),
		Time:     easygo.NewInt64(GetMillSecond()),
		PlayerId: easygo.NewInt64(pid),
	}

	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, tableName)
	defer closeFun()
	err := col.Insert(msg)
	easygo.PanicError(err)

	easygo.Spawn(func() {
		if len(c) == 0 && loginType == 1 {
			MakeRegisterLoginReport(msg) //生成埋点注册登录报表
		}

		var operationChannelReportType int32 = 6
		switch t {
		case 7, 8, 9:
			operationChannelReportType = 1
		}
		pMgr := PlayerBaseMgr.GetRedisPlayerBaseObj(pid)
		channleNo := pMgr.GetChannel()
		MakeOperationChannelReport(operationChannelReportType, pid, channleNo, nil, nil) //生成运营渠道数据汇总报表 已优化到Redis
	})
}

//id查询系统参数设置 //id: 转账功能参数 limit_parameter，头像参数 avatar_parameter
func QuerySysParameterById(id string) *share_message.SysParameter {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SYS_PARAMETER)
	defer closeFun()
	rp := &share_message.SysParameter{}

	errc := col.Find(bson.M{"_id": id}).One(rp)
	if errc != nil && errc != mgo.ErrNotFound {
		panic(errc)
	}
	if errc == mgo.ErrNotFound {
		return nil
	}
	return rp
}

//删除系统参数设置
func DeleteSysParameter() {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SYS_PARAMETER)
	defer closeFun()

	queryBson := bson.M{}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//渠道号查询渠道详情
func QueryOperationByNo(no string) *share_message.OperationChannel {
	data := &share_message.OperationChannel{}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_OPERATION_CHANNEL_USE)
	defer closeFun()
	errc := col.Find(bson.M{"ChannelNo": no}).One(data)

	if errc != nil && errc != mgo.ErrNotFound {
		panic(errc)
	}
	if errc == mgo.ErrNotFound {
		return nil
	}

	return data
}

//渠道号查询渠道列表
func QueryOperationChannleList() []*share_message.OperationChannel {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_OPERATION_CHANNEL_USE)
	defer closeFun()
	var list []*share_message.OperationChannel
	query := col.Find(bson.M{"Status": 1})
	errc := query.All(&list)
	easygo.PanicError(errc)

	return list
}

//判断IP是否访问过  true 访问过 false 未访问过
func IsNewIp(ip string) bool {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_IPLIBRARY)
	defer closeFun()
	ipl := &share_message.IpLibrary{}
	err := col.Find(bson.M{"_id": ip}).One(&ipl)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		ipMsg := &share_message.IpLibrary{Ip: easygo.NewString(ip)}
		_, erru := col.Upsert(bson.M{"_id": ip}, bson.M{"$set": ipMsg})
		easygo.PanicError(erru)

		return false
	}

	return true
}

//获取所有标签
func GetInterestTagAllList() []*share_message.InterestTag {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_INTERESTTAG)
	defer closeFun()

	queryBson := bson.M{"Status": 0}
	var list []*share_message.InterestTag
	query := col.Find(queryBson)

	errc := query.Sort("Sort").All(&list)
	easygo.PanicError(errc)

	return list
}

func GetInterestTag(id int32) *share_message.InterestTag {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_INTERESTTAG)
	defer closeFun()

	var obj *share_message.InterestTag
	err := col.Find(bson.M{"_id": id}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return obj
}

//查询兴趣分类列表查询
func GetInterestTypeList() []*share_message.InterestType {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_INTERESTTYPE)
	defer closeFun()

	var list []*share_message.InterestType
	err := col.Find(bson.M{"Status": 0}).All(&list)
	easygo.PanicError(err)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return list
}

//默认推荐玩家列表
var DefaultAccList = []string{
	"lm708195f6",
	"lm7081967e",
	"lm70819682",
	"lm70819686",
	"lm70819687",
	"lm7081968a",
	"lm7081968d",
	"lm7081a8de",
	"lm7081a8ee",
	"lm7081961e",
	"lm70819675",
	"lm70819473",
	"lm70819449",
	"lm7081966b",
	"lm70819659"}

func GetRecommendPlayerInfo(pid int64) []*client_server.RecommendPlayerInfo {
	rand.Seed(time.Now().UnixNano())
	//推荐的玩家
	player := GetRedisPlayerBase(pid)
	if player == nil {
		panic("没有玩家对象存在")
	}
	labels := player.GetLabelList()
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	var playerList, accountList []*share_message.PlayerBase
	t := GetMillSecond()
	min := t - 3*86400*1000

	queryBson := bson.M{"_id": bson.M{"$ne": pid}, "Label": bson.M{"$elemMatch": bson.M{"$in": labels}},
		"Status": ACCOUNT_NORMAL, "IsRecommendOver": false, "LastOnLineTime": bson.M{"$gte": min}, "Types": 1}

	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": 20}},
		{"$sort": bson.M{"LastOnLineTime": -1}},
	}

	query := col.Pipe(m)
	err1 := query.All(&playerList)
	easygo.PanicError(err1)

	err := col.Find(bson.M{"Account": bson.M{"$in": DefaultAccList}}).All(&accountList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	//真实玩家
	for i := len(playerList) - 1; i > 0; i-- { //随机排序
		num := rand.Intn(i + 1)
		playerList[i], playerList[num] = playerList[num], playerList[i]
	}

	var realList, serverList, robotList []*client_server.RecommendPlayerInfo
	var real, server, robot int
	leng := len(playerList)
	lenac := len(accountList)
	if leng >= 3 { //如果真人玩家超过三分之一
		if leng > 5 {

			if lenac >= 2 {
				real = 4
				server = 3
			} else {
				real = 5
				server = 4
			}
		} else {
			if lenac >= 2 {
				real = leng
				server = 7 - leng
			} else {
				real = leng
				server = 9 - leng
			}

		}
		robot = 0
	} else {
		//num := 3 - leng    //补充多少客服和机器人
		//num1 := num / 2    //补充机器人数量
		//num2 := num - num1 //补充客服数量
		//real = leng
		//server = 3 + num2
		//robot = 3 + num1
		real = leng
		server = 7 - real
		robot = 0
	}
	//真实玩家列表
	for _, info := range playerList[:real] {
		msg := &client_server.RecommendPlayerInfo{
			PlayerId: easygo.NewInt64(info.GetPlayerId()),
			Name:     easygo.NewString(info.GetNickName()),
			Sex:      easygo.NewInt32(info.GetSex()),
			HeadIcon: easygo.NewString(info.GetHeadIcon()),
			Type:     easygo.NewInt32(1),
		}
		realList = append(realList, msg)
	}
	//真实指定玩家列表
	if len(accountList) >= 2 {
		ls := easygo.RandGetNItemFromSlice(accountList, 2)
		for _, p := range ls {
			if info, ok := p.(*share_message.PlayerBase); ok {
				msg := &client_server.RecommendPlayerInfo{
					PlayerId: easygo.NewInt64(info.GetPlayerId()),
					Name:     easygo.NewString(info.GetNickName()),
					Sex:      easygo.NewInt32(info.GetSex()),
					HeadIcon: easygo.NewString(info.GetHeadIcon()),
					Type:     easygo.NewInt32(1),
				}
				realList = append(realList, msg)
				//logs.Info("运营玩家:", msg)
			}
		}
	}
	serverIds := GetPlayerForChat(pid, 0, server, labels)
	for _, pid := range serverIds {
		player := GetRedisPlayerBase(pid)
		msg := &client_server.RecommendPlayerInfo{
			PlayerId: easygo.NewInt64(pid),
			Name:     easygo.NewString(player.GetNickName()),
			Sex:      easygo.NewInt32(player.GetSex()),
			HeadIcon: easygo.NewString(player.GetHeadIcon()),
			Type:     easygo.NewInt32(2),
		}
		serverList = append(serverList, msg)
	}

	if robot > 0 {
		minId := Min_Robot_PlayerId
		manHeadList := GetManyRobotHeadIcon(robot, 1)  //男随机头像列表
		girlHeadList := GetManyRobotHeadIcon(robot, 2) //女头像随机列表
		manNameList := GetManyRobotName(1, robot)
		girlNameList := GetManyRobotName(2, robot)
		for i := minId; i < minId+robot; i++ {
			sex := RandInt(1, 3)
			var name string
			if sex == 1 {
				name = manNameList[i-minId]
			} else {
				name = girlNameList[i-minId]
			}
			msg := &client_server.RecommendPlayerInfo{
				PlayerId: easygo.NewInt64(i),
				Name:     easygo.NewString(name),
				Sex:      easygo.NewInt32(sex),
				Type:     easygo.NewInt32(3),
			}
			var mark string
			var icon int
			if sex == 1 {
				mark = "mavatar"
				icon = manHeadList[i-minId]
			} else {
				mark = "wavatar"
				icon = girlHeadList[i-minId]
			}
			head := fmt.Sprintf("https://im-resource-1253887233.file.myqcloud.com/prod/%s/%d.png", mark, icon)
			msg.HeadIcon = easygo.NewString(head)
			robotList = append(robotList, msg)
		}
	}

	playerInfo := make([]*client_server.RecommendPlayerInfo, 0)
	playerInfo = append(playerInfo, realList...)
	playerInfo = append(playerInfo, serverList...)
	playerInfo = append(playerInfo, robotList...)
	return playerInfo

}

// 优化前
/*func GetRecommendPlayerInfo(pid int64) []*client_server.RecommendPlayerInfo {
	rand.Seed(time.Now().UnixNano())
	//推荐的玩家
	player := GetRedisPlayerBase(pid)
	if player == nil {
		panic("没有玩家对象存在")
	}
	labels := player.GetLabelList()
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	var playerList, accountList []*share_message.PlayerBase
	t := GetMillSecond()
	min := t - 3*86400*1000
	err := col.Find(bson.M{"_id": bson.M{"$ne": pid}, "Label": bson.M{"$elemMatch": bson.M{"$in": labels}},
		"Status": ACCOUNT_NORMAL, "IsRecommendOver": false, "LastOnLineTime": bson.M{"$gte": min}, "Types": 1}).Sort("-LastOnLineTime").Limit(20).All(&playerList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	err = col.Find(bson.M{"Account": bson.M{"$in": DefaultAccList}}).All(&accountList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	//logs.Info("推荐真实用户:", len(playerList), playerList)
	//logs.Info("指定用户:", len(accountList), accountList)
	//真实玩家
	for i := len(playerList) - 1; i > 0; i-- { //随机排序
		num := rand.Intn(i + 1)
		playerList[i], playerList[num] = playerList[num], playerList[i]
	}
	//指定的玩家
	//for i := len(accountList) - 1; i > 0; i-- { //随机排序
	//	num := rand.Intn(i + 1)
	//	accountList[i], accountList[num] = accountList[num], accountList[i]
	//}
	var realList, serverList, robotList []*client_server.RecommendPlayerInfo
	var real, server, robot int
	leng := len(playerList)
	lenac := len(accountList)
	if leng >= 3 { //如果真人玩家超过三分之一
		if leng > 5 {
			//real = 5
			//server = 4
			if lenac >= 2 {
				real = 4
				server = 3
			} else {
				real = 5
				server = 4
			}
		} else {
			if lenac >= 2 {
				real = leng
				server = 7 - leng
			} else {
				real = leng
				server = 9 - leng
			}

		}
		robot = 0
	} else {
		//num := 3 - leng    //补充多少客服和机器人
		//num1 := num / 2    //补充机器人数量
		//num2 := num - num1 //补充客服数量
		//real = leng
		//server = 3 + num2
		//robot = 3 + num1
		real = leng
		server = 7 - real
		robot = 0
	}
	//真实玩家列表
	for _, info := range playerList[:real] {
		msg := &client_server.RecommendPlayerInfo{
			PlayerId: easygo.NewInt64(info.GetPlayerId()),
			Name:     easygo.NewString(info.GetNickName()),
			Sex:      easygo.NewInt32(info.GetSex()),
			HeadIcon: easygo.NewString(info.GetHeadIcon()),
			Type:     easygo.NewInt32(1),
		}
		realList = append(realList, msg)
	}
	//真实指定玩家列表
	if len(accountList) >= 2 {
		ls := easygo.RandGetNItemFromSlice(accountList, 2)
		for _, p := range ls {
			if info, ok := p.(*share_message.PlayerBase); ok {
				msg := &client_server.RecommendPlayerInfo{
					PlayerId: easygo.NewInt64(info.GetPlayerId()),
					Name:     easygo.NewString(info.GetNickName()),
					Sex:      easygo.NewInt32(info.GetSex()),
					HeadIcon: easygo.NewString(info.GetHeadIcon()),
					Type:     easygo.NewInt32(1),
				}
				realList = append(realList, msg)
				//logs.Info("运营玩家:", msg)
			}
		}
	}
	serverIds := GetPlayerForChat(pid, 0, server, labels)
	for _, pid := range serverIds {
		player := GetRedisPlayerBase(pid)
		msg := &client_server.RecommendPlayerInfo{
			PlayerId: easygo.NewInt64(pid),
			Name:     easygo.NewString(player.GetNickName()),
			Sex:      easygo.NewInt32(player.GetSex()),
			HeadIcon: easygo.NewString(player.GetHeadIcon()),
			Type:     easygo.NewInt32(2),
		}
		serverList = append(serverList, msg)
	}

	if robot > 0 {
		minId := Min_Robot_PlayerId
		manHeadList := GetManyRobotHeadIcon(robot, 1)  //男随机头像列表
		girlHeadList := GetManyRobotHeadIcon(robot, 2) //女头像随机列表
		manNameList := GetManyRobotName(1, robot)
		girlNameList := GetManyRobotName(2, robot)
		for i := minId; i < minId+robot; i++ {
			sex := RandInt(1, 3)
			var name string
			if sex == 1 {
				name = manNameList[i-minId]
			} else {
				name = girlNameList[i-minId]
			}
			msg := &client_server.RecommendPlayerInfo{
				PlayerId: easygo.NewInt64(i),
				Name:     easygo.NewString(name),
				Sex:      easygo.NewInt32(sex),
				Type:     easygo.NewInt32(3),
			}
			var mark string
			var icon int
			if sex == 1 {
				mark = "mavatar"
				icon = manHeadList[i-minId]
			} else {
				mark = "wavatar"
				icon = girlHeadList[i-minId]
			}
			head := fmt.Sprintf("https://im-resource-1253887233.file.myqcloud.com/prod/%s/%d.png", mark, icon)
			msg.HeadIcon = easygo.NewString(head)
			robotList = append(robotList, msg)
		}
	}

	playerInfo := []*client_server.RecommendPlayerInfo{}
	playerInfo = append(playerInfo, realList...)
	playerInfo = append(playerInfo, serverList...)
	playerInfo = append(playerInfo, robotList...)
	return playerInfo

}
*/
func GetRecommendTeamInfo() []*client_server.RecommendTeamInfo {
	rand.Seed(time.Now().Unix())
	var teamList []*share_message.TeamData
	col1, closeFun1 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAM_DATA)
	defer closeFun1()
	err1 := col1.Find(bson.M{"IsRecommend": true, "Status": 0}).All(&teamList)
	easygo.PanicError(err1)
	for i := len(teamList) - 1; i > 0; i-- { //随机排序
		num := rand.Intn(i + 1)
		teamList[i], teamList[num] = teamList[num], teamList[i]
	}
	//var teamIds []int64
	//if len(teamList) > 18 {
	//	for _, info := range teamList[:18] {
	//		teamIds = append(teamIds, info.GetId())
	//	}
	//} else {
	//	for _, info := range teamList {
	//		teamIds = append(teamIds, info.GetId())
	//	}
	//}
	teamInfo := []*client_server.RecommendTeamInfo{}
	var newTeamList []*share_message.TeamData
	if len(teamList) > 6 {
		newTeamList = teamList[:6]
	} else {
		newTeamList = teamList
	}

	for _, info := range newTeamList {
		if info.GetMaxMember() == int32(len(info.GetMemberList())) {
			continue
		}
		msg := &client_server.RecommendTeamInfo{
			TeamId:  easygo.NewInt64(info.GetId()),
			Name:    easygo.NewString(info.GetName()),
			OwnerId: easygo.NewInt64(info.GetOwner()),
		}
		var memberInfo []*client_server.RecommendPlayerInfo
		members := info.GetMemberList()
		var newMembers []int64
		if len(members) >= 9 {
			newMembers = members[:9]
		} else {
			newMembers = members
		}
		for _, pid := range newMembers {
			player := GetRedisPlayerBase(pid)
			if player == nil {
				continue
			}
			m := &client_server.RecommendPlayerInfo{
				HeadIcon: easygo.NewString(player.GetHeadIcon()),
				Sex:      easygo.NewInt32(player.GetSex()),
				PlayerId: easygo.NewInt64(pid),
				Name:     easygo.NewString(player.GetNickName()),
			}
			memberInfo = append(memberInfo, m)
		}
		msg.MemberInfo = memberInfo
		teamInfo = append(teamInfo, msg)
	}
	return teamInfo
}

func GetRecommendInfo(pid int64, t int32) *client_server.RecommendInfo {
	msg := &client_server.RecommendInfo{}
	if t == 0 {
		msg.PlayerInfo = GetRecommendPlayerInfo(pid)
		msg.TeamInfo = GetRecommendTeamInfo()
	} else if t == 1 {
		msg.PlayerInfo = GetRecommendPlayerInfo(pid)
	} else {
		msg.TeamInfo = GetRecommendTeamInfo()
	}
	return msg
}

func GetAllRedPacketInfo(pid int64, t int32, year, month int) map[int64]*share_message.RedPacket {
	var year2, month2 int
	if month == 12 {
		year2 = year + 1
		month2 = 1
	} else {
		year2 = year
		month2 = month + 1
	}
	t1 := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).Unix()
	t2 := time.Date(year2, time.Month(month2), 1, 0, 0, 0, 0, time.Local).Unix()
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_RED_PACKET)
	defer closeFun()
	col1, closeFun1 := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_RED_PACKET_LOG)
	defer closeFun1()

	var lst []*share_message.RedPacket
	if t == 1 { //收红包
		//查找收包人是自己的所有红包
		err := col.Find(bson.M{"TargetId": pid, "CreateTime": bson.M{"$gte": t1, "$lte": t2}}).All(&lst)
		easygo.PanicError(err)

		//查找自己抢过红包的群红包
		var redLogds []*share_message.RedPacketLog
		var teamRedIds []int64
		err1 := col1.Find(bson.M{"PlayerId": pid, "CreateTime": bson.M{"$gte": t1 * 1000, "$lte": t2 * 1000}}).All(&redLogds)
		easygo.PanicError(err1)
		for _, log := range redLogds {
			teamRedIds = append(teamRedIds, log.GetRedPacketId())
		}
		if len(teamRedIds) > 0 {
			var newlst []*share_message.RedPacket
			err2 := col.Find(bson.M{"_id": bson.M{"$in": teamRedIds}}).All(&newlst)
			easygo.PanicError(err2)
			lst = append(lst, newlst...)
		}
	} else { //发送包
		err := col.Find(bson.M{"Sender": pid, "CreateTime": bson.M{"$gte": t1, "$lte": t2}}).All(&lst)
		easygo.PanicError(err)

	}
	var redIds []int64 //所有与自己有关的红包id
	for _, info := range lst {
		redIds = append(redIds, info.GetId())
	}
	var logList []interface{}
	logInfo := make(map[int64][]*share_message.RedPacketLog)

	queryBson := bson.M{"RedPacketId": bson.M{"$in": redIds}}
	err2 := col1.Find(queryBson).All(&logList)
	easygo.PanicError(err2)
	for _, m := range logList {
		log := &share_message.RedPacketLog{}
		info := m.(bson.M)
		j, _ := json.Marshal(info)
		json.Unmarshal(j, &log)
		logId := log.GetRedPacketId()
		logInfo[logId] = append(logInfo[logId], log)
	}

	for _, info := range lst {
		for id, loglst := range logInfo {
			if info.GetId() == id {
				info.Logs = loglst
				break
			}
		}
	}

	redPacketInfo := make(map[int64]*share_message.RedPacket)
	for _, m := range lst {
		redPacketInfo[m.GetId()] = m
	}
	return redPacketInfo
}

//批量通过id获取用户信息
func GetPlayerListByIds(ids []int64) []*share_message.PlayerBase {
	var list []*share_message.PlayerBase
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()

	// 关键词查询账号
	queryBson := bson.M{"_id": bson.M{"$in": ids}}
	query := col.Find(queryBson)
	errc := query.All(&list)
	easygo.PanicError(errc)
	return list
}

func GetJGIds(lst []int64) []string {
	var list []*share_message.PlayerBase
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()

	// 关键词查询账号,只有正常的用户才极光推送
	queryBson := bson.M{"_id": bson.M{"$in": lst}, "Status": ACCOUNT_NORMAL}
	query := col.Find(queryBson)
	errc := query.All(&list)
	easygo.PanicError(errc)
	var playerList []string
	for _, m := range list {
		if m.GetToken() != "" {
			playerList = append(playerList, m.GetToken())
		}
	}
	return playerList
}

//初始化自增id的初值
func InitIdGenerator(key string, val int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ID_GENERATOR)
	defer closeFun()
	identity := &Identity{}
	err := col.Find(bson.M{"_id": key}).One(&identity)
	if err == mgo.ErrNotFound {
		identity = &Identity{
			Key:   easygo.NewString(key),
			Value: easygo.NewUint64(val),
		}
		err = col.Insert(identity)
		easygo.PanicError(err)
	}
}

//保存终端设备号
func SavePosDeviceCode(reqMsg *share_message.PosDeviceCode) bool {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_POS_DEVICECODE)
	defer closeFun()

	n, err := col.Find(bson.M{"DeviceCode": reqMsg.GetDeviceCode()}).Count()
	if n == 0 {
		_, err = col.Upsert(bson.M{"_id": reqMsg.GetCreateTime()}, bson.M{"$set": reqMsg})
		easygo.PanicError(err)
		if reqMsg.Channle != nil && reqMsg.GetChannle() != "" {
			channel := QueryOperationByNo(reqMsg.GetChannle()) //查询渠道信息
			if channel != nil {
				SetRedisOperationChannelReportFildVal(easygo.Get0ClockTimestamp(reqMsg.GetCreateTime()), 1, reqMsg.GetChannle(), "ActDevCount") //更新渠道激活设备数量
			}
		}
		return true
	}
	easygo.PanicError(err)
	return false
}

//保存终端设备Idfa码
func SavePosDeviceIdfa(reqMsg *share_message.PosDeviceIdfa) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_POS_DEVICEIDFA)
	defer closeFun()

	n, err := col.Find(bson.M{"_id": reqMsg.GetDeviceIdfa()}).Count()
	if n == 0 {
		_, err = col.Upsert(bson.M{"_id": reqMsg.GetDeviceIdfa()}, bson.M{"$set": reqMsg})
	}
	easygo.PanicError(err)
}

//检查终端设备Idfa码是否存在
func CheckPosDeviceIdfa(idfa string) bool {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_POS_DEVICEIDFA)
	defer closeFun()
	n, err := col.Find(bson.M{"_id": idfa}).Count()
	easygo.PanicError(err)
	result := false
	if n == 0 {
		result = true
	}
	return result
}

//保存终端设备Adv Idfa码
func SavePosDeviceAdvIdfa(reqMsg *share_message.KsPosAdvIdfa, isSave bool) bool {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_POS_ADV_DEVICEIDFA)
	defer closeFun()

	if !isSave {
		n, err := col.Find(bson.M{"_id": reqMsg.GetCodeMd5()}).Count()
		easygo.PanicError(err)
		if n > 0 {
			return false
		}
	}

	_, erri := col.Upsert(bson.M{"_id": reqMsg.GetCodeMd5()}, bson.M{"$set": reqMsg})
	easygo.PanicError(erri)
	return true
}

//查询终端设备Adv码
func QueryPosDeviceAdv(code string) *share_message.KsPosAdvIdfa {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_POS_ADV_DEVICEIDFA)
	defer closeFun()
	var one *share_message.KsPosAdvIdfa
	err := col.Find(bson.M{"Code": code}).One(&one)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("QueryPosDeviceAdv err:", err.Error())
		return nil
	}
	return one
}

//查询终端设备Adv Idfa码
func QueryPosDeviceAdvIdfa(code string) *share_message.KsPosAdvIdfa {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_POS_ADV_DEVICEIDFA)
	defer closeFun()
	one := &share_message.KsPosAdvIdfa{}
	err := col.Find(bson.M{"_id": code}).One(one)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one
}

//广告idfa激活注册更新
func UpdatePosDeviceAdvIdfa(code, fild string, value bool, ida ...string) {
	idfa := append(ida, "")[0]
	logs.Info("code:", code, idfa)
	var old *share_message.KsPosAdvIdfa
	old = QueryPosDeviceAdv(code)
	if old != nil {
		if old.GetIsActive() && fild == "IsActive" {
			logs.Error("设备已经激活过了")
			return
		}
		if old.GetIsRegister() && fild == "IsRegister" {
			logs.Error("设备已经注册过了")
			return
		}
	}
	md5Code := Md5(code)
	if idfa != "" {
		md5Code = Md5(idfa)
	}
	msg := &share_message.KsPosAdvIdfa{}
	if old == nil {
		old = QueryPosDeviceAdvIdfa(md5Code)
		msg.CodeMd5 = easygo.NewString(md5Code)
	} else {
		msg.CodeMd5 = easygo.NewString(old.GetCodeMd5())
	}
	if old == nil {
		logs.Error("没找到idfa")
		return
	}
	nowTime := util.GetMilliTime()
	switch fild {
	case "IsActive":
		if old.GetIsActive() {
			return
		}
		msg.IsActive = easygo.NewBool(value)
		msg.ActiveTime = easygo.NewInt64(nowTime)
		msg.Code = easygo.NewString(code)
		switch old.GetPlatform() {
		case "ks":
			AdvIdfaPosDevCallBack(old.GetPlatform(), 1, nowTime, old.GetCallback())
		case "qtt":
			AdvIdfaPosDevCallBack(old.GetPlatform(), 0, nowTime, old.GetCallback())
		case "youmi":
			AdvIdfaPosDevCallBack(old.GetPlatform(), 0, nowTime, old.GetCallback())
		}

	case "IsRegister":
		if old.GetIsRegister() {
			return
		}
		msg.IsRegister = easygo.NewBool(value)
		msg.RegisterTime = easygo.NewInt64(nowTime)
		msg.Code = easygo.NewString(code)
		switch old.GetPlatform() {
		case "ks":
			AdvIdfaPosDevCallBack(old.GetPlatform(), 2, nowTime, old.GetCallback())
		case "qtt":
			AdvIdfaPosDevCallBack(old.GetPlatform(), 1, nowTime, old.GetCallback())
			//case "youmi":
			//	AdvIdfaPosDevCallBack(old.GetPlatform(), 0, nowTime, old.GetCallback())
		}
	}
	logs.Info("UpdatePosDeviceAdvIdfa 更新设备信息:", fild, msg)
	SavePosDeviceAdvIdfa(msg, true)
}

//玩家ID查询玩家拥有的群ids
func GetTeamIdsByPlayerId(id PLAYER_ID) []int64 {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAMMEMBER)
	defer closeFun()
	var result []*share_message.PersonalTeamData
	err := col.Find(bson.M{"PlayerId": id}).Select(bson.M{"TeamId": 1, "_id": 0}).All(&result)
	easygo.PanicError(err)

	var teamIds []int64
	for _, v := range result {
		teamIds = append(teamIds, v.GetTeamId())
	}

	return teamIds
}

//查询玩家表情数据
func GetEmoticonFromMongo(id PLAYER_ID) []*share_message.PlayerEmoticon {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_EMOTICON)
	defer closeFun()
	result := []*share_message.PlayerEmoticon{}
	err := col.Find(bson.M{"PlayerId": id}).All(&result)
	if err != mgo.ErrNotFound || err != nil {
		easygo.PanicError(err)
	}
	return result
}

//删除玩家表情数据
func DelEmoticonFromMongo(id PLAYER_ID, typeId int32) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_EMOTICON)
	defer closeFun()
	err := col.Remove(bson.M{"PlayerId": id, "TypeId": typeId})
	if err != mgo.ErrNotFound || err != nil {
		easygo.PanicError(err)
	}
}

//查询组合是否存在
func QueryInterestGroupByGroup(group []int32) *share_message.InterestGroup {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_INTERESTGROUP)
	defer closeFun()

	siteOne := &share_message.InterestGroup{}
	err := col.Find(bson.M{"Group": group}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//=昵称
func GetRandNickName() string {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, "nickname")
	defer closeFun()

	m := []bson.M{
		{"$sample": bson.M{"size": 1}},
	}

	query := col.Pipe(m)
	var list []interface{}
	err := query.All(&list)
	easygo.PanicError(err)

	if list == nil {
		return ""
	}

	info := list[0]
	info1 := info.(bson.M)
	if s, ok := info1["_id"]; ok {
		return s.(string)
	}
	return ""
}

//随机获取签名
func GetRandSignature() string {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, "signature")
	defer closeFun()

	m := []bson.M{
		{"$sample": bson.M{"size": 1}},
	}

	query := col.Pipe(m)
	var list []interface{}
	err := query.All(&list)
	easygo.PanicError(err)

	if list == nil {
		return ""
	}

	info := list[0]
	info1 := info.(bson.M)
	if s, ok := info1["_id"]; ok {
		return s.(string)
	}
	return ""
}

//随机获取省
func GetRandProvice() string {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_DATA_REGION)
	defer closeFun()

	m := []bson.M{
		{"$match": bson.M{"Country": "中国"}},
		{"$sample": bson.M{"size": 1}},
	}

	query := col.Pipe(m)
	var list []interface{}
	err := query.All(&list)
	easygo.PanicError(err)

	if list == nil {
		return ""
	}

	info := list[0]
	info1 := info.(bson.M)
	if s, ok := info1["_id"]; ok {
		return s.(string)
	}
	return ""
}

//随机获取市
func GetRandCity(region string) string {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_DATA_CITY)
	defer closeFun()

	m := []bson.M{
		{"$match": bson.M{"Region": region}},
		{"$sample": bson.M{"size": 1}},
	}

	query := col.Pipe(m)
	var list []interface{}
	err := query.All(&list)
	easygo.PanicError(err)

	if list == nil {
		return ""
	}

	info := list[0]
	info1 := info.(bson.M)
	if s, ok := info1["_id"]; ok {
		return s.(string)
	}
	return ""
}

//随机获取IP
func GetRandPlayerIp() string {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_LOG_LOGIN_INFO)
	defer closeFun()

	m := []bson.M{
		{"$sample": bson.M{"size": 1}},
	}

	query := col.Pipe(m)
	var list []interface{}
	err := query.All(&list)
	easygo.PanicError(err)

	if len(list) == 0 {
		return "127.0.0.1"
	}

	info := list[0]
	info1 := info.(bson.M)
	if s, ok := info1["LoginIP"]; ok {
		return s.(string)
	}
	return "127.0.0.1"
}

func GetRandNickNames(nub int) []string {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, "nickname")
	defer closeFun()

	m := []bson.M{
		{"$sample": bson.M{"size": nub}},
	}

	query := col.Pipe(m)
	var list []interface{}
	err := query.All(&list)
	easygo.PanicError(err)

	if list == nil {
		return nil
	}

	var names []string
	for _, s := range list {
		n := s.(bson.M)
		if name, ok := n["_id"]; ok {
			names = append(names, name.(string))
		}

	}

	return names
}

//随机获取标签
func GetRandLable() []int32 {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_INTERESTTAG)
	defer closeFun()

	m := []bson.M{
		{"$match": bson.M{"Status": 0}},
		{"$sample": bson.M{"size": 3}},
	}

	query := col.Pipe(m)
	var list []interface{}
	err := query.All(&list)
	easygo.PanicError(err)

	lb := []int32{}
	if len(list) == 0 {
		return lb
	}

	for _, info := range list {
		info1 := info.(bson.M)
		if s, ok := info1["_id"]; ok {
			a := s.(int)
			lb = append(lb, int32(a))
		}
	}

	return lb
}

//查询屏蔽词库
func GetDirtyWords() []*share_message.DirtyWords {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_DIRTY_WORDS)
	defer closeFun()
	result := []*share_message.DirtyWords{}
	err := col.Find(bson.M{}).All(&result)
	if err != mgo.ErrNotFound || err != nil {
		easygo.PanicError(err)
	}
	return result
}

//查询玩家注册登录日志
func GetRegisterLoginLog(startTime, endTime int64, channel ...string) []*share_message.LoginRegisterInfo {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_LOGIN_REGISTER_LOG)
	defer closeFun()
	result := []*share_message.LoginRegisterInfo{}
	queryBson := bson.M{"Time": bson.M{"$gte": startTime, "$lt": endTime}}
	if len(channel) > 0 {
		queryBson["Channel"] = channel[0]
	}
	query := col.Find(queryBson)
	err := query.All(&result)
	easygo.PanicError(err)
	return result
}

//查询指定玩家注册日志
func GetRegisterLoginLogById(id PLAYER_ID) *share_message.LoginRegisterInfo {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_LOGIN_REGISTER_LOG)
	defer closeFun()

	queryBson := bson.M{"PlayerId": id}
	queryBson["Type"] = bson.M{"$in": []int32{LOGINREGISTER_PHONEREGISTER, LOGINREGISTER_ONEKEYREGISTER, LOGINREGISTER_WECHATREGISTER}}
	siteOne := &share_message.LoginRegisterInfo{}
	err := col.Find(queryBson).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

func GetRegisterLoginLogByType(types int32) []*share_message.LoginRegisterInfo {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_LOGIN_REGISTER_LOG)
	defer closeFun()
	result := []*share_message.LoginRegisterInfo{}
	queryBson := bson.M{"Type": types}
	query := col.Find(queryBson)
	err := query.All(&result)
	easygo.PanicError(err)
	return result
}

//去重查询返回数组 tableName表名， files 去重字段 ，isLog 是否是log数据库
func GetDistinctArray(tableName, files string, isLog ...bool) []interface{} {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, tableName)
	if len(isLog) > 0 && isLog[0] {
		col, closeFun = easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, tableName)
	}
	defer closeFun()
	var result []interface{}
	err := col.Find(bson.M{}).Distinct(files, &result)
	easygo.PanicError(err)
	return result
}

type InterestTypeEx struct {
	Id         int32 `json:"_id"` //玩家id
	Name       string
	UpdateTime int64
	Sort       int32
	Status     int32
}

func ReloadRedisLabelInfo() []*share_message.InterestType {
	lst := GetInterestTypeList()
	fun := func(data *share_message.InterestType) {
		ex := &InterestTypeEx{}
		StructToOtherStruct(data, ex)
		err := easygo.RedisMgr.GetC().HMSet(MakeRedisKey(REDIS_LABEL_INFO, data.GetId()), ex)
		easygo.PanicError(err)
	}
	for _, v := range lst {
		fun(v)
	}
	return lst
}

func GetRedisLabelInfo(skip ...bool) []*share_message.InterestType {
	var lst []*share_message.InterestType
	keys, err := easygo.RedisMgr.GetC().Scan(REDIS_LABEL_INFO)
	easygo.PanicError(err)
	if len(keys) == 0 {
		return ReloadRedisLabelInfo()
	}
	for _, key := range keys {
		ex := &InterestTypeEx{}
		value, err := redis.Values(easygo.RedisMgr.GetC().HGetAll(key))
		easygo.PanicError(err)
		err = redis.ScanStruct(value, ex)
		easygo.PanicError(err)
		newObj := &share_message.InterestType{}
		StructToOtherStruct(ex, newObj)
		lst = append(lst, newObj)
	}
	return lst
}

// 根据id集合兴趣标签信息列表
func QueryLabelByIds(ids []int32) (labelInfoList []*share_message.LabelInfo) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_INTERESTTAG)
	defer closeFun()
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&labelInfoList)
	easygo.PanicError(err)
	return labelInfoList
}

func GetLabelInfo(labelList []int32) []*share_message.LabelInfo {
	// 暂时先从数据库中获取,后续转入缓存,
	return QueryLabelByIds(labelList)
}

// 下面方法是原来获取兴趣标签信息,暂时不删除
//func GetLabelInfo(labellst []int32) []*share_message.LabelInfo {
//	ids := []string{}
//	for _, id := range labellst {
//		ids = append(ids, easygo.AnytoA(id))
//	}
//	values, err := easygo.RedisMgr.GetC().HMGet(REDIS_SQUARE_DYNAMIC, ids...)
//	easygo.PanicError(err)
//	var infoList []*share_message.LabelInfo
//
//	for _, m := range values {
//		if m == nil {
//			continue
//		}
//		var msg *share_message.InterestType
//		err := json.Unmarshal(m.([]byte), &msg)
//		if err != nil {
//			logs.Error(err)
//			continue
//		}
//		info := &share_message.LabelInfo{
//			Id:   easygo.NewInt32(msg.GetId()),
//			Name: easygo.NewString(msg.GetName()),
//		}
//		infoList = append(infoList, info)
//	}
//	return infoList
//}

//保存注销成功账号
func SaveCancelPhone(log *share_message.CancelAccountList) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_CANCEL_ACCOUNT_LIST)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": log.GetPhone()}, log)
	easygo.PanicError(err)
}

//删除注销成功记录
func DelCancelPhone(phone string) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_CANCEL_ACCOUNT_LIST)
	defer closeFun()
	err := col.RemoveId(phone)
	easygo.PanicError(err)
}

//检测账号注销状态
func CheckCancelAccount(phone string) int64 {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_CANCEL_ACCOUNT_LIST)
	defer closeFun()
	var data share_message.CancelAccountList
	err := col.Find(bson.M{"_id": phone}).One(&data)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	//检测账号是否注销超60天
	if data.GetFinishTime()+60*86400000 < GetMillSecond() {
		err := col.RemoveId(phone)
		if err != nil && err != mgo.ErrNotFound {
			easygo.PanicError(err)
		}
		return 0
	}
	return data.GetFinishTime()
}

//通过玩家id查找注销电话号
func GetPhoneByPlayerId(playerId int64) string {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_CANCEL_ACCOUNT)
	defer closeFun()
	var data *share_message.PlayerCancleAccount
	err := col.Find(bson.M{"PlayerId": playerId}).One(&data)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	return data.GetPhone()
}

//获取登录埋点事件返回
func GetLoginEvent(reqMsg *client_login.LoginEventRequst) *client_login.LoginEventResult {
	playerMrg := GetRedisPlayerBase(reqMsg.GetPlayerId())
	if playerMrg == nil {
		return &client_login.LoginEventResult{}
	}
	codecheck := &share_message.PosDeviceCode{
		CreateTime: easygo.NewInt64(util.GetMilliTime()),
		DeviceCode: easygo.NewString(reqMsg.GetDeviceCode()),
		Channle:    easygo.NewString(playerMrg.GetChannel()),
	}
	isAppAct := SavePosDeviceCode(codecheck)

	list := GetRegisterLoginLog(easygo.GetToday0ClockTimestamp()*1000, easygo.GetToday24ClockTimestamp()*1000)
	logCount := 0
	for _, li := range list {
		if li.GetPlayerId() == reqMsg.GetPlayerId() {
			switch li.GetType() {
			case LOGINREGISTER_PASSWDLOGIN:
				logCount += 1
			case LOGINREGISTER_MESSAGELOGIN:
				logCount += 1
			case LOGINREGISTER_ONEKEYLOGIN:
				logCount += 1
			case LOGINREGISTER_WECHATLOGIN:
				logCount += 1
			case LOGINREGISTER_AUTOLOGIN:
				logCount += 1
			}
		}
	}

	isLoginMan := true
	if logCount > 1 {
		isLoginMan = false
	}

	msg := &client_login.LoginEventResult{
		IsAppAct:   easygo.NewBool(isAppAct),
		IsLoginMan: easygo.NewBool(isLoginMan),
	}
	return msg
}

//查询开启的客服分类列表
func GetManagerTypesListForNormal() []*share_message.ManagerTypes {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_MANAGER_TYPES)
	defer closeFun()

	queryBson := bson.M{"Status": 1} //1开启
	query := col.Find(queryBson)
	var list []*share_message.ManagerTypes
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

//更新玩家所在地库
func UpsertLocation(reqMsg *LocationData) {
	if reqMsg == nil {
		return
	}

	country := &share_message.DataCountry{
		Id:   easygo.NewString(reqMsg.Country),
		Code: easygo.NewString(reqMsg.CountryId),
	}
	if country.GetId() != "" {
		col1, closeFun1 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_DATA_COUNTRY)
		defer closeFun1()
		_, err := col1.Upsert(bson.M{"_id": country.GetId()}, bson.M{"$set": country})
		easygo.PanicError(err)
	}

	area := &share_message.DataArea{
		Id:      easygo.NewString(reqMsg.Area),
		Country: easygo.NewString(reqMsg.Country),
	}
	if area.GetId() != "" {
		col2, closeFun2 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_DATA_AREA)
		defer closeFun2()
		_, err2 := col2.Upsert(bson.M{"_id": area.GetId()}, bson.M{"$set": area})
		easygo.PanicError(err2)
	}

	region := &share_message.DataRegion{
		Id:      easygo.NewString(reqMsg.Region),
		Code:    easygo.NewString(reqMsg.RegionId),
		Country: easygo.NewString(reqMsg.Country),
	}
	if region.GetId() != "" {
		col3, closeFun3 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_DATA_REGION)
		defer closeFun3()

		_, err3 := col3.Upsert(bson.M{"_id": region.GetId()}, bson.M{"$set": region})
		easygo.PanicError(err3)
	}

	city := &share_message.DataCity{
		Id:      easygo.NewString(reqMsg.City),
		Code:    easygo.NewString(reqMsg.CityId),
		Country: easygo.NewString(reqMsg.Country),
		Area:    easygo.NewString(reqMsg.Area),
		Region:  easygo.NewString(reqMsg.Region),
	}
	if city.GetId() != "" {
		col4, closeFun4 := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_DATA_CITY)
		defer closeFun4()

		_, err4 := col4.Upsert(bson.M{"_id": city.GetId()}, bson.M{"$set": city})
		easygo.PanicError(err4)
	}

}

//获取群成员信息
func GetPerTeamDataForMongoDB(teamId int64, playerId PLAYER_ID) *share_message.PersonalTeamData {
	var m *share_message.PersonalTeamData
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAMMEMBER)
	defer closeFun()
	err := col.Find(bson.M{"TeamId": teamId, "PlayerId": playerId}).One(&m)
	if err != nil {
		return nil
	}
	return m
}

//写语音视频时长日志
func SaveVideoVoiceDurationLog(types int32, duration int64, sendId, targetId PLAYER_ID) {
	msg := &share_message.VideoVoiceDurationLog{
		Types:      easygo.NewInt32(types),
		Duration:   easygo.NewInt64(duration),
		CreateTime: easygo.NewInt64(easygo.NowTimestamp()),
		SendId:     easygo.NewInt64(sendId),
		TargetId:   easygo.NewInt64(targetId),
	}

	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_VV_DURSTION_LOG)
	defer closeFun()

	err := col.Insert(msg)
	easygo.PanicError(err)

}

//随机获取指定类型指定数量的用户
func GetRandPlayerByTypes(types []int32, nub int32) []*share_message.PlayerBase {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()

	m := []bson.M{
		{"$match": bson.M{"Types": bson.M{"$in": types}}},
		{"$sample": bson.M{"size": nub}},
	}

	query := col.Pipe(m)
	var list []*share_message.PlayerBase
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

//
func MakePhone() string {
	phone := "1" + easygo.AnytoA(RandInt(3, 9)) + easygo.AnytoA(RandInt(100000000, 999999999))
	aMgr := GetRedisAccountByPhone(phone)
	if aMgr != nil {
		MakePhone()
	}
	return phone
}

//获取硬币道具配置
func GetPropsItemsCfg() *client_hall.PropsItemList {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PROPS_ITEM)
	defer closeFun()
	var params []*share_message.PropsItem
	err := col.Find(bson.M{}).All(&params)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return &client_hall.PropsItemList{Items: params}
}

//获取硬币充值配置
func GetCoinRechargeCfg(way int32, pid int64) *client_hall.CoinRechargeList {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_COIN_RECHARGE)
	defer closeFun()
	var params []*share_message.CoinRecharge
	err := col.Find(bson.M{"Platform": way, "Status": 1}).Sort("-Sort").All(&params)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	logObj := GetRedisCoinLogObj()
	for _, item := range params {
		if !logObj.CheckMonthRecharge(pid, item.GetId()) {
			if !CheckPlayerMonthRecharge(pid, item.GetId()) {
				item.MonthFirst = easygo.NewInt64(0)
			}
		}
	}
	return &client_hall.CoinRechargeList{Items: params}
}

//获取指定类型的虚拟商品
func GetCoinShopList(t int32) []*share_message.CoinProduct {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_COIN_PRODUCT)
	defer closeFun()
	var params []*share_message.CoinProduct
	nowTime := util.GetMilliTime()
	or1 := bson.M{"SaleStartTime": nil, "SaleEndTime": nil}
	or2 := bson.M{"SaleStartTime": 0, "SaleEndTime": 0}
	or3 := bson.M{"SaleStartTime": bson.M{"$lte": nowTime}, "SaleEndTime": bson.M{"$gte": nowTime}}
	queryBson := bson.M{"PropsType": t, "Status": COIN_PRODUCT_STATUS_UP, "$or": []bson.M{or1, or2, or3}}
	err := col.Find(queryBson).Sort("-Sort").All(&params)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return params
}

//获取指定充值id信息
func GetCoinRecharge(id int64) *share_message.CoinRecharge {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_COIN_RECHARGE)
	defer closeFun()
	var item *share_message.CoinRecharge
	err := col.Find(bson.M{"_id": id}).One(&item)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return item
}

//获取指定虚拟物品信息
func GetCoinShopItem(id int64) *share_message.CoinProduct {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_COIN_PRODUCT)
	defer closeFun()
	var item *share_message.CoinProduct
	err := col.Find(bson.M{"_id": id}).One(&item)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return item
}

//获取指定道具物品配置信息
func GetPropsItemInfo(id int64) *share_message.PropsItem {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PROPS_ITEM)
	defer closeFun()
	var item *share_message.PropsItem
	err := col.Find(bson.M{"_id": id}).One(&item)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return item
}

//背包删除道具
func RemovePlayerBagItem(id int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BAG_ITEM)
	defer closeFun()
	err := col.RemoveId(bson.M{"_id": id})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
}

//获取每月充值记录
func CheckPlayerMonthRecharge(playerId, id int64) bool {

	n := FindAllCount(MONGODB_NINGMENG, TABLE_COINCHANGELOG, bson.M{"PlayerId": playerId, "SourceType": COIN_TYPE_EXCHANGE_IN, "CreateTime": bson.M{"$gte": easygo.GetMonth0ClockOfTimestamp(easygo.NowTimestamp()) * 1000}, "Extend.RedPacketId": id})
	return n == 0
}

//func GetNearInfoFromDB(pid int64, reqMsg *client_hall.LocationInfoNewReq, num int) []*share_message.PlayerBase {
//	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
//	defer closeFun()
//	where := bson.M{}
//	//where["_id"] = bson.M{"$ne": pid}
//	if reqMsg.GetSex() != PLAYER_SEX_ALL {
//		where["Sex"] = reqMsg.GetSex()
//	}
//	playerList := make([]*share_message.PlayerBase, 0)
//
//	m := []bson.M{
//		{
//
//			"$geoNear": bson.M{
//				//"near":               []float64{reqMsg.GetX(), reqMsg.GetY()}, // 当前坐标
//				"near": bson.M{
//					"type":        "Points",
//					"coordinates": []float64{113.336020, 23.140620}, // 当前坐标
//				},
//				"spherical":          true,           // 计算球面距离
//				"distanceMultiplier": 6378137,        // 地球半径,单位是米,那么的除的记录也是米
//				"maxDistance":        5000 / 6378137, // 过滤条件100公里以内，需要弧度
//				"distanceField":      "Distance",     // 距离字段别名
//			},
//		},
//		//{"$match": where},
//		//{"$sample": bson.M{"size": num}},
//		//{"$sort": bson.M{"IsOnline": 1}},
//	}
//
//	//bytes, _ := json.Marshal(m)
//	//fmt.Println("map---------->", string(bytes))
//	query := col.Pipe(m)
//
//	err := query.All(&playerList)
//	easygo.PanicError(err)
//	return playerList
//}

func GetNearInfoFromDB(friends []int64, accountType int32, reqMsg *client_hall.LocationInfoNewReq, num int) []*share_message.PlayerBase {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	where := bson.M{}
	where["_id"] = bson.M{"$nin": friends}
	where["Types"] = accountType
	where["Status"] = ACCOUNT_NORMAL
	where["Label"] = bson.M{"$ne": nil}
	if reqMsg.GetSex() != PLAYER_SEX_ALL {
		where["Sex"] = reqMsg.GetSex()
	}
	playerList := make([]*share_message.PlayerBase, 0)
	m := []bson.M{}
	if accountType == ACCOUNT_TYPES_PT { // 普通用户,需要经纬度
		ptM := bson.M{
			"$geoNear": bson.M{
				//"includeLocs":   "location",
				"distanceField": "Distance",
				"maxDistance":   NEAR_DISTANCE, // 100公里以内
				"spherical":     true,
				"near": bson.M{
					"type":        "Point",
					"coordinates": []float64{reqMsg.GetX(), reqMsg.GetY()},
				},
			},
		}
		m = append(m, ptM)
	}
	m = append(m, bson.M{"$match": where})
	m = append(m, bson.M{"$sample": bson.M{"size": num}})

	if reqMsg.GetSort() == NEAR_SORT_ONLINE {
		m = append(m, bson.M{"$sort": bson.M{"IsOnline": -1}})
	}
	query := col.Pipe(m)
	err := query.All(&playerList)
	easygo.PanicError(err)
	return playerList
}

// 获取运营号.
func GetOperationByPhones(pid int64, phones []string, reqMsg *client_hall.LocationInfoNewReq) []*share_message.PlayerBase {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	where := bson.M{}
	where["_id"] = bson.M{"$ne": pid}
	where["Phone"] = bson.M{"$in": phones}
	where["Status"] = ACCOUNT_NORMAL
	// 运营号暂时没有性别区分.手机号没法作区分
	if reqMsg.GetSex() != PLAYER_SEX_ALL {
		where["Sex"] = reqMsg.GetSex()
	}
	playerList := make([]*share_message.PlayerBase, 0)
	m := make([]bson.M, 0)
	m = append(m, bson.M{"$match": where})

	if reqMsg.GetSort() == NEAR_SORT_ONLINE {
		m = append(m, bson.M{"$sort": bson.M{"IsOnline": -1}})
	}
	query := col.Pipe(m)
	err := query.All(&playerList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return playerList
}

func GetAllNearLeadFromDB() []*share_message.NearSet {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_NEAR_LEAD)
	defer closeFun()
	where := bson.M{"Status": 1}
	list := make([]*share_message.NearSet, 0)

	err := col.Find(where).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return list

}

// 附近的人,用户推荐区域列表
func GetNearRecommendPlayer(pid int64, page, pageSize int, x, y float64, area string) ([]*share_message.PlayerBase, int) {
	list := make([]*share_message.PlayerBase, 0)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	t := GetMillSecond() - 3*86400*1000 // 3天
	curPage := easygo.If(page > 1, page-1, 0).(int)
	//queryBson := bson.M{"Status": ACCOUNT_NORMAL, "$or": []bson.M{{"IsOnline": true}, {"IsRecommend": true}, {"Area": area}, {"LastLogOutTime": bson.M{"$gt": t}}}}
	bsonArr := []bson.M{{"IsOnline": true}, {"IsRecommend": true}, {"LastLogOutTime": bson.M{"$gt": t}}}
	if area != "" {
		bsonArr = append(bsonArr, bson.M{"Area": area})
	}
	queryBson := bson.M{"Status": ACCOUNT_NORMAL, "_id": bson.M{"$ne": pid}, "$or": bsonArr}
	m := []bson.M{
		{
			"$geoNear": bson.M{
				"distanceField": "Distance",
				"minDistance":   NEAR_NEAR_DISTANCE, // 50公里以外
				"spherical":     true,
				"near": bson.M{
					"type":        "Point",
					"coordinates": []float64{x, y},
				},
			},
		},
		{"$match": queryBson},
	}

	query := col.Pipe(m)
	err := query.All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	count := len(list)
	skip := bson.M{"$skip": curPage * pageSize}
	limit := bson.M{"$limit": pageSize}

	m = append(m, skip)
	m = append(m, limit)
	list1 := make([]*share_message.PlayerBase, 0)
	query1 := col.Pipe(m)
	err = query1.All(&list1)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return list1, count
}

// 添加会话列表
func AddSessionList(sendId, receiveId int64, reqMsg *share_message.NearSessionList) int64 {
	createTime := GetMillSecond()
	id := MakeNewString(sendId, receiveId)
	// 会话列表
	reqMsg.Id = easygo.NewString(id)
	reqMsg.SendPlayerId = easygo.NewInt64(sendId)
	reqMsg.Status = easygo.NewInt32(NEAR_MESSAGE_NORMAL)
	reqMsg.IsRead = easygo.NewBool(false)
	reqMsg.CreateTime = easygo.NewInt64(createTime)
	AddSessionToDB(reqMsg)

	// 埋点使用  判断今天是否是第一条数据打招呼数据.
	var isFirst bool
	if count := GetIsFirstSayHi(sendId, receiveId); count == 0 {
		isFirst = true
	}
	contentId := NextId(TABLE_NEARBY_MESSAGE_NEW_LOG)
	newLog := &share_message.NearMessageNewLog{
		Id:              easygo.NewInt64(contentId),
		SendPlayerId:    easygo.NewInt64(sendId),
		ReceivePlayerId: easygo.NewInt64(receiveId),
		Content:         easygo.NewString(reqMsg.GetContent()),
		ContentType:     easygo.NewInt32(reqMsg.GetContentType()),
		Status:          easygo.NewInt32(NEAR_MESSAGE_NORMAL),
		IsRead:          easygo.NewBool(false),
		CreateTime:      easygo.NewInt64(createTime),
		IsFirst:         easygo.NewBool(isFirst),
	}
	InsertNearMessageNewLog(newLog)
	return contentId
}

// 保存附近的人打招呼会话列表
func AddSessionToDB(req *share_message.NearSessionList) {
	req.UpdateTime = easygo.NewInt64(GetMillSecond())
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_SESSIOIN_LIST)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": req.GetId()}, bson.M{"$set": req})
	easygo.PanicError(err)
}

func NearSessionList(pid int64, reqMsg *client_hall.NearSessionListReq) *client_hall.NearSessionListResp {
	if reqMsg.GetPage() == 0 {
		reqMsg.Page = easygo.NewInt64(DEFAULT_PAGE)
	}
	if reqMsg.GetPageSize() == 0 {
		reqMsg.PageSize = easygo.NewInt64(DEFAULT_PAGE_SIZE)
	}
	if reqMsg.GetPage() == 1 {
		reqMsg.QueryTime = easygo.NewInt64(GetMillSecond())
	}
	sayMessages, count := GetNearSessionListByPageFromDB(pid, int(reqMsg.GetPage()), int(reqMsg.GetPageSize()), reqMsg.GetQueryTime())
	sessionList := make([]*client_hall.NearSessionList, 0)
	// 计算距离
	recPlayer := GetRedisPlayerBase(pid) // 接收方
	friends := recPlayer.GetFriends()
	var distance float64
	recGeoJson := recPlayer.GeRedisPoints()
	for _, v := range sayMessages {
		isFriend := util.Int64InSlice(v.GetSendPlayerId(), friends)
		base := GetRedisPlayerBase(v.GetSendPlayerId())
		sendGeoJson := base.GeRedisPoints()
		if recGeoJson != nil && len(recGeoJson.GetCoordinates()) == 2 && sendGeoJson != nil && len(sendGeoJson.GetCoordinates()) == 2 { // 计算距离
			distance = GetDistance(recGeoJson.GetCoordinates()[0], sendGeoJson.GetCoordinates()[0],
				recGeoJson.GetCoordinates()[1], sendGeoJson.GetCoordinates()[1])
		}
		equipmentObj := GetRedisPlayerEquipmentObj(v.GetSendPlayerId())
		equipment := equipmentObj.GetEquipmentForClient()
		sessionList = append(sessionList, &client_hall.NearSessionList{
			SendPlayerId: easygo.NewInt64(v.GetSendPlayerId()),
			NickName:     easygo.NewString(base.GetNickName()),
			HeadIcon:     easygo.NewString(base.GetHeadIcon()),
			Sex:          easygo.NewInt32(base.GetSex()),
			Content:      easygo.NewString(v.GetContent()),
			IsRead:       easygo.NewBool(v.GetIsRead()),
			ContentType:  easygo.NewInt32(v.GetContentType()),
			Distance:     easygo.NewFloat64(distance),
			IsFriend:     easygo.NewBool(isFriend),
			PropsId:      easygo.NewInt64(equipment.GetQP().GetPropsId()),
		})
	}

	return &client_hall.NearSessionListResp{
		SessionList: sessionList,
		Count:       easygo.NewInt32(count),
		QueryTime:   easygo.NewInt64(reqMsg.GetQueryTime()),
	}
}

func GetNearSessionListByPageFromDB(pid int64, page, pageSize int, queryTime int64) ([]*share_message.NearSessionList, int) {
	curPage := page - 1
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_SESSIOIN_LIST)
	defer closeFun()

	list := make([]*share_message.NearSessionList, 0)
	t := GetMillSecond() - TIME_BEFORE
	query := col.Find(bson.M{"Status": NEAR_MESSAGE_NORMAL, "ReceivePlayerId": pid, "CreateTime": bson.M{"$gt": t, "$lt": queryTime}})
	count, err := query.Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	//err = query.Sort("IsRead", "-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	err = query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return list, count
}

func InsertNearMessageNewLog(msg *share_message.NearMessageNewLog) {
	msg.UpdateTime = easygo.NewInt64(GetMillSecond())
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_MESSAGE_NEW_LOG)
	defer closeFun()
	err := col.Insert(msg)
	easygo.PanicError(err)
}

// 具体到某个人打招呼的列表
func GetNearSayMessageList(receivePlayerId int64, reqMsg *client_hall.SendPlayerMessageListReq) ([]*client_hall.NearSessionList, int) {
	if reqMsg.GetPage() == 0 {
		reqMsg.Page = easygo.NewInt64(DEFAULT_PAGE)
	}

	if reqMsg.GetPageSize() == 0 {
		reqMsg.PageSize = easygo.NewInt64(DEFAULT_PAGE_SIZE)
	}
	sendPlayerId := reqMsg.GetSendPlayerId()
	newLogs, count := GetSayMessageListFromDB(sendPlayerId, receivePlayerId, reqMsg.GetMaxId(), int(reqMsg.GetPage()), int(reqMsg.GetPageSize()))
	sendPlayer := GetRedisPlayerBase(sendPlayerId)
	sessionList := make([]*client_hall.NearSessionList, 0)
	equipmentObj := GetRedisPlayerEquipmentObj(sendPlayerId)
	equipment := equipmentObj.GetEquipmentForClient()
	for _, v := range newLogs {
		sessionList = append(sessionList, &client_hall.NearSessionList{
			Id:           easygo.NewInt64(v.GetId()),
			SendPlayerId: easygo.NewInt64(sendPlayerId),
			NickName:     easygo.NewString(sendPlayer.GetNickName()),
			HeadIcon:     easygo.NewString(sendPlayer.GetHeadIcon()),
			Sex:          easygo.NewInt32(sendPlayer.GetSex()),
			Content:      easygo.NewString(v.GetContent()),
			IsRead:       easygo.NewBool(true),
			ContentType:  easygo.NewInt32(v.GetContentType()),
			CreateTime:   easygo.NewInt64(v.GetCreateTime()),
			PropsId:      easygo.NewInt64(equipment.GetQP().GetPropsId()),
		})
		v.IsRead = easygo.NewBool(true)
		// 设置成已读
		UpsertMessageNewLog(v)
	}
	// 异步判断会话是否还有未读
	fun := func() {
		checkHaveUnread(sendPlayerId, receivePlayerId)
	}
	easygo.Spawn(fun)

	return sessionList, count
}

// 判断会话是否还有未读
func checkHaveUnread(sendPlayerId, receivePlayerId int64) {
	// 判断会话是否还有未读
	if GetPlayerHaveUnRead(sendPlayerId, receivePlayerId) == 0 {
		if s := GetSessionFromDB(MakeNewString(sendPlayerId, receivePlayerId)); s != nil {
			s.IsRead = easygo.NewBool(true)
			AddSessionToDB(s)
		}
	}

	// 判断接收人身上是否还有未读
	if GetPlayerHaveUnReadAll(receivePlayerId) == 0 {
		receivePlayer := GetRedisPlayerBase(receivePlayerId)
		receivePlayer.SetIsNearBy(false)
	}
}
func GetSayMessageListFromDB(sendPlayerId, receivePlayerId, maxId int64, page, pageSize int) ([]*share_message.NearMessageNewLog, int) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_MESSAGE_NEW_LOG)
	defer closeFun()
	curPage := page - 1
	logList := make([]*share_message.NearMessageNewLog, 0)
	t := GetMillSecond() - TIME_BEFORE
	query := col.Find(bson.M{"Status": NEAR_MESSAGE_NORMAL, "SendPlayerId": sendPlayerId, "ReceivePlayerId": receivePlayerId, "CreateTime": bson.M{"$gt": t}})
	if page == DEFAULT_PAGE {
		maxId = 0
	}
	if maxId != 0 {
		query = col.Find(bson.M{"Status": NEAR_MESSAGE_NORMAL, "SendPlayerId": sendPlayerId, "ReceivePlayerId": receivePlayerId, "CreateTime": bson.M{"$gt": t}, "_id": bson.M{"$lte": maxId}})
	}
	count, err := query.Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	err = query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&logList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return logList, count
}

// 修改状态
func UpsertMessageNewLog(req *share_message.NearMessageNewLog) {
	req.UpdateTime = easygo.NewInt64(GetMillSecond())
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_MESSAGE_NEW_LOG)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": req.GetId()}, bson.M{"$set": req})
	easygo.PanicError(err)
}

// 收信人是否还有未读(所有的)
func GetPlayerHaveUnReadAll(receivePlayerId int64) int {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_MESSAGE_NEW_LOG)
	defer closeFun()

	t := GetMillSecond() - TIME_BEFORE
	count, err := col.Find(bson.M{"Status": NEAR_MESSAGE_NORMAL, "ReceivePlayerId": receivePlayerId, "IsRead": false, "CreateTime": bson.M{"$gt": t}}).Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return count
}

// 会话是否还有未读
func GetPlayerHaveUnRead(sendPlayerId, receivePlayerId int64) int {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_MESSAGE_NEW_LOG)
	defer closeFun()

	t := GetMillSecond() - TIME_BEFORE
	count, err := col.Find(bson.M{"Status": NEAR_MESSAGE_NORMAL, "SendPlayerId": sendPlayerId, "ReceivePlayerId": receivePlayerId, "IsRead": false, "CreateTime": bson.M{"$gt": t}}).Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return count
}

// 保存附近的人打招呼会话列表
func GetSessionFromDB(id string) *share_message.NearSessionList {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_SESSIOIN_LIST)
	defer closeFun()
	var data *share_message.NearSessionList
	err := col.Find(bson.M{"_id": id}).One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return data
}

// 获取附近的人聊天内容列表 sendPlayerId  自己, receivePlayerId  对方
//func GetNearChatList(sendPlayerId, receivePlayerId, page, pageSize int64) ([]*client_hall.NearSessionList, int) {
func GetNearChatList(sendPlayerId int64, reqMsg *client_hall.GetNearChatListReq) ([]*client_hall.NearSessionList, int) {
	page := reqMsg.GetPage()
	pageSize := reqMsg.GetPageSize()
	receivePlayerId := reqMsg.GetReceivePlayerId()
	if page == 0 {
		page = DEFAULT_PAGE
	}

	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}

	newLogs, count := GetNearChatListFromDB(sendPlayerId, receivePlayerId, reqMsg.GetMaxId(), int(page), int(pageSize))

	sessionList := make([]*client_hall.NearSessionList, 0)
	for _, v := range newLogs {
		sendPlayer := GetRedisPlayerBase(v.GetSendPlayerId())
		equipmentObj := GetRedisPlayerEquipmentObj(v.GetSendPlayerId())
		equipment := equipmentObj.GetEquipmentForClient()

		sessionList = append(sessionList, &client_hall.NearSessionList{
			Id:           easygo.NewInt64(v.GetId()),
			SendPlayerId: easygo.NewInt64(v.GetSendPlayerId()),
			NickName:     easygo.NewString(sendPlayer.GetNickName()),
			HeadIcon:     easygo.NewString(sendPlayer.GetHeadIcon()),
			Sex:          easygo.NewInt32(sendPlayer.GetSex()),
			Content:      easygo.NewString(v.GetContent()),
			IsRead:       easygo.NewBool(true),
			ContentType:  easygo.NewInt32(v.GetContentType()),
			CreateTime:   easygo.NewInt64(v.GetCreateTime()),
			PropsId:      easygo.NewInt64(equipment.GetQP().GetPropsId()),
		})
		v.IsRead = easygo.NewBool(true)
	}
	// 异步判断会话是否还有未读
	fun := func() {
		// 设置成已读
		for _, v := range newLogs {
			UpsertMessageNewLog(v)
		}
		checkHaveUnread(sendPlayerId, receivePlayerId)
		checkHaveUnread(receivePlayerId, sendPlayerId)
	}

	easygo.Spawn(fun)
	return sessionList, count
}

//获取附近的人聊天内容列表
func GetNearChatListFromDB(sendPlayerId, receivePlayerId, maxId int64, page, pageSize int) ([]*share_message.NearMessageNewLog, int) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_MESSAGE_NEW_LOG)
	defer closeFun()
	curPage := page - 1
	logList := make([]*share_message.NearMessageNewLog, 0)
	t := GetMillSecond() - TIME_BEFORE
	query := col.Find(bson.M{"Status": NEAR_MESSAGE_NORMAL, "SendPlayerId": bson.M{"$in": []int64{sendPlayerId, receivePlayerId}},
		"ReceivePlayerId": bson.M{"$in": []int64{sendPlayerId, receivePlayerId}}, "CreateTime": bson.M{"$gt": t}})

	if page == DEFAULT_PAGE {
		maxId = 0
	}
	if maxId != 0 {
		query = col.Find(bson.M{"_id": bson.M{"$lte": maxId}, "Status": NEAR_MESSAGE_NORMAL, "SendPlayerId": bson.M{"$in": []int64{sendPlayerId, receivePlayerId}},
			"ReceivePlayerId": bson.M{"$in": []int64{sendPlayerId, receivePlayerId}}, "CreateTime": bson.M{"$gt": t}})
	}

	count, err := query.Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	err = query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&logList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return logList, count
}

// 删除聊天内容
func DelNearMessage(sendPlayerId []int64, receivePlayerId int64) {
	// 修改会话列表的状态
	for _, v := range sendPlayerId {
		UpdateSessionStatusById(MakeNewString(v, receivePlayerId))
		// 修改具体的聊天内容.
		UpdateNearMessageStatusById(v, receivePlayerId)
	}

	// 异步检测是否还有未读
	fun := func() {
		for _, v := range sendPlayerId {
			//检测是否还有未读
			checkHaveUnread(v, receivePlayerId)
		}
	}
	easygo.Spawn(fun)
}

func UpdateSessionStatusById(id string) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_SESSIOIN_LIST)
	defer closeFun()
	t := GetMillSecond()
	err := col.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"Status": NEAR_MESSAGE_DELETE, "UpdateTime": t}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
}
func UpdateNearMessageStatusById(sendPlayerId, receivePlayerId int64) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_MESSAGE_NEW_LOG)
	defer closeFun()
	t := GetMillSecond()
	_, err := col.UpdateAll(bson.M{"SendPlayerId": sendPlayerId, "ReceivePlayerId": receivePlayerId}, bson.M{"$set": bson.M{"Status": NEAR_MESSAGE_DELETE, "UpdateTime": t}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
}

// 修改聊天内容已读
func UpdateIsRead(contentId int64) {
	// 修改聊天内容已读
	UpdateIsReadById(contentId)
	// 查询出发送人,接收人
	fun := func() {
		if messageLog := GetMessageLogByIdFromDB(contentId); messageLog != nil {
			// 判断是否还有未读已读
			checkHaveUnread(messageLog.GetSendPlayerId(), messageLog.GetReceivePlayerId())
		}
	}
	easygo.Spawn(fun)

}
func UpdateIsReadById(contentId int64) {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_MESSAGE_NEW_LOG)
	defer closeFun()
	t := GetMillSecond()
	err := col.Update(bson.M{"_id": contentId}, bson.M{"$set": bson.M{"IsRead": true, "UpdateTime": t}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
}

func GetMessageLogByIdFromDB(contentId int64) *share_message.NearMessageNewLog {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_MESSAGE_NEW_LOG)
	defer closeFun()

	var messageLog *share_message.NearMessageNewLog
	err := col.Find(bson.M{"_id": contentId}).One(&messageLog)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return messageLog
}

//写道具获得日志
func AddGetPropsLog(data *share_message.PlayerGetPropsLog) {
	// data := &share_message.PlayerGetPropsLog{
	// 	Id:            easygo.NewInt64(NextId(TABLE_PLAYER_GETPROPS_LOG)),
	// 	PlayerId:      easygo.NewInt64(Pid),
	// 	GivePlayerId:  easygo.NewInt64(givePid),
	// 	PropsId:       easygo.NewInt64(propId),
	// 	PropsNum:      easygo.NewInt64(count),
	// 	GetType:       easygo.NewInt32(way),
	// 	CreateTime:    easygo.NewInt64(easygo.NowTimestamp()),
	// 	EffectiveTime: easygo.NewInt64(days),
	// 	BagId:         easygo.NewInt64(bagId),
	// 	BuyWay:        easygo.NewInt32(bugWay),
	// 	OrderId:       easygo.NewString(orderId),
	// }

	// if len(oprater) > 0 {
	// 	data.Operator = easygo.NewString(oprater[0])
	// 	data.RecycleTime = easygo.NewInt64(easygo.NowTimestamp())
	// }
	var datas []interface{}
	datas = append(datas, data)
	InsertAllMgo(MONGODB_NINGMENG, TABLE_PLAYER_GETPROPS_LOG, datas...)
}

func GetLimitPlayerBase() []*share_message.PlayerBase {
	colVar, closeFunVar := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFunVar()
	players := make([]*share_message.PlayerBase, 0)
	err := colVar.Find(bson.M{"Points1": nil}).Limit(5000).All(&players)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return players
}

// 是否是第一条打招呼信息
func GetIsFirstSayHi(sendPlayerId, receivePlayerId int64) int {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_NEARBY_MESSAGE_NEW_LOG)
	defer closeFun()

	t := easygo.GetToday0ClockTimestamp() * 1000 // 当天的时间戳 毫秒
	query := col.Find(bson.M{"Status": NEAR_MESSAGE_NORMAL, "SendPlayerId": bson.M{"$in": []int64{sendPlayerId, receivePlayerId}},
		"ReceivePlayerId": bson.M{"$in": []int64{sendPlayerId, receivePlayerId}}, "CreateTime": bson.M{"$gt": t}})

	count, err := query.Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	return count
}

//通过会话id获取会话数据
func GetSessionDataByPlayerId(pid int64) []*share_message.ChatSession {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_MONGODB_CHAT_SESSION)
	defer closeFun()
	var session []*share_message.ChatSession
	err := col.Find(bson.M{"PlayerIds": pid}).All(&session)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return session
}

//获取指定ids的会话
func GetAllSessionData(ids []string) []*share_message.ChatSession {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_MONGODB_CHAT_SESSION)
	defer closeFun()
	var sessions []*share_message.ChatSession
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&sessions)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return sessions
}

//获取玩家会话
func GetMySessions(pid int64) *share_message.PlayerChatSession {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_CHAT_SESSION)
	defer closeFun()
	var session *share_message.PlayerChatSession
	err := col.Find(bson.M{"_id": pid}).One(&session)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	} else if err != nil && err.Error() == mgo.ErrNotFound.Error() {
		return &share_message.PlayerChatSession{
			PlayerId:   easygo.NewInt64(pid),
			SessionIds: []string{},
		}
	}
	return session
}

//保存我的会话列表
func SaveMySessions(pid int64, sessions []string) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_CHAT_SESSION)
	defer closeFun()
	playerSession := &share_message.PlayerChatSession{
		PlayerId:   easygo.NewInt64(pid),
		SessionIds: sessions,
	}
	_, err := col.Upsert(bson.M{"_id": pid}, playerSession)
	easygo.PanicError(err)
}

//保存我的会话列表
func DeleteMySessions(pid int64, ids []string) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_CHAT_SESSION)
	defer closeFun()
	err1 := col.Update(bson.M{"_id": pid}, bson.M{"$pullAll": bson.M{"SessionIds": ids}})
	if err1 != nil {
		easygo.PanicError(err1)
	}
}

//群成员查询@分页20一页
func GetTeamAtMember(pid, teamId int64, page int32) []*client_hall.AtData {
	//pageSize := 20
	//curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAMMEMBER)
	defer closeFun()
	var list []*share_message.PersonalTeamData
	//err := col.Find(bson.M{"TeamId": teamId}).Sort("NickName").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	err := col.Find(bson.M{"TeamId": teamId}).Sort("Position").All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	data := make([]*client_hall.AtData, 0)
	teamObj := GetRedisTeamObj(teamId)
	ids := teamObj.GetTeamMemberList()
	pMap := GetAllPlayerBase(ids, false)
	for _, m := range list {
		p := pMap[m.GetPlayerId()]
		name := m.GetNickName()
		if name == "" && p != nil {
			name = p.GetNickName()
		}
		at := &client_hall.AtData{
			PlayerId: easygo.NewInt64(m.GetPlayerId()),
			Name:     easygo.NewString(name),
			Position: easygo.NewInt32(m.GetPosition()),
		}
		if p != nil {
			at.HeadUrl = easygo.NewString(p.GetHeadIcon())
			at.Sex = easygo.NewInt32(p.GetSex())
		}
		data = append(data, at)
	}
	return data
}

//群成员获取，每次30个
func GetTeamMemberDatas(pid, teamId int64, page int32) []*share_message.PersonalTeamData {
	pageSize := 30
	curPage := easygo.If(page > 1, page-1, 0).(int)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAMMEMBER)
	defer closeFun()
	var list []*share_message.PersonalTeamData
	err := col.Find(bson.M{"TeamId": teamId}).Sort("Position").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return list
}

//获取玩家保存列表
func GetMySaveTeamIds(pid int64) []string {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAMMEMBER)
	defer closeFun()
	var list []*share_message.PersonalTeamData
	err := col.Find(bson.M{"PlayerId": pid, "Setting.IsSaveAdd": true}).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	ids := make([]string, 0)
	for _, m := range list {
		ids = append(ids, easygo.AnytoA(m.GetTeamId()))
	}
	return ids
}

//获取我的录音作品列表
func GetMyMixVideo(me, pid int64) []*share_message.PlayerMixVoiceVideo {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_MIX_VOICE_VIDEO)
	defer closeFun()
	var list []*share_message.PlayerMixVoiceVideo
	//ls := make([]int32, 0)
	//if me == pid {
	//ls = []int32{VC_STATUS_UNCHECK, VC_STATUS_CHECKED, VC_STATUS_USING}
	//} else {
	//	ls = []int32{VC_STATUS_CHECKED, VC_STATUS_USING}
	//}
	q := bson.M{"PlayerId": pid, "Status": bson.M{"$ne": VC_STATUS_DELETE}}
	err := col.Find(q).Sort("-CreateTime").All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return list
}

//获取指定玩家的最新录音作品
func GetNewMyMixVideo(pid int64) *share_message.PlayerMixVoiceVideo {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_MIX_VOICE_VIDEO)
	defer closeFun()
	var data *share_message.PlayerMixVoiceVideo
	err := col.Find(bson.M{"PlayerId": pid, "Status": bson.M{"$lt": VC_STATUS_DELETE}}).Sort("-CreateTime").One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return data
}

//修改录音作品状态
func ModifyMyMixVideoStatus(mixId int64, st int32) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_MIX_VOICE_VIDEO)
	defer closeFun()
	err := col.Update(bson.M{"_id": mixId}, bson.M{"$set": bson.M{"Status": st}})
	easygo.PanicError(err)
}
func ModifyMyMixVideoUseing(mixId int64, b bool) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_MIX_VOICE_VIDEO)
	defer closeFun()
	err := col.Update(bson.M{"_id": mixId}, bson.M{"$set": bson.M{"IsUse": b}})
	easygo.PanicError(err)
}

//获取背景录像数据
func GetBgVideoData(ids []int64) []*share_message.BgVoiceVideo {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_BG_VOICE_VIDEO)
	defer closeFun()
	var list []*share_message.BgVoiceVideo
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return list
}

//获取指定id的背景音乐
func GetOneBgVideoData(id int64) *share_message.BgVoiceVideo {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_BG_VOICE_VIDEO)
	defer closeFun()
	var one *share_message.BgVoiceVideo
	err := col.Find(bson.M{"_id": id}).One(&one)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	err = col.Update(bson.M{"_id": id}, bson.M{"$inc": bson.M{"UseCount": 1}})
	if err != nil {
		logs.Error("找不到背景音频信息:", id)
	}
	return one
}

//获取背景录像指定标签的名称
func GetBgVideoTags(tags []int32) map[int32]*share_message.InterestTag {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_BG_VOICE_TAG)
	defer closeFun()
	var list []*share_message.InterestTag
	err := col.Find(bson.M{"_id": bson.M{"$in": tags}}).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	data := make(map[int32]*share_message.InterestTag, 0)
	for _, d := range list {
		data[d.GetId()] = d
	}
	return data
}

//玩家个性化标签
func GetPlayerPersonalityTags(tags []int32) []*share_message.InterestTag {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_CHARACTER_TAG)
	defer closeFun()
	var list []*share_message.InterestTag
	err := col.Find(bson.M{"_id": bson.M{"$in": tags}}).All(&list)
	if err != nil {
		easygo.PanicError(err)
	}
	return list
}
func GetPlayerPersonalityTagsMp(tags []int32) map[int32]*share_message.InterestTag {
	data := GetPlayerPersonalityTags(tags)
	mp := make(map[int32]*share_message.InterestTag, 0)
	for _, d := range data {
		mp[d.GetId()] = d
	}
	return mp
}

//玩家个性化标签
func GetPlayerPersonalityTagAllData(tags []int32) []*share_message.InterestTag {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_CHARACTER_TAG)
	defer closeFun()
	var list []*share_message.InterestTag
	err := col.Find(bson.M{"_id": bson.M{"$in": tags}}).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取玩家个性化数据失败:%v", err)
		return nil
	}
	return list
}

//随机获取指定个数的个性化标签
func GetPagePersonalityTags(num, page int) []*share_message.InterestTag {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_CHARACTER_TAG)
	defer closeFun()
	var list []*share_message.InterestTag
	//只取没有封禁的标签
	err := col.Find(bson.M{"Status": bson.M{"$ne": 1}}).Sort("-_id").Skip((page - 1) * num).Limit(num).All(&list)
	if err != nil {
		easygo.PanicError(err)
	}
	return list
}

//获取星座名字
func GetPlayerConstellationStr(id int32) string {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_STARSIGNS_TAG)
	defer closeFun()
	var data *share_message.InterestTag
	err := col.Find(bson.M{"_id": id}).One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return data.GetName()
}

//获取系统提供的所有星座信息
func GetConfigConstellationFormDB() []*share_message.InterestTag {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_STARSIGNS_TAG)
	defer closeFun()
	var data []*share_message.InterestTag
	err := col.Find(bson.M{}).All(&data)
	if err != nil {
		easygo.PanicError(err)
	}
	return data
}

//语音名片点赞
func CheckVoiceCardZan(pid, target int64) *share_message.PlayerVCZanLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_VC_ZAN_LOG)
	defer closeFun()
	var one *share_message.PlayerVCZanLog
	err := col.Find(bson.M{"PlayerId": pid, "TargetId": target}).One(&one)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return one

}

//删除语音点赞记录
func AddVoiceCardZanNum(logId int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_VC_ZAN_LOG)
	defer closeFun()
	t := time.Now().Unix()
	err := col.Update(bson.M{"_id": logId}, bson.M{"$inc": bson.M{"ZanNum": 1}, "$set": bson.M{"CreateTime": t}})
	easygo.PanicError(err)
}

//增加语音点赞记录
func AddVoiceCardZanLog(log *share_message.PlayerVCZanLog) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_VC_ZAN_LOG)
	defer closeFun()
	err := col.Insert(log)
	if err != nil {
		easygo.PanicError(err)
	}
}

//语音SayHi:只要有一方打过招呼，就认为打过招呼
func CheckSayHiToPlayer(pid, target int64) *share_message.PlayerVCSayHiLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_VC_SAY_HI_LOG)
	defer closeFun()
	var one *share_message.PlayerVCSayHiLog
	//q1 := bson.M{"PlayerId": pid, "TargetId": target}
	//q2 := bson.M{"PlayerId": target, "TargetId": pid}
	//q := bson.M{"$or": []bson.M{q1, q2}}
	err := col.Find(bson.M{"PlayerId": pid, "TargetId": target}).One(&one)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return one

}

//增加语音SayHi记录
func AddVoiceCardSayHiLog(log *share_message.PlayerVCSayHiLog) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_VC_SAY_HI_LOG)
	defer closeFun()
	err := col.Insert(log)
	if err != nil {
		easygo.PanicError(err)
	}
}

//删除语音SayHi记录
func DelVoiceCardSayHiLog(playerId, targetId int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_VC_SAY_HI_LOG)
	defer closeFun()
	err := col.Remove(bson.M{"PlayerId": playerId, "TargetId": targetId})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
}

//增加语音关注记录
func AddAttentionLog(playerId, targetId int64, opt int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_ATTENTION_LOG)
	defer closeFun()
	var one *share_message.PlayerAttentionLog
	t := GetMillSecond()
	err := col.Find(bson.M{"PlayerId": playerId, "TargetId": targetId}).One(&one)
	if err != nil {
		if err == mgo.ErrNotFound {
			one = &share_message.PlayerAttentionLog{
				Id:       easygo.NewInt64(NextId(TABLE_PLAYER_ATTENTION_LOG)),
				PlayerId: easygo.NewInt64(playerId),
				TargetId: easygo.NewInt64(targetId),
				SortTime: easygo.NewInt64(t),
				Opt:      easygo.NewInt32(opt),
			}
			if opt == VC_ATTENTION_HI {
				one.SayHiTime = easygo.NewInt64(t)
			}
			err = col.Insert(one)
		} else {
			logs.Error("添加关注记录失败%v", err)
		}
		return
	}
	upBson := bson.M{}
	upBson["SortTime"] = easygo.NewInt64(t)
	if opt == VC_ATTENTION_HI {
		upBson["Opt"] = VC_ATTENTION_HI
		upBson["SayHiTime"] = easygo.NewInt64(t)
		upBson["SortTime"] = easygo.NewInt64(t)
	}
	err = col.Update(bson.M{"_id": one.GetId()}, bson.M{"$set": upBson})
	if err != nil {
		logs.Error("更新关注记录失败%v", err)
	}
}

//删除语音关注记录
func DelAttentionLog(playerId, targetId int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_ATTENTION_LOG)
	defer closeFun()
	var one *share_message.PlayerAttentionLog
	err := col.Find(bson.M{"PlayerId": playerId, "TargetId": targetId}).One(&one)
	if err != nil {
		logs.Error("查找语音关注记录失败%v", err)
		return
	}
	//此处只有点赞能删除记录，若同时是打招呼则恢复排序时间
	if one.GetOpt() == VC_ATTENTION_LIKE {
		err = col.Remove(bson.M{"PlayerId": playerId, "TargetId": targetId, "Opt": VC_ATTENTION_LIKE})
	} else {
		err = col.Update(bson.M{"_id": one.GetId()}, bson.M{"$set": bson.M{"SortTime": one.GetSayHiTime()}})
	}
	if err != nil {
		logs.Error("删除语音关注记录失败%v", err)
	}
}

//获取语音关注喜欢我的人
func GetAttentionPlayers(playerId int64, page, cnt, opt int) []*share_message.PlayerAttentionLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_ATTENTION_LOG)
	defer closeFun()
	page -= 1
	queryBson := bson.M{}
	if opt == VC_ATTENTION_TO_ME {
		queryBson["TargetId"] = playerId
	} else {
		queryBson["PlayerId"] = playerId
	}
	var data []*share_message.PlayerAttentionLog
	err := col.Find(queryBson).Sort("-SortTime").Skip(page * cnt).Limit(cnt).All(&data)
	if err != nil {
		logs.Error("查找语音关注记录失败%v", err)
		return nil
	}
	return data
}

//获取喜欢我新消息条目
func GetLoveMeNewNum(pid, t int64) (int, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_ATTENTION_LOG)
	defer closeFun()
	n, err := col.Find(bson.M{"TargetId": pid, "SortTime": bson.M{"$gte": t}}).Count()
	easygo.PanicError(err)
	total, err := col.Find(bson.M{"TargetId": pid}).Count()
	return n, total
}

//查询指定录音ID的信息
func GetVoiceCardInfo(id int64) *share_message.PlayerMixVoiceVideo {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_MIX_VOICE_VIDEO)
	defer closeFun()
	var one *share_message.PlayerMixVoiceVideo
	err := col.Find(bson.M{"_id": id}).One(&one)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return one
}

//查询指定录音ID的信息
func GetVoiceCardInfoByIds(ids []int64) []*share_message.PlayerMixVoiceVideo {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_MIX_VOICE_VIDEO)
	defer closeFun()
	var all []*share_message.PlayerMixVoiceVideo
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&all)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取ids(%v)录音数据失败: %v", ids, err)
	}
	return all
}

//获取指定时间内删除玩家的列表
func GetTeamLogDelPlayers(teamId, t int64) []int64 {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_TEAM_CHAT_LOG)
	defer closeFun()
	var list []*share_message.TeamChatLog
	q := bson.M{"Type": TALK_CONTENT_SYSTEM, "TeamId": teamId, "Time": bson.M{"$gte": t}, "TeamMessage.Type": bson.M{"$in": []int32{DEL_PLAYER, EXIT_PLAYER}}}
	err := col.Find(q).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		return nil
	}
	data := make([]int64, 0)
	for _, log := range list {
		message := log.GetTeamMessage()
		if message != nil {
			data = append(data, message.GetPlayerList()...)
		}
	}
	return data
}

//获取指定时间内成员职位变动信息
func GetTeamPlayerManagerChange(teamId, t int64) []*client_hall.TeamChangePos {
	col, closeFun := easygo.MongoLogMgr.GetC(MONGODB_NINGMENG_LOG, TABLE_TEAM_CHAT_LOG)
	defer closeFun()
	var list []*share_message.TeamChatLog
	err := col.Find(bson.M{"Type": TALK_CONTENT_SYSTEM, "TeamId": teamId, "Time": bson.M{"$gte": t}, "TeamMessage.Type": bson.M{"$in": []int32{ADD_MANAGER, DEL_MANAGER}}}).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		return nil
	}
	mIds := make(map[int64]int32)
	for _, log := range list {
		message := log.GetTeamMessage()
		if message != nil {
			for _, id := range message.GetPlayerList() {
				pos := GetTeamPlayerPos(teamId, id)
				mIds[id] = pos
			}
		}
	}
	data := make([]*client_hall.TeamChangePos, 0)
	for pid, p := range mIds {
		m := &client_hall.TeamChangePos{
			PlayerId: easygo.NewInt64(pid),
			Position: easygo.NewInt32(p),
		}
		data = append(data, m)
	}
	return data
}

//获取指定时间，群增加的成员
func GetTeamAddPlayers(teamId, t int64) []*share_message.PersonalTeamData {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAMMEMBER)
	defer closeFun()
	var list []*share_message.PersonalTeamData
	err := col.Find(bson.M{"TeamId": teamId, "Time": bson.M{"$gte": t}}).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error(err.Error())
	}
	return list
}

//检查新增点击报表人数
func OneClickLoginCheckAdd(id PLAYER_ID, btnType int32, filed string) {
	timeNow := easygo.GetToday0ClockTimestamp()
	log := &share_message.ButtonClickLog{
		Id:         easygo.NewString(easygo.AnytoA(timeNow) + easygo.AnytoA(id) + easygo.AnytoA(btnType)),
		CreateTime: easygo.NewInt64(timeNow),
		PlayerId:   easygo.NewInt64(id),
		Type:       easygo.NewInt32(btnType),
	}
	err := InsertMgo(MONGODB_NINGMENG_LOG, TABLE_BUTTON_CLICK_LOG, log) //如果插入失败,说明数据已经存在.
	if err == nil {
		SetRedisButtonClickReportFildVal(easygo.GetToday0ClockTimestamp(), 1, filed)
	}
}

//获取玩家之间的亲密度
func GetPlayerIntimacyInfoData(id string) *share_message.PlayerIntimacy {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_INTIMACY)
	defer closeFun()
	var data *share_message.PlayerIntimacy
	err := col.Find(bson.M{"_id": id}).One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error(err.Error())
	}
	return data
}

//获取自己合成语音作品数量
func GetCountMixVideo(pid int64) int32 {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_MIX_VOICE_VIDEO)
	defer closeFun()
	n, err := col.Find(bson.M{"PlayerId": pid, "Status": bson.M{"$ne": VC_STATUS_DELETE}}).Count()
	if err != nil {
		logs.Error("GetCountMixVideo err:", err.Error())
		return 0
	}
	return int32(n)
}

//录像合成写库
func InsertNewMixVideo(pid, bgId, mixTime int64, newUrl string, t int32) int64 {
	db := &share_message.PlayerMixVoiceVideo{
		Id:          easygo.NewInt64(NextId(TABLE_PLAYER_MIX_VOICE_VIDEO)),
		PlayerId:    easygo.NewInt64(pid),               //玩家id
		BgId:        easygo.NewInt64(bgId),              //背景音频ID
		MixVoiceUrl: easygo.NewString(newUrl),           //合成音频url
		CreateTime:  easygo.NewInt64(time.Now().Unix()), //合成时间
		PlayerType:  easygo.NewInt64(t),                 //玩家类型
		Status:      easygo.NewInt32(0),                 //状态：0-未审核,1-已发布,2-使用中,3-已删除
		MixTime:     easygo.NewInt64(mixTime),
	}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_MIX_VOICE_VIDEO)
	defer closeFun()
	err := col.Insert(db)
	easygo.PanicError(err)
	return db.GetId()
}

// 匹配数量常量
const (
	MATCH_LV_0 = 3 // 特级匹配对象数量
	MATCH_LV_1 = 5 // 一级匹配对象数量
	MATCH_LV_2 = 7 // 二级匹配对象数量
	MATCH_LV_3 = 4 // 三级匹配对象数量
	MATCH_LV_4 = 1 // 四级匹配对象数量
)

// 初始化星座匹配信息
var ConstellationMap map[int32]int32
var ConstellationNameMap map[int32]string

// 获取语音匹配
func GetVoiceCardToDB(uid int64, isFirst bool) []*client_hall.VoiceCard {
	st := GetMillSecond()
	result := make([]*client_hall.VoiceCard, 0)
	// 获取用户
	base := GetRedisPlayerBase(uid)
	if base == nil {
		return result
	}
	count := GetPlayerVoiceCardNum()
	operateIds := make([]int64, 0)
	//随机获取5个运营号
	if count < 500 {
		operateIds = GetRandomPlayerOperate(1, 5)
	} else if count >= 500 && count < 1000 {
		operateIds = GetRandomPlayerOperate(1, 2)
	}
	logs.Info("operateAccounts:", count, operateIds)
	mOperatePlayer := GetAllPlayerBase(operateIds)
	// 特级
	playerVoiceMap := make(map[int64]int64) // [玩家id]语音id
	premiumIds := make([]int64, 0)
	//选中玩家的个性化标签
	personalTags := make(map[int64][]int32, 0)
	premiumIds = append(premiumIds, uid) //过滤自己的id
	if len(operateIds) > 0 {
		premiumIds = append(premiumIds, operateIds...) //过滤运营号id
	}
	premium := PremiumPlayer(base, premiumIds, MATCH_LV_0) //去掉一个运营号
	for _, v := range premium {
		premiumIds = append(premiumIds, v.GetPlayerId())
	}
	if len(mOperatePlayer) > 0 {
		p := GetOneOperatePlayer(mOperatePlayer)
		if p != nil {
			premium = append(premium, p)
		}
	}
	// 排序
	result1 := SortByStart(base.GetConstellation(), premium)
	for _, v := range result1 {
		//playerVoiceMap[v.GetPlayerId()] = v.GetMixId()
		//tags:=make([]*client_hall.PersonTag,0)
		//for _,ta:=range v.GetPersonalityTags(){
		//	tag:=&client_hall.PersonTag{
		//		Id:easygo.NewInt32(ta),
		//		Name:easygo.NewString(GetPlayerPersonalityTags())
		//	}
		//}
		playerVoiceMap[v.GetPlayerId()] = v.GetMixId()
		personalTags[v.GetPlayerId()] = append(personalTags[v.GetPlayerId()], v.GetPersonalityTags()...)
		zanNum := v.GetVCZanNum() + v.GetBsVCZanNum()
		result = append(result, &client_hall.VoiceCard{
			PlayerId:      easygo.NewInt64(v.GetPlayerId()),
			NickName:      easygo.NewString(v.GetNickName()),
			HeadUrl:       easygo.NewString(v.GetHeadIcon()),
			Sex:           easygo.NewInt32(v.GetSex()),
			Constellation: easygo.NewInt32(v.GetConstellation()),
			ZanNum:        easygo.NewInt32(zanNum),
			VoiceUrl:      nil,
			IsOnLine:      easygo.NewBool(v.GetIsOnline()),
			BgUrl:         easygo.NewString(v.GetBgImageUrl()),
			//PersonalityTags: v.GetPersonalityTags(),
		})
	}
	//logs.Info("premiumIds--------->", premiumIds)
	// 一级
	oneNum := MATCH_LV_0 - len(premium) + MATCH_LV_1 + 1 // 一级用户需要的数量
	OnePlayer := MatchOnePlayer(base, premiumIds, oneNum)
	for _, v := range OnePlayer {
		premiumIds = append(premiumIds, v.GetPlayerId())
	}
	if len(mOperatePlayer) > 0 {
		p := GetOneOperatePlayer(mOperatePlayer)
		if p != nil {
			OnePlayer = append(OnePlayer, p)
		}
	}
	result2 := SortByStart(base.GetConstellation(), OnePlayer)
	for _, v := range result2 {
		playerVoiceMap[v.GetPlayerId()] = v.GetMixId()
		personalTags[v.GetPlayerId()] = append(personalTags[v.GetPlayerId()], v.GetPersonalityTags()...)
		zanNum := v.GetVCZanNum() + v.GetBsVCZanNum()
		result = append(result, &client_hall.VoiceCard{
			PlayerId:      easygo.NewInt64(v.GetPlayerId()),
			NickName:      easygo.NewString(v.GetNickName()),
			HeadUrl:       easygo.NewString(v.GetHeadIcon()),
			Sex:           easygo.NewInt32(v.GetSex()),
			Constellation: easygo.NewInt32(v.GetConstellation()),
			ZanNum:        easygo.NewInt32(zanNum),
			VoiceUrl:      nil,
			IsOnLine:      easygo.NewBool(v.GetIsOnline()),
			BgUrl:         easygo.NewString(v.GetBgImageUrl()),
		})
	}
	// 二级
	twoNum := oneNum - len(OnePlayer) + MATCH_LV_2 + 1
	twoPlayer := MatchTwoPlayer(base, premiumIds, twoNum)
	for _, v := range twoPlayer {
		premiumIds = append(premiumIds, v.GetPlayerId())
	}
	if len(mOperatePlayer) > 0 {
		p := GetOneOperatePlayer(mOperatePlayer)
		if p != nil {
			twoPlayer = append(twoPlayer, p)
		}
	}
	result3 := SortByStart(base.GetConstellation(), twoPlayer)
	for _, v := range result3 {
		playerVoiceMap[v.GetPlayerId()] = v.GetMixId()
		personalTags[v.GetPlayerId()] = append(personalTags[v.GetPlayerId()], v.GetPersonalityTags()...)
		zanNum := v.GetVCZanNum() + v.GetBsVCZanNum()
		result = append(result, &client_hall.VoiceCard{
			PlayerId:      easygo.NewInt64(v.GetPlayerId()),
			NickName:      easygo.NewString(v.GetNickName()),
			HeadUrl:       easygo.NewString(v.GetHeadIcon()),
			Sex:           easygo.NewInt32(v.GetSex()),
			Constellation: easygo.NewInt32(v.GetConstellation()),
			ZanNum:        easygo.NewInt32(zanNum),
			VoiceUrl:      nil,
			IsOnLine:      easygo.NewBool(v.GetIsOnline()),
			BgUrl:         easygo.NewString(v.GetBgImageUrl()),
		})
	}
	// 三级
	threeNum := twoNum - len(twoPlayer) + MATCH_LV_3 + 1
	threePlayer := MatchThreePlayer(base, premiumIds, threeNum)
	for _, v := range threePlayer {
		premiumIds = append(premiumIds, v.GetPlayerId())
	}
	if len(mOperatePlayer) > 0 {
		p := GetOneOperatePlayer(mOperatePlayer)
		if p != nil {
			threePlayer = append(threePlayer, p)
		}
	}
	result4 := SortByStart(base.GetConstellation(), threePlayer)
	for _, v := range result4 {
		playerVoiceMap[v.GetPlayerId()] = v.GetMixId()
		personalTags[v.GetPlayerId()] = append(personalTags[v.GetPlayerId()], v.GetPersonalityTags()...)
		zanNum := v.GetVCZanNum() + v.GetBsVCZanNum()
		result = append(result, &client_hall.VoiceCard{
			PlayerId:      easygo.NewInt64(v.GetPlayerId()),
			NickName:      easygo.NewString(v.GetNickName()),
			HeadUrl:       easygo.NewString(v.GetHeadIcon()),
			Sex:           easygo.NewInt32(v.GetSex()),
			Constellation: easygo.NewInt32(v.GetConstellation()),
			ZanNum:        easygo.NewInt32(zanNum),
			VoiceUrl:      nil,
			IsOnLine:      easygo.NewBool(v.GetIsOnline()),
			BgUrl:         easygo.NewString(v.GetBgImageUrl()),
		})
	}
	// 四级
	foreNum := threeNum - len(threePlayer) + MATCH_LV_4 + 1
	forePlayer := MatchForePlayer(base, premiumIds, foreNum)
	for _, v := range forePlayer {
		premiumIds = append(premiumIds, v.GetPlayerId())
	}
	if len(mOperatePlayer) > 0 {
		p := GetOneOperatePlayer(mOperatePlayer)
		if p != nil {
			forePlayer = append(forePlayer, p)
		}
	}
	result5 := SortByStart(base.GetConstellation(), forePlayer)
	for _, v := range result5 {
		playerVoiceMap[v.GetPlayerId()] = v.GetMixId()
		personalTags[v.GetPlayerId()] = append(personalTags[v.GetPlayerId()], v.GetPersonalityTags()...)
		zanNum := v.GetVCZanNum() + v.GetBsVCZanNum()
		result = append(result, &client_hall.VoiceCard{
			PlayerId:      easygo.NewInt64(v.GetPlayerId()),
			NickName:      easygo.NewString(v.GetNickName()),
			HeadUrl:       easygo.NewString(v.GetHeadIcon()),
			Sex:           easygo.NewInt32(v.GetSex()),
			Constellation: easygo.NewInt32(v.GetConstellation()),
			ZanNum:        easygo.NewInt32(zanNum),
			VoiceUrl:      nil,
			IsOnLine:      easygo.NewBool(v.GetIsOnline()),
			BgUrl:         easygo.NewString(v.GetBgImageUrl()),
		})
	}
	logs.Info("找到匹配人数:", len(result), base.GetPlayerId())
	voiceMap := GetPlayerVoice(playerVoiceMap)
	tagIds := make([]int32, 0)
	for _, ids := range personalTags {
		tagIds = append(tagIds, ids...)
	}
	allPersonalTags := GetPlayerPersonalityTagsMp(tagIds)

	for _, v := range result {
		comTags := base.GetCommonTags(v.GetPlayerId())
		v.VoiceUrl = easygo.NewString(voiceMap[v.GetPlayerId()])
		v.PersonalityTags = GetPersonalTagsForClient(personalTags[v.GetPlayerId()], allPersonalTags)
		v.MatchingDegree = easygo.NewInt32(base.GetMatchingDegree(v.GetPlayerId()))
		v.CommonTags = comTags
	}
	//广告
	advList := QueryAdvListToDB(ADV_LOCATION_VOICE_LOVE, base.GetLastLoginIP()) //恋爱匹配广告
	//	logs.Info("广告:", advList)
	if len(advList) == 0 {
		//无需插入广告
		ed := GetMillSecond()
		logs.Info("最终时间为====------=======", ed-st)
		return result
	}
	//需要插入广告
	pos := RandInt(0, len(result))
	lRe := len(result)
	if lRe < 10 {
		adv := GetRandomAdvData(advList, isFirst)
		if adv != nil {
			if isFirst && adv.GetAdv().GetIsTop() {
				result = append([]*client_hall.VoiceCard{adv}, result...)
			} else {
				result = append(result, adv)
			}
		}
	} else {
		//第一个广告
		pos = RandInt(7, 10)
		adv := GetRandomAdvData(advList, isFirst)
		if adv != nil {
			if isFirst && adv.GetAdv().GetIsTop() {
				result = append([]*client_hall.VoiceCard{adv}, result...)
			} else {
				result = easygo.Insert(result, pos, adv).([]*client_hall.VoiceCard)
			}
		}
		//第二个广告
		if lRe > 21 {
			if len(advList) > 0 {
				pos = RandInt(18, 21)
				adv = GetRandomAdvData(advList, false)
				if adv != nil {
					result = easygo.Insert(result, pos, adv).([]*client_hall.VoiceCard)
				}
			}
		}
	}
	ed := GetMillSecond()
	logs.Info("最终时间为====------=======", len(result), ed-st)
	return result
}

//通过权重获取随机获取广告
func GetRandomAdvData(list []*share_message.AdvSetting, isFirst bool) *client_hall.VoiceCard {
	//特殊处理，返回许愿池广告
	//sort.Slice(list, func(i, j int) bool {
	//		return list[i].GetWeights() > list[j].GetWeights() // 降序
	//})
	if isFirst {
		for _, l := range list {
			if l.GetIsTop() {
				adv := &client_hall.VoiceCard{
					Type: easygo.NewInt32(VC_CARD_TYPE_ADV),
					Adv:  l,
				}
				return adv
			}
		}
	}
	weights := make([]float32, 0)
	for _, l := range list {
		weights = append(weights, float32(l.GetWeights()))
	}
	n := WeightedRandomIndex(weights)
	data := list[n]
	list = append(list[:n], list[n+1:]...)
	if data != nil {
		adv := &client_hall.VoiceCard{
			Type: easygo.NewInt32(VC_CARD_TYPE_ADV),
			Adv:  data,
		}
		return adv
	}
	return nil
}
func GetPersonalTagsForClient(tags []int32, mp map[int32]*share_message.InterestTag) []*client_hall.PersonTag {
	data := make([]*client_hall.PersonTag, 0)
	for _, id := range tags {
		p := mp[id]
		if p != nil {
			t := &client_hall.PersonTag{
				Id:   easygo.NewInt32(id),
				Name: easygo.NewString(p.GetName()),
			}
			data = append(data, t)
		}
	}
	return data
}

// 找到对应的语音名片
func GetPlayerVoice(playerVoiceMap map[int64]int64) map[int64]string {
	voiceStringMap := make(map[int64]string)
	if len(playerVoiceMap) == 0 {
		return voiceStringMap
	}
	vids := make([]int64, 0)
	for _, v := range playerVoiceMap {
		vids = append(vids, v)
	}
	// 批量查询语音
	voiceList := GetVoiceCardInfoByIds(vids)
	for k, v := range playerVoiceMap {
		for _, v1 := range voiceList {
			if v == v1.GetId() {
				voiceStringMap[k] = v1.GetMixVoiceUrl()
			}
		}
	}
	return voiceStringMap
}

// 星座排序
func SortByStart(start int32, list []*share_message.PlayerBase) []*share_message.PlayerBase {
	result := make([]*share_message.PlayerBase, 0) // 结果
	if len(list) == 0 {
		return result
	}
	list1 := make([]*share_message.PlayerBase, 0) // 匹配
	list2 := make([]*share_message.PlayerBase, 0) // 不匹配
	for _, v := range list {
		if ConstellationMap[start] == v.GetConstellation() { // 星座匹配
			list1 = append(list1, v)
		} else {
			list2 = append(list2, v)
		}

	}

	if len(list1) > 0 {
		result = append(result, list1...)
	}
	if len(list2) > 0 {
		result = append(result, list2...)
	}
	return result
}

// 特级用户 异性,双方在线,共同兴趣>=2,共同标签>=2
func PremiumPlayer(base *RedisPlayerBaseObj, ids []int64, num int) []*share_message.PlayerBase {
	st := GetMillSecond()
	label := base.GetLabelList()
	personalityTags := base.GetPersonalityTags()
	// 特级 双方在线,共同兴趣
	m := []bson.M{
		{"$match": bson.M{"Label": bson.M{"$exists": true}, "PersonalityTags": bson.M{"$exists": true}}}, // todo 测试注释
		{"$project": bson.M{
			"PersonalityTags": 1, "Label": 1, "HeadIcon": 1, "NickName": 1, "IsOnline": 1, "Sex": 1, "VCZanNum": 1, "BsVCZanNum": 1, "Constellation": 1,
			"aaa": bson.M{"$setIntersection": []interface{}{"$PersonalityTags", personalityTags}}, // todo 测试注释
			"bbb": bson.M{"$setIntersection": []interface{}{"$Label", label}},
		}},

		{"$project": bson.M{"PersonalityTags": 1, "Label": 1, "HeadIcon": 1, "NickName": 1, "IsOnline": 1, "Sex": 1, "VCZanNum": 1, "BsVCZanNum": 1, "Constellation": 1, "count1": bson.M{"$size": "$aaa"}, "count2": bson.M{"$size": "$bbb"}}}, // todo 测试注释
		{"$match": bson.M{"count1": bson.M{"$gte": 2}, "count2": bson.M{"$gte": 2}, "MixId": bson.M{"$gt": 0}}},                                                                                                                                 // todo 测试注释
		{"$match": bson.M{"IsOnline": true, "_id": bson.M{"$nin": ids}, "Sex": bson.M{"$ne": base.GetSex()}}},
		//{"$limit": num}, // todo 样例
		{"$sample": bson.M{"size": num}},
	}

	ls := FindPipeAll(MONGODB_NINGMENG, TABLE_PLAYER_BASE, m, 0, 0)

	var list []*share_message.PlayerBase
	for _, li := range ls {
		one := &share_message.PlayerBase{}
		StructToOtherStruct(li, one)

		list = append(list, one)
	}
	ed := GetMillSecond()
	logs.Info("特级时间------->%d,人数:%d", ed-st, len(ls))
	return list

}

// 一级用户 双方在线、共同个性标签>=2
func MatchOnePlayer(base *RedisPlayerBaseObj, ids []int64, num int) []*share_message.PlayerBase {
	st := GetMillSecond()
	personalityTags := base.GetPersonalityTags()
	// 特级 双方在线,共同兴趣
	m := []bson.M{
		{"$match": bson.M{"PersonalityTags": bson.M{"$exists": true}}}, // todo 测试注释
		{"$project": bson.M{
			"PersonalityTags": 1, "Label": 1, "HeadIcon": 1, "NickName": 1, "IsOnline": 1, "Sex": 1, "VCZanNum": 1, "BsVCZanNum": 1, "Constellation": 1,
			"aaa": bson.M{"$setIntersection": []interface{}{"$PersonalityTags", personalityTags}}, // todo 测试注释
		}},

		{"$project": bson.M{"PersonalityTags": 1, "Label": 1, "HeadIcon": 1, "NickName": 1, "IsOnline": 1, "Sex": 1, "VCZanNum": 1, "BsVCZanNum": 1, "Constellation": 1, "count1": bson.M{"$size": "$aaa"}}}, // todo 测试注释
		{"$match": bson.M{"count1": bson.M{"$gte": 2}, "IsOnline": true, "Sex": bson.M{"$ne": base.GetSex()}, "MixId": bson.M{"$gt": 0}}},                                                                    // todo 测试注释
		//{"$limit": num}, // todo 样例
		{"$sample": bson.M{"size": num}},
	}
	if len(ids) > 0 {
		matchBson := bson.M{"$match": bson.M{"_id": bson.M{"$nin": ids}}}
		m = append(m, matchBson)
	}
	ls := FindPipeAll(MONGODB_NINGMENG, TABLE_PLAYER_BASE, m, 0, 0)

	var list []*share_message.PlayerBase
	for _, li := range ls {
		one := &share_message.PlayerBase{}
		StructToOtherStruct(li, one)

		list = append(list, one)
	}
	ed := GetMillSecond()
	logs.Info("1级时间------->%d,人数:%d", ed-st, len(ls))
	return list

}

// 二级用户 共同个性标签>=2
func MatchTwoPlayer(base *RedisPlayerBaseObj, ids []int64, num int) []*share_message.PlayerBase {
	st := GetMillSecond()
	personalityTags := base.GetPersonalityTags()
	//
	m := []bson.M{
		{"$match": bson.M{"PersonalityTags": bson.M{"$exists": true}}}, // todo 测试注释
		{"$project": bson.M{
			"PersonalityTags": 1, "Label": 1, "HeadIcon": 1, "NickName": 1, "IsOnline": 1, "Sex": 1, "VCZanNum": 1, "BsVCZanNum": 1, "Constellation": 1,
			"aaa": bson.M{"$setIntersection": []interface{}{"$PersonalityTags", personalityTags}}, // todo 测试注释
		}},

		{"$project": bson.M{"PersonalityTags": 1, "Label": 1, "HeadIcon": 1, "NickName": 1, "IsOnline": 1, "Sex": 1, "VCZanNum": 1, "BsVCZanNum": 1, "Constellation": 1, "count1": bson.M{"$size": "$aaa"}}}, // todo 测试注释
		{"$match": bson.M{"count1": bson.M{"$gte": 2}, "MixId": bson.M{"$gt": 0}}}, // todo 测试注释
		//{"$limit": num}, // todo 样例
		{"$sample": bson.M{"size": num}},
	}
	if len(ids) > 0 {
		matchBson := bson.M{"$match": bson.M{"_id": bson.M{"$nin": ids}}}
		m = append(m, matchBson)
	}
	ls := FindPipeAll(MONGODB_NINGMENG, TABLE_PLAYER_BASE, m, 0, 0)

	var list []*share_message.PlayerBase
	for _, li := range ls {
		one := &share_message.PlayerBase{}
		StructToOtherStruct(li, one)

		list = append(list, one)
	}
	ed := GetMillSecond()
	logs.Info("2级时间------->%d,人数:%d", ed-st, len(ls))
	return list
}

// 三级用户 共同兴趣>=1  或  共同个性标签>=1  或  星座匹配
func MatchThreePlayer(base *RedisPlayerBaseObj, ids []int64, num int) []*share_message.PlayerBase {
	st := GetMillSecond()
	personalityTags := base.GetPersonalityTags()
	label := base.GetLabelList()

	//match := bson.M{"$or": []bson.M{{"aaa.0": bson.M{"$exists": 1}}, {"bbb.0": bson.M{"$exists": 1}}, {"Constellation": base.GetConstellation()}}, "_id": bson.M{"$nin": ids}}
	m := []bson.M{

		{"$match": bson.M{"Label": bson.M{"$exists": true}, "PersonalityTags": bson.M{"$exists": true}}}, // todo 测试注释
		{"$project": bson.M{
			"PersonalityTags": 1, "Label": 1, "HeadIcon": 1, "NickName": 1, "IsOnline": 1, "Sex": 1, "VCZanNum": 1, "BsVCZanNum": 1, "Constellation": 1,
			"aaa": bson.M{"$setIntersection": []interface{}{"$PersonalityTags", personalityTags}}, // todo 测试注释
			"bbb": bson.M{"$setIntersection": []interface{}{"$Label", label}},
		}},

		{"$project": bson.M{"PersonalityTags": 1, "Label": 1, "HeadIcon": 1, "NickName": 1, "IsOnline": 1, "Sex": 1, "VCZanNum": 1, "BsVCZanNum": 1, "Constellation": 1, "count1": bson.M{"$size": "$aaa"}, "count2": bson.M{"$size": "$bbb"}}}, // todo 测试注释
		{"$match": bson.M{"$or": []bson.M{{"count1": bson.M{"$gte": 1}}, {"count2": bson.M{"$gte": 1}}}}},
		{"$match": bson.M{"MixId": bson.M{"$gt": 0}}},
		{"$sample": bson.M{"size": num}},
	}
	if len(ids) > 0 {
		matchBson := bson.M{"$match": bson.M{"_id": bson.M{"$nin": ids}}}
		m = append(m, matchBson)
	}
	ls := FindPipeAll(MONGODB_NINGMENG, TABLE_PLAYER_BASE, m, 0, 0)

	var list []*share_message.PlayerBase
	for _, li := range ls {
		one := &share_message.PlayerBase{}
		StructToOtherStruct(li, one)

		list = append(list, one)
	}
	ed := GetMillSecond()
	logs.Info("3级时间------->%d,人数:%d", ed-st, len(ls))
	return list

}

func MatchForePlayer(base *RedisPlayerBaseObj, ids []int64, num int) []*share_message.PlayerBase {
	st := GetMillSecond()
	m := []bson.M{
		{"$match": bson.M{"Sex": bson.M{"$ne": base.GetSex()}, "MixId": bson.M{"$gt": 0}}}, // todo 测试注释
		{"$sample": bson.M{"size": 3000}},
		//{"$match": bson.M{"Constellation": bson.M{"$ne": ConstellationMap[base.GetConstellation()]}}}, // todo 测试注释

	}
	if len(ids) > 0 {
		matchBson := bson.M{"$match": bson.M{"_id": bson.M{"$nin": ids}}}
		m = append(m, matchBson)
	}
	mm := bson.M{"$limit": num}
	m = append(m, mm)
	ls := FindPipeAll(MONGODB_NINGMENG, TABLE_PLAYER_BASE, m, 0, 0)

	var list []*share_message.PlayerBase
	for _, li := range ls {
		one := &share_message.PlayerBase{}
		StructToOtherStruct(li, one)

		list = append(list, one)
	}
	ed := GetMillSecond()
	logs.Info("4级时间------->%d,人数:%d", ed-st, len(ls))
	return list
}

/* func MatchForePlayer(base *RedisPlayerBaseObj, ids []int64, num int) []*share_message.PlayerBase {
	st := GetMillSecond()
	personalityTags := base.GetPersonalityTags()
	label := base.GetLabelList()
	//match :=   bson.M {"$or":[{"aaa.0": {"$exists": 0}},{"bbb.0":  {"$exists": 0}}],"Constellation":{"$ne":1}} ,
	match := bson.M{"$or": []bson.M{{"aaa.0": bson.M{"$exists": 0}}, {"bbb.0": bson.M{"$exists": 0}}, {"Constellation": bson.M{"$ne": base.GetConstellation()}}}, "_id": bson.M{"$nin": ids}}
	list := MatchPlayer(label, personalityTags, match, num)
	ed := GetMillSecond()
	logs.Info("4级时间------->", ed-st)
	return list
}
*/
// 级别用户
func MatchPlayer(label, personalityTags []int32, match bson.M, num int) []*share_message.PlayerBase {
	m := []bson.M{
		{"$project": bson.M{
			"PersonalityTags": 1, "Label": 1, "HeadIcon": 1, "NickName": 1, "IsOnline": 1, "Sex": 1, "VCZanNum": 1, "BsVCZanNum": 1, "Constellation": 1,
			"aaa": bson.M{"$setIntersection": []interface{}{"$PersonalityTags", personalityTags}},
			"bbb": bson.M{"$setIntersection": []interface{}{"$Label", label}},
		}},

		{"$match": match},
		//{"$limit": num}, // todo 样例
		{"$sample": bson.M{"size": num}},
	}

	limit := bson.M{"$limit": num}
	m = append(m, limit)
	ls := FindPipeAll(MONGODB_NINGMENG, TABLE_PLAYER_BASE, m, 0, 0)

	var list []*share_message.PlayerBase
	for _, li := range ls {
		one := &share_message.PlayerBase{}
		StructToOtherStruct(li, one)
		list = append(list, one)
	}
	return list
}

//随机切换一张背景图
func GetRandomBgImageUrl(url string) *share_message.SystemBgImage {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SYSTEM_BG_IMAGE)
	defer closeFun()
	m := []bson.M{
		{"$match": bson.M{"_id": bson.M{"$ne": url}}},
		{"$sample": bson.M{"size": 1}},
	}
	query := col.Pipe(m)
	var one *share_message.SystemBgImage
	err := query.One(&one)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return one
}

//获取亲密度配置
func GetConfigIntimacyFormDB() []*share_message.IntimacyConfig {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_INTIMACY_COINFIG)
	defer closeFun()
	var data []*share_message.IntimacyConfig
	err := col.Find(bson.M{}).All(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error(err.Error())
	}
	return data
}

//获取最熱、最新的标签
func GetVoiceTags(reqType int32, hotNum, laterNum int) (topTags, laterTags []*share_message.InterestTag) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_BG_VOICE_TAG)
	defer closeFun()

	err := col.Find(bson.M{"InterestType": reqType, "Status": 0}).Sort("-Count").Limit(hotNum).All(&topTags)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取最热的标签: %v", err.Error())
		return nil, nil
	}
	topIds := make([]int32, 0)
	for _, v := range topTags {
		topIds = append(topIds, v.GetId())
	}
	//获取所有标签
	if laterNum == -1 {
		err = col.Find(bson.M{"InterestType": reqType, "_id": bson.M{"$nin": topIds}, "Status": 0}).Sort("-UpdateTime").All(&laterTags)
	} else {
		//err = col.Find(bson.M{"InterestType": reqType, "_id": bson.M{"$nin": topIds}}).Sort("-UpdateTime").Limit(laterNum).All(&laterTags)
		m := []bson.M{
			{"$match": bson.M{"InterestType": reqType, "_id": bson.M{"$nin": topIds}, "Status": 0}},
			{"$sample": bson.M{"size": laterNum}},
		}
		query := col.Pipe(m)
		err = query.All(&laterTags)
	}
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取最新的标签: %v", err.Error())
		return topTags, nil
	}
	return topTags, laterTags
}

//获取标签下的背景音频素材
func GetTagVoiceVideo(pid int64, tagId int32, page, pageSize int) []*share_message.BgVoiceVideo {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_BG_VOICE_VIDEO)
	defer closeFun()
	curPage := easygo.If(page > 1, page-1, 0).(int)
	var voices []*share_message.BgVoiceVideo
	q1 := bson.M{"Status": VC_VOIDE_PASS_AUDIT, "Tags": tagId}
	q2 := bson.M{"Status": bson.M{"$lte": VC_VOIDE_PASS_AUDIT}, "PlayerId": pid}
	q := bson.M{"$or": []bson.M{q1, q2}}
	err := col.Find(q).Sort("-UseCount").Skip(curPage * pageSize).Limit(pageSize).All(&voices)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取指定标签音频素材失败：", err)
	}
	return voices
}

//增加背景音频记录
func AddBgVoiceVideo(log *share_message.BgVoiceVideo) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_BG_VOICE_VIDEO)
	defer closeFun()
	err := col.Insert(log)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error(err.Error())
	}
}

//获取搜索的背景音乐
func GetBgVoice(pid int64, reqMsg *client_hall.SearchVoiceVideoReq) []*share_message.BgVoiceVideo {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_BG_VOICE_VIDEO)
	defer closeFun()
	reqType := reqMsg.GetType()
	content := reqMsg.GetContent()
	page := easygo.If(int(reqMsg.GetPage()) > 1, int(reqMsg.GetPage())-1, 0).(int)
	pageSize := 10                                          //默认一页10条
	q1 := bson.M{"Status": VC_VOIDE_PASS_AUDIT}             //获取通过审核的
	q2 := bson.M{"Status": VC_VOIDE_AUDIT, "PlayerId": pid} //获取未通过审核的，但是自己的
	queryBson := bson.M{"$or": []bson.M{q1, q2}}
	//queryBson := bson.M{"Status": VC_VOIDE_PASS_AUDIT} //获取通过审核的
	if reqType == VC_VOIDE_NAME {
		queryBson["Name"] = bson.M{"$regex": bson.RegEx{Pattern: content, Options: "im"}}
	} else if reqType == VC_VOIDE_CONTENT {
		queryBson["Content"] = bson.M{"$regex": bson.RegEx{Pattern: content, Options: "im"}}
	} else if reqType == VC_VOIDE_AUTHOR {
		queryBson["Maker"] = bson.M{"$regex": bson.RegEx{Pattern: content, Options: "im"}}
	} else {
		logs.Error("搜索背景音频素材类型有误")
		return nil
	}
	var voices []*share_message.BgVoiceVideo
	err := col.Find(queryBson).Sort("-UseCount").Skip(page * pageSize).Limit(pageSize).All(&voices)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("搜索背景音频素材类型失败：", err)
	}
	return voices
}

//获取指定录像背景
func GetBgVoiceVideo(pid int64, reqMsg *client_hall.GetVoiceVideoReq) []*share_message.BgVoiceVideo {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_BG_VOICE_VIDEO)
	defer closeFun()
	var voices []*share_message.BgVoiceVideo
	page := easygo.If(int(reqMsg.GetPage()) > 1, int(reqMsg.GetPage())-1, 0).(int)
	pageSize := int(reqMsg.GetPageSize()) //默认一页10条
	if pageSize == 0 {
		pageSize = 10
	}
	q := bson.M{"$or": []bson.M{{"Status": VC_VOIDE_PASS_AUDIT}, {"Status": VC_VOIDE_AUDIT, "PlayerId": pid}}}
	query := bson.M{"$and": []bson.M{bson.M{"Type": reqMsg.GetType()}, q}}
	err := col.Find(query).Sort("_id").Skip(page * pageSize).Limit(pageSize).All(&voices)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("搜索背景素材类型失败：", err)
	}
	return voices
}

//获取热门片段
func GetHotBgVoice(reqType int32, topNum int) []*share_message.BgVoiceVideo {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_BG_VOICE_VIDEO)
	defer closeFun()
	var topVoice []*share_message.BgVoiceVideo

	err := col.Find(bson.M{"Type": reqType, "Status": VC_VOIDE_PASS_AUDIT}).Sort("-UseCount").Limit(topNum).All(&topVoice)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取最热背景音频素材失败：", err)
	}

	return topVoice
}

//随机获取素材作品
func GetRandBgVoice(reqType int32, num int) []*share_message.BgVoiceVideo {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_BG_VOICE_VIDEO)
	defer closeFun()
	var voice []*share_message.BgVoiceVideo

	m := []bson.M{
		bson.M{"$match": bson.M{"Status": VC_VOIDE_PASS_AUDIT, "Type": reqType}},
		bson.M{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&voice)
	if err != nil {
		logs.Error("随机获取音乐作品素材失败")
		return nil
	}
	return voice
}

//获取超过3天加亲密度的数据
func GetExpirationPlayerIntimacy() []*share_message.PlayerIntimacy {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_INTIMACY)
	defer closeFun()
	t := time.Now().Unix() - PLAYER_INTIMACY_REDUCE_DAY
	var data []*share_message.PlayerIntimacy
	q1 := bson.M{"LastTime": bson.M{"$lte": t}}
	q2 := bson.M{"$or": []bson.M{{"IntimacyVal": bson.M{"$gt": 0}}, {"IntimacyLv": bson.M{"$gt": 0}}}}
	q := bson.M{"$and": []bson.M{q1, q2}}
	err := col.Find(q).All(&data)
	easygo.PanicError(err)
	return data
}

//充值玩家亲密度数据
func CleanPlayerIntimacy(id string) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_INTIMACY)
	defer closeFun()
	err := col.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"IntimacyVal": 0, "IntimacyLv": 0, "LastTime": 0}})
	easygo.PanicError(err)
}

//更新指定标签引用次数
func IncBgVoiceTag(ids []int32) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_BG_VOICE_TAG)
	defer closeFun()
	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$inc": bson.M{"Count": 1}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("无法更新标签引用次数")
	}
}

//
func GetSomeNoticePlayers(pid int64, findBson bson.M) []*share_message.PlayerBase {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	datas := make([]*share_message.PlayerBase, 0)
	queryBson := bson.M{}
	if pid > 0 {
		queryBson["_id"] = bson.M{"$gt": pid}
	}
	query := bson.M{"$and": []bson.M{queryBson, findBson}}
	err := col.Find(query).Select(bson.M{"_id": 1}).Sort("_id").Limit(10000).Sort("_id").All(&datas)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return datas
}

//获取指定类型，指定数量的运营号
func GetRandomPlayerOperate(t, num int32) []int64 {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_OPERATE)
	defer closeFun()
	m := []bson.M{
		{"$match": bson.M{"Type": t}},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	var list []*share_message.PlayerOperate
	err := query.All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetRandomPlayerOperate err:", err)
	}
	var ids []int64
	for _, d := range list {
		ids = append(ids, d.GetPlayerId())
	}
	return ids
}

//从已获取的运营号中抽取一个
func GetOneOperatePlayer(m map[PLAYER_ID]*share_message.PlayerBase) *share_message.PlayerBase {
	if len(m) == 0 {
		return nil
	}
	for k, p := range m {
		delete(m, k)
		return p
	}
	return nil
}

//获取设置了语音名片的玩家数量
func GetPlayerVoiceCardNum() int {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_BASE)
	defer closeFun()
	n, err := col.Find(bson.M{"MixId": bson.M{"$gt": 0}}).Count()
	if err != nil {
		logs.Error("GetPlayerVoiceCardNum err:", err.Error())
		return 0
	}
	return n
}

//获取指定充值活动id信息
func GetCoinRechargeActCfg(id int64) *share_message.WishCoinRechargeActivityCfg {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_COIN_RECHARGE_ACT_CFG)
	defer closeFun()
	var item *share_message.WishCoinRechargeActivityCfg
	err := col.Find(bson.M{"_id": id}).One(&item)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return item
}

//获取话题群组信息
func GetTeamsByTopic(topic string) []*share_message.TeamData {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAM_DATA)
	defer closeFun()
	var teams []*share_message.TeamData
	err := col.Find(bson.M{"Topic": topic, "Status": NORMAL}).Sort("-CreateTime").All(&teams)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetTeamsByTopic err:", err)
	}
	return teams
}

//获取当前话题群组数量
func GetTopicTeamNumByTopic(topic string) int {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAM_DATA)
	defer closeFun()
	n, err := col.Find(bson.M{"Topic": topic, "Status": NORMAL}).Count()
	if err != nil {
		logs.Error("GetTeamsByTopic err:", err)
	}
	return n
}

//获取10分钟内没发过动态的群组
func GetAllUnSendDynamicTopicTeam() []*share_message.TeamData {
	var teams []*share_message.TeamData
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TEAM_DATA)
	defer closeFun()
	t := time.Now().Unix()
	q1 := bson.M{"Topic": bson.M{"$ne": nil}}
	q2 := bson.M{"Topic": bson.M{"$ne": ""}}
	q3 := bson.M{"$or": []bson.M{bson.M{"LastDynamicTime": bson.M{"$lt": t - 600}}, bson.M{"LastDynamicTime": bson.M{"$eq": nil}}}}
	q := bson.M{"$and": []bson.M{q1, q2, q3}}
	err := col.Find(q).All(&teams)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetTeamsByTopic err:", err)
	}
	return teams
}

//获取主页菜单项数据返回给前端
func GetIndexTipsToClient() *client_hall.AllMainMenu {
	allMenu := make([]*client_hall.MainMenu, 0)
	var data []*share_message.IndexTips
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_INDEX_TIPS)
	defer closeFun()
	err := col.Find(bson.M{"Status": 1}).All(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetIndexTipsToClient err:", err.Error())
	}
	advIds := make([]int64, 0)
	mAdvIds := make(map[int32]int64, 0) //菜单id映射广告id
	for _, d := range data {
		menu := &client_hall.MainMenu{
			Id:     easygo.NewInt32(d.GetId()),
			Type:   easygo.NewInt32(d.GetTypes()),
			IsOpen: easygo.NewBool(true),
			Slot:   easygo.NewInt32(d.GetWeights()),
		}
		if d.GetAdvId() != 0 {
			advIds = append(advIds, d.GetAdvId())
		}
		mAdvIds[d.GetId()] = d.GetAdvId()
		allMenu = append(allMenu, menu)
	}
	mAdvs := GetAllAdvsByIds(advIds)
	newallMenu := make([]*client_hall.MainMenu, 0)
	for _, menu := range allMenu {
		id := mAdvIds[menu.GetId()]
		if id != 0 && mAdvs[id] == nil {
			continue
		}
		menu.Adv = mAdvs[id]
		newallMenu = append(newallMenu, menu)
	}
	msg := &client_hall.AllMainMenu{
		Menus: newallMenu,
	}
	return msg
}

//获取弹窗广告
func GetAllTipAdvsFromDB() *client_hall.AllTipAdv {
	allMenu := make([]*client_hall.TipAdv, 0)
	var data []*share_message.PopSuspend
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_POP_SUSPEND)
	defer closeFun()
	err := col.Find(bson.M{}).All(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetIndexTipsToClient err:", err.Error())
	}
	advIds := make([]int64, 0)
	mAdvIds := make(map[int32]int64, 0) //
	for _, d := range data {
		menu := &client_hall.TipAdv{
			Page:         easygo.NewInt32(d.GetId()),
			IsJump:       easygo.NewBool(d.GetIsPop()),
			Issuspension: easygo.NewBool(d.GetIsSuspend()),
		}
		advIds = append(advIds, d.GetAdvId())
		mAdvIds[d.GetId()] = d.GetAdvId()
		allMenu = append(allMenu, menu)
	}
	mAdvs := GetAllAdvsByIds(advIds)
	newallMenu := make([]*client_hall.TipAdv, 0)
	for _, menu := range allMenu {
		id := mAdvIds[menu.GetPage()]
		if id != 0 && mAdvs[id] == nil {
			continue
		}
		menu.Adv = mAdvs[id]
		newallMenu = append(newallMenu, menu)
	}
	msg := &client_hall.AllTipAdv{
		Advs: newallMenu,
	}
	return msg
}
