package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
)

const (
	SHOP_ITEMLIST_TYPE_1 int32 = 1 //1商品id(默认)
	SHOP_ITEMLIST_TYPE_2 int32 = 2 //2商品名称
	SHOP_ITEMLIST_TYPE_3 int32 = 3 //3卖家柠檬号

	SHOP_ITEMLIST_STATUS_DEFAULT int32 = 1000 //1000全部(默认)

	SHOP_ORDER_TIME_TYPE_1 int32 = 1 //1创建时间(默认)
	SHOP_ORDER_TIME_TYPE_2 int32 = 2 //2付款时间
	SHOP_ORDER_TIME_TYPE_3 int32 = 3 //3发货时间
	SHOP_ORDER_TIME_TYPE_4 int32 = 4 //4完成时间

	SHOP_ORDER_STATUS_CANCEL  int32 = 8    //8取消(包括1超时,2取消,7后台取消)
	SHOP_ORDER_STATUS_DEFAULT int32 = 1000 //1000全部(默认)

	SHOP_ORDER_TYPE_1 int32 = 1 //1订单id(默认)
	SHOP_ORDER_TYPE_2 int32 = 2 //2商品id
	SHOP_ORDER_TYPE_3 int32 = 3 //3卖家柠檬号
	SHOP_ORDER_TYPE_4 int32 = 4 //4买家柠檬号
	SHOP_ORDER_TYPE_5 int32 = 5 //卖家电话或邮箱

	SHOP_COMMENT_TYPE_DEFAULT int32 = 1000 //1000全部(默认)

	SHOP_POINT_CARD_STATUS_DEFAULT int32 = 1000 //1000全部(默认)

	SHOP_POINT_CARD_TYPE_1 int32 = 1 //1点卡id
	SHOP_POINT_CARD_TYPE_2 int32 = 2 //2点卡名称
	SHOP_POINT_CARD_TYPE_3 int32 = 3 //3卡号
	SHOP_POINT_CARD_TYPE_4 int32 = 4 //4卖家柠檬帐号
	SHOP_POINT_CARD_TYPE_5 int32 = 5 //5订单号
)

