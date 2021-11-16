package for_game

// 盲盒数据

type RedisWishBoxObj struct {
	Id int64 //玩家id
	RedisBase
}

/*
//写入redis
func NewRedisPlayerBagItem(id PLAYER_ID, items ...[]*share_message.PlayerBagItem) *RedisPlayerBagItemObj {
	if id == 0 {
		return nil
	}
	p := &RedisPlayerBagItemObj{
		Id: id,
	}
	obj := append(items, nil)[0]
	return p.Init(obj)
}

func (self *RedisPlayerBagItemObj) Init(obj []*share_message.PlayerBagItem) *RedisPlayerBagItemObj {
	self.RedisBase.Init(self, self.Id, easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_BAG_ITEM)
	self.Sid = PlayerBagItemMgr.GetSid()
	if self.IsExistKey() {
		//redis已经存在key了，但管理器已销毁
		PlayerBagItemMgr.Store(self.Id, self)
		self.AddToExistList(self.Id)
		//启动定时任务
		easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	} else {
		if obj == nil {
			//没有去数据库查询
			obj = self.QueryPlayerBagItem(self.Id)
			//if obj == nil {
			//	return nil
			//}
		}
		self.SetRedisPlayerBagItem(obj)
	}
	return self
}
func (self *RedisPlayerBagItemObj) GetId() interface{} { //override
	return self.Id
}
func (self *RedisPlayerBagItemObj) GetKeyId() string { //override
	return MakeRedisKey(TABLE_PLAYER_BAG_ITEM, self.Id)
}

//重写保存方法
func (self *RedisPlayerBagItemObj) SaveToMongo() {
	items := self.GetRedisPlayerBagItem()
	if len(items) > 0 {
		saveList := make([]*share_message.PlayerBagItem, 0)
		for _, m := range items {
			if m.GetIsSave() {
				m.IsSave = easygo.NewBool(false)
				saveList = append(saveList, m)
				self.UpdateItem(m)
			}
		}
		if len(saveList) > 0 {
			var data []interface{}
			for _, v := range saveList {
				b1 := bson.M{"_id": v.GetId()}
				b2 := v
				data = append(data, b1, b2)
			}
			UpsertAll(self.DB, self.DBName, self.TBName, data)
		}
		self.SetSaveStatus(false)
	}
}

//定时更新数据
func (self *RedisPlayerBagItemObj) UpdateData() { //override
	if !self.IsExistKey() {
		PlayerBagItemMgr.Delete(self.Id) // 释放对象
		self.DelToExistList(self.Id)
		return
	}
	if self.GetSaveStatus() { //需要保存的数据进行存储
		self.SaveToMongo()
	}
	t := GetMillSecond()
	//存活10分钟，用到重新拉取
	if t-self.CreateTime > REDIS_PLAYER_BAG_ITEM_EXIST_TIME { //单位：毫秒
		if self.CheckIsDelRedisKey() {
			self.DelToExistList(self.Id)
			self.DelRedisKey() //redis删除
		}
		PlayerBagItemMgr.Delete(self.Id) // 释放对象
	}
	easygo.AfterFunc(REDIS_SAVE_TIME, self.UpdateData)
}
func (self *RedisPlayerBagItemObj) InitRedis() { //override
	obj := self.QueryPlayerBagItem(self.Id)
	if obj == nil {
		return
	}
	self.SetRedisPlayerBagItem(obj)
}
func (self *RedisPlayerBagItemObj) GetRedisSaveData() interface{} { //override
	data := self.GetRedisPlayerBagItem()
	return data
}

func (self *RedisPlayerBagItemObj) SaveOtherData() { //override
}
func (self *RedisPlayerBagItemObj) QueryPlayerBagItem(id int64) []*share_message.PlayerBagItem {
	var data []*share_message.PlayerBagItem
	col, closeFun := self.DB.GetC(self.DBName, self.TBName)
	defer closeFun()
	//过滤过期的道具
	err := col.Find(bson.M{"PlayerId": id}).All(&data)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return data
}
func (self *RedisPlayerBagItemObj) SetRedisPlayerBagItem(items []*share_message.PlayerBagItem) {
	PlayerBagItemMgr.Store(self.Id, self)
	self.AddToExistList(self.Id)
	//重置过期时间
	self.CreateTime = GetMillSecond()
	//启动定时任务
	easygo.AfterFunc(REDIS_SAVE_TIME, self.Me.UpdateData)
	if self.IsExistKey() {
		//如果数据已经存在redis，直接返回
		return
	}
	itemInfo := make(map[int64]string)
	for _, info := range items {
		s, _ := json.Marshal(info)
		itemInfo[info.GetId()] = string(s)
	}
	if len(itemInfo) == 0 {
		return
	}
	err2 := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), itemInfo)
	easygo.PanicError(err2)
}
func (self *RedisPlayerBagItemObj) GetRedisPlayerBagItem(t ...int32) []*share_message.PlayerBagItem {
	items := make([]*share_message.PlayerBagItem, 0)
	m, err := Int64StringMap(easygo.RedisMgr.GetC().HGetAll(self.GetKeyId()))
	easygo.PanicError(err)
	pos := append(t, 0)[0]
	for _, s := range m {
		var obj *share_message.PlayerBagItem
		_ = json.Unmarshal([]byte(s), &obj)
		//检测道具是否过期
		obj = self.CheckItemExpired(obj)
		if pos == 0 {
			items = append(items, obj)
		} else if pos != 0 && obj.GetPropsType() == pos {
			//obj.IsNew = easygo.NewBool(false)
			items = append(items, obj)
			//self.UpdateItem(obj)
		}
	}
	return items
}

//通过网络id获取道具
func (self *RedisPlayerBagItemObj) GetItemNetId(netId int64) *share_message.PlayerBagItem {
	if netId == 0 {
		return nil
	}
	b, err := easygo.RedisMgr.GetC().HGet(self.GetKeyId(), easygo.AnytoA(netId))
	if err != nil {
		logs.Error("获取道具失败:", netId)
		return nil
	}
	var item *share_message.PlayerBagItem
	//item := &share_message.PlayerBagItem{}
	err = json.Unmarshal(b, &item)
	easygo.PanicError(err)
	//检测道具是否过期了
	item = self.CheckItemExpired(item)
	return item
}

//检测道具是否过期，过期要做过期处理
func (self *RedisPlayerBagItemObj) CheckItemExpired(item *share_message.PlayerBagItem) *share_message.PlayerBagItem {
	if item.GetOverTime() != COIN_PROPS_FOREVER && time.Now().Unix() > item.GetOverTime() && item.GetStatus() != COIN_BAG_ITEM_EXPIRED {
		if item.GetStatus() == COIN_BAG_ITEM_USED {
			//如果道具使用中，要从装备栏卸下
			equipmentObj := GetRedisPlayerEquipmentObj(self.Id)
			equipmentObj.EquipmentDown(item.GetPropsType())
		}
		//修改道具属性
		item.Status = easygo.NewInt32(COIN_BAG_ITEM_EXPIRED)
		item.IsSave = easygo.NewBool(true)
		self.UpdateItem(item)
	}
	return item
}
func (self *RedisPlayerBagItemObj) GetPropsIdByNetId(netId int64) int64 {
	item := self.GetItemNetId(netId)
	return item.GetPropsId()
}

//通过配置id获取背包道具
func (self *RedisPlayerBagItemObj) GetItemPropsId(propsId int64, isForever ...bool) *share_message.PlayerBagItem {
	items := self.GetRedisPlayerBagItem()
	forever := append(isForever, false)[0]
	for _, item := range items {
		if item.GetPropsId() == propsId {
			if forever && item.GetOverTime() == COIN_PROPS_FOREVER {
				return item
			} else if !forever && item.GetOverTime() != COIN_PROPS_FOREVER {
				return item
			}
		}
	}
	return nil
}

//判断道具是否存在
func (self *RedisPlayerBagItemObj) IsExistItem(propsId int64) bool {
	item := self.GetItemPropsId(propsId)
	return item == nil
}

//批量增加道具
func (self *RedisPlayerBagItemObj) AddItems(items []*share_message.CoinProduct, way, bugWay int32, orderId string, operator string, givePlayerId ...int64) []*share_message.PlayerBagItem {
	newItems := make([]*share_message.PlayerBagItem, 0)
	for _, item := range items {
		bagItem := self.AddItem(item, way, bugWay, orderId, operator, givePlayerId...)
		if bagItem != nil {
			newItems = append(newItems, bagItem...)
		}
	}
	//写入到redis
	itemInfo := make(map[int64]string)
	for _, info := range newItems {
		s, _ := json.Marshal(info)
		itemInfo[info.GetId()] = string(s)
	}
	if len(itemInfo) == 0 {
		return nil
	}
	err2 := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), itemInfo)
	easygo.PanicError(err2)
	self.SetSaveSid()
	self.SetSaveStatus(true)
	return newItems
}

//单个增加道具
func (self *RedisPlayerBagItemObj) AddItem(item *share_message.CoinProduct, way, bugWay int32, orderId string, operator string, givePlayerId ...int64) []*share_message.PlayerBagItem {
	gPid := append(givePlayerId, 0)[0]
	propsItem := GetPropsItemInfo(item.GetPropsId())
	if propsItem == nil {
		logs.Error("找不到商品的配置Id:", item.GetPropsId())
		return nil
	}
	if item.GetProductNum() == 0 {
		logs.Error("怎么可以购买数量0的物品呢:", item.GetPropsId())
		return nil
	}
	newItems := make([]*share_message.PlayerBagItem, 0)

	//增加礼包型道具，直接使用
	if propsItem.GetPropsType() == COIN_PROPS_TYPE_LB {
		values := propsItem.GetUseValue()
		ids := make([]int64, 0)
		err := json.Unmarshal([]byte(values), &ids)
		easygo.PanicError(err)
		for _, id := range ids {
			it := GetPropsItemInfo(id)
			newItem := self.CreateItemExt(item, it, way)
			newItem.IsNew = easygo.NewBool(true)
			newItem.IsSave = easygo.NewBool(true)
			//TODO 增加道具Log
			data := &share_message.PlayerGetPropsLog{
				Id:            easygo.NewInt64(NextId(TABLE_PLAYER_GETPROPS_LOG)),
				PlayerId:      easygo.NewInt64(self.Id),
				GivePlayerId:  easygo.NewInt64(gPid),
				PropsId:       easygo.NewInt64(it.GetId()),
				PropsNum:      easygo.NewInt64(item.GetProductNum()),
				GetType:       easygo.NewInt32(way),
				CreateTime:    easygo.NewInt64(easygo.NowTimestamp()),
				EffectiveTime: easygo.NewInt64(item.GetEffectiveTime()),
				BagId:         easygo.NewInt64(newItem.GetId()),
				BuyWay:        easygo.NewInt32(bugWay),
				OrderId:       easygo.NewString(orderId),
				ProductId:     easygo.NewInt64(item.GetId()),
				Operator:      easygo.NewString(operator),
			}
			easygo.Spawn(AddGetPropsLog, data)
			newItems = append(newItems, newItem)
		}
	} else {
		newItem := self.CreateItemExt(item, propsItem, way)
		newItem.IsNew = easygo.NewBool(true)
		newItem.IsSave = easygo.NewBool(true)
		//TODO 增加道具Log
		data := &share_message.PlayerGetPropsLog{
			Id:            easygo.NewInt64(NextId(TABLE_PLAYER_GETPROPS_LOG)),
			PlayerId:      easygo.NewInt64(self.Id),
			GivePlayerId:  easygo.NewInt64(gPid),
			PropsId:       easygo.NewInt64(propsItem.GetId()),
			PropsNum:      easygo.NewInt64(item.GetProductNum()),
			GetType:       easygo.NewInt32(way),
			CreateTime:    easygo.NewInt64(easygo.NowTimestamp()),
			EffectiveTime: easygo.NewInt64(item.GetEffectiveTime()),
			BagId:         easygo.NewInt64(newItem.GetId()),
			BuyWay:        easygo.NewInt32(bugWay),
			OrderId:       easygo.NewString(orderId),
			ProductId:     easygo.NewInt64(item.GetId()),
			Operator:      easygo.NewString(operator),
		}
		easygo.Spawn(AddGetPropsLog, data)
		newItems = append(newItems, newItem)
	}
	return newItems
}
func (self *RedisPlayerBagItemObj) CreateItemExt(item *share_message.CoinProduct, propsItem *share_message.PropsItem, way int32) *share_message.PlayerBagItem {
	if item.GetEffectiveTime() == COIN_PROPS_FOREVER {
		//购买永久道具
		newItem := self.GetItemPropsId(propsItem.GetId(), true)
		if newItem != nil {
			//被回收的永久道具
			if newItem.GetStatus() == COIN_BAG_ITEM_EXPIRED {
				newItem.Status = easygo.NewInt32(COIN_BAG_ITEM_UNUSE)
			}
			newItem.CreateTime = easygo.NewInt64(time.Now().Unix()) //过期道具续费，重置购买时间
			return newItem
		}
		return self.CreateItem(item, propsItem, way)
	} else {
		//购买时效道具
		newItem := self.GetItemPropsId(propsItem.GetId())
		if newItem != nil {
			addTime := item.GetEffectiveTime() * item.GetProductNum() * 86400
			if newItem.GetStatus() == COIN_BAG_ITEM_EXPIRED {
				t := time.Now().Unix()
				newItem.OverTime = easygo.NewInt64(t + addTime)
				newItem.Status = easygo.NewInt32(COIN_BAG_ITEM_UNUSE)
				newItem.CreateTime = easygo.NewInt64(t) //过期道具续费，重置购买时间
			} else {
				newItem.OverTime = easygo.NewInt64(newItem.GetOverTime() + addTime)
			}
			return newItem
		}
		return self.CreateItem(item, propsItem, way)
	}
}

//创建一个新的背包道具
func (self *RedisPlayerBagItemObj) CreateItem(item *share_message.CoinProduct, propsItem *share_message.PropsItem, way int32) *share_message.PlayerBagItem {
	t := time.Now().Unix()
	if item.GetEffectiveTime() == COIN_PROPS_FOREVER {
		t = COIN_PROPS_FOREVER
	} else {
		//单位是天
		t += item.GetEffectiveTime() * item.GetProductNum() * 86400
	}
	newItem := &share_message.PlayerBagItem{
		Id:         easygo.NewInt64(NextId(TABLE_PLAYER_BAG_ITEM)),
		PlayerId:   easygo.NewInt64(self.Id),             // 玩家id
		PropsId:    easygo.NewInt64(propsItem.GetId()),   // 道具id
		Status:     easygo.NewInt32(COIN_BAG_ITEM_UNUSE), // 状态; 1-待使用;2-使用中;3-已使用完毕.4-未使用但失效
		GetType:    easygo.NewInt32(way),                 // 获得类型;1-购买;2-赠送(做任务获得)
		OverTime:   easygo.NewInt64(t),                   // 道具过期时间
		IsSave:     easygo.NewBool(true),                 //是否存储
		PropsType:  easygo.NewInt32(propsItem.GetPropsType()),
		PropsName:  easygo.NewString(propsItem.GetName()),
		CreateTime: easygo.NewInt64(time.Now().Unix()),
	}
	return newItem
}

//使用道具
func (self *RedisPlayerBagItemObj) SetItemStatus(id int64, st int32) {
	item := self.GetItemNetId(id)
	if item == nil {
		return
	}
	item.Status = easygo.NewInt32(st)
	item.IsSave = easygo.NewBool(true)
	//item.IsNew = easygo.NewBool(false) // 装上没有new标签
	self.UpdateItem(item)
}

//检测是否已经购买了永久道具
func (self *RedisPlayerBagItemObj) CheckHadForeverItem(ids []int64) []string {
	names := make([]string, 0)
	items := self.GetRedisPlayerBagItem()
	for _, id := range ids {
		for _, item := range items {
			if id == item.GetPropsId() && item.GetOverTime() == COIN_PROPS_FOREVER && item.GetStatus() != COIN_BAG_ITEM_EXPIRED {
				names = append(names, item.GetPropsName())
			}
		}
	}
	return names
}

//更新redis数据
func (self *RedisPlayerBagItemObj) UpdateItem(item *share_message.PlayerBagItem) {
	itemInfo := make(map[int64]string)
	s, _ := json.Marshal(item)
	itemInfo[item.GetId()] = string(s)
	err2 := easygo.RedisMgr.GetC().HMSet(self.GetKeyId(), itemInfo)
	easygo.PanicError(err2)
	self.SetSaveSid()
	self.SetSaveStatus(true)
}

//道具减少使用时间
func (self *RedisPlayerBagItemObj) ReduceUseTime(netId, t int64) *share_message.PlayerBagItem {
	item := self.GetItemNetId(netId)
	if item == nil {
		return nil
	}
	if t == COIN_PROPS_FOREVER {
		//永久道具，直接删除
		if item.GetStatus() == COIN_BAG_ITEM_USED {
			//如果道具使用中，要从装备栏卸下
			equipmentObj := GetRedisPlayerEquipmentObj(self.Id)
			equipmentObj.EquipmentDown(item.GetPropsType())
		}
		item.Status = easygo.NewInt32(COIN_BAG_ITEM_EXPIRED)
	} else {
		lt := item.GetOverTime() - t
		if lt < time.Now().Unix() {
			//道具过期
			if item.GetStatus() == COIN_BAG_ITEM_USED {
				//如果道具使用中，要从装备栏卸下
				equipmentObj := GetRedisPlayerEquipmentObj(self.Id)
				equipmentObj.EquipmentDown(item.GetPropsType())
			}
			item.Status = easygo.NewInt32(COIN_BAG_ITEM_EXPIRED)
			item.OverTime = easygo.NewInt64(0)
		} else {
			item.OverTime = easygo.NewInt64(lt)
		}

	}
	item.IsSave = easygo.NewBool(true)
	self.UpdateItem(item)
	return item
}

//检测是否有新道具,同步装备
func (self *RedisPlayerBagItemObj) CheckNewItemNotice() []int32 {
	newList := []int32{}
	equipments := make(map[int32]int64)
	items := self.GetRedisPlayerBagItem()
	for _, item := range items {
		if item.GetIsNew() {
			newList = append(newList, item.GetPropsType())
		}
		if item.GetStatus() == COIN_BAG_ITEM_USED {
			equipments[item.GetPropsType()] = item.GetId()
		}
	}
	//同步装备信息
	if len(equipments) > 0 {
		equipmentObj := GetRedisPlayerEquipmentObj(self.Id)
		for t, id := range equipments {
			equipmentObj.Equipment(t, id)
		}
	}
	return newList
}

//对外接口
func GetRedisPlayerBagItemObj(id int64, data ...[]*share_message.PlayerBagItem) *RedisPlayerBagItemObj {
	return PlayerBagItemMgr.GetRedisPlayerBagItemObj(id, data...)
}

//停服保存处理，保存需要存储的数据
func SaveRedisPlayerBagItemToMongo() {
	ids := []int64{}
	GetAllRedisSaveList(TABLE_PLAYER_BAG_ITEM, &ids)
	items := make([]*share_message.PlayerBagItem, 0)
	for _, id := range ids {
		obj := GetRedisPlayerBagItemObj(id)
		if obj != nil {
			data := obj.GetRedisPlayerBagItem()
			items = append(items, data...)
			obj.SetSaveStatus(false)
		}
	}
	if len(items) > 0 {
		saveData := make([]interface{}, 0)
		for _, it := range items {
			saveData = append(saveData, bson.M{"_id": it.GetId()}, it)
		}
		UpsertAll(easygo.MongoMgr, MONGODB_NINGMENG, TABLE_PLAYER_BAG_ITEM, saveData)
	}
}
*/
