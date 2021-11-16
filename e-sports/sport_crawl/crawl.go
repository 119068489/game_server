//电竞爬虫
package sport_crawl

import (
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/astaxie/beego/logs"

	"github.com/PuerkitoBio/goquery"
	"github.com/akqp2019/mgo/bson"
)

//=======================================
//爬取项目名称
const (
	ARTICLE_FNSCORE = "article_fnscore" //蜂鸟电竞资讯
	ARTICLE_TVBCP   = "article_tvbcp"   //鲨鱼比分资讯
	VIDEO_WANPLUS   = "video_wanplus"   //玩加电竞视频
	VIDEO_CHAOFAN   = "video_chaofan"   //超凡电竞视频
	NEWS_CHAOFAN    = "news_chaofan"    //超凡电竞新闻
	NEWS_QQ         = "news_qq"         //腾讯网资讯
	NEWS_YXRB       = "news_yxrb"       //游戏日报资讯
	NEWS_SINA       = "news_sina"       //新浪电竞资讯
)

type CrawlData struct {
	Title   string `json:"Title"`   //标题
	Class   string `json:"Class"`   //游戏类型
	Content string `json:"Content"` //新闻内容
	Img     string `json:"Img"`     //图片地址
	Time    int64  `json:"Time"`    //采集时间
	Video   string `json:"Video"`   //视频地址
	Type    int    //1-新闻 2-视频
	Source  string //采集来源
}

//QQ新闻列表结构
type QQnetListResult struct {
	Data []QQnetListData `json:"data"` //列表数据
}
type QQnetListData struct {
	Title   string       `json:"title"`    //标题
	ExtData ExtDataStrut `json:"ext_data"` //扩展数据
}
type ExtDataStrut struct {
	OmUrl string `json:"om_url"` //url
}

//GET请求
func HttpGet(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		logs.Error(err)
		return nil
	}
	result, err1 := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err1 != nil {
		logs.Error(err1)
		return nil
	}
	return result
}

//获取将要爬取的html文档信息
func GetHtmlDoc(url string) *goquery.Document {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		logs.Error(err)
		return nil
	}
	return doc
}

//上传视频资源到存储桶 class站点
func UploadVideo(class, id, url string) (videourl string, imgUrl string) {
	imgPathName := fmt.Sprintf(class+"_%s_0.jpg", id)
	filePathName := fmt.Sprintf(class+"_%s.mp4", id)
	pathFileName := "sportcrawl/" + filePathName
	li := QQbucket.ObjectPutRemote(pathFileName, url)
	if li != "" {
		videourl = li
		imgUrl = "https://im-resource-1253887233.file.myqcloud.com/sportcrawl/" + imgPathName
	}
	return videourl, imgUrl
}

//蜂鸟电竞资讯爬取
func GetFnCrawlData(url string) *CrawlData {
	doc := GetHtmlDoc(url)
	if doc == nil {
		return nil
	}

	data := &CrawlData{}
	doc.Find(".body .article-detail .container .content").Each(func(i int, s *goquery.Selection) {
		s.Find(".main-body a").Remove()
		questionTitle := s.Find(".title").Text()
		questionContent, _ := s.Find(".main-body").Html()
		questionClass := s.Find(".desc p:first-child").Text()
		questionImg, _ := s.Find(".cover-img img").Attr("src")

		data.Title = questionTitle
		data.Content = questionContent
		data.Class = questionClass
		data.Time = easygo.NowTimestamp()
		data.Img = questionImg
		data.Type = 1
		data.Source = "蜂鸟电竞"
	})

	if data.Title == "" {
		return nil
	}

	return data
}

//蜂鸟电竞资讯列表爬取 class站点
func GetFnCrawlList(class string) {
	job := ReadCrawlJob(class)
	start := int64(600)
	end := int64(0)
	if job.GetValue() != "" {
		start = easygo.AtoInt64(job.GetValue())
	}
fnloop:
	for i := 1; i <= 10; i++ {
		start += int64(i)
		url := fmt.Sprintf("https://www.fnscore.com/information/%d.html", start)
		data := GetFnCrawlData(url)
		if data != nil {
			SaveCrawlDataToDB(data)
			SaveCrawlJob(class, easygo.AnytoA(start))
			end = start
		}
	}

	if end+5 >= start {
		goto fnloop
	}
}

