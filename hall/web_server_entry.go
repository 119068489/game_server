package hall

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/for_game"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

var (
	_H5_LOG_PATH         = "web_api.log"
	_H5_JOB_ID           = "web_jobId"
	_H5_PER_CONNECT_TIME = int64(30)          //定时一分钟30次封
	_H5_MAX_CONNECT_TIME = int64(500)         //一分钟500次立即封
	WEB_API_BLACK_IP     = "web_api_black_ip" //redis key
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

type WebHttpServer struct {
	Service  reflect.Value
	JobId    int64
	Address  map[string]int64 //5分钟内访问接口ip次数，每5分钟清一次
	LastTime int64            //上一次请求的时间
}

func NewWebHttpServer() *WebHttpServer {
	p := &WebHttpServer{}
	p.Init()
	return p
}
func (self *WebHttpServer) Init() {
	self.Address = make(map[string]int64)
}
func (self *WebHttpServer) NextJobId() {
	self.JobId = for_game.NextId(_H5_JOB_ID)
}

//更新上次访问时间
func (self *WebHttpServer) UpdateLastTime() {
	self.LastTime = easygo.NowTimestamp()
}
func (self *WebHttpServer) Serve() {
	//port := easygo.YamlCfg.GetValueAsInt("LISTEN_ADDR_FOR_WEB_API")
	address := for_game.MakeAddress("0.0.0.0", PServerInfo.GetWebApiPort())
	logs.Info("(API 服务) 开始监听: %v", address)
	http.HandleFunc("/miaodao", self.MiaoDaoEntry) //秒到支付回调
	http.HandleFunc("/wxpay", self.WXPayEntry)     // 支付接口 所有支付都是从这个入口进去.
	http.HandleFunc("/tonglian", self.TongLianEntry)
	http.HandleFunc("/wxpayresult", self.WXPayResultEntry) // 小程序支付结果
	http.HandleFunc("/pengju", self.PengJuEntry)           // 鹏聚支付异步回调
	http.HandleFunc("/huiju", self.HuiJuEntry)             // 汇聚支付异步回调
	http.HandleFunc("/huichaopay", self.HuiChaoEntry)      // 汇潮支付异步回调
	http.HandleFunc("/huichaoDF", self.HuiChaoDFEntry)     // 汇潮代付异步回调
	http.HandleFunc("/tongtongpay", self.TongTongPayEntry) // 统统付异步回调
	http.HandleFunc("/checkorder", self.CheckOrderEntry)   // 提供前端查询支付订单
	http.HandleFunc("/coins", self.CoinsEntry)             // 硬币H5接口
	http.HandleFunc("/bp", self.BuryingPointEntry)         // 埋点接口
	http.HandleFunc("/ytspay", self.YunTongShangEntry)     // 云通商回调接口

	//启动定时任务，检测ip
	easygo.AfterFunc(time.Second*60, self.Update)
	err := http.ListenAndServe(address, nil)
	easygo.PanicError(err)
}
func (self *WebHttpServer) SendErrorCode(w http.ResponseWriter, errCode int, errText string) {
	//输入记录
	result := &ErrorResult{
		Code:  errCode,
		Error: errText,
	}
	content, err := json.Marshal(result)
	easygo.PanicError(err)
	self.LogPay("web请求返回:" + string(content))
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

//返回json格式数据
func (self *WebHttpServer) OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := &Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	w.Write(b)
}

//支付打印
func (self *WebHttpServer) LogPay(data interface{}) {
	for_game.WriteFile(_H5_LOG_PATH, easygo.AnytoA(self.JobId)+":", data)
}

//统计ip地址访问次数
func (self *WebHttpServer) CheckAddress(addr string) bool {
	logs.Info("addres:", addr)
	params := strings.Split(addr, ":")
	if self.CheckBlackList(params[0]) {
		return false
	}
	if len(params) == 2 {
		self.Address[params[0]] += 1
		if self.Address[params[0]] >= _H5_MAX_CONNECT_TIME {
			self.AddBlackList(params[0])
		}
	}
	return true
}

//检测ip连接数，并处理拉黑
func (self *WebHttpServer) Update() {
	for ip, num := range self.Address {
		if num >= _H5_PER_CONNECT_TIME {
			self.AddBlackList(ip)
		}
	}
	self.Address = make(map[string]int64)
	easygo.AfterFunc(time.Second*60, self.Update)
}

//增加黑名单ip
func (self *WebHttpServer) AddBlackList(ip string) {
	err := easygo.RedisMgr.GetC().SAdd(WEB_API_BLACK_IP, ip)
	if err != nil {
		logs.Error("AddBlackList err", ip)
	}
}

//检测是否是黑名单
func (self *WebHttpServer) CheckBlackList(ip string) bool {
	return easygo.RedisMgr.GetC().SIsMember(WEB_API_BLACK_IP, ip)
}
