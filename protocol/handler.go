package protocol

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
)

// NewPacket 新封包
func NewPacket(len, id uint16) []byte {
	w := make([]byte, len+2)
	binary.BigEndian.PutUint16(w[:2], len)
	binary.BigEndian.PutUint16(w[2:4], id)
	return w
}

// HandleTODO 空任务
func HandleTODO(data []byte, c io.Writer) error {
	return nil
}

// HandleVersion 版本号请求
func HandleVersion(data []byte, c io.Writer) error {
	log.Printf("recv %s version %s", c, hex.EncodeToString(data))

	// SSMG_Universal
	w := NewPacket(0x16, 0xffff)
	hex.Decode(w[4:], []byte("e86a6acadce806052b29f8962f867cab2a57ad30"))
	if _, err := c.Write(w); err != nil {
		return err
	}
	// SSMG_VERSION_ACK
	w = NewPacket(0xa, 0x0002)
	copy(w[4:], data)
	if _, err := c.Write(w); err != nil {
		return err
	}
	// SSMG_LOGIN_ALLOWED
	w = NewPacket(0xa, 0x001e)
	rand.Read(w[4:])
	if _, err := c.Write(w); err != nil {
		return err
	}
	return nil
}

// HandlePing 心跳
func HandlePing(data []byte, c io.Writer) error {
	// SSMG_PONG
	_, err := c.Write(NewPacket(0x06, 0x000b))
	return err
}

// HandleLogin 登录
func HandleLogin(data []byte, c io.Writer) error {
	s := bufio.NewScanner(bytes.NewReader(data))
	s.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		n := int(data[0])
		if atEOF || n <= 0 {
			return len(data) + 1, nil, nil
		}
		return int(n + 1), data[1:n], nil
	})
	for s.Scan() {
		log.Printf("recv %s login %s", c, s.Text())
	}
	// SSMG_LOGIN_ACK
	_, err := c.Write(NewPacket(0x13, 0x0020))
	/*// SSMG_CHAR_DATA
	w = NewPacket(86, 0x0028)
	if _, err := c.Write(w); err != nil {
		return err
	}
	// SSMG_CHAR_EQUIP
	w = NewPacket(161, 0x0029)
	if _, err := c.Write(w); err != nil {
		return err
	}*/
	_, err = c.Write(NewPacket(0x4, 0x0150))
	return err
}
