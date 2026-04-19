package simc

type Input struct {
	Character Character                   `json:"character"`
	Equipment map[EquipmentSlot]Equipment `json:"equipment"`
	Options   Options                     `json:"options"`
}

var _ Marshaler = (*Input)(nil)
var _ Unmarshaler = (*Input)(nil)

func (i Input) MarshalSimC() ([]byte, error) {
	return Marshal(&i)
}

// UnmarshalSimC parses the simc input string into the Input struct.
func (i *Input) UnmarshalSimC(data []byte) error {
	return Unmarshal(data, i)
}
