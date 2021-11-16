// 为[浏览器]提供的API服务

package shop

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

type Result struct {
	Ret    int
	Reason string
	Data   interface{}
}

const API_MD5_KEY = "nmcl2020!@#$"

type WebHttpServer struct {
	Service reflect.Value
}

func NewWebHttpServer() *WebHttpServer {
	p := &WebHttpServer{}
	p.Init()
	return p
}
func (self *WebHttpServer) Init() {

}
func (self *WebHttpServer) Serve() {
	port := easygo.YamlCfg.GetValueAsInt("LISTEN_ADDR_FOR_WEB_API")
	address := for_game.MakeAddress("0.0.0.0", int32(port))

	logs.Info("(API 服务) 开始监听: %v", address)

	http.HandleFunc("/shop", self.ShopEntry) //商城api

	err := http.ListenAndServe(address, nil) // 第 2 个参数可以传 nil。传进去的 self 必须有实现一个 ServeHTTP 函数
	easygo.PanicError(err)
}

func (self *WebHttpServer) ShopEntry(w http.ResponseWriter, r *http.Request) {
	logs.Info("=====ShopEntry=====")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	params := r.Form
	if !self.ApiCheckSign(params) {
		OutputJson(w, 0, "签名错误", nil)
		return
	}

	t := params.Get("t")
	logs.Info("=====ShopEntry=====t=", t)
	//请求数据类型t:
	// 1查询商品详情数据包装(链接打开等动作),
	// 2用订单号拿支付配置
	// 3用户下订单
	// 4搜索按钮
	// 5支付后的订单详情
	switch t {
	case for_game.SHOP_API_ITEM_DETAIL:
		self.QueryShopItemDetail(w, params)
	case for_game.SHOP_API_NOW_BUY:
		self.DoShopNowBuy(w, params)
	case for_game.SHOP_API_FIRST_PAY:
		self.DoShopFirstPay(w, params)
	case for_game.SHOP_API_SERACH:
		self.DoShopSearch(w, params)
	case for_game.SHOP_API_ORDER_DETAIL:
		self.QueryShopOrderDetail(w, params)
	default:
		OutputJson(w, 0, "请求类型错误", nil)
	}
}

//返回json
func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := &Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	w.Write(b)
}

//请求验签
func (self *WebHttpServer) ApiCheckSign(params url.Values) bool {
	sign := params.Get("sign")
	// logs.Info("================>前端签名", sign)
	params.Del("sign")
	src := params.Encode()
	newSign := for_game.Md5(src + API_MD5_KEY)
	// logs.Info("================>服务器签名", src+API_MD5_KEY, newSign)
	return sign == newSign
}

//1查询商品详情数据包装(链接打开等动作),
func (self *WebHttpServer) QueryShopItemDetail(w http.ResponseWriter, params url.Values) {
	id := params.Get("id")
	var itemId int64 = easygo.StringToInt64noErr(id)

	if itemId == 0 {
		OutputJson(w, 0, "商品ID错误", nil)
		return
	}

	//取得商品(下架了也能取到,参数传卖家)
	errStr, shopItem := ShopInstance.GetShopItem(share_message.BuySell_Type_Seller, itemId)
	//判断是否err
	if errStr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}

	//判断是否点卡类型的商品
	if shopItem.item_type != for_game.SHOP_POINT_CARD_CATEGORY {
		OutputJson(w, 0, "商品ID错误", nil)
		return
	}

	//封装商品详情的返回
	sellerInfo := ShopInstance.SellerInfo(shopItem)
	itemDetail := ShopInstance.ItemDetail(shopItem)
	if nil == sellerInfo {
		OutputJson(w, 0, "商品已经下架", nil)
		return
	}
	if nil == itemDetail {
		OutputJson(w, 0, "商品已经下架", nil)
		return
	}

	//设置商品信息中的同类商品
	relatedErr, relatedList := QueryRelatedShopItem(shopItem)
	if relatedErr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}

	var tempRelatedShopItems []*share_message.RelatedShopItem = make([]*share_message.RelatedShopItem, 0)

	if nil != relatedList && len(relatedList) > 0 {
		for _, value := range relatedList {

			tempRelatedShopItem := &share_message.RelatedShopItem{
				ItemId:    easygo.NewInt64(value.GetItemId()),
				ItemFiles: value.GetItemFiles(),
			}

			tempRelatedShopItems = append(tempRelatedShopItems, tempRelatedShopItem)
		}
	}

	//设置同类
	itemDetail.RelatedShopItems = tempRelatedShopItems
	result := &ItemDetailForH5{ItemDetail: itemDetail, SellerInfo: sellerInfo}
	OutputJson(w, 1, "success", result)

	//记录当天浏览量
	easygo.Spawn(func(pageViewFlag int32, item *ShopItem) {
		if item != nil {
			ShopInstance.UpdatePageViews(pageViewFlag, item)
		} else {
			logs.Debug("记录浏览量时缺少商品信息")
		}
	}, int32(0), shopItem)

	return
}

