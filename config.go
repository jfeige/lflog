package lflog

import (
	"encoding/xml"
	"io/ioutil"
	"strconv"
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
	opendate := time.Now().Format("2006-01-02")
	filters := logging.Filters
	for _, filter := range filters {
		lr := new(LogRecord)
		lr.Type = filter.Type

		lr.Enabled = filter.Enabled
		var outtype int
		tps := strings.Split(lr.Type,"|")
		for _,tp := range tps{
			switch tp {
			case "file":
				outtype += 1
			case "console":
				outtype += 2
			}
		}
		lr.OutType = outtype
		lr.Opendate = opendate
		propertys := filter.Propertys
		logFile := LogFile{}
		for _, property := range propertys {
			switch property.Name {
			case "filename":
				logFile.Filename = property.Value
			case "format":
				logFile.Format = property.Value
			case "rotate":
				logFile.Rotate, _ = strconv.ParseBool(property.Value)
			case "maxsize":
				logFile.Maxsize = property.Value
			case "maxlines":
				logFile.Maxlines = property.Value
			case "daily":
				logFile.Daily, _ = strconv.ParseBool(property.Value)
			}
		}
		lr.Logfile = logFile
		createLogFile(logFile.Filename)

		lr.f,_ = os.OpenFile(logFile.Filename,os.O_APPEND|os.O_WRONLY,0666)
		lr.MessageQueue = make(chan string, MaxQueue)
		var loglevel = -1
		switch filter.Level {
		case "DEBUG":
			loglevel = 0
		case "INFO":
			loglevel = 1
		case "WARNING":
			loglevel = 2
		case "ERROR":
			loglevel = 3
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
