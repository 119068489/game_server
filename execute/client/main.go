/*
Websocket客户端模拟器
*/
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"game_server/e-sports/sport_crawl"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/easygo/util"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	"path"
	"regexp"

	// _ "game_server/pb/brower_backstage"
	"game_server/pb/client_hall"
	"game_server/pb/client_login"
	"game_server/pb/client_server"
	"game_server/pb/share_message"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"

	"golang.org/x/net/websocket"
)

type Command struct {
}

type Client struct {
	C      *websocket.Conn
	Addr   string
	IsConn bool
}

var Cmd = &Command{}
var CmdValue = reflect.ValueOf(Cmd)

var Conn *Client
var QQbucket = util.NewQQbucket() //腾讯云存储桶

//GET请求
func HttpGet(url string) []byte {
	res, err := http.Get(url)
	easygo.PanicError(err)
	result, err1 := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	easygo.PanicError(err1)
	return result
}

func NewWebsocketClient(host, path string) *Client {
	u := url.URL{Scheme: "ws", Host: host, Path: path}
	ws, err := websocket.Dial(u.String(), "", "http://"+host+"/")
	if err != nil {
		easygo.PanicError(err)
	}

	return &Client{
		C:      ws,
		Addr:   host,
		IsConn: true,
	}
}

func (c *Client) PrintFailMessage(methodName string, res []byte) {
	result := &base.Fail{}
	err := result.Unmarshal(res)
	if err != nil {
		easygo.PanicError(err)
	}

	if result != nil {
		logs.Info(methodName, result.GetReason())
	}
}

func (c *Client) GetSendData(methodName string, req easygo.IMessage) []byte {
	msg, err := req.Marshal()
	easygo.PanicError(err)
	request := base.Request{
		MethodName: easygo.NewString(methodName),
		Serialized: msg,
		Timestamp:  easygo.NewInt64(time.Now().Unix()),
	}
	msg2, err := request.Marshal()
	easygo.PanicError(err)
	t := base.PacketType_TYPE_REQUEST
	packet := base.Packet{
		Type:       &t,
		Serialized: msg2,
	}

	data, err := packet.Marshal()
	easygo.PanicError(err)

	return data
}

func (c *Client) SendMessage(methodName string, req easygo.IMessage) {
	body := c.GetSendData(methodName, req)
	err := websocket.Message.Send(c.C, body)
	if err != nil {
		easygo.PanicError(err)
	}
}

func (c *Client) ReceiveMessage(methodName string) []byte {
	var (
		data []byte
		b    base.Packet
		res  base.Response
		req  base.Request
	)
	err := websocket.Message.Receive(c.C, &data)
	if err != nil {
		easygo.PanicError(err)
	}

	err = b.Unmarshal(data)
	if err != nil {
		easygo.PanicError(err)
	}

	err = res.Unmarshal(b.GetSerialized())
	if err != nil {
		err = req.Unmarshal(b.GetSerialized())
		if err != nil {
			easygo.PanicError(err)
		}
		logs.Error("MethodName", req.GetMethodName())
		return req.GetSerialized()
	} else {
		return res.GetSerialized()
	}
}

func (c *Client) ISendMessage(methodName string, req, dist easygo.IMessage) {
	c.SendMessage(methodName, req)
	if dist != nil {
		res := c.ReceiveMessage(methodName)
		err := dist.Unmarshal(res)
		if err != nil {
			c.PrintFailMessage(methodName, res)
		} else {
			logs.Info(methodName, dist)
		}
	}
}

func (c *Client) OnMessage() {
	var msg = make([]byte, 512)
	if !c.IsConn {
		return
	}
	_, err := c.C.Read(msg)
	if err != nil {
		easygo.PanicError(err)
	}
	count := len(msg)
	methodName := string(msg[0])
	data := string(msg[:count])
	logs.Info("Receive: %s\n", methodName, data)

}

func CheckParam(p []string, n int) bool {
	if len(p) < n {
		logs.Error("参数不足")
		return false
	}
	return true
}

func main() {

	Conn = &Client{
		IsConn: false,
	}
	if len(os.Args) > 2 {
		address := os.Args[1] + ":" + os.Args[2]
		logs.Info("正在连接", address)
		Conn = NewWebsocketClient(address, "")
	}

	for {
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		s := input.Text()
		if s != "" {
			params := strings.Split(s, " ")
			methodName := params[0]
			if !strings.HasPrefix(methodName, "Rpc") {
				methodName = "Rpc" + methodName
			}
			params[0] = methodName
			DoCmd(params)
		}
	}
}

