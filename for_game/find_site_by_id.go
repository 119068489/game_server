package for_game

import (
	"game_server/easygo"
	"github.com/akqp2019/mgo/bson"
)

type FindSite struct {
	Info map[PLAYER_ID]SITE

	Mutex easygo.RLock
}

func NewFindSite() *FindSite {
	p := &FindSite{}
	p.Init()
	return p
}

func (self *FindSite) Init() {
	self.Info = make(map[PLAYER_ID]SITE)
}

func (self *FindSite) Find(playerId PLAYER_ID) SITE {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	if site, ok := self.Info[playerId]; ok {
		return site
	}

	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, "player_site")
	defer closeFun()
	var b bson.M
	err := col.Find(bson.M{"_id": playerId}).One(&b)
	easygo.PanicError(err)

	site := b["Site"].(string)
	self.Info[playerId] = site

	return site
}

func (self *FindSite) Insert(playerId PLAYER_ID, site SITE) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	if _, ok := self.Info[playerId]; ok {
		return
	}

	col, closeFun := easygo.MongoMgr.GetC(MONGODB_NINGMENG, "player_site")
	defer closeFun()
	_ = col.Insert(bson.M{"_id": playerId, "Site": site})

	self.Info[playerId] = site
}

var FindSiteById = NewFindSite()
