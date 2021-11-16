// 许愿池逻辑层
package wish

import (
	"encoding/json"
	"errors"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/h5_wish"
	"game_server/pb/share_message"
	"time"

	"github.com/astaxie/beego/logs"
)

type wishData struct {
	WishBoxId     int64
	WishBoxItemId int64
}

//试玩一次
func TryOnceService(id int64) error {
	playerBase := for_game.GetRedisWishPlayer(id)
	if playerBase == nil {
		logs.Error("获取用户信息失败")
		return errors.New("获取用户信息失败")
	}
	playerBase.SetTryOne()
	return nil
}

//获取已收藏的盲盒
func GetCollectedWishBoxList(playerId int64, skip, limit, reqType int) ([]*h5_wish.CollectBox, int, int) {
	var wishCollectList []*share_message.PlayerWishCollection
	onSaleCount := 0
	SaleOutCount := 0

	wishCollectList, onSaleCount, SaleOutCount = for_game.GetCollectedBoxByPlayerId(playerId, skip, limit, reqType)

	if onSaleCount < 1 && SaleOutCount < 1 {
		return nil, 0, 0
	}
	var wishBoxIds []int64
	for _, v := range wishCollectList {
		wishBoxIds = append(wishBoxIds, v.GetWishBoxId())
	}

	wishBoxList, err := for_game.GetWishBoxListByIds(wishBoxIds)
	if err != nil || len(wishBoxList) < 1 {
		return nil, 0, 0
	}
	for k, r := range wishBoxList {
		wishBoxList[k].WishBoxId = easygo.NewInt64(r.GetXId())
	}

	return wishBoxList, onSaleCount, SaleOutCount
}

//获取所有已收藏的盲盒
func GetAllCollectedWishBoxList(playerId int64) ([]*h5_wish.CollectBox, error) {
	wishCollectList, err := for_game.GetAllCollectedBoxByPlayerId(playerId)
	if err != nil {
		return nil, errors.New("获取用户所有盲盒收藏数据失败")
	}

	var wishBoxIds []int64
	for _, v := range wishCollectList {
		wishBoxIds = append(wishBoxIds, v.GetWishBoxId())
	}

	wishBoxList, err1 := for_game.GetWishBoxListByIds(wishBoxIds)
	if err1 != nil {
		return nil, errors.New("查询盲盒列表失败")
	}
	preSaleBox, err2 := for_game.GetPreSaleBoxByItemIds(wishBoxIds)
	if err2 != nil {
		return nil, errors.New("查询预售商品列表失败")
	}
	preSaleBoxMap := make(map[int64]bool)
	for _, v := range preSaleBox {
		preSaleBoxMap[v.GetWishBoxId()] = true
	}

	for k, r := range wishBoxList {
		wishBoxList[k].WishBoxId = easygo.NewInt64(r.GetXId())
		wishBoxList[k].IsPreSale = easygo.NewBool(preSaleBoxMap[r.GetWishBoxId()])
	}

	return wishBoxList, nil
}

//获取待兑换、已兑换、已回收物品
func MyWishService(playerId int64, reqMsg *h5_wish.MyWishReq) ([]*h5_wish.Product, error) {
	var sortStr string
	reqType := reqMsg.GetType()
	page := easygo.If(reqMsg.GetPage() > 1, int(reqMsg.GetPage()-1), 0).(int)
	pageSize := int(reqMsg.GetPageSize())

	switch reqType {
	case for_game.WISH_TO_EXCHANGE, for_game.WISH_TO_CONTROL:
		sortStr = "-CreateTime"
		//err := for_game.EditPlayerItemIsRead(playerId)
		//if err != nil {
		//	return nil, errors.New("设置玩家新物品查阅状态失败")
		//}
	case for_game.WISH_EXCHANGED, for_game.WISH_RECYCLED, for_game.WISH_RECYCLE_CHECK:
		sortStr = "-UpdateTime"
	default:
		logs.Error("MyWishService参数有误: playerId: %v reqType: %v", playerId, reqType)
		return nil, errors.New("参数有误")
	}
	playerWishItem, err := for_game.GetPlayWishItemByPlayerId(playerId, int(reqType), page, pageSize, sortStr)
	if err != nil {
		return nil, errors.New("获取玩家物品所有列表失败")
	}
	productIds := make([]int64, 0)
	//封装返回数据
	respProduct := make([]*h5_wish.Product, 0)
	for _, v := range playerWishItem {
		productIds = append(productIds, v.GetWishItemId())
		retProduct := &h5_wish.Product{
			Id:               easygo.NewInt64(v.GetChallengeItemId()), //210407 盲盒下架后，该id对应的中间表信息会被删除
			ProductType:      easygo.NewInt32(v.GetWishItemStyle()),
			Icon:             easygo.NewString(v.GetWishItemIcon()),
			Name:             easygo.NewString(v.GetProductName()),
			BoxId:            easygo.NewInt64(v.GetWishBoxId()),
			Price:            easygo.NewInt64(v.GetWishItemPrice()), //价格单位均为分
			ProductId:        easygo.NewInt64(v.GetWishItemId()),
			BoxName:          easygo.NewString(v.GetBoxName()),
			PlayerWishItemId: easygo.NewInt64(v.GetId()),
			ExpireTime:       easygo.NewInt64(v.GetExpireTime()),
			Match:            easygo.NewInt32(v.GetBoxMatch()),
			CtrlStatus:       easygo.NewInt32(v.GetStatus()),
			Diamond:          easygo.NewInt64(v.GetWishItemDiamond()),
			GiveType:         easygo.NewInt64(v.GetGiveType()),
		}
		if reqType == for_game.WISH_RECYCLED {
			retProduct.RecyclePrice = easygo.NewInt64(v.GetRecyclePrice())
			retProduct.RecycleType = easygo.NewInt32(v.GetRecycleType())
		}
		respProduct = append(respProduct, retProduct)
	}
	products, err1 := for_game.GetWishItemByIdsFromDB(productIds)
	if err1 != nil {
		logs.Error("获取物品最新状态失败：%v", err1)
	}
	productMap := make(map[int64]*share_message.WishItem) //[物品id]物品信息
	for _, v := range products {
		productMap[v.GetId()] = v
	}

	for _, v := range respProduct {
		if product, ok := productMap[v.GetProductId()]; ok {
			v.IsPreSale = easygo.NewBool(product.GetIsPreSale())
			v.PreHaveTime = easygo.NewInt64(product.GetPreHaveTime())
		}
	}

	return respProduct, nil
}

