package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"

	"fileTransfer/internal/handler"
	"fileTransfer/internal/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found")
	}
	if err := storage.InitR2(); err != nil {
		panic(err)
	}

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	http.HandleFunc("/upload", handler.UploadHandler)
	http.HandleFunc("/download", handler.DownloadHandler)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
