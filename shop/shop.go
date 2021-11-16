package shop

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/share_message"
	"sort"
	"strings"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

type ShopItemList []*ShopItem
type ShopItemMap = map[ITEM_ID]*ShopItem

type LastShopItemsNew []*ShopItem
type LastShopItemsPriceDesc []*ShopItem
type LastShopItemsPriceAsc []*ShopItem
type LastShopItemsSalesDesc []*ShopItem
type LastShopItemsStoreDesc []*ShopItem
type LastShopItemsComDesc []*ShopItem
type ExpressBodyList []*share_message.QueryExpressBody

//超时半小时取消
var ORDER_EXPIRE_TIME int64 = 1800

type ItemFile struct {
	file_url    string
	file_type   int32
	file_width  string
	file_height string
}

type ShopItem struct {
	item_id             ITEM_ID
	price               int32
	origin_price        int32
	title               string
	name                string
	item_files          []ItemFile
	player_id           int64
	nickname            string
	avatar              string
	account             string
	userName            string
	phone               string
	address             string
	detail_address      string
	item_type           int32
	other_type          []string
	create_time         int64
	stock_count         int32
	sex                 int32
	state               int32
	realPayCnt          int32
	fakePayCnt          int32
	realPageViews       int32
	fakePageViews       int32
	realGoodCommCnt     int32
	realFinCommCnt      int32
	fakeGoodCommCnt     int32
	fakeFinCommCnt      int32
	fakeFixGoodCommRate int32
	realCommentCnt      int32
	realStoreCnt        int32
	realFinOrderCnt     int32
	pointCardName       string
}

type Shop struct {
	shop_items       *ShopItemList
	shop_items_by_id *ShopItemMap
}

type ReceiveAddress struct {
	name           string
	phone          string
	region         string
	detail_address string
}

func NewShop() *Shop {
	p := &Shop{}
	p.Init()
	return p
}
func (self *Shop) Init() {
	//初始化redis订单生成的主键id
	//redis分布式key的初始值
	self.InitCreateOrderId()

	var shop_item_list ShopItemList = ShopItemList{}
	self.shop_items = &shop_item_list

	var shop_item_map ShopItemMap = ShopItemMap{}
	self.shop_items_by_id = &shop_item_map

	//改用定时器
	//加载上架商品(第一次启动立即执行)
	easygo.Spawn(self.UpdateItemList)
	//自动收货(第一次启动立即执行)
	easygo.Spawn(self.UpdateReceiveOrder)
}

func (self *Shop) CreateOrderID() int64 {
	b, err := easygo.RedisMgr.GetC().Exist(for_game.SHOP_CREATE_ORDER_ID)
	easygo.PanicError(err)
	if !b {
		self.InitCreateOrderId()
	}
	return easygo.RedisMgr.GetC().StringIncrForInt64(for_game.SHOP_CREATE_ORDER_ID)
}

func (self *Shop) Recommend(
	page int32,
	pageSize int32,
	itemType []int32,
	playerId int64,
	cacheItemTypes []int32,
	cacheSearch []string) *share_message.ItemList {

	var list []*share_message.ShopItem = []*share_message.ShopItem{}
	shopItems := self.shop_items
	sortShopItems := ShopItemList{}

	for _, value := range *shopItems {
		for j, _ := range itemType {
			if (value.item_type == itemType[j] || itemType[j] == 0) && value.stock_count > 0 {
				sortShopItems = append(sortShopItems, value)
				break
			}
		}
	}

	midShopItems := ShopItemList{}

	var blackList []PLAYER_ID = self.GetBlackLists(playerId)

	var addFlag int32 = 0
	for _, value := range sortShopItems {
		for _, black := range blackList {
			if value.player_id == black {
				addFlag = 1
				break
			}
		}
		if addFlag == 0 {
			midShopItems = append(midShopItems, value)
		}

		addFlag = 0
	}

	//排序
	sort.Sort(midShopItems)

	// 过滤出浏览的过的缓存商品和未浏览过的商品
	viewShopItems := ShopItemList{}
	notViewShopItems := ShopItemList{}

	for _, value := range midShopItems {
		var viewFlag bool = false
		if nil != cacheItemTypes && len(cacheItemTypes) > 0 {
			for j, _ := range cacheItemTypes {
				if value.item_type == cacheItemTypes[j] {
					viewShopItems = append(viewShopItems, value)
					viewFlag = true
					break
				}
			}
		}

		if !viewFlag {
			notViewShopItems = append(notViewShopItems, value)
		}
	}

	// TODO 把浏览过的商品按照搜索记录进行排序
	lastShopItems := ShopItemList{}
	lastShopItems = append(lastShopItems, viewShopItems...)
	lastShopItems = append(lastShopItems, notViewShopItems...)

	var count = int32(len(lastShopItems))
	if page*pageSize <= count {

		for i := page * pageSize; i < pageSize*(page+1) && i < count; i++ {
			itemId := lastShopItems[i].item_id
			var itemFile share_message.ItemFile = share_message.ItemFile{
				FileUrl:    &lastShopItems[i].item_files[0].file_url,
				FileType:   &lastShopItems[i].item_files[0].file_type,
				FileWidth:  &lastShopItems[i].item_files[0].file_width,
				FileHeight: &lastShopItems[i].item_files[0].file_height,
			}

			var item share_message.ShopItem = share_message.ShopItem{
				ItemId:     easygo.NewInt64(itemId),
				Price:      easygo.NewInt32(lastShopItems[i].price),
				Title:      easygo.NewString(lastShopItems[i].title),
				ItemFile:   &itemFile,
				StoreCount: easygo.NewInt32(lastShopItems[i].realStoreCnt),
				PlayerId:   easygo.NewInt64(lastShopItems[i].player_id),
				Nickname:   easygo.NewString(lastShopItems[i].nickname),
				Avatar:     easygo.NewString(lastShopItems[i].avatar),
				Account:    easygo.NewString(lastShopItems[i].account),
				Sex:        easygo.NewInt32(lastShopItems[i].sex),
				Name:       easygo.NewString(lastShopItems[i].name),
				State:      easygo.NewInt32(lastShopItems[i].state),
				ItemType:   easygo.NewInt32(lastShopItems[i].item_type),
			}
			list = append(list, &item)
		}

		return &share_message.ItemList{
			Items:    list,
			PageSize: easygo.NewInt32(pageSize),
			Page:     easygo.NewInt32(page),
			Count:    easygo.NewInt32(count)}
	}

	return &share_message.ItemList{
		Items:    list,
		PageSize: easygo.NewInt32(pageSize),
		Page:     easygo.NewInt32(page),
		Count:    easygo.NewInt32(count)}
}

