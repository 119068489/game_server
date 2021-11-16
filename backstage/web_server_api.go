// 为[浏览器]提供的API服务

package backstage

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
)

//结果返回
type Result struct {
	Ret    int
	Reason string
	Data   interface{}
}

//错误码结果
type ErrorResult struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

const API_MD5_KEY = "nmcl2020!@#$"

type WebHttpServer struct {
	Service reflect.Value
}

func NewWebHttpServer() *WebHttpServer {
	p := &WebHttpServer{}
	p.Init()
	return p
}

func (self *WebHttpServer) Init() {}

func (self *WebHttpServer) SendErrorCode(w http.ResponseWriter, errCode int, errText string) {
	//输入记录
	result := &ErrorResult{
		Code:  errCode,
		Error: errText,
	}
	content, err := json.Marshal(result)
	easygo.PanicError(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

//GET请求
func httpGet(url string) []byte {
	res, err := http.Get(url)
	easygo.PanicError(err)
	result, err1 := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	easygo.PanicError(err1)
	return result
}

//返回json
func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := &Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	w.Write(b)
}

func (self *WebHttpServer) Serve() {
	//address := for_game.MakeAddress(PServerInfo.GetIp(), PServerInfo.GetBackStageApiPort())
	//logs.Info("后台端口：", PServerInfo.GetBackStageApiPort())
	address := for_game.MakeAddress("0.0.0.0", PServerInfo.GetBackStageApiPort())

	logs.Info("(API 服务) 开始监听: %v", address)
	http.HandleFunc("/version", self.VersionEntry)           //动态跳转app下载包地址
	http.HandleFunc("/tfserver", self.TfserverEntry)         //转发服务器列表
	http.HandleFunc("/shorturl", self.ShortUrlEntry)         //生成短网址api
	http.HandleFunc("/downpage", self.DownPageEntry)         //落地页资源api
	http.HandleFunc("/pagedownload", self.PageDownLoadEntry) //落地页下载api
	http.HandleFunc("/article", self.ArticleEntry)           //文章api
	http.HandleFunc("/DeviceCode", self.DeviceCodeEntry)     //保存设备api
	http.HandleFunc("/picture", self.CheckPicture)           //违禁图回调
	http.HandleFunc("/teamdata", self.TeamdataEntry)         //查询群资料
	http.HandleFunc("/registered", self.RegisteredEntry)     //注册api 参数t=1获取验证码,2注册 ，phone，code,channelno
	http.HandleFunc("/freezeLog", self.FreezeLogEntry)       //用户冻结日志 参数pagesize，curpage
	http.HandleFunc("/square", self.SquareEntry)             //查询社交广场动态信息
	http.HandleFunc("/idfa", self.SaveIdfaEntry)             //保存设备idfa码
	http.HandleFunc("/checkidfa", self.CheckIdfaEntry)       //检查设备idfa码是否存在
	http.HandleFunc("/activity", self.ActivityEntry)         //活动API
	http.HandleFunc("/advidfa", self.AdvIdfaEntry)           //广告idfa上报API
	http.HandleFunc("/backstage", self.BackstageEntry)       //后台管理api
	http.HandleFunc("/uprecall", self.UpRecallEntry)         //召回用户api
	http.HandleFunc("/youmi", self.YourMiEntry)              //有米api接口
	http.HandleFunc("/upload", self.UploadEntry)             //上传文件
	http.HandleFunc("/sharetopic", self.ShareTopicEntry)     //分享话题
	err := http.ListenAndServe(address, nil)                 // 第 2 个参数可以传 nil。传进去的 self 必须有实现一个 ServeHTTP 函数
	easygo.PanicError(err)
}

//获取客户端版本配置
func (self *WebHttpServer) VersionEntry(w http.ResponseWriter, r *http.Request) {
	jsonFile := easygo.YamlCfg.GetValueAsString("CLIENT_VERSION_DATA")
	data := for_game.FindVersion(jsonFile)
	out, _ := json.MarshalIndent(data, "", "  ")
	// logs.Info("返回数据", string(out))
	w.Write(out)
}

//获取转发服务器列表
func (self *WebHttpServer) TfserverEntry(w http.ResponseWriter, r *http.Request) {
	jsonFile := easygo.YamlCfg.GetValueAsString("CLIENT_TFSERVER_ADDR")
	data := for_game.FindTfserver(jsonFile)
	out, _ := json.MarshalIndent(data, "", "  ")
	// logs.Info("返回数据", string(out))
	w.Write(out)
}

//短网址api
func (self *WebHttpServer) ShortUrlEntry(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	params := r.Form
	key := params.Get("url")
	if key == "" {
		msg := "Url Param 'url' is missing"
		OutputJson(w, 0, msg, nil)
		return
	}

	longurl := url.QueryEscape(key)
	apiurl := "https://api.d5.nz/api/dwz/tcn.php"
	qurl := apiurl + "?url=" + longurl
	result := httpGet(qurl)
	OutputJson(w, 1, "success", result)
}

//落地页资源请求api
func (self *WebHttpServer) DownPageEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ip := r.RemoteAddr
	params := r.Form
	id := params.Get("id")

	bList := []string{"http://www.lemonchat.cn/img/downpage/1.jpg", "http://www.lemonchat.cn/img/downpage/2.jpg", "http://www.lemonchat.cn/img/downpage/3.png"}
	dp := &share_message.DownPage{
		ModId:   easygo.NewInt32(1),
		Icon:    easygo.NewString("http://www.lemonchat.cn/img/downpage/top.png"),
		Banner:  bList,
		BtnText: easygo.NewString("立刻下载"),
		Floot:   easygo.NewString(""),
	}

	if id != "" {
		id := url.QueryEscape(id)
		op_cl := for_game.QueryOperationByNo(id)
		if op_cl != nil {
			dp = op_cl.DpSet
			if !for_game.IsNewIp(ip) {
				easygo.Spawn(func() { for_game.MakeOperationChannelReport(3, 0, id, nil, nil) }) //生成运营渠道数据汇总报表 已优化到Redis
			}

		}
	}

	out, _ := json.MarshalIndent(dp, "", "  ")
	// logs.Info("返回数据", string(out))
	w.Write(out)
}

