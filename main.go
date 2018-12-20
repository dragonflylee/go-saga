package main

import (
	"flag"
	"io"
	"log"
	"net"
	"runtime/debug"

	"github.com/dragonflylee/go-saga/protocol"
	"github.com/dragonflylee/go-saga/server"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9090", "server listen address")
	mode := flag.String("mode", "login", "run mode")
	flag.Parse()

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			debug.PrintStack()
		}
	}()

	r := make(server.Route)
	r.Handle(0x0000, protocol.HandleTODO)
	// CSMG_SEND_VERSION
	r.Handle(0x0001, protocol.HandleVersion)
	// CSMG_PING
	r.Handle(0x000a, protocol.HandlePing)
	// CSMG_LOGIN
	r.Handle(0x001f, protocol.HandleLogin)
	r.Handle(0x002f, func(data []byte, c io.Writer) error {
		_, err := c.Write(protocol.NewPacket(0x6, 0x0030))
		return err
	})
	// CSMG_REQUEST_MAP_SERVER
	r.Handle(0x0031, func(data []byte, c io.Writer) error {
		c.Write(protocol.NewPacket(0x2, 0x0032))
		// SSMG_SEND_TO_MAP_SERVER
		w := protocol.NewPacket(0xb7, 0x0033)
		w[4] = byte(copy(w[5:], []byte("ECO")) + 1)
		w[5+w[4]] = byte(copy(w[6+w[4]:], []byte(*addr)) + 1)
		_, err := c.Write(w)
		c.Write(protocol.NewPacket(0x2, 0x0034))
		return err
	})

	s, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Printf("listen fail: %v", err)
		return
	}
	log.Printf("mode %s listen %s", *mode, *addr)

	for {
		c, err := s.Accept()
		if err != nil {
			log.Printf("accept failed: %v", err)
			continue
		}
		// create goroutine for each connect
		server.HandleClient(c, r)
	}
}
