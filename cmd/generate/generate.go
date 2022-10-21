package generate

import (
	"fmt"
	"math/rand"
	"time"

	"go-challenge/internal/args"
	"go-challenge/internal/printer"
	"go-challenge/internal/writer"
)

func Cmd(config *args.Config) {
	printer.Header()
	fmt.Print("generating test random files ...\n")
	encFn := typeEncoderCfg[config.Type]
	wr := writer.NewWriter(config, encFn.Encoder, config.StartTime)
	currDay := config.StartTime.Day()
	m := make(map[string]int64)
	for current := config.StartTime; current.Before(config.EndTime); current = current.Add(randomMillis()) {
		metrics := NewRandomMetric(current)
		incLevel(metrics, m)
		if currDay != metrics.Timestamp.Day() {
			wr.CreateFile(current)
			currDay = metrics.Timestamp.Day()
		}
		err := wr.Write(encFn.Fn(metrics))
		if err != nil {
			panic(err)
		}
	}
	err := wr.Close()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", m)
}

func NewRandomMetric(timestamp time.Time) Metric {
	return Metric{
		Timestamp: timestamp,
		LevelName: randomLevel(),
		Value:     int(randomMillis().Milliseconds()),
	}
}

func incLevel(metric Metric, m map[string]int64) {
	m[metric.LevelName] = m[metric.LevelName] + int64(metric.Value)
	m["total"] = m["total"] + int64(metric.Value)
}

func randomLevel() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("level_%d", rand.Intn(10)+1)
}

func randomMillis() time.Duration {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(1000)
	return time.Duration(i) * time.Millisecond
}

type EncoderFN struct {
	Encoder writer.Encoder
	Fn      func(metric Metric) any
}

var typeEncoderCfg = map[string]EncoderFN{
	"json": {
		Encoder: writer.NewJSON(), Fn: func(metric Metric) any {
			return metric.ToJSON()
		},
	},
	"csv": {
		Encoder: writer.NewCSV(), Fn: func(metric Metric) any {
			return metric.ToCSV()
		},
	},
	"yaml": {
		Encoder: writer.NewYAML(), Fn: func(metric Metric) any {
			return metric.ToYAML()
		},
	},
}
