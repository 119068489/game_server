//
// 初始化一些需要预先写入到数据库的数据
package mongo_init

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"time"
)

//仅在第一次数据初始化时需要调用
func GetUsers() []*share_message.Manager {
	users := []*share_message.Manager{}

	// 只有第一个配置字段有注释，如果需要，请截图对比修改
	users = append(users, &share_message.Manager{
		Account:     easygo.NewString("admin"),                   //用户名
		Site:        easygo.NewString(for_game.MONGODB_NINGMENG), //数据库
		RealName:    easygo.NewString("超级管理员"),
		Phone:       easygo.NewString("18814139996"),
		Status:      easygo.NewInt32(0), //帐号状态  lg：0 正常，1 冻结
		Role:        easygo.NewInt32(0), //角色 lg：1超管，2站点管理员，3客服
		RoleType:    easygo.NewInt32(0), //角色类型
		CreateTime:  easygo.NewInt64(time.Now().Unix()),
		IsGoogleVer: easygo.NewBool(false), //是否开启谷歌验证器
		BindIp:      easygo.SliceValue{},   //绑定IP
	})

	return users
}

//初始化角色权限表的数据。超级管理员不需要插入具体的权限数据，
// func InitRolePowerMsg() []*share_message.RolePower {
// 	data := []*share_message.RolePower{
// 		{
// 			RoleName: easygo.NewString("管理员"),
// 			RoleType: easygo.NewInt32(0),
// 			Note:     easygo.NewString("系统管理员"),
// 		},
// 	}
// 	return data
// }

//通用额度配置
func InitPayGeneralQuota() *share_message.GeneralQuota {
	set := &share_message.GeneralQuota{
		Id:  easygo.NewInt32(1),
		Min: easygo.NewInt64(100),    //最小金额 单位分
		Max: easygo.NewInt64(500000), //最大金额 单位分
		Q1:  easygo.NewInt64(1000),   //位置1金额 单位分
		Q2:  easygo.NewInt64(5000),   //位置2金额 单位分
		Q3:  easygo.NewInt64(10000),  //位置3金额 单位分
		Q4:  easygo.NewInt64(20000),  //位置4金额 单位分
		Q5:  easygo.NewInt64(50000),  //位置5金额 单位分
		Q6:  easygo.NewInt64(100000), //位置6金额 单位分
	}
	return set
}

//支付类型
func InitPayType() []*share_message.PayType {
	settings := []*share_message.PayType{}
	settings = append(settings, &share_message.PayType{
		Name: easygo.NewString("微信支付"),
	})
	settings = append(settings, &share_message.PayType{
		Name: easygo.NewString("支付宝支付"),
	})
	settings = append(settings, &share_message.PayType{
		Name: easygo.NewString("银联支付"),
	})
	return settings
}

//支付场景
func InitPayScene() []*share_message.PayScene {
	settings := []*share_message.PayScene{}
	settings = append(settings, &share_message.PayScene{
		Name: easygo.NewString("app支付"),
	})
	settings = append(settings, &share_message.PayScene{
		Name: easygo.NewString("小程序"),
	})
	settings = append(settings, &share_message.PayScene{
		Name: easygo.NewString("扫码支付"),
	})
	settings = append(settings, &share_message.PayScene{
		Name: easygo.NewString("公众号"),
	})

	return settings
}

//支付设定===================================================(金额单位为分)
func InitPaymentSetting() []*share_message.PaymentSetting {
	settings := []*share_message.PaymentSetting{}
	settings = append(settings, &share_message.PaymentSetting{
		Name: easygo.NewString("默认入款设置"),
		//共用限制项
		Types:       easygo.NewInt32(1), //类型 (1入款，2出款)
		PlatformTax: easygo.NewInt64(0), //入金费用
	})
	settings = append(settings, &share_message.PaymentSetting{
		Name:        easygo.NewString("默认出款设置"),
		Types:       easygo.NewInt32(2), //类型 (1入款，2出款)
		FeeRate:     easygo.NewInt32(5), //手续费千分比(必须大于0的正整数) 使用时请除以1000
		PlatformTax: easygo.NewInt64(2), //平台手续费
		RealTax:     easygo.NewInt64(0), //我方服务费：我方扣除用户手续费
	})

	return settings
}

//支付平台===================================================
func InitPaymentPlatform() []*share_message.PaymentPlatform {
	settings := []*share_message.PaymentPlatform{}
	settings = append(settings, &share_message.PaymentPlatform{
		Name: easygo.NewString("秒到"),
	})
	settings = append(settings, &share_message.PaymentPlatform{
		Name: easygo.NewString("通联"),
	})
	settings = append(settings, &share_message.PaymentPlatform{
		Name: easygo.NewString("鹏聚代付"),
	})
	settings = append(settings, &share_message.PaymentPlatform{
		Name: easygo.NewString("汇聚"),
	})

	return settings
}

