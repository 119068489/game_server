package sport_apply

import (
	"fmt"
	dal "game_server/e-sports/sport_common_dal"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/client_hall"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"time"
)

func GetSportExChangeConfig(playId int64, rd *client_hall.ESPortsCoinViewResult) {
	//设置兑换配置
	exChangeConfigs := for_game.GetRedisEXChangeConfig()

	if nil != exChangeConfigs && len(exChangeConfigs) > 0 {
		rd.ExChangeList = exChangeConfigs
	}

	//============设置赠送类型  开始=======================
	//默认给予日常赠送
	changeType := for_game.ESPORT_EXCHANGE_TYPE_1
	rd.Type = easygo.NewInt32(changeType)

	//是否首充
	IsFirstExChange(playId, rd)

	//是否活动以及活动是否有效期
	activeConfig := for_game.GetRedisActiveConfig()
	nowTime := time.Now().Unix()

	if nil != activeConfig {
		//活动开启并且在活动期间
		if activeConfig.GetStatus() == 0 &&
			activeConfig.GetStartTime() <= nowTime &&
			nowTime < activeConfig.GetEndTime() {
			changeType = for_game.ESPORT_EXCHANGE_TYPE_3
			rd.Type = easygo.NewInt32(changeType)
		}
	}
	//==========设置赠送类型 结束========================

	//banner设置=================开始=======

	banners := for_game.GetRedisExChangeBanner()
	if nil != banners && len(banners) > 0 {
		if rd.GetType() == for_game.ESPORT_EXCHANGE_TYPE_1 {
			rd.BannerUrl = easygo.NewString(banners[for_game.ESPORT_EXCHANGE_TYPE_1])
		} else if rd.GetType() == for_game.ESPORT_EXCHANGE_TYPE_2 {
			rd.BannerUrl = easygo.NewString(banners[for_game.ESPORT_EXCHANGE_TYPE_2])
		} else if rd.GetType() == for_game.ESPORT_EXCHANGE_TYPE_3 {
			rd.BannerUrl = easygo.NewString(banners[for_game.ESPORT_EXCHANGE_TYPE_3])
		}
	}
	//banner设置=================结束=======
}

func IsFirstExChange(playId int64, rdView *client_hall.ESPortsCoinViewResult) bool {

	//是否首充
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_EXCHANGE_FIRST)
	defer closeFun()

	exChangeFirstQuery := share_message.TableESPortsExChangeFirst{}
	//通过条件查询
	errFirstQuery := col.Find(bson.M{"_id": playId}).One(&exChangeFirstQuery)

	if errFirstQuery != nil && errFirstQuery != mgo.ErrNotFound {
		logs.Error(errFirstQuery)
		s := fmt.Sprintf("=======IsFirstExChange========查询用户兑换首次充值数据失败=======条件:_id:%v", playId)
		logs.Error(s)
		rdView.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rdView.Msg = easygo.NewString("系统异常")
		return false
	}

	//是首冲
	if errFirstQuery == mgo.ErrNotFound {
		rdView.Type = easygo.NewInt32(for_game.ESPORT_EXCHANGE_TYPE_2)
		return true
	}

	if errFirstQuery == nil {
		return false
	}

	return false
}

//兑换后重新设置页面的最新值
func ResetESportExChangeView(playId int64, rd *client_hall.ESPortsCoinExChangeResult, rdView *client_hall.ESPortsCoinViewResult) {
	GetSportExChangeConfig(playId, rdView)
	rd.Type = easygo.NewInt32(rdView.GetType())
	rd.BannerUrl = easygo.NewString(rdView.GetBannerUrl())
	rd.ExChangeList = rdView.GetExChangeList()
}

//是否白名单
func IsESportExChangeWhite(playId int64) bool {
	//是否白名单
	one := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GIVE_WHITELIST, bson.M{"_id": playId})
	if one != nil {
		return true
	}
	return false
}

func InsFirstExChange(playId int64) {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_EXCHANGE_FIRST)
	defer closeFun()

	nowTime := time.Now().Unix()
	insObj := share_message.TableESPortsExChangeFirst{
		PlayerId:   easygo.NewInt64(playId),
		CreateTime: easygo.NewInt64(nowTime),
		UpdateTime: easygo.NewInt64(nowTime),
	}

	errIns := col.Insert(insObj)

	if errIns != nil {
		logs.Error(errIns)
		s := fmt.Sprintf("=======InsFirstExChange========插入首充记录失败=======insObj:%v", insObj)
		logs.Error(s)
	}
}

