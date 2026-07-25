package simc

type SimConfig struct {
	DesiredTargets int
	MaxTime        int
	ReportDetails  int
	Iterations     int
}

var (
	_ ValueMarshaler   = SimConfig{}
	_ ValueUnmarshaler = (*SimConfig)(nil)
)

func (sc SimConfig) MarshalSimcValue() (string, error) {
	panic("not implemented")
}

func (sc SimConfig) UnmarshalSimcValue(value string) error {
	panic("not implemented")
}
