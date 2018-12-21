package login

import (
	"bufio"
	"encoding/json"

	"github.com/dragonflylee/go-saga/model"
	"github.com/dragonflylee/go-saga/network"
	"github.com/golang/glog"
)

// CharCreate 角色创建
func CharCreate(data []byte, c *network.Conn) error {
	var m struct {
		Slot      byte
		Name      string
		Race      model.CharRace
		Sex       model.CharSex
		HairStyle uint16
		HairColor byte
		Face      byte
	}
	network.Unmarshal(data, &m)
	b, _ := json.Marshal(m)
	glog.Infof("char create %v", string(b))

	// SSMG_CHAR_CREATE_ACK
	w := bufio.NewWriter(c)
	network.Packet(w, 0x06, 0x00a1)
	return w.Flush()
}

// CharSelect 角色选择
func CharSelect(data []byte, c *network.Conn) error {
	// SSMG_CHAR_SELECT_ACK
	w := bufio.NewWriter(c)
	network.Packet(w, 0x14, 0x00a8, uint32(10024000))
	return w.Flush()
}

// CharDelete 角色删除
func CharDelete(data []byte, c *network.Conn) error {
	// SSMG_CHAR_DELETE_ACK
	w := bufio.NewWriter(c)
	network.Packet(w, 0x03, 0x00a6, byte(0))
	return w.Flush()
}
