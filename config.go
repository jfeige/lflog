package lflog

import (
	"encoding/xml"
	"io/ioutil"
	"strconv"
	"time"
	"os"
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
		if !filter.Enabled {
			continue
		}
		lr := LogRecord{}
		lr.Tag = filter.Tag
		lr.Type = filter.Type
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
		switch filter.Level {
		case "INFO":
			lr.Level = "INFO"
			logs[infolog] = lr
			go lr.writeLog()
		case "DEBUG":
			lr.Level = "DEBUG"
			logs[debuglog] = lr
			go lr.writeLog()
		case "ERROR":
			lr.Level = "ERROR"
			logs[errorlog] = lr
			go lr.writeLog()
		case "WARNING":
			lr.Level = "WARNING"
			logs[warninglog] = lr
			go lr.writeLog()
		default:
			continue
		}
	}
	return nil
}
