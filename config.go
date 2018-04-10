package lflog

import (
	"encoding/xml"
	"io/ioutil"
	"time"
	"os"
	"strings"
)

type Logging struct {
	Filters []Filter `xml:"filter"`
}
type Filter struct {
	Enabled    bool      `xml:"enabled,attr"`
	Tag       string     `xml:"tag"`
	Type      string     `xml:"type"`
	Level     string     `xml:"level"`
	Propertys []Property `xml:"property"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

func readConfigFile(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	var logging = Logging{}
	err = xml.Unmarshal(data, &logging)
	if err != nil {
		return err
	}
	filters := logging.Filters
	for _, filter := range filters {
		lr := new(LogRecord)

		lr.Enabled = filter.Enabled
		var outtype int
		tps := strings.Split(filter.Type,"|")
		for _,tp := range tps{
			switch tp {
			case "file":
				outtype += 1
			case "console":
				outtype += 2
			}
		}
		lr.OutType = outtype
		lr.Opendate = time.Now().Format("2006-01-02")
		propertys := filter.Propertys
		logFile := LogFile{}
		for _, property := range propertys {
			switch property.Name {
			case "filename":
				logFile.Filename = property.Value
			case "format":
				logFile.Format = property.Value
			}
		}
		lr.Logfile = logFile
		err := createLogFile(logFile.Filename)
		if err != nil{
			//创建日志文件失败
			return err
		}

		lr.f,err = os.OpenFile(logFile.Filename,os.O_APPEND|os.O_WRONLY,0666)
		if err != nil{
			return err
		}
		lr.MessageQueue = make(chan string, MaxQueue)
		var loglevel = -1
		switch filter.Level {
		case "DEBUG":
			loglevel = debuglog
		case "INFO":
			loglevel = infolog
		case "WARNING":
			loglevel = warninglog
		case "ERROR":
			loglevel = errorlog
		default:
			continue
		}
		if loglevel >= 0{
			lr.Level = filter.Level
			logs[loglevel] = lr
			go lr.writeLog()
		}
	}
	return nil
}