func DoCmd(params []string) {
	methodName := params[0]
	method := CmdValue.MethodByName(methodName)
	if !method.IsValid() || method.Kind() != reflect.Func {
		s := fmt.Sprintf("%v 不能识别的命令,方法没有实现", methodName)
		log.Println(s)
		return
	}
	args := make([]reflect.Value, 0, len(params)-1)
	for _, para := range params[1:] {
		v := reflect.ValueOf(para)
		args = append(args, v)
	}

	method.Call(args) // 分发
}

func InitConn() *Client {
	if Conn.IsConn && Conn.Addr != "127.0.0.1:4001" {
		return Conn
	}

	return NewWebsocketClient("127.0.0.1:4001", "")
}

func (c *Client) HeartBeat() {
	c.ISendMessage("RpcHeartbeat", &client_server.NTP{T1: easygo.NewInt64(time.Now().Unix())}, nil)
	easygo.AfterFunc(10*time.Second, c.HeartBeat)
}

func (c *Command) Rpclogin(param ...string) {
	Conn = InitConn()
	methodName := "RpcLogin"
	if len(param) == 0 {
		param = append(param, "admin2", "123456")
	}
	req := &brower_backstage.LoginRequest{
		UserAccount: easygo.NewString(param[0]),
		Password:    easygo.NewString(param[1]),
	}
	res := &brower_backstage.LoginResponse{}

	Conn.ISendMessage(methodName, req, res)
}

func (c *Command) Rpccoin(param ...string) {
	methodName := "RpcCoinItemList"
	req := &brower_backstage.ListRequest{
		CurPage:  easygo.NewInt32(0),
		PageSize: easygo.NewInt32(0),
	}
	res := &brower_backstage.CoinItemListResponse{}

	Conn.ISendMessage(methodName, req, res)
}

func (self *Command) RpcWish(param ...string) {
	methodName := "RpcBrowerTest"
	req := &base.Empty{}
	res := &base.Empty{}

	Conn.ISendMessage(methodName, req, res)
}

//测试获取回收说明配置
func (self *Command) RpcGetRecycleNoteCfg(param ...string) {
	methodName := "RpcGetRecycleNoteCfg"
	req := &base.Empty{}
	res := &brower_backstage.RecycleNoteCfg{}

	Conn.ISendMessage(methodName, req, res)
}

func (c *Command) Rpcgetplayer(param ...string) {
	methodName := "RpcGetPlayerByAccount"
	if b := CheckParam(param, 1); !b {
		return
	}
	req := &brower_backstage.QueryDataById{
		IdStr: easygo.NewString(param[0]),
	}
	res := &share_message.PlayerBase{}

	Conn.ISendMessage(methodName, req, res)
}

func (c *Command) Rpcexit(param ...string) {
	methodName := "RpcLogout"
	req := &base.Empty{}
	Conn.ISendMessage(methodName, req, nil)
}

func (c *Command) Rpcjs(param ...string) {
	for m := 1; m < 10; m++ {
		n := 1
	LOOP:
		if n <= m {
			fmt.Printf("%dx%d=%d ", n, m, m*n)
			n++
			goto LOOP
		} else {
			fmt.Println("")
		}
		n++
	}
}

func (c *Command) Rpcloginapp(param ...string) {
	if Conn.Addr != "127.0.0.1:1111" {
		Conn = NewWebsocketClient("127.0.0.1:1111", "")
	}

	methodName := "RpcLoginHall"
	req := &client_login.LoginMsg{
		Account:  easygo.NewString(param[0]),
		Password: easygo.NewString(param[1]),
		Type:     easygo.NewInt32(2),
	}
	res := &client_login.LoginResult{}
	Conn.ISendMessage(methodName, req, res)
	time.Sleep(1 * time.Second)

	var pa = []string{res.GetAccount(), res.GetToken(), easygo.AnytoA(res.GetPlayerId())}
	c.Rpcloginhall(pa...)

}

