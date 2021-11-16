// 活动监听处理

package backstage

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"net/http"
	"net/url"
	"strconv"
)

//首页返回数据
type IndexData struct {
	PlayerProps   []*share_message.PlayerProps `json:"PlayerProps"`   //用户已收集的卡片
	LuckyCount    int32                        `json:"LuckyCount"`    //用户剩余抽卡次数
	PlayerTask    *for_game.LuckyPlayerTask    `json:"PlayerTask"`    //用户任务列表
	IsOpenLucky   bool                         `json:"IsOpenLucky"`   //用户是否开过奖
	LuckyMoney    int64                        `json:"LuckyMoney"`    //用户开奖金额
	FullCount     int64                        `json:"FullCount"`     //已集齐的人数
	OpenStatus    int32                        `json:"OpenStatus"`    //活动 3预热,1等待,2开奖
	FullNames     []string                     `json:"FullNames"`     //集齐玩家名字
	NowTime0Clock int64                        `json:"NowTime0Clock"` //当前时间0点时间
	FullPlaces    int64                        `json:"FullPlaces"`    //第几位集齐玩家
}

//抽卡返回数据
type DrawCard struct {
	DrawCard *share_message.PlayerProps
	IsFull   bool
}

type ActivityTime struct {
	Activity      *share_message.Activity `json:"Activity"`      //当前时间0点时间
	NowTime0Clock int64                   `json:"NowTime0Clock"` //当前时间0点时间
	FullNames     []string                `json:"FullNames"`     //集齐玩家名字
	FullCount     int64                   `json:"FullCount"`     //已集齐的人数
}

//活动api
//url?t=1&id=1&pid=1885040404&sign=md5(所传参数askill排序后的值+md5key)
func (self *WebHttpServer) ActivityEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	params := r.Form
	if !self.ApiCheckSign(params) {
		OutputJson(w, 0, "签名错误", nil)
		return
	}

	t := params.Get("t") //请求数据类型
	switch t {
	case "1": //集卡活动首页
		self.CardActivity(w, params)
	case "2": //抽卡
		self.DrawCard(w, params)
	case "3": //送卡记录
		self.SendCardsLog(w, params)
	case "4": //赠送卡
		self.SendCards(w, params)
	case "5": //发验证码
		self.SendPhoneCode(w, params)
	case "6": //激活建立绑定关系
		self.ActivateBind(w, params)
	case "7": //开奖
		self.CardOpenEggs(w, params)
	case "8": //签到
		self.CardSignIn(w, params)
	case "9": //分享活动
		self.ShareActivity(w, params)
	case "10": //好友列表
		self.FriendsList(w, params)
	case "11": //获取活动数据
		self.GetActivity(w, params)
	default:
		OutputJson(w, 0, "请求类型错误", nil)
	}
}

