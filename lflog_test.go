package lflog

import (
	"testing"
	"fmt"
	"time"
)

func TestLog(t *testing.T) {


	LoadConfig("config.xml")
	defer Close()
	Info("当前时间:%s",time.Now().Format("2006-01-02 15:04:05"))
	Debug(fmt.Sprintf("uid:%d",12345))
	time.Sleep(3*time.Second)
	Info("finish...")
}
