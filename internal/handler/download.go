package handler

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"fileTransfer/internal/storage"
	"fileTransfer/internal/store"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	rs := store.NewRedisStore()
	meta, err := rs.Get(code)
	if err != nil {
		http.Error(w, "Invalid or expired code", http.StatusNotFound)
		return
	}

	reader, contentType, contentLength, err := storage.DownloadFile(context.Background(), meta.ObjectKey)
	if err != nil {
		http.Error(w, "Cannot download file", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	_ = rs.Delete(code)

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+meta.FileName+`"`)
	if contentLength > 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	}

	if _, err := io.Copy(w, reader); err != nil {
		http.Error(w, "Cannot stream file", http.StatusInternalServerError)
		return
	}
}