func (self *Shop) Show(
	page int32,
	pageSize int32,
	itemType []int32,
	playerId int64,
	cacheItemTypes []int32) *share_message.ItemList {

	var list []*share_message.ShopItem = []*share_message.ShopItem{}

	shopItems := self.shop_items
	sortShopItems := ShopItemList{}

	for _, value := range *shopItems {
		for j, _ := range itemType {
			if (value.item_type == itemType[j] || itemType[j] == 0) && value.stock_count > 0 {
				sortShopItems = append(sortShopItems, value)
				break
			}
		}
	}
	lastShopItems := ShopItemList{}

	var blackList []PLAYER_ID = self.GetBlackLists(playerId)

	var addFlag int32 = 0
	for _, value := range sortShopItems {
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

	//排序
	sort.Sort(lastShopItems)

	var count = int32(len(lastShopItems))
	if page*pageSize <= count {

		for i := page * pageSize; i < pageSize*(page+1) && i < count; i++ {
			itemId := lastShopItems[i].item_id
			var itemFile share_message.ItemFile = share_message.ItemFile{
				FileUrl:    &lastShopItems[i].item_files[0].file_url,
				FileType:   &lastShopItems[i].item_files[0].file_type,
				FileWidth:  &lastShopItems[i].item_files[0].file_width,
				FileHeight: &lastShopItems[i].item_files[0].file_height,
			}

			var item share_message.ShopItem = share_message.ShopItem{
				ItemId:     easygo.NewInt64(itemId),
				Price:      easygo.NewInt32(lastShopItems[i].price),
				Title:      easygo.NewString(lastShopItems[i].title),
				ItemFile:   &itemFile,
				StoreCount: easygo.NewInt32(lastShopItems[i].realStoreCnt),
				PlayerId:   easygo.NewInt64(lastShopItems[i].player_id),
				Nickname:   easygo.NewString(lastShopItems[i].nickname),
				Avatar:     easygo.NewString(lastShopItems[i].avatar),
				Account:    easygo.NewString(lastShopItems[i].account),
				Sex:        easygo.NewInt32(lastShopItems[i].sex),
				Name:       easygo.NewString(lastShopItems[i].name),
				State:      easygo.NewInt32(lastShopItems[i].state),
				ItemType:   easygo.NewInt32(lastShopItems[i].item_type),
				CopyName:   easygo.NewString(lastShopItems[i].name),
			}
			//如果是点卡的时候重新设置商品名称
			if lastShopItems[i].item_type == for_game.SHOP_POINT_CARD_CATEGORY {
				item.Name = easygo.NewString(lastShopItems[i].pointCardName)
			}

			list = append(list, &item)
		}

		//打标list中的你可能喜欢
		if nil != cacheItemTypes && len(cacheItemTypes) > 0 {
			for _, value := range list {
				for j, _ := range cacheItemTypes {
					if value.GetItemType() == cacheItemTypes[j] {
						value.MayEnjoy = easygo.NewBool(true)
						break
					}
				}
			}
		}

		return &share_message.ItemList{
			Items:    list,
			PageSize: easygo.NewInt32(pageSize),
			Page:     easygo.NewInt32(page),
			Count:    easygo.NewInt32(count)}
	}

	return &share_message.ItemList{
		Items:    list,
		PageSize: easygo.NewInt32(pageSize),
		Page:     easygo.NewInt32(page),
		Count:    easygo.NewInt32(count)}
}

func (self *Shop) GetItemFromCache(itemId ITEM_ID) *ShopItem {
	shop_items_by_id := self.shop_items_by_id
	return (*shop_items_by_id)[itemId]
}

func (self *Shop) ItemDetail(shopItem *ShopItem) *share_message.ShopItemDetail {
	if shopItem == nil {
		return nil
	}

	var types share_message.ShopItemType = share_message.ShopItemType{
		Type:      &(*shopItem).item_type,
		OtherType: (*shopItem).other_type,
	}

	var itemFiles = self.ItemFile(shopItem)
	var stockCount = shopItem.stock_count
	var typess int32
	if p := for_game.GetRedisPlayerBase(shopItem.player_id); p != nil {
		typess = p.GetTypes()
	}
	var itemDetail share_message.ShopItemDetail = share_message.ShopItemDetail{
		ItemId:        easygo.NewInt64(shopItem.item_id),
		Price:         easygo.NewInt32(shopItem.price),
		ItemFiles:     itemFiles,
		Title:         easygo.NewString(shopItem.title),
		PlayerId:      easygo.NewInt64(shopItem.player_id),
		Nickname:      easygo.NewString(shopItem.nickname),
		Avatar:        easygo.NewString(shopItem.avatar),
		StoreCount:    easygo.NewInt32(shopItem.realStoreCnt),
		CreateTime:    easygo.NewInt64(shopItem.create_time),
		StockCount:    easygo.NewInt32(stockCount),
		Type:          &types,
		Address:       easygo.NewString(shopItem.address),
		DetailAddress: easygo.NewString(shopItem.detail_address),
		Name:          easygo.NewString(shopItem.name),
		Sex:           easygo.NewInt32(shopItem.sex),
		State:         easygo.NewInt32(shopItem.state),
		UserName:      easygo.NewString(shopItem.userName),
		Phone:         easygo.NewString(shopItem.phone),
		PointCardName: easygo.NewString(shopItem.pointCardName),
		CopyName:      easygo.NewString(shopItem.name),
		Types:         easygo.NewInt32(typess),
	}
	//如果是点卡的时候,为了客户端不改代码把商品名称设置为点卡名称
	if shopItem.item_type == for_game.SHOP_POINT_CARD_CATEGORY {
		itemDetail.Name = easygo.NewString(shopItem.pointCardName)
	}

	return &itemDetail
}
func (self *Shop) ItemFile(shopItem *ShopItem) []*share_message.ItemFile {

	if shopItem == nil {
		return []*share_message.ItemFile{}
	}

	var itemFiles = []*share_message.ItemFile{}

	var mItemFiles []ItemFile = shopItem.item_files
	if nil != mItemFiles && len(mItemFiles) > 0 {
		for i := 0; i < len(mItemFiles); i++ {
			var itemFile *share_message.ItemFile = &share_message.ItemFile{
				FileUrl:    &mItemFiles[i].file_url,
				FileType:   &mItemFiles[i].file_type,
				FileWidth:  &mItemFiles[i].file_width,
				FileHeight: &mItemFiles[i].file_height,
			}

			itemFiles = append(itemFiles, itemFile)
		}
	}

	return itemFiles
}

func (self *Shop) BriefItem(item *ShopItem) *share_message.ShopItem {
	if item == nil {
		return nil
	}
	var itemFile share_message.ItemFile = share_message.ItemFile{
		FileUrl:    &item.item_files[0].file_url,
		FileType:   &item.item_files[0].file_type,
		FileWidth:  &item.item_files[0].file_width,
		FileHeight: &item.item_files[0].file_height,
	}

	return &share_message.ShopItem{
		ItemId:     &item.item_id,
		Price:      &item.price,
		Title:      &item.title,
		ItemFile:   &itemFile,
		StoreCount: &item.realStoreCnt,
		PlayerId:   &item.player_id,
		Nickname:   &item.nickname,
		Avatar:     &item.avatar,
		Account:    &item.account,
		Sex:        &item.sex,
		Name:       &item.name,
		State:      &item.state,
	}

}

func (self *Shop) GetTimeLongforDetail(createTime int64) string {

	var timeLong string = ""

	var nowTime int64 = time.Now().Unix()
	//现在时间距离数据库时间的距离
	var timeLeng int64 = nowTime - createTime

	var tempTime int64 = 0

	if 0 <= timeLeng && timeLeng < 60 {
		timeLong = "刚刚"
	} else if 60 <= timeLeng && timeLeng < 3600 {
		tempTime = timeLeng / 60
		timeLong = easygo.IntToString(int(tempTime)) + "分钟前"
	} else if 3600 <= timeLeng && timeLeng < 86400 {
		tempTime = timeLeng / (60 * 60)
		timeLong = easygo.IntToString(int(tempTime)) + "小时前"
	} else if 86400 <= timeLeng && timeLeng < 2592000 {
		tempTime = timeLeng / (60 * 60 * 24)
		timeLong = easygo.IntToString(int(tempTime)) + "天前"
	} else if 2592000 <= timeLeng && timeLeng < 31104000 {
		tempTime = timeLeng / (60 * 60 * 24 * 30)
		timeLong = easygo.IntToString(int(tempTime)) + "个月前"
	} else if 31104000 <= timeLeng {
		tempTime = timeLeng / (60 * 60 * 24 * 30 * 12)
		timeLong = easygo.IntToString(int(tempTime)) + "年前"
	} else {
		timeLong = "刚刚"
	}

	return timeLong
}

func (self *Shop) GetDayLongforSeller(createTime int64) int32 {

	var dayLong int32 = 1

	var nowTime int64 = time.Now().Unix()
	//现在时间距离数据库时间的距离
	var timeLeng int64 = nowTime - createTime

	var tempDay int64 = 0

	if 0 <= timeLeng && timeLeng < 60 {
		dayLong = 1
	} else if 60 <= timeLeng && timeLeng < 3600 {
		dayLong = 1
	} else if 3600 <= timeLeng && timeLeng < 86400 {
		dayLong = 1
	} else if 86400 <= timeLeng {
		tempDay = timeLeng / (60 * 60 * 24)
		dayLong = int32(tempDay)
	} else {
		dayLong = 1
	}

	return dayLong
}

func (self *Shop) GetTimeLongforComent(createTime int64) string {

	var timeLong string = ""

	var nowTime int64 = time.Now().Unix()
	//现在时间距离数据库时间的距离
	var timeLeng int64 = nowTime - createTime

	var tempTime int64 = 0

	if 0 <= timeLeng && timeLeng < 60 {
		timeLong = "1分钟前"
	} else if 60 <= timeLeng && timeLeng < 3600 {
		tempTime = timeLeng / 60
		timeLong = easygo.IntToString(int(tempTime)) + "分钟前"

	} else if 3600 <= timeLeng && timeLeng < 86400 {
		tempTime = timeLeng / (60 * 60)
		timeLong = easygo.IntToString(int(tempTime)) + "小时前"
	} else if 86400 <= timeLeng && timeLeng < 2592000 {
		tempTime = timeLeng / (60 * 60 * 24)
		timeLong = easygo.IntToString(int(tempTime)) + "天前"
	} else if 2592000 <= timeLeng && timeLeng < 31104000 {
		tempTime = timeLeng / (60 * 60 * 24 * 30)
		timeLong = easygo.IntToString(int(tempTime)) + "个月前"
	} else if 31104000 <= timeLeng {
		tempTime = timeLeng / (60 * 60 * 24 * 30 * 12)
		timeLong = easygo.IntToString(int(tempTime)) + "年前"
	} else {
		timeLong = "1分钟前"
	}

	return timeLong
}

func (self *Shop) SellerInfo(shopItemPar *ShopItem) *share_message.SellerInfo {

	result := for_game.GetRedisPlayerBase(shopItemPar.player_id)
	account := for_game.GetRedisAccountObj(shopItemPar.player_id)

	if nil == result || account == nil {
		return nil
	}
	//注册天数
	var registerDay int32 = self.GetDayLongforSeller(account.GetCreateTime() / 1e3)

	//卖了多少件宝贝
	var sellItemCount int32
	var fakePlayFinOrderCnt int32
	var shopPlayer = share_message.TableShopPlayer{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_PLAYER)
	e := col.Find(bson.M{"_id": shopItemPar.player_id}).Limit(1).One(&shopPlayer)
	closeFun()

	if nil == e {
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
			{"$match": bson.M{"$and": []bson.M{bson.M{"sponsor_id": shopItemPar.player_id},
				bson.M{"$or": []bson.M{bson.M{"state": for_game.SHOP_ORDER_FINISH},
					bson.M{"state": for_game.SHOP_ORDER_EVALUTE}}}}}},
			{"$group": bson.M{"_id": "$sponsor_id", "total": bson.M{"$sum": "$items.count"}}}})

		rst := make([]bson.M, 0)
		e := query.All(&rst)
		var sum int = 0
		if nil == e {
			if rst != nil && len(rst) > 0 {
				sum = (rst[0]["total"]).(int)
			}
		} else {
			sum = 0
		}
		sellItemCount = int32(sum)
	}

	//付款数
	var payCnt int32
	if shopItemPar.fakePayCnt > 0 {
		payCnt = shopItemPar.fakePayCnt
	} else {
		//不实时从订单表中取得数据，从商品表冗余取得数据
		payCnt = shopItemPar.realPayCnt
	}

	//浏览数
	var pageViews int32
	if shopItemPar.fakePageViews > 0 {
		pageViews = shopItemPar.fakePageViews
	} else {
		pageViews = shopItemPar.realPageViews
	}

	var peopleId string = result.GetPeopleId()
	var realName string = result.GetRealName()
	var nameAuthBool bool = self.CheckPeopleAuth(peopleId, realName)
	var nameAuth int32 = 0
	if nameAuthBool {
		nameAuth = 1
	} else {
		nameAuth = 0
	}

	photo := ""
	photoList := result.GetPhoto()
	if nil != photoList && len(photoList) != 0 {
		photo = photoList[0]
	}

	var sellerInfo share_message.SellerInfo = share_message.SellerInfo{
		Nickname:      easygo.NewString(shopItemPar.nickname),
		Avatar:        easygo.NewString(shopItemPar.avatar),
		RegisterDay:   easygo.NewInt32(registerDay),
		SellItemCount: easygo.NewInt32(sellItemCount),
		NameAuth:      easygo.NewInt32(nameAuth),
		PlayerId:      easygo.NewInt64(result.GetPlayerId()),
		Account:       easygo.NewString(result.GetAccount()),
		Phone:         easygo.NewString(result.GetPhone()),
		Photo:         easygo.NewString(photo),
		Signature:     easygo.NewString(result.GetSignature()),
		Sex:           easygo.NewInt32(shopItemPar.sex),
		PaymentCount:  easygo.NewInt32(payCnt),
		PageViews:     easygo.NewInt32(pageViews),
		Types:         easygo.NewInt32(result.GetTypes()),
	}
	return &sellerInfo
}

