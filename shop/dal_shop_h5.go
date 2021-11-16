package shop

import (
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"strconv"
	"strings"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"gopkg.in/gomail.v2"
)

//h5详情页返回回数据
type ItemDetailForH5 struct {
	ItemDetail *share_message.ShopItemDetail //商品信息
	SellerInfo *share_message.SellerInfo     //卖家信息
}

//点击立即购买返回给前端的数据
type ItemForNowBuyH5 struct {
	ItemId        *int64                  //id
	NickName      *string                 //卖家昵称
	ItemFile      *share_message.ItemFile //图片信息
	PointCardName *string                 //点卡名称
	Price         *int32                  //单价
}

//点击立即购买返回给前端的数据
type ItemForFirstPay struct {
	OrderId         *int64                           //bill表中的id
	PlayerId        *int64                           //支付人玩家id
	PayChannle      []*share_message.PlatformChannel //支付渠道信息
	SponsorNickname *string                          //卖家昵称
	Items           *share_message.ShopOrderItem     //商品信息
	CreateTime      *int64                           //下单时间
	PayTime         *int64                           //付款时间
}

//点击查询
type OrderForSearch struct {
	Orders []*share_message.TableShopOrder
}

//解密后卡的信息结构
type DecPointCardInfo struct {
	DecCardNo       *string
	DecCardPassword *string //解密后密码
}

//订单详情
type OrderForDetail struct {
	OrderId           *int64                       //订单编号
	SponsorNickname   *string                      //卖家昵称
	Items             *share_message.ShopOrderItem //商品信息
	DecPointCardInfos []*DecPointCardInfo          //解密后的信息
	CreateTime        *int64                       //下单时间
	PayTime           *int64                       //付款时间
}

//id查询同一卖家在架中的同类相关商品
func QueryRelatedShopItem(shopItemPar *ShopItem) (string, []*share_message.TableShopItem) {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()

	var list []*share_message.TableShopItem = make([]*share_message.TableShopItem, 0)
	if shopItemPar != nil {
		queryBson := bson.M{ //"_id": bson.M{"$ne": shopItemPar.item_id}, //自己也要返回给前端
			"type.type": shopItemPar.item_type,
			"player_id": shopItemPar.player_id,
			"state":     for_game.SHOP_ITEM_SALE}

		err := col.Find(queryBson).All(&list)

		if err != nil && err != mgo.ErrNotFound {
			return "操作失败,刷新重试", nil
		}
	}
	return "", list
}

//通过之前详情页的输入内容和商品id去查询订单
func QueryH5Orders(searchTxt string, itemId int64) (string, []*share_message.TableShopOrder) {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()

	var list []*share_message.TableShopOrder = make([]*share_message.TableShopOrder, 0)

	queryBson := bson.M{"h5_search_con": searchTxt,
		"state": for_game.SHOP_ORDER_EVALUTE}

	if itemId > 0 {
		queryBson["items.item_id"] = itemId
	}

	err := col.Find(queryBson).Sort("-create_time").All(&list)

	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return "操作失败,刷新重试", nil
	}
	return "", list
}

//通过之前详情页的输入内容和商品id去查询订单
func QueryH5OrdersNoPay(searchTxt string, itemId int64, count int) (string, *share_message.TableShopOrder) {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()

	var one *share_message.TableShopOrder

	err := col.Find(bson.M{"h5_search_con": searchTxt,
		"items.item_id": itemId,
		"items.count":   count,
		"state":         for_game.SHOP_ORDER_WAIT_PAY}).One(&one)

	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return "操作失败,刷新重试", nil
	}

	if err == mgo.ErrNotFound {
		return "", nil
	}

	return "", one
}

//通过订单id去查询订单
func QueryOrderDetailForH5(orderId int64) (string, *share_message.TableShopOrder) {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()

	var order *share_message.TableShopOrder = &share_message.TableShopOrder{}

	err := col.Find(bson.M{"_id": orderId}).One(order)

	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return "操作失败,刷新重试", nil
	}

	if err == mgo.ErrNotFound {
		return "订单不存在", nil
	}
	return "", order
}

