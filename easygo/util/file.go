package util

import (
	"bufio"
	"errors"
	"github.com/astaxie/beego/logs"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// GetIntSliceFromFile 从文件中读取[]int结构
// 用\n作为行分隔符, "splitString"作为列分隔符
func GetIntSliceFromFile(file, splitString string) ([]int, error) {
	s := make([]int, 0)
	f, err := os.Open(file)
	if err != nil {
		return s, err
	}
	defer f.Close()

	// 读取文件到buffer里边
	buf := bufio.NewReader(f)
	for {
		// 按照换行读取每一行
		l, err := buf.ReadString('\n')
		// 跳过空行
		if l == "\n" {
			continue
		}

		lineSplit := strings.SplitN(l, splitString, 1024)
		for _, v := range lineSplit {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			value, _ := strconv.Atoi(v)
			s = append(s, value)
		}
		if err != nil {
			break
		}
	}
	return s, nil
}

//电竞服开启
// CreateDateDir 根据当前日期来创建文件夹
func CreateDateDir(Path string) string {
	folderName := time.Now().Format("20060102")
	folderPath := filepath.Join(Path, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.Mkdir(folderPath, 0777) //0777也可以os.ModePerm
		os.Chmod(folderPath, 0777)
	}
	return folderPath
}

//检查文件是否存在
func IsFileExist(filePath string, fileSize int64) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// fmt.Println(info)
		return false
	}

	if fileSize == info.Size() {
		// fmt.Println("安装包已存在！", info.Name(), info.Size(), info.ModTime())
		return true
	}
	err = os.Remove(filePath)
	if err != nil {
		logs.Error("remove", err)
	}
	return false
}

func DownloadFile(url string, localPath string, folderName ...string) error {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	folderName = append(folderName, "download")
	os.Mkdir("./"+folderName[0], os.ModePerm)

	folderPath := filepath.Join(folderName[0], localPath)
	tmpFilePath := folderPath + ".download"
	// fmt.Println(tmpFilePath)
	//创建一个http client
	client := new(http.Client)
	client.Timeout = time.Second * 120 //设置超时时间
	//get方法获取资源
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	//读取服务器返回的文件大小
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		logs.Error("size", err)
	}
	if IsFileExist(localPath, fsize) {
		return errors.New("file is not nil")

	}
	// fmt.Println("fsize", fsize)
	//创建文件
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		return errors.New("body is null")
	}
	defer resp.Body.Close()
	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}

	if err == nil {
		file.Close()
		err = os.Rename(tmpFilePath, folderPath)
		if err != nil {
			logs.Error("rename", err)
		}
	}
	return err
}

func DeleteFile(localPathFileName string, folderName ...string) error {
	folderName = append(folderName, "download")
	localPathFileName = folderName[0] + "/" + localPathFileName

	err := os.Remove(localPathFileName)
	if err != nil {
		return err
	}
	return nil
}
