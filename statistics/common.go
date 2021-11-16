package statistics

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/for_game/greenScan"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"math"
	"time"

	"github.com/akqp2019/protobuf/proto"

	"github.com/akqp2019/mgo"

	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

type ITEM_ID = int64
type PLAYER_ID = int64
type ENDPOINT_ID = easygo.ENDPOINT_ID
type DB_NAME = string
type SERVER_ID = int32
type INSTANCE_ID = int32
type GAME_TYPE = int32
type SITE = string
type LEVEL = int32
type AliItemScanMap = map[ITEM_ID]*AliItemScan

const PLAYER_OFFLINE_TIME = 600 * 1000           //玩家如果离线超10分钟，则清理掉key
const UPDATA_DEL_PLAYER_KEYS = time.Second * 900 //15分检测下
const ACCOUNT_CANCEL_TIME = 86400 * 7 * 1000     //账号注销时间，毫秒 7天自动完成
//const ACCOUNT_CANCEL_TIME = 300 * 1000 //账号注销时间，毫秒 7天自动完成
//const ACCOUNT_CANCEL_TIME = 6 * 1000 //账号注销时间，毫秒10分钟

//商城超时半小时的数据取消
const ORDER_EXPIRE_TIME int64 = 1800

//存储桶定时时间
const COS_IMAGE_DELETE_TIME int64 = 600 //秒

//商城阿里验证出错后录入的类型
//目前只在发布商品的验证才录入
const (
	ALI_AUDIT_ORIGIN_1 = "1" //发布商品
	ALI_AUDIT_ORIGIN_2 = "2" //留言
	ALI_AUDIT_ORIGIN_3 = "3" //评价
	ALI_AUDIT_TYPE_1   = "1" //文本
	ALI_AUDIT_TYPE_2   = "2" //图片
	ALI_AUDIT_TYPE_3   = "3" //视频
)

//商城阿里验证内存结构体
type AliItemScan struct {
	item_id         ITEM_ID
	image_taskIds   []string
	image_task_flag int32
	image_task_time int64
	video_taskIds   []string
	video_task_flag int32
	video_task_time int64
}

//处理长时间未使用rediskeys
func DelRedisTimeOutkeys() {
	easygo.AfterFunc(UPDATA_DEL_PLAYER_KEYS, DelRedisTimeOutkeys)

}

//处理注销账号到期完成逻辑
func DealAccountCancel() {
	var orders []*share_message.PlayerCancleAccount
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CANCEL_ACCOUNT)
	defer closeFun()
	err := col.Find(bson.M{"Status": for_game.ACCOUNT_CANCEL_WAITING}).All(&orders)
	if err != nil && err != mgo.ErrNotFound {
		easygo.PanicError(err)
	}
	t := for_game.GetMillSecond()
	var saveList []interface{}
	for _, v := range orders {
		//过期了
		if v.GetStatus() == for_game.ACCOUNT_CANCEL_WAITING {
			if v.GetCreateTime()+ACCOUNT_CANCEL_TIME < t {
				//修改玩家状态
				logs.Info("注销玩家:", v.GetPlayerId(), v.GetPhone())
				p := for_game.GetRedisPlayerBase(v.GetPlayerId())
				if p == nil {
					continue
				}
				p.CancelAccountFinish()
				v.Status = easygo.NewInt32(for_game.ACCOUNT_CANCEL_FINISH)
				saveList = append(saveList, bson.M{"_id": v.GetId()}, v)
			}
		}
	}
	//如果有修改，则存储
	if len(saveList) > 0 {
		for_game.UpsertAll(MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_CANCEL_ACCOUNT, saveList)
	}
	easygo.AfterFunc(time.Second*600, DealAccountCancel)
	//easygo.AfterFunc(time.Second*200, DealAccountCancel)
}

