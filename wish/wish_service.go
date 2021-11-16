// 许愿池逻辑层
package wish

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/h5_wish"
	"game_server/pb/share_message"
	"runtime"
	"strconv"
	"time"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

const (
	WISH_BOX_NUM          = 3    // 首页顶部盲盒的个数
	DEFAULT_GET_COUNT_NUM = 10   // XXX获得多少个硬币默认 10条
	DEFAULT_MIN_COIN      = 45   //XXX获得多少个硬币 最少
	DEFAULT_MAX_COIN      = 1000 //XXX获得多少个硬币 最多
)

const (
	WISH_BRAND_HOT_NUM = 3  // 热门品牌数量
	WISH_BRAND_NUM     = 10 // 品牌数量
)

const (
	WISH_MIN_TIME = 25  // 最小的分钟数
	WISH_MAX_TIME = 100 // 最大的分钟数
)

const (
	TYPE_PRODUCT = 1 //  商品
	TYPE_BOX     = 2 // 盲盒
)

const (
	IS_RECOMMEND_NUM = 8 // 热门推荐条数
	RANGE_NUM        = 10
)

const (
	WISH_OCCUPY_MINUTE_MIN  = 25 * 60  // 运营号愿望守护最少秒
	WISH_OCCUPY_MINUTE_MAX  = 100 * 60 // 运营号愿望守护最多秒
	WISH_OCCUPY_DIAMOND_MAX = 1000     // 运营号愿望守护获取最少钻石
	WISH_OCCUPY_DIAMOND_MIN = 45       // 运营号愿望守护获取最多钻石
)
const (
//AFTER_TIME_BACK = 14 * 24 * 3600 // 14天过期
//AFTER_TIME_BACK = 30 // 14天过期

)

//奖项类型:1钻石，2实物
const (
	WISH_ACT_PRIZE_TYPE_DIAMOND = 1 // 1钻石，
	WISH_ACT_PRIZE_TYPE_PRODUCT = 2 //2实物
)

const (
	WISH_ACT_TOP_NUM = 1000 // 1000个排行
)

// 许愿池首页查询 service
func QueryBoxService(reqMsg *h5_wish.QueryBoxReq) (*h5_wish.QueryBoxResp, error) {
	queryBox := make([]*h5_wish.QueryBox, 0)
	switch reqMsg.GetType() {
	case TYPE_PRODUCT:
		// 可能会出现多个盲盒
		items := for_game.GetWishBoxItemByItemId(reqMsg.GetId())
		boxIds := make([]int64, 0)
		for _, v := range items {
			boxIds = append(boxIds, v.GetWishBoxId())
		}
		boxes := for_game.GetWishBoxsByIdsFromDB(boxIds)
		for _, v := range boxes {
			queryBox = append(queryBox, &h5_wish.QueryBox{
				Id:       easygo.NewInt64(v.GetId()),
				Name:     easygo.NewString(v.GetName()),
				Image:    easygo.NewString(v.GetIcon()),
				Price:    easygo.NewInt64(v.GetPrice()),
				TotalNum: easygo.NewInt32(v.GetTotalNum()),
				RareNum:  easygo.NewInt32(v.GetRareNum()),
				Label:    easygo.NewInt32(v.GetMatch()),
			})
		}
	case TYPE_BOX:
		boxes := for_game.GetWishBoxsByIdsFromDB([]int64{reqMsg.GetId()})
		for _, v := range boxes {
			queryBox = append(queryBox, &h5_wish.QueryBox{
				Id:       easygo.NewInt64(v.GetId()),
				Name:     easygo.NewString(v.GetName()),
				Image:    easygo.NewString(v.GetIcon()),
				Price:    easygo.NewInt64(v.GetPrice()),
				TotalNum: easygo.NewInt32(v.GetTotalNum()),
				RareNum:  easygo.NewInt32(v.GetRareNum()),
				Label:    easygo.NewInt32(v.GetMatch()),
			})
		}
	default:
		logs.Error("QueryBoxService type 有误,t: ", reqMsg.GetType())
		return nil, errors.New("参数有误")
	}
	result := &h5_wish.QueryBoxResp{
		Boxes: queryBox,
	}
	return result, nil
}

// 所有商品的名字或者盲盒的名字
func QueryBoxProductNameService() (*h5_wish.BoxProductNameResp, error) {
	boxProducts := make([]*h5_wish.BoxProductName, 0)
	// 商品
	wishItems := for_game.GetAllWishItem()
	for _, v := range wishItems {
		boxProducts = append(boxProducts, &h5_wish.BoxProductName{
			Id:   easygo.NewInt64(v.GetId()),
			Name: easygo.NewString(v.GetName()),
			Type: easygo.NewInt32(TYPE_PRODUCT), // 表示商品
		})
	}
	// 盲盒
	wishBox := for_game.GetAllWishBox()
	for _, v := range wishBox {
		boxProducts = append(boxProducts, &h5_wish.BoxProductName{
			Id:   easygo.NewInt64(v.GetId()),
			Name: easygo.NewString(v.GetName()),
			Type: easygo.NewInt32(TYPE_BOX), // 表示盲盒
		})
	}
	result := &h5_wish.BoxProductNameResp{
		BoxProductNames: boxProducts,
	}
	return result, nil
}

// 搜索发现
func SearchFoundService() (*h5_wish.SearchFoundResp, error) {
	SearchFounds := make([]*h5_wish.WishItemType, 0)
	// 商品
	wishItemTypes := for_game.GetRangeRecommendWishItem(IS_RECOMMEND_NUM)
	for _, v := range wishItemTypes {
		SearchFounds = append(SearchFounds, &h5_wish.WishItemType{
			Id:          easygo.NewInt64(v.GetId()),
			Name:        easygo.NewString(v.GetName()),
			IsRecommend: easygo.NewBool(v.GetIsRecommend()),
		})
	}
	result := &h5_wish.SearchFoundResp{
		WishItemTypes: SearchFounds,
	}
	return result, nil
}

// 商品展示区 原接口
//func ProductShowService(uid int64) (*h5_wish.ProductShowResp, error) {
//	// 找到盒子
//	boxes, err := for_game.FindWishBoxByNum(WISH_BOX_NUM)
//	if err != nil {
//		logs.Error("---ProductShowService FindWishBoxByNum err: ", err.Error())
//		return nil, err
//	}
//	bIds := make([]int64, 0)
//	for _, v := range boxes {
//		bIds = append(bIds, v.GetId())
//	}
//	// 批量找出许愿数据
//	wishDataList, _ := for_game.GetWishDataByBidsStatus(uid, bIds, for_game.WISH_CHALLENGE_WAIT)
//	wIdMap := make(map[int64]int64) // map[盲盒id]boxItemId
//	for _, v := range wishDataList {
//		wIdMap[v.GetWishBoxId()] = v.GetWishBoxItemId()
//	}
//	// 盒子中最贵的物品
//	productIds := make([]int64, 0)
//	products := make([]*h5_wish.Product, 0)
//	for _, v := range boxes {
//		// 找到价格最高的物品
//		item, err := for_game.FindMaxPriceBoxItemByBoxId(v.GetId(), productIds)
//		if err != nil {
//			return nil, err
//		}
//		productIds = append(productIds, item.GetWishItemId())
//		// 通过item 找到商品图片,头像
//		wishItem, err := for_game.QueryWishItemById(item.GetWishItemId())
//		if err != nil {
//			return nil, err
//		}
//		products = append(products, &h5_wish.Product{
//			Id:        easygo.NewInt64(item.GetId()),
//			Icon:      easygo.NewString(wishItem.GetIcon()),
//			Name:      easygo.NewString(wishItem.GetName()),
//			BoxId:     easygo.NewInt64(v.GetId()),
//			Price:     easygo.NewInt64(wishItem.GetPrice()),
//			ProductId: easygo.NewInt64(item.GetWishItemId()),
//			IsWish:    easygo.NewBool(item.GetId() == wIdMap[v.GetId()]),
//			Match:     easygo.NewInt32(v.GetMatch()),
//		})
//	}
//
//	result := &h5_wish.ProductShowResp{
//		Products: products,
//	}
//	return result, nil
//}

// 商品展示区
func ProductShowService(uid int64) (*h5_wish.ProductShowResp, error) {
	boxes := for_game.GetRecommendBox()
	boxIds := make([]int64, 0)
	//必须保证盲盒表中的中间表id与物品id一一对应
	boxProductMap := make(map[int64]int64)           //[盲盒id]物品id(排序第一/许愿)
	boxItemMap := make(map[int64]int64)              //[盲盒id]中间表id(排序第一/许愿)
	wishItemMap := make(map[int64]int64)             //[盲盒id]中间表id ：用于判断是否许愿款
	boxMap := make(map[int64]*share_message.WishBox) // //[盲盒id]盲盒信息
	for _, v := range boxes {
		boxIds = append(boxIds, v.GetId())
		boxProductMap[v.GetId()] = v.GetWishItems()[0]
		boxItemMap[v.GetId()] = v.GetItems()[0]
		boxMap[v.GetId()] = v
	}

	// 获取许愿的物品
	wishData, _ := for_game.GetWishDataByPid(uid, for_game.WISH_CHALLENGE_WAIT, boxIds)
	for _, v := range wishData {
		boxProductMap[v.GetWishBoxId()] = v.GetWishItemId()
		boxItemMap[v.GetWishBoxId()] = v.GetWishBoxItemId()
		wishItemMap[v.GetWishBoxId()] = v.GetWishBoxItemId()
	}

	allItemIds := make([]int64, 0)
	for _, v := range boxProductMap {
		allItemIds = append(allItemIds, v)
	}

	//获取所有盲盒物品信息
	allItems, _ := for_game.GetWishItemByIdsFromDB(allItemIds)
	itemMap := make(map[int64]*share_message.WishItem, 0) //[物品id]物品信息
	for _, item := range allItems {
		itemMap[item.GetId()] = item
	}

	products := make([]*h5_wish.Product, 0)
	var boxItemId int64
	var item *share_message.WishItem
	var box *share_message.WishBox
	for boxId, productId := range boxProductMap {
		ok, isWish := false, false
		boxItemId, ok = boxItemMap[boxId]
		if !ok {
			continue
		}
		item, ok = itemMap[productId]
		if !ok {
			continue
		}
		box, ok = boxMap[boxId]
		if !ok {
			continue
		}
		_, ok = wishItemMap[boxId]
		if ok {
			isWish = true
		}
		products = append(products, &h5_wish.Product{
			Id:        easygo.NewInt64(boxItemId),
			ProductId: easygo.NewInt64(item.GetId()),
			Icon:      easygo.NewString(item.GetIcon()),
			Name:      easygo.NewString(item.GetName()),
			BoxId:     easygo.NewInt64(boxId),
			Match:     easygo.NewInt32(box.GetMatch()),
			IsWish:    easygo.NewBool(isWish),
		})

	}

	result := &h5_wish.ProductShowResp{
		Products: products,
	}
	return result, nil
}

// xxx 已获得8888硬币
func GetCoinService() (*h5_wish.GetCoinResp, error) {
	getCoins := make([]*h5_wish.GetCoin, 0)
	occupied := for_game.GetSomeWishOccupied(DEFAULT_GET_COUNT_NUM)
	for _, v := range occupied {
		getCoins = append(getCoins, &h5_wish.GetCoin{
			PlayerId: v.PlayerId,
			Name:     v.NickName,
			HeadUrl:  v.HeadUrl,
			Coin:     easygo.NewInt64(v.GetCoinNum()),
		})
	}
	result := &h5_wish.GetCoinResp{
		GetCoins: getCoins,
	}
	return result, nil

}

// 玩家获得信息获取商品详细信息  // [中间表物品id]用户id,   返回[最底层商品id]用户id
func ProductInfoFromPlayerWishItem(pMap map[int64]int64) (map[int64]int64, []*share_message.WishItem) {
	if len(pMap) <= 0 {
		return nil, nil
	}
	boxItemIds := make([]int64, 0)
	for k, _ := range pMap {
		boxItemIds = append(boxItemIds, k)
	}
	if len(boxItemIds) <= 0 {
		return nil, nil
	}
	// 找到中间表中的物品,获得真正物品id
	items, err := for_game.GetWishBoxItemByIdsFromDB(boxItemIds)
	if err != nil {
		return nil, nil
	}
	productIds := make([]int64, 0)
	boxItemMap := make(map[int64]int64) // [WishItemId]用户id
	for _, v := range items {
		productIds = append(productIds, v.GetWishItemId())
		boxItemMap[v.GetWishItemId()] = pMap[v.GetId()]
	}
	if len(productIds) <= 0 {
		return nil, nil
	}
	products, err := for_game.GetWishItemByIdsFromDB(productIds)
	if err != nil {
		return nil, nil
	}

	return boxItemMap, products
}

// 首页消息播放区
func HomeMessageService() (*h5_wish.HomeMessageResp, error) {
	// 随机获取10条挑战成功的信息
	products := for_game.GetRandProducts(DEFAULT_GET_COUNT_NUM)
	// 封装头像
	pMap := make(map[int64]int64) // [中间表物品id]用户id
	for _, v := range products {
		pMap[v.GetChallengeItemId()] = v.GetPlayerId()
	}
	// 封装用户和商品信息
	productPlayerMap, productList := ProductInfoFromPlayerWishItem(pMap)
	getProducts := make([]*h5_wish.GetProduct, 0)
	for _, v := range productList {
		var nickName, headUrl string
		if base := for_game.GetRedisWishPlayer(productPlayerMap[v.GetId()]); base != nil {
			nickName = base.GetNickName()
			headUrl = base.GetHeadIcon()
		}
		getProducts = append(getProducts, &h5_wish.GetProduct{
			PlayerId:    easygo.NewInt64(productPlayerMap[v.GetId()]),
			PlayerName:  easygo.NewString(nickName),
			HeadUrl:     easygo.NewString(headUrl),
			ProductName: easygo.NewString(v.GetName()),
			Image:       easygo.NewString(v.GetIcon()),
		})
	}
	if len(products) < DEFAULT_GET_COUNT_NUM { // 抽取运营号
		players := for_game.GetRandWishPlayer(DEFAULT_GET_COUNT_NUM - len(products))
		for _, v := range players {
			item := for_game.GetRandWishItem(1)
			getProducts = append(getProducts, &h5_wish.GetProduct{
				PlayerId:    easygo.NewInt64(v.GetId()),
				PlayerName:  easygo.NewString(v.GetNickName()),
				HeadUrl:     easygo.NewString(v.GetHeadUrl()),
				ProductName: easygo.NewString(item[0].GetName()),
				Image:       easygo.NewString(item[0].GetIcon()),
			})
		}
	}

	// 查询挑战成功记录
	wishLogs := for_game.GetRangeWishLog(DEFAULT_GET_COUNT_NUM, true)
	dreMessages := make([]*h5_wish.DareMessage, 0)

	for _, v := range wishLogs {
		var name, headURL string
		if base := for_game.GetRedisWishPlayer(v.GetDareId()); base != nil {
			name = base.GetNickName()
			headURL = base.GetHeadIcon()
		}
		dreMessages = append(dreMessages, &h5_wish.DareMessage{
			PlayerId: easygo.NewInt64(v.GetDareId()),
			Name:     easygo.NewString(name),
			HeadUrl:  easygo.NewString(headURL),
		})
	}

	// 如果不够,抽取运营号
	if len(wishLogs) < DEFAULT_GET_COUNT_NUM {
		// 抽取运营号
		players := for_game.GetRandWishPlayer(DEFAULT_GET_COUNT_NUM - len(wishLogs))
		for _, v := range players {
			dreMessages = append(dreMessages, &h5_wish.DareMessage{
				PlayerId: easygo.NewInt64(v.GetId()),
				Name:     easygo.NewString(v.GetNickName()),
				HeadUrl:  easygo.NewString(v.GetHeadUrl()),
			})
		}
	}

	// 随机获取10条信息获得物品的信息.
	result := &h5_wish.HomeMessageResp{
		DareMessages: dreMessages,
		GetProducts:  getProducts,
	}
	return result, nil
}

// 获取随机十条物品信息
func GetRandProductService(boxId int64) (*h5_wish.RandProductResp, error) {
	// 随机获取10条挑战成功的信息
	var products []*share_message.PlayerWishItem
	if boxId < 1 {
		products = for_game.GetRandProducts(DEFAULT_GET_COUNT_NUM)
	} else {
		products = for_game.GetRandBoxProducts(boxId, DEFAULT_GET_COUNT_NUM)
	}

	wishItemIds := make([]int64, 0)                    // 物品id
	playerIds := make([]int64, 0)                      // 用户id 与wishItemIds一一对应
	pMap := make(map[int64]*share_message.WishItem, 0) // [商品id] 商品数据
	for _, v := range products {
		wishItemIds = append(wishItemIds, v.GetWishItemId())
		playerIds = append(playerIds, v.GetPlayerId())
	}
	// 封装用户和商品信息
	productList, err := for_game.GetWishItemByIdsFromDB(wishItemIds)
	if err != nil {
		return nil, errors.New("获取随机十条物品信息失败")
	}
	for _, v := range productList {
		pMap[v.GetId()] = v
	}

	getProducts := make([]*h5_wish.GetProduct, 0)
	for i, v := range wishItemIds {
		var nickName, headUrl string
		if base := for_game.GetRedisWishPlayer(playerIds[i]); base != nil {
			nickName = base.GetNickName()
			headUrl = base.GetHeadIcon()
		}
		getProducts = append(getProducts, &h5_wish.GetProduct{
			PlayerId:    easygo.NewInt64(playerIds[i]),
			PlayerName:  easygo.NewString(nickName),
			HeadUrl:     easygo.NewString(headUrl),
			ProductName: easygo.NewString(pMap[v].GetName()),
			Image:       easygo.NewString(pMap[v].GetIcon()),
		})
	}
	if len(getProducts) < DEFAULT_GET_COUNT_NUM { // 抽取运营号
		players := for_game.GetRandWishPlayer(DEFAULT_GET_COUNT_NUM - len(getProducts))
		var itemInfo []*share_message.WishItem
		if boxId > 0 {
			box, err := for_game.GetWishBoxJustById(boxId)
			if err != nil {
				return nil, errors.New("获取物品轮播数据失败失败")
			}
			items := box.GetWishItems()
			itemInfo, err = for_game.GetWishItemByIdsFromDB(items)
			if err != nil || len(itemInfo) < 1 {
				return nil, errors.New("获取物品轮播数据失败失败")
			}
		} else {
			itemInfo = for_game.GetRandWishItem(DEFAULT_GET_COUNT_NUM - len(getProducts))
		}
		itemLen := len(itemInfo)
		if itemLen > 0 {
			for _, v := range players {
				num := util.RandIntn(itemLen)
				getProducts = append(getProducts, &h5_wish.GetProduct{
					PlayerId:    easygo.NewInt64(v.GetId()),
					PlayerName:  easygo.NewString(v.GetNickName()),
					HeadUrl:     easygo.NewString(v.GetHeadUrl()),
					ProductName: easygo.NewString(itemInfo[num].GetName()),
					Image:       easygo.NewString(itemInfo[num].GetIcon()),
				})
			}
		}

	}

	// 随机获取10条物品的信息.
	result := &h5_wish.RandProductResp{
		GetProducts: getProducts,
	}
	return result, nil
}

// 获取十条挑战成功记录
func GetDareMessageService() (*h5_wish.DareMessageResp, error) {
	// 查询挑战成功记录
	wishLogs := for_game.GetRangeWishLog(DEFAULT_GET_COUNT_NUM, true)
	dreMessages := make([]*h5_wish.DareMessage, 0)

	var name, headURL string
	for _, v := range wishLogs {
		if base := for_game.GetRedisWishPlayer(v.GetDareId()); base != nil {
			name = base.GetNickName()
			headURL = base.GetHeadIcon()
		}
		dreMessages = append(dreMessages, &h5_wish.DareMessage{
			PlayerId: easygo.NewInt64(v.GetDareId()),
			Name:     easygo.NewString(name),
			HeadUrl:  easygo.NewString(headURL),
		})
	}

	// 随机获取10条挑战成功记录
	result := &h5_wish.DareMessageResp{
		DareMessages: dreMessages,
	}
	return result, nil
}

// 挑战守护者赢硬币 快捷入口的轮播展示入口 (后台配置一定会有的,所以只需要查询有挑战者的盲盒)
//func ProtectorService() (*h5_wish.ProtectorResp, error) {
//	ids := for_game.GetRangWishLogGuardians(DEFAULT_GET_COUNT_NUM)
//	// 通过守护者id找到对应的盲盒id
//	protectors := make([]*h5_wish.Protector, 0)
//	for _, v := range ids {
//		var name, headUrl string
//		base := for_game.GetRedisWishPlayer(v)
//		if base == nil {
//			continue
//		}
//		name = base.GetNickName()
//		headUrl = base.GetHeadIcon()
//		boxs := for_game.GetRangWishLogByGuardian(v)
//		boxIds := make([]int64, 0)
//		for _, box := range boxs {
//			boxIds = append(boxIds, box.GetId())
//		}
//		protectors = append(protectors, &h5_wish.Protector{
//			PlayerId:   easygo.NewInt64(v),
//			PlayerName: easygo.NewString(name),
//			HeadUrl:    easygo.NewString(headUrl),
//			BoxIds:     boxIds,
//		})
//	}
//	result := &h5_wish.ProtectorResp{
//		Protectors: protectors,
//	}
//	return result, nil
//}

// 挑战守护者赢硬币 快捷入口的轮播展示入口 (后台配置一定会有的,所以只需要查询有挑战者的盲盒)
func ProtectorService() *h5_wish.ProtectorDataResp {
	boxList := for_game.GetRandWishBox(DEFAULT_GET_COUNT_NUM)
	list := make([]*h5_wish.ProtectorData, 0)
	for _, v := range boxList {
		var headUrl string
		if v.GetGuardianId() > 0 && v.GetIsGuardian() {
			if p := for_game.GetRedisWishPlayer(v.GetGuardianId()); p != nil {
				headUrl = p.GetHeadIcon()
			}
		}
		list = append(list, &h5_wish.ProtectorData{
			BoxId:   easygo.NewInt64(v.GetId()),
			HeadUrl: easygo.NewString(headUrl),
		})

	}
	return &h5_wish.ProtectorDataResp{
		Protectors: list,
	}
}

// 最新上线,人气盲盒,欧气爆棚 菜单栏
func MenuService() (*h5_wish.MenuResp, error) {
	list := for_game.GetMenuList()
	menus := make([]*h5_wish.Menu, 0)
	for _, v := range list {
		menus = append(menus, &h5_wish.Menu{
			Id:   easygo.NewInt64(v.GetId()),
			Name: easygo.NewString(v.GetName()),
		})
	}
	resp := &h5_wish.MenuResp{
		Menus: menus,
	}
	return resp, nil
}

//综合下面的,商品品牌(苹果,雷神,三星,罗技)
func ProductBrandListService() (*h5_wish.ProductBrandListResp, error) {
	productBrands := make([]*h5_wish.ProductBrand, 0)
	// 先找到3个热门的
	hot := for_game.GetRangeHotWishBrandByNum(WISH_BRAND_HOT_NUM)
	for _, v := range hot {
		productBrands = append(productBrands, &h5_wish.ProductBrand{
			Id:    easygo.NewInt64(v.GetId()),
			Name:  easygo.NewString(v.GetName()),
			IsHot: easygo.NewBool(v.GetIsHot()),
		})
	}
	common := for_game.GetRangeWishBrandByNum(WISH_BRAND_NUM - len(hot))
	for _, v := range common {
		productBrands = append(productBrands, &h5_wish.ProductBrand{
			Id:   easygo.NewInt64(v.GetId()),
			Name: easygo.NewString(v.GetName()),
		})
	}
	return &h5_wish.ProductBrandListResp{
		ProductBrandList: productBrands,
	}, nil
}

// 修改点击数.
func AddClickNum(brandIds, typeIds []int64) {
	if len(brandIds) > 0 {
		for _, v := range brandIds {
			if err := for_game.SetBrandClickNum(v); err != nil {
				logs.Error("SetBrandClickNum 错误,主键id 为: %d,err: %s", v, err.Error())
			}
		}
	}
	if len(typeIds) > 0 {
		for _, v := range typeIds {
			if err := for_game.SetTypeClickNum(v); err != nil {
				logs.Error("SetTypeClickNum 错误,主键id 为: %d,err: %s", v, err.Error())
			}
		}
	}
}

// 盲盒区列表(综合,最新上线,人气盲盒,欧气爆棚条件筛选)
func SearchBoxService(uid int64, reqMsg *h5_wish.SearchBoxReq) (*h5_wish.SearchBoxResp, error) {
	// 异步处理点击数
	easygo.Spawn(AddClickNum, reqMsg.GetWishBrandId(), reqMsg.GetWishItemTypeId())

	page, pageSize := for_game.MakePageAndPageSize(reqMsg.GetPage(), reqMsg.GetPageSize())
	boxes, count := for_game.GetBoxListByPage(reqMsg, page, pageSize)
	boxShowList := make([]*h5_wish.BoxShow, 0)

	// 挑战赛的id列表
	ids := make([]int64, 0)
	for _, v := range boxes {
		if v.GetMatch() == WISH_DARE {
			ids = append(ids, v.GetId())
		}

		makeWishCount := for_game.GetPlayerMakeWish(uid, v.GetId())
		var redisMakeWish bool
		if makeWishCount == 0 {
			redisMakeWish = true
		}
		var haveIsWish bool
		// 判断是否已是新人
		if redisMakeWish {
			haveIsWish = v.GetHaveIsWin()
		}
		boxShowList = append(boxShowList, &h5_wish.BoxShow{
			Id:            easygo.NewInt64(v.GetId()),
			Name:          easygo.NewString(v.GetName()),
			Image:         easygo.NewString(v.GetIcon()),
			Desc:          easygo.NewString(v.GetDesc()),
			Price:         easygo.NewInt64(v.GetPrice()),
			WishFishCount: easygo.NewInt64(v.GetWinNum()),
			Label:         easygo.NewInt32(v.GetMatch()), // 0-非挑战赛,1-挑战赛
			ProductStatus: easygo.NewInt32(v.GetProductStatus()),
			TotalNum:      easygo.NewInt32(v.GetTotalNum()),
			RareNum:       easygo.NewInt32(v.GetRareNum()),
			HaveIsWin:     easygo.NewBool(haveIsWish),
		})
	}
	if len(ids) > 0 {
		occupieds := for_game.GetWishOccupiedByBoxIds(ids)
		for _, v := range occupieds {
			for _, v1 := range boxShowList {
				if v.GetWishBoxId() == v1.GetId() {
					v1.HeadUrl = easygo.NewString(v.GetHeadUrl())
					v1.CoinNum = easygo.NewInt32(v.GetCoinNum())
					v1.OccupiedTime = easygo.NewInt64(time.Now().Unix() - v.GetCreateTime())
					v1.GuardianPlayerId = easygo.NewInt64(v.GetPlayerId())
				}
			}
		}
	}

	return &h5_wish.SearchBoxResp{
		BoxShowList: boxShowList,
		Count:       easygo.NewInt32(count),
	}, nil
}