//落地页下载api
func (self *WebHttpServer) PageDownLoadEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	params := r.Form
	id := params.Get("id")
	if id != "" {
		id := url.QueryEscape(id)
		op_cl := for_game.QueryOperationByNo(id)
		if op_cl != nil {
			easygo.Spawn(func() { for_game.MakeOperationChannelReport(2, 0, id, nil, nil) }) //生成运营渠道数据汇总报表 已优化到Redis
		}
	}
}

//请求验签
func (self *WebHttpServer) ApiCheckSign(params url.Values) bool {
	sign := params.Get("sign")
	// logs.Info("================>前端签名", sign)
	params.Del("sign")
	src := params.Encode()
	newSign := for_game.Md5(src + API_MD5_KEY)
	// logs.Info("================>服务器签名", src+API_MD5_KEY, newSign)
	return sign == newSign
}

//文章api
//url?t=1&id=1&pid=1885040404&sign=md5(所传参数askill排序后的值+md5key)
//pid=1885040404&id=1&t=1nmcl2020!@#$
func (self *WebHttpServer) ArticleEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	params := r.Form
	content := params.Get("content")
	params.Del("content")
	// if !self.ApiCheckSign(params) {
	// 	OutputJson(w, 0, "签名错误", nil)
	// 	return
	// }
	t := params.Get("t") //请求数据类型
	switch t {
	case for_game.ARTICLE_API_READ:
		self.ApiArticleRead(w, params)
	case for_game.ARTICLE_API_ZAN:
		self.ApiArticleZan(w, params)
	case for_game.ARTICLE_API_COMMENT:
		params.Add("content", content)
		self.ApiArticleComment(w, params)
	case for_game.ARTICLE_API_GET_COMMENT:
		self.ApiArticleGetComment(w, params)
	}

}

