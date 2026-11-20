package simc

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	fieldID                  = "id"
	fieldEnchantID           = "enchant_id"
	fieldBonusID             = "bonus_id"
	fieldGemID               = "gem_id"
	fieldContentTuning       = "content_tuning"
	fieldDropLevel           = "drop_level"
	fieldRedirectedBaseStats = "redirected_base_stats"
	fieldCraftedStats        = "crafted_stats"
	fieldGemBonusID          = "gem_bonus_id"
	fieldCraftingQuality     = "crafting_quality"
	fieldTitanDiscID         = "titan_disc_id"
)

type Item struct {
	Slot SlotGroup `json:"slot"`
	Name string    `json:"name,omitempty"`
	ID   int       `json:"id"`

	EnchantID           IDList `json:"enchant_id,omitempty"`
	BonusID             IDList `json:"bonus_id,omitempty"`
	GemID               IDList `json:"gem_id,omitempty"`
	ContentTuning       int    `json:"content_tuning,omitempty"`
	DropLevel           int    `json:"drop_level,omitempty"`
	RedirectedBaseStats int    `json:"redirected_base_stats,omitempty"`
	CraftedStats        IDList `json:"crafted_stats,omitempty"`
	GemBonusID          IDList `json:"gem_bonus_id,omitempty"`
	CraftingQuality     int    `json:"crafting_quality,omitempty"`
	TitanDiscID         int    `json:"titan_disc_id,omitempty"`
}

func newItem(s statement) (Item, error) {
	var i Item
	var err error
	var slot ItemSlot
	if slot, err = newItemSlot(s.Key); err != nil {
		return Item{}, err
	}
	i.Slot = slot.Group()
	i.Name = s.Comment

	for kv := range strings.SplitSeq(s.Value, ",") {
		k, v, _ := strings.Cut(kv, "=")
		switch k {
		case fieldID:
			i.ID, err = strconv.Atoi(v)
		case fieldEnchantID:
			i.EnchantID, err = newIDList(v)
		case fieldBonusID:
			i.BonusID, err = newIDList(v)
		case fieldGemID:
			i.GemID, err = newIDList(v)
		case fieldContentTuning:
			i.ContentTuning, err = strconv.Atoi(v)
		case fieldDropLevel:
			i.DropLevel, err = strconv.Atoi(v)
		case fieldRedirectedBaseStats:
			i.RedirectedBaseStats, err = strconv.Atoi(v)
		case fieldCraftedStats:
			i.CraftedStats, err = newIDList(v)
		case fieldGemBonusID:
			i.GemBonusID, err = newIDList(v)
		case fieldCraftingQuality:
			i.CraftingQuality, err = strconv.Atoi(v)
		case fieldTitanDiscID:
			i.TitanDiscID, err = strconv.Atoi(v)
		}
		if err != nil {
			return Item{}, fmt.Errorf("item %s: %w", k, err)
		}
	}

	return i, nil
}

func isItemKey(key string) bool {
	if _, err := newItemSlot(key); err == nil {
		return true
	}
	return false
}

func (i Item) render(slot ItemSlot) string {
	res := []string{
		fmt.Sprintf(`%s=`, slot),
		fmt.Sprintf(`id=%d`, i.ID),
	}

	if len(i.EnchantID) > 0 {
		res = append(res, fmt.Sprintf(`%s=%s`, fieldEnchantID, i.EnchantID.render()))
	}
	if len(i.BonusID) > 0 {
		res = append(res, fmt.Sprintf(`%s=%s`, fieldBonusID, i.BonusID.render()))
	}
	if len(i.GemID) > 0 {
		res = append(res, fmt.Sprintf(`%s=%s`, fieldGemID, i.GemID.render()))
	}
	if i.ContentTuning != 0 {
		res = append(res, fmt.Sprintf(`%s=%d`, fieldContentTuning, i.ContentTuning))
	}
	if i.DropLevel != 0 {
		res = append(res, fmt.Sprintf(`%s=%d`, fieldDropLevel, i.DropLevel))
	}
	if i.RedirectedBaseStats != 0 {
		res = append(res, fmt.Sprintf(`%s=%d`, fieldRedirectedBaseStats, i.RedirectedBaseStats))
	}
	if len(i.CraftedStats) > 0 {
		res = append(res, fmt.Sprintf(`%s=%s`, fieldCraftedStats, i.CraftedStats.render()))
	}
	if len(i.GemBonusID) > 0 {
		res = append(res, fmt.Sprintf(`%s=%s`, fieldGemBonusID, i.GemBonusID.render()))
	}
	if i.CraftingQuality != 0 {
		res = append(res, fmt.Sprintf(`%s=%d`, fieldCraftingQuality, i.CraftingQuality))
	}
	if i.TitanDiscID != 0 {
		res = append(res, fmt.Sprintf(`%s=%d`, fieldTitanDiscID, i.TitanDiscID))
	}

	return strings.Join(res, ",")
}