//商城超时取消订单
func UpdateOrderList() {

	now := time.Now().Unix()
	//这里必须加一个函数结构体为了defer处理
	func() {

		col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ORDERS)
		defer closeFun()

		list := []*share_message.TableShopOrder{}

		e := col.Find(bson.M{"state": for_game.SHOP_ORDER_WAIT_PAY}).All(&list)

		if e != nil && e != mgo.ErrNotFound {
			logs.Error(e)
		}

		if e == nil && len(list) > 0 {

			for _, value := range list {

				if now > value.GetCreateTime()+ORDER_EXPIRE_TIME &&
					value.GetState() == for_game.SHOP_ORDER_WAIT_PAY {

					//这里必须加一个函数结构体为了defer处理以及分布式锁未取得退出函数做下条数据
					func() {

						//这里锁必须在函数内而且必须是在满足的条件内
						//以每个订单为单位取得锁,取不到直接做下次做
						lockKey := for_game.MakeRedisKey(for_game.SHOP_ORDER_WAIT_PAY_MUTEX, value.GetOrderId())
						//取得分布式锁,跟主动收货互斥（失效时间设置20秒）
						errLock := easygo.RedisMgr.GetC().DoRedisLockNoRetry(lockKey, 20)
						defer easygo.RedisMgr.GetC().DoRedisUnlock(lockKey)

						//如果未取得锁
						if errLock != nil {
							s := fmt.Sprintf("UpdateOrderList定时任务 单key取得订单redis分布式无重试锁失败,redis key is %v", lockKey)
							logs.Error(s)
							logs.Error(errLock)
							//直接退出函数进入下一次循环
							return
						}

						//取得订单对应的商品的分布式锁，此锁因为在定时中不需要重试(恢复库存的竞争)
						//取不到下次定时做不需要重试
						tempItemLockKey := for_game.MakeRedisKey(for_game.SHOP_ITEM_PAY_MUTEX, value.GetItems().GetItemId())
						//取得分布式锁开始1、取得订单对应的商品的分布式锁不需要重试
						errLock2 := easygo.RedisMgr.GetC().DoRedisLockNoRetry(tempItemLockKey, 10)
						defer easygo.RedisMgr.GetC().DoRedisUnlock(tempItemLockKey)

						//如果未取得锁就直接不做了
						if errLock2 != nil {
							s := fmt.Sprintf("UpdateOrderList定时任务 单key取得商品redis分布式无重试锁失败,redis key is %v", tempItemLockKey)
							logs.Error(s)
							logs.Error(errLock2)
							//直接退出函数进入下一次循环
							return
						}

						e := col.Update(
							bson.M{"_id": value.GetOrderId(), "state": for_game.SHOP_ORDER_WAIT_PAY},
							bson.M{"$set": bson.M{"state": for_game.SHOP_ORDER_EXPIRE}})
						if e != nil && e != mgo.ErrNotFound {
							logs.Error(e, value.GetOrderId())
							//直接退出函数进入下一次循环
							return
						}

						if e == mgo.ErrNotFound {
							logs.Error(e, value.GetOrderId())
							s := fmt.Sprintf("UpdateOrderList定时任务中 %v订单状态是待付款才能取消,不用管,属于正常的,直接做下条记录", value.GetOrderId())
							logs.Error(s)
							//直接退出函数进入下一次循环
							return
						}

						//判断是不是拉起了微信或支付宝支付后的订单但是没有在小程序点击取消的订单
						countCle, errCle := for_game.GetPayOrderListByShopOrderId(value.GetOrderId())
						if errCle == "" {
							//执行恢复库存,有几个支付订单就做几个
							for i := 0; i < countCle; i++ {
								//恢复库存并上架判断操作
								recoverStockErr := for_game.ShopRecoverStock(value.GetOrderId())
								if recoverStockErr != "" {
									logs.Error(recoverStockErr)
									//直接退出函数进入下一次循环
									return
								}
							}
						} else {
							logs.Error(errCle)
							//直接退出函数进入下一次循环
							return
						}
					}()

				}
			}
		}
	}()

	easygo.AfterFunc(time.Second*60, UpdateOrderList)
}