//读取文章
func (self *WebHttpServer) ApiArticleRead(w http.ResponseWriter, params url.Values) {
	id := params.Get("id")
	pid := params.Get("pid")
	time := params.Get("time")
	var articleId, playerId, timev int64
	if id != "" {
		articleId, _ = strconv.ParseInt(id, 10, 64)
	}
	if pid != "" {
		playerId, _ = strconv.ParseInt(pid, 10, 64)
	}
	if time != "" {
		timev, _ = strconv.ParseInt(time, 10, 64)
	}
	articleInfo := for_game.ReadArticle(playerId, articleId)
	if timev > 0 {
		articleInfo.CreateTime = easygo.NewInt64(timev)
	}
	if articleInfo == nil || articleInfo.GetState() == 1 {
		OutputJson(w, 0, "文章不存在", nil)
	} else {
		OutputJson(w, 1, "success", articleInfo)

	}
}

//赞文章
func (self *WebHttpServer) ApiArticleZan(w http.ResponseWriter, params url.Values) {
	id := params.Get("id")
	pid := params.Get("pid")
	var articleId, playerId int64
	if id == "" || pid == "" {
		OutputJson(w, 0, "id or pid is null", nil)
		return
	}
	articleId, _ = strconv.ParseInt(id, 10, 64)
	playerId, _ = strconv.ParseInt(pid, 10, 64)
	for_game.ZanArticle(playerId, articleId)
	OutputJson(w, 1, "success", nil)
}

//评论文章
func (self *WebHttpServer) ApiArticleComment(w http.ResponseWriter, params url.Values) {
	id := params.Get("id")
	pid := params.Get("pid")
	content := params.Get("content")
	var articleId, playerId int64
	if id != "" {
		articleId, _ = strconv.ParseInt(id, 10, 64)
	}
	if pid != "" {
		playerId, _ = strconv.ParseInt(pid, 10, 64)
	}
	key, _ := base64.StdEncoding.DecodeString(content)
	isDirty, _ := for_game.PDirtyWordsMgr.CheckWord(string(key))
	if isDirty {
		OutputJson(w, 2, "包含敏感词,请重新编辑", nil)
		return
	}
	for_game.CommentArticle(playerId, articleId, content)
	OutputJson(w, 1, "success", nil)
}

//获取文章评论内容
func (self *WebHttpServer) ApiArticleGetComment(w http.ResponseWriter, params url.Values) {
	id := params.Get("id")
	pageSize := params.Get("pagesize")
	curPage := params.Get("curpage")
	var articleId int64
	var pagesize int64
	var curpage int64
	if id != "" {
		articleId, _ = strconv.ParseInt(id, 10, 64)
	}
	if pageSize != "" {
		pagesize, _ = strconv.ParseInt(pageSize, 10, 64)
	}
	if curPage != "" {
		curpage, _ = strconv.ParseInt(curPage, 10, 64)
	}
	comments := for_game.GetArticleComment(articleId, pagesize, curpage, for_game.ARTICLE_COMMENT_SHOW)
	OutputJson(w, 1, "success", comments)
}

//保存设备api 参数code,channle //上报激活设备
func (self *WebHttpServer) DeviceCodeEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	params := r.Form
	key := params.Get("code")
	idfa := params.Get("idfa")
	if key != "" {
		cal := params.Get("channle")
		channle := url.QueryEscape(cal)
		code := url.QueryEscape(key)
		msg := &share_message.PosDeviceCode{
			CreateTime: easygo.NewInt64(util.GetMilliTime()),
			DeviceCode: easygo.NewString(code),
			Channle:    easygo.NewString(channle),
		}
		for_game.SavePosDeviceCode(msg)
		easygo.Spawn(for_game.UpdatePosDeviceAdvIdfa, code, "IsActive", true, idfa) //广告设备激活
		if idfa != "" {
			matched, err := regexp.MatchString("^[A-Z0-9-]+$", idfa)
			cal := params.Get("source")
			source := ""
			if cal != "" {
				source = url.QueryEscape(cal)
			}

			if matched == true && err == nil {
				msg := &share_message.PosDeviceIdfa{
					CreateTime: easygo.NewInt64(util.GetMilliTime()),
					DeviceIdfa: easygo.NewString(idfa),
					Source:     easygo.NewString(source),
				}
				//保存idfa
				for_game.SavePosDeviceIdfa(msg)
				//easygo.Spawn(for_game.UpdatePosDeviceAdvIdfa, idfa, "IsActive", true, code) //广告设备激活
			}
		}
	} else {
		OutputJson(w, 0, "code不能为空", nil)
	}

}

