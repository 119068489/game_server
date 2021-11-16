//
// 初始化系统参数数据
package mongo_init

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"

	"github.com/astaxie/beego/logs"
)

type InitSysParameterStruct struct {
	PushSet []*share_message.PushSettings
}

//初始化系统转账付款参数
func InitSysParameterCfg() []interface{} {
	//在这里进行初始化
	var sysparameter []interface{}
	sysparameter = append(sysparameter, &share_message.SysParameter{
		Id:                easygo.NewString(for_game.LIMIT_PARAMETER),
		IsTransfer:        easygo.NewBool(true),     //是否允许转账
		TransferOneMax:    easygo.NewInt64(200000),  //单次转账上限 分单位     0表示不限制
		TransferOneDayMax: easygo.NewInt64(500000),  //单日转账上限 分单位     0表示不限制
		IsRedPacket:       easygo.NewBool(true),     //是否允许发红包
		RedPacketMin:      easygo.NewInt64(1),       //单次发红包最小值  分单位     0表示不限制
		RedPacketMax:      easygo.NewInt64(20000),   //单次发红包最大值  分单位     0表示不限制
		IsRecharge:        easygo.NewBool(true),     //是否允许充值 充值限制
		RechargeMin:       easygo.NewInt64(1),       //单次充值最小值  分单位     0表示不限制
		RechargeMax:       easygo.NewInt64(500000),  //单次充值最大值  分单位     0表示不限制
		IsWithdrawal:      easygo.NewBool(true),     //是否允许提现 提现限制
		WithdrawalMin:     easygo.NewInt64(1000),    //单次提现最小值  分单位     0表示不限制
		WithdrawalMax:     easygo.NewInt64(500000),  //单次提现最大值  分单位     0表示不限制
		TeamRedPacketNum:  easygo.NewInt32(100),     //群红包个数  分单位     0表示不限制
		OutSum:            easygo.NewInt64(1000000), //单日累计出款总金额  分单位     0表示不限制
		OutTimes:          easygo.NewInt32(3),       //单日累计出款最大次数
		RiskControl:       easygo.NewInt64(0),       //单次风控额度 分单位     0表示不限制
		IsQRcode:          easygo.NewBool(true),     //是否允许二维码付款
	})
	sysparameter = append(sysparameter, &share_message.SysParameter{
		Id:           easygo.NewString(for_game.AVATAR_PARAMETER),
		MavatarCount: easygo.NewInt32(1425), //男头像
		WavatarCount: easygo.NewInt32(1423), //女头像
	})
	sysparameter = append(sysparameter, &share_message.SysParameter{
		Id:          easygo.NewString(for_game.INTEREST_PARAMETER),
		InterestMin: easygo.NewInt32(1), //选择兴趣标签下限
		InterestMax: easygo.NewInt32(3), //选择兴趣标签上限
	})
	textTags, imageTags := InitControlModerations()
	sysparameter = append(sysparameter, &share_message.SysParameter{
		Id:               easygo.NewString(for_game.OBJ_MODERATIONS),
		TextModerations:  textTags,  //文本屏蔽标签
		ImageModerations: imageTags, //图片屏蔽标签
	})
	sysparameter = append(sysparameter, &share_message.SysParameter{
		Id:           easygo.NewString(for_game.SQUAREHOT_PARAMETER),
		ZanScore:     easygo.NewInt32(1),
		CoinScore:    easygo.NewInt32(5),
		CommentScore: easygo.NewInt32(5),
		HotScore:     easygo.NewInt32(100),
		DampRatio:    easygo.NewInt32(5),
	})
	sysparameter = append(sysparameter, &share_message.SysParameter{
		Id:              easygo.NewString(for_game.ESPORT_PARAMETER),
		EsOneBetGold:    easygo.NewInt64(1000),     //电竞每用户单笔下注硬币上限
		EsOneDayBetGold: easygo.NewInt64(10000),    //电竞每用户单日下注总硬币上限
		EsDaySumGold:    easygo.NewInt64(10000000), //电竞单日所有用户下注总硬币上限
	}) //初始化电竞系统参数
	return sysparameter
}
func InitControlModerations() ([]*share_message.GrabTag, []*share_message.GrabTag) {
	//在这里进行初始化
	textTags := []*share_message.GrabTag{}
	imageTags := []*share_message.GrabTag{}
	textList := map[int32]string{
		20001: "政治",
		20002: "色情",
		20006: "涉毒违法",
		20007: "谩骂",
		20105: "广告引流",
		24001: "暴恐"}
	textValList := map[int32]int32{
		20001: 50,
		20002: 60,
		20006: 50,
		20007: 80,
		20105: 80,
		24001: 20}
	imageList := map[int32]string{
		20001: "政治",
		20002: "色情",
		20006: "涉毒违法",
		20007: "谩骂",
		20103: "性感",
		24001: "暴恐"}
	imageValList := map[int32]int32{
		20001: 90,
		20002: 90,
		20006: 100,
		20007: 100,
		20103: 95,
		24001: 100}
	for k, v := range textList {
		tags := &share_message.GrabTag{
			Id:    easygo.NewInt32(k),
			Name:  easygo.NewString(v),
			Count: easygo.NewInt64(textValList[k]),
		}
		textTags = append(textTags, tags)
	}
	for k, v := range imageList {
		tags := &share_message.GrabTag{
			Id:    easygo.NewInt32(k),
			Name:  easygo.NewString(v),
			Count: easygo.NewInt64(imageValList[k]),
		}
		imageTags = append(imageTags, tags)
	}
	return textTags, imageTags
}

//初始化推送配置
func InitPushSetCfg() []interface{} {
	var jsonData = []byte(`{
		"PushSet": [
			{
				"Id": 1,
				"Title": "虚拟商城",
				"ObjId": 101,
				"ObjTitle": "赠送硬币过期提醒",
				"ObjContent": "您有平台赠送硬币明日0点即将过期",
				"IsPush": true
			},
			{
				"Id": 1,
				"Title": "虚拟商城",
				"ObjId": 102,
				"ObjTitle": "道具过期提醒",
				"ObjContent": "您的%v即将过期,请及时查看",
				"IsPush": true
			},
			{
				"Id": 2,
				"Title": "社交广场",
				"ObjId": 201,
				"ObjTitle": "点赞我的动态",
				"ObjContent": "%v点赞了我的社交广场动态",
				"IsPush": true
			},
			{
				"Id": 2,
				"Title": "社交广场",
				"ObjId": 202,
				"ObjTitle": "评论我的动态",
				"ObjContent": "%v评论了我的社交广场动态",
				"IsPush": true
			},
			{
				"Id": 2,
				"Title": "社交广场",
				"ObjId": 203,
				"ObjTitle": "回复我的动态",
				"ObjContent": "%v回复了我在社交广场的评论",
				"IsPush": true
			},
			{
				"Id": 2,
				"Title": "社交广场",
				"ObjId": 204,
				"ObjTitle": "我关注的人发布新动态",
				"ObjContent": "%v发布了新的社交广场动态",
				"IsPush": true
			}
		]
	}`)

	data := &InitSysParameterStruct{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		logs.Info(err)
	}

	var list []interface{}
	lis := data.PushSet
	for _, li := range lis {
		list = append(list, li)
	}

	one := &share_message.SysParameter{
		Id:      easygo.NewString(for_game.PUSH_PARAMETER),
		PushSet: lis,
	}

	var returnls []interface{}
	returnls = append(returnls, one)
	return returnls
}
