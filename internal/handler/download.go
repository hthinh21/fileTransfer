package handler

import (
	"net/http"
	"os"
	"strings"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/download/")
	if id == "" {
		http.Error(w, "File ID is required", http.StatusBadRequest)
		return
	}

	files, err := os.ReadDir("uploads")
	if err != nil {
		http.Error(w, "Upload directory not found", http.StatusNotFound)
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), id) {
			http.ServeFile(w, r, "uploads/"+file.Name())
			return
		}
	}

	http.Error(w, "File not found", http.StatusNotFound)
}
