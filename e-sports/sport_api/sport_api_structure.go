package sport_api

import "game_server/pb/share_message"

//5.0野子科技API比赛列表返回
type YeZiESPortsGamesInfo struct {
	Data      []*share_message.TableESPortsGame `json:"data"`
	Code      int32                             `json:"code"`
	Msg       string                            `json:"msg"`
	ErrorCode int32                             `json:"error_code"`
	ErrorMsg  string                            `json:"error_msg"`
}

//6.1野子科技API两组对战比赛详情返回
type YeZiESPortsGameDetailInfo struct {
	Data      *share_message.TableESPortsGameDetail `json:"data"`
	Code      int32                                 `json:"code"`
	Msg       string                                `json:"msg"`
	ErrorCode int32                                 `json:"error_code"`
	ErrorMsg  string                                `json:"error_msg"`
}

//6.1野子科技API两组对战比赛详情返回Map中间结构
type YeZiESPortsMiddleDetailMapInfo struct {
	Data      *YeZiESPortsGameDetailMap `json:"data"`
	Code      int32                     `json:"code"`
	Msg       string                    `json:"msg"`
	ErrorCode int32                     `json:"error_code"`
	ErrorMsg  string                    `json:"error_msg"`
}

//野子科技API两组对战比赛详情结构(proto2.0无法定义map、先在程序中转然后再设置到TableESPortsGameDetail表结构中)
type YeZiESPortsGameDetailMap struct {
	TeamAPlayers   map[string]*share_message.APIPlayerDetail `json:"team_a_players"`  //A 队出场队员详情
	TeamBPlayers   map[string]*share_message.APIPlayerDetail `json:"team_b_players"`  //B 队出场队员详情
	LiveURL        map[string]map[string]string              `json:"live_url"`        //直播信号源第一层key:数字,可能会增加;第二层key：name,url,url_h5,name_h5
	WinProbability map[string]map[string]string              `json:"win_probability"` //两队历史交锋第一层key:"this_two_team两队历史交锋胜率"all两个队伍在当前比赛所处的赛事下的胜率;第二层key：队伍team_id
}

//两队对阵相关数据
//6.4野子科技API两组对战比赛历史对阵相关
//8.1 两队胜败统计接口
//8.2 两队天敌克制统计接口
type YeZiESPortsTeamBout struct {
	Data      *share_message.TableESPortsTeamBout `json:"data"`
	Code      int32                               `json:"code"`
	Msg       string                              `json:"msg"`
	ErrorCode int32                               `json:"error_code"`
	ErrorMsg  string                              `json:"error_msg"`
}

//8.1 两队胜败统计接口(先在程序中转然后再设置到TableESPortsTeamBout表结构中)
type YeZiESPortsTeamWinFailInfo struct {
	Data      *YeZiESPortsTeamWinFailData `json:"data"`
	Code      int32                       `json:"code"`
	Msg       string                      `json:"msg"`
	ErrorCode int32                       `json:"error_code"`
	ErrorMsg  string                      `json:"error_msg"`
}

type YeZiESPortsTeamWinFailData struct {
	TeamA *YeZiESPortsTeamWinFaiObject `json:"team_a"` //队伍a连胜败信息
	TeamB *YeZiESPortsTeamWinFaiObject `json:"team_b"` //队伍b连胜败信息
}

type YeZiESPortsTeamWinFaiObject struct {
	IsContinueWin int32  `json:"is_continue_win"` //连胜连败(-1：无 1： 连胜 0：连败)
	Num           int32  `json:"num"`             //连胜记录数
	TeamID        string `json:"team_id"`         //队伍id
}

//8.2 两队天敌克制统计接口(先在程序中转然后再设置到TableESPortsTeamBout表结构中)
type YeZiESPortsTeamNatResInfo struct {
	Data      *YeZiESPortsTeamNatResData `json:"data"`
	Code      int32                      `json:"code"`
	Msg       string                     `json:"msg"`
	ErrorCode int32                      `json:"error_code"`
	ErrorMsg  string                     `json:"error_msg"`
}

type YeZiESPortsTeamNatResData struct {
	TeamA *YeZiESPortsTeamNatResObject `json:"team_a"` //队伍a天敌克制信息
	TeamB *YeZiESPortsTeamNatResObject `json:"team_b"` //队伍b天敌克制信息
}

type YeZiESPortsTeamNatResObject struct {
	NaturalTeam  string `json:"naturalTeam"`  //天敌
	RestrainTeam string `json:"restrainTeam"` //克制
	TeamID       string `json:"team_id"`      //队伍id
}

//9.0 、10.0比赛动态信息接口（早盘、滚盘）
type YeZiESPortsGameGuess struct {
	Data      *share_message.TableESPortsGameGuess `json:"data"`
	Code      int32                                `json:"code"`
	Msg       string                               `json:"msg"`
	ErrorCode int32                                `json:"error_code"`
	ErrorMsg  string                               `json:"error_msg"`
}

//10.1 使用滚盘
type YeZiESPortsUSEGuessROLL struct {
	Data      []string `json:"data"`
	Code      int32    `json:"code"`
	Msg       string   `json:"msg"`
	ErrorCode int32    `json:"error_code"`
	ErrorMsg  string   `json:"error_msg"`
}

//回调用请求体内body
type CallBackBody struct {
	GameId     int32  `json:"game_id"`
	EventId    int32  `json:"event_id"`
	Type       string `json:"type"`
	Func       string `json:"func"`
	UpdateTime int64  `json:"update_time"`
	BetId      int32  `json:"bet_id"` //冲正回调时候用
}

//12.7 冲正回调后调用19.0接口返回结果
type YeZiESPortsGameBetInfo struct {
	Data      *share_message.ApiGuessObject `json:"data"`
	Code      int32                         `json:"code"`
	Msg       string                        `json:"msg"`
	ErrorCode int32                         `json:"error_code"`
	ErrorMsg  string                        `json:"error_msg"`
}

//15.0 取得推流直播地址
type YeZiESPortsGetGameLiveUrl struct {
	Data      *YeZiESPortsGameLiveUrlObject `json:"data"`
	Code      int32                         `json:"code"`
	Msg       string                        `json:"msg"`
	ErrorCode int32                         `json:"error_code"`
	ErrorMsg  string                        `json:"error_msg"`
}

type YeZiESPortsGameLiveUrlObject struct {
	GameId    string                                `json:"game_id"`    //比赛id
	LivePaths *share_message.ESPortsGameLivePathObj `json:"live_paths"` //直播地址
}
