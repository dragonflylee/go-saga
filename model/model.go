package model

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	// 数据库驱动
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db    *gorm.DB
	debug = flag.Bool("d", false, "debug mode")
)

// DBConfig 数据库配置项
type DBConfig struct {
	Type string `json:"type"`
	Host string `json:"host,omitempty"`
	Port uint64 `json:"port,omitempty"`
	User string `json:"user,omitempty"`
	Pass string `json:"pass,omitempty"`
	Name string `json:"name"`
}

// Open 连接数据库
func Open(conf *DBConfig) error {
	var (
		source string
		err    error
	)
	if conf.Type == "mysql" {
		source = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&allowOldPasswords=1",
			conf.User, conf.Pass, conf.Host, conf.Port, conf.Name)
	} else if conf.Type == "postgres" {
		source = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
			conf.User, conf.Pass, conf.Host, conf.Port, conf.Name)
	} else {
		return errors.New("数据库类型不支持")
	}
	if db, err = gorm.Open(conf.Type, source); err != nil {
		return fmt.Errorf("connect database failed: %v", err)
	}
	db.BlockGlobalUpdate(true)
	if debug != nil {
		db.LogMode(*debug)
	}
	if time.Local, err = time.LoadLocation("Asia/Chongqing"); err != nil {
		return fmt.Errorf("load location failed: %v", err)
	}
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}
	return db.AutoMigrate().Error
}
