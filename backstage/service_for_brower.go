// 管理后台为[浏览器]提供的服务

package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"path"
	"time"

	"github.com/astaxie/beego/logs"
)

type ServiceForBrower struct {
}

type cls4 = ServiceForBrower

//电竞服开启
//上传文件到存储桶
func (self *cls4) RpcUploadFile(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.UploadRequest) easygo.IMessage {
	fileName := reqMsg.GetFileName()
	upPath := reqMsg.GetPath()
	file := reqMsg.GetFile()
	isBucket := reqMsg.GetIsBucket()
	// file, _ := base64.StdEncoding.DecodeString(files)
	var url string
	//上传类型 1-本地文件上传,2-网络文件上传
	switch reqMsg.GetType() {
	case 1:
		if isBucket {
			pathfileName := path.Join("backstage", upPath, fileName)
			url = QQbucket.ObjectPutByte(pathfileName, file)
			if url == "" {
				return easygo.NewFailMsg("上传存储桶失败")
			}
		} else {
			url, _, _ = easygo.UploadFile(fileName, file)
		}
	case 2:
		if reqMsg.FileUrl == nil || reqMsg.GetFileUrl() == "" {
			return easygo.NewFailMsg("网络文件地址不能为空")
		}
		if isBucket {
			pathfileName := path.Join("backstage", upPath, fileName)
			url = QQbucket.ObjectPutRemote(pathfileName, reqMsg.GetFileUrl())
			if url == "" {
				return easygo.NewFailMsg("上传存储桶失败")
			}
		} else {
			//get方法获取资源
			util.DownloadFile(reqMsg.GetFileUrl(), fileName, upPath)
			url = path.Join(upPath, fileName)
		}
	}

	return &brower_backstage.UploadResponse{
		Url: easygo.NewString(url),
	}
}

//删除存储桶中的文件
func (self *cls4) RpcDelUploadFile(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.UploadRequest) easygo.IMessage {
	fileName := reqMsg.GetFileName()
	if fileName == "" {
		return easygo.NewFailMsg("文件名不能为空")
	}
	upPath := reqMsg.GetPath()
	if upPath == "" {
		upPath = "upload"
	}
	isBucket := reqMsg.GetIsBucket()
	if reqMsg.IsBucket == nil {
		isBucket = true
	}

	if isBucket {
		pathfileName := path.Join("backstage", upPath, fileName)
		QQbucket.ObjectDel(pathfileName)
	} else {
		err := util.DeleteFile(fileName, upPath)
		if err != nil {
			return easygo.NewFailMsg(err.Error())
		}
	}

	return easygo.EmptyMsg
}

//获取存储桶中的文件列表
func (self *cls4) RpcUploadFileList(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	paths := reqMsg.GetSrtType()
	marker := reqMsg.GetKeyword()
	pagesize := reqMsg.GetPageSize()
	if paths == "" {
		paths = "upload"
	}
	paths = path.Join("backstage", paths)
	lis := QQbucket.GetObjectList(paths, marker, pagesize)
	var list []*brower_backstage.UploadList
	for _, li := range lis {
		one := &brower_backstage.UploadList{}
		for_game.StructToOtherStruct(li, one)
		list = append(list, one)
	}

	return &brower_backstage.UploadListResponse{
		List: list,
	}
}

//谷歌验证器
func (self *cls4) RpcGoogleCode(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.GoogleCodeRequest) easygo.IMessage {
	msg := &brower_backstage.GoogleCodeResponse{}
	switch reqMsg.GetType() {
	case 1: //1获取设置信息，
		secret := user.GetGoogleSecret()
		if secret == "" {
			secret = SetGoogleSecret(ep, user)
		}
		qrcodeUrl := for_game.NewGoogleAuth().GetQrcode(user.GetAccount(), secret)
		msg = &brower_backstage.GoogleCodeResponse{
			Status:    user.IsGoogleVer,
			Secret:    easygo.NewString(secret),
			QrcodeUrl: easygo.NewString(qrcodeUrl),
		}
	case 2: //2重置生成秘钥,
		secret := SetGoogleSecret(ep, user)
		qrcodeUrl := for_game.NewGoogleAuth().GetQrcode(user.GetAccount(), secret)
		msg = &brower_backstage.GoogleCodeResponse{
			Status:    user.IsGoogleVer,
			Secret:    easygo.NewString(secret),
			QrcodeUrl: easygo.NewString(qrcodeUrl),
		}
	case 3: //3开关
		if user.GetRole() <= 1 && user.GetSite() != for_game.MONGODB_NINGMENG {
			// site := for_game.GetSiteConfig().GetSiteByName(user.GetSite())
			// site.IsGoogleVer = reqMsg.Status
			// for_game.GetSiteConfig().EditSiteList(site)
			// RpcHallReloadConfig(for_game.SiteConfigName, user.GetSite()) //通知大厅
			var ms string
			if reqMsg.GetStatus() {
				ms = "开启全站谷歌验证器"
			} else {
				ms = "关闭全站谷歌验证器"
			}
			logs.Info(ms)
		}

		admin := GetUser(user.GetId())
		admin.IsGoogleVer = reqMsg.Status
		EditManage(admin.GetSite(), admin, "")
		var ms string
		if reqMsg.GetStatus() {
			ms = "开启谷歌验证器"
		} else {
			ms = "关闭谷歌验证器"
		}
		logs.Info(ms)
		qrcodeUrl := for_game.NewGoogleAuth().GetQrcode(admin.GetAccount(), admin.GetGoogleSecret())
		msg = &brower_backstage.GoogleCodeResponse{
			Status:    admin.IsGoogleVer,
			Secret:    admin.GoogleSecret,
			QrcodeUrl: easygo.NewString(qrcodeUrl),
		}
	}

	return msg
}

