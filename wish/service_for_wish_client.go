// 大厅服务器为[游戏客户端]提供的服务

package wish

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/h5_wish"
	"time"

	"github.com/astaxie/beego/logs"
)

// 1-设置守护者,2-取消守护者
const (
	WISH_BACKSTAGE_SET_GUARDIAN    = 1 //  1-设置守护者
	WISH_BACKSTAGE_CANCEL_GUARDIAN = 2 // 2-取消守护者
)

//许愿池硬币修改
//func (self *WebHttpForClient) RpcAddCoin(common *base.Common, reqMsg *h5_wish.AddCoinReq) easygo.IMessage {
//	base := for_game.GetRedisWishPlayer(common.GetUserId())
//	if base == nil {
//		logs.Error("获取许愿池用户对象失败,许愿池用户id为: %d", common.GetUserId())
//		return easygo.NewFailMsg("用户信息有误")
//	}
//	resp, err := SendMsgToIdelServer(for_game.SERVER_TYPE_HALL, "RpcAddCoin", reqMsg, base.GetPlayerId())
//	if err != nil {
//		logs.Error("err:", err)
//	}
//	return resp
//}

//  用户零钱变化
//func (self *WebHttpForClient) RpcAddGold(common *base.Common, reqMsg *h5_wish.AddGoldReq) easygo.IMessage {
//	base := for_game.GetRedisWishPlayer(common.GetUserId())
//	if base == nil {
//		logs.Error("获取许愿池用户对象失败,许愿池用户id为: %d", common.GetUserId())
//		return easygo.NewFailMsg("用户信息有误")
//	}
//	resp, err := SendMsgToIdelServer(for_game.SERVER_TYPE_HALL, "RpcAddGold", reqMsg, base.GetPlayerId())
//	if err != nil {
//		logs.Error("err:", err)
//	}
//	return resp
//}

//获取用户信息:分2部分：一部分是许愿池wish_player里面有基础信息，详细信息获取的话，请到大厅获取
func (self *WebHttpForClient) RpcGetUseInfo(common *base.Common, reqMsg *h5_wish.UserInfoReq) easygo.IMessage {
	//TODO
	base := for_game.GetRedisWishPlayer(common.GetUserId())
	if base == nil {
		logs.Error("获取许愿池用户对象失败,许愿池用户id为: %d", common.GetUserId())
		return easygo.NewFailMsg("用户信息有误")
	}

	resp, err := SendMsgToIdelServer(for_game.SERVER_TYPE_HALL, "RpcGetUseInfo", reqMsg, base.GetPlayerId())
	if err != nil {
		logs.Error("err:", err)
	}
	data := resp.(*h5_wish.UserInfoResp)

	// 判断用户是否有月首冲
	var b bool // true-月首充,false-不是月首充
	diamondTime := base.GetLastExchangeDiamondTime()
	m1 := int(time.Unix(diamondTime/1000, 0).Month())
	m2 := int(time.Unix(for_game.GetMillSecond()/1000, 0).Month())
	if m1 != m2 {
		b = true
	}
	data.HeadUrl = easygo.NewString(base.GetHeadIcon())
	data.Name = easygo.NewString(base.GetNickName())
	data.Diamond = easygo.NewInt64(base.GetDiamond())
	data.IsFirst = easygo.NewBool(b)

	return data
}

// 许愿池首页查询(搜索功能)  (已完成)
func (self *WebHttpForClient) RpcQueryBox(common *base.Common, reqMsg *h5_wish.QueryBoxReq) easygo.IMessage {
	logs.Info("===许愿池首页查询 RpcQueryBox pid:%v, reqMsg: %v", common.GetUserId(), reqMsg)
	t := reqMsg.GetType()
	if t == 0 {
		logs.Error("许愿池首页查询 type为空")
		return easygo.NewFailMsg("参数有误")
	}
	if reqMsg.GetId() == 0 {
		logs.Error("许愿池首页查询 id为空")
		return easygo.NewFailMsg("参数有误")
	}
	resp, err := QueryBoxService(reqMsg)
	return RpcReturnCommon("RpcQueryBox", resp, err)
}

// 所有商品的名字或者盲盒的名字(已完成)
func (self *WebHttpForClient) RpcQueryBoxProductName(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===所有商品的名字或者盲盒的名字 RpcQueryBoxProductName pid:%v", common.GetUserId())
	resp, err := QueryBoxProductNameService()
	return RpcReturnCommon("RpcQueryBoxProductName", resp, err)
}

// 搜索发现(已完成)
func (self *WebHttpForClient) RpcSearchFound(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===搜索发现 RpcSearchFound pid:%v", common.GetUserId())
	resp, err := SearchFoundService()
	return RpcReturnCommon("RpcSearchFound", resp, err)
}

// 商品展示区(首页顶部)(已完成)
func (self *WebHttpForClient) RpcProductShow(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===商品展示区 RpcProductShow pid:%v", common.GetUserId())
	resp, err := ProductShowService(common.GetUserId())
	return RpcReturnCommon("RpcProductShow", resp, err)
}

