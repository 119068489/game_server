package wish

import (
	"game_server/easygo"
	"game_server/pb/share_message"
)

//=======================================

type WishBoxObj struct {
	BoxId int64
	Mutex easygo.RLock
}

func NewWishBox(boxId int64) *WishBoxObj {
	p := &WishBoxObj{
		BoxId: boxId,
	}
	WishBoxMgr.Store(boxId, p)
	return p
}

//=====================对外接口====================

// 对外接口，获取盲盒数据.
func GetWishBoxObj(boxId int64) *WishBoxObj {
	WishBoxMgr.Mutex.Lock()
	defer WishBoxMgr.Mutex.Unlock()
	wishBox := WishBoxMgr.LoadWishBox(boxId)
	if wishBox != nil {
		return wishBox
	}
	obj := NewWishBox(boxId)
	return obj
}

// 从redis中获取盲盒数据.
func (self *WishBoxObj) GetWishBoxData() *share_message.WishBox {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return nil
}
