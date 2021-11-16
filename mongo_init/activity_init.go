//
// 初始化数据
package mongo_init

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
)

//=============================================================活动初始化
//活动初始化
func InitActivity() []*share_message.Activity {
	activitys := []*share_message.Activity{}
	activitys = append(activitys, &share_message.Activity{
		Id:        easygo.NewInt64(1),
		Title:     easygo.NewString("参加集卡活动，瓜分50万现金-柠檬畅聊"),
		StartTime: easygo.NewInt64(1601222400), // 2020-09-28 00:00:00
		EndTime:   easygo.NewInt64(1602157800), // 2020-10-08 19:50:00
		Status:    easygo.NewInt32(for_game.ACTIVITY_OPEN),
		Types:     easygo.NewInt32(1),          //1集卡活动
		OpenTime:  easygo.NewInt64(1602158400), // 2020-10-08 20:00:00
		CloseTime: easygo.NewInt64(1602345599), // 2020-10-10 23:59:59
	})
	return activitys
}

//道具初始化
func InitProps() []*share_message.Props {
	props := []*share_message.Props{}
	props = append(props, &share_message.Props{
		Name:       easygo.NewString("和"),
		Status:     easygo.NewInt32(1),  //1启用 2禁用
		Count:      easygo.NewInt64(-1), //-1不限制
		ActivityId: easygo.NewInt64(1),  //活动id
	})
	props = append(props, &share_message.Props{
		Name:       easygo.NewString("柠"),
		Status:     easygo.NewInt32(1),  //1启用 2禁用
		Count:      easygo.NewInt64(-1), //-1不限制
		ActivityId: easygo.NewInt64(1),  //活动id
	})
	props = append(props, &share_message.Props{
		Name:       easygo.NewString("檬"),
		Status:     easygo.NewInt32(1),  //1启用 2禁用
		Count:      easygo.NewInt64(-1), //-1不限制
		ActivityId: easygo.NewInt64(1),  //活动id
	})
	props = append(props, &share_message.Props{
		Name:       easygo.NewString("趣"),
		Status:     easygo.NewInt32(1),   //1启用 2禁用
		Count:      easygo.NewInt64(150), //每天限量150张
		ActivityId: easygo.NewInt64(1),   //活动id
	})
	props = append(props, &share_message.Props{
		Name:       easygo.NewString("旅"),
		Status:     easygo.NewInt32(1),  //1启用 2禁用
		Count:      easygo.NewInt64(-1), //-1不限制
		ActivityId: easygo.NewInt64(1),  //活动id
	})
	props = append(props, &share_message.Props{
		Name:       easygo.NewString("行"),
		Status:     easygo.NewInt32(1),  //1启用 2禁用
		Count:      easygo.NewInt64(-1), //-1不限制
		ActivityId: easygo.NewInt64(1),  //活动id
	})
	return props
}

