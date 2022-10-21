package writer

import "fmt"

type YAML struct{}

func NewYAML() YAML {
	return YAML{}
}

func (c YAML) OnCreateFile() []byte {
	return []byte("---\n")
}

func (c YAML) OnCloseFile() []byte {
	return []byte{}
}

func (c YAML) OnNext() []byte {
	return []byte{}
}

func (c YAML) ToBytes(entity any) ([]byte, error) {
	return []byte(fmt.Sprint(entity)), nil
}
