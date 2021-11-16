// 大厅服务器为[游戏客户端]提供的服务

package shop

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/for_game/greenScan"
	"game_server/pb/share_message"
	"strings"

	"github.com/astaxie/beego/logs"

	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
)

//=============================mainlogic.proto=============================

// 登录大厅都用这一个
//func (self *ServiceForHall) RpcLogin(common *base.Common, reqMsg *client_shop.LoginMsg) easygo.IMessage {
//	logs.Info("RpcLogin,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
//	account := for_game.GetRedisAccountByPhone(reqMsg.GetAccount())
//	if account == nil {
//		res := "进入商城失败，无效的玩家"
//		return easygo.NewFailMsg(res)
//	}
//	base := for_game.GetRedisPlayerBase(account.GetPlayerId())
//	if base == nil {
//		res := "进入商城失败，无效的玩家"
//		return easygo.NewFailMsg(res)
//	}
//	if base.GetToken() != reqMsg.GetToken() {
//		res := LOGIN_TOKEN_WRONG
//		logs.Info("token验证码错误", reqMsg.GetToken(), base.GetToken())
//		return easygo.NewFailMsg(res)
//	}
//	playerId := account.GetPlayerId()
//	oldEp := ClientEpMp.LoadEndpoint(playerId)
//	if oldEp != nil && oldEp != ep {
//		oldEp.Shutdown()
//	}
//	player := PlayerMgr.LoadPlayer(playerId)
//	if player == nil {
//		player = NewPlayer(playerId)
//		player.OnLoadFromDB()
//		PlayerMgr.Store(playerId, player)
//	} else {
//		player.OnLoadFromDB()
//	}
//
//	ep.SetAssociativePlayer(player)
//	msg := player.GetAllPlayerInfo()
//	ep.RpcPlayerLoginResponse(msg)
//	//存储玩家节点
//	ClientEpMp.StoreEndpoint(playerId, ep.GetEndpointId())
//
//	return easygo.EmptyMsg
//}

//退出登录
//func (self *ServiceForHall) RpcLogOut(ep IGameClientEndpoint, who *Player, reqMsg *base.Empty) easygo.IMessage {
//	logs.Info("RpcLogOut,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
//	playerId := who.GetPlayerId()
//	PlayerMgr.Delete(playerId)
//	return easygo.EmptyMsg
//}

//心跳协议
//func (self *ServiceForHall) RpcHeartbeat(common *base.Common, reqMsg *client_server.NTP) easygo.IMessage {
//	reqMsg.T2 = easygo.NewInt64(time.Now().Unix())
//	return reqMsg
//}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcShopItemUpload(common *base.Common, reqMsg *share_message.ShopItemUploadInfo) easygo.IMessage {
	logs.Info("===api RpcShopItemUpload===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	if len(reqMsg.GetName()) == 0 {
		logs.Debug("缺少好物标题信息")
		return &share_message.ShopItemUploadResult{Result: easygo.NewInt32(1), Msg: easygo.NewString(UPLOAD_ITEM_NAME_IS_NULL)}
	}

	if len(reqMsg.GetTitle()) == 0 {
		logs.Debug("商品描述为不能空")
		return &share_message.ShopItemUploadResult{Result: easygo.NewInt32(1), Msg: easygo.NewString(UPLOAD_ITEM_TITLE_IS_NULL)}
	}

	if len(reqMsg.GetItemFiles()) < 2 || len(reqMsg.GetItemFiles()) > 9 {
		logs.Debug("图片只能发布2~9张")
		return &share_message.ShopItemUploadResult{Result: easygo.NewInt32(1), Msg: easygo.NewString(UPLOAD_ITEM_IMAGE_COUNT)}
	}

	if reqMsg.GetStockCount() < 1 || reqMsg.GetStockCount() > 999 {
		logs.Debug("好物库存数量只能为1~999之间")
		return &share_message.ShopItemUploadResult{Result: easygo.NewInt32(1), Msg: easygo.NewString(UPLOAD_ITEM_COUNT_OUT_OF_RANGE)}
	}

	if reqMsg.GetPrice() < 0 || reqMsg.GetPrice() > 500000 {
		logs.Debug("好物价格只能为1~5000之间")
		return &share_message.ShopItemUploadResult{Result: easygo.NewInt32(1), Msg: easygo.NewString(UPLOAD_ITEM_PRICE_VAR)}
	}

	if len(reqMsg.GetAddress()) == 0 || len(reqMsg.GetDetailAddress()) == 0 {
		logs.Debug("位置不能为空")
		return &share_message.ShopItemUploadResult{Result: easygo.NewInt32(1), Msg: easygo.NewString(UPLOAD_ITEM_ADDRESS_IS_NULL)}
	}

	if reqMsg.GetType() == nil {
		reqMsg.Type = &share_message.ShopItemTypeName{}
	}
	who := for_game.GetRedisPlayerBase(common.GetUserId())
	//若已达到20，则不予上架并在该界面弹出tip提示“已超过单人上架最多商品数”
	colVar, closeFunVar := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	query := colVar.Find(
		bson.M{"player_id": who.GetPlayerId(),
			"$or": []bson.M{
				bson.M{"state": for_game.SHOP_ITEM_SALE},
				bson.M{"state": for_game.SHOP_ITEM_IN_AUDIT}}})
	count, errVa := query.Count()
	closeFunVar()
	if errVa != nil {
		logs.Error(errVa)

		//发布时候需要提示发布中需改成客户端接收到的提示
		return &share_message.ShopItemUploadResult{
			Result: easygo.NewInt32(1),
			Msg:    easygo.NewString(DATABASE_ERROR)}
	}
	if count >= 20 {
		logs.Debug("已超过单人上架最多商品数")

		//发布时候需要提示发布中需改成客户端接收到的提示
		return &share_message.ShopItemUploadResult{
			Result: easygo.NewInt32(1),
			Msg:    easygo.NewString(UPLOAD_ITEM_PEOPLE_VALIDATE_COUNT)}
	}

	itemId := easygo.NewInt64(for_game.NextId(for_game.TABLE_SHOP_ITEMS))

	var storeCount int32 = 0
	var state int32 = for_game.SHOP_ITEM_SALE
	if for_game.IS_FORMAL_SERVER {
		state = for_game.SHOP_ITEM_IN_AUDIT
	}

	var lockCount int32 = 0
	var timeNow int64 = time.Now().Unix()

	//重新设置type
	var itemType share_message.ShopItemType = share_message.ShopItemType{
		Type:      easygo.NewInt32(GetItemType(reqMsg.Type.GetTypeName())),
		OtherType: reqMsg.Type.OtherType}

	newItem := &share_message.TableShopItem{
		ItemId:        itemId,
		Price:         easygo.NewInt32(reqMsg.GetPrice()),
		Type:          &itemType,
		ItemFiles:     reqMsg.ItemFiles,
		Title:         easygo.NewString(reqMsg.GetTitle()),
		UserName:      easygo.NewString(reqMsg.GetUserName()),
		Phone:         easygo.NewString(reqMsg.GetPhone()),
		Address:       easygo.NewString(reqMsg.GetAddress()),
		DetailAddress: easygo.NewString(reqMsg.GetDetailAddress()),
		PlayerId:      easygo.NewInt64(who.GetPlayerId()),
		PlayerAccount: easygo.NewString(who.GetAccount()),
		Nickname:      easygo.NewString(who.GetNickName()),
		Avatar:        easygo.NewString(who.GetHeadIcon()),
		State:         easygo.NewInt32(state),
		Sex:           easygo.NewInt32(who.GetSex()),
		StockCount:    easygo.NewInt32(reqMsg.GetStockCount()),
		CreateTime:    easygo.NewInt64(timeNow),
		LockCount:     easygo.NewInt32(lockCount),
		Name:          easygo.NewString(reqMsg.GetName()),
		RealStoreCnt:  easygo.NewInt32(storeCount),
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	e := col.Insert(newItem)
	closeFun()

	if e != nil {
		logs.Error(e)

		//发布时候需要提示发布中需改成客户端接收到的提示
		return &share_message.ShopItemUploadResult{
			Result: easygo.NewInt32(1),
			Msg:    easygo.NewString(DATABASE_ERROR)}
	}

	itemFiles := []ItemFile{}
	if nil != reqMsg.ItemFiles && len(reqMsg.ItemFiles) > 0 {
		for i := 0; i < len(reqMsg.ItemFiles); i++ {
			var itemFile ItemFile = ItemFile{
				file_url:    reqMsg.ItemFiles[i].GetFileUrl(),
				file_type:   reqMsg.ItemFiles[i].GetFileType(),
				file_width:  reqMsg.ItemFiles[i].GetFileWidth(),
				file_height: reqMsg.ItemFiles[i].GetFileHeight()}
			itemFiles = append(itemFiles, itemFile)
		}
	}

	return &share_message.ShopItemUploadResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(UPLOAD_ITEM_SUCCESS),
		ItemId: itemId}
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcShopItemEdit(common *base.Common, reqMsg *share_message.ShopItemEditInfo) easygo.IMessage {
	logs.Info("===api RpcShopItemEdit===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	var item = ShopInstance.GetItemFromCache(reqMsg.GetItemId())

	if item == nil {
		var data share_message.TableShopItem = share_message.TableShopItem{}
		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
		e := col.Find(bson.M{"_id": reqMsg.GetItemId()}).One(&data)
		closeFun()

		if e != nil && e != mgo.ErrNotFound {
			logs.Error(e)
			//发布时候需要提示发布中需改成客户端接收到的提示
			return &share_message.ShopItemUploadResult{
				Result: easygo.NewInt32(1),
				Msg:    easygo.NewString(DATABASE_ERROR)}
		}

		if e == mgo.ErrNotFound {

			//发布时候需要提示发布中需改成客户端接收到的提示
			return &share_message.ShopItemUploadResult{
				Result: easygo.NewInt32(1),
				Msg:    easygo.NewString(EDIT_ITEM_ITEM_NOT_EXIST)}
		}

		itemFiles := []ItemFile{}
		if nil != data.ItemFiles && len(data.ItemFiles) > 0 {
			for i := 0; i < len(data.ItemFiles); i++ {
				var itemFile ItemFile = ItemFile{
					file_url:    data.ItemFiles[i].GetFileUrl(),
					file_type:   data.ItemFiles[i].GetFileType(),
					file_width:  data.ItemFiles[i].GetFileWidth(),
					file_height: data.ItemFiles[i].GetFileHeight()}
				itemFiles = append(itemFiles, itemFile)
			}
		}

		newItem := ShopItem{
			item_id:        data.GetItemId(),
			item_files:     itemFiles,
			title:          data.GetTitle(),
			origin_price:   data.GetOriginPrice(),
			price:          data.GetPrice(),
			userName:       data.GetUserName(),
			phone:          data.GetPhone(),
			address:        data.GetAddress(),
			detail_address: data.GetDetailAddress(),
			avatar:         data.GetAvatar(),
			player_id:      data.GetPlayerId(),
			item_type:      data.GetType().GetType(),
			other_type:     data.GetType().GetOtherType(),
			nickname:       data.GetNickname(),
			create_time:    data.GetCreateTime(),
			stock_count:    data.GetStockCount(),
			account:        data.GetPlayerAccount(),
			sex:            data.GetSex(),
			name:           data.GetName(),
			state:          data.GetState(),
		}
		item = &newItem
	}

	if item == nil {

		//发布时候需要提示发布中需改成客户端接收到的提示
		return &share_message.ShopItemUploadResult{
			Result: easygo.NewInt32(1),
			Msg:    easygo.NewString(EDIT_ITEM_ITEM_NOT_EXIST)}
	}

	if len(reqMsg.Info.GetItemFiles()) < 2 || len(reqMsg.Info.GetItemFiles()) > 9 {
		logs.Debug("照片只能发布2~9张")

		//发布时候需要提示发布中需改成客户端接收到的提示
		return &share_message.ShopItemUploadResult{
			Result: easygo.NewInt32(1),
			Msg:    easygo.NewString(EDIT_ITEM_IMAGE_NOT_ENOUGH)}
	}

	if reqMsg.Info.GetStockCount() < 1 || reqMsg.Info.GetStockCount() > 999 {
		logs.Debug("好物库存数量只能为1~999之间")

		//发布时候需要提示发布中需改成客户端接收到的提示
		return &share_message.ShopItemUploadResult{
			Result: easygo.NewInt32(1),
			Msg:    easygo.NewString(EDIT_ITEM_COUNT_OUT_OF_RANGE)}
	}

	if reqMsg.Info.GetType() == nil {
		reqMsg.Info.Type = &share_message.ShopItemTypeName{}
	}

	newItemId := easygo.NewInt64(for_game.NextId(for_game.TABLE_SHOP_ITEMS))

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	err := col.Update(
		bson.M{"_id": reqMsg.GetItemId()},
		bson.M{"$set": bson.M{"state": for_game.SHOP_ITEM_SOLD_OUT, "sold_out_time": time.Now().Unix()}})
	closeFun()

	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)

		//发布时候需要提示发布中需改成客户端接收到的提示
		return &share_message.ShopItemUploadResult{
			Result: easygo.NewInt32(1),
			Msg:    easygo.NewString(DATABASE_ERROR)}
	}
	var storeCount int32 = 0
	var state int32 = for_game.SHOP_ITEM_SALE
	if for_game.IS_FORMAL_SERVER {
		state = for_game.SHOP_ITEM_IN_AUDIT
	}

	var timeNow int64 = time.Now().Unix()
	var lockCount int32 = 0
	//重新设置type
	var itemType share_message.ShopItemType = share_message.ShopItemType{
		Type:      easygo.NewInt32(GetItemType(reqMsg.Info.Type.GetTypeName())),
		OtherType: reqMsg.Info.Type.OtherType}

	newItem := &share_message.TableShopItem{
		ItemId:        newItemId,
		Price:         reqMsg.Info.Price,
		Type:          &itemType,
		ItemFiles:     reqMsg.Info.ItemFiles,
		Title:         easygo.NewString(reqMsg.Info.GetTitle()),
		UserName:      easygo.NewString(reqMsg.Info.GetUserName()),
		Phone:         easygo.NewString(reqMsg.Info.GetPhone()),
		Address:       easygo.NewString(reqMsg.Info.GetAddress()),
		DetailAddress: easygo.NewString(reqMsg.Info.GetDetailAddress()),
		PlayerId:      easygo.NewInt64(item.player_id),
		Nickname:      easygo.NewString(item.nickname),
		PlayerAccount: easygo.NewString(item.account),
		Avatar:        easygo.NewString(item.avatar),
		State:         easygo.NewInt32(state),
		Sex:           easygo.NewInt32(item.sex),
		StockCount:    easygo.NewInt32(reqMsg.Info.GetStockCount()),
		CreateTime:    easygo.NewInt64(timeNow),
		LockCount:     easygo.NewInt32(lockCount),
		Name:          easygo.NewString(reqMsg.Info.GetName()),
		RealStoreCnt:  easygo.NewInt32(storeCount),
	}

	col, closeFun = MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	err = col.Insert(newItem)
	closeFun()

	if err != nil {
		logs.Error(err)

		//发布时候需要提示发布中需改成客户端接收到的提示
		return &share_message.ShopItemUploadResult{
			Result: easygo.NewInt32(1),
			Msg:    easygo.NewString(DATABASE_ERROR)}
	}

	itemFiles := []ItemFile{}
	if nil != reqMsg.Info.ItemFiles && len(reqMsg.Info.ItemFiles) > 0 {
		for i := 0; i < len(reqMsg.Info.ItemFiles); i++ {
			var itemFile ItemFile = ItemFile{
				file_url:    reqMsg.Info.ItemFiles[i].GetFileUrl(),
				file_type:   reqMsg.Info.ItemFiles[i].GetFileType(),
				file_width:  reqMsg.Info.ItemFiles[i].GetFileWidth(),
				file_height: reqMsg.Info.ItemFiles[i].GetFileHeight()}
			itemFiles = append(itemFiles, itemFile)
		}
	}

	return &share_message.ShopItemUploadResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(EDIT_ITEM_SUCCESS),
		ItemId: newItemId}
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcShopItemDelete(common *base.Common, reqMsg *share_message.ShopItemID) easygo.IMessage {
	logs.Info("===api RpcShopItemDelete===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()
	err := col.Update(
		bson.M{"_id": reqMsg.GetItemId()},
		bson.M{"$set": bson.M{"state": for_game.SHOP_ITEM_DELETE}})

	if err == mgo.ErrNotFound {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DELETE_ITEM_ITEM_NOT_EXIST)
		return easygo.NewFailMsg(DELETE_ITEM_ITEM_NOT_EXIST)
	}

	if err != nil {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)
		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	return &share_message.ShopItemDeleteResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(DELETE_ITEM_SUCCESS)}
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcShopItemSoldOut(common *base.Common, reqMsg *share_message.ShopItemID) easygo.IMessage {
	logs.Info("===api RpcShopItemSoldOut===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()
	err := col.Update(
		bson.M{"_id": reqMsg.GetItemId()},
		bson.M{"$set": bson.M{"state": for_game.SHOP_ITEM_SOLD_OUT, "sold_out_time": time.Now().Unix()}})

	if err == mgo.ErrNotFound {

		SendToHallServerByApi(common.GetUserId(), "RpcToast", SOLD_OUT_ITEM_ITEM_NOT_EXIST)
		return easygo.NewFailMsg(SOLD_OUT_ITEM_SUCCESS)
	}

	if err != nil {

		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)
		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	return &share_message.ShopItemSoldOutResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(SOLD_OUT_ITEM_SUCCESS)}
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcShopItemList(common *base.Common, reqMsg *share_message.ShopInfo) easygo.IMessage {
	logs.Info("===api RpcShopItemList===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	//为了能过通联支付审核 要修改分类
	itemType := [][]int32{
		{0},     //全部
		{10086}, //推荐
		//{27, 28},             //母婴玩具
		{8, 12, 23, 40, 1}, //数码电器
		//{34, 35},                                                 //美妆个护
		{4, 13, 16, 17, 19, 20, 21, 29, 30, 41}, //服装服饰
		{14, 15, 31, 32, 33},                    //家居家具
		//{10},                                                   //图书文具
		//{11}, //宠物用品
		{2, 3, 5, 6, 7, 9, 18, 22, 24, 25, 26, 36, 37, 38, 39, 42, 43, 44, 45}, //其他
	}
	//推介
	if reqMsg.GetType() == 1 {
		return ShopInstance.Recommend(
			reqMsg.GetPage(),
			reqMsg.GetPageSize(),
			itemType[0],
			common.GetUserId(),
			reqMsg.CacheItemTypes,
			reqMsg.CacheSearch)

		//推介以外
	} else {
		var list []*share_message.ShopItem = []*share_message.ShopItem{}

		if int32(len(itemType)) <= reqMsg.GetType() || reqMsg.GetType() < 0 {
			var count int32 = 0
			return &share_message.ItemList{
				Items:    list,
				Page:     reqMsg.Page,
				PageSize: reqMsg.PageSize,
				Count:    &count}
		}

		result := ShopInstance.Show(
			reqMsg.GetPage(),
			reqMsg.GetPageSize(),
			itemType[reqMsg.GetType()],
			common.GetUserId(),
			reqMsg.CacheItemTypes)
		return result
	}
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcShopItemInfo(common *base.Common, reqMsg *share_message.ShopItemInfo) easygo.IMessage {
	logs.Info("===api RpcShopItemInfo===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	var blackList []PLAYER_ID = ShopInstance.GetBlackLists(common.GetUserId())

	if *reqMsg.Flag == share_message.BuySell_Type_Buyer {
		errStr1, shopItem := ShopInstance.GetShopItem(*reqMsg.Flag, reqMsg.GetItemId())
		if errStr1 != "" {
			//客户端为了增加物品详情页的缓存，先显示详情然后取数据所以不能走通用的msg
			return &share_message.ShopItemShowDetail{
				Result: easygo.NewInt32(1),
				Msg:    easygo.NewString(errStr1)}
		} else {

			var addFlag int32 = 0
			for _, black := range blackList {
				if shopItem.player_id == black {
					addFlag = 1
					break
				}
			}
			if addFlag == 1 {
				//客户端为了增加物品详情页的缓存，先显示详情然后取数据所以不能走通用的msg
				return &share_message.ShopItemShowDetail{
					Result: easygo.NewInt32(1),
					Msg:    easygo.NewString(DETAIL_SHOP_ITEM_BLACK_VAR)}
			}

			sellerInfo := ShopInstance.SellerInfo(shopItem)
			itemDetail := ShopInstance.ItemDetail(shopItem)
			if nil == sellerInfo {
				//客户端为了增加物品详情页的缓存，先显示详情然后取数据所以不能走通用的msg
				return &share_message.ShopItemShowDetail{
					Result: easygo.NewInt32(1),
					Msg:    easygo.NewString(DATABASE_ERROR)}
			}
			if nil == itemDetail {
				//客户端为了增加物品详情页的缓存，先显示详情然后取数据所以不能走通用的msg
				return &share_message.ShopItemShowDetail{
					Result: easygo.NewInt32(1),
					Msg:    easygo.NewString(DETAIL_SHOP_ITEM_NOT_EXIST)}
			}

			col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_STORE)
			defer closeFun()
			query := col.Find(bson.M{"player_id": common.GetUserId(), "item_id": itemDetail.GetItemId()})
			count, err := query.Count()

			if err != mgo.ErrNotFound && err != nil {
				SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)
				return easygo.NewFailMsg(DATABASE_ERROR)
			}

			if count > 0 {
				itemDetail.StoreFlag = easygo.NewBool(true)
			}
			var commInfoForDetail = ShopInstance.CommInfoForDetail(shopItem, reqMsg.Flag)

			//记录当天浏览量
			easygo.Spawn(func(pageViewFlag int32, item *ShopItem) {
				if item != nil {
					ShopInstance.UpdatePageViews(pageViewFlag, item)
				} else {
					logs.Debug("记录浏览量时缺少商品信息")
				}
			}, reqMsg.GetPageViewFlag(), shopItem)

			return &share_message.ShopItemShowDetail{
				Result:         easygo.NewInt32(0),
				Msg:            easygo.NewString(""),
				ShopItemDetail: itemDetail,
				SellerInfo:     sellerInfo,
				CommentInfo:    commInfoForDetail}
		}
	} else {
		errStr2, shopItem := ShopInstance.GetShopItem(*reqMsg.Flag, reqMsg.GetItemId())
		if errStr2 != "" {
			//客户端为了增加物品详情页的缓存，先显示详情然后取数据所以不能走通用的msg
			return &share_message.ShopItemShowDetail{
				Result: easygo.NewInt32(1),
				Msg:    easygo.NewString(errStr2)}
		}

		sellerInfo := ShopInstance.SellerInfo(shopItem)
		itemDetail := ShopInstance.ItemDetail(shopItem)

		if nil == sellerInfo {
			//客户端为了增加物品详情页的缓存，先显示详情然后取数据所以不能走通用的msg
			return &share_message.ShopItemShowDetail{
				Result: easygo.NewInt32(1),
				Msg:    easygo.NewString(DATABASE_ERROR)}
		}

		if nil == itemDetail {
			//客户端为了增加物品详情页的缓存，先显示详情然后取数据所以不能走通用的msg
			return &share_message.ShopItemShowDetail{
				Result: easygo.NewInt32(1),
				Msg:    easygo.NewString(DETAIL_SHOP_ITEM_NOT_EXIST)}
		}

		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_STORE)
		defer closeFun()
		query := col.Find(bson.M{"player_id": common.GetUserId(), "item_id": itemDetail.GetItemId()})
		count, err := query.Count()

		if err != mgo.ErrNotFound && err != nil {
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)
			return easygo.NewFailMsg(DATABASE_ERROR)
		}

		if count > 0 {
			itemDetail.StoreFlag = easygo.NewBool(true)
		}

		var commInfoForDetail = ShopInstance.CommInfoForDetail(shopItem, reqMsg.Flag)

		//记录当天浏览量
		easygo.Spawn(func(pageViewFlag int32, item *ShopItem) {
			if item != nil {
				ShopInstance.UpdatePageViews(pageViewFlag, item)
			} else {
				logs.Debug("记录浏览量时缺少商品信息")
			}
		}, reqMsg.GetPageViewFlag(), shopItem)

		return &share_message.ShopItemShowDetail{
			Result:         easygo.NewInt32(0),
			Msg:            easygo.NewString(""),
			ShopItemDetail: itemDetail,
			SellerInfo:     sellerInfo,
			CommentInfo:    commInfoForDetail}
	}
}

