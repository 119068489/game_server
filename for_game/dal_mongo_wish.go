package for_game

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/h5_wish"
	"game_server/pb/share_message"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
)

//初始化500个男，女许愿用户表
func InitWishPlayer() []interface{} {
	var data []interface{}
	manHeadList := GetManyRobotHeadIcon(500, 1)  //男随机头像列表
	girlHeadList := GetManyRobotHeadIcon(500, 2) //女头像随机列表
	manNameList := GetManyRobotName(1, 500)
	girlNameList := GetManyRobotName(2, 500)
	for i := 0; i < 500; i++ {
		headMan := fmt.Sprintf("https://im-resource-1253887233.file.myqcloud.com/prod/mavatar/%d.png", manHeadList[i])
		headGirl := fmt.Sprintf("https://im-resource-1253887233.file.myqcloud.com/prod/wavatar/%d.png", girlHeadList[i])
		manPlayer := &share_message.WishPlayer{
			Id:         easygo.NewInt64(NextId(TABLE_WISH_PLAYER)),
			PlayerId:   easygo.NewInt64(0),
			NickName:   easygo.NewString(manNameList[i]),
			HeadUrl:    easygo.NewString(headMan),
			CreateTime: easygo.NewInt64(time.Now().Unix()),
			Types:      easygo.NewInt32(WISH_PLAYER_TYPE_ROBOT),
		}
		girlPlayer := &share_message.WishPlayer{
			Id:         easygo.NewInt64(NextId(TABLE_WISH_PLAYER)),
			PlayerId:   easygo.NewInt64(0),
			NickName:   easygo.NewString(girlNameList[i]),
			HeadUrl:    easygo.NewString(headGirl),
			CreateTime: easygo.NewInt64(time.Now().Unix()),
			Types:      easygo.NewInt32(WISH_PLAYER_TYPE_ROBOT),
		}
		data = append(data, manPlayer, girlPlayer)
	}
	return data
}

//处理许愿池相关db操作
func InitPlayerWishData() []interface{} {
	var data []interface{}
	data = append(data, &share_message.PlayerWishData{
		Id:            easygo.NewInt64(NextId(TABLE_PLAYER_WISH_DATA)),
		PlayerId:      easygo.NewInt64(888888),
		WishBoxId:     easygo.NewInt64(1),
		WishBoxItemId: easygo.NewInt64(1),
		Status:        easygo.NewInt32(0),
		CreateTime:    easygo.NewInt64(time.Now().Unix()),
	})
	data = append(data, &share_message.PlayerWishData{
		Id:            easygo.NewInt64(NextId(TABLE_PLAYER_WISH_DATA)),
		PlayerId:      easygo.NewInt64(888888),
		WishBoxId:     easygo.NewInt64(1),
		WishBoxItemId: easygo.NewInt64(2),
		Status:        easygo.NewInt32(0),
		CreateTime:    easygo.NewInt64(time.Now().Unix()),
	})
	return data
}

//初始化PlayerWishCollection
func InitPlayerWishCollection() []interface{} {
	var data []interface{}
	data = append(data, &share_message.PlayerWishCollection{
		Id:         easygo.NewInt64(NextId(TABLE_PLAYER_WISH_COLLECTION)),
		PlayerId:   easygo.NewInt64(888888),
		WishBoxId:  easygo.NewInt64(1),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
	})
	data = append(data, &share_message.PlayerWishCollection{
		Id:         easygo.NewInt64(NextId(TABLE_PLAYER_WISH_COLLECTION)),
		PlayerId:   easygo.NewInt64(888888),
		WishBoxId:  easygo.NewInt64(2),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
	})
	return data
}

