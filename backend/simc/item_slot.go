package simc

import "fmt"

type ItemSlot string

const (
	ItemSlotHead     ItemSlot = "head"
	ItemSlotNeck     ItemSlot = "neck"
	ItemSlotShoulder ItemSlot = "shoulder"
	ItemSlotBack     ItemSlot = "back"
	ItemSlotChest    ItemSlot = "chest"
	ItemSlotWrist    ItemSlot = "wrist"
	ItemSlotHands    ItemSlot = "hands"
	ItemSlotWaist    ItemSlot = "waist"
	ItemSlotLegs     ItemSlot = "legs"
	ItemSlotFeet     ItemSlot = "feet"
	ItemSlotFinger1  ItemSlot = "finger1"
	ItemSlotFinger2  ItemSlot = "finger2"
	ItemSlotTrinket1 ItemSlot = "trinket1"
	ItemSlotTrinket2 ItemSlot = "trinket2"
	ItemSlotMainHand ItemSlot = "main_hand"
	ItemSlotOffHand  ItemSlot = "off_hand"
)

func NewItemSlot(s string) (ItemSlot, error) {
	switch ItemSlot(s) {
	case ItemSlotHead,
		ItemSlotNeck,
		ItemSlotShoulder,
		ItemSlotBack,
		ItemSlotChest,
		ItemSlotWrist,
		ItemSlotHands,
		ItemSlotWaist,
		ItemSlotLegs,
		ItemSlotFeet,
		ItemSlotFinger1,
		ItemSlotFinger2,
		ItemSlotTrinket1,
		ItemSlotTrinket2,
		ItemSlotMainHand,
		ItemSlotOffHand:
		return ItemSlot(s), nil
	default:
		return "", fmt.Errorf("%s is not an item slot", s)
	}
}

func (is ItemSlot) Group() SlotGroup {
	switch is {
	case ItemSlotTrinket1, ItemSlotTrinket2:
		return SlotGroupTrinket
	case ItemSlotFinger1, ItemSlotFinger2:
		return SlotGroupFinger
	}
	return SlotGroup(is)
}

type SlotGroup string

const (
	SlotGroupHead     SlotGroup = SlotGroup(ItemSlotHead)
	SlotGroupNeck     SlotGroup = SlotGroup(ItemSlotNeck)
	SlotGroupShoulder SlotGroup = SlotGroup(ItemSlotShoulder)
	SlotGroupBack     SlotGroup = SlotGroup(ItemSlotBack)
	SlotGroupChest    SlotGroup = SlotGroup(ItemSlotChest)
	SlotGroupWrist    SlotGroup = SlotGroup(ItemSlotWrist)
	SlotGroupHands    SlotGroup = SlotGroup(ItemSlotHands)
	SlotGroupWaist    SlotGroup = SlotGroup(ItemSlotWaist)
	SlotGroupLegs     SlotGroup = SlotGroup(ItemSlotLegs)
	SlotGroupFeet     SlotGroup = SlotGroup(ItemSlotFeet)
	SlotGroupFinger   SlotGroup = "finger"
	SlotGroupTrinket  SlotGroup = "trinket"
	SlotGroupMainHand SlotGroup = SlotGroup(ItemSlotMainHand)
	SlotGroupOffHand  SlotGroup = SlotGroup(ItemSlotOffHand)
)
