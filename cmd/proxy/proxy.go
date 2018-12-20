package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/dragonflylee/go-saga/server"
)

var (
	m     = make(map[string]string)
	start = int32(12000)
)

func main() {
	remote := flag.String("remote", "127.0.0.1:9090", "remote address")
	flag.Parse()

	listenPort(*remote)
	select {}
}

func listenPort(remote string) string {
	if port, exist := m[remote]; exist {
		return port
	}
	addr := fmt.Sprintf("127.0.0.1:%d", atomic.AddInt32(&start, 1))
	m[remote] = addr

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen fail: %v", err)
	}
	go func() {
		defer l.Close()
		log.Printf("listen %s start", l.Addr())

		for {
			c, err := l.Accept()
			if err != nil {
				log.Printf("accept failed: %v", err)
				continue
			}
			// create goroutine for each connect
			go handleConn(c, remote)
		}
	}()
	return addr
}

func handleConn(c net.Conn, remote string) {
	var from, to io.Writer
	r, err := net.Dial("tcp", remote)
	if err != nil {
		log.Printf("dial %s fail: %v", c.RemoteAddr(), err)
		c.Close()
		return
	}
	// 握手包
	if err = binary.Write(r, binary.BigEndian, uint64(0x10)); err != nil {
		log.Printf("dial %s fail: %v", c.RemoteAddr(), err)
		r.Close()
		c.Close()
		return
	}
	to = server.NewClient(r, server.HandleFunc(func(data []byte, w io.Writer) error {
		log.Printf("recv %s: %s", r.RemoteAddr(), hex.EncodeToString(data))
		for from == nil {
			time.Sleep(time.Second)
		}
		_, err := from.Write(data)
		if err != nil {
			r.Close()
		}
		return err
	}))
	from = server.HandleClient(c, server.HandleFunc(func(data []byte, w io.Writer) error {
		log.Printf("send %s: %s", r.RemoteAddr(), hex.EncodeToString(data))
		for to == nil {
			time.Sleep(time.Second)
		}
		if binary.BigEndian.Uint16(data[2:4]) == 0x33 {
			buf := bufio.NewScanner(bytes.NewReader(data[4:]))
			buf.Split(func(data []byte, atEOF bool) (int, []byte, error) {
				n := int(data[0])
				if atEOF || n <= 0 {
					return len(data) + 1, nil, nil
				}
				return int(n + 1), data[1:n], nil
			})
			for buf.Scan() {
				log.Println(buf.Text())
			}
		}
		_, err := to.Write(data)
		if err != nil {
			c.Close()
		}
		return err
	}))

}
