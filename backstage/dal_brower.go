package backstage

import (
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/server_server"
	"game_server/pb/share_message"
	"log"
	"net"
	"strings"
	"time"

	"github.com/akqp2019/mgo"
	"github.com/akqp2019/mgo/bson"
	"github.com/astaxie/beego/logs"
)

//管理员登录验证
func UserLogin(account string, password string, ep IBrowerEndpoint) (*share_message.Manager, *base.Fail) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER)
	defer closeFun()
	user := &share_message.Manager{}

	errc := col.Find(bson.M{"Account": account}).One(user)
	if errc != nil && errc != mgo.ErrNotFound {
		panic(errc)
	}
	password = for_game.CreatePasswd(password, user.GetSalt())
	err := col.Find(bson.M{"Account": account, "Password": password}).One(user)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		rspMsg := fmt.Sprintf("账号: %s 登录失败，原因：账号密码错误", account)
		return nil, easygo.NewFailMsg(rspMsg)
	}
	if user.GetStatus() == 1 {
		rspMsg := fmt.Sprintf("账号: %s 登录失败，原因：账号被禁用", account)
		return user, easygo.NewFailMsg(rspMsg)
	}
	//登录IP限制
	if user.GetRole() != 0 && len(user.GetBindIp()) > 0 {
		pass := false
		lip := GetUserIP(ep)
		uips := user.GetBindIp()
		for _, p := range uips {
			if p == lip {
				pass = true
			}
		}

		if !pass {
			rspMsg := fmt.Sprintf("账号: %s 登录失败，原因：登录IP异常", account)
			return user, easygo.NewFailMsg(rspMsg)
		}
	}
	//顶号处理
	sid := PServerInfo.GetSid()
	admin := for_game.GetRedisAdmin(user.GetId())
	if admin != nil {
		if admin.ServerId == sid {
			ReplaceLogin(user.GetId())
		} else {
			//msg := &server_server.AdminInfo{
			//	UserId:   easygo.NewInt64(user.GetId()),
			//	ServerId: easygo.NewInt32(sid),
			//}
			msg := &server_server.PlayerSI{
				PlayerId: easygo.NewInt64(user.GetId()),
			}
			BroadCastMsgToServerNew(for_game.SERVER_TYPE_BACKSTAGE, "RpcReplaceLogin", msg)
			//ChooseOneHall(0, "RpcReplaceLoginToHall", msg)
		}
	}
	ep.SetUser(user)
	BrowerEpMp.StoreEndpoint(user.GetId(), ep.GetEndpointId())
	s := fmt.Sprintf("管理员登录服务器,ep=%d,userid=%d,username=%s\n", ep.GetEndpointId(), user.GetId(), user.GetAccount())
	logs.Info(s)

	timestamp := easygo.NowTimestamp()
	//判断是否是客服帐号
	if user.GetRole() == 2 {
		//查询活跃消息数量
		count := GetActiveIMmessageCount(user)
		obj := &for_game.RedisWaiter{
			UserId:    user.GetId(),
			Role:      1,
			ConnCount: count,
			ServerId:  sid,
			Status:    0,
			Types:     user.GetTypes(),
		}
		for_game.SetRedisWaiter(obj)
		easygo.AfterFunc(5*time.Second, func() {
			Waiters := for_game.GetRedisWaiterList()
			for _, v := range Waiters {
				log.Println("=============在线客服列表", v)
			}
		})

	} else {
		token := for_game.Md5(easygo.AnytoA(user.GetId()) + easygo.AnytoA(timestamp))
		obj := &for_game.RedisAdmin{
			UserId:    user.GetId(),
			Role:      user.GetRole(),
			ServerId:  sid,
			Timestamp: timestamp,
			Token:     token,
		}
		for_game.SetRedisAdmin(obj)
	}
	remoteAddr, _, err := net.SplitHostPort(ep.GetAddr().String()) //获取客户端IP
	easygo.PanicError(err)
	user.LoginCount = easygo.NewInt32(user.GetLoginCount() + 1)
	user.PrevLoginTime = user.LastLoginTime
	user.LastLoginTime = easygo.NewInt64(timestamp)
	user.PrevLoginIP = user.LastLoginIP
	user.LastLoginIP = easygo.NewString(remoteAddr)
	user.IsOnlie = easygo.NewBool(true)
	//更新登录数据
	EditManage(user.GetSite(), user, "login")
	NotifyPlayerOnLine(user) //通知其他后台服务器，用户上线
	return user, nil
}