//商城阿里云验证发送
func DoAliItemScanSend() {

	var newMap AliItemScanMap = AliItemScanMap{}

	var list []*share_message.TableShopItem

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	e := col.Find(bson.M{"state": for_game.SHOP_ITEM_IN_AUDIT}).Sort("-create_time").All(&list)
	closeFun()

	if e == nil {
		if nil != list && len(list) != 0 {
			for _, value := range list {

				time.Sleep(time.Duration(500) * time.Millisecond)

				var aliScanMap *AliItemScan = (*ali_item_scan_by_id)[value.GetItemId()]
				if aliScanMap != nil {
					//判断该商品是否有图片或者视频需要重新验证，重新验证的加到新的map中二次确认
					if (aliScanMap.image_task_flag == 1 &&
						aliScanMap.image_taskIds != nil &&
						len(aliScanMap.image_taskIds) > 0) ||
						(aliScanMap.video_task_flag == 1 &&
							aliScanMap.video_taskIds != nil &&
							len(aliScanMap.video_taskIds) > 0) {
						newMap[value.GetItemId()] = (*ali_item_scan_by_id)[value.GetItemId()]
						continue
					}
				}

				itemFiles := value.GetItemFiles()
				//商品名称
				var itemNames []string = []string{value.GetName()}
				//正文描述
				var texts []string = []string{value.GetTitle()}
				//取得商品名称验证
				rstName, rstNameErrCode, rstNameErrContent := greenScan.GetTextScanResult(itemNames)
				//取得正文描述验证
				rstText, rstErrCode, rstErrContent := greenScan.GetTextScanResult(texts)

				if rstName == 0 {
					ItemTextSoldOut(value.GetItemId(), rstNameErrCode, rstNameErrContent, "商品名称")

					continue

				} else if rstName == 2 {

					//判断是否超过一小时
					if time.Now().Unix() >= (value.GetCreateTime() + 3600) {

						ItemTextSoldOut(value.GetItemId(), "9999", "商品名称审核超过一个小时", "商品名称")
					}
					continue
				} else if rstText == 0 {

					ItemTextSoldOut(value.GetItemId(), rstErrCode, rstErrContent, "商品描述")

					continue

				} else if rstText == 2 {

					//判断是否超过一小时
					if time.Now().Unix() >= (value.GetCreateTime() + 3600) {
						ItemTextSoldOut(value.GetItemId(), "9999", "描述审核超过一个小时", "商品描述")
					}

					continue

				} else {

					if nil != itemFiles && len(itemFiles) > 0 {

						var imageUrls []string = []string{}
						var videoUrls []string = []string{}

						for i := 0; i < len(itemFiles); i++ {
							var itemUrl string = itemFiles[i].GetFileUrl()
							var itemType int32 = itemFiles[i].GetFileType()
							if itemType == 0 {
								imageUrls = append(imageUrls, itemUrl)
							} else if itemType == 1 {
								videoUrls = append(videoUrls, itemUrl)
							}
						}

						var aliItemScan = AliItemScan{
							item_id: value.GetItemId(),
						}
						//图片异步发送
						if nil != imageUrls && len(imageUrls) > 0 {
							imageTaskIds, reqFlag := greenScan.GetImageScanTaskIds(imageUrls)

							if 0 == reqFlag {
								ItemSendSoldOut(value.GetItemId())

								continue

							} else {

								aliItemScan.image_taskIds = imageTaskIds
								aliItemScan.image_task_flag = 1
								aliItemScan.image_task_time = time.Now().Unix()
							}

						} else {
							aliItemScan.image_taskIds = nil
							aliItemScan.image_task_flag = 0
							aliItemScan.image_task_time = time.Now().Unix()
						}
						//视频异步发送
						if nil != videoUrls && len(videoUrls) > 0 {
							videoTaskIds, reqFlag := greenScan.GetVideoScanTaskIds(videoUrls)

							if 0 == reqFlag {
								ItemSendSoldOut(value.GetItemId())

								continue

							} else {

								aliItemScan.video_taskIds = videoTaskIds
								aliItemScan.video_task_flag = 1
								aliItemScan.video_task_time = time.Now().Unix()
							}

						} else {

							aliItemScan.video_taskIds = nil
							aliItemScan.video_task_flag = 0
							aliItemScan.video_task_time = time.Now().Unix()
						}

						newMap[value.GetItemId()] = &aliItemScan
					}

				}
			}
		}
	} else {
		logs.Error(e)
	}
	ali_item_scan_by_id = &newMap

	easygo.AfterFunc(time.Second*30, DoAliItemScanSend)
}

