package sport_common_dal

import (
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"time"
	//   "time"
)

func AddTableESPortsRealTime_test(title string) {
	coverImageUrl := "http://192.168.150.253:6060/image/005.jpg"
	CreateTableESPortsRealTimeInfo(&share_message.TableESPortsRealTimeInfo{
		Id:               nil,
		Status:           easygo.NewInt32(1),
		IssueTime:        easygo.NewInt64(0),
		CoverBigImageUrl: easygo.NewString(coverImageUrl),
		CoverSmallImageUrl: []string{
			"http://192.168.150.253:6060/image/002.jpg", "http://192.168.150.253:6060/image/003.jpg", "http://192.168.150.253:6060/image/004.jpg",
		},
		Title:              easygo.NewString(title),
		Content:            easygo.NewString("<div>abc<img src=\"http://192.168.150.253:6060/image/006.jpg\"></div>"),
		AuthorPlayerId:     easygo.NewInt64(1),
		AuthorAccount:      easygo.NewString(title),
		Author:             easygo.NewString(title),
		DataSource:         easygo.NewString(title),
		LookCount:          easygo.NewInt32(0),
		LookCountSys:       easygo.NewInt32(1),
		ThumbsUpCount:      easygo.NewInt32(0),
		ThumbsUpCountSys:   easygo.NewInt32(2),
		AppLabelID:         easygo.NewInt64(0),
		AppLabelName:       easygo.NewString(title),
		BeginEffectiveTime: easygo.NewInt64(1),
		EffectiveType:      easygo.NewInt64(0),
		MenuId:             easygo.NewInt32(1),
		CommentCount:       easygo.NewInt32(1),
		ShowType:           easygo.NewInt32(1),
		LabelIds:           []int64{1, 2, 3},
	})

}

func AddTableESPortsVideoInfo_test(vt int32, labelid int32, id int64) {
	coverImageUrl := "http://192.168.150.253:6060/image/lol.jpg"
	vdeoUrl := "http://192.168.150.253:6060/video/20180905_101308.mp4"
	ln := for_game.LabelToESportNameMap[labelid]
	title := fmt.Sprintf("%s%d局", ln, id)
	CreateTableESPortsVideoInfo(&share_message.TableESPortsVideoInfo{
		IssueTime:          easygo.NewInt64(0),
		Status:             easygo.NewInt32(1),
		CoverImageUrl:      easygo.NewString(coverImageUrl),
		Title:              easygo.NewString(title),
		VideoUrl:           easygo.NewString(vdeoUrl),
		AuthorPlayerId:     easygo.NewInt64(1),
		AuthorAccount:      easygo.NewString(title),
		Author:             easygo.NewString(title),
		DataSource:         easygo.NewString(title),
		LookCount:          easygo.NewInt32(0),
		LookCountSys:       easygo.NewInt32(1),
		ThumbsUpCount:      easygo.NewInt32(0),
		ThumbsUpCountSys:   easygo.NewInt32(2),
		AppLabelID:         easygo.NewInt64(labelid),
		AppLabelName:       easygo.NewString(ln),
		BeginEffectiveTime: easygo.NewInt64(1),
		EffectiveType:      easygo.NewInt64(0),
		VideoType:          easygo.NewInt64(vt),
		IsRecommend:        easygo.NewInt32(1),
		IsHot:              easygo.NewInt32(1),
		UniqueGameId:       easygo.NewInt64(1),
		MenuId:             easygo.NewInt32(1),
		CommentCount:       easygo.NewInt32(1),
	})

}

func AddRealTimeo_Test() {
	AddTableESPortsRealTime_test("图11111图11111图11111图11111")
	AddTableESPortsRealTime_test("图11111图11111图11111图11112")
	AddTableESPortsRealTime_test("图11111图11111图11111图11113")
	AddTableESPortsRealTime_test("图11111图11111图11111图11114")
	AddTableESPortsRealTime_test("图11111图11111图11111图11115")
	AddTableESPortsRealTime_test("图11111图11111图11111图11116")
	AddTableESPortsRealTime_test("图11111图11111图11111图11117")
	AddTableESPortsRealTime_test("图11111图11111图11111图11118")
	AddTableESPortsRealTime_test("图11111图11111图11111图11119")
	AddTableESPortsRealTime_test("图11111图11111图11111图11110")
	AddTableESPortsRealTime_test("图11111图11111图11111图11111")
	AddTableESPortsRealTime_test("图11111图11111图11111图11112")
	AddTableESPortsRealTime_test("图11111图11111图11111图11113")
}