// 物品品牌列表
func BrandListService() (*h5_wish.BrandListResp, error) {
	// 从 redis 中获取
	/*	value := for_game.GetWishBrand()
		if value != "" {
			return &h5_wish.BrandListResp{
				BrandList: easygo.NewString(value),
			}, nil
		}*/

	result := make(map[string][]*share_message.WishBrand) // key 为0 的数据表示是热门的数据.

	hotList := for_game.GetHotWishBrandList()
	hotIds := make([]int64, 0)
	for _, v := range hotList {
		hotIds = append(hotIds, v.GetId())
	}
	if len(hotList) > 0 {
		result["0"] = hotList
	}

	list, list1 := for_game.GetNotHotWishBrandList(hotIds, RANGE_NUM)
	if len(hotList) > 0 {
		result["1"] = list1 // 热门
	}
	mMap := make(map[string]string) // [品牌名字]字母
	mss := make([]string, 0)        // 字母的数组
	for _, m := range list {
		mMap[m.GetName()] = m.GetType()
		// 判断是否在数组
		if !util.InStringSlice(m.GetType(), mss) {
			mss = append(mss, m.GetType())
		}
	}

	for _, m11 := range mss {
		for _, v := range list {
			if m11 == v.GetType() {
				value, ok := result[m11]
				if !ok {
					value = make([]*share_message.WishBrand, 0)
				}
				value = append(value, v)
				result[m11] = value
			}
		}
	}
	marshal, _ := json.Marshal(result)
	// 设置进redis
	//for_game.SetWishBrand(string(marshal))
	return &h5_wish.BrandListResp{
		BrandList: easygo.NewString(string(marshal)),
	}, nil
}

// 物品类别列表
func ProductTypeListService() (*h5_wish.TypeListResp, error) {
	// 从redis中查找
	/*	value := for_game.GetWishItemType()
		if value != "" {
			return &h5_wish.TypeListResp{
				TypeList: easygo.NewString(value),
			}, nil
		}*/
	result := make(map[string][]*share_message.WishItemType)

	hotList := for_game.GetHotWishTypeList()
	hotIds := make([]int64, 0)
	for _, v := range hotList {
		hotIds = append(hotIds, v.GetId())
	}
	if len(hotList) > 0 {
		result["0"] = hotList // 热门
	}
	list, list1 := for_game.GetNotHotWishTypeList(hotIds, RANGE_NUM)
	if len(list1) > 0 {
		result["1"] = list1 // 权重排序
	}
	mMap := make(map[string]string)
	mss := make([]string, 0)
	for _, m := range list {
		mMap[m.GetName()] = m.GetType()
		// 判断是否在数组
		if !util.InStringSlice(m.GetType(), mss) {
			mss = append(mss, m.GetType())
		}
	}

	for _, m11 := range mss {
		for _, v := range list {
			if m11 == v.GetType() {
				value, ok := result[m11]
				if !ok {
					value = make([]*share_message.WishItemType, 0)
				}
				value = append(value, v)
				result[m11] = value
			}
		}
	}
	marshal, _ := json.Marshal(result)
	//for_game.SetWishItemType(string(marshal))
	return &h5_wish.TypeListResp{
		TypeList: easygo.NewString(string(marshal)),
	}, nil
}

// 挑战赛主界面-消息播放区
func DareRecommendService() (*h5_wish.DareRecommendResp, error) {
	protectors := make([]*h5_wish.Protector, 0)
	// 随机从 WishOccupied 表中拉去数据
	//occupied := for_game.GetRangeWishOccupied(RANGE_NUM)
	occupied := for_game.GetSomeWishOccupied(RANGE_NUM)
	for _, v := range occupied {
		protectors = append(protectors, &h5_wish.Protector{
			PlayerId:      easygo.NewInt64(v.GetPlayerId()),
			PlayerName:    easygo.NewString(v.GetNickName()),
			HeadUrl:       easygo.NewString(v.GetHeadUrl()),
			BoxIds:        []int64{v.GetWishBoxId()},
			ProtectorTime: easygo.NewInt64(v.GetOccupiedTime()),
			Coin:          easygo.NewInt64(v.GetCoinNum()),
		})
	}

	return &h5_wish.DareRecommendResp{
		Protectors: protectors,
	}, nil
}

// 排行榜列表
func RankingsService() (*h5_wish.RankingResp, error) {
	topLogs := for_game.GetTop10FromDB()
	pids := make([]int64, 0)
	for _, v := range topLogs {
		pids = append(pids, v.GetId())
	}
	//
	players, _ := for_game.GetWishPlayerByIds(pids)

	kinds := make([]*h5_wish.Ranking, 0)
	for _, v := range topLogs {
		var name, headUrl string
		for _, p := range players {
			if p.GetId() == v.GetPlayerId() {
				name = p.GetNickName()
				headUrl = p.GetHeadUrl()
			}
		}
		kinds = append(kinds, &h5_wish.Ranking{
			PlayerId:       easygo.NewInt64(v.GetPlayerId()),
			ProtectorCount: easygo.NewInt64(v.GetWishNum()),
			CoinCount:      easygo.NewInt64(v.GetCoinNum()),
			Name:           easygo.NewString(name),
			HeadUrl:        easygo.NewString(headUrl),
		})

	}
	return &h5_wish.RankingResp{
		Rankings: kinds,
	}, nil
}

// 我的战绩头部的挑战成功数和硬币数
func MyRecordService(pid int64) (*h5_wish.MyRecordResp, error) {
	//wishLog
	wishLog := for_game.GetAllWishOccupiedByPid(pid)
	var coin int32
	for _, v := range wishLog {
		coin += v.GetCoinNum()
	}
	return &h5_wish.MyRecordResp{
		DareCount:      easygo.NewInt32(len(wishLog)),
		TotalCoinCount: easygo.NewInt64(coin),
	}, nil
}

// 我的战绩列表
func MyDareService(pid int64, reqMsg *h5_wish.MyDareReq) (*h5_wish.MyDareResp, error) {
	page, pageSize := for_game.MakePageAndPageSize(reqMsg.GetPage(), reqMsg.GetPageSize())
	wishOccupied, count := for_game.GetAllWishOccupiedByPage(pid, page, pageSize)
	boxIds := make([]int64, 0)
	for _, v := range wishOccupied {
		boxIds = append(boxIds, v.GetWishBoxId())
	}
	m := GetIconByBoxIds(boxIds)
	myDare := make([]*h5_wish.MyDare, 0)
	for _, v := range wishOccupied {
		// 计算占领时间
		t := v.GetOccupiedTime()
		if v.GetStatus() == for_game.WISH_OCCUPIED_STATUS_UP { // 占领中重新计算时间
			t = time.Now().Unix() - v.GetCreateTime()
		}

		myDare = append(myDare, &h5_wish.MyDare{
			BoxId:         easygo.NewInt64(v.GetWishBoxId()),
			Image:         easygo.NewString(m[v.GetWishBoxId()]),
			ProtectorTime: easygo.NewInt64(t),
			CoinCount:     easygo.NewInt64(v.GetCoinNum()),
		})
	}

	return &h5_wish.MyDareResp{
		Dares: myDare,
		Count: easygo.NewInt32(count),
	}, nil
}

// 通过盲盒id列表找到对应的图片
func GetIconByBoxIds(boxIds []int64) map[int64]string {
	m := make(map[int64]string) // [boxId]icon
	if len(boxIds) <= 0 {
		return m
	}
	boxes := for_game.GetWishBoxsByIdsFromDB(boxIds)
	for _, v := range boxes {
		m[v.GetId()] = v.GetIcon()
	}
	return m
}

// xxx 发起了挑战列表列表
func DareListService(reqMsg *h5_wish.DareReq) (*h5_wish.DareResp, error) {
	result := make([]*h5_wish.WhoDare, 0)
	// 找到当前盲盒挑战成功的人
	item := for_game.GetRangeWishWishOccupiedByNum(reqMsg.GetBoxId(), 1)
	num := DEFAULT_GET_COUNT_NUM
	if len(item) != 0 {
		protectorTime := time.Now().Unix() - item[0].GetCreateTime()
		if item[0].GetOccupiedTime() > 0 {
			protectorTime = item[0].GetOccupiedTime()
		}
		num = DEFAULT_GET_COUNT_NUM - 1
		result = append(result, &h5_wish.WhoDare{
			PlayerId:      easygo.NewInt64(item[0].GetPlayerId()),
			Name:          easygo.NewString(item[0].GetNickName()),
			IsSuccess:     easygo.NewBool(true),
			HeadIcon:      easygo.NewString(item[0].GetHeadUrl()),
			ProtectorTime: easygo.NewInt64(protectorTime),
		})
	}
	// 指定条数随机选取对应的盲盒的挑战记录
	wishLogs := for_game.GetRangeWishLogByBoxId(reqMsg.GetBoxId(), num)
	for _, v := range wishLogs {
		result = append(result, &h5_wish.WhoDare{
			PlayerId:  easygo.NewInt64(v.GetDareId()),
			Name:      easygo.NewString(v.GetDareName()),
			IsSuccess: easygo.NewBool(v.GetResult()),
			HeadIcon:  easygo.NewString(v.GetDareHeadIcon()),
		})
	}
	return &h5_wish.DareResp{
		Dares: result,
	}, nil
}

func BoxInfoService(uid, boxId int64, reqType int32) (*h5_wish.BoxResp, error) {
	result := &h5_wish.BoxResp{}
	//获取盲盒
	CheckGuardianIsEp(boxId) // 更新盲盒守护者信息.
	box, err := for_game.GetWishBoxJustById(boxId)

	if err != nil {
		return result, err
	}

	if box.GetId() == 0 {
		logs.Error("盲盒数据为空,boxId为: %d", boxId)
		return &h5_wish.BoxResp{}, errors.New("数据有误")
	}

	// 如果盲盒是大亏情况下,修改盲盒的状态未补货中
	status := GetPoolStatus(box.GetWishPoolId())
	if status == for_game.POOL_STATUS_BIGLOSS {
		box.Status = easygo.NewInt32(2) //状态:0下架，1上架;2-积极补货中
		_ = for_game.UpsertBox(box.GetId(), box)
	}
	result.BoxId = easygo.NewInt64(boxId)
	if reqType == WISH_DARE { //  1-挑战赛 2-非挑战赛
		if defender := for_game.GetBoxDefenderByBoxId(boxId); defender != nil {
			result.Protector = easygo.NewString(defender.GetNickName())
			result.ProtectorId = easygo.NewInt64(defender.GetPlayerId())
			result.ProtectorTime = easygo.NewInt64(defender.GetOccupiedTime())
			result.ProtectorHeadUrl = easygo.NewString(defender.GetHeadUrl())
			result.CreateTime = easygo.NewInt64(defender.GetCreateTime())
		}
	}

	// 查询收藏表
	collection := for_game.GetPlayerWishCollection(uid, boxId)
	result.IsCollection = easygo.NewBool(collection.GetId() != 0)
	result.Status = easygo.NewInt32(box.GetStatus())
	//获取用户许愿款
	wishData, _ := for_game.GetWishDataByStatus(uid, boxId, for_game.WISH_CHALLENGE_WAIT)
	wishId := wishData.GetWishBoxItemId() // 107

	//根据盲盒ID找到物品配置信息
	boxItems := box.GetItems()
	if len(boxItems) == 0 {
		logs.Error("盲盒中的物品为空,盲盒id为: %d", boxId)
		return result, errors.New("盲盒为空")
	}
	wishBoxItems, _ := for_game.GetWishBoxItemByIdsFromDB(boxItems) // 盲盒物品列表(中间表)
	if len(wishBoxItems) == 0 {
		logs.Error("盲盒中的物品为空,盲盒id为: %d", boxId)
		return result, errors.New("盲盒为空")
	}
	itemIds := make([]int64, 0)
	boxItemMap := make(map[int64]*share_message.WishBoxItem) //[最底层物品的id]中间表物品
	for _, v := range wishBoxItems {
		itemIds = append(itemIds, v.GetWishItemId())
		boxItemMap[v.GetWishItemId()] = v
	}
	items, _ := for_game.GetWishItemByIdsFromDB(itemIds)
	if len(items) == 0 {
		logs.Error("盲盒中间表中最底层的的物品为空,盲盒id为: %d", boxId)
		return result, errors.New("盲盒为空")
	}
	products := make([]*h5_wish.Product, 0)
	newItems := make([]*share_message.WishItem, 0)
	for _, v := range itemIds {
		for _, v1 := range items {
			if v == v1.GetId() {
				newItems = append(newItems, v1)
			}
		}
	}
	makeWishCount := for_game.GetPlayerMakeWish(uid, boxId)
	var redisMakeWish bool
	if makeWishCount == 0 {
		redisMakeWish = true
	}
	var isMakeWish bool
	for _, v := range newItems {
		boxItem := boxItemMap[v.GetId()]
		if redisMakeWish {
			isMakeWish = boxItem.GetIsWin()
		}
		products = append(products, &h5_wish.Product{
			Id:          easygo.NewInt64(boxItem.GetId()),
			ProductType: easygo.NewInt32(boxItem.GetStyle()),
			Icon:        easygo.NewString(v.GetIcon()),
			Name:        easygo.NewString(v.GetName()),
			Price:       easygo.NewInt64(v.GetPrice()),
			ProductId:   easygo.NewInt64(v.GetId()),
			PreHaveTime: easygo.NewInt64(v.GetPreHaveTime()),
			IsPreSale:   easygo.NewBool(v.GetIsPreSale()),
			BoxName:     easygo.NewString(box.GetName()),
			IsWish:      easygo.NewBool(wishId == boxItem.GetId()),
			BoxId:       easygo.NewInt64(boxId),
			IsMakeWish:  easygo.NewBool(isMakeWish),
		})
	}
	result.BoxPrice = easygo.NewInt64(box.GetPrice())
	result.ProductList = products
	result.BoxIcon = easygo.NewString(box.GetIcon())
	result.ProductStatus = easygo.NewInt32(box.GetProductStatus())
	return result, nil
}

// 挑战记录,占领时长都是这个接口.
func DareRecordService(reqMsg *h5_wish.DareRecordReq) (*h5_wish.DareRecordResp, error) {
	page, pageSize := for_game.MakePageAndPageSize(reqMsg.GetPage(), reqMsg.GetPageSize())
	dareRecordResp := &h5_wish.DareRecordResp{}
	switch reqMsg.GetType() {
	case WISH_DARE_RECORED: //挑战记录
		wishLogList, wishLogCount := for_game.GetWishLogList(reqMsg.GetBoxId(), page, pageSize)
		wishLogListResult := make([]*h5_wish.WishLog, 0)
		for _, v := range wishLogList {
			wishLogListResult = append(wishLogListResult, &h5_wish.WishLog{
				Id:              easygo.NewInt64(v.GetId()),
				WishBoxId:       easygo.NewInt64(v.GetWishBoxId()),
				DareId:          easygo.NewInt64(v.GetDareId()),
				DareName:        easygo.NewString(v.GetDareName()),
				BeDareId:        easygo.NewInt64(v.GetBeDareId()),
				BeDareName:      easygo.NewString(v.GetBeDareName()),
				CreateTime:      easygo.NewInt64(v.GetCreateTime()),
				Result:          easygo.NewBool(v.GetResult()),
				ChallengeItemId: easygo.NewInt64(v.GetChallengeItemId()),
				DareHeadIcon:    easygo.NewString(v.GetDareHeadIcon()),
				DefendTime:      easygo.NewInt64(v.GetDefendTime()),
			})
		}
		dareRecordResp.WishLogList = wishLogListResult
		dareRecordResp.WishLogCount = easygo.NewInt32(wishLogCount)
	case WISH_DARE_HOLD_TIME: // 占领时长
		//wishOccupiedList, wishOccupiedCount := for_game.GetWishOccupiedList(reqMsg.GetBoxId(), page, pageSize)
		// todo 同步占领时长
		box1, _ := for_game.GetWishBox(reqMsg.GetBoxId(), []int64{for_game.WISH_DOWN_STATUS, for_game.WISH_PUT_ON_STATUS, for_game.WISH_ADD_PRODUCT_STATUS})
		if box1 != nil && box1.GetGuardianId() > 0 {
			for_game.UpOccupied1(reqMsg.GetBoxId(), box1.GetGuardianId())
		}

		wishOccupiedList, wishOccupiedCount := for_game.GetWishSumOccupied(reqMsg.GetBoxId(), page, pageSize)
		wishOccupiedListResult := make([]*h5_wish.WishOccupied, 0)
		for _, v := range wishOccupiedList {
			wishOccupiedListResult = append(wishOccupiedListResult, &h5_wish.WishOccupied{
				Id:        easygo.NewInt64(v.GetId()),
				WishBoxId: easygo.NewInt64(v.GetWishBoxId()),
				NickName:  easygo.NewString(v.GetNickName()),
				HeadUrl:   easygo.NewString(v.GetHeadUrl()),
				PlayerId:  easygo.NewInt64(v.GetPlayerId()),
				//CreateTime:   easygo.NewInt64(v.GetCreateTime()),
				//EndTime:      easygo.NewInt64(v.GetEndTime()),
				OccupiedTime: easygo.NewInt64(v.GetOccupiedTime()),
				//Status:       easygo.NewInt32(v.GetStatus()),
				//CoinNum:      easygo.NewInt32(v.GetCoinNum()),
			})
		}
		dareRecordResp.WishOccupiedList = wishOccupiedListResult
		dareRecordResp.WishOccupiedCount = easygo.NewInt32(wishOccupiedCount)
	}
	return dareRecordResp, nil
}

// 商品详情
func ProductDetailService(productId int64) (*h5_wish.ProductDetail, error) {
	// 中间表查询物品信息
	boxItem, err := for_game.GetWishBoxItem(productId)
	if err != nil {
		logs.Error("中间表没有数据,中间表id为: %d,err: %s", productId, err.Error())
		return nil, err
	}
	wishItem := for_game.GetWishItemByIdFromDB(boxItem.GetWishItemId())

	return &h5_wish.ProductDetail{
		ProductId:   easygo.NewInt64(productId),
		ProductName: easygo.NewString(wishItem.GetName()),
		Image:       easygo.NewString(wishItem.GetIcon()),
		ProductType: easygo.NewInt32(boxItem.GetStyle()),
		Desc:        easygo.NewString(wishItem.GetDesc()),
		//Material:    easygo.NewString(wishItem.GetMaterial()),
		Long:  easygo.NewInt64(wishItem.GetLength()),
		Width: easygo.NewInt64(wishItem.GetWide()),
		High:  easygo.NewInt64(wishItem.GetHigh()),
		Price: easygo.NewInt64(boxItem.GetPrice()),
	}, nil
}

// 许愿/修改愿望
func WishService(reqMsg *h5_wish.WishReq, userId int64) (*h5_wish.WishResp, error) {
	box, _ := for_game.GetWishBox(reqMsg.GetBoxId(), []int64{for_game.WISH_DOWN_STATUS, for_game.WISH_PUT_ON_STATUS, for_game.WISH_ADD_PRODUCT_STATUS})
	// 许愿的物品是否在该盲盒中
	if !easygo.Contain(box.GetItems(), reqMsg.GetProductId()) {
		logs.Error("许愿/修改愿望,opType= %d(1-许愿;2-修改愿望) 该物品不在此盲盒中,盲盒id为: %d,中间表的id为: %d", reqMsg.GetOpType(), reqMsg.GetBoxId(), reqMsg.GetProductId())
		return nil, errors.New("该物品不在此盲盒中")
	}
	switch reqMsg.GetOpType() {
	case WISH_OP_TYPE_ADD: // 1-许愿;2-修改愿望
		fun := func() {
			if redisWishPlayer := for_game.GetRedisWishPlayer(userId); redisWishPlayer != nil {
				redisWishPlayer.SetNotOneWish(true)
			}
		}
		easygo.Spawn(fun) // 修改不是首次许愿
		// 判断是否是重复许愿了
		wishData, _ := for_game.GetWishDataByStatus(userId, reqMsg.GetBoxId(), for_game.WISH_CHALLENGE_WAIT)
		//if wishData.GetId() != 0 && wishData.GetWishBoxItemId() == reqMsg.GetProductId() {
		if wishData.GetId() != 0 {
			logs.Error("已许愿过此愿望,用户id: %d,盲盒id: %d,物品id为: %d", userId, reqMsg.GetBoxId(), reqMsg.GetProductId())
			return nil, errors.New(for_game.WISH_ERR_1)
		}
		boxItem, _ := for_game.GetWishBoxItem(reqMsg.GetProductId())
		if boxItem.GetId() == 0 {
			logs.Error("许愿,根据物品中间表id查找物品失败,中间表id为: %d", reqMsg.GetProductId())
			return nil, errors.New("物品id有误")
		}
		var productUrl string
		var isMakeWish bool
		var wishItemId int64
		if boxItem != nil {
			if item := for_game.GetWishItemByIdFromDB(boxItem.GetWishItemId()); item.GetId() != 0 {
				productUrl = item.GetIcon()
			}
			isMakeWish = boxItem.GetIsWin()
			wishItemId = boxItem.GetWishItemId()
		}

		data := &share_message.PlayerWishData{
			PlayerId:      easygo.NewInt64(userId),
			WishBoxId:     easygo.NewInt64(reqMsg.GetBoxId()),
			WishBoxItemId: easygo.NewInt64(reqMsg.GetProductId()),
			Match:         easygo.NewInt32(box.GetMatch()),
			ProductUrl:    easygo.NewString(productUrl),
			IsMakeWish:    easygo.NewBool(isMakeWish),

			WishItemId: easygo.NewInt64(wishItemId),
		}
		err := for_game.AddWishData(data)
		if err != nil {
			return nil, err
		}
		// 记录一击必中许愿记录进redis
		for_game.SetDareFrequency1(0, userId, reqMsg.GetBoxId(), reqMsg.GetProductId())
	case WISH_OP_TYPE_EDIT:
		data, _ := for_game.GetWishDataByStatus(userId, reqMsg.GetBoxId(), for_game.WISH_CHALLENGE_WAIT)
		if data.GetId() == 0 {
			logs.Error("修改愿望,没有以前的许愿记录,用户id: %d,盲盒id: %d", userId, reqMsg.GetBoxId())
			return nil, errors.New("参数有误")
		}
		productId := data.GetWishBoxItemId()
		if productId == reqMsg.GetProductId() { // 取消操作
			for_game.DelWishData(data.GetId())
			// 删除一击必中许愿记录
			if err := for_game.DelDareFrequency1(userId, reqMsg.GetBoxId(), reqMsg.GetProductId()); err != nil {
				logs.Error("delDareFrequency err: %s", err.Error())
			}
		} else { // 修改
			var productUrl string
			var wishItemId int64
			if boxItem := for_game.GetWishBoxItemByIdFromDB(reqMsg.GetProductId()); boxItem != nil {
				wishItemId = boxItem.GetWishItemId()
				if item := for_game.GetWishItemByIdFromDB(wishItemId); item.GetId() != 0 {
					productUrl = item.GetIcon()
				}
			}
			err := for_game.UpWishBoxItemId(userId, reqMsg.GetBoxId(), reqMsg.GetProductId(), wishItemId, for_game.WISH_CHALLENGE_WAIT, productUrl)
			if err != nil {
				return nil, err
			}
			// 删除以前的
			if err := for_game.DelDareFrequency1(userId, reqMsg.GetBoxId(), productId); err != nil {
				logs.Error("delDareFrequency err: %s", err.Error())
			}
			// 记录一击必中许愿记录进redis
			for_game.SetDareFrequency1(0, userId, reqMsg.GetBoxId(), reqMsg.GetProductId())
		}
	}
	respMsg := &h5_wish.WishResp{
		Result: easygo.NewInt32(1),
	}
	return respMsg, nil
}

// 校验盲盒中是否有需要补货的物品
func CheckBoxStatus(box *share_message.WishBox) easygo.IMessage {
	lst, _ := for_game.GetWishBoxItemByIdsAndNumFromDB(box.GetItems())
	if len(lst) > 0 { // 长度大于0 表示有需要定时修改积极补货中的盲盒
		// 修改盲盒的状态为积极补货中
		box.Status = easygo.NewInt32(2)   //状态:0下架，1上架;2-积极补货中
		box.IsTask = easygo.NewBool(true) // 设置定时任务中.
		logs.Info("需要单次补货,盲盒的状态修改成积极补货中,boxId: %d", box.GetId())
		_ = for_game.UpsertBox(box.GetId(), box)
		var lastTaskTime int32
		var saveData []interface{}
		for _, v := range lst {
			taskTime := v.GetPerTime() * 60 // 单位秒
			if v.GetTaskTime() > 0 {        // 停服后重启导致的
				if int32(v.GetTaskTime()) > lastTaskTime {
					lastTaskTime = int32(v.GetTaskTime())
				}
			} else {
				v.TaskTime = easygo.NewInt64(taskTime)
				saveData = append(saveData, bson.M{"_id": v.GetId()}, v)
				if taskTime > lastTaskTime {
					lastTaskTime = taskTime
				}
			}
		}
		// 批量修改时间
		if len(saveData) > 0 {
			for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX_ITEM, saveData)
		}

		// 开启定时任务修改上架
		fun := func() {
			taskUpdateBoxStatus(box, lst)
		}
		logs.Info("%d 秒后执行定时任务", lastTaskTime)
		easygo.AfterFunc(time.Duration(lastTaskTime)*time.Second, fun)
		return easygo.NewFailMsg("积极补货中...")
	}
	return nil
}

// 修改上架
func taskUpdateBoxStatus(box *share_message.WishBox, lst []*share_message.WishBoxItem) {
	logs.Info("积极补货中,定时修改上架开始,boxId: %d", box.GetId())
	box.Status = easygo.NewInt32(1) //状态:0下架，1上架;2-积极补货中
	box.IsTask = easygo.NewBool(false)
	_ = for_game.UpsertBox(box.GetId(), box)
	// 修改物品的当前数量和清空定时任务时间.
	var saveData []interface{}
	for _, v := range lst {
		v.LocalNum = easygo.NewInt32(v.GetPerNum())
		v.TaskTime = easygo.NewInt64(0)
		saveData = append(saveData, bson.M{"_id": v.GetId()}, v)
	}
	// 批量修改时间
	if len(saveData) > 0 {
		for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_BOX_ITEM, saveData)
	}
}

