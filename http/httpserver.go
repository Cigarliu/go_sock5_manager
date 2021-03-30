package http

import (
	"errors"

	"net/http"
	"time"
)

func LoginSys(user, pass string) (string, interface{}) {
	var token string
	url := "ssr.comeboy.cn?"
	kUser := "user="
	kPass := "pass="
	reqUrl := url + kUser + user + "&" + kPass + pass
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(reqUrl)
	if err != nil {
		return token, errors.New("请求验证帐号密码时发生错误")
	}
	defer resp.Body.Close()

	return token, nil
}