//已弃用。 获取所有待兑换、已兑换、已回收物品
func MyAllWishService(playerId int64, reqType int32) ([]*h5_wish.Product, error) {
	var sortStr string
	switch reqType {
	case for_game.WISH_TO_EXCHANGE:
		sortStr = "-CreateTime"
		err := for_game.EditPlayerItemIsRead(playerId)
		if err != nil {
			return nil, errors.New("设置玩家新物品查阅状态失败")
		}
	case for_game.WISH_EXCHANGED, for_game.WISH_RECYCLED:
		sortStr = "-UpdateTime"
	default:
		logs.Error("MyAllWishService参数有误: playerId: %v reqType: %v", playerId, reqType)
		return nil, errors.New("参数有误")
	}
	playerWishItem, err := for_game.GetAllPlayWishItemByPlayerId(playerId, int(reqType), sortStr)
	if err != nil {
		return nil, errors.New("获取玩家物品所有列表失败")
	}
	productIds := make([]int64, 0)
	//封装返回数据
	respProduct := make([]*h5_wish.Product, 0)
	for _, v := range playerWishItem {
		retProduct := &h5_wish.Product{
			Id:               easygo.NewInt64(v.GetChallengeItemId()), //210407 盲盒下架后，该id对应的中间表信息会被删除
			ProductType:      easygo.NewInt32(v.GetWishItemStyle()),
			Icon:             easygo.NewString(v.GetWishItemIcon()),
			Name:             easygo.NewString(v.GetProductName()),
			BoxId:            easygo.NewInt64(v.GetWishBoxId()),
			Price:            easygo.NewInt64(v.GetWishItemPrice()),
			ProductId:        easygo.NewInt64(v.GetWishItemId()),
			BoxName:          easygo.NewString(v.GetBoxName()),
			PlayerWishItemId: easygo.NewInt64(v.GetId()),
			ExpireTime:       easygo.NewInt64(v.GetExpireTime()),
			Match:            easygo.NewInt32(v.GetBoxMatch()),
			CtrlStatus:       easygo.NewInt32(v.GetStatus()),
			Diamond:          easygo.NewInt64(v.GetWishItemDiamond()),
		}
		if reqType == for_game.WISH_RECYCLED {
			retProduct.RecyclePrice = easygo.NewInt64(v.GetRecyclePrice())
		}
		respProduct = append(respProduct, retProduct)
	}
	products, err1 := for_game.GetWishItemByIdsFromDB(productIds)
	if err1 != nil {
		logs.Error("获取物品最新状态失败：%v", err1)
	}
	productMap := make(map[int64]*share_message.WishItem) //[物品id]物品信息
	for _, v := range products {
		productMap[v.GetId()] = v
	}

	for _, v := range respProduct {
		if product, ok := productMap[v.GetProductId()]; ok {
			v.IsPreSale = easygo.NewBool(product.GetIsPreSale())
			v.PreHaveTime = easygo.NewInt64(product.GetPreHaveTime())
		}
	}

	return respProduct, nil
}

//分页获取已许愿的盲盒
func GetWishDataList(playerId int64, skip, limit, reqType int) ([]*h5_wish.CollectBox, int, int) {
	wishList, onSaleTotal, saleOutTotal := for_game.GetPlayerWishDataByPlayerId(playerId, skip, limit, reqType)
	if onSaleTotal == 0 && saleOutTotal == 0 {
		return nil, 0, 0
	}

	var wishBoxIds []int64
	var wishDataMap []wishData
	for _, v := range wishList {
		wishBoxIds = append(wishBoxIds, v.GetWishBoxId())
		wishDataMap = append(wishDataMap, wishData{
			WishBoxId:     v.GetWishBoxId(),
			WishBoxItemId: v.GetWishBoxItemId(),
		})
	}

	wishBoxList, _ := for_game.GetWishBoxListByIds(wishBoxIds)
	for _, r := range wishDataMap {
		for k, rr := range wishBoxList {
			if r.WishBoxId == rr.GetXId() {
				wishBoxList[k].WishBoxItemId = easygo.NewInt64(r.WishBoxItemId)
				wishBoxList[k].WishBoxId = easygo.NewInt64(r.WishBoxId)
				break
			}
		}
	}
	return wishBoxList, onSaleTotal, saleOutTotal
}

