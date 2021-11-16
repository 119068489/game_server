// 管理后台为[浏览器]提供的服务

package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

const (
	NOTIFY_SHOP_OPERATE_1 int32 = 1 //1取消订单
	NOTIFY_SHOP_OPERATE_2 int32 = 2 //2发货
	NOTIFY_SHOP_OPERATE_5 int32 = 5 //5取得物流信息

)

//后台导航跳转以及查询商城商品列表
func (self *cls4) RpcQueryShopItem(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryShopItemRequest) easygo.IMessage {
	list, count := QueryShopItemList(reqMsg)
	msg := &brower_backstage.QueryShopItemResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//点击列表中的商品id跳转商品详情页面
func (self *cls4) RpcQueryShopItemDetailById(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	logs.Info("====RpcQueryShopItemDetailById=====,reqMsg=%v", reqMsg)
	id := reqMsg.GetId64()
	if id == 0 {
		return easygo.NewFailMsg("商品id不能为空")
	}
	result := QueryShopItemDetailById(id)
	if result == nil {
		return easygo.NewFailMsg("商品不存在！")
	}

	return result
}

//操作中下架按钮的确定
func (self *cls4) RpcShopSoldOut(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	id := reqMsg.GetId64()
	if id == 0 {
		return easygo.NewFailMsg("商品id不能为空")
	}
	ShopSoldOut(id)

	msg := fmt.Sprintf("下架商品：下架商品%d", id)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SHOP_MANAGE, msg)

	return easygo.EmptyMsg
}