func (c *Command) Rpcloginhall(param ...string) {
	if Conn.Addr != "127.0.0.1:2001" {
		Conn = NewWebsocketClient("127.0.0.1:2001", "")
	}

	methodName := "RpcLogin"

	// res := &client_login.LoginResult{
	// 	Account:  easygo.NewString("13999"),
	// 	Token:    easygo.NewString("2PILQNgdq3PmmCUQ13999"),
	// 	PlayerId: easygo.NewInt64(1887550225),
	// }

	res := &client_login.LoginResult{
		Account:  easygo.NewString(param[0]),
		Token:    easygo.NewString(param[1]),
		PlayerId: easygo.NewInt64(easygo.AnytoA(param[2])),
	}

	reqh := &client_hall.LoginMsg{
		Account:        res.Account,
		Token:          res.Token,
		RegistrationId: easygo.NewString(""),
		Channel:        easygo.NewString(""),
		LoginType:      easygo.NewInt32(1),
		DeviceType:     easygo.NewInt32(2),
		Type:           easygo.NewInt32(2),
		PlayerId:       res.PlayerId,
		VersionNumber:  nil,
		Brand:          nil,
	}
	resh := &client_server.AllPlayerMsg{}

	Conn.ISendMessage(methodName, reqh, resh)
	Conn.HeartBeat()
	// logs.Info("登录成功")

}

func (c *Command) Rpcgetfn(param ...string) {
	// url := "https://www.fnscore.com/detail/information-643.html"
	// data := sport_crawl.GetFnCrawlData(url)
	// logs.Info(data)

	// url = "http://dj.sina.com.cn/article/iznezxt1898575.shtml"
	// data = sport_crawl.GetXlCrawlData(url)
	// logs.Info(data)

	// url = "http://www.tvbcp.com/detail/50201.html"
	// data = sport_crawl.GetSyCrawlData(url)
	// logs.Info(data)

	url := "http://www.wanplus.com/zh/video/1562436"
	data := sport_crawl.GetWanplusVideoData(url, "1562436")
	logs.Info("t", data.Title, "\ni", data.Img, "\nv", data.Video, "\nc", data.Class, "\ntime", data.Time)
}

func (c *Command) RpcqueryWishActPoolDetail(param ...string) {
	//Conn = InitConn()
	methodName := "RpcQueryWishActPoolDetail"
	req := &brower_backstage.ListRequest{
		CurPage:  easygo.NewInt32(1),
		PageSize: easygo.NewInt32(20),
	}
	res := &brower_backstage.WishActPoolDetail{}
	Conn.ISendMessage(methodName, req, res)
}

//守护者列表
func (c *Command) RpcGuardianList(param ...string) {
	// methodName := "RpcGuardianList"
	// req := &base.Empty{}
	// res := &brower_backstage.GuardianListResp{}
	// Conn.ISendMessage(methodName, req, res)
}

//删除白名单
func (c *Command) RpcDeleteWishAllowList() {
	methodName := "RpcDeleteWishAllowList"
	req := &brower_backstage.QueryDataByIds{
		Ids64: []int64{18828044, 18828045},
	}
	res := &base.Empty{}
	Conn.ISendMessage(methodName, req, res)
}

//添加白名单
func (c *Command) RpcAddWishAllowList(param ...string) {
	methodName := "RpcAddWishAllowList"
	req := &brower_backstage.AddWishAllowListReq{
		Accounts: []string{"lm7080089f", "lm708008a0", "lm708008a1"},
		Remark:   easygo.NewString("白名单备注test"),
	}
	res := &brower_backstage.WishAllowListResp{}
	Conn.ISendMessage(methodName, req, res)
}

//白名单列表
func (c *Command) RpcwishAllowList(param ...string) {
	methodName := "RpcWishAllowList"
	req := &base.Empty{}
	res := &brower_backstage.WishAllowListResp{}
	Conn.ISendMessage(methodName, req, res)
}

// 许愿池活动用户记录-活动用户记录列表
func (c *Command) RpcQueryWishActPlayerRecordList(param ...string) {
	methodName := "RpcQueryWishActPlayerRecordList"
	req := &brower_backstage.ListRequest{
		CurPage:  easygo.NewInt32(1),
		PageSize: easygo.NewInt32(20),
	}
	res := &brower_backstage.WishActPlayerRecordList{}
	Conn.ISendMessage(methodName, req, res)
}

func (c *Command) Rpcuploadobjapi() {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	filename := "photo.jpg"
	targetUrl := "http://192.168.150.8:4501/upload"
	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		logs.Error("error writing to buffer")

	}

	//打开文件句柄操作
	fh, err := os.Open("download/" + filename)
	if err != nil {
		logs.Error("error opening file")
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		logs.Error(err)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		logs.Error(err)
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
	}

	logs.Info("retuen", resp.Status, string(resp_body))
}

