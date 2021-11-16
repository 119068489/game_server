package for_game

import (
	"errors"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/h5_wish"
	"game_server/pb/share_message"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/akqp2019/mgo"

	"github.com/akqp2019/mgo/bson"
)

var mu easygo.Mutex

const (
	AFTER_TIME_GUARDIAN = 24 * 3600 // 守护者的时间,24小时过期.
)

//奖项状态:0待领取，1已领取
const (
	WISH_ACTIVITY_PRIZE_STATUS_0 = 0 // 0待领取
	WISH_ACTIVITY_PRIZE_STATUS_1 = 1 // 1已领取
)

const (
	WISH_ACT_H5_WEEK      = 1 // 前端参数,周榜
	WISH_ACT_H5_MONTH     = 2 // 前端参数,月榜
	WISH_ACT_H5_WEEK_NOW  = 3 // 前端参数本周周榜
	WISH_ACT_H5_MONTH_NOW = 4 // 前端参数本周月榜
)

// 规则类型： 1、次数 2、天数 3、周排名 4、月排名
const (
	WISH_ACTIVITY_DATA_TYPE_1 = 1
	WISH_ACTIVITY_DATA_TYPE_2 = 2
	WISH_ACTIVITY_DATA_TYPE_3 = 3
	WISH_ACTIVITY_DATA_TYPE_4 = 4
)
const (
//WISH_PLAYER_BASE_TYPES_2 = 2 // 许愿池白名单用户
)

// 随机找n个盒子
func FindWishBoxByNum(num int) ([]*share_message.WishBox, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFun()
	queryBson := bson.M{"Status": 1}
	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
		//{"$project": bson.M{"_id": 1}},
	}
	query := col.Pipe(m)
	var list []*share_message.WishBox
	err := query.All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
		return nil, err
	}
	return list, nil
}

// 找到价格最高的物品
func FindMaxPriceBoxItemByBoxId(boxId int64, ids []int64) (*share_message.WishBoxItem, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFun()
	query := bson.M{"WishBoxId": boxId, "Status": 1}
	if len(ids) > 0 {
		query["WishItemId"] = bson.M{"$nin": ids}
	}

	var boxItem *share_message.WishBoxItem
	err := col.Find(query).Sort("-Price").One(&boxItem)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
		return nil, err
	}
	return boxItem, nil
}

// 根据盲盒id抽出小奖的物品
func FindMaxPriceBoxItemByBoxIdAndLv(boxId int64, lv int32) []*share_message.WishBoxItem {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFun()
	query := bson.M{"WishBoxId": boxId, "RewardLv": lv}

	boxItems := make([]*share_message.WishBoxItem, 0)
	err := col.Find(query).All(&boxItems)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())

	}
	return boxItems
}
func FindToolMaxPriceBoxItemByBoxIdAndLv(boxId int64, lv int32) []*share_message.WishBoxItem {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOOL_WISH_BOX_ITEM)
	defer closeFun()
	query := bson.M{"RewardLv": lv, "WishBoxId": boxId}

	boxItems := make([]*share_message.WishBoxItem, 0)
	err := col.Find(query).All(&boxItems)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())

	}
	return boxItems
}

//获取物品配置
func QueryWishItemById(id int64) (*share_message.WishItem, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ITEM)
	defer closeFun()
	var wishItem *share_message.WishItem
	err := col.Find(bson.M{"_id": id}).One(&wishItem)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
		return nil, err
	}
	return wishItem, nil
}

// 随机抽取挑战成功的记录
/*func GetRandPlayer(num int) []*share_message.PlayerBase {
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
*/
// 获得物品的记录
func InsertPlayerWishItemToDb(wishItem *share_message.PlayerWishItem) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_ITEM)
	defer closeFun()
	err := col.Insert(wishItem)
	easygo.PanicError(err)
	return err
}

// 随机抽取获得物品的记录
func GetRandProducts(num int) []*share_message.PlayerWishItem {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_ITEM)
	defer closeFun()
	list := make([]*share_message.PlayerWishItem, 0)
	m := []bson.M{
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)
	return list
}

// 随机抽取获得该盲盒物品的记录
func GetRandBoxProducts(boxId int64, num int) []*share_message.PlayerWishItem {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_ITEM)
	defer closeFun()
	list := make([]*share_message.PlayerWishItem, 0)
	queryBson := bson.M{"WishBoxId": boxId}
	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)
	return list
}

// 根据商品中间表的ids查询中间表信息
func GetWishBoxItemByIdsFromDB(ids []int64) ([]*share_message.WishBoxItem, error) {
	var lst []*share_message.WishBoxItem
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFunc()
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).Sort("PerRate").All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return lst, err
}

// 根据数量查询
func GetWishBoxItemByIdsAndNumFromDB(ids []int64) ([]*share_message.WishBoxItem, error) {
	var lst []*share_message.WishBoxItem
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFunc()
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}, "PerNum": bson.M{"$gt": 0}, "LocalNum": 0}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return lst, err
}

// 根据商品中间表的id查询中间表信息
func GetWishBoxItemByIdFromDB(id int64) *share_message.WishBoxItem {
	var boxItem *share_message.WishBoxItem
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFunc()
	err := col.Find(bson.M{"_id": id}).One(&boxItem)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return boxItem
}
func GetToolWishBoxItemByIdFromDB(id int64) *share_message.WishBoxItem {
	var boxItem *share_message.WishBoxItem
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOOL_WISH_BOX_ITEM)
	defer closeFunc()
	err := col.Find(bson.M{"_id": id}).One(&boxItem)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return boxItem
}

// 商品ids找到真正的商品信息
func GetWishItemByIdsFromDB(ids []int64) ([]*share_message.WishItem, error) {
	var lst []*share_message.WishItem
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ITEM)
	defer closeFunc()
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("商品ids(%v)找到真正的商品信息 err: %s", ids, err.Error())
	}
	return lst, err
}
func GetWishItemByIdFromDB(id int64) *share_message.WishItem {
	item := &share_message.WishItem{}
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ITEM)
	defer closeFunc()
	err := col.Find(bson.M{"_id": id}).One(&item)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return item
}

// 随机抽取获得真正的物品
func GetRandWishItem(num int) []*share_message.WishItem {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ITEM)
	defer closeFun()
	list := make([]*share_message.WishItem, 0)

	m := []bson.M{
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)
	return list
}

// 随机抽取获得盲盒中真正的物品
func GetRandWishBoxItem(boxId int64, num int) []*share_message.WishItem {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ITEM)
	defer closeFun()
	list := make([]*share_message.WishItem, 0)

	m := []bson.M{
		{"$match": bson.M{"_id": boxId}},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)
	return list
}

// 获得物品的记录
func InsertWishLogToDB(wishLog *share_message.WishLog) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_LOG)
	defer closeFun()
	err := col.Insert(wishLog)
	easygo.PanicError(err)
	return err
}

// 随机条数查询挑战记录 WishLog
func GetRangeWishLog(num int, result bool) []*share_message.WishLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_LOG)
	defer closeFun()
	list := make([]*share_message.WishLog, 0)
	m := []bson.M{
		{"$match": bson.M{"Result": result}},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)
	return list
}
func GetRangeWishLogByBoxId(boxId int64, num int) []*share_message.WishLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_LOG)
	defer closeFun()
	list := make([]*share_message.WishLog, 0)
	m := []bson.M{
		{"$match": bson.M{"WishBoxId": boxId}},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)
	return list
}

// 获取最新的挑战者
func GetOneSuccessDarer(boxId int64) *share_message.WishLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_LOG)
	defer closeFun()
	var wishLog *share_message.WishLog
	err := col.Find(bson.M{"_id": boxId, "Result": true}).Sort("-CreateTime").One(&wishLog)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return wishLog
}

// 随机条数查询有 挑战者的盲盒 挑战者的id列表
func GetRangWishLogGuardians(num int) []int64 {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFun()

	m := []bson.M{
		{"$match": bson.M{"Status": WISH_PUT_ON_STATUS, "GuardianId": bson.M{"$gt": 0}}},
		{"$match": bson.M{"Status": WISH_PUT_ON_STATUS}},
		{"$sample": bson.M{"size": num}},
		{"$group": bson.M{"_id": "$GuardianId"}},
	}

	query := col.Pipe(m)
	res := make([]*share_message.WishBox, 0)
	err := query.All(&res)
	easygo.PanicError(err)
	ids := make([]int64, 0)
	for _, value := range res {
		ids = append(ids, value.GetId())
	}
	return ids
}

// 通过守护者id找到对应的盲盒列表
func GetRangWishLogByGuardian(dId int64) []*share_message.WishBox {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFun()
	list := make([]*share_message.WishBox, 0)
	err := col.Find(bson.M{"GuardianId": dId}).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}

	return list
}

