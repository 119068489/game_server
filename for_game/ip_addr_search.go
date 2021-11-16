package for_game

import (
	"encoding/json"
	"game_server/easygo"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

/*
	ip地址查询接口
*/

type IpSearchMgr struct {
	ApiUrl  string
	AppCode string
}

func NewIpSearchMgr() *IpSearchMgr {
	p := &IpSearchMgr{}
	p.Init()
	return p
}

//http://api01.aliyun.venuscn.com/ip?ip=113.66.216.195
func (self *IpSearchMgr) Init() {
	self.ApiUrl = "http://api01.aliyun.venuscn.com/ip"
	self.AppCode = "dafc6328f9fd4ca9a8f07952d9faa196"
}

//查询所在省
func (self *IpSearchMgr) GetRegion(ip string) string {
	data := self.IpSearch(ip)
	if data == nil {
		return ""
	}
	return data.Region
}

//查询ip地址返回详情
func (self *IpSearchMgr) IpSearch(ip string, isSave ...bool) *LocationData {
	isSave = append(isSave, true)
	if strings.ContainsAny(ip, ":") {
		ip = ip[0:strings.LastIndex(ip, ":")]
	}

	vals := url.Values{}
	vals.Add("ip", ip)
	resData := self.HttpsReq(vals)
	params := make(easygo.KWAT)
	err := json.Unmarshal(resData, &params)
	if err != nil {
		logs.Error("IpSearch err:", err)
		return nil
	}
	if params["data"] == nil {
		return nil
	}
	jsData := easygo.AnytoA(params["data"])
	result := &LocationData{}
	err = json.Unmarshal([]byte(jsData), &result)
	if err != nil {
		logs.Error("IpSearch err:", err)
		return nil
	}

	if result != nil && isSave[0] {
		easygo.Spawn(UpsertLocation, result) //更新本地区域库
	}

	return result
}

//外部查询方法
func IpSearch(ip string, isSave ...bool) *LocationData {
	isSave = append(isSave, true)
	mgr := NewIpSearchMgr()
	data := mgr.IpSearch(ip, isSave[0])
	return data
}

//发起Get请求
func (self *IpSearchMgr) HttpsReq(q url.Values) []byte {
	u, _ := url.Parse(self.ApiUrl)
	u.RawQuery = q.Encode()
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", u.String(), nil)
	reqest.Header.Add("Authorization", "APPCODE "+self.AppCode)
	easygo.PanicError(err)
	response, _ := client.Do(reqest)
	result, err1 := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	easygo.PanicError(err1)
	return result
}

// //查询所在省
// func (self *IpSearchMgr) GetRegion(ip string) (string, *base.Fail) {
// 	data := self.IpSearch(ip)
// 	params := make(easygo.KWAT)
// 	err := json.Unmarshal(data, &params)
// 	jsData := easygo.AnytoA(params["data"])
// 	mpData := make(easygo.KWAT)
// 	err = json.Unmarshal([]byte(jsData), &mpData)
// 	easygo.PanicError(err)
// 	return mpData.GetString("region"), easygo.NewFailMsg(params.GetString("msg"), params.GetString("ret"))
// }

// //查询ip地址返回详情
// func (self *IpSearchMgr) IpSearch(ip string) []byte {
// 	if strings.ContainsAny(ip, ":") {
// 		ip = ip[0:strings.LastIndex(ip, ":")]
// 	}

// 	vals := url.Values{}
// 	vals.Add("ip", ip)
// 	resData := self.HttpsReq(vals)
// 	return resData
// }

//是否是屏蔽用户
func IsShieldUser(ip string, isSave ...bool) bool {
	isSave = append(isSave, true)
	loc := IpSearch(ip, isSave[0])
	if loc != nil && (loc.CountryId == "KH" || loc.Region == "广东" || loc.Region == "福建") {
		return true
	}
	return false
}
