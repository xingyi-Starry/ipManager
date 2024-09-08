package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"ipManager/device"
	"net"
	"net/http"
	"strconv"
)

var devices = make(map[string]device.Device)

func genToken(length int) (string, error) {
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func NewBind(w http.ResponseWriter, r *http.Request) {
	// generate token
	token, err := genToken(16)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	device := device.Device{Token: token, ID: strconv.Itoa(len(devices))}
	devices[token] = device

	// log
	fmt.Println("New device: ", device)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
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
	var tkData struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(bodyText, &tkData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	token := tkData.Token

	if device, ok := devices[token]; ok {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		device.IP = ip
		devices[token] = device

		// log
		fmt.Println("Device verified: ", device)

		w.Header().Set("Content-Type", "application/json")
		resp := struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{
			Status:  "success",
			Message: "Device bound successfully",
		}
		json.NewEncoder(w).Encode(resp)
	} else {
		resp := struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{
			Status:  "error",
			Message: "Invalid or expired token",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
}
