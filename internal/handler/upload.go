package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	path := "uploads/" + id + "_" + header.Filename

	if err := os.MkdirAll("uploads", 0o755); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(path)
	if err != nil {
		http.Error(w, "Failed to create file on server", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Uploaded! <br>Download: <a href='/download/%s'>Link</a>", id)
}
