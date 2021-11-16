package wish

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/h5_wish"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"reflect"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/gin-gonic/gin"
)

type WebHttpForServer struct {
	for_game.WebHttpServer
	Service reflect.Value
}

func NewWebHttpForServer(port int32) *WebHttpForServer {
	p := &WebHttpForServer{}
	p.Init(port)
	return p
}

func (self *WebHttpForServer) Init(port int32) {
	services := map[string]interface{}{
		SERVER_NAME: self,
	}
	//TODO:分发消息定义
	upRpc := easygo.CombineRpcMap(client_hall.UpRpc, server_server.UpRpc)
	self.WebHttpServer.Init(port, services, upRpc)
	self.InitRoute()
}

//初始化路由
func (self *WebHttpForServer) InitRoute() {
	self.R.POST("/api", self.ApiEntry)
}

//api入口，路由分发  RpcLogin bysf
func (self *WebHttpForServer) ApiEntry(c *gin.Context) {
	data, b := c.Get("Data")
	if !b {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(1, easygo.NewFailMsg("err ApiEntry 1")))
		return
	}
	request, ok := data.(*base.Request)
	if !ok {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(request.GetRequestId(), easygo.NewFailMsg("err ApiEntry 2")))
		return
	}
	com, b := c.Get("Common")
	if !b {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(request.GetRequestId(), easygo.NewFailMsg("err ApiEntry 3")))
		return
	}
	common, ok := com.(*base.Common)
	if !ok {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(request.GetRequestId(), easygo.NewFailMsg("err ApiEntry 4")))
		return
	}
	result := self.WebHttpServer.DealRequest(0, request, common)

	_, _ = c.Writer.Write(for_game.PacketProtoMsg(request.GetRequestId(), result))

}

//TODO 消息接收分发
func (self *WebHttpForServer) RpcMsgToOtherServer(common *base.Common, reqMsg *share_message.MsgToServer) easygo.IMessage {
	logs.Info("收到其他服务器的请求:", common, reqMsg)
	methodName := reqMsg.GetRpcName()
	var method reflect.Value
	for _, service := range self.Services {
		method = service.MethodByName(methodName)
		if method.IsValid() {
			break
		}
	}
	var msg easygo.IMessage
	if reqMsg.GetMsgName() != "" {
		msg = easygo.NewMessage(reqMsg.GetMsgName())
		err := msg.Unmarshal(reqMsg.GetMsg())
		easygo.PanicError(err)
	}
	if !method.IsValid() || method.Kind() != reflect.Func {
		logs.Info("无效的rpc请求，找不到methodName:", methodName)
		return nil
	}
	args := make([]reflect.Value, 0, 3)
	args = append(args, reflect.ValueOf(common))
	if msg == nil {
		msg = easygo.EmptyMsg
	}
	args = append(args, reflect.ValueOf(msg))
	backMsg := method.Call(args) // 分发到指定的rpc
	if backMsg != nil {
		bb, ok := backMsg[0].Interface().(easygo.IMessage)
		if ok {
			return bb
		}
	}
	return nil
}

