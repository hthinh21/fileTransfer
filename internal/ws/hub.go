package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu  sync.Mutex
	conns map[string]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{conns: make(map[string]*websocket.Conn)}
}

func (h *Hub) Register(code string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.conns[code] = conn
}

func (h *Hub) Unregister(code string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.conns, code)
}

// take lấy kết nối ứng với code (nếu có) và xoá khỏi hub trong cùng 1 lần khoá
func (h *Hub) take(code string) *websocket.Conn {
	h.mu.Lock()
	defer h.mu.Unlock()

	conn := h.conns[code]
	delete(h.conns, code)
	return conn
}

func (h *Hub) NotifyDownloaded(code string) {
	conn := h.take(code)
	if conn == nil {
		return
	}

	_ = conn.WriteJSON(map[string]string{"event": "downloaded"})
	_ = conn.Close()
}

// Expire đóng kết nối khi code đã hết hạn mà chưa được tải (tránh treo kết nối vô thời hạn)
func (h *Hub) Expire(code string) {
	conn := h.take(code)
	if conn == nil {
		return
	}

	_ = conn.Close()
}
