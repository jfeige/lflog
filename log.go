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
	console = Console{}				//控制台输出
)

type Log interface {
	isenable()bool
	write(source, message string)string
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

func info(level int,tag string,args0 interface{}, args ...interface{}){
	var log Log
	var ret string
	message := handleMessage(args0, args...)
	source := handleLineNb()

	for ;level < len(logs);level++{
		log = logs[level]
		if log.isenable(){
			ret = log.write(source, message)
		}
	}
	if console.Enable{
		if _,ok := console.Level[tag];ok{
			fmt.Println(ret)
		}
	}
}

func Debug(args0 interface{}, args ...interface{}) {
	info(debuglog,"debug",args0,args...)
}

func Info(args0 interface{}, args ...interface{}) {
	info(infolog,"info",args0,args...)
}

func Warn(args0 interface{}, args ...interface{}) {
	info(warninglog,"warning",args0,args...)
}

func Error(args0 interface{}, args ...interface{}) {
	info(errorlog,"error",args0,args...)
}

//处理日志内容
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

//获取行号
func handleLineNb() string{
	var source string

	pc, _, lineno, ok := runtime.Caller(2)
	if ok {
		source = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	return source
}

