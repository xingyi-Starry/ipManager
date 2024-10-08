package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"ipManager/device"
	"ipManager/proxy"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var tokens = make(map[string]string) // token -> username
var dm = device.NewDeviceManager()

func genToken(length int) (string, error) {
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func NewBind(w http.ResponseWriter, r *http.Request) {
	// 从请求体中解析出 username
	/*示例请求体：
	{
		"username": "username",
	}
	*/
	bodyText, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var username BindReq
	err = json.Unmarshal(bodyText, &username)
	if err != nil {
		SendErrResp(w, "Invalid request body")
		return
	}

	// generate token
	genSuccess := false
	token := ""
	for !genSuccess {
		token, err = genToken(3)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if _, exists := tokens[token]; !exists {
			genSuccess = true
		}
	}
	tokens[token] = username.Username

	// log
	fmt.Printf("%s requested a new token: %s\n", username.Username, token)

	w.Header().Set("Content-Type", "application/json")
	resp := BindResp{
		Status: "success",
		Token:  token,
	}
	json.NewEncoder(w).Encode(resp)
}

func VerifyBind(w http.ResponseWriter, r *http.Request) {
	// 从请求体中解析出token
	/*示例请求体：
		{
	  		"token": "abc123xyz"
		}
	*/
	bodyText, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var tkData VerifyReq
	err = json.Unmarshal(bodyText, &tkData)
	if err != nil {
		SendErrResp(w, "Invalid request body")
		return
	}
	token := tkData.Token

	if username, exists := tokens[token]; exists {
		// 创建设备和绑定
		ip := r.RemoteAddr
		ip, _, err := net.SplitHostPort(ip)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		d := device.NewDevice(username, len(dm.Devices)+1, ip)
		err = dm.AddDevice(d)
		if err != nil {
			SendErrResp(w, err.Error())
			return
		}
		delete(tokens, token)

		// log
		fmt.Println("Device bound: ", d)

		w.Header().Set("Content-Type", "application/json")
		resp := VerifyResp{
			Status:  "success",
			Message: "Device bound successfully",
		}
		json.NewEncoder(w).Encode(resp)
	} else {
		SendErrResp(w, "Invalid or expired token")
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	// 从请求体中解析出 username
	/*示例请求体：
	{
		"username": "username",
	}
	*/
	bodyText, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println("api.Login:", err)
		return
	}
	var loginData LoginReq
	err = json.Unmarshal(bodyText, &loginData)
	if err != nil {
		SendErrResp(w, "Invalid request body")
		fmt.Println("api.Login:", err)
		return
	}
	username := loginData.Username

	// 从 URL 中解析出设备 ID
	id_raw := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_raw)
	if err != nil {
		SendErrResp(w, "Invalid device ID")
		fmt.Println("api.Login:", err)
		return
	}

	// 验证设备是否符合username
	if device, exists := dm.GetDeviceByID(id); exists {
		if device.Username == username {
			if !device.Logged_in { // 设备未登录
				err = proxy.Login(username, device.IP)
				if err != nil { // 登录失败
					SendErrResp(w, err.Error())
					fmt.Println("api.Login:", err)
				} else { // 登录成功
					device.Logged_in = true
					w.Header().Set("Content-Type", "application/json")
					resp := LoginResp{
						Status:  "success",
						Message: "Device logged in successfully",
					}
					json.NewEncoder(w).Encode(resp)
					fmt.Println("Device logged in: ", device)
				}
			} else { // 设备已登录
				SendErrResp(w, "Device already logged in")
				fmt.Println("api.Login: Device already logged in")
			}
		} else {
			SendErrResp(w, "Invalid username")
			fmt.Println("api.Login: Invalid username")
		}
	} else {
		SendErrResp(w, "Device not found")
		fmt.Println("api.Login: Device not found")
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// 从请求体中解析出 username
	/*示例请求体：
	{
		"username": "username",
	}
	*/
	bodyText, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println("api.Logout:", err)
		return
	}
	var loginData LoginReq
	err = json.Unmarshal(bodyText, &loginData)
	if err != nil {
		SendErrResp(w, "Invalid request body")
		fmt.Println("api.Logout:", err)
		return
	}
	username := loginData.Username

	// 从 URL 中解析出设备 ID
	id_raw := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_raw)
	if err != nil {
		SendErrResp(w, "Invalid device ID")
		fmt.Println("api.Logout:", err)
		return
	}

	// 验证设备是否符合username
	if device, exists := dm.GetDeviceByID(id); exists {
		if device.Username == username {
			if device.Logged_in {
				device.Logged_in = false
				w.Header().Set("Content-Type", "application/json")
				resp := LoginResp{
					Status:  "success",
					Message: "Device logged out successfully",
				}
				json.NewEncoder(w).Encode(resp)
				fmt.Println("Device logged out: ", device)
			} else {
				SendErrResp(w, "Device already logged out")
				fmt.Println("api.Logout: Device already logged out")
			}
		} else {
			SendErrResp(w, "Invalid username")
			fmt.Println("api.Logout: Invalid username")
		}
	} else {
		SendErrResp(w, "Device not found")
		fmt.Println("api.Logout: Device not found")
	}
}

func QueryDevices(w http.ResponseWriter, r *http.Request) { // 查询所有设备
	devices := make([]QueryDevice, 0)
	for _, device := range dm.Devices {
		devices = append(devices, QueryDevice{
			ID:       device.ID,
			IP:       device.IP,
			LoggedIn: device.Logged_in,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	resp := QueryResp{
		Status:  "success",
		Devices: devices,
	}
	json.NewEncoder(w).Encode(resp)
}
