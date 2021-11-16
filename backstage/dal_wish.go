package backstage

import (
	"errors"
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"strings"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/astaxie/beego/logs"

	"github.com/akqp2019/mgo/bson"
)

// 设置页签
func SetMgoPage(page, cur int32) (int, int) {
	pageSize := int(page)
	curPage := int(cur)
	curPage = easygo.If(curPage > 1, curPage-1, 0).(int)

	return pageSize, curPage
}

func SingleRoomInit() {
	// 房间标签初始化
	cfg := GetWishPoolCfgOnce()
	if cfg.GetStatus() == 0 {
		OnceWishPool()
		UpdateWishPoolCfgOnce()
	}
}

// 一次性执行代码
func OnceWishPool() {
	list := GetWishPoolAll()

	if len(list) > 0 {
		for_game.SetCurrentId(for_game.TABLE_WISH_POOL_CFG, 100)
	}

	for _, v := range list {

		one := &brower_backstage.WishPool{
			Id:               easygo.NewInt64(v.GetId()),
			PoolLimit:        easygo.NewInt32(v.GetPoolLimit()),
			Name:             easygo.NewString(v.GetName()),
			ShowInitialValue: easygo.NewInt32(v.GetShowInitialValue()),
			ShowRecycle:      easygo.NewInt32(v.GetShowRecycle()),
			ShowCommission:   easygo.NewInt32(v.GetShowCommission()),
			ShowStartAward:   easygo.NewInt32(v.GetShowStartAward()),
			ShowCloseAward:   easygo.NewInt32(v.GetShowCloseAward()),
			IsDefault:        easygo.NewBool(v.GetIsDefault()),
		}

		one.SmallLoss = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(v.SmallLoss.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(v.SmallLoss.GetShowMinValue()),
		}

		one.SmallWin = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(v.SmallWin.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(v.SmallWin.GetShowMinValue()),
		}

		one.BigLoss = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(v.BigLoss.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(v.BigLoss.GetShowMinValue()),
		}

		one.BigWin = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(v.BigWin.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(v.BigWin.GetShowMinValue()),
		}
		one.Common = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(v.Common.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(v.Common.GetShowMinValue()),
		}

		UpdateWishPool(one)

		new := &share_message.WishPool{
			Id:           easygo.NewInt64(v.GetId()),
			PoolConfigId: easygo.NewInt64(v.GetId()),
		}
		UpdateWishPoolDB(new)
	}

}

// 获取水池初始化状态
func GetWishPoolCfgOnce() *share_message.SingleInitCfg {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	list := &share_message.SingleInitCfg{}
	queryBson := bson.M{"_id": for_game.TABLE_WISH_POOL_CFG_ONCE}
	err := col.Find(queryBson).One(&list)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}

	if err == mgo.ErrNotFound {
		return list
	}

	return list
}

// 更新水池初始化状态
func UpdateWishPoolCfgOnce() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	data := &share_message.SingleInitCfg{
		Status: easygo.NewInt32(1),
	}

	_, err := col.Upsert(bson.M{"_id": for_game.TABLE_WISH_POOL_CFG_ONCE}, bson.M{"$set": data})
	easygo.PanicError(err)
}

//获取用户信息
func GetWishPlayerInfoByPid(playerId int64) *share_message.WishPlayer {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_PLAYER)
	defer closeFun()
	var player *share_message.WishPlayer
	err := col.Find(bson.M{"PlayerId": playerId}).One(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return player
}

//获取用户信息
func GetWishPlayerInfo(playerId int64) *share_message.WishPlayer {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_PLAYER)
	defer closeFun()
	var player *share_message.WishPlayer
	err := col.Find(bson.M{"_id": playerId}).One(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return player
}

//根据id集合获取指定的玩家集合信息
func GetWishPlayersByIds(ids []int64) []*share_message.WishPlayer {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_PLAYER)
	defer closeFun()
	var player []*share_message.WishPlayer
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return player
}

//根据玩家id集合获取wishPlayer列表
func GetWishPlayerListByPlayers(players []int64) []*share_message.WishPlayer {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_PLAYER)
	defer closeFun()
	var wishPlayerList []*share_message.WishPlayer
	err := col.Find(bson.M{"PlayerId": bson.M{"$in": players}}).All(&wishPlayerList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return wishPlayerList
}

//根据PlayerId获取白名单
func GetWishWhiteListByPlayers(players []int64) []*share_message.WishWhite {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_WHITE)
	defer closeFun()
	var wishWhiteList []*share_message.WishWhite
	err := col.Find(bson.M{"PlayerId": players}).Select(bson.M{"Account": 1, "Note": 1}).All(&wishWhiteList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return wishWhiteList
}

func GetWishWhiteListByAccounts(accounts []string) []*share_message.WishWhite {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_WHITE)
	defer closeFun()
	var wishWhiteList []*share_message.WishWhite
	err := col.Find(bson.M{"PlayerId": accounts}).Select(bson.M{"Account": 1, "Note": 1}).All(&wishWhiteList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return wishWhiteList
}

//根据柠檬号获取白名单
func GetWishWhiteInfoByAccount(account string) (*share_message.WishWhite, error) {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_WHITE)
	defer closeFun()
	var wishWhiteInfo *share_message.WishWhite
	err := col.Find(bson.M{"Account": account}).One(&wishWhiteInfo)
	return wishWhiteInfo, err
}

//获取所有白名单
func GetAllWishWhite() []*share_message.WishWhite {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_WHITE)
	defer closeFun()
	var wishWhiteList []*share_message.WishWhite
	err := col.Find(bson.M{}).All(&wishWhiteList)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return wishWhiteList
}

//获取白名单所用用户的id集合
func GetAllIdsInWishWhite() []int64 {
	list := GetAllWishWhite()
	var players []int64
	for _, v := range list {
		players = append(players, v.GetId())
	}
	return players
}

func UpdateWishOrder(id string, status int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER)
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

// 获取许愿池商品
func AddWishItem(data *share_message.WishItem) int64 {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM)
	defer closeFun()

	data.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_ITEM))

	err := col.Insert(data)
	easygo.PanicError(err)

	return data.GetId()
}

// 更新许愿池商品
func UpdateWishItem(data *share_message.WishItem) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM)
	defer closeFun()

	if data.GetId() == 0 {
		data.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_ITEM))
	}

	_, err := col.Upsert(bson.M{"_id": data.GetId()}, bson.M{"$set": data})
	easygo.PanicError(err)
}

// 根据用户id获取收藏盲盒数
func GetWishPlayerInfos(id []int64) []*brower_backstage.DrawRecord {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_PLAYER)
	defer closeFun()

	match := bson.M{"_id": bson.M{"$in": id}}
	group := bson.M{
		"_id":          "$_id",
		"UserId":       bson.M{"$last": "$_id"},
		"UserAccount":  bson.M{"$last": "$Account"},
		"UserNickname": bson.M{"$last": "$NickName"},
		"Phone":        bson.M{"$last": "$Phone"},
	}

	project := bson.M{
		"UserId":      1,
		"AddBoxCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var counts []*brower_backstage.DrawRecord
	err := col.Pipe(pipeCond).All(&counts)
	easygo.PanicError(err)

	return counts
}

//获取用户信息
func GetPlayerInfoByPid(playerId int64) *share_message.PlayerBase {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()
	player := &share_message.PlayerBase{}
	err := col.Find(bson.M{"_id": playerId}).One(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return player
}

// 根据条件获取相应用户id
func GetUserIdByBson(queryBson bson.M) []int64 {
	selBson := bson.M{"_id": 1}
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	var list []*share_message.PlayerBase
	col.Find(queryBson).Select(selBson).All(&list)

	var ids []int64

	for _, v := range list {
		ids = append(ids, v.GetPlayerId())
	}

	return ids
}

//获取水池配置信息
func GetWishPoolCfg(id int64) *share_message.WishPoolCfg {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL_CFG)
	defer closeFun()
	player := &share_message.WishPoolCfg{}
	err := col.Find(bson.M{"_id": id}).One(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return player
}

//获取盲盒水池信息
func GetWishPool(id int64) *share_message.WishPool {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL)
	defer closeFun()
	player := &share_message.WishPool{}
	err := col.Find(bson.M{"_id": id}).One(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return player
}

// 根据条件获取相应用水池id
func GetWishPoolIdByBson(queryBson bson.M) []int64 {
	selBson := bson.M{"_id": 1}
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL)
	defer closeFun()

	var list []*share_message.WishPool
	col.Find(queryBson).Select(selBson).All(&list)

	var ids []int64

	for _, v := range list {
		ids = append(ids, v.GetId())
	}

	return ids
}

// 根据条件获取相应用水池id
func GetWishPoolAll() []*share_message.WishPool {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL)
	defer closeFun()

	var list []*share_message.WishPool
	col.Find(bson.M{}).All(&list)

	return list
}

// 根据名称获取盲盒信息
func GetWishBoxByName(name string) *share_message.WishBox {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX)
	defer closeFun()

	var one *share_message.WishBox
	queryBson := bson.M{"Name": name}
	err := col.Find(queryBson).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one
}

// 根据Id获取盲盒信息
func GetWishBoxById(id int64) *share_message.WishBox {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX)
	defer closeFun()

	var one *share_message.WishBox
	queryBson := bson.M{"_id": id}
	err := col.Find(queryBson).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one
}

// 根据条件获取盲盒信息
func GetWishBoxByBson(qBson bson.M) *share_message.WishBox {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX)
	defer closeFun()

	var one *share_message.WishBox
	queryBson := qBson
	err := col.Find(queryBson).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one
}

// 根据Id获取盲盒信息
func GetWishBoxList(ids []int64) []*share_message.WishBox {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX)
	defer closeFun()

	var list []*share_message.WishBox
	queryBson := bson.M{}
	if len(ids) > 0 {
		queryBson = bson.M{"_id": bson.M{"$in": ids}}
	}

	err := col.Find(queryBson).All(&list)
	easygo.PanicError(err)

	return list
}

//获取可用的盲盒列表
func GetActiveWishBoxList() []*share_message.WishBox {
	boxIds := for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL, "BoxIds", nil)
	var ids []int64
	ids = for_game.InterfersToInt64(boxIds)
	if len(ids) > 0 {
		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX)
		defer closeFun()
		var list []*share_message.WishBox
		queryBson := bson.M{}
		queryBson = bson.M{"_id": bson.M{"$nin": ids}}
		err := col.Find(queryBson).All(&list)
		easygo.PanicError(err)
		return list
	}
	return nil
}

//根据守护者id查询盲盒列表
func GetWishBoxListByGuardianIds(GuardianId int64) []*share_message.WishBox {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX)
	defer closeFun()

	var list []*share_message.WishBox
	queryBson := bson.M{"GuardianId": GuardianId}
	err := col.Find(queryBson).All(&list)
	easygo.PanicError(err)
	return list
}

// 根据Id获取盲盒信息
func GetWishItemById(id int64) *share_message.WishItem {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM)
	defer closeFun()

	var one *share_message.WishItem
	queryBson := bson.M{"_id": id}
	err := col.Find(queryBson).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one
}

// 根据Ids获取许愿池商品
func GetWishItemByIds(id []int64) []*share_message.WishItem {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM)
	defer closeFun()

	var list []*share_message.WishItem
	queryBson := bson.M{"_id": bson.M{"$in": id}}
	err := col.Find(queryBson).All(&list)
	easygo.PanicError(err)

	return list
}

// 根据Id获取盲盒信息
func GetWishBoxItemById(id int64) *share_message.WishBoxItem {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX_ITEM)
	defer closeFun()

	var one *share_message.WishBoxItem
	queryBson := bson.M{"_id": id}
	err := col.Find(queryBson).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one
}

// 获取全部商品
func GetAllWishItem() []*share_message.WishItem {
	var lst []*share_message.WishItem
	col, closeFunc := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM)
	defer closeFunc()
	err := col.Find(bson.M{}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return lst
}

// 获取全部盲盒
func GetAllWishBox() []*share_message.WishBox {
	var lst []*share_message.WishBox
	col, closeFunc := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX)
	defer closeFunc()
	err := col.Find(bson.M{}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return lst
}

// 获取全部盲盒
func GetAllWishBoxOnline() []*share_message.WishBox {
	var lst []*share_message.WishBox
	col, closeFunc := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX)
	defer closeFunc()
	err := col.Find(bson.M{"Status": 1}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return lst
}

// 获取全部盲盒商品
func GetAllWishBoxItem() []*share_message.WishBoxItem {
	var lst []*share_message.WishBoxItem
	col, closeFunc := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX_ITEM)
	defer closeFunc()
	err := col.Find(bson.M{}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return lst
}

// 获取全部盲盒商品
func GetAllWishBoxItemByBoxId(id int64) []*share_message.WishBoxItem {
	var lst []*share_message.WishBoxItem
	col, closeFunc := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX_ITEM)
	defer closeFunc()
	err := col.Find(bson.M{"WishBoxId": id}).All(&lst)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return lst
}

func UpdateBoxGuardianInfo() {
	//col, closeFunc := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OCCUPIED)
	//defer closeFunc()
}

// 获取商品类型
func GetWishTypeList() []*share_message.WishStyle {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_STYLE)
	defer closeFun()

	queryBson := bson.M{}
	var list []*share_message.WishStyle
	err := col.Find(queryBson).All(&list)
	easygo.PanicError(err)

	return list
}

// 获取盲盒列表
func QueryWishBoxList(req *brower_backstage.WishBoxListRequest) ([]*brower_backstage.WishBox, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"_id"}

	if req.GetSort() != "" {
		sort = []string{req.GetSort()}
	}

	queryBson := bson.M{}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	if req.GetKeyword() != "" {
		switch req.GetType() {
		case 1: // 盲盒id
			id := easygo.StringToInt64noErr(req.GetKeyword())
			queryBson["_id"] = id
		case 2: // 盲盒名称
			queryBson["Name"] = req.GetKeyword()
		case 3: // 守护者柠檬号
			user := for_game.GetPlayerByAccount(req.GetKeyword())
			queryBson["GuardianId"] = user.GetPlayerId()
		default:

		}
	}

	// 是否挑战赛
	if req.IsChallenge != nil {
		if req.GetIsChallenge() {
			queryBson["Match"] = 1
		} else {
			queryBson["Match"] = 0
		}
	}

	// 盲盒属性
	if req.Attribute != nil && len(req.GetAttribute()) > 0 {
		queryBson["Menu"] = bson.M{"$in": req.GetAttribute()}
	}

	// 当前水池的状态
	if req.LocalStatus != nil && len(req.GetLocalStatus()) > 0 {
		qBson := bson.M{}
		qBson["LocalStatus"] = bson.M{"$in": req.GetLocalStatus()}
		wishIds := GetWishPoolIdByBson(qBson)
		queryBson["WishPoolId"] = bson.M{"$in": wishIds}
	}

	// 状态
	if req.Status != nil && req.GetStatus() != 1000 {
		queryBson["Status"] = req.GetStatus()
	}

	// 是否有守护者
	if req.IsHasUser != nil {
		if req.GetIsHasUser() {
			queryBson["GuardianId"] = bson.M{"$gt": 0}
		} else {
			queryBson["$or"] = []bson.M{{"GuardianId": nil}, {"GuardianId": 0}}
		}

	}

	if req.HaveIsWin != nil {
		queryBson["HaveIsWin"] = req.GetHaveIsWin()
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishBox
	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	var retList []*brower_backstage.WishBox

	for _, v := range list {
		one := &brower_backstage.WishBox{
			Id:          easygo.NewInt64(v.GetId()),
			Name:        easygo.NewString(v.GetName()),
			Icon:        easygo.NewString(v.GetIcon()),
			GoodsAmount: easygo.NewInt32(v.GetTotalNum()),
			Attribute:   v.GetMenu(),
			UserId:      easygo.NewInt64(v.GetGuardianId()),
			Price:       easygo.NewInt64(v.GetPrice()),
			Status:      easygo.NewInt32(v.GetStatus()),
			CreateTime:  easygo.NewInt64(v.GetCreateTime()),
			UpdateTime:  easygo.NewInt64(v.GetUpdateTime()),
			SortWeight:  easygo.NewInt64(v.GetSortWeight()),
			WishPoolId:  easygo.NewInt64(v.GetWishPoolId()),
			IsRecommend: easygo.NewBool(v.GetIsRecommend()),
			UploadTime:  easygo.NewInt64(v.GetPutOnTime()),
			HaveIsWin:   easygo.NewBool(v.GetHaveIsWin()),
		}

		// 是否挑战赛
		if v.GetMatch() == 1 {
			one.IsChallenge = easygo.NewBool(true)
		} else {
			one.IsChallenge = easygo.NewBool(false)
		}

		user := for_game.GetWishPlayerByPid(one.GetUserId())
		one.UserId = easygo.NewInt64(user.GetPlayerId())
		user2 := GetPlayerInfoByPid(user.GetPlayerId())
		one.UserAccount = easygo.NewString(user2.GetAccount())

		retList = append(retList, one)
	}

	return retList, int32(count)
}

//  获取盲盒包含的商品列表
func GetWishBoxGoodsItemList(req *brower_backstage.ListRequest) ([]*brower_backstage.WishBoxGoodsWin, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	id := req.GetId()
	queryBson := bson.M{"WishBoxId": id, "ItemStatus": 1}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX_ITEM)
	defer closeFun()

	count, _ := col.Find(queryBson).Count()

	lookup := bson.M{"from": for_game.TABLE_WISH_ITEM, "localField": "WishItemId", "foreignField": "_id", "as": "item"}
	match := bson.M{"WishBoxId": id}

	project := bson.M{
		"Id":                    "$_id",
		"GoodsId":               "$WishItemId",
		"Name":                  "$item.Name",
		"Price":                 "$Price",
		"IsInfallible":          "$IsWin",
		"ReplenishAmount":       "$PerNum",
		"ReplenishIntervalTime": "$PerTime",
		"Weight":                "$PerRate",
		"GoodsType":             "$Style",

		"BigLoss":   "$BigLoss",
		"SmallLoss": "$SmallLoss",
		"Common":    "$Common",
		"BigWin":    "$BigWin",
		"SmallWin":  "$SmallWin",

		"CommonAddWeight":   "$CommonAddWeight",
		"BigWinAddWeight":   "$BigWinAddWeight",
		"SmallWinAddWeight": "$SmallWinAddWeight",

		"PerRate":  "$PerRate",
		"RewardLv": "$RewardLv",
		"Diamond":  "$Diamond",

		"BoxItemId": "$_id",
	}
	sort := bson.M{"Style": 1, "CreateTime": -1}
	pipeCond := []bson.M{
		{"$match": match},
		{"$lookup": lookup},
		{"$unwind": "$item"},
		{"$project": project},
		{"$sort": sort},
	}

	if pageSize > 0 {
		pipeCond = append(pipeCond, bson.M{"$skip": curPage * pageSize})
		pipeCond = append(pipeCond, bson.M{"$limit": pageSize})
	}

	var list []*brower_backstage.WishBoxGoodsWin
	err := col.Pipe(pipeCond).All(&list)
	easygo.PanicError(err)

	return list, int32(count)
}

//根据盲盒Id查询盲盒下的商品
func GetGoodsItemByBoxId(boxId int64) []*brower_backstage.WishBoxGoodsWin {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX_ITEM)
	defer closeFun()

	lookup := bson.M{"from": for_game.TABLE_WISH_ITEM, "localField": "WishItemId", "foreignField": "_id", "as": "item"}
	match := bson.M{"WishBoxId": boxId}

	project := bson.M{
		"Id":      "$_id",
		"GoodsId": "$WishItemId",
		"Name":    "$item.Name",
		"Price":   "$Price",
		/*"IsInfallible":          "$IsWin",
		"ReplenishAmount":       "$PerNum",
		"ReplenishIntervalTime": "$PerTime",
		"Weight":                "$PerRate",
		"GoodsType":             "$Style",

		"BigLoss":   "$BigLoss",
		"SmallLoss": "$SmallLoss",
		"Common":    "$Common",
		"BigWin":    "$BigWin",
		"SmallWin":  "$SmallWin",

		"CommonAddWeight":   "$CommonAddWeight",
		"BigWinAddWeight":   "$BigWinAddWeight",
		"SmallWinAddWeight": "$SmallWinAddWeight",

		"PerRate":  "$PerRate",
		"RewardLv": "$RewardLv",
		"Diamond":  "$Diamond",*/

		"BoxItemId": "$_id",
	}
	sort := bson.M{"Style": 1, "CreateTime": -1}
	pipeCond := []bson.M{
		{"$match": match},
		{"$lookup": lookup},
		{"$unwind": "$item"},
		{"$project": project},
		{"$sort": sort},
	}

	var list []*brower_backstage.WishBoxGoodsWin
	err := col.Pipe(pipeCond).All(&list)
	easygo.PanicError(err)

	return list
}

// 获取所属盲盒商品信息
func GetWishBoxGoodsItemById(wishId, itemId int64) {
	queryBson := bson.M{"WishBoxId": wishId, "WishItemId": itemId}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX_ITEM)
	defer closeFun()

	var list []*brower_backstage.WishBoxGoodsItem
	err := col.Find(queryBson).All(&list)
	easygo.PanicError(err)
}

// 根据商品id获取所有相应盲盒商品id
func GetWishBoxItemIdsByItemId(itemId int64) []int64 {
	queryBson := bson.M{"WishItemId": itemId}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX_ITEM)
	defer closeFun()

	var (
		ids  []int64
		list []*share_message.WishBoxItem
	)
	err := col.Find(queryBson).Select(bson.M{"_id": 1}).All(&list)
	easygo.PanicError(err)

	for i := range list {
		ids = append(ids, list[i].GetId())
	}

	return ids
}

// 新增/更新盲盒
func UpdateWishBox(data *brower_backstage.WishBox) error {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX)
	defer closeFun()

	upData := &share_message.WishBox{
		Id:            easygo.NewInt64(data.GetId()),
		Name:          easygo.NewString(data.GetName()),
		Icon:          easygo.NewString(data.GetIcon()),
		Price:         easygo.NewInt64(data.GetPrice()),
		Status:        easygo.NewInt32(data.GetStatus()),
		SortWeight:    easygo.NewInt64(data.GetSortWeight()),
		GuardianId:    easygo.NewInt64(data.GetUserId()),
		Menu:          data.GetAttribute(),
		Items:         data.GetItems(),
		Brands:        data.GetBrands(),
		Styles:        data.GetStyles(),
		Types:         data.GetTypes(),
		ProductStatus: easygo.NewInt32(data.GetProductStatus()),
		WishItems:     data.GetWishItems(),
		UpdateTime:    easygo.NewInt64(data.GetUpdateTime()),
		PutOnTime:     easygo.NewInt64(data.GetUploadTime()),
		CreateTime:    easygo.NewInt64(data.GetCreateTime()),
		TotalNum:      easygo.NewInt32(data.GetGoodsAmount()),
		RareNum:       easygo.NewInt32(data.GetRareNum()),
		IsRecommend:   easygo.NewBool(data.GetIsRecommend()),
		HaveIsWin:     easygo.NewBool(data.GetHaveIsWin()),
	}

	if data.WishPoolId != nil && data.GetWishPoolId() > 0 {
		upData.WishPoolId = easygo.NewInt64(data.GetWishPoolId())
	}

	if data.GuardianOverTime != nil && data.GetGuardianOverTime() > 0 {
		upData.GuardianOverTime = easygo.NewInt64(data.GetGuardianOverTime())
	}

	if data.IsGuardian != nil {
		upData.IsGuardian = easygo.NewBool(data.GetIsGuardian())
	}

	if data.GetIsChallenge() {
		upData.Match = easygo.NewInt32(1)
	} else {
		upData.Match = easygo.NewInt32(0)
	}

	// 新增

	_, err1 := col.Upsert(bson.M{"_id": upData.GetId()}, bson.M{"$set": upData})
	easygo.PanicError(err1)
	return err1
}

func UpdateWishPoolDB(upData *share_message.WishPool) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL)
	defer closeFun()

	_, err1 := col.Upsert(bson.M{"_id": upData.GetId()}, bson.M{"$set": upData})
	easygo.PanicError(err1)
}

// 删除水池日志
func DeleteWishPoolLog(bid int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL_LOG)
	defer closeFun()

	_, err1 := col.RemoveAll(bson.M{"PoolId": bid})
	easygo.PanicError(err1)
}

func DeleteToolWishPoolLog(bid int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TOOL_WISH_POOL_LOG)
	defer closeFun()

	_, err1 := col.RemoveAll(bson.M{"PoolId": bid})
	easygo.PanicError(err1)
}

// 删除水池抽水日志
func DeleteWishPoolPumpLog(bid int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL_PUMP_LOG)
	defer closeFun()

	_, err1 := col.RemoveAll(bson.M{"PoolId": bid})
	easygo.PanicError(err1)
}

func DeleteToolWishPoolPumpLog(bid int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TOOL_WISH_POOL_PUMP_LOG)
	defer closeFun()

	_, err1 := col.RemoveAll(bson.M{"PoolId": bid})
	easygo.PanicError(err1)
}

func DeleteWishBoxItem(bid int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX_ITEM)
	defer closeFun()

	_, err1 := col.RemoveAll(bson.M{"WishBoxId": bid})
	easygo.PanicError(err1)
}

// 盲盒包含的商品
func UpdateWishBoxGoodsItem(data *brower_backstage.WishBoxGoodsWin) error {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX_ITEM)
	defer closeFun()

	goods := GetGoodsById(data.GetGoodsId())
	if goods == nil {
		return errors.New("商品信息错误")
	}

	upData := &share_message.WishBoxItem{
		Id:                 easygo.NewInt64(data.GetId()),
		WishItemId:         easygo.NewInt64(data.GetGoodsId()),
		WishBoxId:          easygo.NewInt64(data.GetWishBoxId()),
		Style:              easygo.NewInt32(data.GetGoodsType()),
		IsWin:              easygo.NewBool(data.GetIsInfallible()),
		PerNum:             easygo.NewInt32(data.GetReplenishAmount()),
		PerTime:            easygo.NewInt32(data.GetReplenishIntervalTime()),
		PerRate:            easygo.NewInt32(data.GetPerRate()),
		Price:              easygo.NewInt64(data.GetPrice()),
		BigLoss:            easygo.NewInt32(data.GetBigLoss()),
		SmallLoss:          easygo.NewInt32(data.GetSmallLoss()),
		Common:             easygo.NewInt32(data.GetCommon()),
		BigWin:             easygo.NewInt32(data.GetBigWin()),
		SmallWin:           easygo.NewInt32(data.GetSmallWin()),
		CommonAddWeight:    easygo.NewInt32(data.GetCommonAddWeight()),
		BigWinAddWeight:    easygo.NewInt32(data.GetBigWinAddWeight()),
		SmallWinAddWeight:  easygo.NewInt32(data.GetSmallWinAddWeight()),
		RewardLv:           easygo.NewInt32(data.GetRewardLv()),
		LocalNum:           easygo.NewInt32(data.GetReplenishAmount()),
		Diamond:            easygo.NewInt64(data.GetDiamond()),
		PredictArrivalTime: easygo.NewInt64(data.GetArrivalTime()),
		CreateTime:         easygo.NewInt64(data.GetCreateTime()),
	}

	// 判断是否是预售商品
	if goods.GetIsPreSale() {
		upData.Status = easygo.NewInt32(2)
	} else {
		upData.Status = easygo.NewInt32(1)
	}

	// 新增
	if upData.GetId() == 0 {
		upData.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_BOX_ITEM))
	}

	_, err1 := col.Upsert(bson.M{"_id": upData.GetId()}, bson.M{"$set": upData})
	easygo.PanicError(err1)

	return nil
}