//管理员登录api验证
func UserLoginApi(account, password, ip string) (*share_message.Manager, string, *base.Fail) {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER)
	defer closeFun()
	user := &share_message.Manager{}
	errc := col.Find(bson.M{"Account": account}).One(user)
	if errc != nil && errc != mgo.ErrNotFound {
		panic(errc)
	}
	password = for_game.CreatePasswd(password, user.GetSalt())
	err := col.Find(bson.M{"Account": account, "Password": password}).One(user)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		rspMsg := fmt.Sprintf("账号: %s 登录失败，原因：账号密码错误", account)
		return nil, "", easygo.NewFailMsg(rspMsg)
	}
	if user.GetStatus() == 1 {
		rspMsg := fmt.Sprintf("账号: %s 登录失败，原因：账号被禁用", account)
		return user, "", easygo.NewFailMsg(rspMsg)
	}
	//登录IP限制
	if user.GetRole() != 0 && len(user.GetBindIp()) > 0 {
		pass := false
		uips := user.GetBindIp()
		for _, p := range uips {
			if p == ip {
				pass = true
			}
		}

		if !pass {
			rspMsg := fmt.Sprintf("账号: %s 登录失败，原因：登录IP异常", account)
			return user, "", easygo.NewFailMsg(rspMsg)
		}
	}

	//顶号处理
	sid := PServerInfo.GetSid()
	admin := for_game.GetRedisAdmin(user.GetId())
	if admin != nil {
		if admin.ServerId == sid {
			ReplaceLogin(user.GetId())
		} else {
			msg := &server_server.AdminInfo{
				UserId:   easygo.NewInt64(user.GetId()),
				ServerId: easygo.NewInt32(sid),
			}
			ChooseOneHall(0, "RpcReplaceLoginToHall", msg)
		}
	}

	timestamp := easygo.NowTimestamp()
	token := for_game.Md5(easygo.AnytoA(user.GetId()) + easygo.AnytoA(timestamp))

	obj := &for_game.RedisAdmin{
		UserId:    user.GetId(),
		Role:      user.GetRole(),
		ServerId:  sid,
		Timestamp: timestamp,
		Token:     token,
	}
	for_game.SetRedisAdmin(obj)

	user.LoginCount = easygo.NewInt32(user.GetLoginCount() + 1)
	user.PrevLoginTime = user.LastLoginTime
	user.LastLoginTime = easygo.NewInt64(timestamp)
	user.PrevLoginIP = user.LastLoginIP
	user.LastLoginIP = easygo.NewString(ip)
	user.IsOnlie = easygo.NewBool(true)
	//更新登录数据
	EditManage(user.GetSite(), user, "login")
	return user, token, nil
}

//获取谷歌验证器秘钥
func SetGoogleSecret(ep IBrowerEndpoint, user *share_message.Manager) string {
	secret := for_game.NewGoogleAuth().GetSecret()
	user.GoogleSecret = easygo.NewString(secret)
	EditManage(user.GetSite(), user, "")
	var ms string
	if user.GetGoogleSecret() == "" {
		ms = "生成谷歌验证器秘钥"
	} else {
		ms = "重置谷歌验证器秘钥"
	}
	logs.Info(ms)
	return secret
}

//顶号处理
func ReplaceLogin(id USER_ID) {
	ep := BrowerEpMp.LoadEndpoint(id)
	if ep != nil {
		msg := &brower_backstage.ErrMessage{Err: easygo.NewString("帐号已在别处登录,请确认帐号安全")}
		ep.RpcReplacePush(msg) //通知前端下线
		ep.Shutdown()          //强制下线
	}
}

//id查询管理员
func GetUser(userid USER_ID) *share_message.Manager {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER)
	defer closeFun()
	user := &share_message.Manager{}

	errc := col.Find(bson.M{"_id": userid}).One(user)
	if errc != nil && errc != mgo.ErrNotFound {
		panic(errc)
	}
	if errc == mgo.ErrNotFound {
		return nil
	}
	return user
}

func GetUserByIds(ids []int64) []*share_message.Manager {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER)
	defer closeFun()

	users := []*share_message.Manager{}
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&users)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return users
}

//ID查询账号基础信息 Player_base
func GetPlayerBaseBackstageById(site string, playerid int64) *share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	player := &share_message.PlayerBase{}
	err := col.Find(bson.M{"_id": playerid}).One(player)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return player
}

//层级id查询账号基础信息列表 Player_base
func GetPlayerBaseByLayer(site string, id int32) []*share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	playerList := []*share_message.PlayerBase{}
	err := col.Find(bson.M{"Layer": id}).All(&playerList)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return playerList
}

//ID查询账号基础信息 Player_base
func GetPlayerBaseBkById(site string, playerid int64) *share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	player := &share_message.PlayerBase{}
	err := col.Find(bson.M{"_id": playerid}).One(player)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return player
}

//ID查询账号基础信息 Player_base
func GetPlayerBaseByIds(playerids []int64) []*share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	players := []*share_message.PlayerBase{}
	err := col.Find(bson.M{"_id": bson.M{"$in": playerids}}).All(&players)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return players
}

