package main

import (
	"game_server/e-sports/sport_common_dal"
	_ "game_server/e-sports/sport_common_dal"
	"game_server/easygo"
	"game_server/for_game"
)

func main() {
	initializer := for_game.NewInitializer()
	defer func() { // 若是异常了,确保异步日志有成功写盘
		logger := initializer.GetBeeLogger()
		if logger != nil {
			logger.Flush()
		}
	}()

	dict := easygo.KWAT{
		"logName":  "sport_test",
		"yamlPath": "config_share.yaml",
	}
	initializer.Execute(dict)

	for_game.InitRedisObjManager(33331)

	//sport_common_dal.AddTableESPortsFlowInfo_Test()
	//sport_common_dal.AddTableESPortsFlowInfo_Test()
	//sport_common_dal.AddLabel_test()
	//sport_common_dal.GetTableESPortsFlowLiveFollow_test()
	sport_common_dal.TestAddOrderMsg(1887436984)
	//sport_common_dal.AddCarousel(301, "调到哪1", "http://192.168.150.253:6060/image/lol.jpg")
	/*
		svr := sport_apply.ServiceForClient{}
		commn := &base.Common{
			UserId:     easygo.NewInt64(1887436001),
			Token:      nil,
			Flag:       nil,
			ServerType: nil,
		}
		fMsg := &client_hall.ESportInfoRequest{
			GameTypeId: easygo.NewInt64(for_game.ESPORT_FLOW_LIVE_FOLLOW_HISTORY),
			DataId:     easygo.NewInt64(10),
		}
		svr.RpcESportAddFollowLive(commn, fMsg)
	*/

	/*
		reqMsg := &client_hall.ESportInfoRequest{
			MenuId:     easygo.NewInt32(301),
			GameTypeId: easygo.NewInt64(1),
			DataId:     nil,
			ExtId:      nil,
		}
		rd := svr.RpcESportGetHomeInfo(commn, reqMsg)
		for _, v := range rd.LabelList {
			logs.Info("%s Id:%d lbid：%d wg:%d,tid:%d", v.GetTitle(), v.GetId(), v.GetLabelId(), v.GetWeight(), v.GetLabelType())
		}
		pageMsg := &client_hall.ESportPageRequest{
			MenuId:     easygo.NewInt32(1),
			TypeId:     easygo.NewInt64(3),
			LabelId:    easygo.NewInt64(10001),
			Page:       easygo.NewInt32(1),
			PageSize:   easygo.NewInt32(10),
			OrderField: nil,
			AscOrDesc:  nil,
		}
		prd := svr.RpcESportGetRealtimeList(commn, pageMsg)
		for _, v := range prd.List {
			logs.Info("%s Id:%d CoverBigImageUrl：%s wg:%d,tid:%d", v.GetTitle(), v.GetId(), v.GetCoverBigImageUrl())
		}
		//logs.Info(rd)
	*/
	/*
		infoMsg := &client_hall.ESportInfoRequest{
			MenuId:     easygo.NewInt32(301),
			GameTypeId: easygo.NewInt64(10001),
			DataId:     easygo.NewInt64(2),
			ExtId:      easygo.NewInt64(1),
		}

		logs.Info(svr.RpcESportGetRealtimeInfo(commn, infoMsg)) //获取单个资讯
		rComment := &client_hall.ESportCommentInfo{
			MenuId:    easygo.NewInt32(301),
			ParentId:  easygo.NewInt64(3),
			CommentId: nil,
			Content:   easygo.NewString("评论个毛线123"),
		}
		svr.RpcESportSendComment(commn, rComment)
		r2Comment := &client_hall.ESportCommentInfo{
			MenuId:    easygo.NewInt32(301),
			ParentId:  easygo.NewInt64(3),
			CommentId: easygo.NewInt64(2),
			Content:   easygo.NewString("评论个二级毛线123"),
		}

		svr.RpcESportSendComment(commn, r2Comment)

			gComment := &client_hall.ESportCommentRequest{
				MenuId:    easygo.NewInt32(301),
				ParentId:  easygo.NewInt64(3),
				CommentId: nil,
				Page:      easygo.NewInt32(1),
				PageSize:  easygo.NewInt32(10),
			}
			logs.Info(svr.RpcESportGetComment(commn, gComment))
			g2Comment := &client_hall.ESportCommentRequest{
				MenuId:    easygo.NewInt32(301),
				ParentId:  easygo.NewInt64(3),
				CommentId: easygo.NewInt64(2),
				Page:      easygo.NewInt32(1),
				PageSize:  easygo.NewInt32(10),
			}
			logs.Info(svr.RpcESportGetComment(commn, g2Comment))

			logs.Info(svr.RpcESportGetVideoList(commn, &client_hall.ESportVideoPageRequest{
				VideoType: easygo.NewInt32(1),
				MenuId:    easygo.NewInt32(301),
				TypeId:    easygo.NewInt64(3),
				LabelId:   easygo.NewInt64(10001),
				Page:      easygo.NewInt32(1),
				PageSize:  easygo.NewInt32(10),
			}))
	*/
	/*logs.Info(svr.RpcESportThumbsUp(commn, &client_hall.ESportInfoRequest{
		MenuId: easygo.NewInt32(301),
		DataId: easygo.NewInt64(2),
	}))

	logs.Info(svr.RpcESportThumbsUp(commn, &client_hall.ESportInfoRequest{
		MenuId: easygo.NewInt32(301),
		DataId: easygo.NewInt64(3),
		ExtId:  easygo.NewInt64(4),
	}))*/

}
