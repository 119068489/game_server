// 违禁图监听回调处理

package backstage

import (
	"encoding/json"
	"game_server/easygo"
	"github.com/astaxie/beego/logs"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"io/ioutil"
	"net/http"
	"strings"
)

func (self *WebHttpServer) CheckPicture(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body) //获取post的数据
	logs.Info("收到返回数据:", string(data))
	self.ParseData(data)
}
func (self *WebHttpServer) ParseData(data []byte) {
	params := make(easygo.KWAT)
	err := json.Unmarshal([]byte(data), &params)
	easygo.PanicError(err)
	dd := make(easygo.KWAT)
	js, err := json.Marshal(params["data"])
	err = json.Unmarshal(js, &dd)
	easygo.PanicError(err)
	url := dd.GetString("url")
	if url != "" {
		url = strings.Replace(url, "im-resource-1253887233.picgz.myqcloud.com", "im-resource-1253887233.file.myqcloud.com", 1)
		//cdn缓存刷新接口
		self.FlushPictureUrl(url)
	}
}

//刷新指定url缓存
func (self *WebHttpServer) FlushPictureUrl(urls string) {
	c, err := v20180606.NewClient(common.NewCredential("AKIDOYR4Xst9ZicIwJrJ7ex2rPlqgY9VbIj1", "DDCVaPxb2evwpgJfleZEm4RXPAe7KOCk"), "ap-guangzhou", profile.NewClientProfile())
	easygo.PanicError(err)
	req := v20180606.NewPurgeUrlsCacheRequest()
	req.Urls = append(req.Urls, easygo.NewString(urls))
	_, err = c.PurgeUrlsCache(req)
	easygo.PanicError(err)
	logs.Info("刷新完成")
}