//昵称模糊查询账号信息
func GetPlayerLikeNickname(site string, nickname string) []*share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	/*player := &share_message.PlayerBase{}
	err := col.Find(bson.M{"NickName": nickname}).One(player)*/

	player := []*share_message.PlayerBase{}
	err := col.Find(bson.M{"NickName": bson.M{"$regex": bson.RegEx{Pattern: nickname, Options: "im"}}}).All(&player)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}

	return player
}

//昵称查询账号信息   *待优化 昵称支持唯一性
func QuryPlayerByNickname(nickname string) *share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	player := &share_message.PlayerBase{}
	err := col.Find(bson.M{"NickName": nickname}).One(player)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return player
}

//根据登录IP查询玩家列表
func QuryPlayerByIp(site string, ip string) []*share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	player := []*share_message.PlayerBase{}
	err := col.Find(bson.M{"LastLoginIP": ip}).All(&player)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return player
}

//根据登录IP查询玩家
func QuryPlayerById(site string, id int64) *share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	player := &share_message.PlayerBase{}
	err := col.Find(bson.M{"_id": id}).One(&player)
	if err != nil && err != mgo.ErrNotFound {
		return nil
	}
	if err == mgo.ErrNotFound {
		return nil
	}
	return player
}

// 根据Id 删除mongo 里面的数据
func DelDataById(document string, idList []int64) error {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, document)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"_id": bson.M{"$in": idList}})
	return err
}

//获取用户IP
func GetUserIP(ep IBrowerEndpoint) string {
	remoteAddr, _, err := net.SplitHostPort(ep.GetAddr().String())
	easygo.PanicError(err)
	return remoteAddr
}

// 编辑冻结IP
func EditFreezeIpList(site string, reqMsg *share_message.FreezeIpList) {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_FREEZEIP)
	defer closeFun()

	_, err := col.Upsert(bson.M{"_id": reqMsg.GetId()}, bson.M{"$set": reqMsg})
	easygo.PanicError(err)
}

// 新增冻结IP
func ListUpFreezeIpList(site string, reqMsg []*share_message.FreezeIpList) {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_FREEZEIP)
	defer closeFun()
	// 转换类型并存到新数组
	var il []interface{}
	for _, rr := range reqMsg {
		il = append(il, rr)
	}
	// 批量插入数据库
	err := col.Insert(il...)
	easygo.PanicError(err)
}

//批量删除冻结IP
func BatchDelFreezeIp(site string, id []int64) {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_FREEZEIP)
	defer closeFun()
	_, err := col.RemoveAll(bson.M{"_id": bson.M{"$in": id}})
	easygo.PanicError(err)
}

// 获取今日注册人数
func GetRegisterNumber(site string) int32 {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_PLAYER_BASE)
	defer closeFun()

	//获取时间戳
	startTime := easygo.GetToday0ClockTimestamp()
	count, err := col.Find(bson.M{"CreateTime": bson.M{"$gte": startTime}}).Count()

	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}

	return int32(count)
}

// 获取指定时间范围注册的玩家列表
func GetRegisterPlayerListByTime(site string, startTime int64, endTime int64) []*share_message.PlayerBase {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_PLAYER_BASE)
	defer closeFun()
	list := []*share_message.PlayerBase{}
	err := col.Find(bson.M{"CreateTime": bson.M{"$gte": startTime, "$lte": endTime}}).All(&list)

	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}

	return list
}

// 查询在线状态
func GetStatus(site string, status int32) ([]*share_message.PlayerBase, int) {
	col, closeFun := MongoMgr.GetC(site, for_game.TABLE_PLAYER_BASE)
	defer closeFun()
	list := []*share_message.PlayerBase{}
	err := col.Find(bson.M{"IsOnlie": status, "Status": 1}).All(&list)

	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}

	count := len(list)

	return list, count
}

//根据用户权限判断是否全部显示手机号码
func SetUserPhone(phone string) string {
	if len(phone) < 11 {
		return phone
	}
	old := ""
	for k, v := range phone {
		if k > 4 {
			old = old + string(v)
		}
	}
	phone = strings.Replace(phone, old, "******", -1)
	//slice := strings.Split(phone, "")
	//str := strings.Join(slice[0:5], "") + "******"
	return phone
}

//判断用户是否有权限去查看完整的手机号码
func VailRole(site string) bool {
	col, closeFun := MongoMgr.GetC(for_game.MONGODB_NINGMENG, for_game.TABLE_MANAGER)
	defer closeFun()
	user := &share_message.Manager{}
	err := col.Find(bson.M{"Site": site}).One(&user)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		return false
	}
	if user.GetRole() > 1 {
		return false
	}
	return true
}
