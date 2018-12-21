package main

import (
	"flag"

	"github.com/dragonflylee/go-saga/app"
	"github.com/dragonflylee/go-saga/model"
	"github.com/golang/glog"
)

func main() {
	db := flag.String("db", "saga.db", "path to database")
	login := flag.String("login", ":9090", "login listen address")
	flag.Parse()

	if err := model.Open(*db); err != nil {
		glog.Fatalf("open db failed: %v", err)
	}
	go app.StartLogin(*login)
	glog.Info("start ok")
	select {}
}
