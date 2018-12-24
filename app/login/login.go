package login

import (
	"bufio"
	"encoding/hex"

	"github.com/dragonflylee/go-saga/network"
	"github.com/golang/glog"
)

// Version 版本号请求
func Version(data []byte, c *network.Conn) error {

	glog.Infof("login %s version %s", c, hex.EncodeToString(data))

	w := bufio.NewWriter(c)
	// SSMG_VERSION_ACK
	network.Packet(w, 0x0A, 0x0002, data)
	// SSMG_LOGIN_ALLOWED
	network.Packet(w, 0x0A, 0x001E, data)
	return w.Flush()
}

// Ping 心跳
func Ping(data []byte, c *network.Conn) error {
	// SSMG_PONG
	w := bufio.NewWriter(c)
	network.Packet(w, 0x06, 0x000B)
	return w.Flush()
}

// Login 登录
func Login(data []byte, c *network.Conn) error {
	var m struct {
		Name string
		Pass string
	}
	network.Unmarshal(data, &m)

	glog.Infof("login %s user %s pass %s", c, m.Name, m.Pass)
	w := bufio.NewWriter(c)
	// SSMG_LOGIN_ACK
	network.Packet(w, 0x13, 0x0020)
	SendCharData(w)
	return w.Flush()
}

// RequestMap 请求地图
func RequestMap(data []byte, c *network.Conn) error {
	// SSMG_SEND_TO_MAP_SERVER
	w := bufio.NewWriter(c)
	network.Packet(w, 0xb7, 0x0033, byte(0x10), "127.0.0.1", uint32(12001))
	return w.Flush()
}