//xxx 已获得8888硬币 (已完成)
func (self *WebHttpForClient) RpcGetCoin(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===xxx 已获得8888硬币 RpcGetCoin pid:%v", common.GetUserId())
	resp, err := GetCoinService()
	return RpcReturnCommon("RpcGetCoin", resp, err)
}

// 首页消息播放区(已完成)
func (self *WebHttpForClient) RpcHomeMessage(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 首页消息播放区 RpcHomeMessage pid:%v", common.GetUserId())
	resp, err := HomeMessageService()
	return RpcReturnCommon("RpcHomeMessage", resp, err)
}

// 获取随机十条物品信息
func (self *WebHttpForClient) RpcGetRandProduct(common *base.Common, reqMsg *h5_wish.DareReq) easygo.IMessage {
	logs.Info("=== 获取随机十条物品信息 RpcGetRandProduct pid:%v, req: %v", common.GetUserId(), reqMsg)

	resp, err := GetRandProductService(reqMsg.GetBoxId())
	return RpcReturnCommon("RpcGetRandProduct", resp, err)
}

// 获取十条挑战成功记录
func (self *WebHttpForClient) RpcGetDareMessage(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 获取十条挑战成功记录 RpcGetDareMessage pid:%v", common.GetUserId())
	resp, err := GetDareMessageService()
	return RpcReturnCommon("RpcGetDareMessage", resp, err)
}

// 挑战守护者赢硬币 快捷入口的轮播展示入口 (后台配置一定会有的) 已完成
func (self *WebHttpForClient) RpcProtector(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 挑战守护者快捷入口的轮播展示入口 RpcProtector pid:%v", common.GetUserId())
	//resp, err := ProtectorService()
	resp := ProtectorService()
	return resp
}

// 最新上线,人气盲盒,欧气爆棚 菜单栏 已完成
func (self *WebHttpForClient) RpcMenu(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 最新上线,人气盲盒,欧气爆棚 菜单栏 RpcMenu pid:%v", common.GetUserId())
	resp, err := MenuService()
	return RpcReturnCommon("RpcMenu", resp, err)
}

// 综合下面的,商品品牌(苹果,雷神,三星,罗技)
func (self *WebHttpForClient) RpcProductBrand(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 综合下面的,商品品牌(苹果,雷神,三星,罗技) RpcProductBrand pid:%v", common.GetUserId())
	resp, err := ProductBrandListService()
	return RpcReturnCommon("RpcProductBrand", resp, err)
}

// 商品展示列表区 盲盒区列表(综合,最新上线,人气盲盒,欧气爆棚条件筛选) 已完成
func (self *WebHttpForClient) RpcSearchBox(common *base.Common, reqMsg *h5_wish.SearchBoxReq) easygo.IMessage {
	logs.Info("=== 盲盒区列表 RpcSearchBox pid:%v,reqMsg: %v", common.GetUserId(), reqMsg)
	resp, err := SearchBoxService(common.GetUserId(), reqMsg)
	return RpcReturnCommon("RpcSearchBox", resp, err)
}

// 物品品牌列表 (热门数据待确定,字母的完成了) {"D":[{"_id":2,"Name":"电脑","Type":"D"},{"_id":5,"Name":"电机","Type":"D"}],"E":[{"_id":3,"Name":"耳机","Type":"E"},{"_id":4,"Name":"耳环","Type":"E"}],"S":[{"_id":1,"Name":"手机","Type":"S"}]}
func (self *WebHttpForClient) RpcBrandList(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 品牌列表 RpcBrandList pid:%v ", common.GetUserId())
	resp, err := BrandListService()
	return RpcReturnCommon("RpcBrandList", resp, err)
}

// 物品类别列表 (热门数据待确定,字母的完成了) {"D":[{"_id":2,"Name":"电脑","Type":"D"},{"_id":5,"Name":"电机","Type":"D"}],"E":[{"_id":3,"Name":"耳机","Type":"E"},{"_id":4,"Name":"耳环","Type":"E"}],"S":[{"_id":1,"Name":"手机","Type":"S"}]}
func (self *WebHttpForClient) RpcProductTypeList(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 物品类别列表 RpcProductTypeList pid:%v ", common.GetUserId())
	resp, err := ProductTypeListService()
	return RpcReturnCommon("RpcProductTypeList", resp, err)
}

// 挑战赛主界面-消息播放区 已完成
func (self *WebHttpForClient) RpcDareRecommend(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 挑战赛主界面-消息播放区 RpcDareRecommend pid:%v ", common.GetUserId())
	resp, err := DareRecommendService()
	return RpcReturnCommon("RpcDareRecommend", resp, err)
}

// 排行榜列表 已完成
func (self *WebHttpForClient) RpcRankings(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 排行榜列表 RpcRankings pid:%v,reqMsg: %v ", common.GetUserId(), reqMsg)
	resp, err := RankingsService()
	return RpcReturnCommon("RpcRankings", resp, err)
}

// 我的战绩头部的挑战成功数和硬币数 (完成)
func (self *WebHttpForClient) RpcMyRecord(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 我的战绩头部的挑战成功数和硬币数 RpcMyRecord pid:%v", common.GetUserId())
	resp, err := MyRecordService(common.GetUserId())
	return RpcReturnCommon("RpcMyRecord", resp, err)
}

