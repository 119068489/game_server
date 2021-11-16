package common

type ClinetInfo struct {
	SdkVersion  string `json:"sdkVersion"`
	CfgVersion  string `json:"cfgVersion"`
	UserType    string `json:"userType"`
	UserId      string `json:"userId"`
	UserNick    string `json:"userNick"`
	Avatar      string `json:"avatar"`
	Imei        string `json:"imei"`
	Imsi        string `json:"imsi"`
	Umid        string `json:"umid"`
	Ip          string `json:"ip"`
	Os          string `json:"os"`
	Channel     string `json:"channel"`
	HostAppName string `json:"hostAppName"`
	HostPackage string `json:"hostPackage"`
	HostVersion string `json:"hostVersion"`
}

type ImageTask struct {
	DataId string `json:"dataId"` //图片数据ID。需要保证在一次请求中所有的ID不重复。
	Url    string `json:"url"`    //图片URL
}

type ImageBizData struct {
	BizType string      `json:"bizType"` //该字段用于标识业务场景
	Scenes  []string    `json:"scenes"`  //指定图片检测场景,	参考阿里云图片场景设置
	Tasks   []ImageTask `json:"tasks"`   //图片检测任务列表
}

type TextTask struct {
	DataId  string `json:"dataId"`  //文本数据ID。需要保证在一次请求中所有的ID不重复。
	Content string `json:"content"` //文本内容
}

type TextBizData struct {
	BizType string     `json:"bizType"` //该字段用于标识业务场景
	Scenes  []string   `json:"scenes"`  //指定文本检测场景,	参考阿里云文本场景设置(取值：antispam，表示文本垃圾内容检测)
	Tasks   []TextTask `json:"tasks"`   //文本检测任务列表
}

type VideoTask struct {
	DataId string `json:"dataId"` //视频数据ID。需要保证在一次请求中所有的ID不重复。
	Url    string `json:"url"`    //视频URL
}

type VideoBizData struct {
	BizType string      `json:"bizType"` //该字段用于标识业务场景
	Scenes  []string    `json:"scenes"`  //指定视频检测场景,	参考阿里云视频场景设置
	Tasks   []VideoTask `json:"tasks"`   //视频检测任务列表
}
