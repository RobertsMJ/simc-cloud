package simc

// TODO
type Output string

func (o Output) MarshalSimC() ([]byte, error) {
	return []byte(o), nil
}

func (o *Output) UnmarshalSimC(data []byte) error {
	*o = Output(data)
	return nil
}