//更新盲盒
func UpWishBox(id int64, WishBox *share_message.WishBox, num ...int) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFun()
	var err error
	if len(num) > 0 {
		err = col.Update(bson.M{"_id": id}, bson.M{"$set": WishBox, "$inc": bson.M{"WinNum": num[0]}})
	} else {

		err = col.Update(bson.M{"_id": id}, bson.M{"$set": WishBox})
	}
	return err
}

// 获取菜单列表
func GetMenuList() []*share_message.WishMenu {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_MENU)
	defer closeFun()

	menus := make([]*share_message.WishMenu, 0)
	err := col.Find(bson.M{}).All(&menus)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return menus
}

// 热门品牌数据
func GetHotWishBrandList() []*share_message.WishBrand {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BRAND)
	defer closeFun()
	brands := make([]*share_message.WishBrand, 0)
	err := col.Find(bson.M{"IsHot": true, "Status": 1}).Sort("-HostWeight").All(&brands)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return brands
}

// 品牌数据
func GetNotHotWishBrandList(ids []int64, num int) ([]*share_message.WishBrand, []*share_message.WishBrand) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BRAND)
	defer closeFun()
	brands := make([]*share_message.WishBrand, 0)
	brands1 := make([]*share_message.WishBrand, 0)
	query := bson.M{"Status": 1}
	if len(ids) > 0 {
		query["_id"] = bson.M{"$nin": ids}
	}
	err := col.Find(query).Sort("Type").All(&brands)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	err = col.Find(query).Sort("-HostWeight").Limit(num).All(&brands1)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return brands, brands1
}

// 热门物品类型
func GetHotWishTypeList() []*share_message.WishItemType {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ITEM_TYPE)
	defer closeFun()
	types := make([]*share_message.WishItemType, 0)
	err := col.Find(bson.M{"IsHot": true, "Status": 1}).Sort("-HostWeight").All(&types)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return types
}

// 物品类型
func GetNotHotWishTypeList(ids []int64, num int) ([]*share_message.WishItemType, []*share_message.WishItemType) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ITEM_TYPE)
	defer closeFun()
	query := bson.M{"Status": 1}
	if len(ids) > 0 {
		query["_id"] = bson.M{"$nin": ids}
	}

	types := make([]*share_message.WishItemType, 0)
	err := col.Find(query).Sort("Type").All(&types)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	types1 := make([]*share_message.WishItemType, 0)
	err = col.Find(query).Sort("-HostWeight").Limit(num).All(&types1)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return types, types1
}

// 获得指定条数的热门品牌
func GetRangeHotWishBrandByNum(num int) []*share_message.WishBrand {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BRAND)
	defer closeFun()
	list := make([]*share_message.WishBrand, 0)
	m := []bson.M{
		{"$match": bson.M{"IsHot": true}},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)
	return list
}

// 获得指定条数的占领时间数据
func GetRangeWishWishOccupiedByNum(boxId int64, num int) []*share_message.WishOccupied {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFun()
	list := make([]*share_message.WishOccupied, 0)
	m := []bson.M{
		{"$match": bson.M{"WishBoxId": boxId}},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)
	return list
}

// 获得指定条数的非热门品牌
func GetRangeWishBrandByNum(num int) []*share_message.WishBrand {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BRAND)
	defer closeFun()
	list := make([]*share_message.WishBrand, 0)
	m := []bson.M{
		{"$match": bson.M{"IsHot": bson.M{"$ne": true}}},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)
	return list
}

//添加许愿
func AddWishData(data *share_message.PlayerWishData) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_DATA)
	defer closeFun()
	data.Id = easygo.NewInt64(NextId(TABLE_PLAYER_WISH_DATA))
	data.CreateTime = easygo.NewInt64(time.Now().Unix())
	data.Status = easygo.NewInt32(0)
	err := col.Insert(data)
	return err
}

//修改许愿物品
func UpWishBoxItemId(playerId, wishBoxId, wishBoxItemId, wishItemId int64, status int32, productUrl string) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_DATA)
	defer closeFun()
	err := col.Update(bson.M{"PlayerId": playerId, "WishBoxId": wishBoxId, "Status": status}, bson.M{"$set": bson.M{"WishBoxItemId": wishBoxItemId, "CreateTime": time.Now().Unix(), "ProductUrl": productUrl, "WishItemId": wishItemId}})
	return err
}

//修改许愿
func UpWishData(playerId, wishBoxId int64, wishData *share_message.PlayerWishData) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_DATA)
	defer closeFun()
	err := col.Update(bson.M{"PlayerId": playerId, "WishBoxId": wishBoxId}, bson.M{"$set": wishData})
	return err
}
func UpWishDataById(id int64, wishData *share_message.PlayerWishData) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_DATA)
	defer closeFun()
	err := col.Update(bson.M{"_id": id}, bson.M{"$set": wishData})
	return err
}

//获取许愿信息
func GetWishData(playerId, wishBoxId int64) (*share_message.PlayerWishData, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_DATA)
	defer closeFun()
	d := &share_message.PlayerWishData{}
	e := col.Find(bson.M{"PlayerId": playerId, "WishBoxId": wishBoxId}).One(&d)
	if e != nil && e.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", e.Error())
	}
	return d, e
}
func GetWishDataByStatus(playerId, wishBoxId int64, status int32) (*share_message.PlayerWishData, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_DATA)
	defer closeFun()
	d := &share_message.PlayerWishData{}
	e := col.Find(bson.M{"PlayerId": playerId, "WishBoxId": wishBoxId, "Status": status}).Sort("-CreateTime").One(&d)
	if e != nil && e.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", e.Error())
	}
	return d, e
}
func GetWishDataByBidsStatus(playerId int64, wishBoxId []int64, status int32) ([]*share_message.PlayerWishData, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_DATA)
	defer closeFun()
	d := make([]*share_message.PlayerWishData, 0)
	e := col.Find(bson.M{"PlayerId": playerId, "WishBoxId": bson.M{"$in": wishBoxId}, "Status": status}).All(&d)
	if e != nil && e.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", e.Error())
	}
	return d, e
}

//获取盲盒信息
func GetWishBox(WishBoxId int64, status []int64) (*share_message.WishBox, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFun()
	d := &share_message.WishBox{}
	e := col.Find(bson.M{"_id": WishBoxId, "Status": bson.M{"$in": status}}).One(&d)
	if e != nil && e.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", e.Error())
	}
	return d, e
}

//获取盲盒信息
func GetWishBoxJustById(WishBoxId int64) (*share_message.WishBox, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFun()
	d := &share_message.WishBox{}
	e := col.Find(bson.M{"_id": WishBoxId}).One(&d)
	if e != nil && e.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", e.Error())
	}
	return d, e
}

//获取盲盒出售物品
func GetWishBoxItem(id int64) (*share_message.WishBoxItem, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFun()
	d := &share_message.WishBoxItem{}
	e := col.Find(bson.M{"_id": id}).One(&d)
	if e != nil && e.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", e.Error())
	}
	return d, e
}
func GetToolWishBoxItem(id int64) (*share_message.WishBoxItem, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOOL_WISH_BOX_ITEM)
	defer closeFun()
	d := &share_message.WishBoxItem{}
	e := col.Find(bson.M{"_id": id}).One(&d)
	if e != nil && e.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", e.Error())
	}
	return d, e
}

//添加盲盒挑战记录
func AddWishLog(data *share_message.WishLog) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_LOG)
	defer closeFun()
	data.Id = easygo.NewInt64(NextId(TABLE_WISH_LOG))
	data.CreateTime = easygo.NewInt64(time.Now().Unix())
	err := col.Insert(data)
	return err
}

//获取盲盒挑战列表
func GetWishLogList(wishBoxId int64, page, pageSize int) ([]*share_message.WishLog, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_LOG)
	defer closeFun()
	list := make([]*share_message.WishLog, 0, pageSize)
	whereBson := bson.M{"WishBoxId": wishBoxId}
	col.Find(whereBson).Sort("-_id").Skip((page - 1) * pageSize).Limit(pageSize).All(&list)
	count, _ := col.Find(whereBson).Count()
	return list, count
}

//添加盲盒挑战占领时长表
func AddOccupied(occupied *share_message.WishOccupied) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFun()
	occupied.Id = easygo.NewInt64(NextId(TABLE_WISH_OCCUPIED))
	occupied.CreateTime = easygo.NewInt64(time.Now().Unix())
	err := col.Insert(occupied)
	return err
}

