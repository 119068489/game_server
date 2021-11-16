// 管理后台为[浏览器]提供的服务
//客服服务

package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/share_message"
)

//登录获取未读消息条数
func (self *cls4) RpcGetIMmessageCount(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	if user.GetRole() != 2 {
		return easygo.NewFailMsg("只有客服账号才能查询消息条数")
	}

	oldMsg := GetWaiterMsg(user.GetId())
	var count int32
	if len(oldMsg) > 0 {
		for _, v := range oldMsg {
			count += v.GetSnew()
		}
	}
	msg := &brower_backstage.IMmessageResponse{
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//客服查询正在沟通的消息列表
func (self *cls4) RpcGetWaiterMsg(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	if user.GetRole() != 2 {
		return easygo.NewFailMsg("只有客服账号才能查询消息")
	}
	list := GetWaiterMsg(user.GetId())
	lis := []*share_message.IMmessage{}
	for _, c := range list {
		content := c.GetContent()
		start := 1
		if len(content)-1 > 0 {
			start = len(content) - 1
			c.Content = content[start:]
		}

		// c.Content = []*share_message.IMcontent{}
		lis = append(lis, c)
	}
	msg := &brower_backstage.IMmessageNopageResponse{
		List: lis,
	}

	return msg
}

// 查询指定条数的消息
func (self *cls4) RpcGetWaiterMsgByMid(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	if user.GetRole() != 2 {
		return easygo.NewFailMsg("只有客服账号才能查询消息")
	}

	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		return easygo.NewFailMsg("消息ID不能为空")
	}

	im := &share_message.IMmessage{}
	im = for_game.QueryIMmessageByMid(reqMsg.GetId(), for_game.WAITER_MESSAGE_ING)
	if im != nil {
		new := im.GetCnew()
		if new > 0 && reqMsg.GetType() == 1 {
			content := im.GetContent()
			start := len(content) - int(new)
			im.Content = content[start:]
		}
		for_game.UpdateMessageRead(im.GetId(), 2) //前端查看消息，修改消息已读状态
	} else {
		im = &share_message.IMmessage{}
	}

	return im
}

//客服发消息给玩家
func (self *cls4) RpcWaiterSendMsgToPlayer(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.IMmessage) easygo.IMessage {
	if user.GetRole() != 2 {
		return easygo.NewFailMsg("只有客服账号才能发送消息")
	}

	if reqMsg.PlayerId == nil || reqMsg.GetPlayerId() == 0 {
		return easygo.NewFailMsg("用户Id错误")
	}

	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		return easygo.NewFailMsg("消息Id错误")
	} else {
		im := for_game.QueryIMmessageByMid(reqMsg.GetId())
		if im == nil {
			return easygo.NewFailMsg("消息Id错误")
		} else if im.GetStatus() >= for_game.WAITER_MESSAGE_END {
			return easygo.NewFailMsg("本次服务已结束")
		}
		reqMsg.Cnew = easygo.NewInt32(im.GetCnew() + 1)
	}

	content := reqMsg.GetContent()[0]
	if content.GetMtype() != for_game.WAITER_MSG_TYPE_S {
		return easygo.NewFailMsg("消息发送类型错误")
	}

	if content.Ctype == nil {
		return easygo.NewFailMsg("消息类型不能为空")
	}

	reqMsg.WaiterId = easygo.NewInt64(user.GetId())
	reqMsg.WaiterName = easygo.NewString(user.GetRealName())

	SendIMmessage(reqMsg)
	return easygo.EmptyMsg
}

//客服发送结束消息给玩家
func (self *cls4) RpcWaiterOverMsgToPlayer(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.IMmessage) easygo.IMessage {
	if user.GetRole() != 2 {
		return easygo.NewFailMsg("只有客服账号才能发送消息")
	}

	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		return easygo.NewFailMsg("消息Id错误")
	} else {
		im := for_game.QueryIMmessageByMid(reqMsg.GetId())
		if im == nil {
			return easygo.NewFailMsg("消息不存在")
		} else if im.GetStatus() == 1 {
			reqMsg.PlayerId = easygo.NewInt64(im.GetPlayerId())
			reqMsg.Status = easygo.NewInt32(2)
			reqMsg.UpdateTime = easygo.NewInt64(util.GetMilliTime())
			for_game.OverWaiterMessage(reqMsg.GetId())        //结束消息
			for_game.UpdateRedisWaiterCount(user.GetId(), -1) //减掉客服连接数
			SendToPlayer(im.GetPlayerId(), "RpcEndMessageToPlayer", reqMsg)
		}
	}

	return reqMsg
}

