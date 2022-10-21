package writer

import "encoding/json"

type JSON struct{}

func NewJSON() JSON {
	return JSON{}
}

func (J JSON) OnCreateFile() []byte {
	return []byte("[\n")
}

func (J JSON) OnCloseFile() []byte {
	return []byte("\n]")
}

func (J JSON) OnNext() []byte {
	return []byte(",\n")
}

func (J JSON) ToBytes(entity any) ([]byte, error) {
	return json.MarshalIndent(entity, "", "  ")
}
