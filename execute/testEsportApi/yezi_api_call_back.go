package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type testBody struct {
	GameId     int32  `json:"game_id"`
	EventId    int32  `json:"event_id"`
	Type       string `json:"type"`
	Func       string `json:"func"`
	UpdateTime int64  `json:"update_time"`
}

func main() {
	yeziURL := "http://127.0.0.1:10711/notice"
	//初始化参数
	param := url.Values{}
	var timeStr string = strconv.FormatInt(time.Now().Unix(), 10)

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("func", "game_info")
	param.Set("game_id", "52776")
	param.Set("type", "update")
	param.Set("timestamp", timeStr)
	param.Set("event_id", "24")

	body := testBody{
		GameId:     52776,
		EventId:    24,
		Type:       "update",
		Func:       "game_info",
		UpdateTime: 1611661675,
	}

	data, _ := json.Marshal(body)
	rst, _ := doBytesPost(yeziURL, data, param)
	fmt.Println("rst")
	if rst == "success" {
		fmt.Println("=========success===========")
	}
}

//body提交二进制数据
func doBytesPost(url string, data []byte, param url.Values) (string, error) {
	body := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url, body)
	easygo.PanicError(err)
	request.Header.Set("Connection", "Keep-Alive")
	reqParam := param.Encode()
	signParam := for_game.Md5(reqParam)

	request.Header.Set("xxe-request", reqParam)
	request.Header.Set("xxe-sign", signParam)
	var resp *http.Response

	resp, err = http.DefaultClient.Do(request)
	easygo.PanicError(err)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err)
	return string(b), err
}
