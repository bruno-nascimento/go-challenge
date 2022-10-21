package args

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"
)

//go:embed testdata/empty_args_err_msg.txt
var emptyArgsErrMsg string

//go:embed testdata/only_dir.txt
var onlyDirMsg string

//go:embed testdata/invalid_type.txt
var invalidType string

//go:embed testdata/only_dir_type.txt
var onlyDirType string

//go:embed testdata/invalid_start_time.txt
var invalidStartTime string

//go:embed testdata/onlyDirTypeStartTime.txt
var onlyDirTypeStartTime string

//go:embed testdata/end_before_start.txt
var endBeforeStart string

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		assertFn func(tt *testing.T, cfg *Config, err *ValidationError)
	}{
		{
			name: "no args",
			args: []string{"cmd"},
			assertFn: func(tt *testing.T, cfg *Config, err *ValidationError) {
				if emptyArgsErrMsg != err.Error() {
					tt.Error(fmt.Errorf("expected: '%s' got: '%s'", emptyArgsErrMsg, err.Error()))
				}
			},
		},
		{
			name: "only dir",
			args: []string{"cmd", "--directory", "/tmp"},
			assertFn: func(tt *testing.T, cfg *Config, err *ValidationError) {
				if onlyDirMsg != err.Error() {
					tt.Error(fmt.Errorf("expected: '%s' got: '%s'", emptyArgsErrMsg, err.Error()))
				}
			},
		},
		{
			name: "only short dir",
			args: []string{"cmd", "-d", "/tmp"},
			assertFn: func(tt *testing.T, cfg *Config, err *ValidationError) {
				if onlyDirMsg != err.Error() {
					tt.Error(fmt.Errorf("expected: '%s' got: '%s'", emptyArgsErrMsg, err.Error()))
				}
			},
		},
		{
			name: "only short dir and invalid type",
			args: []string{"cmd", "-d", "/dev/null", "-type", "toml"},
			assertFn: func(tt *testing.T, cfg *Config, err *ValidationError) {
				if invalidType != err.Error() {
					tt.Error(fmt.Errorf("expected: '%s' got: '%s'", emptyArgsErrMsg, err.Error()))
				}
			},
		},
		{
			name: "only short dir and invalid short type",
			args: []string{"cmd", "-d", "/dev/null", "-t", "toml"},
			assertFn: func(tt *testing.T, cfg *Config, err *ValidationError) {
				if invalidType != err.Error() {
					tt.Error(fmt.Errorf("expected: '%s' got: '%s'", emptyArgsErrMsg, err.Error()))
				}
			},
		},
		{
			name: "only short dir and short type",
			args: []string{"cmd", "-d", "/dev/null", "-t", "json"},
			assertFn: func(tt *testing.T, cfg *Config, err *ValidationError) {
				if onlyDirType != err.Error() {
					tt.Error(fmt.Errorf("expected: '%s' got: '%s'", emptyArgsErrMsg, err.Error()))
				}
			},
		},
		{
			name: "only short dir, short type and invalid startTime",
			args: []string{"cmd", "-d", "/dev/null", "-t", "json", "--startTime", "invalid"},
			assertFn: func(tt *testing.T, cfg *Config, err *ValidationError) {
				if invalidStartTime != err.Error() {
					tt.Error(fmt.Errorf("expected: '%s' got: '%s'", emptyArgsErrMsg, err.Error()))
				}
			},
		},
		{
			name: "only short dir, short type and startTime",
			args: []string{"cmd", "-d", "/dev/null", "-t", "json", "--startTime", "2022-01-01T00:00:00.00Z"},
			assertFn: func(tt *testing.T, cfg *Config, err *ValidationError) {
				if onlyDirTypeStartTime != err.Error() {
					tt.Error(fmt.Errorf("expected: '%s' got: '%s'", emptyArgsErrMsg, err.Error()))
				}
			},
		},
		{
			name: "only short dir, short type, startTime and endtime before starttime",
			args: []string{"cmd", "-d", "/dev/null", "-t", "json", "--startTime", "2022-01-01T00:00:00.00Z", "--endTime", "2021-01-01T00:00:00.00Z"},
			assertFn: func(tt *testing.T, cfg *Config, err *ValidationError) {
				if endBeforeStart != err.Error() {
					tt.Error(fmt.Errorf("expected: '%s' got: '%s'", emptyArgsErrMsg, err.Error()))
				}
			},
		},
		{
			name: "valid",
			args: []string{"cmd", "-d", "/dev/null", "-t", "json", "--startTime", "2022-01-01T00:00:00.00Z", "--endTime", "2023-01-01T00:00:00.00Z"},
			assertFn: func(tt *testing.T, cfg *Config, err *ValidationError) {
				if err != nil {
					tt.Error("error should be nil")
				}
				if cfg.Directory != "/dev/null" {
					tt.Error("wrong dir in cfg")
				}
				if cfg.Type != "json" {
					tt.Error("wrong type in cfg")
				}
				fmt.Println(">>>>>>>>>>> ", cfg.StartTime.Format(time.RFC3339))
				if cfg.StartTime.Format(time.RFC3339) != "2022-01-01T00:00:00Z" {
					tt.Error("wrong starttime in cfg")
				}
				if cfg.EndTime.Format(time.RFC3339) != "2023-01-01T00:00:00Z" {
					tt.Error("wrong endtime in cfg")
				}
				if cfg.OutPutFileName != "out" {
					tt.Error("wrong outputfile in cfg")
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = test.args
			config, err := Parse()
			test.assertFn(tt, config, err)
		})
	}
}
