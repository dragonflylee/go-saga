package network

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"

	"github.com/golang/glog"
)

// Handler 消息处理
type Handler interface {
	Serve(data []byte, c *Conn) error
}

// HandleFunc 消息回调
type HandleFunc func(data []byte, c *Conn) error

// Serve 处理消息
func (f HandleFunc) Serve(data []byte, c *Conn) error {
	return f(data, c)
}

// Route 路由
type Route map[uint16]Handler

// Serve 处理消息
func (r Route) Serve(data []byte, c *Conn) error {
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
			glog.Warningf("recv %s unkwnon 0x%04x: %s", c, id, hex.EncodeToString(m.Bytes()))
		} else if err := h.Serve(m.Bytes(), c); err != nil {
			glog.Warningf("recv %s msg 0x%04x failed: %v", c, id, err)
		} else if id > 0 {
			// glog.Infof("recv %s msg 0x%04x: %d", c, id, len(m.Bytes()))
		}
	}
	return m.Err()
}

// Handle 处理消息
func (r Route) Handle(id uint16, f HandleFunc) {
	r[id] = f
}

// TODO 空任务
func TODO(data []byte, c io.Writer) error {
	return nil
}
