package simc

type Input struct {
	Character Character                   `json:"character"`
	Equipment map[EquipmentSlot]Equipment `json:"equipment"`
	Options   Options                     `json:"options"`
}

var _ SimCMarshaler = (*Input)(nil)
var _ SimCUnmarshaler = (*Input)(nil)

func (i Input) MarshalSimC() ([]byte, error) {
	return Marshal(i)
}

func (i *Input) UnmarshalSimC(data []byte) error {
	input, err := Unmarshal(data)
	if err != nil {
		return err
	}
	*i = input
	return nil
}