func GetImageVideoScanRst() {

	aliItemScanById := ali_item_scan_by_id

	if nil != aliItemScanById && len(*aliItemScanById) > 0 {

		for itemId := range *aliItemScanById {

			time.Sleep(time.Duration(500) * time.Millisecond)

			var imageTaskIds []string = (*aliItemScanById)[itemId].image_taskIds
			var imageTaskFlag int32 = (*aliItemScanById)[itemId].image_task_flag
			var imageTaskTime int64 = (*aliItemScanById)[itemId].image_task_time

			//判断图片验证
			if (nil == imageTaskIds || len(imageTaskIds) == 0) && imageTaskFlag == 1 {
				//图片结果还没有出来 等待下一次
				continue
			} else if (nil == imageTaskIds || len(imageTaskIds) == 0) && imageTaskFlag == 0 {
				//说明没有要图片验证,直接做视频验证
				var videoTaskIds []string = (*aliItemScanById)[itemId].video_taskIds
				var videoTaskFlag int32 = (*aliItemScanById)[itemId].video_task_flag
				var videoTaskTime int64 = (*aliItemScanById)[itemId].video_task_time

				//判断视频验证
				if (nil == videoTaskIds || len(videoTaskIds) == 0) && videoTaskFlag == 1 {
					//视频结果还没有出来 等待下一次
					continue
				} else if (nil == videoTaskIds || len(videoTaskIds) == 0) && videoTaskFlag == 0 {
					//说明没有要视频验证,直接上架
					col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
					e := col.Update(
						bson.M{"_id": itemId},
						bson.M{"$set": bson.M{"state": for_game.SHOP_ITEM_SALE}})
					closeFun()

					if e != nil {
						logs.Error(e)
					}
					//直接做下个商品
					continue
				}

				var nowTime int64 = time.Now().Unix()
				if (nowTime - videoTaskTime) >= 3600 {
					ItemVideoSoldOut(itemId, "9999", "视频审核超过一个小时")
					continue
				}

				rstVideo, rstErrCode, rstErrContent := greenScan.GetVideoRstByTaskIds(videoTaskIds)

				//视频审核未通过 下架
				if rstVideo == 0 {
					ItemVideoSoldOut(itemId, rstErrCode, rstErrContent)
					continue

					//下次继续审核
				} else if rstVideo == 2 {

					continue

					//图片,视频审核都通过,上架
				} else {

					col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
					e := col.Update(
						bson.M{"_id": itemId},
						bson.M{"$set": bson.M{"state": for_game.SHOP_ITEM_SALE}})
					closeFun()

					if e != nil {
						logs.Error(e)
					}

				}

				//做完视频验证直接做下个商品
				continue
			}

			//先判断图片验证是否超一个小时
			var nowTime int64 = time.Now().Unix()
			if (nowTime - imageTaskTime) >= 3600 {
				//图片审核超过一个小时直接下架处理
				ItemImageSoldOut(itemId, "9999", "图片审核超过一个小时")
				continue
			}

			rstImage, rstErrCode, rstErrContent := greenScan.GetImageRstByTaskIds(imageTaskIds)

			//图片审核未通过 下架
			if rstImage == 0 {
				ItemImageSoldOut(itemId, rstErrCode, rstErrContent)
				continue

				//下次继续审核
			} else if rstImage == 2 {

				continue

				//图片审核通过,继续做视频审核
			} else {

				var videoTaskIds []string = (*aliItemScanById)[itemId].video_taskIds
				var videoTaskFlag int32 = (*aliItemScanById)[itemId].video_task_flag
				var videoTaskTime int64 = (*aliItemScanById)[itemId].video_task_time

				//先判断图片验证
				if (nil == videoTaskIds || len(videoTaskIds) == 0) && videoTaskFlag == 1 {
					//视频结果还没有出来 等待下一次
					continue
				} else if (nil == videoTaskIds || len(videoTaskIds) == 0) && videoTaskFlag == 0 {
					//说明没有要视频验证,直接上架
					col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
					e := col.Update(
						bson.M{"_id": itemId},
						bson.M{"$set": bson.M{"state": for_game.SHOP_ITEM_SALE}})
					closeFun()

					if e != nil {
						logs.Error(e)
					}
					//直接做下个商品
					continue
				}

				//先判断是否超一个小时
				var nowTime int64 = time.Now().Unix()
				if (nowTime - videoTaskTime) >= 3600 {
					ItemVideoSoldOut(itemId, "9999", "视频审核超过一个小时")
					continue
				}

				rstVideo, rstErrCode, rstErrContent := greenScan.GetVideoRstByTaskIds(videoTaskIds)

				//视频审核未通过 下架
				if rstVideo == 0 {
					ItemVideoSoldOut(itemId, rstErrCode, rstErrContent)
					continue

					//下次继续审核
				} else if rstVideo == 2 {

					continue

					//图片,视频审核都通过,上架
				} else {

					col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
					e := col.Update(
						bson.M{"_id": itemId},
						bson.M{"$set": bson.M{"state": for_game.SHOP_ITEM_SALE}})
					closeFun()

					if e != nil {
						logs.Error(e)
					}

				}
			}
		}
	}
	easygo.AfterFunc(time.Second*20, GetImageVideoScanRst)
}

//阿里验证出错数据录入
func InsAliAuditErr(
	itemId *int64,
	origin *string,
	auditType *string,
	errCode *string,
	errContent *string,
	nowTime *int64) {

	aliFailInfo := share_message.TableShopAliAuditFail{
		ItemId:     itemId,
		Origin:     origin,
		Type:       auditType,
		ErrorCode:  errCode,
		Content:    errContent,
		CreateTime: nowTime,
	}

	colAudit, closeFunAudit := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ALI_AUDIT_FAIL)
	eAudit := colAudit.Insert(aliFailInfo)
	closeFunAudit()

	if eAudit != nil {
		logs.Error(eAudit)
	}
}