func InitPropsRate() []*share_message.PropsRate {
	propsRate := []*share_message.PropsRate{}
	propsRate = append(propsRate, &share_message.PropsRate{
		CreateTime: easygo.NewInt64(1601222400),
		Rate: []*share_message.Props{
			{
				Id:   easygo.NewInt64(1),
				Rate: easygo.NewInt32(45),
			},
			{
				Id:   easygo.NewInt64(2),
				Rate: easygo.NewInt32(45),
			},
			{
				Id:   easygo.NewInt64(3),
				Rate: easygo.NewInt32(5),
			},
			{
				Id:   easygo.NewInt64(4),
				Rate: easygo.NewInt32(0),
			},
			{
				Id:   easygo.NewInt64(5),
				Rate: easygo.NewInt32(5),
			},
			{
				Id:   easygo.NewInt64(6),
				Rate: easygo.NewInt32(0),
			},
		},
	})
	propsRate = append(propsRate, &share_message.PropsRate{
		CreateTime: easygo.NewInt64(1601308800),
		Rate: []*share_message.Props{
			{
				Id:   easygo.NewInt64(1),
				Rate: easygo.NewInt32(45),
			},
			{
				Id:   easygo.NewInt64(2),
				Rate: easygo.NewInt32(45),
			},
			{
				Id:   easygo.NewInt64(3),
				Rate: easygo.NewInt32(5),
			},
			{
				Id:   easygo.NewInt64(4),
				Rate: easygo.NewInt32(0),
			},
			{
				Id:   easygo.NewInt64(5),
				Rate: easygo.NewInt32(5),
			},
			{
				Id:   easygo.NewInt64(6),
				Rate: easygo.NewInt32(0),
			},
		},
	})
	propsRate = append(propsRate, &share_message.PropsRate{
		CreateTime: easygo.NewInt64(1601395200),
		Rate: []*share_message.Props{
			{
				Id:   easygo.NewInt64(1),
				Rate: easygo.NewInt32(45),
			},
			{
				Id:   easygo.NewInt64(2),
				Rate: easygo.NewInt32(45),
			},
			{
				Id:   easygo.NewInt64(3),
				Rate: easygo.NewInt32(5),
			},
			{
				Id:   easygo.NewInt64(4),
				Rate: easygo.NewInt32(0),
			},
			{
				Id:   easygo.NewInt64(5),
				Rate: easygo.NewInt32(5),
			},
			{
				Id:   easygo.NewInt64(6),
				Rate: easygo.NewInt32(0),
			},
		},
	})
	propsRate = append(propsRate, &share_message.PropsRate{
		CreateTime: easygo.NewInt64(1601481600),
		Rate: []*share_message.Props{
			{
				Id:   easygo.NewInt64(1),
				Rate: easygo.NewInt32(30),
			},
			{
				Id:   easygo.NewInt64(2),
				Rate: easygo.NewInt32(30),
			},
			{
				Id:   easygo.NewInt64(3),
				Rate: easygo.NewInt32(15),
			},
			{
				Id:   easygo.NewInt64(4),
				Rate: easygo.NewInt32(2),
			},
			{
				Id:   easygo.NewInt64(5),
				Rate: easygo.NewInt32(15),
			},
			{
				Id:   easygo.NewInt64(6),
				Rate: easygo.NewInt32(8),
			},
		},
	})
	propsRate = append(propsRate, &share_message.PropsRate{
		CreateTime: easygo.NewInt64(1601568000),
		Rate: []*share_message.Props{
			{
				Id:   easygo.NewInt64(1),
				Rate: easygo.NewInt32(30),
			},
			{
				Id:   easygo.NewInt64(2),
				Rate: easygo.NewInt32(30),
			},
			{
				Id:   easygo.NewInt64(3),
				Rate: easygo.NewInt32(15),
			},
			{
				Id:   easygo.NewInt64(4),
				Rate: easygo.NewInt32(2),
			},
			{
				Id:   easygo.NewInt64(5),
				Rate: easygo.NewInt32(15),
			},
			{
				Id:   easygo.NewInt64(6),
				Rate: easygo.NewInt32(8),
			},
		},
	})
	propsRate = append(propsRate, &share_message.PropsRate{
		CreateTime: easygo.NewInt64(1601654400),
		Rate: []*share_message.Props{
			{
				Id:   easygo.NewInt64(1),
				Rate: easygo.NewInt32(30),
			},
			{
				Id:   easygo.NewInt64(2),
				Rate: easygo.NewInt32(30),
			},
			{
				Id:   easygo.NewInt64(3),
				Rate: easygo.NewInt32(15),
			},
			{
				Id:   easygo.NewInt64(4),
				Rate: easygo.NewInt32(2),
			},
			{
				Id:   easygo.NewInt64(5),
				Rate: easygo.NewInt32(15),
			},
			{
				Id:   easygo.NewInt64(6),
				Rate: easygo.NewInt32(8),
			},
		},
	})
	propsRate = append(propsRate, &share_message.PropsRate{
		CreateTime: easygo.NewInt64(1601740800),
		Rate: []*share_message.Props{
			{
				Id:   easygo.NewInt64(1),
				Rate: easygo.NewInt32(30),
			},
			{
				Id:   easygo.NewInt64(2),
				Rate: easygo.NewInt32(30),
			},
			{
				Id:   easygo.NewInt64(3),
				Rate: easygo.NewInt32(15),
			},
			{
				Id:   easygo.NewInt64(4),
				Rate: easygo.NewInt32(2),
			},
			{
				Id:   easygo.NewInt64(5),
				Rate: easygo.NewInt32(15),
			},
			{
				Id:   easygo.NewInt64(6),
				Rate: easygo.NewInt32(8),
			},
		},
	})
	propsRate = append(propsRate, &share_message.PropsRate{
		CreateTime: easygo.NewInt64(1601827200),
		Rate: []*share_message.Props{
			{
				Id:   easygo.NewInt64(1),
				Rate: easygo.NewInt32(30),
			},
			{
				Id:   easygo.NewInt64(2),
				Rate: easygo.NewInt32(30),
			},
			{
				Id:   easygo.NewInt64(3),
				Rate: easygo.NewInt32(15),
			},
			{
				Id:   easygo.NewInt64(4),
				Rate: easygo.NewInt32(2),
			},
			{
				Id:   easygo.NewInt64(5),
				Rate: easygo.NewInt32(15),
			},
			{
				Id:   easygo.NewInt64(6),
				Rate: easygo.NewInt32(8),
			},
		},
	})
	propsRate = append(propsRate, &share_message.PropsRate{
		CreateTime: easygo.NewInt64(1601913600),
		Rate: []*share_message.Props{
			{
				Id:   easygo.NewInt64(1),
				Rate: easygo.NewInt32(30),
			},
			{
				Id:   easygo.NewInt64(2),
				Rate: easygo.NewInt32(30),
			},
			{
				Id:   easygo.NewInt64(3),
				Rate: easygo.NewInt32(15),
			},
			{
				Id:   easygo.NewInt64(4),
				Rate: easygo.NewInt32(2),
			},
			{
				Id:   easygo.NewInt64(5),
				Rate: easygo.NewInt32(15),
			},
			{
				Id:   easygo.NewInt64(6),
				Rate: easygo.NewInt32(8),
			},
		},
	})
	propsRate = append(propsRate, &share_message.PropsRate{
		CreateTime: easygo.NewInt64(1602000000),
		Rate: []*share_message.Props{
			{
				Id:   easygo.NewInt64(1),
				Rate: easygo.NewInt32(20),
			},
			{
				Id:   easygo.NewInt64(2),
				Rate: easygo.NewInt32(20),
			},
			{
				Id:   easygo.NewInt64(3),
				Rate: easygo.NewInt32(20),
			},
			{
				Id:   easygo.NewInt64(4),
				Rate: easygo.NewInt32(5),
			},
			{
				Id:   easygo.NewInt64(5),
				Rate: easygo.NewInt32(20),
			},
			{
				Id:   easygo.NewInt64(6),
				Rate: easygo.NewInt32(15),
			},
		},
	})
	propsRate = append(propsRate, &share_message.PropsRate{
		CreateTime: easygo.NewInt64(1602086400),
		Rate: []*share_message.Props{
			{
				Id:   easygo.NewInt64(1),
				Rate: easygo.NewInt32(20),
			},
			{
				Id:   easygo.NewInt64(2),
				Rate: easygo.NewInt32(20),
			},
			{
				Id:   easygo.NewInt64(3),
				Rate: easygo.NewInt32(20),
			},
			{
				Id:   easygo.NewInt64(4),
				Rate: easygo.NewInt32(5),
			},
			{
				Id:   easygo.NewInt64(5),
				Rate: easygo.NewInt32(20),
			},
			{
				Id:   easygo.NewInt64(6),
				Rate: easygo.NewInt32(15),
			},
		},
	})
	return propsRate
}

