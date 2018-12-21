package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/dragonflylee/go-saga/network"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9090", "server address")
	flag.Parse()

	r := make(network.Route)
	r.Handle(0x0000, func(data []byte, c *network.Conn) error {
		w := bufio.NewWriter(c)
		network.Packet(w, 0x0a, 0x0001, uint64(0x3e8015f3771))
		return w.Flush()
	})
	r.Handle(0x0002, func(data []byte, c *network.Conn) error {
		return nil
	})
	r.Handle(0x0020, func(data []byte, c *network.Conn) error {
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
	r.Handle(0x001e, func(data []byte, c *network.Conn) error {
		var m struct {
			Front, Back uint32
		}
		network.Unmarshal(data, &m)
		log.Printf("allow 0x%08x, 0x%08x", m.Front, m.Back)

		hash := md5.Sum([]byte("lilongfei"))
		token := sha1.Sum([]byte(fmt.Sprintf("%d%s%d", binary.BigEndian.Uint32(data[:4]),
			hex.EncodeToString(hash[:]), binary.BigEndian.Uint32(data[4:]))))

		w := bufio.NewWriter(c)
		network.Packet(w, 0x50, 0x001f, "dragonfly2", hex.EncodeToString(token[:]))
		log.Printf("send login: %d", w.Buffered())

		network.Packet(w, 0x02, 0x002f)
		return w.Flush()
	})
	r.Handle(0x000b, func(data []byte, c *network.Conn) error {
		return nil
	})
	r.Handle(0x0030, func(data []byte, c *network.Conn) error {
		w := bufio.NewWriter(c)
		network.Packet(w, 0x02, 0x0031)
		return w.Flush()
	})
	r.Handle(0x0033, func(data []byte, c *network.Conn) error {
		var m struct {
			Server string
			Addr   string
		}
		network.Unmarshal(data, &m)
		n := strings.LastIndex(m.Addr, ",")
		if n > 0 && len(m.Addr) > n {
			m.Addr = m.Addr[n+1:]
		}
		log.Printf("recv %s map %s", c, m.Addr)
		return nil
	})

	startClient(*addr, r)
	select {}
}

func startClient(address string, r network.Route) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dial %s ok", conn.RemoteAddr())

	// 握手包
	if err = binary.Write(conn, binary.BigEndian, uint64(0x10)); err != nil {
		log.Fatal(err)
	}
	network.NewClient(conn, r)
}