// 请求收货地址
func (self *ServiceForHall) RpcReceiveAddressList(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===api RpcReceiveAddressList===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	list := []share_message.TableReceiveAddress{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_RECEIVE_ADDRESS)
	defer closeFun()

	e := col.Find(bson.M{"player_id": common.GetUserId()}).Sort("-create_time").All(&list)
	if e != mgo.ErrNotFound && e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)
		return easygo.NewFailMsg(DATABASE_ERROR)
	}
	infoList := make([]*share_message.ReceiveAddressInfo, 0)
	defaultInfoList := make([]*share_message.ReceiveAddressInfo, 0)
	lastInfoList := make([]*share_message.ReceiveAddressInfo, 0)

	for _, value := range list {
		info := share_message.ReceiveAddressInfo{
			AddressId: value.AddressId,
			Address: &share_message.ReceiveAddress{
				Name:          value.Name,
				Region:        value.Region,
				Phone:         value.Phone,
				DetailAddress: value.DetailAddress,
				DefaultFlag:   value.DefaultFlag}}

		if value.GetDefaultFlag() == 1 {
			defaultInfoList = append(defaultInfoList, &info)
		} else {
			infoList = append(infoList, &info)
		}
	}

	lastInfoList = append(lastInfoList, defaultInfoList...)
	lastInfoList = append(lastInfoList, infoList...)

	return &share_message.ReceiveAddressList{List: lastInfoList}
}

// 编辑收货地址
func (self *ServiceForHall) RpcReceiveAddressEdit(common *base.Common, reqMsg *share_message.ReceiveAddressInfo) easygo.IMessage {
	logs.Info("===api RpcReceiveAddressEdit===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_RECEIVE_ADDRESS)
	defer closeFun()

	//判断是否设置成默认地址
	if reqMsg.Address.GetDefaultFlag() == 1 {
		// 将数据库中的所有设置为不是默认地址
		_, err := col.UpdateAll(bson.M{}, bson.M{"$set": bson.M{"default_flag": 0}})

		if err != nil && err != mgo.ErrNotFound {
			logs.Error(err)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)
			return easygo.NewFailMsg(DATABASE_ERROR)
		}
	}
	nowTime := time.Now().Unix()
	e := col.UpdateId(
		reqMsg.GetAddressId(),
		bson.M{"$set": bson.M{"name": reqMsg.Address.GetName(),
			"phone":          reqMsg.Address.GetPhone(),
			"region":         reqMsg.Address.GetRegion(),
			"detail_address": reqMsg.Address.GetDetailAddress(),
			"default_flag":   reqMsg.Address.GetDefaultFlag(),
			"create_time":    nowTime}})
	if e != nil && e != mgo.ErrNotFound {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	return &share_message.ReceiveAddressEditResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(RECEIVE_ADDRESS_EDIT_SUCCESS)}
}

// 添加收货地址
func (self *ServiceForHall) RpcReceiveAddressAdd(common *base.Common, reqMsg *share_message.ReceiveAddress) easygo.IMessage {
	logs.Info("===api RpcReceiveAddressAdd===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_RECEIVE_ADDRESS)
	defer closeFun()

	//判断是否设置成默认地址
	if reqMsg.GetDefaultFlag() == 1 {
		// 将数据库中的所有设置为不是默认地址
		_, err := col.UpdateAll(
			bson.M{},
			bson.M{"$set": bson.M{"default_flag": 0}})

		if err != nil && err != mgo.ErrNotFound {
			logs.Error(err)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)
		}
	}

	addressId := easygo.NewInt64(for_game.NextId(for_game.TABLE_RECEIVE_ADDRESS))

	nowTime := time.Now().Unix()

	address := &share_message.TableReceiveAddress{
		PlayerId:      easygo.NewInt64(common.GetUserId()),
		AddressId:     addressId,
		Name:          easygo.NewString(reqMsg.GetName()),
		Region:        easygo.NewString(reqMsg.GetRegion()),
		Phone:         easygo.NewString(reqMsg.GetPhone()),
		DetailAddress: easygo.NewString(reqMsg.GetDetailAddress()),
		DefaultFlag:   easygo.NewInt32(reqMsg.GetDefaultFlag()),
		CreateTime:    easygo.NewInt64(nowTime)}

	e := col.Insert(address)
	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	return &share_message.ReceiveAddressAddResult{
		Result:    easygo.NewInt32(0),
		Msg:       easygo.NewString(RECEIVE_ADDRESS_ADD_SUCCESS),
		AddressId: addressId}
}

