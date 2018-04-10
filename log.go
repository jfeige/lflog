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
	close()
}

func LoadConfig(configFile string) {
	//加载配置文件,xml
	err := readConfigFile(configFile)
	if err != nil {
		fmt.Printf("lflog read config file has error:%v", err)
		os.Exit(1)
	}
}

func Close(){
	for _,f := range logs{
		f.close()
	}
}

func Debug(args0 interface{}, args ...interface{}) {

	message := handleMessage(args0, args...)
	source := handleLineNb()

	for level := debuglog;level < len(logs);level++{
		logs[level].write(source, message)
	}
}

func Info(args0 interface{}, args ...interface{}) {

	message := handleMessage(args0, args...)
	source := handleLineNb()

	for level := infolog;level < len(logs);level++{
		logs[level].write(source, message)
	}
}

func Warn(args0 interface{}, args ...interface{}) {

	message := handleMessage(args0, args...)
	source := handleLineNb()

	for level := warninglog;level < len(logs);level++{
		logs[level].write(source, message)
	}
}

func Error(args0 interface{}, args ...interface{}) {

	message := handleMessage(args0, args...)
	source := handleLineNb()

	for level := errorlog;level < len(logs);level++{
		logs[level].write(source, message)
	}
}


func handleMessage(args0 interface{},args ...interface{})string{
	var message string

	switch first := args0.(type) {
	case string:
		message = first
		if len(args) > 0 {
			message = fmt.Sprintf(first, args)
		}
	default:
		message = fmt.Sprintf(fmt.Sprint(args0)+strings.Repeat(" %v", len(args)),args)
	}

	return message
}


func handleLineNb() string{
	var source string

	pc, _, lineno, ok := runtime.Caller(2)
	if ok {
		source = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	return source
}