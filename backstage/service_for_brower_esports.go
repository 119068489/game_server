// 管理后台为[浏览器]提供的服务

package backstage

//电竞管理
import (
	"fmt"
	"game_server/e-sports/sport_common_dal"
	"game_server/e-sports/sport_crawl"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"strconv"
	"strings"
	"time"

	"github.com/akqp2019/mgo"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

const (
	//注单操作
	BET_SLIP_OPERATE_1 string = "1" //无效
	BET_SLIP_OPERATE_2 string = "2" //违规
)

//启服加载电竞定时任务
func TimeSendEsports() {
	logs.Info("加载定时电竞新闻发布任务")
	lin, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_NEWS, bson.M{"Status": for_game.ESPORTS_NEWS_STATUS_0}, 0, 0)
	for _, item := range lin {
		AddSendEsportsTime(ES_NEWS, item)
	}

	logs.Info("加载定时电竞视频发布任务")
	liv, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_VIDEO, bson.M{"Status": for_game.ESPORTS_NEWS_STATUS_0}, 0, 0)
	for _, item := range liv {
		AddSendEsportsTime(ES_VIDEO, item)
	}

	logs.Info("加载定时电竞系统消息发布任务")
	lis, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_SYS_MSG, bson.M{"Status": for_game.ESPORTS_STATUS_0, "EffectiveType": ES_SEND_TYPE_FUTURE}, 0, 0)
	for _, item := range lis {
		AddSendEsportsTime(ES_SYSMSG, item)
	}

	logs.Info("加载定时电竞活动关闭任务")
	one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_LUCKY_ACTIVITY, bson.M{"Types": 2})
	if one != nil {
		AddActivityCloseTime(one)
	}

}

//添加定时电竞充值赠送活动
func AddActivityCloseTime(item interface{}) {
	one := &share_message.Activity{}
	for_game.StructToOtherStruct(item, one)
	if EsportsActCloseTimeMgr.GetTimerById(one.GetId()) != nil {
		EsportsActCloseTimeMgr.DelTimerList(one.GetId())
	}
	triggerTime := time.Duration(one.GetEndTime()-time.Now().Unix()) * time.Second
	if triggerTime >= 0 {
		timer := easygo.AfterFunc(triggerTime, func() {
			UpdateEsportsActStatus(for_game.MONGODB_NINGMENG, for_game.TABLE_LUCKY_ACTIVITY, one.GetId(), 1)
		})
		SendEsportsNewsTimeMgr.AddTimerList(one.GetId(), timer)
	} else {
		UpdateEsportsActStatus(for_game.MONGODB_NINGMENG, for_game.TABLE_LUCKY_ACTIVITY, one.GetId(), 1)
	}
}

//修改活动状态
func UpdateEsportsActStatus(db, table string, id int64, status int32) {
	queryBson := bson.M{"_id": id}
	updateBson := bson.M{"$set": bson.M{"Status": status}}
	for_game.FindAndModify(db, table, queryBson, updateBson, false)
}

//添加发送新闻视频任务定时器 t=1-新闻 2-视频 3-系统消息
func AddSendEsportsTime(t int, item interface{}) {
	switch t {
	case ES_NEWS:
		one := &share_message.TableESPortsRealTimeInfo{}
		for_game.StructToOtherStruct(item, one)
		if SendEsportsNewsTimeMgr.GetTimerById(one.GetId()) != nil {
			SendEsportsNewsTimeMgr.DelTimerList(one.GetId())
		}
		triggerTime := time.Duration(one.GetBeginEffectiveTime()-time.Now().Unix()) * time.Second
		if triggerTime >= 0 {
			timer := easygo.AfterFunc(triggerTime, func() {
				UpdateEsportsStatus(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_NEWS, one.GetId(), for_game.ESPORTS_NEWS_STATUS_1)
			})
			SendEsportsNewsTimeMgr.AddTimerList(one.GetId(), timer)
		} else {
			UpdateEsportsStatus(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_NEWS, one.GetId(), for_game.ESPORTS_NEWS_STATUS_4)
		}
	case ES_VIDEO:
		one := &share_message.TableESPortsVideoInfo{}
		for_game.StructToOtherStruct(item, one)
		if SendEsportsVideoTimeMgr.GetTimerById(one.GetId()) != nil {
			SendEsportsVideoTimeMgr.DelTimerList(one.GetId())
		}
		triggerTime := time.Duration(one.GetBeginEffectiveTime()-time.Now().Unix()) * time.Second
		if triggerTime >= 0 {
			timer := easygo.AfterFunc(triggerTime, func() {
				UpdateEsportsStatus(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_VIDEO, one.GetId(), for_game.ESPORTS_NEWS_STATUS_1)
			})
			SendEsportsVideoTimeMgr.AddTimerList(one.GetId(), timer)
		} else {
			UpdateEsportsStatus(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_VIDEO, one.GetId(), for_game.ESPORTS_NEWS_STATUS_4)
		}
	case ES_SYSMSG:
		one := &share_message.TableESPortsSysMsg{}
		for_game.StructToOtherStruct(item, one)
		if SendEsportsSysMsgTimeMgr.GetTimerById(one.GetId()) != nil {
			SendEsportsSysMsgTimeMgr.DelTimerList(one.GetId())
		}
		triggerTime := time.Duration(one.GetEffectiveTime()-time.Now().Unix()) * time.Second
		if triggerTime >= 0 {
			timer := easygo.AfterFunc(triggerTime, func() {
				UpdateEsportsStatus(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_SYS_MSG, one.GetId(), for_game.ESPORTS_STATUS_1)
				ChooseOneHall(0, "RpcSendSportSysNoticeToHall", one)
			})
			SendEsportsSysMsgTimeMgr.AddTimerList(one.GetId(), timer)
		} else {
			UpdateEsportsStatus(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_SYS_MSG, one.GetId(), for_game.ESPORTS_STATUS_2)
		}
	}
}

//修改状态
func UpdateEsportsStatus(db, table string, id int64, status int32) {
	queryBson := bson.M{"_id": id}
	updateBson := bson.M{"$set": bson.M{"Status": status}}
	for_game.FindAndModify(db, table, queryBson, updateBson, false)
}

//爬虫进度
func (s *cls4) RpcCrawlJobList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-Time"}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_CRAWL_JOB, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.TableCrawlJob
	for _, li := range lis {
		one := &share_message.TableCrawlJob{}
		for_game.StructToOtherStruct(li, one)
		o := one.GetId()
		game := ""
	reSwitch:
		switch o {
		case sport_crawl.ARTICLE_FNSCORE:
			one.Name = easygo.NewString("蜂鸟电竞资讯")
		case sport_crawl.ARTICLE_TVBCP:
			one.Name = easygo.NewString("鲨鱼比分资讯")
		case sport_crawl.VIDEO_WANPLUS:
			one.Name = easygo.NewString("玩加电竞视频" + game)
		case sport_crawl.VIDEO_CHAOFAN:
			one.Name = easygo.NewString("超凡电竞视频" + game)
		case sport_crawl.NEWS_CHAOFAN:
			one.Name = easygo.NewString("超凡电竞资讯" + game)
		case sport_crawl.NEWS_YXRB:
			one.Name = easygo.NewString("游戏日报资讯")
		case sport_crawl.NEWS_QQ:
			one.Name = easygo.NewString("腾讯网资讯" + game)
		case sport_crawl.NEWS_SINA:
			one.Name = easygo.NewString("新浪电竞资讯")
		default:
			i := strings.LastIndex(o, "_")
			if i == -1 {
				continue
			} else {
				game = o[i+1:]
				o = o[:i]
				goto reSwitch
			}
		}
		list = append(list, one)
	}

	msg := &brower_backstage.CrawlJobResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//手动爬取数据
func (s *cls4) RpcCrawlPull(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	findBson := bson.M{}
	if reqMsg.IdStr != nil && reqMsg.GetIdStr() != "" {
		findBson["_id"] = reqMsg.GetIdStr()
	}

	lis, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_CRAWL_JOB, findBson, 0, 0)
	for _, li := range lis {
		easygo.Spawn(func() {
			one := &share_message.TableCrawlJob{}
			for_game.StructToOtherStruct(li, one)
			o := one.GetId()
			game := ""
		reSwitch:
			switch o {
			case sport_crawl.ARTICLE_FNSCORE:
				one.Name = easygo.NewString("蜂鸟电竞资讯")
			case sport_crawl.ARTICLE_TVBCP:
				one.Name = easygo.NewString("鲨鱼比分资讯")
			case sport_crawl.VIDEO_WANPLUS:
				one.Name = easygo.NewString("玩加电竞视频" + game)
			case sport_crawl.VIDEO_CHAOFAN:
				one.Name = easygo.NewString("超凡电竞视频" + game)
			case sport_crawl.NEWS_CHAOFAN:
				one.Name = easygo.NewString("超凡电竞资讯" + game)
			case sport_crawl.NEWS_YXRB:
				one.Name = easygo.NewString("游戏日报资讯")
			case sport_crawl.NEWS_QQ:
				one.Name = easygo.NewString("腾讯网资讯" + game)
			case sport_crawl.NEWS_SINA:
				one.Name = easygo.NewString("新浪电竞资讯")
			default:
				i := strings.LastIndex(o, "_")
				if i == -1 {
					return
				} else {
					game = o[i+1:]
					o = o[:i]
					goto reSwitch
				}
			}

			_, err := SendMsgRandToServerNew(for_game.SERVER_TYPE_SPORT_CRAWL, "RpcCrawlPullForBackstage", one)
			if err == nil {
				ep.RpcCrawlPush(one)
			}
		})
	}
	return easygo.EmptyMsg
}

//查询app内游戏标签
func (s *cls4) RpcGetAppLabel(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{"LabelType": 3} //查询游戏标签
	sort := []string{"-Weight"}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		key := bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "im"}}
		findBson["Title"] = key
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_LABEL, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.TableESPortsLabel
	for _, li := range lis {
		one := &share_message.TableESPortsLabel{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.SysLabelResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//修改发布app内游戏标签
func (s *cls4) RpcSaveAppLabel(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.TableESPortsLabel) easygo.IMessage {
	msg := fmt.Sprintf("修改游戏标签:%d", reqMsg.GetId())
	if reqMsg.LabelId == nil || reqMsg.GetLabelId() == 0 {
		return easygo.NewFailMsg("游戏ID不能为空")
	}

	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_LABEL))
		msg = fmt.Sprintf("添加游戏标签:%d", reqMsg.GetId())
	}

	if reqMsg.IconUrl == nil || reqMsg.GetIconUrl() == "" {
		return easygo.NewFailMsg("游戏图标不能为空")
	}

	queryBson := bson.M{"LabelId": reqMsg.GetLabelId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_LABEL, queryBson, updateBson, true)
	for_game.SetRedisGameLabel()
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)
	return easygo.EmptyMsg
}

//查询自定义标签
func (s *cls4) RpcGetSysLabel(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{"LabelType": 2} //查询系统标签
	sort := []string{"-Weight"}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		key := bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "im"}}
		findBson["Title"] = key
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_LABEL, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.TableESPortsLabel
	for _, li := range lis {
		one := &share_message.TableESPortsLabel{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}

	msg := &brower_backstage.SysLabelResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//修改发布自定义标签
func (s *cls4) RpcSaveSysLabel(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.TableESPortsLabel) easygo.IMessage {
	msg := fmt.Sprintf("修改电竞自定义标签:%d", reqMsg.GetId())
	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_LABEL))
		reqMsg.LabelType = easygo.NewInt32(2)
		msg = fmt.Sprintf("添加电竞自定义标签:%d", reqMsg.GetId())
	}

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_LABEL, queryBson, updateBson, true)

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)
	return easygo.EmptyMsg
}

