package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//lol交战统计
type BattleBody struct {
	QueryType    int64                `json:"query_type"`
	TeamId       []int64              `json:"team_id"`
	StatsType    int64                `json:"stats_type"`
	StatsCount   int64                `json:"stats_count"`
	StatsDate    int64                `json:"stats_date"`
	SpecialStats []BattleSpecialStats `json:"special_stats"`
}

//lol联赛统计
type LeagueBody struct {
	LeagueId     int64                `json:"league_id"`
	SpecialStats []BattleSpecialStats `json:"special_stats"`
}

//lol近期统计
type RecentBody struct {
	QueryType    int64                `json:"query_type"`
	SeriesId     int64                `json:"series_id"`
	StatsType    int64                `json:"stats_type"`
	StatsCount   int64                `json:"stats_count"`
	StatsDate    int64                `json:"stats_date"`
	SpecialStats []BattleSpecialStats `json:"special_stats"`
}

type RecentBody2 struct {
	QueryType    int64                `json:"query_type"`
	TeamId       int64                `json:"team_id"`
	StatsType    int64                `json:"stats_type"`
	StatsCount   int64                `json:"stats_count"`
	StatsDate    int64                `json:"stats_date"`
	SpecialStats []BattleSpecialStats `json:"special_stats"`
}

type BattleSpecialStats struct {
	Type  int64   `json:"type"`
	Value []int64 `json:"value"`
}

func main() {
	var timeStrParam string = strconv.FormatInt(time.Now().Unix(), 10)
	//////lol交战统计数据
	poxiaoURL := "https://openapi.dawnbyte.com/api/stats/battle/lol?game_id=1&time_stamp="
	poxiaoURL = poxiaoURL + timeStrParam
	//timeStr := time.Now().Unix()
	//
	body := RecentBody{
		QueryType:  1,
		SeriesId:   1606972685,
		StatsType:  1,
		StatsCount: 20,
	}

	////lol联赛统计
	//poxiaoURL := "https://openapi.dawnbyte.com/api/stats/league/lol?game_id=1&time_stamp="
	//poxiaoURL = poxiaoURL + timeStrParam
	//
	//body := LeagueBody{LeagueId: 5471}

	////lol近期统计
	//poxiaoURL := "https://openapi.dawnbyte.com/api/stats/recent/lol?game_id=1&time_stamp="
	//poxiaoURL = poxiaoURL + timeStrParam

	//body := RecentBody{
	//	QueryType:  1,
	//	SeriesId:   1606972685,
	//	StatsType:  1,
	//	StatsCount: 12,
	//}

	//body := RecentBody2{
	//	QueryType:  2,
	//	TeamId:     334,
	//	StatsType:  1,
	//	StatsCount: 6,
	//	StatsDate:  timeStr,
	//}

	tempSpecialStatses := make([]BattleSpecialStats, 0)
	for key, value := range for_game.LOLSpecialStatsCode {
		tempStats := BattleSpecialStats{
			Type:  key,
			Value: value,
		}
		tempSpecialStatses = append(tempSpecialStatses, tempStats)
	}

	body.SpecialStats = tempSpecialStatses
	data, _ := json.Marshal(body)
	rst, _ := doPoxiaoBytesPost(poxiaoURL, data)
	fmt.Println("=======rst=======", rst)

}

//body提交二进制数据
func doPoxiaoBytesPost(url string, data []byte) (string, error) {
	body := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url, body)
	easygo.PanicError(err)
	request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("token", "nKf635U5zDb8uzWRPcoB1p7lsIOAdkfcfkKcqxBzxQZrUoCGL1")

	var resp *http.Response

	resp, err = http.DefaultClient.Do(request)
	easygo.PanicError(err)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err)
	return string(b), err
}