// 我的战绩列表 (完成)
func (self *WebHttpForClient) RpcMyDare(common *base.Common, reqMsg *h5_wish.MyDareReq) easygo.IMessage {
	logs.Info("=== 我的战绩列表 RpcMyDare pid:%v,reqMsg: %v ", common.GetUserId(), reqMsg)
	resp, err := MyDareService(common.GetUserId(), reqMsg)
	return RpcReturnCommon("RpcMyDare", resp, err)
}

// xxx 发起了挑战列表列表
func (self *WebHttpForClient) RpcDareList(common *base.Common, reqMsg *h5_wish.DareReq) easygo.IMessage {
	logs.Info("=== xxx 发起了挑战列表 RpcDareList pid:%v,ip :%v,reqMsg: %v ", common.GetUserId(), common.GetIp(), reqMsg)
	boxId := reqMsg.GetBoxId()
	if boxId == 0 {
		logs.Error("发起了挑战列表 RpcDareList boxId 为 0")
		return easygo.NewFailMsg("参数有误")
	}
	resp, err := DareListService(reqMsg)
	return RpcReturnCommon("RpcDareList", resp, err)
}

// 盲盒信息(守护者,守护时间)
func (self *WebHttpForClient) RpcBoxInfo(common *base.Common, reqMsg *h5_wish.BoxReq) easygo.IMessage {
	logs.Info("=== 盲盒信息(守护者,守护时间) RpcBoxInfo pid:%v,reqMsg: %v ", common.GetUserId(), reqMsg)
	boxId := reqMsg.GetBoxId()
	reqType := reqMsg.GetType()

	if boxId == 0 || (reqType != WISH_DARE && reqType != WISH_NO_DARE) {
		logs.Error("---发起了挑战列表 RpcDareList boxId 为 0")
		return easygo.NewFailMsg("参数有误")
	}
	resp, err := BoxInfoService(common.GetUserId(), boxId, reqType)
	return RpcReturnCommon("RpcBoxInfo", resp, err)
}

// 挑战记录,占领时长都是这个接口.
func (self *WebHttpForClient) RpcDareRecord(common *base.Common, reqMsg *h5_wish.DareRecordReq) easygo.IMessage {
	logs.Info("=== 挑战记录,占领时长都是这个接口. RpcDareRecord pid:%v,reqMsg: %v ", common.GetUserId(), reqMsg)
	boxId := reqMsg.GetBoxId()
	t := reqMsg.GetType() //查询类型;1-挑战记录;2-占领时长
	switch 0 {
	case int(boxId):
		logs.Error("--- 挑战记录,占领时长都是这个接口 boxId 为 0")
		return easygo.NewFailMsg("参数有误")
	case int(t):
		logs.Error("--- 挑战记录,占领时长都是这个接口 Type 为 0")
		return easygo.NewFailMsg("参数有误")
	}
	resp, err := DareRecordService(reqMsg)

	return RpcReturnCommon("RpcDareRecord", resp, err)
}

// 商品详情
func (self *WebHttpForClient) RpcProductDetail(common *base.Common, reqMsg *h5_wish.ProductDetailReq) easygo.IMessage {
	logs.Info("=== 商品详情 RpcProductDetail pid:%v,reqMsg: %v ", common.GetUserId(), reqMsg)
	productId := reqMsg.GetProductId()
	if productId == 0 {
		logs.Error("---商品详情 productId 为 0")
		return easygo.NewFailMsg("参数有误")
	}
	//TODO 逻辑处理
	resp, err := ProductDetailService(productId)
	return RpcReturnCommon("RpcProductDetail", resp, err)
}

// 许愿/修改愿望
func (self *WebHttpForClient) RpcWish(common *base.Common, reqMsg *h5_wish.WishReq) easygo.IMessage {
	logs.Info("=== 许愿/修改愿望 RpcWish pid:%v,reqMsg: %v ", common.GetUserId(), reqMsg)
	boxId := reqMsg.GetBoxId()
	productId := reqMsg.GetProductId()
	if boxId == 0 || productId == 0 {
		logs.Error("许愿/修改愿望 参数有误, BoxId: %d,ProductId: %d", boxId, productId)
		return easygo.NewFailMsg("参数有误")
	}
	opType := reqMsg.GetOpType() //1-许愿;2-修改愿望
	if opType != WISH_OP_TYPE_ADD && opType != WISH_OP_TYPE_EDIT {
		logs.Error("--- 许愿/修改愿望 类型有误, opType: %d", opType)
		return easygo.NewFailMsg("参数有误")
	}
	respMsg, err := WishService(reqMsg, common.GetUserId())
	return RpcReturnCommon("RpcWish", respMsg, err)
}