//活动首页
func (self *WebHttpServer) CardActivity(w http.ResponseWriter, params url.Values) {
	id := params.Get("id")
	pid := params.Get("pid")
	// logs.Debug(params)
	// 参数校验
	m := make(map[string]string)
	m["id"] = id
	if errStr := VerifyParams(m, "CardActivity"); errStr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}
	activityId, _ := strconv.ParseInt(id, 10, 64)
	act := for_game.GetActivityFromDB(activityId)
	if act == nil || act.GetStatus() == for_game.ACTIVITY_CLOSE {
		OutputJson(w, 0, "活动不存在", nil)
		return
	}

	startTime := act.GetStartTime()
	if startTime > 9999999999 {
		startTime = startTime / 1000
	}
	if startTime > easygo.NowTimestamp() {
		OutputJson(w, 0, "活动未开始", nil)
		return
	}
	endTime := act.GetCloseTime()
	if endTime > 9999999999 {
		endTime = endTime / 1000
	}
	if endTime < easygo.NowTimestamp() {
		OutputJson(w, 0, "活动已结束", nil)
		return
	}

	data := &IndexData{}
	playerId, _ := strconv.ParseInt(pid, 10, 64)
	playerProps, err := for_game.GetPlayerPropsList(playerId)
	if err != nil {
		OutputJson(w, 0, err.GetReason(), nil)
		return
	}

	data.PlayerProps = playerProps

	playerTask, err := for_game.GetLuckyPlayerTask(playerId)
	if err != nil {
		OutputJson(w, 0, err.GetReason(), nil)
		return
	}
	data.PlayerTask = playerTask
	data.FullCount = for_game.GetLocalFullCount()
	data.LuckyCount = for_game.GetLuckyCount(playerId)
	if act.GetStartTime() > easygo.NowTimestamp() {
		data.OpenStatus = 3
	}
	lp := for_game.GetLuckyPlayerFromDB(playerId)
	if act.GetOpenTime() < easygo.NowTimestamp() && easygo.NowTimestamp() < act.GetCloseTime() {
		if lp.GetIsOpen() { // 判断该玩家是否已抽过奖了
			data.IsOpenLucky = true
			data.LuckyMoney = lp.GetLuckyMoney()
		}
		data.OpenStatus = 2
	}

	waitTimeS := act.GetOpenTime() - 10*60
	if waitTimeS < easygo.NowTimestamp() && easygo.NowTimestamp() < act.GetOpenTime() {
		data.OpenStatus = 1
	}
	data.FullNames = for_game.GetFullLuckyPlayerList()
	data.NowTime0Clock = easygo.NowTimestamp()
	if lp != nil {
		data.FullPlaces = lp.GetFullPlaces()
	}

	OutputJson(w, 1, "success", data)
}

//抽卡
func (self *WebHttpServer) DrawCard(w http.ResponseWriter, params url.Values) {
	id := params.Get("id")
	pid := params.Get("pid")
	// 参数校验
	m := make(map[string]string)
	m["id"] = id
	m["pid"] = pid
	if errStr := VerifyParams(m, "DrawCard"); errStr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}
	aid, _ := strconv.ParseInt(id, 10, 64)
	playerId, _ := strconv.ParseInt(pid, 10, 64)
	playerProps, _ := for_game.GetPlayerPropsList(playerId)
	propscount := len(playerProps)
	drawCard, err := for_game.LuckyCard(playerId, aid)
	if err != nil {
		OutputJson(w, 0, err.GetReason(), nil)
		return
	}
	isFull := false
	if propscount < 6 {
		newplayerProps, _ := for_game.GetPlayerPropsList(playerId)
		if len(newplayerProps) == 6 {
			isFull = true
		}
	}

	data := &DrawCard{
		DrawCard: drawCard,
		IsFull:   isFull,
	}

	OutputJson(w, 1, "success", data)

	easygo.Spawn(func() {
		switch data.DrawCard.GetPropsId() {
		case for_game.ID_HE:
			for_game.UpdateActivityReport("CardHe", 1)
		case for_game.ID_NING:
			for_game.UpdateActivityReport("CardNing", 1)
		case for_game.ID_MENG:
			for_game.UpdateActivityReport("CardMeng", 1)
		case for_game.ID_QU:
			for_game.UpdateActivityReport("CardQu", 1)
		case for_game.ID_LV:
			for_game.UpdateActivityReport("CardLv", 1)
		case for_game.ID_XING:
			for_game.UpdateActivityReport("CardXing", 1)
		}
	})

}