//新浪电竞资讯爬取
func GetSinaCrawlData(id string) *CrawlData {
	url := fmt.Sprintf("https://dj.sina.com.cn/article/%s.shtml", id)
	doc := GetHtmlDoc(url)
	if doc == nil {
		return nil
	}

	data := &CrawlData{}
	doc.Find(".acticle_main").Each(func(i int, s *goquery.Selection) {
		questionTitle := s.Find(".acticle_top h1").Text()
		questionContent, _ := s.Find(".acticle_body div").Html()
		data.Title = questionTitle
		data.Content = questionContent
		data.Time = easygo.NowTimestamp()
		data.Type = 1
		data.Source = "新浪电竞"
	})

	if data.Title == "" {
		return nil
	}

	return data
}

//新浪电竞资讯列表爬取
func GetSinaCrawlNewsList(jobid string) {
	vurl := "https://dj.sina.com.cn/information"
	doc := GetHtmlDoc(vurl)
	if doc == nil {
		return
	}

	var jobstart string
	job := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_CRAWL_JOB, bson.M{"_id": jobid})
	if job != nil {
		jobstart = job.(bson.M)["Value"].(string)
	}

	doc.Find(".E_sports_news_list").Each(func(i int, s *goquery.Selection) {
		s.Find(".main_list_onepic").EachWithBreak(func(j int, d *goquery.Selection) bool {
			vurl, _ := d.Find("h3 a").Attr("href")
			filenameWithSuffix := path.Base(vurl)
			fileSuffix := path.Ext(filenameWithSuffix)
			filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)

			if jobstart == filenameOnly {
				return false
			}

			data := GetSinaCrawlData(filenameOnly)
			if data == nil && j == 0 {
				return false
			}
			if data != nil {
				SaveCrawlDataToDB(data)
				if j == 0 {
					SaveCrawlJob(jobid, filenameOnly)
				}
			}
			return true
		})
	})
}

//鲨鱼比分资讯爬取
func GetSyCrawlData(url string) *CrawlData {
	doc := GetHtmlDoc(url)
	if doc == nil {
		return nil
	}

	data := &CrawlData{}
	doc.Find(".news-description-box > .clearfix-row").Each(func(i int, s *goquery.Selection) {
		questionTitle := s.Find("h1").Text()
		questionContent, _ := s.Find(".description").Html()
		data.Title = questionTitle
		data.Content = questionContent
		data.Time = easygo.NowTimestamp()
		data.Type = 1
		data.Source = "鲨鱼比分"
	})

	if data.Title == "" {
		return nil
	}

	return data
}

//玩加电竞视频爬取
func GetWanplusVideoData(id, game string, isUpload ...bool) *CrawlData {
	url := fmt.Sprintf("http://www.wanplus.com/%s/video/%s", game, id)
	doc := GetHtmlDoc(url)
	if doc == nil {
		return nil
	}

	data := &CrawlData{}
	doc.Find("body .body-inner").Each(func(i int, s *goquery.Selection) {
		questionTitle := s.Find(".banner .user-text #shareTitle").Text()
		questionvideo, _ := s.Find(".content .ov .video-player video").Attr("src")
		questionClass := s.Find(".banner .user-text .user-n .user-name").Text()
		questionImg, _ := s.Find(".content .ov .video-player video").Attr("poster")

		data.Title = questionTitle
		data.Video = questionvideo
		data.Class = questionClass
		data.Time = easygo.NowTimestamp()
		data.Img = questionImg
		data.Type = 2
		data.Source = "玩加电竞"
	})

	if data.Title == "" || data.Video == "" {
		return nil
	}

	isUpload = append(isUpload, false)
	if isUpload[0] {
		vUrl, iUrl := UploadVideo(VIDEO_WANPLUS, id, data.Video)
		data.Video = vUrl
		data.Img = iUrl
	}

	return data
}

