package backstage

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"log"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"

	"fmt"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
)

//查询系统管理员列表
func GetManagerList(user *share_message.Manager, reqMsg *brower_backstage.GetPlayerListRequest) ([]*share_message.Manager, int) {
	var list []*share_message.Manager
	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER)
	defer closeFun()

	queryBson := bson.M{}
	// 判断有日期才按日期查询
	if reqMsg.GetBeginTimestamp() != 0 && reqMsg.GetEndTimestamp() != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": reqMsg.GetBeginTimestamp(), "$lte": reqMsg.GetEndTimestamp()}
	}

	switch reqMsg.GetListType() {
	case 1:
		queryBson["Status"] = 0
	case 2:
		queryBson["Status"] = 1
	case 3:
		queryBson["IsOnlie"] = true
	case 4:
		queryBson["IsOnlie"] = false
	}

	switch reqMsg.GetType() {
	case 1:
		queryBson["Account"] = easygo.If(reqMsg.GetKeyword() != "", reqMsg.GetKeyword(), bson.M{"$ne": nil})
	case 2:
		queryBson["RealName"] = easygo.If(reqMsg.GetKeyword() != "", reqMsg.GetKeyword(), bson.M{"$ne": nil})
	case 3:
		i, _ := strconv.Atoi(reqMsg.GetKeyword()) //搜索查询不需要返回错误
		queryBson["_id"] = easygo.If(i != 0, i, bson.M{"$ne": nil})
	}
	if user.GetAccount() != "admin" {
		queryBson["Account"] = bson.M{"$ne": "admin"}
	}

	if reqMsg.GetRole() < 2 {
		queryBson["Role"] = bson.M{"$ne": 2}
	} else {
		queryBson["Role"] = reqMsg.GetRole()
	}

	if reqMsg.WaiterType != nil && reqMsg.GetWaiterType() != 0 {
		queryBson["Types"] = reqMsg.GetWaiterType()
	}

	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//修改管理员
func EditManage(site string, reqMsg *share_message.Manager, et string) {
	if et == "edit" && reqMsg.GetPassword() != "" {
		reqMsg.Password = easygo.NewString(for_game.CreatePasswd(reqMsg.GetPassword(), reqMsg.GetSalt()))
	}
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//创建管理员
func AddManage(siteId string, reqMsg *share_message.Manager) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER)
	defer closeFun()

	salt := for_game.RandString(8)
	reqMsg.Id = easygo.NewInt64(10000 + for_game.NextId(for_game.TABLE_MANAGER))
	reqMsg.Password = easygo.NewString(for_game.CreatePasswd(reqMsg.GetPassword(), salt))
	reqMsg.CreateTime = easygo.NewInt64(time.Now().Unix())
	reqMsg.Status = easygo.NewInt32(0)
	reqMsg.Salt = easygo.NewString(salt)
	reqMsg.Site = easygo.NewString(siteId)
	reqMsg.IsGoogleVer = easygo.NewBool(false)
	if reqMsg.Role == nil {
		reqMsg.Role = easygo.NewInt32(1)
	}

	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//帐号查询管理员
func QueryManage(account string) *share_message.Manager {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER)
	defer closeFun()

	siteOne := &share_message.Manager{}
	err := col.Find(bson.M{"Account": account}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//通过ID查询管理员
func QueryManageByID(id int64) *share_message.Manager {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER)
	defer closeFun()

	siteOne := &share_message.Manager{}
	err := col.Find(bson.M{"_id": id}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//通过昵称查询管理员
func QueryManageByName(realname string) *share_message.Manager {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER)
	defer closeFun()

	siteOne := &share_message.Manager{}
	err := col.Find(bson.M{"RealName": realname}).One(siteOne)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return siteOne
}

//查询客服分类列表
func GetManagerTypesList(reqMsg *brower_backstage.ListRequest) ([]*share_message.ManagerTypes, int) {

	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER_TYPES)
	defer closeFun()

	queryBson := bson.M{}
	queryBson["Name"] = easygo.If(reqMsg.GetKeyword() != "", reqMsg.GetKeyword(), bson.M{"$ne": nil})
	queryBson["Status"] = easygo.If(reqMsg.GetStatus() > 0, reqMsg.GetStatus(), bson.M{"$ne": nil})

	var list []*share_message.ManagerTypes
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	errc := query.Sort("-_id").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

//客服分类下拉列表
func GetManagerTypesListNopage() []*brower_backstage.KeyValueTag {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER_TYPES)
	defer closeFun()

	queryBson := bson.M{}
	query := col.Find(queryBson)
	var list []*share_message.ManagerTypes
	err := query.All(&list)
	easygo.PanicError(err)

	var lis []*brower_backstage.KeyValueTag

	for _, i := range list {
		li := &brower_backstage.KeyValueTag{
			//Key:   easygo.NewString(easygo.IntToString(int(i.GetId()))),
			Key:   i.Id,
			Value: i.Name,
		}
		lis = append(lis, li)
	}

	return lis
}

//修改客服分类
func EditManageTypes(reqMsg *share_message.ManagerTypes) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER_TYPES)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