//获取所有已许愿的盲盒
func GetAllWishDataList(playerId int64) ([]*h5_wish.CollectBox, error) {
	wishList, err := for_game.GetPlayerAllWishDataByPlayerId(playerId)
	if err != nil {
		return nil, errors.New("查询用户所有的已许愿列表失败")
	}
	wishBoxList := make([]*h5_wish.CollectBox, 0)
	if len(wishList) == 0 {
		return wishBoxList, nil
	}

	var wishBoxIds []int64
	var wishDataMap []wishData
	for _, v := range wishList {
		wishBoxIds = append(wishBoxIds, v.GetWishBoxId())
		wishDataMap = append(wishDataMap, wishData{
			WishBoxId:     v.GetWishBoxId(),
			WishBoxItemId: v.GetWishBoxItemId(),
		})
	}
	if len(wishBoxIds) == 0 {
		return wishBoxList, nil
	}
	wishBoxList, err = for_game.GetWishBoxListByIds(wishBoxIds)
	if err != nil {
		return nil, errors.New("查询收藏盲盒列表数据失败")
	}
	for _, r := range wishDataMap {
		for k, rr := range wishBoxList {
			if r.WishBoxId == rr.GetXId() {
				wishBoxList[k].WishBoxItemId = easygo.NewInt64(r.WishBoxItemId)
				wishBoxList[k].WishBoxId = easygo.NewInt64(r.WishBoxId)
				break
			}
		}
	}
	return wishBoxList, nil
}

//收货地址列表
func GetAddressListByUid(uid int64, page, pageSize int) []*h5_wish.WishAddress {
	player := for_game.GetRedisWishPlayer(uid)
	if player == nil {
		logs.Error("获取用户信息失败")
		return nil
	}
	addressList := player.GetAddressList()
	maxLen := 0
	addressLen := len(addressList)
	getLen := page * pageSize
	if addressLen > getLen {
		maxLen = getLen
	} else {
		maxLen = addressLen
	}
	//someAddress := make([]*share_message.WishAddress, 0)
	//for i := (page - 1) * pageSize; i < maxLen; i++ {
	//	copy(someAddress, addressList)
	//}
	someAddress := addressList[(page-1)*pageSize : maxLen]

	jsonBytes, _ := json.Marshal(someAddress)
	var res []*h5_wish.WishAddress
	_ = json.Unmarshal(jsonBytes, &res)
	return res
}

//添加收货地址
func AddAddress(playerId int64, address *h5_wish.WishAddress) bool {
	player := for_game.GetRedisWishPlayer(playerId)
	if player == nil {
		logs.Error("获取用户信息有误")
		return false
	}
	playerAddress := player.GetAddressList()
	address.AddressId = easygo.NewInt64(for_game.GetMillSecond())
	//若新地址为默认地址，需要将之前地址设置非默认
	if address.GetIfDefault() == true {
		for _, v := range playerAddress {
			if v.GetIfDefault() == true {
				v.IfDefault = easygo.NewBool(false)
			}
		}
	}
	newAddress := &share_message.WishAddress{}
	for_game.StructToOtherStruct(address, newAddress)
	playerAddress = append(playerAddress, newAddress)
	player.SetAddressList(playerAddress)
	return true
	/*address.AddressId = easygo.NewInt64(time.Now().Unix())
	if address.GetIfDefault() == true {
		err := for_game.ClearUserDefaultAddress(playerId)
		if err != nil {
			return err
		}
	}
	err := for_game.InsertAddressByUid(playerId, address)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		logs.Error("添加收货地址: playerId: %v, err: %v", playerId, err)
		return err
	}
	return nil*/
}

//编辑收货地址
func EditAddress(playerId int64, address *h5_wish.WishAddress) bool {
	player := for_game.GetRedisWishPlayer(playerId)
	if player == nil {
		logs.Error("获取用户信息有误")
		return false
	}
	playerAddress := player.GetAddressList()
	for _, v := range playerAddress {
		if v.GetAddressId() == address.GetAddressId() {
			if address.GetIfDefault() == true { //取消之前的默认地址
				for _, v := range playerAddress {
					v.IfDefault = easygo.NewBool(false)
				}
			}
			v.IfDefault = address.IfDefault
			v.Name = address.Name
			v.Phone = address.Phone
			v.Detail = address.Detail
			v.Province = address.Province
			v.City = address.City
			v.Area = address.Area
			break
		}
	}
	player.SetAddressList(playerAddress)
	return true
	//err := for_game.EditAddressByUid(playerId, address)
}

//移除收货地址
func RemoveAddress(playerId, addressId int64) bool {
	player := for_game.GetRedisWishPlayer(playerId)
	if player == nil {
		logs.Error("获取用户信息有误")
		return false
	}
	playerAddress := player.GetAddressList()
	for i := 0; i < len(playerAddress); i++ {
		if playerAddress[i].GetAddressId() == addressId {
			playerAddress = append(playerAddress[:i], playerAddress[i+1:]...)
		}
	}
	player.SetAddressList(playerAddress)
	//err := for_game.RemoveAddressByUid(playerId, addressId)
	return true
}