func AddTableESPortsFlowInfo_Test_Item(count int64, vt, laid int32) {
	idx := for_game.NextId(for_game.TABLE_ESPORTS_VIDEO)
	for i := idx; i < idx+count; i++ {
		AddTableESPortsVideoInfo_test(vt, laid, i)
	}
}
func AddTableESPortsFlowInfo_Test() {
	/*
		AddTableESPortsSysMsg(1)
		AddTableESPortsSysMsg(2)
		AddTableESPortsSysMsg(3)
		AddTableESPortsSysMsg(4)
	*/
	/*
		AddTableESPortsFlowInfo(for_game.ESPORT_FLOW_LIVE_HISTORY, 1887436001, 1)
		AddTableESPortsFlowInfo(for_game.ESPORT_FLOW_LIVE_HISTORY, 1887436001, 3)
		AddTableESPortsFlowInfo(for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY, 1887436001, 5)
		AddTableESPortsFlowInfo(for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY, 1887436001, 6)
		///*
			AddTableESPortsFlowInfo_Test_Item(90, 3, for_game.ESPORTS_LABEL_WZRY)
			AddTableESPortsFlowInfo_Test_Item(90, 2, for_game.ESPORTS_LABEL_WZRY)
			AddTableESPortsFlowInfo_Test_Item(120, 2, for_game.ESPORTS_LABEL_DOTA2)
			AddTableESPortsFlowInfo_Test_Item(140, 2, for_game.ESPORTS_LABEL_LOL)
			AddTableESPortsFlowInfo_Test_Item(100, 2, for_game.ESPORTS_LABEL_CSGO)
			AddTableESPortsFlowInfo_Test_Item(120, 2, for_game.ESPORTS_LABEL_OTHER)
			AddTableESPortsFlowInfo_Test_Item(110, 1, for_game.ESPORTS_LABEL_WZRY)
			AddTableESPortsFlowInfo_Test_Item(130, 1, for_game.ESPORTS_LABEL_DOTA2)
			AddTableESPortsFlowInfo_Test_Item(140, 1, for_game.ESPORTS_LABEL_LOL)
			AddTableESPortsFlowInfo_Test_Item(150, 1, for_game.ESPORTS_LABEL_CSGO)
			AddTableESPortsFlowInfo_Test_Item(130, 1, for_game.ESPORTS_LABEL_OTHER)
	*/
	///*
	AddTableESPortsFlowInfo(for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY, 2, 1)
	AddTableESPortsFlowInfo(for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY, 2, 2)
	AddTableESPortsFlowInfo(for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY, 2, 3)
	AddTableESPortsFlowInfo(for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY, 2, 4)
	AddTableESPortsFlowInfo(for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY, 2, 5)
	AddTableESPortsFlowInfo(for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY, 2, 6)
	//*/
}

//获取放映厅关注列表
func GetTableESPortsFlowLiveFollow_test() {

	GetTableESPortsFlowLiveFollow_t(for_game.TABLE_ESPORTS_FLOW_LIVE_HISTORY, for_game.TABLE_ESPORTS_VIDEO, 1, 10, "CreateTime", 1887436001)

}
func PadLeftString(size int, str string) string {
	temp := str
	num := fmt.Sprintf("%v", size-len(temp))
	fill := fmt.Sprintf("%0"+num+"s", "")
	result := fmt.Sprintf("%v%v", fill, temp)
	return result
}

func GetOrderId(idx int64) int64 {
	now := time.Now()
	str := fmt.Sprintf("%s%s", now.Format("20060102150405"), PadLeftString(2, fmt.Sprintf("%d", idx)))
	return easygo.AtoInt64(str)
}
func AddOrderMsg(pid, idx int64, status string) {
	resultAmount := int64(0)
	if status == for_game.GAME_GUESS_BET_RESULT_2 {
		resultAmount = for_game.RangeRand(idx, 100)
	}
	or := share_message.TableESPortsGameOrderSysMsg{
		OrderId:      easygo.NewInt64(GetOrderId(idx)),
		UniqueGameId: easygo.NewInt64(0),
		BetTime:      easygo.NewInt64(easygo.NowTimestamp()),
		Odds:         easygo.NewString("1.8"),
		BetResult:    easygo.NewString(status),
		BetTitle:     easygo.NewString("传说中的标题"),
		BetNum:       easygo.NewString("12"),
		BetName:      easygo.NewString("比赛项"),
		GameName:     easygo.NewString("游戏名"),
		CreateTime:   easygo.NewInt64(easygo.NowTimestamp()),
		UpdateTime:   easygo.NewInt64(easygo.NowTimestamp()),
		ResultAmount: easygo.NewInt64(resultAmount),
		PlayerId:     easygo.NewInt64(pid),
		BetAmount:    easygo.NewInt64(idx),
	}
	if status == for_game.GAME_GUESS_BET_RESULT_1 {
		CreateTableESPortsGameOrderSysMsg(&or)
	} else {
		CreateTableESPortsGameEndOrderSysMsg(&or)
	}

}
func TestAddSysMsg() {
	info := share_message.TableESPortsSysMsg{
		RecipientType:   easygo.NewInt64(0),
		Title:           easygo.NewString("系統消息"),
		Content:         easygo.NewString("系統消息系統消息"),
		Status:          easygo.NewInt32(1),
		CreateTime:      easygo.NewInt64(easygo.NowTimestamp()),
		JumpInfo:        nil,
		EffectiveTime:   easygo.NewInt64(1618455044),
		EffectiveType:   easygo.NewInt64(1),
		IsPush:          nil,
		IsMessageCenter: easygo.NewBool(true),
		FailureTime:     easygo.NewInt64(easygo.NowTimestamp() + 60),
	}
	CreateTableESPortsSysMsg(&info)

}
func TestAddOrderMsg(plyId int64) {

	//plyId := int64(1887436984)
	for i := int64(0); i < 3; i++ {
		AddOrderMsg(plyId, i, for_game.GAME_GUESS_BET_RESULT_1)
	}
	time.Sleep(time.Second)
	for i := int64(0); i < 4; i++ {
		AddOrderMsg(plyId, i, for_game.GAME_GUESS_BET_RESULT_2)
	}
	for i := int64(0); i < 5; i++ {
		AddOrderMsg(plyId, i, for_game.GAME_GUESS_BET_RESULT_3)
	}
	for i := int64(0); i < 6; i++ {
		AddOrderMsg(plyId, i, for_game.GAME_GUESS_BET_RESULT_4)
	}
	for i := int64(0); i < 6; i++ {
		AddOrderMsg(plyId, i, for_game.GAME_GUESS_BET_RESULT_5)
	}
}
