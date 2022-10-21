package args

import (
	"flag"
	"fmt"
	"time"

	"go-challenge/internal/printer"
)

type Config struct {
	Directory         string
	Type              string
	StartTime         time.Time
	EndTime           time.Time
	OutPutFileType    string
	OutPutFileName    string
	GenerateTestFiles bool
	strStartTime      string
	strEndTime        string
}

func Parse() (*Config, *ValidationError) {
	cfg := &Config{}
	loadFlag(&cfg.Directory, "d", "directory", "")
	loadFlag(&cfg.Type, "t", "type", "")
	loadFlag(&cfg.strStartTime, "", "startTime", "")
	loadFlag(&cfg.strEndTime, "", "endTime", "")
	loadFlag(&cfg.OutPutFileType, "", "outputFileType", "json")
	loadFlag(&cfg.OutPutFileName, "", "outputFileName", "out")
	flag.BoolVar(&cfg.GenerateTestFiles, "generate", false, "")
	flag.Usage = help
	flag.Parse()
	if err := validate(cfg); err != nil {
		printer.Header()
		println(err.Error())
		return nil, err
	}
	return cfg, nil
}

func loadFlag(holder *string, short, long, value string) {
	if short != "" {
		flag.StringVar(holder, short, value, "")
	}
	flag.StringVar(holder, long, value, "")
}

func help() {
	printer.Header()
	fmt.Printf(`
--directory	-d	Required	Directory path, the directory contains single type of file, it can be csv or json
--type		-t	Required	Type of the input files, supported format: json and csv
--startTime		Required	Starting time to scan the data in the format of rfc3339, inclusive
--endTime		Required	Ending time to scan the data in the format of rfc3339, exclusive
--outputFileName	Optional	Name of the output file of the summary. Default is 'out'
--outputFileType	Optional	Output type of the summary, supported value: json(default) and yaml
--generate		Optional	Generate test files with random values following the configuration from the other args.
--help		-h			Prints help information 
`)
}