//保存设备idfa码 参数idfa
func (self *WebHttpServer) SaveIdfaEntry(w http.ResponseWriter, r *http.Request) {
	logs.Error("11111111")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("content-type", "application/json")
	//r.ParseForm()
	//params := r.Form
	//key := params.Get("idfa")
	//logs.Info("设备上传idfa:", key)
	//if key != "" {
	//	idfa := url.QueryEscape(key)
	//	matched, err := regexp.MatchString("^[A-Z0-9-]+$", idfa)
	//	cal := params.Get("source")
	//	source := ""
	//	if cal != "" {
	//		source = url.QueryEscape(cal)
	//	}
	//
	//	if matched == true && err == nil {
	//		msg := &share_message.PosDeviceIdfa{
	//			CreateTime: easygo.NewInt64(util.GetMilliTime()),
	//			DeviceIdfa: easygo.NewString(idfa),
	//			Source:     easygo.NewString(source),
	//		}
	//		for_game.SavePosDeviceIdfa(msg)
	//		easygo.Spawn(for_game.UpdatePosDeviceAdvIdfa, idfa, "IsActive", true) //广告设备激活
	//	}
	//
	//	OutputJson(w, 1, "成功", idfa)
	//} else {
	//	OutputJson(w, 0, "idfa不能为空", nil)
	//}
}

//检查设备idfa码是否存在 参数idfa
func (self *WebHttpServer) CheckIdfaEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json")
	r.ParseForm()
	params := r.Form
	key := params.Get("idfa")

	if key != "" {
		idfa := url.QueryEscape(key)
		matched, err := regexp.MatchString("^[A-Z0-9-]+$", idfa)
		result := "1"
		if matched == true && err == nil {
			if for_game.CheckPosDeviceIdfa(idfa) {
				result = "0"
			}
		}

		ma := make(MAP_STRING)
		ma[idfa] = result
		out, _ := json.Marshal(ma)
		w.Write(out)
	} else {
		OutputJson(w, 0, "idfa不能为空", nil)
	}
}

//查询群资料 参数id
func (self *WebHttpServer) TeamdataEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json")
	r.ParseForm()
	params := r.Form
	key := params.Get("id")
	if key != "" {
		id := url.QueryEscape(key)
		team := QueryTeambyId(easygo.AtoInt64(id))
		if team != nil {
			teamdate := share_message.TeamData{
				Id:       team.Id,
				TeamChat: team.TeamChat,
				Name:     team.Name,
				GongGao:  team.GongGao,
			}
			OutputJson(w, 1, "成功", teamdate)
		}
	} else {
		OutputJson(w, 0, "id错误", nil)
	}
}

//发验证码 参数phone
func (self *WebHttpServer) SendPhoneCode(w http.ResponseWriter, params url.Values) {
	phone := params.Get("phone")
	// 参数校验
	m := make(map[string]string)
	m["phone"] = phone
	if errStr := VerifyParams(m, "SendPhoneCode"); errStr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}

	data := for_game.MessageMarkInfo.GetMessageMarkInfo(for_game.CLIENT_CODE_REGISTER, phone)
	if data != nil {
		leaveTime := time.Now().Unix() - data.Timestamp
		if leaveTime <= 120 {
			OutputJson(w, 0, "验证码已发送", nil)
			return
		}
	}

	codes := for_game.SendCodeToClientUser(for_game.CLIENT_CODE_REGISTER, phone, "")
	if codes == "" {
		OutputJson(w, 0, "验证码发送失败", nil)
		return
	}

	OutputJson(w, 1, "success", nil)
}