func (self *Shop) CommInfoForDetail(shopItemPar *ShopItem, flag *share_message.BuySell_Type) *share_message.CommInfoForDetail {

	commInfoForDetail := &share_message.CommInfoForDetail{}

	//留言总数
	var commCount int32
	//最热门的留言信息
	commentInfo := share_message.CommentInfo{}
	//好评率
	var goodCommentRate int32

	newestComm := &share_message.TableItemComment{}
	var realMaxComm *share_message.TableItemComment
	var fakeMaxComm *share_message.TableItemComment
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ITEM_COMMENT)
	defer closeFun()
	realMaxQuery := col.Find(bson.M{"item_id": shopItemPar.item_id, "status": bson.M{"$ne": for_game.SHOP_COMMENT_DELETE}}).Sort("-real_like_count")
	fakeMaxQuery := col.Find(bson.M{"item_id": shopItemPar.item_id, "status": bson.M{"$ne": for_game.SHOP_COMMENT_DELETE}}).Sort("-fake_like_count")

	cnt, errVa := realMaxQuery.Count()
	if errVa != nil {
		logs.Error(errVa)
		return commInfoForDetail
	}
	commCount = int32(cnt)

	realMaxErrQuery := realMaxQuery.Limit(1).One(&realMaxComm)
	if realMaxErrQuery != nil && realMaxErrQuery != mgo.ErrNotFound {
		logs.Error(realMaxErrQuery)
		return commInfoForDetail
	}

	fakeMaxErrQuery := fakeMaxQuery.Limit(1).One(&fakeMaxComm)
	if fakeMaxErrQuery != nil && fakeMaxErrQuery != mgo.ErrNotFound {
		logs.Error(fakeMaxErrQuery)
		return commInfoForDetail
	}

	if realMaxComm == nil && fakeMaxComm == nil {
		return commInfoForDetail
	} else if realMaxComm != nil && fakeMaxComm == nil {
		newestComm = realMaxComm
	} else if realMaxComm == nil && fakeMaxComm != nil {
		newestComm = fakeMaxComm
	} else if realMaxComm != nil && fakeMaxComm != nil {
		if realMaxComm.GetRealLikeCount() >= fakeMaxComm.GetFakeLikeCount() {
			newestComm = realMaxComm
		} else {
			newestComm = fakeMaxComm
		}
	}

	//设置好评率
	//固定好评率设置不为0的时候 显示固定好评率
	if shopItemPar.fakeFixGoodCommRate > 0 {
		goodCommentRate = shopItemPar.fakeFixGoodCommRate

		//假的好评数和假的评价总数的判断
	} else if shopItemPar.fakeGoodCommCnt > 0 && shopItemPar.fakeFinCommCnt > 0 {
		goodCommentRate = shopItemPar.fakeGoodCommCnt * 100 / shopItemPar.fakeFinCommCnt
	} else if shopItemPar.realGoodCommCnt > 0 && shopItemPar.realFinCommCnt > 0 {
		goodCommentRate = shopItemPar.realGoodCommCnt * 100 / shopItemPar.realFinCommCnt
	} else {
		goodCommentRate = 0
	}

	//设置最新留言信息
	if newestComm != nil && newestComm.PlayerId != nil {
		//点赞数
		var likeCnt int32

		if newestComm.GetFakeLikeCount() > 0 {
			likeCnt = newestComm.GetFakeLikeCount()
		} else {
			likeCnt = newestComm.GetRealLikeCount()
		}

		//通过买家卖家视角来显示昵称
		var nickName string
		if *flag == share_message.BuySell_Type_Buyer {
			nickName = GetMarkNickName(newestComm.GetNickname())
		} else {
			nickName = newestComm.GetNickname()
		}

		commentInfo = share_message.CommentInfo{
			CommentId: newestComm.CommentId,
			PlayerId:  newestComm.PlayerId,
			Avatar:    newestComm.Avatar,
			Nickname:  easygo.NewString(nickName),
			Content:   newestComm.Content,
			ItemId:    newestComm.ItemId,
			TimeLong:  easygo.NewString(GetYMDTime(newestComm.GetCreateTime())),
			Sex:       newestComm.Sex,
			StarLevel: newestComm.StarLevel,
			LikeCount: easygo.NewInt32(likeCnt),
		}
	}
	commInfoForDetail = &share_message.CommInfoForDetail{
		CommentInfo:     &commentInfo,
		CommentCount:    easygo.NewInt32(commCount),
		GoodCommentRate: easygo.NewInt32(goodCommentRate),
	}
	return commInfoForDetail

}

func (self *Shop) BuyCallBack(order_id int64) {

	var bill *share_message.TableBill = &share_message.TableBill{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_BILLS)
	defer closeFun()

	e := col.Find(bson.M{"_id": order_id}).One(bill)
	if e != nil && e != mgo.ErrNotFound {
		logs.Error(e)
		//恢复库存
		rst1 := for_game.ShopRecoverStock(order_id)
		if rst1 != "" {
			logs.Error(rst1)
			return
		}
		return
	}

	if e == mgo.ErrNotFound {
		delStr := self.DelOrder(order_id)
		if delStr != "" {
			//恢复库存
			rst1 := for_game.ShopRecoverStock(order_id)
			if rst1 != "" {
				logs.Error(rst1)
				return
			}
			return
		}
	} else {

		for _, id := range bill.OrderList {
			delStr := self.DelOrder(id)
			if delStr != "" {
				//恢复库存
				rst1 := for_game.ShopRecoverStock(id)
				if rst1 != "" {
					logs.Error(rst1)
					return
				}
				return
			}
		}
		//更新bill总订单的状态(第一次立即执行,出错走定时器重试执行)
		SaveDataToDBForBuyCallBack(order_id, 0)
	}
}

//付款定时存储
func SaveDataToDBForBuyCallBack(orderId int64, t time.Duration) {
	t += 2 * time.Second    //现在是2秒间隔一次，间隔多久自定义
	if t > 10*time.Second { //现在是执行10秒跳出，多久跳出自定义
		return
	}
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_BILLS)
	defer closeFun()

	e := col.Update(
		bson.M{"_id": orderId},
		bson.M{"$set": bson.M{"state": for_game.SHOP_ORDER_WAIT_SEND}})
	if e != nil {
		logs.Error(e)
		fun := func() {
			SaveDataToDBForBuyCallBack(orderId, t)
		}
		easygo.AfterFunc(t, fun)
	}
}