//获取盲盒商品的所有信息
func GetBoxItemAllDateByIds(wishBoxItemIds []int64) (map[int64]*share_message.WishBox,
	map[int64]*share_message.WishBoxItem, map[int64]*share_message.WishItem, error) {
	wishBoxItemList, err := for_game.GetWishBoxItemByIdsFromDB(wishBoxItemIds)
	if err != nil {
		logs.Error("查询盲盒详细数据列表数据失败: wishBoxItemIds: %v, errors: %v", wishBoxItemIds, err)
		return nil, nil, nil, errors.New("查询盲盒详细数据列表数据失败")
	}

	wishBoxIds := make([]int64, 0)
	productIds := make([]int64, 0)
	boxItemMap := make(map[int64]*share_message.WishBoxItem)
	for _, v := range wishBoxItemList {
		boxItemMap[v.GetId()] = v
		wishBoxIds = append(wishBoxIds, v.GetWishBoxId())
		productIds = append(productIds, v.GetWishItemId())
	}
	boxMap := make(map[int64]*share_message.WishBox)
	wishBoxList, err := for_game.GetWishBoxByIds(wishBoxIds)
	if err != nil {
		return nil, nil, nil, errors.New("查询盲盒列表数据失败")
	}
	for _, v := range wishBoxList {
		boxMap[v.GetId()] = v
	}
	productMap := make(map[int64]*share_message.WishItem)
	productList, err := for_game.GetWishItemByIdsFromDB(productIds)
	if err != nil {
		logs.Error("查询商品列表数据失败: productIds: %v, errors: %v", productIds, err)
		return nil, nil, nil, errors.New("查询商品列表数据失败")
	}
	for _, v := range productList {
		productMap[v.GetId()] = v
	}
	return boxMap, boxItemMap, productMap, nil
}

// 收藏或者取消收藏
func CollectionBoxService(pid int64, reqMsg *h5_wish.CollectionBoxReq) error {
	switch reqMsg.GetOpType() {
	case 1: // 1-收藏
		// 判断是否已收藏
		if pwc := for_game.GetPlayerWishCollection(pid, reqMsg.GetIdList()[0]); pwc.GetId() != 0 {
			return errors.New("已经收藏过了")
		}
		// 获取盲盒数据
		data := &share_message.PlayerWishCollection{
			Id:         easygo.NewInt64(for_game.NextId(for_game.TABLE_PLAYER_WISH_COLLECTION)),
			PlayerId:   easygo.NewInt64(pid),
			WishBoxId:  easygo.NewInt64(reqMsg.GetIdList()[0]),
			CreateTime: easygo.NewInt64(time.Now().Unix()),
		}
		return for_game.InsertPlayerWishCollection(data)
	case 2: // 2-删除
		return for_game.RemovePlayerWishCollectionByIds(pid, reqMsg.GetIdList())
	default:
		logs.Error("收藏或者取消收藏 opType 有误,opType: %d", reqMsg.GetOpType())
	}
	return nil
}

// 盲盒商品详情
func BoxProductDetailService(boxId int64) ([]*h5_wish.ProductDetail, error) {
	// 中间表查询物品信息
	boxItem, err := for_game.QueryWishBoxItemByWishBoxIds([]int64{boxId})
	if err != nil {
		return nil, errors.New("获取盲盒信息失败")
	}
	productIds := make([]int64, 0)
	boxItemMap := make(map[int64]*share_message.WishBoxItem, 0) //商品ID,中间表信息
	for _, v := range boxItem {
		boxItemMap[v.GetWishItemId()] = v
		productIds = append(productIds, v.GetWishItemId())
	}

	wishItemArr, err1 := for_game.GetWishItemByIdsFromDB(productIds)
	if err1 != nil {
		return nil, errors.New("获取盲盒商品信息失败")
	}
	result := make([]*h5_wish.ProductDetail, 0)

	for _, v := range wishItemArr {
		boxItem := boxItemMap[v.GetId()]
		result = append(result, &h5_wish.ProductDetail{
			ProductId:   easygo.NewInt64(v.GetId()),
			ProductName: easygo.NewString(v.GetName()),
			Image:       easygo.NewString(v.GetIcon()),
			ProductType: easygo.NewInt32(boxItem.GetStyle()),
			Desc:        easygo.NewString(v.GetDesc()),
			Long:        easygo.NewInt64(v.GetLength()),
			Width:       easygo.NewInt64(v.GetWide()),
			High:        easygo.NewInt64(v.GetHigh()),
			Price:       easygo.NewInt64(boxItem.GetPrice()),
		})
	}
	return result, nil
}

