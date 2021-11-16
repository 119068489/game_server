package for_game

import (
	"fmt"
	"game_server/easygo"
	"game_server/pb/share_message"
	"log"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
)

var _ = fmt.Sprintf
var _ = log.Println
var _ = easygo.Underline

//================================================================================支付设置
//查询通用额度配置
func GetGeneralQuota() *share_message.GeneralQuota {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_GENERAL_QUOTA)
	defer closeFun()
	var obj *share_message.GeneralQuota
	err := col.Find(bson.M{}).One(&obj)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return obj
}

//查询支付类型
func GetPayType() []*share_message.PayType {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PAYTYPE)
	defer closeFun()

	var list []*share_message.PayType
	queryBson := bson.M{}

	q := col.Find(queryBson)
	err := q.All(&list)
	easygo.PanicError(err)

	return list
}

//id查询支付源类型
func GettSourceTypeById(id int32) *share_message.SourceType {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SOURCETYPE)
	defer closeFun()

	siteOne := &share_message.SourceType{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//查询当前可用的支付、代付通道配置
func QueryPlatformChannelList(t int32) []*share_message.PlatformChannel {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLATFORM_CHANNEL)
	defer closeFun()

	queryBson := bson.M{}
	queryBson["Status"] = t

	var list []*share_message.PlatformChannel
	query := col.Find(queryBson)
	errc := query.Sort("-_id").All(&list)
	easygo.PanicError(errc)

	return list
}

//查询当前可用的支付设定配置
func QueryPaymentSettingList() []*share_message.PaymentSetting {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PAYMENTSETTING)
	defer closeFun()

	queryBson := bson.M{}
	var list []*share_message.PaymentSetting
	query := col.Find(queryBson)
	errc := query.Sort("-_id").All(&list)
	easygo.PanicError(errc)

	return list
}

//id查询支付平台配置
func QuerPaymentPlatformById(id int32) *share_message.PaymentPlatform {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PAYMENTPLATFORM)
	defer closeFun()

	siteOne := &share_message.PaymentPlatform{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//id查询支付设定
func QueryPaymentSettingById(id int32) *share_message.PaymentSetting {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PAYMENTSETTING)
	defer closeFun()

	siteOne := &share_message.PaymentSetting{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//id查询支付场景
func QueryPaySceneById(id int32) *share_message.PayScene {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PAYSCENE)
	defer closeFun()

	siteOne := &share_message.PayScene{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//id查询支付方式
func QueryPayTypeById(id int32) *share_message.PayType {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PAYTYPE)
	defer closeFun()

	siteOne := &share_message.PayType{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}
