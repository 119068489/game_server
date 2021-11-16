package for_game

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"game_server/pb/share_message"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"

	"github.com/akqp2019/protobuf/proto"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"

	"game_server/easygo"
	"game_server/easygo/base"
	"time"
)

type IWebHttpServer interface {
	InitPath()
	SendToServer(server *share_message.ServerInfo, methodName string, msg easygo.IMessage) (bool, []byte)
}

type WebHttpServer struct {
	Port       int32       //监听的端口
	R          *gin.Engine //gin对象
	Services   map[string]reflect.Value
	ServiceMap map[string]easygo.Pair
	Me         IWebHttpServer
}

func Middle(c *gin.Context) {
	ip := c.ClientIP()
	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		_, _ = c.Writer.Write(PacketProtoMsg(1, easygo.NewFailMsg("err Middle 0")))
		c.Abort()
	}
	packet := &base.Packet{}
	err = packet.Unmarshal(bytes)
	if err != nil {
		logs.Error("err:", err.Error())
		logs.Debug("收到的逻辑包有误,反序列化出错,收到字节流是%v", string(bytes))
		_, _ = c.Writer.Write(PacketProtoMsg(1, easygo.NewFailMsg("err Middle 1")))
		c.Abort()

	}
	if packet.GetType() == base.PacketType_TYPE_REQUEST {
		request := &base.Request{}
		err := request.Unmarshal(packet.GetSerialized())
		if err != nil {
			logs.Warn("请求的逻辑包有误,反序列化出错 %v", packet.GetSerialized())
			_, _ = c.Writer.Write(PacketProtoMsg(1, easygo.NewFailMsg("err Middle 2")))
			c.Abort()
		}
		//if request.GetMethodName() != "RpcLogin" {
		com, b := c.Get("Common")
		common, ok := com.(*base.Common)
		if !b || !ok {
			if request.GetCommon() != nil {
				common = request.GetCommon()
				common.Ip = easygo.NewString(ip)
				c.Set("Common", common)
			} else {
				_, _ = c.Writer.Write(PacketProtoMsg(1, easygo.NewFailMsg("err Middle 3")))
				c.Abort()
			}
		}
		if common.GetToken() == "" && common.GetFlag() != MSG_FLAG { //服务器间通讯不需要token
			logs.Info("token 为空")
			_, _ = c.Writer.Write(PacketProtoMsg(1, easygo.NewFailMsg("err Middle 5")))
			c.Abort()
		}
		//}
		c.Set("Data", request)

	} else {
		_, _ = c.Writer.Write(PacketProtoMsg(1, easygo.NewFailMsg("err Middle 6")))
		c.Abort()
	}
}
func CommonMiddle(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	s := c.Request.Header.Get("Common")
	byte, err := base64.StdEncoding.DecodeString(s)
	easygo.PanicError(err)
	common := &base.Common{}

	if len(byte) == 0 {
		return
	}
	err = common.Unmarshal(byte)
	if err != nil {
		_, _ = c.Writer.Write(PacketProtoMsg(1, easygo.NewFailMsg("err CommonMiddle 1")))
		c.Abort()
	}
	ip := c.ClientIP()
	common.Ip = easygo.NewString(ip)
	c.Set("Common", common)
}

func (self *WebHttpServer) Init(port int32, services map[string]interface{}, serviceMap map[string]easygo.Pair) {
	self.Port = port
	self.R = gin.Default()
	self.R.Use(CommonMiddle)
	self.R.Use(Middle)
	self.Services = make(map[string]reflect.Value)
	for modName, service := range services {
		v := reflect.ValueOf(service)
		if v.Kind() == reflect.Ptr {
			if v.Elem().Kind() != reflect.Struct {
				panic("参数 service 必须是个 struct 对象或 struct 指针")
			}
		} else if v.Kind() != reflect.Struct {
			panic("参数 service 必须是个 struct 对象或 struct 指针")
		}

		self.Services[modName] = v
	}

	self.ServiceMap = serviceMap
}
func (self *WebHttpServer) GetGinRoute() *gin.Engine {
	return self.R
}

//初始化路由:子类实现
func (self *WebHttpServer) InitPath() { //override
	panic("该方法由子类实现")
}
func (self *WebHttpServer) Serve() {
	address := net.JoinHostPort("0.0.0.0", easygo.AnytoA(self.Port))
	logs.Info("address:", address)
	err := self.R.Run(address)
	easygo.PanicError(err)
}

//打包发送消息
func PacketProtoMsg(requestId uint64, msg easygo.IMessage) []byte {
	isSuccess := true
	if _, ok := msg.(*base.Fail); ok {
		isSuccess = false
	}
	serialized, err := msg.Marshal()
	easygo.PanicError(err)
	st := base.ResponseType_TYPE_SUCCESS
	if !isSuccess {
		st = base.ResponseType_TYPE_FAIL
	}
	response := base.Response{ResponseId: &requestId, SubType: &st, Serialized: serialized, MsgName: easygo.NewString(proto.MessageName(msg))}
	bs := easygo.Marshal(&response)
	t := base.PacketType_TYPE_RESPONSE
	packet := base.Packet{Type: &t, Serialized: bs}
	bs = easygo.Marshal(&packet)
	return bs
}

