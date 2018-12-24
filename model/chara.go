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

// CharData 角色数据
type CharData struct {
	Name      [4]string
	Race      [4]CharRace
	Unkown1   [4]byte
	Sex       [4]CharSex
	HairStyle [4]uint16
	HairColor [4]byte
	Wig       [4]uint16
	Exist     [4]byte
	Face      [4]byte
	Zero      byte    `json:"-"`
	Convert   [4]byte `json:"-"`
	Unkown2   [4]byte `json:"-"`
	Unkown3   [4]byte `json:"-"`
	Job       [4]byte
	Map       [4]uint32
	Level     [4]byte
	Job1      [4]byte
	Quest     [4]uint16
	Job2X     [4]byte
	Job2T     [4]byte
	Job3      [4]byte
}

// CharEquip 角色装备
type CharEquip [14]uint32

// CharInfo 角色数据
type CharInfo struct {
	Name      string
	Race      CharRace
	Unkown1   byte `json:"-"`
	Sex       CharSex
	HairStyle uint16
	HairColor byte
	Wig       uint16
	Exist     byte
	Face      byte
	Unkown2   uint16 `json:"-"`
	MapID     uint32
	X         byte
	Y         byte
	Dir       byte // 方向
	HP        uint32
	MaxHP     uint32
	MP        uint32
	MaxMP     uint32
	SP        uint32
	MaxSP     uint32
	EP        uint32
	MaxEP     uint32
	Status    [8]uint16 // 角色状态

}
