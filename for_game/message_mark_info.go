package for_game

import (
	"game_server/easygo"
	"time"
)

type MarkInfo struct {
	Mark      string
	Timestamp int64
	//ImgCheckCode int32
}

//func (self *MarkInfo) GetImgCheckCode() int32 {
//	return self.ImgCheckCode
//}

type MessageMarkInfoMgr struct {
	MessageMarkInfo map[int32]map[string]*MarkInfo
	//ImgCheckCode    map[string]*MarkInfo
	Mutex easygo.RLock
}

func NewMessageMarkInfoMgr() *MessageMarkInfoMgr {
	p := &MessageMarkInfoMgr{}
	p.Init()
	return p
}

func (self *MessageMarkInfoMgr) Init() {
	self.MessageMarkInfo = make(map[int32]map[string]*MarkInfo)
	//self.ImgCheckCode = make(map[string]*MarkInfo)
	easygo.AfterFunc(time.Second, self.MessageCheckValid)
}

func (self *MessageMarkInfoMgr) AddMessageMarkInfo(t int32, phone string, code string) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	info := &MarkInfo{
		Mark:      code,
		Timestamp: time.Now().Unix(),
	}
	var message map[string]*MarkInfo
	if _, ok := self.MessageMarkInfo[t]; !ok {
		message = make(map[string]*MarkInfo)
	} else {
		message = self.MessageMarkInfo[t]
	}
	message[phone] = info
	self.MessageMarkInfo[t] = message
}

func (self *MessageMarkInfoMgr) GetAllMessageMarkInfo() map[int32]map[string]*MarkInfo {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	copyList := make(map[int32]map[string]*MarkInfo)
	for key, value := range self.MessageMarkInfo {
		copyList[key] = value
	}
	return copyList
}

func (self *MessageMarkInfoMgr) GetMessageMarkInfo(t int32, phone string) *MarkInfo {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return self.MessageMarkInfo[t][phone]
}

func (self *MessageMarkInfoMgr) CheckPhoneVaild(phone string) bool {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	t := time.Now().Unix()
	for _, info := range self.MessageMarkInfo {
		for p, m := range info {
			if p == phone {
				if t-m.Timestamp < 60 { //如果有验证码发送没超过60秒
					return false
				}
			}
		}
	}
	return true
}

//func (self *MessageMarkInfoMgr) AddImgCheckCode(phone string, img int32) {
//	self.Mutex.Lock()
//	defer self.Mutex.Unlock()
//
//	self.ImgCheckCode[phone] = &MarkInfo{
//		ImgCheckCode: img,
//		Timestamp:    time.Now().Unix(),
//	}
//}

//func (self *MessageMarkInfoMgr) GetImgCheckCode(phone string) *MarkInfo {
//	self.Mutex.Lock()
//	defer self.Mutex.Unlock()
//	return self.ImgCheckCode[phone]
//}

//func (self *MessageMarkInfoMgr) GetAllImgCheckCode() map[string]*MarkInfo {
//	self.Mutex.Lock()
//	defer self.Mutex.Unlock()
//	copyList := make(map[string]*MarkInfo)
//	for key, value := range self.ImgCheckCode {
//		copyList[key] = value
//	}
//
//	return copyList
//}

//func (self *MessageMarkInfoMgr) DeleteImgCheckCode(phones []string) {
//	self.Mutex.Lock()
//	defer self.Mutex.Unlock()
//	if len(phones) > 0 {
//		for _, phone := range phones {
//			delete(self.ImgCheckCode, phone)
//		}
//	}
//}

func (self *MessageMarkInfoMgr) MessageCheckValid() {
	ti := time.Now().Unix()
	for t, msg := range self.GetAllMessageMarkInfo() {
		for phone, info := range msg {
			timestamp := info.Timestamp
			if ti-timestamp >= int64(MessageCheckTime) {
				delete(self.MessageMarkInfo[t], phone)
			}
		}
	}

	easygo.AfterFunc(time.Second, self.MessageCheckValid)
}

//func (self *MessageMarkInfoMgr) CheckCode(code int32, phone string) bool {
//	checkCode := self.ImgCheckCode[phone]
//	if checkCode != nil {
//		if code < checkCode.GetImgCheckCode()+5 && code > checkCode.GetImgCheckCode()-5 {
//			return true
//		}
//	}
//	return false
//}
