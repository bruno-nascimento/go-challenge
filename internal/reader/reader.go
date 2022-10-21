package reader

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"go-challenge/internal/args"
)

type Decoder interface {
	File(file *os.File) (Decoder, error)
	ParseAll(sendParsedMetric func(item Metric))
}

type Metric struct {
	Timestamp time.Time `json:"timestamp"`
	LevelName string    `json:"level_name"`
	Value     int       `json:"value"`
}

type Reader struct {
	config        *args.Config
	fileChannel   chan string
	wg            sync.WaitGroup
	M             map[string]int64
	resultChannel chan Metric
	decoder       Decoder
}

func NewReader(config *args.Config, decoder Decoder) *Reader {
	r := &Reader{config: config, fileChannel: make(chan string), M: make(map[string]int64), resultChannel: make(chan Metric, 1_000_000), decoder: decoder}
	r.startWorkers()
	return r
}

func (r *Reader) Read() error {
	fileNames, err := r.filterFileNames()
	if err != nil {
		return err
	}

	for _, name := range fileNames {
		r.wg.Add(1)
		r.fileChannel <- name
	}

	r.wg.Wait()
	close(r.fileChannel)
	close(r.resultChannel)

	return nil
}

func (r *Reader) startWorkers() {
	for i := range make([]int, runtime.NumCPU()) {
		go func(j int) {
			for name := range r.fileChannel {
				file, er := os.Open(path.Join(r.config.Directory, name))
				if er != nil {
					panic(er)
				}
				dec, err := r.decoder.File(file)
				if err != nil {
					panic(err)
				}
				dec.ParseAll(r.receiveDecodedMetric)
				r.wg.Done()
			}
		}(i)
	}

	go func() {
		for metric := range r.resultChannel {
			r.M[metric.LevelName] += int64(metric.Value)
			r.wg.Done()
		}
	}()
}

func (r *Reader) receiveDecodedMetric(m Metric) {
	if (m.Timestamp.After(r.config.StartTime) || m.Timestamp.Equal(r.config.StartTime)) && m.Timestamp.Before(r.config.EndTime) {
		r.wg.Add(1)
		r.resultChannel <- m
	}
}

func (r *Reader) filterFileNames() ([]string, error) {
	f, err := os.Open(r.config.Directory)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	files, err := f.Readdir(0)
	files = r.checkExtensions(files)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if len(files) < 1 {
		return nil, fmt.Errorf("there were no files with the extension '%s' found in the specified directory (%s)", r.config.Type, r.config.Directory)
	}
	var fileNames []string
	isInRange := r.fileInRangeCompareFn()
	for _, v := range files {
		if isInRange(v.Name()) {
			fileNames = append(fileNames, v.Name())
		}
	}
	return fileNames, nil
}

func (r *Reader) checkExtensions(files []os.FileInfo) []os.FileInfo {
	var fs []os.FileInfo
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), r.config.Type) {
			fs = append(fs, f)
		}
	}
	return fs
}

func (r *Reader) fileInRangeCompareFn() func(fileName string) bool {
	start := r.toMidnight(r.config.StartTime)
	end := r.toMidnight(r.config.EndTime)
	return func(fileName string) bool {
		// requirement: Metric file and the time range parameter are in the same (current) year, starting from January 1st.
		// that is why I am adding the current year to the name of the file here
		f := fmt.Sprintf("%s-%d", fileName[:strings.Index(fileName, ".")], time.Now().Year())
		parse, err := time.Parse("02-Jan-2006", f)
		if err != nil {
			fmt.Printf("WARNING - ignoring invalid named file: %s\n", fileName)
			return false
		}
		return (parse.After(start) || parse.Equal(start)) && (parse.Before(end) || parse.Equal(end))
	}
}

func (r *Reader) toMidnight(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}