// 减少收货地址
func (self *ServiceForHall) RpcReceiveAddressDelete(common *base.Common, reqMsg *share_message.ReceiveAddressID) easygo.IMessage {
	logs.Info("===api RpcReceiveAddressDelete===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_RECEIVE_ADDRESS)
	defer closeFun()
	e := col.RemoveId(reqMsg.GetAddressId())

	if e != nil && e != mgo.ErrNotFound {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	return &share_message.ReceiveAddressRemoveResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(RECEIVE_ADDRESS_DELETE_SUCCESS)}
}

// 请求发货地址
func (self *ServiceForHall) RpcDeliverAddressList(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===api RpcDeliverAddressList===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	list := []share_message.TableDeliverAddress{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_DELIVER_ADDRESS)
	defer closeFun()
	e := col.Find(bson.M{"player_id": common.GetUserId()}).Sort("-create_time").All(&list)

	if e != mgo.ErrNotFound && e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	infoList := []*share_message.DeliverAddressInfo{}
	defaultInfoList := []*share_message.DeliverAddressInfo{}
	lastInfoList := []*share_message.DeliverAddressInfo{}

	for _, value := range list {

		info := share_message.DeliverAddressInfo{
			AddressId: value.AddressId,
			Address: &share_message.DeliverAddress{
				Name:          value.Name,
				Region:        value.Region,
				Phone:         value.Phone,
				DetailAddress: value.DetailAddress,
				DefaultFlag:   value.DefaultFlag}}

		if value.GetDefaultFlag() == 1 {
			defaultInfoList = append(defaultInfoList, &info)
		} else {
			infoList = append(infoList, &info)
		}
	}

	lastInfoList = append(lastInfoList, defaultInfoList...)
	lastInfoList = append(lastInfoList, infoList...)

	return &share_message.DeliverAddressList{List: lastInfoList}
}

// 编辑发货地址
func (self *ServiceForHall) RpcDeliverAddressEdit(common *base.Common, reqMsg *share_message.DeliverAddressInfo) easygo.IMessage {
	logs.Info("===api RpcDeliverAddressEdit===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_DELIVER_ADDRESS)
	defer closeFun()

	//判断是否设置成默认地址
	if reqMsg.Address.GetDefaultFlag() == 1 {
		// 将数据库中的所有设置为不是默认地址
		_, err := col.UpdateAll(
			bson.M{},
			bson.M{"$set": bson.M{"default_flag": 0}})

		if err != nil && err != mgo.ErrNotFound {
			logs.Error(err)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)
		}
	}

	nowTime := time.Now().Unix()

	e := col.UpdateId(reqMsg.GetAddressId(),
		bson.M{"$set": bson.M{
			"name":           reqMsg.Address.GetName(),
			"phone":          reqMsg.Address.GetPhone(),
			"region":         reqMsg.Address.GetRegion(),
			"detail_address": reqMsg.Address.GetDetailAddress(),
			"default_flag":   reqMsg.Address.GetDefaultFlag(),
			"create_time":    nowTime}})
	if e != nil && e != mgo.ErrNotFound {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	return &share_message.DeliverAddressEditResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(DELIVER_ADDRESS_EDIT_SUCCESS)}
}

// 添加发货地址
func (self *ServiceForHall) RpcDeliverAddressAdd(common *base.Common, reqMsg *share_message.DeliverAddress) easygo.IMessage {
	logs.Info("===api RpcDeliverAddressAdd===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_DELIVER_ADDRESS)
	defer closeFun()

	//判断是否设置成默认地址
	if reqMsg.GetDefaultFlag() == 1 {
		// 将数据库中的所有设置为不是默认地址
		_, err := col.UpdateAll(
			bson.M{},
			bson.M{"$set": bson.M{"default_flag": 0}})

		if err != nil && err != mgo.ErrNotFound {
			logs.Error(err)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)
		}
	}

	nowTime := time.Now().Unix()
	addressId := easygo.NewInt64(for_game.NextId(for_game.TABLE_DELIVER_ADDRESS))
	address := &share_message.TableDeliverAddress{
		PlayerId:      easygo.NewInt64(common.GetUserId()),
		AddressId:     addressId,
		Name:          easygo.NewString(reqMsg.GetName()),
		Region:        easygo.NewString(reqMsg.GetRegion()),
		Phone:         easygo.NewString(reqMsg.GetPhone()),
		DetailAddress: easygo.NewString(reqMsg.GetDetailAddress()),
		DefaultFlag:   easygo.NewInt32(reqMsg.GetDefaultFlag()),
		CreateTime:    easygo.NewInt64(nowTime)}

	e := col.Insert(address)
	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	return &share_message.DeliverAddressAddResult{
		Result:    easygo.NewInt32(0),
		Msg:       easygo.NewString(DELIVER_ADDRESS_ADD_SUCCESS),
		AddressId: addressId}
}

// 减少发货地址
func (self *ServiceForHall) RpcDeliverAddressDelete(common *base.Common, reqMsg *share_message.DeliverAddressID) easygo.IMessage {
	logs.Info("===api RpcDeliverAddressDelete===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_DELIVER_ADDRESS)
	defer closeFun()

	e := col.RemoveId(reqMsg.GetAddressId())

	if e != nil && e != mgo.ErrNotFound {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	return &share_message.DeliverAddressRemoveResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(DELIVER_ADDRESS_DELETE_SUCCESS)}
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcShopItemCommentUpload(common *base.Common, reqMsg *share_message.UploadComment) easygo.IMessage {
	logs.Info("===api RpcShopItemCommentUpload===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	if reqMsg.GetContent() == "" {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", UPLOAD_ITEM_COMMENT_NOT_BLANK)

		return easygo.NewFailMsg(UPLOAD_ITEM_COMMENT_NOT_BLANK)
	}
	//留言内容审核
	var texts []string = []string{*reqMsg.Content}
	rstText, rstErrCode, rstErrContent := greenScan.GetTextScanResult(texts)

	if rstText == 0 {
		logs.Debug("留言内容审核失败")
		if "" != rstErrCode {
			s := fmt.Sprintf("留言审核失败-玩家id:%v;商品id:%v;留言内容:%v;错误码:%v;错误内容:%v",
				common.GetUserId(), reqMsg.GetItemId(), texts, rstErrCode, rstErrContent)

			for_game.WriteFile("shop_audit.log", s)
		}
		SendToHallServerByApi(common.GetUserId(), "RpcToast", UPLOAD_ITEM_COMMENT_AUDIT_FAIL)

		return easygo.NewFailMsg(UPLOAD_ITEM_COMMENT_AUDIT_FAIL)

	} else if rstText == 2 {

		error := fmt.Sprintf("留言内容验证网络出错%v", reqMsg.GetItemId())
		logs.Error(error)
		for_game.WriteFile("shop_audit.log", error)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	var timeNow int64 = time.Now().Unix()
	commentId := easygo.NewInt64(for_game.NextId(for_game.TABLE_ITEM_COMMENT))
	who := for_game.GetRedisPlayerBase(common.GetUserId())
	newItem := &share_message.TableItemComment{
		CommentId:     commentId,
		ItemId:        easygo.NewInt64(reqMsg.GetItemId()),
		PlayerId:      easygo.NewInt64(common.GetUserId()),
		Nickname:      easygo.NewString(who.GetNickName()),
		Avatar:        easygo.NewString(who.GetHeadIcon()),
		Sex:           easygo.NewInt32(who.GetSex()),
		Content:       easygo.NewString(reqMsg.GetContent()),
		CreateTime:    easygo.NewInt64(timeNow),
		StarLevel:     easygo.NewInt32(for_game.SHOP_COMMENT_LEVEL_COMMON),
		RealLikeCount: easygo.NewInt32(0),
		FakeLikeCount: easygo.NewInt32(0),
		Status:        easygo.NewInt32(for_game.SHOP_COMMENT_NO_REPLY),
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ITEM_COMMENT)
	defer closeFun()
	e := col.Insert(newItem)

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	var commInfo *share_message.CommInfoForDetail
	// 买家0 卖家1
	var peopleFlag share_message.BuySell_Type
	if reqMsg.GetSponsor_Id() != who.GetPlayerId() {
		peopleFlag = share_message.BuySell_Type_Buyer
	} else {
		peopleFlag = share_message.BuySell_Type_Seller
	}

	errStr1, shopItem := ShopInstance.GetShopItem(peopleFlag, reqMsg.GetItemId())

	if errStr1 != "" {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", errStr1)

		return easygo.NewFailMsg(errStr1)

	} else {
		commInfo = ShopInstance.CommInfoForDetail(shopItem, &peopleFlag)
	}

	// 更新商品表相关的留言总数
	easygo.Spawn(func(itemIdPara int64) {

		if itemIdPara != 0 {

			// 更新商品表相关的留言总数
			AddAllCommentCnt(itemIdPara)

		} else {
			logs.Debug("更新商品表相关的留言总数,缺少商品ID")
		}

	}, reqMsg.GetItemId())

	return &share_message.UploadCommentResult{
		Result:  easygo.NewInt32(0),
		Msg:     easygo.NewString(UPLOAD_ITEM_COMMENT_SUCCESS),
		Comment: commInfo}
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcShopItemCommentList(common *base.Common, reqMsg *share_message.ShopCommentList) easygo.IMessage {
	logs.Info("===api RpcShopItemCommentList===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	var count int32 = 0
	page := reqMsg.GetPage()
	pageSize := reqMsg.GetPageSize()

	var list []*share_message.TableItemComment
	commentInfoList := []*share_message.CommentInfo{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ITEM_COMMENT)
	defer closeFun()

	//全部中的留言总数
	var allCommCount int32
	cntAll, errVaAll := col.Find(bson.M{"item_id": reqMsg.GetItemId(), "status": bson.M{"$ne": for_game.SHOP_COMMENT_DELETE}}).Count()
	allCommCount = int32(cntAll)

	if errVaAll != nil {

		logs.Error(errVaAll)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	selectStr := make(bson.M)
	if *reqMsg.QueryCon == share_message.CommQuery_Con_ALL {

		selectStr = bson.M{"item_id": reqMsg.GetItemId(), "status": bson.M{"$ne": for_game.SHOP_COMMENT_DELETE}}
	} else if *reqMsg.QueryCon == share_message.CommQuery_Con_NEW {

		selectStr = bson.M{"item_id": reqMsg.GetItemId(), "status": bson.M{"$ne": for_game.SHOP_COMMENT_DELETE}}
	} else if *reqMsg.QueryCon == share_message.CommQuery_Con_GOOD {

		selectStr = bson.M{"item_id": reqMsg.GetItemId(), "star_level": for_game.SHOP_COMMENT_LEVEL_GOOD, "status": bson.M{"$ne": for_game.SHOP_COMMENT_DELETE}}
	} else if *reqMsg.QueryCon == share_message.CommQuery_Con_MIDDLE {

		selectStr = bson.M{"item_id": reqMsg.GetItemId(), "star_level": for_game.SHOP_COMMENT_LEVEL_MID, "status": bson.M{"$ne": for_game.SHOP_COMMENT_DELETE}}
	} else if *reqMsg.QueryCon == share_message.CommQuery_Con_BAD {

		selectStr = bson.M{"item_id": reqMsg.GetItemId(), "star_level": for_game.SHOP_COMMENT_LEVEL_BAD, "status": bson.M{"$ne": for_game.SHOP_COMMENT_DELETE}}
	}
	query := col.Find(selectStr)
	cnt, errVa := query.Count()
	count = int32(cnt)

	if errVa != nil {
		logs.Error(errVa)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	if *reqMsg.QueryCon == share_message.CommQuery_Con_ALL {

		e := query.Limit(int(pageSize)).Skip(int(page*pageSize)).Sort(
			"-fake_like_count", "-real_like_count").All(&list)

		if e != mgo.ErrNotFound && e != nil {
			logs.Error(e)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)
		}

	} else if *reqMsg.QueryCon == share_message.CommQuery_Con_NEW {
		e := query.Limit(int(pageSize)).Skip(int(page * pageSize)).Sort("-create_time").All(&list)

		if e != mgo.ErrNotFound && e != nil {
			logs.Error(e)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)
		}

	} else if *reqMsg.QueryCon == share_message.CommQuery_Con_GOOD {

		e := query.Limit(int(pageSize)).Skip(int(page * pageSize)).Sort("-create_time").All(&list)
		if e != mgo.ErrNotFound && e != nil {
			logs.Error(e)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)
		}
	} else if *reqMsg.QueryCon == share_message.CommQuery_Con_MIDDLE {

		e := query.Limit(int(pageSize)).Skip(int(page * pageSize)).Sort("-create_time").All(&list)
		if e != mgo.ErrNotFound && e != nil {
			logs.Error(e)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)
		}
	} else if *reqMsg.QueryCon == share_message.CommQuery_Con_BAD {

		e := query.Limit(int(pageSize)).Skip(int(page * pageSize)).Sort("-create_time").All(&list)
		if e != mgo.ErrNotFound && e != nil {
			logs.Error(e)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)
		}
	}

	// 取得该用户的点赞信息
	colLike, closeFunLike := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_LIKE)
	defer closeFunLike()

	if list != nil {
		for _, value := range list {
			//转换显示的时间
			var timeLong string = GetYMDTime(*value.CreateTime)
			//买家视角需要遮掩用户的昵称
			var nickName string
			if common.GetUserId() != reqMsg.GetSponsor_Id() {
				nickName = GetMarkNickName(value.GetNickname())
			} else {
				nickName = value.GetNickname()
			}

			//计算点赞
			var likeCount int32
			if value.GetFakeLikeCount() > 0 {
				likeCount = value.GetFakeLikeCount()
			} else {
				likeCount = value.GetRealLikeCount()
			}
			var dataList []*share_message.TableLikeRecord
			err := colLike.Find(
				bson.M{"comment_id": value.GetCommentId(),
					"player_id": common.GetUserId(),
					"like_flag": true}).All(&dataList)

			if err != nil && err != mgo.ErrNotFound {
				logs.Error(err)
			}

			var isLike bool
			if dataList != nil && len(dataList) > 0 {
				for _, valueLike := range dataList {

					if common.GetUserId() == valueLike.GetPlayerId() {
						isLike = true
						break
					}
				}
			}

			newItem := share_message.CommentInfo{
				CommentId: value.CommentId,
				PlayerId:  value.PlayerId,
				Avatar:    value.Avatar,
				Nickname:  easygo.NewString(nickName),
				Content:   value.Content,
				ItemId:    value.ItemId,
				TimeLong:  easygo.NewString(timeLong),
				Sex:       value.Sex,
				StarLevel: value.StarLevel,
				LikeCount: easygo.NewInt32(likeCount),
				Status:    value.Status,
				IsLike:    easygo.NewBool(isLike),
			}
			commentInfoList = append(commentInfoList, &newItem)
		}
	}
	return &share_message.ShopCommentListResult{
		Result:       easygo.NewInt32(0),
		Msg:          easygo.NewString(""),
		Comments:     commentInfoList,
		AllCommCount: easygo.NewInt32(allCommCount),
		Page:         easygo.NewInt32(page),
		PageSize:     easygo.NewInt32(pageSize),
		Count:        easygo.NewInt32(count),
	}
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcLikeComment(common *base.Common, reqMsg *share_message.LikeCommentInfo) easygo.IMessage {
	logs.Info("===api RpcLikeComment===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	easygo.Spawn(func(commentId int64, likeType bool, playId int64) {

		func() {

			col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ITEM_COMMENT)
			defer closeFun()

			colLike, closeFunLike := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_LIKE)
			defer closeFunLike()

			var data share_message.TableItemComment = share_message.TableItemComment{}
			err := col.FindId(reqMsg.GetCommentId()).One(&data)

			if err != nil && err != mgo.ErrNotFound {
				logs.Error(err)
			}

			if err == nil {
				//取消点赞
				if likeType == false {

					if data.GetFakeLikeCount() > 0 {
						if data.GetRealLikeCount() <= 0 {

							err1 := col.Update(
								bson.M{"_id": reqMsg.GetCommentId()},
								bson.M{"$inc": bson.M{"fake_like_count": -1}})

							if err1 != nil {
								logs.Error(err1)
							}
						} else {
							err1 := col.Update(
								bson.M{"_id": reqMsg.GetCommentId()},
								bson.M{"$inc": bson.M{"real_like_count": -1, "fake_like_count": -1}})

							if err1 != nil {
								logs.Error(err1)
							}
						}

					} else {

						if data.GetRealLikeCount() > 0 {
							err1 := col.Update(bson.M{"_id": reqMsg.GetCommentId()},
								bson.M{"$inc": bson.M{"real_like_count": -1}})

							if err1 != nil {
								logs.Error(err1)

							}
						}
					}

					errLike := colLike.Remove(
						bson.M{"comment_id": commentId, "player_id": playId})

					if errLike != nil && errLike != mgo.ErrNotFound {
						logs.Error(errLike)
					}

					//点赞
				} else {

					if data.GetFakeLikeCount() > 0 {
						err1 := col.Update(
							bson.M{"_id": reqMsg.GetCommentId()},
							bson.M{"$inc": bson.M{"real_like_count": 1, "fake_like_count": 1}})

						if err1 != nil {
							logs.Error(err1)
						}
					} else {
						err1 := col.Update(
							bson.M{"_id": reqMsg.GetCommentId()},
							bson.M{"$inc": bson.M{"real_like_count": 1}})

						if err1 != nil {
							logs.Error(err1)
						}
					}

					nowTime := time.Now().Unix()

					//往点赞表中记录数据
					likeRecord := share_message.TableLikeRecord{
						CommentId:  easygo.NewInt64(commentId),
						PlayerId:   easygo.NewInt64(playId),
						LikeFlag:   easygo.NewBool(likeType),
						CreateTime: easygo.NewInt64(nowTime),
						ItemId:     easygo.NewInt64(data.GetItemId()),
					}
					_, errLike := colLike.Upsert(
						bson.M{"comment_id": commentId, "player_id": playId},
						bson.M{"$set": likeRecord})

					if errLike != nil && errLike != mgo.ErrNotFound {
						logs.Error(errLike)
					}
				}
			}
		}()

	}, reqMsg.GetCommentId(), reqMsg.GetLikeType(), common.GetUserId())

	return &share_message.LikeCommentResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(""),
	}

}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcRemoveItemFromStore(common *base.Common, reqMsg *share_message.ShopItemID) easygo.IMessage {
	logs.Info("===api RpcRemoveItemFromStore===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_STORE)
	defer closeFun()

	var playerStore share_message.TableShopPlayerStore = share_message.TableShopPlayerStore{}

	err1 := col.Find(bson.M{"player_id": common.GetUserId(), "item_id": reqMsg.GetItemId(), "store_type": 0}).One(&playerStore)

	if err1 != mgo.ErrNotFound && err1 != nil {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	_, err2 := col.RemoveAll(bson.M{"player_id": common.GetUserId(), "item_id": reqMsg.GetItemId()})

	if err2 != mgo.ErrNotFound && err2 != nil {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	//修改商品表中的收藏数
	easygo.Spawn(func(itemIdPara int64) {

		func() {

			col1, closeFun1 := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
			defer closeFun1()

			var item share_message.TableShopItem = share_message.TableShopItem{}

			err3 := col1.Find(bson.M{"_id": itemIdPara}).One(&item)
			if err3 != nil {
				logs.Error(err3)
			} else {
				if item.GetRealStoreCnt() > 0 {

					err4 := col1.Update(bson.M{"_id": itemIdPara},
						bson.M{"$inc": bson.M{"real_storeCnt": -1}})

					if err4 != nil {
						logs.Error(err4)
					}
				}
			}
		}()

	}, reqMsg.GetItemId())

	return &share_message.RemoveStoreResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(REMOVE_ITEM_FROM_STORE_SUCCESS)}
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcAddItemToStore(common *base.Common, reqMsg *share_message.ShopItemID) easygo.IMessage {
	logs.Info("===api RpcAddItemToStore===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	shopItem := ShopInstance.GetItemFromCache(reqMsg.GetItemId())

	if shopItem == nil {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", ADD_CART_NOT_SALE)

		return easygo.NewFailMsg(ADD_CART_NOT_SALE)
	}
	var itemDetail = ShopInstance.ItemDetail(shopItem)
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_STORE)
	defer closeFun()

	var timeNow int64 = time.Now().Unix()
	itemFile := share_message.ItemFile{}

	if itemDetail.ItemFiles != nil && len(itemDetail.ItemFiles) > 0 {
		itemFile = *itemDetail.ItemFiles[0]
	}

	var storeInfo *share_message.TableShopPlayerStore = &share_message.TableShopPlayerStore{
		PlayerId:       easygo.NewInt64(common.GetUserId()),
		ItemId:         easygo.NewInt64(reqMsg.GetItemId()),
		Name:           easygo.NewString(itemDetail.GetName()),
		Title:          easygo.NewString(itemDetail.GetTitle()),
		Price:          easygo.NewInt32(itemDetail.GetPrice()),
		ItemFile:       &itemFile,
		CreateTime:     easygo.NewInt64(timeNow),
		SellerPlayerId: easygo.NewInt64(itemDetail.GetPlayerId()),
		StoreType:      easygo.NewInt32(0),
	}

	_, e := col.Upsert(
		bson.M{"player_id": common.GetUserId(), "item_id": reqMsg.GetItemId(), "store_type": 0},
		bson.M{"$set": storeInfo})

	if e != mgo.ErrNotFound && e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	//修改商品表中的收藏数
	easygo.Spawn(func(itemIdPar int64) {

		func() {

			col1, closeFun1 := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
			defer closeFun1()

			err1 := col1.Update(
				bson.M{"_id": itemIdPar},
				bson.M{"$inc": bson.M{"real_storeCnt": 1}})

			if err1 != nil {
				logs.Error(err1)
			}
		}()

	}, reqMsg.GetItemId())

	return &share_message.AddStoreResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(ADD_STORE_SUCCESS),
		ItemId: easygo.NewInt64(reqMsg.GetItemId()),
	}
}

// 从购物车页面批量收藏商品
func (self *ServiceForHall) RpcBatchAddItemToStore(common *base.Common, reqMsg *share_message.ItemIdList) easygo.IMessage {
	logs.Info("===api RpcBatchAddItemToStore===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	if len(reqMsg.ItemIds) == 0 {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", BATCH_REMOVE_STORE_NOT_SELECT)

		return easygo.NewFailMsg(BATCH_REMOVE_STORE_NOT_SELECT)
	}

	//判断商品是否已经下架 并且取得商品的内容
	var storeInfos []*share_message.TableShopPlayerStore = []*share_message.TableShopPlayerStore{}
	for _, value := range reqMsg.GetItemIds() {
		shopItem := ShopInstance.GetItemFromCache(value)

		if shopItem == nil {

			SendToHallServerByApi(common.GetUserId(), "RpcToast", BATCH_REMOVE_STORE_NOT_EXIST)

			return easygo.NewFailMsg(BATCH_REMOVE_STORE_NOT_EXIST)
		}
		var itemDetail = ShopInstance.ItemDetail(shopItem)

		var timeNow int64 = time.Now().Unix()
		itemFile := share_message.ItemFile{}

		if itemDetail.ItemFiles != nil && len(itemDetail.ItemFiles) > 0 {
			itemFile = *itemDetail.ItemFiles[0]
		}

		var storeInfo *share_message.TableShopPlayerStore = &share_message.TableShopPlayerStore{
			PlayerId:       easygo.NewInt64(common.GetUserId()),
			ItemId:         easygo.NewInt64(value),
			Name:           easygo.NewString(itemDetail.GetName()),
			Title:          easygo.NewString(itemDetail.GetTitle()),
			Price:          easygo.NewInt32(itemDetail.GetPrice()),
			ItemFile:       &itemFile,
			CreateTime:     easygo.NewInt64(timeNow),
			SellerPlayerId: easygo.NewInt64(itemDetail.GetPlayerId()),
			StoreType:      easygo.NewInt32(0),
		}

		storeInfos = append(storeInfos, storeInfo)
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_STORE)
	defer closeFun()

	var playerStores []*share_message.TableShopPlayerStore = []*share_message.TableShopPlayerStore{}

	err := col.Find(bson.M{"player_id": common.GetUserId(), "store_type": 0}).All(&playerStores)

	if err != mgo.ErrNotFound && err != nil {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}
	//取得未搜藏的数据
	var needAddDatas []*share_message.TableShopPlayerStore = []*share_message.TableShopPlayerStore{}
	var needAddItemIds []int64 = []int64{}

	for _, reqValue := range storeInfos {
		var tempFlag bool = false
		if playerStores != nil && len(playerStores) >= 0 {
			for _, queryValue := range playerStores {
				if reqValue.GetItemId() == queryValue.GetItemId() {
					tempFlag = true
					break
				}
			}
		}
		if !tempFlag {
			needAddDatas = append(needAddDatas, reqValue)
			needAddItemIds = append(needAddItemIds, reqValue.GetItemId())
		}
	}

	var insLst []interface{}
	if nil != needAddDatas && len(needAddDatas) > 0 {
		for _, insValue := range needAddDatas {
			insLst = append(insLst, insValue)
		}
	}

	//批量插入
	if nil != insLst && len(insLst) > 0 {
		col.Insert(insLst...)
	}

	colCart, closeFunCart := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_CART)
	defer closeFunCart()
	//批量删除购物车
	_, errItem := colCart.RemoveAll(bson.M{"player_id": common.GetUserId(), "item_id": bson.M{"$in": reqMsg.GetItemIds()}})

	if errItem != mgo.ErrNotFound && errItem != nil {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	//修改商品表中的收藏数
	easygo.Spawn(func(itemIdPars []int64) {

		func() {

			col1, closeFun1 := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
			defer closeFun1()

			_, err1 := col1.UpdateAll(
				bson.M{"_id": bson.M{"$in": itemIdPars}},
				bson.M{"$inc": bson.M{"real_storeCnt": 1}})

			if err1 != nil {
				logs.Error(err1)
			}
		}()

	}, needAddItemIds)

	return &share_message.BatchAddStoreResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(BATCH_REMOVE_STORE_SUCCESS)}
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcStoreInfo(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===api RpcStoreInfo===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	storeItems := []*share_message.ShopItem{}
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_STORE)
	defer closeFun()
	var list []*share_message.TableShopPlayerStore

	e := col.Find(bson.M{"player_id": common.GetUserId()}).Sort("-create_time").All(&list)

	if e != nil && e != mgo.ErrNotFound {
		logs.Error(e)
	}

	if e == nil && list != nil {

		var itemIds []int64 = []int64{}
		for _, value := range list {
			itemIds = append(itemIds, value.GetItemId())
		}

		var deletedItems []*share_message.TableShopItem

		// 判断是否有删除的商品
		colItem, closeFunItem := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
		defer closeFunItem()
		errItem := colItem.Find(bson.M{"_id": bson.M{"$in": itemIds}, "state": for_game.SHOP_ITEM_DELETE}).All(&deletedItems)

		if errItem != nil && errItem != mgo.ErrNotFound {
			logs.Error(e)
		}

		for _, value := range list {

			var doFlag bool = false
			if errItem == nil && len(deletedItems) > 0 {
				for _, delValue := range deletedItems {
					if value.GetItemId() == delValue.GetItemId() {
						doFlag = true
						break
					}
				}
			}

			if doFlag {
				continue
			}

			shopItem := ShopInstance.GetItemFromCache(*value.ItemId)
			//判断是否下架
			var itemFlag bool
			if nil != shopItem {
				itemFlag = false
			} else {
				itemFlag = true
			}
			var types int32
			if p := for_game.GetRedisPlayerBase(value.GetPlayerId()); p != nil {
				types = p.GetTypes()
			}
			newItem := share_message.ShopItem{
				ItemId:   easygo.NewInt64(value.GetItemId()),
				Name:     easygo.NewString(value.GetName()),
				Title:    easygo.NewString(value.GetTitle()),
				ItemFile: value.ItemFile,
				Price:    easygo.NewInt32(value.GetPrice()),
				ItemFlag: easygo.NewBool(itemFlag),
				PlayerId: easygo.NewInt64(value.GetSellerPlayerId()),
				Types:    easygo.NewInt32(types),
			}
			storeItems = append(storeItems, &newItem)
		}
	}
	return &share_message.StoreItemList{
		Result:     easygo.NewInt32(0),
		Msg:        easygo.NewString(""),
		StoreItems: storeItems,
	}
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcCartInfo(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===api RpcCartInfo===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	cartItemInfos := []*share_message.CartItemInfo{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_CART)
	defer closeFun()
	var list []*share_message.TablePlayerCart

	e := col.Find(bson.M{"player_id": common.GetUserId()}).Sort("create_time").All(&list)

	var tempSellerMap map[int64]string = map[int64]string{}

	if e == nil && list != nil {
		for _, value := range list {
			tempSellerMap[value.GetSellerPlayerId()] = value.GetSellerNickName()
		}
	}

	if e == nil && list != nil {
		for sellPlayId, sellNickName := range tempSellerMap {
			var cartItemInfo share_message.CartItemInfo
			cartItemList := []*share_message.CartItem{}
			for _, value := range list {

				if sellPlayId == value.GetSellerPlayerId() {

					var stockErr string  //库存不足错误内容
					var noSaleErr string //下架错误内容
					var blackErr string  //黑名单错误内容

					shopItem := ShopInstance.GetItemFromCache(*value.ItemId)
					//先判断库存库存不足判断
					var flag int32
					if nil != shopItem {
						if shopItem.stock_count <= 0 {
							//if shopItem.stock_count-shopItem.lock_count <= 0 {
							flag = 1
							stockErr = QUERY_CART_NOT_STOCK_WARN
						} else {
							// 黑名单判断
							var blackList []PLAYER_ID = ShopInstance.GetBlackLists(common.GetUserId())
							var blackFlag int32 = 0
							for _, black := range blackList {
								if shopItem.player_id == black {
									blackFlag = 1
									break
								}
							}
							//显示黑名单
							if blackFlag == 1 {
								flag = 1
								blackErr = QUERY_CART_NOT_BLACK_WARN
							}
						}
					} else {
						//下架判断
						flag = 1
						noSaleErr = QUERY_CART_NOT_SALE_WARN
					}

					var errContent string
					if stockErr != "" {
						errContent = stockErr
					} else if noSaleErr != "" {
						errContent = noSaleErr
					} else if blackErr != "" {
						errContent = blackErr
					}

					newItem := share_message.CartItem{
						ItemId:     easygo.NewInt64(value.GetItemId()),
						Title:      easygo.NewString(value.GetTitle()),
						ItemFile:   value.ItemFile,
						Price:      easygo.NewInt32(value.GetPrice()),
						AddCount:   easygo.NewInt32(value.GetAddCount()),
						Name:       easygo.NewString(value.GetName()),
						Flag:       easygo.NewInt32(flag),
						ErrContent: easygo.NewString(errContent),
					}
					cartItemList = append(cartItemList, &newItem)

				}
			}

			if cartItemList != nil && len(cartItemList) > 0 {

				cartItemInfo = share_message.CartItemInfo{
					CartItems:      cartItemList,
					SellerPlayerId: easygo.NewInt64(sellPlayId),
					SellerNickName: easygo.NewString(sellNickName),
				}

				cartItemInfos = append(cartItemInfos, &cartItemInfo)
			}
		}
	}

	return &share_message.CartItemInfoList{CartItemInfos: cartItemInfos}
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcAddItemToCart(common *base.Common, reqMsg *share_message.ShopItemID) easygo.IMessage {
	logs.Info("===api RpcAddItemToCart===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	shopItem := ShopInstance.GetItemFromCache(reqMsg.GetItemId())

	if shopItem == nil {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", ADD_CART_NOT_SALE)

		return easygo.NewFailMsg(ADD_CART_NOT_SALE)
	}
	var itemDetail = ShopInstance.ItemDetail(shopItem)
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_CART)
	defer closeFun()

	var cartInfo *share_message.TablePlayerCart = &share_message.TablePlayerCart{}

	e := col.Find(bson.M{
		"player_id": common.GetUserId(),
		"item_id":   reqMsg.GetItemId()}).One(cartInfo)

	if e != mgo.ErrNotFound && e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	if e != mgo.ErrNotFound {
		var addCount int32 = 1
		addCount = addCount + cartInfo.GetAddCount()

		//判断加购的数量是否超过了库存数量
		if addCount > shopItem.stock_count {
			SendToHallServerByApi(common.GetUserId(), "RpcToast", ADD_CART_COUNT_OVER)

			return easygo.NewFailMsg(ADD_CART_COUNT_OVER)
		}

		e := col.Update(
			bson.M{"player_id": common.GetUserId(), "item_id": reqMsg.GetItemId()},
			bson.M{"$inc": bson.M{"add_count": 1}})

		if e == mgo.ErrNotFound {
			SendToHallServerByApi(common.GetUserId(), "RpcToast", ADD_CART_NOT_EXIST)

			return easygo.NewFailMsg(ADD_CART_NOT_EXIST)
		}

		if e != nil {
			logs.Error(e)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)
		}

		return &share_message.AddCartResult{
			Result:   easygo.NewInt32(0),
			Msg:      easygo.NewString(ADD_CART_SUCCESS),
			ItemId:   easygo.NewInt64(reqMsg.GetItemId()),
			AddCount: easygo.NewInt32(addCount)}

	} else {
		var timeNow int64 = time.Now().Unix()
		itemFile := share_message.ItemFile{}

		if itemDetail.ItemFiles != nil && len(itemDetail.ItemFiles) > 0 {
			itemFile = *itemDetail.ItemFiles[0]
		}
		newItem := &share_message.TablePlayerCart{
			PlayerId:   easygo.NewInt64(common.GetUserId()),
			ItemId:     easygo.NewInt64(reqMsg.GetItemId()),
			Title:      easygo.NewString(itemDetail.GetTitle()),
			Price:      easygo.NewInt32(itemDetail.Price),
			ItemFile:   &itemFile,
			AddCount:   easygo.NewInt32(1),
			CreateTime: easygo.NewInt64(timeNow),
			//OriginPrice: ItemDetail.,
			Name:           easygo.NewString(itemDetail.GetName()),
			SellerPlayerId: easygo.NewInt64(itemDetail.GetPlayerId()),
			SellerNickName: easygo.NewString(itemDetail.GetNickname()),
		}

		e := col.Insert(newItem)

		if e != nil {
			logs.Error(e)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)
		}

		return &share_message.AddCartResult{
			Result:   easygo.NewInt32(0),
			Msg:      easygo.NewString(ADD_CART_SUCCESS),
			ItemId:   easygo.NewInt64(reqMsg.GetItemId()),
			AddCount: easygo.NewInt32(1)}
	}
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcSubItemToCart(common *base.Common, reqMsg *share_message.ShopItemID) easygo.IMessage {
	logs.Info("===api RpcSubItemToCart===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	shopItem := ShopInstance.GetItemFromCache(reqMsg.GetItemId())

	if shopItem == nil {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", SUB_CART_NOT_SALE)

		return easygo.NewFailMsg(SUB_CART_NOT_SALE)
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_CART)
	defer closeFun()

	var cartInfo *share_message.TablePlayerCart = &share_message.TablePlayerCart{}

	e := col.Find(bson.M{"player_id": common.GetUserId(), "item_id": reqMsg.GetItemId()}).One(cartInfo)

	if e == mgo.ErrNotFound || cartInfo == nil {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", SUB_CART_NOT_EXIST)

		return easygo.NewFailMsg(SUB_CART_NOT_EXIST)
	}
	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	var addCount int32 = 1
	if cartInfo != nil && cartInfo.GetAddCount() > 1 {
		addCount = cartInfo.GetAddCount() - 1

		e := col.Update(
			bson.M{"player_id": common.GetUserId(), "item_id": reqMsg.GetItemId()},
			bson.M{"$inc": bson.M{"add_count": -1}})

		if e == mgo.ErrNotFound {
			SendToHallServerByApi(common.GetUserId(), "RpcToast", SUB_CART_NOT_EXIST)

			return easygo.NewFailMsg(SUB_CART_NOT_EXIST)
		}

		if e != nil {
			logs.Error(e)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)
		}
	}

	return &share_message.AddCartResult{
		Result:   easygo.NewInt32(0),
		Msg:      easygo.NewString(SUB_CART_SUCCESS),
		ItemId:   easygo.NewInt64(reqMsg.GetItemId()),
		AddCount: easygo.NewInt32(addCount)}
}

// 接收登陆服务器发送过来消息
func (self *ServiceForHall) RpcRemoveItemFromCart(common *base.Common, reqMsg *share_message.ItemIdList) easygo.IMessage {
	logs.Info("===api RpcRemoveItemFromCart===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	if len(reqMsg.ItemIds) == 0 {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", REMOVE_ITEM_FROM_CART_NOT_EXIST)

		return easygo.NewFailMsg(REMOVE_ITEM_FROM_CART_NOT_EXIST)
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_CART)
	defer closeFun()

	_, err := col.RemoveAll(bson.M{"player_id": common.GetUserId(), "item_id": bson.M{"$in": reqMsg.GetItemIds()}})

	if err != mgo.ErrNotFound && err != nil {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	return &share_message.RemoveCartResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(REMOVE_ITEM_FROM_CART_SUCCESS)}
}

