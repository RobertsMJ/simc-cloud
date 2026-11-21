package gearset

import (
	"log/slog"
	"maps"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/simc"
)

const (
	slotTalent simc.SlotGroup = "talent"
	slotOmnium simc.SlotGroup = "omnium"
)

type slotRequirement struct {
	Slot  simc.SlotGroup
	Count int
}

var slotOrder = func() []slotRequirement {
	var order []slotRequirement
	for group, slots := range simc.SlotsByGroup {
		order = append(order, slotRequirement{
			Slot:  group,
			Count: len(slots),
		})
	}
	order = append(order, slotRequirement{
		Slot:  slotTalent,
		Count: 1,
	})
	order = append(order, slotRequirement{
		Slot:  slotOmnium,
		Count: 1,
	})
	return order
}()

type slotCandidate struct {
	slot         simc.SlotGroup
	item         *simc.Item
	talent       *simc.Talent
	omniumTalent *simc.OmniumTalents
}

type Service struct {
	log *slog.Logger
}

func NewService(log *slog.Logger) *Service {
	return &Service{
		log: log,
	}
}

// TODO: parse, don't validate
func (s *Service) CountGearsets(baseline simc.Loadout, options models.Candidates) int {
	// do math
	// simply multiple number of options for each slot
	// trinkets and rings are n choose 2
	// (item sets) * talents * omniums
	countBySlot := map[simc.SlotGroup]int{}
	for _, item := range options.Items {
		countBySlot[item.Slot]++
	}

	total := 1
	for slot, count := range countBySlot {
		if count == 0 {
			continue
		}

		if slot == simc.SlotGroupTrinket || slot == simc.SlotGroupFinger {
			// n choose 2 for trinkets and rings
			// add 2 to count to account for the baseline items in these slots
			count = count + 2
			total *= (count * (count - 1) / 2)
		} else {
			count = count + 1
			total *= count
		}
	}

	if len(options.Talents) > 0 {
		total *= len(options.Talents) + 1
	}
	if len(options.OmniumTalents) > 0 {
		total *= len(options.OmniumTalents) + 1
	}
	return total
}

// GenerateGearsets does not include the baseline loadout in the returned gearsets
func (s *Service) GenerateGearsets(baseline simc.Loadout, options models.Candidates) []simc.Loadout {
	bySlot := map[simc.SlotGroup][]slotCandidate{}
	// add candidates
	for i, item := range options.Items {
		// there is a candidate for this item slot - initialize with the baseline
		if len(bySlot[item.Slot]) == 0 {
			bySlot[item.Slot] = []slotCandidate{}
			// handle multi-slot items (e.g., rings, trinkets)
			for _, slot := range item.Slot.Slots() {
				equipped := baseline.Items[slot]
				bySlot[item.Slot] = append(bySlot[item.Slot], slotCandidate{slot: item.Slot, item: &equipped})
			}
		}
		bySlot[item.Slot] = append(bySlot[item.Slot], slotCandidate{slot: item.Slot, item: &options.Items[i]})
	}

	// if there are any talent options, initialize the slot with the baseline talent
	if len(options.Talents) > 0 {
		bySlot[slotTalent] = []slotCandidate{
			{slot: slotTalent, talent: &baseline.Talents},
		}
		for _, talent := range options.Talents {
			bySlot[slotTalent] = append(bySlot[slotTalent], slotCandidate{slot: slotTalent, talent: &talent})
		}
	}

	// if there are any omnium talent options, initialize the slot with the baseline omnium talent
	if len(options.OmniumTalents) > 0 {
		bySlot[slotOmnium] = []slotCandidate{
			{slot: slotOmnium, omniumTalent: &baseline.OmniumTalents},
		}
		for _, omnium := range options.OmniumTalents {
			bySlot[slotOmnium] = append(bySlot[slotOmnium], slotCandidate{slot: slotOmnium, omniumTalent: &omnium})
		}
	}

	var results []simc.Loadout
	current := simc.Loadout{
		Items:         maps.Clone(baseline.Items),
		Talents:       baseline.Talents,
		OmniumTalents: baseline.OmniumTalents,
	}

	var backtrack func(slotIdx, startIdx, picked int)
	backtrack = func(slotIdx, startIdx, picked int) {
		if slotIdx >= len(slotOrder) {
			results = append(results, simc.Loadout{
				Items:         maps.Clone(current.Items),
				Talents:       current.Talents,
				OmniumTalents: current.OmniumTalents,
			})
			return
		}

		req := slotOrder[slotIdx]
		candidates := bySlot[req.Slot]

		// no options for this slot, move on
		if len(candidates) == 0 {
			backtrack(slotIdx+1, 0, 0)
		}

		// if we have picked enough items for this slot, move to the next slot
		if picked >= req.Count {
			backtrack(slotIdx+1, 0, 0)
			return
		}

		for i := startIdx; i < len(candidates); i++ {
			candidate := candidates[i]
			switch candidate.slot {
			case slotTalent:
				if candidate.talent == nil {
					s.log.Warn("nil talent candidate") // TODO: need more info, this shouldn't happen
					continue
				}
				current.Talents = *candidate.talent
			case slotOmnium:
				if candidate.omniumTalent == nil {
					s.log.Warn("nil omnium talent candidate") // TODO: need more info, this shouldn't happen
					continue
				}
				current.OmniumTalents = *candidate.omniumTalent
			case simc.SlotGroupFinger:
				if candidate.item == nil {
					s.log.Warn("nil finger item candidate") // TODO: need more info, this shouldn't happen
					continue
				}
				item := *candidate.item
				slots := item.Slot.Slots()
				slot := slots[picked]
				current.Items[slot] = item
			case simc.SlotGroupTrinket:
				if candidate.item == nil {
					s.log.Warn("nil trinket item candidate") // TODO: need more info, this shouldn't happen
					continue
				}
				item := *candidate.item
				slots := item.Slot.Slots()
				slot := slots[picked]
				current.Items[slot] = item
			default:
				if candidate.item == nil {
					s.log.Warn("nil item candidate") // TODO: need more info, this shouldn't happen
					continue
				}
				item := *candidate.item
				current.Items[item.Slot.Slots()[0]] = item
			}
			backtrack(slotIdx, i+1, picked+1)
		}
	}

	backtrack(0, 0, 0)

	return results
}
