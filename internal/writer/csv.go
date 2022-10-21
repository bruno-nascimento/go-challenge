package writer

import "fmt"

type CSV struct{}

func NewCSV() CSV {
	return CSV{}
}

func (c CSV) OnCreateFile() []byte {
	return []byte("timestamp,level_name,value\n")
}

func (c CSV) OnCloseFile() []byte {
	return []byte{}
}

func (c CSV) OnNext() []byte {
	return []byte{}
}

func (c CSV) ToBytes(entity any) ([]byte, error) {
	return []byte(fmt.Sprint(entity)), nil
}