//获取验证码
func (self *cls4) RpcGetCode(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.SigninRequest) easygo.IMessage {
	phone := reqMsg.GetPhone()
	t := reqMsg.GetTypes()
	if reqMsg.Phone == nil || phone == "" {
		return easygo.NewFailMsg("手机号不能为空")
	}

	if reqMsg.Types == nil || t == 0 {
		return easygo.NewFailMsg("验证码类型不能为空")
	}

	data := for_game.MessageMarkInfo.GetMessageMarkInfo(t, phone)
	if data != nil {
		leaveTime := time.Now().Unix() - data.Timestamp
		if leaveTime <= 120 {
			return easygo.NewFailMsg("验证码已发送!")
		}
	}

	codes := for_game.SendCodeToClientUser(t, phone, reqMsg.GetAreaCode())
	if codes == "" {
		return easygo.NewFailMsg("验证码发送失败!")
	}

	msg := &brower_backstage.CodeResponse{
		Code: easygo.NewString(codes),
	}
	return msg
}

//特殊界面获取验证码
func (self *cls4) RpcGetBsCode(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.SigninRequest) easygo.IMessage {
	phone := reqMsg.GetPhone()
	t := reqMsg.GetTypes()
	if reqMsg.Phone == nil || phone == "" {
		return easygo.NewFailMsg("手机号不能为空")
	}

	if phone != user.GetPhone() {
		return easygo.NewFailMsg("手机号码权限不足,请联系管理员")
	}

	if reqMsg.Types == nil || t == 0 {
		return easygo.NewFailMsg("验证码类型不能为空")
	}

	data := for_game.MessageMarkInfo.GetMessageMarkInfo(t, phone)
	if data != nil {
		leaveTime := time.Now().Unix() - data.Timestamp
		if leaveTime <= 120 {
			return easygo.NewFailMsg("验证码已发送!")
		}
	}

	codes := for_game.SendCodeToClientUser(t, phone, reqMsg.GetAreaCode())
	if codes == "" {
		return easygo.NewFailMsg("验证码发送失败!")
	}

	msg := &brower_backstage.CodeResponse{
		Code: easygo.NewString(codes),
	}
	return msg
}

//特殊界面校验验证码
func (self *cls4) RpcVerCode(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.SigninRequest) easygo.IMessage {
	if reqMsg.Code == nil {
		return easygo.NewFailMsg("验证码不能为空")
	}
	if reqMsg.Phone == nil {
		return easygo.NewFailMsg("手机号不能为空")
	}
	if reqMsg.Types == nil {
		return easygo.NewFailMsg("验证码类型不能为空")
	}
	data := for_game.MessageMarkInfo.GetMessageMarkInfo(reqMsg.GetTypes(), reqMsg.GetPhone())
	if data == nil {
		res := "验证码不存在"
		return easygo.NewFailMsg(res)
	}
	if data.Mark != reqMsg.GetCode() {
		res := "验证码不正确"
		return easygo.NewFailMsg(res)
	}
	return easygo.EmptyMsg
}

