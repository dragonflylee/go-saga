package app

import (
	"net"
	"runtime/debug"

	"github.com/dragonflylee/go-saga/app/play"
	"github.com/dragonflylee/go-saga/network"
	"github.com/golang/glog"
)

// StartMap 开始监听地图服务
func StartMap(addr net.Addr) {
	defer func() {
		if err := recover(); err != nil {
			glog.Error(err)
			debug.PrintStack()
		}
	}()

	r := make(network.Route)
	// CSMG_SEND_VERSION
	r.Handle(0x000A, play.Version)
	// CSMG_PING
	r.Handle(0x0032, play.Ping)
	// CSMG_LOGIN
	r.Handle(0x0010, play.Login)
	// CSMG_LOGOUT
	r.Handle(0x001F, play.Logout)
	// CSMG_CHAR_SLOT
	r.Handle(0x01FD, play.CharSlot)
	// CSMG_CHAT_MOTION
	r.Handle(0x121B, play.ChatMotion)

	s, err := net.Listen(addr.Network(), addr.String())
	if err != nil {
		glog.Warningf("listen fail: %v", err)
		return
	}
	glog.Infof("start map listen %s", addr)

	for {
		c, err := s.Accept()
		if err != nil {
			glog.Warningf("map accept failed: %v", err)
			continue
		}
		// create goroutine for each connect
		network.HandleClient(c, r)
	}
}