// 备份抽奖的函数,修改砖石抽奖在这里修改
//func DoDareService(reqMsg *h5_wish.DoDareReq, userId int64, wishData *share_message.PlayerWishData) (*h5_wish.DoDareResp, error) {
//	boxId := reqMsg.GetWishBoxId()
//	wishBox, err := for_game.GetWishBox(boxId, []int64{for_game.WISH_PUT_ON_STATUS})
//	if err != nil {
//		return nil, errors.New("数据错误")
//	}
//	// 拿到兑换比例,扣钱之前确保流程能跑通使用
//	cfg := for_game.GetCurrencyCfg()
//	if cfg == nil {
//		logs.Error("获取兑换配置失败")
//		return nil, errors.New("获取兑换配置失败")
//	}
//
//	c := cfg.GetCoin()
//	d := cfg.GetDiamond()
//	if c == 0 || d == 0 {
//		logs.Error("配置表配置的兑换汇率为空")
//		return nil, errors.New("兑换汇率为空")
//	}
//	base := for_game.GetRedisWishPlayer(userId)
//	if base == nil {
//		logs.Error("获取许愿池用户对象失败,许愿池用户id为: %d", userId)
//		return nil, errors.New("用户信息有误")
//	}
//	wishBoxPrice := wishBox.GetPrice()
//	imPid := base.GetPlayerId()
//	extendLog := &share_message.GoldExtendLog{
//		PlayerId: easygo.NewInt64(imPid),
//	}
//
//	err11, _ := base.AddDiamond(0-wishBoxPrice, "抽奖", for_game.DIAMOND_TYPE_PLAYER_OUT, extendLog)
//	if err11 != nil {
//		return nil, errors.New(err11.GetReason())
//	}
//
//	IsLucky := false
//	IsOnce := false
//	randItemId := MakeWishLucky(userId, wishBox, wishData, cfg) // 挑战赛抽奖 得出中间表的id
//	if randItemId == 0 {
//		logs.Error("抽奖得到的物品id为0,需要人工排查,userId= %d,wishBox= %v,wishData= %v", userId, wishBox, wishData)
//		logs.Info("-----恢复用户的金额------")
//
//		err11, _ := base.AddDiamond(wishBoxPrice, "许愿池抽奖失败钻石返回", for_game.DIAMOND_TYPE_WISH_FAILD, extendLog)
//		if err11 != nil {
//			return nil, errors.New(err11.GetReason())
//		}
//
//		logs.Info("-----恢复用户的金额结束------")
//		return nil, errors.New("抽奖失败")
//	}
//	// 异步修改当前的库存量
//	easygo.Spawn(IncLocalNum, reqMsg.GetWishBoxId(), randItemId)
//	// 异步修改水池状态
//	fun := func() {
//		status := GetPoolStatus(wishBox.GetWishPoolId())
//		poolObj := GetPoolObj(wishBox.GetWishPoolId())
//		poolObj.SetLocalStatus(wishBox.GetWishPoolId(), int64(status))
//		if status == for_game.POOL_STATUS_BIGLOSS {
//			wishBox.Status = easygo.NewInt32(2)
//			logs.Info("水池为大亏状态,盲盒的状态修改成积极补货中,boxId: %d", wishBox.GetId())
//			_ = for_game.UpsertBox(wishBox.GetId(), wishBox)
//		}
//	}
//	easygo.Spawn(fun)
//	// 异步小助手推送
//	//easygo.Spawn(pushDareResult, randItemId)
//
//	wishBoxItem, _ := for_game.GetWishBoxItem(randItemId) // 中奖信息
//
//	wishItem, _ := for_game.QueryWishItemById(wishBoxItem.GetWishItemId())
//	dareId := userId
//	//playerPlayerInfo := for_game.GetWishPlayerInfo(dareId) //挑战者信息
//	var beDareId int64
//	var beDareNickName string
//	if randItemId == wishData.GetWishBoxItemId() {
//		IsLucky = true
//	}
//	//挑战赛
//	if reqMsg.GetDareType() == WISH_DARE {
//		Protector(wishBox, wishData, IsLucky)
//		beDareId = wishBox.GetGuardianId()
//		guardianPlayerInfo := for_game.GetWishPlayerInfo(wishBox.GetGuardianId()) //守护者信息
//		beDareNickName = guardianPlayerInfo.GetNickName()
//	}
//
//	result := false
//	if IsLucky {
//		result = true
//	}
//
//	var headIcon, nickName string
//	if p := for_game.GetRedisWishPlayer(wishData.GetPlayerId()); p != nil {
//		headIcon = p.GetHeadIcon()
//		nickName = p.GetNickName()
//	}
//	poolObj := GetPoolObj(wishBox.GetWishPoolId())
//	pool := poolObj.GetPollInfoFromRedis()
//	WishLog := &share_message.WishLog{
//		WishBoxId:         easygo.NewInt64(boxId),
//		DareId:            easygo.NewInt64(dareId),
//		DareName:          easygo.NewString(nickName),
//		BeDareId:          easygo.NewInt64(beDareId),
//		BeDareName:        easygo.NewString(beDareNickName),
//		Result:            easygo.NewBool(result),
//		ChallengeItemId:   easygo.NewInt64(wishBoxItem.GetId()),
//		DareHeadIcon:      easygo.NewString(headIcon),
//		DarePrice:         easygo.NewInt64(wishBox.GetPrice()),
//		WishItemId:        easygo.NewInt64(wishBoxItem.GetWishItemId()),
//		DareType:          easygo.NewInt32(reqMsg.GetDareType()),
//		ChallengeItemName: easygo.NewString(wishItem.GetName()),
//		PoolLocalStatus:   easygo.NewInt64(pool.GetLocalStatus()),
//		PoolIsOpenAward:   easygo.NewBool(pool.GetIsOpenAward()),
//		PoolId:            easygo.NewInt64(wishBox.GetWishPoolId()),
//		PoolIncomeValue:   easygo.NewInt64(pool.GetIncomeValue()),
//	}
//	for_game.AddWishLog(WishLog) //添加盲盒挑战记录
//
//	id := for_game.NextId(for_game.TABLE_PLAYER_WISH_ITEM)
//
//	playerWishItem := &share_message.PlayerWishItem{
//		Id:              easygo.NewInt64(id),
//		PlayerId:        easygo.NewInt64(userId),
//		ChallengeItemId: easygo.NewInt64(wishBoxItem.GetId()),
//		Status:          easygo.NewInt32(0),
//		WishBoxId:       easygo.NewInt64(wishBox.GetId()),
//		IsRead:          easygo.NewBool(false),
//		WishItemId:      easygo.NewInt64(wishBoxItem.GetWishItemId()),
//		WishItemPrice:   easygo.NewInt64(wishItem.GetPrice()),
//		DareDiamond:     easygo.NewInt64(wishBox.GetPrice()),
//		ProductName:     easygo.NewString(wishItem.GetName()),
//		WishItemDiamond: easygo.NewInt64(wishItem.GetDiamond()),
//		WishItemIcon:    easygo.NewString(wishItem.GetIcon()),
//		WishItemStyle:   easygo.NewInt32(wishBoxItem.GetStyle()),
//		BoxName:         easygo.NewString(wishBox.GetName()),
//		BoxIcon:         easygo.NewString(wishBox.GetIcon()),
//		BoxMatch:        easygo.NewInt32(wishBox.GetMatch()),
//	}
//	err = for_game.AddPlayerWishItem(playerWishItem) //添加玩家物品
//	easygo.PanicError(err)
//
//	//WishDataItemList = append(WishDataItemList, playerWishItem) // todo csv 使用
//
//	err = for_game.SetDareFrequency(userId) //挑战成功，增加挑战次数
//	easygo.PanicError(err)
//
//	frequency := for_game.GetDareFrequency(userId)
//	if frequency <= 1 && IsLucky { //第一次就中
//		IsOnce = true
//	}
//
//	respData := &h5_wish.DoDareResp{
//		IsLucky:          easygo.NewBool(IsLucky),
//		ProductId:        easygo.NewInt64(wishBoxItem.GetWishItemId()),
//		ProductName:      easygo.NewString(wishItem.GetName()),
//		Image:            easygo.NewString(wishItem.GetIcon()),
//		ProductType:      easygo.NewInt32(wishBoxItem.GetStyle()),
//		IsOnce:           easygo.NewBool(IsOnce),
//		Status:           easygo.NewInt32(wishBoxItem.GetStatus()),
//		ForecastTime:     easygo.NewInt64(wishBoxItem.GetPredictArrivalTime()),
//		PlayerWishItemId: easygo.NewInt64(id),
//		ProductPrice:     easygo.NewInt64(wishItem.GetPrice()),
//		ProductDiamond:   easygo.NewInt64(wishItem.GetDiamond()),
//	}
//
//	return respData, nil
//}

func DoDareService1(reqMsg *h5_wish.DoDareReq, userId int64, wishData *share_message.PlayerWishData, goroutineID uint64) (*h5_wish.DoDareResp, error) {
	boxId := reqMsg.GetWishBoxId()
	wishBox, err := for_game.GetWishBox(boxId, []int64{for_game.WISH_PUT_ON_STATUS})
	if err != nil {
		return nil, errors.New("数据错误")
	}
	// 拿到兑换比例,扣钱之前确保流程能跑通使用
	cfg := for_game.GetCurrencyCfg()
	if cfg == nil {
		logs.Error("goroutineID: %v,获取兑换配置失败", goroutineID)
		return nil, errors.New("获取兑换配置失败")
	}

	c := cfg.GetCoin()
	d := cfg.GetDiamond()
	if c == 0 || d == 0 {
		logs.Error("goroutineID: %v,配置表配置的兑换汇率为空", goroutineID)
		return nil, errors.New("兑换汇率为空")
	}
	base := for_game.GetRedisWishPlayer(userId)
	if base == nil {
		logs.Error("goroutineID: %v,获取许愿池用户对象失败,许愿池用户id为: %d", goroutineID, userId)
		return nil, errors.New("用户信息有误")
	}
	wishBoxPrice := wishBox.GetPrice()
	imPid := base.GetPlayerId()
	extendLog := &share_message.GoldExtendLog{
		PlayerId: easygo.NewInt64(imPid),
	}

	err11, _ := base.AddDiamond(0-wishBoxPrice, "抽奖", for_game.DIAMOND_TYPE_PLAYER_OUT, extendLog)
	if err11 != nil {
		return nil, errors.New(err11.GetReason())
	}

	IsLucky := false
	IsOnce := false
	dareStatus, randItemId := MakeWishLucky(base, userId, wishBox, wishData, goroutineID) // 挑战赛抽奖 得出中间表的id
	if randItemId == 0 {
		logs.Error("goroutineID: %v,抽奖得到的物品id为0,需要人工排查,userId= %d,wishBox= %v,wishData= %v", goroutineID, userId, wishBox, wishData)
		logs.Info("goroutineID: %v,-----恢复用户的金额------", goroutineID)

		err11, _ := base.AddDiamond(wishBoxPrice, "许愿池抽奖失败钻石返回", for_game.DIAMOND_TYPE_WISH_FAILD, extendLog)
		if err11 != nil {
			return nil, errors.New(err11.GetReason())
		}

		logs.Info("goroutineID: %v,-----恢复用户的金额结束------", goroutineID)
		return nil, errors.New("抽奖失败")
	}
	// 异步修改当前的库存量
	easygo.Spawn(IncLocalNum, reqMsg.GetWishBoxId(), randItemId)
	// 异步处理许愿池活动
	easygo.Spawn(WishActService, userId, boxId, wishBox.GetPrice(), goroutineID)
	//  修改水池状态
	status := GetPoolStatus(wishBox.GetWishPoolId(), goroutineID)
	poolObj1 := GetPoolObj(wishBox.GetWishPoolId())
	poolObj1.SetLocalStatus(wishBox.GetWishPoolId(), int64(status))
	if status == for_game.POOL_STATUS_BIGLOSS {
		wishBox.Status = easygo.NewInt32(2)
		logs.Info("goroutineID: %v,水池为大亏状态,盲盒的状态修改成积极补货中,boxId: %d", goroutineID, wishBox.GetId())
		_ = for_game.UpsertBox(wishBox.GetId(), wishBox)
	}

	wishBoxItem, _ := for_game.GetWishBoxItem(randItemId) // 中奖信息

	wishItem, _ := for_game.QueryWishItemById(wishBoxItem.GetWishItemId())
	dareId := userId
	var beDareId int64
	var beDareNickName string
	if randItemId == wishData.GetWishBoxItemId() {
		IsLucky = true
	}

	result := false
	if IsLucky {
		result = true
		// 修改许愿数据
		UpPlayerWishData := &share_message.PlayerWishData{
			Status:     easygo.NewInt32(1),
			FinishTime: easygo.NewInt64(time.Now().Unix()),
		}
		if err := for_game.UpWishDataById(wishData.GetId(), UpPlayerWishData); err != nil { //修改许愿;
			logs.Error("goroutineID: %v,修改玩家许愿数据失败,更新的内容为: %+v, err: %s", goroutineID, UpPlayerWishData, err.Error())
		}
		// 盲盒许愿数加1
		if err := for_game.SetWishBoxWishNum(boxId); err != nil {
			logs.Error("goroutineID: %v,增加盲盒winNum失败,盲盒id为: %d, err: %s", goroutineID, boxId, err.Error())
		}
	}

	//挑战赛
	if reqMsg.GetDareType() == WISH_DARE {
		Protector(wishBox, wishData, IsLucky)
		beDareId = wishBox.GetGuardianId()
		guardianPlayerInfo := for_game.GetWishPlayerInfo(wishBox.GetGuardianId()) //守护者信息
		beDareNickName = guardianPlayerInfo.GetNickName()
	}

	/*	var headIcon, nickName string
		if p := for_game.GetRedisWishPlayer(wishData.GetPlayerId()); p != nil {
			headIcon = p.GetHeadIcon()
			nickName = p.GetNickName()
		}*/
	WishLog := &share_message.WishLog{
		WishBoxId:            easygo.NewInt64(boxId),
		DareId:               easygo.NewInt64(dareId),
		DareName:             easygo.NewString(base.GetNickName()),
		BeDareId:             easygo.NewInt64(beDareId),
		BeDareName:           easygo.NewString(beDareNickName),
		Result:               easygo.NewBool(result),
		ChallengeItemId:      easygo.NewInt64(wishBoxItem.GetId()),
		DareHeadIcon:         easygo.NewString(base.GetHeadIcon()),
		DarePrice:            easygo.NewInt64(wishBox.GetPrice()),
		WishItemId:           easygo.NewInt64(wishBoxItem.GetWishItemId()),
		DareType:             easygo.NewInt32(reqMsg.GetDareType()),
		ChallengeItemName:    easygo.NewString(wishItem.GetName()),
		PoolLocalStatus:      easygo.NewInt64(dareStatus.LocalStatus),
		PoolIsOpenAward:      easygo.NewBool(dareStatus.IsOpenAward),
		AfterPoolIsOpenAward: easygo.NewBool(dareStatus.AfterIsOpenAward),
		AfterPoolLocalStatus: easygo.NewInt64(dareStatus.AfterLocalStatus),
		PoolId:               easygo.NewInt64(wishBox.GetWishPoolId()),
		PoolIncomeValue:      easygo.NewInt64(dareStatus.PoolIncomeValue),
		AfterPoolIncomeValue: easygo.NewInt64(dareStatus.AfterPoolIncomeValue),
		CreateTimeMill:       easygo.NewInt64(for_game.GetMillSecond()),
		GoroutineID:          easygo.NewInt64(goroutineID),
	}
	for_game.AddWishLog(WishLog) //添加盲盒挑战记录

	id := for_game.NextId(for_game.TABLE_PLAYER_WISH_ITEM)

	playerWishItem := &share_message.PlayerWishItem{
		Id:               easygo.NewInt64(id),
		PlayerId:         easygo.NewInt64(userId),
		ChallengeItemId:  easygo.NewInt64(wishBoxItem.GetId()),
		Status:           easygo.NewInt32(0),
		WishBoxId:        easygo.NewInt64(wishBox.GetId()),
		IsRead:           easygo.NewBool(false),
		WishItemId:       easygo.NewInt64(wishBoxItem.GetWishItemId()),
		WishItemPrice:    easygo.NewInt64(wishItem.GetPrice()),
		DareDiamond:      easygo.NewInt64(wishBox.GetPrice()),
		ProductName:      easygo.NewString(wishItem.GetName()),
		WishItemDiamond:  easygo.NewInt64(wishItem.GetDiamond()),
		WishItemIcon:     easygo.NewString(wishItem.GetIcon()),
		WishItemStyle:    easygo.NewInt32(wishBoxItem.GetStyle()),
		BoxName:          easygo.NewString(wishBox.GetName()),
		BoxIcon:          easygo.NewString(wishBox.GetIcon()),
		BoxMatch:         easygo.NewInt32(wishBox.GetMatch()),
		AfterIsOpenAward: easygo.NewBool(dareStatus.AfterIsOpenAward),  // 抽奖后是否开启放奖 true 放奖 false 关闭
		AfterLocalStatus: easygo.NewInt64(dareStatus.AfterLocalStatus), //  抽奖后当前水池的状态;1-大亏;2-小亏;3-普通;4-大赢;5-小盈
		IsOpenAward:      easygo.NewBool(dareStatus.IsOpenAward),       // 是否开启放奖 true 放奖 false 关闭
		LocalStatus:      easygo.NewInt64(dareStatus.LocalStatus),      // 当前水池的状态;1-大亏;2-小亏;3-普通;4-大赢;5-小盈
		GoroutineID:      easygo.NewInt64(goroutineID),
		CreateTimeMill:   easygo.NewInt64(for_game.GetMillSecond()),
	}
	err = for_game.AddPlayerWishItem(playerWishItem) //添加玩家物品
	easygo.PanicError(err)

	/*
		// 异步小助手推送 todo 策划突然取消推送功能
		pushReq := &h5_wish.WishNoticeAssistantReq{
			PlayerId:     easygo.NewInt64(imPid),
			WishPlayerId: easygo.NewInt64(userId),
			ProductName:  easygo.NewString(wishItem.GetName()),
		}
		easygo.Spawn(pushDareResult, pushReq)
	*/
	err = for_game.SetDareFrequency(userId) //挑战成功，增加挑战次数 用户-盲盒-许愿的物品id做唯一key
	easygo.PanicError(err)
	frequency := for_game.SetDareFrequency1(1, userId, boxId, wishData.GetWishBoxItemId()) // 一击必中的次数.
	//frequency := for_game.GetDareFrequency(userId)
	if frequency <= 1 && IsLucky { //第一次就中
		IsOnce = true
	}

	respData := &h5_wish.DoDareResp{
		IsLucky:          easygo.NewBool(IsLucky),
		ProductId:        easygo.NewInt64(wishBoxItem.GetWishItemId()),
		ProductName:      easygo.NewString(wishItem.GetName()),
		Image:            easygo.NewString(wishItem.GetIcon()),
		ProductType:      easygo.NewInt32(wishBoxItem.GetStyle()),
		IsOnce:           easygo.NewBool(IsOnce),
		Status:           easygo.NewInt32(wishBoxItem.GetStatus()),
		ForecastTime:     easygo.NewInt64(wishBoxItem.GetPredictArrivalTime()),
		PlayerWishItemId: easygo.NewInt64(id),
		ProductPrice:     easygo.NewInt64(wishItem.GetPrice()),
		ProductDiamond:   easygo.NewInt64(wishItem.GetDiamond()),
	}

	return respData, nil
}

// 异步修改当前的库存量
func IncLocalNum(boxId, boxItemId int64) {
	wishBox, err := for_game.GetWishBox(boxId, []int64{for_game.WISH_PUT_ON_STATUS})
	if err != nil {
		return
	}

	result, err := for_game.UpBoxItemLocalNum(boxItemId)
	if err != nil {
		return
	}
	if result > 0 {
		return
	}

	CheckBoxStatus(wishBox)
}

// 平台回收
//func PlantBack(ids []int64, playerId int64) {
//	logs.Info("==========回收定时任务执行===========id: %v,playerId: %v", ids, playerId)
//	err := RecycleGoods(ids, playerId, RECYCLE_BY_SYSTEM) // 系统回收.
//	easygo.PanicError(err)
//}

// 推送中奖到柠檬助手
func pushDareResult(pushReq *h5_wish.WishNoticeAssistantReq) {
	logs.Info("===========推送中奖到柠檬助手=============pushReq: %v", pushReq)
	_, err1 := SendMsgToIdelServer(for_game.SERVER_TYPE_HALL, "RpcWishNoticeAssistant", pushReq, pushReq.GetWishPlayerId())
	if err1 != nil {
		logs.Error("pushDareResult-err:", err1)
	}
}

func Protector(WishBox *share_message.WishBox, playerWishData *share_message.PlayerWishData, IsLucky bool) {

	playerPlayerInfo := for_game.GetWishPlayerInfo(playerWishData.GetPlayerId()) //挑战者信息
	var dayDiamondTop, addDiamondNum int64
	if cfg := for_game.GetWishGuardianCfg(); cfg != nil {
		dayDiamondTop = cfg.GetDayDiamondTop()
		addDiamondNum = cfg.GetOnceDiamondRebate()
	}

	addNum, err := for_game.UpOccupiedCoin(WishBox.GetId(), WishBox.GetGuardianId(), int32(dayDiamondTop), int32(addDiamondNum)) //增加金币
	if err == nil {
		gid := WishBox.GetGuardianId()
		if guardian := for_game.GetRedisWishPlayer(gid); guardian != nil {
			extendLog := &share_message.GoldExtendLog{
				PlayerId: easygo.NewInt64(guardian.GetPlayerId()),
			}
			err, _ := guardian.AddDiamond(int64(addNum), "守护者收益", for_game.DIAMOND_TYPE_WISH_GUARDIAN_IN, extendLog)
			if err != nil {
				logs.Error("给守护者添加钻石失败,用户id为: %d,盲盒id为: %d", WishBox.GetGuardianId(), WishBox.GetId())
			}

			//设置当天被挑战者(守护者)获取的金币总数
			dataLog := &share_message.WishGuardianDiamondLog{
				Id:         easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_GUARDIAN_DIAMOND_LOG)),
				BoxId:      easygo.NewInt64(WishBox.GetId()),
				PlayerId:   easygo.NewInt64(gid),
				HeadIcon:   easygo.NewString(guardian.GetHeadIcon()),
				NickName:   easygo.NewString(guardian.GetNickName()),
				CoinNum:    easygo.NewInt64(addNum),
				CreateTime: easygo.NewInt64(for_game.GetMillSecond()),
			}
			if err := for_game.AddWishGuardianDiamondLog(dataLog); err != nil {
				logs.Error("后台添加守护者.添加榜单信息失败,err: %s", err.Error())
			}
		}
	} else {
		logs.Error("------------->err", err.Error())
	}
	/*

		// 给发起挑战的人添加硬币 // todo 策划说先屏蔽
		if redisWishPlayer := for_game.GetRedisWishPlayer(playerWishData.GetPlayerId()); redisWishPlayer != nil {
			extendLog := &share_message.GoldExtendLog{
				PlayerId: easygo.NewInt64(playerWishData.GetPlayerId()),
			}

			err, _ := redisWishPlayer.AddDiamond(COIN, "抽奖返利", for_game.DIAMOND_TYPE_WISH_DARE_BACK, extendLog)
			if err != nil {
				logs.Error("给发起挑战的人添加钻石失败,用户id为: %d,盲盒id为: %d", WishBox.GetGuardianId(), WishBox.GetId())
			}
		}
	*/

	if IsLucky { //挑战成功
		UpWishBoxData := &share_message.WishBox{
			GuardianId:       easygo.NewInt64(playerWishData.GetPlayerId()),
			GuardianOverTime: easygo.NewInt64(time.Now().Unix() + for_game.AFTER_TIME_GUARDIAN),
			IsGuardian:       easygo.NewBool(true),
		}
		if err := for_game.UpWishBox(playerWishData.GetWishBoxId(), UpWishBoxData); err != nil { //挑战成功，设置为守护者
			logs.Error("挑战成功,设置成守护者失败,更新的内容为:UpWishBoxData: %+v, err: %s", UpWishBoxData, err.Error())
		}
		if err := for_game.UpOccupied(playerWishData.GetWishBoxId(), WishBox.GetGuardianId()); err != nil { //设置挑战者结束时间
			logs.Error("挑战成功,设置挑战者结束时间,  err: %s", err.Error())
		}
		playerOccupied := &share_message.WishOccupied{
			WishBoxId: easygo.NewInt64(playerWishData.GetWishBoxId()),
			NickName:  easygo.NewString(playerPlayerInfo.GetNickName()),
			HeadUrl:   easygo.NewString(playerPlayerInfo.GetHeadUrl()),
			PlayerId:  easygo.NewInt64(playerWishData.GetPlayerId()),
			Status:    easygo.NewInt32(1),
		}

		if err := for_game.AddOccupied(playerOccupied); err != nil { //挑战者时长记录
			logs.Error("挑战成功设置占领时长表失败,playerOccupied: %+v, err: %s", playerOccupied, err.Error())
		}

		dataLog := &share_message.WishGuardianDiamondLog{
			Id:         easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_GUARDIAN_DIAMOND_LOG)),
			BoxId:      easygo.NewInt64(WishBox.GetId()),
			PlayerId:   easygo.NewInt64(playerWishData.GetPlayerId()),
			HeadIcon:   easygo.NewString(playerPlayerInfo.GetHeadUrl()),
			NickName:   easygo.NewString(playerPlayerInfo.GetNickName()),
			WishNum:    easygo.NewInt64(1),
			CreateTime: easygo.NewInt64(for_game.GetMillSecond()),
		}
		if err := for_game.AddWishGuardianDiamondLog(dataLog); err != nil {
			logs.Error("添加挑战者挑战的成功次数,err: %s", err.Error())
		}

	} else { // 如果当前的守护者是自己的话,清空自己的守护者位置
		if WishBox.GetGuardianId() == playerWishData.GetPlayerId() {
			UpdateBoxGuardian(WishBox)
		}
	}

}