func SendMail(mailTo []string, subject string, body string) error {
	//定义邮箱服务器连接信息，如果是网易邮箱 pass填密码，qq邮箱填授权码

	mailConn := map[string]string{
		"user": easygo.YamlCfg.GetValueAsString("SEND_MAIL_ADDRESS"),
		"pass": easygo.YamlCfg.GetValueAsString("SEND_MAIL_PASS"),
		"host": easygo.YamlCfg.GetValueAsString("SEND_MAIL_HOST"),
		"port": easygo.YamlCfg.GetValueAsString("SEND_MAIL_PORT"),
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()

	m.SetAddressHeader("From", mailConn["user"], "点卡购买")
	m.SetHeader("To", mailTo...)    //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err

}

func DoSendMail(orderId int64) string {
	logs.Info("========DoSendMail=======入口，orderId:", orderId)
	//再次查询是为了取得绑定点卡的信息
	order := share_message.TableShopOrder{}

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()
	e := col.Find(bson.M{"_id": orderId}).One(&order)

	if e != nil && e != mgo.ErrNotFound {
		logs.Error("查询数据库失败 TABLE_SHOP_ORDERS，err:", e)
		return "操作失败"
	}

	if e == mgo.ErrNotFound {
		logs.Error("订单不存在", e)
		return "订单不存在"
	}
	if order.GetH5SearchCon() != "" &&
		strings.Contains(order.GetH5SearchCon(), "@") &&
		order.GetItems() != nil &&
		order.GetItems().GetItemType() == for_game.SHOP_POINT_CARD_CATEGORY &&
		order.GetItems().GetPointCardInfos() != nil &&
		len(order.GetItems().GetPointCardInfos()) > 0 {

		//定义收件人
		mailTo := []string{
			order.GetH5SearchCon(),
		}
		//邮件主题为"点卡购买成功提醒"
		subject := "点卡购买成功提醒"
		// 邮件正文
		body := fmt.Sprintf("您已经成功购买%v点卡<br>点卡信息:<br>", order.GetItems().GetPointCardName())

		for _, value := range order.GetItems().GetPointCardInfos() {
			//对key进行加密后取得解密的真正Key
			passKey := for_game.AesEncrypt(for_game.AES_KEY_SHOP_CARD, []byte(value.GetKey()))
			//对密码解密
			bytesPass, err := for_game.AesDecrypt(value.GetCardPassword(), []byte(passKey))

			if err != nil {
				logs.Error("对密码解密--->", err)
			}

			body = body + fmt.Sprintf("卡号:%v<br>卡密:%v<br><br>", value.GetCardNo(), string(bytesPass))

		}

		err := SendMail(mailTo, subject, body)
		if err != nil {
			logs.Error(err)
			s := fmt.Sprintf("点卡订单id%v对应的商品id%v邮件%v发送失败", order.GetOrderId(), order.GetItems().GetItemId(), order.GetH5SearchCon())
			logs.Error(s)
			return "邮件发送失败"
		}
		logs.Info("发送邮件成功")

	} else {

		logs.Info("条件不满足的日志，order.GetH5SearchCon()：%v,order.GetH5SearchCon():%v,order.GetItems():%v,order.GetItems().GetItemType():%v,order.GetItems().GetPointCardInfos():%v,len(order.GetItems().GetPointCardInfos()):%v",
			order.GetH5SearchCon(), order.GetH5SearchCon(), order.GetItems(), order.GetItems().GetItemType(), order.GetItems().GetPointCardInfos(), len(order.GetItems().GetPointCardInfos()))
	}

	return ""
}

func GetShopItemFromDB(itemId int64) (string, *ShopItem) {

	errStr := ""
	var shopItem *ShopItem

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()
	var item share_message.TableShopItem = share_message.TableShopItem{}

	e := col.Find(bson.M{"_id": itemId, "state": for_game.SHOP_ITEM_SALE, "stock_count": bson.M{"$gt": 0}}).One(&item)
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

	return errStr, shopItem
}
