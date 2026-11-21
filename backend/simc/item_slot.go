package simc

import (
	"encoding"
	"fmt"
)

type ItemSlot string

var _ encoding.TextUnmarshaler = (*ItemSlot)(nil)

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

var SlotsByGroup = map[SlotGroup][]ItemSlot{
	SlotGroupHead:     {ItemSlotHead},
	SlotGroupNeck:     {ItemSlotNeck},
	SlotGroupShoulder: {ItemSlotShoulder},
	SlotGroupBack:     {ItemSlotBack},
	SlotGroupChest:    {ItemSlotChest},
	SlotGroupWrist:    {ItemSlotWrist},
	SlotGroupHands:    {ItemSlotHands},
	SlotGroupWaist:    {ItemSlotWaist},
	SlotGroupLegs:     {ItemSlotLegs},
	SlotGroupFeet:     {ItemSlotFeet},
	SlotGroupFinger:   {ItemSlotFinger1, ItemSlotFinger2},
	SlotGroupTrinket:  {ItemSlotTrinket1, ItemSlotTrinket2},
	SlotGroupMainHand: {ItemSlotMainHand},
	SlotGroupOffHand:  {ItemSlotOffHand},
}

var groupBySlot = func() map[ItemSlot]SlotGroup {
	m := make(map[ItemSlot]SlotGroup)
	for g, slots := range SlotsByGroup {
		for _, s := range slots {
			m[s] = g
		}
	}
	return m
}()

func newItemSlot(s string) (ItemSlot, error) {
	_, ok := groupBySlot[ItemSlot(s)]
	if !ok {
		return "", fmt.Errorf("%s is not an item slot", s)
	}
	return ItemSlot(s), nil
}

func (i *ItemSlot) UnmarshalText(text []byte) error {
	if _, ok := groupBySlot[ItemSlot(text)]; !ok {
		return fmt.Errorf("%s is not an item slot", text)
	}
	*i = ItemSlot(text)
	return nil
}

func (i ItemSlot) Group() SlotGroup {
	return groupBySlot[i]
}

func (s *SlotGroup) UnmarshalText(text []byte) error {
	if _, ok := SlotsByGroup[SlotGroup(text)]; !ok {
		return fmt.Errorf("%s is not a slot group", text)
	}
	*s = SlotGroup(text)
	return nil
}

func (s SlotGroup) Slots() []ItemSlot {
	return SlotsByGroup[s]
}
