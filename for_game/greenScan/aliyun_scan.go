package greenScan

import (
	"encoding/json"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/for_game/greenScan/common"
	"github.com/astaxie/beego/logs"
)

const (
	ACCESS_KEY_ID                = "LTAI4FmFpwJyTiBWWJQM3at9"
	ACCESS_KEY_SECRET            = "inWptQcWaLhSpbR72XF3tXLgjvwIDP"
	IMAGE_ASYNC_SCAN_PATH        = "/green/image/asyncscan" //批量异步无回调的图片发送地址,需要业务轮询调用查结果
	VIDEO_ASYNC_SCAN_PATH        = "/green/video/asyncscan" //批量异步无回调的视频发送地址,需要业务轮询调用查结果
	RESULT_IMAGE_ASYNC_SCAN_PATH = "/green/image/results"   //批量异步发送图片结果查询地址
	RESULT_VIDEO_ASYNC_SCAN_PATH = "/green/video/results"   //批量异步发送视频结果查询地址
	TEXT_SCAN_PATH               = "/green/text/scan"       //同步文本实时调用发送地址
	IMAGE_SCAN_PATH              = "/green/image/scan"      //同步图片实时调用发送地址
	//TODO:视频同步调用地址目前没有封装
	SCAN_IP = "127.0.0.1"
)

//========解析阿里所需要的实体定义========开始
//返回定义
type AliRsp struct {
	Code      int32     `json:"code"`      //公用:错误码，和HTTP的status code一致。2xx：表示成功。 4xx：表示请求有误。
	Msg       string    `json:"msg"`       //公用:错误描述信息。
	RequestId string    `json:"requestId"` //公用:调用请求id。
	Data      []AliData `json:"data"`      //公用:阿里返回的内容
}

type AliData struct {
	Code        int32       `json:"code"`    //公用:错误码2xx：表示成功。 4xx：表示请求有误。5xx：表示后端有误。	(200请求成功。280任务正在执行中参照阿里云上文档)
	TextContent string      `json:"content"` //注意：文本解析的时候用到的特定域:被检测文本，和调用请求中的待检测文本对应。
	DataId      string      `json:"dataId"`  //公用:检测对象对应的数据ID。说明 如果在请求参数中传入了dataId，则此处返回对应的dataId
	Msg         string      `json:"msg"`     //公用:错误描述信息。
	TaskId      string      `json:"taskId"`  //公用:本次检测任务的ID。
	Results     []AliResult `json:"results"` //公用:返回结果。调用成功时（code=200），返回结果中包含一个或多个元素。每个元素是个结构体，具体结构描述请参见result
	ImageUrl    string      `json:"url"`     //注意：图片解析的时候用到的特定域:图片URL地址
}
type AliResult struct {
	Label      string          `json:"label"`      //公用:检测结果的分类:参照文本和图片,视频的结果
	Rate       float64         `json:"rate"`       //公用:结果属于当前分类的概率，取值范围：0.00~100.00。值越高，表示越有可能属于当前分类。
	Scene      string          `json:"scene"`      //公用:检测场景，和调用请求中的场景对应。
	Suggestion string          `json:"suggestion"` //公用:建议您执行的后续操作
	Details    []AliTextDetail `json:"details"`    //注意：文本解析的时候用到的特定域:文本的一些详细内容,用于替换命中词,取得真正的命中词
}

//注意：文本解析的时候用到的特定域:
type AliTextDetail struct {
	Label string `json:"label"` //文本命中风险的分类。取值：
	//文本命中风险的分类。取值：
	//spam：含垃圾信息
	//ad：广告
	//politics：涉政
	//terrorism：暴恐
	//abuse：辱骂
	//porn：色情
	//flood：灌水
	//contraband：违禁
	//meaningless：无意义
	//customized：自定义（例如命中自定义关键词）
	Contexts []AliTextContext `json:"contexts"` //命中该风险的上下文信息。具体结构描述请参见context。
}

type AliTextContext struct {
	Context string `json:"context"` //检测文本命中的风险内容的上下文信息。如果命中了您自定义的风险文本库，则会返回命中的文本内容（关键词或相似文本）。
}

