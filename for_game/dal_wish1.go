package for_game

import (
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/brower_backstage"
	"game_server/pb/h5_wish"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo"

	//"gopkg.in/mgo.v2"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

//兑换状态
const (
	WISH_TO_EXCHANGE   = 0 // 待兑换
	WISH_EXCHANGED     = 1 // 已兑换
	WISH_RECYCLED      = 2 //已回收
	WISH_RECYCLE_CHECK = 3 //回收审核中
	WISH_TO_CONTROL    = 4 //待兑换与回收审核中
)

const (
	RECYCLE_TO_DIAMOND = 0 //回收到钻石
	RECYCLE_TO_CARD    = 1 //回收到银行卡
)

const (
	RECYCLE_BY_PLAYER = 1 //用户回收
	RECYCLE_BY_SYSTEM = 2 //系统回收
	RECYCLE_BY_CHECK  = 3 //人工审核
)

const (
	WISH_SALEOUT = 0 // 下架
	WISH_ONSALE  = 1 // 上架
)

const (
	PAGE       = 1
	PAGE_LIMIT = 20
)

const (
	WISH_OUT_MONEY_TIME   = "WishOutMoneyTime"   //单日出款次数
	WISH_OUT_MONEY_SUM    = "WishOutMoneySum"    //单日出款金额
	WISH_OUT_DIAMOND_TIME = "WishOutDiamondTime" //单日钻石出款次数
	WISH_OUT_DIAMOND_SUM  = "WishOutDiamondSum"  //单日出款钻石总数
)

//玩家物品操作
const (
	DEAL_SUCCESS      = 0 //成功
	DEAL_FAULT        = 1 //一般失败
	DEAL_PLAYER_LIMIT = 2 //白名单，运营好限制
)

//许愿池埋点
const (
	WISH_REPORT_ACCESS_WISH     = 1 //访问许愿池
	WISH_REPORT_ACCESS_EXCHANGE = 2 //访问兑换钻石页
	WISH_REPORT_VEXCHANGE       = 3 //成功兑换钻石
	WISH_REPORT_ACCESS_DERE     = 4 //访问挑战页
	WISH_REPORT_CHALLENGE       = 5 //点击挑战盲盒人数
)

//回收处理状态
const (
	RECYCLE_STATUS_ERROR   = 1 //失败
	RECYCLE_STATUS_SUCCESS = 2 //成功
	RECYCLE_STATUS_CHECK   = 3 //人工审核
	RECYCLE_STATUS_FREEZE  = 4 //冻结
)

//兑换处理状态
const (
	EXCHANGE_STATUS_ERROR     = 1 //一般失败
	EXCHANGE_STATUS_SUCCESS   = 2 //成功
	EXCHANGE_STATUS_CHECK     = 3 //邮费不足
	EXCHANGE_STATUS_FREEZE    = 4 //冻结
	EXCHANGE_STATUS_ISPRESALE = 5 //全为预售商品
)

//设置已试玩
func SetFirstPlay(pid int64) error {
	col, closeFunItem := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFunItem()
	err := col.Update(bson.M{"_id": pid}, bson.M{"$set": bson.M{"IsTryOne": true}})
	if err != nil {
		logs.Error("设置玩家(%v)已试玩err: %v", pid, err.Error())
	}
	return err
}

//根据Id查询PlayerWishItem
func GetPlayerWishItemByIds(pid int64, ids []int64) ([]*share_message.PlayerWishItem, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_ITEM)
	defer closeFun()
	data := make([]*share_message.PlayerWishItem, 0)
	//校验是否已兑换回收
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}, "Status": 0, "PlayerId": pid}).Sort("-CreateTime").All(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("根据Id:%v查询PlayerWishItem: %v", ids, err.Error())
	}
	return data, err
}
func GetWishPlayerByIds(ids []int64) ([]*share_message.WishPlayer, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFun()
	data := make([]*share_message.WishPlayer, 0)
	//校验是否已兑换回收
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("根据Id:%v查询GetWishPlayerByIds: %v", ids, err.Error())
	}
	return data, err
}

