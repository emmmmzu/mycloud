package main

import (
	"encoding/json"
	"net/http"
	"os"
)

// Checks if the API works
func handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "API Response", Status: "Success"}
	json.NewEncoder(w).Encode(response)
}

// Lists the data of a directory
func handleList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	path := query.Get("path")

	if path == "" {
		http.Error(w, "No Path Available", http.StatusBadRequest)
		return
	}

	fullPath := RootFolder + path

	entries, err := os.ReadDir(fullPath)

	if err != nil {
		http.Error(w, "Path Not Found", http.StatusBadRequest)
		return
	}

	fileList := make([]map[string]interface{}, 0)

	for i := 0; i < len(entries); i++ {
		info, err := entries[i].Info()

		if err != nil {
			return
		}

		fileName := entries[i].Name()

		var fileType string
		if entries[i].IsDir() {
			fileType = "folder"
		} else {
			fileType = "file"
		}

		fileSize := info.Size()

		fileModified := info.ModTime().Format(TimeFormat)

		fileObj := map[string]interface{}{
			"name":     fileName,
			"type":     fileType,
			"size":     fileSize,
			"modified": fileModified,
		}

		fileList = append(fileList, fileObj)

	}

	jsonData, err := json.Marshal(fileList)

	if err != nil {
		return
	}

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

	fullPath := RootFolder + path

	info, err := os.Stat(fullPath)

	if err != nil {
		return
	}

	var fileType string

	if info.IsDir() {
		fileType = "folder"
	} else {
		fileType = "file"
	}

	fileSize := info.Size()

	fileModified := info.ModTime().Format(TimeFormat)

	fileObj := map[string]interface{}{
		"path":     fullPath,
		"type":     fileType,
		"size":     fileSize,
		"modified": fileModified,
	}

	jsonFile, err := json.Marshal(fileObj)

	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonFile))
}
