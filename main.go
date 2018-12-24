package main

import (
	"flag"
	"net"

	"github.com/dragonflylee/go-saga/app"
	"github.com/dragonflylee/go-saga/model"
	"github.com/golang/glog"
)

func main() {
	db := flag.String("db", "saga.db", "path to database")
	loginPort := flag.Int("login", 12000, "login listen port")
	mapPort := flag.Int("map", 12001, "map listen port")
	flag.Parse()

	if err := model.Open(*db); err != nil {
		glog.Fatalf("open db failed: %v", err)
	}
	go app.StartLogin(&net.TCPAddr{Port: *loginPort})
	go app.StartMap(&net.TCPAddr{Port: *mapPort})
	glog.Info("start ok")
	select {}
}