//付款定时存储相关订单
func SaveDataToDBForDelOrderOrders(orderId int64, itemType int32, h5BuyFlag string, t time.Duration) {
	t += 2 * time.Second    //现在是2秒间隔一次，间隔多久自定义
	if t > 10*time.Second { //现在是执行10秒跳出，多久跳出自定义
		return
	}
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()

	//是h5过来的购买且是点卡类型的订单
	if itemType == for_game.SHOP_POINT_CARD_CATEGORY && h5BuyFlag != "" {
		e := col.Update(
			bson.M{"_id": orderId, "state": for_game.SHOP_ORDER_WAIT_PAY},
			bson.M{"$set": bson.M{"state": for_game.SHOP_ORDER_EVALUTE,
				"pay_time":             time.Now().Unix(),
				"send_time":            time.Now().Unix(),
				"receive_time":         time.Now().Unix(),
				"finish_time":          time.Now().Unix(),
				"update_time":          time.Now().Unix(),
				"receiver_notify_flag": true,
				"sponsor_notify_flag":  true}})

		if e != nil {
			logs.Error(e)
			fun := func() {
				SaveDataToDBForDelOrderOrders(orderId, itemType, h5BuyFlag, t)
			}
			easygo.AfterFunc(t, fun)
		}

		//app端购买的点卡或者普通商品
	} else {
		e := col.Update(
			bson.M{"_id": orderId, "state": for_game.SHOP_ORDER_WAIT_PAY},
			bson.M{"$set": bson.M{"state": for_game.SHOP_ORDER_WAIT_SEND,
				"pay_time":             time.Now().Unix(),
				"update_time":          time.Now().Unix(),
				"receiver_notify_flag": true,
				"sponsor_notify_flag":  true}})

		if e != nil {
			logs.Error(e)
			fun := func() {
				SaveDataToDBForDelOrderOrders(orderId, itemType, h5BuyFlag, t)
			}
			easygo.AfterFunc(t, fun)
		}
	}
}

//如果是点卡的时候去更新导入库中的库存并且将点卡信息挂到订单信息中
func SaveDataToShopPointCard(order_id int64, t time.Duration) {
	t += 2 * time.Second    //现在是2秒间隔一次，间隔多久自定义
	if t > 10*time.Second { //现在是执行10秒跳出，多久跳出自定义
		return
	}

	order := share_message.TableShopOrder{}

	colOrder, closeFunOrder := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFunOrder()

	colCard, closeFunCard := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_POINT_CARD)
	defer closeFunCard()

	var itemPointCards []*share_message.ShopPointCardInfo = make([]*share_message.ShopPointCardInfo, 0)
	var tempCardIds []int64 = []int64{}

	eOrderQuery := colOrder.Find(bson.M{"_id": order_id}).One(&order)

	if eOrderQuery != nil {
		logs.Error(eOrderQuery)
		fun := func() {
			SaveDataToShopPointCard(order_id, t)
		}
		easygo.AfterFunc(t, fun)
		return
	} else {
		if order.GetItems().GetPointCardInfos() == nil || len(order.GetItems().GetPointCardInfos()) <= 0 {
			//通过物品的个数取得卡密个数
			pointCardList := for_game.GetPointCardByBuyInfos(order.GetSponsorAccount(),
				order.GetItems().GetPointCardName(),
				order.GetItems().GetCount())
			if nil != pointCardList && len(pointCardList) > 0 {
				for _, valuePoint := range pointCardList {
					itemPointCard := share_message.ShopPointCardInfo{
						CardId:       easygo.NewInt64(valuePoint.GetCardId()),
						CardNo:       easygo.NewString(valuePoint.GetCardNo()),
						CardPassword: easygo.NewString(valuePoint.GetCardPassword()),
						Key:          easygo.NewString(valuePoint.GetKey()),
					}

					itemPointCards = append(itemPointCards, &itemPointCard)
					tempCardIds = append(tempCardIds, valuePoint.GetCardId())
				}
			}
		} else {
			for _, valuePoint := range order.GetItems().GetPointCardInfos() {
				itemPointCard := share_message.ShopPointCardInfo{
					CardId:       easygo.NewInt64(valuePoint.GetCardId()),
					CardNo:       easygo.NewString(valuePoint.GetCardNo()),
					CardPassword: easygo.NewString(valuePoint.GetCardPassword()),
					Key:          easygo.NewString(valuePoint.GetKey()),
				}

				itemPointCards = append(itemPointCards, &itemPointCard)
				tempCardIds = append(tempCardIds, valuePoint.GetCardId())
			}
		}

		//更新导入库中点卡状态和订单id
		_, errCard := colCard.UpdateAll(bson.M{"_id": bson.M{"$in": tempCardIds}}, bson.M{"$set": bson.M{"card_status": for_game.SHOP_POINT_CARD_SELLOUT, "order_no": order_id}})

		//绑定点卡到订单
		eOrder := colOrder.Update(
			bson.M{"_id": order_id},
			bson.M{"$set": bson.M{"items.pointCardInfos": itemPointCards}})

		if eOrder != nil || errCard != nil {
			logs.Error(eOrder)
			logs.Error(errCard)
			fun := func() {
				SaveDataToShopPointCard(order_id, t)
			}
			easygo.AfterFunc(t, fun)
		}
	}
}

////付款定时存储相关商品
//func SaveDataToDBForDelOrderItems(itemId int64, t time.Duration) {
//	t += 2 * time.Second    //现在是2秒间隔一次，间隔多久自定义
//	if t > 10*time.Second { //现在是执行10秒跳出，多久跳出自定义
//		return
//	}
//	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
//	defer closeFun()
//
//	e := col.Update(bson.M{"_id": itemId}, bson.M{"$set": bson.M{"state": for_game.SHOP_ITEM_SOLD_OUT, "sold_out_time": time.Now().Unix()}})
//
//	if e != nil {
//		logs.Error(e)
//		fun := func() {
//			SaveDataToDBForDelOrderItems(itemId, t)
//		}
//		easygo.AfterFunc(t, fun)
//	}
//}

func (self *Shop) DelOrder(order_id int64) string {

	order := share_message.TableShopOrder{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()
	e := col.Find(bson.M{"_id": order_id, "state": for_game.SHOP_ORDER_WAIT_PAY}).One(&order)

	if e != nil && e != mgo.ErrNotFound {
		logs.Error(e)
		return "操作失败"
	}

	if e == mgo.ErrNotFound {
		logs.Error(e)
		return "订单不存在"
	}

	s := fmt.Sprintf("订单%v处理开始", order_id)
	for_game.WriteFile("shop_order.log", s)

	//如果是点卡的时候去更新导入库中的库存并且将点卡信息挂到订单信息中
	if order.GetItems() != nil && order.GetItems().GetItemType() == for_game.SHOP_POINT_CARD_CATEGORY {

		s = fmt.Sprintf("点卡订单%v 绑定点卡信息到订单开始", order_id)
		for_game.WriteFile("shop_order.log", s)

		//支付后绑定点卡信息到订单(第一次立即启动执行,出错走定时器执行)
		SaveDataToShopPointCard(order_id, 0)

		s = fmt.Sprintf("点卡订单%v 绑定点卡信息到订单结束", order_id)
		for_game.WriteFile("shop_order.log", s)

	}

	s = fmt.Sprintf("订单%v 状态修改开始", order_id)
	for_game.WriteFile("shop_order.log", s)

	//更新订单状态(第一次立即启动执行,出错走定时器执行)
	SaveDataToDBForDelOrderOrders(order_id, order.GetItems().GetItemType(), order.GetH5SearchCon(), 0)

	s = fmt.Sprintf("订单%v 状态修改结束", order_id)
	for_game.WriteFile("shop_order.log", s)

	s = fmt.Sprintf("订单%v处理结束", order_id)
	for_game.WriteFile("shop_order.log", s)

	//app端购买的时候,不管是点卡还是普通商品都push通知
	//在app端购买的时候(包括开放点卡商品的时候)
	if order.GetH5SearchCon() == "" {
		//买家成功付款后通知商家
		easygo.Spawn(func(orderPara share_message.TableShopOrder) {

			if &orderPara != nil {

				var content string = MESSAGE_TO_SELLER_PAY
				typeValue := share_message.BuySell_Type_Seller

				ShopInstance.InsMessageNotify(
					easygo.NewString(content),
					&typeValue,
					&orderPara)

				var jgContent string = MESSAGE_TO_SELLER_PAY_PUSH

				ShopInstance.JGMessageNotify(jgContent, orderPara.GetSponsorId(), orderPara.GetOrderId(), typeValue)

				// 修改商品表中真实和虚假的付款数
				AddPayCnt(orderPara.GetItems().GetItemId())

				//买家付款 商城订单红点推送
				SendMsgToHallClientNew([]int64{orderPara.GetReceiverId()}, "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
					OrderId: easygo.NewInt64(orderPara.GetOrderId())})
				/*				SendToPlayer(orderPara.GetReceiverId(), "RpcShopOrderNotify",
								&share_message.ShopOrderNotifyInfoWithWho{
									PlayerId: easygo.NewInt64(orderPara.GetReceiverId()),
									OrderId:  easygo.NewInt64(orderPara.GetOrderId()),
								})*/

				//买家付款 商城订单红点推送
				SendMsgToHallClientNew([]int64{orderPara.GetSponsorId()}, "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
					OrderId: easygo.NewInt64(orderPara.GetOrderId())})
				/*		SendToPlayer(orderPara.GetSponsorId(), "RpcShopOrderNotify",
						&share_message.ShopOrderNotifyInfoWithWho{
							PlayerId: easygo.NewInt64(orderPara.GetSponsorId()),
							OrderId:  easygo.NewInt64(orderPara.GetOrderId()),
						})*/
			} else {
				logs.Debug("买家成功付款后通知商家,缺少订单")
			}
		}, order)
	} else {

		//h5购买的商品订单已经评价状态要给卖家加钱
		SendMsgToServerNewEx(order.GetSponsorId(),
			"RpcShopPaySeller",
			&share_message.PaySellerInfo{
				OrderId:    order.OrderId,
				Money:      easygo.NewInt32(order.Items.GetPrice() * order.Items.GetCount()),
				Sponsor_Id: order.SponsorId,
				ReceiverId: order.ReceiverId,
				PayType:    easygo.NewInt32(0),
			})

		//h5购买成功状态已经是已经评价所以在这里增加付款数和订单成交数
		//买家成功付款后通知商家
		easygo.Spawn(func(itemIdPara int64) {

			// 修改商品表中真实和虚假的付款数
			AddPayCnt(itemIdPara)

			//修改该商品完成的订单数
			AddFinOrderCnt(itemIdPara)

			//是邮件并且是点卡才邮件通知
			if strings.Contains(order.GetH5SearchCon(), "@") &&
				order.GetItems() != nil &&
				order.GetItems().GetItemType() == for_game.SHOP_POINT_CARD_CATEGORY {
				easygo.Spawn(func(orderIdPara int64) {
					errStr := DoSendMail(order.GetOrderId())
					if errStr != "" {
						logs.Error(errStr)
					}
				}, order.GetOrderId())

			}

		}, order.GetItems().GetItemId())

		//处理商城app的push推送
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
			} else {
				logs.Debug("h5购买后发通知买家,缺少订单")
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

				//自动收货 商城订单红点推送
				SendMsgToHallClientNew([]int64{orderPara.GetReceiverId()}, "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
					OrderId: easygo.NewInt64(orderPara.GetOrderId())})

				//自动收货 商城订单红点推送
				SendMsgToHallClientNew([]int64{orderPara.GetSponsorId()}, "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
					OrderId: easygo.NewInt64(orderPara.GetOrderId())})

			} else {
				logs.Debug("h5购买后发通知卖家,缺少订单")
			}
		}, order)

	}

	return ""
}