//许愿池首页菜单
func InitWishMenu() []interface{} {
	var topicTypes []interface{}
	topicTypes = append(topicTypes, &share_message.WishMenu{
		Id:   easygo.NewInt64(1),
		Name: easygo.NewString("最新上线"),
	})
	topicTypes = append(topicTypes, &share_message.WishMenu{
		Id:   easygo.NewInt64(2),
		Name: easygo.NewString("人气盲盒"),
	})
	topicTypes = append(topicTypes, &share_message.WishMenu{
		Id:   easygo.NewInt64(3),
		Name: easygo.NewString("欧气爆棚"),
	})
	return topicTypes
}

//许愿池物品品牌
func InitWishBrand() []interface{} {
	var topicTypes []interface{}
	topicTypes = append(topicTypes, &share_message.WishBrand{
		Id:   easygo.NewInt64(1),
		Name: easygo.NewString("apple"),
	})
	topicTypes = append(topicTypes, &share_message.WishBrand{
		Id:   easygo.NewInt64(2),
		Name: easygo.NewString("联想"),
	})
	topicTypes = append(topicTypes, &share_message.WishBrand{
		Id:   easygo.NewInt64(3),
		Name: easygo.NewString("华为"),
	})
	topicTypes = append(topicTypes, &share_message.WishBrand{
		Id:   easygo.NewInt64(4),
		Name: easygo.NewString("小米"),
	})
	return topicTypes
}

//许愿池物品类型
func InitWishItemType() []interface{} {
	var topicTypes []interface{}
	topicTypes = append(topicTypes, &share_message.WishItemType{
		Id:   easygo.NewInt64(1),
		Name: easygo.NewString("手机"),
	})
	topicTypes = append(topicTypes, &share_message.WishItemType{
		Id:   easygo.NewInt64(2),
		Name: easygo.NewString("电脑"),
	})
	topicTypes = append(topicTypes, &share_message.WishItemType{
		Id:   easygo.NewInt64(3),
		Name: easygo.NewString("耳机"),
	})
	return topicTypes
}
func InitWishStyle() []interface{} {
	var topicTypes []interface{}
	topicTypes = append(topicTypes, &share_message.WishStyle{
		Id:   easygo.NewInt64(1),
		Name: easygo.NewString("普通款"),
	})
	topicTypes = append(topicTypes, &share_message.WishStyle{
		Id:   easygo.NewInt64(2),
		Name: easygo.NewString("典藏款"),
	})
	topicTypes = append(topicTypes, &share_message.WishStyle{
		Id:   easygo.NewInt64(3),
		Name: easygo.NewString("梦幻款"),
	})
	topicTypes = append(topicTypes, &share_message.WishStyle{
		Id:   easygo.NewInt64(4),
		Name: easygo.NewString("超凡款"),
	})
	return topicTypes
}

//许愿池物品初始化
func InitWishItem() []interface{} {
	var topicTypes []interface{}
	topicTypes = append(topicTypes, &share_message.WishItem{
		Id:            easygo.NewInt64(1),
		Name:          easygo.NewString("apple 6sp"),
		Icon:          easygo.NewString(""),
		Desc:          easygo.NewString("苹果6s plus"),
		Brand:         easygo.NewInt32(1),
		Type:          easygo.NewInt32(1),
		Price:         easygo.NewInt64(3000),
		RecoveryPrice: easygo.NewInt64(2400),
		BigSize:       easygo.NewString("5.5英寸"),
	})
	topicTypes = append(topicTypes, &share_message.WishItem{
		Id:            easygo.NewInt64(2),
		Name:          easygo.NewString("apple X"),
		Icon:          easygo.NewString(""),
		Desc:          easygo.NewString("苹果X"),
		Brand:         easygo.NewInt32(1),
		Type:          easygo.NewInt32(1),
		Price:         easygo.NewInt64(6000),
		RecoveryPrice: easygo.NewInt64(4800),
		BigSize:       easygo.NewString("5.5英寸"),
	})
	topicTypes = append(topicTypes, &share_message.WishItem{
		Id:            easygo.NewInt64(3),
		Name:          easygo.NewString("apple 12"),
		Icon:          easygo.NewString(""),
		Desc:          easygo.NewString("苹果12"),
		Brand:         easygo.NewInt32(1),
		Type:          easygo.NewInt32(1),
		Price:         easygo.NewInt64(8000),
		RecoveryPrice: easygo.NewInt64(6400),
		BigSize:       easygo.NewString("5.5英寸"),
	})
	topicTypes = append(topicTypes, &share_message.WishItem{
		Id:            easygo.NewInt64(4),
		Name:          easygo.NewString("Redmi 9A"),
		Icon:          easygo.NewString(""),
		Desc:          easygo.NewString("Redmi 9A"),
		Brand:         easygo.NewInt32(4),
		Type:          easygo.NewInt32(1),
		Price:         easygo.NewInt64(1000),
		RecoveryPrice: easygo.NewInt64(800),
		BigSize:       easygo.NewString("5.5英寸"),
	})
	return topicTypes
}