//送卡
func (self *WebHttpServer) SendCards(w http.ResponseWriter, params url.Values) {
	id := params.Get("id")
	pid := params.Get("pid")
	cid := params.Get("cid")
	aid := params.Get("aid")
	// 参数校验
	m := make(map[string]string)
	m["id"] = id
	m["pid"] = pid
	m["cid"] = cid
	m["aid"] = aid
	if errStr := VerifyParams(m, "SendCards"); errStr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}
	playerId, _ := strconv.ParseInt(id, 10, 64)
	frendId, _ := strconv.ParseInt(pid, 10, 64)
	propsId, _ := strconv.ParseInt(cid, 10, 64)
	activityId, _ := strconv.ParseInt(aid, 10, 64)
	if playerId <= 0 || frendId <= 0 || propsId <= 0 || activityId <= 0 {
		OutputJson(w, 0, "参数id、pid、cid、aid不能为空", nil)
		return
	}

	err := for_game.GiveCard(activityId, playerId, frendId, propsId)
	if err != nil {
		OutputJson(w, 0, err.GetReason(), nil)
		return
	}

	easygo.Spawn(func() {
		req := &share_message.SystemNotice{
			Id:      easygo.NewInt64(frendId),
			Title:   easygo.NewString("好友送了一张柠檬卡给你"),
			Content: easygo.NewString("好友送了一张柠檬卡给你，快去集卡活动详情页看看吧~"),
		}
		SendToPlayer(frendId, "RpcActivityNotic", req)
	})

	OutputJson(w, 1, "success", nil)
}

//送卡记录
func (self *WebHttpServer) SendCardsLog(w http.ResponseWriter, params url.Values) {
	pid := params.Get("pid")
	// 参数校验
	m := make(map[string]string)
	m["pid"] = pid
	if errStr := VerifyParams(m, "SendCardsLog"); errStr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}
	playerId, _ := strconv.ParseInt(pid, 10, 64)
	data := for_game.GetGiveList(playerId)
	OutputJson(w, 1, "success", data)
}

//激活建立绑定关系
func (self *WebHttpServer) ActivateBind(w http.ResponseWriter, params url.Values) {
	sid := params.Get("sid") //邀请人
	phone := params.Get("phone")
	code := params.Get("code")
	id := params.Get("id")
	// 参数校验
	m := make(map[string]string)
	m["sid"] = sid
	m["phone"] = phone
	m["code"] = code
	m["id"] = id
	if errStr := VerifyParams(m, "ActivateBind"); errStr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}

	aid, _ := strconv.ParseInt(id, 10, 64)
	err := for_game.CheckMessageCode(phone, code, for_game.CLIENT_CODE_REGISTER)
	if err != nil {
		OutputJson(w, 0, "验证码错误", nil)
		return
	}
	playerId, _ := strconv.ParseInt(sid, 10, 64)
	errmsg := for_game.InviteNewFriend(aid, playerId, phone)
	if errmsg != nil {
		OutputJson(w, 0, errmsg.GetReason(), nil)
		return
	}

	OutputJson(w, 1, "success", nil)
}

//开奖
func (self *WebHttpServer) CardOpenEggs(w http.ResponseWriter, params url.Values) {
	pid := params.Get("pid")
	id := params.Get("id")
	// 参数校验
	m := make(map[string]string)
	m["pid"] = pid
	m["id"] = id
	if errStr := VerifyParams(m, "CardOpenEggs"); errStr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}

	playerId, _ := strconv.ParseInt(pid, 10, 64)
	aid, _ := strconv.ParseInt(id, 10, 64)

	data, err := for_game.OpenLucky(playerId, aid)
	if err != nil {
		OutputJson(w, 0, err.GetReason(), nil)
		return
	}
	OutputJson(w, 1, "success", data)

	req := &server_server.Recharge{
		PlayerId:     easygo.NewInt64(playerId),
		RechargeGold: easygo.NewInt64(data.GetLuckyMoney()),
		OrderId:      easygo.NewString("0"),
		SourceType:   easygo.NewInt32(409), //临时增加活动奖励类型
	}

	SendToPlayer(playerId, "RpcActivityAddGold", req)

	easygo.Spawn(for_game.UpdateActivityReport, "LuckPlayerCount", 1) //活动报表增加开奖人数
}

