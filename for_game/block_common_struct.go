package for_game

import "game_server/pb/share_message"

type TFserver struct {
	Hostlist []string
}

//推文
type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

//金币日志结构
type GoldChangeLog struct {
	LogId      int64                         `json:"_id,omitempty" bson:"_id,omitempty"`
	PlayerId   int64                         `bson:"PlayerId,omitempty"`   //用户ID
	ChangeGold float64                       `bson:"ChangeGold,omitempty"` //变化金币
	SourceType int32                         `bson:"SourceType,omitempty"` //源类型(支付类型)
	PayType    int32                         `bson:"PayType,omitempty"`    //1收入，2支出
	CurGold    float64                       `bson:"CurGold,omitempty"`    //变化前携带金币
	Gold       float64                       `bson:"Gold,omitempty"`       //变化后携带金币
	Note       string                        `bson:"Note,omitempty"`       //备注
	CreateTime int64                         `bson:"CreateTime,omitempty"` //流水时间
	Extend     *share_message.RechargeExtend `bson:"Extend,omitempty"`     //扩展数据
}

//硬币日志结构
type CoinChangeLog struct {
	LogId      int64                        `json:"_id,omitempty" bson:"_id,omitempty"`
	PlayerId   int64                        `bson:"PlayerId,omitempty"`   //用户ID
	ChangeCoin float64                      `bson:"ChangeCoin,omitempty"` //变化硬币
	SourceType int32                        `bson:"SourceType,omitempty"` //源类型(支付类型)
	PayType    int32                        `bson:"PayType,omitempty"`    //1收入，2支出
	CurCoin    float64                      `bson:"CurCoin,omitempty"`    //变化前携带硬币
	Coin       float64                      `bson:"Coin,omitempty"`       //变化后携带硬币
	Note       string                       `bson:"Note,omitempty"`       //备注
	CreateTime int64                        `bson:"CreateTime,omitempty"` //流水时间
	Extend     *share_message.GoldExtendLog `bson:"Extend,omitempty"`     //扩展数据
}

//电竞币日志结构
type ESportCoinChangeLog struct {
	LogId            int64                        `json:"_id,omitempty" bson:"_id,omitempty"`
	PlayerId         int64                        `bson:"PlayerId,omitempty"`   //用户ID
	ChangeESportCoin float64                      `bson:"ChangeCoin,omitempty"` //变化电竞币
	SourceType       int32                        `bson:"SourceType,omitempty"` //源类型(支付类型)
	PayType          int32                        `bson:"PayType,omitempty"`    //1收入，2支出
	CurESportCoin    float64                      `bson:"CurCoin,omitempty"`    //变化前携带电竞币
	ESportCoin       float64                      `bson:"Coin,omitempty"`       //变化后携带电竞币
	Note             string                       `bson:"Note,omitempty"`       //备注
	CreateTime       int64                        `bson:"CreateTime,omitempty"` //流水时间
	Extend           *share_message.GoldExtendLog `bson:"Extend,omitempty"`     //扩展数据
}

//金币日志结构
type GoldLog struct {
	share_message.GoldChangeLog `bson:",inline,omitempty"`
	Extend                      *share_message.GoldExtendLog `bson:"Extend,omitempty"` //扩展数据
}
type CommonGold struct {
	share_message.GoldChangeLog `bson:",inline,omitempty"` // inline 类型不能用指针
	Extend                      interface{}                `bson:"Extend,omitempty"`
}

//硬币日志结构
type CoinLog struct {
	share_message.CoinChangeLog `bson:",inline,omitempty"`
	Extend                      *share_message.GoldExtendLog `bson:"Extend,omitempty"` //扩展数据
}
type CommonCoin struct {
	share_message.CoinChangeLog `bson:",inline,omitempty"` // inline 类型不能用指针
	Extend                      interface{}                `bson:"Extend,omitempty"`
}

//电竞币日志结构
type ESportCoinLog struct {
	share_message.ESportCoinChangeLog `bson:",inline,omitempty"`
	Extend                            *share_message.GoldExtendLog `bson:"Extend,omitempty"` //扩展数据
}
type CommonESportCoin struct {
	share_message.ESportCoinChangeLog `bson:",inline,omitempty"` // inline 类型不能用指针
	Extend                            interface{}                `bson:"Extend,omitempty"`
}

