//后台管理相关

package backstage

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"net/http"
)

func (self *WebHttpServer) BackstageEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	// w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept,WG-App-Version, WG-Device-Id, WG-Network-Type, WG-Vendor, WG-OS-Type, WG-OS-Version, WG-Device-Model, WG-CPU, WG-Sid, WG-App-Id, WG-Token")
	t := r.PostFormValue("t")

	switch t {
	case "1": //登录后台
		ip := for_game.GetUserIp(r)
		user := r.PostFormValue("user")
		pass := r.PostFormValue("pass")
		self.Login(w, user, pass, ip)
	case "2": //查询支付设置
		id := r.PostFormValue("id")
		token := r.PostFormValue("token")
		self.QueryPaySetting(w, id, token)
	case "3": //修改支付设置
		id := r.PostFormValue("id")
		token := r.PostFormValue("token")
		data := r.PostFormValue("data")
		self.UpdatePaySetting(w, id, token, data)
	default:
		OutputJson(w, 0, "fail", nil)
	}
}

//登录
func (self *WebHttpServer) Login(w http.ResponseWriter, user, pass, ip string) {
	admin, token, err := UserLoginApi(user, pass, ip)
	if err != nil {
		OutputJson(w, 0, err.GetReason(), nil)
		return
	}

	if !admin.GetIsLoginH5() {
		rspMsg := fmt.Sprintf("账号: %s 登录失败，原因：权限不足", user)
		OutputJson(w, 0, rspMsg, nil)
		return
	}

	result := &brower_backstage.LoginResponse{
		User:  admin,
		Token: easygo.NewString(token),
	}

	AddBackstageLog(user, ip, for_game.LOGIN_BACKSTAGE, "管理员登录H5")
	// logs.Debug("登录用户:", for_game.GetRedisAdmin(admin.GetId()))
	OutputJson(w, 1, "success", result)
}

//查询支付设置
func (self *WebHttpServer) QueryPaySetting(w http.ResponseWriter, id, token string) {
	// logs.Debug(fmt.Sprintf("QueryPaySetting==>id:%s,token:%s", id, token))
	uid := easygo.AtoInt64(id)
	tokenResult := self.CheckToken(uid, token)
	if !tokenResult {
		OutputJson(w, 0, "登录过期,请重新登录", nil)
		return
	}

	result := for_game.QuerySysParameterById(for_game.LIMIT_PARAMETER)
	if result == nil {
		OutputJson(w, 0, "找不到配置", nil)
		return
	}

	date := &share_message.SysParameter{
		Id:           easygo.NewString(result.GetId()),
		IsTransfer:   easygo.NewBool(result.GetIsTransfer()),
		IsRedPacket:  easygo.NewBool(result.GetIsRedPacket()),
		IsRecharge:   easygo.NewBool(result.GetIsRecharge()),
		IsWithdrawal: easygo.NewBool(result.GetIsWithdrawal()),
		IsQRcode:     easygo.NewBool(result.GetIsQRcode()),
	}

	OutputJson(w, 1, "success", date)
}

//修改支付设置
func (self *WebHttpServer) UpdatePaySetting(w http.ResponseWriter, id, token, data string) {
	// logs.Debug(fmt.Sprintf("UpdatePaySetting==>id:%s,token:%s,data:%s", id, token, data))
	uid := easygo.AtoInt64(id)
	tokenResult := self.CheckToken(uid, token)
	if !tokenResult {
		OutputJson(w, 0, "token错误,请重新登录", nil)
		return
	}

	msg := &share_message.SysParameter{}
	err := json.Unmarshal([]byte(data), &msg)
	easygo.PanicError(err)

	result := EditSysParameter(msg)
	if result != nil {
		OutputJson(w, 0, result.GetReason(), nil)
		return
	}

	OutputJson(w, 1, "success", nil)
}

//检查token
func (self *WebHttpServer) CheckToken(uid int64, token string) bool {
	admin := for_game.GetRedisAdmin(uid)
	if admin == nil {
		return false
	}

	if easygo.NowTimestamp()-admin.Timestamp > 600 {
		return false
	} else {
		if token != admin.Token {
			return false
		}
	}

	return true
}