//func Protector(WishBox *share_message.WishBox, playerWishData *share_message.PlayerWishData, IsLucky bool) {
//
//	playerPlayerInfo := for_game.GetWishPlayerInfo(playerWishData.GetPlayerId()) //挑战者信息
//	var dayDiamondTop, addDiamondNum int64
//	if cfg := for_game.GetWishGuardianCfg(); cfg != nil {
//		dayDiamondTop = cfg.GetDayDiamondTop()
//		addDiamondNum = cfg.GetOnceDiamondRebate()
//	}
//
//	err := for_game.UpOccupiedCoin(WishBox.GetId(), WishBox.GetGuardianId(), int32(dayDiamondTop), int32(addDiamondNum)) //增加金币
//	if err == nil {
//		gid := WishBox.GetGuardianId()
//		if guardian := for_game.GetRedisWishPlayer(gid); guardian != nil {
//			extendLog := &share_message.GoldExtendLog{
//				PlayerId: easygo.NewInt64(guardian.GetPlayerId()),
//			}
//			err, _ := guardian.AddDiamond(addDiamondNum, "守护者收益", for_game.DIAMOND_TYPE_WISH_GUARDIAN_IN, extendLog)
//			if err != nil {
//				logs.Error("给守护者添加钻石失败,用户id为: %d,盲盒id为: %d", WishBox.GetGuardianId(), WishBox.GetId())
//			}
//		}
//
//	} else {
//		logs.Error("------------->err", err.Error())
//	}
//	/*
//
//		// 给发起挑战的人添加硬币 // todo 策划说先屏蔽
//		if redisWishPlayer := for_game.GetRedisWishPlayer(playerWishData.GetPlayerId()); redisWishPlayer != nil {
//			extendLog := &share_message.GoldExtendLog{
//				PlayerId: easygo.NewInt64(playerWishData.GetPlayerId()),
//			}
//
//			err, _ := redisWishPlayer.AddDiamond(COIN, "抽奖返利", for_game.DIAMOND_TYPE_WISH_DARE_BACK, extendLog)
//			if err != nil {
//				logs.Error("给发起挑战的人添加钻石失败,用户id为: %d,盲盒id为: %d", WishBox.GetGuardianId(), WishBox.GetId())
//			}
//		}
//	*/
//	//添加榜单数据
//	wishTopLog, err := for_game.GetWishTopLog(playerWishData.GetPlayerId())
//	if err != nil && err.Error() == mgo.ErrNotFound.Error() {
//		wishTopLog = &share_message.WishTopLog{
//			PlayerId: easygo.NewInt64(playerWishData.GetPlayerId()),
//			HeadIcon: easygo.NewString(playerPlayerInfo.GetHeadUrl()),
//			NickName: easygo.NewString(playerPlayerInfo.GetNickName()),
//		}
//		for_game.AddWishTopLog(playerWishData.GetPlayerId(), wishTopLog)
//	}
//	if IsLucky { //挑战成功
//		UpPlayerWishData := &share_message.PlayerWishData{
//			Status:     easygo.NewInt32(1),
//			FinishTime: easygo.NewInt64(time.Now().Unix()),
//		}
//		if err := for_game.UpWishDataById(playerWishData.GetId(), UpPlayerWishData); err != nil { //修改许愿;
//			logs.Error("修改玩家许愿数据失败,更新的内容为: %+v, err: %s", UpPlayerWishData, err.Error())
//		}
//		UpWishBoxData := &share_message.WishBox{
//			GuardianId:       easygo.NewInt64(playerWishData.GetPlayerId()),
//			GuardianOverTime: easygo.NewInt64(time.Now().Unix() + for_game.AFTER_TIME_GUARDIAN),
//			IsGuardian:       easygo.NewBool(true),
//		}
//		if err := for_game.UpWishBox(playerWishData.GetWishBoxId(), UpWishBoxData, 1); err != nil { //挑战成功，设置为守护者
//			logs.Error("挑战成功,设置成守护者失败,更新的内容为:UpWishBoxData: %+v, err: %s", UpWishBoxData, err.Error())
//		}
//		if err := for_game.UpOccupied(playerWishData.GetWishBoxId(), WishBox.GetGuardianId()); err != nil { //设置挑战者结束时间
//			logs.Error("挑战成功,设置挑战者结束时间,  err: %s", err.Error())
//		}
//		playerOccupied := &share_message.WishOccupied{
//			WishBoxId: easygo.NewInt64(playerWishData.GetWishBoxId()),
//			NickName:  easygo.NewString(playerPlayerInfo.GetNickName()),
//			HeadUrl:   easygo.NewString(playerPlayerInfo.GetHeadUrl()),
//			PlayerId:  easygo.NewInt64(playerWishData.GetPlayerId()),
//			Status:    easygo.NewInt32(1),
//		}
//
//		if err := for_game.AddOccupied(playerOccupied); err != nil { //挑战者时长记录
//			logs.Error("挑战成功设置占领时长表失败,playerOccupied: %+v, err: %s", playerOccupied, err.Error())
//		}
//		if err := for_game.SetWishTopLogWishNum(playerWishData.GetPlayerId()); err != nil { // 设置当天挑战者挑战的成功次数
//			logs.Error("挑战成功 设置当天挑战者挑战的成功次数  err: %s", err.Error())
//		}
//	} else { // 如果当前的守护者是自己的话,清空自己的守护者位置
//		if WishBox.GetGuardianId() == playerWishData.GetPlayerId() {
//			UpdateBoxGuardian(WishBox)
//		}
//	}
//
//	for_game.SetWishTopLogCoin(WishBox.GetGuardianId()) //设置当天被挑战者(守护者)获取的金币总数
//}

//  处理盲盒守护者信息
func UpdateBoxGuardian(box *share_message.WishBox) {
	if box.GetId() == 0 {
		return
	}
	if !box.GetIsGuardian() {
		return
	}
	guardianId := box.GetGuardianId()
	box.IsGuardian = easygo.NewBool(false)
	box.GuardianOverTime = easygo.NewInt64(0)
	box.GuardianId = easygo.NewInt64(0)
	if err := for_game.UpWishBox(box.GetId(), box); err != nil {
		logs.Error("定时修改盲盒的守护者信息失败,err: %s", err.Error())
	}

	// 修改守护时间表信息
	if err := for_game.UpOccupied(box.GetId(), guardianId); err != nil {
		logs.Error("修改守护时间表信息失败,boxId: %d, err: %s", box.GetId(), err.Error())
	}
}

func MakeWishLucky(pBase *for_game.RedisWishPlayerObj, pid int64, wishBox *share_message.WishBox, wishData *share_message.PlayerWishData, goroutineID uint64) (*DareStatus, int64) {
	// 判断用户是否是首次抽奖
	/**
	许愿必中
	1.非挑战赛
	2.用户_盲盒id
	*/

	//wp := for_game.GetWishPlayerByPid(pid)
	makeWish := for_game.GetPlayerMakeWish(pid, wishBox.GetId())
	b := makeWish == 0

	if b && wishData.GetIsMakeWish() && wishBox.GetMatch() == 0 { // 对应的盲盒许愿必中,许愿必中,match:0 普通赛
		return IsMakeWish(pBase, pid, wishBox, wishData, goroutineID)
	}
	return RandWishBoxItem(pBase, pid, wishBox, goroutineID)
}

/**
1.找到许愿的物品
2.扣除水池的价格
3.返回物品id
*/
// 许愿必中
func IsMakeWish(pBase *for_game.RedisWishPlayerObj, pid int64, wishBox *share_message.WishBox, wishData *share_message.PlayerWishData, goroutineID uint64) (*DareStatus, int64) {
	logs.Info("许愿必中开始=======================")
	st := for_game.GetMillSecond()
	boxId := wishBox.GetId()
	poolId := wishBox.GetWishPoolId()
	dStatus := &DareStatus{}
	if boxId == 0 || poolId == 0 {
		logs.Error("goroutineID %v,盲盒id或者水池id为空,boxId: %d,poolId: %d", goroutineID, boxId, poolId)
		return dStatus, 0
	}
	poolObj := GetPoolObj(poolId)
	// 判断金额是否超过抽水的阀值
	pool := poolObj.GetPollInfoFromRedis()
	if pool == nil {
		logs.Error("goroutineID %v,从redis中获取水池数据失败,用户id为: %d,水池id为: %d", goroutineID, pid, poolId)
		return dStatus, 0
	}
	beforePool := poolObj.GetPollInfoFromRedis()
	//poolObj.Mutex1.Lock()
	//defer poolObj.Mutex1.Unlock()
	//poolObj := GePoolObj()
	//  inc 挑战的金额,返回当前的水池的余额,判断是否符合抽水状态,符合就立马抽水,记录抽水日志.返回的值根据当前水池余额得到水池状态
	boxPrice := wishBox.GetPrice() // 砖石
	var pollPrice, recycle, commission int64
	var err error
	whiteList := for_game.GetWishWhiteList() // 白名单列表
	whiteIds := make([]int64, 0)
	for _, v := range whiteList {
		whiteIds = append(whiteIds, v.GetId())
	}
	// 不是白名单并且不是运营号
	if !util.Int64InSlice(pid, whiteIds) {
		//if pBase.GetTypes() != for_game.WISH_PLAYER_BASE_TYPES_2  {
		//logs.Info("盲盒的价格 %d 个硬币兑换了 %d 分人民币", boxCoin, boxPrice)
		pollPrice, err = poolObj.UpdatePollPriceToRedis(poolId, boxPrice, "IncomeValue")
		if err != nil {
			logs.Error("goroutineID %v,更改水池价格失败,用户id为: %d,盲盒id为: %d, 水池id为: %d,err: %s", goroutineID, pid, boxId, poolId, err.Error())
			return dStatus, 0
		}
		pool = poolObj.GetPollInfoFromRedis()
		// 记录水池的流水
		wishPoolLog := &share_message.WishPoolLog{
			BoxId:            easygo.NewInt64(boxId),
			PoolId:           easygo.NewInt64(poolId),
			PlayerId:         easygo.NewInt64(pid),
			BeforeValue:      easygo.NewInt64(pool.GetIncomeValue() - boxPrice),
			AfterValue:       easygo.NewInt64(pool.GetIncomeValue()),
			Value:            easygo.NewInt64(boxPrice),
			Type:             easygo.NewInt64(for_game.WISH_POOL_LOG_TYPE_2), //2-挑战增加
			LocalStatus:      easygo.NewInt64(beforePool.GetLocalStatus()),
			IsOpenAward:      easygo.NewBool(beforePool.GetIsOpenAward()),
			AfterLocalStatus: easygo.NewInt64(pool.GetLocalStatus()),
			AfterIsOpenAward: easygo.NewBool(pool.GetIsOpenAward()),
			CreateTimeMill:   easygo.NewInt64(for_game.GetMillSecond()),
			GoroutineID:      easygo.NewInt64(goroutineID),
		}
		easygo.Spawn(for_game.AddWishPoolLog, wishPoolLog)

		recycle = pool.GetRecycle() // 回收阀值(存库)
		logs.Info("goroutineID %v,投钱进水池后的金额为:%d,用户id为:%d, 抽水的阀值为: %d", goroutineID, pollPrice, pid, recycle)
		commission = pool.GetCommission() // 抽水金额

		pool = poolObj.GetPollInfoFromRedis()
		pollPrice = pool.GetIncomeValue()
		if pollPrice >= recycle { // 符合就立马抽水
			poolObj.Mutex1.Lock()
			pool = poolObj.GetPollInfoFromRedis()
			pollPrice = pool.GetIncomeValue()
			if pollPrice >= recycle {
				p, err := pump(poolId, commission, pid, boxId, recycle, beforePool.GetLocalStatus(), beforePool.GetIsOpenAward(), poolObj, pool, goroutineID)
				if err != nil {
					logs.Error("goroutineID: %v,抽水公共方法失败了,err: %s", goroutineID, err.Error())
					return dStatus, 0
				}
				pollPrice = p

			}
			poolObj.Mutex1.Unlock()
		}
	} else {
		logs.Info("goroutineID: %v,此用户为白名单用户,不参与水池水量变化,id为: %d", goroutineID, pid)
	}

	// 从数据库中查找
	item := for_game.GetWishBoxItemByIdFromDB(wishData.GetWishBoxItemId())
	if item == nil {
		logs.Error("goroutineID %v,获取wishBoxItem 为nil,用户id为: %d,_id为: %d", goroutineID, pid, wishData.GetWishBoxItemId())
		return dStatus, 0
	}
	// 修改不是第一次抽奖了
	if redisWishPlayer := for_game.GetRedisWishPlayer(pid); redisWishPlayer != nil {
		redisWishPlayer.SetNoOne(true)
	}

	for_game.SetPlayerMakeWish(1, pid, boxId)

	//物品价格兑换砖石后扣水池余额
	//diamond, err := MoneyToDiamond(item.GetPrice(), cfg)
	//if err != nil || diamond == 0 {
	//	logs.Error("人民币兑换砖石出错,尽快人工介入排查,diamond :%v,cfg:%v", diamond, cfg)
	//}
	diamond := item.GetDiamond()
	logs.Info("goroutineID %v,物品id(中间表的): %d 的钻石价格为: %d", goroutineID, wishData.GetWishBoxItemId(), diamond)
	beforePool = poolObj.GetPollInfoFromRedis()
	// todo 白名单 水池减水位的都要检查
	var wishPoolLog1BeforeValue int64
	if !util.Int64InSlice(pid, whiteIds) {
		//if pBase.GetTypes() != for_game.WISH_PLAYER_BASE_TYPES_2 {
		// 从水池中减去物品的价格
		//pollPrice, err = poolObj.UpdatePollPriceToRedis(poolId, 0-item.GetPrice(), "IncomeValue")
		pollPrice, err = poolObj.UpdatePollPriceToRedis(poolId, 0-diamond, "IncomeValue")
		if err != nil {
			logs.Error("goroutineID %v,发奖后更改水池价格失败,用户id为: %d,盲盒id为: %d, 水池id为: %d,err: %s", goroutineID, pid, boxId, poolId, err.Error())
			return dStatus, 0
		}
		pool = poolObj.GetPollInfoFromRedis()
		// 得出物品后 记录水池的流水
		wishPoolLog1 := &share_message.WishPoolLog{
			BoxId:            easygo.NewInt64(boxId),
			PoolId:           easygo.NewInt64(poolId),
			PlayerId:         easygo.NewInt64(pid),
			BeforeValue:      easygo.NewInt64(pollPrice + diamond),
			AfterValue:       easygo.NewInt64(pollPrice),
			Value:            easygo.NewInt64(0 - diamond),
			Type:             easygo.NewInt64(for_game.WISH_POOL_LOG_TYPE_3), // 3-得到物品后扣除
			IsOpenAward:      easygo.NewBool(beforePool.GetIsOpenAward()),
			LocalStatus:      easygo.NewInt64(beforePool.GetLocalStatus()),
			AfterIsOpenAward: easygo.NewBool(pool.GetIsOpenAward()),
			AfterLocalStatus: easygo.NewInt64(pool.GetLocalStatus()),
			CreateTimeMill:   easygo.NewInt64(for_game.GetMillSecond()),
			GoroutineID:      easygo.NewInt64(goroutineID),
		}
		wishPoolLog1BeforeValue = wishPoolLog1.GetBeforeValue()
		easygo.Spawn(for_game.AddWishPoolLog, wishPoolLog1)
	} else {
		logs.Info("goroutineID: %v,此用户为白名单用户,不参与水池水量变化,id为: %d", goroutineID, pid)
		wishPoolLog1BeforeValue = beforePool.GetIncomeValue()
		pool = poolObj.GetPollInfoFromRedis()
		pollPrice = pool.GetIncomeValue()
	}

	dStatus.LocalStatus = beforePool.GetLocalStatus()
	dStatus.IsOpenAward = beforePool.GetIsOpenAward()
	dStatus.AfterLocalStatus = pool.GetLocalStatus()
	dStatus.AfterIsOpenAward = pool.GetIsOpenAward()
	dStatus.PoolIncomeValue = wishPoolLog1BeforeValue
	dStatus.AfterPoolIncomeValue = pollPrice
	//   修改是是否需要附加权重.
	isOpen := checkIsOpenAward(pollPrice, pool.GetStartAward(), pool.GetCloseAward(), pool.GetIsOpenAward())
	poolObj.SetIsOpenAward(poolId, isOpen)
	if !util.Int64InSlice(pid, whiteIds) {
		//if pBase.GetTypes() != for_game.WISH_PLAYER_BASE_TYPES_2 {
		// 再次校验抽水条件.
		if pool.GetIncomeValue() >= recycle {
			poolObj.Mutex1.Lock()
			pool = poolObj.GetPollInfoFromRedis()
			if pool.GetIncomeValue() >= recycle {
				if _, err := pump(poolId, commission, pid, boxId, recycle, beforePool.GetLocalStatus(), beforePool.GetIsOpenAward(), poolObj, pool, goroutineID); err != nil {
					logs.Error("goroutineID: %v,再次校验抽水,抽水公共方法失败了,err: %s", goroutineID, err.Error())
				}
			}
			poolObj.Mutex1.Unlock()
		}
	} else {
		logs.Info("goroutineID: %v,此用户为白名单用户,不参与水池水量变化,id为: %d", goroutineID, pid)
	}

	et := for_game.GetMillSecond()
	logs.Debug("goroutineID %v,抽奖耗时为: %v", goroutineID, et-st)
	return dStatus, wishData.GetWishBoxItemId()
}

// 硬币换算人民币, 返回人民币,单位:分
func CoinToMoney(coin int64) (int64, error) {
	// 从配置表中获取兑换比例
	cfg := for_game.GetCurrencyCfg()
	if cfg == nil {
		return 0, errors.New("兑换失败")
	}
	money := cfg.GetMoney() // 单位:元
	c := cfg.GetCoin()
	if money == 0 || c == 0 {
		logs.Error("配置表配置的兑换汇率为空")
		return 0, errors.New("兑换汇率为空")
	}
	price := easygo.Decimal(easygo.AtoFloat64(easygo.AnytoA((money*100/c)*int32(coin))), 0)
	result := easygo.AtoInt64(easygo.AnytoA(price))
	return result, nil
}

// 人民币兑换砖石 人民币单位为分.
func MoneyToDiamond(rmb int64, cfg *share_message.WishCurrencyConversionCfg) (int64, error) {
	if rmb <= 0 {
		logs.Error("需要兑换的金额 rmb 为空")
		return 0, errors.New("需要兑换的金额为空")
	}
	//cfg := for_game.GetCurrencyCfg()
	if cfg == nil {
		logs.Error("获取兑换配置失败")
		return 0, errors.New("兑换失败")
	}
	//money := cfg.GetMoney() // 单位:元
	c := cfg.GetCoin()
	d := cfg.GetDiamond()
	if c == 0 || d == 0 {
		logs.Error("配置表配置的兑换汇率为空")
		return 0, errors.New("兑换汇率为空")
	}

	// todo 获取人民币与硬币兑换的比例,现在先写死=====
	m1 := 1 * 100 // 分
	c1 := 6       // 1元 6个硬币
	/**
	100  6
	2000  ?
	*/
	rmbFloat := easygo.AtoFloat64(easygo.AnytoA(rmb))  // 人民币float
	moneyFloat := easygo.AtoFloat64(easygo.AnytoA(m1)) // 配置表 人民币float
	coin1Float := easygo.AtoFloat64(easygo.AnytoA(c1)) // 硬币 float
	coinNum := easygo.Decimal((rmbFloat*coin1Float)/moneyFloat, 0)
	logs.Info("%v 分人民币兑换了%v 个硬币", rmb, coinNum)
	// todo 获取人民币与硬币兑换的比例=====

	coinFloat := easygo.AtoFloat64(easygo.AnytoA(c))    // 硬币 float
	diamondFloat := easygo.AtoFloat64(easygo.AnytoA(d)) // 砖石 float

	diamondNum := easygo.Decimal((coinNum*diamondFloat)/coinFloat, 0)
	logs.Info(" %v 个硬币 兑换了 %v 个砖石", coinNum, diamondNum)
	return int64(diamondNum), nil
}

//随机从盲盒中获取一个物品
/**
1 inc 挑战的金额,返回当前的水池的余额,判断是否符合抽水状态,符合就立马抽水,记录抽水日志.返回的值根据当前水池余额得到水池状态
2 判断pool 获取附加权重的状态,根据当前的水量是否达到了附加权重的的状态.
3 通过权重获取盲盒物品
4 通过盲盒物品分2步走:
	<1>普通物品就直接结束发货
	<2>如果是贵重物品,先上锁取水池的值(锁住当前水池不许更改.),能不能开大奖(当前盈利是否都大于该大奖物品的价格),如果属于,就开大奖,否则随机开小奖,更新水池放奖后的水量,解锁
更新水池的金额
*/

func RandWishBoxItem(pBase *for_game.RedisWishPlayerObj, pid int64, wishBox *share_message.WishBox, goroutineID uint64) (*DareStatus, int64) {
	boxId := wishBox.GetId()
	poolId := wishBox.GetWishPoolId()
	dStatus := &DareStatus{}
	if boxId == 0 || poolId == 0 {
		logs.Error("goroutineID: %v,抽奖失败原因,盲盒id或者水池id为空,boxId: %d,poolId: %d", goroutineID, boxId, poolId)
		return dStatus, 0
	}
	poolObj := GetPoolObj(poolId)
	//poolObj.Mutex1.Lock()
	//defer poolObj.Mutex1.Unlock()
	//poolObj.Mutex2.Lock()
	//defer poolObj.Mutex2.Unlock()
	//  inc 挑战的金额,返回当前的水池的余额,判断是否符合抽水状态,符合就立马抽水,记录抽水日志.返回的值根据当前水池余额得到水池状态
	var pool *share_message.WishPool
	var pollPrice, recycle, commission int64
	var err error
	whiteList := for_game.GetWishWhiteList() // 白名单列表
	whiteIds := make([]int64, 0)
	for _, v := range whiteList {
		whiteIds = append(whiteIds, v.GetId())
	}
	// 不是白名单并且不是运营号
	if !util.Int64InSlice(pid, whiteIds) {
		//if pBase.GetTypes() != for_game.WISH_PLAYER_BASE_TYPES_2 {
		boxPrice := wishBox.GetPrice()
		beforePool := poolObj.GetPollInfoFromRedis()
		pollPrice, err = poolObj.UpdatePollPriceToRedis(poolId, boxPrice, "IncomeValue")
		if err != nil {
			logs.Error("goroutineID: %v,抽奖失败原因,更改水池价格失败,用户id为: %d,盲盒id为: %d, 水池id为: %d,err: %s", goroutineID, pid, boxId, poolId, err.Error())
			return dStatus, 0
		}
		pool = poolObj.GetPollInfoFromRedis()
		if pool == nil {
			logs.Error("goroutineID :%v,抽奖失败原因,从redis中获取水池数据失败,用户id为: %d,水池id为: %d", goroutineID, pid, poolId)
			return dStatus, 0
		}

		// 记录水池的流水
		wishPoolLog := &share_message.WishPoolLog{
			BoxId:            easygo.NewInt64(boxId),
			PoolId:           easygo.NewInt64(poolId),
			PlayerId:         easygo.NewInt64(pid),
			BeforeValue:      easygo.NewInt64(pollPrice - boxPrice),
			AfterValue:       easygo.NewInt64(pollPrice),
			Value:            easygo.NewInt64(boxPrice),
			Type:             easygo.NewInt64(for_game.WISH_POOL_LOG_TYPE_2), //2-挑战增加
			LocalStatus:      easygo.NewInt64(beforePool.GetLocalStatus()),
			IsOpenAward:      easygo.NewBool(beforePool.GetIsOpenAward()),
			AfterLocalStatus: easygo.NewInt64(pool.GetLocalStatus()),
			AfterIsOpenAward: easygo.NewBool(pool.GetIsOpenAward()),
			CreateTimeMill:   easygo.NewInt64(for_game.GetMillSecond()),
			GoroutineID:      easygo.NewInt64(goroutineID),
		}
		easygo.Spawn(for_game.AddWishPoolLog, wishPoolLog)

		// 判断金额是否超过抽水的阀值
		recycle = pool.GetRecycle() // 回收阀值(存库)
		logs.Info("goroutineID: %v,投钱进水池后的金额为:%d,用户id为:%d, 抽水的阀值为: %d", goroutineID, pollPrice, pid, recycle)
		commission = pool.GetCommission() // 抽水金额
		pool = poolObj.GetPollInfoFromRedis()
		pollPrice = pool.GetIncomeValue()
		if pollPrice >= recycle { // 符合就立马抽水
			poolObj.Mutex1.Lock()
			pool = poolObj.GetPollInfoFromRedis()
			pollPrice = pool.GetIncomeValue()
			if pollPrice >= recycle {
				p, err := pump(poolId, commission, pid, boxId, recycle, beforePool.GetLocalStatus(), beforePool.GetIsOpenAward(), poolObj, pool, goroutineID)
				if err != nil {
					logs.Error("goroutineID: %v,抽奖失败原因,抽水公共方法失败了,err: %s", goroutineID, err.Error())
					return dStatus, 0
				}
				pollPrice = p
			}
			poolObj.Mutex1.Unlock()
		}
	} else {
		logs.Info("goroutineID: %v,此用户为白名单用户,不参与水池水量变化,id为: %d", goroutineID, pid)
	}

	pool = poolObj.GetPollInfoFromRedis()
	pollPrice = pool.GetIncomeValue()
	logs.Info("goroutineID: %v,计算水池状态时当前水池水量为: %d,盲盒id为: %d,水池id为: %d", goroutineID, pollPrice, boxId, poolId)
	f := easygo.Decimal(easygo.AtoFloat64(easygo.AnytoA(pollPrice))/easygo.AtoFloat64(easygo.AnytoA(pool.GetPoolLimit())), 2)
	poolStatusLimitFloat64 := f * 100
	poolStatusLimit := int64(easygo.Decimal(poolStatusLimitFloat64, 0))
	logs.Info("goroutineID: %v,盲盒id为: %d,水池id为:%d,水池的金额为:%d, 计算状态前的百分数: %d", goroutineID, boxId, poolId, pollPrice, poolStatusLimit)
	// 判断是在哪个状态
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
	if bigLoss != nil && (poolStatusLimit <= bigLoss.GetShowMinValue()) {
		poolStatus = for_game.POOL_STATUS_BIGLOSS
	}
	logs.Info("goroutineID:%v,1-大亏;2-小亏;3-普通;4-大盈;5-小盈,水池的状态为: %v", goroutineID, poolStatus)
	// 判断是否需要附加权重
	isOpenAward := pool.GetIsOpenAward()
	// 根据盲盒状态找出物品
	wishBoxItems := for_game.GetProductByPoolStatus(boxId, poolStatus, isOpenAward) // 中间表的
	if len(wishBoxItems) == 0 {
		logs.Error("goroutineID: %v,抽奖失败原因,没有当前状态的物品pid: %d,boxId: %d, poolStatus: %d", goroutineID, pid, boxId, poolStatus)
		return dStatus, 0
	}

	// 权重的map[权重]中间表物品的id
	//sortWeightMap, sortWeightList := getWeightByStatus(isOpenAward, poolStatus, wishBoxItems)
	indexWeightMap, sortWeightList := getWeightByStatus(isOpenAward, poolStatus, wishBoxItems, goroutineID)
	if len(indexWeightMap) == 0 {
		logs.Error("goroutineID: %v,抽奖失败原因,getWeightByStatus()返回的 indexWeightMap 长度为0", goroutineID)
		return dStatus, 0
	}
	if len(sortWeightList) == 0 {
		logs.Error("goroutineID: %v,抽奖失败原因,getWeightByStatus()返回的 sortWeightList 长度为0", goroutineID)
		return dStatus, 0
	}
	// 计算权重,抽出物品
	boxItemId := GetProductByWeight(indexWeightMap, sortWeightList, isOpenAward, goroutineID)
	if boxItemId == 0 {
		logs.Error("goroutineID: %v,抽奖失败原因,抽奖失败得到的物品id为0,数据有误,用户id为: %d,indexWeightMap: %v,sortWeightList: %v", goroutineID, pid, indexWeightMap, sortWeightList)
		return dStatus, 0
	}
	// 从数据库中查找
	item := for_game.GetWishBoxItemByIdFromDB(boxItemId)
	if item == nil {
		logs.Error("goroutineID: %v,抽奖失败原因,获取wishBoxItem 为nil,用户id为: %d,_id为: %d", goroutineID, pid, boxItemId)
		return dStatus, 0
	}

	rewardLv := item.GetRewardLv() //1、小奖 2、大奖
	var afterPrice int64
	var wishItemId int64 // 中间表的主键id

	diamond := item.GetDiamond()
	logs.Info("goroutineID: %v,物品id(中间表的): %d 的钻石价格为: %d", goroutineID, boxItemId, diamond)
	if rewardLv == for_game.ITEM_REWARDLV_2 {
		//如果是贵重物品,先上锁取水池的值(锁住当前水池不许更改.),能不能开大奖(当前盈利是否都大于该大奖物品的价格),如果属于,就开大奖,否则随机开小奖,更新水池放奖后的水量,解锁
		pp, isBigReward := poolObj.GetReward(pBase, pid, poolId, diamond, "IncomeValue", goroutineID)
		if !isBigReward { // 不能开大奖,随机抽小奖
			// 随机抽小奖的物品.
			items := for_game.FindMaxPriceBoxItemByBoxIdAndLv(boxId, for_game.ITEM_REWARDLV_1)
			if len(items) == 0 {
				logs.Error("goroutineID: %v,抽奖失败原因,不能开大奖,也没有普通的物品,boxId: %d", goroutineID, boxId)
				return dStatus, 0
			}
			r := for_game.RandInt(0, len(items))
			afterItem := items[r]

			diamond1 := afterItem.GetDiamond()
			logs.Info("goroutineID: %v,物品id(中间表的): %d 的钻石价格为: %d", goroutineID, afterItem.GetId(), diamond1)
			afterPrice = diamond1
			wishItemId = afterItem.GetId()
			// todo 白名单 判断是否是白名单用户表,白名单用户不扣水池
			// 不是白名单并且不是运营号
			if !util.Int64InSlice(pid, whiteIds) {
				//if pBase.GetTypes() != for_game.WISH_PLAYER_BASE_TYPES_2 {
				// 从水池中减去物品的价格
				pollPrice, err = poolObj.UpdatePollPriceToRedis(poolId, 0-afterPrice, "IncomeValue")
				logs.Info("goroutineID: %v,中了大奖但只能开小奖,物品价格为: %v,扣除物品后的水量为: %v", goroutineID, afterPrice, pollPrice)
				if err != nil {
					logs.Error("goroutineID: %v,抽奖失败原因,发奖后更改水池价格失败,用户id为: %d,盲盒id为: %d, 水池id为: %d,err: %s", goroutineID, pid, boxId, poolId, err.Error())
					return dStatus, 0
				}
			} else {
				logs.Info("goroutineID: %v,此用户为白名单用户,不参与水池水量变化,id为: %d", goroutineID, pid)
			}

		} else { // 大奖已经扣了大奖的的价格了.
			afterPrice = diamond
			wishItemId = item.GetId()
			if pp > 0 {
				pollPrice = pp
			}
			logs.Info("goroutineID: %v,中了大奖,用户id为: %d,物品id为: %d,扣除大奖物品后水池水量为: %d", goroutineID, pid, wishItemId, pollPrice)

		}
	} else {
		afterPrice = diamond
		wishItemId = item.GetId()
		// todo 白名单 判断是否是白名单用户表,白名单用户不扣水池
		// 不是白名单并且不是运营号
		if !util.Int64InSlice(pid, whiteIds) {
			//if pBase.GetTypes() != for_game.WISH_PLAYER_BASE_TYPES_2 {
			// 从水池中减去物品的价格
			pollPrice, err = poolObj.UpdatePollPriceToRedis(poolId, 0-afterPrice, "IncomeValue")
			logs.Info("goroutineID: %v,开小奖,物品价格为: %v,扣除物品后的水量为: %v", goroutineID, afterPrice, pollPrice)
			if err != nil {
				logs.Error("goroutineID: %v,抽奖失败原因,发奖后更改水池价格失败,用户id为: %d,盲盒id为: %d, 水池id为: %d,err: %s", goroutineID, pid, boxId, poolId, err.Error())
				return dStatus, 0
			}
		} else {
			logs.Info("goroutineID: %v,此用户为白名单用户,不参与水池水量变化,id为: %d", goroutineID, pid)
		}

	}

	pool = poolObj.GetPollInfoFromRedis()
	var wishPoolLog1BeforeValue int64
	// todo 白名单 判断是否是白名单用户表,白名单用户不扣水池
	// 不是白名单并且不是运营号
	if !util.Int64InSlice(pid, whiteIds) {
		//if pBase.GetTypes() == 0 && !util.Int64InSlice(pid, whiteIds) {
		//if pBase.GetTypes() != for_game.WISH_PLAYER_BASE_TYPES_2 {
		// 得出物品后 记录水池的流水
		wishPoolLog1 := &share_message.WishPoolLog{
			BoxId:            easygo.NewInt64(boxId),
			PoolId:           easygo.NewInt64(poolId),
			PlayerId:         easygo.NewInt64(pid),
			BeforeValue:      easygo.NewInt64(pollPrice + afterPrice),
			AfterValue:       easygo.NewInt64(pollPrice),
			Value:            easygo.NewInt64(0 - afterPrice),
			Type:             easygo.NewInt64(for_game.WISH_POOL_LOG_TYPE_3), // 3-得到物品后扣除
			LocalStatus:      easygo.NewInt64(poolStatus),
			IsOpenAward:      easygo.NewBool(isOpenAward),
			AfterLocalStatus: easygo.NewInt64(pool.GetLocalStatus()),
			AfterIsOpenAward: easygo.NewBool(pool.GetIsOpenAward()),
			CreateTimeMill:   easygo.NewInt64(for_game.GetMillSecond()),
			GoroutineID:      easygo.NewInt64(goroutineID),
		}
		wishPoolLog1BeforeValue = wishPoolLog1.GetBeforeValue()
		easygo.Spawn(for_game.AddWishPoolLog, wishPoolLog1)
	} else {
		logs.Info("goroutineID: %v,此用户为白名单用户,不参与水池水量变化,id为: %d", goroutineID, pid)
		wishPoolLog1BeforeValue = pollPrice
		pollPrice = pool.GetIncomeValue()

	}

	dStatus.AfterIsOpenAward = pool.GetIsOpenAward()
	dStatus.AfterLocalStatus = pool.GetLocalStatus()
	dStatus.IsOpenAward = isOpenAward
	dStatus.LocalStatus = int64(poolStatus)
	dStatus.PoolIncomeValue = wishPoolLog1BeforeValue
	dStatus.AfterPoolIncomeValue = pollPrice

	//   修改是是否需要附加权重.
	isOpen := checkIsOpenAward(pool.GetIncomeValue(), pool.GetStartAward(), pool.GetCloseAward(), pool.GetIsOpenAward())
	poolObj.SetIsOpenAward(poolId, isOpen)
	logs.Info("goroutineID: %v,用户id: %d, 抽到的物品id为: %d,扣除物品后水池的水量为: %v", goroutineID, pid, wishItemId, pollPrice)

	// 不是白名单并且不是运营号
	if !util.Int64InSlice(pid, whiteIds) {
		//if pBase.GetTypes() != for_game.WISH_PLAYER_BASE_TYPES_2 {
		// 再次校验是否要抽水
		if pool.GetIncomeValue() >= recycle {
			poolObj.Mutex1.Lock()
			pool = poolObj.GetPollInfoFromRedis()
			if pool.GetIncomeValue() >= recycle {
				if _, err := pump(poolId, commission, pid, boxId, recycle, int64(poolStatus), isOpenAward, poolObj, pool, goroutineID); err != nil {
					logs.Error("goroutineID: %v,再次校验抽水,抽水公共方法失败了,err: %s", goroutineID, err.Error())
				}

			}
			poolObj.Mutex1.Unlock()
		}
	} else {
		logs.Info("goroutineID: %v,此用户为白名单用户,不参与水池水量变化,id为: %d", goroutineID, pid)
	}

	// 修改不是第一次抽奖了
	if redisWishPlayer := for_game.GetRedisWishPlayer(pid); redisWishPlayer != nil {
		redisWishPlayer.SetNoOne(true)
	}

	return dStatus, wishItemId
}