// 发起挑战
func (self *WebHttpForClient) RpcDoDare(common *base.Common, reqMsg *h5_wish.DoDareReq) easygo.IMessage {
	st := for_game.GetMillSecond()
	goroutineID := GetGoroutineID()

	/*
		// ==========获取分布式锁开始==========
		redisLockKey := for_game.MakeRedisKey(for_game.WISH_REDIS_LOCK_KEY, common.GetUserId())
		errLock := for_game.DoRedisLockWithRetry(redisLockKey, 20) // 20秒后key过期
		defer for_game.DoRedisUnlock(redisLockKey)
		if errLock != nil {
			logs.Error("发起挑战,获取分布式锁错误,redisKey: %s,err: %s", redisLockKey, errLock.Error())
			return easygo.NewFailMsg("系统出错")
		}
		logs.Info("goroutineID: %v,获取分布式锁成功", goroutineID)
		//  ==========获取分布式锁结束==========
	*/

	logs.Info("goroutineID: %v,=== 发起挑战 RpDoDare, pid:%v,Ip: %v,reqMsg: %v ", goroutineID, common.GetUserId(), common.GetIp(), reqMsg)

	whiteList := for_game.GetWishWhiteList() // 白名单列表
	whiteIds := make([]int64, 0)
	for _, v := range whiteList {
		whiteIds = append(whiteIds, v.GetId())
	}

	b1 := for_game.GetRedisWishPlayer(common.GetUserId())
	if b1 == nil {
		logs.Error("goroutineID: %v,获取许愿池用户对象失败,许愿池用户id为: %d", goroutineID, common.GetUserId())
		return easygo.NewFailMsg("参数有误")
	}

	// 只有运营号和普通号才能抽奖.
	if b1.GetTypes() >= 1 && !util.Int64InSlice(common.GetUserId(), whiteIds) {
		logs.Error("运营号不能参加抽奖,types: %d", b1.GetTypes)
		return easygo.NewFailMsg("此账号不允许抽奖")
	}

	// 1-挑战赛;2-非挑战赛
	dType := reqMsg.GetDareType()
	if dType != WISH_DARE && dType != WISH_NO_DARE {
		logs.Error("goroutineID: %v,--- 发起挑战, dareType: %d", goroutineID, dType)
		return easygo.NewFailMsg("参数有误")
	}
	//获取盲盒
	box, err := for_game.GetWishBox(reqMsg.GetWishBoxId(), []int64{for_game.WISH_DOWN_STATUS, for_game.WISH_PUT_ON_STATUS, for_game.WISH_ADD_PRODUCT_STATUS})
	if err != nil {
		logs.Error("goroutineID: %v,发起挑战查找盲盒失败, err: ", goroutineID, err.Error())
		return easygo.NewFailMsg("该盲盒数据有误")
	}
	if box.GetStatus() == for_game.WISH_DOWN_STATUS {
		logs.Error("goroutineID: %v,该盲盒已下架 ", goroutineID)
		return easygo.NewFailMsg("该盲盒已下架")
	}
	if box.GetStatus() == for_game.WISH_ADD_PRODUCT_STATUS {
		logs.Error("goroutineID: %v,该盲盒补货中 ", goroutineID)
		return easygo.NewFailMsg("该盲盒补货中")
	}
	if box.GetId() == 0 {
		logs.Error("goroutineID: %v,盲盒数据为空,boxId为: %d", goroutineID, reqMsg.GetWishBoxId())
		return easygo.NewFailMsg("参数有误")
	}
	// 盲盒的状态
	if box.GetStatus() == 2 { // 2 积极补货中的盲盒不能挑战
		logs.Error("goroutineID: %v,该盲盒补货中... ", goroutineID)
		return easygo.NewFailMsg("积极补货中...")
	}
	// 判断用户钻石数量是否足够
	base1 := for_game.GetRedisWishPlayer(common.GetUserId())
	if base1 == nil {
		logs.Error("goroutineID: %v,获取许愿池用户对象失败,许愿池用户id为: %d", goroutineID, common.GetUserId())
		return easygo.NewFailMsg("用户信息有误")
	}
	if base1.GetDiamond() < box.GetPrice() {
		logs.Error("goroutineID: %v,钻石不足 ", goroutineID)
		return easygo.NewFailMsg("钻石不足")
	}

	if reqMsg.GetDareType() == WISH_DARE {
		if box.GetMatch() != WISH_DARE {
			logs.Error("goroutineID: %v,用户发起的是挑战赛，而盲盒不是挑战赛盲盒,boxId为: %d", goroutineID, reqMsg.GetWishBoxId())
			return easygo.NewFailMsg("参数有误")
		}
		// 如果自己是守护者了,那就不能再挑战了
		//if box.GetGuardianId() == common.GetUserId() {
		//	logs.Error("用户现在是守护者,不可以自己挑战自己,守护者id为: %d,当前抽奖的用户id为: %d", box.GetGuardianId(), common.GetUserId())
		//	return easygo.NewFailMsg("你已成为守护者,不能发起挑战")
		//}
	}
	if len(box.GetItems()) == 0 {
		logs.Error("goroutineID: %v,盲盒中物品为空,不能挑战,boxId为: %d", goroutineID, reqMsg.GetWishBoxId())
		return easygo.NewFailMsg("参数有误")
	}

	confData := for_game.GetWishCoolDownConfigFromDB()

	//每个用户每日最多挑战抽奖30次
	if confData != nil && confData.GetIsOpen() {
		frequency := for_game.GetDareFrequency(common.GetUserId())
		if int64(frequency) >= confData.GetDayLimit() {
			logs.Error("goroutineID: %v,超过单日抽奖次数了,单日抽奖次数: %d,现在的次数: %d", goroutineID, confData.GetDayLimit(), frequency)
			return easygo.NewFailMsg("囤货太多，明日再来吧~")
		}

		createTime, b := checkIsLimit(common.GetUserId(), confData)
		if b {
			overTime := createTime + confData.GetCoolDownTime() - time.Now().Unix()
			m := make(map[string]string)
			m["msg"] = "囤货太多，先休息一下吧"
			m["createTime"] = easygo.AnytoA(overTime)
			bytes, _ := json.Marshal(m)
			logs.Error("goroutineID: %v,用户进入冷却期了,返回的结果: %s,confData: %+v", goroutineID, string(bytes), confData)
			return easygo.NewFailMsg(string(bytes))
		}
	}

	status := GetPoolStatus(box.GetWishPoolId(), goroutineID)
	if status == for_game.POOL_STATUS_BIGLOSS {
		logs.Error("goroutineID: %v,抽奖,检测到盲盒在大亏状态,返回积极补货中", goroutineID)
		return easygo.NewFailMsg("积极补货中...")
	}
	// 如果是挑战赛,如果没有许愿,不可以发起挑战
	wishData, _ := for_game.GetWishDataByStatus(common.GetUserId(), reqMsg.GetWishBoxId(), for_game.WISH_CHALLENGE_WAIT)
	if reqMsg.GetDareType() == WISH_DARE {
		if wishData.GetId() == 0 {
			logs.Error("goroutineID: %v,用户是挑战赛,没有许愿,用户id为: %d,盲盒id为: %d", goroutineID, common.GetUserId(), reqMsg.GetWishBoxId())
			return easygo.NewFailMsg("先进行许愿")
		}
	}
	// 校验盲盒中是否有需要补货的物品
	if err := CheckBoxStatus(box); err != nil {
		return err
	}
	resp, err := DoDareService1(reqMsg, common.GetUserId(), wishData, goroutineID)
	ed := for_game.GetMillSecond()
	logs.Info("goroutineID: %v,抽奖的时间为:------->%v毫秒", goroutineID, ed-st)
	if err == nil {
		easygo.Spawn(for_game.AddPayPlayerLocationLog, base1.GetPlayerId(), common.GetIp())
	}
	return RpcReturnCommon("RpDoDare", resp, err)
}

