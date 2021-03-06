package app

import (
	"net"
	"runtime/debug"

	"github.com/dragonflylee/go-saga/app/login"
	"github.com/dragonflylee/go-saga/network"
	"github.com/golang/glog"
)

// StartLogin 开始监听接入服务
func StartLogin(addr net.Addr) {
	defer func() {
		if err := recover(); err != nil {
			glog.Error(err)
			debug.PrintStack()
		}
	}()

	r := make(network.Route)
	// CSMG_SEND_VERSION
	r.Handle(0x0001, login.Version)
	// CSMG_PING
	r.Handle(0x000A, login.Ping)
	// CSMG_LOGIN
	r.Handle(0x001F, login.Login)
	// CSMG_CHAR_CREATE
	r.Handle(0x00A0, login.CharCreate)
	// CSMG_CHAR_DELETE
	r.Handle(0x00A5, login.CharDelete)
	// CSMG_CHAR_SELECT
	r.Handle(0x00A7, login.CharSelect)
	// CSMG_REQUEST_MAP_SERVER
	r.Handle(0x0032, login.RequestMap)

	s, err := net.Listen(addr.Network(), addr.String())
	if err != nil {
		glog.Warningf("listen fail: %v", err)
		return
	}
	glog.Infof("start login listen %s", addr)

	for {
		c, err := s.Accept()
		if err != nil {
			glog.Warningf("login accept failed: %v", err)
			continue
		}
		// create goroutine for each connect
		network.HandleClient(c, r)
	}
}
