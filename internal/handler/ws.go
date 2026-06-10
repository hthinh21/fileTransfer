package handler

import (
	"net/http"
	"os"
	"time"
	

	"github.com/gorilla/websocket"

	"fileTransfer/internal/storage"
	"fileTransfer/internal/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: checkOrigin,
}

// checkOrigin chỉ chấp nhận kết nối WS từ domain được cấu hình qua biến môi trường ALLOWED_ORIGIN
// (mặc định là localhost:8080 cho môi trường dev). Khi deploy, chỉ cần đổi giá trị biến này.
func checkOrigin(r *http.Request) bool {
	allowed := os.Getenv("ALLOWED_ORIGIN")
	if allowed == "" {
		allowed = os.DevNull
	}
	return r.Header.Get("Origin") == allowed
}

type WSHandler struct {
	hub *ws.Hub
}

func NewWSHandler(hub *ws.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	h.hub.Register(code, conn)

	// tự đóng kết nối khi code hết hạn (TTL trùng với TTL của Redis) mà chưa có ai tải
	timer := time.AfterFunc(storage.CodeTTL, func() {
		h.hub.Expire(code)
	})

	// giữ kết nối mở & thoát khi client đóng kết nối
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			timer.Stop()
			h.hub.Unregister(code)
			conn.Close()
			return
		}
	}
}