// 更多挑战
func (self *WebHttpForClient) RpcBoxList(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 更多挑战 RpcBoxList pid:%v ", common.GetUserId())
	resp, err := BoxListService()
	return RpcReturnCommon("RpcBoxList", resp, err)
}

// 挑战赛主界面守护者轮播数据（默认十条）
func (self *WebHttpForClient) RpcDefenderCarousel(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 挑战赛主界面守护者轮播数据 RpcDefenderCarousel pid:%v ", common.GetUserId())
	result := DefenderCarouselService()
	if result == nil {
		return easygo.NewFailMsg("获取挑战赛主界面守护者数据失败")
	}
	return result
}

// 获取许愿款轮播数据（默认十条）
func (self *WebHttpForClient) RpcGotWishCarousel(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 获取许愿款轮播数据 RpcGotWishCarousel pid:%v ", common.GetUserId())
	result, err := GotWishCarouselService()
	return RpcReturnCommon("RpcBackProduct", result, err)
}

// 硬币兑换砖石
func (self *WebHttpForClient) RpcCoinToDiamond(common *base.Common, reqMsg *h5_wish.CoinToDiamondReq) easygo.IMessage {
	logs.Info("=== 硬币兑换砖石 RpcCoinToDiamond pid:%v,reqMsg:= %+v ", common.GetUserId(), reqMsg)
	id := reqMsg.GetId()
	if id <= 0 {
		logs.Error("用户兑换id为空")
		return easygo.NewFailMsg("参数有误")
	}

	// 获取柠檬的playerId
	wishPlayer := for_game.GetRedisWishPlayer(common.GetUserId())
	if wishPlayer == nil {
		logs.Error("获取许愿池用户失败")
		return easygo.NewFailMsg("参数有误")
	}
	data1 := for_game.GetDiamondRechargeById(id)
	if data1 == nil {
		logs.Error("找不到对应的兑换配置记录,id为: %d", id)
		return easygo.NewFailMsg("参数有误")
	}

	diamond, coin, ChangeDiamond, err := CoinToDiamondService(common.GetUserId(), id, data1)
	if err != nil {
		return err
	}
	//兑换成功埋点
	easygo.Spawn(AddReportWishLogService, common.GetUserId(), for_game.WISH_REPORT_VEXCHANGE)
	resp := &h5_wish.CoinToDiamondResq{
		Result:       easygo.NewInt32(1),
		Coin:         easygo.NewInt64(coin),
		Diamond:      easygo.NewInt64(diamond),
		DiamondCount: easygo.NewInt64(ChangeDiamond),
	}
	return resp
}

