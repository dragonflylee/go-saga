package login

import (
	"bufio"
	"encoding/hex"

	"github.com/dragonflylee/go-saga/model"
	"github.com/dragonflylee/go-saga/network"
	"github.com/golang/glog"
)

// Version 版本号请求
func Version(data []byte, c *network.Conn) error {
	glog.Infof("recv %s version %s", c, hex.EncodeToString(data))

	w := bufio.NewWriter(c)
	// SSMG_VERSION_ACK
	network.Packet(w, 0x0a, 0x0002, data)
	// SSMG_LOGIN_ALLOWED
	network.Packet(w, 0x0a, 0x001e, data)
	return w.Flush()
}

// Ping 心跳
func Ping(data []byte, c *network.Conn) error {
	// SSMG_PONG
	w := bufio.NewWriter(c)
	network.Packet(w, 0x06, 0x000b)
	return w.Flush()
}

// CharData 角色数据
type CharData struct {
	Name      [4]string
	Race      [4]model.CharRace
	Unkown1   [4]byte
	Sex       [4]model.CharSex
	HairStyle [4]uint16
	HairColor [4]byte
	Wig       [4]uint16
	Exist     [4]byte
	Face      [4]byte
	Zero      byte    `json:"-"`
	Convert   [4]byte `json:"-"`
	Unkown2   [4]byte `json:"-"`
	Unkown3   [4]byte `json:"-"`
	Job       [4]byte
	Map       [4]uint32
	Level     [4]byte
	Job1      [4]byte
	Quest     [4]uint16
	Job2X     [4]byte
	Job2T     [4]byte
	Job3      [4]byte
}

// CharEquip 角色装备
type CharEquip [14]uint32

// Login 登录
func Login(data []byte, c *network.Conn) error {
	var m struct {
		Name string
		Pass string
	}
	network.Unmarshal(data, &m)

	glog.Infof("recv %s login %s pass %s", c, m.Name, m.Pass)
	w := bufio.NewWriter(c)
	// SSMG_LOGIN_ACK
	network.Packet(w, 0x13, 0x0020)
	// SSMG_CHAR_DATA
	network.Packet(w, 0xc0, 0x0028, &CharData{
		Name:      [4]string{"鲷鱼烧"},
		Race:      [4]model.CharRace{model.Titania},
		Sex:       [4]model.CharSex{model.Female},
		HairStyle: [4]uint16{7},
		HairColor: [4]byte{60},
		Wig:       [4]uint16{0xff},
		Exist:     [4]byte{0xff},
		Face:      [4]byte{14},
		Job:       [4]byte{47},
		Convert:   [4]byte{43},
		Map:       [4]uint32{10065000},
		Level:     [4]byte{110},
		Job1:      [4]byte{50},
		Quest:     [4]uint16{5},
		Job2X:     [4]byte{50},
		Job2T:     [4]byte{50},
		Job3:      [4]byte{50},
	})
	// SSMG_CHAR_EQUIP
	var eq CharEquip
	eq[model.UpperBody] = 90000030
	//eq[model.LowerBody] = 50249503
	eq[model.Shoes] = 50248000
	// eq[model.Back] = 10020114
	network.Packet(w, 0xa1, 0x0029, eq)

	return w.Flush()
}

// RequestMap 请求地图
func RequestMap(data []byte, c *network.Conn) error {
	// SSMG_SEND_TO_MAP_SERVER
	w := bufio.NewWriter(c)
	network.Packet(w, 0xb7, 0x0033, byte(0x10), "127.0.0.1", uint32(9090))
	return w.Flush()
}
