package reader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"go-challenge/internal/args"
)

type CSV struct {
	config  *args.Config
	decoder *csv.Reader
}

func NewCSV(config *args.Config) *CSV {
	return &CSV{config: config}
}

func (c CSV) File(file *os.File) (Decoder, error) {
	newDecoder := NewCSV(c.config)
	newDecoder.decoder = csv.NewReader(file)
	// skipping header
	_, _ = newDecoder.decoder.Read()
	return newDecoder, nil
}

func (c CSV) ParseAll(sendParsedMetric func(item Metric)) {
	for {
		line, err := c.decoder.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			panic(err)
		}
		metric, err := c.csvToMetric(line)
		if err != nil {
			fmt.Printf("WARNING - error parsing csv line to metrics: '%s'; err: %s", line, err.Error())
		}
		sendParsedMetric(*metric)
	}
}

func (c CSV) csvToMetric(line []string) (*Metric, error) {
	timeStamp, err := time.Parse(time.RFC3339, line[0])
	if err != nil {
		return nil, err
	}
	value, err := strconv.Atoi(line[2])
	if err != nil {
		return nil, err
	}
	return &Metric{
		Timestamp: timeStamp,
		LevelName: line[1],
		Value:     value,
	}, nil
}