// 商品ids找到预售商品信息
func GetPreSaleProductByIds(ids []int64) []*share_message.WishItem {
	var lst []*share_message.WishItem
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ITEM)
	defer closeFunc()
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}, "IsPreSale": true}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("商品ids(%v)查询预售物品失败 err: %s", ids, err.Error())
	}
	return lst
}

//根据用户Id查询许愿池
func GetWishPlayerByPlayerId(playerId int64) (*share_message.WishPlayer, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFun()
	var data *share_message.WishPlayer
	err := col.Find(bson.M{"_id": playerId}).One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("根据用户Id(%v)查询许愿池err: %v", playerId, err.Error())
	}
	return data, err
}

//新增收货地址
func InsertAddressByUid(uid int64, address *h5_wish.WishAddress) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFun()
	err := col.Update(bson.M{"_id": uid}, bson.M{"$push": bson.M{"Address": address}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("新增玩家(%v)的收货地址: %v", uid, err.Error())
	}
	return err
}

//清空用户默认地址状态
func ClearUserDefaultAddress(uid int64) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFun()

	count, err := col.Find(bson.M{"_id": uid, "Address.IfDefault": true}).Count()
	if err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			return nil
		}
		logs.Error("查询用户添加的地址: player(%v), err: %v", uid, err.Error())
		return err
	}
	if count > 0 {
		err = col.Update(bson.M{"_id": uid}, bson.M{"$set": bson.M{"Address.$[].IfDefault": false}})
		if err != nil {
			logs.Error("清空用户默认地址状态: %v", err.Error())
		}
	}
	return err
}

//编辑收货地址
func EditAddressByUid(uid int64, address *h5_wish.WishAddress) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFun()
	//若修改新地址为默认地址，则将默认地址清空
	if address.GetIfDefault() == true {
		err := ClearUserDefaultAddress(uid)
		if err != nil {
			return err
		}
	}
	setBson := bson.M{
		"Address.$.Name":      address.GetName(),
		"Address.$.Phone":     address.GetPhone(),
		"Address.$.Detail":    address.GetDetail(),
		"Address.$.IfDefault": address.GetIfDefault(),
		"Address.$.Province":  address.GetProvince(),
		"Address.$.City":      address.GetCity(),
		"Address.$.Area":      address.GetArea(),
	}
	err := col.Update(bson.M{"_id": uid, "Address.AddressId": address.GetAddressId()}, bson.M{"$set": setBson})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("编辑收货地址: %v", err.Error())
	}
	return err
}

//删除收货地址
func RemoveAddressByUid(uid, addressId int64) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFun()
	err := col.Update(bson.M{"_id": uid}, bson.M{"$pull": bson.M{"Address": bson.M{"AddressId": addressId}}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("删除收货地址err: %v", err.Error())
	}
	return err
}

//根据中间表id集合查询盲盒商品列表
func GetWishItemListByWishBoxItemIds(WishBoxItemIds []int64) ([]*share_message.WishItem, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFun()
	list := make([]*share_message.WishItem, PAGE_LIMIT)
	match := bson.M{"_id": bson.M{"$in": WishBoxItemIds}}
	lookup := bson.M{"from": TABLE_WISH_ITEM, "localField": "WishItemId", "foreignField": "_id", "as": "goods"}
	project := bson.M{
		"Id":            "$goods._id",
		"Name":          "$goods.Name",
		"Icon":          "$goods.Icon",
		"Price":         "$goods.Price",
		"RecoveryPrice": "$goods.RecoveryPrice",
	}
	pipeCond := []bson.M{
		{"$match": match},
		{"$lookup": lookup},
		{"$unwind": "$goods"},
		{"$project": project},
	}

	err := col.Pipe(pipeCond).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("根据中间表id集合查询盲盒商品列表err: %v", err.Error())
	}
	return list, err
}

