// 大厅服务器为[游戏客户端]提供的服务

package login

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/pb/client_server"

	"github.com/astaxie/beego/logs"
)

/*func (self *WebHttpForClient) CheckLoginState(player *for_game.RedisPlayerBaseObj, msg *client_login.LoginResult) bool {
	if player == nil {
		code := client_login.LoginMark_LOGIN_ERROR_ACCOUNT
		msg.Result = &code
		return false
	}
	if player.GetStatus() == for_game.ACCOUNT_USER_FROZEN || player.GetStatus() == for_game.ACCOUNT_ADMIN_FROZEN {
		code := client_login.LoginMark_LOGIN_ERROR_FREEZEACCOUT
		msg.Result = &code
		return false
	}
	if player.GetStatus() == for_game.ACCOUNT_CANCELING {
		code := client_login.LoginMark_LOGIN_ERROR_ACCOUNT_CANCELING
		msg.Result = &code
		return false
	}
	if player.GetStatus() == for_game.ACCOUNT_CANCELED {
		code := client_login.LoginMark_LOGIN_ERROR_ACCOUNT_CANCELED
		phone := player.GetPhone()
		if phone == "" {
			phone = for_game.GetPhoneByPlayerId(player.GetPlayerId())
		}
		overTime := for_game.CheckCancelAccount(phone)
		if overTime > 0 {
			//检测账号是否在注销期未满60天
			s := "此账号" + easygo.Stamp2StrExt(overTime) + "成功注销，自注销之日起6个月内不得重新注册"
			msg.Result = &code
			msg.ErrMsg = easygo.NewString(s)
			return false
		}

	}
	return true
}

// 登录都用这一个  注册人数,登录次数.
func (self *WebHttpForClient) RpcLoginHall(common *base.Common, reqMsg *client_login.LoginMsg) easygo.IMessage {
	// logs.Info("登陆服登陆开始,id为-------->", reqMsg.GetPlayerId()) // 别删，永久留存
	logs.Info("RpcLoginHall,msg=%v", reqMsg) // 别删，永久留存
	t := reqMsg.GetType()
	account := reqMsg.GetAccount()
	msg := &client_login.LoginResult{
		Account:   easygo.NewString(account),
		LoginType: easygo.NewInt32(t),
	}
	//服务器关闭拒绝
	if IsStopServer {
		code := client_login.LoginMark_LOGIN_ERROR_STOPSERVER
		msg.Result = &code
		return msg
	}
	add := ep.GetClientAddr()
	token := easygo.RandStringRunes(16) + account
	if add == "" {
		add = ep.GetConnection().RemoteAddr().String()
	}
	loginInfo := &share_message.LogLoginInfo{
		Account:   easygo.NewString(account),
		PF:        easygo.NewString(reqMsg.GetMark()),
		LoginTime: easygo.NewInt64(time.Now().Unix()),
		LoginIP:   easygo.NewString(add),
	}

	defer LogLoginInfo(loginInfo)

	if t == OneKey_Login { //一键登录
		phone := reqMsg.GetOpenId()
		if phone == "" { //是否传了手机号
			onkeyToken := reqMsg.GetOneKeyToken()
			phone = for_game.GetJGOneKeyLoginPhone(onkeyToken, reqMsg.GetApkCode())
			if phone == "" {
				code := client_login.LoginMark_LOGIN_ERROR_ONEKEYPHONE
				msg.Result = &code
				return msg
			}
		}
		Info := for_game.GetRedisAccountByPhone(phone)
		var pid PLAYER_ID
		if Info == nil { //如果没有这个用户就创建用户
			overTime := for_game.CheckCancelAccount(phone)
			logs.Info("overTime:", overTime)
			if overTime > 0 {
				//检测账号是否在注销期未满60天
				s := "此账号" + easygo.Stamp2StrExt(overTime) + "成功注销，自注销之日起6个月内不得重新注册"
				msg.ErrMsg = easygo.NewString(s)
				logs.Info("账号未满60天")
				code := client_login.LoginMark_LOGIN_ERROR_ACCOUNT_CANCELED
				msg.Result = &code
				return msg
			}
			data := &share_message.CreateAccountData{
				Phone:    easygo.NewString(phone),
				PassWord: easygo.NewString(""),
				Ip:       easygo.NewString(add),
				IsOnline: easygo.NewBool(true),
			}
			var b bool
			b, pid = for_game.CreateAccount(data)
			if !b {
				panic("一键登录创建玩家失败")
			}
			Info = for_game.GetRedisAccountByPhone(phone)
			for_game.AddStatisticsInfo(for_game.LOGINREGISTER_ONEKEYREGISTER, pid,1)
			msg.IsUserReg = easygo.NewBool(true) // todo 注册操作 注册人数埋点
		} else {
			pid = Info.GetPlayerId()
			//  登录操作
			msg.IsLoginFreq = easygo.NewBool(true)
		}
		if pid == 0 {
			panic("人物id怎么会是0")
		}
		oldEp := ClientEpMp.LoadEndpoint(pid)
		if oldEp != nil && oldEp != ep {
			logs.Error("正在登录中1", pid)
			return easygo.NewFailMsg("正在登陆中", for_game.FAIL_MSG_CODE_1001)

		}
		player := for_game.GetRedisPlayerBase(pid)
		//检测登录状态
		if !self.CheckLoginState(player, msg) {
			return msg
		}
		ClientEpMp.StoreEndpoint(pid, ep.GetEndpointId())
		player.SetAutoLoginInfo(token)
		msg = GetLoginMsg(pid, phone, token, t)
		loginInfo.PlayerId = easygo.NewInt64(pid)
		loginInfo.State = easygo.NewInt32(0)
		loginInfo.RegTime = easygo.NewInt64(Info.GetCreateTime())
		//  添加埋点信息
		loginEventResp := for_game.GetLoginEvent(&client_login.LoginEventRequst{
			PlayerId:   easygo.NewInt64(pid),
			DeviceCode: easygo.NewString(reqMsg.GetMark()),
		})
		// logs.Info("GetLoginEvent----------->%+v", loginEventResp)
		msg.IsAppAct = loginEventResp.IsAppAct
		msg.IsLoginMan = loginEventResp.IsLoginMan
		return msg
	} else if t == Visitor_Login { //游客登录
		data := &share_message.CreateAccountData{
			Phone:     easygo.NewString(account),
			PassWord:  easygo.NewString(""),
			IsVisitor: easygo.NewBool(true),
			Ip:        easygo.NewString(add),
			IsOnline:  easygo.NewBool(true),
		}
		b, pid := for_game.CreateAccount(data)
		if !b {
			panic("游客登录创建玩家失败")
		}
		oldEp := ClientEpMp.LoadEndpoint(pid)
		if oldEp != nil && oldEp != ep {
			logs.Error("正在登录中2", pid)
			return easygo.NewFailMsg("正在登陆中", for_game.FAIL_MSG_CODE_1001)
		}
		ClientEpMp.StoreEndpoint(pid, ep.GetEndpointId())
		Info := for_game.GetRedisAccountByPhone(account)
		msg = GetLoginMsg(pid, account, token, t)
		loginInfo.PlayerId = easygo.NewInt64(pid)
		loginInfo.State = easygo.NewInt32(0)
		loginInfo.RegTime = easygo.NewInt64(Info.GetCreateTime())
		//  添加埋点信息,注册人数,登录次数自己填充.
		msg.IsUserReg = easygo.NewBool(true)
		loginEventResp := for_game.GetLoginEvent(&client_login.LoginEventRequst{
			PlayerId:   easygo.NewInt64(pid),
			DeviceCode: easygo.NewString(reqMsg.GetMark()),
		})
		msg.IsAppAct = loginEventResp.IsAppAct
		msg.IsLoginMan = loginEventResp.IsLoginMan
		// logs.Info("GetLoginEvent----------->%+v", loginEventResp)
		return msg
	} else if t == Wechat_Login {
		var info *share_message.PlayerAccount
		var wechatToken, openId, unionId string
		var pid PLAYER_ID
		var phone string
		var createTime int64
		if reqMsg.GetUnionId() != "" {
			//微信openid登录，找不到账号直接返回
			info = for_game.GetPlayerInfoForUnionId(reqMsg.GetUnionId())
			if info == nil {
				if reqMsg.GetPhone() == "" {
					code := client_login.LoginMark_LOGIN_ERROR_WECHATTOKEN
					msg = &client_login.LoginResult{
						Result: &code,
					}
					return msg
				} else {
					//微信绑定手机号登录
					if for_game.IS_FORMAL_SERVER {
						if err := for_game.CheckMessageCode(reqMsg.GetPhone(), reqMsg.GetLoginCode(), for_game.CLIENT_CODE_BINDPHONE); err != nil {
							//验证码错误
							code := client_login.LoginMark_LOGIN_ERROR_MESSAGE
							msg = &client_login.LoginResult{
								Result: &code,
							}
							return msg
						}
					}
					acc := for_game.GetRedisAccountByPhone(reqMsg.GetPhone())
					if acc != nil {
						if acc.GetUnionId() != "" {
							logs.Info("手机号已经被绑定其他账号")
							code := client_login.LoginMark_LOGIN_ERROR_BIND_PHONE_REPEAT
							msg = &client_login.LoginResult{
								Result: &code,
							}
							return msg
						}
						//把openid绑定到指定存在的账号上
						acc.SetOpenId(reqMsg.GetOpenId())
						acc.SetUnionId(reqMsg.GetUnionId())
						acc.SetIsBind(true)
						acc.SaveToMongo()
						info = acc.GetRedisAccount()
					} else {
						logs.Info("手机号未绑定过账号")
						//创建新账号并绑定手机号
						//如果没有绑定手机号，则先绑定手机号
						unionId = reqMsg.GetUnionId()
						_, sex, headIcon := for_game.GetWeChatUserInfo(reqMsg.GetWeChatToken(), reqMsg.GetOpenId())
						name := for_game.GetRandNickName()
						if headIcon == "" { //如果没有头像
							icon := rand.Intn(5) + 1
							if sex == 0 || sex == 2 {
								sex = 2
								headIcon = fmt.Sprintf("https://im-resource-1253887233.cos.accelerate.myqcloud.com/defaulticon/girl_%d.png", icon)
							} else {
								headIcon = fmt.Sprintf("https://im-resource-1253887233.cos.accelerate.myqcloud.com/defaulticon/boy_%d.png", icon)
							}
						}
						pid = for_game.CreateAccountForWechat(reqMsg.GetPhone(), name, headIcon, unionId, sex, reqMsg.GetAreaCode())
						phone = reqMsg.GetPhone()
						m := for_game.GetRedisAccountObj(pid)
						createTime = m.GetCreateTime()
						for_game.AddStatisticsInfo(for_game.LOGINREGISTER_WECHATREGISTER, pid,1)
						//  注册操作
						msg.IsUserReg = easygo.NewBool(true)
					}
				}

			}
		} else {
			wechatToken, openId, unionId = for_game.GetWeChatInfo(reqMsg.GetWechatCode(), reqMsg.GetApkCode())
			logs.Info("拿到微信openid:", wechatToken, openId, unionId)
			if wechatToken == "" || openId == "" || unionId == "" {
				code := client_login.LoginMark_LOGIN_ERROR_WECHATTOKEN
				msg = &client_login.LoginResult{
					Result: &code,
				}
				return msg
			}
			info = for_game.GetPlayerInfoForUnionId(unionId)
			//TODO:新的微信登录，要绑定手机号才能登录,后续去掉这个判定
			if info == nil {
				//新微信登录
				code := client_login.LoginMark_LOGIN_ERROR_BIND_PHONE
				msg = &client_login.LoginResult{
					Result:      &code,
					OpenId:      easygo.NewString(openId),
					WeChatToken: easygo.NewString(wechatToken),
					UnionId:     easygo.NewString(unionId),
				}
				return msg
			}
		}
		if info != nil {
			pid = info.GetPlayerId()
			phone = info.GetAccount()
			createTime = info.GetCreateTime()
			//  登录操作
			msg.IsLoginFreq = easygo.NewBool(true)
		}

		if pid == 0 {
			panic("玩家id怎么是0")
		}
		oldEp := ClientEpMp.LoadEndpoint(pid)
		if oldEp != nil && oldEp != ep {
			logs.Error("正在登录中3", pid)
			return easygo.NewFailMsg("正在登陆中", for_game.FAIL_MSG_CODE_1001)
		}
		player := for_game.GetRedisPlayerBase(pid)
		if !self.CheckLoginState(player, msg) {
			return msg
		}
		//if player.GetStatus() == 1 || player.GetStatus() == 2 {
		//	code := client_login.LoginMark_LOGIN_ERROR_FREEZEACCOUT
		//	msg = &client_login.LoginResult{
		//		Result: &code,
		//	}
		//	return msg
		//}
		ClientEpMp.StoreEndpoint(pid, ep.GetEndpointId())
		player.SetAutoLoginInfo(token)
		msg = GetLoginMsg(pid, phone, token, t)
		loginInfo.PlayerId = easygo.NewInt64(pid)
		loginInfo.State = easygo.NewInt32(0)
		loginInfo.RegTime = easygo.NewInt64(createTime)
		//  添加埋点信息
		loginEventResp := for_game.GetLoginEvent(&client_login.LoginEventRequst{
			PlayerId:   easygo.NewInt64(pid),
			DeviceCode: easygo.NewString(reqMsg.GetMark()),
		})
		msg.IsAppAct = loginEventResp.IsAppAct
		msg.IsLoginMan = loginEventResp.IsLoginMan
		// logs.Info("GetLoginEvent----------->%+v", loginEventResp)
		return msg
	} else if t == Auto_Login { // 自动登录
		autoToken := reqMsg.GetToken()
		pid := reqMsg.GetPlayerId()
		oldEp := ClientEpMp.LoadEndpoint(pid)
		if oldEp != nil && oldEp != ep {
			logs.Error("正在登录中4", pid)
			return easygo.NewFailMsg("正在登陆中", for_game.FAIL_MSG_CODE_1001)
		}
		info := for_game.GetRedisAccountObj(pid)
		if info == nil {
			code := client_login.LoginMark_LOGIN_ERROR_ACCOUNT
			msg = &client_login.LoginResult{
				Result: &code,
			}
			return msg
		}
		player := for_game.GetRedisPlayerBase(pid)
		if !self.CheckLoginState(player, msg) {
			return msg
		}
		//if player.GetStatus() == 1 || player.GetStatus() == 2 {
		//	code := client_login.LoginMark_LOGIN_ERROR_FREEZEACCOUT
		//	msg = &client_login.LoginResult{
		//		Result: &code,
		//	}
		//	return msg
		//}
		ClientEpMp.StoreEndpoint(pid, ep.GetEndpointId())
		stoken := player.GetAutoLoginToken()
		ti := player.GetAutoLoginTime()
		if stoken != autoToken || time.Now().Unix()-ti > 7*86400 {
			code := client_login.LoginMark_LOGIN_ERROR_AUTOTOKEN
			msg = &client_login.LoginResult{
				Result: &code,
			}
			return msg
		}
		player.SetAutoLoginInfo(autoToken)
		msg = GetLoginMsg(pid, account, autoToken, t)
		loginInfo.PlayerId = easygo.NewInt64(pid)
		loginInfo.State = easygo.NewInt32(0)
		loginInfo.RegTime = easygo.NewInt64(info.GetCreateTime())
		//  添加埋点信息
		msg.IsLoginFreq = easygo.NewBool(true)
		loginEventResp := for_game.GetLoginEvent(&client_login.LoginEventRequst{
			PlayerId:   easygo.NewInt64(pid),
			DeviceCode: easygo.NewString(reqMsg.GetMark()),
		})
		msg.IsAppAct = loginEventResp.IsAppAct
		msg.IsLoginMan = loginEventResp.IsLoginMan
		// logs.Info("GetLoginEvent----------->%+v", loginEventResp)
		return msg
	}

	loginAuth := for_game.LoginIpAuth(reqMsg.GetLoginIp())
	if !loginAuth {
		code := client_login.LoginMark_LOGIN_ERROR_FREEZEIP
		msg.Result = &code
		return msg
	}
	var playerId int64
	var player *for_game.RedisPlayerBaseObj
	password := reqMsg.GetPassword()
	loginCode := reqMsg.GetLoginCode()
	info := for_game.GetRedisAccountByPhone(account)
	if info == nil { //如果没有这个账号
		if for_game.IS_FORMAL_SERVER && !MessageAuth(for_game.CLIENT_CODE_LOGIN, account, loginCode) {
			code := client_login.LoginMark_LOGIN_ERROR_MESSAGE
			msg.Result = &code
			return msg
		}
		data := &share_message.CreateAccountData{
			Phone:     easygo.NewString(account),
			PassWord:  easygo.NewString(password),
			IsVisitor: easygo.NewBool(false),
			Ip:        easygo.NewString(reqMsg.GetLoginIp()),
			IsOnline:  easygo.NewBool(true),
			AreaCode:  easygo.NewString(reqMsg.GetAreaCode()),
		}
		b, pid := for_game.CreateAccount(data)
		if !b {
			//账号注册失败
			code := client_login.LoginMark_REGISTER_ERROR_CREATEACCOUNT
			msg.Result = &code
			return msg
		}
		info = for_game.GetRedisAccountObj(pid)
		for_game.AddStatisticsInfo(for_game.LOGINREGISTER_PHONEREGISTER, pid,1)
		playerId = pid
		player = for_game.GetRedisPlayerBase(playerId)
		//   注册操作
		msg.IsUserReg = easygo.NewBool(true)
	} else {
		playerId = info.GetPlayerId()
		player = for_game.GetRedisPlayerBase(playerId)

		if t == Password_Login { //如果是密码登录
			if info.GetPassword() != for_game.Md5(password) {
				code := client_login.LoginMark_LOGIN_ERROR_PASSWORD
				msg.Result = &code
				return msg
			}
		} else {
			if player.GetTypes() != 1 { //不是普通用户  用密码登陆
				if info.GetPassword() != for_game.Md5(loginCode) {
					code := client_login.LoginMark_LOGIN_ERROR_PASSWORD
					msg.Result = &code
					return msg
				}
			} else { //普通用户 用验证码登陆
				if for_game.IS_FORMAL_SERVER && !MessageAuth(for_game.CLIENT_CODE_LOGIN, account, loginCode) {
					code := client_login.LoginMark_LOGIN_ERROR_MESSAGE
					msg.Result = &code
					return msg
				}
			}
		}
		oldEp := ClientEpMp.LoadEndpoint(playerId)
		if oldEp != nil && oldEp != ep {
			logs.Error("正在登录中4", playerId)
			return easygo.NewFailMsg("正在登陆中", for_game.FAIL_MSG_CODE_1001)
		}
		ClientEpMp.StoreEndpoint(playerId, ep.GetEndpointId())
		//   登录操作
		msg.IsLoginFreq = easygo.NewBool(true)
	}
	if !self.CheckLoginState(player, msg) {
		return msg
	}
	//if player.GetStatus() == 1 || player.GetStatus() == 2 {
	//	code := client_login.LoginMark_LOGIN_ERROR_FREEZEACCOUT
	//	msg.Result = &code
	//	return msg
	//}

	//accountInfo := for_game.GetPlayerAccount(playerId)
	//accountInfo.SetToken(token)

	player.SetToken(token)
	player.SetDeviceCode(reqMsg.GetMark()) //记录用户最后登录的设备码
	player.SetAutoLoginInfo(token)
	code := client_login.LoginMark_LOGIN_SUCCESS
	msg.Result = &code
	serverInfo := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_HALL)
	var address string
	var port int32
	//根据连接类型获取对应的端口
	if ep.GetConType() == for_game.CONN_TYPE_WEBSOCKET {
		port = serverInfo.GetClientWSPort()
	} else if ep.GetConType() == for_game.CONN_TYPE_TCP {
		port = serverInfo.GetClientTCPPort()
	}
	if for_game.IS_TFSERVER {
		address = easygo.AnytoA(port)
	} else {
		address = for_game.MakeAddress(serverInfo.GetExternalIp(), port)
	}
	logs.Info("返回客户端地址：", address)
	msg.Address = easygo.NewString(address)
	msg.Token = easygo.NewString(token)
	msg.LoginType = easygo.NewInt32(t)
	msg.PlayerId = easygo.NewInt64(playerId)
	loginInfo.PlayerId = easygo.NewInt64(playerId)
	loginInfo.State = easygo.NewInt32(0)
	loginInfo.RegTime = easygo.NewInt64(info.GetCreateTime())
	//  添加埋点信息
	msg.IsLoginFreq = easygo.NewBool(true)
	loginEventResp := for_game.GetLoginEvent(&client_login.LoginEventRequst{
		PlayerId:   easygo.NewInt64(playerId),
		DeviceCode: easygo.NewString(reqMsg.GetMark()),
	})
	// logs.Info("GetLoginEvent----------->%+v", loginEventResp)
	msg.IsAppAct = loginEventResp.IsAppAct
	msg.IsLoginMan = loginEventResp.IsLoginMan
	// logs.Info("=======================登陆服登陆结束,id为-------->", reqMsg.GetPlayerId()) // 别删，永久留存
	if msg.GetIsUserReg() {
		easygo.Spawn(for_game.UpdatePosDeviceAdvIdfa, reqMsg.GetMark(), "IsRegister", true) //广告设备注册
	}

	return msg
}

*/
//客户端请求验证码
func (self *WebHttpForClient) RpcClientGetCode(common *base.Common, reqMsg *client_server.GetCodeRequest) easygo.IMessage {
	logs.Info("====api请求短信验证码=====:", reqMsg)
	/*phone := reqMsg.GetPhone()
	t := reqMsg.GetType()
	if !for_game.IS_FORMAL_SERVER {
		return nil
	}

	if !for_game.MessageMarkInfo.CheckPhoneVaild(phone) {
		return easygo.NewFailMsg("你操作频繁过快，请稍后再试")
	}

	data := for_game.MessageMarkInfo.GetMessageMarkInfo(t, phone)
	if data != nil {
		leaveTime := time.Now().Unix() - data.Timestamp
		if leaveTime <= 55 {
			return easygo.NewFailMsg("验证码已发送!")
		}
	}
	codes := for_game.SendCodeToClientUser(t, phone, reqMsg.GetAreaCode())
	if codes != "" {
		return nil
	}*/
	return easygo.NewFailMsg("验证码发送失败！")
}