func (c *Command) Rpcuploadobj() {
	//读原图片
	// ff, _ := os.Open("download/photo.jpg")
	// defer ff.Close()
	// sourcebuffer := make([]byte, 500000)
	// n, _ := ff.Read(sourcebuffer)
	// file := base64.StdEncoding.EncodeToString(sourcebuffer[:n])

	url := "https://tpc.googlesyndication.com/simgad/11036827312291413221?sqp=4sqPyQQ7QjkqNxABHQAAtEIgASgBMAk4A0DwkwlYAWBfcAKAAQGIAQGdAQAAgD-oAQGwAYCt4gS4AV_FAS2ynT4&rs=AOga4qk6Q5Of02G7R4nv_7LZiT5Z4WFifw"
	methodName := "RpcUploadFile"
	req := &brower_backstage.UploadRequest{
		FileName: easygo.NewString("wlt.jpg"),
		Path:     easygo.NewString("upload"),
		// File:     sourcebuffer[:n],
		IsBucket: easygo.NewBool(true),
		Type:     easygo.NewInt32(1),
		FileUrl:  easygo.NewString(url),
	}
	res := &brower_backstage.UploadResponse{}

	Conn.ISendMessage(methodName, req, res)

}

func (c *Command) Rpcdelobj(param ...string) {
	// QQbucket.ObjectDel("backstage/upload/001.jpg")
	// QQbucket.ObjectDel("backstage/upload/abc.jpg")
	// QQbucket.ObjectDel("backstage/upload/def.jpg")
	// QQbucket.ObjectDel("backstage/upload/hhc.jpg")
	// QQbucket.ObjectDel("backstage/upload/xqd.jpg")
	// QQbucket.ObjectDel("backstage/upload/photo.jpg")
	var key string
	var path string
	var isBucket bool = true
	if len(param) > 0 {
		key = param[0]
	}
	if len(param) > 1 {
		path = param[1]
	}
	if len(param) > 2 && param[2] == "false" {
		isBucket = false
	}

	methodName := "RpcDelUploadFile"
	req := &brower_backstage.UploadRequest{
		FileName: easygo.NewString(key),
		Path:     easygo.NewString(path),
		IsBucket: easygo.NewBool(isBucket),
	}
	res := &base.Empty{}

	Conn.ISendMessage(methodName, req, res)
}

func (c *Command) Rpcobjlist(param ...string) {
	pagesize := 10
	path := "upload"
	if len(param) > 0 {
		p := easygo.StringToIntnoErr(param[0])
		if p > 0 {
			pagesize = p
		}
		if len(param) > 1 {
			path = param[1]
		}
	}

	methodName := "RpcUploadFileList"
	req := &brower_backstage.ListRequest{
		Keyword:  easygo.NewString(""),
		SrtType:  easygo.NewString(path),
		PageSize: easygo.NewInt32(pagesize),
		CurPage:  easygo.NewInt32(0),
	}
	res := &brower_backstage.UploadListResponse{}

	Conn.ISendMessage(methodName, req, res)

}

func (c *Command) Rpcobjls(param ...string) {
	if len(param) < 2 {
		logs.Error("参数不足")
		return
	}
	if len(param) < 3 {
		param = append(param, "")
	}
	size := int32(easygo.StringToIntnoErr(param[2]))
	i := 0
goon:
	lis := QQbucket.GetObjectList(param[0], param[1], size)
	for _, li := range lis {
		one := &brower_backstage.UploadList{}
		for_game.StructToOtherStruct(li, one)
		logs.Info(one)
		param[1] = one.GetTitle()
		i++
	}
	if len(lis) == 10 {
		goto goon
	}
	logs.Debug("count", i)
}

func (c *Command) Rpccheckorder(param ...string) {
	methodName := "RpcCheckOrder"
	req := &brower_backstage.OptOrderRequest{
		Oid: easygo.NewString(param[0]),
	}
	res := &base.Empty{}

	Conn.ISendMessage(methodName, req, res)
}

func (c *Command) RpcApplyTopicMaster(param ...string) {
	methodName := "RpcApplyTopicMaster"
	req := &brower_backstage.QueryDataById{
		Id64: easygo.NewInt64(param[0]),
		Id32: easygo.NewInt32(param[1]),
	}
	res := &base.Empty{}

	Conn.ISendMessage(methodName, req, res)
}

