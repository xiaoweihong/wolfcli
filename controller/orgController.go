package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"go.uber.org/zap"
	"net/http"
	"strings"
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
		if org.Id == "0000" {
			fmt.Printf("组织id:%v\t\t\t\t\t组织名称:%v\n", org.Id, org.Name)
		} else {
			fmt.Printf("组织id:%v\t组织名称:%v\n", org.Id, org.Name)
		}
		//fmt.Printf("组织名称:%v",org.Id)
		//fmt.Println(org.Id, org.Name, org.SuperiorOrgId, org.OrgCode)
	}
}

func AddOrg(token, orgName, SuperiorOrgId string) {
	if strings.TrimSpace(orgName) == "" {
		zap.L().Error("组织名称不能为空")
		return
	}
	post := httplib.Post(fmt.Sprintf("http://%s/proxy/loki/api/org", global.IP))
	post.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	orgReq := make(map[string]string)
	orgReq["Name"] = orgName
	orgReq["SuperiorOrgId"] = SuperiorOrgId
	orgReqM, _ := json.Marshal(&orgReq)
	body := post.Body(orgReqM)
	response, err := body.Response()
	if err != nil {
		fmt.Println(err)
		return
	}
	if response.StatusCode == http.StatusNoContent {
		zap.L().Info("组织添加成功", zap.String("组织名称", orgName))
	} else {
		zap.Error(errors.New("未知错误"))
		fmt.Println(response)
		return
	}
}