/*
//转发客户端连接报道
func (self *WebHttpForClient) RpcTFToServer(common *base.Common, reqMsg *client_server.ClientInfo) easygo.IMessage {
	ep.AddAddrs(reqMsg.GetIp())
	//	logs.Info("玩家真是地址:", ep.GetClientAddr())
	return nil
}

//检查验证码
func (self *WebHttpForClient) RpcCheckMessageCode(common *base.Common, reqMsg *client_server.CodeResponse) easygo.IMessage {
	if !for_game.IS_FORMAL_SERVER {
		return nil
	}
	code := reqMsg.GetCode()
	t := reqMsg.GetType()
	phone := reqMsg.GetPhone()
	data := for_game.MessageMarkInfo.GetMessageMarkInfo(t, phone)
	if data == nil {
		res := "验证码不存在"
		return easygo.NewFailMsg(res)
	}
	if data.Mark != code {
		res := "验证码不正确"
		return easygo.NewFailMsg(res)
	}
	return nil
}

//修改登录密码
func (self *WebHttpForClient) RpcForgetLoginPassword(common *base.Common, reqMsg *client_login.LoginMsg) easygo.IMessage {
	password := reqMsg.GetPassword()
	account := reqMsg.GetAccount()
	accountInfo := for_game.GetRedisAccountByPhone(account)
	if accountInfo == nil {
		res := "不存在该账号"
		return easygo.NewFailMsg(res)
	}
	if accountInfo.GetPassword() == for_game.Md5(password) {
		res := "新密码与旧密码相同"
		return easygo.NewFailMsg(res)
	}
	accountInfo.SetPassword(password)
	return nil
}

func (self *WebHttpForClient) RpcCheckAccountVaild(common *base.Common, reqMsg *client_server.CheckInfo) easygo.IMessage {
	logs.Info("检测账号有效性:", reqMsg)
	account := reqMsg.GetAccount()
	accountInfo := for_game.GetRedisAccountByPhone(account)
	logs.Info("accountInfo:", accountInfo)
	//info := accountInfo.GetPlayerAccount()
	if accountInfo == nil {
		overTime := for_game.CheckCancelAccount(account)
		logs.Info("overTime:", overTime)
		if overTime > 0 {
			//检测账号是否在注销期未满60天
			s := "此账号" + easygo.Stamp2StrExt(overTime) + "成功注销，自注销之日起6个月内不得重新注册"
			reqMsg.ErrMsg = easygo.NewString(s)
			reqMsg.Vaild = easygo.NewBool(true)
			logs.Info("账号未满60天")
			reqMsg.State = easygo.NewInt32(for_game.ACCOUNT_CANCELED)
		} else {
			reqMsg.Vaild = easygo.NewBool(false)
			reqMsg.State = easygo.NewInt32(for_game.ACCOUNT_NORMAL)
		}
	} else {
		//检测是否可以绑定这个手机号
		if reqMsg.GetIsCheckPhone() {
			logs.Info("accountInfo:", accountInfo.GetIsBind())
			if accountInfo.GetUnionId() != "" {
				reqMsg.Vaild = easygo.NewBool(true)
				return reqMsg
			} else {
				reqMsg.Vaild = easygo.NewBool(false)
				return reqMsg
			}
		}
		pid := accountInfo.GetPlayerId()
		base := for_game.GetRedisPlayerBase(pid)
		if base == nil {
			//异常玩家数据，标志位冻结
			reqMsg.Vaild = easygo.NewBool(false)
			reqMsg.State = easygo.NewInt32(for_game.ACCOUNT_USER_FROZEN)
		} else {
			reqMsg.HeadIcon = easygo.NewString(base.GetHeadIcon())
			reqMsg.Sex = easygo.NewInt32(base.GetSex())
			reqMsg.Vaild = easygo.NewBool(true)
			reqMsg.State = easygo.NewInt32(base.GetStatus())
		}
	}
	return reqMsg
}

//取消注销账号
func (self *WebHttpForClient) RpcAccountCancel(common *base.Common, reqMsg *client_login.AccountCancel) easygo.IMessage {
	logs.Info("RpcAccountCancel", reqMsg)
	phone := reqMsg.GetAccount()
	if reqMsg.GetAccountType() == LOGIN_TYPE_WX {
		wechatToken, openId, unionId := for_game.GetWeChatInfo(reqMsg.GetAccount(), reqMsg.GetApkCode())
		if wechatToken == "" || openId == "" || unionId == "" {
			return easygo.NewFailMsg("无效的微信code")
		}
		info := for_game.GetPlayerInfoForUnionId(unionId)
		if info == nil {
			return easygo.NewFailMsg("无效的微信code")
		}
		phone = info.GetAccount()
		reqMsg.UnionId = easygo.NewString(unionId)
	} else if reqMsg.GetAccountType() == LOGIN_TYPE_ONEKEY {
		phone = for_game.GetJGOneKeyLoginPhone(reqMsg.GetAccount(), reqMsg.GetApkCode())
		if phone == "" {
			return easygo.NewFailMsg("一键登录失败")
		}
		reqMsg.UnionId = easygo.NewString(phone)
	}
	logs.Info("phone:", phone)
	accountInfo := for_game.GetRedisAccountByPhone(phone)
	if accountInfo == nil {
		return easygo.NewFailMsg("无效的手机号")
	}
	base := for_game.GetRedisPlayerBase(accountInfo.GetPlayerId())
	if base == nil {
		return easygo.NewFailMsg("账号数据异常")
	}
	if base.GetStatus() != for_game.ACCOUNT_CANCELING {
		return easygo.NewFailMsg("操作失败，账号已处于正常状态")
	}
	//设置为正常状态
	base.SetStatus(for_game.ACCOUNT_NORMAL)
	//修改注销账号订单表
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_CANCEL_ACCOUNT)
	defer closeFun()
	err := col.Update(bson.M{"PlayerId": base.GetPlayerId(), "Status": for_game.ACCOUNT_CANCEL_WAITING}, bson.M{"$set": bson.M{"Status": for_game.ACCOUNT_CANCEL_CANCEL}})
	easygo.PanicError(err)
	return reqMsg
}

func GetLoginMsg(pid PLAYER_ID, phone, token string, t int32) *client_login.LoginResult {
	info := PServerInfoMgr.GetIdelServer(for_game.SERVER_TYPE_HALL)
	code := client_login.LoginMark_LOGIN_SUCCESS
	var address string
	if for_game.IS_TFSERVER {
		address = easygo.AnytoA(info.GetClientWSPort())
	} else {
		address = for_game.MakeAddress(info.GetExternalIp(), info.GetClientWSPort())
	}
	player := for_game.GetRedisPlayerBase(pid)
	player.SetToken(token)
	msg := &client_login.LoginResult{
		Account:   easygo.NewString(phone),
		Result:    &code,
		Address:   easygo.NewString(address),
		Token:     easygo.NewString(token),
		LoginType: easygo.NewInt32(t),
		PlayerId:  easygo.NewInt64(pid),
	}
	return msg
}
*/