//引导项列表
func (s *cls4) RpcGetCarouselList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-Weight"}
	if reqMsg.Status != nil && reqMsg.GetStatus() != 0 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.Type != nil && reqMsg.GetType() != 0 {
		menuId := 0
		switch reqMsg.GetType() {
		case 1:
			menuId = for_game.ESPORTMENU_REALTIME
		case 2:
			menuId = for_game.ESPORTMENU_RECREATION
		case 3:
			menuId = for_game.ESPORTMENU_SHOP
		}
		findBson["MenuId"] = menuId
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_CAROUSEL, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.TableESPortsCarousel
	for _, li := range lis {
		one := &share_message.TableESPortsCarousel{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}

	msg := &brower_backstage.CarouselResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//修改引导项
func (s *cls4) RpcSaveCarousel(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.TableESPortsCarousel) easygo.IMessage {
	msg := fmt.Sprintf("修改引导项:%d", reqMsg.GetId())
	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_CAROUSEL))
		msg = fmt.Sprintf("添加引导项:%d", reqMsg.GetId())
	}

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_CAROUSEL, queryBson, updateBson, true)
	for_game.SetRedisExChangeBanner()
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)
	return easygo.EmptyMsg
}

//获取新闻资源列表
func (s *cls4) RpcNewsSource(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{"Status": bson.M{"$ne": for_game.ESPORTS_NEWS_STATUS_2}}
	sort := []string{"-CreateTime"}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		key := bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "im"}}
		switch reqMsg.GetType() {
		case 1:
			findBson["Title"] = key
		case 2:
			findBson["Content"] = key
		default:
			return easygo.NewFailMsg("搜索类型错误,1-标题,2-内容")
		}
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["AppLabelID"] = reqMsg.GetListType()
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_NEWS_SOURCE, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.TableESPortsRealTimeInfo
	for _, li := range lis {
		one := &share_message.TableESPortsRealTimeInfo{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.NewsSourceResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//获取视频资源列表
func (s *cls4) RpcVideoSource(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{"Status": bson.M{"$ne": for_game.ESPORTS_NEWS_STATUS_2}}
	sort := []string{"-CreateTime"}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		key := bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "im"}}
		switch reqMsg.GetType() {
		case 1:
			findBson["Title"] = key
		case 2:
			findBson["Content"] = key
		default:
			return easygo.NewFailMsg("搜索类型错误,1-标题,2-内容")
		}
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["AppLabelID"] = reqMsg.GetListType()
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_VIDEO_SOURCE, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.TableESPortsVideoInfo
	for _, li := range lis {
		one := &share_message.TableESPortsVideoInfo{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.VideoSourceResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//获取新闻资源列表
func (s *cls4) RpcNewsList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{}
	if reqMsg.GetSort() != "" {
		sort = append(sort, reqMsg.GetSort())
	} else {
		sort = append(sort, "-BeginEffectiveTime")
	}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		player := QueryPlayerbyPhone(reqMsg.GetKeyword())
		if player != nil {
			findBson["AuthorPlayerId"] = player.GetPlayerId()
		} else {
			findBson["AuthorAccount"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["AppLabelID"] = reqMsg.GetListType()
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() < 1000 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_NEWS, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.TableESPortsRealTimeInfo
	for _, li := range lis {
		one := &share_message.TableESPortsRealTimeInfo{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.NewsSourceResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//发布修改新闻资讯
func (s *cls4) RpcSaveNews(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.TableESPortsRealTimeInfo) easygo.IMessage {
	msg := fmt.Sprintf("修改新闻资讯:%d", reqMsg.GetId())
	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_NEWS))
		msg = fmt.Sprintf("添加新闻资讯:%d", reqMsg.GetId())
		if reqMsg.BeginEffectiveTime == nil {
			reqMsg.BeginEffectiveTime = easygo.NewInt64(easygo.NowTimestamp()) //资讯资源管理发布时间不随修改时间变化
		}
	}

	reqMsg.Content = easygo.NewString(QQbucket.ReplaceImg(reqMsg.GetContent())) //将文章中的图片转存到存储桶

	if reqMsg.GetEffectiveType() == ES_SEND_TYPE_FUTURE {
		AddSendEsportsTime(ES_NEWS, reqMsg)
	} else {
		reqMsg.Status = easygo.NewInt32(for_game.ESPORTS_NEWS_STATUS_1)
	}
	// newContent := QQbucket.ReplaceImg(reqMsg.GetContent())
	// reqMsg.Content = easygo.NewString(newContent)
	reqMsg.UpdateTime = easygo.NewInt64(easygo.NowTimestamp())
	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_NEWS, queryBson, updateBson, true)

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)
	return easygo.EmptyMsg
}

//删除新闻资讯资源
func (s *cls4) RpcDelNewsSource(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	queryBson := bson.M{"_id": reqMsg.GetId64()}
	updateBson := bson.M{"$set": bson.M{"Status": for_game.ESPORTS_NEWS_STATUS_2}}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_NEWS_SOURCE, queryBson, updateBson, true)

	return easygo.EmptyMsg
}

//删除新闻资讯
func (s *cls4) RpcDelNews(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}

	findBson := bson.M{"_id": bson.M{"$in": idList}}
	updateBson := bson.M{"$set": bson.M{"Status": for_game.ESPORTS_NEWS_STATUS_2}}
	for_game.UpdateAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_NEWS, findBson, updateBson)

	var ids string
	count := len(idList)
	for i, t := range idList {
		ids += easygo.AnytoA(t)
		if i < count-1 {
			ids += ","

		}

		if SendEsportsNewsTimeMgr.GetTimerById(t) != nil {
			SendEsportsNewsTimeMgr.DelTimerList(t)
		}
	}

	msg := fmt.Sprintf("批量删除新闻资讯: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)

	return easygo.EmptyMsg
}

//获取视频列表
func (s *cls4) RpcVideoList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{}
	if reqMsg.GetSort() != "" {
		sort = append(sort, reqMsg.GetSort())
	} else {
		sort = append(sort, "-BeginEffectiveTime")
	}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		player := QueryPlayerbyPhone(reqMsg.GetKeyword())
		if player != nil {
			findBson["AuthorPlayerId"] = player.GetPlayerId()
		} else {
			findBson["AuthorAccount"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["AppLabelID"] = reqMsg.GetListType()
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() < 1000 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.GetType() != 0 && reqMsg.Type != nil {
		findBson["VideoType"] = reqMsg.GetType()
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}

	if reqMsg.DownType != nil && reqMsg.GetDownType() != 0 {
		findBson["AuthorPlayerType"] = reqMsg.GetDownType()
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_VIDEO, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.TableESPortsVideoInfo
	for _, li := range lis {
		one := &share_message.TableESPortsVideoInfo{}
		for_game.StructToOtherStruct(li, one)
		if one.GetVideoType() == 2 {
			game := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME, bson.M{"_id": one.GetUniqueGameId()})
			if game != nil {
				if game.(bson.M)["MatchName"] != nil {
					one.UniqueGameName = easygo.NewString(easygo.AnytoA(game.(bson.M)["MatchName"])) ////比赛名称(赛事名+赛事阶段名+赛事阶段id)
				}
			}
			player := QueryPlayerbyId(one.GetAuthorPlayerId())
			if player != nil {
				one.AuthorPlayerType = easygo.NewInt32(player.GetTypes())
				one.AuthorAccount = easygo.NewString(player.GetAccount())
				one.AuthorNickName = easygo.NewString(player.GetNickName())
			}
		}
		list = append(list, one)
	}
	msg := &brower_backstage.VideoSourceResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//发布修改视频
func (s *cls4) RpcSaveVideo(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.TableESPortsVideoInfo) easygo.IMessage {
	msg := fmt.Sprintf("修改视频:%d", reqMsg.GetId())
	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_VIDEO))
		msg = fmt.Sprintf("添加视频:%d", reqMsg.GetId())
		if reqMsg.BeginEffectiveTime == nil {
			reqMsg.BeginEffectiveTime = easygo.NewInt64(easygo.NowTimestamp()) //视频资源管理发布时间不随修改时间变化
		}
	}

	if reqMsg.AuthorPlayerType == nil || reqMsg.GetAuthorPlayerType() == 0 {
		player := QueryPlayerbyId(reqMsg.GetAuthorPlayerId())
		if player != nil {
			reqMsg.AuthorPlayerType = easygo.NewInt32(player.GetTypes())
		}
	}

	if reqMsg.GetEffectiveType() == ES_SEND_TYPE_FUTURE {
		AddSendEsportsTime(ES_VIDEO, reqMsg)
	} else {
		reqMsg.UpdateTime = easygo.NewInt64(easygo.NowTimestamp())
		reqMsg.Status = easygo.NewInt32(for_game.ESPORTS_NEWS_STATUS_1)
	}
	//如果是放映厅
	if reqMsg.GetVideoType() == 2 && reqMsg.GetUniqueGameId() > 0 {
		gameInfo := sport_common_dal.GetMatchInfo(reqMsg.GetUniqueGameId())
		reqMsg.UniqueGameInfo = gameInfo
	}

	reqMsg.Operator = easygo.NewString(user.GetAccount())
	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_VIDEO, queryBson, updateBson, true)

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)
	return easygo.EmptyMsg
}

//删除视频资源
func (s *cls4) RpcDelVideoSource(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	queryBson := bson.M{"_id": reqMsg.GetId64()}
	updateBson := bson.M{"$set": bson.M{"Status": for_game.ESPORTS_NEWS_STATUS_2}}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_VIDEO_SOURCE, queryBson, updateBson, true)

	return easygo.EmptyMsg
}

//删除视频
func (s *cls4) RpcDelVideo(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}

	findBson := bson.M{"_id": bson.M{"$in": idList}}
	updateBson := bson.M{"$set": bson.M{"Status": for_game.ESPORTS_NEWS_STATUS_2}}
	for_game.UpdateAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_VIDEO, findBson, updateBson)

	var ids string
	count := len(idList)
	for i, t := range idList {
		ids += easygo.AnytoA(t)
		if i < count-1 {
			ids += ","
		}

		if SendEsportsVideoTimeMgr.GetTimerById(t) != nil {
			SendEsportsVideoTimeMgr.DelTimerList(t)
		}
	}

	msg := fmt.Sprintf("批量删除视频: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)

	return easygo.EmptyMsg
}

//审核直播(放映厅)
func (s *cls4) RpcChekVideo(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	if reqMsg.GetId64() == 0 || reqMsg.Id64 == nil {
		return easygo.NewFailMsg("Id64参数不能为空")
	}

	one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_VIDEO, bson.M{"_id": reqMsg.GetId64(), "Status": for_game.ESPORTS_NEWS_STATUS_0})
	if one == nil {
		return easygo.NewFailMsg("未找到需要审核的直播间")
	}

	status := for_game.ESPORTS_NEWS_STATUS_0
	switch reqMsg.GetId32() {
	case 1: //通过
		status = for_game.ESPORTS_NEWS_STATUS_1
	case 2: //拒绝
		status = for_game.ESPORTS_NEWS_STATUS_2
	default:
		return easygo.NewFailMsg("审核操作错误")
	}

	findBson := bson.M{"_id": reqMsg.GetId64()}
	upBson := bson.M{"$set": bson.M{"Status": status, "Operator": user.GetAccount()}}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_VIDEO, findBson, upBson, false)

	msg := fmt.Sprintf("审核直播间[%d]", reqMsg.GetId64())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)

	return easygo.EmptyMsg
}

//封禁直播(放映厅)
func (s *cls4) RpcBanVideo(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要禁用的项")
	}

	status := reqMsg.GetIds32()
	if len(status) == 0 {
		return easygo.NewFailMsg("状态参数不能为空")
	}

	if status[0] == 1 && (reqMsg.Note == nil || reqMsg.GetNote() == "") {
		return easygo.NewFailMsg("请先填写封禁备注")
	}
	msgf := "批量封禁直播: %s"
	updataBson := bson.M{"Operator": user.GetAccount()}
	switch status[0] {
	case 1:
		updataBson["Note"] = reqMsg.GetNote()
		updataBson["Status"] = for_game.ESPORTS_NEWS_STATUS_3
	case 2:
		msgf = "批量解禁直播: %s"
		updataBson["Note"] = ""
		updataBson["Status"] = for_game.ESPORTS_NEWS_STATUS_1
	default:
		return easygo.NewFailMsg("状态参数错误")
	}

	findBson := bson.M{"_id": bson.M{"$in": idList}}
	updateBson := bson.M{"$set": updataBson}
	for_game.UpdateAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_VIDEO, findBson, updateBson)

	var ids string
	count := len(idList)
	for i, t := range idList {
		ids += easygo.AnytoA(t)
		if i < count-1 {
			ids += ","
		}

		if SendEsportsVideoTimeMgr.GetTimerById(t) != nil {
			SendEsportsVideoTimeMgr.DelTimerList(t)
		}

		req := &client_hall.ESportDataStatusInfo{
			MenuId: easygo.NewInt32(for_game.ESPORTMENU_LIVE),
			DataId: easygo.NewInt64(t),
			Status: easygo.NewInt32(status[0]),
		}
		SendMsgRandToServerNew(for_game.SERVER_TYPE_SPORT_APPLY, "RpcESportDataStatusInfo", req)
	}

	msg := fmt.Sprintf(msgf, ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)

	return easygo.EmptyMsg
}

