package model

import (
	"flag"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	// 数据库驱动
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db    *gorm.DB
	debug = flag.Bool("d", false, "debug mode")
)

// Open 连接数据库
func Open(path string) (err error) {
	if db, err = gorm.Open("sqlite3", path); err != nil {
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
	return db.AutoMigrate(&User{}).Error
}
