package hall

import (
	"encoding/base64"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/client_server"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"strconv"
	"strings"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"

	"github.com/astaxie/beego/logs"
)

//=======================================
const PLAYER_ONLINE_NUM = "player_online_num"

type Player struct {
	*for_game.RedisPlayerBaseObj
	Mutex easygo.RLock
}

func NewPlayer(playerId PLAYER_ID) *Player {
	p := &Player{}
	p.Init(playerId)
	return p
}

func (self *Player) Init(playerId PLAYER_ID) {
	self.RedisPlayerBaseObj = for_game.GetRedisPlayerBase(playerId)
}

func (self *Player) OnLoadFromDB() {
	self.SetIsOnLine(true)
	//self.SetEndPoind(self)
	self.UpdateLogInTimestamp()
	self.UpdateLoginTimes()
	// accountInfo := for_game.GetRedisPlayerAccount(self.GetPhone())
	// if accountInfo.GetPayPassword() != "" {
	// self.IsPayPassword = true
	// }
	// if accountInfo.GetPassword() != "" {
	// 	self.IsLoginPassword = true
	// }

}
func (self *Player) GetClientEndpoint() for_game.IClientEndpoint { // override
	return ClientEpMp.LoadEndpoint(self.GetPlayerId())
}

// func (self *Player) GetIsPayPassword() bool {
// 	self.Mutex.Lock()
// 	defer self.Mutex.Unlock()
// 	return self.IsPayPassword
// }

// func (self *Player) SetIsPayPassword(b bool) {
// 	self.Mutex.Lock()
// 	defer self.Mutex.Unlock()
// 	self.IsPayPassword = b
// }

//func (self *Player) SetSafePassword(password string) {
//	self.Mutex.Lock()
//	defer self.Mutex.Unlock()
//	pwd := for_game.Md5(password)
//	self.SetSafePassword(pwd)
//	if !self.GetIsSafePassword() {
//		self.SetIsSafePassword(true)
//	}
//}

func (self *Player) DelSafePassword() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.SetSafePassword("")
}

func (self *Player) GetBankMsg(id string) *share_message.BankInfo {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for _, info := range self.GetBankInfos() {
		if id == info.GetBankId() {
			return info
		}
	}
	return nil
}

func (self *Player) SetPeopleAuth(id, name string) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if id != "" && name != "" {
		self.SetPeopleId(id)
		self.SetRealName(name)
		self.SetAuthTime(for_game.GetMillSecond())
	}
}

func (self *Player) CheckPeopleAuth() bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if self.GetPeopleId() != "" && self.GetRealName() != "" {
		return true
	}
	return false
}

func (self *Player) GetBankId() []string {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var bankIds []string
	for _, bank := range self.GetBankInfos() {
		bankIds = append(bankIds, bank.GetBankId())
	}
	return bankIds
}