// 后台添加钻石
func (self *WebHttpForServer) RpcBackstageAddDiamond(common *base.Common, reqMsg *h5_wish.BackStageAddDiamondReq) easygo.IMessage {
	logs.Info("后台添加钻石,common: %v,reqMsg: %v", common, reqMsg)
	resp := &h5_wish.BackStageAddDiamondResp{
		Result: easygo.NewInt32(2),
	}
	if reqMsg.GetPlayerId() == 0 {
		logs.Error("传递的参数有误,用户id为空")
		return resp
	}
	if reqMsg.GetDiamond() == 0 {
		logs.Error("后台添加钻石参数有误,钻石数量为0")
		return resp
	}
	// 判断是否有玩家存在
	player := for_game.GetWishPlayerByAccount(reqMsg.GetChannel(), reqMsg.GetAccount())
	if player == nil {
		req := &h5_wish.LoginReq{
			Account:  easygo.NewString(reqMsg.GetAccount()),
			Channel:  easygo.NewInt32(reqMsg.GetChannel()),
			NickName: easygo.NewString(reqMsg.GetNickName()),
			HeadUrl:  easygo.NewString(reqMsg.GetHeadUrl()),
			PlayerId: easygo.NewInt64(reqMsg.GetPlayerId()),
			Token:    easygo.NewString(reqMsg.GetToken()),
		}
		newPlayer, err := for_game.CreatePlayerInfo(req)
		if err != nil {
			return resp
		}
		player = newPlayer
	}
	// 获取柠檬的playerId
	wishPlayer := for_game.GetRedisWishPlayer(player.GetId())
	if wishPlayer == nil {
		logs.Error("获取许愿池用户失败")
		return easygo.NewFailMsg("参数有误")
	}

	extendLog := &share_message.GoldExtendLog{
		PlayerId: easygo.NewInt64(wishPlayer.GetPlayerId()),
	}
	var reason string
	var sourceType int32
	if reqMsg.GetDiamond() > 0 {
		reason = "系统赠送"
		sourceType = for_game.DIAMOND_TYPE_BACK_GIVE
	} else {
		reason = "系统扣除"
		sourceType = for_game.DIAMOND_TYPE_BACK_RECYCLE
	}
	err, _ := wishPlayer.AddDiamond(reqMsg.GetDiamond(), reason, sourceType, extendLog)
	if err != nil {
		return err
	}
	resp.Result = easygo.NewInt32(1)
	return resp
}

// 后台更新钻石
func (self *WebHttpForServer) RpcBackstageUpdateDiamond(common *base.Common, reqMsg *h5_wish.BackStageUpdateDiamondReq) easygo.IMessage {
	logs.Info("后台更新钻石,common: %v,reqMsg: %v", common, reqMsg)
	resp := &h5_wish.BackStageAddDiamondResp{
		Result: easygo.NewInt32(2),
	}
	if reqMsg.GetPlayerId() == 0 {
		logs.Error("传递的参数有误,用户id为空")
		return resp
	}
	if reqMsg.GetDiamond() == 0 {
		logs.Error("后台添加钻石参数有误,钻石数量为0")
		return resp
	}
	// 判断是否有玩家存在
	player := for_game.GetWishPlayerByAccount(reqMsg.GetChannel(), reqMsg.GetAccount())
	if player == nil {
		req := &h5_wish.LoginReq{
			Account:  easygo.NewString(reqMsg.GetAccount()),
			Channel:  easygo.NewInt32(reqMsg.GetChannel()),
			NickName: easygo.NewString(reqMsg.GetNickName()),
			HeadUrl:  easygo.NewString(reqMsg.GetHeadUrl()),
			PlayerId: easygo.NewInt64(reqMsg.GetPlayerId()),
			Token:    easygo.NewString(reqMsg.GetToken()),
		}
		newPlayer, err := for_game.CreatePlayerInfo(req)
		if err != nil {
			return resp
		}
		player = newPlayer
	}
	// 获取柠檬的playerId
	wishPlayer := for_game.GetRedisWishPlayer(player.GetId())
	if wishPlayer == nil {
		logs.Error("获取许愿池用户失败")
		return easygo.NewFailMsg("参数有误")
	}

	extendLog := &share_message.GoldExtendLog{
		PlayerId: easygo.NewInt64(wishPlayer.GetPlayerId()),
	}
	var reason string
	var sourceType int32
	reason = reqMsg.GetReason()
	sourceType = reqMsg.GetSourceType()
	err, _ := wishPlayer.AddDiamond(reqMsg.GetDiamond(), reason, sourceType, extendLog)
	if err != nil {
		return err
	}
	resp.Result = easygo.NewInt32(1)
	return resp
}

