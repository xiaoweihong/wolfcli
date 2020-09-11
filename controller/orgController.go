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

func GetOrgs(token, orgid string) {
	post := httplib.Post(fmt.Sprintf("http://%s/proxy/loki/api/orgs/list", global.IP))
	post.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	req := struct {
		Limit  int64
		Offset int64
		OrgID  string
	}{
		0,
		0,
		orgid,
	}
	byteS, _ := json.Marshal(&req)
	response, _ := post.Body(byteS).Response()
	if response.StatusCode == http.StatusUnauthorized {
		TokenSave()
		zap.L().Error("token已经过期，请重新运行")
		return
	}
	s, err := post.Bytes()
	if err != nil {
		zap.L().Error("获取org错误", zap.Any("err", err))
	}
	result := struct {
		Total int64
		Rets  []model.Org
	}{}
	json.Unmarshal(s, &result)
	if result.Total == 0 {
		zap.L().Info("orgid不存在")
		return
	}
	for _, org := range result.Rets {
		fmt.Println(org.Name, org.SuperiorOrgId, org.OrgCode)
	}
}
