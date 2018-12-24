package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"net"

	"github.com/dragonflylee/go-saga/network"
	"github.com/golang/glog"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9091", "server address")
	flag.Parse()

	r := make(network.Route)
	r.Handle(0x0002, func(data []byte, c *network.Conn) error {
		return nil
	})
	r.Handle(0x0020, func(data []byte, c *network.Conn) error {
		result := binary.BigEndian.Uint32(data[:4])
		switch int32(result) {
		case 0:
			glog.Infof("login ok")
		case -2:
			glog.Infof("unknown account")
		case -3:
			glog.Infof("password incorrect")
		case -4:
			glog.Infof("account locked")
		case -5:
			glog.Infof("account reset")
		case -6:
			glog.Infof("maintenance")
		default:
			glog.Infof("login error: 0x%08x", result)
		}
		return nil
	})
	r.Handle(0x001e, func(data []byte, c *network.Conn) error {
		var m struct {
			Front, Back uint32
		}
		network.Unmarshal(data, &m)
		glog.Infof("allow 0x%08x, 0x%08x", m.Front, m.Back)

		hash := md5.Sum([]byte("lilongfei"))
		token := sha1.Sum([]byte(fmt.Sprintf("%d%s%d", binary.BigEndian.Uint32(data[:4]),
			hex.EncodeToString(hash[:]), binary.BigEndian.Uint32(data[4:]))))

		w := bufio.NewWriter(c)
		network.Packet(w, 0x50, 0x001f, "dragonfly2", hex.EncodeToString(token[:]))
		glog.Infof("send login: %d", w.Buffered())

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
		glog.Infof("recv %s map %s", c, m.Addr)
		return nil
	})

	startClient(*addr, r)
}

func startClient(address string, r network.Route) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Infof("dial %s ok", conn.RemoteAddr())

	// 握手包
	if err = binary.Write(conn, binary.BigEndian, uint64(0x10)); err != nil {
		glog.Fatal(err)
	}
	c := network.NewClient(conn)
	// 初始包
	w := bufio.NewWriter(c)
	network.Packet(w, 0x0a, 0x0001, uint64(0x3e8015f3771))
	if err = w.Flush(); err != nil {
		glog.Fatal(err)
	}
	c.Run(r)
}