// 后台设置成守护者
func (self *WebHttpForServer) RpcBackstageSetGuardian(common *base.Common, reqMsg *h5_wish.BackstageSetGuardianReq) easygo.IMessage {
	logs.Info("后台设置成守护者", common, reqMsg)
	resp := &h5_wish.BackstageSetGuardianResp{
		Result: easygo.NewInt32(2),
	}
	opType := reqMsg.GetOpType() // 1-设置守护者,2-取消守护者
	if opType != WISH_BACKSTAGE_SET_GUARDIAN && opType != WISH_BACKSTAGE_CANCEL_GUARDIAN {
		logs.Error("操作类型有误,传递的参数为: %d", opType)
		return resp
	}

	// 判断是否有玩家存在
	player := for_game.GetWishPlayerByAccount(reqMsg.GetChannel(), reqMsg.GetAccount())
	if player == nil {
		logs.Error("后台设置守护者查询许愿池用户信息失败,channel: %v,account: %v", reqMsg.GetChannel(), reqMsg.GetAccount())
		return resp
	}
	//if player == nil {
	//	req := &h5_wish.LoginReq{
	//		Account:  easygo.NewString(reqMsg.GetAccount()),
	//		Channel:  easygo.NewInt32(reqMsg.GetChannel()),
	//		NickName: easygo.NewString(reqMsg.GetNickName()),
	//		HeadUrl:  easygo.NewString(reqMsg.GetHeadUrl()),
	//		PlayerId: easygo.NewInt64(reqMsg.GetPlayerId()),
	//		Token:    easygo.NewString(reqMsg.GetToken()),
	//	}
	//	newPlayer, err := for_game.CreatePlayerInfo(req)
	//	if err != nil {
	//		return resp
	//	}
	//	player = newPlayer
	//}

	//box, err := for_game.GetWishBox(reqMsg.GetBoxId(), []int64{for_game.WISH_DOWN_STATUS})
	box, err := for_game.GetWishBoxJustById(reqMsg.GetBoxId())
	if err != nil {
		logs.Error("GetWishBoxJustById err: %v", err.Error())
		return resp
	}
	if box.GetGuardianId() == player.GetId() {
		logs.Error("后台,当前用户是守护者")
		resp.Result = easygo.NewInt32(1)
		return resp
	}

	if box.GetId() == 0 {
		logs.Error("获取盲盒失败")
		return resp
	}
	switch opType {
	case WISH_BACKSTAGE_SET_GUARDIAN:
		// 获取柠檬的playerId
		wishPlayer := for_game.GetRedisWishPlayer(player.GetId())
		if wishPlayer == nil {
			logs.Error("获取许愿池用户失败")
			return easygo.NewFailMsg("参数有误")
		}

		// 判断用户是否在白名单内
		whiteList := for_game.GetWishWhiteList() // 白名单列表
		whiteIds := make([]int64, 0)
		for _, v := range whiteList {
			whiteIds = append(whiteIds, v.GetId())
		}
		if !util.Int64InSlice(player.GetId(), whiteIds) {
			logs.Error("用户不在白名单内,不可设置守护者,pid: %d", player.GetId())
			return resp
		}

		gid := box.GetGuardianId()
		var beDareNickName string
		var dayDiamondTop, addDiamondNum int64
		if cfg := for_game.GetWishGuardianCfg(); cfg != nil {
			dayDiamondTop = cfg.GetDayDiamondTop()
			addDiamondNum = cfg.GetOnceDiamondRebate()
		}
		if gid > 0 {
			if addNum, err := for_game.UpOccupiedCoin(box.GetId(), box.GetGuardianId(), int32(dayDiamondTop), int32(addDiamondNum)); err == nil { //增加金币
				if guardian := for_game.GetRedisWishPlayer(gid); guardian != nil {
					extendLog := &share_message.GoldExtendLog{
						PlayerId: easygo.NewInt64(guardian.GetPlayerId()),
					}
					err, _ := guardian.AddDiamond(int64(addNum), "守护者收益", for_game.DIAMOND_TYPE_WISH_GUARDIAN_IN, extendLog)
					if err != nil {
						logs.Error("给守护者添加钻石失败,用户id为: %d,盲盒id为: %d", box.GetGuardianId(), box.GetId())
					}
					//设置当天被挑战者(守护者)获取的金币总数
					dataLog := &share_message.WishGuardianDiamondLog{
						Id:         easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_GUARDIAN_DIAMOND_LOG)),
						BoxId:      easygo.NewInt64(reqMsg.GetBoxId()),
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
			}
			if gPlayer := for_game.GetWishPlayerByPid(gid); gPlayer.GetId() > 0 {
				beDareNickName = gPlayer.GetNickName()
			}
		}
		UpWishBoxData := &share_message.WishBox{
			GuardianId:       easygo.NewInt64(player.GetId()),
			GuardianOverTime: easygo.NewInt64(time.Now().Unix() + for_game.AFTER_TIME_GUARDIAN),
			IsGuardian:       easygo.NewBool(true),
			GuardType:        easygo.NewInt32(1), // 0-用户自己挑战成为的守护者,1-后台设置的守护者
		}
		if err := for_game.UpWishBox(reqMsg.GetBoxId(), UpWishBoxData, 1); err != nil { //挑战成功，设置为守护者
			logs.Error("挑战成功,设置成守护者失败,更新的内容为:UpWishBoxData: %+v, err: %s", UpWishBoxData, err.Error())
		}

		if err := for_game.UpOccupied(reqMsg.GetBoxId(), gid); err != nil { //设置挑战者结束时间
			logs.Error("挑战成功,设置挑战者结束时间,  err: %s", err.Error())
		}

		playerOccupied := &share_message.WishOccupied{
			WishBoxId: easygo.NewInt64(reqMsg.GetBoxId()),
			NickName:  easygo.NewString(player.GetNickName()),
			HeadUrl:   easygo.NewString(player.GetHeadUrl()),
			PlayerId:  easygo.NewInt64(player.GetId()),
			Status:    easygo.NewInt32(1),
		}

		if err := for_game.AddOccupied(playerOccupied); err != nil { //挑战者时长记录
			logs.Error("挑战成功设置占领时长表失败,playerOccupied: %+v, err: %s", playerOccupied, err.Error())
		}

		dataLog := &share_message.WishGuardianDiamondLog{
			Id:         easygo.NewInt64(for_game.NextId(for_game.TABLE_WISH_GUARDIAN_DIAMOND_LOG)),
			BoxId:      easygo.NewInt64(reqMsg.GetBoxId()),
			PlayerId:   easygo.NewInt64(player.GetPlayerId()),
			HeadIcon:   easygo.NewString(player.GetHeadUrl()),
			NickName:   easygo.NewString(player.GetNickName()),
			WishNum:    easygo.NewInt64(1),
			CreateTime: easygo.NewInt64(for_game.GetMillSecond()),
		}
		if err := for_game.AddWishGuardianDiamondLog(dataLog); err != nil {
			logs.Error("后台添加守护者.添加挑战者挑战的成功次数,err: %s", err.Error())
		}

		poolObj := GetPoolObj(box.GetWishPoolId())
		pool := poolObj.GetPollInfoFromRedis()
		WishLog := &share_message.WishLog{
			WishBoxId:  easygo.NewInt64(reqMsg.GetBoxId()),
			DareId:     easygo.NewInt64(player.GetId()),
			DareName:   easygo.NewString(player.GetNickName()),
			BeDareId:   easygo.NewInt64(gid),
			BeDareName: easygo.NewString(beDareNickName),
			Result:     easygo.NewBool(true),
			//ChallengeItemId:   easygo.NewInt64(wishBoxItem.GetId()),
			DareHeadIcon: easygo.NewString(player.GetHeadUrl()),
			DarePrice:    easygo.NewInt64(box.GetPrice()),
			//WishItemId:        easygo.NewInt64(wishBoxItem.GetWishItemId()),
			DareType: easygo.NewInt32(1), // 1-挑战赛;2-非挑战赛
			//ChallengeItemName: easygo.NewString(wishItem.GetName()),
			PoolLocalStatus: easygo.NewInt64(pool.GetLocalStatus()),
			PoolIsOpenAward: easygo.NewBool(pool.GetIsOpenAward()),
			PoolId:          easygo.NewInt64(box.GetWishPoolId()),
			PoolIncomeValue: easygo.NewInt64(pool.GetIncomeValue()),
			GuardType:       easygo.NewInt32(1), //0-用户自己挑战成为的守护者,1-后台设置的守护者产生的记录.
		}
		for_game.AddWishLog(WishLog) //添加盲盒挑战记录
	case WISH_BACKSTAGE_CANCEL_GUARDIAN:
		if !box.GetIsGuardian() || box.GetGuardianId() == 0 || box.GetGuardianOverTime() <= 0 {
			return resp
		}
		// 获得守护时间
		guardianId := box.GetGuardianId()
		occupiedInfo, _ := for_game.GetOccupied(reqMsg.GetBoxId(), guardianId)
		occupiedTime := time.Now().Unix() - occupiedInfo.GetCreateTime()
		// 判断时间是否超过一天,如果是的话,不能处理了
		if occupiedInfo.GetCreateTime()+for_game.AFTER_TIME_GUARDIAN < time.Now().Unix() {
			//return errors.New("该守护者超过24小时了,挑战不需要去修改了")
			occupiedTime = for_game.AFTER_TIME_GUARDIAN
		}

		// 清空盲盒的守护者
		box.IsGuardian = easygo.NewBool(false)
		box.GuardianOverTime = easygo.NewInt64(0)
		box.GuardianId = easygo.NewInt64(0)
		if err := for_game.UpWishBox(reqMsg.GetBoxId(), box); err != nil {
			logs.Error("定时修改盲盒的守护者信息失败,err: %s", err.Error())
		}
		// 修改守护者占领的时间
		if err := for_game.UpOccupiedEx(reqMsg.GetBoxId(), guardianId, occupiedTime); err != nil {
			logs.Error("修改守护时间表信息失败,boxId: %d, err: %s", reqMsg.GetBoxId(), err.Error())
		}

	}

	resp.Result = easygo.NewInt32(1)
	return resp
}

