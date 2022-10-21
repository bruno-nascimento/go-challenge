package summary

import (
	"fmt"
	"sort"

	"go-challenge/internal/args"
	"go-challenge/internal/printer"
	"go-challenge/internal/reader"
	"go-challenge/internal/writer"
)

func Cmd(config *args.Config) {
	printer.Header()
	r := reader.NewReader(config, GetDecoder(config))
	err := r.Read()
	if err != nil {
		println(err.Error())
		return
	}
	sums := summaryList(r)
	encFn := typeEncoderCfg[config.OutPutFileType]
	wr := writer.NewWriter(config, encFn.Encoder)
	table := printer.NewTableBuilder()
	for _, sum := range sums {
		err := wr.Write(encFn.Fn(sum))
		if err != nil {
			panic(err)
		}
		table.AddRow(sum.Level, sum.Total)
	}
	err = wr.Close()
	if err != nil {
		panic(err)
	}
	println(table.Build())
}

func summaryList(r *reader.Reader) []Summary {
	var sums []Summary
	for k, v := range r.M {
		sums = append(sums, Summary{Level: k, Total: v})
	}
	sort.Slice(sums[:], func(i, j int) bool {
		return sums[i].Total > sums[j].Total
	})
	return sums
}

func GetDecoder(config *args.Config) reader.Decoder {
	if config.Type == "json" {
		return reader.NewJSON(config)
	}
	if config.Type == "csv" {
		return reader.NewCSV(config)
	}
	panic(fmt.Sprint("error getting decoder for invalid type: ", config.Type))
}

type EncoderFN struct {
	Encoder writer.Encoder
	Fn      func(sum Summary) any
}

var typeEncoderCfg = map[string]EncoderFN{
	"json": {
		Encoder: writer.NewJSON(), Fn: func(sum Summary) any {
			return sum.ToJSON()
		},
	},
	"yaml": {
		Encoder: writer.NewYAML(), Fn: func(sum Summary) any {
			return sum.ToYAML()
		},
	},
}
