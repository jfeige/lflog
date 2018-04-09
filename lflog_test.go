package lflog

import (
	"testing"
	"fmt"
	"time"
)

func TestLog(t *testing.T) {


	LoadConfig("config.xml")
	Info(fmt.Sprintf("大家好，才是真的好%s.","小明"))
	Debug(fmt.Sprintf("这是Debug:%s","大家好"))
	time.Sleep(3*time.Second)
	Info("finish...")
}