//其他服务器获取钻石
func (self *WebHttpForServer) GetPlayerDiamond(common *base.Common, reqMsg *server_server.PlayerSI) easygo.IMessage {
	logs.Info("GetPlayerDiamond:", reqMsg)
	player := for_game.GetWishPlayerInfoByImId(reqMsg.GetPlayerId())
	count := int64(0)
	if player != nil {
		obj := for_game.GetRedisWishPlayer(player.GetId())
		if obj != nil {
			count = player.GetDiamond()
		}
	}
	reqMsg.Count = easygo.NewInt64(count)
	logs.Info("返回数据:", reqMsg)
	return reqMsg
}

//修改许愿池用户Account,处理修改IM手机号
func (self *WebHttpForServer) RpcSetWishAccount(common *base.Common, reqMsg *server_server.PlayerSI) easygo.IMessage {
	logs.Info("RpcSetWishAccount:", reqMsg)
	player := for_game.GetWishPlayerInfoByImId(reqMsg.GetPlayerId())
	if player != nil {
		obj := for_game.GetRedisWishPlayer(player.GetId())
		if obj != nil {
			obj.SetAccount(reqMsg.GetAccount())
		}
	}
	return reqMsg
}

// 后台许愿/修改愿望
func (self *WebHttpForServer) RpcWish(common *base.Common, reqMsg *h5_wish.WishReq) easygo.IMessage {
	logs.Info("后台许愿/修改愿望 RpcWish:", reqMsg)
	boxId := reqMsg.GetBoxId()
	productId := reqMsg.GetProductId()
	if boxId == 0 || productId == 0 {
		logs.Error("后台许愿/修改愿望 参数有误, BoxId: %d,ProductId: %d", boxId, productId)
		return easygo.NewFailMsg("参数有误")
	}
	opType := reqMsg.GetOpType() //1-许愿;2-修改愿望
	if opType != WISH_OP_TYPE_ADD && opType != WISH_OP_TYPE_EDIT {
		logs.Error("--- 后台许愿/修改愿望 类型有误, opType: %d", opType)
		return easygo.NewFailMsg("参数有误")
	}
	respMsg, err := WishService(reqMsg, common.GetUserId())
	return RpcReturnCommon("RpcWish", respMsg, err)
}