//客服管理列表
func (self *cls4) RpcWaiterPerformanceList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	manage := &share_message.Manager{}
	if reqMsg.GetType() == 2 {
		manage = QueryManageByName(reqMsg.GetKeyword())
		reqMsg.Keyword = easygo.NewString(easygo.AnytoA(manage.GetId()))
	}

	list, count := QueryWaiterPerformanceList(reqMsg)
	var ids []int64
	for _, v := range list {
		if for_game.IsContains(v.GetWaiterId(), ids) == -1 {
			ids = append(ids, v.GetWaiterId())
		}
	}

	mapUser := make(map[USER_ID]*share_message.Manager)
	managers := GetUserByIds(ids)
	for _, manager := range managers {
		mapUser[manager.GetId()] = manager
	}

	var lis []*share_message.WaiterPerformance
	for _, i := range list {
		i.WaiterName = easygo.NewString(mapUser[i.GetWaiterId()].GetRealName())
		i.WaiterRole = easygo.NewInt32(mapUser[i.GetWaiterId()].GetRoleType())
		lis = append(lis, i)
	}

	msg := &brower_backstage.WaiterPerformanceResponse{
		List:      lis,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//客服绩效查询
func (self *cls4) RpcWaiterPerformance(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	if reqMsg.GetId64() == 0 || reqMsg.Id64 == nil {
		return easygo.NewFailMsg("Id不能为空")
	}

	msg := for_game.QueryWaiterPerformance(reqMsg.GetId64())
	if msg == nil {
		waiter := QueryManageByID(reqMsg.GetId64())
		if waiter == nil {
			return easygo.NewFailMsg("Id错误")
		}

		msg = &share_message.WaiterPerformance{
			WaiterId:   waiter.Id,
			WaiterName: easygo.NewString(waiter.GetRealName()),
			WaiterRole: easygo.NewInt32(waiter.GetRoleType()),
			ConNum:     easygo.NewInt32(0),
			GradeNum:   easygo.NewInt32(0),
			SumGrade:   easygo.NewInt32(0),
		}
	}

	return msg
}

//客服聊天记录列表
func (self *cls4) RpcWaiterChatLogList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := QueryWaiterChatLogList(reqMsg)
	msg := &brower_backstage.IMmessageResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//客服常见问题列表
func (self *cls4) RpcWaiterFAQList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := QueryWaiterFAQList(reqMsg)
	msg := &brower_backstage.WaiterFAQResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//客服常见问题修改
func (self *cls4) RpcEditWaiterFAQ(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.WaiterFAQ) easygo.IMessage {
	if reqMsg.Title == nil {
		return easygo.NewFailMsg("标题不能为空")
	}

	if reqMsg.Content == nil {
		return easygo.NewFailMsg("回答不能为空")
	}

	msg := fmt.Sprintf("修改常见问题:%d", reqMsg.GetId())
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WAITER_FAQ))
		reqMsg.CreateTime = easygo.NewInt64(easygo.NowTimestamp())
		msg = fmt.Sprintf("添加常见问题:%d", reqMsg.GetId())
	}
	reqMsg.UpdateTime = easygo.NewInt64(easygo.NowTimestamp())
	EditWaiterFAQ(reqMsg)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WAITER_MANAGE, msg)

	return easygo.EmptyMsg
}

//删除客服常见问题
func (self *cls4) RpcDelWaiterFAQ(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds32()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}
	DelWaiterFAQ(idList)

	var ids string
	idsarr := reqMsg.GetIds32()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}
	msg := fmt.Sprintf("批量删除客服常见问题: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WAITER_MANAGE, msg)

	return easygo.EmptyMsg
}

//客服常用语列表
func (self *cls4) RpcWaiterFastReply(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := QueryWaiterFastReply(reqMsg)
	msg := &brower_backstage.WaiterFastReplyResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//客服常用语列表不分页
func (self *cls4) RpcWaiterFastReplyNopage(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	list, count := QueryWaiterFastReplyNopage()
	msg := &brower_backstage.WaiterFastReplyResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//客服常用语列表
func (self *cls4) RpcEditWaiterFastReply(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.WaiterFastReply) easygo.IMessage {
	if reqMsg.Content == nil {
		return easygo.NewFailMsg("内容不能为空")
	}

	msg := fmt.Sprintf("修改常用语:%d", reqMsg.GetId())
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_WAITER_FASTREPLY))
		reqMsg.CreateTime = easygo.NewInt64(easygo.NowTimestamp())
		msg = fmt.Sprintf("添加常用语:%d", reqMsg.GetId())
	}
	reqMsg.UpdateTime = easygo.NewInt64(easygo.NowTimestamp())
	EditWaiterFastReply(reqMsg)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WAITER_MANAGE, msg)

	return easygo.EmptyMsg
}

//删除客服常用语
func (self *cls4) RpcDelWaiterFastReply(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds32()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}
	DelWaiterFastReply(idList)

	var ids string
	idsarr := reqMsg.GetIds32()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}
	msg := fmt.Sprintf("批量删除客服常用语: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.WAITER_MANAGE, msg)

	return easygo.EmptyMsg
}

// 客服设置接待状态
func (self *cls4) RpcWaiterReception(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	if reqMsg.GetId64() == 0 || reqMsg.Id64 == nil {
		return easygo.NewFailMsg("Id不能为空")
	}
	waiter := for_game.GetRedisWaiter(reqMsg.GetId64())
	if waiter == nil {
		return easygo.NewFailMsg("客服不在线")
	} else if waiter.Status == 1 {
		for_game.SetRedisWaiterStatus(reqMsg.GetId64(), 0)
	}

	return easygo.EmptyMsg
}

// 客服设置休息状态
func (self *cls4) RpcWaiterRest(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	if reqMsg.GetId64() == 0 || reqMsg.Id64 == nil {
		return easygo.NewFailMsg("Id不能为空")
	}
	waiter := for_game.GetRedisWaiter(reqMsg.GetId64())
	if waiter == nil {
		return easygo.NewFailMsg("客服不在线")
	} else if waiter.Status == 0 {
		for_game.SetRedisWaiterStatus(reqMsg.GetId64(), 1)
	}

	if GetActiveIMmessageCount(user) > 0 {
		return easygo.NewFailMsg("您有尚未结束的对话")
	}

	return easygo.EmptyMsg
}