// 返回挑战前和挑战后的状态
type DareStatus struct {
	AfterIsOpenAward     bool  // 抽奖后是否开启放奖 true 放奖 false 关闭
	AfterLocalStatus     int64 //  抽奖后当前水池的状态;1-大亏;2-小亏;3-普通;4-大赢;5-小盈
	IsOpenAward          bool  // 是否开启放奖 true 放奖 false 关闭
	LocalStatus          int64 // 当前水池的状态;1-大亏;2-小亏;3-普通;4-大赢;5-小盈
	PoolIncomeValue      int64 // 抽奖前水池的收入
	AfterPoolIncomeValue int64 // 抽奖后水池的收入
}

// 抽水的方法
func pump(poolId, commission, pid, boxId, recycle, localStatus int64, isOpenAward bool, poolObj *PoolObj, pool *share_message.WishPool, goroutineID uint64) (int64, error) {
	// 符合就立马抽水
	afterCommissionPrice, err := poolObj.UpdatePollPriceToRedis(poolId, 0-commission, "IncomeValue") // 抽水后的余额
	if err != nil {
		logs.Error("抽水更改水池价格失败,用户id为: %d,盲盒id为: %d, 水池id为: %d,err: %s", pid, boxId, poolId, err.Error())
		return 0, errors.New("抽水更改水池价格失败")
	}
	pollPrice := afterCommissionPrice // 水池的余额,计算水池状态使用
	logs.Info("goroutineID:%v,用户id为: %d,盲盒id为: %d,水池id为:%d,抽水后的金额为:%d", goroutineID, pid, boxId, poolId, afterCommissionPrice)
	// 确保可以抽水
	if afterCommissionPrice+commission < recycle { // 恢复水池水量,不抽水
		logs.Error("goroutineID: %v,抽水后的水量二次验证不满足抽水,抽水后的水量为: %v,校验后的水量为: %v,抽水阀值为: %v,抽水水量为: %v",
			goroutineID, afterCommissionPrice, afterCommissionPrice+commission, recycle, commission)
		pp, err := poolObj.UpdatePollPriceToRedis(poolId, commission, "IncomeValue") // 抽水后的余额
		if err != nil {
			logs.Error("goroutineID: %v,水位不够不能抽水,恢复水池价格失败,盲盒id为: %d, 水池id为: %d,err: %s", goroutineID, boxId, poolId, err.Error())
			return 0, errors.New("水位不够不能抽水")
		}
		pollPrice = pp // 水池的余额,计算水池状态使用
	} else {
		logs.Info("goroutineID: %v,抽水后设置开启放奖开始", goroutineID)
		//   修改是是否需要附加权重.
		isOpen := checkIsOpenAward(pool.GetIncomeValue(), pool.GetStartAward(), pool.GetCloseAward(), pool.GetIsOpenAward())
		poolObj.SetIsOpenAward(poolId, isOpen)
		// todo 并发测试打印
		pool1 := poolObj.GetPollInfoFromRedis()
		bytes, _ := json.Marshal(pool1)
		logs.Info("goroutineID: %v,抽水后设置开启放奖,开启放奖后的数据为--->%v", goroutineID, string(bytes))
		// todo 记录抽水日志.
		wishPoolPumpLog := &share_message.WishPoolPumpLog{
			BoxId:          easygo.NewInt64(boxId),
			PoolId:         easygo.NewInt64(poolId),
			Price:          easygo.NewInt64(commission),
			IncomeValue:    easygo.NewInt64(afterCommissionPrice + commission),
			GoroutineID:    easygo.NewInt64(goroutineID),
			CreateTimeMill: easygo.NewInt64(for_game.GetMillSecond()),
		}
		easygo.Spawn(for_game.AddWishPoolPumpLog, wishPoolPumpLog)
		pool = poolObj.GetPollInfoFromRedis()
		// 记录水池的流水
		wishPoolLog := &share_message.WishPoolLog{
			BoxId:            easygo.NewInt64(boxId),
			PoolId:           easygo.NewInt64(poolId),
			PlayerId:         easygo.NewInt64(pid),
			BeforeValue:      easygo.NewInt64(pollPrice + commission),
			AfterValue:       easygo.NewInt64(pollPrice),
			Value:            easygo.NewInt64(0 - commission),
			Type:             easygo.NewInt64(for_game.WISH_POOL_LOG_TYPE_1), //1-抽水扣除
			LocalStatus:      easygo.NewInt64(localStatus),
			IsOpenAward:      easygo.NewBool(isOpenAward),
			AfterLocalStatus: easygo.NewInt64(pool.GetLocalStatus()),
			AfterIsOpenAward: easygo.NewBool(pool.GetIsOpenAward()),
			CreateTimeMill:   easygo.NewInt64(for_game.GetMillSecond()),
			GoroutineID:      easygo.NewInt64(goroutineID),
		}
		easygo.Spawn(for_game.AddWishPoolLog, wishPoolLog)
	}
	return pollPrice, nil
}

// 后台工具抽水的方法
func toolPump(boxId, poolId, commission, pid, recycle, localStatus int64, isOpenAward bool, poolObj *PoolObj, pool *share_message.WishPool) (int64, error) {
	// 符合就立马抽水
	afterCommissionPrice, err := poolObj.UpdatePollPriceToRedis(poolId, 0-commission, "IncomeValue") // 抽水后的余额
	if err != nil {
		logs.Error("后台抽奖工具 抽水更改水池价格失败,用户id为: %d, 水池id为: %d,err: %s", pid, poolId, err.Error())
		return 0, errors.New("后台抽奖工具 抽水更改水池价格失败")
	}
	pollPrice := afterCommissionPrice // 水池的余额,计算水池状态使用
	logs.Info("后台抽奖工具 用户id为: %d, 水池id为:%d,抽水后的金额为:%d", pid, poolId, afterCommissionPrice)
	// 确保可以抽水
	if afterCommissionPrice+commission < recycle { // 恢复水池水量,不抽水
		logs.Error("后台抽奖工具 抽水后的水量二次验证不满足抽水,抽水后的水量为: %v,校验后的水量为: %v,抽水阀值为: %v,抽水水量为: %v",
			afterCommissionPrice, afterCommissionPrice+commission, recycle, commission)
		pp, err := poolObj.UpdatePollPriceToRedis(poolId, commission, "IncomeValue") // 抽水后的余额
		if err != nil {
			logs.Error("后台抽奖工具 水位不够不能抽水,恢复水池价格失败,水池id为: %d,err: %s", poolId, err.Error())
			return 0, errors.New("后台抽奖工具 水位不够不能抽水")
		}
		pollPrice = pp // 水池的余额,计算水池状态使用
	} else {
		logs.Info("后台抽奖工具 抽水后设置开启放奖开始")
		//   修改是是否需要附加权重.
		isOpen := checkIsOpenAward(pool.GetIncomeValue(), pool.GetStartAward(), pool.GetCloseAward(), pool.GetIsOpenAward())
		poolObj.SetIsOpenAward(poolId, isOpen)
		pool1 := poolObj.GetPollInfoFromRedis()
		bytes, _ := json.Marshal(pool1)
		logs.Info("后台抽奖工具 抽水后设置开启放奖,开启放奖后的数据为--->%v", string(bytes))
		wishPoolPumpLog := &share_message.WishPoolPumpLog{
			PoolId:         easygo.NewInt64(poolId),
			Price:          easygo.NewInt64(commission),
			IncomeValue:    easygo.NewInt64(afterCommissionPrice + commission),
			CreateTimeMill: easygo.NewInt64(for_game.GetMillSecond()),
			BoxId:          easygo.NewInt64(boxId),
		}
		easygo.Spawn(for_game.AddToolWishPoolPumpLog, wishPoolPumpLog)
		//  获取并修改水池状态
		//toolSetPoolStatus(poolId)
		pool = poolObj.GetPollInfoFromRedis()
		// 记录水池的流水
		wishPoolLog := &share_message.WishPoolLog{
			PoolId:           easygo.NewInt64(poolId),
			PlayerId:         easygo.NewInt64(pid),
			BeforeValue:      easygo.NewInt64(pollPrice + commission),
			AfterValue:       easygo.NewInt64(pollPrice),
			Value:            easygo.NewInt64(0 - commission),
			Type:             easygo.NewInt64(for_game.WISH_POOL_LOG_TYPE_1), //1-抽水扣除
			LocalStatus:      easygo.NewInt64(localStatus),
			IsOpenAward:      easygo.NewBool(isOpenAward),
			AfterLocalStatus: easygo.NewInt64(pool.GetLocalStatus()),
			AfterIsOpenAward: easygo.NewBool(pool.GetIsOpenAward()),
			CreateTimeMill:   easygo.NewInt64(for_game.GetMillSecond()),
			BoxId:            easygo.NewInt64(boxId),
		}
		easygo.Spawn(for_game.AddToolWishPoolLog, wishPoolLog)
	}
	return pollPrice, nil
}

// 根据状态获取权重
//func getWeightByStatus(isOpenAward bool, status int32, wishBoxItems []*share_message.WishBoxItem) (map[int64]int64, []int32) {
//func getWeightByStatus(isOpenAward bool, status int32, wishBoxItems []*share_message.WishBoxItem) (map[int]int64, []int32) {
//	//sortWeightMap := make(map[int64]int64) // 权重的map[权重]中间表物品的id
//	indexWeightMap := make(map[int]int64) // 数组索引
//	sortWeightList := make([]int32, 0)
//	if len(wishBoxItems) == 0 {
//		//return sortWeightMap, sortWeightList
//		return indexWeightMap, sortWeightList
//	}
//
//	for _, v := range wishBoxItems {
//		switch status {
//		case for_game.POOL_STATUS_BIGLOSS: // 大亏
//			bigLoss := v.GetBigLoss()
//			if isOpenAward {
//				//sortWeightMap[int64(bigLoss)] = v.GetId()
//
//				sortWeightList = append(sortWeightList, bigLoss)
//				indexWeightMap[len(sortWeightList)-1] = v.GetId()
//			} else {
//				//sortWeightMap[int64(bigLoss)] = v.GetId()
//				sortWeightList = append(sortWeightList, bigLoss)
//				indexWeightMap[len(sortWeightList)-1] = v.GetId()
//			}
//
//		case for_game.POOL_STATUS_SMALLLOSS: // 小亏
//			smallLoss := v.GetSmallLoss()
//			if isOpenAward {
//				//sortWeightMap[int64(smallLoss)] = v.GetId()
//				sortWeightList = append(sortWeightList, smallLoss)
//				indexWeightMap[len(sortWeightList)-1] = v.GetId()
//			} else {
//				//sortWeightMap[int64(smallLoss)] = v.GetId()
//				sortWeightList = append(sortWeightList, smallLoss)
//				indexWeightMap[len(sortWeightList)-1] = v.GetId()
//			}
//
//		case for_game.POOL_STATUS_COMMON: // 普通
//			common := v.GetCommon()
//			addWeight := v.GetCommonAddWeight()
//			if isOpenAward {
//				//sortWeightMap[int64(common+addWeight)] = v.GetId()
//				sortWeightList = append(sortWeightList, common+addWeight)
//				indexWeightMap[len(sortWeightList)-1] = v.GetId()
//			} else {
//				//sortWeightMap[int64(common)] = v.GetId()
//				sortWeightList = append(sortWeightList, common)
//				indexWeightMap[len(sortWeightList)-1] = v.GetId()
//			}
//
//		case for_game.POOL_STATUS_BIGWIN: // 大盈
//			bigWin := v.GetBigWin()
//			addWeight := v.GetBigWinAddWeight()
//			if isOpenAward {
//				//sortWeightMap[int64(bigWin+addWeight)] = v.GetId()
//				sortWeightList = append(sortWeightList, bigWin+addWeight)
//				indexWeightMap[len(sortWeightList)-1] = v.GetId()
//			} else {
//				//sortWeightMap[int64(bigWin)] = v.GetId()
//				sortWeightList = append(sortWeightList, bigWin)
//				indexWeightMap[len(sortWeightList)-1] = v.GetId()
//			}
//
//		case for_game.POOL_STATUS_SMALLWIN: // 小盈
//			smallWin := v.GetSmallWin()
//			addWeight := v.GetSmallWinAddWeight()
//			if isOpenAward {
//				//sortWeightMap[int64(smallWin+addWeight)] = v.GetId()
//				sortWeightList = append(sortWeightList, smallWin+addWeight)
//				indexWeightMap[len(sortWeightList)-1] = v.GetId()
//			} else {
//				//sortWeightMap[int64(smallWin)] = v.GetId()
//				sortWeightList = append(sortWeightList, smallWin)
//				indexWeightMap[len(sortWeightList)-1] = v.GetId()
//			}
//
//		}
//	}
//	return indexWeightMap, sortWeightList
//}

// todo 三爷说,如果是放奖,就用附加权重的,不需要累加了
func getWeightByStatus(isOpenAward bool, status int32, wishBoxItems []*share_message.WishBoxItem, goroutineID uint64) (map[int]int64, []int32) {
	logs.Warn("goroutineID %v,获取权重开始,是否开启放奖状态: %v,此时水池的状态为: %v", goroutineID, isOpenAward, status)
	//sortWeightMap := make(map[int64]int64) // 权重的map[权重]中间表物品的id
	indexWeightMap := make(map[int]int64) // 数组索引
	sortWeightList := make([]int32, 0)
	if len(wishBoxItems) == 0 {
		//return sortWeightMap, sortWeightList
		return indexWeightMap, sortWeightList
	}

	for _, v := range wishBoxItems {
		switch status {
		case for_game.POOL_STATUS_BIGLOSS: // 大亏
			bigLoss := v.GetBigLoss()
			if isOpenAward {
				//sortWeightMap[int64(bigLoss)] = v.GetId()

				sortWeightList = append(sortWeightList, bigLoss)
				indexWeightMap[len(sortWeightList)-1] = v.GetId()
			} else {
				//sortWeightMap[int64(bigLoss)] = v.GetId()
				sortWeightList = append(sortWeightList, bigLoss)
				indexWeightMap[len(sortWeightList)-1] = v.GetId()
			}

		case for_game.POOL_STATUS_SMALLLOSS: // 小亏
			smallLoss := v.GetSmallLoss()
			if isOpenAward {
				//sortWeightMap[int64(smallLoss)] = v.GetId()
				sortWeightList = append(sortWeightList, smallLoss)
				indexWeightMap[len(sortWeightList)-1] = v.GetId()
			} else {
				//sortWeightMap[int64(smallLoss)] = v.GetId()
				sortWeightList = append(sortWeightList, smallLoss)
				indexWeightMap[len(sortWeightList)-1] = v.GetId()
			}

		case for_game.POOL_STATUS_COMMON: // 普通
			common := v.GetCommon()
			addWeight := v.GetCommonAddWeight()
			if isOpenAward {
				//sortWeightMap[int64(common+addWeight)] = v.GetId()
				sortWeightList = append(sortWeightList, addWeight)
				indexWeightMap[len(sortWeightList)-1] = v.GetId()
			} else {
				//sortWeightMap[int64(common)] = v.GetId()
				sortWeightList = append(sortWeightList, common)
				indexWeightMap[len(sortWeightList)-1] = v.GetId()
			}

		case for_game.POOL_STATUS_BIGWIN: // 大盈
			bigWin := v.GetBigWin()
			addWeight := v.GetBigWinAddWeight()
			if isOpenAward {
				//sortWeightMap[int64(bigWin+addWeight)] = v.GetId()
				sortWeightList = append(sortWeightList, addWeight)
				indexWeightMap[len(sortWeightList)-1] = v.GetId()
			} else {
				//sortWeightMap[int64(bigWin)] = v.GetId()
				sortWeightList = append(sortWeightList, bigWin)
				indexWeightMap[len(sortWeightList)-1] = v.GetId()
			}

		case for_game.POOL_STATUS_SMALLWIN: // 小盈
			smallWin := v.GetSmallWin()
			addWeight := v.GetSmallWinAddWeight()
			if isOpenAward {
				//sortWeightMap[int64(smallWin+addWeight)] = v.GetId()
				sortWeightList = append(sortWeightList, addWeight)
				indexWeightMap[len(sortWeightList)-1] = v.GetId()
			} else {
				//sortWeightMap[int64(smallWin)] = v.GetId()
				sortWeightList = append(sortWeightList, smallWin)
				indexWeightMap[len(sortWeightList)-1] = v.GetId()
			}

		}
	}
	return indexWeightMap, sortWeightList
}

// 根据权重获取物品 得到中间表的主键id
//func GetProductByWeight(sortWeightMap map[int64]int64, sortWeightList []int32) int64 {
func GetProductByWeight(sortWeightMap map[int]int64, sortWeightList []int32, isOpenAward bool, goroutineID uint64) int64 {
	if len(sortWeightList) == 0 {
		return 0
	}
	var count int32
	for _, v := range sortWeightList {
		count += v
	}
	rate := make([]float32, 0)
	for _, key := range sortWeightList {
		// 计算百分比，保留两位小数
		f := easygo.Decimal(easygo.AtoFloat64(easygo.AnytoA(key))/easygo.AtoFloat64(easygo.AnytoA(count)), 6)
		rate = append(rate, float32(f))
	}

	index := for_game.WeightedRandomIndex(rate)

	// 找到对应的物品id
	logs.Debug("goroutineID: %v,sortWeightList---------->%v", goroutineID, sortWeightList)
	logs.Debug("goroutineID: %v,放奖开关是否打开了----->%v,index---------->%v", goroutineID, isOpenAward, index)
	logs.Debug("goroutineID: %v,sortWeightMap---------->%v", goroutineID, sortWeightMap)

	value := sortWeightMap[index]
	return value
}

// 操作开启或关闭放奖
func checkIsOpenAward(poolPrice, startAward, closeAward int64, localStatus bool) bool {
	if poolPrice >= startAward {
		localStatus = true
	}
	if poolPrice <= closeAward {
		localStatus = false
	}
	return localStatus
}

// 守护者轮播数据
func DefenderCarouselService() *h5_wish.DefenderMsgResp {
	// 随机获取10条挑战成功记录
	result := &h5_wish.DefenderMsgResp{}
	// 查询挑战成功记录
	wishLogs := for_game.GetWishOccupiedCurTime(DEFAULT_GET_COUNT_NUM)
	defenders := make([]*h5_wish.DefenderMsg, 0)
	for _, v := range wishLogs {
		var headURL string
		if base := for_game.GetRedisWishPlayer(v.GetPlayerId()); base != nil {
			headURL = base.GetHeadIcon()
		}
		defenders = append(defenders, &h5_wish.DefenderMsg{
			HeadUrl:      easygo.NewString(headURL),
			OccupiedTime: easygo.NewInt64(v.GetOccupiedTime()),
			CoinNum:      easygo.NewInt32(v.GetCoinNum()),
			CreateTime:   easygo.NewInt64(v.GetCreateTime()),
		})
	}
	result.Msg = defenders
	return result
}