func (c *Command) Rpccrawllist(param ...string) {
	methodName := "RpcCrawlJobList"
	req := &brower_backstage.ListRequest{
		CurPage:  easygo.NewInt32(1),
		PageSize: easygo.NewInt32(10),
	}
	res := &brower_backstage.CrawlJobResponse{}

	if len(param) > 0 {
		req.SrtType = easygo.NewString(param[0])

	}
	if len(param) > 1 {
		req.Type = easygo.NewInt32(easygo.StringToIntnoErr(param[1]))
	}

	Conn.ISendMessage(methodName, req, res)
}
func (c *Command) Rpcgetdata() {
	url := "https://img1.famulei.com/match_pre/11738.json"
	result := HttpGet(url)
	one := for_game.ResultMsg{}
	err := json.Unmarshal(result, &one)
	if err != nil {
		logs.Error("one")
	}

	dis := &share_message.RecentData{}
	for_game.StructToOtherStruct(one.Data, &dis)

	logs.Debug(one.Data)
}

func (c *Command) Rpcgetrr() {
	methodName := "RpcQueryRechargeEsAct"

	req := &base.Empty{}
	res := &share_message.Activity{}

	Conn.ISendMessage(methodName, req, res)
}
func (c *Command) Rpctest1(param ...string) {
	doc, err := goquery.NewDocument("https://www.chaofan.com/video/lol")
	if err != nil {
		easygo.PanicError(err)
	}

	doc.Find("body .content-box .hot-data-box .left-box").Each(func(i int, s *goquery.Selection) {
		s.Find(".video-list li").EachWithBreak(func(j int, d *goquery.Selection) bool {
			vurl, _ := d.Find(".link").Attr("href")
			cP, _ := url.Parse(vurl)
			if cP.Scheme == "" {
				cP.Scheme = "https"
			}
			logs.Info("url", cP.String())
			img, _ := d.Find(".link .img-view").Attr("style")
			strs := strings.SplitN(img, "?", 2)
			if len(strs) == 2 {
				strs = strings.SplitN(strs[0], "('", 2)
			}
			if len(strs) == 2 {
				img = strs[1]
			}
			logs.Info("img", img)
			filenameWithSuffix := path.Base(cP.String())
			fileSuffix := path.Ext(filenameWithSuffix)
			filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
			logs.Info(filenameOnly)
			if filenameOnly == "edksosofj9z6bwv" {
				return false
			}
			return true
		})
	})
}

func (c *Command) Rpctestqqlist(param ...string) {
	doc, err := goquery.NewDocument("http://news.yxrb.net/202105/12219315.html")
	if err != nil {
		easygo.PanicError(err)
	}

	data := &sport_crawl.CrawlData{}
	questionTitle := doc.Find(".article-title p").Text()
	questionContent, _ := doc.Find("article").Html()

	data.Title = questionTitle
	data.Content = questionContent
	data.Time = easygo.NowTimestamp()
	data.Type = 1
	data.Source = "游戏日报"

	if data.Title == "" {
		return
	}

	logs.Info(data)
}

func (c *Command) Rpctestqq(param ...string) {
	path := "https://page.om.qq.com/page/ORWXMv-cDq7NLVnTENMdfU8w0"
	doc := sport_crawl.GetHtmlDoc(path)
	if doc == nil {
		return
	}
	data := &sport_crawl.CrawlData{}
	questionTitle := doc.Find("#content .header .title").Text()
	questionContent, _ := doc.Find("#content .article").Html()

	data.Title = questionTitle
	data.Content = questionContent
	data.Time = easygo.NowTimestamp()
	data.Type = 1
	data.Source = "腾讯网"

	logs.Info(data)
}

