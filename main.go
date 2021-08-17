package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
	"strings"
	"time"
)

var (
	login loginResp
	tasks getTaskResp
	trash trashResp
	ip string
	username string
	password string
)

func loadConfig()error{
	txt,err:=ioutil.ReadFile("config.txt")
	if err!=nil{
		Error.Printf("读取配置错误,error:%v\n",err)
		return err
	}
	config:=string(txt)
	configList:=strings.Split(config,"\r\n")
	//ip=8.8.8.8
	ip=configList[0][3:]
	//username=admin
	username=configList[1][9:]
	//password=tji1f3GNK8k-1Xqx4AkT
	password=configList[2][9:]
	return nil
}

func main() {
	run()
}


func run(){
	err:=loadConfig()
	if err!=nil{
		return
	}
	client := resty.New()
	requst:=client.R()
	header := map[string]string{
		"Accept":        "application/json",
		"Authorization": "Basic YWRtaW46YWRtaW4=",
		"User-Agent" : "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
	}
	//登录获取token
	Info.Println("登录获取Token")
	data := map[string]string{
		"username": username,
		"password": password,
	}
	err = getToken(requst, header,data)
	if err!=nil {
		Error.Printf("登录出错\n")
		return
	}
	header["Authorization"] = fmt.Sprintf("Bearer %s", login.Token)

	//获取任务列表
	Info.Println("获取任务列表")
	err= getList(requst,header)
	if err!=nil {
		Error.Printf("获取任务列表出错\n")
		return
	}
	duplicatedId,count := getDuplicateTasks()
	Info.Printf("原任务列表共有%d条任务\n",len(tasks.Data))
	Info.Printf("去重后剩下%d条任务\n",len(tasks.Data)-count)
	// 删除重复任务
	Info.Println("删除重复任务")
	rawData := getRawData(duplicatedId)
	err=deleteDuplicateTasks(requst,header,rawData)
	if err!=nil{
		Error.Printf("删除重复任务出错\n")
		return
	}
	Info.Println("-------分割线--------")
	return
}

func deleteDuplicateTasks(requst *resty.Request,header map[string]string,rawData string)error {
	timeUnix := time.Now().UnixNano() / 1e6
	url := fmt.Sprintf("http://%s:5700/api/crons?t=%d",ip,timeUnix)
	resp, err := requst.
		SetHeaders(header).
		SetBody(rawData).
		Delete(url)
	if err!=nil{
		Error.Printf("[Net error]delete duplicate tasks failed, err:%v\n", err)
		return err
	}
	err = json.Unmarshal(resp.Body(), &trash)
	if err != nil {
		Error.Printf("json.Unmarshal failed, err:%v\n", err)
		return err
	}
	return nil
}
func getRawData(duplicatedId []string) string {
	rawData := "["
	for index, id := range duplicatedId {
		if id == "" {
			break
		}
		rawData += fmt.Sprintf("\"%s\"", id)
		if duplicatedId[index+1] != "" {
			rawData += ", "
		}
	}
	rawData += "]"
	return rawData
}

func getDuplicateTasks() ([]string,int){
	duplicatedId := make([]string, len(tasks.Data)+1)
	allNames := make(map[string]int, len(tasks.Data))
	count := 0
	for _, task := range tasks.Data {
		key := task.Name
		_, ok := allNames[key]
		if ok {
			duplicatedId[count] = task.Id
			count++
		} else {
			allNames[key] = 0
		}
	}
	return duplicatedId,count
}

func getList(requst *resty.Request,header map[string]string) error {
	timeUnix := time.Now().UnixNano() / 1e6
	url := fmt.Sprintf("http://%s:5700/api/crons?searchValue=&t=%d",ip,timeUnix)
	resp, err := requst.
		SetHeaders(header).
		EnableTrace().
		Get(url)
	if err!=nil{
		Error.Printf("[Net error]getList failed, err:%v\n", err)
		return err
	}
	err = json.Unmarshal(resp.Body(), &tasks)
	if err != nil {
		Error.Printf("json.Unmarshal failed, err:%v\n", err)
		return err
	}
	return nil
}

func getToken(requst *resty.Request, header map[string]string,data map[string]string) error{
	timeUnix := time.Now().UnixNano() / 1e6
	url := fmt.Sprintf("http://%s:5700/api/login?t=%d",ip,timeUnix)
	resp, err := requst.
		SetHeaders(header).
		SetBody(data).
		Post(url)
	if err!=nil{
		Error.Printf("[Net error]login failed, err:%v\n", err)
		return err
	}
	err = json.Unmarshal(resp.Body(), &login)
	if err != nil {
		Error.Printf("json.Unmarshal failed, err:%v\n", err)
		return err
	}
	return nil
}