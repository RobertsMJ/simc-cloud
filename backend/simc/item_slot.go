package simc

import (
	"encoding"
	"fmt"
)

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

var _ encoding.TextUnmarshaler = (*ItemSlot)(nil)

func (i *ItemSlot) UnmarshalText(text []byte) error {
	slot, err := newItemSlot(string(text))
	if err != nil {
		return err
	}
	*i = slot
	return nil
}

func newItemSlot(s string) (ItemSlot, error) {
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

func (i ItemSlot) Group() SlotGroup {
	switch i {
	case ItemSlotTrinket1, ItemSlotTrinket2:
		return SlotGroupTrinket
	case ItemSlotFinger1, ItemSlotFinger2:
		return SlotGroupFinger
	}
	return SlotGroup(i)
}

type SlotGroup string

var _ encoding.TextUnmarshaler = (*SlotGroup)(nil)

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

func (s *SlotGroup) UnmarshalText(text []byte) error {
	switch SlotGroup(text) {
	case SlotGroupHead,
		SlotGroupNeck,
		SlotGroupShoulder,
		SlotGroupBack,
		SlotGroupChest,
		SlotGroupWrist,
		SlotGroupHands,
		SlotGroupWaist,
		SlotGroupLegs,
		SlotGroupFeet,
		SlotGroupFinger,
		SlotGroupTrinket,
		SlotGroupMainHand,
		SlotGroupOffHand:
		*s = SlotGroup(text)
		return nil
	default:
		return fmt.Errorf("%s is not a slot group", text)
	}
}

func (s SlotGroup) Slots() []ItemSlot {
	switch s {
	case SlotGroupTrinket:
		return []ItemSlot{ItemSlotTrinket1, ItemSlotTrinket2}
	case SlotGroupFinger:
		return []ItemSlot{ItemSlotFinger1, ItemSlotFinger2}
	case SlotGroupHead:
		return []ItemSlot{ItemSlotHead}
	case SlotGroupNeck:
		return []ItemSlot{ItemSlotNeck}
	case SlotGroupShoulder:
		return []ItemSlot{ItemSlotShoulder}
	case SlotGroupBack:
		return []ItemSlot{ItemSlotBack}
	case SlotGroupChest:
		return []ItemSlot{ItemSlotChest}
	case SlotGroupWrist:
		return []ItemSlot{ItemSlotWrist}
	case SlotGroupHands:
		return []ItemSlot{ItemSlotHands}
	case SlotGroupWaist:
		return []ItemSlot{ItemSlotWaist}
	case SlotGroupLegs:
		return []ItemSlot{ItemSlotLegs}
	case SlotGroupFeet:
		return []ItemSlot{ItemSlotFeet}
	case SlotGroupMainHand:
		return []ItemSlot{ItemSlotMainHand}
	case SlotGroupOffHand:
		return []ItemSlot{ItemSlotOffHand}
	default:
		return nil
	}
}