//================以下是阿里云转腾讯结果的中间结构========================
//=====解析阿里TEXT文本 阿里到腾讯转文本的值的中间实体======
//1、腾讯的TEXT文本EvilType 恶意类型 100：正常 20001：政治 20002：色情 20006：涉毒违法 20007：谩骂 20105：广告引流 24001：暴恐
// 2、腾讯的TEXT文本EvilFlag是否恶意 0：正常 1：可疑
// 3、腾讯的TEXT文本EvilLabel恶意标签，Normal：正常，Polity：涉政，Porn：色情，Illegal：违法，Abuse：谩骂，Terror：暴恐，Ad：广告，Custom：自定义关键词
type AliToTenTextBody struct {
	EvilFlag  int32
	EvilLabel string
	EvilType  int32
}

//解析阿里TEXT文本 通过阿里的label对应到腾讯的label和type
// 阿里的label normal：正常文本 spam：含垃圾信息 ad：广告 politics：涉政 terrorism：暴恐 abuse：
// 辱骂 porn：色情 flood：灌水 contraband：违禁 meaningless：无意义 customized：自定义（例如命中自定义关键词）
//注意阿里的spam：含垃圾信息和meaningless：无意义   对应到腾讯的正常

//1、腾讯的TEXT文本EvilType 恶意类型 100：正常 20001：政治 20002：色情 20006：涉毒违法 20007：谩骂 20105：广告引流 24001：暴恐
// 2、腾讯的TEXT文本EvilFlag是否恶意 0：正常 1：可疑
// 3、腾讯的TEXT文本EvilLabel恶意标签，Normal：正常，Polity：涉政，Porn：色情，Illegal：违法，Abuse：谩骂，Terror：暴恐，Ad：广告，Custom：自定义关键词
var AliToTenTextBodyMap = map[string]*AliToTenTextBody{
	"normal":      &AliToTenTextBody{EvilFlag: 0, EvilLabel: "Normal", EvilType: 100},
	"ad":          &AliToTenTextBody{EvilFlag: 1, EvilLabel: "Ad", EvilType: 20105},
	"politics":    &AliToTenTextBody{EvilFlag: 1, EvilLabel: "Polity", EvilType: 20001},
	"terrorism":   &AliToTenTextBody{EvilFlag: 1, EvilLabel: "Terror", EvilType: 24001},
	"abuse":       &AliToTenTextBody{EvilFlag: 1, EvilLabel: "Abuse", EvilType: 20007},
	"porn":        &AliToTenTextBody{EvilFlag: 1, EvilLabel: "Porn", EvilType: 20002},
	"contraband":  &AliToTenTextBody{EvilFlag: 1, EvilLabel: "Illegal", EvilType: 20006},
	"spam":        &AliToTenTextBody{EvilFlag: 0, EvilLabel: "Normal", EvilType: 100},
	"meaningless": &AliToTenTextBody{EvilFlag: 0, EvilLabel: "Normal", EvilType: 100},
}

//解析阿里TEXT文本 阿里的suggestion对应腾讯的suggestion
//=====阿里的suggestion pass：文本正常 review：需要人工审核 block：文本违规，可以直接删除或者做限制处理
//=====腾讯的suggestion 建议值,Block：打击,Review：待复审,Normal：正常
var AliToTenTextSuggestionMap = map[string]string{
	"pass":   "Normal",
	"review": "Review",
	"block":  "Block",
}

//=====解析阿里Image图片 阿里到腾讯转图片的值的中间实体======
//1、腾讯的Image图片EvilType 恶意类型 100：正常 20001：政治 20002：色情 20006：涉毒违法 20007：谩骂 20103：性感 24001：暴恐
// 2、腾讯的Image图片EvilFlag是否恶意 0：正常 1：可疑
type AliToTenImageBody struct {
	EvilFlag int32
	EvilType int32
}

//解析阿里Image图片 通过阿里的label对应到腾讯的EvilFlag和EvilType
// 阿里的label 具体参考阿里文档

