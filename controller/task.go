package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
	"wolfcli/global"
	"wolfcli/model"
)

type TaskResult struct {
	SensorId     string
	DetTypesInfo []int
	Config       model.SensorConfig
	Type         int
	SensorName   string
	Uts          time.Time
}

func (t *TaskResult) String() string {
	return ""
}

func TaskStatusSave(token string) {
	var result = struct {
		Total int
		Rets  []TaskResult
	}{}
	post := httplib.Get(fmt.Sprintf("http://%s/proxy/loki/api/tasks?ResultType=1&TaskType=3&Limit=1000&Offset=0", global.IP))
	post.Header("Authorization", fmt.Sprintf("Bearer %s", token))

	response, err := post.Response()
	if err != nil {
		zap.L().Error("请求失败", zap.Error(err))
		return
	}
	if response.StatusCode == http.StatusUnauthorized {
		TokenSave()
		zap.L().Error("token已经过期，请重新运行")
		return
	}
	resBytes, err := post.Bytes()
	if err != nil {
		zap.L().Error("获取任务列表失败", zap.Error(err))
		return
	}
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		zap.L().Error("序列化失败", zap.Error(err))
		return
	}
	var str bytes.Buffer
	err = json.Indent(&str, resBytes, "", " ")
	if err != nil {
		zap.L().Error("格式化失败", zap.Error(err))
		return
	}
	err = ioutil.WriteFile("task.txt", str.Bytes(), 0755)
	if err != nil {
		fmt.Println(err)
	}

}

func TaskResume(token string) {
	post := httplib.Post(fmt.Sprintf("http://%s/proxy/loki/api/tasks", global.IP))
	post.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	var taskReq = struct {
		Total int
		Rets  []TaskResult
	}{}

	file, err := ioutil.ReadFile("task.txt")
	if err != nil {
		zap.L().Error("打开文件失败", zap.Error(err))
		return
	}
	json.Unmarshal(file, &taskReq)
	byteS, _ := json.Marshal(&taskReq.Rets)
	response, err := post.Body(byteS).Response()
	if err != nil {
		zap.L().Error("请求失败", zap.Error(err))
		return
	}
	if response.StatusCode == http.StatusUnauthorized {
		TokenSave()
		zap.L().Error("token已经过期，请重新运行")
		return
	}
	bytes, err := post.Bytes()
	if err != nil {
		zap.L().Error("请求失败", zap.Error(err))
	}
	zap.S().Info(string(bytes))
}