// 回收盲盒商品
func RecycleGoods(player *for_game.RedisWishPlayerObj, playerWishItems []*share_message.PlayerWishItem, recycleType int,
	recycleRatio int32, reqMsg *h5_wish.WishBoxReq) (int, error) {
	result := for_game.RECYCLE_STATUS_ERROR
	playerId := player.GetId()
	var recyclePriceTotal int64
	var productPriceTotal int64
	recycleMap := make(map[int64]int64) //[玩家物品表id]回收价格

	for _, r := range playerWishItems {
		price := r.GetWishItemPrice()
		//金额单位为分，四舍五入需精准到元
		recyclePrice := int64(easygo.Decimal(float64(price)*float64(recycleRatio)*0.0001, 0) * 100)
		productPriceTotal += price
		recyclePriceTotal += recyclePrice
		recycleMap[r.GetId()] = recyclePrice
	}

	bankCardId := reqMsg.GetBankCardId()
	recycleStatus := RECYCLE_CHECKING
	recycleMsg := &h5_wish.RecycleToHall{
		PlayerId:   easygo.NewInt64(player.GetPlayerId()),
		BankCardId: easygo.NewString(bankCardId),
		Price:      easygo.NewInt64(recyclePriceTotal),
	}

	//直接回收到用户银行卡
	orderMsg, serverErr := SendMsgToIdelServer(for_game.SERVER_TYPE_HALL, "RpcRecycleToBandCard", recycleMsg, player.GetPlayerId())
	if serverErr != nil {
		result = for_game.RECYCLE_STATUS_ERROR
		logs.Error("用户(%v)直接回收到用户银行卡失败：%v", playerId, serverErr.GetReason())
		return result, errors.New(serverErr.GetReason())
	}
	var paymentOrderId string
	if res1, ok := orderMsg.(*h5_wish.OrderMsgResp); ok {
		paymentOrderId = res1.GetOrderId()
		result = int(res1.GetStatus())
	}
	//用户回收表状态
	if result == for_game.RECYCLE_STATUS_CHECK {
		//进入人工审核，修改回收方式、回收状态
		recycleType = for_game.RECYCLE_BY_CHECK
		recycleStatus = RECYCLE_CHECKING
	} else if result == for_game.RECYCLE_STATUS_SUCCESS {
		//回收成功，修改回收状态
		recycleStatus = RECYCLE_RECYCLED
	} else {
		logs.Error("用户(%v)数据异常，回收到用户银行卡失败", playerId)
		return result, errors.New(RECYCLE_ERR_STR)
	}

	recycleItemList := make([]*share_message.WishRecycleItem, 0)
	//填写回收订单
	for _, v := range playerWishItems {
		recycleItem := &share_message.WishRecycleItem{
			BoxDrawPrice:   easygo.NewInt64(v.GetDareDiamond()),
			WinTime:        easygo.NewInt64(v.GetCreateTime()),
			WishBoxId:      easygo.NewInt64(v.GetWishBoxId()),
			ProductId:      easygo.NewInt64(v.GetWishItemId()),
			Style:          easygo.NewInt32(v.GetWishItemStyle()),
			ProductName:    easygo.NewString(v.GetProductName()),
			ProductIcon:    easygo.NewString(v.GetWishItemIcon()),
			PlayerItemId:   easygo.NewInt64(v.GetId()),
			ProductPrice:   easygo.NewInt64(v.GetWishItemPrice()),
			WishBoxItemId:  easygo.NewInt64(v.GetChallengeItemId()),
			ProductDiamond: easygo.NewInt64(v.GetWishItemDiamond()),
			RecyclePrice:   easygo.NewInt64(recycleMap[v.GetId()]),
			RecycleDiamond: easygo.NewInt64(easygo.Decimal(float64(v.GetWishItemDiamond())*float64(recycleRatio)*0.01, 0)),
			GiveType:       easygo.NewInt64(v.GetGiveType()),
		}
		recycleItemList = append(recycleItemList, recycleItem)
	}

	data := &share_message.WishRecycleOrder{
		UserId:            easygo.NewInt64(player.GetPlayerId()),
		PlayerId:          easygo.NewInt64(playerId),
		Type:              easygo.NewInt32(recycleType),
		Status:            easygo.NewInt32(recycleStatus),
		RecycleItemList:   recycleItemList,
		ProductPriceTotal: easygo.NewInt64(productPriceTotal),
		RecycleNote:       easygo.NewInt32(reqMsg.GetRecycleNote()),
		RecyclePriceTotal: easygo.NewInt64(recyclePriceTotal),
		BankCardId:        easygo.NewString(bankCardId),
		Channel:           easygo.NewInt32(1),
		PaymentOrderId:    easygo.NewString(paymentOrderId),
	}

	//用户回收表状态
	playerItemStatus := for_game.WISH_RECYCLE_CHECK
	if result == for_game.RECYCLE_STATUS_SUCCESS {
		data.RecycleTime = easygo.NewInt64(easygo.NowTimestamp())
		playerItemStatus = for_game.WISH_RECYCLED
	}
	err := for_game.InsertPlayerRecycleLog(data)
	if err != nil {
		logs.Error("插入玩家(IM_ID:%v)回收日志失败，请查看提现订单Id：%v, 发生的错误: %v", playerId, paymentOrderId, err.Error())
		return for_game.RECYCLE_STATUS_ERROR, errors.New(RECYCLE_ERR_STR)
	}
	//设置玩家物品表数据
	for_game.RecycleToPlayerItem(recycleMap, playerItemStatus, for_game.RECYCLE_TO_CARD)

	return result, nil
}

