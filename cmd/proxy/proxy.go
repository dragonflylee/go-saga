package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/dragonflylee/go-saga/network"
	"github.com/golang/glog"
)

var (
	m     = make(map[string]*net.TCPAddr)
	lock  sync.Mutex
	start = int32(12000)
)

func main() {
	remote := flag.String("remote", "127.0.0.1:9090", "remote address")
	flag.Parse()

	listenPort(*remote)
	select {}
}

func listenPort(remote string) *net.TCPAddr {
	lock.Lock()
	defer lock.Unlock()
	if port, exist := m[remote]; exist {
		return port
	}

	addr := &net.TCPAddr{
		IP:   net.IP{127, 0, 0, 1},
		Port: int(atomic.AddInt32(&start, 1)),
	}
	m[remote] = addr

	glog.Warningf("proxy %s => %s start", addr, remote)

	l, err := net.Listen("tcp", addr.String())
	if err != nil {
		log.Fatalf("listen fail: %v", err)
	}
	go func() {
		defer l.Close()

		for {
			c, err := l.Accept()
			if err != nil {
				glog.Warningf("accept failed: %v", err)
				continue
			}
			// create goroutine for each connect
			go handleConn(c, remote)
		}
	}()
	return addr
}

func handleConn(c net.Conn, remote string) {
	r, err := net.Dial("tcp", remote)
	if err != nil {
		glog.Warningf("dial %s fail: %v", c.RemoteAddr(), err)
		c.Close()
		return
	}

	var to, from io.Writer
	from = network.HandleClient(c, network.HandleFunc(func(data []byte, w *network.Conn) error {
		glog.Infof("send %s: %s", r.RemoteAddr(), hex.EncodeToString(data))
		_, err := to.Write(data)
		if err != nil {
			c.Close()
		}
		return err
	}))

	to = network.NewClient(r, network.HandleFunc(func(data []byte, w *network.Conn) error {
		glog.Infof("recv %s: %s", r.RemoteAddr(), hex.EncodeToString(data))
		if len(data) > 4 && binary.BigEndian.Uint16(data[2:4]) == 0x33 {
			if data[5] != 0x10 {
				var m struct {
					Server string
					Addr   string
				}
				network.Unmarshal(data[4:], &m)
				n := strings.Split(m.Addr[1:], ",")
				for i := 0; i < len(n); i++ {
					n[i] = listenPort(strings.TrimSpace(n[i])).String()
				}
				m.Addr = strings.Join(n, ",")

				b := &bytes.Buffer{}
				network.Packet(b, uint16(len(data)-2), 0x33, m)
				data = b.Bytes()
			} else {
				var m struct {
					ID   byte
					Addr string
					Port uint32
				}
				network.Unmarshal(data[4:], &m)
				n := listenPort(fmt.Sprintf("%s:%d", m.Addr, m.Port))
				m.Addr = n.IP.String()
				m.Port = uint32(n.Port)

				b := &bytes.Buffer{}
				network.Packet(b, uint16(len(data)-2), 0x33, m)
				data = b.Bytes()
			}

		}
		_, err := from.Write(data)
		if err != nil {
			r.Close()
		}
		return err
	}))

	// 握手包
	if err = binary.Write(r, binary.BigEndian, uint64(0x10)); err != nil {
		glog.Warningf("dial %s fail: %v", c.RemoteAddr(), err)
		r.Close()
		c.Close()
		return
	}

}
