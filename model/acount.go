package model

// User 账户
type User struct {
	ID              int64  `gorm:"primary_key;auto_increment"`
	Username        string `gorm:"size:30;not null"`
	Password        string `gorm:"size:32;not null"`
	DelPass         string `gorm:"size:32;not null"`
	Banned          uint8  `gorm:"size:3;not null;default 0"`
	GMLevel         uint8  `gorm:"size:3;not null;default 0"`
	Bank            int64  `gorm:"size:10;not null;default 0"`
	VshopPoints     int64  `gorm:"size:10;not null;default 0"`
	UsedVshopPoints int64  `gorm:"size:10;not null;default 0"`
	LastIP          string `gorm:"size:20"`
}

// Select 查询
func (u *User) Select() error {
	return db.New().Take(u).Error
}

// Update 更新
func (u *User) Update() error {
	return db.New().Model(u).Updates(u).Error
}