//同一时间,集群环境下每台服务器都要做，客户端看到的数据基于每台服务器自身加载内存的数据
//这只是首页列表显示要是有部分误差没有关系最终每个用户都一致
func (self *Shop) UpdateItemList() {

	func() {
		var newList ShopItemList = ShopItemList{}
		var newMap ShopItemMap = ShopItemMap{}

		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
		defer closeFun()

		var list []*share_message.TableShopItem

		//加载在架但不是点卡的商品到内存中
		e := col.Find(bson.M{"state": for_game.SHOP_ITEM_SALE,
			"stock_count": bson.M{"$gt": 0},
			"type.type":   bson.M{"$ne": 45}}).All(&list)

		if e == nil {
			for _, value := range list {
				itemFiles := []ItemFile{}
				if nil != value.ItemFiles && len(value.ItemFiles) > 0 {
					for i := 0; i < len(value.ItemFiles); i++ {
						var itemFile ItemFile = ItemFile{
							file_url:    value.ItemFiles[i].GetFileUrl(),
							file_type:   value.ItemFiles[i].GetFileType(),
							file_width:  value.ItemFiles[i].GetFileWidth(),
							file_height: value.ItemFiles[i].GetFileHeight()}
						itemFiles = append(itemFiles, itemFile)
					}
				}

				newItem := ShopItem{
					item_id:             value.GetItemId(),
					item_files:          itemFiles,
					title:               value.GetTitle(),
					origin_price:        value.GetOriginPrice(),
					price:               value.GetPrice(),
					userName:            value.GetUserName(),
					phone:               value.GetPhone(),
					address:             value.GetAddress(),
					detail_address:      value.GetDetailAddress(),
					avatar:              value.GetAvatar(),
					player_id:           value.GetPlayerId(),
					nickname:            value.GetNickname(),
					create_time:         value.GetCreateTime(),
					stock_count:         value.GetStockCount(),
					account:             value.GetPlayerAccount(),
					name:                value.GetName(),
					sex:                 value.GetSex(),
					state:               value.GetState(),
					realPayCnt:          value.GetRealPayCnt(),
					fakePayCnt:          value.GetFakePayCnt(),
					realPageViews:       value.GetRealPageViews(),
					fakePageViews:       value.GetFakePageViews(),
					realGoodCommCnt:     value.GetRealGoodCommCnt(),
					fakeGoodCommCnt:     value.GetFakeGoodCommCnt(),
					realFinCommCnt:      value.GetRealFinCommCnt(),
					fakeFinCommCnt:      value.GetFakeFinCommCnt(),
					fakeFixGoodCommRate: value.GetFakeFixGoodCommRate(),
					realCommentCnt:      value.GetRealCommentCnt(),
					realStoreCnt:        value.GetRealStoreCnt(),
					realFinOrderCnt:     value.GetRealFinOrderCnt(),
					pointCardName:       value.GetPointCardName(),
				}
				if nil != value.GetType() {
					newItem.item_type = value.GetType().GetType()
					newItem.other_type = value.GetType().GetOtherType()
				} else {
					newItem.item_type = 0 //全部
					newItem.other_type = []string{}
				}

				newList = append(newList, &newItem)
				newMap[newItem.item_id] = &newItem
			}

			self.shop_items = &newList
			self.shop_items_by_id = &newMap

		} else {
			logs.Error(e)
		}

	}()

	easygo.AfterFunc(time.Second*30, self.UpdateItemList)
}

//自动收货
//这个定时任务没法移到目前static的单服务上要跟大厅通信等
//同一时间，集群环境下每台服务器上只能一台做该业务，用分布式不重试锁，没有取到订单的锁的下次做
//保证服务器之间的处理是单线的
func (self *Shop) UpdateReceiveOrder() {

	//超过七天自动收货,10分钟做一次自动收货
	//加函数是为了defer处理和取得分布式锁失败直接退出函数
	func() {

		//如果这个定时不移到单服务上去必须加个定时自己取得自己的互斥分布式锁(失效时间设置1200秒即使20分钟)
		//这个锁是不重试的
		//集群环境下为了避免多次自动收货导致给商家多次加钱
		errLock1 := easygo.RedisMgr.GetC().DoRedisLockNoRetry(for_game.SHOP_AUTO_RECEIVE_MUTEX_SERVER, 1200)
		defer easygo.RedisMgr.GetC().DoRedisUnlock(for_game.SHOP_AUTO_RECEIVE_MUTEX_SERVER)

		//如果未取得锁
		if errLock1 != nil {
			s := fmt.Sprintf("UpdateReceiveOrder 单key取得redis分布式无重试锁失败,redis key is %v", for_game.SHOP_AUTO_RECEIVE_MUTEX_SERVER)
			logs.Error(s)
			logs.Error(errLock1)
			//直接退出函数进入下一次循环
			return
		}

		now := time.Now().Unix()
		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
		defer closeFun()

		list := []*share_message.TableShopOrder{}

		e := col.Find(bson.M{"state": for_game.SHOP_ORDER_WAIT_RECEIVE}).All(&list)

		if e != nil {

			logs.Error(e)

		} else {

			//取到锁就执行业务
			for _, value := range list {
				if now > value.GetReceiveTime() && value.GetState() == for_game.SHOP_ORDER_WAIT_RECEIVE {
					//加函数为了defer处理
					func() {

						//这个锁必须在函数中而且要在if的条件中取得
						//以每个订单为单位取得锁,取不到直接做下次做
						lockKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_RECEIVE_MUTEX, value.GetOrderId())
						//取得分布式锁,跟主动收货互斥（失效时间设置20秒）
						errLock := easygo.RedisMgr.GetC().DoRedisLockNoRetry(lockKey, 20)
						defer easygo.RedisMgr.GetC().DoRedisUnlock(lockKey)

						//如果重试后还未取得锁
						if errLock != nil {
							s := fmt.Sprintf("UpdateReceiveOrder定时任务 单key取得redis分布式无重试锁失败,redis key is %v", lockKey)
							logs.Error(s)
							logs.Error(errLock)
							//直接退出函数进入下一次循环
							return
						}

						var nowTime int64 = time.Now().Unix()

						e := col.Update(
							bson.M{"_id": value.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_RECEIVE},
							bson.M{"$set": bson.M{"state": for_game.SHOP_ORDER_FINISH,
								"receive_time":         nowTime,
								"finish_time":          nowTime,
								"receiver_notify_flag": true,
								"sponsor_notify_flag":  true,
							}})

						if e != nil {
							logs.Error(e, value.GetOrderId())
							if e == mgo.ErrNotFound {
								s := fmt.Sprintf("UpdateReceiveOrder定时任务 %v订单用户收货操作,在状态发生变化,这里不做了,不用管,属于正常的", value.GetOrderId())
								logs.Error(s)
							}
							//直接退出函数进入下一次循环
							return
						} else {

							easygo.Spawn(func() {
								for_game.MakePlayerBehaviorReport(4, 0, nil, value, nil, nil) //生成用户行为报表商城订单完成相关字段 已优化到Redis
								// for_game.MakeOperationChannelReport(4, value.GetReceiverId(), "", value, nil) //生成运营渠道数据汇总报表 已优化到Redis
							})

							SendMsgToServerNewEx(value.GetSponsorId(),
								"RpcShopPaySeller",
								&share_message.PaySellerInfo{
									OrderId:    value.OrderId,
									Money:      easygo.NewInt32(value.Items.GetPrice() * value.Items.GetCount()),
									Sponsor_Id: value.SponsorId,
									ReceiverId: value.ReceiverId,
									PayType:    easygo.NewInt32(0),
								})

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
								} else {
									logs.Debug("确认收货后发通知买家,缺少订单")
								}

							}, value)

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

									//自动收货 商城订单红点推送
									SendMsgToHallClientNew([]int64{orderPara.GetReceiverId()}, "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
										OrderId: easygo.NewInt64(orderPara.GetOrderId())})
									/*			SendToPlayer(orderPara.GetReceiverId(), "RpcShopOrderNotify",
												&share_message.ShopOrderNotifyInfoWithWho{
													PlayerId: easygo.NewInt64(orderPara.GetReceiverId()),
													OrderId:  easygo.NewInt64(orderPara.GetOrderId()),
												})*/

									//自动收货 商城订单红点推送
									SendMsgToHallClientNew([]int64{orderPara.GetSponsorId()}, "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
										OrderId: easygo.NewInt64(orderPara.GetOrderId())})
									/*	SendToPlayer(orderPara.GetSponsorId(), "RpcShopOrderNotify",
										&share_message.ShopOrderNotifyInfoWithWho{
											PlayerId: easygo.NewInt64(orderPara.GetSponsorId()),
											OrderId:  easygo.NewInt64(orderPara.GetOrderId()),
										})*/
								} else {
									logs.Debug("确认收货后发通知卖家,缺少订单")
								}
							}, value)
						}
					}()

				}
			}
		}

	}()

	easygo.AfterFunc(time.Second*60*10, self.UpdateReceiveOrder)
}

