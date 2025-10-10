package main

import (
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Entry point of the Server
func main() {
	http.HandleFunc("/api", handleAPI)

	http.HandleFunc("/list", handleList)

	http.HandleFunc("/stat", handleStat)

	http.HandleFunc("/upload", handleUpload)

	http.HandleFunc("/download", handleDownload)

	fmt.Println("Starting server at port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
