package wish

import (
	"encoding/base64"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/h5_wish"
	"net/url"

	"github.com/astaxie/beego/logs"

	"github.com/gin-gonic/gin"
)

const WISH_LOGIN_MUTEX_LOCK = "wish:login_mutex"

type WebHttpForClient struct {
	for_game.WebHttpServer
}

func NewWebHttpForClient(port int32) *WebHttpForClient {
	p := &WebHttpForClient{}
	p.Init(port)
	return p
}

func (self *WebHttpForClient) Init(port int32) {
	services := map[string]interface{}{
		SERVER_NAME: self,
	}
	self.WebHttpServer.Init(port, services, h5_wish.UpRpc)
	self.InitRoute()
}

//初始化路由
func (self *WebHttpForClient) InitRoute() {
	self.R.POST("/api", self.ApiEntry)
}

//api入口，路由分发  RpcLogin bysf
func (self *WebHttpForClient) ApiEntry(c *gin.Context) {
	data, b := c.Get("Data")
	if !b {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(1, easygo.NewFailMsg("err ApiEntry 1")))
		return
	}
	request, ok := data.(*base.Request)
	if !ok {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(1, easygo.NewFailMsg("err ApiEntry 2")))
		return
	}
	com, b := c.Get("Common")
	if !b {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(1, easygo.NewFailMsg("err ApiEntry 3")))
		return
	}
	common, ok := com.(*base.Common)
	if !ok {
		_, _ = c.Writer.Write(for_game.PacketProtoMsg(1, easygo.NewFailMsg("err ApiEntry 4")))
		return
	}
	logs.Info("commom:", common)

	// 登陆接口不校验
	if request.GetMethodName() != "RpcLogin" {
		// 验证token
		pBase := for_game.GetRedisWishPlayer(common.GetUserId())
		if pBase == nil {
			logs.Error("中间件获取用户信息失败")
			reply := base64.StdEncoding.EncodeToString(for_game.PacketProtoMsg(1, easygo.NewFailMsg("无效的用户")))
			_, _ = c.Writer.WriteString(url.QueryEscape(reply))

			return
		}
		if pBase.GetToke() != common.GetToken() {
			logs.Error("中间件许愿池token校验失败,许愿池数据的token为: %s,前端传递的token为: %s", pBase.GetToke(), common.GetToken())
			reply := base64.StdEncoding.EncodeToString(for_game.PacketProtoMsg(1, easygo.NewFailMsg("鉴权失败")))
			_, _ = c.Writer.WriteString(url.QueryEscape(reply))
			return
		}
	}

	result := self.WebHttpServer.DealRequest(0, request, common)
	//resp := base64.StdEncoding.EncodeToString(for_game.PacketProtoMsg(1, result))
	//logs.Info("响应给前端数据:", resp)
	//_, _ = c.Writer.WriteString(string(for_game.PacketProtoMsg(1, result)))

	reply := base64.StdEncoding.EncodeToString(for_game.PacketProtoMsg(1, result))
	_, _ = c.Writer.WriteString(url.QueryEscape(reply))
	//_, _ = c.Writer.Write(for_game.PacketProtoMsg(1, result))

}

//登录的分布式锁
func (self *WebHttpForClient) CheckLogin(id string) bool {
	b := easygo.RedisMgr.GetC().SIsMember(WISH_LOGIN_MUTEX_LOCK, id)
	return b
}
func (self *WebHttpForClient) LockLogin(id string) {
	err := easygo.RedisMgr.GetC().SAdd(WISH_LOGIN_MUTEX_LOCK, id)
	easygo.PanicError(err)
}
func (self *WebHttpForClient) UnLockLogin(id string) {
	err := easygo.RedisMgr.GetC().SRem(WISH_LOGIN_MUTEX_LOCK, id)
	easygo.PanicError(err)
}

