package for_game

import (
	"game_server/easygo"
	"game_server/pb/share_message"
)

//映射已经实例化的群信息到指定服务器上
type PayChannelManager struct {
	PaymentSettingList  []*share_message.PaymentSetting
	PlatformChannelList []*share_message.PlatformChannel
	GeneralQuota        *share_message.GeneralQuota
	Mutex               easygo.RLock
}

func NewPayChannelManager() *PayChannelManager {
	p := &PayChannelManager{}
	p.Init()
	return p
}
func (self *PayChannelManager) Init() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.PaymentSettingList = QueryPaymentSettingList() //出入款的各种限制 单日限额 限次等等
	self.PlatformChannelList = QueryPlatformChannelList(1)
	self.GeneralQuota = GetGeneralQuota()
	// logs.Info("渠道:", self.PlatformChannelList)
	// logs.Info("额度:", self.GeneralQuota)
	// logs.Info("配置:", self.PaymentSettingList)
}

//渠道信息设置
func (self *PayChannelManager) SetPlatformChannelList(list []*share_message.PlatformChannel) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.PlatformChannelList = list
}
func (self *PayChannelManager) GetPlatformChannelList(t ...int32) []*share_message.PlatformChannel {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	tp := append(t, 0)[0]
	if tp != 0 {
		newList := make([]*share_message.PlatformChannel, 0)
		for _, c := range self.PlatformChannelList {
			if c.GetTypes() == tp {
				newList = append(newList, c)
			}
		}
		return newList
	}
	//全部
	return self.PlatformChannelList
}
func (self *PayChannelManager) SetGeneralQuota(data *share_message.GeneralQuota) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.GeneralQuota = data
}
func (self *PayChannelManager) GetGeneralQuota() *share_message.GeneralQuota {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return self.GeneralQuota
}

//获取渠道信息：id:渠道id
func (self *PayChannelManager) GetPlatformChannel(id int32) *share_message.PlatformChannel {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for _, v := range self.PlatformChannelList {
		if v.GetId() == id {
			return v
		}
	}
	return nil
}

//获取配置信息 id:渠道id
func (self *PayChannelManager) GetPaymentSettingList() []*share_message.PaymentSetting {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	list := []*share_message.PaymentSetting{}
	for _, channel := range self.PlatformChannelList {
		if channel != nil {
			for _, v := range self.PaymentSettingList {
				if v.GetId() == channel.GetPaymentSettingId() {
					if !easygo.Contain(list, v) {
						list = append(list, v)
					}
					break
				}
			}
		}
	}
	return list
}

//获取当前代付渠道
func (self *PayChannelManager) GetCurPayChannel() *share_message.PlatformChannel {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	for _, channel := range self.PlatformChannelList {
		if channel != nil {
			if channel.GetTypes() == 2 && channel.GetStatus() == 1 {
				//出款
				return channel
			}
		}
	}
	return nil
}

//获取当前渠道提现配置
func (self *PayChannelManager) GetCurPaymentSetting() *share_message.PaymentSetting {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	channel := self.GetCurPayChannel()
	for _, v := range self.PaymentSettingList {
		if channel.GetPaymentSettingId() == v.GetId() {
			//提现配置
			return v
		}
	}
	return nil
}

//获取支付渠道信息
//func (self *PayChannelManager) GetPayChannelData() []*share_message.PayData {
//	self.Mutex.Lock()
//	defer self.Mutex.Unlock()
//	var list []*share_message.PayData
//	for _, v := range self.PlatformChannelList {
//		pay := &share_message.PayData{
//			PayId:    easygo.NewInt32(v.PlatformId),
//			PayWay:   easygo.NewInt32(v.PayTypeId),
//			PaySence: easygo.NewInt32(v.PaySceneId),
//		}
//		list = append(list, pay)
//	}
//	return list
//}