func (self *Player) ChangePlayerInfo(msg *client_server.ChangePlayerInfo) string { //修改个人信息
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	t := msg.GetType()
	switch t {
	case 1:
		nikeName := msg.GetValue1()
		evilType, _ := for_game.PDirtyWordsMgr.CheckWord(nikeName)
		if evilType {
			return "名字包含敏感字，请重新编辑"
		}
		self.SetNickName(nikeName)
		return ""
	case 2:
		headIcon := msg.GetValue1()
		//头像检测
		evilType := ImageModeration(headIcon, 0, 0)
		if evilType != 100 { //敏感图
			return "图片包含敏感信息，请重新上传"
		}
		self.SetHeadIcon(headIcon)
		return ""
	case 3:
		headIcon := msg.GetValue1()
		if headIcon != "" {
			//头像检测
			evilType := ImageModeration(headIcon, 0, 0)
			if evilType != 100 { //敏感图
				return "图片包含敏感信息，请重新上传"
			}
			self.SetHeadIcon(headIcon)
		}
		sex := msg.GetValue()
		self.SetSex(sex)
		return ""
	case 4:
		photo := msg.GetPhoto()
		//检测相册图片
		for _, v := range photo {
			evilType := ImageModeration(v, 0, 0)
			if evilType != 100 { //敏感图
				return "图片包含敏感信息，请重新上传"
			}
		}
		self.SetPhoto(photo)
		return ""
	case 5:
		email := msg.GetValue1()
		self.SetEmail(email)
		return ""
	case 6:
		signture := msg.GetValue1()
		evilType, _ := for_game.PDirtyWordsMgr.CheckWord(signture)
		if evilType {
			return "签名包含敏感字，请重新编辑"
		}
		self.SetSignature(signture)
		return ""
	case 7:
		phone := msg.GetValue1()
		account := for_game.GetRedisAccountByPhone(phone)
		if account != nil {
			return "该手机号码已经被绑定"
		}
		period := for_game.GetPlayerPeriod(self.Id)
		if period.HaltYearPeriod.Fetch(for_game.CHANGE_PHONE) != nil {
			return "需要间隔半年时间才能重新绑定手机号码"
		}
		self.SetPhone(phone)
		return ""
	case 8:
		provice := msg.GetValue1()
		self.SetProvice(provice)
		return ""
	case 9:
		city := msg.GetValue1()
		self.SetCity(city)
		return ""
	case 10:
		setting := msg.GetPlayerSetting()
		if setting.IsNewMessage != nil {
			value := setting.GetIsNewMessage()
			self.SetIsNewMessage(value)
			return ""
		}
		if setting.IsMusic != nil {
			value := setting.GetIsMusic()
			self.SetIsMusic(value)
			return ""
		}
		if setting.IsShake != nil {
			value := setting.GetIsShake()
			self.SetIsShake(value)
			return ""
		}
		if setting.IsAddFriend != nil {
			value := setting.GetIsAddFriend()
			self.SetIsAddFriend(value)
			return ""
		}
		if setting.IsPhone != nil {
			value := setting.GetIsPhone()
			self.SetIsPhone(value)
			return ""
		}
		if setting.IsAccount != nil {
			value := setting.GetIsAccount()
			self.SetIsAccount(value)
			return ""
		}
		if setting.IsTeamChat != nil {
			value := setting.GetIsTeamChat()
			self.SetIsTeamChat(value)
			return ""
		}
		if setting.IsCode != nil {
			value := setting.GetIsCode()
			self.SetIsCode(value)
			return ""
		}
		if setting.IsCard != nil {
			value := setting.GetIsCard()
			self.SetIsCard(value)
			return ""
		}
		if setting.IsSafeProtect != nil {
			value := setting.GetIsSafeProtect()
			self.SetIsSafeProtect(value)
			return ""
		}
		if setting.IsMessageShow != nil {
			value := setting.GetIsMessageShow()

			self.SetIsMessageShow(value)
			return ""
		}

		if setting.IsOpenSquare != nil {
			value := setting.GetIsOpenSquare()
			self.SetIsOpenSquare(value)
			return ""
		}

		if setting.IsOpenZanOrComment != nil {
			value := setting.GetIsOpenZanOrComment()
			self.SetIsOpenZanOrComment(value)
			return ""
		}

		if setting.IsOpenRecoverComment != nil {
			value := setting.GetIsOpenRecoverComment()
			self.SetIsOpenRecoverComment(value)
			return ""
		}

		if setting.IsOpenMyAttention != nil {
			value := setting.GetIsOpenMyAttention()
			self.SetIsOpenMyAttention(value)
			return ""
		}
		if setting.IsOpenRecommend != nil {
			value := setting.GetIsOpenRecommend()
			self.SetIsOpenRecommend(value)
			return ""
		}
		if setting.IsOpenCoinShop != nil {
			value := setting.GetIsOpenCoinShop()
			self.SetIsOpenCoinShop(value)
			return ""
		}
		if setting.IsBanSayHi != nil {
			value := setting.GetIsBanSayHi()
			self.SetIsBanSayHi(value)
			return ""
		}

	case 11:
		area := msg.GetValue1()
		self.SetArea(area)
		return ""
	case 12:
		bgURL := msg.GetBackgroundImageURL()
		evilType := ImageModeration(bgURL, 0, 0)
		if evilType != 100 { //敏感图
			return "选取的照片违规无法设置为背景"
		}
		self.SetBackgroundImageURL(bgURL)
		return ""
	}
	return fmt.Sprintf("未知错误类型,%v", msg)
}

