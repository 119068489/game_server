package for_game

import (
	"context"
	"encoding/json"
	"game_server/easygo"
	"game_server/pb/share_message"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"time"
)

//etcd连接管理
type Client3KVManager struct {
	ServerType string                    //服务器类型:login,hall,game,shop,backstatge,square
	ServerInfo *share_message.ServerInfo //服务器信息
	PClient    *clientv3.Client
	PClient3KV clientv3.KV
	Lease      clientv3.Lease     //租约
	LeaseId    clientv3.LeaseID   //租约id
	CancleFun  context.CancelFunc //取消租约
	Mutex      easygo.RLock
}

func NewClient3KVManager(serverType string, serverInfo *share_message.ServerInfo) *Client3KVManager { // services map[string]interface{},
	p := &Client3KVManager{}
	p.Init(serverType, serverInfo)
	return p
}

//初始化
func (self *Client3KVManager) Init(serverType string, serverInfo *share_message.ServerInfo) {
	self.ServerType = serverType
	self.ServerInfo = serverInfo
}

//连接ETCD服务器
func (self *Client3KVManager) StartClintTV3() {
	//创建连接
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var err error
	self.PClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{SERVER_CENTER_ADDR},
		DialTimeout: 5 * time.Second,
	})
	logs.Info("Client3KV 连接服务器成功")
	easygo.PanicError(err)
	self.CreateLease()
	self.SetLeaseTime(10)
	self.UpdateLeaseTime()
	//创建KV
	self.PClient3KV = clientv3.NewKV(self.PClient)
	serverInfo, err1 := json.Marshal(self.ServerInfo)
	easygo.PanicError(err1)
	//通过租约put
	_, err = self.PClient3KV.Put(context.TODO(), ETCD_SERVER_PATH+self.ServerType+"/"+easygo.AnytoA(self.ServerInfo.GetSid()), string(serverInfo), clientv3.WithLease(self.LeaseId))
	if err != nil {
		logs.Info("put 失败：%s", err.Error())
		panic(err)
	}
	logs.Info("Client3KV put服务器信息:", self.ServerInfo)
}

//创建租约
func (self *Client3KVManager) CreateLease() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.Lease = clientv3.NewLease(self.PClient)
}

//设置租约时间
func (self *Client3KVManager) SetLeaseTime(t int64) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	leaseResp, err := self.Lease.Grant(context.TODO(), t)
	if err != nil {
		logs.Info("设置租约失败")
		panic(err)
	}
	self.LeaseId = leaseResp.ID
}

////设置续租
func (self *Client3KVManager) UpdateLeaseTime() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	var ctx context.Context
	ctx, self.CancleFun = context.WithCancel(context.TODO())
	leaseRespChan, err := self.Lease.KeepAlive(ctx, self.LeaseId)
	if err != nil {
		logs.Info("续租失败:", err.Error())
		panic(err)
	}
	//监听租约
	easygo.Spawn(func() {
		for {
			select {
			case leaseKeepResp := <-leaseRespChan:
				if leaseKeepResp == nil {
					logs.Info("已经关闭续租功能")
					return
				} else {
					goto END
				}
			}
		END:
			time.Sleep(500 * time.Millisecond)
		}
	})
}
func (self *Client3KVManager) CancleLease() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.CancleFun()
	_, err := self.Lease.Revoke(context.TODO(), self.LeaseId)
	if err != nil {
		logs.Info("撤销租约失败:%s", err.Error())
		panic(err)
	}
	logs.Info("撤销租约成功")
}

//监听某个key值变化
func (self *Client3KVManager) WatchClientKV(key string) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	easygo.Spawn(func() {
		wc := self.PClient.Watch(context.TODO(), key, clientv3.WithPrevKV())
		for v := range wc {
			for _, e := range v.Events {
				logs.Info("type:%v kv:%v  prevKey:%v \n ", e.Type, string(e.Kv.Key), e.PrevKv)
			}
		}
	})
}

//关闭Client
func (self *Client3KVManager) Close() {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	self.PClient.Close()
}
func (self *Client3KVManager) GetClient() *clientv3.Client {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return self.PClient
}
func (self *Client3KVManager) GetClientKV() clientv3.KV {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	return self.PClient3KV
}
