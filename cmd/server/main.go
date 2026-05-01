package main

import (
	"fmt"
	"net/http"

	"fileTransfer/internal/handler"
)

func main() {
	// Serve static files from the web directory.
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	// API routes
	http.HandleFunc("/upload", handler.UploadHandler)
	http.HandleFunc("/download/", handler.DownloadHandler)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
