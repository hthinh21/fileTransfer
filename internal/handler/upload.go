package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"

	"fileTransfer/internal/sharecode"
	"fileTransfer/internal/storage"
)

type UploadHandler struct {
	store *storage.RedisStore
}

func NewUploadHandler(store *storage.RedisStore) *UploadHandler {
	return &UploadHandler{store: store}
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large (max 10MB)", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	id := uuid.New().String()
	objectKey := filepath.ToSlash(filepath.Join("uploads", id+"_"+header.Filename))

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	if err = storage.UploadFile(context.Background(), objectKey, file, contentType); err != nil {
		log.Printf("upload to R2 failed: key=%s file=%s err=%v", objectKey, header.Filename, err)
		http.Error(w, "Upload failed, please try again", http.StatusInternalServerError)
		return
	}

	var code string
	for attempt := 0; attempt < 5; attempt++ {
		candidate := sharecode.Generate()

		ok, err := h.store.SaveIfAbsent(candidate, storage.FileMeta{
			ObjectKey: objectKey,
			FileName:  header.Filename,
		})
		if err != nil {
			log.Printf("save code failed: code=%s key=%s err=%v", candidate, objectKey, err)
			http.Error(w, "Upload failed, please try again", http.StatusInternalServerError)
			return
		}

		if ok {
			code = candidate
			break
		}
	}

	if code == "" {
		log.Printf("could not generate unique code: key=%s", objectKey)
		http.Error(w, "Upload failed, please try again", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"code": code})
}
