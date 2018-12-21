package model

// EquipSlot 装备位置
type EquipSlot int

const (
	// Head 头部装备
	Head EquipSlot = iota
	HeadAcce
	Face
	FaceAcce
	ChestAcce
	UpperBody
	LowerBody
	Body
	RightHand
	LeftHand
	Shoes
	Socks
	Pet
)

// ItemType 道具类型
type ItemType int
