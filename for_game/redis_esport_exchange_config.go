package for_game

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

const (
	//电竞币兑换配置参数 key(常驻redis的key)
	ESPORT_REDIS_EXCHANGE_CONFIG = "redis_esport:exchange_config"

	//电竞币活动参数 key(常驻redis的key)
	ESPORT_REDIS_EXCHANGE_ACTIVE_CONFIG = "redis_esport:exchange_active_config"

	//电竞币兑换banner参数 key(常驻redis的key)
	ESPORT_REDIS_EXCHANGE_BANNER_CONFIG = "redis_esport:exchange_banner_config"
)

//设置redis兑换banner
func SetRedisExChangeBanner() {

	banners := GetDBExChangeBanner()

	if nil != banners && len(banners) > 0 {

		//json序列化
		data, errJs := json.Marshal(&banners)
		easygo.PanicError(errJs)
		//redis 常驻key
		err := easygo.RedisMgr.GetC().Set(ESPORT_REDIS_EXCHANGE_BANNER_CONFIG, string(data))
		easygo.PanicError(err)
	} else {
		keys := []interface{}{ESPORT_REDIS_EXCHANGE_BANNER_CONFIG}

		easygo.RedisMgr.GetC().Delete(keys...)
	}
}

//取得redis兑换banner
func GetRedisExChangeBanner() map[int32]string {

	b, err := easygo.RedisMgr.GetC().Exist(ESPORT_REDIS_EXCHANGE_BANNER_CONFIG)
	easygo.PanicError(err)

	var banners map[int32]string
	if !b {
		dbBanners := GetDBExChangeBanner()

		if nil != dbBanners && len(dbBanners) > 0 {

			//json序列化
			data, errJs := json.Marshal(&dbBanners)
			easygo.PanicError(errJs)
			err := easygo.RedisMgr.GetC().Set(ESPORT_REDIS_EXCHANGE_BANNER_CONFIG, string(data))
			easygo.PanicError(err)

			banners = dbBanners
		}
	} else {
		obj := make(map[int32]string, 0)
		value, err := easygo.RedisMgr.GetC().Get(ESPORT_REDIS_EXCHANGE_BANNER_CONFIG)
		easygo.PanicError(err)
		errJs := json.Unmarshal([]byte(value), &obj)
		easygo.PanicError(errJs)

		banners = obj
	}
	return banners
}

//数据库兑换banner
func GetDBExChangeBanner() map[int32]string {

	//从数据库中取得数据
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ESPORTS_CAROUSEL)
	defer closeFun()

	bannerMap := make(map[int32]string, 0)

	carousel1 := share_message.TableESPortsCarousel{}
	err1 := col.Find(bson.M{"Types": ESPORT_EXCHANGE_TYPE_1, "Status": ESPORTS_BANNER_STATUS_1, "MenuId": ESPORTMENU_SHOP}).Sort("-Weight").Limit(1).One(&carousel1)
	if err1 != nil && err1 != mgo.ErrNotFound {
		logs.Error(err1)
		return bannerMap
	}

	if err1 == nil {
		bannerMap[ESPORT_EXCHANGE_TYPE_1] = carousel1.GetImageUrl()
	}

	carousel2 := share_message.TableESPortsCarousel{}
	err2 := col.Find(bson.M{"Types": ESPORT_EXCHANGE_TYPE_2, "Status": ESPORTS_BANNER_STATUS_1, "MenuId": ESPORTMENU_SHOP}).Sort("-Weight").Limit(1).One(&carousel2)
	if err2 != nil && err2 != mgo.ErrNotFound {
		logs.Error(err2)
		return bannerMap
	}

	if err2 == nil {
		bannerMap[ESPORT_EXCHANGE_TYPE_2] = carousel2.GetImageUrl()
	}

	carousel3 := share_message.TableESPortsCarousel{}
	err3 := col.Find(bson.M{"Types": ESPORT_EXCHANGE_TYPE_3, "Status": ESPORTS_BANNER_STATUS_1, "MenuId": ESPORTMENU_SHOP}).Sort("-Weight").Limit(1).One(&carousel3)
	if err3 != nil && err3 != mgo.ErrNotFound {
		logs.Error(err3)
		return bannerMap
	}

	if err3 == nil {
		bannerMap[ESPORT_EXCHANGE_TYPE_3] = carousel3.GetImageUrl()
	}

	return bannerMap

}

//数据库活动参数
func GetDBActiveConfig() *share_message.Activity {

	//从数据库中取得数据
	one := FindOne(MONGODB_NINGMENG, TABLE_LUCKY_ACTIVITY, bson.M{"Types": 2})
	activeConfig := &share_message.Activity{}
	if one != nil {
		StructToOtherStruct(one, activeConfig)
	} else {
		activeConfig = nil
	}

	return activeConfig
}

