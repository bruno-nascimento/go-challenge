package generate

import (
	"fmt"
	"time"
)

type Metric struct {
	Timestamp time.Time `json:"timestamp"`
	LevelName string    `json:"level_name"`
	Value     int       `json:"value"`
}

func (m Metric) ToCSV() any {
	return fmt.Sprintf("%s,%s,%d\n", m.Timestamp.Format(time.RFC3339), m.LevelName, m.Value)
}

func (m Metric) ToJSON() any {
	return m
}

func (m Metric) ToYAML() any {
	return fmt.Sprintf("- level_name: %s\n  value: %d\n  timestamp: %s\n", m.LevelName, m.Value, m.Timestamp.Format(time.RFC3339))
}