//获取赛事列表
func (s *cls4) RpcGetGameList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{}
	if reqMsg.Sort == nil && reqMsg.GetSort() == "" {
		sort = append(sort, "begin_time_int")
	} else {
		sort = append(sort, reqMsg.GetSort())
	}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetDownType() {
		case 1:
			findBson["match_name"] = reqMsg.GetKeyword()
		case 2:
			findBson["bo"] = reqMsg.GetKeyword()
		case 3:
			findBson["match_stage"] = reqMsg.GetKeyword()
		case 4:
			findBson["$or"] = []bson.M{{"team_a.name": reqMsg.GetKeyword()}, {"team_a.name_en": reqMsg.GetKeyword()}, {"team_b.name_en": reqMsg.GetKeyword()}, {"team_b.name": reqMsg.GetKeyword()}}
		case 5: //比赛ID
			findBson["_id"] = easygo.StringToInt64noErr(reqMsg.GetKeyword())
		}
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["app_label_id"] = reqMsg.GetListType()
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() != 0 {
		findBson["release_flag"] = reqMsg.GetStatus()
	}

	if reqMsg.Type != nil {
		if reqMsg.GetType() == 100 {
			findBson["game_status"] = bson.M{"$ne": "2"} //查询未开始和进行中的比赛
		} else {
			//比赛状态 0 未开始，1 进行中，2 已结束(api字段)(0和1的时候结合begin_time判断)
			status := easygo.AnytoA(reqMsg.GetType())
			switch status {
			case for_game.GAME_STATUS_0:
				findBson["game_status"] = bson.M{"$ne": "2"}
				findBson["begin_time_int"] = bson.M{"$gt": easygo.NowTimestamp()}
			case for_game.GAME_STATUS_1:
				findBson["game_status"] = bson.M{"$ne": "2"}
				findBson["begin_time_int"] = bson.M{"$lt": easygo.NowTimestamp()}
			case for_game.GAME_STATUS_2:
				findBson["game_status"] = status
			}
		}
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		if reqMsg.SrtType != nil && reqMsg.GetSrtType() == "1" {
			findBson["begin_time_int"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		} else {
			findBson["create_time"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.TableESPortsGame
	for _, li := range lis {
		one := &share_message.TableESPortsGame{}
		for_game.StructToOtherStruct(li, one)
		//手动赋予比赛状态值
		if one.GetGameStatus() != for_game.GAME_STATUS_2 {
			one.GameStatus = easygo.NewString(for_game.GetGameStatus(one.GetBeginTime(), one.GetGameStatus()))
		}

		list = append(list, one)
	}
	msg := &brower_backstage.GameListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return msg
}

//发布赛事
/*
func (self *cls4) RpcReleaseGame(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {

	if reqMsg.Id64 == nil || reqMsg.GetId64() == 0 {
		return easygo.NewFailMsg("参数错误Id64")
	}

	queryBson := bson.M{"_id": reqMsg.GetId64()}
	updateBson := bson.M{"$set": bson.M{"release_flag": for_game.GAME_RELEASE_FLAG_2, "update_time": easygo.NowTimestamp()}}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME, queryBson, updateBson, true)

	msg := fmt.Sprintf("发布赛事:%d", reqMsg.GetId64())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)
	return easygo.EmptyMsg
}*/

//赛事竞猜
func (s *cls4) RpcGetGameGuess(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.GameGuessRequest) easygo.IMessage {

	guessDetailObject1 := for_game.GetDBGuessDetail(reqMsg.GetLabelId(), reqMsg.GetGameId(), reqMsg.GetApiOrigin(), for_game.GAME_IS_MORN_ROLL_1)
	guessDetailObject2 := for_game.GetDBGuessDetail(reqMsg.GetLabelId(), reqMsg.GetGameId(), reqMsg.GetApiOrigin(), for_game.GAME_IS_MORN_ROLL_2)

	if nil == guessDetailObject1 {
		guessDetailObject1 = &share_message.GameGuessDetailObject{}
	}

	if nil == guessDetailObject2 {
		guessDetailObject2 = &share_message.GameGuessDetailObject{}
	}

	one := &share_message.TableESPortsGameDetail{}
	result := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_DETAIL, bson.M{"app_label_id": reqMsg.GetLabelId(), "game_id": reqMsg.GetGameId(), "api_origin": reqMsg.GetApiOrigin()})
	if result != nil {
		for_game.StructToOtherStruct(result, one)
	}

	return &brower_backstage.GameGuessResponse{
		List1Id:   easygo.NewInt64(guessDetailObject1.GetUniqueGameGuessId()),
		List1:     guessDetailObject1.GetGuessOddsNums(),
		List2Id:   easygo.NewInt64(guessDetailObject2.GetUniqueGameGuessId()),
		List2:     guessDetailObject2.GetGuessOddsNums(),
		LivePaths: one.GetLivePaths(),
	}
}

//查询比赛队伍信息
func (s *cls4) RpcGetGameTeamInfo(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.GameGuessRequest) easygo.IMessage {
	id, gameid, origin := reqMsg.GetLabelId(), reqMsg.GetGameId(), reqMsg.GetApiOrigin()
	if id == 0 || gameid == "" || origin == 0 {
		return easygo.NewFailMsg("参数错误")
	}

	one := &share_message.TableESPortsGameDetail{}
	result := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_DETAIL, bson.M{"app_label_id": id, "game_id": gameid, "api_origin": origin})
	if result != nil {
		for_game.StructToOtherStruct(result, one)
	}

	return &brower_backstage.GameTeamInfoResponse{
		TeamA: one.GetApiTeamAPlayers(),
		TeamB: one.GetApiTeamBPlayers(),
	}
}

//发布更新赛事竞猜状态
func (s *cls4) RpcEditGameGuess(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.EditGameGuessRequest) easygo.IMessage {
	logs.Debug("RpcEditGameGuess", reqMsg)
	if reqMsg.GetId() == 0 {
		return easygo.NewFailMsg("赛事ID错误")
	}

	if reqMsg.HistoryId != nil {
		oldone := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME, bson.M{"_id": reqMsg.GetId()})
		oldoneGame := &share_message.TableESPortsGame{}
		for_game.StructToOtherStruct(oldone, oldoneGame)

		queryBson := bson.M{"_id": reqMsg.GetId()}
		updateBson := bson.M{"$set": bson.M{"history_id": reqMsg.GetHistoryId()}}
		newone := for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME, queryBson, updateBson, true)
		newoneGame := &share_message.TableESPortsGame{}
		for_game.StructToOtherStruct(newone, newoneGame)

		//修改历史战绩id后,爬取更新历史战绩
		if newoneGame.GetHistoryId() != oldoneGame.GetHistoryId() {
			easygo.Spawn(sport_crawl.CrawlScoreHistoryDataRun, reqMsg.GetHistoryId())
		}
	}

	if reqMsg.GetOpt() == "add" {
		queryBson := bson.M{"_id": reqMsg.GetId()}
		updateBson := bson.M{"$set": bson.M{"release_flag": for_game.GAME_RELEASE_FLAG_2, "update_time": easygo.NowTimestamp(), "create_time": easygo.NowTimestamp()}}
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME, queryBson, updateBson, true)

		msg := fmt.Sprintf("发布赛事:%d", reqMsg.GetId())
		AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)
	} else {
		queryBson := bson.M{"_id": reqMsg.GetId()}
		updateBson := bson.M{"$set": bson.M{"update_time": easygo.NowTimestamp()}}
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME, queryBson, updateBson, true)
	}

	//判断比赛是否开奖 开奖了不需要设置
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME)
	defer closeFun()
	sportGameQuery := share_message.TableESPortsGame{}
	//通过条件查询
	errQuery := col.Find(bson.M{"_id": reqMsg.GetId()}).One(&sportGameQuery)

	if errQuery != nil && errQuery != mgo.ErrNotFound {
		logs.Error(errQuery)
		return easygo.NewFailMsg("redis设置时、查询比赛失败")
	}

	if errQuery == mgo.ErrNotFound {
		return easygo.NewFailMsg("redis设置未找到比赛")
	}

	guess := reqMsg.GetGuesses()
	if len(guess) == 0 {
		return easygo.EmptyMsg
	}
	for _, g := range guess {
		id := g.GetId()
		kv := g.GetGuess()
		if len(kv) == 0 {
			continue
		}

		flags := make(map[string]int32)
		showflags := make(map[string]int32)
		for _, v := range kv {
			if int32(v.GetValue()) == for_game.GAME_APP_GUESS_FLAG_2 {
				flags[v.GetName()] = for_game.GAME_APP_GUESS_FLAG_2

			} else {
				flags[v.GetName()] = for_game.GAME_APP_GUESS_FLAG_1

			}
			if int32(v.GetExtend()) == for_game.GAME_APP_GUESS_FLAG_2 {
				showflags[v.GetName()] = for_game.GAME_APP_GUESS_FLAG_2
			} else {
				showflags[v.GetName()] = for_game.GAME_APP_GUESS_FLAG_1
			}
		}

		gg := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_GUESS, bson.M{"_id": id})
		if gg == nil {
			return easygo.NewFailMsg("竞猜不存在")
		}
		one := &share_message.TableESPortsGameGuess{}
		for_game.StructToOtherStruct(gg, one)

		ggs := one.GetGuess()
		if len(ggs) == 0 {
			return easygo.NewFailMsg("竞猜项不存在")
		}

		for _, g := range ggs {
			if flags[g.GetBetId()] > 0 {
				g.AppGuessFlag = easygo.NewInt32(flags[g.GetBetId()])
				g.AppGuessViewFlag = easygo.NewInt32(showflags[g.GetBetId()])
			}
		}

		for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME_GUESS, bson.M{"_id": id}, bson.M{"$set": bson.M{"guess": ggs}}, false)
	}

	//设置redis=============开始===================
	if errQuery == nil {
		if for_game.GetGameStatus(sportGameQuery.GetBeginTime(), sportGameQuery.GetGameStatus()) == for_game.GAME_STATUS_2 {
			return easygo.NewFailMsg("比赛已结束")
		}

		//设置redis(把最新的historyId设置到redis)
		for_game.SetRedisGameDetailHead(sportGameQuery.GetId())
		//设置redis早盘信息
		for_game.SetRedisGuessMornDetail(reqMsg.GetId(), sportGameQuery.GetAppLabelId(), sportGameQuery.GetGameId(), sportGameQuery.GetApiOrigin())
		//设置redis滚盘信息
		for_game.SetRedisGuessRollDetail(reqMsg.GetId(), sportGameQuery.GetAppLabelId(), sportGameQuery.GetGameId(), sportGameQuery.GetApiOrigin())
	}
	//设置redis=============结束===================

	msg := fmt.Sprintf("修改赛事[%d]竞猜项状态:", reqMsg.GetId())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)
	return easygo.EmptyMsg
}