//阿里验证图片失败下架处理
func ItemImageSoldOut(
	itemId int64,
	errCode string,
	errContent string) {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	e := col.Update(
		bson.M{"_id": itemId},
		bson.M{"$set": bson.M{"state": for_game.SHOP_ITEM_FAIL_AUDIT, "sold_out_time": time.Now().Unix()}})
	closeFun()

	if e != nil {

		logs.Error(e)
	}

	if "" != errCode {
		s := fmt.Sprintf("图片审核失败-商品id:%v;错误码:%v;错误内容:%v", itemId, errCode, errContent)
		for_game.WriteFile("shop_audit.log", s)

		easygo.Spawn(func(itemId int64, errCode string, errContent string) {
			origin := ALI_AUDIT_ORIGIN_1
			auditType := ALI_AUDIT_TYPE_2
			nowTime := time.Now().Unix()

			InsAliAuditErr(easygo.NewInt64(itemId),
				easygo.NewString(origin),
				easygo.NewString(auditType),
				easygo.NewString(errCode),
				easygo.NewString(errContent),
				easygo.NewInt64(nowTime))

		}, itemId, errCode, errContent)
	}
}

//阿里验证视频失败下架处理
func ItemVideoSoldOut(
	itemId int64,
	errCode string,
	errContent string) {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	e := col.Update(
		bson.M{"_id": itemId},
		bson.M{"$set": bson.M{"state": for_game.SHOP_ITEM_FAIL_AUDIT, "sold_out_time": time.Now().Unix()}})
	closeFun()

	if e != nil {
		logs.Error(e)
	}

	if "" != errCode {

		s := fmt.Sprintf("视频审核失败-商品id:%v;错误码:%v;错误内容:%v", itemId, errCode, errContent)
		for_game.WriteFile("shop_audit.log", s)

		easygo.Spawn(func(itemId int64, errCode string, errContent string) {
			origin := ALI_AUDIT_ORIGIN_1
			auditType := ALI_AUDIT_TYPE_3
			nowTime := time.Now().Unix()

			InsAliAuditErr(easygo.NewInt64(itemId),
				easygo.NewString(origin),
				easygo.NewString(auditType),
				easygo.NewString(errCode),
				easygo.NewString(errContent),
				easygo.NewInt64(nowTime))

		}, itemId, errCode, errContent)
	}
}

//阿里验证文本失败下架处理
func ItemTextSoldOut(
	itemId int64,
	errCode string,
	errContent string,
	name string) {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	e := col.Update(
		bson.M{"_id": itemId},
		bson.M{"$set": bson.M{"state": for_game.SHOP_ITEM_FAIL_AUDIT, "sold_out_time": time.Now().Unix()}})
	closeFun()

	if e != nil {
		logs.Error(e)
	}

	if "" != errCode {
		s := fmt.Sprintf("审核失败-商品id:%v;错误码:%v;错误内容:%v",
			itemId,
			errCode,
			errContent)

		//把哪里调用的名称加上
		s = name + s

		for_game.WriteFile("shop_audit.log", s)

		easygo.Spawn(func(itemId int64, errCode string, errContent string) {

			origin := ALI_AUDIT_ORIGIN_1
			auditType := ALI_AUDIT_TYPE_1
			nowTime := time.Now().Unix()
			InsAliAuditErr(easygo.NewInt64(itemId),
				easygo.NewString(origin),
				easygo.NewString(auditType),
				easygo.NewString(errCode),
				easygo.NewString(errContent),
				easygo.NewInt64(nowTime))

		}, itemId, errCode, errContent)
	}
}

//阿里验证图片视频发送失败下架处理
func ItemSendSoldOut(itemId int64) {

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_SHOP_ITEMS)
	e := col.Update(
		bson.M{"_id": itemId},
		bson.M{"$set": bson.M{"state": for_game.SHOP_ITEM_FAIL_AUDIT, "sold_out_time": time.Now().Unix()}})
	closeFun()

	if e != nil {
		logs.Error(e)
	}
}

//腾讯云违禁图片从存储桶删除
func UpdateDelImageFromTX() {
	logs.Info("违规图片定时删除")
	col, closeFun := MongoLogMgr.GetC(for_game.MONGODB_NINGMENG_LOG, for_game.TABLE_PLAYER_TALK_LOG)
	defer closeFun()
	var logsInfo []*share_message.PlayerTalkLog
	t := time.Now().Unix()
	err := col.Find(bson.M{"EvilType": 20006, "CreateTime": bson.M{"$gte": t - COS_IMAGE_DELETE_TIME}}).All(&logsInfo)
	if err != nil && err != mgo.ErrNotFound {
		logs.Error("处理为空:", err)
	} else {
		logs.Info("logsInfo:", logsInfo)
		urls := []string{}
		for _, log := range logsInfo {
			urls = append(urls, log.GetConnect())
		}
		//TODO:通知存储桶删除指定url
		if len(urls) > 0 {
			for_game.DeleteMulti(urls)
		}
		logs.Info("处理违规图片:", urls)
	}
	easygo.AfterFunc(time.Second*time.Duration(COS_IMAGE_DELETE_TIME), UpdateDelImageFromTX)
}

