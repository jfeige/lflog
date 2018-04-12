package lflog

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)
const(
	MaxQueue = 10000
	MaxLog = 999
)
const (
	errorlog = iota
	warninglog
	infolog
	debuglog
)

var (
	logs = make([]Log, 4)
)

type Log interface {
	isenable()bool
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

func info(level int,args0 interface{}, args ...interface{}){
	var log Log
	for ;level < len(logs);level++{
		log = logs[level]
		if log.isenable(){
			message := handleMessage(args0, args...)
			source := handleLineNb()
			log.write(source, message)
		}
	}
}

func Debug(args0 interface{}, args ...interface{}) {
	info(debuglog,args0,args...)
}

func Info(args0 interface{}, args ...interface{}) {
	info(infolog,args0,args...)
}

func Warn(args0 interface{}, args ...interface{}) {
	info(warninglog,args0,args...)
}

func Error(args0 interface{}, args ...interface{}) {
	info(errorlog,args0,args...)
}


func handleMessage(args0 interface{},args ...interface{})string{
	var message string

	switch first := args0.(type) {
	case string:
		message = first
		if len(args) > 0 {
			message = fmt.Sprintf(first, args...)
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