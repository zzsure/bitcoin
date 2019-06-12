package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"gitlab.azbit.cn/web/bitcoin/conf"
	goLog "log"
	"os"
	"runtime"
)

type RequestID string
type Event string

type jsonLog struct{}

type jsonLogData struct {
	RequestID RequestID `json:"requestId"`
	Module    string    `json:"module"`
	Event     Event     `json:"event"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Idc       string    `json:"idc"`
}

func newJsonLog() *jsonLog {
	j := new(jsonLog)
	return j
}

func (j *jsonLog) Log(level logging.Level, depth int, record *logging.Record) error {
	var jl jsonLogData
	var ok bool
	jl.Module = record.Module
	jl.RequestID, ok = record.Args[0].(RequestID)
	if ok {
		record.Args = record.Args[1:]
	}
	jl.Event, ok = record.Args[0].(Event)
	if ok {
		record.Args = record.Args[1:]
	}
	jl.Message = record.Formatted(depth + 1)
	jl.Level = level.String()
	jl.Idc = os.Getenv("AIBEE_KUBERNETES_IDC")
	jlStr, _ := json.Marshal(jl)
	os.Stdout.Write(append(jlStr, '\n'))

	return nil
}

type consoleLog struct {
	log *goLog.Logger
}

func newConsoleLog() *consoleLog {
	c := new(consoleLog)
	c.log = goLog.New(os.Stdout, "[AzbitServer] ", 0)
	return c
}

func (c *consoleLog) Log(level logging.Level, depth int, record *logging.Record) error {
	module := record.Module
	message := record.Formatted(depth + 1)
	c.log.SetPrefix(fmt.Sprintf("[%s] ", module))
	c.log.Println(message)
	return nil
}

func Init() {
	if conf.Config.Log.Type == "json" {
		jsonLogFormat := logging.MustStringFormatter(`%{shortfile} %{shortfunc} %{message}`)
		jL := logging.NewBackendFormatter(newJsonLog(), jsonLogFormat)
		jsonLogLeveled := logging.AddModuleLevel(jL)
		jsonLogLeveled.SetLevel(logging.Level(conf.Config.Log.Level), "")
		logging.SetBackend(jsonLogLeveled)
	} else {
		format := logging.MustStringFormatter(
			`%{color}%{time:2006-01-02 15:04:05.000} %{shortfile} %{shortfunc} ▶ %{level:.4s} %{color:reset} %{message}`,
		)
		consoleLogFormatter := logging.NewBackendFormatter(newConsoleLog(), format)
		consoleLogLeveled := logging.AddModuleLevel(consoleLogFormatter)
		consoleLogLeveled.SetLevel(logging.Level(conf.Config.Log.Level), "")
		logging.SetBackend(consoleLogLeveled)
	}
}

func WrapError(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return errors.New(fmt.Sprintf("%s:%d %s", file, line, err.Error()))
}