func (self *Player) CheckRegistrationIdOrChannel(id, channel string) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	if self.GetRegistrationId() != id {
		self.SetRegistrationId(id)
	}
	if self.GetChannel() == "" && channel != "" {
		self.SetChannel(channel)
	}
}

func (self *Player) GetPageCollectInfo(page, num, t int32) []*share_message.CollectInfo {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	start := (page - 1) * num
	end := page * num
	collects := self.GetCollectInfo()
	if t == 0 {
		if int(end) < len(collects) {
			return collects[start:end]
		} else {
			return collects[start:]
		}
	} else {
		var lst []*share_message.CollectInfo
		for _, msg := range collects {
			info := msg.GetCollect()
			for _, m := range info {
				if t == m.GetType() {
					lst = append(lst, msg)
					break
				}
			}

		}
		if int(end) < len(lst) {
			return lst[start:end]
		} else {
			return lst[start:]
		}
	}
	return nil
}

func (self *Player) GetAllCollectIndexList() []int32 {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var lst []int32
	for _, info := range self.GetCollectInfo() {
		lst = append(lst, info.GetIndex())
	}
	return lst
}

func (self *Player) GetMaxCollectIndex() int32 {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	collects := self.GetCollectInfo()
	if len(collects) == 0 {
		return 0
	}
	return collects[len(collects)-1].GetIndex()
}

func (self *Player) GetCollectInfoForIndex(lst []int32) *client_hall.AllCollectInfo {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var Info []*share_message.CollectInfo
	for _, index := range lst {
		for _, info := range self.GetCollectInfo() {
			if info.GetIndex() == index {
				Info = append(Info, info)
				break
			}
		}
	}
	msg := &client_hall.AllCollectInfo{
		CollectInfo: Info,
	}
	return msg
}

func (self *Player) GetSearchCollectInfo(content string) []*share_message.CollectInfo {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var lst []*share_message.CollectInfo
	for _, msg := range self.GetCollectInfo() {
		info := msg.GetCollect()
		for _, m := range info {
			if strings.Contains(m.GetContent(), content) {
				lst = append(lst, msg)
				break
			}
		}
	}
	return lst
}

func (self *Player) GetBlackInfo() *client_server.AllPlayerInfo {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var alllst []*client_server.PlayerMsg
	// 获取base
	base := for_game.GetFriendBase(self.GetPlayerId())
	for _, pid := range self.GetBlackList() {
		reName := base.GetFriendsReName(pid)
		msg := GetFriendInfo(pid, reName)
		alllst = append(alllst, msg)
	}
	msg := &client_server.AllPlayerInfo{
		PlayerMsg: alllst,
	}
	return msg
}

//从许愿池服务器获取钻石数
func GetDiamondFromWishServer(playerId int64) int64 {
	srv := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_WISH)
	if srv == nil {
		logs.Error("许愿池服务器尚未开启")
		return 0
	}
	msg := &server_server.PlayerSI{
		PlayerId: easygo.NewInt64(playerId),
	}
	resp, err := SendMsgToServerNew(srv.GetSid(), "GetPlayerDiamond", msg)
	if err != nil {
		logs.Error("去许愿池获取钻石出错")
		return 0
	}
	if re, ok := resp.(*server_server.PlayerSI); ok {
		return re.GetCount()
	}
	return 0
}

