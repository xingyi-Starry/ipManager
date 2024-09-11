package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func HttpRequest(method string, url string, body string, header string) ([]byte, error) {
	client := &http.Client{}
	var data io.Reader = nil
	if method == "POST" {
		data = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if header != "" {
		req.Header.Set("Authorization", "Bearer "+header)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyText, nil
}

func Login(username, IP string) error {
	// 1. 构造请求体
	/*示例请求体：
	{
		"username": "username",
		"ip": "ip",
	}
	*/
	reqBody := LoginReq{
		Username: username,
		IP:       IP,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// 2. 发送请求
	respBody, err := HttpRequest("POST", "http://yxms.byr.ink/api/login", string(body), "")
	if err != nil {
		return err
	}

	// 3. 处理响应
	var resp LoginResp
	err = resp.ResolveJwt(respBody)
	if err != nil {
		return err
	}
	success := resp.Success
	if success {
		return nil
	} else {
		err = errors.New(resp.Message)
		return err
	}

}

func Logout(username, IP string) error {
	// 1. 构造请求体
	/*示例请求体：
	{
		"username": "username",
		"ip": "ip",
	}
	*/
	reqBody := LogoutReq{
		Username: username,
		IP:       IP,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// 2. 发送请求
	respBody, err := HttpRequest("POST", "http://yxms.byr.ink/api/logout", string(body), "")
	if err != nil {
		return err
	}

	// 3. 处理响应
	var resp LogoutResp
	err = resp.ResolveJwt(respBody)
	if err != nil {
		return err
	}
	success := resp.Success
	if success {
		return nil
	} else {
		err = fmt.Errorf("Logout failed")
		return err
	}

}