func ExChangeESportCoins(playId int64,
	eSportCoin int64,
	rd *client_hall.ESPortsCoinExChangeResult,
	rdView *client_hall.ESPortsCoinViewResult,
	exChangeFlag int32) *client_hall.ESPortsCoinExChangeResult {

	var msg string
	if exChangeFlag == for_game.ESPORTS_EXCHANGE_NORMAL {
		msg = fmt.Sprintf("兑换获得电竞币[%d]个", eSportCoin)
	} else if exChangeFlag == for_game.ESPORTS_EXCHANGE_WHITE {
		msg = fmt.Sprintf("兑换白名单赠送电竞币[%d]个", eSportCoin)
	} else if exChangeFlag == for_game.ESPORTS_EXCHANGE_FIRSRT {
		msg = fmt.Sprintf("兑换首充赠送电竞币[%d]个", eSportCoin)
	} else if exChangeFlag == for_game.ESPORTS_EXCHANGE_DAY {
		msg = fmt.Sprintf("兑换日常赠送电竞币[%d]个", eSportCoin)
	} else if exChangeFlag == for_game.ESPORTS_EXCHANGE_ACTIVE {
		msg = fmt.Sprintf("兑换活动赠送电竞币[%d]个", eSportCoin)
	}

	var st int32
	if exChangeFlag == for_game.ESPORTS_EXCHANGE_NORMAL {
		st = for_game.ESPORTCOIN_TYPE_EXCHANGE_IN
	} else {
		st = for_game.ESPORTCOIN_TYPE_EXCHANGE_GIVE_IN
	}

	//取得流水订单号
	streamOrderId := for_game.RedisCreateOrderNo(for_game.GOLD_CHANGE_TYPE_IN, st)
	req := &share_message.ESportCoinRecharge{
		PlayerId:     easygo.NewInt64(playId),
		RechargeCoin: easygo.NewInt64(eSportCoin),
		SourceType:   easygo.NewInt32(st),
		Note:         easygo.NewString(msg),
		ExtendLog: &share_message.GoldExtendLog{
			OrderId: easygo.NewString(streamOrderId), //流水的订单号
		},
	}

	result, err := dal.SendMsgToServerNewEx(PServerInfoMgr, playId, "RpcESportSendChangeESportCoins", req) //通知大厅
	if err != nil {
		logs.Error(err.GetReason())
		rd.Code = easygo.NewInt32(for_game.C_SYS_ERROR)
		rd.Msg = easygo.NewString(err.GetReason())
		//重置最新页面
		ResetESportExChangeView(playId, rd, rdView)
		return rd
	}
	if nil != result {
		rst2, ok2 := result.(*client_hall.ESportCommonResult)
		if ok2 && nil != rst2 {
			if rst2.GetCode() == for_game.C_SETTLEMENT_MONEY_FAIL {
				rd.Code = easygo.NewInt32(for_game.C_SETTLEMENT_MONEY_FAIL)
				rd.Msg = easygo.NewString(rst2.GetMsg())
				//重置最新页面
				ResetESportExChangeView(playId, rd, rdView)
				return rd
			}
		}
	}

	rd.Code = easygo.NewInt32(for_game.C_OPT_SUCCESS)
	rd.Msg = easygo.NewString("兑换成功")

	if exChangeFlag != for_game.ESPORTS_EXCHANGE_NORMAL {
		rd.GiveTotalCoins = easygo.NewInt64(eSportCoin)
	}
	//重置最新页面
	ResetESportExChangeView(playId, rd, rdView)

	return rd
}

//活动、通过页面的值算出要随机赠送的电竞币
func GetActiveESportCoins(exchangeRates []*share_message.ExchangeRate, esportCoins int64) int64 {
	if nil == exchangeRates || len(exchangeRates) <= 0 {
		return int64(0)
	}
	ratios := make([]int32, 0)
	for _, value := range exchangeRates {
		odds := value.GetOdds()
		if odds > 0 {
			var i int32
			for i = 0; i < odds; i++ {
				ratios = append(ratios, value.GetRatio())
			}
		}
	}

	//打乱ratios
	tempRatios := DoShuffleSilence(ratios)
	//随机取得数组中一个值
	randomRatio := tempRatios[rand.Intn(len(tempRatios))]

	if randomRatio > 0 {
		//计算电竞币

		tempCoins := esportCoins * int64(randomRatio) / int64(100)

		return tempCoins
	}
	return int64(0)
}

func DoShuffleSilence(ratiosParam []int32) []int32 {
	ratios := make([]int32, len(ratiosParam))
	copy(ratios, ratiosParam)
	rand.Shuffle(len(ratios), func(i int, j int) {
		ratios[i], ratios[j] = ratios[j], ratios[i]
	})
	return ratios
}
