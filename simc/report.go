package simc

import (
	"encoding/json"
	"fmt"
	"io"
)

type Report struct {
	Sim struct {
		Statistics Statistics `json:"statistics"`
	} `json:"sim"`
}

func ParseReport(r io.Reader) (Report, error) {
	var report Report
	if err := json.NewDecoder(r).Decode(&report); err != nil {
		return Report{}, fmt.Errorf("decode simc report: %w", err)
	}
	return report, nil
}

type Statistics struct {
	ElapsedCpuSeconds    float64    `json:"elapsed_cpu_seconds"`
	ElapsedTimeSeconds   float64    `json:"elapsed_time_seconds"`
	InitTimeSeconds      float64    `json:"init_time_seconds"`
	MergeTimeSeconds     float64    `json:"merge_time_seconds"`
	AnalyzeTimeSeconds   float64    `json:"analyze_time_seconds"`
	TotalEventsProcessed int64      `json:"total_events_processed"`
	SimulationLength     StatSample `json:"simulation_length"`
	RaidDps              StatSample `json:"raid_dps"`
	TotalDmg             StatSample `json:"total_dmg"`
}

type StatSample struct {
	Sum          float64 `json:"sum"`
	Count        int     `json:"count"`
	Mean         float64 `json:"mean"`
	Min          float64 `json:"min,omitempty"`
	Max          float64 `json:"max,omitempty"`
	Median       float64 `json:"median,omitempty"`
	Variance     float64 `json:"variance,omitempty"`
	StdDev       float64 `json:"std_dev,omitempty"`
	MeanVariance float64 `json:"mean_variance,omitempty"`
	MeanStdDev   float64 `json:"mean_std_dev,omitempty"`
}
