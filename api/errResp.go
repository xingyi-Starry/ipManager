package api

import (
	"encoding/json"
	"net/http"
)

func SendErrResp(w http.ResponseWriter, message string) {
	resp := ErrResp{
		Status:  "error",
		Message: message,
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(respBytes)
}