//注册账号 参数phone,code,channelno
func (self *WebHttpServer) Registered(w http.ResponseWriter, params url.Values) {
	logs.Debug(params)
	phone := params.Get("phone")
	code := params.Get("code")
	channelno := params.Get("channelno")
	// 参数校验
	m := make(map[string]string)
	m["phone"] = phone
	m["code"] = code
	m["channelno"] = channelno
	if errStr := VerifyParams(m, "Registered"); errStr != "" {
		OutputJson(w, 0, errStr, nil)
		return
	}

	err := for_game.CheckMessageCode(phone, code, for_game.CLIENT_CODE_REGISTER)
	if err != nil {
		OutputJson(w, 0, "验证码错误", nil)
		return
	}

	op_cl := for_game.QueryOperationByNo(channelno)
	if op_cl == nil {
		OutputJson(w, 0, "渠道号错误", nil)
		return
	}

	acc := for_game.GetRedisAccountByPhone(phone)
	if acc != nil {
		OutputJson(w, 0, "账号重复！请修改账号重新确认添加", nil)
		return
	}

	ip := "127.0.0.1"
	data := &share_message.CreateAccountData{
		Phone:     easygo.NewString(phone),
		PassWord:  easygo.NewString("0000"),
		IsVisitor: easygo.NewBool(false),
		Ip:        easygo.NewString(ip),
		IsOnline:  easygo.NewBool(false),
		Types:     easygo.NewInt32(1),
	}

	b, playerId := for_game.CreateAccount(data)
	if !b {
		OutputJson(w, 0, "账号创建失败,请稍后再试", nil)
		return
	}

	player := for_game.GetRedisPlayerBase(playerId)
	player.SetNickName(player.GetAccount())
	player.SetDeviceType(3)
	player.SetCreateTime()
	player.SetChannel(channelno)
	player.SaveToMongo()
	OutputJson(w, 1, "success", nil)
	easygo.Spawn(for_game.MakeOperationChannelReport, 1, playerId, channelno, nil, nil) //生成运营渠道数据汇总报表 已优化到Redis
}

//注册api
func (self *WebHttpServer) RegisteredEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	params := r.Form
	if !self.ApiCheckSign(params) {
		OutputJson(w, 0, "签名错误", nil)
		return
	}
	logs.Info(params)
	t := params.Get("t") //请求数据类型
	switch t {
	case "1": //获取验证码
		self.SendPhoneCode(w, params)
	case "2": //注册
		self.Registered(w, params)
	default:
		OutputJson(w, 0, "请求类型错误", nil)
	}
}

//用户冻结日志 参数 pagesize,curpage
func (self *WebHttpServer) FreezeLogEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json")
	r.ParseForm()
	params := r.Form
	pageS := params.Get("pagesize")
	curS := params.Get("curpage")
	page := 10
	cur := 1
	if pageS != "" {
		page = easygo.StringToIntnoErr(pageS)

	}

	if curS != "" {
		cur = easygo.StringToIntnoErr(curS)
	}
	list, count := GetPlayerFreezeLogsList(int32(cur), int32(page))

	msg := &brower_backstage.PlayerFreezeLogResponse{
		List:      list,
		PageCount: easygo.NewInt32(count),
	}
	OutputJson(w, 1, "成功", msg)
}