func (self *Player) GetPlayerInfo() *client_server.PlayerMsg {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	b := len(self.GetPersonalityTags()) > 0
	diamond := GetDiamondFromWishServer(self.Id)
	msg := &client_server.PlayerMsg{
		PlayerId:           easygo.NewInt64(self.Id),
		Gold:               easygo.NewInt64(self.GetGold()),
		HeadIcon:           easygo.NewString(self.GetHeadIcon()),
		NickName:           easygo.NewString(self.GetNickName()),
		Sex:                easygo.NewInt32(self.GetSex()),
		Account:            easygo.NewString(self.GetAccount()),
		Phone:              easygo.NewString(self.GetPhone()),
		Email:              easygo.NewString(self.GetEmail()),
		PeopleID:           easygo.NewString(self.GetPeopleId()),
		BankInfo:           self.GetBankInfos(),
		Signature:          easygo.NewString(self.GetSignature()),
		Provice:            easygo.NewString(self.GetProvice()),
		City:               easygo.NewString(self.GetCity()),
		Area:               easygo.NewString(self.GetArea()),
		IsPayPassword:      easygo.NewBool(self.GetIsPayPassword()),
		PlayerSetting:      self.GetPlayerSetting(),
		RealName:           easygo.NewString(self.GetRealName()),
		BlackList:          self.GetBlackList(),
		ClearLocalLogTime:  easygo.NewInt64(self.GetClearLocalLogTime()),
		IsVisitor:          easygo.NewBool(self.GetIsVisitor()),
		IsLoginPassword:    easygo.NewBool(self.GetIsLoginPassword()),
		FreeTimes:          easygo.NewInt32(self.GetFreeTimes()),
		IsBindWechat:       easygo.NewBool(self.GetIsBindWechat()),
		Emoticons:          self.GetEmoticons(),
		LabelInfo:          for_game.GetLabelInfo(self.GetLabelList()),
		AreaCode:           easygo.NewString(self.GetAreaCode()),
		BackgroundImageURL: easygo.NewString(self.GetBackgroundImageURL()),
		Coin:               easygo.NewInt64(self.GetCoin()),
		BCoin:              easygo.NewInt64(self.GetBCoin()),
		YoungPassWord:      easygo.NewString(self.GetYoungPassWord()),
		Types:              easygo.NewInt32(self.GetTypes()),
		IsCanRoam:          easygo.NewBool(self.GetIsCanRoam()),
		Constellation:      easygo.NewInt32(self.GetConstellation()),
		MixId:              easygo.NewInt64(self.GetMixId()),
		ESportCoin:         easygo.NewInt64(self.GetESportCoin()),
		IsSetPersonalTags:  easygo.NewBool(b),
		Diamond:            easygo.NewInt64(diamond),
	}
	return msg
}

func (self *Player) GetFriendsInfo(createTime ...int64) []*client_server.PlayerMsg {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var allMsg []*client_server.PlayerMsg
	base := for_game.GetFriendBase(self.Id)
	if base == nil {
		return nil
	}
	time := append(createTime, 0)[0]
	pMap := for_game.GetAllPlayerBase(base.GetFriendIds(), false)
	//获取好友信息优化
	for _, friend := range base.GetFriends() {
		if friend.GetCreateTime() < time {
			//不是新增的好友不需要查找
			continue
		}
		player, ok := pMap[friend.GetPlayerId()]
		if !ok {
			logs.Info("找不到群成员信息")
			continue
		}
		//if easygo.Contain(self.GetBlackList(), player.GetPlayerId()) {
		//	logs.Info("黑名单玩家")
		//	continue
		//}
		msg := GetFriendInfoEx(player, friend.GetReName())
		msg.FriendSetting = friend.GetSetting()
		t := share_message.AddFriend_Type(friend.GetType())
		msg.AddType = &t
		msg.BackgroundImageURL = easygo.NewString(player.GetBackgroundImageURL())
		allMsg = append(allMsg, msg)
	}
	return allMsg
}

func (self *Player) GetAllPlayerInfo(login_type int32) *client_server.AllPlayerMsg {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	msg := &client_server.AllPlayerMsg{
		Myself:            self.GetPlayerInfo(),
		AssistantInfoList: BuildUnreadAssistantList(self, login_type),
	}
	//fq := for_game.GetFriendBase(self.Id)
	//if fq != nil {
	//	msg.AllAddPlayerMsg = fq.GetNewVersionAllFriendRequestForOne()
	//}
	//msg.TweetsListResponse = self.GetSweets(1, self.GetLastLogOutTime(), util.GetMilliTime())
	//logs.Info("上次下线时间：", self.GetLastLogOutTime())
	return msg
}

