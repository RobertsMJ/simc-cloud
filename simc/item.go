package simc

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"sync"

	"github.com/RobertsMJ/simc-cloud/simc/internal"
	"github.com/samber/lo"
)

var (
	ErrInvalidItemSlot = errors.New("invalid item slot")
	ErrInvalidValue    = errors.New("invalid value for type")
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

func ItemSlotFromString(s string) (ItemSlot, bool) {
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
		return ItemSlot(s), true
	default:
		return "", false
	}
}

type SlotGroup string

const (
	SlotGroupHead     ItemSlot  = "head"
	SlotGroupNeck     ItemSlot  = "neck"
	SlotGroupShoulder ItemSlot  = "shoulder"
	SlotGroupBack     ItemSlot  = "back"
	SlotGroupChest    ItemSlot  = "chest"
	SlotGroupWrist    ItemSlot  = "wrist"
	SlotGroupHands    ItemSlot  = "hands"
	SlotGroupWaist    ItemSlot  = "waist"
	SlotGroupLegs     ItemSlot  = "legs"
	SlotGroupFeet     ItemSlot  = "feet"
	SlotGroupFinger   SlotGroup = "finger"
	SlotGroupTrinket  SlotGroup = "trinket"
	SlotGroupMainHand ItemSlot  = "main_hand"
	SlotGroupOffHand  ItemSlot  = "off_hand"
)

func (is ItemSlot) GetSlotGroup() SlotGroup {
	switch is {
	case ItemSlotTrinket1, ItemSlotTrinket2:
		return SlotGroupTrinket
	case ItemSlotFinger1, ItemSlotFinger2:
		return SlotGroupFinger
	}
	return ""
}

type Item struct {
	Slot                ItemSlot        `json:"slot" simc:"-"`
	SlotGroup           SlotGroup       `json:"slot_group,omitempty" simc:"-"`
	Name                string          `json:"name,omitempty" simc:"-"`
	ID                  int             `json:"id" simc:"id"`
	EnchantID           internal.IDList `json:"enchant_id,omitempty" simc:"enchant_id,omitempty"`
	BonusID             internal.IDList `json:"bonus_id,omitempty" simc:"bonus_id,omitempty"`
	GemID               internal.IDList `json:"gem_id,omitempty" simc:"gem_id,omitempty"`
	ContentTuning       *int            `json:"content_tuning,omitempty" simc:"content_tuning,omitempty"`
	DropLevel           *int            `json:"drop_level,omitempty" simc:"drop_level,omitempty"`
	RedirectedBaseStats *int            `json:"redirected_base_stats,omitempty" simc:"redirected_base_stats,omitempty"`
	CraftedStats        internal.IDList `json:"crafted_stats,omitempty" simc:"crafted_stats,omitempty"`
	GemBonusID          internal.IDList `json:"gem_bonus_id,omitempty" simc:"gem_bonus_id,omitempty"`
	CraftingQuality     *int            `json:"crafting_quality,omitempty" simc:"crafting_quality,omitempty"`
	TitanDiscID         *int            `json:"titan_disc_id,omitempty" simc:"titan_disc_id,omitempty"`
}

func (i Item) GetSlotGroup() SlotGroup {
	return i.Slot.GetSlotGroup()
}

var (
	_ internal.ValueMarshaler       = Item{}
	_ internal.ValueUnmarshaler     = (*Item)(nil)
	_ internal.StatementUnmarshaler = (*Item)(nil)
)

var itemFieldIdx = sync.OnceValue(internal.IndexSimcFields[Item])

func isItemKey(key string) bool {
	_, ok := ItemSlotFromString(key)
	return ok
}

func (i *Item) UnmarshalStatement(s internal.Statement) error {
	if i == nil {
		return fmt.Errorf("character unmarshal: %w", internal.ErrUnmarshalIntoNilPtr)
	}

	slot, ok := ItemSlotFromString(s.Key)
	if !ok {
		return ErrInvalidItemSlot
	}
	i.Slot = slot
	i.SlotGroup = i.Slot.GetSlotGroup()
	i.Name = s.Comment

	vals := lo.Map(lo.Filter(lo.Map(strings.Split(s.Value, ","), func(v string, _ int) string {
		return strings.TrimSpace(v)
	}), func(v string, _ int) bool {
		return v != ""
	}), func(v string, _ int) Parameter {
		return NewParameter(v)
	})

	fields := itemFieldIdx()
	for _, v := range vals {
		idx, ok := fields[v.Key]
		if !ok {
			slog.Debug("attempted to assign unmapped field to an item", "key", v.Key, "value", v.Value)
			continue
		}

		field := reflect.ValueOf(i).Elem().Field(idx)
		if err := internal.AssignToField(field, v.Value); err != nil {
			return fmt.Errorf("item unmarshal into field(%v): %w", s.Key, err)
		}
	}
	return nil
}

func (i Item) MarshalSimcValue() (string, error) {
	panic("not implemented")
}

func (i *Item) UnmarshalSimcValue(value string) error {
	panic("not implemented")
}
