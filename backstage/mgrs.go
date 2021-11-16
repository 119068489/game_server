//
package backstage

import (
	"game_server/easygo"
	"sync"
)

// 浏览器 endpoint 管理器

type BrowerEndpointManager struct {
	sync.Map
}

func NewBrowerEndpointManager() *BrowerEndpointManager {
	p := &BrowerEndpointManager{}
	return p
}

func (self *BrowerEndpointManager) LoadEndpoint(endpointId ENDPOINT_ID) *BrowerEndpoint {
	value, ok := self.Load(endpointId)
	if ok {
		return value.(*BrowerEndpoint)
	}
	return nil
}

//============================================================================================================

type BrowerEndpointMapping struct {
	easygo.Mapping
}

func NewBrowerEndpointMapping() *BrowerEndpointMapping {
	p := &BrowerEndpointMapping{}
	p.Init(p)
	return p
}

func (self *BrowerEndpointMapping) LoadEndpoint(userId USER_ID) IBrowerEndpoint {
	value := self.Mapping.Load(userId)
	if value != nil {
		return value.(IBrowerEndpoint)
	}
	return nil
}

func (self *BrowerEndpointMapping) StoreEndpoint(userId USER_ID, endpointId ENDPOINT_ID) {
	if endpointId == 0 {
		panic("endpoint id 不可能是 0")
	}
	fetch := func() interface{} {
		v := BrowerEpMgr.LoadEndpoint(endpointId)
		if v == nil {
			return nil
		}
		return v
	}
	self.Store(userId, fetch)
}

func (self *BrowerEndpointMapping) Delete(userId USER_ID) { // overwrite 为了让使用者清楚参数的类型
	self.Mapping.Delete(userId)
}

//----------------------------------------------
