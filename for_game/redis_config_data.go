package for_game

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
)

//配置数据存redis供后续使用
const (
	REDIS_CONFIG_INTIMACY      = "config:intimacy"      //亲密度配置
	REDIS_CONFIG_CONSTELLATION = "config:constellation" //星座

	REDIS_CONFIG_WISH_PAYMENT = "wish:config:payment"  //支付风控配置
	REDIS_CONFIG_WISH_PAYWARN = "wish:config:pay_warn" //支付回收预警
)

//初始化所有配置
func InitAllConfig() {
	InitConfigIntimacy()
	InitConfigConstellation()
	InitConfigWishPayment()
	InitConfigWishPayWarn()
}

//================================亲密度配置========================================start
//初始化亲密度redis配置
func InitConfigIntimacy() {
	configs := GetConfigIntimacyFormDB()
	SetConfigIntimacy(configs)
}

//设置亲密度
func SetConfigIntimacy(configs []*share_message.IntimacyConfig) {
	mConfig := make(map[int32]string)
	for _, info := range configs {
		s, _ := json.Marshal(info)
		mConfig[info.GetLv()] = string(s)
	}
	if len(mConfig) == 0 {
		return
	}
	err := easygo.RedisMgr.GetC().HMSet(REDIS_CONFIG_INTIMACY, mConfig)
	easygo.PanicError(err)
}

// 获取亲密度redis配置
func GetConfigIntimacy(lv int32) *share_message.IntimacyConfig {
	if lv > PLAYER_INTIMACY_MAX_LV {
		lv = PLAYER_INTIMACY_MAX_LV //目前只有5个等级，暂时这样处理
	}
	b, _ := easygo.RedisMgr.GetC().Exist(REDIS_CONFIG_INTIMACY)
	if !b {
		InitConfigIntimacy()
	}
	val, err := easygo.RedisMgr.GetC().HGet(REDIS_CONFIG_INTIMACY, easygo.AnytoA(lv))
	easygo.PanicError(err)
	var config *share_message.IntimacyConfig
	err = json.Unmarshal(val, &config)
	easygo.PanicError(err)
	return config
}

//================================亲密度配置========================================end
//================================星座配置========================================start
func InitConfigConstellation() {
	configs := GetConfigConstellationFormDB()
	SetConfigConstellation(configs)
}

//设置星座配置
func SetConfigConstellation(configs []*share_message.InterestTag) {
	mConfig := make(map[int32]string)
	for _, info := range configs {
		s, _ := json.Marshal(info)
		mConfig[info.GetId()] = string(s)
	}
	if len(mConfig) == 0 {
		return
	}
	err := easygo.RedisMgr.GetC().HMSet(REDIS_CONFIG_CONSTELLATION, mConfig)
	easygo.PanicError(err)
}

// 获取星座redis配置
func GetConfigConstellation(id int32) *share_message.InterestTag {
	b, _ := easygo.RedisMgr.GetC().Exist(REDIS_CONFIG_CONSTELLATION)
	if !b {
		InitConfigConstellation()
	}
	val, err := easygo.RedisMgr.GetC().HGet(REDIS_CONFIG_CONSTELLATION, easygo.AnytoA(id))
	easygo.PanicError(err)
	var config *share_message.InterestTag
	err = json.Unmarshal(val, &config)
	easygo.PanicError(err)
	return config
}

//获取星座名称
func GetConfigConstellationName(id int32) string {
	if id == 0 {
		return ""
	}
	b, _ := easygo.RedisMgr.GetC().Exist(REDIS_CONFIG_CONSTELLATION)
	if !b {
		InitConfigConstellation()
	}
	val, err := easygo.RedisMgr.GetC().HGet(REDIS_CONFIG_CONSTELLATION, easygo.AnytoA(id))
	if err != nil {
		logs.Error("不存在的星座信息id:", id)
		return ""
	}
	easygo.PanicError(err)
	var config *share_message.InterestTag
	err = json.Unmarshal(val, &config)
	easygo.PanicError(err)
	return config.GetName()
}

//获取星座简称
func GetConfigConstellationSortName(id int32) string {
	if id == 0 {
		return ""
	}
	b, _ := easygo.RedisMgr.GetC().Exist(REDIS_CONFIG_CONSTELLATION)
	if !b {
		InitConfigConstellation()
	}
	val, err := easygo.RedisMgr.GetC().HGet(REDIS_CONFIG_CONSTELLATION, easygo.AnytoA(id))
	if err != nil {
		logs.Error("不存在的星座信息id:", id)
		return ""
	}
	var config *share_message.InterestTag
	err = json.Unmarshal(val, &config)
	easygo.PanicError(err)
	return config.GetSortName()
}

//================================星座配置========================================end
//================================支付风控配置========================================start
//初始化支付风控配置
func InitConfigWishPayment() {
	configs := GetConfigWishPaymentFormDB()
	if configs == nil {
		return
	}
	SetConfigWishPayment(configs)

}

//设置支付风控配置
func SetConfigWishPayment(configs *share_message.WishRecycleSection) {
	mConfig := make(map[string]string)
	s, _ := json.Marshal(configs)
	mConfig["payment"] = string(s)
	err := easygo.RedisMgr.GetC().HMSet(REDIS_CONFIG_WISH_PAYMENT, mConfig)
	easygo.PanicError(err)
}

// 获取支付风控redis配置
func GetConfigWishPayment() *share_message.WishRecycleSection {
	b, _ := easygo.RedisMgr.GetC().Exist(REDIS_CONFIG_WISH_PAYMENT)
	if !b {
		InitConfigWishPayment()
	}
	val, err := easygo.RedisMgr.GetC().HGet(REDIS_CONFIG_WISH_PAYMENT, "payment")
	easygo.PanicError(err)
	var config *share_message.WishRecycleSection
	err = json.Unmarshal(val, &config)
	easygo.PanicError(err)
	return config
}

//================================支付风控配置========================================end

//================================支付预警配置========================================start
//初始化支付预警配置
func InitConfigWishPayWarn() {
	configs := GetWishPayWarnCfg()
	if configs == nil {
		return
	}
	SetConfigWishPayWarn(configs)
}

//设置支付预警配置
func SetConfigWishPayWarn(configs *brower_backstage.WishPayWarnCfg) {
	mConfig := make(map[string]string)
	s, _ := json.Marshal(configs)
	mConfig["pay_warn"] = string(s)
	err := easygo.RedisMgr.GetC().HMSet(REDIS_CONFIG_WISH_PAYWARN, mConfig)
	easygo.PanicError(err)
}

// 获取支付预警redis配置
func GetConfigWishPayWarn() *brower_backstage.WishPayWarnCfg {
	b, _ := easygo.RedisMgr.GetC().Exist(REDIS_CONFIG_WISH_PAYWARN)
	if !b {
		InitConfigWishPayWarn()
	}
	val, err := easygo.RedisMgr.GetC().HGet(REDIS_CONFIG_WISH_PAYWARN, "pay_warn")
	easygo.PanicError(err)
	var config *brower_backstage.WishPayWarnCfg
	err = json.Unmarshal(val, &config)
	easygo.PanicError(err)
	return config
}

//================================支付预警配置========================================end
