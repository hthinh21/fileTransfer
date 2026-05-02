package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"fileTransfer/internal/store"
	"fileTransfer/internal/utils"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	// limit size
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File too large (max 10MB)", 400)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", 400)
		return
	}
	defer file.Close()

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	path := "uploads/" + id + "_" + header.Filename

	dst, err := os.Create(path)
	if err != nil {
		http.Error(w, "Cannot save file", 500)
		return
	}
	defer dst.Close()

	io.Copy(dst, file)

	// generate code
	code := utils.GenerateCode()

	// save to Redis
	rs := store.NewRedisStore()
	err = rs.Save(code, path)
	if err != nil {
		http.Error(w, "Redis error", 500)
		return
	}

	fmt.Fprintf(w, "Your code: <b>%s</b>", code)
}