//玩加电竞视频列表爬取
func GetWanplusVideoList(jobid string) {
	var gameList []string
	i := strings.LastIndex(jobid, "_")
	if jobid[:i] != VIDEO_WANPLUS {
		gameList = append(gameList, "lol", "kog", "csgo")
	} else {
		gameList = append(gameList, jobid[i+1:])
	}

	for _, game := range gameList {
		jobid = VIDEO_WANPLUS
		vurl := fmt.Sprintf("http://www.wanplus.com/%s/video", game)
		doc := GetHtmlDoc(vurl)
		if doc == nil {
			return
		}

		jobid += "_" + game //重装jobid
		var jobstart string
		job := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_CRAWL_JOB, bson.M{"_id": jobid})
		if job != nil {
			jobstart = job.(bson.M)["Value"].(string)
		}

		doc.Find("#video-ul").Each(func(i int, s *goquery.Selection) {
			s.Find("li").EachWithBreak(func(j int, d *goquery.Selection) bool {
				vurl, _ := d.Find(".v-info a").Attr("href")
				filenameWithSuffix := path.Base(vurl)
				fileSuffix := path.Ext(filenameWithSuffix)
				filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
				if jobstart == filenameOnly {
					return false
				}

				data := GetWanplusVideoData(filenameOnly, game, true)
				if data == nil && j == 0 {
					return false
				}
				if data != nil {
					SaveCrawlDataToDB(data)
					if j == 0 {
						SaveCrawlJob(jobid, filenameOnly)
					}
				}
				return true
			})
		})
	}
}

//超凡电竞视频爬取
func GetChaofanVideoData(id, game, img, title, jobid string, isUpload ...bool) *CrawlData {
	url := "https://www.chaofan.com/video/player?id=" + id
	doc := GetHtmlDoc(url)
	if doc == nil {
		return nil
	}
	data := &CrawlData{}
	questionvideo, _ := doc.Find("#video-box").Attr("data-url")
	data.Title = title
	data.Video = questionvideo
	data.Class = game
	data.Time = easygo.NowTimestamp()
	data.Img = img
	data.Type = 2
	data.Source = "超凡电竞"

	if data.Title == "" || data.Video == "" {
		return nil
	}

	isUpload = append(isUpload, false)
	if isUpload[0] {
		vUrl, iUrl := UploadVideo(VIDEO_WANPLUS, id, data.Video)
		data.Video = vUrl
		data.Img = iUrl
	}

	return data
}

//超凡电竞视频爬取列表
func GetChaoFanVideoList(jobid string) {
	var gameList []string
	i := strings.LastIndex(jobid, "_")
	if jobid[:i] != VIDEO_CHAOFAN {
		gameList = append(gameList, "lol", "kog", "dota2", "csgo")
	} else {
		gameList = append(gameList, jobid[i+1:])
	}

	for _, game := range gameList {
		jobid = VIDEO_CHAOFAN
		vurl := "https://www.chaofan.com/video/" + game
		doc := GetHtmlDoc(vurl)
		if doc == nil {
			return
		}

		jobid += "_" + game //重装jobid
		var jobstart string
		job := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_CRAWL_JOB, bson.M{"_id": jobid})
		if job != nil {
			jobstart = job.(bson.M)["Value"].(string)
		}

		doc.Find("body .content-box .hot-data-box .left-box").Each(func(i int, s *goquery.Selection) {
			s.Find(".video-list li").EachWithBreak(func(j int, d *goquery.Selection) bool {
				vurl, _ := d.Find(".link").Attr("href")
				cP, _ := url.Parse(vurl)
				if cP.Scheme == "" {
					cP.Scheme = "https"
				}
				img, _ := d.Find(".link .img-view").Attr("style")
				strs := strings.SplitN(img, "?", 2)
				if len(strs) == 2 {
					strs = strings.SplitN(strs[0], "('", 2)
				}
				if len(strs) == 2 {
					img = strs[1]
				}
				filenameWithSuffix := path.Base(cP.String())
				fileSuffix := path.Ext(filenameWithSuffix)
				filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
				title := d.Find(".title").Text()
				if jobstart == filenameOnly {
					return false
				}

				data := GetChaofanVideoData(filenameOnly, game, img, title, jobid, true)
				if data == nil && j == 0 {
					return false
				}
				if data != nil {
					SaveCrawlDataToDB(data)
					if j == 0 {
						SaveCrawlJob(jobid, filenameOnly)
					}
				}
				return true
			})
		})
	}
}