// 守护者轮播数据
func GotWishCarouselService() (*h5_wish.GotWishPlayerResp, error) {
	result := &h5_wish.GotWishPlayerResp{}
	// 查询挑战成功记录
	playerWishData := for_game.GotWishPlayerLogByCurTime(DEFAULT_GET_COUNT_NUM)
	wishPlayer := make([]*h5_wish.GotWishPlayer, 0)

	dataLen := len(playerWishData)
	if dataLen > 1 {
		for _, v := range playerWishData {
			var headURL string
			var nickName string
			if base := for_game.GetRedisWishPlayer(v.GetPlayerId()); base != nil {
				headURL = base.GetHeadIcon()
				nickName = base.GetNickName()
			}
			wishPlayer = append(wishPlayer, &h5_wish.GotWishPlayer{
				HeadUrl:     easygo.NewString(headURL),
				NickName:    easygo.NewString(nickName),
				ProductIcon: easygo.NewString(v.GetProductUrl()),
			})
		}
	}
	// 如果不够,抽取运营号
	remain := DEFAULT_GET_COUNT_NUM - dataLen
	if remain > 0 {
		// 抽取运营号
		players := for_game.GetRandWishPlayer(remain)
		iconArray, err := for_game.GetRangeItemIcon(remain)
		if err != nil || len(iconArray) < remain {
			logs.Error("GotWishCarouselService err: 随机盲盒商品图片出错")
			return nil, errors.New("获取守护者轮播数据出错")
		}

		for i, v := range players {
			wishPlayer = append(wishPlayer, &h5_wish.GotWishPlayer{
				HeadUrl:     easygo.NewString(v.GetHeadUrl()),
				NickName:    easygo.NewString(v.GetNickName()),
				ProductIcon: easygo.NewString(iconArray[i]),
			})
		}
	}

	// 随机获取10条挑战成功记录
	result.Msg = wishPlayer
	return result, nil
}

func checkIsLimit(pid int64, confData *share_message.WishCoolDownConfig) (int64, bool) {
	var result bool
	if confData == nil {
		return 0, result
	}
	if !confData.GetIsOpen() {
		return 0, result
	}

	// 判断是否在冷却期间
	if t, b := for_game.IsExistPlayerCoolDownTimeFromRedis(pid); b {
		result = true
		return t, result
	}
	count := for_game.AddPlayerDareCountFromRedis(pid, confData.GetContinuousTime())
	if count-1 >= confData.GetContinuousTimes() { // 超过限制了,设置冷却期
		result = true
		for_game.SetPlayerCoolDownTimeFromRedis(pid, confData.GetCoolDownTime())
	}
	return time.Now().Unix(), result
}

// 更多挑战
func BoxListService() (*h5_wish.BoxListResp, error) {
	// 随机查询10条盲盒数据.
	boxes, err := for_game.FindWishBoxByNum(DEFAULT_GET_COUNT_NUM)
	if err != nil {
		logs.Error("随机查询10条盲盒数据失败 err: ", err.Error())
		return nil, err
	}
	boxList := make([]*h5_wish.QueryBox, 0)
	for _, v := range boxes {
		boxList = append(boxList, &h5_wish.QueryBox{
			Id:            easygo.NewInt64(v.GetId()),
			Name:          easygo.NewString(v.GetName()),
			Image:         easygo.NewString(v.GetIcon()),
			Desc:          easygo.NewString(v.GetDesc()),
			Price:         easygo.NewInt64(v.GetPrice()),
			ProductStatus: easygo.NewInt32(v.GetProductStatus()),
			Label:         easygo.NewInt32(v.GetMatch()),
		})
	}
	return &h5_wish.BoxListResp{
		BoxList: boxList,
	}, nil
}

// 硬币兑换砖石
func CoinToDiamondService(pid, id int64, data1 *share_message.DiamondRecharge) (int64, int64, int64, easygo.IMessage) {
	first := data1.GetMonthFirst()                //首充数量
	if first > 0 && data1.GetGiveDiamond() != 0 { // 判断这个月是否是首充,二者不能并存
		logs.Error("该兑换配置数据有误,不可以同事存在月首充和赠送硬币")
		return 0, 0, 0, easygo.NewFailMsg("参数有误")
	}
	// 获取柠檬的playerId
	wishPlayer := for_game.GetRedisWishPlayer(pid)
	if wishPlayer == nil {
		logs.Error("获取许愿池用户失败")
		return 0, 0, 0, easygo.NewFailMsg("参数有误")
	}
	diamond := data1.GetDiamond()
	t := time.Now().Unix()
	if t >= data1.GetStartTime() && t <= data1.GetEndTime() {
		diamond += data1.GetGiveDiamond()
	}
	if first > 0 { // 判断这个月是否是首充,二者不能并存
		logObj := for_game.GetRedisDiamondLogObj()
		if !logObj.CheckMonthRecharge(pid, id) && for_game.CheckPlayerDiamondMonthRecharge(pid, id) {
			diamond += first
		}

	}
	// 扣除用户的硬币,优先扣除绑定硬币,在扣除非绑定硬币
	changDiamond := wishPlayer.GetDiamond() + diamond
	coinData := &h5_wish.AddCoinReq{
		UserId:     easygo.NewInt64(wishPlayer.GetPlayerId()),
		Coin:       easygo.NewInt64(0 - data1.GetCoinPrice()),
		SourceType: easygo.NewInt32(for_game.COIN_TYPE_WISH_PAY),
		Diamond:    easygo.NewInt64(changDiamond),
	}
	resp, err1 := SendMsgToIdelServer(for_game.SERVER_TYPE_HALL, "RpcAddCoin", coinData, wishPlayer.GetPlayerId())
	logs.Info("硬币兑换砖石,RpcAddCoin 返回resp: %v,err: %v", resp, err1)
	if err1 != nil {
		logs.Error("DoDareService-err:", err1)
		return 0, 0, 0, easygo.NewFailMsg("扣费失败")
	}
	data := resp.(*h5_wish.AddCoinResp)
	if data.GetResult() == 2 { // 处理扣费失败
		return 0, 0, 0, easygo.NewFailMsg("扣费失败")
	}
	extendLog := &share_message.GoldExtendLog{
		RedPacketId: easygo.NewInt64(id), // 兑换钻石的配置表 id
		PlayerId:    easygo.NewInt64(wishPlayer.GetPlayerId()),
	}
	err, newDiamond := wishPlayer.AddDiamond(diamond, "兑换", for_game.DIAMOND_TYPE_EXCHANGE_IN, extendLog)
	if err != nil {
		return 0, 0, 0, err
	}
	// 设置最新兑换钻石的时间
	wishPlayer.SetLastExchangeDiamondTime(for_game.GetMillSecond())
	return newDiamond, data.GetCoin(), diamond, nil
}

func DiamondRechargeListService(uid int64) []*h5_wish.DiamondRecharge {
	result := make([]*h5_wish.DiamondRecharge, 0)
	list := for_game.GetDiamondRechargeList()
	if len(list) == 0 {
		return result
	}
	logObj := for_game.GetRedisDiamondLogObj()
	for _, v := range list {
		var isFirst bool
		if v.GetMonthFirst() > 0 && !logObj.CheckMonthRecharge(uid, v.GetId()) && for_game.CheckPlayerDiamondMonthRecharge(uid, v.GetId()) {
			isFirst = true
		}
		result = append(result, &h5_wish.DiamondRecharge{
			Id:           easygo.NewInt64(v.GetId()),
			Diamond:      easygo.NewInt64(v.GetDiamond()),
			CoinPrice:    easygo.NewInt64(v.GetCoinPrice()),
			MonthFirst:   easygo.NewInt64(v.GetMonthFirst()),
			Rebate:       easygo.NewInt32(v.GetRebate()),
			StartTime:    easygo.NewInt64(v.GetStartTime()),
			EndTime:      easygo.NewInt64(v.GetEndTime()),
			Status:       easygo.NewInt32(v.GetStatus()),
			Sort:         easygo.NewInt32(v.GetSort()),
			DisPrice:     easygo.NewInt64(v.GetDisPrice()),
			GiveDiamond:  easygo.NewInt64(v.GetGiveDiamond()),
			IsMonthFirst: easygo.NewBool(isFirst),
		})
	}
	return result
}

func GetDiamondChangeLogByPageService(uid int64, reqMsg *h5_wish.DiamondChangeLogReq) *h5_wish.DiamondChangeLogResp {
	// 存数据库
	for_game.SaveDiamondChangeLogToMongoDB()
	page, pageSize := for_game.MakePageAndPageSize(reqMsg.GetPage(), reqMsg.GetPageSize())
	diamondLogs, count := for_game.GetDiamondChangeLogByPage(uid, reqMsg.GetType(), page, pageSize)
	list := make([]*h5_wish.DiamondChangeLog, 0)
	for _, v := range diamondLogs {
		list = append(list, &h5_wish.DiamondChangeLog{
			LogId:         easygo.NewInt64(v.GetLogId()),
			PlayerId:      easygo.NewInt64(v.GetPlayerId()),
			ChangeDiamond: easygo.NewInt64(v.GetChangeDiamond()),
			SourceType:    easygo.NewInt32(v.GetSourceType()),
			PayType:       easygo.NewInt32(v.GetPayType()),
			CurDiamond:    easygo.NewInt64(v.GetCurDiamond()),
			Diamond:       easygo.NewInt64(v.GetDiamond()),
			Note:          easygo.NewString(v.GetNote()),
			CreateTime:    easygo.NewInt64(v.GetCreateTime()),
		})
	}
	return &h5_wish.DiamondChangeLogResp{
		DiamondChangeLogList: list,
		Count:                easygo.NewInt32(count),
	}
}

// 校验守护者是否过期了.
func CheckGuardianIsEp(boxId int64) {
	box, err := for_game.GetWishBox(boxId, []int64{for_game.WISH_PUT_ON_STATUS})
	if err != nil {
		return
	}
	if box.GetId() == 0 {
		return
	}
	if !box.GetIsGuardian() || box.GetGuardianId() == 0 || box.GetGuardianOverTime() <= 0 {
		return
	}
	// 判断时间是否大于24小时
	afterTime := box.GetGuardianOverTime() - time.Now().Unix()
	logs.Info("afterTime==============>", afterTime)
	if afterTime > 0 {
		return
	}

	// 清空盲盒的守护者
	guardianId := box.GetGuardianId()
	box.IsGuardian = easygo.NewBool(false)
	box.GuardianOverTime = easygo.NewInt64(0)
	box.GuardianId = easygo.NewInt64(0)
	if err := for_game.UpWishBox(boxId, box); err != nil {
		logs.Error("定时修改盲盒的守护者信息失败,err: %s", err.Error())
	}
	// 修改守护者占领的时间
	if err := for_game.UpOccupiedEx(boxId, guardianId, for_game.AFTER_TIME_GUARDIAN); err != nil {
		logs.Error("修改守护时间表信息失败,boxId: %d, err: %s", boxId, err.Error())
	}

}

//var WishDataItemList = make([]*share_message.PlayerWishItem, 0)
//
//func SaveCsv() {
//	logs.Info("WishDataItemList===>", len(WishDataItemList))
//	for_game.WriteCsv("得到的物品.csv", WishDataItemList) // todo 操作csv,跑大数据使用.
//	WishDataItemList = make([]*share_message.PlayerWishItem, 0)
//}

/*
todo 以下为连抽测试代码.
// 批量抽奖,数据统一
func BatchDoDareService(reqMsg *h5_wish.DoDareReq, userId int64, wishData *share_message.PlayerWishData, wishItemChan chan *share_message.PlayerWishItem, countChan chan int) (*h5_wish.DoDareResp, error) {
	boxId := reqMsg.GetWishBoxId()
	wishBox, err := for_game.GetWishBox(boxId)
	if err != nil {
		return nil, errors.New("数据错误")
	}
	// 拿到兑换比例,扣钱之前确保流程能跑通使用
	cfg := for_game.GetCurrencyCfg()
	if cfg == nil {
		logs.Error("获取兑换配置失败")
		return nil, errors.New("获取兑换配置失败")
	}

	c := cfg.GetCoin()
	d := cfg.GetDiamond()
	if c == 0 || d == 0 {
		logs.Error("配置表配置的兑换汇率为空")
		return nil, errors.New("兑换汇率为空")
	}
	base := for_game.GetRedisWishPlayer(userId)
	if base == nil {
		logs.Error("获取许愿池用户对象失败,许愿池用户id为: %d", userId)
		return nil, errors.New("用户信息有误")
	}
	wishBoxPrice := wishBox.GetPrice()
	imPid := base.GetPlayerId()
	extendLog := &share_message.GoldExtendLog{
		PlayerId: easygo.NewInt64(imPid),
	}

	err11, _ := base.AddDiamond(0-wishBoxPrice, "抽奖", for_game.DIAMOND_TYPE_PLAYER_OUT, extendLog)
	if err11 != nil {
		return nil, errors.New(err11.GetReason())
	}

	IsLucky := false
	IsOnce := false
	randItemId := MakeWishLucky(userId, wishBox, wishData, cfg) // 挑战赛抽奖 得出中间表的id
	if randItemId == 0 {
		logs.Error("抽奖得到的物品id为0,需要人工排查,userId= %d,wishBox= %v,wishData= %v", userId, wishBox, wishData)
		logs.Info("-----恢复用户的金额------")

		err11, _ := base.AddDiamond(wishBoxPrice, "许愿池抽奖失败钻石返回", for_game.DIAMOND_TYPE_WISH_FAILD, extendLog)
		if err11 != nil {
			return nil, errors.New(err11.GetReason())
		}

		logs.Info("-----恢复用户的金额结束------")
		return nil, errors.New("抽奖失败")
	}
	// 异步修改当前的库存量
	easygo.Spawn(IncLocalNum, reqMsg.GetWishBoxId(), randItemId)
	// 异步小助手推送
	//easygo.Spawn(pushDareResult, randItemId)

	wishBoxItem, _ := for_game.GetWishBoxItem(randItemId) // 中奖信息

	wishItem, _ := for_game.QueryWishItemById(wishBoxItem.GetWishItemId())
	dareId := userId
	//playerPlayerInfo := for_game.GetWishPlayerInfo(dareId) //挑战者信息
	var beDareId int64
	var beDareNickName string
	//挑战赛
	if reqMsg.GetDareType() == WISH_DARE {
		if randItemId == wishData.GetWishBoxItemId() {
			IsLucky = true
		}
		Protector(wishBox, wishData, wishBoxItem, IsLucky)
		beDareId = wishBox.GetGuardianId()
		guardianPlayerInfo := for_game.GetWishPlayerInfo(wishBox.GetGuardianId()) //守护者信息
		beDareNickName = guardianPlayerInfo.GetNickName()
	}

	result := false
	if IsLucky {
		result = true
	}

	var headIcon, nickName string
	if p := for_game.GetRedisWishPlayer(wishData.GetPlayerId()); p != nil {
		headIcon = p.GetHeadIcon()
		nickName = p.GetNickName()
	}

	WishLog := &share_message.WishLog{
		WishBoxId:       easygo.NewInt64(boxId),
		DareId:          easygo.NewInt64(dareId),
		DareName:        easygo.NewString(nickName),
		BeDareId:        easygo.NewInt64(beDareId),
		BeDareName:      easygo.NewString(beDareNickName),
		Result:          easygo.NewBool(result),
		ChallengeItemId: easygo.NewInt64(wishBoxItem.GetId()),
		DareHeadIcon:    easygo.NewString(headIcon),
		DarePrice:       easygo.NewInt64(wishBox.GetPrice()),
		WishItemId:      easygo.NewInt64(wishBoxItem.GetWishItemId()),
		DareType:        easygo.NewInt32(reqMsg.GetDareType()),
	}
	for_game.AddWishLog(WishLog) //添加盲盒挑战记录

	//expireTime := time.Now().Unix() + AFTER_TIME_BACK
	id := for_game.NextId(for_game.TABLE_PLAYER_WISH_ITEM)
	playerWishItem := &share_message.PlayerWishItem{
		Id:              easygo.NewInt64(id),
		PlayerId:        easygo.NewInt64(userId),
		ChallengeItemId: easygo.NewInt64(wishBoxItem.GetId()),
		Status:          easygo.NewInt32(0),
		WishBoxId:       easygo.NewInt64(wishBox.GetId()),
		IsRead:          easygo.NewBool(false),
		WishItemId:      easygo.NewInt64(wishBoxItem.GetWishItemId()),
		WishItemPrice:   easygo.NewInt64(wishBoxItem.GetPrice()),
		DareDiamond:     easygo.NewInt64(wishBox.GetPrice()),
		ProductName:     easygo.NewString(wishItem.GetName()),
		WishItemDiamond: easygo.NewInt64(wishBoxItem.GetDiamond()),
	}
	err = for_game.AddPlayerWishItem(playerWishItem) //添加玩家物品
	easygo.PanicError(err)
	wishItemChan <- playerWishItem // 批量统计
	countChan <- 1

	//WishDataItemList = append(WishDataItemList, playerWishItem) // todo csv 使用

	err = for_game.SetDareFrequency(userId) //挑战成功，增加挑战次数
	easygo.PanicError(err)

	frequency := for_game.GetDareFrequency(userId)
	if frequency <= 1 && IsLucky { //第一次就中
		IsOnce = true
	}

	// 异步修改水池状态
	fun := func() {
		status := GetPoolStatus(wishBox.GetWishPoolId())
		poolObj := GetPoolObj(wishBox.GetWishPoolId())
		poolObj.SetLocalStatus(wishBox.GetWishPoolId(), int64(status))
	}
	easygo.Spawn(fun)

	respData := &h5_wish.DoDareResp{
		IsLucky:      easygo.NewBool(IsLucky),
		ProductId:    easygo.NewInt64(wishBoxItem.GetWishItemId()),
		ProductName:  easygo.NewString(wishItem.GetName()),
		Image:        easygo.NewString(wishItem.GetIcon()),
		ProductType:  easygo.NewInt32(wishBoxItem.GetStyle()),
		IsOnce:       easygo.NewBool(IsOnce),
		Status:       easygo.NewInt32(wishBoxItem.GetStatus()),
		ForecastTime: easygo.NewInt64(wishBoxItem.GetPredictArrivalTime()),
	}

	return respData, nil
}

// 测试连抽
func MatchTest() {
	st := for_game.GetMillSecond()
	doDareReq := &h5_wish.DoDareReq{
		DareType:  easygo.NewInt32(2),
		WishBoxId: easygo.NewInt64(18),
	}
	wishItemChan := make(chan *share_message.PlayerWishItem)
	countChan := make(chan int)
	resultChan := make(chan []*share_message.PlayerWishItem)
	wishData, _ := for_game.GetWishDataByStatus(18801001, 18, for_game.WISH_CHALLENGE_WAIT)

	for i := 0; i < 100; i++ {
		easygo.Spawn(BatchDoDareService, doDareReq, int64(18801001), wishData, wishItemChan, countChan)
	}
	easygo.Spawn(ForChan, 100, wishItemChan, resultChan)

	data, ok := <-resultChan
	logs.Info("==========>", ok)
	bytes, _ := json.Marshal(data)
	ed := for_game.GetMillSecond()
	logs.Info("---------->", ed-st, string(bytes))
	logs.Info("100连抽总耗时===========>", ed-st)
	return
}
func ForChan(count int, wishItemChan chan *share_message.PlayerWishItem, resultChan chan []*share_message.PlayerWishItem) {
	result := make([]*share_message.PlayerWishItem, 0)
	for {
		select {
		case data := <-wishItemChan:
			logs.Info("1111111----------wishItemChan-----------111111111")
			result = append(result, data)
			if len(result) == count {
				resultChan <- result
			}
		}
	}
}
*/

// todo 策划需求测试
func BatchDareCH(reqMsg *h5_wish.DoDareReq, userId int64, wishData *share_message.PlayerWishData) {
	var count int
	m := make(map[string]int) // [物品名字] 次数

	for {
		count++

		resp, err := DoDareService1(reqMsg, userId, wishData, GetGoroutineID())
		if err != nil {
			logs.Error("---------->err", err.Error())
			return
		}
		c := m[resp.GetProductName()] + 1

		m[resp.GetProductName()] = c

		// 抽到大奖就结束
		if resp.GetProductId() == 60 || resp.GetProductId() == 61 {
			break
		}

	}
	logs.Info("挑战总次数为:----->", count)
	bytes, _ := json.Marshal(m)
	logs.Info("map:----->", string(bytes))
}

// 获取协程的id,并发使用
func GetGoroutineID() uint64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// 定时同步top
func SumGuardianCoinNum() {
	logs.Info("整点同步挑战记录排行榜开始=========================")
	// 先去查询上次同步的时间
	data := for_game.GetLastCoinTopLog()
	var st int64
	if data != nil {
		st = data.GetUpdateCoinTime()
	}
	list := for_game.GetGuardianCoinNumList(st, for_game.GetMillSecond())
	if len(list) > 0 {
		for _, v := range list {
			pid := v["_id"].(int64)
			coinNum := v["Count"].(int64)
			data := &share_message.WishGuardianTopLog{
				Id:             easygo.NewInt64(pid),
				PlayerId:       easygo.NewInt64(pid),
				UpdateCoinTime: easygo.NewInt64(for_game.GetMillSecond()),
				UpdateWishTime: nil,
			}
			for_game.UpdateWishGuardianTopLog(data, coinNum, 0)
		}
	}

	// 挑战成功数.
	data1 := for_game.GetLastWishTopLog()
	var st1 int64
	if data != nil {
		st1 = data1.GetUpdateWishTime()
	}
	list1 := for_game.GetGuardianWishNumList(st1, for_game.GetMillSecond())
	for _, v := range list1 {
		pid := v["_id"].(int64)
		wishNum := v["Count"].(int64)
		data := &share_message.WishGuardianTopLog{
			Id:             easygo.NewInt64(pid),
			PlayerId:       easygo.NewInt64(pid),
			UpdateWishTime: easygo.NewInt64(for_game.GetMillSecond()),
		}
		for_game.UpdateWishGuardianTopLog(data, 0, wishNum)
	}
	logs.Info("整点同步挑战记录排行榜结束=========================")
}

// 许愿池抽奖活动
/**
1.去活动表中获取活动判断活动时间是否开启
2.根据盲盒找到对应的奖池
3.根据奖池找到对应的规则列表.
4.遍历规则列表确定真正的规则,签到,次数,总金额等统计
*/
func WishActService(pid, boxId, diamond int64, goroutineID uint64) {
	logs.Info("===================goroutineID=%v,pid=%v, boxId=%v, diamond=%v,活动接口开始 WishActService==========================", goroutineID, pid, boxId, diamond)
	if boxId == 0 {
		logs.Error("goroutineID: %v,盲盒id为为0", goroutineID)
		return
	}
	isOpenCount, isOpenDay, isOpenWeekMonth := for_game.CheckWishActIsOpen()
	logs.Info("===================goroutineID=%v,isOpenCount=%v, isOpenDay=%v, isOpenWeekMonth=%v ==========================", goroutineID, isOpenCount, isOpenDay, isOpenWeekMonth)
	//统计投入钻石总额
	if isOpenWeekMonth {
		if err := sumDiamond(pid, diamond); err != nil {
			logs.Error("UpsertWishPlayerActivity sumDiamond err: %s", err.Error())
			return
		}
	}

	// 根据盲盒找到对应的奖池
	actPool := for_game.GetWishActPoolByBoxId(boxId)
	if actPool == nil {
		logs.Error("goroutineID: %v,根据盲盒找到对应的奖池失败", goroutineID)
		return
	}
	ruleList := for_game.GeWishActPoolRuleByPId(actPool.GetId())
	if len(ruleList) == 0 {
		return
	}
	var hasAddLun bool
	for _, v := range ruleList {
		t := v.GetType()
		if t == 0 || t == for_game.WISH_ACTIVITY_DATA_TYPE_3 || t == for_game.WISH_ACTIVITY_DATA_TYPE_4 { // 3-周排名,4-月排名
			continue
		}
		switch t {
		case for_game.WISH_ACTIVITY_DATA_TYPE_1: // 次数
			if isOpenCount {
				logs.Info("次数sumCount次数------>")
				sumCount(pid, v)
			}
		case for_game.WISH_ACTIVITY_DATA_TYPE_2: // 天数
			if isOpenDay {
				if b := sumDay(pid, v, actPool.GetId()); b {
					hasAddLun = b
				}
			}
		}
	}
	if isOpenDay && hasAddLun {
		/*		act := for_game.GeWishPlayerActivityByPId(pid)
				if act == nil {
					return
				}
				act.DayLun = easygo.NewInt32(act.GetDayLun() + 1)
				// 修改数据.
				if err := for_game.UpsertWishPlayerActivity(act); err != nil {
					logs.Error("UpsertWishPlayerActivity err: %s", err.Error())
					return
				}*/
	}

}

// 统计天数
func sumDay(pid int64, rule *share_message.WishActPoolRule, actPoolId int64) bool {
	logs.Info("===========统计天数开始=========")
	// 抽奖天数统计
	dayLog := &share_message.WishDayActivityLog{
		Id:            easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_DAY_ACTIVITY_LOG)),
		PlayerId:      easygo.NewInt64(pid),
		WishActPoolId: easygo.NewInt64(actPoolId),
		CreateTime:    easygo.NewInt64(easygo.GetToday0ClockTimestamp() * 1000),
	}
	for_game.InsertWishDayActivityLog(dayLog)
	// 根据用户id查找抽奖信息
	act := for_game.GeWishPlayerActivityByPId(pid)
	if act == nil {
		return false
	}
	// 判断是否存在当前的奖池数据
	var b bool
	for _, v := range act.GetData() {
		if v.GetPoolRuleId() == rule.GetId() {
			b = true
			break
		}
	}
	if !b {
		updateAct(rule, act)
		// todo 检测是否到了获奖记录
		act1 := for_game.GeWishPlayerActivityByPId(pid)
		if act1 != nil {
			for _, v := range act1.GetData() {
				if v.GetPoolRuleId() != rule.GetId() {
					continue
				}
				key := rule.GetKey()
				logs.Info("规则id为: %d,限制的天数为:%d,当前的天数为:%d", rule.GetId(), key, v.GetValue())

				if v.GetValue() < int64(key) {
					continue
				}
				var nickName string
				if p := for_game.GetRedisWishPlayer(pid); p != nil {
					nickName = p.GetNickName()
				}
				// 插入一条获奖情况记录
				actLog := for_game.NextId(for_game.TABLE_WISH_ACTIVITY_PRIZE_LOG)
				//if insertActLog(actLog, pid, rule, nickName, act1.GetDayLun()) {
				if insertActLog(actLog, pid, rule, nickName, v.GetDayLun()) {
					continue
				}
			}
		}
		return false
	}
	var nickName string
	if p := for_game.GetRedisWishPlayer(pid); p != nil {
		nickName = p.GetNickName()
	}

	var bb bool
	//dayLun := act.GetDayLun()
	var dayLun int32
	// 存在
	for _, v := range act.GetData() {
		if v.GetPoolRuleId() != rule.GetId() {
			continue
		}
		// 判断v里面的时间是不是当天
		if easygo.IsTodayTimestamp(v.GetUpdateTime()) {
			return false
		}
		// 判断天数时间是否需要加 1
		logs.Info("间隔的天数----------->", easygo.GetDifferenceDay(v.GetUpdateTime(), for_game.GetMillSecond()))
		if easygo.GetDifferenceDay(v.GetUpdateTime(), for_game.GetMillSecond()) == 1 { // 加 1 天
			v.Value = easygo.NewInt64(v.GetValue() + 1)
			v.UpdateTime = easygo.NewInt64(for_game.GetMillSecond())
			dayLun = v.GetDayLun()
		} else { // 重置天数
			v.Value = easygo.NewInt64(1)
			v.UpdateTime = easygo.NewInt64(for_game.GetMillSecond())
			bb = true
			v.DayLun = easygo.NewInt32(v.GetDayLun() + 1) // 轮数加1
			dayLun = v.GetDayLun()
		}
		// 修改数据.
		if err := for_game.UpsertWishPlayerActivity(act); err != nil {
			logs.Error("UpsertWishPlayerActivity err: %s", err.Error())
			continue
		}
		key := rule.GetKey()
		logs.Info("规则id为: %d,限制的天数为:%d,当前的天数为:%d", rule.GetId(), key, v.GetValue())

		if v.GetValue() < int64(key) {
			continue
		}

		// 插入一条获奖情况记录
		actLog := for_game.NextId(for_game.TABLE_WISH_ACTIVITY_PRIZE_LOG)
		//if insertActLog(actLog, pid, rule, nickName, act.GetDayLun()) {
		if insertActLog(actLog, pid, rule, nickName, dayLun) {
			continue
		}
	}

	return bb
}

