// 管理后台为[浏览器]提供的服务

package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
)

//查询管理员列表
func (self *cls4) RpcManagerList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.GetPlayerListRequest) easygo.IMessage {
	list, count := GetManagerList(user, reqMsg)

	msg := &brower_backstage.GetManagerListResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return msg
}

//修改管理员
func (self *cls4) RpcEditManager(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.Manager) easygo.IMessage {
	admin := GetUser(reqMsg.GetId())
	if reqMsg.GetAccount() != admin.GetAccount() {
		return easygo.NewFailMsg("管理员账号不能修改")
	}

	//修改之前先强制下线
	sid := PServerInfo.GetSid()
	euser := for_game.GetRedisAdmin(reqMsg.GetId())
	if euser != nil {
		if euser.ServerId == sid {
			ReplaceLogin(reqMsg.GetId())
		} else {
			msg := &server_server.AdminInfo{
				UserId:   easygo.NewInt64(reqMsg.GetId()),
				ServerId: easygo.NewInt32(sid),
			}
			ChooseOneHall(0, "RpcReplaceLoginToHall", msg)
		}
	}

	if reqMsg.Password == nil || reqMsg.GetPassword() == "" {
		reqMsg.Password = admin.Password
		EditManage(user.GetSite(), reqMsg, "")
	} else {
		EditManage(user.GetSite(), reqMsg, "edit")
	}
	// ep.SetUser(reqMsg)
	msg := fmt.Sprintf("修改管理员:%s资料", reqMsg.GetAccount())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.USER_MANAGE, msg)

	return easygo.EmptyMsg
}

//创建管理员
func (self *cls4) RpcAddManager(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.Manager) easygo.IMessage {
	if reqMsg.Account == nil {
		return easygo.NewFailMsg("管理员帐户不能为空")
	}
	if reqMsg.Password == nil {
		return easygo.NewFailMsg("管理员密码不能为空")
	}

	admin := QueryManage(reqMsg.GetAccount())
	if admin != nil && admin.GetSite() == user.GetSite() {
		return easygo.NewFailMsg("账号重复！请修改账号名重新确认添加")
	}
	if admin != nil && admin.GetSite() != user.GetSite() {
		return easygo.NewFailMsg("非法账号！请修改账号名重新确认添加")
	}

	s := reqMsg.GetSite()
	if reqMsg.GetSite() == "" {
		s = user.GetSite()
	}

	if user.GetRole() > 0 && reqMsg.GetRoleType() == 1 {
		return easygo.NewFailMsg("权限不足无法创建管理员帐号")
	}

	AddManage(s, reqMsg)

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.USER_MANAGE, "创建管理员:"+reqMsg.GetAccount())
	return easygo.EmptyMsg
}

//冻结管理员
func (self *cls4) RpcAdminFreeze(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	UpAdminStatus(reqMsg.GetIds32(), 1)

	for _, u := range reqMsg.GetIds32() {
		mp := BrowerEpMp.LoadEndpoint(int64(u))
		if mp != nil {
			msg := &brower_backstage.ErrMessage{Err: easygo.NewString("帐号被冻结")}
			mp.RpcReplacePush(msg) //通知前端下线
			mp.Shutdown()          //强制下线
		}
	}

	var ids string
	idsarr := reqMsg.GetIds32()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}

	msg := fmt.Sprintf("批量冻结管理员: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.USER_MANAGE, msg)
	return easygo.EmptyMsg
}

//解冻管理员
func (self *cls4) RpcAdminUnFreeze(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	UpAdminStatus(reqMsg.GetIds32(), 0)

	var ids string
	idsarr := reqMsg.GetIds32()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}

	msg := fmt.Sprintf("批量解冻管理员: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.USER_MANAGE, msg)
	return easygo.EmptyMsg
}

//查询客服分类列表
func (self *cls4) RpcQueryManagerTypes(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, count := GetManagerTypesList(reqMsg)

	msg := &brower_backstage.ManagerTypesResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}

	return msg
}

//客服分类下拉
func (self *cls4) RpcManagerTypesKeyValue(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	list := GetManagerTypesListNopage()
	return &brower_backstage.KeyValueResponseTag{
		List: list,
	}
}

//修改客服分类
func (self *cls4) RpcEditManagerTypes(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.ManagerTypes) easygo.IMessage {
	msg := fmt.Sprintf("修改客服类型:%d", reqMsg.GetId())
	if reqMsg.GetName() == "" || reqMsg.Name == nil {
		return easygo.NewFailMsg("分类名称不能为空")
	}
	if reqMsg.GetStatus() == 0 || reqMsg.Status == nil {
		return easygo.NewFailMsg("状态不能为空")
	}
	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_MANAGER_TYPES))
		msg = fmt.Sprintf("添加客服类型:%d", reqMsg.GetId())
	}
	EditManageTypes(reqMsg)

	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.USER_MANAGE, msg)

	return easygo.EmptyMsg
}