//支付平台通道===================================================(金额单位为分)
func InitPlatformChannel() []*share_message.PlatformChannel {
	settings := []*share_message.PlatformChannel{}
	settings = append(settings, &share_message.PlatformChannel{
		Name:             easygo.NewString("秒到微信支付"), //*string
		Types:            easygo.NewInt32(1),         //*int32
		PlatformId:       easygo.NewInt32(1),         //*int32
		PayMin:           easygo.NewInt64(1),         //*int64
		PayMax:           easygo.NewInt64(10000),     //*int64
		PayTypeId:        easygo.NewInt32(1),         //*int32
		PaymentSettingId: easygo.NewInt32(1),         //*int32
		StopAmount:       easygo.NewInt64(50000),     //*int64
		Weights:          easygo.NewInt32(1),         //*int32
		Status:           easygo.NewInt32(2),         //*int32
		PaySceneId:       easygo.NewInt32(2),         //*int32

	})
	settings = append(settings, &share_message.PlatformChannel{
		Name:             easygo.NewString("通联微信支付"), //*string
		Types:            easygo.NewInt32(1),         //*int32
		PlatformId:       easygo.NewInt32(2),         //*int32
		PayMin:           easygo.NewInt64(1),         //*int64
		PayMax:           easygo.NewInt64(1000),      //*int64
		PayTypeId:        easygo.NewInt32(1),         //*int32
		PaymentSettingId: easygo.NewInt32(1),         //*int32
		StopAmount:       easygo.NewInt64(50000),     //*int64
		Weights:          easygo.NewInt32(1),         //*int32
		Status:           easygo.NewInt32(1),         //*int32
		PaySceneId:       easygo.NewInt32(2),         //*int32

	})
	settings = append(settings, &share_message.PlatformChannel{
		Name:             easygo.NewString("鹏聚代付"), //*string
		Types:            easygo.NewInt32(2),       //*int32
		PlatformId:       easygo.NewInt32(3),       //*int32
		PayMin:           easygo.NewInt64(1),       //*int64
		PayMax:           easygo.NewInt64(1000),    //*int64
		PayTypeId:        easygo.NewInt32(3),       //*int32
		PaymentSettingId: easygo.NewInt32(2),       //*int32
		StopAmount:       easygo.NewInt64(50000),   //*int64
		Weights:          easygo.NewInt32(1),       //*int32
		Status:           easygo.NewInt32(2),       //*int32
		PaySceneId:       easygo.NewInt32(1),       //*int32

	})
	settings = append(settings, &share_message.PlatformChannel{
		Name:             easygo.NewString("汇聚支付"), //*string
		Types:            easygo.NewInt32(1),       //*int32
		PlatformId:       easygo.NewInt32(4),       //*int32
		PayMin:           easygo.NewInt64(1),       //*int64
		PayMax:           easygo.NewInt64(1000),    //*int64
		PayTypeId:        easygo.NewInt32(3),       //*int32
		PaymentSettingId: easygo.NewInt32(1),       //*int32
		StopAmount:       easygo.NewInt64(50000),   //*int64
		Weights:          easygo.NewInt32(1),       //*int32
		Status:           easygo.NewInt32(1),       //*int32
		PaySceneId:       easygo.NewInt32(1),       //*int32

	})
	settings = append(settings, &share_message.PlatformChannel{
		Name:             easygo.NewString("汇聚代付"), //*string
		Types:            easygo.NewInt32(2),       //*int32
		PlatformId:       easygo.NewInt32(4),       //*int32
		PayMin:           easygo.NewInt64(1),       //*int64
		PayMax:           easygo.NewInt64(1000),    //*int64
		PayTypeId:        easygo.NewInt32(3),       //*int32
		PaymentSettingId: easygo.NewInt32(2),       //*int32
		StopAmount:       easygo.NewInt64(50000),   //*int64
		Weights:          easygo.NewInt32(2),       //*int32
		Status:           easygo.NewInt32(1),       //*int32
		PaySceneId:       easygo.NewInt32(1),       //*int32

	})

	return settings
}

//客服分类
func InitManageTypes() []interface{} {
	//在这里进行初始化
	var managerTypes []interface{}
	managerTypes = append(managerTypes, &share_message.ManagerTypes{
		Name:   easygo.NewString("账号安全"),
		Status: easygo.NewInt32(1),
	})
	managerTypes = append(managerTypes, &share_message.ManagerTypes{
		Name:   easygo.NewString("充值转账"),
		Status: easygo.NewInt32(1),
	})
	managerTypes = append(managerTypes, &share_message.ManagerTypes{
		Name:   easygo.NewString("赌博 诈骗"),
		Status: easygo.NewInt32(1),
	})
	managerTypes = append(managerTypes, &share_message.ManagerTypes{
		Name:   easygo.NewString("色情 违法"),
		Status: easygo.NewInt32(1),
	})
	managerTypes = append(managerTypes, &share_message.ManagerTypes{
		Name:   easygo.NewString("好物商城"),
		Status: easygo.NewInt32(1),
	})
	managerTypes = append(managerTypes, &share_message.ManagerTypes{
		Name:   easygo.NewString("社交广场"),
		Status: easygo.NewInt32(1),
	})
	managerTypes = append(managerTypes, &share_message.ManagerTypes{
		Name:   easygo.NewString("其他问题"),
		Status: easygo.NewInt32(1),
	})
	return managerTypes
}