func (self *Player) CheckPayPassWord(ps string) bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	return self.GetPayPassword() == ps
}

////设置连接商场服务器id
//func (self *Player) SetShopServerId(id SERVER_ID) {
//	self.Mutex.Lock()
//	defer self.Mutex.Unlock()
//	self.ShopServerId = id
//}
//
////获取连接服务器商场id
//func (self *Player) GetShopServerId() SERVER_ID {
//	self.Mutex.Lock()
//	defer self.Mutex.Unlock()
//	return self.ShopServerId
//}
func (self *Player) GetEndpoint() IGameClientEndpoint {
	return ClientEpMp.LoadEndpoint(self.GetPlayerId())
}

//发送消息给商场
func (self *Player) SendMsgToShop(methodName string, msg easygo.IMessage) easygo.IMessage {

	srv := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_SHOP)
	if srv == nil {
		logs.Error("no shop server")
		return nil
	}
	backMsg, err := SendMsgToServerNew(srv.GetSid(), methodName, msg)
	if err != nil {
		return err
	}
	return backMsg
}

// 玩家所在子游戏服务器的 endpoint
// func (self *Player) GetSubGameEndpoint() *SubGameEndpoint {
// 	return SubGameEpMp.LoadEndpoint(self.SubGameServerId)
// }

func GetFriendInfo(pid PLAYER_ID, NameList ...string) *client_server.PlayerMsg {
	name := append(NameList, "")[0]
	base := for_game.GetRedisPlayerBase(pid)
	phone := base64.StdEncoding.EncodeToString([]byte(base.GetPhone()))
	msg := &client_server.PlayerMsg{
		PlayerId:  easygo.NewInt64(base.GetPlayerId()),
		HeadIcon:  easygo.NewString(base.GetHeadIcon()),
		NickName:  easygo.NewString(base.GetNickName()),
		Sex:       easygo.NewInt32(base.GetSex()),
		Account:   easygo.NewString(base.GetAccount()),
		Photo:     base.GetPhoto(),
		Email:     easygo.NewString(base.GetEmail()),
		Signature: easygo.NewString(base.GetSignature()),
		Provice:   easygo.NewString(base.GetProvice()),
		City:      easygo.NewString(base.GetCity()),
		ReName:    easygo.NewString(name),
		Phone:     easygo.NewString(phone),
		Types:     easygo.NewInt32(base.GetTypes()),
	}
	return msg
}
func GetFriendInfoEx(player *share_message.PlayerBase, NameList ...string) *client_server.PlayerMsg {
	name := append(NameList, "")[0]
	//base := for_game.GetRedisPlayerBase(pid)

	base := for_game.GetRedisPlayerBase(player.GetPlayerId(), player)
	phone := base64.StdEncoding.EncodeToString([]byte(base.GetPhone()))
	msg := &client_server.PlayerMsg{
		PlayerId:  easygo.NewInt64(base.GetPlayerId()),
		HeadIcon:  easygo.NewString(base.GetHeadIcon()),
		NickName:  easygo.NewString(base.GetNickName()),
		Sex:       easygo.NewInt32(base.GetSex()),
		Account:   easygo.NewString(base.GetAccount()),
		Photo:     base.GetPhoto(),
		Email:     easygo.NewString(base.GetEmail()),
		Signature: easygo.NewString(base.GetSignature()),
		Provice:   easygo.NewString(base.GetProvice()),
		City:      easygo.NewString(base.GetCity()),
		ReName:    easygo.NewString(name),
		Phone:     easygo.NewString(phone),
		Types:     easygo.NewInt32(base.GetTypes()),
	}
	return msg
}

// 对外接口，获取玩家数据
func GetPlayerObj(pid PLAYER_ID) *Player {
	player := PlayerMgr.LoadPlayer(pid)
	if player != nil {
		return player
	}
	obj := NewPlayer(pid)
	//obj.OnLoadFromDB()
	return obj
}

