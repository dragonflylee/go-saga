package play

import (
	"bufio"
	"encoding/hex"

	"github.com/dragonflylee/go-saga/model"
	"github.com/dragonflylee/go-saga/network"
	"github.com/golang/glog"
)

// Version 版本号请求
func Version(data []byte, c *network.Conn) error {
	glog.Infof("map %s version %s", c, hex.EncodeToString(data))

	w := bufio.NewWriter(c)
	// SSMG_VERSION_ACK
	network.Packet(w, 0x0A, 0x000B, data)
	// SSMG_LOGIN_ALLOWED
	network.Packet(w, 0x0A, 0x000F, data)
	return w.Flush()
}

// Ping 心跳
func Ping(data []byte, c *network.Conn) error {
	// SSMG_PONG
	w := bufio.NewWriter(c)
	network.Packet(w, 0x02, 0x0033)
	return w.Flush()
}

// Login 登录
func Login(data []byte, c *network.Conn) error {
	var m struct {
		Name string
		Pass string
	}
	network.Unmarshal(data, &m)

	glog.Infof("map %s user %s pass %s", c, m.Name, m.Pass)
	w := bufio.NewWriter(c)
	// SSMG_LOGIN_ACK
	network.Packet(w, 0x0c, 0x0011, uint32(0), uint16(0x0100), uint32(0x486EB420))

	return w.Flush()
}

const (
	LogoutOK     byte = 0
	LogoutCancle byte = 0xF9
	LogoutFailed byte = 0xFF
)

// Logout 登出
func Logout(data []byte, c *network.Conn) error {
	glog.Infof("map %s logout %d", c, data[0])
	// SSMG_LOGOUT
	w := bufio.NewWriter(c)
	network.Packet(w, 0x03, 0x0020, LogoutCancle)
	return w.Flush()
}

// CharSlot 角色登陆
func CharSlot(data []byte, c *network.Conn) error {
	var m struct {
		ID   uint16
		Slot uint8
	}
	network.Unmarshal(data, &m)

	glog.Infof("map %s login %v", c, m)
	w := bufio.NewWriter(c)
	// SSMG_ACTOR_SPEED
	network.Packet(w, 0x08, 0x1239,
		uint32(1),  // ActorID
		uint16(10), // Speed
	)
	// SSMG_ACTOR_MODE
	network.Packet(w, 0x0e, 0x0FA7,
		uint32(1), // ActorID
		uint32(2), // Mode1
		uint32(0), // Mode2
	)
	// SSMG_ACTOR_OPTION
	network.Packet(w, 0x06, 0x1A5F,
		uint32(1), // ActorID
		uint32(2), // Option
	)
	// SSMG_PLAYER_INFO
	network.Packet(w, 0xDE, 0x01FF,
		uint32(2), // ActorID
		uint32(6), // CharID
		&model.CharInfo{
			Name:      "鲷鱼烧",
			Race:      model.Titania,
			Sex:       model.Female,
			HairStyle: 7,
			HairColor: 60,
			Wig:       0xff,
			Exist:     0xff,
			Face:      14,
			MapID:     10065000,
			X:         50,
			Y:         50,
		},
	)
	return w.Flush()
}
