// 大厅服务器为[游戏客户端]提供的服务

package sport_apply

import (
	"fmt"
	dal "game_server/e-sports/sport_common_dal"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"reflect"
	"time"
)

type ServiceForClient struct {
	Service reflect.Value
}
type sfc = ServiceForClient

func GetPlayerId(common *base.Common) int64 {
	//logs.Info("common", common)
	playerId := common.GetUserId()
	return playerId
}

func (self *sfc) RpcESportEnter(common *base.Common, reqMsg *client_hall.ESportCommonResult) easygo.IMessage {
	logs.Info("===api RpcESportEnter===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存
	rd := client_hall.ESportCommonResult{
		Code: easygo.NewInt32(1),
		Msg:  easygo.NewString("跑错地方了吧，小老弟"),
	}
	return &rd
}

//获取某个首页的数据 使用 MenuId 菜单Id
func (self *sfc) RpcESportGetHomeInfo(common *base.Common, reqMsg *client_hall.ESportInfoRequest) *client_hall.ESportMenuHomeInfo {

	logs.Info("RpcESportGetHomeInfo 提交", reqMsg)

	rd := &client_hall.ESportMenuHomeInfo{}
	//playerId := GetPlayerId(common)
	//pinfo := dal.GetTableESPortsPlayerById(playerId)

	//CraeteESPortsPlayer(playerId) //加载用户

	menuId := reqMsg.GetMenuId()
	//rd.GameLabelList = dal.GetTableESPortsGameLabelList()
	rd.LabelList = dal.GetTableESPortsLabelList(menuId)
	//rd.CarouselList = dal.GetTableESPortsCarouselByMenuId(menuId, 1)
	//pinfo := GetOrCraeteESPortsPlayer(playerId) //for_game.GetRedisESportPlayerObj(playerId)
	//if pinfo != nil {
	//	pbase := for_game.GetRedisPlayerBase(playerId)
	//	if pbase != nil {
	//	c := dal.GetSysMsgCount(pinfo.GetLastPullTime(), playerId, pbase.GetDeviceType())
	//	rd.SysMsgCount = easygo.NewInt32(c)
	//	}
	//}

	//logs.Info("RpcESportGetHomeInfo 返回", rd)
	return rd
}

//获取资讯列表 使用TypeId =游戏标签ID   GameTpyeId
func (self *sfc) RpcESportGetRealtimeList(common *base.Common, reqMsg *client_hall.ESportPageRequest) *client_hall.ESportRealtimeListResult {

	logs.Info("RpcESportGetRealtimeList 提交", reqMsg)
	rd := &client_hall.ESportRealtimeListResult{}
	page := reqMsg.GetPage()
	pageSize := reqMsg.GetPageSize()
	sort := "-BeginEffectiveTime"
	if reqMsg.AscOrDesc != nil && reqMsg.GetAscOrDesc() != "" && reqMsg.OrderField != nil && reqMsg.GetOrderField() != "" {
		sort = reqMsg.GetAscOrDesc() + reqMsg.GetOrderField()
	}

	typeId := reqMsg.GetTypeId()
	labelId := reqMsg.GetLabelId()
	list, count := dal.GetESPortsRealTimeItemList(page, pageSize, sort, typeId, labelId)

	if len(list) > 0 {
		plyid := GetPlayerId(common)
		tuobj := for_game.GetRedisESportThumbsUpObj(for_game.ESPORTMENU_REALTIME, 0, plyid)

		thumbsUpList := tuobj.GetThumbsList()
		for _, v := range list {
			id := v.GetId()
			if tuobj.IsInThumbsList(thumbsUpList, id) {
				v.IsThumbsUp = easygo.NewInt32(1)
			} else {
				v.IsThumbsUp = easygo.NewInt32(0)
			}
		}
	}
	rd.List = list
	rd.Total = easygo.NewInt32(count)
	//logs.Info("RpcESportGetRealtimeList 返回", rd)
	return rd
}

//获取资讯数据
func (self *sfc) RpcESportGetRealtimeInfo(common *base.Common, reqMsg *client_hall.ESportInfoRequest) easygo.IMessage {
	logs.Info("RpcESportGetRealtimeInfo 提交", reqMsg)
	rd := &client_hall.ESportRealTimeResult{}
	id := reqMsg.GetDataId()
	if id < 1 {
		rd.Msg = easygo.NewString("资讯ID不能小于0")
		rd.Code = easygo.NewInt32(for_game.C_INFO_NOT_EXISTS)
		return rd
	}

	data := dal.GetTableESPortsRealTimeInfoById(reqMsg.GetDataId())
	if data == nil || data.GetId() < 1 {
		rd.Msg = easygo.NewString("资讯不存在")
		rd.Code = easygo.NewInt32(for_game.C_INFO_NOT_EXISTS)
		return rd
	}

	plyid := GetPlayerId(common)
	_, b := dal.IsThumbsUp(for_game.ESPORTMENU_REALTIME, 0, data.GetId(), plyid)
	if b {
		data.IsThumbsUp = easygo.NewInt32(1)
	} else {
		data.IsThumbsUp = easygo.NewInt32(0)
	}
	rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
	rd.Data = data
	//logs.Info("RpcESportGetRealtimeInfo 返回", rd)
	return rd
}

//点赞操作（所有点赞）
func (self *sfc) RpcESportThumbsUp(common *base.Common, reqMsg *client_hall.ESportInfoRequest) *client_hall.ESportThumbsUpResult {
	logs.Info("RpcESportThumbsUp 提交", reqMsg)
	plyid := GetPlayerId(common)
	rd := &client_hall.ESportThumbsUpResult{}
	menuid := reqMsg.GetMenuId() //菜单Id 资讯，还是视频
	//gameTypeId := reqMsg.GetGameTypeId()
	preId := reqMsg.GetDataId() //如果是回复 请填写 父ID
	exId := reqMsg.GetExtId()   //当前视频或者资讯或者回复的ID
	table := ""
	//menuid = for_game.ESPORTMENU_RECREATION
	if menuid == for_game.ESPORTMENU_RECREATION { //娱乐视频
		if preId > 0 {
			table = for_game.TABLE_ESPORTS_COMMENT_VIDEO
		} else {
			table = for_game.TABLE_ESPORTS_VIDEO
		}
	} else if menuid == for_game.ESPORTMENU_REALTIME { //资讯
		if preId > 0 {
			table = for_game.TABLE_ESPORTS_COMMENT_NEWS
		} else {
			table = for_game.TABLE_ESPORTS_NEWS
		}
	}
	if len(table) > 0 {
		code, msg, c, thc := dal.DataOptThumbsUp(table, preId, exId, plyid, menuid)
		rd.Code = easygo.NewInt32(code)
		rd.Msg = easygo.NewString(msg)
		rd.IsThumbsUp = easygo.NewInt32(c)
		rd.ThumbsUpCount = easygo.NewInt32(thc)
	} else {
		rd.Code = easygo.NewInt32(for_game.C_INFO_NOT_EXISTS)
		rd.Msg = easygo.NewString("菜单类型不存在!")

	}
	logs.Info("RpcESportThumbsUp 返回", rd)
	return rd
}

//获取评论
func (self *sfc) RpcESportGetComment(common *base.Common, reqMsg *client_hall.ESportCommentRequest) *client_hall.ESportCommentReplyListResult {

	logs.Info("RpcESportGetComment 提交", reqMsg)
	m := reqMsg
	menuId := m.GetMenuId()
	parentId := m.GetParentId() //新聞ID或者視頻ID
	commentId := m.GetCommentId()
	page := m.GetPage()
	pageSize := m.GetPageSize()
	table := ""
	orderby := "-CreateTime"
	if commentId > 0 { ///大于0那么就是二级评论
		if for_game.ESPORTMENU_REALTIME == menuId {
			table = for_game.TABLE_ESPORTS_COMMENT_NEWS_REPLY
		} else if for_game.ESPORTMENU_RECREATION == menuId {
			table = for_game.TABLE_ESPORTS_COMMENT_VIDEO_REPLY
		}
		orderby = "+CreateTime"
	} else {
		if for_game.ESPORTMENU_REALTIME == menuId {
			table = for_game.TABLE_ESPORTS_COMMENT_NEWS
		} else if for_game.ESPORTMENU_RECREATION == menuId {
			table = for_game.TABLE_ESPORTS_COMMENT_VIDEO
		}
	}
	lst, total := dal.GetTableESportCommentList(table, page, pageSize, parentId, commentId, orderby)

	var playerIdlist map[int64]*share_message.PlayerBase
	pidlist := []int64{}
	if commentId < 1 && len(lst) > 0 { //如果是一级回复
		plyid := GetPlayerId(common)
		tuobj := for_game.GetRedisESportThumbsUpObj(menuId, parentId, plyid)
		thumbsUpList := tuobj.GetThumbsList()
		for _, v := range lst {
			id := v.GetId()
			if tuobj.IsInThumbsList(thumbsUpList, id) {
				v.IsThumbsUp = easygo.NewInt32(1)
			} else {
				v.IsThumbsUp = easygo.NewInt32(0)
			}

			pidlist = append(pidlist, v.GetPlayerId())
			if v.GetReplyPlayerId() > 0 {
				pidlist = append(pidlist, v.GetPlayerId(), v.GetReplyPlayerId())
			}
		}
	} else {
		for _, v := range lst {
			pidlist = append(pidlist, v.GetPlayerId())
			if v.GetReplyPlayerId() > 0 {
				pidlist = append(pidlist, v.GetPlayerId(), v.GetReplyPlayerId())
			}
		}
	}
	playerIdlist = for_game.GetAllPlayerBase(pidlist, false)

	for _, v := range lst {
		pinfo := playerIdlist[v.GetPlayerId()]
		if pinfo != nil {
			v.PlayerIconUrl = pinfo.HeadIcon
			v.PlayerNickName = pinfo.NickName
		}
		if v.GetReplyPlayerId() > 0 {
			pinfo = playerIdlist[v.GetReplyPlayerId()]
			if pinfo != nil {
				v.ReplyPlayerNickName = pinfo.NickName
			}
		}
	}

	rd := &client_hall.ESportCommentReplyListResult{
		Total: easygo.NewInt32(total),
		List:  lst,
	}
	logs.Info("RpcESportGetComment 返回", len(lst), "共", total)
	//for _, v := range lst {
	//	logs.Info(v)
	//}
	return rd
}

//发送评论
func (self *sfc) RpcESportSendComment(common *base.Common, reqMsg *client_hall.ESportCommentInfo) *client_hall.ESportCommonResult {
	logs.Info("RpcESportSendComment 提交", reqMsg)
	rd := &client_hall.ESportCommonResult{}
	m := reqMsg
	plyid := GetPlayerId(common)

	if plyid < 1 {
		rd.Code = easygo.NewInt32(for_game.C_NOT_LOGIN)
		rd.Msg = easygo.NewString("用户信息不存在")
		return rd
	}
	code := int32(0)
	msg := ""
	dataid := int64(0)

	//判断评论是否违规
	isDirty, dWord := for_game.PDirtyWordsMgr.CheckWord(string(m.GetContent()))

	if isDirty {
		//屏蔽不让发送
		rd.Code = easygo.NewInt32(for_game.C_VIOLATE_CONTENT)
		s := fmt.Sprintf("包含违规词:%s", dWord)
		rd.Msg = easygo.NewString(s)
		rd.DataId = easygo.NewInt64(dataid)
		return rd
	}

	pinfo := for_game.NewRedisPlayerBase(plyid)
	if pinfo == nil {
		rd.Code = easygo.NewInt32(for_game.C_NOT_LOGIN)
		rd.Msg = easygo.NewString("用户信息不存在")
		return rd
	}
	//nickname := pinfo.GetNickName()
	//icon := pinfo.GetHeadIcon()

	if m.GetCommentId() > 0 { //评论id不为了0的时候，我当你是某个评论的回复
		info := &share_message.TableESportComment{
			Content:  m.Content,
			PlayerId: easygo.NewInt64(plyid),
			//PlayerNickName: easygo.NewString(nickname),
			ParentId:  m.ParentId,
			MenuId:    m.MenuId,
			CommentId: m.CommentId,
			//PlayerIconUrl:  easygo.NewString(icon),
			ReplyPlayerId: m.ReplyPlayerId,
			Status:        easygo.NewInt32(1),
		}
		code, msg, dataid = dal.AddSportCommentReply(info)
	} else {
		info := &share_message.TableESportComment{
			Content:       m.Content,
			ThumbsUpCount: easygo.NewInt32(0),
			PlayerId:      easygo.NewInt64(plyid),
			//PlayerNickName: easygo.NewString(nickname),
			ParentId: m.ParentId,
			MenuId:   m.MenuId,
			//AppLabelID:     easygo.NewInt64(0),
			ReplyCount: easygo.NewInt32(0),
			//PlayerIconUrl:  easygo.NewString(icon),
			Status: easygo.NewInt32(1),
		}
		code, msg, dataid = dal.AddSportComment(info)
	}

	rd.Code = easygo.NewInt32(code)
	rd.Msg = easygo.NewString(msg)
	rd.DataId = easygo.NewInt64(dataid)
	//logs.Info("RpcESportSendComment 返回", rd)
	return rd
}

//删除评论
func (self *sfc) RpcESportDeleteComment(common *base.Common, reqMsg *client_hall.ESportDeleteCommentInfo) *client_hall.ESportCommonResult {

	rd := &client_hall.ESportCommonResult{}

	if reqMsg.GetCommentType() == 1 {
		code, msg := dal.DeleteSportCommentReply(reqMsg.GetMenuId(), reqMsg.GetParentId(), reqMsg.GetPCommentId(), reqMsg.GetCommentId())
		rd.Code = easygo.NewInt32(code)
		rd.Msg = easygo.NewString(msg)
	} else if reqMsg.GetCommentType() == 2 {
		code, msg := dal.DeleteSportComment(reqMsg.GetMenuId(), reqMsg.GetParentId(), reqMsg.GetCommentId())
		rd.Code = easygo.NewInt32(code)
		rd.Msg = easygo.NewString(msg)
	} else {
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString("类型错误")
	}
	return rd
}

//获取标签配置
func (self *sfc) RpcESportGetAllLabelList(common *base.Common, reqMsg *client_hall.ESportInfoRequest) *client_hall.ESportLabelList {
	rd := &client_hall.ESportLabelList{}
	rd.LabelList = dal.GetTableESPortsLabelList(reqMsg.GetMenuId())
	return rd
}

//轮播图列表
func (self *sfc) RpcESportGetCarouselList(common *base.Common, reqMsg *client_hall.ESportInfoRequest) *client_hall.ESportCarouselList {
	logs.Info("RpcESportGetCarouselList 提交", reqMsg)
	rd := &client_hall.ESportCarouselList{}
	list := dal.GetTableESPortsCarouselByMenuId(reqMsg.GetMenuId(), 1)
	rd.CarouselList = list
	logs.Info("RpcESportGetCarouselList 返回", rd)
	return rd
}

/*
func CreateESPortsPlayer(plyId int64) *for_game.RedisESportPlayerObj {

	var pinfo *for_game.RedisESportPlayerObj
	pdata := &share_message.TableESPortsPlayer{
		Id:                easygo.NewInt64(plyId),
		Status:            easygo.NewInt32(for_game.ESPORT_PLAYER_STATUS_1),
		LastPullTime:      easygo.NewInt64(0),
		CurrentRoomLiveId: easygo.NewInt64(0),
	}
	code, _ := dal.CreateTableESPortsPlayer(pdata)
	if code == for_game.C_OPT_SUCCESS {
		pinfo = for_game.NewRedisESportPlayerObj(plyId, pdata)
		logs.Info("创建电竞用户: %d", plyId)
	}
	return pinfo
}
*/

//加载系统消息
func LoadSysMsg(pinfo *for_game.RedisESportPlayerObj, recipientType int32, rd *client_hall.ESPortsSysMsgList) {

	list, code := dal.GetTableESportAllSysList(pinfo.GetLastPullTime(), recipientType)
	logs.Info("GetTableESportAllSysList", list)
	if code == for_game.C_OPT_SUCCESS {
		if list != nil || len(list) > 0 {
			rd.SysMsgList = list
		}
		pinfo.SetLastPullTime(easygo.NowTimestamp())
	} else {
		logs.Info("获取系统消息数据失败 用户ID：", pinfo.PlayerId)
	}
}

//加载比赛订单
func LoadGameOrderSysMsg(pinfo *for_game.RedisESportPlayerObj, rd *client_hall.ESPortsSysMsgList) {
	plyid := pinfo.PlayerId
	list1, code := dal.GetTableESportAllGameOrderSysMsgList(plyid)
	logs.Info("GetTableESportAllGameOrderSysMsgList", list1)
	if code == for_game.C_OPT_SUCCESS {
		if list1 != nil || len(list1) > 0 {
			rd.UnPayed = list1
			dal.DeleteTableESPortsGameOrderSysMsg(plyid)
		}
	}
	list2, code := dal.GetTableESportAllGameOrderEndSysMsgList(plyid)
	logs.Info("GetTableESportAllGameOrderEndSysMsgList", list2)
	if code == for_game.C_OPT_SUCCESS {
		if list2 != nil || len(list2) > 0 {
			rd.Payed = list2
			dal.ClaerTableESPortsGameOrderEndSysMsg(plyid)
		}
	}
}

//获取系统消息列表
func (self *sfc) RpcESportGetSysMsgList(common *base.Common, reqMsg *client_hall.ESportInfoRequest) *client_hall.ESPortsSysMsgList {
	rd := &client_hall.ESPortsSysMsgList{}

	logs.Info("RpcESportGetSysMsgList ", reqMsg)
	plyId := GetPlayerId(common)
	isNew := false
	pinfo := for_game.GetRedisESportPlayerObj(plyId)
	if pinfo == nil {
		isNew = true
		pinfo = GetOrCraeteESPortsPlayer(plyId) //CreateESPortsPlayer(plyId)
	}
	if pinfo != nil {
		rd = &client_hall.ESPortsSysMsgList{
			PlayerId: easygo.NewInt64(plyId),
		}
		pbase := for_game.GetRedisPlayerBase(plyId)
		if pbase != nil {
			//接收者类型 0 全体，1 IOS,2 Android
			recipientType := int32(0)
			////设备类型 1 IOS，2 Android，3 PC
			recipientType = pbase.GetDeviceType()
			LoadSysMsg(pinfo, recipientType, rd) //加载系统消息
		}

		//
		if !isNew { //不是新人就查订单消息
			LoadGameOrderSysMsg(pinfo, rd)
		}
	} else {
		logs.Info("NewRedisESportPlayerObj", "nil")
	}
	logs.Info("RpcESportGetSysMsgList 返回", rd)

	return rd
}

//加载作者昵称
func LoadVideoAuthor(list []*share_message.TableESPortsVideoInfo) []*share_message.TableESPortsVideoInfo {
	pidlist := []int64{}
	for _, v := range list {
		pidlist = append(pidlist, v.GetAuthorPlayerId())
	}
	playerIdlist := for_game.GetAllPlayerBase(pidlist, false)

	for _, v := range list {
		pinfo := playerIdlist[v.GetAuthorPlayerId()]
		if pinfo != nil {
			v.Author = pinfo.NickName
		}
	}
	return list
}

//获取娱乐视频或者放映厅列表
func (self *sfc) RpcESportGetVideoList(common *base.Common, reqMsg *client_hall.ESportVideoPageRequest) *client_hall.ESportVideoListResult {
	rd := &client_hall.ESportVideoListResult{}
	logs.Info("RpcESportGetVideoList 提交", reqMsg)
	page := reqMsg.GetPage()
	pageSize := reqMsg.GetPageSize()
	sort := "-BeginEffectiveTime"
	if reqMsg.AscOrDesc != nil && reqMsg.GetAscOrDesc() != "" && reqMsg.OrderField != nil && reqMsg.GetOrderField() != "" {
		sort = reqMsg.GetAscOrDesc() + reqMsg.GetOrderField()
	}
	typeId := reqMsg.GetTypeId()
	labelId := reqMsg.GetLabelId()
	videoType := reqMsg.GetVideoType()
	plyId := GetPlayerId(common)

	list, count := dal.GetTableESPortsVideoItemList(page, pageSize, sort, typeId, labelId, videoType, plyId)
	list = LoadVideoAuthor(list)
	rd.List = list
	rd.Total = easygo.NewInt32(count)
	logs.Info("RpcESportGetVideoList 返回 當前：(%d) ,共：%d", len(list), count)
	return rd
}

//是否关注
func GetIsAuthorFollow(fans []int64, playerId int64) int32 {
	b32 := int32(2)
	for _, v := range fans {
		if v == playerId {
			b32 = 1
		}
	}
	return b32
}

//是否关注放映厅
func GetIsRoomFollow(liveId int64, playerId int64) int32 {
	f_onj := for_game.GetRedisESportLiveFollowObj(playerId)
	if f_onj.IsFollow(liveId) {
		return 1
	} else {
		return 2
	}
}

//进入该房间并且通知其他人
func EnterRoom(liveId int64, playerId int64) {
	mypinfo := for_game.GetRedisPlayerBase(playerId)
	if mypinfo != nil {
		room := for_game.GetRedisLiveRoomPlayerObj(liveId)
		if room != nil {
			esport_p := GetOrCraeteESPortsPlayer(playerId) //for_game.GetRedisESportPlayerObj(playerId)
			if esport_p != nil {
				room.EnterRoom(playerId)
				esport_p.SetCurrentRoomLiveId(liveId)
			} else {
				logs.Error("找不到入房間的用戶 ： %d", playerId)
			}
			msg := &share_message.TableESPortsLiveRoomMsgLog{
				NickName:       easygo.NewString(mypinfo.GetNickName()),
				HeadIcon:       easygo.NewString(mypinfo.GetHeadIcon()),
				LiveId:         easygo.NewInt64(liveId),
				DataType:       easygo.NewInt32(for_game.LIVE_ROOM_OPT_2), //进入直播间
				SenderPlayerId: easygo.NewInt64(playerId),
			}
			plist := room.GetPlayerIds()
			smsg := func() {
				SendMsgToHallClientNew(plist, "RpcESportNewRoomMsg", msg)
			}
			easygo.Spawn(smsg)
		} else {
			logs.Error("找不到房間 ： %d", liveId)
		}

	} else {
		logs.Error("找不到im用户 ： %d", playerId)
	}
}

//加载放映厅作者信息
func LoadAuthorInfo(data *share_message.TableESPortsVideoInfo, playerId int64) *share_message.TableESPortsVideoInfo {
	apinfo := for_game.GetRedisPlayerBase(data.GetAuthorPlayerId())
	if apinfo != nil {

		pHeadIcon := apinfo.GetHeadIcon()
		pNickName := apinfo.GetNickName()
		fs := apinfo.GetFans()
		data.FanCount = easygo.NewInt32(len(fs))
		data.IsAuthorFollow = easygo.NewInt32(GetIsAuthorFollow(fs, playerId))
		data.Author = easygo.NewString(pNickName)
		data.PlayerIconUrl = easygo.NewString(pHeadIcon)
	}
	return data
}

//获取娱乐视频或者放映厅数据
func (self *sfc) RpcESportGetVideoInfo(common *base.Common, reqMsg *client_hall.ESportVideoRequest) easygo.IMessage {
	logs.Info("RpcESportGetVideoInfo 提交", reqMsg)
	rd := &client_hall.ESportVideoResult{}
	id := reqMsg.GetDataId()

	if id < 1 {
		rd.Msg = easygo.NewString("视频ID不能小于0")
		rd.Code = easygo.NewInt32(for_game.C_INFO_NOT_EXISTS)
		return rd
	}
	data := dal.GetTableESPortsVideoInfoById(reqMsg.GetDataId())
	if data == nil || data.GetId() < 1 {
		rd.Msg = easygo.NewString("数据不存在")
		rd.Code = easygo.NewInt32(for_game.C_INFO_NOT_EXISTS)
	} else {
		liveId := data.GetId()
		if data.GetStatus() != for_game.ESPORTS_NEWS_STATUS_1 {
			data.VideoUrl = easygo.NewString("")
		}
		plyid := GetPlayerId(common)
		flowType := int64(for_game.ESPORT_FLOW_VIDEO_HISTORY)
		menuid := for_game.ESPORTMENU_RECREATION
		if data.GetVideoType() == int64(for_game.ESPORTS_VIDEO_TYPE_2) {
			menuid = for_game.ESPORTMENU_LIVE
			flowType = int64(for_game.ESPORT_FLOW_LIVE_HISTORY)
			data = LoadAuthorInfo(data, plyid) //加载作者信息
			data.IsFollow = easygo.NewInt32(GetIsRoomFollow(liveId, plyid))
			//data.IsFollow = easygo.NewInt32(2)
			EnterRoom(liveId, plyid) //设置用户进入放映厅

		} else {
			_, isThumbsUp := dal.IsThumbsUp(int32(menuid), 0, reqMsg.GetDataId(), plyid)
			if isThumbsUp {
				data.IsThumbsUp = easygo.NewInt32(1)
			} else {
				data.IsThumbsUp = easygo.NewInt32(0)
			}
		}
		if data.GetAuthorPlayerId() != plyid {
			dal.AddTableESPortsFlowInfo(flowType, plyid, liveId) //添加播放历史
		}

		rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
		rd.Data = data
	}
	logs.Info("RpcESportGetVideoInfo 返回", rd)
	return rd
}

//离开直播放映厅
func (self *sfc) RpcESportLeaveLive(common *base.Common, reqMsg *client_hall.ESportCommonResult) *client_hall.ESportCommonResult {

	rd := &client_hall.ESportCommonResult{}
	room := for_game.GetRedisLiveRoomPlayerObj(reqMsg.GetDataId())
	plyid := GetPlayerId(common)
	if room != nil {
		room.LeaveRoom(plyid)
	} else {
		logs.Error("离开房间时，无法获取放映厅房间信息", reqMsg.GetDataId())
	}
	pobj := for_game.GetRedisESportPlayerObj(plyid)
	if pobj != nil {
		pobj.SetCurrentRoomLiveId(0)
	}
	rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
	return rd
}

//获取观看历史视频或者放映厅列表（近期）
func (self *sfc) RpcESportGetMyHistoryVideoList(common *base.Common, reqMsg *client_hall.ESportVideoPageRequest) *client_hall.ESportVideoListResult {
	rd := &client_hall.ESportVideoListResult{}
	typeid := reqMsg.GetTypeId()
	logs.Info("RpcESportGetMyHistoryVideoList 提交", reqMsg)
	page := reqMsg.GetPage()
	pageSize := reqMsg.GetPageSize()
	plyid := GetPlayerId(common)
	lst, c := dal.GetFlowVideoOrLiveListByPlayerId(typeid, page, pageSize, "CreateTime", plyid)
	lst = LoadVideoAuthor(lst)
	rd.List = lst
	rd.Total = easygo.NewInt32(c)
	logs.Info("RpcESportGetMyHistoryVideoList 返回", rd)
	return rd
}

//获取放映厅首页的数据
func (self *sfc) RpcESportGetLiveHomeInfo(common *base.Common, reqMsg *client_hall.ESportInfoRequest) *client_hall.ESportLiveHomeInfo {
	logs.Info("RpcESportGetLiveHomeInfo 提交", reqMsg)

	rd := &client_hall.ESportLiveHomeInfo{}
	plyid := GetPlayerId(common)
	rd.LabelList = dal.GetTableESPortsLabelList(reqMsg.GetMenuId())
	rd.MyLiveInfo = dal.GetMyliveInfoByPlayerId(plyid)
	logs.Info("RpcESportGetLiveHomeInfo 返回", rd)
	return rd
}

//添加关注放映厅
func (self *sfc) RpcESportAddFollowLive(common *base.Common, reqMsg *client_hall.ESportInfoRequest) *client_hall.ESportCommonResult {
	logs.Info("RpcESportAddFollowLive 提交", reqMsg)
	playerId := GetPlayerId(common)
	dataType := reqMsg.GetGameTypeId()
	dataId := reqMsg.GetExtId()

	code := int32(for_game.C_SYS_ERROR)
	msg := "系統錯誤"

	obj := for_game.GetRedisESportLiveFollowObj(playerId)
	if obj != nil {
		if obj.IsFollow(dataId) {
			b := dal.DeleteTableESPortsFlowInfoEx(dataType, playerId, dataId)
			if b {
				obj.CancelFollow(dataId)
				code, msg = for_game.C_OPT_SUCCESS, ""
				code, msg = dal.UpdateFedAddition_xv(for_game.TABLE_ESPORTS_VIDEO, "FlowCount", dataId, -1)
			} else {
				code, msg = for_game.C_SYS_ERROR, "系統錯誤"
			}
		} else {
			code, msg = dal.AddTableESPortsFlowfollowInfo(dataType, playerId, dataId)
			if code == for_game.C_OPT_SUCCESS {
				obj.AddFollow(dataId)
				//dal.UpdateFedAddition(for_game.TABLE_ESPORTS_VIDEO, "FlowCount", dataId, 1)
				code, msg = dal.UpdateFedAddition_xv(for_game.TABLE_ESPORTS_VIDEO, "FlowCount", dataId, 1)
			}
		}
	}
	rd := client_hall.ESportCommonResult{
		Code: easygo.NewInt32(code),
		Msg:  easygo.NewString(msg),
	}
	logs.Info("RpcESportAddFollowLive 返回", rd)
	return &rd
}

//申请放映厅
func (self *sfc) RpcESportApplyOpenLive(common *base.Common, m *client_hall.ESportMyLiveRoomInfo) *client_hall.ESportCommonResult {
	logs.Info("RpcESportApplyOpenLive 提交", m)
	rd := &client_hall.ESportCommonResult{}
	playerId := GetPlayerId(common)
	title := m.GetTitle()
	content := m.GetContent()
	uniqueGameId := m.GetUniqueGameId()
	applabelId := m.GetAppLabelID()
	videoUrl := m.GetVideoUrl()
	imageUrl := m.GetCoverImageUrl()
	status := int32(1) //0未发布 1(审核通过) 2审核拒绝 3审核不通过 4已过期
	note := ""         //原因
	logs.Info("RpcESportApplyOpenLive 開始審核")
	//违规词判定
	//判断放映厅名称是否违规
	isDirtyTitle, dWordTitle := for_game.PDirtyWordsMgr.CheckWord(title)

	if isDirtyTitle {
		//屏蔽不让发送
		rd.Code = easygo.NewInt32(for_game.C_VIOLATE_CONTENT)
		s := fmt.Sprintf("名称包含违规词:%s,请重新编辑", dWordTitle)
		rd.Msg = easygo.NewString(s)
		logs.Info("RpcESportApplyOpenLive 返回", rd)
		return rd
	}

	num := ImageModeration(imageUrl)
	if num != 100 {
		//屏蔽不让发送
		rd.Code = easygo.NewInt32(for_game.C_VIOLATE_CONTENT)
		rd.Msg = easygo.NewString("封面图违规,请重新编辑")
		logs.Info("RpcESportApplyOpenLive 返回", rd)
		return rd
	}

	//判断公告是否违规
	isDirtyContent, dWordContent := for_game.PDirtyWordsMgr.CheckWord(content)

	if isDirtyContent {
		//屏蔽不让发送
		rd.Code = easygo.NewInt32(for_game.C_VIOLATE_CONTENT)
		s := fmt.Sprintf("公告包含违规词:%s,请重新编辑", dWordContent)
		rd.Msg = easygo.NewString(s)
		logs.Info("RpcESportApplyOpenLive 返回", rd)
		return rd
	}
	logs.Info("RpcESportApplyOpenLive 結束審核")
	code, msg := dal.ApplyMylive(playerId, title, content, uniqueGameId, applabelId, videoUrl, imageUrl, status, note, m.GetUniqueGameName())
	//if code == for_game.C_OPT_SUCCESS { //写入数据库成功
	//代表通过
	//}
	rd.Msg = easygo.NewString(msg)
	rd.Code = easygo.NewInt32(code)
	logs.Info("RpcESportApplyOpenLive 返回", rd)
	return rd
}

//放映厅发言
func (self *sfc) RpcESportSendLiveRoomMsg(common *base.Common, reqMsg *client_hall.ESportCommentInfo) *client_hall.ESportCommonResult {
	rd := &client_hall.ESportCommonResult{}
	playerId := GetPlayerId(common)
	liveId := reqMsg.GetCommentId()
	if liveId < 1 {
		rd.Code = easygo.NewInt32(for_game.C_INFO_NOT_EXISTS)
		rd.Msg = easygo.NewString(fmt.Sprintf("房间号不能为 %d", liveId))
		return rd
	}
	obj := for_game.GetRedisPlayerBase(playerId)
	if obj != nil {

		msg := &share_message.TableESPortsLiveRoomMsgLog{
			Id:             nil,
			NickName:       easygo.NewString(obj.GetNickName()),
			Content:        easygo.NewString(reqMsg.GetContent()),
			HeadIcon:       easygo.NewString(obj.GetHeadIcon()),
			LiveId:         easygo.NewInt64(liveId),
			DataType:       easygo.NewInt32(for_game.LIVE_ROOM_OPT_1),
			SenderPlayerId: easygo.NewInt64(playerId),
			CreateTime:     easygo.NewInt64(easygo.NowTimestamp()),
		}

		room := for_game.GetRedisLiveRoomPlayerObj(liveId)
		plist := room.GetPlayerIds()
		logs.Info("RpcESportSendLiveRoomMsg", plist)
		//sm := func() {
		if len(plist) < 1 {
			rd.Code = easygo.NewInt32(for_game.C_INFO_NOT_EXISTS)
			rd.Msg = easygo.NewString(fmt.Sprintf("房间%d不存在", liveId))
			return rd
		}
		roomobj := for_game.GetRedisESportRoomChatLogObj(liveId)
		id := for_game.NextId(for_game.TABLE_ESPORTS_ROOM_CHAT_MSG_LOG)
		if roomobj != nil {
			msg.Id = easygo.NewInt64(id)
			roomobj.AddRedisESportRoomChatLog(msg)
		} else {
			dal.AddTableESPortsRoomChatMsg(msg, id)
		}
		SendMsgToHallClientNew(plist, "RpcESportNewRoomMsg", msg)
		logs.Info(obj.GetNickName(), liveId, "说："+reqMsg.GetContent())

		rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
		//}
		//easygo.Spawn(sm)
	} else {
		rd.Code = easygo.NewInt32(for_game.C_INFO_NOT_EXISTS)
		rd.Msg = easygo.NewString("用户不存在")
	}

	return rd
}

//获取比赛名字和Id列表
func (self *sfc) RpcESportGetESPortsGameViewList(common *base.Common, reqMsg *client_hall.ESportPageRequest) *client_hall.ESPortsGameItemViewResult {

	logs.Info("RpcESportGetESPortsGameViewList", reqMsg)
	rd := &client_hall.ESPortsGameItemViewResult{}
	lst := []*share_message.ESPortsGameItemView{}
	gamelst, total := dal.GetESPortsGameItemViewList(reqMsg.GetPage(), reqMsg.GetPageSize(), reqMsg.GetLabelId())

	if total > 0 { //测试数据
		for _, v := range gamelst {
			tan := "A队未知名称"
			if v.GetTeamA() != nil {
				tan = v.GetTeamA().GetName()
			}
			tbn := "B队未知名称"
			if v.GetTeamB() != nil {
				tbn = v.GetTeamB().GetName()
			}
			matchVsName := for_game.GetMatchVSName(v.GetMatchName(), v.GetMatchStage(), v.GetBo(), tan, tbn)
			lst = append(lst, &share_message.ESPortsGameItemView{
				Id:          v.Id,
				MatchVsName: easygo.NewString(matchVsName),
			})
		}
	}

	rd.Total = easygo.NewInt32(total)
	rd.List = lst
	logs.Info("RpcESportGetESPortsGameViewList 返回", rd)
	return rd
}

//获取比赛相关放映厅
func (self *sfc) RpcESportGetGameVideoList(common *base.Common, reqMsg *client_hall.ESportGameViewPageRequest) *client_hall.ESportVideoListResult {
	rd := &client_hall.ESportVideoListResult{}
	logs.Info("RpcESportGetGameVideoList 提交", reqMsg)
	page := reqMsg.GetPage()
	pageSize := reqMsg.GetPageSize()
	sort := "-CreateTime"

	list, count := dal.GetTableESPortsGameVideoItemList(page, pageSize, sort, reqMsg.GetUniqueGameId())
	list = LoadVideoAuthor(list)
	rd.List = list
	rd.Total = easygo.NewInt32(count)
	logs.Info("RpcESportGetGameVideoList 返回", rd)
	return rd
}

func BpsLog(f interface{}, v ...interface{}) {
	//if for_game.IS_FORMAL_SERVER {
	logs.Info(f, v...)
	//}
}

//埋点点击
func (self *sfc) RpcESPortsBpsClick(common *base.Common, reqMsg *client_hall.ESPortsBpsClickRequest) *client_hall.ESportCommonResult {
	rd := &client_hall.ESportCommonResult{}
	BpsLog("RpcESPortsNpsClick 提交", reqMsg)
	plyid := GetPlayerId(common)
	if reqMsg.GetPageType() == for_game.ESPORT_BPS_PAGE_TYPE_1 && reqMsg.GetNavigationId() == for_game.ESPORT_MODLE_4 { //进入电竞模块
		GetOrCraeteESPortsPlayer(plyid) //如果用户不存在就创建
	}
	if reqMsg.GetPageType() != for_game.ESPORT_BPS_PAGE_TYPE_1 && reqMsg.GetNavigationId() != for_game.ESPORT_MODLE_4 {
		logs.Info("<----------------------------------------------------->")
		logs.Warn("<----------", reqMsg, "---------->")
		logs.Info("<----------------------------------------------------->")
	}

	ckobj := for_game.GetRedisESportBpsClickLogObj(plyid)

	ckobj.BpsClick(reqMsg)
	BpsLog("RpcESPortsNpsClick 返回", rd)
	return rd
}

//埋点点击列表
func (self *sfc) RpcESPortsBpsClickList(common *base.Common, reqMsg *client_hall.ESPortsBpsClickListRequest) *client_hall.ESportCommonResult {
	rd := &client_hall.ESportCommonResult{}
	BpsLog("RpcESPortsBpsClickList 提交", reqMsg)

	plyid := GetPlayerId(common)
	ckobj := for_game.GetRedisESportBpsClickLogObj(plyid)
	if ckobj != nil {
		ckobj.BpsClickEx(reqMsg)
	} else {
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString("系统错误，找不到点击redis缓存")
	}
	BpsLog("RpcESPortsBpsClickList 返回", rd)
	return rd
}

//埋点停留时长
func (self *sfc) RpcESPortsBpsDuration(common *base.Common, reqMsg *client_hall.ESPortsBpsDurationRequest) *client_hall.ESportCommonResult {
	rd := &client_hall.ESportCommonResult{}
	BpsLog("RpcESPortsNpsDuration 提交", reqMsg)
	plyid := GetPlayerId(common)
	//easygoNo
	pinfo := for_game.GetRedisESportBpsDurationLogObj(plyid)
	if pinfo != nil && reqMsg.List != nil {
		for _, it := range reqMsg.List {
			if it.GetOpt() == 1 {
				pinfo.BeginCurrentBpsDuration(it.GetPageType(), it.GetMenuId(), it.GetLabelId(), it.GetExTabId(), it.GetDataId(), it.GetExId(), it.GetNavigationId())
			}
			if it.GetOpt() == 2 {
				pinfo.EndCurrentBpsDuration(it.GetPageType()) //, reqMsg.GetMenuId(), reqMsg.GetLabelId(), reqMsg.GetExTabId(), reqMsg.GetDataId(), reqMsg.GetExId(), reqMsg.GetNavigation())
			}
		}
	}
	BpsLog("RpcESPortsNpsDuration 返回", rd)
	return rd
}

//电竞币兑换页
func (self *sfc) RpcESPortsCoinView(common *base.Common, reqMsg *base.Empty) *client_hall.ESPortsCoinViewResult {
	logs.Info("===RpcESPortsCoinView===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	playId := common.GetUserId()

	rd := &client_hall.ESPortsCoinViewResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}

	//取得配置信息
	GetSportExChangeConfig(playId, rd)

	if rd.GetCode() != for_game.C_OPT_SUCCESS {
		logs.Error("=========RpcESPortsCoinView 返回===========", rd)
		return rd
	}

	rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
	rd.Msg = easygo.NewString("")

	logs.Info("=========RpcESPortsCoinView 返回===========", rd)

	return rd
}

//电竞币兑换动作
func (self *sfc) RpcESPortsCoinExChange(common *base.Common,
	reqMsg *client_hall.ESPortsCoinExChangeRequest) *client_hall.ESPortsCoinExChangeResult {
	logs.Info("=========RpcESPortsCoinExChange========,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	//类型
	exchangeType := reqMsg.GetType()
	//扣除的硬币
	coins := reqMsg.GetExChangeObject().GetCoin()
	//电竞币额度
	esportCoins := reqMsg.GetExChangeObject().GetESportCoin()

	//首充
	firstGive := reqMsg.GetExChangeObject().GetFirstGive()

	//日常送
	dailyGive := reqMsg.GetExChangeObject().GetDailyGive()

	//活动
	exchangeRates := reqMsg.GetExChangeObject().GetRate()

	playId := common.GetUserId()

	rd := &client_hall.ESPortsCoinExChangeResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString("兑换成功"),
	}

	//返回给前端显示用
	rdView := &client_hall.ESPortsCoinViewResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}

	//活动是否停止
	if exchangeType == for_game.ESPORT_EXCHANGE_TYPE_3 {

		//是否活动以及活动是否有效期
		activeConfig := for_game.GetRedisActiveConfig()
		nowTime := time.Now().Unix()

		if nil != activeConfig {
			//活动关闭并且不在活动期间
			if activeConfig.GetStatus() == 1 ||
				activeConfig.GetEndTime() <= nowTime || activeConfig.GetStartTime() > nowTime {

				rd.Code = easygo.NewInt32(for_game.C_ACTIVE_STOP)
				rd.Msg = easygo.NewString("活动停止")
				//重置最新页面
				ResetESportExChangeView(playId, rd, rdView)
				return rd
			}
		}
	}

	//先扣硬币
	//==============先扣硬币    开始====================
	st := for_game.COIN_TYPE_ESPORT_EXCHANGE_OUT
	msg := fmt.Sprintf("电竞兑换扣除硬币[%d]个", coins)

	var tempCoins int64
	if coins > 0 {
		tempCoins = -coins
	} else {
		tempCoins = coins
	}

	//取得流水订单号
	streamOrderId := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_OUT, st)
	req := &share_message.ESportCoinRecharge{
		PlayerId:     easygo.NewInt64(common.GetUserId()),
		RechargeCoin: easygo.NewInt64(tempCoins),
		SourceType:   easygo.NewInt32(st),
		Note:         easygo.NewString(msg),
		ExtendLog: &share_message.GoldExtendLog{
			OrderId: easygo.NewString(streamOrderId), //流水的订单号
		},
	}

	result, err := dal.SendMsgToServerNewEx(PServerInfoMgr, playId, "RpcESportSendChangeCoins", req) //通知大厅
	//传输出错
	if err != nil {
		logs.Error(err.GetReason())
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString(err.GetReason())
		//重置最新页面
		ResetESportExChangeView(playId, rd, rdView)
		return rd
	}
	//结果出错
	if nil != result {
		rst, ok := result.(*client_hall.ESportCommonResult)
		if ok && nil != rst {
			if rst.GetCode() == for_game.C_DEDUCT_MONEY_FAIL {
				rd.Code = easygo.NewInt32(for_game.C_DEDUCT_MONEY_FAIL)
				rd.Msg = easygo.NewString(rst.GetMsg())
				//重置最新页面
				ResetESportExChangeView(playId, rd, rdView)
				return rd
			}
		}
	}
	//================先扣硬币   结束========================

	//================先结算正常的兑换电竞币额度  开始====================
	ExChangeESportCoins(playId, esportCoins, rd, rdView, for_game.ESPORTS_EXCHANGE_NORMAL)
	//如果没有、插入一条
	if IsFirstExChange(playId, rdView) {
		InsFirstExChange(playId)
	}
	//================先结算正常的兑换电竞币额度   结束===================

	//白名单优先级最高百分百赠送=====开始======
	isWhite := IsESportExChangeWhite(playId)
	if isWhite {
		//这里要直接返回
		return ExChangeESportCoins(playId, esportCoins, rd, rdView, for_game.ESPORTS_EXCHANGE_WHITE)
	}
	//白名单优先级最高百分百赠送=====结束======

	//赠送类型
	switch exchangeType {
	case for_game.ESPORT_EXCHANGE_TYPE_3:
		//通过设置的活动相关参数算出要送多少币
		activeCoins := GetActiveESportCoins(exchangeRates, esportCoins)

		if activeCoins > 0 {
			return ExChangeESportCoins(playId, activeCoins, rd, rdView, for_game.ESPORTS_EXCHANGE_ACTIVE)
		}

	case for_game.ESPORT_EXCHANGE_TYPE_2:
		//首充
		if firstGive > 0 {
			return ExChangeESportCoins(playId, firstGive, rd, rdView, for_game.ESPORTS_EXCHANGE_FIRSRT)
		}
	case for_game.ESPORT_EXCHANGE_TYPE_1:
		//日常赠送
		if dailyGive > 0 {
			return ExChangeESportCoins(playId, dailyGive, rd, rdView, for_game.ESPORTS_EXCHANGE_DAY)
		}
	default:
	}
	rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
	rd.Msg = easygo.NewString("兑换成功")

	logs.Info("=========RpcESPortsCoinExChange 返回===========", rd)

	return rd
}

