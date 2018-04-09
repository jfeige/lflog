package lflog

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)
const(
	MaxQueue = 10000
)
const (
	debuglog = iota
	infolog
	warninglog
	errorlog
)

var (
	logs = make([]Log, 4)
)

type Log interface {
	write(source, message string)
}

func LoadConfig(configFile string) {
	//加载配置文件,xml
	err := readConfigFile(configFile)
	if err != nil {
		fmt.Printf("lflog read config file has error:%v", err)
		os.Exit(1)
	}
}

func Info(args0 interface{}, args ...interface{}) {
	var message, source string

	switch first := args0.(type) {
	case string:
		message = first
		if len(args) > 0 {
			message = fmt.Sprintf(first, args)
		}
	default:
		message = fmt.Sprintf(fmt.Sprint(args0)+strings.Repeat(" %v", len(args)),args)
	}

	pc, _, lineno, ok := runtime.Caller(1)
	if ok {
		source = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}
	logs[infolog].write(source, message)
}

func Error(args0 interface{}, args ...interface{}) {
	var message, source string

	switch first := args0.(type) {
	case string:
		message = first
		if len(args) > 0 {
			message = fmt.Sprintf(first, args)
		}
	default:
		message = fmt.Sprintf(fmt.Sprint(args0)+strings.Repeat(" %v", len(args)),args)
	}

	pc, _, lineno, ok := runtime.Caller(1)
	if ok {
		source = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	logs[errorlog].write(source, message)
}

func Debug(args0 interface{}, args ...interface{}) {
	var message, source string

	switch first := args0.(type) {
	case string:
		message = first
		if len(args) > 0 {
			message = fmt.Sprintf(first, args)
		}
	default:
		message = fmt.Sprintf(fmt.Sprint(args0)+strings.Repeat(" %v", len(args)),args)
	}

	pc, _, lineno, ok := runtime.Caller(1)
	if ok {
		source = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	logs[debuglog].write(source, message)
}

func Warn(args0 interface{}, args ...interface{}) {
	var message, source string

	switch first := args0.(type) {
	case string:
		message = first
		if len(args) > 0 {
			message = fmt.Sprintf(first, args)
		}
	default:
		message = fmt.Sprintf(fmt.Sprint(args0)+strings.Repeat(" %v", len(args)),args)
	}

	pc, _, lineno, ok := runtime.Caller(1)
	if ok {
		source = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	logs[warninglog].write(source, message)
}