//1、腾讯的Image图片EvilType 恶意类型 100：正常 20001：政治 20002：色情 20006：涉毒违法 20007：谩骂 20103：性感 24001：暴恐
// 2、腾讯的Image图片EvilFlag是否恶意 0：正常 1：可疑
var AliToTenImageBodyMap = map[string]*AliToTenImageBody{
	"normal":      &AliToTenImageBody{EvilFlag: 0, EvilType: 100},
	"sexy":        &AliToTenImageBody{EvilFlag: 1, EvilType: 20103},
	"porn":        &AliToTenImageBody{EvilFlag: 1, EvilType: 20002},
	"bloody":      &AliToTenImageBody{EvilFlag: 1, EvilType: 24001},
	"explosion":   &AliToTenImageBody{EvilFlag: 1, EvilType: 24001},
	"outfit":      &AliToTenImageBody{EvilFlag: 1, EvilType: 24001},
	"logo":        &AliToTenImageBody{EvilFlag: 1, EvilType: 24001},
	"weapon":      &AliToTenImageBody{EvilFlag: 1, EvilType: 24001},
	"politics":    &AliToTenImageBody{EvilFlag: 1, EvilType: 20001},
	"violence":    &AliToTenImageBody{EvilFlag: 1, EvilType: 24001},
	"crowd":       &AliToTenImageBody{EvilFlag: 1, EvilType: 24001},
	"parade":      &AliToTenImageBody{EvilFlag: 1, EvilType: 24001},
	"carcrash":    &AliToTenImageBody{EvilFlag: 1, EvilType: 24001},
	"flag":        &AliToTenImageBody{EvilFlag: 1, EvilType: 24001},
	"location":    &AliToTenImageBody{EvilFlag: 1, EvilType: 24001},
	"others":      &AliToTenImageBody{EvilFlag: 0, EvilType: 100},
	"abuse":       &AliToTenImageBody{EvilFlag: 1, EvilType: 20007},
	"terrorism":   &AliToTenImageBody{EvilFlag: 1, EvilType: 24001},
	"contraband":  &AliToTenImageBody{EvilFlag: 1, EvilType: 20006},
	"spam":        &AliToTenImageBody{EvilFlag: 0, EvilType: 100},
	"npx":         &AliToTenImageBody{EvilFlag: 0, EvilType: 100},
	"qrcode":      &AliToTenImageBody{EvilFlag: 0, EvilType: 100},
	"programCode": &AliToTenImageBody{EvilFlag: 0, EvilType: 100},
	"ad":          &AliToTenImageBody{EvilFlag: 0, EvilType: 100},
}

//=============解析阿里TEXT文本和图片所需要的实体定义========结束

//图片验证维度
var imageScenes []string = []string{"porn", "terrorism", "ad"}

//视频验证维度
var videoScenes []string = []string{"porn", "terrorism", "live"}

//文本验证维度
var textScenes []string = []string{"antispam"}