//超凡电竞新闻爬取
func GetChaofanNewsData(id string, isUpload ...bool) *CrawlData {
	url := fmt.Sprintf("https://www.chaofan.com/news/%s.html", id)
	doc := GetHtmlDoc(url)
	if doc == nil {
		return nil
	}
	data := &CrawlData{}
	questionTitle := doc.Find(".news-layout-center h1").Text()
	questionContent, _ := doc.Find(".news-layout-center .description").Html()

	data.Title = questionTitle
	data.Content = questionContent
	data.Time = easygo.NowTimestamp()
	data.Type = 1
	data.Source = "超凡电竞"

	if data.Title == "" {
		return nil
	}

	return data
}

//超凡电竞新闻爬取列表
func GetChaoFanNewsList(jobid string) {
	var gameList []string

	i := strings.LastIndex(jobid, "_")
	if jobid[:i] != NEWS_CHAOFAN {
		gameList = append(gameList, "lol", "kog", "dota2", "csgo", "pubg", "underlords", "tft")
	} else {
		gameList = append(gameList, jobid[i+1:])
	}

	for _, game := range gameList {
		jobid = NEWS_CHAOFAN
		vurl := "https://www.chaofan.com/news/" + game
		doc := GetHtmlDoc(vurl)
		if doc == nil {
			return
		}
		jobid += "_" + game //重装jobid
		var jobstart string
		job := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_CRAWL_JOB, bson.M{"_id": jobid})
		if job != nil {
			jobstart = job.(bson.M)["Value"].(string)
		}

		doc.Find("body .content-box .list-box").Each(func(i int, s *goquery.Selection) {
			s.Find(".clearfix-row").EachWithBreak(func(j int, d *goquery.Selection) bool {
				vurl, _ := d.Find("a").Attr("href")
				filenameWithSuffix := path.Base(vurl)
				fileSuffix := path.Ext(filenameWithSuffix)
				filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)

				if jobstart == filenameOnly {
					return false
				}

				data := GetChaofanNewsData(filenameOnly, true)
				if data == nil && j == 0 {
					return false
				}
				if data != nil {
					SaveCrawlDataToDB(data)
					if j == 0 {
						SaveCrawlJob(jobid, filenameOnly)
					}
				}
				return true
			})
		})
	}
}

//腾讯电竞新闻爬取
func GetQQNewsData(id, title string) *CrawlData {
	url := "https://page.om.qq.com/page/" + id
	doc := GetHtmlDoc(url)
	if doc == nil {
		return nil
	}
	data := &CrawlData{}
	questionTitle := doc.Find("#content .header .title").Text()
	questionContent, _ := doc.Find("#content .article").Html()

	data.Title = questionTitle
	data.Content = questionContent
	data.Time = easygo.NowTimestamp()
	data.Type = 1
	data.Source = "腾讯网"

	if data.Title == "" {
		return nil
	}

	return data
}

//腾讯电竞列表爬取 class站点
func QQnetList(class string) {
	var authors []string //作者id
	i := strings.LastIndex(class, "_")
	if class[:i] != NEWS_QQ {
		authors = append(authors, "6102318", "15446343")
	} else {
		authors = append(authors, class[i+1:])
	}

	for _, a := range authors {
		jobid := NEWS_QQ
		url := fmt.Sprintf("https://pacaio.match.qq.com/om/mediaArticles?mid=%s&num=10", a)

		result := HttpGet(url)
		one := QQnetListResult{}
		err := json.Unmarshal(result, &one)
		if err != nil {
			logs.Error(err)
			continue
		}

		var jobstart string
		jobid += "_" + a //重装jobid
		job := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_CRAWL_JOB, bson.M{"_id": jobid})
		if job != nil {
			jobstart = job.(bson.M)["Value"].(string)
		}

		for i, u := range one.Data {
			var id string
			strs := strings.SplitN(u.ExtData.OmUrl, "page/", 2)
			if len(strs) == 2 {
				id = strs[1]
			}

			if jobstart == id {
				break
			}

			data := GetQQNewsData(id, u.Title)
			if data == nil && i == 0 {
				break
			}
			if data != nil {
				SaveCrawlDataToDB(data)
				if i == 0 {
					SaveCrawlJob(jobid, id)
				}
			}

		}
	}
}

