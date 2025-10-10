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
	path := r.URL.Query().Get("path")

	if path == "" {
		writeError(w, http.StatusBadRequest, "missing 'path' query parameter")
		return
	}

	folderPath, errPath := safeFolderPath(RootFolder, path)
	if errPath != nil {
		writeError(w, http.StatusForbidden, "invalid path")
		return
	}

	entries, err := os.ReadDir(folderPath)

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
	path := r.URL.Query().Get("path")

	if path == "" {
		writeError(w, http.StatusBadRequest, "missing 'path' query parameter")
		return
	}

	folderPath, errPath := safeFolderPath(RootFolder, path)
	if errPath != nil {
		writeError(w, http.StatusForbidden, "invalid path")
		return
	}

	info, err := os.Stat(folderPath)

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
		"path":     folderPath,
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

// Uploads a file
func handleUpload(w http.ResponseWriter, r *http.Request) {
	// Checking that only the POST Method is used
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed, use POST")
		return
	}

	path := r.FormValue("path")

	if path == "" {
		writeError(w, http.StatusBadRequest, "missing 'path' query parameter")
		return
	}

	folderPath, errPath := safeFolderPath(RootFolder, path)
	if errPath != nil {
		writeError(w, http.StatusForbidden, "invalid path")
		return
	}

	// Parsing
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("failed to parse form: %v", err))
		return
	}

	// Getting the file to upload
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("failed to get file: %v", err))
		return
	}
	defer file.Close()

	// Creating the full path for the file to get copied into
	fullPath := filepath.Join(folderPath, header.Filename)

	// Check if directory exists, Create directory if needed
	if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create directory: %v", err))
		return
	}

	// Prevent overwriting
	if _, err := os.Stat(fullPath); err == nil {
		writeError(w, http.StatusConflict, "file already exists")
		return
	}

	// Creating a the uploaded file
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

// Downloads a file
func handleDownload(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		writeError(w, http.StatusBadRequest, "missing 'path' query parameter")
		return
	}

	fullPath, errPath := safeFolderPath(RootFolder, path)
	if errPath != nil {
		writeError(w, http.StatusForbidden, "invalid path")
		return
	}

	// Open the file
	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			writeError(w, http.StatusNotFound, "file not found")
		} else {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to open file: %v", err))
		}
		return
	}
	defer file.Close()

	// Set headers for download
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(fullPath))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Serve the file content
	if _, err := io.Copy(w, file); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to send file: %v", err))
		return
	}
}

// Deletes a file
func handleDelete(w http.ResponseWriter, r *http.Request) {
	// Checking that only the POST Method is used
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed, use POST")
		return
	}

	path := r.FormValue("path")

	if path == "" {
		writeError(w, http.StatusBadRequest, "missing 'path' query parameter")
		return
	}

	fullPath, errPath := safeFolderPath(RootFolder, path)
	if errPath != nil {
		writeError(w, http.StatusForbidden, "invalid path")
		return
	}

	_, err := os.Stat(fullPath)
	if err != nil {
		writeError(w, http.StatusNotFound, "path does not exist")
	}

	info, err := os.Stat(fullPath)

	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("file or folder not found: %v", err))
		return
	}

	if info.IsDir() {
		err := os.RemoveAll(fullPath)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("Path Error: %v", err))
			return
		}
	} else {
		err := os.Remove(fullPath)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("Path Error: %v", err))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "deleted successfully",
		"path":    fullPath,
	})
}