//查询商城商品列表
func QueryShopItemList(reqMsg *brower_backstage.QueryShopItemRequest) ([]*brower_backstage.QueryShopItemObject, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 &&
		reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
		queryBson["create_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	} else if (reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0) &&
		(reqMsg.EndTimestamp == nil || reqMsg.GetEndTimestamp() == 0) {
		queryBson["create_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp()}
	} else if (reqMsg.BeginTimestamp == nil || reqMsg.GetBeginTimestamp() == 0) &&
		(reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0) {
		queryBson["create_time"] = bson.M{"$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() != SHOP_ITEMLIST_STATUS_DEFAULT {
		queryBson["state"] = reqMsg.GetStatus()
	} else if reqMsg.Status != nil && reqMsg.GetStatus() == SHOP_ITEMLIST_STATUS_DEFAULT {
		queryBson["$or"] = []bson.M{
			{"state": for_game.SHOP_ITEM_SALE},
			{"state": for_game.SHOP_ITEM_SOLD_OUT}}
	}

	if reqMsg.GetKeyword() != "" && reqMsg.Keyword != nil {
		switch reqMsg.GetTypes() {
		case SHOP_ITEMLIST_TYPE_1:
			queryBson["_id"] = easygo.StringToInt64noErr(reqMsg.GetKeyword())
		case SHOP_ITEMLIST_TYPE_2:
			queryBson["name"] = reqMsg.GetKeyword()
		case SHOP_ITEMLIST_TYPE_3:
			queryBson["player_account"] = reqMsg.GetKeyword()
		}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	sort := "-_id"

	var list []*share_message.TableShopItem
	errc := query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	if errc != nil && errc != mgo.ErrNotFound {
		easygo.PanicError(errc)
	}

	// responseList := make([]*brower_backstage.QueryShopItemObject, pageSize)
	responseList := []*brower_backstage.QueryShopItemObject{}
	//组装成页面显示的数据
	if nil != list && len(list) > 0 {
		for _, value := range list {

			itemTypeName := value.GetType().GetOtherType()[0]

			if itemTypeName == "" {
				itemTypeName = "--"
			}

			response := &brower_backstage.QueryShopItemObject{
				ItemId:        easygo.NewInt64(value.GetItemId()),
				Name:          easygo.NewString(value.GetName()),
				ItemTypeName:  easygo.NewString(itemTypeName),
				Price:         easygo.NewInt32(value.GetPrice()),
				StockCount:    easygo.NewInt32(value.GetStockCount()),
				PlayerAccount: easygo.NewString(value.GetPlayerAccount()),
				State:         easygo.NewInt32(value.GetState()),
				CreateTime:    easygo.NewInt64(value.GetCreateTime()),
			}
			responseList = append(responseList, response)
		}
	}

	return responseList, count
}

//id查询商品详情商品详情页面用
func QueryShopItemDetailById(Id SHOP_ITEM_ID) *brower_backstage.QueryShopItemDetailResponse {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()
	shopItem := &share_message.TableShopItem{}

	err := col.Find(bson.M{"_id": Id}).One(shopItem)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	rp := &brower_backstage.QueryShopItemDetailResponse{}
	if nil != shopItem {

		rp.ItemId = easygo.NewInt64(shopItem.GetItemId())
		rp.Name = easygo.NewString(shopItem.GetName())

		//组装图片URL
		var itemFiles = []*brower_backstage.ShopItemFile{}

		var mItemFiles []*share_message.ItemFile = shopItem.GetItemFiles()
		if nil != mItemFiles && len(mItemFiles) > 0 {
			for i := 0; i < len(mItemFiles); i++ {
				var itemFile *brower_backstage.ShopItemFile = &brower_backstage.ShopItemFile{
					FileUrl:    easygo.NewString(mItemFiles[i].GetFileUrl()),
					FileType:   easygo.NewInt32(mItemFiles[i].GetFileType()),
					FileWidth:  easygo.NewString(mItemFiles[i].GetFileWidth()),
					FileHeight: easygo.NewString(mItemFiles[i].GetFileHeight()),
				}

				itemFiles = append(itemFiles, itemFile)
			}
		}
		rp.ItemFiles = itemFiles
		rp.PlayerAccount = easygo.NewString(shopItem.GetPlayerAccount())
		rp.Price = easygo.NewInt32(shopItem.GetPrice())
		//设置商品的分类 45为点卡
		if nil != shopItem.GetType() {
			rp.ItemType = easygo.NewInt32(shopItem.GetType().GetType())
		}

		//分割品类标签和常用标签
		itemTypeName := "--"
		categoryLabel := ""
		commonUseLabel := ""
		if nil != shopItem.GetType() &&
			shopItem.GetType().GetOtherType() != nil &&
			len(shopItem.GetType().GetOtherType()) > 1 {

			itemTypeName = shopItem.GetType().GetOtherType()[0]
			if itemTypeName == "" {
				itemTypeName = "--"
			}

			categoryLabel = shopItem.GetType().GetOtherType()[1]

			if len(shopItem.GetType().GetOtherType()) > 2 {
				for i := 2; i < len(shopItem.GetType().GetOtherType()); i++ {
					tempLabel := shopItem.GetType().GetOtherType()[i]
					if i == 2 {
						commonUseLabel = commonUseLabel + tempLabel
					} else {
						commonUseLabel = commonUseLabel + "/" + tempLabel
					}
				}
			}

		}
		rp.ItemTypeName = easygo.NewString(itemTypeName)
		rp.CategoryLabel = easygo.NewString(categoryLabel)
		rp.CommonUseLabel = easygo.NewString(commonUseLabel)

		//取得好评率
		var goodCommentRate int32
		//设置好评率
		//固定好评率设置不为0的时候 显示固定好评率
		if shopItem.GetFakeFixGoodCommRate() > 0 {
			goodCommentRate = shopItem.GetFakeFixGoodCommRate()

			//假的好评数和假的评价总数的判断
			//假的好评和假的好评完成数目前没有地方设置 这里备用
		} else if shopItem.GetFakeGoodCommCnt() > 0 && shopItem.GetFakeFinCommCnt() > 0 {
			goodCommentRate = shopItem.GetFakeGoodCommCnt() * 100 / shopItem.GetFakeFinCommCnt()
		} else if shopItem.GetRealGoodCommCnt() > 0 && shopItem.GetRealFinCommCnt() > 0 {
			goodCommentRate = shopItem.GetRealGoodCommCnt() * 100 / shopItem.GetRealFinCommCnt()
		} else {
			goodCommentRate = 0
		}

		rp.GoodCommentRate = easygo.NewInt32(goodCommentRate)
		rp.StockCount = easygo.NewInt32(shopItem.GetStockCount())

		//付款数
		var payCnt int32
		if shopItem.GetFakePayCnt() > 0 {
			payCnt = shopItem.GetFakePayCnt()
		} else {
			//不实时从订单表中取得数据，从商品表冗余取得数据
			payCnt = shopItem.GetRealPayCnt()
		}
		rp.PaymentCount = easygo.NewInt32(payCnt)

		//浏览数
		var pageViews int32
		if shopItem.GetFakePageViews() > 0 {
			pageViews = shopItem.GetFakePageViews()
		} else {
			pageViews = shopItem.GetRealPageViews()
		}
		rp.PageViews = easygo.NewInt32(pageViews)

		//卖了多少件宝贝
		var sellItemCount int32
		var fakePlayFinOrderCnt int32
		var shopPlayer = share_message.TableShopPlayer{}

		colPlay, closeFuPlay := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_PLAYER)
		defer closeFuPlay()
		errPlay := colPlay.Find(bson.M{"_id": shopItem.GetPlayerId()}).Limit(1).One(&shopPlayer)

		if nil == errPlay {
			if shopPlayer.GetFakePlayFinOrderCnt() > 0 {
				fakePlayFinOrderCnt = shopPlayer.GetFakePlayFinOrderCnt()
			} else {
				fakePlayFinOrderCnt = 0
			}
		} else {
			fakePlayFinOrderCnt = 0
		}

		if fakePlayFinOrderCnt > 0 {
			sellItemCount = fakePlayFinOrderCnt
		} else {
			//计算该用户在平台一共卖出的宝贝数

			col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
			defer closeFun()
			query := col.Pipe([]bson.M{
				{"$match": bson.M{"$and": []bson.M{{"sponsor_id": shopItem.GetPlayerId()},
					{"$or": []bson.M{{"state": for_game.SHOP_ORDER_FINISH},
						{"state": for_game.SHOP_ORDER_EVALUTE}}}}}},
				{"$group": bson.M{"_id": "$sponsor_id", "total": bson.M{"$sum": "$items.count"}}}})

			rst := make([]bson.M, 0)
			e := query.All(&rst)
			var sum int = 0
			if nil == e {
				c := len(rst)
				if rst != nil && c > 0 {
					sum = (rst[0]["total"]).(int)
				}
			} else {
				sum = 0
			}
			sellItemCount = int32(sum)
		}
		rp.SellItemCount = easygo.NewInt32(sellItemCount)
		rp.CreateTime = easygo.NewInt64(shopItem.GetCreateTime())
		rp.SoldOutTime = easygo.NewInt64(shopItem.GetSoldOutTime())
		rp.State = easygo.NewInt32(shopItem.GetState())
		rp.Title = easygo.NewString(shopItem.GetTitle())
		//点卡名称
		rp.PointCardName = easygo.NewString(shopItem.GetPointCardName())
	} else {
		return nil
	}
	return rp
}

//列表页下架按钮
func ShopSoldOut(Id SHOP_ITEM_ID) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()
	var nowTime int64 = time.Now().Unix()
	err := col.Update(
		bson.M{"_id": Id},
		bson.M{"$set": bson.M{"state": for_game.SHOP_ITEM_SOLD_OUT,
			"sold_out_time": nowTime,
		}})

	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
}

//发布商品
func ReleaseShopItem(reqMsg *share_message.TableShopItem) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()
	err := col.Insert(reqMsg)
	easygo.PanicError(err)
}

//修改跳转取得的数据
func GetEditShopItemDetailById(Id SHOP_ITEM_ID) *brower_backstage.ReleaseEditShopItemObject {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()
	shopItem := &share_message.TableShopItem{}

	err := col.Find(bson.M{"_id": Id}).One(shopItem)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	rp := &brower_backstage.ReleaseEditShopItemObject{}
	if nil != shopItem {

		rp.ItemId = easygo.NewInt64(shopItem.GetItemId())
		rp.Name = easygo.NewString(shopItem.GetName())

		//组装图片URL
		var itemFiles = []*brower_backstage.ShopItemFile{}

		var mItemFiles []*share_message.ItemFile = shopItem.GetItemFiles()
		if nil != mItemFiles && len(mItemFiles) > 0 {
			for i := 0; i < len(mItemFiles); i++ {
				var itemFile *brower_backstage.ShopItemFile = &brower_backstage.ShopItemFile{
					FileUrl:    easygo.NewString(mItemFiles[i].GetFileUrl()),
					FileType:   easygo.NewInt32(mItemFiles[i].GetFileType()),
					FileWidth:  easygo.NewString(mItemFiles[i].GetFileWidth()),
					FileHeight: easygo.NewString(mItemFiles[i].GetFileHeight()),
				}

				itemFiles = append(itemFiles, itemFile)
			}
		}
		rp.ItemFiles = itemFiles
		rp.PlayerAccount = easygo.NewString(shopItem.GetPlayerAccount())

		//设置商品分类
		if nil != shopItem.GetType() {
			rp.ItemType = easygo.NewInt32(shopItem.GetType().GetType())
		}

		//取得品类标签和常用标签
		categoryLabel := "请选择"
		commonUseLabel := []string{}
		if nil != shopItem.GetType() &&
			shopItem.GetType().GetOtherType() != nil &&
			len(shopItem.GetType().GetOtherType()) > 1 {

			categoryLabel = shopItem.GetType().GetOtherType()[1]
			if len(shopItem.GetType().GetOtherType()) > 2 {
				for i := 2; i < len(shopItem.GetType().GetOtherType()); i++ {
					commonUseLabel = append(commonUseLabel, shopItem.GetType().GetOtherType()[i])

				}
			}

		}

		rp.ItemCategory = easygo.NewString(categoryLabel)
		rp.CommonUseLabel = commonUseLabel
		rp.Price = easygo.NewInt32(shopItem.GetPrice())
		rp.StockCount = easygo.NewInt32(shopItem.GetStockCount())
		rp.UserName = easygo.NewString(shopItem.GetUserName())
		rp.Phone = easygo.NewString(shopItem.GetPhone())
		rp.Address = easygo.NewString(shopItem.GetAddress())
		rp.DetailAddress = easygo.NewString(shopItem.GetDetailAddress())
		rp.FakePaymentCount = easygo.NewInt32(shopItem.GetFakePayCnt())
		rp.RealPaymentCount = easygo.NewInt32(shopItem.GetRealPayCnt())
		rp.FakePageViews = easygo.NewInt32(shopItem.GetFakePageViews())
		rp.RealPageViews = easygo.NewInt32(shopItem.GetRealPageViews())

		//卖了多少件宝贝
		var sellItemCount int32
		var fakePlayFinOrderCnt int32
		var shopPlayer = share_message.TableShopPlayer{}

		colPlay, closeFuPlay := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_PLAYER)
		defer closeFuPlay()
		errPlay := colPlay.Find(bson.M{"_id": shopItem.GetPlayerId()}).Limit(1).One(&shopPlayer)

		if nil == errPlay {
			if shopPlayer.GetFakePlayFinOrderCnt() > 0 {
				fakePlayFinOrderCnt = shopPlayer.GetFakePlayFinOrderCnt()
			} else {
				fakePlayFinOrderCnt = 0
			}
		} else {
			fakePlayFinOrderCnt = 0
		}

		rp.FakeSellItemCount = easygo.NewInt32(fakePlayFinOrderCnt)

		//计算该用户在平台一共卖出的实际的宝贝数
		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
		defer closeFun()
		query := col.Pipe([]bson.M{
			{"$match": bson.M{"$and": []bson.M{{"sponsor_id": shopItem.GetPlayerId()},
				{"$or": []bson.M{{"state": for_game.SHOP_ORDER_FINISH},
					{"state": for_game.SHOP_ORDER_EVALUTE}}}}}},
			{"$group": bson.M{"_id": "$sponsor_id", "total": bson.M{"$sum": "$items.count"}}}})

		rst := make([]bson.M, 0)
		e := query.All(&rst)
		var sum int = 0
		if nil == e {
			c := len(rst)
			if rst != nil && c > 0 {
				sum = (rst[0]["total"]).(int)
			}
		} else {
			sum = 0
		}
		sellItemCount = int32(sum)
		rp.RealSellItemCount = easygo.NewInt32(sellItemCount)

		rp.FakeGoodCommentRate = easygo.NewInt32(shopItem.GetFakeFixGoodCommRate())
		//实际的好评
		var goodCommentRate int32
		if shopItem.GetRealGoodCommCnt() > 0 && shopItem.GetRealFinCommCnt() > 0 {
			goodCommentRate = shopItem.GetRealGoodCommCnt() * 100 / shopItem.GetRealFinCommCnt()
		} else {
			goodCommentRate = 0
		}
		rp.RealGoodCommentRate = easygo.NewInt32(goodCommentRate)
		rp.Title = easygo.NewString(shopItem.GetTitle())
		rp.State = easygo.NewInt32(shopItem.GetState())

		//设置点卡名称
		rp.PointCardName = easygo.NewString(shopItem.GetPointCardName())

	} else {
		return nil
	}
	return rp
}