//设置挑战者结束时间
func UpOccupied(wishBoxId, playerId int64) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFun()
	occupiedInfo, _ := GetOccupied(wishBoxId, playerId)
	occupiedTime := time.Now().Unix() - occupiedInfo.GetCreateTime()
	// 判断时间是否超过一天,如果是的话,不能处理了
	logs.Info("相加后的时间为:", occupiedInfo.GetCreateTime()+AFTER_TIME_GUARDIAN)
	logs.Info("当前时间为:", time.Now().Unix())
	if occupiedInfo.GetCreateTime()+AFTER_TIME_GUARDIAN < time.Now().Unix() {
		//return errors.New("该守护者超过24小时了,挑战不需要去修改了")
		occupiedTime = AFTER_TIME_GUARDIAN
	}
	occupied := &share_message.WishOccupied{
		EndTime:      easygo.NewInt64(time.Now().Unix()),
		OccupiedTime: easygo.NewInt64(occupiedTime),
		Status:       easygo.NewInt32(2),
	}
	err := col.Update(bson.M{"PlayerId": playerId, "WishBoxId": wishBoxId, "Status": 1}, bson.M{"$set": occupied})
	if err == nil {
		// 同步总的占领时间
		oc, _ := GetOccupiedNoStauts(wishBoxId, playerId)
		data := &share_message.WishSumOccupied{
			WishBoxId: easygo.NewInt64(wishBoxId),
			NickName:  easygo.NewString(oc.GetNickName()),
			HeadUrl:   easygo.NewString(oc.GetHeadUrl()),
			PlayerId:  easygo.NewInt64(playerId),
		}
		UpdateWishSumOccupied(data, occupiedTime)
	}
	return err
}
func UpOccupied1(wishBoxId, playerId int64) {

	occupiedInfo, _ := GetOccupied(wishBoxId, playerId)
	occupiedTime := time.Now().Unix() - occupiedInfo.GetCreateTime()
	// 判断时间是否超过一天,如果是的话,不能处理了
	logs.Info("相加后的时间为:", occupiedInfo.GetCreateTime()+AFTER_TIME_GUARDIAN)
	logs.Info("当前时间为:", time.Now().Unix())
	if occupiedInfo.GetCreateTime()+AFTER_TIME_GUARDIAN < time.Now().Unix() {
		//return errors.New("该守护者超过24小时了,挑战不需要去修改了")
		occupiedTime = AFTER_TIME_GUARDIAN
	}

	logs.Info("最后计算的占领时间occupiedTime---------->%d秒", occupiedTime)
	// 同步总的占领时间
	oc, _ := GetOccupiedNoStauts(wishBoxId, playerId)
	data := &share_message.WishSumOccupied{
		WishBoxId:    easygo.NewInt64(wishBoxId),
		NickName:     easygo.NewString(oc.GetNickName()),
		HeadUrl:      easygo.NewString(oc.GetHeadUrl()),
		PlayerId:     easygo.NewInt64(playerId),
		OccupiedTime: easygo.NewInt64(occupiedTime),
	}
	UpdateWishSumOccupied1(data)
}

//设置挑战者结束时间
func UpOccupiedEx(wishBoxId, playerId int64, t int64) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFun()

	occupied := &share_message.WishOccupied{
		EndTime:      easygo.NewInt64(time.Now().Unix()),
		OccupiedTime: easygo.NewInt64(t),
		Status:       easygo.NewInt32(2),
	}
	err := col.Update(bson.M{"PlayerId": playerId, "WishBoxId": wishBoxId, "Status": 1}, bson.M{"$set": occupied})
	if err == nil {
		// 同步总的占领时间
		oc, _ := GetOccupiedNoStauts(wishBoxId, playerId)
		data := &share_message.WishSumOccupied{
			WishBoxId: easygo.NewInt64(wishBoxId),
			NickName:  easygo.NewString(oc.GetNickName()),
			HeadUrl:   easygo.NewString(oc.GetHeadUrl()),
			PlayerId:  easygo.NewInt64(playerId),
		}
		UpdateWishSumOccupied(data, t)
	}
	return err
}

//获取盲盒挑战占领时长列表
func GetWishOccupiedList(wishBoxId int64, page, pageSize int) ([]*share_message.WishOccupied, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFun()
	list := make([]*share_message.WishOccupied, 0, pageSize)
	whereBson := bson.M{"WishBoxId": wishBoxId}
	col.Find(whereBson).Sort("-_id").Skip((page - 1) * pageSize).Limit(pageSize).All(&list)
	count, _ := col.Find(whereBson).Count()
	return list, count
}

func GetOccupied(wishBoxId, playerId int64) (*share_message.WishOccupied, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFun()
	d := &share_message.WishOccupied{}
	e := col.Find(bson.M{"PlayerId": playerId, "WishBoxId": wishBoxId, "Status": 1}).One(&d)
	return d, e
}
func GetOccupiedNoStauts(wishBoxId, playerId int64) (*share_message.WishOccupied, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFun()
	d := &share_message.WishOccupied{}
	e := col.Find(bson.M{"PlayerId": playerId, "WishBoxId": wishBoxId}).One(&d)
	return d, e
}

//更新守护者期间获得的硬币数
func UpOccupiedCoin(wishBoxId, playerId int64, diamondLimit, addDiamondNum int32) (int32, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFun()
	occupiedInfo, _ := GetOccupied(wishBoxId, playerId)
	if occupiedInfo == nil {
		return 0, errors.New("没有守护者")
	}
	if occupiedInfo.GetCoinNum() >= diamondLimit {
		logs.Error("当日收益钻石已达上限,上限值为: %d", diamondLimit)
		return 0, errors.New("当日收益钻石已达上限")
	}
	if occupiedInfo.GetCreateTime()+24*3600 < time.Now().Unix() {
		return 0, errors.New("守护者有效时间为24小时")
	}
	addNum := addDiamondNum
	if occupiedInfo.GetCoinNum()+addDiamondNum >= diamondLimit {
		addNum = diamondLimit - occupiedInfo.GetCoinNum()
	}
	logs.Info("玩家: %d 在盲盒id为: %d中成为守护者获得的硬币为: %d,接下来需要奖励的硬币为: %d", playerId, wishBoxId, occupiedInfo.GetCoinNum(), addNum)
	err := col.Update(bson.M{"PlayerId": playerId, "WishBoxId": wishBoxId, "Status": 1}, bson.M{"$inc": bson.M{"CoinNum": addNum}})
	if err != nil {
		return 0, err
	}
	return addNum, nil

}

//添加玩家物品
func AddPlayerWishItem(playerWishItem *share_message.PlayerWishItem) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_ITEM)
	defer closeFun()
	curTime := time.Now().Unix()

	playerWishItem.CreateTime = easygo.NewInt64(curTime)

	err := col.Insert(playerWishItem)
	return err
}
func AddToolPlayerWishItem(playerWishItem *share_message.PlayerWishItem) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOOL_PLAYER_WISH_ITEM)
	defer closeFun()
	curTime := time.Now().Unix()

	playerWishItem.CreateTime = easygo.NewInt64(curTime)

	err := col.Insert(playerWishItem)
	return err
}

// 添加守护者获得的钻石流水
func AddWishGuardianDiamondLog(data *share_message.WishGuardianDiamondLog) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_GUARDIAN_DIAMOND_LOG)
	defer closeFun()
	err := col.Insert(data)
	return err
}

func GetAllWishLogByPid(pid int64) []*share_message.WishLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_LOG)
	defer closeFun()
	wishLogs := make([]*share_message.WishLog, 0)
	err := col.Find(bson.M{"DareId": pid, "Result": true}).All(&wishLogs)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return wishLogs
}
func GetAllWishLogByPage(pid int64, page, pageSize int) ([]*share_message.WishLog, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_LOG)
	defer closeFun()
	wishLogs := make([]*share_message.WishLog, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)

	query := col.Find(bson.M{"DareId": pid, "Result": true})
	count, err := query.Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	err = query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&wishLogs)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return wishLogs, count

}

// 根据商品中间表的ids查询中间表信息
func GetWishBoxsByIdsFromDB(ids []int64) []*share_message.WishBox {
	var lst []*share_message.WishBox
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFunc()
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return lst
}

func GetAllWishItem() []*share_message.WishItem {
	var lst []*share_message.WishItem
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ITEM)
	defer closeFunc()
	err := col.Find(bson.M{}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return lst
}

// 随机获取固定条数的推荐类别
func GetRangeRecommendWishItem(num int) []*share_message.WishItemType {
	var lst []*share_message.WishItemType
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ITEM_TYPE)
	defer closeFunc()
	queryBson := bson.M{"IsRecommend": true}
	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return lst
}