//注册
func (self *cls4) RpcSignin(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.SigninRequest) easygo.IMessage {
	// site := reqMsg.GetSite()
	// player := GetPlayerByPhone(site, reqMsg.GetPhone())
	// if player != nil {
	// 	return easygo.NewFailMsg("该手机号已使用")
	// }

	// parent := &share_message.PlayerBase{}

	// // 判断是否是官方发放的成为一级代理的二维码
	// if reqMsg.GetParentId() != 0 {
	// 	parent = GetPlayerBaseById(site, reqMsg.GetParentId()) //上级代理base
	// }

	// remoteAddr := GetUserIP(ep)
	// createIP := easygo.NewString(remoteAddr)
	// createTime := easygo.NewInt64(time.Now().Unix())

	// playerBase := &share_message.PlayerBase{
	// 	PlayerId: easygo.NewInt64(0),
	// 	// HeadIcon:           easygo.NewInt32(1),    //头像
	// 	Sex:  easygo.NewInt32(2),   //性别
	// 	Gold: easygo.NewFloat64(0), //玩家携带的金币
	// 	// Vip:                easygo.NewInt32(1),    //Vip等级
	// 	// Layer:              easygo.NewInt32(1),    //层级
	// 	IsRobot: easygo.NewBool(false), //是否是机器人
	// 	// IsOnline:           easygo.NewInt32(0),    //玩家状态 0 离线, 1 闲置, 2 游戏中
	// 	// HeadIconFrame:      easygo.NewInt32(1),    //玩家头相框
	// 	// NickNameModifyTime: createTime, //昵称修改时的时间戳
	// 	// SafeboxGold:        easygo.NewFloat64(0), //保险箱金币
	// 	// SafeboxPwd:         easygo.NewString(""), //保险箱密码

	// 	// VipAdvanceTime:          createTime, //vip晋级时间
	// 	// GetedUpRewardVipLvSlice: []int32{},  //已领vip晋级奖励的等级的列表
	// 	CreateTime: createTime, //创建时间

	// 	// LoginDays:       easygo.NewInt32(0), //登录天数
	// 	// LoginTimes:      easygo.NewInt32(0), //登录总次数
	// 	// TodayLoginTimes: easygo.NewInt32(0), //今日登录次数
	// 	// OnlineTime:      easygo.NewInt32(0), //在线总时长  单位分钟
	// 	// TodayOnlineTime: easygo.NewInt32(0), //当天在线时长  单位分钟

	// 	// LastLoginTime:       easygo.NewInt64(0),   //最后登录时间
	// 	// LastLoginIP: createIP, //最后登录IP
	// 	// LastLoginDeviceCode: easygo.NewString(""), //最后登录机器码
	// }

	// notifymsg := &backstage_notify.CreatePlayer{
	// 	Site:       easygo.NewString(site),
	// 	PlayerBase: playerBase,
	// 	ParentId:   parent.PlayerId,
	// 	CreateIp:   createIP,
	// 	Phone:      easygo.NewString(reqMsg.GetPhone()),
	// }

	// playerId := CreatePlayerForNotify(notifymsg) //告诉通知服务器创建帐号
	// playerExtend.PlayerId = easygo.NewInt64(playerId)
	// SetPlayerExtend(site, playerExtend) //写玩家扩展表

	// siteCfg := for_game.GetSiteConfig().GetSiteByName(site)

	//下载网址重定向
	msg := &brower_backstage.SigninResponse{}

	return msg
}

//登录
func (self *cls4) RpcLogin(ep IBrowerEndpoint, ctx interface{}, reqMsg *brower_backstage.LoginRequest) easygo.IMessage {
	userAccount := reqMsg.GetUserAccount()
	password := reqMsg.GetPassword()
	admin, err := UserLogin(userAccount, password, ep)
	if err != nil {
		return easygo.NewFailMsg(err.GetReason())
	}

	if admin.GetIsGoogleVer() {
		code := reqMsg.GetCode()
		secret := admin.GetGoogleSecret()
		result, errr := for_game.NewGoogleAuth().VerifyCode(secret, code)
		easygo.PanicError(errr)
		if !result {
			rspMsg := fmt.Sprintf("账号: %s 登录失败，原因：谷歌身份证失败", userAccount)
			return easygo.NewFailMsg(rspMsg)
		}
	}

	roletype := &share_message.RolePower{}
	if admin.GetRole() > 0 {
		roletype = GetPowerRouter(admin.GetRoleType())
		if roletype == nil {
			rspMsg := fmt.Sprintf("账号: %s 登录失败，原因：权限不足", userAccount)
			return easygo.NewFailMsg(rspMsg)
		}
	}

	result := &brower_backstage.LoginResponse{
		User:  admin,
		Power: roletype,
	}

	AddBackstageLog(userAccount, GetUserIP(ep), for_game.LOGIN_BACKSTAGE, "管理员登录")
	return result
}

//登出
func (self *cls4) RpcLogout(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	mp := BrowerEpMp.LoadEndpoint(user.GetId())
	if mp != nil {
		mp.Shutdown()
	}
	AddBackstageLog(user.GetAccount(), GetUserIP(ep), for_game.SIGNOUT_BACKSTAGE, "管理员登出")
	return easygo.EmptyMsg
}

//查询屏蔽分数
func (self *cls4) RpcCheckShieldScore(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.CheckScoreRequest) easygo.IMessage {
	msg := CheckOjbScore(reqMsg.GetKey(), reqMsg.GetType())
	if msg == nil {
		return easygo.NewFailMsg("系统正忙，请稍后再试")
	}

	return msg
}