// 获取盲盒中奖配置列表
func QueryWishBoxWinCfgList(req *brower_backstage.ListRequest) ([]*brower_backstage.WishBoxWinCfg, int32) {
	//pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())

	id := req.GetId()
	queryBson := bson.M{"WishBoxId": id}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM)
	defer closeFun()

	count, _ := col.Find(queryBson).Count()

	lookupWin := bson.M{"from": for_game.TABLE_WISH_BOX_ITEM_WIN_CFG, "localField": "WishItemId", "foreignField": "_id", "as": "win"}
	lookupItem := bson.M{"from": for_game.TABLE_WISH_BOX_ITEM, "localField": "_id", "foreignField": "_id", "as": "item"}
	match := bson.M{"WishBoxId": id, "ItemStatus": 1}

	project := bson.M{
		"GoodsId":   "WishItemId",
		"Name":      "$Name",
		"Price":     "$Price",
		"GoodsType": "$item.Style",
		"BigLoss":   "$win.BigLoss",
		"SmallLoss": "$win.SmallLoss",
		"Common":    "$win.Common",
		"BigWin":    "$win.BigWin",
		"SmallWin":  "$win.SmallWin",
	}
	sort := bson.M{"-Price": 1}
	pipeCond := []bson.M{
		{"$match": match},
		{"$lookup": lookupWin},
		{"$unwind": "win"},
		{"$lookup": lookupItem},
		{"$unwind": "item"},
		{"$project": project},
		{"$sort": sort},
		//{"$skip": curPage * pageSize},
		//{"$limit": pageSize},
	}

	var list []*brower_backstage.WishBoxWinCfg
	err := col.Pipe(pipeCond).All(&list)
	easygo.PanicError(err)

	return list, int32(count)
}

// TODO
// 获取盲盒中奖配置列表
func GetWishBoxWinCfgListByBoxId(id int64) []*brower_backstage.WishBoxWinCfg {
	var list []*brower_backstage.WishBoxWinCfg

	return list
}

// 新增/更新盲盒中奖配置列表
func UpdateWishBoxWinCfgList(data *brower_backstage.WishBoxWinCfg) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM)
	defer closeFun()

	upWishBox := &share_message.WishBoItemWinCfg{
		WishBoxId:  easygo.NewInt64(data.GetWishBoxId()),
		WishItemId: easygo.NewInt64(data.GetGoodsId()),
		Id:         easygo.NewInt64(data.GetWishBoxItemId()),
		BigLoss:    easygo.NewInt32(data.GetBigLoss()),
		SmallLoss:  easygo.NewInt32(data.GetSmallLoss()),
		Common:     easygo.NewInt32(data.GetCommon()),
		BigWin:     easygo.NewInt32(data.GetBigWin()),
		SmallWin:   easygo.NewInt32(data.GetSmallWin()),
	}

	// 新增
	if upWishBox.GetId() == 0 {
		return
	}

	err1 := col.Update(bson.M{"_id": upWishBox.GetId()}, bson.M{"$set": upWishBox})
	easygo.PanicError(err1)
}

// 获取商品列表
func QueryWishGoodsList(req *brower_backstage.WishBoxGoodsListRequest) ([]*brower_backstage.WishBoxGoods, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"_id"}

	if req.GetSort() != "" {
		sort = []string{req.GetSort()}
	}

	// 过滤掉，许愿池活动配置的实物奖励
	queryBson := bson.M{"UseType": bson.M{"$ne": 1}}
	if req.GetBeginTimestamp() != 0 {
		switch req.GetTimeType() {
		case 1:
			queryBson["UploadTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
		case 2:
			queryBson["PreHaveTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
		}

	}

	if req.GetKeyword() != "" {
		switch req.GetType() {
		case 1: // 商品id
			id := easygo.StringToInt64noErr(req.GetKeyword())
			queryBson["_id"] = id
		case 2: // 商品名称
			queryBson["Name"] = req.GetKeyword()
		default:

		}
	}

	// 状态
	if req.Status != nil && req.GetStatus() != 1000 {
		queryBson["Status"] = req.GetStatus()
	}

	// 是否预售
	if req.IsPreSale != nil {
		if req.GetIsPreSale() {
			queryBson["Status"] = 2
		} else {
			queryBson["Status"] = bson.M{"$in": []int32{0, 1}}
		}
	}

	// 品牌
	if len(req.GetWishBrandId()) > 0 {
		queryBson["Brand"] = bson.M{"$in": req.GetWishBrandId()}
	}

	// 分类
	if len(req.GetWishItemTypeId()) > 0 {
		queryBson["Type"] = bson.M{"$in": req.GetWishItemTypeId()}
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishItem
	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	var retList []*brower_backstage.WishBoxGoods

	for _, v := range list {
		one := &brower_backstage.WishBoxGoods{
			Id:             easygo.NewInt64(v.GetId()),
			Name:           easygo.NewString(v.GetName()),
			Icon:           easygo.NewString(v.GetIcon()),
			Price:          easygo.NewInt64(v.GetPrice()),
			Status:         easygo.NewInt32(v.GetStatus()),
			WishBrandId:    easygo.NewInt32(v.GetBrand()),
			WishItemTypeId: easygo.NewInt32(v.GetType()),
			IsPreSale:      easygo.NewBool(v.GetIsPreSale()),
			ArrivalTime:    easygo.NewInt64(v.GetPreHaveTime()),
			StockAmount:    easygo.NewInt32(v.GetStockAmount()),
			UploadTime:     easygo.NewInt64(v.GetUploadTime()),
			SoldOutTime:    easygo.NewInt64(v.GetSoldOutTime()),
			Describe:       easygo.NewString(v.GetDesc()),
			Diamond:        easygo.NewInt64(v.GetDiamond()),
			UpdateTime:     easygo.NewInt64(v.GetUpdateTime()),
		}

		retList = append(retList, one)
	}

	return retList, int32(count)
}

// 获取商品
func GetGoodsById(id int64) *share_message.WishItem {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM)
	defer closeFun()

	var one *share_message.WishItem
	queryBson := bson.M{"_id": id}
	err := col.Find(queryBson).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one
}

// 检查是否有盲盒包含该商品
func CheckBoxHasGoods(id int64) bool {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX)
	defer closeFun()

	var one *share_message.WishBox
	queryBson := bson.M{"wishItems": id, "Status": 1}
	err := col.Find(queryBson).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return false
	}

	if one.GetId() == 0 {
		return false
	}

	return true
}

// 检查是否有盲盒包含该商品
func GetBoxByItemId(id int64) []*share_message.WishBox {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX)
	defer closeFun()

	var list []*share_message.WishBox
	queryBson := bson.M{"wishItems": id, "Status": 1}
	err := col.Find(queryBson).All(&list)
	easygo.PanicError(err)

	return list
}

// 新增/更新商品
func UpdateWishGoods(data *brower_backstage.WishBoxGoods) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM)
	defer closeFun()

	curTime := time.Now().Unix()
	upWishBox := &share_message.WishItem{
		Id:          easygo.NewInt64(data.GetId()),
		Name:        easygo.NewString(data.GetName()),
		Icon:        easygo.NewString(data.GetIcon()),
		Price:       easygo.NewInt64(data.GetPrice()),
		Status:      easygo.NewInt32(data.GetStatus()),
		Brand:       easygo.NewInt32(data.GetWishBrandId()),
		Type:        easygo.NewInt32(data.GetWishItemTypeId()),
		StockAmount: easygo.NewInt32(data.GetStockAmount()),
		PreHaveTime: easygo.NewInt64(data.GetArrivalTime()),
		IsPreSale:   easygo.NewBool(data.GetIsPreSale()),
		Desc:        easygo.NewString(data.GetDescribe()),
		UploadTime:  easygo.NewInt64(data.GetUploadTime()),
		SoldOutTime: easygo.NewInt64(data.GetSoldOutTime()),
		Diamond:     easygo.NewInt64(data.GetDiamond()),
		UpdateTime:  easygo.NewInt64(data.GetUpdateTime()),
	}

	// 新增
	if upWishBox.GetId() == 0 {
		upWishBox.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_ITEM))
		data.CreateTime = easygo.NewInt64(curTime)
		err := col.Insert(upWishBox)
		easygo.PanicError(err)

		return
	}
	data.UpdateTime = easygo.NewInt64(curTime)

	err1 := col.Update(bson.M{"_id": upWishBox.GetId()}, bson.M{"$set": upWishBox})
	easygo.PanicError(err1)
}

// 获取商品品牌列表
func QueryWishGoodsBrandList(req *brower_backstage.ListRequest) ([]*brower_backstage.WishGoodsBrand, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"-HostWeight"}

	queryBson := bson.M{}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	if req.GetKeyword() != "" {
		switch req.GetType() {
		case 1: // 品牌id
			id := easygo.StringToInt64noErr(req.GetKeyword())
			queryBson["_id"] = id
		case 2: // 品牌名称
			queryBson["Name"] = req.GetKeyword()
		default:

		}
	}

	// 状态
	if req.GetStatus() != 1000 {
		queryBson["Status"] = req.GetStatus()
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BRAND)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishBrand
	var errc error

	if req.GetPageSize() == 0 {
		errc = query.Sort(sort...).All(&list)
	} else {
		errc = query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	}
	easygo.PanicError(errc)

	var retList []*brower_backstage.WishGoodsBrand

	for _, v := range list {
		one := &brower_backstage.WishGoodsBrand{
			Id:         easygo.NewInt32(v.GetId()),
			Name:       easygo.NewString(v.GetName()),
			Status:     easygo.NewInt32(v.GetStatus()),
			IsHot:      easygo.NewBool(v.GetIsHot()),
			HotWeight:  easygo.NewInt32(v.GetHostWeight()),
			CreateTime: easygo.NewInt64(v.GetCreateTime()),
			UpdateTime: easygo.NewInt64(v.GetUpdateTime()),
			ClickCount: easygo.NewInt32(v.GetClickCount()),
		}

		retList = append(retList, one)
	}

	return retList, int32(count)
}

// 获取商品列表
func GetWishGoodsBrandList() []*share_message.WishBrand {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BRAND)
	defer closeFun()
	query := col.Find(bson.M{"Status": 1})
	var list []*share_message.WishBrand
	errc := query.All(&list)
	easygo.PanicError(errc)

	return list
}

// 新增/更新商品品牌
func UpdateWishGoodsBrand(data *brower_backstage.WishGoodsBrand) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BRAND)
	defer closeFun()

	upData := &share_message.WishBrand{
		Id:         easygo.NewInt64(data.GetId()),
		Name:       easygo.NewString(data.GetName()),
		Status:     easygo.NewInt32(data.GetStatus()),
		IsHot:      easygo.NewBool(data.GetIsHot()),
		HostWeight: easygo.NewInt32(data.GetHotWeight()),
		Type:       easygo.NewString(data.GetInitial()),
	}

	curTime := time.Now().Unix()
	// 新增
	if upData.GetId() == 0 {
		upData.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_BRAND))
		upData.CreateTime = easygo.NewInt64(curTime)
		err := col.Insert(upData)
		easygo.PanicError(err)

		return
	}

	upData.UpdateTime = easygo.NewInt64(curTime)

	err1 := col.Update(bson.M{"_id": upData.GetId()}, bson.M{"$set": upData})
	easygo.PanicError(err1)
}

// 获取商品类型列表
func QueryWishGoodsTypeList(req *brower_backstage.ListRequest) ([]*brower_backstage.WishGoodsType, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"_id"}

	queryBson := bson.M{}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	if req.GetKeyword() != "" {
		switch req.GetType() {
		case 1: // 类型id
			id := easygo.StringToInt64noErr(req.GetKeyword())
			queryBson["_id"] = id
		case 2: // 类型名称
			queryBson["Name"] = req.GetKeyword()
		default:

		}
	}

	// 状态
	if req.GetStatus() != 1000 {
		queryBson["Status"] = req.GetStatus()
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM_TYPE)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishItemType
	var errc error
	if req.GetPageSize() == 0 {
		errc = query.Sort(sort...).All(&list)
	} else {
		errc = query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	}
	easygo.PanicError(errc)

	var retList []*brower_backstage.WishGoodsType

	for _, v := range list {
		one := &brower_backstage.WishGoodsType{
			Id:          easygo.NewInt32(v.GetId()),
			Name:        easygo.NewString(v.GetName()),
			Status:      easygo.NewInt32(v.GetStatus()),
			CreateTime:  easygo.NewInt64(v.GetCreateTime()),
			UpdateTime:  easygo.NewInt64(v.GetUpdateTime()),
			ClickCount:  easygo.NewInt32(v.GetClickCount()),
			HotWeight:   easygo.NewInt32(v.GetHostWeight()),
			Initial:     easygo.NewString(v.GetType()),
			IsRecommend: easygo.NewBool(v.GetIsRecommend()),
			IsHot:       easygo.NewBool(v.GetIsHot()),
		}

		retList = append(retList, one)
	}

	return retList, int32(count)
}

// 获取商品类型列表
func GetWishGoodsTypeList() []*share_message.WishItemType {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM_TYPE)
	defer closeFun()
	query := col.Find(bson.M{"Status": 1})
	var list []*share_message.WishItemType
	errc := query.All(&list)
	easygo.PanicError(errc)

	return list
}

// 新增/更新商品类型
func UpdateWishGoodsType(data *brower_backstage.WishGoodsType) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM_TYPE)
	defer closeFun()

	upData := &share_message.WishItemType{
		Id:          easygo.NewInt64(data.GetId()),
		Name:        easygo.NewString(data.GetName()),
		Status:      easygo.NewInt32(data.GetStatus()),
		HostWeight:  easygo.NewInt32(data.GetHotWeight()),
		Type:        easygo.NewString(data.GetInitial()),
		IsRecommend: easygo.NewBool(data.GetIsRecommend()),
		IsHot:       easygo.NewBool(data.GetIsHot()),
	}

	curTime := time.Now().Unix()
	upData.UpdateTime = easygo.NewInt64(curTime)
	// 新增
	if upData.GetId() == 0 {
		upData.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_ITEM_TYPE))
		upData.CreateTime = easygo.NewInt64(curTime)
		err := col.Insert(upData)
		easygo.PanicError(err)

		return
	}

	err1 := col.Update(bson.M{"_id": upData.GetId()}, bson.M{"$set": upData})
	easygo.PanicError(err1)
}

// 获取发货订单列表
func QueryWishDeliveryOrderList(req *brower_backstage.ListRequest) ([]*brower_backstage.WishDeliveryOrder, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"-CreateTime"}

	queryBson := bson.M{}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	if req.GetKeyword() != "" {
		switch req.GetType() {
		case 1: // 订单id
			id := easygo.StringToInt64noErr(req.GetKeyword())
			queryBson["_id"] = id
		case 2: // 商品id
			id := easygo.StringToInt64noErr(req.GetKeyword())
			queryBson["ProductId"] = id
		case 3: // 中奖柠檬号
			user := for_game.GetPlayerByAccount(req.GetKeyword())
			wishP := GetWishPlayerInfoByPid(user.GetPlayerId())
			queryBson["PlayerId"] = wishP.GetId()
		default:

		}
	}

	// 状态
	if req.Status != nil {
		switch req.GetStatus() {
		case 0: // 待发货
			queryBson["Status"] = 0
		case 1: // 已发货
			queryBson["Status"] = 1
		case 2: // 已取消
			queryBson["Status"] = 2
		default:

		}

	}

	db, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG)
	defer closeFun()

	var retList []*brower_backstage.WishDeliveryOrder
	var list []*share_message.PlayerExchangeLog

	query := db.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	for _, v := range list {
		one := &brower_backstage.WishDeliveryOrder{
			OrderId:        easygo.NewInt64(v.GetId()),
			GoodsId:        easygo.NewInt64(v.GetProductId()),
			GoodsPrice:     easygo.NewInt64(v.GetProductDiamond()),
			PriceTotal:     easygo.NewInt64(v.GetBoxDrawPrice()),
			CreateTime:     easygo.NewInt64(v.GetCreateTime()),
			DeliveryTime:   easygo.NewInt64(v.GetDeliveryTime()),
			Status:         easygo.NewInt32(v.GetStatus()),
			DeliverName:    easygo.NewString(v.GetReceiver()),
			DeliverPhone:   easygo.NewString(v.GetPhone()),
			DeliverAddress: easygo.NewString(v.GetAddress()),
			UserAccount:    easygo.NewString(v.GetPlayerAccount()),
			UserId:         easygo.NewInt64(v.GetPlayerId()),
			Note:           easygo.NewString(v.GetNote()),
			Company:        easygo.NewString(v.GetCompany()),
			CompanyCode:    easygo.NewString(v.GetCompanyCode()),
			Odd:            easygo.NewString(v.GetOdd()),
			GoodsName:      easygo.NewString(v.GetProductName()),
			UpdateTime:     easygo.NewInt64(v.GetUpdateTime()),
			Operator:       easygo.NewString(v.GetOperator()),
		}

		retList = append(retList, one)
	}

	return retList, int32(count)
}

// 获取回收订单数量
func GetWishRecycleOrderCount(sTime int64) int32 {

	queryBson := bson.M{}
	queryBson["RecycleTime"] = bson.M{"$gte": sTime}

	db, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()

	query := db.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	return int32(count)

}

// 获取发货订单数量
func GetWishDeliveryOrderCount(sTime int64) int32 {

	queryBson := bson.M{"Status": 1}
	queryBson["CreateTime"] = bson.M{"$gte": sTime}

	db, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG)
	defer closeFun()

	query := db.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	return int32(count)

}

// 获取发货订单数量
func GetWishDeliveryOrderCountByBoxItemId(itemId, sTime int64) int32 {

	queryBson := bson.M{"WishBoxItem": itemId, "Status": 1}
	queryBson["CreateTime"] = bson.M{"$gte": sTime}

	db, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG)
	defer closeFun()

	query := db.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	return int32(count)

}

// 获取发货订单人数
func GetWishDeliveryOrderPlayerCount(sTime int64) *share_message.WishPoolReport {

	match := bson.M{"Status": 1}
	match["CreateTime"] = bson.M{"$gte": sTime}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG)
	defer closeFun()

	group := bson.M{
		"_id":                "$PlayerId",
		"DeliverPlayerCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"DeliverPlayerCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	one := &share_message.WishPoolReport{}
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	return one

}

// 获取发货订单人数
func GetWishDeliveryOrderPlayerCountByBoxItemId(itemId, sTime int64) *share_message.WishPoolReport {

	match := bson.M{"WishBoxItem": itemId, "Status": 1}
	match["CreateTime"] = bson.M{"$gte": sTime}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG)
	defer closeFun()

	group := bson.M{
		"_id":                "$PlayerId",
		"DeliverPlayerCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"DeliverPlayerCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	one := &share_message.WishPoolReport{}
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	return one

}

// 获取回收订单人数
func GetWishRecyclePlayerCount(sTime int64) *share_message.WishPoolReport {

	match := bson.M{}
	match["RecycleTime"] = bson.M{"$gte": sTime}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()

	group := bson.M{
		"_id":                "$PlayerId",
		"RecyclePlayerCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"RecyclePlayerCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	one := &share_message.WishPoolReport{}
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	return one

}

// 获取回收订单人数
func GetWishRecyclePlayerCountByBoxItemId(itemId, sTime int64) *share_message.WishPoolReport {

	match := bson.M{}
	match["RecycleTime"] = bson.M{"$gte": sTime}
	match["WishRecycleItem.WishBoxItemId"] = itemId
	match2 := bson.M{}
	match2["WishRecycleItem.WishBoxItemId"] = itemId
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()

	group := bson.M{
		"_id":                "$PlayerId",
		"RecyclePlayerCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"RecyclePlayerCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$match": match2},
		{"$project": project},
	}

	one := &share_message.WishPoolReport{}
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	return one

}

// 兑换钻石人数
func GetWishConvertDiamondPlayerCount(sTime int64) *share_message.WishPoolReport {

	match := bson.M{}
	match["RecycleTime"] = bson.M{"$gte": sTime}
	match["SourceType"] = for_game.DIAMOND_TYPE_EXCHANGE_IN
	match["PayType"] = 1

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_DIAMOND_RECHARGE)
	defer closeFun()

	group := bson.M{
		"_id":                "$PlayerId",
		"RecyclePlayerCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"RecyclePlayerCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	one := &share_message.WishPoolReport{}
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	return one

}

// 兑换钻石总金额+次数
func GetWishConvertDiamond(sTime int64) *share_message.WishPoolReport {

	match := bson.M{}
	match["RecycleTime"] = bson.M{"$gte": sTime}
	match["SourceType"] = for_game.DIAMOND_TYPE_EXCHANGE_IN
	match["PayType"] = 1

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_DIAMOND_RECHARGE)
	defer closeFun()

	group := bson.M{
		"_id":                 nil,
		"ConvertDiamondTotal": bson.M{"$sum": "$ChangeDiamond"},
		"ConvertDiamondCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"ConvertDiamondTotal": 1,
		"ConvertDiamondCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	one := &share_message.WishPoolReport{}
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	return one

}

// 更新订单状态
func UpdateDeliveryOrderStatus(data *brower_backstage.UpdateStatusRequest, operator string) {
	db, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG)
	defer closeFun()

	status := 0
	switch data.GetStatus() {
	case 1:
		status = 1
	case 2:
		status = 2
	case 3:
		status = 3
	default:

	}
	note := data.GetNote()
	curTime := time.Now().Unix()
	upData := bson.M{
		"Status":     easygo.NewInt32(status),
		"Note":       easygo.NewString(note),
		"UpdateTime": easygo.NewInt64(curTime),
		"Operator":   easygo.NewString(operator),
	}

	err := db.Update(bson.M{"_id": data.GetId()}, bson.M{"$set": upData})
	easygo.PanicError(err)

}

// 填写发货信息
func UpdateDeliveryOrderCourierInfo(data *brower_backstage.UpdateDeliveryOrderCourierInfo, operator string) {
	db, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG)
	defer closeFun()

	curTime := time.Now().Unix()
	upData := &share_message.PlayerExchangeLog{
		Id:           easygo.NewInt64(data.GetOrderId()),
		Company:      easygo.NewString(data.GetCompany()),
		CompanyCode:  easygo.NewString(data.GetCompanyCode()),
		Odd:          easygo.NewString(data.GetOdd()),
		Status:       easygo.NewInt32(1),
		DeliveryTime: easygo.NewInt64(curTime),
		UpdateTime:   easygo.NewInt64(curTime),
		Operator:     easygo.NewString(operator),
	}

	err := db.Update(bson.M{"_id": upData.GetId()}, bson.M{"$set": upData})
	easygo.PanicError(err)
}

//获取兑换信息
func GetWishPlayerExchangeLog(id int64) *share_message.PlayerExchangeLog {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG)
	defer closeFun()
	player := &share_message.PlayerExchangeLog{}
	err := col.Find(bson.M{"_id": id}).One(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return player
}

func GetWishPlayerExchangeLogListByOrderId(id string) []*share_message.PlayerExchangeLog {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG)
	defer closeFun()

	var list []*share_message.PlayerExchangeLog
	queryBson := bson.M{"OrderId": id}
	query := col.Find(queryBson)
	err := query.All(&list)
	easygo.PanicError(err)

	return list
}

