package shop

import (
	"game_server/easygo"
	"time"
)

func InitTest() {
	easygo.Spawn(TestUpdate)
}

func TestUpdate() {

	for {

		//ShopItemListTest()
		//col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
		//players := []*share_message.PlayerBase{}
		//err := col.Find(bson.M{"BlackList": 6}).All(&players)
		//closeFun()
		//if(err == nil && len(players) > 0){
		//	for _, black := range players {
		//		logs.Info(black.GetPlayerId())
		//	}
		//
		//}

		//ShopItemAddItemToStoreTest()
		time.Sleep(10 * time.Second)
		//logs.Debug()
	}
}

//
//func ShopItemAddItemToStoreTest() {
//	var itemId int64 = 42
//	var playId int64 = 421887436005
//	var who Player = Player{}
//	who.PlayerId = playId
//	var reqMsg *share_message.ShopItemID = &share_message.ShopItemID{
//		ItemId: &itemId,
//	}
//
//	instance := ServiceForGameClient{}
//	logs.Debug(instance.RpcAddItemToStore(nil, &who, reqMsg))
//}

// 测试用例文件，用于日常测试，别删
//func ShopItemListTest() {
//	var t int32 = 0
//	var page int32 = 0
//	var page_size int32 = 20
//	var reqMsg *share_message.ShopInfo = &share_message.ShopInfo{Page: &page, PageSize: &page_size, Type: &t}
//
//	instance := ServiceForHall{}
//	logs.Debug(instance.RpcShopItemList(nil, nil, reqMsg))
//}

