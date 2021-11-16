// 管理后台为[浏览器]提供的服务
//运营渠道管理

package backstage

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"time"
)

//渠道列表
func (self *cls4) RpGetChannelList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	list := for_game.GetChannelListNopage()

	return &brower_backstage.KeyValueResponse{
		List: list,
	}
}

//查询运营渠道列表
func (self *cls4) RpcOperationChannelList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetOperationList(reqMsg)

	return &brower_backstage.OperationListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//修改运营渠道
func (self *cls4) RpcEditOperationChannel(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.OperationChannel) easygo.IMessage {
	msg := "修改运营渠道:"
	if reqMsg.Id == nil && reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_OPERATION_CHANNEL))
		reqMsg.CreateTime = easygo.NewInt64(time.Now().Unix())
		ban := []string{}
		dpSet := &share_message.DownPage{
			ModId:   easygo.NewInt32(1),   //模板ID：1默认模板，
			Icon:    easygo.NewString(""), //顶部图片url
			Banner:  ban,                  //轮播图
			BtnText: easygo.NewString(""), //下载按钮文字
			Floot:   easygo.NewString(""), //底部文字：粤ICP备17162070号-2
		}
		reqMsg.DpSet = dpSet
		msg = "添加运营渠道:"

		EditOperationChannelUse(reqMsg)
	} else {
		reqMsg.UpdateTime = easygo.NewInt64(time.Now().Unix())
		channel := for_game.QueryOperationByNo(reqMsg.GetChannelNo())
		if channel.GetChannelNo() == reqMsg.GetChannelNo() && channel.GetCooperation() == reqMsg.GetCooperation() && channel.GetPrice() == reqMsg.GetPrice() && reqMsg.GetStatus() == channel.GetStatus() && channel.GetRate() == reqMsg.GetRate() {
			EditOperationChannelUse(reqMsg)
		}
	}

	EditOperationChannel(reqMsg)

	msg = msg + easygo.IntToString(int(reqMsg.GetId()))
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.OPERATION_MANAGE, msg)

	return easygo.EmptyMsg
}
