package simc

import (
	"errors"
	"reflect"
	"strings"
)

type EquipmentSlot string

// https://github.com/simulationcraft/simc/wiki/Equipment#slots-keywords
const (
	EquipmentSlotMetaGem   EquipmentSlot = "meta_gem"
	EquipmentSlotHead      EquipmentSlot = "head"
	EquipmentSlotNeck      EquipmentSlot = "neck"
	EquipmentSlotShoulder  EquipmentSlot = "shoulder"
	EquipmentSlotShoulders EquipmentSlot = "shoulders"
	EquipmentSlotBack      EquipmentSlot = "back"
	EquipmentSlotShirt     EquipmentSlot = "shirt"
	EquipmentSlotChest     EquipmentSlot = "chest"
	EquipmentSlotWaist     EquipmentSlot = "waist"
	EquipmentSlotWrist     EquipmentSlot = "wrist"
	EquipmentSlotWrists    EquipmentSlot = "wrists"
	EquipmentSlotHand      EquipmentSlot = "hand"
	EquipmentSlotHands     EquipmentSlot = "hands"
	EquipmentSlotLegs      EquipmentSlot = "legs"
	EquipmentSlotFeet      EquipmentSlot = "feet"
	EquipmentSlotFinger1   EquipmentSlot = "finger1"
	EquipmentSlotRing1     EquipmentSlot = "ring1"
	EquipmentSlotFinger2   EquipmentSlot = "finger2"
	EquipmentSlotRing2     EquipmentSlot = "ring2"
	EquipmentSlotTrinket1  EquipmentSlot = "trinket1"
	EquipmentSlotTrinket2  EquipmentSlot = "trinket2"
	EquipmentSlotMainHand  EquipmentSlot = "main_hand"
	EquipmentSlotOffHand   EquipmentSlot = "off_hand"
	EquipmentSlotRanged    EquipmentSlot = "ranged"
	EquipmentSlotTabard    EquipmentSlot = "tabard"
)

var validSlots = map[string]struct{}{
	string(EquipmentSlotMetaGem):   {},
	string(EquipmentSlotHead):      {},
	string(EquipmentSlotNeck):      {},
	string(EquipmentSlotShoulder):  {},
	string(EquipmentSlotShoulders): {},
	string(EquipmentSlotBack):      {},
	string(EquipmentSlotShirt):     {},
	string(EquipmentSlotChest):     {},
	string(EquipmentSlotWaist):     {},
	string(EquipmentSlotWrist):     {},
	string(EquipmentSlotWrists):    {},
	string(EquipmentSlotHand):      {},
	string(EquipmentSlotHands):     {},
	string(EquipmentSlotLegs):      {},
	string(EquipmentSlotFeet):      {},
	string(EquipmentSlotFinger1):   {},
	string(EquipmentSlotRing1):     {},
	string(EquipmentSlotFinger2):   {},
	string(EquipmentSlotRing2):     {},
	string(EquipmentSlotTrinket1):  {},
	string(EquipmentSlotTrinket2):  {},
	string(EquipmentSlotMainHand):  {},
	string(EquipmentSlotOffHand):   {},
	string(EquipmentSlotRanged):    {},
	string(EquipmentSlotTabard):    {},
}

// Check if a string is a line defining equipment by checking for the presence of a valid slot keyword and an id= field
func IsEquipmentValue(v string) bool {
	keyword := strings.Split(v, "=")[0] // Take the part before the first '=' as the keyword
	_, exists := validSlots[keyword]
	return exists
}

type Equipment struct {
	Slot            EquipmentSlot `json:"slot" simc:"keyword"`
	ItemName        string        `json:"item_name" simc:"-"`
	ID              int           `json:"id" simc:"id"`
	EnchantID       *int          `json:"enchant_id" simc:"enchant_id"`
	GemID           *int          `json:"gem_id" simc:"gem_id"`
	BonusID         *[]int        `json:"bonus_id" simc:"bonus_id"`
	CraftedStats    *[]int        `json:"crafted_stats" simc:"crafted_stats"`
	CraftingQuality *int          `json:"crafting_quality" simc:"crafting_quality"`
}

var _ Marshaler = (*Equipment)(nil)   // Ensure SimcEquipment implements SimCUnmarshaler
var _ Unmarshaler = (*Equipment)(nil) // Ensure SimcEquipment implements SimCUnmarshaler

func (e Equipment) MarshalSimC() ([]byte, error) {
	parts := []string{}
	val := reflect.ValueOf(e)
	typ := reflect.TypeOf(e)
	// Special case for the slot and item name, which are combined in the first part of the string, e.g. "head=some_helm"
	parts = append(parts, string(e.Slot)+"="+e.ItemName)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("simc")
		if tag == "" || tag == "keyword" || tag == "-" {
			continue
		}
		part, err := marshalField(field, tag)
		if err != nil {
			return nil, err
		}
		if len(part) > 0 {
			parts = append(parts, string(part))
		}
	}
	return []byte(strings.Join(parts, ",")), nil
}

func (e *Equipment) UnmarshalSimC(data []byte) error {
	parts := strings.Split(string(data), ",")
	if len(parts) == 0 {
		return errors.New("invalid equipment format")
	}
	// The first part should be the slot keyword + equipment piece name, e.g. "head=some_helm"
	slotPart := strings.Split(parts[0], "=")
	if len(slotPart) != 2 {
		return errors.New("invalid equipment format")
	}
	if _, exists := validSlots[slotPart[0]]; !exists {
		return errors.New("invalid equipment slot: " + slotPart[0])
	}
	e.Slot = EquipmentSlot(slotPart[0])
	e.ItemName = slotPart[1]

	val := reflect.ValueOf(e).Elem()
	typ := val.Type()
	for _, part := range parts[1:] {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key, value := kv[0], kv[1]
		for i := 0; i < val.NumField(); i++ {
			fieldType := typ.Field(i)
			if fieldType.Tag.Get("simc") == key {
				field := val.Field(i)
				if err := unmarshalField(field, key, value); err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}