//每小时衰减热门分值
func CoolingSquare() {
	logs.Info("整点热门分衰减")
	config := for_game.QuerySysParameterById(for_game.SQUAREHOT_PARAMETER)

	if config.GetDampRatio() <= 0 {
		return
	}

	list := for_game.GetHotDynamic()
	if len(list) == 0 {
		return
	}

	for _, li := range list {
		lessScore := -int32(math.Ceil(float64(li.GetHostScore()) * float64(config.GetDampRatio()) / float64(100)))
		// logs.Debug(fmt.Sprintf("原始分(%d),要减少的分(%d)", li.GetHostScore(), lessScore))
		find := bson.M{"_id": li.GetLogId()}
		update := bson.M{"$inc": bson.M{"HostScore": lessScore}}
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC, find, update, true)
		// 修改了动态,把redis删掉
		for_game.DelSquareDynamicById(li.GetLogId())
	}
}

//每小时更新热门话题
func UpdateHotTopic() {
	logs.Info("整点更新热门话题")
	queryBson := bson.M{"$and": []bson.M{{"HotCount": bson.M{"$ne": nil}}, {"HotCount": bson.M{"$gt": 0}}}}
	types, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC_TYPE, queryBson, 0, 0)
	config := for_game.QuerySysParameterById("topichot_parameter")
	var ids []int64
	for _, t := range types {
		if count, ok := t.(bson.M)["HotCount"]; ok {
			m := []bson.M{
				{"$match": bson.M{"TopicTypeId": t.(bson.M)["_id"]}},
				{"$addFields": bson.M{"Sort": bson.M{"$add": []string{"$ViewingNum", "$FansNum", "$ParticipationNum", "$AddViewingNum", "$AddParticipationNum", "$AddFansNum"}}}},
				{"$match": bson.M{"Sort": bson.M{"$gte": config.GetHotScore()}}},
				{"$sort": bson.M{"Sort": -1}},
				{"$project": bson.M{"Sort": 0}},
			}
			list := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC, m, count.(int), 0)

			for _, li := range list {
				id := li.(bson.M)["_id"].(int64)
				find := bson.M{"_id": id}
				update := bson.M{"$set": bson.M{"IsHot": true}}
				for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC, find, update, false)
				ids = append(ids, id)
			}
		}
	}

	falseList, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC, bson.M{"IsHot": true, "_id": bson.M{"$nin": ids}}, 0, 0)
	for _, fli := range falseList {
		falseFind := bson.M{"_id": fli.(bson.M)["_id"]}
		falseupdate := bson.M{"$set": bson.M{"IsHot": false}}
		for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_TOPIC, falseFind, falseupdate, false)
	}
}

//服务器间通讯通用
func SendMsgToServerNew(sid SERVER_ID, methodName string, msg easygo.IMessage, pid ...int64) (easygo.IMessage, *base.Fail) {
	srv := PServerInfoMgr.GetServerInfo(sid)
	if srv == nil {
		return nil, easygo.NewFailMsg("无效的服务器id =" + easygo.AnytoA(sid))
		//
	}
	var msgByte []byte
	if msg != nil {
		b, err := msg.Marshal()
		easygo.PanicError(err)
		msgByte = b
	} else {
		msgByte = []byte{}
	}

	playerId := append(pid, 0)[0]
	req := &share_message.MsgToServer{
		PlayerId: easygo.NewInt64(playerId),
		RpcName:  easygo.NewString(methodName),
		MsgName:  easygo.NewString(proto.MessageName(msg)),
		Msg:      msgByte,
	}
	return PWebApiForServer.SendToServer(srv, "RpcMsgToOtherServer", req)
}

//广播给指定类型服务器
func BroadCastMsgToServerNew(t int32, methodName string, msg easygo.IMessage, pid ...int64) {
	servers := PServerInfoMgr.GetAllServers(t)
	for _, srv := range servers {
		if srv == nil {
			continue
		}
		var msgByte []byte
		if msg != nil {
			b, err := msg.Marshal()
			easygo.PanicError(err)
			msgByte = b
		} else {
			msgByte = []byte{}
		}

		playerId := append(pid, 0)[0]
		req := &share_message.MsgToServer{
			PlayerId: easygo.NewInt64(playerId),
			RpcName:  easygo.NewString(methodName),
			MsgName:  easygo.NewString(proto.MessageName(msg)),
			Msg:      msgByte,
		}
		PWebApiForServer.SendToServer(srv, "RpcMsgToOtherServer", req)
	}
}