//新增用户兑换盲盒商品记录
func InsertPlayerExchangeLog(data []interface{}) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_EXCHANGE_LOG)
	defer closeFun()

	err := col.Insert(data...)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("批量插入用户兑换记录失败")
	}
	return err
}

//新增用户回收盲盒商品记录
func InsertPlayerRecycleLog(data *share_message.WishRecycleOrder) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()
	data.InitTime = easygo.NewInt64(time.Now().Unix())
	data.Id = easygo.NewInt64(NextId(TABLE_WISH_RECYCLE_ORDER))
	err := col.Insert(data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("新增用户(%v)回收盲盒商品记录err: %v", data.Id, err.Error())
	}
	return err
}

//查询用户添加的地址
func GetAddressByPlayerIdAndId(playerId, addressId int64) (*share_message.WishAddress, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFun()
	condBson := bson.M{"_id": playerId, "Address.AddressId": addressId}
	data := share_message.WishPlayer{}
	err := col.Find(condBson).One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("查询用户添加的地址: player(%v), err: %v", playerId, err.Error())
		return nil, err
	}
	var addressArr *share_message.WishAddress
	if len(data.GetAddress()) > 0 {
		addressArr = data.GetAddress()[0]
	}
	return addressArr, err
}

//查询玩家收藏的盲盒集合
func GetCollectedBoxByPlayerId(playerId int64, skip, limit, reqType int) ([]*share_message.PlayerWishCollection, int, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_COLLECTION)
	defer closeFun()

	onSaleCount := 0
	if reqType == WISH_ONSALE {
		total, err := col.Find(bson.M{"PlayerId": playerId, "SaleStatus": WISH_ONSALE}).Count()
		if err != nil && err.Error() != mgo.ErrNotFound.Error() {
			logs.Error("查询玩家收藏盲盒上架的个数: player(%v), err: %v", playerId, err.Error())
			return nil, 0, 0
		}
		onSaleCount = total
	}

	SaleOutCount, err := col.Find(bson.M{"PlayerId": playerId, "SaleStatus": WISH_SALEOUT}).Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("查询玩家收藏盲盒下架的个数: player(%v), err: %v", playerId, err.Error())
		return nil, 0, 0
	}

	if onSaleCount > 0 || SaleOutCount > 0 {
		data := make([]*share_message.PlayerWishCollection, PAGE_LIMIT)
		if limit > PAGE_LIMIT {
			limit = PAGE_LIMIT
		}
		err = col.Find(bson.M{"PlayerId": playerId, "SaleStatus": reqType}).Sort("-CreateTime").Skip(skip).Limit(limit).All(&data)
		if err != nil && err.Error() != mgo.ErrNotFound.Error() {
			logs.Error("分页查询玩家收藏盲盒数据: PlayerId: %v, SaleStatus: %v ,err: %v", playerId, reqType, err.Error())
		}
		return data, onSaleCount, SaleOutCount
	}
	return nil, 0, 0
}

//查询玩家收藏的盲盒集合
func GetAllCollectedBoxByPlayerId(playerId int64) ([]*share_message.PlayerWishCollection, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_COLLECTION)
	defer closeFun()

	data := make([]*share_message.PlayerWishCollection, 0)
	err := col.Find(bson.M{"PlayerId": playerId}).Sort("-CreateTime").All(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("查询玩家(%v)收藏的盲盒集合 err: %v", playerId, err.Error())
	}
	return data, err
}

// 收藏物品的记录
func InsertPlayerWishCollection(data *share_message.PlayerWishCollection) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_COLLECTION)
	defer closeFun()
	err := col.Insert(data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("收藏物品的记录 err: %v", err.Error())
	}
	return err
}

//获取玩家收藏的盲盒
func GetPlayerWishCollection(pid, boxId int64) *share_message.PlayerWishCollection {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_COLLECTION)
	defer closeFun()
	pwc := &share_message.PlayerWishCollection{}
	err := col.Find(bson.M{"PlayerId": pid, "WishBoxId": boxId}).One(&pwc)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取玩家(%v)收藏的盲盒 err: %v", pid, err.Error())
	}
	return pwc
}