//阿里异步通过taskIds去取得图片验证的结果
//请求参数：taskIds：异步发送的时候取得的任务列表
//返回结果说明:
//int32 0 图片审核未通过  1审核通过 2审核处理中或网络请求出错,需下次验证
//string 错误码code记录Log用 ,string 封装后的错误内容存储数据库用
func GetImageRstByTaskIds(taskIds []string) (int32, string, string) {

	profile := common.Profile{AccessKeyId: ACCESS_KEY_ID, AccessKeySecret: ACCESS_KEY_SECRET}

	path := RESULT_IMAGE_ASYNC_SCAN_PATH

	clientInfo := common.ClinetInfo{Ip: SCAN_IP}

	var client common.AliYunClient = common.ScanClient{Profile: profile}

	if nil == client {
		return 2, "", ""
	}
	var rstClientRes = ""

	func() {
		rstClientRes = client.GetRstResponse(path, clientInfo, taskIds)
	}()

	if "" == rstClientRes {
		return 2, "", ""
	}
	result := AliRsp{}
	err := json.Unmarshal([]byte(rstClientRes), &result)

	if nil == err {
		if result.Code == 200 {

			//任务的循环
			for i := 0; i < len(taskIds); i++ {
				if len(result.Data) > 0 && result.Data[i].Code == 200 {
					//循环的次数是发请求scenes := []string{"porn","terrorism","ad"}的个数
					//场景的个数
					for k := 0; k < len(imageScenes); k++ {
						var label string = ""
						if len(result.Data[i].Results) > 0 {
							label = result.Data[i].Results[k].Label
						}
						if "" == label {
							return 2, "", ""
						} else if "normal" != label {
							return 0, "10000", "图片审核失败-label:" + label + ";-taskid:" + taskIds[i]
						} else {
							continue
						}
					}
				} else if len(result.Data) > 0 && result.Data[i].Code == 480 {

					return 0, "480", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 400 {

					return 0, "400", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 401 {

					return 0, "401", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 403 {

					return 0, "403", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 404 {

					return 0, "404", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 586 {

					return 0, "586", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 587 {

					return 0, "587", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 588 {

					return 0, "588", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 589 {

					return 0, "589", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 590 {

					return 0, "590", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 591 {

					return 0, "591", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 592 {

					return 0, "592", "网络出错"

				} else {
					//除了以上状态码其余的需要重试
					return 2, "", ""
				}
			}
		} else if result.Code == 400 {
			return 0, "400", "网络出错"
		} else if result.Code == 401 {
			return 0, "401", "网络出错"
		} else if result.Code == 403 {
			return 0, "403", "网络出错"
		} else if result.Code == 404 {
			return 0, "404", "网络出错"
		} else if result.Code == 480 {
			return 0, "480", "网络出错"
		} else {
			return 2, "", ""
		}
		return 1, "", ""
	} else {
		logs.Error(err)
		return 2, "", ""
	}
	return 1, "", ""
}

//阿里异步通过taskIds去取得视频验证的结果
//请求参数：taskIds：异步发送的时候取得的任务列表
//返回结果说明:
// int32 0 视频审核未通过  1审核通过 2审核处理中或网络请求出错,需下次验证
// string 错误码code记录Log用 ,string 封装后的错误内容存储数据库用
func GetVideoRstByTaskIds(taskIds []string) (int32, string, string) {

	profile := common.Profile{AccessKeyId: ACCESS_KEY_ID, AccessKeySecret: ACCESS_KEY_SECRET}

	path := RESULT_VIDEO_ASYNC_SCAN_PATH

	clientInfo := common.ClinetInfo{Ip: SCAN_IP}

	var client common.AliYunClient = common.ScanClient{Profile: profile}
	if nil == client {
		return 2, "", ""
	}
	var rstClientRes = ""
	func() {
		rstClientRes = client.GetRstResponse(path, clientInfo, taskIds)
	}()

	if "" == rstClientRes {
		return 2, "", ""
	}

	result := AliRsp{}
	err := json.Unmarshal([]byte(rstClientRes), &result)

	if nil == err {
		if result.Code == 200 {
			for i := 0; i < len(taskIds); i++ {

				if len(result.Data) > 0 && result.Data[i].Code == 200 {

					//循环的次数是发请求scenes := []string{"porn","terrorism","live"}的个数
					for k := 0; k < len(videoScenes); k++ {
						var label string = ""
						if len(result.Data[i].Results) > 0 {
							label = result.Data[i].Results[k].Label
						}
						if "" == label {
							return 2, "", ""
						} else if "normal" != label {
							return 0, "10000", "视频审核失败-label:" + label + ";-taskid:" + taskIds[i]
						} else {
							continue
						}
					}
				} else if len(result.Data) > 0 && result.Data[i].Code == 480 {

					return 0, "480", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 400 {

					return 0, "400", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 401 {

					return 0, "401", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 403 {

					return 0, "403", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 404 {

					return 0, "404", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 586 {

					return 0, "586", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 587 {

					return 0, "587", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 588 {

					return 0, "588", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 589 {

					return 0, "589", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 590 {

					return 0, "590", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 591 {

					return 0, "591", "网络出错"

				} else if len(result.Data) > 0 && result.Data[i].Code == 592 {

					return 0, "592", "网络出错"

				} else {
					return 2, "", ""
				}
			}
		} else if result.Code == 400 {
			return 0, "400", "网络出错"
		} else if result.Code == 401 {
			return 0, "401", "网络出错"
		} else if result.Code == 403 {
			return 0, "403", "网络出错"
		} else if result.Code == 404 {
			return 0, "404", "网络出错"
		} else if result.Code == 480 {
			return 0, "480", "网络出错"
		} else {
			return 2, "", "'"
		}

		return 1, "", ""

	} else {
		logs.Error(err)
		return 2, "", ""
	}
	return 1, "", ""
}

//阿里异步发送图片请求取得异步调用结果的taskIds
//请求参数:图片URL切片
//返回值说明:[]string:检测的任务id列表 int32:0请求出错,1请求成功或者需要二次重新验证
func GetImageScanTaskIds(image_urls []string) ([]string, int32) {

	rstTaksks := []string{}
	if nil == image_urls || len(image_urls) == 0 {
		return rstTaksks, 1
	}
	profile := common.Profile{AccessKeyId: ACCESS_KEY_ID, AccessKeySecret: ACCESS_KEY_SECRET}

	path := IMAGE_ASYNC_SCAN_PATH

	clientInfo := common.ClinetInfo{Ip: SCAN_IP}

	// 构造请求数据
	bizType := "Green"
	//porn图片智能鉴黄,terrorism图片暴恐涉政识别,ad图文违规识别
	scenes := imageScenes

	tasks := []common.ImageTask{}
	for i := 0; i < len(image_urls); i++ {
		task := common.ImageTask{DataId: common.Rand().Hex(), Url: image_urls[i]}
		tasks = append(tasks, task)
	}
	bizData := common.ImageBizData{bizType, scenes, tasks}

	var client common.AliYunClient = common.ScanClient{Profile: profile}
	if nil == client {
		return rstTaksks, 1
	}
	var imageClientRes = ""

	func() {
		imageClientRes = client.GetImageResponse(path, clientInfo, bizData)
	}()

	if "" == imageClientRes {
		return rstTaksks, 1
	}

	result := AliRsp{}
	err := json.Unmarshal([]byte(imageClientRes), &result)

	if nil == err {
		if result.Code == 200 {
			for i := 0; i < len(image_urls); i++ {
				if len(result.Data) > 0 && result.Data[i].Code == 200 {

					var taskId string = result.Data[i].TaskId

					//组装得到的taskIds
					if "" == taskId {
						rstTaksks = []string{}
						return rstTaksks, 1
					} else {
						rstTaksks = append(rstTaksks, taskId)
					}
				} else if len(result.Data) > 0 && result.Data[i].Code == 400 {

					return rstTaksks, 0

				} else if len(result.Data) > 0 && result.Data[i].Code == 401 {

					return rstTaksks, 0

				} else if len(result.Data) > 0 && result.Data[i].Code == 403 {

					return rstTaksks, 0

				} else if len(result.Data) > 0 && result.Data[i].Code == 480 {

					return rstTaksks, 0

				} else {
					rstTaksks = []string{}
					return rstTaksks, 1
				}
			}
		} else if result.Code == 400 {

			return rstTaksks, 0

		} else if result.Code == 401 {

			return rstTaksks, 0

		} else if result.Code == 403 {

			return rstTaksks, 0

		} else if result.Code == 404 {

			return rstTaksks, 0

		} else if result.Code == 480 {

			return rstTaksks, 0

		} else {
			//需要重试
			return rstTaksks, 1
		}
	} else {
		//需要重试
		logs.Error(err)
		return rstTaksks, 1
	}
	return rstTaksks, 1
}

//阿里异步发送视频请求取得异步调用结果的taskIds
//请求参数:视频URL切片
//返回值说明:[]string:检测的任务id列表 int32:0请求出错,1请求成功或者需要二次重新验证
func GetVideoScanTaskIds(video_urls []string) ([]string, int32) {
	rstTaksks := []string{}
	if nil == video_urls || len(video_urls) == 0 {
		return rstTaksks, 1
	}
	profile := common.Profile{AccessKeyId: ACCESS_KEY_ID, AccessKeySecret: ACCESS_KEY_SECRET}

	path := VIDEO_ASYNC_SCAN_PATH

	clientInfo := common.ClinetInfo{Ip: SCAN_IP}

	// 构造请求数据
	bizType := "Green"
	//"porn"识别短视频是否为色情视频。,"terrorism"识别短视频是否为暴恐涉政视频。,"live"识别短视频中的不良场景。
	scenes := videoScenes

	tasks := []common.VideoTask{}
	for i := 0; i < len(video_urls); i++ {
		task := common.VideoTask{DataId: common.Rand().Hex(), Url: video_urls[i]}
		tasks = append(tasks, task)
	}

	bizData := common.VideoBizData{bizType, scenes, tasks}

	var client common.AliYunClient = common.ScanClient{Profile: profile}

	if nil == client {
		return rstTaksks, 1
	}

	var vidoeClientRes = ""
	func() {
		vidoeClientRes = client.GetVideoResponse(path, clientInfo, bizData)
	}()

	if "" == vidoeClientRes {
		return rstTaksks, 1
	}
	result := AliRsp{}
	err := json.Unmarshal([]byte(vidoeClientRes), &result)

	if nil == err {
		if result.Code == 200 {
			for i := 0; i < len(video_urls); i++ {
				if len(result.Data) > 0 && result.Data[i].Code == 200 {
					var taskId string = result.Data[i].TaskId

					if "" == taskId {
						rstTaksks = []string{}
						return rstTaksks, 1
					} else {
						rstTaksks = append(rstTaksks, taskId)
					}
				} else if len(result.Data) > 0 && result.Data[i].Code == 400 {
					return rstTaksks, 0
				} else if len(result.Data) > 0 && result.Data[i].Code == 401 {
					return rstTaksks, 0
				} else if len(result.Data) > 0 && result.Data[i].Code == 403 {
					return rstTaksks, 0
				} else if len(result.Data) > 0 && result.Data[i].Code == 404 {
					return rstTaksks, 0
				} else if len(result.Data) > 0 && result.Data[i].Code == 480 {
					return rstTaksks, 0
				} else {
					rstTaksks = []string{}
					return rstTaksks, 1
				}
			}
		} else if result.Code == 400 {
			return rstTaksks, 0
		} else if result.Code == 401 {
			return rstTaksks, 0
		} else if result.Code == 403 {
			return rstTaksks, 0
		} else if result.Code == 404 {
			return rstTaksks, 0
		} else if result.Code == 480 {
			return rstTaksks, 0
		} else {
			return rstTaksks, 1
		}
	} else {
		logs.Error(err)
		return rstTaksks, 1
	}
	return rstTaksks, 1
}

//阿里文本同步实时调用取得结果
//请求参数说明:文本内容的切片,可以设置多个文本
//返回结果参数说明:
//int32: 0 文本审核未通过  1审核通过  2审核处理中或网络请求出错,需下次验证(因为是文本基本不会出现2,异步调用图片会出现2)
//string: 错误码  code值记录Log用 ,string :封装后的错误内容存储数据库用
func GetTextScanResult(texts []string) (int32, string, string) {
	if nil == texts || len(texts) == 0 {
		return 1, "", ""
	}
	profile := common.Profile{AccessKeyId: ACCESS_KEY_ID, AccessKeySecret: ACCESS_KEY_SECRET}

	path := TEXT_SCAN_PATH

	clientInfo := common.ClinetInfo{Ip: SCAN_IP}

	// 构造请求数据
	bizType := "Green"
	//antispam圾文本检测
	scenes := textScenes

	tasks := []common.TextTask{}
	for i := 0; i < len(texts); i++ {
		task := common.TextTask{DataId: common.Rand().Hex(), Content: texts[i]}
		tasks = append(tasks, task)
	}

	bizData := common.TextBizData{bizType, scenes, tasks}

	var client common.AliYunClient = common.ScanClient{Profile: profile}
	if nil == client {
		return 2, "", ""
	}

	var textClientRes = ""
	func() {
		textClientRes = client.GetTextResponse(path, clientInfo, bizData)
	}()

	if "" == textClientRes {
		return 2, "", ""
	}

	result := AliRsp{}
	err := json.Unmarshal([]byte(textClientRes), &result)

	if nil == err {

		if result.Code == 200 {

			for i := 0; i < len(tasks); i++ {
				if len(result.Data) > 0 && result.Data[i].Code == 200 {

					var label string = ""
					//文本的时候,只有一种场景所以这个结果只有一个所以取得第0个判断
					if len(result.Data[i].Results) > 0 {
						label = (result.Data[i].Results)[0].Label
					}

					if "" == label {
						return 2, "", ""
					} else if "normal" != label {
						return 0, "10000", "文本审核失败-label:" + label + ";-dataid:" + tasks[i].DataId
					} else {
						continue
					}
				} else if len(result.Data) > 0 && result.Data[i].Code == 400 {
					return 0, "400", "网络出错"
				} else if len(result.Data) > 0 && result.Data[i].Code == 401 {
					return 0, "401", "网络出错"
				} else if len(result.Data) > 0 && result.Data[i].Code == 403 {
					return 0, "403", "网络出错"
				} else if len(result.Data) > 0 && result.Data[i].Code == 404 {
					return 0, "404", "网络出错"
				} else if len(result.Data) > 0 && result.Data[i].Code == 480 {
					return 0, "480", "网络出错"
				} else if len(result.Data) > 0 && result.Data[i].Code == 586 {
					return 0, "586", "网络出错"
				} else if len(result.Data) > 0 && result.Data[i].Code == 587 {
					return 0, "587", "网络出错"
				} else if len(result.Data) > 0 && result.Data[i].Code == 588 {
					return 0, "588", "网络出错"
				} else if len(result.Data) > 0 && result.Data[i].Code == 589 {
					return 0, "589", "网络出错"
				} else if len(result.Data) > 0 && result.Data[i].Code == 590 {
					return 0, "590", "网络出错"
				} else if len(result.Data) > 0 && result.Data[i].Code == 591 {
					return 0, "591", "网络出错"
				} else if len(result.Data) > 0 && result.Data[i].Code == 592 {
					return 0, "592", "网络出错"
				} else {
					//这些需要重试,文本基本是不会出现的
					return 2, "", ""
				}
			}

		} else if result.Code == 400 {
			return 0, "400", "网络出错"
		} else if result.Code == 401 {
			return 0, "401", "网络出错"
		} else if result.Code == 403 {
			return 0, "403", "网络出错"
		} else if result.Code == 404 {
			return 0, "404", "网络出错"
		} else if result.Code == 480 {
			return 0, "480", "网络出错"
		} else {
			return 2, "", ""
		}
		return 1, "", ""
	} else {
		logs.Error(err)
		return 2, "", ""
	}
	return 1, "", ""
}

//将阿里云文本调用结果转化为腾讯云的结果
//敏感词屏蔽:content为string字符串
//腾讯返回的EvilType 100:正常 20001：政治 20002：色情 20006：涉毒违法 20007：谩骂  20105：广告引流 24001：暴恐
// EvilFlag是否恶意 0：正常 1：可疑
//EvilLabel // 恶意标签，Normal：正常，Polity：涉政，Porn：色情，Illegal：违法，Abuse：谩骂，Terror：暴恐，Ad：广告，Custom：自定义关键词
func GetTextScanToTenCentRst(text string) *for_game.TextModerationRsp {
	if "" == text {
		return nil
	}

	profile := common.Profile{AccessKeyId: ACCESS_KEY_ID, AccessKeySecret: ACCESS_KEY_SECRET}

	path := TEXT_SCAN_PATH

	clientInfo := common.ClinetInfo{Ip: SCAN_IP}

	// 构造请求数据
	bizType := "Green"

	task := common.TextTask{DataId: common.Rand().Hex(), Content: text}
	tasks := []common.TextTask{task}

	bizData := common.TextBizData{bizType, textScenes, tasks}

	var client common.AliYunClient = common.ScanClient{Profile: profile}
	if nil == client {
		return nil
	}
	var textClientRes = ""
	func() {
		textClientRes = client.GetTextResponse(path, clientInfo, bizData)
	}()

	if "" == textClientRes {
		return nil
	}

	result := AliRsp{}
	err := json.Unmarshal([]byte(textClientRes), &result)

	easygo.PanicError(err)

	//封装腾讯云的结果
	bsResult := &for_game.TextModerationRsp{}

	if result.Code != 200 {
		return nil
	} else {

		if result.Data != nil && len(result.Data) > 0 {
			aliTextData := result.Data[0]
			if aliTextData.Code == 200 {

				hitWord := aliTextData.TextContent
				if aliTextData.Results != nil && len(aliTextData.Results) > 0 {
					aliTextResult := aliTextData.Results[0]
					aliToTenTextBody := AliToTenTextBodyMap[aliTextResult.Label]
					if aliToTenTextBody == nil {
						return nil
					} else {

						// 重新取得命中词
						if aliTextResult.Details != nil && len(aliTextResult.Details) > 0 {

							for i := 0; i < len(aliTextResult.Details); i++ {

								if aliTextResult.Label == aliTextResult.Details[i].Label {

									contexts := aliTextResult.Details[i].Contexts
									if nil != contexts && len(contexts) > 0 && "" != contexts[0].Context {
										hitWord = contexts[0].Context
										break
									}
								}
							}
						}

						//设置腾讯的结果集
						detailResultObj := for_game.DetailResultObj{
							EvilLabel:  aliToTenTextBody.EvilLabel,                          // 恶意标签，Normal：正常，Polity：涉政，Porn：色情，Illegal：违法，Abuse：谩骂，Terror：暴恐，Ad：广告，Custom：自定义关键词
							EvilType:   aliToTenTextBody.EvilType,                           //恶意类型
							Keywords:   []string{hitWord},                                   //命中的关键词
							Score:      int32(aliTextResult.Rate),                           //命中的模型分值
							Suggestion: AliToTenTextSuggestionMap[aliTextResult.Suggestion], //建议值:Block：打击,Review：待复审,Normal：正常
						}
						bsResult.DetailResult = []for_game.DetailResultObj{detailResultObj}
						bsResult.Keywords = detailResultObj.Keywords
						bsResult.EvilFlag = aliToTenTextBody.EvilFlag
						bsResult.EvilType = detailResultObj.EvilType

					}
				} else {
					return nil
				}
			} else {
				return nil
			}
		} else {
			return nil
		}
	}

	return bsResult
}

//将阿里云图片调用结果转化为腾讯云的结果
//对外提供的同步调用方法 直接调用解析阿里图片的方法
//腾讯返回的EvilType  100：正常 20001：政治 20002：色情 20006：涉毒违法 20007：谩骂 20103：性感 24001：暴恐
// EvilFlag是否恶意 0：正常 1：可疑
func GetImageScanToTenCentRst(imageUrl string) *for_game.ImageModerationRsp {

	if "" == imageUrl {
		return nil
	}

	profile := common.Profile{AccessKeyId: ACCESS_KEY_ID, AccessKeySecret: ACCESS_KEY_SECRET}
	path := IMAGE_SCAN_PATH

	clientInfo := common.ClinetInfo{Ip: SCAN_IP}

	// 构造请求数据
	bizType := "Green"
	task := common.ImageTask{DataId: common.Rand().Hex(), Url: imageUrl}
	tasks := []common.ImageTask{task}

	bizData := common.ImageBizData{bizType, imageScenes, tasks}

	var client common.AliYunClient = common.ScanClient{Profile: profile}
	if nil == client {
		return nil
	}

	var imageClientRes = ""

	func() {
		imageClientRes = client.GetImageResponse(path, clientInfo, bizData)
	}()

	if "" == imageClientRes {
		return nil
	}

	result := AliRsp{}
	err := json.Unmarshal([]byte(imageClientRes), &result)

	easygo.PanicError(err)

	//封装腾讯云的结果
	bsResult := &for_game.ImageModerationRsp{}

	if result.Code != 200 {
		return nil
	} else {

		if result.Data != nil && len(result.Data) > 0 {
			aliImageData := result.Data[0]
			if aliImageData.Code == 200 {
				results := aliImageData.Results
				if results != nil && len(results) > 0 {

					for i := 0; i < len(results); i++ {

						label := results[i].Label
						if label != "normal" {

							aliToTenImageBody := AliToTenImageBodyMap[label]

							if nil == aliToTenImageBody {
								return nil
							} else {
								bsResult.EvilType = aliToTenImageBody.EvilType
								bsResult.EvilFlag = aliToTenImageBody.EvilFlag
								//1、腾讯的Image图片EvilType 恶意类型 100：正常 20001：政治 20002：色情 20006：涉毒违法 20007：谩骂 20103：性感 24001：暴恐
								//性感
								if bsResult.EvilType == 20103 {
									imageHotDetect := for_game.ImageHotDetect{
										ImageCommon: for_game.ImageCommon{
											EvilType: bsResult.EvilType,      //类型
											HitFlag:  bsResult.EvilFlag,      //判定：0正常，1可疑
											Score:    int32(results[i].Rate), //得分
										}}
									bsResult.HotDetect = imageHotDetect
									//涉毒违法
								} else if bsResult.EvilType == 20006 {
									imageIllegalDetect := for_game.ImageIllegalDetect{
										ImageCommon: for_game.ImageCommon{
											EvilType: bsResult.EvilType,      //类型
											HitFlag:  bsResult.EvilFlag,      //判定：0正常，1可疑
											Score:    int32(results[i].Rate), //得分
										}}
									bsResult.IllegalDetect = imageIllegalDetect

									//政治
								} else if bsResult.EvilType == 20001 {
									imagePolityDetect := for_game.ImagePolityDetect{
										ImageCommon: for_game.ImageCommon{
											EvilType: bsResult.EvilType,      //类型
											HitFlag:  bsResult.EvilFlag,      //判定：0正常，1可疑
											Score:    int32(results[i].Rate), //得分
										}}
									bsResult.PolityDetect = imagePolityDetect
									//色情
								} else if bsResult.EvilType == 20002 {
									imagePornDetect := for_game.ImagePornDetect{
										ImageCommon: for_game.ImageCommon{
											EvilType: bsResult.EvilType,      //类型
											HitFlag:  bsResult.EvilFlag,      //判定：0正常，1可疑
											Score:    int32(results[i].Rate), //得分
										}}
									bsResult.PornDetect = imagePornDetect
									//暴恐
								} else if bsResult.EvilType == 24001 {
									imageTerrorDetect := for_game.ImageTerrorDetect{
										ImageCommon: for_game.ImageCommon{
											EvilType: bsResult.EvilType,      //类型
											HitFlag:  bsResult.EvilFlag,      //判定：0正常，1可疑
											Score:    int32(results[i].Rate), //得分
										}}
									bsResult.TerrorDetect = imageTerrorDetect
								}

							}

							break
						}
					}

				} else {
					return nil
				}
			} else {
				return nil
			}
		} else {
			return nil
		}
	}

	return bsResult

}
