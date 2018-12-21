package model

// CharRace 种族
type CharRace byte

const (
	// Emil 埃米尔
	Emil CharRace = iota
	// Titania 塔尼亚
	Titania
	// Dominion 道米尼
	Dominion
	// DEM
	DEM
)

// CharSex 性别
type CharSex byte

const (
	// Male 男性
	Male CharSex = iota
	// Female 女性
	Female
)

const (
	Cabalist byte = 73
)