//================许愿池活动初始化
func InitWishCoinRechargeActivityCfg() []interface{} {
	var cfg []interface{}
	cfg = append(cfg, &share_message.WishCoinRechargeActivityCfg{
		Id:           easygo.NewInt64(72),
		Amount:       easygo.NewInt64(1200),
		FirstDiamond: easygo.NewInt64(60),
		FirstEsCoin:  easygo.NewInt64(300),
		FirstRatio:   easygo.NewInt64(4170),
		DailyDiamond: easygo.NewInt64(10),
		DailyEsCoin:  easygo.NewInt64(50),
		DailyRatio:   easygo.NewInt64(690),
	})
	cfg = append(cfg, &share_message.WishCoinRechargeActivityCfg{
		Id:           easygo.NewInt64(288),
		Amount:       easygo.NewInt64(4800),
		FirstDiamond: easygo.NewInt64(96),
		FirstEsCoin:  easygo.NewInt64(480),
		FirstRatio:   easygo.NewInt64(1670),
		DailyDiamond: easygo.NewInt64(45),
		DailyEsCoin:  easygo.NewInt64(225),
		DailyRatio:   easygo.NewInt64(780),
	})
	cfg = append(cfg, &share_message.WishCoinRechargeActivityCfg{
		Id:           easygo.NewInt64(588),
		Amount:       easygo.NewInt64(9800),
		FirstDiamond: easygo.NewInt64(199),
		FirstEsCoin:  easygo.NewInt64(995),
		FirstRatio:   easygo.NewInt64(1690),
		DailyDiamond: easygo.NewInt64(96),
		DailyEsCoin:  easygo.NewInt64(480),
		DailyRatio:   easygo.NewInt64(820),
	})
	cfg = append(cfg, &share_message.WishCoinRechargeActivityCfg{
		Id:           easygo.NewInt64(1188),
		Amount:       easygo.NewInt64(19800),
		FirstDiamond: easygo.NewInt64(419),
		FirstEsCoin:  easygo.NewInt64(2095),
		FirstRatio:   easygo.NewInt64(1760),
		DailyDiamond: easygo.NewInt64(198),
		DailyEsCoin:  easygo.NewInt64(990),
		DailyRatio:   easygo.NewInt64(830),
	})
	cfg = append(cfg, &share_message.WishCoinRechargeActivityCfg{
		Id:           easygo.NewInt64(2988),
		Amount:       easygo.NewInt64(49800),
		FirstDiamond: easygo.NewInt64(1066),
		FirstEsCoin:  easygo.NewInt64(5330),
		FirstRatio:   easygo.NewInt64(1780),
		DailyDiamond: easygo.NewInt64(548),
		DailyEsCoin:  easygo.NewInt64(2740),
		DailyRatio:   easygo.NewInt64(920),
	})
	cfg = append(cfg, &share_message.WishCoinRechargeActivityCfg{
		Id:           easygo.NewInt64(5988),
		Amount:       easygo.NewInt64(99800),
		FirstDiamond: easygo.NewInt64(2199),
		FirstEsCoin:  easygo.NewInt64(10995),
		FirstRatio:   easygo.NewInt64(1840),
		DailyDiamond: easygo.NewInt64(1198),
		DailyEsCoin:  easygo.NewInt64(5990),
		DailyRatio:   easygo.NewInt64(1000),
	})
	cfg = append(cfg, &share_message.WishCoinRechargeActivityCfg{
		Id:           easygo.NewInt64(17988),
		Amount:       easygo.NewInt64(299800),
		FirstDiamond: easygo.NewInt64(6699),
		FirstEsCoin:  easygo.NewInt64(33495),
		FirstRatio:   easygo.NewInt64(1860),
		DailyDiamond: easygo.NewInt64(3798),
		DailyEsCoin:  easygo.NewInt64(18990),
		DailyRatio:   easygo.NewInt64(1060),
	})
	cfg = append(cfg, &share_message.WishCoinRechargeActivityCfg{
		Id:           easygo.NewInt64(29988),
		Amount:       easygo.NewInt64(499800),
		FirstDiamond: easygo.NewInt64(11399),
		FirstEsCoin:  easygo.NewInt64(56995),
		FirstRatio:   easygo.NewInt64(1900),
		DailyDiamond: easygo.NewInt64(6598),
		DailyEsCoin:  easygo.NewInt64(32990),
		DailyRatio:   easygo.NewInt64(1100),
	})
	cfg = append(cfg, &share_message.WishCoinRechargeActivityCfg{
		Id:           easygo.NewInt64(47988),
		Amount:       easygo.NewInt64(799800),
		FirstDiamond: easygo.NewInt64(18699),
		FirstEsCoin:  easygo.NewInt64(93495),
		FirstRatio:   easygo.NewInt64(1950),
		DailyDiamond: easygo.NewInt64(10898),
		DailyEsCoin:  easygo.NewInt64(54490),
		DailyRatio:   easygo.NewInt64(1140),
	})

	return cfg
}