/*
func (self *ServiceForHall) RpcCreateOrder(common *base.Common, reqMsg *share_message.BuyItemInfo) easygo.IMessage {
	logs.Info("===api RpcCreateOrder===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	if len(reqMsg.Items) <= 0 {
		logs.Info("没有买任何商品")
		SendToHallServerByApi(common.GetUserId(), "RpcToast", CREATE_ORDER_BUY_NULL)

		return easygo.NewFailMsg(CREATE_ORDER_BUY_NULL)
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()

	for _, value := range reqMsg.Items {

		var item share_message.TableShopItem = share_message.TableShopItem{}
		e := col.Find(bson.M{"_id": value.GetItemId(), "state": for_game.SHOP_ITEM_SALE}).One(&item)

		if e == mgo.ErrNotFound {
			s := fmt.Sprintf(CREATE_ORDER_ITEM_NO_SALE, value.GetItemName())
			logs.Info(s)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", s)

			return easygo.NewFailMsg(s)
		}

		if e != nil {
			logs.Error(e)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)
		}

		//库存
		if item.GetStockCount() < value.GetCount() {

			logs.Info("商品(%s)：库存不足,还剩%v库存",
				item.GetName(),
				easygo.AnytoA(int64(item.GetStockCount())))

			SendToHallServerByApi(common.GetUserId(), "RpcToast", fmt.Sprintf(CREATE_ORDER_ITEM_NO_STOCK,
				item.GetName(),
				easygo.AnytoA(int64(item.GetStockCount()))))

			return easygo.NewFailMsg(fmt.Sprintf(CREATE_ORDER_ITEM_NO_STOCK,
				item.GetName(),
				easygo.AnytoA(int64(item.GetStockCount()))))
		}

		//黑名单
		var blackList []PLAYER_ID = ShopInstance.GetBlackLists(common.GetUserId())
		var blackFlag int32 = 0

		for _, black := range blackList {
			if item.GetPlayerId() == black {
				blackFlag = 1
				break
			}
		}
		//显示黑名单
		if blackFlag == 1 {
			SendToHallServerByApi(common.GetUserId(), "RpcToast", fmt.Sprintf(CREATE_ORDER_ITEM_BLACK, item.GetName()))

			return easygo.NewFailMsg(fmt.Sprintf(CREATE_ORDER_ITEM_BLACK, item.GetName()))
		}
	}

	now := time.Now().Unix()
	var state int32 = for_game.SHOP_ORDER_WAIT_PAY
	var deleteBuy int32 = 0
	var deleteSell int32 = 0
	var delayReceive int32 = 0
	var totalPrice int32 = 0
	bill := share_message.TableBill{Price: &totalPrice}

	colOrder, closeFunOrder := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFunOrder()
	who := for_game.GetRedisPlayerBase(common.GetUserId())
	for _, value := range reqMsg.Items {

		shopItem := ShopInstance.GetItemFromCache(value.GetItemId())

		if shopItem != nil {

			orderId := ShopInstance.CreateOrderID()

			itemFile := share_message.ItemFile{
				FileUrl:  &shopItem.item_files[0].file_url,
				FileType: &shopItem.item_files[0].file_type}

			item := share_message.ShopOrderItem{
				ItemId:   value.ItemId,
				Name:     &shopItem.name,
				Price:    &shopItem.price,
				ItemFile: &itemFile,
				Count:    value.Count,
				Title:    &shopItem.title}

			order := share_message.TableShopOrder{
				OrderId:          easygo.NewInt64(orderId),
				SponsorId:        easygo.NewInt64(shopItem.player_id),
				SponsorSex:       easygo.NewInt32(shopItem.sex),
				SponsorAvatar:    easygo.NewString(shopItem.avatar),
				SponsorNickname:  easygo.NewString(shopItem.nickname),
				ReceiverId:       easygo.NewInt64(common.GetUserId()),
				ReceiverSex:      easygo.NewInt32(who.GetSex()),
				ReceiverNickname: easygo.NewString(who.GetNickName()),
				ReceiverAvatar:   easygo.NewString(who.GetHeadIcon()),
				Items:            &item,
				State:            easygo.NewInt32(state),
				DelayReceive:     easygo.NewInt32(delayReceive),
				DeleteBuy:        easygo.NewInt32(deleteBuy),
				DeleteSell:       easygo.NewInt32(deleteSell),
				DeliverAddress: &share_message.DeliverAddress{
					Name:          easygo.NewString(shopItem.userName),
					Phone:         easygo.NewString(shopItem.phone),
					Region:        easygo.NewString(shopItem.address),
					DetailAddress: easygo.NewString(shopItem.detail_address),
				},
				ReceiveAddress: &share_message.ReceiveAddress{
					Name:          easygo.NewString(reqMsg.Address.GetName()),
					Phone:         easygo.NewString(reqMsg.Address.GetPhone()),
					Region:        easygo.NewString(reqMsg.Address.GetRegion()),
					DetailAddress: easygo.NewString(reqMsg.Address.GetDetailAddress()),
				},
				CreateTime:         easygo.NewInt64(now),
				ReceiveTime:        easygo.NewInt64(now),
				Remark:             easygo.NewString(value.GetRemark()),
				SponsorAccount:     easygo.NewString(shopItem.account),
				ReceiverAccount:    easygo.NewString(who.GetAccount()),
				ReceiverNotifyFlag: easygo.NewBool(true),
				SponsorNotifyFlag:  easygo.NewBool(true),
				UpdateTime:         easygo.NewInt64(now),
			}

			e := colOrder.Insert(order)

			if e != nil {
				logs.Error(e)
				SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

				return easygo.NewFailMsg(DATABASE_ERROR)
			}
			totalPrice += shopItem.price * value.GetCount()

			//创建订单就通知
			//订单通知

			easygo.Spawn(func(orderParam share_message.TableShopOrder) {
				if &orderParam != nil {

					var content string = MESSAGE_TO_SELLER_NEW
					typeValue := share_message.BuySell_Type_Seller

					//商城消息通知
					ShopInstance.InsMessageNotify(
						easygo.NewString(content),
						&typeValue,
						&orderParam)

					//商城极光推送
					var jgContent string = MESSAGE_TO_SELLER_NEW_PUSH
					ShopInstance.JGMessageNotify(jgContent, orderParam.GetSponsorId(), orderParam.GetOrderId(), typeValue)

					//创建订单通知买家 商城订单红点推送
					SendToPlayer(orderParam.GetReceiverId(), "RpcShopOrderNotify",
						&share_message.ShopOrderNotifyInfoWithWho{
							PlayerId: easygo.NewInt64(orderParam.GetReceiverId()),
							OrderId:  easygo.NewInt64(orderParam.GetOrderId()),
						})
					//创建订单通知卖家 商城订单红点推送
					SendToPlayer(orderParam.GetSponsorId(), "RpcShopOrderNotify",
						&share_message.ShopOrderNotifyInfoWithWho{
							PlayerId: easygo.NewInt64(orderParam.GetSponsorId()),
							OrderId:  easygo.NewInt64(orderParam.GetOrderId()),
						})

				} else {
					logs.Debug("创建订单发通知,缺少订单")
				}
			}, order)

			bill.OrderList = append(bill.OrderList, orderId)

		}
	}
	orderId := ShopInstance.CreateOrderID()
	bill.OrderId = &orderId
	bill.State = easygo.NewInt32(for_game.SHOP_ORDER_WAIT_PAY)

	fun := func() {
		SaveDataToDBForBills(bill, 0)
	}
	easygo.AfterFunc(0, fun)

	var tempItemIds []int64 = []int64{}
	for _, value := range reqMsg.Items {
		tempItemIds = append(tempItemIds, value.GetItemId())
	}

	easygo.Spawn(func(itemIds []int64) {

		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_CART)
		_, err := col.RemoveAll(
			bson.M{"player_id": common.GetUserId(),
				"item_id": bson.M{"$in": itemIds}})
		closeFun()

		if err != nil && err != mgo.ErrNotFound {
			logs.Error(err)
		}
	}, tempItemIds)

	logs.Info("创建订单成功")
	// easygo.Spawn(func() {
	// 	for_game.SetRedisOperationChannelReportFildVal(util.GetMilliTime(), 1, who.GetChannel(), "ShopOrderCount") //渠道汇总报表添加下单数量
	// })
	return &share_message.BuyItemResult{
		Result:  easygo.NewInt32(0),
		Msg:     easygo.NewString(CREATE_ORDER_SUCCESS),
		OrderId: &orderId}
}
*/

//生成订单
func (self *ServiceForHall) RpcCreateOrder(common *base.Common, reqMsg *share_message.BuyItemInfo) easygo.IMessage {
	logs.Info("RpcCreateOrder,msg=%v", reqMsg) // 别删，永久留存

	if len(reqMsg.Items) <= 0 {
		logs.Info("没有买任何商品")
		SendToHallServerByApi(common.GetUserId(), "RpcToast", CREATE_ORDER_BUY_NULL)
		return easygo.NewFailMsg(CREATE_ORDER_BUY_NULL)
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()
	who := for_game.GetRedisPlayerBase(common.GetUserId())
	for _, value := range reqMsg.Items {

		var item share_message.TableShopItem = share_message.TableShopItem{}
		e := col.Find(bson.M{"_id": value.GetItemId(), "state": for_game.SHOP_ITEM_SALE}).One(&item)

		if e == mgo.ErrNotFound {
			s := fmt.Sprintf(CREATE_ORDER_ITEM_NO_SALE, value.GetItemName())
			logs.Info(s)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", s)
			return easygo.NewFailMsg(s)
		}

		if e != nil {
			logs.Error(e)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)
			return easygo.NewFailMsg(DATABASE_ERROR)
		}

		//库存
		if item.GetStockCount() < value.GetCount() {

			logs.Info("商品(%s)：库存不足,还剩%v库存",
				value.GetItemName(),
				easygo.AnytoA(int64(item.GetStockCount())))

			SendToHallServerByApi(common.GetUserId(), "RpcToast", fmt.Sprintf(CREATE_ORDER_ITEM_NO_STOCK,
				value.GetItemName(),
				easygo.AnytoA(int64(item.GetStockCount()))))
			return easygo.NewFailMsg(fmt.Sprintf(CREATE_ORDER_ITEM_NO_STOCK,
				value.GetItemName(),
				easygo.AnytoA(int64(item.GetStockCount()))))
		}

		//如果是点卡 再判断一次实际导入库中的库存是否足够
		if item.GetType() != nil && item.GetType().GetType() == for_game.SHOP_POINT_CARD_CATEGORY {

			//通过物品的个数取得卡密个数
			pointCardList := for_game.GetPointCardByBuyInfos(item.GetPlayerAccount(), item.GetPointCardName(), value.GetCount())
			if nil != pointCardList && (len(pointCardList) == 0 || int32(len(pointCardList)) < value.GetCount()) {
				logs.Info("商品(%s)：库存不足,还剩%v库存",
					item.GetPointCardName(),
					easygo.AnytoA(int64(len(pointCardList))))

				SendToHallServerByApi(common.GetUserId(), "RpcToast", fmt.Sprintf(CREATE_ORDER_ITEM_NO_STOCK,
					item.GetPointCardName(),
					easygo.AnytoA(int64(len(pointCardList)))))
				return easygo.NewFailMsg(fmt.Sprintf(CREATE_ORDER_ITEM_NO_STOCK,
					item.GetPointCardName(),
					easygo.AnytoA(int64(len(pointCardList)))))
			}
		}
		//黑名单
		var blackList []PLAYER_ID = ShopInstance.GetBlackLists(common.GetUserId())
		var blackFlag int32 = 0

		for _, black := range blackList {
			if item.GetPlayerId() == black {
				blackFlag = 1
				break
			}
		}
		//显示黑名单
		if blackFlag == 1 {
			SendToHallServerByApi(common.GetUserId(), "RpcToast", fmt.Sprintf(CREATE_ORDER_ITEM_BLACK, item.GetName()))
			return easygo.NewFailMsg(fmt.Sprintf(CREATE_ORDER_ITEM_BLACK, item.GetName()))
		}
	}

	now := time.Now().Unix()
	var state int32 = for_game.SHOP_ORDER_WAIT_PAY
	var deleteBuy int32 = 0
	var deleteSell int32 = 0
	var delayReceive int32 = 0
	var totalPrice int32 = 0
	bill := share_message.TableBill{Price: &totalPrice}

	colOrder, closeFunOrder := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFunOrder()

	for _, value := range reqMsg.Items {

		errStr, shopItem := ShopInstance.GetShopItem(share_message.BuySell_Type_Seller, value.GetItemId())
		if errStr != "" {
			logs.Error(errStr)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", errStr)
			return easygo.NewFailMsg(errStr)
		}

		if nil == shopItem {
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DETAIL_SHOP_ITEM_NOT_EXIST)
			return easygo.NewFailMsg(DETAIL_SHOP_ITEM_NOT_EXIST)
		}

		if shopItem != nil {

			orderId := ShopInstance.CreateOrderID()

			itemFile := share_message.ItemFile{
				FileUrl:  &shopItem.item_files[0].file_url,
				FileType: &shopItem.item_files[0].file_type}

			item := share_message.ShopOrderItem{
				ItemId:        easygo.NewInt64(value.GetItemId()),
				Name:          easygo.NewString(shopItem.name),
				Price:         easygo.NewInt32(shopItem.price),
				ItemFile:      &itemFile,
				Count:         easygo.NewInt32(value.GetCount()),
				Title:         easygo.NewString(shopItem.title),
				PointCardName: easygo.NewString(shopItem.pointCardName),
				CopyName:      easygo.NewString(shopItem.name),
				ItemType:      easygo.NewInt32(shopItem.item_type),
			}

			//点卡的时候,为了客户端不修改代码,把商品名称设置为点卡名称
			if shopItem.item_type == for_game.SHOP_POINT_CARD_CATEGORY {
				item.Name = easygo.NewString(shopItem.pointCardName)
			}

			order := share_message.TableShopOrder{
				OrderId:          easygo.NewInt64(orderId),
				SponsorId:        easygo.NewInt64(shopItem.player_id),
				SponsorSex:       easygo.NewInt32(shopItem.sex),
				SponsorAvatar:    easygo.NewString(shopItem.avatar),
				SponsorNickname:  easygo.NewString(shopItem.nickname),
				ReceiverId:       easygo.NewInt64(who.GetPlayerId()),
				ReceiverSex:      easygo.NewInt32(who.GetSex()),
				ReceiverNickname: easygo.NewString(who.GetNickName()),
				ReceiverAvatar:   easygo.NewString(who.GetHeadIcon()),
				Items:            &item,
				State:            easygo.NewInt32(state),
				DelayReceive:     easygo.NewInt32(delayReceive),
				DeleteBuy:        easygo.NewInt32(deleteBuy),
				DeleteSell:       easygo.NewInt32(deleteSell),
				DeliverAddress: &share_message.DeliverAddress{
					Name:          easygo.NewString(shopItem.userName),
					Phone:         easygo.NewString(shopItem.phone),
					Region:        easygo.NewString(shopItem.address),
					DetailAddress: easygo.NewString(shopItem.detail_address),
				},
				ReceiveAddress: &share_message.ReceiveAddress{
					Name:          easygo.NewString(reqMsg.Address.GetName()),
					Phone:         easygo.NewString(reqMsg.Address.GetPhone()),
					Region:        easygo.NewString(reqMsg.Address.GetRegion()),
					DetailAddress: easygo.NewString(reqMsg.Address.GetDetailAddress()),
				},
				CreateTime:         easygo.NewInt64(now),
				ReceiveTime:        easygo.NewInt64(now),
				Remark:             easygo.NewString(value.GetRemark()),
				SponsorAccount:     easygo.NewString(shopItem.account),
				ReceiverAccount:    easygo.NewString(who.GetAccount()),
				ReceiverNotifyFlag: easygo.NewBool(true),
				SponsorNotifyFlag:  easygo.NewBool(true),
				UpdateTime:         easygo.NewInt64(now),
			}

			e := colOrder.Insert(order)

			if e != nil {
				logs.Error(e)
				SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)
				return easygo.NewFailMsg(DATABASE_ERROR)
			}
			totalPrice += shopItem.price * value.GetCount()

			//创建订单就通知
			//订单通知

			easygo.Spawn(func(orderParam share_message.TableShopOrder) {
				if &orderParam != nil {

					var content string = MESSAGE_TO_SELLER_NEW
					typeValue := share_message.BuySell_Type_Seller

					//商城消息通知
					ShopInstance.InsMessageNotify(easygo.NewString(content), &typeValue, &orderParam)

					//商城极光推送
					var jgContent string = MESSAGE_TO_SELLER_NEW_PUSH
					ShopInstance.JGMessageNotify(jgContent, orderParam.GetSponsorId(), orderParam.GetOrderId(), typeValue)

					//创建订单通知买家 商城订单红点推送
					logs.Info("创建订单通知买家 商城订单红点推送")
					SendMsgToHallClientNew([]int64{orderParam.GetReceiverId()}, "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
						OrderId: easygo.NewInt64(orderParam.GetOrderId())})
					/*		SendToPlayer(orderParam.GetReceiverId(), "RpcShopOrderNotify",
							&share_message.ShopOrderNotifyInfoWithWho{
								PlayerId: easygo.NewInt64(orderParam.GetReceiverId()),
								OrderId:  easygo.NewInt64(orderParam.GetOrderId()),
							})*/
					//创建订单通知卖家 商城订单红点推送
					SendMsgToHallClientNew([]int64{orderParam.GetSponsorId()}, "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
						OrderId: easygo.NewInt64(orderParam.GetOrderId())})
					/*		SendToPlayer(orderParam.GetSponsorId(), "RpcShopOrderNotify",
							&share_message.ShopOrderNotifyInfoWithWho{
								PlayerId: easygo.NewInt64(orderParam.GetSponsorId()),
								OrderId:  easygo.NewInt64(orderParam.GetOrderId()),
							})*/

				} else {
					logs.Debug("创建订单发通知,缺少订单")
				}
			}, order)

			bill.OrderList = append(bill.OrderList, orderId)

		}
	}
	orderId := ShopInstance.CreateOrderID()
	bill.OrderId = &orderId
	bill.State = easygo.NewInt32(for_game.SHOP_ORDER_WAIT_PAY)

	colBill, closeFunBill := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_BILLS)
	defer closeFunBill()

	_, eBill := colBill.Upsert(
		bson.M{"_id": bill.GetOrderId()},
		bson.M{"$set": bill})

	if eBill != nil {
		logs.Error(eBill)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)
		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	var tempItemIds []int64 = []int64{}
	for _, value := range reqMsg.Items {
		tempItemIds = append(tempItemIds, value.GetItemId())
	}

	easygo.Spawn(func(itemIds []int64) {

		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_CART)
		_, err := col.RemoveAll(
			bson.M{"player_id": common.GetUserId(),
				"item_id": bson.M{"$in": itemIds}})
		closeFun()

		if err != nil && err != mgo.ErrNotFound {
			logs.Error(err)
		}
	}, tempItemIds)

	logs.Info("创建订单成功")
	// easygo.Spawn(func() {
	// 	for_game.SetRedisOperationChannelReportFildVal(util.GetMilliTime(), 1, who.GetChannel(), "ShopOrderCount") //渠道汇总报表添加下单数量
	// })
	return &share_message.BuyItemResult{
		Result:  easygo.NewInt32(0),
		Msg:     easygo.NewString(CREATE_ORDER_SUCCESS),
		OrderId: &orderId}
}