//商品发布页面和修改页面上 商品分类和品类标签的下拉内容取得
func (self *cls4) RpcGetShopItemTypeDropDown(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {

	var dropDownType []*brower_backstage.KeyValueTag = []*brower_backstage.KeyValueTag{}
	var dropDownCategory []*brower_backstage.KeyValue = []*brower_backstage.KeyValue{}

	//=======设置商品分类开始==========================
	//对应手机客户端的全部  的时候是0
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(0),
		Value: easygo.NewString("请选择"),
	})

	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(1),
		Value: easygo.NewString("手机"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(2),
		Value: easygo.NewString("农用物资"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(3),
		Value: easygo.NewString("生鲜水果"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(4),
		Value: easygo.NewString("童鞋"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(5),
		Value: easygo.NewString("园艺植物"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(6),
		Value: easygo.NewString("五金工具"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(7),
		Value: easygo.NewString("游戏"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(8),
		Value: easygo.NewString("电子零件"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(9),
		Value: easygo.NewString("动漫/周边"),
	})

	//为了应对通联支付审核,目前页面删除,但不改变其他分类的type值
	//dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
	//	Key: easygo.NewInt32(10),
	//	Value: easygo.NewString("图书"),
	//})
	//dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
	//	Key: easygo.NewInt32(11),
	//	Value: easygo.NewString("宠物/用品"),
	//})

	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(12),
		Value: easygo.NewString("网络设备"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(13),
		Value: easygo.NewString("服饰配件"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(14),
		Value: easygo.NewString("家装/建材"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(15),
		Value: easygo.NewString("家纺布艺"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(16),
		Value: easygo.NewString("珠宝首饰"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(17),
		Value: easygo.NewString("钟表眼镜"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(18),
		Value: easygo.NewString("古董收藏"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(19),
		Value: easygo.NewString("女士鞋靴"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(20),
		Value: easygo.NewString("箱包"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(21),
		Value: easygo.NewString("男士鞋靴"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(22),
		Value: easygo.NewString("办公用品"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(23),
		Value: easygo.NewString("游戏设备"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(24),
		Value: easygo.NewString("运动户外"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(25),
		Value: easygo.NewString("实体卡/券/票"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(26),
		Value: easygo.NewString("工艺礼品"),
	})
	//为了应对通联支付审核,目前页面删除,但不改变其他分类的type值
	//dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
	//	Key: easygo.NewInt32(27),
	//	Value: easygo.NewString("玩具乐器"),
	//})
	//dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
	//	Key: easygo.NewInt32(28),
	//	Value: easygo.NewString("母婴用品"),
	//})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(29),
		Value: easygo.NewString("童装"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(30),
		Value: easygo.NewString("女士服装"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(31),
		Value: easygo.NewString("运动户外"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(32),
		Value: easygo.NewString("居家用品"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(33),
		Value: easygo.NewString("家用电器"),
	})

	//为了应对通联支付审核,目前页面删除,但不改变其他分类的type值
	//dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
	//	Key: easygo.NewInt32(34),
	//	Value: easygo.NewString("个护美妆"),
	//})
	//dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
	//	Key: easygo.NewInt32(35),
	//	Value: easygo.NewString("保健护理"),
	//})

	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(36),
		Value: easygo.NewString("摩托车/用品"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(37),
		Value: easygo.NewString("自行车/用品"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(38),
		Value: easygo.NewString("汽车/用品"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(39),
		Value: easygo.NewString("电动车/用品"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(40),
		Value: easygo.NewString("3C数码"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(41),
		Value: easygo.NewString("男士服装"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(42),
		Value: easygo.NewString("其他闲置"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(43),
		Value: easygo.NewString("音像"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(44),
		Value: easygo.NewString("演艺/表演类门票"),
	})
	dropDownType = append(dropDownType, &brower_backstage.KeyValueTag{
		Key:   easygo.NewInt32(45),
		Value: easygo.NewString("点卡"),
	})
	//=======设置商品分类结束==========================

	//=======设置品类标签开始=========================
	dropDownCategory = append(dropDownCategory, &brower_backstage.KeyValue{
		Key:   easygo.NewString("请选择"),
		Value: easygo.NewString("请选择"),
	})
	dropDownCategory = append(dropDownCategory, &brower_backstage.KeyValue{
		Key:   easygo.NewString("电影票"),
		Value: easygo.NewString("电影票"),
	})
	dropDownCategory = append(dropDownCategory, &brower_backstage.KeyValue{
		Key:   easygo.NewString("口红"),
		Value: easygo.NewString("口红"),
	})
	dropDownCategory = append(dropDownCategory, &brower_backstage.KeyValue{
		Key:   easygo.NewString("运动鞋"),
		Value: easygo.NewString("运动鞋"),
	})
	dropDownCategory = append(dropDownCategory, &brower_backstage.KeyValue{
		Key:   easygo.NewString("女鞋"),
		Value: easygo.NewString("女鞋"),
	})
	dropDownCategory = append(dropDownCategory, &brower_backstage.KeyValue{
		Key:   easygo.NewString("篮球鞋"),
		Value: easygo.NewString("篮球鞋"),
	})
	dropDownCategory = append(dropDownCategory, &brower_backstage.KeyValue{
		Key:   easygo.NewString("耳机"),
		Value: easygo.NewString("耳机"),
	})
	dropDownCategory = append(dropDownCategory, &brower_backstage.KeyValue{
		Key:   easygo.NewString("唇釉"),
		Value: easygo.NewString("唇釉"),
	})
	dropDownCategory = append(dropDownCategory, &brower_backstage.KeyValue{
		Key:   easygo.NewString("手机"),
		Value: easygo.NewString("手机"),
	})
	dropDownCategory = append(dropDownCategory, &brower_backstage.KeyValue{
		Key:   easygo.NewString("iphone"),
		Value: easygo.NewString("iphone"),
	})
	dropDownCategory = append(dropDownCategory, &brower_backstage.KeyValue{
		Key:   easygo.NewString("高跟鞋"),
		Value: easygo.NewString("高跟鞋"),
	})
	dropDownCategory = append(dropDownCategory, &brower_backstage.KeyValue{
		Key:   easygo.NewString("卡券"),
		Value: easygo.NewString("卡券"),
	})
	//=======设置品类标签结束=========================

	return &brower_backstage.GetShopItemTypeDropDownResponse{
		DropDownItemType:     dropDownType,
		DropDownItemCategory: dropDownCategory,
	}
}

//发布商品确定按钮
func (self *cls4) RpcReleaseShopItem(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ReleaseEditShopItemObject) easygo.IMessage {

	if reqMsg.GetName() == "" {
		return easygo.NewFailMsg("商品名称不能为空！")
	}
	if reqMsg.GetPlayerAccount() == "" {
		return easygo.NewFailMsg("卖家柠檬账号不能为空！")
	}

	player := QueryPlayerbyAccount(reqMsg.GetPlayerAccount())
	if player == nil {
		return easygo.NewFailMsg("玩家柠檬号不存在!")
	}

	switch player.GetStatus() {
	case for_game.ACCOUNT_USER_FROZEN, for_game.ACCOUNT_ADMIN_FROZEN: // 1 //用户冻结
		return easygo.NewFailMsg("玩家柠檬号已冻结!")
	case for_game.ACCOUNT_CANCELING: // 3 //注销中
		return easygo.NewFailMsg("玩家柠檬号注销中!")
	case for_game.ACCOUNT_CANCELED: // 4 //已注销
		return easygo.NewFailMsg("玩家柠檬号已注销!")
	}

	//页面传过来的价格是乘以100的 按照分计算
	if reqMsg.GetPrice() < 0 || reqMsg.GetPrice() > 500000 {
		return easygo.NewFailMsg("单价必须在0~5000之间!")
	}

	//点卡的时候不判断库存
	if reqMsg.GetItemType() != for_game.SHOP_POINT_CARD_CATEGORY && (reqMsg.GetStockCount() < 1 || reqMsg.GetStockCount() > 9999) {
		return easygo.NewFailMsg("库存必须在1~9999之间!")
	}

	if reqMsg.GetUserName() == "" {
		return easygo.NewFailMsg("发货人真实姓名不能为空!")
	}

	if reqMsg.GetPhone() == "" {
		return easygo.NewFailMsg("发货人手机号不能为空!")
	}

	//点卡的时候不判断地址
	if reqMsg.GetItemType() != for_game.SHOP_POINT_CARD_CATEGORY && reqMsg.GetAddress() == "" {
		return easygo.NewFailMsg("发货地址不能为空!")
	}

	//点卡的时候不判断地址
	if reqMsg.GetItemType() != for_game.SHOP_POINT_CARD_CATEGORY && reqMsg.GetDetailAddress() == "" {
		return easygo.NewFailMsg("详细地址不能为空!")
	}

	if reqMsg.GetTitle() == "" {
		return easygo.NewFailMsg("描述不能为空!")
	}

	if len(reqMsg.GetItemFiles()) < 2 || len(reqMsg.GetItemFiles()) > 9 {
		return easygo.NewFailMsg("广告图片在2~9张之间")
	}

	stockCnt := int32(0)
	if reqMsg.GetItemType() == for_game.SHOP_POINT_CARD_CATEGORY {
		//通过点卡名称和柠檬帐号取得该发布商品的点卡库存
		stockCnt = QueryShopPointCardStock(reqMsg.GetPointCardName(), reqMsg.GetPlayerAccount())
		if stockCnt <= 0 {
			return easygo.NewFailMsg("商品点卡的待售库存不足!")
		}
	}

	//上架判断该商家是否存在相同点卡的商品
	//取得该商家存在库存的点卡的上架中的商品
	if reqMsg.GetState() == for_game.SHOP_ITEM_SALE {
		colItem, closeFunItem := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
		defer closeFunItem()

		var listItem []*share_message.TableShopItem = []*share_message.TableShopItem{}
		errItem := colItem.Find(bson.M{"player_account": reqMsg.GetPlayerAccount(), "state": for_game.SHOP_ITEM_SALE, "point_card_name": reqMsg.GetPointCardName()}).All(&listItem)

		if errItem != nil && errItem != mgo.ErrNotFound {
			easygo.PanicError(errItem)
		}

		if nil != listItem && len(listItem) > 0 {
			return easygo.NewFailMsg("同类点卡已存在上架商品,不能上架!")
		}
	}

	//重新设置type
	var itemType share_message.ShopItemType = share_message.ShopItemType{}
	itemType.Type = easygo.NewInt32(reqMsg.GetItemType())
	//拼接OtherType
	//商品分类名称
	var itemTypeName = GetItemTypeName(reqMsg.GetItemType())
	//品类标签
	var itemCategory = reqMsg.GetItemCategory()

	if itemCategory == "请选择" {
		itemCategory = ""
	}
	//常用选项
	var commonUseLabel []string = reqMsg.GetCommonUseLabel()
	var otherType []string = []string{}
	otherType = append(otherType, itemTypeName)
	otherType = append(otherType, itemCategory)

	if nil != commonUseLabel && len(commonUseLabel) > 0 {

		for _, value := range commonUseLabel {
			otherType = append(otherType, value)
		}
	}
	itemType.OtherType = otherType

	//组装图片url
	itemFiles := []*share_message.ItemFile{}
	if nil != reqMsg.ItemFiles && len(reqMsg.ItemFiles) > 0 {
		for i := 0; i < len(reqMsg.ItemFiles); i++ {
			var itemFile *share_message.ItemFile = &share_message.ItemFile{
				FileUrl:    easygo.NewString(reqMsg.GetItemFiles()[i].GetFileUrl()),
				FileType:   easygo.NewInt32(reqMsg.GetItemFiles()[i].GetFileType()),
				FileWidth:  easygo.NewString(reqMsg.GetItemFiles()[i].GetFileWidth()),
				FileHeight: easygo.NewString(reqMsg.GetItemFiles()[i].GetFileHeight())}
			itemFiles = append(itemFiles, itemFile)
		}
	}
	newItem := &share_message.TableShopItem{
		ItemId:        easygo.NewInt64(for_game.NextId(for_game.TABLE_SHOP_ITEMS)),
		Price:         easygo.NewInt32(reqMsg.GetPrice()),
		Type:          &itemType,
		ItemFiles:     itemFiles,
		Title:         easygo.NewString(reqMsg.GetTitle()),
		UserName:      easygo.NewString(reqMsg.GetUserName()),
		Phone:         easygo.NewString(reqMsg.GetPhone()),
		PlayerId:      easygo.NewInt64(player.GetPlayerId()),
		PlayerAccount: easygo.NewString(player.GetAccount()),
		Nickname:      easygo.NewString(player.GetNickName()),
		Avatar:        easygo.NewString(player.GetHeadIcon()),
		State:         easygo.NewInt32(reqMsg.GetState()),
		Sex:           easygo.NewInt32(player.GetSex()),
		CreateTime:    easygo.NewInt64(time.Now().Unix()),
		LockCount:     easygo.NewInt32(0),
		Name:          easygo.NewString(reqMsg.GetName()),
		RealStoreCnt:  easygo.NewInt32(0),
		FakePayCnt:    easygo.NewInt32(reqMsg.GetFakePaymentCount()),
		FakePageViews: easygo.NewInt32(reqMsg.GetFakePageViews()),
		PointCardName: easygo.NewString(reqMsg.GetPointCardName()),
	}

	//点卡的时候判断
	if reqMsg.GetItemType() != for_game.SHOP_POINT_CARD_CATEGORY {
		newItem.Address = easygo.NewString(reqMsg.GetAddress())
		newItem.DetailAddress = easygo.NewString(reqMsg.GetDetailAddress())
		newItem.StockCount = easygo.NewInt32(reqMsg.GetStockCount())

	} else {
		//点卡的时候从导入库中取得
		newItem.StockCount = easygo.NewInt32(stockCnt)
	}

	//判断状态是下架的时候要设置下架时间
	if reqMsg.GetState() == for_game.SHOP_ITEM_SOLD_OUT {
		newItem.SoldOutTime = easygo.NewInt64(time.Now().Unix())
	}
	if reqMsg.GetFakeSellItemCount() >= 0 {
		shopAuth := share_message.TableShopPlayer{
			PlayerId:            easygo.NewInt64(player.GetPlayerId()),
			UploadAuthFlag:      easygo.NewInt32(1),
			CreateTime:          easygo.NewInt64(time.Now().Unix()),
			FakePlayFinOrderCnt: easygo.NewInt32(reqMsg.GetFakeSellItemCount())}

		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_PLAYER)
		_, e := col.Upsert(
			bson.M{"_id": player.GetPlayerId()},
			bson.M{"$set": shopAuth})
		closeFun()

		if e != nil {
			logs.Error(e)
		}
	}

	msg := fmt.Sprintf("商城管理后台发布商品: %d", newItem.GetItemId())

	ReleaseShopItem(newItem)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SHOP_MANAGE, msg)

	return easygo.EmptyMsg
}

//点击修改按钮跳转取得的商品的数据
func (self *cls4) RpcGetEditShopItemDetailById(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {

	id := reqMsg.GetId64()
	if id == 0 {
		return easygo.NewFailMsg("商品id不能为空")
	}
	result := GetEditShopItemDetailById(id)
	if result == nil {
		return easygo.NewFailMsg("商品不存在！")
	}

	return result
}

//修改商品确定按钮
func (self *cls4) RpcEditShopItem(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ReleaseEditShopItemObject) easygo.IMessage {

	if reqMsg.ItemId == nil || reqMsg.GetItemId() == 0 {
		return easygo.NewFailMsg("商品id不能为空!")
	}
	result := QueryShopItemById(reqMsg.GetItemId())
	if result == nil {
		return easygo.NewFailMsg("商品不存在!")
	}
	if result.GetState() != for_game.SHOP_ITEM_SOLD_OUT {
		return easygo.NewFailMsg("只有下架商品能修改信息!")
	}

	if reqMsg.GetName() == "" {
		return easygo.NewFailMsg("商品名称不能为空！")
	}

	//页面传过来的价格是乘以100的 按照分计算
	if reqMsg.GetPrice() < 0 || reqMsg.GetPrice() > 500000 {
		return easygo.NewFailMsg("单价必须在0~5000之间!")
	}

	//点卡的时候不判断库存
	if reqMsg.GetItemType() != for_game.SHOP_POINT_CARD_CATEGORY && (reqMsg.GetStockCount() < 1 || reqMsg.GetStockCount() > 9999) {
		return easygo.NewFailMsg("库存必须在1~9999之间!")
	}
	//点卡的时候不判断库存
	if reqMsg.GetItemType() != for_game.SHOP_POINT_CARD_CATEGORY && (reqMsg.GetStockCount() < result.GetStockCount()) {
		return easygo.NewFailMsg("库存不能小于现有的库存数量!")
	}

	if reqMsg.GetUserName() == "" {
		return easygo.NewFailMsg("发货人真实姓名不能为空!")
	}

	if reqMsg.GetPhone() == "" {
		return easygo.NewFailMsg("发货人手机号不能为空!")
	}

	//点卡的时候不判断地址
	if reqMsg.GetItemType() != for_game.SHOP_POINT_CARD_CATEGORY && reqMsg.GetAddress() == "" {
		return easygo.NewFailMsg("发货地址不能为空!")
	}

	//点卡的时候不判断地址
	if reqMsg.GetItemType() != for_game.SHOP_POINT_CARD_CATEGORY && reqMsg.GetDetailAddress() == "" {
		return easygo.NewFailMsg("详细地址不能为空!")
	}

	if reqMsg.GetTitle() == "" {
		return easygo.NewFailMsg("描述不能为空!")
	}

	if len(reqMsg.GetItemFiles()) < 2 || len(reqMsg.GetItemFiles()) > 9 {
		return easygo.NewFailMsg("广告图片在2~9张之间")
	}

	stockCnt := int32(0)
	if reqMsg.GetItemType() == for_game.SHOP_POINT_CARD_CATEGORY {
		//通过点卡名称和柠檬帐号取得该发布商品的点卡库存
		stockCnt = QueryShopPointCardStock(reqMsg.GetPointCardName(), reqMsg.GetPlayerAccount())
		if stockCnt <= 0 {
			if reqMsg.GetState() == for_game.SHOP_ITEM_SALE {
				return easygo.NewFailMsg("商品点卡的待售库存不足,不能上架!")
			}
		}
	}

	//上架判断该商家是否存在相同点卡的商品,该商品下是否有待支付商品订单未完成的支付订单
	//取得该商家存在库存的点卡的上架中的商品
	if reqMsg.GetState() == for_game.SHOP_ITEM_SALE {
		colItem, closeFunItem := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
		defer closeFunItem()

		var listItem []*share_message.TableShopItem = []*share_message.TableShopItem{}
		errItem := colItem.Find(bson.M{"player_account": reqMsg.GetPlayerAccount(), "state": for_game.SHOP_ITEM_SALE, "point_card_name": reqMsg.GetPointCardName()}).All(&listItem)

		if errItem != nil && errItem != mgo.ErrNotFound {
			easygo.PanicError(errItem)
		}

		if nil != listItem && len(listItem) > 0 {
			return easygo.NewFailMsg("同类点卡已存在上架商品,不能上架!")
		}

		//该商品下是否有待支付商品订单未完成的支付订单
		flag, tempErr := VerdictItemOrders(reqMsg.GetItemId())
		if tempErr != "" {
			return easygo.NewFailMsg(tempErr)
		} else {
			if flag {
				return easygo.NewFailMsg("当前有待支付且拉起三方支付的订单")
			}
		}
	}

	//重新设置type
	var itemType share_message.ShopItemType = share_message.ShopItemType{}
	itemType.Type = easygo.NewInt32(reqMsg.GetItemType())
	//拼接OtherType
	//商品分类名称
	var itemTypeName = GetItemTypeName(reqMsg.GetItemType())
	//品类标签
	var itemCategory = reqMsg.GetItemCategory()

	if itemCategory == "请选择" {
		itemCategory = ""
	}
	//常用选项
	var commonUseLabel []string = reqMsg.GetCommonUseLabel()
	var otherType []string = []string{}
	otherType = append(otherType, itemTypeName)
	otherType = append(otherType, itemCategory)

	if nil != commonUseLabel && len(commonUseLabel) > 0 {

		for _, value := range commonUseLabel {
			otherType = append(otherType, value)
		}
	}
	itemType.OtherType = otherType

	//组装图片url
	itemFiles := []*share_message.ItemFile{}
	if nil != reqMsg.ItemFiles && len(reqMsg.ItemFiles) > 0 {
		for i := 0; i < len(reqMsg.ItemFiles); i++ {
			var itemFile *share_message.ItemFile = &share_message.ItemFile{
				FileUrl:    easygo.NewString(reqMsg.GetItemFiles()[i].GetFileUrl()),
				FileType:   easygo.NewInt32(reqMsg.GetItemFiles()[i].GetFileType()),
				FileWidth:  easygo.NewString(reqMsg.GetItemFiles()[i].GetFileWidth()),
				FileHeight: easygo.NewString(reqMsg.GetItemFiles()[i].GetFileHeight())}
			itemFiles = append(itemFiles, itemFile)
		}
	}

	result.Name = easygo.NewString(reqMsg.GetName())
	//重新设置type
	result.Type = &itemType
	result.Price = easygo.NewInt32(reqMsg.GetPrice())
	//点卡的时候判断
	if reqMsg.GetItemType() != for_game.SHOP_POINT_CARD_CATEGORY {
		result.StockCount = easygo.NewInt32(reqMsg.GetStockCount())
	} else {
		//点卡的时候从导入库中取得
		result.StockCount = easygo.NewInt32(stockCnt)
	}

	result.UserName = easygo.NewString(reqMsg.UserName)
	result.Phone = easygo.NewString(reqMsg.GetPhone())
	//点卡的时候判断
	if reqMsg.GetItemType() != for_game.SHOP_POINT_CARD_CATEGORY {
		result.Address = easygo.NewString(reqMsg.GetAddress())
		result.DetailAddress = easygo.NewString(reqMsg.GetDetailAddress())
	}
	result.FakePayCnt = easygo.NewInt32(reqMsg.GetFakePaymentCount())
	result.FakePageViews = easygo.NewInt32(reqMsg.GetFakePageViews())
	result.FakeFixGoodCommRate = easygo.NewInt32(reqMsg.GetFakeGoodCommentRate())
	result.Title = easygo.NewString(reqMsg.GetTitle())
	result.State = easygo.NewInt32(reqMsg.GetState())
	result.ItemFiles = itemFiles

	//判断状态是下架的时候要设置下架时间
	if reqMsg.GetState() == for_game.SHOP_ITEM_SOLD_OUT {
		result.SoldOutTime = easygo.NewInt64(time.Now().Unix())
		result.CreateTime = easygo.NewInt64(time.Now().Unix())
	} else {
		result.CreateTime = easygo.NewInt64(time.Now().Unix())
	}

	if reqMsg.GetFakeSellItemCount() >= 0 {
		shopAuth := share_message.TableShopPlayer{
			PlayerId:            easygo.NewInt64(result.GetPlayerId()),
			UploadAuthFlag:      easygo.NewInt32(1),
			CreateTime:          easygo.NewInt64(time.Now().Unix()),
			FakePlayFinOrderCnt: easygo.NewInt32(reqMsg.GetFakeSellItemCount())}

		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_PLAYER)
		_, e := col.Upsert(
			bson.M{"_id": result.GetPlayerId()},
			bson.M{"$set": shopAuth})
		closeFun()

		if e != nil {
			logs.Error(e)
		}
	}

	EditShopItem(result)

	msg := fmt.Sprintf("后台修改商品: %d", result.GetItemId())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SHOP_MANAGE, msg)

	return easygo.EmptyMsg
}

//商品列表点击留言查看跳转到查询留言列表
func (self *cls4) RpcQueryShopComment(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryShopCommentRequest) easygo.IMessage {
	if reqMsg.ItemId == nil || reqMsg.GetItemId() == 0 {
		return easygo.NewFailMsg("商品id不能为空!")
	}

	list, count := QueryShopComment(reqMsg)
	msg := &brower_backstage.QueryShopCommentResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//留言修改点赞数页面的确定按钮
func (self *cls4) RpcEditShopComment(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.EditShopCommentRequest) easygo.IMessage {
	if reqMsg.CommentId == nil || reqMsg.GetCommentId() == 0 {
		return easygo.NewFailMsg("留言id错误！")
	}
	EditShopComment(reqMsg)

	msg := fmt.Sprintf("上架商品：修改点赞数-留言id%d", reqMsg.GetCommentId())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SHOP_MANAGE, msg)

	return easygo.EmptyMsg
}

//删除留言
func (self *cls4) RpcDeleteShopComment(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	id := reqMsg.GetId64()
	if id == 0 {
		return easygo.NewFailMsg("留言id不能为空")
	}
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ITEM_COMMENT)
	defer closeFun()

	var itemComment share_message.TableItemComment = share_message.TableItemComment{}

	errQuery := col.Find(bson.M{"_id": id}).One(&itemComment)
	if errQuery != nil && errQuery != mgo.ErrNotFound {
		return easygo.NewFailMsg("操作失败!")
	}

	if errQuery == mgo.ErrNotFound {
		return easygo.NewFailMsg("数据不存在!")
	}

	if itemComment.GetStatus() == for_game.SHOP_COMMENT_DELETE {

		return easygo.NewFailMsg("重复操作!")
	}

	DeleteShopComment(id)

	msg := fmt.Sprintf("删除留言：留言id%d", id)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SHOP_MANAGE, msg)

	return easygo.EmptyMsg
}

//后台导航跳转以及查询商城订单列表
func (self *cls4) RpcQueryShopOrder(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryShopOrderRequest) easygo.IMessage {
	list, count := QueryShopOrder(reqMsg)

	msg := &brower_backstage.QueryShopOrderResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//待发货的订单 发货页面取得快递公司下拉列表的内容
func (self *cls4) RpcGetExpressComDropDown(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {

	var dropDownExpressCom []*brower_backstage.ShopOrderExpressCom = []*brower_backstage.ShopOrderExpressCom{}
	var expressName string

	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("1000"),
		Name: easygo.NewString("请选择"),
	})

	expressName, _ = GetExpressNamePhone("ht")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("ht"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("ems")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("ems"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("sf")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("sf"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("sto")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("sto"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("tt")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("tt"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("yt")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("yt"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("yd")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("yd"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("youzheng")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("youzheng"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("zto")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("zto"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("bsky")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("bsky"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("db")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("db"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("emsg")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("emsg"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("jd")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("jd"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("suning")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("suning"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("yzgn")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("yzgn"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("zhongyou")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("zhongyou"),
		Name: easygo.NewString(expressName),
	})

	expressName, _ = GetExpressNamePhone("ztoky")
	dropDownExpressCom = append(dropDownExpressCom, &brower_backstage.ShopOrderExpressCom{
		Code: easygo.NewString("ztoky"),
		Name: easygo.NewString(expressName),
	})

	return &brower_backstage.GetExpressComDropDownResponse{
		DropDownExpressCom: dropDownExpressCom,
	}
}

//待付款取消商城订单确定按钮(待付款取消)
func (self *cls4) RpcCancelShopOrder(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.CancelShopOrderRequest) easygo.IMessage {

	oid := reqMsg.GetOrderId()
	if oid == 0 {
		return easygo.NewFailMsg("订单号错误")
	}

	shopOrder := QueryShopOrderById(oid)
	if shopOrder == nil {
		return easygo.NewFailMsg("该订单不存在")
	}

	if shopOrder.GetState() == for_game.SHOP_ORDER_BACKSTAGE_CANCLE {

		return easygo.NewFailMsg("重复操作!")
	}

	if shopOrder.GetState() != for_game.SHOP_ORDER_WAIT_PAY {

		return easygo.NewFailMsg("订单状态变化,请刷新列表!")
	}

	//以每个订单为单位取得锁,取不到说明订单状态在改变中
	lockKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_WAIT_PAY_MUTEX, oid)
	//取得分布式锁（失效时间设置10秒）
	errLock := easygo.RedisMgr.GetC().DoRedisLockNoRetry(lockKey, 10)
	defer easygo.RedisMgr.GetC().DoRedisUnlock(lockKey)

	//如果未取得锁
	if errLock != nil {
		s := fmt.Sprintf("RpcCancelShopOrder 单key取得订单redis分布式无重试锁失败,redis key is %v", lockKey)
		logs.Error(s)
		logs.Error(errLock)
		return easygo.NewFailMsg("操作失败,刷新重试")
	}

	//取得订单对应的商品的分布式锁，此锁需要重试(恢复库存的竞争)
	tempItemLockKey := for_game.MakeRedisKey(for_game.SHOP_ITEM_PAY_MUTEX, shopOrder.GetItems().GetItemId())
	//取得分布式锁开始1、取得订单对应的商品的分布式锁(阻塞重试）
	//1、取得订单对应的商品的分布式锁，此锁需要重试，直到重试次数结束提示退出
	errLock2 := easygo.RedisMgr.GetC().DoRedisLockWithRetry(tempItemLockKey, 10)
	defer easygo.RedisMgr.GetC().DoRedisUnlock(tempItemLockKey)

	//如果重试后还未取得锁就直接不做了
	if errLock2 != nil {
		s := fmt.Sprintf("RpcCancelShopOrder 单key取得商品redis分布式重试锁失败redis key is %v", tempItemLockKey)
		logs.Error(s)
		logs.Error(errLock2)
		return easygo.NewFailMsg("取消失败,刷新重试！")
	}

	//更新取消原因
	rst := UpdateOrderCancelReason(reqMsg, shopOrder.GetState())

	if rst != "" {
		return easygo.NewFailMsg(rst)
	}

	//通知大厅取消商城订单
	req := &server_server.ShopOrderRequest{
		OrderId: easygo.NewInt64(oid),
		Types:   easygo.NewInt32(NOTIFY_SHOP_OPERATE_1),
	}
	ChooseOneHall(0, "RpcBsOpShopOrder", req)

	msg := fmt.Sprintf("订单管理：后台取消待付款订单%d", oid)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SHOP_MANAGE, msg)

	return easygo.EmptyMsg
}

//待发货取消商城订单确定按钮(待发货取消)
func (self *cls4) RpcCancelShopOrderForWaitSend(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.CancelShopOrderRequest) easygo.IMessage {

	oid := reqMsg.GetOrderId()
	if oid == 0 {
		return easygo.NewFailMsg("订单号错误")
	}

	shopOrder := QueryShopOrderById(oid)
	if shopOrder == nil {
		return easygo.NewFailMsg("该订单不存在")
	}

	if shopOrder.GetState() == for_game.SHOP_ORDER_BACKSTAGE_CANCLE {

		return easygo.NewFailMsg("重复操作!")
	}

	if shopOrder.GetState() != for_game.SHOP_ORDER_WAIT_SEND {

		return easygo.NewFailMsg("订单状态变化,请刷新列表!")
	}

	//以每个订单为单位取得锁,取不到说明订单状态在改变中
	lockKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_SEND_MUTEX, oid)
	//取得分布式锁（失效时间设置10秒）
	errLock := easygo.RedisMgr.GetC().DoRedisLockNoRetry(lockKey, 10)
	defer easygo.RedisMgr.GetC().DoRedisUnlock(lockKey)

	//如果未取得锁
	if errLock != nil {
		s := fmt.Sprintf("RpcCancelShopOrderForWaitSend 单key取得订单redis分布式无重试锁失败,redis key is %v", lockKey)
		logs.Error(s)
		logs.Error(errLock)
		return easygo.NewFailMsg("订单可能已经发货，请刷新")
	}

	//取得订单对应的商品的分布式锁，此锁需要重试(恢复库存的竞争)
	tempItemLockKey := for_game.MakeRedisKey(for_game.SHOP_ITEM_PAY_MUTEX, shopOrder.GetItems().GetItemId())
	//取得分布式锁开始1、取得订单对应的商品的分布式锁(阻塞重试）
	//1、取得订单对应的商品的分布式锁，此锁需要重试，直到重试次数结束提示退出
	errLock2 := easygo.RedisMgr.GetC().DoRedisLockWithRetry(tempItemLockKey, 10)
	defer easygo.RedisMgr.GetC().DoRedisUnlock(tempItemLockKey)

	//如果重试后还未取得锁就直接不做了
	if errLock2 != nil {
		s := fmt.Sprintf("RpcCancelShopOrderForWaitSend 单key取得商品redis分布式重试锁失败redis key is %v", tempItemLockKey)
		logs.Error(s)
		logs.Error(errLock2)
		return easygo.NewFailMsg("取消失败,刷新重试！")
	}

	//更新取消原因
	rst := UpdateOrderCancelReason(reqMsg, shopOrder.GetState())

	if rst != "" {
		return easygo.NewFailMsg(rst)
	}

	//通知大厅取消商城订单
	req := &server_server.ShopOrderRequest{
		OrderId: easygo.NewInt64(oid),
		Types:   easygo.NewInt32(NOTIFY_SHOP_OPERATE_1),
	}
	ChooseOneHall(0, "RpcBsOpShopOrder", req)

	msg := fmt.Sprintf("订单管理：后台取消待发货订单%d", oid)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SHOP_MANAGE, msg)

	return easygo.EmptyMsg
}

//待收货取消商城订单确定按钮(待收货取消)
func (self *cls4) RpcCancelShopOrderForWaitReceive(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.CancelShopOrderRequest) easygo.IMessage {

	oid := reqMsg.GetOrderId()
	if oid == 0 {
		return easygo.NewFailMsg("订单号错误")
	}

	shopOrder := QueryShopOrderById(oid)
	if shopOrder == nil {
		return easygo.NewFailMsg("该订单不存在")
	}

	if shopOrder.GetState() == for_game.SHOP_ORDER_BACKSTAGE_CANCLE {

		return easygo.NewFailMsg("重复操作!")
	}

	if shopOrder.GetState() != for_game.SHOP_ORDER_WAIT_RECEIVE {

		return easygo.NewFailMsg("订单状态变化,请刷新列表!")
	}

	//以每个订单为单位取得锁,取不到说明订单在变化中直接返回
	lockKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_RECEIVE_MUTEX, reqMsg.GetOrderId())
	//取得分布式锁,跟主动收货互斥（失效时间设置10秒）
	errLock := easygo.RedisMgr.GetC().DoRedisLockNoRetry(lockKey, 10)
	defer easygo.RedisMgr.GetC().DoRedisUnlock(lockKey)

	//如果未取得锁
	if errLock != nil {
		s := fmt.Sprintf("RpcCancelShopOrderForWaitReceive 单key取得订单redis分布式无重试锁失败,redis key is %v", lockKey)
		logs.Error(s)
		logs.Error(errLock)
		return easygo.NewFailMsg("订单可能已经收货，请刷新")
	}

	//取得订单对应的商品的分布式锁，此锁需要重试(恢复库存的竞争)
	tempItemLockKey := for_game.MakeRedisKey(for_game.SHOP_ITEM_PAY_MUTEX, shopOrder.GetItems().GetItemId())
	//取得分布式锁开始1、取得订单对应的商品的分布式锁(阻塞重试）
	//1、取得订单对应的商品的分布式锁，此锁需要重试，直到重试次数结束提示退出
	errLock2 := easygo.RedisMgr.GetC().DoRedisLockWithRetry(tempItemLockKey, 10)
	defer easygo.RedisMgr.GetC().DoRedisUnlock(tempItemLockKey)

	//如果重试后还未取得锁就直接不做了
	if errLock2 != nil {
		s := fmt.Sprintf("RpcCancelShopOrderForWaitReceive 单key取得商品redis分布式重试锁失败redis key is %v", tempItemLockKey)
		logs.Error(s)
		logs.Error(errLock2)
		return easygo.NewFailMsg("取消失败,刷新重试！")
	}

	//更新取消原因
	rst := UpdateOrderCancelReason(reqMsg, shopOrder.GetState())

	if rst != "" {
		return easygo.NewFailMsg(rst)
	}

	//通知大厅取消商城订单
	req := &server_server.ShopOrderRequest{
		OrderId: easygo.NewInt64(oid),
		Types:   easygo.NewInt32(NOTIFY_SHOP_OPERATE_1),
	}
	ChooseOneHall(0, "RpcBsOpShopOrder", req)

	msg := fmt.Sprintf("订单管理：后台取消待收货订单%d", oid)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SHOP_MANAGE, msg)

	return easygo.EmptyMsg
}

//确认发货商城订单
func (self *cls4) RpcSendShopOrder(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.SendShopOrderRequest) easygo.IMessage {

	if reqMsg.ExpressCom == nil || reqMsg.GetExpressCom() == "1000" {
		return easygo.NewFailMsg("请选择快递公司")
	}

	kd := reqMsg.GetExpressCode()
	if kd == "" {
		return easygo.NewFailMsg("快递单号不能为空")
	}

	oid := reqMsg.GetOrderId()
	if oid == 0 {
		return easygo.NewFailMsg("订单号错误")
	}

	//以每个订单为单位取得锁,取不到说明订单状态在改变中
	lockKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_SEND_MUTEX, oid)
	//取得分布式锁（失效时间设置10秒）
	errLock := easygo.RedisMgr.GetC().DoRedisLockNoRetry(lockKey, 10)
	defer easygo.RedisMgr.GetC().DoRedisUnlock(lockKey)

	//如果未取得锁
	if errLock != nil {
		s := fmt.Sprintf("RpcSendShopOrder 单key取得redis分布式无重试锁失败,redis key is %v", lockKey)
		logs.Error(s)
		logs.Error(errLock)
		return easygo.NewFailMsg("订单可能已经发货，请刷新")
	}

	shopOrder := QueryShopOrderById(oid)
	if shopOrder == nil {
		return easygo.NewFailMsg("该订单不存在")
	}

	if shopOrder.GetState() == for_game.SHOP_ORDER_WAIT_RECEIVE {

		return easygo.NewFailMsg("重复操作!")
	}

	//0待付款 1超时 2取消 3待发货 4待收货 5已完成 6评价 7后台取消
	switch shopOrder.GetState() {
	case for_game.SHOP_ORDER_FINISH, for_game.SHOP_ORDER_EVALUTE:
		return easygo.NewFailMsg("该订单已完成,请刷新列表")
	case for_game.SHOP_ORDER_EXPIRE, for_game.SHOP_ORDER_CANCEL, for_game.SHOP_ORDER_BACKSTAGE_CANCLE:
		return easygo.NewFailMsg("该订单已取消,请刷新列表")
	case for_game.SHOP_ORDER_WAIT_PAY:
		return easygo.NewFailMsg("该订单还未付款,请刷新列表")
	case for_game.SHOP_ORDER_WAIT_RECEIVE:
		return easygo.NewFailMsg("该订单已经发货,请刷新列表")
	}

	rst := UpdateOrderExpressInfo(reqMsg)

	if rst != "" {
		return easygo.NewFailMsg(rst)
	}

	//通知大厅发货商城订单
	req := &server_server.ShopOrderRequest{
		OrderId: easygo.NewInt64(oid),
		Types:   easygo.NewInt32(NOTIFY_SHOP_OPERATE_2),
	}
	ChooseOneHall(0, "RpcBsOpShopOrder", req)

	msg := fmt.Sprintf("订单管理：后台发货订单%d", oid)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SHOP_MANAGE, msg)

	return easygo.EmptyMsg
}

//取得物流信息
func (self *cls4) RpcQueryShopOrderExpress(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	oid := reqMsg.GetId64()
	if oid == 0 {
		return easygo.NewFailMsg("订单号错误!")
	}

	shopOrder := QueryShopOrderById(oid)
	if shopOrder == nil {
		return easygo.NewFailMsg("该订单不存在!")
	}

	//通知大厅取得物流信息
	req := &server_server.ShopOrderRequest{
		OrderId: easygo.NewInt64(oid),
		Types:   easygo.NewInt32(NOTIFY_SHOP_OPERATE_5),
		UserId:  easygo.NewInt64(user.GetId()),
	}
	ChooseOneHall(0, "RpcBsOpShopOrder", req)

	msg := fmt.Sprintf("订单管理：取得物流信息%d", oid)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SHOP_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询商城收货地址
func (self *cls4) RpcQueryShopReceiveAddress(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	id := reqMsg.GetId64()
	if id < 1 {
		return easygo.NewFailMsg("用户ID错误 ")
	}

	list, count := QueryShopReceiveAddress(id)
	msg := &brower_backstage.QueryShopReceiveAddressResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//查询商城发货地址
func (self *cls4) RpcQueryShopDeliverAddress(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	id := reqMsg.GetId64()
	if id < 1 {
		return easygo.NewFailMsg("用户ID错误 ")
	}

	list, count := QueryShopDeliverAddress(id)
	msg := &brower_backstage.QueryShopDeliverAddressResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//商城点卡导入
func (self *cls4) RpcImportShopPointCard(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ImportShopPointCardRequest) easygo.IMessage {

	//导入信息没有值
	if nil == reqMsg.GetPointCardList() || len(reqMsg.GetPointCardList()) <= 0 {
		return easygo.NewFailMsg("点卡信息不存在")
	}

	//导入点卡信息正确性检测true信息无误  false存在信息错误
	infoCheckRst := ShopPointCardInfoCheck(reqMsg)
	if !infoCheckRst {
		return easygo.NewFailMsg("存在部分点卡信息不全，点卡导入失败")
	}

	//检测导入信息中是否存在重复卡号
	fileRepeatedMsg := GetFileRepeatedCardMsg(reqMsg)
	if nil != fileRepeatedMsg && len(fileRepeatedMsg) > 0 {
		return &brower_backstage.ImportShopPointCardResponse{
			Result: easygo.NewInt32(1),
			Msg:    fileRepeatedMsg}
	}

	tempList := reqMsg.GetPointCardList()
	for _, tempLi := range tempList {
		base := QueryPlayerbyAccount(tempLi.GetSellerAccount())
		if base == nil {
			s := fmt.Sprintf("卖家柠檬号[%s]不存在", tempLi.GetSellerAccount())
			return easygo.NewFailMsg(s)
		}
	}

	//导入的卡号是否在数据库中已经存在
	dbRepeatedMsg := GetDbRepeatedCardMsg(reqMsg)
	if nil != dbRepeatedMsg && len(dbRepeatedMsg) > 0 {
		return &brower_backstage.ImportShopPointCardResponse{
			Result: easygo.NewInt32(1),
			Msg:    dbRepeatedMsg}
	}

	//检测正确开始导入卡号
	//组装数据库对象
	shopCards := make([]*share_message.TableShopPointCard, 0)
	for _, value := range reqMsg.GetPointCardList() {
		cardId := easygo.NewInt64(for_game.NextId(for_game.TABLE_SHOP_POINT_CARD))

		//先取得一个随即key
		key := for_game.GenerateAesKey()
		//对key进行加密处理
		passKey := for_game.AesEncrypt(for_game.AES_KEY_SHOP_CARD, []byte(key))
		//用加密的key对卡密码再次加密
		cardPassword := for_game.AesEncrypt([]byte(compressStr(value.GetCardPassword())), []byte(passKey))
		shopCard := &share_message.TableShopPointCard{
			CardId:        cardId,
			CardName:      easygo.NewString(compressStr(value.GetCardName())),
			CardNo:        easygo.NewString(compressStr(value.GetCardNo())),
			CardPassword:  easygo.NewString(cardPassword),
			SellerAccount: easygo.NewString(compressStr(value.GetSellerAccount())),
			CardStatus:    easygo.NewInt32(for_game.SHOP_POINT_CARD_SALE),
			CreateTime:    easygo.NewInt64(time.Now().Unix()),
			Key:           easygo.NewString(key),
		}

		shopCards = append(shopCards, shopCard)
	}

	//批量插入数据库中
	var insLst []interface{}
	if nil != shopCards && len(shopCards) > 0 {
		for _, insValue := range shopCards {
			insLst = append(insLst, insValue)
		}
	}

	//批量插入
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_POINT_CARD)
	defer closeFun()

	if nil != insLst && len(insLst) > 0 {
		col.Insert(insLst...)
	}

	msg := fmt.Sprintf("点卡管理：导入点卡")
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SHOP_MANAGE, msg)

	return &brower_backstage.ImportShopPointCardResponse{
		Result: easygo.NewInt32(0)}
}

//后台导航跳转以及查询点卡列表
func (self *cls4) RpcQueryShopPointCard(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryShopPointCardRequest) easygo.IMessage {

	list, count := QueryShopPointCard(reqMsg)

	msg := &brower_backstage.QueryShopPointCardResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//商品发布页面 通过卖家帐号取得点卡名称的下拉内容
func (self *cls4) RpcGetShopPointCardDropDown(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.GetShopPointCardDropDownRequest) easygo.IMessage {

	//通过卖家帐号取得点卡名称列表(过滤重复的名称)
	if reqMsg.GetSellerAccount() == "" {
		return easygo.NewFailMsg("取得点卡名称下拉时,卖家帐号不能为空")
	}
	resCardNames := QueryShopPointCardDropDown(reqMsg)

	return &brower_backstage.GetShopPointCardDropDownResponse{
		DropDownShopPointCard: resCardNames,
	}
}
