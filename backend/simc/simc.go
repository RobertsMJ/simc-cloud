package simc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/RobertsMJ/simc-cloud/backend/simc/internal"
)

// SimulationProfile represents a full single-character simulation profile
type SimulationProfile struct {
	Character Character   `json:"character"`
	Baseline  Loadout     `json:"baseline"`
	Options   []Loadout   `json:"options"`
	Config    Config      `json:"config"`
	Extra     []Parameter `json:"extra,omitempty"`
}

var _ internal.ValueMarshaler = SimulationProfile{}

func (p SimulationProfile) MarshalSimc() (string, error) {
	character, err := p.Character.MarshalSimc()
	if err != nil {
		return "", errors.Join(err, errors.New("simulation profile: failed to marshal character"))
	}

	baseline, err := p.Baseline.MarshalSimc()
	if err != nil {
		return "", errors.Join(err, errors.New("simulation profile: failed to marshal baseline"))
	}

	options := make([]string, 0, len(p.Options))
	for _, opt := range p.Options {
		optStr, err := opt.MarshalSimc()
		if err != nil {
			return "", errors.Join(err, fmt.Errorf("simulation profile: failed to marshal option: %v", opt))
		}
		options = append(options, optStr)
	}
	optionsStr := strings.Join(options, "\n")

	config, err := p.Config.MarshalSimc()
	if err != nil {
		return "", errors.Join(err, errors.New("simulation profile: failed to marshal config"))
	}

	return strings.Join([]string{character, baseline, optionsStr, config}, "\n"), nil
}