// 获取兑换列表
func (self *WebHttpForClient) RpcDiamondRechargeList(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 获取兑换列表 RpcDiamondRechargeList pid:%v", common.GetUserId())
	result := DiamondRechargeListService(common.GetUserId())
	//访问兑换列表埋点
	easygo.Spawn(AddReportWishLogService, common.GetUserId(), for_game.WISH_REPORT_ACCESS_EXCHANGE)

	resp := &h5_wish.DiamondRechargeResp{
		DiamondRechargeList: result,
	}
	return resp
}

// 钻石消费记录
func (self *WebHttpForClient) RpcDiamondChangeLogList(common *base.Common, reqMsg *h5_wish.DiamondChangeLogReq) easygo.IMessage {
	logs.Info("=== 钻石消费记录 RpcDiamondChangeLogList pid:%v,reqMsg: %v", common.GetUserId(), reqMsg)
	t := reqMsg.GetType()
	if t != 1 && t != 2 { //  1 获取,2消耗
		logs.Error("前端传的查询类型有误,Type: %d", t)
		return easygo.NewFailMsg("参数有误")
	}

	resp := GetDiamondChangeLogByPageService(common.GetUserId(), reqMsg)
	return resp
}

// 价格区间
func (self *WebHttpForClient) RpcGetPriceSection(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 价格区间 RpcGetPriceSection pid:%v", common.GetUserId())
	section := for_game.GetPriceSection()
	resp := &h5_wish.PriceSectionResp{
		PriceSection: &h5_wish.PriceSection{
			OneMin:   easygo.NewInt32(section.GetOneMin()),
			OneMax:   easygo.NewInt32(section.GetOneMax()),
			TwoMin:   easygo.NewInt32(section.GetTwoMin()),
			TwoMax:   easygo.NewInt32(section.GetTwoMax()),
			ThreeMin: easygo.NewInt32(section.GetThreeMin()),
		},
	}
	return resp
}

// 玩法配置
func (self *WebHttpForClient) RpcPlayCfg(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("=== 玩法配置 RpcPlayCfg pid:%v", common.GetUserId())
	var dayDiamondTop, dayLimit, onceDiamondRebate int64
	if confData := for_game.GetWishCoolDownConfigFromDB(); confData != nil {
		dayLimit = confData.GetDayLimit()
	}
	if cfg := for_game.GetWishGuardianCfg(); cfg != nil {
		dayDiamondTop = cfg.GetDayDiamondTop()
		onceDiamondRebate = cfg.GetOnceDiamondRebate()
	}
	resp := &h5_wish.PlayCfgResp{
		DayDiamondTop:     easygo.NewInt64(dayDiamondTop),
		DayLimit:          easygo.NewInt64(dayLimit),
		OnceDiamondRebate: easygo.NewInt64(onceDiamondRebate),
	}
	bytes, _ := json.Marshal(resp)
	logs.Info("玩法规则配置------>", string(bytes))
	return resp
}

//批量抽奖测试
func (self *WebHttpForClient) RpcBatchDare(common *base.Common, reqMsg *h5_wish.BatchDareReq) easygo.IMessage {
	logs.Info("=== 批量抽奖测试 RpcBatchDare pid:%v,reqMsg: %v", common.GetUserId(), reqMsg)
	st := for_game.GetMillSecond()
	resp := &h5_wish.BatchDareResp{
		Result: easygo.NewInt32(1),
	}
	count := reqMsg.GetCount()
	if count == 0 {
		logs.Error("抽奖次数为空")
		return easygo.NewFailMsg("抽奖次数为空")
	}

	boxId := reqMsg.GetBoxId()
	if boxId == 0 {
		logs.Error("盲盒id为0")
		return easygo.NewFailMsg("盲盒id为0")
	}
	// 校验钻石是否充足
	info := for_game.GetWishPlayerInfo(common.GetUserId())

	// 盲盒价格
	//获取盲盒
	box, err := for_game.GetWishBox(boxId, []int64{for_game.WISH_PUT_ON_STATUS})
	if err != nil {
		logs.Error("发起挑战查找盲盒失败,err: ", err.Error())
		return easygo.NewFailMsg(err.Error())
	}
	if box.GetId() == 0 {
		logs.Error("盲盒数据为空,boxId为: %d", boxId)
		return easygo.NewFailMsg("参数有误")
	}
	totalDiamond := box.GetPrice() * int64(count)
	if totalDiamond > info.GetDiamond() {
		logs.Error("砖石不足,抽奖次数为: %d,每次钻石数为: %d, 当前人物身上的钻石为: %d", count, box.GetPrice(), info.GetDiamond())
		return easygo.NewFailMsg("钻石不足")
	}

	doDareReq := &h5_wish.DoDareReq{
		DareType:  easygo.NewInt32(2),
		WishBoxId: easygo.NewInt64(reqMsg.GetBoxId()),
	}
	for i := 0; i < int(count); i++ {
		//time.Sleep(10 * time.Millisecond)
		//easygo.Spawn(self.RpcDoDare, common, doDareReq)
		self.RpcDoDare(common, doDareReq)
	}
	ed := for_game.GetMillSecond()
	logs.Info("批量抽奖总时间为:---------->", ed-st)
	return resp

}