//func ReadPersonalMessage(pid, otherId PLAYER_ID, logIds []int64) { //发送者默认已读 所以不计字段
//	key := GetChatLogKey(pid, otherId)
//	tableName := for_game.GetMongoTableName(key, for_game.TABLE_PERSONAL_CHAT_LOG)
//	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_PERSON_LOG, tableName)
//	defer closeFun()
//	for _, logId := range logIds {
//		err1 := col.Update(bson.M{"_id": logId}, bson.M{"$set": bson.M{"IsRead": true}}) //未读消息
//		easygo.PanicError(err1)
//	}
//}
//
//func WithdrawMessage(pid, otherId PLAYER_ID, logId int64) (bool, int64) { //撤回消息
//	key := GetChatLogKey(pid, otherId)
//	tableName := for_game.GetMongoTableName(key, for_game.TABLE_PERSONAL_CHAT_LOG)
//	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_PERSON_LOG, tableName)
//	defer closeFun()
//	var chat *share_message.PersonalChatLog
//	err := col.Find(bson.M{"_id": logId}).One(&chat)
//	easygo.PanicError(err)
//	if for_game.GetMillSecond()-chat.GetTime() > 5*60*1000 {
//		return false, 0
//	}
//	player := GetPlayerObj(pid)
//	content := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`"%s"撤回了一条消息`, player.GetNickName())))
//	err1 := col.Update(bson.M{"_id": logId}, bson.M{"$set": bson.M{"Content": content, "Type": for_game.TALK_CONTENT_WITHDRAW}}) //未读消息
//	easygo.PanicError(err1)
//	return true, chat.GetTime()
//}

//根据id批量查询文章
func QueryArticleByIds(ids []int64) (articleList []*share_message.Article) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ARTICLE)
	defer closeFun()
	var articles []*share_message.Article
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&articles)
	easygo.PanicError(err)
	return articles
}

//查询未推送消息
func (self *Player) GetSweets() *client_server.TweetsListResponse {
	articleUrl := easygo.YamlCfg.GetValueAsString("CLIENT_ARTICLE_URL") //测试服
	now := util.GetMilliTime()

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_TWEETS)
	defer closeFun()

	playerTweets := &share_message.PlayerTweets{}

	tweetsListResponse := &client_server.TweetsListResponse{}
	articleListRes := make([]*client_server.ArticleListResponse, 0)

	queryBson := bson.M{"_id": self.GetPlayerId()}
	err := col.Find(queryBson).One(playerTweets)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}

	ids := playerTweets.GetTweetsIdList() //获取未推送的图文信息

	if len(ids) > 0 {
		tweetsList := QuerySweetsByIds(ids)
		for _, tweets := range tweetsList { //遍历推文

			tweetsInfo := &client_server.ArticleListResponse{}
			tweetsInfo.TweetsId = easygo.NewInt64(tweets.GetID())
			tweetsInfo.ArticleListId = easygo.NewInt64(tweets.GetSendTime())

			articleIds := []int64{}
			articleList1 := tweets.GetArticle()
			for _, article := range articleList1 { //遍历文章
				articleIds = append(articleIds, article.GetID())
			}
			articleList2 := QueryArticleByIds(articleIds) //得到文章列表

			articleList := make([]*client_server.ArticleResponse, 0)
			for _, article := range articleList2 {

				articleAdd := articleUrl + "?id=" + strconv.FormatInt(article.GetID(), 10) + "&t=1&pid=" //测试服
				if article.GetTransArticleUrl() != "" && article.TransArticleUrl != nil {
					articleAdd = article.GetTransArticleUrl()
				}

				articleRes := &client_server.ArticleResponse{
					Id:          easygo.NewInt64(article.GetID()),
					Title:       easygo.NewString(article.GetTitle()),
					Icon:        easygo.NewString(article.GetIcon()),
					ArticleAdd:  easygo.NewString(articleAdd),
					ArticleType: easygo.NewInt32(article.GetArticleType()),
					Location:    easygo.NewInt32(article.GetLocation()),
					IsMain:      easygo.NewInt32(article.GetIsMain()),
					Profile:     easygo.NewString(article.GetProfile()),
					ObjectId:    easygo.NewInt64(article.GetObjectId()),
				}
				articleList = append(articleList, articleRes)
			}

			tweetsInfo.ArticleList = articleList
			articleListRes = append(articleListRes, tweetsInfo)
		}
	}

	tweetsListResponse.TweetsId = easygo.NewInt64(now)
	tweetsListResponse.TweetsList = articleListRes
	return tweetsListResponse
}