func (self *ServiceForHall) RpcOrderList(common *base.Common, reqMsg *share_message.OrderInfo) easygo.IMessage {
	logs.Info("===api RpcOrderList===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	var count int32 = 0
	page := reqMsg.GetPage()
	pageSize := reqMsg.GetPageSize()

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()

	list := []*share_message.TableShopOrder{}

	selectString := make(bson.M)

	if reqMsg.GetItemType() == share_message.SearchOrder_Con_ALL_ORDER {

		if reqMsg.GetType() == share_message.BuySell_Type_Seller {
			selectString = bson.M{
				"sponsor_id":  common.GetUserId(),
				"delete_sell": 0}
		} else {
			selectString = bson.M{
				"receiver_id": common.GetUserId(),
				"delete_buy":  0}
		}

	} else if reqMsg.GetItemType() == share_message.SearchOrder_Con_WAIT_PAY {
		if reqMsg.GetType() == share_message.BuySell_Type_Seller {

			selectString = bson.M{
				"sponsor_id":  common.GetUserId(),
				"state":       for_game.SHOP_ORDER_WAIT_PAY,
				"delete_sell": 0}
		} else {

			selectString = bson.M{
				"receiver_id": common.GetUserId(),
				"state":       for_game.SHOP_ORDER_WAIT_PAY,
				"delete_buy":  0}
		}
	} else if reqMsg.GetItemType() == share_message.SearchOrder_Con_WAIT_SEND {
		if reqMsg.GetType() == share_message.BuySell_Type_Seller {

			selectString = bson.M{
				"sponsor_id":  common.GetUserId(),
				"state":       for_game.SHOP_ORDER_WAIT_SEND,
				"delete_sell": 0}
		} else {

			selectString = bson.M{
				"receiver_id": common.GetUserId(),
				"state":       for_game.SHOP_ORDER_WAIT_SEND,
				"delete_buy":  0}
		}
	} else if reqMsg.GetItemType() == share_message.SearchOrder_Con_WAIT_RECEIVE {
		if reqMsg.GetType() == share_message.BuySell_Type_Seller {

			selectString = bson.M{
				"sponsor_id":  common.GetUserId(),
				"state":       for_game.SHOP_ORDER_WAIT_RECEIVE,
				"delete_sell": 0}
		} else {

			selectString = bson.M{
				"receiver_id": common.GetUserId(),
				"state":       for_game.SHOP_ORDER_WAIT_RECEIVE,
				"delete_buy":  0}
		}
	} else if reqMsg.GetItemType() == share_message.SearchOrder_Con_FINISH_ORDER {
		if reqMsg.GetType() == share_message.BuySell_Type_Seller {

			selectString = bson.M{
				"sponsor_id":  common.GetUserId(),
				"delete_sell": 0,
				"$or": []bson.M{
					bson.M{"state": for_game.SHOP_ORDER_FINISH},
					bson.M{"state": for_game.SHOP_ORDER_EVALUTE},
				}}
		} else {

			selectString = bson.M{
				"receiver_id": common.GetUserId(),
				"delete_buy":  0,
				"$or": []bson.M{
					bson.M{"state": for_game.SHOP_ORDER_FINISH},
					bson.M{"state": for_game.SHOP_ORDER_EVALUTE},
				}}
		}
	}

	query := col.Find(selectString)
	cnt, errVa := query.Count()
	count = int32(cnt)

	if errVa != nil {

		logs.Error(errVa)

		return &share_message.OrderItemList{
			Page:     easygo.NewInt32(page),
			PageSize: easygo.NewInt32(pageSize),
			Count:    easygo.NewInt32(count)}
	}

	e := query.Limit(int(pageSize)).Skip(int(page * pageSize)).Sort("-update_time").All(&list)

	if e != nil {

		logs.Error(e)

		return &share_message.OrderItemList{
			Page:     easygo.NewInt32(page),
			PageSize: easygo.NewInt32(pageSize),
			Count:    easygo.NewInt32(count)}
	}
	orders := []*share_message.OrderItem{}

	for _, value := range list {

		var avatar string
		var nickName string
		var sex int32

		if reqMsg.GetType() == share_message.BuySell_Type_Seller {
			avatar = value.GetReceiverAvatar()
			nickName = value.GetReceiverNickname()
			sex = value.GetReceiverSex()
		} else {
			avatar = value.GetSponsorAvatar()
			nickName = value.GetSponsorNickname()
			sex = value.GetSponsorSex()
		}

		item := share_message.BuyItem{
			ItemId:    easygo.NewInt64(value.Items.GetItemId()),
			ItemFile:  value.Items.ItemFile,
			Name:      easygo.NewString(value.Items.GetName()),
			Price:     easygo.NewInt32(value.Items.GetPrice()),
			SponsorId: easygo.NewInt64(value.GetSponsorId()),
			Count:     easygo.NewInt32(value.Items.GetCount()),
			Avatar:    easygo.NewString(avatar),
			Nickname:  easygo.NewString(nickName),
			Sex:       easygo.NewInt32(sex),
			Title:     easygo.NewString(value.Items.GetTitle()),
			DeliverAddress: &share_message.DeliverAddress{
				Name:          easygo.NewString(value.DeliverAddress.GetName()),
				Phone:         easygo.NewString(value.DeliverAddress.GetPhone()),
				Region:        easygo.NewString(value.DeliverAddress.GetRegion()),
				DetailAddress: easygo.NewString(value.DeliverAddress.GetDetailAddress()),
			},
		}

		orderItem := share_message.OrderItem{
			OrderId: easygo.NewInt64(value.GetOrderId()),
			Item:    &item,
			State:   easygo.NewInt32(value.GetState()),
			Address: &share_message.ReceiveAddress{
				DetailAddress: easygo.NewString(value.ReceiveAddress.GetDetailAddress()),
				Name:          easygo.NewString(value.ReceiveAddress.GetName()),
				Region:        easygo.NewString(value.ReceiveAddress.GetRegion()),
				Phone:         easygo.NewString(value.ReceiveAddress.GetPhone())},
			ExpressCode:   easygo.NewString(value.GetExpressCode()),
			ExpressCom:    easygo.NewString(value.GetExpressCom()),
			CreateTime:    easygo.NewInt64(value.GetCreateTime()),
			ServerNowTime: easygo.NewInt64(time.Now().Unix()),
			DelayReceive:  easygo.NewInt32(value.GetDelayReceive()),
		}

		orders = append(orders, &orderItem)
	}
	return &share_message.OrderItemList{
		Items:    orders,
		Page:     easygo.NewInt32(page),
		PageSize: easygo.NewInt32(pageSize),
		Count:    easygo.NewInt32(count)}
}

func (self *ServiceForHall) RpcOrderInfo(common *base.Common, reqMsg *share_message.OrderDetailInfoPara) easygo.IMessage {
	logs.Info("===api RpcOrderInfo===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()
	order := share_message.TableShopOrder{}
	var e error

	if *reqMsg.Type == share_message.BuySell_Type_Buyer {
		e = col.Find(bson.M{"_id": reqMsg.GetOrderId(), "delete_buy": 0}).One(&order)
	} else {
		e = col.Find(bson.M{"_id": reqMsg.GetOrderId(), "delete_sell": 0}).One(&order)
	}

	if e == mgo.ErrNotFound {

		SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_DELETE_ORDER_NOT_FOUND)

		return easygo.NewFailMsg(ORDER_DELETE_ORDER_NOT_FOUND)
	}

	if e != nil {

		logs.Error(e)

		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	var blackList []PLAYER_ID = ShopInstance.GetBlackLists(common.GetUserId())

	var addFlag int32 = 0
	for _, black := range blackList {
		if order.GetSponsorId() == black || order.GetReceiverId() == black {
			addFlag = 1
			break
		}
	}

	if addFlag == 1 {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DETAIL_ORDER_BLACK_VAR)

		return easygo.NewFailMsg(DETAIL_ORDER_BLACK_VAR)
	}

	var avatar string
	var nickName string
	var sex int32

	if reqMsg.GetType() == share_message.BuySell_Type_Seller {
		avatar = order.GetReceiverAvatar()
		nickName = order.GetReceiverNickname()
		sex = order.GetReceiverSex()
	} else {
		avatar = order.GetSponsorAvatar()
		nickName = order.GetSponsorNickname()
		sex = order.GetSponsorSex()
	}

	item := share_message.BuyItem{
		ItemId:    easygo.NewInt64(order.Items.GetItemId()),
		ItemFile:  order.Items.ItemFile,
		Name:      easygo.NewString(order.Items.GetName()),
		Price:     easygo.NewInt32(order.Items.GetPrice()),
		SponsorId: easygo.NewInt64(order.GetSponsorId()),
		Count:     easygo.NewInt32(order.Items.GetCount()),
		Avatar:    easygo.NewString(avatar),
		Nickname:  easygo.NewString(nickName),
		Sex:       easygo.NewInt32(sex),
		//OriginPrice: easygo.NewInt32(order.Items.GetOriginPrice()),
		Title: easygo.NewString(order.Items.GetTitle()),
		DeliverAddress: &share_message.DeliverAddress{
			Name:          easygo.NewString(order.DeliverAddress.GetName()),
			Phone:         easygo.NewString(order.DeliverAddress.GetPhone()),
			Region:        easygo.NewString(order.DeliverAddress.GetRegion()),
			DetailAddress: easygo.NewString(order.DeliverAddress.GetDetailAddress()),
		},
	}

	var expressBodyList []*share_message.QueryExpressBody = []*share_message.QueryExpressBody{}
	var errCode string
	var expressName string
	var expressPhone string

	//取最新物流信息
	if order.GetExpressCode() != "" && order.GetExpressCom() != "" {

		expressName, expressPhone = GetExpressNamePhone(order.GetExpressCom())

		//不是顺丰的时候
		if order.GetExpressCom() != "sf" {
			expressBodyList, errCode, _ = GetExpressInfos(
				order.GetOrderId(),
				order.GetExpressCom(),
				order.GetExpressCode(),
				"",
				"",
				0,
			)
		} else {
			expressBodyList, errCode, _ = GetExpressInfos(
				order.GetOrderId(),
				order.GetExpressCom(),
				order.GetExpressCode(),
				order.DeliverAddress.GetPhone(),
				order.ReceiveAddress.GetPhone(),
				0,
			)
		}

		if errCode != "" {
			if errCode == EXPRESS_QUERY_ERROR_CODE_998 {

				s := fmt.Sprintf("查询缓存数据库出错：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.GetDeliverAddress().GetPhone(),
					order.GetReceiveAddress().GetPhone())

				logs.Error(s)
				SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

				return easygo.NewFailMsg(DATABASE_ERROR)

				//接口那里出错不直接通知客户端返回一个空物流打印下err
			} else if errCode == EXPRESS_QUERY_ERROR_CODE_999 {

				s := fmt.Sprintf("快递查询请求快递接口出错：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.DeliverAddress.GetPhone(),
					order.ReceiveAddress.GetPhone())

				logs.Error(s)

			} else if errCode == EXPRESS_QUERY_ERROR_CODE_1 {

				s := fmt.Sprintf("快递查询快递公司错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.DeliverAddress.GetPhone(),
					order.ReceiveAddress.GetPhone())

				logs.Error(s)

			} else if errCode == EXPRESS_QUERY_ERROR_CODE_2 {

				s := fmt.Sprintf("快递查询运单号错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.DeliverAddress.GetPhone(),
					order.ReceiveAddress.GetPhone())

				logs.Error(s)

			} else if errCode == EXPRESS_QUERY_ERROR_CODE_3 {

				s := fmt.Sprintf("快递查询失败：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.DeliverAddress.GetPhone(),
					order.ReceiveAddress.GetPhone())

				logs.Error(s)

			} else if errCode == EXPRESS_QUERY_ERROR_CODE_4 {

				s := fmt.Sprintf("快递查询查不到物流信息：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.DeliverAddress.GetPhone(),
					order.ReceiveAddress.GetPhone())

				logs.Error(s)

			} else if errCode == EXPRESS_QUERY_ERROR_CODE_5 {

				s := fmt.Sprintf("快递查询寄件人或收件人手机尾号错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.DeliverAddress.GetPhone(),
					order.ReceiveAddress.GetPhone())

				logs.Error(s)

			} else {
				s := fmt.Sprintf("快递查询其他错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
					errCode,
					order.GetExpressCode(),
					order.GetExpressCom(),
					order.DeliverAddress.GetPhone(),
					order.ReceiveAddress.GetPhone())

				logs.Error(s)
			}
		}
	}

	orderDetailInfo := &share_message.OrderDetailInfo{
		OrderId: easygo.NewInt64(order.GetOrderId()),
		Item:    &item,
		State:   easygo.NewInt32(order.GetState()),
		Address: &share_message.ReceiveAddress{
			DetailAddress: easygo.NewString(order.ReceiveAddress.GetDetailAddress()),
			Name:          easygo.NewString(order.ReceiveAddress.GetName()),
			Region:        easygo.NewString(order.ReceiveAddress.GetRegion()),
			Phone:         easygo.NewString(order.ReceiveAddress.GetPhone())},
		ExpressCode:   easygo.NewString(order.GetExpressCode()),
		ExpressCom:    easygo.NewString(order.GetExpressCom()),
		ExpressName:   easygo.NewString(expressName),
		CreateTime:    easygo.NewInt64(order.GetCreateTime()),
		ServerNowTime: easygo.NewInt64(time.Now().Unix()),
		PayTime:       easygo.NewInt64(order.GetPayTime()),
		SendTime:      easygo.NewInt64(order.GetSendTime()),
		FinishTime:    easygo.NewInt64(order.GetFinishTime()),
		ExpressInfos:  expressBodyList,
		ExpressPhone:  easygo.NewString(expressPhone)}

	return &share_message.OrderDetailInfoShow{
		Result:          easygo.NewInt32(0),
		Msg:             easygo.NewString(""),
		OrderDetailInfo: orderDetailInfo}
}

func (self *ServiceForHall) RpcCancelOrder(common *base.Common, reqMsg *share_message.OrderID) easygo.IMessage {
	logs.Info("===api RpcCancelOrder===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	var order *share_message.TableShopOrder = &share_message.TableShopOrder{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()

	e := col.Find(bson.M{"_id": reqMsg.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_PAY}).One(order)

	if e == mgo.ErrNotFound {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", CANCEL_ORDER_STATE_CHANGE)

		return easygo.NewFailMsg(CANCEL_ORDER_STATE_CHANGE)
	}

	if e != nil {

		logs.Error(e)

		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	if order.GetState() == for_game.SHOP_ORDER_CANCEL {

		SendToHallServerByApi(common.GetUserId(), "RpcToast", "重复操作！")

		return easygo.NewFailMsg("重复操作！")
	}

	if order.GetState() != for_game.SHOP_ORDER_WAIT_PAY {

		SendToHallServerByApi(common.GetUserId(), "RpcToast", CANCEL_ORDER_STATE_CHANGE)

		return easygo.NewFailMsg(CANCEL_ORDER_STATE_CHANGE)

	} else {

		//以每个订单为单位取得锁,取不到说明订单状态在改变中
		//1、取得订单的分布式锁(不需要重试）
		lockKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_WAIT_PAY_MUTEX, reqMsg.GetOrderId())
		//取得分布式锁（失效时间设置10秒）
		errLock := easygo.RedisMgr.GetC().DoRedisLockNoRetry(lockKey, 10)
		defer easygo.RedisMgr.GetC().DoRedisUnlock(lockKey)

		//如果未取得锁
		if errLock != nil {
			s := fmt.Sprintf("RpcCancelOrder 单key取得订单redis分布式无重试锁失败redis key is %v", lockKey)
			logs.Error(s)
			logs.Error(errLock)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", "取消失败,刷新重试！")

			return easygo.NewFailMsg("取消失败,刷新重试！")
		}

		//取得订单对应的商品的分布式锁，此锁需要重试(恢复库存的竞争)
		tempItemLockKey := for_game.MakeRedisKey(for_game.SHOP_ITEM_PAY_MUTEX, order.GetItems().GetItemId())
		//2、取得订单对应的商品的分布式锁，此锁需要重试，直到重试次数结束提示退出
		errLock2 := easygo.RedisMgr.GetC().DoRedisLockWithRetry(tempItemLockKey, 10)
		defer easygo.RedisMgr.GetC().DoRedisUnlock(tempItemLockKey)

		//如果重试后还未取得锁就直接不做了
		if errLock2 != nil {
			s := fmt.Sprintf("RpcCancelOrder 单key取得商品redis分布式重试锁失败redis key is %v", tempItemLockKey)
			logs.Error(s)
			logs.Error(errLock2)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", "取消失败,刷新重试！")

			return easygo.NewFailMsg("取消失败,刷新重试！")
		}

		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
		defer closeFun()

		e := col.Update(
			bson.M{"_id": reqMsg.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_PAY},
			bson.M{"$set": bson.M{"state": for_game.SHOP_ORDER_CANCEL, "update_time": time.Now().Unix()}})

		if e == mgo.ErrNotFound {
			logs.Error(e)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", CANCEL_ORDER_STATE_CHANGE)

			return easygo.NewFailMsg(CANCEL_ORDER_STATE_CHANGE)
		}

		if e != nil {
			logs.Error(e)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)
		}
	}

	return &share_message.CancelOrderResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(CANCEL_ORDER_SUCCESS)}
}

