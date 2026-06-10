package handler

import (
	"context"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"

	"fileTransfer/internal/storage"
	"fileTransfer/internal/ws"
)

type DownloadHandler struct {
	store *storage.RedisStore
	hub   *ws.Hub
}

func NewDownloadHandler(store *storage.RedisStore, hub *ws.Hub) *DownloadHandler {
	return &DownloadHandler{store: store, hub: hub}
}

func (h *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	meta, err := h.store.Get(code)
	if err != nil {
		http.Error(w, "Invalid or expired code", http.StatusNotFound)
		return
	}

	if meta.Downloaded {
		http.Error(w, "Invalid or expired code", http.StatusNotFound)
		return
	}

	reader, contentType, contentLength, err := storage.DownloadFile(context.Background(), meta.ObjectKey)
	if err != nil {
		log.Printf("download from R2 failed: key=%s err=%v", meta.ObjectKey, err)
		http.Error(w, "Cannot download file", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	meta.Downloaded = true
	meta.DownloadedAt = time.Now()
	h.store.Save(code, meta)

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition",
		mime.FormatMediaType("attachment", map[string]string{"filename": meta.FileName}))
	if contentLength > 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	}
	_, err = io.Copy(w, reader)
	if err == nil {
		h.hub.NotifyDownloaded(code)
	}

	h.store.Delete(code)
}
