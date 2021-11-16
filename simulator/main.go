package main

import (
	"bufio"
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/client_server"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

var _ = fmt.Sprintf
var _ = log.Println

//==========================================================================
var token string
var playerId int64
var hallConnect *HallConnector
var loginConnect *LoginConnector

var instanceId int32

type Command struct {
}
type SERVER_ID = for_game.SERVER_ID
type ENDPOINT_ID = int32

var Cmd = &Command{}
var CmdValue = reflect.ValueOf(Cmd)

func main() {
	initializer := NewInitializer()
	dict := easygo.KWAT{
		"logName":  "simulator",
		"yamlPath": "config_login.yaml",
	}
	initializer.Execute(dict)
	//日志
	_ = logs.SetLogger(logs.AdapterConsole)
	logs.SetLevel(logs.LevelDebug)
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
	logs.Async()

	length := len(os.Args)
	if length < 2 {
		log.Println("Usage:<command> [<Host>] <Port>")
		return
	}
	var address string
	if length == 2 {
		port := os.Args[1]
		address = "localhost:" + port
	} else {
		host := os.Args[1]
		port := os.Args[2]
		address = host + ":" + port
	}
	loginConnect = NewLoginConnector(address)
	loginConnect.ConnectOnce()
	//心跳
	HeartBeat()
	for {
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		s := input.Text()
		if s != "" {
			DoCmd(s)
		}
	}
}

func DoCmd(s string) {
	defer easygo.RecoverAndLog()

	params := strings.Split(s, " ")
	methodName := params[0]
	if !strings.HasPrefix(methodName, "Rpc") {
		methodName = "Rpc" + methodName
	}
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
func HeartBeat() {
	if hallConnect != nil {
		ep := hallConnect.FetchEndpoint()
		if ep != nil {
			ep.RpcHeartbeat(&client_server.NTP{
				T1: easygo.NewInt64(time.Now().Unix()),
			})
		}
	}
	if loginConnect != nil {
		ep := loginConnect.FetchEndpoint()
		if ep != nil {
			ep.RpcHeartbeat(&client_server.NTP{
				T1: easygo.NewInt64(time.Now().Unix()),
			})
		}
	}
	easygo.AfterFunc(10*time.Second, HeartBeat)
}