func (self *ServiceForHall) RpcDeleteOrder(common *base.Common, reqMsg *share_message.OrderID) easygo.IMessage {

	logs.Info("===api RpcDeleteOrder===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	var order *share_message.TableShopOrder = &share_message.TableShopOrder{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()

	e := col.Find(bson.M{"_id": reqMsg.GetOrderId()}).One(order)

	if e == mgo.ErrNotFound {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_DELETE_ORDER_NOT_FOUND)

		return easygo.NewFailMsg(ORDER_DELETE_ORDER_NOT_FOUND)
	}

	if e != nil {

		logs.Error(e)

		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	if order.GetState() == for_game.SHOP_ORDER_FINISH ||
		order.GetState() == for_game.SHOP_ORDER_CANCEL ||
		order.GetState() == for_game.SHOP_ORDER_EXPIRE ||
		order.GetState() == for_game.SHOP_ORDER_BACKSTAGE_CANCLE ||
		order.GetState() == for_game.SHOP_ORDER_EVALUTE {

		if common.GetUserId() == order.GetReceiverId() {

			e = col.Update(
				bson.M{"_id": reqMsg.GetOrderId()},
				bson.M{"$set": bson.M{"delete_buy": 1, "update_time": time.Now().Unix()}})

		} else if common.GetUserId() == order.GetSponsorId() {

			e = col.Update(
				bson.M{"_id": reqMsg.GetOrderId()},
				bson.M{"$set": bson.M{"delete_sell": 1, "update_time": time.Now().Unix()}})

		} else {

			SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_DELETE_ORDER_NOT_OWNER)

			return easygo.NewFailMsg(ORDER_DELETE_ORDER_NOT_OWNER)

		}

		if e == mgo.ErrNotFound {
			SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_DELETE_ORDER_NOT_OWNER)

			return easygo.NewFailMsg(ORDER_DELETE_ORDER_NOT_OWNER)
		}

		if e != nil && e != mgo.ErrNotFound {
			logs.Error(e)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)
		}

		return &share_message.DeleteOrderResult{
			Msg:    easygo.NewString(ORDER_DELETE_SUCCESS),
			Result: easygo.NewInt32(0)}

	}

	SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_DELETE_ORDER_STATE_WRONG)

	return easygo.NewFailMsg(ORDER_DELETE_ORDER_STATE_WRONG)
}

func (self *ServiceForHall) RpcItemListForMyReleaseOnline(common *base.Common, reqMsg *share_message.MyReleaseInfo) easygo.IMessage {
	logs.Info("===api RpcItemListForMyReleaseOnline===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	var count int32 = 0
	page := reqMsg.GetPage()
	pageSize := reqMsg.GetPageSize()

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()

	var list []*share_message.TableShopItem

	query := col.Find(
		bson.M{"player_id": common.GetUserId(),
			"$or": []bson.M{
				bson.M{"state": for_game.SHOP_ITEM_SALE},
				bson.M{"state": for_game.SHOP_ITEM_IN_AUDIT}}})

	cnt, errVa := query.Count()
	count = int32(cnt)

	if errVa != nil {

		logs.Error(errVa)

		return &share_message.ItemListForMyRelease{
			Page:     &page,
			PageSize: &pageSize,
			Count:    &count}
	}

	e := query.Limit(int(pageSize)).Skip(int(page * pageSize)).Sort("-create_time").All(&list)

	if e != mgo.ErrNotFound && e != nil {

		logs.Error(e)

		return &share_message.ItemListForMyRelease{Page: &page,
			PageSize: &pageSize,
			Count:    &count}
	}
	newList := []*share_message.ShopItem{}

	if e == nil && list != nil {
		for _, value := range list {
			var item_file share_message.ItemFile = share_message.ItemFile{}
			if nil != value.ItemFiles && len(value.ItemFiles) != 0 {
				item_file = *value.ItemFiles[0]
			}
			p := for_game.GetRedisPlayerBase(value.GetPlayerId())
			newItem := share_message.ShopItem{
				ItemId:     value.ItemId,
				Price:      value.Price,
				ItemFile:   &item_file,
				Title:      value.Title,
				StoreCount: value.RealStoreCnt,
				PlayerId:   value.PlayerId,
				Nickname:   value.Nickname,
				Avatar:     value.Avatar,
				Account:    value.PlayerAccount,
				Sex:        value.Sex,
				Name:       value.Name,
				State:      value.State,
				Types:      easygo.NewInt32(p.GetTypes()),
			}
			newList = append(newList, &newItem)
		}
		return &share_message.ItemListForMyRelease{
			Items:    newList,
			Page:     &page,
			PageSize: &pageSize,
			Count:    &count}

	} else {

		return &share_message.ItemListForMyRelease{
			Page:     &page,
			PageSize: &pageSize,
			Count:    &count}
	}
}

func (self *ServiceForHall) RpcItemListForMyReleaseOffline(common *base.Common, reqMsg *share_message.MyReleaseInfo) easygo.IMessage {
	logs.Info("===api RpcItemListForMyReleaseOffline===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	var count int32 = 0
	page := reqMsg.GetPage()
	pageSize := reqMsg.GetPageSize()

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()

	var list []*share_message.TableShopItem

	query := col.Find(
		bson.M{"player_id": common.GetUserId(),
			"$or": []bson.M{
				bson.M{"state": for_game.SHOP_ITEM_SOLD_OUT},
				bson.M{"state": for_game.SHOP_ITEM_FAIL_AUDIT}}})

	cnt, errVa := query.Count()
	count = int32(cnt)
	if errVa != nil {
		logs.Error(errVa)
		return &share_message.ItemListForMyRelease{
			Page:     &page,
			PageSize: &pageSize,
			Count:    &count}
	}

	e := query.Limit(int(pageSize)).Skip(int(page * pageSize)).Sort("-sold_out_time").All(&list)

	if e != mgo.ErrNotFound && e != nil {
		logs.Error(e)
		return &share_message.ItemListForMyRelease{
			Page:     &page,
			PageSize: &pageSize,
			Count:    &count}
	}

	newList := []*share_message.ShopItem{}

	if e == nil && list != nil {

		for _, value := range list {

			var itemFile share_message.ItemFile = share_message.ItemFile{}

			if nil != value.ItemFiles && len(value.ItemFiles) != 0 {
				itemFile = *value.ItemFiles[0]
			}

			shopItem := ShopInstance.GetItemFromCache(value.GetItemId())
			//判断是否下架
			var itemFlag bool
			if nil != shopItem {
				itemFlag = false
			} else {
				itemFlag = true
			}
			var types int32
			if p := for_game.GetRedisPlayerBase(value.GetPlayerId()); p != nil {
				types = p.GetTypes()
			}
			newItem := share_message.ShopItem{
				ItemId:     value.ItemId,
				Price:      value.Price,
				ItemFile:   &itemFile,
				Title:      value.Title,
				StoreCount: value.RealStoreCnt,
				PlayerId:   value.PlayerId,
				Nickname:   value.Nickname,
				Avatar:     value.Avatar,
				Account:    value.PlayerAccount,
				Sex:        value.Sex,
				Name:       value.Name,
				State:      value.State,
				ItemFlag:   easygo.NewBool(itemFlag),
				Types:      easygo.NewInt32(types),
			}

			newList = append(newList, &newItem)
		}

		return &share_message.ItemListForMyRelease{
			Items:    newList,
			Page:     &page,
			PageSize: &pageSize,
			Count:    &count}

	} else {

		return &share_message.ItemListForMyRelease{
			Page:     &page,
			PageSize: &pageSize,
			Count:    &count}
	}
}

func (self *ServiceForHall) RpcItemSearch(common *base.Common, reqMsg *share_message.SearchInfo) easygo.IMessage {
	logs.Info("===api RpcItemSearch===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	shopItems := ShopInstance.shop_items
	var cnt int32 = 0

	if reqMsg.GetContent() == "" {

		return &share_message.SearchResult{
			PageSize: reqMsg.PageSize,
			Page:     reqMsg.PageSize,
			Count:    &cnt}
	}

	re, err := regexp.Compile(reqMsg.GetContent())

	if err != nil {

		logs.Error("输入了不法字符%s", reqMsg.GetContent())
		return &share_message.SearchResult{PageSize: reqMsg.PageSize, Page: reqMsg.PageSize, Count: &cnt}
	}

	items := []*ShopItem{}

	for _, item := range *shopItems {

		if item.stock_count > 0 {
			if re.MatchString(item.title) || re.MatchString(item.name) {
				items = append(items, item)
				continue
			}

			//因为结构体中定义的是数组，不是数组指针，上面的循环会出现同一个值
			if nil != item.other_type && len(item.other_type) > 0 {

				for i := 0; i < len(item.other_type); i++ {
					if re.MatchString(item.other_type[i]) {
						items = append(items, item)
						break
					}
				}
			}
		}
	}

	var addFlag int32 = 0

	//价格升序
	if reqMsg.GetSearchFlag() == share_message.Search_Type_PriceAsc {

		lastShopItems := LastShopItemsPriceAsc{}

		var blackList []PLAYER_ID = ShopInstance.GetBlackLists(common.GetUserId())

		addFlag = 0
		for _, value := range items {
			for _, black := range blackList {
				if value.player_id == black {
					addFlag = 1
					break
				}
			}
			if addFlag == 0 {
				lastShopItems = append(lastShopItems, value)
			}

			addFlag = 0
		}

		sort.Sort(lastShopItems)

		var list []*share_message.ShopItem = []*share_message.ShopItem{}
		var count = int32(len(lastShopItems))

		page := reqMsg.GetPage()
		pageSize := reqMsg.GetPageSize()

		if page*pageSize <= count {

			var list []*share_message.ShopItem
			for i := page * pageSize; i < pageSize*(page+1) && i < count; i++ {
				list = append(list, ShopInstance.BriefItem(lastShopItems[i]))
			}

			return &share_message.SearchResult{
				Items:    list,
				PageSize: &pageSize,
				Page:     &page,
				Count:    &count}
		}
		return &share_message.SearchResult{
			Items:    list,
			PageSize: &pageSize,
			Page:     &page,
			Count:    &count}

		//价格降序
	} else if reqMsg.GetSearchFlag() == share_message.Search_Type_PriceDesc {

		lastShopItems := LastShopItemsPriceDesc{}

		var blackList []PLAYER_ID = ShopInstance.GetBlackLists(common.GetUserId())

		addFlag = 0
		for _, value := range items {
			for _, black := range blackList {
				if value.player_id == black {
					addFlag = 1
					break
				}
			}
			if addFlag == 0 {
				lastShopItems = append(lastShopItems, value)
			}

			addFlag = 0
		}

		sort.Sort(lastShopItems)

		var list []*share_message.ShopItem = []*share_message.ShopItem{}
		var count = int32(len(lastShopItems))

		page := reqMsg.GetPage()
		pageSize := reqMsg.GetPageSize()

		if page*pageSize <= count {

			var list []*share_message.ShopItem

			for i := page * pageSize; i < pageSize*(page+1) && i < count; i++ {

				list = append(list, ShopInstance.BriefItem(lastShopItems[i]))

			}

			return &share_message.SearchResult{
				Items:    list,
				PageSize: &pageSize,
				Page:     &page,
				Count:    &count}
		}

		return &share_message.SearchResult{
			Items:    list,
			PageSize: &pageSize,
			Page:     &page,
			Count:    &count}

	} else if reqMsg.GetSearchFlag() == share_message.Search_Type_NewSort {

		lastShopItems := LastShopItemsNew{}

		var blackList []PLAYER_ID = ShopInstance.GetBlackLists(common.GetUserId())

		addFlag = 0
		for _, value := range items {
			for _, black := range blackList {
				if value.player_id == black {
					addFlag = 1
					break
				}
			}
			if addFlag == 0 {
				lastShopItems = append(lastShopItems, value)
			}

			addFlag = 0
		}

		sort.Sort(lastShopItems)

		var list []*share_message.ShopItem = []*share_message.ShopItem{}
		var count = int32(len(lastShopItems))

		page := reqMsg.GetPage()
		pageSize := reqMsg.GetPageSize()

		if page*pageSize <= count {

			var list []*share_message.ShopItem

			for i := page * pageSize; i < pageSize*(page+1) && i < count; i++ {
				list = append(list, ShopInstance.BriefItem(lastShopItems[i]))
			}

			return &share_message.SearchResult{
				Items:    list,
				PageSize: &pageSize,
				Page:     &page,
				Count:    &count}
		}

		return &share_message.SearchResult{
			Items:    list,
			PageSize: &pageSize,
			Page:     &page,
			Count:    &count}

	} else if reqMsg.GetSearchFlag() == share_message.Search_Type_SalesSort {

		lastShopItems := LastShopItemsSalesDesc{}

		var blackList []PLAYER_ID = ShopInstance.GetBlackLists(common.GetUserId())

		addFlag = 0
		for _, value := range items {
			for _, black := range blackList {
				if value.player_id == black {
					addFlag = 1
					break
				}
			}
			if addFlag == 0 {
				lastShopItems = append(lastShopItems, value)
			}

			addFlag = 0
		}

		sort.Sort(lastShopItems)

		var list []*share_message.ShopItem = []*share_message.ShopItem{}
		var count = int32(len(lastShopItems))

		page := reqMsg.GetPage()
		pageSize := reqMsg.GetPageSize()

		if page*pageSize <= count {

			var list []*share_message.ShopItem

			for i := page * pageSize; i < pageSize*(page+1) && i < count; i++ {
				list = append(list, ShopInstance.BriefItem(lastShopItems[i]))
			}

			return &share_message.SearchResult{
				Items:    list,
				PageSize: &pageSize,
				Page:     &page,
				Count:    &count}
		}
		return &share_message.SearchResult{
			Items:    list,
			PageSize: &pageSize,
			Page:     &page,
			Count:    &count}

	} else if reqMsg.GetSearchFlag() == share_message.Search_Type_StoreSort {

		lastShopItems := LastShopItemsStoreDesc{}

		var blackList []PLAYER_ID = ShopInstance.GetBlackLists(common.GetUserId())

		addFlag = 0
		for _, value := range items {
			for _, black := range blackList {
				if value.player_id == black {
					addFlag = 1
					break
				}
			}
			if addFlag == 0 {
				lastShopItems = append(lastShopItems, value)
			}

			addFlag = 0
		}

		sort.Sort(lastShopItems)

		var list []*share_message.ShopItem = []*share_message.ShopItem{}
		var count = int32(len(lastShopItems))

		page := reqMsg.GetPage()
		pageSize := reqMsg.GetPageSize()

		if page*pageSize <= count {

			var list []*share_message.ShopItem

			for i := page * pageSize; i < pageSize*(page+1) && i < count; i++ {
				list = append(list, ShopInstance.BriefItem(lastShopItems[i]))
			}

			return &share_message.SearchResult{
				Items:    list,
				PageSize: &pageSize,
				Page:     &page,
				Count:    &count}
		}
		return &share_message.SearchResult{
			Items:    list,
			PageSize: &pageSize,
			Page:     &page,
			Count:    &count}

	} else {

		lastShopItems := LastShopItemsComDesc{}
		var blackList []PLAYER_ID = ShopInstance.GetBlackLists(common.GetUserId())

		addFlag = 0
		for _, value := range items {
			for _, black := range blackList {
				if value.player_id == black {
					addFlag = 1
					break
				}
			}
			if addFlag == 0 {
				lastShopItems = append(lastShopItems, value)
			}

			addFlag = 0
		}

		sort.Sort(lastShopItems)

		var list []*share_message.ShopItem = []*share_message.ShopItem{}
		var count = int32(len(lastShopItems))

		page := reqMsg.GetPage()
		pageSize := reqMsg.GetPageSize()

		if page*pageSize <= count {

			var list []*share_message.ShopItem

			for i := page * pageSize; i < pageSize*(page+1) && i < count; i++ {
				list = append(list, ShopInstance.BriefItem(lastShopItems[i]))
			}

			return &share_message.SearchResult{
				Items:    list,
				PageSize: &pageSize,
				Page:     &page,
				Count:    &count}
		}
		return &share_message.SearchResult{
			Items:    list,
			PageSize: &pageSize,
			Page:     &page,
			Count:    &count}
	}
}

func (self *ServiceForHall) RpcExpressCodeUpload(common *base.Common, reqMsg *share_message.ExpressInfo) easygo.IMessage {
	logs.Info("===api RpcExpressCodeUpload===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	//以每个订单为单位取得锁,取不到说明订单状态在改变中
	lockKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_SEND_MUTEX, reqMsg.GetOrderId())
	//取得分布式锁（失效时间设置10秒）
	errLock := easygo.RedisMgr.GetC().DoRedisLockNoRetry(lockKey, 10)
	defer easygo.RedisMgr.GetC().DoRedisUnlock(lockKey)

	//如果未取得锁
	if errLock != nil {
		s := fmt.Sprintf("RpcExpressCodeUpload 单key取得redis分布式无重试锁失败,redis key is %v", lockKey)
		logs.Error(s)
		logs.Error(errLock)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_EXPRESS_UPLOAD_REPEATED)

		return easygo.NewFailMsg(ORDER_EXPRESS_UPLOAD_REPEATED)
	}

	var order *share_message.TableShopOrder = &share_message.TableShopOrder{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()

	e := col.Find(bson.M{"_id": reqMsg.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_SEND}).One(order)

	if e == mgo.ErrNotFound {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_EXPRESS_UPLOAD_REPEATED)

		return easygo.NewFailMsg(ORDER_EXPRESS_UPLOAD_REPEATED)
	}

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	if order.GetState() != for_game.SHOP_ORDER_WAIT_SEND {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_EXPRESS_UPLOAD_REPEATED)

		return easygo.NewFailMsg(ORDER_EXPRESS_UPLOAD_REPEATED)
	}

	// TODO
	//确认快递单号和公司是否正确
	//由于有快递单号但是未发货这种情况存在，接口返回的是204303(这个错误包含很多情况)
	//204303主要是单号不存在，所以无法判断该返回值是有快递未发货，还是快递单号和公司组合填错了
	_, errCode, reason := GetExpressInfos(
		order.GetOrderId(),
		reqMsg.GetCom(),
		reqMsg.GetCode(),
		reqMsg.GetSendPhone(),
		order.GetReceiveAddress().GetPhone(),
		1,
	)

	if errCode != "" {
		if errCode == EXPRESS_QUERY_ERROR_CODE_999 {

			s := fmt.Sprintf("快递查询请求快递接口出错：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				reqMsg.GetCom(),
				reqMsg.GetCode(),
				reqMsg.GetSendPhone(),
				order.GetReceiveAddress().GetPhone(),
			)

			logs.Error(s)

			SendToHallServerByApi(common.GetUserId(), "RpcToast", EXPRESS_QUERY_ERROR_MSG_999)

			return easygo.NewFailMsg(EXPRESS_QUERY_ERROR_MSG_999)

		} else if errCode == EXPRESS_QUERY_ERROR_CODE_1 {

			s := fmt.Sprintf("快递查询快递公司错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				reqMsg.GetCom(),
				reqMsg.GetCode(),
				reqMsg.GetSendPhone(),
				order.GetReceiveAddress().GetPhone(),
			)

			logs.Info(s)

			return &share_message.ExpressCodeResult{
				Result: easygo.NewInt32(1),
				Msg:    easygo.NewString(EXPRESS_QUERY_ERROR_MSG)}

		} else if errCode == EXPRESS_QUERY_ERROR_CODE_2 {

			s := fmt.Sprintf("快递查询运单号错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				reqMsg.GetCom(),
				reqMsg.GetCode(),
				reqMsg.GetSendPhone(),
				order.GetReceiveAddress().GetPhone(),
			)

			logs.Info(s)

			return &share_message.ExpressCodeResult{
				Result: easygo.NewInt32(1),
				Msg:    easygo.NewString(EXPRESS_QUERY_ERROR_MSG)}

		} else if errCode == EXPRESS_QUERY_ERROR_CODE_3 {

			s := fmt.Sprintf("快递查询失败：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				reqMsg.GetCom(),
				reqMsg.GetCode(),
				reqMsg.GetSendPhone(),
				order.GetReceiveAddress().GetPhone(),
			)
			logs.Info(s)

			//这个code只能人为部分解析
			if strings.Contains(reason, "参数错误") {
				return &share_message.ExpressCodeResult{
					Result: easygo.NewInt32(1),
					Msg:    easygo.NewString(EXPRESS_QUERY_ERROR_MSG)}
			}

		} else if errCode == EXPRESS_QUERY_ERROR_CODE_4 {

			s := fmt.Sprintf("快递查询查不到物流信息：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				reqMsg.GetCom(),
				reqMsg.GetCode(),
				reqMsg.GetSendPhone(),
				order.GetReceiveAddress().GetPhone(),
			)

			logs.Info(s)

		} else if errCode == EXPRESS_QUERY_ERROR_CODE_5 {

			s := fmt.Sprintf("快递查询寄件人或收件人手机尾号错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				reqMsg.GetCom(),
				reqMsg.GetCode(),
				reqMsg.GetSendPhone(),
				order.GetReceiveAddress().GetPhone(),
			)

			logs.Info(s)

		} else {
			s := fmt.Sprintf("快递查询其他错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				reqMsg.GetCom(),
				reqMsg.GetCode(),
				reqMsg.GetSendPhone(),
				order.GetReceiveAddress().GetPhone(),
			)

			logs.Info(s)

		}
	}

	expressName, _ := GetExpressNamePhone(reqMsg.GetCom())
	e = col.Update(
		bson.M{"_id": reqMsg.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_SEND},
		bson.M{"$set": bson.M{
			"express_code":         reqMsg.GetCode(),
			"express_com":          reqMsg.GetCom(),
			"state":                for_game.SHOP_ORDER_WAIT_RECEIVE,
			"receive_time":         time.Now().Unix() + 7*24*3600,
			"send_time":            time.Now().Unix(),
			"express_name":         expressName,
			"receiver_notify_flag": true,
			"sponsor_notify_flag":  true,
			"update_time":          time.Now().Unix(),
		}})

	if e == mgo.ErrNotFound {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_EXPRESS_UPLOAD_REPEATED)

		return easygo.NewFailMsg(ORDER_EXPRESS_UPLOAD_REPEATED)
	}

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	easygo.Spawn(func(orderPara *share_message.TableShopOrder) {

		if nil != orderPara {

			var content string = MESSAGE_TO_BUYER_SEND
			typeValue := share_message.BuySell_Type_Buyer

			ShopInstance.InsMessageNotify(
				easygo.NewString(content),
				&typeValue,
				orderPara)

			var jgContent string = MESSAGE_TO_BUYER_SEND_PUSH
			ShopInstance.JGMessageNotify(jgContent, orderPara.GetReceiverId(), orderPara.GetOrderId(), typeValue)

			//提交物流 商城订单红点推送

			SendMsgToHallClientNew([]int64{orderPara.GetReceiverId()}, "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
				OrderId: easygo.NewInt64(orderPara.GetOrderId())})

			/*		SendToPlayer(orderPara.GetReceiverId(), "RpcShopOrderNotify",
					&share_message.ShopOrderNotifyInfoWithWho{
						PlayerId: easygo.NewInt64(orderPara.GetReceiverId()),
						OrderId:  easygo.NewInt64(orderPara.GetOrderId()),
					})
			*/
			//提交物流 商城订单红点推送

			SendMsgToHallClientNew([]int64{orderPara.GetSponsorId()}, "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
				OrderId: easygo.NewInt64(orderPara.GetOrderId())})

			/*			SendToPlayer(orderPara.GetSponsorId(), "RpcShopOrderNotify",
						&share_message.ShopOrderNotifyInfoWithWho{
							PlayerId: easygo.NewInt64(orderPara.GetSponsorId()),
							OrderId:  easygo.NewInt64(orderPara.GetOrderId()),
						})
			*/
		} else {
			logs.Debug("提交物流后发通知,缺少订单")
		}
	}, order)

	return &share_message.ExpressCodeResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(ORDER_EXPRESS_UPLOAD_SUCCESS)}

}