//获取用户信息
func GetWishRecycleOrder(id int64) *share_message.WishRecycleOrder {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()
	player := &share_message.WishRecycleOrder{}
	err := col.Find(bson.M{"_id": id}).One(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return player
}

//获取用户信息
func GetWishRecycleOrderByPaymentOrderId(id string) *share_message.WishRecycleOrder {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()
	player := &share_message.WishRecycleOrder{}
	err := col.Find(bson.M{"PaymentOrderId": id}).One(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return player
}

// 更新订单状态
func UpdateWishRecycleOrderStatus(data *brower_backstage.UpdateStatusRequest, operator string) {
	db, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()

	status := int32(0)
	switch data.GetStatus() {
	case 0, 1, 2, 3:
		status = data.GetStatus()
	}

	note := data.GetNote()
	curTime := time.Now().Unix()
	upData := &share_message.WishRecycleOrder{
		Id:          easygo.NewInt64(data.GetId()),
		Status:      easygo.NewInt32(status),
		RefusalNote: easygo.NewString(note),
		Operator:    easygo.NewString(operator),
	}

	switch status {
	case 1:
		upData.RecycleTime = easygo.NewInt64(curTime)
		upData.UpdateTime = easygo.NewInt64(curTime)
	case 2:
		upData.UpdateTime = easygo.NewInt64(curTime)
	default:

	}

	err := db.Update(bson.M{"_id": upData.GetId()}, bson.M{"$set": upData})
	easygo.PanicError(err)

}

// 获取回收订单列表
func QueryWishRecycleOrderList(req *brower_backstage.ListRequest) ([]*brower_backstage.WishRecycleOrder, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"-InitTime"}

	queryBson := bson.M{"Channel": 2}
	if req.GetBeginTimestamp() != 0 {
		queryBson["RecycleTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	if req.GetKeyword() != "" {
		switch req.GetType() {
		case 1: // 类型id
			id := easygo.StringToInt64noErr(req.GetKeyword())
			queryBson["_id"] = id
		case 2: // 类型名称
			id := easygo.StringToInt64noErr(req.GetKeyword())
			queryBson["RecycleItemList.ProductId"] = id
		case 3: // 回收账号
			user := for_game.GetPlayerByAccount(req.GetKeyword())
			queryBson["UserId"] = user.GetPlayerId()
		default:

		}
	}

	// 状态
	if req.ListType != nil && req.GetListType() != 1000 {
		switch req.GetListType() {
		case 0: // 待审核
			queryBson["Status"] = 0
		case 1: // 已回收
			queryBson["Status"] = 1
		case 2: // 已拒绝
			queryBson["Status"] = 2
		default:

		}

	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishRecycleOrder
	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	var retList []*brower_backstage.WishRecycleOrder

	for _, v := range list {
		one := &brower_backstage.WishRecycleOrder{
			OrderId:        easygo.NewInt64(v.GetId()),
			RecyclePrice:   easygo.NewInt64(v.GetRecyclePriceTotal()),
			RecycleTime:    easygo.NewInt64(v.GetRecycleTime()),
			UserAccount:    easygo.NewString(v.GetUserAccount()),
			UserId:         easygo.NewInt64(v.GetUserId()),
			Type:           easygo.NewInt32(v.GetType()),
			UpdateTime:     easygo.NewInt64(v.GetUpdateTime()),
			RefusalNote:    easygo.NewString(v.GetRefusalNote()),
			InitTime:       easygo.NewInt64(v.GetInitTime()),
			Status:         easygo.NewInt32(v.GetStatus()),
			RecycleDiamond: easygo.NewInt64(v.GetRecycleDiamond()),
			Operator:       easygo.NewString(v.GetOperator()),
			Note:           easygo.NewString(v.GetRemarks()),
		}

		reason := GetWishRecycleReason(v.GetRecycleNote())
		one.UserReason = easygo.NewString(reason.GetReason())

		retList = append(retList, one)
	}

	return retList, int32(count)
}

//获取用户信息
func GetWishRecycleReason(id int32) *share_message.RecycleReason {
	col, closeFun := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_REASON)
	defer closeFun()
	player := &share_message.RecycleReason{}
	err := col.Find(bson.M{"_id": id}).One(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return player
}

// 获取回收订单列表
func GetWishRecycleItemListByPaymentOrderId(str string) []*share_message.WishRecycleItem {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()

	var one *share_message.WishRecycleOrder
	//str := strconv.FormatInt(id,10)
	queryBson := bson.M{"PaymentOrderId": str}
	query := col.Find(queryBson)
	err := query.One(&one)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}

	return one.GetRecycleItemList()
}

func GetWishRecycleItemList(id int64) []*share_message.WishRecycleItem {
	if id == 0 {
		return nil
	}
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()

	var one *share_message.WishRecycleOrder
	queryBson := bson.M{"_id": id}
	query := col.Find(queryBson)
	err := query.One(&one)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}

	return one.GetRecycleItemList()
}

//订单列表查询
func QueryWishOrderList(reqMsg *brower_backstage.QueryWishOrderRequest) ([]*share_message.Order, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ORDER)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{"OrderType": 1}
	if reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {

		switch reqMsg.GetTimeType() {
		case 1:
			queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
		case 2:
			queryBson["OverTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
		default:

		}
	}

	//不查关键词
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetTypes() {
		case 1: //订单号
			queryBson["_id"] = reqMsg.GetKeyword()
		case 2: //柠檬号
			player := QueryPlayerbyAccount(reqMsg.GetKeyword())
			queryBson["PlayerId"] = player.GetPlayerId()
		case 3:
			queryBson["Operator"] = reqMsg.GetKeyword()
		case 4:
			queryBson["AccountNo"] = reqMsg.GetKeyword()
		default:
			easygo.NewFailMsg("查询条件有误")
		}
	}

	if reqMsg.PayStatus != nil && reqMsg.GetPayStatus() != 1000 {
		queryBson["PayStatus"] = reqMsg.GetPayStatus()

	}

	//if reqMsg.SourceType != nil && reqMsg.GetSourceType() > 0 {
	//	// types := reqMsg.GetSourceType()
	//	// queryBson["SourceType"] = bson.M{"$in": types}
	//	queryBson["SourceType"] = reqMsg.GetSourceType()
	//
	//}
	//
	//if reqMsg.PayType != nil && reqMsg.GetPayType() != 0 {
	//	queryBson["PayChannel"] = reqMsg.GetPayType()
	//}

	if reqMsg.Status != nil && reqMsg.GetStatus() != 1000 {
		queryBson["Status"] = reqMsg.GetStatus()
	}

	//if reqMsg.ChangeType != nil && reqMsg.GetChangeType() != 0 {
	//	queryBson["ChangeType"] = reqMsg.GetChangeType()
	//}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.Order
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

// 获取许愿池报表列表
func QueryWishPoolReportList(req *brower_backstage.ListRequest) ([]*share_message.WishPoolReport, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"-_id"}

	queryBson := bson.M{}
	if req.GetBeginTimestamp() != 0 {
		queryBson["_id"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	//默认日报
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_POOL)
	timeType := req.GetTimeType()
	switch timeType {
	case 7:
		col, closeFun = MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_POOL_WEEK)
	case 30:
		col, closeFun = MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_POOL_MONTH)
	}

	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishPoolReport
	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)
}

//获取ruleId对应的许愿池-规则 map
func GetAllActPoolRulesName(typesInt []int32) map[int64]string {
	poolMap := make(map[int64]string, 0)
	wishPools, countPool := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL, bson.M{}, 0, 0)
	if countPool > 0 {
		for _, pool := range wishPools {
			poolId := pool.(bson.M)["_id"].(int64)
			poolMap[poolId] = pool.(bson.M)["Name"].(string)
		}
	}
	res := make(map[int64]string, 0)
	rules, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ACT_POOL_RULE, bson.M{"Type": bson.M{"$in": typesInt}}, 0, 0)
	if count > 0 {
		for _, li := range rules {
			dest := &share_message.WishActPoolRule{}
			for_game.StructToOtherStruct(li, dest)
			var str string
			if dest.GetType() == 1 {
				str = fmt.Sprintf("%d次", dest.GetKey())
			} else if dest.GetType() == 2 {
				str = fmt.Sprintf("%d天", dest.GetKey())
			}
			res[dest.GetId()] = fmt.Sprintf("%s-%s", poolMap[dest.GetWishActPoolId()], str)
		}
	}
	return res
}

//许愿池活动列表
func QueryWishPoolActivityReportList(req *brower_backstage.ListRequest) ([]*share_message.WishActivityReport, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"-_id"}

	queryBson := bson.M{}
	if req.GetBeginTimestamp() != 0 {
		queryBson["_id"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_ACTIVITY)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishActivityReport
	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)
	rulesNameMap := GetAllActPoolRulesName([]int32{1, 2})
	for k, v := range list {
		//次数
		counterData := v.GetCounterData()
		for k1, v1 := range counterData {
			counterData[k1].Key = easygo.NewString(rulesNameMap[v1.GetWishActPoolRuleId()])
		}
		//天数
		dayCountData := v.GetDayCountData()
		for k2, v2 := range dayCountData {
			dayCountData[k2].Key = easygo.NewString(rulesNameMap[v2.GetWishActPoolRuleId()])
		}
		list[k].CounterData = counterData
		list[k].DayCountData = dayCountData
	}
	return list, int32(count)
}

// 获取许愿池盲盒报表列表
func QueryWishBoxReportList(req *brower_backstage.ListRequest) ([]*share_message.WishBoxReport, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"-CreateTime"}

	queryBson := bson.M{"CreateTime": req.GetId()}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX)
	timeType := req.GetTimeType()
	switch timeType {
	case 7:
		col, closeFun = MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_WEEK)
	case 30:
		col, closeFun = MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_MONTH)
	}
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishBoxReport
	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)
}

// 获取许愿池盲盒报表列表
func QueryWishBoxReportTempList(req *brower_backstage.ListRequest) ([]*share_message.WishBoxReport, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"-CreateTime"}

	queryBson := bson.M{"CreateTime": req.GetId()}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_TEMP)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishBoxReport
	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)
}

// 删除许愿池盲盒报表列表
func DeleteWishBoxReportTempList(req *brower_backstage.ListRequest) {
	queryBson := bson.M{"CreateTime": req.GetId()}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}
	queryBson["WishBoxId"] = req.GetType()
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_TEMP)
	defer closeFun()

	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)

}

// 获取许愿池盲盒详情报表列表
func QueryWishBoxDetailReportList(req *brower_backstage.ListRequest) ([]*share_message.WishBoxDetailReport, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"ItemId"}

	queryBson := bson.M{}
	queryBson["CreateTime"] = req.GetId()
	queryBson["WishBoxId"] = req.GetType()

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_DETAIL)
	timeType := req.GetTimeType()
	switch timeType {
	case 7:
		col, closeFun = MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_DETAIL_WEEK)
	case 30:
		col, closeFun = MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_DETAIL_MONTH)
	}
	switch timeType {

	}
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishBoxDetailReport
	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)
}

// 获取许愿池盲盒详情报表列表
func QueryWishBoxDetailReportTempList(req *brower_backstage.ListRequest) ([]*share_message.WishBoxDetailReport, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"ItemId"}

	queryBson := bson.M{}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	queryBson["WishBoxId"] = req.GetType()

	if req.GetKeyword() != "" {
		switch req.GetType() {
		case 1: // 盲盒报表id
			boxId := easygo.StringToInt64noErr(req.GetKeyword())
			queryBson["WishBoxId"] = boxId
		default:

		}
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_DETAIL_TEMP)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishBoxDetailReport
	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)
}

// 获取许愿池盲盒详情报表列表
func DeleteWishBoxDetailReportTempList(req *brower_backstage.ListRequest) {
	queryBson := bson.M{}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}
	queryBson["WishBoxId"] = req.GetType()
	if req.GetKeyword() != "" {
		switch req.GetType() {
		case 1: // 盲盒报表id
			boxId := easygo.StringToInt64noErr(req.GetKeyword())
			queryBson["WishBoxId"] = boxId
		default:

		}
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_DETAIL_TEMP)
	defer closeFun()

	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)

}

// 获取许愿池商品报表列表
func QueryWishItemReportList(req *brower_backstage.ListRequest) ([]*share_message.WishItemReport, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"ItemId"}

	queryBson := bson.M{"CreateTime": req.GetId()}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_ITEM)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishItemReport
	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)
}

// 获取玩家物品列表(挑战产生物品)列表
func QueryPlayerWishItemList(req *brower_backstage.ListRequest) ([]*share_message.PlayerWishItem, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"-CreateTime"}

	queryBson := bson.M{"WishBoxId": req.GetId()}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.PlayerWishItem
	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)
}

// 获取水池流水日志列表
func QueryWishPoolLogList(req *brower_backstage.ListRequest) ([]*share_message.WishPoolLog, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"-CreateTime"}

	queryBson := bson.M{"BoxId": req.GetId()}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL_LOG)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishPoolLog
	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)
}

// 获取水池抽水日志列表
func QueryWishPoolPumpLogList(req *brower_backstage.ListRequest) ([]*share_message.WishPoolPumpLog, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"-CreateTime"}

	queryBson := bson.M{"BoxId": req.GetId()}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL_PUMP_LOG)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishPoolPumpLog
	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)
}

// 获取用户收藏盲盒数量
func GetPlayerWishCollectionCountByPid(id int64) int32 {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM) // TODO
	defer closeFun()

	queryBson := bson.M{"PlayerId": id}
	count, err := col.Find(queryBson).Count()
	easygo.PanicError(err)

	return int32(count)
}

// 获取用户许愿物品数量
func GetPlayerWishDataCountByPid(id int64) int32 {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA)
	defer closeFun()

	queryBson := bson.M{"PlayerId": id}
	queryBson["Status"] = bson.M{"$in": []int{0, 1}}
	count, err := col.Find(queryBson).Count()
	easygo.PanicError(err)

	return int32(count)
}

// 获取抽奖记录列表
func QueryDrawRecordList(req *brower_backstage.ListRequest) *brower_backstage.DrawRecordList {

	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG)
	defer closeFun()

	//lookup := bson.M{"from": for_game.TABLE_WISH_PLAYER, "localField": "DareId", "foreignField": "_id", "as": "user"}
	//lookupAdd := bson.M{"from": for_game.TABLE_PLAYER_WISH_COLLECTION, "localField": "DareId", "foreignField": "PlayerId", "as": "add"}
	//lookupWish := bson.M{"from": for_game.TABLE_PLAYER_WISH_DATA, "localField": "DareId", "foreignField": "PlayerId", "as": "wish"}
	//cond1 := []interface{}{bson.M{"$eq": []interface{}{"$wish.PlayerId", "$DareId"}}, 1, 0}
	match := bson.M{}

	if req.GetKeyword() != "" {
		switch req.GetType() {
		case 1: // 柠檬号
			user := for_game.GetPlayerByAccount(req.GetKeyword())
			wishP := GetWishPlayerInfoByPid(user.GetPlayerId())
			match["DareId"] = wishP.GetId()
		case 2: // 用户昵称
			user := for_game.GetPlayerByNickName(req.GetKeyword())
			wishP := GetWishPlayerInfoByPid(user.GetPlayerId())
			match["DareId"] = wishP.GetId()
		case 3: // 手机号码
			user := for_game.GetPlayerByPhone(req.GetKeyword())
			wishP := GetWishPlayerInfoByPid(user.GetPlayerId())
			match["DareId"] = wishP.GetId()
		default:

		}
	}

	if req.GetUserType() > 0 {
		var ids []int64
		//1普通用户 2运营号 3白名单
		switch req.GetUserType() {
		case 1:
			wplis, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_PLAYER, bson.M{"Types": 0}, 0, 0)
			for _, p := range wplis {
				ids = append(ids, p.(bson.M)["_id"].(int64))
			}
			match["DareId"] = bson.M{"$in": ids}
		case 2:
			wplis, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_PLAYER, bson.M{"Types": bson.M{"$gte": 2}}, 0, 0)
			for _, p := range wplis {
				ids = append(ids, p.(bson.M)["_id"].(int64))
			}
			match["DareId"] = bson.M{"$in": ids}
		case 3:
			players := GetAllIdsInWishWhite()
			match["DareId"] = bson.M{"$in": players}
		}
	}

	if req.GetBeginTimestamp() != 0 {
		match["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	group := bson.M{
		"_id":           "$DareId",
		"DrawGoldTotal": bson.M{"$sum": "$DarePrice"},
		"LastDrawTime":  bson.M{"$last": "$CreateTime"},
		"DrawCount":     bson.M{"$sum": 1},
		//"AddBoxCount":   bson.M{"$sum": "$add._id"},
		//"WishItemCount": bson.M{"$sum": bson.M{"$cond":cond1}},
		"UserId": bson.M{"$last": "$DareId"},
		//"UserAccount":  bson.M{"$last": "$user.Account"},
		//"UserNickname": bson.M{"$last": "$user.NickName"},
		//"Phone":        bson.M{"$last": "$user.Phone"},
	}

	project := bson.M{
		"UserId": 1,
		//"UserAccount":   1,
		//"UserNickname":  1,
		//"Phone":         1,
		"DrawGoldTotal": 1,
		"LastDrawTime":  1,
		"DrawCount":     1,
		//"AddBoxCount":   1,
		//"WishItemCount": 1,
	}

	countType := "PageCount"

	countCond := []bson.M{
		{"$match": match},
		//{"$lookup": lookup},
		//{"$unwind": "$user"},
		//{"$lookup": lookupAdd},
		//{"$unwind": "$add"},
		//{"$lookup": lookupWish},
		//{"$unwind": "$wish"},
		{"$group": group},
		{"$count": countType},
	}

	retList := &brower_backstage.DrawRecordList{}
	err := col.Pipe(countCond).One(&retList)

	sort := bson.M{"LastDrawTime": -1}
	if req.GetSort() != "" {
		if strings.Contains(req.GetSort(), "-") {
			s := strings.TrimLeft(req.GetSort(), "-")
			sort = bson.M{s: -1}
		} else {
			sort = bson.M{req.GetSort(): 1}
		}
	}

	pipeCond := []bson.M{
		{"$match": match},
		//{"$lookup": lookup},
		//{"$unwind": "$user"},
		//{"$lookup": lookupAdd},
		//{"$unwind": "$add"},
		//{"$lookup": lookupWish},
		//{"$unwind": "$wish"},
		{"$group": group},
		{"$project": project},
		{"$sort": sort},
		{"$skip": curPage * pageSize},
		{"$limit": pageSize},
	}

	var list []*brower_backstage.DrawRecord
	err = col.Pipe(pipeCond).All(&list)
	easygo.PanicError(err)

	retList.List = list
	return retList
}

// 根据用户id获取收藏盲盒数
func GetAddBoxCountByUserId(id int64) int32 {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_COLLECTION)
	defer closeFun()

	match := bson.M{"PlayerId": id}
	count, err := col.Find(match).Count()
	easygo.PanicError(err)

	return int32(count)
}

// 根据用户id获取收藏盲盒数
func GetAddBoxCountByUserIds(id []int64) []*brower_backstage.DrawRecord {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_COLLECTION)
	defer closeFun()

	match := bson.M{"PlayerId": bson.M{"$in": id}}
	group := bson.M{
		"_id":         "$PlayerId",
		"AddBoxCount": bson.M{"$sum": 1},
		"UserId":      bson.M{"$last": "$PlayerId"},
	}

	project := bson.M{
		"UserId":      1,
		"AddBoxCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var counts []*brower_backstage.DrawRecord
	err := col.Pipe(pipeCond).All(&counts)
	easygo.PanicError(err)

	return counts
}

// 根据用户id获取许愿数
func GetWishDataCountByUserId(id int64) int32 {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA)
	defer closeFun()

	//lookupBox := bson.M{"from": for_game.TABLE_WISH_BOX, "localField": "WishBoxId", "foreignField": "_id", "as": "box"}
	//lookupItem := bson.M{"from": for_game.TABLE_WISH_ITEM, "localField": "WishItemId", "foreignField": "_id", "as": "item"}
	match := bson.M{"PlayerId": id}

	countType := "PageCount"

	countCond := []bson.M{
		{"$match": match},
		//{"$lookup": lookupBox},
		//{"$unwind": "$box"},
		//{"$lookup": lookupItem},
		//{"$unwind": "$item"},
		{"$count": countType},
	}

	retList := &brower_backstage.WishGoodsRecordList{}
	err := col.Pipe(countCond).One(&retList)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return int32(0)
	}

	return retList.GetPageCount()
}

// 根据用户id获取收藏盲盒数
func GetWishDataCountByUserIds(id []int64) []*brower_backstage.DrawRecord {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA)
	defer closeFun()

	match := bson.M{"PlayerId": bson.M{"$in": id}}
	group := bson.M{
		"_id":           "$PlayerId",
		"UserId":        bson.M{"$last": "$PlayerId"},
		"WishItemCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"UserId":        1,
		"WishItemCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var counts []*brower_backstage.DrawRecord
	err := col.Pipe(pipeCond).All(&counts)
	easygo.PanicError(err)

	return counts
}

// 根据用户id获取现有物品
func GetPlayerWishItem(id int64) int32 {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	match := bson.M{"PlayerId": id, "Status": 0}

	countType := "PageCount"

	countCond := []bson.M{
		{"$match": match},
		{"$count": countType},
	}

	retList := &brower_backstage.WishGoodsRecordList{}
	err := col.Pipe(countCond).One(&retList)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return int32(0)
	}

	return retList.GetPageCount()
}

// 根据用户id获取现有物品数量
func GetPlayerWishItems(id []int64) []*brower_backstage.DrawRecord {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()
	match := bson.M{"PlayerId": bson.M{"$in": id}, "Status": 0}

	group := bson.M{
		"_id":           "$PlayerId",
		"UserId":        bson.M{"$last": "$PlayerId"},
		"HaveItemCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"UserId":        1,
		"HaveItemCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var counts []*brower_backstage.DrawRecord
	err := col.Pipe(pipeCond).All(&counts)
	easygo.PanicError(err)

	return counts
}

// 根据用户id获取扣除物品数量
func GetPlayerDelWishItemsCount(id []int64) []*brower_backstage.DrawRecord {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()
	match := bson.M{"PlayerId": bson.M{"$in": id}, "Status": 4}

	group := bson.M{
		"_id":          "$PlayerId",
		"UserId":       bson.M{"$last": "$PlayerId"},
		"DelItemCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"UserId":       1,
		"DelItemCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var counts []*brower_backstage.DrawRecord
	err := col.Pipe(pipeCond).All(&counts)
	easygo.PanicError(err)

	return counts
}

// 获取收藏盲盒记录列表
func QueryAddBoxRecordList(req *brower_backstage.ListRequest) *brower_backstage.AddBoxRecordList {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())

	match := bson.M{"PlayerId": req.GetId()}
	if req.GetBeginTimestamp() != 0 {
		match["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_COLLECTION)
	defer closeFun()

	lookup := bson.M{"from": for_game.TABLE_WISH_BOX, "localField": "WishBoxId", "foreignField": "_id", "as": "box"}

	countType := "PageCount"

	countCond := []bson.M{
		{"$match": match},
		{"$lookup": lookup},
		{"$unwind": "$box"},
		{"$count": countType},
	}

	retList := &brower_backstage.AddBoxRecordList{}
	err := col.Pipe(countCond).One(&retList)

	project := bson.M{
		"UserId":     "$PlayerId",
		"BoxId":      "$WishBoxId",
		"BoxName":    "$box.Name",
		"CreateTime": "$CreateTime",
	}
	pipeCond := []bson.M{
		{"$match": match},
		{"$lookup": lookup},
		{"$unwind": "$box"},
		{"$project": project},
		{"$sort": bson.M{"CreateTime": -1}},
		{"$skip": curPage * pageSize},
		{"$limit": pageSize},
	}

	var list []*brower_backstage.AddBoxRecord
	err = col.Pipe(pipeCond).All(&list)
	easygo.PanicError(err)
	retList.List = list
	return retList
}

// 获取许愿物品记录列表
func QueryWishGoodsRecordList(req *brower_backstage.ListRequest) *brower_backstage.WishGoodsRecordList {

	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA)
	defer closeFun()

	//lookupBox := bson.M{"from": for_game.TABLE_WISH_BOX, "localField": "WishBoxId", "foreignField": "_id", "as": "box"}
	//lookupItem := bson.M{"from": for_game.TABLE_WISH_ITEM, "localField": "WishItemId", "foreignField": "_id", "as": "item"}
	match := bson.M{"PlayerId": req.GetId()}

	countType := "PageCount"

	countCond := []bson.M{
		{"$match": match},
		//{"$lookup": lookupBox},
		//{"$unwind": "$box"},
		//{"$lookup": lookupItem},
		//{"$unwind": "$item"},
		{"$count": countType},
	}

	retList := &brower_backstage.WishGoodsRecordList{}
	err := col.Pipe(countCond).One(&retList)

	project := bson.M{
		"GoodsId":    "$WishItemId",
		"BoxId":      "$WishBoxId",
		"CreateTime": "$CreateTime",
	}
	pipeCond := []bson.M{
		{"$match": match},
		//{"$lookup": lookupBox},
		//{"$unwind": "$box"},
		//{"$lookup": lookupItem},
		//{"$unwind": "$item"},
		{"$project": project},
		{"$sort": bson.M{"CreateTime": -1}},
		{"$skip": curPage * pageSize},
		{"$limit": pageSize},
	}

	var list []*brower_backstage.WishGoodsRecord
	err = col.Pipe(pipeCond).All(&list)
	easygo.PanicError(err)

	for i := range list {
		box := GetWishBoxById(list[i].GetBoxId())
		list[i].BoxName = easygo.NewString(box.GetName())
		item := GetWishItemById(list[i].GetGoodsId())
		list[i].GoodsName = easygo.NewString(item.GetName())
	}

	retList.List = list
	return retList
}

// 盲盒抽奖记录列表
func DrawBoxRecordList(req *brower_backstage.ListRequest) *brower_backstage.DrawBoxRecordList {

	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG)
	defer closeFun()

	match := bson.M{"DareId": req.GetId()}

	if req.GetEndTimestamp() > 0 {
		match["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	group := bson.M{
		"_id":       "$WishBoxId",
		"DrawCount": bson.M{"$sum": 1},
		"LastTime":  bson.M{"$last": "$CreateTime"},
		"BoxId":     bson.M{"$last": "$WishBoxId"},
		"UserId":    bson.M{"$last": "$DareId"},
	}

	project := bson.M{
		"UserId":    1,
		"BoxId":     1,
		"LastTime":  1,
		"DrawCount": 1,
	}

	countType := "PageCount"

	countCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$count": countType},
	}

	retList := &brower_backstage.DrawBoxRecordList{}
	err := col.Pipe(countCond).One(&retList)

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
		{"$sort": bson.M{"CreateTime": -1}},
		{"$skip": curPage * pageSize},
		{"$limit": pageSize},
	}

	var list []*brower_backstage.DrawBoxRecord
	err = col.Pipe(pipeCond).All(&list)
	easygo.PanicError(err)
	retList.List = list
	return retList
}