//回收商品以钻石结算
func RecycleGoodsToDiamond(player *for_game.RedisWishPlayerObj, playerWishItems []*share_message.PlayerWishItem, recycleType int,
	recycleRatio int32, reqMsg *h5_wish.WishBoxReq) (int, error) {
	result := for_game.RECYCLE_STATUS_ERROR //返回操作结果，只有返回的error不为nil才有效
	playerId := player.GetId()
	var productPriceTotal int64
	var totalDiamond int64
	recycleMap := make(map[int64]int64) //[玩家物品表ID]回收钻石
	for _, r := range playerWishItems {
		price := r.GetWishItemPrice()
		recycleDiamond := int64(easygo.Decimal(float64(r.GetWishItemDiamond())*float64(recycleRatio)*0.01, 0))
		productPriceTotal += price
		totalDiamond += recycleDiamond
		recycleMap[r.GetId()] = recycleDiamond
	}

	//风控值检测
	config := for_game.GetConfigWishPayment()
	if config == nil {
		logs.Error("RecycleGoods获取config为空")
		return result, errors.New(RECYCLE_ERR_STR)
	}
	recycleStatus := RECYCLE_CHECKING
	extendLog := &share_message.GoldExtendLog{
		PlayerId: easygo.NewInt64(player.GetPlayerId()),
	}
	if !config.GetStatus() {
		return result, errors.New("回收失败，回收功能暂未开启")
	}
	//每日回收次数限制
	period := for_game.GetPlayerPeriod(player.GetPlayerId())
	curTimes := period.DayPeriod.FetchInt32(for_game.WISH_OUT_DIAMOND_TIME)
	if curTimes >= config.GetDayDiamondTopCount() {
		return result, errors.New("回收失败，今天回收次数已达" + easygo.AnytoA(config.GetDayDiamondTopCount()) + "次")
	}
	//每日回收总额限制
	curSum := period.DayPeriod.FetchInt64(for_game.WISH_OUT_DIAMOND_SUM)
	if totalDiamond+curSum > config.GetDayDiamondTop() {
		return result, errors.New("回收失败，今天回收钻石总额超" + easygo.AnytoA(float64(config.GetDayDiamondTop())))
	}
	var remarks string
	//低于安全阀值则自动回收，否则设置人工审核
	if totalDiamond < config.GetOrderThreshold() {
		result = for_game.RECYCLE_STATUS_SUCCESS
		recycleStatus = RECYCLE_RECYCLED
		remarks = "交易完成"
		err, _ := player.AddDiamond(totalDiamond, "回收钻石", for_game.DIAMOND_TYPE_WISH_BACK, extendLog)
		if err != nil {
			logs.Error("用户(%v)回收时增加钻石失败：%v", playerId, err)
			return for_game.RECYCLE_STATUS_ERROR, errors.New(RECYCLE_ERR_STR)
		}
	} else {
		//进入人工审核
		recycleType = for_game.RECYCLE_BY_CHECK
		result = for_game.RECYCLE_STATUS_CHECK
		remarks = "回收额度超过风控制，进行人工审核"
	}

	recycleItemList := make([]*share_message.WishRecycleItem, 0)
	//填写回收订单
	for _, v := range playerWishItems {
		recycleItem := &share_message.WishRecycleItem{
			BoxDrawPrice:   easygo.NewInt64(v.GetDareDiamond()),
			WinTime:        easygo.NewInt64(v.GetCreateTime()),
			WishBoxId:      easygo.NewInt64(v.GetWishBoxId()),
			ProductId:      easygo.NewInt64(v.GetWishItemId()),
			Style:          easygo.NewInt32(v.GetWishItemStyle()),
			ProductName:    easygo.NewString(v.GetProductName()),
			ProductIcon:    easygo.NewString(v.GetWishItemIcon()),
			PlayerItemId:   easygo.NewInt64(v.GetId()),
			ProductPrice:   easygo.NewInt64(v.GetWishItemPrice()),
			WishBoxItemId:  easygo.NewInt64(v.GetChallengeItemId()),
			ProductDiamond: easygo.NewInt64(v.GetWishItemDiamond()),
			RecycleDiamond: easygo.NewInt64(recycleMap[v.GetId()]),
			RecyclePrice:   easygo.NewInt64(easygo.NewInt64(easygo.Decimal(float64(v.GetWishItemPrice())*float64(recycleRatio)*0.01, 0))),
			GiveType:       easygo.NewInt64(v.GetGiveType()),
		}
		recycleItemList = append(recycleItemList, recycleItem)
	}

	data := &share_message.WishRecycleOrder{
		UserId:            easygo.NewInt64(player.GetPlayerId()),
		PlayerId:          easygo.NewInt64(playerId),
		Type:              easygo.NewInt32(recycleType),
		Status:            easygo.NewInt32(recycleStatus),
		RecycleItemList:   recycleItemList,
		ProductPriceTotal: easygo.NewInt64(productPriceTotal),
		RecycleNote:       easygo.NewInt32(reqMsg.GetRecycleNote()),
		RecycleDiamond:    easygo.NewInt64(totalDiamond),
		Channel:           easygo.NewInt32(2),
		Remarks:           easygo.NewString(remarks),
	}

	//用户回收表状态
	playerItemStatus := for_game.WISH_RECYCLE_CHECK
	if result == for_game.RECYCLE_STATUS_SUCCESS {
		data.RecycleTime = easygo.NewInt64(easygo.NowTimestamp())
		playerItemStatus = for_game.WISH_RECYCLED
	}
	err := for_game.InsertPlayerRecycleLog(data)
	if err != nil {
		if result == for_game.RECYCLE_STATUS_SUCCESS {
			err1, _ := player.AddDiamond(0-totalDiamond, "回收失败返回", for_game.DIAMOND_TYPE_WISH_BACK_FAILD, extendLog)
			if err1 != nil {
				logs.Error("插入玩家(IM_ID %v)回收日志失败后返回钻石(%v)不成功：%v", playerId, 0-totalDiamond, err1)
				return for_game.RECYCLE_STATUS_ERROR, errors.New(RECYCLE_ERR_STR)
			}
		}
		logs.Error("插入玩家(IM_ID %v)回收钻石(%v)日志失败：%v", playerId, 0-totalDiamond, err)
		return for_game.RECYCLE_STATUS_ERROR, errors.New(RECYCLE_ERR_STR)
	}
	// 体现预警通知
	easygo.Spawn(for_game.CheckWishWarningSMS, int64(0), totalDiamond)
	//设置玩家物品表数据
	for_game.RecycleToPlayerItem(recycleMap, playerItemStatus, for_game.RECYCLE_TO_DIAMOND)
	period.DayPeriod.AddInteger(for_game.WISH_OUT_DIAMOND_TIME, 1)
	period.DayPeriod.AddInteger(for_game.WISH_OUT_DIAMOND_SUM, totalDiamond)
	return result, nil
}