func GetAllWishBox() []*share_message.WishBox {
	var lst []*share_message.WishBox
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFunc()
	err := col.Find(bson.M{}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return lst
}
func GetWishBoxItemByItemId(itemId int64) []*share_message.WishBoxItem {
	var lst []*share_message.WishBoxItem
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFunc()
	err := col.Find(bson.M{"WishItemId": itemId}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return lst
}

func GetBoxItemListBySaleStatus(saleStatus int32) []*share_message.WishBoxItem {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFun()
	wishBoxItems := make([]*share_message.WishBoxItem, 0)
	err := col.Find(bson.M{"Status": saleStatus}).Sort("-CreateTime").All(&wishBoxItems)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return wishBoxItems

}

func GetBoxListByPage(req *h5_wish.SearchBoxReq, page, pageSize int) ([]*share_message.WishBox, int) {
	wishBoxes := make([]*share_message.WishBox, 0)
	boxIds := make([]int64, 0)
	if req.GetProductStatus() != 0 {
		boxes := GetBoxItemListBySaleStatus(req.GetProductStatus())
		for _, v := range boxes {
			if !util.Int64InSlice(v.GetWishBoxId(), boxIds) {
				boxIds = append(boxIds, v.GetWishBoxId())
			}
		}
	}

	if req.GetProductStatus() != 0 && len(boxIds) == 0 { // 有状态,但是没有盲盒id列表,直接返回
		return wishBoxes, 0
	}

	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFun()
	queryBson := bson.M{"Status": bson.M{"$gt": WISH_DOWN_STATUS}, "items.0": bson.M{"$exists": 1}} // 上架
	curPage := easygo.If(page > 1, page-1, 0).(int)
	if len(boxIds) > 0 {
		queryBson = bson.M{"_id": bson.M{"$in": boxIds}}
	}
	if req.GetMinPrice() != 0 {
		queryBson["Price"] = bson.M{"$gte": req.GetMinPrice()}
	}
	if req.GetMaxPrice() != 0 {
		queryBson["Price"] = bson.M{"$lte": req.GetMaxPrice()}
	}

	if len(req.GetWishBrandId()) > 0 {
		queryBson["Brands"] = bson.M{"$in": req.GetWishBrandId()}
	}
	if len(req.GetWishItemTypeId()) > 0 {
		queryBson["Types"] = bson.M{"$in": req.GetWishItemTypeId()}
	}

	if req.GetLabel() == 1 { // 普通
		queryBson["Match"] = 0 // 0普通赛，1挑战赛
	} else if req.GetLabel() == 2 { // 挑战
		queryBson["Match"] = 1 // 0普通赛，1挑战赛
	}

	if req.GetCondition() != 0 && req.GetCondition() != 1 { // 1-不是最新上线
		queryBson["Menu"] = req.GetCondition()
	}
	query := col.Find(queryBson)
	count, err := query.Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	complex := req.GetComplex()
	switch complex {
	case 0: // 综合
		err = query.Sort("-SortWeight", "-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&wishBoxes)
	case 1: //1-价格降序
		err = query.Sort("-Price").Skip(curPage * pageSize).Limit(pageSize).All(&wishBoxes)
	case 2: // 2-价格升序
		err = query.Sort("Price").Skip(curPage * pageSize).Limit(pageSize).All(&wishBoxes)
	}

	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return wishBoxes, count
}

func GetRangeWishOccupied(num int) []*share_message.WishOccupied {
	var lst []*share_message.WishOccupied
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFunc()
	queryBson := bson.M{"Status": WISH_OCCUPIED_STATUS_UP} // 占领中
	queryBson["CoinNum"] = bson.M{"$gt": 0}
	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return lst
}

func GetSomeWishOccupied(num int) []*share_message.WishOccupied {
	var lst []*share_message.WishOccupied
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFunc()
	m := []bson.M{
		{"$match": bson.M{"CoinNum": bson.M{"$gt": 0}}},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return lst
}

func GetWishOccupiedByBoxIds(boxIds []int64) []*share_message.WishOccupied {
	var lst []*share_message.WishOccupied
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFunc()
	queryBson := bson.M{"WishBoxId": bson.M{"$in": boxIds}, "Status": WISH_OCCUPIED_STATUS_UP} // 占领中

	err := col.Find(queryBson).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return lst
}
func GetAllWishOccupiedByPid(pid int64) []*share_message.WishOccupied {
	lst := make([]*share_message.WishOccupied, 0)
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFunc()
	err := col.Find(bson.M{"PlayerId": pid}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return lst
}
func GetAllWishOccupiedByPage(pid int64, page, pageSize int) ([]*share_message.WishOccupied, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFun()
	wishLogs := make([]*share_message.WishOccupied, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)

	query := col.Find(bson.M{"PlayerId": pid})
	count, err := query.Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	err = query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&wishLogs)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return wishLogs, count

}

// 榜单 WishTopLog
//func UpsetWishTopToDB(data *share_message.WishTopLog, wishNum, coinNum int64) {
//	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_TOP_LOG)
//	defer closeFunc()
//	_, err := col.Upsert(bson.M{"_id": data.GetId()}, bson.M{"$set": data, "$inc": bson.M{"WishNum": wishNum, "CoinNum": coinNum}})
//	easygo.PanicError(err)
//}

// 获取10条榜单数据
//func GetTop10FromDB(thatDayTime int64) []*share_message.WishTopLog {
//	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_TOP_LOG)
//	defer closeFunc()
//
//	list := make([]*share_message.WishTopLog, 0)
//	err := col.Find(bson.M{"ThatDayTime": thatDayTime}).Sort("-CoinNum").Limit(10).All(&list)
//	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
//		logs.Error("dal_wish err: %s", err.Error())
//	}
//	return list
//
//}
// 获取10条榜单数据
func GetTop10FromDB() []*share_message.WishGuardianTopLog {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_GUARDIAN_TOP_LOG)
	defer closeFunc()

	list := make([]*share_message.WishGuardianTopLog, 0)
	err := col.Find(bson.M{}).Sort("-CoinNum").Limit(10).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return list

}

// 获取 水池信息
func GetPoolInfoFromDB(poolId int64) *share_message.WishPool {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_POOL)
	defer closeFunc()
	var result *share_message.WishPool
	err := col.Find(bson.M{"_id": poolId}).One(&result)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return result
}

// 根据水池状态获取盲盒中物品信息
func GetProductByPoolStatus(boxId int64, status int32, isOpenAward bool) []*share_message.WishBoxItem {
	list := make([]*share_message.WishBoxItem, 0)
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFunc()
	queryBson := bson.M{}
	if isOpenAward {
		// 如果是放奖,只有普通,大赢,小盈才有数据的
		if status != POOL_STATUS_COMMON && status != POOL_STATUS_BIGWIN && status != POOL_STATUS_SMALLWIN {
			logs.Error("开奖状态,水池状态有误,boxId: %d,当前状态为: %d", boxId, status)
			return list
		}
		switch status {
		case POOL_STATUS_COMMON: // 普通
			queryBson = bson.M{"CommonAddWeight": bson.M{"$gt": 0}}
		case POOL_STATUS_BIGWIN: // 大盈
			queryBson = bson.M{"BigWinAddWeight": bson.M{"$gt": 0}}
		case POOL_STATUS_SMALLWIN: // 小盈
			queryBson = bson.M{"SmallWinAddWeight": bson.M{"$gt": 0}}
		default:
			logs.Error("水池状态有误,当前状态为: %d", status)

		}
	} else {
		switch status {
		case POOL_STATUS_BIGLOSS: // 大亏
			queryBson = bson.M{"BigLoss": bson.M{"$gt": 0}}
		case POOL_STATUS_SMALLLOSS: // 小亏
			queryBson = bson.M{"SmallLoss": bson.M{"$gt": 0}}
		case POOL_STATUS_COMMON: // 普通
			queryBson = bson.M{"Common": bson.M{"$gt": 0}}
		case POOL_STATUS_BIGWIN: // 大盈
			queryBson = bson.M{"BigWin": bson.M{"$gt": 0}}
		case POOL_STATUS_SMALLWIN: // 小盈
			queryBson = bson.M{"SmallWin": bson.M{"$gt": 0}}
		default:
			logs.Error("水池状态有误,当前状态为: %d", status)
		}
	}

	queryBson["WishBoxId"] = boxId
	err := col.Find(queryBson).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return list
}

func GetToolProductByPoolStatus(boxId int64, status int32, isOpenAward bool) []*share_message.WishBoxItem {
	list := make([]*share_message.WishBoxItem, 0)
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOOL_WISH_BOX_ITEM)
	defer closeFunc()
	queryBson := bson.M{}
	if isOpenAward {
		// 如果是放奖,只有普通,大赢,小盈才有数据的
		if status != POOL_STATUS_COMMON && status != POOL_STATUS_BIGWIN && status != POOL_STATUS_SMALLWIN {
			logs.Error("后台抽奖工具 开奖状态,水池状态有误, 当前状态为: %d", status)
			return list
		}
		switch status {
		case POOL_STATUS_COMMON: // 普通
			queryBson = bson.M{"CommonAddWeight": bson.M{"$gt": 0}}
		case POOL_STATUS_BIGWIN: // 大盈
			queryBson = bson.M{"BigWinAddWeight": bson.M{"$gt": 0}}
		case POOL_STATUS_SMALLWIN: // 小盈
			queryBson = bson.M{"SmallWinAddWeight": bson.M{"$gt": 0}}
		default:
			logs.Error("后台抽奖工具 水池状态有误,当前状态为: %d", status)

		}
	} else {
		switch status {
		case POOL_STATUS_BIGLOSS: // 大亏
			queryBson = bson.M{"BigLoss": bson.M{"$gt": 0}}
		case POOL_STATUS_SMALLLOSS: // 小亏
			queryBson = bson.M{"SmallLoss": bson.M{"$gt": 0}}
		case POOL_STATUS_COMMON: // 普通
			queryBson = bson.M{"Common": bson.M{"$gt": 0}}
		case POOL_STATUS_BIGWIN: // 大盈
			queryBson = bson.M{"BigWin": bson.M{"$gt": 0}}
		case POOL_STATUS_SMALLWIN: // 小盈
			queryBson = bson.M{"SmallWin": bson.M{"$gt": 0}}
		default:
			logs.Error("后台抽奖工具 水池状态有误,当前状态为: %d", status)
		}
	}
	queryBson["WishBoxId"] = boxId
	err := col.Find(queryBson).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("后台抽奖工具 dal_wish err: %s", err.Error())
	}
	return list
}

