package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
		writeError(w, http.StatusBadRequest, "missing 'path' query parameter")
		return
	}

	fullPath := RootFolder + path

	entries, err := os.ReadDir(fullPath)

	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("directory not found or inaccessible: %v", err))
		return
	}

	fileList := make([]map[string]interface{}, 0)

	for i := 0; i < len(entries); i++ {
		entry := entries[i]

		info, err := entry.Info()

		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to read file info for '%s': %v", entry.Name(), err))
			return
		}

		fileName := entry.Name()

		var fileType string
		if entry.IsDir() {
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
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to convert into JSON: %v", err))
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
		writeError(w, http.StatusBadRequest, "missing 'path' query parameter")
		return
	}

	fullPath := RootFolder + path

	info, err := os.Stat(fullPath)

	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("file or folder not found: %v", err))
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
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to convert into JSON: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonFile))
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed, use POST")
		return
	}

	query := r.URL.Query()
	path := query.Get("path")

	if path == "" {
		writeError(w, http.StatusBadRequest, "missing 'path' query parameter")
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("failed to parse form: %v", err))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("failed to get file: %v", err))
		return
	}
	defer file.Close()

	fullPath := filepath.Join(RootFolder, path, header.Filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create file: %v", err))
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to save file: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "file uploaded successfully",
		"filename": header.Filename,
		"size":     header.Size,
	})

}
