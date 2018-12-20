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

// NewClient 处理连接
func NewClient(conn net.Conn, h Handler) io.Writer {
	c := &wrapper{
		auther: client(0),
		conn:   conn,
		sendCh: make(chan []byte, 512),
	}
	go c.readPump(h)
	log.Printf("handle client: %s", conn.RemoteAddr())
	return c
}

type client int

func (c client) Init(r io.Reader) error {
	_, err := r.Read(make([]byte, 4))
	return err
}

func (c client) Auth(s *bufio.Scanner, w *bufio.Writer) ([]byte, error) {
	// 读取私钥
	var (
		m   [3]*big.Int
		p   *big.Int
		ok  bool
		err error
	)
	for i := 0; i < len(m); i++ {
		if !s.Scan() {
			return nil, s.Err()
		}
		if m[i], ok = new(big.Int).SetString(s.Text(), 16); !ok {
			return nil, fmt.Errorf("bad token %s", hex.EncodeToString(s.Bytes()))
		}
		log.Printf("recv m[%d] bit %d", i, m[i].BitLen())
	}
	// 生成交换密钥
	if p, err = rand.Prime(rand.Reader, 0x100); err != nil {
		return nil, err
	}
	// 发送Token
	sig := new(big.Int).Exp(m[0], p, m[1]).Text(16)
	if err = binary.Write(w, binary.BigEndian, uint32(len(sig))); err != nil {
		return nil, err
	}
	if _, err = w.WriteString(sig); err != nil {
		return nil, err
	}
	if err = w.Flush(); err != nil {
		return nil, err
	}
	return new(big.Int).Exp(m[2], p, m[1]).Bytes(), nil
}