//查询比赛实时数据
func (s *cls4) RpcGetGameRealTimeData(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.GameGuessRequest) easygo.IMessage {
	if reqMsg.LabelId == nil {
		return easygo.NewFailMsg("LabelId不能为空")
	}
	if reqMsg.GameId == nil {
		return easygo.NewFailMsg("GameId不能为空")
	}
	if reqMsg.ApiOrigin == nil {
		return easygo.NewFailMsg("ApiOrigin不能为空")
	}
	labelId := reqMsg.GetLabelId()
	gameId := easygo.StringToInt64noErr(reqMsg.GetGameId())
	apiOrigin := reqMsg.GetApiOrigin()

	msg := &brower_backstage.GameRealTimeResponse{}
	switch labelId {
	case for_game.ESPORTS_LABEL_WZRY:
		result, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_WZRY_REAL_TIME_DATA, bson.M{"game_id": gameId, "api_origin": apiOrigin}, 0, 0)
		list := []*share_message.TableESPortsWZRYRealTimeData{}
		for _, d := range result {
			one := &share_message.TableESPortsWZRYRealTimeData{}
			for_game.StructToOtherStruct(d, one)
			// ret := &share_message.TableESPortsWZRYRealTimeData{
			// 	Id:               easygo.NewInt64(one.GetId()),
			// 	GameRound:        easygo.NewInt32(one.GetGameRound()),
			// 	FirstTower:       easygo.NewInt32(one.GetFirstTower()),
			// 	FirstSmallDragon: easygo.NewInt32(one.GetFirstSmallDragon()),
			// 	FirstFiveKill:    easygo.NewInt32(one.GetFirstFiveKill()),
			// 	FirstBigDragon:   easygo.NewInt32(one.GetFirstBigDragon()),
			// 	FirstTenKill:     easygo.NewInt32(one.GetFirstTenKill()),
			// 	TeamA: &share_message.ApiWZRYTeam{
			// 		TowerState:   easygo.NewInt32(one.TeamA.GetTowerState()),
			// 		Drakes:       easygo.NewInt32(one.TeamA.GetDrakes()),
			// 		NahsorBarons: easygo.NewInt32(one.TeamA.GetNahsorBarons()),
			// 	},
			// 	TeamB: &share_message.ApiWZRYTeam{
			// 		TowerState:   easygo.NewInt32(one.TeamB.GetTowerState()),
			// 		Drakes:       easygo.NewInt32(one.TeamB.GetDrakes()),
			// 		NahsorBarons: easygo.NewInt32(one.TeamB.GetNahsorBarons()),
			// 	},
			// 	PlayerAInfo: one.GetPlayerAInfo(),
			// 	PlayerBInfo: one.GetPlayerBInfo(),
			// }
			list = append(list, one)
		}
		msg.Wzry = list
	case for_game.ESPORTS_LABEL_LOL:
		result, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_LOL_REAL_TIME_DATA, bson.M{"game_id": gameId, "api_origin": apiOrigin}, 0, 0)
		list := []*share_message.TableESPortsLOLRealTimeData{}
		for _, d := range result {
			one := &share_message.TableESPortsLOLRealTimeData{}
			for_game.StructToOtherStruct(d, one)
			// ret := &share_message.TableESPortsLOLRealTimeData{
			// 	Id:               easygo.NewInt64(one.GetId()),
			// 	GameRound:        easygo.NewInt32(one.GetGameRound()),
			// 	FirstTower:       easygo.NewInt32(one.GetFirstTower()),
			// 	FirstSmallDragon: easygo.NewInt32(one.GetFirstSmallDragon()),
			// 	FirstFiveKill:    easygo.NewInt32(one.GetFirstFiveKill()),
			// 	FirstBigDragon:   easygo.NewInt32(one.GetFirstBigDragon()),
			// 	FirstTenKill:     easygo.NewInt32(one.GetFirstTenKill()),
			// 	TeamA: &share_message.ApiLOLTeam{
			// 		TowerState:   easygo.NewInt32(one.TeamA.GetTowerState()),
			// 		Drakes:       easygo.NewInt32(one.TeamA.GetDrakes()),
			// 		NahsorBarons: easygo.NewInt32(one.TeamA.GetNahsorBarons()),
			// 	},
			// 	TeamB: &share_message.ApiLOLTeam{
			// 		TowerState:   easygo.NewInt32(one.TeamB.GetTowerState()),
			// 		Drakes:       easygo.NewInt32(one.TeamB.GetDrakes()),
			// 		NahsorBarons: easygo.NewInt32(one.TeamB.GetNahsorBarons()),
			// 	},
			// 	PlayerAInfo: one.GetPlayerAInfo(),
			// 	PlayerBInfo: one.GetPlayerBInfo(),
			// }
			list = append(list, one)
		}
		msg.Lol = list
	}
	return msg
}

//修改比赛实时数据
func (s *cls4) RpcEditGameRealTimeData(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.EditGameRealTimeRequest) easygo.IMessage {
	id := reqMsg.GetId()
	realTimeObjects := reqMsg.GetRealTimeObject()
	if reqMsg.Id == nil || id == 0 {
		return easygo.NewFailMsg("Id参数不能为空")
	}

	if realTimeObjects == nil || len(realTimeObjects) <= 0 {
		return easygo.NewFailMsg("没有修改的内容、realTimeObjects参数错误")
	}

	//从redis取比赛的信息
	gameObject := for_game.GetRedisGameDetailHead(id)
	if nil == gameObject {
		s := fmt.Sprintf("缺少比赛id:%v的数据", id)
		return easygo.NewFailMsg(s)
	}
	appLabelId := gameObject.GetAppLabelId()
	gameIdStr := gameObject.GetGameId()
	apiOrigin := gameObject.GetApiOrigin()

	if appLabelId != for_game.ESPORTS_LABEL_LOL && appLabelId != for_game.ESPORTS_LABEL_WZRY {
		s := fmt.Sprintf("比赛id:%v、不是王者荣耀和LOL", id)
		return easygo.NewFailMsg(s)
	}

	gameId, err := strconv.ParseInt(gameIdStr, 10, 32)

	if err != nil {
		s := fmt.Sprintf("=======实时数据参数gameIdStr转换时错误==========:gameId=%v,err=%v", gameIdStr, err)
		logs.Error(s)
		return easygo.NewFailMsg("系统异常")
	}

	//LOL
	if appLabelId != for_game.ESPORTS_LABEL_LOL {
		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_LOL_REAL_TIME_DATA)
		defer closeFun()

		for _, value := range realTimeObjects {

			gameRound := value.GetGameRound()

			query := share_message.TableESPortsLOLRealTimeData{}
			//通过条件查询
			errQuery := col.Find(bson.M{"app_label_id": appLabelId,
				"game_id":    int32(gameId),
				"api_origin": apiOrigin,
				"game_round": gameRound}).One(&query)

			if errQuery != nil && errQuery != mgo.ErrNotFound {
				logs.Error(errQuery)
				s := fmt.Sprintf("======RpcEditGameRealTimeData修改游戏LOL实时数据表查询实时数据失败======查询条件为:app_label_id:%v,===game_id:%v,====api_origin:%v,====game_round:%v",
					appLabelId, gameId, apiOrigin, gameRound)
				logs.Error(s)
				return easygo.NewFailMsg("系统异常")
			}

			if errQuery == mgo.ErrNotFound {
				s := fmt.Sprintf("缺少第%v局数据", gameRound)
				return easygo.NewFailMsg(s)
			}

			//修改数据逻辑
			if errQuery == nil {
				//实时数据第一层数据
				query.FirstTower = easygo.NewInt32(value.GetFirstTower())
				query.FirstSmallDragon = easygo.NewInt32(value.GetFirstSmallDragon())
				query.FirstBigDragon = easygo.NewInt32(value.GetFirstBigDragon())
				query.FirstFiveKill = easygo.NewInt32(value.GetFirstFiveKill())
				query.FirstTenKill = easygo.NewInt32(value.GetFirstTenKill())
				playerAInfo := query.PlayerAInfo
				playerAInfoParam := value.TeamAPlayers
				//实时数据中A队员信息
				if playerAInfo != nil && playerAInfoParam != nil {
					for _, playValue := range playerAInfo {
						for _, playValueParam := range playerAInfoParam {
							if playValue.GetName() == playValueParam.GetName() {
								if playValueParam.GetPhoto() != "" {
									playValue.Photo = easygo.NewString(playValueParam.GetPhoto())
								}
								break
							}
						}
					}
				}
				playerBInfo := query.PlayerBInfo
				playerBInfoParam := value.TeamBPlayers
				//实时数据中B队员信息
				if playerBInfo != nil && playerBInfoParam != nil {
					for _, playValue := range playerBInfo {
						for _, playValueParam := range playerBInfoParam {
							if playValue.GetName() == playValueParam.GetName() {
								if playValueParam.GetPhoto() != "" {
									playValue.Photo = easygo.NewString(playValueParam.GetPhoto())
								}
								break
							}
						}
					}
				}
			}

			//修改数据库
			errUpd := col.Update(bson.M{"_id": query.GetId()},
				bson.M{"$set": query})

			if errUpd != nil {
				logs.Error(errUpd)
				s := fmt.Sprintf("========RpcEditGameRealTimeData LOL游戏实时数据表更新数据失败======更新条件WZRY唯一_id:%v",
					query.GetId())
				logs.Error(s)
				return easygo.NewFailMsg("系统异常")
			}

			//设置reids
			for_game.SetRedisGameRealTime(id, gameRound)
		}
	}

	//WZRY
	if appLabelId != for_game.ESPORTS_LABEL_WZRY {
		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_WZRY_REAL_TIME_DATA)
		defer closeFun()

		for _, value := range realTimeObjects {

			gameRound := value.GetGameRound()

			query := share_message.TableESPortsLOLRealTimeData{}
			//通过条件查询
			errQuery := col.Find(bson.M{"app_label_id": appLabelId,
				"game_id":    int32(gameId),
				"api_origin": apiOrigin,
				"game_round": gameRound}).One(&query)

			if errQuery != nil && errQuery != mgo.ErrNotFound {
				logs.Error(errQuery)
				s := fmt.Sprintf("======RpcEditGameRealTimeData修改游戏WZRY实时数据表查询实时数据失败======查询条件为:app_label_id:%v,===game_id:%v,====api_origin:%v,====game_round:%v",
					appLabelId, gameId, apiOrigin, gameRound)
				logs.Error(s)
				return easygo.NewFailMsg("系统异常")
			}

			if errQuery == mgo.ErrNotFound {
				s := fmt.Sprintf("缺少第%v局数据", gameRound)
				return easygo.NewFailMsg(s)
			}

			//修改数据逻辑
			if errQuery == nil {
				//实时数据第一层数据
				query.FirstTower = easygo.NewInt32(value.GetFirstTower())
				query.FirstSmallDragon = easygo.NewInt32(value.GetFirstSmallDragon())
				query.FirstBigDragon = easygo.NewInt32(value.GetFirstBigDragon())
				query.FirstFiveKill = easygo.NewInt32(value.GetFirstFiveKill())
				query.FirstTenKill = easygo.NewInt32(value.GetFirstTenKill())

				//WZRY需要设置击杀小龙数、击杀大龙数
				teamA := query.TeamA
				if teamA != nil {
					teamA.NahsorBarons = easygo.NewInt32(value.GetTeamANahsorBarons())
					teamA.Drakes = easygo.NewInt32(value.GetTeamADrakes())
				}

				teamB := query.TeamB
				if teamB != nil {
					teamB.NahsorBarons = easygo.NewInt32(value.GetTeamBNahsorBarons)
					teamB.Drakes = easygo.NewInt32(value.GetTeamBDrakes())
				}

				playerAInfo := query.PlayerAInfo
				playerAInfoParam := value.TeamAPlayers
				//实时数据中A队员信息
				if playerAInfo != nil && playerAInfoParam != nil {
					for _, playValue := range playerAInfo {
						for _, playValueParam := range playerAInfoParam {
							if playValue.GetName() == playValueParam.GetName() {
								if playValueParam.GetPhoto() != "" {
									playValue.Photo = easygo.NewString(playValueParam.GetPhoto())
								}
								break
							}
						}
					}
				}
				playerBInfo := query.PlayerBInfo
				playerBInfoParam := value.TeamBPlayers
				//实时数据中B队员信息
				if playerBInfo != nil && playerBInfoParam != nil {
					for _, playValue := range playerBInfo {
						for _, playValueParam := range playerBInfoParam {

							if playValue.GetName() == playValueParam.GetName() {
								if playValueParam.GetPhoto() != "" {
									playValue.Photo = easygo.NewString(playValueParam.GetPhoto())
								}
								break
							}
						}
					}
				}
			}

			//修改数据库
			errUpd := col.Update(bson.M{"_id": query.GetId()},
				bson.M{"$set": query})

			if errUpd != nil {
				logs.Error(errUpd)
				s := fmt.Sprintf("========RpcEditGameRealTimeData WZRY游戏实时数据表更新数据失败======更新条件WZRY唯一_id:%v",
					query.GetId())
				logs.Error(s)
				return easygo.NewFailMsg("系统异常")
			}

			//设置reids
			for_game.SetRedisGameRealTime(id, gameRound)
		}
	}

	//
	//roundLs := []int32{}
	//labelId := reqMsg.GetLabelId()
	//var tableName string
	//var data []interface{}
	//switch labelId {
	//case for_game.ESPORTS_LABEL_WZRY:
	//	tableName = for_game.TABLE_ESPORTS_WZRY_REAL_TIME_DATA
	//	for _, v := range reqMsg.GetWzry() {
	//		b1 := bson.M{"_id": v.Id}
	//		b2 := v
	//		data = append(data, b1, b2)
	//		roundLs = append(roundLs, v.GetGameRound())
	//	}
	//case for_game.ESPORTS_LABEL_LOL:
	//	tableName = for_game.TABLE_ESPORTS_LOL_REAL_TIME_DATA
	//	for _, v := range reqMsg.GetLol() {
	//		b1 := bson.M{"_id": v.Id}
	//		b2 := v
	//		data = append(data, b1, b2)
	//		roundLs = append(roundLs, v.GetGameRound())
	//	}
	//default:
	//	return easygo.NewFailMsg("LabelId参数错误")
	//}
	//
	//for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, tableName, data)

	//for _, r := range roundLs {
	//	for_game.SetRedisGameRealTime(id, r)
	//}

	msg := fmt.Sprintf("修改比赛[%d]实时数据", id)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)

	return easygo.EmptyMsg
}

