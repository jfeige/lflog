package lflog

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Logging struct {
	Filters []Filter `xml:"filter"`
}
type Filter struct {
	Enabled  bool   `xml:"enabled,attr"`
	Tag      string `xml:"tag"`
	Level    string `xml:"level"`
	Filename string `xml:"filename"`
	Format   string `xml:"format"`
}

type Console struct {
	Enable bool
	Level  map[string]string
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
		if filter.Tag == "CONSOLE" {
			console.Enable = filter.Enabled
			levels := strings.Split(filter.Level, "|")
			mlevel := make(map[string]string)
			for _, level := range levels {
				mlevel[level] = level
			}
			console.Level = mlevel
			continue
		}
		lr := new(LogRecord)
		lr.Enabled = filter.Enabled
		lr.Tag = filter.Tag
		lr.Opendate = time.Now().Format("2006-01-02")
		switch lr.Tag {
		case "DEBUG", "INFO", "WARNING", "ERROR":
			lr.Filename = filter.Filename
			lr.Format = filter.Format
		default:
			continue
		}
		createLogFile(lr.Filename)

		lr.f, err = os.OpenFile(lr.Filename, os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		lr.MessageQueue = make(chan string, MaxQueue)
		var loglevel = -1
		switch lr.Tag {
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
		if loglevel >= 0 {
			logs[loglevel] = lr
			go lr.writeLog()
		}
	}
	return nil
}