// 盲盒抽奖记录列表
func PlayerWishItemList(req *brower_backstage.ListRequest) *brower_backstage.HaveItemList {

	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	match := bson.M{"PlayerId": req.GetId(), "Status": 0}

	if req.GetEndTimestamp() > 0 {
		match["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	project := bson.M{
		"GoodsId":      "$WishItemId",
		"GoodsName":    "$ProductName",
		"BoxName":      "$BoxName",
		"BoxId":        "$WishBoxId",
		"CreateTime":   "$CreateTime",
		"PlayerItemId": "$_id",
	}

	countType := "PageCount"

	countCond := []bson.M{
		{"$match": match},
		{"$count": countType},
	}

	retList := &brower_backstage.HaveItemList{}
	err := col.Pipe(countCond).One(&retList)

	pipeCond := []bson.M{
		{"$match": match},
		{"$project": project},
		{"$sort": bson.M{"CreateTime": -1}},
		{"$skip": curPage * pageSize},
		{"$limit": pageSize},
	}

	if pageSize == 0 {
		pipeCond = []bson.M{
			{"$match": match},
			{"$project": project},
			{"$sort": bson.M{"CreateTime": -1}},
		}
	}

	var list []*brower_backstage.HaveItem
	err = col.Pipe(pipeCond).All(&list)
	easygo.PanicError(err)
	retList.List = list
	return retList
}

// 用户扣除物品记录列表
func PlayerWishDelItemList(req *brower_backstage.ListRequest) *brower_backstage.HaveItemList {

	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	match := bson.M{"PlayerId": req.GetId(), "Status": 4}

	if req.GetEndTimestamp() > 0 {
		match["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	project := bson.M{
		"GoodsId":      "$WishItemId",
		"GoodsName":    "$ProductName",
		"BoxName":      "$BoxName",
		"BoxId":        "$WishBoxId",
		"CreateTime":   "$CreateTime",
		"PlayerItemId": "$_id",
		"UpdateTime":   "$UpdateTime",
		"Operator":     "$Operator",
	}

	countType := "PageCount"

	countCond := []bson.M{
		{"$match": match},
		{"$count": countType},
	}

	retList := &brower_backstage.HaveItemList{}
	err := col.Pipe(countCond).One(&retList)

	pipeCond := []bson.M{
		{"$match": match},
		{"$project": project},
		{"$sort": bson.M{"CreateTime": -1}},
		{"$skip": curPage * pageSize},
		{"$limit": pageSize},
	}

	if pageSize == 0 {
		pipeCond = []bson.M{
			{"$match": match},
			{"$project": project},
			{"$sort": bson.M{"CreateTime": -1}},
		}
	}

	var list []*brower_backstage.HaveItem
	err = col.Pipe(pipeCond).All(&list)
	easygo.PanicError(err)
	retList.List = list
	return retList
}

// 扣除用户抽取的物品
func DeletePlayerWishItem(ids []int64, userid string) {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": ids}, "Status": 0}, bson.M{"$set": bson.M{"Status": 4, "UpdateTime": easygo.NowTimestamp(), "Operator": userid}})
	easygo.PanicError(err)
}

// 扣除用户所有的物品
func DeletePlayerWishItemByPid(ids int64, userid string) {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	_, err := col.UpdateAll(bson.M{"PlayerId": bson.M{"$in": ids}, "Status": 0}, bson.M{"$set": bson.M{"Status": 4, "Operator": userid}})
	easygo.PanicError(err)
}

// 中奖记录记录列表
func QueryWinRecordList(req *brower_backstage.ListRequest) *brower_backstage.WinRecordList {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	uid := req.GetId()

	//lookupItem := bson.M{"from": for_game.TABLE_WISH_ITEM, "localField": "WishItemId", "foreignField": "_id", "as": "item"}
	match := bson.M{"PlayerId": uid}

	if req.GetType() != 0 {
		match["WishBoxId"] = req.GetType()
	}

	if req.GetEndTimestamp() > 0 {
		match["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	countType := "PageCount"

	countCond := []bson.M{
		{"$match": match},
		//{"$lookup": lookupItem},
		//{"$unwind": "$item"},
		{"$count": countType},
	}

	retList := &brower_backstage.WinRecordList{}
	err := col.Pipe(countCond).One(&retList)

	project := bson.M{
		"UserId":  "$PlayerId",
		"GoodsId": "$WishItemId",
		//"GoodsName":  "$item.Name",
		"CreateTime": "$CreateTime",
		"BoxItemId":  "$ChallengeItemId",
	}
	pipeCond := []bson.M{
		{"$match": match},
		//{"$lookup": lookupItem},
		//{"$unwind": "$item"},
		{"$project": project},
		{"$sort": bson.M{"CreateTime": -1}},
		{"$skip": curPage * pageSize},
		{"$limit": pageSize},
	}

	var list []*brower_backstage.WinRecord
	err = col.Pipe(pipeCond).All(&list)
	easygo.PanicError(err)

	retList.List = list
	return retList
}

var divisor = int64(100)

// 获取水池列表
func QueryWishPoolList(req *brower_backstage.ListRequest) ([]*brower_backstage.WishPool, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"-CreateTime"}

	queryBson := bson.M{}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	if req.GetKeyword() != "" {
		switch req.GetType() {
		case 1: // 水池id
			id := easygo.StringToInt64noErr(req.GetKeyword())
			queryBson["_id"] = id
		case 2: // 水池名称
			queryBson["Name"] = req.GetKeyword()
		default:

		}
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL_CFG)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishPool
	var errc error

	if req.GetPageSize() == 0 {
		errc = query.Sort(sort...).All(&list)
	} else {
		errc = query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	}
	easygo.PanicError(errc)

	var retList []*brower_backstage.WishPool

	for _, v := range list {
		one := &brower_backstage.WishPool{
			Id:               easygo.NewInt64(v.GetId()),
			PoolLimit:        easygo.NewInt32(v.GetPoolLimit()),
			Name:             easygo.NewString(v.GetName()),
			CreateTime:       easygo.NewInt64(v.GetCreateTime()),
			ShowInitialValue: easygo.NewInt32(v.GetShowInitialValue()),
			ShowRecycle:      easygo.NewInt32(v.GetShowRecycle()),
			ShowCommission:   easygo.NewInt32(v.GetShowCommission()),
			ShowStartAward:   easygo.NewInt32(v.GetShowStartAward()),
			ShowCloseAward:   easygo.NewInt32(v.GetShowCloseAward()),
			IsDefault:        easygo.NewBool(v.GetIsDefault()),
		}

		one.SmallLoss = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(v.SmallLoss.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(v.SmallLoss.GetShowMinValue()),
		}

		one.SmallWin = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(v.SmallWin.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(v.SmallWin.GetShowMinValue()),
		}

		one.BigLoss = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(v.BigLoss.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(v.BigLoss.GetShowMinValue()),
		}

		one.BigWin = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(v.BigWin.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(v.BigWin.GetShowMinValue()),
		}
		one.Common = &brower_backstage.WishPoolStatus{
			ShowMaxValue: easygo.NewInt32(v.Common.GetShowMaxValue()),
			ShowMinValue: easygo.NewInt32(v.Common.GetShowMinValue()),
		}

		retList = append(retList, one)
	}

	return retList, int32(count)
}

// 获取水池列表
func QueryWishPoolIds(req *brower_backstage.ListRequest) ([]int64, int32) {
	pageSize, curPage := SetMgoPage(req.GetPageSize(), req.GetCurPage())
	sort := []string{"-CreateTime"}

	queryBson := bson.M{}
	if req.GetBeginTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": req.GetBeginTimestamp(), "$lte": req.GetEndTimestamp()}
	}

	if req.GetKeyword() != "" {
		switch req.GetType() {
		case 1: // 水池id
			id := easygo.StringToInt64noErr(req.GetKeyword())
			queryBson["_id"] = id
		case 2: // 水池名称
			queryBson["Name"] = req.GetKeyword()
		default:

		}
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL)
	defer closeFun()

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.WishPool
	var errc error

	if req.GetPageSize() == 0 {
		errc = query.Sort(sort...).Select(bson.M{"_id": 1}).All(&list)
	} else {
		errc = query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).Select(bson.M{"_id": 1}).All(&list)
	}
	easygo.PanicError(errc)

	var ids []int64

	for _, v := range list {
		ids = append(ids, v.GetId())
	}

	return ids, int32(count)
}

// 获取子水池列表
func GetSubWishPoolIds(id int64) []int64 {
	queryBson := bson.M{"PoolConfigId": id}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL)
	defer closeFun()

	query := col.Find(queryBson)

	var list []*share_message.WishPool
	var errc error

	errc = query.Select(bson.M{"_id": 1}).All(&list)
	easygo.PanicError(errc)

	var ids []int64

	for _, v := range list {
		ids = append(ids, v.GetId())
	}

	return ids
}

// 新增/更新水池
func UpdateWishPool(data *brower_backstage.WishPool) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL_CFG)
	defer closeFun()

	upData := &share_message.WishPoolCfg{
		Id:               easygo.NewInt64(data.GetId()),
		PoolLimit:        easygo.NewInt64(data.GetPoolLimit()),
		Name:             easygo.NewString(data.GetName()),
		ShowInitialValue: easygo.NewInt64(data.GetShowInitialValue()),
		ShowRecycle:      easygo.NewInt64(data.GetShowRecycle()),
		ShowCommission:   easygo.NewInt64(data.GetShowCommission()),
		ShowStartAward:   easygo.NewInt64(data.GetShowStartAward()),
		ShowCloseAward:   easygo.NewInt64(data.GetShowCloseAward()),
		IsDefault:        easygo.NewBool(data.GetIsDefault()),
	}
	upData.InitialValue = easygo.NewInt64(upData.GetShowInitialValue() * upData.GetPoolLimit() / divisor)
	upData.Recycle = easygo.NewInt64(upData.GetShowRecycle() * upData.GetPoolLimit() / divisor)
	upData.Commission = easygo.NewInt64(upData.GetShowCommission() * upData.GetPoolLimit() / divisor)
	upData.StartAward = easygo.NewInt64(upData.GetShowStartAward() * upData.GetPoolLimit() / divisor)
	upData.CloseAward = easygo.NewInt64(upData.GetShowCloseAward() * upData.GetPoolLimit() / divisor)
	logs.Info("pool:%+v", upData.GetRecycle())
	logs.Info("ShowRecycle", upData.GetShowRecycle()*upData.GetPoolLimit()/divisor)
	logs.Info("ShowRecycle", upData.GetShowRecycle(), upData.GetPoolLimit(), divisor)
	logs.Info("Commission", upData.GetCommission()*upData.GetPoolLimit()/divisor)
	logs.Info("Commission", upData.GetCommission(), upData.GetPoolLimit(), divisor)
	if upData.GetInitialValue() != 0 {
		upData.IncomeValue = easygo.NewInt64(upData.GetInitialValue())
	}

	upData.SmallLoss = &share_message.WishPoolStatus{
		ShowMaxValue: easygo.NewInt64(data.SmallLoss.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt64(data.SmallLoss.GetShowMinValue()),
	}
	upData.SmallLoss.MaxValue = easygo.NewInt64(int64(upData.SmallLoss.GetShowMaxValue()) * upData.GetPoolLimit() / divisor)
	upData.SmallLoss.MinValue = easygo.NewInt64(int64(upData.SmallLoss.GetShowMinValue()) * upData.GetPoolLimit() / divisor)

	upData.SmallWin = &share_message.WishPoolStatus{
		ShowMaxValue: easygo.NewInt64(data.SmallWin.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt64(data.SmallWin.GetShowMinValue()),
	}
	upData.SmallWin.MaxValue = easygo.NewInt64(int64(upData.SmallWin.GetShowMaxValue()) * upData.GetPoolLimit() / divisor)
	upData.SmallWin.MinValue = easygo.NewInt64(int64(upData.SmallWin.GetShowMinValue()) * upData.GetPoolLimit() / divisor)

	upData.BigLoss = &share_message.WishPoolStatus{
		ShowMaxValue: easygo.NewInt64(data.BigLoss.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt64(data.BigLoss.GetShowMinValue()),
	}
	upData.BigLoss.MaxValue = easygo.NewInt64(int64(upData.BigLoss.GetShowMaxValue()) * upData.GetPoolLimit() / divisor)
	upData.BigLoss.MinValue = easygo.NewInt64(int64(upData.BigLoss.GetShowMinValue()) * upData.GetPoolLimit() / divisor)

	upData.BigWin = &share_message.WishPoolStatus{
		ShowMaxValue: easygo.NewInt64(data.BigWin.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt64(data.BigWin.GetShowMinValue()),
	}
	upData.BigWin.MaxValue = easygo.NewInt64(int64(upData.BigWin.GetShowMaxValue()) * upData.GetPoolLimit() / divisor)
	upData.BigWin.MinValue = easygo.NewInt64(int64(upData.BigWin.GetShowMinValue()) * upData.GetPoolLimit() / divisor)

	upData.Common = &share_message.WishPoolStatus{
		ShowMaxValue: easygo.NewInt64(data.Common.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt64(data.Common.GetShowMinValue()),
	}

	upData.Common.MaxValue = easygo.NewInt64(int64(upData.Common.GetShowMaxValue()) * upData.GetPoolLimit() / divisor)
	upData.Common.MinValue = easygo.NewInt64(int64(upData.Common.GetShowMinValue()) * upData.GetPoolLimit() / divisor)

	curTime := time.Now().Unix()
	upData.UpdateTime = easygo.NewInt64(curTime)
	// 新增
	if upData.GetId() == 0 {
		upData.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_POOL_CFG))
		upData.CreateTime = easygo.NewInt64(curTime)
		//err := col.Insert(upData)
		//easygo.PanicError(err)

		//return
	}

	logs.Info("pool:%+v", upData)

	_, err1 := col.Upsert(bson.M{"_id": upData.GetId()}, bson.M{"$set": upData})
	easygo.PanicError(err1)
}

// 新增/更新水池,特别处理了水池盈利状态
func UpdateWishPool2(data *brower_backstage.WishPool) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL)
	defer closeFun()

	upData := &share_message.WishPool{
		Id:               easygo.NewInt64(data.GetId()),
		PoolLimit:        easygo.NewInt64(data.GetPoolLimit()),
		Name:             easygo.NewString(data.GetName()),
		ShowInitialValue: easygo.NewInt64(data.GetShowInitialValue()),
		ShowRecycle:      easygo.NewInt64(data.GetShowRecycle()),
		ShowCommission:   easygo.NewInt64(data.GetShowCommission()),
		ShowStartAward:   easygo.NewInt64(data.GetShowStartAward()),
		ShowCloseAward:   easygo.NewInt64(data.GetShowCloseAward()),
		IsDefault:        easygo.NewBool(data.GetIsDefault()),
		LocalStatus:      easygo.NewInt64(data.GetLocalStatus()),
		PoolConfigId:     easygo.NewInt64(data.GetPoolCfgId()),
		IsOpenAward:      easygo.NewBool(data.GetIsOpenAward()),
	}
	upData.InitialValue = easygo.NewInt64(upData.GetShowInitialValue() * upData.GetPoolLimit() / divisor)
	upData.Recycle = easygo.NewInt64(upData.GetShowRecycle() * upData.GetPoolLimit() / divisor)
	upData.Commission = easygo.NewInt64(upData.GetShowCommission() * upData.GetPoolLimit() / divisor)
	upData.StartAward = easygo.NewInt64(upData.GetShowStartAward() * upData.GetPoolLimit() / divisor)
	upData.CloseAward = easygo.NewInt64(upData.GetShowCloseAward() * upData.GetPoolLimit() / divisor)
	// logs.Info("pool:%+v", upData.GetRecycle())
	// logs.Info("ShowRecycle", upData.GetShowRecycle()*upData.GetPoolLimit()/divisor)
	// logs.Info("ShowRecycle", upData.GetShowRecycle(), upData.GetPoolLimit(), divisor)
	// logs.Info("Commission", upData.GetCommission()*upData.GetPoolLimit()/divisor)
	// logs.Info("Commission", upData.GetCommission(), upData.GetPoolLimit(), divisor)
	if upData.GetInitialValue() != 0 {
		upData.IncomeValue = easygo.NewInt64(upData.GetInitialValue())
	}

	upData.SmallLoss = &share_message.WishPoolStatus{
		ShowMaxValue: easygo.NewInt64(data.SmallLoss.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt64(data.SmallLoss.GetShowMinValue()),
	}
	upData.SmallLoss.MaxValue = easygo.NewInt64(int64(upData.SmallLoss.GetShowMaxValue()) * upData.GetPoolLimit() / divisor)
	upData.SmallLoss.MinValue = easygo.NewInt64(int64(upData.SmallLoss.GetShowMinValue()) * upData.GetPoolLimit() / divisor)

	upData.SmallWin = &share_message.WishPoolStatus{
		ShowMaxValue: easygo.NewInt64(data.SmallWin.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt64(data.SmallWin.GetShowMinValue()),
	}
	upData.SmallWin.MaxValue = easygo.NewInt64(int64(upData.SmallWin.GetShowMaxValue()) * upData.GetPoolLimit() / divisor)
	upData.SmallWin.MinValue = easygo.NewInt64(int64(upData.SmallWin.GetShowMinValue()) * upData.GetPoolLimit() / divisor)

	upData.BigLoss = &share_message.WishPoolStatus{
		ShowMaxValue: easygo.NewInt64(data.BigLoss.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt64(data.BigLoss.GetShowMinValue()),
	}
	upData.BigLoss.MaxValue = easygo.NewInt64(int64(upData.BigLoss.GetShowMaxValue()) * upData.GetPoolLimit() / divisor)
	upData.BigLoss.MinValue = easygo.NewInt64(int64(upData.BigLoss.GetShowMinValue()) * upData.GetPoolLimit() / divisor)

	upData.BigWin = &share_message.WishPoolStatus{
		ShowMaxValue: easygo.NewInt64(data.BigWin.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt64(data.BigWin.GetShowMinValue()),
	}
	upData.BigWin.MaxValue = easygo.NewInt64(int64(upData.BigWin.GetShowMaxValue()) * upData.GetPoolLimit() / divisor)
	upData.BigWin.MinValue = easygo.NewInt64(int64(upData.BigWin.GetShowMinValue()) * upData.GetPoolLimit() / divisor)

	upData.Common = &share_message.WishPoolStatus{
		ShowMaxValue: easygo.NewInt64(data.Common.GetShowMaxValue()),
		ShowMinValue: easygo.NewInt64(data.Common.GetShowMinValue()),
	}

	upData.Common.MaxValue = easygo.NewInt64(int64(upData.Common.GetShowMaxValue()) * upData.GetPoolLimit() / divisor)
	upData.Common.MinValue = easygo.NewInt64(int64(upData.Common.GetShowMinValue()) * upData.GetPoolLimit() / divisor)

	curTime := time.Now().Unix()
	upData.UpdateTime = easygo.NewInt64(curTime)

	status := GetPoolStatus(upData)
	upData.LocalStatus = easygo.NewInt64(status)

	// 新增
	if upData.GetId() == 0 {
		upData.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_POOL))
		upData.CreateTime = easygo.NewInt64(curTime)
		//err := col.Insert(upData)
		//easygo.PanicError(err)
		//
		//return
	}

	logs.Info("pool:%+v", upData)

	_, err1 := col.Upsert(bson.M{"_id": upData.GetId()}, bson.M{"$set": upData})
	easygo.PanicError(err1)
}

// 根据盲盒物品id获取愿望
func GetPlayerWishDataByWishBoxItemId(pid, itemId int64) *share_message.PlayerWishData {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA)
	defer closeFun()

	list := &share_message.PlayerWishData{}
	queryBson := bson.M{"PlayerId": pid, "WishBoxItemId": itemId}
	err := col.Find(queryBson).One(&list)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}

	return list
}

// 根据盲盒物品id获取愿望
func GetPlayerWishDataByWishBoxItemIds(pid int64, itemId []int64) []*share_message.PlayerWishData {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA)
	defer closeFun()

	var list []*share_message.PlayerWishData
	queryBson := bson.M{"PlayerId": pid, "WishBoxItemId": bson.M{"$in": itemId}}
	err := col.Find(queryBson).All(&list)
	easygo.PanicError(err)

	return list
}

// 删除水池
func DeleteWishPool(ids []int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"_id": bson.M{"$in": ids}})
	easygo.PanicError(err)
}

// 获取价格区间参数
func GetPriceSection() *brower_backstage.PriceSection {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	list := &share_message.PriceSection{}
	queryBson := bson.M{"_id": for_game.TABLE_WISH_PRICE_SECTION}
	err := col.Find(queryBson).One(&list)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	res := &brower_backstage.PriceSection{
		OneMin:   easygo.NewInt32(list.GetOneMin()),
		OneMax:   easygo.NewInt32(list.GetOneMax()),
		TwoMin:   easygo.NewInt32(list.GetTwoMin()),
		TwoMax:   easygo.NewInt32(list.GetTwoMax()),
		ThreeMin: easygo.NewInt32(list.GetThreeMin()),
	}
	if err == mgo.ErrNotFound {
		return res
	}

	return res
}

// 更新价格区间参数
func UpdatePriceSection(list *brower_backstage.PriceSection) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	data := &share_message.PriceSection{
		Id:       easygo.NewString(for_game.TABLE_WISH_PRICE_SECTION),
		OneMin:   easygo.NewInt32(list.GetOneMin()),
		OneMax:   easygo.NewInt32(list.GetOneMax()),
		TwoMin:   easygo.NewInt32(list.GetTwoMin()),
		TwoMax:   easygo.NewInt32(list.GetTwoMax()),
		ThreeMin: easygo.NewInt32(list.GetThreeMin()),
	}

	_, err := col.Upsert(bson.M{"_id": data.GetId()}, bson.M{"$set": data})
	easygo.PanicError(err)
}

// 获取邮寄参数设置
func GetMailSection() *brower_backstage.WishMailSection {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	list := &share_message.WishMailSection{}
	queryBson := bson.M{"_id": for_game.TABLE_WISH_MAIL_SECTION}
	err := col.Find(queryBson).One(&list)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	res := &brower_backstage.WishMailSection{
		Postage1:       easygo.NewInt32(list.GetPostage1()),
		Postage2:       easygo.NewInt32(list.GetPostage2()),
		Postage3:       easygo.NewInt32(list.GetPostage3()),
		RemoteAreaList: list.GetRemoteAreaList(),
		FreeNumber:     easygo.NewInt32(list.GetFreeNumber()),
	}
	if err == mgo.ErrNotFound {
		return res
	}

	return res
}

// 更新邮寄参数设置
func UpdateMailSection(list *brower_backstage.WishMailSection) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	data := &share_message.WishMailSection{
		Id:             easygo.NewString(for_game.TABLE_WISH_MAIL_SECTION),
		Postage1:       easygo.NewInt32(list.GetPostage1()),
		Postage2:       easygo.NewInt32(list.GetPostage2()),
		Postage3:       easygo.NewInt32(list.GetPostage3()),
		RemoteAreaList: list.GetRemoteAreaList(),
		FreeNumber:     easygo.NewInt32(list.GetFreeNumber()),
	}

	_, err := col.Upsert(bson.M{"_id": data.GetId()}, bson.M{"$set": data})
	easygo.PanicError(err)
}

//更新回收说明
func SaveRecycleNoteCfg(data *share_message.RecycleNoteCfg) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": data.GetId()}, bson.M{"$set": data})
	easygo.PanicError(err)
}

//查看回收说明
func GetRecycleNoteCfg() *share_message.RecycleNoteCfg {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()
	data := &share_message.RecycleNoteCfg{}
	err := col.Find(bson.M{"_id": for_game.TABLE_WISH_RECYCLE_NOTE_CFG}).One(&data)
	if err != nil {
		return nil
	}
	/*if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}*/
	return data
}

// 获取物品回收参数设置
func GetWishRecycleSection() *brower_backstage.WishRecycleSection {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	list := &share_message.WishRecycleSection{}
	queryBson := bson.M{"_id": for_game.TABLE_WISH_RECYCLE_SECTION}
	err := col.Find(queryBson).One(&list)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	res := &brower_backstage.WishRecycleSection{
		Platform:            easygo.NewInt32(list.GetPlatform()),
		DayTopCount:         easygo.NewInt32(list.GetDayTopCount()),
		Player:              easygo.NewInt32(list.GetPlayer()),
		Status:              easygo.NewBool(list.GetStatus()),
		DayMoneyTopCount:    easygo.NewInt32(list.GetDayMoneyTopCount()),
		DayMoneyTop:         easygo.NewInt64(list.GetDayMoneyTop()),
		OrderThresholdMoney: easygo.NewInt64(list.GetOrderThresholdMoney()),
		DayDiamondTopCount:  easygo.NewInt32(list.GetDayDiamondTopCount()),
		DayDiamondTop:       easygo.NewInt64(list.GetDayDiamondTop()),
		OrderThreshold:      easygo.NewInt64(list.GetOrderThreshold()),

		FeeRate:     easygo.NewInt32(list.GetFeeRate()),
		PlatformTax: easygo.NewInt64(list.GetPlatformTax()),
		RealTax:     easygo.NewInt64(list.GetRealTax()),
	}
	if err == mgo.ErrNotFound {
		return res
	}

	return res
}

// 更新物品回收参数设置
func UpdateWishRecycleSection(list *brower_backstage.WishRecycleSection) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	data := &share_message.WishRecycleSection{
		Id:                  easygo.NewString(for_game.TABLE_WISH_RECYCLE_SECTION),
		Platform:            easygo.NewInt32(list.GetPlatform()),
		DayTopCount:         easygo.NewInt32(list.GetDayTopCount()),
		Player:              easygo.NewInt32(list.GetPlayer()),
		Status:              easygo.NewBool(list.GetStatus()),
		DayMoneyTopCount:    easygo.NewInt32(list.GetDayMoneyTopCount()),
		DayMoneyTop:         easygo.NewInt64(list.GetDayMoneyTop()),
		OrderThresholdMoney: easygo.NewInt64(list.GetOrderThresholdMoney()),
		DayDiamondTopCount:  easygo.NewInt32(list.GetDayDiamondTopCount()),
		DayDiamondTop:       easygo.NewInt64(list.GetDayDiamondTop()),
		OrderThreshold:      easygo.NewInt64(list.GetOrderThreshold()),

		FeeRate:     easygo.NewInt32(list.GetFeeRate()),
		PlatformTax: easygo.NewInt64(list.GetPlatformTax()),
		RealTax:     easygo.NewInt64(list.GetRealTax()),
	}

	_, err := col.Upsert(bson.M{"_id": data.GetId()}, bson.M{"$set": data})
	easygo.PanicError(err)
}

// 获取支付预警
func GetWishPayWarnCfg() *brower_backstage.WishPayWarnCfg {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	list := &share_message.WishPayWarnCfg{}
	queryBson := bson.M{"_id": for_game.TABLE_WISH_PAY_WARN_CFG}
	err := col.Find(queryBson).One(&list)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
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
	if err == mgo.ErrNotFound {
		return res
	}

	return res
}

// 更新支付预警
func UpdateWishPayWarnCfg(list *brower_backstage.WishPayWarnCfg) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	data := &share_message.WishPayWarnCfg{
		WithdrawalTime:        easygo.NewInt64(list.GetWithdrawalTime()),
		WithdrawalTimes:       easygo.NewInt64(list.GetWithdrawalTimes()),
		WithdrawalGoldRate:    easygo.NewInt64(list.GetWithdrawalGoldRate()),
		WithdrawalGold:        easygo.NewInt64(list.GetWithdrawalGold()),
		WithdrawalDiamondRate: easygo.NewInt64(list.GetWithdrawalDiamondRate()),
		WithdrawalDiamond:     easygo.NewInt64(list.GetWithdrawalDiamond()),
		PhoneList:             list.GetPhoneList(),
	}

	_, err := col.Upsert(bson.M{"_id": for_game.TABLE_WISH_PAY_WARN_CFG}, bson.M{"$set": data})
	easygo.PanicError(err)
}

