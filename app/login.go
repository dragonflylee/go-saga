package app

import (
	"net"
	"runtime/debug"

	"github.com/dragonflylee/go-saga/app/login"
	"github.com/dragonflylee/go-saga/network"
	"github.com/golang/glog"
)

// StartLogin 开始监听登录服务
func StartLogin(addr string) {
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
	r.Handle(0x000a, login.Ping)
	// CSMG_LOGIN
	r.Handle(0x001f, login.Login)
	// CSMG_CHAR_CREATE
	r.Handle(0x00a0, login.CharCreate)
	// CSMG_CHAR_DELETE
	r.Handle(0x00a5, login.CharDelete)
	// CSMG_CHAR_SELECT
	r.Handle(0x00a7, login.CharSelect)
	// CSMG_REQUEST_MAP_SERVER
	r.Handle(0x0032, login.RequestMap)

	s, err := net.Listen("tcp", addr)
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