//2用订单号拿支付配置
func (self *WebHttpServer) DoShopNowBuy(w http.ResponseWriter, params url.Values) {
	id := params.Get("id")
	var orderId int64 = easygo.StringToInt64noErr(id)

	if orderId == 0 {
		OutputJson(w, 0, "订单ID错误", nil)
		return
	}

	//通过订单id取得订单详情
	err, order := QueryOrderDetailForH5(orderId)
	if err != "" {
		OutputJson(w, 0, err, nil)
		return
	}
	if nil == order {
		OutputJson(w, 0, "订单不存在", nil)
		return
	}

	if order.GetState() == for_game.SHOP_ORDER_EXPIRE {
		OutputJson(w, 0, "订单已超时", nil)
		return
	}

	if order.GetState() == for_game.SHOP_ORDER_EVALUTE {
		OutputJson(w, 0, "订单已完成", nil)
		return
	}

	if order.GetState() != for_game.SHOP_ORDER_WAIT_PAY {
		OutputJson(w, 0, "订单超时", nil)
		return
	}

	//查询h5开启
	lis := for_game.QueryPlatformChannelList(2)
	plist := []*share_message.PlatformChannel{}
	for _, cl := range lis {
		if cl.GetTypes() == 1 { //查找入款的渠道
			plist = append(plist, cl)
		}
	}

	result := &ItemForFirstPay{
		OrderId:         easygo.NewInt64(orderId),
		SponsorNickname: easygo.NewString(order.GetSponsorNickname()),
		Items:           order.Items,
		CreateTime:      easygo.NewInt64(order.GetCreateTime()),
		PayTime:         easygo.NewInt64(order.GetPayTime()),
		PlayerId:        easygo.NewInt64(order.GetReceiverId()),
		PayChannle:      plist,
	}

	OutputJson(w, 1, "success", result)
}