// 查询角色管理列表
func (self *cls4) RpcQueryRolePower(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	// if user.GetRole() != 0 {
	// 	return easygo.NewFailMsg("权限不足")
	// }

	list, count := QueryRolePowerList(reqMsg)
	msg := &brower_backstage.QueryRolePowerList{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	return msg
}

// 更新角色权限
func (self *cls4) RpcUpdateRolePower(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *share_message.RolePower) easygo.IMessage {
	if user.GetRole() != 0 {
		return easygo.NewFailMsg("权限不足")
	}

	result := CheckRoleName(user.GetSite(), reqMsg)
	if result != 0 {
		return easygo.NewFailMsg("角色名称重复")
	}
	if reqMsg.Id == nil || reqMsg.GetId() == 0 {
		reqMsg.Id = easygo.NewInt32(for_game.NextId(for_game.TABLE_ROLEPOWER))
	}
	if len(reqMsg.GetMenuIds()) == 0 {
		return easygo.NewFailMsg("权限不能为空")
	}

	UpdateRolePwer(user.GetSite(), reqMsg)
	msg := fmt.Sprintf("更新角色Id:%d的权限", reqMsg.GetId())
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ROLE_MANAGE, msg)

	return easygo.EmptyMsg
}

// 删除角色权限
func (self *cls4) RpcDeleteRolePower(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataByIds) easygo.IMessage {
	if user.GetRole() != 0 {
		return easygo.NewFailMsg("权限不足")
	}

	for _, item := range reqMsg.GetIds64() {
		if item == 1 {
			return easygo.NewFailMsg("不能删除管理员角色")
		}
	}

	idList := reqMsg.GetIds64()
	// 判断权限角色是否在被使用
	errStr := CheckAuthGroupByRole(idList)
	if errStr != "" {
		return easygo.NewFailMsg(errStr)
	}

	err := DelDataById(for_game.TABLE_ROLEPOWER, idList)

	easygo.PanicError(err)

	var ids string
	idsarr := reqMsg.GetIds64()
	count := len(idsarr)
	for i := 0; i < count; i++ {
		if i < count {
			ids += easygo.IntToString(int(idsarr[i])) + ","

		} else {
			ids += easygo.IntToString(int(idsarr[i]))
		}
	}

	msg := fmt.Sprintf("批量删除角色权限: %s", ids)
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.ROLE_MANAGE, msg)
	return easygo.EmptyMsg
}

//查询管理员列表下拉配置列表
func (self *cls4) RpcGetRolePowerList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	list := GetRolePowerList()
	return &brower_backstage.KeyValueResponseTag{List: list}
}

//根据角色id获取角色权限
func (self *cls4) RpcGetPowerRouter(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.QueryDataById) easygo.IMessage {
	ids := reqMsg.GetId32()
	if ids == 0 {
		return easygo.NewFailMsg("角色id不能为空")
	}
	result := GetPowerRouter(ids)
	if result == nil {
		return easygo.NewFailMsg("找不到对应的角色信息！ ")
	}

	return result
}

// 查询客户端版本管理
func (self *cls4) RpcQueryVersion(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	jsonFile := easygo.YamlCfg.GetValueAsString("CLIENT_VERSION_DATA")
	data := for_game.FindVersion(jsonFile)
	return data
}

// 更新客户端版本管理
func (self *cls4) RpcUpdateVersion(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.VersionData) easygo.IMessage {
	jsonFile := easygo.YamlCfg.GetValueAsString("CLIENT_VERSION_DATA")
	for_game.UpdateAll(reqMsg, jsonFile)

	return easygo.EmptyMsg
}

// 获取转发服务器列表
func (self *cls4) RpcQueryTfserver(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	jsonFile := easygo.YamlCfg.GetValueAsString("CLIENT_TFSERVER_ADDR")
	data := for_game.FindTfserver(jsonFile)
	ps := []*brower_backstage.Hostlist{}

	for _, i := range data.Hostlist {
		ts := &brower_backstage.Hostlist{
			Ip: easygo.NewString(i),
		}
		ps = append(ps, ts)
	}
	msg := &brower_backstage.Tfserver{
		Hostlist: ps,
	}

	return msg
}

// 更新转发服务器列表
func (self *cls4) RpcUpdateTfserver(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.Tfserver) easygo.IMessage {
	jsonFile := easygo.YamlCfg.GetValueAsString("CLIENT_TFSERVER_ADDR")

	data := &for_game.TFserver{}
	hostlist := data.Hostlist
	for _, i := range reqMsg.Hostlist {
		hostlist = append(hostlist, i.GetIp())
	}
	data.Hostlist = hostlist
	for_game.UpdateTfserver(data, jsonFile)

	return easygo.EmptyMsg
}
