package network

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/golang/glog"
)

// 缓存大小
const maxBufSize = 0x1000

// Conn 连接对象
type Conn struct {
	sync.Map
	auth interface {
		Init(r io.Reader) error
		Auth(s *bufio.Scanner, w *bufio.Writer) ([]byte, error)
	}
	conn   net.Conn
	block  cipher.Block
	sendCh chan []byte
}

// String 实现接口 Stringer
func (c *Conn) String() string {
	return c.conn.RemoteAddr().String()
}

// Close 关闭连接
func (c *Conn) Close() error {
	return c.conn.Close()
}

// Write 实现 io.Writer
func (c *Conn) Write(s []byte) (int, error) {
	select {
	case c.sendCh <- s:
		return len(s), nil
	default:
		return 0, fmt.Errorf("send %s full", c)
	}
}

// readPump 读线程
func (c *Conn) readPump(h Handler) {
	defer c.conn.Close()

	if err := c.auth.Init(c.conn); err != nil {
		glog.Warningf("init %s failed: %v", c, err)
		return
	}
	// 读取数据
	s := bufio.NewScanner(c.conn)
	s.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		n := binary.BigEndian.Uint32(data[:4])
		if n > 0xFFFF {
			return len(data), nil, bufio.ErrFinalToken
		}
		if c.block == nil {
			if len(data) < int(n+4) {
				// 还没读完
				return 0, nil, nil
			}
			return int(n + 4), data[4 : n+4], nil
		}
		if int(n)%c.block.BlockSize() > 0 {
			return 0, nil, fmt.Errorf("block %d no padding", n)
		}
		if len(data) < int(n+8) {
			// 还没读完
			return 0, nil, nil
		}
		// AES 解码
		for i := 8; i < int(n+8); {
			end := i + c.block.BlockSize()
			c.block.Decrypt(data[i:end], data[i:end])
			i = end
		}
		off := binary.BigEndian.Uint32(data[4:8])
		if off+4 > n {
			return 0, nil, fmt.Errorf("offet %d too short %d", off, n)
		}
		return int(n + 8), data[8 : off+8], nil
	})

	// 初始化 AES 密钥
	key, err := c.auth.Auth(s, bufio.NewWriter(c.conn))
	if err != nil {
		glog.Warningf("auth %s failed: %v", c, err)
		return
	}
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
	if c.block, err = aes.NewCipher(key[:16]); err != nil {
		glog.Warningf("init %s failed: %v", c, err)
		return
	}
	glog.Infof("init %s aes: %s", c, hex.EncodeToString(key[:16]))
	// 写线程
	go c.writePump()
	// 初始包
	// if err = h.Serve([]byte{0, 2, 0, 0}, c); err != nil {
	//	glog.Infof("serv %s failed: %v", c, err)
	// }
	// 处理数据包
	for s.Scan() {
		if err = h.Serve(s.Bytes(), c); err != nil {
			glog.Warningf("serv %s failed: %v", c, err)
		}
	}
	glog.Infof("recv %s closed: %v", c, s.Err())
}

// writePump 写线程
func (c *Conn) writePump() {
	defer c.conn.Close()

	buf := make([]byte, maxBufSize)
	for {
		s, ok := <-c.sendCh
		if !ok || c.block == nil {
			return
		}
		// Padding
		padding := c.block.BlockSize() - len(s)%c.block.BlockSize()
		size := len(s) + padding
		binary.BigEndian.PutUint32(buf[:4], uint32(size))
		binary.BigEndian.PutUint32(buf[4:8], uint32(len(s)))
		if padding > 0 {
			s = append(s, bytes.Repeat([]byte{byte(padding)}, padding)...)
		}
		// AES 加密
		for i := 0; i < size; {
			end := i + c.block.BlockSize()
			c.block.Encrypt(buf[i+8:end+8], s[i:end])
			i = end
		}
		if _, err := c.conn.Write(buf[:size+8]); err != nil {
			glog.Warningf("write %s failed: %v", c, err)
			break
		}
	}
}
