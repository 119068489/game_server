package hall

import (
	"encoding/base64"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/h5_wish"
	"game_server/pb/server_server"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"net/url"
)

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
	upRpc := easygo.CombineRpcMap(client_hall.UpRpc, h5_wish.UpRpc)
	self.WebHttpServer.Init(port, services, upRpc)
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

	// 登陆接口不校验
	if request.GetMethodName() != "RpcLogin" {
		// 验证token
		pBase := for_game.GetRedisPlayerBase(common.GetUserId())
		if pBase == nil {
			logs.Error("中间件获取用户信息失败")
			reply := base64.StdEncoding.EncodeToString(for_game.PacketProtoMsg(1, easygo.NewFailMsg("无效的用户")))
			_, _ = c.Writer.WriteString(url.QueryEscape(reply))
			return
		}
		if pBase.GetToken() != common.GetToken() {
			logs.Error("中间件IM-token校验失败,IM 数据的token为: %s,前端传递的token为: %s", pBase.GetToken(), common.GetToken())
			reply := base64.StdEncoding.EncodeToString(for_game.PacketProtoMsg(1, easygo.NewFailMsg("鉴权失败")))
			_, _ = c.Writer.WriteString(url.QueryEscape(reply))
			return
		}
	}

	result := self.WebHttpServer.DealRequest(0, request, common)
	reply := base64.StdEncoding.EncodeToString(for_game.PacketProtoMsg(1, result))
	_, _ = c.Writer.WriteString(url.QueryEscape(reply))
}

//TODO 消息接收分发
func (self *WebHttpForClient) RpcLogin(common *base.Common, reqMsg *client_hall.LoginMsg) easygo.IMessage {
	logs.Info("收到http请求:", common, reqMsg)
	return reqMsg
}

//==============================充值活动==============================

//获取用户充值活动数据
func (self *WebHttpForClient) RpcGetRechargeAct(common *base.Common, reqMsg *base.Empty) easygo.IMessage {
	playerId := common.GetUserId()
	logs.Info("=====获取用户充值活动数据 RpcGetRechargeAct pid:%v", playerId)

	activity := for_game.GetActivityByType(for_game.WISH_ACT_WEEK_RECHARGE)
	if activity == nil || activity.GetStatus() == 1 {
		return easygo.NewFailMsg("活动未开启")
	}
	actST := activity.GetStartTime()
	actET := activity.GetEndTime()
	curT := easygo.NowTimestamp()
	if activity.GetStatus() == 1 || actST > curT || actET < curT {
		return easygo.NewFailMsg("不在活动期限")
	}

	player := for_game.GetRedisPlayerBase(playerId)
	if player == nil {
		return easygo.NewFailMsg("用户不存在")
	}

	rest := &h5_wish.RechargeActResp{}
	//获取用户已充值的道具id
	recharged := make(map[int64]bool, 0) //已充值id，标志
	for _, v := range for_game.GetPlayerRechargeActFirst(playerId).GetLevels() {
		recharged[v] = true
	}
	//获取充值活动数据
	actCfg := for_game.GetWishCoinRechargeActCfg()
	actCfgData := make([]*h5_wish.WishCoinRechargeActivityCfg, 0)
	for _, v := range actCfg {
		cfg := &h5_wish.WishCoinRechargeActivityCfg{
			Id:      v.Id,
			Amount:  v.Amount,
			Ratio:   v.FirstRatio,
			Diamond: easygo.NewInt64(0), //为0则不赠送该币种
			EsCoin:  easygo.NewInt64(0),
		}
		if _, ok := recharged[v.GetId()]; ok {
			cfg.IsRecharge = easygo.NewBool(true)
			if v.GetIsDailyDiamond() {
				cfg.Diamond = v.DailyDiamond
			}
			if v.GetIsDailyEsCoin() {
				cfg.EsCoin = v.DailyEsCoin
			}
			cfg.Ratio = v.DailyRatio
		} else {
			cfg.IsRecharge = easygo.NewBool(false)
			if v.GetIsFirstDiamond() {
				cfg.Diamond = v.FirstDiamond
			}
			if v.GetIsFirstEsCoin() {
				cfg.EsCoin = v.FirstEsCoin
			}
		}
		actCfgData = append(actCfgData, cfg)
	}
	rest.RechargeLevels = actCfgData
	//获取用户拥有的币
	wishSrv := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_WISH)
	msg := &server_server.PlayerSI{
		PlayerId: easygo.NewInt64(playerId),
	}
	diamondResp, err := SendMsgToServerNew(wishSrv.GetSid(), "GetPlayerDiamond", msg)
	if err != nil {
		logs.Error("获取许愿池钻石失败", playerId)
	}
	rst, ok := diamondResp.(*server_server.PlayerSI)
	if ok && nil != rst {
		rest.PlayerDiamond = easygo.NewInt64(rst.GetCount())
	}
	rest.PlayerCoin = easygo.NewInt64(player.GetAllCoin())
	rest.PlayerEsCoin = easygo.NewInt64(player.GetESportCoin())
	return rest
}

//获取用户充值活动获取记录
func (self *WebHttpForClient) RpcGetRechargeLogs(common *base.Common, reqMsg *h5_wish.DataPageReq) easygo.IMessage {
	logs.Info("=====获取用户充值活动获取记录 RpcGetRechargeLogs pid:%v, reqMsg: %v", common.GetUserId(), reqMsg)
	logData := for_game.GetPlayerRechargeActLog(common.GetUserId(), reqMsg)
	if logData == nil {
		return easygo.NewFailMsg("获取用户首充记录失败")
	}

	rechargeLogs := make([]*h5_wish.RechargeLogs, 0)
	for _, v := range logData {
		rechargeLogs = append(rechargeLogs, &h5_wish.RechargeLogs{
			PayMoney:   easygo.NewInt64(v.GetMoney()),
			CreateTime: easygo.NewInt64(v.GetCreateTime()),
			CoinNum:    easygo.NewInt64(v.GetCoin()),
			GiveType:   easygo.NewInt32(v.GetGiveType()),
			GiveNum:    easygo.NewInt64(v.GetGiveCoin()),
		})
	}
	return &h5_wish.RechargeLogsResp{
		RechargeLogs: rechargeLogs,
	}
}

//==============================充值活动==============================
