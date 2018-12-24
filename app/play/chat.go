package play

import (
	"github.com/dragonflylee/go-saga/network"
	"github.com/golang/glog"
)

// ChatMotion 聊天行为
func ChatMotion(data []byte, c *network.Conn) error {
	var m struct {
		Motion uint16
		Loop   uint8
	}
	network.Unmarshal(data, &m)

	glog.Infof("map %s chat motion %v", c, m)
	return nil
}