// 获取冷却配置表
func GetWishCoolDownConfigFromDB() *share_message.WishCoolDownConfig {
	col, closeFunc := easygo.MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFunc()
	result := &share_message.WishCoolDownConfig{}
	err := col.Find(bson.M{"_id": for_game.TABLE_WISH_COOL_DOWN_CONFIG}).One(&result)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return result
}

// 更新冷却配置表
func UpdateWishCoolDownConfigFromDB(data *brower_backstage.WishCoolDownConfig) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	upData := &share_message.WishCoolDownConfig{
		Id:              easygo.NewString(for_game.TABLE_WISH_COOL_DOWN_CONFIG),
		IsOpen:          easygo.NewBool(data.GetIsOpen()),
		ContinuousTime:  easygo.NewInt64(data.GetContinuousTime()),
		ContinuousTimes: easygo.NewInt64(data.GetContinuousTimes()),
		CoolDownTime:    easygo.NewInt64(data.GetCoolDownTime()),
		DayLimit:        easygo.NewInt64(data.GetDayLimit()),
		CreateTime:      easygo.NewInt64(data.GetCreateTime()),
	}

	_, err := col.Upsert(bson.M{"_id": upData.GetId()}, bson.M{"$set": upData})
	easygo.PanicError(err)
}

// 获取物品回收参数设置
func GetWishCurrencyConversionCfg() *brower_backstage.WishCurrencyConversionCfg {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	list := &share_message.WishCurrencyConversionCfg{}
	queryBson := bson.M{"_id": for_game.TABLE_WISH_CURRENCY_CONVERSION_CFG}
	err := col.Find(queryBson).One(&list)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	res := &brower_backstage.WishCurrencyConversionCfg{
		Diamond: easygo.NewInt32(list.GetDiamond()),
		Coin:    easygo.NewInt32(list.GetCoin()),
	}
	if err == mgo.ErrNotFound {
		return res
	}

	return res
}

// 更新物品回收参数设置
func UpdateWishCurrencyConversionCfg(list *brower_backstage.WishCurrencyConversionCfg) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	data := &share_message.WishCurrencyConversionCfg{
		Id:      easygo.NewString(for_game.TABLE_WISH_CURRENCY_CONVERSION_CFG),
		Diamond: easygo.NewInt32(list.GetDiamond()),
		Coin:    easygo.NewInt32(list.GetCoin()),
	}

	_, err := col.Upsert(bson.M{"_id": data.GetId()}, bson.M{"$set": data})
	easygo.PanicError(err)
}

// 获取守护者收益设置
func GetWishGuardianCfg() *brower_backstage.WishGuardianCfg {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	list := &share_message.WishGuardianCfg{}
	queryBson := bson.M{"_id": for_game.TABLE_WISH_GUARDIAN_CFG}
	err := col.Find(queryBson).One(&list)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	res := &brower_backstage.WishGuardianCfg{
		DayDiamondTop:     easygo.NewInt64(list.GetDayDiamondTop()),
		OnceDiamondRebate: easygo.NewInt64(list.GetOnceDiamondRebate()),
	}
	if err == mgo.ErrNotFound {
		return res
	}

	return res
}

// 更新守护者收益设置
func UpdateWishGuardianCfg(list *brower_backstage.WishGuardianCfg) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OTHER_CFG)
	defer closeFun()

	data := &share_message.WishGuardianCfg{
		Id:                easygo.NewString(for_game.TABLE_WISH_GUARDIAN_CFG),
		DayDiamondTop:     easygo.NewInt64(list.GetDayDiamondTop()),
		OnceDiamondRebate: easygo.NewInt64(list.GetOnceDiamondRebate()),
	}

	_, err := col.Upsert(bson.M{"_id": data.GetId()}, bson.M{"$set": data})
	easygo.PanicError(err)
}

// =======================================================================================================许愿池报表

// 生成许愿池相关报表
func MakeWishReports() {
	// easygo.Spawn(MakeWishItemReport) 暂无需求
	easygo.Spawn(MakeWishBoxReports)
	easygo.Spawn(MakeWishActReports)
}

func MakeWishBoxReports() {
	MakeWishBoxDetailReport()
	MakeWishBoxReport()
	MakeWishPoolReport()
}

//许愿池统计报表-周报
func MakeWishBoxReportsOfWeek() {
	MakeWishBoxDetailReportByWeek()
	MakeWishBoxReportByWeek()
	MakeWishPoolReportByWeek()
}

//许愿池统计报表-月报
func MakeWishBoxReportsOfMonth() {
	MakeWishBoxDetailReportByMonth()
	MakeWishBoxReportByMonth()
	MakeWishPoolReportByMonth()
}

//初始化许愿池商品报表
func InitWishItemReport() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_ITEM)
	defer closeFun()

	queryBson := bson.M{}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//更新许愿池商品报表
func UpdateWishItemReport(report *share_message.WishItemReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_ITEM)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetId()}, bson.M{"$set": report})
	easygo.PanicError(err)

}

//查询许愿池商品报表报表
func GetWishItemReport(id, startTime int64) *share_message.WishItemReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_ITEM)
	defer closeFun()
	report := &share_message.WishItemReport{}
	err := col.Find(bson.M{"CreateTime": startTime, "ItemId": id}).One(&report)

	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return report
}

// 聚合查询许愿池商品报表邀请报表
//func PipeWishItemReport(itemId, startTime int64) *share_message.WishItemReport {
//
//
//	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_ITEM)
//	defer closeFun()
//
//	lookupData := bson.M{"from": for_game.TABLE_PLAYER_WISH_DATA, "localField": "_id", "foreignField": "WishBoxItemId", "as": "data"}
//	lookupWish := bson.M{"from": for_game.TABLE_WISH_LOG, "localField": "data.WishBoxItemId", "foreignField": "ChallengeItemId", "as": "wish"}
//	lookupItem := bson.M{"from": for_game.TABLE_PLAYER_WISH_ITEM, "localField": "_id", "foreignField": "ChallengeItemId", "as": "item"}
//	lookupOrder := bson.M{"from": for_game.TABLE_WISH_RECYCLE_ORDER, "localField": "_id", "foreignField": "WishItemId", "as": "order"}
//	match := bson.M{"_id": itemId}
//	match["data.CreateTime"] = bson.M{"$gte": startTime}
//	match["wish.CreateTime"] = bson.M{"$gte": startTime}
//	match["order.RecycleTime"] = bson.M{"$gte": startTime}
//	match["item.CreateTime"] = bson.M{"$gte": startTime}
//	cond1 := []interface{}{bson.M{"$eq": []interface{}{"$item.Status", 1}}, 1, 0}
//	cond2 := []interface{}{bson.M{"$eq": []interface{}{"$item.Status", 0}}, 1, 0}
//	group1 := bson.M{
//		"_id":            "$data.WishBoxItemId",
//		"WishPlayerCount": bson.M{"$sum": "$data._id"},
//		"ConvertCount": bson.M{"$sum": bson.M{"$cond": cond1}},
//		"PendConvertCount": bson.M{"$sum": bson.M{"$cond": cond2}},
//	}
//
//	group2 := bson.M{
//		"_id":           "$wish",
//		"WishAndDrawPlayerCount": bson.M{"$sum": 1},
//	}
//
//
//
//	group3 := bson.M{
//		"_id":           "$item.ChallengeItemId",
//		"WinCount": bson.M{"$sum": 1},
//		"ConvertCount": bson.M{"$sum": bson.M{"$cond": cond1}},
//		"PendConvertCount": bson.M{"$sum": bson.M{"$cond": cond2}},
//	}
//
//	cond3 := []interface{}{bson.M{"$eq": []interface{}{"$order.Type", 1}}, 1, 0}
//	cond4 := []interface{}{bson.M{"$eq": []interface{}{"$order.Type", 2}}, 1, 0}
//	group4 := bson.M{
//		"_id":           "$order.WishItemId",
//		"RecycleCount": bson.M{"$sum": 1},
//		"PlayerRecycleCount": bson.M{"$sum": bson.M{"$cond": cond3}},
//		"OfficialRecycleCount": bson.M{"$sum": bson.M{"$cond": cond4}},
//	}
//
//	countType := "PageCount"
//
//	countCond := []bson.M{
//		{"$match": match},
//		{"$count": countType},
//	}
//
//	var retList *brower_backstage.DrawBoxRecordList
//	err := col.Pipe(countCond).One(&retList)
//
//	project := bson.M{
//		//"ItemId":      "$_id",
//		//"ItemName":    "$Name",
//		"_id": "$_id",
//		"WishPlayerCount": 1,
//		"WishAndDrawPlayerCount": 1,
//		"WinCount": 1,
//		"ConvertCount": 1,
//		"PendConvertCount": 1,
//		"RecycleCount": 1,
//		"PlayerRecycleCount": 1,
//		"OfficialRecycleCount": 1,
//		"data": "$data",
//		"item": "$item",
//	}
//
//	logs.Info(group1)
//	logs.Info(group2)
//	logs.Info(group3)
//	logs.Info(group4)
//	logs.Info(lookupData)
//	logs.Info(lookupWish)
//	logs.Info(lookupItem)
//	logs.Info(lookupOrder)
//	pipeCond := []bson.M{
//
//
//
//		//{"$lookup": lookupWish},
//		//{"$unwind": "$wish"},
//		//{"$match": match},
//		//{"$group": group2},
//		////{"$group": group2},
//		{"$lookup": lookupItem},
//		{"$unwind": "$item"},
//		{"$group": group1},
//		{"$project": project},
//		//{"$lookup": lookupOrder},
//		//{"$unwind": "$order"},
//
//		//{"$match": match},
//		//{"$group": group3},
//		//{"$group": group4},
//		{"$lookup": lookupData},
//		{"$unwind": "$data"},
//
//		{"$project": project},
//	}
//
//	var one *share_message.WishItemReport
//	var list []*share_message.WishItemReport
//	err = col.Pipe(pipeCond).All(&list)
//	easygo.PanicError(err)
//	retList.List = nil
//	return one
//}

// 查询商品回收统计
func GroupWishRecycleOrder(tid, startTime int64) *share_message.WishItemReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()

	match := bson.M{"RecycleItemList.ProductId": tid}
	match["RecycleTime"] = bson.M{"$gt": startTime}
	cond3 := []interface{}{bson.M{"$eq": []interface{}{"$Type", 1}}, 1, 0}
	cond4 := []interface{}{bson.M{"$eq": []interface{}{"$Type", 2}}, 1, 0}
	group := bson.M{
		"_id":                  nil,
		"RecycleCount":         bson.M{"$sum": 1},
		"PlayerRecycleCount":   bson.M{"$sum": bson.M{"$cond": cond3}},
		"OfficialRecycleCount": bson.M{"$sum": bson.M{"$cond": cond4}},
	}

	project := bson.M{
		"RecycleCount":         1,
		"PlayerRecycleCount":   1,
		"OfficialRecycleCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishItemReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询商品回收统计
func GroupWishRecycleOrderByTime(startTime int64) *share_message.WishItemReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()

	match := bson.M{}
	match["RecycleTime"] = bson.M{"$gt": startTime}
	cond3 := []interface{}{bson.M{"$eq": []interface{}{"$Type", 1}}, 1, 0}
	cond4 := []interface{}{bson.M{"$eq": []interface{}{"$Type", 2}}, 1, 0}
	group := bson.M{
		"_id":                  nil,
		"RecycleCount":         bson.M{"$sum": 1},
		"PlayerRecycleCount":   bson.M{"$sum": bson.M{"$cond": cond3}},
		"OfficialRecycleCount": bson.M{"$sum": bson.M{"$cond": cond4}},
	}

	project := bson.M{
		"RecycleCount":         1,
		"PlayerRecycleCount":   1,
		"OfficialRecycleCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishItemReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询商品回收统计
func GroupWishRecycleOrderByBoxItemIdByTime(tid, startTime, entTime int64) *share_message.WishBoxDetailReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()

	match := bson.M{"RecycleItemList.WishBoxItemId": tid}
	match["RecycleTime"] = bson.M{"$gte": startTime, "$lt": entTime}
	cond3 := []interface{}{bson.M{"$eq": []interface{}{"$Type", 1}}, 1, 0}
	cond4 := []interface{}{bson.M{"$eq": []interface{}{"$Type", 2}}, 1, 0}
	cond5 := []interface{}{bson.M{"$eq": []interface{}{"$Type", 1}}, "$RecycleDiamond", 0}
	cond6 := []interface{}{bson.M{"$eq": []interface{}{"$Type", 2}}, "$RecycleDiamond", 0}
	cond7 := []interface{}{bson.M{"$eq": []interface{}{"$RecycleItemList.WishBoxItemId", tid}}, "$RecycleItemList.RecycleDiamond", 0}
	group := bson.M{
		"_id":                      nil,
		"RecycleCount":             bson.M{"$sum": 1},
		"PlayerRecycleCount":       bson.M{"$sum": bson.M{"$cond": cond3}},
		"OfficialRecycleCount":     bson.M{"$sum": bson.M{"$cond": cond4}},
		"RecycleGoldTotal":         bson.M{"$sum": bson.M{"$cond": cond7}},
		"PlayerRecycleGoldTotal":   bson.M{"$sum": bson.M{"$cond": cond5}},
		"OfficialRecycleGoldTotal": bson.M{"$sum": bson.M{"$cond": cond6}},
	}
	//		"RecyclePrice":
	project := bson.M{
		"RecycleCount":             1,
		"PlayerRecycleCount":       1,
		"OfficialRecycleCount":     1,
		"RecycleGoldTotal":         1,
		"PlayerRecycleGoldTotal":   1,
		"OfficialRecycleGoldTotal": 1,
	}

	match2 := bson.M{"RecycleItemList.WishBoxItemId": tid}
	pipeCond := []bson.M{
		{"$match": match},
		{"$unwind": "$RecycleItemList"},
		{"$match": match2},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishBoxDetailReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询商品中奖统计
func GroupPlayerWishItemByTime(tid, startTime, entTime int64) *share_message.WishItemReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	match := bson.M{"WishItemId": tid}
	match["CreateTime"] = bson.M{"$gte": startTime, "$lt": entTime}
	cond3 := []interface{}{bson.M{"$eq": []interface{}{"$Status", 0}}, 1, 0}
	cond4 := []interface{}{bson.M{"$eq": []interface{}{"$Status", 1}}, 1, 0}
	group := bson.M{
		"_id":              nil,
		"WinCount":         bson.M{"$sum": 1},
		"ConvertCount":     bson.M{"$sum": bson.M{"$cond": cond4}},
		"PendConvertCount": bson.M{"$sum": bson.M{"$cond": cond3}},
	}

	project := bson.M{
		"WinCount":         1,
		"ConvertCount":     1,
		"PendConvertCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishItemReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询商品中奖统计
func GroupPlayerWishItemByBoxItemIdByTime(tid, startTime, entTime int64) *share_message.WishBoxDetailReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	match := bson.M{"ChallengeItemId": tid}
	match["CreateTime"] = bson.M{"$gte": startTime, "$lt": entTime}
	cond3 := []interface{}{bson.M{"$eq": []interface{}{"$Status", 0}}, 1, 0}
	cond4 := []interface{}{bson.M{"$eq": []interface{}{"$Status", 1}}, 1, 0}
	cond5 := []interface{}{bson.M{"$eq": []interface{}{"$Status", 0}}, "$WishItemDiamond", 0}
	cond6 := []interface{}{bson.M{"$eq": []interface{}{"$Status", 1}}, "$WishItemDiamond", 0}
	group := bson.M{
		"_id":                    nil,
		"WinCount":               bson.M{"$sum": 1},
		"ConvertCount":           bson.M{"$sum": bson.M{"$cond": cond4}},
		"PendConvertCount":       bson.M{"$sum": bson.M{"$cond": cond3}},
		"ConvertGoodsPriceTotal": bson.M{"$sum": bson.M{"$cond": cond6}},
		"PendConvertPriceTotal":  bson.M{"$sum": bson.M{"$cond": cond5}},
	}

	project := bson.M{
		"WinCount":               1,
		"ConvertCount":           1,
		"PendConvertCount":       1,
		"ConvertGoodsPriceTotal": 1,
		"PendConvertPriceTotal":  1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishBoxDetailReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询商品许愿统计
func GroupPlayerWishDataByTime(tid, startTime, entTime int64) *share_message.WishItemReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA)
	defer closeFun()

	match := bson.M{"WishItemId": tid}
	//lookup := bson.M{"from": for_game.TABLE_PLAYER_WISH_DATA, "localField": "_id", "foreignField": "WishBoxItemId", "as": "item"}
	match["CreateTime"] = bson.M{"$gte": startTime, "$lt": entTime}
	group := bson.M{
		"_id":             nil,
		"WishPlayerCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"WishPlayerCount": 1,
	}

	pipeCond := []bson.M{

		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishItemReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询商品许愿并抽奖统计
func GroupPlayerWishAndDrawDataByBoxItemIdByTime(tid, startTime, entTime int64) *share_message.WishBoxDetailReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG)
	defer closeFun()

	match := bson.M{"ChallengeItemId": tid}
	lookupData := bson.M{"from": for_game.TABLE_PLAYER_WISH_DATA, "localField": "ChallengeItemId", "foreignField": "WishBoxItemId", "as": "data"}
	match["CreateTime"] = bson.M{"$gte": startTime, "$lt": entTime}
	match["data.CreateTime"] = bson.M{"$gte": startTime, "$lt": entTime}
	match["data.Status"] = bson.M{"$ne": 2} // 过滤取消许愿的
	match["$expr"] = bson.M{"$eq": []string{"$DareId", "$data.PlayerId"}}
	cond4 := []interface{}{bson.M{"$eq": []interface{}{"$Result", true}}, 1, 0}
	group := bson.M{
		//"_id":             bson.M{"ChallengeItemId": "$ChallengeItemId"},
		"_id":                    nil,
		"WishAndDrawPlayerCount": bson.M{"$sum": 1},
		"LuckyWishCount":         bson.M{"$sum": bson.M{"$cond": cond4}},
		"LuckyGoldTotal":         bson.M{"$sum": "$DarePrice"},
	}

	project := bson.M{
		"WishAndDrawPlayerCount": 1,
		"LuckyWishCount":         1,
		"LuckyGoldTotal":         1,
	}

	pipeCond := []bson.M{
		{"$lookup": lookupData},
		{"$unwind": "$data"},
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishBoxDetailReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询商品许愿统计
func GroupPlayerWishDataByBoxItemIdByTime(tid, startTime, entTime int64) *share_message.WishItemReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA)
	defer closeFun()

	match := bson.M{"WishBoxItemId": tid}
	match["CreateTime"] = bson.M{"$gte": startTime, "$lt": entTime}
	group := bson.M{
		"_id":             nil,
		"WishPlayerCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"WishPlayerCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishItemReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询商品回收统计
func GroupWishRecycleOrderByBoxItemId(tid, startTime, endTime int64) *share_message.WishBoxDetailReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER)
	defer closeFun()

	match := bson.M{"RecycleItemList.WishBoxItemId": tid}

	if endTime > 0 {
		match["RecycleTime"] = bson.M{"$gte": startTime, "$lt": endTime}
	} else {
		match["RecycleTime"] = bson.M{"$gt": startTime}
	}
	cond3 := []interface{}{bson.M{"$eq": []interface{}{"$Type", 1}}, 1, 0}
	cond4 := []interface{}{bson.M{"$eq": []interface{}{"$Type", 2}}, 1, 0}
	cond5 := []interface{}{bson.M{"$eq": []interface{}{"$Type", 1}}, "$RecycleItemList.RecycleDiamond", 0}
	cond6 := []interface{}{bson.M{"$eq": []interface{}{"$Type", 2}}, "$RecycleItemList.RecycleDiamond", 0}
	cond7 := []interface{}{bson.M{"$eq": []interface{}{"$RecycleItemList.WishBoxItemId", tid}}, "$RecycleItemList.RecycleDiamond", 0}
	cond8 := []interface{}{bson.M{"$eq": []interface{}{"$RecycleItemList.WishBoxItemId", tid}}, "$RecycleItemList.ProductDiamond", 0}
	group := bson.M{
		"_id":                      nil,
		"RecycleCount":             bson.M{"$sum": 1},
		"PlayerRecycleCount":       bson.M{"$sum": bson.M{"$cond": cond3}},
		"OfficialRecycleCount":     bson.M{"$sum": bson.M{"$cond": cond4}},
		"RecycleGoldTotal":         bson.M{"$sum": bson.M{"$cond": cond7}},
		"ProductDiamondTotal":      bson.M{"$sum": bson.M{"$cond": cond8}},
		"PlayerRecycleGoldTotal":   bson.M{"$sum": bson.M{"$cond": cond5}},
		"OfficialRecycleGoldTotal": bson.M{"$sum": bson.M{"$cond": cond6}},
	}
	//		"RecyclePrice":
	project := bson.M{
		"RecycleCount":             1,
		"PlayerRecycleCount":       1,
		"OfficialRecycleCount":     1,
		"RecycleGoldTotal":         1,
		"ProductDiamondTotal":      1,
		"PlayerRecycleGoldTotal":   1,
		"OfficialRecycleGoldTotal": 1,
	}

	match2 := bson.M{"RecycleItemList.WishBoxItemId": tid}
	pipeCond := []bson.M{
		{"$match": match},
		{"$unwind": "$RecycleItemList"},
		{"$match": match2},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishBoxDetailReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询商品中奖统计
func GroupPlayerWishItem(tid, startTime int64) *share_message.WishItemReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	match := bson.M{"WishItemId": tid}
	match["CreateTime"] = bson.M{"$gt": startTime}
	cond3 := []interface{}{bson.M{"$eq": []interface{}{"$Status", 0}}, 1, 0}
	cond4 := []interface{}{bson.M{"$eq": []interface{}{"$Status", 1}}, 1, 0}
	group := bson.M{
		"_id":              nil,
		"WinCount":         bson.M{"$sum": 1},
		"ConvertCount":     bson.M{"$sum": bson.M{"$cond": cond4}},
		"PendConvertCount": bson.M{"$sum": bson.M{"$cond": cond3}},
	}

	project := bson.M{
		"WinCount":         1,
		"ConvertCount":     1,
		"PendConvertCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishItemReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询商品中奖统计
func GroupPlayerWishItemByBoxItemId(tid, startTime, endTime int64) *share_message.WishBoxDetailReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	match := bson.M{"ChallengeItemId": tid}
	if endTime > 0 {
		match["CreateTime"] = bson.M{"$gt": startTime, "$lt": endTime}
	} else {
		match["CreateTime"] = bson.M{"$gt": startTime}
	}

	cond3 := []interface{}{bson.M{"$eq": []interface{}{"$Status", 0}}, 1, 0}
	cond4 := []interface{}{bson.M{"$eq": []interface{}{"$Status", 1}}, 1, 0}
	cond5 := []interface{}{bson.M{"$eq": []interface{}{"$Status", 0}}, "$WishItemDiamond", 0}
	cond6 := []interface{}{bson.M{"$eq": []interface{}{"$Status", 1}}, "$WishItemDiamond", 0}
	cond7 := []interface{}{bson.M{"$eq": []interface{}{"$Status", 3}}, 1, 0} // 审核中的也算在待回收中
	cond8 := []interface{}{bson.M{"$eq": []interface{}{"$Status", 3}}, "$WishItemDiamond", 0}
	group := bson.M{
		"_id":                    nil,
		"WinCount":               bson.M{"$sum": 1},
		"WinItemPriceTotal":      bson.M{"$sum": "$WishItemDiamond"},
		"ConvertCount":           bson.M{"$sum": bson.M{"$cond": cond4}},
		"PendConvertCount":       bson.M{"$sum": bson.M{"$cond": cond3}},
		"ConvertGoodsPriceTotal": bson.M{"$sum": bson.M{"$cond": cond6}},
		"PendConvertPriceTotal":  bson.M{"$sum": bson.M{"$cond": cond5}},

		"AuditConvertCount":      bson.M{"$sum": bson.M{"$cond": cond7}},
		"AuditConvertPriceTotal": bson.M{"$sum": bson.M{"$cond": cond8}},
	}

	project := bson.M{
		"WinCount":               1,
		"WinItemPriceTotal":      1,
		"ConvertCount":           1,
		"PendConvertCount":       1,
		"ConvertGoodsPriceTotal": 1,
		"PendConvertPriceTotal":  1,

		"AuditConvertCount":      1,
		"AuditConvertPriceTotal": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishBoxDetailReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询商品许愿统计
func GroupPlayerWishData(tid, startTime int64) *share_message.WishItemReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA)
	defer closeFun()

	match := bson.M{"WishItemId": tid}
	//lookup := bson.M{"from": for_game.TABLE_PLAYER_WISH_DATA, "localField": "_id", "foreignField": "WishBoxItemId", "as": "item"}
	match["CreateTime"] = bson.M{"$gt": startTime}
	group := bson.M{
		"_id":             nil,
		"WishPlayerCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"WishPlayerCount": 1,
	}

	pipeCond := []bson.M{

		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishItemReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询商品许愿统计
func GroupPlayerWishDataByBoxItemId(tid, startTime, endTime int64) *share_message.WishItemReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA)
	defer closeFun()

	match := bson.M{"WishBoxItemId": tid}
	if endTime > 0 {
		match["CreateTime"] = bson.M{"$gte": startTime, "$lt": endTime}
	} else {
		match["CreateTime"] = bson.M{"$gt": startTime}
	}

	group := bson.M{
		"_id":             "$WishBoxId",
		"WishPlayerCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"WishPlayerCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishItemReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询商品许愿并抽奖统计
func GroupPlayerWishAndDrawData(tid, startTime int64) *share_message.WishItemReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG)
	defer closeFun()

	//ids := GetWishBoxItemIdsByItemId(tid)

	match := bson.M{"WishItemId": tid}
	lookupData := bson.M{"from": for_game.TABLE_PLAYER_WISH_DATA, "localField": "ChallengeItemId", "foreignField": "WishBoxItemId", "as": "data"}
	match["CreateTime"] = bson.M{"$gt": startTime}
	match["data.CreateTime"] = bson.M{"$gt": startTime}
	match["data.Status"] = bson.M{"$ne": 2} // 过滤取消许愿的
	match["$expr"] = bson.M{"$eq": []string{"$DareId", "$data.PlayerId"}}
	group := bson.M{
		//"_id":             bson.M{"ChallengeItemId": "$ChallengeItemId"},
		"_id":                    nil,
		"WishAndDrawPlayerCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"WishAndDrawPlayerCount": 1,
	}

	pipeCond := []bson.M{
		{"$lookup": lookupData},
		{"$unwind": "$data"},
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishItemReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询商品许愿并抽奖统计
func GroupPlayerWishAndDrawDataByBoxItemId(tid, startTime, endTime int64) *share_message.WishBoxDetailReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG)
	defer closeFun()

	match := bson.M{"ChallengeItemId": tid}
	lookupData := bson.M{"from": for_game.TABLE_PLAYER_WISH_DATA, "localField": "ChallengeItemId", "foreignField": "WishBoxItemId", "as": "data"}
	if endTime > 0 {
		match["CreateTime"] = bson.M{"$gte": startTime, "$lt": endTime}
	} else {
		match["CreateTime"] = bson.M{"$gt": startTime}
	}

	match["data.CreateTime"] = bson.M{"$gt": startTime}
	match["data.Status"] = bson.M{"$ne": 2} // 过滤取消许愿的
	match["$expr"] = bson.M{"$eq": []string{"$DareId", "$data.PlayerId"}}
	cond4 := []interface{}{bson.M{"$eq": []interface{}{"$Result", true}}, 1, 0}
	group := bson.M{
		//"_id":             bson.M{"ChallengeItemId": "$ChallengeItemId"},
		"_id":                    nil,
		"WishAndDrawPlayerCount": bson.M{"$sum": 1},
		"LuckyWishCount":         bson.M{"$sum": bson.M{"$cond": cond4}},
		"LuckyGoldTotal":         bson.M{"$sum": "$DarePrice"},
	}

	project := bson.M{
		"WishAndDrawPlayerCount": 1,
		"LuckyWishCount":         1,
		"LuckyGoldTotal":         1,
	}

	pipeCond := []bson.M{
		{"$lookup": lookupData},
		{"$unwind": "$data"},
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishBoxDetailReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询盲盒收藏统计
func GroupPlayerWishCollectionByTime(tid, startTime, entTime int64) *share_message.WishBoxReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_COLLECTION)
	defer closeFun()

	match := bson.M{"WishBoxId": tid}
	//lookup := bson.M{"from": for_game.TABLE_PLAYER_WISH_DATA, "localField": "_id", "foreignField": "WishBoxItemId", "as": "item"}
	match["CreateTime"] = bson.M{"$gte": startTime, "$lt": entTime}

	one := &share_message.WishBoxReport{}
	count, err := col.Find(match).Count()
	easygo.PanicError(err)
	one.AddPlayerCount = easygo.NewInt32(count)

	return one

}

// 查询盲盒收藏统计
func GroupPlayerWishCollection(tid, startTime int64) *share_message.WishBoxReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_COLLECTION)
	defer closeFun()

	match := bson.M{"WishBoxId": tid}
	//lookup := bson.M{"from": for_game.TABLE_PLAYER_WISH_DATA, "localField": "_id", "foreignField": "WishBoxItemId", "as": "item"}
	match["CreateTime"] = bson.M{"$gt": startTime}

	one := &share_message.WishBoxReport{}
	count, err := col.Find(match).Count()
	easygo.PanicError(err)
	one.AddPlayerCount = easygo.NewInt32(count)

	return one

}

// 查询抽奖统计
func GroupWishLogByTime(tid, startTime, entTime int64) *share_message.WishBoxReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG)
	defer closeFun()

	match := bson.M{"WishBoxId": tid}
	match["CreateTime"] = bson.M{"$gte": startTime, "$lt": entTime}
	//count, err2 := col.Find(match).Count()
	//easygo.PanicError(err2)
	group := bson.M{
		"_id":             "$DareId",
		"DrawPlayerCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"DrawPlayerCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	one := &share_message.WishBoxReport{}
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	//one.DrawCount = easygo.NewInt32(count)
	return one

}

// 查询抽奖统计
func GroupWishLog(tid, startTime int64) *share_message.WishBoxReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG)
	defer closeFun()

	match := bson.M{"WishBoxId": tid}
	match["CreateTime"] = bson.M{"$gt": startTime}
	//count, err2 := col.Find(match).Count()
	//easygo.PanicError(err2)
	group := bson.M{
		"_id":             "$DareId",
		"DrawPlayerCount": bson.M{"$sum": 1},
	}

	project := bson.M{
		"DrawPlayerCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	one := &share_message.WishBoxReport{}
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	//one.DrawCount = easygo.NewInt32(count)
	return one

}

// 查询商品抽奖统计
func GroupWishLogByBoxId(tid, startTime, endTime int64) *share_message.WishBoxReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG)
	defer closeFun()

	match := bson.M{"WishBoxId": tid}
	if endTime > 0 {
		match["CreateTime"] = bson.M{"$gte": startTime, "$lt": endTime}
	} else {
		match["CreateTime"] = bson.M{"$gt": startTime}
	}
	group := bson.M{
		//"_id":             bson.M{"ChallengeItemId": "$ChallengeItemId"},
		"_id":           nil,
		"DrawGoldTotal": bson.M{"$sum": "$DarePrice"},
		"DrawCount":     bson.M{"$sum": 1},
	}

	project := bson.M{
		"DrawGoldTotal": 1,
		"DrawCount":     1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishBoxReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one

}

