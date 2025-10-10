package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Helper Function for writing Errors
func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error":  msg,
		"status": http.StatusText(code),
	})
}

// Helper Function for Making a clean path, so you can't get out off the root folder
func safeFolderPath(base, rel string) (string, error) {
	cleanPath := filepath.Clean(rel)
	cleanPath = strings.TrimPrefix(cleanPath, string(os.PathListSeparator))

	folderPath := filepath.Join(base, cleanPath)

	relPath, err := filepath.Rel(base, folderPath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return "", fmt.Errorf("invalid path")
	}

	return folderPath, nil
}