//接收到的都是请求包
func (self *WebHttpServer) RecvPacket(bytes []byte) (bool, easygo.IMessage, int64) {
	packet := &base.Packet{}
	err := packet.Unmarshal(bytes)
	if err != nil {
		logs.Debug("收到的逻辑包有误,反序列化出错,收到字节流是%v", string(bytes))
		return false, nil, 0
	}
	if packet.GetType() == base.PacketType_TYPE_REQUEST {
		request := &base.Request{}
		err := request.Unmarshal(packet.GetSerialized())
		if err != nil {
			logs.Warn("请求的逻辑包有误,反序列化出错 %v", packet.GetSerialized())
			return false, nil, 0
		}
		return true, request, time.Now().Unix()

	} else {
		logs.Debug("对端恶意:收到未知类型的网络包:%v", packet)
		return false, nil, 0
	}
}
func (self *WebHttpServer) DealRequest(recvTimeStamp int64, request *base.Request, common *base.Common) easygo.IMessage {
	methodName, _ := request.GetMethodName(), request.GetRequestId()
	pair, ok := self.ServiceMap[methodName]
	if !ok {
		logs.Warn(`收到 "%s" 请求,但是在 serviceMap 中找不到这个 key.`, methodName)
		return easygo.NewFailMsg("err:MethodName not find")
	}
	requestCls, _ := pair[0], pair[1]

	reqMsg := easygo.NewMessage(requestCls)
	if request.Serialized != nil { // 当消息体没有域时不走这里，比如 base.Empty
		err := reqMsg.Unmarshal(request.Serialized)
		if err != nil {
			s := fmt.Sprintf("rpc调用%v时,反序列化 msg 失败,很可能消息定义不一致.err=%v", methodName, err)
			logs.Debug(s)
			easygo.NewFailMsg("err:" + s)
		}
	}

	var method reflect.Value
	for _, service := range self.Services {
		method = service.MethodByName(methodName)
		if method.IsValid() {
			break
		}
	}
	if !method.IsValid() {
		s := fmt.Sprintf("找不到 %v 方法", methodName)
		logs.Debug(s)
		return easygo.NewFailMsg(s)
	}
	if method.Kind() != reflect.Func {
		s := fmt.Sprintf("%v 不是方法", methodName)
		logs.Debug(s)
		return easygo.NewFailMsg(s)
	}
	args := []reflect.Value{
		reflect.ValueOf(common),
		reflect.ValueOf(reqMsg),
	}
	results := method.Call(args) // 分发到各 rpc 业务处理函数
	//// TODO 记录时间消耗
	if len(results) == 0 {
		s := fmt.Sprintf("%v 方法函数原型错了,返回值类型必须是 IMessage", methodName)
		panic(s)
	}
	result := results[0].Interface()
	if result == nil {
		return &base.NoReturn{}
	}
	var respMsg easygo.IMessage
	respMsg = result.(easygo.IMessage)
	return respMsg
}
func (self *WebHttpServer) SendToServer(server *share_message.ServerInfo, methodName string, msg easygo.IMessage, pid ...int64) (easygo.IMessage, *base.Fail) {
	return SendToServerEx(server, methodName, msg, pid...)
}

//给指定服务器发送消息
func SendToServerEx(server *share_message.ServerInfo, methodName string, msg easygo.IMessage, pid ...int64) (easygo.IMessage, *base.Fail) {
	msg1, err := msg.Marshal()
	easygo.PanicError(err)
	request := base.Request{
		MethodName: easygo.NewString(methodName),
		Serialized: msg1,
		Timestamp:  easygo.NewInt64(time.Now().Unix()),
	}
	msg2, err := request.Marshal()
	t := base.PacketType_TYPE_REQUEST
	packet := base.Packet{
		Type:       &t,
		Serialized: msg2,
	}
	port := server.GetServerApiPort()
	if port == 0 {
		logs.Error("服务器api端口怎么会为0呢")
		return nil, easygo.NewFailMsg("服务器api端口怎么会为0呢")
	}
	u := "http://" + server.GetInternalIP() + ":" + easygo.AnytoA(port) + "/api"
	data, err := packet.Marshal()
	userId := append(pid, 0)[0]
	common := &base.Common{
		Version: easygo.NewString(server.GetVersion()),
		UserId:  easygo.NewInt64(userId),
		Token:   easygo.NewString(""),
		Flag:    easygo.NewInt32(MSG_FLAG),
	}
	bs, err := DoBytesPost(u, data, common)
	if err != nil {
		logs.Error("err:", err)
		return nil, easygo.NewFailMsg(err.Error())
	}
	b := &base.Packet{}
	err = b.Unmarshal(bs)
	if err != nil {
		logs.Error("err:", err)
		return nil, easygo.NewFailMsg(err.Error())
	}
	resp := &base.Response{}
	err = resp.Unmarshal(b.GetSerialized())
	if err != nil {
		logs.Error("err:", err)
		return nil, easygo.NewFailMsg("resp Unmarshal err")
	}
	msgName := resp.GetMsgName()
	if msgName == "" {
		return nil, nil
	}
	rspMsg := easygo.NewMessage(msgName)
	err = rspMsg.Unmarshal(resp.GetSerialized())
	if err != nil {
		return nil, easygo.NewFailMsg(err.Error())
	}
	if resp.GetSubType() == base.ResponseType_TYPE_SUCCESS {
		return rspMsg, nil
	} else {
		return nil, rspMsg.(*base.Fail)
	}
}

//body提交二进制数据
func DoBytesPost(url string, data []byte, common *base.Common) ([]byte, error) {
	body := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		logs.Error("err:", err)
		return nil, err
	}
	request.Header.Set("Connection", "Keep-Alive")
	com, err := common.Marshal()
	if err != nil {
		logs.Error("err:", err)
		return nil, err
	}
	request.Header.Set("Common", base64.StdEncoding.EncodeToString(com))
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		logs.Error("err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("err:", err)
		return nil, err
	}
	return b, err
}