//兑换物品
func ExchangeGoods(player *for_game.RedisWishPlayerObj, playerWishItems []*share_message.PlayerWishItem,
	freeNum int32, address *share_message.WishAddress) int {
	playerId := player.GetId()
	//校验有无预售商品，若有则去除
	productIds := make([]int64, 0)
	for _, r := range playerWishItems {
		productIds = append(productIds, r.GetWishItemId())
	}
	preSaleProduct := for_game.GetPreSaleProductByIds(productIds)
	if len(preSaleProduct) > 0 {
		preSaleIds := make([]int64, 0)
		for _, v := range preSaleProduct {
			for i := 0; i < len(playerWishItems); i++ {
				if playerWishItems[i].GetWishItemId() == v.GetId() {
					preSaleIds = append(preSaleIds, playerWishItems[i].GetId())
					playerWishItems = append(playerWishItems[:i], playerWishItems[i+1:]...)
					i--
				}
			}
		}
		//还原预售商品为待兑换状态
		for_game.ExchangeToPlayerItem(preSaleIds, for_game.WISH_TO_EXCHANGE)
	}
	if len(playerWishItems) < 1 {
		logs.Error("兑换物品为预售状态")
		return for_game.EXCHANGE_STATUS_ISPRESALE
	}

	var postage int32
	extendLog := &share_message.GoldExtendLog{
		PlayerId: easygo.NewInt64(playerId),
	}
	if len(playerWishItems) < int(freeNum) {
		postage = for_game.CheckAreaPostage(address.GetProvince())
		err, _ := player.AddDiamond(0-int64(postage), "运费", for_game.DIAMOND_TYPE_POSTAGE_OUT, extendLog)
		if err != nil {
			logs.Error("用户(%v)扣除邮费失败：%v", playerId, err)
			return for_game.EXCHANGE_STATUS_CHECK
		}
	}

	addressDetail := address.GetProvince() + address.GetCity() + address.GetArea() + address.GetDetail()
	orderId := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_OUT, for_game.DIAMOND_TYPE_POSTAGE_OUT)
	exchangeLogs := make([]interface{}, 0)
	//订单生成后修改玩家物品状态
	editStatusIds := make([]int64, 0)
	for _, v := range playerWishItems {
		editStatusIds = append(editStatusIds, v.GetId())
		exchangeLogs = append(exchangeLogs, &share_message.PlayerExchangeLog{
			Id:             easygo.NewInt64(for_game.NextId(for_game.TABLE_PLAYER_EXCHANGE_LOG)),
			PlayerId:       easygo.NewInt64(playerId),
			UserId:         easygo.NewInt64(player.GetPlayerId()),
			Status:         easygo.NewInt32(for_game.WISH_TO_EXCHANGE), //代发货状态
			Receiver:       easygo.NewString(address.GetName()),
			Phone:          easygo.NewString(address.GetPhone()),
			Address:        easygo.NewString(addressDetail),
			CreateTime:     easygo.NewInt64(easygo.NowTimestamp()),
			OrderId:        easygo.NewString(orderId),
			BoxDrawPrice:   easygo.NewInt64(v.GetDareDiamond()),
			WishBoxItem:    easygo.NewInt64(v.GetChallengeItemId()),
			WishBoxId:      easygo.NewInt64(v.GetWishBoxId()),
			ProductId:      easygo.NewInt64(v.GetWishItemId()),
			ProductName:    easygo.NewString(v.GetProductName()),
			ProductIcon:    easygo.NewString(v.GetWishItemIcon()),
			ProductPrice:   easygo.NewInt64(v.GetWishItemPrice()),
			ProductDiamond: easygo.NewInt64(v.GetWishItemDiamond()),
			PlayerItemId:   easygo.NewInt64(v.GetId()),
			Postage:        easygo.NewInt64(postage),
			GiveType:       easygo.NewInt64(v.GetGiveType()),
		})
	}
	//插入兑换日志
	err := for_game.InsertPlayerExchangeLog(exchangeLogs)
	if err != nil {
		if postage > 0 {
			err, _ := player.AddDiamond(int64(postage), "运费返回", for_game.DIAMOND_TYPE_POSTAGE_FAILD, extendLog)
			if err != nil {
				return for_game.EXCHANGE_STATUS_ERROR
			}
		}
		logs.Error("用户(%v)插入兑换日志失败：%v", playerId, err)
		return for_game.EXCHANGE_STATUS_ERROR
	}
	return for_game.EXCHANGE_STATUS_SUCCESS
}