func (self *ServiceForHall) RpcExpressComInfos(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===api RpcExpressComInfos===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	return GetExpressComInfos()
}

func (self *ServiceForHall) RpcExpressInfos(common *base.Common, reqMsg *share_message.QueryExpressInfo) easygo.IMessage {
	logs.Info("===api RpcExpressInfos===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()
	order := share_message.TableShopOrder{}
	e := col.Find(bson.M{"_id": reqMsg.GetOrderId()}).One(&order)

	if e == mgo.ErrNotFound {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_EXPRESS_QUERY_NOT_FOUND)

		return easygo.NewFailMsg(ORDER_EXPRESS_QUERY_NOT_FOUND)
	}

	if e != nil {

		logs.Error(e)

		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	var errCode string
	var expressName string
	var expressPhone string
	var expressBodyList []*share_message.QueryExpressBody = []*share_message.QueryExpressBody{}

	expressName, expressPhone = GetExpressNamePhone(order.GetExpressCom())

	expressBodyList, errCode, _ = GetExpressInfos(
		order.GetOrderId(),
		order.GetExpressCom(),
		order.GetExpressCode(),
		order.GetDeliverAddress().GetPhone(),
		order.GetReceiveAddress().GetPhone(),
		0,
	)

	if errCode != "" {
		if errCode == EXPRESS_QUERY_ERROR_CODE_998 {

			s := fmt.Sprintf("查询缓存数据库出错：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				order.GetExpressCode(),
				order.GetExpressCom(),
				order.GetDeliverAddress().GetPhone(),
				order.GetReceiveAddress().GetPhone())

			logs.Error(s)
			SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

			return easygo.NewFailMsg(DATABASE_ERROR)

			//接口那里出错不直接通知客户端返回一个空物流打印下err
		} else if errCode == EXPRESS_QUERY_ERROR_CODE_999 {

			s := fmt.Sprintf("快递查询请求快递接口出错：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				order.GetExpressCode(),
				order.GetExpressCom(),
				order.GetDeliverAddress().GetPhone(),
				order.GetReceiveAddress().GetPhone())

			logs.Error(s)

		} else if errCode == EXPRESS_QUERY_ERROR_CODE_1 {

			s := fmt.Sprintf("快递查询快递公司错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				order.GetExpressCode(),
				order.GetExpressCom(),
				order.GetDeliverAddress().GetPhone(),
				order.GetReceiveAddress().GetPhone())

			logs.Error(s)

		} else if errCode == EXPRESS_QUERY_ERROR_CODE_2 {

			s := fmt.Sprintf("快递查询运单号错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				order.GetExpressCom(),
				order.GetExpressCode(),
				order.GetDeliverAddress().GetPhone(),
				order.GetReceiveAddress().GetPhone())

			logs.Error(s)

		} else if errCode == EXPRESS_QUERY_ERROR_CODE_3 {

			s := fmt.Sprintf("快递查询失败：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				order.GetExpressCode(),
				order.GetExpressCom(),
				order.GetDeliverAddress().GetPhone(),
				order.GetReceiveAddress().GetPhone())

			logs.Error(s)

		} else if errCode == EXPRESS_QUERY_ERROR_CODE_4 {

			s := fmt.Sprintf("快递查询查不到物流信息：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				order.GetExpressCode(),
				order.GetExpressCom(),
				order.GetDeliverAddress().GetPhone(),
				order.GetReceiveAddress().GetPhone())

			logs.Error(s)

		} else if errCode == EXPRESS_QUERY_ERROR_CODE_5 {

			s := fmt.Sprintf("快递查询寄件人或收件人手机尾号错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				order.GetExpressCode(),
				order.GetExpressCom(),
				order.GetDeliverAddress().GetPhone(),
				order.GetReceiveAddress().GetPhone())

			logs.Error(s)

		} else {

			s := fmt.Sprintf("快递查询其他错误：错误码(%s)，快递单号(%s)，快递公司(%s)，发件人手机号(%s)，收件人手机号(%s)",
				errCode,
				order.GetExpressCode(),
				order.GetExpressCom(),
				order.GetDeliverAddress().GetPhone(),
				order.GetReceiveAddress().GetPhone())

			logs.Error(s)

		}
	}

	return &share_message.QueryExpressInfosResult{
		Result:       easygo.NewInt32(0),
		Msg:          easygo.NewString(""),
		ExpressInfos: expressBodyList,
		ExpressPhone: easygo.NewString(expressPhone),
		ExpressName:  easygo.NewString(expressName),
	}
}

func (self *ServiceForHall) RpcEditOrderAddress(common *base.Common, reqMsg *share_message.EditOrderAddress) easygo.IMessage {
	logs.Info("===api RpcEditOrderAddress===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	var order *share_message.TableShopOrder = &share_message.TableShopOrder{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()

	e := col.Find(bson.M{"_id": reqMsg.GetOrderId()}).One(order)

	if e == mgo.ErrNotFound {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_EDIT_ADDRESSS_ORDER_NOT_FOUND)

		return easygo.NewFailMsg(ORDER_EDIT_ADDRESSS_ORDER_NOT_FOUND)
	}

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	if order.GetState() != for_game.SHOP_ORDER_WAIT_SEND && order.GetState() != for_game.SHOP_ORDER_WAIT_PAY {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_EDIT_ADDRESSS_STATE_ERROR)

		return easygo.NewFailMsg(ORDER_EDIT_ADDRESSS_STATE_ERROR)
	}
	if order.GetReceiverAddEditCnt() > 4 {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_EDIT_ADDRESSS_COUNT)

		return easygo.NewFailMsg(ORDER_EDIT_ADDRESSS_COUNT)
	}
	var receiveAddress share_message.ReceiveAddress = share_message.ReceiveAddress{
		Name:          easygo.NewString(reqMsg.GetAddress().GetName()),
		Phone:         easygo.NewString(reqMsg.GetAddress().GetPhone()),
		Region:        easygo.NewString(reqMsg.GetAddress().GetRegion()),
		DetailAddress: easygo.NewString(reqMsg.GetAddress().GetDetailAddress()),
	}
	e = col.Update(bson.M{"_id": reqMsg.GetOrderId()},
		bson.M{"$set": bson.M{"receiveAddress": &receiveAddress,
			"receiver_addEditCnt": easygo.NewInt32(order.GetReceiverAddEditCnt() + 1)}})

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	//通知
	easygo.Spawn(func(orderPara *share_message.TableShopOrder) {
		if nil != orderPara {

			var content string = MESSAGE_TO_SELLER_ADDCHANGE
			typeValue := share_message.BuySell_Type_Seller

			ShopInstance.InsMessageNotify(
				easygo.NewString(content),
				&typeValue,
				orderPara)

			var jgContent string = MESSAGE_TO_SELLER_ADDCHANGE_PUSH

			ShopInstance.JGMessageNotify(jgContent, orderPara.GetReceiverId(), orderPara.GetOrderId(), typeValue)

		} else {
			logs.Debug("买家修改地址后发通知,缺少订单")
		}
	}, order)

	return &share_message.EditOrderAddressResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(ORDER_EDIT_ADDRESSS_SUCCESS)}
}

func (self *ServiceForHall) RpcEditDeliverAddress(common *base.Common, reqMsg *share_message.EditDeliverAddress) easygo.IMessage {
	logs.Info("===api RpcEditDeliverAddress===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	var order *share_message.TableShopOrder = &share_message.TableShopOrder{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()

	e := col.Find(bson.M{"_id": reqMsg.GetOrderId()}).One(order)

	if e == mgo.ErrNotFound {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_EDIT_DELIVER_ADDRESSS_ORDER_NOT_FOUND)

		return easygo.NewFailMsg(ORDER_EDIT_DELIVER_ADDRESSS_ORDER_NOT_FOUND)
	}

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	if order.GetState() != for_game.SHOP_ORDER_WAIT_PAY && order.GetState() != for_game.SHOP_ORDER_WAIT_SEND {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", ORDER_EDIT_DELIVER_ADDRESSS_STATE_ERROR)

		return easygo.NewFailMsg(ORDER_EDIT_DELIVER_ADDRESSS_STATE_ERROR)
	}

	var deliverAddress share_message.DeliverAddress = share_message.DeliverAddress{
		Name:          easygo.NewString(reqMsg.GetAddress().GetName()),
		Phone:         easygo.NewString(reqMsg.GetAddress().GetPhone()),
		Region:        easygo.NewString(reqMsg.GetAddress().GetRegion()),
		DetailAddress: easygo.NewString(reqMsg.GetAddress().GetDetailAddress()),
	}
	e = col.Update(bson.M{"_id": reqMsg.GetOrderId()},
		bson.M{"$set": bson.M{"deliverAddress": deliverAddress}})

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	easygo.Spawn(func(itemId int64, deliverAddress *share_message.DeliverAddress) {

		colItem, closeFunItem := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
		e = colItem.Update(bson.M{"_id": itemId},
			bson.M{"$set": bson.M{"user_name": deliverAddress.GetName(),
				"phone":          deliverAddress.GetPhone(),
				"region":         deliverAddress.GetRegion(),
				"detail_address": deliverAddress.GetDetailAddress(),
			}})
		closeFunItem()

		if e != nil {
			logs.Error(e)
		}

	}, order.GetItems().GetItemId(), reqMsg.GetAddress())

	return &share_message.EditOrderAddressResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(ORDER_EDIT_DELIVER_ADDRESSS_SUCCESS)}
}

func (self *ServiceForHall) RpcDelayReceiveItem(common *base.Common, reqMsg *share_message.OrderID) easygo.IMessage {
	logs.Info("===api RpcDelayReceiveItem===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	//以每个订单为单位取得锁,取不到说明订单在变化中直接返回
	lockKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_RECEIVE_MUTEX, reqMsg.GetOrderId())
	//取得分布式锁,跟主动收货互斥（失效时间设置10秒）
	errLock := easygo.RedisMgr.GetC().DoRedisLockNoRetry(lockKey, 10)
	defer easygo.RedisMgr.GetC().DoRedisUnlock(lockKey)

	//如果未取得锁
	if errLock != nil {
		s := fmt.Sprintf("RpcDelayReceiveItem 单key取得redis分布式无重试锁失败,redis key is %v", lockKey)
		logs.Error(s)
		logs.Error(errLock)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DELAY_RECEIVE_ORDER_STATE_ERROR)

		return easygo.NewFailMsg(DELAY_RECEIVE_ORDER_STATE_ERROR)
	}

	var order *share_message.TableShopOrder = &share_message.TableShopOrder{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()

	e := col.Find(bson.M{"_id": reqMsg.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_RECEIVE}).One(order)

	if e == mgo.ErrNotFound {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DELAY_RECEIVE_ORDER_STATE_ERROR)

		return easygo.NewFailMsg(DELAY_RECEIVE_ORDER_STATE_ERROR)
	}

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	if order.GetState() != for_game.SHOP_ORDER_WAIT_RECEIVE {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DELAY_RECEIVE_ORDER_STATE_ERROR)

		return easygo.NewFailMsg(DELAY_RECEIVE_ORDER_STATE_ERROR)
	}

	if order.GetReceiveTime() > time.Now().Unix()+3*24*3600 {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DELAY_RECEIVE_TIME_NOT_ARRIVAL)

		return easygo.NewFailMsg(DELAY_RECEIVE_TIME_NOT_ARRIVAL)
	}
	if order.GetDelayReceive() == 1 {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DELAY_RECEIVE_REPEATE)

		return easygo.NewFailMsg(DELAY_RECEIVE_REPEATE)
	}

	newReceiveTime := order.GetReceiveTime() + 7*24*3600

	e = col.Update(
		bson.M{"_id": reqMsg.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_RECEIVE},
		bson.M{"$set": bson.M{"receive_time": newReceiveTime,
			"delay_receive": 1}})

	if e == mgo.ErrNotFound {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DELAY_RECEIVE_ORDER_STATE_ERROR)

		return easygo.NewFailMsg(DELAY_RECEIVE_ORDER_STATE_ERROR)
	}

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	return &share_message.DelayReceiveResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(DELAY_RECEIVE_SUCCESS)}
}
func (self *ServiceForHall) RpcConfirmReceive(common *base.Common, reqMsg *share_message.OrderID) easygo.IMessage {
	logs.Info("===api RpcConfirmReceive===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	//以每个订单为单位取得锁,取不到说明订单在变化中直接返回
	lockKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_RECEIVE_MUTEX, reqMsg.GetOrderId())
	//取得分布式锁,跟主动收货互斥（失效时间设置10秒）
	errLock := easygo.RedisMgr.GetC().DoRedisLockNoRetry(lockKey, 10)
	defer easygo.RedisMgr.GetC().DoRedisUnlock(lockKey)

	//如果未取得锁
	if errLock != nil {
		s := fmt.Sprintf("RpcConfirmReceive 单key取得redis分布式无重试锁失败,redis key is %v", lockKey)
		logs.Error(s)
		logs.Error(errLock)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", CONFIRM_RECEIVE_STATE_ERROR)

		return easygo.NewFailMsg(CONFIRM_RECEIVE_STATE_ERROR)
	}

	var order *share_message.TableShopOrder = &share_message.TableShopOrder{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()

	e := col.Find(bson.M{"_id": reqMsg.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_RECEIVE}).One(order)

	if e == mgo.ErrNotFound {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", CONFIRM_RECEIVE_STATE_ERROR)

		return easygo.NewFailMsg(CONFIRM_RECEIVE_STATE_ERROR)
	}

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	if order.GetState() == for_game.SHOP_ORDER_FINISH {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", "重复操作")

		return easygo.NewFailMsg("重复操作")
	}

	if order.GetState() != for_game.SHOP_ORDER_WAIT_RECEIVE {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", CONFIRM_RECEIVE_STATE_ERROR)

		return easygo.NewFailMsg(CONFIRM_RECEIVE_STATE_ERROR)
	}

	var nowTime int64 = time.Now().Unix()
	e = col.Update(
		bson.M{"_id": reqMsg.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_RECEIVE},
		bson.M{"$set": bson.M{"state": for_game.SHOP_ORDER_FINISH,
			"receive_time":         nowTime,
			"finish_time":          nowTime,
			"receiver_notify_flag": true,
			"sponsor_notify_flag":  true,
			"update_time":          time.Now().Unix(),
		}})

	if e != nil && e != mgo.ErrNotFound {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	if e == mgo.ErrNotFound {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", CONFIRM_RECEIVE_STATE_ERROR)

		return easygo.NewFailMsg(CONFIRM_RECEIVE_STATE_ERROR)
	}

	easygo.Spawn(func() {
		for_game.MakePlayerBehaviorReport(
			4,
			0,
			nil,
			order,
			nil,
			nil) //生成用户行为报表商城订单完成相关字段 已优化到Redis

		// for_game.MakeOperationChannelReport(
		// 	4,
		// 	order.GetReceiverId(),
		// 	"",
		// 	order,
		// 	nil) //生成运营渠道数据汇总报表 已优化到Redis
	})

	SendMsgToServerNewEx(order.GetSponsorId(), "RpcShopPaySeller", &share_message.PaySellerInfo{
		OrderId:    order.OrderId,
		Money:      easygo.NewInt32(order.Items.GetPrice() * order.Items.GetCount()),
		Sponsor_Id: order.SponsorId,
		ReceiverId: order.ReceiverId,
		PayType:    easygo.NewInt32(0)})
	//通知
	easygo.Spawn(func(orderPara *share_message.TableShopOrder) {
		if nil != orderPara {

			var content string = MESSAGE_TO_BUYER_SIGN
			typeValue := share_message.BuySell_Type_Buyer

			ShopInstance.InsMessageNotify(
				easygo.NewString(content),
				&typeValue,
				orderPara)

			var jgContent string = MESSAGE_TO_BUYER_SIGN_PUSH

			ShopInstance.JGMessageNotify(jgContent, orderPara.GetReceiverId(), orderPara.GetOrderId(), typeValue)

			//确认收货 商城订单红点推送

			SendMsgToHallClientNew([]int64{orderPara.GetReceiverId()}, "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
				OrderId: easygo.NewInt64(orderPara.GetOrderId())})

			/*	SendToPlayer(orderPara.GetReceiverId(), "RpcShopOrderNotify",
				&share_message.ShopOrderNotifyInfoWithWho{
					PlayerId: easygo.NewInt64(orderPara.GetReceiverId()),
					OrderId:  easygo.NewInt64(orderPara.GetOrderId()),
				})*/

			//确认收货 商城订单红点推送

			SendMsgToHallClientNew([]int64{orderPara.GetSponsorId()}, "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
				OrderId: easygo.NewInt64(orderPara.GetOrderId())})

			/*	SendToPlayer(orderPara.GetSponsorId(), "RpcShopOrderNotify",
				&share_message.ShopOrderNotifyInfoWithWho{
					PlayerId: easygo.NewInt64(orderPara.GetSponsorId()),
					OrderId:  easygo.NewInt64(orderPara.GetOrderId()),
				})*/

		} else {
			logs.Debug("确认收货后发通知,缺少订单")
		}
	}, order)

	easygo.Spawn(func(orderPara *share_message.TableShopOrder) {
		if nil != orderPara {

			var content string = MESSAGE_TO_SELLER_SIGN
			typeValue := share_message.BuySell_Type_Seller

			ShopInstance.InsMessageNotify(
				easygo.NewString(content),
				&typeValue,
				orderPara)

			var jgContent string = MESSAGE_TO_SELLER_SIGN_PUSH

			ShopInstance.JGMessageNotify(jgContent, orderPara.GetSponsorId(), orderPara.GetOrderId(), typeValue)

			//修改该商品完成的订单数
			AddFinOrderCnt(orderPara.GetItems().GetItemId())
		} else {
			logs.Debug("确认收货后发通知,缺少订单")
		}
	}, order)

	return &share_message.ConfirmReceiveResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(CONFIRM_RECEIVE_SUCCESS)}
}