//发送给指定大厅指定玩家发送消息
func SendMsgToHallClientNew(playerIds []int64, methodName string, msg easygo.IMessage) {
	serversInfo := make(map[int32][]int64)
	for _, pid := range playerIds { //群发 每个人都发
		player := for_game.GetRedisPlayerBase(pid)
		if player == nil {
			continue
		}
		serversInfo[player.GetSid()] = append(serversInfo[player.GetSid()], pid)
	}
	for sid, pList := range serversInfo {
		srv := PServerInfoMgr.GetServerInfo(sid)
		if srv == nil {
			continue
		}
		var msgByte []byte
		if msg != nil {
			b, err := msg.Marshal()
			easygo.PanicError(err)
			msgByte = b
		} else {
			msgByte = []byte{}
		}

		req := &share_message.MsgToClient{
			PlayerIds: pList,
			RpcName:   easygo.NewString(methodName),
			MsgName:   easygo.NewString(proto.MessageName(msg)),
			Msg:       msgByte,
		}
		PWebApiForServer.SendToServer(srv, "RpcMsgToHallClient", req)
	}
}

//绑定硬币过期检测
func UpdateBCoinExpiration() {
	logs.Info("UpdateBCoinExpiration ----》》")
	var logs []*share_message.PlayerBCoinLog
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BCOIN_LOG)
	defer closeFun()
	t := time.Now().Unix()
	err := col.Find(bson.M{"Status": for_game.BCOIN_STATUS_UNUSE, "OverTime": bson.M{"$lt": t}}).Sort("OverTime").All(&logs)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	//处理到期绑定硬币
	DealBCoinExpiration(logs, col)
}

//处理过期硬币:扣除掉指定绑定硬币
func DealBCoinExpiration(data []*share_message.PlayerBCoinLog, col *mgo.Collection) {
	if len(data) <= 0 {
		return
	}
	logs.Info("处理过期硬币", data)
	total := int64(0)
	var saveLogs []interface{}
	playerBCoin := make(map[int64]int64)
	for _, log := range data {
		if log.GetCurBCoin() > 0 {
			total += log.GetCurBCoin()
			playerBCoin[log.GetPlayerId()] += log.GetCurBCoin()
			log.Status = easygo.NewInt32(for_game.BCOIN_STATUS_EXPIRATION)
			log.CurBCoin = easygo.NewInt64(0)
			saveLogs = append(saveLogs, bson.M{"_id": log.GetId()}, log)
		}
	}
	for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BCOIN_LOG, saveLogs)
	bb, _ := json.Marshal(playerBCoin)
	rpcReq := &server_server.NotifyAddCoinReq{
		NotifyAddCoin: easygo.NewString(string(bb)),
	}
	sInfo := PServerInfoMgr.GetIdelServer(easygo.SERVER_TYPE_HALL)
	logs.Debug("sInfo------->", sInfo)
	if sInfo == nil {
		logs.Error("服务不存在")
		return
	}
	sid := sInfo.GetSid()
	//  通知大厅去修改.
	SendMsgToServerNew(sid, "RpcNotifyAddCoin", rpcReq)

}

//处理1天内将要过期硬币:给玩家发送推送通知
func DealPreBCoinExpiration() {
	logs.Info("处理1天内将要过期硬币:给玩家发送推送通知")
	var exLogs []*share_message.PlayerBCoinLog
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BCOIN_LOG)
	defer closeFun()
	t := time.Now().Unix()
	err := col.Find(bson.M{"Status": for_game.BCOIN_STATUS_UNUSE, "OverTime": bson.M{"$gt": t, "$lt": t + 86400}}).Sort("OverTime").All(&exLogs)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}

	if len(exLogs) <= 0 {
		return
	}
	logs.Info("处理即将过期硬币", exLogs)
	playerBCoin := make(map[int64]int64)
	for _, log := range exLogs {
		playerBCoin[log.GetPlayerId()] += log.GetCurBCoin()
	}
	bb, _ := json.Marshal(playerBCoin)
	rpcReq := &server_server.NotifyAddCoinReq{
		NotifyAddCoin: easygo.NewString(string(bb)),
	}
	//  通知大厅.
	sInfo := PServerInfoMgr.GetIdelServer(easygo.SERVER_TYPE_HALL)
	if sInfo == nil {
		logs.Error("服务不存在")
		return
	}
	SendMsgToServerNew(sInfo.GetSid(), "RpcNoticeAssistant", rpcReq)

}