//批量删除收藏的盲盒
func RemovePlayerWishCollectionByIds(playerId int64, ids []int64) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_COLLECTION)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"PlayerId": playerId, "WishBoxId": bson.M{"$in": ids}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("批量用户(%v)删除收藏的盲盒 err: %v", playerId, err.Error())
	}
	return err
}

//查询用户的已许愿列表（区分上下架）
func GetPlayerWishDataByPlayerId(playerId int64, skip, limit, reqType int) ([]*share_message.PlayerWishData, int, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_DATA)
	defer closeFun()

	onSaleTotal := 0
	//获取上架的许愿物品需要同时要返回下架的物品
	if reqType == WISH_ONSALE {
		total, err := col.Find(bson.M{"PlayerId": playerId, "SaleStatus": WISH_ONSALE}).Count()
		if err != nil && err.Error() != mgo.ErrNotFound.Error() {
			logs.Error("获取玩家许愿的上架盲盒: playerId: %v, err: %v", playerId, err.Error())
			return nil, 0, 0
		}
		onSaleTotal = total
	}

	saleOutTotal, err := col.Find(bson.M{"PlayerId": playerId, "SaleStatus": WISH_SALEOUT}).Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取玩家许愿的下架盲盒：playerId: %v, err: %v", playerId, err.Error())
		return nil, 0, 0
	}
	data := make([]*share_message.PlayerWishData, PAGE_LIMIT)
	if limit > PAGE_LIMIT {
		limit = PAGE_LIMIT
	}
	err = col.Find(bson.M{"PlayerId": playerId, "SaleStatus": reqType}).Sort("-CreateTime").Skip(skip).Limit(limit).All(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("分页查询玩家许愿盲盒：playerId: %v, SaleStatus: %v, skip: %v, limit: %v, err: %v", playerId, reqType, skip, limit, err.Error())
	}
	return data, onSaleTotal, saleOutTotal
}

//查询用户所有的已许愿列表（
func GetPlayerAllWishDataByPlayerId(playerId int64) ([]*share_message.PlayerWishData, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_DATA)
	defer closeFun()

	data := make([]*share_message.PlayerWishData, 0)
	err := col.Find(bson.M{"PlayerId": playerId}).All(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("查询用户(%v)所有的已许愿列表err: %v", playerId, err)
	}
	return data, err
}

//删除已许愿记录
func RemovePlayerWishDataByIds(playerId int64, ids []int64) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_DATA)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"PlayerId": playerId, "WishBoxId": bson.M{"$in": ids}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("删除id: %v 已许愿记录err: %v", playerId, err)
	}
	return err
}

//根据id查询盲盒列表
func GetWishBoxListByIds(ids []int64) (data []*h5_wish.CollectBox, err error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFun()

	err = col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("根据id(%v)查询盲盒列表err: %v", ids, err)
	}
	return data, err
}

//根据盲盒ID查询有无预售物品
func GetPreSaleBoxByItemIds(ids []int64) (data []*h5_wish.CollectBox, err error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFun()
	err = col.Find(bson.M{"WishBoxId": bson.M{"$in": ids}, "Status": 2}).Select(bson.M{"WishBoxId": 1}).All(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("根据盲盒ID(%v)查询有无预售物品err: %v", ids, err)
	}
	return data, err
}

//根据盲盒Id查询盲盒出售物品列表
func QueryWishBoxItemByWishBoxIds(wishBoxIds []int64) (data []*share_message.WishBoxItem, err error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFun()

	err = col.Find(bson.M{"WishBoxId": bson.M{"$in": wishBoxIds}}).Sort("-CreateTime").All(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("根据盲盒Id(%v)查询盲盒出售物品列表err: %v", wishBoxIds, err)
	}
	return data, err
}

//获取待兑换的盲盒个数
func CountToExchangePlayerWish(playerId int64) int {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_ITEM)
	defer closeFun()
	count, _ := col.Find(bson.M{"PlayerId": playerId, "Status": WISH_TO_EXCHANGE}).Count()
	return count
}