type DynamicData struct {
	LogId             int64 `json:"_id"`
	PlayerId          int64
	HeadIcon          string
	Sex               int32
	Content           string
	Photo             []string
	Zan               int32
	IsZan             bool
	IsAtten           bool
	Voice             string
	Video             string
	CreateTime        int64
	CommentNum        int64
	CommentList       *share_message.CommentList
	TrueZan           int32
	Statue            int32
	VoiceTime         int64
	NickName          string
	Account           string
	High              float64
	Weight            float64
	TopOverTime       int64
	IsBsTop           bool
	IsShield          bool
	Note              string
	IsTop             bool
	VideoThumbnailURL string
}

type LocationData struct {
	Area      string `json:"area"`
	City      string `json:"city"`
	CityId    string `json:"city_id"`
	Country   string `json:"country"`
	CountryId string `json:"country_id"`
	Region    string `json:"region"`
	RegionId  string `json:"region_id"`
}

type PipeIntCount struct {
	Id    *int32 `bson:"_id"`
	Count *int64 `bson:"Count"`
}

type PipeStringCount struct {
	Id    string `bson:"_id"`
	Count int64  `bson:"Count"`
}

type PipeStringList struct {
	Id   string  `bson:"_id"`
	List []int64 `bson:"List"`
}

type PlayerGroup struct {
	Id       int64 `bson:"_id"`
	PlayerId int64 `bson:"PlayerId"`
}

// 玩家任务列表
type LuckyPlayerTask struct {
	PlayerId      int64 `json:"PlayerId"`
	IsSignIn      bool  `json:"IsSignIn"`      // 今日是否签到
	IsShare       bool  `json:"IsShare"`       // 今日是否已分享
	IsSendDynamic bool  `json:"IsSendDynamic"` // 今日是否发布了动态
	IsRedpacket   bool  `json:"IsRedpacket"`   // 今日是否发布了红包
}

type ResultMsg struct {
	Data RecentData `json:"data"`
}

type RecentData struct {
	Id                 int64              `json:"_id" bson:"_id,omitempty"`
	TournamentBiaoxian TournamentBiaoxian `json:"tournament_biaoxian" bson:"tournament_biaoxian"`
	StrengthIndex      StrengthIndex      `json:"strength_index" bson:"strength_index"`
	MatchData          []MatchData        `json:"match_data" bson:"match_data"`
	MatchRecord        MatchRecord        `json:"match_record" bson:"match_record"`
	MatchRecordA       MatchRecord        `json:"match_record_a" bson:"match_record_a"`
	MatchRecordB       MatchRecord        `json:"match_record_b" bson:"match_record_b"`
}

type TeamX struct {
	TeamID            string      `json:"teamID" bson:"teamID"`
	KDA               string      `json:"KDA" bson:"KDA"`
	AVERAGEKILLS      string      `json:"AVERAGE_KILLS" bson:"AVERAGE_KILLS"`
	AVERAGEASSISTS    string      `json:"AVERAGE_ASSISTS" bson:"AVERAGE_ASSISTS"`
	AVERAGEDEATHS     string      `json:"AVERAGE_DEATHS" bson:"AVERAGE_DEATHS"`
	MINUTEHITS        string      `json:"MINUTE_HITS" bson:"MINUTE_HITS"`
	MINUTEECONOMIC    string      `json:"MINUTE_ECONOMIC" bson:"MINUTE_ECONOMIC"`
	MINUTEDAMAGEDEALT interface{} `json:"MINUTE_DAMAGEDEALT" bson:"MINUTE_DAMAGEDEALT"`
	SMALLDRAGONRATE   float64     `json:"SMALLDRAGON_RATE" bson:"SMALLDRAGON_RATE"`
	BIGDRAGONRATE     float64     `json:"BIGDRAGON_RATE" bson:"BIGDRAGON_RATE"`
	VICTORYRATE       float64     `json:"VICTORY_RATE" bson:"VICTORY_RATE"`
	VICTORYCOUNT      float64     `json:"VICTORY_COUNT" bson:"VICTORY_COUNT"`
	FAIlCOUNT         float64     `json:"FAIl_COUNT" bson:"FAIl_COUNT"`
	ContinuityCount   float64     `json:"continuity_count" bson:"continuity_count"`
}
type TournamentBiaoxian struct {
	TeamA TeamX `json:"team_a" bson:"team_a"`
	TeamB TeamX `json:"team_b" bson:"team_b"`
}

