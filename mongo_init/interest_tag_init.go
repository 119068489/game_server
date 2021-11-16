//
// 初始化标签数据
package mongo_init

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/easygo/util"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
)

//兴趣类型
func InitInterestTagTypeCfg() []interface{} {
	var types []interface{}
	types = append(types, &share_message.InterestType{
		Id:         easygo.NewInt32(1),
		Name:       easygo.NewString("生活"),
		Sort:       easygo.NewInt32(1),
		UpdateTime: easygo.NewInt64(easygo.NowTimestamp()),
		Status:     easygo.NewInt32(0),
	})
	types = append(types, &share_message.InterestType{
		Id:         easygo.NewInt32(2),
		Name:       easygo.NewString("兴趣"),
		Sort:       easygo.NewInt32(2),
		UpdateTime: easygo.NewInt64(easygo.NowTimestamp()),
		Status:     easygo.NewInt32(0),
	})
	return types
}

//兴趣标签
func InitInterestTagCfg() []interface{} {
	var jsonData = []byte(`{
		"InterestTag": [
			{
				"_id": 1,
				"Name": "恋爱",
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499451000.png",
				"UpdateTime":1600499450,
				"Sort": 1,
				"InterestType": 1,
				"Status": 0,
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499435000.gif"
			},
			{
				"_id": 2,
				"Name": "灵魂社交",
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499468000.png",
				"UpdateTime":1600500153,
				"Sort": 10,
				"InterestType": 1,
				"Status": 0,
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499501000.gif"
			},
			{
				"_id": 3,
				"Name": "美女",
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/1588247767000.png",
				"UpdateTime":1602743093,
				"Sort": 13,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 0
			},
			{
				"_id": 4,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600498955000.png",
				"Name": "游戏",
				"Sort": 14,
				"UpdateTime":1602743115,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600498967000.gif"
			},
			{
				"_id": 5,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499092000.png",
				"Name": "电影",
				"Sort": 6,
				"UpdateTime":1600499807,
				"InterestType": 1,
				"Status": 0,
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499102000.gif"
			},
			{
				"_id": 6,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/1590115800000.png",
				"Name": "动漫",
				"Sort": 16,
				"UpdateTime":1602743133,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 0
			},
			{
				"_id": 7,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499566000.png",
				"Name": "KTV",
				"Sort": 8,
				"UpdateTime":1600512920,
				"InterestType": 1,
				"Status": 0,
				"PlayTime": 2700,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499560000.gif"
			},
			{
				"_id": 8,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499074000.png",
				"Name": "美食",
				"Sort": 18,
				"UpdateTime":1602743145,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 0,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499059000.gif"
			},
			{
				"_id": 9,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499004000.png",
				"Name": "美妆",
				"Sort": 11,
				"UpdateTime":1600499937,
				"InterestType": 1,
				"Status": 0,
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499019000.gif"
			},
			{
				"_id": 10,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/1590115879000.png",
				"Name": "小说",
				"Sort": 110,
				"UpdateTime":1602743182,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 0
			},
			{
				"_id": 11,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499618000.png",
				"Name": "健身",
				"Sort": 3,
				"UpdateTime":1600512898,
				"InterestType": 1,
				"Status": 0,
				"PlayTime": 2700,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499640000.gif"
			},
			{
				"_id": 12,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499263000.png",
				"Name": "金融",
				"Sort": 27,
				"UpdateTime":1602743177,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 0,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499253000.gif"
			},
			{
				"_id": 13,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499222000.png",
				"Name": "电竞",
				"Sort": 13,
				"UpdateTime":1602743104,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499237000.gif"
			},
			{
				"_id": 14,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/1590115986000.png",
				"Name": "舞蹈",
				"Sort": 14,
				"UpdateTime":1602743122,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 0
			},
			{
				"_id": 15,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/1590115997000.png",
				"Name": "搞笑",
				"Sort": 15,
				"UpdateTime":1602743127,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 0
			},
			{
				"_id": 16,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/1590116012000.png",
				"Name": "时尚",
				"Sort": 16,
				"UpdateTime":1602743139,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 0
			},
			{
				"_id": 17,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499337000.png",
				"Name": "旅行",
				"Sort": 7,
				"UpdateTime":1600512875,
				"InterestType": 1,
				"Status": 0,
				"PlayTime": 3200,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499829000.gif"
			},
			{
				"_id": 18,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/1590116049000.png",
				"Name": "摄影",
				"Sort": 18,
				"UpdateTime":1602743152,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 0
			},
			{
				"_id": 19,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/1590116085000.png",
				"Name": "绘画",
				"Sort": 19,
				"UpdateTime":1602743157,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 0
			},
			{
				"_id": 20,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600500095000.png",
				"Name": "宠物",
				"Sort": 2,
				"UpdateTime":1600500094,
				"InterestType": 1,
				"Status": 0,
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600500083000.gif"
			},
			{
				"_id": 21,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499961000.png",
				"Name": "萌娃",
				"Sort": 12,
				"UpdateTime":1600500163,
				"InterestType": 1,
				"Status": 0,
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499969000.gif"
			},
			{
				"_id": 22,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/1590116136000.png",
				"Name": "养生",
				"Sort": 22,
				"UpdateTime":1602743162,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 0
			},
			{
				"_id": 23,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/1590116147000.png",
				"Name": "汽车",
				"Sort": 23,
				"UpdateTime":1602743166,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 0
			},
			{
				"_id": 24,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/1590116163000.png",
				"Name": "科技",
				"Sort": 23,
				"UpdateTime":1602743171,
				"InterestType": 1,
				"Status": 1,
				"PlayTime": 0
			},
			{
				"_id": 25,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600498899000.png",
				"InterestType": 2,
				"Name": "股票",
				"PlayTime": 3500,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600498935000.gif",
				"Sort": 1,
				"Status": 0,
				"UpdateTime":1600512971
			},
			{
				"_id": 26,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499057000.png",
				"InterestType": 2,
				"Name": "金融",
				"PlayTime": 3500,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499032000.gif",
				"Sort": 2,
				"Status": 0,
				"UpdateTime":1600512984
			},
			{
				"_id": 27,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499078000.png",
				"InterestType": 2,
				"Name": "投资理财",
				"PlayTime": 3500,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499086000.gif",
				"Sort": 3,
				"Status": 0,
				"UpdateTime":1600512988
			},
			{
				"_id": 28,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499130000.png",
				"InterestType": 2,
				"Name": "购物",
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499105000.gif",
				"Sort": 4,
				"Status": 0,
				"UpdateTime":1600499136
			},
			{
				"_id": 29,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499153000.png",
				"InterestType": 2,
				"Name": "摄影",
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499160000.gif",
				"Sort": 5,
				"Status": 0,
				"UpdateTime":1600499166
			},
			{
				"_id": 30,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499192000.png",
				"InterestType": 2,
				"Name": "美食",
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499178000.gif",
				"Sort": 6,
				"Status": 0,
				"UpdateTime":1600499196
			},
			{
				"_id": 31,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499232000.png",
				"InterestType": 2,
				"Name": "阅读",
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499241000.gif",
				"Sort": 7,
				"Status": 0,
				"UpdateTime":1600499248
			},
			{
				"_id": 32,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499277000.png",
				"InterestType": 2,
				"Name": "极限运动",
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499261000.gif",
				"Sort": 8,
				"Status": 0,
				"UpdateTime":1600499279
			},
			{
				"_id": 33,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499294000.png",
				"InterestType": 2,
				"Name": "改装车",
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499300000.gif",
				"Sort": 9,
				"Status": 0,
				"UpdateTime":1600499307
			},
			{
				"_id": 34,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499341000.png",
				"InterestType": 2,
				"Name": "说唱",
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499327000.gif",
				"Sort": 10,
				"Status": 0,
				"UpdateTime":1600499343
			},
			{
				"_id": 35,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499361000.png",
				"InterestType": 2,
				"Name": "二次元",
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499368000.gif",
				"Sort": 11,
				"Status": 0,
				"UpdateTime":1600499374
			},
			{
				"_id": 36,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499414000.png",
				"InterestType": 2,
				"Name": "电竞",
				"PlayTime": 2700,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499394000.gif",
				"Sort": 12,
				"Status": 0,
				"UpdateTime":1600513028
			},
			{
				"_id": 37,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499432000.png",
				"InterestType": 2,
				"Name": "球鞋爱好者",
				"PlayTime": 2700,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499439000.gif",
				"Sort": 13,
				"Status": 0,
				"UpdateTime":1600513063
			},
			{
				"_id": 38,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499480000.png",
				"InterestType": 2,
				"Name": "古风",
				"PlayTime": 3500,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499468000.gif",
				"Sort": 14,
				"Status": 0,
				"UpdateTime":1600513052
			},
			{
				"_id": 39,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499503000.png",
				"InterestType": 2,
				"Name": "Lolita",
				"PlayTime": 3500,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499512000.gif",
				"Sort": 15,
				"Status": 0,
				"UpdateTime":1600513092
			},
			{
				"_id": 40,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499545000.png",
				"InterestType": 2,
				"Name": "球类",
				"PlayTime": 3500,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499535000.gif",
				"Sort": 16,
				"Status": 0,
				"UpdateTime":1600513088
			},
			{
				"_id": 41,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499564000.png",
				"InterestType": 2,
				"Name": "游戏",
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499581000.gif",
				"Sort": 17,
				"Status": 0,
				"UpdateTime":1600499587
			},
			{
				"_id": 42,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499615000.png",
				"InterestType": 2,
				"Name": "JK制服",
				"PlayTime": 3500,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499602000.gif",
				"Sort": 18,
				"Status": 0,
				"UpdateTime":1600513084
			},
			{
				"_id": 43,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499742000.png",
				"InterestType": 1,
				"Name": "蹦迪",
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499731000.gif",
				"Sort": 4,
				"Status": 0,
				"UpdateTime":1600499747
			},
			{
				"_id": 44,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499771000.png",
				"InterestType": 1,
				"Name": "Club",
				"PlayTime": 3000,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499777000.gif",
				"Sort": 5,
				"Status": 0,
				"UpdateTime":1600499780
			},
			{
				"_id": 45,
				"Icon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499890000.png",
				"InterestType": 1,
				"Name": "街舞",
				"PlayTime": 2700,
				"PopIcon": "https://im-resource-1253887233.file.myqcloud.com/backstage/picture/1600499874000.gif",
				"Sort": 9,
				"Status": 0,
				"UpdateTime":1600512939
			}
		]
	}`)

	data := &Initstruct{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		logs.Info(err)
	}

	var list []interface{}
	lis := data.InterestTag
	for _, li := range lis {
		list = append(list, li)
	}
	return list
}