//盲盒商品初始化
func InitWishBoxItem() []interface{} {
	var topicTypes []interface{}
	topicTypes = append(topicTypes, &share_message.WishBoxItem{
		Id:         easygo.NewInt64(1),
		WishItemId: easygo.NewInt64(1),
		Status:     easygo.NewInt32(1),
		WishBoxId:  easygo.NewInt64(1),
		Style:      easygo.NewInt32(1),
	})
	topicTypes = append(topicTypes, &share_message.WishBoxItem{
		Id:         easygo.NewInt64(2),
		WishItemId: easygo.NewInt64(2),
		Status:     easygo.NewInt32(1),
		WishBoxId:  easygo.NewInt64(1),
		Style:      easygo.NewInt32(2),
	})
	topicTypes = append(topicTypes, &share_message.WishBoxItem{
		Id:         easygo.NewInt64(3),
		WishItemId: easygo.NewInt64(3),
		Status:     easygo.NewInt32(1),
		WishBoxId:  easygo.NewInt64(1),
		Style:      easygo.NewInt32(3),
	})
	topicTypes = append(topicTypes, &share_message.WishBoxItem{
		Id:         easygo.NewInt64(4),
		WishItemId: easygo.NewInt64(4),
		Status:     easygo.NewInt32(1),
		WishBoxId:  easygo.NewInt64(1),
		Style:      easygo.NewInt32(4),
	})

	topicTypes = append(topicTypes, &share_message.WishBoxItem{
		Id:         easygo.NewInt64(5),
		WishItemId: easygo.NewInt64(1),
		Status:     easygo.NewInt32(1),
		WishBoxId:  easygo.NewInt64(2),
		Style:      easygo.NewInt32(1),
	})
	topicTypes = append(topicTypes, &share_message.WishBoxItem{
		Id:         easygo.NewInt64(6),
		WishItemId: easygo.NewInt64(2),
		Status:     easygo.NewInt32(1),
		WishBoxId:  easygo.NewInt64(2),
		Style:      easygo.NewInt32(2),
	})
	topicTypes = append(topicTypes, &share_message.WishBoxItem{
		Id:         easygo.NewInt64(7),
		WishItemId: easygo.NewInt64(3),
		Status:     easygo.NewInt32(1),
		WishBoxId:  easygo.NewInt64(2),
		Style:      easygo.NewInt32(3),
	})
	topicTypes = append(topicTypes, &share_message.WishBoxItem{
		Id:         easygo.NewInt64(8),
		WishItemId: easygo.NewInt64(4),
		Status:     easygo.NewInt32(1),
		WishBoxId:  easygo.NewInt64(2),
		Style:      easygo.NewInt32(4),
	})
	return topicTypes
}