// 插入一条获奖情况记录
func insertActLog(actLogId, pid int64, rule *share_message.WishActPoolRule, nickName string, dayLun ...int32) bool {
	prizeLog := for_game.GeWishActivityPrizeLogByPId(pid, rule.GetId(), dayLun...)
	if prizeLog != nil {
		if prizeLog.GetStatus() == for_game.WISH_ACTIVITY_PRIZE_STATUS_1 {
			marshal, _ := json.Marshal(prizeLog)
			logs.Error("insert 失败,marshal-------->", string(marshal))
			return true
		}
	} else {
		prizeLog = &share_message.WishActivityPrizeLog{
			Id:                easygo.NewInt64(actLogId),
			PlayerId:          easygo.NewInt64(pid),
			PlayerAccount:     easygo.NewString(nickName),
			Status:            easygo.NewInt32(for_game.WISH_ACTIVITY_PRIZE_STATUS_0),
			WishActPoolRuleId: easygo.NewInt64(rule.GetId()),
			CreateTime:        easygo.NewInt64(for_game.GetMillSecond()),
		}
	}
	if len(dayLun) > 0 {
		prizeLog.DayLun = easygo.NewInt32(dayLun[0])
	}
	at := rule.GetAwardType()
	prizeLog.Type = easygo.NewInt32(rule.GetType())
	prizeLog.ActType = easygo.NewInt64(rule.GetKey())
	prizeLog.PrizeType = easygo.NewInt64(at)
	var prizeValue int64
	var note string
	switch at { // 1-钻石;2-实物
	case WISH_ACT_PRIZE_TYPE_DIAMOND:
		prizeValue = rule.GetDiamond()
		note = fmt.Sprintf("获得钻石*%d", rule.GetDiamond())
	case WISH_ACT_PRIZE_TYPE_PRODUCT:
		prizeValue = rule.GetWishItemId()
		wishItem := for_game.GetWishItemByIdFromDB(rule.GetWishItemId())
		if wishItem.GetName() == "" {
			logs.Error("更改获奖记录,获取物品名字失败,WishItemId: %d", rule.GetWishItemId())
		}
		note = fmt.Sprintf("获得%s", wishItem.GetName())
	}
	prizeLog.PrizeValue = easygo.NewInt64(prizeValue)
	prizeLog.Status = easygo.NewInt32(0)
	prizeLog.Note = easygo.NewString(note)
	if err := for_game.UpsertWishActivityPrizeLog(prizeLog); err != nil {
		logs.Error("UpsertWishActivityPrizeLog err: %s", err.Error())
	}
	return false
}
func insertActLog1(actLogId, pid int64, rule *share_message.WishActPoolRule, nickName string, dayLun ...int32) bool {
	prizeLog := &share_message.WishActivityPrizeLog{
		Id:                easygo.NewInt64(actLogId),
		PlayerId:          easygo.NewInt64(pid),
		PlayerAccount:     easygo.NewString(nickName),
		Status:            easygo.NewInt32(for_game.WISH_ACTIVITY_PRIZE_STATUS_0),
		WishActPoolRuleId: easygo.NewInt64(rule.GetId()),
		CreateTime:        easygo.NewInt64(for_game.GetMillSecond()),
	}

	if len(dayLun) > 0 {
		prizeLog.DayLun = easygo.NewInt32(dayLun[0])
	}
	at := rule.GetAwardType()
	prizeLog.Type = easygo.NewInt32(rule.GetType())
	prizeLog.ActType = easygo.NewInt64(rule.GetKey())
	prizeLog.PrizeType = easygo.NewInt64(at)
	var prizeValue int64
	var note string
	switch at { // 1-钻石;2-实物
	case WISH_ACT_PRIZE_TYPE_DIAMOND:
		prizeValue = rule.GetDiamond()
		note = fmt.Sprintf("获得钻石*%d", rule.GetDiamond())
	case WISH_ACT_PRIZE_TYPE_PRODUCT:
		prizeValue = rule.GetWishItemId()
		wishItem := for_game.GetWishItemByIdFromDB(rule.GetWishItemId())
		if wishItem.GetName() == "" {
			logs.Error("更改获奖记录,获取物品名字失败,WishItemId: %d", rule.GetWishItemId())
		}
		note = fmt.Sprintf("获得%s", wishItem.GetName())
	}
	prizeLog.PrizeValue = easygo.NewInt64(prizeValue)
	prizeLog.Status = easygo.NewInt32(0)
	prizeLog.Note = easygo.NewString(note)
	if err := for_game.UpsertWishActivityPrizeLog(prizeLog); err != nil {
		logs.Error("UpsertWishActivityPrizeLog err: %s", err.Error())
	}
	return false
}

// 统计次数
func sumCount(pid int64, rule *share_message.WishActPoolRule) {
	logs.Info("===========统计次数开始=========")
	// 根据用户id查找抽奖信息
	act := for_game.GeWishPlayerActivityByPId(pid)
	if act == nil {
		return
	}

	// 判断是否存在当前的奖池数据
	var b bool
	for _, v := range act.GetData() {
		if v.GetPoolRuleId() == rule.GetId() {
			b = true
			break
		}
	}
	if !b {
		updateAct(rule, act)
		return
	}

	// 存在
	for _, v := range act.GetData() {
		if v.GetPoolRuleId() != rule.GetId() {
			continue
		}
		v.Value = easygo.NewInt64(v.GetValue() + 1)
		v.UpdateTime = easygo.NewInt64(for_game.GetMillSecond())

		// 修改数据.
		if err := for_game.UpsertWishPlayerActivity(act); err != nil {
			logs.Error("UpsertWishPlayerActivity err: %s", err.Error())
			return
		}
		key := rule.GetKey()
		logs.Info("规则id为: %d,限制的次数数为:%d,当前的次数为:%d", rule.GetId(), key, v.GetValue())
		if v.GetValue() < int64(key) {
			continue
		}

		var nickName string
		if p := for_game.GetRedisWishPlayer(pid); p != nil {
			nickName = p.GetNickName()
		}
		// 插入一条获奖情况记录
		actLog := for_game.NextId(for_game.TABLE_WISH_ACTIVITY_PRIZE_LOG)
		if insertActLog(actLog, pid, rule, nickName) {
			continue
		}
	}
}

// 修改,内部公共方法
func updateAct(rule *share_message.WishActPoolRule, act *share_message.WishPlayerActivity) {
	var dayLun int32
	if rule.GetType() == for_game.WISH_ACTIVITY_DATA_TYPE_2 {
		dayLun = 1
	}
	actData := &share_message.ActivityData{
		PoolRuleId: easygo.NewInt64(rule.GetId()),
		Value:      easygo.NewInt64(1), // 第一次
		Type:       easygo.NewInt32(rule.GetType()),
		UpdateTime: easygo.NewInt64(for_game.GetMillSecond()),
		DayLun:     easygo.NewInt32(dayLun),
	}
	data := act.GetData()
	data = append(data, actData)
	act.Data = data
	// 插入进数据库
	if err := for_game.UpsertWishPlayerActivity(act); err != nil {
		logs.Error("UpsertWishPlayerActivity err: %s", err.Error())
	}
}

// 统计投入钻石总额
func sumDiamond(pid, diamond int64) error {
	// 根据用户id查找抽奖信息
	act := for_game.GeWishPlayerActivityByPId(pid)
	if act == nil { // 新建
		playerAct := &share_message.WishPlayerActivity{
			PlayerId:   easygo.NewInt64(pid),
			Diamond:    easygo.NewInt64(diamond),
			WeekPrize:  easygo.NewInt64(diamond),
			MonthPrize: easygo.NewInt64(diamond),

			//DayLun:            easygo.NewInt32(1), // 第一轮
			UpdateTime:        easygo.NewInt64(for_game.GetMillSecond()),
			UpdateDiamondTime: easygo.NewInt64(for_game.GetMillSecond()),
		}
		// 插入进数据库
		return for_game.UpsertWishPlayerActivity(playerAct)
	}
	// 存在
	act.Diamond = easygo.NewInt64(act.GetDiamond() + diamond)
	act.WeekPrize = easygo.NewInt64(act.GetWeekPrize() + diamond)
	act.MonthPrize = easygo.NewInt64(act.GetMonthPrize() + diamond)
	act.UpdateTime = easygo.NewInt64(for_game.GetMillSecond())
	act.UpdateDiamondTime = easygo.NewInt64(for_game.GetMillSecond())
	return for_game.UpsertWishPlayerActivity(act)
}

// 统计周,月数量 ,t 1-周,2-月
func SumWeekMonth(t int) {
	logs.Info("统计周,月数量 ,t=%d, 1-周,2-月", t)
	if t != for_game.WISH_ACT_H5_WEEK && t != for_game.WISH_ACT_H5_MONTH {
		return
	}
	// 查找所有的用户消费记录
	list := for_game.GetWishPlayerActivityList()
	if len(list) == 0 {
		return
	}
	var dataType int
	switch t {
	case for_game.WISH_ACT_H5_WEEK:
		dataType = for_game.WISH_ACTIVITY_DATA_TYPE_3
		//   清空周榜表数据
		_ = for_game.RemoveAllWeekTop()
		// 	获取前1000个排行的数据
		actList := for_game.GetWishPlayerActivityTop(for_game.WISH_ACTIVITY_DATA_TYPE_3, WISH_ACT_TOP_NUM)
		saveData := make([]interface{}, 0)
		for _, v := range actList {
			saveData = append(saveData, bson.M{"_id": v.GetPlayerId()}, &share_message.WishWeekTop{
				PlayerId:   easygo.NewInt64(v.GetPlayerId()),
				WeekPrize:  easygo.NewInt64(v.GetWeekPrize()),
				CreateTime: easygo.NewInt64(time.Now().Unix()),
			})
		}
		// 插入新表
		if len(saveData) > 0 {
			for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_WEEK_TOP, saveData)
		}
	case for_game.WISH_ACT_H5_MONTH:
		dataType = for_game.WISH_ACTIVITY_DATA_TYPE_4
		// 清空月榜表数据
		_ = for_game.RemoveAllMonthTop()
		actList := for_game.GetWishPlayerActivityTop(for_game.WISH_ACTIVITY_DATA_TYPE_4, WISH_ACT_TOP_NUM)
		saveData := make([]interface{}, 0)
		for _, v := range actList {
			saveData = append(saveData, bson.M{"_id": v.GetPlayerId()}, &share_message.WishMonthTop{
				PlayerId:   easygo.NewInt64(v.GetPlayerId()),
				MonthPrize: easygo.NewInt64(v.GetMonthPrize()),
				CreateTime: easygo.NewInt64(time.Now().Unix()),
			})
		}
		// 插入新表
		if len(saveData) > 0 {
			for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_MONTH_TOP, saveData)
		}
	}

	// 更新上个周榜,月榜的数据.
	UpdateLastWeekMonth(dataType)
	//  插入排行
	giveWeekMonthPrize(t)
}

// 周榜月榜得到奖项 t 1-周,2-月
func giveWeekMonthPrize(t int) {
	logs.Info("周榜月榜获插入奖项记录开始,t %v", t)
	if t != for_game.WISH_ACT_H5_WEEK && t != for_game.WISH_ACT_H5_MONTH {
		return
	}
	var dataType int
	switch t {
	case for_game.WISH_ACT_H5_WEEK:
		dataType = for_game.WISH_ACTIVITY_DATA_TYPE_3
	case for_game.WISH_ACT_H5_MONTH:
		dataType = for_game.WISH_ACTIVITY_DATA_TYPE_4
	}
	// 查询规则列表
	ruleList := for_game.GetWishActPoolRuleListByType(int64(dataType))
	if len(ruleList) == 0 {
		return
	}

	switch dataType {
	case for_game.WISH_ACTIVITY_DATA_TYPE_3:
		weekTopList := for_game.GetWishWeekTop(len(ruleList))
		if len(weekTopList) == 0 {
			return
		}
		saveData := make([]interface{}, 0) // 批量修改活动用户身上的奖项id
		for i, v := range weekTopList {
			var nickName string
			if p := for_game.GetRedisWishPlayer(v.GetPlayerId()); p != nil {
				nickName = p.GetNickName()
			}
			actLog := for_game.NextId(for_game.TABLE_WISH_ACTIVITY_PRIZE_LOG)
			insertActLog1(actLog, v.GetPlayerId(), ruleList[i], nickName)
			v.WeekPrizeId = easygo.NewInt64(actLog)

			saveData = append(saveData, bson.M{"_id": v.GetPlayerId()}, v)
		}
		// 批量修改
		if len(saveData) > 0 {
			for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_WEEK_TOP, saveData)
		}

	case for_game.WISH_ACTIVITY_DATA_TYPE_4:
		monthTopList := for_game.GetWishMonthTop(len(ruleList))
		if len(monthTopList) == 0 {
			return
		}
		saveData := make([]interface{}, 0) // 批量修改活动用户身上的奖项id
		for i, v := range monthTopList {
			var nickName string
			if p := for_game.GetRedisWishPlayer(v.GetPlayerId()); p != nil {
				nickName = p.GetNickName()
			}
			actLog := for_game.NextId(for_game.TABLE_WISH_ACTIVITY_PRIZE_LOG)
			insertActLog1(actLog, v.GetPlayerId(), ruleList[i], nickName)
			v.MonthPrizeId = easygo.NewInt64(actLog)
			saveData = append(saveData, bson.M{"_id": v.GetPlayerId()}, v)
		}
		// 批量修改
		if len(saveData) > 0 {
			for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_MONTH_TOP, saveData)
		}

	default:
		logs.Error("giveWeekMonthPrize 类型有误", dataType)
		return
	}

}

// 奖池列表
func ActPoolListService() []*h5_wish.ActPool {
	list := for_game.GetWishActPoolList()
	result := make([]*h5_wish.ActPool, 0)
	for _, v := range list {
		result = append(result, &h5_wish.ActPool{
			Id:         easygo.NewInt64(v.GetId()),
			Name:       easygo.NewString(v.GetName()),
			BoxNum:     easygo.NewInt32(v.GetBoxNum()),
			CreateTime: easygo.NewInt64(v.GetCreateTime()),
			BoxIds:     v.GetBoxIds(),
		})
	}
	return result
}

// 奖池规则查询
func ActPoolRuleService(actPoolId int64, t int32) []*h5_wish.WishActPoolRule {
	list := for_game.GetWishActPoolRuleListByPoolId(actPoolId, t)
	data := make([]*h5_wish.WishActPoolRule, 0)
	for _, v := range list {
		data = append(data, &h5_wish.WishActPoolRule{
			Id:            easygo.NewInt64(v.GetId()),
			WishActPoolId: easygo.NewInt64(v.GetWishActPoolId()),
			Key:           easygo.NewInt32(v.GetKey()),
			Diamond:       easygo.NewInt64(v.GetDiamond()),
			WishItemId:    easygo.NewInt64(v.GetWishItemId()),
			AwardType:     easygo.NewInt32(v.GetAwardType()),
			Type:          easygo.NewInt32(v.GetType()),
		})
	}
	return data
}

// 累计天数活动查询
func SumDayService(uid, actPoolId int64, t int32) *h5_wish.SumNumResp {
	// 抽奖天数
	var signNum, lastDayTime int64
	// 根据奖池找到规则id
	id1 := for_game.GetWishActPoolRuleListByPoolId1(actPoolId)
	ruleIds := make([]int64, 0)
	for _, v := range id1 {
		ruleIds = append(ruleIds, v.GetId())
	}
	data := for_game.GeWishPlayerActivityByPId(uid)
	var dayLun int32
	if data != nil {
		for _, v := range data.GetData() {
			if !util.Int64InSlice(v.GetPoolRuleId(), ruleIds) {
				continue
			}
			if v.GetType() == t { // 天数
				signNum = v.GetValue()
				lastDayTime = v.GetUpdateTime()
				dayLun = v.GetDayLun()
				break
			}
		}
	}
	// 根据活动奖池id查询出规则列表
	list := for_game.GetWishActPoolRuleListByPoolId(actPoolId, t)
	// 遍历规则列表,查询出获奖记录
	dataList := make([]*h5_wish.WishActivityPrizeLog, 0)
	if t == 2 {
		// 判断是否相隔一天
		//tt := for_game.GetMillSecond() - lastDayTime
		lt := easygo.Get0ClockMillTimestamp(lastDayTime)
		timestamp := easygo.GetToday0ClockTimestamp()
		end := timestamp - lt
		if end < 24*3600000 {
			for _, v := range list {
				if logData := for_game.GetWishActivityPrizeLogByRuleIdAndLun(uid, v.GetId(), dayLun); logData != nil {
					dataList = append(dataList, &h5_wish.WishActivityPrizeLog{
						Id:                easygo.NewInt64(logData.GetId()),
						PlayerId:          easygo.NewInt64(logData.GetPlayerId()),
						Type:              easygo.NewInt32(logData.GetType()),
						ActType:           easygo.NewInt64(logData.GetActType()),
						WishActPoolRuleId: easygo.NewInt64(logData.GetWishActPoolRuleId()),
						Status:            easygo.NewInt32(logData.GetStatus()),
						CreateTime:        easygo.NewInt64(logData.GetCreateTime()),
					})
				}
			}
		} else {
			signNum = 0
		}
	} else {
		for _, v := range list {
			if logData := for_game.GetWishActivityPrizeLogByRuleId(uid, v.GetId()); logData != nil {
				dataList = append(dataList, &h5_wish.WishActivityPrizeLog{
					Id:                easygo.NewInt64(logData.GetId()),
					PlayerId:          easygo.NewInt64(logData.GetPlayerId()),
					Type:              easygo.NewInt32(logData.GetType()),
					ActType:           easygo.NewInt64(logData.GetActType()),
					WishActPoolRuleId: easygo.NewInt64(logData.GetWishActPoolRuleId()),
					Status:            easygo.NewInt32(logData.GetStatus()),
					CreateTime:        easygo.NewInt64(logData.GetCreateTime()),
				})
			}
		}
	}

	timeList := make([]int64, 0)
	dayList := for_game.GetWishDayActivityLogList(uid, actPoolId)
	for _, v := range dayList {
		timeList = append(timeList, v.GetCreateTime())
	}
	resp := &h5_wish.SumNumResp{
		SignNum:                  easygo.NewInt64(signNum),
		WishActivityPrizeLogList: dataList,
		LastDayTime:              easygo.NewInt64(lastDayTime),
		DayTimeList:              timeList,
	}
	return resp
}

func SumMoneyService1(uid int64, reqMsg *h5_wish.SumMoneyReq) easygo.IMessage {
	dataType := reqMsg.GetDataType()
	if dataType != for_game.WISH_ACT_H5_WEEK_NOW && dataType != for_game.WISH_ACT_H5_MONTH_NOW { // 3-现在周榜;4-现在月榜
		logs.Error("SumMoneyService 请求的类型有误3-现在周榜;4-现在月榜,dataType: %d", dataType)
		return easygo.NewFailMsg("类型有误.")
	}
	resp := &h5_wish.SumMoneyResp{}
	page, pageSize := for_game.MakePageAndPageSize(reqMsg.GetPage(), reqMsg.GetPageSize())
	data, count := for_game.GetWishPlayerActivityByPage(dataType, page, pageSize)
	if count >= 1000 { // 限制1000条
		count = 1000
	}
	if len(data) == 0 {
		logs.Error("GetWishPlayerActivityByPage 没有数据,dataType: %v, page: %v, pageSize: %v", dataType, page, pageSize)
		return resp
	}
	// 封装玩家头像昵称
	pIds := make([]int64, 0)
	for _, v := range data {
		pIds = append(pIds, v.GetPlayerId())
	}
	players, err := for_game.GetWishPlayerByIds(pIds)
	if err != nil {
		logs.Error("批量查询用户数据有误")
		return easygo.NewFailMsg("参数有误")
	}
	result := make([]*h5_wish.SumMoneyData, 0)
	for _, v := range data {
		for _, v1 := range players {
			if v.GetPlayerId() == v1.GetId() {
				var d int64
				if dataType == for_game.WISH_ACT_H5_WEEK_NOW {
					d = v.GetWeekPrize()
				} else if dataType == for_game.WISH_ACT_H5_MONTH_NOW {
					d = v.GetMonthPrize()
				}
				result = append(result, &h5_wish.SumMoneyData{
					NickName:         easygo.NewString(v1.GetNickName()),
					HeadIcon:         easygo.NewString(v1.GetHeadUrl()),
					ConSumDiamondNum: easygo.NewInt64(d),
				})
			}
		}
	}

	var key int // 规则表中的周榜,月榜的排名
	if page > 1 {
		key = (page - 1) * pageSize
	} else {
		key = page - 1
	}
	// dataType 1-周榜;2-月榜
	list := for_game.GetWishActPoolRuleListByTypeKey(dataType, key, pageSize)
	// 判断查询的数据,是否在获奖排名内
	k1 := key + 1 // 排名
	for i, v1 := range result {
		for _, v := range list {
			if v.GetKey() != int32(k1+i) {
				continue
			}
			awardType := v.GetAwardType()
			v1.AwardType = easygo.NewInt32(awardType)
			switch awardType { //1、钻石奖励 2、实物奖励
			case WISH_ACT_PRIZE_TYPE_DIAMOND:
				v1.GiveDiamondNum = easygo.NewInt64(v.GetDiamond())
			case WISH_ACT_PRIZE_TYPE_PRODUCT:
				wishItem := for_game.GetWishItemByIdFromDB(v.GetWishItemId())
				v1.ProductIcon = easygo.NewString(wishItem.GetIcon())
				v1.ProductId = easygo.NewString(v.GetWishItemId())
			}
		}

	}
	list1 := for_game.GetWishPlayerActivityByNum(WISH_ACT_TOP_NUM, dataType)
	// 判断自己是第几名
	var topNum int
	var diamondNum int64
	// 自己是第几名
	for i, v := range list1 {
		if v.GetPlayerId() == uid {
			topNum = i + 1
		}
	}
	p := for_game.GetWishPlayerActivityByPid(uid)
	if dataType == for_game.WISH_ACT_H5_WEEK_NOW {
		diamondNum = p.GetWeekPrize()
	} else if dataType == for_game.WISH_ACT_H5_MONTH_NOW {
		diamondNum = p.GetMonthPrize()
	}
	resp.DiamondNum = easygo.NewInt64(diamondNum)
	resp.TopNum = easygo.NewInt64(topNum)
	resp.SumMoneyDataList = result
	resp.TotalCount = easygo.NewInt64(count)
	return resp
}

// 修改状态
func UpdateTopStatus(dataType int64, page, pageSize int) {
	switch dataType {
	case for_game.WISH_ACT_H5_WEEK:
		data, _ := for_game.GetWishWeekTopByPage(page, pageSize)
		if len(data) == 0 {
			return
		}
		var saveData []interface{}
		for _, p := range data {
			if p.GetStatus() > 0 {
				continue
			}
			p.Status = easygo.NewInt64(1)
			saveData = append(saveData, bson.M{"_id": p.GetPlayerId()}, p)
		}
		if len(saveData) > 0 {
			for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_WEEK_TOP, saveData)
		}

	case for_game.WISH_ACT_H5_MONTH:
		data, _ := for_game.GetWishMonthTopByPage(page, pageSize)
		if len(data) == 0 {
			return
		}
		// 修改状态
		var saveData []interface{}
		for _, p := range data {
			if p.GetStatus() > 0 {
				continue
			}
			p.Status = easygo.NewInt64(1)
			saveData = append(saveData, bson.M{"_id": p.GetPlayerId()}, p)
		}
		if len(saveData) > 0 {
			for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_MONTH_TOP, saveData)
		}

	default:
		logs.Error("UpdateTopStatus类型不对")
		return
	}
}

