package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/dragonflylee/go-saga/server"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:12001", "server address")
	flag.Parse()

	r := make(server.Route)
	r.Handle(0x0000, func(data []byte, c io.Writer) error {
		w := newPacket(10, 0x0001)
		hex.Decode(w[4:], []byte("000003e8015f3771"))
		_, err := c.Write(w)
		return err
	})
	r.Handle(0x0002, func(data []byte, c io.Writer) error {
		return nil
	})
	r.Handle(0x0020, func(data []byte, c io.Writer) error {
		result := binary.BigEndian.Uint32(data[:4])
		switch int32(result) {
		case 0:
			log.Printf("login ok")
		case -2:
			log.Printf("unknown account")
		case -3:
			log.Printf("password incorrect")
		case -4:
			log.Printf("account locked")
		case -5:
			log.Printf("account reset")
		case -6:
			log.Printf("maintenance")
		default:
			log.Printf("login error: 0x%08x", result)
		}
		return nil
	})
	r.Handle(0x001e, func(data []byte, c io.Writer) error {
		w := newPacket(55, 0x001f)

		w[4] = byte(copy(w[5:], []byte("dragonfly2")) + 1)
		log.Printf("allow %d: %s", len(data), hex.EncodeToString(data))

		hash := md5.Sum([]byte("lilongfei"))
		token := sha1.Sum([]byte(fmt.Sprintf("%d%s%d", binary.BigEndian.Uint32(data[:4]),
			hex.EncodeToString(hash[:]), binary.BigEndian.Uint32(data[4:]))))
		pass := []byte(hex.EncodeToString(token[:]))
		w[5+w[4]] = byte(copy(w[6+w[4]:], pass) + 1)

		_, err := c.Write(w)
		log.Printf("send login: %d", len(w))

		_, err = c.Write(newPacket(2, 0x002f))
		return err
	})
	r.Handle(0x000b, func(data []byte, c io.Writer) error {
		return nil
	})
	r.Handle(0x0030, func(data []byte, c io.Writer) error {
		_, err := c.Write(newPacket(2, 0x0031))
		return err
	})
	r.Handle(0x0033, func(data []byte, c io.Writer) error {
		s := bufio.NewScanner(bytes.NewReader(data))
		s.Split(func(data []byte, atEOF bool) (int, []byte, error) {
			n := int(data[0])
			if atEOF || n <= 0 {
				return len(data) + 1, nil, nil
			}
			return int(n + 1), data[1:n], nil
		})
		for s.Scan() {
			log.Printf("recv %s map %s", c, s.Text())
		}
		return nil
	})

	startClient(*addr, r)
	select {}
}

func startClient(address string, r server.Route) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dial %s ok", conn.RemoteAddr())

	// 握手包
	if err = binary.Write(conn, binary.BigEndian, uint64(0x10)); err != nil {
		log.Fatal(err)
	}
	server.NewClient(conn, r)
}

// newPacket 新封包
func newPacket(len, id uint16) []byte {
	w := make([]byte, len+2)
	binary.BigEndian.PutUint16(w[:2], len)
	binary.BigEndian.PutUint16(w[2:4], id)
	return w
}
