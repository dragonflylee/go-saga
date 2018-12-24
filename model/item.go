package model

// EquipSlot 装备位置
type EquipSlot int

const (
	// SlotHead 头部装备
	SlotHead EquipSlot = iota
	SlotHeadAcce
	SlotFace
	SlotFaceAcce
	SlotChestAcce
	SlotUpperBody
	SlotLowerBody
	SlotBack
	SlotRightHand
	SlotLeftHand
	SlotShoes
	SlotSocks
	SlotPet
)

// SlotType 道具类型
type SlotType int