//许愿池活动领奖
func (self *WebHttpForClient) RpcGive(common *base.Common, reqMsg *h5_wish.GiveReq) easygo.IMessage {
	logs.Info("===========许愿池活动领奖 RpcGive=========== common=%v,reqMsg=%v", common, reqMsg)
	prizeLogId := reqMsg.GetPrizeLogId()
	if prizeLogId == 0 {
		logs.Error("许愿池活动领奖,prizeLogId 为 0")
		return easygo.NewFailMsg("参数有误")
	}
	return GiveService(common.GetUserId(), prizeLogId)
}

// 奖池列表
func (self *WebHttpForClient) RpcActPoolList(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===========奖池列表 RpcActPoolList=========== common=%v", common)
	list := ActPoolListService()
	resp := &h5_wish.ActPoolResp{
		ActPoolList: list,
	}
	return resp
}

// 奖池规则查询返回
func (self *WebHttpForClient) RpcActPoolRule(common *base.Common, reqMsg *h5_wish.ActPoolRuleReq) easygo.IMessage {
	logs.Info("===========奖池规则查询 RpcActPoolRule=========== common=%v,reqMsg=%v", common, reqMsg)
	data := ActPoolRuleService(reqMsg.GetActPoolId(), reqMsg.GetType())
	resp := &h5_wish.ActPoolRuleResp{
		WishActPoolRuleList: data,
	}
	bytes, _ := json.Marshal(resp)
	logs.Info("奖池规则查询返回------->", string(bytes))
	return resp
}

// 累计天数,次数活动查询
func (self *WebHttpForClient) RpcSumNum(common *base.Common, reqMsg *h5_wish.SumReq) easygo.IMessage {
	logs.Info("===========累计天数,次数活动查询, 1、次数 2、天数 3、周排名 4、月排名, RpcSumNum=========== common=%v,reqMsg=%v", common, reqMsg)
	resp := SumDayService(common.GetUserId(), reqMsg.GetActPoolId(), reqMsg.GetType())
	bytes, _ := json.Marshal(resp)
	logs.Info("------------->", string(bytes))
	return resp
}

// 累计金额活动查询 周榜  月榜
func (self *WebHttpForClient) RpcSumMoney(common *base.Common, reqMsg *h5_wish.SumMoneyReq) easygo.IMessage {
	logs.Info("===========累计金额活动查询周榜月榜 RpcSumMoney=========== common=%v,reqMsg=%v", common, reqMsg)
	var resp easygo.IMessage
	if reqMsg.GetDataType() == for_game.WISH_ACT_H5_WEEK_NOW || reqMsg.GetDataType() == for_game.WISH_ACT_H5_MONTH_NOW {
		resp = SumMoneyService1(common.GetUserId(), reqMsg)
	} else if reqMsg.GetDataType() == for_game.WISH_ACT_H5_WEEK || reqMsg.GetDataType() == for_game.WISH_ACT_H5_MONTH {
		resp = SumMoneyService(common.GetUserId(), reqMsg)
	}
	bytes, _ := json.Marshal(resp)
	logs.Info("RpcSumMoney------------>", string(bytes))
	return resp
}

// 活动名字
func (self *WebHttpForClient) RpcActName(common *base.Common, reqMsg *h5_wish.ActNameReq) easygo.IMessage {
	logs.Info("===========活动名字 RpcActName=========== common=%v,reqMsg=%v", common, reqMsg)
	t := reqMsg.GetType()
	if t != for_game.WISH_ACT_COUNT && t != for_game.WISH_ACT_DAY && t != for_game.WISH_ACT_WEEK_MONTH {
		logs.Error("请求类型有误,3-累计次数活动;4-累计天数活动;5-累计金额活动,Type: ", t)
		return easygo.NewFailMsg("参数有误")
	}
	title := ActNameService(t)
	resp := &h5_wish.ActNameResp{
		Name: easygo.NewString(title),
		Type: easygo.NewInt32(t),
	}
	return resp
}

// 活动开启状态.
func (self *WebHttpForClient) RpcActOpenStatus(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===========活动开启状态 RpcActOpenStatus=========== common=%v", common)
	var status bool
	isOpenCount, isOpenDay, isOpenWeekMonth := for_game.CheckWishActIsOpen()
	if isOpenCount && isOpenDay && isOpenWeekMonth {
		status = true
	}
	resp := &h5_wish.ActOpenStatusResp{
		Status: easygo.NewBool(status),
	}
	return resp
}

// 充值活动开启状态..
func (self *WebHttpForClient) RpcRechargeActStatus(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("===========充值活动开启状态 RpcRechargeActStatus=========== common=%v", common)
	var status bool

	at := for_game.GetActivityByType(6)
	t := time.Now().Unix()
	if t >= at.GetStartTime() && t <= at.GetEndTime() && at.GetStatus() == 0 {
		status = true
	}

	resp := &h5_wish.ActOpenStatusResp{
		Status: easygo.NewBool(status),
	}
	return resp
}