//查询社交广场动态信息 参数 id
func (self *WebHttpServer) SquareEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json")
	r.ParseForm()
	params := r.Form
	id := params.Get("id")
	if id != "" {
		i := easygo.StringToIntnoErr(id)
		if i < 1 {
			OutputJson(w, 0, "id错误", nil)
		}
		logid := int64(i)
		dynamicInfo := for_game.GetRedisDynamic(logid)
		if dynamicInfo != nil {
			player := for_game.GetRedisPlayerBase(dynamicInfo.GetPlayerId())
			msg := for_game.GetRedisDynamicAllInfo(0, logid, dynamicInfo.GetPlayerId(), player.GetAttention())
			result := &for_game.DynamicData{}
			for_game.StructToOtherStruct(dynamicInfo, result)
			// commentList:=GetRedisDynamicCommentInfo(logid, 0)
			commentList := msg.GetCommentList()
			for i, j := range commentList.CommentInfo {
				commPlayer := for_game.GetRedisPlayerBase(j.GetPlayerId())
				commentList.CommentInfo[i].Sex = easygo.NewInt32(commPlayer.GetSex())
				commentList.CommentInfo[i].HeadIcon = easygo.NewString(commPlayer.GetHeadIcon())
				commentList.CommentInfo[i].Name = easygo.NewString(commPlayer.GetNickName())
			}
			result.CommentNum = msg.GetCommentNum()
			result.CommentList = commentList
			result.Zan = msg.GetZan()
			result.HeadIcon = player.GetHeadIcon()
			result.Sex = player.GetSex()
			result.Account = player.GetAccount()
			result.NickName = player.GetNickName()
			OutputJson(w, 1, "成功", result)
		} else {
			OutputJson(w, 0, "动态不存在", nil)
		}
	} else {
		OutputJson(w, 0, "参数错误", nil)
	}
}

//召回用户api
func (self *WebHttpServer) UpRecallEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	params := r.Form
	if !self.ApiCheckSign(params) {
		OutputJson(w, 0, "签名错误", nil)
		return
	}

	t := params.Get("t") //请求数据类型
	switch t {
	case "1": //打开上报
		id := easygo.Get0ClockTimestamp(easygo.NowTimestamp())
		for_game.SetRedisRecallReportFildVal(id, 1, "Pv")
		OutputJson(w, 1, "success", nil)
	case "2": //下载上报
		id := easygo.Get0ClockTimestamp(easygo.NowTimestamp())
		for_game.SetRedisRecallReportFildVal(id, 1, "DownCount")
		OutputJson(w, 1, "success", nil)
	case "3": //每设备打开上报
		id := easygo.Get0ClockTimestamp(easygo.NowTimestamp())
		for_game.SetRedisRecallReportFildVal(id, 1, "Uv")
		OutputJson(w, 1, "success", nil)
	default:
		OutputJson(w, 0, "请求类型错误", nil)
	}
}

//上传文件
func (self *WebHttpServer) UploadEntry(w http.ResponseWriter, r *http.Request) {
	// 5,242,880 bytes == 5 MiB
	maxUploadSize := flag.Int64("upload_limit", 5242880, "max size of uploaded file (byte)")
	tokenFlag := flag.String("token", "", "specify the security token (it is automatically generated if empty)")
	corsEnabled := flag.Bool("cors", false, "if true, add ACAO header to support CORS")
	flag.Parse()

	token := *tokenFlag
	if token == "" {
		count := 10
		b := make([]byte, count)
		if _, err := rand.Read(b); err != nil {
			logs.Error("could not generate token")
		}
		token = fmt.Sprintf("%x", b)
		// logs.Info("token", token, "token generated")
	}

	var url string
	r.ParseForm()
	params := r.Form
	t := params.Get("type") //上传类型 1-上传到服务器
	if t == "1" {
		server := easygo.NewServer("upload", *maxUploadSize, token, *corsEnabled)
		url = server.HandlePost(w, r)
	} else {
		srcFile, info, err := r.FormFile("file")
		if err != nil {
			OutputJson(w, 0, "failed to acquire the uploaded content", nil)
			return
		}
		defer srcFile.Close()

		body, err := ioutil.ReadAll(srcFile)
		if err != nil {
			OutputJson(w, 0, "failed to read the uploaded content", nil)
			return
		}
		filename := info.Filename
		if filename == "" {
			filename = fmt.Sprintf("%x", sha1.Sum(body))
		}

		pathfileName := path.Join("backstage", "upload", filename)
		url = QQbucket.ObjectPutByte(pathfileName, body)
		if url == "" {
			OutputJson(w, 0, "上传存储桶失败", nil)
			return
		}
	}

	OutputJson(w, 1, "成功", url)
}

type YouMiResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