//获取未读待兑换的盲盒个数
func CountUnReadToExchangePlayerWish(playerId int64) int {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_ITEM)
	defer closeFun()
	count, _ := col.Find(bson.M{"PlayerId": playerId, "IsRead": false, "Status": WISH_TO_EXCHANGE}).Count()
	return count
}

//获取用户待兑换，已兑换，已回收总条数
func CountPlayerWish(playerId int64) (toExchange, exchanged, recycle int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_ITEM)
	defer closeFun()
	toExchange, _ = col.Find(bson.M{"PlayerId": playerId, "$or": []bson.M{
		{"Status": WISH_TO_EXCHANGE},
		{"Status": WISH_RECYCLE_CHECK}}}).Count()
	exchanged, _ = col.Find(bson.M{"PlayerId": playerId, "Status": WISH_EXCHANGED}).Count()
	recycle, _ = col.Find(bson.M{"PlayerId": playerId, "Status": WISH_RECYCLED}).Count()

	return toExchange, exchanged, recycle
}

//根据用户ID 获取玩家物品列表
func GetPlayWishItemByPlayerId(playerId int64, reqType, page, pageSize int, sortStr string) ([]*share_message.PlayerWishItem, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_ITEM)
	defer closeFun()
	wishItem := make([]*share_message.PlayerWishItem, 0)
	queryBson := bson.M{"PlayerId": playerId}
	if reqType == WISH_TO_CONTROL {
		queryBson["$or"] = []bson.M{
			{"Status": WISH_TO_EXCHANGE},
			{"Status": WISH_RECYCLE_CHECK}}
	} else {
		queryBson["Status"] = reqType
	}
	err := col.Find(queryBson).Sort(sortStr).Skip(page * pageSize).Limit(pageSize).All(&wishItem)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("根据用户ID%v 获取玩家物品所有列表err: %v", playerId, err)
	}
	return wishItem, err
}

//根据用户ID 获取玩家物品所有列表
func GetAllPlayWishItemByPlayerId(playerId int64, reqType int, sortStr string) ([]*share_message.PlayerWishItem, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_ITEM)
	defer closeFun()
	wishItem := make([]*share_message.PlayerWishItem, 0)
	queryBson := bson.M{"PlayerId": playerId}
	if reqType == WISH_TO_EXCHANGE {
		queryBson["$or"] = []bson.M{
			{"Status": WISH_TO_EXCHANGE},
			{"Status": WISH_RECYCLE_CHECK}}
	} else {
		queryBson["Status"] = reqType
	}
	err := col.Find(queryBson).Sort(sortStr).All(&wishItem)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("根据用户ID%v 获取玩家物品所有列表err: %v", playerId, err)
	}
	return wishItem, err
}

//根据盲盒ID 找到对应的盲盒信息
func GetWishBoxByIds(ids []int64) ([]*share_message.WishBox, error) {
	lst := make([]*share_message.WishBox, 0)
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFunc()
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("根据盲盒ID(%v) 找到对应的盲盒信息err: %v", ids, err)
	}
	return lst, err
}

//编辑玩家已查阅新增物品
func EditPlayerItemIsRead(uid int64) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_ITEM)
	defer closeFun()
	_, err := col.UpdateAll(bson.M{"PlayerId": uid, "IsRead": false, "Status": 0}, bson.M{"$set": bson.M{"IsRead": true}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("编辑玩家(%v)已查阅新增物品err: %v", uid, err)
	}
	return err
}

// 获取当前时间24小时内占领用户信息 WishOccupied
func GetWishOccupiedCurTime(num int) []*share_message.WishOccupied {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFun()
	list := make([]*share_message.WishOccupied, 0)
	m := []bson.M{
		{"$match": bson.M{"CreateTime": bson.M{"$gte": time.Now().Unix() - 24*3600}}},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取当前时间24小时内占领用户信息err: %v", err)
	}
	return list
}