//3用户下订单,
func (self *WebHttpServer) DoShopFirstPay(w http.ResponseWriter, params url.Values) {
	cnt := params.Get("count")
	var count int = easygo.StringToIntnoErr(cnt)

	if count == 0 {
		OutputJson(w, 0, "数量错误", nil)
		return
	}

	//保险起见,再做一次验证
	//页面输入的text,后台做一次正确性验证
	searchText := params.Get("searchText")
	if "" == searchText {
		OutputJson(w, 0, "请输入内容", nil)
		return
	} else {
		//邮箱判断
		if !strings.Contains(searchText, "@") {
			OutputJson(w, 0, "输入正确的邮箱", nil)
			return
		}
		//判断如果不是邮箱,必须是纯数字(把手机号也兼容做了)
		if !strings.Contains(searchText, "@") && easygo.StringToInt64noErr(searchText) == 0 {
			OutputJson(w, 0, "详情页输入内容错误非法", nil)
			return
		}
	}

	id := params.Get("id")
	var itemId int64 = easygo.StringToInt64noErr(id)

	//为了保险起见,这里再做一次验证
	if itemId == 0 {
		OutputJson(w, 0, "商品ID错误", nil)
		return
	}

	//数据库取得在架有库存的商品
	errStr, shopItemFromCache := GetShopItemFromDB(itemId)
	//判断商品是否下架
	if errStr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}

	//判断是否点卡类型的商品
	if shopItemFromCache.item_type != for_game.SHOP_POINT_CARD_CATEGORY {
		OutputJson(w, 0, "商品ID错误", nil)
		return
	}

	if nil == shopItemFromCache {
		OutputJson(w, 0, DETAIL_SHOP_ITEM_NOT_EXIST, nil)
		return
	}

	//创建订单
	logs.Info("H5创建订单开始")

	//库存
	if shopItemFromCache.stock_count < int32(count) {
		s := fmt.Sprintf("商品：库存不足,还剩%v库存",
			easygo.AnytoA(int64(shopItemFromCache.stock_count)))
		OutputJson(w, 0, s, nil)
		return
	}

	//如果是点卡 再判断一次实际导入库中的库存是否足够
	if shopItemFromCache.item_type == for_game.SHOP_POINT_CARD_CATEGORY {

		//通过物品的个数取得卡密个数
		pointCardList := for_game.GetPointCardByBuyInfos(shopItemFromCache.account, shopItemFromCache.pointCardName, int32(count))
		if nil != pointCardList && (len(pointCardList) == 0 || int32(len(pointCardList)) < int32(count)) {

			s := fmt.Sprintf("商品：库存不足,还剩%v库存",
				easygo.AnytoA(int64(len(pointCardList))))
			OutputJson(w, 0, s, nil)
			return
		}
	}

	err, order := QueryH5OrdersNoPay(searchText, itemId, count)
	if err != "" {
		OutputJson(w, 0, err, nil)
		return
	}
	if order != nil {
		result := &OrderForDetail{
			OrderId: easygo.NewInt64(order.GetOrderId()),
		}
		OutputJson(w, 1, "success", result)
		return
	}

	now := time.Now().Unix()
	var state int32 = for_game.SHOP_ORDER_WAIT_PAY
	var deleteBuy int32 = 0
	var deleteSell int32 = 0
	var delayReceive int32 = 0

	colOrder, closeFunOrder := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFunOrder()

	if shopItemFromCache != nil {

		orderId := ShopInstance.CreateOrderID()

		itemFile := share_message.ItemFile{
			FileUrl:  &shopItemFromCache.item_files[0].file_url,
			FileType: &shopItemFromCache.item_files[0].file_type}

		item := share_message.ShopOrderItem{
			ItemId:        easygo.NewInt64(itemId),
			Name:          easygo.NewString(shopItemFromCache.name),
			Price:         easygo.NewInt32(shopItemFromCache.price),
			ItemFile:      &itemFile,
			Count:         easygo.NewInt32(count),
			Title:         easygo.NewString(shopItemFromCache.title),
			PointCardName: easygo.NewString(shopItemFromCache.pointCardName),
			CopyName:      easygo.NewString(shopItemFromCache.name),
			ItemType:      easygo.NewInt32(shopItemFromCache.item_type),
		}

		//点卡的时候,为了客户端不修改代码,把商品名称设置为点卡名称
		if shopItemFromCache.item_type == for_game.SHOP_POINT_CARD_CATEGORY {
			item.Name = easygo.NewString(shopItemFromCache.pointCardName)
		}

		//通过详情页面输入的内容==============================================================开始创建用户
		playAccountInfo := for_game.GetRedisAccountByPhone(searchText)
		if nil == playAccountInfo {

			data := &share_message.CreateAccountData{
				Phone:    easygo.NewString(searchText),
				PassWord: easygo.NewString(""),
				Ip:       easygo.NewString(""),
				IsOnline: easygo.NewBool(false),
			}
			b, _ := for_game.CreateAccount(data)
			if !b {
				OutputJson(w, 0, "操作失败,刷新重试", nil)
				return
			}
			playAccountInfo = for_game.GetRedisAccountByPhone(searchText)

			// player := for_game.GetRedisPlayerBase(playerId)
			// player.SetNickName(for_game.GetRandNickName())
			// player.SetDeviceType(3)
			// player.SetCreateTime()
			// player.SetSex(2)
			// player.SetHeadIcon(for_game.GetDefaultHeadicon(2))
			// player.SaveToMongo()
		}

		//取得playerBase信息
		playerBaseInfo := for_game.GetRedisPlayerBase(playAccountInfo.GetPlayerId())
		//===================================================================================结束创建用户
		order := share_message.TableShopOrder{
			OrderId:          easygo.NewInt64(orderId),
			SponsorId:        easygo.NewInt64(shopItemFromCache.player_id),
			SponsorSex:       easygo.NewInt32(shopItemFromCache.sex),
			SponsorAvatar:    easygo.NewString(shopItemFromCache.avatar),
			SponsorNickname:  easygo.NewString(shopItemFromCache.nickname),
			ReceiverId:       easygo.NewInt64(playerBaseInfo.GetPlayerId()),
			ReceiverSex:      easygo.NewInt32(playerBaseInfo.GetSex()),
			ReceiverNickname: easygo.NewString(playerBaseInfo.GetNickName()),
			ReceiverAvatar:   easygo.NewString(playerBaseInfo.GetHeadIcon()),
			Items:            &item,
			State:            easygo.NewInt32(state),
			DelayReceive:     easygo.NewInt32(delayReceive),
			DeleteBuy:        easygo.NewInt32(deleteBuy),
			DeleteSell:       easygo.NewInt32(deleteSell),
			DeliverAddress: &share_message.DeliverAddress{
				Name:          easygo.NewString(shopItemFromCache.userName),
				Phone:         easygo.NewString(shopItemFromCache.phone),
				Region:        easygo.NewString(shopItemFromCache.address),
				DetailAddress: easygo.NewString(shopItemFromCache.detail_address),
			},
			CreateTime:         easygo.NewInt64(now),
			ReceiveTime:        easygo.NewInt64(now),
			SponsorAccount:     easygo.NewString(shopItemFromCache.account),
			ReceiverAccount:    easygo.NewString(playerBaseInfo.GetAccount()),
			ReceiverNotifyFlag: easygo.NewBool(true),
			SponsorNotifyFlag:  easygo.NewBool(true),
			UpdateTime:         easygo.NewInt64(now),
			H5SearchCon:        easygo.NewString(searchText),
		}

		e := colOrder.Insert(order)

		if e != nil {
			logs.Error(e)
			OutputJson(w, 0, "操作失败,刷新重试", nil)
			return
		}

		//h5创建订单不通知app用户完成支付通知
		////创建订单就通知
		////订单通知
		//easygo.Spawn(func(orderParam share_message.TableShopOrder) {
		//	if &orderParam != nil {
		//
		//		var content string = MESSAGE_TO_SELLER_NEW
		//		typeValue := share_message.BuySell_Type_Seller
		//
		//		//商城消息通知
		//		ShopInstance.InsMessageNotify(
		//			easygo.NewString(content),
		//			&typeValue,
		//			&orderParam)
		//
		//		//商城极光推送
		//		var jgContent string = MESSAGE_TO_SELLER_NEW_PUSH
		//		ShopInstance.JGMessageNotify(jgContent, orderParam.GetSponsorId(), orderParam.GetOrderId(), typeValue)
		//
		//		//创建订单通知买家 商城订单红点推送
		//		SendToHallClient(orderParam.GetReceiverId(), "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
		//			OrderId: easygo.NewInt64(orderParam.GetOrderId())})
		//		/*	SendToPlayer(orderParam.GetReceiverId(), "RpcShopOrderNotify",
		//			&share_message.ShopOrderNotifyInfoWithWho{
		//				PlayerId: easygo.NewInt64(orderParam.GetReceiverId()),
		//				OrderId:  easygo.NewInt64(orderParam.GetOrderId()),
		//			})*/
		//		//创建订单通知卖家 商城订单红点推送
		//		SendToHallClient(orderParam.GetSponsorId(), "RpcShopOrderNotify", &share_message.ShopOrderNotifyInfo{
		//			OrderId: easygo.NewInt64(orderParam.GetOrderId())})
		//		/*	SendToPlayer(orderParam.GetSponsorId(), "RpcShopOrderNotify",
		//			&share_message.ShopOrderNotifyInfoWithWho{
		//				PlayerId: easygo.NewInt64(orderParam.GetSponsorId()),
		//				OrderId:  easygo.NewInt64(orderParam.GetOrderId()),
		//			})*/
		//
		//	} else {
		//		logs.Debug("H5创建订单发通知,缺少订单")
		//	}
		//}, order)

		logs.Info("H5创建订单成功")
		//easygo.Spawn(func() {
		//	for_game.SetRedisOperationChannelReportFildVal(util.GetMilliTime(), 1, playerBaseInfo.GetChannel(), "ShopOrderCount") //渠道汇总报表添加下单数量
		//})
		result := &OrderForDetail{
			OrderId: easygo.NewInt64(order.GetOrderId()),
		}
		OutputJson(w, 1, "success", result)
	}
}

