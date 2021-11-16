package for_game

import (
	"encoding/json"
	"game_server/easygo"
	"io/ioutil"
	"net/http"
)

/*
	获取B站数据接口
*/

type BstationContent struct {
	Replies []struct {
		Content struct {
			Message string `json:"message"`
		} `json:"content"`
	} `json:"replies"`
}

func GetBstationContent(url string) *BstationContent {
	res, err := http.Get(url)
	easygo.PanicError(err)
	resData, err1 := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	easygo.PanicError(err1)

	params := make(easygo.KWAT)
	err2 := json.Unmarshal(resData, &params)
	easygo.PanicError(err2)
	if params["data"] == nil {
		return nil
	}
	jsData := easygo.AnytoA(params["data"])
	result := &BstationContent{}
	err2 = json.Unmarshal([]byte(jsData), &result)
	easygo.PanicError(err2)

	return result
}
