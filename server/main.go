package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Lists the data of a directory
func handleList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	path := query.Get("path")

	if path == "" {
		http.Error(w, "No Path Available", http.StatusBadRequest)
		return
	}

	//test file list
	jsonData := `[
        { "name": "file1.txt", "type": "file", "size": 1234, "modified": "2025-10-08T12:34:56Z" },
        { "name": "docs", "type": "folder", "size": 0, "modified": "2025-10-01T08:00:00Z" }
    ]`

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
}

// Lists the metadata of a file
func handleStat(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	path := query.Get("path")

	if path == "" {
		http.Error(w, "No Path Available", http.StatusBadRequest)
		return
	}

	//test file data
	jsonFile := `{
		"path": "/some/path/file1.txt",
		"type": "file",
		"size": 1234,
		"modified": "2025-10-08T12:34:56Z"
	}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonFile))
}

func main() {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := Response{Message: "API Response", Status: "Success"}
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/list", handleList)

	http.HandleFunc("/stat", handleStat)

	fmt.Println("Starting server at port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
