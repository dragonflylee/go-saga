package server

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
)

var two = big.NewInt(2)

// HandleClient 处理连接
func HandleClient(conn net.Conn, h Handler) io.Writer {
	c := &wrapper{
		auther: serve(0),
		conn:   conn,
		sendCh: make(chan []byte, 512),
	}
	go c.readPump(h)
	log.Printf("handle client: %s", conn.RemoteAddr())
	return c
}

type serve int

func (c serve) Init(r io.Reader) error {
	var sig uint64
	// 读取握手消息
	err := binary.Read(r, binary.BigEndian, &sig)
	if err != nil {
		return err
	}
	if sig != 0x10 {
		return fmt.Errorf("bad sig 0x%.8x", sig)
	}
	return nil
}

func (c serve) Auth(s *bufio.Scanner, w *bufio.Writer) ([]byte, error) {
	var (
		m, p *big.Int
		err  error
	)
	// 生成私钥
	if m, err = rand.Prime(rand.Reader, 0x400); err != nil {
		return nil, err
	}
	if p, err = rand.Prime(rand.Reader, 0x100); err != nil {
		return nil, err
	}
	// 发送Token
	if err = binary.Write(w, binary.BigEndian, uint32(0)); err != nil {
		return nil, err
	}
	for _, s := range []string{two.Text(16), m.Text(16), new(big.Int).Exp(two, p, m).Text(16)} {
		if err = binary.Write(w, binary.BigEndian, uint32(len(s))); err != nil {
			return nil, err
		}
		if _, err = w.WriteString(s); err != nil {
			return nil, err
		}
	}
	if err = w.Flush(); err != nil {
		return nil, err
	}
	// 读取密钥
	if !s.Scan() {
		return nil, s.Err()
	}
	x, ok := new(big.Int).SetString(s.Text(), 16)
	if !ok {
		return nil, fmt.Errorf("bad sig %s", hex.EncodeToString(s.Bytes()))
	}
	key := new(big.Int).Exp(x, p, m).Bytes()
	for i := 0; i < 16; i++ {
		high, low := key[i]>>4, key[i]&0xF
		if high > 9 {
			high -= 9
		}
		if low > 9 {
			low -= 9
		}
		key[i] = high<<4 | low
	}
	return key, nil
}
