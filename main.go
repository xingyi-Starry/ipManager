package main

import (
	"fmt"
	"ipManager/api"
	"net/http"
)

func main() {
	http.HandleFunc("/bind", api.NewBind)
	http.HandleFunc("/verify", api.VerifyBind)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