//评论查询
func (s *cls4) RpcCommentList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	var table string
	switch reqMsg.GetType() {
	case 1:
		table = for_game.TABLE_ESPORTS_COMMENT_NEWS
	case 2:
		table = for_game.TABLE_ESPORTS_COMMENT_NEWS_REPLY
	case 3:
		table = for_game.TABLE_ESPORTS_COMMENT_VIDEO
	case 4:
		table = for_game.TABLE_ESPORTS_COMMENT_VIDEO_REPLY
	default:
		return easygo.NewFailMsg("Type参数错误")
	}

	findBson := bson.M{"ParentId": reqMsg.GetId()}
	sort := []string{"-CreateTime"}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["CommentId"] = reqMsg.GetListType()
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() != 0 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, table, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.TableESportComment
	for _, li := range lis {
		one := &share_message.TableESportComment{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.CommentListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//评论删除
func (s *cls4) RpcDelComment(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.CommentDelRequest) easygo.IMessage {
	idList := reqMsg.GetIds64()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}

	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		return easygo.NewFailMsg("Id参数错误")
	}

	var table, upTable, typeName, comTable string
	switch reqMsg.GetType() {
	case 1:
		table = for_game.TABLE_ESPORTS_COMMENT_NEWS
		upTable = for_game.TABLE_ESPORTS_NEWS
		typeName = "资讯评论"
	case 2:
		table = for_game.TABLE_ESPORTS_COMMENT_NEWS_REPLY
		upTable = for_game.TABLE_ESPORTS_NEWS
		comTable = for_game.TABLE_ESPORTS_COMMENT_NEWS
		typeName = "资讯评论回复"
	case 3:
		table = for_game.TABLE_ESPORTS_COMMENT_VIDEO
		upTable = for_game.TABLE_ESPORTS_VIDEO
		typeName = "视频评论"
	case 4:
		table = for_game.TABLE_ESPORTS_COMMENT_VIDEO_REPLY
		upTable = for_game.TABLE_ESPORTS_VIDEO
		comTable = for_game.TABLE_ESPORTS_COMMENT_VIDEO
		typeName = "视频评论回复"
	default:
		return easygo.NewFailMsg("Type参数错误")
	}

	one := for_game.FindOne(for_game.MONGODB_NINGMENG, upTable, bson.M{"_id": reqMsg.GetId()})
	if one == nil {
		return easygo.NewFailMsg("Id对象不存在")
	}

	findBson := bson.M{"_id": bson.M{"$in": idList}}
	updateBson := bson.M{"$set": bson.M{"Status": for_game.ESPORTS_COMM_STATUS_3}}
	info, _ := for_game.UpdateAllMgo(for_game.MONGODB_NINGMENG, table, findBson, updateBson)
	updated := info.Updated
	if updated > 0 {
		upfindBson := bson.M{"_id": reqMsg.GetId()}
		upCountBson := bson.M{"$inc": bson.M{"CommentCount": -updated}}
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, upTable, upfindBson, upCountBson, false)
	}

	if table == for_game.TABLE_ESPORTS_COMMENT_NEWS_REPLY || table == for_game.TABLE_ESPORTS_COMMENT_VIDEO_REPLY {
		M := []bson.M{
			{"$match": findBson},
			{"$group": bson.M{"_id": "$CommentId", "Count": bson.M{"$sum": 1}}},
		}
		list := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, table, M, 0, 0)

		for _, li := range list {
			one := &share_message.PipeIntCount{}
			for_game.StructToOtherStruct(li, one)
			for_game.FindAndModify(for_game.MONGODB_NINGMENG, comTable, bson.M{"_id": one.GetId()}, bson.M{"$inc": bson.M{"ReplyCount": -one.GetCount()}}, false)
		}
	}

	var ids string
	count := len(idList)
	for i, t := range idList {
		ids += easygo.AnytoA(t)
		if i < count-1 {
			ids += ","
		}

		if SendEsportsNewsTimeMgr.GetTimerById(t) != nil {
			SendEsportsNewsTimeMgr.DelTimerList(t)
		}
	}

	msg := fmt.Sprintf("批量删除%s: %s", typeName, ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)

	return easygo.EmptyMsg
}

//批量上传评论
func (s *cls4) RpcUploadComment(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.CommentUploadRequest) easygo.IMessage {
	id := reqMsg.GetId()
	if reqMsg.Id == nil || id == 0 {
		return easygo.NewFailMsg("Id参数错误")
	}

	var table string
	var table_p string
	var startTime int64
	switch reqMsg.GetType() {
	case 301:
		table = for_game.TABLE_ESPORTS_COMMENT_NEWS
		table_p = for_game.TABLE_ESPORTS_NEWS
	case 302:
		table = for_game.TABLE_ESPORTS_COMMENT_VIDEO
		table_p = for_game.TABLE_ESPORTS_VIDEO
	default:
		return easygo.NewFailMsg("Type参数错误")
	}

	one := for_game.FindOne(for_game.MONGODB_NINGMENG, table_p, bson.M{"_id": id})
	if one == nil {
		return easygo.NewFailMsg("Id参数无效")
	}

	switch reqMsg.GetType() {
	case 301:
		art := &share_message.TableESPortsRealTimeInfo{}
		for_game.StructToOtherStruct(one, art)
		startTime = art.GetBeginEffectiveTime()
	case 302:
		art := &share_message.TableESPortsVideoInfo{}
		for_game.StructToOtherStruct(one, art)
		startTime = art.GetBeginEffectiveTime()
	}

	if startTime == 0 {
		return easygo.NewFailMsg("出现目标对象发布时间为空异常")
	}

	list := reqMsg.GetList()
	if len(list) == 0 {
		return easygo.NewFailMsg("评论内容不能为空")
	}

	count := len(list)
	players := for_game.GetRandPlayerByTypes([]int32{2, 3, 4, 5}, int32(count)) //随机取一个用户
	if len(players) == 0 {
		return easygo.NewFailMsg("没有运营账户，请先添加")
	}

	endTime := easygo.NowTimestamp()
	online := (endTime - startTime) / int64(count)

	var infoList []interface{}
	for i, c := range list {
		player := players[i]
		tindex := for_game.RandInt(0, int(online))
		ctime := startTime + int64(tindex)
		info := &share_message.TableESportComment{
			Id:             easygo.NewInt64(for_game.NextId(table)),
			Content:        easygo.NewString(c),
			ThumbsUpCount:  easygo.NewInt32(0),
			PlayerId:       easygo.NewInt64(player.GetPlayerId()),
			PlayerNickName: easygo.NewString(player.GetNickName()),
			ParentId:       easygo.NewInt64(id),
			MenuId:         easygo.NewInt32(reqMsg.GetType()),
			AppLabelID:     easygo.NewInt64(0),
			ReplyCount:     easygo.NewInt32(0),
			PlayerIconUrl:  easygo.NewString(player.GetHeadIcon()),
			Status:         easygo.NewInt32(1),
			CreateTime:     easygo.NewInt64(ctime),
		}
		startTime = ctime
		infoList = append(infoList, info)
	}

	for_game.InsertAllMgo(for_game.MONGODB_NINGMENG, table, infoList...)
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, table_p, bson.M{"_id": id}, bson.M{"$inc": bson.M{"CommentCount": count}}, false)

	msg := fmt.Sprintf("%d批量添加评论%d条", id, count)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)

	return easygo.EmptyMsg
}

//电竞系统消息(推送管理)
func (s *cls4) RpcSportSysNotice(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	sort := []string{"-CreateTime"}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		findBson["Title"] = reqMsg.GetKeyword()
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["RecipientType"] = reqMsg.GetListType()
	}

	if reqMsg.Status != nil && reqMsg.GetStatus() < 1000 {
		findBson["Status"] = reqMsg.GetStatus()
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_SYS_MSG, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.TableESPortsSysMsg
	for _, li := range lis {
		one := &share_message.TableESPortsSysMsg{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	return &brower_backstage.SportSysNoticeResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//发送电竞系统消息
func (s *cls4) RpcSendSportSysNotice(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.TableESPortsSysMsg) easygo.IMessage {
	if reqMsg.FailureTime == nil {
		reqMsg.FailureTime = easygo.NewInt64(0)
	} else if reqMsg.GetFailureTime() > 0 && reqMsg.GetFailureTime() < easygo.NowTimestamp() {
		return easygo.NewFailMsg("过期时间不能小于当前时间")
	}

	if !reqMsg.GetIsPush() && !reqMsg.GetIsMessageCenter() {
		return easygo.NewFailMsg("缺少必填项")
	}

	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_SYS_MSG))
	}

	if reqMsg.EffectiveTime == nil || reqMsg.GetEffectiveTime() == 0 {
		reqMsg.EffectiveTime = easygo.NewInt64(easygo.NowTimestamp())
		reqMsg.Status = easygo.NewInt32(for_game.ESPORTS_STATUS_1)
	}
	reqMsg.CreateTime = easygo.NewInt64(easygo.NowTimestamp())

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	one := for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_SYS_MSG, queryBson, updateBson, true)
	if one != nil && reqMsg.GetEffectiveType() == ES_SEND_TYPE_FUTURE {
		AddSendEsportsTime(ES_SYSMSG, one)
	} else {
		easygo.Spawn(ChooseOneHall, int32(0), "RpcSendSportSysNoticeToHall", reqMsg)
	}

	msg := fmt.Sprintf("创建电竞系统消息:%d", reqMsg.GetId())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)

	return easygo.EmptyMsg
}