// 查询盲盒抽水
func GroupWishPoolPumpLogByTime(tid, startTime, entTime int64) *share_message.WishBoxReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL_PUMP_LOG)
	defer closeFun()

	match := bson.M{"BoxId": tid}
	match["CreateTime"] = bson.M{"$gte": startTime, "$lt": entTime}

	group := bson.M{
		"_id":             nil,
		"CommissionTotal": bson.M{"$sum": "$Price"},
	}

	project := bson.M{
		"CommissionTotal": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	one := &share_message.WishBoxReport{}
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	//one.DrawCount = easygo.NewInt32(count)
	return one
}

// 查询盲盒抽水
func GroupWishPoolPumpLog(tid, startTime, endTime int64) *share_message.WishBoxReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_POOL_PUMP_LOG)
	defer closeFun()

	match := bson.M{"BoxId": tid}
	if endTime > 0 {
		match["CreateTime"] = bson.M{"$gte": startTime, "$lt": endTime}
	} else {
		match["CreateTime"] = bson.M{"$gt": startTime}
	}

	group := bson.M{
		"_id":             nil,
		"CommissionTotal": bson.M{"$sum": "$Price"},
	}

	project := bson.M{
		"CommissionTotal": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	one := &share_message.WishBoxReport{}
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	//one.DrawCount = easygo.NewInt32(count)
	return one
}

// 生成许愿池商品报表
func MakeWishItemReport() {
	logs.Info("开始更新许愿池商品报表")
	startTime := int64(0)
	endTime := time.Now().Unix()
	job := GetReportJob(for_game.TABLE_REPORT_WISH_ITEM)
	if job != nil {
		startTime = job.GetTime()
	} else {
		InitWishItemReport()
	}

	report := &share_message.WishItemReport{}
	dList := GetAllWishItem()
	if len(dList) == 0 {
		return
	}

	dayTime := easygo.Get0ClockTimestamp(startTime)
	if startTime == 0 {
		dayTime = easygo.Get0ClockTimestamp(endTime)
	}
	for _, item := range dList {
		report = GetWishItemReport(item.GetId(), dayTime)
		if report == nil {
			report = &share_message.WishItemReport{
				Id:         easygo.NewString(easygo.IntToString(int(dayTime)) + easygo.IntToString(int(item.GetId()))),
				CreateTime: easygo.NewInt64(dayTime),
				ItemId:     easygo.NewInt64(item.GetId()),
				ItemName:   easygo.NewString(item.GetName()),
			}
		}

		one1 := GroupWishRecycleOrder(item.GetId(), startTime)
		one2 := GroupPlayerWishItem(item.GetId(), startTime)
		one3 := GroupPlayerWishData(item.GetId(), startTime)
		one4 := GroupPlayerWishAndDrawData(item.GetId(), startTime)
		report.RecycleCount = easygo.NewInt32(one1.GetRecycleCount() + report.GetRecycleCount())
		report.PlayerRecycleCount = easygo.NewInt32(one1.GetPlayerRecycleCount() + report.GetPlayerRecycleCount())
		report.OfficialRecycleCount = easygo.NewInt32(one1.GetOfficialRecycleCount() + report.GetOfficialRecycleCount())
		report.WinCount = easygo.NewInt32(one2.GetWinCount() + report.GetWinCount())
		report.PendConvertCount = easygo.NewInt32(one2.GetPendConvertCount() + report.GetPendConvertCount())
		report.ConvertCount = easygo.NewInt32(one2.GetConvertCount() + report.GetConvertCount())
		report.WishPlayerCount = easygo.NewInt32(one3.GetWishPlayerCount() + report.GetWishPlayerCount())
		report.WishAndDrawPlayerCount = easygo.NewInt32(one4.GetWishAndDrawPlayerCount() + report.GetWishAndDrawPlayerCount())

		UpdateWishItemReport(report)
	}
	logs.Info("完成更新许愿池商品报表")
	//写报表生成进度
	MakeReportJob(for_game.TABLE_REPORT_WISH_ITEM, endTime)
}

func InitWishBoxDetailReport() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_DETAIL)
	defer closeFun()

	queryBson := bson.M{}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//更新许愿池盲盒报表
func UpdateWishBoxDetailReport(report *share_message.WishBoxDetailReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_DETAIL)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetId()}, bson.M{"$set": report})
	easygo.PanicError(err)

}

//更新许愿池盲盒报表
func UpdateWishBoxDetailReportTemp(report *share_message.WishBoxDetailReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_DETAIL_TEMP)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetId()}, bson.M{"$set": report})
	easygo.PanicError(err)

}

//查询许愿池盲盒报表
func GetWishBoxDetailReport(id, startTime int64) *share_message.WishBoxDetailReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_DETAIL)
	defer closeFun()
	report := &share_message.WishBoxDetailReport{}
	err := col.Find(bson.M{"CreateTime": startTime, "ItemId": id}).One(&report)

	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return report
}

// 生成许愿池盲盒详情报表
func MakeWishBoxDetailReportByTime(id, sTime, eTime int64) {
	logs.Info("开始更新许愿池盲盒详情报表")
	startTime := sTime
	endTime := eTime

	//report := &share_message.WishBoxDetailReport{}
	//var list []*share_message.WishBoxDetailReport
	dList := GetAllWishBoxItemByBoxId(id)
	if len(dList) == 0 {
		return
	}
	for _, item := range dList {
		report := &share_message.WishBoxDetailReport{
			Id:         easygo.NewString(easygo.IntToString(int(startTime)) + easygo.IntToString(int(item.GetId()))),
			CreateTime: easygo.NewInt64(startTime),
			ItemId:     easygo.NewInt64(item.GetId()),
			Style:      easygo.NewInt32(item.GetStyle()),
			WishBoxId:  easygo.NewInt64(item.GetWishBoxId()),
		}

		one1 := GroupWishRecycleOrderByBoxItemIdByTime(item.GetId(), startTime, endTime)
		one2 := GroupPlayerWishItemByBoxItemIdByTime(item.GetId(), startTime, endTime)
		one3 := GroupPlayerWishDataByBoxItemIdByTime(item.GetId(), startTime, endTime)
		one4 := GroupPlayerWishAndDrawDataByBoxItemIdByTime(item.GetId(), startTime, endTime)
		report.RecycleCount = easygo.NewInt32(one1.GetRecycleCount() + report.GetRecycleCount())
		report.PlayerRecycleCount = easygo.NewInt32(one1.GetPlayerRecycleCount() + report.GetPlayerRecycleCount())
		report.OfficialRecycleCount = easygo.NewInt32(one1.GetOfficialRecycleCount() + report.GetOfficialRecycleCount())
		//divisor := int64(10)
		//// 回收价格
		//playerRecycleGoldTotal := (one1.GetPlayerRecycleGoldTotal()*divisor*8) / 100
		//officialRecycleGoldTotal := (one1.GetOfficialRecycleGoldTotal()*divisor*7) / 100
		//recycleGoldTotal := playerRecycleGoldTotal + officialRecycleGoldTotal
		report.RecycleGoldTotal = easygo.NewInt64(one1.GetRecycleGoldTotal() + report.GetRecycleGoldTotal())
		report.PlayerRecycleGoldTotal = easygo.NewInt64(one1.GetPlayerRecycleGoldTotal() + report.GetPlayerRecycleGoldTotal())
		report.OfficialRecycleGoldTotal = easygo.NewInt64(one1.GetOfficialRecycleGoldTotal() + report.GetOfficialRecycleGoldTotal())
		report.WinCount = easygo.NewInt32(one2.GetWinCount() + report.GetWinCount())

		report.ConvertCount = easygo.NewInt32(one2.GetConvertCount() + report.GetConvertCount())
		report.ConvertGoodsPriceTotal = easygo.NewInt64(one2.GetConvertGoodsPriceTotal() + report.GetConvertGoodsPriceTotal())
		report.PendConvertCount = easygo.NewInt32(one2.GetPendConvertCount() + report.GetPendConvertCount())
		report.PendConvertPriceTotal = easygo.NewInt64(one2.GetPendConvertPriceTotal() + report.GetPendConvertPriceTotal())
		report.WishPlayerCount = easygo.NewInt32(one3.GetWishPlayerCount() + report.GetWishPlayerCount())
		report.WishAndDrawPlayerCount = easygo.NewInt32(one4.GetWishAndDrawPlayerCount() + report.GetWishAndDrawPlayerCount())
		report.LuckyWishCount = easygo.NewInt32(one4.GetLuckyWishCount() + report.GetLuckyWishCount())
		report.LuckyGoldTotal = easygo.NewInt64(one4.GetLuckyGoldTotal() + report.GetLuckyGoldTotal())

		// 盈亏
		profit := one4.GetLuckyGoldTotal() - one1.GetRecycleGoldTotal() - (one2.GetConvertGoodsPriceTotal() + one2.GetPendConvertPriceTotal())
		report.Profit = easygo.NewInt64(profit + report.GetProfit())

		//list = append(list, report)

		UpdateWishBoxDetailReportTemp(report)
	}
	logs.Info("完成更新许愿池盲盒详情报表")

}

// 生成许愿池盲盒报表
func MakeWishBoxReportByTime(sTime, eTime int64) {
	logs.Info("开始更新许愿池盲盒报表")

	startTime := sTime
	endTime := eTime

	dList := GetAllWishBox()
	if len(dList) == 0 {
		return
	}

	for _, item := range dList {
		report := &share_message.WishBoxReport{
			Id:         easygo.NewString(easygo.IntToString(int(startTime)) + easygo.IntToString(int(item.GetId()))),
			CreateTime: easygo.NewInt64(startTime),
			BoxId:      easygo.NewInt64(item.GetId()),
			BoxName:    easygo.NewString(item.GetName()),
		}

		//one1 := GroupWishRecycleOrderByBoxItemId(item.GetId(), startTime)
		//one2 := GroupPlayerWishItemByBoxItemId(item.GetId(), startTime)
		//one3 := GroupPlayerWishDataByBoxItemId(item.GetId(), startTime)
		//one4 := GroupPlayerWishAndDrawDataByBoxItemId(item.GetId(), startTime)
		one1 := GroupWishBoxDetailReportTemp(item.GetId(), startTime)
		one2 := GroupWishLogByTime(item.GetId(), startTime, endTime)
		one22 := GroupWishLogByBoxId(item.GetId(), startTime, 0)
		one3 := GroupPlayerWishCollectionByTime(item.GetId(), startTime, endTime)
		one4 := GroupWishPoolPumpLogByTime(item.GetId(), startTime, endTime)
		report.DrawPlayerCount = easygo.NewInt32(one2.GetDrawPlayerCount() + report.GetDrawPlayerCount())
		report.DrawCount = easygo.NewInt32(one22.GetDrawCount() + report.GetDrawCount())
		report.LuckyGoldTotal = easygo.NewInt64(one22.GetLuckyGoldTotal() + report.GetLuckyGoldTotal())
		report.AddPlayerCount = easygo.NewInt32(one3.GetAddPlayerCount() + report.GetAddPlayerCount())

		report.RecycleCount = easygo.NewInt32(one1.GetRecycleCount())
		report.RecycleGoldTotal = easygo.NewInt64(one1.GetRecycleGoldTotal())
		report.PlayerRecycleCount = easygo.NewInt32(one1.GetPlayerRecycleCount())
		report.OfficialRecycleCount = easygo.NewInt32(one1.GetOfficialRecycleCount())
		//report.WinCount = easygo.NewInt32(one2.GetWinCount() + report.GetWinCount())
		report.PendConvertCount = easygo.NewInt32(one1.GetPendConvertCount())
		report.ConvertCount = easygo.NewInt32(one1.GetConvertCount())
		report.WishPlayerCount = easygo.NewInt32(one1.GetWishPlayerCount())

		report.LuckyWishCount = easygo.NewInt32(one1.GetLuckyWishCount())
		report.PlayerRecycleGoldTotal = easygo.NewInt64(one1.GetPlayerRecycleGoldTotal())
		report.OfficialRecycleGoldTotal = easygo.NewInt64(one1.GetOfficialRecycleGoldTotal())
		report.ConvertGoodsPriceTotal = easygo.NewInt64(one1.GetConvertGoodsPriceTotal())
		report.PendConvertGoodsPriceTotal = easygo.NewInt64(one1.GetPendConvertGoodsPriceTotal())
		report.Profit = easygo.NewInt64(one1.GetProfit())

		// 官方抽水
		report.CommissionTotal = easygo.NewInt32(one4.GetCommissionTotal())

		UpdateWishBoxReportTemp(report)
	}
	logs.Info("完成更新许愿池盲盒报表")

}

// 生成许愿池盲盒详情报表
func MakeWishBoxDetailReport() {
	logs.Info("开始更新许愿池盲盒详情报表")
	startTime := int64(0)
	endTime := time.Now().Unix()
	job := GetReportJob(for_game.TABLE_REPORT_WISH_BOX_DETAIL)
	if job != nil {
		startTime = job.GetTime()
	} else {
		InitWishBoxDetailReport()
	}

	//report := &share_message.WishBoxDetailReport{}
	dList := GetAllWishBoxItem()
	if len(dList) == 0 {
		return
	}

	dayTime := easygo.Get0ClockTimestamp(startTime)
	if startTime == 0 {
		dayTime = easygo.Get0ClockTimestamp(endTime)
	}
	for _, item := range dList {
		report := GetWishBoxDetailReport(item.GetId(), dayTime)
		if report == nil {
			report = &share_message.WishBoxDetailReport{
				Id:         easygo.NewString(easygo.IntToString(int(dayTime)) + easygo.IntToString(int(item.GetId()))),
				CreateTime: easygo.NewInt64(dayTime),
				ItemId:     easygo.NewInt64(item.GetId()),
				Style:      easygo.NewInt32(item.GetStyle()),
				WishBoxId:  easygo.NewInt64(item.GetWishBoxId()),
			}
		}

		one1 := GroupWishRecycleOrderByBoxItemId(item.GetId(), dayTime, 0)
		one2 := GroupPlayerWishItemByBoxItemId(item.GetId(), dayTime, 0)
		one3 := GroupPlayerWishDataByBoxItemId(item.GetId(), dayTime, 0)
		one4 := GroupPlayerWishAndDrawDataByBoxItemId(item.GetId(), dayTime, 0)
		report.RecycleCount = easygo.NewInt32(one1.GetRecycleCount())
		report.PlayerRecycleCount = easygo.NewInt32(one1.GetPlayerRecycleCount())
		report.OfficialRecycleCount = easygo.NewInt32(one1.GetOfficialRecycleCount())
		//divisor := int64(10)
		//// 回收价格
		//playerRecycleGoldTotal := (one1.GetPlayerRecycleGoldTotal()*divisor*8) / 100
		//officialRecycleGoldTotal := (one1.GetOfficialRecycleGoldTotal()*divisor*7) / 100
		//recycleGoldTotal := playerRecycleGoldTotal + officialRecycleGoldTotal
		report.RecycleGoldTotal = easygo.NewInt64(one1.GetRecycleGoldTotal())
		report.ProductDiamondTotal = easygo.NewInt64(one1.GetProductDiamondTotal())
		report.PlayerRecycleGoldTotal = easygo.NewInt64(one1.GetPlayerRecycleGoldTotal())
		report.OfficialRecycleGoldTotal = easygo.NewInt64(one1.GetOfficialRecycleGoldTotal())
		report.WinCount = easygo.NewInt32(one2.GetWinCount())

		report.ConvertCount = easygo.NewInt32(one2.GetConvertCount())
		report.ConvertGoodsPriceTotal = easygo.NewInt64(one2.GetConvertGoodsPriceTotal())
		report.PendConvertCount = easygo.NewInt32(one2.GetPendConvertCount() + one2.GetAuditConvertCount())                // 待兑换 + 审核中
		report.PendConvertPriceTotal = easygo.NewInt64(one2.GetPendConvertPriceTotal() + one2.GetAuditConvertPriceTotal()) // 待兑换 + 审核中
		report.WishPlayerCount = easygo.NewInt32(one3.GetWishPlayerCount())
		report.WishAndDrawPlayerCount = easygo.NewInt32(one4.GetWishAndDrawPlayerCount())
		report.LuckyWishCount = easygo.NewInt32(one4.GetLuckyWishCount())
		//report.LuckyGoldTotal = easygo.NewInt64(one4.GetLuckyGoldTotal() + report.GetLuckyGoldTotal())
		report.WinItemPriceTotal = easygo.NewInt64(one2.GetWinItemPriceTotal())

		// 盈亏
		profit := one4.GetLuckyGoldTotal() - one1.GetRecycleGoldTotal() - (one2.GetConvertGoodsPriceTotal() + one2.GetPendConvertPriceTotal())
		report.Profit = easygo.NewInt64(profit)

		// 新版需求
		// 发货订单
		deliverCount := GetWishDeliveryOrderCountByBoxItemId(item.GetId(), dayTime)
		report.DeliverCount = easygo.NewInt64(deliverCount)
		one5 := GetWishDeliveryOrderPlayerCountByBoxItemId(item.GetId(), dayTime)
		report.DeliverPlayerCount = easygo.NewInt32(one5.GetDeliverPlayerCount())
		// 回收订单
		one6 := GetWishRecyclePlayerCountByBoxItemId(item.GetId(), dayTime)
		report.RecyclePlayerCount = easygo.NewInt32(one6.GetRecyclePlayerCount())

		UpdateWishBoxDetailReport(report)
	}
	logs.Info("完成更新许愿池盲盒详情报表")
	//写报表生成进度
	MakeReportJob(for_game.TABLE_REPORT_WISH_BOX_DETAIL, endTime)
}

// 生成许愿池盲盒详情报表-周表
func MakeWishBoxDetailReportByWeek() {
	logs.Info("开始更新许愿池盲盒详情报表-周表")

	//report := &share_message.WishBoxDetailReport{}
	dList := GetAllWishBoxItem()
	if len(dList) == 0 {
		return
	}

	lastWeekEndTime := easygo.GetWeek0ClockOfTimestamp(time.Now().Unix())
	lastWeekStartTime := lastWeekEndTime - 3600*24*7
	startTime := lastWeekStartTime
	endTime := lastWeekStartTime

	ctBson := bson.M{"$gte": startTime, "$lt": endTime} //前闭后开
	var saveData []interface{}
	for _, item := range dList {
		wishBoxItemId := item.GetId()
		wishPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA, "PlayerId", bson.M{"WishBoxItemId": wishBoxItemId, "CreateTime": ctBson}))
		recyclePlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER, "PlayerId", bson.M{"WishRecycleItem.WishBoxItemId": wishBoxItemId, "RecycleTime": ctBson}))
		deliverCount := for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG, bson.M{"WishBoxItem": wishBoxItemId, "Status": 1, "CreateTime": ctBson})
		deliverPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG, "PlayerId", bson.M{"WishBoxItem": wishBoxItemId, "Status": 1, "CreateTime": ctBson}))

		one1 := GroupWishRecycleOrderByBoxItemId(wishBoxItemId, startTime, endTime)
		one2 := GroupPlayerWishItemByBoxItemId(wishBoxItemId, startTime, endTime)
		one4 := GroupPlayerWishAndDrawDataByBoxItemId(item.GetId(), startTime, endTime)
		// 盈亏
		profit := one4.GetLuckyGoldTotal() - one1.GetRecycleGoldTotal() - (one2.GetConvertGoodsPriceTotal() + one2.GetPendConvertPriceTotal())
		report := &share_message.WishBoxDetailReportWeek{
			Id:                       easygo.NewString(easygo.IntToString(int(startTime)) + easygo.IntToString(int(item.GetId()))),
			WishPlayerCount:          easygo.NewInt32(wishPlayerCount),
			WishAndDrawPlayerCount:   easygo.NewInt32(one4.GetWishAndDrawPlayerCount()), // 许愿并抽奖人数
			WinCount:                 easygo.NewInt32(one2.GetWinCount()),
			LuckyWishCount:           easygo.NewInt32(one4.GetLuckyWishCount()),
			ConvertCount:             easygo.NewInt32(one2.GetConvertCount()),
			RecycleCount:             easygo.NewInt32(one1.GetRecycleCount()),
			PlayerRecycleCount:       easygo.NewInt32(one1.GetPlayerRecycleCount()),
			OfficialRecycleCount:     easygo.NewInt32(one1.GetOfficialRecycleCount()),
			RecycleGoldTotal:         easygo.NewInt64(one1.GetRecycleGoldTotal()),
			ProductDiamondTotal:      easygo.NewInt64(one1.GetProductDiamondTotal()),
			LuckyGoldTotal:           easygo.NewInt64(one4.GetLuckyGoldTotal()),
			ConvertGoodsPriceTotal:   easygo.NewInt64(one2.GetConvertGoodsPriceTotal()),
			PendConvertCount:         easygo.NewInt32(one2.GetPendConvertCount() + one2.GetAuditConvertCount()), // 待兑换 + 审核中
			PlayerRecycleGoldTotal:   easygo.NewInt64(one1.GetPlayerRecycleGoldTotal()),
			OfficialRecycleGoldTotal: easygo.NewInt64(one1.GetOfficialRecycleGoldTotal()),
			PendConvertPriceTotal:    easygo.NewInt64(one2.GetPendConvertPriceTotal() + one2.GetAuditConvertPriceTotal()), // 待兑换 + 审核中
			Profit:                   easygo.NewInt64(profit),

			RecyclePlayerCount: easygo.NewInt32(recyclePlayerCount),
			DeliverCount:       easygo.NewInt64(deliverCount),
			DeliverPlayerCount: easygo.NewInt32(deliverPlayerCount),
			WishBoxId:          easygo.NewInt64(item.GetWishBoxId()),
			// WishBoxReportId: easygo.NewString(),
			CreateTime: easygo.NewInt64(endTime),
			ItemId:     easygo.NewInt64(item.GetId()),
			Style:      easygo.NewInt32(item.GetStyle()),
			StartTime:  easygo.NewInt64(startTime),
			EndTime:    easygo.NewInt64(endTime),
		}
		saveData = append(saveData, bson.M{"_id": report.GetId()}, report)
	}
	//logs.Info("许愿池盲盒详情报表-周表:", saveData)
	for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_DETAIL_WEEK, saveData)
	logs.Info("完成更新许愿池盲盒详情报表-周表")
}