//后台管理员日志类型
func InitManageLogTypes() []*brower_backstage.KeyValue {
	//在这里进行初始化
	types := []*brower_backstage.KeyValue{}
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.LOGIN_BACKSTAGE),
		Value: easygo.NewString(for_game.LOGIN_BACKSTAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.SIGNOUT_BACKSTAGE),
		Value: easygo.NewString(for_game.SIGNOUT_BACKSTAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.USER_MANAGE),
		Value: easygo.NewString(for_game.USER_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.BSPLAYER_MANAGE),
		Value: easygo.NewString(for_game.BSPLAYER_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.ROLE_MANAGE),
		Value: easygo.NewString(for_game.ROLE_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.TEAM_MANAGE),
		Value: easygo.NewString(for_game.TEAM_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.SITE_MANAGE),
		Value: easygo.NewString(for_game.SITE_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.SYS_MAIL),
		Value: easygo.NewString(for_game.SYS_MAIL),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.SYS_MSG),
		Value: easygo.NewString(for_game.SYS_MSG),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.PAY_MANAGE),
		Value: easygo.NewString(for_game.PAY_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.FEATURES_MANAGE),
		Value: easygo.NewString(for_game.FEATURES_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.SHOP_MANAGE),
		Value: easygo.NewString(for_game.SHOP_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.OPERATION_MANAGE),
		Value: easygo.NewString(for_game.OPERATION_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.WAITER_MANAGE),
		Value: easygo.NewString(for_game.WAITER_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.SQUARE_MANAGE),
		Value: easygo.NewString(for_game.SQUARE_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.ADV_MANAGE),
		Value: easygo.NewString(for_game.ADV_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.MALL_MANAGE),
		Value: easygo.NewString(for_game.MALL_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.COINS_MANAGE),
		Value: easygo.NewString(for_game.COINS_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.PROPS_MANAGE),
		Value: easygo.NewString(for_game.PROPS_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.LOVE_MANAGE),
		Value: easygo.NewString(for_game.LOVE_MANAGE),
	})
	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.ESPORTS_MANAGE),
		Value: easygo.NewString(for_game.ESPORTS_MANAGE),
	})

	types = append(types, &brower_backstage.KeyValue{
		Key:   easygo.NewString(for_game.WISH_MANAGE),
		Value: easygo.NewString(for_game.WISH_MANAGE),
	})

	return types
}

//首页Tips配置
func InitIndexTips() []interface{} {
	return []interface{}{
		&share_message.IndexTips{
			Id:      easygo.NewInt32(1),
			Title:   easygo.NewString("谁喜欢我"),
			Weights: easygo.NewInt32(99),
			Status:  easygo.NewInt32(1),
			Types:   easygo.NewInt32(1),
		},
		&share_message.IndexTips{
			Id:      easygo.NewInt32(2),
			Title:   easygo.NewString("热门话题"),
			Weights: easygo.NewInt32(98),
			Status:  easygo.NewInt32(1),
			Types:   easygo.NewInt32(1),
		},
		&share_message.IndexTips{
			Id:      easygo.NewInt32(3),
			Title:   easygo.NewString("柠檬花田"),
			Weights: easygo.NewInt32(97),
			Status:  easygo.NewInt32(1),
			Types:   easygo.NewInt32(1),
		},
		&share_message.IndexTips{
			Id:      easygo.NewInt32(4),
			Title:   easygo.NewString("附近的人"),
			Weights: easygo.NewInt32(96),
			Status:  easygo.NewInt32(1),
			Types:   easygo.NewInt32(1),
		},
	}
}

//弹窗悬浮球配置
func InitPopSuspend() []interface{} {
	return []interface{}{
		&share_message.PopSuspend{
			Id:        easygo.NewInt32(1),
			Title:     easygo.NewString("恰柠檬"),
			IsPop:     easygo.NewBool(true),
			IsSuspend: easygo.NewBool(true),
		},
		&share_message.PopSuspend{
			Id:        easygo.NewInt32(2),
			Title:     easygo.NewString("消息"),
			IsPop:     easygo.NewBool(true),
			IsSuspend: easygo.NewBool(true),
		},
		&share_message.PopSuspend{
			Id:        easygo.NewInt32(3),
			Title:     easygo.NewString("广场"),
			IsPop:     easygo.NewBool(true),
			IsSuspend: easygo.NewBool(true),
		},
		&share_message.PopSuspend{
			Id:        easygo.NewInt32(4),
			Title:     easygo.NewString("电竞"),
			IsPop:     easygo.NewBool(true),
			IsSuspend: easygo.NewBool(true),
		},
		&share_message.PopSuspend{
			Id:        easygo.NewInt32(5),
			Title:     easygo.NewString("我的"),
			IsPop:     easygo.NewBool(true),
			IsSuspend: easygo.NewBool(true),
		},
	}
}