//注单明细列表
func (s *cls4) RpcBetSlipList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	sort := []string{"-CreateTime"}
	findBson := bson.M{}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetDownType() {
		case 1:
			findBson["UniqueGameId"] = easygo.StringToInt64noErr(reqMsg.GetKeyword())
		case 2:
			findBson["$or"] = []bson.M{{"PlayInfo.Account": reqMsg.GetKeyword()}, {"PlayInfo.Phone": reqMsg.GetKeyword()}}
		case 3:
			findBson["_id"] = easygo.StringToInt64noErr(reqMsg.GetKeyword())
		}
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}
	if reqMsg.Status != nil && reqMsg.GetStatus() != 0 {
		findBson["BetStatus"] = easygo.AnytoA(reqMsg.GetStatus())
	}
	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["AppLabelId"] = reqMsg.GetListType()
	}
	if reqMsg.Type != nil && reqMsg.GetType() != 0 {
		findBson["BetResult"] = easygo.AnytoA(reqMsg.GetType())
	}
	if reqMsg.SrtType != nil && reqMsg.GetSrtType() != "" {
		findBson["BetName"] = bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetSrtType(), Options: "im"}}
	}

	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
	var list []*share_message.TableESPortsGuessBetRecord
	for _, li := range lis {
		one := &share_message.TableESPortsGuessBetRecord{}
		for_game.StructToOtherStruct(li, one)
		if one.PlayInfo.NickName == nil || one.PlayInfo.GetNickName() == "" {
			pmg := for_game.GetRedisPlayerBase(one.PlayInfo.GetPlayId())
			one.PlayInfo.NickName = easygo.NewString(pmg.GetNickName())
		}
		list = append(list, one)
	}

	return &brower_backstage.BetSlipListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//注单输赢统计列表
func (s *cls4) RpcBetWinLosStatistics(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		findBson["$or"] = []bson.M{{"PlayInfo.Account": reqMsg.GetKeyword()}, {"PlayInfo.Phone": reqMsg.GetKeyword()}}
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}

	groupBson := bson.M{"_id": "$PlayInfo.PlayId"}
	groupBson["BetSlips"] = bson.M{"$sum": 1}
	groupBson["BetAmount"] = bson.M{"$sum": "$BetAmount"}
	groupBson["SuccessAmount"] = bson.M{"$sum": "$SuccessAmount"}
	groupBson["FailAmount"] = bson.M{"$sum": "$FailAmount"}
	groupBson["DisableAmount"] = bson.M{"$sum": "$DisableAmount"}
	groupBson["IllegalAmount"] = bson.M{"$sum": "$IllegalAmount"}
	groupBson["Account"] = bson.M{"$addToSet": "$PlayInfo.Account"}
	m := []bson.M{
		{"$match": findBson},
		{"$group": groupBson},
		{"$unwind": "$Account"},
		{"$sort": bson.M{"BetSlips": -1}},
	}

	sumAmout := int64(0)
	ls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD, m, 0, 0)
	count := len(ls)
	for _, item := range ls {
		one := &share_message.BetSlipReport{}
		for_game.StructToOtherStruct(item, one)
		sumAmout = sumAmout + (one.GetSuccessAmount() + one.GetDisableAmount() - one.GetBetAmount())
	}

	lis := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD, m, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()))
	var list []*share_message.BetSlipReport
	for _, li := range lis {
		one := &share_message.BetSlipReport{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}
	msg := &brower_backstage.BetSlipStatisticsResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
		SumAmount: easygo.NewInt64(sumAmout),
	}
	return msg
}

//注单赛事统计列表
func (s *cls4) RpcBetGameStatistics(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			findBson["UniqueGameId"] = easygo.StringToInt64noErr(reqMsg.GetKeyword())
		case 2:
			findBson["GameName"] = reqMsg.GetKeyword()
		}
	}

	if reqMsg.ListType != nil && reqMsg.GetListType() != 0 {
		findBson["AppLabelId"] = reqMsg.GetListType()
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}

	groupBson := bson.M{"_id": "$UniqueGameId"}
	groupBson["BetSlips"] = bson.M{"$sum": 1}
	groupBson["BetAmount"] = bson.M{"$sum": "$BetAmount"}
	groupBson["SuccessAmount"] = bson.M{"$sum": "$SuccessAmount"}
	groupBson["FailAmount"] = bson.M{"$sum": "$FailAmount"}
	groupBson["DisableAmount"] = bson.M{"$sum": "$DisableAmount"}
	groupBson["IllegalAmount"] = bson.M{"$sum": "$IllegalAmount"}
	groupBson["GameInfo"] = bson.M{"$addToSet": "$GameInfo"}
	groupBson["AppLabelName"] = bson.M{"$addToSet": "$AppLabelName"}
	m := []bson.M{
		{"$match": findBson},
		{"$group": groupBson},
		{"$unwind": "$GameInfo"},
		{"$unwind": "$AppLabelName"},
		{"$sort": bson.M{"BetSlips": -1}},
	}

	ls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD, m, 0, 0)
	count := len(ls)

	lis := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD, m, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()))
	var list []*share_message.BetSlipReport
	for _, li := range lis {
		one := &share_message.BetSlipReport{}
		for_game.StructToOtherStruct(li, one)
		findBson["UniqueGameId"] = one.GetId()
		mm := []bson.M{
			{"$match": findBson},
			{"$group": bson.M{"_id": "$PlayInfo.PlayId"}},
		}
		pls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD, mm, 0, 0)

		one.Players = easygo.NewInt64(len(pls))
		list = append(list, one)
	}
	msg := &brower_backstage.BetSlipStatisticsResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

//注单报表折线图
func (s *cls4) RpcBetSlipReportLine(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	groupBson := bson.M{"_id": "$AppLabelID"}
	m := []bson.M{
		{"$match": findBson},
		{"$group": groupBson},
	}
	ls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_BET_SLIP_REPORT, m, 0, 0)
	lines := []*brower_backstage.LineData{}
	for _, r := range ls {
		gameLabelID := r.(bson.M)["_id"].(int64)
		findBson["AppLabelID"] = gameLabelID
		list, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_BET_SLIP_REPORT, findBson, 0, 0)
		line := &brower_backstage.LineData{}
		line.LabelId = easygo.NewInt64(gameLabelID)
		for _, k := range list {
			one := &share_message.BetSlipReport{}
			for_game.StructToOtherStruct(k, one)
			// 折线图数据
			line.TimeData = append(line.TimeData, one.GetCreateTime())
			switch reqMsg.GetListType() {
			case 1:
				line.VelueData = append(line.VelueData, one.GetBetAmount())
			case 2:
				line.VelueData = append(line.VelueData, one.GetSuccessAmount())
			case 3:
				line.VelueData = append(line.VelueData, one.GetFailAmount())
			case 4:
				line.VelueData = append(line.VelueData, one.GetDisableAmount())
			case 5:
				line.VelueData = append(line.VelueData, one.GetIllegalAmount())
			case 6:
				line.VelueData = append(line.VelueData, one.GetSumAmount())
			default:
				return easygo.NewFailMsg("查询类型错误")
			}
		}
		lines = append(lines, line)
	}

	return &brower_backstage.LineChartsResponse{
		Line: lines,
	}
}

//注单报表柱状图
func (s *cls4) RpcBetSlipReportBar(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	findBson := bson.M{}
	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}

	groupBson := bson.M{"_id": "$AppLabelID"}
	groupBson["BetAmount"] = bson.M{"$sum": "$BetAmount"}
	groupBson["SuccessAmount"] = bson.M{"$sum": "$SuccessAmount"}
	groupBson["FailAmount"] = bson.M{"$sum": "$FailAmount"}
	groupBson["DisableAmount"] = bson.M{"$sum": "$DisableAmount"}
	groupBson["IllegalAmount"] = bson.M{"$sum": "$IllegalAmount"}
	groupBson["SumAmount"] = bson.M{"$sum": "$SumAmount"}
	m := []bson.M{
		{"$match": findBson},
		{"$group": groupBson},
	}
	ls := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_BET_SLIP_REPORT, m, 0, 0)
	lines := []*brower_backstage.LineData{}
	for _, r := range ls {
		one := &share_message.BetSlipReport{}
		for_game.StructToOtherStruct(r, one)
		line := &brower_backstage.LineData{}
		line.LabelId = easygo.NewInt64(r.(bson.M)["_id"].(int64)) //返回标题序号
		// 折线图数据
		line.StrData = append(line.StrData, "总投注金额", "总成功金额", "总失败金额", "总返还金额", "总扣除金额", "总盈利金额")
		line.TimeData = append(line.TimeData, one.GetBetAmount(), one.GetSuccessAmount(), one.GetFailAmount(), one.GetDisableAmount(), one.GetIllegalAmount(), one.GetSumAmount())
		lines = append(lines, line)
	}

	return &brower_backstage.LineChartsResponse{
		Line: lines,
	}
}

