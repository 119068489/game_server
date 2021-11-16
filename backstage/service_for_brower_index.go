// 管理后台为[浏览器]提供的服务

package backstage

import (
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/pb/brower_backstage"
	_ "game_server/pb/brower_backstage"
	"game_server/pb/share_message"
	"sort"
)

// 查询客户端版本管理
func (self *cls4) RpcQueryDataOverview(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	regSumCount, logSumCount := GetRegAndLogSumCount(user)
	yestodayStime := easygo.Get0ClockTimestamp(easygo.NowTimestamp() - 86400)
	report := for_game.QueryRegisterLoginReport(yestodayStime, user.GetRole())
	data := &brower_backstage.DataOverview{
		RegCount:      easygo.NewInt64(report.GetRegSumCount()),
		LoginCount:    easygo.NewInt64(report.GetLoginSumCount()),
		RegSumCount:   easygo.NewInt64(regSumCount),
		LoginSumCount: easygo.NewInt64(logSumCount),
		PvCount:       easygo.NewInt64(report.GetPvCount()),
		UvCount:       easygo.NewInt64(report.GetUvCount()),
	}
	return data
}

//查询运营渠道汇总报表曲线图
func (self *cls4) RpcRegisterLoginReportLine(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list, _ := GetRegisterLoginReport(reqMsg, user)
	line := &brower_backstage.LineData{}

	for _, k := range list {
		// 折线图数据
		line.TimeData = append(line.TimeData, k.GetCreateTime())
		switch reqMsg.GetListType() {
		case 1:
			line.VelueData = append(line.VelueData, k.GetRegSumCount())
		case 2:
			line.VelueData = append(line.VelueData, k.GetLoginSumCount())
		case 3:
			line.VelueData = append(line.VelueData, k.GetPvCount())
		case 4:
			line.VelueData = append(line.VelueData, k.GetUvCount())
		default:
			return easygo.NewFailMsg("查询类型错误")
		}
	}

	return &brower_backstage.LineChartResponse{
		Line: line,
	}
}

type tempSort struct {
	Id    int
	Name  string
	Value int64
}

//兴趣爱好柱状图
func (self *cls4) RpcInterestTagLine(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list := GetInterestTagListNopage(reqMsg)
	line := &brower_backstage.LineData{}
	temp := make([]*tempSort, 0)
	var sumCount int64
	var tags []int32
	tagMap := make(map[int32]string)
	for _, k := range list {
		tags = append(tags, k.GetKey())
		tagMap[k.GetKey()] = k.GetValue()
	}
	lis := GetInterestTagSumCount(tags)
	for i, k := range lis {
		m := &tempSort{
			Id:    i,
			Name:  tagMap[int32(k.GetId())],
			Value: k.GetCount(),
		}
		temp = append(temp, m)
	}

	sort.Slice(temp, func(i int, j int) bool {
		return temp[i].Value > temp[j].Value
	})

	for _, v := range temp {
		// 折线图数据
		sumCount += v.Value
		line.VelueData = append(line.VelueData, v.Value)
		line.Name = append(line.Name, v.Name)
	}

	line.Total = easygo.NewInt64(sumCount)

	return &brower_backstage.LineChartResponse{
		Line: line,
	}
}

//手机品牌柱状图
func (self *cls4) RpcPhoneBrandLine(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	brands := []string{"HUAWEI", "iPhone", "Xiaomi", "Meizu", "OPPO", "vivo", "HONOR", "samsung", "Redmi"}
	line := []*brower_backstage.NameValueTag{}
	li := &brower_backstage.NameValueTag{}
	var sumCount int64
	for _, k := range brands {
		velue := GetPhoneBrandSumCount(k)
		sumCount += velue
		li = &brower_backstage.NameValueTag{
			Name:  easygo.NewString(k),
			Value: easygo.NewInt64(velue),
		}
		// 折线图数据
		line = append(line, li)
	}

	sort.Slice(line, func(i int, j int) bool {
		return line[i].GetValue() > line[j].GetValue()
	})

	regCount, _ := GetRegAndLogSumCount(user)
	otherCount := regCount - sumCount
	li = &brower_backstage.NameValueTag{
		Name:  easygo.NewString("其他"),
		Value: easygo.NewInt64(otherCount),
	}
	line = append(line, li)

	return &brower_backstage.NameValueResponseTag{
		List:  line,
		Total: easygo.NewInt64(regCount),
	}
}

