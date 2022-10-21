package cmd

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"go-challenge/internal/args"
)

//go:embed testdata/1/out.yaml
var assertYamlOutput []byte

func newConfig(dir, inputType, outputType string) *args.Config {
	cfg := &args.Config{
		Directory: dir,
		Type:      inputType,
		StartTime: func() time.Time {
			ts, _ := time.Parse(time.RFC3339, "2022-01-01T00:00:00.00Z")
			return ts
		}(),
		EndTime: func() time.Time {
			ts, _ := time.Parse(time.RFC3339, "2022-02-01T00:00:00.01Z")
			return ts
		}(),
		OutPutFileType:    outputType,
		OutPutFileName:    "out",
		GenerateTestFiles: false,
	}
	return cfg
}

func TestSummaryCMD(t *testing.T) {
	assertMap := map[string]int{"level_4": 2727, "level_7": 2114, "level_2": 1993, "level_9": 1783, "level_10": 1777, "level_8": 916, "level_5": 814, "level_1": 516, "level_3": 379, "level_6": 297}

	tests := []struct {
		name        string
		config      *args.Config
		outFileName string
		assertFn    func(t *testing.T, outFileContent []byte)
	}{
		{
			name:        "json input -> json output",
			config:      newConfig("testdata/1/json", "json", "json"),
			outFileName: "out.json",
			assertFn: func(ttt *testing.T, outFileContent []byte) {
				var results []Result
				err := json.Unmarshal(outFileContent, &results)
				if err != nil {
					ttt.Error(err)
				}
				for _, r := range results {
					if assertMap[r.LevelName] != r.TotalValue {
						ttt.Error(fmt.Sprintf("expected: '%d' got: '%d' for level: '%s'", assertMap[r.LevelName], r.TotalValue, r.LevelName))
					}
				}
			},
		},
		{
			name:        "csv input -> json output",
			config:      newConfig("testdata/1/csv", "csv", "json"),
			outFileName: "out.json",
			assertFn: func(ttt *testing.T, outFileContent []byte) {
				var results []Result
				err := json.Unmarshal(outFileContent, &results)
				if err != nil {
					ttt.Error(err)
				}
				for _, r := range results {
					if assertMap[r.LevelName] != r.TotalValue {
						ttt.Error(fmt.Sprintf("expected: '%d' got: '%d' for level: '%s'", assertMap[r.LevelName], r.TotalValue, r.LevelName))
					}
				}
			},
		},
		{
			name:        "csv input -> yaml output",
			config:      newConfig("testdata/1/csv", "csv", "yaml"),
			outFileName: "out.yaml",
			assertFn: func(ttt *testing.T, outFileContent []byte) {
				if bytes.Compare(outFileContent, assertYamlOutput) != 0 {
					ttt.Error("golden file does not match with the generate file")
				}
			},
		},
		{
			name:        "json input -> yaml output",
			config:      newConfig("testdata/1/json", "json", "yaml"),
			outFileName: "out.yaml",
			assertFn: func(ttt *testing.T, outFileContent []byte) {
				if bytes.Compare(outFileContent, assertYamlOutput) != 0 {
					ttt.Error("golden file does not match with the generate file")
				}
			},
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(tt *testing.T) {
			Run(tst.config)
			outFileContent, err := ioutil.ReadFile(tst.outFileName)
			if err != nil {
				t.Error(err)
			}
			tst.assertFn(tt, outFileContent)
		})
	}
}

type Result struct {
	LevelName  string `json:"level_name"`
	TotalValue int    `json:"total_value"`
}