//修改商品详情
func EditShopItem(reqMsg *share_message.TableShopItem) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetItemId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//id查询商品是否存在
func QueryShopItemById(Id TEAM_ID) *share_message.TableShopItem {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()
	rp := &share_message.TableShopItem{}

	errc := col.Find(bson.M{"_id": Id}).One(rp)
	if errc != nil && errc != mgo.ErrNotFound {
		panic(errc)
	}
	if errc == mgo.ErrNotFound {
		return nil
	}
	return rp
}

//查询商城留言列表
func QueryShopComment(reqMsg *brower_backstage.QueryShopCommentRequest) ([]*share_message.TableItemComment, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ITEM_COMMENT)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	queryBson["item_id"] = reqMsg.GetItemId()
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 &&
		reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
		queryBson["create_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	} else if (reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0) &&
		(reqMsg.EndTimestamp == nil || reqMsg.GetEndTimestamp() == 0) {
		queryBson["create_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp()}
	} else if (reqMsg.BeginTimestamp == nil || reqMsg.GetBeginTimestamp() == 0) &&
		(reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0) {
		queryBson["create_time"] = bson.M{"$lte": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.ComType != nil && reqMsg.GetComType() != SHOP_COMMENT_TYPE_DEFAULT {
		queryBson["star_level"] = reqMsg.GetComType()
	}

	if reqMsg.Nickname != nil && reqMsg.GetNickname() != "" {
		queryBson["nickname"] = reqMsg.GetNickname()
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	sort := "-_id"

	var list []*share_message.TableItemComment
	errc := query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	if errc != nil && errc != mgo.ErrNotFound {
		easygo.PanicError(errc)
	}

	return list, count
}

//留言点赞数修改
func EditShopComment(reqMsg *brower_backstage.EditShopCommentRequest) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ITEM_COMMENT)
	defer closeFun()
	err := col.Update(
		bson.M{"_id": reqMsg.GetCommentId()},
		bson.M{"$set": bson.M{"fake_like_count": reqMsg.GetFakeLikeCount()}})
	easygo.PanicError(err)
}

//id查询订单是否存在
func QueryShopOrderById(Id SHOP_ORDER_ID) *share_message.TableShopOrder {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()
	rp := &share_message.TableShopOrder{}

	errc := col.Find(bson.M{"_id": Id}).One(rp)
	if errc != nil && errc != mgo.ErrNotFound {
		panic(errc)
	}
	if errc == mgo.ErrNotFound {
		return nil
	}
	return rp
}

//查询商城订单列表
func QueryShopOrder(reqMsg *brower_backstage.QueryShopOrderRequest) ([]*share_message.TableShopOrder, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}
	if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 {
		queryBson["create_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}
	switch reqMsg.GetTimeTypes() {
	case SHOP_ORDER_TIME_TYPE_1:
		if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 &&
			reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
			queryBson["create_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
		} else if (reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0) &&
			(reqMsg.EndTimestamp == nil || reqMsg.GetEndTimestamp() == 0) {
			queryBson["create_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp()}
		} else if (reqMsg.BeginTimestamp == nil || reqMsg.GetBeginTimestamp() == 0) &&
			(reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0) {
			queryBson["create_time"] = bson.M{"$lte": reqMsg.GetEndTimestamp()}
		}
	case SHOP_ORDER_TIME_TYPE_2:
		if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 &&
			reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
			queryBson["pay_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
		} else if (reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0) &&
			(reqMsg.EndTimestamp == nil || reqMsg.GetEndTimestamp() == 0) {
			queryBson["pay_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp()}
		} else if (reqMsg.BeginTimestamp == nil || reqMsg.GetBeginTimestamp() == 0) &&
			(reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0) {
			queryBson["pay_time"] = bson.M{"$lte": reqMsg.GetEndTimestamp()}
		}
	case SHOP_ORDER_TIME_TYPE_3:
		if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 &&
			reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
			queryBson["send_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
		} else if (reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0) &&
			(reqMsg.EndTimestamp == nil || reqMsg.GetEndTimestamp() == 0) {
			queryBson["send_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp()}
		} else if (reqMsg.BeginTimestamp == nil || reqMsg.GetBeginTimestamp() == 0) &&
			(reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0) {
			queryBson["send_time"] = bson.M{"$lte": reqMsg.GetEndTimestamp()}
		}
	case SHOP_ORDER_TIME_TYPE_4:
		if reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0 &&
			reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0 {
			queryBson["finish_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
		} else if (reqMsg.BeginTimestamp != nil && reqMsg.GetBeginTimestamp() != 0) &&
			(reqMsg.EndTimestamp == nil || reqMsg.GetEndTimestamp() == 0) {
			queryBson["finish_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp()}
		} else if (reqMsg.BeginTimestamp == nil || reqMsg.GetBeginTimestamp() == 0) &&
			(reqMsg.EndTimestamp != nil && reqMsg.GetEndTimestamp() != 0) {
			queryBson["finish_time"] = bson.M{"$lte": reqMsg.GetEndTimestamp()}
		}
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() != SHOP_ORDER_STATUS_DEFAULT {
		if reqMsg.GetStatus() == SHOP_ORDER_STATUS_CANCEL {
			queryBson["$or"] = []bson.M{
				{"state": for_game.SHOP_ORDER_EXPIRE},
				{"state": for_game.SHOP_ORDER_CANCEL},
				{"state": for_game.SHOP_ORDER_BACKSTAGE_CANCLE}}
		} else {
			queryBson["state"] = reqMsg.GetStatus()
		}
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetTypes() {
		case SHOP_ORDER_TYPE_1:
			queryBson["_id"] = easygo.StringToInt64noErr(reqMsg.GetKeyword())
		case SHOP_ORDER_TYPE_2:
			queryBson["items.item_id"] = easygo.StringToInt64noErr(reqMsg.GetKeyword())
		case SHOP_ORDER_TYPE_3:
			queryBson["sponsor_account"] = reqMsg.GetKeyword()
		case SHOP_ORDER_TYPE_4:
			queryBson["receiver_account"] = reqMsg.GetKeyword()
		case SHOP_ORDER_TYPE_5:
			queryBson["h5_search_con"] = reqMsg.GetKeyword()
		}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	sort := "-_id"

	var list []*share_message.TableShopOrder
	errc := query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	if errc != nil && errc != mgo.ErrNotFound {
		easygo.PanicError(errc)
	}

	return list, count
}

//通知大厅到商城发货商城订单前做的处理:更新快递单号和快递公司
func UpdateOrderExpressInfo(reqMsg *brower_backstage.SendShopOrderRequest) string {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()
	expressName, _ := GetExpressNamePhone(reqMsg.GetExpressCom())
	err := col.Update(
		bson.M{"_id": reqMsg.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_SEND},
		bson.M{"$set": bson.M{"express_code": reqMsg.GetExpressCode(),
			"express_com":          reqMsg.GetExpressCom(),
			"express_name":         expressName,
			"receiver_notify_flag": true,
			"sponsor_notify_flag":  true,
		}})
	if err == mgo.ErrNotFound {

		return "该订单已经发货,请刷新列表"
	}

	if err != nil {
		easygo.PanicError(err)
	}
	return ""
}

//通知大厅到商城取消商城订单前做的处理 更新取消原因
func UpdateOrderCancelReason(reqMsg *brower_backstage.CancelShopOrderRequest, orderState int32) string {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()

	err := col.Update(
		bson.M{"_id": reqMsg.GetOrderId(), "state": orderState},
		bson.M{"$set": bson.M{"cancel_reason": reqMsg.GetCancelReason()}})

	if err == mgo.ErrNotFound {

		return "该订单状态变化,请刷新列表"
	}

	if err != nil {
		easygo.PanicError(err)
	}
	return ""
}

//查询商城收货地址
func QueryShopReceiveAddress(id PLAYER_ID) ([]*share_message.TableReceiveAddress, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_RECEIVE_ADDRESS)
	defer closeFun()

	queryBson := bson.M{"player_id": id}
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	sort := "-_id"

	var list []*share_message.TableReceiveAddress
	errc := query.Sort(sort).All(&list)
	if errc != nil && errc != mgo.ErrNotFound {
		easygo.PanicError(errc)
	}

	return list, count
}

//查询商城发货地址
func QueryShopDeliverAddress(id PLAYER_ID) ([]*share_message.TableDeliverAddress, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_DELIVER_ADDRESS)
	defer closeFun()

	queryBson := bson.M{"player_id": id}
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	sort := "-_id"

	var list []*share_message.TableDeliverAddress
	errc := query.Sort(sort).All(&list)
	if errc != nil && errc != mgo.ErrNotFound {
		easygo.PanicError(errc)
	}

	return list, count
}

//删除评论
func DeleteShopComment(Id int64) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ITEM_COMMENT)
	defer closeFun()

	var itemComment share_message.TableItemComment = share_message.TableItemComment{}

	errQuery := col.Find(bson.M{"_id": Id}).One(&itemComment)

	if errQuery == nil {
		err := col.Update(
			bson.M{"_id": Id},
			bson.M{"$set": bson.M{"status": for_game.SHOP_COMMENT_DELETE}})

		if err != nil && err != mgo.ErrNotFound {
			easygo.PanicError(err)
		}

		if itemComment.GetItemId() != 0 {

			//处理好恢复商品表中的数据
			easygo.Spawn(func(itemIdPara int64, startLevel int32) {

				colItem, closeFunItem := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
				defer closeFunItem()

				var item share_message.TableShopItem = share_message.TableShopItem{}
				eItem := colItem.Find(bson.M{"_id": itemIdPara}).One(&item)
				if eItem != nil {
					logs.Error(eItem)
				}

				err1 := colItem.UpdateId(item.GetItemId(),
					bson.M{"$inc": bson.M{"real_commentCnt": -1}})
				logs.Error(err1)

				if startLevel != for_game.SHOP_COMMENT_LEVEL_COMMON {
					if item.GetFakeFinCommCnt() <= 0 {

						err2 := colItem.Update(
							bson.M{"_id": item.GetItemId()},
							bson.M{"$inc": bson.M{"real_finCommCnt": -1}})
						logs.Error(err2)
					} else {
						err3 := colItem.UpdateId(item.GetItemId(),
							bson.M{"$inc": bson.M{"real_finCommCnt": -1, "fake_finCommCnt": -1}})

						logs.Error(err3)
					}
				}

				if startLevel == for_game.SHOP_COMMENT_LEVEL_GOOD {
					if item.GetFakeGoodCommCnt() <= 0 {

						err4 := colItem.Update(
							bson.M{"_id": item.GetItemId()},
							bson.M{"$inc": bson.M{"real_goodCommCnt": -1}})

						logs.Error(err4)
					} else {
						err5 := colItem.UpdateId(item.GetItemId(),
							bson.M{"$inc": bson.M{"real_goodCommCnt": -1, "fake_goodCommCnt": -1}})

						logs.Error(err5)
					}
				}

			}, itemComment.GetItemId(), itemComment.GetStarLevel())

		} else {
			logs.Debug("更新商品表相关的留言数,缺少商品ID")
		}
	}
}

//导入点卡信息正确性检测
//true正确  false不正确
func ShopPointCardInfoCheck(reqMsg *brower_backstage.ImportShopPointCardRequest) bool {
	tempList := reqMsg.GetPointCardList()

	if nil != tempList && len(tempList) > 0 {
		for _, value := range tempList {
			if compressStr(value.GetCardName()) == "" ||
				compressStr(value.GetCardNo()) == "" ||
				compressStr(value.GetCardPassword()) == "" ||
				compressStr(value.GetSellerAccount()) == "" {

				return false
			}
		}
	}

	return true
}

//判断导入文件中是否存在重复卡号
func GetFileRepeatedCardMsg(reqMsg *brower_backstage.ImportShopPointCardRequest) []string {
	msg := make([]string, 0)
	tempList := reqMsg.GetPointCardList()
	tempRepeatedCards := make([]string, 0)

	if nil != tempList && len(tempList) > 0 {
		for i := 0; i < len(tempList); i++ {
			for j := i + 1; j < len(tempList); j++ {
				if compressStr(tempList[i].GetCardNo()) == compressStr(tempList[j].GetCardNo()) {
					tempCheck := CheckFileRepeatedCard(tempRepeatedCards, compressStr(tempList[i].GetCardNo()))
					if !tempCheck {
						tempRepeatedCards = append(tempRepeatedCards, compressStr(tempList[i].GetCardNo()))
						s := fmt.Sprintf("导入重复卡号%v，该文档导入失败", compressStr(tempList[i].GetCardNo()))
						msg = append(msg, s)
						break
					}
				}
			}
		}
	}

	return msg
}

//false 不存在 true已经存在
func CheckFileRepeatedCard(paraStrs []string, paraStr string) bool {
	if nil != paraStrs {
		if len(paraStrs) <= 0 {
			return false
		} else {
			for _, value := range paraStrs {
				if value == paraStr {
					return true
				}

			}
		}
	}

	return false
}

//利用正则表达式压缩字符串，去除空格或制表符,回车,换行
func compressStr(str string) string {
	tempStr := str
	if tempStr == "" {
		return ""
	}
	// 去除空格
	tempStr = strings.Replace(tempStr, " ", "", -1)
	// 去除换行回车
	reg := regexp.MustCompile("\\r+")
	tempStr = reg.ReplaceAllString(tempStr, "")
	reg1 := regexp.MustCompile("\\n+")
	tempStr = reg1.ReplaceAllString(tempStr, "")
	reg2 := regexp.MustCompile("\\r\\n+")
	tempStr = reg2.ReplaceAllString(tempStr, "")

	//匹配一个或多个空白符的正则表达式
	reg3 := regexp.MustCompile("\\s+")
	return reg3.ReplaceAllString(tempStr, "")
}

//判断数据库是否已经存在文件中的卡号
func GetDbRepeatedCardMsg(reqMsg *brower_backstage.ImportShopPointCardRequest) []string {
	msg := make([]string, 0)
	tempList := reqMsg.GetPointCardList()
	tempCards := make([]string, 0)

	if nil != tempList && len(tempList) > 0 {
		for i := 0; i < len(tempList); i++ {
			tempCards = append(tempCards, compressStr(tempList[i].GetCardNo()))
		}
		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_POINT_CARD)
		defer closeFun()

		var list []*share_message.TableShopPointCard

		e := col.Find(bson.M{"card_no": bson.M{"$in": tempCards}}).All(&list)
		if e != nil && e != mgo.ErrNotFound {
			easygo.PanicError(e)
		}

		if e == mgo.ErrNotFound {
			return msg
		}

		if nil != list && len(list) > 0 {
			for _, value := range list {
				s := fmt.Sprintf("导入已存在卡号%v，该文档导入失败", value.GetCardNo())
				msg = append(msg, s)
			}
		}
	}

	return msg
}

//查询点卡列表
func QueryShopPointCard(reqMsg *brower_backstage.QueryShopPointCardRequest) ([]*brower_backstage.ResShopPointCardObject, int) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_POINT_CARD)
	defer closeFun()

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	queryBson := bson.M{}

	if reqMsg.Status != nil && reqMsg.GetStatus() != SHOP_POINT_CARD_STATUS_DEFAULT {
		queryBson["card_status"] = reqMsg.GetStatus()
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetTypes() {
		case SHOP_POINT_CARD_TYPE_1:
			queryBson["_id"] = easygo.StringToInt64noErr(reqMsg.GetKeyword())
		case SHOP_POINT_CARD_TYPE_2:
			queryBson["card_name"] = reqMsg.GetKeyword()
		case SHOP_POINT_CARD_TYPE_3:
			queryBson["card_no"] = reqMsg.GetKeyword()
		case SHOP_POINT_CARD_TYPE_4:
			queryBson["seller_account"] = reqMsg.GetKeyword()
		case SHOP_POINT_CARD_TYPE_5:
			queryBson["order_no"] = easygo.StringToInt64noErr(reqMsg.GetKeyword())
		}
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)

	sort := "-_id"

	var list []*share_message.TableShopPointCard
	errc := query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	if errc != nil && errc != mgo.ErrNotFound {
		easygo.PanicError(errc)
	}

	var resList []*brower_backstage.ResShopPointCardObject = make([]*brower_backstage.ResShopPointCardObject, 0)
	//封装返回对象
	if len(list) > 0 {
		for _, value := range list {
			resShopPointCardObject := &brower_backstage.ResShopPointCardObject{
				CardId:        easygo.NewInt64(value.GetCardId()),
				CardName:      easygo.NewString(value.GetCardName()),
				CardNo:        easygo.NewString(value.GetCardNo()),
				CardPassword:  easygo.NewString(value.GetCardPassword()),
				SellerAccount: easygo.NewString(value.GetSellerAccount()),
				CardStatus:    easygo.NewInt32(value.GetCardStatus()),
				OrderNo:       easygo.NewInt64(value.GetOrderNo()),
			}
			resList = append(resList, resShopPointCardObject)
		}
	}
	return resList, count
}

//通过卖家帐号取得待售库存的点卡名称列表(过滤重复的名称),并且过滤到正在上架的该商家商品的相同点卡的名称
func QueryShopPointCardDropDown(reqMsg *brower_backstage.GetShopPointCardDropDownRequest) []*brower_backstage.KeyValue {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_POINT_CARD)
	defer closeFun()

	//拼装返回对象
	resCardNames := make([]*brower_backstage.KeyValue, 0)

	var list []*share_message.TableShopPointCard
	err := col.Find(bson.M{"seller_account": reqMsg.GetSellerAccount(), "card_status": for_game.SHOP_POINT_CARD_SALE}).All(&list)

	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	if err == mgo.ErrNotFound {
		return resCardNames
	}
	//去重复的点卡名称
	tempRepeatedCardNames := make([]string, 0)

	if nil != list && len(list) > 0 {
		for i := 0; i < len(list); i++ {

			tempCheck := CheckFileRepeatedCard(tempRepeatedCardNames, list[i].GetCardName())
			if tempCheck {
				continue
			}

			var tempFlag1 bool
			for j := i + 1; j < len(list); j++ {
				if list[i] == list[j] {
					tempCheck := CheckFileRepeatedCard(tempRepeatedCardNames, list[i].GetCardName())
					if !tempCheck {
						tempRepeatedCardNames = append(tempRepeatedCardNames, list[i].GetCardName())
					}

					tempFlag1 = true
					break
				}
			}

			if !tempFlag1 {
				tempRepeatedCardNames = append(tempRepeatedCardNames, list[i].GetCardName())
			}
		}
	}

	//取得该商家存在库存的点卡的上架中的商品
	colItem, closeFunItem := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFunItem()

	var listItem []*share_message.TableShopItem = []*share_message.TableShopItem{}
	errItem := colItem.Find(bson.M{"player_account": reqMsg.GetSellerAccount(), "state": for_game.SHOP_ITEM_SALE, "point_card_name": bson.M{"$in": tempRepeatedCardNames}}).All(&listItem)

	if errItem != nil && errItem != mgo.ErrNotFound {
		easygo.PanicError(errItem)
	}

	//过滤掉上架中的点卡名称
	//去重复的点卡名称
	tempRepeatedCardNames1 := make([]string, 0)
	if nil != listItem && len(listItem) > 0 {
		for _, value := range tempRepeatedCardNames {
			var tempFlag bool
			for _, value1 := range listItem {
				if value == value1.GetPointCardName() {
					tempFlag = true
					break
				}
			}
			if !tempFlag {
				tempRepeatedCardNames1 = append(tempRepeatedCardNames1, value)
			}
		}
	}

	if nil != tempRepeatedCardNames1 && len(tempRepeatedCardNames1) > 0 {
		for _, value := range tempRepeatedCardNames1 {
			keyValue := &brower_backstage.KeyValue{
				Key:   easygo.NewString(value),
				Value: easygo.NewString(value),
			}
			resCardNames = append(resCardNames, keyValue)
		}
	} else {
		for _, value := range tempRepeatedCardNames {
			keyValue := &brower_backstage.KeyValue{
				Key:   easygo.NewString(value),
				Value: easygo.NewString(value),
			}
			resCardNames = append(resCardNames, keyValue)
		}
	}

	return resCardNames
}

//通过卖家帐号和点卡名称取得该卖家该点卡的待售库存
func QueryShopPointCardStock(cardName, sellerAccount string) int32 {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_POINT_CARD)
	defer closeFun()

	query := col.Find(bson.M{"seller_account": sellerAccount, "card_name": cardName, "card_status": for_game.SHOP_POINT_CARD_SALE})
	stockCnt, err := query.Count()
	easygo.PanicError(err)

	return int32(stockCnt)
}

//通过快递代码取得快递名和电话号码
func GetExpressNamePhone(com string) (string, string) {

	switch com {
	case "sf":
		return "顺丰", "95338"
	case "sto":
		return "申通", "400-889-5543"
	case "yt":
		return "圆通", "95554"
	case "yd":
		return "韵达", "95546"
	case "tt":
		return "天天", "400-188-8888"
	case "ems":
		return "EMS", "11183"
	case "zto":
		return "中通", "95311"
	case "ht":
		return "百世快递（汇通）", "95320"
	case "db":
		return "德邦", "95353"
	case "jd":
		return "京东快递", "400-603-3600"
	case "zjs":
		return "宅急送", "4006-789-000"
	case "emsg":
		return "EMS国际", "11183"
	case "yzgn":
		return "邮政国内（挂号信）", "11183"
	case "ztky":
		return "中铁快运", "95572"
	case "zhongyou":
		return "中邮物流", "11183"
	case "ztoky":
		return "中通快运", "4000-270-270"
	case "youzheng":
		return "邮政快递", "11185"
	case "bsky":
		return "百世快运", "400-8856-561"
	case "suning":
		return "苏宁快递", "95315"
	default:
		return "", ""
	}
}

func GetItemTypeName(itemType int32) string {

	switch itemType {
	case 1:
		return "手机"
	case 2:
		return "农用物资"
	case 3:
		return "生鲜水果"
	case 4:
		return "童鞋"
	case 5:
		return "园艺植物"
	case 6:
		return "五金工具"
	case 7:
		return "游泳"
	case 8:
		return "电子零件"
	case 9:
		return "动漫/周边"
	case 10:
		return "图书"
	case 11:
		return "宠物/用品"
	case 12:
		return "网络设备"
	case 13:
		return "服饰配件"
	case 14:
		return "家装/建材"
	case 15:
		return "家纺布艺"
	case 16:
		return "珠宝首饰"
	case 17:
		return "钟表眼镜"
	case 18:
		return "古董收藏"
	case 19:
		return "女士鞋靴"
	case 20:
		return "箱包"
	case 21:
		return "男士鞋靴"
	case 22:
		return "办公用品"
	case 23:
		return "游戏设备"
	case 24:
		return "运动户外"
	case 25:
		return "实体卡/券/票"
	case 26:
		return "工艺礼品"
	case 27:
		return "玩具乐器"
	case 28:
		return "母婴用品"
	case 29:
		return "童装"
	case 30:
		return "女士服装"
	case 31:
		return "家具"
	case 32:
		return "居家日用"
	case 33:
		return "家用电器"
	case 34:
		return "个护美妆"
	case 35:
		return "保健护理"
	case 36:
		return "摩托车/用品"
	case 37:
		return "自行车/用品"
	case 38:
		return "汽车/用品"
	case 39:
		return "电动车/用品"
	case 40:
		return "3C数码"
	case 41:
		return "男士服装"
	case 42:
		return "其他闲置"
	case 43:
		return "音像"
	case 44:
		return "演艺/表演类门票"
	case 45:
		return "点卡"
	default:
		return ""
	}
}

//商品下架再次上架判断该商品下是否有待支付的商城订单以及对应的商城订单的支付订单的状态
//true存在,false不存在
func VerdictItemOrders(itemId int64) (bool, string) {

	//取得该商品下的待支付的订单
	colShopOrder, closeFunShopOrder := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFunShopOrder()

	var list []*share_message.TableShopOrder

	eShopOrder := colShopOrder.Find(bson.M{"state": for_game.SHOP_ORDER_WAIT_PAY}).All(&list)

	if eShopOrder != nil && eShopOrder != mgo.ErrNotFound {
		logs.Error(eShopOrder)
		return false, "操作失败"
	}
	if eShopOrder == mgo.ErrNotFound || list == nil || len(list) <= 0 {
		return false, ""
	}

	//判断各个支付订单的状态来确定能否上架
	//组装生成的商城订单用来查询支付订单
	var tempShopOrders []int64 = make([]int64, 0)

	for _, value := range list {
		tempShopOrders = append(tempShopOrders, value.GetOrderId())
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ORDER)
	defer closeFun()

	query := col.Find(bson.M{"PayTargetId": bson.M{"$in": tempShopOrders}, "PayWay": for_game.PAY_TYPE_SHOP, "$or": []bson.M{
		{"PayStatus": for_game.PAY_ST_WAITTING},
		{"PayStatus": for_game.PAY_ST_REFUSE}}})
	count, err := query.Count()

	if err != nil && err != mgo.ErrNotFound {

		logs.Error(err)

		return false, "操作失败"
	}

	if count > 0 {
		return true, ""
	}

	//如果是app批量过去的那要看bill表重新操作
	colBill, closeFunBill := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_BILLS)
	defer closeFunBill()

	var bills []*share_message.TableBill
	//2判断传入的商城子订单是否存在bill中订单的支付订单也需要当作库存恢复
	errBill := colBill.Find(bson.M{"order_list": bson.M{"$in": tempShopOrders}, "state": for_game.SHOP_ORDER_WAIT_PAY}).All(&bills)
	if errBill != nil && errBill != mgo.ErrNotFound {

		logs.Error(errBill)

		return false, "操作失败"
	}

	//这里没有找到可以返回
	if errBill == mgo.ErrNotFound || bills == nil || len(bills) <= 0 {

		return false, ""
	}

	//找到继续判断支付订单
	//组装bill的订单id
	var tempBillOrders []int64 = make([]int64, 0)

	for _, value := range bills {
		tempBillOrders = append(tempBillOrders, value.GetOrderId())
	}

	if errBill != mgo.ErrNotFound && bills != nil && len(bills) > 0 {

		queryByBill := col.Find(bson.M{"PayTargetId": bson.M{"$in": tempBillOrders}, "PayWay": for_game.PAY_TYPE_SHOP, "$or": []bson.M{
			{"PayStatus": for_game.PAY_ST_WAITTING},
			{"PayStatus": for_game.PAY_ST_REFUSE}}})
		countByBill, errByBill := queryByBill.Count()
		if errByBill != nil && errByBill != mgo.ErrNotFound {

			logs.Error(err)

			return false, "操作失败"
		}

		if countByBill > 0 {
			return true, ""
		}
	}

	return false, ""
}