// 查询当前时间24小时内获取愿望款的用户
func GotWishPlayerLogByCurTime(num int) []*share_message.PlayerWishData {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_DATA)
	defer closeFun()
	list := make([]*share_message.PlayerWishData, 0)
	m := []bson.M{
		{"$match": bson.M{"Status": WISH_CHALLENGE_SUCCESS, "FinishTime": bson.M{"$gte": time.Now().Unix() - (24 * 3600)}}},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("查询当前时间24小时内获取愿望款的用户err: %v", err)
	}
	return list
}

//随机指定数量获取物品图片
func GetRangeItemIcon(num int) ([]string, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ITEM)
	defer closeFun()
	wishItem := make([]*share_message.WishItem, 0)
	m := []bson.M{
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&wishItem)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("随机指定数量获取物品图片: %v", err)
		return nil, err
	}
	result := make([]string, 0)

	for _, v := range wishItem {
		result = append(result, v.GetIcon())
	}

	return result, err
}

//获取回收比例
func GetRecycleRatio(recycleType int) (int32, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OTHER_CFG)
	defer closeFun()
	section := &share_message.WishRecycleSection{}
	err := col.Find(bson.M{"_id": TABLE_WISH_RECYCLE_SECTION}).One(&section)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取回收比例: %v", err)
		return 0, err
	}
	var ratio int32
	if recycleType == RECYCLE_BY_PLAYER {
		ratio = section.GetPlayer()
	} else {
		ratio = section.GetPlatform()
	}
	return ratio, nil
}

//获取回收安全阀值
func GetRecycleSafeValve() (int64, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OTHER_CFG)
	defer closeFun()
	section := &share_message.WishRecycleSection{}
	err := col.Find(bson.M{"_id": TABLE_WISH_RECYCLE_SECTION}).One(&section)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取回收安全阀值失败: %v", err)
		return 0, err
	}

	return section.GetOrderThreshold(), nil
}

//设置玩家物品表回收价格
func RecycleToPlayerItem(data map[int64]int64, status, recycleType int) {
	var update []interface{}
	for k, v := range data { //[玩家表ID]回收价格
		uData := &share_message.PlayerWishItem{
			Status:       easygo.NewInt32(status),
			UpdateTime:   easygo.NewInt64(time.Now().Unix()),
			RecyclePrice: easygo.NewInt64(v),
			RecycleType:  easygo.NewInt32(recycleType),
		}
		update = append(update, bson.M{"_id": k}, bson.M{"$set": uData})
	}
	if len(update) < 1 {
		return
	}
	UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_WISH_ITEM, update)
}

//获取邮费表所有信息
func GetWishMailSection() (*h5_wish.PostageResp, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OTHER_CFG)
	defer closeFun()
	var postage *h5_wish.PostageResp
	err := col.Find(bson.M{"_id": TABLE_WISH_MAIL_SECTION}).One(&postage)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取邮费err: %v", err)
	}
	return postage, err
}

//获取地区邮费  ----看是否可以通过redis优化
func CheckAreaPostage(province string) int32 {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OTHER_CFG)
	defer closeFun()

	var postage *h5_wish.PostageResp
	err := col.Find(bson.M{"_id": TABLE_WISH_MAIL_SECTION}).One(&postage)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取邮费err: %v", err)
		return 0
	}
	if province == "浙江省" || province == "江苏省" || province == "上海市" {
		return postage.GetPostage1()
	} else if IsContainsStr(province, postage.GetRemoteAreaList()) == -1 {
		return postage.GetPostage2()
	}
	return postage.GetPostage3()
}

//获取回收理由
func GetRecycleReason() ([]*share_message.RecycleReason, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_RECYCLE_REASON)
	defer closeFun()

	var reason []*share_message.RecycleReason
	err := col.Find(nil).All(&reason)
	if err != nil {
		logs.Error("获取回收理由失败err: %v", err)
		return nil, err
	}

	return reason, nil
}