//有米api请求
func (self *WebHttpServer) YourMiEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	var code int = 0
	var msg string = "success"
	defer func() {
		resp := &YouMiResp{
			Code: code,
			Msg:  msg,
		}
		b, err := json.Marshal(resp)
		if err != nil {
			s := "{\"code\":1000,\"msg\":\"系统异常\"}"
			w.Write([]byte(s))
		}
		w.Write(b)
	}()
	params := r.Form
	idfa := params.Get("idfa")
	imei := params.Get("imei")
	source := params.Get("source")
	//appid := params.Get("appId")
	callBackUrl := params.Get("callback_url")
	cb, _ := url.PathUnescape(callBackUrl)
	if idfa == "" && imei == "" {
		code = 1001
		msg = "设备ID不能为空"
		return
	}
	if callBackUrl == "" {
		code = 1002
		msg = "回调地址不能为空"
		return
	}
	if source != "youmi" {
		code = 1003
		msg = "渠道号不正确"
		return
	}
	//if appid != "iosPack" && appid != "androidPack" {
	//	code = 1004
	//	msg = "渠道包不正确"
	//	return
	//}
	cd := idfa
	oS := 0
	if cd == "" {
		cd = imei
		oS = 1
	}
	logs.Info("params:", params)
	ip := r.RemoteAddr
	record := &share_message.KsPosAdvIdfa{
		CodeMd5:    easygo.NewString(for_game.Md5(cd)),
		CreateTime: easygo.NewInt64(for_game.GetMillSecond()),
		OsType:     easygo.NewInt32(oS),
		Ip:         easygo.NewString(ip),
		Callback:   easygo.NewString(cb),
		IsActive:   easygo.NewBool(false),
		IsRegister: easygo.NewBool(false),
		Platform:   easygo.NewString(source),
	}
	b := for_game.SavePosDeviceAdvIdfa(record, false)
	if !b {
		code = 1005
		msg = "已有点击记录"
		return
	}
}

//分享落地页
func (self *WebHttpServer) ShareTopicEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	params := r.Form

	t := params.Get("topicId") //请求数据类型
	topicId := easygo.StringToInt64noErr(t)

	// 判断该话题是否存在
	topic := for_game.GetTopicByIdNoStatusFromDB(topicId)
	if topic == nil {
		logs.Error("话题主页动态列表 RpcGetTopicMainPageList,该话题不存在, topicId: ", topicId)
		return
	}
	if topic.GetStatus() == for_game.TOPIC_STATUS_CLOSE {
		logs.Error("违规话题整改中,状态为已关闭,id为: ", topicId)
		return
	}
	page := int64(1)
	pageSize := int64(15)
	dynamicList, _ := for_game.GetNewDynamicByTopicId(topicId, page, pageSize)
	playerIds := make([]int64, 0)
	for _, v := range dynamicList {
		playerIds = append(playerIds, v.GetPlayerId())
	}
	players := for_game.GetPlayerListByIds(playerIds)
	playerMap := make(map[int64]*share_message.PlayerBase, len(players))
	for _, v := range players {
		playerMap[v.GetPlayerId()] = v
	}

	m := make([]map[string]interface{}, 0)

	for _, v := range dynamicList {
		dynamic := make(map[string]interface{})
		dynamic["Content"] = v.Content
		dynamic["Photo"] = v.Photo
		dynamic["Zan"] = for_game.GetRedisDynamicZanNum(v.GetLogId())
		dynamic["Voice"] = v.Voice
		dynamic["Video"] = v.Video
		dynamic["CommentNum"] = for_game.GetRedisDynamicCommentNum(v.GetLogId())
		dynamic["VoiceTime"] = v.VoiceTime
		dynamic["High"] = v.High
		dynamic["Weight"] = v.Weight
		dynamic["VideoThumbnailURL"] = v.VideoThumbnailURL
		dynamic["SendTime"] = v.SendTime
		if player, ok := playerMap[v.GetPlayerId()]; ok {
			dynamic["NickName"] = player.NickName
			dynamic["Sex"] = player.Sex
			dynamic["HeadIcon"] = player.HeadIcon
		}
		m = append(m, dynamic)
	}
	OutputJson(w, 0, "dynamicList", m)
}