//游戏日报新闻爬取
func GetYxrbNewsData(url string) *CrawlData {
	doc := GetHtmlDoc(url)
	if doc == nil {
		return nil
	}
	data := &CrawlData{}
	questionTitle := doc.Find(".article-title p").Text()
	questionContent, _ := doc.Find("article").Html()

	data.Title = questionTitle
	data.Content = questionContent
	data.Time = easygo.NowTimestamp()
	data.Type = 1
	data.Source = "游戏日报"

	if data.Title == "" {
		return nil
	}

	return data
}

//游戏日报新闻列表爬取
func GetYxrbNewsList(id string) {
	doc := GetHtmlDoc("http://news.yxrb.net/aggrogame/")
	if doc == nil {
		return
	}
	var jobstart string
	job := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_CRAWL_JOB, bson.M{"_id": id})
	if job != nil {
		jobstart = job.(bson.M)["Value"].(string)
	}

	doc.Find(".content .news-list-module").Each(func(i int, s *goquery.Selection) {
		s.Find("article").EachWithBreak(func(j int, d *goquery.Selection) bool {
			vurl, _ := d.Find("a").Attr("href")
			cP, _ := url.Parse(vurl)
			filenameWithSuffix := path.Base(cP.String())
			fileSuffix := path.Ext(filenameWithSuffix)
			filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)

			if filenameOnly == jobstart {
				return false
			}

			data := GetYxrbNewsData(vurl)
			if data == nil && j == 0 {
				return false
			}
			if data != nil {
				SaveCrawlDataToDB(data)
				if j == 0 {
					SaveCrawlJob(id, filenameOnly)
				}
			}
			return true
		})
	})
}

//电竞爬虫爬取数据
func CrawlDataRun(obj ...string) {
	logs.Info(obj, "爬虫爬取数据开始")
	if len(obj) == 0 {
		obj = append(obj, ARTICLE_FNSCORE, ARTICLE_TVBCP, VIDEO_WANPLUS, VIDEO_CHAOFAN, NEWS_CHAOFAN, NEWS_QQ, NEWS_YXRB, NEWS_SINA)
	}

	// regFnscore, _ := regexp.Compile("^" + ARTICLE_FNSCORE)
	// for _, o := range obj {
	// 	job := ReadCrawlJob(o)
	// 	if regFnscore.MatchString(o) {
	// 	}
	// }

	for _, o := range obj {
		id := o
	reSwitch:
		switch id {
		case ARTICLE_FNSCORE:
			GetFnCrawlList(ARTICLE_FNSCORE)
		case ARTICLE_TVBCP:
			job := ReadCrawlJob(ARTICLE_TVBCP)
			start := int64(50197)
			end := int64(0)
			if job.GetValue() != "" {
				start = easygo.AtoInt64(job.GetValue())
			}

		tvbcploop:
			for i := 1; i <= 10; i++ {
				start += int64(i)
				url := fmt.Sprintf("http://www.tvbcp.com/detail/%d.html", start)
				data := GetSyCrawlData(url)
				if data != nil {
					SaveCrawlDataToDB(data)
					SaveCrawlJob(o, easygo.AnytoA(start))
					end = start
				}
			}

			if end+5 >= start {
				goto tvbcploop
			}
		case VIDEO_WANPLUS:
			GetWanplusVideoList(o)
		case VIDEO_CHAOFAN:
			GetChaoFanVideoList(o)
		case NEWS_CHAOFAN:
			GetChaoFanNewsList(o)
		case NEWS_QQ:
			QQnetList(o)
		case NEWS_YXRB:
			GetYxrbNewsList(NEWS_YXRB)
		case NEWS_SINA:
			GetSinaCrawlNewsList(NEWS_SINA)
		default:
			i := strings.LastIndex(id, "_")
			if i == -1 {
				continue
			} else {
				id = id[:i]
				goto reSwitch
			}
		}
	}
	logs.Info("爬虫爬取数据结束")
}

//读取爬虫进度 st爬取项目名称
func ReadCrawlJob(st string) *share_message.TableCrawlJob {
	msg := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_CRAWL_JOB, bson.M{"_id": st})
	one := &share_message.TableCrawlJob{}
	for_game.StructToOtherStruct(msg, one)
	return one
}