//添加水池抽水日志
func AddWishPoolPumpLog(data *share_message.WishPoolPumpLog) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_POOL_PUMP_LOG)
	defer closeFun()
	data.Id = easygo.NewInt64(NextId(TABLE_WISH_POOL_PUMP_LOG))
	data.CreateTime = easygo.NewInt64(time.Now().Unix())
	err := col.Insert(data)
	easygo.PanicError(err)
}
func AddToolWishPoolPumpLog(data *share_message.WishPoolPumpLog) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOOL_WISH_POOL_PUMP_LOG)
	defer closeFun()
	data.Id = easygo.NewInt64(NextId(TABLE_TOOL_WISH_POOL_PUMP_LOG))
	data.CreateTime = easygo.NewInt64(time.Now().Unix())
	err := col.Insert(data)
	easygo.PanicError(err)
}

//水池流水日志
func AddWishPoolLog(data *share_message.WishPoolLog) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_POOL_LOG)
	defer closeFun()
	data.Id = easygo.NewInt64(NextId(TABLE_WISH_POOL_LOG))
	data.CreateTime = easygo.NewInt64(time.Now().Unix())
	err := col.Insert(data)
	easygo.PanicError(err)
}
func AddToolWishPoolLog(data *share_message.WishPoolLog) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_TOOL_WISH_POOL_LOG)
	defer closeFun()
	data.Id = easygo.NewInt64(NextId(TABLE_TOOL_WISH_POOL_LOG))
	data.CreateTime = easygo.NewInt64(time.Now().Unix())
	err := col.Insert(data)
	easygo.PanicError(err)
}

//根据盲盒ID获取挑战者信息
func GetBoxDefenderByBoxId(boxId int64) *share_message.WishOccupied {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OCCUPIED)
	defer closeFun()
	var defender *share_message.WishOccupied
	err := col.Find(bson.M{"WishBoxId": boxId, "Status": 1}).Sort("-_id").One(&defender)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return defender
}

//根据盲盒ID获取商品展示区列表
func GetWishBoxItemByBoxId(boxId int64) []*share_message.WishBoxItem {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFun()
	var boxItem []*share_message.WishBoxItem
	err := col.Find(bson.M{"WishBoxId": boxId}).All(&boxItem)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return boxItem
}

// 新增水池数据
func AddWishPool(data *share_message.WishPool) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_POOL)
	defer closeFun()
	data.Id = easygo.NewInt64(NextId(TABLE_WISH_POOL))
	data.CreateTime = easygo.NewInt64(time.Now().Unix())
	err := col.Insert(data)
	easygo.PanicError(err)
}

// 删除许愿的数据.
func DelWishData(id int64) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_DATA)
	defer closeFun()
	err := col.Remove(bson.M{"_id": id})
	easygo.PanicError(err)
}

// 获取抽奖玩家信息
func GetWishPlayerByPid(pid int64) *share_message.WishPlayer {
	col, closeFunItem := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFunItem()
	wp := &share_message.WishPlayer{}
	err := col.Find(bson.M{"_id": pid}).One(&wp)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return wp
}

// 获取抽奖玩家信息
func GetWishPlayerActivityByPid(pid int64) *share_message.WishPlayerActivity {
	col, closeFunItem := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER_ACTIVITY)
	defer closeFunItem()
	wp := &share_message.WishPlayerActivity{}
	err := col.Find(bson.M{"_id": pid}).One(&wp)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return wp
}

// 获取玩家信息
func GetWishPlayerByAccount(channel int32, account string) *share_message.WishPlayer {
	col, closeFunItem := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFunItem()
	var wp *share_message.WishPlayer
	err := col.Find(bson.M{"Account": account, "Channel": channel}).One(&wp)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return wp
}

func UpsertBox(bId int64, data *share_message.WishBox) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": bId}, data)
	easygo.PanicError(err)
	return err
}

// 获取冷却配置表
func GetWishCoolDownConfigFromDB() *share_message.WishCoolDownConfig {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OTHER_CFG)
	defer closeFunc()
	var result *share_message.WishCoolDownConfig
	err := col.Find(bson.M{"_id": TABLE_WISH_COOL_DOWN_CONFIG}).One(&result)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return result
}
func AddWishCoolDownConfigFromDB(data *share_message.WishCoolDownConfig) {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OTHER_CFG)
	defer closeFunc()
	data.Id = easygo.NewString(TABLE_WISH_COOL_DOWN_CONFIG)
	data.CreateTime = easygo.NewInt64(time.Now().Unix())
	err := col.Insert(data)
	easygo.PanicError(err)
}

// 获取守护者收益设置
func GetWishGuardianCfg() *share_message.WishGuardianCfg {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OTHER_CFG)
	defer closeFun()

	var cfg *share_message.WishGuardianCfg
	queryBson := bson.M{"_id": TABLE_WISH_GUARDIAN_CFG}
	err := col.Find(queryBson).One(&cfg)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}

	return cfg
}

// 更新物品回收参数设置
func UpdateWishRecycleSection(list *share_message.WishRecycleSection) {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OTHER_CFG)
	defer closeFunc()
	data := &share_message.WishRecycleSection{
		Id:       easygo.NewString(TABLE_WISH_RECYCLE_SECTION),
		Player:   easygo.NewInt32(list.GetPlayer()),
		Platform: easygo.NewInt32(list.GetPlatform()),
	}

	_, err := col.Upsert(bson.M{"_id": data.GetId()}, bson.M{"$set": data})
	easygo.PanicError(err)
}

func UpdateTaskTime(id int64, taskTime int64) error {
	col, closeFunItem := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFunItem()
	err := col.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"TaskTime": taskTime}})
	return err
}

func UpBoxItemLocalNum(id int64) (int32, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX_ITEM)
	defer closeFun()
	var item *share_message.WishBoxItem
	err := col.Find(bson.M{"_id": id, "PerNum": bson.M{"$gt": 0}}).One(&item)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		return 0, err
	}
	if item == nil {
		return 0, errors.New("not fond")
	}

	identity := &share_message.WishBoxItem{}
	_, err = col.Find(bson.M{"_id": id}).Apply(mgo.Change{
		Update:    bson.M{"$inc": bson.M{"LocalNum": -1}},
		Upsert:    true,
		ReturnNew: true,
	}, &identity)
	if err != nil {
		return 0, err
	}

	result := identity.GetLocalNum()
	return result, nil
}

// 获取积极补货中定时任务中的盲盒
func GetTaskBox() []*share_message.WishBox {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFun()
	list := make([]*share_message.WishBox, 0)
	err := col.Find(bson.M{"Status": 2, "IsTask": true}).All(&list) // 积极补货中
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return list
}

// 获取兑换比例配置
func GetCurrencyCfg() *share_message.WishCurrencyConversionCfg {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OTHER_CFG)
	defer closeFun()

	var cfg *share_message.WishCurrencyConversionCfg
	queryBson := bson.M{"_id": TABLE_WISH_CURRENCY_CONVERSION_CFG}
	err := col.Find(queryBson).One(&cfg)
	if err != nil && err != mgo.ErrNotFound {
		logs.Error("查询兑换配置失败,err: %s", err.Error())
		return nil
	}

	return cfg
}