//更新用户待兑换商品状态
func ExchangeToPlayerItem(playerWishItemIds []int64, editType int) {
	var data []interface{}
	for _, v := range playerWishItemIds {
		uData := &share_message.PlayerWishItem{
			Status:     easygo.NewInt32(editType),
			UpdateTime: easygo.NewInt64(time.Now().Unix()),
		}
		data = append(data, bson.M{"_id": v}, bson.M{"$set": uData})
	}
	if len(data) == 0 {
		return
	}
	UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_WISH_ITEM, data)
}

//获取回收支付风控配置
func GetConfigWishPaymentFormDB() *share_message.WishRecycleSection {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OTHER_CFG)
	defer closeFun()
	section := &share_message.WishRecycleSection{}
	err := col.Find(bson.M{"_id": TABLE_WISH_RECYCLE_SECTION}).One(&section)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取回收支付风控配置: %v", err)
		return nil
	}
	return section
}

// 获取支付预警
func GetWishPayWarnCfg() *brower_backstage.WishPayWarnCfg {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OTHER_CFG)
	defer closeFun()

	list := &share_message.WishPayWarnCfg{}
	queryBson := bson.M{"_id": TABLE_WISH_PAY_WARN_CFG}
	err := col.Find(queryBson).One(&list)
	if err != nil && err != mgo.ErrNotFound {
		logs.Error("获取支付预警失败: ", err)
		return nil
	}
	res := &brower_backstage.WishPayWarnCfg{
		WithdrawalTime:        easygo.NewInt64(list.GetWithdrawalTime()),
		WithdrawalTimes:       easygo.NewInt64(list.GetWithdrawalTimes()),
		WithdrawalGoldRate:    easygo.NewInt64(list.GetWithdrawalGoldRate()),
		WithdrawalGold:        easygo.NewInt64(list.GetWithdrawalGold()),
		WithdrawalDiamondRate: easygo.NewInt64(list.GetWithdrawalDiamondRate()),
		WithdrawalDiamond:     easygo.NewInt64(list.GetWithdrawalDiamond()),
		PhoneList:             list.GetPhoneList(),
	}
	return res
}

// 获取回收订单
func GetWishRecycleOrderByPaymentOrderId(id string) *share_message.WishRecycleOrder {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()
	player := &share_message.WishRecycleOrder{}
	err := col.Find(bson.M{"PaymentOrderId": id}).One(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return player
}

// 更新回收订单状态
func UpdateWishOrder(id string, status int32) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()
	data := &share_message.WishRecycleOrder{
		Status: easygo.NewInt32(status),
	}
	if status == 1 {
		curTime := time.Now().Unix()
		data.RecycleTime = easygo.NewInt64(curTime)
	}

	err := col.Update(bson.M{"PaymentOrderId": id}, bson.M{"$set": data})
	easygo.PanicError(err)
}

// 更新用户物品信息状态
func UpdatePlayerWishItemStatusByOrder(ids []int64, status int32) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": ids}, "Status": 3}, bson.M{"$set": bson.M{"Status": status}})
	easygo.PanicError(err)
}

// 更新相应回收订单的状态已经回收物品状态
func UpdateWishOrderByPayOrder(oid string) {
	o := GetRedisOrderObj(oid)
	if o == nil {
		return
	}
	wishO := GetWishRecycleOrderByPaymentOrderId(o.GetOrderId())
	var ids []int64
	for _, v := range wishO.GetRecycleItemList() {
		ids = append(ids, v.GetPlayerItemId())
	}
	if o.GetStatus() == ORDER_ST_FINISH {
		UpdateWishOrder(o.GetOrderId(), 1)
		// 更新所回收物品的状态
		UpdatePlayerWishItemStatusByOrder(ids, 2)
	} else if o.GetStatus() == ORDER_ST_CANCEL || o.GetStatus() == ORDER_ST_REFUSE {
		// 更新所回收物品的状态
		UpdatePlayerWishItemStatusByOrder(ids, 0)
	}
}

