//
// 初始化开关列表数据
package mongo_init

import (
	"game_server/easygo"
	"game_server/pb/share_message"
)

func InitSourcetypeCfg() []interface{} {
	//在这里进行初始化
	var sourcetype []interface{}
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(100),
		Value:   easygo.NewString("人工入款"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(101),
		Value:   easygo.NewString("充值"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(102),
		Value:   easygo.NewString("收红包"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(103),
		Value:   easygo.NewString("转入"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(104),
		Value:   easygo.NewString("二维码收款"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(111),
		Value:   easygo.NewString("红包退款"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(112),
		Value:   easygo.NewString("转账退款"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(113),
		Value:   easygo.NewString("商家退款"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(114),
		Value:   easygo.NewString("提现退款"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(115),
		Value:   easygo.NewString("商城卖家货款"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(1),
	})

	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(200),
		Value:   easygo.NewString("人工出款"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(201),
		Value:   easygo.NewString("提现"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(202),
		Value:   easygo.NewString("发红包"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(203),
		Value:   easygo.NewString("转出"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(204),
		Value:   easygo.NewString("二维码付款"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(215),
		Value:   easygo.NewString("罚没"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(216),
		Value:   easygo.NewString("手续费"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(217),
		Value:   easygo.NewString("商城消费"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(218),
		Value:   easygo.NewString("平台税收"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(219),
		Value:   easygo.NewString("兑换硬币"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(1),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(220),
		Value:   easygo.NewString("虚拟商场消费"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(1),
	})
	//活动
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(409),
		Value:   easygo.NewString("活动奖励"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(1),
	})

	// 硬币加
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(500),
		Value:   easygo.NewString("系统赠送"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(2),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(501),
		Value:   easygo.NewString("兑换"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(2),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(502),
		Value:   easygo.NewString("被投币"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(2),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(503),
		Value:   easygo.NewString("活动奖励"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(2),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(505),
		Value:   easygo.NewString("许愿池回收"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(2),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(506),
		Value:   easygo.NewString("许愿池守护者收益"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(2),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(507),
		Value:   easygo.NewString("许愿池平台回收"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(2),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(508),
		Value:   easygo.NewString("许愿池抽奖返利"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(2),
	})
	// 硬币减
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(600),
		Value:   easygo.NewString("系统回收"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(2),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(601),
		Value:   easygo.NewString("虚拟商城消费"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(2),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(602),
		Value:   easygo.NewString("投币"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(2),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(603),
		Value:   easygo.NewString("过期回收"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(2),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(604),
		Value:   easygo.NewString("系统罚没"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(2),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(605),
		Value:   easygo.NewString("许愿池兑换钻石"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(2),
	})

	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(606),
		Value:   easygo.NewString("电竞兑换"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(2),
	})

	// 钻石加
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(700),
		Value:   easygo.NewString("兑换"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(3),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(701),
		Value:   easygo.NewString("许愿池抽奖返利"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(3),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(702),
		Value:   easygo.NewString("许愿池抽奖失败钻石返回"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(3),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(703),
		Value:   easygo.NewString("守护者收益"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(3),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(704),
		Value:   easygo.NewString("回收"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(3),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(705),
		Value:   easygo.NewString("邮费返回"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(3),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(706),
		Value:   easygo.NewString("系统赠送"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(3),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(707),
		Value:   easygo.NewString("许愿池活动所得"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(3),
	})

	// 钻石减
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(800),
		Value:   easygo.NewString("抽奖"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(3),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(801),
		Value:   easygo.NewString("运费"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(3),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(802),
		Value:   easygo.NewString("回收失败返回"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(3),
	})
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(803),
		Value:   easygo.NewString("系统扣除"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(3),
	})
	//电竞币加
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(900),
		Value:   easygo.NewString("竞猜返还"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(4),
	})

	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(901),
		Value:   easygo.NewString("电竞兑换"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(4),
	})

	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(902),
		Value:   easygo.NewString("电竞兑换赠送"),
		Type:    easygo.NewInt32(1),
		Channel: easygo.NewInt32(4),
	})

	//电竞币减
	sourcetype = append(sourcetype, &share_message.SourceType{
		Key:     easygo.NewInt32(1000),
		Value:   easygo.NewString("竞猜投注"),
		Type:    easygo.NewInt32(2),
		Channel: easygo.NewInt32(4),
	})
	return sourcetype
}
