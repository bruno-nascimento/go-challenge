package args

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	typeValues        = map[string]struct{}{"json": {}, "csv": {}}
	outputTypesValues = map[string]struct{}{"json": {}, "yaml": {}}
)

type ValidationError struct {
	msgs []string
}

func (v *ValidationError) append(msg string) {
	v.msgs = append(v.msgs, msg)
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("Invalid arguments:\n%s", strings.Join(v.msgs, "\n"))
}

func (v *ValidationError) check() *ValidationError {
	if v.msgs == nil {
		return nil
	}
	return v
}

func validate(cfg *Config) *ValidationError {
	err := &ValidationError{}
	if required(cfg.Directory, "directory", err) {
		if _, er := os.Stat(cfg.Directory); er != nil {
			err.append("directory: " + er.Error())
		}
	}
	if required(cfg.Type, "type", err) {
		permittedValue(cfg.Type, "type", typeValues, err)
	}
	if required(cfg.strStartTime, "startTime", err) {
		parseDate(&cfg.StartTime, cfg.strStartTime, "startTime", err)
	}
	if required(cfg.strEndTime, "endTime", err) {
		parseDate(&cfg.EndTime, cfg.strEndTime, "endTime", err)
	}
	permittedValue(cfg.OutPutFileType, "outputFileType", outputTypesValues, err)
	startBeforeFinish(cfg.StartTime, cfg.EndTime, err)
	return err.check()
}

func required(field, name string, err *ValidationError) bool {
	empty := field == ""
	if empty {
		err.append(fmt.Sprintf("%s: argument is required", name))
	}
	return !empty
}

func permittedValue(value, name string, permitted map[string]struct{}, err *ValidationError) bool {
	_, ok := permitted[value]
	if !ok {
		err.append(fmt.Sprintf("%s: '%s' is not a valid value; permitted values are '%s'", name, value, enumeratePermittedValues(permitted)))
	}
	return ok
}

func enumeratePermittedValues(permitted map[string]struct{}) string {
	var keys []string
	for k := range permitted {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func parseDate(date *time.Time, strDate, name string, err *ValidationError) bool {
	parsed, er := time.Parse(time.RFC3339, strDate)
	if er != nil {
		err.append(fmt.Sprintf("%s: invalid date format - %s", name, er.Error()))
		return false
	}
	*date = parsed
	return true
}

func startBeforeFinish(start, end time.Time, err *ValidationError) bool {
	valid := start.Before(end)
	if !valid {
		err.append("invalid dates: startTime should be before endTime")
	}
	return valid
}