//注单操作
func (s *cls4) RpcBetSlipOperate(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.RpcBetSlipOperateRequest) easygo.IMessage {

	//判断参数
	if nil == reqMsg.GetId() || len(reqMsg.GetId()) <= 0 {
		return easygo.NewFailMsg("请传入注单订单号")
	}
	if reqMsg.GetOpt() != BET_SLIP_OPERATE_1 && reqMsg.GetOpt() != BET_SLIP_OPERATE_2 {
		return easygo.NewFailMsg("参数操作类型错误")
	}

	for _, betOrd := range reqMsg.GetId() {
		// 通过订单号查询投注订单
		one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD, bson.M{"_id": betOrd,
			"BetResult": for_game.GAME_GUESS_BET_RESULT_1})
		if one == nil {
			s := fmt.Sprintf("注单号:%v的未结算状态的注单不存在", betOrd)
			return easygo.NewFailMsg(s)
		}

		betRecord := &share_message.TableESPortsGuessBetRecord{}
		for_game.StructToOtherStruct(one, betRecord)

		if betRecord.GetBetStatus() != for_game.GAME_GUESS_BET_STATUS_1 {
			s := fmt.Sprintf("注单号:%v的注单状态已变化、刷新重试", betOrd)
			return easygo.NewFailMsg(s)
		}
		//注单无效操作
		if reqMsg.GetOpt() == BET_SLIP_OPERATE_1 {
			queryBson := bson.M{"_id": betOrd, "BetResult": for_game.GAME_GUESS_BET_RESULT_1}

			//设置值
			betRecord.DisableAmount = easygo.NewInt64(betRecord.GetBetAmount())
			betRecord.Reason = easygo.NewString(for_game.GAME_GUESS_BET_DISABLE_REASON_3)
			betRecord.BetStatus = easygo.NewString(for_game.GAME_GUESS_BET_STATUS_3)
			betRecord.BetResult = easygo.NewString(for_game.GAME_GUESS_BET_RESULT_4)
			betRecord.UpdateTime = easygo.NewInt64(time.Now().Unix())

			updateBson := bson.M{"$set": betRecord}
			//更新数据库
			modifyOne := for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD, queryBson, updateBson, false)

			if nil == modifyOne {
				s := fmt.Sprintf("注单号:%v的未结算状态的注单不存在", betOrd)
				return easygo.NewFailMsg(s)
			}

			st := for_game.ESPORTCOIN_TYPE_GUESS_BACK_IN
			msg := fmt.Sprintf("后台无效返还电竞币[%d]个", betRecord.GetDisableAmount())

			//取得流水订单号
			streamOrderId := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_IN, st)
			req := &share_message.ESportCoinRecharge{
				PlayerId:     easygo.NewInt64(betRecord.GetPlayInfo().GetPlayId()),
				RechargeCoin: easygo.NewInt64(betRecord.GetDisableAmount()),
				SourceType:   easygo.NewInt32(st),
				Note:         easygo.NewString(msg),
				ExtendLog: &share_message.GoldExtendLog{
					OrderId:    easygo.NewString(streamOrderId),          //流水的订单号
					MerchantId: easygo.NewString(betRecord.GetOrderId()), //这里设置电竞的订单号
				},
			}
			result := SendToPlayer(betRecord.GetPlayInfo().GetPlayId(), "RpcESportSendChangeESportCoins", req) //通知大厅
			err := for_game.ParseReturnDataErr(result)
			if err != nil {
				s := fmt.Sprintf(" 后台无效操作的注单号:%v的注单%v", betOrd, err.GetReason())
				return easygo.NewFailMsg(s)
			}

			//通知用户
			eSPortsGameOrderSysMsg := &share_message.TableESPortsGameOrderSysMsg{
				OrderId:      easygo.NewInt64(betRecord.GetOrderId()),
				UniqueGameId: easygo.NewInt64(betRecord.GetUniqueGameId()),
				BetTime:      easygo.NewInt64(betRecord.GetCreateTime()),
				Odds:         easygo.NewString(betRecord.GetOdds()),
				BetResult:    easygo.NewString(betRecord.GetBetResult()),
				BetTitle:     easygo.NewString(betRecord.GetBetTitle()),
				BetNum:       easygo.NewString(betRecord.GetBetNum()),
				BetName:      easygo.NewString(betRecord.GetBetName()),
				ResultAmount: easygo.NewInt64(betRecord.GetDisableAmount()),
				PlayerId:     easygo.NewInt64(betRecord.GetPlayInfo().GetPlayId()),
				BetAmount:    easygo.NewInt64(betRecord.GetBetAmount()),
			}
			//重新设置比赛名称
			if nil != betRecord.GetGameInfo() {
				eSPortsGameOrderSysMsg.GameName = easygo.NewString(betRecord.GetGameInfo().GetGameName() + " " +
					betRecord.GetGameInfo().GetTeamAName() +
					" VS " +
					betRecord.GetGameInfo().GetTeamBName())
			}

			rstNotify, errNotify := SendMsgRandToServerNew(for_game.SERVER_TYPE_SPORT_APPLY, "RpcBetSlipOperateNotify", eSPortsGameOrderSysMsg) //随机通知一台电竞应用服务器

			if errNotify != nil {
				s := fmt.Sprintf("注单号:%v的注单%v", betOrd, errNotify.GetReason())
				return easygo.NewFailMsg(s)
			}

			err1 := for_game.ParseReturnDataErr(rstNotify)
			if err1 != nil {
				s := fmt.Sprintf("注单号:%v的注单%v", betOrd, err1.GetReason())
				return easygo.NewFailMsg(s)
			}

			//注单违规操作
		} else if reqMsg.GetOpt() == BET_SLIP_OPERATE_2 {
			queryBson := bson.M{"_id": betOrd, "BetResult": for_game.GAME_GUESS_BET_RESULT_1}

			//设置值
			betRecord.IllegalAmount = easygo.NewInt64(betRecord.GetBetAmount())
			betRecord.Reason = easygo.NewString(for_game.GAME_GUESS_BET_DISABLE_REASON_3)
			betRecord.BetStatus = easygo.NewString(for_game.GAME_GUESS_BET_STATUS_4)
			betRecord.BetResult = easygo.NewString(for_game.GAME_GUESS_BET_RESULT_5)
			betRecord.UpdateTime = easygo.NewInt64(time.Now().Unix())

			updateBson := bson.M{"$set": betRecord}
			//更新数据库
			modifyOne := for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GUESS_BET_RECORD, queryBson, updateBson, false)

			if nil == modifyOne {
				s := fmt.Sprintf("注单号:%v的注单不存在", betOrd)
				return easygo.NewFailMsg(s)
			}

			//通知用户
			eSPortsGameOrderSysMsg := &share_message.TableESPortsGameOrderSysMsg{
				OrderId:      easygo.NewInt64(betRecord.GetOrderId()),
				UniqueGameId: easygo.NewInt64(betRecord.GetUniqueGameId()),
				BetTime:      easygo.NewInt64(betRecord.GetCreateTime()),
				Odds:         easygo.NewString(betRecord.GetOdds()),
				BetResult:    easygo.NewString(betRecord.GetBetResult()),
				BetTitle:     easygo.NewString(betRecord.GetBetTitle()),
				BetNum:       easygo.NewString(betRecord.GetBetNum()),
				BetName:      easygo.NewString(betRecord.GetBetName()),
				ResultAmount: easygo.NewInt64(0),
				PlayerId:     easygo.NewInt64(betRecord.GetPlayInfo().GetPlayId()),
				BetAmount:    easygo.NewInt64(betRecord.GetBetAmount()),
			}
			//重新设置比赛名称
			if nil != betRecord.GetGameInfo() {
				eSPortsGameOrderSysMsg.GameName = easygo.NewString(betRecord.GetGameInfo().GetGameName() + " " +
					betRecord.GetGameInfo().GetTeamAName() +
					" VS " +
					betRecord.GetGameInfo().GetTeamBName())
			}

			rstNotify, errNotify := SendMsgRandToServerNew(for_game.SERVER_TYPE_SPORT_APPLY, "RpcBetSlipOperateNotify", eSPortsGameOrderSysMsg) //随机通知一台电竞应用服务器

			if errNotify != nil {
				s := fmt.Sprintf("注单号:%v的注单%v", betOrd, errNotify.GetReason())
				return easygo.NewFailMsg(s)
			}

			err1 := for_game.ParseReturnDataErr(rstNotify)
			if err1 != nil {
				s := fmt.Sprintf("注单号:%v的注单%v", betOrd, err1.GetReason())
				return easygo.NewFailMsg(s)
			}
		}
	}

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, "电竞注单操作")

	return easygo.EmptyMsg
}

//充值赠送白名单
func (s *cls4) RpcGiveWhiteList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	lis, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GIVE_WHITELIST, bson.M{}, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()))
	var ids PLAYER_IDS
	for _, i := range lis {
		ids = append(ids, i.(bson.M)["_id"].(int64))
	}
	players := QueryplayerlistByIds(ids)
	pmap := make(map[PLAYER_ID]*share_message.PlayerBase)
	for _, p := range players {
		pmap[p.GetPlayerId()] = p
	}
	list := []*share_message.TableESportsGiveWhiteList{}
	for _, l := range lis {
		one := &share_message.TableESportsGiveWhiteList{}
		for_game.StructToOtherStruct(l, one)
		one.NickName = easygo.NewString(pmap[one.GetPlayerId()].GetNickName())
		one.Account = easygo.NewString(pmap[one.GetPlayerId()].GetAccount())
		list = append(list, one)
	}

	return &brower_backstage.GiveWhiteListRes{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
}

//添加充值赠送白名单
func (s *cls4) RpcAddGiveWhiteList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	psList := reqMsg.GetIdsStr()
	if len(psList) == 0 {
		return easygo.NewFailMsg("添加空气不好吧")
	}
	queryBson := bson.M{"Account": bson.M{"$in": psList}}
	lis, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE, queryBson, 0, 0)

	if len(lis) == 0 {
		return easygo.NewFailMsg("所有柠檬号都是无效的")
	}

	note := reqMsg.GetNote()

	var data []interface{}
	for _, v := range lis {
		id := v.(bson.M)["_id"].(int64)
		b1 := bson.M{"_id": id}
		b2 := &share_message.TableESportsGiveWhiteList{PlayerId: easygo.NewInt64(id), Note: easygo.NewString(note)}
		data = append(data, b1, b2)
	}
	for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GIVE_WHITELIST, data)

	var ids string
	idsarr := reqMsg.GetIdsStr()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += idsarr[i] + ","

		} else {
			ids += idsarr[i]
		}
	}
	msg := fmt.Sprintf("批量添加赠送白名单: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)

	return easygo.EmptyMsg
}

//删除充值赠送白名单
func (s *cls4) RpcDelGiveWhiteList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}

	findBson := bson.M{"_id": bson.M{"$in": idList}}
	for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GIVE_WHITELIST, findBson)

	var ids string
	count := len(idList)
	for i, t := range idList {
		ids += easygo.AnytoA(t)
		if i < count-1 {
			ids += ","
		}
	}

	msg := fmt.Sprintf("批量删除赠送白名单: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)

	return easygo.EmptyMsg
}

//查询电竞币兑换活动设置
func (s *cls4) RpcQueryRechargeEsAct(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_LUCKY_ACTIVITY, bson.M{"Types": 2})
	returnMsg := &share_message.Activity{}
	if one != nil {
		for_game.StructToOtherStruct(one, returnMsg)
	}
	return returnMsg
}

//修改电竞币兑换活动设置
func (s *cls4) RpcUpdateRechargeEsAct(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.Activity) easygo.IMessage {
	if reqMsg.StartTime == nil {
		return easygo.NewFailMsg("活动开始时间不能为空")
	}
	if reqMsg.EndTime == nil {
		return easygo.NewFailMsg("活动结束时间不能为空")
	}

	if reqMsg.GetEndTime() < easygo.NowTimestamp() {
		return easygo.NewFailMsg("活动结束时间不能小于当前时间")
	}

	if reqMsg.Status == nil {
		return easygo.NewFailMsg("活动开启状态不能为空")
	}

	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt64(for_game.NextId(for_game.TABLE_LUCKY_ACTIVITY))
	}
	if reqMsg.Title == nil {
		reqMsg.Title = easygo.NewString("电竞币兑换赠送活动")
	}
	if reqMsg.Types == nil {
		reqMsg.Types = easygo.NewInt32(2) //1集卡活动,2电竞币兑换赠送活动
	}

	one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_LUCKY_ACTIVITY, bson.M{"Types": 2})
	if one != nil {
		id := one.(bson.M)["_id"].(int64)
		reqMsg.Id = easygo.NewInt64(id)
	}

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_LUCKY_ACTIVITY, queryBson, updateBson, true)
	for_game.SetRedisActiveConfig()
	AddActivityCloseTime(reqMsg)

	msg := "修改电竞币兑换活动"
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)
	return easygo.EmptyMsg
}

//查询电竞币兑换配置列表
func (s *cls4) RpcRechargeEsCfg(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	lis, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_EXCHANGE_CFG, bson.M{}, 0, 0)
	list := []*share_message.TableESportsExchangeCfg{}
	for _, l := range lis {
		one := &share_message.TableESportsExchangeCfg{}
		for_game.StructToOtherStruct(l, one)
		list = append(list, one)
	}

	return &brower_backstage.RechargeEsCfgRes{
		List: list,
	}
}

//保存电竞币兑换配置
func (s *cls4) RpcSaveRechargeEsCfg(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.TableESportsExchangeCfg) easygo.IMessage {
	msg := fmt.Sprintf("保存电竞币兑换配置:%d", reqMsg.GetId())
	if reqMsg.Id == nil || reqMsg.GetId() < 1 {
		return easygo.NewFailMsg("电竞币数量不能小于1")
	}

	if reqMsg.Coin == nil || reqMsg.GetCoin() < 1 {
		return easygo.NewFailMsg("硬币价格不能小于1")
	}

	queryBson := bson.M{"_id": reqMsg.GetId()}
	updateBson := bson.M{"$set": reqMsg}
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_EXCHANGE_CFG, queryBson, updateBson, true)
	for_game.SetRedisEXChangeConfigs()
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)
	return easygo.EmptyMsg
}

//删除电竞币兑换配置
func (s *cls4) RpcDelRechargeEsCfg(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	idList := reqMsg.GetIds64()
	if len(idList) == 0 {
		return easygo.NewFailMsg("请先选择要删除的项")
	}

	findBson := bson.M{"_id": bson.M{"$in": idList}}
	for_game.DelAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_EXCHANGE_CFG, findBson)

	var ids string
	count := len(idList)
	for i, t := range idList {
		ids += easygo.AnytoA(t)
		if i < count-1 {
			ids += ","
		}
	}
	for_game.SetRedisEXChangeConfigs()
	msg := fmt.Sprintf("批量删除电竞币兑换配置: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ESPORTS_MANAGE, msg)

	return easygo.EmptyMsg
}

