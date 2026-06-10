package main

import (
    "log"
    "net/http"

    "github.com/joho/godotenv"

    "fileTransfer/internal/handler"
    "fileTransfer/internal/storage"
    "fileTransfer/internal/ws"
)

func main() {
    _ = godotenv.Load()

    if err := storage.InitR2(); err != nil {
        log.Fatalf("R2 init failed: %v", err)
    }

    store, err := storage.NewRedisStore()
    if err != nil {
        log.Fatalf("Redis init failed: %v", err)
    }

    hub := ws.NewHub()

    mux := http.NewServeMux()
    mux.Handle("/", http.FileServer(http.Dir("./web")))
    mux.Handle("/upload", handler.NewUploadHandler(store))
    mux.Handle("/download", handler.NewDownloadHandler(store, hub))
    mux.Handle("/ws", handler.NewWSHandler(hub))

    log.Println("Server running on http://localhost:8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}