// 获取充值活动配置
func GetWishCoinRechargeActCfg() []*share_message.WishCoinRechargeActivityCfg {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_COIN_RECHARGE_ACT_CFG)
	defer closeFun()
	var data []*share_message.WishCoinRechargeActivityCfg
	err := col.Find(nil).Sort("_id").All(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取充值活动配置失败：%v", err.Error())
		return nil
	}
	return data
}

// 获取用户充值记录
func GetPlayerRechargeActLog(pid int64, pageReq *h5_wish.DataPageReq) []*share_message.PlayerRechargeActLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER_RECHARGE_ACT_LOG)
	defer closeFun()
	page := int(pageReq.GetPage() - 1)
	pageSize := int(pageReq.GetPageSize())

	var data []*share_message.PlayerRechargeActLog
	err := col.Find(bson.M{"PlayerId": pid}).Sort("-CreateTime").Skip(page * pageSize).Limit(pageSize).All(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取用户首充记录：%v", err.Error())
		return nil
	}
	return data
}

// 用户充值记录
func AddPlayerRechargeActLog(pid int64, money, coin, giveCoin int64, giveType int32) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER_RECHARGE_ACT_LOG)
	defer closeFun()
	err := col.Insert(bson.M{"_id": NextId(TABLE_WISH_PLAYER_RECHARGE_ACT_LOG), "PlayerId": pid,
		"CreateTime": easygo.NowTimestamp(), "Money": money, "Coin": coin,
		"GiveCoin": giveCoin, "GiveType": giveType})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("增加用户充值记录失败: %v", err.Error())
	}
}

// 获取用户首充记录
func GetPlayerRechargeActFirst(pid int64) *share_message.PlayerRechargeActFirst {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER_RECHARGE_ACT_FIRST)
	defer closeFun()
	var data *share_message.PlayerRechargeActFirst
	err := col.Find(bson.M{"_id": pid}).One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取用户首充记录：%v", err.Error())
		return nil
	}
	return data
}

// 用户首充记录
func AddPlayerRechargeActFirst(pid int64, level int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER_RECHARGE_ACT_FIRST)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": pid}, bson.M{"$push": bson.M{"Levels": level}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("增加用户首充记录失败: %v", err.Error())
	}
}

// 获取用户访问记录表
func GetWishPlayerAccessFormDB(pid int64) *share_message.WishPlayerAccessLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER_ACCESS_LOG)
	defer closeFun()
	var data *share_message.WishPlayerAccessLog
	err := col.Find(bson.M{"_id": pid}).One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("获取用户访问记录表失败: %v", err.Error())
		return nil
	}
	return data
}

// 更新指定用户访问记录表
func UpsertWishPlayerAccessFormDB(pid, val int64, key string) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER_ACCESS_LOG)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": pid}, bson.M{"$set": bson.M{key: val}})
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("更新指定用户访问记录表失败: %v", err.Error())
	}
}

//插入埋点日志
func InsertWishBPLog(playerId int64, bpType int32) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BURYING_POINT_LOG)
	defer closeFun()
	data := &share_message.WishBuryingPointLog{
		Id:        easygo.NewInt64(NextId(TABLE_WISH_BURYING_POINT_LOG)),
		PlayerId:  easygo.NewInt64(playerId),
		EventType: easygo.NewInt32(bpType),
		Time:      easygo.NewInt64(easygo.NowTimestamp()),
	}
	err := col.Insert(data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("新增用户(%v)埋点日志err: %v", data.Id, err.Error())
	}
}

//随机获取回收说明
func GetOneRecycleNote() string {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OTHER_CFG)
	defer closeFun()
	noteCfg := &share_message.RecycleNoteCfg{}
	note := ""
	err := col.Find(bson.M{"_id": TABLE_WISH_RECYCLE_NOTE_CFG}).One(&noteCfg)
	if err != nil {
		logs.Info("随机获取回收说明失败:%v", err)
		return note
	}
	noteLen := len(noteCfg.GetText())
	if noteLen > 0 {
		index := util.RandIntn(noteLen)
		note = noteCfg.GetText()[index]
	}
	return note
}