// 后台发起挑战
func (self *WebHttpForServer) RpcDoDare(common *base.Common, reqMsg *h5_wish.DoDareReq) easygo.IMessage {
	logs.Info("后台发起挑战 RpcDoDare:", reqMsg)
	st := for_game.GetMillSecond()
	goroutineID := GetGoroutineID()

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
		return easygo.NewFailMsg("该盲盒已下架")
	}
	if box.GetStatus() == for_game.WISH_ADD_PRODUCT_STATUS {
		return easygo.NewFailMsg("该盲盒补货中")
	}
	if box.GetId() == 0 {
		logs.Error("goroutineID: %v,盲盒数据为空,boxId为: %d", goroutineID, reqMsg.GetWishBoxId())
		return easygo.NewFailMsg("参数有误")
	}
	// 盲盒的状态
	if box.GetStatus() == 2 { // 2 积极补货中的盲盒不能挑战
		return easygo.NewFailMsg("积极补货中...")
	}
	// 判断用户钻石数量是否足够
	base1 := for_game.GetRedisWishPlayer(common.GetUserId())
	if base1 == nil {
		logs.Error("goroutineID: %v,获取许愿池用户对象失败,许愿池用户id为: %d", goroutineID, common.GetUserId())
		return easygo.NewFailMsg("用户信息有误")
	}
	if base1.GetDiamond() < box.GetPrice() {
		return easygo.NewFailMsg("钻石不足")
	}

	if reqMsg.GetDareType() == WISH_DARE {
		if box.GetMatch() != WISH_DARE {
			logs.Error("goroutineID: %v,用户发起的是挑战赛，而盲盒不是挑战赛盲盒,boxId为: %d", goroutineID, reqMsg.GetWishBoxId())
			return easygo.NewFailMsg("参数有误")
		}

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
			return easygo.NewFailMsg("囤货太多，明日再来吧~")
		}
	}

	createTime, b := checkIsLimit(common.GetUserId(), confData)
	if b {
		m := make(map[string]string)
		m["msg"] = "囤货太多，先休息一下吧"
		m["createTime"] = easygo.AnytoA(createTime + confData.GetCoolDownTime() - time.Now().Unix())
		bytes, _ := json.Marshal(m)
		return easygo.NewFailMsg(string(bytes))
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

// 后台工具抽奖
func (self *WebHttpForServer) RpcBackstageDareTool(common *base.Common, reqMsg *h5_wish.BackstageDareToolReq) easygo.IMessage {
	logs.Info("后台工具抽奖 RpcBackstageDareTool reqMsg: %v", reqMsg)
	resp := &h5_wish.BackstageDareToolResp{
		Result:    easygo.NewInt64(2),
		DareCount: easygo.NewInt64(0),
	}

	poolId := reqMsg.GetPoolId()
	if poolId == 0 {
		logs.Error("RpcBackstageDareTool poolId 为 0")
		return resp
	}
	count := reqMsg.GetCount()
	if count == 0 {
		logs.Error("RpcBackstageDareTool count 为 0")
		return resp
	}
	diamond := reqMsg.GetDiamond()
	if diamond == 0 {
		logs.Error("RpcBackstageDareTool diamond 为 0")
		return resp
	}
	userId := reqMsg.GetUserId()
	if userId == 0 {
		logs.Error("RpcBackstageDareTool userId 为 0")
		return resp
	}
	boxId := reqMsg.GetWishBoxId()
	if boxId == 0 {
		logs.Error("RpcBackstageDareTool boxId 为 0")
		return resp
	}
	// service
	dareCount := BackstageDareToolService(boxId, common.GetUserId(), poolId, count, diamond, userId)
	if dareCount != int(count) {
		logs.Error("BackstageDareToolService 没有达到抽奖次数")
	}
	resp.Result = easygo.NewInt64(1)
	resp.DareCount = easygo.NewInt64(dareCount)
	return resp
}

// 后台清理抽奖工具数据
func (self *WebHttpForServer) RpcBackstageClearToolData(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	logs.Info("后台工具抽奖 RpcBackstageClearToolData reqMsg: %v", reqMsg)

	resp := &h5_wish.BackstageClearToolDataResp{
		Result: nil,
	}
	return resp
}