func (self *Shop) CheckPeopleAuth(peopleId string, realName string) bool {
	if peopleId != "" && realName != "" {
		return true
	}
	return false
}

func (self *Shop) GetBlackLists(playerId int64) []PLAYER_ID {
	playInfo := for_game.GetRedisPlayerBase(playerId)
	var blackList []PLAYER_ID = []PLAYER_ID{}

	if nil != playInfo {
		blackList = playInfo.GetBlackList()
	}

	colVar, closeFunVar := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	players := []*share_message.PlayerBase{}
	errVar := colVar.Find(bson.M{"BlackList": playerId}).All(&players)
	closeFunVar()

	if errVar == nil && len(players) > 0 {
		for _, black := range players {
			blackList = append(blackList, black.GetPlayerId())
		}
	}
	return blackList
}

func (self *Shop) GetShortAddress(addressSrc string) string {

	address := addressSrc
	addressArray := strings.Split(address, "-")

	if len(addressArray) >= 3 {
		address = addressArray[1]
	}
	return address
}

func (self *Shop) InsMessageNotify(content *string, typeValue *share_message.BuySell_Type, orderPara *share_message.TableShopOrder) {

	//订单通知
	var timeNow int64 = time.Now().Unix()
	messageId := easygo.NewInt64(for_game.NextId(for_game.TABLE_SHOP_MESSAGE))

	tableMessage := share_message.TableShopMessage{
		MessageId:        messageId,
		UserType:         easygo.NewInt32(int32(*typeValue)),
		SponsorPlayerId:  easygo.NewInt64(orderPara.GetSponsorId()),
		SponsorNickname:  easygo.NewString(orderPara.GetSponsorNickname()),
		SponsorAvatar:    easygo.NewString(orderPara.GetSponsorAvatar()),
		SponsorSex:       easygo.NewInt32(orderPara.GetSponsorSex()),
		ReceiverPlayerId: easygo.NewInt64(orderPara.GetReceiverId()),
		ReceiverNickname: easygo.NewString(orderPara.GetReceiverNickname()),
		ReceiverAvatar:   easygo.NewString(orderPara.GetReceiverAvatar()),
		ReceiverSex:      easygo.NewInt32(orderPara.GetReceiverSex()),
		File:             orderPara.Items.ItemFile,
		ItemName:         easygo.NewString(orderPara.Items.GetName()),
		ItemTitle:        easygo.NewString(orderPara.Items.GetTitle()),
		Content:          content,
		CreateTime:       easygo.NewInt64(timeNow),
		OrderId:          easygo.NewInt64(orderPara.GetOrderId()),
		ViewFlag:         easygo.NewBool(false),
		CopyName:         easygo.NewString(orderPara.Items.GetName()),
	}
	//如果是点卡就重新设置显示的商品名称
	if orderPara.GetItems() != nil && orderPara.GetItems().GetItemType() == for_game.SHOP_POINT_CARD_CATEGORY {
		tableMessage.ItemName = easygo.NewString(orderPara.Items.GetPointCardName())
	}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_MESSAGE)
	e := col.Insert(tableMessage)
	closeFun()

	if e == nil {

		var toPlayerId int64
		var nicknamePer string
		var avatarPer string

		if *typeValue == share_message.BuySell_Type_Buyer {
			toPlayerId = orderPara.GetReceiverId()
			nicknamePer = orderPara.GetSponsorNickname()
			avatarPer = orderPara.GetSponsorAvatar()

		} else {

			toPlayerId = orderPara.GetSponsorId()
			nicknamePer = orderPara.GetReceiverNickname()
			avatarPer = orderPara.GetReceiverAvatar()

		}
		shopMessage := &share_message.ShopItemMessage{

			MessageId:  messageId,
			Type:       typeValue,
			File:       orderPara.Items.ItemFile,
			Nickname:   easygo.NewString(nicknamePer),
			Avatar:     easygo.NewString(avatarPer),
			ItemName:   easygo.NewString(orderPara.Items.GetName()),
			ItemTitle:  easygo.NewString(orderPara.Items.GetTitle()),
			Content:    content,
			CreateTime: easygo.NewInt64(timeNow),
			OrderId:    easygo.NewInt64(orderPara.GetOrderId()),
			ShowTime:   easygo.NewString(util.FormatUnixTime(timeNow)),
			ViewFlag:   easygo.NewBool(false),
			CopyName:   easygo.NewString(orderPara.Items.GetName()),
		}
		//如果是点卡就重新设置显示的商品名称
		if orderPara.GetItems() != nil && orderPara.GetItems().GetItemType() == for_game.SHOP_POINT_CARD_CATEGORY {
			shopMessage.ItemName = easygo.NewString(orderPara.Items.GetPointCardName())
		}

		req := &share_message.ShopItemMessageInfo{
			Type:        typeValue,
			ShopMessage: shopMessage}

		SendMsgToHallClientNew([]int64{toPlayerId}, "RpcShopItemMessageNotify", req)
	} else {
		logs.Error(e)
	}
}

func (self *Shop) JGMessageNotify(content string,
	playerId int64,
	orderId int64,
	operaType share_message.BuySell_Type) {

	ids := for_game.GetJGIds([]int64{playerId})
	m := for_game.PushMessage{
		Title:       content,
		Content:     content,
		ContentType: for_game.JG_TYPE_SHOP,
		OrderId:     easygo.AnytoA(orderId),
		OperaType:   int32(operaType),
	}
	for_game.JGSendMessage(ids, m)
}

func (self *Shop) InsAliAuditErr(
	itemId *int64,
	origin *string,
	auditType *string,
	errCode *string,
	errContent *string,
	nowTime *int64) {

	aliFailInfo := share_message.TableShopAliAuditFail{
		ItemId:     itemId,
		Origin:     origin,
		Type:       auditType,
		ErrorCode:  errCode,
		Content:    errContent,
		CreateTime: nowTime,
	}

	colAudit, closeFunAudit := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ALI_AUDIT_FAIL)
	eAudit := colAudit.Insert(aliFailInfo)
	closeFunAudit()

	if eAudit != nil {
		logs.Error(eAudit)
	}
}

func (self *Shop) UpdatePageViews(pageViewFlag int32, shopItemPar *ShopItem) {
	//记录当天浏览量
	if pageViewFlag == 0 {
		//记录浏览量
		colPage, closeFunPage := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
		defer closeFunPage()

		//记录真实的浏览量
		//如果后台修改了假的浏览量则同时增加假的浏览量
		if shopItemPar.fakePageViews > 0 {
			errPage := colPage.Update(
				bson.M{"_id": shopItemPar.item_id},
				bson.M{"$inc": bson.M{"real_pageViews": 1, "fake_pageViews": 1}})
			if errPage != nil && errPage != mgo.ErrNotFound {
				logs.Error(errPage)
			}
			//只记录真实浏览量
		} else {
			errRealPage := colPage.Update(
				bson.M{"_id": shopItemPar.item_id},
				bson.M{"$inc": bson.M{"real_pageViews": 1}})
			if errRealPage != nil && errRealPage != mgo.ErrNotFound {
				logs.Error(errRealPage)
			}
		}
	}
}