func (c *Command) Rpctest2(param ...string) {
	// paths := "https://www.chaofan.com/news/lol"
	// doc := sport_crawl.GetHtmlDoc(paths)
	// if doc == nil {
	// 	return
	// }

	// doc.Find("body .content-box .list-box").Each(func(i int, s *goquery.Selection) {
	// 	s.Find(".clearfix-row").EachWithBreak(func(j int, d *goquery.Selection) bool {
	// 		vurl, _ := d.Find("a").Attr("href")
	// 		filenameWithSuffix := path.Base(vurl)
	// 		fileSuffix := path.Ext(filenameWithSuffix)
	// 		filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	// 		title := d.Find(".v-a-m").Text()
	// 		logs.Info(title, filenameOnly)
	// 		return true

	// 	})
	// })

	// paths := "https://www.chaofan.com/news/ek8sos9smk9wuoe.html"
	// doc := sport_crawl.GetHtmlDoc(paths)
	// if doc == nil {
	// 	return
	// }

	// data := &sport_crawl.CrawlData{}
	// questionTitle := doc.Find(".news-layout-center h1").Text()
	// questionContent, _ := doc.Find(".news-layout-center .description").Html()

	// data.Title = questionTitle
	// data.Content = questionContent
	// data.Time = easygo.NowTimestamp()
	// data.Type = 1
	// data.Source = "超凡电竞"

	// if data.Title == "" {
	// 	return
	// }
	// logs.Info(data)
	rep, _ := regexp.Compile("^" + "news_chaofan")
	if rep.MatchString("news_chaofan_lol") {
		logs.Info("匹配成功")
		return
	}

}

func (c *Command) Rpctest3(param ...string) {
	// paths := "https://dj.sina.com.cn/information"
	// doc := sport_crawl.GetHtmlDoc(paths)
	// if doc == nil {
	// 	return
	// }

	// doc.Find(".E_sports_news_list").Each(func(i int, s *goquery.Selection) {
	// 	s.Find(".main_list_onepic").EachWithBreak(func(j int, d *goquery.Selection) bool {
	// 		vurl, _ := d.Find("h3 a").Attr("href")
	// 		filenameWithSuffix := path.Base(vurl)
	// 		fileSuffix := path.Ext(filenameWithSuffix)
	// 		filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	// 		title := d.Find("h3 a").Text()
	// 		logs.Info(title, filenameOnly)
	// 		return true

	// 	})
	// })

	data := sport_crawl.GetSinaCrawlData("kmyaawc5498012")
	logs.Info(data)

}

func (c *Command) Rpctest(param ...string) {
	paths := "https://bbs.hupu.com/lol-postdate"
	doc := sport_crawl.GetHtmlDoc(paths)
	if doc == nil {
		return
	}

	t := doc.Find(".for-list li")
	for i := 0; i < t.Length(); i++ {
		vurl, _ := t.Eq(i).Find("a[class=truetit]").Attr("href")
		filenameWithSuffix := path.Base(vurl)
		fileSuffix := path.Ext(filenameWithSuffix)
		filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
		title := t.Eq(i).Find("a[class=truetit]").Text()
		time := t.Eq(i).Find("a:nth-child(3)").Text() + " " + t.Eq(i).Find(".endreply a").Text() + ":00"
		a := easygo.GetTimeStrToTimestamp(time)
		if a < 0 {
			continue
		}
		logs.Info("tt", time, "a", a, "i", filenameOnly, "t:", title)
	}

	// doc.Find(".for-list").Each(func(i int, s *goquery.Selection) {
	// 	s.Find("li").EachWithBreak(func(j int, d *goquery.Selection) bool {
	// 		vurl, _ := d.Find("a[class=truetit]").Attr("href")
	// 		filenameWithSuffix := path.Base(vurl)
	// 		fileSuffix := path.Ext(filenameWithSuffix)
	// 		filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	// 		title := d.Find("a[class=truetit]").Text()
	// 		time := d.Find("a:nth-child(3)").Text() + " " + d.Find(".endreply a").Text() + ":00"
	// 		a := easygo.GetTimeStrToTimestamp(time)
	// 		logs.Info("t:", title, "i", filenameOnly, "tt", time, "a", a)
	// 		return true
	// 	})
	// 	// tt, _ := s.Find("li").Html()
	// 	// logs.Info(tt)
	// })

	// paths := "https://www.chaofan.com/news/ek8sos9smk9wuoe.html"
	// doc := sport_crawl.GetHtmlDoc(paths)
	// if doc == nil {
	// 	return
	// }

	// data := &sport_crawl.CrawlData{}
	// questionTitle := doc.Find(".news-layout-center h1").Text()
	// questionContent, _ := doc.Find(".news-layout-center .description").Html()

	// data.Title = questionTitle
	// data.Content = questionContent
	// data.Time = easygo.NowTimestamp()
	// data.Type = 1
	// data.Source = "超凡电竞"

	// if data.Title == "" {
	// 	return
	// }
	// logs.Info(data)

}