func MakeWishBoxDetailReportByMonth() {
	logs.Info("开始更新许愿池盲盒详情报表-月报")
	//report := &share_message.WishBoxDetailReport{}
	dList := GetAllWishBoxItem()
	if len(dList) == 0 {
		return
	}
	now := time.Now()
	lastMonthFirstDay := now.AddDate(0, -1, -now.Day()+1)
	lastMonthStartTime := time.Date(lastMonthFirstDay.Year(), lastMonthFirstDay.Month(), lastMonthFirstDay.Day(), 0, 0, 0, 0, now.Location()).Unix()
	lastMonthEndTime := easygo.GetMonth0ClockTimestamp()

	startTime := lastMonthStartTime
	endTime := lastMonthEndTime

	ctBson := bson.M{"$gte": startTime, "$lt": endTime} //前闭后开
	var saveData []interface{}
	for _, item := range dList {
		wishBoxItemId := item.GetId()
		wishPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA, "PlayerId", bson.M{"WishBoxItemId": wishBoxItemId, "CreateTime": ctBson}))
		recyclePlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER, "PlayerId", bson.M{"WishRecycleItem.WishBoxItemId": wishBoxItemId, "RecycleTime": ctBson}))
		deliverCount := for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG, bson.M{"WishBoxItem": wishBoxItemId, "Status": 1, "CreateTime": ctBson})
		deliverPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG, "PlayerId", bson.M{"WishBoxItem": wishBoxItemId, "Status": 1, "CreateTime": ctBson}))

		one1 := GroupWishRecycleOrderByBoxItemId(wishBoxItemId, startTime, endTime)
		one2 := GroupPlayerWishItemByBoxItemId(wishBoxItemId, startTime, endTime)
		one4 := GroupPlayerWishAndDrawDataByBoxItemId(item.GetId(), startTime, endTime)
		// 盈亏
		profit := one4.GetLuckyGoldTotal() - one1.GetRecycleGoldTotal() - (one2.GetConvertGoodsPriceTotal() + one2.GetPendConvertPriceTotal())
		report := &share_message.WishBoxDetailReportMonth{
			Id:                       easygo.NewString(easygo.IntToString(int(startTime)) + easygo.IntToString(int(item.GetId()))),
			WishPlayerCount:          easygo.NewInt32(wishPlayerCount),
			WishAndDrawPlayerCount:   easygo.NewInt32(one4.GetWishAndDrawPlayerCount()), // 许愿并抽奖人数
			WinCount:                 easygo.NewInt32(one2.GetWinCount()),
			LuckyWishCount:           easygo.NewInt32(one4.GetLuckyWishCount()),
			ConvertCount:             easygo.NewInt32(one2.GetConvertCount()),
			RecycleCount:             easygo.NewInt32(one1.GetRecycleCount()),
			PlayerRecycleCount:       easygo.NewInt32(one1.GetPlayerRecycleCount()),
			OfficialRecycleCount:     easygo.NewInt32(one1.GetOfficialRecycleCount()),
			RecycleGoldTotal:         easygo.NewInt64(one1.GetRecycleGoldTotal()),
			ProductDiamondTotal:      easygo.NewInt64(one1.GetProductDiamondTotal()),
			LuckyGoldTotal:           easygo.NewInt64(one4.GetLuckyGoldTotal()),
			ConvertGoodsPriceTotal:   easygo.NewInt64(one2.GetConvertGoodsPriceTotal()),
			PendConvertCount:         easygo.NewInt32(one2.GetPendConvertCount() + one2.GetAuditConvertCount()), // 待兑换 + 审核中
			PlayerRecycleGoldTotal:   easygo.NewInt64(one1.GetPlayerRecycleGoldTotal()),
			OfficialRecycleGoldTotal: easygo.NewInt64(one1.GetOfficialRecycleGoldTotal()),
			PendConvertPriceTotal:    easygo.NewInt64(one2.GetPendConvertPriceTotal() + one2.GetAuditConvertPriceTotal()), // 待兑换 + 审核中
			Profit:                   easygo.NewInt64(profit),

			RecyclePlayerCount: easygo.NewInt32(recyclePlayerCount),
			DeliverCount:       easygo.NewInt64(deliverCount),
			DeliverPlayerCount: easygo.NewInt32(deliverPlayerCount),
			WishBoxId:          easygo.NewInt64(item.GetWishBoxId()),
			// WishBoxReportId: easygo.NewString(),
			CreateTime: easygo.NewInt64(endTime),
			ItemId:     easygo.NewInt64(item.GetId()),
			Style:      easygo.NewInt32(item.GetStyle()),
			StartTime:  easygo.NewInt64(startTime),
			EndTime:    easygo.NewInt64(endTime),
		}
		saveData = append(saveData, bson.M{"_id": report.GetId()}, report)
	}
	//logs.Info("许愿池盲盒详情报表-月报saveData: ", saveData)
	for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_DETAIL_MONTH, saveData)
	logs.Info("完成更新许愿池盲盒详情报表-月报")
}

// 聚合查询盲盒详情报表
func GroupWishBoxDetailReport(bid, starTime, endTime int64) *share_message.WishBoxReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_DETAIL)
	defer closeFun()

	match := bson.M{"WishBoxId": bid}
	if endTime > 0 {
		match["CreateTime"] = bson.M{"$gte": starTime, "$lt": endTime}
	} else {
		match["CreateTime"] = starTime
	}
	group := bson.M{
		"_id":                        "$WishBoxId",
		"RecycleCount":               bson.M{"$sum": "$RecycleCount"},
		"PlayerRecycleCount":         bson.M{"$sum": "$PlayerRecycleCount"},
		"OfficialRecycleCount":       bson.M{"$sum": "$OfficialRecycleCount"},
		"WishPlayerCount":            bson.M{"$sum": "$WishPlayerCount"},
		"WishAndDrawPlayerCount":     bson.M{"$sum": "$WishAndDrawPlayerCount"},
		"WinCount":                   bson.M{"$sum": "$WinCount"},
		"LuckyWishCount":             bson.M{"$sum": "$LuckyWishCount"},
		"ConvertCount":               bson.M{"$sum": "$ConvertCount"},
		"PendConvertCount":           bson.M{"$sum": "$PendConvertCount"},
		"PlayerRecycleGoldTotal":     bson.M{"$sum": "$PlayerRecycleGoldTotal"},
		"OfficialRecycleGoldTotal":   bson.M{"$sum": "$OfficialRecycleGoldTotal"},
		"ConvertGoodsPriceTotal":     bson.M{"$sum": "$ConvertGoodsPriceTotal"},
		"PendConvertGoodsPriceTotal": bson.M{"$sum": "$PendConvertPriceTotal"},
		"Profit":                     bson.M{"$sum": "$Profit"},

		"RecycleGoldTotal":    bson.M{"$sum": "$RecycleGoldTotal"},
		"ProductDiamondTotal": bson.M{"$sum": "$ProductDiamondTotal"},
		"DrawCount":           bson.M{"$sum": "$LuckyWishCount"},
		"DrawGoldTotal":       bson.M{"$sum": "$LuckyGoldTotal"},
		"WinItemPriceTotal":   bson.M{"$sum": "$WinItemPriceTotal"},

		"DeliverCount":       bson.M{"$sum": "$DeliverCount"},
		"DeliverPlayerCount": bson.M{"$sum": "$DeliverPlayerCount"},
		"RecyclePlayerCount": bson.M{"$sum": "$RecyclePlayerCount"},
	}

	project := bson.M{
		"RecycleCount":               1,
		"PlayerRecycleCount":         1,
		"OfficialRecycleCount":       1,
		"WishPlayerCount":            1,
		"WishAndDrawPlayerCount":     1,
		"WinCount":                   1,
		"LuckyWishCount":             1,
		"ConvertCount":               1,
		"PendConvertCount":           1,
		"PlayerRecycleGoldTotal":     1,
		"OfficialRecycleGoldTotal":   1,
		"ConvertGoodsPriceTotal":     1,
		"PendConvertGoodsPriceTotal": 1,
		"Profit":                     1,

		"RecycleGoldTotal":    1,
		"ProductDiamondTotal": 1,
		"DrawCount":           1,
		"DrawGoldTotal":       1,
		"WinItemPriceTotal":   1,

		"DeliverCount":       1,
		"DeliverPlayerCount": 1,
		"RecyclePlayerCount": 1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishBoxReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one
}

// 聚合查询盲盒详情报表
func GroupWishBoxDetailReportTemp(bid, starTime int64) *share_message.WishBoxReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_DETAIL_TEMP)
	defer closeFun()

	match := bson.M{"WishBoxId": bid}
	match["CreateTime"] = starTime
	group := bson.M{
		"_id":                        "$WishBoxId",
		"RecycleCount":               bson.M{"$sum": "$RecycleCount"},
		"PlayerRecycleCount":         bson.M{"$sum": "$PlayerRecycleCount"},
		"OfficialRecycleCount":       bson.M{"$sum": "$OfficialRecycleCount"},
		"WishPlayerCount":            bson.M{"$sum": "$WishPlayerCount"},
		"WishAndDrawPlayerCount":     bson.M{"$sum": "$WishAndDrawPlayerCount"},
		"WinCount":                   bson.M{"$sum": "$WinCount"},
		"LuckyWishCount":             bson.M{"$sum": "$LuckyWishCount"},
		"ConvertCount":               bson.M{"$sum": "$ConvertCount"},
		"PendConvertCount":           bson.M{"$sum": "$PendConvertCount"},
		"PlayerRecycleGoldTotal":     bson.M{"$sum": "$PlayerRecycleGoldTotal"},
		"OfficialRecycleGoldTotal":   bson.M{"$sum": "$OfficialRecycleGoldTotal"},
		"ConvertGoodsPriceTotal":     bson.M{"$sum": "$ConvertGoodsPriceTotal"},
		"PendConvertGoodsPriceTotal": bson.M{"$sum": "$PendConvertPriceTotal"},
		"Profit":                     bson.M{"$sum": "$Profit"},

		"RecycleGoldTotal": bson.M{"$sum": "$RecycleGoldTotal"},
		"DrawCount":        bson.M{"$sum": "$LuckyWishCount"},
		"DrawGoldTotal":    bson.M{"$sum": "$LuckyGoldTotal"},
	}

	project := bson.M{
		"RecycleCount":               1,
		"PlayerRecycleCount":         1,
		"OfficialRecycleCount":       1,
		"WishPlayerCount":            1,
		"WishAndDrawPlayerCount":     1,
		"WinCount":                   1,
		"LuckyWishCount":             1,
		"ConvertCount":               1,
		"PendConvertCount":           1,
		"PlayerRecycleGoldTotal":     1,
		"OfficialRecycleGoldTotal":   1,
		"ConvertGoodsPriceTotal":     1,
		"PendConvertGoodsPriceTotal": 1,
		"Profit":                     1,

		"RecycleGoldTotal": 1,
		"DrawCount":        1,
		"DrawGoldTotal":    1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishBoxReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one
}

func InitWishBoxReport() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX)
	defer closeFun()

	queryBson := bson.M{}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//更新许愿池盲盒报表
func UpdateWishBoxReportTemp(report *share_message.WishBoxReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_TEMP)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetId()}, bson.M{"$set": report})
	easygo.PanicError(err)

}

//更新许愿池盲盒报表
func UpdateWishBoxReport(report *share_message.WishBoxReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetId()}, bson.M{"$set": report})
	easygo.PanicError(err)

}

//查询许愿池盲盒报表
func GetWishBoxReport(id, startTime int64) *share_message.WishBoxReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX)
	defer closeFun()
	report := &share_message.WishBoxReport{}
	err := col.Find(bson.M{"CreateTime": startTime, "BoxId": id}).One(&report)

	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return report
}

// 生成许愿池盲盒报表
func MakeWishBoxReport() {
	logs.Info("开始更新许愿池盲盒报表")

	startTime := int64(0)
	endTime := time.Now().Unix()
	job := GetReportJob(for_game.TABLE_REPORT_WISH_BOX)
	if job != nil {
		startTime = job.GetTime()
	} else {
		InitWishBoxReport()
	}

	//report := &share_message.WishBoxReport{}
	dList := GetAllWishBox()
	if len(dList) == 0 {
		return
	}

	dayTime := easygo.Get0ClockTimestamp(startTime)
	if startTime == 0 {
		dayTime = easygo.Get0ClockTimestamp(endTime)
	}
	for _, item := range dList {
		report := GetWishBoxReport(item.GetId(), dayTime)
		if report == nil {
			report = &share_message.WishBoxReport{
				Id:         easygo.NewString(easygo.IntToString(int(dayTime)) + easygo.IntToString(int(item.GetId()))),
				CreateTime: easygo.NewInt64(dayTime),
				BoxId:      easygo.NewInt64(item.GetId()),
				BoxName:    easygo.NewString(item.GetName()),
			}
		}

		//one1 := GroupWishRecycleOrderByBoxItemId(item.GetId(), startTime)
		//one2 := GroupPlayerWishItemByBoxItemId(item.GetId(), startTime)
		//one3 := GroupPlayerWishDataByBoxItemId(item.GetId(), startTime)
		//one4 := GroupPlayerWishAndDrawDataByBoxItemId(item.GetId(), startTime)
		one1 := GroupWishBoxDetailReport(item.GetId(), dayTime, 0)
		one2 := GroupWishLog(item.GetId(), dayTime)
		one22 := GroupWishLogByBoxId(item.GetId(), dayTime, 0)
		one3 := GroupPlayerWishCollection(item.GetId(), dayTime)
		one4 := GroupWishPoolPumpLog(item.GetId(), dayTime, 0)
		report.DrawPlayerCount = easygo.NewInt32(one2.GetDrawPlayerCount())
		report.DrawCount = easygo.NewInt32(one22.GetDrawCount())
		report.LuckyGoldTotal = easygo.NewInt64(one22.GetLuckyGoldTotal())
		report.AddPlayerCount = easygo.NewInt32(one3.GetAddPlayerCount())

		report.RecycleCount = easygo.NewInt32(one1.GetRecycleCount())
		report.RecycleGoldTotal = easygo.NewInt64(one1.GetRecycleGoldTotal())
		report.ProductDiamondTotal = easygo.NewInt64(one1.GetProductDiamondTotal())

		report.PlayerRecycleCount = easygo.NewInt32(one1.GetPlayerRecycleCount())
		report.OfficialRecycleCount = easygo.NewInt32(one1.GetOfficialRecycleCount())
		//report.WinCount = easygo.NewInt32(one2.GetWinCount() + report.GetWinCount())
		report.PendConvertCount = easygo.NewInt32(one1.GetPendConvertCount())
		report.ConvertCount = easygo.NewInt32(one1.GetConvertCount())
		report.WishPlayerCount = easygo.NewInt32(one1.GetWishPlayerCount())

		report.LuckyWishCount = easygo.NewInt32(one1.GetLuckyWishCount())
		report.PlayerRecycleGoldTotal = easygo.NewInt64(one1.GetPlayerRecycleGoldTotal())
		report.OfficialRecycleGoldTotal = easygo.NewInt64(one1.GetOfficialRecycleGoldTotal())
		report.ConvertGoodsPriceTotal = easygo.NewInt64(one1.GetConvertGoodsPriceTotal())
		report.PendConvertGoodsPriceTotal = easygo.NewInt64(one1.GetPendConvertGoodsPriceTotal())
		report.Profit = easygo.NewInt64(one1.GetProfit())

		//report.ConvertDiamondPlayerCount = easygo.NewInt32()
		report.WinItemPriceTotal = easygo.NewInt64(one1.GetWinItemPriceTotal())

		// 官方抽水
		report.CommissionTotal = easygo.NewInt32(one4.GetCommissionTotal())

		// 新版需求
		// 发货订单
		report.DeliverCount = easygo.NewInt64(one1.GetDeliverCount())
		report.DeliverPlayerCount = easygo.NewInt32(one1.GetDeliverPlayerCount())
		// 回收订单
		report.RecyclePlayerCount = easygo.NewInt32(one1.GetRecyclePlayerCount())

		UpdateWishBoxReport(report)
	}
	logs.Info("完成更新许愿池盲盒报表")
	//写报表生成进度
	MakeReportJob(for_game.TABLE_REPORT_WISH_BOX, endTime)

}

// 生成许愿池盲盒报表-周表
func MakeWishBoxReportByWeek() {
	logs.Info("开始更新许愿池盲盒报表-周表")
	dList := GetAllWishBox()
	if len(dList) == 0 {
		return
	}
	lastWeekEndTime := easygo.GetWeek0ClockOfTimestamp(time.Now().Unix())
	lastWeekStartTime := lastWeekEndTime - 3600*24*7
	var saveData []interface{}
	startTime := lastWeekStartTime
	endTime := lastWeekEndTime
	ct := bson.M{"$gte": startTime, "$lt": endTime}
	for _, item := range dList {
		wishBoxId := item.GetId()
		//wishBoxDetailReportWeek := GroupWishBoxReportDetailByWeek(wishBoxId,lastWeekStartTime,lastWeekEndTime)
		wishPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA, "PlayerId", bson.M{"WishBoxId": wishBoxId, "CreateTime": ct}))
		addPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_COLLECTION, "PlayerId", bson.M{"WishBoxId": wishBoxId, "CreateTime": ct}))
		//luckyPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG, "DareId", bson.M{"CreateTime": ct }))
		DrawPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG, "DareId", bson.M{"WishBoxId": wishBoxId, "CreateTime": ct}))
		one1 := GroupWishBoxDetailReport(wishBoxId, startTime, endTime)
		one22 := GroupWishLogByBoxId(wishBoxId, startTime, endTime)
		one4 := GroupWishPoolPumpLog(wishBoxId, startTime, endTime)
		/*diamondExchangePeopleBson := bson.M{"PayType": 1, "SourceType": for_game.DIAMOND_TYPE_EXCHANGE_IN, "RecycleTime": ct}
		convertDiamondPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_DIAMOND_CHANGELOG, "PlayerId", diamondExchangePeopleBson))*/
		wishBoxReportWeek := &share_message.WishBoxReportWeek{
			Id:                         easygo.NewString(easygo.IntToString(int(lastWeekStartTime)) + easygo.IntToString(int(item.GetId()))),
			WishPlayerCount:            easygo.NewInt32(wishPlayerCount),                      //许愿人数
			AddPlayerCount:             easygo.NewInt32(addPlayerCount),                       //收藏人数
			DrawPlayerCount:            easygo.NewInt32(DrawPlayerCount),                      // 抽奖人数
			DrawCount:                  easygo.NewInt32(one22.GetDrawCount()),                 // 抽奖次数
			LuckyWishCount:             easygo.NewInt32(one1.GetLuckyWishCount()),             // 抽中许愿款次数
			ConvertCount:               easygo.NewInt32(one1.GetConvertCount()),               // 兑换次数
			RecycleCount:               easygo.NewInt32(one1.GetRecycleCount()),               // 回收次数
			PlayerRecycleCount:         easygo.NewInt32(one1.GetPlayerRecycleCount()),         // 用户回收 次数
			OfficialRecycleCount:       easygo.NewInt32(one1.GetOfficialRecycleCount()),       // 平台回收 次数
			RecycleGoldTotal:           easygo.NewInt64(one1.GetRecycleGoldTotal()),           // 回收合计（钻石）
			ProductDiamondTotal:        easygo.NewInt64(one1.GetProductDiamondTotal()),        // 商品价格（钻石）
			LuckyGoldTotal:             easygo.NewInt64(one22.GetLuckyGoldTotal()),            // 抽奖合计（钻石）
			ConvertGoodsPriceTotal:     easygo.NewInt64(one1.GetConvertGoodsPriceTotal()),     // 兑换商品总额（钻石）
			CommissionTotal:            easygo.NewInt32(one4.GetCommissionTotal()),            //官方抽水（分）
			PendConvertGoodsPriceTotal: easygo.NewInt64(one1.GetPendConvertGoodsPriceTotal()), // 待兑换商品总额（钻石）
			PendConvertCount:           easygo.NewInt32(one1.GetPendConvertCount()),           // 待兑换次数
			PlayerRecycleGoldTotal:     easygo.NewInt64(one1.GetPlayerRecycleGoldTotal()),     // 用户回收回收合计（钻石）
			OfficialRecycleGoldTotal:   easygo.NewInt64(one1.GetOfficialRecycleGoldTotal()),   // 平台回收回收合计（钻石）
			Profit:                     easygo.NewInt64(one1.GetProfit()),                     // 盈利（分）
			//ConvertDiamondPlayerCount:  easygo.NewInt32(convertDiamondPlayerCount),            // 兑换钻石人数
			RecyclePlayerCount: easygo.NewInt32(one1.GetRecyclePlayerCount()), // 回收订单人数
			DeliverCount:       easygo.NewInt64(one1.GetDeliverCount()),       // 发货订单数
			DeliverPlayerCount: easygo.NewInt32(one1.GetDeliverPlayerCount()), // 发货订单人数

			CreateTime:        easygo.NewInt64(lastWeekEndTime),
			BoxId:             easygo.NewInt64(item.GetId()),
			BoxName:           easygo.NewString(item.GetName()),
			WinItemPriceTotal: easygo.NewInt64(one1.GetWinItemPriceTotal()), // 出奖钻石金额（钻石）
			StartTime:         easygo.NewInt64(startTime),
			EndTime:           easygo.NewInt64(endTime),
		}
		saveData = append(saveData, bson.M{"_id": wishBoxReportWeek.GetId()}, wishBoxReportWeek)
	}

	//批量生成
	for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_WEEK, saveData)
	logs.Info("完成更新许愿池盲盒报表")

}

// 生成许愿池盲盒报表-月报
func MakeWishBoxReportByMonth() {
	logs.Info("开始更新许愿池盲盒报表-月报")
	dList := GetAllWishBox()
	if len(dList) == 0 {
		return
	}
	now := time.Now()
	lastMonthFirstDay := now.AddDate(0, -1, -now.Day()+1)
	lastMonthStartTime := time.Date(lastMonthFirstDay.Year(), lastMonthFirstDay.Month(), lastMonthFirstDay.Day(), 0, 0, 0, 0, now.Location()).Unix()
	lastMonthEndTime := easygo.GetMonth0ClockTimestamp()

	startTime := lastMonthStartTime
	endTime := lastMonthEndTime
	var saveData []interface{}
	ct := bson.M{"$gte": startTime, "$lt": endTime}
	for _, item := range dList {
		wishBoxId := item.GetId()
		//wishBoxDetailReportWeek := GroupWishBoxReportDetailByWeek(wishBoxId,lastWeekStartTime,lastWeekEndTime)
		wishPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA, "PlayerId", bson.M{"WishBoxId": wishBoxId, "CreateTime": ct}))
		addPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_COLLECTION, "PlayerId", bson.M{"WishBoxId": wishBoxId, "CreateTime": ct}))
		//luckyPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG, "DareId", bson.M{"CreateTime": ct }))
		DrawPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG, "DareId", bson.M{"WishBoxId": wishBoxId, "CreateTime": ct}))
		one1 := GroupWishBoxDetailReport(wishBoxId, startTime, endTime)
		one22 := GroupWishLogByBoxId(wishBoxId, startTime, endTime)
		one4 := GroupWishPoolPumpLog(wishBoxId, startTime, endTime)
		report := &share_message.WishBoxReportMonth{
			Id:                         easygo.NewString(easygo.IntToString(int(startTime)) + easygo.IntToString(int(item.GetId()))),
			WishPlayerCount:            easygo.NewInt32(wishPlayerCount),                      //许愿人数
			AddPlayerCount:             easygo.NewInt32(addPlayerCount),                       //收藏人数
			DrawPlayerCount:            easygo.NewInt32(DrawPlayerCount),                      // 抽奖人数
			DrawCount:                  easygo.NewInt32(one22.GetDrawCount()),                 // 抽奖次数
			LuckyWishCount:             easygo.NewInt32(one1.GetLuckyWishCount()),             // 抽中许愿款次数
			ConvertCount:               easygo.NewInt32(one1.GetConvertCount()),               // 兑换次数
			RecycleCount:               easygo.NewInt32(one1.GetRecycleCount()),               // 回收次数
			PlayerRecycleCount:         easygo.NewInt32(one1.GetPlayerRecycleCount()),         // 用户回收 次数
			OfficialRecycleCount:       easygo.NewInt32(one1.GetOfficialRecycleCount()),       // 平台回收 次数
			RecycleGoldTotal:           easygo.NewInt64(one1.GetRecycleGoldTotal()),           // 回收合计（钻石）
			ProductDiamondTotal:        easygo.NewInt64(one1.GetProductDiamondTotal()),        // 商品价格（钻石）
			LuckyGoldTotal:             easygo.NewInt64(one22.GetLuckyGoldTotal()),            // 抽奖合计（钻石）
			ConvertGoodsPriceTotal:     easygo.NewInt64(one1.GetConvertGoodsPriceTotal()),     // 兑换商品总额（钻石）
			CommissionTotal:            easygo.NewInt32(one4.GetCommissionTotal()),            //官方抽水（分）
			PendConvertGoodsPriceTotal: easygo.NewInt64(one1.GetPendConvertGoodsPriceTotal()), // 待兑换商品总额（钻石）
			PendConvertCount:           easygo.NewInt32(one1.GetPendConvertCount()),           // 待兑换次数
			PlayerRecycleGoldTotal:     easygo.NewInt64(one1.GetPlayerRecycleGoldTotal()),     // 用户回收回收合计（钻石）
			OfficialRecycleGoldTotal:   easygo.NewInt64(one1.GetOfficialRecycleGoldTotal()),   // 平台回收回收合计（钻石）
			Profit:                     easygo.NewInt64(one1.GetProfit()),                     // 盈利（分）
			RecyclePlayerCount:         easygo.NewInt32(one1.GetRecyclePlayerCount()),         // 回收订单人数
			DeliverCount:               easygo.NewInt64(one1.GetDeliverCount()),               // 发货订单数
			DeliverPlayerCount:         easygo.NewInt32(one1.GetDeliverPlayerCount()),         // 发货订单人数

			CreateTime:        easygo.NewInt64(endTime),
			BoxId:             easygo.NewInt64(item.GetId()),
			BoxName:           easygo.NewString(item.GetName()),
			WinItemPriceTotal: easygo.NewInt64(one1.GetWinItemPriceTotal()), // 出奖钻石金额（钻石）
			StartTime:         easygo.NewInt64(startTime),
			EndTime:           easygo.NewInt64(endTime),
		}
		saveData = append(saveData, bson.M{"_id": report.GetId()}, report)
	}
	//logs.Info("更新许愿池盲盒报表-月报saveData:", saveData)
	//批量生成
	for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX_MONTH, saveData)
	logs.Info("完成更新许愿池盲盒报表-月报")

}

