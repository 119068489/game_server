package sport_common_dal

import (
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
	"time"
)

func AddCarousel(menuid int32, title, ImageUrl string) (int32, string) {
	info := &share_message.TableESPortsCarousel{
		Id:          nil,
		Title:       easygo.NewString(title),
		UpdateTime:  easygo.NewInt64(time.Now().Unix()),
		CreateTime:  easygo.NewInt64(time.Now().Unix()),
		Status:      easygo.NewInt32(1),
		ImageUrl:    easygo.NewString(ImageUrl),
		ActionCount: easygo.NewInt32(0),
		MenuId:      easygo.NewInt32(menuid),
		/*JumpType:    easygo.NewInt32(1),
		JumpObjId:   easygo.NewInt64(10001),
		JumpObject:  easygo.NewInt32(5),
		JumpUrl:     easygo.NewString("http://www.baidu.com"),*/
		Weight: easygo.NewInt32(5),
	}

	return CreateTableESPortsCarousel(info)

}

func CreateTableESPortsCarousel(info *share_message.TableESPortsCarousel) (int32, string) {

	id := for_game.NextId(for_game.TABLE_ESPORTS_CAROUSEL)
	timeUnix := time.Now().Unix()
	info.UpdateTime = easygo.NewInt64(timeUnix)
	info.CreateTime = easygo.NewInt64(timeUnix)
	col, closeFun := GetC(for_game.TABLE_ESPORTS_CAROUSEL)
	defer closeFun()

	info.Id = easygo.NewInt64(id)
	_, err := col.Upsert(bson.M{"_id": id}, bson.M{"$set": info})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	return for_game.C_OPT_SUCCESS, "創建成功"
}
func DeleteTableESPortsCarousel(id int64) bool {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_CAROUSEL)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"_id": id})
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}

func UpdateTableESPortsCarousel(info *share_message.TableESPortsCarousel) (int32, string) {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_CAROUSEL)
	defer closeFun()
	updatedata := bson.M{}
	if info.Title != nil {
		updatedata["Title"] = info.GetTitle()
	}
	updatedata["UpdateTime"] = time.Now().Unix()
	if info.Status != nil {
		updatedata["Status"] = info.GetStatus()
	}
	if info.ImageUrl != nil {
		updatedata["ImageUrl"] = info.GetImageUrl()
	}

	if info.ActionCount != nil {
		updatedata["ActionCount"] = info.GetActionCount()
	}
	if info.MenuId != nil {
		updatedata["ContentType"] = info.GetMenuId()
	}

	cinfo, err := col.Upsert(bson.M{"_id": info.GetId()}, bson.M{"$set": updatedata})
	if err != nil {
		logs.Error(err)
		return for_game.C_SYS_ERROR, "系统错误"
	}
	if cinfo.Updated > 0 {
		return for_game.C_OPT_SUCCESS, "修改成功"
	} else {
		return for_game.C_INFO_NOT_EXISTS, "数据不存在，修改失败"
	}
}
func GetTableESPortsCarouselById(id int64) *share_message.TableESPortsCarousel {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_CAROUSEL)
	defer closeFun()
	data := &share_message.TableESPortsCarousel{}
	err := col.Find(bson.M{"_id": id}).One(&data)
	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return nil
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return data
}
func GetTableESPortsCarouselByMenuId(menuid int32, status int32) []*share_message.TableESPortsCarousel {
	col, closeFun := GetC(for_game.TABLE_ESPORTS_CAROUSEL)
	defer closeFun()

	var data []*share_message.TableESPortsCarousel
	err := col.Find(bson.M{"MenuId": menuid, "Status": status}).All(&data)
	if err != nil && err != mgo.ErrNotFound {
		logs.Error(err)
		return nil
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return data
}
func GetTableESPortsCarouselList1(offset, limit int, sort string, keyword string, begin, end int64, Status int32) ([]*share_message.TableESPortsCarousel, int) {
	var list []*share_message.TableESPortsCarousel
	col, closeFun := GetC(for_game.TABLE_ESPORTS_CAROUSEL)
	defer closeFun()
	queryBson := bson.M{}
	if keyword != "" {
		queryBson["Title"] = bson.M{"$regex": "^" + keyword + "+"}
	}
	if Status != 0 {
		queryBson["Status"] = Status
	}
	if begin != 0 && end != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": begin, "$lte": end}
	}
	query := col.Find(queryBson)
	count, err := query.Count()
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	err = query.Sort(sort).Skip(offset).Limit(limit).All(&list)
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	if list == nil {
		list = []*share_message.TableESPortsCarousel{}
	}
	return list, count
}
func GetTableESPortsCarouselList2(cPage, pSize int, sort string, keyword string, begin, end int64, Status int32) ([]*share_message.TableESPortsCarousel, int) {
	pageSize := int(pSize)
	curPage := easygo.If(int(cPage) > 1, int(cPage)-1, 0).(int)
	var list []*share_message.TableESPortsCarousel
	col, closeFun := GetC(for_game.TABLE_ESPORTS_CAROUSEL)
	defer closeFun()
	queryBson := bson.M{}
	if keyword != "" {
		queryBson["Title"] = bson.M{"$regex": "^" + keyword + "+"}
	}
	if Status != 0 {
		queryBson["Status"] = Status
	}
	if begin != 0 && end != 0 {
		queryBson["CreateTime"] = bson.M{"$gte": begin, "$lte": end}
	}
	query := col.Find(queryBson)
	count, err := query.Count()
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	err = query.Sort(sort).Skip(curPage * pageSize).Limit(pageSize).All(&list)
	if err != nil {
		logs.Error(err)
		return nil, 0
	}
	return list, count
}
