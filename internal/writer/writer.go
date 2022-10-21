package writer

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"go-challenge/internal/args"
)

type Writer[enc Encoder] struct {
	serializer enc
	w          *bufio.Writer
	f          *os.File
	config     *args.Config
	firstItem  bool
}

type Encoder interface {
	OnCreateFile() []byte
	OnCloseFile() []byte
	OnNext() []byte
	ToBytes(entity any) ([]byte, error)
}

func NewWriter[enc Encoder](config *args.Config, e enc, date ...time.Time) Writer[enc] {
	writer := Writer[enc]{
		serializer: e,
		config:     config,
		firstItem:  true,
	}
	writer.CreateFile(date...)
	return writer
}

func (wr *Writer[s]) Write(x any) error {
	b, err := wr.serializer.ToBytes(x)
	if err != nil {
		panic(err)
	}
	if !wr.firstItem {
		b = append(wr.serializer.OnNext(), b...)
	}
	_, err = wr.w.Write(b)
	wr.firstItem = false
	return err
}

func (wr *Writer[s]) CreateFile(date ...time.Time) {
	err := wr.Close()
	if err != nil {
		panic(err)
	}
	filePath := wr.config.Directory
	name := fmt.Sprintf("%s.%s", wr.config.OutPutFileName, wr.config.OutPutFileType)
	if len(date) > 0 && wr.config.GenerateTestFiles {
		filePath = path.Join(wr.config.Directory, wr.config.Type)
		name = fmt.Sprintf("%s/%s.%s", filePath, strings.ToLower(date[0].Format("02-Jan")), wr.config.Type)
	}
	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	wr.f, err = os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	wr.w = bufio.NewWriter(wr.f)
	if len(date) > 0 {
		fmt.Printf("%s\n", name)
	}
	_, err = wr.w.Write(wr.serializer.OnCreateFile())
	if err != nil {
		panic(err)
	}
	wr.firstItem = true
}

func (wr *Writer[s]) Close() error {
	if wr.f == nil {
		return nil
	}
	_, err := wr.w.Write(wr.serializer.OnCloseFile())
	if err != nil {
		return err
	}
	err = wr.w.Flush()
	if err != nil {
		return err
	}
	err = wr.f.Close()
	if err != nil {
		return err
	}
	return nil
}
