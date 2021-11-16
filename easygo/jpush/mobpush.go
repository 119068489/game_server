package jpushclient

// Push   推送对象
type Push struct {
	Source     string      `json:"source,omitempty"` //枚举值 ”webapi“
	Appkey     string      `json:"appkey,omitempty"`
	PushTarget *PushTarget `json:"pushTarget,omitempty"`
	PushNotify *PushNotify `json:"pushNotify,omitempty"`
}

//	PushNotify	推送展示细节
type PushNotify struct {
	Plats          []int32                  `json:"plats"`          //设备类型: 1 Android;2 Ios
	IosProduction  int32                    `json:"iosProduction"`  //plat = 2,0 测试;1 生成环境
	OfflineSeconds int32                    `json:"offlineSeconds"` //离线消息保存时间
	Content        string                   `json:"content"`        //推送内容 必填
	Title          string                   `json:"title"`          //推送标题
	Type           int32                    `json:"type"`           //推送类型:1 通知;2 自定义 必填
	TaskCron       int32                    `json:"taskCron"`       //是否定时任务: 0 否;1 是
	TaskTime       int32                    `json:"taskTime"`       //定时发送时间:
	IosNotify      *IosNotify               `json:"iosNotify"`      //IOS设置
	ExtrasMapList  []map[string]interface{} `json:"extrasMapList"`  //扩展
}

//	PushTarget	推送目标
type PushTarget struct {
	Target int32    `json:"target"` //目标类型: 1 广播;2 别名;3 标签;4 regid 必填
	Alias  []string `json:"alias"`  //别名
}

//	Login	返回内容
type Login struct {
	Status int         `json:"status"`
	Error  string      `json:"error"`
	Res    interface{} `json:"res"`
}

type IosNotify struct {
	Badge     int32 `json:"badge"`
	BadgeType int32 `json:"badgeType"`
}

func (push *Push) SetSource(source string) {
	push.Source = source
}

func (push *Push) SetAppkey(appkey string) {
	push.Appkey = appkey
}

func (push *Push) SetPushTarget(pushTarget *PushTarget) {
	push.PushTarget = pushTarget
}

func (push *Push) SetPushNotify(pushNotify *PushNotify) {
	push.PushNotify = pushNotify
}

func (pushNotify *PushNotify) SetPlats(plats []int32) {
	pushNotify.Plats = plats
}

func (pushNotify *PushNotify) SetContent(content string) {
	pushNotify.Content = content
}

func (pushNotify *PushNotify) SetTitle(title string) {
	pushNotify.Title = title
}

func (pushNotify *PushNotify) SetIosProduction(Type int32) {
	pushNotify.IosProduction = Type
}
func (pushNotify *PushNotify) SetType(Type int32) {
	pushNotify.Type = Type
}

func (pushNotify *PushNotify) SetExtrasMapList(extrasMapList []map[string]interface{}) {
	pushNotify.ExtrasMapList = extrasMapList
}

func (pushNotify *PushNotify) SetIosNotify(iosNotify *IosNotify) {
	pushNotify.IosNotify = iosNotify
}

func (pushTarget *PushTarget) SetTarget(target int32) {
	pushTarget.Target = target
}

func (pushTarget *PushTarget) SetAlias(alias []string) {
	pushTarget.Alias = alias
}

func (iosNotify *IosNotify) SetBadge(badge int32) {
	iosNotify.Badge = badge
}

func (iosNotify *IosNotify) SetBadgeType(badgeType int32) {
	iosNotify.BadgeType = badgeType
}