//根据id返回推文列表
func QuerySweetsByIds(ids []int64) []*share_message.Tweets {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_TWEETS)
	defer closeFun()
	var list []*share_message.Tweets
	queryBson := bson.M{"_id": bson.M{"$in": ids}}
	err := col.Find(queryBson).All(&list)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	return list
}

//根据id删除用户推文id
func (self *Player) DelSweetsByIds(ids []int64) easygo.IMessage {
	now := util.GetMilliTime()
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_TWEETS)
	defer closeFun()
	//logs.Info("删除ids:", ids)
	//logs.Info("用户id:", self.GetPlayerId())
	err := col.Update(bson.M{"_id": self.GetPlayerId()}, bson.M{"$pull": bson.M{"TweetsIdList": bson.M{"$in": ids}}, "$set": bson.M{"UpdateTime": now}})
	if err != nil {
		logs.Info("删除推文找不到ids:", ids)
		//easygo.PanicError(err)
	}
	return easygo.EmptyMsg
}

func GetAllPlayers(tweets *share_message.Tweets) []*share_message.PlayerBase {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	conditon := []bson.M{}
	queryBsonAnd := bson.M{}
	labels := []bson.M{}
	//判断设备类型
	//只有正常的账号玩家才推送
	conditon = append(conditon, bson.M{"Status": for_game.ACCOUNT_NORMAL})
	if tweets.GetUserType() != 0 && tweets.UserType != nil {
		conditon = append(conditon, bson.M{"DeviceType": tweets.GetUserType()})
	}
	//
	if tweets.GetAllLabel() != 0 {

		if len(tweets.GetList()) > 0 {
			labels = append(labels, bson.M{"Label": bson.M{"$elemMatch": bson.M{"$in": tweets.GetList()}}})
		}

		if len(tweets.GetCustomLabel()) > 0 {
			labels = append(labels, bson.M{"CustomTag": bson.M{"$elemMatch": bson.M{"$in": tweets.GetCustomLabel()}}})
		}

		//抓取词只有一个
		if len(tweets.GetCatchLabel()) > 0 {
			labels = append(labels, bson.M{"GrabTag": tweets.GetCatchLabel()[0]})
			logs.Info("抓取词:", tweets.GetCatchLabel()[0])
		}

	}

	if len(labels) > 0 {
		conditon = append(conditon, bson.M{"$or": labels})
	}

	if len(conditon) != 0 {
		queryBsonAnd["$and"] = conditon
	}
	var players []*share_message.PlayerBase
	err := col.Find(queryBsonAnd).Select(bson.M{"Label": 1, "GrabTag": 1, "CustomTag": 1, "Token": 1, "IsOnline": 1}).All(&players)
	easygo.PanicError(err)
	return players
}

//检测是否已经绑定手机号和实名认证，并且账号和创建超过30天
func (self *Player) CheckThirtyDays() bool {
	//创建账号30天才能注销
	createTime := for_game.GetMillSecond() - self.GetCreateTime()
	if createTime < 86400000*30 {
		return false
	}
	return true
}

//检测玩家账号状态是否正常
func (self *Player) CheckAccountStatus() bool {
	if self.GetStatus() != 0 {
		return false
	}
	return true
}

//检测玩家是否有投诉或者被投诉
func (self *Player) CheckDisputeState() bool {
	return true
}

//检测玩家充值支付未完成订单
func (self *Player) CheckTradeState() bool {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ORDER)
	defer closeFun()
	var orders []*share_message.Order
	err := col.Find(bson.M{"PlayerId": self.Id, "Status": for_game.ORDER_ST_WAITTING}).All(&orders)
	easygo.PanicError(err)
	if len(orders) > 0 {
		return false
	}
	return true
}