func InitTopicCfg() []interface{} {
	var topicTypes []interface{}
	topicTypes = append(topicTypes, &share_message.TopicType{
		Id:         easygo.NewInt64(0),
		Name:       easygo.NewString("自定义"),
		TopicClass: easygo.NewInt32(2),
		CreateTime: easygo.NewInt64(util.GetMilliTime()),
		UpdateTime: easygo.NewInt64(util.GetMilliTime()),
		Sort:       easygo.NewInt64(0),
		Status:     easygo.NewInt32(1),
	})
	topicTypes = append(topicTypes, &share_message.TopicType{
		Name:       easygo.NewString("热门"),
		TopicClass: easygo.NewInt32(1),
		CreateTime: easygo.NewInt64(util.GetMilliTime()),
		UpdateTime: easygo.NewInt64(util.GetMilliTime()),
		Sort:       easygo.NewInt64(1),
		Status:     easygo.NewInt32(1),
	})
	return topicTypes
}

//星座标签
func InitStarSignsTag() []interface{} {
	var jsonData = []byte(`{
		"InterestTag": [
			{
				"_id": 1,
				"Name": "自信满满白羊座",
				"SortName":"白羊座",
				"Icon":"https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/png_1620297077000png"
			},
			{
				"_id": 2,
				"Name": "慢热达人金牛座",
				"SortName":"金牛座",
				"Icon":"https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/png_1620297125000png"
			},
			{
				"_id": 3,
				"Name": "好奇宝宝双子座",
				"SortName":"双子座",
				"Icon":"https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/png_1620297192000png"
			},
			{
				"_id": 4,
				"Name": "温柔体贴巨蟹座",
				"SortName":"巨蟹座",
				"Icon":"https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/png_1620297139000png"
			},
			{
				"_id": 5,
				"Name": "慷慨大方狮子座",
				"SortName":"狮子座",
				"Icon":"https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/png_1620297174000png"
			},
			{
				"_id": 6,
				"Name": "完美主义处女座",
				"SortName":"处女座",
				"Icon":"https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/png_1620297106000png"
			},
			{
				"_id": 7,
				"Name": "世界和平天枰座",
				"SortName":"天枰座",
				"Icon":"https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/png_1620297216000png"
			},
			{
				"_id": 8,
				"Name": "爱憎分明天蝎座",
				"SortName":"天蝎座",
				"Icon":"https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/png_1620297225000png"
			},
			{
				"_id": 9,
				"Name": "崇尚自由射手座",
				"SortName":"射手座",
				"Icon":"https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/png_1620297162000png"
			},
			{
				"_id": 10,
				"Name": "沉着冷静魔羯座",
				"SortName":"魔羯座",
				"Icon":"https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/png_1620297153000png"
			},
			{
				"_id": 11,
				"Name": "聪明过人水瓶座",
				"SortName":"水瓶座",
				"Icon":"https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/png_1620297206000png"
			},
			{
				"_id": 12,
				"Name": "梦幻达人双鱼座",
				"SortName":"双鱼座",
				"Icon":"https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/png_1620297183000png"
			}	
		]
	}`)

	data := &Initstruct{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		logs.Info(err)
	}

	var list []interface{}
	lis := data.InterestTag
	for _, li := range lis {
		list = append(list, li)
	}
	return list
}

