package simc

import "fmt"

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
	return fmt.Sprintf(`
desired_targets=%d
max_time=%d
report_details=%d
iterations=%d
`,
		sc.DesiredTargets, sc.MaxTime, sc.ReportDetails, sc.Iterations), nil
}

func (sc SimConfig) UnmarshalSimcValue(value string) error {
	panic("not implemented")
}
