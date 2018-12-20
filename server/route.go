package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
)

// Handler 消息处理
type Handler interface {
	Serve(data []byte, w io.Writer) error
}

// HandleFunc 消息回调
type HandleFunc func(data []byte, w io.Writer) error

// Serve 处理消息
func (f HandleFunc) Serve(data []byte, w io.Writer) error {
	return f(data, w)
}

// Route 路由
type Route map[uint16]Handler

// Serve 处理消息
func (r Route) Serve(data []byte, w io.Writer) error {
	var id uint16
	m := bufio.NewScanner(bytes.NewReader(data))
	m.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		n := binary.BigEndian.Uint16(data[:2])
		if atEOF || int(n+2) > len(data) || n < 2 {
			return len(data), nil, nil
		}
		id = binary.BigEndian.Uint16(data[2:4])
		return int(n + 2), data[4 : n+2], nil
	})
	for m.Scan() {
		if h, exist := r[id]; !exist {
			log.Printf("recv %s unkwnon 0x%04x: %s", w, id, hex.EncodeToString(m.Bytes()))
		} else if err := h.Serve(m.Bytes(), w); err != nil {
			log.Printf("recv %s msg 0x%04x failed: %v", w, id, err)
		} else if id > 0 {
			log.Printf("recv %s msg 0x%04x: %d", w, id, len(m.Bytes()))
		}
	}
	return m.Err()
}

// Handle 处理消息
func (r Route) Handle(id uint16, f HandleFunc) {
	r[id] = f
}
