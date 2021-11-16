package for_game

import (
	"game_server/easygo"
	"game_server/pb/share_message"

	"github.com/akqp2019/mgo/bson"
)

//系统功能参数管理器
type SysParameterManager struct {
	SysParameterList map[string]*share_message.SysParameter
	Mutex            easygo.Mutex
}

func NewSysParameterManager() *SysParameterManager {
	p := &SysParameterManager{}
	p.Init()
	return p
}
func (self *SysParameterManager) Init() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.SysParameterList = make(map[string]*share_message.SysParameter)
	list := QuerySysParameterList()
	for _, v := range list {
		self.SysParameterList[v.GetId()] = v
	}
	// logs.Info("系统参数:", self.SysParameterList)
}

func (self *SysParameterManager) UpLoad(id string) {
	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, TABLE_SYS_PARAMETER)
	defer closeFun()

	queryBson := bson.M{"_id": id}
	var val *share_message.SysParameter
	err := col.Find(queryBson).One(&val)
	easygo.PanicError(err)
	self.SysParameterList[id] = val
}

func (self *SysParameterManager) GetSysParameter(id string) *share_message.SysParameter {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	v, ok := self.SysParameterList[id]
	if !ok {
		return nil
	}
	return v
}

//1、text屏蔽，2、image屏蔽
func (self *SysParameterManager) GetTextModeration(t int32) int32 {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	v, ok := self.SysParameterList[OBJ_MODERATIONS]
	if !ok {
		return 0
	}
	for _, v := range v.GetTextModerations() {
		if v.GetId() == t {
			return int32(v.GetCount())
		}
	}
	return 0
}
func (self *SysParameterManager) GetImageModeration(t int32) int32 {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	v, ok := self.SysParameterList[OBJ_MODERATIONS]
	if !ok {
		return 0
	}
	for _, v := range v.GetImageModerations() {
		if v.GetId() == t {
			return int32(v.GetCount())
		}
	}
	return 0
}