//冻结解冻管理员
func UpAdminStatus(adminid []int32, status int32) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER)
	defer closeFun()
	_, err := col.UpdateAll(bson.M{"_id": bson.M{"$in": adminid}}, bson.M{"$set": bson.M{"Status": status}})
	easygo.PanicError(err)
}

//给角色权限表插入admin，客服，运营的初始化数据
func InitRolePower(InsertData []*share_message.RolePower) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ROLEPOWER)
	defer closeFun()
	query := col.Find(bson.M{})
	count, err := query.Count()
	easygo.PanicError(err)
	if count == 0 {
		var il []interface{}
		for _, rr := range InsertData {
			rr.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_ROLEPOWER))
			il = append(il, rr)
		}

		err := col.Insert(il...) //批量插入到数据库
		easygo.PanicError(err)

		s := "初始化角色权限表"
		logs.Info(s)
	}

}

// 查询角色列表
func QueryRolePowerList(reqMsg *brower_backstage.ListRequest) ([]*share_message.RolePower, int) {
	pageSize := int(reqMsg.GetPageSize())
	curPage := easygo.If(int(reqMsg.GetCurPage()) > 1, int(reqMsg.GetCurPage())-1, 0).(int)

	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ROLEPOWER)
	defer closeFun()

	queryBson := bson.M{}
	query := col.Find(queryBson)
	count, err := query.Count()
	easygo.PanicError(err)
	var list []*share_message.RolePower
	errc := query.Sort("-CreateTime").Skip(curPage * pageSize).Limit(pageSize).All(&list)
	easygo.PanicError(errc)

	return list, count
}

// Id查询角色
func QueryRolePowerById(site SITE, Id int32) *share_message.RolePower {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_ROLEPOWER)
	defer closeFun()
	rp := &share_message.RolePower{}

	errc := col.Find(bson.M{"_id": Id}).One(rp)
	if errc != nil && errc != mgo.ErrNotFound {
		panic(errc)
	}
	if errc == mgo.ErrNotFound {
		return nil
	}
	return rp
}

// Id查询角色是否拥有权限
func QueryPermissionById(site SITE, Id int32, permission string) bool {
	role := QueryRolePowerById(site, Id)
	if role == nil {
		logs.Error("無效的授权账号")
		return false
	}
	for _, menu_id := range role.GetMenuIds() {
		if menu_id == permission {
			return true
		}
	}

	return false
}

//检查角色名称
func CheckRoleName(site string, reqMsg *share_message.RolePower) int {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_ROLEPOWER)
	defer closeFun()
	queryBson := bson.M{}
	reqId := reqMsg.GetId()
	queryBson["RoleName"] = reqMsg.GetRoleName()
	if reqId != 0 {
		queryBson["_id"] = bson.M{"$ne": reqId}
	}
	count, err := col.Find(queryBson).Count()
	easygo.PanicError(err)
	return count
}

// 判断权限角色是否在被使用
func CheckAuthGroupByRole(ids []int64) string {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER)
	defer closeFun()

	queryBson := bson.M{}

	var list []*share_message.Manager
	roleName := ""
	queryBson["RoleType"] = bson.M{"$in": ids}
	err := col.Find(queryBson).All(&list)
	easygo.PanicError(err)

	if len(list) > 0 {
		for _, i := range list {
			role := GetPowerRouter(i.GetRoleType())
			roleName += fmt.Sprintf("%s ", role.GetRoleName())
		}

		roleName = fmt.Sprintf("%s-角色正在被管理账号使用，无法删除", roleName)
	}

	return roleName

}

//查询管理员列表下拉配置列表
func GetRolePowerList() []*brower_backstage.KeyValueTag {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ROLEPOWER)
	defer closeFun()
	var roles []*share_message.RolePower
	err := col.Find(bson.M{}).Select(bson.M{"_id": 1, "RoleName": 1}).All(&roles)
	easygo.PanicError(err)
	var list []*brower_backstage.KeyValueTag
	for _, item := range roles {
		l := &brower_backstage.KeyValueTag{
			Key:   easygo.NewInt32(item.GetId()),
			Value: item.RoleName,
		}

		list = append(list, l)
	}

	return list
}

// 更新角色权限
func UpdateRolePwer(site string, reqMsg *share_message.RolePower) {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_ROLEPOWER)
	defer closeFun()
	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})

	easygo.PanicError(err)
}

//获取角色权限
func GetPowerRouter(role int32) *share_message.RolePower {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_ROLEPOWER)
	defer closeFun()
	var result *share_message.RolePower
	err := col.Find(bson.M{"_id": role}).One(&result)
	if err == mgo.ErrNotFound {
		return nil
	}
	return result

}

//检查客服活跃消息数量并更新
func CheckWaiterCount() {
	Waiters := for_game.GetRedisWaiterList()
	for _, v := range Waiters {
		user := GetUser(v.UserId)
		count := GetActiveIMmessageCount(user)           //查询活跃消息数量
		for_game.ReloadRedisWaiterCount(v.UserId, count) //修改客服接待数量

		log.Println("=============每整10分钟更新客服接待数:", for_game.GetRedisWaiter(v.UserId))
	}
}
