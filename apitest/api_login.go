package apitest

import (
	"bytes"
	"encoding/base64"
	"game_server/easygo"
	"game_server/easygo/base"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/astaxie/beego/logs"
)

func GetMsgCode() {

}

func SendToServer(methodName string, msg easygo.IMessage, pid ...int64) (easygo.IMessage, *base.Fail) {
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

	u := "http://127.0.0.1:1001/api"
	data, err := packet.Marshal()
	userId := append(pid, 0)[0]
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(userId),
		Token:   easygo.NewString(""),
		//Flag:    easygo.NewInt32(MSG_FLAG),
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