func (self *ServiceForHall) RpcShopItemEvaluteUpload(common *base.Common, reqMsg *share_message.UploadEvalute) easygo.IMessage {
	logs.Info("===api RpcShopItemEvaluteUpload===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	if reqMsg.GetContent() == "" {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", UPLOAD_ITEM_EVALUTE_NOT_BLANK)

		return easygo.NewFailMsg(UPLOAD_ITEM_EVALUTE_NOT_BLANK)
	}

	//商品评价内容审核
	var texts []string = []string{*reqMsg.Content}

	rstText, rstErrCode, rstErrContent := greenScan.GetTextScanResult(texts)

	if rstText == 0 {

		logs.Debug("商品评价内容审核失败")
		if "" != rstErrCode {
			s := fmt.Sprintf("评价审核失败-玩家id:%v;商品id:%v,;留言内容:%v;错误码:%v;错误内容:%v",
				common.GetUserId(), reqMsg.GetItemId(), texts, rstErrCode, rstErrContent)
			for_game.WriteFile("shop_audit.log", s)
		}
		SendToHallServerByApi(common.GetUserId(), "RpcToast", UPLOAD_ITEM_EVALUTE_AUDIT_FAIL)

		return easygo.NewFailMsg(UPLOAD_ITEM_EVALUTE_AUDIT_FAIL)

	} else if rstText == 2 {

		error := fmt.Sprintf("留言内容验证网络出错%v", reqMsg.GetItemId())
		logs.Error(error)
		for_game.WriteFile("shop_audit.log", error)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	order := share_message.TableShopOrder{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	e := col.Find(bson.M{"_id": reqMsg.GetOrderId()}).One(&order)
	closeFun()

	if e == mgo.ErrNotFound {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", UPLOAD_ITEM_REPEATED)

		return easygo.NewFailMsg(UPLOAD_ITEM_REPEATED)
	}

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	if order.GetState() == for_game.SHOP_ORDER_EVALUTE {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", UPLOAD_ITEM_REPEATED)

		return easygo.NewFailMsg(UPLOAD_ITEM_REPEATED)
	}

	col, closeFun = MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	e = col.Update(
		bson.M{"_id": reqMsg.GetOrderId()},
		bson.M{"$set": bson.M{"state": for_game.SHOP_ORDER_EVALUTE, "update_time": time.Now().Unix()}})
	closeFun()

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	var timeNow int64 = time.Now().Unix()
	commentId := easygo.NewInt64(for_game.NextId(for_game.TABLE_ITEM_COMMENT))
	who := for_game.GetRedisPlayerBase(common.GetUserId())
	newItem := &share_message.TableItemComment{
		CommentId:     commentId,
		ItemId:        easygo.NewInt64(reqMsg.GetItemId()),
		PlayerId:      easygo.NewInt64(common.GetUserId()),
		Nickname:      easygo.NewString(who.GetNickName()),
		Avatar:        easygo.NewString(who.GetHeadIcon()),
		Content:       easygo.NewString(reqMsg.GetContent()),
		Sex:           easygo.NewInt32(who.GetSex()),
		CreateTime:    easygo.NewInt64(timeNow),
		StarLevel:     easygo.NewInt32(reqMsg.GetStarLevel()),
		RealLikeCount: easygo.NewInt32(0),
		FakeLikeCount: easygo.NewInt32(0),
		Status:        easygo.NewInt32(for_game.SHOP_COMMENT_NO_REPLY),
	}

	col, closeFun = MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ITEM_COMMENT)
	defer closeFun()

	e = col.Insert(newItem)

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	// 更新商品表相关的留言总数,完成评价数，好评数
	easygo.Spawn(func(itemIdPara int64, startLevel int32) {

		if itemIdPara != 0 {

			// 修改商品表中留言总数
			AddAllCommentCnt(itemIdPara)

			// 修改该商品的真实已经完成的评价数
			AddFinCommCnt(itemIdPara)

			// 修改商品表中好评数
			if startLevel == 3 {
				AddGoodCommCnt(itemIdPara)
			}

		} else {
			logs.Debug("更新商品表相关的留言总数,缺少商品ID")
		}

	}, reqMsg.GetItemId(), reqMsg.GetStarLevel())

	return &share_message.UploadCommentResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(UPLOAD_ITEM_EVALUTE_SUCCESS)}
}

func (self *ServiceForHall) RpcShopMessageFlgUpd(common *base.Common, reqMsg *share_message.MessageIdList) easygo.IMessage {
	logs.Info("===api RpcShopMessageFlgUpd===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_MESSAGE)
	defer closeFun()
	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": reqMsg.GetMessageIds()}}, bson.M{"$set": bson.M{"view_flag": true}})
	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
	}
	return easygo.EmptyMsg
}

func (self *ServiceForHall) RpcShopItemMessage(common *base.Common, reqMsg *share_message.ShopItemMessageListInfo) easygo.IMessage {
	logs.Info("===api RpcShopItemMessage===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_MESSAGE)
	defer closeFun()
	var list []*share_message.TableShopMessage

	if reqMsg.Type != nil {

		//0买家（我购买的）
		if share_message.BuySell_Type_Buyer == reqMsg.GetType() {
			query := col.Find(bson.M{
				"receiver_player_id": common.GetUserId(),
				"user_type":          int32(share_message.BuySell_Type_Buyer),
				"create_time":        bson.M{"$gt": reqMsg.GetIncrementTime()},
				"view_flag":          false})
			cnt, err := query.Count()

			if err != nil && err != mgo.ErrNotFound {
				logs.Error(err)
				return &share_message.ShopItemMessageList{IncrementFlag: easygo.NewBool(false)}
			}

			if cnt <= 0 {
				return &share_message.ShopItemMessageList{IncrementFlag: easygo.NewBool(false)}
			}

			e := col.Find(
				bson.M{
					"receiver_player_id": common.GetUserId(),
					"user_type":          int32(share_message.BuySell_Type_Buyer),
					"create_time":        bson.M{"$gt": reqMsg.GetIncrementTime()}}).Sort("-create_time").All(&list)

			newList := []*share_message.ShopItemMessage{}

			if e != mgo.ErrNotFound && e != nil {
				logs.Error(e)
				return &share_message.ShopItemMessageList{IncrementFlag: easygo.NewBool(false)}
			}

			if e == nil && list != nil {

				var tempType share_message.BuySell_Type = share_message.BuySell_Type_Buyer
				for _, value := range list {
					newItem := share_message.ShopItemMessage{
						MessageId:  value.MessageId,
						Type:       &tempType,
						File:       value.File,
						Nickname:   value.SponsorNickname,
						Avatar:     value.SponsorAvatar,
						ItemName:   value.ItemName,
						ItemTitle:  value.ItemTitle,
						Content:    value.Content,
						CreateTime: value.CreateTime,
						OrderId:    value.OrderId,
						ShowTime:   easygo.NewString(util.FormatUnixTime(value.GetCreateTime())),
						ViewFlag:   value.ViewFlag,
					}
					newList = append(newList, &newItem)
				}

				return &share_message.ShopItemMessageList{List: newList, IncrementFlag: easygo.NewBool(true)}
			} else {

				return &share_message.ShopItemMessageList{IncrementFlag: easygo.NewBool(false)}

			}

			//1卖家（我卖出的）
		} else {
			query := col.Find(bson.M{
				"sponsor_player_id": common.GetUserId(),
				"user_type":         int32(share_message.BuySell_Type_Seller),
				"create_time":       bson.M{"$gt": reqMsg.GetIncrementTime()},
				"view_flag":         false})
			cnt, err := query.Count()

			if err != nil && err != mgo.ErrNotFound {
				logs.Error(err)
				return &share_message.ShopItemMessageList{IncrementFlag: easygo.NewBool(false)}
			}

			if cnt <= 0 {
				return &share_message.ShopItemMessageList{IncrementFlag: easygo.NewBool(false)}
			}

			e := col.Find(
				bson.M{
					"sponsor_player_id": common.GetUserId(),
					"user_type":         int32(share_message.BuySell_Type_Seller),
					"create_time":       bson.M{"$gt": reqMsg.GetIncrementTime()}}).Sort("-create_time").All(&list)

			if e != mgo.ErrNotFound && e != nil {
				logs.Error(e)
				return &share_message.ShopItemMessageList{IncrementFlag: easygo.NewBool(false)}
			}

			newList := []*share_message.ShopItemMessage{}

			if e == nil && list != nil {

				var tempType share_message.BuySell_Type = share_message.BuySell_Type_Seller
				for _, value := range list {
					newItem := share_message.ShopItemMessage{
						MessageId:  value.MessageId,
						Type:       &tempType,
						File:       value.File,
						Nickname:   value.ReceiverNickname,
						Avatar:     value.ReceiverAvatar,
						ItemName:   value.ItemName,
						ItemTitle:  value.ItemTitle,
						Content:    value.Content,
						CreateTime: value.CreateTime,
						OrderId:    value.OrderId,
						ShowTime:   easygo.NewString(GetYMDTime(value.GetCreateTime())),
						ViewFlag:   value.ViewFlag,
					}
					newList = append(newList, &newItem)
				}

				return &share_message.ShopItemMessageList{
					List:          newList,
					IncrementFlag: easygo.NewBool(true)}

			} else {
				return &share_message.ShopItemMessageList{IncrementFlag: easygo.NewBool(false)}
			}
		}

	} else {
		return &share_message.ShopItemMessageList{IncrementFlag: easygo.NewBool(false)}
	}
}

func (self *ServiceForHall) RpcNotifySendItem(common *base.Common, reqMsg *share_message.OrderID) easygo.IMessage {
	logs.Info("===api RpcNotifySendItem===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	shopOrder := share_message.TableShopOrder{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	e := col.Find(bson.M{"_id": reqMsg.OrderId}).Limit(1).One(&shopOrder)
	closeFun()

	if e == nil {
		var content string = MESSAGE_TO_SELLER_REMIND
		typeValue := share_message.BuySell_Type_Seller

		ShopInstance.InsMessageNotify(
			easygo.NewString(content),
			&typeValue,
			&shopOrder)

	} else {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", NOTIFY_USER_FAIL)

		return easygo.NewFailMsg(NOTIFY_USER_FAIL)
	}

	//push推送

	easygo.Spawn(func(orderPara share_message.TableShopOrder) {

		typeValue := share_message.BuySell_Type_Seller

		var jgContent string = MESSAGE_TO_SELLER_REMIND_PUSH

		ShopInstance.JGMessageNotify(jgContent, orderPara.GetSponsorId(), orderPara.GetOrderId(), typeValue)

	}, shopOrder)

	return &share_message.NotifySendItemResult{
		Result: easygo.NewInt32(0),
		Msg:    easygo.NewString(NOTIFY_USER__SHIPPING_SUCCESS)}
}

func (self *ServiceForHall) RpcGetShopOrderNotifyInfos(common *base.Common, reqMsg *share_message.PlayerID) easygo.IMessage {
	logs.Info("===api RpcGetShopOrderNotifyInfos===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()
	var list []*share_message.TableShopOrder

	err := col.Find(bson.M{
		"$or": []bson.M{
			bson.M{"sponsor_id": reqMsg.GetPlayerId(), "sponsor_notify_flag": true, "delete_sell": 0},
			bson.M{"receiver_id": reqMsg.GetPlayerId(), "receiver_notify_flag": true, "delete_buy": 0},
		}}).All(&list)

	if err != nil {
		logs.Error(err)
		return easygo.EmptyMsg
	}

	orderList := []int64{}
	if nil != list && len(list) > 0 {
		for _, value := range list {
			orderList = append(orderList, value.GetOrderId())
		}
	}
	return &share_message.ShopOrderIdList{OrderIds: orderList}
}

func (self *ServiceForHall) RpcShopOrderNotifyFlgUpd(common *base.Common, reqMsg *share_message.ShopOrderNotifyFlgUpdInfo) easygo.IMessage {
	logs.Info("===api RpcShopOrderNotifyFlgUpd===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()
	if share_message.BuySell_Type_Buyer == reqMsg.GetBuySell_Type() {
		_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": reqMsg.GetOrderIds()}},
			bson.M{"$set": bson.M{"receiver_notify_flag": false}})
		if err != nil && err != mgo.ErrNotFound {
			logs.Error(err)
		}
	} else {
		_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": reqMsg.GetOrderIds()}},
			bson.M{"$set": bson.M{"sponsor_notify_flag": false}})
		if err != nil && err != mgo.ErrNotFound {
			logs.Error(err)
		}
	}

	return easygo.EmptyMsg
}

func (self *ServiceForHall) RpcShopUploadAuth(common *base.Common, reqMsg *share_message.PlayerID) easygo.IMessage {
	logs.Info("===api RpcShopUploadAuth===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	var shopAuth = share_message.TableShopPlayer{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_PLAYER)
	e := col.Find(bson.M{"_id": reqMsg.PlayerId}).Limit(1).One(&shopAuth)
	closeFun()

	if e == mgo.ErrNotFound {
		var uploadAuthFlag int32 = 0
		return &share_message.UploadAuthResult{UploadAuthFlag: &uploadAuthFlag}
	}

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	if shopAuth.GetUploadAuthFlag() == 1 {
		var uploadAuthFlag int32 = 1
		return &share_message.UploadAuthResult{UploadAuthFlag: &uploadAuthFlag}
	}

	var uploadAuthFlag int32 = 0
	return &share_message.UploadAuthResult{UploadAuthFlag: &uploadAuthFlag}
}

func (self *ServiceForHall) RpcShopUploadAuthConfirm(common *base.Common, reqMsg *share_message.PlayerID) easygo.IMessage {
	logs.Info("===api RpcShopUploadAuthConfirm===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	shopAuth := share_message.TableShopPlayer{
		PlayerId:            reqMsg.PlayerId,
		UploadAuthFlag:      easygo.NewInt32(1),
		CreateTime:          easygo.NewInt64(time.Now().Unix()),
		FakePlayFinOrderCnt: easygo.NewInt32(0)}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_PLAYER)
	_, e := col.Upsert(
		bson.M{"_id": reqMsg.GetPlayerId()},
		bson.M{"$set": shopAuth})
	closeFun()

	if e != nil {
		logs.Error(e)
		SendToHallServerByApi(common.GetUserId(), "RpcToast", DATABASE_ERROR)

		return easygo.NewFailMsg(DATABASE_ERROR)
	}

	return easygo.EmptyMsg
}

func (self *ServiceForHall) RpcSettlementBtn(common *base.Common, reqMsg *share_message.SettlementInfo) easygo.IMessage {
	logs.Info("===api RpcSettlementBtn===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	if len(reqMsg.Items) == 0 {
		SendToHallServerByApi(common.GetUserId(), "RpcToast", SETTLEMENT_BUY_NULL)

		return easygo.NewFailMsg(SETTLEMENT_BUY_NULL)
	}

	var noSaleMessages []string = []string{}
	var noStockMessages []string = []string{}
	var blackMessages []string = []string{}

	for _, value := range reqMsg.Items {
		shopItem := ShopInstance.GetItemFromCache(value.GetItemId())

		if shopItem == nil {
			noSaleMessages = append(noSaleMessages, fmt.Sprintf(SETTLEMENT_ITEM_PARA_NO_SALE,
				value.GetItemName()))
		}
	}

	if len(noSaleMessages) > 0 {
		return &share_message.SettlementResult{
			Result:         easygo.NewInt32(1),
			NoSaleMessages: noSaleMessages}
	}

	for _, value := range reqMsg.Items {
		shopItem := ShopInstance.GetItemFromCache(value.GetItemId())

		if shopItem.stock_count < value.GetCount() {

			stockCount := shopItem.stock_count
			noStockMessages = append(noStockMessages,
				fmt.Sprintf(SETTLEMENT_ITEM_NO_STOCK,
					shopItem.name,
					easygo.AnytoA(int64(stockCount))))
		}

		// 黑名单判断
		var blackList []PLAYER_ID = ShopInstance.GetBlackLists(common.GetUserId())
		var blackFlag int32 = 0
		for _, black := range blackList {
			if shopItem.player_id == black {
				blackFlag = 1
				break
			}
		}
		//显示黑名单
		if blackFlag == 1 {
			blackMessages = append(blackMessages, fmt.Sprintf(SETTLEMENT_ITEM_BLACK,
				shopItem.name))
		}
	}

	if len(noStockMessages) > 0 {

		return &share_message.SettlementResult{
			Result:          easygo.NewInt32(1),
			NoStockMessages: noStockMessages}
	}

	if len(blackMessages) > 0 {

		return &share_message.SettlementResult{
			Result:        easygo.NewInt32(1),
			BlackMessages: blackMessages}
	}

	return &share_message.SettlementResult{
		Result: easygo.NewInt32(0)}
}
func (list LastShopItemsComDesc) Len() int { return len(list) }
func (list LastShopItemsComDesc) Swap(i, j int) {
	s := list[j]
	list[j] = list[i]
	list[i] = s
}

func (list LastShopItemsComDesc) Less(i, j int) bool {
	return list[i].realCommentCnt > list[j].realCommentCnt
}

func (list LastShopItemsStoreDesc) Len() int { return len(list) }
func (list LastShopItemsStoreDesc) Swap(i, j int) {
	s := list[j]
	list[j] = list[i]
	list[i] = s
}

func (list LastShopItemsStoreDesc) Less(i, j int) bool {
	return list[i].realStoreCnt > list[j].realStoreCnt
}

func (list LastShopItemsSalesDesc) Len() int { return len(list) }
func (list LastShopItemsSalesDesc) Swap(i, j int) {
	s := list[j]
	list[j] = list[i]
	list[i] = s
}

func (list LastShopItemsSalesDesc) Less(i, j int) bool {
	return list[i].realFinOrderCnt > list[j].realFinOrderCnt
}

func (list LastShopItemsNew) Len() int { return len(list) }
func (list LastShopItemsNew) Swap(i, j int) {
	s := list[j]
	list[j] = list[i]
	list[i] = s
}

func (list LastShopItemsNew) Less(i, j int) bool {
	return list[i].create_time > list[j].create_time
}

func (list LastShopItemsPriceDesc) Len() int { return len(list) }
func (list LastShopItemsPriceDesc) Swap(i, j int) {
	s := list[j]
	list[j] = list[i]
	list[i] = s
}

func (list LastShopItemsPriceDesc) Less(i, j int) bool {
	return list[i].price < list[j].price
}

func (list LastShopItemsPriceAsc) Len() int { return len(list) }
func (list LastShopItemsPriceAsc) Swap(i, j int) {
	s := list[j]
	list[j] = list[i]
	list[i] = s
}

func (list LastShopItemsPriceAsc) Less(i, j int) bool {
	return list[i].price > list[j].price
}