/*
func ShopItemUploadTest() {
	var title = "带上飞机哦"
	var image string = "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1575886058607&di=7e5e61f23425fbed91c1d59230fc6df8&imgtype=0&src=http%3A%2F%2Fwww.bkill.com%2Fu%2Fupload%2F2018%2F11%2F01%2F012157281452.jpg"
	var file_type int32 = 0
	var item_file share_message.ItemFile = share_message.ItemFile{
		FileUrl:  &image,
		FileType: &file_type,
	}
	var item_files []*share_message.ItemFile = []*share_message.ItemFile{&item_file}
	var price int32 = 12
	var origin_price int32 = 12
	var address = "dasd"
	var detail_address = "dsadas"
	var player_id int64 = 12
	var nickname = "啦啦啦"
	var avata = ""
	var item_type int32 = 1
	var other_type = []string{"1", "2"}
	var stock_count int32 = 5
	var t = share_message.ShopItemType{Type: &item_type, OtherType: other_type}
	var reqMsg *share_message.ShopItemUploadWithWhoInfo = &share_message.ShopItemUploadWithWhoInfo{Title: &title, ItemFiles: item_files,
		Price: &price, OriginPrice: &origin_price, Address: &address, PlayerId: &player_id,
		Nickname: &nickname, Avata: &avata, Type: &t, StockCount: &stock_count, DetailAddress: &detail_address}

	instance := ServiceForGameClient{}
	logs.Debug(instance.RpcShopItemUpload(nil, nil, reqMsg))
}

func ShopItemCommentUploadTest() {
	var itemId int32 = 2
	var content = "测试留言上传功能"
	var player_id int64 = 12
	var avatar = "http://img.huabanxiu.com/system/2019/7/22/img1563771744301_770.jpg"
	var nickname = "测试1"

	var reqMsg *share_message.UploadCommentWithWhoInfo = &share_message.UploadCommentWithWhoInfo{
		ItemId: &itemId, Content: &content, PlayerId: &player_id, Avatar: &avatar, Nickname: &nickname,
	}

	instance := ServiceForGameClient{}
	logs.Debug(instance.RpcShopItemCommentUpload(nil, nil, reqMsg))
}

func ShopItemEditTest() {
	var title = "21123"
	var image string = "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1575886058607&di=7e5e61f23425fbed91c1d59230fc6df8&imgtype=0&src=http%3A%2F%2Fwww.bkill.com%2Fu%2Fupload%2F2018%2F11%2F01%2F012157281452.jpg"
	var file_type int32 = 0
	var item_file share_message.ItemFile = share_message.ItemFile{
		FileUrl:  &image,
		FileType: &file_type,
	}
	var item_files []*share_message.ItemFile = []*share_message.ItemFile{&item_file}
	var price int32 = 12
	var origin_price int32 = 12
	var address = "dasddasdas"
	var item_id int32 = 2
	var t = share_message.ShopItemType{}
	var info = share_message.ShopItemUploadInfo{Title: &title, ItemFiles: item_files, Price: &price, OriginPrice: &origin_price, Address: &address, Type: &t}
	var reqMsg *share_message.ShopItemEditInfo = &share_message.ShopItemEditInfo{Info: &info, ItemId: &item_id}

	instance := ServiceForGameClient{}
	logs.Debug(instance.RpcShopItemEdit(nil, nil, reqMsg))
}

//func ShopItemShowDetailTest() {
//
//	var item_id int32 = 4
//	var flag int32 = 1
//	var reqMsg *share_message.ShopItemInfo = &share_message.ShopItemInfo{ItemId: &item_id, Flag: &flag}
//
//	instance := ServiceForHall{}
//	logs.Debug(instance.RpcShopItemInfo(nil, nil, reqMsg))
//}

func CreateOrderTest() {
	var item_id int32 = 32
	var remark string = "111"
	items := []*share_message.BuyItemID{&share_message.BuyItemID{ItemId: &item_id, Remark: &remark}}
	var player_id int64 = 5
	var reqMsg *share_message.BuyItemInfoWithWho = &share_message.BuyItemInfoWithWho{Items: items, PlayerId: &player_id}

	instance := ServiceForGameClient{}
	logs.Debug(instance.RpcCreateOrder(nil, nil, reqMsg))
}

func AddItemToCartTest() {

	var item_id int32 = 16
	var player_id int64 = 100
	var reqMsg *share_message.ShopItemIDWithWhoInfo = &share_message.ShopItemIDWithWhoInfo{ItemId: &item_id,
		PlayerId: &player_id}

	instance := ServiceForGameClient{}
	logs.Debug(instance.RpcAddItemToCart(nil, nil, reqMsg))
}

func RemoveItemFromCartTest() {

	item_id := []int32{16, 15}
	var player_id int64 = 18874369
	var reqMsg *share_message.RemoveCartWithWhoInfo = &share_message.RemoveCartWithWhoInfo{ItemIds: item_id,
		PlayerId: &player_id}

	instance := ServiceForGameClient{}
	logs.Debug(instance.RpcRemoveItemFromCart(nil, nil, reqMsg))
}

func TestReceiveAddressAdd() {
	var player_id int64 = 2
	var name = "dad"
	var phone = "dasdas"
	var region = "发生发射点"
	var detail_address = "房贷首付"
	var address share_message.ReceiveAddress = share_message.ReceiveAddress{Name: &name, Phone: &phone, Region: &region, DetailAddress: &detail_address}
	reqMsg := &share_message.ReceiveAddressWithWho{PlayerId: &player_id, Address: &address}

	instance := ServiceForGameClient{}
	logs.Debug(instance.RpcReceiveAddressAdd(nil, nil, reqMsg))
}

func TestRpcReceiveAddressList() {
	var player_id int64 = 2

	reqMsg := &share_message.PlayerID{PlayerId: &player_id}

	instance := ServiceForGameClient{}
	logs.Debug(instance.RpcReceiveAddressList(nil, nil, reqMsg))
}

func TestReceiveAddressEdit() {
	var address_id int32 = 15794
	var name = "dad"
	var phone = "1111"
	var region = "1111"
	var detail_address = "1111"
	var address share_message.ReceiveAddress = share_message.ReceiveAddress{Name: &name, Phone: &phone, Region: &region, DetailAddress: &detail_address}
	reqMsg := &share_message.ReceiveAddressInfo{AddressId: &address_id, Address: &address}

	instance := ServiceForGameClient{}
	logs.Debug(instance.RpcReceiveAddressEdit(nil, nil, reqMsg))
}

func TestReceiveAddressDelete() {
	var address_id int32 = 15794

	reqMsg := &share_message.ReceiveAddressID{AddressId: &address_id}

	instance := ServiceForGameClient{}
	logs.Debug(instance.RpcReceiveAddressDelete(nil, nil, reqMsg))
}

func TestRpcOrderList() {
	var player_id int64 = 1
	var tt int32 = 1
	var item_type int32 = 0

	reqMsg := &share_message.OrderInfoWithWho{PlayerId: &player_id, Type: &tt, ItemType: &item_type}

	instance := ServiceForGameClient{}
	logs.Debug(instance.RpcOrderList(nil, nil, reqMsg))
}

func TestRpcShopItemEvaluteUpload() {
	var item_id int32 = 1
	var context string = "1111"
	var player_id int64 = 2
	var avata string = "1111"
	var nickname string = "1111"
	var order_id int64 = 1578559621001
	reqMsg := &share_message.UploadEvaluteWithWhoInfo{ItemId: &item_id, Content: &context, Avatar:&avata, PlayerId:&player_id,Nickname:&nickname, OrderId:&order_id}

	instance := ServiceForGameClient{}
	logs.Debug(instance.RpcShopItemEvaluteUpload(nil, nil, reqMsg))
}
*/