// 根据商品中间表的ids查询中间表信息
func GetIsGuardianWishBoxList() []*share_message.WishBox {
	var lst []*share_message.WishBox
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFunc()
	err := col.Find(bson.M{"IsGuardian": true, "Status": 1, "GuardianId": bson.M{"$gt": 0}}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return lst
}

// 获取砖石兑换列表
func GetDiamondRechargeList() []*share_message.DiamondRecharge {
	lst := make([]*share_message.DiamondRecharge, 0)
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_DIAMOND_RECHARGE)
	defer closeFunc()
	err := col.Find(bson.M{"Status": 1}).Sort("-Sort").All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return lst
}
func InsertDiamondRecharge(data *share_message.DiamondRecharge) error {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_DIAMOND_RECHARGE)
	defer closeFunc()
	return col.Insert(data)
}

func GetDiamondRechargeById(id int64) *share_message.DiamondRecharge {
	var data *share_message.DiamondRecharge
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_DIAMOND_RECHARGE)
	defer closeFunc()
	err := col.Find(bson.M{"_id": id, "Status": 1}).One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
		return data
	}
	return data
}

// 随机抽取玩家
func GetRandWishPlayer(num int) []*share_message.WishPlayer {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFun()
	list := make([]*share_message.WishPlayer, 0)

	queryBson := bson.M{}
	queryBson = bson.M{"Types": WISH_PLAYER_TYPE_ROBOT, "HeadIcon": bson.M{"$ne": ""}} //有效的运营号

	m := []bson.M{
		{"$match": queryBson},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)
	return list
}

func GetDiamondChangeLogByPage(uid int64, t int32, page, pageSize int) ([]*DiamondLog, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_DIAMOND_CHANGELOG)
	defer closeFun()
	var pageList []*DiamondLog
	curPage := easygo.If(page > 1, page-1, 0).(int)
	queryBson := bson.M{"PlayerId": uid, "PayType": t}
	query := col.Find(queryBson)
	count, err := query.Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	err = query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&pageList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}

	return pageList, count
}

// 获取价格区间参数
func GetPriceSection() *share_message.PriceSection {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_OTHER_CFG)
	defer closeFunc()

	data := &share_message.PriceSection{}
	queryBson := bson.M{"_id": TABLE_WISH_PRICE_SECTION}
	err := col.Find(queryBson).One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetPriceSection err: %v", err.Error())
	}

	return data
}

// 找到推荐的盲盒
func GetRecommendBox() []*share_message.WishBox {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFun()
	list := make([]*share_message.WishBox, 0)
	err := col.Find(bson.M{"Status": bson.M{"$ne": 0}, "IsRecommend": true}).Sort("-SortWeight").All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("dal_wish err: %s", err.Error())
	}
	return list
}

func GetWishDataByPid(playerId int64, status int32, boxIds []int64) ([]*share_message.PlayerWishData, error) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_PLAYER_WISH_DATA)
	defer closeFun()
	result := make([]*share_message.PlayerWishData, 0)
	e := col.Find(bson.M{"PlayerId": playerId, "Status": status, "WishBoxId": bson.M{"$in": boxIds}}).Sort("-CreateTime").All(&result)
	if e != nil && e.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishDataByPid err: %s", e.Error())
	}
	return result, e
}

//设置当天挑战者挑战的成功次数
func SetWishBoxWishNum(id int64) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFun()
	return col.Update(bson.M{"_id": id}, bson.M{"$inc": bson.M{"WinNum": 1}})
}

func SetBrandClickNum(id int64) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BRAND)
	defer closeFun()
	return col.Update(bson.M{"_id": id}, bson.M{"$inc": bson.M{"ClickCount": 1}})
}
func SetTypeClickNum(id int64) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ITEM_TYPE)
	defer closeFun()
	return col.Update(bson.M{"_id": id}, bson.M{"$inc": bson.M{"ClickCount": 1}})
}

// 获取排行榜的数据列表
func GetGuardianCoinNumList(st, ed int64) []bson.M {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_GUARDIAN_DIAMOND_LOG)
	defer closeFun()

	queryBson := []bson.M{
		{"$match": bson.M{"CoinNum": bson.M{"$gt": 0}, "$and": []bson.M{{"CreateTime": bson.M{"$gt": st}}, {"CreateTime": bson.M{"$lt": ed}}}}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": "$CoinNum"}}},
	}
	rst := make([]bson.M, 0)
	err := col.Pipe(queryBson).All(&rst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	return rst
}

// 获取排行榜的数据列表
func GetGuardianWishNumList(st, ed int64) []bson.M {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_GUARDIAN_DIAMOND_LOG)
	defer closeFun()

	queryBson := []bson.M{
		{"$match": bson.M{"WishNum": bson.M{"$gt": 0}, "$and": []bson.M{{"CreateTime": bson.M{"$gt": st}}, {"CreateTime": bson.M{"$lt": ed}}}}},
		{"$group": bson.M{"_id": "$PlayerId", "Count": bson.M{"$sum": "$WishNum"}}},
	}
	rst := make([]bson.M, 0)
	err := col.Pipe(queryBson).All(&rst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	return rst
}

func GetLastCoinTopLog() *share_message.WishGuardianTopLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_GUARDIAN_TOP_LOG)
	defer closeFun()
	var data *share_message.WishGuardianTopLog
	err := col.Find(bson.M{}).Sort("-UpdateCoinTime").One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return data
}
func GetLastWishTopLog() *share_message.WishGuardianTopLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_GUARDIAN_TOP_LOG)
	defer closeFun()
	var data *share_message.WishGuardianTopLog
	err := col.Find(bson.M{}).Sort("-UpdateWishTime").One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return data
}

func GetGuardianWinNumList(t int64, status int64) []*share_message.WishGuardianDiamondLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_GUARDIAN_DIAMOND_LOG)
	defer closeFun()
	dLogList := make([]*share_message.WishGuardianDiamondLog, 0)
	err := col.Find(bson.M{"Status": status, "CreateTime": bson.M{"$lt": t}, "WishNum": bson.M{"$gt": 0}}).All(&dLogList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return dLogList
}

func UpdateWishGuardianTopLog(data *share_message.WishGuardianTopLog, coinNum, wishNum int64) {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_GUARDIAN_TOP_LOG)
	defer closeFunc()

	_, err := col.Upsert(bson.M{"_id": data.GetId()}, bson.M{"$set": data, "$inc": bson.M{"WishNum": wishNum, "CoinNum": coinNum}})
	easygo.PanicError(err)
}

// 统计总的占领时长
func UpdateWishSumOccupied(data *share_message.WishSumOccupied, occupiedTime int64) {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_SUM_OCCUPIED)
	defer closeFunc()

	_, err := col.Upsert(bson.M{"WishBoxId": data.GetWishBoxId(), "PlayerId": data.GetPlayerId()}, bson.M{"$set": data, "$inc": bson.M{"OccupiedTime": occupiedTime}})
	easygo.PanicError(err)
}
func UpdateWishSumOccupied1(data *share_message.WishSumOccupied) {
	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_SUM_OCCUPIED)
	defer closeFunc()

	_, err := col.Upsert(bson.M{"WishBoxId": data.GetWishBoxId(), "PlayerId": data.GetPlayerId()}, bson.M{"$set": data})
	easygo.PanicError(err)
}

func GetWishSumOccupied(wishBoxId int64, page, pageSize int) ([]*share_message.WishSumOccupied, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_SUM_OCCUPIED)
	defer closeFun()
	list := make([]*share_message.WishSumOccupied, 0, pageSize)
	whereBson := bson.M{"WishBoxId": wishBoxId}
	col.Find(whereBson).Sort("-OccupiedTime").Skip((page - 1) * pageSize).Limit(pageSize).All(&list)
	count, _ := col.Find(whereBson).Count()
	return list, count
}

//func UpsetWishTopToDB(data *share_message.WishTopLog, wishNum, coinNum int64) {
//	col, closeFunc := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_TOP_LOG)
//	defer closeFunc()
//	_, err := col.Upsert(bson.M{"_id": data.GetId()}, bson.M{"$set": data, "$inc": bson.M{"WishNum": wishNum, "CoinNum": coinNum}})
//	easygo.PanicError(err)
//}

// 随机抽取上架的盲盒
func GetRandWishBox(num int) []*share_message.WishBox {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_BOX)
	defer closeFun()
	list := make([]*share_message.WishBox, 0)

	m := []bson.M{
		{"$match": bson.M{"Status": 1, "Match": 1}},
		{"$sample": bson.M{"size": num}},
	}
	query := col.Pipe(m)
	err := query.All(&list)
	easygo.PanicError(err)
	return list
}

