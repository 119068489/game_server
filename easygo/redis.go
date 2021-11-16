package easygo

import (
	"github.com/astaxie/beego/logs"
)

type IRedisManager interface {
	GetC() RedisCache
}

type RedisManager struct {
	Me IRedisManager
	RedisCache
	Mutex RLock
}

func NewRedisManager(yaml IYamlConfig) *RedisManager { // services map[string]interface{},
	p := &RedisManager{}
	p.Init(p, yaml)
	return p
}

//初始化
func (self *RedisManager) Init(me IRedisManager, yaml IYamlConfig) {
	self.Me = me
	host := yaml.GetValueAsString("REDIS_SERVER_ADDR")
	pass := yaml.GetValueAsString("REDIS_SERVER_PASS")
	db := yaml.GetValueAsInt("REDIS_SERVER_DATABASE")
	self.RedisCache = NewRedisCache(db, host, REDIS_DEFAULT, pass)
	logs.Info("连接redis服务器成功:" + host)

}

func (self *RedisManager) GetC() RedisCache {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return self.RedisCache
}