//盲盒初始化
func InitWishBox() []interface{} {
	var topicTypes []interface{}
	topicTypes = append(topicTypes, &share_message.WishBox{
		Id:         easygo.NewInt64(1),
		Name:       easygo.NewString("盲盒1"),
		Icon:       easygo.NewString(""),
		Menu:       []int32{0, 1},
		Items:      []int64{1, 2, 3, 4},
		Desc:       easygo.NewString("测试1盲盒"),
		Index:      easygo.NewString("盲盒1,apple 6sp,apple X,apple 12,Redmi 9A"),
		Match:      easygo.NewInt32(0),
		TotalNum:   easygo.NewInt32(4),
		RareNum:    easygo.NewInt32(3),
		Price:      easygo.NewInt64(1000),
		Status:     easygo.NewInt32(1),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
		Brands:     []int64{1, 4},
		Types:      []int64{1},
		Styles:     []int32{1, 2, 3, 4},
		WinNum:     easygo.NewInt64(0),
	})
	topicTypes = append(topicTypes, &share_message.WishBox{
		Id:         easygo.NewInt64(2),
		Name:       easygo.NewString("盲盒2"),
		Icon:       easygo.NewString(""),
		Menu:       []int32{2, 3},
		Items:      []int64{5, 6, 7, 8},
		Desc:       easygo.NewString("测试2盲盒"),
		Index:      easygo.NewString("盲盒2,apple 6sp,apple X,apple 12,Redmi 9A"),
		Match:      easygo.NewInt32(0),
		TotalNum:   easygo.NewInt32(4),
		RareNum:    easygo.NewInt32(3),
		Price:      easygo.NewInt64(1000),
		Status:     easygo.NewInt32(1),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
		Brands:     []int64{1, 4},
		Types:      []int64{1},
		Styles:     []int32{1, 2, 3, 4},
		WinNum:     easygo.NewInt64(0),
	})
	return topicTypes
}

//修改用户信息
func UpdatePlayerInfoSid(playerId int64, sid int32) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFun()
	err := col.Update(bson.M{"PlayerId": playerId}, bson.M{"$set": bson.M{"HallSid": sid}})
	if err != nil {
		easygo.PanicError(err)
	}
}

//获取用户信息
func GetWishPlayerInfo(playerId int64) *share_message.WishPlayer {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFun()
	var player *share_message.WishPlayer
	err := col.Find(bson.M{"_id": playerId}).One(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return player
}

//获取用户信息
func GetWishPlayerInfoByImId(playerId int64) *share_message.WishPlayer {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFun()
	var player *share_message.WishPlayer
	err := col.Find(bson.M{"PlayerId": playerId}).One(&player)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		easygo.PanicError(err)
	}
	return player
}

//创建一个许愿池用户信息
func CreatePlayerInfo(msg *h5_wish.LoginReq) (*share_message.WishPlayer, *base.Fail) {
	player := &share_message.WishPlayer{
		Id:         easygo.NewInt64(NextId(TABLE_WISH_PLAYER)),
		Account:    easygo.NewString(msg.GetAccount()),
		Channel:    easygo.NewInt32(msg.GetChannel()),
		NickName:   easygo.NewString(msg.GetNickName()),
		HeadUrl:    easygo.NewString(msg.GetHeadUrl()),
		PlayerId:   easygo.NewInt64(msg.GetPlayerId()),
		Token:      easygo.NewString(msg.GetToken()),
		IsTryOne:   easygo.NewBool(false),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
		Types:      easygo.NewInt32(msg.GetTypes()),
	}
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_WISH_PLAYER)
	defer closeFun()
	err := col.Insert(player)
	if err != nil {
		return nil, easygo.NewFailMsg(err.Error())
	}
	return player, nil
}

//许愿池物品类型
func InitWishRecycleReason() []interface{} {
	var reason []interface{}
	reason = append(reason, &share_message.RecycleReason{
		Id:     easygo.NewInt64(1),
		Reason: easygo.NewString("已有多件相似的相同物品"),
	})
	reason = append(reason, &share_message.RecycleReason{
		Id:     easygo.NewInt64(2),
		Reason: easygo.NewString("不喜欢"),
	})
	reason = append(reason, &share_message.RecycleReason{
		Id:     easygo.NewInt64(3),
		Reason: easygo.NewString("不想要"),
	})
	return reason
}