//集卡签到
func (self *WebHttpServer) CardSignIn(w http.ResponseWriter, params url.Values) {
	id := params.Get("id")
	pid := params.Get("pid")
	// 参数校验
	m := make(map[string]string)
	m["pid"] = pid
	m["id"] = id
	if errStr := VerifyParams(m, "CardSignIn"); errStr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}

	aid, _ := strconv.ParseInt(id, 10, 64)
	playerId, _ := strconv.ParseInt(pid, 10, 64)
	err := for_game.LuckySignIn(aid, playerId)
	if err != nil {
		OutputJson(w, 0, err.GetReason(), nil)
		return
	}
	easygo.Spawn(for_game.UpdateActivityReport, "TaskSignIn", 1) //活动报表增加完成签到数
	OutputJson(w, 1, "success", nil)
}

//分享活动
func (self *WebHttpServer) ShareActivity(w http.ResponseWriter, params url.Values) {
	pid := params.Get("pid")
	id := params.Get("id") //活动id
	// 参数校验
	m := make(map[string]string)
	m["pid"] = pid
	m["id"] = id
	if errStr := VerifyParams(m, "ShareActivity"); errStr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}

	playerId, _ := strconv.ParseInt(pid, 10, 64)
	aid, _ := strconv.ParseInt(id, 10, 64)
	easygo.Spawn(for_game.UpdateActivityReport, "ShareTimes", 1) //活动报表增加分享次数
	err := for_game.LuckyShare(aid, playerId)
	if err != nil {
		OutputJson(w, 0, err.GetReason(), nil)
		return
	}

	OutputJson(w, 1, "success", nil)
	easygo.Spawn(for_game.UpdateActivityReport, "TaskShare", 1) //活动报表增加分享任务完成
}

//获取好友列表
func (self *WebHttpServer) FriendsList(w http.ResponseWriter, params url.Values) {
	pid := params.Get("pid")
	// 参数校验
	m := make(map[string]string)
	m["pid"] = pid
	if errStr := VerifyParams(m, "FriendsList"); errStr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}
	playerId, _ := strconv.ParseInt(pid, 10, 64)
	if playerId <= 0 {
		OutputJson(w, 0, "用户不能为空", nil)
		return
	}
	pmgr := for_game.GetRedisPlayerBase(playerId)
	if pmgr == nil {
		OutputJson(w, 0, "用户id错误", nil)
		return
	}
	ids := pmgr.GetFriends()
	data := []*share_message.PlayerBase{}
	lis := GetPlayerBaseByIds(ids)
	for _, i := range lis {
		one := &share_message.PlayerBase{
			PlayerId: easygo.NewInt64(i.GetPlayerId()),
			NickName: easygo.NewString(i.GetNickName()),
			HeadIcon: easygo.NewString(i.GetHeadIcon()),
		}
		data = append(data, one)
	}
	OutputJson(w, 1, "success", data)
}

//获取活动数据
func (self *WebHttpServer) GetActivity(w http.ResponseWriter, params url.Values) {
	id := params.Get("id")
	// 参数校验
	m := make(map[string]string)
	m["id"] = id
	if errStr := VerifyParams(m, "GetActivity"); errStr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}
	activityId, _ := strconv.ParseInt(id, 10, 64)
	if activityId <= 0 {
		OutputJson(w, 0, "活动ID错误", nil)
		return
	}
	act := for_game.GetActivityFromDB(activityId)
	if act == nil || act.GetStatus() == for_game.ACTIVITY_CLOSE {
		OutputJson(w, 0, "活动不存在", nil)
		return
	}

	data := &ActivityTime{
		Activity:      act,
		NowTime0Clock: easygo.NowTimestamp(),
		FullNames:     for_game.GetFullLuckyPlayerList(),
		FullCount:     for_game.GetLocalFullCount(),
	}

	OutputJson(w, 1, "success", data)
}
