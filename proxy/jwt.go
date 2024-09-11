package proxy

import "encoding/json"

type LoginReq struct {
	Username string `json:"username"`
	IP       string `json:"ip"`
}

type LoginResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (l *LoginResp) ResolveJwt(bodyText []byte) error {
	err := json.Unmarshal(bodyText, l)
	if err != nil {
		return err
	}
	return nil
}

type LogoutReq struct {
	Username string `json:"username"`
	IP       string `json:"ip"`
}

type LogoutResp struct {
	Success bool `json:"success"`
}

func (l *LogoutResp) ResolveJwt(bodyText []byte) error {
	err := json.Unmarshal(bodyText, l)
	if err != nil {
		return err
	}
	return nil
}