//电竞埋点报表查询
func (s *cls4) RpcPointsReportList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	msg := &brower_backstage.PointsReportRes{}
	var reportTable string
	var lis []interface{}
	var count int = 0
	sort := []string{"-_id"}
	findBson := bson.M{}
	switch reqMsg.GetSrtType() {
	case "Basis":
		switch reqMsg.GetListType() {
		case 1:
			reportTable = for_game.TABLE_ESPORTS_BASIS_POINTS_REPORT_DAY
		case 2:
			reportTable = for_game.TABLE_ESPORTS_BASIS_POINTS_REPORT_WEEK
		case 3:
			reportTable = for_game.TABLE_ESPORTS_BASIS_POINTS_REPORT_MONTH
		default:
			return easygo.NewFailMsg("ListType参数错误")
		}
		if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
			findBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		}
		lis, count = for_game.FindAll(for_game.MONGODB_NINGMENG, reportTable, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
		var list []*share_message.BasisPointsReport
		for _, li := range lis {
			one := &share_message.BasisPointsReport{}
			for_game.StructToOtherStruct(li, one)
			list = append(list, one)
		}
		msg.Basis = list
	case "Menu":
		switch reqMsg.GetListType() {
		case 1:
			reportTable = for_game.TABLE_ESPORTS_MENU_POINTS_REPORT_DAY
		case 2:
			reportTable = for_game.TABLE_ESPORTS_MENU_POINTS_REPORT_WEEK
		case 3:
			reportTable = for_game.TABLE_ESPORTS_MENU_POINTS_REPORT_MONTH
		default:
			return easygo.NewFailMsg("ListType参数错误")
		}
		if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
			findBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		}
		lis, count = for_game.FindAll(for_game.MONGODB_NINGMENG, reportTable, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
		var list []*share_message.MenuPointsReport
		for _, li := range lis {
			one := &share_message.MenuPointsReport{}
			for_game.StructToOtherStruct(li, one)
			list = append(list, one)
		}
		msg.Menu = list
	case "Label":
		switch reqMsg.GetListType() {
		case 1:
			reportTable = for_game.TABLE_ESPORTS_LABEL_POINTS_REPORT_DAY
		case 2:
			reportTable = for_game.TABLE_ESPORTS_LABEL_POINTS_REPORT_WEEK
		case 3:
			reportTable = for_game.TABLE_ESPORTS_LABEL_POINTS_REPORT_MONTH
		default:
			return easygo.NewFailMsg("ListType参数错误")
		}
		if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
			findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		}
		if reqMsg.Type != nil {
			findBson["MenuId"] = reqMsg.GetType()
		}
		sort = append(sort, "-CreateTime")
		lis, count = for_game.FindAll(for_game.MONGODB_NINGMENG, reportTable, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
		var list []*share_message.LabelPointsReport
		for _, li := range lis {
			one := &share_message.LabelPointsReport{}
			for_game.StructToOtherStruct(li, one)
			list = append(list, one)
		}
		msg.Label = list
	case "NewsAmuse":
		switch reqMsg.GetListType() {
		case 1:
			reportTable = for_game.TABLE_ESPORTS_NEWS_AMUSE_POINTS_REPORT_DAY
		case 2:
			reportTable = for_game.TABLE_ESPORTS_NEWS_AMUSE_POINTS_REPORT_WEEK
		case 3:
			reportTable = for_game.TABLE_ESPORTS_NEWS_AMUSE_POINTS_REPORT_MONTH
		default:
			return easygo.NewFailMsg("ListType参数错误")
		}
		if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
			findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		}
		if reqMsg.Type != nil {
			findBson["MenuId"] = reqMsg.GetType()
		}
		sort = append(sort, "-CreateTime")
		lis, count = for_game.FindAll(for_game.MONGODB_NINGMENG, reportTable, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
		var list []*share_message.NewsAmusePointsReport
		for _, li := range lis {
			one := &share_message.NewsAmusePointsReport{}
			for_game.StructToOtherStruct(li, one)
			list = append(list, one)
		}
		msg.NewsAmuse = list
	case "VdoHall":
		switch reqMsg.GetListType() {
		case 1:
			reportTable = for_game.TABLE_ESPORTS_VDOHALL_POINTS_REPORT_DAY
		case 2:
			reportTable = for_game.TABLE_ESPORTS_VDOHALL_POINTS_REPORT_WEEK
		case 3:
			reportTable = for_game.TABLE_ESPORTS_VDOHALL_POINTS_REPORT_MONTH
		default:
			return easygo.NewFailMsg("ListType参数错误")
		}
		if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
			findBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		}
		lis, count = for_game.FindAll(for_game.MONGODB_NINGMENG, reportTable, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
		var list []*share_message.VdoHallPointsReport
		for _, li := range lis {
			one := &share_message.VdoHallPointsReport{}
			for_game.StructToOtherStruct(li, one)
			list = append(list, one)
		}
		msg.VdoHall = list
	case "ApplyVdoHall":
		switch reqMsg.GetListType() {
		case 1:
			reportTable = for_game.TABLE_ESPORTS_APPLYVDOHALL_POINTS_REPORT_DAY
		case 2:
			reportTable = for_game.TABLE_ESPORTS_APPLYVDOHALL_POINTS_REPORT_WEEK
		case 3:
			reportTable = for_game.TABLE_ESPORTS_APPLYVDOHALL_POINTS_REPORT_MONTH
		default:
			return easygo.NewFailMsg("ListType参数错误")
		}
		if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
			findBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		}
		lis, count = for_game.FindAll(for_game.MONGODB_NINGMENG, reportTable, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
		var list []*share_message.ApplyVdoHallPointsReport
		for _, li := range lis {
			one := &share_message.ApplyVdoHallPointsReport{}
			for_game.StructToOtherStruct(li, one)
			list = append(list, one)
		}
		msg.ApplyVdoHall = list
	case "MatchLs":
		switch reqMsg.GetListType() {
		case 1:
			reportTable = for_game.TABLE_ESPORTS_MATCHLS_POINTS_REPORT_DAY
		case 2:
			reportTable = for_game.TABLE_ESPORTS_MATCHLS_POINTS_REPORT_WEEK
		case 3:
			reportTable = for_game.TABLE_ESPORTS_MATCHLS_POINTS_REPORT_MONTH
		default:
			return easygo.NewFailMsg("ListType参数错误")
		}
		if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
			findBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		}
		lis, count = for_game.FindAll(for_game.MONGODB_NINGMENG, reportTable, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
		var list []*share_message.MatchLsPointsReport
		for _, li := range lis {
			one := &share_message.MatchLsPointsReport{}
			for_game.StructToOtherStruct(li, one)
			list = append(list, one)
		}
		msg.MatchLs = list
	case "MatchDil":
		switch reqMsg.GetListType() {
		case 1:
			reportTable = for_game.TABLE_ESPORTS_MATCHDIL_POINTS_REPORT_DAY
		case 2:
			reportTable = for_game.TABLE_ESPORTS_MATCHDIL_POINTS_REPORT_WEEK
		case 3:
			reportTable = for_game.TABLE_ESPORTS_MATCHDIL_POINTS_REPORT_MONTH
		default:
			return easygo.NewFailMsg("ListType参数错误")
		}
		if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
			findBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		}
		lis, count = for_game.FindAll(for_game.MONGODB_NINGMENG, reportTable, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
		var list []*share_message.MatchDilPointsReport
		for _, li := range lis {
			one := &share_message.MatchDilPointsReport{}
			for_game.StructToOtherStruct(li, one)
			list = append(list, one)
		}
		msg.MatchDil = list
	case "Guess":
		switch reqMsg.GetListType() {
		case 1:
			reportTable = for_game.TABLE_ESPORTS_GUESS_POINTS_REPORT_DAY
		case 2:
			reportTable = for_game.TABLE_ESPORTS_GUESS_POINTS_REPORT_WEEK
		case 3:
			reportTable = for_game.TABLE_ESPORTS_GUESS_POINTS_REPORT_MONTH
		default:
			return easygo.NewFailMsg("ListType参数错误")
		}
		if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
			findBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		}
		lis, count = for_game.FindAll(for_game.MONGODB_NINGMENG, reportTable, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
		var list []*share_message.GuessPointsReport
		for _, li := range lis {
			one := &share_message.GuessPointsReport{}
			for_game.StructToOtherStruct(li, one)
			list = append(list, one)
		}
		msg.Guess = list
	case "Msg":
		switch reqMsg.GetListType() {
		case 1:
			reportTable = for_game.TABLE_ESPORTS_MSG_POINTS_REPORT_DAY
		case 2:
			reportTable = for_game.TABLE_ESPORTS_MSG_POINTS_REPORT_WEEK
		case 3:
			reportTable = for_game.TABLE_ESPORTS_MSG_POINTS_REPORT_MONTH
		default:
			return easygo.NewFailMsg("ListType参数错误")
		}
		if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
			findBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		}
		lis, count = for_game.FindAll(for_game.MONGODB_NINGMENG, reportTable, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
		var list []*share_message.MsgPointsReport
		for _, li := range lis {
			one := &share_message.MsgPointsReport{}
			for_game.StructToOtherStruct(li, one)
			list = append(list, one)
		}
		msg.Msg = list
	case "EsportCoin":
		switch reqMsg.GetListType() {
		case 1:
			reportTable = for_game.TABLE_ESPORTS_COIN_POINTS_REPORT_DAY
		case 2:
			reportTable = for_game.TABLE_ESPORTS_COIN_POINTS_REPORT_WEEK
		case 3:
			reportTable = for_game.TABLE_ESPORTS_COIN_POINTS_REPORT_MONTH
		default:
			return easygo.NewFailMsg("ListType参数错误")
		}
		if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
			findBson["_id"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
		}
		lis, count = for_game.FindAll(for_game.MONGODB_NINGMENG, reportTable, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)
		var list []*share_message.EsportCoinPointsReport
		for _, li := range lis {
			one := &share_message.EsportCoinPointsReport{}
			for_game.StructToOtherStruct(li, one)
			list = append(list, one)
		}
		msg.EsportCoin = list
	default:
		return easygo.NewFailMsg("SrtType参数错误")
	}

	msg.PageCount = easygo.NewInt32(count)
	return msg
}

//电竞流水查询
func (s *cls4) RpcQueryeSportCoinLog(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	if reqMsg.GetCurPage() == 1 {
		for_game.SaveESportCoinChangeLogToMongoDB()
	}

	findBson := bson.M{}
	sort := []string{"-_id"}
	if reqMsg.GetListType() > 0 {
		findBson["PayType"] = reqMsg.GetListType()
	}
	if reqMsg.GetDownType() > 0 {
		findBson["SourceType"] = reqMsg.GetDownType()
	}

	if reqMsg.Keyword != nil && reqMsg.GetKeyword() != "" {
		switch reqMsg.GetType() {
		case 1:
			player := QueryPlayerByAccountOrPhone(reqMsg.GetKeyword())
			if player != nil {
				findBson["PlayerId"] = player.GetPlayerId()
			}
		case 2:
			findBson["Extend.OrderId"] = reqMsg.GetKeyword()
		case 3:
			findBson["Extend.MerchantId"] = bson.M{"$regex": bson.RegEx{Pattern: reqMsg.GetKeyword(), Options: "i"}}
		}
	}

	if reqMsg.GetBeginTimestamp() > 0 && reqMsg.GetEndTimestamp() > 0 {
		findBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lt": reqMsg.GetEndTimestamp()}
	}

	list, count := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTCHANGELOG, findBson, int(reqMsg.GetPageSize()), int(reqMsg.GetCurPage()), sort...)

	var playerids []int64
	for _, item := range list {
		one := &for_game.ESportCoinChangeLog{}
		for_game.StructToOtherStruct(item, one)
		if for_game.IsContains(one.PlayerId, playerids) == -1 {
			playerids = append(playerids, one.PlayerId)
		}
	}

	players := for_game.GetAllPlayerBase(playerids)

	reqType := &brower_backstage.SourceTypeRequest{}
	soucetype := QuerySouceTypeList(reqType)
	sts := make(map[int32]*share_message.SourceType)
	for _, s := range soucetype {
		sts[s.GetKey()] = s
	}

	var msg []*brower_backstage.SportCoinLogList
	for _, l := range list {
		line := &for_game.ESportCoinChangeLog{}
		for_game.StructToOtherStruct(l, line)

		one := &brower_backstage.SportCoinLogList{
			InLine: &share_message.ESportCoinChangeLog{
				LogId:            easygo.NewInt64(line.LogId),
				PlayerId:         easygo.NewInt64(line.PlayerId),
				Account:          easygo.NewString(players[line.PlayerId].GetAccount()),
				ChangeESportCoin: easygo.NewInt64(line.ChangeESportCoin),
				PayType:          easygo.NewInt32(line.PayType),
				SourceType:       easygo.NewInt32(line.SourceType),
				SourceTypeName:   easygo.NewString(sts[line.SourceType].GetValue()),
				CurESportCoin:    easygo.NewInt64(line.CurESportCoin),
				ESportCoin:       easygo.NewInt64(line.ESportCoin),
				Note:             easygo.NewString(line.Note),
				CreateTime:       easygo.NewInt64(line.CreateTime),
			},
			Extend: line.Extend,
		}
		msg = append(msg, one)
	}

	return &brower_backstage.SportCoinLogResponse{
		List:      msg,
		PageCount: easygo.NewInt32(count),
	}
}