//处理下架商城限时售卖的商品
func UpdateMallItemSaleStatus() {
	logs.Info("每10分钟更新虚拟商城商品售卖状态")
	nowTime := util.GetMilliTime()
	queryBson := bson.M{"Status": for_game.COIN_PRODUCT_STATUS_UP, "SaleEndTime": bson.M{"$lte": nowTime, "$gt": 0}}
	list, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_PRODUCT, queryBson, 0, 0)
	var ids []int64
	for _, li := range list {
		ids = append(ids, li.(bson.M)["_id"].(int64))
	}

	updateFind := bson.M{"_id": bson.M{"$in": ids}}
	updateBson := bson.M{"$set": bson.M{"Status": for_game.COIN_PRODUCT_STATUS_DOWN}}
	for_game.UpdateAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_COIN_PRODUCT, updateFind, updateBson)
}

// 查询任务身上快过期的物品
func GetPlayerExpProduct() []*share_message.PlayerBagItem {
	start := time.Now().Unix()       // 提前一天通知
	end := time.Now().Unix() + 86400 // 提前一天通知
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BAG_ITEM)
	defer closeFun()
	list := make([]*share_message.PlayerBagItem, 0)
	err := col.Find(bson.M{"Status": bson.M{"$ne": 3}, "OverTime": bson.M{"$lte": end, "$gte": start}}).All(&list)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return list
}

// 通知物品过期
func NoticeProductExp() {
	logs.Info("NoticeProductExp 过期物品检测")
	productList := GetPlayerExpProduct()
	m := make(map[int64]string)
	for _, v := range productList {
		propsName := v.GetPropsName()
		if v1, ok := m[v.GetPlayerId()]; ok {
			propsName = fmt.Sprintf("%s、%s", v1, propsName)
		}
		m[v.GetPlayerId()] = propsName
	}
	if len(m) == 0 {
		return
	}
	// 执行推送
	bb, _ := json.Marshal(m)
	rpcReq := &server_server.NotifyAddCoinReq{
		NotifyAddCoin: easygo.NewString(string(bb)),
	}

	//  通知大厅.
	sInfo := PServerInfoMgr.GetIdelServer(easygo.SERVER_TYPE_HALL)
	if sInfo == nil {
		logs.Error("服务不存在")
		return
	}
	SendMsgToServerNew(sInfo.GetSid(), "RpcNoticeProductExp", rpcReq)

}

//玩家亲密度定时处理
func UpdatePlayerIntimacy() {
	logs.Info("UpdatePlayerIntimacy 处理超过3天以上没联系的好友亲密度")
	//先把redis数据存储到mongo
	for_game.SavePlayerIntimacyoMongoDB()
	//查找有亲密度值，并且超过3天没联系的亲密度数据，进行减少对应值
	intimacys := for_game.GetExpirationPlayerIntimacy()
	logs.Info("UpdatePlayerIntimacy 记录条数:", len(intimacys))
	saveData := make([]interface{}, 0)
	for _, intimacy := range intimacys {
		obj := for_game.GetRedisPlayerIntimacyObj(intimacy.GetId())
		if obj == nil {
			continue
		}
		obj.PerDayReduce()
		data := obj.GetRedisPlayerIntimacy()
		saveData = append(saveData, bson.M{"_id": data.GetId()}, data)
	}
	//批量存储修改的数据
	if len(saveData) > 0 {
		for_game.UpsertAll(MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_INTIMACY, saveData)
	}
	logs.Info("UpdatePlayerIntimacy 处理完成")
}

//每分钟处理过期话题内置顶
func UpdateTopicTopExp() {
	m := []bson.M{
		{"$project": bson.M{"TopicTopSet": 1}},
		{"$unwind": "$TopicTopSet"},
		{"$match": bson.M{"TopicTopSet.IsTopicTop": true, "TopicTopSet.TopicTopOverTime": bson.M{"$lt": easygo.NowTimestamp() * 1000}}},
	}
	lis := for_game.FindPipeAll(for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC, m, 0, 0)
	var data []interface{}
	for _, v := range lis {
		b1 := bson.M{"_id": v.(bson.M)["_id"], "TopicTopSet.TopicId": v.(bson.M)["TopicTopSet"].(bson.M)["TopicId"]}
		b2 := bson.M{"$set": bson.M{"TopicTopSet.$.IsTopicTop": false}}
		data = append(data, b1, b2)
		for_game.DelSquareDynamicById(v.(bson.M)["_id"].(int64)) //删除redis数据
	}
	for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_SQUARE_DYNAMIC, data)
	easygo.AfterFunc(1*time.Minute, UpdateTopicTopExp)
}

//定时给话题群组发送动态
func UpdateTopicTeamDynamic() {

}