//访问许愿池埋点
func AccessWishBPoint(pid int64, playerAccess *share_message.WishPlayerAccessLog) {
	cur0Time := easygo.GetToday0ClockTimestamp()
	if playerAccess == nil || playerAccess.GetWishTime() == 0 {
		//更新上次访问许愿池日期，并修改新用户登录次数，设置记录
		for_game.SetWishPlayerAccessLogToRedis(pid, cur0Time, "WishTime")
		for_game.SetWishPlayerAccessLogToRedis(pid, 1, "RetainedDay")
		for_game.SetRedisWishLogReportFieldVal(cur0Time, 1, "NewPlayer")
	} else {
		leadTime := cur0Time - playerAccess.GetWishTime()
		if leadTime > 0 {
			var filed string
			//更新上次访问许愿池日期，并修改老用户登录次数
			for_game.SetWishPlayerAccessLogToRedis(pid, cur0Time, "WishTime")
			if leadTime > easygo.A_DAY_SECOND {
				for_game.SetWishPlayerAccessLogToRedis(pid, 1, "RetainedDay")
			} else {
				retainedDay := for_game.IncrWishPlayerAccessLogToRedis(pid, 1, "RetainedDay")
				switch retainedDay {
				case 2:
					filed = "TwoDayKeep"
				case 3:
					filed = "ThreeDayKeep"
				case 7:
					filed = "SevenDayKeep"
				case 15:
					filed = "FifteenDayKeep"
				case 30:
					filed = "ThirtyDayKeep"
				default:
					filed = ""
				}
				if filed != "" {
					for_game.SetRedisWishLogReportFieldVal(cur0Time, 1, filed)
				}
			}
			for_game.SetRedisWishLogReportFieldVal(cur0Time, 1, "OldPlayer")
		}
	}
	for_game.InsertWishBPLog(pid, for_game.BP_WISH_ACCESS)
	for_game.SetRedisWishLogReportFieldVal(cur0Time, 1, "WishTime")
}

// 添加许愿池埋点记录服务
func AddReportWishLogService(pid int64, reqType int) {
	playerAccess := for_game.GetWishPlayerAccessLogFromRedis(pid)
	//访问许愿池
	if reqType == for_game.WISH_REPORT_ACCESS_WISH {
		AccessWishBPoint(pid, playerAccess)
		return
	}
	//无效用户无需处理
	if playerAccess == nil {
		return
	}

	cur0Time := easygo.GetToday0ClockTimestamp()
	switch reqType {
	case for_game.WISH_REPORT_ACCESS_EXCHANGE:
		if cur0Time-playerAccess.GetExchangeTime() > 0 {
			for_game.SetWishPlayerAccessLogToRedis(pid, cur0Time, "ExchangeTime")
			for_game.SetRedisWishLogReportFieldVal(cur0Time, 1, "ExchangeMen")
			for_game.InsertWishBPLog(pid, for_game.BP_WISH_EXCHANGE)
		}
	case for_game.WISH_REPORT_VEXCHANGE:
		if cur0Time-playerAccess.GetVExchangeTime() > 0 {
			for_game.SetWishPlayerAccessLogToRedis(pid, cur0Time, "VExchangeTime")
			for_game.SetRedisWishLogReportFieldVal(cur0Time, 1, "VExchangeMen")
		}
		for_game.InsertWishBPLog(pid, for_game.BP_WISH_VEXCHANGE)
		for_game.SetRedisWishLogReportFieldVal(cur0Time, 1, "VExchangeTime")
	case for_game.WISH_REPORT_ACCESS_DERE:
		if cur0Time-playerAccess.GetDareTime() > 0 {
			for_game.SetWishPlayerAccessLogToRedis(pid, cur0Time, "DareTime")
			for_game.SetRedisWishLogReportFieldVal(cur0Time, 1, "DareMen")
			for_game.InsertWishBPLog(pid, for_game.BP_WISH_DARE)
		}
	case for_game.WISH_REPORT_CHALLENGE:
		if cur0Time-playerAccess.GetChallengeTime() > 0 {
			for_game.SetWishPlayerAccessLogToRedis(pid, cur0Time, "ChallengeTime")
			for_game.SetRedisWishLogReportFieldVal(cur0Time, 1, "ChallengeMen")
			for_game.InsertWishBPLog(pid, for_game.BP_WISH_CHALLENGE)
		}
	}
}