//保存爬取进度 st爬取项目名称，v进度值
func SaveCrawlJob(st string, v interface{}) {
	for_game.FindAndModify(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_CRAWL_JOB, bson.M{"_id": st}, bson.M{"$set": bson.M{"Value": v, "Time": easygo.NowTimestamp()}}, true)
}

//保存爬取数据到数据库
func SaveCrawlDataToDB(data ...*CrawlData) {
	if len(data) == 0 {
		return
	}
	var atcs []interface{}
	var vid []interface{}
	for _, d := range data {
		switch d.Type {
		case 1:
			one := &share_message.TableESPortsRealTimeInfo{
				Id:           easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_NEWS_SOURCE)),
				CreateTime:   easygo.NewInt64(easygo.NowTimestamp()),
				Status:       easygo.NewInt32(for_game.ESPORTS_NEWS_STATUS_0),
				IssueTime:    easygo.NewInt64(d.Time),
				Title:        easygo.NewString(d.Title),
				Content:      easygo.NewString(d.Content),
				DataSource:   easygo.NewString(d.Source),
				AppLabelName: easygo.NewString(d.Class),
			}
			result := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_NEWS_SOURCE, bson.M{"Title": one.GetTitle()})
			if result == nil {
				atcs = append(atcs, one)
			}
		case 2:
			one := &share_message.TableESPortsVideoInfo{
				Id:            easygo.NewInt64(for_game.NextId(for_game.TABLE_ESPORTS_VIDEO_SOURCE)),
				CreateTime:    easygo.NewInt64(easygo.NowTimestamp()),
				Status:        easygo.NewInt32(for_game.ESPORTS_NEWS_STATUS_0),
				IssueTime:     easygo.NewInt64(d.Time),
				Title:         easygo.NewString(d.Title),
				VideoUrl:      easygo.NewString(d.Video),
				DataSource:    easygo.NewString(d.Source),
				AppLabelName:  easygo.NewString(d.Class),
				VideoType:     easygo.NewInt64(for_game.ESPORTS_VIDEO_TYPE_1),
				CoverImageUrl: easygo.NewString(d.Img),
			}
			result := for_game.FindOne(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_VIDEO_SOURCE, bson.M{"Title": one.GetTitle()})
			if result == nil {
				vid = append(vid, one)
			}
		}

	}

	if len(atcs) > 0 {
		for_game.InsertAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_NEWS_SOURCE, atcs...)
	}

	if len(vid) > 0 {
		for_game.InsertAllMgo(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_VIDEO_SOURCE, vid...)
	}
}

//比赛历史数据
func CrawlScoreHistoryDataRun(obj ...int64) {
	logs.Info("爬虫爬取比赛历史数据开始")

	gameList := make([]int64, 0)
	if len(obj) == 0 {
		list, _ := for_game.FindAll(for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_GAME, bson.M{"game_status": bson.M{"$ne": for_game.GAME_STATUS_2}, "history_id": bson.M{"$gt": 0}}, 0, 0)
		if len(list) == 0 {
			return
		}
		for _, li := range list {
			gameList = append(gameList, li.(bson.M)["history_id"].(int64))
		}
	} else {
		gameList = append(gameList, obj...)
	}

	saveDate := []for_game.RecentData{}
	for _, id := range gameList {
		url := fmt.Sprintf("https://img1.famulei.com/match_pre/%v.json", id)
		result := HttpGet(url)
		one := for_game.ResultMsg{}
		err := json.Unmarshal(result, &one)
		if err != nil {
			logs.Error(err)
			continue
		}
		one.Data.Id = id
		saveDate = append(saveDate, one.Data)
	}

	if len(saveDate) > 0 {
		var data []interface{}
		for _, v := range saveDate {
			b1 := bson.M{"_id": v.Id}
			b2 := v
			data = append(data, b1, b2)
		}
		for_game.UpsertAll(easygo.MongoMgr, for_game.MONGODB_NINGMENG, for_game.TABLE_ESPORTS_HISTORY, data)
	}

	logs.Info("爬虫爬取比赛历史数据结束")
}