//上网热度柱状图
func (self *cls4) RpcPlayerOnlineLine(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list := GetPlayerOnlineLineReport(reqMsg)
	line := &brower_backstage.LineData{}
	var sumCount int64
	var clock0, clock1, clock2, clock3, clock4, clock5, clock6, clock7, clock8, clock9, clock10, clock11, clock12, clock13, clock14, clock15, clock16, clock17, clock18, clock19, clock20, clock21, clock22, clock23 int64
	for _, k := range list {
		clock0 += k.GetClock0()
		clock1 += k.GetClock1()
		clock2 += k.GetClock2()
		clock3 += k.GetClock3()
		clock4 += k.GetClock4()
		clock5 += k.GetClock5()
		clock6 += k.GetClock6()
		clock7 += k.GetClock7()
		clock8 += k.GetClock8()
		clock9 += k.GetClock9()
		clock10 += k.GetClock10()
		clock11 += k.GetClock11()
		clock12 += k.GetClock12()
		clock13 += k.GetClock13()
		clock14 += k.GetClock14()
		clock15 += k.GetClock15()
		clock16 += k.GetClock16()
		clock17 += k.GetClock17()
		clock18 += k.GetClock18()
		clock19 += k.GetClock19()
		clock20 += k.GetClock20()
		clock21 += k.GetClock21()
		clock22 += k.GetClock22()
		clock23 += k.GetClock23()
	}
	sumCount = (clock0 + clock1 + clock2 + clock3 + clock4 + clock5 + clock6 + clock7 + clock8 + clock9 + clock10 + clock11 + clock12 + clock13 + clock14 + clock15 + clock16 + clock17 + clock18 + clock19 + clock20 + clock21 + clock22 + clock23)
	line.VelueData = append(line.VelueData, clock0, clock1, clock2, clock3, clock4, clock5, clock6, clock7, clock8, clock9, clock10, clock11, clock12, clock13, clock14, clock15, clock16, clock17, clock18, clock19, clock20, clock21, clock22, clock23)
	for i := 0; i < 24; i++ {
		line.Name = append(line.Name, easygo.AnytoA(i)+"时")
	}
	line.Total = easygo.NewInt64(sumCount)

	return &brower_backstage.LineChartResponse{
		Line: line,
	}
}

//用户登录地区分布图
func (self *cls4) RpcPlayerLogLocation(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *brower_backstage.ListRequest) easygo.IMessage {
	list := GetPlayerLogLocation(reqMsg)
	line := []*brower_backstage.NameValueTag{}
	var total int64 = 0
	for _, k := range list {
		lineOne := &brower_backstage.NameValueTag{
			Name:  easygo.NewString(k.Id),
			Value: easygo.NewInt64(k.Count),
		}
		total += k.Count
		line = append(line, lineOne)
	}

	sort.Slice(line, func(i int, j int) bool {
		return line[i].GetValue() > line[j].GetValue()
	})

	return &brower_backstage.NameValueResponseTag{
		List:  line,
		Total: easygo.NewInt64(total),
	}
}

//用户画像
func (self *cls4) RpcPlayerPortrait(ep IBrowerEndpoint, user *share_message.Manager, reqMsg *base.Empty) easygo.IMessage {
	line := []*brower_backstage.NameValueTag{}
	var sumCount int64
	list := GetPlayerRegLocation()
	for _, k := range list {
		sumCount += k.Count
		lineOne := &brower_backstage.NameValueTag{
			Name:  easygo.NewString(k.Id),
			Value: easygo.NewInt64(k.Count),
		}
		line = append(line, lineOne)
	}

	sort.Slice(line, func(i int, j int) bool {
		return line[i].GetValue() > line[j].GetValue()
	})

	lis := GetPlayerGenderCount()
	var man, woman int64
	for _, s := range lis {
		if *s.Id == 1 {
			man = *s.Count
		}
		if *s.Id == 2 {
			woman = *s.Count
		}
	}

	return &brower_backstage.PlayerPortraitResponse{
		List:       line,
		Total:      easygo.NewInt64(sumCount),
		ManCount:   easygo.NewInt64(man),
		WomanCount: easygo.NewInt64(woman),
	}
}