//TODO 进入许愿池
func (self *WebHttpForClient) RpcLogin(common *base.Common, reqMsg *h5_wish.LoginReq) easygo.IMessage {
	logs.Info("RpcLogin:common=%+v:,reqMsg = %+v:", common, reqMsg)

	resp := &h5_wish.LoginResp{
		Result: easygo.NewInt32(2),
	}
	if reqMsg.GetChannel() == CHANNEL_NINGMANG { //来自柠檬im的用户
		logs.Info("im登录")
		req := &h5_wish.UserInfoReq{
			UserId: easygo.NewInt64(reqMsg.GetPlayerId()),
		}
		re, err := SendMsgToIdelServer(for_game.SERVER_TYPE_HALL, "RpcGetUseInfo", req, reqMsg.GetPlayerId())
		if err != nil {
			logs.Error("SendMsgToIdelServer err:", err.GetReason())
			return resp
		}
		//if res1, ok := re.(*h5_wish.UserInfoResp); ok {
		//	reqMsg.NickName = easygo.NewString(res1.GetName())
		//	reqMsg.HeadUrl = easygo.NewString(res1.GetHeadUrl())
		//	reqMsg.Account = easygo.NewString(res1.GetAccount())
		//} else {
		//	logs.Info("数据异常:", re)
		//	return resp
		//}

		res1, ok := re.(*h5_wish.UserInfoResp)
		if !ok {
			logs.Info("数据异常:", re)
			return resp
		}
		reqMsg.NickName = easygo.NewString(res1.GetName())
		reqMsg.HeadUrl = easygo.NewString(res1.GetHeadUrl())
		reqMsg.Account = easygo.NewString(res1.GetAccount())
		reqMsg.Types = easygo.NewInt32(res1.GetTypes())

	} else {
		//生成新的token
		token := easygo.RandStringRunes(16) + reqMsg.GetAccount()
		logs.Info("页面登录 token=", token)
		reqMsg.Token = easygo.NewString(token)
	}
	loginId := for_game.MakeNewString(reqMsg.GetChannel(), reqMsg.GetAccount())
	if self.CheckLogin(loginId) {
		logs.Info("玩家登录中:", loginId)
		return resp
	}
	self.LockLogin(loginId)         //订单上锁
	defer self.UnLockLogin(loginId) //处理完解锁
	player := for_game.GetWishPlayerByAccount(reqMsg.GetChannel(), reqMsg.GetAccount())
	if player == nil {
		newPlayer, err := for_game.CreatePlayerInfo(reqMsg)
		if err != nil {
			return resp
		}
		player = newPlayer
	}
	if reqMsg.GetChannel() == CHANNEL_NINGMANG { //校验token
		result, err1 := SendMsgToIdelServer(for_game.SERVER_TYPE_HALL, "RpcCheckWishToken", reqMsg)
		logs.Error("校验token:", err1)
		res, ok := result.(*h5_wish.LoginResp)
		if ok {
			//修改大厅
			for_game.UpdatePlayerInfoSid(reqMsg.GetPlayerId(), res.GetHallSid())
			wishPlayer := for_game.GetRedisWishPlayer(player.GetId())
			wishPlayer.SetToken(res.GetToken())
		}
		resp.Result = res.Result
		resp.Token = res.Token
	} else {
		resp.Result = easygo.NewInt32(1)
		resp.Token = player.Token
	}

	wishPid := player.GetId()
	// 判断账号是否冻结
	if CheckIsFreeze(wishPid) {
		logs.Error("登录许愿池失败,账号被冻结.wishPid: %d", wishPid)
		resp.Result = easygo.NewInt32(2)
		resp.Reason = easygo.NewString("您的许愿池账号已被冻结，请联系官方客服")
		return resp
	}

	resp.NotOneWish = easygo.NewBool(player.GetNotOneWish())
	resp.IsTryOne = easygo.NewBool(player.GetIsTryOne())
	resp.UserId = easygo.NewInt64(wishPid)
	resp.UserRole = player.Types
	whiteList := for_game.GetWishWhiteList() // 白名单列表
	for _, v := range whiteList {
		if wishPid == v.GetId() {
			resp.UserRole = easygo.NewInt32(7)
			break
		}
	}
	logs.Info("登陆结果:", resp)
	//访问埋点
	easygo.Spawn(AddReportWishLogService, wishPid, for_game.WISH_REPORT_ACCESS_WISH)
	return resp
}
