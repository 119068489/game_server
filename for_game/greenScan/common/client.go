package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ScanClient struct {
	Profile Profile
}

func (self ScanClient) GetImageResponse(path string, clinetInfo ClinetInfo, bizData ImageBizData) string {
	clientInfoJson, _ := json.Marshal(clinetInfo)
	bizDataJson, _ := json.Marshal(bizData)

	client := &http.Client{}
	req, err := http.NewRequest(method, host+path+"?clientInfo="+url.QueryEscape(string(clientInfoJson)), strings.NewReader(string(bizDataJson)))

	if err != nil {
		// handle error
		return ErrorResult(err)
	} else {
		addRequestHeader(string(bizDataJson), req, string(clientInfoJson), path, self.Profile.AccessKeyId, self.Profile.AccessKeySecret)

		response, err1 := client.Do(req)
		//defer response.Body.Close()
		if err1 != nil {
			// handle error
			return ErrorResult(err1)
		}

		body, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			// handle error
			return ErrorResult(err)
		} else {
			return string(body)
		}
	}
}

func (self ScanClient) GetTextResponse(path string, clinetInfo ClinetInfo, bizData TextBizData) string {
	clientInfoJson, _ := json.Marshal(clinetInfo)
	bizDataJson, _ := json.Marshal(bizData)

	client := &http.Client{}
	req, err := http.NewRequest(method, host+path+"?clientInfo="+url.QueryEscape(string(clientInfoJson)), strings.NewReader(string(bizDataJson)))

	if err != nil {
		// handle error
		return ErrorResult(err)
	} else {
		addRequestHeader(string(bizDataJson), req, string(clientInfoJson), path, self.Profile.AccessKeyId, self.Profile.AccessKeySecret)

		response, err1 := client.Do(req)
		if err1 != nil {
			// handle error
			return ErrorResult(err1)
		}

		body, err := ioutil.ReadAll(response.Body)

		response.Body.Close()
		if err != nil {
			// handle error
			return ErrorResult(err)
		} else {
			return string(body)
		}
	}
}

func (self ScanClient) GetVideoResponse(path string, clinetInfo ClinetInfo, bizData VideoBizData) string {
	clientInfoJson, _ := json.Marshal(clinetInfo)
	bizDataJson, _ := json.Marshal(bizData)

	client := &http.Client{}
	req, err := http.NewRequest(method, host+path+"?clientInfo="+url.QueryEscape(string(clientInfoJson)), strings.NewReader(string(bizDataJson)))

	if err != nil {
		// handle error
		return ErrorResult(err)
	} else {
		addRequestHeader(string(bizDataJson), req, string(clientInfoJson), path, self.Profile.AccessKeyId, self.Profile.AccessKeySecret)

		response, err1 := client.Do(req)
		if err1 != nil {
			// handle error
			return ErrorResult(err1)
		}
		//defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)

		response.Body.Close()
		if err != nil {
			// handle error
			return ErrorResult(err)
		} else {
			return string(body)
		}
	}
}

func (self ScanClient) GetRstResponse(path string, clinetInfo ClinetInfo, bizData []string) string {
	clientInfoJson, _ := json.Marshal(clinetInfo)
	bizDataJson, _ := json.Marshal(bizData)

	client := &http.Client{}
	req, err := http.NewRequest(method, host+path+"?clientInfo="+url.QueryEscape(string(clientInfoJson)), strings.NewReader(string(bizDataJson)))

	if err != nil {
		// handle error
		return ErrorResult(err)
	} else {
		addRequestHeader(string(bizDataJson), req, string(clientInfoJson), path, self.Profile.AccessKeyId, self.Profile.AccessKeySecret)

		response, err1 := client.Do(req)
		if err1 != nil {
			// handle error
			return ErrorResult(err1)
		}
		//defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)

		response.Body.Close()
		if err != nil {
			// handle error
			return ErrorResult(err)
		} else {
			return string(body)
		}
	}
}

type AliYunClient interface {
	GetImageResponse(path string, clinetInfo ClinetInfo, bizData ImageBizData) string

	GetTextResponse(path string, clinetInfo ClinetInfo, bizData TextBizData) string

	GetVideoResponse(path string, clinetInfo ClinetInfo, bizData VideoBizData) string

	GetRstResponse(path string, clinetInfo ClinetInfo, bizData []string) string
}