//检测玩家商场是否有未完成交易
func (self *Player) CheckShopState() bool {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
	defer closeFun()
	var orders []*share_message.TableShopOrder
	b1 := bson.M{"sponsor_id": self.Id}
	b2 := bson.M{"receiver_id": self.Id}
	ls := []int64{for_game.SHOP_ORDER_WAIT_PAY, for_game.SHOP_ORDER_WAIT_SEND, for_game.SHOP_ORDER_WAIT_RECEIVE}
	//err := col.Find(bson.M{"$or": []bson.M{b1, b2}, "Status": bson.M{"$in": ls}}).All(&orders)
	err := col.Find(bson.M{"$or": []bson.M{b1, b2}, "state": bson.M{"$in": ls}}).All(&orders)
	easygo.PanicError(err)
	if len(orders) > 0 {
		return false
	}
	if !self.CheckShopItemState() {
		return false
	}
	return true
}

//检测玩家余额是否为零,并且银行卡都解绑了
func (self *Player) CheckBalanceState() bool {
	//余额是否为0
	if self.GetGold() > 0 {
		logs.Info("零钱不为0")
		return false
	}
	//是否解绑所有银行卡
	if len(self.GetBankInfos()) > 0 {
		logs.Info("银行卡未解绑")
		return false
	}
	//是否存在正在收发的红包
	if !for_game.CheckPlayerRedPacket(self.Id) {
		logs.Info("正在发红包")
		return false
	}
	//是否存在正在收发的转账
	if !for_game.CheckPlayerTransferMoney(self.Id) {
		logs.Info("正在转账")
		return false
	}
	//是否存在正在充值或者提现订单
	if !self.CheckTradeState() {
		logs.Info("有充值未完成订单")
		return false
	}
	return true
}

//检测好友，和群状态
func (self *Player) CheckFriendTeamState() bool {
	ids := self.GetFriends()
	for _, id := range self.GetBlackList() {
		ids = easygo.Del(ids, id).([]int64)
	}
	if len(ids) > 0 || len(self.GetTeamIds()) > 0 {
		return false
	}
	return true
}

//检测玩家商场是否存在在架商品
func (self *Player) CheckShopItemState() bool {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	defer closeFun()
	var items []*share_message.TableShopItem
	err := col.Find(bson.M{"player_id": self.Id, "state": for_game.SHOP_ITEM_SALE}).All(&items)
	easygo.PanicError(err)
	if len(items) > 0 {
		return false
	}
	return true
}

//获取在线人数
func GetPlayerOnLineNum() int64 {
	b, err := easygo.RedisMgr.GetC().Exist(PLAYER_ONLINE_NUM)
	if err != nil {
		logs.Error("GetPlayerOnLineNum err:", err)
		return 0
	}
	if b {
		num, err := easygo.RedisMgr.GetC().Get(PLAYER_ONLINE_NUM)
		if err != nil {
			logs.Error("GetPlayerOnLineNum err:", err)
			return 0
		}
		return easygo.AtoInt64(num)
	}
	t := time.Now().Hour()
	num := 0
	randNum := for_game.RandInt(0, 300)
	if t >= 2 && t < 8 {
		num = 400 + randNum
	} else if t >= 8 && t < 10 {
		num = 1000 + randNum
	} else if t >= 10 && t < 14 {
		num = 1600 + randNum
	} else if t >= 14 && t < 18 {
		num = 1000 + randNum
	} else if t >= 18 && t < 22 {
		num = 1400 + randNum
	} else if t >= 22 && t < 24 {
		num = 2000 + randNum
	} else if t >= 0 && t < 2 {
		num = 1000 + randNum
	}
	err = easygo.RedisMgr.GetC().Set(PLAYER_ONLINE_NUM, num)
	if err != nil {
		logs.Error("GetPlayerOnLineNum err:", err)
		return 0
	}
	easygo.RedisMgr.GetC().Expire(PLAYER_ONLINE_NUM, 300)
	return int64(num)
}
