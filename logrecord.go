package lflog

import (
	"strings"
	"time"
	"os"
	"fmt"
)

type LogRecord struct {
	Enabled 	 bool
	loglevel     int
	Level        string //file|console	<!-- level is (:?FINEST|FINE|DEBUG|TRACE|INFO|WARNING|ERROR) -->
	Type         string
	OutType		 int	//1:file;2:console;3:file+console
	Created      time.Time
	Source       string
	Message      string
	Logfile      LogFile
	Opendate	 string
	f 			*os.File
	MessageQueue chan string
}

type LogFile struct {
	Filename string
	Format   string
	Rotate   bool
	Maxsize  string
	Maxlines string
	Daily    bool
}

func (l *LogRecord) write(source, message string) {
	//组装字符串，写入队列
	//[2018/04/03 15:07:44 CST] [INFO] [action.addDynamic:230] addDynamic--------------%!(EXTRA string=110, int=1583649)
	date := time.Now().Format("2006-01-02 15:04:05")
	src := l.Logfile.Format
	src = strings.Replace(src, "%D", date, 1)
	src = strings.Replace(src, "%L", l.Level, 1)
	src = strings.Replace(src, "%S", source, 1)
	src = strings.Replace(src, "%M", message, 1)

	l.MessageQueue <- src
}

func (l *LogRecord) close(){
	if l.f != nil{
		l.f.Close()
	}
}

func (l *LogRecord) isenable()bool{
	return l.Enabled
}

//info
func (l *LogRecord) writeLog() {
	for {
		select {
		case message := <-l.MessageQueue:
			//写日志文件
			l.checkLogDate()
			fmt.Fprintln(l.f,message)
			if l.OutType > 1{
				fmt.Println(message)
			}
		}
	}
}

//写日志时，判断是否已跨日，如果已跨日，则备份日志
func (l *LogRecord) checkLogDate(){
	if time.Now().Format("2006-01-02") != l.Opendate{
		var err error
		var num = 1
		//当前日期和程序建立时的日期不符，将日志备份为昨天的日期
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		//重命名带 1，2，3后缀的日志文件
		for ;err == nil && num < 999;num++{
			new_name := l.Logfile.Filename + fmt.Sprintf(".%s.%03d",yesterday,num)
			old_name := l.Logfile.Filename + fmt.Sprintf(".%d",num)
			_,err = os.Lstat(old_name)
			if err == nil{
				os.Rename(old_name,new_name)
			}else{
				break
			}
		}
		//关闭当前
		l.f.Close()
		_,err = os.Lstat(l.Logfile.Filename)
		if err == nil{
			new_name := l.Logfile.Filename + fmt.Sprintf(".%s.%03d",yesterday,num)
			os.Rename(l.Logfile.Filename,new_name)

		}
		l.f,err = os.OpenFile(l.Logfile.Filename,os.O_APPEND|os.O_WRONLY|os.O_CREATE,0666)
	}
}

//日志文件创建(应用启动时调用)
func createLogFile(filename string){
	_,err := os.Lstat(filename)
	if err != nil{
		//文件不存在
		os.Create(filename)
		return
	}
	num := 1
	fname := ""
	for ;err == nil && num <= 999;num++{
		fname = filename + fmt.Sprintf(".%d",num)
		_,err = os.Lstat(fname)
		if err != nil{
			//更改日志文件名称
			os.Rename(filename,fname)
			os.Create(filename)
			break
		}
	}
	if err == nil{
		//当天已到999个日志，抛出异常
	}
}