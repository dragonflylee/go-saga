package network

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"net"

	"github.com/golang/glog"
)

// NewClient 处理连接
func NewClient(conn net.Conn, h Handler) *Conn {
	c := &Conn{
		auth:   clientAuth(0),
		conn:   conn,
		sendCh: make(chan []byte, 512),
	}
	go c.readPump(h)
	glog.Infof("new client: %s", conn.RemoteAddr())
	return c
}

type clientAuth int

func (c clientAuth) Init(r io.Reader) error {
	_, err := r.Read(make([]byte, 4))
	return err
}

func (c clientAuth) Auth(s *bufio.Scanner, w *bufio.Writer) ([]byte, error) {
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
		glog.Infof("recv m[%d] bit %d", i, m[i].BitLen())
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