//4搜索按钮
func (self *WebHttpServer) DoShopSearch(w http.ResponseWriter, params url.Values) {
	//保险起见,再做一次验证
	//页面输入的text,后台做一次正确性验证
	searchText := params.Get("searchText")
	if "" == searchText {
		OutputJson(w, 0, "请输入搜索条件", nil)
		return
	} else {
		//邮箱判断
		if !strings.Contains(searchText, "@") {
			OutputJson(w, 0, "输入正确的邮箱", nil)
			return
		}
		//判断如果不是邮箱,必须是纯数字(把手机号也兼容做了)
		if !strings.Contains(searchText, "@") && easygo.StringToInt64noErr(searchText) == 0 {
			OutputJson(w, 0, "详情页输入内容错误非法", nil)
			return
		}
	}

	id := params.Get("id")
	var itemId int64 = easygo.StringToInt64noErr(id)

	//通过之前详情页的输入内容和商品id去查询订单
	err, orders := QueryH5Orders(searchText, itemId)
	if err != "" {
		OutputJson(w, 0, err, nil)
		return
	}

	result := &OrderForSearch{
		Orders: orders,
	}

	OutputJson(w, 1, "success", result)
}

//5支付后的订单详情
func (self *WebHttpServer) QueryShopOrderDetail(w http.ResponseWriter, params url.Values) {
	id := params.Get("id")
	var orderId int64 = easygo.StringToInt64noErr(id)

	if orderId == 0 {
		OutputJson(w, 0, "订单ID错误", nil)
		return
	}

	//通过订单id取得订单详情
	err, order := QueryOrderDetailForH5(orderId)
	if err != "" {
		OutputJson(w, 0, err, nil)
		return
	}
	if nil == order {
		OutputJson(w, 0, "订单不存在", nil)
		return
	}

	result := &OrderForDetail{
		OrderId:         easygo.NewInt64(order.GetOrderId()),
		SponsorNickname: easygo.NewString(order.GetSponsorNickname()),
		Items:           order.Items,
		CreateTime:      easygo.NewInt64(order.GetCreateTime()),
		PayTime:         easygo.NewInt64(order.GetPayTime()),
	}

	if order.GetState() == for_game.SHOP_ORDER_EVALUTE {
		var decPointCardInfos []*DecPointCardInfo = make([]*DecPointCardInfo, 0)

		//点卡信息解密
		if order.GetItems() != nil &&
			order.GetItems().GetPointCardInfos() != nil &&
			len(order.GetItems().GetPointCardInfos()) > 0 {
			for _, value := range order.GetItems().GetPointCardInfos() {
				//对key进行加密后取得解密的真正Key
				passKey := for_game.AesEncrypt(for_game.AES_KEY_SHOP_CARD, []byte(value.GetKey()))
				//对密码解密
				bytesPass, err := for_game.AesDecrypt(value.GetCardPassword(), []byte(passKey))

				if err != nil {
					logs.Error(err)
					OutputJson(w, 0, "操作失败,刷新重试", nil)
					return
				}
				pointCardInfo := DecPointCardInfo{
					DecCardNo:       easygo.NewString(value.GetCardNo()),
					DecCardPassword: easygo.NewString(string(bytesPass)),
				}
				decPointCardInfos = append(decPointCardInfos, &pointCardInfo)
			}

			result.DecPointCardInfos = decPointCardInfos
		}
	}

	OutputJson(w, 1, "success", result)
}