func SumMoneyService(uid int64, reqMsg *h5_wish.SumMoneyReq) easygo.IMessage {
	resp := &h5_wish.SumMoneyResp{}
	t := reqMsg.GetDataType()
	page, pageSize := for_game.MakePageAndPageSize(reqMsg.GetPage(), reqMsg.GetPageSize())
	//var dataType int
	switch t {
	case for_game.WISH_ACT_H5_WEEK:
		//dataType = for_game.WISH_ACTIVITY_DATA_TYPE_3
		data, count := for_game.GetWishWeekTopByPage(page, pageSize)
		if len(data) == 0 {
			return resp
		}
		pIds := make([]int64, 0)
		for _, v := range data {
			pIds = append(pIds, v.GetPlayerId())
		}
		players, err := for_game.GetWishPlayerByIds(pIds)
		if err != nil {
			logs.Error("批量查询用户数据有误")
			return easygo.NewFailMsg("参数有误")
		}

		result := make([]*h5_wish.SumMoneyData, 0)
		prizeLogIds := make([]int64, 0) // 奖项id列表
		for _, v := range data {
			for _, v1 := range players {
				if v.GetPlayerId() == v1.GetId() {
					result = append(result, &h5_wish.SumMoneyData{
						PlayerId:         easygo.NewInt64(v.GetPlayerId()),
						NickName:         easygo.NewString(v1.GetNickName()),
						HeadIcon:         easygo.NewString(v1.GetHeadUrl()),
						ConSumDiamondNum: easygo.NewInt64(v.GetWeekPrize()),
					})
				}
			}
			prizeLogIds = append(prizeLogIds, v.GetWeekPrizeId())
		}
		if len(prizeLogIds) > 0 {
			prizeLogs := for_game.GeWishActivityPrizeLogByIds(prizeLogIds)
			for _, v := range prizeLogs {
				// 封装奖品信息
				for _, v1 := range result {
					if v1.GetPlayerId() != v.GetPlayerId() {
						continue
					}
					switch v.GetPrizeType() {
					case WISH_ACT_PRIZE_TYPE_DIAMOND:
						v1.GiveDiamondNum = easygo.NewInt64(v.GetPrizeValue())
					case WISH_ACT_PRIZE_TYPE_PRODUCT:
						wishItem := for_game.GetWishItemByIdFromDB(v.GetPrizeValue())
						v1.ProductIcon = easygo.NewString(wishItem.GetIcon())
						v1.ProductId = easygo.NewString(v.GetPrizeValue())
					}
					// 封装奖项id
					v1.PrizeLogId = easygo.NewInt64(v.GetId())
				}
			}
		}
		var topNum int
		var diamondNum int64
		// 自己是第几名
		list := for_game.GetWishWeekTop(WISH_ACT_TOP_NUM)
		for i, v := range list {
			if v.GetPlayerId() == uid {
				topNum = i + 1
				diamondNum = v.GetWeekPrize()
				if IsGtGetWeek12ClockTimestamp() {
					// 是否领奖
					if prizeLog := for_game.GetWishActivityPrizeLogById(v.GetWeekPrizeId()); prizeLog != nil {
						resp.Status = easygo.NewInt32(prizeLog.GetStatus())
					}
				} else {
					resp.Status = easygo.NewInt32(2)
				}

			}
		}
		if diamondNum == 0 {
			p := for_game.GetWishPlayerActivityByPid(uid)
			diamondNum = p.GetLastWeekDiamond()
		}
		resp.DiamondNum = easygo.NewInt64(diamondNum)
		resp.TopNum = easygo.NewInt64(topNum)
		resp.SumMoneyDataList = result
		resp.TotalCount = easygo.NewInt64(count)

	case for_game.WISH_ACT_H5_MONTH:
		data, count := for_game.GetWishMonthTopByPage(page, pageSize)
		if len(data) == 0 {
			return resp
		}
		pIds := make([]int64, 0)
		for _, v := range data {
			pIds = append(pIds, v.GetPlayerId())
		}
		players, err := for_game.GetWishPlayerByIds(pIds)
		if err != nil {
			logs.Error("批量查询用户数据有误")
			return easygo.NewFailMsg("参数有误")
		}

		result := make([]*h5_wish.SumMoneyData, 0)
		prizeLogIds := make([]int64, 0) // 奖项id列表
		for _, v := range data {
			for _, v1 := range players {
				if v.GetPlayerId() == v1.GetId() {
					result = append(result, &h5_wish.SumMoneyData{
						PlayerId:         easygo.NewInt64(v.GetPlayerId()),
						NickName:         easygo.NewString(v1.GetNickName()),
						HeadIcon:         easygo.NewString(v1.GetHeadUrl()),
						ConSumDiamondNum: easygo.NewInt64(v.GetMonthPrize()),
					})
				}
			}
			prizeLogIds = append(prizeLogIds, v.GetMonthPrizeId())
		}

		if len(prizeLogIds) > 0 {
			prizeLogs := for_game.GeWishActivityPrizeLogByIds(prizeLogIds)
			for _, v := range prizeLogs {
				// 封装奖品信息
				for _, v1 := range result {
					if v1.GetPlayerId() != v.GetPlayerId() {
						continue
					}
					switch v.GetPrizeType() {
					case WISH_ACT_PRIZE_TYPE_DIAMOND:
						v1.GiveDiamondNum = easygo.NewInt64(v.GetPrizeValue())
					case WISH_ACT_PRIZE_TYPE_PRODUCT:
						wishItem := for_game.GetWishItemByIdFromDB(v.GetPrizeValue())
						v1.ProductIcon = easygo.NewString(wishItem.GetIcon())
						v1.ProductId = easygo.NewString(v.GetPrizeValue())
					}
					// 封装奖项id
					v1.PrizeLogId = easygo.NewInt64(v.GetId())
				}
			}
		}

		var topNum int
		var diamondNum int64
		// 自己是第几名
		list := for_game.GetWishMonthTop(WISH_ACT_TOP_NUM)
		for i, v := range list {
			if v.GetPlayerId() == uid {
				topNum = i + 1
				diamondNum = v.GetMonthPrize()
				if IsGtGetMonth12ClockTimestamp() {
					// 是否领奖
					if prizeLog := for_game.GetWishActivityPrizeLogById(v.GetMonthPrizeId()); prizeLog != nil {
						resp.Status = easygo.NewInt32(prizeLog.GetStatus())
					}
				} else {
					resp.Status = easygo.NewInt32(2)
				}
			}
		}
		if diamondNum == 0 {
			p := for_game.GetWishPlayerActivityByPid(uid)
			diamondNum = p.GetLastMonthDiamond()
		}
		resp.DiamondNum = easygo.NewInt64(diamondNum)
		resp.TopNum = easygo.NewInt64(topNum)
		resp.SumMoneyDataList = result
		resp.TotalCount = easygo.NewInt64(count)

	default:
		//logs.Error("SumMoneyService 请求的类型有误1-周榜;2-月榜,dataType: %d", dataType)
		return easygo.NewFailMsg("类型有误.")
	}
	return resp
}

// 领奖活动
func GiveService(pid, prizeLogId int64) easygo.IMessage {
	resp := &h5_wish.GiveResp{
		Result: easygo.NewInt32(2),
	}
	// 判断是否存在奖项
	prizeLog := for_game.GetWishActivityPrizeLogById(prizeLogId)
	if prizeLog == nil {
		logs.Error("根据id查询奖项失败,prizeLogId", prizeLogId)
		return resp
	}
	// 判断是否是总金额的奖项
	if prizeLog.GetType() == for_game.WISH_ACTIVITY_DATA_TYPE_3 {
		// 判断时间是否大于周一中午12点
		if !IsGtGetWeek12ClockTimestamp() {
			logs.Error("时间不到本周一中午12点,当前时间为: %d", time.Now().Unix())
			return resp
		}
	}
	if prizeLog.GetType() == for_game.WISH_ACTIVITY_DATA_TYPE_4 {
		// 判断时间是否大于本月1号中午12点
		if !IsGtGetMonth12ClockTimestamp() {
			logs.Error("时间不到本月中午12点,当前时间为: %d", time.Now().Unix())
			return resp
		}
	}
	if prizeLog.GetStatus() == for_game.WISH_ACTIVITY_PRIZE_STATUS_1 {
		logs.Error("领奖接口,该奖项已被领取,prizeLogId", prizeLogId)
		return resp
	}
	if prizeLog.GetPlayerId() != pid {
		logs.Error("奖项跟用户不匹配,uid: %v,prizeLogId: %v", pid, prizeLogId)
		return resp
	}
	var playerDiamond int64
	// 判断是奖品类型
	switch prizeLog.GetPrizeType() {
	case WISH_ACT_PRIZE_TYPE_DIAMOND:
		diamond := prizeLog.GetPrizeValue()
		if diamond <= 0 {
			logs.Error("奖项的钻石数量为0,diamond: %d", diamond)
			return resp
		}
		// 修改用户钻石数
		base := for_game.GetRedisWishPlayer(pid)
		if base == nil {
			logs.Error("获取许愿池用户对象失败,许愿池用户id为: %d", pid)
			return resp
		}
		imPid := base.GetPlayerId()
		extendLog := &share_message.GoldExtendLog{
			PlayerId: easygo.NewInt64(imPid),
		}

		err11, _ := base.AddDiamond(diamond, "许愿池活动所得", for_game.DIAMOND_TYPE_WISH_ACT, extendLog)
		if err11 != nil {
			logs.Error("许愿池活动领奖,添加钻石失败,err: %v", err11)
			return resp
		}
		playerDiamond = base.GetDiamond()
	case WISH_ACT_PRIZE_TYPE_PRODUCT:
		// 在许愿池物品表中添加用户得到你的物品
		wishItem := for_game.GetWishItemByIdFromDB(prizeLog.GetPrizeValue())
		if wishItem.GetPrice() == 0 {
			marshal, _ := json.Marshal(wishItem)
			logs.Error("许愿池活动的物品信息有误,wishItem: %s", string(marshal))
		}
		playerWishItem := &share_message.PlayerWishItem{
			Id:              easygo.NewInt64(for_game.NextId(for_game.TABLE_PLAYER_WISH_ITEM)),
			PlayerId:        easygo.NewInt64(pid),
			Status:          easygo.NewInt32(0), // 带兑换
			IsRead:          easygo.NewBool(false),
			WishItemId:      easygo.NewInt64(prizeLog.GetPrizeValue()),
			WishItemPrice:   easygo.NewInt64(wishItem.GetPrice()),
			ProductName:     easygo.NewString(wishItem.GetName()),
			WishItemDiamond: easygo.NewInt64(wishItem.GetDiamond()),
			WishItemIcon:    easygo.NewString(wishItem.GetIcon()),
			GiveType:        easygo.NewInt64(1),
			CreateTimeMill:  easygo.NewInt64(for_game.GetMillSecond()),
		}
		if err := for_game.AddPlayerWishItem(playerWishItem); err != nil { //添加玩家物品
			return resp
		}

	default:
		logs.Error("奖项类型有误,PrizeType", prizeLog.GetPrizeType())
		return resp
	}

	// 修改领奖状态领奖时间
	prizeLog.Status = easygo.NewInt32(for_game.WISH_ACTIVITY_PRIZE_STATUS_1)
	prizeLog.FinishTime = easygo.NewInt64(for_game.GetMillSecond())
	if err := for_game.UpsertWishActivityPrizeLog(prizeLog); err != nil {
		marshal, _ := json.Marshal(prizeLog)
		logs.Error("许愿池活动领奖,修改领奖状态失败,需要人工处理,prizeLog: %s", string(marshal))
		return resp
	}

	if playerDiamond > 0 {
		resp.Diamond = easygo.NewInt64(playerDiamond)
	}

	resp.Result = easygo.NewInt32(1)
	return resp
}

// 批量修改上个周期周榜月榜的数据
func UpdateLastWeekMonth(t int) {
	if t != for_game.WISH_ACTIVITY_DATA_TYPE_3 && t != for_game.WISH_ACTIVITY_DATA_TYPE_4 {
		return
	}
	// 5000条5000条的查询
	var maxId int64 = 18801000
	//增加索引
	for {
		actList := for_game.GetWishPlayerActivityByGTPid(maxId)
		logs.Info("批量处理许愿池活动用户的排行榜:", maxId, len(actList))
		if len(actList) == 0 {
			break
		}
		var saveData []interface{}
		for _, p := range actList {
			if t == for_game.WISH_ACTIVITY_DATA_TYPE_3 {
				p.LastWeekDiamond = easygo.NewInt64(p.GetWeekPrize())
				p.WeekPrize = easygo.NewInt64(0)
			} else if t == for_game.WISH_ACTIVITY_DATA_TYPE_4 {
				p.LastMonthDiamond = easygo.NewInt64(p.GetMonthPrize())
				p.MonthPrize = easygo.NewInt64(0)
			}

			saveData = append(saveData, bson.M{"_id": p.GetPlayerId()}, p)
		}
		if len(saveData) > 0 {
			for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_WISH_PLAYER_ACTIVITY, saveData)
		}
		maxId = actList[len(actList)-1].GetPlayerId()
		logs.Info("批量处理许愿池活动用户的排行榜完成:", maxId, len(actList))
	}
}

func ActNameService(t int32) string {
	act := for_game.GetActivityByType(t)
	if act == nil {
		logs.Error("获取活动数据失败,t: %d", t)
		return ""
	}
	return act.GetTitle()
}

// 当前时间是否大于本月开始12点
func IsGtGetMonth12ClockTimestamp() bool {
	var ok bool
	t := time.Now().Unix()
	if t > easygo.GetMonth12ClockTimestamp() {
		ok = true
	}
	return ok
}

// 当前时间是否大于本周开始12点
func IsGtGetWeek12ClockTimestamp() bool {
	var ok bool
	t := time.Now().Unix()
	if t > easygo.GetWeek12ClockTimestamp() {
		ok = true
	}
	return ok
}

// 后台工具抽奖
func BackstageDareToolService(boxId, uid, poolId, count, diamond, userId int64) int {

	if count == 0 || poolId == 0 || diamond == 0 || boxId == 0 {
		logs.Error("BackstageDareToolService 0")
		return 0
	}
	var dareCount int
	for i := 0; i < int(count); i++ {
		if err := BackDareTool(uid, boxId, poolId, diamond); err != nil {
			logs.Error("BackstageDareToolService err: %s", err.Error())
			break
		} else {
			dareCount++
			easygo.Spawn(SendMsgToIdelServer, for_game.SERVER_TYPE_BACKSTAGE, "RpcWishToolLucky", &h5_wish.BackstageDareToolReq{UserId: easygo.NewInt64(userId)})
		}
	}
	logs.Info("应该抽奖 %d次,实际抽奖了: %d", count, dareCount)
	return dareCount
}

// 抽奖工具核心逻辑
func BackDareTool(uid, boxId, poolId, dm int64) error {
	logs.Info("BackDareTool抽奖工具核心逻辑开始")
	dStatus := &DareStatus{}
	poolObj := GetPoolObj(poolId)
	//  inc 挑战的金额,返回当前的水池的余额,判断是否符合抽水状态,符合就立马抽水,记录抽水日志.返回的值根据当前水池余额得到水池状态
	var pool *share_message.WishPool
	var pollPrice, recycle, commission int64
	var err error

	boxPrice := dm // 抽奖砖石数量
	beforePool := poolObj.GetPollInfoFromRedis()
	pollPrice, err = poolObj.UpdatePollPriceToRedis(poolId, boxPrice, "IncomeValue")
	if err != nil {
		logs.Error("后台抽奖工具 抽奖失败原因,更改水池价格失败,用户id为: %d, 水池id为: %d,err: %s", uid, poolId, err.Error())
		return err
	}
	//  获取并修改水池状态
	//toolSetPoolStatus(poolId)
	pool = poolObj.GetPollInfoFromRedis()
	if pool == nil {
		logs.Error("后台抽奖工具 抽奖失败原因,从redis中获取水池数据失败,用户id为: %d,水池id为: %d", uid, poolId)
		return errors.New("后台抽奖工具  水池为空")
	}

	// 记录水池的流水
	wishPoolLog := &share_message.WishPoolLog{
		PoolId:           easygo.NewInt64(poolId),
		PlayerId:         easygo.NewInt64(uid),
		BeforeValue:      easygo.NewInt64(pollPrice - boxPrice),
		AfterValue:       easygo.NewInt64(pollPrice),
		Value:            easygo.NewInt64(boxPrice),
		Type:             easygo.NewInt64(for_game.WISH_POOL_LOG_TYPE_2), //2-挑战增加
		LocalStatus:      easygo.NewInt64(beforePool.GetLocalStatus()),
		IsOpenAward:      easygo.NewBool(beforePool.GetIsOpenAward()),
		AfterLocalStatus: easygo.NewInt64(pool.GetLocalStatus()),
		AfterIsOpenAward: easygo.NewBool(pool.GetIsOpenAward()),
		CreateTimeMill:   easygo.NewInt64(for_game.GetMillSecond()),
		BoxId:            easygo.NewInt64(boxId),
	}
	easygo.Spawn(for_game.AddToolWishPoolLog, wishPoolLog)

	// 判断金额是否超过抽水的阀值
	recycle = pool.GetRecycle() // 回收阀值(存库)
	logs.Info("后台抽奖工具 投钱进水池后的金额为:%d,用户id为:%d, 抽水的阀值为: %d", pollPrice, uid, recycle)
	commission = pool.GetCommission() // 抽水金额
	pool = poolObj.GetPollInfoFromRedis()
	pollPrice = pool.GetIncomeValue()
	if pollPrice >= recycle { // 符合就立马抽水
		poolObj.Mutex1.Lock()
		pool = poolObj.GetPollInfoFromRedis()
		pollPrice = pool.GetIncomeValue()
		if pollPrice >= recycle {
			// 抽水
			p, err := toolPump(boxId, poolId, commission, uid, recycle, beforePool.GetLocalStatus(), beforePool.GetIsOpenAward(), poolObj, pool)
			if err != nil {
				logs.Error("后台抽奖工具 抽奖失败原因,抽水公共方法失败了,err: %s", err.Error())
				return err
			}
			pollPrice = p
		}
		poolObj.Mutex1.Unlock()
	}

	pool = poolObj.GetPollInfoFromRedis()
	pollPrice = pool.GetIncomeValue()
	logs.Info("后台抽奖工具 计算水池状态时当前水池水量为: %d,水池id为: %d", pollPrice, poolId)
	f := easygo.Decimal(easygo.AtoFloat64(easygo.AnytoA(pollPrice))/easygo.AtoFloat64(easygo.AnytoA(pool.GetPoolLimit())), 2)
	poolStatusLimitFloat64 := f * 100
	poolStatusLimit := int64(easygo.Decimal(poolStatusLimitFloat64, 0))
	logs.Info("后台抽奖工具 水池id为:%d,水池的金额为:%d, 计算状态前的百分数: %d", poolId, pollPrice, poolStatusLimit)
	// 判断是在哪个状态
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
	if bigLoss != nil && (poolStatusLimit <= bigLoss.GetShowMinValue()) {
		poolStatus = for_game.POOL_STATUS_BIGLOSS
	}
	logs.Info("后台抽奖工具 1-大亏;2-小亏;3-普通;4-大盈;5-小盈,水池的状态为: %v", poolStatus)
	// 判断是否需要附加权重
	isOpenAward := pool.GetIsOpenAward()
	// 根据盲盒状态找出物品
	wishBoxItems := for_game.GetToolProductByPoolStatus(boxId, poolStatus, isOpenAward) // 中间表的
	if len(wishBoxItems) == 0 {
		logs.Error("后台抽奖工具 抽奖失败原因,没有当前状态的物品pid: %d, poolStatus: %d", uid, poolStatus)
		return errors.New("后台抽奖工具 没有当前状态的物品")
	}

	// 权重的map[权重]中间表物品的id
	indexWeightMap, sortWeightList := getWeightByStatus(isOpenAward, poolStatus, wishBoxItems, 0)
	if len(indexWeightMap) == 0 {
		logs.Error("后台抽奖工具 抽奖失败原因,getWeightByStatus()返回的 indexWeightMap 长度为0")
		return errors.New("后台抽奖工具 indexWeightMap 长度为0")
	}
	if len(sortWeightList) == 0 {
		logs.Error("后台抽奖工具 抽奖失败原因,getWeightByStatus()返回的 sortWeightList 长度为0")
		return errors.New("后台抽奖工具 sortWeightList 长度为0")
	}
	// 计算权重,抽出物品
	boxItemId := GetProductByWeight(indexWeightMap, sortWeightList, isOpenAward, 0)
	if boxItemId == 0 {
		logs.Error("后台抽奖工具 抽奖失败原因,抽奖失败得到的物品id为0,数据有误,用户id为: %d,indexWeightMap: %v,sortWeightList: %v", uid, indexWeightMap, sortWeightList)
		return errors.New("后台抽奖工具 计算权重,抽出物品失败")
	}
	// 从数据库中查找
	item := for_game.GetToolWishBoxItemByIdFromDB(boxItemId)
	if item == nil {
		logs.Error("后台抽奖工具 抽奖失败原因,获取wishBoxItem 为nil,用户id为: %d,_id为: %d", uid, boxItemId)
		return errors.New("后台抽奖工具 得到物品失败")
	}

	rewardLv := item.GetRewardLv() //1、小奖 2、大奖
	var afterPrice int64
	var wishItemId int64 // 中间表的主键id

	diamond := item.GetDiamond()
	logs.Info("后台抽奖工具 物品id(中间表的): %d 的钻石价格为: %d", boxItemId, diamond)
	if rewardLv == for_game.ITEM_REWARDLV_2 {
		//如果是贵重物品,先上锁取水池的值(锁住当前水池不许更改.),能不能开大奖(当前盈利是否都大于该大奖物品的价格),如果属于,就开大奖,否则随机开小奖,更新水池放奖后的水量,解锁
		pp, isBigReward := poolObj.GetToolReward(poolId, diamond, "IncomeValue")
		if !isBigReward { // 不能开大奖,随机抽小奖
			// 随机抽小奖的物品.
			items := for_game.FindToolMaxPriceBoxItemByBoxIdAndLv(boxId, for_game.ITEM_REWARDLV_1)
			if len(items) == 0 {
				logs.Error("后台抽奖工具 抽奖失败原因,不能开大奖,也没有普通的物品")
				return errors.New("后台抽奖工具 不能开大奖,也没有普通的物品")
			}
			r := for_game.RandInt(0, len(items))
			afterItem := items[r]

			diamond1 := afterItem.GetDiamond()
			logs.Info("后台抽奖工具 物品id(中间表的): %d 的钻石价格为: %d", afterItem.GetId(), diamond1)
			afterPrice = diamond1
			wishItemId = afterItem.GetId()
			pollPrice, err = poolObj.UpdatePollPriceToRedis(poolId, 0-afterPrice, "IncomeValue")
			logs.Info("后台抽奖工具 中了大奖但只能开小奖,物品价格为: %v,扣除物品后的水量为: %v", afterPrice, pollPrice)
			if err != nil {
				logs.Error("后台抽奖工具 抽奖失败原因,发奖后更改水池价格失败,用户id为: %d, 水池id为: %d,err: %s", uid, poolId, err.Error())
				return err
			}
			//  获取并修改水池状态
			//toolSetPoolStatus(poolId)

		} else { // 大奖已经扣了大奖的的价格了.
			afterPrice = diamond
			wishItemId = item.GetId()
			if pp > 0 {
				pollPrice = pp
			}
			logs.Info("后台抽奖工具 中了大奖,用户id为: %d,物品id为: %d,扣除大奖物品后水池水量为: %d", uid, wishItemId, pollPrice)

		}
	} else {
		afterPrice = diamond
		wishItemId = item.GetId()
		// 从水池中减去物品的价格
		pollPrice, err = poolObj.UpdatePollPriceToRedis(poolId, 0-afterPrice, "IncomeValue")
		logs.Info("后台抽奖工具 开小奖,物品价格为: %v,扣除物品后的水量为: %v", afterPrice, pollPrice)
		if err != nil {
			logs.Error("后台抽奖工具 抽奖失败原因,发奖后更改水池价格失败,用户id为: %d, 水池id为: %d,err: %s", uid, poolId, err.Error())
			return err
		}
		//  获取并修改水池状态
		//toolSetPoolStatus(poolId)
	}

	pool = poolObj.GetPollInfoFromRedis()
	var wishPoolLog1BeforeValue int64
	// 得出物品后 记录水池的流水
	wishPoolLog1 := &share_message.WishPoolLog{
		PoolId:           easygo.NewInt64(poolId),
		PlayerId:         easygo.NewInt64(uid),
		BeforeValue:      easygo.NewInt64(pollPrice + afterPrice),
		AfterValue:       easygo.NewInt64(pollPrice),
		Value:            easygo.NewInt64(0 - afterPrice),
		Type:             easygo.NewInt64(for_game.WISH_POOL_LOG_TYPE_3), // 3-得到物品后扣除
		LocalStatus:      easygo.NewInt64(poolStatus),
		IsOpenAward:      easygo.NewBool(isOpenAward),
		AfterLocalStatus: easygo.NewInt64(pool.GetLocalStatus()),
		AfterIsOpenAward: easygo.NewBool(pool.GetIsOpenAward()),
		CreateTimeMill:   easygo.NewInt64(for_game.GetMillSecond()),
		BoxId:            easygo.NewInt64(boxId),
	}
	wishPoolLog1BeforeValue = wishPoolLog1.GetBeforeValue()
	easygo.Spawn(for_game.AddToolWishPoolLog, wishPoolLog1)

	dStatus.AfterIsOpenAward = pool.GetIsOpenAward()
	dStatus.AfterLocalStatus = pool.GetLocalStatus()
	dStatus.IsOpenAward = isOpenAward
	dStatus.LocalStatus = int64(poolStatus)
	dStatus.PoolIncomeValue = wishPoolLog1BeforeValue
	dStatus.AfterPoolIncomeValue = pollPrice

	//   修改是是否需要附加权重.
	isOpen := checkIsOpenAward(pool.GetIncomeValue(), pool.GetStartAward(), pool.GetCloseAward(), pool.GetIsOpenAward())
	poolObj.SetIsOpenAward(poolId, isOpen)
	logs.Info("后台抽奖工具 用户id: %d, 抽到的物品id为: %d,扣除物品后水池的水量为: %v", uid, wishItemId, pollPrice)

	// 再次校验是否要抽水
	if pool.GetIncomeValue() >= recycle {
		poolObj.Mutex1.Lock()
		pool = poolObj.GetPollInfoFromRedis()
		if pool.GetIncomeValue() >= recycle {
			if _, err := toolPump(boxId, poolId, commission, uid, recycle, int64(poolStatus), isOpenAward, poolObj, pool); err != nil {
				logs.Error("后台抽奖工具 再次校验抽水,抽水公共方法失败了,err: %s", err.Error())
			}

		}
		poolObj.Mutex1.Unlock()
	}

	// 记录用户得到的物品
	if wishItemId == 0 {
		logs.Error("后台抽奖工具 抽奖得到的物品id为0,需要人工排查,userId= %d", uid)
		return errors.New("后台抽奖工具 抽奖得到的物品id为0")
	}

	//  获取并修改水池状态
	toolSetPoolStatus(poolId)

	wishBoxItem, _ := for_game.GetToolWishBoxItem(wishItemId) // 中奖信息

	id := for_game.NextId(for_game.TABLE_TOOL_PLAYER_WISH_ITEM)
	playerWishItem := &share_message.PlayerWishItem{
		Id:               easygo.NewInt64(id),
		PlayerId:         easygo.NewInt64(uid),
		ChallengeItemId:  easygo.NewInt64(wishBoxItem.GetId()),
		WishItemPrice:    easygo.NewInt64(wishBoxItem.GetPrice()), // 商品价格人民币
		DareDiamond:      easygo.NewInt64(dm),                     // 挑战时候的钻石
		ProductName:      easygo.NewString(wishBoxItem.GetWishItemName()),
		WishItemDiamond:  easygo.NewInt64(wishBoxItem.GetDiamond()), // 物品钻石数
		WishItemStyle:    easygo.NewInt32(wishBoxItem.GetStyle()),
		AfterIsOpenAward: easygo.NewBool(dStatus.AfterIsOpenAward),  // 抽奖后是否开启放奖 true 放奖 false 关闭
		AfterLocalStatus: easygo.NewInt64(dStatus.AfterLocalStatus), //  抽奖后当前水池的状态;1-大亏;2-小亏;3-普通;4-大赢;5-小盈
		IsOpenAward:      easygo.NewBool(dStatus.IsOpenAward),       // 是否开启放奖 true 放奖 false 关闭
		LocalStatus:      easygo.NewInt64(dStatus.LocalStatus),      // 当前水池的状态;1-大亏;2-小亏;3-普通;4-大赢;5-小盈
		CreateTimeMill:   easygo.NewInt64(for_game.GetMillSecond()),
		WishBoxId:        easygo.NewInt64(boxId),
	}
	err = for_game.AddToolPlayerWishItem(playerWishItem) //添加玩家物品
	easygo.PanicError(err)
	return nil
}

// 修改水池状态
func toolSetPoolStatus(poolId int64) {
	status := GetPoolStatus(poolId, 0)
	poolObj1 := GetPoolObj(poolId)
	poolObj1.SetLocalStatus(poolId, int64(status))
}

// 校验账号是否冻结. true-为冻结,false-不冻结.
func CheckIsFreeze(pid int64) bool {
	p := for_game.GetRedisWishPlayer(pid)
	if p == nil {
		return true
	}
	if !p.GetIsFreeze() {
		return false
	}
	// 判断是否到期
	if p.GetFreezeTime() == 1 { // 永久冻结
		return true
	}
	if time.Now().Unix() > p.GetFreezeTime() {
		// 设置false
		p.SetIsFreeze(false)
		return false
	}
	return true
}
