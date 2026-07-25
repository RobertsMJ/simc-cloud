package simc

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
	// TODO:MJR add the rest
)

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

type Item struct {
	Slot                ItemSlot  `json:"slot" simc:",key"`
	SlotGroup           SlotGroup `json:"slot_group,omitempty" simc:"-"`
	ID                  int       `json:"id" simc:"id"`
	EnchantID           IDList    `json:"enchant_id,omitempty" simc:"enchant_id,omitempty"`
	BonusID             IDList    `json:"bonus_id,omitempty" simc:"bonus_id,omitempty"`
	GemID               IDList    `json:"gem_id,omitempty" simc:"gem_id,omitempty"`
	ContentTuning       *int      `json:"content_tuning,omitempty" simc:"content_tuning,omitempty"`
	DropLevel           *int      `json:"drop_level,omitempty" simc:"drop_level,omitempty"`
	RedirectedBaseStats *int      `json:"redirected_base_stats,omitempty" simc:"redirected_base_stats,omitempty"`
	CraftedStats        IDList    `json:"crafted_stats,omitempty" simc:"crafted_stats,omitempty"`
	GemBonusID          IDList    `json:"gem_bonus_id,omitempty" simc:"gem_bonus_id,omitempty"`
	CraftingQuality     *int      `json:"crafting_quality,omitempty" simc:"crafting_quality,omitempty"`
	TitanDiscID         *int      `json:"titan_disc_id,omitempty" simc:"titan_disc_id,omitempty"`
}

var (
	_ ValueMarshaler   = Item{}
	_ ValueUnmarshaler = (*Item)(nil)
)

func (i Item) MarshalSimcValue() (string, error) {
	panic("not implemented")
}

func (i Item) UnmarshalSimcValue(value string) error {
	panic("not implemented")
}