//电竞币兑换流水
func (self *sfc) RpcESPortsCoinExChangeRecord(common *base.Common, reqMsg *client_hall.ESPortsCoinExChangeRecordRequest) *client_hall.ESPortsCoinExChangeRecordResult {
	logs.Info("===RpcESPortsCoinExChangeRecord===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	//先压库
	for_game.SaveESportCoinChangeLogToMongoDB()

	rd := &client_hall.ESPortsCoinExChangeRecordResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}

	//取得用户的兑换流水
	findBson := bson.M{}
	findBson["PlayerId"] = common.GetUserId()
	findBson["$or"] = []bson.M{{"SourceType": for_game.ESPORTCOIN_TYPE_EXCHANGE_IN}, {"SourceType": for_game.ESPORTCOIN_TYPE_EXCHANGE_GIVE_IN}}
	sort := []string{"-CreateTime"}
	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTCHANGELOG, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetPage()), sort...)
	var list []*client_hall.ExChangeRecordObject

	for _, li := range lis {
		var title string
		if int32(li.(bson.M)["SourceType"].(int)) == for_game.ESPORTCOIN_TYPE_EXCHANGE_IN {
			title = "兑换"
		} else if int32(li.(bson.M)["SourceType"].(int)) == for_game.ESPORTCOIN_TYPE_EXCHANGE_GIVE_IN {
			title = "赠送"
		}

		one := &client_hall.ExChangeRecordObject{
			ChangeESportCoin: easygo.NewInt64(li.(bson.M)["ChangeESportCoin"].(int64)),
			CreateTime:       easygo.NewInt64(li.(bson.M)["CreateTime"].(int64)),
			Title:            easygo.NewString(title),
		}
		list = append(list, one)
	}

	rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
	rd.Msg = easygo.NewString("")
	rd.Total = easygo.NewInt32(count)
	rd.ExChangeRecordList = list

	logs.Info("=========RpcESPortsCoinExChangeRecord 返回===========", rd)

	return rd
}

//取得api平台来源
func (self *sfc) RpcESPortsApiOrigin(common *base.Common, reqMsg *base.Empty) *client_hall.RpcESPortsApiOriginResult {
	logs.Info("===RpcESPortsApiOrigin===,common=%v,msg=%v", common, reqMsg) // 别删，永久留存

	rd := &client_hall.RpcESPortsApiOriginResult{
		Code: easygo.NewInt32(for_game.C_OPT_SUCCESS),
		Msg:  easygo.NewString(""),
	}

	// TODO 等待2版后台功能
	//默认给与野子科技
	rd.ApiOrigin = easygo.NewInt32(for_game.ESPORTS_API_ORIGIN_ID_YEZI)
	logs.Info("=========RpcESPortsApiOrigin 返回===========", rd)

	return rd
}
