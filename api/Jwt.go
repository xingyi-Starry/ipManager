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

type QueryDevice struct {
	ID       int    `json:"id"`
	IP       string `json:"ip"`
	LoggedIn bool   `json:"logged_in"`
}

type QueryResp struct {
	Status  string        `json:"status"`
	Devices []QueryDevice `json:"devices"`
}

type ErrResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