//匹配引导语
func InitMatchGuideCfg() []interface{} {
	var list []interface{}
	list = append(list, &share_message.CommStrId{
		Id: easygo.NewString("缘份在敲门，快分享一下你们的择偶标准吧~"),
	})
	list = append(list, &share_message.CommStrId{
		Id: easygo.NewString("分享一件你觉得在恋爱中最浪漫的事情吧~"),
	})
	list = append(list, &share_message.CommStrId{
		Id: easygo.NewString("说说你最希望和另一半做什么事情吧~"),
	})

	return list
}

func InitSayHiCfg() []interface{} {
	var list []interface{}
	list = append(list, &share_message.CommStrId{
		Id: easygo.NewString("为你的声音打CALL"),
	})
	list = append(list, &share_message.CommStrId{
		Id: easygo.NewString("声音的天赋型选手出现了"),
	})
	list = append(list, &share_message.CommStrId{
		Id: easygo.NewString("你知道吗，我很少喜欢别人的声音"),
	})
	list = append(list, &share_message.CommStrId{
		Id: easygo.NewString("你的声音我很喜欢，可以交个朋友吗"),
	})

	return list
}

//初始化亲密度配置
func InitIntimacyConfig() []interface{} {
	var list []interface{}
	for i := 0; i < 6; i++ {
		conf := &share_message.IntimacyConfig{
			Lv:     easygo.NewInt32(i),
			MaxVal: easygo.NewInt64((i + 1) * 1000),
		}
		if i == 1 {
			conf.PerDayVal = easygo.NewInt32(20)
		} else {
			conf.PerDayVal = easygo.NewInt32(10)
		}
		list = append(list, conf)
	}
	return list
}

