package main

import (
	"fmt"
	"ipManager/api"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/bind", api.NewBind).Methods("POST")
	router.HandleFunc("/verify", api.VerifyBind).Methods("POST")
	router.HandleFunc("/devices/{id}/login", api.Login).Methods("POST")
	router.HandleFunc("/devices/{id}/logout", api.Logout).Methods("POST")

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", router)
}
