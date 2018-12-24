package login

import (
	"bufio"
	"encoding/json"
	"io"

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
	network.Packet(w, 0x14, 0x00A8, uint32(10024000))
	return w.Flush()
}

// CharDelete 角色删除
func CharDelete(data []byte, c *network.Conn) error {
	// SSMG_CHAR_DELETE_ACK
	w := bufio.NewWriter(c)
	network.Packet(w, 0x03, 0x00A6, byte(0))
	SendCharData(w)
	return w.Flush()
}

// SendCharData 发送角色数据
func SendCharData(w io.Writer) {
	// SSMG_CHAR_DATA
	network.Packet(w, 0xC0, 0x0028, &model.CharData{
		Name:      [4]string{"鲷鱼烧"},
		Race:      [4]model.CharRace{model.Titania},
		Sex:       [4]model.CharSex{model.Female},
		HairStyle: [4]uint16{7},
		HairColor: [4]byte{60},
		Wig:       [4]uint16{0xff},
		Exist:     [4]byte{0xff},
		Face:      [4]byte{14},
		Job:       [4]byte{67},
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
	var eq model.CharEquip
	eq[model.SlotUpperBody] = 60102950
	//eq[model.SlotLowerBody] = 50010456
	eq[model.SlotShoes] = 50066750
	eq[model.SlotBack] = 10057650
	//eq[model.SlotSocks] = 50011101
	eq[model.SlotRightHand] = 61040250
	network.Packet(w, 0xA1, 0x0029, eq)
}
