package models

type ItemSlot string

const (
	ItemSlotHead ItemSlot = "head"
	// TODO:MJR add the rest 4head
)

type Item struct {
	Slot            ItemSlot `json:"slot"`
	ID              int      `json:"id"`
	ItemLevel       *int     `json:"item_level,omitempty"`
	EnchantID       *[]int   `json:"enchant_id,omitempty"`
	BonusID         *[]int   `json:"bonus_id,omitempty"`
	GemID           *[]int   `json:"gem_id,omitempty"`
	CraftedStats    *[]int   `json:"crafted_stats,omitempty"`
	CraftingQuality *[]int   `json:"crafting_quality,omitempty"`
}