// todo 测试,删除
//func (self *WebHttpForClient) RpcBatchDare(common *base.Common, reqMsg *h5_wish.BatchDareReq) easygo.IMessage {
//	logs.Info("=== 批量抽奖测试 RpcBatchDare pid:%v", common.GetUserId())
//	count := reqMsg.GetCount()
//	if count == 0 {
//		logs.Error("抽奖次数为空")
//		return nil
//	}
//
//	// 检验盲盒是否在大亏状态
//	//获取盲盒
//	box, err := for_game.GetWishBox(reqMsg.GetBoxId())
//	if err != nil {
//		logs.Error("发起挑战查找盲盒失败,err: ", err.Error())
//		return easygo.NewFailMsg(err.Error())
//	}
//	if box.GetId() == 0 {
//		logs.Error("盲盒数据为空,boxId为: %d", reqMsg.GetBoxId())
//		return easygo.NewFailMsg("参数有误")
//	}
//	/*	if reqMsg.GetDareType() == WISH_DARE {
//		if box.GetMatch() != WISH_DARE {
//			logs.Error("用户发起的是挑战赛，而盲盒不是挑战赛盲盒,boxId为: %d", reqMsg.GetWishBoxId())
//			return easygo.NewFailMsg("参数有误")
//		}
//		// 如果自己是守护者了,那就不能再挑战了
//		if box.GetGuardianId() == common.GetUserId() {
//			logs.Error("用户现在是守护者,不可以自己挑战自己,守护者id为: %d,当前抽奖的用户id为: %d", box.GetGuardianId(), common.GetUserId())
//			return easygo.NewFailMsg("你已成为守护者,不能发起挑战")
//		}
//	}*/
//	if len(box.GetItems()) == 0 {
//		logs.Error("盲盒中物品为空,不能挑战,boxId为: %d", reqMsg.GetBoxId())
//		return easygo.NewFailMsg("参数有误")
//	}
//
//	confData := for_game.GetWishCoolDownConfigFromDB()
//
//	//每个用户每日最多挑战抽奖30次
//	if confData != nil && confData.GetIsOpen() {
//		frequency := for_game.GetDareFrequency(common.GetUserId())
//		if int64(frequency) >= confData.GetDayLimit() {
//			return easygo.NewFailMsg("囤货太多，明日再来吧~")
//		}
//	}
//
//	createTime, b := checkIsLimit(common.GetUserId(), confData)
//	if b {
//		m := make(map[string]string)
//		m["msg"] = "囤货太多，先休息一下吧"
//		m["createTime"] = easygo.AnytoA(createTime + confData.GetCoolDownTime() - time.Now().Unix())
//		bytes, _ := json.Marshal(m)
//		return easygo.NewFailMsg(string(bytes))
//	}
//
//	status := GetPoolStatus(box.GetWishPoolId())
//	if status == for_game.POOL_STATUS_BIGLOSS {
//		logs.Error("抽奖,检测到盲盒在大亏状态,返回积极补货中")
//		return easygo.NewFailMsg("积极补货中...")
//	}
//	// 如果是挑战赛,如果没有许愿,不可以发起挑战
//	wishData, _ := for_game.GetWishDataByStatus(common.GetUserId(), reqMsg.GetBoxId(), for_game.WISH_CHALLENGE_WAIT)
//	/*	if reqMsg.GetDareType() == WISH_DARE {
//		if wishData.GetId() == 0 {
//			logs.Error("用户是挑战赛,没有许愿,用户id为: %d,盲盒id为: %d", common.GetUserId(), reqMsg.GetWishBoxId())
//			return easygo.NewFailMsg("先进行许愿")
//		}
//	}*/
//	// 校验盲盒中是否有需要补货的物品
//	if err := CheckBoxStatus(box); err != nil {
//		return err
//	}
//	doDareReq := &h5_wish.DoDareReq{
//		DareType:  easygo.NewInt32(2),
//		WishBoxId: easygo.NewInt64(reqMsg.GetBoxId()),
//	}
//
//	wishItemChan := make(chan *share_message.PlayerWishItem)
//	countChan := make(chan int)
//
//	for i := 0; i < int(count); i++ {
//		easygo.Spawn(BatchDoDareService, doDareReq, common.GetUserId(), wishData, wishItemChan, countChan)
//	}
//	wishItemList := make([]*share_message.PlayerWishItem, 0)
//	wishItemListChan := make(chan []*share_message.PlayerWishItem)
//	for {
//		select {
//		case <-countChan:
//			logs.Info("========countChan===============")
//			count--
//		case data := <-wishItemChan:
//			logs.Info("========wishItemChan===============")
//			wishItemList = append(wishItemList, data)
//			if len(wishItemList) == int(reqMsg.GetCount()) {
//				wishItemListChan <- wishItemList
//				break
//			}
//		}
//	}
//	wc := <-wishItemListChan
//	bytes, _ := json.Marshal(wc)
//	fmt.Println("-------------->", string(bytes))
//	return nil
//}