type StrengthIndex struct {
	HandWinTeamA              string `json:"hand_win_team_a" bson:"hand_win_team_a"`
	HandWinTeamB              string `json:"hand_win_team_b" bson:"hand_win_team_b"`
	HandLoseTeamA             string `json:"hand_lose_team_a" bson:"hand_lose_team_a"`
	HandLoseTeamB             string `json:"hand_lose_team_b" bson:"hand_lose_team_b"`
	RecordWinTeamA            string `json:"record_win_team_a" bson:"record_win_team_a"`
	RecordWinTeamB            string `json:"record_win_team_b" bson:"record_win_team_b"`
	RecordLoseTeamA           string `json:"record_lose_team_a" bson:"record_lose_team_a"`
	RecordLoseTeamB           string `json:"record_lose_team_b" bson:"record_lose_team_b"`
	AverageKillsTeamA         string `json:"average_kills_team_a" bson:"average_kills_team_a"`
	AverageKillsTeamB         string `json:"average_kills_team_b" bson:"average_kills_team_b"`
	AverageTowerTeamA         string `json:"average_tower_team_a" bson:"average_tower_team_a"`
	AverageTowerTeamB         string `json:"average_tower_team_b" bson:"average_tower_team_b"`
	AverageMoneyTeamA         string `json:"average_money_team_a" bson:"average_money_team_a"`
	AverageMoneyTeamB         string `json:"average_money_team_b" bson:"average_money_team_b"`
	ScoreTeamA                string `json:"score_team_a" bson:"score_team_a"`
	ScoreTeamB                string `json:"score_team_b" bson:"score_team_b"`
	VictoryRateA              string `json:"victory_rate_a" bson:"victory_rate_a"`
	VictoryRateB              string `json:"victory_rate_b" bson:"victory_rate_b"`
	AverageTimeA              string `json:"average_time_a" bson:"average_time_a"`
	AverageTimeB              string `json:"average_time_b" bson:"average_time_b"`
	AverageAssistsTeamA       string `json:"average_assists_team_a" bson:"average_assists_team_a"`
	AverageAssistsTeamB       string `json:"average_assists_team_b" bson:"average_assists_team_b"`
	AverageDeathsTeamA        string `json:"average_deaths_team_a" bson:"average_deaths_team_a"`
	AverageDeathsTeamB        string `json:"average_deaths_team_b" bson:"average_deaths_team_b"`
	AverageKdaTeamA           string `json:"average_kda_team_a" bson:"average_kda_team_a"`
	AverageKdaTeamB           string `json:"average_kda_team_b" bson:"average_kda_team_b"`
	FirstBloodKillTeamA       string `json:"firstBloodKill_team_a" bson:"firstBloodKill_team_a"`
	FirstBloodKillTeamB       string `json:"firstBloodKill_team_b" bson:"firstBloodKill_team_b"`
	MinuteDamageTeamA         string `json:"minute_damage_team_a" bson:"minute_damage_team_a"`
	MinuteDamageTeamB         string `json:"minute_damage_team_b" bson:"minute_damage_team_b"`
	FirstTowerKillTeamA       string `json:"firstTowerKill_team_a" bson:"firstTowerKill_team_a"`
	FirstTowerKillTeamB       string `json:"firstTowerKill_team_b" bson:"firstTowerKill_team_b"`
	AverageMoneyDiffTeamA     string `json:"average_money_diff_team_a" bson:"average_money_diff_team_a"`
	AverageMoneyDiffTeamB     string `json:"average_money_diff_team_b" bson:"average_money_diff_team_b"`
	MinuteMoneyTeamA          string `json:"minute_money_team_a" bson:"minute_money_team_a"`
	MinuteMoneyTeamB          string `json:"minute_money_team_b" bson:"minute_money_team_b"`
	MinuteHitsTeamA           string `json:"minute_hits_team_a" bson:"minute_hits_team_a"`
	MinuteHitsTeamB           string `json:"minute_hits_team_b" bson:"minute_hits_team_b"`
	AverageDragonTeamA        string `json:"average_dragon_team_a" bson:"average_dragon_team_a"`
	AverageDragonTeamB        string `json:"average_dragon_team_b" bson:"average_dragon_team_b"`
	AverageBaronTeamA         string `json:"average_baron_team_a" bson:"average_baron_team_a"`
	AverageBaronTeamB         string `json:"average_baron_team_b" bson:"average_baron_team_b"`
	RateDragonTeamA           string `json:"rate_dragon_team_a" bson:"rate_dragon_team_a"`
	RateDragonTeamB           string `json:"rate_dragon_team_b" bson:"rate_dragon_team_b"`
	RateBaronTeamA            string `json:"rate_baron_team_a" bson:"rate_baron_team_a"`
	RateBaronTeamB            string `json:"rate_baron_team_b" bson:"rate_baron_team_b"`
	MinuteWardsPlacedTeamA    string `json:"minute_wardsPlaced_team_a" bson:"minute_wardsPlaced_team_a"`
	MinuteWardsPlacedTeamB    string `json:"minute_wardsPlaced_team_b" bson:"minute_wardsPlaced_team_b"`
	MinuteWardsKilledTeamA    string `json:"minute_wardsKilled_team_a" bson:"minute_wardsKilled_team_a"`
	MinuteWardsKilledTeamB    string `json:"minute_wardsKilled_team_b" bson:"minute_wardsKilled_team_b"`
	AverageBeTurretKillsTeamA string `json:"average_be_turretKills_team_a" bson:"average_be_turretKills_team_a"`
	AverageBeTurretKillsTeamB string `json:"average_be_turretKills_team_b" bson:"average_be_turretKills_team_b"`
	RateFullBureauTeamA       string `json:"rate_full_bureau_team_a" bson:"rate_full_bureau_team_a"`
	RateFullBureauTeamB       string `json:"rate_full_bureau_team_b" bson:"rate_full_bureau_team_b"`
	TotalDragonTeamA          string `json:"total_dragon_team_a" bson:"total_dragon_team_a"`
	TotalBaronTeamA           string `json:"total_baron_team_a" bson:"total_baron_team_a"`
	TotalDragonTeamB          string `json:"total_dragon_team_b" bson:"total_dragon_team_b"`
	TotalBaronTeamB           string `json:"total_baron_team_b" bson:"total_baron_team_b"`
}
type HeroWinLose struct {
	HeroID    string `json:"heroID" bson:"heroID"`
	HeroImage string `json:"hero_image" bson:"hero_image"`
	Win       int32  `json:"win" bson:"win"`
	Lose      int32  `json:"lose" bson:"lose"`
}
type MatchData struct {
	PlayerID                 string        `json:"playerID" bson:"playerID"`
	Nickname                 string        `json:"nickname" bson:"nickname"`
	PositionID               string        `json:"positionID" bson:"positionID"`
	TeamID                   string        `json:"teamID" bson:"teamID"`
	PlayerImageThumb         string        `json:"player_image_thumb" bson:"player_image_thumb"`
	PositionName             string        `json:"position_name" bson:"position_name"`
	CountryID                string        `json:"country_id" bson:"country_id"`
	CountryImage             string        `json:"country_image" bson:"country_image"`
	TeamShortName            string        `json:"team_short_name" bson:"team_short_name"`
	TeamImageThumb           string        `json:"team_image_thumb" bson:"team_image_thumb"`
	PlayerChineseName        string        `json:"player_chinese_name" bson:"player_chinese_name"`
	Total                    string        `json:"total" bson:"total"`
	DamageDealPercent        string        `json:"DamageDealPercent" bson:"DamageDealPercent"`
	TeamPercent              string        `json:"TeamPercent" bson:"TeamPercent"`
	MINUTEDAMAGEDEALT        interface{}   `json:"MINUTE_DAMAGEDEALT" bson:"MINUTE_DAMAGEDEALT"`
	AVERAGEKILLS             string        `json:"AVERAGE_KILLS" bson:"AVERAGE_KILLS"`
	AVERAGEASSISTS           string        `json:"AVERAGE_ASSISTS" bson:"AVERAGE_ASSISTS"`
	AVERAGEDEATHS            string        `json:"AVERAGE_DEATHS" bson:"AVERAGE_DEATHS"`
	KDA                      string        `json:"KDA" bson:"KDA"`
	MINUTEECONOMIC           interface{}   `json:"MINUTE_ECONOMIC" bson:"MINUTE_ECONOMIC"`
	MINUTEWARDSPLACED        string        `json:"MINUTE_WARDSPLACED" bson:"MINUTE_WARDSPLACED"`
	AVERAGEMinionsKilled     string        `json:"AVERAGE_MinionsKilled" bson:"AVERAGE_MinionsKilled"`
	AVERAGELife              string        `json:"AVERAGE_Life" bson:"AVERAGE_Life"`
	TotalDamageTaken         interface{}   `json:"totalDamageTaken" bson:"totalDamageTaken"`
	NeutralMinionsKilled     string        `json:"neutralMinionsKilled" bson:"neutralMinionsKilled"`
	WardsPlaced              string        `json:"wardsPlaced" bson:"wardsPlaced"`
	WardsKilled              string        `json:"wardsKilled" bson:"wardsKilled"`
	StatusID                 string        `json:"statusID" bson:"statusID"`
	WinCount                 interface{}   `json:"win_count" bson:"win_count"`
	LoseCount                interface{}   `json:"lose_count" bson:"lose_count"`
	LastMatchTime            string        `json:"last_match_time" bson:"last_match_time"`
	WinLose                  []string      `json:"win_lose" bson:"win_lose"`
	VICTORYRATE              string        `json:"VICTORY_RATE" bson:"VICTORY_RATE"`
	TeamType                 string        `json:"team_type" bson:"team_type"`
	NeutralMinionsKilledRate string        `json:"neutralMinionsKilled_rate" bson:"neutralMinionsKilled_rate"`
	HeroWinLose              []HeroWinLose `json:"hero_win_lose" bson:"hero_win_lose"`
}
type MatchRecordResultList struct {
	ResultID      string `json:"resultID" bson:"resultID"`
	WinTeamID     string `json:"win_teamID" bson:"win_teamID"`
	TeamName      string `json:"team_name" bson:"team_name"`
	TeamShortName string `json:"team_short_name" bson:"team_short_name"`
	TeamImage     string `json:"team_image" bson:"team_image"`
	Bo            string `json:"bo" bson:"bo"`
}
type MatchRecordList struct {
	WinTeamID      string                  `json:"win_team_id" bson:"win_team_id"`
	WinTeamName    string                  `json:"win_team_name" bson:"win_team_name"`
	StartTime      string                  `json:"start_time" bson:"start_time"`
	MatchID        string                  `json:"matchID" bson:"matchID"`
	Title          string                  `json:"title" bson:"title"`
	TeamIDA        string                  `json:"teamID_a" bson:"teamID_a"`
	TeamIDB        string                  `json:"teamID_b" bson:"teamID_b"`
	TeamAImage     string                  `json:"team_a_image" bson:"team_a_image"`
	TeamBImage     string                  `json:"team_b_image" bson:"team_b_image"`
	TeamAShortName string                  `json:"team_a_short_name" bson:"team_a_short_name"`
	TeamBShortName string                  `json:"team_b_short_name" bson:"team_b_short_name"`
	TeamAWin       string                  `json:"team_a_win" bson:"team_a_win"`
	TeamBWin       string                  `json:"team_b_win" bson:"team_b_win"`
	ResultList     []MatchRecordResultList `json:"result_list" bson:"result_list"`
}
type MatchRecord struct {
	TeamAWinCount int32             `json:"team_a_win_count" bson:"team_a_win_count"`
	TeamBWinCount int32             `json:"team_b_win_count" bson:"team_b_win_count"`
	List          []MatchRecordList `json:"list" bson:"list"`
}