// 写付费用户地理位置分布日志
func AddPayPlayerLocationLog(id PLAYER_ID, ip string) {
	logs.Info("========写付费用户地理位置分布日志AddPayPlayerLocationLog===========,id: %d,ip: %v", id, ip)
	pmgr := GetRedisPlayerBase(id)
	if pmgr == nil {
		return
	}
	dayTime := easygo.GetToday0ClockTimestamp()
	dType := pmgr.GetDeviceType()

	if ip == "" {
		return
	}
	data := IpSearch(ip)
	if data == nil {
		return
	}
	var position, piece string = data.Region, data.CountryId

	logid := easygo.AnytoA(dayTime) + easygo.AnytoA(id) + easygo.AnytoA(dType) + position
	log := &share_message.PayPlayerLocationLog{
		Id:         easygo.NewString(logid),
		DayTime:    easygo.NewInt64(dayTime),
		Position:   easygo.NewString(position),
		Piece:      easygo.NewString(piece),
		DeviceType: easygo.NewInt32(dType),
		PlayerId:   easygo.NewInt64(id),
	}

	FindAndModify(MONGODB_NINGMENG_LOG, TABLE_WISH_PAYPLAYER_LOCATION_LOG, bson.M{"_id": log.GetId()}, log, true)
}

// 根据盲盒id找到对应的奖池信息
func GetWishActPoolByBoxId(boxId int64) *share_message.WishActPool {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACT_POOL)
	defer closeFun()
	var actPool *share_message.WishActPool
	err := col.Find(bson.M{"BoxIds": boxId}).One(&actPool)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishActPoolByBoxId err: %s", err.Error())
		return actPool
	}
	return actPool
}
func GeWishActPoolRuleByPId(pId int64) []*share_message.WishActPoolRule {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACT_POOL_RULE)
	defer closeFun()
	poolRule := make([]*share_message.WishActPoolRule, 0)
	err := col.Find(bson.M{"WishActPoolId": pId}).All(&poolRule)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GeWishActPoolRuleByPId err: %s", err.Error())
		return poolRule
	}
	return poolRule
}
func InsertWishActPool(data *share_message.WishActPool) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACT_POOL)
	defer closeFun()
	return col.Insert(data)
}
func InsertWishActPoolRule(data *share_message.WishActPoolRule) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACT_POOL_RULE)
	defer closeFun()
	return col.Insert(data)
}

// 玩家活动数据表
func GeWishPlayerActivityByPId(pId int64) *share_message.WishPlayerActivity {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER_ACTIVITY)
	defer closeFun()
	var activity *share_message.WishPlayerActivity
	err := col.Find(bson.M{"_id": pId}).One(&activity)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GeWishActPoolRuleByPId err: %s", err.Error())
		return activity
	}
	return activity
}

func UpsertWishPlayerActivity(data *share_message.WishPlayerActivity) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER_ACTIVITY)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": data.GetPlayerId()}, data)
	easygo.PanicError(err)
	return err
}

func GeWishActivityPrizeLogByPId(pId, wishActPoolId int64, dayLun ...int32) *share_message.WishActivityPrizeLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACTIVITY_PRIZE_LOG)
	defer closeFun()
	var activityPrizeLog *share_message.WishActivityPrizeLog
	q := bson.M{"PlayerId": pId, "WishActPoolRuleId": wishActPoolId}
	if len(dayLun) > 0 {
		q["DayLun"] = dayLun[0]
	}
	err := col.Find(q).One(&activityPrizeLog)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GeWishActivityPrizeLogByPId err: %s", err.Error())
		return activityPrizeLog
	}
	return activityPrizeLog
}

func UpsertWishActivityPrizeLog(data *share_message.WishActivityPrizeLog) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACTIVITY_PRIZE_LOG)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": data.GetId()}, data)
	easygo.PanicError(err)
	return err
}

func GetWishPlayerActivityList() []*share_message.WishPlayerActivity {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER_ACTIVITY)
	defer closeFun()
	data := make([]*share_message.WishPlayerActivity, 0)
	if err := col.Find(bson.M{}).All(&data); err != nil {
		easygo.PanicError(err)
	}

	return data
}

func GetWishActPoolList() []*share_message.WishActPool {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACT_POOL)
	defer closeFun()
	list := make([]*share_message.WishActPool, 0)
	err := col.Find(bson.M{}).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishActPoolList err: %s", err.Error())
		return list
	}
	return list
}

// t  1连续抽奖次数,2、天数
func GetWishActPoolRuleListByPoolId(poolId int64, t int32) []*share_message.WishActPoolRule {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACT_POOL_RULE)
	defer closeFun()
	list := make([]*share_message.WishActPoolRule, 0)
	err := col.Find(bson.M{"WishActPoolId": poolId, "Type": t}).Sort("Key").All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishActPoolRuleList err: %s", err.Error())
		return list
	}
	return list
}
func GetWishActPoolRuleListByPoolId1(poolId int64) []*share_message.WishActPoolRule {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACT_POOL_RULE)
	defer closeFun()
	list := make([]*share_message.WishActPoolRule, 0)
	err := col.Find(bson.M{"WishActPoolId": poolId}).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishActPoolRuleList err: %s", err.Error())
		return list
	}
	return list
}

//根据奖池查询
func GetPoolRuleIdsByPoolId(poolId int64) []int64 {
	list := GetWishActPoolRuleListByPoolId1(poolId)
	var ruleIds []int64
	for _, v := range list {
		ruleIds = append(ruleIds, v.GetId())
	}
	return ruleIds
}

//根据奖池和类型查询累计次数
func SumContinuedCount(pId, actPoolId int64) (int32, int32) {
	var drawTotal, dayTotal int32
	wishPlayerActivity := GeWishPlayerActivityByPId(pId)
	ruleIds := GetPoolRuleIdsByPoolId(actPoolId)
	//logs.Info("rules--", ruleIds)
	activityDatas := wishPlayerActivity.GetData()
	for _, activityData := range activityDatas {
		//不属于当前奖池的过滤掉
		if !util.Int64InSlice(activityData.GetPoolRuleId(), ruleIds) {
			continue
		}

		//不是天数和次数的也过滤掉
		typeNum := activityData.GetType()
		if typeNum == 1 || typeNum == 2 {
			if typeNum == 1 {
				//logs.Info("activityData-----", activityData)
				if drawTotal < int32(activityData.GetValue()) {
					drawTotal = int32(activityData.GetValue()) //累计次数
				}
			} else {
				dayTotal = int32(activityData.GetValue()) //累计天数
			}
		}
	}
	return drawTotal, dayTotal
}

func GetWishActivityPrizeLogByRuleId(uid, ruleId int64) *share_message.WishActivityPrizeLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACTIVITY_PRIZE_LOG)
	defer closeFun()
	var data *share_message.WishActivityPrizeLog
	err := col.Find(bson.M{"WishActPoolRuleId": ruleId, "PlayerId": uid}).Sort("-CreateTime").One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishActPoolRuleList err: %s", err.Error())
		return data
	}
	return data
}
func GetWishActivityPrizeLogByRuleIdAndLun(uid, ruleId int64, dayLun int32) *share_message.WishActivityPrizeLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACTIVITY_PRIZE_LOG)
	defer closeFun()
	var data *share_message.WishActivityPrizeLog
	err := col.Find(bson.M{"WishActPoolRuleId": ruleId, "PlayerId": uid, "DayLun": dayLun}).Sort("-CreateTime").One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishActPoolRuleList err: %s", err.Error())
		return data
	}
	return data
}

// 获取排名周榜月榜   dataType 1-周榜;2-月榜
func GetWishPlayerActivityByPage(dataType int64, page, pageSize int) ([]*share_message.WishPlayerActivity, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER_ACTIVITY)
	defer closeFun()
	activityList := make([]*share_message.WishPlayerActivity, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)
	q := bson.M{}
	var sort string
	switch dataType { // 1-周榜;2-月榜
	case WISH_ACT_H5_WEEK_NOW:
		q["WeekPrize"] = bson.M{"$gt": 0}
		sort = "-WeekPrize"
	case WISH_ACT_H5_MONTH_NOW:
		q["MonthPrize"] = bson.M{"$gt": 0}
		sort = "-MonthPrize"
	}

	query := col.Find(q)
	count, err := query.Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishPlayerActivityByPage err: %s", err.Error())
	}
	err = query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&activityList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishPlayerActivityByPage err: %s", err.Error())
	}
	return activityList, count
}