func (self *Shop) GetShopItem(peopleFlag share_message.BuySell_Type, itemId int64) (string, *ShopItem) {

	errStr := ""
	var shopItem *ShopItem

	if peopleFlag == share_message.BuySell_Type_Buyer {
		shopItem = ShopInstance.GetItemFromCache(itemId)
		if shopItem == nil {
			errStr = DETAIL_SHOP_ITEM_NOT_SALE
		}
		return errStr, shopItem

	} else {

		shopItem = ShopInstance.GetItemFromCache(itemId)

		if shopItem == nil {
			col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
			defer closeFun()
			var item share_message.TableShopItem = share_message.TableShopItem{}

			e := col.Find(bson.M{"_id": itemId}).One(&item)
			if e == mgo.ErrNotFound {
				errStr = DETAIL_SHOP_ITEM_NOT_EXIST
				return errStr, shopItem
			}

			if e != nil {
				logs.Error(e)
				errStr = DATABASE_ERROR
				return errStr, shopItem
			}

			itemFiles := []ItemFile{}
			if nil != item.ItemFiles && len(item.ItemFiles) > 0 {
				for i := 0; i < len(item.ItemFiles); i++ {
					var itemFile ItemFile = ItemFile{
						file_url:  item.ItemFiles[i].GetFileUrl(),
						file_type: item.ItemFiles[i].GetFileType()}
					itemFiles = append(itemFiles, itemFile)
				}
			}

			newItem := ShopItem{
				item_id:             item.GetItemId(),
				item_files:          itemFiles,
				title:               item.GetTitle(),
				origin_price:        item.GetOriginPrice(),
				price:               item.GetPrice(),
				userName:            item.GetUserName(),
				phone:               item.GetPhone(),
				address:             item.GetAddress(),
				detail_address:      item.GetDetailAddress(),
				avatar:              item.GetAvatar(),
				player_id:           item.GetPlayerId(),
				item_type:           item.GetType().GetType(),
				other_type:          item.GetType().GetOtherType(),
				nickname:            item.GetNickname(),
				create_time:         item.GetCreateTime(),
				stock_count:         item.GetStockCount(),
				account:             item.GetPlayerAccount(),
				sex:                 item.GetSex(),
				name:                item.GetName(),
				state:               item.GetState(),
				realPayCnt:          item.GetRealPayCnt(),
				fakePayCnt:          item.GetFakePayCnt(),
				realPageViews:       item.GetRealPageViews(),
				fakePageViews:       item.GetFakePageViews(),
				realGoodCommCnt:     item.GetRealGoodCommCnt(),
				fakeGoodCommCnt:     item.GetFakeGoodCommCnt(),
				realFinCommCnt:      item.GetRealFinCommCnt(),
				fakeFinCommCnt:      item.GetFakeFinCommCnt(),
				fakeFixGoodCommRate: item.GetFakeFixGoodCommRate(),
				realCommentCnt:      item.GetRealCommentCnt(),
				realStoreCnt:        item.GetRealStoreCnt(),
				realFinOrderCnt:     item.GetRealFinOrderCnt(),
				pointCardName:       item.GetPointCardName(),
			}

			shopItem = &newItem
		}
	}

	return errStr, shopItem
}

func GetMarkNickName(nickName string) string {
	reNickName := nickName
	//通过买家卖家视角来显示昵称
	if len(reNickName) >= 1 {
		reNickName = string([]rune(reNickName)[:1])
		reNickName = fmt.Sprintf("%v%v", reNickName, "***")
	}

	return reNickName
}

func (list ShopItemList) Len() int { return len(list) }

func (list ShopItemList) Swap(i, j int) {
	s := list[j]
	list[j] = list[i]
	list[i] = s
}

func (list ShopItemList) Less(i, j int) bool {
	if list[i].realStoreCnt == list[j].realStoreCnt {
		return list[i].create_time > list[j].create_time
	}

	return list[i].realStoreCnt > list[j].realStoreCnt
}

