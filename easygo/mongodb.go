package easygo

import (
	"errors"
	"fmt"
	"runtime"

	"strings"

	"github.com/akqp2019/mgo"
)

type IMongoDBManager interface {
	InitOneDB(kwargs map[string]string) *mgo.Session
	InitSomeDB(urls map[DB_NAME]map[string]string)
	GetDB(databaseName string) (*mgo.Database, func())
	GetC(databaseName string, collectionName string) (*mgo.Collection, func())

	AddDbSession(dbName string, kwargs map[string]string)
	GetMongoSessions() map[string]*mgo.Session
	ExistDatabase(databaseName string) bool
}

type MongoDBManager struct {
	Me            IMongoDBManager
	MONGO_SESSION map[string]*mgo.Session
}

func NewMongoDBManager() *MongoDBManager {
	p := &MongoDBManager{}
	p.Init(p)
	return p
}

func (self *MongoDBManager) Init(me IMongoDBManager) {
	self.Me = me
	self.MONGO_SESSION = make(map[string]*mgo.Session)
}

func (self *MongoDBManager) GetC(databaseName string, collectionName string) (col *mgo.Collection, closeFun func()) {
	db, f := self.Me.GetDB(databaseName)
	c := db.C(collectionName)
	return c, f
}

func (self *MongoDBManager) GetDB(databaseName string) (*mgo.Database, func()) {
	if databaseName == "" {
		panic("数据库名字必须传给我")
	}
	session, ok := self.MONGO_SESSION[databaseName]
	if !ok {
		str := fmt.Sprintf("不存在 key 为 %v 的 MongoDB 连接，也许是未正确连接", databaseName)
		panic(str)
	}
	session = session.Copy()

	finalizer := func(c *mgo.Database) {
		session.Close()
	}

	db := session.DB(databaseName)
	runtime.SetFinalizer(db, finalizer) // 在最坏的情况下，使用者没有 defer closeFunc 或直接调用 closeFunc,仍然靠 GC 实现关闭

	hasClose := false // upvalue. 避免重复 close
	closeF := func() {
		if hasClose {
			return
		}
		hasClose = true
		runtime.SetFinalizer(db, nil)
		session.Close()
	}
	return db, closeF
}

func (self *MongoDBManager) InitOneDB(kwargs map[string]string) *mgo.Session {
	url := "mongodb://user:password@host:port?maxPoolSize=max_pool_size"
	// 以下都是有效格式
	// url := "mongodb://user:password@host:port"
	// url := "mongodb://host:port"
	// url := "localhost:40001?maxPoolSize=512"

	var args []string
	for k, v := range kwargs {
		args = append(args, k)
		args = append(args, v)
	}
	s := strings.NewReplacer(args...).Replace(url)

	session, err := mgo.Dial(s)
	// 也可以使用 mgo.DialWithInfo(info *DialInfo)
	if err != nil {
		s := fmt.Sprintf("MongoDB 启动了吗？用户名，密码对了吗？%v;url=%s", err, s)
		e := errors.New(s)
		panic(e)
	}
	return session
}

// 初始化多个数据库
func (self *MongoDBManager) InitSomeDB(urls map[DB_NAME]map[string]string) {
	// 数据库名与 key 名保持一致，方便管理
	for databaseName, kwargs := range urls {
		self.Me.AddDbSession(databaseName, kwargs)
	}
}

func (self *MongoDBManager) AddDbSession(dbName string, kwargs map[string]string) {
	// logs.Info("mongo初始化:", dbName)
	session := self.Me.InitOneDB(kwargs)
	self.MONGO_SESSION[dbName] = session
}

func (self *MongoDBManager) ExistDatabase(databaseName string) bool {
	_, ok := self.MONGO_SESSION[databaseName]
	return ok
}

func (self *MongoDBManager) GetMongoSessions() map[string]*mgo.Session {
	return self.MONGO_SESSION
}

//------------------------------------------------------

func MongoStats() []string {
	var slice []string
	stats := mgo.GetStats()
	slice = append(slice, fmt.Sprintf("Clusters: %d", stats.Clusters))
	slice = append(slice, fmt.Sprintf("MasterConns: %d", stats.MasterConns))
	slice = append(slice, fmt.Sprintf("SlaveConns: %d", stats.SlaveConns))

	slice = append(slice, fmt.Sprintf("SentOps: %d", stats.SentOps))
	slice = append(slice, fmt.Sprintf("ReceivedOps: %d", stats.ReceivedOps))
	slice = append(slice, fmt.Sprintf("ReceivedDocs: %d", stats.ReceivedDocs))

	slice = append(slice, fmt.Sprintf("SocketsAlive: %d", stats.SocketsAlive))
	slice = append(slice, fmt.Sprintf("SocketsInUse: %d", stats.SocketsInUse))
	slice = append(slice, fmt.Sprintf("SocketRefs: %d", stats.SocketRefs))
	return slice
}