func GetWishPlayerActivityByNum(num int, dataType int64) []*share_message.WishPlayerActivity {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER_ACTIVITY)
	defer closeFun()
	activityList := make([]*share_message.WishPlayerActivity, 0)

	q := bson.M{}
	var sort string
	switch dataType { // 3-周榜;4-月榜
	case WISH_ACT_H5_WEEK_NOW:
		q["WeekPrize"] = bson.M{"$gt": 0}
		sort = "-WeekPrize"
	case WISH_ACT_H5_MONTH_NOW:
		q["MonthPrize"] = bson.M{"$gt": 0}
		sort = "-MonthPrize"
	}

	err := col.Find(q).Sort(sort).Limit(num).All(&activityList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishPlayerActivityByPage err: %s", err.Error())
	}
	return activityList
}

// 3、周排名 4、月排名,key 排名
func GetWishActPoolRuleListByTypeKey(dataType int64, key, num int) []*share_message.WishActPoolRule {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACT_POOL_RULE)
	defer closeFun()
	list := make([]*share_message.WishActPoolRule, 0)
	var t int32
	switch dataType {
	case WISH_ACT_H5_WEEK_NOW:
		t = WISH_ACTIVITY_DATA_TYPE_3
	case WISH_ACT_H5_MONTH_NOW:
		t = WISH_ACTIVITY_DATA_TYPE_4
	}
	q := bson.M{"Type": t, "Key": bson.M{"$gt": key}}
	err := col.Find(q).Sort("Key").Limit(num).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishActPoolRuleList err: %s", err.Error())
		return list
	}
	return list
}
func GetWishActPoolRuleListByType(dataType int64) []*share_message.WishActPoolRule {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACT_POOL_RULE)
	defer closeFun()
	list := make([]*share_message.WishActPoolRule, 0)

	q := bson.M{"Type": dataType}
	err := col.Find(q).Sort("Key").All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishActPoolRuleList err: %s", err.Error())
		return list
	}
	return list
}

// 根据id查询奖项
func GetWishActivityPrizeLogById(id int64) *share_message.WishActivityPrizeLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACTIVITY_PRIZE_LOG)
	defer closeFun()
	var data *share_message.WishActivityPrizeLog
	err := col.Find(bson.M{"_id": id}).One(&data)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishActivityPrizeLogById err: %s", err.Error())
		return data
	}
	return data
}

// 根据周榜,月榜获取前面几行
func GetWishPlayerActivityTop(dataType, num int) []*share_message.WishPlayerActivity {
	col, closeFunItem := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER_ACTIVITY)
	defer closeFunItem()
	list := make([]*share_message.WishPlayerActivity, 0)
	var sort string
	var grep string
	switch dataType {
	case WISH_ACTIVITY_DATA_TYPE_3:
		sort = "-WeekPrize"
		grep = "WeekPrize"
	case WISH_ACTIVITY_DATA_TYPE_4:
		sort = "-MonthPrize"
		grep = "MonthPrize"
	default:
		return list
	}
	err := col.Find(bson.M{grep: bson.M{"$gt": 0}}).Sort(sort, "-UpdateDiamondTime").Limit(num).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return list
}

func GeWishActivityPrizeLogByIds(ids []int64) []*share_message.WishActivityPrizeLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_ACTIVITY_PRIZE_LOG)
	defer closeFun()
	list := make([]*share_message.WishActivityPrizeLog, 0)
	q := bson.M{"_id": bson.M{"$in": ids}}

	err := col.Find(q).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GeWishActivityPrizeLogByIds err: %s", err.Error())
	}
	return list
}

// 删除周榜数据
func RemoveAllWeekTop() error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_WEEK_TOP)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{})
	easygo.PanicError(err)
	return err
}

// 删除月数据
func RemoveAllMonthTop() error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_MONTH_TOP)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{})
	easygo.PanicError(err)
	return err
}

func GetWishWeekTop(num int) []*share_message.WishWeekTop {
	col, closeFunItem := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_WEEK_TOP)
	defer closeFunItem()
	list := make([]*share_message.WishWeekTop, 0)

	err := col.Find(bson.M{}).Sort("-WeekPrize").Limit(num).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return list
}
func GetWishMonthTop(num int) []*share_message.WishMonthTop {
	col, closeFunItem := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_MONTH_TOP)
	defer closeFunItem()
	list := make([]*share_message.WishMonthTop, 0)

	err := col.Find(bson.M{}).Sort("-MonthPrize").Limit(num).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return list
}

// 获取排名 月榜
func GetWishMonthTopByPage(page, pageSize int) ([]*share_message.WishMonthTop, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_MONTH_TOP)
	defer closeFun()
	activityList := make([]*share_message.WishMonthTop, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)
	q := bson.M{}
	q["MonthPrize"] = bson.M{"$gt": 0}
	sort := "-MonthPrize"

	query := col.Find(q)
	count, err := query.Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishMonthTopByPage err: %s", err.Error())
	}
	err = query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&activityList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishMonthTopByPage err: %s", err.Error())
	}
	return activityList, count
}

// 获取排名 周榜
func GetWishWeekTopByPage(page, pageSize int) ([]*share_message.WishWeekTop, int) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_WEEK_TOP)
	defer closeFun()
	activityList := make([]*share_message.WishWeekTop, 0)
	curPage := easygo.If(page > 1, page-1, 0).(int)
	q := bson.M{}
	q["WeekPrize"] = bson.M{"$gt": 0}
	sort := "-WeekPrize"

	query := col.Find(q)
	count, err := query.Count()
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishWeekTopByPage err: %s", err.Error())
	}
	err = query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&activityList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("GetWishWeekTopByPage err: %s", err.Error())
	}
	return activityList, count
}

// 批量查询
func GetWishPlayerActivityByGTPid(pid int64) []*share_message.WishPlayerActivity {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER_ACTIVITY)
	defer closeFun()
	datas := make([]*share_message.WishPlayerActivity, 0)
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

// 许愿池 白名单列表
func GetWishWhiteList() []*share_message.WishWhite {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_WHITE)
	defer closeFun()
	datas := make([]*share_message.WishWhite, 0)
	queryBson := bson.M{}
	err := col.Find(queryBson).All(&datas)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return datas
}
func InsertWishWhite(data *share_message.WishWhite) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_WHITE)
	defer closeFun()
	return col.Insert(data)
}

//
func InsertWishDayActivityLog(data *share_message.WishDayActivityLog) error {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_DAY_ACTIVITY_LOG)
	defer closeFun()
	dLog := GetWishDayActivityLog(data)
	if dLog != nil && dLog.GetPlayerId() > 0 {
		logs.Error("不计入时长,查出的内容是: %+v", dLog)
		return nil
	}
	return col.Insert(data)
}
func GetWishDayActivityLog(data *share_message.WishDayActivityLog) *share_message.WishDayActivityLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_DAY_ACTIVITY_LOG)
	defer closeFun()
	var d *share_message.WishDayActivityLog
	err := col.Find(bson.M{"PlayerId": data.GetPlayerId(), "CreateTime": data.GetCreateTime(), "WishActPoolId": data.GetWishActPoolId()}).One(&d)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return d
}
func GetWishDayActivityLogList(pid, actPoolId int64) []*share_message.WishDayActivityLog {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_DAY_ACTIVITY_LOG)
	defer closeFun()
	list := make([]*share_message.WishDayActivityLog, 0)
	err := col.Find(bson.M{"PlayerId": pid, "WishActPoolId": actPoolId}).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return list
}

// 检测许愿池活动是否开启
func CheckWishActIsOpen() (bool, bool, bool) {
	var isOpenCount, isOpenDay, isOpenWeekMonth bool
	actTypeList := GetActivityByTypes([]int32{WISH_ACT_COUNT, WISH_ACT_DAY, WISH_ACT_WEEK_MONTH})
	if len(actTypeList) == 0 {
		logs.Error("没有活动配置")
		return isOpenCount, isOpenDay, isOpenWeekMonth
	}
	for _, v := range actTypeList {
		if v.GetStatus() > 0 {
			logs.Error("actType:%d 活动状态是关闭状态", v.GetTypes())
			continue
		}
		t := time.Now().Unix()
		if t < v.GetStartTime() || t > v.GetEndTime() {
			logs.Error("actType:%d 活动时间为开始,配置的开始时间为: %d,结束时间为: %d,当前时间为: %d", v.GetStartTime(), v.GetEndTime(), t)
			continue
		}
		switch v.GetTypes() {
		case WISH_ACT_COUNT:
			isOpenCount = true
		case WISH_ACT_DAY:
			isOpenDay = true
		case WISH_ACT_WEEK_MONTH:
			isOpenWeekMonth = true
		}
	}
	return isOpenCount, isOpenDay, isOpenWeekMonth
}