//初始化个性化标签
func InitSysPersonalityTags() []interface{} {
	var list []interface{}
	for i := 0; i < 30; i++ {
		conf := &share_message.InterestTag{
			Id:   easygo.NewInt32(i + 1),
			Name: easygo.NewString("个性化标签" + easygo.AnytoA(i)),
		}
		list = append(list, conf)
	}
	return list
}

//系统背景标签
//func InitSysBgVoiceTags() []interface{} {
//	var list []interface{}
//	for i := 0; i < 30; i++ {
//		conf := &share_message.InterestTag{
//			Id:   easygo.NewInt32(for_game.NextId(for_game.TABLE_BG_VOICE_TAG)),
//			Name: easygo.NewString("系统标签" + easygo.AnytoA(i)),
//		}
//		list = append(list, conf)
//	}
//	return list
//}
//系统背景图
func InitSysBgImageTags() []interface{} {
	var list []interface{}
	list = append(list, &share_message.SystemBgImage{Url: easygo.NewString("https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/bon27.jpeg")})
	list = append(list, &share_message.SystemBgImage{Url: easygo.NewString("https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/bon3.jpeg")})
	list = append(list, &share_message.SystemBgImage{Url: easygo.NewString("https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/bon8.jpeg")})
	list = append(list, &share_message.SystemBgImage{Url: easygo.NewString("https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/jkl1.png")})
	list = append(list, &share_message.SystemBgImage{Url: easygo.NewString("https://im-resource-1253887233.file.myqcloud.com/backstage/imgResource/jkl17.jpeg")})
	return list

}

//恋爱匹配运营号初始化
//匹配运营号初始化
func InitOperatePlayer() []interface{} {
	players := make([]interface{}, 0)
	players = append(players, &share_message.PlayerOperate{
		Account:  easygo.NewString("lm7081fdfb"),
		PlayerId: easygo.NewInt64(1887567355),
		Type:     easygo.NewInt32(1),
	})
	players = append(players, &share_message.PlayerOperate{
		Account:  easygo.NewString("lm7081fdfc"),
		PlayerId: easygo.NewInt64(1887567356),
		Type:     easygo.NewInt32(1),
	})
	players = append(players, &share_message.PlayerOperate{
		Account:  easygo.NewString("lm7081fdfd"),
		PlayerId: easygo.NewInt64(1887567357),
		Type:     easygo.NewInt32(1),
	})
	players = append(players, &share_message.PlayerOperate{
		Account:  easygo.NewString("lm7081fdfe"),
		PlayerId: easygo.NewInt64(1887567358),
		Type:     easygo.NewInt32(1),
	})
	players = append(players, &share_message.PlayerOperate{
		Account:  easygo.NewString("lm7081fdff"),
		PlayerId: easygo.NewInt64(1887567359),
		Type:     easygo.NewInt32(1),
	})
	return players
}
