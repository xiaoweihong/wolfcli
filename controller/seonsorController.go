package controller

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"go.uber.org/zap"
	"net/http"
	"wolfcli/global"
	"wolfcli/model"
)

// GetSensors 获取设备列表
func GetSensors(token string, orgid string) {
	get := httplib.Get(fmt.Sprintf("http://%s/proxy/loki/api/sensors", global.IP))
	get.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	response, _ := get.Param("orgid", orgid).Response()
	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		TokenSave()
		zap.L().Error("token已经过期，请重新运行")
		return
	}
	s, err := get.Bytes()
	if err != nil {
		zap.L().Error("获取sensor错误", zap.Any("err", err))
		return
	}
	result := struct {
		Total int64
		Rets  []model.Sensor
	}{}
	json.Unmarshal(s, &result)
	for _, sensor := range result.Rets {
		fmt.Println(sensor)
	}
}