//设置redis活动的配置参数
func SetRedisActiveConfig() {

	activeConfig := GetDBActiveConfig()

	if nil != activeConfig {

		//json序列化
		data, errJs := json.Marshal(activeConfig)
		easygo.PanicError(errJs)
		err := easygo.RedisMgr.GetC().Set(ESPORT_REDIS_EXCHANGE_ACTIVE_CONFIG, string(data))
		easygo.PanicError(err)
	} else {
		keys := []interface{}{ESPORT_REDIS_EXCHANGE_ACTIVE_CONFIG}

		easygo.RedisMgr.GetC().Delete(keys...)
	}
}

//取得redis活动的配置参数
func GetRedisActiveConfig() *share_message.Activity {

	b, err := easygo.RedisMgr.GetC().Exist(ESPORT_REDIS_EXCHANGE_ACTIVE_CONFIG)
	easygo.PanicError(err)

	var activeConfig *share_message.Activity
	if !b {
		dbActiveConfig := GetDBActiveConfig()

		if nil != dbActiveConfig {

			//json序列化
			data, errJs := json.Marshal(dbActiveConfig)
			easygo.PanicError(errJs)
			err := easygo.RedisMgr.GetC().Set(ESPORT_REDIS_EXCHANGE_ACTIVE_CONFIG, string(data))
			easygo.PanicError(err)

			activeConfig = dbActiveConfig
		}

	} else {
		obj := &share_message.Activity{}
		value, err := easygo.RedisMgr.GetC().Get(ESPORT_REDIS_EXCHANGE_ACTIVE_CONFIG)
		easygo.PanicError(err)
		errJs := json.Unmarshal([]byte(value), obj)
		easygo.PanicError(errJs)

		activeConfig = obj
	}
	return activeConfig
}

//数据库兑换配置参数
func GetDBEXChangeConfigs() []*client_hall.ExChangeObject {

	//从数据库中取得数据
	dbEXChangeConfigs := make([]*share_message.TableESportsExchangeCfg, 0)
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_ESPORTS_EXCHANGE_CFG)
	defer closeFun()

	var rtList []*client_hall.ExChangeObject

	err := col.Find(bson.M{}).Sort("_id").All(&dbEXChangeConfigs)

	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return rtList
	}

	if err == mgo.ErrNotFound {
		return rtList
	}

	if err == nil && nil != dbEXChangeConfigs && len(dbEXChangeConfigs) > 0 {
		rtList = make([]*client_hall.ExChangeObject, 0)

		for _, value := range dbEXChangeConfigs {
			rtList = append(rtList, &client_hall.ExChangeObject{
				ESportCoin: easygo.NewInt64(value.GetId()),
				Coin:       easygo.NewInt64(value.GetCoin()),
				DailyGive:  easygo.NewInt64(value.GetDailyGive()),
				FirstGive:  easygo.NewInt64(value.GetFirstGive()),
				Rate:       value.GetRate(),
			})
		}
	}

	return rtList
}

//设置redis兑换配置参数
func SetRedisEXChangeConfigs() {

	dbEXChangeConfigs := GetDBEXChangeConfigs()

	if nil != dbEXChangeConfigs && len(dbEXChangeConfigs) > 0 {

		//json序列化
		data, errJs := json.Marshal(&dbEXChangeConfigs)
		easygo.PanicError(errJs)
		err := easygo.RedisMgr.GetC().Set(ESPORT_REDIS_EXCHANGE_CONFIG, string(data))
		easygo.PanicError(err)
	} else {
		keys := []interface{}{ESPORT_REDIS_EXCHANGE_CONFIG}

		easygo.RedisMgr.GetC().Delete(keys...)
	}
}

//取得redis兑换配置参数
func GetRedisEXChangeConfig() []*client_hall.ExChangeObject {

	b, err := easygo.RedisMgr.GetC().Exist(ESPORT_REDIS_EXCHANGE_CONFIG)
	easygo.PanicError(err)

	var exChangeConfigs []*client_hall.ExChangeObject
	if !b {
		dbEXChangeConfigs := GetDBEXChangeConfigs()

		if nil != dbEXChangeConfigs && len(dbEXChangeConfigs) > 0 {

			//json序列化
			data, errJs := json.Marshal(&dbEXChangeConfigs)
			easygo.PanicError(errJs)
			err := easygo.RedisMgr.GetC().Set(ESPORT_REDIS_EXCHANGE_CONFIG, string(data))
			easygo.PanicError(err)

			exChangeConfigs = dbEXChangeConfigs
		}

	} else {
		obj := make([]*client_hall.ExChangeObject, 0)
		value, err := easygo.RedisMgr.GetC().Get(ESPORT_REDIS_EXCHANGE_CONFIG)
		easygo.PanicError(err)
		errJs := json.Unmarshal([]byte(value), &obj)
		easygo.PanicError(errJs)

		exChangeConfigs = obj
	}
	return exChangeConfigs
}
