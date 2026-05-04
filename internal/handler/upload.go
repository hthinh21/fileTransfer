package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"fileTransfer/internal/storage"
	"fileTransfer/internal/store"
	"fileTransfer/internal/utils"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	// limit size
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File too large (max 10MB)", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	objectKey := filepath.ToSlash(filepath.Join("uploads", id+"_"+header.Filename))

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	err = storage.UploadFile(context.Background(), objectKey, file, contentType)
	if err != nil {
		log.Printf("upload to R2 failed: key=%s file=%s err=%v", objectKey, header.Filename, err)
		http.Error(w, fmt.Sprintf("Cannot upload file: %v", err), http.StatusInternalServerError)
		return
	}

	code := utils.GenerateCode()

	rs := store.NewRedisStore()
	err = rs.Save(code, store.FileMeta{
		ObjectKey: objectKey,
		FileName:  header.Filename,
	})
	if err != nil {
		log.Printf("save upload code failed: code=%s key=%s err=%v", code, objectKey, err)
		http.Error(w, fmt.Sprintf("Redis error: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Your code: <b>%s</b>", code)
}
