package reader

import (
	"encoding/json"
	"os"

	"go-challenge/internal/args"
)

type JSON struct {
	config  *args.Config
	decoder *json.Decoder
}

func NewJSON(config *args.Config) *JSON {
	j := &JSON{config: config}
	return j
}

func (j JSON) File(file *os.File) (Decoder, error) {
	newDecoder := NewJSON(j.config)
	newDecoder.decoder = json.NewDecoder(file)
	t, err := newDecoder.decoder.Token()
	if err != nil {
		return nil, err
	}
	if t != json.Delim('[') {
		return nil, err
	}
	return newDecoder, nil
}

func (j JSON) ParseAll(sendParsedMetric func(item Metric)) {
	for j.decoder.More() {
		var m Metric
		if e := j.decoder.Decode(&m); e == nil {
			sendParsedMetric(m)
		}
	}
}