func AddPayCnt(itemId int64) {

	var item share_message.TableShopItem = share_message.TableShopItem{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()

	e := col.Find(bson.M{"_id": itemId}).One(&item)

	if e != nil {
		logs.Error(e)
	}

	if e == nil {
		if item.GetFakePayCnt() <= 0 {

			err1 := col.Update(
				bson.M{"_id": itemId},
				bson.M{"$inc": bson.M{"real_payCnt": 1}})

			if err1 != nil {
				logs.Error(err1)
			}
		} else {
			err1 := col.Update(
				bson.M{"_id": itemId},
				bson.M{"$inc": bson.M{"real_payCnt": 1, "fake_payCnt": 1}})

			if err1 != nil {
				logs.Error(err1)
			}
		}
	}
}

func SubPayCnt(itemId int64) {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()

	var item share_message.TableShopItem = share_message.TableShopItem{}

	e := col.Find(bson.M{"_id": itemId}).One(&item)

	if e != nil {
		logs.Error(e)
	}

	if e == nil {
		if item.GetFakePayCnt() > 0 {
			if item.GetRealPayCnt() <= 0 {

				err1 := col.UpdateId(itemId,
					bson.M{"$inc": bson.M{"fake_payCnt": -1}})

				if err1 != nil && err1 != mgo.ErrNotFound {
					logs.Error(err1)
				}

			} else {
				err1 := col.UpdateId(itemId,
					bson.M{"$inc": bson.M{"real_payCnt": -1, "fake_payCnt": -1}})

				if err1 != nil && err1 != mgo.ErrNotFound {
					logs.Error(err1)
				}
			}

		} else {

			if item.GetRealPayCnt() > 0 {
				err1 := col.UpdateId(itemId,
					bson.M{"$inc": bson.M{"real_payCnt": -1}})

				if err1 != nil && err1 != mgo.ErrNotFound {
					logs.Error(err1)
				}
			}
		}
	}
}

func AddFinOrderCnt(itemId int64) {

	var item share_message.TableShopItem = share_message.TableShopItem{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()

	e := col.Find(bson.M{"_id": itemId}).One(&item)

	if e != nil {
		logs.Error(e)
	}

	if e == nil {
		err1 := col.Update(
			bson.M{"_id": itemId},
			bson.M{"$inc": bson.M{"real_finOrderCnt": 1}})

		if err1 != nil {
			logs.Error(err1)
		}

		//  判断用户的假订单完成数
		//只有用户在app发布商品或者后台设置了假数据才会有数据
		var shopPlay share_message.TableShopPlayer = share_message.TableShopPlayer{}
		colPlay, closeFunPlay := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_PLAYER)
		defer closeFunPlay()

		ePlay := colPlay.Find(bson.M{"_id": item.GetPlayerId()}).One(&shopPlay)

		if ePlay != nil {
			logs.Error(err1)
		} else {
			if shopPlay.GetFakePlayFinOrderCnt() > 0 {
				shopAuth := share_message.TableShopPlayer{
					PlayerId:            easygo.NewInt64(item.GetPlayerId()),
					UploadAuthFlag:      easygo.NewInt32(1),
					CreateTime:          easygo.NewInt64(time.Now().Unix()),
					FakePlayFinOrderCnt: easygo.NewInt32(shopPlay.GetFakePlayFinOrderCnt() + 1)}

				_, ePlay1 := colPlay.Upsert(
					bson.M{"_id": item.GetPlayerId()},
					bson.M{"$set": shopAuth})

				if ePlay1 != nil {
					logs.Error(e)
				}
			}
		}
	}
}

func AddGoodCommCnt(itemId int64) {

	var item share_message.TableShopItem = share_message.TableShopItem{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()

	e := col.Find(bson.M{"_id": itemId}).One(&item)

	if e != nil {
		logs.Error(e)
	}

	if e == nil {
		if item.GetFakeGoodCommCnt() <= 0 {

			err1 := col.Update(
				bson.M{"_id": itemId},
				bson.M{"$inc": bson.M{"real_goodCommCnt": 1}})

			if err1 != nil {
				logs.Error(err1)
			}
		} else {
			err1 := col.UpdateId(itemId,
				bson.M{"$inc": bson.M{"real_goodCommCnt": 1, "fake_goodCommCnt": 1}})

			if err1 != nil {
				logs.Error(err1)
			}
		}
	}
}

func AddFinCommCnt(itemId int64) {

	var item share_message.TableShopItem = share_message.TableShopItem{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()

	e := col.Find(bson.M{"_id": itemId}).One(&item)

	if e != nil {
		logs.Error(e)
	}

	if e == nil {
		if item.GetFakeFinCommCnt() <= 0 {

			err1 := col.Update(
				bson.M{"_id": itemId},
				bson.M{"$inc": bson.M{"real_finCommCnt": 1}})

			if err1 != nil {
				logs.Error(err1)
			}
		} else {
			err1 := col.UpdateId(itemId,
				bson.M{"$inc": bson.M{"real_finCommCnt": 1, "fake_finCommCnt": 1}})

			if err1 != nil {
				logs.Error(err1)
			}
		}
	}
}

func AddAllCommentCnt(itemId int64) {

	var item share_message.TableShopItem = share_message.TableShopItem{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()

	e := col.Find(bson.M{"_id": itemId}).One(&item)

	if e != nil {
		logs.Error(e)
	}

	if e == nil {
		err1 := col.UpdateId(itemId,
			bson.M{"$inc": bson.M{"real_commentCnt": 1}})

		if err1 != nil {
			logs.Error(err1)
		}
	}
}

func GetItemType(typeName string) int32 {

	switch typeName {
	case "手机":
		return 1
	case "农用物资":
		return 2
	case "生鲜水果":
		return 3
	case "童鞋":
		return 4
	case "园艺植物":
		return 5
	case "五金工具":
		return 6
	case "游泳":
		return 7
	case "电子零件":
		return 8
	case "动漫/周边":
		return 9
	case "图书":
		return 10
	case "宠物/用品":
		return 11
	case "网络设备":
		return 12
	case "服饰配件":
		return 13
	case "家装/建材":
		return 14
	case "家纺布艺":
		return 15
	case "珠宝首饰":
		return 16
	case "钟表眼镜":
		return 17
	case "古董收藏":
		return 18
	case "女士鞋靴":
		return 19
	case "箱包":
		return 20
	case "男士鞋靴":
		return 21
	case "办公用品":
		return 22
	case "游戏设备":
		return 23
	case "运动户外":
		return 24
	case "实体卡/券/票":
		return 25
	case "工艺礼品":
		return 26
	case "玩具乐器":
		return 27
	case "母婴用品":
		return 28
	case "童装":
		return 29
	case "女士服装":
		return 30
	case "家具":
		return 31
	case "居家日用":
		return 32
	case "家用电器":
		return 33
	case "个护美妆":
		return 34
	case "保健护理":
		return 35
	case "摩托车/用品":
		return 36
	case "自行车/用品":
		return 37
	case "汽车/用品":
		return 38
	case "电动车/用品":
		return 39
	case "3C数码":
		return 40
	case "男士服装":
		return 41
	case "其他闲置":
		return 42
	case "音像":
		return 43
	case "演艺/表演类门票":
		return 44
	case "点卡":
		return 45
	default:
		return 0
	}
}

func GetExpressInfos(
	orderId int64,
	com string,
	no string,
	senderPhone string,
	receiverPhone string,
	serviceFlag int32) ([]*share_message.QueryExpressBody, string, string) {

	var expressBodyList ExpressBodyList = ExpressBodyList{}
	cacheExpress := share_message.TableShopCacheExpress{}

	//先从数据库中取得数据
	if serviceFlag == 0 {
		col1, closeFun1 := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_CACHE_EXPRESS)
		defer closeFun1()

		err1 := col1.Find(bson.M{"_id": orderId}).One(&cacheExpress)

		if err1 != mgo.ErrNotFound && err1 != nil {
			return expressBodyList, EXPRESS_QUERY_ERROR_CODE_998, ""
		}
	}

	//然后判断是否过一个小时 没有用缓存的的数据
	if serviceFlag == 0 &&
		cacheExpress.GetExpressList() != nil &&
		len(cacheExpress.GetExpressList()) > 0 &&
		(time.Now().Unix() < (cacheExpress.GetCreateTime()+3600) ||
			cacheExpress.GetExpressList()[0].GetStatus() == share_message.Express_Status_SIGNED) {

		expressBodyList = cacheExpress.GetExpressList()

	} else {
		//从第三方去获得数据
		var sPhoneLastFour string
		var rPhoneLastFour string

		//顺丰物流的时候要判断手机号
		if com == "sf" {

			if senderPhone == "" || len(senderPhone) < 4 {
				return nil, EXPRESS_QUERY_ERROR_CODE_5, "发货人手机号错误"
			}

			if receiverPhone == "" || len(receiverPhone) < 4 {
				return nil, EXPRESS_QUERY_ERROR_CODE_5, "收货人手机号错误"
			}
			sPhoneLastFour = senderPhone[len(senderPhone)-4 : len(senderPhone)]
			rPhoneLastFour = receiverPhone[len(receiverPhone)-4 : len(receiverPhone)]
		}

		netReturn, errCode := ExpressQuery(com, no, sPhoneLastFour, rPhoneLastFour)

		if errCode != "" {

			return nil, errCode, ""
		}

		if nil != netReturn &&
			nil != netReturn["error_code"] &&
			netReturn["error_code"].(float64) != 0 &&
			nil != netReturn["reason"] {

			if netReturn["error_code"].(float64) == 204301 {
				return nil, EXPRESS_QUERY_ERROR_CODE_1, netReturn["reason"].(string)
			} else if netReturn["error_code"].(float64) == 204302 {
				return nil, EXPRESS_QUERY_ERROR_CODE_2, netReturn["reason"].(string)
			} else if netReturn["error_code"].(float64) == 204303 {
				return nil, EXPRESS_QUERY_ERROR_CODE_3, netReturn["reason"].(string)
			} else if netReturn["error_code"].(float64) == 204304 {
				return nil, EXPRESS_QUERY_ERROR_CODE_4, netReturn["reason"].(string)
			} else if netReturn["error_code"].(float64) == 204305 {
				return nil, EXPRESS_QUERY_ERROR_CODE_5, netReturn["reason"].(string)
			} else {
				return nil, easygo.AnytoA(int64(netReturn["error_code"].(float64))), netReturn["reason"].(string)
			}

		} else if nil != netReturn &&
			nil != netReturn["error_code"] &&
			netReturn["error_code"].(float64) == 0 &&
			netReturn["result"] != nil {

			//就算正常返回该值也会出现nil的情况
			statusDetail := netReturn["result"].(map[string]interface{})["status_detail"]

			if netReturn["result"].(map[string]interface{})["list"] != nil {

				for k, value := range netReturn["result"].(map[string]interface{})["list"].([]interface{}) {
					var dateTime interface{}
					var remark interface{}
					if nil != value {
						dateTime = value.(map[string]interface{})["datetime"]
						remark = value.(map[string]interface{})["remark"]
					}
					var expressBody share_message.QueryExpressBody = share_message.QueryExpressBody{
						DateTime: easygo.NewString(dateTime),
						Remark:   easygo.NewString(remark),
					}
					if k == len(netReturn["result"].(map[string]interface{})["list"].([]interface{}))-1 && statusDetail != nil {
						expressBody.Status = GetNewestExpressStatus(statusDetail.(string))
					} else if k == len(netReturn["result"].(map[string]interface{})["list"].([]interface{}))-1 && statusDetail == nil {
						expressBody.Status = GetNewestExpressStatus("PENDING")
					}
					expressBodyList = append(expressBodyList, &expressBody)
				}

				//排序
				sort.Sort(expressBodyList)

				//设置物流的状态
				for _, value := range expressBodyList {

					tempStatus := GetMiddleExpressStatus(value.GetRemark())
					if nil != tempStatus {
						value.Status = tempStatus
					}
				}

				//线程缓存物流信息
				if nil != expressBodyList && len(expressBodyList) > 0 {

					easygo.Spawn(func(orderIdPara int64, expressBodyListPara []*share_message.QueryExpressBody) {

						cacheExpress := share_message.TableShopCacheExpress{
							OrderId:     easygo.NewInt64(orderIdPara),
							ExpressList: expressBodyList,
							CreateTime:  easygo.NewInt64(time.Now().Unix()),
						}

						col2, closeFun2 := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_CACHE_EXPRESS)
						defer closeFun2()
						_, err2 := col2.Upsert(
							bson.M{"_id": orderIdPara},
							bson.M{"$set": cacheExpress})

						if nil != err2 {
							s := fmt.Sprintf("取得物流信息缓存物流信息时数据库出错:%s", err2)
							logs.Error(s)
						}

					}, orderId, expressBodyList)
				}
			}
		}
	}

	return expressBodyList, "", ""
}

func (list ExpressBodyList) Len() int { return len(list) }

func (list ExpressBodyList) Swap(i, j int) {
	s := list[j]
	list[j] = list[i]
	list[i] = s
}

func (list ExpressBodyList) Less(i, j int) bool {

	return list[i].GetDateTime() > list[j].GetDateTime()
}

//订单id生成的初始化
func (self *Shop) InitCreateOrderId() {
	//取得订单的订单id和订单对应的商品id
	var bill share_message.TableBill = share_message.TableBill{}
	var shopOrder share_message.TableShopOrder = share_message.TableShopOrder{}

	colBill, closeFunBill := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_BILLS)
	defer closeFunBill()

	colOrder, closeFunOrder := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFunOrder()

	errBill := colBill.Find(bson.M{}).Sort("-_id").Limit(1).One(&bill)
	errOrder := colOrder.Find(bson.M{}).Sort("-_id").Limit(1).One(&shopOrder)
	if errBill != nil && errBill != mgo.ErrNotFound {
		panic(errBill)
	}
	if errOrder != nil && errOrder != mgo.ErrNotFound {
		panic(errOrder)
	}
	if errBill == mgo.ErrNotFound && errOrder == mgo.ErrNotFound {
		err := easygo.RedisMgr.GetC().StringSet(for_game.SHOP_CREATE_ORDER_ID, util.GetMilliTime()*1000)
		easygo.PanicError(err)
	} else if errBill == nil && errOrder == mgo.ErrNotFound {
		err := easygo.RedisMgr.GetC().StringSet(for_game.SHOP_CREATE_ORDER_ID, bill.GetOrderId())
		easygo.PanicError(err)
	} else if errBill == mgo.ErrNotFound && errOrder == nil {
		err := easygo.RedisMgr.GetC().StringSet(for_game.SHOP_CREATE_ORDER_ID, shopOrder.GetOrderId())
		easygo.PanicError(err)
	} else if errBill == nil && errOrder == nil {
		if bill.GetOrderId() >= shopOrder.GetOrderId() {
			err := easygo.RedisMgr.GetC().StringSet(for_game.SHOP_CREATE_ORDER_ID, bill.GetOrderId())
			easygo.PanicError(err)
		} else {
			err := easygo.RedisMgr.GetC().StringSet(for_game.SHOP_CREATE_ORDER_ID, shopOrder.GetOrderId())
			easygo.PanicError(err)
		}
	}
}