// 聚合查询盲盒报表
func GroupWishBoxReport(starTime, endTime int64) *share_message.WishPoolReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_BOX)
	defer closeFun()

	match := bson.M{}
	if endTime > 0 {
		match["CreateTime"] = bson.M{"$gte": starTime, "$lt": endTime}
	} else {
		match["CreateTime"] = starTime
	}
	group := bson.M{
		"_id":                        nil,
		"RecycleCount":               bson.M{"$sum": "$RecycleCount"},
		"PlayerRecycleCount":         bson.M{"$sum": "$PlayerRecycleCount"},
		"OfficialRecycleCount":       bson.M{"$sum": "$OfficialRecycleCount"},
		"WishPlayerCount":            bson.M{"$sum": "$WishPlayerCount"},
		"WishAndDrawPlayerCount":     bson.M{"$sum": "$WishAndDrawPlayerCount"},
		"WinCount":                   bson.M{"$sum": "$WinCount"},
		"LuckyWishCount":             bson.M{"$sum": "$LuckyWishCount"},
		"ConvertCount":               bson.M{"$sum": "$ConvertCount"},
		"PendConvertCount":           bson.M{"$sum": "$PendConvertCount"},
		"PlayerRecycleGoldTotal":     bson.M{"$sum": "$PlayerRecycleGoldTotal"},
		"OfficialRecycleGoldTotal":   bson.M{"$sum": "$OfficialRecycleGoldTotal"},
		"ConvertGoodsPriceTotal":     bson.M{"$sum": "$ConvertGoodsPriceTotal"},
		"PendConvertGoodsPriceTotal": bson.M{"$sum": "$PendConvertGoodsPriceTotal"},
		"Profit":                     bson.M{"$sum": "$Profit"},
		"LuckyCount":                 bson.M{"$sum": "$DrawCount"},
		"LuckyPlayerCount":           bson.M{"$sum": "$DrawPlayerCount"},
		"LuckyGoldTotal":             bson.M{"$sum": "$DrawGoldTotal"},
		"AddPlayerCount":             bson.M{"$sum": "$AddPlayerCount"},
		"RecycleGoldTotal":           bson.M{"$sum": "$RecycleGoldTotal"},
		"ProductDiamondTotal":        bson.M{"$sum": "$ProductDiamondTotal"},
		"CommissionTotal":            bson.M{"$sum": "$CommissionTotal"},
		"WinItemPriceTotal":          bson.M{"$sum": "$WinItemPriceTotal"},
	}

	project := bson.M{
		"RecycleCount":               1,
		"PlayerRecycleCount":         1,
		"OfficialRecycleCount":       1,
		"WishPlayerCount":            1,
		"WishAndDrawPlayerCount":     1,
		"WinCount":                   1,
		"LuckyWishCount":             1,
		"ConvertCount":               1,
		"PendConvertCount":           1,
		"PlayerRecycleGoldTotal":     1,
		"OfficialRecycleGoldTotal":   1,
		"ConvertGoodsPriceTotal":     1,
		"PendConvertGoodsPriceTotal": 1,
		"Profit":                     1,
		"LuckyCount":                 1,
		"LuckyPlayerCount":           1,
		"LuckyGoldTotal":             1,
		"AddPlayerCount":             1,
		"RecycleGoldTotal":           1,
		"ProductDiamondTotal":        1,
		"CommissionTotal":            1,
		"WinItemPriceTotal":          1,
	}

	pipeCond := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$project": project},
	}

	var one *share_message.WishPoolReport
	err := col.Pipe(pipeCond).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return one
}

func InitWishPoolReport() {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_POOL)
	defer closeFun()

	queryBson := bson.M{}
	_, err := col.RemoveAll(queryBson)
	easygo.PanicError(err)
}

//更新许愿池报表
func UpdateWishPoolReport(report *share_message.WishPoolReport) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_POOL)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetCreateTime()}, bson.M{"$set": report})
	easygo.PanicError(err)
}

//查询许愿池报表
func GetWishPoolReport(startTime int64) *share_message.WishPoolReport {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_POOL)
	defer closeFun()
	report := &share_message.WishPoolReport{}
	err := col.Find(bson.M{"_id": startTime}).One(&report)

	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return report
}

// 生成许愿池报表
func MakeWishPoolReport() {
	logs.Info("开始更新许愿池报表")
	startTime := int64(0)
	endTime := time.Now().Unix()
	job := GetReportJob(for_game.TABLE_REPORT_WISH_POOL)
	if job != nil {
		startTime = job.GetTime()
	} else {
		InitWishPoolReport()
	}

	report := &share_message.WishPoolReport{}

	dayTime := easygo.Get0ClockTimestamp(startTime)
	if startTime == 0 {
		dayTime = easygo.Get0ClockTimestamp(endTime)
		startTime = dayTime
	}
	report = GetWishPoolReport(dayTime)
	if report == nil {
		report = &share_message.WishPoolReport{
			CreateTime: easygo.NewInt64(dayTime),
		}
	}

	one1 := GroupWishBoxReport(dayTime, 0)
	report.LuckyPlayerCount = easygo.NewInt32(one1.GetLuckyPlayerCount())
	report.LuckyCount = easygo.NewInt32(one1.GetLuckyCount())
	report.LuckyGoldTotal = easygo.NewInt64(one1.GetLuckyGoldTotal())
	report.AddPlayerCount = easygo.NewInt32(one1.GetAddPlayerCount())

	report.RecycleCount = easygo.NewInt32(one1.GetRecycleCount())
	report.RecycleGoldTotal = easygo.NewInt64(one1.GetRecycleGoldTotal())
	report.ProductDiamondTotal = easygo.NewInt64(one1.GetProductDiamondTotal())
	report.PlayerRecycleCount = easygo.NewInt32(one1.GetPlayerRecycleCount())
	report.OfficialRecycleCount = easygo.NewInt32(one1.GetOfficialRecycleCount())
	report.PendConvertCount = easygo.NewInt32(one1.GetPendConvertCount())
	report.ConvertCount = easygo.NewInt32(one1.GetConvertCount())
	report.WishPlayerCount = easygo.NewInt32(one1.GetWishPlayerCount())
	report.LuckyWishCount = easygo.NewInt32(one1.GetLuckyWishCount())
	report.PlayerRecycleGoldTotal = easygo.NewInt64(one1.GetPlayerRecycleGoldTotal())
	report.OfficialRecycleGoldTotal = easygo.NewInt64(one1.GetOfficialRecycleGoldTotal())
	report.ConvertGoodsPriceTotal = easygo.NewInt64(one1.GetConvertGoodsPriceTotal())
	report.PendConvertGoodsPriceTotal = easygo.NewInt64(one1.GetPendConvertGoodsPriceTotal())
	report.Profit = easygo.NewInt64(one1.GetProfit())

	report.WinItemPriceTotal = easygo.NewInt64(one1.GetWinItemPriceTotal())

	// 官方抽水
	report.CommissionTotal = easygo.NewInt32(one1.GetCommissionTotal())

	// 新版需求
	// 发货订单
	deliverCount := GetWishDeliveryOrderCount(dayTime)
	report.DeliverCount = easygo.NewInt64(deliverCount)
	one2 := GetWishDeliveryOrderPlayerCount(dayTime)
	report.DeliverPlayerCount = easygo.NewInt32(one2.GetDeliverPlayerCount())
	// 回收订单
	one3 := GetWishRecyclePlayerCount(dayTime)
	report.RecyclePlayerCount = easygo.NewInt32(one3.GetRecyclePlayerCount())
	recycleCount := GetWishRecycleOrderCount(dayTime) // 需求更改，由原来的回收数改为回收订单数，所以这里要改
	report.RecycleCount = easygo.NewInt32(recycleCount)

	// 获取许愿池埋点信息
	wishLog := for_game.GetRedisWishLogReport(dayTime)
	wishReport := wishLog.GetRedisWishLogReport()
	report.InPoolCount = easygo.NewInt64(wishReport.GetWishTime())
	report.InPoolPlayerCount = easygo.NewInt32(wishReport.GetNewPlayer() + wishReport.GetOldPlayer())
	// 兑换钻石
	/*one4 := GetWishConvertDiamondPlayerCount(dayTime)
	report.ConvertDiamondPlayerCount = easygo.NewInt32(one4.GetConvertDiamondPlayerCount())*/
	diamondExchangePeopleBson := bson.M{"PayType": 1, "SourceType": for_game.DIAMOND_TYPE_EXCHANGE_IN, "RecycleTime": dayTime}
	convertDiamondPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_DIAMOND_CHANGELOG, "PlayerId", diamondExchangePeopleBson))
	report.ConvertDiamondPlayerCount = easygo.NewInt32(convertDiamondPlayerCount)
	one5 := GetWishConvertDiamond(dayTime)
	report.ConvertDiamondTotal = easygo.NewInt64(one5.GetRecyclePlayerCount())
	report.ConvertDiamondCount = easygo.NewInt64(one3.GetConvertDiamondCount())

	UpdateWishPoolReport(report)
	logs.Info("完成更新许愿池报表")
	//写报表生成进度
	MakeReportJob(for_game.TABLE_REPORT_WISH_POOL, endTime)
}

// 生成许愿池报表-周表
func MakeWishPoolReportByWeek() {
	logs.Info("开始更新许愿池报表-周表")
	lastWeekEndTime := easygo.GetWeek0ClockOfTimestamp(time.Now().Unix())
	lastWeekStartTime := lastWeekEndTime - 3600*24*7
	startTime := lastWeekStartTime
	endTime := lastWeekEndTime
	ctBson := bson.M{"$gte": startTime, "$lt": endTime}
	one1 := GroupWishBoxReport(startTime, endTime)
	wishPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA, "PlayerId", bson.M{"CreateTime": ctBson}))
	addPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_COLLECTION, "PlayerId", bson.M{"CreateTime": ctBson}))
	luckyPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG, "DareId", bson.M{"CreateTime": ctBson}))
	recyclePlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER, "PlayerId", bson.M{"RecycleTime": ctBson}))
	InPoolCount := CountWishTime(startTime, endTime)
	inPoolPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BURYING_POINT_LOG, "PlayerId", bson.M{"Time": ctBson}))
	convertDiamondTotal := CountExchangeDiamond(startTime, endTime)
	exchangeDiamondPeople := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_DIAMOND_CHANGELOG, "PlayerId", bson.M{"SourceType": for_game.DIAMOND_TYPE_EXCHANGE_IN, "PayType": 1, "CreateTime": ctBson}))
	convertDiamondCount := for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_DIAMOND_CHANGELOG, bson.M{"SourceType": for_game.DIAMOND_TYPE_EXCHANGE_IN})
	deliverCount := for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG, bson.M{"Status": 1, "CreateTime": ctBson})
	deliverPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG, "PlayerId", bson.M{"Status": 1, "CreateTime": bson.M{"$gte": lastWeekStartTime, "$lt": lastWeekStartTime}}))
	report := &share_message.WishPoolReportWeek{
		CreateTime:                 easygo.NewInt64(lastWeekStartTime),
		WishPlayerCount:            easygo.NewInt32(wishPlayerCount),          //许愿人数
		AddPlayerCount:             easygo.NewInt32(addPlayerCount),           //收藏人数
		LuckyPlayerCount:           easygo.NewInt32(luckyPlayerCount),         //抽奖人数
		LuckyCount:                 easygo.NewInt32(one1.GetLuckyCount()),     //抽奖次数
		LuckyWishCount:             easygo.NewInt32(one1.GetLuckyWishCount()), //抽中许愿款次数
		ConvertCount:               easygo.NewInt32(one1.GetConvertCount()),   // 兑换次数
		RecycleCount:               easygo.NewInt32(one1.GetRecycleCount()),
		PlayerRecycleCount:         easygo.NewInt32(one1.GetPlayerRecycleCount()),
		OfficialRecycleCount:       easygo.NewInt32(one1.GetOfficialRecycleCount()),
		RecycleGoldTotal:           easygo.NewInt64(one1.GetRecycleGoldTotal()),
		ProductDiamondTotal:        easygo.NewInt64(one1.GetProductDiamondTotal()),
		LuckyGoldTotal:             easygo.NewInt64(one1.GetLuckyGoldTotal()),
		ConvertGoodsPriceTotal:     easygo.NewInt64(one1.GetConvertGoodsPriceTotal()),
		CommissionTotal:            easygo.NewInt32(one1.GetCommissionTotal()),
		PendConvertGoodsPriceTotal: easygo.NewInt64(one1.GetPendConvertGoodsPriceTotal()),
		PendConvertCount:           easygo.NewInt32(one1.GetPendConvertCount()),
		PlayerRecycleGoldTotal:     easygo.NewInt64(one1.GetPlayerRecycleGoldTotal()),
		OfficialRecycleGoldTotal:   easygo.NewInt64(one1.GetOfficialRecycleGoldTotal()),
		Profit:                     easygo.NewInt64(one1.GetProfit()),

		ConvertDiamondTotal:       easygo.NewInt64(convertDiamondTotal),   // 兑换钻石总金额
		ConvertDiamondPlayerCount: easygo.NewInt32(exchangeDiamondPeople), //兑换钻石人数
		ConvertDiamondCount:       easygo.NewInt64(convertDiamondCount),   // 兑换钻石次数

		RecyclePlayerCount: easygo.NewInt32(recyclePlayerCount), //回收订单人数
		DeliverCount:       easygo.NewInt64(deliverCount),       // 发货订单数
		DeliverPlayerCount: easygo.NewInt32(deliverPlayerCount), // 发货订单人数
		InPoolCount:        easygo.NewInt64(InPoolCount),        // 进入许愿池次数
		InPoolPlayerCount:  easygo.NewInt32(inPoolPlayerCount),  // 进入许愿池人数
		//WinItemPriceTotal: easygo.NewInt64(0), // 出奖钻石金额（钻石）前端计算
		StartTime: easygo.NewInt64(lastWeekStartTime),
		EndTime:   easygo.NewInt64(lastWeekEndTime),
	}
	//logs.Info("更新许愿池报表-周报: report", report)
	for_game.InsertMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_POOL_WEEK, report)
	//UpdateWishPoolReport(report)
	logs.Info("完成更新许愿池报表-周报")
}

// 生成许愿池报表-月表
func MakeWishPoolReportByMonth() {
	logs.Info("开始更新许愿池报表-月报")
	now := time.Now()
	lastMonthFirstDay := now.AddDate(0, -1, -now.Day()+1)
	lastMonthStartTime := time.Date(lastMonthFirstDay.Year(), lastMonthFirstDay.Month(), lastMonthFirstDay.Day(), 0, 0, 0, 0, now.Location()).Unix()
	lastMonthEndTime := easygo.GetMonth0ClockTimestamp()

	startTime := lastMonthStartTime
	endTime := lastMonthEndTime

	ctBson := bson.M{"$gte": startTime, "$lt": endTime}
	one1 := GroupWishBoxReport(startTime, endTime)
	wishPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA, "PlayerId", bson.M{"CreateTime": ctBson}))
	addPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_COLLECTION, "PlayerId", bson.M{"CreateTime": ctBson}))
	luckyPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_LOG, "DareId", bson.M{"CreateTime": ctBson}))
	recyclePlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_RECYCLE_ORDER, "PlayerId", bson.M{"RecycleTime": ctBson}))
	InPoolCount := CountWishTime(startTime, endTime)
	inPoolPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BURYING_POINT_LOG, "PlayerId", bson.M{"Time": ctBson}))
	convertDiamondTotal := CountExchangeDiamond(startTime, endTime)
	exchangeDiamondPeople := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_DIAMOND_CHANGELOG, "PlayerId", bson.M{"SourceType": for_game.DIAMOND_TYPE_EXCHANGE_IN, "PayType": 1, "CreateTime": ctBson}))
	convertDiamondCount := for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_DIAMOND_CHANGELOG, bson.M{"SourceType": for_game.DIAMOND_TYPE_EXCHANGE_IN})
	deliverCount := for_game.FindAllCount(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG, bson.M{"Status": 1, "CreateTime": ctBson})
	deliverPlayerCount := len(for_game.Distinct(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG, "PlayerId", bson.M{"Status": 1, "CreateTime": bson.M{"$gte": startTime, "$lt": endTime}}))
	report := &share_message.WishPoolReportMonth{
		CreateTime:                 easygo.NewInt64(startTime),
		WishPlayerCount:            easygo.NewInt32(wishPlayerCount),          //许愿人数
		AddPlayerCount:             easygo.NewInt32(addPlayerCount),           //收藏人数
		LuckyPlayerCount:           easygo.NewInt32(luckyPlayerCount),         //抽奖人数
		LuckyCount:                 easygo.NewInt32(one1.GetLuckyCount()),     //抽奖次数
		LuckyWishCount:             easygo.NewInt32(one1.GetLuckyWishCount()), //抽中许愿款次数
		ConvertCount:               easygo.NewInt32(one1.GetConvertCount()),   // 兑换次数
		RecycleCount:               easygo.NewInt32(one1.GetRecycleCount()),
		PlayerRecycleCount:         easygo.NewInt32(one1.GetPlayerRecycleCount()),
		OfficialRecycleCount:       easygo.NewInt32(one1.GetOfficialRecycleCount()),
		RecycleGoldTotal:           easygo.NewInt64(one1.GetRecycleGoldTotal()),
		ProductDiamondTotal:        easygo.NewInt64(one1.GetProductDiamondTotal()),
		LuckyGoldTotal:             easygo.NewInt64(one1.GetLuckyGoldTotal()),
		ConvertGoodsPriceTotal:     easygo.NewInt64(one1.GetConvertGoodsPriceTotal()),
		CommissionTotal:            easygo.NewInt32(one1.GetCommissionTotal()),
		PendConvertGoodsPriceTotal: easygo.NewInt64(one1.GetPendConvertGoodsPriceTotal()),
		PendConvertCount:           easygo.NewInt32(one1.GetPendConvertCount()),
		PlayerRecycleGoldTotal:     easygo.NewInt64(one1.GetPlayerRecycleGoldTotal()),
		OfficialRecycleGoldTotal:   easygo.NewInt64(one1.GetOfficialRecycleGoldTotal()),
		Profit:                     easygo.NewInt64(one1.GetProfit()),

		ConvertDiamondTotal:       easygo.NewInt64(convertDiamondTotal),   // 兑换钻石总金额
		ConvertDiamondPlayerCount: easygo.NewInt32(exchangeDiamondPeople), //兑换钻石人数
		ConvertDiamondCount:       easygo.NewInt64(convertDiamondCount),   // 兑换钻石次数

		RecyclePlayerCount: easygo.NewInt32(recyclePlayerCount), //回收订单人数
		DeliverCount:       easygo.NewInt64(deliverCount),       // 发货订单数
		DeliverPlayerCount: easygo.NewInt32(deliverPlayerCount), // 发货订单人数
		InPoolCount:        easygo.NewInt64(InPoolCount),        // 进入许愿池次数
		InPoolPlayerCount:  easygo.NewInt32(inPoolPlayerCount),  // 进入许愿池人数
		//WinItemPriceTotal: easygo.NewInt64(0), // 出奖钻石金额（钻石）前端计算
		StartTime: easygo.NewInt64(startTime),
		EndTime:   easygo.NewInt64(endTime),
	}
	//logs.Info("更新许愿池报表-月报：report", report)
	for_game.InsertMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_REPORT_WISH_POOL_MONTH, report)
	//UpdateWishPoolReport(report)
	logs.Info("完成更新许愿池报表-月报")
}

// +++++++++++++++++++++++++++++许愿池
// 获取水池当前的状态
func GetPoolStatus(pool *share_message.WishPool) int32 {
	f := easygo.Decimal(easygo.AtoFloat64(easygo.AnytoA(pool.GetIncomeValue()))/easygo.AtoFloat64(easygo.AnytoA(pool.GetPoolLimit())), 2)
	poolStatusLimitFloat64 := f * 100
	//poolStatusLimit := easygo.AtoInt64(easygo.AnytoA(poolStatusLimitFloat64))
	poolStatusLimit := int64(easygo.Decimal(poolStatusLimitFloat64, 0))
	logs.Info("poolStatusLimit---------->", poolStatusLimit)
	var poolStatus int32 // 1-大亏;2-小亏;3-普通;4-大盈;5-小盈
	bigLoss := pool.GetBigLoss()
	if bigLoss != nil && (poolStatusLimit <= bigLoss.GetShowMaxValue() && poolStatusLimit >= bigLoss.GetShowMinValue()) {
		poolStatus = for_game.POOL_STATUS_BIGLOSS
	}
	smallLoss := pool.GetSmallLoss()
	if smallLoss != nil && (poolStatusLimit <= smallLoss.GetShowMaxValue() && poolStatusLimit >= smallLoss.GetShowMinValue()) {
		poolStatus = for_game.POOL_STATUS_SMALLLOSS
	}
	common := pool.GetCommon()
	if common != nil && (poolStatusLimit <= common.GetShowMaxValue() && poolStatusLimit >= common.GetShowMinValue()) {
		poolStatus = for_game.POOL_STATUS_COMMON
	}
	bigWin := pool.GetBigWin()
	if bigWin != nil && (poolStatusLimit <= bigWin.GetShowMaxValue() && poolStatusLimit >= bigWin.GetShowMinValue()) {
		poolStatus = for_game.POOL_STATUS_BIGWIN
	}
	smallWin := pool.GetSmallWin()
	if smallWin != nil && (poolStatusLimit <= smallWin.GetShowMaxValue() && poolStatusLimit >= smallWin.GetShowMinValue()) {
		poolStatus = for_game.POOL_STATUS_SMALLWIN
	}
	if bigWin != nil && poolStatusLimit >= bigWin.GetShowMaxValue() { // 大于100%
		poolStatus = for_game.POOL_STATUS_BIGWIN
	}
	return poolStatus
}

// ++++++++++++++++++盲盒
//盲盒挑战占领时长表
func UpdateWishOccupied(report *share_message.WishOccupied) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OCCUPIED)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": report.GetId()}, bson.M{"$set": report})
	easygo.PanicError(err)

}

func GetWishOccupiedByBoxId(boxId, playerId int64) *share_message.WishOccupied {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_OCCUPIED)
	defer closeFun()

	one := &share_message.WishOccupied{}
	err := col.Find(bson.M{"WishBoxId": boxId, "PlayerId": playerId}).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}

	if err == mgo.ErrNotFound {
		return one
	}

	return one
}

// 删除盲盒商品相关用户许愿数据
func DeleteWishPlayerData(boxId int64, itemId []int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_DATA)
	defer closeFun()

	_, err := col.RemoveAll(bson.M{"WishBoxId": boxId, "WishBoxItemId": bson.M{"$in": itemId}})
	easygo.PanicError(err)
}

// 获取用户物品信息
func GetPlayerWishItemByOrder(itemId, pId int64) *share_message.PlayerWishItem {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	one := &share_message.PlayerWishItem{}
	err := col.Find(bson.M{"ChallengeItemId": itemId, "PlayerId": pId}).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}

	if err == mgo.ErrNotFound {
		return one
	}

	return one
}

// 更新用户物品信息状态
func UpdatePlayerWishItemStatusByOrder(ids []int64, status int32, Operator ...string) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	upBson := bson.M{"Status": status, "UpdateTime": easygo.NowTimestamp()}
	if status == 4 && len(Operator) > 0 {
		upBson["Operator"] = Operator[0]
	}

	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": ids}, "Status": 3}, bson.M{"$set": upBson})
	easygo.PanicError(err)
}

// 更新用户物品信息状态
func UpdatePlayerWishItemStatusByPlayerExchangeLog(ids []int64, status int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_WISH_ITEM)
	defer closeFun()

	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": bson.M{"Status": status, "UpdateTime": time.Now().Unix()}})
	easygo.PanicError(err)
}

// 更新用户物品信息状态
func UpdatePlayerExchangeLog(ids []int64, status int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_EXCHANGE_LOG)
	defer closeFun()

	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": ids}, "Status": 3}, bson.M{"$set": bson.M{"Status": status}})
	easygo.PanicError(err)
}

// 获取钻石流水日志
func QueryDiamondChangeLog(reqMsg *brower_backstage.ListRequest) ([]*share_message.DiamondChangeLog, int32) {
	pageSize, curPage := SetMgoPage(reqMsg.GetPageSize(), reqMsg.GetCurPage())
	findBson := bson.M{}
	sort := []string{"-_id"}
	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp() * 1000, "$lt": reqMsg.GetEndTimestamp() * 1000}
	}
	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["PayType"] = reqMsg.GetListType()
	}
	if reqMsg.DownType != nil && reqMsg.GetDownType() != 0 {
		findBson["SourceType"] = reqMsg.GetDownType()
	}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			id := easygo.StringToInt64noErr(reqMsg.GetKeyword())
			findBson["_id"] = id
		case 2:
			base := QueryPlayerbyAccount(reqMsg.GetKeyword())
			if base != nil {
				wishP := GetWishPlayerInfoByPid(base.GetPlayerId())
				findBson["PlayerId"] = wishP.GetId()
			}
		}
	}
	if reqMsg.UserType != nil {
		oPids := GetOperateIds()
		switch reqMsg.GetUserType() {
		case 1:
			findBson["Extend.PlayerId"] = bson.M{"$nin": oPids}
		case 2:
			findBson["Extend.PlayerId"] = bson.M{"$in": oPids}
		}
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_DIAMOND_CHANGELOG)
	defer closeFun()

	query := col.Find(findBson)
	count, err := query.Count()
	easygo.PanicError(err)

	var list []*share_message.DiamondChangeLog
	errc := query.Sort(sort...).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, int32(count)
}
