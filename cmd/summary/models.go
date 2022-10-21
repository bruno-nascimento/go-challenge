package summary

import "fmt"

type Summary struct {
	Level string `json:"level_name"`
	Total int64  `json:"total_value"`
}

func (s Summary) ToJSON() any {
	return s
}

func (s Summary) ToYAML() any {
	return fmt.Sprintf("- level_name: %s\n  value: %d\n", s.Level, s.Total)
}
