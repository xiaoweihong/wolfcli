package controller

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"wolfcli/global"
)

const tokenfile = ".token"

func ProductToken() string {
	post := httplib.Post(fmt.Sprintf("http://%s:3000/login", global.IP))
	account := make(map[string]string)
	account["username"] = "admin"
	account["password"] = "admin@2013"
	accountM, _ := json.Marshal(&account)
	body := post.Body(accountM)
	result := make(map[string]string)
	result["token"] = ""
	bytes, _ := body.Bytes()
	_ = json.Unmarshal(bytes, &result)

	list.New()
	return result["token"]
}

func GetToken() string {
	if !CheckTokenFile() {
		return TokenSave()
	} else {
		tokenByte, _ := ioutil.ReadFile(tokenfile)
		return string(tokenByte)
	}
}

func CheckTokenFile() bool {
	_, err := os.Stat(tokenfile)
	if err != nil {
		zap.L().Info("token文件不存在，创建文件")
		return false
	}
	return true
}

func TokenSave() string {
	token := ProductToken()
	err := ioutil.WriteFile(tokenfile, []byte(token), 0644)
	if err != nil {
		zap.L().Error("token文件创建失败")
		return ""
	}
	return token
}
