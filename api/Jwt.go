package api

type BindReq struct {
	Username string `json:"username"`
}

type BindResp struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

type VerifyReq struct {
	Token string `json:"token"`
}

type VerifyResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type DeleteResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type LoginReq struct {
	Username string `json:"username"`
}

type LoginResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type LogoutResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ErrResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
