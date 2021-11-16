package util

import (
	"bytes"
	"context"
	"fmt"
	"game_server/easygo"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

/*腾讯存储桶*/
//电竞服开启
const (
	host       = "https://im-resource-1253887233.file.myqcloud.com/"
	bucketName = "im-resource-1253887233"
	region     = "ap-guangzhou"
	secretId   = "AKIDOYR4Xst9ZicIwJrJ7ex2rPlqgY9VbIj1" // 可传固定密钥或者临时密钥
	secretKey  = "DDCVaPxb2evwpgJfleZEm4RXPAe7KOCk"     // 可传固定密钥或者临时密钥
)

type Bucket struct {
	B *cos.Client
}

type ObjectList struct {
	Title *string
	Time  *string
	Size  *int64
}

func NewQQbucket() *Bucket {
	u := cos.NewBucketURL(bucketName, region, true)
	b := &cos.BaseURL{BucketURL: u}
	// 1.永久密钥
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})

	return &Bucket{
		B: client,
	}

}

//查询存储桶
func (self *Bucket) GetBuckets() {
	s, _, err := self.B.Service.Get(context.Background())
	if err != nil {
		panic(err)
	}

	for _, b := range s.Buckets {
		fmt.Printf("%#v\n", b)
	}
}

//查询存储桶Object列表 path 路径, pagesize 请求条数 ,marker 上次请求最后一个key
func (self *Bucket) GetObjectList(path, marker string, pagesize int32) []*ObjectList {
	opt := &cos.BucketGetOptions{
		Prefix:  path,
		MaxKeys: int(pagesize),
	}

	if marker != "" {
		opt.Marker = marker
	}

	v, _, err := self.B.Bucket.Get(context.Background(), opt)
	if err != nil {
		panic(err)
	}

	var lis []*ObjectList
	for _, c := range v.Contents {
		li := &ObjectList{
			Title: easygo.NewString(c.Key),
			Time:  easygo.NewString(c.LastModified),
			Size:  easygo.NewInt64(c.Size),
		}
		lis = append(lis, li)
	}

	return lis
}

//上传本地Object对象到存储桶 pathFileName 路径+文件名=test/objectPut.go, pagesize 请求条数
func (self *Bucket) ObjectPut(pathFileName, localPathFileName string) string {

	// 对象键（pathFileName）是对象在存储桶中的唯一标识。
	// 例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go

	// // 1.通过字符串上传对象
	// f := strings.NewReader("test")

	// _, err := self.B.Object.Put(context.Background(), pathFileName, f, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// 2.通过本地文件上传对象
	_, err := self.B.Object.PutFromFile(context.Background(), pathFileName, localPathFileName, nil)
	if err != nil {
		logs.Error(err)
		return ""
	}

	return host + pathFileName
}

//删除存储桶Object对象
func (self *Bucket) ObjectDel(pathFileName string) {
	_, err := self.B.Object.Delete(context.Background(), pathFileName)
	if err != nil {
		panic(err)
	}
}

//上传远端Object对象到存储桶 pathFileName 路径+文件名=test/objectPut.go
func (self *Bucket) ObjectPutRemote(pathFileName, fileUrl string) string {
	res, err := http.Get(fileUrl)
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	_, err = self.B.Object.Put(context.Background(), pathFileName, res.Body, nil)
	if err != nil {
		logs.Error(err)
		return ""
	}

	return host + pathFileName
}

//上传base64 Object对象到存储桶 pathFileName 路径+文件名=test/objectPut.go, file 文件base64码
func (self *Bucket) ObjectPutByte(pathFileName string, file []byte) string {
	rdata := bytes.NewReader(file)
	if rdata.Size() == 0 {
		return ""
	}
	_, err2 := self.B.Object.Put(context.Background(), pathFileName, rdata, nil)
	if err2 != nil {
		logs.Error(err2)
		return ""
	}

	return host + pathFileName
}

//替换网页内容中的图片为存储桶图片
func (self *Bucket) ReplaceImg(content string) string {
	r := strings.NewReader(content)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logs.Error(err.Error())
		return content //发生错误返回原文
	}

	locl, _ := url.Parse(host)
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		img, _ := s.Attr("src")
		u, _ := url.Parse(img)
		if u.Host != locl.Host {
			fileName := filepath.Base(img)
			pathfileName := path.Join("backstage", "upload", fileName)
			url := self.ObjectPutRemote(pathfileName, img)
			if url != "" {
				s.RemoveAttr("src")
				s.SetAttr("src", url)
			}
		}
	})

	dContent, _ := doc.Find("body").Html()

	return dContent
}
