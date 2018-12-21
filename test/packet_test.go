package test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/dragonflylee/go-saga/network"
)

func Test_Marshal(t *testing.T) {
	var m struct {
		Name string
		Pass string
	}
	m.Name = "user"
	m.Pass = "pass"
	w := &bytes.Buffer{}
	if _, err := network.Marshal(w, m); err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(w.Bytes()))
}

func Test_Packet(t *testing.T) {
	var m struct {
		Server string
		Addr   string
	}
	m.Server = "ECO-RE"
	m.Addr = "127.0.0.1:12002"
	w := &bytes.Buffer{}
	network.Packet(w, 0xb7, 0x33, &m)
}
