package jpushclient

// PushV2   推送对象
type PushV2 struct {
	Appkey       string                 `json:"source,omitempty"` //枚举值 ”webapi“
	Plats        []int32                `json:"plats,omitempty"`
	Target       int                    `json:"target,omitempty"`
	Content      string                 `json:"content,omitempty"`
	Type         int                    `json:"type,omitempty"`
	IosBadge     int                    `json:"iosBadge,omitempty"`
	Alias        []string               `json:"alias,omitempty"`
	AndroidTitle string                 `json:"androidTitle,omitempty"`
	IosTitle     string                 `json:iosTitle,omitempty"`
	Extras       map[string]interface{} `json:"extras,omitempty"`
}

func (push *PushV2) SetAppkey(appkey string) {
	push.Appkey = appkey
}

func (push *PushV2) SetPlats(plats []int32) {
	push.Plats = plats
}

func (push *PushV2) SetTarget(target int) {
	push.Target = target
}

func (push *PushV2) SetContent(content string) {
	push.Content = content
}

func (push *PushV2) SetType(Type int) {
	push.Type = Type
}

func (push *PushV2) SetExtras(extras map[string]interface{}) {
	push.Extras = extras
}

func (push *PushV2) SetAlias(alias []string) {
	push.Alias = alias
}

func (push *PushV2) SetAndroidTitle(androidTitle string) {
	push.AndroidTitle = androidTitle
}

func (push *PushV2) SetIosTitle(iosTitle string) {
	push.IosTitle = iosTitle
}

func (push *PushV2) SetIosBadge(iosBadge int) {
	push.IosBadge = iosBadge
}
